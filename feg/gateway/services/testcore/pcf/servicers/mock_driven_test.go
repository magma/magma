/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package servicers_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/sbi"
	"magma/feg/gateway/services/testcore/pcf/servicers"
	orcprotos "magma/orc8r/lib/go/protos"
)

const (
	LOCAL_ADDR      = "127.0.0.1:0"
	BASE_PATH       = "/sm-policy-control/v1"
	HTTP_HOST       = "http://localhost:0"
	API_ROOT        = HTTP_HOST + BASE_PATH
	TOKEN_URL       = HTTP_HOST + "/token"
	CLIENT_ID       = "feg_magma_client"
	CLIENT_SECRET   = "feg_magma_secret"
	NOTIFY_API_ROOT = "https://magma-feg.magam.com/npcf-smpolicycontrol/v1"
	RULE_ID1        = "rule1"
	UM_DATA1        = "um_data1"
	POLICY_ID       = "12345"
)

var (
	UsageTx1    uint64 = 3000000
	UsageRx1    uint64 = 7000000
	UsageTotal1 uint64 = UsageTx1 + UsageRx1
)

func TestPCRFExpectations(t *testing.T) {
	n7Config := getTestN7Config(t)
	defaultAns := ""

	expectedPolicyContext := `{
		"accessType":      "3GPP_ACCESS",
		"dnn":             "apn.magma.com",
		"gpsi":            "9876543210",
		"ipv4Address":     "10.1.2.3",
		"pduSessionId":    10,
		"pduSessionType":  "IPV4",
		"ratType":         "NR",
		"supi":            "123456789012345",
		"notificationUri": "https://magma-feg.magam.com/npcf-smpolicycontrol/v1"
	}`
	createAns := `{
		"context": %s,
		"policy": %s
	}`
	createAns = fmt.Sprintf(createAns, expectedPolicyContext, createPolicyDecision())

	expectedUpdatePolicyContext := `{
		"accessType":      "3GPP_ACCESS",
		"ratType":         "NR",
		"accuUsageReports": [%s]
	}`
	expectedUpdatePolicyContext = fmt.Sprintf(expectedUpdatePolicyContext, createAccuUsageReport())
	updateAns := createPolicyUpdateAnswer()
	expectedDeletePolicyContext := fmt.Sprintf(`{"accuUsageReports": [%s]}`, createAccuUsageReport())

	mockPcf := startServerWithExpectations(
		t, n7Config,
		[]*protos.N7Expectation{
			{
				RequestType:     protos.N7Expectation_CREATE,
				ExpectedRequest: expectedPolicyContext,
				Answer:          createAns,
			},
			{
				RequestType:     protos.N7Expectation_UPDATE,
				ExpectedRequest: expectedUpdatePolicyContext,
				Answer:          updateAns,
			},
			{
				RequestType:     protos.N7Expectation_TERMINATE,
				ExpectedRequest: expectedDeletePolicyContext,
				Answer:          "",
			},
		},
		protos.UnexpectedRequestBehavior_CONTINUE_WITH_ERROR,
		defaultAns,
	)

	listenAddr := mockPcf.Server.ListenerAddr().String()
	// Create SM Policy
	resp, err := sendCreateSmPolicy(listenAddr, expectedPolicyContext)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, createAns, string(body))

	// Update SM Policy
	resp, err = sendUpdateSmPolicy(listenAddr, POLICY_ID, expectedUpdatePolicyContext)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, updateAns, string(body))

	// Delete SM Policy
	resp, err = sendDeleteSmPolicy(listenAddr, POLICY_ID, expectedDeletePolicyContext)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Test if the requests are correct
	expectedResult := []*protos.ExpectationResult{
		{ExpectationMet: true, ExpectationIndex: 0},
		{ExpectationMet: true, ExpectationIndex: 1},
		{ExpectationMet: true, ExpectationIndex: 2},
	}
	result, err := mockPcf.AssertExpectations(context.Background(), &orcprotos.Void{})
	assert.Nil(t, err)
	assert.ElementsMatch(t, expectedResult, result.Results)
}

func startServerWithExpectations(
	t *testing.T,
	serverConfig *sbi.NotifierConfig,
	expectations []*protos.N7Expectation,
	failureBehavior protos.UnexpectedRequestBehavior,
	defaultAns string,
) *servicers.MockPCFServer {
	mockPcf, err := servicers.NewMockPCFServer(serverConfig)
	require.NoError(t, err)
	ctx := context.Background()
	mockPcf.SetPCFConfigs(ctx, &protos.PCFConfigs{UseMockDriver: true})
	mockPcf.SetExpectations(ctx, &protos.N7Expectations{
		Expectations:              expectations,
		UnexpectedRequestBehavior: failureBehavior,
		DefaultAnswer:             defaultAns,
	})
	return mockPcf
}

