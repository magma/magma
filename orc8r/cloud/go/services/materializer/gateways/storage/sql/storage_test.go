/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package sql_test

import (
	"fmt"
	"testing"

	magma_protos "magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/config/registry"
	magmadprotos "magma/orc8r/cloud/go/services/magmad/protos"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage"
	storage_sql "magma/orc8r/cloud/go/services/materializer/gateways/storage/sql"
	sql_protos "magma/orc8r/cloud/go/services/materializer/gateways/storage/sql/protos"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage/test_utils"

	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestSqlGatewayViewStorage_GetGatewayViewsForNetwork(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening stub DB conn: %s", err)
	}
	defer db.Close()

	registry.ClearRegistryForTesting()
	registry.RegisterConfigManager(test_utils.NewConfig1Manager())
	registry.RegisterConfigManager(test_utils.NewConfig2Manager())

	store := storage_sql.NewSqlGatewayViewStorage(db)

	// happy path
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectQuery("SELECT gateway_id, status, record, configs, \"offset\" FROM network_gateway_views").
		WithArgs().
		WillReturnRows(
			sqlmock.NewRows([]string{"gateway_id", "status", "record", "configs", "offset"}).
				AddRow("gw1", getMockStatusBytes(t, "gw1"), getMockRecordBytes(t, "hw1"), getMockConfigsBytes(t), 42).
				AddRow("gw2", getMockStatusBytes(t, "gw2"), getMockRecordBytes(t, "hw2"), getMockConfigsBytes(t), 43).
				AddRow("gw3", []byte{}, []byte{}, []byte{}, 0),
		)
	mock.ExpectCommit()

	actual, err := store.GetGatewayViewsForNetwork("network")
	assert.NoError(t, err)
	expected := map[string]*storage.GatewayState{
		"gw1": {
			GatewayID: "gw1",
			Record:    getMockRecordProto(t, "hw1"),
			Status:    getMockStatusProto(t, "gw1"),
			Config:    getMockConfigs(),
			Offset:    42,
		},
		"gw2": {
			GatewayID: "gw2",
			Record:    getMockRecordProto(t, "hw2"),
			Status:    getMockStatusProto(t, "gw2"),
			Config:    getMockConfigs(),
			Offset:    43,
		},
		"gw3": {
			GatewayID: "gw3",
			Config:    map[string]interface{}{},
		},
	}
	assert.Equal(t, expected, actual)
	assert.NoError(t, mock.ExpectationsWereMet())

	// SQL error
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectQuery("SELECT gateway_id, status, record, configs, \"offset\" FROM network_gateway_views").
		WithArgs().
		WillReturnError(fmt.Errorf("Error 42 foobar"))
	mock.ExpectRollback()

	_, err = store.GetGatewayViewsForNetwork("network")
	assert.Error(t, err, "Storage query error: Error 42 foobar")
	assert.NoError(t, mock.ExpectationsWereMet())

	registry.ClearRegistryForTesting()
}

