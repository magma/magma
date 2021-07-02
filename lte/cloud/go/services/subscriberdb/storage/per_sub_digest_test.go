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

	lte_protos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/subscriberdb"
	"magma/lte/cloud/go/services/subscriberdb/storage"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/stretchr/testify/assert"
)

func TestPerSubDigestLookup(t *testing.T) {
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	fact := blobstore.NewEntStorage(subscriberdb.PerSubDigestTableBlobstore, db, sqorc.GetSqlBuilder())
	assert.NoError(t, fact.InitializeFactory())
	s := storage.NewPerSubDigestLookup(fact)

	t.Run("empty initially", func(t *testing.T) {
		digest, err := s.GetDigest("n0")
		assert.NoError(t, err)
		checkPerSubDigests(t, []*lte_protos.SubscriberDigestByID{}, digest)
	})

	t.Run("basic insert", func(t *testing.T) {
		expected := []*lte_protos.SubscriberDigestByID{
			{
				Sid:    &lte_protos.SubscriberID{Id: "00000", Type: lte_protos.SubscriberID_IMSI},
				Digest: &lte_protos.Digest{Md5Base64Digest: "apple"},
			},
			{
				Sid:    &lte_protos.SubscriberID{Id: "00001", Type: lte_protos.SubscriberID_IMSI},
				Digest: &lte_protos.Digest{Md5Base64Digest: "lemon"},
			},
			{
				Sid:    &lte_protos.SubscriberID{Id: "00001", Type: lte_protos.SubscriberID_IMSI},
				Digest: &lte_protos.Digest{Md5Base64Digest: "peach"},
			},
		}
		err = s.SetDigest("n0", expected)
		assert.NoError(t, err)

		got, err := s.GetDigest("n0")
		assert.NoError(t, err)
		checkPerSubDigests(t, expected, got)
	})

	t.Run("upsert", func(t *testing.T) {
		expected := []*lte_protos.SubscriberDigestByID{
			{
				Sid:    &lte_protos.SubscriberID{Id: "00001", Type: lte_protos.SubscriberID_IMSI},
				Digest: &lte_protos.Digest{Md5Base64Digest: "turtle"},
			},
			{
				Sid:    &lte_protos.SubscriberID{Id: "00003", Type: lte_protos.SubscriberID_IMSI},
				Digest: &lte_protos.Digest{Md5Base64Digest: "donkey"},
			},
			{
				Sid:    &lte_protos.SubscriberID{Id: "00004", Type: lte_protos.SubscriberID_IMSI},
				Digest: &lte_protos.Digest{Md5Base64Digest: "monkey"},
			},
		}
		// The upserted set should completely replace the original set
		err = s.SetDigest("n1", expected)
		assert.NoError(t, err)
		got, err := s.GetDigest("n1")
		assert.NoError(t, err)
		checkPerSubDigests(t, expected, got)

		err = s.SetDigest("n0", expected)
		assert.NoError(t, err)
		got, err = s.GetDigest("n0")
		assert.NoError(t, err)
		checkPerSubDigests(t, expected, got)
	})

	t.Run("delete many", func(t *testing.T) {
		err = s.DeleteDigests([]string{"n0", "n1"})
		assert.NoError(t, err)

		got, err := s.GetDigest("n0")
		assert.NoError(t, err)
		checkPerSubDigests(t, []*lte_protos.SubscriberDigestByID{}, got)

		got, err = s.GetDigest("n1")
		assert.NoError(t, err)
		checkPerSubDigests(t, []*lte_protos.SubscriberDigestByID{}, got)

		// Deleting digests of a non-existent network shouldn't cause an error
		err = s.DeleteDigests([]string{"n2"})
		assert.NoError(t, err)
	})
}

func checkPerSubDigests(t *testing.T, expected []*lte_protos.SubscriberDigestByID, got []*lte_protos.SubscriberDigestByID) {
	assert.Equal(t, len(expected), len(got))
	for ind := range expected {
		assert.Equal(t, expected[ind].Digest.GetMd5Base64Digest(), got[ind].Digest.GetMd5Base64Digest())
		assert.Equal(t, expected[ind].Sid.Id, got[ind].Sid.Id)
	}
}
