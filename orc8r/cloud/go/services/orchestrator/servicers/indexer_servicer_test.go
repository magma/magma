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
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	configurator_test_utils "magma/orc8r/cloud/go/services/configurator/test_utils"
	device_test_init "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/directoryd"
	directoryd_test_init "magma/orc8r/cloud/go/services/directoryd/test_init"
	directoryd_types "magma/orc8r/cloud/go/services/directoryd/types"
	"magma/orc8r/cloud/go/services/orchestrator"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	orchestrator_test_init "magma/orc8r/cloud/go/services/orchestrator/test_init"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/indexer"
	state_test_init "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/services/state/test_utils"
	state_types "magma/orc8r/cloud/go/services/state/types"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

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

func TestIndexerRecordIDs(t *testing.T) {
	const (
		version indexer.Version = 1 // copied from indexer_servicer.go

		hwid0 = "hwid0"
		hwid1 = "hwid1"
		imsi0 = "some_imsi_0"
		imsi1 = "some_imsi_1"
		imsi2 = "some_imsi_2"
		sid0  = "some_session_id_0"
		sid1  = "some_session_id_1"
		sid2  = "some_session_id_2"
		nid0  = "some_network_id_0"
	)
	var (
		types = []string{orc8r.DirectoryRecordType} // copied from indexer_servicer.go
	)
	configurator_test_init.StartTestService(t)
	device_test_init.StartTestService(t)
	directoryd_test_init.StartTestService(t)
	orchestrator_test_init.StartTestService(t)
	state_test_init.StartTestService(t)

	configurator_test_utils.RegisterNetwork(t, nid0, "DirectoryD Service Test")
	configurator_test_utils.RegisterGateway(t, nid0, hwid0, &models.GatewayDevice{HardwareID: hwid0})
	configurator_test_utils.RegisterGateway(t, nid0, hwid1, &models.GatewayDevice{HardwareID: hwid1})

	idx := indexer.NewRemoteIndexer(orchestrator.ServiceName, version, types...)

	record0 := &directoryd_types.DirectoryRecord{
		LocationHistory: []string{hwid0},
		Identifiers: map[string]interface{}{
			directoryd_types.RecordKeySessionID: sid0,
		},
	}
	record1 := &directoryd_types.DirectoryRecord{
		LocationHistory: []string{hwid1},
		Identifiers: map[string]interface{}{
			directoryd_types.RecordKeySessionID: sid1,
		},
	}
	record2 := &directoryd_types.DirectoryRecord{
		LocationHistory: []string{hwid0},
		Identifiers: map[string]interface{}{
			directoryd_types.RecordKeySessionID: sid2,
		},
	}

	id0 := state_types.ID{
		Type:     orc8r.DirectoryRecordType,
		DeviceID: imsi0,
	}
	id1 := state_types.ID{
		Type:     orc8r.DirectoryRecordType,
		DeviceID: imsi1,
	}
	id2 := state_types.ID{
		Type:     orc8r.DirectoryRecordType,
		DeviceID: imsi2,
	}
	st0 := state_types.State{
		ReportedState:      record0,
		Version:            44,
		TimeMs:             42,
		CertExpirationTime: 43,
	}
	st1 := state_types.State{
		ReportedState:      record1,
		Version:            47,
		TimeMs:             42,
		CertExpirationTime: 43,
	}
	st2 := state_types.State{
		ReportedState:      record2,
		Version:            49,
		TimeMs:             42,
		CertExpirationTime: 43,
	}
	reportDirectorydState(t, hwid0, imsi0)
	reportDirectorydState(t, hwid0, imsi2)
	reportDirectorydState(t, hwid1, imsi1)

	// Indexer subscription matches directory records
	assert.True(t, len(idx.GetTypes()) > 0)
	assert.True(t, idx.GetTypes()[0] == orc8r.DirectoryRecordType)

	// Index st0, result is hwid0->[imsi0] reverse mapping
	errs, err := idx.Index(nid0, state_types.SerializedStatesByID{id0: serialize(t, st0, orc8r.DirectoryRecordType)})
	assert.NoError(t, err)
	assert.Empty(t, errs)
	expectedRecordIDs := []string{imsi0}
	recordIDs, err := directoryd.GetHWIDToDirectoryRecordIDs(nid0, hwid0)
	assert.NoError(t, err)
	assert.Equal(t, expectedRecordIDs, recordIDs)

	// Index st1 and st2 result is hwid0->[imsi0, imsi2], hwid1->[imsi1] reverse mapping
	errs, err = idx.Index(nid0, state_types.SerializedStatesByID{
		id1: serialize(t, st1, orc8r.DirectoryRecordType),
		id2: serialize(t, st2, orc8r.DirectoryRecordType),
	})
	assert.NoError(t, err)
	assert.Empty(t, errs)

	expectedRecordIDs = []string{imsi0, imsi2}
	recordIDs, err = directoryd.GetHWIDToDirectoryRecordIDs(nid0, hwid0)
	assert.NoError(t, err)
	assert.Equal(t, expectedRecordIDs, recordIDs)

	expectedRecordIDs = []string{imsi1}
	recordIDs, err = directoryd.GetHWIDToDirectoryRecordIDs(nid0, hwid1)
	assert.NoError(t, err)
	assert.Equal(t, expectedRecordIDs, recordIDs)

	// Ensure deleted state is pruned
	deleteDirectorydState(t, nid0, hwid1, imsi1)
	expectedRecordIDs = nil
	recordIDs, err = directoryd.GetHWIDToDirectoryRecordIDs(nid0, hwid1)
	assert.NoError(t, err)
	assert.Equal(t, expectedRecordIDs, recordIDs)

	// No errs when when can't deserialize state -- just logs
	errs, err = idx.Index(nid0, state_types.SerializedStatesByID{id0: state_types.SerializedState{SerializedReportedState: []byte("0xdeadbeef")}})
	assert.NoError(t, err)
	assert.Empty(t, errs)
	expectedRecordIDs = []string{imsi0, imsi2}
	recordIDs, err = directoryd.GetHWIDToDirectoryRecordIDs(nid0, hwid0)
	assert.NoError(t, err)
	assert.Equal(t, expectedRecordIDs, recordIDs)

	// Err when can deserialize but is wrong type
	id0.Type = orc8r.GatewayStateType
	st0.ReportedState = &models.GatewayStatus{Meta: map[string]string{"foo": "bar"}}
	errs, err = idx.Index(nid0, state_types.SerializedStatesByID{id0: serialize(t, st0, orc8r.GatewayStateType)})
	assert.NoError(t, err)
	assert.Error(t, errs[id0])
}

