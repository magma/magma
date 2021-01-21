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

package serdes

import (
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serde"
	ctraced_models "magma/orc8r/cloud/go/services/ctraced/obsidian/models"
	directoryd_types "magma/orc8r/cloud/go/services/directoryd/types"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/state"
)

const (
	deviceDomain = "device"
)

var (
	// Network contains the base orc8r serdes for configurator network configs
	Network = models.NetworkSerdes
	// Entity contains the base orc8r serdes for configurator network entities
	Entity = models.EntitySerdes.
		MustMerge(ctraced_models.EntitySerdes)
	// State contains the base orc8r serdes for the state service
	State = serde.NewRegistry(
		state.NewStateSerde(orc8r.GatewayStateType, &models.GatewayStatus{}),
		state.NewStateSerde(orc8r.StringMapSerdeType, &state.StringToStringMap{}),
		state.NewStateSerde(orc8r.DirectoryRecordType, &directoryd_types.DirectoryRecord{}),
	)
	// Device contains the base orc8r serdes for the device service
	Device = serde.NewRegistry(
		serde.NewBinarySerde(deviceDomain, orc8r.AccessGatewayRecordType, &models.GatewayDevice{}),
	)
)
