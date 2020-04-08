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
#include <utility>
#include <unordered_set>
#include <vector>

#include "RuleStore.h"
#include "SessionState.h"
#include "StoredState.h"
#include "magma_logging.h"
#include "CreditKey.h"

namespace magma {

std::unique_ptr<SessionState> SessionState::unmarshal(
  const StoredSessionState &marshaled, StaticRuleStore &rule_store)
{
  return std::make_unique<SessionState>(marshaled, rule_store);
}

StoredSessionState SessionState::marshal()
{
  StoredSessionState marshaled{};

  marshaled.config = marshal_config();
  marshaled.charging_pool = charging_pool_.marshal();
  marshaled.monitor_pool = monitor_pool_.marshal();
  marshaled.imsi = imsi_;
  marshaled.session_id = session_id_;
  marshaled.core_session_id = core_session_id_;
  marshaled.subscriber_quota_state = subscriber_quota_state_;
  marshaled.tgpp_context = tgpp_context_;
  marshaled.request_number = request_number_;

  for (auto& rule_id : active_static_rules_) {
    marshaled.static_rule_ids.push_back(rule_id);
  }
  std::vector<PolicyRule> dynamic_rules;
  dynamic_rules_.get_rules(dynamic_rules);
  marshaled.dynamic_rules = std::move(dynamic_rules);

  return marshaled;
}

StoredSessionConfig SessionState::marshal_config()
{
  StoredSessionConfig config{};
  config.ue_ipv4 = config_.ue_ipv4;
  config.spgw_ipv4 = config_.spgw_ipv4;
  config.msisdn = config_.msisdn;
  config.apn = config_.apn;
  config.imei = config_.apn;
  config.plmn_id = config_.plmn_id;
  config.imsi_plmn_id = config_.imsi_plmn_id;
  config.user_location = config_.user_location;
  config.rat_type = config_.rat_type;
  config.mac_addr = config_.mac_addr; // MAC Address for WLAN
  config.hardware_addr = config_.hardware_addr; // MAC Address for WLAN (binary)
  config.radius_session_id = config_.radius_session_id;
  config.bearer_id = config_.bearer_id;
  StoredQoSInfo qos_info{};
  qos_info.enabled = config_.qos_info.enabled;
  qos_info.qci = config_.qos_info.qci;
  config.qos_info = qos_info;
  return config;
}

void SessionState::unmarshal_config(const StoredSessionConfig& marshaled)
{
  SessionState::Config cfg{};
  cfg.ue_ipv4 = marshaled.ue_ipv4;
  cfg.spgw_ipv4 = marshaled.spgw_ipv4;
  cfg.msisdn = marshaled.msisdn;
  cfg.apn = marshaled.apn;
  cfg.imei = marshaled.apn;
  cfg.plmn_id = marshaled.plmn_id;
  cfg.imsi_plmn_id = marshaled.imsi_plmn_id;
  cfg.user_location = marshaled.user_location;
  cfg.rat_type = marshaled.rat_type;
  cfg.mac_addr = marshaled.mac_addr; // MAC Address for WLAN
  cfg.hardware_addr = marshaled.hardware_addr; // MAC Address for WLAN (binary)
  cfg.radius_session_id = marshaled.radius_session_id;
  cfg.bearer_id = marshaled.bearer_id;
  SessionState::QoSInfo qos_info{};
  qos_info.enabled = marshaled.qos_info.enabled;
  qos_info.qci = marshaled.qos_info.qci;
  cfg.qos_info = qos_info;
  config_ = cfg;
}

SessionState::SessionState(
  const StoredSessionState &marshaled,
  StaticRuleStore &rule_store):
  request_number_(marshaled.request_number),
  curr_state_(SESSION_ACTIVE),
  charging_pool_(std::move(*ChargingCreditPool::unmarshal(marshaled.charging_pool))),
  monitor_pool_(std::move(*UsageMonitoringCreditPool::unmarshal(marshaled.monitor_pool))),
  static_rules_(rule_store)
{
SessionState::Config cfg{};
  StoredSessionConfig marshaled_cfg = marshaled.config;
  cfg.ue_ipv4 = marshaled_cfg.ue_ipv4;
  cfg.spgw_ipv4 = marshaled_cfg.spgw_ipv4;
  cfg.msisdn = marshaled_cfg.msisdn;
  cfg.apn = marshaled_cfg.apn;
  cfg.imei = marshaled_cfg.apn;
  cfg.plmn_id = marshaled_cfg.plmn_id;
  cfg.imsi_plmn_id = marshaled_cfg.imsi_plmn_id;
  cfg.user_location = marshaled_cfg.user_location;
  cfg.rat_type = marshaled_cfg.rat_type;
  cfg.mac_addr = marshaled_cfg.mac_addr; // MAC Address for WLAN
  cfg.hardware_addr = marshaled_cfg.hardware_addr; // MAC Address for WLAN (binary)
  cfg.radius_session_id = marshaled_cfg.radius_session_id;
  cfg.bearer_id = marshaled_cfg.bearer_id;
  SessionState::QoSInfo qos_info{};
  qos_info.enabled = marshaled_cfg.qos_info.enabled;
  qos_info.qci = marshaled_cfg.qos_info.qci;
  cfg.qos_info = qos_info;
  config_ = cfg;

  imsi_ = marshaled.imsi;
  session_id_ = marshaled.session_id;
  core_session_id_ = marshaled.core_session_id;
  subscriber_quota_state_ = marshaled.subscriber_quota_state;
  tgpp_context_ = marshaled.tgpp_context;

  for (const std::string& rule_id : marshaled.static_rule_ids)
  {
    active_static_rules_.push_back(rule_id);
  }
  for (auto& rule : marshaled.dynamic_rules)
  {
    dynamic_rules_.insert_rule(rule);
  }
}

SessionState::SessionState(
  const std::string& imsi,
  const std::string& session_id,
  const std::string& core_session_id,
  const SessionState::Config& cfg,
  StaticRuleStore& rule_store,
  const magma::lte::TgppContext& tgpp_context):
  imsi_(imsi),
  session_id_(session_id),
  core_session_id_(core_session_id),
  config_(cfg),
  // Request number set to 2, because request 1 is INIT call
  request_number_(2),
  curr_state_(SESSION_ACTIVE),
  charging_pool_(imsi),
  monitor_pool_(imsi),
  tgpp_context_(tgpp_context),
  static_rules_(rule_store)
{
}

void SessionState::new_report()
{
  if (curr_state_ == SESSION_TERMINATING_FLOW_ACTIVE) {
    curr_state_ = SESSION_TERMINATING_AGGREGATING_STATS;
  }
}

void SessionState::finish_report()
{
  if (curr_state_ == SESSION_TERMINATING_AGGREGATING_STATS) {
    curr_state_ = SESSION_TERMINATING_FLOW_DELETED;
  }
}

void SessionState::add_used_credit(
  const std::string& rule_id,
  uint64_t used_tx,
  uint64_t used_rx,
  SessionStateUpdateCriteria& update_criteria)
{
  if (curr_state_ == SESSION_TERMINATING_AGGREGATING_STATS) {
    curr_state_ = SESSION_TERMINATING_FLOW_ACTIVE;
  }

  CreditKey charging_key;
  if (get_charging_key_for_rule_id(rule_id, &charging_key)) {
    MLOG(MINFO) << "Updating used charging credit for Rule=" << rule_id
                 << " Rating Group=" << charging_key.rating_group
                 << " Service Identifier=" << charging_key.service_identifier;
    charging_pool_.add_used_credit(charging_key, used_tx, used_rx, update_criteria);
  }
  std::string monitoring_key;
  if (get_monitoring_key_for_rule_id(rule_id, &monitoring_key)) {
    MLOG(MINFO) << "Updating used monitoring credit for Rule=" << rule_id
                 << " Monitoring Key=" << monitoring_key;
    monitor_pool_.add_used_credit(monitoring_key, used_tx, used_rx, update_criteria);
  }
  auto session_level_key_p = monitor_pool_.get_session_level_key();
  if (
    session_level_key_p != nullptr && monitoring_key != *session_level_key_p) {
    // Update session level key if its different
    monitor_pool_.add_used_credit(*session_level_key_p, used_tx, used_rx, update_criteria);
  }
}

void SessionState::set_subscriber_quota_state(
    const magma::lte::SubscriberQuotaUpdate_Type state,
    SessionStateUpdateCriteria& update_criteria)
{
  update_criteria.updated_subscriber_quota_state = state;
  subscriber_quota_state_ = state;
}

bool SessionState::active_monitored_rules_exist()
{
  return total_monitored_rules_count() > 0;
}

void SessionState::get_updates_from_charging_pool(
  UpdateSessionRequest& update_request_out,
  std::vector<std::unique_ptr<ServiceAction>>* actions_out,
  SessionStateUpdateCriteria& update_criteria,
  const bool force_update)
{
  // charging updates
  std::vector<CreditUsage> charging_updates;
  charging_pool_.get_updates(
    imsi_, config_.ue_ipv4, static_rules_, &dynamic_rules_, &charging_updates, actions_out, update_criteria, force_update);
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
  SessionStateUpdateCriteria& update_criteria,
  const bool force_update)
{
  // monitor updates
  std::vector<UsageMonitorUpdate> monitor_updates;
  monitor_pool_.get_updates(
    imsi_, config_.ue_ipv4, static_rules_, &dynamic_rules_, &monitor_updates, actions_out, update_criteria, force_update);
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
    request_number_++;
  }
}

