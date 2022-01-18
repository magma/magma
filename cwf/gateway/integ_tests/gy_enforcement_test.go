// build all gy

/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package integration

import (
	"fmt"
	"math"
	"testing"
	"time"

	cwfprotos "magma/cwf/cloud/go/protos"
	"magma/feg/cloud/go/protos"
	fegProtos "magma/feg/cloud/go/protos"
	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/session_proxy/credit_control/gy"
	lteprotos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/policydb/obsidian/models"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ocsTestSetup(t *testing.T) (*TestRunner, *RuleManager, *cwfprotos.UEConfig) {
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

	// PCRF Setup: apply a dynamic rule that points to the static rules above
	err = ruleManager.AddRulesToPCRF(ue.Imsi, []string{"static-pass-all-ocs1", "static-pass-all-ocs2"}, nil)
	assert.NoError(t, err)
	return tr, ruleManager, ues[0]
}

func ocsTestSetupSingleRule(t *testing.T) (*TestRunner, *RuleManager, *cwfprotos.UEConfig) {
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

	// set a pass all rule to be installed by ocs with a rating group 1
	err = ruleManager.AddStaticPassAllToDB("static-pass-all-ocs1", "", 1, models.PolicyRuleConfigTrackingTypeONLYOCS, 10)
	assert.NoError(t, err)

	tr.WaitForPoliciesToSync()

	// PCRF Setup: apply a dynamic rule that points to the static rules above
	err = ruleManager.AddRulesToPCRF(ue.Imsi, []string{"static-pass-all-ocs1"}, nil)
	assert.NoError(t, err)
	return tr, ruleManager, ues[0]
}

func provisionRestrictRules(t *testing.T, tr *TestRunner, ruleManager *RuleManager) {
	// Set a block all rule to be installed by the final unit action
	err := ruleManager.AddStaticRuleToDB(
		getStaticDenyAll("restrict-deny-all", "mkey-ocs", 0, models.PolicyRuleConfigTrackingTypeONLYPCRF, 200),
	)
	assert.NoError(t, err)

	// set a pass rule for traffic from TrafficCltIPP
	err = ruleManager.AddStaticRuleToDB(
		getStaticPassTraffic("restrict-pass-user", TrafficCltIP, MATCH_ALL, "mkey-ocs", 0, models.PolicyRuleConfigTrackingTypeONLYPCRF, 100, nil),
	)
	assert.NoError(t, err)

	tr.WaitForPoliciesToSync()
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
	tr, ruleManager, ue := ocsTestSetup(t)
	imsi := ue.GetImsi()
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
		ResultCode:    diam.Success,
	}
	initRequest := protos.NewGyCCRequest(imsi, protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGyCCAnswer(diam.Success).SetQuotaGrant(quotaGrant)
	initExpectation := protos.NewGyCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	// We expect an update request with some usage update (probably around 80-100% of the given quota)
	finalQuotaGrant := &fegprotos.QuotaGrant{
		RatingGroup: 1,
		GrantedServiceUnit: &fegprotos.Octets{
			TotalOctets: 3 * MegaBytes,
		},
		IsFinalCredit: true,
		FinalUnitIndication: &fegprotos.FinalUnitIndication{
			FinalUnitAction: fegprotos.FinalUnitAction_Terminate,
		},
		ResultCode: 2001,
	}
	updateRequest1 := protos.NewGyCCRequest(imsi, protos.CCRequestType_UPDATE)
	updateAnswer1 := protos.NewGyCCAnswer(diam.Success).SetQuotaGrant(finalQuotaGrant)
	updateExpectation1 := protos.NewGyCreditControlExpectation().Expect(updateRequest1).Return(updateAnswer1)
	expectations := []*protos.GyCreditControlExpectation{initExpectation, updateExpectation1}

	// On unexpected requests, just return the default update answer
	assert.NoError(t, setOCSExpectations(expectations, updateAnswer1))
	tr.AuthenticateAndAssertSuccess(imsi)
	// First wait until we see the original static-pass-all-ocs2 show up
	assert.Eventually(t,
		tr.WaitForEnforcementStatsForRule(imsi, "static-pass-all-ocs2"), time.Minute, 2*time.Second)

	// we need to generate over 80% of the quota to trigger a CCR update
	req := &cwfprotos.GenTrafficRequest{
		Imsi:    imsi,
		Volume:  &wrappers.StringValue{Value: "4.5M"},
		Bitrate: &wrappers.StringValue{Value: "30M"},
		Timeout: 60,
	}
	_, err := tr.GenULTraffic(req)
	// Assert that enforcement_stats rules are properly installed and the right
	// amount of data was passed through
	assert.NoError(t, err)
	assert.Eventually(t,
		tr.WaitForEnforcementStatsForRuleGreaterThan(imsi, "static-pass-all-ocs2", 3*MegaBytes), time.Minute, 2*time.Second)

	// Assert that a CCR-I and at least one CCR-U were sent up to the OCS
	tr.AssertAllGyExpectationsMetNoError()

	// When we use up all of the quota, we expect a termination request to go up.
	terminateRequest := protos.NewGyCCRequest(imsi, protos.CCRequestType_TERMINATION)
	terminateAnswer := protos.NewGyCCAnswer(diam.Success)
	terminateExpectation := protos.NewGyCreditControlExpectation().Expect(terminateRequest).Return(terminateAnswer)
	expectations = []*protos.GyCreditControlExpectation{terminateExpectation}
	assert.NoError(t, setOCSExpectations(expectations, nil))

	// We need to generate over 100% of the quota to trigger a session termination
	req = &cwfprotos.GenTrafficRequest{
		Imsi:   imsi,
		Volume: &wrappers.StringValue{Value: "10M"},
	}
	// This could error out if the session terminates during traffic generation
	tr.GenULTraffic(req)

	// Wait for flow deletion due to quota exhaustion -> termination
	tr.AssertEventuallyAllRulesRemovedAfterDisconnect(imsi)

	// Assert that we saw a Terminate request
	tr.AssertAllGyExpectationsMetNoError()
}

