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

package syncstore

import "magma/orc8r/lib/go/protos"

// GetLeafDigestsDiff computes the data changeset according to two lists of leaf digests,
// ordered by their IDs. It returns
// 1. The set of objects that have been added/modified, with the new digests.
// 2. An ordered list of the objects that have been removed.
func GetLeafDigestsDiff(prev []*protos.LeafDigest, next []*protos.LeafDigest) (map[string]string, []string) {
	n, m, i, j := len(prev), len(next), 0, 0
	toRenew := map[string]string{}
	var deleted []string

	for i < n && j < m {
		iId, jId := prev[i].Id, next[j].Id
		if iId == jId {
			if prev[i].Digest.Md5Base64Digest != next[j].Digest.Md5Base64Digest {
				toRenew[jId] = next[j].Digest.Md5Base64Digest
			}
			i++
			j++
		} else if iId > jId {
			toRenew[jId] = next[j].Digest.Md5Base64Digest
			j++
		} else {
			deleted = append(deleted, iId)
			i++
		}
	}

	for ; i < n; i++ {
		prevId := prev[i].Id
		deleted = append(deleted, prevId)
	}
	for ; j < m; j++ {
		nextId := next[j].Id
		toRenew[nextId] = next[j].Digest.Md5Base64Digest
	}

	return toRenew, deleted
}
