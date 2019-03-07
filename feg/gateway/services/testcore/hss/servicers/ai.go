/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"bytes"
	"errors"
	"fmt"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/s6a_proxy/servicers"
	"magma/feg/gateway/services/testcore/hss/crypto"
	"magma/feg/gateway/services/testcore/hss/storage"
	"magma/lte/cloud/go/protos"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/golang/glog"
)

const (
	// indBits is the number of bits reserved for IND (one of the two parts of SQN).
	// See 3GPP TS 33.102 Appendix C.1.1.1 and C.3.
	indBits = 5

	// indMask is a bit mask where a bit is 1 if and only if it is a part of ind.
	indMask = (1 << indBits) - 1

	// seqMask is a bit mask where a bit is 1 if and only if it is a part of seq.
	seqMask = (1 << 48) - 1 - indMask

	// lteResyncInfoBytes is the expected size of the lte resync info in bytes.
	// The first 16 bytes store RAND and the next 14 bytes store AUTS.
	lteResyncInfoBytes = crypto.RandChallengeBytes + crypto.ExpectedAutsBytes

	// maxSeqDelta is the maximum allowed increase to SQN.
	// eg. if x was the last accepted SQN, then the next SQN must
	// be greater than x and less than (x + maxSeqDelta) to be accepted.
	// See 3GPP TS 33.102 Appendix C.2.1.
	maxSeqDelta = 1 << 28
)

// NewAIA outputs a authentication information answer (AIA) to reply to an
// authentication information request (AIR) message.
func NewAIA(srv *HomeSubscriberServer, msg *diam.Message) (*diam.Message, error) {
	if err := ValidateAIR(msg); err != nil {
		return msg.Answer(diam.MissingAVP), err
	}

	var air servicers.AIR
	if err := msg.Unmarshal(&air); err != nil {
		return msg.Answer(diam.UnableToComply), fmt.Errorf("AIR Unmarshal failed for message: %v failed: %v", msg, err)
	}

	subscriber, err := srv.store.GetSubscriberData(air.UserName)
	if err != nil {
		if _, ok := err.(storage.UnknownSubscriberError); ok {
			return ConstructFailureAnswer(msg, air.SessionID, srv.Config.Server, uint32(fegprotos.ErrorCode_USER_UNKNOWN)), err
		}
		return ConstructFailureAnswer(msg, air.SessionID, srv.Config.Server, uint32(fegprotos.ErrorCode_AUTHENTICATION_DATA_UNAVAILABLE)), err
	}

	err = srv.ResyncLteAuthSeq(subscriber, air.RequestedEUTRANAuthInfo.ResyncInfo.Serialize())
	if err != nil {
		return ConvertAuthErrorToFailureMessage(err, msg, air.SessionID, srv.Config.Server), err
	}

	const plmnOffsetBytes = 1
	plmn := air.VisitedPLMNID.Serialize()[plmnOffsetBytes:]

	var vectors = make([]*crypto.EutranVector, 0, air.RequestedEUTRANAuthInfo.NumVectors)
	for i := datatype.Unsigned32(0); i < air.RequestedEUTRANAuthInfo.NumVectors; i++ {
		vector, err := srv.GenerateLteAuthVector(subscriber, plmn)
		if err != nil {
			// If we have already generated an auth vector successfully, then we can
			// return it. Otherwise, we must signal an error.
			// See 3GPP TS 29.272 section 5.2.3.1.3.
			if i == 0 {
				return ConvertAuthErrorToFailureMessage(err, msg, air.SessionID, srv.Config.Server), err
			}
			glog.Errorf("failed to generate lte auth vector: %v", err)
			break
		}
		vectors = append(vectors, vector)
	}

	return srv.NewSuccessfulAIA(msg, air.SessionID, vectors), nil
}

