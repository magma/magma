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
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mockPolicyDB "magma/feg/gateway/policydb/mocks"
	n7_sbi "magma/feg/gateway/sbi/specs/TS29512NpcfSMPolicyControl"
	sbi "magma/feg/gateway/sbi/specs/TS29571CommonData"
	"magma/feg/gateway/services/n7_n40_proxy/n7"
	mockN7 "magma/feg/gateway/services/n7_n40_proxy/n7/mocks"
	"magma/feg/gateway/services/n7_n40_proxy/servicers"
	relay_mocks "magma/feg/gateway/services/session_proxy/relay/mocks"
	"magma/lte/cloud/go/protos"
)

const (
	URL1            = "https://mockpcf/npcf-smpolicycontrol/v1"
	TOKEN_URL       = "https://mockpcf/oauth2/token"
	CLIENT_ID       = "feg_magma_client"
	CLIENT_SECRET   = "feg_magma_secret"
	LOCAL_ADDR      = "127.0.0.1:10100"
	NOTIFY_API_ROOT = "https://magma-feg.magam.com/npcf-smpolicycontrol/v1"
	IMSI1           = "IMSI123456789012345"
	IMSI2           = "IMSI123456789012346"
	IMSI1_NOPREFIX  = "123456789012345"
	IMSI2_NOPREFIX  = "123456789012346"
	SESS_ID1        = IMSI1 + "-1234"
	SESS_ID2        = IMSI2 + "-1234"
	UE_IPV4         = "10.1.2.3"
	PDU_SESSION_ID  = 10
	GPSI1           = "9876543210"
	APN1            = "apn.magma.com"
	RULE_ID1        = "rule1"
	UM_DATA1        = "um_data1"
)

var (
	AccessType3gpp     = sbi.AccessType("3GPP_ACCESS")
	DnnMagma           = sbi.Dnn(APN1)
	Gpsi1              = sbi.Gpsi(GPSI1)
	UeIpv4             = sbi.Ipv4Addr(UE_IPV4)
	PduSessionTypeIpv4 = sbi.PduSessionType("IPV4")
	RatTypeNR          = sbi.RatType("NR")
	UeTzIST            = sbi.TimeZone("+05:30")
	SmPolicyUrl        = "https://example.com//npcf-smpolicycontrol/v1/sm-policies/12345"
	ActTime            = time.Unix(1634906551, 0)
	DeactTime          = time.Unix(1634913751, 0)
)

func TestCreateSession(t *testing.T) {
	srv, mockDb, mockN7 := createCentralSessionControllerForTest(t, false)
	defer srv.Close()

	mockN7.On("PostSmPoliciesWithResponse", mock.Anything,
		mock.MatchedBy(matchSmPolicyContextData(IMSI1_NOPREFIX)),
	).Return(createSmPolicyResponse(t), nil).Once()

	mockDb.On("GetOmnipresentRules").Return(
		[]string{"omni_rule_1", "omni_rule_2"}, []string{"base_10"}).Once()
	mockDb.On("GetRuleIDsForBaseNames", []string{"base_10"}).Return([]string{"base_rule_1", "base_rule_2"})

	csr := defaultCreateSessionRequest()
	response, err := srv.CreateSession(context.Background(), csr)
	require.NoError(t, err)
	mockN7.AssertExpectations(t)

	assert.Equal(t, SESS_ID1, response.SessionId)
	allRuleIds := []string{}
	for _, staticRule := range response.StaticRules {
		allRuleIds = append(allRuleIds, staticRule.RuleId)
	}
	assert.ElementsMatch(t, []string{"static_rule1", "omni_rule_1", "omni_rule_2", "base_rule_1", "base_rule_2"}, allRuleIds)
	require.Equal(t, 1, len(response.DynamicRules))
	assert.Equal(t, defaultDynamicRule(RULE_ID1, UM_DATA1), response.DynamicRules[0])
	require.Equal(t, 1, len(response.UsageMonitors))
	assert.Equal(t, defaultUsageMonitors(UM_DATA1), response.UsageMonitors[0])
	assert.Equal(t, &protos.TgppContext{GxDestHost: SmPolicyUrl}, response.TgppCtx)
	assert.Equal(t, []protos.EventTrigger{protos.EventTrigger_REVALIDATION_TIMEOUT}, response.EventTriggers)
	assert.True(t, response.Online)
}

