/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package models

import (
	"fmt"
	"reflect"

	"magma/lte/cloud/go/protos"
	orcprotos "magma/orc8r/cloud/go/protos"

	"github.com/go-openapi/strfmt"
	"github.com/golang/protobuf/proto"
)

var formatsRegistry = strfmt.NewFormats()

// PolicyRules's FromProto fills in models.PolicyRules struct from
// passed protos.PolicyRule
func (policyRule *PolicyRule) FromProto(pfrm proto.Message) error {
	flowRuleProto, ok := pfrm.(*protos.PolicyRule)
	if !ok {
		return fmt.Errorf(
			"Invalid Source Type %s, *protos.PolicyRule expected",
			reflect.TypeOf(pfrm))
	}
	if policyRule != nil {
		if flowRuleProto != nil {
			orcprotos.FillIn(flowRuleProto, policyRule)
			if flowRuleProto.FlowList != nil && policyRule.FlowList == nil {
				policyRule.FlowList = flowListFromProto(flowRuleProto.FlowList)
			}
			policyRule.Priority = &flowRuleProto.Priority
			trackingName, ok := protos.PolicyRule_TrackingType_name[int32(flowRuleProto.TrackingType)]
			if ok {
				policyRule.TrackingType = trackingName
			}
			if flowRuleProto.Redirect != nil {
				policyRule.Redirect = redirectInfoFromProto(flowRuleProto.Redirect)
			}
			policyRule.MonitoringKey = &flowRuleProto.MonitoringKey
			policyRule.RatingGroup = &flowRuleProto.RatingGroup
			return policyRule.Verify()
		}
	}
	return nil
}

// PolicyRule's ToProto fills in passed protos.PolicyRule struct from
// receiver's protos.PolicyRule
func (policyRule *PolicyRule) ToProto(pfrm proto.Message) error {
	flowRuleProto, ok := pfrm.(*protos.PolicyRule)
	if !ok {
		return fmt.Errorf(
			"Invalid Destination Type %s, *protos.PolicyRule expected",
			reflect.TypeOf(pfrm))
	}
	if policyRule != nil || flowRuleProto != nil {
		orcprotos.FillIn(policyRule, flowRuleProto)
		trackingVal, ok := protos.PolicyRule_TrackingType_value[policyRule.TrackingType]
		if ok {
			flowRuleProto.TrackingType = protos.PolicyRule_TrackingType(trackingVal)
		}
		if flowRuleProto.FlowList == nil {
			flowRuleProto.FlowList = flowListToProto(policyRule.FlowList)
		}
		if policyRule.Redirect != nil {
			flowRuleProto.Redirect = redirectInfoToProto(policyRule.Redirect)
		}
		if policyRule.Priority != nil {
			flowRuleProto.Priority = *policyRule.Priority
		}
		if policyRule.MonitoringKey != nil {
			flowRuleProto.MonitoringKey = *policyRule.MonitoringKey
		}
		if policyRule.RatingGroup != nil {
			flowRuleProto.RatingGroup = *policyRule.RatingGroup
		}
	}
	return nil
}

func redirectInfoFromProto(redirectProto *protos.RedirectInformation) *RedirectInformation {
	modelInfo := &RedirectInformation{}
	orcprotos.FillIn(redirectProto, modelInfo)
	supportName, ok := protos.RedirectInformation_Support_name[int32(redirectProto.Support)]
	if ok {
		modelInfo.Support = supportName
	}
	addrTypeName, ok := protos.RedirectInformation_AddressType_name[int32(redirectProto.AddressType)]
	if ok {
		modelInfo.AddressType = addrTypeName
	}
	return modelInfo
}

func redirectInfoToProto(redirectModel *RedirectInformation) *protos.RedirectInformation {
	redirectProto := &protos.RedirectInformation{}
	orcprotos.FillIn(redirectModel, redirectProto)
	supportVal, ok := protos.RedirectInformation_Support_value[redirectModel.Support]
	if ok {
		redirectProto.Support = protos.RedirectInformation_Support(supportVal)
	}
	addrTypeVal, ok := protos.RedirectInformation_AddressType_value[redirectModel.AddressType]
	if ok {
		redirectProto.AddressType = protos.RedirectInformation_AddressType(addrTypeVal)
	}
	return redirectProto
}

// Fill models.PolicyRules.FlowList From protos.PolicyRule.FlowList
func flowListFromProto(flowList []*protos.FlowDescription) []*FlowDescription {
	var s []*FlowDescription
	for i, flow := range flowList {
		s = append(s, new(FlowDescription))
		orcprotos.FillIn(flow, s[i])
		match := flow.Match
		orcprotos.FillIn(match, s[i].Match)
		protoName, ok := protos.FlowMatch_IPProto_name[int32(match.IpProto)]
		if ok {
			s[i].Match.IPProto = &protoName
		}
		directionName, ok := protos.FlowMatch_Direction_name[int32(match.Direction)]
		if ok {
			s[i].Match.Direction = directionName
		}
		actionName, ok := protos.FlowDescription_Action_name[int32(flow.Action)]
		if ok {
			s[i].Action = &actionName
		}
	}
	return s
}

// Fill protos.PolicyRule.FlowList From passed protos.PolicyRule.FlowList
func flowListToProto(flowList []*FlowDescription) []*protos.FlowDescription {
	var s []*protos.FlowDescription
	for i, flow := range flowList {
		s = append(s, new(protos.FlowDescription))
		orcprotos.FillIn(flow, s[i])
		match := flow.Match
		orcprotos.FillIn(match, s[i].Match)
		if match.IPProto != nil {
			protoVal, ok := protos.FlowMatch_IPProto_value[*match.IPProto]
			if ok {
				s[i].Match.IpProto = protos.FlowMatch_IPProto(protoVal)
			}
		}
		directionVal, ok := protos.FlowMatch_Direction_value[match.Direction]
		if ok {
			s[i].Match.Direction = protos.FlowMatch_Direction(directionVal)
		}
		if flow.Action != nil {
			actionVal, ok := protos.FlowDescription_Action_value[*flow.Action]
			if ok {
				s[i].Action = protos.FlowDescription_Action(actionVal)
			}
		}
	}
	return s
}

// Verify validates given PolicyRule
func (policyRule *PolicyRule) Verify() error {
	if policyRule == nil {
		return fmt.Errorf("Nil PolicyRule pointer")
	}
	err := policyRule.Validate(formatsRegistry)
	if policyRule.ID == "" {
		return fmt.Errorf("Missing PolicyRule ID")
	}
	if err != nil {
		return fmt.Errorf("Flow rule validation error: %s", err)
	}
	return nil
}
