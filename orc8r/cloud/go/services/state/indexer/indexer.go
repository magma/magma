/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package indexer

import (
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/storage"
)

// Version of the indexer. Capped to uint32 to fit into Postgres/Maria integer (int32).
type Version uint32

// StateErrors is a mapping of state type+key to error experienced indexing the state.
// Type is state type, key is state reporter ID.
type StateErrors map[storage.TypeAndKey]error

// Indexer creates a set of secondary indices for consumption by a service.
// Each Indexer should
// 	- be owned by a single service
//	- store per-version data in a properly isolated manner,
// 	  e.g. different SQL tables for different indexer versions
//	- have its generated data exposed by the owning service,
//	  i.e. only one other service should access the generated data directly via the storage interface.
type Indexer interface {
	// GetID returns the unique identifier for the indexer.
	// Unique ID should be alphanumeric with underscores, prefixed by the owning service,
	// e.g. directoryd_sessionid.
	GetID() string

	// GetVersion returns the current version for the indexer.
	// Incrementing the version in a release will result in a reindex.
	// An indexer's version is required to be non-decreasing across successive releases.
	GetVersion() Version

	// GetSubscriptions defines the composite keys this indexer is interested in.
	GetSubscriptions() []Subscription

	// PrepareReindex prepares for a reindex operation.
	// Each version should use e.g. a separate SQL table, so preparing for
	// a reindex would include creating new table(s).
	// isFirstReindex is set if this is the first time this indexer has been registered.
	PrepareReindex(from, to Version, isFirstReindex bool) error

	// CompleteReindex indicates the reindex operation is complete.
	// Any internal state relevant only to the from version can subsequently be
	// safely removed, e.g. dropping old SQL tables.
	CompleteReindex(from, to Version) error

	// Index updates secondary indices based on the added/updated states.
	Index(networkID, reporterHWID string, states []state.State) (StateErrors, error)

	// IndexRemove updates secondary indices based on the removed states.
	// NOTE: for now, we will defer IndexRemove to future efforts.
	//IndexRemove(reporterHWID string, states []State) (map[TypeAndKey]error, error)
}
