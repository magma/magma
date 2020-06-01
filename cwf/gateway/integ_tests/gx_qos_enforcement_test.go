// +build all qos

/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package integration

import (
	"encoding/json"
	"fmt"
	cwfprotos "magma/cwf/cloud/go/protos"
	"magma/feg/cloud/go/protos"
	fegProtos "magma/feg/cloud/go/protos"
	"magma/lte/cloud/go/plugin/models"
	lteProtos "magma/lte/cloud/go/protos"

	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
)

const (
	ErrMargin = 10
)

func verifyEgressRate(t *testing.T, tr *TestRunner, req *cwfprotos.GenTrafficRequest, expRate float64) {
	resp, err := tr.GenULTraffic(req)
	if err != nil {
		t.Fatalf("error %v generating traffic", err)
	}
	// Wait for the traffic to go through
	time.Sleep(6 * time.Second)
	if resp != nil {
		var perfResp map[string]interface{}
		json.Unmarshal([]byte(resp.Output), &perfResp)
		respEndRecd := perfResp["end"].(map[string]interface{})
		respEndRcvMap := respEndRecd["sum_received"].(map[string]interface{})
		b := respEndRcvMap["bits_per_second"].(float64)

		errRate := math.Abs((b-expRate)/expRate) * 100
		fmt.Printf("bit rate observed at server %.0fbps, err rate %.2f%%\n", b, errRate)
		if (b > expRate) && (errRate > ErrMargin) {
			fmt.Printf("recd bps %f exp bps %f\n", b, expRate)
			assert.Fail(t, "error greater than acceptable margin")
		}
	}
}

//TestGxUplinkTrafficQosEnforcement
// This test verifies the QOS configuration(uplink) present in the rules
// - Set an expectation for a  CCR-I to be sent up to PCRF, to which it will
//   respond with a rule install (static-ULQos) with QOS config setting with
//   maximum uplink bitrate.
// - Generate traffic and verify if the traffic observed bitrate matches the configured
// bitrate
func TestGxUplinkTrafficQosEnforcement(t *testing.T) {
	fmt.Println("\nRunning TestGxUplinkTrafficQosEnforcement")
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	tr := NewTestRunner(t)
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)
	assert.NoError(t, usePCRFMockDriver())
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
		clearPCRFMockDriver()
	}()

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)
	imsi := ues[0].GetImsi()

	ki := rand.Intn(1000000)
	monitorKey := fmt.Sprintf("monitor-ULQos-%d", ki)
	ruleKey := fmt.Sprintf("static-ULQos-%d", ki)

	uplinkBwMax := uint32(1000000)
	qos := &models.FlowQos{MaxReqBwUl: &uplinkBwMax}
	rule := getStaticPassAll(ruleKey, monitorKey, 0, models.PolicyRuleTrackingTypeONLYPCRF, 3, qos)

	err = ruleManager.AddStaticRuleToDB(rule)
	assert.NoError(t, err)
	tr.WaitForPoliciesToSync()

	usageMonitorInfo := getUsageInformation(monitorKey, 10*MegaBytes)
	initRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGxCCAnswer(diam.Success).
		SetStaticRuleInstalls([]string{ruleKey}, []string{}).
		SetUsageMonitorInfo(usageMonitorInfo)
	initExpectation := protos.NewGxCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	// On unexpected requests, just return the default update answer
	assert.NoError(t, setPCRFExpectations([]*protos.GxCreditControlExpectation{initExpectation},
		protos.NewGxCCAnswer(diam.Success)))

	tr.AuthenticateAndAssertSuccess(imsi)
	req := &cwfprotos.GenTrafficRequest{
		Imsi:   imsi,
		Volume: &wrappers.StringValue{Value: *swag.String("5M")},
	}
	verifyEgressRate(t, tr, req, float64(uplinkBwMax))

	// Assert that enforcement_stats rules are properly installed and the right
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	record := recordsBySubID["IMSI"+imsi][ruleKey]
	assert.NotNil(t, record, fmt.Sprintf("No policy usage record for imsi: %v", imsi))

	tr.DisconnectAndAssertSuccess(imsi)
	fmt.Println("wait for flows to get deactivated")
	time.Sleep(3 * time.Second)
}

