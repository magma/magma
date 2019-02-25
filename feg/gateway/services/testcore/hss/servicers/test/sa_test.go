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

func TestNewSAA_SuccessfulResponse(t *testing.T) {
	sar := createSAR("sub1")
	server := newTestHomeSubscriberServer(t)
	response, err := servicers.NewSAA(server, sar)
	assert.NoError(t, err)

	var saa definitions.SAA
	err = response.Unmarshal(&saa)
	assert.NoError(t, err)
	assert.Equal(t, "magma;123_1234", saa.SessionID)
	assert.Equal(t, diam.Success, int(saa.ResultCode))
	assert.Equal(t, uint32(diam.Success), saa.ExperimentalResult.ExperimentalResultCode)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), saa.OriginHost)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), saa.OriginRealm)
}

func TestNewSAA_MissingAVP(t *testing.T) {
	sar := createBaseSAR()
	sar.NewAVP(avp.UserName, avp.Mbit, diameter.Vendor3GPP, datatype.UTF8String("sub1"))
	server := newTestHomeSubscriberServer(t)
	response, err := servicers.NewSAA(server, sar)
	assert.EqualError(t, err, "Missing server assignment type in message")

	var saa definitions.SAA
	err = response.Unmarshal(&saa)
	assert.NoError(t, err)
	assert.Equal(t, diam.MissingAVP, int(saa.ResultCode))
}

func TestValidateSAR_MissingUserName(t *testing.T) {
	sar := createBaseSAR()
	sar.NewAVP(avp.ServerAssignmentType, avp.Mbit, diameter.Vendor3GPP, datatype.Enumerated(definitions.ServerAssignmentType_REGISTRATION))
	err := servicers.ValidateSAR(sar)
	assert.EqualError(t, err, "Missing IMSI in message")
}

func TestValidateSAR_MissingServerAssignmentType(t *testing.T) {
	sar := createBaseSAR()
	sar.NewAVP(avp.UserName, avp.Mbit, diameter.Vendor3GPP, datatype.UTF8String("sub1"))
	err := servicers.ValidateSAR(sar)
	assert.EqualError(t, err, "Missing server assignment type in message")
}

func TestValidateSAR_NilMessage(t *testing.T) {
	err := servicers.ValidateSAR(nil)
	assert.EqualError(t, err, "Message is nil")
}

func TestValidateSAR_Success(t *testing.T) {
	sar := createSAR("sub1")
	err := servicers.ValidateSAR(sar)
	assert.NoError(t, err)
}

func createBaseSAR() *diam.Message {
	sar := diameter.NewProxiableRequest(diam.ServerAssignment, diam.TGPP_SWX_APP_ID, dict.Default)
	sar.NewAVP(avp.SessionID, avp.Mbit, diameter.Vendor3GPP, datatype.UTF8String("magma;123_1234"))
	sar.NewAVP(avp.OriginHost, avp.Mbit, diameter.Vendor3GPP, datatype.DiameterIdentity("magma.com"))
	sar.NewAVP(avp.OriginRealm, avp.Mbit, diameter.Vendor3GPP, datatype.DiameterIdentity("magma.com"))
	return sar
}

func createSAR(userName string) *diam.Message {
	sar := createBaseSAR()
	sar.NewAVP(avp.UserName, avp.Mbit, diameter.Vendor3GPP, datatype.UTF8String(userName))
	sar.NewAVP(avp.ServerAssignmentType, avp.Mbit, diameter.Vendor3GPP, datatype.Enumerated(definitions.ServerAssignmentType_REGISTRATION))
	return sar
}
