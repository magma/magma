/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package directoryd_test

import (
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/serde"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	configurator_test_utils "magma/orc8r/cloud/go/services/configurator/test_utils"
	"magma/orc8r/cloud/go/services/device"
	device_test_init "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/directoryd"
	directoryd_test_init "magma/orc8r/cloud/go/services/directoryd/test_init"
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
	sidWithoutPrefix  = "155129"
	sidWithIMSIPrefix = "IMSI156304337849371-" + sidWithoutPrefix
)

func TestGetSessionID(t *testing.T) {
	record := &directoryd.DirectoryRecord{
		LocationHistory: []string{hwid0}, // imsi0->hwid0
		Identifiers: map[string]interface{}{
			directoryd.RecordKeySessionID: sid0, // imsi0->sid0
		},
	}

	// Default path
	sid, err := record.GetSessionID()
	assert.NoError(t, err)
	assert.Equal(t, sid0, sid)

	// IMSI-prefixed session ID should remove prefix
	record.Identifiers[directoryd.RecordKeySessionID] = sidWithIMSIPrefix
	sid, err = record.GetSessionID()
	assert.NoError(t, err)
	assert.Equal(t, sidWithoutPrefix, sid)

	// Err on non-string sid
	record.Identifiers[directoryd.RecordKeySessionID] = 42
	_, err = record.GetSessionID()
	assert.Error(t, err)

	// Empty string on no sid
	delete(record.Identifiers, directoryd.RecordKeySessionID)
	sid, err = record.GetSessionID()
	assert.NoError(t, err)
	assert.Equal(t, "", sid)
}

func TestDirectorydMethods(t *testing.T) {
	directoryd_test_init.StartTestService(t)

	// Empty initially
	_, err := directoryd.GetSessionIDForIMSI(nid0, imsi0)
	assert.Error(t, err)
	_, err = directoryd.GetHostnameForHWID(hwid0)
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
}

func TestDirectorydStateMethods(t *testing.T) {
	configurator_test_init.StartTestService(t)
	device_test_init.StartTestService(t)

	directoryd_test_init.StartTestService(t)
	state_test_init.StartTestService(t)

	stateClient, err := getStateServiceClient(t)
	assert.NoError(t, err)

	err = serde.RegisterSerdes(
		state.NewStateSerde(orc8r.DirectoryRecordType, &directoryd.DirectoryRecord{}),
		serde.NewBinarySerde(device.SerdeDomain, orc8r.AccessGatewayRecordType, &models.GatewayDevice{}),
	)
	assert.NoError(t, err)

	configurator_test_utils.RegisterNetwork(t, nid0, "DirectoryD Service Test")
	configurator_test_utils.RegisterGateway(t, nid0, hwid0, &models.GatewayDevice{HardwareID: hwid0})
	ctx := test_utils.GetContextWithCertificate(t, hwid0)

	record := &directoryd.DirectoryRecord{
		LocationHistory: []string{hwid0}, // imsi0->hwid0
		Identifiers: map[string]interface{}{
			directoryd.RecordKeySessionID: sid0, // imsi0->sid0
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
	err = state.DeleteStates(nid0, []state_types.ID{stateID})
	assert.NoError(t, err)

	// Get imsi0->hwid0, should be gone
	hwid, err = directoryd.GetHWIDForIMSI(nid0, imsi0)
	assert.Error(t, err)

	// Get imsi0->sid0, should be gone
	sid, err = directoryd.GetSessionIDForIMSI(nid0, imsi0)
	assert.Error(t, err)
}

func getStateServiceClient(t *testing.T) (protos.StateServiceClient, error) {
	conn, err := registry.GetConnection(state.ServiceName)
	assert.NoError(t, err)
	return protos.NewStateServiceClient(conn), err
}
