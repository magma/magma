/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "SessionState.h"
#include "SessionStore.h"
#include "StoredState.h"
#include "magma_logging.h"

namespace magma {
namespace lte {

SessionStore::SessionStore(std::shared_ptr<StaticRuleStore> rule_store):
  rule_store_(rule_store), store_client_(rule_store)
{
}

SessionMap SessionStore::read_sessions(const SessionRead& req)
{
  // First allocate some request numbers
  auto subscriber_ids = std::vector<std::string> {};
  for (const auto& it : req) {
    subscriber_ids.push_back(it.first);
  }
  auto session_map = store_client_.read_sessions(subscriber_ids);
  auto session_map_2 = store_client_.read_sessions(subscriber_ids);

  // For all sessions of the subscriber, increment the request numbers
  for (const auto& it : req) {
    for (auto& session : session_map_2[it.first]) {
      session->increment_request_number(it.second);
    }
  }
  store_client_.write_sessions(std::move(session_map_2));
  return session_map;
}

bool SessionStore::create_sessions(
  const std::string& subscriber_id,
  std::vector<std::unique_ptr<SessionState>> sessions)
{
  auto session_map = SessionMap {};
  session_map[subscriber_id] = std::move(sessions);
  store_client_.write_sessions(std::move(session_map));
  return true;
}

bool SessionStore::update_sessions(const SessionUpdate& update_criteria)
{
  // Read the current state
  auto subscriber_ids = std::vector<std::string> {};
  for (const auto& it : update_criteria) {
    subscriber_ids.push_back(it.first);
  }
  auto session_map = store_client_.read_sessions(subscriber_ids);

  // Now attempt to modify the state
  for (auto& it : session_map) {
    for (auto& session : it.second) {
      auto updates = update_criteria.find(it.first)->second;
      if (updates.find(session->get_session_id()) != updates.end()) {
        if (!merge_into_session(session, updates[session->get_session_id()])) {
          return false;
        }
      }
    }
  }
  return store_client_.write_sessions(std::move(session_map));
}

bool SessionStore::merge_into_session(
  std::unique_ptr<SessionState>& session,
  const SessionStateUpdateCriteria& update_criteria)
{
  // Static rules
  for (const auto& rule_id : update_criteria.static_rules_to_install) {
    if (session->is_static_rule_installed(rule_id)) {
      return false;
    }
    session->activate_static_rule(rule_id);
  }
  for (const auto& rule_id : update_criteria.static_rules_to_uninstall) {
    if (!session->is_static_rule_installed(rule_id)) {
      return false;
    }
    session->deactivate_static_rule(rule_id);
  }

  // Dynamic rules
  for (const auto& rule : update_criteria.dynamic_rules_to_install) {
    if (session->is_dynamic_rule_installed(rule.id())) {
      return false;
    }
    session->insert_dynamic_rule(rule);
  }
  PolicyRule* _ = {};
  for (const auto& rule_id : update_criteria.dynamic_rules_to_uninstall) {
    if (!session->is_dynamic_rule_installed(rule_id)) {
      return false;
    }
    session->remove_dynamic_rule(rule_id, _);
  }

  // Charging credit
  for (const auto& it : update_criteria.charging_credit_map) {
    auto key = it.first;
    auto credit_update = it.second;
    session->get_charging_pool().merge_credit_update(key, credit_update);
  }
  for (const auto& it : update_criteria.charging_credit_to_install) {
    auto key = it.first;
    auto stored_credit = it.second;
    session->get_charging_pool().add_credit(
      key, SessionCredit::unmarshal(stored_credit, CHARGING));
  }

  // Monitoring credit
  for (const auto& it : update_criteria.monitor_credit_map) {
    auto key = it.first;
    auto credit_update = it.second;
    session->get_monitor_pool().merge_credit_update(key, credit_update);
  }
  for (const auto& it : update_criteria.monitor_credit_to_install) {
    auto key = it.first;
    auto stored_monitor = it.second;
    session->get_monitor_pool().add_monitor(
      key, UsageMonitoringCreditPool::unmarshal_monitor(stored_monitor));
  }
  return true;
}

} // namespace lte
} // namespace magma