// NewSuccessfulAIA outputs a successful authentication information answer (AIA) to reply to an
// authentication information request (AIR) message. It populates AIA with all of the mandatory fields
// and adds the authentication vectors.
func (srv *HomeSubscriberServer) NewSuccessfulAIA(msg *diam.Message, sessionID datatype.UTF8String, vectors []*crypto.EutranVector) *diam.Message {
	answer := ConstructSuccessAnswer(msg, sessionID, srv.Config.Server)
	for _, vector := range vectors {
		answer.NewAVP(avp.AuthenticationInfo, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.EUTRANVector, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
					AVP: []*diam.AVP{
						diam.NewAVP(avp.RAND, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(vector.Rand[:])),
						diam.NewAVP(avp.XRES, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(vector.Xres[:])),
						diam.NewAVP(avp.AUTN, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(vector.Autn[:])),
						diam.NewAVP(avp.KASME, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(vector.Kasme[:])),
					},
				}),
			},
		})
	}
	return answer
}

// ResyncLteAuthSeq validates a re-synchronization request and computes the SEQ
// from the AUTS sent by U-SIM.
// See 3GPP TS 33.102 section 6.3.5.
func (srv *HomeSubscriberServer) ResyncLteAuthSeq(subscriber *protos.SubscriberData, resyncInfo []byte) error {
	if AllZero(resyncInfo) {
		return nil
	}
	if len(resyncInfo) != lteResyncInfoBytes {
		return NewAuthRejectedError(fmt.Sprintf("resync info incorrect length. expected %v bytes, but got %v bytes", lteResyncInfoBytes, len(resyncInfo)))
	}
	lte := subscriber.Lte
	if err := ValidateLteSubscription(lte); err != nil {
		return NewAuthRejectedError(err.Error())
	}

	// Use dummy AMF for re-synchronization. See 3GPP TS 33.102 section 6.3.3.
	milenage, err := crypto.NewMilenageCipher(make([]byte, crypto.ExpectedAmfBytes))
	if err != nil {
		return NewAuthDataUnavailableError(err.Error())
	}
	rand := resyncInfo[:crypto.RandChallengeBytes]
	auts := resyncInfo[crypto.RandChallengeBytes:]
	opc, err := srv.GetOrGenerateOpc(lte)
	if err != nil {
		return err
	}
	sqnMs, macS, err := milenage.GenerateResync(auts, subscriber.Lte.AuthKey, opc, rand)
	if err != nil {
		return NewAuthDataUnavailableError(err.Error())
	}
	if !bytes.Equal(macS[:], auts[crypto.ExpectedAutsBytes-len(macS):]) {
		return NewAuthRejectedError("Invalid resync authentication code")
	}

	return srv.SetNextLteAuthSqnAfterResync(subscriber, sqnMs)
}

// SetNextLteAuthSqnAfterResync tries to set the subscribers State.LteAuthNextSeq field
// to the next sequence number after `seq`.
// See 3GPP TS 33.102 Appendix C.3.
func (srv *HomeSubscriberServer) SetNextLteAuthSqnAfterResync(subscriber *protos.SubscriberData, sqn uint64) error {
	seq, _ := SplitSqn(sqn)
	currentSeq := subscriber.State.LteAuthNextSeq - 1
	if seq < currentSeq {
		seqDelta := currentSeq - seq
		if seqDelta <= maxSeqDelta {
			// This error indicates that the last sequence number should have been
			// accepted by the USIM but wasn't (this should never happen).
			return NewAuthRejectedError(fmt.Sprintf("Re-sync delta in range but UE rejected auth: %d", seqDelta))
		}
	}
	return srv.SetNextLteAuthSeq(subscriber, seq+1)
}

// GenerateLteAuthVector returns the lte auth vector for the subscriber.
// Inputs:
//   subscriber: The subscriber data for the subscriber we want to generate auth vectors for
//   plmn: 24 bit network identifier
//   index: the index of the current vector being generated
func (srv *HomeSubscriberServer) GenerateLteAuthVector(subscriber *protos.SubscriberData, plmn []byte) (*crypto.EutranVector, error) {
	lte := subscriber.Lte
	if err := ValidateLteSubscription(lte); err != nil {
		return nil, NewAuthRejectedError(err.Error())
	}
	if subscriber.State == nil {
		return nil, NewAuthRejectedError("Subscriber data missing subscriber state")
	}

	opc, err := srv.GetOrGenerateOpc(lte)
	if err != nil {
		return nil, err
	}
	err = srv.IncreaseSQN(subscriber)
	if err != nil {
		return nil, err
	}
	sqn := SeqToSqn(subscriber.State.LteAuthNextSeq, srv.AuthSqnInd)
	vector, err := srv.Milenage.GenerateEutranVector(lte.AuthKey, opc, sqn, plmn)
	if err != nil {
		return vector, NewAuthRejectedError(err.Error())
	}
	return vector, err
}

