/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package mock_pcrf

import (
	"fmt"

	"magma/feg/cloud/go/protos"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/golang/glog"
)

type GxExpectation struct {
	*protos.GxCreditControlExpectation
}

type GxAnswer struct {
	*protos.GxCreditControlAnswer
}

type requestPK struct {
	imsi        string
	requestType protos.CCRequestType
}

func (r requestPK) String() string {
	return fmt.Sprintf("Imsi: %v, Type: %v", r.imsi, r.requestType)
}

func (e GxExpectation) GetAnswer() interface{} {
	return GxAnswer{e.Answer}
}

func (e GxExpectation) DoesMatch(message interface{}) error {
	expected := e.ExpectedRequest
	ccr := message.(ccrMessage)
	expectedPK := requestPK{imsi: expected.Imsi, requestType: expected.RequestType}
	actualImsi, _ := ccr.GetIMSI()
	actualPK := requestPK{imsi: actualImsi, requestType: protos.CCRequestType(ccr.RequestType)}
	// For better readability of errors, we will check for the IMSI and the request type first.
	if expectedPK != actualPK {
		return fmt.Errorf("Expected: %v, Received: %v", expectedPK, actualPK)
	}
	expectedUsageReports := expected.GetUsageMonitoringReports()
	if !compareUsageMonitorsAgainstExpected(ccr.UsageMonitors, expectedUsageReports, expected.GetUsageReportDelta()) {
		return fmt.Errorf("For Request=%v, Expected: %v, Received: %v", actualPK, expectedUsageReports, ccr.UsageMonitors)
	}
	return nil
}

func (answer GxAnswer) toAVPs() ([]*diam.AVP, uint32) {
	avps := []*diam.AVP{}
	ruleInstalls := answer.GetRuleInstalls()
	if ruleInstalls != nil {
		ruleInstallAVPs := toRuleInstallAVPs(
			answer.RuleInstalls.GetRuleNames(),
			answer.RuleInstalls.GetRuleBaseNames(),
			answer.RuleInstalls.GetRuleDefinitions())
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
	if monitorInstalls != nil {
		for _, monitor := range monitorInstalls {
			octets := monitor.GetOctets()
			if octets == nil {
				glog.Errorf("Monitor Octets is nil, skipping.")
				continue
			}
			avps = append(avps, toUsageMonitoringInfoAVP(string(monitor.MonitoringKey), octets, monitor.MonitoringLevel))
		}
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
		if !equalWithinDelta(actualMonitor.UsedServiceUnit.TotalOctets, expectedMonitor.GetOctets().GetTotalOctets(), delta) {
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

func equalWithinDelta(a, b, delta uint64) bool {
	if b >= a && b-a <= delta {
		return true
	}
	if a >= b && a-b <= delta {
		return true
	}
	return false
}
