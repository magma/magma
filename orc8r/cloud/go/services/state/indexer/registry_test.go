/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package indexer_test

import (
	"github.com/stretchr/testify/assert"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/indexer/mocks"
	"testing"
)

var id0 = "some_id_0"
var id1 = "some_id_1"

func TestRegisterIndexers(t *testing.T) {
	var idx0 *mocks.Indexer
	var idx1 *mocks.Indexer

	// Try to register multiple indexers with same id
	idx0 = &mocks.Indexer{}
	idx1 = &mocks.Indexer{}
	idx0.On("GetID").Return(id0).Times(2)
	idx1.On("GetID").Return(id0).Once()

	err := indexer.RegisterIndexers(idx0, idx1)
	assert.Error(t, err)
	idx0.AssertExpectations(t)
	idx1.AssertExpectations(t)

	// Success
	idx0 = &mocks.Indexer{}
	idx1 = &mocks.Indexer{}
	idx0.On("GetID").Return(id0).Once()
	idx1.On("GetID").Return(id1).Once()

	err = indexer.RegisterIndexers(idx0, idx1)
	assert.NoError(t, err)
	idx0.AssertExpectations(t)
	idx1.AssertExpectations(t)
}
