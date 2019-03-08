/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test

import (
	"context"
	"magma/lte/cloud/go/services/eps_authentication/crypto"
	"testing"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	definitions "magma/feg/gateway/services/s6a_proxy/servicers"
	"magma/feg/gateway/services/testcore/hss/servicers"
	"magma/feg/gateway/services/testcore/hss/storage"
	lteprotos "magma/lte/cloud/go/protos"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/dict"
	"github.com/stretchr/testify/assert"
)

func TestNewAIA_MissingSessionID(t *testing.T) {
	m := diameter.NewProxiableRequest(diam.AuthenticationInformation, diam.TGPP_S6A_APP_ID, dict.Default)
	server := newTestHomeSubscriberServer(t)
	response, err := servicers.NewAIA(server, m)
	assert.Error(t, err)

	// Check that the AIA is a failure message.
	var aia definitions.AIA
	err = response.Unmarshal(&aia)
	assert.NoError(t, err)
	assert.Equal(t, diam.MissingAVP, int(aia.ResultCode))
}

func TestNewAIA_UnknownIMSI(t *testing.T) {
	air := createAIR("sub_unknown")
	server := newTestHomeSubscriberServer(t)
	response, err := servicers.NewAIA(server, air)
	assert.Exactly(t, storage.NewUnknownSubscriberError("sub_unknown"), err)

	// Check that the AIA is a failure message.
	var aia definitions.AIA
	err = response.Unmarshal(&aia)
	assert.NoError(t, err)
	assert.Equal(t, uint32(protos.ErrorCode_USER_UNKNOWN), aia.ExperimentalResult.ExperimentalResultCode)
}

func TestNewAIA_SuccessfulResponse(t *testing.T) {
	server := newTestHomeSubscriberServer(t)
	amf := []byte("\x80\x00")
	rand := []byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b\x0c\r\x0e\x0f")
	milenage, err := crypto.NewMockMilenageCipher(amf, rand)
	assert.NoError(t, err)
	server.Milenage = milenage

	air := createAIR("sub1")
	response, err := servicers.NewAIA(server, air)
	assert.NoError(t, err)

	// Check that the AIA has all the expected data.
	var aia definitions.AIA
	err = response.Unmarshal(&aia)
	assert.NoError(t, err)
	assert.Equal(t, "magma;123_1234", aia.SessionID)
	assert.Equal(t, diam.Success, int(aia.ResultCode))
	assert.Equal(t, uint32(diam.Success), aia.ExperimentalResult.ExperimentalResultCode)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), aia.OriginHost)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), aia.OriginRealm)
	assert.Equal(t, 1, len(aia.AIs))

	ai := aia.AIs[0]
	assert.Equal(t, 1, len(ai.EUtranVectors))

	vec := ai.EUtranVectors[0]
	assert.Equal(t, datatype.OctetString(rand), vec.RAND)
	assert.Equal(t, datatype.OctetString([]byte("\x2d\xaf\x87\x3d\x73\xf3\x10\xc6")), vec.XRES)
	assert.Equal(t, datatype.OctetString([]byte{0x6f, 0xbf, 0xa3, 0x83, 0x95, 0x0, 0x80, 0x0, 0x9f, 0xbc, 0xe8, 0xd3, 0x47, 0xe, 0x82, 0xd5}), vec.AUTN)
	assert.Equal(t, datatype.OctetString([]byte{0x62, 0x23, 0xd, 0x4d, 0x26, 0xec, 0xa, 0x12, 0x35, 0x54, 0x6, 0x85, 0x5, 0x5a, 0x94, 0xf8, 0x61, 0x53, 0x71, 0x4b, 0xd9, 0x42, 0xe, 0x64, 0xf1, 0x2f, 0x55, 0xd5, 0x84, 0xca, 0xd9, 0x6}), vec.KASME)

	subscriber, err := server.GetSubscriberData(context.Background(), &lteprotos.SubscriberID{Id: "sub1"})
	assert.NoError(t, err)
	assert.Equal(t, uint64(7351), subscriber.State.LteAuthNextSeq)
}

func TestNewAIA_MultipleVectors(t *testing.T) {
	server := newTestHomeSubscriberServer(t)
	air := createAIRExtended("sub1", 3)
	response, err := servicers.NewAIA(server, air)
	assert.NoError(t, err)

	var aia definitions.AIA
	err = response.Unmarshal(&aia)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(aia.AIs))

	for i := 0; i < len(aia.AIs); i++ {
		assert.Equal(t, 1, len(aia.AIs[i].EUtranVectors))
		vector := aia.AIs[i].EUtranVectors[0]
		assert.Equal(t, crypto.RandChallengeBytes, len(vector.RAND))
		assert.Equal(t, crypto.XresBytes, len(vector.XRES))
		assert.Equal(t, crypto.AutnBytes, len(vector.AUTN))
		assert.Equal(t, crypto.KasmeBytes, len(vector.KASME))

		for j := i + 1; j < len(aia.AIs); j++ {
			assert.NotEqual(t, aia.AIs[i], aia.AIs[j])
		}
	}
}

