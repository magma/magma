/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test

import (
	"testing"

	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/testcore/hss/crypto"
	"magma/feg/gateway/services/testcore/hss/servicers"
	"magma/feg/gateway/services/testcore/hss/storage"

	"magma/feg/cloud/go/protos"
	definitions "magma/feg/gateway/services/swx_proxy/servicers"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/dict"
	"github.com/stretchr/testify/assert"
)

func TestNewMAA_SuccessfulResponse(t *testing.T) {
	mar := createMARWithSingleAuthItem("sub1")
	server := newTestHomeSubscriberServer(t)
	response, err := servicers.NewMAA(server, mar)
	assert.NoError(t, err)

	var maa definitions.MAA
	err = response.Unmarshal(&maa)
	assert.NoError(t, err)
	assert.Equal(t, "magma;123_1234", maa.SessionID)
	assert.Equal(t, diam.Success, int(maa.ResultCode))
	assert.Equal(t, uint32(diam.Success), maa.ExperimentalResult.ExperimentalResultCode)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), maa.OriginHost)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), maa.OriginRealm)
	checkSIPAuthVectors(t, maa, 1)
}

func TestNewMAA_UnknownIMSI(t *testing.T) {
	mar := createMARWithSingleAuthItem("sub_unknown")
	server := newTestHomeSubscriberServer(t)
	response, err := servicers.NewMAA(server, mar)
	assert.Exactly(t, storage.NewUnknownSubscriberError("sub_unknown"), err)

	// Check that the MAA is a failure message.
	var maa definitions.MAA
	err = response.Unmarshal(&maa)
	assert.NoError(t, err)
	assert.Equal(t, uint32(protos.ErrorCode_USER_UNKNOWN), maa.ExperimentalResult.ExperimentalResultCode)
}

func TestNewMAA_MissingAuthKey(t *testing.T) {
	server := newTestHomeSubscriberServer(t)

	mar := createMARWithSingleAuthItem("missing_auth_key")
	response, err := servicers.NewMAA(server, mar)
	assert.Exactly(t, servicers.NewAuthRejectedError("incorrect key size. Expected 16 bytes, but got 0 bytes"), err)

	// Check that the MAA has the expected error.
	var maa definitions.MAA
	err = response.Unmarshal(&maa)
	assert.NoError(t, err)
	assert.Equal(t, uint32(protos.ErrorCode_AUTHORIZATION_REJECTED), maa.ExperimentalResult.ExperimentalResultCode)
	checkSIPAuthVectors(t, maa, 0)
	assert.Equal(t, "magma;123_1234", maa.SessionID)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), maa.OriginHost)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), maa.OriginRealm)
}

func TestNewMAA_MultipleVectors(t *testing.T) {
	server := newTestHomeSubscriberServer(t)
	mar := createMARExtended("sub1", 3)
	response, err := servicers.NewMAA(server, mar)
	assert.NoError(t, err)

	var maa definitions.MAA
	err = response.Unmarshal(&maa)
	assert.NoError(t, err)
	checkSIPAuthVectors(t, maa, 3)
}

func TestNewMAA_MissingAVP(t *testing.T) {
	mar := createBaseMAR()
	mar.NewAVP(avp.RATType, avp.Mbit, diameter.Vendor3GPP, datatype.Unsigned32(0))
	mar.NewAVP(avp.SIPNumberAuthItems, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(1))
	mar.NewAVP(avp.SIPAuthDataItem, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.SIPAuthenticationScheme, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.UTF8String("EAP-AKA")),
		},
	})

	server := newTestHomeSubscriberServer(t)
	response, err := servicers.NewMAA(server, mar)
	assert.EqualError(t, err, "Missing IMSI in message")

	var maa definitions.MAA
	err = response.Unmarshal(&maa)
	assert.NoError(t, err)
	assert.Equal(t, uint32(diam.MissingAVP), maa.ResultCode)
}

func TestValidateMAR_Success(t *testing.T) {
	mar := createMARWithSingleAuthItem("sub1")
	err := servicers.ValidateMAR(mar)
	assert.NoError(t, err)
}

func TestValidateMAR_MissingUserName(t *testing.T) {
	mar := createBaseMAR()
	mar.NewAVP(avp.RATType, avp.Mbit, diameter.Vendor3GPP, datatype.Unsigned32(0))
	mar.NewAVP(avp.SIPNumberAuthItems, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(1))
	mar.NewAVP(avp.SIPAuthDataItem, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.SIPAuthenticationScheme, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.UTF8String("EAP-AKA")),
		},
	})

	err := servicers.ValidateMAR(mar)
	assert.EqualError(t, err, "Missing IMSI in message")
}

