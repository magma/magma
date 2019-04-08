/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package blobstore

import (
	"fmt"

	"magma/orc8r/cloud/go/blobstore/protos"
)

// TypeAndKey is an identifier for a blob
type TypeAndKey struct {
	Type, Key string
}

// FromProto fills in this TypeAndKey from protos.TypeAndKey
func (tk TypeAndKey) FromProto(other protos.TypeAndKey) {
	tk.Type = other.Type
	tk.Key = other.Key
}

// ToProto returns a protos.TypeAndKey equivalent to this TypeAndKey
func (tk TypeAndKey) ToProto() protos.TypeAndKey {
	return protos.TypeAndKey{
		Type: tk.Type,
		Key:  tk.Key,
	}
}

func (tk TypeAndKey) String() string {
	return fmt.Sprintf("%s-%s", tk.Type, tk.Key)
}

// Blob encapsulates a blob for storage
type Blob struct {
	Type    string
	Key     string
	Value   []byte
	Version uint64
}

// BlobStorageFactory is an API to create a storage API bound to a transaction.
type BlobStorageFactory interface {
	// StartTransaction opens a transaction for all following blob storage
	// operations, and returns a TransactionalBlobStorage instance tied to the
	// opened transaction.
	StartTransaction() (TransactionalBlobStorage, error)
}

// TransactionalBlobStorage is the client API for blob storage operations
// within the context of a transaction.
type TransactionalBlobStorage interface {

	// Commit commits the existing transaction. If an error is returned from
	// the backing storage while committing, the transaction will be rolled
	// back.
	Commit() error

	// Rollback rolls back the existing transaction.
	Rollback() error

	// ListKeys returns all the blob keys stored for the network and type.
	ListKeys(networkID string, typeVal string) ([]string, error)

	// Get loads a specific blob from storage.
	// If there is no blob matching the given ID, ErrNotFound from
	// magma/orc8r/cloud/go/errors will be returned.
	Get(networkID string, id TypeAndKey) (Blob, error)

	// Get loads and returns a collection of blobs matching the specified IDs.
	// If there is no blob corresponding to a TypeAndKey, the returned list
	// will not have a corresponding Blob.
	GetMany(networkID string, ids []TypeAndKey) ([]Blob, error)

	// CreateOrUpdate writes blobs to the storage. Blobs are either updated
	// in-place or created. The Version field of Blobs passed in here is
	// ignored - all version incrementation is done internally inside the
	// storage implementation.
	CreateOrUpdate(networkID string, blobs []Blob) error

	// Delete deletes specified blobs from storage.
	Delete(networkID string, ids []TypeAndKey) error
}

// GetTableName returns the full table name for a networkID and a base table
func GetTableName(networkID string, baseTableName string) string {
	return fmt.Sprintf("%s_%s", networkID, baseTableName)
}

// GetBlobsByTypeAndKey returns a computed view of a list of blobs as a map of
// blobs keyed by blob TypeAndKey.
func GetBlobsByTypeAndKey(blobs []Blob) map[TypeAndKey]Blob {
	ret := make(map[TypeAndKey]Blob, len(blobs))
	for _, blob := range blobs {
		ret[TypeAndKey{Type: blob.Type, Key: blob.Key}] = blob
	}
	return ret
}