func TestNewAIA_MissingAuthKey(t *testing.T) {
	server := newTestHomeSubscriberServer(t)

	air := createAIR("missing_auth_key")
	response, err := servicers.NewAIA(server, air)
	assert.Exactly(t, servicers.NewAuthRejectedError("incorrect key size. Expected 16 bytes, but got 0 bytes"), err)

	// Check that the AIA has the expected error.
	var aia definitions.AIA
	err = response.Unmarshal(&aia)
	assert.NoError(t, err)
	assert.Equal(t, uint32(protos.ErrorCode_AUTHORIZATION_REJECTED), aia.ExperimentalResult.ExperimentalResultCode)
	assert.Equal(t, 0, len(aia.AIs))
	assert.Equal(t, "magma;123_1234", aia.SessionID)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), aia.OriginHost)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), aia.OriginRealm)
}

func TestValidateAIR_MissingUserName(t *testing.T) {
	m := createBaseAIR()
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("magma;123_1234"))
	m.NewAVP(avp.VisitedPLMNID, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(0))
	authInfo := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(
				avp.NumberOfRequestedVectors,
				avp.Vbit|avp.Mbit,
				diameter.Vendor3GPP,
				datatype.Unsigned32(1)),
			diam.NewAVP(
				avp.ImmediateResponsePreferred, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, datatype.Unsigned32(0)),
		},
	}
	m.NewAVP(avp.RequestedEUTRANAuthenticationInfo, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, authInfo)

	assert.EqualError(t, servicers.ValidateAIR(m), "Missing IMSI in message")
}

func TestValidateAIR_MissingVistedPLMNID(t *testing.T) {
	m := createBaseAIR()
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("magma;123_1234"))
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String("magma"))
	authInfo := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(
				avp.NumberOfRequestedVectors,
				avp.Vbit|avp.Mbit,
				diameter.Vendor3GPP,
				datatype.Unsigned32(1)),
			diam.NewAVP(
				avp.ImmediateResponsePreferred, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, datatype.Unsigned32(0)),
		},
	}
	m.NewAVP(avp.RequestedEUTRANAuthenticationInfo, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, authInfo)

	assert.EqualError(t, servicers.ValidateAIR(m), "Missing Visited PLMN ID in message")
}

func TestValidateAIR_MissingEUTRANInfo(t *testing.T) {
	m := createBaseAIR()
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("magma;123_1234"))
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String("magma"))
	m.NewAVP(avp.VisitedPLMNID, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(0))

	assert.EqualError(t, servicers.ValidateAIR(m), "Missing requested E-UTRAN authentication info in message")
}

func TestValidateAIR_MissingSessionId(t *testing.T) {
	m := createBaseAIR()
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String("magma"))
	m.NewAVP(avp.VisitedPLMNID, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(0))
	authInfo := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(
				avp.NumberOfRequestedVectors,
				avp.Vbit|avp.Mbit,
				diameter.Vendor3GPP,
				datatype.Unsigned32(1)),
			diam.NewAVP(
				avp.ImmediateResponsePreferred, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, datatype.Unsigned32(0)),
		},
	}
	m.NewAVP(avp.RequestedEUTRANAuthenticationInfo, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, authInfo)

	assert.EqualError(t, servicers.ValidateAIR(m), "Missing SessionID in message")
}

func TestValidateAIR_Success(t *testing.T) {
	air := createAIR("sub1")
	assert.NoError(t, servicers.ValidateAIR(air))
}

// createBaseAIR outputs a mock authentication information request with only a
// few AVPs added.
func createBaseAIR() *diam.Message {
	air := diameter.NewProxiableRequest(diam.AuthenticationInformation, diam.TGPP_S6A_APP_ID, dict.Default)
	air.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("magma.com"))
	air.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("magma.com"))
	air.NewAVP(avp.AuthSessionState, avp.Mbit, 0, datatype.Enumerated(1))
	return air
}

// createAIR outputs a mock authentication information request.
func createAIR(userName string) *diam.Message {
	return createAIRExtended(userName, 1)
}