func getTestN7Config(t *testing.T) *sbi.NotifierConfig {
	return &sbi.NotifierConfig{
		LocalAddr:     LOCAL_ADDR,
		NotifyApiRoot: API_ROOT,
	}
}

func sendCreateSmPolicy(pcfAddr string, payload string) (*http.Response, error) {
	postUrl := fmt.Sprintf("http://%s%s/sm-policies", pcfAddr, BASE_PATH)
	return http.Post(postUrl, "application/json", bytes.NewBuffer([]byte(payload)))
}

func sendUpdateSmPolicy(pcfAddr string, smPolicyId string, payload string) (*http.Response, error) {
	postUrl := fmt.Sprintf("http://%s%s/sm-policies/%s/update", pcfAddr, BASE_PATH, smPolicyId)
	return http.Post(postUrl, "application/json", bytes.NewBuffer([]byte(payload)))
}

func sendDeleteSmPolicy(pcfAddr string, smPolicyId string, payload string) (*http.Response, error) {
	postUrl := fmt.Sprintf("http://%s%s/sm-policies/%s/delete", pcfAddr, BASE_PATH, smPolicyId)
	return http.Post(postUrl, "application/json", bytes.NewBuffer([]byte(payload)))
}

func createPolicyDecision() string {
	policyDecisionStr := `{
		"pccRules": {
			"rule1": {
				"pccRuleId": "%s",
				"flowInfos": [{
					"flowDescription": "permit in ip from 0.0.0.0/0 to 4.2.2.4"
				}],
				"precedence": 1,
				"refQosData": ["qos_data1"],
				"refTcData": ["tc_data1"],
				"refChgData": ["chg_data1"],
				"refUmData": ["%s"],
				"refCondData": "cond_data1"
			},
			"static_rule1": {
				"pccRuleId": "static_rule1",
				"refCondData": "cond_data1"
			}
		},
		"qosDesc": {
			"qos_data1": {
				"qodId": "qos_data1",
				"5qi": 5,
				"maxbrUl": "200000",
				"maxbrDl": "500000",
				"gbrUl": "100000",
				"gbrDl": "250000"
			}
		},
		"chgDecs": {
			"chg_data1": {
				"chgId": "chg_data1",
				"online": true,
				"ratingGroup": 1,
				"serviceId": 12
			}
		},
		"traffContDecs": {
			"tc_data1": {
				"tcId": "tc_data1",
				"redirectInfo": {
					"redirectEnabled": true,
					"redirectAddressType": "URL",
					"redirectServerAddress": "https://redirect.example.com/tc"
				}
			}
		},
		"umDecs": {
			"%s": {
				"umId": "%s",
				"volumeThreshold": 4000000,
				"volumeThresholdUplink": 1500000,
				"volumeThresholdDownlink": 3500000
			}
		},
		"conds": {
			"cond_data1": {
				"condId": "cond_data1",
				"activationTime": "2021-10-22T12:42:31Z",
				"deactivationTime": "2021-10-22T14:42:31Z"
			}
		},
		"policyCtrlReqTriggers": ["RE_TIMEOUT"],
		"revalidationTime": "2021-10-22T14:42:31Z",
		"online": true
	}`
	return fmt.Sprintf(policyDecisionStr, RULE_ID1, UM_DATA1, UM_DATA1, UM_DATA1)
}

func createPolicyUpdateAnswer() string {
	policyDecisionStr := `{
		"pccRules": {
			"rule1": {
				"pccRuleId": "%s",
				"flowInfos": [{
					"flowDescription": "permit in ip from 0.0.0.0/0 to 4.2.2.4"
				}],
				"precedence": 1,
				"refUmData": ["%s"]
			}
		},
		"umDecs": {
			"%s": {
				"umId": "%s",
				"volumeThreshold": 4000000,
				"volumeThresholdUplink": 1500000,
				"volumeThresholdDownlink": 3500000
			}
		}
		"policyCtrlReqTriggers": ["RE_TIMEOUT"],
		"revalidationTime": "2021-10-22T14:42:31Z",
		"online": true
	}`
	return fmt.Sprintf(policyDecisionStr, RULE_ID1, UM_DATA1, UM_DATA1, UM_DATA1)
}

func createAccuUsageReport() string {
	usageReport := `{
		"refUmIds": "%s",
		"volUsageUplink": %d,
		"volUsageDownlink": %d,
		"volUsage": %d
	}`
	return fmt.Sprintf(usageReport, UM_DATA1, UsageTx1, UsageRx1, UsageTotal1)
}
