/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage_test

import (
	"errors"
	"testing"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/config/streaming/storage"
	storage_protos "magma/orc8r/cloud/go/services/config/streaming/storage/protos"
	"magma/orc8r/cloud/go/services/config/streaming/storage/test_protos"
	"magma/orc8r/cloud/go/test_utils"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestDatastoreMconfigStorage_GetMconfig(t *testing.T) {
	db := test_utils.NewMockDatastore()

	// Setup fixtures
	mcfg, _ := test_protos.GetDefaultMconfig(t)
	test_utils.SetupTestFixtures(
		t, db,
		storage.GetMconfigViewTableName("network"),
		map[string]interface{}{
			"gateway1": &storage_protos.StoredMconfig{Configs: mcfg, Offset: 1},
		},
		serializeStoredMconfig,
	)

	store := storage.NewDatastoreMconfigStorage(db)

	actual, err := store.GetMconfig("network", "gateway2")
	assert.NoError(t, err)
	assert.Nil(t, actual)

	actual, err = store.GetMconfig("network", "gateway1")
	assert.NoError(t, err)
	expected := &storage.StoredMconfig{NetworkId: "network", GatewayId: "gateway1", Mconfig: mcfg, Offset: 1}
	assert.Equal(t, expected, actual)
}

func TestDatastoreMconfigStorage_GetMconfigs(t *testing.T) {
	db := test_utils.NewMockDatastore()

	// Setup fixtures
	mcfg1, _ := test_protos.GetDefaultMconfig(t)
	mcfg2, _ := test_protos.GetMconfig(
		t,
		map[string]proto.Message{"foo": &test_protos.Config1{Field: "bar"}},
	)
	test_utils.SetupTestFixtures(
		t, db,
		storage.GetMconfigViewTableName("network"),
		map[string]interface{}{
			"gateway1": &storage_protos.StoredMconfig{Configs: mcfg1, Offset: 1},
			"gateway2": &storage_protos.StoredMconfig{Configs: mcfg2, Offset: 2},
		},
		serializeStoredMconfig,
	)

	store := storage.NewDatastoreMconfigStorage(db)

	actual, err := store.GetMconfigs("network", []string{"gateway1", "gateway2"})
	expected := map[string]*storage.StoredMconfig{
		"gateway1": {NetworkId: "network", GatewayId: "gateway1", Mconfig: mcfg1, Offset: 1},
		"gateway2": {NetworkId: "network", GatewayId: "gateway2", Mconfig: mcfg2, Offset: 2},
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	actual, err = store.GetMconfigs("network", []string{"gateway2"})
	expected = map[string]*storage.StoredMconfig{
		"gateway2": {NetworkId: "network", GatewayId: "gateway2", Mconfig: mcfg2, Offset: 2},
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// One key DNE
	actual, err = store.GetMconfigs("network", []string{"gateway2", "gateway3"})
	expected = map[string]*storage.StoredMconfig{
		"gateway2": {NetworkId: "network", GatewayId: "gateway2", Mconfig: mcfg2, Offset: 2},
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// All keys DNE
	actual, err = store.GetMconfigs("network", []string{"gateway3", "gateway4"})
	assert.NoError(t, err)
	assert.Equal(t, map[string]*storage.StoredMconfig{}, actual)
}

func TestDatastoreMconfigStorage_CreateOrUpdateMconfigs(t *testing.T) {
	db := test_utils.NewMockDatastore()

	// Setup fixtures
	mcfg1, _ := test_protos.GetDefaultMconfig(t)
	mcfg2, _ := test_protos.GetMconfig(
		t,
		map[string]proto.Message{"foo": &test_protos.Config1{Field: "bar"}},
	)
	test_utils.SetupTestFixtures(
		t, db,
		storage.GetMconfigViewTableName("network"),
		map[string]interface{}{
			"gateway1": &storage_protos.StoredMconfig{Configs: mcfg1, Offset: 1},
			"gateway2": &storage_protos.StoredMconfig{Configs: mcfg2, Offset: 2},
		},
		serializeStoredMconfig,
	)

	store := storage.NewDatastoreMconfigStorage(db)

	// Create 1, update 1
	err := store.CreateOrUpdateMconfigs(
		"network",
		[]*storage.MconfigUpdateCriteria{
			{GatewayId: "gateway1", NewMconfig: mcfg2, Offset: 4},
			{GatewayId: "gateway3", NewMconfig: mcfg1, Offset: 3},
		},
	)
	assert.NoError(t, err)
	test_utils.AssertDatastoreHasRows(
		t, db,
		storage.GetMconfigViewTableName("network"),
		map[string]interface{}{
			"gateway1": &storage_protos.StoredMconfig{Configs: mcfg2, Offset: 4},
			"gateway2": &storage_protos.StoredMconfig{Configs: mcfg2, Offset: 2},
			"gateway3": &storage_protos.StoredMconfig{Configs: mcfg1, Offset: 3},
		},
		deserializeStoredMconfig,
	)
}

func TestDatastoreMconfigStorage_DeleteMconfigs(t *testing.T) {
	db := test_utils.NewMockDatastore()

	// Setup fixtures
	mcfg1, _ := test_protos.GetDefaultMconfig(t)
	mcfg2, _ := test_protos.GetMconfig(
		t,
		map[string]proto.Message{"foo": &test_protos.Config1{Field: "bar"}},
	)
	test_utils.SetupTestFixtures(
		t, db,
		storage.GetMconfigViewTableName("network"),
		map[string]interface{}{
			"gateway1": &storage_protos.StoredMconfig{Configs: mcfg1, Offset: 1},
			"gateway2": &storage_protos.StoredMconfig{Configs: mcfg2, Offset: 2},
		},
		serializeStoredMconfig,
	)

	store := storage.NewDatastoreMconfigStorage(db)

	err := store.DeleteMconfigs("network", []string{"gateway1", "gateway3"})
	assert.NoError(t, err)
	test_utils.AssertDatastoreHasRows(
		t, db,
		storage.GetMconfigViewTableName("network"),
		map[string]interface{}{
			"gateway2": &storage_protos.StoredMconfig{Configs: mcfg2, Offset: 2},
		},
		deserializeStoredMconfig,
	)
}

func serializeStoredMconfig(storedMconfig interface{}) ([]byte, error) {
	casted, ok := storedMconfig.(*storage_protos.StoredMconfig)
	if !ok {
		return nil, errors.New("Expected protobuf *StoredMconfig")
	}
	return protos.MarshalIntern(casted)
}

func deserializeStoredMconfig(marshaled []byte) (interface{}, error) {
	ret := &storage_protos.StoredMconfig{}
	err := protos.Unmarshal(marshaled, ret)
	return ret, err
}
