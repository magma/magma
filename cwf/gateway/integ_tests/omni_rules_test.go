/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package integration

import (
	"fmt"
	"testing"
	"time"

	cwfprotos "magma/cwf/cloud/go/protos"
	"magma/feg/cloud/go/protos"
	"magma/lte/cloud/go/plugin/models"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
)

// - Insert two static rules into the DB. (static-block-all and omni-pass-all)
// - The omnipresent rule, with higher priority than the block all, will be set
//   as omnipresent/network-wide.
// - Set an expectation for a CCR-I to be sent up to PCRF, to which it will
//   respond with a rule install (static-block-all).
//   Generate traffic and assert the CCR-I is received.
//   Assert that the traffic goes through. This means the network wide rules
//   gets installed properly.
// - Trigger a Gx RAR with a rule removal for the block all rule. Assert the
//   answer is successful. Since the only rule with a usage monitor is removed,
//   the session will terminate. Assert that policy usage is empty.
func TestOmnipresentRules(t *testing.T) {
	fmt.Println("\nRunning TestOmnipresentRules...")
	tr := NewTestRunner(t)
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)
	assert.NoError(t, usePCRFMockDriver())
	defer func() {
		// Delete omni rules
		assert.NoError(t, ruleManager.RemoveOmniPresentRulesFromDB("omni"))
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearPCRFMockDriver())
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)
	imsi := ues[0].GetImsi()

	// Set a block all rule to be installed by the PCRF
	err = ruleManager.AddStaticRuleToDB(getStaticDenyAll("static-block-all", "mkey1", 0, models.PolicyRuleConfigTrackingTypeONLYPCRF, 30))
	assert.NoError(t, err)
	// Override with an omni pass all static rule with a higher priority
	err = ruleManager.AddStaticPassAllToDB("omni-pass-all-1", "", 0, models.PolicyRuleTrackingTypeNOTRACKING, 20)
	assert.NoError(t, err)
	// Apply a network wide rule that points to the static rule above
	err = ruleManager.AddOmniPresentRulesToDB("omni", []string{"omni-pass-all-1"}, []string{""})
	assert.NoError(t, err)
	tr.WaitForPoliciesToSync()

	usageMonitorInfo := []*protos.UsageMonitoringInformation{
		{
			MonitoringLevel: protos.MonitoringLevel_RuleLevel,
			MonitoringKey:   []byte("mkey1"),
			Octets:          &protos.Octets{TotalOctets: 1 * MegaBytes},
		},
	}
	initRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGxCCAnswer(diam.Success).
		SetStaticRuleInstalls([]string{"static-block-all"}, []string{}).
		SetUsageMonitorInfos(usageMonitorInfo)
	initExpectation := protos.NewGxCreditControlExpectation().Expect(initRequest).Return(initAnswer)
	expectations := []*protos.GxCreditControlExpectation{initExpectation}
	assert.NoError(t, setPCRFExpectations(expectations, nil)) // we don't expect any update requests

	tr.AuthenticateAndAssertSuccess(imsi)

	req := &cwfprotos.GenTrafficRequest{Imsi: imsi, Volume: &wrappers.StringValue{Value: *swag.String("200k")}}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	omniRecord := recordsBySubID["IMSI"+imsi]["omni-pass-all-1"]
	blockAllRecord := recordsBySubID["IMSI"+imsi]["static-block-all"]
	assert.NotNil(t, omniRecord, fmt.Sprintf("No policy usage omniRecord for imsi: %v", imsi))
	assert.NotNil(t, blockAllRecord, fmt.Sprintf("Block all record was not installed for imsi %v", imsi))

	assert.True(t, omniRecord.BytesTx > uint64(0), fmt.Sprintf("%s did not pass any data", omniRecord.RuleId))
	assert.Equal(t, uint64(0x0), blockAllRecord.BytesTx)

	// Trigger a ReAuth with rule removals of monitored rules
	target := &protos.PolicyReAuthTarget{
		Imsi: imsi,
		RulesToRemove: &protos.RuleRemovals{
			RuleNames: []string{"static-block-all"},
		},
	}
	fmt.Printf("Sending a ReAuthRequest with target %v\n", target)
	raa, err := sendPolicyReAuthRequest(target)
	tr.WaitForReAuthToProcess()

	// Check ReAuth success
	assert.NoError(t, err)
	assert.Contains(t, raa.SessionId, "IMSI"+imsi)
	assert.Equal(t, uint32(diam.Success), raa.ResultCode)

	// With all monitored rules gone, the session should terminate
	recordsBySubID, err = tr.GetPolicyUsage()
	assert.NoError(t, err)
	assert.Empty(t, recordsBySubID)

	// trigger disconnection
	tr.DisconnectAndAssertSuccess(imsi)
	fmt.Println("wait for flows to get deactivated")
	time.Sleep(3 * time.Second)
}
