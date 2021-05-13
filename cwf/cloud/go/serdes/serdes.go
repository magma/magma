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
	"magma/cwf/cloud/go/cwf"
	"magma/cwf/cloud/go/services/cwf/obsidian/models"
	feg_models "magma/feg/cloud/go/services/feg/obsidian/models"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	policydb_models "magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/state"
)

var (
	// Network contains the full set of configurator network config serdes
	// used in the CWF module
	Network = serdes.Network.
		MustMerge(models.NetworkSerdes).
		// FeG serdes
		MustMerge(feg_models.NetworkSerdes).
		// LTE serdes
		MustMerge(lte_models.NetworkSerdes).
		MustMerge(policydb_models.NetworkSerdes)
	// Entity contains the full set of configurator network entity serdes used
	// in the CWF module
	Entity = serdes.Entity.
		MustMerge(models.EntitySerdes).
		MustMerge(feg_models.EntitySerdes)
	// State contains the full set of state serdes used in the CWF module
	State = serdes.State.MustMerge(serde.NewRegistry(
		state.NewStateSerde(cwf.CwfSubscriberDirectoryType, &models.CwfSubscriberDirectoryRecord{}),
		state.NewStateSerde(cwf.CwfHAPairStatusType, &models.CarrierWifiHaPairStatus{}),
		state.NewStateSerde(cwf.CwfGatewayHealthType, &models.CarrierWifiGatewayHealthStatus{}),
	))
	// Device contains the full set of device serdes used in the CWF module
	Device = serdes.Device
)
