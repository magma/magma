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
	"magma/orc8r/cloud/go/services/certifier/protos"
	"magma/orc8r/cloud/go/services/certifier/storage"
	"magma/orc8r/cloud/go/sqorc"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestCertifierStorageBlobstore_Integation(t *testing.T) {
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	fact := blobstore.NewEntStorage(storage.CertifierTableBlobstore, db, sqorc.GetSqlBuilder())
	err = fact.InitializeFactory()
	assert.NoError(t, err)
	store := storage.NewCertifierBlobstore(fact)
	testCertifierStorageImpl(t, store)
}

func testCertifierStorageImpl(t *testing.T, store storage.CertifierStorage) {
	sn0 := "serial_number_0"
	sn1 := "serial_number_1"
	sn2 := "serial_number_2"
	info0 := &protos.CertificateInfo{
		Id:        nil,
		NotBefore: &timestamp.Timestamp{Seconds: 0xdead, Nanos: 0xbeef},
		NotAfter:  &timestamp.Timestamp{Seconds: 0xbeef, Nanos: 0x4444},
		CertType:  1,
	}
	info1 := &protos.CertificateInfo{
		Id:        nil,
		NotBefore: &timestamp.Timestamp{Seconds: 0x1111, Nanos: 0x2222},
		NotAfter:  &timestamp.Timestamp{Seconds: 0x3333, Nanos: 0xdddd},
		CertType:  2,
	}
	info2 := &protos.CertificateInfo{
		Id:        nil,
		NotBefore: &timestamp.Timestamp{Seconds: 0x9999, Nanos: 0x2222},
		NotAfter:  &timestamp.Timestamp{Seconds: 0x3333, Nanos: 0xaaaa},
		CertType:  3,
	}

	// Empty initially
	sns, err := store.ListSerialNumbers()
	assert.NoError(t, err)
	assert.Len(t, sns, 0)

	// Put and Get info0
	err = store.PutCertInfo(sn0, info0)
	assert.NoError(t, err)

	info, err := store.GetCertInfo(sn0)
	assert.NoError(t, err)
	assert.True(t, proto.Equal(info, info0))

	// Put and Get info1
	err = store.PutCertInfo(sn1, info1)
	assert.NoError(t, err)

	info, err = store.GetCertInfo(sn1)
	assert.NoError(t, err)
	assert.True(t, proto.Equal(info, info1))

	// Put info2, GetMany infos 0 and 1
	err = store.PutCertInfo(sn2, info2)
	assert.NoError(t, err)
	infos, err := store.GetManyCertInfo([]string{sn0, sn1})
	assert.NoError(t, err)
	assert.Len(t, infos, 2)
	assert.True(t, proto.Equal(infos[sn0], info0))
	assert.True(t, proto.Equal(infos[sn1], info1))

	// Delete info0, Get info0, GetMany infos 0 and 1
	err = store.DeleteCertInfo(sn0)
	assert.NoError(t, err)
	_, err = store.GetCertInfo(sn0)
	assert.EqualError(t, err, merrors.ErrNotFound.Error())
	infos, err = store.GetManyCertInfo([]string{sn0, sn1})
	assert.NoError(t, err)
	assert.Len(t, infos, 1)
	assert.True(t, proto.Equal(infos[sn1], info1))

	// ListSerialNumbers -- sns 1 and 2 remain
	sns, err = store.ListSerialNumbers()
	assert.NoError(t, err)
	assert.Len(t, sns, 2)
	assert.Contains(t, sns, sn1)
	assert.Contains(t, sns, sn2)

	// GetAll -- infos 1 and 2 remain
	infos, err = store.GetAllCertInfo()
	assert.NoError(t, err)
	assert.Len(t, infos, 2)
	assert.True(t, proto.Equal(infos[sn1], info1))
	assert.True(t, proto.Equal(infos[sn2], info2))

	// GetAll -- add back info0, now infos 0, 1, and 2 remain
	err = store.PutCertInfo(sn0, info0)
	assert.NoError(t, err)
	infos, err = store.GetAllCertInfo()
	assert.NoError(t, err)
	assert.Len(t, infos, 3)
	assert.True(t, proto.Equal(infos[sn0], info0))
	assert.True(t, proto.Equal(infos[sn1], info1))
	assert.True(t, proto.Equal(infos[sn2], info2))
}
