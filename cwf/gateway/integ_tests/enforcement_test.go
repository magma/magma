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

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticateUplinkTrafficWithEnforcement(t *testing.T) {
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
	err = ruleManager.AddUsageMonitor(imsi, "mkey1", 4096, 1024)
	assert.NoError(t, err)
	err = ruleManager.AddStaticPassAll("static-pass-all", "mkey1")
	assert.NoError(t, err)
	err = ruleManager.AddDynamicRules(imsi, []string{"static-pass-all"}, nil)
	assert.NoError(t, err)

	// wait for the rules to be synced into sessiond
	time.Sleep(1 * time.Second)

	radiusP, err := tr.Authenticate(imsi)
	assert.NoError(t, err)

	eapMessage := radiusP.Attributes.Get(rfc2869.EAPMessage_Type)
	assert.NotNil(t, eapMessage)
	assert.True(t, reflect.DeepEqual(int(eapMessage[0]), eap.SuccessCode))

	// TODO assert CCR-I
	err = tr.GenULTraffic(imsi, swag.String("2048K"))
	assert.NoError(t, err)

	// Wait for the traffic to go through
	time.Sleep(2 * time.Second)

	// Assert that enforcement_stats rules are properly installed and the right
	// amount of data was passed through
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	record := recordsBySubID["IMSI"+imsi]
	assert.NotNil(t, record)
	assert.Equal(t, "static-pass-all", record.RuleId)
	// We should not be seeing > 1024k data here
	assert.True(t, record.BytesTx > uint64(0), fmt.Sprintf("%s did not pass any data", record.RuleId))
	assert.True(t, record.BytesTx <= uint64(2048000), fmt.Sprintf("policy usage: %v", record))
	// TODO Talk to PCRF and verify appropriate CCRs propagate up

	// Clear hss, ocs, and pcrf
	assert.NoError(t, ruleManager.RemoveInstalledRules())
	assert.NoError(t, tr.CleanUp())
}