void SessionState::get_updates(
  UpdateSessionRequest& update_request_out,
  std::vector<std::unique_ptr<ServiceAction>>* actions_out,
  SessionStateUpdateCriteria& update_criteria,
  const bool force_update)
{
  if (curr_state_ != SESSION_ACTIVE) return;
  get_updates_from_charging_pool(update_request_out, actions_out, update_criteria, force_update);
  get_updates_from_monitor_pool(update_request_out, actions_out, update_criteria, force_update);
}

void SessionState::start_termination(SessionStateUpdateCriteria& update_criteria)
{
  curr_state_ = SESSION_TERMINATING_FLOW_ACTIVE;
}

void SessionState::set_termination_callback(
  std::function<void(SessionTerminateRequest)> on_termination_callback)
{
  on_termination_callback_ = on_termination_callback;
}

bool SessionState::can_complete_termination() const
{
  return curr_state_ == SESSION_TERMINATING_FLOW_DELETED;
}

void SessionState::mark_as_awaiting_termination(
  SessionStateUpdateCriteria& update_criteria)
{
  curr_state_ = SESSION_TERMINATION_SCHEDULED;
}

SubscriberQuotaUpdate_Type SessionState::get_subscriber_quota_state() const
{
  return subscriber_quota_state_;
}

void SessionState::complete_termination(
  SessionStateUpdateCriteria& update_criteria)
{
  if (curr_state_ == SESSION_TERMINATED) {
    // session is already terminated. Do nothing.
    return;
  }
  if (!can_complete_termination()) {
    MLOG(MERROR) << "Encountered unexpected state(" << curr_state_
                 << ") while terminating session for IMSI " << imsi_
                 << " and session id " << session_id_
                 << ". Forcefully terminating session.";
  }
  // mark entire session as terminated
  curr_state_ = SESSION_TERMINATED;
  SessionTerminateRequest termination;
  termination.set_sid(imsi_);
  termination.set_session_id(session_id_);
  termination.set_request_number(request_number_);
  termination.set_ue_ipv4(config_.ue_ipv4);
  termination.set_msisdn(config_.msisdn);
  termination.set_spgw_ipv4(config_.spgw_ipv4);
  termination.set_apn(config_.apn);
  termination.set_imei(config_.imei);
  termination.set_plmn_id(config_.plmn_id);
  termination.set_imsi_plmn_id(config_.imsi_plmn_id);
  termination.set_user_location(config_.user_location);
  termination.set_hardware_addr(config_.hardware_addr);
  termination.set_rat_type(config_.rat_type);
  fill_protos_tgpp_context(termination.mutable_tgpp_ctx());
  monitor_pool_.get_termination_updates(&termination, update_criteria);
  charging_pool_.get_termination_updates(&termination, update_criteria);
  try {
    on_termination_callback_(termination);
  } catch (std::bad_function_call&) {
    MLOG(MERROR) << "Missing termination callback function while terminating "
                    "session for IMSI "
                 << imsi_ << " and session id " << session_id_;
  }
}

