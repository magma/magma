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
	"crypto/hmac"
	"crypto/md5"
	"encoding/binary"
	"fmt"

	"fbc/cwf/radius/modules/eap/packet"
	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2865"
	"fbc/lib/go/radius/rfc2866"
	"fbc/lib/go/radius/rfc2869"
	"magma/feg/gateway/services/eap"

	"github.com/pkg/errors"
)

// todo Replace constants with configurable fields
const (
	Auth = "\x73\xea\x5e\xdf\x10\x25\x45\x3b\x21\x15\xdb\xc2\xa9\x8a\x7c\x99"
)

// HandleRadius routes the Radius packet to the UE with the specified imsi.
func (srv *UESimServer) HandleRadius(imsi string, calledStationID string, p *radius.Packet) (*radius.Packet, error) {
	// todo Validate the packet. (Requires keeping state)

	// Extract EAP packet.
	eapMessage, err := packet.NewPacketFromRadius(p)
	if err != nil {
		err = errors.Wrap(err, "Error extracting EAP message from Radius packet")
		return nil, err
	}
	eapBytes, err := eapMessage.Bytes()
	if err != nil {
		err = errors.Wrap(err, "Error converting EAP packet to bytes")
		return nil, err
	}

	// Get the specified UE from the blobstore.
	ue, err := getUE(srv.store, imsi)
	if err != nil {
		return nil, err
	}

	// Generate EAP response.
	eapRes, err := srv.HandleEap(ue, eapBytes)
	if err != nil {
		return nil, err
	}

	// Wrap EAP response in Radius packet.
	res, err := srv.EapToRadius(eapRes, imsi, calledStationID, p.Identifier+1)
	if err != nil {
		return nil, err
	}

	return res, err
}

// EapToRadius puts an Eap packet payload in a Radius packet.
func (srv *UESimServer) EapToRadius(eapP eap.Packet, imsi string, calledStationID string, identifier uint8) (*radius.Packet, error) {
	radiusP := radius.New(radius.CodeAccessRequest, []byte(srv.cfg.radiusSecret))
	radiusP.Identifier = identifier

	// Hardcode in the auth.
	copy(radiusP.Authenticator[:], []byte(Auth)[:])
	radiusP.Attributes[rfc2865.UserName_Type] = []radius.Attribute{
		[]byte(imsi + IdentityPostfix),
	}
	// TODO: Fetch UE MAC addr and use as CallingStationID
	err := rfc2865.CallingStationID_SetString(radiusP, srv.cfg.brMac)
	if err != nil {
		return nil, err
	}
	err = rfc2865.CalledStationID_SetString(radiusP, calledStationID)
	if err != nil {
		return nil, err
	}
	err = rfc2869.EAPMessage_Set(radiusP, eapP)
	if err != nil {
		return nil, err
	}

	encoded, err := radiusP.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "Error encoding Radius packet")
	}

	// Add Message-Authenticator Attribute.
	encoded = srv.addMessageAuthenticator(encoded)

	// Parse to Radius packet.
	res, err := radius.Parse(encoded, []byte(srv.cfg.radiusSecret))
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	return res, nil
}

// MakeAccountStopRequest creates an Accounting Stop radius packet
func (srv *UESimServer) MakeAccountingStopRequest(calledStationID string) (*radius.Packet, error) {
	radiusP := radius.New(radius.CodeAccountingRequest, []byte(srv.cfg.radiusSecret))
	err := rfc2866.AcctStatusType_Set(radiusP, rfc2866.AcctStatusType_Value_Stop)
	if err != nil {
		return nil, err
	}
	err = rfc2865.CallingStationID_SetString(radiusP, srv.cfg.brMac)
	if err != nil {
		return nil, err
	}
	err = rfc2865.CalledStationID_SetString(radiusP, calledStationID)
	return radiusP, err
}

// addMessageAuthenticator calculates and adds the Message-Authenticator
// Attribute to a RADIUS packet.
func (srv *UESimServer) addMessageAuthenticator(encoded []byte) []byte {
	// Calculate new size
	size := uint16(len(encoded)) + radius.MessageAuthenticatorAttrLength
	binary.BigEndian.PutUint16(encoded[2:4], uint16(size))

	// Append the empty Message-Authenticator Attribute to the packet
	encoded = append(
		encoded,
		uint8(rfc2869.MessageAuthenticator_Type),
		uint8(radius.MessageAuthenticatorAttrLength),
	)
	encoded = append(encoded, make([]byte, 16)...)

	// Calculate Message-Authenticator and overwrite.
	hash := hmac.New(md5.New, []byte(srv.cfg.radiusSecret))
	hash.Write(encoded)
	encoded = hash.Sum(encoded[:len(encoded)-16])

	return encoded
}
