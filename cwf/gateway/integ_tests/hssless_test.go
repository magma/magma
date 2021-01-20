// +build hssless

/*Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package integration

import (
	"fmt"
	"testing"
	"time"

	cwfprotos "magma/cwf/cloud/go/protos"
	"magma/cwf/gateway/registry"
	fegprotos "magma/feg/cloud/go/protos"
	"magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/orc8r/cloud/go/test_utils"

	"magma/cwf/gateway/services/uesim/servicers"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

const (
	DefaultApn     = "test"
	DefaultMsisdn  = "5100001234"
	DefaultRatType = 1
)

func TestHsslessAuthenticateUe(t *testing.T) {
	fmt.Println("\nRunning TestAuthenticateUe HSSLess...")

	ueSimServer, err := setupHssLessTestEnv(t)
	assert.NoError(t, err)

	tr := NewTestRunner(t)
	ruleManager, err := NewRuleManager()
	assert.NoError(t, err)
	defer func() {
		// Delete omni rules
		assert.NoError(t, ruleManager.RemoveOmniPresentRulesFromDB("omni"))
		// Clear hss, ocs, and pcrf
		assert.NoError(t, ruleManager.RemoveInstalledRules())
		assert.NoError(t, tr.CleanUp())
	}()

	ues, err := tr.ConfigUEs(1)
	assert.NoError(t, err)

	imsi := ues[0].GetImsi()

	// Configure static rule in policyDB and PCRF
	err = ruleManager.AddStaticPassAllToDBAndPCRFforIMSIs([]string{imsi}, "omni-pass-all-1", "1", 1, models.PolicyRuleConfigTrackingTypeONLYOCS, 20)
	assert.NoError(t, err)

	// Set Credits on OCS
	setCreditOnOCS(
		&fegprotos.CreditInfo{
			Imsi:        imsi,
			ChargingKey: 1,
			Volume:      &fegprotos.Octets{TotalOctets: 500 * KiloBytes},
			UnitType:    fegprotos.CreditInfo_Bytes,
		},
	)

	tr.WaitForPoliciesToSync()

	ues[0].HsslessCfg = &cwfprotos.AuthenticateRequestHssLess{
		Msisdn: DefaultMsisdn,
		Rat:    DefaultRatType,
		Apn:    DefaultApn,
	}

	_, err = ueSimServer.AddUE(context.Background(), ues[0])
	assert.NoError(t, err)

	authReq := &cwfprotos.AuthenticateRequest{Imsi: imsi}

	res, err := ueSimServer.Authenticate(context.Background(), authReq)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	time.Sleep(5 * time.Second)

	discReq := &cwfprotos.DisconnectRequest{Imsi: imsi}
	discRes, err := ueSimServer.Disconnect(context.Background(), discReq)
	assert.NoError(t, err)
	assert.NotNil(t, discRes)

	time.Sleep(5 * time.Second)

}

func setupHssLessTestEnv(t *testing.T) (*servicers.UESimServerHssLess, error) {

	registry.AddService("SESSIOND", "127.0.0.1", 50065)

	store := test_utils.NewSQLBlobstore(t, "hssless_test_blobstore")
	server, err := servicers.NewUESimServerHssLess(store)
	return server, err
}