void SessionState::complete_termination(
  SessionReporter& reporter,
  SessionStateUpdateCriteria& update_criteria)
{
  if (curr_state_ == SESSION_TERMINATED) {
    // session is already terminated. Do nothing.
    return;
  }
  if (!can_complete_termination()) {
    MLOG(MERROR) << "Encountered unexpected state(" << curr_state_
                 << ") while terminating session for IMSI " << imsi_
                 << " and session id " << session_id_
                 << ". Forcefully terminating session.";
  }
  // mark entire session as terminated
  curr_state_ = SESSION_TERMINATED;
  SessionTerminateRequest termination;
  termination.set_sid(imsi_);
  termination.set_session_id(session_id_);
  termination.set_request_number(request_number_);
  termination.set_ue_ipv4(config_.ue_ipv4);
  termination.set_msisdn(config_.msisdn);
  termination.set_spgw_ipv4(config_.spgw_ipv4);
  termination.set_apn(config_.apn);
  termination.set_imei(config_.imei);
  termination.set_plmn_id(config_.plmn_id);
  termination.set_imsi_plmn_id(config_.imsi_plmn_id);
  termination.set_user_location(config_.user_location);
  termination.set_hardware_addr(config_.hardware_addr);
  termination.set_rat_type(config_.rat_type);
  fill_protos_tgpp_context(termination.mutable_tgpp_ctx());
  monitor_pool_.get_termination_updates(&termination, update_criteria);
  charging_pool_.get_termination_updates(&termination, update_criteria);
  try {
    on_termination_callback_(termination);
  } catch (std::bad_function_call&) {
    on_termination_callback_ = [&reporter](SessionTerminateRequest term_req) {
      // report to cloud
      auto logging_cb =
          SessionReporter::get_terminate_logging_cb(term_req);
      reporter.report_terminate_session(term_req, logging_cb);
    };
    on_termination_callback_(termination);
  }
}

