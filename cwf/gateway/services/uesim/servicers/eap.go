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
	cwfprotos "magma/cwf/cloud/go/protos"
	fegprotos "magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap"

	"github.com/pkg/errors"
)

// todo Replace constants with configurable fields
const (
	IdentityPostfix = "@wlan.mnc001.mcc001.3gppnetwork.org"
)

// HandleEAP routes the EAP request to the UE with the specified imsi.
func (srv *UESimServer) HandleEap(ue *cwfprotos.UEConfig, req eap.Packet) (eap.Packet, error) {
	err := req.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "Error validating EAP packet")
	}

	switch fegprotos.EapType(req.Type()) {
	case fegprotos.EapType_Identity:
		return srv.eapIdentityRequest(ue, req)
	case fegprotos.EapType_AKA:
		return srv.handleEapAka(ue, req)
	}
	return nil, errors.Errorf("Unsupported Eap Type: %d", req[eap.EapMsgMethodType])
}

func (srv *UESimServer) eapIdentityRequest(ue *cwfprotos.UEConfig, req eap.Packet) (res eap.Packet, err error) {
	// Create the response EAP packet with the identity attribute.
	p := eap.NewPacket(
		eap.ResponseCode,
		req.Identifier(),
		append(
			[]byte{uint8(fegprotos.EapType_Identity)},
			[]byte("\x30"+ue.GetImsi()+IdentityPostfix)...,
		),
	)

	return p, nil
}
