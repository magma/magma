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
	cwfprotos "magma/cwf/cloud/go/protos"
	"magma/feg/cloud/go/protos"
	"magma/lte/cloud/go/plugin/models"
	"testing"
	"time"

	"github.com/emakeev/go-diameter/diam"
	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/ptypes/wrappers"

	"github.com/stretchr/testify/assert"
)

const (
	KiloBytes = 1024
	MegaBytes = 1024 * KiloBytes
	Buffer    = 50 * KiloBytes
)

func TestBasicUplinkTrafficWithEnforcement(t *testing.T) {
	fmt.Println("\nRunning TestBasicUplinkTrafficWithEnforcement...")
	tr := NewTestRunner()
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)
	assert.NoError(t, usePCRFMockDriver())

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)
	imsi := ues[0].GetImsi()

	err = ruleManager.AddStaticPassAllToDB("ul-enforcement-static-pass-all", "mkey1", 0, models.PolicyRuleTrackingTypeONLYPCRF, 3)
	assert.NoError(t, err)

	usageMonitorInfo := []*protos.UsageMonitoringInformation{
		{
			MonitoringLevel: protos.MonitoringLevel_RuleLevel,
			MonitoringKey:   []byte("mkey1"),
			Octets:          &protos.Octets{TotalOctets: 250 * KiloBytes},
		},
	}

	initRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_INITIAL, 1)
	initAnswer := protos.NewGxCCAnswer(diam.Success).
		SetStaticRuleInstalls([]string{"ul-enforcement-static-pass-all"}, []string{}).
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

	// wait for the rules to be synced into sessiond
	time.Sleep(1 * time.Second)

	tr.AuthenticateAndAssertSuccess(t, imsi)

	req := &cwfprotos.GenTrafficRequest{Imsi: imsi, Volume: &wrappers.StringValue{Value: *swag.String("500K")}}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)

	// Wait for the traffic to go through
	time.Sleep(6 * time.Second)

	// Assert that enforcement_stats rules are properly installed and the right
	// amount of data was passed through
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	record := recordsBySubID["IMSI"+imsi]["ul-enforcement-static-pass-all"]
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
	time.Sleep(3 * time.Second)

	// Assert that we saw a Terminate request
	resultByIndex, errByIndex, err = getAssertExpectationsResult()
	assert.NoError(t, err)
	assert.Empty(t, errByIndex)
	expectedResult = []*protos.ExpectationResult{
		{ExpectationIndex: 0, ExpectationMet: true},
	}
	assert.ElementsMatch(t, expectedResult, resultByIndex)

	// Clear hss, ocs, and pcrf
	assert.NoError(t, clearPCRFMockDriver())
	assert.NoError(t, ruleManager.RemoveInstalledRules())
	assert.NoError(t, tr.CleanUp())
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

	tr := NewTestRunner()
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)
	assert.NoError(t, usePCRFMockDriver())

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

	// wait for the rules to be synced into sessiond
	time.Sleep(1 * time.Second)

	tr.AuthenticateAndAssertSuccess(t, imsi)

	req := &cwfprotos.GenTrafficRequest{Imsi: imsi, Volume: &wrappers.StringValue{Value: "250K"}}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)

	// At this point both static-pass-all-1 & static-pass-all-3 are installed.
	// Since static-pass-all-1 has higher precedence, it will get hit.
	// Wait for some traffic to go through
	time.Sleep(1 * time.Second)

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
	time.Sleep(3 * time.Second)

	fmt.Println("Generating traffic again to put data through static-pass-all-2")
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	// Wait for some traffic to go through
	time.Sleep(1 * time.Second)

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
	record2 := recordsBySubID[prependIMSIPrefix(imsi)]["static-pass-all-2"]
	assert.NotNil(t, record1, fmt.Sprintf("No policy usage record for imsi: %v rule=static-pass-all-2", imsi))
	if record2 != nil {
		// This rule should have passed some traffic
		assert.True(t, record2.BytesTx > 0, fmt.Sprintf("%s did not pass any data", record2.RuleId))
	}

	_, err = tr.Disconnect(imsi)
	assert.NoError(t, err)
	time.Sleep(3 * time.Second)

	// Clear hss, ocs, and pcrf
	assert.NoError(t, clearPCRFMockDriver())
	assert.NoError(t, ruleManager.RemoveInstalledRules())
	assert.NoError(t, tr.CleanUp())
}
