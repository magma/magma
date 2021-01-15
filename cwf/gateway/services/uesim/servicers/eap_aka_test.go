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
	"reflect"
	"testing"

	"magma/feg/gateway/services/eap"

	"github.com/stretchr/testify/assert"
)

const (
	EapAkaIdentityRequestPacket  = "\x01\xe8\x00\x0c\x17\x05\x00\x00\x0a\x01\x00\x00"
	EapAkaIdentityResponsePacket = "\x02\xe8\x00\x40\x17\x05\x00\x00\x0e\x0e\x00\x33\x30\x30\x30\x31" +
		"\x30\x31\x30\x30\x30\x30\x30\x30\x30\x30\x39\x31\x40\x77\x6c\x61" +
		"\x6e\x2e\x6d\x6e\x63\x30\x30\x31\x2e\x6d\x63\x63\x30\x30\x31\x2e" +
		"\x33\x67\x70\x70\x6e\x65\x74\x77\x6f\x72\x6b\x2e\x6f\x72\x67\x00"

	EapAkaChallengeRequestPacket = "\x01\xea\x00\x44\x17\x01\x00\x00\x01\x05\x00\x00\xee\xb3\x53\x6c" +
		"\x2f\xc3\x68\xfe\x3a\xfb\xd5\x5c\xfe\xf9\x6b\x29\x02\x05\x00\x00" +
		"\x94\x73\x37\x74\x82\xbd\x67\x41\x51\x11\x05\x57\x68\x17\xaa\x23" +
		"\x0b\x05\x00\x00\xda\x14\xa9\xce\x0e\x66\xaf\x38\x7b\x9f\xc1\xe6" +
		"\xf0\x31\x5e\x00"
	EapAkaChallengeResponsePacket = "\x02\xea\x00\x40\x17\x01\x00\x00\x03\x03\x00\x40\xdc\x89\x15\x16" +
		"\x8d\xd2\xeb\x56\x86\x06\x00\x00\x86\xe8\x20\x4d\xc6\xe1\xe3\xd8" +
		"\x94\x44\x3c\x26\xa7\xc6\x5d\xee\x3c\x42\xab\xf8\x0b\x05\x00\x00" +
		"\x13\x00\x7f\xe9\x86\xfc\xc1\x54\xf5\xca\x2b\xa7\x23\x88\x6d\x5b"
)

func TestEapAkaIdentityRequest(t *testing.T) {
	server, ue, err := setupTest(t)
	assert.NoError(t, err)

	res, err := server.HandleEap(ue, eap.Packet(EapAkaIdentityRequestPacket))
	assert.NoError(t, err)
	assert.True(
		t,
		reflect.DeepEqual([]byte(res), []byte(EapAkaIdentityResponsePacket)),
		"Actual packet didn't match expected packet\nexpected: %x\nactual:   %x\n",
		EapAkaIdentityResponsePacket,
		res,
	)
}

func TestEapAkaChallengeRequest(t *testing.T) {
	server, ue, err := setupTest(t)
	assert.NoError(t, err)

	res, err := server.HandleEap(ue, eap.Packet(EapAkaChallengeRequestPacket))
	assert.NoError(t, err)
	assert.True(
		t,
		reflect.DeepEqual([]byte(res), []byte(EapAkaChallengeResponsePacket)),
		"Actual packet didn't match expected packet\nexpected: %x\nactual:   %x\n",
		EapAkaChallengeResponsePacket,
		res,
	)
}
