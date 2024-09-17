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

package protos

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
)

func NewGxCreditControlExpectation() *GxCreditControlExpectation {
	return &GxCreditControlExpectation{}
}

func (m *GxCreditControlExpectation) Expect(ccr *GxCreditControlRequest) *GxCreditControlExpectation {
	m.ExpectedRequest = ccr
	return m
}

func (m *GxCreditControlExpectation) Return(cca *GxCreditControlAnswer) *GxCreditControlExpectation {
	m.Answer = cca
	return m
}

func NewGxCCRequest(imsi string, requestType CCRequestType) *GxCreditControlRequest {
	return &GxCreditControlRequest{Imsi: imsi, RequestType: requestType}
}

func NewGxCCAnswer(resultCode uint32) *GxCreditControlAnswer {
	return &GxCreditControlAnswer{ResultCode: resultCode}
}

func (m *GxCreditControlAnswer) initializeRuleInstallsIfNil() {
	if m.RuleInstalls == nil {
		m.RuleInstalls = &RuleInstalls{}
	}
}

func (m *GxCreditControlAnswer) SetUsageMonitorInfo(monitor *UsageMonitoringInformation) *GxCreditControlAnswer {
	if m.UsageMonitoringInfos == nil {
		m.UsageMonitoringInfos = []*UsageMonitoringInformation{}
	}
	m.UsageMonitoringInfos = append(m.UsageMonitoringInfos, monitor)
	return m
}

func (m *GxCreditControlAnswer) SetStaticRuleInstalls(ruleIDs, baseNames []string) *GxCreditControlAnswer {
	m.initializeRuleInstallsIfNil()
	m.RuleInstalls.RuleNames = ruleIDs
	m.RuleInstalls.RuleBaseNames = baseNames
	return m
}

func (m *GxCreditControlAnswer) SetRuleActivationTime(activationTime *timestamp.Timestamp) *GxCreditControlAnswer {
	m.initializeRuleInstallsIfNil()
	m.RuleInstalls.ActivationTime = activationTime
	return m
}

func (m *GxCreditControlAnswer) SetRuleDeactivationTime(deactivationTime *timestamp.Timestamp) *GxCreditControlAnswer {
	m.initializeRuleInstallsIfNil()
	m.RuleInstalls.DeactivationTime = deactivationTime
	return m
}

func (m *GxCreditControlAnswer) SetDynamicRuleInstall(rule *RuleDefinition) *GxCreditControlAnswer {
	m.initializeRuleInstallsIfNil()
	m.RuleInstalls.RuleDefinitions = append(m.RuleInstalls.RuleDefinitions, rule)
	return m
}

func (m *GxCreditControlAnswer) SetStaticRuleRemovals(rulesIDs, baseNames []string) *GxCreditControlAnswer {
	if m.RuleRemovals == nil {
		m.RuleRemovals = &RuleRemovals{}
	}
	m.RuleRemovals.RuleNames = rulesIDs
	m.RuleRemovals.RuleBaseNames = baseNames
	return m
}

func (m *GxCreditControlAnswer) SetEventTriggers(eventTriggers []uint32) *GxCreditControlAnswer {
	m.EventTriggers = eventTriggers
	return m
}

func (m *GxCreditControlAnswer) SetRevalidationTime(revalidationTime *timestamp.Timestamp) *GxCreditControlAnswer {
	m.RevalidationTime = revalidationTime
	return m
}

func (m *GxCreditControlRequest) SetUsageMonitorReport(report *UsageMonitoringInformation) *GxCreditControlRequest {
	if m.UsageMonitoringReports == nil {
		m.UsageMonitoringReports = []*UsageMonitoringInformation{}
	}
	m.UsageMonitoringReports = append(m.UsageMonitoringReports, report)
	return m
}

func (m *GxCreditControlRequest) SetUsageReportDelta(delta uint64) *GxCreditControlRequest {
	m.UsageReportDelta = delta
	return m
}

func (m *GxCreditControlRequest) SetEventTrigger(eventTrigger int32) *GxCreditControlRequest {
	m.EventTrigger = &wrappers.Int32Value{Value: eventTrigger}
	return m
}

func NewGyCreditControlExpectation() *GyCreditControlExpectation {
	return &GyCreditControlExpectation{}
}

func (m *GyCreditControlExpectation) Expect(ccr *GyCreditControlRequest) *GyCreditControlExpectation {
	m.ExpectedRequest = ccr
	return m
}

func (m *GyCreditControlExpectation) Return(cca *GyCreditControlAnswer) *GyCreditControlExpectation {
	m.Answer = cca
	return m
}

func NewGyCCRequest(imsi string, requestType CCRequestType) *GyCreditControlRequest {
	return &GyCreditControlRequest{Imsi: imsi, RequestType: requestType}
}

func NewGyCCAnswer(resultCode uint32) *GyCreditControlAnswer {
	return &GyCreditControlAnswer{ResultCode: resultCode}
}

func (m *GyCreditControlAnswer) SetQuotaGrant(grant *QuotaGrant) *GyCreditControlAnswer {
	if m.QuotaGrants == nil {
		m.QuotaGrants = []*QuotaGrant{}
	}
	m.QuotaGrants = append(m.QuotaGrants, grant)
	return m
}

func (m *GyCreditControlRequest) SetRequestNumber(requestNumber int32) *GyCreditControlRequest {
	m.RequestNumber = &wrappers.Int32Value{Value: requestNumber}
	return m
}

func (m *GyCreditControlAnswer) SetLinkFailure(linkFailure bool) *GyCreditControlAnswer {
	m.LinkFailure = linkFailure
	return m
}

func (m *GyCreditControlRequest) SetMSCC(mscc *MultipleServicesCreditControl) *GyCreditControlRequest {
	if m.Mscc == nil {
		m.Mscc = []*MultipleServicesCreditControl{}
	}
	m.Mscc = append(m.Mscc, mscc)
	return m
}

func (m *GyCreditControlRequest) SetMSCCDelta(delta uint64) *GyCreditControlRequest {
	m.UsageReportDelta = delta
	return m
}
