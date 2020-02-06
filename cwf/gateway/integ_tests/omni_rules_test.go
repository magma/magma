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
	// Add a pass all static rule not tied to any monitor usage
	err = ruleManager.AddStaticPassAll("omni-pass-all", "", models.PolicyRuleTrackingTypeNOTRACKING)
	assert.NoError(t, err)

	// Apply a network wide rule that points to the static rule above
	err = ruleManager.AddOmniPresentRules("onmi", []string{"omni-pass-all"}, []string{""})
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

	record := recordsBySubID["IMSI"+ue.GetImsi()]

	assert.NotNil(t, record, fmt.Sprintf("No policy usage record for imsi: %v", ue.GetImsi()))
	assert.Equal(t, "omni-pass-all", record.GetRuleId())

	assert.True(t, record.BytesTx > uint64(0), fmt.Sprintf("%s did not pass any data", record.RuleId))
	assert.True(t, record.BytesTx <= uint64(200*KiloBytes+Buffer), fmt.Sprintf("policy usage: %v", record))

	_, err = tr.Disconnect(ue.GetImsi())
	assert.NoError(t, err)

	// Clear hss, ocs, and pcrf
	assert.NoError(t, ruleManager.RemoveInstalledRules())
	assert.NoError(t, tr.CleanUp())
}
