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
	"magma/orc8r/cloud/go/test_utils"

	"github.com/stretchr/testify/assert"
)

func TestPerSubDigestStore(t *testing.T) {
	fact := test_utils.NewSQLBlobstore(t, subscriberdb.PerSubDigestTableBlobstore)
	assert.NoError(t, fact.InitializeFactory())
	s := storage.NewPerSubDigestStore(fact)

	t.Run("empty initially", func(t *testing.T) {
		digest, err := s.GetDigest("n0")
		assert.NoError(t, err)
		checkPerSubDigests(t, []*lte_protos.SubscriberDigestWithID{}, digest)
	})

	t.Run("basic insert", func(t *testing.T) {
		expected := []*lte_protos.SubscriberDigestWithID{
			{
				Sid:    &lte_protos.SubscriberID{Id: "00001", Type: lte_protos.SubscriberID_IMSI},
				Digest: &lte_protos.Digest{Md5Base64Digest: "apple"},
			},
			{
				Sid:    &lte_protos.SubscriberID{Id: "00002", Type: lte_protos.SubscriberID_IMSI},
				Digest: &lte_protos.Digest{Md5Base64Digest: "banana"},
			},
			{
				Sid:    &lte_protos.SubscriberID{Id: "00003", Type: lte_protos.SubscriberID_IMSI},
				Digest: &lte_protos.Digest{Md5Base64Digest: "cherry"},
			},
		}
		err := s.SetDigest("n0", expected)
		assert.NoError(t, err)

		got, err := s.GetDigest("n0")
		assert.NoError(t, err)
		checkPerSubDigests(t, expected, got)
	})

	t.Run("upsert", func(t *testing.T) {
		expected := []*lte_protos.SubscriberDigestWithID{
			{
				Sid:    &lte_protos.SubscriberID{Id: "00001", Type: lte_protos.SubscriberID_IMSI},
				Digest: &lte_protos.Digest{Md5Base64Digest: "apple2"},
			},
			{
				Sid:    &lte_protos.SubscriberID{Id: "00003", Type: lte_protos.SubscriberID_IMSI},
				Digest: &lte_protos.Digest{Md5Base64Digest: "cherry2"},
			},
			{
				Sid:    &lte_protos.SubscriberID{Id: "00004", Type: lte_protos.SubscriberID_IMSI},
				Digest: &lte_protos.Digest{Md5Base64Digest: "dragonfruit"},
			},
		}
		// The upserted set should completely replace the original set
		err := s.SetDigest("n1", expected)
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
		err := s.DeleteDigests([]string{"n0", "n1"})
		assert.NoError(t, err)

		got, err := s.GetDigest("n0")
		assert.NoError(t, err)
		checkPerSubDigests(t, []*lte_protos.SubscriberDigestWithID{}, got)

		got, err = s.GetDigest("n1")
		assert.NoError(t, err)
		checkPerSubDigests(t, []*lte_protos.SubscriberDigestWithID{}, got)

		// Deleting digests of a non-existent network shouldn't cause an error
		err = s.DeleteDigests([]string{"n2"})
		assert.NoError(t, err)
	})
}

func checkPerSubDigests(t *testing.T, expected []*lte_protos.SubscriberDigestWithID, got []*lte_protos.SubscriberDigestWithID) {
	assert.Equal(t, len(expected), len(got))
	for ind := range expected {
		assert.Equal(t, expected[ind].Digest.GetMd5Base64Digest(), got[ind].Digest.GetMd5Base64Digest())
		assert.Equal(t, expected[ind].Sid.Id, got[ind].Sid.Id)
	}
}
