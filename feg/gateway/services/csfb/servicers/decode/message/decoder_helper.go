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

package message

import (
	"errors"
	"fmt"

	"magma/feg/gateway/services/csfb/servicers/decode"
	"magma/feg/gateway/services/csfb/servicers/decode/ie"
)

// IMSIAndTMSI wraps up IMSI or TMSI when they are optional in a message
type IMSIAndTMSI struct {
	IMSI string
	TMSI []byte
}

func validateMessageLength(chunk []byte, minLength int, maxLength int) error {
	if len(chunk) >= minLength && len(chunk) <= maxLength {
		return nil
	}
	errorMsg := fmt.Sprintf(
		"wrong chunk size, min length: %d, max length: %d, actual length of the chunk: %d",
		minLength,
		maxLength,
		len(chunk),
	)
	return errors.New(errorMsg)
}

func validateMessageMinLength(chunk []byte, minLength int) error {
	if len(chunk) >= minLength {
		return nil
	}
	errorMsg := fmt.Sprintf(
		"wrong chunk size, min length: %d, actual length of the chunk: %d",
		minLength,
		len(chunk),
	)
	return errors.New(errorMsg)
}

func readMobileIdentity(chunk []byte, imsiLength int, readIdx int) (IMSIAndTMSI, error) {
	minLength := decode.IELengthMessageType + imsiLength + decode.IELengthLocationAreaIdentifier
	if len(chunk) == minLength {
		return IMSIAndTMSI{}, nil
	}

	if len(chunk) < minLength {
		errorMsg := fmt.Sprintf(
			"wrong chunk size, length without mobile identity: %d, actual length of the chunk: %d",
			minLength,
			len(chunk),
		)
		return IMSIAndTMSI{}, errors.New(errorMsg)
	}

	mobileIdentity, _, err := ie.DecodeLimitedLengthIE(
		chunk[readIdx:],
		decode.IELengthMobileIdentityMin,
		decode.IELengthMobileIdentityMax,
		decode.IEIMobileIdentity,
	)
	if err != nil {
		return IMSIAndTMSI{}, err
	}
	identityType := mobileIdentity[0] & decode.MobileIdentityTypeMask
	if identityType == decode.MobileIdentityIMSI {
		// IMSI
		imsi, err := ie.ExtractIMSIString(mobileIdentity)
		if err != nil {
			return IMSIAndTMSI{}, err
		}
		return IMSIAndTMSI{IMSI: imsi}, nil
	} else if identityType == decode.MobileIdentityTMSI {
		// TMSI
		if mobileIdentity[0] != byte(decode.MobileIdentityTMSIFirstByte) {
			errorMsg := fmt.Sprintf(
				"byte 3 of mobile identity field as TMSI should be 0xF4, not 0x%X",
				mobileIdentity[0],
			)
			return IMSIAndTMSI{}, errors.New(errorMsg)
		}
		return IMSIAndTMSI{TMSI: mobileIdentity[1:]}, nil
	} else {
		errorMsg := fmt.Sprintf(
			"cannot recognize the identity type %x for mobile identity field",
			identityType,
		)
		return IMSIAndTMSI{}, errors.New(errorMsg)
	}
}

func readLAI(chunk []byte, imsiLength int, readIdx int) ([]byte, error) {
	minLength := decode.IELengthMessageType + imsiLength + decode.LengthRejectCause
	if len(chunk) == minLength+decode.IELengthLocationAreaIdentifier {
		lai, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthLocationAreaIdentifier, decode.IEILocationAreaIdentifier)
		if err != nil {
			return nil, err
		}
		return lai, nil
	} else if len(chunk) == minLength {
		return nil, nil
	} else {
		errorMsg := fmt.Sprintf(
			"wrong chunk size, length without LAI: %d, length with LAI: %d, actual length of the chunk: %d",
			minLength,
			minLength+decode.IELengthLocationAreaIdentifier,
			len(chunk),
		)
		return nil, errors.New(errorMsg)
	}
}
