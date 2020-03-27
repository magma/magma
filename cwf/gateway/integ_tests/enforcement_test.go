/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package integ_tests

import (
	"fmt"
	"testing"
	"time"

	cwfprotos "magma/cwf/cloud/go/protos"
	"magma/feg/cloud/go/protos"
	"magma/lte/cloud/go/plugin/models"

	"github.com/emakeev/go-diameter/diam"
	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
)

const (
	KiloBytes = 1024
	MegaBytes = 1024 * KiloBytes
	Buffer    = 50 * KiloBytes
)

// - Set an expectation for a  CCR-I to be sent up to PCRF, to which it will
//   respond with a rule install (usage-enforcement-static-pass-all), 250KB of
//   quota.
//   Generate traffic and assert the CCR-I is received.
// - Set an expectation for a CCR-U with >80% of data usage to be sent up to
// 	 PCRF, to which it will response with more quota.
//   Generate traffic and assert the CCR-U is received.
// - Generate traffic to put traffic through the newly installed rule.
//   Assert that there's > 0 data usage in the rule.
// - Expect a CCR-T, trigger a UE disconnect, and assert the CCR-T is received.
func TestUsageReportEnforcement(t *testing.T) {
	fmt.Println("\nRunning TestUsageReportEnforcement...")
	tr := NewTestRunner(t)
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)
	assert.NoError(t, usePCRFMockDriver())
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearPCRFMockDriver())
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)
	imsi := ues[0].GetImsi()

	err = ruleManager.AddStaticPassAllToDB("usage-enforcement-static-pass-all", "mkey1", 0, models.PolicyRuleTrackingTypeONLYPCRF, 3)
	assert.NoError(t, err)
	tr.WaitForPoliciesToSync()

	usageMonitorInfo := []*protos.UsageMonitoringInformation{
		{
			MonitoringLevel: protos.MonitoringLevel_RuleLevel,
			MonitoringKey:   []byte("mkey1"),
			Octets:          &protos.Octets{TotalOctets: 250 * KiloBytes},
		},
	}

	initRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_INITIAL, 1)
	initAnswer := protos.NewGxCCAnswer(diam.Success).
		SetStaticRuleInstalls([]string{"usage-enforcement-static-pass-all"}, []string{}).
		SetUsageMonitorInfos(usageMonitorInfo)
	initExpectation := protos.NewGxCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	// We expect an update request with some usage update (probably around 80-100% of the given quota)
	updateRequest1 := protos.NewGxCCRequest(imsi, protos.CCRequestType_UPDATE, 2).
		SetUsageMonitorReports(usageMonitorInfo).
		SetUsageReportDelta(250 * KiloBytes * 0.2)
	updateAnswer1 := protos.NewGxCCAnswer(diam.Success).SetUsageMonitorInfos(usageMonitorInfo)
	updateExpectation1 := protos.NewGxCreditControlExpectation().Expect(updateRequest1).Return(updateAnswer1)
	expectations := []*protos.GxCreditControlExpectation{initExpectation, updateExpectation1}
	// On unexpected requests, just return the default update answer
	assert.NoError(t, setPCRFExpectations(expectations, updateAnswer1))

	tr.AuthenticateAndAssertSuccess(imsi)

	req := &cwfprotos.GenTrafficRequest{Imsi: imsi, Volume: &wrappers.StringValue{Value: *swag.String("500K")}}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	// Assert that enforcement_stats rules are properly installed and the right
	// amount of data was passed through
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	record := recordsBySubID["IMSI"+imsi]["usage-enforcement-static-pass-all"]
	assert.NotNil(t, record, fmt.Sprintf("No policy usage record for imsi: %v", imsi))
	if record != nil {
		// We should not be seeing > 1024k data here
		assert.True(t, record.BytesTx > uint64(0), fmt.Sprintf("%s did not pass any data", record.RuleId))
		assert.True(t, record.BytesTx <= uint64(500*KiloBytes+Buffer), fmt.Sprintf("policy usage: %v", record))
	}

	// Assert that reasonable CCR-I and at least one CCR-U were sent up to the PCRF
	resultByIndex, errByIndex, err := getAssertExpectationsResult()
	assert.NoError(t, err)
	assert.Empty(t, errByIndex)
	expectedResult := []*protos.ExpectationResult{
		{ExpectationIndex: 0, ExpectationMet: true},
		{ExpectationIndex: 1, ExpectationMet: true},
	}
	assert.ElementsMatch(t, expectedResult, resultByIndex)

	// When we initiate a UE disconnect, we expect a terminate request to go up
	terminateRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_TERMINATION, 3)
	terminateAnswer := protos.NewGxCCAnswer(diam.Success)
	terminateExpectation := protos.NewGxCreditControlExpectation().Expect(terminateRequest).Return(terminateAnswer)
	expectations = []*protos.GxCreditControlExpectation{terminateExpectation}
	assert.NoError(t, setPCRFExpectations(expectations, nil))

	_, err = tr.Disconnect(imsi)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	// Assert that we saw a Terminate request
	resultByIndex, errByIndex, err = getAssertExpectationsResult()
	assert.NoError(t, err)
	assert.Empty(t, errByIndex)
	expectedResult = []*protos.ExpectationResult{
		{ExpectationIndex: 0, ExpectationMet: true},
	}
	assert.ElementsMatch(t, expectedResult, resultByIndex)
}

