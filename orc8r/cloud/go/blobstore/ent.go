/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package blobstore

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"magma/orc8r/cloud/go/blobstore/ent"
	"magma/orc8r/cloud/go/blobstore/ent/blob"
	"magma/orc8r/cloud/go/blobstore/ent/predicate"
	magmaerrors "magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"

	entsql "github.com/facebookincubator/ent/dialect/sql"
)

func NewEntStorage(tableName string, db *sql.DB, builder sqorc.StatementBuilder) BlobStorageFactory {
	dialect, ok := os.LookupEnv("SQL_DRIVER")
	if !ok {
		dialect = "postgres"
	}
	drv := entsql.OpenDB(dialect, db)
	client := ent.NewClient(ent.Driver(drv))
	// ent is created and initialized once per service (process).
	// therefore, it's safe to set the table used by the builders.
	blob.Table = tableName
	return &entFactory{tableName: tableName, db: db, client: client, builder: builder}
}

type entFactory struct {
	tableName string
	db        *sql.DB
	client    *ent.Client
	builder   sqorc.StatementBuilder
}

func (f *entFactory) InitializeFactory() error {
	return NewSQLBlobStorageFactory(f.tableName, f.db, f.builder).InitializeFactory()
}

func (f *entFactory) StartTransaction(opts *storage.TxOptions) (TransactionalBlobStorage, error) {
	tx, err := f.client.BeginTx(context.Background(), getSqlOpts(opts))
	if err != nil {
		return nil, err
	}
	return &entStorage{Tx: tx}, nil
}

type entStorage struct {
	*ent.Tx
}

func (e *entStorage) ListKeys(networkID string, typ string) ([]string, error) {
	ctx := context.Background()
	return e.Blob.Query().
		Where(blob.NetworkID(networkID), blob.Type(typ)).
		Select(blob.FieldKey).
		Strings(ctx)
}

func (e *entStorage) Get(networkID string, id storage.TypeAndKey) (Blob, error) {
	blobs, err := e.GetMany(networkID, []storage.TypeAndKey{id})
	if err != nil {
		return Blob{}, err
	}
	if len(blobs) == 0 {
		return Blob{}, magmaerrors.ErrNotFound
	}
	return blobs[0], nil
}

func (e *entStorage) GetMany(networkID string, ids []storage.TypeAndKey) ([]Blob, error) {
	ctx := context.Background()
	var blobs []Blob
	err := e.Blob.Query().
		Where(P(networkID, ids)).
		Select(blob.FieldKey, blob.FieldType, blob.FieldValue, blob.FieldVersion).
		Scan(ctx, &blobs)
	if err != nil {
		return nil, err
	}
	return blobs, nil
}

func (e *entStorage) IncrementVersion(networkID string, id storage.TypeAndKey) error {
	ctx := context.Background()
	switch _, err := e.Get(networkID, id); {
	case err == magmaerrors.ErrNotFound:
		_, err = e.Blob.Create().
			SetKey(id.Key).
			SetType(id.Type).
			SetNetworkID(networkID).
			SetVersion(1).
			Save(ctx)
		return err
	case err != nil: // err != not found.
		return err
	default:
		return e.Blob.Update().
			Where(blob.NetworkID(networkID), blob.Type(id.Type), blob.Key(id.Key)).
			AddVersion(1).
			Exec(ctx)
	}
}

func (e *entStorage) Delete(networkID string, ids []storage.TypeAndKey) error {
	ctx := context.Background()
	_, err := e.Blob.Delete().
		Where(P(networkID, ids)).
		Exec(ctx)
	return err
}

func (e *entStorage) CreateOrUpdate(networkID string, blobs []Blob) error {
	ctx := context.Background()
	existingBlobs, err := e.GetMany(networkID, getBlobIDs(blobs))
	if err != nil {
		return fmt.Errorf("error reading existing blobs: %s", err)
	}
	changeSet := partitionBlobsToCreateAndChange(blobs, existingBlobs)
	for _, id := range getSortedTypeAndKeys(changeSet.blobsToChange) {
		change := changeSet.blobsToChange[id]
		version := change.old.Version + 1
		if change.new.Version != 0 {
			version = change.new.Version
		}
		err := e.Blob.Update().
			SetVersion(version).
			SetValue(change.new.Value).
			Where(P(networkID, []storage.TypeAndKey{id})).
			Exec(ctx)
		if err != nil {
			return err
		}
	}
	for _, b := range changeSet.blobsToCreate {
		_, err = e.Blob.Create().
			SetKey(b.Key).
			SetType(b.Type).
			SetNetworkID(networkID).
			SetValue(b.Value).
			SetVersion(b.Version).
			Save(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *entStorage) GetExistingKeys(keys []string, filter SearchFilter) ([]string, error) {
	ctx := context.Background()
	preds := make([]predicate.Blob, 0, len(keys))
	for _, key := range keys {
		and := []predicate.Blob{
			blob.Key(key),
		}
		if nid := filter.NetworkID; nid != nil {
			and = append(and, blob.NetworkID(*nid))
		}
		preds = append(preds, blob.And(and...))
	}
	return e.Blob.Query().
		Where(blob.Or(preds...)).
		GroupBy(blob.FieldKey).
		Strings(ctx)
}

func P(networkID string, ids []storage.TypeAndKey) predicate.Blob {
	preds := make([]predicate.Blob, 0, len(ids))
	for _, id := range ids {
		preds = append(preds, blob.And(blob.NetworkID(networkID), blob.Type(id.Type), blob.Key(id.Key)))
	}
	return blob.Or(preds...)
}
