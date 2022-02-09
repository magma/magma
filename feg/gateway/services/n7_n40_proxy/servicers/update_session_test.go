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
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	sbi_NpcfSMPolicyControl "magma/feg/gateway/sbi/specs/TS29512NpcfSMPolicyControl"
	"magma/feg/gateway/services/n7_n40_proxy/n7"
	"magma/lte/cloud/go/protos"
)

const (
	IMSI3_NOPREFIX     = "543210987654321"
	MON_KEY3           = "mon_key_3"
	MON_KEY4           = "mon_key_4"
	SM_POLICY_ID2      = "271827"
	SM_POLICY_ID3      = "315149"
	BASE_SM_POLICY_URL = "https://example.com//npcf-smpolicycontrol/v1/sm-policies/"
	SESS_ID1_NO_PREFIX = IMSI1_NOPREFIX + "-1234"
	SESS_ID2_NO_PREFIX = IMSI2_NOPREFIX + "-1234"
	SESS_ID3_NO_PREFIX = IMSI3_NOPREFIX + "-1234"
)

var (
	SmPolicyUrl2 = BASE_SM_POLICY_URL + SM_POLICY_ID2
	SmPolicyUrl3 = BASE_SM_POLICY_URL + SM_POLICY_ID3
)

func TestOneUpdateSession(t *testing.T) {
	srv, _, mockN7 := createCentralSessionControllerForTest(t, false)
	defer srv.Close()

	mockN7.On("PostSmPoliciesSmPolicyIdUpdateWithResponse",
		mock.Anything,
		mock.MatchedBy(func(argSmPolicyId string) bool {
			return argSmPolicyId == SM_POLICY_ID2
		}),
		mock.MatchedBy(matchSmPolicyUpdate()),
	).Return(createSingleSmPolicyUpdateResponse(t), nil).Once()

	updateReqs := &protos.UpdateSessionRequest{
		UsageMonitors: createSingleUmUpdateReqProto(),
	}
	response, err := srv.UpdateSession(context.Background(), updateReqs)
	require.NoError(t, err)
	mockN7.AssertExpectations(t)
	require.Equal(t, 2, len(response.UsageMonitorResponses))
	for _, res := range response.UsageMonitorResponses {
		assert.True(t, res.Success)
	}
}

func TestMultiUpdateSession(t *testing.T) {
	srv, _, mockN7 := createCentralSessionControllerForTest(t, false)
	defer srv.Close()

	mockN7.On("PostSmPoliciesSmPolicyIdUpdateWithResponse",
		mock.Anything,
		mock.MatchedBy(func(argSmPolicyId string) bool {
			return argSmPolicyId == SM_POLICY_ID2 || argSmPolicyId == SM_POLICY_ID3
		}),
		mock.MatchedBy(matchSmPolicyUpdate()),
	).Return(newSmPolicyUpdateResponse(t), nil).Times(2)

	updateReqs := &protos.UpdateSessionRequest{
		UsageMonitors: createUmUpdateReqProto(),
	}
	response, err := srv.UpdateSession(context.Background(), updateReqs)
	require.NoError(t, err)
	mockN7.AssertExpectations(t)
	require.Equal(t, 2, len(response.UsageMonitorResponses))
	for _, res := range response.UsageMonitorResponses {
		assert.True(t, res.Success)
	}
}

func TestMultiUpdateTimeoutSession(t *testing.T) {
	srv, _, mockN7 := createCentralSessionControllerForTest(t, false)
	defer srv.Close()

	mockN7.On("PostSmPoliciesSmPolicyIdUpdateWithResponse",
		mock.Anything,
		mock.MatchedBy(func(argSmPolicyId string) bool {
			return argSmPolicyId == SM_POLICY_ID2 || argSmPolicyId == SM_POLICY_ID3
		}),
		mock.MatchedBy(matchSmPolicyUpdate()),
	).Return(newSmPolicyUpdateResponse(t), nil).Once()
	mockN7.On("PostSmPoliciesSmPolicyIdUpdateWithResponse",
		mock.Anything,
		mock.MatchedBy(func(argSmPolicyId string) bool {
			return argSmPolicyId == SM_POLICY_ID2 || argSmPolicyId == SM_POLICY_ID3
		}),
		mock.MatchedBy(matchSmPolicyUpdate()),
	).Return(nil, &url.Error{Err: context.DeadlineExceeded}).Once()

	updateReqs := &protos.UpdateSessionRequest{
		UsageMonitors: createUmUpdateReqProto(),
	}
	response, err := srv.UpdateSession(context.Background(), updateReqs)
	require.NoError(t, err)
	mockN7.AssertExpectations(t)
	require.Equal(t, 2, len(response.UsageMonitorResponses))
	successCount := 0
	failureCount := 0
	for _, res := range response.UsageMonitorResponses {
		if res.Success {
			successCount++
		} else {
			failureCount++
		}
	}
	assert.Equal(t, 1, successCount)
	assert.Equal(t, 1, failureCount)
}

