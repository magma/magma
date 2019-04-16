/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package blobstore

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"

	magmaerrors "magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/storage"
)

// NewSQLBlobStorageFactory returns a BlobStorageFactory implementation which
// will return storage APIs backed by SQL.
func NewSQLBlobStorageFactory(baseTableName string, db *sql.DB) BlobStorageFactory {
	return &sqlBlobStoreFactory{baseTableName: baseTableName, db: db}
}

type sqlBlobStoreFactory struct {
	baseTableName string
	db            *sql.DB
}

func (fact *sqlBlobStoreFactory) StartTransaction() (TransactionalBlobStorage, error) {
	tx, err := fact.db.Begin()
	if err != nil {
		return nil, err
	}
	return &sqlBlobStorage{baseTableName: fact.baseTableName, tx: tx}, nil
}

type sqlBlobStorage struct {
	baseTableName string
	tx            *sql.Tx
}

func (store *sqlBlobStorage) Commit() error {
	if store.tx == nil {
		return errors.New("There is no current transaction to commit")
	}

	err := store.tx.Commit()
	store.tx = nil
	return err
}

func (store *sqlBlobStorage) Rollback() error {
	if store.tx == nil {
		return errors.New("There is no current transaction to rollback")
	}

	err := store.tx.Rollback()
	store.tx = nil
	return err
}

func (store *sqlBlobStorage) ListKeys(networkID string, typeVal string) ([]string, error) {
	ret := []string{}
	if err := store.validateTx(); err != nil {
		return ret, err
	}

	tableName := GetTableName(networkID, store.baseTableName)
	if err := store.initTable(tableName); err != nil {
		return ret, err
	}

	queryFormat := "SELECT key FROM %s WHERE type = $1"
	rows, err := store.tx.Query(fmt.Sprintf(queryFormat, tableName), typeVal)
	if err != nil {
		return ret, err
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		err = rows.Scan(&key)
		if err != nil {
			return []string{}, err
		}
		ret = append(ret, key)
	}
	return ret, nil
}

func (store *sqlBlobStorage) Get(networkID string, id storage.TypeAndKey) (Blob, error) {
	// defer table init, tx validation to GetMany
	multiRet, err := store.GetMany(networkID, []storage.TypeAndKey{id})
	if err != nil {
		return Blob{}, err
	}
	if len(multiRet) == 0 {
		return Blob{}, magmaerrors.ErrNotFound
	}
	return multiRet[0], nil
}

func (store *sqlBlobStorage) GetMany(networkID string, ids []storage.TypeAndKey) ([]Blob, error) {
	emptyRet := []Blob{}
	if err := store.validateTx(); err != nil {
		return emptyRet, err
	}

	tableName := GetTableName(networkID, store.baseTableName)
	if err := store.initTable(tableName); err != nil {
		return emptyRet, err
	}

	query := fmt.Sprintf(
		"SELECT type, key, value, version FROM %s WHERE %s",
		tableName,
		getCompositeWhereInArgList(1, len(ids)),
	)
	rows, err := store.tx.Query(query, typeAndKeysToArgs(ids)...)
	if err != nil {
		return emptyRet, err
	}
	defer rows.Close()

	scannedRows := []Blob{}
	for rows.Next() {
		var t, k string
		var val []byte
		var version uint64

		err = rows.Scan(&t, &k, &val, &version)
		if err != nil {
			return emptyRet, err
		}
		scannedRows = append(scannedRows, Blob{Type: t, Key: k, Value: val, Version: version})
	}
	return scannedRows, nil
}

func (store *sqlBlobStorage) CreateOrUpdate(networkID string, blobs []Blob) error {
	// defer table init, tx validation to GetMany
	existingBlobs, err := store.GetMany(networkID, getBlobIDs(blobs))
	if err != nil {
		return fmt.Errorf("Error reading existing blobs: %s", err)
	}
	tableName := GetTableName(networkID, store.baseTableName)
	blobsToCreateAndChange := partitionBlobsToCreateAndChange(blobs, existingBlobs)

	if len(blobsToCreateAndChange.blobsToChange) > 0 {
		if err := store.updateExistingBlobs(tableName, blobsToCreateAndChange.blobsToChange); err != nil {
			return err
		}
	}
	if len(blobsToCreateAndChange.blobsToCreate) > 0 {
		if err := store.insertNewBlobs(tableName, blobsToCreateAndChange.blobsToCreate); err != nil {
			return err
		}
	}

	return nil
}

