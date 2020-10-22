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
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/state"
)

var (
	// CWFStateSerdes contains the CWF-specific state serdes
	CWFStateSerdes = serde.NewRegistry(
		state.NewStateSerde(cwf.CwfSubscriberDirectoryType, &models.CwfSubscriberDirectoryRecord{}),
		state.NewStateSerde(cwf.CwfHAPairStatusType, &models.CarrierWifiHaPairStatus{}),
		state.NewStateSerde(cwf.CwfGatewayHealthType, &models.CarrierWifiGatewayHealthStatus{}),
	)

	// StateSerdes contains the full set of state serdes used in the LTE module
	StateSerdes = CWFStateSerdes.MustMerge(serdes.StateSerdes)
)
