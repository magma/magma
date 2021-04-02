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

package indexer

import (
	"fmt"

	state_types "magma/orc8r/cloud/go/services/state/types"
)

// Indexer creates a set of secondary indices for consumption by a service.
// Each Indexer should
// 	- be owned by a single service
//	- have its generated data exposed by the owning service, i.e. only one other service
//	  should access the generated data directly via the storage interface.
// Notes
//	- There is an unlikely but existent race condition during a reindex
//	  operation, where Index could be called with an outdated version of a state.
//	  If indexers care about preventing this race condition:
//		- add a Reindex method to indexer interface, called in-place of Index
//		  during reindex operations
//		- individual indexers should track received state IDs per version and
//		  drop Reindex-ed states with stale versions.
type Indexer interface {
	// GetID returns the unique identifier for the indexer.
	// For remote indexers, unique ID should be the service name.
	GetID() string

	// GetVersion returns the current version for the indexer.
	// Incrementing the version in a release will result in a reindex.
	// An indexer's version is required to be
	//	- nonzero
	// 	- non-decreasing across successive releases
	GetVersion() Version

	// GetTypes defines the types of states this indexer is interested in.
	GetTypes() []string

	// PrepareReindex prepares for a reindex operation.
	// isFirstReindex is set if this is the first time this indexer has been registered.
	PrepareReindex(from, to Version, isFirstReindex bool) error

	// CompleteReindex indicates the reindex operation is complete.
	CompleteReindex(from, to Version) error

	// Index updates secondary indices based on the added/updated states.
	Index(networkID string, states state_types.SerializedStatesByID) (state_types.StateErrors, error)

	// TODO(4/10/20): consider adding support for removing states from an indexer
	// IndexRemove updates secondary indices based on the removed states.
	//IndexRemove(states state_types.SerializedStatesByID) (state_types.StateErrors, error)
}

// Version of the indexer. Capped to uint32 to fit into Postgres/Maria integer (int32).
type Version uint32

// NewIndexerVersion returns a new indexer version, ensuring it fits into
// the required integer size.
func NewIndexerVersion(version int64) (Version, error) {
	capped := Version(version)
	if int64(capped) != version {
		return 0, fmt.Errorf("indexer version %v too large for %T", version, Version(0))
	}
	return capped, nil
}

// Versions represents the discrepancy between an indexer's versions,
// actual vs. desired.
type Versions struct {
	// IndexerID is the ID of the indexer.
	// ID should be the owning service's name.
	IndexerID string
	// Actual version of the indexer.
	Actual Version
	// Desired version of the indexer.
	Desired Version
}

func (v *Versions) String() string {
	return fmt.Sprintf("{id: %s, actual: %d, desired: %d}", v.IndexerID, v.Actual, v.Desired)
}
