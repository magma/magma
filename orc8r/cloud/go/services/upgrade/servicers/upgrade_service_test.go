/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers_test

import (
	"fmt"
	"testing"

	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/protos"
	upgrade_protos "magma/orc8r/cloud/go/services/upgrade/protos"
	"magma/orc8r/cloud/go/services/upgrade/servicers"
	"magma/orc8r/cloud/go/test_utils"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
)

//------------------------------------------------------------------------------
// Release management tests
//------------------------------------------------------------------------------
func TestUpgradeService_CreateReleaseChannel(t *testing.T) {
	tableKey := "releases"
	ctx := context.Background()
	ds := test_utils.NewMockDatastore()
	srv := servicers.NewUpgradeService(ds)

	seedChannel := &upgrade_protos.ReleaseChannel{
		SupportedVersions: []string{"1.0.0-1", "1.0.1-2"},
	}
	_, err := srv.CreateReleaseChannel(ctx, &upgrade_protos.CreateOrUpdateReleaseChannelRequest{
		ChannelName: "stable",
		Channel:     seedChannel,
	})
	assert.NoError(t, err)
	assertDatastoreHasReleaseManagementRow(t, ds, tableKey, "stable", seedChannel)

	// Create already existing channel
	_, err = srv.CreateReleaseChannel(ctx, &upgrade_protos.CreateOrUpdateReleaseChannelRequest{
		ChannelName: "stable",
		Channel:     &upgrade_protos.ReleaseChannel{SupportedVersions: []string{"1.1.1-0"}},
	})
	assert.Error(t, err, "Release channel stable already exists", codes.AlreadyExists)
	assertDatastoreHasReleaseManagementRow(t, ds, tableKey, "stable", seedChannel)
}

func TestUpgradeService_GetReleaseChannel(t *testing.T) {
	tableKey := "releases"
	ctx := context.Background()
	ds := test_utils.NewMockDatastore()
	stableVers := []string{"1.0.0-1", "1.0.1-2", "1.1.3-4"}
	betaVers := []string{"1.2.0-0", "1.2.1-1"}
	fixtures := map[string]*upgrade_protos.ReleaseChannel{
		"stable": {SupportedVersions: stableVers},
		"beta":   {SupportedVersions: betaVers},
	}
	setupReleaseManagementFixtures(t, ds, tableKey, fixtures)

	srv := servicers.NewUpgradeService(ds)
	req := &upgrade_protos.GetReleaseChannelRequest{}

	req.ChannelName = "stable"
	actual, err := srv.GetReleaseChannel(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, *fixtures["stable"], *actual)

	req.ChannelName = "beta"
	actual, err = srv.GetReleaseChannel(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, *fixtures["beta"], *actual)

	req.ChannelName = "nochannel"
	actual, err = srv.GetReleaseChannel(ctx, req)
	assert.Error(t, err, "Error fetching release channel nochannel", codes.Aborted)
}

func TestUpgradeService_UpdateReleaseChannel(t *testing.T) {
	tableKey := "releases"
	ctx := context.Background()
	ds := test_utils.NewMockDatastore()
	fixtures := map[string]*upgrade_protos.ReleaseChannel{
		"stable": {SupportedVersions: []string{"1.0.0-1", "1.0.1-2", "1.1.3-4"}},
	}
	setupReleaseManagementFixtures(t, ds, tableKey, fixtures)

	srv := servicers.NewUpgradeService(ds)
	req := &upgrade_protos.CreateOrUpdateReleaseChannelRequest{}

	req.ChannelName = "stable"
	req.Channel = &upgrade_protos.ReleaseChannel{
		SupportedVersions: []string{"1.0.0-1", "1.0.1-2", "1.1.3-4", "1.1.4-1"},
	}
	_, err := srv.UpdateReleaseChannel(ctx, req)
	assert.NoError(t, err)
	assertDatastoreHasReleaseManagementRow(t, ds, tableKey, "stable", req.Channel)

	req.ChannelName = "nochannel"
	req.Channel = &upgrade_protos.ReleaseChannel{
		SupportedVersions: []string{"1.0.0-0"},
	}
	_, err = srv.UpdateReleaseChannel(ctx, req)
	assert.Error(t, err, "Can't update nonexistent channel", codes.FailedPrecondition)
}

