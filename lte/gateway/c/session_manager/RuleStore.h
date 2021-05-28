/**
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
#pragma once

#include <lte/protos/pipelined.grpc.pb.h>
#include <lte/protos/policydb.pb.h>

#include <functional>
#include <memory>
#include <mutex>
#include <string>
#include <unordered_map>
#include <vector>

#include "CreditKey.h"
#include "includes/GRPCReceiver.h"

using grpc::Status;

namespace magma {
using namespace lte;
/**
 * Template class for keeping track of a map of one key to many policy rules
 */
template<
    typename KeyType, typename hash = std::hash<KeyType>,
    typename equal = std::equal_to<KeyType>>
class PoliciesByKeyMap {
 public:
  PoliciesByKeyMap() {}
  PoliciesByKeyMap(hash hasher, equal eq) : rules_by_key_(4, hasher, eq) {}

  void insert(const KeyType& key, std::shared_ptr<PolicyRule> rule_p);

  void remove(const KeyType& key, std::shared_ptr<PolicyRule> rule_p);

  uint32_t policy_count();

  bool get_rule_ids_for_key(
      const KeyType& key, std::vector<std::string>& rules_out);

  bool get_rule_definitions_for_key(
      const KeyType& key, std::vector<PolicyRule>& rules_out);

 private:
  std::unordered_map<
      KeyType, std::vector<std::shared_ptr<PolicyRule>>, hash, equal>
      rules_by_key_;
};

/**
 * RuleChargingKeyMapper is a class for querying a bi-directional map of
 * rule_id <-> charging_key
 */
class PolicyRuleBiMap {
 public:
  PolicyRuleBiMap() : rules_by_charging_key_(&ccHash, &ccEqual) {}
  /**
   * Clear the maps and add in the given rules
   */
  virtual void sync_rules(const std::vector<PolicyRule>& rules);

  virtual void insert_rule(const PolicyRule& rule);

  // Get the rule definition associated with the given rule_id
  // If the rule is found, copy the rule into the output parameter and return
  // true. Otherwise, return false.
  // If the output rule param is NULL, the rule object is not copied.
  virtual bool get_rule(const std::string& rule_id, PolicyRule* rule_out);

  virtual bool get_rules_by_ids(
      const std::vector<std::string>& rule_ids,
      std::vector<PolicyRule>& rules_out);

  // Remove a rule from the store by ID. Returns true if the rule ID was found.
  // The removed rule will be copied into rule_out.
  // If the output rule param is NULL, the rule object is not copied.
  virtual bool remove_rule(const std::string& rule_id, PolicyRule* rule_out);

  /**
   * Get the charging key for a particular rule id. The charging key is set in
   * the out parameter charging_key
   * @returns false if it doesn't exist, true if so
   */
  virtual bool get_charging_key_for_rule_id(
      const std::string& rule_id, CreditKey* charging_key);

  virtual bool get_monitoring_key_for_rule_id(
      const std::string& rule_id, std::string* monitoring_key);

  /**
   * Get all the rules for a given key. Rule ids are copied into rules_out
   */
  virtual bool get_rule_ids_for_charging_key(
      const CreditKey& charging_key, std::vector<std::string>& rules_out);

  /**
   * Get all the rules for a given key. Rule ids are copied into rules_out
   */
  virtual bool get_rule_definitions_for_charging_key(
      const CreditKey& charging_key, std::vector<PolicyRule>& rules_out);

  /**
   * Get the number of rules tracked by a monitoring key
   */
  virtual uint32_t monitored_rules_count();

  virtual bool get_rule_ids(std::vector<std::string>& rules_ids_out);

  virtual bool get_rules(std::vector<PolicyRule>& rules_out);

 protected:
  // guards all three maps below
  std::mutex map_mutex_;
  // rule_id -> PolicyRule
  std::unordered_map<std::string, std::shared_ptr<PolicyRule>>
      rules_by_rule_id_;
  // charging key -> [PolicyRule]
  PoliciesByKeyMap<CreditKey, decltype(&ccHash), decltype(&ccEqual)>
      rules_by_charging_key_;
  // monitoring key -> [PolicyRule]
  PoliciesByKeyMap<std::string> rules_by_monitoring_key_;
};

/**
 * StaticRuleStore holds the rules that are defined in policydb
 */
class StaticRuleStore : public PolicyRuleBiMap {};

/**
 * DynamicRuleStore manages dynamic rules for a subscriber
 */
class DynamicRuleStore : public PolicyRuleBiMap {};

class ConvergedRuleStore : public PolicyRuleBiMap {
 public:
  uint32_t pdr_rule_count(void) { return rules_by_pdr_key_.size(); }
  uint32_t far_rule_count(void) { return rules_by_far_key_.size(); }
  ConvergedRuleStore() {}
  void insert_rule(uint32_t rule_id, const SetGroupPDR& rule);
  void insert_rule(uint32_t rule_id, const SetGroupFAR& rule);
  bool remove_rule(uint32_t rule_id, SetGroupPDR* rule);
  bool remove_rule(uint32_t rule_id, SetGroupFAR* rule);
  bool get_rule(uint32_t rule_id, SetGroupPDR* rule);
  bool get_rule(uint32_t rule_id, SetGroupFAR* rule);

 private:
  std::unordered_map<uint32_t, std::shared_ptr<SetGroupPDR>> rules_by_pdr_key_;
  std::unordered_map<uint32_t, std::shared_ptr<SetGroupFAR>> rules_by_far_key_;
};
}  // namespace magma
