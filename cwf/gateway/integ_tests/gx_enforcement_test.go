// +build all gx

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
	cwfprotos "magma/cwf/cloud/go/protos"
	"magma/feg/cloud/go/protos"
	fegProtos "magma/feg/cloud/go/protos"
	lteProtos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/policydb/obsidian/models"

	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
)

// - Set an expectation for a  CCR-I to be sent up to PCRF, to which it will
//   respond with a rule install (usage-enforcement-static-pass-all), 1MB of
//   quota.
//   Generate traffic and assert the CCR-I is received.
// - Set an expectation for a CCR-U with >80% of data usage to be sent up to
// 	 PCRF, to which it will response with more quota.
//   Generate traffic and assert the CCR-U is received.
// - Generate traffic to put traffic through the newly installed rule.
//   Assert that there's > 0 data usage in the rule.
// - Expect a CCR-T, trigger a UE disconnect, and assert the CCR-T is received.
func TestGxUsageReportEnforcement(t *testing.T) {
	fmt.Println("\nRunning TestGxUsageReportEnforcement...")
	tr := NewTestRunner(t)
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)
	assert.NoError(t, usePCRFMockDriver())
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearPCRFMockDriver())
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)
	imsi := ues[0].GetImsi()

	err = ruleManager.AddStaticPassAllToDB("usage-enforcement-static-pass-all", "mkey1", 0, models.PolicyRuleTrackingTypeONLYPCRF, 3)
	assert.NoError(t, err)
	tr.WaitForPoliciesToSync()

	usageMonitorInfo := getUsageInformation("mkey1", 1*MegaBytes)

	initRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGxCCAnswer(diam.Success).
		SetStaticRuleInstalls([]string{"usage-enforcement-static-pass-all"}, []string{}).
		SetUsageMonitorInfo(usageMonitorInfo)
	initExpectation := protos.NewGxCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	// We expect an update request with some usage update (probably around 80-100% of the given quota)
	updateRequest1 := protos.NewGxCCRequest(imsi, protos.CCRequestType_UPDATE).
		SetUsageMonitorReport(usageMonitorInfo).
		SetUsageReportDelta(uint64(math.Round(0.2 * 1 * MegaBytes))).
		SetEventTrigger(int32(lteProtos.EventTrigger_USAGE_REPORT))
	updateAnswer1 := protos.NewGxCCAnswer(diam.Success).SetUsageMonitorInfo(usageMonitorInfo)
	updateExpectation1 := protos.NewGxCreditControlExpectation().Expect(updateRequest1).Return(updateAnswer1)
	expectations := []*protos.GxCreditControlExpectation{initExpectation, updateExpectation1}
	// On unexpected requests, just return the default update answer
	assert.NoError(t, setPCRFExpectations(expectations, updateAnswer1))

	tr.AuthenticateAndAssertSuccess(imsi)
	// First wait until we see the original static-pass-all-ocs2 show up
	assert.Eventually(t,
		tr.WaitForEnforcementStatsForRule(imsi, "usage-enforcement-static-pass-all"), time.Minute, 2*time.Second)
	fmt.Println("CCR-I exchanged installed usage-enforcement-static-pass-all")

	req := &cwfprotos.GenTrafficRequest{Imsi: imsi, Volume: &wrappers.StringValue{Value: "900K"}}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()
	req = &cwfprotos.GenTrafficRequest{Imsi: imsi, Volume: &wrappers.StringValue{Value: "200K"}}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	// Assert that enforcement_stats rules are properly installed and the right
	// amount of data was passed through
	tr.AssertPolicyUsage(imsi, "usage-enforcement-static-pass-all", 1, uint64(math.Round(1.2*MegaBytes+Buffer)))

	// Assert that a CCR-I and at least one CCR-U were sent up to the PCRF
	tr.AssertAllGxExpectationsMetNoError()

	// When we initiate a UE disconnect, we expect a terminate request to go up
	terminateRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_TERMINATION)
	terminateAnswer := protos.NewGxCCAnswer(diam.Success)
	terminateExpectation := protos.NewGxCreditControlExpectation().Expect(terminateRequest).Return(terminateAnswer)
	expectations = []*protos.GxCreditControlExpectation{terminateExpectation}
	assert.NoError(t, setPCRFExpectations(expectations, nil))

	tr.DisconnectAndAssertSuccess(imsi)
	tr.AssertEventuallyAllRulesRemovedAfterDisconnect(imsi)
	// Assert that we saw a Terminate request
	tr.AssertAllGxExpectationsMetNoError()
}

