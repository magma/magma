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

package diameter

import (
	"encoding/binary"
	"fmt"
)

// EncodePLMNID encodes a PLMN into its byte encoding
// See https://www.arib.or.jp/english/html/overview/doc/STD-T63v10_70/5_Appendix/Rel11/29/29272-bb0.pdf#page=73
func EncodePLMNID(plmn string) ([]byte, error) {
	if len(plmn) < 5 || len(plmn) > 6 {
		return []byte{}, fmt.Errorf("Invalid PLMN length: %s", plmn)
	}
	rawPLMN := make([]byte, 6)
	for i, r := range plmn {
		// transform rune digit to byte
		rawPLMN[i] = byte(r - '0')
	}
	encodedPLMN := []byte{
		rawPLMN[0] | (rawPLMN[1] << 4),
		rawPLMN[2] | (rawPLMN[5] << 4),
		rawPLMN[3] | (rawPLMN[4] << 4),
	}
	if len(plmn) < 6 {
		encodedPLMN[1] |= 0xF0
	}
	return encodedPLMN, nil
}

// EncodeUserLocation encodes a PLMN ID, TAI, and ECI into the correct encoding
// for 3GPP-User-Location-Info. Normally this value is provided by the EPC, but
// this function can be used for command lines.
// Encoding defined in 3GPP 29.061 Section 16.4.7.2
func EncodeUserLocation(plmn string, tai uint16, eci uint32) ([]byte, error) {
	encodedPLMN, err := EncodePLMNID(plmn)
	if err != nil {
		return []byte{}, err
	}

	returnBytes := []byte{}

	// Location Type
	var locationType byte = 0x82 // TAI & ECGI
	returnBytes = append(returnBytes, locationType)

	// Tracking Area Identity
	returnBytes = append(returnBytes, encodedPLMN...)

	taiBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(taiBytes, tai)
	returnBytes = append(returnBytes, taiBytes...)

	// E-UTRAN Cell Global Identifier
	returnBytes = append(returnBytes, encodedPLMN...)

	eciBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(eciBytes, eci)
	returnBytes = append(returnBytes, eciBytes...)

	return returnBytes, nil
}
