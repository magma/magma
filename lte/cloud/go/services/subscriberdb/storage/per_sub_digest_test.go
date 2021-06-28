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

	"magma/lte/cloud/go/services/subscriberdb/storage"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/stretchr/testify/assert"
)

func TestPerSubDigestLookup(t *testing.T) {
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	s := storage.NewPerSubDigestLookup(db, sqorc.GetSqlBuilder())
	assert.NoError(t, s.Initialize())

	t.Run("empty initially", func(t *testing.T) {
		digest, err := s.GetDigest("n0")
		assert.NoError(t, err)
		checkPerSubDigests(t, map[string]string{}, digest)

		networkIDs, err := storage.GetAllNetworks(s)
		assert.NoError(t, err)
		assert.Equal(t, []string{}, networkIDs)
	})

	t.Run("basic insert", func(t *testing.T) {
		err = s.SetDigest("n0", storage.PerSubDigestUpsertArgs{
			ToRenew: map[string]string{
				"IMSI0001": "apple",
				"IMSI0002": "lemon",
				"IMSI0003": "peach",
			},
		})
		assert.NoError(t, err)
		err = s.SetDigest("n1", storage.PerSubDigestUpsertArgs{
			ToRenew: map[string]string{
				"IMSI1111": "banana",
				"IMSI1112": "durian",
			},
		})
		assert.NoError(t, err)
		err = s.SetDigest("n2", storage.PerSubDigestUpsertArgs{
			ToRenew: map[string]string{"IMSI2221": "cherry"},
		})
		assert.NoError(t, err)

		networkIDs, err := storage.GetAllNetworks(s)
		assert.NoError(t, err)
		assert.Equal(t, []string{"n0", "n1", "n2"}, networkIDs)

		digest, err := s.GetDigest("n0")
		assert.NoError(t, err)
		expected := map[string]string{
			"IMSI0001": "apple",
			"IMSI0002": "lemon",
			"IMSI0003": "peach",
		}
		checkPerSubDigests(t, expected, digest)

		digest, err = s.GetDigest("n1")
		assert.NoError(t, err)
		expected = map[string]string{
			"IMSI1111": "banana",
			"IMSI1112": "durian",
		}
		checkPerSubDigests(t, expected, digest)

		digest, err = s.GetDigest("n2")
		assert.NoError(t, err)
		expected = map[string]string{"IMSI2221": "cherry"}
		checkPerSubDigests(t, expected, digest)
	})

	t.Run("upsert", func(t *testing.T) {
		err = s.SetDigest("n0", storage.PerSubDigestUpsertArgs{
			ToRenew: map[string]string{"IMSI0001": "orange", "IMSI0004": "papaya"},
			Deleted: []string{"IMSI0002"},
		})
		assert.NoError(t, err)
		digest, err := s.GetDigest("n0")
		assert.NoError(t, err)
		expected := map[string]string{
			"IMSI0001": "orange",
			"IMSI0003": "peach",
			"IMSI0004": "papaya",
		}
		checkPerSubDigests(t, expected, digest)

		err = s.SetDigest("n1", storage.PerSubDigestUpsertArgs{
			ToRenew: map[string]string{"IMSI1113": "starfruit", "IMSI1114": "cactus"},
			Deleted: []string{"IMSI1111", "IMSI1112"},
		})
		assert.NoError(t, err)
		digest, err = s.GetDigest("n1")
		assert.NoError(t, err)
		expected = map[string]string{
			"IMSI1113": "starfruit",
			"IMSI1114": "cactus",
		}
		checkPerSubDigests(t, expected, digest)
	})

	t.Run("delete", func(t *testing.T) {
		err = s.DeleteDigests([]string{"n1", "n2"})
		assert.NoError(t, err)

		networks, err := storage.GetAllNetworks(s)
		assert.NoError(t, err)
		assert.Equal(t, []string{"n0"}, networks)

		digest, err := s.GetDigest("n1")
		assert.NoError(t, err)
		checkPerSubDigests(t, map[string]string{}, digest)
		digest, err = s.GetDigest("n0")
		assert.NoError(t, err)
		expected := map[string]string{
			"IMSI0001": "orange",
			"IMSI0003": "peach",
			"IMSI0004": "papaya",
		}
		checkPerSubDigests(t, expected, digest)
	})
}

func checkPerSubDigests(t *testing.T, expected map[string]string, digest interface{}) {
	digestsBySubscriber, ok := digest.(map[string]string)
	assert.True(t, ok)
	assert.Equal(t, len(expected), len(digestsBySubscriber))
	for k := range expected {
		assert.Equal(t, expected[k], digestsBySubscriber[k])
	}
}
