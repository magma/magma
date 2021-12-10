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

package n7

import (
	b64 "encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/protobuf/types/known/timestamppb"

	"magma/feg/gateway/policydb"
	common5g "magma/feg/gateway/sbi/specs/TS29122CommonData"
	n7_sbi "magma/feg/gateway/sbi/specs/TS29512NpcfSMPolicyControl"
	sbi "magma/feg/gateway/sbi/specs/TS29571CommonData"
	"magma/lte/cloud/go/protos"
)

var FailureCodeProtoToN7Map = []n7_sbi.FailureCode{
	"UNUSED",                 // 0
	"UNK_RULE_ID",            // 1
	"RA_GR_ERR",              // 2
	"SER_ID_ERR",             // 3
	"NF_MAL",                 // 4
	"RES_LIM",                // 5
	"MAX_NR_QoS_FLOW",        // 6
	"UNKNOWN",                // 7
	"UNKNOWN",                // 8
	"MISS_FLOW_INFO",         // 10
	"RES_ALLO_FAIL",          // 11
	"UNSUCC_QOS_VAL",         // 12
	"INCOR_FLOW_INFO",        // 13
	"PS_TO_CS_HAN",           // 14
	"APP_ID_ERR",             // 15
	"NO_QOS_FLOW_BOUND",      // 16
	"FILTER_RES",             // 17
	"UNKNOWN",                // 18
	"MISS_REDI_SER_ADDR",     // 19
	"CM_END_USER_SER_DENIED", // 20
	"CM_CREDIT_CON_NOT_APP",  // 21
	"CM_AUTH_REJ",            // 22
	"CM_USER_UNK",            // 23
	"CM_RAT_FAILED",          // 24
	"UNKNOWN",                // 25
	"UNKNOWN",                // 26
}

//
// From Proto to SBI types
//

func GetSmPolicyContextDataN7(
	request *protos.CreateSessionRequest,
	notifyApiRoot string,
) *n7_sbi.PostSmPoliciesJSONRequestBody {
	imsi := removeIMSIPrefix(request.GetCommonContext().GetSid().Id)
	common := request.GetCommonContext()
	ratType := common.GetRatType()

	reqBody := &n7_sbi.PostSmPoliciesJSONRequestBody{
		Supi:              sbi.Supi(imsi),
		Ipv4Address:       getSbiIpv4(common.GetUeIpv4()),
		Ipv6AddressPrefix: getSbiIpv6(common.GetUeIpv6()),
		Dnn:               sbi.Dnn(common.GetApn()),
		Gpsi:              getSbiGpsi(string(common.GetMsisdn())),
		RatType:           getSbiRatType(ratType),
		AccessType:        getSbiAccessType(ratType),
		UeTimeZone:        GetSbiTimeZone(request.GetAccessTimezone()),
		NotificationUri:   GenNotifyUrl(notifyApiRoot, request.SessionId),
	}
	if request.RatSpecificContext != nil {
		ratSpecific := request.GetRatSpecificContext().GetContext()
		switch context := ratSpecific.(type) {
		case *protos.RatSpecificContext_M5GsmSessionContext:
			m5gCtx := context.M5GsmSessionContext
			reqBody.PduSessionId = sbi.PduSessionId(m5gCtx.GetPduSessionId())
			reqBody.Gpsi = getSbiGpsi(m5gCtx.GetGpsi())
			reqBody.PduSessionType = getSbiPduSesionType(m5gCtx.GetPduSessionType())
		}
	}

	return reqBody
}

// IsReAuthSuccess checks if the PolicyReAuthAnswer received is successful
func IsReAuthSuccess(raa *protos.PolicyReAuthAnswer) bool {
	if raa.Result == protos.ReAuthResult_UPDATE_INITIATED ||
		raa.Result == protos.ReAuthResult_UPDATE_NOT_NEEDED ||
		len(raa.FailedRules) == 0 {
		// no errors to send
		return true
	}
	return false
}

