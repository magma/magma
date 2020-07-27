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

package authstate

import (
	"fbc/cwf/radius/modules/eap/packet"

	"fbc/lib/go/radius"
)

// Container A storable container for protocol state
type Container struct {
	LogCorrelationID uint64         `json:"correlation_id"` // Request correlation id
	EapType          packet.EAPType `json:"eap_type"`       // EAP type of the auth session
	ProtocolState    string         `json:"protocol_state"` // EAP-* Protocol-specific state
	RadiusSessionID  *string        `json:"session_id"`     // RADIUS Session ID
}

// Manager an interface for EAP state management storage
type Manager interface {
	Set(authReq *radius.Packet, eaptype packet.EAPType, state Container) error
	Get(authReq *radius.Packet, eaptype packet.EAPType) (*Container, error)
	Reset(authReq *radius.Packet, eapType packet.EAPType) error
}
