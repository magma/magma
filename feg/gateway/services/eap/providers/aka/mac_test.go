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
package aka

import (
	"reflect"
	"testing"

	"magma/feg/gateway/services/eap"
)

const (
	testData = "\x01\x02\x00\xbc\x17\x01\x00\x00\x01\x05\x00\x00\x01\x23\x45\x67" +
		"\x89\xab\xcd\xef\x01\x23\x45\x67\x89\xab\xcd\xef\x02\x05\x00\x00" +
		"\x54\xab\x64\x4a\x90\x51\xb9\xb9\x5e\x85\xc1\x22\x3e\x0e\xf1\x4c" +
		"\x81\x05\x00\x00\xd1\xef\x2a\xdf\x8a\xf9\x74\xf1\xe2\x5f\xac\x28" +
		"\x58\xbc\xe4\x9e\x82\x19\x00\x00\x22\x24\x46\x74\xd6\x10\x1b\x1e" +
		"\xd7\xc8\xfa\x8d\x8c\x43\x87\x37\xd0\x49\x72\xac\x8a\x7a\x28\x64" +
		"\xb6\x39\x20\xb0\x7c\x25\xc4\xbf\xd4\x69\x2e\x88\xe2\x18\xd9\xd6" +
		"\xdf\x20\xe3\x05\x94\x5c\x25\x97\x23\xd4\x6a\x59\x5b\xf7\x1b\x25" +
		"\x2e\x8a\x47\xe1\x45\x0f\xb2\x3f\x40\xc1\x1b\x22\xeb\xf3\x69\x86" +
		"\xd6\x61\xb1\xa9\x98\xf1\xb8\x16\x50\xe6\x5c\x73\xd5\x66\xf1\xea" +
		"\x31\xd6\x68\x5d\x87\x36\x7d\xb4\x0b\x05\x00\x00" +
		// 128 bit MAC part
		"\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"

	expectedMac  = "\x55\x1f\xec\x03\xe0\xb1\xcc\x85\x31\x48\xb7\x5d\xf2\x57\x93\x65"
	IK           = "\xd5\x37\x0f\x13\x79\x6f\x2f\x61\x5c\xbe\x15\xef\x9f\x42\x0a\x98"
	CK           = "\xa8\x35\xcf\x22\xb0\xf4\x3e\x15\x19\xd6\xfd\x23\x4c\x00\xd7\x93"
	origIdentity = "\x30\x30\x30\x31\x30\x31\x30\x30\x30\x30\x30\x30" +
		"\x30\x30\x35\x35\x40\x77\x6c\x61\x6e\x2e\x6d\x6e\x63\x30\x30\x31" +
		"\x2e\x6d\x63\x63\x30\x30\x31\x2e\x33\x67\x70\x70\x6e\x65\x74\x77" +
		"\x6f\x72\x6b\x2e\x6f\x72\x67"
)

var (
	testHmac   = [20]byte{222, 124, 155, 133, 184, 183, 138, 166, 188, 138, 122, 54, 247, 10, 144, 112, 28, 157, 180, 217}
	ueTestData = []byte{2, 2, 0, 40, 23, 1, 0, 0, 11, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 0,
		64, 41, 92, 0, 234, 227, 136, 147, 13}
	ueMac                 = []byte{253, 43, 80, 189, 50, 36, 122, 215, 50, 157, 157, 38, 65, 96, 68, 70}
	MS_MPPE_Send_Key_Salt = []byte{0x9b, 0x87}
	MS_MPPE_Recv_Key_Salt = []byte{0x95, 0x63}
	authenticator         = []byte{
		0x9f, 0xe8, 0xff, 0xcb, 0xc9, 0xd4, 0x85, 0x97, 0xb9, 0x5b, 0x79, 0x7c, 0x2d, 0xf5, 0x43, 0x31}
	sharedSecret = []byte("1qaz2wsx")

	// Expected MS_MPPE_Send_Key
	MS_MPPE_Send_Key = []byte{0x9b, 0x87, 0x83, 0x49, 0x6a, 0x78, 0xcc, 0xaa, 0x34, 0x4e, 0x45, 0x51, 0x7f, 0x15,
		0x37, 0xf9, 0x30, 0x94, 0x26, 0x07, 0x60, 0x68, 0x97, 0xf0, 0xb5, 0x69, 0xab, 0x1d, 0x61, 0x9d, 0x8b, 0xa9,
		0x85, 0x3c, 0xc8, 0xaf, 0x68, 0x4b, 0xaa, 0x8f, 0x8f, 0x77, 0x5f, 0x68, 0x94, 0xf0, 0xcd, 0xc6, 0xc9, 0x2f}
	MS_MPPE_Recv_Key = []byte{0x95, 0x63, 0x3c, 0x3a, 0xa5, 0x8b, 0x48, 0xbe, 0xde, 0x6d, 0x2c, 0x1a, 0x91, 0x70,
		0x71, 0xf5, 0x63, 0xd4, 0xed, 0x7f, 0xba, 0xb3, 0xec, 0x61, 0xed, 0x7e, 0x3a, 0xf4, 0x82, 0x06, 0x58, 0x71,
		0x8c, 0xf7, 0xee, 0x86, 0x81, 0x0d, 0xf4, 0xf9, 0xf4, 0xb7, 0xb9, 0xdd, 0x14, 0xca, 0xc3, 0xbd, 0x95, 0x80}
)

