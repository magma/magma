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
	"time"

	"magma/lte/cloud/go/services/subscriberdb/storage"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/stretchr/testify/assert"
)

func TestDigestStore(t *testing.T) {
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	s := storage.NewDigestStore(db, sqorc.GetSqlBuilder())
	assert.NoError(t, s.Initialize())

	t.Run("return default value when empty", func(t *testing.T) {
		digest, err := storage.GetDigest(s, "n0")
		assert.NoError(t, err)
		assert.Equal(t, "", digest)

		networkIDs, err := storage.GetAllNetworks(s)
		assert.NoError(t, err)
		assert.Equal(t, []string{}, networkIDs)
	})

	t.Run("basic insert", func(t *testing.T) {
		err = s.SetDigest("n0", "apple")
		assert.NoError(t, err)
		err = s.SetDigest("n1", "banana")
		assert.NoError(t, err)
		err = s.SetDigest("n2", "cherry")
		assert.NoError(t, err)

		networkIDs, err := storage.GetAllNetworks(s)
		assert.NoError(t, err)
		assert.Equal(t, []string{"n0", "n1", "n2"}, networkIDs)

		digest, err := storage.GetDigest(s, "n0")
		assert.NoError(t, err)
		assert.Equal(t, "apple", digest)
		digest, err = storage.GetDigest(s, "n1")
		assert.NoError(t, err)
		assert.Equal(t, "banana", digest)
		digest, err = storage.GetDigest(s, "n2")
		assert.NoError(t, err)
		assert.Equal(t, "cherry", digest)
	})

	t.Run("upsert", func(t *testing.T) {
		err = s.SetDigest("n0", "apple2")
		assert.NoError(t, err)
		digest, err := storage.GetDigest(s, "n0")
		assert.NoError(t, err)
		assert.Equal(t, "apple2", digest)

		err = s.SetDigest("n0", "apple3")
		assert.NoError(t, err)
		digest, err = storage.GetDigest(s, "n0")
		assert.NoError(t, err)
		assert.Equal(t, "apple3", digest)
	})

	t.Run("get outdated", func(t *testing.T) {
		outdatedNetworks, err := storage.GetOutdatedNetworks(s, time.Now().Unix()+10000)
		assert.NoError(t, err)
		assert.Equal(t, []string{"n0", "n1", "n2"}, outdatedNetworks)

		outdatedNetworks, err = storage.GetOutdatedNetworks(s, time.Now().Unix()-10000)
		assert.NoError(t, err)
		assert.Equal(t, []string{}, outdatedNetworks)
	})

	t.Run("delete", func(t *testing.T) {
		err = s.DeleteDigests([]string{"n1", "n2"})
		assert.NoError(t, err)

		networks, err := storage.GetAllNetworks(s)
		assert.NoError(t, err)
		assert.Equal(t, []string{"n0"}, networks)

		digest, err := storage.GetDigest(s, "n1")
		assert.NoError(t, err)
		assert.Empty(t, digest)
		digest, err = storage.GetDigest(s, "n0")
		assert.NoError(t, err)
		assert.Equal(t, "apple3", digest)
	})
}
