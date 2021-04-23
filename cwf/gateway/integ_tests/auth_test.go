// +build all authenticate

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

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
)

func getCalledStationIDs() []string {
	return []string{"98-DE-D0-84-B5-47:CWF-TP-LINK_B547_5G",
		"78-FF-FF-84-B5-99:CWF-TP-LINK_B547_5G"}
}

// - Initialize 3 UEs and initiate Authentication. Assert that it is successful.
// - Disconnect all UEs.
func TestAuthenticateMultipleUEs(t *testing.T) {
	fmt.Println("\nRunning TestAuthenticate...")
	tr := NewTestRunner(t)
	ues, err := tr.ConfigUEs(3)
	assert.NoError(t, err)
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, tr.CleanUp())
	}()

	for _, ue := range ues {
		tr.AuthenticateAndAssertSuccess(ue.GetImsi())
		tr.DisconnectAndAssertSuccess(ue.GetImsi())
	}
	time.Sleep(1 * time.Second)
}

// - Initialize 2 UEs and initiate Authentication using diffrent APs.
//   Assert that it is successful.
// - Disconnect all UEs.
func TestAuthenticateWithDifferentAPs(t *testing.T) {
	fmt.Println("\nRunning TestAuthenticateWithDifferentAPs...")
	tr := NewTestRunner(t)
	ues, err := tr.ConfigUEs(2)
	assert.NoError(t, err)
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, tr.CleanUp())
	}()

	CalledStationIDs := getCalledStationIDs()
	for i, ue := range ues {
		tr.AuthenticateWithCalledIDAndAssertSuccess(ue.GetImsi(), CalledStationIDs[i])
		_, err = tr.Disconnect(ue.GetImsi(), CalledStationIDs[i])
		assert.NoError(t, err)
	}
	time.Sleep(1 * time.Second)
}

// - Expect a Gx CCR-I to come into PCRF, and return with Authentication Reject.
// - Configure a UE and trigger an authentication. Assert that the expectation was
//   met, and the authentication failed.
// - Expect a Gx CCR-I to come into PCRF, and return with Success. Expect a Gy
//   CCR-I to come into OCS, and return with Authentication Reject.
// - Trigger an authentication. Assert that all expectations were met, and the
//   authentication failed.
func TestAuthenticateFail(t *testing.T) {
	fmt.Println("\nRunning TestAuthenticateFail...")
	tr := NewTestRunner(t)

	assert.NoError(t, useOCSMockDriver())
	assert.NoError(t, usePCRFMockDriver())
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearOCSMockDriver())
		assert.NoError(t, clearPCRFMockDriver())
		assert.NoError(t, tr.CleanUp())
	}()

	ues, err := tr.ConfigUEs(2)
	assert.NoError(t, err)

	// ----- Gx CCR-I fail -> Authentication fails -----
	imsi := ues[0].GetImsi()
	gxInitReq := protos.NewGxCCRequest(imsi, protos.CCRequestType_INITIAL)
	gxInitAns := protos.NewGxCCAnswer(diam.AuthenticationRejected)
	gxInitExpectation := protos.NewGxCreditControlExpectation().Expect(gxInitReq).Return(gxInitAns)

	defaultGxAns := protos.NewGxCCAnswer(diam.AuthenticationRejected)
	assert.NoError(t, setPCRFExpectations([]*protos.GxCreditControlExpectation{gxInitExpectation}, defaultGxAns))

	tr.AuthenticateAndAssertFail(imsi)
	tr.AssertAllGxExpectationsMetNoError()

	// Since CCR/A-I failed, pipelined should see no rules installed
	tr.AssertPolicyEnforcementRecordIsNil(imsi)

	// ----- Gx CCR-I success && Gy CCR-I fail -> Authentication fails -----
	imsi = ues[1].GetImsi()
	gxInitReq = protos.NewGxCCRequest(imsi, protos.CCRequestType_INITIAL)
	gxInitAns = protos.NewGxCCAnswer(diam.Success).
		SetDynamicRuleInstall(getPassAllRuleDefinition("rule1", "", swag.Uint32(1), 0))
	gxInitExpectation = gxInitExpectation.Expect(gxInitReq).Return(gxInitAns)
	assert.NoError(t, setPCRFExpectations([]*protos.GxCreditControlExpectation{gxInitExpectation}, defaultGxAns))
	// Fail on Gy
	gyInitReq := protos.NewGyCCRequest(imsi, protos.CCRequestType_INITIAL)
	gyInitAns := protos.NewGyCCAnswer(diam.AuthenticationRejected)
	gyInitExpectation := protos.NewGyCreditControlExpectation().Expect(gyInitReq).Return(gyInitAns)
	defaultGyAns := gyInitAns
	assert.NoError(t, setOCSExpectations([]*protos.GyCreditControlExpectation{gyInitExpectation}, defaultGyAns))

	tr.AuthenticateAndAssertFail(imsi)
	// assert gx & gy init was received
	tr.AssertAllGxExpectationsMetNoError()
	tr.AssertAllGyExpectationsMetNoError()

	// Since CCR/A-I failed, pipelined should see no rules installed
	tr.AssertPolicyEnforcementRecordIsNil(imsi)
}