//TestGxDownlinkTrafficQosEnforcement
// This test verifies the QOS configuration(downlink) present in the rules
// - Set an expectation for a  CCR-I to be sent up to PCRF, to which it will
//   respond with a rule install (static-DLQos) with QOS config setting with
//   maximum downlink bitrate.
// - Generate traffic from server to client and verify if the traffic observed bitrate
//   matches the configured bitrate
func TestGxDownlinkTrafficQosEnforcement(t *testing.T) {
	fmt.Println("\nRunning TestGxDownlinkTrafficQosEnforcement")
	tr := NewTestRunner(t)
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)
	assert.NoError(t, usePCRFMockDriver())
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
		clearPCRFMockDriver()
	}()

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)
	imsi := ues[0].GetImsi()

	ki := rand.Intn(1000000)
	monitorKey := fmt.Sprintf("monitor-DLQos-%d", ki)
	ruleKey := fmt.Sprintf("static-DLQos-%d", ki)

	downlinkBwMax := uint32(1000000)
	qos := &models.FlowQos{MaxReqBwDl: &downlinkBwMax}
	rule := getStaticPassAll(ruleKey, monitorKey, 0, models.PolicyRuleTrackingTypeONLYPCRF, 3, qos)

	err = ruleManager.AddStaticRuleToDB(rule)
	assert.NoError(t, err)
	tr.WaitForPoliciesToSync()

	usageMonitorInfo := getUsageInformation(monitorKey, 10*MegaBytes)
	initRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGxCCAnswer(diam.Success).
		SetStaticRuleInstalls([]string{ruleKey}, []string{}).
		SetUsageMonitorInfo(usageMonitorInfo)
	initExpectation := protos.NewGxCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	// On unexpected requests, just return the default update answer
	assert.NoError(t, setPCRFExpectations([]*protos.GxCreditControlExpectation{initExpectation},
		protos.NewGxCCAnswer(diam.Success)))

	tr.AuthenticateAndAssertSuccess(imsi)
	req := &cwfprotos.GenTrafficRequest{
		Imsi:        imsi,
		ReverseMode: true,
		Volume:      &wrappers.StringValue{Value: *swag.String("5M")},
	}
	verifyEgressRate(t, tr, req, float64(downlinkBwMax))

	// Assert that enforcement_stats rules are properly installed and the right
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	record := recordsBySubID["IMSI"+imsi][ruleKey]
	assert.NotNil(t, record, fmt.Sprintf("No policy usage record for imsi: %v", imsi))

	tr.DisconnectAndAssertSuccess(imsi)
	fmt.Println("wait for flows to get deactivated")
	time.Sleep(3 * time.Second)
}

