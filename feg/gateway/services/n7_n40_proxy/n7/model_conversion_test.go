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

package n7_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	n7_sbi "magma/feg/gateway/sbi/specs/TS29512NpcfSMPolicyControl"
	sbi "magma/feg/gateway/sbi/specs/TS29571CommonData"
	"magma/feg/gateway/services/n7_n40_proxy/n7"
	"magma/lte/cloud/go/protos"
)

const (
	IMSI1              = "123456789012345"
	SESSION_ID1        = IMSI1 + "-1234"
	UE_IPV4            = "10.1.2.3"
	PDU_SESSION_ID     = 10
	GPSI1              = "9876543210"
	APN1               = "apn.magma.com"
	MON_KEY1           = "mon_key_1"
	MON_KEY2           = "mon_key_2"
	IMSI2              = "543210987654321"
	SESSION_ID2        = IMSI2 + "-1234"
	MON_KEY3           = "mon_key_3"
	MON_KEY4           = "mon_key_4"
	SM_POLICY_ID2      = "271827"
	SM_POLICY_ID3      = "315149"
	BASE_SM_POLICY_URL = "https://example.com//npcf-smpolicycontrol/v1/sm-policies/"
)

var (
	AccessType3gpp            = sbi.AccessType("3GPP_ACCESS")
	DnnMagma                  = sbi.Dnn(APN1)
	Gpsi1                     = sbi.Gpsi(GPSI1)
	UeIpv4                    = sbi.Ipv4Addr(UE_IPV4)
	PduSessionTypeIpv4        = sbi.PduSessionType("IPV4")
	RatTypeNR                 = sbi.RatType("NR")
	UeTzIST                   = sbi.TimeZone("+05:30")
	SmPolicyUrl               = "https://example.com//npcf-smpolicycontrol/v1/sm-policies/12345"
	ActTime                   = time.Unix(1634906551, 0)
	DeactTime                 = time.Unix(1634913751, 0)
	UsageTx1           uint64 = 3000000
	UsageRx1           uint64 = 7000000
	UsageTotal1        uint64 = UsageTx1 + UsageRx1
	UsageTx2           uint64 = 8000000
	UsageRx2           uint64 = 14000000
	UsageTotal2        uint64 = UsageTx2 + UsageRx2
	UnkRuleId                 = n7_sbi.FailureCodeUNKRULEID
	IncorrectFlow             = n7_sbi.FailureCodeINCORFLOWINFO
	SmPolicyUrl2              = BASE_SM_POLICY_URL + SM_POLICY_ID2
	SmPolicyUrl3              = BASE_SM_POLICY_URL + SM_POLICY_ID3
)

func TestSmPolicyContextFromProto(t *testing.T) {
	csrProto := &protos.CreateSessionRequest{
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
		SessionId:      SESSION_ID1,
		AccessTimezone: &protos.Timezone{OffsetMinutes: 330},
	}
	expected := &n7_sbi.PostSmPoliciesJSONRequestBody{
		AccessType:      &AccessType3gpp,
		Dnn:             DnnMagma,
		Gpsi:            &Gpsi1,
		Ipv4Address:     &UeIpv4,
		PduSessionId:    sbi.PduSessionId(PDU_SESSION_ID),
		PduSessionType:  PduSessionTypeIpv4,
		RatType:         &RatTypeNR,
		Supi:            sbi.Supi(IMSI1),
		UeTimeZone:      &UeTzIST,
		NotificationUri: n7.GenNotifyUrl(NOTIFY_API_ROOT, SESSION_ID1),
	}

	reqBody := n7.GetSmPolicyContextDataN7(csrProto, NOTIFY_API_ROOT)
	// Check if JSON conversion is successufl
	_, err := json.Marshal(reqBody)
	require.NoError(t, err)
	assert.Equal(t, expected, reqBody)
}