func TestGyCreditValidityTime(t *testing.T) {
	fmt.Println("\nRunning TestGyCreditValidityTime...")

	tr, ruleManager, ue := ocsTestSetup(t)
	imsi := ue.GetImsi()
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
		ValidityTime:  3, // seconds
		IsFinalCredit: false,
		ResultCode:    2001,
	}
	initRequest := protos.NewGyCCRequest(imsi, protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGyCCAnswer(diam.Success).SetQuotaGrant(quotaGrant)
	initExpectation := protos.NewGyCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	// We expect an update request with some usage update but not the full quota < 5MB
	mscc := &fegprotos.MultipleServicesCreditControl{
		RatingGroup:     1,
		UsedServiceUnit: &fegprotos.Octets{TotalOctets: 500 * KiloBytes},
		UpdateType:      int32(lteprotos.CreditUsage_VALIDITY_TIMER_EXPIRED),
	}
	// accept 1k - 999k
	// TODO Would be good to add asserting on less than, greater than, etc.
	updateRequest1 := protos.NewGyCCRequest(imsi, protos.CCRequestType_UPDATE).
		SetMSCC(mscc).SetMSCCDelta(499 * KiloBytes)
	updateAnswer1 := protos.NewGyCCAnswer(diam.Success).SetQuotaGrant(quotaGrant)
	updateExpectation1 := protos.NewGyCreditControlExpectation().Expect(updateRequest1).Return(updateAnswer1)
	expectations := []*protos.GyCreditControlExpectation{initExpectation, updateExpectation1}

	// On unexpected requests, just return the default update answer
	assert.NoError(t, setOCSExpectations(expectations, updateAnswer1))

	tr.AuthenticateAndAssertSuccess(imsi)

	require.Eventually(t,
		tr.WaitForEnforcementStatsForRule(imsi, "static-pass-all-ocs2"),
		20*time.Second, 2*time.Second)
	// Generate some traffic but not enough to trigger a quota update request
	// We want the update type to be VALIDITY TIMER EXPIRED
	req := &cwfprotos.GenTrafficRequest{
		Imsi:    imsi,
		Volume:  &wrappers.StringValue{Value: "500K"},
		Bitrate: &wrappers.StringValue{Value: "10M"},
		Timeout: 60,
	}
	_, err := tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	time.Sleep(time.Second * 5)
	tr.AssertAllGyExpectationsMetNoError()
	tr.DisconnectAndAssertSuccess(imsi)
	tr.AssertEventuallyAllRulesRemovedAfterDisconnect(imsi)
}

// - Set an expectation for a CCR-I to be sent up to OCS, to which it will
//   respond with a quota grant of 4M.
//   Generate traffic and assert the CCR-I is received.
// - Generate 5M traffic to exceed 100% of the quota and trigger session termination
// - Assert that UE flows are deleted.
// - Expect a CCR-T, trigger a UE disconnect, and assert the CCR-T is received.
func TestGyCreditExhaustionWithoutCCRU(t *testing.T) {
	fmt.Println("\nRunning TestGyCreditExhaustionWithoutCCRU...")

	tr, ruleManager, ue := ocsTestSetup(t)
	imsi := ue.GetImsi()
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearOCSMockDriver())
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()

	finalUnitIndication := fegprotos.FinalUnitIndication{
		FinalUnitAction: fegprotos.FinalUnitAction_Terminate,
	}
	quotaGrant := &fegprotos.QuotaGrant{
		RatingGroup: 1,
		GrantedServiceUnit: &fegprotos.Octets{
			TotalOctets: 4 * MegaBytes,
		},
		IsFinalCredit:       true,
		FinalUnitIndication: &finalUnitIndication,
		ResultCode:          2001,
	}
	initRequest := protos.NewGyCCRequest(imsi, protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGyCCAnswer(diam.Success).SetQuotaGrant(quotaGrant)
	initExpectation := protos.NewGyCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	defaultUpdateAnswer := protos.NewGyCCAnswer(diam.Success)
	expectations := []*protos.GyCreditControlExpectation{initExpectation}

	// On unexpected requests, just return the default update answer
	assert.NoError(t, setOCSExpectations(expectations, defaultUpdateAnswer))
	tr.AuthenticateAndAssertSuccess(imsi)
	// First wait until we see the original static-pass-all-ocs2 show up
	assert.Eventually(t, tr.WaitForEnforcementStatsForRule(imsi, "static-pass-all-ocs2"), time.Minute, 2*time.Second)
	fmt.Println("CCR-I exchanged installed static-pass-all-ocs2")

	// Assert that a CCR-I was sent to OCS
	tr.AssertAllGyExpectationsMetNoError()

	// When we initiate a UE disconnect, we expect a terminate request to go up
	terminateRequest := protos.NewGyCCRequest(imsi, protos.CCRequestType_TERMINATION)
	terminateAnswer := protos.NewGyCCAnswer(diam.Success)
	terminateExpectation := protos.NewGyCreditControlExpectation().Expect(terminateRequest).Return(terminateAnswer)
	expectations = []*protos.GyCreditControlExpectation{terminateExpectation}
	assert.NoError(t, setOCSExpectations(expectations, nil))

	// we need to generate over 100% of the quota to trigger a session termination
	req := &cwfprotos.GenTrafficRequest{
		Imsi:    imsi,
		Volume:  &wrappers.StringValue{Value: "10M"},
		Timeout: 60,
	}
	// This could error out if the session terminates during traffic generation
	tr.GenULTraffic(req)
	tr.AssertEventuallyAllRulesRemovedAfterDisconnect(imsi)

	// Assert that we saw a Terminate request
	tr.AssertAllGyExpectationsMetNoError()
}

