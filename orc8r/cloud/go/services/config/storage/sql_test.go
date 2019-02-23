/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage_test

import (
	"database/sql"
	"errors"
	"testing"

	"magma/orc8r/cloud/go/services/config/storage"

	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestSqlConfigStorage_GetConfig(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening stub DB conn: %s", err)
	}
	defer db.Close()

	// happy path
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectQuery("SELECT value, version FROM network_configurations").
		WithArgs("type", "key").
		WillReturnRows(
			sqlmock.NewRows([]string{"value", "version"}).
				AddRow([]byte("value"), 42),
		)
	mock.ExpectCommit()

	store := storage.NewSqlConfigurationStorage(db)
	actual, err := store.GetConfig("network", "type", "key")
	assert.NoError(t, err)
	assert.Equal(t, &storage.ConfigValue{Value: []byte("value"), Version: 42}, actual)
	assert.NoError(t, mock.ExpectationsWereMet())

	// error case
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectQuery("SELECT value, version FROM network_configurations").
		WithArgs("type2", "key2").
		WillReturnError(errors.New("Mock select error"))
	mock.ExpectRollback()

	_, err = store.GetConfig("network", "type2", "key2")
	assert.Error(t, err)
	assert.Equal(t, "Mock select error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSqlConfigStorage_GetConfigs(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening stub DB conn: %s", err)
	}
	defer db.Close()
	store := storage.NewSqlConfigurationStorage(db)

	// no filter
	_, err = store.GetConfigs("network", &storage.FilterCriteria{})
	assert.EqualError(t, err, "At least one field of filter criteria must be specified")

	// happy path type only
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectQuery("SELECT type, key, value, version FROM network_configurations").
		WithArgs("type").
		WillReturnRows(
			sqlmock.NewRows([]string{"type", "key", "value", "version"}).
				AddRow("type", "key1", []byte("value1"), 1).
				AddRow("type", "key2", []byte("value2"), 2),
		)
	mock.ExpectCommit()

	actual, err := store.GetConfigs("network", &storage.FilterCriteria{Type: "type"})
	assert.NoError(t, err)
	expected := map[storage.TypeAndKey]*storage.ConfigValue{
		{Type: "type", Key: "key1"}: {Value: []byte("value1"), Version: 1},
		{Type: "type", Key: "key2"}: {Value: []byte("value2"), Version: 2},
	}
	assert.Equal(t, expected, actual)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	// happy path key only
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectQuery("SELECT type, key, value, version FROM network_configurations").
		WithArgs("key").
		WillReturnRows(
			sqlmock.NewRows([]string{"type", "key", "value", "version"}).
				AddRow("type1", "key", []byte("value1"), 1).
				AddRow("type2", "key", []byte("value2"), 2),
		)
	mock.ExpectCommit()

	actual, err = store.GetConfigs("network", &storage.FilterCriteria{Key: "key"})
	assert.NoError(t, err)
	expected = map[storage.TypeAndKey]*storage.ConfigValue{
		{Type: "type1", Key: "key"}: {Value: []byte("value1"), Version: 1},
		{Type: "type2", Key: "key"}: {Value: []byte("value2"), Version: 2},
	}
	assert.Equal(t, expected, actual)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	// happy path type and key
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectQuery("SELECT type, key, value, version FROM network_configurations").
		WithArgs("type", "key").
		WillReturnRows(
			sqlmock.NewRows([]string{"type", "key", "value", "version"}).
				AddRow("type", "key", []byte("value1"), 1),
		)
	mock.ExpectCommit()

	actual, err = store.GetConfigs("network", &storage.FilterCriteria{Type: "type", Key: "key"})
	assert.NoError(t, err)
	expected = map[storage.TypeAndKey]*storage.ConfigValue{
		{Type: "type", Key: "key"}: {Value: []byte("value1"), Version: 1},
	}
	assert.Equal(t, expected, actual)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	// Error query
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectQuery("SELECT type, key, value, version FROM network_configurations").
		WithArgs("type", "key").
		WillReturnError(errors.New("Mock query error"))
	mock.ExpectRollback()

	_, err = store.GetConfigs("network", &storage.FilterCriteria{Type: "type", Key: "key"})
	assert.EqualError(t, err, "Mock query error")
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSqlConfigStorage_ListKeysForType(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening stub DB conn: %s", err)
	}
	defer db.Close()
	store := storage.NewSqlConfigurationStorage(db)

	// happy path
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectQuery("SELECT key FROM network_configurations").
		WithArgs("type").
		WillReturnRows(
			sqlmock.NewRows([]string{"key"}).AddRow("key1").AddRow("key2"),
		)
	mock.ExpectCommit()

	actual, err := store.ListKeysForType("network", "type")
	assert.NoError(t, err)
	assert.Equal(t, []string{"key1", "key2"}, actual)
	assert.NoError(t, mock.ExpectationsWereMet())

	// error
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectQuery("SELECT key FROM network_configurations").
		WithArgs("type").
		WillReturnError(errors.New("Mock query error"))
	mock.ExpectRollback()

	_, err = store.ListKeysForType("network", "type")
	assert.EqualError(t, err, "Mock query error")
}

