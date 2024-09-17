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

// Package handlers provided AKA Response handlers for supported AKA subtypes
package aka

import (
	"crypto/hmac"
	"crypto/sha1"
	"math/big"

	"magma/feg/gateway/services/eap"
)

const (
	chunkSize = 64
	gSumSize  = 20
)

// GenMac calculates AKA MAC given data & K_auth (see: https://tools.ietf.org/html/rfc4187#section-10.15)
func GenMac(data, K_aut []byte) []byte {
	return HmacSha1(data, K_aut)[:16]
}

// HmacSha1 - SHA1 based HMAC
func HmacSha1(data, key []byte) []byte {
	h := hmac.New(sha1.New, key)
	h.Write(data)
	return h.Sum(nil)[:]
}

// MakeAKAKeys returns generated K_encr, K_aut, MSK, EMSK keys for AKA Authentication (RFC 4187, section 7)
func MakeAKAKeys(identity, IK, CK []byte) (K_encr, K_aut, MSK, EMSK []byte) {
	x := XSum(MK(identity, IK, CK))
	return x[:16], x[16:32], x[32:96], x[96:160]
}

// MK calculates & returns AKA Master Key: MK = SHA1(Identity|IK|CK)
func MK(identity, IK, CK []byte) []byte {
	d := sha1.New()
	d.Write(identity)
	d.Write(IK)
	d.Write(CK)
	return d.Sum(nil)
}

/* XSum generates 160 byte long byte slice of concatenated x_0..X_3 calculated according to RFC 4187, Appendix A
 *
 * let XKEY := MK,
 *
 * Step 3: For j = 0 to 3 do
 *   a. XVAL = XKEY
 *   b. w_0 = SHA1_Based_G(XVAL)
 *   c. XKEY = (1 + XKEY + w_0) mod 2^160
 *   d. XVAL = XKEY
 *   e. w_1 = SHA1_Based_G(XVAL)
 *   f. XKEY = (1 + XKEY + w_1) mod 2^160
 * 3.3 x_j = w_0|w_1
 */
func XSum(xkey []byte) []byte {
	x := make([]byte, 0, 160)
	for j := 0; j < 4; j++ {
		h := [5]uint32{0x67452301, 0xEFCDAB89, 0x98BADCFE, 0x10325476, 0xC3D2E1F0}
		w0 := GSum(&h, xkey)
		x = append(x, w0...) // x_j = w_0 | w_1 => x |= w_0
		xkey = updateXKey(xkey, w0)
		h = [5]uint32{0x67452301, 0xEFCDAB89, 0x98BADCFE, 0x10325476, 0xC3D2E1F0}
		w1 := GSum(&h, xkey)
		x = append(x, w1...) // x |= w_1
		xkey = updateXKey(xkey, w1)
	}
	if len(x) != 160 {
		panic(len(x))
	}
	// fmt.Printf("XSum (len=%d):\n%v\n", len(x), x)
	return x
}

var oneInt = big.NewInt(1)

// updateXKey: xkey = (1 + xkey + w_i) % 2^160
func updateXKey(xkey, wi []byte) []byte {
	var xkeyInt, wInt, xki *big.Int = big.NewInt(0), big.NewInt(0), big.NewInt(0)
	xki.Add(oneInt, xkeyInt.SetBytes(xkey))
	xki.Add(xki, wInt.SetBytes(wi))

	xkey = xki.Bytes()
	xlen := len(xkey)
	if xlen > 20 {
		xkey = xkey[xlen-20:]
	} else if xlen < 20 {
		xkey = append(make([]byte, 20-xlen), xkey...)
	}
	return xkey
}

// GSum SHA-1 based G function digest from FIPS Publication 186-2.
func GSum(h *[5]uint32, data []byte) []byte {
	var (
		digest [gSumSize]byte
		M0     = make([]byte, chunkSize)
	)
	copy(M0, data)
	block(h, M0)

	putUint32(digest[0:], h[0])
	putUint32(digest[4:], h[1])
	putUint32(digest[8:], h[2])
	putUint32(digest[12:], h[3])
	putUint32(digest[16:], h[4])
	return digest[:]
}