// - Set an expectation for a  CCR-I to be sent up to PCRF, to which it will
//   respond with a rule install (static-pass-all-1).
//   Generate traffic and assert the CCR-I is received.
// - Set an expectation for a CCR-U to be sent up to PCRF, to which it will
//   respond with a rule removal (static-pass-all-1) and rule install (static-pass-all-2).
//   Generate traffic and assert the CCR-U is received.
// - Generate traffic to put traffic through the newly installed rule.
//   Assert that there's > 0 data usage in the rule.
func TestGxMidSessionRuleRemovalWithCCA_U(t *testing.T) {
	fmt.Println("\nRunning TestGxMidSessionRuleRemovalWithCCA_U...")

	tr := NewTestRunner(t)
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)
	assert.NoError(t, usePCRFMockDriver())
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearPCRFMockDriver())
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)
	imsi := ues[0].GetImsi()

	err = ruleManager.AddStaticPassAllToDB("static-pass-all-1", "mkey1", 0, models.PolicyRuleTrackingTypeONLYPCRF, 100)
	assert.NoError(t, err)
	err = ruleManager.AddStaticPassAllToDB("static-pass-all-2", "mkey2", 0, models.PolicyRuleTrackingTypeONLYPCRF, 150)
	assert.NoError(t, err)
	err = ruleManager.AddStaticPassAllToDB("static-pass-all-3", "mkey2", 0, models.PolicyRuleTrackingTypeONLYPCRF, 200)
	assert.NoError(t, err)
	err = ruleManager.AddBaseNameMappingToDB("base-1", []string{"static-pass-all-3"})
	tr.WaitForPoliciesToSync()

	usageMonitorInfo := getUsageInformation("mkey1", 2*MegaBytes)

	initRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGxCCAnswer(diam.Success).
		SetStaticRuleInstalls([]string{"static-pass-all-1"}, []string{"base-1"}).
		SetUsageMonitorInfo(usageMonitorInfo)
	initExpectation := protos.NewGxCreditControlExpectation().Expect(initRequest).Return(initAnswer)
	// Remove the high priority Rule
	defaultUpdateAnswer := protos.NewGxCCAnswer(diam.Success).SetUsageMonitorInfo(usageMonitorInfo)
	expectations := []*protos.GxCreditControlExpectation{initExpectation}
	// On unexpected requests, just return some quota
	assert.NoError(t, setPCRFExpectations(expectations, defaultUpdateAnswer))

	tr.AuthenticateAndAssertSuccess(imsi)
	assert.Eventually(t, tr.WaitForEnforcementStatsForRule(imsi, "static-pass-all-1", "static-pass-all-3"), time.Minute, 2*time.Second)

	// Pass a small amount, but not enough to trigger a CCR-U
	req := &cwfprotos.GenTrafficRequest{
		Imsi:    imsi,
		Volume:  &wrappers.StringValue{Value: "1M"},
		Bitrate: &wrappers.StringValue{Value: "30M"},
		Timeout: 30,
	}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)

	// At this point both static-pass-all-1 & static-pass-all-3 are installed.
	// Since static-pass-all-1 has higher precedence, it will get hit.
	assert.Eventually(t,
		tr.WaitForEnforcementStatsForRuleGreaterThan(imsi, "static-pass-all-1", 1), time.Minute, 2*time.Second)

	// Assert that a CCR-I was sent up to the PCRF
	tr.AssertAllGxExpectationsMetNoError()

	updateRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_UPDATE).
		SetUsageMonitorReport(usageMonitorInfo).
		SetUsageReportDelta(1 * MegaBytes).
		SetEventTrigger(int32(lteProtos.EventTrigger_USAGE_REPORT))
	updateAnswer := protos.NewGxCCAnswer(diam.Success).SetUsageMonitorInfo(usageMonitorInfo).
		SetStaticRuleInstalls([]string{"static-pass-all-2"}, []string{}).
		SetStaticRuleRemovals([]string{"static-pass-all-1"}, []string{"base-1"})
	updateExpectation := protos.NewGxCreditControlExpectation().Expect(updateRequest).Return(updateAnswer)
	expectations = []*protos.GxCreditControlExpectation{updateExpectation}
	// On unexpected requests, just return some quota
	assert.NoError(t, setPCRFExpectations(expectations, defaultUpdateAnswer))

	fmt.Println("Generating traffic again to trigger a CCR/A-U so that 'static-pass-all-1' gets removed")
	// Generate traffic to trigger the CCR-U so that the rule removal/install happens
	req = &cwfprotos.GenTrafficRequest{
		Imsi:    imsi,
		Volume:  &wrappers.StringValue{Value: "1M"},
		Bitrate: &wrappers.StringValue{Value: "20M"},
		Timeout: 60,
	}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)

	assert.Eventually(t, tr.WaitForNoEnforcementStatsForRule(imsi, "static-pass-all-1", "static-pass-all-3"), 1*time.Minute, 2*time.Second)
	assert.Eventually(t, tr.WaitForEnforcementStatsForRule(imsi, "static-pass-all-2"), 1*time.Minute, 2*time.Second)

	fmt.Println("Generating traffic again to put data through static-pass-all-2")
	req = &cwfprotos.GenTrafficRequest{
		Imsi:    imsi,
		Volume:  &wrappers.StringValue{Value: "1M"},
		Bitrate: &wrappers.StringValue{Value: "30M"},
		Timeout: 30,
	}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	assert.Eventually(t,
		tr.WaitForEnforcementStatsForRuleGreaterThan(imsi, "static-pass-all-2", 1), time.Minute, 2*time.Second)

	// Assert that we sent back a CCA-Update with RuleRemovals
	tr.AssertAllGxExpectationsMetNoError()

	tr.DisconnectAndAssertSuccess(imsi)
	tr.AssertEventuallyAllRulesRemovedAfterDisconnect(imsi)
}

