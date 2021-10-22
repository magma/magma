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
	"magma/orc8r/cloud/go/serdes"
	"magma/wifi/cloud/go/services/wifi/obsidian/models"
)

var (
	// Network contains the full set of configurator network config serdes
	// used in the wifi module
	Network = serdes.Network.MustMerge(models.NetworkSerdes)
	// Entity contains the full set of configurator network entity serdes used
	// in the wifi module
	Entity = serdes.Entity.MustMerge(models.EntitySerdes)
	// Device contains the full set of device serdes used in the wifi module
	Device = serdes.Device
)
