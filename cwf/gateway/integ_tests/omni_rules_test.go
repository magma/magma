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
	"magma/feg/gateway/services/eap"
	"magma/lte/cloud/go/plugin/models"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticateUplinkTrafficWithOmniRules(t *testing.T) {
	fmt.Printf("Running TestAuthenticateUplinkTrafficWithOmniRules...\n")

	tr := NewTestRunner()
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)

	ue := ues[0]
	// Set a block all rule to be installed by the PCRF
	err = ruleManager.AddStaticRule(getStaticDenyAll("static-block-all", "mkey1", 0, models.PolicyRuleConfigTrackingTypeONLYPCRF, 30))
	assert.NoError(t, err)
	err = ruleManager.AddDynamicRules(ue.Imsi, []string{"static-block-all"}, nil)
	assert.NoError(t, err)

	// Override with an omni pass all static rule with a higher priority
	err = ruleManager.AddStaticPassAll("omni-pass-all-1", "", 0, models.PolicyRuleTrackingTypeNOTRACKING, 20)
	assert.NoError(t, err)
	// Apply a network wide rule that points to the static rule above
	err = ruleManager.AddOmniPresentRules("onmi", []string{"omni-pass-all-1"}, []string{""})
	assert.NoError(t, err)

	// Wait for rules propagation
	time.Sleep(2 * time.Second)
	radiusP, err := tr.Authenticate(ue.GetImsi())
	assert.NoError(t, err)

	eapMessage := radiusP.Attributes.Get(rfc2869.EAPMessage_Type)
	assert.NotNil(t, eapMessage)
	assert.True(t, reflect.DeepEqual(int(eapMessage[0]), eap.SuccessCode))

	err = tr.GenULTraffic(ue.GetImsi(), swag.String("200K"))
	assert.NoError(t, err)

	// Wait for traffic to go through
	time.Sleep(6 * time.Second)
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)

	omniRecord := recordsBySubID["IMSI"+ue.GetImsi()]["omni-pass-all-1"]
	blockAllRecord := recordsBySubID["IMSI"+ue.GetImsi()]["static-block-all"]
	assert.NotNil(t, omniRecord, fmt.Sprintf("No policy usage omniRecord for imsi: %v", ue.GetImsi()))
	assert.NotNil(t, blockAllRecord, fmt.Sprintf("Block all record was not installed for imsi %v", ue.GetImsi()))

	assert.True(t, omniRecord.BytesTx > uint64(0), fmt.Sprintf("%s did not pass any data", omniRecord.RuleId))
	assert.True(t, omniRecord.BytesTx <= uint64(200*KiloBytes+Buffer), fmt.Sprintf("policy usage: %v", omniRecord))
	assert.Equal(t, uint64(0x0), blockAllRecord.BytesTx)
	_, err = tr.Disconnect(ue.GetImsi())
	assert.NoError(t, err)

	// Clear hss, ocs, and pcrf
	assert.NoError(t, ruleManager.RemoveInstalledRules())
	assert.NoError(t, tr.CleanUp())
}