// - Set an expectation for a  CCR-I to be sent up to PCRF, to which it will
//   respond with a rule install (static-pass-all-1).
// - Set an expectation for a CCR-U to be sent up to PCRF, to which it will
//   respond with a rule install (static-pass-all-2), with activation (now + X sec)
//   and deactivation time (activation + Y sec) specified.
// - Generate traffic to trigger a CCR-U. Check policy usage and assert
//   static-pass-all-2 is not installed.
// - Sleep for X seconds and check policy usage again. Assert that
//   static-pass-all-2 is installed.
// - Sleep for Y seconds and check policy usage again. Assert that
//   static-pass-all-2 is uninstalled.
// Note: things might get weird if there are clock skews
func TestGxRuleInstallTime(t *testing.T) {
	t.Skip()
	fmt.Println("\nRunning TestGxRuleInstallTime...")

	tr := NewTestRunner(t)
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)
	assert.NoError(t, usePCRFMockDriver())
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearPCRFMockDriver())
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)
	imsi := ues[0].GetImsi()

	err = ruleManager.AddStaticPassAllToDB("static-pass-all-1", "mkey1", 0, models.PolicyRuleTrackingTypeONLYPCRF, 200)
	assert.NoError(t, err)
	err = ruleManager.AddStaticPassAllToDB("static-pass-all-2", "mkey1", 0, models.PolicyRuleTrackingTypeONLYPCRF, 100)
	assert.NoError(t, err)
	tr.WaitForPoliciesToSync()

	usageMonitorInfo := getUsageInformation("mkey1", 250*KiloBytes)
	initRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGxCCAnswer(diam.Success).
		SetStaticRuleInstalls([]string{"static-pass-all-1"}, nil).
		SetUsageMonitorInfo(usageMonitorInfo)
	initExpectation := protos.NewGxCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	now := time.Now().Round(1 * time.Second)
	timeUntilActivation := 8 * time.Second
	activation := now.Add(timeUntilActivation)
	pActivation, err := ptypes.TimestampProto(activation)
	assert.NoError(t, err)
	timeUntilDeactivation := 8 * time.Second
	deactivation := activation.Add(timeUntilDeactivation)
	pDeactivation, err := ptypes.TimestampProto(deactivation)
	assert.NoError(t, err)

	updateRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_UPDATE)
	updateAnswer := protos.NewGxCCAnswer(diam.Success).
		SetUsageMonitorInfo(usageMonitorInfo).
		SetStaticRuleInstalls([]string{"static-pass-all-2"}, nil).
		SetRuleActivationTime(pActivation).
		SetRuleDeactivationTime(pDeactivation)
	updateExpectation := protos.NewGxCreditControlExpectation().Expect(updateRequest).Return(updateAnswer)
	defaultUpdateAnswer := protos.NewGxCCAnswer(diam.Success).SetUsageMonitorInfo(usageMonitorInfo)
	expectations := []*protos.GxCreditControlExpectation{initExpectation, updateExpectation}
	// On unexpected requests, just return some quota
	assert.NoError(t, setPCRFExpectations(expectations, defaultUpdateAnswer))

	tr.AuthenticateAndAssertSuccess(imsi)

	// Generate over the given quota
	req := &cwfprotos.GenTrafficRequest{Imsi: imsi, Volume: &wrappers.StringValue{Value: "300K"}}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()
	recordsBySubID, err := tr.GetPolicyUsage()
	// only static-pass-all-1 should be installed
	assert.NotNil(t, recordsBySubID[prependIMSIPrefix(imsi)]["static-pass-all-1"])
	assert.Nil(t, recordsBySubID[prependIMSIPrefix(imsi)]["static-pass-all-2"])

	fmt.Printf("Waiting %v for rule activation\n", timeUntilActivation)
	time.Sleep(timeUntilActivation)
	recordsBySubID, err = tr.GetPolicyUsage()
	// both rules should exist
	assert.NotNil(t, recordsBySubID[prependIMSIPrefix(imsi)]["static-pass-all-1"])
	assert.NotNil(t, recordsBySubID[prependIMSIPrefix(imsi)]["static-pass-all-2"])

	fmt.Printf("Waiting %v for rule deactivation\n", timeUntilDeactivation)
	time.Sleep(timeUntilDeactivation)
	recordsBySubID, err = tr.GetPolicyUsage()
	// static-pass-all-2 should be gone
	assert.NotNil(t, recordsBySubID[prependIMSIPrefix(imsi)]["static-pass-all-1"])
	assert.Nil(t, recordsBySubID[prependIMSIPrefix(imsi)]["static-pass-all-2"])

	tr.AssertAllGxExpectationsMetNoError()

	tr.DisconnectAndAssertSuccess(imsi)
	tr.AssertEventuallyAllRulesRemovedAfterDisconnect(imsi)
}

