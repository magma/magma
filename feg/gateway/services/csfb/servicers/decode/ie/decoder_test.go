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

package ie_test

import (
	"fmt"
	"testing"

	"magma/feg/gateway/services/csfb/servicers/decode"
	"magma/feg/gateway/services/csfb/servicers/decode/ie"
	"magma/feg/gateway/services/csfb/servicers/decode/test_utils"

	"github.com/stretchr/testify/assert"
)

func TestDecodeIMSI(t *testing.T) {
	// IMSI = 1234567
	chunk, _ := test_utils.ConstructIMSI("1234567")

	IMSI, length, err := ie.DecodeIMSI(chunk)
	assert.NoError(t, err)
	assert.Equal(t, "1234567", IMSI)
	assert.Equal(t, 6, length)

	// IMSI = 123456
	chunk, _ = test_utils.ConstructIMSI("123456")

	IMSI, length, err = ie.DecodeIMSI(chunk)
	assert.NoError(t, err)
	assert.Equal(t, "123456", IMSI)
	assert.Equal(t, 6, length)

	// wrong IEI
	chunk = []byte{}
	chunk = append(chunk, byte(decode.IEITMSI))
	chunk = append(chunk, byte(4))
	chunk = append(chunk, byte(0x11))
	chunk = append(chunk, []byte{byte(0x11), byte(0x11)}...)
	chunk = append(chunk, byte(0xF1))

	IMSI, length, err = ie.DecodeIMSI(chunk)
	errorMsg := fmt.Sprintf(
		"IEI is wrong, should be 0x%02x, not 0x%02x",
		byte(decode.IEIIMSI),
		chunk[0],
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, "", IMSI)
	assert.Equal(t, -1, length)

	// chunk too short
	chunk = []byte{}
	chunk = append(chunk, byte(decode.IEIIMSI))
	chunk = append(chunk, byte(4))

	IMSI, length, err = ie.DecodeIMSI(chunk)
	errorMsg = fmt.Sprintf(
		"failed to decode IMSI: chunk too short, \n"+
			"min length of information element: %d, "+
			"number of undecoded bytes: %d",
		decode.LengthIEI+decode.LengthLengthIndicator+4,
		len(chunk),
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, "", IMSI)
	assert.Equal(t, -1, length)
}

func TestDecodeFixedLengthIE(t *testing.T) {
	// correct SGs cause
	cause, err := ie.DecodeFixedLengthIE(test_utils.ConstructDefaultSGsCause(), decode.IELengthSGsCause, decode.IEISGsCause)
	assert.NoError(t, err)
	assert.Equal(t, []byte{byte(0x11)}, cause)

	// wrong IEI
	wrongChunk1 := test_utils.ConstructDefaultSGsCause()
	wrongChunk1[0] = byte(decode.IEIIMSI)
	cause, err = ie.DecodeFixedLengthIE(wrongChunk1, decode.IELengthSGsCause, decode.IEISGsCause)
	errorMsg := fmt.Sprintf(
		"IEI is wrong, should be 0x%02x, not 0x%02x",
		byte(decode.IEISGsCause),
		wrongChunk1[0],
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, []byte{}, cause)

	// chunk too short
	wrongChunk2 := test_utils.ConstructDefaultSGsCause()
	cause, err = ie.DecodeFixedLengthIE(wrongChunk2[:2], decode.IELengthSGsCause, decode.IEISGsCause)
	errorMsg = fmt.Sprintf(
		"failed to decode SGsCause: chunk too short, \n"+
			"min length of information element: %d, "+
			"number of undecoded bytes: %d",
		decode.IELengthSGsCause,
		2,
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, []byte{}, cause)

	// wrong length indicator
	wrongChunk3 := test_utils.ConstructDefaultSGsCause()
	wrongChunk3[1] = byte(0)
	cause, err = ie.DecodeFixedLengthIE(wrongChunk3, decode.IELengthSGsCause, decode.IEISGsCause)
	errorMsg = fmt.Sprintf(
		"failed to decode SGsCause: wrong length indicator, \n"+
			"length indicator should be %d, not %d",
		decode.IELengthSGsCause-decode.LengthIEI-decode.LengthLengthIndicator,
		0,
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, []byte{}, cause)
}

