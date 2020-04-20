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

	fegprotos "magma/feg/cloud/go/protos"

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
	tr.WaitForPoliciesToSync()
	tr.AuthenticateAndAssertFail(ue.Imsi)

	// Since CCA-I was never received, there should be no rules installed
	recordsBySubID, err := tr.GetPolicyUsage()
	assert.NoError(t, err)
	assert.Empty(t, recordsBySubID["IMSI"+ue.Imsi])
}
