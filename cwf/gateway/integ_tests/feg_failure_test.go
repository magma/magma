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

	cwfprotos "magma/cwf/cloud/go/protos"
	"magma/feg/cloud/go/protos"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
)

// - Shutdown the ingress container to break the link between the CWAG and FEG
// - Asset that authentication fails and that no rules were insalled
func TestLinkFailureCWAGtoFEG(t *testing.T) {
	fmt.Println("\nRunning TestLinkFailureCWAGtoFEG...")

	tr := NewTestRunner(t)
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)
	assert.NoError(t, tr.PauseService("ingress"))

	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearOCSMockDriver())
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
		assert.NoError(t, tr.RestartService("ingress"))
	}()

	ues, err := tr.ConfigUEs(1)
	ue := ues[0]
	imsi := ue.Imsi
	assert.NoError(t, err)
	setNewOCSConfig(
		&protos.OCSConfig{
			MaxUsageOctets: &protos.Octets{TotalOctets: GyMaxUsageBytes},
			MaxUsageTime:   GyMaxUsageTime,
			ValidityTime:   GyValidityTime,
			UseMockDriver:  true,
		},
	)

	tr.WaitForPoliciesToSync()
	tr.AuthenticateAndAssertFail(imsi)

	// Since CCA-I was never received, there should be no rules installed
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	assert.Empty(t, recordsBySubID[prependIMSIPrefix(imsi)])
}

// FeG connection breaks mid way through a session
func TestLinkFailureCWAGtoFEGMidSession(t *testing.T) {
	fmt.Println("\nRunning TestLinkFailureCWAGtoFEGMidSession...")

	tr := NewTestRunner(t)
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)

	assert.NoError(t, usePCRFMockDriver())

	defer func() {
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearPCRFMockDriver())
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
		assert.NoError(t, tr.RestartService("ingress"))
	}()
	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)
	ue := ues[0]
	imsi := ue.GetImsi()
	usageMonitorInfo := getUsageInformation("mkey1", 250*KiloBytes)

	initRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGxCCAnswer(diam.Success).
		SetDynamicRuleInstall(getPassAllRuleDefinition("dynamic-pass-all", "mkey1", nil, 100)).
		SetUsageMonitorInfo(usageMonitorInfo)
	initExpectation := protos.NewGxCreditControlExpectation().Expect(initRequest).Return(initAnswer)
	// return success with credit on unexpected requests
	defaultAnswer := protos.NewGxCCAnswer(2001).SetUsageMonitorInfo(usageMonitorInfo)
	assert.NoError(t, setPCRFExpectations([]*protos.GxCreditControlExpectation{initExpectation}, defaultAnswer))

	tr.AuthenticateAndAssertSuccess(imsi)

	req := &cwfprotos.GenTrafficRequest{Imsi: imsi, Volume: &wrappers.StringValue{Value: "100K"}}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	assert.NoError(t, tr.PauseService("ingress"))

	req = &cwfprotos.GenTrafficRequest{Imsi: imsi, Volume: &wrappers.StringValue{Value: "500K"}}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()
	tr.DisconnectAndAssertSuccess(ue.GetImsi())
}
