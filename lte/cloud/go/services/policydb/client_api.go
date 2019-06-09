/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package client provides a thin client for contacting the policydb service.
// This can be used by apps to discover and contact the service, without knowing about
// the RPC implementation.
package policydb

import (
	"fmt"

	"magma/lte/cloud/go/protos"
	orcprotos "magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/registry"

	"github.com/golang/glog"
	"golang.org/x/net/context"
)

const EntityType = "policy"

const ServiceName = "POLICYDB"

// Utility function to get a RPC connection to the policydb service
func getPolicydbClient() (
	protos.PolicyDBControllerClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		glog.Errorf("Policydb client initialization error: %s", err)
		return nil, fmt.Errorf(
			"Policydb client initialization error: %s", err)
	}
	return protos.NewPolicyDBControllerClient(conn), err
}

// AddRule add a new rule.
// The rule must not be existing already.
func AddRule(networkId string, rule *protos.PolicyRule) error {
	ruleData := &protos.PolicyRuleData{
		NetworkId: &orcprotos.NetworkID{Id: networkId},
		Rule:      rule,
	}
	client, err := getPolicydbClient()
	if err != nil {
		return err
	}

	if _, err = client.AddRule(context.Background(), ruleData); err != nil {
		glog.Errorf("[Network: %s] AddRule error: %s", networkId, err)
		return err
	}
	return nil
}

// GetRule get the rule data.
func GetRule(networkId string, ruleId string) (
	*protos.PolicyRule, error) {
	client, err := getPolicydbClient()
	if err != nil {
		return nil, err
	}

	lookup := &protos.PolicyRuleLookup{
		NetworkId: &orcprotos.NetworkID{Id: networkId},
		RuleId:    ruleId}
	data, err := client.GetRule(context.Background(), lookup)
	if err != nil {
		glog.Errorf("[Network: %s, Sub: %s] GetRule error: %s",
			networkId, ruleId, err)
		return nil, err
	}
	return data, nil
}

// DeleteRule delete the rule.
func DeleteRule(networkId string, ruleId string) error {
	client, err := getPolicydbClient()
	if err != nil {
		return err
	}

	lookup := &protos.PolicyRuleLookup{
		NetworkId: &orcprotos.NetworkID{Id: networkId},
		RuleId:    ruleId}
	if _, err := client.DeleteRule(context.Background(), lookup); err != nil {
		glog.Errorf("[Network: %s, Sub: %s] DeleteSubscribererror: %s",
			networkId, ruleId, err)
		return err
	}
	return nil
}

// UpdateRule update the policy rule.
func UpdateRule(networkId string, rule *protos.PolicyRule) error {
	ruleData := &protos.PolicyRuleData{
		NetworkId: &orcprotos.NetworkID{Id: networkId},
		Rule:      rule,
	}
	client, err := getPolicydbClient()
	if err != nil {
		return err
	}

	if _, err = client.UpdateRule(context.Background(), ruleData); err != nil {
		glog.Errorf("[Network: %s] UpdateRule error: %s", networkId, err)
		return err
	}
	return nil
}

// GetAllRules get all policy rule objects
func GetAllRules(networkId string) ([]*protos.PolicyRule, error) {
	client, err := getPolicydbClient()
	if err != nil {
		return nil, err
	}

	ruleSet, err := client.ListRules(
		context.Background(),
		&orcprotos.NetworkID{Id: networkId})
	if err != nil {
		glog.Errorf("ListSubscribers error: %s", err)
		return nil, err
	}
	return ruleSet.GetRules(), nil
}

// ListRuleIds get all policy rule objects
func ListRuleIds(networkId string) ([]string, error) {
	rules, err := GetAllRules(networkId)
	if err != nil {
		return nil, err
	}
	ruleIds := make([]string, 0, len(rules))
	for _, rule := range rules {
		ruleIds = append(ruleIds, rule.Id)
	}
	return ruleIds, nil
}

//
// Base Name API
//
// AddBaseName adds new Charging Rule Base Name Record or Updates
// an existing RecordCorresponding to the given network & base name
// Returns the the existing base name record if present
func AddBaseName(networkId, baseName string, ruleNames []string) ([]string, error) {
	client, err := getPolicydbClient()
	if err != nil {
		return nil, err
	}

	old, err := client.AddBaseName(
		context.Background(),
		&protos.ChargingRuleBaseNameRequest{
			Lookup: &protos.ChargingRuleBaseNameLookup{
				NetworkID: &orcprotos.NetworkID{Id: networkId},
				Name:      baseName,
			},
			Record: &protos.ChargingRuleNameSet{RuleNames: ruleNames},
		})
	if err != nil {
		return nil, err
	}
	if len(old.RuleNames) == 0 {
		return nil, nil
	}
	return old.RuleNames, nil
}

// DeleteBaseName deletes an existing Charging Rule Base Name Record
func DeleteBaseName(networkId, baseName string) error {
	client, err := getPolicydbClient()
	if err == nil {
		_, err = client.DeleteBaseName(
			context.Background(),
			&protos.ChargingRuleBaseNameLookup{
				NetworkID: &orcprotos.NetworkID{Id: networkId},
				Name:      baseName})
	}
	return err
}

// GetBaseName returns the Charging Rule Name List corresponding to the given base name on the given network
func GetBaseName(networkId, baseName string) ([]string, error) {
	client, err := getPolicydbClient()
	if err != nil {
		return nil, err
	}
	ruleNamesSet, err := client.GetBaseName(
		context.Background(),
		&protos.ChargingRuleBaseNameLookup{
			NetworkID: &orcprotos.NetworkID{Id: networkId},
			Name:      baseName,
		})
	if err != nil {
		return nil, err
	}
	return ruleNamesSet.RuleNames, err
}

// ListBaseNames returns a list of all Base Names for the network, the Rule Name Lists
// associated with each base name can be retrieved using separate GetBaseName() call
func ListBaseNames(networkId string) ([]string, error) {
	client, err := getPolicydbClient()
	if err != nil {
		return nil, err
	}
	baseNamesSet, err := client.ListBaseNames(context.Background(), &orcprotos.NetworkID{Id: networkId})
	return baseNamesSet.RuleNames, err
}

// GetAllBaseNames returns a list of all Base Names Records for the network
func GetAllBaseNames(networkId string) ([]*protos.ChargingRuleBaseNameRecord, error) {
	client, err := getPolicydbClient()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	baseNamesSet, err := client.ListBaseNames(ctx, &orcprotos.NetworkID{Id: networkId})
	if err != nil {
		return nil, err
	}

	res := make([]*protos.ChargingRuleBaseNameRecord, len(baseNamesSet.RuleNames))
	for i, baseName := range baseNamesSet.RuleNames {
		ruleNamesSet, err := client.GetBaseName(
			ctx,
			&protos.ChargingRuleBaseNameLookup{
				NetworkID: &orcprotos.NetworkID{Id: networkId},
				Name:      baseName,
			})
		if err != nil {
			return nil, err
		}
		res[i] = &protos.ChargingRuleBaseNameRecord{
			Name:         baseName,
			RuleNamesSet: &protos.ChargingRuleNameSet{RuleNames: ruleNamesSet.GetRuleNames()},
		}
	}

	return res, err
}