func TestDecodeVariableLengthIE(t *testing.T) {
	// correct MM information
	mmInfo, ieLength, err := ie.DecodeVariableLengthIE(test_utils.ConstructDefaultMMInformation(), decode.IELengthMMInformationMin, decode.IEIMMInformation)
	assert.NoError(t, err)
	assert.Equal(t, []byte{byte(0x11), byte(0x12), byte(0x13), byte(0x14)}, mmInfo)
	assert.Equal(t, 6, ieLength)

	// wrong IEI
	wrongChunk1 := test_utils.ConstructDefaultMMInformation()
	wrongChunk1[0] = byte(decode.IEIIMSI)
	mmInfo, ieLength, err = ie.DecodeVariableLengthIE(wrongChunk1, decode.IELengthMMInformationMin, decode.IEIMMInformation)
	errorMsg := fmt.Sprintf(
		"IEI is wrong, should be 0x%02x, not 0x%02x",
		byte(decode.IEIMMInformation),
		wrongChunk1[0],
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, []byte{}, mmInfo)
	assert.Equal(t, -1, ieLength)

	// chunk too short
	wrongChunk2 := test_utils.ConstructDefaultMMInformation()
	wrongChunk2[1] = byte(0)
	mmInfo, ieLength, err = ie.DecodeVariableLengthIE(wrongChunk2[:2], decode.IELengthMMInformationMin, decode.IEIMMInformation)
	errorMsg = fmt.Sprintf(
		"failed to decode MMInformation: chunk too short, \n"+
			"min length of information element: %d, "+
			"number of undecoded bytes: %d",
		decode.IELengthMMInformationMin,
		2,
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, []byte{}, mmInfo)
	assert.Equal(t, -1, ieLength)

	// chunk too short or wrong length indicator
	wrongChunk3 := test_utils.ConstructDefaultMMInformation()
	mmInfo, ieLength, err = ie.DecodeVariableLengthIE(wrongChunk3[:3], decode.IELengthMMInformationMin, decode.IEIMMInformation)
	errorMsg = fmt.Sprintf(
		"failed to decode MMInformation: chunk too short or wrong length indicator, \n"+
			"total length of information element specified by length indicator: %d, "+
			"number of undecoded bytes: %d",
		decode.LengthIEI+decode.LengthLengthIndicator+4,
		3,
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, []byte{}, mmInfo)
	assert.Equal(t, -1, ieLength)
}

func TestDecodeLimitedLengthIE(t *testing.T) {
	mandatoryFieldLength := decode.LengthIEI + decode.LengthLengthIndicator

	// correct CLI
	cli, ieLength, err := ie.DecodeLimitedLengthIE(
		test_utils.ConstructDefaultIE(decode.IEICLI, decode.IELengthCLIMin-mandatoryFieldLength),
		decode.IELengthCLIMin,
		decode.IELengthCLIMax,
		decode.IEICLI,
	)
	assert.NoError(t, err)
	assert.Equal(t, []byte{byte(0x11)}, cli)
	assert.Equal(t, 3, ieLength)

	// wrong IEI
	cli, ieLength, err = ie.DecodeLimitedLengthIE(
		test_utils.ConstructDefaultIE(decode.IEIIMSI, decode.IELengthCLIMin-mandatoryFieldLength),
		decode.IELengthCLIMin,
		decode.IELengthCLIMax,
		decode.IEICLI,
	)
	errorMsg := fmt.Sprintf(
		"IEI is wrong, should be 0x%02x, not 0x%02x",
		byte(decode.IEICLI),
		byte(decode.IEIIMSI),
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, []byte{}, cli)
	assert.Equal(t, -1, ieLength)

	// chunk too short
	wrongChunk := test_utils.ConstructDefaultIE(decode.IEICLI, decode.IELengthCLIMax-mandatoryFieldLength)
	cli, ieLength, err = ie.DecodeLimitedLengthIE(
		wrongChunk[:2],
		decode.IELengthCLIMin,
		decode.IELengthCLIMax,
		decode.IEICLI,
	)
	errorMsg = fmt.Sprintf(
		"failed to decode CLI: chunk too short, \n"+
			"min length of information element: %d, "+
			"number of undecoded bytes: %d",
		decode.IELengthCLIMin,
		2,
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, []byte{}, cli)
	assert.Equal(t, -1, ieLength)

	// wrong length indicator
	cli, ieLength, err = ie.DecodeLimitedLengthIE(
		test_utils.ConstructDefaultIE(decode.IEICLI, decode.IELengthCLIMax-mandatoryFieldLength+1),
		decode.IELengthCLIMin,
		decode.IELengthCLIMax,
		decode.IEICLI,
	)
	errorMsg = fmt.Sprintf(
		"failed to decode CLI: wrong length indicator, \n"+
			"min value: %d, max value: %d, "+
			"length indicator: %d",
		decode.IELengthCLIMin-mandatoryFieldLength,
		decode.IELengthCLIMax-mandatoryFieldLength,
		decode.IELengthCLIMax-mandatoryFieldLength+1,
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, []byte{}, cli)
	assert.Equal(t, -1, ieLength)

	// chunk too short or wrong length indicator
	cli, ieLength, err = ie.DecodeLimitedLengthIE(
		wrongChunk[:5],
		decode.IELengthCLIMin,
		decode.IELengthCLIMax,
		decode.IEICLI,
	)
	errorMsg = fmt.Sprintf(
		"failed to decode CLI: chunk too short or wrong length indicator, \n"+
			"total length of information element specified by length indicator: %d, "+
			"number of undecoded bytes: %d",
		decode.IELengthCLIMax,
		5,
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, []byte{}, cli)
	assert.Equal(t, -1, ieLength)
}
