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
	fegProtos "magma/feg/cloud/go/protos"
	fegprotos "magma/feg/cloud/go/protos"
	"magma/lte/cloud/go/plugin/models"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
)

func ocsCreditExhaustionTestSetup(t *testing.T) (*TestRunner, *RuleManager, *cwfprotos.UEConfig) {
	tr := NewTestRunner(t)
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)
	setNewOCSConfig(
		&fegprotos.OCSConfig{
			MaxUsageOctets: &fegprotos.Octets{TotalOctets: GyMaxUsageBytes},
			MaxUsageTime:   GyMaxUsageTime,
			ValidityTime:   GyValidityTime,
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

// - Set an expectation for a CCR-I to be sent up to OCS, to which it will
//   NOT respond with any answer.
// - Asset that authentication fails and that no rules were insalled
func TestGyLinksFailureOCStoFEG(t *testing.T) {
	fmt.Println("\nRunning TestGyLinksFailureOCStoFEG...")

	tr, ruleManager, ue := ocsCreditExhaustionTestSetup(t)
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearOCSMockDriver())
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()

	initRequest := protos.NewGyCCRequest(ue.GetImsi(), protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGyCCAnswer(0).SetLinkFailure(true)
	initExpectation := protos.NewGyCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	expectations := []*protos.GyCreditControlExpectation{initExpectation}
	// On unexpected requests, just return the default update answer
	assert.NoError(t, setOCSExpectations(expectations, nil))
	tr.AuthenticateAndAssertFail(ue.Imsi)

	resultByIndex, errByIndex, err := getOCSAssertExpectationsResult()
	assert.NoError(t, err)
	assert.Empty(t, errByIndex)
	expectedResult := []*protos.ExpectationResult{{ExpectationIndex: 0, ExpectationMet: true}}
	assert.ElementsMatch(t, expectedResult, resultByIndex)

	// Since CCA-I was never received, there should be no rules installed
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	assert.Empty(t, recordsBySubID["IMSI"+ue.Imsi])
}

// - Set an expectation for a CCR-I to be sent up to OCS, to which it will
//   respond with a quota grant of 4M and final action set to redirect.
//   Generate traffic and assert the CCR-I is received.
// - Generate 5M traffic to exceed 100% of the quota and validate that session was not terminated
// - Assert that UE flows are NOT deleted.
// - Expect a CCR-T, trigger a UE disconnect, and assert the CCR-T is received.
// NOTE : the test is only verifying that session was not terminated. Improvment is needed to validate
//   that ovs rule is well added and traffic is being redirected.
func TestGyCreditExhaustionRedirect(t *testing.T) {
	fmt.Println("\nRunning TestGyCreditExhaustionRedirect...")

	tr, ruleManager, ue := ocsCreditExhaustionTestSetup(t)
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearOCSMockDriver())
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()

	redirectSrv := fegprotos.RedirectServer{
		RedirectServerAddress: "2.2.2.2",
	}
	quotaGrant := &fegprotos.QuotaGrant{
		RatingGroup: 1,
		GrantedServiceUnit: &fegprotos.Octets{
			TotalOctets: 4 * MegaBytes,
		},
		IsFinalCredit:   true,
		FinalUnitAction: fegprotos.FinalUnitAction_Redirect,
		RedirectServer:  &redirectSrv,
		ResultCode:      2001,
	}

	initRequest := protos.NewGyCCRequest(ue.GetImsi(), protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGyCCAnswer(diam.Success).
		SetQuotaGrant(quotaGrant)
	initExpectation := protos.NewGyCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	defaultUpdateAnswer := protos.NewGyCCAnswer(diam.Success)
	expectations := []*protos.GyCreditControlExpectation{initExpectation}

	// On unexpected requests, just return the default update answer
	assert.NoError(t, setOCSExpectations(expectations, defaultUpdateAnswer))
	tr.AuthenticateAndAssertSuccess(ue.GetImsi())

	// Update directoryd record to include client IP
	err := updateDirectorydRecord("IMSI"+ue.GetImsi(), "ipv4_addr", TrafficCltIP)
	assert.NoError(t, err)

	// we need to generate over 100% of the quota to trigger a session redirection
	req := &cwfprotos.GenTrafficRequest{Imsi: ue.GetImsi(), Volume: &wrappers.StringValue{Value: "5M"}}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	// Check that UE mac flow was not removed
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	record := recordsBySubID["IMSI"+ue.GetImsi()]["static-pass-all-ocs2"]
	assert.NotNil(t, record, fmt.Sprintf("Policy usage record for imsi: %v was not removed", ue.GetImsi()))
	if record != nil {
		// We should not be seeing > 4M data here
		assert.True(t, record.BytesTx > uint64(0), fmt.Sprintf("%s did not pass any data", record.RuleId))
		assert.True(t, record.BytesTx <= uint64(5*MegaBytes+Buffer), fmt.Sprintf("policy usage: %v", record))
	}

	// Assert that a CCR-I to the OCS
	tr.AssertAllGyExpectationsMetNoError()

	// When we initiate a UE disconnect, we expect a terminate request to go up
	terminateRequest := protos.NewGyCCRequest(ue.GetImsi(), protos.CCRequestType_TERMINATION)
	terminateAnswer := protos.NewGyCCAnswer(diam.Success)
	terminateExpectation := protos.NewGyCreditControlExpectation().Expect(terminateRequest).Return(terminateAnswer)
	expectations = []*protos.GyCreditControlExpectation{terminateExpectation}
	assert.NoError(t, setOCSExpectations(expectations, nil))

	// trigger disconnection
	tr.DisconnectAndAssertSuccess(ue.GetImsi())
	tr.WaitForEnforcementStatsToSync()

	// Assert that we saw a Terminate request
	time.Sleep(3 * time.Second)
	tr.AssertAllGyExpectationsMetNoError()
}

func TestGyCreditUpdateCommandLevelFail(t *testing.T) {
	fmt.Println("\nRunning TestGyCreditUpdateFail...")

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
		IsFinalCredit: false,
		ResultCode:    diam.Success,
	}
	initRequest := protos.NewGyCCRequest(ue.GetImsi(), protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGyCCAnswer(diam.Success).SetQuotaGrant(quotaGrant)
	initExpectation := protos.NewGyCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	// Return a permanent failure on Update
	updateRequest := protos.NewGyCCRequest(ue.GetImsi(), protos.CCRequestType_UPDATE)
	// The CCR/A-U exchange fails
	updateAnswer := protos.NewGyCCAnswer(diam.UnableToComply).
		SetQuotaGrant(&fegprotos.QuotaGrant{ResultCode: diam.AuthorizationRejected})
	updateExpectation := protos.NewGyCreditControlExpectation().Expect(updateRequest).Return(updateAnswer)
	// The failure above in CCR/A-U should trigger a termination
	terminateRequest := protos.NewGyCCRequest(ue.GetImsi(), protos.CCRequestType_TERMINATION)
	terminateAnswer := protos.NewGyCCAnswer(diam.Success)
	terminateExpectation := protos.NewGyCreditControlExpectation().Expect(terminateRequest).Return(terminateAnswer)

	expectations := []*protos.GyCreditControlExpectation{initExpectation, updateExpectation, terminateExpectation}
	assert.NoError(t, setOCSExpectations(expectations, nil))

	tr.AuthenticateAndAssertSuccess(ue.GetImsi())
	// Trigger a ReAuth to force an update request
	raa, err := sendChargingReAuthRequest(ue.GetImsi(), 1)
	tr.WaitForReAuthToProcess()

	// Check ReAuth success
	assert.NoError(t, err)
	assert.Contains(t, raa.SessionId, "IMSI"+ue.GetImsi())
	assert.Equal(t, diam.LimitedSuccess, int(raa.ResultCode))

	// Wait for a termination to propagate
	time.Sleep(5 * time.Second)
	tr.WaitForEnforcementStatsToSync()

	// Assert that a CCR-I/U/T was sent to OCS
	tr.AssertAllGyExpectationsMetNoError()

	tr.AssertPolicyEnforcementRecordIsNil(ue.GetImsi())
}

// This test verifies the abort session request
// Here we initially setup a session and install a pass all rule
// We then invoke abort session request from ocs and expect the
// ASR to complete without any error and all the rules associated with
// that session to be cleaned up
func TestGyAbortSessionRequest(t *testing.T) {
	t.Log("Testing TestGyAbortSessionRequest")

	tr := NewTestRunner(t)
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()

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
	tr.WaitForPoliciesToSync()

	// Apply a dynamic rule that points to the static rules above
	err = ruleManager.AddRulesToPCRF(imsi, []string{"static-pass-all-ocs1", "static-pass-all-ocs2"}, nil)
	assert.NoError(t, err)

	tr.AuthenticateAndAssertSuccess(imsi)

	// Generate over 80% of the quota to trigger a CCR Update
	req := &cwfprotos.GenTrafficRequest{
		Imsi:   imsi,
		Volume: &wrappers.StringValue{Value: "4.5M"}}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	// Check that UE mac flow is installed and traffic is less than the quota
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	record := recordsBySubID["IMSI"+imsi]["static-pass-all-ocs2"]
	assert.NotNil(t, record, fmt.Sprintf("Policy usage record for imsi: %v was removed", imsi))
	assert.True(t, record.BytesTx > uint64(0), fmt.Sprintf("%s did not pass any data", record.RuleId))
	assert.True(t, record.BytesTx <= uint64(5*MegaBytes+Buffer), fmt.Sprintf("policy usage: %v", record))

	asa, err := sendChargingAbortSession(
		&fegProtos.AbortSessionRequest{
			Imsi: imsi,
		},
	)
	assert.NoError(t, err)

	// Check for Limited ASR success - There is only limited success here
	// since radius will not do the teardown, radius specifically (COA_DYNAMIC)
	// module throws this error here. coa_dynamic module isn't enabled during
	// authentication and hence it isn't aware of the sessionID used when
	// processing disconnect
	assert.Contains(t, asa.SessionId, "IMSI"+imsi)
	assert.Equal(t, uint32(diam.LimitedSuccess), asa.ResultCode)

	// check if all session related info is cleaned up
	checkSessionAborted := func() bool {
		recordsBySubID, err = tr.GetPolicyUsage()
		assert.NoError(t, err)
		return recordsBySubID["IMSI"+imsi]["static-pass-all-ocs2"] == nil
	}
	assert.Eventually(t, checkSessionAborted, 2*time.Minute, 5*time.Second,
		"request not terminated as expected")

	// trigger disconnection
	tr.DisconnectAndAssertSuccess(imsi)
	fmt.Println("wait for flows to get deactivated")
	time.Sleep(3 * time.Second)
}
