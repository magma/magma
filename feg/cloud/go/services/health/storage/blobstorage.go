/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package storage

import (
	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/services/health"
	"magma/feg/cloud/go/services/health/util"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/storage"

	"github.com/golang/glog"
)

// GetHealth fetches health status for the given networkID and gatewayID from
// the provided TransactionalBlobStorage
func GetHealth(store blobstore.TransactionalBlobStorage, networkID string, gatewayID string) (*fegprotos.HealthStats, error) {
	healthTypeAndKey := storage.TypeAndKey{
		Type: health.HealthStatusType,
		Key:  gatewayID,
	}
	healthBlob, err := store.Get(networkID, healthTypeAndKey)
	if err != nil {
		return nil, err
	}
	retHealth := &fegprotos.HealthStats{}
	err = protos.Unmarshal(healthBlob.Value, retHealth)
	return retHealth, err
}

// UpdateHealth updates the given gateway's health status in the provided
// TransactionalBlobStorage
func UpdateHealth(store blobstore.TransactionalBlobStorage, networkID string, gatewayID string, healthStats *fegprotos.HealthStats) error {
	healthBlob, err := util.HealthToBlob(gatewayID, healthStats)
	if err != nil {
		return err
	}
	return store.CreateOrUpdate(networkID, []blobstore.Blob{healthBlob})
}

// UpdateClusterState updates the given cluster's state in the provided
// TransactionalBlobStorage
func UpdateClusterState(store blobstore.TransactionalBlobStorage, networkID string, clusterID string, logicalID string) error {
	clusterBlob, err := util.ClusterToBlob(clusterID, logicalID)
	if err != nil {
		return err
	}
	return store.CreateOrUpdate(networkID, []blobstore.Blob{clusterBlob})
}

// GetClusterState retrieves the stored clusterState for the provided networkID
// and logicalID from the provided TransactionalBlobStorage. The clusterState
// is initialized if it doesn't already exist
func GetClusterState(store blobstore.TransactionalBlobStorage, networkID string, logicalID string) (*fegprotos.ClusterState, error) {
	keys := []string{networkID}
	filter := blobstore.SearchFilter{
		NetworkID: &networkID,
	}
	foundKeys, err := store.GetExistingKeys(keys, filter)
	if err != nil {
		return nil, err
	}
	if len(foundKeys) == 0 {
		glog.V(2).Infof("Initializing clusterState for networkID: %s with active: %s", networkID, logicalID)
		err = UpdateClusterState(store, networkID, networkID, logicalID)
		if err != nil {
			return nil, err
		}
	}
	clusterID := networkID
	clusterTypeAndKey := storage.TypeAndKey{
		Type: health.ClusterStatusType,
		Key:  clusterID,
	}
	clusterBlob, err := store.Get(networkID, clusterTypeAndKey)
	if err != nil {
		return nil, err
	}
	retClusterState := &fegprotos.ClusterState{}
	err = protos.Unmarshal(clusterBlob.Value, retClusterState)
	return retClusterState, err
}
