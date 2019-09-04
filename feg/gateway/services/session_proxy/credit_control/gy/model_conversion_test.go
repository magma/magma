/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package gy_test

import (
	"testing"

	"magma/feg/gateway/services/session_proxy/credit_control/gy"
	"magma/lte/cloud/go/protos"

	"github.com/stretchr/testify/assert"
)

func TestRedirectServer_ToProto(t *testing.T) {
	var convertedRedirectServer *protos.RedirectServer = nil
	convertedRedirectServer = (&gy.RedirectServer{
		RedirectAddressType:   gy.IPV4Address,
		RedirectServerAddress: "www.magma.com",
	}).ToProto()

	assert.Equal(t, protos.RedirectServer_IPV4, convertedRedirectServer.RedirectAddressType)
	assert.Equal(t, "www.magma.com", convertedRedirectServer.RedirectServerAddress)

	var nilRedirectServer *gy.RedirectServer = nil
	convertedRedirectServer = nilRedirectServer.ToProto()

	assert.Equal(t, &protos.RedirectServer{}, convertedRedirectServer)
}
