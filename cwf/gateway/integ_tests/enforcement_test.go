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
	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/eap"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticateUplinkTrafficWithEnforcement(t *testing.T) {
	tr, _ := NewTestRunner()
	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)

	ue := ues[0]
	err = tr.AddPCRFRules(getAllAcceptingPCRFRule(ue.GetImsi()))
	assert.NoError(t, err)
	err = tr.AddPCRFUsageMonitors(getUsageMonitor(ue.GetImsi()))
	assert.NoError(t, err)
	fmt.Printf("************************* Successfully added PCRF rules and monitors for IMSI: %s\n", ue.GetImsi())

	radiusP, err := tr.Authenticate(ue.GetImsi())
	assert.NoError(t, err)

	eapMessage := radiusP.Attributes.Get(rfc2869.EAPMessage_Type)
	assert.NotNil(t, eapMessage)
	assert.True(t, reflect.DeepEqual(int(eapMessage[0]), eap.SuccessCode))

	// Send a small amount of data to start the session
	err = tr.GenULTraffic(ue.GetImsi(), swag.String("2048K"))
	assert.NoError(t, err)

	// Wait for the traffic to go through
	time.Sleep(2 * time.Second)

	// Assert that enforcement_stats rules are properly installed and the right
	// amount of data was passed through
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	subID := "IMSI" + ue.Imsi
	record := recordsBySubID[subID]
	assert.NotNil(t, recordsBySubID[subID])
	assert.Equal(t, makeRuleIDFromIMSI(ue.Imsi), record.RuleId)
	// We should not be seeing > 1024k data here
	assert.True(t, record.BytesTx > uint64(0))
	assert.True(t, record.BytesTx <= uint64(2048000))
	// TODO Talk to PCRF and verify appropriate CCRs propagate up

	// Clear hss, ocs, and pcrf
	assert.NoError(t, tr.CleanUp())
}

func getAllAcceptingPCRFRule(imsi string) *protos.AccountRules {
	return &protos.AccountRules{
		Imsi:          imsi,
		RuleNames:     []string{},
		RuleBaseNames: []string{},
		RuleDefinitions: []*protos.RuleDefinition{
			{
				MonitoringKey:    "mkey1",
				ChargineRuleName: makeRuleIDFromIMSI(imsi),
				Precedence:       100,
				FlowDescriptions: []string{"permit out ip from any to any", "permit in ip from any to any"},
			},
		},
	}
}

func getUsageMonitor(imsi string) *protos.UsageMonitorInfo {
	return &protos.UsageMonitorInfo{
		Imsi: imsi,
		UsageMonitorCredits: []*protos.UsageMonitorCredit{
			{
				MonitoringKey:   "mkey1",
				Volume:          4096,
				ReturnBytes:     1024,
				MonitoringLevel: protos.UsageMonitorCredit_RuleLevel,
			},
		},
	}
}
