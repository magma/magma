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

package integration

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
	staticRuleDefs []*lteProtos.PolicyRule
	// List of basename -> rule names mapping successfully inserted into the
	// policyDB store
	baseNameMappings []*lteProtos.ChargingRuleBaseNameRecord
	// List of dynamic rules successfully installed into PCRF
	accountRules []*fegProtos.AccountRules
	// List of usage monitors successfully installed into PCRF
	monitors []*fegProtos.UsageMonitorConfiguration
	// List of network wide static rules successfully inserted into the policyDB store
	omniPresentRules []*lteProtos.AssignedPolicies
	// Wrapper around redis operations for policyDB objects
	policyDBWrapper *policyDBWrapper
	// Instance name of the PCRF this rule manager is attached to
	pcrfInstance string
}

// NewRuleManager initialized the struct
func NewRuleManager() (*RuleManager, error) {
	return NewRuleManagerPerInstance(MockPCRFRemote)
}

// NewRuleManagerPerInstance initialized the struct per PCRFinstance
func NewRuleManagerPerInstance(pcrfInstance string) (*RuleManager, error) {
	policyDBWrapper, err := initializePolicyDBWrapper()
	if err != nil {
		return nil, err
	}
	return &RuleManager{
		policyDBWrapper: policyDBWrapper,
		pcrfInstance:    pcrfInstance,
	}, nil
}

// AddStaticPassAllToDBAndPCRF adds a static rule that passes all traffic to policyDB
// storage and to the PCRF instance
func (manager *RuleManager) AddStaticPassAllToDBAndPCRFforIMSIs(
	IMSIs []string, ruleID string, monitoringKey string, ratingGroup uint32, trackingType string, priority uint32,
) error {
	fmt.Printf("************* Adding a Pass-All static rule to DB and PCRF: %s\n", ruleID)
	staticPassAll := getStaticPassAll(ruleID, monitoringKey, ratingGroup, trackingType, priority, nil)

	err := manager.insertStaticRuleIntoRedis(staticPassAll)
	if err != nil {
		return err
	}
	for _, imsi := range IMSIs {
		err = manager.AddRulesToPCRF(imsi, []string{ruleID}, []string{})
		if err != nil {
			return err
		}
	}
	return nil
}

// AddStaticPassAllToDB adds a static rule that passes all traffic to policyDB
// storage
func (manager *RuleManager) AddStaticPassAllToDB(ruleID string, monitoringKey string, ratingGroup uint32, trackingType string, priority uint32) error {
	fmt.Printf("************* Adding a Pass-All static rule: %s, priority: %d, mkey: %s, rg: %d, trackingType: %s\n",
		ruleID, priority, monitoringKey, ratingGroup, trackingType)
	staticPassAll := getStaticPassAll(ruleID, monitoringKey, ratingGroup, trackingType, priority, nil)
	return manager.insertStaticRuleIntoRedis(staticPassAll)
}

// AddStaticRuleToDB adds the static rule to policyDB storage
func (manager *RuleManager) AddStaticRuleToDB(rule *lteProtos.PolicyRule) error {
	fmt.Printf("************* Adding a static rule: %s, priority: %d, mkey: %s, rg: %d, trackingType: %s\n",
		rule.Id, rule.Priority, rule.MonitoringKey, rule.RatingGroup, rule.TrackingType)
	return manager.insertStaticRuleIntoRedis(rule)
}

// AddDynamicPassAllToPCRF adds a dynamic rule that passes all traffic into PCRF
func (manager *RuleManager) AddDynamicPassAllToPCRF(imsi, ruleID, monitoringKey string) error {
	fmt.Printf("************* Adding Pass-All Dynamic Rule for UE with IMSI: %s, ruleID: %s\n", imsi, ruleID)
	dynamicPassAll := getAccountRulesWithDynamicPassAll(imsi, ruleID, monitoringKey)
	return manager.addAccountRules(dynamicPassAll)
}

// AddRulesToPCRF adds the dynamic rule into PCRF
func (manager *RuleManager) AddRulesToPCRF(imsi string, ruleNames, baseNames []string) error {
	fmt.Printf("************* Adding PCRF Rule for UE with IMSI: %s"+
		" with ruleNames=%v, baseNames=%v\n", imsi, ruleNames, baseNames)
	rules := makeAccountRules(imsi, ruleNames, baseNames)
	return manager.addAccountRules(rules)
}

func (manager *RuleManager) AddBaseNameMappingToDB(basename string, ruleNames []string) error {
	fmt.Printf("************* Adding a base name mapping of %s -> %v\n", basename, ruleNames)
	record := &lteProtos.ChargingRuleBaseNameRecord{
		Name:         basename,
		RuleNamesSet: &lteProtos.ChargingRuleNameSet{RuleNames: ruleNames},
	}
	return manager.insertBaseNameMappingIntoRedis(record)
}

// AddOmniPresentRulesToDB adds the network wide static rule to policyDB storage
func (manager *RuleManager) AddOmniPresentRulesToDB(keyId string, ruleNames, baseNames []string) error {
	fmt.Printf("************* Adding a network wide rule\n")
	rule := makeAssignedRules(ruleNames, baseNames)
	return manager.insertOmniPresentRuleIntoRedis(keyId, rule)
}

