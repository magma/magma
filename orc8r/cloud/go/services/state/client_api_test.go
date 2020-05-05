/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package state_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/serde"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	configurator_test "magma/orc8r/cloud/go/services/configurator/test_utils"
	"magma/orc8r/cloud/go/services/device"
	device_test_init "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/state"
	state_test_init "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/services/state/test_utils"
	"magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
	"github.com/stretchr/testify/assert"
)

const testAgHwId = "Test-AGW-Hw-Id"

func TestStateService(t *testing.T) {
	configurator_test_init.StartTestService(t)
	device_test_init.StartTestService(t)
	// Set up test networkID, hwID, and encode into context
	state_test_init.StartTestService(t)
	err := serde.RegisterSerdes(
		state.NewStateSerde("test-serde", &Name{}),
		serde.NewBinarySerde(device.SerdeDomain, orc8r.AccessGatewayRecordType, &models.GatewayDevice{}))
	assert.NoError(t, err)

	networkID := "state_service_test_network"
	configurator_test.RegisterNetwork(t, networkID, "State Service Test")
	gatewayID := testAgHwId
	configurator_test.RegisterGateway(t, networkID, gatewayID, &models.GatewayDevice{HardwareID: testAgHwId})
	ctx := test_utils.GetContextWithCertificate(t, testAgHwId)

	// Create States, IDs, values
	value0 := Name{Name: "name0"}
	value1 := Name{Name: "name1"}
	value2 := NameAndAge{Name: "name2", Age: 20}
	bundle0 := makeStateBundle("test-serde", "key0", value0)
	bundle1 := makeVersionedStateBundle("test-serde", "key1", value1, 10)
	bundle2 := makeVersionedStateBundle("test-serde", "key2", value2, 12)

	// Check contract for empty network
	states, err := state.GetStates(networkID, []state.ID{bundle0.ID})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(states))

	// Report and read back
	_, err = reportStates(ctx, bundle0, bundle1)
	assert.NoError(t, err)
	states, err = state.GetStates(networkID, []state.ID{bundle0.ID, bundle1.ID})
	assert.NoError(t, err)
	testGetStatesResponse(t, states, bundle0, bundle1)
	assert.Equal(t, uint64(0), states[bundle0.ID].Version)
	assert.Equal(t, uint64(10), states[bundle1.ID].Version)

	// Update states, ensuring version is set properly
	bundle1.state.Version = 15
	_, err = reportStates(ctx, bundle0, bundle1)
	assert.NoError(t, err)
	states, err = state.GetStates(networkID, []state.ID{bundle0.ID, bundle1.ID})
	assert.NoError(t, err)
	testGetStatesResponse(t, states, bundle0, bundle1)
	assert.Equal(t, uint64(1), states[bundle0.ID].Version)
	assert.Equal(t, uint64(15), states[bundle1.ID].Version)

	// Sync states
	bundle0.state.Version = 1  // synced
	bundle1.state.Version = 20 // unsynced
	res, err := syncStates(ctx, bundle0, bundle1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res.GetUnsyncedStates()))
	assert.Equal(t, bundle1.ID.DeviceID, res.GetUnsyncedStates()[0].Id.DeviceID)
	assert.Equal(t, bundle1.ID.Type, res.GetUnsyncedStates()[0].Id.Type)
	assert.Equal(t, uint64(15), res.GetUnsyncedStates()[0].Version)

	// Report a state with fields the corresponding serde does not expect
	_, err = reportStates(ctx, bundle2)
	assert.NoError(t, err)
	states, err = state.GetStates(networkID, []state.ID{bundle2.ID})
	assert.NoError(t, err)
	testGetStatesResponse(t, states, bundle2)
	assert.Equal(t, uint64(12), states[bundle2.ID].Version)

	// Delete and read back
	err = state.DeleteStates(networkID, []state.ID{bundle0.ID, bundle2.ID})
	assert.NoError(t, err)
	states, err = state.GetStates(networkID, []state.ID{bundle0.ID, bundle1.ID, bundle2.ID})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(states))
	testGetStatesResponse(t, states, bundle1)

	// Send a valid state and a state with no corresponding serde
	unserializableBundle := makeStateBundle("nonexistent-serde", "key3", value0)
	invalidBundle := makeStateBundle("test-serde", "key1", Name{Name: "BADNAME"})
	resp, err := reportStates(ctx, bundle0, unserializableBundle, invalidBundle)
	assert.NoError(t, err)
	assert.Equal(t, "nonexistent-serde", resp.UnreportedStates[0].Type)
	assert.Equal(t, "No Serde found for type nonexistent-serde", resp.UnreportedStates[0].Error)
	assert.Equal(t, "test-serde", resp.UnreportedStates[1].Type)
	assert.Equal(t, "this name: BADNAME is not allowed", resp.UnreportedStates[1].Error)
	// Valid state should still be reported
	states, err = state.GetStates(networkID, []state.ID{bundle0.ID, bundle1.ID, bundle2.ID})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(states))
	testGetStatesResponse(t, states, bundle0, bundle1)
}

