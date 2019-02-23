/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package datastore

import (
	"database/sql"
	"fmt"

	"magma/orc8r/cloud/go/sql_utils"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type SqlDb struct {
	db *sql.DB
}

type SqlQueryable interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func NewSqlDb(driver string, source string) (*SqlDb, error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		return nil, err
	}

	store := new(SqlDb)
	store.db = db
	return store, nil
}

func initTable(queryable SqlQueryable, table string) error {
	_, err := queryable.Exec(fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s (key text PRIMARY KEY, value bytea,
		generation_number INTEGER NOT NULL DEFAULT 0,
		deleted BOOLEAN NOT NULL DEFAULT FALSE)`, table))
	return err
}

func (store *SqlDb) Put(table string, key string, value []byte) error {
	// Create a transaction for the lookup and insert/update operation.
	// This also guarantees the atomicity of generation number increment.
	tx, err := store.db.Begin()
	if err != nil {
		return err
	}

	err = initTable(tx, table)
	if err != nil {
		return err
	}

	// Check if the data is already present and query for its generation number
	var generationNumber uint64
	err = tx.QueryRow(fmt.Sprintf(
		"SELECT generation_number FROM %s WHERE key = $1",
		table), key).Scan(&generationNumber)

	if err != nil {
		// Insert the new data
		_, err = tx.Exec(fmt.Sprintf(
			"INSERT INTO %s (key, value) VALUES($1, $2)", table), key, value)
	} else {
		// Update existing data and increment generation number
		_, err = tx.Exec(fmt.Sprintf(
			"UPDATE %s SET value = $1, generation_number = $2 WHERE key = $3",
			table), value, generationNumber+1, key)
	}

	if err != nil {
		// Error occured with the operation. Rollback the transaction.
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (store *SqlDb) PutMany(table string, valuesToPut map[string][]byte) (map[string]error, error) {
	rowKeys := make([]string, len(valuesToPut))
	for k := range valuesToPut {
		rowKeys = append(rowKeys, k)
	}

	err := initTable(store.db, table)
	if err != nil {
		return map[string]error{}, err
	}

	existingRows, err := store.GetMany(table, rowKeys)
	if err != nil {
		return map[string]error{}, err
	}

	updateStmt, err := store.db.Prepare(fmt.Sprintf(
		"UPDATE %s SET value = $1, generation_number = $2 WHERE key = $3", table))
	if err != nil {
		return map[string]error{}, err
	}
	defer updateStmt.Close()

	insertStmt, err := store.db.Prepare(fmt.Sprintf(
		"INSERT INTO %s (key, value) VALUES($1, $2)", table))
	if err != nil {
		return map[string]error{}, err
	}
	defer insertStmt.Close()

	errorMap := make(map[string]error)
	for key, newValue := range valuesToPut {
		if existingValue, keyExists := existingRows[key]; keyExists {
			newGeneration := existingValue.Generation + 1
			_, err = updateStmt.Exec(newValue, newGeneration, key)
		} else {
			_, err = insertStmt.Exec(key, newValue)
		}

		if err != nil {
			errorMap[key] = err
		}
	}

	return errorMap, nil
}

func (store *SqlDb) Get(table string, key string) ([]byte, uint64, error) {
	var value []byte
	var generationNumber uint64
	if err := initTable(store.db, table); err != nil {
		return value, generationNumber, err
	}
	err := store.db.QueryRow(fmt.Sprintf(
		"SELECT value, generation_number FROM %s WHERE key = $1",
		table), key).Scan(&value, &generationNumber)
	if err == sql.ErrNoRows {
		return value, generationNumber, ErrNotFound
	}
	return value, generationNumber, err
}

func (store *SqlDb) GetMany(table string, keys []string) (map[string]ValueWrapper, error) {
	valuesByKey := make(map[string]ValueWrapper)
	if err := initTable(store.db, table); err != nil {
		return valuesByKey, err
	}
	if len(keys) == 0 {
		return valuesByKey, nil
	}

	queryString, queryArgs := getSelectInQueryAndArgs(table, keys)
	rows, err := store.db.Query(queryString, queryArgs...)
	if err != nil {
		return valuesByKey, err
	}
	defer rows.Close()
	return getSqlRowsAsMap(rows)
}

func getSelectInQueryAndArgs(table string, keys []string) (string, []interface{}) {
	inList := sql_utils.GetPlaceholderArgList(1, len(keys))
	queryString := fmt.Sprintf(
		"SELECT key, value, generation_number FROM %s WHERE key IN %s",
		table, inList)
	queryArgs := make([]interface{}, len(keys))
	for i := range keys {
		queryArgs[i] = keys[i]
	}
	return queryString, queryArgs
}

func getSqlRowsAsMap(rows *sql.Rows) (map[string]ValueWrapper, error) {
	var valuesByKey = make(map[string]ValueWrapper)

	for rows.Next() {
		var key string
		var value []byte
		var generationNumber uint64

		err := rows.Scan(&key, &value, &generationNumber)
		if err != nil {
			return map[string]ValueWrapper{}, err
		}

		valuesByKey[key] = ValueWrapper{
			Value:      value,
			Generation: generationNumber,
		}
	}

	return valuesByKey, nil
}

func (store *SqlDb) Delete(table string, key string) error {
	_, err := store.db.Exec(fmt.Sprintf(
		"DELETE FROM %s WHERE key = $1", table), key)
	return err
}

func (store *SqlDb) DeleteMany(table string, keys []string) (map[string]error, error) {
	err := initTable(store.db, table)
	if err != nil {
		return map[string]error{}, err
	}

	stmt, err := store.db.Prepare(
		fmt.Sprintf("DELETE FROM %s WHERE key = $1", table))
	if err != nil {
		return map[string]error{}, err
	}

	errMap := make(map[string]error)
	for _, k := range keys {
		_, err = stmt.Exec(k)
		if err != nil {
			errMap[k] = err
		}
	}
	return errMap, nil
}

func (store *SqlDb) ListKeys(table string) ([]string, error) {
	if err := initTable(store.db, table); err != nil {
		return nil, err
	}

	rows, err := store.db.Query(fmt.Sprintf("SELECT key FROM %s", table))
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0)
	for rows.Next() {
		var key string
		if err = rows.Scan(&key); err != nil {
			rows.Close()
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, rows.Err()
}

func (store *SqlDb) DeleteTable(table string) error {
	_, err := store.db.Exec("DROP TABLE IF EXISTS " + table)
	return err
}

func (store *SqlDb) DoesKeyExist(table string, key string) (bool, error) {
	var placeHolder uint64
	if err := initTable(store.db, table); err != nil {
		return false, err
	}
	err := store.db.QueryRow(
		fmt.Sprintf("SELECT 1 FROM %s WHERE key = $1 LIMIT 1", table),
		key,
	).Scan(&placeHolder)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
