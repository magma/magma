/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"

	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/sql_utils"
	"magma/orc8r/cloud/go/storage"
)

type sqlConfigStorage struct {
	db *sql.DB
}

func NewSqlConfigurationStorage(db *sql.DB) ConfigurationStorage {
	return &sqlConfigStorage{db: db}
}

const tableName = "configurations"

func GetConfigTableName(networkId string) string {
	return datastore.GetTableName(networkId, tableName)
}

// This is a mega-hack before our inter-service streaming architecture is in
// place. Execute a CREATE TABLE IF NOT EXISTS on every query so we don't
// query a non-existent table.
func initTable(tx *sql.Tx, table string) error {
	queryFormat := `
		CREATE TABLE IF NOT EXISTS %s
		(
			type text NOT NULL,
			key text NOT NULL,
			value bytea,
			version INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY (type, key)
		)
	`
	_, err := tx.Exec(fmt.Sprintf(queryFormat, table))
	return err
}

func (store *sqlConfigStorage) GetConfig(networkId string, configType string, key string) (*ConfigValue, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		value, version, err := getConfig(tx, networkId, configType, key)
		if err == sql.ErrNoRows {
			return &ConfigValue{}, nil
		}
		if err != nil {
			return nil, err
		}
		return &ConfigValue{Value: value, Version: version}, nil
	}

	ret, err := sql_utils.ExecInTx(store.db, getInitFn(networkId), txFn)
	if err != nil {
		return nil, err
	}
	return ret.(*ConfigValue), err
}

func (store *sqlConfigStorage) GetConfigs(networkId string, criteria *FilterCriteria) (map[storage.TypeAndKey]*ConfigValue, error) {
	err := validateFilterCriteria(criteria)
	if err != nil {
		return nil, err
	}

	txFn := func(tx *sql.Tx) (interface{}, error) {
		tableName := GetConfigTableName(networkId)
		queryFormat := getQueryFormatStringForFilterCriteria(criteria)
		rows, err := tx.Query(fmt.Sprintf(queryFormat, tableName), getArgsFromFilterCriteria(criteria)...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		scannedRows := map[storage.TypeAndKey]*ConfigValue{}
		for rows.Next() {
			var configType, key string
			var value []byte
			var version uint64

			err = rows.Scan(&configType, &key, &value, &version)
			if err != nil {
				return nil, err
			}
			scannedRows[storage.TypeAndKey{Type: configType, Key: key}] = &ConfigValue{Value: value, Version: version}
		}
		return scannedRows, nil
	}

	ret, err := sql_utils.ExecInTx(store.db, getInitFn(networkId), txFn)
	if err != nil {
		return nil, err
	}
	return ret.(map[storage.TypeAndKey]*ConfigValue), err
}

func (store *sqlConfigStorage) ListKeysForType(networkId string, configType string) ([]string, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		tableName := GetConfigTableName(networkId)
		queryFormat := "SELECT key FROM %s WHERE type = $1"
		rows, err := tx.Query(fmt.Sprintf(queryFormat, tableName), configType)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		scannedRows := make([]string, 0)
		for rows.Next() {
			var key string
			err = rows.Scan(&key)
			if err != nil {
				return nil, err
			}
			scannedRows = append(scannedRows, key)
		}
		return scannedRows, nil
	}

	ret, err := sql_utils.ExecInTx(store.db, getInitFn(networkId), txFn)
	if err != nil {
		return nil, err
	}
	return ret.([]string), err
}

func (store *sqlConfigStorage) CreateConfig(networkId string, configType string, key string, value []byte) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		// Check for existing record
		tableName := GetConfigTableName(networkId)
		exists, err := doesConfigExist(tx, networkId, configType, key)
		if err != nil {
			return nil, err
		}
		if exists {
			err = fmt.Errorf("Creating already existing config with type %s and key %s", configType, key)
			return nil, err
		}

		queryFormat := "INSERT INTO %s (type, key, value) VALUES($1, $2, $3)"
		_, err = tx.Exec(fmt.Sprintf(queryFormat, tableName), configType, key, value)
		return nil, err
	}

	_, err := sql_utils.ExecInTx(store.db, getInitFn(networkId), txFn)
	return err
}

func (store *sqlConfigStorage) UpdateConfig(networkId string, configType string, key string, newValue []byte) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		tableName := GetConfigTableName(networkId)
		// Get current generation number and update row
		_, version, err := getConfig(tx, networkId, configType, key)
		if err != nil {
			if err == sql.ErrNoRows {
				err = fmt.Errorf("Updating nonexistent config with type %s and key %s", configType, key)
			}
			return nil, err
		}
		queryFormat := "UPDATE %s SET value = $1, version = $2 WHERE type = $3 AND key = $4"
		_, err = tx.Exec(fmt.Sprintf(queryFormat, tableName), newValue, version+1, configType, key)
		return nil, err
	}

	_, err := sql_utils.ExecInTx(store.db, getInitFn(networkId), txFn)
	return err
}

