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
#include <glog/logging.h>

#include <memory>
#include <string>
#include <vector>

#include "RuleStore.h"
#include "includes/ServiceRegistrySingleton.h"

namespace magma {

template<typename KeyType, typename hash, typename equal>
void PoliciesByKeyMap<KeyType, hash, equal>::insert(
    const KeyType& key, std::shared_ptr<PolicyRule> rule_p) {
  auto iter = rules_by_key_.find(key);
  if (iter == rules_by_key_.end()) {
    rules_by_key_[key] = {rule_p};
    return;
  }
  iter->second.push_back(rule_p);
}

template<typename KeyType, typename hash, typename equal>
void PoliciesByKeyMap<KeyType, hash, equal>::remove(
    const KeyType& key, std::shared_ptr<PolicyRule> rule_p) {
  auto iter = rules_by_key_.find(key);
  if (iter == rules_by_key_.end()) {
    return;
  }

  auto rules = iter->second;
  auto found = std::find(rules.begin(), rules.end(), rule_p);
  if (found == rules.end()) {
    return;
  }
  rules.erase(found);
  rules_by_key_[key] = rules;
}

template<typename KeyType, typename hash, typename equal>
bool PoliciesByKeyMap<KeyType, hash, equal>::get_rule_ids_for_key(
    const KeyType& key, std::vector<std::string>& rules_out) {
  auto iter = rules_by_key_.find(key);
  if (iter == rules_by_key_.end()) {
    return false;
  }

  for (const auto& rule : iter->second) {
    rules_out.push_back(rule->id());
  }
  return true;
}

template<typename KeyType, typename hash, typename equal>
bool PoliciesByKeyMap<KeyType, hash, equal>::get_rule_definitions_for_key(
    const KeyType& key, std::vector<PolicyRule>& rules_out) {
  auto iter = rules_by_key_.find(key);
  if (iter == rules_by_key_.end()) {
    return false;
  }

  for (const auto& rule : iter->second) {
    rules_out.push_back(*rule);
  }
  return true;
}

template<typename KeyType, typename hash, typename equal>
uint32_t PoliciesByKeyMap<KeyType, hash, equal>::policy_count() {
  uint32_t count = 0;
  for (auto const& kv : rules_by_key_) {
    count += kv.second.size();
  }
  return count;
}

static bool should_track_charging_key(PolicyRule::TrackingType tracking_type) {
  return tracking_type == PolicyRule::ONLY_OCS ||
         tracking_type == PolicyRule::OCS_AND_PCRF;
}

static bool should_track_monitoring_key(
    PolicyRule::TrackingType tracking_type) {
  return tracking_type == PolicyRule::ONLY_PCRF ||
         tracking_type == PolicyRule::OCS_AND_PCRF;
}

void PolicyRuleBiMap::sync_rules(const std::vector<PolicyRule>& rules) {
  std::lock_guard<std::mutex> lock(map_mutex_);
  rules_by_rule_id_.clear();
  rules_by_charging_key_ =
      PoliciesByKeyMap<CreditKey, decltype(&ccHash), decltype(&ccEqual)>(
          &ccHash, &ccEqual);
  rules_by_monitoring_key_ = PoliciesByKeyMap<std::string>();
  for (const auto& rule : rules) {
    auto rule_p                  = std::make_shared<PolicyRule>(rule);
    rules_by_rule_id_[rule.id()] = rule_p;
    if (should_track_charging_key(rule.tracking_type())) {
      rules_by_charging_key_.insert(CreditKey(rule), rule_p);
    }
    if (should_track_monitoring_key(rule.tracking_type())) {
      rules_by_monitoring_key_.insert(rule.monitoring_key(), rule_p);
    }
  }
}

void PolicyRuleBiMap::insert_rule(const PolicyRule& rule) {
  auto rule_p = std::make_shared<PolicyRule>(rule);
  std::lock_guard<std::mutex> lock(map_mutex_);
  rules_by_rule_id_[rule.id()] = rule_p;
  if (should_track_charging_key(rule.tracking_type())) {
    rules_by_charging_key_.insert(CreditKey(rule), rule_p);
  }
  if (should_track_monitoring_key(rule.tracking_type())) {
    rules_by_monitoring_key_.insert(rule.monitoring_key(), rule_p);
  }
}

bool PolicyRuleBiMap::get_rule(
    const std::string& rule_id, PolicyRule* rule_out) {
  std::lock_guard<std::mutex> lock(map_mutex_);
  auto it = rules_by_rule_id_.find(rule_id);
  if (it == rules_by_rule_id_.end()) {
    return false;
  }
  if (rule_out != NULL) {
    rule_out->CopyFrom(*it->second);
  }
  return true;
}

bool PolicyRuleBiMap::get_rules_by_ids(
    const std::vector<std::string>& rule_ids,
    std::vector<PolicyRule>& rules_out) {
  std::lock_guard<std::mutex> lock(map_mutex_);
  for (const std::string rule_id : rule_ids) {
    auto it = rules_by_rule_id_.find(rule_id);
    if (it == rules_by_rule_id_.end()) {
      return false;
    }
    rules_out.push_back(*it->second);
  }
  return true;
}

bool PolicyRuleBiMap::remove_rule(
    const std::string& rule_id, PolicyRule* rule_out) {
  std::lock_guard<std::mutex> lock(map_mutex_);
  auto it = rules_by_rule_id_.find(rule_id);
  if (it == rules_by_rule_id_.end()) {
    return false;
  }

  auto rule_ptr = it->second;
  if (rule_out != NULL) {
    rule_out->CopyFrom(*rule_ptr);
  }

  // Remove the rule from all mappings
  rules_by_rule_id_.erase(it);
  if (should_track_charging_key(rule_ptr->tracking_type())) {
    rules_by_charging_key_.remove(CreditKey(rule_ptr.get()), rule_ptr);
  }
  if (should_track_monitoring_key(rule_ptr->tracking_type())) {
    rules_by_monitoring_key_.remove(rule_ptr->monitoring_key(), rule_ptr);
  }

  return true;
}

bool PolicyRuleBiMap::get_charging_key_for_rule_id(
    const std::string& rule_id, CreditKey* charging_key) {
  std::lock_guard<std::mutex> lock(map_mutex_);
  auto it = rules_by_rule_id_.find(rule_id);
  if (it == rules_by_rule_id_.end()) {
    return false;
  }
  if (should_track_charging_key(it->second->tracking_type())) {
    charging_key->set(it->second.get());
    return true;
  }
  return false;
}

bool PolicyRuleBiMap::get_monitoring_key_for_rule_id(
    const std::string& rule_id, std::string* monitoring_key) {
  std::lock_guard<std::mutex> lock(map_mutex_);
  auto it = rules_by_rule_id_.find(rule_id);
  if (it == rules_by_rule_id_.end() ||
      !should_track_monitoring_key(it->second->tracking_type())) {
    return false;
  }
  // nullptr means the caller does not care about retrieving the value
  if (monitoring_key != nullptr) {
    monitoring_key->assign(it->second->monitoring_key());
  }
  return true;
}

bool PolicyRuleBiMap::get_rule_ids_for_charging_key(
    const CreditKey& charging_key, std::vector<std::string>& rules_out) {
  std::lock_guard<std::mutex> lock(map_mutex_);
  bool success =
      rules_by_charging_key_.get_rule_ids_for_key(charging_key, rules_out);
  return success;
}

bool PolicyRuleBiMap::get_rule_definitions_for_charging_key(
    const CreditKey& charging_key, std::vector<PolicyRule>& rules_out) {
  std::lock_guard<std::mutex> lock(map_mutex_);
  bool success = rules_by_charging_key_.get_rule_definitions_for_key(
      charging_key, rules_out);
  return success;
}

uint32_t PolicyRuleBiMap::monitored_rules_count() {
  std::lock_guard<std::mutex> lock(map_mutex_);
  return rules_by_monitoring_key_.policy_count();
}

bool PolicyRuleBiMap::get_rule_ids(std::vector<std::string>& rules_ids_out) {
  std::lock_guard<std::mutex> lock(map_mutex_);
  for (auto kv : rules_by_rule_id_) {
    rules_ids_out.push_back(kv.first);
  }
  return true;
}

bool PolicyRuleBiMap::get_rules(std::vector<PolicyRule>& rules_out) {
  std::lock_guard<std::mutex> lock(map_mutex_);
  for (auto kv : rules_by_rule_id_) {
    rules_out.push_back(*kv.second);
  }
  return true;
}

void ConvergedRuleStore::insert_rule(uint32_t id, const SetGroupPDR& rule) {
  auto rule_p = std::make_shared<SetGroupPDR>(rule);
  std::lock_guard<std::mutex> lock(map_mutex_);
  rules_by_pdr_key_[id] = rule_p;
}

void ConvergedRuleStore::insert_rule(uint32_t id, const SetGroupFAR& rule) {
  auto rule_p = std::make_shared<SetGroupFAR>(rule);
  std::lock_guard<std::mutex> lock(map_mutex_);
  rules_by_far_key_[id] = rule_p;
}

bool ConvergedRuleStore::remove_rule(uint32_t rule_id, SetGroupPDR* rule_out) {
  std::lock_guard<std::mutex> lock(map_mutex_);
  auto it = rules_by_pdr_key_.find(rule_id);
  if (it != rules_by_pdr_key_.end()) {
    auto rule_ptr = it->second;
    if (rule_out != NULL) {
      rule_out->CopyFrom(*rule_ptr);
    }
    rules_by_pdr_key_.erase(rule_id);
    return true;
  }
  return false;
}

bool ConvergedRuleStore::remove_rule(uint32_t rule_id, SetGroupFAR* rule_out) {
  std::lock_guard<std::mutex> lock(map_mutex_);
  auto it = rules_by_far_key_.find(rule_id);
  if (it != rules_by_far_key_.end()) {
    auto rule_ptr = it->second;
    if (rule_out != NULL) {
      rule_out->CopyFrom(*rule_ptr);
    }
    rules_by_far_key_.erase(rule_id);
    return true;
  }
  return false;
}

bool ConvergedRuleStore::get_rule(uint32_t rule_id, SetGroupPDR* rule_out) {
  std::lock_guard<std::mutex> lock(map_mutex_);
  auto it = rules_by_pdr_key_.find(rule_id);
  if (it == rules_by_pdr_key_.end()) {
    return false;
  }
  if (rule_out != NULL) {
    rule_out->CopyFrom(*it->second);
  }
  return true;
}

bool ConvergedRuleStore::get_rule(uint32_t rule_id, SetGroupFAR* rule_out) {
  std::lock_guard<std::mutex> lock(map_mutex_);
  auto it = rules_by_far_key_.find(rule_id);
  if (it == rules_by_far_key_.end()) {
    return false;
  }
  if (rule_out != NULL) {
    rule_out->CopyFrom(*it->second);
  }
  return true;
}

}  // namespace magma
