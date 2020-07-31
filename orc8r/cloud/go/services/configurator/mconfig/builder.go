/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package mconfig

import (
	"magma/orc8r/cloud/go/services/configurator/storage"

	"github.com/golang/protobuf/ptypes/any"
)

type ConfigsByKey map[string]*any.Any

// Builder creates a partial mconfig for a gateway within a network.
type Builder interface {
	// Build returns a partial mconfig containing the gateway configs for which
	// this builder is responsible.
	//
	// Parameters:
	//	- network	-- network containing the gateway
	//	- graph		-- entity graph associated with the gateway
	//	- gatewayID	-- HWID of the gateway
	Build(network *storage.Network, graph *storage.EntityGraph, gatewayID string) (ConfigsByKey, error)
}