func TestUpgradeService_DeleteReleaseChannel(t *testing.T) {
	tableKey := "releases"
	ctx := context.Background()
	ds := test_utils.NewMockDatastore()
	fixtures := map[string]*upgrade_protos.ReleaseChannel{
		"stable": {SupportedVersions: []string{}},
		"beta":   {SupportedVersions: []string{"1.2.0-0", "1.2.1-1"}},
	}
	setupReleaseManagementFixtures(t, ds, tableKey, fixtures)

	srv := servicers.NewUpgradeService(ds)
	req := &upgrade_protos.DeleteReleaseChannelRequest{}

	req.ChannelName = "stable"
	_, err := srv.DeleteReleaseChannel(ctx, req)
	assert.NoError(t, err)
	assertDatastoreHasReleaseManagementRow(t, ds, tableKey, "beta", fixtures["beta"])
	assertDatastoreDoesNotHaveReleaseManagementRow(t, ds, tableKey, "stable")

	// Delete non-existent channel
	req.ChannelName = "stable"
	_, err = srv.DeleteReleaseChannel(ctx, req)
	assert.Error(t, err, "Can't delete nonexistent channel")
	assertDatastoreHasReleaseManagementRow(t, ds, tableKey, "beta", fixtures["beta"])
}

func setupReleaseManagementFixtures(
	t *testing.T,
	ds datastore.Api,
	tableKey string,
	fixtures map[string]*upgrade_protos.ReleaseChannel,
) {
	for k, protoVal := range fixtures {
		marshaledProto, err := protos.MarshalIntern(protoVal)
		assert.NoError(t, err)
		ds.Put(tableKey, k, marshaledProto)
	}
}

func assertDatastoreHasReleaseManagementRow(t *testing.T, ds datastore.Api, tableKey string, key string, expectedVal *upgrade_protos.ReleaseChannel) {
	marshaledProto, _, err := ds.Get(tableKey, key)
	assert.NoError(t, err)
	actualProto := upgrade_protos.ReleaseChannel{}
	err = protos.Unmarshal(marshaledProto, &actualProto)
	assert.NoError(t, err)
	assert.Equal(t, *expectedVal, actualProto)
}

func assertDatastoreDoesNotHaveReleaseManagementRow(t *testing.T, ds datastore.Api, tableKey string, key string) {
	allKeys, err := ds.ListKeys(tableKey)
	assert.NoError(t, err)
	for _, k := range allKeys {
		if k == key {
			assert.Fail(
				t,
				fmt.Sprintf("Found table key %s which is not supposed to exist", key))
		}
	}
}

//------------------------------------------------------------------------------
// Tier versioning tests
//------------------------------------------------------------------------------

func TestUpgradeService_GetTiers(t *testing.T) {
	tableKey := "network_tierVersions"
	ctx := context.Background()
	ds := test_utils.NewMockDatastore()
	fixtures := map[string]*upgrade_protos.TierInfo{
		"t1": {Name: "t1", Version: "1.1.0-0", Images: []*upgrade_protos.ImageSpec{{Name: "fred", Order: 111}}},
		"t2": {Name: "t2", Version: "1.2.0-0", Images: []*upgrade_protos.ImageSpec{{Name: "bill", Order: 222}}},
	}
	setupTierVersioningFixtures(t, ds, tableKey, fixtures)

	srv := servicers.NewUpgradeService(ds)

	// No tier filter
	actual, err := srv.GetTiers(ctx, &upgrade_protos.GetTiersRequest{
		NetworkId:  "network",
		TierFilter: []string{},
	})
	assert.NoError(t, err)
	assert.Equal(t, fixtures, actual.GetTiers())

	// Tier filter
	actual, err = srv.GetTiers(ctx, &upgrade_protos.GetTiersRequest{
		NetworkId:  "network",
		TierFilter: []string{"t1"},
	})
	assert.NoError(t, err)
	assert.Equal(
		t,
		map[string]*upgrade_protos.TierInfo{"t1": fixtures["t1"]},
		actual.GetTiers())

	// Tier filter with a nonexistent tier
	actual, err = srv.GetTiers(ctx, &upgrade_protos.GetTiersRequest{
		NetworkId:  "network",
		TierFilter: []string{"t1", "t3"},
	})
	assert.NoError(t, err)
	assert.Equal(
		t,
		map[string]*upgrade_protos.TierInfo{"t1": fixtures["t1"]},
		actual.GetTiers(),
	)
}

