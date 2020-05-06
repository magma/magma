/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package servicers_test

import (
	"context"
	"testing"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/blobstore/mocks"
	"magma/orc8r/cloud/go/services/state/servicers"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TODO: fill out more test cases for servicer test with mocked storage

func TestStateServicer_GetStates(t *testing.T) {
	// mock setup: expect 1 RPC to result in a search, the other to a concrete
	// GetMany
	mockStore := &mocks.TransactionalBlobStorage{}
	mockStore.On("Search",
		blobstore.CreateSearchFilter(strPtr("network1"), []string{"t1", "t2"}, []string{"k1", "k2"}),
		blobstore.GetDefaultLoadCriteria(),
	).
		Return(map[string][]blobstore.Blob{
			"network1": {
				{Type: "t1", Key: "k1", Value: []byte("v1"), Version: 42},
				{Type: "t2", Key: "k2", Value: []byte("v2"), Version: 43},
			},
		}, nil)
	mockStore.On("GetMany", "network1", []storage.TypeAndKey{{Type: "t1", Key: "k1"}, {Type: "t2", Key: "k2"}}).
		Return([]blobstore.Blob{
			{Type: "t1", Key: "k1", Value: []byte("v1"), Version: 42},
			{Type: "t2", Key: "k2", Value: []byte("v2"), Version: 43},
		}, nil)

	mockStore.On("Commit").Return(nil)

	fact := &mocks.BlobStorageFactory{}
	fact.On("StartTransaction", mock.Anything).Return(mockStore, nil)

	srv, err := servicers.NewStateServicer(fact)
	assert.NoError(t, err)

	actual, err := srv.GetStates(context.Background(), &protos.GetStatesRequest{
		NetworkID:  "network1",
		TypeFilter: []string{"t1", "t2"},
		IdFilter:   []string{"k1", "k2"},
		LoadValues: true,
	})
	assert.NoError(t, err)
	expected := &protos.GetStatesResponse{
		States: []*protos.State{
			{
				Type:     "t1",
				DeviceID: "k1",
				Value:    []byte("v1"),
				Version:  42,
			},
			{
				Type:     "t2",
				DeviceID: "k2",
				Value:    []byte("v2"),
				Version:  43,
			},
		},
	}
	assert.Equal(t, expected, actual)

	// Prefer concrete GetMany over Search
	actual, err = srv.GetStates(context.Background(), &protos.GetStatesRequest{
		NetworkID: "network1",
		Ids: []*protos.StateID{
			{Type: "t1", DeviceID: "k1"},
			{Type: "t2", DeviceID: "k2"},
		},
		TypeFilter: []string{"t1", "t2"},
		IdFilter:   []string{"k1", "k2"},
		LoadValues: false,
	})
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	mockStore.AssertExpectations(t)
	fact.AssertExpectations(t)
}

func strPtr(s string) *string {
	return &s
}
