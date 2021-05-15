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

package gx_test

import (
	"sort"
	"strings"
	"testing"
	"time"

	"magma/feg/gateway/policydb/mocks"
	"magma/feg/gateway/services/session_proxy/credit_control"
	"magma/feg/gateway/services/session_proxy/credit_control/gx"
	"magma/lte/cloud/go/protos"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
)

func TestReAuthRequest_ToProto(t *testing.T) {
	// Check nil, 1-element, multiple elements, and empty arrays
	monitoringKey := []byte("monitor")
	monitoringKey2 := []byte("monitor2")
	monitoringKey3 := []byte("monitor3")
	bearerID := "bearer1"
	var ratingGroup uint32 = 42
	var totalOctets uint64 = 2048
	var monitorSupport0 = gx.UsageMonitoringDisabled
	var monitorReport0 = gx.UsageMonitoringReport
	var qci uint32 = 1
	var monitoringLevel = gx.SessionLevel
	currentTime := time.Now()
	protoTimestamp, err := ptypes.TimestampProto(currentTime)
	assert.NoError(t, err)
	in := &gx.PolicyReAuthRequest{
		SessionID: "IMSI001010000000001-1234",
		RulesToRemove: []*gx.RuleRemoveAVP{
			{RuleNames: []string{"remove1", "remove2"}, RuleBaseNames: []string{"baseRemove1"}},
			{RuleNames: nil, RuleBaseNames: nil},
			{RuleNames: []string{"remove3"}, RuleBaseNames: []string{}},
			{RuleNames: []string{}, RuleBaseNames: []string{"baseRemove2", "baseRemove3"}},
		},
		RulesToInstall: []*gx.RuleInstallAVP{
			{RuleNames: []string{"install1", "install2"}, RuleBaseNames: []string{"baseInstall1"}, RuleDefinitions: nil},
			{
				RuleNames:     nil,
				RuleBaseNames: nil,
				RuleDefinitions: []*gx.RuleDefinition{
					{RuleName: "dynamic1",
						MonitoringKey: monitoringKey,
						Precedence:    100,
						RatingGroup:   &ratingGroup},
				},
			},
			{RuleNames: []string{"install3"}, RuleBaseNames: []string{}},
			{RuleNames: []string{}, RuleBaseNames: []string{"baseInstall2", "baseInstall3"}},
		},
		EventTriggers:    []gx.EventTrigger{gx.UsageReportTrigger, gx.RevalidationTimeout},
		RevalidationTime: &currentTime,
		UsageMonitors: []*gx.UsageMonitoringInfo{
			{
				MonitoringKey: monitoringKey,
				GrantedServiceUnit: &credit_control.GrantedServiceUnit{
					TotalOctets: &totalOctets,
				},
				Level: monitoringLevel,
			},
			{
				MonitoringKey: monitoringKey2,
				GrantedServiceUnit: &credit_control.GrantedServiceUnit{
					InputOctets:  &totalOctets,
					OutputOctets: &totalOctets,
					TotalOctets:  &totalOctets,
				},
				Level:   monitoringLevel,
				Support: &monitorSupport0,
			},
			{
				MonitoringKey:      monitoringKey3,
				GrantedServiceUnit: nil,
				Level:              monitoringLevel,
				Report:             &monitorReport0,
			},
		},
		Qos: &gx.QosInformation{
			BearerIdentifier: bearerID,
			Qci:              &qci,
		},
	}
	policyClient := &mocks.PolicyDBClient{}
	policyClient.On("GetRuleIDsForBaseNames", []string{"baseRemove1", "baseRemove2", "baseRemove3"}).
		Return([]string{"remove42", "remove43", "remove44"})
	policyClient.On("GetRuleIDsForBaseNames", []string{"baseInstall1"}).
		Return([]string{})
	policyClient.On("GetRuleIDsForBaseNames", []string{"baseInstall2", "baseInstall3"}).
		Return([]string{"install42", "install43"})

	actual := in.ToProto("IMSI001010000000001", "magma;1234;1234;IMSI001010000000001", policyClient)
	expected := &protos.PolicyReAuthRequest{
		SessionId:     "magma;1234;1234;IMSI001010000000001",
		Imsi:          "IMSI001010000000001",
		RulesToRemove: []string{"remove1", "remove2", "remove3", "remove42", "remove43", "remove44"},
		RulesToInstall: []*protos.StaticRuleInstall{
			{
				RuleId: "install1",
			},
			{
				RuleId: "install2",
			},
			{
				RuleId: "install3",
			},
			{
				RuleId: "install42",
			},
			{
				RuleId: "install43",
			},
		},
		DynamicRulesToInstall: []*protos.DynamicRuleInstall{
			{
				PolicyRule: &protos.PolicyRule{
					Id:            "dynamic1",
					RatingGroup:   42,
					MonitoringKey: monitoringKey,
					Priority:      100,
					TrackingType:  protos.PolicyRule_OCS_AND_PCRF,
					Redirect:      &protos.RedirectInformation{},
				},
			},
		},
		EventTriggers: []protos.EventTrigger{
			protos.EventTrigger_UNSUPPORTED,
			protos.EventTrigger_REVALIDATION_TIMEOUT,
		},
		RevalidationTime: protoTimestamp,
		UsageMonitoringCredits: []*protos.UsageMonitoringCredit{
			{
				Action:        protos.UsageMonitoringCredit_CONTINUE,
				MonitoringKey: monitoringKey,
				GrantedUnits: &protos.GrantedUnits{
					Total: &protos.CreditUnit{IsValid: true, Volume: totalOctets},
					Tx:    &protos.CreditUnit{IsValid: false},
					Rx:    &protos.CreditUnit{IsValid: false},
				},
				Level: protos.MonitoringLevel(monitoringLevel),
			},
			{
				Action:        protos.UsageMonitoringCredit_DISABLE,
				MonitoringKey: monitoringKey2,
				GrantedUnits: &protos.GrantedUnits{
					Total: &protos.CreditUnit{IsValid: true, Volume: totalOctets},
					Tx:    &protos.CreditUnit{IsValid: true, Volume: totalOctets},
					Rx:    &protos.CreditUnit{IsValid: true, Volume: totalOctets},
				},
				Level: protos.MonitoringLevel(monitoringLevel),
			},
			{
				Action: protos.UsageMonitoringCredit_FORCE,
				GrantedUnits: &protos.GrantedUnits{
					Total: &protos.CreditUnit{IsValid: false},
					Tx:    &protos.CreditUnit{IsValid: false},
					Rx:    &protos.CreditUnit{IsValid: false},
				},
				MonitoringKey: monitoringKey3,
				Level:         protos.MonitoringLevel(monitoringLevel),
			},
		},
		QosInfo: &protos.QoSInformation{
			BearerId: bearerID,
			Qci:      protos.QCI_QCI_1,
		},
	}
	assert.Equal(t, expected, actual)
	policyClient.AssertExpectations(t)
}