ChargingCreditPool& SessionState::get_charging_pool()
{
  return charging_pool_;
}

UsageMonitoringCreditPool& SessionState::get_monitor_pool()
{
  return monitor_pool_;
}

SessionState::TotalCreditUsage SessionState::get_total_credit_usage()
{
  // Collate unique charging/monitoring keys used by rules
  std::unordered_set<CreditKey, decltype(&ccHash), decltype(&ccEqual)> used_charging_keys
    (4, ccHash, ccEqual);
  std::unordered_set<std::string> used_monitoring_keys;

  std::vector<std::reference_wrapper<PolicyRuleBiMap>> bimaps{
    static_rules_, dynamic_rules_
  };

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
      if (should_track_monitoring_key) used_monitoring_keys.insert(monitoring_key);
    }
  }

  // Sum up usage
  TotalCreditUsage usage{
    .monitoring_tx = 0,
    .monitoring_rx = 0,
    .charging_tx = 0,
    .charging_rx = 0,
  };
  for (auto monitoring_key : used_monitoring_keys) {
    usage.monitoring_tx += get_monitor_pool().get_credit(monitoring_key, USED_TX);
    usage.monitoring_rx += get_monitor_pool().get_credit(monitoring_key, USED_RX);
  }
  for (auto charging_key : used_charging_keys) {
    usage.charging_tx += get_charging_pool().get_credit(charging_key, USED_TX);
    usage.charging_rx += get_charging_pool().get_credit(charging_key, USED_RX);
  }
  return usage;
}


