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

import (
	"magma/orc8r/lib/go/protos"
)

type SyncStore interface {
	SyncStoreReader

	// CollectGarbage drops all store contents that aren't linked to
	// one of the passed networks.
	CollectGarbage(trackedNetworks []string)

	// SetDigest sets the digest tree for a network.
	SetDigest(network string, digests *protos.DigestTree) error

	// UpdateCache returns a CacheWriter object that'll be used to update
	// the cache of a network.
	UpdateCache(network string) (CacheWriter, error)
}

type SyncStoreReader interface {
	Initialize() error

	// GetDigests returns all digests last-updated before the passed unix time,
	// keyed by network. Caveats:
	// 1. If networks is empty, returns digests for all networks.
	// 2. lastUpdatedBefore is recorded in unix seconds. Filters for all digests that
	// were last updated earlier than this time.
	GetDigests(networks []string, lastUpdatedBefore int64, loadLeaves bool) (DigestTrees, error)

	// GetCachedByID and GetCachedByPage return cached objects by their IDs or
	// the page specified by the token. The returned objects are ordered by their
	// IDs in ascending order.
	GetCachedByID(network string, ids []string) ([][]byte, error)
	GetCachedByPage(network string, token string, pageSize uint64) ([][]byte, string, error)

	// GetLastResync returns the last resync time of a particular gateway.
	GetLastResync(network string, gateway string) (int64, error)

	// RecordResync tracks the last resync time of a gateway.
	RecordResync(network string, gateway string, t int64) error
}

type CacheWriter interface {
	// InsertMany adds objects to the list staged for the batch update.
	//
	// NOTE: Caller of the function should enforce that the max size of the
	// insertion aligns reasonably with the max page size of its corresponding
	// load source.
	InsertMany(objects map[string][]byte) error
	// Apply completes the batch cache update.
	Apply() error
}

type DigestTrees map[string]*protos.DigestTree

func (digestTrees DigestTrees) Networks() []string {
	var networks []string
	for network := range digestTrees {
		networks = append(networks, network)
	}
	return networks
}