// BuildPartialSuccessReportN7 returns a PartialSuccessReport with the list of rules that had
// failed ot install. If there are no FailedRules returns an empty RuleReports.
func BuildPartialSuccessReportN7(raa *protos.PolicyReAuthAnswer) *n7_sbi.PartialSuccessReport {
	ruleReports := []n7_sbi.RuleReport{}
	for ruleId, failureCode := range raa.FailedRules {
		ruleReport := n7_sbi.RuleReport{
			PccRuleIds:  []string{ruleId},
			RuleStatus:  n7_sbi.RuleStatusINACTIVE,
			FailureCode: getFailureCodeN7(failureCode),
		}
		ruleReports = append(ruleReports, ruleReport)
	}

	return &n7_sbi.PartialSuccessReport{
		FailureCause: n7_sbi.FailureCausePCCRULEEVENT,
		RuleReports:  &ruleReports,
	}
}

func GetSmPolicyDeleteReqBody(
	request *protos.SessionTerminateRequest,
) *n7_sbi.PostSmPoliciesSmPolicyIdDeleteJSONRequestBody {
	return &n7_sbi.PostSmPoliciesSmPolicyIdDeleteJSONRequestBody{
		AccuUsageReports: getUmReportN7(request.MonitorUsages),
	}
}

func getSbiRatType(ratType protos.RATType) *sbi.RatType {
	var sbiRatType sbi.RatType
	switch ratType {
	case protos.RATType_TGPP_LTE:
		sbiRatType = sbi.RatTypeEUTRA
	case protos.RATType_TGPP_NR:
		sbiRatType = sbi.RatTypeNR
	case protos.RATType_TGPP_WLAN:
		sbiRatType = sbi.RatTypeWLAN
	default:
		return nil
	}
	return &sbiRatType
}

func getSbiAccessType(ratType protos.RATType) *sbi.AccessType {
	var accessType sbi.AccessType
	switch ratType {
	case protos.RATType_TGPP_LTE:
		accessType = sbi.AccessTypeN3GPPACCESS
	case protos.RATType_TGPP_NR:
		accessType = sbi.AccessTypeN3GPPACCESS
	case protos.RATType_TGPP_WLAN:
		accessType = sbi.AccessTypeNON3GPPACCESS
	default:
		return nil
	}
	return &accessType
}

func GetSbiTimeZone(tzProto *protos.Timezone) *sbi.TimeZone {
	if tzProto == nil {
		return nil
	}
	offsetMin := tzProto.GetOffsetMinutes()
	isNegative := false
	if offsetMin < 0 {
		isNegative = true
		offsetMin = -offsetMin
	}
	tz := fmt.Sprintf("%02d:%02d", offsetMin/60, offsetMin%60)
	if isNegative {
		tz = "-" + tz
	} else {
		tz = "+" + tz
	}
	tzStr := sbi.TimeZone(tz)
	return &tzStr
}

func getSbiPduSesionType(pduSessionType protos.PduSessionType) sbi.PduSessionType {
	switch pduSessionType {
	case protos.PduSessionType_IPV4:
		return sbi.PduSessionTypeIPV4
	case protos.PduSessionType_IPV4IPV6:
		return sbi.PduSessionTypeIPV4V6
	case protos.PduSessionType_IPV6:
		return sbi.PduSessionTypeIPV6
	case protos.PduSessionType_UNSTRUCTURED:
		return sbi.PduSessionTypeUNSTRUCTURED
	default:
		return sbi.PduSessionTypeIPV4
	}
}

func getFailureCodeN7(failureCodeProto protos.PolicyReAuthAnswer_FailureCode) *n7_sbi.FailureCode {
	failureCodeNum := int(failureCodeProto)
	if failureCodeNum > len(FailureCodeProtoToN7Map) {
		return nil
	}
	retFailCode := FailureCodeProtoToN7Map[failureCodeNum]
	return &retFailCode
}

func getUmReportN7(umUpdates []*protos.UsageMonitorUpdate) *[]n7_sbi.AccuUsageReport {
	reports := []n7_sbi.AccuUsageReport{}

	for _, umUpdate := range umUpdates {
		reports = append(reports, getAccUsageReportN7(umUpdate))
	}

	return &reports
}

func getAccUsageReportN7(umUpdate *protos.UsageMonitorUpdate) n7_sbi.AccuUsageReport {
	return n7_sbi.AccuUsageReport{
		RefUmIds:         string(umUpdate.MonitoringKey),
		VolUsageDownlink: GetSbiVolume(umUpdate.BytesRx), // Output == Rx == Downlink
		VolUsageUplink:   GetSbiVolume(umUpdate.BytesTx), // Input == Tx == Uplink
		VolUsage:         GetSbiVolume(umUpdate.BytesRx + umUpdate.BytesTx),
	}
}

