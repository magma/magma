/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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
	blobs := map[string]blobstore.Blobs{
		"nid0": {{Type: "typeA", Key: "keyA"}},
		"nid1": {{Type: "typeB", Key: "keyB"}},
	}
	ids := state_types.IDsByNetwork{
		"nid0": {{Type: "typeA", DeviceID: "keyA"}},
		"nid1": {{Type: "typeB", DeviceID: "keyB"}},
	}

	store := &mocks.TransactionalBlobStorage{}
	store.On("Search",
		blobstore.CreateSearchFilter(nil, nil, nil, nil),
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