// - Set an expectation for a CCR-I to be sent up to OCS, to which it will
//   NOT respond with any answer.
// - Asset that authentication fails and that no rules were installed
func TestGyLinksFailureOCStoFEG(t *testing.T) {
	fmt.Println("\nRunning TestGyLinksFailureOCStoFEG...")

	tr, ruleManager, ue := ocsTestSetup(t)
	imsi := ue.GetImsi()
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearOCSMockDriver())
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()

	initRequest := protos.NewGyCCRequest(imsi, protos.CCRequestType_INITIAL)
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
// - Generate 5M traffic to exceed 100% of the quota to trigger redirection.
// - When redirection happens, redirect rule is installed on top of the actual rules
// 	 Assert that UE flows are NOT deleted and data was passed
// - Assert redirect rule is installed
// - Send a Charging ReAuth request to top up quota and assert that the
//   response is successful
// - Assert that CCR-U was is generated
// - Assert the redirect rule is gone
// - Generate  traffic and assert that UE flows are NOT deleted and data was passed.
// - Expect a CCR-T, trigger a UE disconnect, and assert the CCR-T is received.
// NOTE : the test is only verifying that session was not terminated.
//        Improvement is needed to validate that ovs rule is well added and
//        traffic is being redirected.
func TestGyCreditExhaustionRedirect(t *testing.T) {
	fmt.Println("\nRunning TestGyCreditExhaustionRedirect...")

	tr, ruleManager, ue := ocsTestSetup(t)
	imsi := ue.GetImsi()
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
		IsFinalCredit: true,
		FinalUnitIndication: &fegprotos.FinalUnitIndication{
			FinalUnitAction: fegprotos.FinalUnitAction_Redirect,
			RedirectServer: &fegprotos.RedirectServer{
				RedirectServerAddress: "2.2.2.2",
			},
		},
		ResultCode: diameter.SuccessCode,
	}

	initRequest := protos.NewGyCCRequest(imsi, protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGyCCAnswer(diam.Success).
		SetQuotaGrant(quotaGrant)
	initExpectation := protos.NewGyCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	expectedMSCC := &protos.MultipleServicesCreditControl{
		RatingGroup: 1,
		UpdateType:  int32(gy.FORCED_REAUTHORISATION),
	}
	// We expect an update request with some usage update after reauth
	updateRequest := protos.NewGyCCRequest(imsi, protos.CCRequestType_UPDATE).
		SetMSCC(expectedMSCC)
	updateAnswer := protos.NewGyCCAnswer(diam.Success).SetQuotaGrant(quotaGrant)
	updateExpectation := protos.NewGyCreditControlExpectation().Expect(updateRequest).
		Return(updateAnswer)
	expectations := []*protos.GyCreditControlExpectation{initExpectation, updateExpectation}

	// On unexpected requests, just return the default update answer
	assert.NoError(t, setOCSExpectations(expectations, updateAnswer))
	tr.AuthenticateAndAssertSuccess(imsi)
	// First wait until we see the original static-pass-all-ocs2 show up
	assert.Eventually(t,
		tr.WaitForEnforcementStatsForRule(imsi, "static-pass-all-ocs2"), time.Minute, 2*time.Second)

	// we need to generate over 100% of the quota to trigger a session redirection
	req := &cwfprotos.GenTrafficRequest{
		Imsi:   imsi,
		Volume: &wrappers.StringValue{Value: "10M"},
		//Bitrate: &wrappers.StringValue{Value: "100M"},
		Timeout: 60,
	}

	//time.Sleep(500 * time.Microsecond)
	_, err := tr.GenULTraffic(req)
	assert.NoError(t, err)

	// Check that enforcement stats flow was not removed and data was passed
	assert.Eventually(t,
		tr.WaitForEnforcementStatsForRuleGreaterThan(imsi, "static-pass-all-ocs2", 3*MegaBytes), time.Minute, 2*time.Second)
	// Wait for service deactivation
	assert.Eventually(t,
		tr.WaitForEnforcementStatsForRule(imsi, "redirect"), time.Minute, 2*time.Second)

	// Send ReAuth Request to update quota
	raa, err := sendChargingReAuthRequest(imsi, 1)
	assert.NoError(t, err)
	assert.Eventually(t, tr.WaitForChargingReAuthToProcess(raa, imsi), time.Minute, 2*time.Second)

	assert.NotNil(t, raa)
	if raa != nil {
		// Check ReAuth success
		assert.Equal(t, diam.LimitedSuccess, int(raa.ResultCode))
	}

	// Assert that a CCR-I and CCR-U were sent to the OCS
	tr.AssertAllGyExpectationsMetNoError()

	// we need to generate more traffic
	req = &cwfprotos.GenTrafficRequest{
		Imsi:    imsi,
		Volume:  &wrappers.StringValue{Value: "2M"},
		Bitrate: &wrappers.StringValue{Value: "30M"},
		Timeout: 60,
	}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	// Check that enforcement stats flow was not removed and data was passed
	assert.Eventually(t,
		tr.WaitForEnforcementStatsForRuleGreaterThan(imsi, "static-pass-all-ocs2", 1*MegaBytes), time.Minute, 2*time.Second)

	// When we initiate a UE disconnect, we expect a terminate request to go up
	terminateRequest := protos.NewGyCCRequest(imsi, protos.CCRequestType_TERMINATION)
	terminateAnswer := protos.NewGyCCAnswer(diam.Success)
	terminateExpectation := protos.NewGyCreditControlExpectation().Expect(terminateRequest).Return(terminateAnswer)
	expectations = []*protos.GyCreditControlExpectation{terminateExpectation}
	assert.NoError(t, setOCSExpectations(expectations, nil))

	// trigger disconnection
	tr.DisconnectAndAssertSuccess(imsi)
	tr.AssertEventuallyAllRulesRemovedAfterDisconnect(imsi)

	// Assert that we saw a Terminate request
	tr.AssertAllGyExpectationsMetNoError()
}

