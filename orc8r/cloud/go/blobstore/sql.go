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
func NewSQLBlobStorageFactory(tableName string, db *sql.DB) BlobStorageFactory {
	return &sqlBlobStoreFactory{tableName: tableName, db: db}
}

type sqlBlobStoreFactory struct {
	tableName string
	db        *sql.DB
}

func (fact *sqlBlobStoreFactory) StartTransaction() (TransactionalBlobStorage, error) {
	tx, err := fact.db.Begin()
	if err != nil {
		return nil, err
	}
	return &sqlBlobStorage{tableName: fact.tableName, tx: tx}, nil
}

func (fact *sqlBlobStoreFactory) InitializeFactory() error {
	tx, err := fact.db.Begin()
	if err != nil {
		return err
	}
	err = initTable(tx, fact.tableName)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

type sqlBlobStorage struct {
	tableName string
	tx        *sql.Tx
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

	queryFormat := "SELECT key FROM %s WHERE (network_id = $1 AND type = $2)"
	rows, err := store.tx.Query(
		fmt.Sprintf(queryFormat, store.tableName), networkID, typeVal,
	)
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

	query := fmt.Sprintf(
		"SELECT type, key, value, version FROM %s WHERE %s",
		store.tableName,
		getCompositeWhereInArgList(networkID, len(ids), 1),
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
	blobsToCreateAndChange := partitionBlobsToCreateAndChange(blobs, existingBlobs)

	if len(blobsToCreateAndChange.blobsToChange) > 0 {
		err := store.updateExistingBlobs(
			store.tableName,
			networkID,
			blobsToCreateAndChange.blobsToChange,
		)
		if err != nil {
			return err
		}
	}
	if len(blobsToCreateAndChange.blobsToCreate) > 0 {
		err := store.insertNewBlobs(
			store.tableName,
			networkID,
			blobsToCreateAndChange.blobsToCreate,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (store *sqlBlobStorage) Delete(networkID string, ids []storage.TypeAndKey) error {
	if err := store.validateTx(); err != nil {
		return err
	}

	query := fmt.Sprintf(
		"DELETE FROM %s WHERE %s",
		store.tableName,
		getCompositeWhereInArgList(networkID, len(ids), 1),
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

func initTable(tx *sql.Tx, tableName string) error {
	queryFormat := `
		CREATE TABLE IF NOT EXISTS %s
		(
			network_id text NOT NULL,
			type text NOT NULL,
			key text NOT NULL,
			value bytea,
			version INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY (network_id, type, key)
		)
	`
	_, err := tx.Exec(fmt.Sprintf(queryFormat, tableName))
	return err
}

func (store *sqlBlobStorage) updateExistingBlobs(
	tableName string,
	networkID string,
	blobsToChange map[storage.TypeAndKey]blobChange,
) error {
	queryFormat := "UPDATE %s SET value = $1, version = $2 WHERE network_id = $3 AND type = $4 AND key = $5"
	updateQuery := fmt.Sprintf(queryFormat, tableName)
	updateStmt, err := store.tx.Prepare(updateQuery)
	if err != nil {
		return fmt.Errorf("Error preparing update statement: %s", err)
	}
	defer updateStmt.Close()
	// Sort keys for deterministic behavior in tests
	for _, blobID := range getSortedTypeAndKeys(blobsToChange) {
		change := blobsToChange[blobID]
		_, err := updateStmt.Exec(change.new.Value, change.old.Version+1, networkID, blobID.Type, blobID.Key)
		if err != nil {
			return fmt.Errorf("Error updating blob (%s, %s, %s): %s", networkID, blobID.Type, blobID.Key, err)
		}
	}
	return nil
}

func (store *sqlBlobStorage) insertNewBlobs(tableName string, networkID string, blobs []Blob) error {
	queryFormat := "INSERT INTO %s (network_id, type, key, value) VALUES($1, $2, $3, $4)"
	insertQuery := fmt.Sprintf(queryFormat, tableName)
	insertStmt, err := store.tx.Prepare(insertQuery)
	if err != nil {
		return fmt.Errorf("Error preparing insert statement: %s", err)
	}
	defer insertStmt.Close()

	for _, blob := range blobs {
		_, err := insertStmt.Exec(networkID, blob.Type, blob.Key, blob.Value)
		if err != nil {
			return fmt.Errorf("Error creating blob (%s, %s, %s): %s", networkID, blob.Type, blob.Key, err)
		}
	}
	return nil
}

func getCompositeWhereInArgList(networkID string, numArgs int, startIdx int) string {
	retBuilder := strings.Builder{}
	retBuilder.WriteString("(")

	endIdx := startIdx + (2 * numArgs)
	for i := startIdx; i < endIdx; i += 2 {
		queryFormat := "(network_id = '%s' AND type = $%d AND key = $%d)"
		retBuilder.WriteString(fmt.Sprintf(queryFormat, networkID, i, i+1))
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