// GetOrGenerateOpc returns lte.AuthOpc and generates if it isn't stored in the proto
func (srv *HomeSubscriberServer) GetOrGenerateOpc(lte *protos.LTESubscription) ([]byte, error) {
	if lte == nil || len(lte.AuthOpc) == 0 {
		opc, err := crypto.GenerateOpc(lte.AuthKey, srv.Config.LteAuthOp)
		if err != nil {
			err = NewAuthDataUnavailableError(err.Error())
		}
		return opc[:], err
	}
	return lte.AuthOpc, nil
}

// SeqToSqn computes the 48 bit SQN given a seq given the formula defined in
// 3GPP TS 33.102 Annex C.3.2. The length of IND is 5 bits.
// SQN = SEQ || IND
// Inputs:
//    seq: the sequence number
//    index: the index of the current vector being generated
// Output: The 48 bit SQN
func SeqToSqn(seq, index uint64) uint64 {
	return (seq << indBits & seqMask) + (index & indMask)
}

// SplitSqn computes the SEQ and IND given a 48 bit SQN using the formula defined in
// 3GPP TS 33.102 Annex C.3.2. The length of IND is 5 bits.
// SQN = SEQ || IND
// Inputs:
//    seq: the 48 bit SQN
// Outputs: SEQ and IND
func SplitSqn(sqn uint64) (uint64, uint64) {
	return sqn >> indBits, sqn & indMask
}

// IncreaseSQN increases both components of SQN (SEQ and IND) by one.
func (srv *HomeSubscriberServer) IncreaseSQN(subscriber *protos.SubscriberData) error {
	if subscriber.State == nil {
		subscriber.State = &protos.SubscriberState{}
	}
	return srv.SetNextLteAuthSeq(subscriber, subscriber.State.LteAuthNextSeq+1)
}

// SetNextLteAuthSeq sets the State.LteAuthNextSeq field of the subscriber data
// and updates the storage to reflect the change.
func (srv *HomeSubscriberServer) SetNextLteAuthSeq(subscriber *protos.SubscriberData, nextLteAuthSeq uint64) error {
	if subscriber.State == nil {
		subscriber.State = &protos.SubscriberState{}
	}

	subscriber.State.LteAuthNextSeq = nextLteAuthSeq
	err := srv.store.UpdateSubscriber(subscriber)
	if err != nil {
		return NewAuthDataUnavailableError(err.Error())
	}
	return nil
}

// ValidateAIR returns an error if the message is missing any mandatory AVPs.
// Mandatory AVPs are specified in 3GPP TS 29.272 Table 5.2.3.1.1/1
func ValidateAIR(msg *diam.Message) error {
	_, err := msg.FindAVP(avp.UserName, 0)
	if err != nil {
		return errors.New("Missing IMSI in message")
	}
	_, err = msg.FindAVP(avp.VisitedPLMNID, diameter.Vendor3GPP)
	if err != nil {
		return errors.New("Missing Visited PLMN ID in message")
	}
	_, err = msg.FindAVP(avp.RequestedEUTRANAuthenticationInfo, diameter.Vendor3GPP)
	if err != nil {
		return errors.New("Missing requested E-UTRAN authentication info in message")
	}
	_, err = msg.FindAVP(avp.SessionID, 0)
	if err != nil {
		return errors.New("Missing SessionID in message")
	}
	return nil
}

// ValidateLteSubscription returns an error if and only if the lte proto is not
// configured up to use the milenage authentication algorithm.
func ValidateLteSubscription(lte *protos.LTESubscription) error {
	if lte == nil {
		return fmt.Errorf("Subscriber data missing LTE subscription")
	}
	if lte.State != protos.LTESubscription_ACTIVE {
		return fmt.Errorf("LTE Service not active")
	}
	if lte.AuthAlgo != protos.LTESubscription_MILENAGE {
		return fmt.Errorf("Unsupported crypto algorithm: %v", lte.AuthAlgo)
	}
	return nil
}