// createAIRExtended outputs a mock authentication information request.
// It allows specifying more options than createAIR.
func createAIRExtended(userName string, numRequestedVectors uint32) *diam.Message {
	m := createBaseAIR()
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("magma;123_1234"))
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String(userName))
	m.NewAVP(avp.VisitedPLMNID, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(0))
	authInfo := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(
				avp.NumberOfRequestedVectors,
				avp.Vbit|avp.Mbit,
				diameter.Vendor3GPP,
				datatype.Unsigned32(numRequestedVectors)),
			diam.NewAVP(
				avp.ImmediateResponsePreferred, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, datatype.Unsigned32(0)),
		},
	}
	m.NewAVP(avp.RequestedEUTRANAuthenticationInfo, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, authInfo)
	return m
}

func TestSeqToSqn(t *testing.T) {
	assert.Equal(t, uint64(0x1FE000), servicers.SeqToSqn(0xFF00, 0))
	assert.Equal(t, uint64(0xFFFFFFFFFA00), servicers.SeqToSqn(0xFFFFFFFFFFD0, 0))
	assert.Equal(t, uint64(0x142), servicers.SeqToSqn(0xA, 2))
	assert.Equal(t, uint64(0xFFFFFFFFF805), servicers.SeqToSqn(0xFFFFFFFFFFC0, 5))
}

func TestSplitSqn(t *testing.T) {
	sqn, ind := servicers.SplitSqn(0x1FE001)
	assert.Equal(t, uint64(0xFF00), sqn)
	assert.Equal(t, uint64(0x1), ind)

	sqn, ind = servicers.SplitSqn(0xFFFFFFFFFA1F)
	assert.Equal(t, uint64(0x7FFFFFFFFD0), sqn)
	assert.Equal(t, uint64(0x1F), ind)
}

func TestGetOrGenerateOpc(t *testing.T) {
	server := newTestHomeSubscriberServer(t)

	lte := &lteprotos.LTESubscription{AuthOpc: []byte("\xcdc\xcbq\x95J\x9fNH\xa5\x99N7\xa0+\xaf")}
	opc, err := server.GetOrGenerateOpc(lte)
	assert.NoError(t, err)
	assert.Equal(t, lte.AuthOpc, opc)

	lte = &lteprotos.LTESubscription{AuthKey: []byte("\x46\x5b\x5c\xe8\xb1\x99\xb4\x9f\xaa\x5f\x0a\x2e\xe2\x38\xa6\xbc")}
	opc, err = server.GetOrGenerateOpc(lte)
	assert.NoError(t, err)
	expectedOpc, err := crypto.GenerateOpc(lte.AuthKey, []byte(server.Config.LteAuthOp))
	assert.NoError(t, err)
	assert.Equal(t, expectedOpc[:], opc)
}

func TestGenerateLteAuthVector_MissingLTE(t *testing.T) {
	server := newTestHomeSubscriberServer(t)
	plmn := []byte("\x02\xf8\x59")

	subscriber := &lteprotos.SubscriberData{State: &lteprotos.SubscriberState{}}
	_, err := server.GenerateLteAuthVector(subscriber, plmn)
	assert.Exactly(t, servicers.NewAuthRejectedError("Subscriber data missing LTE subscription"), err)
}

func TestGenerateLteAuthVector_MissingSubscriberState(t *testing.T) {
	server := newTestHomeSubscriberServer(t)
	plmn := []byte("\x02\xf8\x59")

	subscriber := &lteprotos.SubscriberData{
		Lte: &lteprotos.LTESubscription{
			State:    lteprotos.LTESubscription_ACTIVE,
			AuthAlgo: lteprotos.LTESubscription_MILENAGE,
		},
	}
	_, err := server.GenerateLteAuthVector(subscriber, plmn)
	assert.Exactly(t, servicers.NewAuthRejectedError("Subscriber data missing subscriber state"), err)
}

func TestGenerateLteAuthVector_InactiveLTESubscription(t *testing.T) {
	server := newTestHomeSubscriberServer(t)
	plmn := []byte("\x02\xf8\x59")

	subscriber := &lteprotos.SubscriberData{
		Lte: &lteprotos.LTESubscription{
			State:    lteprotos.LTESubscription_INACTIVE,
			AuthAlgo: lteprotos.LTESubscription_MILENAGE,
		},
		State: &lteprotos.SubscriberState{},
	}
	_, err := server.GenerateLteAuthVector(subscriber, plmn)
	assert.Exactly(t, servicers.NewAuthRejectedError("LTE Service not active"), err)
}