func TestCreateSessionTimeout(t *testing.T) {
	srv, _, mockN7 := createCentralSessionControllerForTest(t, false)
	defer srv.Close()

	mockN7.On("PostSmPoliciesWithResponse", mock.Anything,
		mock.MatchedBy(matchSmPolicyContextData(IMSI1_NOPREFIX)),
	).Return(nil, &url.Error{Err: context.DeadlineExceeded}).Once()

	csr := defaultCreateSessionRequest()
	response, err := srv.CreateSession(context.Background(), csr)
	require.Error(t, err)
	mockN7.AssertExpectations(t)
	assert.Nil(t, response)
}

func TestCreateSessionErrResp(t *testing.T) {
	srv, _, mockN7 := createCentralSessionControllerForTest(t, false)
	defer srv.Close()

	mockN7.On("PostSmPoliciesWithResponse", mock.Anything,
		mock.MatchedBy(matchSmPolicyContextData(IMSI1_NOPREFIX)),
	).Return(&n7_sbi.PostSmPoliciesResponse{HTTPResponse: &http.Response{StatusCode: 400}}, nil).Once()

	csr := defaultCreateSessionRequest()
	response, err := srv.CreateSession(context.Background(), csr)
	require.Error(t, err)
	mockN7.AssertExpectations(t)
	assert.Nil(t, response)
}

func TestDisableN7Response(t *testing.T) {
	srv, mockDb, mockN7 := createCentralSessionControllerForTest(t, true)
	defer srv.Close()

	mockDb.On("GetOmnipresentRules").Return(
		[]string{"omni_rule_1", "omni_rule_2"}, []string{"base_10"}).Once()
	mockDb.On("GetRuleIDsForBaseNames", []string{"base_10"}).Return([]string{"base_rule_1", "base_rule_2"})

	csr := defaultCreateSessionRequest()
	response, err := srv.CreateSession(context.Background(), csr)
	require.NoError(t, err)
	mockDb.AssertExpectations(t)
	mockN7.AssertNotCalled(t, "PostSmPoliciesWithResponse")

	assert.Equal(t, SESS_ID1, response.SessionId)
	allRuleIds := []string{}
	for _, staticRule := range response.StaticRules {
		allRuleIds = append(allRuleIds, staticRule.RuleId)
	}
	assert.ElementsMatch(t, []string{"omni_rule_1", "omni_rule_2", "base_rule_1", "base_rule_2"}, allRuleIds)
}

func createCentralSessionControllerForTest(t *testing.T, disableN7 bool) (*servicers.CentralSessionController, *mockPolicyDB.PolicyDBClient, *mockN7.ClientWithResponsesInterface) {
	testN7Conf := getTestN7Config(t)
	testN7Conf.DisableN7 = disableN7
	mockPolicyDBClient := &mockPolicyDB.PolicyDBClient{}
	mockPolicyClient := &mockN7.ClientWithResponsesInterface{}
	_, cloudRegistry := relay_mocks.StartMockSessionProxyResponder(t)

	srv, err := servicers.NewCentralSessionController(testN7Conf, mockPolicyDBClient, mockPolicyClient, cloudRegistry)
	require.NoError(t, err)
	return srv, mockPolicyDBClient, mockPolicyClient
}

func getTestN7Config(t *testing.T) *n7.N7Config {
	apiRoot, err := url.ParseRequestURI(URL1)
	require.NoError(t, err)
	return &n7.N7Config{
		DisableN7: false,
		Server: n7.PCFConfig{
			ApiRoot:      *apiRoot,
			TokenUrl:     TOKEN_URL,
			ClientId:     CLIENT_ID,
			ClientSecret: CLIENT_SECRET,
		},
		Client: n7.N7ClientConfig{
			LocalAddr:     LOCAL_ADDR,
			NotifyApiRoot: NOTIFY_API_ROOT,
		},
	}
}

func defaultCreateSessionRequest() *protos.CreateSessionRequest {
	return &protos.CreateSessionRequest{
		CommonContext: &protos.CommonSessionContext{
			Sid:     &protos.SubscriberID{Id: IMSI1},
			RatType: protos.RATType_TGPP_NR,
			UeIpv4:  UE_IPV4,
			Apn:     APN1,
		},
		RatSpecificContext: &protos.RatSpecificContext{
			Context: &protos.RatSpecificContext_M5GsmSessionContext{
				M5GsmSessionContext: &protos.M5GSMSessionContext{
					PduSessionId:   uint32(PDU_SESSION_ID),
					Gpsi:           GPSI1,
					PduSessionType: protos.PduSessionType_IPV4,
				},
			},
		},
		SessionId:      SESS_ID1,
		AccessTimezone: &protos.Timezone{OffsetMinutes: 330},
	}
}

