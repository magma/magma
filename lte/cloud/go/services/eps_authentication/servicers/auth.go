/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"bytes"
	"fmt"

	"magma/lte/cloud/go/crypto"
	"magma/lte/cloud/go/protos"

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

// GenerateLteAuthVectors generates at most `numVectors` lte auth vectors.
// Inputs:
//   numVectors: The maximum number of vectors to generate
//   milenage: The cipher to use to generate the vector
//   subscriber: The subscriber data for the subscriber we want to generate auth vectors for
//   plmn: 24 bit network identifier
//   authSqnInd: the IND of the current vector being generated
// Returns: The E-UTRAN vectors and the next value to set the subscriber's LteAuthNextSeq to (or an error).
func GenerateLteAuthVectors(numVectors uint32, milenage *crypto.MilenageCipher, subscriber *protos.SubscriberData, plmn, lteAuthOp []byte, authSqnInd uint64) ([]*crypto.EutranVector, uint64, error) {
	var vectors = make([]*crypto.EutranVector, 0, numVectors)
	lteAuthNextSeq := subscriber.GetState().GetLteAuthNextSeq()
	for i := uint32(0); i < numVectors; i++ {
		vector, nextSeq, err := GenerateLteAuthVector(milenage, subscriber, plmn, lteAuthOp, authSqnInd)
		lteAuthNextSeq = nextSeq
		if err != nil {
			// If we have already generated an auth vector successfully, then we can
			// return it. Otherwise, we must signal an error.
			// See 3GPP TS 29.272 section 5.2.3.1.3.
			if i == 0 {
				return nil, 0, err
			}
			glog.Errorf("failed to generate lte auth vector: %v", err)
			break
		}
		vectors = append(vectors, vector)
	}
	return vectors, lteAuthNextSeq, nil
}

// GenerateLteAuthVector returns the lte auth vector for the subscriber.
// Inputs:
//   milenage: The cipher to use to generate the vector
//   subscriber: The subscriber data for the subscriber we want to generate auth vectors for
//   plmn: 24 bit network identifier
//   authSqnInd: the IND of the current vector being generated
// Returns: A E-UTRAN vector and the next value to set the subscriber's LteAuthNextSeq to (or an error).
func GenerateLteAuthVector(milenage *crypto.MilenageCipher, subscriber *protos.SubscriberData, plmn, lteAuthOp []byte, authSqnInd uint64) (*crypto.EutranVector, uint64, error) {
	lte := subscriber.Lte
	if err := ValidateLteSubscription(lte); err != nil {
		return nil, 0, NewAuthRejectedError(err.Error())
	}
	if subscriber.State == nil {
		return nil, 0, NewAuthRejectedError("Subscriber data missing subscriber state")
	}

	opc, err := GetOrGenerateOpc(lte, lteAuthOp)
	if err != nil {
		return nil, 0, err
	}

	sqn := SeqToSqn(subscriber.State.LteAuthNextSeq, authSqnInd)
	vector, err := milenage.GenerateEutranVector(lte.AuthKey, opc, sqn, plmn)
	if err != nil {
		return vector, 0, NewAuthRejectedError(err.Error())
	}
	return vector, subscriber.State.LteAuthNextSeq + 1, err
}

// ResyncLteAuthSeq validates a re-synchronization request and computes the SEQ
// from the AUTS sent by U-SIM. The next value of lteAuthNextSeq (or an error) is returned.
// See 3GPP TS 33.102 section 6.3.5.
func ResyncLteAuthSeq(subscriber *protos.SubscriberData, resyncInfo, lteAuthOp []byte) (uint64, error) {
	if subscriber.State == nil {
		return 0, NewAuthDataUnavailableError("subscriber state is nil")
	}

	if IsAllZero(resyncInfo) {
		return subscriber.State.LteAuthNextSeq, nil
	}
	if len(resyncInfo) != lteResyncInfoBytes {
		err := NewAuthRejectedError(fmt.Sprintf("resync info incorrect length. expected %v bytes, but got %v bytes", lteResyncInfoBytes, len(resyncInfo)))
		return 0, err
	}
	lte := subscriber.Lte
	if err := ValidateLteSubscription(lte); err != nil {
		return 0, NewAuthRejectedError(err.Error())
	}

	// Use dummy AMF for re-synchronization. See 3GPP TS 33.102 section 6.3.3.
	milenage, err := crypto.NewMilenageCipher(make([]byte, crypto.ExpectedAmfBytes))
	if err != nil {
		return 0, NewAuthDataUnavailableError(err.Error())
	}
	rand := resyncInfo[:crypto.RandChallengeBytes]
	auts := resyncInfo[crypto.RandChallengeBytes:]
	opc, err := GetOrGenerateOpc(lte, lteAuthOp)
	if err != nil {
		return 0, err
	}
	sqnMs, macS, err := milenage.GenerateResync(auts, subscriber.Lte.AuthKey, opc, rand)
	if err != nil {
		return 0, NewAuthDataUnavailableError(err.Error())
	}
	if !bytes.Equal(macS[:], auts[crypto.ExpectedAutsBytes-len(macS):]) {
		return 0, NewAuthRejectedError("Invalid resync authentication code")
	}

	return GetNextLteAuthSqnAfterResync(subscriber.State, sqnMs)
}

// GetNextLteAuthSqnAfterResync returns the value of the next sequence number after
// sqn or an error if a resync should not occur.
// See 3GPP TS 33.102 Appendix C.3.
func GetNextLteAuthSqnAfterResync(state *protos.SubscriberState, sqn uint64) (uint64, error) {
	if state == nil {
		return 0, NewAuthDataUnavailableError("subscriber state was nil")
	}

	seq, _ := SplitSqn(sqn)
	currentSeq := state.LteAuthNextSeq - 1
	if seq < currentSeq {
		seqDelta := currentSeq - seq
		if seqDelta <= maxSeqDelta {
			// This error indicates that the last sequence number should have been
			// accepted by the USIM but wasn't (this should never happen).
			return 0, NewAuthRejectedError(fmt.Sprintf("Re-sync delta in range but UE rejected auth: %d", seqDelta))
		}
	}

	return seq + 1, nil
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

// GetOrGenerateOpc returns lte.AuthOpc and generates if it isn't stored in the proto
func GetOrGenerateOpc(lte *protos.LTESubscription, lteAuthOp []byte) ([]byte, error) {
	if lte == nil || len(lte.AuthOpc) == 0 {
		opc, err := crypto.GenerateOpc(lte.AuthKey, lteAuthOp)
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

// IsAllZero returns true if and only if the slice contains only zero bytes.
func IsAllZero(bytes []byte) bool {
	for _, b := range bytes {
		if b != 0 {
			return false
		}
	}
	return true
}
