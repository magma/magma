/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package integ_tests

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"fbc/lib/go/radius/rfc2869"
	cwfprotos "magma/cwf/cloud/go/protos"
	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/eap"
	"magma/lte/cloud/go/plugin/models"
	lteProtos "magma/lte/cloud/go/protos"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
)

const (
	KiloBytes = 1024
	MegaBytes = 1024 * KiloBytes
	Buffer    = 50 * KiloBytes
)

func TestAuthenticateUplinkTrafficWithEnforcement(t *testing.T) {
	fmt.Printf("\nRunning TestAuthenticateUplinkTrafficWithEnforcement...\n")
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

	radiusP, err := tr.Authenticate(imsi)
	assert.NoError(t, err)

	eapMessage := radiusP.Attributes.Get(rfc2869.EAPMessage_Type)
	assert.NotNil(t, eapMessage, fmt.Sprintf("EAP Message from authentication is nil"))
	assert.True(t, reflect.DeepEqual(int(eapMessage[0]), eap.SuccessCode), fmt.Sprintf("UE Authentication did not return success"))

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
	assert.NoError(t, clearPCRFMockDriver())

	// Clear hss, ocs, and pcrf
	assert.NoError(t, ruleManager.RemoveInstalledRules())
	assert.NoError(t, tr.CleanUp())
}

func TestAuthenticateUplinkTrafficWithQosEnforcement(t *testing.T) {
	fmt.Printf("\nRunning TestAuthenticateUplinkTrafficWithQosEnforcement...\n")
	tr := NewTestRunner()
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)
	imsi := ues[0].GetImsi()

	err = ruleManager.AddUsageMonitor(imsi, "mqos1", 1000*MegaBytes, 250*MegaBytes)
	assert.NoError(t, err)

	rule := getStaticPassAll("static-qos-1", "mqos1", 0, models.PolicyRuleTrackingTypeONLYPCRF, 3)
	rule.Qos = &lteProtos.FlowQos{
		MaxReqBwUl: uint32(1000000),
		GbrUl:      uint32(12000),
	}

	err = ruleManager.AddStaticRuleToDB(rule)
	assert.NoError(t, err)
	err = ruleManager.AddRulesToPCRF(imsi, []string{"static-qos-1"}, nil)
	assert.NoError(t, err)

	// wait for the rules to be synced into sessiond
	time.Sleep(1 * time.Second)

	radiusP, err := tr.Authenticate(imsi)
	assert.NoError(t, err)

	eapMessage := radiusP.Attributes.Get(rfc2869.EAPMessage_Type)
	assert.NotNil(t, eapMessage, fmt.Sprintf("EAP Message from authentication is nil"))
	assert.True(t, reflect.DeepEqual(int(eapMessage[0]), eap.SuccessCode), fmt.Sprintf("UE Authentication did not return success"))

	req := &cwfprotos.GenTrafficRequest{
		Imsi:       imsi,
		Bitrate:    &wrappers.StringValue{Value: "5m"},
		TimeInSecs: uint64(10)}

	resp, err := tr.GenULTraffic(req)
	assert.NoError(t, err)

	// Wait for the traffic to go through
	time.Sleep(6 * time.Second)

	// Assert that enforcement_stats rules are properly installed and the right
	// amount of data was passed through
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	record := recordsBySubID["IMSI"+imsi]["static-qos-1"]
	assert.NotNil(t, record, fmt.Sprintf("No policy usage record for imsi: %v", imsi))

	if resp != nil {
		var perfResp map[string]interface{}
		json.Unmarshal([]byte(resp.Output), &perfResp)

		intervalResp := perfResp["intervals"].([]interface{})
		assert.Equal(t, len(intervalResp), 10)

		// verify that starting bit rate was > 500k
		firstIntvl := intervalResp[0].(map[string]interface{})
		firstIntvlSumMap := firstIntvl["sum"].(map[string]interface{})
		b := firstIntvlSumMap["bits_per_second"].(float64)
		fmt.Println("initial bit rate transmitted by traffic gen", b)
		assert.GreaterOrEqual(t, b, float64(500*1024))

		// Ensure that the overall bitrate recd by server was <= 128k
		respEndRecd := perfResp["end"].(map[string]interface{})
		respEndRcvMap := respEndRecd["sum_received"].(map[string]interface{})
		b = respEndRcvMap["bits_per_second"].(float64)
		fmt.Println("bit rate observed at server ", b)
		assert.LessOrEqual(t, b, float64(1000000))
	}
	_, err = tr.Disconnect(imsi)
	assert.NoError(t, err)
	time.Sleep(3 * time.Second)
	// Clear hss, ocs, and pcrf
	assert.NoError(t, ruleManager.RemoveInstalledRules())
	assert.NoError(t, tr.CleanUp())
	fmt.Println("wait for flows to get deactivated")
	time.Sleep(10 * time.Second)
}

