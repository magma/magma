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

	"github.com/pkg/errors"
)

const (
	// DirectorydTableBlobstore is the table where blobstore stores directoryd's hwid_to_hostname data.
	DirectorydTableBlobstore = "directoryd_blobstore"

	// DirectorydDefaultType is the default type field for blobstore storage.
	DirectorydDefaultType = "hwid_to_hostname"

	// Blobstore needs a network ID, but directoryd is network-agnostic so we
	// use a placeholder value.
	placeholderNetworkID = "placeholder_network"
)

// NewDirectorydBlobstore returns a directoryd storage implementation
// backed by the provided blobstore factory.
// NOTE: the datastore impl uses tableID as the table name, while here the
// blobstore impl uses tableID as the type field within a single table.
func NewDirectorydBlobstore(factory blobstore.BlobStorageFactory) DirectorydStorage {
	return &directorydBlobstore{factory: factory}
}

type directorydBlobstore struct {
	factory blobstore.BlobStorageFactory
}

func (d *directorydBlobstore) GetHostname(hwid string) (string, error) {
	store, err := d.factory.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return "", errors.Wrap(err, "failed to start transaction")
	}
	defer store.Rollback()

	blob, err := store.Get(
		placeholderNetworkID,
		storage.TypeAndKey{Type: DirectorydDefaultType, Key: hwid},
	)
	if err == merrors.ErrNotFound {
		return "", err
	}
	if err != nil {
		return "", errors.Wrap(err, "failed to get hostname")
	}

	hostname := string(blob.Value)
	return hostname, store.Commit()
}

func (d *directorydBlobstore) PutHostname(hwid, hostname string) error {
	store, err := d.factory.StartTransaction(&storage.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	defer store.Rollback()

	blob := blobstore.Blob{Type: DirectorydDefaultType, Key: hwid, Value: []byte(hostname)}
	err = store.CreateOrUpdate(placeholderNetworkID, []blobstore.Blob{blob})
	if err != nil {
		return errors.Wrap(err, "failed to create or update location record")
	}
	return store.Commit()
}
