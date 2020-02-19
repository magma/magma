/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package integ_tests

import (
	"fmt"

	fegProtos "magma/feg/cloud/go/protos"
	lteProtos "magma/lte/cloud/go/protos"
)

// RuleManager keeps track of rules and monitors for the integration
// test. It keeps track of all successfully installed static/dynamic rules
// along with usage monitors. The dynamic rules and usage monitors are added
// into the mock PCRF service. The static rules are usually streamed down into
// gateway policyDB from the cloud. Since this integration test does not cover
// the cloud component, we will directly insert the static policies into the
// redis database by using policyDBWrapper.
type RuleManager struct {
	// List of static rules successfully inserted into the policyDB store
	staticRules []*lteProtos.PolicyRule
	// List of dynamic rules successfully installed into PCRF
	dynamicRules []*fegProtos.AccountRules
	// List of usage monitors successfully installed into PCRF
	monitors []*fegProtos.UsageMonitorInfo
	// List of network wide static rules successfully inserted into the policyDB store
	omniPresentRules []*lteProtos.AssignedPolicies
	// Wrapper around redis operations for policyDB objects
	policyDBWrapper *policyDBWrapper
}

// NewRuleManager initialized the struct
func NewRuleManager() (*RuleManager, error) {
	policyDBWrapper, err := initializePolicyDBWrapper()
	if err != nil {
		return nil, err
	}
	return &RuleManager{
		policyDBWrapper: policyDBWrapper,
	}, nil
}

// AddStaticPassAll adds a static rule that passes all traffic to policyDB
// storage
func (manager *RuleManager) AddStaticPassAll(ruleID string, monitoringKey string, trackingType string, priority uint32) error {
	fmt.Printf("************************* Adding a Pass-All static rule: %s\n", ruleID)
	staticPassAll := getStaticPassAll(ruleID, monitoringKey, trackingType, priority)
	return manager.insertStaticRuleIntoRedis(staticPassAll)
}

// AddStaticRule adds the static rule to policyDB storage
func (manager *RuleManager) AddStaticRule(rule *lteProtos.PolicyRule) error {
	fmt.Printf("************************* Adding a static rule: %s\n", rule.Id)
	return manager.insertStaticRuleIntoRedis(rule)
}

// AddDynamicPassAll adds a dynamic rule that passes all traffic into PCRF
func (manager *RuleManager) AddDynamicPassAll(imsi, ruleID, monitoringKey string) error {
	fmt.Printf("************************* Adding Pass-All Dynamic Rule for UE with IMSI: %s, ruleID: %s\n", imsi, ruleID)
	dynamicPassAll := getDynamicPassAll(imsi, ruleID, monitoringKey)
	return manager.addDynamicRules(dynamicPassAll)
}

// AddDynamicRules adds the dynamic rule into PCRF
func (manager *RuleManager) AddDynamicRules(imsi string, ruleNames, baseNames []string) error {
	fmt.Printf("************************* Adding PCRF Rule for UE with IMSI: %s"+
		" with ruleNames=%v, baseNames=%v\n", imsi, ruleNames, baseNames)
	rules := makeDynamicRules(imsi, ruleNames, baseNames)
	return manager.addDynamicRules(rules)
}

// AddOmniPresentRules adds the network wide static rule to policyDB storage
func (manager *RuleManager) AddOmniPresentRules(keyId string, ruleNames, baseNames []string) error {
	fmt.Printf("************************* Adding a network wide rule\n")
	rule := makeAssignedRules(ruleNames, baseNames)
	return manager.insertOmniPresentRuleIntoRedis(keyId, rule)
}

