/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <functional>
#include <string>
#include <unordered_set>
#include <utility>
#include <vector>
#include <google/protobuf/timestamp.pb.h>
#include <google/protobuf/util/time_util.h>

#include "CreditKey.h"
#include "RuleStore.h"
#include "SessionState.h"
#include "StoredState.h"
#include "magma_logging.h"

namespace magma {

std::string session_fsm_state_to_str(SessionFsmState state);

std::unique_ptr<SessionState> SessionState::unmarshal(
    const StoredSessionState& marshaled, StaticRuleStore& rule_store) {
  return std::make_unique<SessionState>(marshaled, rule_store);
}

StoredSessionState SessionState::marshal() {
  StoredSessionState marshaled{};

  marshaled.fsm_state              = curr_state_;
  marshaled.config                 = config_;
  marshaled.charging_pool          = charging_pool_.marshal();
  marshaled.monitor_pool           = monitor_pool_.marshal();
  marshaled.imsi                   = imsi_;
  marshaled.session_id             = session_id_;
  marshaled.core_session_id        = core_session_id_;
  marshaled.subscriber_quota_state = subscriber_quota_state_;
  marshaled.tgpp_context           = tgpp_context_;
  marshaled.request_number         = request_number_;

  for (auto& rule_id : active_static_rules_) {
    marshaled.static_rule_ids.push_back(rule_id);
  }
  std::vector<PolicyRule> dynamic_rules;
  dynamic_rules_.get_rules(dynamic_rules);
  marshaled.dynamic_rules = std::move(dynamic_rules);

  for (auto& rule_id : scheduled_static_rules_) {
    marshaled.scheduled_static_rules.insert(rule_id);
  }
  std::vector<PolicyRule> scheduled_dynamic_rules;
  scheduled_dynamic_rules_.get_rules(scheduled_dynamic_rules);
  marshaled.scheduled_dynamic_rules = std::move(scheduled_dynamic_rules);
  for (auto& it : rule_lifetimes_) {
    marshaled.rule_lifetimes[it.first] = it.second;
  }

  return marshaled;
}

SessionState::SessionState(
    const StoredSessionState& marshaled, StaticRuleStore& rule_store)
    : request_number_(marshaled.request_number),
      curr_state_(marshaled.fsm_state),
      config_(marshaled.config),
      imsi_(marshaled.imsi),
      session_id_(marshaled.session_id),
      core_session_id_(marshaled.core_session_id),
      subscriber_quota_state_(marshaled.subscriber_quota_state),
      tgpp_context_(marshaled.tgpp_context),
      charging_pool_(
          std::move(*ChargingCreditPool::unmarshal(marshaled.charging_pool))),
      monitor_pool_(std::move(
          *UsageMonitoringCreditPool::unmarshal(marshaled.monitor_pool))),
      static_rules_(rule_store) {
  for (const std::string& rule_id : marshaled.static_rule_ids) {
    active_static_rules_.push_back(rule_id);
  }
  for (auto& rule : marshaled.dynamic_rules) {
    dynamic_rules_.insert_rule(rule);
  }

  for (const std::string& rule_id : marshaled.scheduled_static_rules) {
    scheduled_static_rules_.insert(rule_id);
  }
  for (auto& rule : marshaled.scheduled_dynamic_rules) {
    scheduled_dynamic_rules_.insert_rule(rule);
  }
  for (auto& it : marshaled.rule_lifetimes) {
    rule_lifetimes_[it.first] = it.second;
  }
  for (auto& rule : marshaled.gy_dynamic_rules) {
    gy_dynamic_rules_.insert_rule(rule);
  }
}

SessionState::SessionState(
    const std::string& imsi, const std::string& session_id,
    const std::string& core_session_id, const SessionConfig& cfg,
    StaticRuleStore& rule_store, const magma::lte::TgppContext& tgpp_context)
    : imsi_(imsi),
      session_id_(session_id),
      core_session_id_(core_session_id),
      config_(cfg),
      // Request number set to 1, because request 0 is INIT call
      request_number_(1),
      curr_state_(SESSION_ACTIVE),
      charging_pool_(imsi),
      monitor_pool_(imsi),
      tgpp_context_(tgpp_context),
      static_rules_(rule_store) {}

void SessionState::new_report(SessionStateUpdateCriteria& update_criteria) {
  if (curr_state_ == SESSION_TERMINATING_FLOW_ACTIVE) {
    set_fsm_state(SESSION_TERMINATING_AGGREGATING_STATS, update_criteria);
  }
}

void SessionState::finish_report(SessionStateUpdateCriteria& update_criteria) {
  if (curr_state_ == SESSION_TERMINATING_AGGREGATING_STATS) {
    set_fsm_state(SESSION_TERMINATING_FLOW_DELETED, update_criteria);
  }
}

void SessionState::add_used_credit(
    const std::string& rule_id, uint64_t used_tx, uint64_t used_rx,
    SessionStateUpdateCriteria& update_criteria) {
  if (curr_state_ == SESSION_TERMINATING_AGGREGATING_STATS) {
    set_fsm_state(SESSION_TERMINATING_FLOW_ACTIVE,
                  update_criteria);
  }

  CreditKey charging_key;
  if (get_charging_key_for_rule_id(rule_id, &charging_key)) {
    MLOG(MINFO) << "Updating used charging credit for Rule=" << rule_id
                << " Rating Group=" << charging_key.rating_group
                << " Service Identifier=" << charging_key.service_identifier;
    charging_pool_.add_used_credit(
        charging_key, used_tx, used_rx, update_criteria);
  }
  std::string monitoring_key;
  if (get_monitoring_key_for_rule_id(rule_id, &monitoring_key)) {
    MLOG(MINFO) << "Updating used monitoring credit for Rule=" << rule_id
                << " Monitoring Key=" << monitoring_key;
    monitor_pool_.add_used_credit(
        monitoring_key, used_tx, used_rx, update_criteria);
  }
  auto session_level_key_p = monitor_pool_.get_session_level_key();
  if (session_level_key_p != nullptr &&
      monitoring_key != *session_level_key_p) {
    // Update session level key if its different
    monitor_pool_.add_used_credit(
        *session_level_key_p, used_tx, used_rx, update_criteria);
  }
}

void SessionState::set_subscriber_quota_state(
    const magma::lte::SubscriberQuotaUpdate_Type state,
    SessionStateUpdateCriteria& update_criteria) {
  update_criteria.updated_subscriber_quota_state = state;
  subscriber_quota_state_                        = state;
}

bool SessionState::active_monitored_rules_exist() {
  return total_monitored_rules_count() > 0;
}

SessionFsmState SessionState::get_state() {
  return curr_state_;
}

bool SessionState::is_terminating() {
  if (is_active() || curr_state_ == SESSION_TERMINATION_SCHEDULED) {
    return false;
  }
  return true;
}

void SessionState::get_updates_from_charging_pool(
    UpdateSessionRequest& update_request_out,
    std::vector<std::unique_ptr<ServiceAction>>* actions_out,
    SessionStateUpdateCriteria& update_criteria, const bool force_update) {
  // charging updates
  std::vector<CreditUsage> charging_updates;
  charging_pool_.get_updates(
      imsi_, config_.ue_ipv4, static_rules_, &dynamic_rules_, &charging_updates,
      actions_out, update_criteria, force_update);
  for (const auto& update : charging_updates) {
    auto new_req = update_request_out.mutable_updates()->Add();
    new_req->set_session_id(session_id_);
    new_req->set_request_number(request_number_);
    new_req->set_sid(imsi_);
    new_req->set_msisdn(config_.msisdn);
    new_req->set_ue_ipv4(config_.ue_ipv4);
    new_req->set_spgw_ipv4(config_.spgw_ipv4);
    new_req->set_apn(config_.apn);
    new_req->set_imei(config_.imei);
    new_req->set_plmn_id(config_.plmn_id);
    new_req->set_imsi_plmn_id(config_.imsi_plmn_id);
    new_req->set_user_location(config_.user_location);
    new_req->set_hardware_addr(config_.hardware_addr);
    new_req->set_rat_type(config_.rat_type);
    fill_protos_tgpp_context(new_req->mutable_tgpp_ctx());
    new_req->mutable_usage()->CopyFrom(update);
    request_number_++;
  }
}

void SessionState::get_updates_from_monitor_pool(
    UpdateSessionRequest& update_request_out,
    std::vector<std::unique_ptr<ServiceAction>>* actions_out,
    SessionStateUpdateCriteria& update_criteria, const bool force_update) {
  // monitor updates
  std::vector<UsageMonitorUpdate> monitor_updates;
  monitor_pool_.get_updates(
      imsi_, config_.ue_ipv4, static_rules_, &dynamic_rules_, &monitor_updates,
      actions_out, update_criteria, force_update);
  for (const auto& update : monitor_updates) {
    auto new_req = update_request_out.mutable_usage_monitors()->Add();
    new_req->set_session_id(session_id_);
    new_req->set_request_number(request_number_);
    new_req->set_sid(imsi_);
    new_req->set_ue_ipv4(config_.ue_ipv4);
    new_req->set_hardware_addr(config_.hardware_addr);
    new_req->set_rat_type(config_.rat_type);
    fill_protos_tgpp_context(new_req->mutable_tgpp_ctx());
    new_req->mutable_update()->CopyFrom(update);
    new_req->set_event_trigger(USAGE_REPORT);
    request_number_++;
  }
}

void SessionState::get_updates(
    UpdateSessionRequest& update_request_out,
    std::vector<std::unique_ptr<ServiceAction>>* actions_out,
    SessionStateUpdateCriteria& update_criteria, const bool force_update) {
  if (curr_state_ != SESSION_ACTIVE) return;
  get_updates_from_charging_pool(
      update_request_out, actions_out, update_criteria, force_update);
  get_updates_from_monitor_pool(
      update_request_out, actions_out, update_criteria, force_update);
}

void SessionState::start_termination(
    SessionStateUpdateCriteria& update_criteria) {
  set_fsm_state(SESSION_TERMINATING_FLOW_ACTIVE, update_criteria);
}

bool SessionState::can_complete_termination() const {
  return curr_state_ == SESSION_TERMINATING_FLOW_DELETED;
}

void SessionState::mark_as_awaiting_termination(
    SessionStateUpdateCriteria& update_criteria) {
  set_fsm_state(SESSION_TERMINATION_SCHEDULED, update_criteria);
}

SubscriberQuotaUpdate_Type SessionState::get_subscriber_quota_state() const {
  return subscriber_quota_state_;
}

void SessionState::complete_termination(
    SessionReporter& reporter, SessionStateUpdateCriteria& update_criteria) {
  switch (curr_state_) {
    case SESSION_ACTIVE:
      MLOG(MERROR) << imsi_ << " Encountered unexpected state 'ACTIVE' when "
                   << "forcefully completing termination. ";
      return;
    case SESSION_TERMINATED:
      // session is already terminated. Do nothing.
      return;
    case SESSION_TERMINATING_FLOW_ACTIVE:
    case SESSION_TERMINATING_AGGREGATING_STATS:
      MLOG(MINFO) << imsi_ << " Forcefully terminating session since it did "
                  << "not receive usage from pipelined in time.";
    default: // Continue termination but no logs are necessary for other states
      break;
  }
  // mark entire session as terminated
  set_fsm_state(SESSION_TERMINATED, update_criteria);
  auto termination_req = make_termination_request(update_criteria);
  auto logging_cb = SessionReporter::get_terminate_logging_cb(termination_req);
  reporter.report_terminate_session(termination_req, logging_cb);
}

SessionTerminateRequest SessionState::make_termination_request(
  SessionStateUpdateCriteria& update_criteria) {
  SessionTerminateRequest req;
  req.set_sid(imsi_);
  req.set_session_id(session_id_);
  req.set_request_number(request_number_);
  req.set_ue_ipv4(config_.ue_ipv4);
  req.set_msisdn(config_.msisdn);
  req.set_spgw_ipv4(config_.spgw_ipv4);
  req.set_apn(config_.apn);
  req.set_imei(config_.imei);
  req.set_plmn_id(config_.plmn_id);
  req.set_imsi_plmn_id(config_.imsi_plmn_id);
  req.set_user_location(config_.user_location);
  req.set_hardware_addr(config_.hardware_addr);
  req.set_rat_type(config_.rat_type);
  fill_protos_tgpp_context(req.mutable_tgpp_ctx());
  monitor_pool_.get_termination_updates(&req, update_criteria);
  charging_pool_.get_termination_updates(&req, update_criteria);
  return req;
}

ChargingCreditPool& SessionState::get_charging_pool() {
  return charging_pool_;
}

UsageMonitoringCreditPool& SessionState::get_monitor_pool() {
  return monitor_pool_;
}

SessionState::TotalCreditUsage SessionState::get_total_credit_usage() {
  // Collate unique charging/monitoring keys used by rules
  std::unordered_set<CreditKey, decltype(&ccHash), decltype(&ccEqual)>
      used_charging_keys(4, ccHash, ccEqual);
  std::unordered_set<std::string> used_monitoring_keys;

  std::vector<std::reference_wrapper<PolicyRuleBiMap>> bimaps{static_rules_,
                                                              dynamic_rules_};

  for (auto bimap : bimaps) {
    PolicyRuleBiMap& rules = bimap;
    std::vector<std::string> rule_ids{};
    std::vector<std::string>& rule_ids_ptr = rule_ids;
    rules.get_rule_ids(rule_ids_ptr);

    for (auto rule_id : rule_ids) {
      CreditKey charging_key;
      bool should_track_charging_key =
          rules.get_charging_key_for_rule_id(rule_id, &charging_key);
      std::string monitoring_key;
      bool should_track_monitoring_key =
          rules.get_monitoring_key_for_rule_id(rule_id, &monitoring_key);

      if (should_track_charging_key) used_charging_keys.insert(charging_key);
      if (should_track_monitoring_key)
        used_monitoring_keys.insert(monitoring_key);
    }
  }

  // Sum up usage
  TotalCreditUsage usage{
      .monitoring_tx = 0,
      .monitoring_rx = 0,
      .charging_tx   = 0,
      .charging_rx   = 0,
  };
  for (auto monitoring_key : used_monitoring_keys) {
    usage.monitoring_tx +=
        get_monitor_pool().get_credit(monitoring_key, USED_TX);
    usage.monitoring_rx +=
        get_monitor_pool().get_credit(monitoring_key, USED_RX);
  }
  for (auto charging_key : used_charging_keys) {
    usage.charging_tx += get_charging_pool().get_credit(charging_key, USED_TX);
    usage.charging_rx += get_charging_pool().get_credit(charging_key, USED_RX);
  }
  return usage;
}

bool SessionState::is_same_config(const SessionConfig& new_config) const {
  return config_.ue_ipv4.compare(new_config.ue_ipv4) == 0 &&
         config_.spgw_ipv4.compare(new_config.spgw_ipv4) == 0 &&
         config_.msisdn.compare(new_config.msisdn) == 0 &&
         config_.apn.compare(new_config.apn) == 0 &&
         config_.imei.compare(new_config.imei) == 0 &&
         config_.plmn_id.compare(new_config.plmn_id) == 0 &&
         config_.imsi_plmn_id.compare(new_config.imsi_plmn_id) == 0 &&
         config_.user_location.compare(new_config.user_location) == 0 &&
         config_.rat_type == new_config.rat_type &&
         config_.hardware_addr.compare(new_config.hardware_addr) == 0 &&
         config_.radius_session_id.compare(new_config.radius_session_id) == 0 &&
         config_.bearer_id == new_config.bearer_id;
}

std::string SessionState::get_session_id() const {
  return session_id_;
}

SessionConfig SessionState::get_config() {
  return config_;
}

void SessionState::set_config(const SessionConfig& config) {
  config_ = config;
}

bool SessionState::is_radius_cwf_session() const {
  return (config_.rat_type == RATType::TGPP_WLAN);
}

void SessionState::get_session_info(SessionState::SessionInfo& info) {
  info.imsi    = imsi_;
  info.ip_addr = config_.ue_ipv4;
  get_dynamic_rules().get_rules(info.dynamic_rules);
  get_gy_dynamic_rules().get_rules(info.gy_dynamic_rules);
  info.static_rules = active_static_rules_;
}

void SessionState::set_tgpp_context(
    const magma::lte::TgppContext& tgpp_context,
    SessionStateUpdateCriteria& update_criteria) {
  update_criteria.updated_tgpp_context = tgpp_context;
  tgpp_context_                        = tgpp_context;
}

void SessionState::fill_protos_tgpp_context(
    magma::lte::TgppContext* tgpp_context) const {
  *tgpp_context = tgpp_context_;
}

uint32_t SessionState::get_request_number() {
  return request_number_;
}

void SessionState::increment_request_number(uint32_t incr) {
  request_number_ += incr;
}

bool SessionState::get_charging_key_for_rule_id(
    const std::string& rule_id, CreditKey* charging_key) {
  // first check dynamic rules and then static rules
  if (dynamic_rules_.get_charging_key_for_rule_id(rule_id, charging_key)) {
    return true;
  }
  return static_rules_.get_charging_key_for_rule_id(rule_id, charging_key);
}

bool SessionState::get_monitoring_key_for_rule_id(
    const std::string& rule_id, std::string* monitoring_key) {
  // first check dynamic rules and then static rules
  if (dynamic_rules_.get_monitoring_key_for_rule_id(rule_id, monitoring_key)) {
    return true;
  }
  return static_rules_.get_monitoring_key_for_rule_id(rule_id, monitoring_key);
}

bool SessionState::is_dynamic_rule_scheduled(const std::string& rule_id) {
  auto _ = new PolicyRule();
  return scheduled_dynamic_rules_.get_rule(rule_id, _);
}

bool SessionState::is_static_rule_scheduled(const std::string& rule_id) {
  return scheduled_static_rules_.count(rule_id) == 1;
}

bool SessionState::is_dynamic_rule_installed(const std::string& rule_id) {
  auto _ = new PolicyRule();
  return dynamic_rules_.get_rule(rule_id, _);
}

bool SessionState::is_gy_dynamic_rule_installed(const std::string& rule_id) {
  auto _ = new PolicyRule();
  return gy_dynamic_rules_.get_rule(rule_id, _);
}

bool SessionState::is_static_rule_installed(const std::string& rule_id) {
  return std::find(
             active_static_rules_.begin(), active_static_rules_.end(),
             rule_id) != active_static_rules_.end();
}

void SessionState::insert_dynamic_rule(
    const PolicyRule& rule, RuleLifetime& lifetime,
    SessionStateUpdateCriteria& update_criteria) {
  if (is_dynamic_rule_installed(rule.id())) {
    return;
  }
  rule_lifetimes_[rule.id()] = lifetime;
  dynamic_rules_.insert_rule(rule);
  update_criteria.dynamic_rules_to_install.push_back(rule);
  update_criteria.new_rule_lifetimes[rule.id()] = lifetime;
}

void SessionState::insert_gy_dynamic_rule(
    const PolicyRule& rule,  SessionStateUpdateCriteria& update_criteria) {
  if (is_gy_dynamic_rule_installed(rule.id())) {
    return;
  }
  update_criteria.dynamic_rules_to_install.push_back(rule);
  gy_dynamic_rules_.insert_rule(rule);
}

void SessionState::activate_static_rule(
    const std::string& rule_id, RuleLifetime& lifetime,
    SessionStateUpdateCriteria& update_criteria) {
  rule_lifetimes_[rule_id] = lifetime;
  active_static_rules_.push_back(rule_id);
  update_criteria.static_rules_to_install.insert(rule_id);
  update_criteria.new_rule_lifetimes[rule_id] = lifetime;
}

bool SessionState::remove_dynamic_rule(
    const std::string& rule_id, PolicyRule* rule_out,
    SessionStateUpdateCriteria& update_criteria) {
  bool removed = dynamic_rules_.remove_rule(rule_id, rule_out);
  if (removed) {
    update_criteria.dynamic_rules_to_uninstall.insert(rule_id);
  }
  return removed;
}

bool SessionState::remove_scheduled_dynamic_rule(
    const std::string& rule_id, PolicyRule* rule_out,
    SessionStateUpdateCriteria& update_criteria) {
  bool removed = scheduled_dynamic_rules_.remove_rule(rule_id, rule_out);
  if (removed) {
    update_criteria.dynamic_rules_to_uninstall.insert(rule_id);
  }
  return removed;
}

bool SessionState::remove_gy_dynamic_rule(
  const std::string& rule_id, PolicyRule *rule_out,
  SessionStateUpdateCriteria& update_criteria)
{
  bool removed = gy_dynamic_rules_.remove_rule(rule_id, rule_out);
  if (removed) {
    update_criteria.dynamic_rules_to_uninstall.insert(rule_id);
  }
  return removed;
}

bool SessionState::deactivate_static_rule(
    const std::string& rule_id, SessionStateUpdateCriteria& update_criteria) {
  auto it = std::find(
      active_static_rules_.begin(), active_static_rules_.end(), rule_id);
  if (it == active_static_rules_.end()) {
    return false;
  }
  update_criteria.static_rules_to_uninstall.insert(rule_id);
  active_static_rules_.erase(it);
  return true;
}

bool SessionState::deactivate_scheduled_static_rule(
    const std::string& rule_id, SessionStateUpdateCriteria& update_criteria) {
  if (scheduled_static_rules_.count(rule_id) == 0) {
    return false;
  }
  scheduled_static_rules_.erase(rule_id);
  return true;
}

void SessionState::sync_rules_to_time(
    std::time_t current_time, SessionStateUpdateCriteria& update_criteria) {
  PolicyRule _rule_unused;
  // Update active static rules
  for (const std::string& rule_id : active_static_rules_) {
    if (should_rule_be_deactivated(rule_id, current_time)) {
      deactivate_static_rule(rule_id, update_criteria);
    }
  }
  // Update scheduled static rules
  std::set<std::string> scheduled_rule_ids = scheduled_static_rules_;
  for (const std::string& rule_id : scheduled_rule_ids) {
    if (should_rule_be_active(rule_id, current_time)) {
      install_scheduled_static_rule(rule_id, update_criteria);
    } else if (should_rule_be_deactivated(rule_id, current_time)) {
      scheduled_static_rules_.erase(rule_id);
      update_criteria.static_rules_to_uninstall.insert(rule_id);
    }
  }
  // Update active dynamic rules
  std::vector<std::string> dynamic_rule_ids;
  dynamic_rules_.get_rule_ids(dynamic_rule_ids);
  for (const std::string& rule_id : dynamic_rule_ids) {
    if (should_rule_be_deactivated(rule_id, current_time)) {
      remove_dynamic_rule(rule_id, &_rule_unused, update_criteria);
    }
  }
  // Update scheduled dynamic rules
  scheduled_dynamic_rules_.get_rule_ids(dynamic_rule_ids);
  for (const std::string& rule_id : dynamic_rule_ids) {
    if (should_rule_be_active(rule_id, current_time)) {
      install_scheduled_dynamic_rule(rule_id, update_criteria);
    } else if (should_rule_be_deactivated(rule_id, current_time)) {
      remove_scheduled_dynamic_rule(rule_id, &_rule_unused, update_criteria);
    }
  }
}

std::vector<std::string>& SessionState::get_static_rules() {
  return active_static_rules_;
}

std::set<std::string>& SessionState::get_scheduled_static_rules() {
  return scheduled_static_rules_;
}

DynamicRuleStore& SessionState::get_dynamic_rules() {
  return dynamic_rules_;
}

DynamicRuleStore& SessionState::get_scheduled_dynamic_rules() {
  return scheduled_dynamic_rules_;
}

RuleLifetime& SessionState::get_rule_lifetime(const std::string& rule_id) {
  return rule_lifetimes_[rule_id];
}

DynamicRuleStore& SessionState::get_gy_dynamic_rules()
{
  return gy_dynamic_rules_;
}

uint32_t SessionState::total_monitored_rules_count() {
  uint32_t monitored_dynamic_rules = dynamic_rules_.monitored_rules_count();
  uint32_t monitored_static_rules  = 0;
  for (auto& rule_id : active_static_rules_) {
    std::string _;
    auto is_monitored =
        static_rules_.get_monitoring_key_for_rule_id(rule_id, &_);
    if (is_monitored) {
      monitored_static_rules++;
    }
  }
  return monitored_dynamic_rules + monitored_static_rules;
}

void SessionState::schedule_dynamic_rule(
    const PolicyRule& rule, RuleLifetime& lifetime,
    SessionStateUpdateCriteria& update_criteria) {
  update_criteria.new_rule_lifetimes[rule.id()] = lifetime;
  update_criteria.new_scheduled_dynamic_rules.push_back(rule);
  rule_lifetimes_[rule.id()] = lifetime;
  scheduled_dynamic_rules_.insert_rule(rule);
}

void SessionState::schedule_static_rule(
    const std::string& rule_id, RuleLifetime& lifetime,
    SessionStateUpdateCriteria& update_criteria) {
  update_criteria.new_rule_lifetimes[rule_id] = lifetime;
  update_criteria.new_scheduled_static_rules.insert(rule_id);
  rule_lifetimes_[rule_id] = lifetime;
  scheduled_static_rules_.insert(rule_id);
}

void SessionState::install_scheduled_dynamic_rule(
    const std::string& rule_id, SessionStateUpdateCriteria& update_criteria) {
  PolicyRule dynamic_rule;
  bool removed = scheduled_dynamic_rules_.remove_rule(rule_id, &dynamic_rule);
  if (!removed) {
    MLOG(MERROR) << "Failed to mark a scheduled dynamic rule as installed "
                 << "with rule_id: " << rule_id;
    return;
  }
  update_criteria.dynamic_rules_to_install.push_back(dynamic_rule);
  dynamic_rules_.insert_rule(dynamic_rule);
}

void SessionState::install_scheduled_static_rule(
    const std::string& rule_id, SessionStateUpdateCriteria& update_criteria) {
  auto it = scheduled_static_rules_.find(rule_id);
  if (it == scheduled_static_rules_.end()) {
    MLOG(MERROR) << "Failed to mark a scheduled static rule as installed "
                    "with rule_id: "
                 << rule_id;
  }
  update_criteria.static_rules_to_install.insert(rule_id);
  scheduled_static_rules_.erase(rule_id);
  active_static_rules_.push_back(rule_id);
}

uint32_t SessionState::get_credit_key_count() {
  return charging_pool_.get_credit_key_count() +
         monitor_pool_.get_credit_key_count();
}

bool SessionState::is_active() {
  return curr_state_ == SESSION_ACTIVE;
}

void SessionState::set_fsm_state(SessionFsmState new_state,
                                 SessionStateUpdateCriteria& uc) {
  // Only log and reflect change into update criteria if the state is new
  if (curr_state_ != new_state) {
    MLOG(MDEBUG) << "Session " << session_id_ << " FSM state change from "
                 << session_fsm_state_to_str(curr_state_) << " to "
                 << session_fsm_state_to_str(new_state);
    curr_state_ = new_state;
    uc.is_fsm_updated = true;
    uc.updated_fsm_state = new_state;
  }
}

std::string session_fsm_state_to_str(SessionFsmState state) {
  switch (state) {
  case SESSION_ACTIVE:
    return "SESSION_ACTIVE";
  case SESSION_TERMINATING_FLOW_ACTIVE:
    return "SESSION_TERMINATING_FLOW_ACTIVE";
  case SESSION_TERMINATING_AGGREGATING_STATS:
    return "SESSION_TERMINATING_AGGREGATING_STATS";
  case SESSION_TERMINATING_FLOW_DELETED:
    return "SESSION_TERMINATING_FLOW_DELETED";
  case SESSION_TERMINATED:
    return "SESSION_TERMINATED";
  case SESSION_TERMINATION_SCHEDULED:
    return "SESSION_TERMINATION_SCHEDULED";
  default:
    return "INVALID SESSION FSM STATE";
  }
}

bool SessionState::should_rule_be_active(
    const std::string& rule_id, std::time_t time) {
  auto lifetime = rule_lifetimes_[rule_id];
  bool deactivated =
      (lifetime.deactivation_time > 0) && (lifetime.deactivation_time < time);
  return lifetime.activation_time < time && !deactivated;
}

bool SessionState::should_rule_be_deactivated(
    const std::string& rule_id, std::time_t time) {
  auto lifetime = rule_lifetimes_[rule_id];
  return lifetime.deactivation_time > 0 && lifetime.deactivation_time < time;
}

StaticRuleInstall SessionState::get_static_rule_install(
  const std::string& rule_id, const RuleLifetime& lifetime) {
  StaticRuleInstall rule_install{};
  rule_install.set_rule_id(rule_id);
  rule_install.mutable_activation_time()->set_seconds(lifetime.activation_time);
  rule_install.mutable_deactivation_time()->set_seconds(lifetime.deactivation_time);
  return rule_install;
}

DynamicRuleInstall SessionState::get_dynamic_rule_install(
  const std::string& rule_id, const RuleLifetime& lifetime) {
  DynamicRuleInstall rule_install{};
  PolicyRule* policy_rule = rule_install.mutable_policy_rule();
  if (!dynamic_rules_.get_rule(rule_id, policy_rule)) {
    scheduled_dynamic_rules_.get_rule(rule_id, policy_rule);
  }
  rule_install.mutable_activation_time()->set_seconds(lifetime.activation_time);
  rule_install.mutable_deactivation_time()->set_seconds(lifetime.deactivation_time);
  return rule_install;
}
}  // namespace magma
