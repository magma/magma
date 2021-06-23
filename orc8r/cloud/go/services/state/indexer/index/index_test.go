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

package index_test

import (
	"testing"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/indexer/index"
	"magma/orc8r/cloud/go/services/state/indexer/mocks"
	state_test_init "magma/orc8r/cloud/go/services/state/test_init"
	state_types "magma/orc8r/cloud/go/services/state/types"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func init() {
	//_ = flag.Set("alsologtostderr", "true") // uncomment to view logs during test
}

func TestIndexImpl_HappyPath(t *testing.T) {
	const (
		maxRetry = 3 // copied from index.go

		nid0 = "some_networkid_0"

		iid0 = "some_indexerid_0"
		iid1 = "some_indexerid_1"
		iid2 = "some_indexerid_2"
		iid3 = "some_indexerid_3"
	)
	var (
		someErr = errors.New("some_error")
	)

	clock.SkipSleeps(t)
	defer clock.ResumeSleeps(t)

	id0 := state_types.ID{Type: orc8r.GatewayStateType}
	id1 := state_types.ID{Type: orc8r.StringMapSerdeType}
	reported0 := &models.GatewayStatus{Meta: map[string]string{"foo": "bar"}}
	reported1 := &state.StringToStringMap{"apple": "banana"}
	st0 := state_types.State{ReportedState: reported0}
	st1 := state_types.State{ReportedState: reported1}

	indexTwo := state_types.StatesByID{id0: st0, id1: st1}
	indexOne := state_types.StatesByID{id1: st1}
	in := state_types.StatesByID{id0: st0, id1: st1}

	idx0 := getIndexer(iid0, []string{orc8r.GatewayStateType, orc8r.StringMapSerdeType})
	idx1 := getIndexer(iid1, []string{orc8r.StringMapSerdeType})
	idx2 := getIndexer(iid2, []string{"type_with_no_reported_states"})
	idx3 := getIndexer(iid3, []string{})

	idx0.On("Index", nid0, serialize(t, indexTwo)).Return(state_types.StateErrors{id0: someErr}, nil).Once()
	idx1.On("Index", nid0, serialize(t, indexOne)).Return(nil, someErr).Times(maxRetry)
	idx0.On("GetVersion").Return(indexer.Version(42))
	idx1.On("GetVersion").Return(indexer.Version(42))
	idx2.On("GetVersion").Return(indexer.Version(42))
	idx3.On("GetVersion").Return(indexer.Version(42))

	indexer.DeregisterAllForTest(t)
	state_test_init.StartNewTestIndexer(t, idx0)
	state_test_init.StartNewTestIndexer(t, idx1)
	state_test_init.StartNewTestIndexer(t, idx2)
	state_test_init.StartNewTestIndexer(t, idx3)

	// All indexing occurs as expected
	actual, err := index.Index(nid0, serialize(t, in))
	assert.NoError(t, err)
	assert.Len(t, actual, 1) // from idx1's overarching err return
	e := actual[0].Error()
	assert.Contains(t, e, iid1)
	assert.Contains(t, e, index.ErrIndex)
	assert.Contains(t, e, someErr.Error())
	idx0.AssertExpectations(t)
	idx1.AssertExpectations(t)
	idx2.AssertExpectations(t)
	idx3.AssertExpectations(t)
}

func getIndexer(id string, types []string) *mocks.Indexer {
	idx := &mocks.Indexer{}
	idx.On("GetID").Return(id)
	idx.On("GetTypes").Return(types)
	return idx
}

func serialize(t *testing.T, states state_types.StatesByID) state_types.SerializedStatesByID {
	serialized := state_types.SerializedStatesByID{}
	for id, st := range states {
		s := state_types.SerializedState{
			Version:    st.Version,
			ReporterID: st.ReporterID,
			TimeMs:     st.TimeMs,
		}
		rep, err := serde.Serialize(st.ReportedState, id.Type, serdes.State)
		assert.NoError(t, err)
		s.SerializedReportedState = rep
		serialized[id] = s
	}
	return serialized
}