func matchSmPolicyUpdate() interface{} {
	expectedReqs := map[string]*sbi_NpcfSMPolicyControl.PostSmPoliciesSmPolicyIdUpdateJSONRequestBody{
		SM_POLICY_ID2: {
			RatType:    &RatTypeNR,
			AccessType: &AccessType3gpp,
			AccuUsageReports: &[]sbi_NpcfSMPolicyControl.AccuUsageReport{
				{
					RefUmIds:         MON_KEY1,
					VolUsageUplink:   n7.GetSbiVolume(UsageTx1),
					VolUsageDownlink: n7.GetSbiVolume(UsageRx1),
					VolUsage:         n7.GetSbiVolume(UsageTotal1),
				},
				{
					RefUmIds:         MON_KEY2,
					VolUsageUplink:   n7.GetSbiVolume(UsageTx2),
					VolUsageDownlink: n7.GetSbiVolume(UsageRx2),
					VolUsage:         n7.GetSbiVolume(UsageTotal2),
				},
			},
		},
		SM_POLICY_ID3: {
			RatType:    &RatTypeNR,
			AccessType: &AccessType3gpp,
			AccuUsageReports: &[]sbi_NpcfSMPolicyControl.AccuUsageReport{
				{
					RefUmIds:         MON_KEY3,
					VolUsageUplink:   n7.GetSbiVolume(UsageTx2),
					VolUsageDownlink: n7.GetSbiVolume(UsageRx2),
					VolUsage:         n7.GetSbiVolume(UsageTotal2),
				},
			},
		},
	}

	return func(reqBody sbi_NpcfSMPolicyControl.PostSmPoliciesSmPolicyIdUpdateJSONRequestBody) bool {
		return reflect.DeepEqual(expectedReqs[SM_POLICY_ID2], &reqBody) ||
			reflect.DeepEqual(expectedReqs[SM_POLICY_ID3], &reqBody)
	}
}

func createUmUpdateReqProto() []*protos.UsageMonitoringUpdateRequest {
	return []*protos.UsageMonitoringUpdateRequest{
		{
			Update: &protos.UsageMonitorUpdate{
				MonitoringKey: []byte(MON_KEY1),
				Level:         protos.MonitoringLevel_PCC_RULE_LEVEL,
				BytesTx:       UsageTx1,
				BytesRx:       UsageRx1,
			},
			SessionId:    SESS_ID1_NO_PREFIX,
			Sid:          IMSI1,
			TgppCtx:      &protos.TgppContext{GxDestHost: SmPolicyUrl2},
			EventTrigger: protos.EventTrigger_USAGE_REPORT,
			RatType:      protos.RATType_TGPP_NR,
		},
		{
			Update: &protos.UsageMonitorUpdate{
				MonitoringKey: []byte(MON_KEY2),
				Level:         protos.MonitoringLevel_PCC_RULE_LEVEL,
				BytesTx:       UsageTx2,
				BytesRx:       UsageRx2,
			},
			SessionId:    SESS_ID1_NO_PREFIX,
			Sid:          IMSI1,
			TgppCtx:      &protos.TgppContext{GxDestHost: SmPolicyUrl2},
			EventTrigger: protos.EventTrigger_USAGE_REPORT,
			RatType:      protos.RATType_TGPP_NR,
		},
		{
			Update: &protos.UsageMonitorUpdate{
				MonitoringKey: []byte(MON_KEY3),
				Level:         protos.MonitoringLevel_PCC_RULE_LEVEL,
				BytesTx:       UsageTx2,
				BytesRx:       UsageRx2,
			},
			SessionId:    SESS_ID3_NO_PREFIX,
			Sid:          IMSI3_NOPREFIX,
			TgppCtx:      &protos.TgppContext{GxDestHost: SmPolicyUrl3},
			EventTrigger: protos.EventTrigger_USAGE_REPORT,
			RatType:      protos.RATType_TGPP_NR,
		},
	}
}

