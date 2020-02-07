/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package storage_test

import (
	"errors"
	"testing"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/blobstore/mocks"
	merrors "magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/services/certifier/protos"
	cstorage "magma/orc8r/cloud/go/services/certifier/storage"
	"magma/orc8r/cloud/go/storage"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCertifierBlobstore_GetCertInfo(t *testing.T) {
	var blobFactMock *mocks.BlobStorageFactory
	var blobStoreMock *mocks.TransactionalBlobStorage
	someErr := errors.New("generic error")

	serialNumber := "serial_number"
	info := &protos.CertificateInfo{
		Id:        nil,
		NotBefore: &timestamp.Timestamp{Seconds: 0xdead, Nanos: 0xbeef},
		NotAfter:  &timestamp.Timestamp{Seconds: 0xbeef, Nanos: 0xfeed},
		CertType:  1,
	}
	marshaledInfo, err := proto.Marshal(info)
	assert.NoError(t, err)
	tks := []storage.TypeAndKey{{Type: cstorage.CertInfoType, Key: serialNumber}}
	blobs := []blobstore.Blob{{Type: cstorage.CertInfoType, Key: serialNumber, Value: marshaledInfo}}

	// Fail to start transaction
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := cstorage.NewCertifierBlobstore(blobFactMock)

	_, err = store.GetCertInfo(serialNumber)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.Get fails with ErrNotFound
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("GetMany", mock.Anything, tks).Return([]blobstore.Blob{}, nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	_, err = store.GetCertInfo(serialNumber)
	assert.Equal(t, err, merrors.ErrNotFound)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.Get fails with error other than ErrNotFound
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("GetMany", mock.Anything, tks).Return([]blobstore.Blob{}, someErr).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	_, err = store.GetCertInfo(serialNumber)
	assert.Error(t, err)
	assert.NotEqual(t, err, merrors.ErrNotFound)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Success
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("GetMany", mock.Anything, tks).Return(blobs, nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	infoRecvd, err := store.GetCertInfo(serialNumber)
	assert.NoError(t, err)
	assert.True(t, proto.Equal(info, infoRecvd))
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}

func TestCertifierBlobstore_GetMany(t *testing.T) {
	var blobFactMock *mocks.BlobStorageFactory
	var blobStoreMock *mocks.TransactionalBlobStorage
	someErr := errors.New("generic error")

	info0 := &protos.CertificateInfo{
		Id:        nil,
		NotBefore: &timestamp.Timestamp{Seconds: 0xdead, Nanos: 0xbeef},
		NotAfter:  &timestamp.Timestamp{Seconds: 0xbeef, Nanos: 0xfeed},
		CertType:  1,
	}
	info1 := &protos.CertificateInfo{
		Id:        nil,
		NotBefore: &timestamp.Timestamp{Seconds: 0xaaaa, Nanos: 0xbbbb},
		NotAfter:  &timestamp.Timestamp{Seconds: 0xcccc, Nanos: 0xdddd},
		CertType:  2,
	}
	marshaledInfo0, err := proto.Marshal(info0)
	assert.NoError(t, err)
	marshaledInfo1, err := proto.Marshal(info1)
	assert.NoError(t, err)

	serialNumbers := []string{"serial_number_0", "serial_number_1"}
	tks := []storage.TypeAndKey{
		{Type: cstorage.CertInfoType, Key: serialNumbers[0]},
		{Type: cstorage.CertInfoType, Key: serialNumbers[1]},
	}
	blobs := []blobstore.Blob{
		{Type: cstorage.CertInfoType, Key: serialNumbers[0], Value: marshaledInfo0},
		{Type: cstorage.CertInfoType, Key: serialNumbers[1], Value: marshaledInfo1},
	}
	infos := map[string]*protos.CertificateInfo{
		serialNumbers[0]: info0,
		serialNumbers[1]: info1,
	}

	// Fail to start transaction
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := cstorage.NewCertifierBlobstore(blobFactMock)

	_, err = store.GetManyCertInfo(serialNumbers)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.GetMany fails
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("GetMany", mock.Anything, tks).Return([]blobstore.Blob{}, someErr).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	_, err = store.GetManyCertInfo(serialNumbers)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.GetMany succeeds with empty return
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("GetMany", mock.Anything, tks).Return([]blobstore.Blob{}, nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	infosRecvds, err := store.GetManyCertInfo(serialNumbers)
	assert.NoError(t, err)
	assert.Equal(t, map[string]*protos.CertificateInfo{}, infosRecvds)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Success
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("GetMany", mock.Anything, tks).Return(blobs, nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	infosRecvds, err = store.GetManyCertInfo(serialNumbers)
	assert.NoError(t, err)
	assert.Len(t, infosRecvds, 2)
	for k := range infos {
		a, b := infos[k], infosRecvds[k]
		assert.True(t, proto.Equal(a, b))
	}
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}

func TestCertifierBlobstore_GetAll(t *testing.T) {
	var blobFactMock *mocks.BlobStorageFactory
	var blobStoreMock *mocks.TransactionalBlobStorage
	someErr := errors.New("generic error")

	info0 := &protos.CertificateInfo{
		Id:        nil,
		NotBefore: &timestamp.Timestamp{Seconds: 0xdead, Nanos: 0xbeef},
		NotAfter:  &timestamp.Timestamp{Seconds: 0xbeef, Nanos: 0xfeed},
		CertType:  1,
	}
	info1 := &protos.CertificateInfo{
		Id:        nil,
		NotBefore: &timestamp.Timestamp{Seconds: 0xaaaa, Nanos: 0xbbbb},
		NotAfter:  &timestamp.Timestamp{Seconds: 0xcccc, Nanos: 0xdddd},
		CertType:  2,
	}
	marshaledInfo0, err := proto.Marshal(info0)
	assert.NoError(t, err)
	marshaledInfo1, err := proto.Marshal(info1)
	assert.NoError(t, err)

	serialNumbers := []string{"serial_number_0", "serial_number_1"}
	tks := []storage.TypeAndKey{
		{Type: cstorage.CertInfoType, Key: serialNumbers[0]},
		{Type: cstorage.CertInfoType, Key: serialNumbers[1]},
	}
	blobs := []blobstore.Blob{
		{Type: cstorage.CertInfoType, Key: serialNumbers[0], Value: marshaledInfo0},
		{Type: cstorage.CertInfoType, Key: serialNumbers[1], Value: marshaledInfo1},
	}
	infos := map[string]*protos.CertificateInfo{
		serialNumbers[0]: info0,
		serialNumbers[1]: info1,
	}

	// Fail to start transaction
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := cstorage.NewCertifierBlobstore(blobFactMock)

	_, err = store.GetAllCertInfo()
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.ListKeys fails
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("ListKeys", mock.Anything, cstorage.CertInfoType).
		Return([]string{}, someErr).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	_, err = store.GetAllCertInfo()
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.ListKeys succeeds with empty return
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("ListKeys", mock.Anything, cstorage.CertInfoType).Return([]string{}, nil).Once()
	blobStoreMock.On("GetMany", mock.Anything, []storage.TypeAndKey{}).Return([]blobstore.Blob{}, nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	infosRecvds, err := store.GetAllCertInfo()
	assert.NoError(t, err)
	assert.Equal(t, map[string]*protos.CertificateInfo{}, infosRecvds)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.GetMany fails
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("ListKeys", mock.Anything, cstorage.CertInfoType).Return(serialNumbers, nil).Once()
	blobStoreMock.On("GetMany", mock.Anything, tks).Return([]blobstore.Blob{}, someErr).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	_, err = store.GetAllCertInfo()
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Success
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("ListKeys", mock.Anything, cstorage.CertInfoType).Return(serialNumbers, nil).Once()
	blobStoreMock.On("GetMany", mock.Anything, tks).Return(blobs, nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	infosRecvds, err = store.GetAllCertInfo()
	assert.NoError(t, err)
	assert.Len(t, infosRecvds, 2)
	for k := range infos {
		a, b := infos[k], infosRecvds[k]
		assert.True(t, proto.Equal(a, b))
	}
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}

func TestCertifierBlobstore_Put(t *testing.T) {
	var blobFactMock *mocks.BlobStorageFactory
	var blobStoreMock *mocks.TransactionalBlobStorage
	someErr := errors.New("generic error")

	serialNumber := "serial_number"
	info := &protos.CertificateInfo{
		Id:        nil,
		NotBefore: &timestamp.Timestamp{Seconds: 0xdead, Nanos: 0xbeef},
		NotAfter:  &timestamp.Timestamp{Seconds: 0xbeef, Nanos: 0xfeed},
		CertType:  1,
	}
	marshaledInfos, err := proto.Marshal(info)
	assert.NoError(t, err)
	blob := blobstore.Blob{Type: cstorage.CertInfoType, Key: serialNumber, Value: marshaledInfos}

	// Fail to start transaction
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := cstorage.NewCertifierBlobstore(blobFactMock)

	err = store.PutCertInfo(serialNumber, info)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.Put fails
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("CreateOrUpdate", mock.Anything, []blobstore.Blob{blob}).Return(someErr).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	err = store.PutCertInfo(serialNumber, info)
	assert.Error(t, err)
	blobStoreMock.AssertExpectations(t)

	// Success
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("CreateOrUpdate", mock.Anything, []blobstore.Blob{blob}).
		Return(nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	err = store.PutCertInfo(serialNumber, info)
	assert.NoError(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}

func TestCertifierBlobstore_Delete(t *testing.T) {
	var blobFactMock *mocks.BlobStorageFactory
	var blobStoreMock *mocks.TransactionalBlobStorage
	someErr := errors.New("generic error")

	serialNumber := "serial_number"
	tks := []storage.TypeAndKey{storage.TypeAndKey{Type: cstorage.CertInfoType, Key: serialNumber}}

	// Fail to start transaction
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := cstorage.NewCertifierBlobstore(blobFactMock)

	err := store.DeleteCertInfo(serialNumber)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.Delete fails
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Delete", mock.Anything, tks).Return(someErr).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	err = store.DeleteCertInfo(serialNumber)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Success
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Delete", mock.Anything, tks).Return(nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	err = store.DeleteCertInfo(serialNumber)
	assert.NoError(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}

func TestCertifierBlobstore_GetSerialNumbers(t *testing.T) {
	var blobFactMock *mocks.BlobStorageFactory
	var blobStoreMock *mocks.TransactionalBlobStorage
	someErr := errors.New("generic error")
	serialNumbers := []string{"serial_number_0", "serial_number_1"}

	// Fail to start transaction
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := cstorage.NewCertifierBlobstore(blobFactMock)

	_, err := store.ListSerialNumbers()
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.GetSerialNumbers fails
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("ListKeys", mock.Anything, cstorage.CertInfoType).
		Return([]string{}, someErr).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	_, err = store.ListSerialNumbers()
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Success
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("ListKeys", mock.Anything, cstorage.CertInfoType).
		Return(serialNumbers, nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	serialNumbersRecvd, err := store.ListSerialNumbers()
	assert.NoError(t, err)
	assert.Equal(t, serialNumbers, serialNumbersRecvd)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}
