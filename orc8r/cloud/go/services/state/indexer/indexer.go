/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package indexer

import (
	"fmt"

	state_types "magma/orc8r/cloud/go/services/state/types"

	"github.com/thoas/go-funk"
)

// Version of the indexer. Capped to uint32 to fit into Postgres/Maria integer (int32).
type Version uint32

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
	Index(networkID string, states state_types.StatesByID) (state_types.StateErrors, error)

	// TODO(4/10/20): consider adding support for removing states from an indexer
	// IndexRemove updates secondary indices based on the removed states.
	//IndexRemove(states state.StatesByID) (StateErrors, error)
}

// FilterIDs to the subset that match one of the state types.
func FilterIDs(types []string, ids []state_types.ID) []state_types.ID {
	var ret []state_types.ID
	for _, id := range ids {
		if funk.Contains(types, id.Type) {
			ret = append(ret, id)
		}
	}
	return ret
}

// FilterStates to the subset that match one of the state types.
func FilterStates(types []string, states state_types.StatesByID) state_types.StatesByID {
	ret := state_types.StatesByID{}
	for id, st := range states {
		if funk.Contains(types, id.Type) {
			ret[id] = st
		}
	}
	return ret
}

// Versions represents the discrepancy between an indexer's versions, actual vs. desired.
type Versions struct {
	IndexerID string
	Actual    Version
	Desired   Version
}

func (v *Versions) String() string {
	return fmt.Sprintf("{id: %s, actual: %d, desired: %d}", v.IndexerID, v.Actual, v.Desired)
}
