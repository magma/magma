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

package common

import (
	"crypto/md5"
	"errors"
	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2548"
	"math/rand"
)

// MicrosoftVendor bytes depicting Microsoft vendor for RADIUS Vendor-Specific attributes
var MicrosoftVendor = []byte{0x00, 0x00, 0x01, 0x37}

func generateMPPESalt(random int) []byte {
	return []byte{uint8((random&0xFF00)>>8) | 0x80, uint8(random & 0xFF)}
}

// GetKeyingAttributes Generates RADIUS keying materials
func GetKeyingAttributes(msk []byte, radiusSecret []byte, reqAuthenticator []byte) ([]radius.Attribute, error) {
	// Generate MS-MPPE-Recv-Key attribute
	recvKeyAttribute, err := GenerateMPPEAttribute(
		msk[:32],
		radiusSecret,
		reqAuthenticator,
		generateMPPESalt(rand.Int()),
		rfc2548.MSMPPERecvKey_Type,
	)
	if err != nil {
		return nil, err
	}

	// Generate MS-MPPE-Send-Key attribute
	sendKeyAttribute, err := GenerateMPPEAttribute(
		msk[32:],
		radiusSecret,
		reqAuthenticator,
		generateMPPESalt(rand.Int()),
		rfc2548.MSMPPESendKey_Type,
	)
	if err != nil {
		return nil, err
	}

	return []radius.Attribute{
		radius.Attribute(recvKeyAttribute),
		radius.Attribute(sendKeyAttribute),
	}, nil
}

// GenerateMPPEAttribute Wraps the MPPE key with RADIUS attributes as indicated
// in rfc2866 (Vendor-Specific attribute)
func GenerateMPPEAttribute(key []byte, s []byte, r []byte, a []byte, t radius.Type) ([]byte, error) {
	if len(a) != 2 {
		return nil, errors.New("salt must be exactly 2 bytes")
	}
	// Generate the cypher key C
	C := GenerateMPPEKey(key, s, r, a)

	// Wrap C (cypher key) with RADIUS Vendor Specific attribute
	attrBytes := append([]byte{
		byte(t),
		0x00, /* reserved for length */
		a[0],
		a[1],
	}, C...)
	attrBytes[1] = byte(len(attrBytes))
	attrBytes = append(MicrosoftVendor, attrBytes...)

	return attrBytes, nil
}

// GenerateMPPEKey follows RFC2548 section 2.4.3 and 2.4.4 to generate MPPE
// keying material. all parameter names are taken from the RFC so it's easier
// to follow the code
func GenerateMPPEKey(key []byte, s []byte, r []byte, a []byte) []byte {
	// Construct plaintext ()
	P := append([]byte{byte(len(key))}, key...)
	P = append(P, make([]byte, 0xF-(len(P)&0xF))...)

	// Calculate C
	C := []byte{}
	b := getMD5(s, r, a)
	for idx := 0; idx < len(P); idx += 16 {
		c := xor(P[idx:idx+16], b)
		C = append(C, c...)
		b = getMD5(s, c)
	}

	return C
}

// getMD5 calculate md5 of the given buffers, concatenated
func getMD5(bufs ...[]byte) []byte {
	hash := md5.New()
	for _, b := range bufs {
		hash.Write(b)
	}
	return hash.Sum(nil)
}

func xor(a []byte, b []byte) []byte {
	// Get the buffer with minimum length
	l := len(a)
	if len(b) < l {
		l = len(b)
	}

	// Iterate the buffer & xor each value
	res := make([]byte, l)
	for i := 0; i < l; i++ {
		res[i] = a[i] ^ b[i]
	}
	return res
}

func split(b []byte, chunkSize int) [][]byte {
	if len(b)%chunkSize != 0 {
		return nil
	}
	chunks := make([][]byte, 0, len(b)/chunkSize)
	for i := 0; i < len(b); i += chunkSize {
		chunks = append(chunks, b[i:i+chunkSize])
	}
	return chunks
}