func TestSqlGatewayViewStorage_GetGatewayViews(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening stub DB conn: %s", err)
	}
	defer db.Close()

	registry.ClearRegistryForTesting()
	registry.RegisterConfigManager(test_utils.NewConfig1Manager())
	registry.RegisterConfigManager(test_utils.NewConfig2Manager())

	store := storage_sql.NewSqlGatewayViewStorage(db)

	// happy path - but have DB return miss gw2
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectQuery("SELECT gateway_id, status, record, configs, \"offset\" FROM network_gateway_views").
		WithArgs("gw1", "gw2", "gw3").
		WillReturnRows(
			sqlmock.NewRows([]string{"gateway_id", "status", "record", "configs", "offset"}).
				AddRow("gw1", getMockStatusBytes(t, "gw1"), getMockRecordBytes(t, "hw1"), getMockConfigsBytes(t), 42).
				AddRow("gw3", []byte{}, []byte{}, []byte{}, 0),
		)
	mock.ExpectCommit()

	actual, err := store.GetGatewayViews("network", []string{"gw1", "gw2", "gw3"})
	assert.NoError(t, err)
	expected := map[string]*storage.GatewayState{
		"gw1": {
			GatewayID: "gw1",
			Record:    getMockRecordProto(t, "hw1"),
			Status:    getMockStatusProto(t, "gw1"),
			Config:    getMockConfigs(),
			Offset:    42,
		},
		"gw3": {
			GatewayID: "gw3",
			Config:    map[string]interface{}{},
		},
	}
	assert.Equal(t, expected, actual)
	assert.NoError(t, mock.ExpectationsWereMet())

	// SQL error
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectQuery("SELECT gateway_id, status, record, configs, \"offset\" FROM network_gateway_views").
		WithArgs("gw1", "gw2", "gw3").
		WillReturnError(fmt.Errorf("Error 42 foobar"))
	mock.ExpectRollback()

	_, err = store.GetGatewayViews("network", []string{"gw1", "gw2", "gw3"})
	assert.Error(t, err, "Storage query error: Error 42 foobar")
	assert.NoError(t, mock.ExpectationsWereMet())

	registry.ClearRegistryForTesting()
}

