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

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

type SqlDb struct {
	db      *sql.DB
	builder sq.StatementBuilderType
}

func NewSqlDb(driver string, source string, sqlBuilder sq.StatementBuilderType) (*SqlDb, error) {
	db, err := sql_utils.Open(driver, source)
	if err != nil {
		return nil, err
	}

	return &SqlDb{
		db:      db,
		builder: sqlBuilder,
	}, nil
}

func getInitFn(table string) func(*sql.Tx) error {
	return func(tx *sql.Tx) error {
		_, err := tx.Exec(
			fmt.Sprintf(
				`CREATE TABLE IF NOT EXISTS %s (
				key text PRIMARY KEY,
				value bytea,
				generation_number INTEGER NOT NULL DEFAULT 0,
				deleted BOOLEAN NOT NULL DEFAULT FALSE
			)`, table,
			),
		)
		if err != nil {
			return errors.Wrap(err, "failed to init table")
		}
		return nil
	}
}

func (store *SqlDb) Put(table string, key string, value []byte) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		// Check if the data is already present and query for its generation number
		var generationNumber uint64
		err := store.builder.Select("generation_number").
			From(table).
			Where(sq.Eq{"key": key}).
			RunWith(tx).
			QueryRow().
			Scan(&generationNumber)
		if err != nil && err != sql.ErrNoRows {
			return nil, errors.Wrap(err, "failed to query for existing generation number")
		}

		rowExists := err == nil
		if rowExists {
			return store.builder.Update(table).
				Set("value", value).
				Set("generation_number", generationNumber+1).
				Where(sq.Eq{"key": key}).
				RunWith(tx).
				Exec()
		} else {
			return store.builder.Insert(table).
				Columns("key", "value").
				Values(key, value).
				RunWith(tx).
				Exec()
		}
	}
	_, err := sql_utils.ExecInTx(store.db, getInitFn(table), txFn)
	return err
}

func (store *SqlDb) PutMany(table string, valuesToPut map[string][]byte) (map[string]error, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		ret := map[string]error{}
		rowKeys := make([]string, len(valuesToPut))
		for k := range valuesToPut {
			rowKeys = append(rowKeys, k)
		}

		existingRows, err := store.getMany(tx, table, rowKeys)
		if err != nil {
			return ret, errors.Wrap(err, "failed to query for existing rows")
		}

		rowsToUpdate := [][3]interface{}{} // (val, gen, key)
		rowsToInsert := [][2]interface{}{} // (key, val)
		for key, newValue := range valuesToPut {
			if existingValue, keyExists := existingRows[key]; keyExists {
				rowsToUpdate = append(rowsToUpdate, [3]interface{}{newValue, existingValue.Generation + 1, key})
			} else {
				rowsToInsert = append(rowsToInsert, [2]interface{}{key, newValue})
			}
		}

		// Update existing rows
		if !funk.IsEmpty(rowsToUpdate) {
			// mock values for the update because we just want the sql string to
			// prepare with the Tx
			updateQuery, _, err := store.builder.Update(table).
				Set("value", "?").
				Set("generation_number", "?").
				Where(sq.Eq{"key": "?"}).
				ToSql()
			if err != nil {
				return ret, errors.Wrap(err, "failed to build update query")
			}

			stmts, err := sql_utils.PrepareStatements(tx, []string{updateQuery})
			if err != nil {
				return ret, errors.Wrap(err, "failed to prepare update statement")
			}
			defer sql_utils.GetCloseStatementsDeferFunc(stmts, "PutMany")()

			updateStmt := stmts[0]
			for _, row := range rowsToUpdate {
				_, err := updateStmt.Exec(row[0], row[1], row[2])
				if err != nil {
					ret[row[2].(string)] = err
				}
			}
		}

		// Insert fresh rows
		if !funk.IsEmpty(rowsToInsert) {
			insertBuilder := store.builder.Insert(table).
				Columns("key", "value")
			for _, row := range rowsToInsert {
				insertBuilder = insertBuilder.Values(row[0], row[1])
			}
			_, err := insertBuilder.RunWith(tx).Exec()
			if err != nil {
				return ret, errors.Wrap(err, "failed to create new entries")
			}
		}

		if funk.IsEmpty(ret) {
			return ret, nil
		} else {
			return ret, errors.New("failed to write entries, see return value for specific errors")
		}
	}

	ret, err := sql_utils.ExecInTx(store.db, getInitFn(table), txFn)
	return ret.(map[string]error), err
}

