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

package syncstore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/syncstore"
	"magma/orc8r/lib/go/protos"
)

func TestGetLeafDigestsDiff(t *testing.T) {
	t.Run("both empty", func(t *testing.T) {
		var prev, next []*protos.LeafDigest
		toRenew, deleted := syncstore.GetLeafDigestsDiff(prev, next)
		assert.Empty(t, toRenew)
		assert.Empty(t, deleted)
	})
	t.Run("all empty", func(t *testing.T) {
		prev := []*protos.LeafDigest{
			{Id: "IMSI00000", Digest: &protos.Digest{Md5Base64Digest: "apple"}},
			{Id: "IMSI00001", Digest: &protos.Digest{Md5Base64Digest: "banana"}},
			{Id: "IMSI00002", Digest: &protos.Digest{Md5Base64Digest: "cherry"}},
		}

		var next []*protos.LeafDigest
		toRenew, deleted := syncstore.GetLeafDigestsDiff(prev, next)

		assert.Empty(t, toRenew)
		assert.Equal(t, []string{"IMSI00000", "IMSI00001", "IMSI00002"}, deleted)
	})
	t.Run("tracked empty", func(t *testing.T) {
		var prev []*protos.LeafDigest
		next := []*protos.LeafDigest{
			{Id: "IMSI00000", Digest: &protos.Digest{Md5Base64Digest: "apple"}},
			{Id: "IMSI00001", Digest: &protos.Digest{Md5Base64Digest: "banana"}},
			{Id: "IMSI00002", Digest: &protos.Digest{Md5Base64Digest: "cherry"}},
		}
		expectedToRenew := map[string]string{
			"IMSI00000": "apple",
			"IMSI00001": "banana",
			"IMSI00002": "cherry",
		}
		toRenew, deleted := syncstore.GetLeafDigestsDiff(prev, next)
		assert.Equal(t, expectedToRenew, toRenew)
		assert.Empty(t, deleted)
	})
	t.Run("both not empty, basic", func(t *testing.T) {
		prev := []*protos.LeafDigest{
			{Id: "IMSI00001", Digest: &protos.Digest{Md5Base64Digest: "apple"}},
			{Id: "IMSI00002", Digest: &protos.Digest{Md5Base64Digest: "banana"}},
			{Id: "IMSI00004", Digest: &protos.Digest{Md5Base64Digest: "dragonfruit"}},
		}
		next := []*protos.LeafDigest{
			{Id: "IMSI00001", Digest: &protos.Digest{Md5Base64Digest: "apple"}},
			{Id: "IMSI00002", Digest: &protos.Digest{Md5Base64Digest: "banana2"}},
			{Id: "IMSI00003", Digest: &protos.Digest{Md5Base64Digest: "cherry"}},
		}
		expectedToRenew := map[string]string{
			"IMSI00002": "banana2",
			"IMSI00003": "cherry",
		}
		toRenew, deleted := syncstore.GetLeafDigestsDiff(prev, next)
		assert.Equal(t, expectedToRenew, toRenew)
		assert.Equal(t, []string{"IMSI00004"}, deleted)
	})
	t.Run("both not empty, involved", func(t *testing.T) {
		prev := []*protos.LeafDigest{
			{Id: "IMSI00001", Digest: &protos.Digest{Md5Base64Digest: "apple"}},
			{Id: "IMSI00002", Digest: &protos.Digest{Md5Base64Digest: "banana"}},
			{Id: "IMSI00004", Digest: &protos.Digest{Md5Base64Digest: "dragonfruit"}},
			{Id: "IMSI00006", Digest: &protos.Digest{Md5Base64Digest: "fig"}},
			{Id: "IMSI00007", Digest: &protos.Digest{Md5Base64Digest: "grape"}},
		}
		next := []*protos.LeafDigest{
			{Id: "IMSI00001", Digest: &protos.Digest{Md5Base64Digest: "apple"}},
			{Id: "IMSI00002", Digest: &protos.Digest{Md5Base64Digest: "banana"}},
			{Id: "IMSI00003", Digest: &protos.Digest{Md5Base64Digest: "cherry"}},
			{Id: "IMSI00005", Digest: &protos.Digest{Md5Base64Digest: "eggplant"}},
			{Id: "IMSI00006", Digest: &protos.Digest{Md5Base64Digest: "fig2"}},
			{Id: "IMSI00008", Digest: &protos.Digest{Md5Base64Digest: "honeydew"}},
		}
		expectedToRenew := map[string]string{
			"IMSI00003": "cherry",
			"IMSI00005": "eggplant",
			"IMSI00006": "fig2",
			"IMSI00008": "honeydew",
		}
		toRenew, deleted := syncstore.GetLeafDigestsDiff(prev, next)
		assert.Equal(t, expectedToRenew, toRenew)
		assert.Equal(t, []string{"IMSI00004", "IMSI00007"}, deleted)
	})
}
