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

package test_utils

import (
	"strconv"

	"magma/feg/gateway/services/csfb/servicers/decode"
)

func ConstructIMSI(IMSIstr string) ([]byte, error) {
	var imsi []byte
	imsi = append(imsi, byte(decode.IEIIMSI))
	imsi = append(imsi, byte(len(IMSIstr)/2+1))
	// third byte with the parity bit
	digit, err := strconv.Atoi(IMSIstr[0:1])
	if err != nil {
		return []byte{}, err
	}
	var thirdByte byte
	if len(IMSIstr)%2 == 1 {
		thirdByte = byte(9) + byte(digit<<4)
	} else {
		thirdByte = byte(1) + byte(digit<<4)
	}
	imsi = append(imsi, thirdByte)
	// the rest of bytes
	for idx := 1; idx+1 < len(IMSIstr); idx += 2 {
		firstDigit, err := strconv.Atoi(string(IMSIstr[idx]))
		if err != nil {
			return []byte{}, err
		}
		secondDigit, err := strconv.Atoi(string(IMSIstr[idx+1]))
		if err != nil {
			return []byte{}, err
		}
		imsi = append(imsi, byte((secondDigit<<4)+firstDigit))
	}
	if len(IMSIstr)%2 == 0 {
		digit, err = strconv.Atoi(IMSIstr[len(IMSIstr)-1:])
		if err != nil {
			return []byte{}, err
		}
		imsi = append(imsi, byte(0xF0+digit))
	}

	return imsi, nil
}

func ConstructDefaultTMSI() []byte {
	var tmsi []byte
	tmsi = append(tmsi, byte(decode.IEITMSI))
	tmsi = append(tmsi, byte(decode.IELengthTMSI-decode.LengthIEI-decode.LengthLengthIndicator))
	tmsi = append(tmsi, []byte{byte(0x11), byte(0x12), byte(0x13), byte(0x14)}...)

	return tmsi
}

func ConstructMobileIdentity(IMSIstr string, TMSI []byte) ([]byte, error) {
	if IMSIstr != "" {
		imsi, err := ConstructIMSI(IMSIstr)
		if err != nil {
			return []byte{}, err
		}
		imsi[0] = byte(decode.IEIMobileIdentity)
		return imsi, nil
	} else {
		tmsi := []byte{byte(decode.IEIMobileIdentity)}
		tmsi = append(tmsi, byte(5)) // length indicator
		tmsi = append(tmsi, byte(0xF4))
		tmsi = append(tmsi, TMSI...)
		return tmsi, nil
	}
}

func ConstructDefaultMMEName() []byte {
	var mmeName []byte
	mmeName = append(mmeName, byte(decode.IEIMMEName))
	mmeName = append(mmeName, byte(decode.IELengthMMEName-decode.LengthIEI-decode.LengthLengthIndicator))
	mmeName = append(mmeName, []byte("abcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcde")...)

	return mmeName
}

func ConstructDefaultVLRName() []byte {
	var vlrName []byte
	vlrName = append(vlrName, byte(decode.IEIVLRName))
	vlrName = append(vlrName, byte(16))
	vlrName = append(vlrName, []byte("www.facebook.com")...)

	return vlrName
}

func ConstructDefaultMMInformation() []byte {
	var mmInfo []byte
	mmInfo = append(mmInfo, byte(decode.IEIMMInformation))
	mmInfo = append(mmInfo, byte(4))
	mmInfo = append(mmInfo, []byte{byte(0x11), byte(0x12), byte(0x13), byte(0x14)}...)

	return mmInfo
}

func ConstructDefaultLocationAreaIdentifier() []byte {
	var lai []byte
	lai = append(lai, byte(decode.IEILocationAreaIdentifier))
	lai = append(lai, byte(5))
	lai = append(lai, []byte{byte(0x11), byte(0x12), byte(0x13), byte(0x14), byte(0x15)}...)

	return lai
}

func ConstructDefaultSGsCause() []byte {
	var cause []byte
	cause = append(cause, byte(decode.IEISGsCause))
	cause = append(cause, byte(1))
	cause = append(cause, byte(0x11))

	return cause
}

func DefaultVal(lengthIndicator int) []byte {
	var chunk []byte
	for i := 0; i < lengthIndicator; i += 1 {
		chunk = append(chunk, byte(0x11))
	}
	return chunk
}

func ConstructDefaultIE(iei decode.InformationElementIdentifier, lengthIndicator int) []byte {
	var chunk []byte
	chunk = append(chunk, byte(iei))
	chunk = append(chunk, byte(lengthIndicator))
	for i := 0; i < lengthIndicator; i += 1 {
		chunk = append(chunk, byte(0x11))
	}

	return chunk
}
