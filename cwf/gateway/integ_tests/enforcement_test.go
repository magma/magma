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

const (
	KiloBytes = 1024
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
	err = tr.GenULTraffic(imsi, swag.String("500K"))
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
	time.Sleep(3 * time.Second)
	// Clear hss, ocs, and pcrf
	assert.NoError(t, ruleManager.RemoveInstalledRules())
	assert.NoError(t, tr.CleanUp())
}
