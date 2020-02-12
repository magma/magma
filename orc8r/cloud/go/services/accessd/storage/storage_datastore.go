/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package storage

import (
	"magma/orc8r/cloud/go/datastore"
	accessprotos "magma/orc8r/cloud/go/services/accessd/protos"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	accessdTableDatastore = "access_control"
)

type accessdDatastore struct {
	store datastore.Api
}

// NewAccessdDatastore returns an initialized instance of accessdDatastore as AccessdStorage.
func NewAccessdDatastore(store datastore.Api) AccessdStorage {
	return &accessdDatastore{store: store}
}

func (a *accessdDatastore) ListAllIdentity() ([]*protos.Identity, error) {
	idHashes, err := a.store.ListKeys(accessdTableDatastore)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list keys: %s", err)
	}

	marshaledACLs, err := a.store.GetMany(accessdTableDatastore, idHashes)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get many acls: %s", err)
	}

	var ret []*protos.Identity
	for _, mACLWrapper := range marshaledACLs {
		acl := &accessprotos.AccessControl_List{}
		err = proto.Unmarshal(mACLWrapper.Value, acl)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to unmarshal acl: %s", err)
		}
		ret = append(ret, acl.Operator)
	}

	return ret, nil
}

func (a *accessdDatastore) GetACL(id *protos.Identity) (*accessprotos.AccessControl_List, error) {
	acls, err := a.GetManyACL([]*protos.Identity{id})
	if err != nil {
		return nil, err
	}
	for _, acl := range acls {
		return acl, nil
	}
	return nil, status.Errorf(codes.NotFound, "get ACL error for Operator %s: %s", id.HashString(), err)
}

func (a *accessdDatastore) GetManyACL(ids []*protos.Identity) ([]*accessprotos.AccessControl_List, error) {
	if ids == nil {
		return nil, status.Error(codes.InvalidArgument, "nil Identity list")
	}

	idHashes := make([]string, 0, len(ids))
	for _, id := range ids {
		if id == nil {
			return nil, status.Error(codes.InvalidArgument, "nil Identity")
		}
		idHashes = append(idHashes, id.HashString())
	}
	marshaledACLs, err := a.store.GetMany(accessdTableDatastore, idHashes)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get many acls: %s", err)
	}

	var ret []*accessprotos.AccessControl_List
	for _, mACLWrapper := range marshaledACLs {
		acl := &accessprotos.AccessControl_List{}
		err = proto.Unmarshal(mACLWrapper.Value, acl)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to unmarshal acl: %s", err)
		}
		ret = append(ret, acl)
	}

	return ret, nil
}

func (a *accessdDatastore) PutACL(id *protos.Identity, acl *accessprotos.AccessControl_List) error {
	if id == nil {
		return status.Error(codes.InvalidArgument, "nil Identity")
	}
	if acl == nil {
		return status.Error(codes.InvalidArgument, "nil AccessControl_List")
	}

	marshaledACL, err := proto.Marshal(acl)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to marshal acl: %s", err)
	}

	err = a.store.Put(accessdTableDatastore, id.HashString(), marshaledACL)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to put acl: %s", err)
	}

	return nil
}

// NOTE: datastore-implemented UpdateACLWithEntities is not atomic.
func (a *accessdDatastore) UpdateACLWithEntities(id *protos.Identity, entities []*accessprotos.AccessControl_Entity) error {
	if id == nil {
		return status.Error(codes.InvalidArgument, "nil Identity")
	}
	if entities == nil {
		return status.Error(codes.InvalidArgument, "nil AccessControl_Entity slice")
	}

	acl, err := a.GetACL(id)
	if err != nil {
		return err
	}

	err = accessprotos.AddToACL(acl, entities)
	if err != nil {
		return err
	}

	err = a.PutACL(id, acl)
	if err != nil {
		return err
	}

	return nil
}

func (a *accessdDatastore) DeleteACL(id *protos.Identity) error {
	if id == nil {
		return status.Error(codes.InvalidArgument, "nil Identity")
	}

	err := a.store.Delete(accessdTableDatastore, id.HashString())
	if err != nil {
		return status.Errorf(codes.Internal, "failed to delete acl: %s", err)
	}

	return nil
}