func putUint32(x []byte, s uint32) {
	_ = x[3]
	x[0] = byte(s >> 24)
	x[1] = byte(s >> 16)
	x[2] = byte(s >> 8)
	x[3] = byte(s)
}

// block is derived from std Go implementation of the SHA-1 block function
func block(h *[5]uint32, p []byte) {
	const (
		K0 = 0x5A827999
		K1 = 0x6ED9EBA1
		K2 = 0x8F1BBCDC
		K3 = 0xCA62C1D6
	)
	var w [16]uint32

	h0, h1, h2, h3, h4 := h[0], h[1], h[2], h[3], h[4]
	for len(p) >= chunkSize {
		for i := 0; i < 16; i++ {
			j := i * 4
			w[i] = uint32(p[j])<<24 | uint32(p[j+1])<<16 | uint32(p[j+2])<<8 | uint32(p[j+3])
		}
		a, b, c, d, e := h0, h1, h2, h3, h4

		i := 0
		for ; i < 16; i++ {
			f := b&c | (^b)&d
			a5 := a<<5 | a>>(32-5)
			b30 := b<<30 | b>>(32-30)
			t := a5 + f + e + w[i&0xf] + K0
			a, b, c, d, e = t, a, b30, c, d
		}
		for ; i < 20; i++ {
			tmp := w[(i-3)&0xf] ^ w[(i-8)&0xf] ^ w[(i-14)&0xf] ^ w[(i)&0xf]
			w[i&0xf] = tmp<<1 | tmp>>(32-1)

			f := b&c | (^b)&d
			a5 := a<<5 | a>>(32-5)
			b30 := b<<30 | b>>(32-30)
			t := a5 + f + e + w[i&0xf] + K0
			a, b, c, d, e = t, a, b30, c, d
		}
		for ; i < 40; i++ {
			tmp := w[(i-3)&0xf] ^ w[(i-8)&0xf] ^ w[(i-14)&0xf] ^ w[(i)&0xf]
			w[i&0xf] = tmp<<1 | tmp>>(32-1)
			f := b ^ c ^ d
			a5 := a<<5 | a>>(32-5)
			b30 := b<<30 | b>>(32-30)
			t := a5 + f + e + w[i&0xf] + K1
			a, b, c, d, e = t, a, b30, c, d
		}
		for ; i < 60; i++ {
			tmp := w[(i-3)&0xf] ^ w[(i-8)&0xf] ^ w[(i-14)&0xf] ^ w[(i)&0xf]
			w[i&0xf] = tmp<<1 | tmp>>(32-1)
			f := ((b | c) & d) | (b & c)

			a5 := a<<5 | a>>(32-5)
			b30 := b<<30 | b>>(32-30)
			t := a5 + f + e + w[i&0xf] + K2
			a, b, c, d, e = t, a, b30, c, d
		}
		for ; i < 80; i++ {
			tmp := w[(i-3)&0xf] ^ w[(i-8)&0xf] ^ w[(i-14)&0xf] ^ w[(i)&0xf]
			w[i&0xf] = tmp<<1 | tmp>>(32-1)
			f := b ^ c ^ d
			a5 := a<<5 | a>>(32-5)
			b30 := b<<30 | b>>(32-30)
			t := a5 + f + e + w[i&0xf] + K3
			a, b, c, d, e = t, a, b30, c, d
		}
		h0 += a
		h1 += b
		h2 += c
		h3 += d
		h4 += e

		p = p[chunkSize:]
	}
	h[0], h[1], h[2], h[3], h[4] = h0, h1, h2, h3, h4
}

// AppendMac appends AT_MAC attribute to eap npacket, signs the packet & returns the new, signed packet
// returns error if provided EAP Packet was malformed
func AppendMac(p eap.Packet, K_aut []byte) (eap.Packet, error) {
	p = p.Truncate()
	atMacOffset := len(p) + ATT_HDR_LEN
	p, err := p.Append(eap.NewAttribute(AT_MAC, append([]byte{0, 0}, make([]byte, MAC_LEN)...)))
	if err != nil {
		return p, err
	}
	mac := GenMac(p, K_aut)
	// Set AT_MAC
	copy(p[atMacOffset:], mac)
	return p, nil
}
