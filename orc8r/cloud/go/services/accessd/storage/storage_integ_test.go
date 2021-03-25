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

package storage_test

import (
	"testing"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/identity"
	accessprotos "magma/orc8r/cloud/go/services/accessd/protos"
	"magma/orc8r/cloud/go/services/accessd/storage"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/proto"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestAccessdStorageBlobstore_Integation(t *testing.T) {
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	fact := blobstore.NewEntStorage(storage.AccessdTableBlobstore, db, sqorc.GetSqlBuilder())
	err = fact.InitializeFactory()
	assert.NoError(t, err)
	store := storage.NewAccessdBlobstore(fact)
	testAccessdStorageImpl(t, store)
}

func testAccessdStorageImpl(t *testing.T, store storage.AccessdStorage) {
	ids := []*protos.Identity{
		identity.NewOperator("test_operator_0"),
		identity.NewOperator("test_operator_1"),
		identity.NewOperator("test_operator_2"),
	}
	idHashes := []string{
		ids[0].HashString(),
		ids[1].HashString(),
		ids[2].HashString(),
	}
	perms := []accessprotos.AccessControl_Permission{
		accessprotos.AccessControl_READ | accessprotos.AccessControl_WRITE,
		accessprotos.AccessControl_READ,
		accessprotos.AccessControl_WRITE,
	}
	entities := []map[string]*accessprotos.AccessControl_Entity{
		{idHashes[0]: {Id: ids[0], Permissions: perms[0]}},
		{idHashes[1]: {Id: ids[1], Permissions: perms[1]}},
		{idHashes[2]: {Id: ids[2], Permissions: perms[2]}},
	}
	acls := map[string]*accessprotos.AccessControl_List{
		idHashes[0]: {Operator: ids[0], Entities: entities[0]},
		idHashes[1]: {Operator: ids[1], Entities: entities[1]},
		idHashes[2]: {Operator: ids[2], Entities: entities[2]},
	}

	// Empty initially
	idsRecvd, err := store.ListAllIdentity()
	assert.NoError(t, err)
	assert.Len(t, idsRecvd, 0)

	// Put and Get acl0
	err = store.PutACL(ids[0], acls[idHashes[0]])
	assert.NoError(t, err)

	aclRecvd, err := store.GetACL(ids[0])
	assert.NoError(t, err)
	assert.True(t, proto.Equal(acls[idHashes[0]], aclRecvd))

	// Put and Get acl1
	err = store.PutACL(ids[1], acls[idHashes[1]])
	assert.NoError(t, err)

	aclRecvd, err = store.GetACL(ids[1])
	assert.NoError(t, err)
	assert.True(t, proto.Equal(acls[idHashes[1]], aclRecvd))

	// Put acl2, GetMany acls 0 and 1
	err = store.PutACL(ids[2], acls[idHashes[2]])
	assert.NoError(t, err)
	aclsRecvd, err := store.GetManyACL(ids[0:2])
	assert.NoError(t, err)
	assert.Len(t, aclsRecvd, 2)

	for _, acl := range aclsRecvd {
		opkeyRecvd := acl.Operator.HashString()
		assert.True(t, proto.Equal(acls[opkeyRecvd], acl))
	}

	// Delete acl0, Get acl0, GetMany acls 0 and 1
	err = store.DeleteACL(ids[0])
	assert.NoError(t, err)
	_, err = store.GetACL(ids[0])
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "NotFound")
	aclsRecvd, err = store.GetManyACL(ids[0:2])
	assert.NoError(t, err)
	assert.Len(t, aclsRecvd, 1)
	assert.True(t, proto.Equal(aclsRecvd[0], acls[idHashes[1]]))

	// ListIdentityHashes -- acls 1 and 2 remain
	idsRecvd, err = store.ListAllIdentity()
	assert.NoError(t, err)
	assert.Len(t, idsRecvd, 2)
	hashToIdRecvd := protos.GetHashToIdentity(idsRecvd)
	assert.True(t, proto.Equal(hashToIdRecvd[idHashes[1]], ids[1]))
	assert.True(t, proto.Equal(hashToIdRecvd[idHashes[2]], ids[2]))

	// ListIdentityHashes -- add back acl0, now acls 0, 1, and 2 remain
	err = store.PutACL(ids[0], acls[idHashes[0]])
	assert.NoError(t, err)
	idsRecvd, err = store.ListAllIdentity()
	assert.NoError(t, err)
	assert.Len(t, idsRecvd, 3)
	hashToIdRecvd = protos.GetHashToIdentity(idsRecvd)
	assert.True(t, proto.Equal(hashToIdRecvd[idHashes[0]], ids[0]))
	assert.True(t, proto.Equal(hashToIdRecvd[idHashes[1]], ids[1]))
	assert.True(t, proto.Equal(hashToIdRecvd[idHashes[2]], ids[2]))

	// UpdateACLWithEntities -- allow id0 to read id1
	ents := []*accessprotos.AccessControl_Entity{{Id: ids[1], Permissions: accessprotos.AccessControl_READ}}
	acl1, err := store.GetACL(ids[0])
	assert.NoError(t, err)
	_, id0HasPermsForID1 := acl1.Entities[idHashes[1]]
	assert.False(t, id0HasPermsForID1)

	err = store.UpdateACLWithEntities(ids[0], ents)
	assert.NoError(t, err)

	acl1, err = store.GetACL(ids[0])
	assert.NoError(t, err)
	id0PermForID1, ok := acl1.Entities[idHashes[1]]
	assert.True(t, ok)
	assert.Equal(t, id0PermForID1.Id.HashString(), idHashes[1])
	assert.Equal(t, id0PermForID1.Permissions, accessprotos.AccessControl_READ)

	// Nil arguments return err, don't cause panic
	_, err = store.GetACL(nil)
	assert.Error(t, err)

	_, err = store.GetManyACL(nil)
	assert.Error(t, err)
	_, err = store.GetManyACL([]*protos.Identity{nil})
	assert.Error(t, err)

	err = store.PutACL(nil, acls[idHashes[0]])
	assert.Error(t, err)
	err = store.PutACL(ids[0], nil)
	assert.Error(t, err)

	err = store.UpdateACLWithEntities(nil, ents)
	assert.Error(t, err)
	err = store.UpdateACLWithEntities(ids[0], nil)
	assert.Error(t, err)
	err = store.UpdateACLWithEntities(ids[0], []*accessprotos.AccessControl_Entity{nil})
	assert.Error(t, err)

	err = store.DeleteACL(nil)
	assert.Error(t, err)
}