func GetSbiVolume(bytes uint64) *common5g.Volume {
	outBytes := common5g.Volume(int64(bytes))
	return &outBytes
}

// From SBI types to proto

func GetCreateSessionResponseProto(
	request *protos.CreateSessionRequest,
	smPolicyDecision *n7_sbi.SmPolicyDecision,
	smPolicyUrl string,
) *protos.CreateSessionResponse {
	staticRules, dynamicRules, _ := getPccRules(smPolicyDecision)
	eventTriggers, revalidationTime := getEventTriggersInfoProto(
		smPolicyDecision.PolicyCtrlReqTriggers, smPolicyDecision.RevalidationTime)

	tgppCtx := &protos.TgppContext{GxDestHost: smPolicyUrl}
	monitors := getUsageMonitorsProto(request, tgppCtx, eventTriggers, revalidationTime, smPolicyDecision)

	rpcResp := &protos.CreateSessionResponse{
		StaticRules:      staticRules,
		DynamicRules:     dynamicRules,
		TgppCtx:          tgppCtx,
		SessionId:        request.SessionId,
		EventTriggers:    eventTriggers,
		RevalidationTime: revalidationTime,
		UsageMonitors:    monitors,
	}
	if smPolicyDecision.Online != nil {
		rpcResp.Online = *smPolicyDecision.Online
	}
	if smPolicyDecision.Offline != nil {
		rpcResp.Offline = *smPolicyDecision.Offline
	}
	return rpcResp
}

func GetPolicyReauthRequestProto(
	sessionId string,
	imsi string,
	smPolicyDecision *n7_sbi.SmPolicyDecision,
) *protos.PolicyReAuthRequest {
	staticRules, dynamicRules, rulesToRemove := getPccRules(smPolicyDecision)
	eventTriggers, revalidationTime := getEventTriggersInfoProto(
		smPolicyDecision.PolicyCtrlReqTriggers, smPolicyDecision.RevalidationTime)
	umCredits := getUsageMonitorsProtoForUpdateNotify(smPolicyDecision)
	return &protos.PolicyReAuthRequest{
		SessionId:              sessionId,
		Imsi:                   imsi,
		RulesToInstall:         staticRules,
		DynamicRulesToInstall:  dynamicRules,
		RulesToRemove:          rulesToRemove,
		EventTriggers:          eventTriggers,
		RevalidationTime:       revalidationTime,
		UsageMonitoringCredits: umCredits,
	}
}

func getPccRules(
	policyDecision *n7_sbi.SmPolicyDecision,
) ([]*protos.StaticRuleInstall, []*protos.DynamicRuleInstall, []string) {
	var (
		staticRules   []*protos.StaticRuleInstall
		dynamicRules  []*protos.DynamicRuleInstall
		rulesToRemove []string
	)

	if policyDecision.PccRules == nil {
		return staticRules, dynamicRules, rulesToRemove
	}

	for ruleId, pccRule := range policyDecision.PccRules.AdditionalProperties {
		if pccRule.PccRuleId == "" {
			// This rule to be removed
			rulesToRemove = append(rulesToRemove, ruleId)
			continue
		}
		activationTime, deactivationTime := getActivateDeactivationTime(policyDecision, &pccRule)
		policyRule := getPolicyRuleProto(policyDecision, &pccRule)
		if policyRule == nil {
			// This is a static rule. Static rule only has pccRuleId and RefCondData
			staticRule := &protos.StaticRuleInstall{
				RuleId:           ruleId,
				ActivationTime:   activationTime,
				DeactivationTime: deactivationTime,
			}
			staticRules = append(staticRules, staticRule)
			continue
		}

		dynamicRule := &protos.DynamicRuleInstall{
			PolicyRule:       policyRule,
			ActivationTime:   activationTime,
			DeactivationTime: deactivationTime,
		}
		dynamicRules = append(dynamicRules, dynamicRule)
	}

	return staticRules, dynamicRules, rulesToRemove
}

