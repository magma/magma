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

package mock_pcrf

import (
	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	lteprotos "magma/lte/cloud/go/protos"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func toStaticRuleNameRemovalAVP(ruleName string) *diam.AVP {
	return diam.NewAVP(avp.ChargingRuleRemove, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.ChargingRuleName, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(ruleName)),
		},
	})
}

func toStaticBaseNameRemovalAVP(baseName string) *diam.AVP {
	return diam.NewAVP(avp.ChargingRuleRemove, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.ChargingRuleBaseName, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(baseName)),
		},
	})
}

func toStaticRuleNameInstallAVP(ruleName string, activationTime, deactivationTime *timestamp.Timestamp) *diam.AVP {
	ruleInstallAVPs := []*diam.AVP{
		diam.NewAVP(avp.ChargingRuleName, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(ruleName)),
	}
	return diam.NewAVP(avp.ChargingRuleInstall, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: appendRuleInstallTimestamps(ruleInstallAVPs, activationTime, deactivationTime),
	})
}

func toStaticBaseNameInstallAVP(baseName string, activationTime, deactivationTime *timestamp.Timestamp) *diam.AVP {
	ruleInstallAVPs := []*diam.AVP{
		diam.NewAVP(avp.ChargingRuleBaseName, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(baseName)),
	}
	return diam.NewAVP(avp.ChargingRuleInstall, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: appendRuleInstallTimestamps(ruleInstallAVPs, activationTime, deactivationTime),
	})
}

func appendRuleInstallTimestamps(ruleInstallAVPs []*diam.AVP, activationTime, deactivationTime *timestamp.Timestamp) []*diam.AVP {
	if aTime, err := ptypes.Timestamp(activationTime); activationTime != nil && err == nil {
		ruleInstallAVPs = append(ruleInstallAVPs,
			diam.NewAVP(avp.RuleActivationTime, avp.Mbit, diameter.Vendor3GPP, datatype.Time(aTime)),
		)
	}
	if dTime, err := ptypes.Timestamp(deactivationTime); deactivationTime != nil && err == nil {
		ruleInstallAVPs = append(ruleInstallAVPs,
			diam.NewAVP(avp.RuleDeactivationTime, avp.Mbit, diameter.Vendor3GPP, datatype.Time(dTime)),
		)
	}
	return ruleInstallAVPs
}

func toDynamicRuleInstallAVP(rule *protos.RuleDefinition, activationTime *timestamp.Timestamp, deactivationTime *timestamp.Timestamp) *diam.AVP {
	ruleDefinitionAVPs := []*diam.AVP{
		diam.NewAVP(avp.ChargingRuleName, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(rule.RuleName)),
		diam.NewAVP(avp.Precedence, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(rule.Precedence)),
	}
	if rule.RatingGroup != 0 {
		ruleDefinitionAVPs = append(ruleDefinitionAVPs, diam.NewAVP(avp.RatingGroup, avp.Mbit, 0, datatype.Unsigned32(rule.RatingGroup)))
	}
	if rule.MonitoringKey != "" {
		ruleDefinitionAVPs = append(
			ruleDefinitionAVPs,
			diam.NewAVP(avp.MonitoringKey, avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(rule.MonitoringKey)),
		)
	}
	if rule.QosInformation != nil && (rule.QosInformation.MaxReqBwUl != 0 || rule.QosInformation.MaxReqBwDl != 0) {
		ruleDefinitionAVPs = append(ruleDefinitionAVPs, toQosAVP(rule.QosInformation))
	}
	for _, flowDescription := range rule.FlowDescriptions {
		ruleDefinitionAVPs = append(
			ruleDefinitionAVPs,
			diam.NewAVP(avp.FlowDescription, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.IPFilterRule(flowDescription)),
		)
	}
	if rule.RedirectInformation != nil {
		ruleDefinitionAVPs = append(
			ruleDefinitionAVPs,
			diam.NewAVP(avp.RedirectInformation, avp.Mbit, diameter.Vendor3GPP, &diam.GroupedAVP{
				AVP: toRedirectionAVP(rule.RedirectInformation),
			}),
		)
	}

	ruleInstallAVPs := []*diam.AVP{
		diam.NewAVP(avp.ChargingRuleDefinition, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
			AVP: ruleDefinitionAVPs,
		}),
	}
	return diam.NewAVP(avp.ChargingRuleInstall, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: appendRuleInstallTimestamps(ruleInstallAVPs, activationTime, deactivationTime),
	})
}

func toQosAVP(qos *lteprotos.FlowQos) *diam.AVP {
	qosAVPs := []*diam.AVP{}
	if qos.MaxReqBwUl != 0 {
		qosAVPs = append(
			qosAVPs,
			diam.NewAVP(avp.MaxRequestedBandwidthUL, avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(qos.MaxReqBwUl)),
		)
	}
	if qos.MaxReqBwDl != 0 {
		qosAVPs = append(
			qosAVPs,
			diam.NewAVP(avp.MaxRequestedBandwidthDL, avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(qos.MaxReqBwDl)),
		)
	}
	return diam.NewAVP(avp.QoSInformation, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{AVP: qosAVPs})
}

