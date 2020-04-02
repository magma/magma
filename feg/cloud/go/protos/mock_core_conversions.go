/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package protos

import "github.com/golang/protobuf/ptypes/timestamp"

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

func NewGxCCRequest(imsi string, requestType CCRequestType, requestNumber uint32) *GxCreditControlRequest {
	return &GxCreditControlRequest{Imsi: imsi, RequestType: requestType, RequestNumber: requestNumber}
}

func NewGxCCAnswer(resultCode uint32) *GxCreditControlAnswer {
	return &GxCreditControlAnswer{ResultCode: resultCode}
}

func (m *GxCreditControlAnswer) SetUsageMonitorInfos(monitors []*UsageMonitoringInformation) *GxCreditControlAnswer {
	m.UsageMonitoringInfos = monitors
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

func (m *GxCreditControlAnswer) SetDynamicRuleInstalls(rules []*RuleDefinition) *GxCreditControlAnswer {
	m.initializeRuleInstallsIfNil()
	m.RuleInstalls.RuleDefinitions = rules
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

func (m *GxCreditControlRequest) SetUsageMonitorReports(reports []*UsageMonitoringInformation) *GxCreditControlRequest {
	m.UsageMonitoringReports = reports
	return m
}

func (m *GxCreditControlRequest) SetUsageReportDelta(delta uint64) *GxCreditControlRequest {
	m.UsageReportDelta = delta
	return m
}

func (m *GxCreditControlAnswer) initializeRuleInstallsIfNil() {
	if m.RuleInstalls == nil {
		m.RuleInstalls = &RuleInstalls{}
	}
}
