/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package indexer

import (
	"magma/orc8r/cloud/go/services/state"
)

// Version of the indexer. Capped to uint32 to fit into Postgres/Maria integer (int32).
type Version uint32

// StateErrors is a mapping of state ID to error experienced indexing the state.
type StateErrors map[state.ID]error

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
	// Unique ID should be alphanumeric with underscores, prefixed by the owning service,
	// e.g. directoryd_sessionid.
	GetID() string

	// GetVersion returns the current version for the indexer.
	// Incrementing the version in a release will result in a reindex.
	// An indexer's version is required to be
	//	- nonzero
	// 	- non-decreasing across successive releases
	GetVersion() Version

	// GetSubscriptions defines the composite keys this indexer is interested in.
	GetSubscriptions() []Subscription

	// PrepareReindex prepares for a reindex operation.
	// isFirstReindex is set if this is the first time this indexer has been registered.
	PrepareReindex(from, to Version, isFirstReindex bool) error

	// CompleteReindex indicates the reindex operation is complete.
	CompleteReindex(from, to Version) error

	// Index updates secondary indices based on the added/updated states.
	Index(networkID string, states state.StatesByID) (StateErrors, error)

	// TODO(4/10/20): consider adding support for removing states from an indexer
	// IndexRemove updates secondary indices based on the removed states.
	//IndexRemove(states state.StatesByID) (StateErrors, error)
}
