// +build all gy

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
	fegprotos "magma/feg/cloud/go/protos"
	"magma/lte/cloud/go/plugin/models"

	"github.com/fiorix/go-diameter/v4/diam"
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
	tr := NewTestRunner(t)
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)
	setNewOCSConfig(
		&fegprotos.OCSConfig{
			MaxUsageOctets: &fegprotos.Octets{TotalOctets: MaxUsageBytes},
			MaxUsageTime:   MaxUsageTime,
			ValidityTime:   ValidityTime,
			UseMockDriver:  true,
		},
	)

	ue := ues[0]
	// Set a pass all rule to be installed by pcrf with a monitoring key to trigger updates
	err = ruleManager.AddUsageMonitor(ue.Imsi, "mkey-ocs", 20*KiloBytes, 10*KiloBytes)
	assert.NoError(t, err)
	err = ruleManager.AddStaticPassAllToDB("static-pass-all-ocs1", "mkey-ocs", 0, models.PolicyRuleConfigTrackingTypeONLYPCRF, 20)
	assert.NoError(t, err)

	// set a pass all rule to be installed by ocs with a rating group 1
	err = ruleManager.AddStaticPassAllToDB("static-pass-all-ocs2", "", 1, models.PolicyRuleConfigTrackingTypeONLYOCS, 10)
	assert.NoError(t, err)

	tr.WaitForPoliciesToSync()

	// Apply a dynamic rule that points to the static rules above
	err = ruleManager.AddRulesToPCRF(ue.Imsi, []string{"static-pass-all-ocs1", "static-pass-all-ocs2"}, nil)
	assert.NoError(t, err)
	return tr, ruleManager, ues[0]
}

// - Set an expectation for a CCR-I to be sent up to OCS, to which it will
//   respond with a quota grant of 4M.
//   Generate traffic and assert the CCR-I is received.
// - Set an expectation for a CCR-U with >80% of data usage to be sent up to
// 	 OCS, to which it will response with more quota.
//   Generate traffic and assert the CCR-U is received with final quota grant.
// - Generate 5M traffic to exceed 100% of the quota and trigger session termination
// - Assert that UE flows are deleted.
// - Expect a CCR-T, trigger a UE disconnect, and assert the CCR-T is received.
func TestGyCreditExhaustionWithCRRU(t *testing.T) {
	fmt.Println("\nRunning TestGyCreditExhaustionWithCRRU...")

	tr, ruleManager, ue := ocsCreditExhaustionTestSetup(t)
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearOCSMockDriver())
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()

	quotaGrant := &fegprotos.QuotaGrant{
		RatingGroup: 1,
		GrantedServiceUnit: &fegprotos.Octets{
			TotalOctets: 5 * MegaBytes,
		},
		IsFinalCredit: false,
		ResultCode:    2001,
	}
	initRequest := protos.NewGyCCRequest(ue.GetImsi(), protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGyCCAnswer(diam.Success).SetQuotaGrant(quotaGrant)
	initExpectation := protos.NewGyCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	// We expect an update request with some usage update (probably around 80-100% of the given quota)
	finalQuotaGrant := &fegprotos.QuotaGrant{
		RatingGroup: 1,
		GrantedServiceUnit: &fegprotos.Octets{
			TotalOctets: 2 * MegaBytes,
		},
		IsFinalCredit:   true,
		FinalUnitAction: fegprotos.FinalUnitAction_Terminate,
		ResultCode:      2001,
	}
	updateRequest1 := protos.NewGyCCRequest(ue.GetImsi(), protos.CCRequestType_UPDATE)
	updateAnswer1 := protos.NewGyCCAnswer(diam.Success).SetQuotaGrant(finalQuotaGrant)
	updateExpectation1 := protos.NewGyCreditControlExpectation().Expect(updateRequest1).Return(updateAnswer1)
	expectations := []*protos.GyCreditControlExpectation{initExpectation, updateExpectation1}

	// On unexpected requests, just return the default update answer
	assert.NoError(t, setOCSExpectations(expectations, updateAnswer1))
	tr.AuthenticateAndAssertSuccess(ue.GetImsi())

	// we need to generate over 80% of the quota to trigger a CCR update
	req := &cwfprotos.GenTrafficRequest{Imsi: ue.GetImsi(), Volume: &wrappers.StringValue{Value: *swag.String("5M")}}
	_, err := tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	// Assert that enforcement_stats rules are properly installed and the right
	// amount of data was passed through
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	record := recordsBySubID["IMSI"+ue.GetImsi()]["static-pass-all-ocs2"]
	assert.NotNil(t, record, fmt.Sprintf("Policy usage record for imsi: %v was removed", ue.GetImsi()))
	if record != nil {
		// We should not be seeing > 1024k data here
		assert.True(t, record.BytesTx > uint64(0), fmt.Sprintf("%s did not pass any data", record.RuleId))
		assert.True(t, record.BytesTx <= uint64(5*MegaBytes+Buffer), fmt.Sprintf("policy usage: %v", record))
	}

	// Assert that a CCR-I and at least one CCR-U were sent up to the OCS
	tr.AssertAllGyExpectationsMetNoError()

	// When we use up all of the quota, we expect a termination request to go up.
	terminateRequest := protos.NewGyCCRequest(ue.GetImsi(), protos.CCRequestType_TERMINATION)
	terminateAnswer := protos.NewGyCCAnswer(diam.Success)
	terminateExpectation := protos.NewGyCreditControlExpectation().Expect(terminateRequest).Return(terminateAnswer)
	expectations = []*protos.GyCreditControlExpectation{terminateExpectation}
	assert.NoError(t, setOCSExpectations(expectations, nil))

	// We need to generate over 100% of the quota to trigger a session termination
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	// Wait for flow deletion due to quota exhaustion
	tr.WaitForEnforcementStatsToSync()

	// Check that UE mac flow is removed
	recordsBySubID, err = tr.GetPolicyUsage()
	assert.NoError(t, err)
	record = recordsBySubID["IMSI"+ue.GetImsi()]["static-pass-all-ocs2"]
	assert.Nil(t, record, fmt.Sprintf("Policy usage record for imsi: %v was not removed", ue.GetImsi()))

	// Assert that we saw a Terminate request
	tr.AssertAllGyExpectationsMetNoError()
}

