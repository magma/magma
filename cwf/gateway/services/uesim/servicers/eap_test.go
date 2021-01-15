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
	"context"
	"reflect"
	"testing"

	cwfprotos "magma/cwf/cloud/go/protos"
	"magma/cwf/gateway/services/uesim/servicers"
	fegprotos "magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap"
	"magma/orc8r/cloud/go/test_utils"

	"github.com/stretchr/testify/assert"
)

const (
	EapIdentityRequestPacket  = "\x01\xe7\x00\x05\x01"
	EapIdentityResponsePacket = "\x02\xe7\x00\x38\x01\x30\x30\x30\x31\x30\x31\x30\x30\x30\x30\x30" +
		"\x30\x30\x30\x39\x31\x40\x77\x6c\x61\x6e\x2e\x6d\x6e\x63\x30\x30" +
		"\x31\x2e\x6d\x63\x63\x30\x30\x31\x2e\x33\x67\x70\x70\x6e\x65\x74" +
		"\x77\x6f\x72\x6b\x2e\x6f\x72\x67"

	Imsi = "\x30\x30\x31\x30\x31\x30\x30\x30\x30\x30\x30\x30\x30\x39\x31"
	Key  = "\x8B\xAF\x47\x3F\x2F\x8F\xD0\x94\x87\xCC\xCB\xD7\x09\x7C\x68\x62"
	Opc  = "\x8e\x27\xb6\xaf\x0e\x69\x2e\x75\x0f\x32\x66\x7a\x3b\x14\x60\x5d"
	Seq  = 31
)

func setupTest(t *testing.T) (*servicers.UESimServer, *cwfprotos.UEConfig, error) {
	store := test_utils.NewSQLBlobstore(t, "useim_eap_test_blobstore")

	server, err := servicers.NewUESimServer(store)
	if err != nil {
		return server, &cwfprotos.UEConfig{}, err
	}

	ue := &cwfprotos.UEConfig{Imsi: Imsi, AuthKey: []byte(Key), AuthOpc: []byte(Opc), Seq: Seq}
	_, err = server.AddUE(context.Background(), ue)
	return server, ue, err
}

func TestEapIdentityRequest(t *testing.T) {
	server, ue, err := setupTest(t)
	assert.NoError(t, err)

	res, err := server.HandleEap(ue, eap.Packet(EapIdentityRequestPacket))
	assert.NoError(t, err)
	assert.True(
		t,
		reflect.DeepEqual([]byte(res), []byte(EapIdentityResponsePacket)),
		"Actual packet didn't match expected packet\nexpected: %x\nactual:   %x\n",
		EapIdentityResponsePacket,
		res,
	)
}

func TestInvalidEapPacket(t *testing.T) {
	server, ue, err := setupTest(t)
	assert.NoError(t, err)

	// Make packet and set its length to zero.
	badPacket := eap.NewPacket(eap.RequestCode, 1, []byte{})
	badPacket[eap.EapMsgLenHigh] = 0
	badPacket[eap.EapMsgLenLow] = 0

	_, err = server.HandleEap(ue, badPacket)
	assert.EqualError(t, err, "Error validating EAP packet: Invalid Packet Length: header => 0, actual => 4")
}

func TestUnsupportedEapType(t *testing.T) {
	server, ue, err := setupTest(t)
	assert.NoError(t, err)

	// Make packet and set its type to an unsupported type.
	badTypePacket := eap.NewPacket(eap.RequestCode, 1, []byte{uint8(fegprotos.EapType_Reserved)})

	_, err = server.HandleEap(ue, badTypePacket)
	assert.EqualError(t, err, "Unsupported Eap Type: 0")
}
