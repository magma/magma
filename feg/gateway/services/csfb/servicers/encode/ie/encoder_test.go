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
	"magma/feg/gateway/services/csfb/servicers/encode/ie"

	"github.com/stretchr/testify/assert"
)

func TestEncodeIMSI(t *testing.T) {
	// successfully encode
	imsi := "111111111111"
	_, err := ie.EncodeIMSI(imsi)
	assert.NoError(t, err)

	// value field too short
	imsi = "1"
	encodedIMSI, err := ie.EncodeIMSI(imsi)
	mandatoryFieldLength := decode.LengthIEI + decode.LengthLengthIndicator
	errorMsg := fmt.Sprintf(
		"failed to encode IMSI, value field length violation, min length: %d, max length: %d, actual length: %d",
		decode.IELengthIMSIMin-mandatoryFieldLength,
		decode.IELengthIMSIMax-mandatoryFieldLength,
		len(imsi)/2+1,
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, []byte{}, encodedIMSI)

	// value field too long
	imsi = "111111111111111111111"
	encodedIMSI, err = ie.EncodeIMSI(imsi)
	errorMsg = fmt.Sprintf(
		"failed to encode IMSI, value field length violation, min length: %d, max length: %d, actual length: %d",
		decode.IELengthIMSIMin-mandatoryFieldLength,
		decode.IELengthIMSIMax-mandatoryFieldLength,
		len(imsi)/2+1,
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, []byte{}, encodedIMSI)
}

func TestEncodeMMEName(t *testing.T) {
	mmeName := ".mmec01.mmegi0001.mme.EPC.mnc001.mcc001.3gppnetwork.org"
	encodedMMEName, err := ie.EncodeMMEName(mmeName)
	assert.NoError(t, err)

	expectedEncodedMMEName, err := ie.EncodeFixedLengthIE(
		decode.IEIMMEName,
		decode.IELengthMMEName,
		[]byte("mmec01.mmegi0001.mme.EPC.mnc001.mcc001.3gppnetwork.org "),
	)
	assert.NoError(t, err)
	// replace the ending space with 0x00
	expectedEncodedMMEName[len(expectedEncodedMMEName)-1] = byte(0x00)

	assert.Equal(t, expectedEncodedMMEName, encodedMMEName)
}

func TestEncodeFixedLengthIE(t *testing.T) {
	// successfully encode TMSI
	tmsi := []byte{byte(0x11), byte(0x12), byte(0x13), byte(0x14)}
	_, err := ie.EncodeFixedLengthIE(decode.IEITMSI, decode.IELengthTMSI, tmsi)
	assert.NoError(t, err)

	// wrong length
	encodedTMSI, err := ie.EncodeFixedLengthIE(decode.IEITMSI, decode.IELengthTMSI, tmsi[:2])
	mandatoryFieldLength := decode.LengthIEI + decode.LengthLengthIndicator
	errorMsg := fmt.Sprintf(
		"failed to encode TMSI, value field length violation, length of value field should be %d, actual length: %d",
		decode.IELengthTMSI-mandatoryFieldLength,
		2,
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, []byte{}, encodedTMSI)
}

func TestEncodeLimitedLengthIE(t *testing.T) {
	// successfully encode NAS message container
	nasMessageContainer := []byte{byte(0x11), byte(0x12), byte(0x13), byte(0x14)}
	_, err := ie.EncodeLimitedLengthIE(
		decode.IEINASMessageContainer,
		decode.IELengthNASMessageContainerMin,
		decode.IELengthNASMessageContainerMax,
		nasMessageContainer,
	)
	assert.NoError(t, err)

	// value field too short
	nasMessageContainer = []byte{}
	encodedNASMessageContainer, err := ie.EncodeLimitedLengthIE(
		decode.IEINASMessageContainer,
		decode.IELengthNASMessageContainerMin,
		decode.IELengthNASMessageContainerMax,
		nasMessageContainer,
	)
	mandatoryFieldLength := decode.LengthIEI + decode.LengthLengthIndicator
	errorMsg := fmt.Sprintf(
		"failed to encode NASMessageContainer, value field length violation, min length: %d, max length: %d, actual length: %d",
		decode.IELengthNASMessageContainerMin-mandatoryFieldLength,
		decode.IELengthNASMessageContainerMax-mandatoryFieldLength,
		len(nasMessageContainer),
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, []byte{}, encodedNASMessageContainer)

	// value field too long
	nasMessageContainer = make([]byte, decode.IELengthNASMessageContainerMax-mandatoryFieldLength+1)
	encodedNASMessageContainer, err = ie.EncodeLimitedLengthIE(
		decode.IEINASMessageContainer,
		decode.IELengthNASMessageContainerMin,
		decode.IELengthNASMessageContainerMax,
		nasMessageContainer,
	)
	mandatoryFieldLength = decode.LengthIEI + decode.LengthLengthIndicator
	errorMsg = fmt.Sprintf(
		"failed to encode NASMessageContainer, value field length violation, min length: %d, max length: %d, actual length: %d",
		decode.IELengthNASMessageContainerMin-mandatoryFieldLength,
		decode.IELengthNASMessageContainerMax-mandatoryFieldLength,
		len(nasMessageContainer),
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, []byte{}, encodedNASMessageContainer)
}

func TestEncodeVariableLengthIE(t *testing.T) {
	// successfully encode VLR name
	vlrName := "www.facebook.com"
	_, err := ie.EncodeVariableLengthIE(decode.IEIVLRName, decode.IELengthVLRNameMin, []byte(vlrName))
	assert.NoError(t, err)

	// value field too short
	vlrName = ""
	encodedVLRName, err := ie.EncodeVariableLengthIE(decode.IEIVLRName, decode.IELengthVLRNameMin, []byte(vlrName))
	mandatoryFieldLength := decode.LengthIEI + decode.LengthLengthIndicator
	errorMsg := fmt.Sprintf(
		"failed to encode VLRName, value field length violation, min length: %d, actual length: %d",
		decode.IELengthVLRNameMin-mandatoryFieldLength,
		len(vlrName),
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, []byte{}, encodedVLRName)
}
