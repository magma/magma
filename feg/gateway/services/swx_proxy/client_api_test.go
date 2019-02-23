/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package swx_proxy_test

import (
	"strconv"
	"testing"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/swx_proxy"
	"magma/feg/gateway/services/swx_proxy/servicers/test"
	"magma/feg/gateway/services/swx_proxy/test_init"

	"github.com/stretchr/testify/assert"
)

func TestSwxProxyClient(t *testing.T) {
	err := test_init.StartTestService(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	expectedUsername := test.BASE_IMSI
	expectedNumVectors := 5
	expectedAuthScheme := protos.AuthenticationScheme_EAP_AKA
	authReq := &protos.AuthenticationRequest{
		UserName:             expectedUsername,
		SipNumAuthVectors:    uint32(expectedNumVectors),
		AuthenticationScheme: expectedAuthScheme,
	}

	// Authentication request - MAR
	authRes, err := swx_proxy.Authenticate(authReq)
	if err != nil {
		t.Fatalf("GRPC MAR Error: %v", err)
		return
	}
	t.Logf("GRPC MAA: %#+v", *authRes)
	assert.Equal(t, expectedUsername, authRes.GetUserName())
	assert.Equal(t, uint32(expectedNumVectors), authReq.GetSipNumAuthVectors())
	for i, v := range authRes.SipAuthVectors {
		assert.Equal(t, protos.AuthenticationScheme_EAP_AKA, v.GetAuthenticationScheme())
		assert.Equal(t, []byte(test.DefaultSIPAuthenticate+strconv.Itoa(int(i+14))), v.GetRandAutn())
		assert.Equal(t, []byte(test.DefaultSIPAuthorization), v.GetXres())
		assert.Equal(t, []byte(test.DefaultCK), v.GetConfidentialityKey())
		assert.Equal(t, []byte(test.DefaultIK), v.GetIntegrityKey())
	}

	// Registration request - SAR
	regReq := &protos.RegistrationRequest{
		UserName: expectedUsername,
	}
	regRes, err := swx_proxy.Register(regReq)
	if err != nil {
		t.Fatalf("GRPC SAR Error: %v", err)
		return
	}
	assert.Equal(t, &protos.RegistrationAnswer{}, regRes)
	t.Logf("GRPC SAA: %#+v", *regRes)

	// Test client error handling
	authRes, err = swx_proxy.Authenticate(nil)
	assert.EqualError(t, err, "Invalid AuthenticationRequest provided")
	assert.Nil(t, authRes)

	regRes, err = swx_proxy.Register(nil)
	assert.EqualError(t, err, "Invalid RegistrationRequest provided")
	assert.Nil(t, regRes)
}
