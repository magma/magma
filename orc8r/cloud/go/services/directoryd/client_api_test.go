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

package directoryd_test

import (
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	configurator_test_utils "magma/orc8r/cloud/go/services/configurator/test_utils"
	device_test_init "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/directoryd"
	directoryd_test_init "magma/orc8r/cloud/go/services/directoryd/test_init"
	"magma/orc8r/cloud/go/services/directoryd/types"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/state"
	state_test_init "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/services/state/test_utils"
	state_types "magma/orc8r/cloud/go/services/state/types"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

	"github.com/stretchr/testify/assert"
)

const (
	hn0               = "some_hostname_0"
	hn1               = "some_hostname_1"
	hwid0             = "some_hardware_id_0"
	hwid1             = "some_hardware_id_1"
	imsi0             = "some_imsi_0"
	nid0              = "some_network_id_0"
	sid0              = "some_session_id_0"
	teid0             = "10"
	teids0            = "10,20,30"
	sidWithoutPrefix  = "155129"
	sidWithIMSIPrefix = "IMSI156304337849371-" + sidWithoutPrefix
)

var (
	teid0Slice = []string{"10", "20", "30"}
)

func TestGetSessionID(t *testing.T) {
	record := &types.DirectoryRecord{
		LocationHistory: []string{hwid0}, // imsi0->hwid0
		Identifiers: map[string]interface{}{
			types.RecordKeySessionID: sid0, // imsi0->sid0
		},
	}

	// Default path
	sid, err := record.GetSessionID()
	assert.NoError(t, err)
	assert.Equal(t, sid0, sid)

	// IMSI-prefixed session ID should remove prefix
	record.Identifiers[types.RecordKeySessionID] = sidWithIMSIPrefix
	sid, err = record.GetSessionID()
	assert.NoError(t, err)
	assert.Equal(t, sidWithoutPrefix, sid)

	// Err on non-string sid
	record.Identifiers[types.RecordKeySessionID] = 42
	_, err = record.GetSessionID()
	assert.Error(t, err)

	// Empty string on no sid
	delete(record.Identifiers, types.RecordKeySessionID)
	sid, err = record.GetSessionID()
	assert.NoError(t, err)
	assert.Equal(t, "", sid)
}

func TestGetHWIDForSgwCTeid(t *testing.T) {
	record := &types.DirectoryRecord{
		LocationHistory: []string{hwid0},
		Identifiers: map[string]interface{}{
			types.RecordKeySpgCTeid: teids0,
		},
	}

	// Default path
	teids, err := record.GetSgwCTeids()
	assert.NoError(t, err)
	assert.Equal(t, teid0Slice, teids)

	// Err on non-string teid
	record.Identifiers[types.RecordKeySpgCTeid] = 10
	_, err = record.GetSgwCTeids()
	assert.Error(t, err)

	// Error on empty teid
	record.Identifiers = map[string]interface{}{}
	teids, err = record.GetSgwCTeids()
	assert.NoError(t, err)
	assert.Exactly(t, []string{}, teids)
}

func TestDirectorydMethods(t *testing.T) {
	directoryd_test_init.StartTestService(t)

	// Empty initially
	_, err := directoryd.GetSessionIDForIMSI(nid0, imsi0)
	assert.Error(t, err)
	_, err = directoryd.GetHostnameForHWID(hwid0)
	assert.Error(t, err)
	_, err = directoryd.GetHWIDForSgwCTeid(nid0, teid0)
	assert.Error(t, err)

	// Put sid0->imsi0
	err = directoryd.MapSessionIDsToIMSIs(nid0, map[string]string{sid0: imsi0})
	assert.NoError(t, err)

	// Put Many hwid0->hn0
	err = directoryd.MapHWIDsToHostnames(map[string]string{hwid0: hn0})
	assert.NoError(t, err)

	// Put Single hwid1->hn1
	err = directoryd.MapHWIDToHostname(hwid1, hn1)
	assert.NoError(t, err)

	// Get sid0->imsi0
	imsi, err := directoryd.GetIMSIForSessionID(nid0, sid0)
	assert.NoError(t, err)
	assert.Equal(t, imsi, imsi0)

	// Get hwid0->hn0
	hn, err := directoryd.GetHostnameForHWID(hwid0)
	assert.NoError(t, err)
	assert.Equal(t, hn0, hn)

	// Get hwid1->hn1
	hn, err = directoryd.GetHostnameForHWID(hwid1)
	assert.NoError(t, err)
	assert.Equal(t, hn1, hn)

	// Put teid->hwid
	err = directoryd.MapSgwCTeidToHWID(nid0, map[string]string{teid0: hwid0})
	assert.NoError(t, err)

	// Get teid->hwid
	hwid, err := directoryd.GetHWIDForSgwCTeid(nid0, teid0)
	assert.NoError(t, err)
	assert.Equal(t, hwid0, hwid)
}