//TestGxAbortSessionRequest
// This test verifies the abort session request
// Here we initially setup a session and install a pass all rule
// We then invoke abort session request from pcrf and expect the
// ASR to complete without any error and all the rules associated with
// that session to be cleaned up
func TestGxAbortSessionRequest(t *testing.T) {
	fmt.Println("\nRunning TestGxAbortSessionRequest")
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
	imsi := ues[0].GetImsi()
	ki := rand.Intn(1000000)
	monitorKey := fmt.Sprintf("monitor-ASR-%d", ki)
	ruleKey := fmt.Sprintf("static-ASR1_CCI-%d", ki)
	rule := getStaticPassAll(ruleKey, monitorKey, 0, models.PolicyRuleTrackingTypeONLYPCRF, 3, nil)
	err = ruleManager.AddStaticRuleToDB(rule)

	err = ruleManager.AddUsageMonitor(imsi, monitorKey, 2*MegaBytes, 1*MegaBytes)
	err = ruleManager.AddRulesToPCRF(imsi, []string{ruleKey}, []string{})
	assert.NoError(t, err)
	tr.WaitForPoliciesToSync()

	tr.AuthenticateAndAssertSuccess(imsi)
	tr.WaitForEnforcementStatsToSync()

	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	assert.NotEmpty(t, recordsBySubID[prependIMSIPrefix(imsi)][ruleKey])

	asa, err := sendPolicyAbortSession(
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
	tr.AssertEventuallyAllRulesRemovedAfterDisconnect(imsi)
}

// - Set an expectation for a CCR-I to be sent up to PCRF, to which it will
//   respond with a rule install (revalidation-time-static-pass-all) with
//   a revalidation time set to now + 10s and an event trigger REVALIDATION_TIMEOUT.
// - Set an expectation for a CCR-U to be sent up to PCRF, upon revalidation
//   timer expiration
// - Check policy usage and assert revalidation-time-static-pass-all is installed and
//   no traffic passed
// Note: things might get weird if there are clock skews
func TestGxRevalidationTime(t *testing.T) {
	t.Skip("Skipping due to current test flakiness, investigating...")
	fmt.Println("\nRunning TestGxRevalidationTime...")

	tr := NewTestRunner(t)
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)
	assert.NoError(t, usePCRFMockDriver())
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearPCRFMockDriver())
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)
	imsi := ues[0].GetImsi()

	err = ruleManager.AddStaticPassAllToDB("revalidation-time-static-pass-all", "mkey1", 0, models.PolicyRuleTrackingTypeONLYPCRF, 1)
	assert.NoError(t, err)
	tr.WaitForPoliciesToSync()

	timeUntilRevalidation := 8 * time.Second
	now := time.Now().Round(1 * time.Second)
	revalidationTime, err := ptypes.TimestampProto(now.Add(timeUntilRevalidation))
	assert.NoError(t, err)

	initRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGxCCAnswer(diam.Success).
		SetStaticRuleInstalls([]string{"revalidation-time-static-pass-all"}, []string{}).
		SetRevalidationTime(revalidationTime).
		SetEventTriggers([]uint32{RevalidationTimeoutEvent})
	initExpectation := protos.NewGxCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	// We expect an update request with some usage update after revalidation timer expires
	updateRequest1 := protos.NewGxCCRequest(imsi, protos.CCRequestType_UPDATE).
		SetEventTrigger(int32(lteProtos.EventTrigger_REVALIDATION_TIMEOUT))
	updateAnswer1 := protos.NewGxCCAnswer(diam.Success)
	updateExpectation1 := protos.NewGxCreditControlExpectation().Expect(updateRequest1).Return(updateAnswer1)
	expectations := []*protos.GxCreditControlExpectation{initExpectation, updateExpectation1}
	// On unexpected requests, just return the default update answer
	assert.NoError(t, setPCRFExpectations(expectations, updateAnswer1))

	tr.AuthenticateAndAssertSuccess(imsi)
	// First wait until we see the original static-pass-all-ocs2 show up
	assert.Eventually(t,
		tr.WaitForEnforcementStatsForRule(imsi, "revalidation-time-static-pass-all"),
		10*time.Second, 2*time.Second)

	waitingTime := timeUntilRevalidation + (5 * time.Second)
	fmt.Printf("Waiting %v seconds for revalidation timer expiration\n", waitingTime)
	time.Sleep(waitingTime) // give an extra few seconds for error

	// Assert that a CCR-I and at least one CCR-U were sent up to the PCRF
	tr.AssertAllGxExpectationsMetNoError()

	// When we initiate a UE disconnect, we expect a terminate request to go up
	terminateRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_TERMINATION)
	terminateAnswer := protos.NewGxCCAnswer(diam.Success)
	terminateExpectation := protos.NewGxCreditControlExpectation().Expect(terminateRequest).Return(terminateAnswer)
	expectations = []*protos.GxCreditControlExpectation{terminateExpectation}
	assert.NoError(t, setPCRFExpectations(expectations, nil))

	tr.DisconnectAndAssertSuccess(imsi)
	tr.AssertEventuallyAllRulesRemovedAfterDisconnect(imsi)

	// Assert that we saw a Terminate request
	tr.AssertAllGxExpectationsMetNoError()
}
