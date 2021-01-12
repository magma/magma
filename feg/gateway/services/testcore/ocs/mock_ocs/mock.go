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

package mock_ocs

import (
	"fmt"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/session_proxy/credit_control/gy"
	"magma/feg/gateway/services/testcore/mock_driver"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/golang/glog"
)

func (c ccrCredit) String() string {
	unit := c.UsedServiceUnit
	return fmt.Sprintf("RG=%v, Reason=%v Usage=(%v,%v,%v), Used-Service-Unit-Reason=%v",
		c.RatingGroup, c.ReportingReason, unit.InputOctets, unit.OutputOctets, unit.TotalOctets, unit.ReportingReason)
}

// Here we wrap the protobuf definitions to easily define instance methods
type GyExpectation struct {
	*protos.GyCreditControlExpectation
}

type GyAnswer struct {
	*protos.GyCreditControlAnswer
}

func (e GyExpectation) GetAnswer() interface{} {
	return GyAnswer{e.Answer}
}

func (e GyExpectation) DoesMatch(message interface{}) error {
	expected := e.ExpectedRequest
	ccr := message.(ccrMessage)
	expectedPK := mock_driver.NewCCRequestPK(expected.Imsi, expected.RequestType)
	actualImsi := ccr.GetIMSI()
	actualPK := mock_driver.NewCCRequestPK(actualImsi, protos.CCRequestType(ccr.RequestType))
	// For better readability of errors, we will check for the IMSI and the request type first.
	if expectedPK != actualPK {
		return fmt.Errorf("Expected: %v, Received: %v", expectedPK, actualPK)
	}
	expectedRN := expected.GetRequestNumber()
	if expectedRN != nil {
		if err := mock_driver.CompareRequestNumber(actualPK, expectedRN, ccr.RequestNumber); err != nil {
			return err
		}
	}
	msccByKey := toActualCreditByKey(ccr.MSCC)
	if !compareMsccAgainstExpected(msccByKey, expected.GetMscc(), expected.GetUsageReportDelta()) {
		return fmt.Errorf("For Request=%v, Expected: %v, Received: %v", actualPK, expected.GetMscc(), msccByKey)
	}
	return nil
}

func (answer GyAnswer) toAVPs() ([]*diam.AVP, uint32) {
	avps := make([]*diam.AVP, 0, len(answer.GetQuotaGrants()))
	for _, grant := range answer.QuotaGrants {
		avps = append(
			avps,
			toGrantedUnitsAVP(
				grant.GetResultCode(),
				grant.GetValidityTime(),
				grant.GetGrantedServiceUnit(),
				grant.GetIsFinalCredit(),
				grant.GetRatingGroup(),
				grant.GetFinalUnitIndication().GetFinalUnitAction(),
				grant.GetFinalUnitIndication().GetRedirectServer().GetRedirectServerAddress(),
				grant.GetFinalUnitIndication().GetRestrictRules()))
	}
	return avps, answer.GetResultCode()
}

func compareMsccAgainstExpected(actualMscc map[uint32]ccrCredit, expectedMscc []*protos.MultipleServicesCreditControl, delta uint64) bool {
	if expectedMscc == nil {
		return true
	}
	expectedCreditByKey := toExpectedCreditByKey(expectedMscc)
	for rg, expectedCredit := range expectedCreditByKey {
		actualCredit, exists := actualMscc[rg]
		if !exists {
			return false
		}
		// If there is no expectation set for UsedServiceUnit, don't assert
		if expectedCredit.UsedServiceUnit != nil {
			actualTotal := actualCredit.UsedServiceUnit.TotalOctets
			expectedTotal := expectedCredit.UsedServiceUnit.TotalOctets
			if !mock_driver.EqualWithinDelta(actualTotal, expectedTotal, delta) {
				return false
			}
		}
		switch gy.UsedCreditsType(expectedCredit.UpdateType) {
		case gy.VALIDITY_TIMER_EXPIRED, gy.FINAL, gy.FORCED_REAUTHORISATION:
			return expectedCredit.UpdateType == int32(actualCredit.ReportingReason)
		case gy.QUOTA_EXHAUSTED:
			return expectedCredit.UpdateType == int32(actualCredit.UsedServiceUnit.ReportingReason)
		default:
			return true
		}
	}
	return true
}

func toExpectedCreditByKey(mscc []*protos.MultipleServicesCreditControl) map[uint32]*protos.MultipleServicesCreditControl {
	msccByRG := map[uint32]*protos.MultipleServicesCreditControl{}
	for _, credit := range mscc {
		msccByRG[credit.RatingGroup] = credit
	}
	return msccByRG
}

func toActualCreditByKey(mscc []*ccrCredit) map[uint32]ccrCredit {
	msccByRG := map[uint32]ccrCredit{}
	if mscc == nil {
		return msccByRG
	}
	for _, credit := range mscc {
		if credit == nil {
			glog.Errorf("Received a nil MSCC... Skipping")
			continue
		}
		msccByRG[credit.RatingGroup] = *credit
	}
	return msccByRG
}
