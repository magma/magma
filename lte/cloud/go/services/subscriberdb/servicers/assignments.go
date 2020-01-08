/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

/*
policyd assignments servicer provides the gRPC interface for the REST and
services to interact with assignments from policy rules and subscribers.
*/
package servicers

import (
	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/protos"
	orcprotos "magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PolicyAssignmentServer struct{}

func NewPolicyAssignmentServer() *PolicyAssignmentServer {
	return &PolicyAssignmentServer{}
}

func (srv *PolicyAssignmentServer) EnableStaticRules(ctx context.Context, req *protos.EnableStaticRuleRequest) (*orcprotos.Void, error) {
	networkID, err := getNetworkID(ctx)
	if !doesSubscriberAndRulesExist(networkID, req.Imsi, req.RuleIds, req.BaseNames) {
		return nil, status.Errorf(codes.InvalidArgument, "Either a subscriber or one more rules/basenames are not found")
	}
	updates := []configurator.EntityUpdateCriteria{}
	for _, ruleID := range req.RuleIds {
		updates = append(updates, getRuleUpdateForEnable(ruleID, req.Imsi))
	}
	for _, baseName := range req.BaseNames {
		updates = append(updates, getBaseNameUpdateForEnable(baseName, req.Imsi))
	}
	_, err = configurator.UpdateEntities(networkID, updates)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "Failed to enable")
	}
	return nil, nil
}

func (srv *PolicyAssignmentServer) DisableStaticRules(ctx context.Context, req *protos.DisableStaticRuleRequest) (*orcprotos.Void, error) {
	networkID, err := getNetworkID(ctx)
	if !doesSubscriberAndRulesExist(networkID, req.Imsi, req.RuleIds, req.BaseNames) {
		return nil, status.Errorf(codes.InvalidArgument, "Either a subscriber or one more rules/basenames are not found")
	}
	updates := []configurator.EntityUpdateCriteria{}
	for _, ruleID := range req.RuleIds {
		updates = append(updates, getRuleUpdateForDisableRule(ruleID, req.Imsi))
	}
	for _, baseName := range req.BaseNames {
		updates = append(updates, getBaseNameUpdateForDisable(baseName, req.Imsi))
	}
	_, err = configurator.UpdateEntities(networkID, updates)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "Failed to disable")
	}
	return nil, nil
}

func doesSubscriberAndRulesExist(networkID string, subscriberID string, ruleIDs []string, baseNames []string) bool {
	ids := []storage.TypeAndKey{storage.TypeAndKey{Type: lte.SubscriberEntityType, Key: subscriberID}}
	for _, ruleID := range ruleIDs {
		ids = append(ids, storage.TypeAndKey{Type: lte.PolicyRuleEntityType, Key: ruleID})
	}
	for _, baseName := range baseNames {
		ids = append(ids, storage.TypeAndKey{Type: lte.BaseNameEntityType, Key: baseName})
	}
	exists, err := configurator.DoEntitiesExist(networkID, ids)
	if err != nil {
		return false
	}
	return exists
}

func getRuleUpdateForEnable(ruleID string, subscriberID string) configurator.EntityUpdateCriteria {
	ret := configurator.EntityUpdateCriteria{
		Type: lte.PolicyRuleEntityType,
		Key:  ruleID,
	}
	ret.AssociationsToAdd = append(ret.AssociationsToAdd, storage.TypeAndKey{Type: lte.SubscriberEntityType, Key: subscriberID})
	return ret
}

func getRuleUpdateForDisableRule(ruleID string, subscriberID string) configurator.EntityUpdateCriteria {
	ret := configurator.EntityUpdateCriteria{
		Type: lte.PolicyRuleEntityType,
		Key:  ruleID,
	}
	ret.AssociationsToDelete = append(ret.AssociationsToDelete, storage.TypeAndKey{Type: lte.SubscriberEntityType, Key: subscriberID})
	return ret
}

func getBaseNameUpdateForEnable(baseName string, subscriberID string) configurator.EntityUpdateCriteria {
	ret := configurator.EntityUpdateCriteria{
		Type: lte.BaseNameEntityType,
		Key:  baseName,
	}
	ret.AssociationsToAdd = append(ret.AssociationsToAdd, storage.TypeAndKey{Type: lte.SubscriberEntityType, Key: subscriberID})
	return ret
}

func getBaseNameUpdateForDisable(baseName string, subscriberID string) configurator.EntityUpdateCriteria {
	ret := configurator.EntityUpdateCriteria{
		Type: lte.BaseNameEntityType,
		Key:  baseName,
	}
	ret.AssociationsToDelete = append(ret.AssociationsToDelete, storage.TypeAndKey{Type: lte.SubscriberEntityType, Key: subscriberID})
	return ret
}

func getNetworkID(ctx context.Context) (string, error) {
	id, err := orcprotos.GetGatewayIdentity(ctx)
	if err != nil {
		return "", err
	}
	return id.GetNetworkId(), nil
}
