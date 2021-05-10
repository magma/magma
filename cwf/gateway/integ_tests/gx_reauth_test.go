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
	"testing"
	"time"

	cwfprotos "magma/cwf/cloud/go/protos"
	fegprotos "magma/feg/cloud/go/protos"
	"magma/lte/cloud/go/services/policydb/obsidian/models"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
)

func gxReAuthTestSetup(t *testing.T) (*TestRunner, *RuleManager, *cwfprotos.UEConfig) {
	tr := NewTestRunner(t)
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)
	ue := ues[0]

	// Install two static rules and a rule base with an additional static rule
	err = ruleManager.AddUsageMonitor(ue.GetImsi(), "raakey1", 500*KiloBytes, 250*KiloBytes)
	assert.NoError(t, err)
	err = ruleManager.AddStaticPassAllToDB("static-pass-all-raa1", "raakey1", 0, models.PolicyRuleTrackingTypeONLYPCRF, 100)
	assert.NoError(t, err)
	err = ruleManager.AddStaticPassAllToDB("static-pass-all-raa2", "raakey2", 0, models.PolicyRuleTrackingTypeONLYPCRF, 150)
	assert.NoError(t, err)
	err = ruleManager.AddStaticPassAllToDB("static-pass-all-raa3", "raakey2", 0, models.PolicyRuleTrackingTypeONLYPCRF, 200)
	assert.NoError(t, err)
	err = ruleManager.AddBaseNameMappingToDB("base-raa1", []string{"static-pass-all-raa3"})
	assert.NoError(t, err)

	// Apply a dynamic rule that points to the static rules above
	err = ruleManager.AddRulesToPCRF(ue.GetImsi(), []string{"static-pass-all-raa1", "static-pass-all-raa2"}, []string{"base-raa1"})
	assert.NoError(t, err)

	tr.WaitForPoliciesToSync()
	return tr, ruleManager, ue
}

// - Install two static rules "static-pass-all-raa1" and "static-pass-all-raa2"
//   and a rule base "base-raa1"
// - Generate traffic and assert that there's > 0 data usage for the rule with the
//   highest priority.
// - Send a PCRF ReAuth request to delete "static-pass-all-raa2" and "base-raa1"
//   and assert that the response is successful
// - Assert that the requested rules were removed
func TestGxReAuthWithMidSessionPolicyRemoval(t *testing.T) {
	fmt.Println("\nRunning TestGxReAuthWithMidSessionPolicyRemoval...")

	tr, ruleManager, ue := gxReAuthTestSetup(t)
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()

	imsi := ue.GetImsi()

	tr.AuthenticateAndAssertSuccess(imsi)

	// Generate over 80% of the quota to trigger a CCR Update
	req := &cwfprotos.GenTrafficRequest{
		Imsi:    imsi,
		Volume:  &wrappers.StringValue{Value: "450K"},
		Bitrate: &wrappers.StringValue{Value: "20M"},
		Timeout: 60,
	}
	_, err := tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	// Check that enforcement flow is installed and traffic is less than the quota
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	record := recordsBySubID["IMSI"+imsi]["static-pass-all-raa1"]
	assert.NotNil(t, record,
		fmt.Sprintf("Policy usage record for imsi: %v rule: 'static-pass-all-raa1' does not exist", imsi))
	assert.True(t, record.BytesTx > uint64(0), fmt.Sprintf("%s did not pass any data", record.RuleId))
	if record.BytesTx == uint64(0) && record.BytesRx == uint64(0) {
		dumpPipelinedState(tr)
	}
	// Send ReAuth Request to update quota
	rulesRemoval := &fegprotos.RuleRemovals{
		RuleNames:     []string{"static-pass-all-raa2"},
		RuleBaseNames: []string{"base-raa1"},
	}
	raa, err := sendPolicyReAuthRequest(
		&fegprotos.PolicyReAuthTarget{Imsi: imsi, RulesToRemove: rulesRemoval},
	)
	assert.NoError(t, err)

	assert.Eventually(t, tr.WaitForPolicyReAuthToProcess(raa, imsi), time.Minute, 2*time.Second)

	assert.NotNil(t, raa)
	if raa != nil {
		assert.Equal(t, diam.Success, int(raa.ResultCode))
	}

	// Check that enforcement flows were deleted for rule 2 and 3
	tr.WaitForEnforcementStatsToSync()

	recordsBySubID, err = tr.GetPolicyUsage()
	assert.NoError(t, err)

	record1 := recordsBySubID["IMSI"+imsi]["static-pass-all-raa1"]
	assert.NotNil(t, record1, fmt.Sprintf("Policy usage record for imsi: %v was removed", imsi))
	record2 := recordsBySubID["IMSI"+imsi]["static-pass-all-raa2"]
	assert.Nil(t, record2, fmt.Sprintf("Policy usage record for imsi: %v was not removed", imsi))
	record3 := recordsBySubID["IMSI"+imsi]["static-pass-all-raa3"]
	assert.Nil(t, record3, fmt.Sprintf("Policy usage record for imsi: %v was not removed", imsi))

	// Trigger disconnection
	tr.DisconnectAndAssertSuccess(imsi)
	tr.AssertEventuallyAllRulesRemovedAfterDisconnect(imsi)
}

