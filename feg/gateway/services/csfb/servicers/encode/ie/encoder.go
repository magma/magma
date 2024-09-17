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

package ie

import (
	"errors"
	"fmt"
	"strconv"

	"magma/feg/gateway/services/csfb/servicers/decode"
)

func EncodeIMSI(imsi string) ([]byte, error) {
	encodedIMSI := []byte{byte(decode.IEIIMSI)}
	lengthIndicator := len(imsi)/2 + 1
	mandatoryFieldsLength := decode.LengthIEI + decode.LengthLengthIndicator
	minLength := decode.IELengthIMSIMin - mandatoryFieldsLength
	maxLength := decode.IELengthIMSIMax - mandatoryFieldsLength
	if lengthIndicator < minLength || lengthIndicator > maxLength {
		errorMsg := fmt.Sprintf(
			"failed to encode IMSI, value field length violation, min length: %d, max length: %d, actual length: %d",
			minLength,
			maxLength,
			lengthIndicator,
		)
		return []byte{}, errors.New(errorMsg)
	}
	encodedIMSI = append(encodedIMSI, byte(lengthIndicator))
	// third byte with the parity bit
	digit, err := strconv.Atoi(imsi[0:1])
	if err != nil {
		return []byte{}, err
	}
	var thirdByte byte
	if len(imsi)%2 == 1 {
		thirdByte = byte(9) + byte(digit<<4)
	} else {
		thirdByte = byte(1) + byte(digit<<4)
	}
	encodedIMSI = append(encodedIMSI, thirdByte)
	// the rest of bytes
	for idx := 1; idx+1 < len(imsi); idx += 2 {
		firstDigit, err := strconv.Atoi(string(imsi[idx]))
		if err != nil {
			return []byte{}, err
		}
		secondDigit, err := strconv.Atoi(string(imsi[idx+1]))
		if err != nil {
			return []byte{}, err
		}
		encodedIMSI = append(encodedIMSI, byte((secondDigit<<4)+firstDigit))
	}
	if len(imsi)%2 == 0 {
		digit, err = strconv.Atoi(imsi[len(imsi)-1:])
		if err != nil {
			return []byte{}, err
		}
		encodedIMSI = append(encodedIMSI, byte(0xF0+digit))
	}

	return encodedIMSI, nil
}

func EncodeMMEName(mmeName string) ([]byte, error) {
	lengthIndicator := len(mmeName)
	mandatoryFieldLength := decode.LengthIEI + decode.LengthLengthIndicator
	if lengthIndicator != decode.IELengthMMEName-mandatoryFieldLength {
		errorsMsg := fmt.Sprintf(
			"failed to encode MME Name, value field length violation, length of value field should be %d, actual length: %d",
			decode.IELengthMMEName-mandatoryFieldLength,
			lengthIndicator,
		)
		return []byte{}, errors.New(errorsMsg)
	}

	// MME name in rpc message will look like:
	// .mmec01.mmegi0001.mme.EPC.mnc001.mcc001.3gppnetwork.org
	encodedMMEName := constructMessage(
		decode.IEIMMEName,
		lengthIndicator,
		[]byte(mmeName[1:]),
	)

	// H-core MSC needs a 0x00 at the end of MME name
	encodedMMEName = append(encodedMMEName, byte(0x00))

	return encodedMMEName, nil
}

func EncodeFixedLengthIE(iei decode.InformationElementIdentifier, ieLength int, valueField []byte) ([]byte, error) {
	lengthIndicator := len(valueField)
	mandatoryFieldLength := decode.LengthIEI + decode.LengthLengthIndicator
	if lengthIndicator != ieLength-mandatoryFieldLength {
		errorMsg := fmt.Sprintf(
			"failed to encode %s, value field length violation, length of value field should be %d, actual length: %d",
			decode.IEINamesByCode[iei],
			ieLength-mandatoryFieldLength,
			lengthIndicator,
		)
		return []byte{}, errors.New(errorMsg)
	}

	return constructMessage(iei, lengthIndicator, valueField), nil
}

func EncodeLimitedLengthIE(iei decode.InformationElementIdentifier, minIELength int, maxIELength int, valueField []byte) ([]byte, error) {
	lengthIndicator := len(valueField)
	mandatoryFieldLength := decode.LengthIEI + decode.LengthLengthIndicator
	if lengthIndicator < minIELength-mandatoryFieldLength || lengthIndicator > maxIELength-mandatoryFieldLength {
		errorMsg := fmt.Sprintf(
			"failed to encode %s, value field length violation, min length: %d, max length: %d, actual length: %d",
			decode.IEINamesByCode[iei],
			minIELength-mandatoryFieldLength,
			maxIELength-mandatoryFieldLength,
			lengthIndicator,
		)
		return []byte{}, errors.New(errorMsg)
	}
	return constructMessage(iei, lengthIndicator, valueField), nil
}

func EncodeVariableLengthIE(iei decode.InformationElementIdentifier, minIELength int, valueField []byte) ([]byte, error) {
	lengthIndicator := len(valueField)
	mandatoryFieldLength := decode.LengthIEI + decode.LengthLengthIndicator
	if lengthIndicator < minIELength-mandatoryFieldLength {
		errorMsg := fmt.Sprintf(
			"failed to encode %s, value field length violation, min length: %d, actual length: %d",
			decode.IEINamesByCode[iei],
			minIELength-mandatoryFieldLength,
			lengthIndicator,
		)
		return []byte{}, errors.New(errorMsg)
	}
	return constructMessage(iei, lengthIndicator, valueField), nil
}

func constructMessage(iei decode.InformationElementIdentifier, lengthIndicator int, valueField []byte) []byte {
	var message []byte
	message = append(message, byte(iei))
	message = append(message, byte(lengthIndicator))
	message = append(message, valueField...)

	return message
}
