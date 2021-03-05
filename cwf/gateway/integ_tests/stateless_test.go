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

	cwfProtos "magma/cwf/cloud/go/protos"
	"magma/feg/cloud/go/protos"
	lteProtos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/policydb/obsidian/models"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
)

// - Set an expectation for a  CCR-I to be sent up to PCRF/OCS, to which it will
//   respond with a rule install (usage-enforcement-static-pass-all), 250KB of
//   quota (both Gx and Gy)
// - Assert that the authentication went through && CCR-I is received
// - Restart SessionD
// - Set an expectation for a CCR-U with >80% of data usage to be sent up to
// 	 PCRF/OCS, to which it will response with more quota
// - Generate traffic
// - Assert CCR-U is received
// - Restart SessionD
// - Assert that there's > 0 data usage in the rule
// - Expect a CCR-T, trigger a UE disconnect, and assert the CCR-T is received
func TestBasicEnforcementWithSessionDRestarts(t *testing.T) {
	t.Skip()
	fmt.Println("\nRunning TestBasicEnforcementWithSessionDRestarts...")
	tr := NewTestRunner(t)
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)
	assert.NoError(t, usePCRFMockDriver())
	assert.NoError(t, useOCSMockDriver())
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearOCSMockDriver())
		assert.NoError(t, clearPCRFMockDriver())
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)
	imsi := ues[0].GetImsi()

	err = ruleManager.AddStaticPassAllToDB("usage-enforcement-static-pass-all", "mkey1", 32, models.PolicyRuleTrackingTypeOCSANDPCRF, 3)
	assert.NoError(t, err)
	tr.WaitForPoliciesToSync()

	usageMonitorInfo := getUsageInformation("mkey1", 250*KiloBytes)
	quotaGrant := &protos.QuotaGrant{
		RatingGroup:        32,
		GrantedServiceUnit: &protos.Octets{TotalOctets: 250 * KiloBytes},
		IsFinalCredit:      false,
		ResultCode:         2001,
	}

	setInitExpectations(t, imsi, usageMonitorInfo, quotaGrant)
	tr.AuthenticateAndAssertSuccess(imsi)
	// Assert that a CCR-I was sent up to the PCRF/OCS
	tr.AssertAllGxExpectationsMetNoError()
	tr.AssertAllGyExpectationsMetNoError()

	// After authentication restart SessionD
	assert.NoError(t, tr.RestartService("sessiond"))
	fmt.Println("Waiting for SessionD to restore context...")
	time.Sleep(4 * time.Second)
	// All session state/config should be restored from Redis by now...

	// We expect an update request with some usage update (probably around 80-100% of the given quota)
	setUpdateExpectations(t, imsi, usageMonitorInfo, quotaGrant)
	req := &cwfProtos.GenTrafficRequest{Imsi: imsi, Volume: &wrappers.StringValue{Value: *swag.String("500K")}}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	// Assert that enforcement_stats rules are properly installed and the right
	// amount of data was passed through
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	record := recordsBySubID["IMSI"+imsi]["usage-enforcement-static-pass-all"]
	assert.NotNil(t, record, fmt.Sprintf("No policy usage record for imsi: %v", imsi))
	if record != nil {
		// We should not be seeing > 1024k data here
		assert.True(t, record.BytesTx > uint64(0), fmt.Sprintf("%s did not pass any data", record.RuleId))
		assert.True(t, record.BytesTx <= uint64(500*KiloBytes+Buffer), fmt.Sprintf("policy usage: %v", record))
	}

	// Assert that CCR-Us were sent up to the PCRF/OCS
	tr.AssertAllGxExpectationsMetNoError()
	tr.AssertAllGyExpectationsMetNoError()

	// After Update request restart SessionD
	assert.NoError(t, tr.RestartService("sessiond"))
	fmt.Println("Waiting for SessionD to restore context...")
	time.Sleep(1 * time.Second)
	// All session state/config should be restored from Redis by now...

	// When we initiate a UE disconnect, we expect a terminate request to go up
	setTerminateExpectations(t, imsi)

	tr.DisconnectAndAssertSuccess(imsi)
	tr.WaitForEnforcementStatsToSync()

	// Wait for CCR-T to propagate up
	time.Sleep(2 * time.Second)

	// Assert that we saw Gx/Gy Terminate requests
	tr.AssertAllGxExpectationsMetNoError()
	tr.AssertAllGyExpectationsMetNoError()
}

func setInitExpectations(t *testing.T, imsi string, usageMonitor *protos.UsageMonitoringInformation,
	quotaGrant *protos.QuotaGrant) {
	gxRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_INITIAL)
	gxAnswer := protos.NewGxCCAnswer(diam.Success).
		SetStaticRuleInstalls([]string{"usage-enforcement-static-pass-all"}, []string{}).
		SetUsageMonitorInfo(usageMonitor)
	gxExpectation := protos.NewGxCreditControlExpectation().Expect(gxRequest).Return(gxAnswer)

	gyRequest := protos.NewGyCCRequest(imsi, protos.CCRequestType_INITIAL)
	gyAnswer := protos.NewGyCCAnswer(diam.Success).SetQuotaGrant(quotaGrant)
	gyExpectation := protos.NewGyCreditControlExpectation().Expect(gyRequest).Return(gyAnswer)

	assert.NoError(t, setPCRFExpectations([]*protos.GxCreditControlExpectation{gxExpectation}, nil))
	assert.NoError(t, setOCSExpectations([]*protos.GyCreditControlExpectation{gyExpectation}, nil))
}

func setUpdateExpectations(t *testing.T, imsi string, usageMonitor *protos.UsageMonitoringInformation,
	quotaGrant *protos.QuotaGrant) {
	gxRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_UPDATE).
		SetUsageMonitorReport(usageMonitor).
		SetUsageReportDelta(250 * KiloBytes * 0.2).
		SetEventTrigger(int32(lteProtos.EventTrigger_USAGE_REPORT))
	gxAnswer := protos.NewGxCCAnswer(diam.Success).SetUsageMonitorInfo(usageMonitor)
	gxExpectation := protos.NewGxCreditControlExpectation().Expect(gxRequest).Return(gxAnswer)

	gyRequest := protos.NewGyCCRequest(imsi, protos.CCRequestType_UPDATE)
	gyAnswer := protos.NewGyCCAnswer(diam.Success).SetQuotaGrant(quotaGrant)
	gyExpectation := protos.NewGyCreditControlExpectation().Expect(gyRequest).Return(gyAnswer)

	assert.NoError(t, setPCRFExpectations([]*protos.GxCreditControlExpectation{gxExpectation}, gxAnswer))
	assert.NoError(t, setOCSExpectations([]*protos.GyCreditControlExpectation{gyExpectation}, gyAnswer))
}

func setTerminateExpectations(t *testing.T, imsi string) {
	gxRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_TERMINATION)
	gxAnswer := protos.NewGxCCAnswer(diam.Success)
	gxExpectation := protos.NewGxCreditControlExpectation().Expect(gxRequest).Return(gxAnswer)

	gyRequest := protos.NewGyCCRequest(imsi, protos.CCRequestType_TERMINATION)
	gyAnswer := protos.NewGyCCAnswer(diam.Success)
	gyExpectation := protos.NewGyCreditControlExpectation().Expect(gyRequest).Return(gyAnswer)

	assert.NoError(t, setPCRFExpectations([]*protos.GxCreditControlExpectation{gxExpectation}, gxAnswer))
	assert.NoError(t, setOCSExpectations([]*protos.GyCreditControlExpectation{gyExpectation}, gyAnswer))
}
