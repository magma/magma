// +build all

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
	"magma/feg/cloud/go/protos"
	fegprotos "magma/feg/cloud/go/protos"
	"magma/lte/cloud/go/services/policydb/obsidian/models"

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
// - Set an expectation for a CCR-I to be sent up to OCS, wo which it will responde
// 	 with some credit installed
// - Generate traffic and assert the CCR-I is received.
// - Assert that the traffic goes through. This means the network wide rules
//   gets installed properly.
// - Trigger a Gx RAR with a rule removal for the block all rule. Assert the
//   answer is successful.
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
		assert.NoError(t, clearOCSMockDriver())
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()

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
	imsi := ues[0].GetImsi()

	// Set a block all rule to be installed by the PCRF
	err = ruleManager.AddStaticRuleToDB(getStaticDenyAll("static-block-all", "mkey1", 0, models.PolicyRuleConfigTrackingTypeONLYPCRF, 30))
	assert.NoError(t, err)
	// Override with an omni pass all static rule with a higher priority and with a charging key
	err = ruleManager.AddStaticPassAllToDB("omni-pass-all-1", "", 1, models.PolicyRuleConfigTrackingTypeONLYOCS, 20)
	assert.NoError(t, err)
	// Apply a network wide rule that points to the static rule above
	err = ruleManager.AddOmniPresentRulesToDB("omni", []string{"omni-pass-all-1"}, []string{""})
	assert.NoError(t, err)
	tr.WaitForPoliciesToSync()

	// Gx - PCRF config
	usageMonitorInfo := getUsageInformation("mkey1", 5*MegaBytes)
	gxInitRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_INITIAL)
	gxInitAnswer := protos.NewGxCCAnswer(diam.Success).
		SetStaticRuleInstalls([]string{"static-block-all"}, []string{}).
		SetUsageMonitorInfo(usageMonitorInfo)
	gxInitExpectation := protos.NewGxCreditControlExpectation().Expect(gxInitRequest).Return(gxInitAnswer)
	gxExpectations := []*protos.GxCreditControlExpectation{gxInitExpectation}
	assert.NoError(t, setPCRFExpectations(gxExpectations, nil)) // we don't expect any update requests

	// Gy - OCS config
	quotaGrant := &fegprotos.QuotaGrant{
		RatingGroup: 1,
		GrantedServiceUnit: &fegprotos.Octets{
			TotalOctets: 1 * MegaBytes,
		},
		IsFinalCredit: false,
		ResultCode:    diam.Success,
	}
	gyInitRequest := protos.NewGyCCRequest(imsi, protos.CCRequestType_INITIAL)
	gyInitAnswer := protos.NewGyCCAnswer(diam.Success).SetQuotaGrant(quotaGrant)
	gyInitExpectation := protos.NewGyCreditControlExpectation().Expect(gyInitRequest).Return(gyInitAnswer)
	gyExpectations := []*protos.GyCreditControlExpectation{gyInitExpectation}
	assert.NoError(t, setOCSExpectations(gyExpectations, nil))

	tr.AuthenticateAndAssertSuccess(imsi)

	req := &cwfprotos.GenTrafficRequest{Imsi: imsi, Volume: &wrappers.StringValue{Value: "200k"}}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	omniRecord := recordsBySubID["IMSI"+imsi]["omni-pass-all-1"]
	blockAllRecord := recordsBySubID["IMSI"+imsi]["static-block-all"]
	assert.NotNil(t, omniRecord, fmt.Sprintf("No policy usage omniRecord for imsi: %v", imsi))
	assert.NotNil(t, blockAllRecord, fmt.Sprintf("Block all record was not installed for imsi %v", imsi))

	if omniRecord != nil {
		assert.True(t, omniRecord.BytesTx > uint64(0), fmt.Sprintf("%s did not pass any data", omniRecord.RuleId))
	}
	if blockAllRecord != nil {
		assert.Equal(t, uint64(0x0), blockAllRecord.BytesTx)
	}

	tr.AssertAllGyExpectationsMetNoError()
	tr.AssertAllGxExpectationsMetNoError()

	// Trigger a ReAuth with rule removals of monitored rules
	target := &protos.PolicyReAuthTarget{
		Imsi: imsi,
		RulesToRemove: &protos.RuleRemovals{
			RuleNames: []string{"static-block-all"},
		},
	}
	fmt.Printf("Sending a ReAuthRequest with target %v\n", target)
	raa, err := sendPolicyReAuthRequest(target)
	assert.Eventually(t, tr.WaitForPolicyReAuthToProcess(raa, imsi), time.Minute, 2*time.Second)

	// Check ReAuth success
	assert.NotNil(t, raa)
	if raa != nil {
		assert.Equal(t, diam.Success, int(raa.ResultCode))
	}
	// trigger disconnection
	tr.DisconnectAndAssertSuccess(imsi)
	tr.AssertEventuallyAllRulesRemovedAfterDisconnect(imsi)
}

