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
	"magma/feg/gateway/services/eap"
	"magma/lte/cloud/go/plugin/models"

	lteProtos "magma/lte/cloud/go/protos"

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
	fmt.Printf("Running TestAuthenticateUplinkTrafficWithEnforcement...\n")
	tr := NewTestRunner()
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)
	imsi := ues[0].GetImsi()

	// setup policy rules & monitor
	// 1. Install a usage monitor
	// 2. Install a static rule that passes all traffic tied to the usage monitor above
	// 3. Install a dynamic rule that points to the static rule above
	err = ruleManager.AddUsageMonitor(imsi, "mkey1", 1000*KiloBytes, 250*KiloBytes)
	assert.NoError(t, err)
	err = ruleManager.AddStaticPassAllToDB("static-pass-all", "mkey1", 0, models.PolicyRuleTrackingTypeONLYPCRF, 3)
	assert.NoError(t, err)
	err = ruleManager.AddRulesToPCRF(imsi, []string{"static-pass-all"}, nil)
	assert.NoError(t, err)

	// wait for the rules to be synced into sessiond
	time.Sleep(1 * time.Second)

	radiusP, err := tr.Authenticate(imsi)
	assert.NoError(t, err)

	eapMessage := radiusP.Attributes.Get(rfc2869.EAPMessage_Type)
	assert.NotNil(t, eapMessage, fmt.Sprintf("EAP Message from authentication is nil"))
	assert.True(t, reflect.DeepEqual(int(eapMessage[0]), eap.SuccessCode), fmt.Sprintf("UE Authentication did not return success"))

	// TODO assert CCR-I
	req := &cwfprotos.GenTrafficRequest{Imsi: imsi, Volume: &wrappers.StringValue{Value: *swag.String("500K")}}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)

	// Wait for the traffic to go through
	time.Sleep(6 * time.Second)

	// Assert that enforcement_stats rules are properly installed and the right
	// amount of data was passed through
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	record := recordsBySubID["IMSI"+imsi]["static-pass-all"]
	assert.NotNil(t, record, fmt.Sprintf("No policy usage record for imsi: %v", imsi))
	// We should not be seeing > 1024k data here
	assert.True(t, record.BytesTx > uint64(0), fmt.Sprintf("%s did not pass any data", record.RuleId))
	assert.True(t, record.BytesTx <= uint64(500*KiloBytes+Buffer), fmt.Sprintf("policy usage: %v", record))
	// TODO Talk to PCRF and verify appropriate CCRs propagate up
	_, err = tr.Disconnect(imsi)
	assert.NoError(t, err)

	// Clear hss, ocs, and pcrf
	assert.NoError(t, ruleManager.RemoveInstalledRules())
	assert.NoError(t, tr.CleanUp())
}

func TestAuthenticateUplinkTrafficWithQosEnforcement(t *testing.T) {
	fmt.Printf("Running TestAuthenticateUplinkTrafficWithQosEnforcement...\n")
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
	fmt.Printf("Running TestAuthenticateDownlinkTrafficWithQosEnforcement...\n")
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
}