func TestGenerateLteAuthVector_UnknownLTEAuthAlgo(t *testing.T) {
	server := newTestHomeSubscriberServer(t)
	plmn := []byte("\x02\xf8\x59")

	subscriber := &lteprotos.SubscriberData{
		Lte: &lteprotos.LTESubscription{
			State:    lteprotos.LTESubscription_ACTIVE,
			AuthAlgo: 10,
		},
		State: &lteprotos.SubscriberState{},
	}
	_, err := server.GenerateLteAuthVector(subscriber, plmn)
	assert.Exactly(t, servicers.NewAuthRejectedError("Unsupported crypto algorithm: 10"), err)
}

func TestGenerateLteAuthVector_Success(t *testing.T) {
	server := newTestHomeSubscriberServer(t)
	server.AuthSqnInd = 23
	rand := []byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b\x0c\r\x0e\x0f")
	milenage, err := crypto.NewMockMilenageCipher([]byte("\x80\x00"), rand)
	assert.NoError(t, err)
	server.Milenage = milenage
	plmn := []byte("\x02\xf8\x59")

	subscriber := &lteprotos.SubscriberData{
		Sid: &lteprotos.SubscriberID{Id: "sub1"},
		Lte: &lteprotos.LTESubscription{
			State:    lteprotos.LTESubscription_ACTIVE,
			AuthAlgo: lteprotos.LTESubscription_MILENAGE,
			AuthKey:  []byte("\x8b\xafG?/\x8f\xd0\x94\x87\xcc\xcb\xd7\t|hb"),
			AuthOpc:  []byte("\x8e'\xb6\xaf\x0ei.u\x0f2fz;\x14`]"),
		},
		State: &lteprotos.SubscriberState{LteAuthNextSeq: 228},
	}
	vector, err := server.GenerateLteAuthVector(subscriber, plmn)
	assert.NoError(t, err)

	assert.Equal(t, rand, vector.Rand[:])
	assert.Equal(t, []byte("\x2d\xaf\x87\x3d\x73\xf3\x10\xc6"), vector.Xres[:])
	assert.Equal(t, []byte("o\xbf\xa3\x80\x1fW\x80\x00{\xdeY\x88n\x96\xe4\xfe"), vector.Autn[:])
	assert.Equal(t, []byte("\x87H\xc1\xc0\xa2\x82o\xa4\x05\xb1\xe2~\xa1\x04CJ\xe5V\xc7e\xe8\xf0a\xeb\xdb\x8a\xe2\x86\xc4F\x16\xc2"), vector.Kasme[:])
}

func TestNewSuccessfulAIA(t *testing.T) {
	server := newTestHomeSubscriberServer(t)
	serverCfg := server.Config.Server

	msg := createAIR("user1")
	var air definitions.AIR
	err := msg.Unmarshal(&air)
	assert.NoError(t, err)

	vector := &crypto.EutranVector{}
	copy(vector.Rand[:], []byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b\x0c\r\x0e\x0f"))
	copy(vector.Xres[:], []byte("\x2d\xaf\x87\x3d\x73\xf3\x10\xc6"))
	copy(vector.Autn[:], []byte("o\xbf\xa3\x80\x1fW\x80\x00{\xdeY\x88n\x96\xe4\xfe"))
	copy(vector.Kasme[:], []byte("\x87H\xc1\xc0\xa2\x82o\xa4\x05\xb1\xe2~\xa1\x04CJ\xe5V\xc7e\xe8\xf0a\xeb\xdb\x8a\xe2\x86\xc4F\x16\xc2"))

	response := server.NewSuccessfulAIA(msg, air.SessionID, []*crypto.EutranVector{vector})
	var aia definitions.AIA
	err = response.Unmarshal(&aia)
	assert.NoError(t, err)

	assert.Equal(t, uint32(diam.Success), aia.ResultCode)
	assert.Equal(t, uint32(diam.Success), aia.ExperimentalResult.ExperimentalResultCode)
	assert.Equal(t, air.SessionID, datatype.UTF8String(aia.SessionID))
	assert.Equal(t, datatype.DiameterIdentity(serverCfg.DestHost), aia.OriginHost)
	assert.Equal(t, datatype.DiameterIdentity(serverCfg.DestRealm), aia.OriginRealm)
	assert.Equal(t, 1, len(aia.AIs))

	ai := aia.AIs[0]
	assert.Equal(t, 1, len(ai.EUtranVectors))

	vec := ai.EUtranVectors[0]
	assert.Equal(t, datatype.OctetString(vector.Rand[:]), vec.RAND)
	assert.Equal(t, datatype.OctetString(vector.Xres[:]), vec.XRES)
	assert.Equal(t, datatype.OctetString(vector.Autn[:]), vec.AUTN)
	assert.Equal(t, datatype.OctetString(vector.Kasme[:]), vec.KASME)
}

func TestResyncLteAuthSeq(t *testing.T) {
	server := newTestHomeSubscriberServer(t)
	subscriber, err := server.GetSubscriberData(context.Background(), &lteprotos.SubscriberID{Id: "sub1"})
	assert.NoError(t, err)

	err = server.ResyncLteAuthSeq(subscriber, nil)
	assert.NoError(t, err)

	err = server.ResyncLteAuthSeq(subscriber, make([]byte, 30))
	assert.NoError(t, err)

	resyncInfo := make([]byte, 50)
	resyncInfo[25] = 1
	err = server.ResyncLteAuthSeq(subscriber, resyncInfo)
	assert.Exactly(t, servicers.NewAuthRejectedError("resync info incorrect length. expected 30 bytes, but got 50 bytes"), err)

	resyncInfo = make([]byte, 30)
	resyncInfo[0] = 0xFF
	err = server.ResyncLteAuthSeq(subscriber, resyncInfo)
	assert.Exactly(t, servicers.NewAuthRejectedError("Invalid resync authentication code"), err)

	macS := []byte{132, 178, 239, 23, 199, 61, 138, 176}
	copy(resyncInfo[22:], macS)
	err = server.ResyncLteAuthSeq(subscriber, resyncInfo)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0x4204c05f18a9001), subscriber.State.LteAuthNextSeq)
}

