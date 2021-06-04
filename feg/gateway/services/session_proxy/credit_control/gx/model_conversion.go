/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gx

import (
	"time"

	"magma/feg/gateway/policydb"
	"magma/feg/gateway/services/session_proxy/credit_control"
	"magma/lte/cloud/go/protos"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/go-openapi/swag"
	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func (ccr *CreditControlRequest) FromUsageMonitorUpdate(update *protos.UsageMonitoringUpdateRequest) *CreditControlRequest {
	ccr.SessionID = update.SessionId
	ccr.TgppCtx = update.GetTgppCtx()
	ccr.RequestNumber = update.RequestNumber
	ccr.Type = credit_control.CRTUpdate
	ccr.IMSI = credit_control.RemoveIMSIPrefix(update.Sid)
	ccr.IPAddr = update.UeIpv4
	ccr.HardwareAddr = update.HardwareAddr
	if update.EventTrigger == protos.EventTrigger_USAGE_REPORT {
		ccr.UsageReports = []*UsageReport{(&UsageReport{}).FromUsageMonitorUpdate(update.Update)}
	}
	ccr.RATType = GetRATType(update.RatType)
	ccr.IPCANType = GetIPCANType(update.RatType)
	ccr.EventTrigger = EventTrigger(update.EventTrigger)
	ccr.ChargingCharacteristics = update.ChargingCharacteristics
	return ccr
}

func (qos *QosRequestInfo) FromProtos(pQos *protos.QosInformationRequest) *QosRequestInfo {
	switch pQos.BrUnit {

	// 3gpp 29.212, 4.5.30 Extended bandwidth support for EPC supporting Dual Connectivity
	case protos.QosInformationRequest_KBPS:
		qos.ApnExtendedAggMaxBitRateDL = pQos.GetApnAmbrDl()
		qos.ApnExtendedAggMaxBitRateUL = pQos.GetApnAmbrUl()
	default:
		qos.ApnAggMaxBitRateDL = pQos.GetApnAmbrDl()
		qos.ApnAggMaxBitRateUL = pQos.GetApnAmbrUl()
	}

	qos.PriLevel = pQos.GetPriorityLevel()
	qos.PreCapability = pQos.GetPreemptionCapability()
	qos.PreVulnerability = pQos.GetPreemptionVulnerability()
	return qos
}

func (rd *RuleDefinition) ToProto() *protos.PolicyRule {
	return &protos.PolicyRule{
		Id:            rd.RuleName,
		RatingGroup:   swag.Uint32Value(rd.RatingGroup),
		MonitoringKey: rd.MonitoringKey,
		Priority:      rd.Precedence,
		Redirect:      rd.RedirectInformation.ToProto(),
		FlowList:      rd.GetFlowList(),
		Qos:           rd.Qos.ToProto(),
		TrackingType:  rd.GetTrackingType(),
	}
}

func (q *QosInformation) ToProto() *protos.FlowQos {
	var qos *protos.FlowQos
	if q != nil {
		qos = &protos.FlowQos{
			MaxReqBwUl: swag.Uint32Value(q.MaxReqBwUL),
			MaxReqBwDl: swag.Uint32Value(q.MaxReqBwDL),
			GbrUl:      swag.Uint32Value(q.GbrUL),
			GbrDl:      swag.Uint32Value(q.GbrDL),
			Qci:        protos.FlowQos_Qci(swag.Uint32Value(q.Qci)),
		}
	}
	return qos
}

func (rd *RuleDefinition) GetTrackingType() protos.PolicyRule_TrackingType {
	monKeyPresent := len(rd.MonitoringKey) > 0
	if monKeyPresent && rd.RatingGroup != nil {
		return protos.PolicyRule_OCS_AND_PCRF
	} else if monKeyPresent && rd.RatingGroup == nil {
		return protos.PolicyRule_ONLY_PCRF
	} else if (!monKeyPresent) && rd.RatingGroup != nil {
		return protos.PolicyRule_ONLY_OCS
	} else {
		return protos.PolicyRule_NO_TRACKING
	}
}

func (r *RedirectInformation) ToProto() *protos.RedirectInformation {
	if r == nil {
		return &protos.RedirectInformation{}
	}
	return &protos.RedirectInformation{
		Support:       protos.RedirectInformation_Support(r.RedirectSupport),
		AddressType:   protos.RedirectInformation_AddressType(r.RedirectAddressType),
		ServerAddress: r.RedirectServerAddress,
	}
}

func (rd *RuleDefinition) GetFlowList() []*protos.FlowDescription {
	allFlowStrings := rd.FlowDescriptions[:]
	for _, info := range rd.FlowInformations {
		allFlowStrings = append(allFlowStrings, info.FlowDescription)
	}
	var flowList []*protos.FlowDescription
	for _, flowString := range allFlowStrings {
		flow, err := policydb.GetFlowDescriptionFromFlowString(flowString)
		if err != nil {
			glog.Errorf("Could not get flow for description %s : %s", flowString, err)
		} else {
			flowList = append(flowList, flow)
		}
	}
	return flowList
}

func (rar *PolicyReAuthRequest) ToProto(imsi, sid string, policyDBClient policydb.PolicyDBClient) *protos.PolicyReAuthRequest {
	var rulesToRemove, baseNamesToRemove []string

	for _, ruleRemove := range rar.RulesToRemove {
		rulesToRemove = append(rulesToRemove, ruleRemove.RuleNames...)
		baseNamesToRemove = append(baseNamesToRemove, ruleRemove.RuleBaseNames...)
	}

	baseNameRuleIDsToRemove := policyDBClient.GetRuleIDsForBaseNames(baseNamesToRemove)
	rulesToRemove = append(rulesToRemove, baseNameRuleIDsToRemove...)

	staticRulesToInstall, dynamicRulesToInstall := ParseRuleInstallAVPs(
		policyDBClient,
		rar.RulesToInstall,
	)

	eventTriggers, revalidationTime := GetEventTriggersRelatedInfo(rar.EventTriggers, rar.RevalidationTime)
	usageMonitoringCredits := getUsageMonitoringCredits(rar.UsageMonitors)
	qosInfo := getQoSInfo(rar.Qos)

	return &protos.PolicyReAuthRequest{
		SessionId:              sid,
		Imsi:                   imsi,
		RulesToRemove:          rulesToRemove,
		RulesToInstall:         staticRulesToInstall,
		DynamicRulesToInstall:  dynamicRulesToInstall,
		EventTriggers:          eventTriggers,
		RevalidationTime:       revalidationTime,
		UsageMonitoringCredits: usageMonitoringCredits,
		QosInfo:                qosInfo,
	}
}

func (raa *PolicyReAuthAnswer) FromProto(sessionID string, answer *protos.PolicyReAuthAnswer) *PolicyReAuthAnswer {
	raa.SessionID = sessionID
	raa.ResultCode = diam.Success
	raa.RuleReports = make([]*ChargingRuleReport, 0, len(answer.FailedRules))
	for ruleName, code := range answer.FailedRules {
		raa.RuleReports = append(
			raa.RuleReports,
			&ChargingRuleReport{RuleNames: []string{ruleName}, FailureCode: RuleFailureCode(code)},
		)
	}
	return raa
}

func ConvertToProtoTimestamp(unixTime *time.Time) *timestamp.Timestamp {
	if unixTime == nil {
		return nil
	}
	protoTimestamp, err := ptypes.TimestampProto(*unixTime)
	if err != nil {
		glog.Errorf("Unable to convert time.Time to google.protobuf.Timestamp: %s", err)
		return nil
	}
	return protoTimestamp
}

func ParseRuleInstallAVPs(
	policyDBClient policydb.PolicyDBClient,
	ruleInstalls []*RuleInstallAVP,
) ([]*protos.StaticRuleInstall, []*protos.DynamicRuleInstall) {
	staticRulesToInstall := make([]*protos.StaticRuleInstall, 0, len(ruleInstalls))
	dynamicRulesToInstall := make([]*protos.DynamicRuleInstall, 0, len(ruleInstalls))
	for _, ruleInstall := range ruleInstalls {
		activationTime := ConvertToProtoTimestamp(ruleInstall.RuleActivationTime)
		deactivationTime := ConvertToProtoTimestamp(ruleInstall.RuleDeactivationTime)

		for _, staticRuleName := range ruleInstall.RuleNames {
			staticRulesToInstall = append(
				staticRulesToInstall,
				&protos.StaticRuleInstall{
					RuleId:           staticRuleName,
					ActivationTime:   activationTime,
					DeactivationTime: deactivationTime,
				},
			)
		}

		if len(ruleInstall.RuleBaseNames) != 0 {
			baseNameRuleIdsToInstall := policyDBClient.GetRuleIDsForBaseNames(ruleInstall.RuleBaseNames)
			for _, baseNameRuleId := range baseNameRuleIdsToInstall {
				staticRulesToInstall = append(
					staticRulesToInstall,
					&protos.StaticRuleInstall{
						RuleId:           baseNameRuleId,
						ActivationTime:   activationTime,
						DeactivationTime: deactivationTime,
					},
				)
			}
		}

		for _, def := range ruleInstall.RuleDefinitions {
			dynamicRulesToInstall = append(
				dynamicRulesToInstall,
				&protos.DynamicRuleInstall{
					PolicyRule:       def.ToProto(),
					ActivationTime:   activationTime,
					DeactivationTime: deactivationTime,
				},
			)
		}
	}
	return staticRulesToInstall, dynamicRulesToInstall
}

func ParseRuleRemoveAVPs(policyDBClient policydb.PolicyDBClient, rulesToRemoveAVP []*RuleRemoveAVP) []string {
	var ruleNames []string
	for _, rule := range rulesToRemoveAVP {
		ruleNames = append(ruleNames, rule.RuleNames...)
		if len(rule.RuleBaseNames) > 0 {
			ruleNames = append(ruleNames, policyDBClient.GetRuleIDsForBaseNames(rule.RuleBaseNames)...)
		}
	}
	return ruleNames
}

func GetEventTriggersRelatedInfo(
	eventTriggers []EventTrigger,
	revalidationTime *time.Time,
) ([]protos.EventTrigger, *timestamp.Timestamp) {
	protoEventTriggers := make([]protos.EventTrigger, 0, len(eventTriggers))
	var protoRevalidationTime *timestamp.Timestamp
	for _, eventTrigger := range eventTriggers {
		switch eventTrigger {
		case RevalidationTimeout:
			protoRevalidationTime = ConvertToProtoTimestamp(revalidationTime)
			protoEventTriggers = append(protoEventTriggers, protos.EventTrigger(eventTrigger))
		default:
			protoEventTriggers = append(protoEventTriggers, protos.EventTrigger_UNSUPPORTED)
		}
	}
	return protoEventTriggers, protoRevalidationTime
}

func getUsageMonitoringCredits(usageMonitors []*UsageMonitoringInfo) []*protos.UsageMonitoringCredit {
	usageMonitoringCredits := make([]*protos.UsageMonitoringCredit, 0, len(usageMonitors))
	for _, monitor := range usageMonitors {
		usageMonitoringCredits = append(
			usageMonitoringCredits,
			monitor.ToUsageMonitoringCredit(),
		)
	}
	return usageMonitoringCredits
}

func getQoSInfo(qosInfo *QosInformation) *protos.QoSInformation {
	if qosInfo == nil {
		return nil
	}
	res := &protos.QoSInformation{
		BearerId: qosInfo.BearerIdentifier,
	}
	if qosInfo.Qci != nil {
		res.Qci = protos.QCI(*qosInfo.Qci)
	}
	return res
}

func (report *UsageReport) FromUsageMonitorUpdate(update *protos.UsageMonitorUpdate) *UsageReport {
	if update == nil {
		return report
	}
	report.MonitoringKey = update.MonitoringKey
	report.Level = MonitoringLevel(update.Level)
	report.InputOctets = update.BytesTx
	report.OutputOctets = update.BytesRx // receive == output
	report.TotalOctets = update.BytesTx + update.BytesRx
	return report
}

func (monitor *UsageMonitoringInfo) ToUsageMonitoringCredit() *protos.UsageMonitoringCredit {
	return &protos.UsageMonitoringCredit{
		Action:        monitor.ToUsageMonitoringAction(),
		MonitoringKey: monitor.MonitoringKey,
		GrantedUnits:  monitor.GrantedServiceUnit.ToProto(),
		Level:         protos.MonitoringLevel(monitor.Level),
	}
}

// 3GPP TS 29.212
func (monitor *UsageMonitoringInfo) ToUsageMonitoringAction() protos.UsageMonitoringCredit_Action {
	if monitor.Report != nil && *monitor.Report == 0x0 {
		// 4.5.17.5 PCRF Requested Usage Report
		// `AVP: Usage-Monitoring-Report`
		return protos.UsageMonitoringCredit_FORCE
	}
	if monitor.Support != nil && *monitor.Support == 0x0 {
		// 4.5.17.3 Usage Monitoring Disabled
		// `AVP: Usage-Monitoring-Support`
		return protos.UsageMonitoringCredit_DISABLE
	}
	return protos.UsageMonitoringCredit_CONTINUE
}

func GetRATType(pRATType protos.RATType) credit_control.RATType {
	switch pRATType {
	case protos.RATType_TGPP_LTE:
		return credit_control.RAT_EUTRAN
	case protos.RATType_TGPP_WLAN:
		return credit_control.RAT_WLAN
	default:
		return credit_control.RAT_EUTRAN
	}
}

// Since we don't specify the IP CAN type at session initialization, and we
// only support WLAN and EUTRAN, we will infer the IP CAN type from RAT type.
func GetIPCANType(pRATType protos.RATType) credit_control.IPCANType {
	switch pRATType {
	case protos.RATType_TGPP_LTE:
		return credit_control.IPCAN_3GPP
	case protos.RATType_TGPP_WLAN:
		return credit_control.IPCAN_Non3GPP
	default:
		return credit_control.IPCAN_Non3GPP
	}
}