func TestGyCreditUpdateCommandLevelFail(t *testing.T) {
	fmt.Println("\nRunning TestGyCreditUpdateFail...")

	tr, ruleManager, ue := ocsTestSetup(t)
	imsi := ue.GetImsi()
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
	initRequest := protos.NewGyCCRequest(imsi, protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGyCCAnswer(diam.Success).SetQuotaGrant(quotaGrant)
	initExpectation := protos.NewGyCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	// Return a permanent failure on Update
	updateRequest := protos.NewGyCCRequest(imsi, protos.CCRequestType_UPDATE)
	// The CCR/A-U exchange fails
	updateAnswer := protos.NewGyCCAnswer(diam.UnableToComply).
		SetQuotaGrant(&fegprotos.QuotaGrant{ResultCode: diam.AuthorizationRejected})
	updateExpectation := protos.NewGyCreditControlExpectation().Expect(updateRequest).Return(updateAnswer)
	// The failure above in CCR/A-U should trigger a termination
	terminateRequest := protos.NewGyCCRequest(imsi, protos.CCRequestType_TERMINATION)
	terminateAnswer := protos.NewGyCCAnswer(diam.Success)
	terminateExpectation := protos.NewGyCreditControlExpectation().Expect(terminateRequest).Return(terminateAnswer)

	expectations := []*protos.GyCreditControlExpectation{initExpectation, updateExpectation, terminateExpectation}
	assert.NoError(t, setOCSExpectations(expectations, nil))

	tr.AuthenticateAndAssertSuccess(imsi)
	// Trigger a ReAuth to force an update request
	raa, err := sendChargingReAuthRequest(imsi, 1)
	assert.NoError(t, err)
	assert.Eventually(t, tr.WaitForChargingReAuthToProcess(raa, imsi), time.Minute, 2*time.Second)

	assert.NotNil(t, raa)
	if raa != nil {
		// Check ReAuth success
		assert.Equal(t, diam.LimitedSuccess, int(raa.ResultCode))
	}

	// Wait for a termination to propagate
	tr.AssertEventuallyAllRulesRemovedAfterDisconnect(imsi)

	// Assert that a CCR-I/U/T was sent to OCS
	tr.AssertAllGyExpectationsMetNoError()
}

// This test verifies the abort session request
// Here we initially setup a session and install a pass all rule
// We then invoke abort session request from ocs and expect the
// ASR to complete without any error and all the rules associated with
// that session to be cleaned up
func TestGyAbortSessionRequest(t *testing.T) {
	fmt.Println("\nTesting TestGyAbortSessionRequest...")

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
		Imsi:    imsi,
		Volume:  &wrappers.StringValue{Value: "4.5M"},
		Bitrate: &wrappers.StringValue{Value: "40M"},
		Timeout: 60,
	}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)

	// Check that enforcement stats flow is installed and traffic is less than the quota
	assert.Eventually(t,
		tr.WaitForEnforcementStatsForRuleGreaterThan(imsi, "static-pass-all-ocs2", 3*MegaBytes), time.Minute, 2*time.Second)
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
	assert.NotNil(t, asa)
	if asa != nil {
		assert.Contains(t, asa.SessionId, "IMSI"+imsi)
		assert.Equal(t, uint32(diam.LimitedSuccess), asa.ResultCode)
	}

	// check if all session related info is cleaned up
	tr.AssertEventuallyAllRulesRemovedAfterDisconnect(imsi)
}

