/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package test_utils

import (
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/stretchr/testify/assert"
)

func RegisterNetwork(t *testing.T, networkID string, networkName string) {
	err := configurator.CreateNetwork(
		configurator.Network{
			ID:   networkID,
			Name: networkName,
		})
	assert.NoError(t, err)
}

func RegisterGateway(t *testing.T, networkID string, gatewayID string) {
	gw := configurator.NetworkEntity{
		Key:  gatewayID,
		Type: orc8r.MagmadGatewayType,
	}
	_, err := configurator.CreateEntity(networkID, gw)
	assert.NoError(t, err)
}
