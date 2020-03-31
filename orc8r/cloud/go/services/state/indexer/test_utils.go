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

type testIndexer struct {
	id      string
	version Version
}

// NewTestIndexer returns a do-nothing indexer with the specified ID.
func NewTestIndexer(id string, version Version) Indexer {
	return &testIndexer{id: id, version: version}
}

func (t *testIndexer) GetID() string {
	return t.id
}

func (t *testIndexer) GetVersion() Version {
	return t.version
}

func (t *testIndexer) GetSubscriptions() []Subscription                           { return nil }
func (t *testIndexer) PrepareReindex(from, to Version, isFirstReindex bool) error { return nil }
func (t *testIndexer) CompleteReindex(from, to Version) error                     { return nil }
func (t *testIndexer) Index(networkID, reporterHWID string, states []state.State) (StateErrors, error) {
	return nil, nil
}
