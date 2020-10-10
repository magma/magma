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

package models

import (
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/wifi/cloud/go/wifi"
)

var (
	// NetworkSerdes contains the package's configurator network config serdes
	NetworkSerdes = serde.NewRegistry(
		configurator.NewNetworkConfigSerde(wifi.WifiNetworkType, &NetworkWifiConfigs{}),
	)
	// EntitySerdes contains the package's configurator network entity serdes
	EntitySerdes = serde.NewRegistry(
		configurator.NewNetworkEntityConfigSerde(wifi.WifiGatewayType, &GatewayWifiConfigs{}),
		configurator.NewNetworkEntityConfigSerde(wifi.MeshEntityType, &MeshWifiConfigs{}),
	)
)
