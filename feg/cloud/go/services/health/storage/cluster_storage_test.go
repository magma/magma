/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage_test

import (
	"testing"

	"magma/feg/cloud/go/services/health/storage"
	"magma/orc8r/cloud/go/test_utils"

	"github.com/stretchr/testify/assert"
)

const TestFegClusterID = "test_clusterID1"
const TestFegClusterID2 = "test_clusterID2"

func TestClusterStorage(t *testing.T) {
	store, err := storage.NewClusterStore(test_utils.NewMockDatastore())
	assert.NoError(t, err)

	exists, err := store.DoesKeyExist(TestFeGNetworkID, TestFegClusterID)
	assert.NoError(t, err)
	assert.False(t, exists)

	ret, err := store.GetClusterState(TestFeGNetworkID, TestFegClusterID)
	assert.EqualError(
		t,
		err,
		"Get ClusterState Error for network: test_networkID1, cluster: test_clusterID1; No record for query",
	)
	assert.Nil(t, ret)

	err = store.UpdateClusterState(TestFeGNetworkID, TestFegClusterID, TestFegLogicalID1)
	assert.NoError(t, err)

	exists, err = store.DoesKeyExist(TestFeGNetworkID, TestFegClusterID)
	assert.NoError(t, err)
	assert.True(t, exists)

	clusterState, err := store.GetClusterState(TestFeGNetworkID, "")
	assert.Error(t, err)
	assert.Nil(t, clusterState)

	clusterState, err = store.GetClusterState(TestFeGNetworkID, TestFegClusterID)
	assert.NoError(t, err)
	assert.Equal(t, TestFegLogicalID1, clusterState.ActiveGatewayLogicalId)

	err = store.UpdateClusterState(TestFeGNetworkID, TestFegClusterID2, TestFegLogicalID2)
	assert.NoError(t, err)

	clusterState2, err := store.GetClusterState(TestFeGNetworkID, TestFegClusterID2)
	assert.NoError(t, err)
	assert.Equal(t, TestFegLogicalID2, clusterState2.ActiveGatewayLogicalId)

	err = store.UpdateClusterState(TestFeGNetworkID, TestFegClusterID, TestFegLogicalID2)
	assert.NoError(t, err)

	updatedState, err := store.GetClusterState(TestFeGNetworkID, TestFegClusterID)
	assert.NoError(t, err)
	assert.Equal(t, TestFegLogicalID2, updatedState.ActiveGatewayLogicalId)
}
