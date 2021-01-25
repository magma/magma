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

package servicers_test

import (
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/directoryd"
	directoryd_test_init "magma/orc8r/cloud/go/services/directoryd/test_init"
	directoryd_types "magma/orc8r/cloud/go/services/directoryd/types"
	"magma/orc8r/cloud/go/services/orchestrator"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	orchestrator_test_init "magma/orc8r/cloud/go/services/orchestrator/test_init"
	"magma/orc8r/cloud/go/services/state/indexer"
	state_types "magma/orc8r/cloud/go/services/state/types"

	"github.com/stretchr/testify/assert"
)

func TestIndexerSessionID(t *testing.T) {
	const (
		version indexer.Version = 1 // copied from indexer_servicer.go

		imsi0 = "some_imsi_0"
		imsi1 = "some_imsi_1"
		nid0  = "some_network_id_0"
		sid0  = "some_session_id_0"
		sid1  = "some_session_id_1"
	)
	var (
		types = []string{orc8r.DirectoryRecordType} // copied from indexer_servicer.go
	)

	directoryd_test_init.StartTestService(t)
	orchestrator_test_init.StartTestService(t)
	idx := indexer.NewRemoteIndexer(orchestrator.ServiceName, version, types...)

	record := &directoryd_types.DirectoryRecord{
		Identifiers: map[string]interface{}{
			directoryd_types.RecordKeySessionID: sid0, // imsi0->sid0
		},
		LocationHistory: []string{"apple"},
	}

	id := state_types.ID{
		Type:     orc8r.DirectoryRecordType,
		DeviceID: imsi0,
	}
	st := state_types.State{
		ReportedState:      record,
		Version:            44,
		TimeMs:             42,
		CertExpirationTime: 43,
	}

	// Indexer subscription matches directory records
	assert.True(t, len(idx.GetTypes()) > 0)
	assert.True(t, idx.GetTypes()[0] == orc8r.DirectoryRecordType)

	// Index the imsi0->sid0 state, result is sid0->imsi0 reverse mapping
	errs, err := idx.Index(nid0, state_types.SerializedStatesByID{id: serialize(t, st, orc8r.DirectoryRecordType)})
	assert.NoError(t, err)
	assert.Empty(t, errs)
	imsi, err := directoryd.GetIMSIForSessionID(nid0, sid0)
	assert.NoError(t, err)
	assert.Equal(t, imsi0, imsi)

	// Update sid -- index imsi0->sid1, result is sid1->imsi0 reverse mapping
	// Note that we specifically don't test for the presence of {sid0 -> ?}, as we allow stale derived state to persist.
	st.ReportedState.(*directoryd_types.DirectoryRecord).Identifiers[directoryd_types.RecordKeySessionID] = sid1
	errs, err = idx.Index(nid0, state_types.SerializedStatesByID{id: serialize(t, st, orc8r.DirectoryRecordType)})
	assert.NoError(t, err)
	assert.Empty(t, errs)
	imsi, err = directoryd.GetIMSIForSessionID(nid0, sid1)
	assert.NoError(t, err)
	assert.Equal(t, imsi0, imsi)

	// Update imsi -- index imsi1->sid1, result is sid1->imsi1 reverse mapping
	// Note that we specifically don't test for the presence of {sid0 -> ?}, as we allow stale derived state to persist.
	id.DeviceID = imsi1
	errs, err = idx.Index(nid0, state_types.SerializedStatesByID{id: serialize(t, st, orc8r.DirectoryRecordType)})
	assert.NoError(t, err)
	assert.Empty(t, errs)
	imsi, err = directoryd.GetIMSIForSessionID(nid0, sid1)
	assert.NoError(t, err)
	assert.Equal(t, imsi1, imsi)

	// No errs when when can't deserialize state -- just logs
	errs, err = idx.Index(nid0, state_types.SerializedStatesByID{id: state_types.SerializedState{SerializedReportedState: []byte("0xdeadbeef")}})
	assert.NoError(t, err)
	assert.Empty(t, errs)
	imsi, err = directoryd.GetIMSIForSessionID(nid0, sid1)
	assert.NoError(t, err)
	assert.Equal(t, imsi1, imsi)

	// Err when can deserialize but is wrong type
	id.Type = orc8r.GatewayStateType
	st.ReportedState = &models.GatewayStatus{Meta: map[string]string{"foo": "bar"}}
	errs, err = idx.Index(nid0, state_types.SerializedStatesByID{id: serialize(t, st, orc8r.GatewayStateType)})
	assert.NoError(t, err)
	assert.Error(t, errs[id])
	imsi, err = directoryd.GetIMSIForSessionID(nid0, sid1)
	assert.NoError(t, err)
	assert.Equal(t, imsi1, imsi)
}

func serialize(t *testing.T, st state_types.State, typ string) state_types.SerializedState {
	s := state_types.SerializedState{
		Version:            st.Version,
		ReporterID:         st.ReporterID,
		TimeMs:             st.TimeMs,
		CertExpirationTime: st.CertExpirationTime,
	}
	rep, err := serde.Serialize(st.ReportedState, typ, serdes.State)
	assert.NoError(t, err)
	s.SerializedReportedState = rep
	return s
}
