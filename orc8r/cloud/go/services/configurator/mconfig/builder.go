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

package mconfig

import (
	"magma/orc8r/cloud/go/services/configurator/storage"
)

type ConfigsByKey map[string][]byte

// Builder creates a partial mconfig for a gateway within a network.
type Builder interface {
	// Build returns a partial mconfig containing the gateway configs for which
	// this builder is responsible.
	//
	// Parameters:
	//	- network	-- network containing the gateway
	//	- graph		-- entity graph associated with the gateway
	//	- gatewayID	-- HWID of the gateway
	Build(network *storage.Network, graph *storage.EntityGraph, gatewayID string) (ConfigsByKey, error)
}
