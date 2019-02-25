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
	"magma/feg/gateway/services/testcore/hss/servicers"

	definitions "magma/feg/gateway/services/swx_proxy/servicers"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/dict"
	"github.com/stretchr/testify/assert"
)

func TestNewMAA_SuccessfulResponse(t *testing.T) {
	mar := createMAR("sub1")
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
	mar := createMAR("sub1")
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

func createMAR(userName string) *diam.Message {
	mar := createBaseMAR()
	mar.NewAVP(avp.UserName, avp.Mbit, diameter.Vendor3GPP, datatype.UTF8String(userName))
	mar.NewAVP(avp.RATType, avp.Mbit, diameter.Vendor3GPP, datatype.Unsigned32(0))
	mar.NewAVP(avp.SIPNumberAuthItems, avp.Mbit|avp.Vbit, 0, datatype.Unsigned32(1))
	mar.NewAVP(avp.SIPAuthDataItem, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.SIPAuthenticationScheme, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.UTF8String("EAP-AKA")),
		},
	})
	return mar
}