func getActivateDeactivationTime(
	policyDecision *n7_sbi.SmPolicyDecision,
	pccRule *n7_sbi.PccRule,
) (activationTime, deactivationTime *timestamp.Timestamp) {
	if pccRule.RefCondData == nil || *pccRule.RefCondData == "" {
		return
	}
	condData, found := policyDecision.Conds.Get(*pccRule.RefCondData)
	if !found {
		return
	}
	activationTime = ConvertToProtoTimeStamp(condData.ActivationTime)
	deactivationTime = ConvertToProtoTimeStamp(condData.DeactivationTime)
	return
}

func getPolicyRuleProto(
	policyDecision *n7_sbi.SmPolicyDecision,
	pccRule *n7_sbi.PccRule,
) *protos.PolicyRule {
	// For all Ref*Data Maximum one entry is expected.
	// These are declared as arrays for future compatibility
	policyRule := &protos.PolicyRule{
		Id:            pccRule.PccRuleId,
		Priority:      getSbiUinteger(pccRule.Precedence),
		MonitoringKey: getMonKeyFromPccRule(policyDecision, pccRule),
		Redirect:      getTcDataFromPccRule(policyDecision, pccRule),
		Qos:           getQosDataFromPccRule(policyDecision, pccRule),
		FlowList:      getFlowListFromPccRule(pccRule),
	}
	insertChargingDataToPolicyRule(policyDecision, pccRule, policyRule)
	policyRule.TrackingType = getTrackingTypeProto(len(policyRule.MonitoringKey) != 0, policyRule.RatingGroup != 0)

	if len(policyRule.MonitoringKey) == 0 && policyRule.Redirect == nil &&
		policyRule.Qos == nil && len(policyRule.FlowList) == 0 && policyRule.RatingGroup == 0 {
		// must be a static rule if these are not present
		return nil
	}
	return policyRule
}

func insertChargingDataToPolicyRule(
	policyDecision *n7_sbi.SmPolicyDecision,
	pccRule *n7_sbi.PccRule,
	policyRule *protos.PolicyRule,
) {
	if pccRule.RefChgData == nil || len(*pccRule.RefChgData) == 0 || policyDecision.ChgDecs == nil {
		return
	}
	chgData, found := policyDecision.ChgDecs.Get((*pccRule.RefChgData)[0])
	if !found {
		return
	}
	policyRule.RatingGroup = uint32(*chgData.RatingGroup)
	policyRule.ServiceIdentifier = &protos.ServiceIdentifier{
		Value: uint32(*chgData.ServiceId),
	}
	policyRule.Offline = getSbiBool(chgData.Offline)
	policyRule.Online = getSbiBool(chgData.Online)
}

func getMonKeyFromPccRule(
	policyDecision *n7_sbi.SmPolicyDecision,
	pccRule *n7_sbi.PccRule,
) []byte {
	if pccRule.RefUmData != nil && len(*pccRule.RefUmData) > 0 && policyDecision.UmDecs != nil {
		umData, found := policyDecision.UmDecs.Get((*pccRule.RefUmData)[0])
		if found {
			return []byte(umData.UmId)
		}
	}
	return []byte{}
}

// getTcDataFromPccRule extracts Traffic Control data (Redirect Information) from the SmPolicyDecision
// and returns a proto Redirect Information
func getTcDataFromPccRule(
	policyDecision *n7_sbi.SmPolicyDecision,
	pccRule *n7_sbi.PccRule,
) *protos.RedirectInformation {
	if pccRule.RefTcData != nil && len(*pccRule.RefTcData) > 0 && policyDecision.TraffContDecs != nil {
		tcData, found := policyDecision.TraffContDecs.Get((*pccRule.RefTcData)[0])
		if found {
			return getRedirectInfoProto(tcData.RedirectInfo)
		}
	}
	return nil
}

func getQosDataFromPccRule(
	policyDecision *n7_sbi.SmPolicyDecision,
	pccRule *n7_sbi.PccRule,
) *protos.FlowQos {
	if pccRule.RefQosData != nil && len(*pccRule.RefQosData) > 0 && policyDecision.QosDecs != nil {
		qosData, found := policyDecision.QosDecs.Get((*pccRule.RefQosData)[0])
		if found {
			return getFlowQosProto(&qosData)
		}
	}
	return nil
}

