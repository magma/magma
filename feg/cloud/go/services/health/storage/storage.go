/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package storage

import "magma/feg/cloud/go/protos"

// HealthBlobstore defines a storage interface for the health service. This
// interface defines create/update and read functionality while abstracting
// away any underlying storage transaction mechanics for clients.
type HealthBlobstore interface {
	GetHealth(networkID string, gatewayID string) (*protos.HealthStats, error)

	UpdateHealth(networkID string, gatewayID string, health *protos.HealthStats) error

	GetClusterState(networkID string, clusterID string) (*protos.ClusterState, error)

	UpdateClusterState(networkID string, clusterID string, logicalID string) error
}
