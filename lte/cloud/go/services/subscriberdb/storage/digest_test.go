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

func TestDigestLookup(t *testing.T) {
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	s := storage.NewDigestLookup(db, sqorc.GetSqlBuilder())
	assert.NoError(t, s.Initialize())

	t.Run("empty initially", func(t *testing.T) {
		digest, timestamp, err := storage.GetDigest(s, "n0")
		assert.NoError(t, err)
		assert.Equal(t, digest, "")
		assert.Equal(t, timestamp, int64(0))

		networkIDs, err := storage.GetAllNetworks(s)
		assert.NoError(t, err)
		assert.Equal(t, []string{}, networkIDs)
	})

	t.Run("basic insert", func(t *testing.T) {
		err = s.SetDigest("n0", "apple")
		assert.NoError(t, err)
		err = s.SetDigest("n1", "lemon")
		assert.NoError(t, err)
		err = s.SetDigest("n2", "peach")
		assert.NoError(t, err)

		networkIDs, err := storage.GetAllNetworks(s)
		assert.NoError(t, err)
		assert.Equal(t, []string{"n0", "n1", "n2"}, networkIDs)

		digest, _, err := storage.GetDigest(s, "n0")
		assert.NoError(t, err)
		assert.Equal(t, "apple", digest)
		digest, _, err = storage.GetDigest(s, "n1")
		assert.NoError(t, err)
		assert.Equal(t, "lemon", digest)
		digest, _, err = storage.GetDigest(s, "n2")
		assert.NoError(t, err)
		assert.Equal(t, "peach", digest)
	})

	t.Run("upsert", func(t *testing.T) {
		err = s.SetDigest("n0", "banana")
		assert.NoError(t, err)
		digest, _, err := storage.GetDigest(s, "n0")
		assert.NoError(t, err)
		assert.Equal(t, "banana", digest)

		err = s.SetDigest("n0", "watermelon")
		assert.NoError(t, err)
		digest, _, err = storage.GetDigest(s, "n0")
		assert.NoError(t, err)
		assert.Equal(t, "watermelon", digest)
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

		digest, _, err := storage.GetDigest(s, "n1")
		assert.NoError(t, err)
		assert.Equal(t, "", digest)
		digest, _, err = storage.GetDigest(s, "n0")
		assert.NoError(t, err)
		assert.Equal(t, "watermelon", digest)
	})
}