func getFlowListFromPccRule(pccRule *n7_sbi.PccRule) []*protos.FlowDescription {
	var flowList []*protos.FlowDescription
	if pccRule.FlowInfos == nil {
		return flowList
	}
	for _, flowInfo := range *pccRule.FlowInfos {
		if flowInfo.FlowDescription == nil {
			continue
		}
		flow, err := policydb.GetFlowDescriptionFromFlowString(string(*flowInfo.FlowDescription))
		if err != nil {
			glog.Errorf("Could not get flow for description %s : %s", *flowInfo.FlowDescription, err)
			continue
		}
		flowList = append(flowList, flow)
	}
	return flowList
}

func getTrackingTypeProto(monKeyPresent bool, ratingGrpPresent bool) protos.PolicyRule_TrackingType {
	if monKeyPresent && ratingGrpPresent {
		return protos.PolicyRule_OCS_AND_PCRF
	} else if monKeyPresent && !ratingGrpPresent {
		return protos.PolicyRule_ONLY_PCRF
	} else if !monKeyPresent && ratingGrpPresent {
		return protos.PolicyRule_ONLY_OCS
	} else {
		return protos.PolicyRule_NO_TRACKING
	}
}

func getRedirectInfoProto(redirInfo *n7_sbi.RedirectInformation) *protos.RedirectInformation {
	redirProto := &protos.RedirectInformation{}
	if redirInfo == nil {
		return redirProto
	}
	if redirInfo.RedirectEnabled != nil {
		redirProto.Support = protos.RedirectInformation_DISABLED
		if *redirInfo.RedirectEnabled {
			redirProto.Support = protos.RedirectInformation_ENABLED
		}
	}
	if redirInfo.RedirectAddressType != nil {
		redirProto.AddressType = getRedirectAddrTypeProto(redirInfo.RedirectAddressType)
	}
	if redirInfo.RedirectServerAddress != nil {
		redirProto.ServerAddress = *redirInfo.RedirectServerAddress
	}

	return redirProto
}

func getRedirectAddrTypeProto(
	redirAddrType *n7_sbi.RedirectAddressType,
) protos.RedirectInformation_AddressType {
	if redirAddrType == nil {
		return protos.RedirectInformation_IPv4
	}
	switch *redirAddrType {
	case n7_sbi.RedirectAddressTypeIPV4ADDR:
		return protos.RedirectInformation_IPv4
	case n7_sbi.RedirectAddressTypeIPV6ADDR:
		return protos.RedirectInformation_IPv6
	case n7_sbi.RedirectAddressTypeSIPURI:
		return protos.RedirectInformation_SIP_URI
	case n7_sbi.RedirectAddressTypeURL:
		return protos.RedirectInformation_URL
	default:
		glog.Error("Unknown Redirect Address Type, defaulting to IPv4(default value)")
		return protos.RedirectInformation_IPv4
	}
}

func getFlowQosProto(qosData *n7_sbi.QosData) *protos.FlowQos {
	if qosData == nil {
		return nil
	}
	return &protos.FlowQos{
		MaxReqBwDl: bitRateStringToUint32(qosData.MaxbrDl),
		MaxReqBwUl: bitRateStringToUint32(qosData.MaxbrUl),
		GbrUl:      bitRateStringToUint32(qosData.GbrUl),
		GbrDl:      bitRateStringToUint32(qosData.GbrDl),
		Qci:        n5qiToQci(qosData.N5qi),
	}
}

func getEventTriggersInfoProto(
	triggers *[]n7_sbi.PolicyControlRequestTrigger,
	revalTimeout *time.Time,
) ([]protos.EventTrigger, *timestamp.Timestamp) {
	var (
		protoRevalidationTime *timestamp.Timestamp
		protoTriggers         []protos.EventTrigger
	)
	if triggers == nil {
		return protoTriggers, protoRevalidationTime
	}
	for _, trigger := range *triggers {
		switch trigger {
		case n7_sbi.PolicyControlRequestTriggerRETIMEOUT:
			protoTriggers = append(protoTriggers, protos.EventTrigger_REVALIDATION_TIMEOUT)
			protoRevalidationTime = ConvertToProtoTimeStamp(revalTimeout)
		}
	}
	return protoTriggers, protoRevalidationTime
}

