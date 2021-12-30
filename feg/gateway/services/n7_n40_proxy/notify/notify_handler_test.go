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
package notify_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	n7_sbi "magma/feg/gateway/sbi/specs/TS29512NpcfSMPolicyControl"
	"magma/feg/gateway/services/n7_n40_proxy/n7"
	"magma/feg/gateway/services/n7_n40_proxy/notify"
	relay_mocks "magma/feg/gateway/services/session_proxy/relay/mocks"
	"magma/lte/cloud/go/protos"
)

const (
	LOCAL_ADDR = "127.0.0.1:0"
	BASE_PATH  = "/sm-policy-control/v1"
	HTTP_HOST  = "http://localhost"
	API_ROOT   = HTTP_HOST + BASE_PATH
	IMSI1      = "123456789012345"
	SESS_ID    = IMSI1 + "-987654321"
)

var (
	ActTime       = time.Unix(1634906551, 0)
	DeactTime     = time.Unix(1634913751, 0)
	UnkRuleId     = n7_sbi.FailureCodeUNKRULEID
	IncorrectFlow = n7_sbi.FailureCodeINCORFLOWINFO
)

func TestUpdateNotify(t *testing.T) {
	sm, cloudRegistry := relay_mocks.StartMockSessionProxyResponder(t)
	notifyHandler, err := notify.NewStartedNotificationHandlerWithHandlers(getClientConfig(), cloudRegistry)
	require.NoError(t, err)
	require.NoError(t, err)
	defer notifyHandler.Shutdown()

	// happy path
	sm.On("PolicyReAuth", mock.Anything, expectedPolicyReauth()).
		Return(&protos.PolicyReAuthAnswer{SessionId: SESS_ID}, nil).Once()

	resp, err := postUpdateNotify(notifyHandler.NotifyServer.ListenAddr, genUpdateNotifyStr())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Some failures, partial error report
	sm.On("PolicyReAuth", mock.Anything, expectedPolicyReauth()).
		Return(&protos.PolicyReAuthAnswer{
			SessionId: SESS_ID,
			Result:    protos.ReAuthResult_OTHER_FAILURE,
			FailedRules: map[string]protos.PolicyReAuthAnswer_FailureCode{
				"rule1": protos.PolicyReAuthAnswer_UNKNOWN_RULE_NAME,
				"rule2": protos.PolicyReAuthAnswer_INCORRECT_FLOW_INFORMATION,
			},
		}, nil).Once()

	resp, err = postUpdateNotify(notifyHandler.NotifyServer.ListenAddr, genUpdateNotifyStr())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var report n7_sbi.PartialSuccessReport
	err = json.Unmarshal(body, &report)
	assert.NoError(t, err)

	expectedReports := []n7_sbi.RuleReport{
		{
			PccRuleIds:  []string{"rule1"},
			RuleStatus:  n7_sbi.RuleStatusINACTIVE,
			FailureCode: &UnkRuleId,
		},
		{
			PccRuleIds:  []string{"rule2"},
			RuleStatus:  n7_sbi.RuleStatusINACTIVE,
			FailureCode: &IncorrectFlow,
		},
	}
	assert.Equal(t, n7_sbi.FailureCausePCCRULEEVENT, report.FailureCause)
	assert.ElementsMatch(t, expectedReports, *report.RuleReports)
}

func getClientConfig() *n7.N7ClientConfig {
	return &n7.N7ClientConfig{
		LocalAddr:     LOCAL_ADDR,
		NotifyApiRoot: API_ROOT,
	}
}

func postUpdateNotify(notifAddr string, payload string) (*http.Response, error) {
	apiRoot := fmt.Sprintf("http://%s%s", notifAddr, BASE_PATH)
	postUrl := string(n7.GenNotifyUrl(apiRoot, SESS_ID)) + "/update"
	return http.Post(postUrl, "application/json", bytes.NewBuffer([]byte(payload)))
}

func genUpdateNotifyStr() string {
	return `{
		"resourceUri": "https://mock.pcrf//sm-policy-control/v1/sm-policies/12345",
		"smPolicyDecision": {
			"pccRules": {
				"static_rule1": {
					"pccRuleId": "static_rule1",
					"refCondData": "cond_data1"
				},
				"remove_rule1": null
			},
			"conds": {
				"cond_data1": {
					"condId": "cond_data1",
					"activationTime": "2021-10-22T12:42:31Z",
					"deactivationTime": "2021-10-22T14:42:31Z"
				}
			},
			"policyCtrlReqTriggers": ["RE_TIMEOUT"],
			"revalidationTime": "2021-10-22T14:42:31Z"
		}
	}`
}

func expectedPolicyReauth() *protos.PolicyReAuthRequest {
	return &protos.PolicyReAuthRequest{
		RulesToInstall: []*protos.StaticRuleInstall{{
			RuleId:           "static_rule1",
			ActivationTime:   n7.ConvertToProtoTimeStamp(&ActTime),
			DeactivationTime: n7.ConvertToProtoTimeStamp(&DeactTime),
		}},
		RulesToRemove:    []string{"remove_rule1"},
		EventTriggers:    []protos.EventTrigger{protos.EventTrigger_REVALIDATION_TIMEOUT},
		RevalidationTime: n7.ConvertToProtoTimeStamp(&DeactTime),
		SessionId:        SESS_ID,
		Imsi:             IMSI1,
	}
}
