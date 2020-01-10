/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage

import "magma/feg/cloud/go/protos"

// HealthStorage defines a persistence interface for the health service
// Implementors are expected to manage health updates given as a HealthStats object
// uniquely keyed by a string tuple of (networkID, gatewayID).
type HealthStorage interface {
	// Fetch HealthStats object for a given networkID and gatewayID. If no such value exists,
	// this will return nil and the associated error
	GetHealth(networkID string, gatewayID string) (*protos.HealthStats, error)

	UpdateHealth(networkID string, gatewayID string, health *protos.HealthStats) error
}

// ClusterStorage defines a persistence interface for the health service to be used
// to achieve High Availability. Implementors are expected to manage cluster state updates given
// as a ClusterState object uniquely keyed by a string tuple of (networkID, clusterID, gatewayID)
// where gatewayID is the logicalID of the gateway.
type ClusterStorage interface {
	GetClusterState(networkID string, clusterID string) (*protos.ClusterState, error)

	UpdateClusterState(networkID string, clusterID string, logicalID string) error

	DoesKeyExist(networkID string, clusterID string) (bool, error)
}

// HealthBlobstore defines a storage interface for the health service. This
// interface defines create/update and read functionality while abstracting
// away any underlying storage transaction mechanics for clients.
type HealthBlobstore interface {
	GetHealth(networkID string, gatewayID string) (*protos.HealthStats, error)

	UpdateHealth(networkID string, gatewayID string, health *protos.HealthStats) error

	GetClusterState(networkID string, clusterID string) (*protos.ClusterState, error)

	UpdateClusterState(networkID string, clusterID string, logicalID string) error
}
