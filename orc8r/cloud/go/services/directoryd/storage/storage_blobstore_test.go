/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package storage_test

import (
	"errors"
	"testing"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/blobstore/mocks"
	dstorage "magma/orc8r/cloud/go/services/directoryd/storage"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	networkID = "placeholder_network"
)

func TestDirectorydBlobstoreStorage_GetHostname(t *testing.T) {
	var blobFactMock *mocks.BlobStorageFactory
	var blobStoreMock *mocks.TransactionalBlobStorage
	someErr := errors.New("generic error")

	hwid := "some_hwid"
	tk := storage.TypeAndKey{Type: dstorage.DirectorydDefaultType, Key: hwid}

	hostname := "some_hostname"
	blob := blobstore.Blob{
		Type:    tk.Type,
		Key:     tk.Key,
		Value:   []byte(hostname),
		Version: 0,
	}

	// Fail to start transaction
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := dstorage.NewDirectorydBlobstore(blobFactMock)

	_, err := store.GetHostname(hwid)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.Get fails with ErrNotFound
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Get", networkID, tk).Return(blobstore.Blob{}, merrors.ErrNotFound).Once()
	store = dstorage.NewDirectorydBlobstore(blobFactMock)

	_, err = store.GetHostname(hwid)
	assert.Exactly(t, merrors.ErrNotFound, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.Get fails with error other than ErrNotFound
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Get", networkID, tk).Return(blobstore.Blob{}, someErr).Once()
	store = dstorage.NewDirectorydBlobstore(blobFactMock)

	_, err = store.GetHostname(hwid)
	assert.Error(t, err)
	assert.NotEqual(t, merrors.ErrNotFound, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Success
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Get", networkID, tk).Return(blob, nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = dstorage.NewDirectorydBlobstore(blobFactMock)

	hostnameRecvd, err := store.GetHostname(hwid)
	assert.Equal(t, hostname, hostnameRecvd)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}

func TestDirectorydBlobstoreStorage_PutRecord(t *testing.T) {
	var blobFactMock *mocks.BlobStorageFactory
	var blobStoreMock *mocks.TransactionalBlobStorage
	someErr := errors.New("generic error")

	hwid := "some_hwid"
	tk := storage.TypeAndKey{Type: dstorage.DirectorydDefaultType, Key: hwid}

	hostname := "some_hostname"
	blob := blobstore.Blob{
		Type:    tk.Type,
		Key:     tk.Key,
		Value:   []byte(hostname),
		Version: 0,
	}

	// Fail to start transaction
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := dstorage.NewDirectorydBlobstore(blobFactMock)

	err := store.PutHostname(hwid, hostname)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.PutRecord fails
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("CreateOrUpdate", networkID, mock.Anything, mock.Anything).
		Return(someErr).Once()
	store = dstorage.NewDirectorydBlobstore(blobFactMock)

	err = store.PutHostname(hwid, hostname)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Success
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("CreateOrUpdate", networkID, []blobstore.Blob{blob}).
		Return(nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = dstorage.NewDirectorydBlobstore(blobFactMock)

	err = store.PutHostname(hwid, hostname)
	assert.NoError(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}
