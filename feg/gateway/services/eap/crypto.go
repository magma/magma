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

package eap

import "crypto/md5"

// EncodeMsMppeKey implements RFC 2048 encoding for MS-MPPE-Send-Key & MS-MPPE-Recv-Key radius keys,
// returns padded & encoded key.
// See: https://www.ietf.org/rfc/rfc2548.txt 2.4.2, 2.4.3
func EncodeMsMppeKey(salt, key, authenticatorKey, sharedSecret []byte) []byte {
	l := len(key) + 1
	reminder := l % 16
	if reminder != 0 {
		l += 16 - reminder
	}
	p := make([]byte, l)
	p[0] = byte(len(key))
	copy(p[1:], key)
	// b(1) = MD5(S + R + A)    c(1) = p(1) xor b(1)   C = c(1)
	hash := md5.New()
	hash.Write(sharedSecret)
	hash.Write(authenticatorKey)
	hash.Write(salt)
	b := hash.Sum(nil)
	c := XorBytes(p[:16], b)
	C := c
	for pstart := 16; pstart < l; pstart += 16 {
		// b(i) = MD5(S + c(i-1))   c(i) = p(i) xor b(i)   C = C + c(i)
		hash.Reset()
		hash.Write(sharedSecret)
		hash.Write(c)
		b = hash.Sum(nil)
		c = XorBytes(p[pstart:pstart+16], b)
		C = append(C, c...)
	}
	return C
}

// XorBytes returns a new byte slice c where c[i] = a[i] ^ b[i]
func XorBytes(a, b []byte) []byte {
	l := len(a)
	if len(b) < l {
		l = len(b)
	}
	c := make([]byte, l)
	for i := 0; i < l; i++ {
		c[i] = a[i] ^ b[i]
	}
	return c
}