// - Set an expectation for a  CCR-I to be sent up to PCRF, to which it will
//   respond with a rule install (static-pass-all-1).
//   Generate traffic and assert the CCR-I is received.
// - Set an expectation for a CCR-U to be sent up to PCRF, to which it will
//   respond with a rule removal (static-pass-all-1) and rule install (static-pass-all-2).
//   Generate traffic and assert the CCR-U is received.
// - Generate traffic to put traffic through the newly installed rule.
//   Assert that there's > 0 data usage in the rule.
func TestMidSessionRuleRemovalWithCCA_U(t *testing.T) {
	fmt.Println("\nRunning TestMidSessionRuleRemovalWithCCA_U...")

	tr := NewTestRunner(t)
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)
	assert.NoError(t, usePCRFMockDriver())
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearPCRFMockDriver())
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)
	imsi := ues[0].GetImsi()

	err = ruleManager.AddStaticPassAllToDB("static-pass-all-1", "mkey1", 0, models.PolicyRuleTrackingTypeONLYPCRF, 100)
	assert.NoError(t, err)
	err = ruleManager.AddStaticPassAllToDB("static-pass-all-2", "mkey2", 0, models.PolicyRuleTrackingTypeONLYPCRF, 150)
	assert.NoError(t, err)
	err = ruleManager.AddStaticPassAllToDB("static-pass-all-3", "mkey2", 0, models.PolicyRuleTrackingTypeONLYPCRF, 200)
	assert.NoError(t, err)
	err = ruleManager.AddBaseNameMappingToDB("base-1", []string{"static-pass-all-3"})
	tr.WaitForPoliciesToSync()

	usageMonitorInfo := []*protos.UsageMonitoringInformation{
		{
			MonitoringLevel: protos.MonitoringLevel_RuleLevel,
			MonitoringKey:   []byte("mkey1"),
			Octets:          &protos.Octets{TotalOctets: 250 * KiloBytes},
		},
	}

	initRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_INITIAL, 1)
	initAnswer := protos.NewGxCCAnswer(diam.Success).
		SetStaticRuleInstalls([]string{"static-pass-all-1"}, []string{"base-1"}).
		SetUsageMonitorInfos(usageMonitorInfo)
	initExpectation := protos.NewGxCreditControlExpectation().Expect(initRequest).Return(initAnswer)
	// Remove the high priority Rule
	defaultUpdateAnswer := protos.NewGxCCAnswer(diam.Success).SetUsageMonitorInfos(usageMonitorInfo)
	expectations := []*protos.GxCreditControlExpectation{initExpectation}
	// On unexpected requests, just return some quota
	assert.NoError(t, setPCRFExpectations(expectations, defaultUpdateAnswer))

	tr.AuthenticateAndAssertSuccess(imsi)

	req := &cwfprotos.GenTrafficRequest{Imsi: imsi, Volume: &wrappers.StringValue{Value: "250K"}}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)

	// At this point both static-pass-all-1 & static-pass-all-3 are installed.
	// Since static-pass-all-1 has higher precedence, it will get hit.
	tr.WaitForEnforcementStatsToSync()

	// Assert that enforcement_stats rules are properly installed and the right
	// amount of data was passed through
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	record1 := recordsBySubID[prependIMSIPrefix(imsi)]["static-pass-all-1"]
	if record1 != nil {
		assert.True(t, record1.BytesTx > uint64(0), fmt.Sprintf("%s did not pass any data", record1.RuleId))
	}
	assert.NotNil(t, record1, fmt.Sprintf("No policy usage record for imsi: %v rule=static-pass-all-1", imsi))

	// Assert that a CCR-I was sent up to the PCRF
	resultByIndex, errByIndex, err := getAssertExpectationsResult()
	assert.NoError(t, err)
	assert.Empty(t, errByIndex)
	expectedResult := []*protos.ExpectationResult{
		{ExpectationIndex: 0, ExpectationMet: true},
	}
	assert.ElementsMatch(t, expectedResult, resultByIndex)

	updateRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_UPDATE, 2).
		SetUsageMonitorReports(usageMonitorInfo).
		SetUsageReportDelta(250 * KiloBytes * 0.5)
	updateAnswer := protos.NewGxCCAnswer(diam.Success).SetUsageMonitorInfos(usageMonitorInfo).
		SetStaticRuleInstalls([]string{"static-pass-all-2"}, []string{}).
		SetStaticRuleRemovals([]string{"static-pass-all-1"}, []string{"base-1"})
	updateExpectation := protos.NewGxCreditControlExpectation().Expect(updateRequest).Return(updateAnswer)
	expectations = []*protos.GxCreditControlExpectation{updateExpectation}
	// On unexpected requests, just return some quota
	assert.NoError(t, setPCRFExpectations(expectations, defaultUpdateAnswer))

	fmt.Println("Generating traffic again to trigger a CCR/A-U so that 'static-pass-all-1' gets removed")
	// Generate traffic to trigger the CCR-U so that the rule removal/install happens
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	fmt.Println("Generating traffic again to put data through static-pass-all-2")
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	// Assert that we sent back a CCA-Update with RuleRemovals
	resultByIndex, errByIndex, err = getAssertExpectationsResult()
	assert.NoError(t, err)
	assert.Empty(t, errByIndex)
	expectedResult = []*protos.ExpectationResult{
		{ExpectationIndex: 0, ExpectationMet: true},
	}
	assert.ElementsMatch(t, expectedResult, resultByIndex)

	recordsBySubID, err = tr.GetPolicyUsage()
	assert.NoError(t, err)
	assert.NotNil(t, recordsBySubID[prependIMSIPrefix(imsi)]["static-pass-all-2"], fmt.Sprintf("No policy usage record for imsi: %v rule=static-pass-all-2", imsi))

	_, err = tr.Disconnect(imsi)
	assert.NoError(t, err)
	time.Sleep(3 * time.Second)
}

