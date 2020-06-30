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
#include "EnumToString.h"
#include "RuleStore.h"
#include "SessionState.h"
#include "StoredState.h"
#include "magma_logging.h"

namespace magma {

std::unique_ptr<SessionState> SessionState::unmarshal(
    const StoredSessionState& marshaled, StaticRuleStore& rule_store) {
  return std::make_unique<SessionState>(marshaled, rule_store);
}

StoredSessionState SessionState::marshal() {
  StoredSessionState marshaled{};

  marshaled.fsm_state              = curr_state_;
  marshaled.config                 = config_;
  marshaled.imsi                   = imsi_;
  marshaled.session_id             = session_id_;
  marshaled.core_session_id        = core_session_id_;
  marshaled.subscriber_quota_state = subscriber_quota_state_;
  marshaled.tgpp_context           = tgpp_context_;
  marshaled.request_number         = request_number_;
  marshaled.pending_event_triggers = pending_event_triggers_;
  marshaled.revalidation_time      = revalidation_time_;

  marshaled.monitor_map = StoredMonitorMap();
  for (auto &monitor_pair : monitor_map_) {
    StoredMonitor monitor{};
    monitor.credit = monitor_pair.second->credit.marshal();
    monitor.level = monitor_pair.second->level;
    marshaled.monitor_map[monitor_pair.first] = monitor;
  }

  if (session_level_key_ != nullptr) {
    marshaled.session_level_key = *session_level_key_;
  } else {
    marshaled.session_level_key = "";
  }

  marshaled.credit_map = StoredChargingCreditMap(4, &ccHash, &ccEqual);
  for (auto &credit_pair : credit_map_) {
    auto key = CreditKey();
    key.rating_group = credit_pair.first.rating_group;
    key.service_identifier = credit_pair.first.service_identifier;
    marshaled.credit_map[key] = credit_pair.second->marshal();
  }

  for (auto& rule_id : active_static_rules_) {
    marshaled.static_rule_ids.push_back(rule_id);
  }
  std::vector<PolicyRule> dynamic_rules;
  dynamic_rules_.get_rules(dynamic_rules);
  marshaled.dynamic_rules = std::move(dynamic_rules);

  std::vector<PolicyRule> gy_dynamic_rules;
  gy_dynamic_rules_.get_rules(gy_dynamic_rules);
  marshaled.gy_dynamic_rules = std::move(gy_dynamic_rules);

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
    : imsi_(marshaled.imsi),
      session_id_(marshaled.session_id),
      core_session_id_(marshaled.core_session_id),
      request_number_(marshaled.request_number),
      curr_state_(marshaled.fsm_state),
      config_(marshaled.config),
      subscriber_quota_state_(marshaled.subscriber_quota_state),
      tgpp_context_(marshaled.tgpp_context),
      static_rules_(rule_store),
      pending_event_triggers_(marshaled.pending_event_triggers),
      revalidation_time_(marshaled.revalidation_time),
      credit_map_(4, &ccHash, &ccEqual) {

  session_level_key_ =
      std::make_unique<std::string>(marshaled.session_level_key);
  for (auto it : marshaled.monitor_map) {
    Monitor monitor;
    monitor.credit = SessionCredit::unmarshal(it.second.credit, MONITORING);
    monitor.level = it.second.level;

    monitor_map_[it.first] = std::make_unique<Monitor>(monitor);
  }

  for (const auto &it : marshaled.credit_map) {
    credit_map_[it.first] =
      std::make_unique<ChargingGrant>(ChargingGrant::unmarshal(it.second));
  }

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
      // Request number set to 1, because request 0 is INIT call
      request_number_(1),
      curr_state_(SESSION_ACTIVE),
      config_(cfg),
      tgpp_context_(tgpp_context),
      static_rules_(rule_store),
      credit_map_(4, &ccHash, &ccEqual) {}

static CreditUsage
get_usage_proto_from_struct(const SessionCredit::Usage &usage_in,
                            CreditUsage::UpdateType proto_update_type,
                            const CreditKey &charging_key) {
  CreditUsage usage;
  usage.set_bytes_tx(usage_in.bytes_tx);
  usage.set_bytes_rx(usage_in.bytes_rx);
  usage.set_type(proto_update_type);
  charging_key.set_credit_usage(&usage);
  return usage;
}

static CreditUsage::UpdateType
convert_update_type_to_proto(CreditUpdateType update_type) {
  switch (update_type) {
    case CREDIT_QUOTA_EXHAUSTED:
      return CreditUsage::QUOTA_EXHAUSTED;
    case CREDIT_REAUTH_REQUIRED:
      return CreditUsage::REAUTH_REQUIRED;
    case CREDIT_VALIDITY_TIMER_EXPIRED:
      return CreditUsage::VALIDITY_TIMER_EXPIRED;
    default:
      MLOG(MERROR) << "Converting invalid update type " << update_type;
      return CreditUsage::QUOTA_EXHAUSTED;
  }
}

static UsageMonitorUpdate make_usage_monitor_update(
  const SessionCredit::Usage &usage_in, const std::string &monitoring_key,
  MonitoringLevel level) {
  UsageMonitorUpdate update;
  update.set_bytes_tx(usage_in.bytes_tx);
  update.set_bytes_rx(usage_in.bytes_rx);
  update.set_level(level);
  update.set_monitoring_key(monitoring_key);
  return update;
}

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

SessionCreditUpdateCriteria* SessionState::get_credit_uc(
  const CreditKey &key, SessionStateUpdateCriteria &uc) {
  if (uc.charging_credit_map.find(key) == uc.charging_credit_map.end()) {
    uc.charging_credit_map[key] = credit_map_[key]->credit.get_update_criteria();
  }
  return &(uc.charging_credit_map[key]);
}

void SessionState::add_rule_usage(
    const std::string& rule_id, uint64_t used_tx, uint64_t used_rx,
    SessionStateUpdateCriteria& update_criteria) {
  if (curr_state_ == SESSION_TERMINATING_AGGREGATING_STATS) {
    set_fsm_state(SESSION_TERMINATING_FLOW_ACTIVE,
                  update_criteria);
  }

  CreditKey charging_key;
  if (dynamic_rules_.get_charging_key_for_rule_id(rule_id, &charging_key) ||
      static_rules_.get_charging_key_for_rule_id(rule_id, &charging_key)) {
    MLOG(MINFO) << "Updating used charging credit for Rule=" << rule_id
                << " Rating Group=" << charging_key.rating_group
                << " Service Identifier=" << charging_key.service_identifier;
    auto it = credit_map_.find(charging_key);
    if (it != credit_map_.end()) {
       auto credit_uc = get_credit_uc(charging_key, update_criteria);
       it->second->credit.add_used_credit(used_tx, used_rx, *credit_uc);
    }
  }
  std::string monitoring_key;
  if (dynamic_rules_.get_monitoring_key_for_rule_id(rule_id, &monitoring_key) ||
      static_rules_.get_monitoring_key_for_rule_id(rule_id, &monitoring_key)) {
    MLOG(MINFO) << "Updating used monitoring credit for Rule=" << rule_id
                << " Monitoring Key=" << monitoring_key;
    add_to_monitor(monitoring_key, used_tx, used_rx, update_criteria);
  }
  auto session_level_key_p = get_session_level_key();
  if (session_level_key_p != nullptr &&
      monitoring_key != *session_level_key_p) {
    // Update session level key if its different
    add_to_monitor(*session_level_key_p, used_tx, used_rx, update_criteria);
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

void SessionState::get_monitor_updates(
    UpdateSessionRequest& update_request_out,
    std::vector<std::unique_ptr<ServiceAction>>* actions_out,
    SessionStateUpdateCriteria& update_criteria) {
  for (auto& monitor_pair : monitor_map_) {
    auto mkey = monitor_pair.first;
    auto& credit = monitor_pair.second->credit;
    auto credit_uc = get_monitor_uc(mkey, update_criteria);
    auto update_type = credit.get_update_type();
    if (update_type == CREDIT_NO_UPDATE) {
      continue;
    }
    MLOG(MDEBUG) << "Subscriber " << imsi_ << " monitoring key "
                 << mkey << " updating due to type "
                 << update_type;
    auto usage = credit.get_usage_for_reporting(*credit_uc);
    auto update = make_usage_monitor_update(
        usage, mkey, monitor_pair.second->level);
    auto new_req = update_request_out.mutable_usage_monitors()->Add();

    add_common_fields_to_usage_monitor_update(new_req);
    new_req->mutable_update()->CopyFrom(update);
    new_req->set_event_trigger(USAGE_REPORT);
    request_number_++;
    update_criteria.request_number_increment++;
  }
  // todo We should also handle other event triggers here too
  auto it = pending_event_triggers_.find(REVALIDATION_TIMEOUT);
  if (it != pending_event_triggers_.end() && it->second == READY) {
    auto new_req = update_request_out.mutable_usage_monitors()->Add();
    add_common_fields_to_usage_monitor_update(new_req);
    new_req->set_event_trigger(REVALIDATION_TIMEOUT);
    request_number_++;
    update_criteria.request_number_increment++;
    // todo we might want to make sure that the update went successfully before
    // clearing here
    remove_event_trigger(REVALIDATION_TIMEOUT, update_criteria);
  }
}

void SessionState::add_common_fields_to_usage_monitor_update(
  UsageMonitoringUpdateRequest* req) {
    req->set_session_id(session_id_);
    req->set_request_number(request_number_);
    req->set_sid(imsi_);
    req->set_ue_ipv4(config_.ue_ipv4);
    req->set_hardware_addr(config_.hardware_addr);
    req->set_rat_type(config_.rat_type);
    fill_protos_tgpp_context(req->mutable_tgpp_ctx());
}

void SessionState::get_updates(
    UpdateSessionRequest& update_request_out,
    std::vector<std::unique_ptr<ServiceAction>>* actions_out,
    SessionStateUpdateCriteria& update_criteria) {
  if (curr_state_ != SESSION_ACTIVE) return;
  get_charging_updates(update_request_out, actions_out, update_criteria);
  get_monitor_updates(update_request_out, actions_out, update_criteria);
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
  // gx monitors
  for (auto &credit_pair : monitor_map_) {
    auto credit_uc = get_monitor_uc(credit_pair.first, update_criteria);
    req.mutable_monitor_usages()->Add()->CopyFrom(
        make_usage_monitor_update(
            credit_pair.second->credit.get_all_unreported_usage_for_reporting(
                *credit_uc),
            credit_pair.first, credit_pair.second->level));
  }
  // gy credits
  for (auto &credit_pair : credit_map_) {
  auto credit_uc = get_credit_uc(credit_pair.first, update_criteria);
  req.mutable_credit_usages()->Add()->CopyFrom(
      get_usage_proto_from_struct(
          credit_pair.second->credit.get_all_unreported_usage_for_reporting(
              *credit_uc),
          CreditUsage::TERMINATED, credit_pair.first));
  }
  return req;
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
    usage.monitoring_tx += get_monitor(monitoring_key, USED_TX);
    usage.monitoring_rx += get_monitor(monitoring_key, USED_RX);
  }
  for (auto charging_key : used_charging_keys) {
    auto it = credit_map_.find(charging_key);
    if (it != credit_map_.end()) {
      usage.charging_tx += it->second->credit.get_credit(USED_TX);
      usage.charging_rx += it->second->credit.get_credit(USED_RX);
    }
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

SessionConfig SessionState::get_config() const {
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

bool SessionState::is_dynamic_rule_scheduled(const std::string& rule_id) {
  return scheduled_dynamic_rules_.get_rule(rule_id, NULL);
}

bool SessionState::is_static_rule_scheduled(const std::string& rule_id) {
  return scheduled_static_rules_.count(rule_id) == 1;
}

bool SessionState::is_dynamic_rule_installed(const std::string& rule_id) {
  return dynamic_rules_.get_rule(rule_id, NULL);
}

bool SessionState::is_gy_dynamic_rule_installed(const std::string& rule_id) {
  return gy_dynamic_rules_.get_rule(rule_id, NULL);
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
    const PolicyRule& rule, RuleLifetime& lifetime,
    SessionStateUpdateCriteria& update_criteria) {
  if (is_gy_dynamic_rule_installed(rule.id())) {
    MLOG(MDEBUG) << "Tried to insert "<< rule.id()
                 <<" (gy dynamic rule), but it already existed";
    return;
  }
  rule_lifetimes_[rule.id()] = lifetime;
  gy_dynamic_rules_.insert_rule(rule);
  update_criteria.gy_dynamic_rules_to_install.push_back(rule);
  update_criteria.new_rule_lifetimes[rule.id()] = lifetime;
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
    update_criteria.gy_dynamic_rules_to_uninstall.insert(rule_id);
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
      remove_dynamic_rule(rule_id, NULL, update_criteria);
    }
  }
  // Update scheduled dynamic rules
  dynamic_rule_ids.clear();
  scheduled_dynamic_rules_.get_rule_ids(dynamic_rule_ids);
  for (const std::string& rule_id : dynamic_rule_ids) {
    if (should_rule_be_active(rule_id, current_time)) {
      install_scheduled_dynamic_rule(rule_id, update_criteria);
    } else if (should_rule_be_deactivated(rule_id, current_time)) {
      remove_scheduled_dynamic_rule(rule_id, NULL, update_criteria);
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
  return credit_map_.size() + monitor_map_.size();
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

// Charging Credits
static FinalActionInfo get_final_action_info(
  const magma::lte::ChargingCredit &credit) {
  FinalActionInfo final_action_info;
  if (credit.is_final()) {
    final_action_info.final_action = credit.final_action();
    if (credit.final_action() == ChargingCredit_FinalAction_REDIRECT) {
      final_action_info.redirect_server = credit.redirect_server();
    }
  }
  return final_action_info;
}

bool SessionState::reset_reporting_charging_credit(
    const CreditKey &key, SessionStateUpdateCriteria &update_criteria) {
  auto it = credit_map_.find(key);
  if (it == credit_map_.end()) {
    MLOG(MERROR) << "Could not reset credit for IMSI" << imsi_
                 << " and charging key " << key << " because it wasn't found";
    return false;
  }
  auto credit_uc = get_credit_uc(key, update_criteria);
  it->second->credit.reset_reporting_credit(*credit_uc);
  return true;
}

bool SessionState::receive_charging_credit(
    const CreditUpdateResponse &update,
    SessionStateUpdateCriteria &update_criteria) {
  auto it = credit_map_.find(CreditKey(update));
  if (it == credit_map_.end()) {
    // new credit
    return init_charging_credit(update, update_criteria);
  }
  auto& grant = it->second;
  auto credit_uc = get_credit_uc(CreditKey(update), update_criteria);
  if (!update.success()) {
    // update unsuccessful, reset credit and return
    MLOG(MDEBUG) << "Rececive_Credit_Update: Unsuccessfull";
    grant->credit.mark_failure(update.result_code(), *credit_uc);
    return false;
  }
  const auto &gsu = update.credit().granted_units();
  MLOG(MDEBUG) << "Received charging credit of " << gsu.total().volume()
               << " total bytes, " << gsu.tx().volume() << " tx bytes, and "
               << gsu.rx().volume() << " rx bytes "
               << "for subscriber " << imsi_ << " rating group "
               << update.charging_key();
  grant->credit.receive_credit(
    gsu, update.credit().validity_time(), update.credit().is_final(),
    get_final_action_info(update.credit()), *credit_uc);
  return true;
}

bool SessionState::init_charging_credit(
    const CreditUpdateResponse &update,
    SessionStateUpdateCriteria &update_criteria) {
  if (!update.success()) {
    // init failed, don't track key
    MLOG(MERROR) << "Credit init failed for imsi " << imsi_
                 << " and charging key " << update.charging_key();
    return false;
  }
  MLOG(MINFO) << "Initialized a charging credit for imsi " << imsi_
              << " and charging key " << update.charging_key();

  auto charging_grant = std::make_unique<ChargingGrant>();
  charging_grant->credit =
    SessionCredit(CreditType::CHARGING, SERVICE_ENABLED, update.limit_type());

  SessionCreditUpdateCriteria credit_uc{};
  auto grant = update.credit();
  charging_grant->credit.receive_credit(
    grant.granted_units(), grant.validity_time(), grant.is_final(),
    get_final_action_info(update.credit()), credit_uc);

  update_criteria.charging_credit_to_install[CreditKey(update)] =
    charging_grant->marshal();
  credit_map_[CreditKey(update)] = std::move(charging_grant);
  return true;
}

uint64_t SessionState::get_charging_credit(const CreditKey &key,
                                        Bucket bucket) const {
  auto it = credit_map_.find(key);
  if (it == credit_map_.end()) {
    return 0;
  }
  return it->second->credit.get_credit(bucket);
}

ReAuthResult SessionState::reauth_key(const CreditKey &charging_key,
                               SessionStateUpdateCriteria &update_criteria) {
  auto it = credit_map_.find(charging_key);
  if (it != credit_map_.end()) {
    // if credit is already reporting, don't initiate update
    if (it->second->credit.is_reporting()) {
      return ReAuthResult::UPDATE_NOT_NEEDED;
    }
    auto uc = it->second->credit.get_update_criteria();
    it->second->credit.reauth(uc);
    update_criteria.charging_credit_map[charging_key] = uc;
    return ReAuthResult::UPDATE_INITIATED;
  }
  // charging_key cannot be found, initialize credit and engage reauth
  auto charging_grant = std::make_unique<ChargingGrant>();
  charging_grant->credit =
    SessionCredit(CreditType::CHARGING, SERVICE_DISABLED);
  SessionCreditUpdateCriteria _{};
  charging_grant->credit.reauth(_);
  update_criteria.charging_credit_to_install[charging_key] =
    charging_grant->marshal();
  credit_map_[charging_key] = std::move(charging_grant);
  return ReAuthResult::UPDATE_INITIATED;
}

ReAuthResult
SessionState::reauth_all(SessionStateUpdateCriteria &update_criteria) {
  auto res = ReAuthResult::UPDATE_NOT_NEEDED;
  for (auto &credit_pair : credit_map_) {
    // Only update credits that aren't reporting
    if (!credit_pair.second->credit.is_reporting()) {
      auto uc = credit_pair.second->credit.get_update_criteria();
      credit_pair.second->credit.reauth(uc);
      update_criteria.charging_credit_map[credit_pair.first] = uc;
      res = ReAuthResult::UPDATE_INITIATED;
    }
  }
  return res;
}

void SessionState::merge_charging_credit_update(
    const CreditKey &key, SessionCreditUpdateCriteria &credit_update) {
  auto it = credit_map_.find(key);
  if (it == credit_map_.end()) {
    return;
  }
  auto& credit = it->second->credit;
  credit.set_is_final_grant_and_final_action(credit_update.is_final, credit_update.final_action_info, credit_update);
  credit.set_reauth(credit_update.reauth_state, credit_update);
  credit.set_service_state(credit_update.service_state, credit_update);
  credit.set_expiry_time(credit_update.expiry_time, credit_update);
  credit.set_grant_tracking_type(credit_update.grant_tracking_type, credit_update);
  for (int i = USED_TX; i != MAX_VALUES; i++) {
    Bucket bucket = static_cast<Bucket>(i);
    credit.add_credit(
      credit_update.bucket_deltas.find(bucket)->second, bucket, credit_update);
  }
}

void SessionState::set_charging_credit(
    const CreditKey &key, ChargingGrant charging_grant,
    SessionStateUpdateCriteria &uc) {
  credit_map_[key] = std::make_unique<ChargingGrant>(charging_grant);
  uc.charging_credit_to_install[key] = credit_map_[key]->marshal();
}

CreditUsageUpdate SessionState::make_credit_usage_update_req(CreditUsage& usage) const {
  CreditUsageUpdate req;
  req.set_session_id(session_id_);
  req.set_request_number(request_number_);
  req.set_sid(imsi_);
  req.set_msisdn(config_.msisdn);
  req.set_ue_ipv4(config_.ue_ipv4);
  req.set_spgw_ipv4(config_.spgw_ipv4);
  req.set_apn(config_.apn);
  req.set_imei(config_.imei);
  req.set_plmn_id(config_.plmn_id);
  req.set_imsi_plmn_id(config_.imsi_plmn_id);
  req.set_user_location(config_.user_location);
  req.set_hardware_addr(config_.hardware_addr);
  req.set_rat_type(config_.rat_type);
  fill_protos_tgpp_context(req.mutable_tgpp_ctx());
  req.mutable_usage()->CopyFrom(usage);
  return req;
}

void SessionState::get_charging_updates(
    UpdateSessionRequest& update_request_out,
    std::vector<std::unique_ptr<ServiceAction>>* actions_out,
    SessionStateUpdateCriteria& uc) {
  for (auto &credit_pair : credit_map_) {
    auto& key = credit_pair.first;
    auto& grant = credit_pair.second;
    auto credit_uc = get_credit_uc(key, uc);

    auto action_type = grant->credit.get_action(*credit_uc);
    auto action = std::make_unique<ServiceAction>(action_type);
    switch (action_type) {
      case CONTINUE_SERVICE:
        {
          auto update_type = grant->credit.get_update_type();
          if (update_type == CREDIT_NO_UPDATE) {
            break;
          }
          // Create Update struct
          MLOG(MDEBUG) << "Subscriber " << imsi_ << " rating group "
                   << key << " updating due to type "
                   << update_type;
          auto usage = grant->credit.get_usage_for_reporting(*credit_uc);
          auto p_update_type = convert_update_type_to_proto(update_type);
          auto update = get_usage_proto_from_struct(usage, p_update_type, key);
          auto credit_req = make_credit_usage_update_req(update);
          update_request_out.mutable_updates()->Add()->CopyFrom(credit_req);
          request_number_++;
          uc.request_number_increment++;
        }
        break;
      case REDIRECT:
        if (credit_uc->service_state == SERVICE_REDIRECTED) {
          MLOG(MDEBUG) << "Redirection already activated.";
          continue;
        }
        grant->credit.set_service_state(SERVICE_REDIRECTED, *credit_uc);
        action->set_redirect_server(grant->credit.get_redirect_server());
      case TERMINATE_SERVICE:
      case ACTIVATE_SERVICE:
      case RESTRICT_ACCESS:
        MLOG(MDEBUG) << "Subscriber " << imsi_ << " rating group "
                     << key << " action type " << action_type;
        action->set_credit_key(key);
        action->set_imsi(imsi_);
        action->set_ip_addr(config_.ue_ipv4);
        static_rules_.get_rule_ids_for_charging_key(
          key, *action->get_mutable_rule_ids());
        dynamic_rules_.get_rule_definitions_for_charging_key(
          key, *action->get_mutable_rule_definitions());
        actions_out->push_back(std::move(action));
        break;
    }
  }
}

// Monitors
bool SessionState::receive_monitor(
    const UsageMonitoringUpdateResponse &update,
    SessionStateUpdateCriteria &update_criteria) {
  if (update.success() &&
      update.credit().level() == MonitoringLevel::SESSION_LEVEL) {
    update_session_level_key(update, update_criteria);
  }
  auto it = monitor_map_.find(update.credit().monitoring_key());
  if (it == monitor_map_.end()) {
    // new credit
    return init_new_monitor(update, update_criteria);
  }
  auto credit_uc =
      get_monitor_uc(update.credit().monitoring_key(), update_criteria);
  if (!update.success()) {
    it->second->credit.mark_failure(update.result_code(), *credit_uc);
    return false;
  }
  const auto &gsu = update.credit().granted_units();
  MLOG(MDEBUG) << "Received monitor credit of " << gsu.total().volume()
               << " total bytes, " << gsu.tx().volume() << " tx bytes, and "
               << gsu.rx().volume() << " rx bytes "
               << "for subscriber " << imsi_ << " monitoring key "
               << update.credit().monitoring_key();
  FinalActionInfo final_action_info;
  it->second->credit.receive_credit(
    gsu, 0, false, final_action_info, *credit_uc);
  if (update.credit().action() == UsageMonitoringCredit::DISABLE) {
    monitor_map_.erase(update.credit().monitoring_key());
  }
  return true;
}

void SessionState::merge_monitor_updates(
  const std::string &key, SessionCreditUpdateCriteria &update) {
  auto it = monitor_map_.find(key);
  if (it == monitor_map_.end()) {
    return;
  }

  it->second->credit.set_is_final_grant_and_final_action(
        update.is_final, update.final_action_info, update);
  it->second->credit.set_reauth(update.reauth_state, update);
  it->second->credit.set_service_state(update.service_state, update);
  it->second->credit.set_expiry_time(update.expiry_time, update);
  for (int i = USED_TX; i != MAX_VALUES; i++) {
    Bucket bucket = static_cast<Bucket>(i);
    it->second->credit.add_credit(
        update.bucket_deltas.find(bucket)->second, bucket, update);
  }
}

uint64_t SessionState::get_monitor(const std::string &key, Bucket bucket) const {
  auto it = monitor_map_.find(key);
  if (it == monitor_map_.end()) {
    return 0;
  }
  return it->second->credit.get_credit(bucket);
}

bool SessionState::add_to_monitor(
  const std::string &key, uint64_t used_tx,
  uint64_t used_rx, SessionStateUpdateCriteria &uc) {
  auto it = monitor_map_.find(key);
  if (it == monitor_map_.end()) {
    return false;
  }
  auto credit_uc = get_monitor_uc(key, uc);
  it->second->credit.add_used_credit(used_tx, used_rx, *credit_uc);
  return true;
}

void SessionState::set_monitor(
  const std::string &key,
  std::unique_ptr<Monitor> monitor,
  SessionStateUpdateCriteria &update_criteria) {
  update_criteria.monitor_credit_to_install[key] = monitor->marshal();
  monitor_map_[key] = std::move(monitor);
}

bool SessionState::reset_reporting_monitor(
  const std::string &key, SessionStateUpdateCriteria &update_criteria) {
    auto it = monitor_map_.find(key);
  if (it == monitor_map_.end()) {
    MLOG(MERROR) << "Could not reset credit for IMSI" << imsi_
                 << " and monitoring key " << key << " because it wasn't found";
    return false;
  }
  auto credit_uc = get_monitor_uc(key, update_criteria);
  it->second->credit.reset_reporting_credit(*credit_uc);
  return true;
}

std::unique_ptr<std::string> SessionState::get_session_level_key() const{
  if (session_level_key_ == nullptr){
    return nullptr;
  }
  return std::make_unique<std::string>(*session_level_key_);
}

bool SessionState::init_new_monitor(
    const UsageMonitoringUpdateResponse &update,
    SessionStateUpdateCriteria &update_criteria) {
  if (!update.success()) {
    MLOG(MERROR) << "Monitoring init failed for imsi " << imsi_
                 << " and monitoring key " << update.credit().monitoring_key();
    return false;
  }
  if (update.credit().action() == UsageMonitoringCredit::DISABLE) {
    MLOG(MWARNING) << "Monitoring init has action disabled for subscriber "
                   << imsi_ << " and monitoring key "
                   << update.credit().monitoring_key();
    return false;
  }
  MLOG(MDEBUG) << "Initialized a monitoring credit for imsi" << imsi_
               << " and monitoring key " << update.credit().monitoring_key();
  auto monitor = std::make_unique<Monitor>();
  monitor->level = update.credit().level();
  // validity time and final units not used for monitors
  auto _ = SessionCreditUpdateCriteria{};
  FinalActionInfo final_action_info;
  auto gsu = update.credit().granted_units();
  monitor->credit.receive_credit(gsu, 0, false, final_action_info, _);

  update_criteria.monitor_credit_to_install[update.credit().monitoring_key()] =
      monitor->marshal();
  monitor_map_[update.credit().monitoring_key()] = std::move(monitor);
  return true;
}

void SessionState::update_session_level_key(
    const UsageMonitoringUpdateResponse &update,
    SessionStateUpdateCriteria &update_criteria) {
  const auto &new_key = update.credit().monitoring_key();
  if (session_level_key_ != nullptr && *session_level_key_ != new_key) {
    MLOG(MWARNING) << "Session level monitoring key already exists, updating";
  }
  if (update.credit().action() == UsageMonitoringCredit::DISABLE) {
    session_level_key_ = nullptr;
    // TODO: set in UpdateCriteria
  } else {
    session_level_key_ = std::make_unique<std::string>(new_key);
    // TODO: set in UpdateCriteria
  }
}

SessionCreditUpdateCriteria* SessionState::get_monitor_uc(
  const std::string &key, SessionStateUpdateCriteria &uc) {
  if (uc.monitor_credit_map.find(key) == uc.monitor_credit_map.end()) {
    uc.monitor_credit_map[key] =
        monitor_map_[key]->credit.get_update_criteria();
  }
  return &(uc.monitor_credit_map[key]);
}

// Event Triggers
void SessionState::add_new_event_trigger(
  magma::lte::EventTrigger trigger,
  SessionStateUpdateCriteria& update_criteria) {
    MLOG(MINFO) << "Event Trigger " << trigger << " is pending for "
                << session_id_;
    set_event_trigger(trigger, PENDING, update_criteria);
}

void SessionState::mark_event_trigger_as_triggered(
  magma::lte::EventTrigger trigger,
  SessionStateUpdateCriteria& update_criteria) {
    auto it = pending_event_triggers_.find(trigger);
    if (it == pending_event_triggers_.end() ||
        pending_event_triggers_[trigger] != PENDING) {
      MLOG(MWARNING) << "Event Trigger " << trigger << " requested to be "
                     << "triggered is not pending for " << session_id_;
    }
    MLOG(MINFO) << "Event Trigger " << trigger << " is ready to update for "
                << session_id_;
    set_event_trigger(trigger, READY, update_criteria);
}

void SessionState::remove_event_trigger(
  magma::lte::EventTrigger trigger,
  SessionStateUpdateCriteria& update_criteria) {
    MLOG(MINFO) << "Event Trigger " << trigger << " is removed for "
                << session_id_;
    pending_event_triggers_.erase(trigger);
    set_event_trigger(trigger, CLEARED, update_criteria);
}

void SessionState::set_event_trigger(
  magma::lte::EventTrigger trigger,
  const EventTriggerState value,
  SessionStateUpdateCriteria& update_criteria) {
    pending_event_triggers_[trigger] = value;
    update_criteria.is_pending_event_triggers_updated = true;
    update_criteria.pending_event_triggers[trigger] = value;
}

void SessionState::set_revalidation_time(
  const google::protobuf::Timestamp& time,
  SessionStateUpdateCriteria& update_criteria) {
  revalidation_time_ = time;
  update_criteria.revalidation_time = time;
}

bool SessionState::is_credit_state_redirected(
    const CreditKey &charging_key) const {
  auto it = credit_map_.find(charging_key);
  if (it == credit_map_.end()) {
    return false;
  }
  return it->second->credit.is_service_redirected();
}
}  // namespace magma