func (store *sqlConfigStorage) DeleteConfig(networkId string, configType string, key string) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		tableName := GetConfigTableName(networkId)
		exists, err := doesConfigExist(tx, networkId, configType, key)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, fmt.Errorf("Deleting nonexistent config with type %s and key %s", configType, key)
		}
		queryFormat := "DELETE FROM %s WHERE type = $1 and key = $2"
		_, err = tx.Exec(fmt.Sprintf(queryFormat, tableName), configType, key)
		return nil, err
	}

	_, err := sql_utils.ExecInTx(store.db, getInitFn(networkId), txFn)
	return err
}

func (store *sqlConfigStorage) DeleteConfigs(networkId string, criteria *FilterCriteria) error {
	err := validateFilterCriteria(criteria)
	if err != nil {
		return err
	}

	txFn := func(tx *sql.Tx) (interface{}, error) {
		action := fmt.Sprintf("DELETE FROM %s", GetConfigTableName(networkId))
		where := getWhereClauseFromFilterCriteria(criteria)
		_, err = tx.Exec(fmt.Sprintf("%s %s", action, where), getArgsFromFilterCriteria(criteria)...)
		return nil, err
	}

	_, err = sql_utils.ExecInTx(store.db, getInitFn(networkId), txFn)
	return err
}

func (store *sqlConfigStorage) DeleteConfigsForNetwork(networkId string) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		tableName := GetConfigTableName(networkId)
		queryFormat := "DROP TABLE IF EXISTS %s"
		_, err := tx.Exec(fmt.Sprintf(queryFormat, tableName))
		return nil, err
	}

	_, err := sql_utils.ExecInTx(store.db, getInitFn(networkId), txFn)
	return err
}

func getInitFn(networkID string) func(*sql.Tx) error {
	return func(tx *sql.Tx) error {
		return initTable(tx, GetConfigTableName(networkID))
	}
}

func getQueryFormatStringForFilterCriteria(criteria *FilterCriteria) string {
	selectClause := "SELECT type, key, value, version FROM %s"
	where := getWhereClauseFromFilterCriteria(criteria)
	return fmt.Sprintf("%s %s", selectClause, where)
}

// Returns "WHERE ..."
// Note no leading space
func getWhereClauseFromFilterCriteria(criteria *FilterCriteria) string {
	argIdx := 1
	var buf bytes.Buffer
	buf.WriteString("WHERE ")

	if len(criteria.Type) > 0 {
		buf.WriteString(fmt.Sprintf("type = $%d", argIdx))
		argIdx += 1
	}
	if len(criteria.Key) > 0 {
		if argIdx > 1 {
			buf.WriteString(" AND ")
		}
		buf.WriteString(fmt.Sprintf("key = $%d", argIdx))
	}
	return buf.String()
}

func getArgsFromFilterCriteria(criteria *FilterCriteria) []interface{} {
	if len(criteria.Type) > 0 && len(criteria.Key) > 0 {
		return []interface{}{criteria.Type, criteria.Key}
	} else if len(criteria.Type) > 0 && len(criteria.Key) == 0 {
		return []interface{}{criteria.Type}
	} else if len(criteria.Type) == 0 && len(criteria.Key) > 0 {
		return []interface{}{criteria.Key}
	} else {
		return []interface{}{}
	}
}

func doesConfigExist(tx *sql.Tx, networkId string, configType string, key string) (bool, error) {
	tableName := GetConfigTableName(networkId)

	var result uint64
	queryFormat := "SELECT 1 FROM %s WHERE type = $1 AND key = $2 LIMIT 1"
	err := tx.QueryRow(fmt.Sprintf(queryFormat, tableName), configType, key).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, err
	}
}

func getConfig(tx *sql.Tx, networkId string, configType string, key string) ([]byte, uint64, error) {
	tableName := GetConfigTableName(networkId)
	var value []byte
	var version uint64
	queryFormat := "SELECT value, version FROM %s WHERE type = $1 AND key = $2"
	err := tx.QueryRow(fmt.Sprintf(queryFormat, tableName), configType, key).Scan(&value, &version)
	return value, version, err
}

func validateFilterCriteria(criteria *FilterCriteria) error {
	if len(criteria.Key) == 0 && len(criteria.Type) == 0 {
		return errors.New("At least one field of filter criteria must be specified")
	}
	return nil
}