// - Install two static rules "static-pass-all-raa1" and "static-pass-all-raa2"
//   and a rule base "base-raa1"
// - Generate traffic and assert that there's > 0 data usage for the rule with the
//   highest priority.
// - Send a PCRF ReAuth request to delete all the installed rules and assert
//   that the response is successful
// - Assert that the requested rules were removed
// - Assert that session was deleted
func TestGxReAuthWithMidSessionPoliciesRemoval(t *testing.T) {
	fmt.Println("\nRunning TestGxReAuthWithMidSessionPoliciesRemoval...")

	tr, ruleManager, ue := gxReAuthTestSetup(t)
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()

	imsi := ue.GetImsi()

	tr.AuthenticateAndAssertSuccess(imsi)
	// Wait for flows to be installed
	assert.Eventually(t, tr.WaitForEnforcementStatsForRule(
		imsi, "static-pass-all-raa1", "static-pass-all-raa2"), time.Minute, 2*time.Second)

	// Generate over 80% of the quota to trigger a CCR Update
	req := &cwfprotos.GenTrafficRequest{
		Imsi:    imsi,
		Volume:  &wrappers.StringValue{Value: "450K"},
		Bitrate: &wrappers.StringValue{Value: "20M"},
		Timeout: 60,
	}
	_, err := tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	// Check that enforcement flow is installed and traffic is less than the quota
	tr.AssertPolicyUsage(imsi, "static-pass-all-raa1", 100*KiloBytes, 500*KiloBytes+Buffer)

	// Send ReAuth Request to update quota
	rulesRemoval := &fegprotos.RuleRemovals{
		RuleNames:     []string{"static-pass-all-raa1", "static-pass-all-raa2"},
		RuleBaseNames: []string{"base-raa1"},
	}
	raa, err := sendPolicyReAuthRequest(
		&fegprotos.PolicyReAuthTarget{Imsi: imsi, RulesToRemove: rulesRemoval},
	)
	assert.NoError(t, err)
	assert.Eventually(t, tr.WaitForPolicyReAuthToProcess(raa, imsi), time.Minute, 2*time.Second)
	assert.NotNil(t, raa)
	if raa != nil {
		assert.Equal(t, diam.Success, int(raa.ResultCode))
	}
	// Check that all UE mac flows are deleted
	tr.WaitForEnforcementStatsToSync()

	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)

	record1 := recordsBySubID["IMSI"+imsi]["static-pass-all-raa1"]
	assert.Nil(t, record1, fmt.Sprintf("Policy usage record for imsi: %v was not removed", imsi))
	record2 := recordsBySubID["IMSI"+imsi]["static-pass-all-raa2"]
	assert.Nil(t, record2, fmt.Sprintf("Policy usage record for imsi: %v was not removed", imsi))
	record3 := recordsBySubID["IMSI"+imsi]["static-pass-all-raa3"]
	assert.Nil(t, record3, fmt.Sprintf("Policy usage record for imsi: %v was not removed", imsi))

	// trigger disconnection
	tr.DisconnectAndAssertSuccess(imsi)
	tr.AssertEventuallyAllRulesRemovedAfterDisconnect(imsi)
}

