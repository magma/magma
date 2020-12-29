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

	"magma/feg/cloud/go/protos"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/stretchr/testify/assert"
)

// - Set an expectation for a CCR-I to be sent up to PCRF, to which it will
//   respond with a rule install for a pass-all dynamic rule and 250KB of
//   quota.
//   Trigger a authentication and assert the CCR-I is received.
//   Wait for the Idle Timer to kick off and trigger a disconnect
//   The Idle Timer length is set in AAA_service's mconfig
func TestIdleTimer(t *testing.T) {
	fmt.Println("\nRunning TestIdleTimer...")
	tr := NewTestRunner(t)
	assert.NoError(t, usePCRFMockDriver())
	// Overwrite the mconfig to have a shorter idle timer length
	err := tr.OverwriteMConfig("./gateway.mconfig.idle_timer", "aaa_server")
	assert.NoError(t, err)
	assert.NoError(t, tr.RestartService("aaa_server"))
	// give it a second after the restart...
	time.Sleep(2 * time.Second)
	defer func() {
		err = tr.OverwriteMConfig("gateway.mconfig", "aaa_server")
		assert.NoError(t, err)
		assert.NoError(t, tr.RestartService("aaa_server"))
		// Clear hss, ocs, and pcrf
		assert.NoError(t, clearPCRFMockDriver())
		assert.NoError(t, tr.CleanUp())
	}()

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)

	imsi := ues[0].GetImsi()
	usageMonitorInfo := getUsageInformation("mkey1", 250*KiloBytes)

	initRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGxCCAnswer(diam.Success).
		SetDynamicRuleInstall(getPassAllRuleDefinition("dynamic-pass-all", "mkey1", nil, 100)).
		SetUsageMonitorInfo(usageMonitorInfo)
	initExpectation := protos.NewGxCreditControlExpectation().Expect(initRequest).Return(initAnswer)
	// return success with credit on unexpected requests
	defaultAnswer := protos.NewGxCCAnswer(2001).SetUsageMonitorInfo(usageMonitorInfo)
	assert.NoError(t, setPCRFExpectations([]*protos.GxCreditControlExpectation{initExpectation}, defaultAnswer))

	tr.AuthenticateAndAssertSuccessWithRetries(imsi, 5)

	tr.AssertAllGxExpectationsMetNoError()

	// If there is no activity for > Idle Timer, AAA service will send a
	// terminate request
	terminateRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_TERMINATION)
	terminateAnswer := protos.NewGxCCAnswer(diam.Success)
	terminateExpectation := protos.NewGxCreditControlExpectation().Expect(terminateRequest).Return(terminateAnswer)
	expectations := []*protos.GxCreditControlExpectation{terminateExpectation}
	assert.NoError(t, setPCRFExpectations(expectations, nil))

	// Wait for Session Idle Timer to kick off + Wait for Termination to go through
	// Idle Timer = 3 seconds
	fmt.Println("Waiting 3 seconds for the idle timer to kick in...")
	time.Sleep(3 * time.Second)
	tr.AssertEventuallyAllRulesRemovedAfterDisconnect(imsi)

	tr.AssertAllGxExpectationsMetNoError()
}