func createSingleUmUpdateReqProto() []*protos.UsageMonitoringUpdateRequest {
	return []*protos.UsageMonitoringUpdateRequest{
		{
			Update: &protos.UsageMonitorUpdate{
				MonitoringKey: []byte(MON_KEY1),
				Level:         protos.MonitoringLevel_PCC_RULE_LEVEL,
				BytesTx:       UsageTx1,
				BytesRx:       UsageRx1,
			},
			SessionId:    SESS_ID1_NO_PREFIX,
			Sid:          IMSI1,
			TgppCtx:      &protos.TgppContext{GxDestHost: SmPolicyUrl2},
			EventTrigger: protos.EventTrigger_USAGE_REPORT,
			RatType:      protos.RATType_TGPP_NR,
		},
		{
			Update: &protos.UsageMonitorUpdate{
				MonitoringKey: []byte(MON_KEY2),
				Level:         protos.MonitoringLevel_PCC_RULE_LEVEL,
				BytesTx:       UsageTx2,
				BytesRx:       UsageRx2,
			},
			SessionId:    SESS_ID1_NO_PREFIX,
			Sid:          IMSI1,
			TgppCtx:      &protos.TgppContext{GxDestHost: SmPolicyUrl2},
			EventTrigger: protos.EventTrigger_USAGE_REPORT,
			RatType:      protos.RATType_TGPP_NR,
		},
	}
}

func createSingleSmPolicyUpdateResponse(t *testing.T) *sbi_NpcfSMPolicyControl.PostSmPoliciesSmPolicyIdUpdateResponse {
	policyDecisionStr := `{
		"pccRules": {
			"rule1": {
				"pccRuleId": "rule1",
				"flowInfos": [{
					"flowDescription": "permit in ip from 0.0.0.0/0 to 4.2.2.4"
				}],
				"precedence": 1,
				"refUmData": ["mon_key1"],
				"refCondData": "cond_data1"
			},
			"rule2": {
				"pccRuleId": "rule2",
				"flowInfos": [{
					"flowDescription": "permit in ip from 0.0.0.0/0 to 4.2.2.4"
				}],
				"precedence": 1,
				"refUmData": ["mon_key2"],
				"refCondData": "cond_data1"
			}
		},
		"umDecs": {
			"mon_key1": {
				"umId": "mon_key1",
				"volumeThreshold": 4000000,
				"volumeThresholdUplink": 1500000,
				"volumeThresholdDownlink": 3500000
			},
			"mon_key2": {
				"umId": "mon_key2",
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
	// Unmarshal json to openapi struct
	var policyDecision sbi_NpcfSMPolicyControl.SmPolicyDecision
	err := json.Unmarshal([]byte(policyDecisionStr), &policyDecision)
	require.NoError(t, err)

	return &sbi_NpcfSMPolicyControl.PostSmPoliciesSmPolicyIdUpdateResponse{
		JSON200:      &policyDecision,
		HTTPResponse: &http.Response{StatusCode: 200},
	}
}

func newSmPolicyUpdateResponse(t *testing.T) *sbi_NpcfSMPolicyControl.PostSmPoliciesSmPolicyIdUpdateResponse {
	policyDecisionStr := `{
		"pccRules": {
			"rule1": {
				"pccRuleId": "rule1",
				"flowInfos": [{
					"flowDescription": "permit in ip from 0.0.0.0/0 to 4.2.2.4"
				}],
				"precedence": 1,
				"refQosData": ["qos_data1"],
				"refChgData": ["chg_data1"],
				"refUmData": ["um_data1"],
				"refCondData": "cond_data1"
			},
			"static_rule1": {
				"pccRuleId": "static_rule1",
				"refCondData": "cond_data1"
			},
			"remove_rule1": null,
			"remove_rule2": null
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
		"umDecs": {
			"um_data1": {
				"umId": "um_data1",
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
	// Unmarshal json to openapi struct
	var policyDecision sbi_NpcfSMPolicyControl.SmPolicyDecision
	err := json.Unmarshal([]byte(policyDecisionStr), &policyDecision)
	require.NoError(t, err)

	return &sbi_NpcfSMPolicyControl.PostSmPoliciesSmPolicyIdUpdateResponse{
		JSON200:      &policyDecision,
		HTTPResponse: &http.Response{StatusCode: 200},
	}
}