func (store *sqlBlobStorage) Delete(networkID string, ids []storage.TypeAndKey) error {
	if err := store.validateTx(); err != nil {
		return err
	}

	tableName := GetTableName(networkID, store.baseTableName)
	if err := store.initTable(tableName); err != nil {
		return err
	}

	query := fmt.Sprintf(
		"DELETE FROM %s WHERE %s",
		tableName,
		getCompositeWhereInArgList(1, len(ids)),
	)
	_, err := store.tx.Exec(query, typeAndKeysToArgs(ids)...)
	return err
}

func (store *sqlBlobStorage) validateTx() error {
	if store.tx == nil {
		return errors.New("No transaction is available")
	}
	return nil
}

func (store *sqlBlobStorage) initTable(fullTableName string) error {
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
	_, err := store.tx.Exec(fmt.Sprintf(queryFormat, fullTableName))
	return err
}

func (store *sqlBlobStorage) updateExistingBlobs(tableName string, blobsToChange map[storage.TypeAndKey]blobChange) error {
	updateQuery := fmt.Sprintf("UPDATE %s SET value = $1, version = $2 WHERE type = $3 AND key = $4", tableName)
	updateStmt, err := store.tx.Prepare(updateQuery)
	if err != nil {
		return fmt.Errorf("Error preparing update statement: %s", err)
	}
	defer updateStmt.Close()

	// Sort keys for deterministic behavior in tests
	for _, blobID := range getSortedTypeAndKeys(blobsToChange) {
		change := blobsToChange[blobID]
		_, err := updateStmt.Exec(change.new.Value, change.old.Version+1, blobID.Type, blobID.Key)
		if err != nil {
			return fmt.Errorf("Error updating blob (%s, %s): %s", blobID.Type, blobID.Key, err)
		}
	}
	return nil
}

func (store *sqlBlobStorage) insertNewBlobs(tableName string, blobs []Blob) error {
	insertQuery := fmt.Sprintf("INSERT INTO %s (type, key, value) VALUES($1, $2, $3)", tableName)
	insertStmt, err := store.tx.Prepare(insertQuery)
	if err != nil {
		return fmt.Errorf("Error preparing insert statement: %s", err)
	}
	defer insertStmt.Close()

	for _, blob := range blobs {
		_, err := insertStmt.Exec(blob.Type, blob.Key, blob.Value)
		if err != nil {
			return fmt.Errorf("Error creating blob (%s, %s): %s", blob.Type, blob.Key, err)
		}
	}
	return nil
}

func getCompositeWhereInArgList(startIdx int, numArgs int) string {
	retBuilder := strings.Builder{}
	retBuilder.WriteString("(")

	endIdx := startIdx + (2 * numArgs)
	for i := startIdx; i < endIdx; i += 2 {
		retBuilder.WriteString(fmt.Sprintf("(type = $%d AND key = $%d)", i, i+1))
		if i < endIdx-2 {
			retBuilder.WriteString(" OR ")
		}
	}

	retBuilder.WriteString(")")
	return retBuilder.String()
}

func typeAndKeysToArgs(ids []storage.TypeAndKey) []interface{} {
	ret := make([]interface{}, 0, len(ids)*2)
	for _, tk := range ids {
		ret = append(ret, tk.Type)
		ret = append(ret, tk.Key)
	}
	return ret
}

func getBlobIDs(blobs []Blob) []storage.TypeAndKey {
	ret := make([]storage.TypeAndKey, 0, len(blobs))
	for _, blob := range blobs {
		ret = append(ret, storage.TypeAndKey{Type: blob.Type, Key: blob.Key})
	}
	return ret
}

type blobChange struct {
	old Blob
	new Blob
}

type blobsToCreateAndChange struct {
	blobsToCreate []Blob
	blobsToChange map[storage.TypeAndKey]blobChange
}

func partitionBlobsToCreateAndChange(blobsToUpdate []Blob, existingBlobs []Blob) blobsToCreateAndChange {
	ret := blobsToCreateAndChange{
		blobsToCreate: []Blob{},
		blobsToChange: map[storage.TypeAndKey]blobChange{},
	}
	existingBlobsByID := GetBlobsByTypeAndKey(existingBlobs)

	for _, blob := range blobsToUpdate {
		blobID := storage.TypeAndKey{Type: blob.Type, Key: blob.Key}
		oldBlob, exists := existingBlobsByID[blobID]
		if exists {
			ret.blobsToChange[blobID] = blobChange{old: oldBlob, new: blob}
		} else {
			ret.blobsToCreate = append(ret.blobsToCreate, blob)
		}
	}
	return ret
}

func getSortedTypeAndKeys(blobsToChange map[storage.TypeAndKey]blobChange) []storage.TypeAndKey {
	ret := make([]storage.TypeAndKey, 0, len(blobsToChange))
	for k := range blobsToChange {
		ret = append(ret, k)
	}
	sort.Slice(ret, func(i, j int) bool { return ret[i].String() < ret[j].String() })
	return ret
}