// - Install two static rules "static-pass-all-raa1" and "static-pass-all-raa2"
//   and a rule base "base-raa1"
// - Generate traffic and assert that there's > 0 data usage for the rule with the
//   highest priority.
// - Send a PCRF ReAuth request to install a new pass all rule with higher priority
//   and assert that the response is successful
// - Generate traffic and assert that there's > 0 data usage for the newly installed
//   rule.
func TestGxReAuthWithMidSessionPolicyInstall(t *testing.T) {
	fmt.Println("\nRunning TestGxReAuthWithMidSessionPolicyInstall...")

	tr, ruleManager, ue := gxReAuthTestSetup(t)
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()
	imsi := ue.GetImsi()

	tr.AuthenticateAndAssertSuccess(imsi)
	// Wait for flows to be installed
	assert.Eventually(t, tr.WaitForEnforcementStatsForRule(
		imsi, "static-pass-all-raa1", "static-pass-all-raa2"), time.Minute, 2*time.Second)

	// Generate over 80% of the quota to trigger a CCR Update
	req := &cwfprotos.GenTrafficRequest{
		Imsi:    imsi,
		Volume:  &wrappers.StringValue{Value: "450K"},
		Bitrate: &wrappers.StringValue{Value: "20M"},
		Timeout: 60,
	}
	_, err := tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	// Check that enforcement flow is installed and traffic is less than the quota
	tr.AssertPolicyUsage(imsi, "static-pass-all-raa1", 0, 500*KiloBytes+Buffer)

	// Add a monitoring key
	err = ruleManager.AddUsageMonitor(imsi, "raakey3", 500*KiloBytes, 250*KiloBytes)
	assert.NoError(t, err)

	// Install a Pass-All Rule with higher priority using PolicyReAuth
	usageMonitoring := []*fegprotos.UsageMonitoringInformation{
		getUsageInformation("raakey3", 250*KiloBytes),
	}
	ruleDefinition := []*fegprotos.RuleDefinition{
		{
			RuleName:         "pcrf-reauth-raa1",
			Precedence:       50,
			MonitoringKey:    "raakey3",
			FlowDescriptions: []string{"permit out ip from any to any", "permit in ip from any to any"},
		},
	}
	ruleInstall := &fegprotos.RuleInstalls{
		RuleDefinitions: ruleDefinition,
	}
	raa, err := sendPolicyReAuthRequest(
		&fegprotos.PolicyReAuthTarget{
			Imsi:                 imsi,
			RulesToInstall:       ruleInstall,
			UsageMonitoringInfos: usageMonitoring,
		},
	)
	assert.NoError(t, err)
	assert.Eventually(t, tr.WaitForPolicyReAuthToProcess(raa, imsi), time.Minute, 2*time.Second)

	assert.NotNil(t, raa)
	if raa != nil {
		assert.Equal(t, diam.Success, int(raa.ResultCode))
	}
	// Generate more traffic
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	// Check that enforcement flow is installed and traffic is less than the quota
	tr.AssertPolicyUsage(imsi, "pcrf-reauth-raa1", 0, 500*KiloBytes+Buffer)

	// trigger disconnection
	tr.DisconnectAndAssertSuccess(imsi)
	tr.AssertEventuallyAllRulesRemovedAfterDisconnect(imsi)
}

// - Install two static rules "static-pass-all-raa1" and "static-pass-all-raa2"
//   and a rule base "base-raa1"
// - Generate traffic and assert that there's > 0 data usage for the rule with the
//   highest priority.
// - Send a PCRF ReAuth request to install a new pass all rule with second higher priority
//   and remove the rule with the highest priority
// - Assert that the response is successful
// - Generate traffic and assert that there's > 0 data usage for the newly installed
//   rule.
func TestGxReAuthWithMidSessionPolicyInstallAndRemoval(t *testing.T) {
	fmt.Println("\nRunning TestGxReAuthWithMidSessionPolicyInstallAndRemoval...")

	tr, ruleManager, ue := gxReAuthTestSetup(t)
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()
	imsi := ue.GetImsi()

	tr.AuthenticateAndAssertSuccess(imsi)
	// Wait for flows to be installed
	assert.Eventually(t, tr.WaitForEnforcementStatsForRule(
		imsi, "static-pass-all-raa1", "static-pass-all-raa2"), time.Minute, 2*time.Second)

	// Generate over 80% of the quota to trigger a CCR Update
	req := &cwfprotos.GenTrafficRequest{
		Imsi:    imsi,
		Volume:  &wrappers.StringValue{Value: "450K"},
		Bitrate: &wrappers.StringValue{Value: "20M"},
		Timeout: 60,
	}
	_, err := tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	// Check that enforcement flow is installed and traffic is less than the quota
	tr.AssertPolicyUsage(imsi, "static-pass-all-raa1", 0, 500*KiloBytes+Buffer)

	// Remove the rule with the highest priority
	rulesRemoval := &fegprotos.RuleRemovals{
		RuleNames:     []string{"static-pass-all-raa1"},
		RuleBaseNames: []string{""},
	}

	// Add a monitoring key
	err = ruleManager.AddUsageMonitor(ue.GetImsi(), "raakey4", 500*KiloBytes, 250*KiloBytes)
	assert.NoError(t, err)

	// Install a Pass-All Rule with higher priority using PolicyReAuth
	usageMonitoring := []*fegprotos.UsageMonitoringInformation{
		getUsageInformation("raakey5", 250*KiloBytes),
	}
	ruleDefinition := []*fegprotos.RuleDefinition{
		{
			RuleName:         "pcrf-reauth-raa2",
			Precedence:       125,
			MonitoringKey:    "raakey4",
			FlowDescriptions: []string{"permit out ip from any to any", "permit in ip from any to any"},
		},
	}
	ruleInstall := &fegprotos.RuleInstalls{
		RuleDefinitions: ruleDefinition,
	}
	raa, err := sendPolicyReAuthRequest(
		&fegprotos.PolicyReAuthTarget{
			Imsi:                 imsi,
			RulesToInstall:       ruleInstall,
			RulesToRemove:        rulesRemoval,
			UsageMonitoringInfos: usageMonitoring,
		},
	)
	assert.NoError(t, err)
	assert.Eventually(t, tr.WaitForPolicyReAuthToProcess(raa, imsi), time.Minute, 2*time.Second)

	assert.NotNil(t, raa)
	if raa != nil {
		assert.Equal(t, diam.Success, int(raa.ResultCode))
	}
	// Generate more traffic
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	// Check that enforcement flow is installed and traffic is less than the quota
	tr.AssertPolicyUsage(imsi, "pcrf-reauth-raa2", 0, 500*KiloBytes+Buffer)

	// trigger disconnection
	tr.DisconnectAndAssertSuccess(imsi)
	tr.AssertEventuallyAllRulesRemovedAfterDisconnect(imsi)
}

