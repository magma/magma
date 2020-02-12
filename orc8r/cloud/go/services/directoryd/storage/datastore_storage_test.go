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

	"magma/orc8r/cloud/go/services/directoryd/storage"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
)

func TestPersistenceService_GetRecord(t *testing.T) {
	db := test_utils.NewMockDatastore()
	store := storage.GetDirectorydPersistenceService(db)

	// Fixtures
	location1 := &protos.LocationRecord{
		Location: "gw1",
	}
	location2 := &protos.LocationRecord{
		Location: "gw2",
	}
	location_map := map[string]interface{}{
		"sid1": location1,
		"sid2": location2,
	}
	test_utils.SetupTestFixtures(
		t, db,
		protos.TableID_IMSI_TO_HWID.String(),
		location_map,
		serializeLocationRecord,
	)

	actual, err := store.GetRecord(protos.TableID_IMSI_TO_HWID, "sid1")
	assert.NoError(t, err)
	assert.Equal(t, location1, actual)

	actual, err = store.GetRecord(protos.TableID_IMSI_TO_HWID, "sid2")
	assert.NoError(t, err)
	assert.Equal(t, location2, actual)

	_, err = store.GetRecord(protos.TableID_IMSI_TO_HWID, "sid3")
	assert.Error(t, err)
}

func TestPersistenceService_DeleteRecord(t *testing.T) {
	db := test_utils.NewMockDatastore()
	store := storage.GetDirectorydPersistenceService(db)

	// Fixtures
	location1 := &protos.LocationRecord{
		Location: "gw1",
	}
	location2 := &protos.LocationRecord{
		Location: "gw2",
	}
	location_map := map[string]interface{}{
		"sid1": location1,
		"sid2": location2,
	}
	test_utils.SetupTestFixtures(
		t, db,
		protos.TableID_IMSI_TO_HWID.String(),
		location_map,
		serializeLocationRecord,
	)

	err := store.DeleteRecord(protos.TableID_IMSI_TO_HWID, "sid1")
	assert.NoError(t, err)

	_, err = store.GetRecord(protos.TableID_IMSI_TO_HWID, "sid1")
	assert.Error(t, err)

	err = store.DeleteRecord(protos.TableID_IMSI_TO_HWID, "sid1")
	assert.Error(t, err)
}

func TestPersistenceService_UpdateOrCreateRecord(t *testing.T) {
	db := test_utils.NewMockDatastore()
	store := storage.GetDirectorydPersistenceService(db)

	// Fixtures
	location1 := &protos.LocationRecord{
		Location: "gw1",
	}
	location2 := &protos.LocationRecord{
		Location: "gw2",
	}
	location_map := map[string]interface{}{
		"sid1": location1,
		"sid2": location2,
	}
	test_utils.SetupTestFixtures(
		t, db,
		protos.TableID_IMSI_TO_HWID.String(),
		location_map,
		serializeLocationRecord,
	)

	// Add a new record
	location3 := &protos.LocationRecord{
		Location: "gw3",
	}
	location_map["sid3"] = location3
	err := store.UpdateOrCreateRecord(protos.TableID_IMSI_TO_HWID, "sid3", location3)
	assert.NoError(t, err)
	assertDatastoreWritesSucceeded(t, db, location_map)

	// Update an existing flow
	location2.Location = "gw4"
	location_map["sid2"] = location2
	err = store.UpdateOrCreateRecord(protos.TableID_IMSI_TO_HWID, "sid2", location2)
	assert.NoError(t, err)
	assertDatastoreWritesSucceeded(t, db, location_map)

	// Add location from scratch
	location_map = map[string]interface{}{
		"sid1": location1,
	}
	err = store.UpdateOrCreateRecord(protos.TableID_IMSI_TO_HWID, "sid1", location1)
	assert.NoError(t, err)
	assertDatastoreWritesSucceeded(t, db, location_map)
}

func assertDatastoreWritesSucceeded(t *testing.T, store *test_utils.MockDatastore, location_map map[string]interface{}) {
	test_utils.AssertDatastoreHasRows(
		t, store,
		protos.TableID_IMSI_TO_HWID.String(),
		location_map,
		deserializeLocationRecord,
	)
}

func serializeLocationRecord(location interface{}) ([]byte, error) {
	locationCasted, ok := location.(*protos.LocationRecord)
	if !ok {
		return nil, errors.New("Expected *protos.LocationRecord")
	}
	return protos.MarshalIntern(locationCasted)
}

func deserializeLocationRecord(marshaled []byte) (interface{}, error) {
	ret := &protos.LocationRecord{}
	err := protos.Unmarshal(marshaled, ret)
	return ret, err
}