// TODO: test disabled for now. Need to modify mconfig to enable/disable Gx
func TestGxDisabledOmnipresentRules(t *testing.T) {
	t.Skip()
	fmt.Println("\nRunning TestOmnipresentRulesGxDisabled...")
	tr := NewTestRunner(t)
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)
	assert.NoError(t, usePCRFMockDriver())
	defer func() {
		// Delete omni rules
		assert.NoError(t, ruleManager.RemoveOmniPresentRulesFromDB("omni"))
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearPCRFMockDriver())
		assert.NoError(t, clearOCSMockDriver())
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()

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
	imsi := ues[0].GetImsi()

	err = ruleManager.AddStaticPassAllToDB("omni-pass-all-1", "", 1, models.PolicyRuleConfigTrackingTypeONLYOCS, 20)
	assert.NoError(t, err)
	// Apply a network wide rule that points to the static rule above
	err = ruleManager.AddOmniPresentRulesToDB("omni", []string{"omni-pass-all-1"}, []string{""})
	assert.NoError(t, err)
	tr.WaitForPoliciesToSync()

	// Gx - PCRF config
	assert.NoError(t, setPCRFExpectations(nil, nil)) // we don't expect any update requests

	// Gy - OCS config
	quotaGrant := &fegprotos.QuotaGrant{
		RatingGroup: 1,
		GrantedServiceUnit: &fegprotos.Octets{
			TotalOctets: 1 * MegaBytes,
		},
		IsFinalCredit: false,
		ResultCode:    diam.Success,
	}
	gyInitRequest := protos.NewGyCCRequest(imsi, protos.CCRequestType_INITIAL)
	gyInitAnswer := protos.NewGyCCAnswer(diam.Success).SetQuotaGrant(quotaGrant)
	gyInitExpectation := protos.NewGyCreditControlExpectation().Expect(gyInitRequest).Return(gyInitAnswer)
	gyExpectations := []*protos.GyCreditControlExpectation{gyInitExpectation}
	assert.NoError(t, setOCSExpectations(gyExpectations, nil))

	tr.AuthenticateAndAssertSuccess(imsi)

	req := &cwfprotos.GenTrafficRequest{Imsi: imsi, Volume: &wrappers.StringValue{Value: *swag.String("200K")}}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	omniRecord := recordsBySubID["IMSI"+imsi]["omni-pass-all-1"]
	assert.NotNil(t, omniRecord, fmt.Sprintf("No policy usage omniRecord for imsi: %v", imsi))
	assert.True(t, omniRecord.BytesTx > uint64(0), fmt.Sprintf("%s did not pass any data", omniRecord.RuleId))

	tr.AssertAllGyExpectationsMetNoError()
	tr.AssertAllGxExpectationsMetNoError()

	// trigger disconnection
	tr.DisconnectAndAssertSuccess(imsi)
	tr.AssertEventuallyAllRulesRemovedAfterDisconnect(imsi)
}