type stateBundle struct {
	state *protos.State
	ID    state.ID
}

func makeVersionedStateBundle(typeVal, key string, value interface{}, version uint64) stateBundle {
	stateBundle := makeStateBundle(typeVal, key, value)
	stateBundle.state.Version = version
	return stateBundle
}

func makeStateBundle(typeVal, key string, value interface{}) stateBundle {
	marshaledValue, _ := json.Marshal(value)
	ID := state.ID{Type: typeVal, DeviceID: key}
	st := protos.State{Type: typeVal, DeviceID: key, Value: marshaledValue}
	return stateBundle{state: &st, ID: ID}
}

type NameAndAge struct {
	// name
	Name string `json:"name"`
	// age
	Age int `json:"age"`
}

type Name struct {
	// name
	Name string `json:"name"`
}

func (*Name) GetDomain() string {
	return state.SerdeDomain
}

func (*Name) GetType() string {
	return "test-serde"
}

func (m *Name) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)

}

func (m *Name) UnmarshalBinary(message []byte) error {
	res := Name{}
	err := json.Unmarshal(message, &res)
	*m = res
	return err
}

func (m *Name) ValidateModel() error {
	if m.Name == "BADNAME" {
		return fmt.Errorf("this name: %s is not allowed", m.Name)
	}
	return nil
}

func getClient() (protos.StateServiceClient, error) {
	conn, err := registry.GetConnection(state.ServiceName)
	if err != nil {
		initErr := errors.NewInitError(err, state.ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewStateServiceClient(conn), err
}

func reportStates(ctx context.Context, bundles ...stateBundle) (*protos.ReportStatesResponse, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}
	response, err := client.ReportStates(ctx, makeReportStatesRequest(bundles))
	return response, err
}

func syncStates(ctx context.Context, bundles ...stateBundle) (*protos.SyncStatesResponse, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}
	response, err := client.SyncStates(ctx, makeSyncStatesRequest(bundles))
	return response, err
}

func testGetStatesResponse(t *testing.T, states map[state.ID]state.State, bundles ...stateBundle) {
	for _, bundle := range bundles {
		value := states[bundle.ID]
		iState, err := serde.Deserialize(state.SerdeDomain, bundle.ID.Type, bundle.state.Value)
		assert.NoError(t, err)
		assert.Equal(t, iState, value.ReportedState)
	}
}

func makeReportStatesRequest(bundles []stateBundle) *protos.ReportStatesRequest {
	res := protos.ReportStatesRequest{}
	res.States = makeStates(bundles)
	return &res
}

func makeSyncStatesRequest(bundles []stateBundle) *protos.SyncStatesRequest {
	res := protos.SyncStatesRequest{}
	var states []*protos.IDAndVersion
	for _, bundle := range bundles {
		st := &protos.IDAndVersion{
			Id: &protos.StateID{
				Type:     bundle.ID.Type,
				DeviceID: bundle.ID.DeviceID,
			},
			Version: bundle.state.Version,
		}
		states = append(states, st)
	}
	res.States = states
	return &res
}

func makeStates(bundles []stateBundle) []*protos.State {
	var states []*protos.State
	for _, bundle := range bundles {
		states = append(states, bundle.state)
	}
	return states
}
