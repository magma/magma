/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package protos_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/feg/cloud/go/services/controller/protos"
)

func TestValidateGatewayConfig(t *testing.T) {
	config := protos.NewDefaultProtosGatewayConfig()
	err := protos.ValidateGatewayConfig(config)
	assert.NoError(t, err)

	err = protos.ValidateGatewayConfig(nil)
	assert.Error(t, err)
}

func TestValidateNetworkConfig(t *testing.T) {
	config := protos.NewDefaultProtosNetworkConfig()
	err := protos.ValidateNetworkConfig(config)
	assert.NoError(t, err)

	err = protos.ValidateNetworkConfig(nil)
	assert.Error(t, err)
}