// GetInstalledRulesByIMSI returns all dynamic rule ids and static rules
// referenced by dynamic rules keyed by the IMSI they are attached to.
func (manager *RuleManager) GetInstalledRulesByIMSI() map[string][]string {
	installedRulesByIMSI := map[string][]string{}
	for _, dynamicRules := range manager.dynamicRules {
		rules, exists := installedRulesByIMSI[dynamicRules.Imsi]
		if !exists {
			rules = []string{}
		}
		for _, ruleID := range dynamicRules.RuleNames {
			rules = append(rules, ruleID)
		}
		for _, dynamicRule := range dynamicRules.RuleDefinitions {
			rules = append(rules, dynamicRule.RuleName)
		}
		installedRulesByIMSI[dynamicRules.Imsi] = rules
	}
	return installedRulesByIMSI
}

// RemoveInstalledRules removes previously installed rules from PCRF and policyDB
func (manager *RuleManager) RemoveInstalledRules() error {
	rulesIDsByIMSI := manager.GetInstalledRulesByIMSI()
	for imsi, ruleIDs := range rulesIDsByIMSI {
		err := deactivateSubscriberFlows(imsi, ruleIDs)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddUsageMonitor constructs a usage monitor according to the parameters and
// inserts it into PCRF
func (manager *RuleManager) AddUsageMonitor(imsi, monitoringKey string, volume, bytesPerGrant uint64) error {
	fmt.Printf("************************* Adding PCRF Usage Monitor for UE with IMSI: %s\n", imsi)
	usageMonitor := makeUsageMonitor(imsi, monitoringKey, volume, bytesPerGrant)
	manager.monitors = append(manager.monitors, usageMonitor)
	return addPCRFUsageMonitors(usageMonitor)
}

func (manager *RuleManager) insertStaticRuleIntoRedis(rule *lteProtos.PolicyRule) error {
	err := manager.policyDBWrapper.policyMap.Set(rule.Id, rule)
	if err != nil {
		manager.staticRules = append(manager.staticRules, rule)
	}
	return err
}

func (manager *RuleManager) insertOmniPresentRuleIntoRedis(keyID string, rule *lteProtos.AssignedPolicies) error {
	err := manager.policyDBWrapper.omniPresentRules.Set(keyID, rule)
	if err != nil {
		manager.omniPresentRules = append(manager.omniPresentRules, rule)
	}
	return err
}

func (manager *RuleManager) addDynamicRules(rules *fegProtos.AccountRules) error {
	err := addPCRFRules(rules)
	if err != nil {
		manager.dynamicRules = append(manager.dynamicRules, rules)
	}
	return err
}

func getDynamicPassAll(imsi, ruleID, monitoringKey string) *fegProtos.AccountRules {
	return &fegProtos.AccountRules{
		Imsi:          imsi,
		RuleNames:     []string{},
		RuleBaseNames: []string{},
		RuleDefinitions: []*fegProtos.RuleDefinition{
			{
				RuleName:         ruleID,
				Precedence:       100,
				FlowDescriptions: []string{"permit out ip from any to any", "permit in ip from any to any"},
				MonitoringKey:    monitoringKey,
			},
		},
	}
}

func makeDynamicRules(imsi string, ruleNames []string, baseNames []string) *fegProtos.AccountRules {
	return &fegProtos.AccountRules{
		Imsi:          imsi,
		RuleNames:     ruleNames,
		RuleBaseNames: baseNames,
	}
}

func makeAssignedRules(ruleNames []string, baseNames []string) *lteProtos.AssignedPolicies {
	return &lteProtos.AssignedPolicies{
		AssignedPolicies:  ruleNames,
		AssignedBaseNames: baseNames,
	}
}

func makeUsageMonitor(imsi, monitoringKey string, volume, bytesPerGrant uint64) *fegProtos.UsageMonitorInfo {
	return &fegProtos.UsageMonitorInfo{
		Imsi: imsi,
		UsageMonitorCredits: []*fegProtos.UsageMonitorCredit{
			{
				MonitoringKey:   monitoringKey,
				Volume:          volume,
				ReturnBytes:     bytesPerGrant,
				MonitoringLevel: fegProtos.UsageMonitorCredit_RuleLevel,
			},
		},
	}
}
