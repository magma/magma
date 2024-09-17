/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package servicers

import (
	"errors"
	"regexp"

	"magma/cwf/cloud/go/protos"
)

// validateUEData ensures that a UE data proto is not nil and that it contains
// a valid IMSI, key, and opc.
func validateUEData(ue *protos.UEConfig) error {
	if ue == nil {
		return errors.New("Invalid Argument: UE data cannot be nil")
	}
	errIMSI := validateUEIMSI(ue.GetImsi())
	if errIMSI != nil {
		return errIMSI
	}
	errkey := validateUEKey(ue.GetAuthKey())
	if errkey != nil {
		return errkey
	}
	erropc := validateUEOpc(ue.GetAuthOpc())
	if erropc != nil {
		return erropc
	}
	return nil
}

// validateUEIMSI ensures that a UE's IMSI can be stored.
func validateUEIMSI(imsi string) error {
	if len(imsi) < 5 || len(imsi) > 15 {
		return errors.New("Invalid Argument: IMSI must be between 5 and 15 digits long")
	}
	isOnlyDigits, err := regexp.MatchString(`^[0-9]*$`, imsi)
	if err != nil || !isOnlyDigits {
		return errors.New("Invalid Argument: IMSI must only be digits")
	}
	return nil
}

// validateUEKey ensures that a UE's key can be stored.
func validateUEKey(k []byte) error {
	if k == nil {
		return errors.New("Invalid Argument: key cannot be nil")
	}
	if len(k) != 16 {
		return errors.New("Invalid Argument: key must be 16 bytes")
	}
	return nil
}

// validateUEOpc ensures that a UE's opc can be stored.
func validateUEOpc(opc []byte) error {
	if opc == nil {
		return errors.New("Invalid Argument: opc cannot be nil")
	}
	if len(opc) != 16 {
		return errors.New("Invalid Argument: opc must be 16 bytes")
	}
	return nil
}

// validateUEMSISDN ensures that a UE data is valid for HssLess authentication.
func validateUEMSISDN(msisdn string) error {
	isOnlyDigits, err := regexp.MatchString(`^[0-9]*$`, msisdn)
	if err != nil || !isOnlyDigits {
		return errors.New("Invalid Argument: MSISDN must only be digits")
	}
	return nil
}

// validateUEDataForHssLess ensures that a UE data proto is not nil and that it contains
// a valid MSISDN.
func validateUEDataForHssLess(ue *protos.UEConfig) error {
	if ue == nil {
		return errors.New("Invalid Argument: UE data cannot be nil")
	}
	/*Validate MSISDN */
	errMSISDN := validateUEMSISDN(ue.GetHsslessCfg().GetMsisdn())
	if errMSISDN != nil {
		return errMSISDN
	}

	return nil
}