//TestGxQosDowngradeWithCCAUpdate
// This test verifies QOS downgrade which can be caused due to CCA-update
// - Set an expectation for a  CCR-I to be sent up to PCRF, to which it will
//   respond with a rule install (static-qos-CCAI) with QOS config setting with
//   maximum uplink bitrate.
// - Generate traffic and verify if the traffic observed bitrate matches the initially
// configured bitrate
// - Set an expectation for a  CCR-U to be sent up to PCRF, which will respond with a
// rule install static-qos-CCAU which will downgrade the QOS setting for the uplink
// - Generate traffic and verify if the traffic observed bitrate matches the newly
// downgraded bitrate
// - Send another CCA-update which upgrades the QOS through a dynamic rule and verify
// that the observed bitrate maches the newly configured bitrate
func TestGxQosDowngradeWithCCAUpdate(t *testing.T) {
	fmt.Println("\nRunning TestGxQosDowngradeWithCCAUpdate")
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
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

	ki := rand.Intn(1000000)
	monitorKey := fmt.Sprintf("monitor-qos-ccaupdate-%d", ki)
	rule1Key := fmt.Sprintf("static-qos-CCAI-%d", ki)
	rule2Key := fmt.Sprintf("static-qos-CCAU-%d", ki+1)

	// Add 2 static rules to db, one with higher qos and one with lower qos
	uplinkBwInitial := uint32(2000000)
	uplinkBwFinal := uint32(500000)

	rule1 := getStaticPassAll(rule1Key, monitorKey, 0,
		models.PolicyRuleTrackingTypeONLYPCRF, 3, &models.FlowQos{MaxReqBwUl: &uplinkBwInitial})

	rule2 := getStaticPassAll(rule2Key, monitorKey, 0,
		models.PolicyRuleTrackingTypeONLYPCRF, 2, &models.FlowQos{MaxReqBwUl: &uplinkBwFinal})

	for _, r := range []*lteProtos.PolicyRule{rule1, rule2} {
		err = ruleManager.AddStaticRuleToDB(r)
		assert.NoError(t, err)
	}
	tr.WaitForPoliciesToSync()

	// usage monitor for init and upgrade
	usageMonitorInfo := getUsageInformation(monitorKey, 5*MegaBytes)
	initRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGxCCAnswer(diam.Success).
		SetStaticRuleInstalls([]string{rule1Key}, []string{}).
		SetUsageMonitorInfo(usageMonitorInfo)
	initExpectation := protos.NewGxCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	// We expect an update request with some usage update (probably around 80-100% of the given quota)
	var c float64 = 0.3 * 5 * MegaBytes
	updateRequest1 := protos.NewGxCCRequest(imsi, protos.CCRequestType_UPDATE).
		SetUsageMonitorReport(usageMonitorInfo).
		SetUsageReportDelta(uint64(c))
	updateAnswer1 := protos.NewGxCCAnswer(diam.Success).
		SetStaticRuleInstalls([]string{rule2Key}, []string{}).
		SetUsageMonitorInfo(getUsageInformation(monitorKey, 10*MegaBytes))
	updateExpectation1 := protos.NewGxCreditControlExpectation().Expect(updateRequest1).Return(updateAnswer1)

	expectations := []*protos.GxCreditControlExpectation{initExpectation, updateExpectation1}

	// On unexpected requests, just return the default update answer
	assert.NoError(t, setPCRFExpectations(expectations, protos.NewGxCCAnswer(diam.Success)))

	tr.AuthenticateAndAssertSuccess(imsi)

	req := &cwfprotos.GenTrafficRequest{Imsi: imsi, Volume: &wrappers.StringValue{Value: *swag.String("5M")}}
	verifyEgressRate(t, tr, req, float64(uplinkBwInitial))

	// wait for the update to kick in
	time.Sleep(3 * time.Second)

	// verify with lower bitrate and check if constraints are met
	req = &cwfprotos.GenTrafficRequest{Imsi: imsi, Volume: &wrappers.StringValue{Value: *swag.String("5M")}}
	verifyEgressRate(t, tr, req, float64(uplinkBwFinal))

	// wait for the update to kick in
	time.Sleep(3 * time.Second)

	// Assert that enforcement_stats rules are properly installed and the right
	// amount of data was passed through
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	record := recordsBySubID["IMSI"+imsi][rule1Key]
	assert.NotNil(t, record, fmt.Sprintf("No policy usage record for imsi: %v", imsi))

	record = recordsBySubID["IMSI"+imsi][rule2Key]
	assert.NotNil(t, record, fmt.Sprintf("No policy usage record for imsi: %v", imsi))

	// Assert that a CCR-I and at least one CCR-U were sent up to the PCRF
	tr.AssertAllGxExpectationsMetNoError()

	// When we initiate a UE disconnect, we expect a terminate request to go up
	terminateRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_TERMINATION)
	terminateAnswer := protos.NewGxCCAnswer(diam.Success)
	terminateExpectation := protos.NewGxCreditControlExpectation().Expect(terminateRequest).Return(terminateAnswer)
	expectations = []*protos.GxCreditControlExpectation{terminateExpectation}
	assert.NoError(t, setPCRFExpectations(expectations, nil))

	tr.DisconnectAndAssertSuccess(imsi)
	fmt.Println("wait for flows to get deactivated")
	time.Sleep(6 * time.Second)

	// Assert that we saw a Terminate request
	tr.AssertAllGxExpectationsMetNoError()
}