func TestDirectorydStateMethods(t *testing.T) {
	configurator_test_init.StartTestService(t)
	device_test_init.StartTestService(t)

	directoryd_test_init.StartTestService(t)
	state_test_init.StartTestService(t)

	stateClient, err := getStateServiceClient(t)
	assert.NoError(t, err)

	configurator_test_utils.RegisterNetwork(t, nid0, "DirectoryD Service Test")
	configurator_test_utils.RegisterGateway(t, nid0, hwid0, &models.GatewayDevice{HardwareID: hwid0})
	ctx := test_utils.GetContextWithCertificate(t, hwid0)

	record := &types.DirectoryRecord{
		LocationHistory: []string{hwid0}, // imsi0->hwid0
		Identifiers: map[string]interface{}{
			types.RecordKeySessionID: sid0, // imsi0->sid0
		},
	}
	serializedRecord, err := record.MarshalBinary()
	assert.NoError(t, err)

	st := &protos.State{
		Type:     orc8r.DirectoryRecordType,
		DeviceID: imsi0,
		Value:    serializedRecord,
	}
	stateID := state_types.ID{
		Type:     st.Type,
		DeviceID: st.DeviceID,
	}

	// Empty initially
	_, err = directoryd.GetHWIDForIMSI(nid0, imsi0)
	assert.Error(t, err)
	_, err = directoryd.GetSessionIDForIMSI(nid0, imsi0)
	assert.Error(t, err)

	// Report state
	reqReport := &protos.ReportStatesRequest{States: []*protos.State{st}}
	res, err := stateClient.ReportStates(ctx, reqReport)
	assert.NoError(t, err)
	assert.Empty(t, res.UnreportedStates)

	// Get imsi0->hwid0
	hwid, err := directoryd.GetHWIDForIMSI(nid0, imsi0)
	assert.NoError(t, err)
	assert.Equal(t, hwid0, hwid)

	// Get imsi0->sid0
	sid, err := directoryd.GetSessionIDForIMSI(nid0, imsi0)
	assert.NoError(t, err)
	assert.Equal(t, sid0, sid)

	// Delete state
	err = state.DeleteStates(nid0, state_types.IDs{stateID})
	assert.NoError(t, err)

	// Get imsi0->hwid0, should be gone
	_, err = directoryd.GetHWIDForIMSI(nid0, imsi0)
	assert.Error(t, err)

	// Get imsi0->sid0, should be gone
	_, err = directoryd.GetSessionIDForIMSI(nid0, imsi0)
	assert.Error(t, err)
}

func TestDirectorydUpdateMethods(t *testing.T) {
	configurator_test_init.StartTestService(t)
	device_test_init.StartTestService(t)

	directoryd_test_init.StartTestService(t)
	state_test_init.StartTestService(t)

	ddUpdaterClient, err := getDirectorydUpdaterClient(t)
	assert.NoError(t, err)

	configurator_test_utils.RegisterNetwork(t, nid0, "DirectoryD Service Test")
	configurator_test_utils.RegisterGateway(t, nid0, hwid0, &models.GatewayDevice{HardwareID: hwid0})
	ctx := test_utils.GetContextWithCertificate(t, hwid0)

	// Empty initially
	_, err = directoryd.GetHWIDForIMSI(nid0, imsi0)
	assert.Error(t, err)
	_, err = directoryd.GetSessionIDForIMSI(nid0, imsi0)
	assert.Error(t, err)

	// Update
	_, err = ddUpdaterClient.UpdateRecord(ctx, &protos.UpdateRecordRequest{
		Id:       imsi0,
		Location: hwid0,
		Fields:   map[string]string{types.RecordKeySessionID: sid0},
	})
	assert.NoError(t, err)

	// Get imsi0->hwid0
	hwid, err := directoryd.GetHWIDForIMSI(nid0, imsi0)
	assert.NoError(t, err)
	assert.Equal(t, hwid0, hwid)

	// Get imsi0->sid0
	sid, err := directoryd.GetSessionIDForIMSI(nid0, imsi0)
	assert.NoError(t, err)
	assert.Equal(t, sid0, sid)

	// Get Field
	field, err := ddUpdaterClient.GetDirectoryField(
		ctx, &protos.GetDirectoryFieldRequest{Id: imsi0, FieldKey: types.RecordKeySessionID})
	assert.NoError(t, err)
	assert.Equal(t, types.RecordKeySessionID, field.GetKey())
	assert.Equal(t, sid0, field.GetValue())

	records, err := ddUpdaterClient.GetAllDirectoryRecords(ctx, &protos.Void{})
	assert.NoError(t, err)
	assert.Equal(t, int(1), len(records.GetRecords()))

	// Delete
	_, err = ddUpdaterClient.DeleteRecord(ctx, &protos.DeleteRecordRequest{Id: imsi0})
	assert.NoError(t, err)

	// Get imsi0->hwid0, should be gone
	_, err = directoryd.GetHWIDForIMSI(nid0, imsi0)
	assert.Error(t, err)

	// Get imsi0->sid0, should be gone
	_, err = directoryd.GetSessionIDForIMSI(nid0, imsi0)
	assert.Error(t, err)
}

func getStateServiceClient(t *testing.T) (protos.StateServiceClient, error) {
	conn, err := registry.GetConnection(state.ServiceName)
	assert.NoError(t, err)
	return protos.NewStateServiceClient(conn), err
}

func getDirectorydUpdaterClient(t *testing.T) (protos.GatewayDirectoryServiceClient, error) {
	conn, err := registry.GetConnection(directoryd.ServiceName)
	assert.NoError(t, err)
	return protos.NewGatewayDirectoryServiceClient(conn), err
}
