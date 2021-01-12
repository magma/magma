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

// Package handlers provided SIM Response handlers for supported SIM subtypes
package sim

import (
	"crypto/hmac"
	"crypto/sha1"

	"magma/feg/gateway/services/eap"
	aka_crypto "magma/feg/gateway/services/eap/providers/aka"
)

// GenMac calculates SIM MAC given data, nonce & K_auth
func GenMac(data, nonce, K_aut []byte) []byte {
	return HmacSha1(K_aut, data, nonce)[:16]
}

// GenMac calculates SIM MAC given data sres vector & K_auth
func GenChallengeMac(data []byte, sres [][]byte, K_aut []byte) []byte {
	return HmacSha1(K_aut, data, sres...)[:16]
}

// HmacSha1 - SHA1 based HMAC
func HmacSha1(key, data []byte, extra ...[]byte) []byte {
	h := hmac.New(sha1.New, key)
	h.Write(data)
	for _, e := range extra {
		if len(e) > 0 {
			h.Write(e)
		}
	}
	return h.Sum(nil)[:]
}

// MakeKeys returns generated K_encr, K_aut, MSK, EMSK keys for SIM Authentication (RFC 4187, section 7)
func MakeKeys(identity, nonce, versionList, selectedVersion []byte, kc [][]byte) (K_encr, K_aut, MSK, EMSK []byte) {
	x := aka_crypto.XSum(MK(identity, nonce, versionList, selectedVersion, kc))
	return x[:16], x[16:32], x[32:96], x[96:160]
}

// MK calculates & returns SIM Master Key: MK = SHA1(Identity|n*Kc| NONCE_MT| Version List| Selected Version)
func MK(identity, nonce, versionList, selectedVersion []byte, kc [][]byte) []byte {
	d := sha1.New()
	d.Write(identity)
	for _, kci := range kc {
		d.Write(kci)
	}
	d.Write(nonce)
	d.Write(versionList)
	d.Write(selectedVersion)
	return d.Sum(nil)
}

// AppendMac appends AT_MAC attribute to eap packet, signs the packet & returns the new, signed packet
// returns error if provided EAP Packet was malformed
func AppendMac(p eap.Packet, K_aut []byte) (eap.Packet, error) {
	return aka_crypto.AppendMac(p, K_aut)
}

//  GsmFromUmts1 generates GSM-Milenage (3GPP TS 55.205) auth triplet from the given UMTS quintuplet
//  using recommended SRES Derivation Function #1
//
//  Inputs:
//    ik   - 128-bit integrity key
//    ck   - 128-bit confidentiality key
//    xres - 64-bit signed response
//
//  Outputs:
//    kc   - 64-bit Kc
//    sres - 32-bit SRES
func GsmFromUmts1(ck, ik, xres []byte) (kc, sres []byte) {
	kc, sres = make([]byte, 8), make([]byte, 4)
	for i, i8 := 0, 8; i < 8; i, i8 = i+1, i8+1 {
		kc[i] = ck[i] ^ ck[i8] ^ ik[i] ^ ik[i8]
	}
	for i := 0; i < 4; i++ {
		sres[i] = xres[i] ^ xres[i+4]
	}
	return
}

//  GsmFromUmts2 generates GSM-Milenage (3GPP TS 55.205) auth triplet from the given UMTS quintuplet
//  using recommended SRES Derivation Function #2
//
//  Inputs:
//    ik   - 128-bit integrity key
//    ck   - 128-bit confidentiality key
//    xres - 64-bit signed response
//
//  Outputs:
//    kc   - 64-bit Kc
//    sres - 32-bit SRES
func GsmFromUmts2(ck, ik, xres []byte) (kc, sres []byte) {
	kc, sres = make([]byte, 8), make([]byte, 4)
	for i, i8 := 0, 8; i < 8; i, i8 = i+1, i8+1 {
		kc[i] = ck[i] ^ ck[i8] ^ ik[i] ^ ik[i8]
	}
	copy(sres, xres)
	return
}
