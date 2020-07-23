/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package servicers

import (
	"fbc/lib/go/radius"
	"magma/feg/gateway/services/eap"
)

// todo use a config to assign this value
const (
	EapIdentityRequestPacket = "\x01\x00\x00\x05\x01"
)

// CreateEAPIdentityRequest simulates starting the EAP-AKA authentication by sending a UE an
// EAP Identity Request packet.
func (srv *UESimServer) CreateEAPIdentityRequest(imsi, calledStationID string) (*radius.Packet, error) {
	ue, err := getUE(srv.store, imsi)
	if err != nil {
		return nil, err
	}

	eapReponse, err := srv.HandleEap(ue, eap.Packet(EapIdentityRequestPacket))
	if err != nil {
		return nil, err
	}

	// Set packet Identifier to 0.
	return srv.EapToRadius(eapReponse, imsi, calledStationID, 0)
}
