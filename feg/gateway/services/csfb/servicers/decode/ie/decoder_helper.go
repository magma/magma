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
	"bytes"
	"errors"
	"fmt"
	"strconv"

	"magma/feg/gateway/services/csfb/servicers/decode"
)

func ExtractIMSIString(valueField []byte) (string, error) {
	parityBit := getIMSIParity(valueField[0])

	var stringBuffer bytes.Buffer
	_, err := stringBuffer.WriteString(getIMSIDigitFromLeftHalfOfByte(valueField[0]))
	if err != nil {
		return "", err
	}

	for _, value := range valueField[1 : len(valueField)-1] {
		_, err = stringBuffer.WriteString(getIMSIDigitsFromByte(value))
		if err != nil {
			return "", err
		}
	}

	if parityBit == 1 {
		_, err = stringBuffer.WriteString(getIMSIDigitsFromByte(valueField[len(valueField)-1]))
		if err != nil {
			return "", err
		}
	} else {
		_, err = stringBuffer.WriteString(getIMSIDigitFromRightHalfOfByte(valueField[len(valueField)-1]))
		if err != nil {
			return "", err
		}
	}

	return stringBuffer.String(), nil
}

func validateFixedLengthIE(restOfChunk []byte, ieLength int, iei decode.InformationElementIdentifier) error {
	err := validateIEIdentifierAndLength(restOfChunk, ieLength, iei)
	if err != nil {
		return err
	}

	lengthIndicator := int(restOfChunk[1])
	correctLength := ieLength - decode.LengthIEI - decode.LengthLengthIndicator
	if lengthIndicator != correctLength {
		errorMsg := fmt.Sprintf(
			"failed to decode %s: "+
				"wrong length indicator, \n"+
				"length indicator should be %d, not %d",
			decode.IEINamesByCode[iei],
			correctLength,
			lengthIndicator,
		)
		return errors.New(errorMsg)
	}

	return nil
}

func validateVariableLengthIE(restOfChunk []byte, minLength int, iei decode.InformationElementIdentifier) (int, error) {
	err := validateIEIdentifierAndLength(restOfChunk, minLength, iei)
	if err != nil {
		return -1, err
	}

	lengthIndicator := int(restOfChunk[1])
	ieLength := decode.LengthIEI + decode.LengthLengthIndicator + lengthIndicator
	if len(restOfChunk) < ieLength {
		errorMsg := fmt.Sprintf(
			"failed to decode %s: "+
				"chunk too short or wrong length indicator, \n"+
				"total length of information element specified by length indicator: %d, "+
				"number of undecoded bytes: %d",
			decode.IEINamesByCode[iei],
			ieLength,
			len(restOfChunk),
		)
		return -1, errors.New(errorMsg)
	}

	return ieLength, nil
}

func validateLimitedLengthIE(restOfChunk []byte, minLength int, maxLength int, iei decode.InformationElementIdentifier) (int, error) {
	err := validateIEIdentifierAndLength(restOfChunk, minLength, iei)
	if err != nil {
		return -1, err
	}

	lengthIndicator := int(restOfChunk[1])
	mandatoryFieldLength := decode.LengthIEI + decode.LengthLengthIndicator
	if lengthIndicator < minLength-mandatoryFieldLength || lengthIndicator > maxLength-mandatoryFieldLength {
		errorMsg := fmt.Sprintf(
			"failed to decode %s: "+
				"wrong length indicator, \n"+
				"min value: %d, max value: %d, length indicator: %d",
			decode.IEINamesByCode[iei],
			minLength-mandatoryFieldLength,
			maxLength-mandatoryFieldLength,
			lengthIndicator,
		)
		return -1, errors.New(errorMsg)
	}
	ieLength := decode.LengthIEI + decode.LengthLengthIndicator + lengthIndicator
	if len(restOfChunk) < ieLength {
		errorMsg := fmt.Sprintf(
			"failed to decode %s: "+
				"chunk too short or wrong length indicator, \n"+
				"total length of information element specified by length indicator: %d, "+
				"number of undecoded bytes: %d",
			decode.IEINamesByCode[iei],
			ieLength,
			len(restOfChunk),
		)
		return -1, errors.New(errorMsg)
	}

	return ieLength, nil
}

func validateIEIdentifierAndLength(restOfChunk []byte, minLength int, iei decode.InformationElementIdentifier) error {
	if len(restOfChunk) < minLength {
		errorMsg := fmt.Sprintf(
			"failed to decode %s: chunk too short, \n"+
				"min length of information element: %d, number of undecoded bytes: %d",
			decode.IEINamesByCode[iei],
			minLength,
			len(restOfChunk),
		)
		return errors.New(errorMsg)
	}

	if decode.InformationElementIdentifier(restOfChunk[0]) != iei {
		errorMsg := fmt.Sprintf(
			"IEI is wrong, should be 0x%02x, not 0x%02x",
			byte(iei),
			restOfChunk[0],
		)
		return errors.New(errorMsg)
	}

	return nil
}

func getIMSIParity(b byte) byte {
	return b & 0x08 >> 3
}

func getIMSIDigitFromLeftHalfOfByte(b byte) string {
	return strconv.Itoa(int(b >> 4))
}

func getIMSIDigitFromRightHalfOfByte(b byte) string {
	return strconv.Itoa(int(b & 0x0F))
}

func getIMSIDigitsFromByte(b byte) string {
	return getIMSIDigitFromRightHalfOfByte(b) + getIMSIDigitFromLeftHalfOfByte(b)
}
