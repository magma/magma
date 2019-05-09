/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 *  LICENSE file in the root directory of this source tree.
 */

package tests

import (
	"context"
	"encoding/json"
	"testing"

	"magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/service/middleware/unary/interceptors/tests"
	"magma/orc8r/cloud/go/services/magmad"
	magmad_protos "magma/orc8r/cloud/go/services/magmad/protos"
	magmad_test_init "magma/orc8r/cloud/go/services/magmad/test_init"
	"magma/orc8r/cloud/go/services/state"
	test_service "magma/orc8r/cloud/go/services/state/test_init"

	"github.com/golang/glog"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	typeName   = "typeName"
	testAgHwId = "Test-AGW-Hw-Id"
)

type stateBundle struct {
	value interface{}
	state *protos.State
	ID    state.StateID
}

func makeStateBundle(typeVal string, key string, value interface{}) stateBundle {
	marshaledValue, _ := json.Marshal(value)
	ID := state.StateID{Type: typeVal, DeviceID: key}
	state := protos.State{Type: typeVal, DeviceID: key, Value: marshaledValue}
	return stateBundle{state: &state, ID: ID, value: value}
}

func TestStateService(t *testing.T) {
	// Set up test networkID, hwID, and encode into context
	magmad_test_init.StartTestService(t)
	networkID, err := magmad.RegisterNetwork(
		&magmad_protos.MagmadNetworkRecord{Name: "State Service Test"},
		"state_service_test_network")
	hwId := protos.AccessGatewayID{Id: testAgHwId}
	magmad.RegisterGateway(
		networkID,
		&magmad_protos.AccessGatewayRecord{HwId: &hwId, Name: "Test GW Name"})
	csn := tests.StartMockGwAccessControl(t, []string{testAgHwId})
	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs(identity.CLIENT_CERT_SN_KEY, csn[0]))

	// Create States, IDs, values
	value0 := Name{Name: "name0"}
	value1 := Name{Name: "name1"}
	value2 := NameAndAge{Name: "name2", Age: 20}
	bundle0 := makeStateBundle(typeName, "key0", value0)
	bundle1 := makeStateBundle(typeName, "key1", value1)
	bundle2 := makeStateBundle(typeName, "key2", value2)

	test_service.StartTestService(t)
	err = serde.RegisterSerdes(&Serde{})
	assert.NoError(t, err)

	// Check contract for empty network
	//response, err := client.GetStates(ctx, makeGetStatesRequest(networkID, bundle0))
	states, err := state.GetStates(networkID, []state.StateID{bundle0.ID})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(states))

	// Report and read back
	err = reportStates(ctx, bundle0, bundle1)
	assert.NoError(t, err)
	states, err = state.GetStates(networkID, []state.StateID{bundle0.ID, bundle1.ID})
	assert.NoError(t, err)
	testGetStatesResponse(t, states, bundle0, bundle1)

	// Report a state with fields the corresponding serde does not expect
	err = reportStates(ctx, bundle2)
	assert.NoError(t, err)
	states, err = state.GetStates(networkID, []state.StateID{bundle2.ID})
	assert.NoError(t, err)
	testGetStatesResponse(t, states, bundle2)

	// Delete and read back
	err = state.DeleteStates(networkID, []state.StateID{bundle0.ID, bundle2.ID})
	assert.NoError(t, err)
	states, err = state.GetStates(networkID, []state.StateID{bundle0.ID, bundle1.ID, bundle2.ID})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(states))
	testGetStatesResponse(t, states, bundle1)
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

type Serde struct {
}

func (*Serde) GetDomain() string {
	return state.SerdeDomain
}

func (*Serde) GetType() string {
	return typeName
}

func (*Serde) Serialize(in interface{}) ([]byte, error) {
	return json.Marshal(in)

}

func (*Serde) Deserialize(message []byte) (interface{}, error) {
	res := Name{}
	err := json.Unmarshal(message, &res)
	return res, err
}

func getClient() (protos.StateServiceClient, *grpc.ClientConn, error) {
	conn, err := registry.GetConnection(state.ServiceName)
	if err != nil {
		initErr := errors.NewInitError(err, state.ServiceName)
		glog.Error(initErr)
		return nil, nil, initErr
	}
	return protos.NewStateServiceClient(conn), conn, err
}

func reportStates(ctx context.Context, bundles ...stateBundle) error {
	client, conn, err := getClient()
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = client.ReportStates(ctx, makeReportStatesRequest(bundles))
	return err
}

func testGetStatesResponse(t *testing.T, states map[state.StateID]state.StateValue, bundles ...stateBundle) {
	for _, bundle := range bundles {
		value := states[bundle.ID]
		assert.Equal(t, bundle.state.Value, value.ReportedValue)
	}
}

func makeReportStatesRequest(bundles []stateBundle) *protos.ReportStatesRequest {
	res := protos.ReportStatesRequest{}
	res.States = makeStates(bundles)
	return &res
}

func makeStates(bundles []stateBundle) []*protos.State {
	states := []*protos.State{}
	for _, bundle := range bundles {
		states = append(states, bundle.state)
	}
	return states
}
