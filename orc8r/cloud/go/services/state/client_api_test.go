/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package state_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"magma/orc8r/cloud/go/serde"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	configurator_test "magma/orc8r/cloud/go/services/configurator/test_utils"
	device_test_init "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/state"
	state_test_init "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/services/state/test_utils"
	state_types "magma/orc8r/cloud/go/services/state/types"
	"magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
)

const (
	testAgHwId = "Test-AGW-Hw-Id"
)

var (
	stateSerdes = serde.NewRegistry(
		state.NewStateSerde("test-serde", &Name{}),
	)
)

func init() {
	//_ = flag.Set("alsologtostderr", "true") // uncomment to view logs during test
}

func TestStateService(t *testing.T) {
	configurator_test_init.StartTestService(t)
	device_test_init.StartTestService(t)
	// Set up test networkID, hwID, and encode into context
	state_test_init.StartTestService(t)

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
	states, err := state.GetStates(networkID, state_types.IDs{bundle0.ID}, stateSerdes)
	assert.NoError(t, err)
	assert.Empty(t, states)

	// Report and read back
	repRes, err := reportStates(ctx, bundle0, bundle1)
	assert.NoError(t, err)
	assert.Empty(t, repRes.UnreportedStates)
	states, err = state.GetStates(networkID, state_types.IDs{bundle0.ID, bundle1.ID}, stateSerdes)
	assert.NoError(t, err)
	testGetStatesResponse(t, states, bundle0, bundle1)
	assert.Equal(t, uint64(0), states[bundle0.ID].Version)
	assert.Equal(t, uint64(10), states[bundle1.ID].Version)

	// Update states, ensuring version is set properly
	bundle1.state.Version = 15
	repRes, err = reportStates(ctx, bundle0, bundle1)
	assert.NoError(t, err)
	assert.Empty(t, repRes.UnreportedStates)
	states, err = state.GetStates(networkID, state_types.IDs{bundle0.ID, bundle1.ID}, stateSerdes)
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
	repRes, err = reportStates(ctx, bundle2)
	assert.NoError(t, err)
	assert.Empty(t, repRes.UnreportedStates)
	states, err = state.GetStates(networkID, state_types.IDs{bundle2.ID}, stateSerdes)
	assert.NoError(t, err)
	testGetStatesResponse(t, states, bundle2)
	assert.Equal(t, uint64(12), states[bundle2.ID].Version)

	// Delete and read back
	err = state.DeleteStates(networkID, state_types.IDs{bundle0.ID, bundle2.ID})
	assert.NoError(t, err)
	states, err = state.GetStates(networkID, state_types.IDs{bundle0.ID, bundle1.ID, bundle2.ID}, stateSerdes)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(states))
	testGetStatesResponse(t, states, bundle1)

	// Send a valid and invalid state
	// key0: added here
	// key1: overwritten with invalid
	// key2: added with no serde
	// => should only receive key0
	unserializableBundle := makeStateBundle("nonexistent-serde", "key3", value0)
	invalidBundle := makeStateBundle("test-serde", "key1", Name{Name: "BADNAME"})
	repRes, err = reportStates(ctx, bundle0, unserializableBundle, invalidBundle)
	assert.NoError(t, err)
	assert.Empty(t, repRes.UnreportedStates) // validity is checked by the consumer
	// Only valid state should be accessible
	states, err = state.GetStates(networkID, state_types.IDs{bundle0.ID, bundle1.ID, bundle2.ID, unserializableBundle.ID}, stateSerdes)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(states))
	testGetStatesResponse(t, states, bundle0)
}

type stateBundle struct {
	state *protos.State
	ID    state_types.ID
}

func makeVersionedStateBundle(typeVal, key string, value interface{}, version uint64) stateBundle {
	stateBundle := makeStateBundle(typeVal, key, value)
	stateBundle.state.Version = version
	return stateBundle
}

func makeStateBundle(typeVal, key string, value interface{}) stateBundle {
	marshaledValue, _ := json.Marshal(value)
	ID := state_types.ID{Type: typeVal, DeviceID: key}
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

func reportStates(ctx context.Context, bundles ...stateBundle) (*protos.ReportStatesResponse, error) {
	client, err := state.GetStateClient()
	if err != nil {
		return nil, err
	}
	response, err := client.ReportStates(ctx, makeReportStatesRequest(bundles))
	return response, err
}

func syncStates(ctx context.Context, bundles ...stateBundle) (*protos.SyncStatesResponse, error) {
	client, err := state.GetStateClient()
	if err != nil {
		return nil, err
	}
	response, err := client.SyncStates(ctx, makeSyncStatesRequest(bundles))
	return response, err
}

func testGetStatesResponse(t *testing.T, states map[state_types.ID]state_types.State, bundles ...stateBundle) {
	for _, bundle := range bundles {
		value := states[bundle.ID]
		iState, err := serde.Deserialize(bundle.state.Value, bundle.ID.Type, stateSerdes)
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
