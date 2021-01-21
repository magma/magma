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

package indexer_test

import (
	"sort"
	"testing"

	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/indexer/mocks"

	"github.com/stretchr/testify/assert"
)

func TestRegisterRemote(t *testing.T) {
	want0 := indexer.NewRemoteIndexer("some_service_0", 42, "type0", "type1")
	want1 := indexer.NewRemoteIndexer("some_service_1", 420)
	want2 := indexer.NewRemoteIndexer("some_service_2", 424, "type2")

	indexer.DeregisterAllForTest(t)

	t.Run("empty initially", func(t *testing.T) {
		got, err := indexer.GetIndexers()
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("set and get one", func(t *testing.T) {
		register(t, want0)

		got, err := indexer.GetIndexers()
		assert.NoError(t, err)
		assert.Equal(t, []indexer.Indexer{want0}, got)
		gotOne, err := indexer.GetIndexer(want0.GetID())
		assert.NoError(t, err)
		assert.Equal(t, want0, gotOne)
	})

	t.Run("set and get two more", func(t *testing.T) {
		register(t, want1, want2)

		got, err := indexer.GetIndexers()
		assert.NoError(t, err)
		sort.Slice(got, func(i, j int) bool { return got[i].GetID() < got[j].GetID() })
		assert.Equal(t, []indexer.Indexer{want0, want1, want2}, got)
		got1, err := indexer.GetIndexer(want1.GetID())
		assert.NoError(t, err)
		assert.Equal(t, want1, got1)
		got2, err := indexer.GetIndexer(want2.GetID())
		assert.NoError(t, err)
		assert.Equal(t, want2, got2)
	})

	t.Run("fail overwrite same name", func(t *testing.T) {
		register(t, want2)

		got, err := indexer.GetIndexers()
		assert.NoError(t, err)
		sort.Slice(got, func(i, j int) bool { return got[i].GetID() < got[j].GetID() })
		assert.Equal(t, []indexer.Indexer{want0, want1, want2}, got)
	})

	t.Run("get indexers for state type", func(t *testing.T) {
		got, err := indexer.GetIndexersForState("type2")
		assert.NoError(t, err)
		assert.Equal(t, []indexer.Indexer{want2}, got)
	})
}

func register(t *testing.T, indexers ...indexer.Indexer) {
	for _, x := range indexers {
		mocks.NewMockIndexer(t, x.GetID(), x.GetVersion(), x.GetTypes(), nil, nil, nil) // registers the indexer
	}
}