func TestMacGeneration(t *testing.T) {
	t.Logf("Inputs:\n\tIdentity: %s\n\tIK: %v\n\tCK: %v", origIdentity, []byte(IK), []byte(CK))
	hmac := HmacSha1([]byte("The quick brown fox jumps over the lazy dog"), []byte("key"))
	t.Logf("Generated HMAC: %v", hmac)
	if !reflect.DeepEqual(hmac, testHmac[:]) {
		t.Fatalf(
			"HMACs don't match.\n\tGenerated HMAC(%d): %v\n\tExpected  HMAC(%d): %v",
			len(hmac), hmac, len(testHmac), testHmac)
	}

	K_encr, K_aut, MSK, EMSK := MakeAKAKeys([]byte(origIdentity), []byte(IK), []byte(CK))
	t.Logf("Generated keys:\n\tK_encr=%v\n\tK_aut=%v\n\tMSK=%v\n\tEMSK=%v", K_encr, K_aut, MSK, EMSK)

	t.Logf("Generated MS_MPPE_Recv_Key:\n\t%x",
		eap.EncodeMsMppeKey(MS_MPPE_Recv_Key_Salt, MSK[0:32], authenticator, sharedSecret))
	t.Logf("Generated MS_MPPE_Send_Key:\n\t%x",
		eap.EncodeMsMppeKey(MS_MPPE_Send_Key_Salt, MSK[32:64], authenticator, sharedSecret))

	if len(K_encr) != 16 {
		t.Fatalf("Invalid K_encr Len: %d", len(K_encr))
	}
	if len(K_aut) != 16 {
		t.Fatalf("Invalid K_aut Len: %d", len(K_aut))
	}
	if len(MSK) != 64 {
		t.Fatalf("Invalid MSK Len: %d", len(MSK))
	}
	if len(MSK) != 64 {
		t.Fatalf("Invalid EMSK Len: %d", len(EMSK))
	}

	mac := GenMac([]byte(testData), K_aut)

	if len(mac) != 16 {
		t.Fatalf("Invalid MAC Len: %d", len(mac))
	}

	// The moment of truth, compare generated MAC with expected
	if !reflect.DeepEqual(mac, []byte(expectedMac)) {
		t.Fatalf(
			"MACs don't match.\n\tGenerated MAC(%d): %v\n\tExpected  MAC(%d): %v",
			len(mac), mac, len(expectedMac), []byte(expectedMac))
	}

	mac = GenMac(ueTestData, K_aut)

	if !reflect.DeepEqual(mac, ueMac) {
		t.Fatalf(
			"MACs don't match.\n\tGenerated UE MAC(%d): %v\n\tExpected  UE MAC(%d): %v",
			len(mac), mac, len(ueMac), ueMac)
	}

	genMS_MPPE_Send_Key := append(
		MS_MPPE_Send_Key_Salt,
		eap.EncodeMsMppeKey(MS_MPPE_Send_Key_Salt, MSK[32:64], authenticator, sharedSecret)...)
	if !reflect.DeepEqual(genMS_MPPE_Send_Key, MS_MPPE_Send_Key) {
		t.Fatalf(
			"MS_MPPE_Send_Keys mismatch.\n\tGenerated MS_MPPE_Send_Key(%d): %v\n\tExpected  MS_MPPE_Send_Key(%d): %v",
			len(genMS_MPPE_Send_Key), genMS_MPPE_Send_Key, len(MS_MPPE_Send_Key), MS_MPPE_Send_Key)
	}

	genMS_MPPE_Recv_Key := append(
		MS_MPPE_Recv_Key_Salt,
		eap.EncodeMsMppeKey(MS_MPPE_Recv_Key_Salt, MSK[0:32], authenticator, sharedSecret)...)
	if !reflect.DeepEqual(genMS_MPPE_Recv_Key, MS_MPPE_Recv_Key) {
		t.Fatalf(
			"MS_MPPE_Recv_Keys mismatch.\n\tGenerated MS_MPPE_Recv_Key(%d): %v\n\tExpected  MS_MPPE_Recv_Key(%d): %v",
			len(genMS_MPPE_Recv_Key), genMS_MPPE_Recv_Key, len(MS_MPPE_Recv_Key), MS_MPPE_Recv_Key)
	}

}