func TestValidateMAR_MissingSIPNumberAuthItems(t *testing.T) {
	mar := createBaseMAR()
	mar.NewAVP(avp.UserName, avp.Mbit, diameter.Vendor3GPP, datatype.UTF8String("sub1"))
	mar.NewAVP(avp.RATType, avp.Mbit, diameter.Vendor3GPP, datatype.Unsigned32(0))
	mar.NewAVP(avp.SIPAuthDataItem, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.SIPAuthenticationScheme, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.UTF8String("EAP-AKA")),
		},
	})

	err := servicers.ValidateMAR(mar)
	assert.EqualError(t, err, "Missing SIP-Number-Auth-Items in message")
}

func TestValidateMAR_MissingSIPAuthDataItem(t *testing.T) {
	mar := createBaseMAR()
	mar.NewAVP(avp.UserName, avp.Mbit, diameter.Vendor3GPP, datatype.UTF8String("sub1"))
	mar.NewAVP(avp.RATType, avp.Mbit, diameter.Vendor3GPP, datatype.Unsigned32(0))
	mar.NewAVP(avp.SIPNumberAuthItems, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(1))
	mar.NewAVP(avp.SIPAuthenticationScheme, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.UTF8String("EAP-AKA"))

	err := servicers.ValidateMAR(mar)
	assert.EqualError(t, err, "Missing SIP-Auth-Data-Item in message")
}

func TestValidateMAR_MissingSIPAuthenticationScheme(t *testing.T) {
	mar := createBaseMAR()
	mar.NewAVP(avp.UserName, avp.Mbit, diameter.Vendor3GPP, datatype.UTF8String("sub1"))
	mar.NewAVP(avp.RATType, avp.Mbit, diameter.Vendor3GPP, datatype.Unsigned32(0))
	mar.NewAVP(avp.SIPNumberAuthItems, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(1))
	mar.NewAVP(avp.SIPAuthDataItem, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{})

	err := servicers.ValidateMAR(mar)
	assert.EqualError(t, err, "Missing SIP-Authentication-Scheme in message")
}

func TestValidateMAR_MissingRATType(t *testing.T) {
	mar := createBaseMAR()
	mar.NewAVP(avp.UserName, avp.Mbit, diameter.Vendor3GPP, datatype.UTF8String("sub1"))
	mar.NewAVP(avp.SIPNumberAuthItems, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(1))
	mar.NewAVP(avp.SIPAuthDataItem, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.SIPAuthenticationScheme, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.UTF8String("EAP-AKA")),
		},
	})

	err := servicers.ValidateMAR(mar)
	assert.EqualError(t, err, "Missing RAT type in message")
}

func TestValidateMAR_NilMessage(t *testing.T) {
	err := servicers.ValidateMAR(nil)
	assert.EqualError(t, err, "Message is nil")
}

func createBaseMAR() *diam.Message {
	mar := diameter.NewProxiableRequest(diam.MultimediaAuthentication, diam.TGPP_SWX_APP_ID, dict.Default)
	mar.NewAVP(avp.SessionID, avp.Mbit, diameter.Vendor3GPP, datatype.UTF8String("magma;123_1234"))
	mar.NewAVP(avp.OriginHost, avp.Mbit, diameter.Vendor3GPP, datatype.DiameterIdentity("magma.com"))
	mar.NewAVP(avp.OriginRealm, avp.Mbit, diameter.Vendor3GPP, datatype.DiameterIdentity("magma.com"))
	return mar
}

func createMARWithSingleAuthItem(userName string) *diam.Message {
	return createMARExtended(userName, 1)
}

func createMARExtended(userName string, numberAuthItems uint32) *diam.Message {
	mar := createBaseMAR()
	mar.NewAVP(avp.UserName, avp.Mbit, diameter.Vendor3GPP, datatype.UTF8String(userName))
	mar.NewAVP(avp.RATType, avp.Mbit, diameter.Vendor3GPP, datatype.Unsigned32(0))
	mar.NewAVP(avp.SIPNumberAuthItems, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(numberAuthItems))
	mar.NewAVP(avp.SIPAuthDataItem, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.SIPAuthenticationScheme, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.UTF8String("EAP-AKA")),
		},
	})
	return mar
}

func checkSIPAuthVectors(t *testing.T, maa definitions.MAA, expectedNumVectors uint32) {
	assert.Equal(t, int(expectedNumVectors), len(maa.SIPAuthDataItems))
	assert.Equal(t, expectedNumVectors, maa.SIPNumberAuthItems)

	for _, vector := range maa.SIPAuthDataItems {
		assert.Equal(t, definitions.SipAuthScheme_EAP_AKA, vector.AuthScheme)
		assert.Equal(t, crypto.RandChallengeBytes+crypto.AutnBytes, len(vector.Authenticate))
		assert.Equal(t, crypto.XresBytes, len(vector.Authorization))
		assert.Equal(t, crypto.ConfidentialityKeyBytes, len(vector.ConfidentialityKey))
		assert.Equal(t, crypto.IntegrityKeyBytes, len(vector.IntegrityKey))
	}
}