// - Set an expectation for a CCR-I to be sent up to OCS, to which it will
//   respond with a quota grant of 4M and final action set to redirect.
//   Generate traffic and assert the CCR-I is received.
// - Generate 5M traffic to exceed 100% of the quota to trigger service restriction.
// - Assert that UE flows are NOT deleted and data was passed.
// - Generate an additional 2M traffic and assert that only Gy flows matched.
// - Send a Charging ReAuth request to top up quota and assert that the
//   response is successful
// - Assert that CCR-U was is generated
// - Generate 2M traffic and assert that UE flows are NOT deleted and data was passed.
func TestGyCreditExhaustionRestrict(t *testing.T) {
	fmt.Println("\nRunning TestGyCreditExhaustionRestrict...")

	tr, ruleManager, ue := ocsTestSetup(t)
	imsi := ue.GetImsi()
	defer func() {
		// clear hss, ocs, and pcrf
		assert.NoError(t, clearOCSMockDriver())
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()

	provisionRestrictRules(t, tr, ruleManager)

	quotaGrant := &fegprotos.QuotaGrant{
		RatingGroup: 1,
		GrantedServiceUnit: &fegprotos.Octets{
			TotalOctets: 4 * MegaBytes,
		},
		IsFinalCredit: true,
		FinalUnitIndication: &fegprotos.FinalUnitIndication{
			FinalUnitAction: fegprotos.FinalUnitAction_Restrict,
			RestrictRules:   []string{"restrict-pass-user", "restrict-deny-all"},
		},
		ResultCode: 2001,
	}

	initRequest := protos.NewGyCCRequest(imsi, protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGyCCAnswer(diam.Success).
		SetQuotaGrant(quotaGrant)
	initExpectation := protos.NewGyCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	expectedMSCC := &protos.MultipleServicesCreditControl{
		RatingGroup: 1,
		UpdateType:  int32(gy.FORCED_REAUTHORISATION),
	}
	// We expect an update request with some usage update after reauth
	updateRequest := protos.NewGyCCRequest(imsi, protos.CCRequestType_UPDATE).
		SetMSCC(expectedMSCC)
	updateAnswer := protos.NewGyCCAnswer(diam.Success).SetQuotaGrant(quotaGrant)
	updateExpectation := protos.NewGyCreditControlExpectation().Expect(updateRequest).
		Return(updateAnswer)
	expectations := []*protos.GyCreditControlExpectation{initExpectation, updateExpectation}

	// On unexpected requests, just return the default update answer
	assert.NoError(t, setOCSExpectations(expectations, updateAnswer))
	tr.AuthenticateAndAssertSuccess(imsi)

	// we need to generate over 100% of the quota to trigger a session redirection
	req := &cwfprotos.GenTrafficRequest{
		Imsi:    imsi,
		Volume:  &wrappers.StringValue{Value: "5M"},
		Bitrate: &wrappers.StringValue{Value: "60M"},
		Timeout: 60,
	}
	_, err := tr.GenULTraffic(req)
	assert.NoError(t, err)
	// First wait until we see the original static-pass-all-ocs2 show up
	assert.Eventually(t, tr.WaitForEnforcementStatsForRule(imsi, "static-pass-all-ocs2"), time.Minute, 2*time.Second)
	// Wait for service deactivation
	assert.Eventually(t, tr.WaitForEnforcementStatsForRule(imsi, "restrict-pass-user"), time.Minute, 2*time.Second)

	// we need to generate more traffic and validate it goes through restrict rule
	req = &cwfprotos.GenTrafficRequest{
		Imsi:    imsi,
		Volume:  &wrappers.StringValue{Value: "2M"},
		Bitrate: &wrappers.StringValue{Value: "60M"},
		Timeout: 60,
	}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	// Check that the og stats flow was not removed and flow data hit restrict rule
	tr.AssertPolicyUsage(imsi, "static-pass-all-ocs2", 1, 6*MegaBytes+Buffer)
	// Check that data went through the restrict rule
	tr.AssertPolicyUsage(imsi, "restrict-pass-user", 1, 3*MegaBytes+Buffer)

	// Send ReAuth Request to update quota
	raa, err := sendChargingReAuthRequest(imsi, 1)
	assert.NoError(t, err)
	assert.Eventually(t, tr.WaitForChargingReAuthToProcess(raa, imsi), time.Minute, 2*time.Second)
	assert.NotNil(t, raa)
	if raa != nil {
		// Check ReAuth success
		assert.Equal(t, diam.LimitedSuccess, int(raa.ResultCode))
	}

	// Assert that a CCR-I and CCR-U were sent to the OCS
	tr.AssertAllGyExpectationsMetNoError()

	// Wait for service activation : restrict rule is removed
	assert.Eventually(t, tr.WaitForNoEnforcementStatsForRule(imsi, "restrict-pass-user"), time.Minute, 2*time.Second)

	// we need to generate more traffic to hit restrict rule
	req = &cwfprotos.GenTrafficRequest{
		Imsi:    imsi,
		Volume:  &wrappers.StringValue{Value: "2M"},
		Bitrate: &wrappers.StringValue{Value: "60M"},
		Timeout: 60,
	}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	// Check that og stats flow has seen more data than before
	// The enforcement stats record is not cumulative at this point since
	// the rule is installed again. Assert that it is above 1.
	tr.AssertPolicyUsage(imsi, "static-pass-all-ocs2", 1, 3*MegaBytes+Buffer)

	// trigger disconnection
	tr.DisconnectAndAssertSuccess(imsi)
	tr.AssertEventuallyAllRulesRemovedAfterDisconnect(imsi)
}

// - Send a CCA-I with valid credit for a RG but with 4012 error code (transient)
// - Assert that UE flows for that RG are deleted
// - Generate an additional 2M traffic and assert that only Gy flows matched.
// - Assert that Redirect flows are installed
// - Send a Charging ReAuth request to top up quota and assert that the
//   response is successful
// - Assert that CCR-U was is generated
// - Generate 2M traffic and assert that UE flows are reinstalled for RG
//   and traffic goes through them.
func TestGyCreditTransientErrorRestrict(t *testing.T) {
	fmt.Println("\nRunning TestGyCreditTransientErrorRestrict...")

	tr, ruleManager, ue := ocsTestSetupSingleRule(t)
	imsi := ue.GetImsi()
	defer func() {
		// clear hss, ocs, and pcrf
		assert.NoError(t, clearPCRFMockDriver())
		assert.NoError(t, clearOCSMockDriver())
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()

	provisionRestrictRules(t, tr, ruleManager)

	finalUnitIndication := fegprotos.FinalUnitIndication{
		FinalUnitAction: fegprotos.FinalUnitAction_Restrict,
		RestrictRules:   []string{"restrict-pass-user", "restrict-deny-all"},
	}
	quotaGrant_Init := &fegprotos.QuotaGrant{
		RatingGroup: 1,
		GrantedServiceUnit: &fegprotos.Octets{
			TotalOctets: 0 * MegaBytes,
		},
		IsFinalCredit:       true,
		FinalUnitIndication: &finalUnitIndication,
		ResultCode:          diameter.DiameterCreditLimitReached,
	}

	// CCR-I
	initRequest := protos.NewGyCCRequest(imsi, protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGyCCAnswer(diam.Success).
		SetQuotaGrant(quotaGrant_Init)
	initExpectation := protos.NewGyCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	// reauth
	expectedMSCC_forReauth := &protos.MultipleServicesCreditControl{
		RatingGroup: 1,
		UpdateType:  int32(gy.FORCED_REAUTHORISATION),
	}
	reauthRequest := protos.NewGyCCRequest(imsi, protos.CCRequestType_UPDATE).
		SetMSCC(expectedMSCC_forReauth)

	quotaGrant_Reauth := &fegprotos.QuotaGrant{
		RatingGroup: 1,
		GrantedServiceUnit: &fegprotos.Octets{
			TotalOctets: 10 * MegaBytes,
		},
		IsFinalCredit:       true,
		FinalUnitIndication: &finalUnitIndication,
		ResultCode:          diam.Success,
	}

	reauthAnswer := protos.NewGyCCAnswer(diam.Success).SetQuotaGrant(quotaGrant_Reauth)
	reauthExpectation := protos.NewGyCreditControlExpectation().Expect(reauthRequest).
		Return(reauthAnswer)

	expectations := []*protos.GyCreditControlExpectation{initExpectation, reauthExpectation}
	assert.NoError(t, setOCSExpectations(expectations, nil))
	tr.AuthenticateAndAssertSuccess(imsi)

	// by this point we should be already restricted since credit was suspended
	assert.Eventually(t, tr.WaitForEnforcementStatsForRule(imsi, "restrict-pass-user"), time.Minute, 2*time.Second)

	// we need to generate traffic and validate it goes through restrict rule
	req := &cwfprotos.GenTrafficRequest{
		Imsi:    imsi,
		Volume:  &wrappers.StringValue{Value: "2M"},
		Timeout: 20,
	}
	_, err := tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	// Check that enforcement stats flow was not removed and flow data hit restrict rule
	tr.AssertPolicyUsage(imsi, "restrict-pass-user", uint64(math.Round(1.5*MegaBytes)), 3*MegaBytes+Buffer)
	// check static rule is gone
	assert.Eventually(t, tr.WaitForNoEnforcementStatsForRule(imsi, "static-pass-all-ocs1"), time.Minute, 2*time.Second)

	// Send ReAuth Request to update quota
	raa, err := sendChargingReAuthRequestEntireSession(imsi)
	assert.NoError(t, err)
	assert.Eventually(t, tr.WaitForChargingReAuthToProcess(raa, imsi), time.Minute, 2*time.Second)

	assert.NotNil(t, raa)
	if raa != nil {
		// Check ReAuth success
		assert.Equal(t, diam.LimitedSuccess, int(raa.ResultCode))
	}

	// Assert that a CCR-I and reauth were sent
	tr.AssertAllGyExpectationsMetNoError()

	// Wait for service activation
	assert.Eventually(t, tr.WaitForNoEnforcementStatsForRule(imsi, "restrict-pass-user"), time.Minute, 2*time.Second)

	req = &cwfprotos.GenTrafficRequest{
		Imsi:    imsi,
		Volume:  &wrappers.StringValue{Value: "2M"},
		Timeout: 20,
	}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	// TODO: uncomment once we fix passing the ip to pipelined for cwf
	// Check that enforcement stats flow was not removed and data passed
	//tr.AssertPolicyUsage(imsi, "static-pass-all-ocs1", uint64(math.Round(1.5*MegaBytes)), 3*MegaBytes+Buffer)
	//assert.Nil(t, policyUsage["IMSI"+imsi]["restrict-pass-user"], fmt.Sprintf("Policy usage restrict-pass-user for imsi: %v was NOT removed", imsi))

	// trigger disconnection
	tr.DisconnectAndAssertSuccess(imsi)
	tr.AssertEventuallyAllRulesRemovedAfterDisconnect(imsi)
}

// - Set an expectation for a CCR-I to be sent up to OCS, to which it will
//   respond with a quota grant of 4M with two rules.
// - Generate traffic and assert the CCR-I is received.
// - Set an expectation for a CCR-U with >80% of data usage to be sent up to
// 	 OCS, to which it will respond with an ERROR CODE
// - Send an CCA-U with a 4012 code transient failure which should trigger suspend that credit
// - Assert that UE flows for one rule are delete.
// - Assert that UE flows for the other rule are still valid
func TestGyWithTransientErrorCode(t *testing.T) {
	fmt.Println("\nRunning TestGyWithTransientErrorCode...")

	tr, ruleManager, ue := ocsTestSetup(t)
	imsi := ue.GetImsi()
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearOCSMockDriver())
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()

	// CCR-I
	quotaGrant := &fegprotos.QuotaGrant{
		RatingGroup: 1,
		GrantedServiceUnit: &fegprotos.Octets{
			TotalOctets: 1 * MegaBytes,
		},
		IsFinalCredit: false,
		ResultCode:    diam.Success,
	}
	initRequest := protos.NewGyCCRequest(imsi, protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGyCCAnswer(diam.Success).SetQuotaGrant(quotaGrant)
	initExpectation := protos.NewGyCreditControlExpectation().Expect(initRequest).Return(initAnswer)
	// grant with DiameterCreditLimitReached
	quotaGrantCreditLimitReached := &fegprotos.QuotaGrant{
		RatingGroup: 1,
		GrantedServiceUnit: &fegprotos.Octets{
			TotalOctets: 0 * MegaBytes,
		},
		IsFinalCredit: false,
		ResultCode:    diameter.DiameterCreditLimitReached,
	}
	// CCR-U  with ERROR CODE 4012 (DiameterCreditLimitReached)
	updateRequest1 := protos.NewGyCCRequest(imsi, protos.CCRequestType_UPDATE)
	updateAnswer1 := protos.NewGyCCAnswer(diam.Success).SetQuotaGrant(quotaGrantCreditLimitReached)
	updateExpectation1 := protos.NewGyCreditControlExpectation().Expect(updateRequest1).Return(updateAnswer1)
	// Load expectations into OCS
	expectations := []*protos.GyCreditControlExpectation{initExpectation, updateExpectation1}
	assert.NoError(t, setOCSExpectations(expectations, nil)) // We only expect one single CCR-U to be sent
	tr.AuthenticateAndAssertSuccess(imsi)

	// we need to generate over 80% but less than 100%  trigger a CCR update without triggering termination
	req := &cwfprotos.GenTrafficRequest{Imsi: imsi}
	_, err := tr.GenULTrafficBasedOnPolicyUsage(req, "static-pass-all-ocs2", 1*MegaBytes, 40*time.Second)
	assert.NoError(t, err)

	// On suspension rules must be removed from pipelined but not from the session
	fmt.Println("RG 1 should now be suspended")
	assert.Eventually(t, tr.WaitForNoEnforcementStatsForRule(imsi, "static-pass-all-ocs2"), 20*time.Second, 2*time.Second)

	tr.AssertAllGyExpectationsMetNoError()

	// un-suspend credit
	// set expectations on the OCS
	updateRequestAfterRar := protos.NewGyCCRequest(imsi, protos.CCRequestType_UPDATE)
	updateAnswerAfterRar := protos.NewGyCCAnswer(diam.Success).SetQuotaGrant(quotaGrant)
	updateExpectationAfterRar := protos.NewGyCreditControlExpectation().Expect(updateRequestAfterRar).Return(updateAnswerAfterRar)
	expectations = []*protos.GyCreditControlExpectation{updateExpectationAfterRar}
	assert.NoError(t, setOCSExpectations(expectations, nil))

	// Send ReAuth Request to update quota
	raa, err := sendChargingReAuthRequest(imsi, 1)
	assert.NoError(t, err)
	assert.Eventually(t, tr.WaitForChargingReAuthToProcess(raa, imsi), time.Minute, 2*time.Second)
	assert.NotNil(t, raa)
	if raa != nil {
		// Check ReAuth success
		assert.Equal(t, diam.LimitedSuccess, int(raa.ResultCode))
	}
	// check the rule is back to the stats
	assert.Eventually(t, tr.WaitForEnforcementStatsForRule(imsi, "static-pass-all-ocs2"), 20*time.Second, 2*time.Second)

	// Assert CCR-U after Rar were sent to the OCS
	tr.AssertAllGyExpectationsMetNoError()

	fmt.Println("RG 1 should now be UN-suspended")
	// check rule is able to track traffic again
	req = &cwfprotos.GenTrafficRequest{Imsi: imsi}
	_, err = tr.GenULTrafficBasedOnPolicyUsage(req, "static-pass-all-ocs2", 200 * KiloBytes, 40 * time.Second)
	assert.NoError(t, err)

	// trigger disconnection
	tr.DisconnectAndAssertSuccess(imsi)
	tr.AssertEventuallyAllRulesRemovedAfterDisconnect(imsi)
}

// - Set an expectation for a CCR-I to be sent up to OCS, to which it will
//   respond with a quota grant of 4M.
//   Generate traffic and assert the CCR-I is received.
// - Generate traffic over 80% and under 100% not to trigger termination
// - Send an CCA-U with a 5xxx code which should trigger termination
// - Assert that UE flows are deleted.
// - Expect a CCR-T, trigger a UE disconnect, and assert the CCR-T is received.
func TestGyWithPermanentErrorCode(t *testing.T) {
	fmt.Println("\nRunning TestGyWithPermanentErrorCode...")

	tr, ruleManager, ue := ocsTestSetup(t)
	imsi := ue.GetImsi()
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearOCSMockDriver())
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()

	// CCR-I
	quotaGrant := &fegprotos.QuotaGrant{
		RatingGroup: 1,
		GrantedServiceUnit: &fegprotos.Octets{
			TotalOctets: 5 * MegaBytes,
		},
		IsFinalCredit: false,
		ResultCode:    diam.Success,
	}
	initRequest := protos.NewGyCCRequest(imsi, protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGyCCAnswer(diam.Success).SetQuotaGrant(quotaGrant)
	initExpectation := protos.NewGyCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	// grant with any 5xxx error (permanent error)
	quotaGrantCreditLimitReached := &fegprotos.QuotaGrant{
		RatingGroup: 1,
		GrantedServiceUnit: &fegprotos.Octets{
			TotalOctets: 0 * MegaBytes,
		},
		IsFinalCredit: false,
		ResultCode:    diameter.DiameterRatingFailed,
	}

	// CCR-U  with ERROR CODE 4012 (DiameterCreditLimitReached)
	updateRequest1 := protos.NewGyCCRequest(imsi, protos.CCRequestType_UPDATE)
	updateAnswer1 := protos.NewGyCCAnswer(diam.Success).SetQuotaGrant(quotaGrantCreditLimitReached)
	updateExpectation1 := protos.NewGyCreditControlExpectation().Expect(updateRequest1).Return(updateAnswer1)

	// CCR-T
	terminateRequest := protos.NewGyCCRequest(imsi, protos.CCRequestType_TERMINATION)
	terminateAnswer := protos.NewGyCCAnswer(diam.Success)
	terminateExpectation := protos.NewGyCreditControlExpectation().Expect(terminateRequest).Return(terminateAnswer)

	// Load expectations into OCS
	expectations := []*protos.GyCreditControlExpectation{initExpectation, updateExpectation1, terminateExpectation}
	assert.NoError(t, setOCSExpectations(expectations, nil)) // We only expect one single CCR-U to be sent
	tr.AuthenticateAndAssertSuccess(imsi)
	// First wait until we see the original static-pass-all-ocs2 show up
	assert.Eventually(t,
		tr.WaitForEnforcementStatsForRule(imsi, "static-pass-all-ocs2"),
		10*time.Second, 2*time.Second)

	// we need to generate over 80% but less than 100%  trigger a CCR update without triggering termination
	req := &cwfprotos.GenTrafficRequest{
		Imsi:   imsi,
		Volume: &wrappers.StringValue{Value: "10M"},
	}
	tr.GenULTraffic(req)

	// Check that enforcement stats flow is removed
	tr.AssertEventuallyAllRulesRemovedAfterDisconnect(imsi)

	// Assert that we saw a Terminate request
	tr.AssertAllGyExpectationsMetNoError()
}
