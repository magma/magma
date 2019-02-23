/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage

import (
	"fmt"
	"time"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/protos"
)

type clusterStorage struct {
	store datastore.Api
}

const ClusterStateTableName string = "clusterState"

func clusterStateTable(networkID string) string {
	return datastore.GetTableName(networkID, ClusterStateTableName)
}

// Create a new Cluster Store
func NewClusterStore(ds datastore.Api) (ClusterStorage, error) {
	cs := &clusterStorage{ds}
	if cs.store == nil {
		return nil, fmt.Errorf("Nil Cluster Store datastore")
	}
	return cs, nil
}

func (c *clusterStorage) UpdateClusterState(networkID, clusterID, logicalID string) error {
	if len(networkID) == 0 || len(clusterID) == 0 || len(logicalID) == 0 {
		return fmt.Errorf(
			"Empty cluster parameter: networkID: %s, clusterID: %s, logicalID: %s",
			networkID,
			clusterID,
			logicalID,
		)
	}
	clusterState := &fegprotos.ClusterState{
		ActiveGatewayLogicalId: logicalID,
		Time:                   uint64(time.Now().UnixNano()) / uint64(time.Millisecond),
	}

	marshaledClusterState, err := protos.MarshalIntern(clusterState)
	if err != nil {
		return fmt.Errorf("Cluster Store: Marshalling error for network: %s, cluster: %s; %s",
			networkID,
			clusterID,
			err,
		)
	}
	err = c.store.Put(clusterStateTable(networkID), clusterID, marshaledClusterState)
	if err != nil {
		return fmt.Errorf("Cluster Store Write error for network: %s, cluster: %s; %s",
			networkID,
			clusterID,
			err,
		)
	}
	return nil
}

func (c *clusterStorage) GetClusterState(networkID, clusterID string) (*fegprotos.ClusterState, error) {
	if len(networkID) == 0 || len(clusterID) == 0 {
		return nil, fmt.Errorf("Empty cluster parameter: networkID: %s, clusterID %s",
			networkID,
			clusterID,
		)
	}

	marshaledClusterState, _, err := c.store.Get(clusterStateTable(networkID), clusterID)
	if err != nil {
		return nil, fmt.Errorf("Get ClusterState Error for network: %s, cluster: %s; %s",
			networkID,
			clusterID,
			err,
		)
	}
	clusterState := new(fegprotos.ClusterState)
	err = protos.Unmarshal(marshaledClusterState, clusterState)
	if err != nil {
		return nil, fmt.Errorf("Cluster Store Unmarshaling Error: %s", err)
	}
	return clusterState, nil
}

func (c *clusterStorage) DoesKeyExist(networkID, clusterID string) (bool, error) {
	if len(networkID) == 0 || len(clusterID) == 0 {
		return false, fmt.Errorf("Empty cluster parameter: networkID: %s, clusterID %s",
			networkID,
			clusterID,
		)

	}
	return c.store.DoesKeyExist(clusterStateTable(networkID), clusterID)
}