func matchSmPolicyContextData(imsi string) interface{} {
	expected := &n7_sbi.PostSmPoliciesJSONRequestBody{
		AccessType:      &AccessType3gpp,
		Dnn:             DnnMagma,
		Gpsi:            &Gpsi1,
		Ipv4Address:     &UeIpv4,
		PduSessionId:    sbi.PduSessionId(PDU_SESSION_ID),
		PduSessionType:  PduSessionTypeIpv4,
		RatType:         &RatTypeNR,
		Supi:            sbi.Supi(imsi),
		UeTimeZone:      &UeTzIST,
		NotificationUri: n7.GenNotifyUrl(NOTIFY_API_ROOT, SESS_ID1),
	}

	return func(reqBody n7_sbi.PostSmPoliciesJSONRequestBody) bool {
		return reflect.DeepEqual(expected, &reqBody)
	}
}

func createSmPolicyResponse(t *testing.T) *n7_sbi.PostSmPoliciesResponse {
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
	policyDecisionStr = fmt.Sprintf(policyDecisionStr, RULE_ID1, UM_DATA1, UM_DATA1, UM_DATA1)
	// Unmarshal json to openapi struct
	var policyDecision n7_sbi.SmPolicyDecision
	err := json.Unmarshal([]byte(policyDecisionStr), &policyDecision)
	require.NoError(t, err)

	return &n7_sbi.PostSmPoliciesResponse{
		JSON201: &policyDecision,
		HTTPResponse: &http.Response{
			StatusCode: 201,
			Header: map[string][]string{
				"Location": {SmPolicyUrl},
			},
		},
	}
}

func defaultDynamicRule(ruleId string, monKey string) *protos.DynamicRuleInstall {
	return &protos.DynamicRuleInstall{
		PolicyRule: &protos.PolicyRule{
			Id:            ruleId,
			Priority:      1,
			RatingGroup:   1,
			MonitoringKey: []byte(monKey),
			Redirect: &protos.RedirectInformation{
				Support:       protos.RedirectInformation_ENABLED,
				ServerAddress: "https://redirect.example.com/tc",
				AddressType:   protos.RedirectInformation_URL,
			},
			FlowList: []*protos.FlowDescription{{
				Match: &protos.FlowMatch{
					IpSrc: &protos.IPAddress{Address: []byte("0.0.0.0/0")},
					IpDst: &protos.IPAddress{Address: []byte("4.2.2.4")},
				},
			}},
			TrackingType:      protos.PolicyRule_OCS_AND_PCRF,
			ServiceIdentifier: &protos.ServiceIdentifier{Value: 12},
			Online:            true,
			Offline:           false,
		},
		ActivationTime:   n7.ConvertToProtoTimeStamp(&ActTime),
		DeactivationTime: n7.ConvertToProtoTimeStamp(&DeactTime),
	}
}

func defaultUsageMonitors(monKey string) *protos.UsageMonitoringUpdateResponse {
	return &protos.UsageMonitoringUpdateResponse{
		Credit: &protos.UsageMonitoringCredit{
			Action:        protos.UsageMonitoringCredit_CONTINUE,
			MonitoringKey: []byte(monKey),
			Level:         protos.MonitoringLevel_PCC_RULE_LEVEL,
			GrantedUnits: &protos.GrantedUnits{
				Total: &protos.CreditUnit{IsValid: true, Volume: 4000000},
				Tx:    &protos.CreditUnit{IsValid: true, Volume: 1500000},
				Rx:    &protos.CreditUnit{IsValid: true, Volume: 3500000},
			},
		},
		SessionId:        SESS_ID1,
		Sid:              IMSI1,
		Success:          true,
		EventTriggers:    []protos.EventTrigger{protos.EventTrigger_REVALIDATION_TIMEOUT},
		RevalidationTime: n7.ConvertToProtoTimeStamp(&DeactTime),
		TgppCtx:          &protos.TgppContext{GxDestHost: SmPolicyUrl},
	}
}