func (store *SqlDb) Get(table string, key string) ([]byte, uint64, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		var value []byte
		var generationNumber uint64
		err := store.builder.Select("value", "generation_number").
			From(table).
			Where(sq.Eq{"key": key}).
			RunWith(tx).
			QueryRow().Scan(&value, &generationNumber)
		if err == sql.ErrNoRows {
			return ValueWrapper{}, ErrNotFound
		}
		return ValueWrapper{Value: value, Generation: generationNumber}, err
	}

	ret, err := sql_utils.ExecInTx(store.db, getInitFn(table), txFn)
	if err != nil {
		return nil, 0, err
	}
	vw := ret.(ValueWrapper)
	return vw.Value, vw.Generation, nil
}

func (store *SqlDb) GetMany(table string, keys []string) (map[string]ValueWrapper, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		return store.getMany(tx, table, keys)
	}
	ret, err := sql_utils.ExecInTx(store.db, getInitFn(table), txFn)
	return ret.(map[string]ValueWrapper), err
}

func (store *SqlDb) Delete(table string, key string) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		return store.builder.Delete(table).Where(sq.Eq{"key": key}).RunWith(tx).Exec()
	}
	_, err := sql_utils.ExecInTx(store.db, getInitFn(table), txFn)
	return err
}

func (store *SqlDb) DeleteMany(table string, keys []string) (map[string]error, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		return store.builder.Delete(table).Where(sq.Eq{"key": keys}).RunWith(tx).Exec()
	}
	_, err := sql_utils.ExecInTx(store.db, getInitFn(table), txFn)
	return map[string]error{}, err
}

func (store *SqlDb) ListKeys(table string) ([]string, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		rows, err := store.builder.Select("key").From(table).RunWith(tx).Query()
		if err != nil {
			return []string{}, errors.Wrap(err, "failed to query for keys")
		}
		defer sql_utils.CloseRowsLogOnError(rows, "ListKeys")

		keys := []string{}
		for rows.Next() {
			var key string
			if err = rows.Scan(&key); err != nil {
				return []string{}, errors.Wrap(err, "failed to read key")
			}
			keys = append(keys, key)
		}
		return keys, nil
	}

	ret, err := sql_utils.ExecInTx(store.db, getInitFn(table), txFn)
	return ret.([]string), err
}

func (store *SqlDb) DeleteTable(table string) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		return tx.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
	}
	// No initFn param because why would we create a table that we're dropping
	_, err := sql_utils.ExecInTx(store.db, func(*sql.Tx) error { return nil }, txFn)
	return err
}

func (store *SqlDb) DoesKeyExist(table string, key string) (bool, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		var placeHolder uint64
		err := store.builder.Select("1").From(table).
			Where(sq.Eq{"key": key}).
			Limit(1).
			RunWith(tx).
			QueryRow().Scan(&placeHolder)
		if err != nil {
			if err == sql.ErrNoRows {
				return false, nil
			}
			return false, err
		}
		return true, nil
	}
	ret, err := sql_utils.ExecInTx(store.db, getInitFn(table), txFn)
	return ret.(bool), err
}

func (store *SqlDb) getMany(tx *sql.Tx, table string, keys []string) (map[string]ValueWrapper, error) {
	valuesByKey := make(map[string]ValueWrapper)
	if len(keys) == 0 {
		return valuesByKey, nil
	}

	rows, err := store.builder.Select("key", "value", "generation_number").
		From(table).
		Where(sq.Eq{"key": keys}).
		RunWith(tx).
		Query()
	if err != nil {
		return valuesByKey, err
	}
	defer sql_utils.CloseRowsLogOnError(rows, "getMany")
	return getSqlRowsAsMap(rows)
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
