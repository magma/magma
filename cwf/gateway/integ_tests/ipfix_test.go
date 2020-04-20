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
	"io/ioutil"
	cwfprotos "magma/cwf/cloud/go/protos"
	"magma/feg/cloud/go/protos"
	"magma/lte/cloud/go/services/policydb/obsidian/models"
	"strings"
	"testing"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
)

// - Set an expectation for a  CCR-I to be sent up to PCRF, to which it will
//   respond with a rule install (usage-enforcement-static-pass-all), 250KB of
//   quota.
//   Generate traffic and assert the CCR-I is received.
// - Set an expectation for a CCR-U with >80% of data usage to be sent up to
// 	 PCRF, to which it will response with more quota.
//   Generate traffic and assert the CCR-U is received.
// - Generate traffic to put traffic through the newly installed rule.
//   Assert that there's > 0 data usage in the rule.
// - Expect a CCR-T, trigger a UE disconnect, and assert the CCR-T is received.
// - Assert that IPDR records were properly exported
func TestIpfixEnforcement(t *testing.T) {
	t.Skip("Skipping test due to DPI changes")
	fmt.Println("\nRunning IPFIX TEST...")
	tr := NewTestRunner(t)

	// Enable IPFIX exporting
	enableIpfixExport()
	err := tr.RestartService("pipelined")
	if err != nil {
		fmt.Printf("error restarting pipelined %v", err)
		assert.Fail(t, "failed restarting pipelined")
	}

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

	usageMonitorInfo := getUsageInformation("mkey1", 250*KiloBytes)

	initRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_INITIAL)
	initAnswer := protos.NewGxCCAnswer(diam.Success).
		SetStaticRuleInstalls([]string{"usage-enforcement-static-pass-all"}, []string{}).
		SetUsageMonitorInfo(usageMonitorInfo)
	initExpectation := protos.NewGxCreditControlExpectation().Expect(initRequest).Return(initAnswer)

	// We expect an update request with some usage update (probably around 80-100% of the given quota)
	updateRequest1 := protos.NewGxCCRequest(imsi, protos.CCRequestType_UPDATE).
		SetUsageMonitorReport(usageMonitorInfo).
		SetUsageReportDelta(250 * KiloBytes * 0.2)
	updateAnswer1 := protos.NewGxCCAnswer(diam.Success).SetUsageMonitorInfo(usageMonitorInfo)
	updateExpectation1 := protos.NewGxCreditControlExpectation().Expect(updateRequest1).Return(updateAnswer1)
	expectations := []*protos.GxCreditControlExpectation{initExpectation, updateExpectation1}
	// On unexpected requests, just return the default update answer
	assert.NoError(t, setPCRFExpectations(expectations, updateAnswer1))

	tr.AuthenticateAndAssertSuccessWithRetries(imsi, 5)

	req := &cwfprotos.GenTrafficRequest{Imsi: imsi, Volume: &wrappers.StringValue{Value: *swag.String("500K")}}
	_, err = tr.GenULTraffic(req)
	assert.NoError(t, err)
	tr.WaitForEnforcementStatsToSync()

	// When we initiate a UE disconnect, we expect a terminate request to go up
	terminateRequest := protos.NewGxCCRequest(imsi, protos.CCRequestType_TERMINATION)
	terminateAnswer := protos.NewGxCCAnswer(diam.Success)
	terminateExpectation := protos.NewGxCreditControlExpectation().Expect(terminateRequest).Return(terminateAnswer)
	expectations = []*protos.GxCreditControlExpectation{terminateExpectation}
	assert.NoError(t, setPCRFExpectations(expectations, nil))

	tr.DisconnectAndAssertSuccess(imsi)
	tr.WaitForEnforcementStatsToSync()

	// Wait for CCR-T to propagate up
	time.Sleep(3 * time.Second)

	// Collector saves files in 'YearMonthDate' format
	//year, month, day := time.Now().Date()
	ipfixDir := fmt.Sprintf("/home/vagrant/magma/records/%s/1/", IPDRControllerIP)

	// Check if ipdr export file was created
	files, err := ioutil.ReadDir(ipfixDir)
	assert.NoError(t, err)
	assert.NotEqual(t, len(files), 0)

	ipfixFileDir := fmt.Sprintf("%s/%s", ipfixDir, files[len(files)-1].Name())
	file, err := ioutil.ReadFile(ipfixFileDir)
	assert.NoError(t, err)

	ipdrExport := string(file)
	imsiStr, err := getEncodedIMSI(imsi)
	assert.NoError(t, err)

	imsiIPFIXRecord := fmt.Sprintf("%s, %s, %s, %s", imsiStr, ipfixMSISDN,
		ipfixApnMacAddress, ipfixApnName)

	// Check if export contains the data from this test
	assert.True(t, strings.Contains(ipdrExport, imsiIPFIXRecord))

	// Disable IPFIX exporting
	disableIpfixExport()
	err = tr.RestartService("pipelined")
	if err != nil {
		fmt.Printf("error restarting pipelined %v", err)
		assert.Fail(t, "failed restarting pipelined")
	}
}

func enableIpfixExport() {
	replacePipelinedConfigValue("#'ipfix'", "'ipfix'")
}

func disableIpfixExport() {
	replacePipelinedConfigValue("'ipfix'", "#'ipfix'")
}

func replacePipelinedConfigValue(old string, new string) {
	path := "/home/vagrant/magma/cwf/gateway/integ_tests/pipelined.yml"
	read, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	newContents := strings.Replace(string(read), old, new, -1)
	err = ioutil.WriteFile(path, []byte(newContents), 0)
	if err != nil {
		panic(err)
	}
}
