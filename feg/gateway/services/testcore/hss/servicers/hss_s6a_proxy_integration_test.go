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
	"magma/feg/gateway/services/s6a_proxy/servicers"
	"magma/feg/gateway/services/testcore/hss/servicers/test_utils"
	"magma/lte/cloud/go/crypto"

	"github.com/stretchr/testify/assert"
)

func TestAIR_Successful(t *testing.T) {
	s6aProxy := getTestS6aProxy(t)
	air := &protos.AuthenticationInformationRequest{
		UserName:                  "sub1",
		VisitedPlmn:               []byte{0, 0, 0},
		NumRequestedEutranVectors: 1,
	}

	aia, err := s6aProxy.AuthenticationInformation(context.Background(), air)
	assert.NoError(t, err)
	assert.Equal(t, protos.ErrorCode_UNDEFINED, aia.ErrorCode)

	assert.Equal(t, 1, len(aia.EutranVectors))
	vector := aia.EutranVectors[0]
	assert.Equal(t, crypto.RandChallengeBytes, len(vector.Rand))
	assert.Equal(t, crypto.XresBytes, len(vector.Xres))
	assert.Equal(t, crypto.AutnBytes, len(vector.Autn))
	assert.Equal(t, crypto.KasmeBytes, len(vector.Kasme))
}

func TestAIR_UnknownIMSI(t *testing.T) {
	s6aProxy := getTestS6aProxy(t)
	air := &protos.AuthenticationInformationRequest{
		UserName:                  "sub_unknown",
		VisitedPlmn:               []byte{0, 0, 0},
		NumRequestedEutranVectors: 1,
	}

	aia, err := s6aProxy.AuthenticationInformation(context.Background(), air)
	assert.NoError(t, err)
	assert.Equal(t, protos.ErrorCode_USER_UNKNOWN, aia.ErrorCode)
	assert.Equal(t, 0, len(aia.EutranVectors))
}

func TestULR_Successful(t *testing.T) {
	s6aProxy := getTestS6aProxy(t)
	ulr := &protos.UpdateLocationRequest{
		UserName:    "sub1",
		VisitedPlmn: []byte{0, 0, 0},
	}

	ula, err := s6aProxy.UpdateLocation(context.Background(), ulr)
	assert.NoError(t, err)
	assert.Equal(t, protos.ErrorCode_UNDEFINED, ula.ErrorCode)
	assert.Equal(t, uint32(test_utils.DefaultMaxUlBitRate), ula.GetTotalAmbr().GetMaxBandwidthUl())
	assert.Equal(t, uint32(test_utils.DefaultMaxDlBitRate), ula.GetTotalAmbr().GetMaxBandwidthDl())
	assert.Equal(t, []byte("12345"), ula.Msisdn)

	assert.Equal(t, 1, len(ula.Apn))
	apn := ula.Apn[0]
	assert.Equal(t, "oai.ipv4", apn.ServiceSelection)
	assert.Equal(t, uint32(test_utils.DefaultMaxUlBitRate), apn.GetAmbr().GetMaxBandwidthUl())
	assert.Equal(t, uint32(test_utils.DefaultMaxDlBitRate), apn.GetAmbr().GetMaxBandwidthDl())
	assert.Equal(t, int32(9), apn.GetQosProfile().GetClassId())
	assert.Equal(t, true, apn.GetQosProfile().GetPreemptionVulnerability())
	assert.Equal(t, uint32(15), apn.GetQosProfile().GetPriorityLevel())
	assert.Equal(t, false, apn.GetQosProfile().GetPreemptionCapability())
}

func TestULR_UnknownIMSI(t *testing.T) {
	s6aProxy := getTestS6aProxy(t)
	ulr := &protos.UpdateLocationRequest{
		UserName:    "sub_unknown",
		VisitedPlmn: []byte{0, 0, 0},
	}

	ula, err := s6aProxy.UpdateLocation(context.Background(), ulr)
	assert.NoError(t, err)
	assert.Equal(t, protos.ErrorCode_USER_UNKNOWN, ula.ErrorCode)
	assert.Equal(t, 0, len(ula.Apn))
}

// getTestS6aProxy creates a s6a proxy server and test hss diameter
// server which are configured to communicate with each other.
func getTestS6aProxy(t *testing.T) protos.S6AProxyServer {
	hss := getTestHSSDiameterServer(t)
	serverCfg := hss.Config.Server

	// Create an s6a proxy server.
	clientCfg := &diameter.DiameterClientConfig{
		Host:             serverCfg.DestHost,
		Realm:            serverCfg.DestRealm,
		ProductName:      "magma",
		AppID:            0,
		AuthAppID:        0,
		Retransmits:      3,
		WatchdogInterval: 10,
		RetryCount:       3,
	}
	diameterServerCfg := &diameter.DiameterServerConfig{
		DiameterServerConnConfig: diameter.DiameterServerConnConfig{
			Addr:      serverCfg.Address,
			Protocol:  serverCfg.Protocol,
			LocalAddr: serverCfg.LocalAddress},
		DestHost:  serverCfg.DestHost,
		DestRealm: serverCfg.DestRealm,
	}
	s6aProxy, err := servicers.NewS6aProxy(clientCfg, diameterServerCfg)
	assert.NoError(t, err)

	return s6aProxy
}
