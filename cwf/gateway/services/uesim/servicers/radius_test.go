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

package servicers_test

import (
	"encoding/binary"
	"reflect"
	"testing"

	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2865"
	"fbc/lib/go/radius/rfc2869"
	"magma/cwf/gateway/services/uesim/servicers"
	"magma/feg/gateway/services/eap"

	"github.com/stretchr/testify/assert"
)

// Radius packets taken from cwf_2_aps.pcap.
const (
	RadiusAccessChallengeEapAkaIdentityRequestPacket = "\x0b\x08\x00\x34\xa0\x70\x9d\x6f\xe0\x8e\x92\x50\x99\x6c\x18\xa6" +
		"\x89\xa1\xa4\x65\x4f\x0e\x01\xe8\x00\x0c\x17\x05\x00\x00\x0a\x01" +
		"\x00\x00\x50\x12\x2c\x3b\x23\x8e\xc6\x11\x9d\xf5\xbb\x72\xb0\xdd" +
		"\x59\xd6\x2a\xbd"
	RadiusAccessRequestEapAkaIdentityResponsePacket = "\x01\x09\x00\xd6\x73\xea\x5e\xdf\x10\x25\x45\x3b\x21\x15\xdb\xc2" +
		"\xa9\x8a\x7c\x99\x4f\x42\x02\xe8\x00\x40\x17\x05\x00\x00\x0e\x0e\x00" +
		"\x33\x30\x30\x30\x31\x30\x31\x30\x30\x30\x30\x30\x30\x30\x30\x39" +
		"\x31\x40\x77\x6c\x61\x6e\x2e\x6d\x6e\x63\x30\x30\x31\x2e\x6d\x63" +
		"\x63\x30\x30\x31\x2e\x33\x67\x70\x70\x6e\x65\x74\x77\x6f\x72\x6b" +
		"\x2e\x6f\x72\x67\x00\x50\x12\x82\x1d\xf8\xb2\xab\xf1\x22\x1c\x2d" +
		"\x9c\x2e\xe7\xb1\xf3\xee\x2a\x1f\x13\x37\x36\x2d\x30\x32\x2d\x35" +
		"\x42\x2d\x38\x30\x2d\x45\x43\x2d\x34\x34\x1e\x27\x39\x38\x2d\x44" +
		"\x45\x2d\x44\x30\x2d\x38\x34\x2d\x42\x35\x2d\x34\x37\x3a\x43\x57" +
		"\x46\x2d\x54\x50\x2d\x4c\x49\x4e\x4b\x5f\x42\x35\x34\x37\x5f\x35" +
		"\x47\x01\x34\x30\x30\x31\x30\x31\x30\x30\x30\x30\x30\x30\x30\x30" +
		"\x39\x31\x40\x77\x6c\x61\x6e\x2e\x6d\x6e\x63\x30\x30\x31\x2e\x6d" +
		"\x63\x63\x30\x30\x31\x2e\x33\x67\x70\x70\x6e\x65\x74\x77\x6f\x72" +
		"\x6b\x2e\x6f\x72\x67"

	CalledStationID1 = "98-DE-D0-84-B5-47:CWF-TP-LINK_B547_5G"
	Secret           = "123456"
)

func TestRadius(t *testing.T) {
	server, _, err := setupTest(t)
	assert.NoError(t, err)

	p, err := radius.Parse([]byte(RadiusAccessChallengeEapAkaIdentityRequestPacket), []byte(Secret))
	assert.NoError(t, err)
	actual, err := server.HandleRadius(Imsi, CalledStationID1, p)
	assert.NoError(t, err)

	expected, err := radius.Parse([]byte(RadiusAccessRequestEapAkaIdentityResponsePacket), []byte(Secret))
	assert.NoError(t, err)
	assert.Equal(t, expected.Code, actual.Code)
	assert.Equal(t, expected.Authenticator, actual.Authenticator)
	assert.Equal(t, expected.Identifier, actual.Identifier)
	assert.Equal(t, expected.Get(rfc2865.UserName_Type), actual.Get(rfc2865.UserName_Type))
	assert.Equal(t, expected.Get(rfc2865.CallingStationID_Type), actual.Get(rfc2865.CallingStationID_Type))
	assert.Equal(t, expected.Get(rfc2865.CalledStationID_Type), actual.Get(rfc2865.CalledStationID_Type))
	assert.Equal(t, expected.Get(rfc2869.EAPMessage_Type), actual.Get(rfc2869.EAPMessage_Type))
}

func TestUserNotFound(t *testing.T) {
	server, _, err := setupTest(t)
	assert.NoError(t, err)

	p, err := radius.Parse([]byte(RadiusAccessChallengeEapAkaIdentityRequestPacket), []byte(Secret))
	assert.NoError(t, err)

	_, err = server.HandleRadius("012345678901234", CalledStationID1, p)
	assert.EqualError(t, err, "Error getting UE with specified IMSI: Not found")
}

func TestMissingEapPacket(t *testing.T) {
	server, _, err := setupTest(t)
	assert.NoError(t, err)

	p := radius.New(radius.CodeAccessChallenge, []byte(Secret))

	_, err = server.HandleRadius(Imsi, CalledStationID1, p)
	assert.EqualError(t, err, "Error extracting EAP message from Radius packet: no EAP-Message attribute found")
}

func TestEapToRadius(t *testing.T) {
	server, _, err := setupTest(t)
	assert.NoError(t, err)
	eapMessage := []byte(EapIdentityResponsePacket)
	expectedIdentifier := uint8(1)
	expectedImsi := "001010000000091"
	expectedCalledStationID := "98-DE-D0-84-B5-47:CWF-TP-LINK_B547_5G"
	radiusP, err := server.EapToRadius(eap.Packet(eapMessage), expectedImsi, expectedCalledStationID, expectedIdentifier)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(radius.CodeAccessRequest, radiusP.Code))
	assert.True(t, reflect.DeepEqual(expectedIdentifier, radiusP.Identifier))

	encoded, err := radiusP.Encode()
	assert.NoError(t, err)
	packetLen := binary.BigEndian.Uint16(encoded[2:4])
	assert.True(t, reflect.DeepEqual(uint16(len(encoded)), packetLen))

	eapAttr := []byte(radiusP.Attributes.Get(rfc2869.EAPMessage_Type))
	assert.True(t, reflect.DeepEqual(eapAttr, eapMessage))

	expectedUserName := expectedImsi + servicers.IdentityPostfix
	actualUserName := radiusP.Attributes.Get(rfc2865.UserName_Type)
	assert.Equal(t, expectedUserName, radius.String(actualUserName))
}
