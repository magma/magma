/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package storage represents storage of materialized views of gateways
package storage

import (
	"magma/orc8r/cloud/go/protos"
	magmadprotos "magma/orc8r/cloud/go/services/magmad/protos"
)

// GatewayState represents the current state of a gateway, including
// information on configuration parameters, status, and record
type GatewayState struct {
	// ID of the gateway
	GatewayID string
	// Configuration of the gateway, represented as a map from configuration types
	// to configuration objects
	Config map[string]interface{}
	// Status of the gateway
	Status *protos.GatewayStatus
	// Gateway record
	Record *magmadprotos.AccessGatewayRecord
	// Offset is the stream offset of this state; this ensures sequential consistency of
	// state processing
	Offset int64
}

// GatewayUpdateParams contains information from an update to a gateway state. Each
// parameter is nullable, and only non-null parameters will be used to update the gateway
// state
type GatewayUpdateParams struct {
	// Only the keys in this NewConfig map will be used to update the config of the GatewayState
	NewConfig map[string]interface{}
	// Configurations to delete
	ConfigsToDelete []string
	// New status of the gateway
	NewStatus *protos.GatewayStatus
	// New gateway record
	NewRecord *magmadprotos.AccessGatewayRecord
	// Offset of the update. If this offset is less than the offset of the current gateway state,
	// then the update will be rejected
	Offset int64
}

// GatewayViewStorage is an interface for gateway materialized views specifying basic
// CRUD operations for the materialized view
type GatewayViewStorage interface {
	// Do initialization work for storage
	InitTables() error
	// Get the states of all gateways in this network
	GetGatewayViewsForNetwork(networkID string) (map[string]*GatewayState, error)
	// Get the state of a specific gateway
	GetGatewayViews(networkID string, gatewayIDs []string) (map[string]*GatewayState, error)
	// Update the state of the specified gateways, or create the gateways if their states
	// are not yet tracked
	UpdateOrCreateGatewayViews(networkID string, updates map[string]*GatewayUpdateParams) error
	// Delete materialized views of the gateways specified in gatewayIDs
	DeleteGatewayViews(networkID string, gatewayIDs []string) error
}
