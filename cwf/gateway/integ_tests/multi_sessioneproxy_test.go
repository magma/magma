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
	"math"
	"sync"
	"testing"
	"time"

	cwfprotos "magma/cwf/cloud/go/protos"
	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/multiplex"
	"magma/lte/cloud/go/services/policydb/obsidian/models"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
)

const (
	numInstances = 2
)

type multipleScenarioElement struct {
	pcrfName    string
	ocsName     string
	ruleManager *RuleManager
	IMSIs       []string
}

func generateMultipleScenarioAndAddSubscribers(t *testing.T, tr *TestRunner, numUEs int) []*multipleScenarioElement {
	// get the instance per IMSI based on the algorithm to distribuite IMSIs on FEG
	IMSIs := generateRandomIMSIS(numUEs, nil)
	IMSIsPerInstance := make([][]string, numInstances, numInstances)
	// create a multiplexor with the same value as main.go on session_proxy
	mux, err := multiplex.NewStaticMultiplexByIMSI(numInstances)
	assert.NoError(t, err)
	for _, imsi := range IMSIs {
		ctx := multiplex.NewContext().WithIMSI(imsi)
		i, err := mux.GetIndex(ctx)
		assert.NoError(t, err)
		IMSIsPerInstance[i] = append(IMSIsPerInstance[i], imsi)
	}
	// Add IMSIs to each instance
	scenario := make([]*multipleScenarioElement, 0, numInstances)
	for i, IMSIs := range IMSIsPerInstance {
		// get names of the servers for this specific instance
		pcrfName, ocsName, err := getPCRFandOCSnamePerInstance(i)
		assert.NoError(t, err)

		// Create accounts in HSS, PCRF and OCS for each IMSI
		ues, err := tr.ConfigUEsPerInstance(IMSIs, pcrfName, ocsName)
		assert.NoError(t, err)
		assert.Equal(t, len(ues), len(IMSIs))
		assert.True(t, checkIMSIsListsAreEqual(ues, IMSIs))

		// Create dynamic rules on AGW database
		ruleManager, err := NewRuleManagerPerInstance(pcrfName)
		assert.NoError(t, err)
		ratingGroup := uint32(1)
		err = ruleManager.AddStaticPassAllToDB("static-pass-all", "mkey1",
			ratingGroup, models.PolicyRuleTrackingTypeOCSANDPCRF, 10)
		assert.NoError(t, err)

		tr.WaitForPoliciesToSync()

		// main OCS config
		assert.NoError(t, setNewOCSConfigPerInstance(
			ocsName,
			&protos.OCSConfig{
				MaxUsageOctets: &protos.Octets{TotalOctets: 10 * MegaBytes},
			},
		),
		)

		// add static rules to PCRF and credit to OCS
		for _, imsi := range IMSIs {
			// Add the dynamic ruleID to the PCRF (only name definition is on AGW, since this is dynamic)
			err = ruleManager.AddRulesToPCRF(imsi, []string{"static-pass-all"}, nil)
			assert.NoError(t, err)

			// Credit addition to OCS with same chargingKey as the rule ratingGroup
			assert.NoError(t, setCreditOnOCSPerInstance(
				ocsName,
				&protos.CreditInfo{
					Imsi:        imsi,
					ChargingKey: 1,
					Volume:      &protos.Octets{TotalOctets: 10 * MegaBytes},
					UnitType:    protos.CreditInfo_Bytes},
			),
			)
			// TODO: Verify OCSs reports proper credit
			//infos, err := getCreditOnOCSPerInstance(ocsName, imsi)
			//fmt.Printf("\t ---> credit left: %s --- %+v\n", err, infos)
		}
		scenario = append(scenario,
			&multipleScenarioElement{
				pcrfName,
				ocsName,
				ruleManager,
				IMSIs,
			})
	}
	return scenario
}

func checkIMSIsListsAreEqual(UEs []*cwfprotos.UEConfig, IMSIs []string) bool {
	set := make(map[string]bool)
	for _, imsi := range IMSIs {
		set[imsi] = true
	}
	for _, ue := range UEs {
		_, present := set[ue.GetImsi()]
		if !present {
			return false
		}
	}
	return true
}

func getPCRFandOCSnamePerInstance(instanceId int) (pcrfName string, ocsName string, err error) {
	switch instanceId {
	case 0:
		pcrfName, ocsName = MockPCRFRemote, MockOCSRemote
	case 1:
		pcrfName, ocsName = MockPCRFRemote2, MockOCSRemote2
	default:
		err = fmt.Errorf("Instance number %d not valid", instanceId)
	}
	return
}