func TestSetNextLteAuthSqnAfterResync(t *testing.T) {
	server := newTestHomeSubscriberServer(t)

	id := &lteprotos.SubscriberID{Id: "sub1"}
	subscriber, err := server.GetSubscriberData(context.Background(), id)
	assert.NoError(t, err)

	err = server.SetNextLteAuthSeq(subscriber, 1<<30)
	assert.NoError(t, err)

	err = server.SetNextLteAuthSqnAfterResync(subscriber, servicers.SeqToSqn(1<<30-1<<10, 2))
	assert.Exactly(t, servicers.NewAuthRejectedError("Re-sync delta in range but UE rejected auth: 1023"), err)

	err = server.SetNextLteAuthSqnAfterResync(subscriber, servicers.SeqToSqn(1<<30-1, 3))
	assert.NoError(t, err)
	assert.Equal(t, uint64(1<<30), subscriber.State.LteAuthNextSeq)
}

func TestSetNextLteAuthSeq(t *testing.T) {
	server := newTestHomeSubscriberServer(t)

	id := &lteprotos.SubscriberID{Id: "sub1"}
	subscriber, err := server.GetSubscriberData(context.Background(), id)
	assert.NoError(t, err)

	err = server.SetNextLteAuthSeq(subscriber, 100)
	assert.NoError(t, err)
	assert.Equal(t, uint64(100), subscriber.State.LteAuthNextSeq)

	subscriber, err = server.GetSubscriberData(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, uint64(100), subscriber.State.LteAuthNextSeq)
}

func TestIncrementSQN(t *testing.T) {
	server := newTestHomeSubscriberServer(t)

	id := &lteprotos.SubscriberID{Id: "sub1"}
	subscriber, err := server.GetSubscriberData(context.Background(), id)
	assert.NoError(t, err)

	err = server.SetNextLteAuthSeq(subscriber, 50)
	assert.NoError(t, err)

	err = server.IncreaseSQN(subscriber)
	assert.NoError(t, err)
	assert.Equal(t, uint64(51), subscriber.State.LteAuthNextSeq)

	subscriber, err = server.GetSubscriberData(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, uint64(51), subscriber.State.LteAuthNextSeq)
}

func TestValidateLteSubscription(t *testing.T) {
	err := servicers.ValidateLteSubscription(nil)
	assert.EqualError(t, err, "Subscriber data missing LTE subscription")

	lte := &lteprotos.LTESubscription{
		State:    lteprotos.LTESubscription_INACTIVE,
		AuthAlgo: lteprotos.LTESubscription_MILENAGE,
	}
	err = servicers.ValidateLteSubscription(lte)
	assert.EqualError(t, err, "LTE Service not active")

	lte = &lteprotos.LTESubscription{
		State:    lteprotos.LTESubscription_ACTIVE,
		AuthAlgo: 50,
	}
	err = servicers.ValidateLteSubscription(lte)
	assert.EqualError(t, err, "Unsupported crypto algorithm: 50")

	lte = &lteprotos.LTESubscription{
		State:    lteprotos.LTESubscription_ACTIVE,
		AuthAlgo: lteprotos.LTESubscription_MILENAGE,
	}
	err = servicers.ValidateLteSubscription(lte)
	assert.NoError(t, err)
}
