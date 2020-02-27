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

// Indexer creates a set of secondary indices for consumption by a service.
// Each Indexer should
// 	- be owned by a single service
//	- have a unique ID, alphanumeric with underscores, prefixed by the owning service,
//	  e.g. directoryd_session_id
//	- store per-version data in a properly isolated manner,
// 	  e.g. different SQL tables for different indexer versions
//	- have its generated data exposed by the owning service,
//	  i.e. only one other service should access the generated data directly via the storage interface.
type Indexer interface {
	// GetID returns the unique identifier for the indexer.
	GetID() string

	// GetVersion returns the current version for the indexer.
	// Incrementing the version in a release will result in a reindex.
	GetVersion() uint64

	// GetSubscriptions defines the composite keys this indexer is interested in.
	GetSubscriptions() []Subscription

	// PrepareReindex prepares for a reindex operation.
	// Each version should use e.g. a separate SQL table, so preparing for
	// a reindex would include creating new table(s).
	// isFirstReindex is set if this is the first time this indexer has been registered.
	PrepareReindex(from, to uint64, isFirstReindex bool)

	// CompleteReindex indicates the reindex operation is complete.
	// Any internal state relevant only to the from version can subsequently be
	// safely removed, e.g. dropping old SQL tables.
	CompleteReindex(from, to uint64)

	// Index updates secondary indices based on the added/updated states.
	Index(reporterHWID string, states []state.State) (map[storage.TypeAndKey]error, error)

	// IndexRemove updates secondary indices based on the removed states.
	// NOTE: for now, we will defer IndexRemove to future efforts.
	//IndexRemove(reporterHWID string, states []State) (map[TypeAndKey]error, error)
}

// Subscription denotes a set of primary keys.
type Subscription struct {
	Type    string
	Matcher KeyMatcher
}
