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

// Package test provides common definitions and function for eap related tests
package test

// Test Unit Data
type Data struct {
	EapIdentityResp,
	ExpectedChallengeReq,
	EapChallengeResp,

	RandAutn,
	Xres,
	ConfidentialityKey,
	IntegrityKey []byte

	IMSI,
	MSISDN string
}

const (
	// Test IMSI #1
	IMSI1 = "001010000000055"
	// Test IMSI #2
	IMSI2 = "001010000000043"
)

var (
	// Test Units Map
	Units = map[string]*Data{
		IMSI1: {
			EapIdentityResp: []byte("\x02\x01\x00\x40\x17\x05\x00\x00\x0e\x0e\x00\x33\x30\x30\x30\x31" +
				"\x30\x31\x30\x30\x30\x30\x30\x30\x30\x30\x35\x35\x40\x77\x6c\x61" +
				"\x6e\x2e\x6d\x6e\x63\x30\x30\x31\x2e\x6d\x63\x63\x30\x30\x31\x2e" +
				"\x33\x67\x70\x70\x6e\x65\x74\x77\x6f\x72\x6b\x2e\x6f\x72\x67\x00"),
			ExpectedChallengeReq: []byte{
				1, 2, 0, 68, 23, 1, 0, 0, 1, 5, 0, 0, 1, 35, 69, 103, 137, 171, 205, 239, 1, 35, 69, 103, 137,
				171, 205, 239, 2, 5, 0, 0, 84, 171, 100, 74, 144, 81, 185, 185, 94, 133, 193, 34, 62, 14, 241,
				76, 11, 5, 0, 0, 187, 28, 77, 175, 111, 216, 83, 74, 247, 124, 169, 254, 40, 141, 169, 189,
			},
			EapChallengeResp: []byte("\x02\x02\x00\x28\x17\x01\x00\x00\x0b\x05\x00\x00\xfd\x2b\x50\xbd" +
				"\x32\x24\x7a\xd7\x32\x9d\x9d\x26\x41\x60\x44\x46\x03\x03\x00\x40\x29\x5c\x00\xea\xe3\x88\x93\x0d"),
			RandAutn: []byte("\x01\x23\x45\x67\x89\xab\xcd\xef\x01\x23\x45\x67\x89\xab\xcd\xef" +
				"\x54\xab\x64\x4a\x90\x51\xb9\xb9\x5e\x85\xc1\x22\x3e\x0e\xf1\x4c"),
			Xres:               []byte("\x29\x5c\x00\xea\xe3\x88\x93\x0d"),
			ConfidentialityKey: []byte("\xa8\x35\xcf\x22\xb0\xf4\x3e\x15\x19\xd6\xfd\x23\x4c\x00\xd7\x93"),
			IntegrityKey:       []byte("\xd5\x37\x0f\x13\x79\x6f\x2f\x61\x5c\xbe\x15\xef\x9f\x42\x0a\x98"),
			IMSI:               IMSI1,
			MSISDN:             "123456789",
		},
		IMSI2: {
			EapIdentityResp: []byte("\x02\x02\x00\x40\x17\x05\x00\x00\x0e\x0e\x00\x33\x30\x30\x30\x31\x30" +
				"\x31\x30\x30\x30\x30\x30\x30\x30\x30\x34\x33\x40\x77\x6c\x61\x6e\x2e\x6d\x6e\x63\x30\x30\x31\x2e" +
				"\x6d\x63\x63\x30\x30\x31\x2e\x33\x67\x70\x70\x6e\x65\x74\x77\x6f\x72\x6b\x2e\x6f\x72\x67\x00"),
			ExpectedChallengeReq: []byte{
				1, 3, 0, 68, 23, 1, 0, 0, 1, 5, 0, 0, 254, 220, 186, 152, 118, 84, 50, 16, 254, 220, 186, 152, 118,
				84, 50, 16, 2, 5, 0, 0, 85, 108, 69, 100, 0, 217, 185, 185, 215, 177, 57, 81, 156, 159, 118, 136,
				11, 5, 0, 0, 9, 176, 57, 57, 175, 141, 130, 36, 60, 20, 41, 206, 233, 71, 100, 170,
			},
			EapChallengeResp: []byte("\x02\x03\x00\x28\x17\x01\x00\x00\x0b\x05\x00\x00\x10\xff\x67\x8d\x06" +
				"\xf2\x59\x09\x1b\x6f\x81\x9e\x5a\x62\x7a\x28\x03\x03\x00\x40\xe7\x17\xf3\x2f\x5d\xc8\xa9\x9b"),
			RandAutn: []byte("\xfe\xdc\xba\x98\x76\x54\x32\x10\xfe\xdc\xba\x98\x76\x54\x32" +
				"\x10\x55\x6c\x45\x64\x00\xd9\xb9\xb9\xd7\xb1\x39\x51\x9c\x9f\x76\x88"),
			Xres:               []byte("\xe7\x17\xf3\x2f\x5d\xc8\xa9\x9b"),
			ConfidentialityKey: []byte("\x21\xb1\x64\x48\x9b\xf5\x04\x7e\xae\x88\xc4\xcd\x7c\xcd\xe3\xc2"),
			IntegrityKey:       []byte("\xf4\xcb\x01\x9b\xed\xc8\x4d\x63\xc6\xce\xa7\xe2\xb0\x77\xfd\xb0"),
			IMSI:               IMSI2,
			MSISDN:             "456789012",
		},
	}

	// EAP Success Packet
	SuccessEAP = []byte{3, 2, 0, 4}
)
