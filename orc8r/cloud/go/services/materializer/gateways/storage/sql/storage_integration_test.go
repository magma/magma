/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package sql_test

import (
	"database/sql"
	"testing"

	"magma/orc8r/cloud/go/services/config/registry"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage"
	storage_sql "magma/orc8r/cloud/go/services/materializer/gateways/storage/sql"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage/test_utils"

	"github.com/stretchr/testify/assert"
)

// Integration test for sql storage impl which does some basic workflow tests
// on an in-memory sqlite3 DB.
func TestSqlGatewayViewStorage_Integration(t *testing.T) {
	// Use an in-memory sqlite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Could not initialize sqlite DB: %s", err)
	}

	registry.ClearRegistryForTesting()
	registry.RegisterConfigManager(test_utils.NewConfig1Manager())
	registry.RegisterConfigManager(test_utils.NewConfig2Manager())

	store := storage_sql.NewSqlGatewayViewStorage(db)

	// Empty database contract
	actual, err := store.GetGatewayViewsForNetwork("network")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(actual))

	actual, err = store.GetGatewayViews("network", []string{"gw1", "gw2"})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(actual))

	err = store.DeleteGatewayViews("network", []string{"gw1", "gw2"})
	assert.NoError(t, err)

	// Create 2 gateways
	updates := map[string]*storage.GatewayUpdateParams{
		"gw1": {
			NewStatus: getMockStatusProto(t, "gw1"),
			NewRecord: getMockRecordProto(t, "hw1"),
			NewConfig: getMockConfigs(),
			Offset:    42,
		},
		"gw2": {
			NewRecord: getMockRecordProto(t, "hw2"),
			Offset:    43,
		},
	}
	err = store.UpdateOrCreateGatewayViews("network", updates)
	assert.NoError(t, err)

	// Read back the gateways
	actual, err = store.GetGatewayViewsForNetwork("network")
	assert.NoError(t, err)
	expected := map[string]*storage.GatewayState{
		"gw1": {
			GatewayID: "gw1",
			Status:    getMockStatusProto(t, "gw1"),
			Record:    getMockRecordProto(t, "hw1"),
			Config:    getMockConfigs(),
			Offset:    42,
		},
		"gw2": {
			GatewayID: "gw2",
			Record:    getMockRecordProto(t, "hw2"),
			Config:    map[string]interface{}{},
			Offset:    43,
		},
	}
	assert.Equal(t, expected, actual)

	actual, err = store.GetGatewayViews("network", []string{"gw1", "gw2"})
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
	actual, err = store.GetGatewayViews("network", []string{"gw1"})
	assert.NoError(t, err)
	assert.Equal(t, map[string]*storage.GatewayState{"gw1": expected["gw1"]}, actual)
	actual, err = store.GetGatewayViews("network", []string{"gw2"})
	assert.NoError(t, err)
	assert.Equal(t, map[string]*storage.GatewayState{"gw2": expected["gw2"]}, actual)

	// Update gateways - cover config update and deletion for gw1,
	// config addition for gw2
	// Also create a new gateway gw3
	finalGw1Configs := map[string]interface{}{
		test_utils.NewConfig2Manager().GetConfigType(): &test_utils.Conf2{
			Value1: []string{"updated1", "updated2"},
			Value2: -100,
		},
	}
	updates = map[string]*storage.GatewayUpdateParams{
		"gw1": {
			NewConfig:       finalGw1Configs,
			ConfigsToDelete: []string{test_utils.NewConfig1Manager().GetConfigType()},
			Offset:          100,
		},
		"gw2": {
			NewConfig: getMockConfigs(),
			NewStatus: getMockStatusProto(t, "gw2"),
			Offset:    101,
		},
		"gw3": {
			Offset: 102,
		},
	}
	err = store.UpdateOrCreateGatewayViews("network", updates)
	assert.NoError(t, err)

	// Read back gateways
	actual, err = store.GetGatewayViewsForNetwork("network")
	assert.NoError(t, err)
	expected = map[string]*storage.GatewayState{
		"gw1": {
			GatewayID: "gw1",
			Status:    getMockStatusProto(t, "gw1"),
			Record:    getMockRecordProto(t, "hw1"),
			Config:    finalGw1Configs,
			Offset:    100,
		},
		"gw2": {
			GatewayID: "gw2",
			Status:    getMockStatusProto(t, "gw2"),
			Record:    getMockRecordProto(t, "hw2"),
			Config:    getMockConfigs(),
			Offset:    101,
		},
		"gw3": {
			GatewayID: "gw3",
			Config:    map[string]interface{}{},
			Offset:    102,
		},
	}
	assert.Equal(t, expected, actual)

	// Delete gw1, gw2
	err = store.DeleteGatewayViews("network", []string{"gw1", "gw2"})
	assert.NoError(t, err)

	// Read back gateways
	actual, err = store.GetGatewayViewsForNetwork("network")
	assert.NoError(t, err)
	expected = map[string]*storage.GatewayState{
		"gw3": {
			GatewayID: "gw3",
			Config:    map[string]interface{}{},
			Offset:    102,
		},
	}
	assert.Equal(t, expected, actual)

	actual, err = store.GetGatewayViews("network", []string{"gw1", "gw2"})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(actual))
}
