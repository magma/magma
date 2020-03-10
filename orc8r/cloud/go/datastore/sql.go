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

	"magma/orc8r/cloud/go/sqorc"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

const (
	// escaped for mysql compat
	keyCol     = "\"key\""
	valueCol   = "value"
	genCol     = "generation_number"
	deletedCol = "deleted"
)

type SqlDb struct {
	db      *sql.DB
	builder sqorc.StatementBuilder
}

func NewSqlDb(driver string, source string, sqlBuilder sqorc.StatementBuilder) (*SqlDb, error) {
	db, err := sqorc.Open(driver, source)
	if err != nil {
		return nil, err
	}

	return &SqlDb{
		db:      db,
		builder: sqlBuilder,
	}, nil
}

func (store *SqlDb) getInitFn(table string) func(*sql.Tx) error {
	return func(tx *sql.Tx) error {
		_, err := store.builder.CreateTable(table).
			IfNotExists().
			// table builder escapes all columns by default
			Column(keyCol).Type(sqorc.ColumnTypeText).PrimaryKey().EndColumn().
			Column(valueCol).Type(sqorc.ColumnTypeBytes).EndColumn().
			Column(genCol).Type(sqorc.ColumnTypeInt).NotNull().Default(0).EndColumn().
			Column(deletedCol).Type(sqorc.ColumnTypeBool).NotNull().Default("FALSE").EndColumn().
			RunWith(tx).
			Exec()
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
		err := store.builder.Select(genCol).
			From(table).
			Where(sq.Eq{keyCol: key}).
			RunWith(tx).
			QueryRow().
			Scan(&generationNumber)
		if err != nil && err != sql.ErrNoRows {
			return nil, errors.Wrap(err, "failed to query for existing generation number")
		}

		rowExists := err == nil
		if rowExists {
			return store.builder.Update(table).
				Set(valueCol, value).
				Set(genCol, generationNumber+1).
				Where(sq.Eq{keyCol: key}).
				RunWith(tx).
				Exec()
		} else {
			return store.builder.Insert(table).
				Columns(keyCol, valueCol).
				Values(key, value).
				RunWith(tx).
				Exec()
		}
	}
	_, err := sqorc.ExecInTx(store.db, store.getInitFn(table), txFn)
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

		// Let squirrel cache prepared statements for us on update
		sc := sq.NewStmtCache(tx)
		defer sqorc.ClearStatementCacheLogOnError(sc, "PutMany")

		// Update existing rows
		for _, row := range rowsToUpdate {
			_, err := store.builder.Update(table).
				Set(valueCol, row[0]).
				Set(genCol, row[1]).
				Where(sq.Eq{keyCol: row[2]}).
				RunWith(sc).
				Exec()
			if err != nil {
				ret[row[2].(string)] = err
			}
		}

		// Insert fresh rows
		if !funk.IsEmpty(rowsToInsert) {
			insertBuilder := store.builder.Insert(table).
				Columns(keyCol, valueCol)
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

	ret, err := sqorc.ExecInTx(store.db, store.getInitFn(table), txFn)
	return ret.(map[string]error), err
}

func (store *SqlDb) Get(table string, key string) ([]byte, uint64, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		var value []byte
		var generationNumber uint64
		err := store.builder.Select(valueCol, genCol).
			From(table).
			Where(sq.Eq{keyCol: key}).
			RunWith(tx).
			QueryRow().Scan(&value, &generationNumber)
		if err == sql.ErrNoRows {
			return ValueWrapper{}, ErrNotFound
		}
		return ValueWrapper{Value: value, Generation: generationNumber}, err
	}

	ret, err := sqorc.ExecInTx(store.db, store.getInitFn(table), txFn)
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
	ret, err := sqorc.ExecInTx(store.db, store.getInitFn(table), txFn)
	return ret.(map[string]ValueWrapper), err
}

func (store *SqlDb) Delete(table string, key string) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		return store.builder.Delete(table).Where(sq.Eq{keyCol: key}).RunWith(tx).Exec()
	}
	_, err := sqorc.ExecInTx(store.db, store.getInitFn(table), txFn)
	return err
}

func (store *SqlDb) DeleteMany(table string, keys []string) (map[string]error, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		return store.builder.Delete(table).Where(sq.Eq{keyCol: keys}).RunWith(tx).Exec()
	}
	_, err := sqorc.ExecInTx(store.db, store.getInitFn(table), txFn)
	return map[string]error{}, err
}

func (store *SqlDb) ListKeys(table string) ([]string, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		rows, err := store.builder.Select(keyCol).From(table).RunWith(tx).Query()
		if err != nil {
			return []string{}, errors.Wrap(err, "failed to query for keys")
		}
		defer sqorc.CloseRowsLogOnError(rows, "ListKeys")

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

	ret, err := sqorc.ExecInTx(store.db, store.getInitFn(table), txFn)
	return ret.([]string), err
}

func (store *SqlDb) DeleteTable(table string) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		return tx.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
	}
	// No initFn param because why would we create a table that we're dropping
	_, err := sqorc.ExecInTx(store.db, func(*sql.Tx) error { return nil }, txFn)
	return err
}

func (store *SqlDb) DoesKeyExist(table string, key string) (bool, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		var placeHolder uint64
		err := store.builder.Select("1").From(table).
			Where(sq.Eq{keyCol: key}).
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
	ret, err := sqorc.ExecInTx(store.db, store.getInitFn(table), txFn)
	return ret.(bool), err
}

func (store *SqlDb) getMany(tx *sql.Tx, table string, keys []string) (map[string]ValueWrapper, error) {
	valuesByKey := map[string]ValueWrapper{}
	if len(keys) == 0 {
		return valuesByKey, nil
	}

	rows, err := store.builder.Select(keyCol, valueCol, genCol).
		From(table).
		Where(sq.Eq{keyCol: keys}).
		RunWith(tx).
		Query()
	if err != nil {
		return valuesByKey, err
	}
	defer sqorc.CloseRowsLogOnError(rows, "getMany")
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
