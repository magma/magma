/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage

import (
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"

	"github.com/pkg/errors"
)

const (
	// Blobstore needs a network ID, but directoryd is network-agnostic so we
	// will use a placeholder value
	placeholderNetworkId = "placeholder_network"
)

// NewDirectorydBlobstoreStorage returns a directoryd storage implementation
// backed by the provided blobstore factory.
func NewDirectorydBlobstoreStorage(factory blobstore.BlobStorageFactory) DirectorydPersistenceService {
	return &directorydBlobstoreStorage{factory: factory}
}

type directorydBlobstoreStorage struct {
	factory blobstore.BlobStorageFactory
}

func (d *directorydBlobstoreStorage) GetRecord(tableId protos.TableID, recordId string) (*protos.LocationRecord, error) {
	typeVal, ok := protos.TableID_name[int32(tableId)]
	if !ok {
		return nil, errors.Errorf("unrecognized table ID %v", tableId)
	}

	store, err := d.factory.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, errors.Wrap(err, "failed to start transaction")
	}
	defer store.Rollback()

	recordBlob, err := store.Get(placeholderNetworkId, storage.TypeAndKey{Type: typeVal, Key: recordId})
	if err == merrors.ErrNotFound {
		return nil, err
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to load location record")
	}

	ret := &protos.LocationRecord{}
	if err := protos.Unmarshal(recordBlob.Value, ret); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal location record")
	}
	return ret, store.Commit()
}

func (d *directorydBlobstoreStorage) UpdateOrCreateRecord(tableId protos.TableID, recordId string, record *protos.LocationRecord) error {
	typeVal, ok := protos.TableID_name[int32(tableId)]
	if !ok {
		return errors.Errorf("unrecognized table ID %v", tableId)
	}

	store, err := d.factory.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	defer store.Rollback()

	value, err := protos.MarshalIntern(record)
	if err != nil {
		return errors.Wrap(err, "failed to marshal location record")
	}

	blob := blobstore.Blob{Type: typeVal, Key: recordId, Value: value}
	err = store.CreateOrUpdate(placeholderNetworkId, []blobstore.Blob{blob})
	if err != nil {
		return errors.Wrap(err, "failed to create or update location record")
	}
	return store.Commit()
}

func (d *directorydBlobstoreStorage) DeleteRecord(tableId protos.TableID, recordId string) error {
	typeVal, ok := protos.TableID_name[int32(tableId)]
	if !ok {
		return errors.Errorf("unrecognized table ID %v", tableId)
	}

	store, err := d.factory.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	defer store.Rollback()

	typeAndKey := storage.TypeAndKey{Type: typeVal, Key: recordId}
	err = store.Delete(placeholderNetworkId, []storage.TypeAndKey{typeAndKey})
	if err != nil {
		return errors.Wrap(err, "failed to delete location record")
	}
	return store.Commit()
}
