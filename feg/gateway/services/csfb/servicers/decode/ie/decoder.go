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
	"magma/feg/gateway/services/csfb/servicers/decode"
)

func DecodeIMSI(restOfChunk []byte) (string, int, error) {
	valueField, ieLength, err := DecodeLimitedLengthIE(
		restOfChunk,
		decode.IELengthIMSIMin,
		decode.IELengthIMSIMax,
		decode.IEIIMSI,
	)
	if err != nil {
		return "", -1, err
	}

	imsi, err := ExtractIMSIString(valueField)
	if err != nil {
		return "", -1, err
	}

	return imsi, ieLength, nil
}

// Decoder for TMSI, MME Name, LAI, SGs Cause, CLI, Service Indicator, SS Code, LCS Indicator, ChannelNeeded, eMLPPPriority
func DecodeFixedLengthIE(restOfChunk []byte, ieLength int, iei decode.InformationElementIdentifier) ([]byte, error) {
	err := validateFixedLengthIE(restOfChunk, ieLength, iei)
	if err != nil {
		return []byte{}, err
	}

	return restOfChunk[decode.LengthIEI+decode.LengthLengthIndicator : ieLength], nil
}

// Decoder for VLR Name, MM Information, LCS Client Identity
func DecodeVariableLengthIE(restOfChunk []byte, minLength int, iei decode.InformationElementIdentifier) ([]byte, int, error) {
	ieLength, err := validateVariableLengthIE(restOfChunk, minLength, iei)
	if err != nil {
		return []byte{}, -1, err
	}

	return restOfChunk[decode.LengthIEI+decode.LengthLengthIndicator : ieLength], ieLength, nil
}

// Decoder for CLI and NAS Message Container
func DecodeLimitedLengthIE(restOfChunk []byte, minLength int, maxLength int, iei decode.InformationElementIdentifier) ([]byte, int, error) {
	ieLength, err := validateLimitedLengthIE(restOfChunk, minLength, maxLength, iei)
	if err != nil {
		return []byte{}, -1, err
	}

	return restOfChunk[decode.LengthIEI+decode.LengthLengthIndicator : ieLength], ieLength, nil
}
