/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package indexer_test

import (
	"testing"

	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/indexer/mocks"

	"github.com/stretchr/testify/assert"
)

func TestRegisterIndexers(t *testing.T) {
	id0 := "some_id_0"
	id1 := "some_id_1"
	var idx0 *mocks.Indexer
	var idx1 *mocks.Indexer

	// Try to register multiple indexers with same id
	idx0 = &mocks.Indexer{}
	idx1 = &mocks.Indexer{}
	idx0.On("GetID").Return(id0).Times(2)
	idx1.On("GetID").Return(id0).Once()

	err := indexer.RegisterAll(idx0, idx1)
	assert.Error(t, err)
	idx0.AssertExpectations(t)
	idx1.AssertExpectations(t)

	// Success
	idx0 = &mocks.Indexer{}
	idx1 = &mocks.Indexer{}
	idx0.On("GetID").Return(id0).Once()
	idx1.On("GetID").Return(id1).Once()

	err = indexer.RegisterAll(idx0, idx1)
	assert.NoError(t, err)
	idx0.AssertExpectations(t)
	idx1.AssertExpectations(t)

	idx, err := indexer.GetIndexer(id0)
	assert.NoError(t, err)
	assert.Equal(t, idx0, idx)
	idx, err = indexer.GetIndexer(id1)
	assert.NoError(t, err)
	assert.Equal(t, idx1, idx)

	idxs := indexer.GetAllIndexers()
	assert.Contains(t, idxs, idx0)
	assert.Contains(t, idxs, idx1)
}