// TestMultiSessionProxyMonitorAndUsageReportEnforcement is an experimental
// test to try multiple OCS and PCRF servers. Currenty it only supports 1 UE
// - Create one UE and add monitoring key and credit
// - Attach UE, tranfer data, detach
// - Check that the Monitored data by the PCRF instance is good
// TODO:
//  * Support for multiple UEs (depends on UEsim service)
//  * Check OCS credit has been reported (right now sessiond sends CCR after accounts
//    are deleted from OCS and PCRF)
func TestMultiSessionProxyMonitorAndUsageReportEnforcement(t *testing.T) {
	fmt.Println("\nRunning TestMultiSessionProxyMonitorAndUsageReportEnforcement...")
	tr := NewTestRunnerWithTwoPCRFandOCS(t)

	// Overwrite the mconfig to enable multiproxy with OCS2 and PCRF2
	err := tr.OverwriteMConfig("./gateway.mconfig.multi_session_proxy", "session_proxy")
	assert.NoError(t, err)
	assert.NoError(t, tr.RestartService("ocs"))
	assert.NoError(t, tr.RestartService("pcrf"))
	assert.NoError(t, tr.RestartService("ocs2"))
	assert.NoError(t, tr.RestartService("pcrf2"))
	assert.NoError(t, tr.RestartService("ingress"))
	time.Sleep(3 * time.Second)
	assert.NoError(t, tr.RestartService("session_proxy"))

	// TODO: this only works with 1 user because UEsim can only use one single MAC address
	scenario := generateMultipleScenarioAndAddSubscribers(t, tr, 1)
	defer func() {
		// Clear hss, ocs, and pcrf
		for _, scenarioElmnt := range scenario {
			assert.NoError(t, scenarioElmnt.ruleManager.RemoveInstalledRules())
		}
		assert.NoError(t, tr.CleanUp())

		err = tr.OverwriteMConfig("gateway.mconfig", "session_proxy")
		assert.NoError(t, err)
		assert.NoError(t, tr.RestartService("ocs"))
		assert.NoError(t, tr.RestartService("pcrf"))
		assert.NoError(t, tr.RestartService("ingress"))
		time.Sleep(3 * time.Second)
		assert.NoError(t, tr.RestartService("session_proxy"))
		// ignore ocs2 and pcrf2 errors
		tr.RestartService("ocs2")
		tr.RestartService("pcrf2")
	}()

	tr.WaitForPoliciesToSync()

	var wg sync.WaitGroup
	for _, element := range scenario {
		for _, imsiPTR := range element.IMSIs {
			wg.Add(1)

			imsi := imsiPTR
			go func() {
				defer wg.Done()
				tr.AuthenticateAndAssertSuccess(imsi)
				// this wait can be remove
				tr.WaitForEnforcementStatsToSync()
				req := &cwfprotos.GenTrafficRequest{
					Imsi:    imsi,
					Volume:  &wrappers.StringValue{Value: "2M"},
					Bitrate: &wrappers.StringValue{Value: "60M"},
					Timeout: 30,
				}
				_, err := tr.GenULTraffic(req)
				assert.NoError(t, err)
			}()
		}
	}
	wg.Wait()
	// this wait CAN NOT be removed. Extra wait to make sure sessiond reported all traffic.
	tr.WaitForEnforcementStatsToSync()
	tr.WaitForEnforcementStatsToSync()

	for _, element := range scenario {
		for _, imsi := range element.IMSIs {
			recordsBySubID, err := tr.GetPolicyUsage()
			assert.NoError(t, err)

			// Check pipelined let the UE to send traffic
			record := recordsBySubID["IMSI"+imsi]["static-pass-all"]
			assert.NotNil(t, record, fmt.Sprintf("No policy usage record for imsi: %v", imsi))
			if record != nil {
				// We should not be seeing > 1024k data here
				assert.True(t, record.BytesTx > uint64(0), fmt.Sprintf("%s did not pass any data", record.RuleId))
				assert.NoError(t, err)
				assert.True(t, record.BytesTx <= uint64(math.Round(2.5*MegaBytes)+Buffer), fmt.Sprintf("policy usage: %v", record))
				// TODO: make sure OCS records its proper usage and it matches with what we monitored
				//infos, err := getCreditOnOCSPerInstance(element.ocsName, imsi)
				//fmt.Printf("\t ---> policy usage: %v\n", record)
				//fmt.Printf("\t ---> credit left: %+v\n", infos)
			}
			// Detach this UE
			tr.DisconnectAndAssertSuccess(imsi)

			// TODO: check CCR-T is sent to the right instance
		}
	}
	// Wait for CCR-T to propagate up
	time.Sleep(2 * time.Second)
}
