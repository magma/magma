/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package indexer_test

import (
	"sort"
	"testing"

	assert "github.com/stretchr/testify/require"

	"magma/orc8r/cloud/go/services/state/indexer"
)

func TestRegisterRemote(t *testing.T) {
	want0 := indexer.NewRemoteIndexer("some_service_0", 42, "type0", "type1")
	want1 := indexer.NewRemoteIndexer("some_service_1", 420)
	want2 := indexer.NewRemoteIndexer("some_service_2", 424, "type2")

	indexer.DeregisterAllForTest(t)

	t.Run("empty initially", func(t *testing.T) {
		got := indexer.GetIndexers()
		assert.Empty(t, got)
	})

	t.Run("set and get one", func(t *testing.T) {
		err := indexer.RegisterIndexers(want0)
		assert.NoError(t, err)

		got := indexer.GetIndexers()
		assert.Equal(t, []indexer.Indexer{want0}, got)
		gotOne := indexer.GetIndexer(want0.GetID())
		assert.Equal(t, want0, gotOne)
	})

	t.Run("set and get two more", func(t *testing.T) {
		err := indexer.RegisterIndexers(want1, want2)
		assert.NoError(t, err)

		got := indexer.GetIndexers()
		sort.Slice(got, func(i, j int) bool { return got[i].GetID() < got[j].GetID() })
		assert.Equal(t, []indexer.Indexer{want0, want1, want2}, got)
		got1 := indexer.GetIndexer(want1.GetID())
		assert.Equal(t, want1, got1)
		got2 := indexer.GetIndexer(want2.GetID())
		assert.Equal(t, want2, got2)
	})

	t.Run("fail overwrite same name", func(t *testing.T) {
		err := indexer.RegisterIndexers(want2)
		assert.Error(t, err)

		got := indexer.GetIndexers()
		sort.Slice(got, func(i, j int) bool { return got[i].GetID() < got[j].GetID() })
		assert.Equal(t, []indexer.Indexer{want0, want1, want2}, got)
	})

	t.Run("get indexers for state type", func(t *testing.T) {
		got := indexer.GetIndexersForState("type2")
		assert.Equal(t, []indexer.Indexer{want2}, got)
	})
}
