/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package blobstore_test

import (
	"sort"
	"testing"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/storage"
	magmaerrors "magma/orc8r/lib/go/errors"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func integration(t *testing.T, fact blobstore.BlobStorageFactory) {
	// Check the contract for an empty data store
	err := fact.InitializeFactory()
	store, err := fact.StartTransaction(nil)
	assert.NoError(t, err)
	listActual, err := store.ListKeys("network", "type")
	assert.NoError(t, err)
	assert.Empty(t, listActual)

	getActual, err := store.Get("network", storage.TypeAndKey{Type: "t", Key: "k"})
	assert.True(t, err == magmaerrors.ErrNotFound)
	assert.Equal(t, blobstore.Blob{}, getActual)

	getManyActual, err := store.GetMany(
		"network",
		[]storage.TypeAndKey{
			{Type: "t1", Key: "k1"},
			{Type: "t2", Key: "k2"},
		},
	)
	assert.NoError(t, err)
	assert.Empty(t, getManyActual)
	assert.NoError(t, store.Commit())

	// Workflow test
	store1, err := fact.StartTransaction(nil)
	assert.NoError(t, err)

	// Create blobs on 2 networks
	// network1: (t1, t2) X (k1, k2)
	err = store1.CreateOrUpdate("network1", []blobstore.Blob{
		{Type: "t1", Key: "k1", Value: []byte("v1")},
		{Type: "t1", Key: "k2", Value: []byte("v2")},
		{Type: "t2", Key: "k1", Value: []byte("v3"), Version: 2},
		{Type: "t2", Key: "k2", Value: []byte("v4"), Version: 1},
	})
	assert.NoError(t, err)
	assert.NoError(t, store1.Commit())

	// network2: (t3) X (k3, k4)
	store2, err := fact.StartTransaction(nil)
	assert.NoError(t, err)
	err = store2.CreateOrUpdate("network2", []blobstore.Blob{
		{Type: "t3", Key: "k3", Value: []byte("v5")},
		{Type: "t3", Key: "k4", Value: []byte("v6")},
	})
	assert.NoError(t, err)
	assert.NoError(t, store2.Commit())

	// Read tests
	store, err = fact.StartTransaction(nil)
	assert.NoError(t, err)

	byNetworkActual, err := blobstore.ListKeysByNetwork(store)
	assert.NoError(t, err)
	for _, v := range byNetworkActual {
		sort.Slice(v, getTKsComparator(v))
	}
	byNetworkExpected := map[string][]storage.TypeAndKey{
		"network1": {
			{Type: "t1", Key: "k1"},
			{Type: "t1", Key: "k2"},
			{Type: "t2", Key: "k1"},
			{Type: "t2", Key: "k2"},
		},
		"network2": {
			{Type: "t3", Key: "k3"},
			{Type: "t3", Key: "k4"},
		},
	}
	assert.Equal(t, byNetworkExpected, byNetworkActual)

	listActual, err = store.ListKeys("network1", "t1")
	assert.NoError(t, err)
	assert.Equal(t, []string{"k1", "k2"}, listActual)

	getManyActual, err = store.GetMany("network1", []storage.TypeAndKey{
		{Type: "t1", Key: "k1"},
		{Type: "t1", Key: "k2"},
		{Type: "t2", Key: "k1"},
		{Type: "t2", Key: "k2"},
	})
	assert.NoError(t, err)
	sort.Slice(getManyActual, getBlobsComparator(getManyActual))
	assert.Equal(
		t,
		[]blobstore.Blob{
			{Type: "t1", Key: "k1", Value: []byte("v1"), Version: 0},
			{Type: "t1", Key: "k2", Value: []byte("v2"), Version: 0},
			{Type: "t2", Key: "k1", Value: []byte("v3"), Version: 2},
			{Type: "t2", Key: "k2", Value: []byte("v4"), Version: 1},
		},
		getManyActual,
	)

	getManyActual, err = store.GetMany("network2", []storage.TypeAndKey{
		{Type: "t3", Key: "k3"},
		{Type: "t3", Key: "k4"},
	})
	assert.NoError(t, err)
	sort.Slice(getManyActual, getBlobsComparator(getManyActual))
	assert.Equal(
		t,
		[]blobstore.Blob{
			{Type: "t3", Key: "k3", Value: []byte("v5"), Version: 0},
			{Type: "t3", Key: "k4", Value: []byte("v6"), Version: 0},
		},
		getManyActual,
	)

	getActual, err = store.Get("network1", storage.TypeAndKey{Type: "t1", Key: "k2"})
	assert.NoError(t, err)
	assert.Equal(t, blobstore.Blob{Type: "t1", Key: "k2", Value: []byte("v2"), Version: 0}, getActual)

	assert.NoError(t, store.Commit())

	// Search tests
	runSearchTestCases(t, fact)

	// Update with creation, read back
	store, err = fact.StartTransaction(nil)
	assert.NoError(t, err)

	err = store.CreateOrUpdate("network1", []blobstore.Blob{
		{Type: "t1", Key: "k1", Value: []byte("hello"), Version: 20},
		{Type: "t9", Key: "k9", Value: []byte("world")},
	})
	assert.NoError(t, err)

	getManyActual, err = store.GetMany("network1", []storage.TypeAndKey{
		{Type: "t1", Key: "k1"},
		{Type: "t9", Key: "k9"},
	})
	assert.NoError(t, err)
	sort.Slice(getManyActual, getBlobsComparator(getManyActual))
	assert.Equal(
		t,
		[]blobstore.Blob{
			{Type: "t1", Key: "k1", Value: []byte("hello"), Version: 20},
			{Type: "t9", Key: "k9", Value: []byte("world"), Version: 0},
		},
		getManyActual,
	)

	assert.NoError(t, store.Commit())

	// Test GetExistingKeys
	store, err = fact.StartTransaction(nil)
	assert.NoError(t, err)
	existingKeys, err := store.GetExistingKeys([]string{"k1", "k9", "k8"}, blobstore.SearchFilter{})
	assert.NoError(t, err)
	assert.Equal(t, []string{"k1", "k9"}, existingKeys)

	network2 := "network2"
	existingKeys, err = store.GetExistingKeys([]string{"k1", "k3", "k4", "k9", "k8"}, blobstore.SearchFilter{NetworkID: &network2})
	assert.NoError(t, err)
	assert.Equal(t, []string{"k3", "k4"}, existingKeys)
	assert.NoError(t, store.Commit())

	// Operation after commit
	_, err = store.Get("network1", storage.TypeAndKey{Type: "t1", Key: "k1"})
	assert.Error(t, err)

	// Delete multiple
	store, err = fact.StartTransaction(nil)
	assert.NoError(t, err)

	err = store.Delete("network1", []storage.TypeAndKey{
		{Type: "t1", Key: "k1"},
		{Type: "t2", Key: "k2"},
	})
	assert.NoError(t, err)

	getManyActual, err = store.GetMany("network1", []storage.TypeAndKey{
		{Type: "t1", Key: "k1"},
		{Type: "t2", Key: "k2"},
		{Type: "t9", Key: "k9"},
	})
	assert.NoError(t, err)
	assert.Equal(t, []blobstore.Blob{{Type: "t9", Key: "k9", Value: []byte("world"), Version: 0}}, getManyActual)

	assert.NoError(t, store.Commit())

	// Delete multiple, rollback, read back
	store, err = fact.StartTransaction(nil)
	assert.NoError(t, err)

	err = store.Delete("network2", []storage.TypeAndKey{
		{Type: "t3", Key: "k3"},
	})
	assert.NoError(t, err)

	// Read back within the tx, should be gone
	getManyActual, err = store.GetMany("network2", []storage.TypeAndKey{
		{Type: "t3", Key: "k3"},
	})
	assert.NoError(t, err)
	assert.Empty(t, getManyActual)
	assert.NoError(t, store.Rollback())

	store, err = fact.StartTransaction(nil)
	assert.NoError(t, err)

	getManyActual, err = store.GetMany("network2", []storage.TypeAndKey{
		{Type: "t3", Key: "k3"},
	})
	assert.NoError(t, err)
	assert.Equal(t, []blobstore.Blob{{Type: "t3", Key: "k3", Value: []byte("v5"), Version: 0}}, getManyActual)
	assert.NoError(t, store.Commit())

	// Increment version
	store, err = fact.StartTransaction(nil)
	assert.NoError(t, err)

	// Non-existent type/key
	err = store.IncrementVersion("network2", storage.TypeAndKey{Type: "t7", Key: "k1"})
	assert.NoError(t, err)

	getManyActual, err = store.GetMany("network2", []storage.TypeAndKey{
		{Type: "t7", Key: "k1"},
	})
	assert.NoError(t, err)
	assert.Equal(t, []blobstore.Blob{{Type: "t7", Key: "k1", Version: 1}}, getManyActual)

	// Increment existing type/key twice
	err = store.IncrementVersion("network2", storage.TypeAndKey{Type: "t3", Key: "k3"})
	assert.NoError(t, err)
	err = store.IncrementVersion("network2", storage.TypeAndKey{Type: "t3", Key: "k3"})
	assert.NoError(t, err)

	getManyActual, err = store.GetMany("network2", []storage.TypeAndKey{
		{Type: "t3", Key: "k3"},
	})
	assert.NoError(t, err)
	assert.Equal(t, []blobstore.Blob{{Type: "t3", Key: "k3", Value: []byte("v5"), Version: 2}}, getManyActual)
}

type searchTestCase struct {
	nid      *string
	types    []string
	keys     []string
	criteria *blobstore.LoadCriteria

	expected map[string][]blobstore.Blob
}

func runSearchTestCases(t *testing.T, fact blobstore.BlobStorageFactory) {
	store, err := fact.StartTransaction(nil)
	assert.NoError(t, err)

	allNetworkSearchTestCases := []searchTestCase{
		{
			// emtpy search filter
			expected: map[string][]blobstore.Blob{
				"network1": {
					{Type: "t1", Key: "k1", Value: []byte("v1"), Version: 0},
					{Type: "t1", Key: "k2", Value: []byte("v2"), Version: 0},
					{Type: "t2", Key: "k1", Value: []byte("v3"), Version: 2},
					{Type: "t2", Key: "k2", Value: []byte("v4"), Version: 1},
				},
				"network2": {
					{Type: "t3", Key: "k3", Value: []byte("v5"), Version: 0},
					{Type: "t3", Key: "k4", Value: []byte("v6"), Version: 0},
				},
			},
		},
		{
			types: []string{"t1", "t3"},
			expected: map[string][]blobstore.Blob{
				"network1": {
					{Type: "t1", Key: "k1", Value: []byte("v1"), Version: 0},
					{Type: "t1", Key: "k2", Value: []byte("v2"), Version: 0},
				},
				"network2": {
					{Type: "t3", Key: "k3", Value: []byte("v5"), Version: 0},
					{Type: "t3", Key: "k4", Value: []byte("v6"), Version: 0},
				},
			},
		},
		{
			// with load criteria
			criteria: &blobstore.LoadCriteria{LoadValue: false},
			types:    []string{"t1", "t3"},
			expected: map[string][]blobstore.Blob{
				"network1": {
					{Type: "t1", Key: "k1", Value: nil, Version: 0},
					{Type: "t1", Key: "k2", Value: nil, Version: 0},
				},
				"network2": {
					{Type: "t3", Key: "k3", Value: nil, Version: 0},
					{Type: "t3", Key: "k4", Value: nil, Version: 0},
				},
			},
		},
		{
			types: []string{"t3"},
			expected: map[string][]blobstore.Blob{
				"network2": {
					{Type: "t3", Key: "k3", Value: []byte("v5"), Version: 0},
					{Type: "t3", Key: "k4", Value: []byte("v6"), Version: 0},
				},
			},
		},
		{
			keys: []string{"k1", "k3"},
			expected: map[string][]blobstore.Blob{
				"network1": {
					{Type: "t1", Key: "k1", Value: []byte("v1"), Version: 0},
					{Type: "t2", Key: "k1", Value: []byte("v3"), Version: 2},
				},
				"network2": {
					{Type: "t3", Key: "k3", Value: []byte("v5"), Version: 0},
				},
			},
		},
		{
			keys: []string{"k3"},
			expected: map[string][]blobstore.Blob{
				"network2": {
					{Type: "t3", Key: "k3", Value: []byte("v5"), Version: 0},
				},
			},
		},
		{
			types: []string{"t1", "t3"},
			keys:  []string{"k1"},
			expected: map[string][]blobstore.Blob{
				"network1": {
					{Type: "t1", Key: "k1", Value: []byte("v1"), Version: 0},
				},
			},
		},
		{
			types:    []string{"t4"},
			keys:     []string{"k1", "k2", "k3", "k4"},
			expected: map[string][]blobstore.Blob{},
		},
	}
	for _, tc := range allNetworkSearchTestCases {
		runSearchTestCase(t, store, tc)
	}

	specificNetworkSearchTestCases := []searchTestCase{
		{
			nid:   strPtr("network1"),
			types: []string{"t1", "t3"},
			expected: map[string][]blobstore.Blob{
				"network1": {
					{Type: "t1", Key: "k1", Value: []byte("v1"), Version: 0},
					{Type: "t1", Key: "k2", Value: []byte("v2"), Version: 0},
				},
			},
		},
		{
			nid:   strPtr("network2"),
			types: []string{"t3"},
			expected: map[string][]blobstore.Blob{
				"network2": {
					{Type: "t3", Key: "k3", Value: []byte("v5"), Version: 0},
					{Type: "t3", Key: "k4", Value: []byte("v6"), Version: 0},
				},
			},
		},
		{
			nid:  strPtr("network1"),
			keys: []string{"k1", "k3"},
			expected: map[string][]blobstore.Blob{
				"network1": {
					{Type: "t1", Key: "k1", Value: []byte("v1"), Version: 0},
					{Type: "t2", Key: "k1", Value: []byte("v3"), Version: 2},
				},
			},
		},
		{
			nid:  strPtr("network2"),
			keys: []string{"k3", "k4"},
			expected: map[string][]blobstore.Blob{
				"network2": {
					{Type: "t3", Key: "k3", Value: []byte("v5"), Version: 0},
					{Type: "t3", Key: "k4", Value: []byte("v6"), Version: 0},
				},
			},
		},
		{
			nid:   strPtr("network1"),
			types: []string{"t1", "t2"},
			keys:  []string{"k1"},
			expected: map[string][]blobstore.Blob{
				"network1": {
					{Type: "t1", Key: "k1", Value: []byte("v1"), Version: 0},
					{Type: "t2", Key: "k1", Value: []byte("v3"), Version: 2},
				},
			},
		},
		{
			nid:   strPtr("network2"),
			types: []string{"t3"},
			keys:  []string{"k1", "k2", "k3", "k4"},
			expected: map[string][]blobstore.Blob{
				"network2": {
					{Type: "t3", Key: "k3", Value: []byte("v5"), Version: 0},
					{Type: "t3", Key: "k4", Value: []byte("v6"), Version: 0},
				},
			},
		},
		{
			nid:      strPtr("network3"),
			types:    []string{"t1", "t2", "t3"},
			keys:     []string{"k1", "k2", "k3", "k4"},
			expected: map[string][]blobstore.Blob{},
		},
	}
	for _, tc := range specificNetworkSearchTestCases {
		runSearchTestCase(t, store, tc)
	}

	assert.NoError(t, store.Commit())
}

func runSearchTestCase(t *testing.T, store blobstore.TransactionalBlobStorage, tc searchTestCase) {
	var criteria blobstore.LoadCriteria
	if tc.criteria != nil {
		criteria = *tc.criteria
	} else {
		criteria = blobstore.GetDefaultLoadCriteria()
	}

	searchActual, err := store.Search(blobstore.CreateSearchFilter(tc.nid, tc.types, tc.keys), criteria)
	assert.NoError(t, err)
	sortSearchOutput(searchActual)
	assert.Equal(t, tc.expected, searchActual)
}

func getTKsComparator(tks []storage.TypeAndKey) func(i, j int) bool {
	return func(i, j int) bool {
		return tks[i].Type+tks[i].Key < tks[j].Type+tks[j].Key
	}
}

func getBlobsComparator(blobs []blobstore.Blob) func(i, j int) bool {
	return func(i, j int) bool {
		return blobs[i].Type+blobs[i].Key < blobs[j].Type+blobs[j].Key
	}
}

func sortSearchOutput(searchActual map[string][]blobstore.Blob) {
	for _, blobs := range searchActual {
		sort.Slice(blobs, getBlobsComparator(blobs))
	}
}

func strPtr(s string) *string {
	return &s
}
