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

	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/indexer/mocks"
	"magma/orc8r/cloud/go/services/state/indexer/reindex"
	"magma/orc8r/cloud/go/sqorc"
)

func TestVersioner(t *testing.T) {
	versioner := initVersioner(t)

	// Empty initially
	v, err := versioner.GetIndexerVersions()
	assert.NoError(t, err)
	assert.Empty(t, v)

	// Write some versions, ensure they stuck
	want := []*indexer.Versions{
		{IndexerID: id0, Actual: zero, Desired: version0},
		{IndexerID: id1, Actual: zero, Desired: version1},
		{IndexerID: id2, Actual: zero, Desired: version2},
	}

	// Start and register indexer servicers
	mocks.NewMockIndexer(t, id0, version0, nil, nil, nil, nil)
	mocks.NewMockIndexer(t, id1, version1, nil, nil, nil, nil)
	mocks.NewMockIndexer(t, id2, version2, nil, nil, nil, nil)

	assert.NoError(t, err)
	got, err := versioner.GetIndexerVersions()
	assert.NoError(t, err)
	assert.Equal(t, want, got)

	// Update one actual version
	err = versioner.SetIndexerActualVersion(id2, version2)
	assert.NoError(t, err)
	gotv, err := reindex.GetIndexerVersion(versioner, id2)
	assert.NoError(t, err)
	assert.Equal(t, version2, gotv.Actual)

	// Bump desired version for same indexer
	mocks.NewMockIndexer(t, id2, version2a, nil, nil, nil, nil)
	assert.NoError(t, err)
	got, err = versioner.GetIndexerVersions()
	assert.NoError(t, err)
	want = []*indexer.Versions{
		{IndexerID: id0, Actual: zero, Desired: version0},
		{IndexerID: id1, Actual: zero, Desired: version1},
		{IndexerID: id2, Actual: version2, Desired: version2a},
	}
	assert.Equal(t, want, got)
}

func initVersioner(t *testing.T) reindex.Versioner {
	indexer.DeregisterAllForTest(t)
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	v := reindex.NewVersioner(db, sqorc.GetSqlBuilder())
	err = v.Initialize()
	assert.NoError(t, err)
	return v
}