func TestReAuthAnswer_FromProto(t *testing.T) {
	in := &protos.PolicyReAuthAnswer{
		SessionId: "foo",
		FailedRules: map[string]protos.PolicyReAuthAnswer_FailureCode{
			"bar": protos.PolicyReAuthAnswer_CM_AUTHORIZATION_REJECTED,
			"baz": protos.PolicyReAuthAnswer_AN_GW_FAILED,
		},
	}
	actual := (&gx.PolicyReAuthAnswer{}).FromProto("sesh", in)

	// sort the rules so we get a deterministic test
	sortFun := func(i, j int) bool {
		first := actual.RuleReports[i]
		second := actual.RuleReports[j]

		concattedFirst := strings.Join(first.RuleNames, "") + strings.Join(first.RuleBaseNames, "")
		concattedSecond := strings.Join(second.RuleNames, "") + strings.Join(second.RuleBaseNames, "")
		return concattedFirst < concattedSecond
	}
	sort.Slice(actual.RuleReports, sortFun)

	expected := &gx.PolicyReAuthAnswer{
		SessionID:  "sesh",
		ResultCode: diam.Success,
		RuleReports: []*gx.ChargingRuleReport{
			{RuleNames: []string{"bar"}, FailureCode: gx.CMAuthorizationRejected},
			{RuleNames: []string{"baz"}, FailureCode: gx.ANGWFailed},
		},
	}
	assert.Equal(t, expected, actual)
}

func TestRuleDefinition_ToProto(t *testing.T) {
	// Check nil, 1-element, multiple elements, and empty arrays
	monitoringKey := []byte("monitor")
	var ratingGroup uint32 = 10
	var ruleOut *protos.PolicyRule = nil

	ruleOut = (&gx.RuleDefinition{
		RuleName:      "rgonly",
		MonitoringKey: nil,
		RatingGroup:   &ratingGroup,
	}).ToProto()
	assert.Equal(t, []byte(nil), ruleOut.MonitoringKey)
	assert.Equal(t, uint32(10), ruleOut.RatingGroup)
	assert.Equal(t, protos.PolicyRule_ONLY_OCS, ruleOut.TrackingType)

	ruleOut = (&gx.RuleDefinition{
		RuleName:      "mkonly",
		MonitoringKey: monitoringKey,
		RatingGroup:   nil,
	}).ToProto()
	assert.Equal(t, []byte("monitor"), ruleOut.MonitoringKey)
	assert.Equal(t, uint32(0), ruleOut.RatingGroup)
	assert.Equal(t, protos.PolicyRule_ONLY_PCRF, ruleOut.TrackingType)

	ruleOut = (&gx.RuleDefinition{
		RuleName:      "both",
		MonitoringKey: monitoringKey,
		RatingGroup:   &ratingGroup,
	}).ToProto()
	assert.Equal(t, []byte("monitor"), ruleOut.MonitoringKey)
	assert.Equal(t, uint32(10), ruleOut.RatingGroup)
	assert.Equal(t, protos.PolicyRule_OCS_AND_PCRF, ruleOut.TrackingType)

	ruleOut = (&gx.RuleDefinition{
		RuleName:      "neither",
		MonitoringKey: nil,
		RatingGroup:   nil,
	}).ToProto()
	assert.Equal(t, []byte(nil), ruleOut.MonitoringKey)
	assert.Equal(t, uint32(0), ruleOut.RatingGroup)
	assert.Equal(t, protos.PolicyRule_NO_TRACKING, ruleOut.TrackingType)
}
