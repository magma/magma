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

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
)

const (
	ReAuthMaxUsageBytes   = 5 * MegaBytes
	ReAuthMaxUsageTimeSec = 1000 // in second
	ReAuthValidityTime    = 60   // in second
)

func TestAuthenticateUplinkWithOCSChargingReAuth(t *testing.T) {
	fmt.Printf("Running TestAuthenticateUplinkWithOcsChargingReAuth...\n")

	tr := NewTestRunner()
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)

	err = setNewOCSConfig(
		&fegprotos.OCSConfig{
			MaxUsageOctets: &fegprotos.Octets{TotalOctets: ReAuthMaxUsageBytes},
			MaxUsageTime:   ReAuthMaxUsageTimeSec,
			ValidityTime:   ReAuthValidityTime,
		},
	)
	assert.NoError(t, err)

	imsi := ues[0].GetImsi()
	setCreditOnOCS(
		&fegprotos.CreditInfo{
			Imsi:        imsi,
			ChargingKey: 1,
			Volume:      &fegprotos.Octets{TotalOctets: 7 * MegaBytes},
			UnitType:    fegprotos.CreditInfo_Bytes,
		},
	)
	// Set a pass all rule to be installed by pcrf with a monitoring key to trigger updates
	err = ruleManager.AddUsageMonitor(imsi, "mkey-ocs", 500*KiloBytes, 100*KiloBytes)
	assert.NoError(t, err)
	err = ruleManager.AddStaticPassAllToDB("static-pass-all-ocs1", "mkey-ocs", 0, models.PolicyRuleConfigTrackingTypeONLYPCRF, 20)
	assert.NoError(t, err)

	// set a pass all rule to be installed by ocs with a rating group 1
	ratingGroup := uint32(1)
	err = ruleManager.AddStaticPassAllToDB("static-pass-all-ocs2", "", ratingGroup, models.PolicyRuleConfigTrackingTypeONLYOCS, 10)
	assert.NoError(t, err)

	// Apply a dynamic rule that points to the static rules above
	err = ruleManager.AddRulesToPCRF(imsi, []string{"static-pass-all-ocs1", "static-pass-all-ocs2"}, nil)
	assert.NoError(t, err)

	// Wait for rules propagation
	time.Sleep(2 * time.Second)
	radiusP, err := tr.Authenticate(imsi)
	assert.NoError(t, err)

	eapMessage := radiusP.Attributes.Get(rfc2869.EAPMessage_Type)
	assert.NotNil(t, eapMessage)
	assert.True(t, reflect.DeepEqual(int(eapMessage[0]), eap.SuccessCode))
	time.Sleep(2 * time.Second)

	// Generate over 80% of the quota to trigger a CCR Update
	req := &cwfprotos.GenTrafficRequest{
		Imsi:   imsi,
		Volume: &wrappers.StringValue{Value: "4.5M"}}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)

	// Wait for traffic to go through
	time.Sleep(3 * time.Second)

	// Check that UE mac flow is installed and traffic is less than the quota
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	record := recordsBySubID["IMSI"+imsi]["static-pass-all-ocs2"]
	assert.NotNil(t, record, fmt.Sprintf("Policy usage record for imsi: %v was removed", imsi))
	assert.True(t, record.BytesTx > uint64(0), fmt.Sprintf("%s did not pass any data", record.RuleId))
	assert.True(t, record.BytesTx <= uint64(5*MegaBytes+Buffer), fmt.Sprintf("policy usage: %v", record))

	// Top UP extra credits (10M total)
	err = setCreditOnOCS(
		&fegprotos.CreditInfo{
			Imsi:        imsi,
			ChargingKey: ratingGroup,
			Volume:      &fegprotos.Octets{TotalOctets: 3 * MegaBytes},
			UnitType:    fegprotos.CreditInfo_Bytes,
		},
	)
	assert.NoError(t, err)

	// Send ReAuth Request to update quota
	raa, err := sendChargingReAuthRequest(imsi, ratingGroup)

	// Wait for RAR to be processed
	time.Sleep(2 * time.Second)

	// Check ReAuth success
	assert.NoError(t, err)
	assert.Contains(t, raa.SessionId, "IMSI"+imsi)
	assert.Equal(t, uint32(diam.LimitedSuccess), raa.ResultCode)

	// Generate over 7M of data to check that initial quota was updated
	req = &cwfprotos.GenTrafficRequest{Imsi: imsi, Volume: &wrappers.StringValue{Value: "5M"}}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)

	// Wait for traffic to go through
	time.Sleep(10 * time.Second)

	// Check that initial quota was exceeded
	recordsBySubID, err = tr.GetPolicyUsage()
	assert.NoError(t, err)
	record = recordsBySubID["IMSI"+imsi]["static-pass-all-ocs2"]
	assert.NotNil(t, record, fmt.Sprintf("Policy usage record for imsi: %v was removed", imsi))
	assert.True(t, record.BytesTx > uint64(8*MegaBytes+Buffer), fmt.Sprintf("did not pass data over initial quota %v", record))
	assert.True(t, record.BytesTx <= uint64(10*MegaBytes+Buffer), fmt.Sprintf("policy usage: %v", record))

	// trigger disconnection
	_, err = tr.Disconnect(imsi)
	assert.NoError(t, err)
	time.Sleep(3 * time.Second)

	// Clear hss, ocs, and pcrf
	assert.NoError(t, ruleManager.RemoveInstalledRules())
	assert.NoError(t, tr.CleanUp())
	fmt.Println("wait for flows to get deactivated")
	time.Sleep(10 * time.Second)
}
