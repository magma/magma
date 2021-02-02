/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package storage

import (
	"magma/orc8r/cloud/go/blobstore"
	accessprotos "magma/orc8r/cloud/go/services/accessd/protos"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// AccessdTableBlobstore is the table where blobstore stores accessd records.
	AccessdTableBlobstore = "access_control_blobstore"

	// AccessdDefaultType is the default type blobstore uses for accessd protos.
	AccessdDefaultType = "access_control"

	// Blobstore needs a network ID, but accessd is network-agnostic so we
	// will use a placeholder value.
	placeholderNetworkID = "placeholder_network"
)

type accessdBlobstore struct {
	factory blobstore.BlobStorageFactory
}

// NewAccessdBlobstore returns an initialized instance of accessdBlobstore as AccessdStorage.
func NewAccessdBlobstore(factory blobstore.BlobStorageFactory) AccessdStorage {
	return &accessdBlobstore{factory: factory}
}

func (a *accessdBlobstore) ListAllIdentity() ([]*protos.Identity, error) {
	var ids []*protos.Identity

	store, err := a.factory.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "failed to start transaction: %s", err)
	}
	defer store.Rollback()

	idHashes, err := blobstore.ListKeys(store, placeholderNetworkID, AccessdDefaultType)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list keys: %s", err)
	}

	if len(idHashes) == 0 {
		return ids, store.Commit()
	}

	tks := storage.MakeTKs(AccessdDefaultType, idHashes)
	blobs, err := store.GetMany(placeholderNetworkID, tks)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get many acls: %s", err)
	}

	ids = make([]*protos.Identity, 0, len(tks))
	for _, blob := range blobs {
		acl := &accessprotos.AccessControl_List{}
		err = proto.Unmarshal(blob.Value, acl)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to unmarshal acl: %s", err)
		}
		ids = append(ids, acl.Operator)
	}
	err = store.Commit()
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "failed to commit transaction: %s", err)
	}
	return ids, nil
}

func (a *accessdBlobstore) GetACL(id *protos.Identity) (*accessprotos.AccessControl_List, error) {
	acls, err := a.GetManyACL([]*protos.Identity{id})
	if err != nil {
		return nil, err
	}
	for _, acl := range acls {
		return acl, nil
	}
	return nil, status.Errorf(codes.NotFound, "get ACL error for Operator %s: %s", id.HashString(), err)
}

func (a *accessdBlobstore) GetManyACL(ids []*protos.Identity) ([]*accessprotos.AccessControl_List, error) {
	if ids == nil {
		return nil, status.Error(codes.InvalidArgument, "nil Identity list")
	}

	store, err := a.factory.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "failed to start transaction: %s", err)
	}
	defer store.Rollback()

	idHashes := make([]string, 0, len(ids))
	for _, id := range ids {
		if id == nil {
			return nil, status.Errorf(codes.InvalidArgument, "nil Identity: %s", err)
		}
		idHashes = append(idHashes, id.HashString())
	}
	tks := storage.MakeTKs(AccessdDefaultType, idHashes)
	blobs, err := store.GetMany(placeholderNetworkID, tks)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get many acls: %s", err)
	}

	ret := make([]*accessprotos.AccessControl_List, 0, len(idHashes))
	for _, blob := range blobs {
		acl := &accessprotos.AccessControl_List{}
		err = proto.Unmarshal(blob.Value, acl)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to unmarshal acl: %s", err)
		}
		ret = append(ret, acl)
	}

	err = store.Commit()
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "failed to commit transaction: %s", err)
	}
	return ret, nil
}

func (a *accessdBlobstore) PutACL(id *protos.Identity, acl *accessprotos.AccessControl_List) error {
	if id == nil {
		return status.Error(codes.InvalidArgument, "nil Identity")
	}
	if acl == nil {
		return status.Error(codes.InvalidArgument, "nil AccessControl_List")
	}

	store, err := a.factory.StartTransaction(&storage.TxOptions{})
	if err != nil {
		return status.Errorf(codes.Unavailable, "failed to start transaction: %s", err)
	}
	defer store.Rollback()

	marshaledACL, err := proto.Marshal(acl)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to marshal acl: %s", err)
	}

	blob := blobstore.Blob{Type: AccessdDefaultType, Key: id.HashString(), Value: marshaledACL}
	err = store.CreateOrUpdate(placeholderNetworkID, blobstore.Blobs{blob})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to put acl: %s", err)
	}

	err = store.Commit()
	if err != nil {
		return status.Errorf(codes.Unavailable, "failed to commit transaction: %s", err)
	}
	return nil
}

func (a *accessdBlobstore) UpdateACLWithEntities(id *protos.Identity, entities []*accessprotos.AccessControl_Entity) error {
	if id == nil {
		return status.Error(codes.InvalidArgument, "nil Identity")
	}
	if entities == nil {
		return status.Error(codes.InvalidArgument, "nil AccessControl_Entity slice")
	}

	store, err := a.factory.StartTransaction(&storage.TxOptions{})
	if err != nil {
		return status.Errorf(codes.Unavailable, "failed to start transaction: %s", err)
	}
	defer store.Rollback()

	blobRecvd, err := store.Get(placeholderNetworkID, storage.TypeAndKey{Type: AccessdDefaultType, Key: id.HashString()})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to get acl: %s", err)
	}

	acl := &accessprotos.AccessControl_List{}
	err = proto.Unmarshal(blobRecvd.Value, acl)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to unmarshal acl: %s", err)
	}

	err = accessprotos.AddToACL(acl, entities)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to update acl with entities: %s", err)
	}

	marshaledACL, err := proto.Marshal(acl)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to marshal acl: %s", err)
	}

	blobPut := blobstore.Blob{Type: AccessdDefaultType, Key: id.HashString(), Value: marshaledACL}
	err = store.CreateOrUpdate(placeholderNetworkID, blobstore.Blobs{blobPut})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to put acl: %s", err)
	}

	err = store.Commit()
	if err != nil {
		return status.Errorf(codes.Unavailable, "failed to commit transaction: %s", err)
	}
	return nil

}

func (a *accessdBlobstore) DeleteACL(id *protos.Identity) error {
	if id == nil {
		return status.Error(codes.InvalidArgument, "nil Identity")
	}

	store, err := a.factory.StartTransaction(&storage.TxOptions{})
	if err != nil {
		return status.Errorf(codes.Unavailable, "failed to start transaction: %s", err)
	}
	defer store.Rollback()

	tk := storage.TypeAndKey{Type: AccessdDefaultType, Key: id.HashString()}
	err = store.Delete(placeholderNetworkID, []storage.TypeAndKey{tk})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to delete acl: %s", err)
	}

	err = store.Commit()
	if err != nil {
		return status.Errorf(codes.Unavailable, "failed to commit transaction: %s", err)
	}
	return nil
}