func TestSqlConfigStorage_CreateConfig(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening stub DB conn: %s", err)
	}
	defer db.Close()
	store := storage.NewSqlConfigurationStorage(db)

	// config exists already
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectQuery("SELECT 1 FROM network_configurations").
		WithArgs("type", "key").
		WillReturnRows(sqlmock.NewRows([]string{"result"}).AddRow(1))
	mock.ExpectRollback()

	err = store.CreateConfig("network", "type", "key", []byte("value"))
	assert.EqualError(t, err, "Creating already existing config with type type and key key")
	assert.NoError(t, mock.ExpectationsWereMet())

	// happy path
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectQuery("SELECT 1 FROM network_configurations").
		WithArgs("type", "key").
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO network_configurations").
		WithArgs("type", "key", []byte("value")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = store.CreateConfig("network", "type", "key", []byte("value"))
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// error in does exist query
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectQuery("SELECT 1 FROM network_configurations").
		WithArgs("type", "key").
		WillReturnError(errors.New("Mock query error"))
	mock.ExpectRollback()

	err = store.CreateConfig("network", "type", "key", []byte("value"))
	assert.EqualError(t, err, "Mock query error")
	assert.NoError(t, mock.ExpectationsWereMet())

	// error in insert
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectQuery("SELECT 1 FROM network_configurations").
		WithArgs("type", "key").
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO network_configurations").
		WithArgs("type", "key", []byte("value")).
		WillReturnError(errors.New("Mock exec error"))
	mock.ExpectRollback()

	err = store.CreateConfig("network", "type", "key", []byte("value"))
	assert.EqualError(t, err, "Mock exec error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSqlConfigStorage_UpdateConfig(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening stub DB conn: %s", err)
	}
	defer db.Close()
	store := storage.NewSqlConfigurationStorage(db)

	// config DNE
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectQuery("SELECT value, version FROM network_configurations").
		WithArgs("type", "key").
		WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	err = store.UpdateConfig("network", "type", "key", []byte("value"))
	assert.EqualError(t, err, "Updating nonexistent config with type type and key key")
	assert.NoError(t, mock.ExpectationsWereMet())

	// happy path
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectQuery("SELECT value, version FROM network_configurations").
		WithArgs("type", "key").
		WillReturnRows(sqlmock.NewRows([]string{"value", "version"}).AddRow([]byte("oldValue"), 1))
	mock.ExpectExec("UPDATE network_configurations").
		WithArgs([]byte("value"), 2, "type", "key").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = store.UpdateConfig("network", "type", "key", []byte("value"))
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// error in select
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectQuery("SELECT value, version FROM network_configurations").
		WithArgs("type", "key").
		WillReturnError(errors.New("Mock query error"))
	mock.ExpectRollback()

	err = store.UpdateConfig("network", "type", "key", []byte("value"))
	assert.EqualError(t, err, "Mock query error")
	assert.NoError(t, mock.ExpectationsWereMet())

	// error in update
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectQuery("SELECT value, version FROM network_configurations").
		WithArgs("type", "key").
		WillReturnRows(sqlmock.NewRows([]string{"value", "version"}).AddRow([]byte("oldValue"), 1))
	mock.ExpectExec("UPDATE network_configurations").
		WithArgs([]byte("value"), 2, "type", "key").
		WillReturnError(errors.New("Mock exec error"))
	mock.ExpectRollback()

	err = store.UpdateConfig("network", "type", "key", []byte("value"))
	assert.EqualError(t, err, "Mock exec error")
	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestSqlConfigStorage_DeleteConfig(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening stub DB conn: %s", err)
	}
	defer db.Close()
	store := storage.NewSqlConfigurationStorage(db)

	// config DNE
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectQuery("SELECT 1 FROM network_configurations").
		WithArgs("type", "key").
		WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	err = store.DeleteConfig("network", "type", "key")
	assert.EqualError(t, err, "Deleting nonexistent config with type type and key key")
	assert.NoError(t, mock.ExpectationsWereMet())

	// happy path
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectQuery("SELECT 1 FROM network_configurations").
		WithArgs("type", "key").
		WillReturnRows(sqlmock.NewRows([]string{"result"}).AddRow(1))
	mock.ExpectExec("DELETE FROM network_configurations").
		WithArgs("type", "key").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = store.DeleteConfig("network", "type", "key")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// error in does exist query
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectQuery("SELECT 1 FROM network_configurations").
		WithArgs("type", "key").
		WillReturnError(errors.New("Mock query error"))
	mock.ExpectRollback()

	err = store.DeleteConfig("network", "type", "key")
	assert.EqualError(t, err, "Mock query error")
	assert.NoError(t, mock.ExpectationsWereMet())

	// error in delete
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectQuery("SELECT 1 FROM network_configurations").
		WithArgs("type", "key").
		WillReturnRows(sqlmock.NewRows([]string{"result"}).AddRow(1))
	mock.ExpectExec("DELETE FROM network_configurations").
		WithArgs("type", "key").
		WillReturnError(errors.New("Mock exec error"))
	mock.ExpectRollback()

	err = store.DeleteConfig("network", "type", "key")
	assert.EqualError(t, err, "Mock exec error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSqlConfigStorage_DeleteConfigs(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening stub DB conn: %s", err)
	}
	defer db.Close()
	store := storage.NewSqlConfigurationStorage(db)

	// no filter
	_, err = store.GetConfigs("network", &storage.FilterCriteria{})
	assert.EqualError(t, err, "At least one field of filter criteria must be specified")

	// happy path type only
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectExec("DELETE FROM network_configurations").
		WithArgs("type").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = store.DeleteConfigs("network", &storage.FilterCriteria{Type: "type"})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// happy path key only
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectExec("DELETE FROM network_configurations").
		WithArgs("key").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = store.DeleteConfigs("network", &storage.FilterCriteria{Key: "key"})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// happy path type and key
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectExec("DELETE FROM network_configurations").
		WithArgs("type", "key").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = store.DeleteConfigs("network", &storage.FilterCriteria{Type: "type", Key: "key"})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Error query
	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectExec("DELETE FROM network_configurations").
		WithArgs("type", "key").
		WillReturnError(errors.New("Mock query error"))
	mock.ExpectRollback()

	err = store.DeleteConfigs("network", &storage.FilterCriteria{Type: "type", Key: "key"})
	assert.EqualError(t, err, "Mock query error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSqlConfigStorage_DeleteConfigsForNetwork(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening stub DB conn: %s", err)
	}
	defer db.Close()
	store := storage.NewSqlConfigurationStorage(db)

	mock.ExpectBegin()
	expectCreateTable(mock)
	mock.ExpectExec("DROP TABLE IF EXISTS network_configurations").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = store.DeleteConfigsForNetwork("network")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func expectCreateTable(mock sqlmock.Sqlmock) {
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS network_configurations").WillReturnResult(sqlmock.NewResult(1, 1))
}
