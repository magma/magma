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
	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/protos"
	accessprotos "magma/orc8r/cloud/go/services/accessd/protos"
	astorage "magma/orc8r/cloud/go/services/accessd/storage"
	"magma/orc8r/cloud/go/storage"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAccessdBlobstore_ListAllIdentity(t *testing.T) {
	var blobFactMock *mocks.BlobStorageFactory
	var blobStoreMock *mocks.TransactionalBlobStorage
	someErr := errors.New("generic error")

	ids := []*protos.Identity{
		identity.NewOperator("test_operator_0"),
		identity.NewOperator("test_operator_1"),
	}
	idHashes := []string{
		ids[0].HashString(),
		ids[1].HashString(),
	}
	hashToId := protos.GetHashToIdentity(ids)
	perms := []accessprotos.AccessControl_Permission{
		accessprotos.AccessControl_READ | accessprotos.AccessControl_WRITE,
		accessprotos.AccessControl_READ,
	}
	entities := []map[string]*accessprotos.AccessControl_Entity{
		{idHashes[0]: {Id: ids[0], Permissions: perms[0]}},
		{idHashes[1]: {Id: ids[1], Permissions: perms[1]}},
	}
	acls := []*accessprotos.AccessControl_List{
		{Operator: ids[0], Entities: entities[0]},
		{Operator: ids[1], Entities: entities[1]},
	}
	marshaledACL0, err := proto.Marshal(acls[0])
	assert.NoError(t, err)
	marshaledACL1, err := proto.Marshal(acls[1])
	assert.NoError(t, err)

	tks := []storage.TypeAndKey{
		{Type: astorage.AccessdDefaultType, Key: idHashes[0]},
		{Type: astorage.AccessdDefaultType, Key: idHashes[1]},
	}
	blobs := []blobstore.Blob{
		{Type: astorage.AccessdDefaultType, Key: idHashes[0], Value: marshaledACL0},
		{Type: astorage.AccessdDefaultType, Key: idHashes[1], Value: marshaledACL1},
	}

	// Fail to start transaction
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := astorage.NewAccessdBlobstore(blobFactMock)

	_, err = store.ListAllIdentity()
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.ListKeys fails
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("ListKeys", mock.Anything, astorage.AccessdDefaultType).
		Return([]string{}, someErr).Once()
	store = astorage.NewAccessdBlobstore(blobFactMock)

	_, err = store.ListAllIdentity()
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.ListKeys succeeds with empty return
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("ListKeys", mock.Anything, astorage.AccessdDefaultType).Return([]string{}, nil).Once()
	blobStoreMock.On("GetMany", mock.Anything, []storage.TypeAndKey{}).
		Return([]blobstore.Blob{}, nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = astorage.NewAccessdBlobstore(blobFactMock)

	idsRecvd, err := store.ListAllIdentity()
	assert.NoError(t, err)
	assert.Empty(t, idsRecvd)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.GetMany fails
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("ListKeys", mock.Anything, astorage.AccessdDefaultType).Return(idHashes, nil).Once()
	blobStoreMock.On("GetMany", mock.Anything, tks).Return([]blobstore.Blob{}, someErr).Once()
	store = astorage.NewAccessdBlobstore(blobFactMock)

	_, err = store.ListAllIdentity()
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Success
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("ListKeys", mock.Anything, astorage.AccessdDefaultType).Return(idHashes, nil).Once()
	blobStoreMock.On("GetMany", mock.Anything, tks).Return(blobs, nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = astorage.NewAccessdBlobstore(blobFactMock)

	idsRecvd, err = store.ListAllIdentity()
	assert.NoError(t, err)
	assert.Len(t, idsRecvd, 2)
	for _, idRecvd := range idsRecvd {
		assert.True(t, proto.Equal(hashToId[idRecvd.HashString()], idRecvd))
	}
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}

func TestAccessdBlobstore_GetACL(t *testing.T) {
	var blobFactMock *mocks.BlobStorageFactory
	var blobStoreMock *mocks.TransactionalBlobStorage
	someErr := errors.New("generic error")

	id := identity.NewOperator("testOperator")
	idHash := id.HashString()
	perm := accessprotos.AccessControl_READ | accessprotos.AccessControl_WRITE
	entities := map[string]*accessprotos.AccessControl_Entity{idHash: {Id: id, Permissions: perm}}
	acl := &accessprotos.AccessControl_List{Operator: id, Entities: entities}

	marshaledACL, err := proto.Marshal(acl)
	assert.NoError(t, err)
	tks := []storage.TypeAndKey{{Type: astorage.AccessdDefaultType, Key: idHash}}
	blobs := []blobstore.Blob{{Type: astorage.AccessdDefaultType, Key: idHash, Value: marshaledACL}}

	// Fail to start transaction
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := astorage.NewAccessdBlobstore(blobFactMock)

	_, err = store.GetACL(id)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Call with nil id
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	store = astorage.NewAccessdBlobstore(blobFactMock)

	_, err = store.GetACL(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidArgument")
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.Get fails with ErrNotFound
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("GetMany", mock.Anything, tks).Return([]blobstore.Blob{}, nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = astorage.NewAccessdBlobstore(blobFactMock)

	_, err = store.GetACL(id)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "NotFound")
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.Get fails with error other than ErrNotFound
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("GetMany", mock.Anything, tks).Return([]blobstore.Blob{}, someErr).Once()
	store = astorage.NewAccessdBlobstore(blobFactMock)

	_, err = store.GetACL(id)
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
	store = astorage.NewAccessdBlobstore(blobFactMock)

	aclRecvd, err := store.GetACL(id)
	assert.NoError(t, err)
	assert.True(t, proto.Equal(acl, aclRecvd))
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}

func TestAccessdBlobstore_GetManyACL(t *testing.T) {
	var blobFactMock *mocks.BlobStorageFactory
	var blobStoreMock *mocks.TransactionalBlobStorage
	someErr := errors.New("generic error")

	ids := []*protos.Identity{
		identity.NewOperator("test_operator_0"),
		identity.NewOperator("test_operator_1"),
	}
	idHashes := []string{
		ids[0].HashString(),
		ids[1].HashString(),
	}
	perms := []accessprotos.AccessControl_Permission{
		accessprotos.AccessControl_READ | accessprotos.AccessControl_WRITE,
		accessprotos.AccessControl_READ,
	}
	entities := []map[string]*accessprotos.AccessControl_Entity{
		{idHashes[0]: {Id: ids[0], Permissions: perms[0]}},
		{idHashes[1]: {Id: ids[1], Permissions: perms[1]}},
	}
	acls := []*accessprotos.AccessControl_List{
		{Operator: ids[0], Entities: entities[0]},
		{Operator: ids[1], Entities: entities[1]},
	}
	marshaledACL0, err := proto.Marshal(acls[0])
	assert.NoError(t, err)
	marshaledACL1, err := proto.Marshal(acls[1])
	assert.NoError(t, err)

	tks := []storage.TypeAndKey{
		{Type: astorage.AccessdDefaultType, Key: idHashes[0]},
		{Type: astorage.AccessdDefaultType, Key: idHashes[1]},
	}
	blobs := []blobstore.Blob{
		{Type: astorage.AccessdDefaultType, Key: idHashes[0], Value: marshaledACL0},
		{Type: astorage.AccessdDefaultType, Key: idHashes[1], Value: marshaledACL1},
	}
	aclsByidHash := map[string]*accessprotos.AccessControl_List{
		idHashes[0]: acls[0],
		idHashes[1]: acls[1],
	}

	// Call with nil ids
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	store := astorage.NewAccessdBlobstore(blobFactMock)

	_, err = store.GetManyACL(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidArgument")
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Fail to start transaction
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store = astorage.NewAccessdBlobstore(blobFactMock)

	_, err = store.GetManyACL(ids)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.GetMany fails
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("GetMany", mock.Anything, tks).Return([]blobstore.Blob{}, someErr).Once()
	store = astorage.NewAccessdBlobstore(blobFactMock)

	_, err = store.GetManyACL(ids)
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
	store = astorage.NewAccessdBlobstore(blobFactMock)

	aclsRecvd, err := store.GetManyACL(ids)
	assert.NoError(t, err)
	assert.Empty(t, aclsRecvd)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Success
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("GetMany", mock.Anything, tks).Return(blobs, nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = astorage.NewAccessdBlobstore(blobFactMock)

	aclsRecvd, err = store.GetManyACL(ids)
	assert.NoError(t, err)
	for _, aclRecvd := range aclsRecvd {
		assert.True(t, proto.Equal(aclRecvd, aclsByidHash[aclRecvd.Operator.HashString()]))
	}
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}

func TestAccessdBlobstore_PutACL(t *testing.T) {
	var blobFactMock *mocks.BlobStorageFactory
	var blobStoreMock *mocks.TransactionalBlobStorage
	someErr := errors.New("generic error")

	id := identity.NewOperator("testOperator")
	idHash := id.HashString()
	perm := accessprotos.AccessControl_READ | accessprotos.AccessControl_WRITE
	entities := map[string]*accessprotos.AccessControl_Entity{idHash: {Id: id, Permissions: perm}}
	acl := &accessprotos.AccessControl_List{Operator: id, Entities: entities}

	marshaledACL, err := proto.Marshal(acl)
	assert.NoError(t, err)
	blobs := []blobstore.Blob{{Type: astorage.AccessdDefaultType, Key: idHash, Value: marshaledACL}}

	// Call with nil id
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	store := astorage.NewAccessdBlobstore(blobFactMock)

	err = store.PutACL(nil, acl)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidArgument")
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Call with nil acl
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	store = astorage.NewAccessdBlobstore(blobFactMock)

	err = store.PutACL(id, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidArgument")
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Fail to start transaction
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store = astorage.NewAccessdBlobstore(blobFactMock)

	err = store.PutACL(id, acl)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.Put fails
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("CreateOrUpdate", mock.Anything, blobs).Return(someErr).Once()
	store = astorage.NewAccessdBlobstore(blobFactMock)

	err = store.PutACL(id, acl)
	assert.Error(t, err)
	blobStoreMock.AssertExpectations(t)

	// Success
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("CreateOrUpdate", mock.Anything, blobs).
		Return(nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = astorage.NewAccessdBlobstore(blobFactMock)

	err = store.PutACL(id, acl)
	assert.NoError(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}

func TestAccessdBlobstore_UpdateACLWithEntities(t *testing.T) {
	var blobFactMock *mocks.BlobStorageFactory
	var blobStoreMock *mocks.TransactionalBlobStorage
	someErr := errors.New("generic error")

	id := identity.NewOperator("testOperator")
	idHash := id.HashString()
	perm := accessprotos.AccessControl_READ | accessprotos.AccessControl_WRITE
	entities := []*accessprotos.AccessControl_Entity{{Id: id, Permissions: perm}}
	hashToEntity := map[string]*accessprotos.AccessControl_Entity{idHash: entities[0]}
	aclInitial := &accessprotos.AccessControl_List{Operator: id}
	aclFinal := &accessprotos.AccessControl_List{Operator: id, Entities: hashToEntity}

	marshaledACLInitial, err := proto.Marshal(aclInitial)
	assert.NoError(t, err)
	marshaledACLFinal, err := proto.Marshal(aclFinal)
	assert.NoError(t, err)
	tk := storage.TypeAndKey{Type: astorage.AccessdDefaultType, Key: idHash}
	blobsInitial := []blobstore.Blob{{Type: astorage.AccessdDefaultType, Key: idHash, Value: marshaledACLInitial}}
	blobsFinal := []blobstore.Blob{{Type: astorage.AccessdDefaultType, Key: idHash, Value: marshaledACLFinal}}

	// Fail to start transaction
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := astorage.NewAccessdBlobstore(blobFactMock)

	err = store.UpdateACLWithEntities(id, entities)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Call with nil id
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	store = astorage.NewAccessdBlobstore(blobFactMock)

	err = store.UpdateACLWithEntities(nil, entities)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidArgument")
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Call with nil entities slice
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	store = astorage.NewAccessdBlobstore(blobFactMock)

	err = store.UpdateACLWithEntities(id, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidArgument")
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Call with nil element in entities slice
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Get", mock.Anything, tk).Return(blobsInitial[0], nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	store = astorage.NewAccessdBlobstore(blobFactMock)

	err = store.UpdateACLWithEntities(id, []*accessprotos.AccessControl_Entity{nil})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidArgument")
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.Get fails
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Get", mock.Anything, tk).Return(blobstore.Blob{}, someErr).Once()
	store = astorage.NewAccessdBlobstore(blobFactMock)

	err = store.UpdateACLWithEntities(id, entities)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.Put fails
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Get", mock.Anything, tk).Return(blobsInitial[0], nil).Once()
	blobStoreMock.On("CreateOrUpdate", mock.Anything, blobsFinal).Return(someErr).Once()
	store = astorage.NewAccessdBlobstore(blobFactMock)

	err = store.UpdateACLWithEntities(id, entities)
	assert.Error(t, err)
	blobStoreMock.AssertExpectations(t)

	// Success
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Get", mock.Anything, tk).Return(blobsInitial[0], nil).Once()
	blobStoreMock.On("CreateOrUpdate", mock.Anything, blobsFinal).Return(nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = astorage.NewAccessdBlobstore(blobFactMock)

	err = store.UpdateACLWithEntities(id, entities)
	assert.NoError(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}

func TestAccessdBlobstore_DeleteACL(t *testing.T) {
	var blobFactMock *mocks.BlobStorageFactory
	var blobStoreMock *mocks.TransactionalBlobStorage
	someErr := errors.New("generic error")

	id := identity.NewOperator("testOperator")
	tks := []storage.TypeAndKey{{Type: astorage.AccessdDefaultType, Key: id.HashString()}}

	// Call with nil id
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	store := astorage.NewAccessdBlobstore(blobFactMock)

	err := store.DeleteACL(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidArgument")
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Fail to start transaction
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store = astorage.NewAccessdBlobstore(blobFactMock)

	err = store.DeleteACL(id)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.Delete fails
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Delete", mock.Anything, tks).Return(someErr).Once()
	store = astorage.NewAccessdBlobstore(blobFactMock)

	err = store.DeleteACL(id)
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
	store = astorage.NewAccessdBlobstore(blobFactMock)

	err = store.DeleteACL(id)
	assert.NoError(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}