// - Set an expectation for a  CCR-I to be sent up to PCRF, to which it will
//   respond with a rule install (static-pass-all-1).
// - Set an expectation for a CCR-U to be sent up to PCRF, to which it will
//   respond with a rule install (static-pass-all-2), with activation (now + X sec)
//   and deactivation time (activation + Y sec) specified.
// - Generate traffic to trigger a CCR-U. Check policy usage and assert
//   static-pass-all-2 is not installed.
// - Sleep for X seconds and check policy usage again. Assert that
//   static-pass-all-2 is installed.
// - Sleep for Y seconds and check policy usage again. Assert that
//   static-pass-all-2 is uninstalled.
// Note: things might get weird if there are clock skews
func TestRuleInstallTime(t *testing.T) {
	fmt.Println("\nRunning TestRuleInstallTime...")

	tr := NewTestRunner(t)
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)
	assert.NoError(t, usePCRFMockDriver())
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearPCRFMockDriver())
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)
	imsi := ues[0].GetImsi()

	err = ruleManager.AddStaticPassAllToDB("static-pass-all-1", "mkey1", 0, models.PolicyRuleTrackingTypeONLYPCRF, 200)
	assert.NoError(t, err)
	err = ruleManager.AddStaticPassAllToDB("static-pass-all-2", "mkey1", 0, models.PolicyRuleTrackingTypeONLYPCRF, 100)
	assert.NoError(t, err)
	tr.WaitForPoliciesToSync()

	usageMonitorInfo := []*protos.UsageMonitoringInformation{
		{
			MonitoringLevel: protos.MonitoringLevel_RuleLevel,
			MonitoringKey:   []byte("mkey1"),
			Octets:          &protos.Octets{TotalOctets: 250 * KiloBytes},
		},
	}
	initRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_INITIAL, 1)
	initAnswer := protos.NewGxCCAnswer(diam.Success).
		SetStaticRuleInstalls([]string{"static-pass-all-1"}, nil).
		SetUsageMonitorInfos(usageMonitorInfo)
	initExpectation := protos.NewGxCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	now := time.Now().Round(1 * time.Second)
	timeUntilActivation := 8 * time.Second
	activation := now.Add(timeUntilActivation)
	pActivation, err := ptypes.TimestampProto(activation)
	assert.NoError(t, err)
	timeUntilDeactivation := 8 * time.Second
	deactivation := activation.Add(timeUntilDeactivation)
	pDeactivation, err := ptypes.TimestampProto(deactivation)
	assert.NoError(t, err)

	updateRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_UPDATE, 2)
	updateAnswer := protos.NewGxCCAnswer(diam.Success).
		SetUsageMonitorInfos(usageMonitorInfo).
		SetStaticRuleInstalls([]string{"static-pass-all-2"}, nil).
		SetRuleActivationTime(pActivation).
		SetRuleDeactivationTime(pDeactivation)
	updateExpectation := protos.NewGxCreditControlExpectation().Expect(updateRequest).Return(updateAnswer)
	defaultUpdateAnswer := protos.NewGxCCAnswer(diam.Success).SetUsageMonitorInfos(usageMonitorInfo)
	expectations := []*protos.GxCreditControlExpectation{initExpectation, updateExpectation}
	// On unexpected requests, just return some quota
	assert.NoError(t, setPCRFExpectations(expectations, defaultUpdateAnswer))

	tr.AuthenticateAndAssertSuccess(imsi)

	// Generate over the given quota
	req := &cwfprotos.GenTrafficRequest{Imsi: imsi, Volume: &wrappers.StringValue{Value: "300K"}}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()
	recordsBySubID, err := tr.GetPolicyUsage()
	// only static-pass-all-1 should be installed
	assert.NotNil(t, recordsBySubID[prependIMSIPrefix(imsi)]["static-pass-all-1"])
	assert.Nil(t, recordsBySubID[prependIMSIPrefix(imsi)]["static-pass-all-2"])

	fmt.Printf("Waiting %v for rule activation\n", timeUntilActivation)
	time.Sleep(timeUntilActivation)
	recordsBySubID, err = tr.GetPolicyUsage()
	// both rules should exist
	assert.NotNil(t, recordsBySubID[prependIMSIPrefix(imsi)]["static-pass-all-1"])
	assert.NotNil(t, recordsBySubID[prependIMSIPrefix(imsi)]["static-pass-all-2"])

	fmt.Printf("Waiting %v for rule deactivation\n", timeUntilDeactivation)
	time.Sleep(timeUntilDeactivation)
	recordsBySubID, err = tr.GetPolicyUsage()
	// static-pass-all-2 should be gone
	assert.NotNil(t, recordsBySubID[prependIMSIPrefix(imsi)]["static-pass-all-1"])
	assert.Nil(t, recordsBySubID[prependIMSIPrefix(imsi)]["static-pass-all-2"])

	resultByIndex, errByIndex, err := getAssertExpectationsResult()
	assert.NoError(t, err)
	assert.Empty(t, errByIndex)
	expectedResult := []*protos.ExpectationResult{
		{ExpectationIndex: 0, ExpectationMet: true},
		{ExpectationIndex: 1, ExpectationMet: true},
	}
	assert.ElementsMatch(t, expectedResult, resultByIndex)

	_, err = tr.Disconnect(imsi)
	assert.NoError(t, err)
	time.Sleep(3 * time.Second)
}
