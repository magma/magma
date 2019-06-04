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
	"fmt"
	"sort"

	magmaerrors "magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"

	sq "github.com/Masterminds/squirrel"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

// NewSQLBlobStorageFactory returns a BlobStorageFactory implementation which
// will return storage APIs backed by SQL.
func NewSQLBlobStorageFactory(tableName string, db *sql.DB, sqlBuilder sqorc.StatementBuilder) BlobStorageFactory {
	return &sqlBlobStoreFactory{tableName: tableName, db: db, builder: sqlBuilder}
}

type sqlBlobStoreFactory struct {
	tableName string
	db        *sql.DB
	builder   sqorc.StatementBuilder
}

func (fact *sqlBlobStoreFactory) StartTransaction() (TransactionalBlobStorage, error) {
	tx, err := fact.db.Begin()
	if err != nil {
		return nil, err
	}
	return &sqlBlobStorage{tableName: fact.tableName, tx: tx, builder: fact.builder}, nil
}

func (fact *sqlBlobStoreFactory) InitializeFactory() error {
	tx, err := fact.db.Begin()
	if err != nil {
		return err
	}
	err = fact.initTable(tx, fact.tableName)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			glog.Errorf("error rolling back transaction initializing blobstore factory: %s", err)
		}

		return err
	}
	return tx.Commit()
}

func (fact *sqlBlobStoreFactory) initTable(tx *sql.Tx, tableName string) error {
	_, err := fact.builder.CreateTable(tableName).
		IfNotExists().
		Column("network_id").Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column("type").Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column("key").Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column("value").Type(sqorc.ColumnTypeBytes).EndColumn().
		Column("version").Type(sqorc.ColumnTypeInt).NotNull().Default(0).EndColumn().
		PrimaryKey("network_id", "type", "key").
		RunWith(tx).
		Exec()
	return err
}

type sqlBlobStorage struct {
	tableName string
	tx        *sql.Tx
	builder   sqorc.StatementBuilder
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

	rows, err := store.builder.Select("key").From(store.tableName).
		Where(sq.Eq{"network_id": networkID, "type": typeVal}).
		RunWith(store.tx).
		Query()
	if err != nil {
		return ret, err
	}
	defer sqorc.CloseRowsLogOnError(rows, "ListKeys")

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

	whereCondition := getWhereCondition(networkID, ids)
	rows, err := sq.Select("type", "key", "value", "version").From(store.tableName).
		Where(whereCondition).
		RunWith(store.tx).
		Query()
	if err != nil {
		return emptyRet, err
	}
	defer sqorc.CloseRowsLogOnError(rows, "GetMany")

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
	// defer tx validation to GetMany
	existingBlobs, err := store.GetMany(networkID, getBlobIDs(blobs))
	if err != nil {
		return fmt.Errorf("Error reading existing blobs: %s", err)
	}
	blobsToCreateAndChange := partitionBlobsToCreateAndChange(blobs, existingBlobs)

	if len(blobsToCreateAndChange.blobsToChange) > 0 {
		err := store.updateExistingBlobs(networkID, blobsToCreateAndChange.blobsToChange)
		if err != nil {
			return err
		}
	}
	if len(blobsToCreateAndChange.blobsToCreate) > 0 {
		err := store.insertNewBlobs(networkID, blobsToCreateAndChange.blobsToCreate)
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

	whereCondition := getWhereCondition(networkID, ids)
	_, err := store.builder.Delete(store.tableName).
		Where(whereCondition).
		RunWith(store.tx).
		Exec()
	return err
}

func (store *sqlBlobStorage) validateTx() error {
	if store.tx == nil {
		return errors.New("No transaction is available")
	}
	return nil
}

func (store *sqlBlobStorage) updateExistingBlobs(networkID string, blobsToChange map[storage.TypeAndKey]blobChange) error {
	// Let squirrel cache prepared statements for us (there should only be 1)
	sc := sq.NewStmtCache(store.tx)
	defer sqorc.ClearStatementCacheLogOnError(sc, "updateExistingBlobs")

	// Sort keys for deterministic behavior in tests
	for _, blobID := range getSortedTypeAndKeys(blobsToChange) {
		change := blobsToChange[blobID]
		_, err := store.builder.Update(store.tableName).
			Set("value", change.new.Value).
			Set("version", change.old.Version+1).
			Where(
				// Use explicit sq.And to preserve ordering of WHERE clause items
				sq.And{
					sq.Eq{"network_id": networkID},
					sq.Eq{"type": blobID.Type},
					sq.Eq{"key": blobID.Key},
				},
			).
			RunWith(sc).
			Exec()
		if err != nil {
			return fmt.Errorf("Error updating blob (%s, %s, %s): %s", networkID, blobID.Type, blobID.Key, err)
		}
	}
	return nil
}

func (store *sqlBlobStorage) insertNewBlobs(networkID string, blobs []Blob) error {
	insertBuilder := store.builder.Insert(store.tableName).
		Columns("network_id", "type", "key", "value")
	for _, blob := range blobs {
		insertBuilder = insertBuilder.Values(networkID, blob.Type, blob.Key, blob.Value)
	}
	_, err := insertBuilder.RunWith(store.tx).Exec()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error creating blobs"))
	}
	return nil
}

func getWhereCondition(networkID string, ids []storage.TypeAndKey) sq.Or {
	whereConditions := make(sq.Or, 0, len(ids))
	for _, id := range ids {
		// Use explicit sq.And to preserve ordering of clauses for testing
		whereConditions = append(whereConditions, sq.And{
			sq.Eq{"network_id": networkID},
			sq.Eq{"type": id.Type},
			sq.Eq{"key": id.Key},
		})
	}
	return whereConditions
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
