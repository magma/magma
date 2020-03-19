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
	"reflect"
	"testing"
	"time"

	"fbc/lib/go/radius/rfc2869"
	cwfprotos "magma/cwf/cloud/go/protos"
	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/eap"
	"magma/lte/cloud/go/plugin/models"

	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
)

const (
	MaxUsageBytes = 5 * 1024 * KiloBytes
	MaxUsageTime  = 1000 // in second
	ValidityTime  = 60   // in second
)

func ocsCreditExhaustionTestSetup(t *testing.T) (*TestRunner, *RuleManager, *cwfprotos.UEConfig) {
	tr := NewTestRunner()
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)
	setNewOCSConfig(
		&fegprotos.OCSConfig{
			MaxUsageOctets: &fegprotos.Octets{TotalOctets: MaxUsageBytes},
			MaxUsageTime:   MaxUsageTime,
			ValidityTime:   ValidityTime,
		},
	)

	ue := ues[0]
	setCreditOnOCS(
		&fegprotos.CreditInfo{
			Imsi:        ue.Imsi,
			ChargingKey: 1,
			Volume:      &fegprotos.Octets{TotalOctets: 7 * 1024 * KiloBytes},
			UnitType:    fegprotos.CreditInfo_Bytes,
		},
	)

	// Set a pass all rule to be installed by pcrf with a monitoring key to trigger updates
	err = ruleManager.AddUsageMonitor(ue.Imsi, "mkey-ocs", 20*KiloBytes, 10*KiloBytes)
	assert.NoError(t, err)
	err = ruleManager.AddStaticPassAllToDB("static-pass-all-ocs1", "mkey-ocs", 0, models.PolicyRuleConfigTrackingTypeONLYPCRF, 20)
	assert.NoError(t, err)

	// set a pass all rule to be installed by ocs with a rating group 1
	err = ruleManager.AddStaticPassAllToDB("static-pass-all-ocs2", "", 1, models.PolicyRuleConfigTrackingTypeONLYOCS, 10)
	assert.NoError(t, err)

	// Apply a dynamic rule that points to the static rules above
	err = ruleManager.AddRulesToPCRF(ue.Imsi, []string{"static-pass-all-ocs1", "static-pass-all-ocs2"}, nil)
	assert.NoError(t, err)

	return tr, ruleManager, ues[0]
}

func TestAuthenticateOcsCreditExhaustedWithCRRU(t *testing.T) {
	fmt.Println("\nRunning TestAuthenticateOcsCreditExhaustedWithCRRU...")

	tr, ruleManager, ue := ocsCreditExhaustionTestSetup(t)
	setCreditOnOCS(
		&fegprotos.CreditInfo{
			Imsi:        ue.Imsi,
			ChargingKey: 1,
			Volume:      &fegprotos.Octets{TotalOctets: 7 * 1024 * KiloBytes},
			UnitType:    fegprotos.CreditInfo_Bytes,
		},
	)

	// Wait for rules propagation
	time.Sleep(2 * time.Second)
	radiusP, err := tr.Authenticate(ue.GetImsi())
	assert.NoError(t, err)

	eapMessage := radiusP.Attributes.Get(rfc2869.EAPMessage_Type)
	assert.NotNil(t, eapMessage)
	assert.True(t, reflect.DeepEqual(int(eapMessage[0]), eap.SuccessCode))
	time.Sleep(2 * time.Second)

	// we need to generate over 80% of the quota to trigger a CCR update
	req := &cwfprotos.GenTrafficRequest{Imsi: ue.GetImsi(), Volume: &wrappers.StringValue{Value: *swag.String("5M")}}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	time.Sleep(3 * time.Second)

	// we need to generate over 100% of the quota to trigger a session termination
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)

	// Wait for traffic to go through
	time.Sleep(5 * time.Second)

	// Check that UE mac flow is removed
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	record := recordsBySubID["IMSI"+ue.GetImsi()]["static-pass-all-ocs2"]
	assert.Nil(t, record, fmt.Sprintf("Policy usage record for imsi: %v was not removed", ue.GetImsi()))

	// trigger disconnection
	_, err = tr.Disconnect(ue.GetImsi())
	assert.NoError(t, err)

	// Clear hss, ocs, and pcrf
	assert.NoError(t, ruleManager.RemoveInstalledRules())
	assert.NoError(t, tr.CleanUp())
}

func TestAuthenticateOcsCreditExhaustedWithoutCRRU(t *testing.T) {
	fmt.Println("\nRunning TestAuthenticateOcsCreditExhaustedWithoutCCRU...")

	tr, ruleManager, ue := ocsCreditExhaustionTestSetup(t)
	setCreditOnOCS(
		&fegprotos.CreditInfo{
			Imsi:        ue.Imsi,
			ChargingKey: 1,
			Volume:      &fegprotos.Octets{TotalOctets: 4 * 1024 * KiloBytes},
			UnitType:    fegprotos.CreditInfo_Bytes,
		},
	)

	// Wait for rules propagation
	time.Sleep(2 * time.Second)
	radiusP, err := tr.Authenticate(ue.GetImsi())
	assert.NoError(t, err)

	eapMessage := radiusP.Attributes.Get(rfc2869.EAPMessage_Type)
	assert.NotNil(t, eapMessage)
	assert.True(t, reflect.DeepEqual(int(eapMessage[0]), eap.SuccessCode))
	time.Sleep(2 * time.Second)

	// we need to generate over 100% of the quota to trigger a session termination
	req := &cwfprotos.GenTrafficRequest{Imsi: ue.GetImsi(), Volume: &wrappers.StringValue{Value: *swag.String("5M")}}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)

	// Wait for traffic to go through
	time.Sleep(5 * time.Second)

	// Check that UE mac flow is removed
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	record := recordsBySubID["IMSI"+ue.GetImsi()]["static-pass-all-ocs2"]
	assert.Nil(t, record, fmt.Sprintf("Policy usage record for imsi: %v was not removed", ue.GetImsi()))

	// trigger disconnection
	_, err = tr.Disconnect(ue.GetImsi())
	assert.NoError(t, err)

	// Clear hss, ocs, and pcrf
	assert.NoError(t, ruleManager.RemoveInstalledRules())
	assert.NoError(t, tr.CleanUp())
}
