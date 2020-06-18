/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package reindex_test

import (
	"testing"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/blobstore/mocks"
	"magma/orc8r/cloud/go/services/state/indexer/reindex"
	state_types "magma/orc8r/cloud/go/services/state/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestIndexerServicer_GetAllIDs(t *testing.T) {
	blobs := map[string][]blobstore.Blob{
		"nid0": {{Type: "typeA", Key: "keyA"}},
		"nid1": {{Type: "typeB", Key: "keyB"}},
	}
	ids := state_types.IDsByNetwork{
		"nid0": {{Type: "typeA", DeviceID: "keyA"}},
		"nid1": {{Type: "typeB", DeviceID: "keyB"}},
	}

	store := &mocks.TransactionalBlobStorage{}
	store.On("Search",
		blobstore.CreateSearchFilter(nil, nil, nil),
		blobstore.LoadCriteria{LoadValue: false},
	).Return(blobs, nil)
	store.On("Commit").Return(nil)
	fact := &mocks.BlobStorageFactory{}
	fact.On("StartTransaction", mock.Anything).Return(store, nil)

	st := reindex.NewStore(fact)
	got, err := st.GetAllIDs()
	assert.NoError(t, err)
	assert.Equal(t, ids, got)
}