// - Install two static rules "static-pass-all-raa1" and "static-pass-all-raa2"
//   and a rule base "base-raa1"
// - Generate traffic and assert that there's > 0 data usage for the rule with the
//   highest priority.
// - Send a PCRF ReAuth request to refill quota for a session
// - Assert that the response is successful
// - Generate traffic and assert that there's > 0 data usage for the newly installed
//   rule.
// - Asserting that quota was updated is still needed.
func TestGxReAuthQuotaRefill(t *testing.T) {
	fmt.Println("\nRunning TestGxReAuthQuotaRefill...")

	tr, ruleManager, ue := gxReAuthTestSetup(t)
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()
	imsi := ue.GetImsi()

	tr.AuthenticateAndAssertSuccess(imsi)
	// Wait for flows to be installed
	assert.Eventually(t, tr.WaitForEnforcementStatsForRule(
		imsi, "static-pass-all-raa1", "static-pass-all-raa2"), time.Minute, 2*time.Second)

	// Generate over 80% of the quota to trigger a CCR Update
	req := &cwfprotos.GenTrafficRequest{
		Imsi:    imsi,
		Volume:  &wrappers.StringValue{Value: "500K"},
		Bitrate: &wrappers.StringValue{Value: "20M"},
		Timeout: 60,
	}
	_, err := tr.GenULTraffic(req)
	assert.NoError(t, err)
	assert.Eventually(t, tr.WaitForEnforcementStatsForRule(imsi, "static-pass-all-raa1"), time.Minute, 2*time.Second)

	// Check that enforcement flow is installed and traffic is less than the quota
	tr.AssertPolicyUsage(imsi, "static-pass-all-raa1", 0, 500*KiloBytes+Buffer)

	// Install a Pass-All Rule with higher priority using PolicyReAuth
	usageMonitoring := []*fegprotos.UsageMonitoringInformation{getUsageInformation("raakey1", 250*KiloBytes)}
	raa, err := sendPolicyReAuthRequest(
		&fegprotos.PolicyReAuthTarget{
			Imsi:                 imsi,
			UsageMonitoringInfos: usageMonitoring,
		},
	)
	assert.NoError(t, err)
	assert.Eventually(t, tr.WaitForPolicyReAuthToProcess(raa, imsi), time.Minute, 2*time.Second)

	assert.NotNil(t, raa)
	if raa != nil {
		assert.Equal(t, diam.Success, int(raa.ResultCode))
	}
	// Generate more traffic
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	// Check that enforcement flow is installed and traffic is less than the quota
	// Usage monitoring does not activate or deactivate services when the quota is up.
	// thus a method to check the current quota is needed to verify the success of this test.
	tr.AssertPolicyUsage(imsi, "static-pass-all-raa1", 500, 1*MegaBytes+Buffer)

	// trigger disconnection
	tr.DisconnectAndAssertSuccess(imsi)
	tr.AssertEventuallyAllRulesRemovedAfterDisconnect(imsi)
}
