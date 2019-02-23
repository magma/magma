/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage

import "magma/orc8r/cloud/go/protos"

// Wraps a stored mconfig and the offset of the stream element that it was
// last computed from
type StoredMconfig struct {
	// ID of the network the gateway this mconfig configures belongs to
	NetworkId string

	// ID of the gateway this mconfig configures
	GatewayId string

	// Mconfig value
	Mconfig *protos.GatewayConfigs

	// The offset of the last stream element that this mconfig was calculated
	// from
	Offset int64
}

// Specifies an update to an existing mconfig or a creation of a new mconfig.
// All fields are mandatory.
type MconfigUpdateCriteria struct {
	// ID of the gateway this update applies to
	GatewayId string

	// New mconfig value for this gateway
	NewMconfig *protos.GatewayConfigs

	// The offset of the stream element that this update was created from
	Offset int64
}

type MconfigStorage interface {
	// Retrieve an mconfig for a specific gateway
	// Returns nil if no such mconfig exists
	GetMconfig(networkId string, gatewayId string) (*StoredMconfig, error)

	// Retrieve mconfigs for a list of gateways inside a network
	// The returned map will be keyed by gateway ID
	GetMconfigs(networkId string, gatewayIds []string) (map[string]*StoredMconfig, error)

	// Store new mconfigs. This will replace existing mconfigs and create new
	// mconfigs for previously unstored gateway IDs.
	CreateOrUpdateMconfigs(networkId string, updates []*MconfigUpdateCriteria) error

	// Delete mconfigs for a list of gateways inside a network
	DeleteMconfigs(networkId string, gatewayIds []string) error
}