// RemoveOmniPresentRulesFromDB adds the network wide static rule to policyDB storage
func (manager *RuleManager) RemoveOmniPresentRulesFromDB(keyId string) error {
	fmt.Printf("************* Removing a network wide rule\n")
	return manager.removeOmniPresentRuleIntoRedis(keyId)
}

// GetInstalledRulesByIMSI returns all dynamic rule ids and static rules
// referenced by dynamic rules keyed by the IMSI they are attached to.
func (manager *RuleManager) GetInstalledRulesByIMSI() map[string][]string {
	installedRulesByIMSI := map[string][]string{}
	for _, accountRules := range manager.accountRules {
		rules, exists := installedRulesByIMSI[accountRules.Imsi]
		if !exists {
			rules = []string{}
		}
		rules = append(rules, accountRules.StaticRuleNames...)
		for _, dynamicRule := range accountRules.DynamicRuleDefinitions {
			rules = append(rules, dynamicRule.RuleName)
		}
		installedRulesByIMSI[accountRules.Imsi] = rules
	}
	return installedRulesByIMSI
}

// RemoveInstalledRules removes previously installed rules from PCRF and policyDB
func (manager *RuleManager) RemoveInstalledRules() error {
	rulesIDsByIMSI := manager.GetInstalledRulesByIMSI()
	for imsi := range rulesIDsByIMSI {
		err := deactivateAllFlowsPerSub(imsi)
		if err != nil {
			return err
		}
	}
	manager.policyDBWrapper.policyMap.DeleteAll()
	manager.policyDBWrapper.baseNameMap.DeleteAll()
	return nil
}

// AddUsageMonitor constructs a usage monitor according to the parameters and
// inserts it into PCRF
func (manager *RuleManager) AddUsageMonitor(imsi, monitoringKey string, volume, bytesPerGrant uint64) error {
	fmt.Printf("************* Adding PCRF Usage Monitor for UE with IMSI: %s\n", imsi)
	usageMonitor := makeUsageMonitor(imsi, monitoringKey, volume, bytesPerGrant)
	err := addPCRFUsageMonitorsPerInstance(manager.pcrfInstance, usageMonitor)
	if err != nil {
		return err
	}
	manager.monitors = append(manager.monitors, usageMonitor)
	return nil
}

func (manager *RuleManager) insertStaticRuleIntoRedis(rule *lteProtos.PolicyRule) error {
	err := manager.policyDBWrapper.policyMap.Set(rule.Id, rule)
	if err != nil {
		return err
	}
	manager.staticRuleDefs = append(manager.staticRuleDefs, rule)
	return nil
}

func (manager *RuleManager) insertBaseNameMappingIntoRedis(record *lteProtos.ChargingRuleBaseNameRecord) error {
	err := manager.policyDBWrapper.baseNameMap.Set(record.GetName(), record.GetRuleNamesSet())
	if err != nil {
		return err
	}
	manager.baseNameMappings = append(manager.baseNameMappings, record)
	return nil
}

func (manager *RuleManager) insertOmniPresentRuleIntoRedis(keyID string, rule *lteProtos.AssignedPolicies) error {
	err := manager.policyDBWrapper.omniPresentRules.Set(keyID, rule)
	if err != nil {
		return err
	}
	manager.omniPresentRules = append(manager.omniPresentRules, rule)
	return nil
}

func (manager *RuleManager) removeOmniPresentRuleIntoRedis(keyID string) error {
	err := manager.policyDBWrapper.omniPresentRules.Delete(keyID)
	if err != nil {
		return err
	}
	return nil
}

func (manager *RuleManager) addAccountRules(rules *fegProtos.AccountRules) error {
	err := addPCRFRulesPerInstance(manager.pcrfInstance, rules)
	if err != nil {
		return err
	}
	manager.accountRules = append(manager.accountRules, rules)
	return nil
}

func getAccountRulesWithDynamicPassAll(imsi, ruleID, monitoringKey string) *fegProtos.AccountRules {
	return &fegProtos.AccountRules{
		Imsi:                imsi,
		StaticRuleNames:     []string{},
		StaticRuleBaseNames: []string{},
		DynamicRuleDefinitions: []*fegProtos.RuleDefinition{
			getPassAllRuleDefinition(ruleID, monitoringKey, nil, 100),
		},
	}
}

func makeAccountRules(imsi string, ruleNames []string, baseNames []string) *fegProtos.AccountRules {
	return &fegProtos.AccountRules{
		Imsi:                imsi,
		StaticRuleNames:     ruleNames,
		StaticRuleBaseNames: baseNames,
	}
}

func makeAssignedRules(ruleNames []string, baseNames []string) *lteProtos.AssignedPolicies {
	return &lteProtos.AssignedPolicies{
		AssignedPolicies:  ruleNames,
		AssignedBaseNames: baseNames,
	}
}

func makeUsageMonitor(imsi, monitoringKey string, volume, bytesPerGrant uint64) *fegProtos.UsageMonitorConfiguration {
	return &fegProtos.UsageMonitorConfiguration{
		Imsi: imsi,
		UsageMonitorCredits: []*fegProtos.UsageMonitor{
			{
				MonitorInfoPerRequest: &fegProtos.UsageMonitoringInformation{
					MonitoringKey:   []byte(monitoringKey),
					MonitoringLevel: fegProtos.MonitoringLevel_RuleLevel,
					Octets:          &fegProtos.Octets{TotalOctets: bytesPerGrant},
				},
				TotalQuota: &fegProtos.Octets{TotalOctets: volume},
			},
		},
	}
}