bool SessionState::is_same_config(const Config& new_config) const
{
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

std::string SessionState::get_session_id() const
{
  return session_id_;
}

std::string SessionState::get_subscriber_ip_addr() const
{
  return config_.ue_ipv4;
}

std::string SessionState::get_mac_addr() const
{
  return config_.mac_addr;
}

std::string SessionState::get_msisdn() const
{
  return config_.msisdn;
}

std::string SessionState::get_apn() const
{
  return config_.apn;
}

void SessionState::set_config(const Config& config) {
  config_ = config;
}

bool SessionState::is_radius_cwf_session() const
{
  return (config_.rat_type == RATType::TGPP_WLAN);
}

std::string SessionState::get_radius_session_id() const
{
  return config_.radius_session_id;
}

void SessionState::get_session_info(SessionState::SessionInfo& info)
{
  info.imsi = imsi_;
  info.ip_addr = config_.ue_ipv4;
  get_dynamic_rules().get_rules(info.dynamic_rules);
  info.static_rules = active_static_rules_;
}

uint32_t SessionState::get_qci() const
{
  if (!config_.qos_info.enabled) {
    MLOG(MWARNING) << "QoS is not enabled.";
    return 0;
  }
  return config_.qos_info.qci;
}

uint32_t SessionState::get_bearer_id() const
{
  return config_.bearer_id;
}

bool SessionState::qos_enabled() const
{
  return config_.qos_info.enabled;
}

void SessionState::set_tgpp_context(
  const magma::lte::TgppContext& tgpp_context,
  SessionStateUpdateCriteria& update_criteria)
{
  update_criteria.updated_tgpp_context = tgpp_context;
  tgpp_context_ = tgpp_context;
}

void SessionState::fill_protos_tgpp_context(
  magma::lte::TgppContext* tgpp_context) const
{
  *tgpp_context = tgpp_context_;
}

uint32_t SessionState::get_request_number() {
  return request_number_;
}

void SessionState::increment_request_number(uint32_t incr) {
  request_number_ += incr;
}

bool SessionState::get_charging_key_for_rule_id(
  const std::string& rule_id,
  CreditKey* charging_key)
{
  // first check dynamic rules and then static rules
  if (dynamic_rules_.get_charging_key_for_rule_id(rule_id, charging_key)) {
    return true;
  }
  return static_rules_.get_charging_key_for_rule_id(rule_id, charging_key);
}

bool SessionState::get_monitoring_key_for_rule_id(
  const std::string& rule_id,
  std::string* monitoring_key)
{
  // first check dynamic rules and then static rules
  if (dynamic_rules_.get_monitoring_key_for_rule_id(rule_id, monitoring_key)) {
    return true;
  }
  return static_rules_.get_monitoring_key_for_rule_id(rule_id, monitoring_key);
}

bool SessionState::is_dynamic_rule_installed(const std::string& rule_id)
{
  auto _ = new PolicyRule();
  return dynamic_rules_.get_rule(rule_id, _);
}

bool SessionState::is_static_rule_installed(const std::string& rule_id)
{
  return std::find(
    active_static_rules_.begin(),
    active_static_rules_.end(),
    rule_id) != active_static_rules_.end();
}

void SessionState::insert_dynamic_rule(
  const PolicyRule& rule,
  SessionStateUpdateCriteria& update_criteria)
{
  if (is_dynamic_rule_installed(rule.id())) {
    return;
  }
  update_criteria.dynamic_rules_to_install.push_back(rule);
  dynamic_rules_.insert_rule(rule);
}

void SessionState::activate_static_rule(
  const std::string& rule_id,
  SessionStateUpdateCriteria& update_criteria)
{
  update_criteria.static_rules_to_install.insert(rule_id);
  active_static_rules_.push_back(rule_id);
}

bool SessionState::remove_dynamic_rule(
  const std::string& rule_id,
  PolicyRule *rule_out,
  SessionStateUpdateCriteria& update_criteria)
{
  bool removed = dynamic_rules_.remove_rule(rule_id, rule_out);
  if (removed) {
    update_criteria.dynamic_rules_to_uninstall.insert(rule_id);
  }
  return removed;
}

bool SessionState::deactivate_static_rule(
  const std::string& rule_id,
  SessionStateUpdateCriteria& update_criteria)
{
  auto it = std::find(active_static_rules_.begin(), active_static_rules_.end(),
                      rule_id);
  if (it == active_static_rules_.end()) {
    return false;
  }
  update_criteria.static_rules_to_uninstall.insert(rule_id);
  active_static_rules_.erase(it);
  return true;
}

DynamicRuleStore& SessionState::get_dynamic_rules()
{
  return dynamic_rules_;
}

uint32_t SessionState::total_monitored_rules_count()
{
  uint32_t monitored_dynamic_rules = dynamic_rules_.monitored_rules_count();
  uint32_t monitored_static_rules = 0;
  for (auto& rule_id : active_static_rules_)
  {
    std::string mkey; // ignore value
    auto is_monitored = static_rules_.get_monitoring_key_for_rule_id(
      rule_id,& mkey);
    if (is_monitored) {
      monitored_static_rules++;
    }
  }
  return monitored_dynamic_rules + monitored_static_rules;
}

uint32_t SessionState::get_credit_key_count()
{
  return charging_pool_.get_credit_key_count() + monitor_pool_.get_credit_key_count();
}

} // namespace magma