func testAuthenticateDownlinkTrafficWithQosEnforcement(t *testing.T) {
	fmt.Printf("\nRunning TestAuthenticateDownlinkTrafficWithQosEnforcement...\n")
	tr := NewTestRunner()
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)
	imsi := ues[0].GetImsi()

	err = ruleManager.AddUsageMonitor(imsi, "mqos2", 1000*MegaBytes, 250*MegaBytes)
	assert.NoError(t, err)

	rule := getStaticPassAll("static-qos-2", "mqos2", 0, models.PolicyRuleTrackingTypeONLYPCRF, 3)
	rule.Qos = &lteProtos.FlowQos{
		MaxReqBwDl: uint32(1000000),
		GbrDl:      uint32(12000),
	}

	err = ruleManager.AddStaticRuleToDB(rule)
	assert.NoError(t, err)
	err = ruleManager.AddRulesToPCRF(imsi, []string{"static-qos-2"}, nil)
	assert.NoError(t, err)

	// wait for the rules to be synced into sessiond
	time.Sleep(3 * time.Second)

	radiusP, err := tr.Authenticate(imsi)
	assert.NoError(t, err)

	eapMessage := radiusP.Attributes.Get(rfc2869.EAPMessage_Type)
	assert.NotNil(t, eapMessage, fmt.Sprintf("EAP Message from authentication is nil"))
	assert.True(t, reflect.DeepEqual(int(eapMessage[0]), eap.SuccessCode), fmt.Sprintf("UE Authentication did not return success"))

	req := &cwfprotos.GenTrafficRequest{
		Imsi:        imsi,
		Bitrate:     &wrappers.StringValue{Value: "5m"},
		TimeInSecs:  uint64(10),
		ReverseMode: true,
	}

	resp, err := tr.GenULTraffic(req)
	assert.NoError(t, err)

	// Wait for the traffic to go through
	time.Sleep(6 * time.Second)

	// Assert that enforcement_stats rules are properly installed and the right
	// amount of data was passed through
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	record := recordsBySubID["IMSI"+imsi]["static-qos-2"]
	assert.NotNil(t, record, fmt.Sprintf("No policy usage record for imsi: %v", imsi))

	if resp != nil {
		var perfResp map[string]interface{}
		json.Unmarshal([]byte(resp.Output), &perfResp)

		// Ensure that the overall bitrate recd by server was <= 128k
		respEndRecd := perfResp["end"].(map[string]interface{})
		respEndRcvMap := respEndRecd["sum_received"].(map[string]interface{})
		b := respEndRcvMap["bits_per_second"].(float64)
		fmt.Println("bit rate observed at server ", b)
		assert.LessOrEqual(t, b, float64(1000000))
	}
	_, err = tr.Disconnect(imsi)
	assert.NoError(t, err)
	time.Sleep(3 * time.Second)

	// Clear hss, ocs, and pcrf
	assert.NoError(t, ruleManager.RemoveInstalledRules())
	assert.NoError(t, tr.CleanUp())
	clearPCRFMockDriver()
}
