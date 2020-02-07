/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package blobstore

import (
	"magma/orc8r/cloud/go/storage"
)

// Blob encapsulates a blob for storage
type Blob struct {
	Type    string
	Key     string
	Value   []byte
	Version uint64
}

type SearchFilter struct {
	NetworkID *string
}

// BlobStorageFactory is an API to create a storage API bound to a transaction.
type BlobStorageFactory interface {
	InitializeFactory() error
	// StartTransaction opens a transaction for all following blob storage
	// operations, and returns a TransactionalBlobStorage instance tied to the
	// opened transaction.
	StartTransaction(opts *storage.TxOptions) (TransactionalBlobStorage, error)
}

// TransactionalBlobStorage is the client API for blob storage operations
// within the context of a transaction.
type TransactionalBlobStorage interface {

	// Commit commits the existing transaction. If an error is returned from
	// the backing storage while committing, the transaction will be rolled
	// back.
	Commit() error

	// Rollback rolls back the existing transaction. If the targeted
	// transaction has already been committed, Rollback has no effect and
	// returns an error.
	Rollback() error

	// ListKeys returns all the blob keys stored for the network and type.
	ListKeys(networkID string, typeVal string) ([]string, error)

	// Get loads a specific blob from storage.
	// If there is no blob matching the given ID, ErrNotFound from
	// magma/orc8r/cloud/go/errors will be returned.
	Get(networkID string, id storage.TypeAndKey) (Blob, error)

	// Get loads and returns a collection of blobs matching the specified IDs.
	// If there is no blob corresponding to a TypeAndKey, the returned list
	// will not have a corresponding Blob.
	GetMany(networkID string, ids []storage.TypeAndKey) ([]Blob, error)

	// CreateOrUpdate writes blobs to the storage. Blobs are either updated
	// in-place or created. The Version field of Blobs passed here will be used
	// if it is not set to 0. Otherwise version incrementation will be handled
	// internally inside the storage implementation.
	CreateOrUpdate(networkID string, blobs []Blob) error

	// GetExistingKeys takes in a list of keys and returns a list of keys
	// that exist from the input. The filter specifies whether to look at the
	// entire storage or just in a network.
	GetExistingKeys(keys []string, filter SearchFilter) ([]string, error)

	// Delete deletes specified blobs from storage.
	Delete(networkID string, ids []storage.TypeAndKey) error

	// IncrementVersion is an atomic upsert (INSERT DO ON CONFLICT) that
	// increments the version column or inserts 1 if it does not exist.
	IncrementVersion(networkID string, id storage.TypeAndKey) error
}

// GetTKsFromKeys returns the passed keys mapped as TypeAndKey, with the passed
// type applied to each.
func GetTKsFromKeys(typ string, keys []string) []storage.TypeAndKey {
	tks := make([]storage.TypeAndKey, 0, len(keys))
	for _, k := range keys {
		tks = append(tks, storage.TypeAndKey{Type: typ, Key: k})
	}
	return tks
}

// GetBlobsByTypeAndKey returns a computed view of a list of blobs as a map of
// blobs keyed by blob TypeAndKey.
func GetBlobsByTypeAndKey(blobs []Blob) map[storage.TypeAndKey]Blob {
	ret := make(map[storage.TypeAndKey]Blob, len(blobs))
	for _, blob := range blobs {
		ret[storage.TypeAndKey{Type: blob.Type, Key: blob.Key}] = blob
	}
	return ret
}