func getUsageMonitorsProto(
	request *protos.CreateSessionRequest,
	tgppCtx *protos.TgppContext,
	eventTriggers []protos.EventTrigger,
	revalidationTime *timestamp.Timestamp,
	smPolicyDecision *n7_sbi.SmPolicyDecision,
) []*protos.UsageMonitoringUpdateResponse {
	var monitors []*protos.UsageMonitoringUpdateResponse
	// Usage monitorting at pcc rule level
	if smPolicyDecision.PccRules != nil {
		for _, pccRule := range smPolicyDecision.PccRules.AdditionalProperties {
			if pccRule.RefUmData == nil || len((*pccRule.RefUmData)[0]) == 0 {
				continue
			}
			umEntry, found := smPolicyDecision.UmDecs.Get((*pccRule.RefUmData)[0])
			if !found {
				continue
			}
			umResp := newUsageMonitoringUpdateResponse(request, tgppCtx, eventTriggers, revalidationTime)
			umResp.Credit = getUsageMonitoringCreditsProto((*pccRule.RefUmData)[0], &umEntry, protos.MonitoringLevel_PCC_RULE_LEVEL)
			monitors = append(monitors, umResp)
		}
	}
	// Usage monitoring at session level
	if smPolicyDecision.SessRules != nil {
		for _, sessRule := range smPolicyDecision.SessRules.AdditionalProperties {
			if sessRule.RefUmData == nil || len(*sessRule.RefUmData) == 0 {
				continue
			}
			umEntry, found := smPolicyDecision.UmDecs.Get(*sessRule.RefUmData)
			if !found {
				continue
			}
			umResp := newUsageMonitoringUpdateResponse(request, tgppCtx, eventTriggers, revalidationTime)
			umResp.Credit = getUsageMonitoringCreditsProto(*sessRule.RefUmData, &umEntry, protos.MonitoringLevel_SESSION_LEVEL)
			monitors = append(monitors, umResp)
		}
	}

	return monitors
}

func newUsageMonitoringUpdateResponse(
	request *protos.CreateSessionRequest,
	tgppCtx *protos.TgppContext,
	eventTriggers []protos.EventTrigger,
	revalidationTime *timestamp.Timestamp,
) *protos.UsageMonitoringUpdateResponse {
	return &protos.UsageMonitoringUpdateResponse{
		SessionId:        request.SessionId,
		TgppCtx:          tgppCtx,
		Sid:              request.GetCommonContext().GetSid().Id,
		Success:          true,
		EventTriggers:    eventTriggers,
		RevalidationTime: revalidationTime,
	}
}

func getUsageMonitoringCreditsProto(
	umId string,
	umData *n7_sbi.UsageMonitoringData,
	monLevel protos.MonitoringLevel,
) *protos.UsageMonitoringCredit {
	monAction := protos.UsageMonitoringCredit_CONTINUE
	if len(umData.UmId) == 0 {
		// UmData not present for ID. It is supposed to be removed
		monAction = protos.UsageMonitoringCredit_DISABLE
	}
	return &protos.UsageMonitoringCredit{
		MonitoringKey: []byte(umId),
		Action:        monAction,
		Level:         monLevel,
		GrantedUnits:  getGrantedUnitsProto(umData),
	}
}

func getGrantedUnitsProto(umData *n7_sbi.UsageMonitoringData) *protos.GrantedUnits {
	if umData == nil {
		return &protos.GrantedUnits{
			Total: &protos.CreditUnit{IsValid: false},
			Tx:    &protos.CreditUnit{IsValid: false},
			Rx:    &protos.CreditUnit{IsValid: false},
		}
	}
	gsuProto := &protos.GrantedUnits{
		Total: getCreditUnitsProto(umData.VolumeThreshold),
		Tx:    getCreditUnitsProto(umData.VolumeThresholdUplink),   // Input == Tx == Uplink
		Rx:    getCreditUnitsProto(umData.VolumeThresholdDownlink), // Output == Rx == Downlink
	}

	return gsuProto
}

func getCreditUnitsProto(volume *common5g.VolumeRm) *protos.CreditUnit {
	if volume == nil {
		return &protos.CreditUnit{IsValid: false}
	}
	return &protos.CreditUnit{IsValid: true, Volume: uint64(*volume)}
}

