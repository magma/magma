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

package mock_pcrf

import (
	"fmt"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/testcore/mock_driver"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/golang/glog"
)

// Here we wrap the protobuf definitions to easily define instance methods
type GxExpectation struct {
	*protos.GxCreditControlExpectation
}

type GxAnswer struct {
	*protos.GxCreditControlAnswer
}

func (e GxExpectation) GetAnswer() interface{} {
	return GxAnswer{e.Answer}
}

func (e GxExpectation) DoesMatch(message interface{}) error {
	expected := e.ExpectedRequest
	ccr := message.(ccrMessage)
	expectedPK := mock_driver.NewCCRequestPK(expected.Imsi, expected.RequestType)
	actualImsi, _ := ccr.GetIMSI()
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
	expectedUsageReports := expected.GetUsageMonitoringReports()
	if !compareUsageMonitorsAgainstExpected(ccr.UsageMonitors, expectedUsageReports, expected.GetUsageReportDelta()) {
		return fmt.Errorf("For Request=%v, Expected: %v, Received: %v", actualPK, expectedUsageReports, ccr.UsageMonitors)
	}
	expectedET := expected.GetEventTrigger()
	if expectedET != nil && expectedET.GetValue() != int32(ccr.EventTrigger) {
		return fmt.Errorf("For Request=%v, Expected EventTrigger: %v, Received: %v", actualPK, expectedET.Value, ccr.EventTrigger)
	}
	return nil
}

func (answer GxAnswer) toAVPs() ([]*diam.AVP, uint32) {
	var avps []*diam.AVP

	ruleInstalls := answer.GetRuleInstalls()
	if ruleInstalls != nil {
		ruleInstallAVPs := toRuleInstallAVPs(
			answer.RuleInstalls.GetRuleNames(),
			answer.RuleInstalls.GetRuleBaseNames(),
			answer.RuleInstalls.GetRuleDefinitions(),
			answer.RuleInstalls.ActivationTime,
			answer.RuleInstalls.DeactivationTime,
		)
		avps = append(avps, ruleInstallAVPs...)
	}

	ruleRemovals := answer.GetRuleRemovals()
	if ruleRemovals != nil {
		ruleRemovalAVPs := toRuleRemovalAVPs(
			ruleRemovals.GetRuleNames(),
			ruleRemovals.GetRuleBaseNames())
		avps = append(avps, ruleRemovalAVPs...)
	}

	monitorInstalls := answer.GetUsageMonitoringInfos()
	for _, monitor := range monitorInstalls {
		octets := monitor.GetOctets()
		if octets == nil {
			glog.Errorf("Monitor Octets is nil, skipping.")
			continue
		}
		avps = append(avps, toUsageMonitoringInfoAVP(string(monitor.MonitoringKey), octets, monitor.MonitoringLevel))
	}

	eventTriggers := answer.GetEventTriggers()
	if eventTriggers != nil {
		avps = append(avps, toEventTriggersAVPs(eventTriggers)...)
	}
	revalidationTime := answer.GetRevalidationTime()
	if revalidationTime != nil {
		avps = append(avps, toRevalidationTimeAVPs(revalidationTime))
	}
	return avps, answer.GetResultCode()
}

func (usageMonitor *usageMonitorRequestAVP) toProtosUsageMonitorInfo() *protos.UsageMonitoringInformation {
	return &protos.UsageMonitoringInformation{
		MonitoringKey:   []byte(usageMonitor.MonitoringKey),
		MonitoringLevel: protos.MonitoringLevel(usageMonitor.Level),
		Octets: &protos.Octets{
			TotalOctets:  usageMonitor.UsedServiceUnit.TotalOctets,
			InputOctets:  usageMonitor.UsedServiceUnit.InputOctets,
			OutputOctets: usageMonitor.UsedServiceUnit.OutputOctets,
		},
	}
}

func compareUsageMonitorsAgainstExpected(actual []*usageMonitorRequestAVP, expected []*protos.UsageMonitoringInformation, delta uint64) bool {
	if expected == nil {
		return true
	}
	actualMonitorByKey := toUsageMonitorByMkey(actual)
	expectedMonitorByKey := toProtosMonitorByKey(expected)
	for mKey, expectedMonitor := range expectedMonitorByKey {
		actualMonitor, exists := actualMonitorByKey[mKey]
		if !exists {
			return false
		}
		if protos.MonitoringLevel(actualMonitor.Level) != expectedMonitor.MonitoringLevel {
			return false
		}
		actualTotal := actualMonitor.UsedServiceUnit.TotalOctets
		expectedTotal := expectedMonitor.GetOctets().GetTotalOctets()
		if !mock_driver.EqualWithinDelta(actualTotal, expectedTotal, delta) {
			return false
		}
	}
	return true
}

func toProtosMonitorByKey(monitors []*protos.UsageMonitoringInformation) map[string]*protos.UsageMonitoringInformation {
	result := map[string]*protos.UsageMonitoringInformation{}
	for _, monitor := range monitors {
		result[string(monitor.MonitoringKey)] = monitor
	}
	return result
}