// - Set an expectation for a CCR-I to be sent up to OCS, to which it will
//   respond with a quota grant of 4M.
//   Generate traffic and assert the CCR-I is received.
// - Generate 5M traffic to exceed 100% of the quota and trigger session termination
// - Assert that UE flows are deleted.
// - Expect a CCR-T, trigger a UE disconnect, and assert the CCR-T is received.
func TestGyCreditExhaustionWithoutCRRU(t *testing.T) {
	fmt.Println("\nRunning TestGyCreditExhaustionWithoutCRRU...")

	tr, ruleManager, ue := ocsCreditExhaustionTestSetup(t)
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearOCSMockDriver())
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()

	quotaGrant := &fegprotos.QuotaGrant{
		RatingGroup: 1,
		GrantedServiceUnit: &fegprotos.Octets{
			TotalOctets: 4 * MegaBytes,
		},
		IsFinalCredit:   true,
		FinalUnitAction: fegprotos.FinalUnitAction_Terminate,
		ResultCode:      2001,
	}
	initRequest := protos.NewGyCCRequest(ue.GetImsi(), protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGyCCAnswer(diam.Success).SetQuotaGrant(quotaGrant)
	initExpectation := protos.NewGyCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	defaultUpdateAnswer := protos.NewGyCCAnswer(diam.Success)
	expectations := []*protos.GyCreditControlExpectation{initExpectation}

	// On unexpected requests, just return the default update answer
	assert.NoError(t, setOCSExpectations(expectations, defaultUpdateAnswer))
	tr.AuthenticateAndAssertSuccess(ue.GetImsi())

	// Assert that a CCR-I was sent to OCS
	tr.AssertAllGyExpectationsMetNoError()

	// When we initiate a UE disconnect, we expect a terminate request to go up
	terminateRequest := protos.NewGyCCRequest(ue.GetImsi(), protos.CCRequestType_TERMINATION)
	terminateAnswer := protos.NewGyCCAnswer(diam.Success)
	terminateExpectation := protos.NewGyCreditControlExpectation().Expect(terminateRequest).Return(terminateAnswer)
	expectations = []*protos.GyCreditControlExpectation{terminateExpectation}
	assert.NoError(t, setOCSExpectations(expectations, nil))

	// we need to generate over 100% of the quota to trigger a session termination
	req := &cwfprotos.GenTrafficRequest{Imsi: ue.GetImsi(), Volume: &wrappers.StringValue{Value: "5M"}}
	_, err := tr.GenULTraffic(req)
	assert.NoError(t, err)
	time.Sleep(5 * time.Second)
	tr.WaitForEnforcementStatsToSync()

	// Assert that we saw a Terminate request
	tr.AssertAllGyExpectationsMetNoError()

	// Check that enforcement stat flow is removed
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	record := recordsBySubID["IMSI"+ue.GetImsi()]["static-pass-all-ocs2"]
	assert.Nil(t, record, fmt.Sprintf("Policy usage record for imsi: %v was removed", ue.GetImsi()))
}