func getUsageMonitorsProtoForUpdateNotify(smPolicyDecision *n7_sbi.SmPolicyDecision) []*protos.UsageMonitoringCredit {
	monitors := []*protos.UsageMonitoringCredit{}
	// Usage monitorting at pcc rule level
	if smPolicyDecision.PccRules != nil {
		for _, pccRule := range smPolicyDecision.PccRules.AdditionalProperties {
			if pccRule.RefUmData == nil || len((*pccRule.RefUmData)[0]) == 0 {
				continue
			}
			umEntry, found := smPolicyDecision.UmDecs.Get((*pccRule.RefUmData)[0])
			if !found {
				continue
			}
			monitors = append(monitors, getUsageMonitoringCreditsProto((*pccRule.RefUmData)[0], &umEntry, protos.MonitoringLevel_PCC_RULE_LEVEL))
		}
	}
	// Usage monitoring at session level
	if smPolicyDecision.SessRules != nil {
		for _, sessRule := range smPolicyDecision.SessRules.AdditionalProperties {
			if sessRule.RefUmData == nil || len(*sessRule.RefUmData) == 0 {
				continue
			}
			umEntry, found := smPolicyDecision.UmDecs.Get(*sessRule.RefUmData)
			if !found {
				continue
			}
			monitors = append(monitors, getUsageMonitoringCreditsProto(*sessRule.RefUmData, &umEntry, protos.MonitoringLevel_SESSION_LEVEL))
		}
	}

	return monitors
}

//
// Utility functions
//

func ConvertToProtoTimeStamp(srcTime *time.Time) *timestamp.Timestamp {
	if srcTime == nil {
		return nil
	}
	return timestamppb.New(*srcTime)
}

func GenNotifyUrl(apiRoot string, sessionId string) sbi.Uri {
	return sbi.Uri(fmt.Sprintf("%s/%s", apiRoot, b64.URLEncoding.EncodeToString([]byte(sessionId))))
}

// GetSmPolicyId returns the SmPolicyId from the TgppContext.
// PolicyUrl is of the form https://{pcf-host}/npcf-smpolicycontrol/v1/sm-policies/{smPolicyId}
func GetSmPolicyId(tgppCtx *protos.TgppContext) (string, error) {
	if tgppCtx == nil {
		return "", fmt.Errorf("couldn't get url from TgppContext: nil TgppCtx")
	}
	policyUrl := tgppCtx.GetGxDestHost()
	if len(policyUrl) == 0 {
		return "", fmt.Errorf("empty PolicyUrl in TgppCtx")
	}
	parsedUrl, err := url.Parse(policyUrl)
	if err != nil {
		return "", fmt.Errorf("policyUrl parse error: %s", err)
	}
	return path.Base(parsedUrl.Path), nil
}

func getSbiIpv4(ipv4 string) *sbi.Ipv4Addr {
	if len(ipv4) == 0 {
		return nil
	}
	ret := sbi.Ipv4Addr(ipv4)
	return &ret
}

func getSbiIpv6(ipv6 string) *sbi.Ipv6Prefix {
	if len(ipv6) == 0 {
		return nil
	}
	ret := sbi.Ipv6Prefix(ipv6)
	return &ret
}

func getSbiGpsi(gpsi string) *sbi.Gpsi {
	if len(gpsi) == 0 {
		return nil
	}
	ret := sbi.Gpsi(gpsi)
	return &ret
}

func getSbiBool(val *bool) bool {
	if val == nil {
		return false
	}
	return *val
}

func getSbiUinteger(val *sbi.Uinteger) uint32 {
	if val == nil {
		return 0
	}
	return uint32(*val)
}

func bitRateStringToUint32(br *sbi.BitRateRm) uint32 {
	if br == nil {
		return 0
	}
	convInt, _ := strconv.ParseUint(string(*br), 10, 32)
	return uint32(convInt)
}

func n5qiToQci(n5Qi *sbi.N5Qi) protos.FlowQos_Qci {
	if n5Qi == nil {
		return protos.FlowQos_QCI_0
	}
	return protos.FlowQos_Qci(*n5Qi)
}
