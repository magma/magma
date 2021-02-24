/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package blobstore_test

import (
	"testing"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
	magmaerrors "magma/orc8r/lib/go/errors"

	"github.com/stretchr/testify/require"
)

func TestMigration(t *testing.T) {
	db, err := sqorc.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	fact := blobstore.NewSQLBlobStorageFactory("states", db, sqorc.GetSqlBuilder())
	err = fact.InitializeFactory()
	require.NoError(t, err)
	storev1, err := fact.StartTransaction(nil)
	require.NoError(t, err)
	blobs := blobstore.Blobs{
		{Type: "type1", Key: "key1", Value: []byte("value")},
		{Type: "type2", Key: "key2", Value: []byte("value")},
		{Type: "type1", Key: "key2", Value: []byte("value")},
	}
	err = storev1.CreateOrUpdate("id1", blobs)
	require.NoError(t, err)

	many, err := storev1.GetMany("id1", []storage.TypeAndKey{
		{Type: "type1", Key: "key1"},
		{Type: "type2", Key: "key2"},
	})
	require.NoError(t, err)
	require.Len(t, many, 2)
	require.Equal(t, blobs[:2], many)

	keys, err := storev1.GetExistingKeys([]string{"key1"}, blobstore.SearchFilter{})
	require.NoError(t, err)
	require.Equal(t, []string{"key1"}, keys)

	err = storev1.Commit()
	require.NoError(t, err)

	entfact := blobstore.NewEntStorage("states", db, nil)
	storev2, err := entfact.StartTransaction(nil)
	require.NoError(t, err)
	blobs, err = storev2.GetMany("id1", []storage.TypeAndKey{
		{Type: "type1", Key: "key1"},
		{Type: "type2", Key: "key2"},
	})
	require.NoError(t, err)
	require.Len(t, many, 2)
	require.Equal(t, blobs[:2], many)

	blob, err := storev2.Get("id1", storage.TypeAndKey{Type: "type1", Key: "key1"})
	require.NoError(t, err)
	require.Equal(t, blobs[0], blob)

	blob, err = storev2.Get("id1", storage.TypeAndKey{Type: "type2", Key: "key2"})
	require.NoError(t, err)
	require.Equal(t, blobs[1], blob)

	err = storev2.IncrementVersion("id1", storage.TypeAndKey{Type: "type3", Key: "key1"})
	require.NoError(t, err)
	blob, err = storev2.Get("id1", storage.TypeAndKey{Type: "type3", Key: "key1"})
	require.NoError(t, err)
	require.Equal(t, blobstore.Blob{Type: "type3", Key: "key1", Version: 1}, blob)

	err = storev2.IncrementVersion("id1", storage.TypeAndKey{Type: "type3", Key: "key1"})
	require.NoError(t, err)
	blob, err = storev2.Get("id1", storage.TypeAndKey{Type: "type3", Key: "key1"})
	require.NoError(t, err)
	require.Equal(t, blobstore.Blob{Type: "type3", Key: "key1", Version: 2}, blob)

	err = storev2.Delete("id1", []storage.TypeAndKey{{Type: "type3", Key: "key1"}})
	require.NoError(t, err)
	_, err = storev2.Get("id1", storage.TypeAndKey{Type: "type3", Key: "key1"})
	require.Equal(t, magmaerrors.ErrNotFound, err)

	err = storev2.CreateOrUpdate("id1", blobstore.Blobs{
		{Type: "type1", Key: "key1", Value: []byte("world")},
		{Type: "type3", Key: "key1", Value: []byte("value")},
	})
	require.NoError(t, err)
	blob, err = storev2.Get("id1", storage.TypeAndKey{Type: "type3", Key: "key1"})
	require.NoError(t, err)
	require.Equal(t, blobstore.Blob{Type: "type3", Key: "key1", Value: []byte("value")}, blob)

	blob, err = storev2.Get("id1", storage.TypeAndKey{Type: "type1", Key: "key1"})
	require.NoError(t, err)
	require.Equal(t, blobstore.Blob{Type: "type1", Key: "key1", Value: []byte("world"), Version: 1}, blob)

	keys, err = storev2.GetExistingKeys([]string{"key1"}, blobstore.SearchFilter{})
	require.NoError(t, err)
	require.Equal(t, []string{"key1"}, keys)
}

func TestIntegration(t *testing.T) {
	db, err := sqorc.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	fact := blobstore.NewEntStorage("states", db, sqorc.GetSqlBuilder())
	integration(t, fact)
}
