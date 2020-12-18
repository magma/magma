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

// Package servicers implements individual NH routed FeG services
package servicers

import (
	"strings"

	"magma/feg/cloud/go/services/feg_relay/gw_to_feg_relay"
	"magma/orc8r/cloud/go/services/dispatcher/gateway_registry"
)

const (
	DiamUnableToDeliverErr = 3002

	// FeG Relay Services
	FegS6aProxy     gateway_registry.GwServiceType = "s6a_proxy"
	FegSessionProxy gateway_registry.GwServiceType = "session_proxy"
	FegHello        gateway_registry.GwServiceType = "feg_hello"
	FegSwxProxy     gateway_registry.GwServiceType = "swx_proxy"
)

// RelayRouter implements generic routing logic and currently just embeds gw_to_feg_relay.Router functionality
type RelayRouter struct {
	gw_to_feg_relay.Router
}

// NewRelayRouter creates & returns a new RelayRouter
func NewRelayRouter() *RelayRouter {
	return &RelayRouter{Router: *gw_to_feg_relay.NewRouter()}
}

func getPlmnId6(imsi string) string {
	imsi = strings.TrimPrefix(strings.TrimSpace(imsi), "IMSI")
	if len(imsi) > gw_to_feg_relay.MaxPlmnIdLen {
		imsi = imsi[:gw_to_feg_relay.MaxPlmnIdLen]
	}
	return imsi
}