func TestCreateSessionResponseProto(t *testing.T) {
	policyDecisionStr := `{
		"pccRules": {
			"rule1": {
				"pccRuleId": "rule1",
				"flowInfos": [{
					"flowDescription": "permit in ip from 0.0.0.0/0 to 4.2.2.4"
				}],
				"precedence": 1,
				"refQosData": ["qos_data1"],
				"refTcData": ["tc_data1"],
				"refChgData": ["chg_data1"],
				"refUmData": ["um_data1"],
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
	var policyDecision n7_sbi.SmPolicyDecision
	err := json.Unmarshal([]byte(policyDecisionStr), &policyDecision)
	require.NoError(t, err)

	csrProto := &protos.CreateSessionRequest{
		CommonContext: &protos.CommonSessionContext{
			Sid: &protos.SubscriberID{Id: IMSI1},
		},
		SessionId: SESSION_ID1,
	}
	expected := &protos.CreateSessionResponse{
		SessionId: csrProto.SessionId,
		StaticRules: []*protos.StaticRuleInstall{{
			RuleId:           "static_rule1",
			ActivationTime:   n7.ConvertToProtoTimeStamp(&ActTime),
			DeactivationTime: n7.ConvertToProtoTimeStamp(&DeactTime),
		}},
		DynamicRules: []*protos.DynamicRuleInstall{{
			PolicyRule: &protos.PolicyRule{
				Id:            "rule1",
				Priority:      1,
				RatingGroup:   1,
				MonitoringKey: []byte("um_data1"),
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
		}},
		UsageMonitors: []*protos.UsageMonitoringUpdateResponse{{
			Credit: &protos.UsageMonitoringCredit{
				Action:        protos.UsageMonitoringCredit_CONTINUE,
				MonitoringKey: []byte("um_data1"),
				Level:         protos.MonitoringLevel_PCC_RULE_LEVEL,
				GrantedUnits: &protos.GrantedUnits{
					Total: &protos.CreditUnit{IsValid: true, Volume: 4000000},
					Tx:    &protos.CreditUnit{IsValid: true, Volume: 1500000},
					Rx:    &protos.CreditUnit{IsValid: true, Volume: 3500000},
				},
			},
			SessionId:        csrProto.SessionId,
			Sid:              csrProto.CommonContext.GetSid().Id,
			Success:          true,
			EventTriggers:    []protos.EventTrigger{protos.EventTrigger_REVALIDATION_TIMEOUT},
			RevalidationTime: n7.ConvertToProtoTimeStamp(&DeactTime),
			TgppCtx:          &protos.TgppContext{GxDestHost: SmPolicyUrl},
		}},
		TgppCtx:          &protos.TgppContext{GxDestHost: SmPolicyUrl},
		EventTriggers:    []protos.EventTrigger{protos.EventTrigger_REVALIDATION_TIMEOUT},
		RevalidationTime: n7.ConvertToProtoTimeStamp(&DeactTime),
		Online:           true,
	}

	csResp := n7.GetCreateSessionResponseProto(csrProto, &policyDecision, SmPolicyUrl)
	assert.Equal(t, expected, csResp)
}

func TestSmPoliciesEmptyFields(t *testing.T) {
	reqBody := &n7_sbi.PostSmPoliciesJSONRequestBody{
		Dnn:            DnnMagma,
		Gpsi:           nil,
		Ipv4Address:    nil,
		PduSessionId:   sbi.PduSessionId(PDU_SESSION_ID),
		PduSessionType: PduSessionTypeIpv4,
		RatType:        nil,
		Supi:           sbi.Supi(IMSI1),
		UeTimeZone:     nil,
	}
	expectedStr := `{"dnn":"apn.magma.com","notificationUri":"","pduSessionId":10,"pduSessionType":"IPV4","sliceInfo":{"sst":0},"supi":"123456789012345","traceReq":null}`
	jsonReq, err := json.Marshal(reqBody)
	require.NoError(t, err)
	assert.Equal(t, expectedStr, string(jsonReq))
}

func TestGetSbiTimezone(t *testing.T) {
	tz1 := protos.Timezone{OffsetMinutes: 30}
	sbiTz := n7.GetSbiTimeZone(&tz1)
	assert.Equal(t, "+00:30", string(*sbiTz))

	tz1 = protos.Timezone{OffsetMinutes: -(3 * 60)}
	sbiTz = n7.GetSbiTimeZone(&tz1)
	assert.Equal(t, "-03:00", string(*sbiTz))

	tz1 = protos.Timezone{OffsetMinutes: -(11*60 + 30)}
	sbiTz = n7.GetSbiTimeZone(&tz1)
	assert.Equal(t, "-11:30", string(*sbiTz))

	tz1 = protos.Timezone{OffsetMinutes: 10 * 60}
	sbiTz = n7.GetSbiTimeZone(&tz1)
	assert.Equal(t, "+10:00", string(*sbiTz))

	tz1 = protos.Timezone{OffsetMinutes: 0}
	sbiTz = n7.GetSbiTimeZone(&tz1)
	assert.Equal(t, "+00:00", string(*sbiTz))
}

func TestGetPolicyRARProto(t *testing.T) {
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
			"remove_rule1": null
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
	var policyDecision n7_sbi.SmPolicyDecision
	err := json.Unmarshal([]byte(policyDecisionStr), &policyDecision)
	require.NoError(t, err)

	policyRar := n7.GetPolicyReauthRequestProto(SESSION_ID1, IMSI1, &policyDecision)
	expected := &protos.PolicyReAuthRequest{
		RulesToInstall: []*protos.StaticRuleInstall{{
			RuleId:           "static_rule1",
			ActivationTime:   n7.ConvertToProtoTimeStamp(&ActTime),
			DeactivationTime: n7.ConvertToProtoTimeStamp(&DeactTime),
		}},
		DynamicRulesToInstall: []*protos.DynamicRuleInstall{{
			PolicyRule: &protos.PolicyRule{
				Id:            "rule1",
				Priority:      1,
				RatingGroup:   1,
				MonitoringKey: []byte("um_data1"),
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
		}},
		RulesToRemove:    []string{"remove_rule1"},
		EventTriggers:    []protos.EventTrigger{protos.EventTrigger_REVALIDATION_TIMEOUT},
		RevalidationTime: n7.ConvertToProtoTimeStamp(&DeactTime),
		UsageMonitoringCredits: []*protos.UsageMonitoringCredit{{
			Action:        protos.UsageMonitoringCredit_CONTINUE,
			MonitoringKey: []byte("um_data1"),
			Level:         protos.MonitoringLevel_PCC_RULE_LEVEL,
			GrantedUnits: &protos.GrantedUnits{
				Total: &protos.CreditUnit{IsValid: true, Volume: 4000000},
				Tx:    &protos.CreditUnit{IsValid: true, Volume: 1500000},
				Rx:    &protos.CreditUnit{IsValid: true, Volume: 3500000},
			},
		}},
		SessionId: SESSION_ID1,
		Imsi:      IMSI1,
	}
	assert.Equal(t, expected, policyRar)
}

func TestGetPartialSuccessReportN7(t *testing.T) {
	raa := &protos.PolicyReAuthAnswer{
		Result: protos.ReAuthResult_OTHER_FAILURE,
		FailedRules: map[string]protos.PolicyReAuthAnswer_FailureCode{
			"rule1": protos.PolicyReAuthAnswer_UNKNOWN_RULE_NAME,
			"rule2": protos.PolicyReAuthAnswer_INCORRECT_FLOW_INFORMATION,
		},
	}
	partialSuccessRep := n7.BuildPartialSuccessReportN7(raa)
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
	assert.Equal(t, n7_sbi.FailureCausePCCRULEEVENT, partialSuccessRep.FailureCause)
	assert.ElementsMatch(t, expectedReports, *partialSuccessRep.RuleReports)
}

func TestSmPolicyDeleteFromProto(t *testing.T) {
	termSessProto := &protos.SessionTerminateRequest{
		SessionId: SESSION_ID1,
		CommonContext: &protos.CommonSessionContext{
			Sid:     &protos.SubscriberID{Id: IMSI1},
			RatType: protos.RATType_TGPP_NR,
			UeIpv4:  UE_IPV4,
		},
		TgppCtx: &protos.TgppContext{GxDestHost: SmPolicyUrl},
		MonitorUsages: []*protos.UsageMonitorUpdate{
			{
				MonitoringKey: []byte(MON_KEY1),
				BytesTx:       UsageTx1,
				BytesRx:       UsageRx1,
			},
			{
				MonitoringKey: []byte(MON_KEY2),
				BytesTx:       UsageTx2,
				BytesRx:       UsageRx2,
			},
		},
	}
	expected := n7_sbi.PostSmPoliciesSmPolicyIdDeleteJSONRequestBody{
		AccuUsageReports: &[]n7_sbi.AccuUsageReport{
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
	}
	reqBody := n7.GetSmPolicyDeleteReqBody(termSessProto)
	// Check if JSON conversion is successful
	_, err := json.Marshal(reqBody)
	require.NoError(t, err)
	assert.Equal(t, &expected, reqBody)
}

func TestSmPolicyUpdateFromProto(t *testing.T) {
	umUpdateReqProto := []*protos.UsageMonitoringUpdateRequest{
		{
			Update: &protos.UsageMonitorUpdate{
				MonitoringKey: []byte(MON_KEY1),
				Level:         protos.MonitoringLevel_PCC_RULE_LEVEL,
				BytesTx:       UsageTx1,
				BytesRx:       UsageRx1,
			},
			SessionId:    SESSION_ID1,
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
			SessionId:    SESSION_ID1,
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
			SessionId:    SESSION_ID2,
			Sid:          IMSI2,
			TgppCtx:      &protos.TgppContext{GxDestHost: SmPolicyUrl3},
			EventTrigger: protos.EventTrigger_USAGE_REPORT,
			RatType:      protos.RATType_TGPP_NR,
		},
	}
	expectedCtxs := map[string]n7.SmPolicyUpdateReqCtx{
		SESSION_ID1: {
			SmPolicyId:    SM_POLICY_ID2,
			SessionId:     SESSION_ID1,
			IMSI:          IMSI1,
			MonitoringKey: []byte(MON_KEY1),
			Level:         protos.MonitoringLevel_PCC_RULE_LEVEL,
			TgppCtx:       &protos.TgppContext{GxDestHost: SmPolicyUrl2},
			ReqBody: &n7_sbi.PostSmPoliciesSmPolicyIdUpdateJSONRequestBody{
				RatType:    &RatTypeNR,
				AccessType: &AccessType3gpp,
				AccuUsageReports: &[]n7_sbi.AccuUsageReport{
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
		},
		SESSION_ID2: {
			SmPolicyId:    SM_POLICY_ID3,
			SessionId:     SESSION_ID2,
			IMSI:          IMSI2,
			MonitoringKey: []byte(MON_KEY3),
			Level:         protos.MonitoringLevel_PCC_RULE_LEVEL,
			TgppCtx:       &protos.TgppContext{GxDestHost: SmPolicyUrl3},
			ReqBody: &n7_sbi.PostSmPoliciesSmPolicyIdUpdateJSONRequestBody{
				RatType:    &RatTypeNR,
				AccessType: &AccessType3gpp,
				AccuUsageReports: &[]n7_sbi.AccuUsageReport{
					{
						RefUmIds:         MON_KEY3,
						VolUsageUplink:   n7.GetSbiVolume(UsageTx2),
						VolUsageDownlink: n7.GetSbiVolume(UsageRx2),
						VolUsage:         n7.GetSbiVolume(UsageTotal2),
					},
				},
			},
		},
	}

	updateCtxs := n7.GetSmPolicyUpdateRequestsN7(umUpdateReqProto)
	require.Equal(t, 2, len(updateCtxs))

	for _, updateCtx := range updateCtxs {
		_, err := json.Marshal(updateCtx.ReqBody)
		assert.NoError(t, err)
		expected, found := expectedCtxs[updateCtx.SessionId]
		assert.True(t, found)
		assert.Equal(t, &expected, updateCtx)
	}
}

func TestSmPolicyUpdateResponseFromProto(t *testing.T) {
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
	var policyDecision n7_sbi.SmPolicyDecision
	err := json.Unmarshal([]byte(policyDecisionStr), &policyDecision)
	require.NoError(t, err)

	smUpdateCtx := &n7.SmPolicyUpdateReqCtx{
		SmPolicyId:    SM_POLICY_ID3,
		SessionId:     SESSION_ID2,
		IMSI:          IMSI2,
		MonitoringKey: []byte("um_data1"),
		Level:         protos.MonitoringLevel_PCC_RULE_LEVEL,
		TgppCtx:       &protos.TgppContext{GxDestHost: SmPolicyUrl2},
	}
	responsesProto := n7.GetUsageMonitoringResponsesProto(smUpdateCtx, &policyDecision)

	expectedStaticRules := []string{"static_rule1"}
	for i, staticRule := range responsesProto[0].StaticRulesToInstall {
		assert.Equal(t, expectedStaticRule(expectedStaticRules[i]), staticRule)
	}
	expectedDynamicRules := []struct {
		ruleId string
		monKey string
	}{
		{ruleId: "rule1", monKey: "um_data1"},
	}
	for i, dynamicRule := range responsesProto[0].DynamicRulesToInstall {
		expRule := expectedDynamicRules[i]
		assert.Equal(t, expectedDynamicRule(expRule.ruleId, expRule.monKey), dynamicRule)
	}
	assert.ElementsMatch(t, []string{"remove_rule1", "remove_rule2"}, responsesProto[0].RulesToRemove)
	assert.Equal(t, []protos.EventTrigger{protos.EventTrigger_REVALIDATION_TIMEOUT}, responsesProto[0].EventTriggers)
	assert.Equal(t, expectedUmCredit("um_data1"), responsesProto[0].Credit)
	assert.Equal(t, n7.ConvertToProtoTimeStamp(&DeactTime), responsesProto[0].RevalidationTime)
	assert.Equal(t, smUpdateCtx.SessionId, responsesProto[0].SessionId)
	assert.Equal(t, smUpdateCtx.IMSI, responsesProto[0].Sid)
	assert.Equal(t, &protos.TgppContext{GxDestHost: SmPolicyUrl2}, responsesProto[0].TgppCtx)
	assert.True(t, responsesProto[0].Success)
}

func expectedStaticRule(ruleId string) *protos.StaticRuleInstall {
	return &protos.StaticRuleInstall{
		RuleId:           ruleId,
		ActivationTime:   n7.ConvertToProtoTimeStamp(&ActTime),
		DeactivationTime: n7.ConvertToProtoTimeStamp(&DeactTime),
	}
}

func expectedDynamicRule(ruleId string, monKey string) *protos.DynamicRuleInstall {
	return &protos.DynamicRuleInstall{
		PolicyRule: &protos.PolicyRule{
			Id:            ruleId,
			Priority:      1,
			RatingGroup:   1,
			MonitoringKey: []byte(monKey),
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

func expectedUmCredit(monKey string) *protos.UsageMonitoringCredit {
	return &protos.UsageMonitoringCredit{
		Action:        protos.UsageMonitoringCredit_CONTINUE,
		MonitoringKey: []byte(monKey),
		Level:         protos.MonitoringLevel_PCC_RULE_LEVEL,
		GrantedUnits: &protos.GrantedUnits{
			Total: &protos.CreditUnit{IsValid: true, Volume: 4000000},
			Tx:    &protos.CreditUnit{IsValid: true, Volume: 1500000},
			Rx:    &protos.CreditUnit{IsValid: true, Volume: 3500000},
		},
	}
}