//TestGxQosDowngradeWithReAuth
// This test verifies QOS downgrade which can be caused due to ReAuth Request
// - Set an expectation for a  CCR-I to be sent up to PCRF, to which it will
//   respond with a rule install (static-qos-CCAI) with QOS config setting with
//   maximum uplink bitrate.
// - Generate traffic and verify if the traffic observed bitrate matches the initially
// configured bitrate
// - Set an expectation for a  CCR-U to be sent up to PCRF, which will respond with a
// rule install static-qos-CCAU which will downgrade the QOS setting for the uplink
// - Generate traffic and verify if the traffic observed bitrate matches the newly
// downgraded bitrate
func TestGxQosDowngradeWithReAuth(t *testing.T) {
	fmt.Println("\nRunning TestGxQosDowngradeWithReAuth")
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	tr := NewTestRunner(t)
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)
	imsi := ues[0].GetImsi()

	ki := rand.Intn(1000000)
	monitorKey := fmt.Sprintf("monitor-RAR-%d", ki)
	rule1Key := fmt.Sprintf("static-CCAI-%d", ki)
	rule2Key := fmt.Sprintf("static-RAR-%d", ki)
	uplinkBwInitial := uint32(2000000)
	uplinkBwFinal := uint32(500000)
	rule1 := getStaticPassAll(rule1Key, monitorKey, 0, models.PolicyRuleTrackingTypeONLYPCRF, 3,
		&models.FlowQos{MaxReqBwUl: &uplinkBwInitial})
	rule2 := getStaticPassAll(rule2Key, monitorKey, 0,
		models.PolicyRuleTrackingTypeONLYPCRF, 1, &models.FlowQos{MaxReqBwUl: &uplinkBwFinal})

	for _, r := range []*lteProtos.PolicyRule{rule1, rule2} {
		err = ruleManager.AddStaticRuleToDB(r)
		assert.NoError(t, err)
	}
	tr.WaitForPoliciesToSync()

	err = ruleManager.AddUsageMonitor(imsi, monitorKey, 20*MegaBytes, 1*MegaBytes)
	err = ruleManager.AddRulesToPCRF(imsi, []string{rule1Key}, []string{})
	assert.NoError(t, err)

	tr.AuthenticateAndAssertSuccess(imsi)
	req := &cwfprotos.GenTrafficRequest{
		Imsi:   imsi,
		Volume: &wrappers.StringValue{Value: *swag.String("5M")},
	}
	verifyEgressRate(t, tr, req, float64(uplinkBwInitial))

	// Install a static rule with lower qos
	rarUsageMonitor := getUsageInformation(monitorKey, 50*MegaBytes)
	raa, err := sendPolicyReAuthRequest(
		&fegProtos.PolicyReAuthTarget{
			Imsi:                 imsi,
			RulesToInstall:       &fegProtos.RuleInstalls{RuleNames: []string{rule2Key}},
			UsageMonitoringInfos: []*fegProtos.UsageMonitoringInformation{rarUsageMonitor},
		},
	)
	assert.NoError(t, err)
	tr.WaitForReAuthToProcess()

	// Check ReAuth success
	assert.Contains(t, raa.SessionId, "IMSI"+imsi)
	assert.Equal(t, diam.Success, int(raa.ResultCode))

	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	verifyEgressRate(t, tr, req, float64(uplinkBwFinal))

	// Assert that enforcement_stats rules are properly installed and the right
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	record := recordsBySubID["IMSI"+imsi][rule1Key]
	assert.NotNil(t, record, fmt.Sprintf("No policy usage record for imsi: %v", imsi))

	record = recordsBySubID["IMSI"+imsi][rule2Key]
	assert.NotNil(t, record, fmt.Sprintf("No policy usage record for imsi: %v", imsi))

	tr.DisconnectAndAssertSuccess(imsi)
	fmt.Println("wait for flows to get deactivated")
	time.Sleep(3 * time.Second)
}