func reportDirectorydState(t *testing.T, hwid string, imsi string) {
	stateClient, err := getStateServiceClient(t)
	assert.NoError(t, err)
	ctx := test_utils.GetContextWithCertificate(t, hwid)
	record := &directoryd_types.DirectoryRecord{
		LocationHistory: []string{hwid},
		Identifiers: map[string]interface{}{
			directoryd_types.RecordKeySessionID: "sid_foo",
		},
	}
	serializedRecord, err := record.MarshalBinary()
	assert.NoError(t, err)

	st := &protos.State{
		Type:     orc8r.DirectoryRecordType,
		DeviceID: imsi,
		Value:    serializedRecord,
	}
	reqReport := &protos.ReportStatesRequest{States: []*protos.State{st}}
	res, err := stateClient.ReportStates(ctx, reqReport)
	assert.NoError(t, err)
	assert.Empty(t, res.UnreportedStates)
}

func deleteDirectorydState(t *testing.T, nid string, hwid string, imsi string) {
	stateClient, err := getStateServiceClient(t)
	assert.NoError(t, err)
	ctx := test_utils.GetContextWithCertificate(t, hwid)
	stId := &protos.StateID{
		Type:     orc8r.DirectoryRecordType,
		DeviceID: imsi,
	}
	reqReport := &protos.DeleteStatesRequest{NetworkID: nid, Ids: []*protos.StateID{stId}}
	_, err = stateClient.DeleteStates(ctx, reqReport)
	assert.NoError(t, err)
}

func getStateServiceClient(t *testing.T) (protos.StateServiceClient, error) {
	conn, err := registry.GetConnection(state.ServiceName)
	assert.NoError(t, err)
	return protos.NewStateServiceClient(conn), err
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
