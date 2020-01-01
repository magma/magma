/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers_test

import (
	"context"
	"testing"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	definitions "magma/feg/gateway/services/s6a_proxy/servicers"
	hss "magma/feg/gateway/services/testcore/hss/servicers"
	"magma/feg/gateway/services/testcore/hss/servicers/test"
	"magma/feg/gateway/services/testcore/hss/storage"
	"magma/lte/cloud/go/crypto"
	lteprotos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/eps_authentication/servicers"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/stretchr/testify/assert"
)

func TestNewAIA_MissingSessionID(t *testing.T) {
	m := diameter.NewProxiableRequest(diam.AuthenticationInformation, diam.TGPP_S6A_APP_ID, dict.Default)
	server := test.NewTestHomeSubscriberServer(t)
	response, err := hss.NewAIA(server, m)
	assert.Error(t, err)

	// Check that the AIA is a failure message.
	var aia definitions.AIA
	err = response.Unmarshal(&aia)
	assert.NoError(t, err)
	assert.Equal(t, diam.MissingAVP, int(aia.ResultCode))
}

func TestNewAIA_UnknownIMSI(t *testing.T) {
	air := createAIR("sub_unknown")
	server := test.NewTestHomeSubscriberServer(t)
	response, err := hss.NewAIA(server, air)
	assert.Exactly(t, storage.NewUnknownSubscriberError("sub_unknown"), err)

	// Check that the AIA is a failure message.
	var aia definitions.AIA
	err = response.Unmarshal(&aia)
	assert.NoError(t, err)
	assert.Equal(t, uint32(protos.ErrorCode_USER_UNKNOWN), aia.ExperimentalResult.ExperimentalResultCode)
}

func TestNewAIA_SuccessfulResponse(t *testing.T) {
	server := test.NewTestHomeSubscriberServer(t)
	amf := []byte("\x80\x00")
	rand := []byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b\x0c\r\x0e\x0f")
	milenage, err := crypto.NewMockMilenageCipher(amf, rand)
	assert.NoError(t, err)
	server.Milenage = milenage

	air := createAIR("sub1")
	response, err := hss.NewAIA(server, air)
	assert.NoError(t, err)

	// Check that the AIA has all the expected data.
	var aia definitions.AIA
	err = response.Unmarshal(&aia)
	assert.NoError(t, err)
	assert.Equal(t, "magma;123_1234", aia.SessionID)
	assert.Equal(t, diam.Success, int(aia.ResultCode))
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), aia.OriginHost)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), aia.OriginRealm)
	assert.Equal(t, 1, len(aia.AIs))

	ai := aia.AIs[0]
	assert.Equal(t, 1, len(ai.EUtranVectors))

	vec := ai.EUtranVectors[0]
	assert.Equal(t, datatype.OctetString(rand), vec.RAND)
	assert.Equal(t, datatype.OctetString([]byte("\x2d\xaf\x87\x3d\x73\xf3\x10\xc6")), vec.XRES)
	assert.Equal(t, datatype.OctetString([]byte("o\xbf\xa3\x83\x95 \x80\x00\xb1\x1f \xbd\xdc\xf5\xeeS")), vec.AUTN)
	assert.Equal(t, datatype.OctetString([]byte("Q\xd0g\xde?\x95\xecB\x94\xf8\xe7\xc4\x0f\x92\x81i\x8e\\Cu\xc1\xe5\xab\x1a\xc0\xe6z\x117\nkz")), vec.KASME)

	subscriber, err := server.GetSubscriberData(context.Background(), &lteprotos.SubscriberID{Id: "sub1"})
	assert.NoError(t, err)
	assert.Equal(t, uint64(7351), subscriber.State.LteAuthNextSeq)
}

func TestNewAIA_MultipleVectors(t *testing.T) {
	server := test.NewTestHomeSubscriberServer(t)
	air := createAIRExtended("sub1", 3)
	response, err := hss.NewAIA(server, air)
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
	server := test.NewTestHomeSubscriberServer(t)

	air := createAIR("missing_auth_key")
	response, err := hss.NewAIA(server, air)
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

	assert.EqualError(t, hss.ValidateAIR(m), "Missing IMSI in message")
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

	assert.EqualError(t, hss.ValidateAIR(m), "Missing Visited PLMN ID in message")
}

func TestValidateAIR_MissingEUTRANInfo(t *testing.T) {
	m := createBaseAIR()
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("magma;123_1234"))
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String("magma"))
	m.NewAVP(avp.VisitedPLMNID, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(0))

	assert.EqualError(t, hss.ValidateAIR(m), "Missing requested E-UTRAN authentication info in message")
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

	assert.EqualError(t, hss.ValidateAIR(m), "Missing SessionID in message")
}

func TestValidateAIR_Success(t *testing.T) {
	air := createAIR("sub1")
	assert.NoError(t, hss.ValidateAIR(air))
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

func TestNewSuccessfulAIA(t *testing.T) {
	server := test.NewTestHomeSubscriberServer(t)
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
