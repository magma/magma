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
	lte_plugin "magma/lte/cloud/go/plugin"
	"magma/lte/cloud/go/serdes"
	lte_service "magma/lte/cloud/go/services/lte"
	"magma/lte/cloud/go/services/lte/obsidian/models"
	lte_test_init "magma/lte/cloud/go/services/lte/test_init"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/state/indexer"
	state_types "magma/orc8r/cloud/go/services/state/types"

	"github.com/go-openapi/swag"
	assert "github.com/stretchr/testify/require"
)

func TestIndexerEnodebState(t *testing.T) {
	const (
		version indexer.Version = 1 // copied from indexer_servicer.go
	)
	var (
		types = []string{lte.EnodebStateType} // copied from indexer_servicer.go
	)
	assert.NoError(t, plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{}))
	assert.NoError(t, plugin.RegisterPluginForTests(t, &lte_plugin.LteOrchestratorPlugin{}))

	lte_test_init.StartTestService(t)
	idx := indexer.NewRemoteIndexer(lte_service.ServiceName, version, types...)

	id1 := state_types.ID{Type: lte.EnodebStateType, DeviceID: "123"}
	id2 := state_types.ID{Type: lte.EnodebStateType, DeviceID: "123"}
	id3 := state_types.ID{Type: lte.MobilitydStateType, DeviceID: "555"}

	enbState1 := models.NewDefaultEnodebStatus()
	enbState2 := models.NewDefaultEnodebStatus()
	serializedState1 := serialize(t, enbState1)
	enbState2.MmeConnected = swag.Bool(false)
	serializedState2 := serialize(t, enbState2)
	stateGw1 := state_types.SerializedStatesByID{
		id1: {SerializedReportedState: serializedState1, ReporterID: "g1"},
	}
	stateGw2 := state_types.SerializedStatesByID{
		id2: {SerializedReportedState: serializedState2, ReporterID: "g2"},
	}

	clock.SetAndFreezeClock(t, time.Now())
	// Index the imsi0->sid0 state, result is sid0->imsi0 reverse mapping
	errs, err := idx.Index("nid0", stateGw1)
	assert.NoError(t, err)
	assert.Empty(t, errs)
	errs, err = idx.Index("nid0", stateGw2)
	assert.NoError(t, err)
	assert.Empty(t, errs)
	gotA, err := lte_service.GetEnodebState("nid0", "g1", "123")
	assert.NoError(t, err)
	assert.Equal(t, enbState1, gotA)
	gotB, err := lte_service.GetEnodebState("nid0", "g2", "123")
	assert.NoError(t, err)
	assert.Equal(t, enbState2, gotB)

	// Correctly handle per-state errs
	states := state_types.SerializedStatesByID{
		id1: {SerializedReportedState: serializedState2, ReporterID: "g1"},
		id3: {SerializedReportedState: serializedState2, ReporterID: "g3"},
	}
	errs, err = idx.Index("nid0", states)
	assert.NoError(t, err)
	assert.Error(t, errs[id3])
	gotC, err := lte_service.GetEnodebState("nid0", "g1", "123")
	assert.NoError(t, err)
	assert.Equal(t, enbState2, gotC)
}

func serialize(t *testing.T, enodebState *models.EnodebState) []byte {
	bytes, err := serde.Serialize(enodebState, lte.EnodebStateType, serdes.State)
	assert.NoError(t, err)
	return bytes
}
