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

// package radius implements AAA server's radius interface for accounting & authentication
package radius

import (
	"crypto/hmac"
	"crypto/md5"
	"errors"
	"math/rand"

	"layeh.com/radius"
	"layeh.com/radius/rfc2869"
)

const (
	MSMPPESendKey_Type radius.Type = 16
	MSMPPERecvKey_Type radius.Type = 17
)

// MicrosoftVendor bytes depicting Microsoft vendor for RADIUS Vendor-Specific attributes
var MicrosoftVendor = []byte{0x00, 0x00, 0x01, 0x37}

// AddMessageAuthenticatorAttr generates and adds MessageAuthenticator AVP to the packet
func AddMessageAuthenticatorAttr(p *radius.Packet) error {
	if p == nil {
		return nil
	}
	zeroedAuthenticatorAttr := [16]byte{}
	authenticatorAttr := zeroedAuthenticatorAttr[:]
	p.Del(rfc2869.MessageAuthenticator_Type)
	p.Add(rfc2869.MessageAuthenticator_Type, authenticatorAttr) // add zeroed MA attr
	encoded, err := p.Encode()
	if err != nil {
		return err
	}
	// calculate MD5 hash for MessageAuthenticator
	hash := hmac.New(md5.New, p.Secret)
	hash.Write(encoded[:4])
	hash.Write(p.Authenticator[:])
	hash.Write(encoded[20:])
	encoded = hash.Sum(authenticatorAttr[0:0])

	p.Set(rfc2869.MessageAuthenticator_Type, encoded)
	return nil
}

// GetKeyingAttributes Generates RADIUS keying materials
func GetKeyingAttributes(msk []byte, secret []byte, authenticator []byte) (rcv, snd radius.Attribute, err error) {
	// Generate MS-MPPE-Recv-Key attribute
	rcv, err = GenerateMPPEAttribute(msk[:32], secret, authenticator, generateMPPESalt(rand.Int()), MSMPPERecvKey_Type)
	if err == nil {
		// Generate MS-MPPE-Send-Key attribute
		snd, err =
			GenerateMPPEAttribute(msk[32:], secret, authenticator, generateMPPESalt(rand.Int()), MSMPPESendKey_Type)
	}
	return
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

func generateMPPESalt(random int) []byte {
	return []byte{uint8((random&0xFF00)>>8) | 0x80, uint8(random & 0xFF)}
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
	if lb := len(b); lb < l {
		l = lb
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