func TestSqlGatewayViewStorage_UpdateOrCreateGatewayViews_HappyPath(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening stub DB conn: %s", err)
	}
	defer db.Close()

	registry.ClearRegistryForTesting()
	registry.RegisterConfigManager(test_utils.NewConfig1Manager())
	registry.RegisterConfigManager(test_utils.NewConfig2Manager())

	store := storage_sql.NewSqlGatewayViewStorage(db)

	// Happy path - create with full criteria, create with only 1 field,
	// update with full criteria, update with only 1 field
	mock.ExpectBegin()
	expectCreateTable(mock)

	// Setup some fixture data first
	// Loaded gw3 configs will only have Config1 type
	selectedGw3Configs := getMockConfigsProto(t)
	delete(selectedGw3Configs.Configs, test_utils.NewConfig2Manager().GetConfigType())
	selectedGw3ConfigsBytes, err := magma_protos.MarshalIntern(selectedGw3Configs)
	assert.NoError(t, err)

	// We'll update gw3 to add the Config2 type and delete the Config1 type
	finalGw3Configs := map[string]interface{}{
		test_utils.NewConfig2Manager().GetConfigType(): &test_utils.Conf2{
			Value1: []string{"updated1", "updated2"},
			Value2: -100,
		},
	}
	finalGw3ConfigsBytes := marshalMockConfigs(t, finalGw3Configs)

	// We expect a select first to get existing views
	// Return results only for the gateways which are being updated
	mock.ExpectQuery("SELECT gateway_id, status, record, configs, \"offset\" FROM network_gateway_views").
		WithArgs("gw1", "gw2", "gw3", "gw4").
		WillReturnRows(
			sqlmock.NewRows([]string{"gateway_id", "status", "record", "configs", "offset"}).
				AddRow("gw3", getMockStatusBytes(t, "gw3"), getMockRecordBytes(t, "hw3"), selectedGw3ConfigsBytes, 42).
				AddRow("gw4", getMockStatusBytes(t, "gw4"), getMockRecordBytes(t, "hw4"), getMockConfigsBytes(t), 43),
		)

	// Upsert query duplicates args (placeholders for INSERT and UPDATE clauses)
	mock.ExpectExec("INSERT INTO network_gateway_views").
		WithArgs(
			"gw1", getMockStatusBytes(t, "gw1"), getMockRecordBytes(t, "hw1"), getMockConfigsBytes(t), 1,
			getMockStatusBytes(t, "gw1"), getMockRecordBytes(t, "hw1"), getMockConfigsBytes(t), 1, "gw1",
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO network_gateway_views").
		WithArgs(
			"gw2", getMockConfigsBytes(t), 2,
			getMockConfigsBytes(t), 2, "gw2",
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Cover add new config and config deletion for gw3
	mock.ExpectExec("INSERT INTO network_gateway_views").
		WithArgs(
			"gw3", getMockStatusBytes(t, "gw3"), getMockRecordBytes(t, "hw3"), finalGw3ConfigsBytes, 100,
			getMockStatusBytes(t, "gw3"), getMockRecordBytes(t, "hw3"), finalGw3ConfigsBytes, 100, "gw3",
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO network_gateway_views").
		WithArgs(
			"gw4", getMockStatusBytes(t, "gw4"), 101,
			getMockStatusBytes(t, "gw4"), 101, "gw4",
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	params := map[string]*storage.GatewayUpdateParams{
		"gw1": {
			NewConfig: getMockConfigs(),
			NewRecord: getMockRecordProto(t, "hw1"),
			NewStatus: getMockStatusProto(t, "gw1"),
			Offset:    1,
		},
		"gw2": {
			NewConfig: getMockConfigs(),
			Offset:    2,
		},
		"gw3": {
			NewConfig:       finalGw3Configs, // we can just re-use this final expected value
			ConfigsToDelete: []string{test_utils.NewConfig1Manager().GetConfigType()},
			NewRecord:       getMockRecordProto(t, "hw3"),
			NewStatus:       getMockStatusProto(t, "gw3"),
			Offset:          100,
		},
		"gw4": {
			NewStatus: getMockStatusProto(t, "gw4"),
			Offset:    101,
		},
	}
	err = store.UpdateOrCreateGatewayViews("network", params)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSqlGatewayViewStorage_UpdateOrCreateGatewayViews_Errors(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening stub DB conn: %s", err)
	}
	defer db.Close()

	registry.ClearRegistryForTesting()
	registry.RegisterConfigManager(test_utils.NewConfig1Manager())
	registry.RegisterConfigManager(test_utils.NewConfig2Manager())

	store := storage_sql.NewSqlGatewayViewStorage(db)

	// Select error
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectQuery("SELECT gateway_id, status, record, configs, \"offset\" FROM network_gateway_views").
		WithArgs("gw1", "gw2").
		WillReturnError(fmt.Errorf("Error foobar 42"))
	mock.ExpectRollback()

	err = store.UpdateOrCreateGatewayViews("network", map[string]*storage.GatewayUpdateParams{
		"gw1": {Offset: 100},
		"gw2": {Offset: 101},
	})
	assert.Error(t, err, "Error loading existing gateway views: Error foobar 42")
	assert.NoError(t, mock.ExpectationsWereMet())

	// Upsert error on second upsert
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectQuery("SELECT gateway_id, status, record, configs, \"offset\" FROM network_gateway_views").
		WithArgs("gw1", "gw2").
		WillReturnRows(sqlmock.NewRows([]string{"gateway_id", "status", "record", "configs", "offset"}))
	mock.ExpectExec("INSERT INTO network_gateway_views").
		WithArgs(
			"gw1", getMockStatusBytes(t, "gw1"), 100,
			getMockStatusBytes(t, "gw1"), 100, "gw1",
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO network_gateway_views").
		WithArgs(
			"gw2", getMockStatusBytes(t, "gw2"), 101,
			getMockStatusBytes(t, "gw2"), 101, "gw2",
		).
		WillReturnError(fmt.Errorf("Error barbaz 43"))
	mock.ExpectRollback()

	err = store.UpdateOrCreateGatewayViews("network", map[string]*storage.GatewayUpdateParams{
		"gw1": {NewStatus: getMockStatusProto(t, "gw1"), Offset: 100},
		"gw2": {NewStatus: getMockStatusProto(t, "gw2"), Offset: 101},
	})
	assert.Error(t, err, "Error updating gateway gw2: Error barbaz 43")
	assert.NoError(t, mock.ExpectationsWereMet())

	registry.ClearRegistryForTesting()
}

func TestSqlGatewayViewStorage_DeleteGatewayViews(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening stub DB conn: %s", err)
	}
	defer db.Close()

	store := storage_sql.NewSqlGatewayViewStorage(db)

	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectExec("DELETE FROM network_gateway_views").WithArgs("gw1", "gw2").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = store.DeleteGatewayViews("network", []string{"gw1", "gw2"})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectExec("DELETE FROM network_gateway_views").WithArgs("gw3", "gw4").WillReturnError(fmt.Errorf("Error bazfoo 44"))
	mock.ExpectRollback()

	err = store.DeleteGatewayViews("network", []string{"gw3", "gw4"})
	assert.Error(t, err, "Error deleting gateway views: Error bazfoo 44")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func expectCreateTable(mock sqlmock.Sqlmock) {
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS network_gateway_views").WillReturnResult(sqlmock.NewResult(1, 1))
}

func getMockStatusBytes(t *testing.T, gwID string) []byte {
	_, ret := test_utils.GetMockStatus(t, gwID)
	return ret
}

func getMockStatusProto(t *testing.T, gwID string) *magma_protos.GatewayStatus {
	ret, _ := test_utils.GetMockStatus(t, gwID)
	return ret
}

func getMockRecordBytes(t *testing.T, hwID string) []byte {
	_, ret := test_utils.GetMockRecord(t, hwID)
	return ret
}

func getMockRecordProto(t *testing.T, hwID string) *magmadprotos.AccessGatewayRecord {
	ret, _ := test_utils.GetMockRecord(t, hwID)
	return ret
}

func getMockConfigsBytes(t *testing.T) []byte {
	_, ret := getMockViewConfigs(t)
	return ret
}

func getMockConfigsProto(t *testing.T) *sql_protos.ViewConfigs {
	ret, _ := getMockViewConfigs(t)
	return ret
}

func getMockViewConfigs(t *testing.T) (*sql_protos.ViewConfigs, []byte) {
	cfg1 := &test_utils.Conf1{Value1: 1, Value2: "foo", Value3: []byte("bar")}
	cfg2 := &test_utils.Conf2{Value1: []string{"foo", "bar"}, Value2: 1}

	cfg1Bytes, err := test_utils.NewConfig1Manager().MarshalConfig(cfg1)
	assert.NoError(t, err)
	cfg2Bytes, err := test_utils.NewConfig2Manager().MarshalConfig(cfg2)
	assert.NoError(t, err)

	ret := &sql_protos.ViewConfigs{
		Configs: map[string][]byte{
			test_utils.NewConfig1Manager().GetConfigType(): cfg1Bytes,
			test_utils.NewConfig2Manager().GetConfigType(): cfg2Bytes,
		},
	}
	retBytes, err := magma_protos.MarshalIntern(ret)
	assert.NoError(t, err)
	return ret, retBytes
}

func getMockConfigs() map[string]interface{} {
	cfg1 := &test_utils.Conf1{Value1: 1, Value2: "foo", Value3: []byte("bar")}
	cfg2 := &test_utils.Conf2{Value1: []string{"foo", "bar"}, Value2: 1}

	return map[string]interface{}{
		test_utils.NewConfig1Manager().GetConfigType(): cfg1,
		test_utils.NewConfig2Manager().GetConfigType(): cfg2,
	}
}

func marshalMockConfigs(t *testing.T, cfgs map[string]interface{}) []byte {
	viewCfgs := &sql_protos.ViewConfigs{Configs: map[string][]byte{}}
	for k, v := range cfgs {
		marshaledV, err := registry.MarshalConfig(k, v)
		assert.NoError(t, err)
		viewCfgs.Configs[k] = marshaledV
	}

	ret, err := magma_protos.MarshalIntern(viewCfgs)
	assert.NoError(t, err)
	return ret
}
