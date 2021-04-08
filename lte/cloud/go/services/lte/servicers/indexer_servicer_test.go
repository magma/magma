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
	"time"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	lte_service "magma/lte/cloud/go/services/lte"
	"magma/lte/cloud/go/services/lte/obsidian/models"
	lte_test_init "magma/lte/cloud/go/services/lte/test_init"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	models2 "magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/state/indexer"
	state_types "magma/orc8r/cloud/go/services/state/types"
	"magma/orc8r/cloud/go/storage"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

func TestIndexerEnodebState(t *testing.T) {
	const (
		version indexer.Version = 1 // copied from indexer_servicer.go
	)
	var (
		types = []string{lte.EnodebStateType} // copied from indexer_servicer.go
	)
	configurator_test_init.StartTestService(t)
	lte_test_init.StartTestService(t)
	idx := indexer.NewRemoteIndexer(lte_service.ServiceName, version, types...)

	id1 := state_types.ID{Type: lte.EnodebStateType, DeviceID: "123"}
	id2 := state_types.ID{Type: lte.EnodebStateType, DeviceID: "123"}
	id3 := state_types.ID{Type: lte.MobilitydStateType, DeviceID: "555"}

	networkID := "nid0"
	gatewayID1 := "g1"
	hwID1 := "hw1"
	gatewayID2 := "g2"
	hwID2 := "hw2"
	enbSN := "123"

	// Setup gw ents to be fetched during indexing
	seedNetwork(t, networkID)
	seedTier(t, networkID)
	seedGateway(t, networkID, gatewayID1, hwID1)
	seedGateway(t, networkID, gatewayID2, hwID2)

	enbState1 := models.NewDefaultEnodebStatus()
	enbState2 := models.NewDefaultEnodebStatus()
	serializedState1 := serialize(t, enbState1)
	enbState2.MmeConnected = swag.Bool(false)
	serializedState2 := serialize(t, enbState2)
	stateGw1 := state_types.SerializedStatesByID{
		id1: {SerializedReportedState: serializedState1, ReporterID: hwID1},
	}
	stateGw2 := state_types.SerializedStatesByID{
		id2: {SerializedReportedState: serializedState2, ReporterID: hwID2},
	}

	clock.SetAndFreezeClock(t, time.Now())
	// Index the imsi0->sid0 state, result is sid0->imsi0 reverse mapping
	errs, err := idx.Index(networkID, stateGw1)
	assert.NoError(t, err)
	assert.Empty(t, errs)
	errs, err = idx.Index(networkID, stateGw2)
	assert.NoError(t, err)
	assert.Empty(t, errs)
	gotA, err := lte_service.GetEnodebState(networkID, gatewayID1, enbSN)
	assert.NoError(t, err)
	assert.Equal(t, enbState1, gotA)
	gotB, err := lte_service.GetEnodebState(networkID, gatewayID2, enbSN)
	assert.NoError(t, err)
	assert.Equal(t, enbState2, gotB)

	// Correctly handle per-state errs
	states := state_types.SerializedStatesByID{
		id1: {SerializedReportedState: serializedState2, ReporterID: hwID1},
		id3: {SerializedReportedState: serializedState2, ReporterID: "hw3"},
	}
	errs, err = idx.Index(networkID, states)
	assert.NoError(t, err)
	assert.Error(t, errs[id3])
	gotC, err := lte_service.GetEnodebState(networkID, gatewayID1, enbSN)
	assert.NoError(t, err)
	assert.Equal(t, enbState2, gotC)
}

func seedNetwork(t *testing.T, networkID string) {
	err := configurator.CreateNetwork(configurator.Network{ID: networkID}, serdes.Network)
	assert.NoError(t, err)
}

func seedGateway(t *testing.T, networkID string, gatewayID string, hwID string) {
	_, err := configurator.CreateEntity(
		networkID,
		configurator.NetworkEntity{
			Type:         orc8r.MagmadGatewayType,
			Key:          gatewayID,
			Config:       &models2.MagmadGatewayConfigs{},
			PhysicalID:   hwID,
			Associations: []storage.TypeAndKey{{Type: orc8r.UpgradeTierEntityType, Key: "t0"}},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)
}

func seedTier(t *testing.T, networkID string) {
	// setup fixtures in backend
	_, err := configurator.CreateEntities(
		networkID,
		[]configurator.NetworkEntity{
			{Type: orc8r.UpgradeTierEntityType, Key: "t0"},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)
}
func serialize(t *testing.T, enodebState *models.EnodebState) []byte {
	bytes, err := serde.Serialize(enodebState, lte.EnodebStateType, serdes.State)
	assert.NoError(t, err)
	return bytes
}