func toRedirectionAVP(redirection *lteprotos.RedirectInformation) []*diam.AVP {
	return []*diam.AVP{
		diam.NewAVP(avp.RedirectSupport, avp.Mbit, diameter.Vendor3GPP, datatype.Enumerated(redirection.Support)),
		diam.NewAVP(avp.RedirectAddressType, avp.Mbit, 0, datatype.Enumerated(redirection.AddressType)),
		diam.NewAVP(avp.RedirectServerAddress, avp.Mbit, 0, datatype.UTF8String(redirection.ServerAddress)),
	}
}

func toUsageMonitoringInfoAVP(monitoringKey string, quotaGrant *protos.Octets, level protos.MonitoringLevel) *diam.AVP {
	return diam.NewAVP(avp.UsageMonitoringInformation, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.MonitoringKey, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(monitoringKey)),
			diam.NewAVP(avp.GrantedServiceUnit, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: toGrantedServiceUnitAVP(quotaGrant),
			}),
			diam.NewAVP(avp.UsageMonitoringLevel, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(level)),
		},
	})
}

func toGrantedServiceUnitAVP(quotaGrant *protos.Octets) []*diam.AVP {
	res := []*diam.AVP{}
	if quotaGrant.GetTotalOctets() != 0 {
		res = append(res, diam.NewAVP(avp.CCTotalOctets, avp.Mbit, 0, datatype.Unsigned64(quotaGrant.GetTotalOctets())))
	}
	if quotaGrant.GetInputOctets() != 0 {
		res = append(res, diam.NewAVP(avp.CCInputOctets, avp.Mbit, 0, datatype.Unsigned64(quotaGrant.GetInputOctets())))
	}
	if quotaGrant.GetOutputOctets() != 0 {
		res = append(res, diam.NewAVP(avp.CCOutputOctets, avp.Mbit, 0, datatype.Unsigned64(quotaGrant.GetOutputOctets())))
	}
	return res
}

func toRuleInstallAVPs(ruleNames, ruleBaseNames []string, ruleDefs []*protos.RuleDefinition, activationTime, deactivationTime *timestamp.Timestamp) []*diam.AVP {
	avps := make([]*diam.AVP, 0, len(ruleNames)+len(ruleBaseNames)+len(ruleDefs))
	for _, ruleName := range ruleNames {
		avps = append(avps, toStaticRuleNameInstallAVP(ruleName, activationTime, deactivationTime))
	}

	for _, baseName := range ruleBaseNames {
		avps = append(avps, toStaticBaseNameInstallAVP(baseName, activationTime, deactivationTime))
	}

	for _, rule := range ruleDefs {
		avps = append(avps, toDynamicRuleInstallAVP(rule, activationTime, deactivationTime))
	}
	return avps
}

func toUsageMonitorAVPs(monitors map[string]*protos.UsageMonitor) []*diam.AVP {
	avps := make([]*diam.AVP, 0, len(monitors))
	for key, monitor := range monitors {
		avps = append(avps,
			toUsageMonitoringInfoAVP(key, getQuotaGrant(monitor), monitor.GetMonitorInfoPerRequest().GetMonitoringLevel()))
	}
	return avps
}

func toRuleRemovalAVPs(ruleNames, ruleBaseNames []string) []*diam.AVP {
	avps := make([]*diam.AVP, 0, len(ruleNames)+len(ruleBaseNames))
	for _, ruleName := range ruleNames {
		avps = append(avps, toStaticRuleNameRemovalAVP(ruleName))
	}

	for _, baseName := range ruleBaseNames {
		avps = append(avps, toStaticBaseNameRemovalAVP(baseName))
	}
	return avps
}

func toUsageMonitorByMkey(monitors []*usageMonitorRequestAVP) map[string]*usageMonitorRequestAVP {
	monitorByKey := map[string]*usageMonitorRequestAVP{}
	for _, monitor := range monitors {
		monitorByKey[monitor.MonitoringKey] = monitor
	}
	return monitorByKey
}

func toRevalidationTimeAVPs(revalidationTime *timestamp.Timestamp) *diam.AVP {
	if rTime, err := ptypes.Timestamp(revalidationTime); revalidationTime != nil && err == nil {
		return diam.NewAVP(avp.RevalidationTime, avp.Mbit, diameter.Vendor3GPP, datatype.Time(rTime))
	}
	return nil
}

func toEventTriggersAVPs(eventTriggers []uint32) []*diam.AVP {
	avps := make([]*diam.AVP, 0, len(eventTriggers))
	for _, event := range eventTriggers {
		avps = append(avps, diam.NewAVP(avp.EventTrigger, avp.Mbit, diameter.Vendor3GPP, datatype.Unsigned32(event)))
	}
	return avps
}
