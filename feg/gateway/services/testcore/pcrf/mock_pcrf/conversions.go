/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package mock_pcrf

import (
	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	lteprotos "magma/lte/cloud/go/protos"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
)

func toStaticRuleNameInstallAVP(ruleName string) *diam.AVP {
	return diam.NewAVP(avp.ChargingRuleInstall, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.ChargingRuleName, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(ruleName)),
		},
	})
}

func toStaticBaseNameInstallAVP(baseName string) *diam.AVP {
	return diam.NewAVP(avp.ChargingRuleInstall, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.ChargingRuleBaseName, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.UTF8String(baseName)),
		},
	})
}

func toDynamicRuleInstallAVP(rule *protos.RuleDefinition) *diam.AVP {
	installAVPs := []*diam.AVP{
		diam.NewAVP(avp.ChargingRuleName, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(rule.RuleName)),
		diam.NewAVP(avp.Precedence, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(rule.Precedence)),
	}
	if rule.RatingGroup != 0 {
		installAVPs = append(installAVPs, diam.NewAVP(avp.RatingGroup, avp.Mbit, 0, datatype.Unsigned32(rule.RatingGroup)))
	}
	if rule.MonitoringKey != "" {
		installAVPs = append(
			installAVPs,
			diam.NewAVP(avp.MonitoringKey, avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(rule.MonitoringKey)),
		)
	}
	if rule.QosInformation != nil && (rule.QosInformation.MaxReqBwUl != 0 || rule.QosInformation.MaxReqBwDl != 0) {
		installAVPs = append(installAVPs, toQosAVP(rule.QosInformation))
	}
	for _, flowDescription := range rule.FlowDescriptions {
		installAVPs = append(
			installAVPs,
			diam.NewAVP(avp.FlowDescription, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.IPFilterRule(flowDescription)),
		)
	}
	if rule.RedirectInformation != nil {
		installAVPs = append(
			installAVPs,
			diam.NewAVP(avp.RedirectInformation, avp.Mbit, diameter.Vendor3GPP, &diam.GroupedAVP{
				AVP: toRedirectionAVP(rule.RedirectInformation),
			}),
		)
	}
	return diam.NewAVP(avp.ChargingRuleInstall, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.ChargingRuleDefinition, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
				AVP: installAVPs,
			}),
		},
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

func toUsageMonitoringInfoAVP(monitoringKey string, quotaGrant protos.Octets, level protos.UsageMonitorCredit_MonitoringLevel) *diam.AVP {
	return diam.NewAVP(avp.UsageMonitoringInformation, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.MonitoringKey, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(monitoringKey)),
			diam.NewAVP(avp.GrantedServiceUnit, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.CCTotalOctets, avp.Mbit, 0, datatype.Unsigned64(quotaGrant.TotalOctets)),
					diam.NewAVP(avp.CCInputOctets, avp.Mbit, 0, datatype.Unsigned64(quotaGrant.InputOctets)),
					diam.NewAVP(avp.CCOutputOctets, avp.Mbit, 0, datatype.Unsigned64(quotaGrant.OutputOctets)),
				},
			}),
			diam.NewAVP(avp.UsageMonitoringLevel, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(level)),
		},
	})
}

func toRuleInstallAVPs(
	ruleNames []string,
	ruleBaseNames []string,
	ruleDefs []*protos.RuleDefinition,
) []*diam.AVP {
	avps := make([]*diam.AVP, 0, len(ruleNames)+len(ruleBaseNames)+len(ruleDefs))
	for _, ruleName := range ruleNames {
		avps = append(avps, toStaticRuleNameInstallAVP(ruleName))
	}

	for _, baseName := range ruleBaseNames {
		avps = append(avps, toStaticBaseNameInstallAVP(baseName))
	}

	for _, rule := range ruleDefs {
		avps = append(avps, toDynamicRuleInstallAVP(rule))
	}
	return avps
}

func toUsageMonitoringInfoAVPs(monitors map[string]*protos.UsageMonitorCredit) []*diam.AVP {
	avps := make([]*diam.AVP, 0, len(monitors))
	for key, monitor := range monitors {
		avps = append(avps, toUsageMonitoringInfoAVP(key, getQuotaGrant(monitor), monitor.MonitoringLevel))
	}
	return avps
}