// - Set an expectation for a CCR-I to be sent up to PCRF, to which it will
//   respond with a rule install for a pass-all dynamic rule and 250KB of
//   quota.
//   Trigger a authentication and assert the CCR-I is received.
// - Generate traffic to put traffic through the newly installed rule.
func TestAuthenticateUplinkTraffic(t *testing.T) {
	fmt.Println("\nRunning TestAuthenticateUplinkTraffic...")
	tr := NewTestRunner(t)
	assert.NoError(t, usePCRFMockDriver())
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearPCRFMockDriver())
		assert.NoError(t, tr.CleanUp())
	}()

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)

	imsi := ues[0].GetImsi()
	usageMonitorInfo := getUsageInformation("mkey1", 2*MegaBytes)

	initRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGxCCAnswer(diam.Success).
		SetDynamicRuleInstall(getPassAllRuleDefinition("dynamic-pass-all", "mkey1", nil, 100)).
		SetUsageMonitorInfo(usageMonitorInfo)
	initExpectation := protos.NewGxCreditControlExpectation().Expect(initRequest).Return(initAnswer)
	// return success with credit on unexpected requests
	defaultAnswer := protos.NewGxCCAnswer(2001).SetUsageMonitorInfo(usageMonitorInfo)
	assert.NoError(t, setPCRFExpectations([]*protos.GxCreditControlExpectation{initExpectation}, defaultAnswer))

	tr.AuthenticateAndAssertSuccess(imsi)

	req := &cwfprotos.GenTrafficRequest{
		Imsi:    imsi,
		Volume:  &wrappers.StringValue{Value: "1M"},
		Bitrate: &wrappers.StringValue{Value: "60M"},
		Timeout: 30,
	}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)

	tr.AssertAllGxExpectationsMetNoError()

	tr.DisconnectAndAssertSuccess(imsi)
	tr.AssertEventuallyAllRulesRemovedAfterDisconnect(imsi)
}

// - Authenticate a UE through a first AP then switch to use a second AP
// - Set an expectation for a CCR-I to be sent up to PCRF, to which it will
//   respond with a rule install for a pass-all dynamic rule and 250KB of quota.
// - Trigger UE authentications through AP1 and generate traffic to  put it
//   through the newly installed rule.
// - Trigger UE authentications through AP2 and assert that
//   only one CCR-I is received. Sessiond must re-use the same session
//   during the handover.
// - Generate traffic to put traffic through the newly installed rule.
func TestAuthenticateMultipleAPsUplinkTraffic(t *testing.T) {
	fmt.Println("\nRunning TestAuthenticateMultipleAPsUplinkTraffic...")
	tr := NewTestRunner(t)
	assert.NoError(t, usePCRFMockDriver())
	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearPCRFMockDriver())
		assert.NoError(t, tr.CleanUp())
	}()

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)

	imsi := ues[0].GetImsi()
	usageMonitorInfo := getUsageInformation("mkey1", 2*MegaBytes)
	initRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGxCCAnswer(diam.Success).
		SetDynamicRuleInstall(getPassAllRuleDefinition("dynamic-pass-all", "mkey1", nil, 100)).
		SetUsageMonitorInfo(usageMonitorInfo)
	initExpectation := protos.NewGxCreditControlExpectation().Expect(initRequest).Return(initAnswer)
	// return success with credit on unexpected requests
	defaultAnswer := protos.NewGxCCAnswer(2001).SetUsageMonitorInfo(usageMonitorInfo)
	assert.NoError(t, setPCRFExpectations([]*protos.GxCreditControlExpectation{initExpectation}, defaultAnswer))

	CalledStationIDs := getCalledStationIDs()
	tr.AuthenticateWithCalledIDAndAssertSuccess(imsi, CalledStationIDs[0])

	req := &cwfprotos.GenTrafficRequest{
		Imsi:    imsi,
		Volume:  &wrappers.StringValue{Value: "1M"},
		Bitrate: &wrappers.StringValue{Value: "60M"},
		Timeout: 30,
	}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)

	tr.AuthenticateWithCalledIDAndAssertSuccess(imsi, CalledStationIDs[1])

	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)

	tr.AssertAllGxExpectationsMetNoError()

	_, err = tr.Disconnect(imsi, CalledStationIDs[1])
	assert.NoError(t, err)
	tr.AssertEventuallyAllRulesRemovedAfterDisconnect(imsi)
}