func TestUpgradeService_CreateTiers(t *testing.T) {
	tableKey := "network_tierVersions"
	ctx := context.Background()
	ds := test_utils.NewMockDatastore()
	fixtures := map[string]*upgrade_protos.TierInfo{
		"t1": {Version: "1.1.0-0", Images: []*upgrade_protos.ImageSpec{{Name: "fred", Order: 111}}},
		"t2": {Version: "1.2.0-0", Images: []*upgrade_protos.ImageSpec{{Name: "bill", Order: 222}}},
	}
	setupTierVersioningFixtures(t, ds, tableKey, fixtures)

	srv := servicers.NewUpgradeService(ds)
	req := &upgrade_protos.CreateTierRequest{NetworkId: "network"}

	// Add existing tier should error
	req.TierId = "t1"
	req.TierInfo = &upgrade_protos.TierInfo{Name: "t1", Version: "1.3.0-0"}
	_, err := srv.CreateTier(ctx, req)
	assert.Error(t, err, "Can't create existing tier", codes.FailedPrecondition)

	req.TierId = "t3"
	req.TierInfo = &upgrade_protos.TierInfo{Name: "t3", Version: "1.3.0-0", Images: []*upgrade_protos.ImageSpec{{Name: "susie", Order: 333}}}

	_, err = srv.CreateTier(ctx, req)
	assert.NoError(t, err)

	expected := fixtures
	expected["t3"] = req.TierInfo
	assertDatastoreHasTierVersioningRows(t, ds, tableKey, expected)
}

func TestUpgradeService_UpdateTiers(t *testing.T) {
	tableKey := "network_tierVersions"
	ctx := context.Background()
	ds := test_utils.NewMockDatastore()
	fixtures := map[string]*upgrade_protos.TierInfo{
		"t1": {Name: "t1", Version: "1.1.0-0", Images: []*upgrade_protos.ImageSpec{{Name: "fred", Order: 111}}},
		"t2": {Name: "t2", Version: "1.2.0-0", Images: []*upgrade_protos.ImageSpec{{Name: "bill", Order: 222}}},
	}
	setupTierVersioningFixtures(t, ds, tableKey, fixtures)

	srv := servicers.NewUpgradeService(ds)
	req := &upgrade_protos.UpdateTierRequest{NetworkId: "network"}

	// Update nonexistent tier
	req.TierId = "t3"
	req.UpdatedTier = &upgrade_protos.TierInfo{Name: "t3", Version: "1.3.0-0"}
	_, err := srv.UpdateTier(ctx, req)
	assert.Error(t, err, "Can't update tier that doesn't exist", codes.FailedPrecondition)

	req.TierId = "t1"
	req.UpdatedTier = &upgrade_protos.TierInfo{Name: "t1v2", Version: "1.2.0-0", Images: []*upgrade_protos.ImageSpec{{Name: "sally", Order: 333}}}
	_, err = srv.UpdateTier(ctx, req)
	assert.NoError(t, err)

	expected := fixtures
	expected["t1"] = req.UpdatedTier
	assertDatastoreHasTierVersioningRows(t, ds, tableKey, expected)
}

func TestUpgradeService_DeleteTiers(t *testing.T) {
	tableKey := "network_tierVersions"
	ctx := context.Background()
	ds := test_utils.NewMockDatastore()
	fixtures := map[string]*upgrade_protos.TierInfo{
		"t1": {Name: "t1", Version: "1.1.0-0", Images: []*upgrade_protos.ImageSpec{{Name: "fred", Order: 111}}},
		"t2": {Name: "t2", Version: "1.2.0-0", Images: []*upgrade_protos.ImageSpec{{Name: "bill", Order: 222}}},
	}
	setupTierVersioningFixtures(t, ds, tableKey, fixtures)

	srv := servicers.NewUpgradeService(ds)
	req := &upgrade_protos.DeleteTierRequest{NetworkId: "network"}

	// Delete nonexistent tier
	req.TierIdToDelete = "t3"
	_, err := srv.DeleteTier(ctx, req)
	assert.Error(t, err, "Can't delete tier that doesn't exist", codes.FailedPrecondition)
	assertDatastoreHasTierVersioningRows(t, ds, tableKey, fixtures)

	req.TierIdToDelete = "t1"
	_, err = srv.DeleteTier(ctx, req)
	assert.NoError(t, err)

	expected := fixtures
	delete(expected, "t1")
	assertDatastoreHasTierVersioningRows(t, ds, tableKey, expected)
}

func setupTierVersioningFixtures(t *testing.T, ds datastore.Api, tableKey string, fixtures map[string]*upgrade_protos.TierInfo) {
	for k, v := range fixtures {
		marshalled, err := protos.MarshalIntern(v)
		assert.NoError(t, err)
		err = ds.Put(tableKey, k, marshalled)
		assert.NoError(t, err)
	}
}

func assertDatastoreHasTierVersioningRows(
	t *testing.T,
	ds datastore.Api,
	tableKey string,
	expectedRows map[string]*upgrade_protos.TierInfo,
) {
	for k, expected := range expectedRows {
		actual, _, err := ds.Get(tableKey, k)
		assert.NoError(t, err)
		actualUnmarshaled := &upgrade_protos.TierInfo{}
		err = protos.Unmarshal(actual, actualUnmarshaled)
		assert.Equal(t, *expected, *actualUnmarshaled)
	}
}
