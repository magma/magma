/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include "CreditPool.h"
#include "StoredState.h"
#include <limits>

#include "magma_logging.h"

namespace magma {

std::unique_ptr<ChargingCreditPool>
ChargingCreditPool::unmarshal(const StoredChargingCreditPool &marshaled) {
  auto pool = std::make_unique<ChargingCreditPool>(marshaled.imsi);
  for (const auto &it : marshaled.credit_map) {
    pool->credit_map_[it.first] = SessionCredit::unmarshal(it.second, CHARGING);
  }
  return pool;
}

StoredChargingCreditPool ChargingCreditPool::marshal() {
  StoredChargingCreditPool marshaled{};
  marshaled.imsi = imsi_;
  auto credit_map =
      std::unordered_map<CreditKey, StoredSessionCredit, decltype(&ccHash),
                         decltype(&ccEqual)>(4, &ccHash, &ccEqual);
  for (auto &credit_pair : credit_map_) {
    auto key = CreditKey();
    key.rating_group = credit_pair.first.rating_group;
    key.service_identifier = credit_pair.first.service_identifier;
    key.use_sid = credit_pair.first.use_sid;
    credit_map[key] = credit_pair.second->marshal();
  }
  marshaled.credit_map = credit_map;
  return marshaled;
}

ChargingCreditPool::ChargingCreditPool(const std::string &imsi)
    : credit_map_(4, &ccHash, &ccEqual), imsi_(imsi) {}

bool ChargingCreditPool::add_used_credit(
    const CreditKey &key, uint64_t used_tx, uint64_t used_rx,
    SessionStateUpdateCriteria &update_criteria) {
  auto it = credit_map_.find(key);
  if (it == credit_map_.end()) {
    return false;
  }
  auto credit_uc = get_credit_update(key, update_criteria);
  it->second->add_used_credit(used_tx, used_rx, *credit_uc);
  return true;
}

bool ChargingCreditPool::reset_reporting_credit(
    const CreditKey &key, SessionStateUpdateCriteria &update_criteria) {
  auto it = credit_map_.find(key);
  if (it == credit_map_.end()) {
    MLOG(MERROR) << "Could not reset credit for IMSI" << imsi_
                 << " and charging key " << key << " because it wasn't found";
    return false;
  }
  auto credit_uc = get_credit_update(key, update_criteria);
  it->second->reset_reporting_credit(*credit_uc);
  return true;
}

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
  }
  MLOG(MERROR) << "Converting invalid update type " << update_type;
  return CreditUsage::QUOTA_EXHAUSTED;
}

void ChargingCreditPool::populate_output_actions(
    std::string imsi, std::string ip_addr, CreditKey key,
    StaticRuleStore &static_rules, DynamicRuleStore *dynamic_rules,
    std::unique_ptr<ServiceAction> &action,
    std::vector<std::unique_ptr<ServiceAction>> *actions_out) const {
  action->set_imsi(imsi);
  action->set_ip_addr(ip_addr);
  static_rules.get_rule_ids_for_charging_key(key,
                                             *action->get_mutable_rule_ids());
  dynamic_rules->get_rule_definitions_for_charging_key(
      key, *action->get_mutable_rule_definitions());
  actions_out->push_back(std::move(action));
}

void UsageMonitoringCreditPool::populate_output_actions(
    std::string imsi, std::string ip_addr, std::string key,
    StaticRuleStore &static_rules, DynamicRuleStore *dynamic_rules,
    std::unique_ptr<ServiceAction> &action,
    std::vector<std::unique_ptr<ServiceAction>> *actions_out) const {
  action->set_imsi(imsi);
  action->set_ip_addr(ip_addr);
  static_rules.get_rule_ids_for_monitoring_key(key,
                                               *action->get_mutable_rule_ids());
  dynamic_rules->get_rule_definitions_for_monitoring_key(
      key, *action->get_mutable_rule_definitions());
  actions_out->push_back(std::move(action));
}

void ChargingCreditPool::get_updates(
    std::string imsi, std::string ip_addr, StaticRuleStore &static_rules,
    DynamicRuleStore *dynamic_rules, std::vector<CreditUsage> *updates_out,
    std::vector<std::unique_ptr<ServiceAction>> *actions_out,
    SessionStateUpdateCriteria &update_criteria, const bool force_update) {
  for (auto &credit_pair : credit_map_) {
    auto &credit = *(credit_pair.second);
    auto credit_uc = get_credit_update(credit_pair.first, update_criteria);
    auto action_type = credit.get_action(*credit_uc);
    if (action_type != CONTINUE_SERVICE) {
      MLOG(MDEBUG) << "Subscriber " << imsi_ << " rating group "
                   << credit_pair.first << " action type " << action_type;
      auto action = std::make_unique<ServiceAction>(action_type);
      if (action_type == REDIRECT) {
        action->set_credit_key(credit_pair.first);
        action->set_redirect_server(credit.get_redirect_server());
      }
      populate_output_actions(imsi, ip_addr, credit_pair.first, static_rules,
                              dynamic_rules, action, actions_out);
    } else {
      auto update_type = credit.get_update_type();
      if (update_type != CREDIT_NO_UPDATE || force_update) {
        MLOG(MDEBUG) << "Subscriber " << imsi_ << " rating group "
                     << credit_pair.first << " updating due to type "
                     << update_type;
        updates_out->push_back(get_usage_proto_from_struct(
            credit.get_usage_for_reporting(*credit_uc),
            convert_update_type_to_proto(update_type), credit_pair.first));
      }
    }
  }
}

bool ChargingCreditPool::get_termination_updates(
    SessionTerminateRequest *termination_out,
    SessionStateUpdateCriteria &update_criteria) {
  for (auto &credit_pair : credit_map_) {
    auto credit_uc = get_credit_update(credit_pair.first, update_criteria);
    termination_out->mutable_credit_usages()->Add()->CopyFrom(
        get_usage_proto_from_struct(
            credit_pair.second->get_all_unreported_usage_for_reporting(
                *credit_uc),
            CreditUsage::TERMINATED, credit_pair.first));
  }
  return true;
}

static uint64_t get_granted_units(const CreditUnit &unit,
                                  uint64_t default_val) {
  return unit.is_valid() ? unit.volume() : default_val;
}

static FinalActionInfo get_final_action_info(const ChargingCredit &credit) {
  FinalActionInfo final_action_info;
  if (credit.is_final()) {
    final_action_info.final_action = credit.final_action();
    if (credit.final_action() == ChargingCredit_FinalAction_REDIRECT) {
      final_action_info.redirect_server = credit.redirect_server();
    }
  }

  return final_action_info;
}

static void receive_charging_credit_with_default(
    SessionCredit &credit, const GrantedUnits &gsu, uint64_t default_volume,
    const ChargingCredit &charging_credit,
    SessionCreditUpdateCriteria &update_criteria) {
  uint64_t total = get_granted_units(gsu.total(), default_volume);
  uint64_t tx = get_granted_units(gsu.tx(), default_volume);
  uint64_t rx = get_granted_units(gsu.rx(), default_volume);

  credit.receive_credit(total, tx, rx, charging_credit.validity_time(),
                        charging_credit.is_final(),
                        get_final_action_info(charging_credit),
                        update_criteria);
}

bool ChargingCreditPool::init_new_credit(
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
  // unless defined, volume is defined as the maximum possible value
  // uint64_t default_volume = std::numeric_limits<uint64_t>::max();
  /*
   * Setting it to 0 otherwise it will lead to data transfer even if there
   * in no GSUs present in Gy:CCA-I
   */
  uint64_t default_volume = 0;
  std::unique_ptr<SessionCredit> credit;
  if (update.limit_type() == CreditUpdateResponse::FINITE) {
    credit = std::make_unique<SessionCredit>(CreditType::CHARGING);
  } else {
    credit = std::make_unique<SessionCredit>(CreditType::CHARGING,
                                             SERVICE_ENABLED, true);
  }
  SessionCreditUpdateCriteria credit_uc{};
  receive_charging_credit_with_default(*credit, update.credit().granted_units(),
                                       default_volume, update.credit(),
                                       credit_uc);
  StoredSessionCredit stored_credit = credit->marshal();
  update_criteria.charging_credit_to_install[CreditKey(update)] =
      credit->marshal();
  credit_map_[CreditKey(update)] = std::move(credit);
  return true;
}

bool ChargingCreditPool::receive_credit(
    const CreditUpdateResponse &update,
    SessionStateUpdateCriteria &update_criteria) {
  auto it = credit_map_.find(CreditKey(update));
  if (it == credit_map_.end()) {
    // new credit
    return init_new_credit(update, update_criteria);
  }
  auto credit_uc = get_credit_update(CreditKey(update), update_criteria);
  if (!update.success()) {
    // update unsuccessful, reset credit and return
    MLOG(MDEBUG) << "Rececive_Credit_Update: Unsuccessfull";
    it->second->mark_failure(update.result_code(), *credit_uc);
    return false;
  }
  const auto &gsu = update.credit().granted_units();
  MLOG(MDEBUG) << "Received charging credit of " << gsu.total().volume()
               << " total bytes, " << gsu.tx().volume() << " tx bytes, and "
               << gsu.rx().volume() << " rx bytes "
               << "for subscriber " << imsi_ << " rating group "
               << update.charging_key();
  uint64_t default_volume = 0; // default to not increasing credit
  receive_charging_credit_with_default(*(it->second), gsu, default_volume,
                                       update.credit(), *credit_uc);
  return true;
}

uint64_t ChargingCreditPool::get_credit(const CreditKey &key,
                                        Bucket bucket) const {
  auto it = credit_map_.find(key);
  if (it == credit_map_.end()) {
    return 0;
  }
  return it->second->get_credit(bucket);
}

void ChargingCreditPool::add_credit(
    const CreditKey &key, std::unique_ptr<SessionCredit> credit,
    SessionStateUpdateCriteria &update_criteria) {
  update_criteria.charging_credit_to_install[key] = credit->marshal();
  credit_map_[key] = std::move(credit);
}

SessionCreditUpdateCriteria *
ChargingCreditPool::get_credit_update(const CreditKey &key,
                                      SessionStateUpdateCriteria &uc) {
  if (uc.charging_credit_map.find(key) == uc.charging_credit_map.end()) {
    uc.charging_credit_map[key] = credit_map_[key]->get_update_criteria();
  }
  return &(uc.charging_credit_map[key]);
}

void ChargingCreditPool::merge_credit_update(
    const CreditKey &key, SessionCreditUpdateCriteria &credit_update) {
  auto it = credit_map_.find(key);
  if (it == credit_map_.end()) {
    return;
  }
  it->second->set_is_final_grant_and_final_action(
        credit_update.is_final, credit_update.final_action_info, credit_update);
  it->second->set_reauth(credit_update.reauth_state, credit_update);
  it->second->set_service_state(credit_update.service_state, credit_update);
  it->second->set_expiry_time(credit_update.expiry_time, credit_update);
  for (int i = USED_TX; i != MAX_VALUES; i++) {
    Bucket bucket = static_cast<Bucket>(i);
    it->second->add_credit(credit_update.bucket_deltas.find(bucket)->second,
                           bucket, credit_update);
  }
}

uint32_t ChargingCreditPool::get_credit_key_count() const {
  return credit_map_.size();
}

ReAuthResult
ChargingCreditPool::reauth_key(const CreditKey &charging_key,
                               SessionStateUpdateCriteria &update_criteria) {
  auto it = credit_map_.find(charging_key);
  if (it != credit_map_.end()) {
    // if credit is already reporting, don't initiate update
    if (it->second->is_reporting()) {
      return ReAuthResult::UPDATE_NOT_NEEDED;
    }
    auto uc = it->second->get_update_criteria();
    it->second->reauth(uc);
    update_criteria.charging_credit_map[charging_key] = uc;
    return ReAuthResult::UPDATE_INITIATED;
  }
  // charging_key cannot be found, initialize credit and engage reauth
  auto credit =
      std::make_unique<SessionCredit>(CreditType::CHARGING, SERVICE_DISABLED);
  SessionCreditUpdateCriteria _{};
  credit->reauth(_);
  update_criteria.charging_credit_to_install[charging_key] = credit->marshal();
  credit_map_[charging_key] = std::move(credit);
  return ReAuthResult::UPDATE_INITIATED;
}

ReAuthResult
ChargingCreditPool::reauth_all(SessionStateUpdateCriteria &update_criteria) {
  auto res = ReAuthResult::UPDATE_NOT_NEEDED;
  for (auto &credit_pair : credit_map_) {
    // Only update credits that aren't reporting
    if (!credit_pair.second->is_reporting()) {
      auto uc = credit_pair.second->get_update_criteria();
      credit_pair.second->reauth(uc);
      update_criteria.charging_credit_map[credit_pair.first] = uc;
      res = ReAuthResult::UPDATE_INITIATED;
    }
  }
  return res;
}

StoredMonitor UsageMonitoringCreditPool::marshal_monitor(
    std::unique_ptr<UsageMonitoringCreditPool::Monitor> &monitor) {
  StoredMonitor marshaled{};
  marshaled.credit = monitor->credit.marshal();
  marshaled.level = monitor->level;
  return marshaled;
}

std::unique_ptr<UsageMonitoringCreditPool::Monitor>
UsageMonitoringCreditPool::unmarshal_monitor(const StoredMonitor &marshaled) {
  UsageMonitoringCreditPool::Monitor monitor;
  monitor.credit = *SessionCredit::unmarshal(marshaled.credit, MONITORING);
  monitor.level = marshaled.level;
  return std::make_unique<UsageMonitoringCreditPool::Monitor>(monitor);
}

std::unique_ptr<UsageMonitoringCreditPool> UsageMonitoringCreditPool::unmarshal(
    const StoredUsageMonitoringCreditPool &marshaled) {
  auto pool = std::make_unique<UsageMonitoringCreditPool>(marshaled.imsi);
  pool->session_level_key_ =
      std::make_unique<std::string>(marshaled.session_level_key);
  for (auto it : marshaled.monitor_map) {
    Monitor monitor;
    monitor.credit = *SessionCredit::unmarshal(it.second.credit, MONITORING);
    monitor.level = it.second.level;

    pool->monitor_map_[it.first] = std::make_unique<Monitor>(monitor);
  }
  return pool;
}

StoredUsageMonitoringCreditPool UsageMonitoringCreditPool::marshal() {
  StoredUsageMonitoringCreditPool marshaled{};
  marshaled.imsi = imsi_;
  if (session_level_key_ != nullptr) {
    marshaled.session_level_key = *session_level_key_;
  } else {
    marshaled.session_level_key = "";
  }
  for (auto &it : monitor_map_) {
    StoredMonitor monitor{};
    monitor.credit = it.second->credit.marshal();
    monitor.level = it.second->level;
    marshaled.monitor_map[it.first] = monitor;
  }
  return marshaled;
}

UsageMonitoringCreditPool::UsageMonitoringCreditPool(const std::string &imsi)
    : imsi_(imsi), session_level_key_(nullptr) {}

static void receive_monitoring_credit_with_default(
    SessionCredit &credit, const GrantedUnits &gsu, uint64_t default_volume,
    SessionCreditUpdateCriteria &update_criteria) {
  uint64_t total = get_granted_units(gsu.total(), default_volume);
  uint64_t tx = get_granted_units(gsu.tx(), default_volume);
  uint64_t rx = get_granted_units(gsu.rx(), default_volume);

  FinalActionInfo final_action_info;

  credit.receive_credit(total, tx, rx, 0, false, final_action_info,
                        update_criteria);
}

bool UsageMonitoringCreditPool::add_used_credit(
    const std::string &key, uint64_t used_tx, uint64_t used_rx,
    SessionStateUpdateCriteria &update_criteria) {
  auto it = monitor_map_.find(key);
  if (it == monitor_map_.end()) {
    return false;
  }
  auto credit_uc = get_credit_update(key, update_criteria);
  it->second->credit.add_used_credit(used_tx, used_rx, *credit_uc);
  return true;
}

bool UsageMonitoringCreditPool::reset_reporting_credit(
    const std::string &key, SessionStateUpdateCriteria &update_criteria) {
  auto it = monitor_map_.find(key);
  if (it == monitor_map_.end()) {
    MLOG(MERROR) << "Could not reset credit for IMSI" << imsi_
                 << " and monitoring key " << key << " because it wasn't found";
    return false;
  }
  auto credit_uc = get_credit_update(key, update_criteria);
  it->second->credit.reset_reporting_credit(*credit_uc);
  return true;
}

static UsageMonitorUpdate
get_monitor_update_from_struct(const SessionCredit::Usage &usage_in,
                               const std::string &monitoring_key,
                               MonitoringLevel level) {
  UsageMonitorUpdate update;
  update.set_bytes_tx(usage_in.bytes_tx);
  update.set_bytes_rx(usage_in.bytes_rx);
  update.set_level(level);
  update.set_monitoring_key(monitoring_key);
  return update;
}

void UsageMonitoringCreditPool::get_updates(
    std::string imsi, std::string ip_addr, StaticRuleStore &static_rules,
    DynamicRuleStore *dynamic_rules,
    std::vector<UsageMonitorUpdate> *updates_out,
    std::vector<std::unique_ptr<ServiceAction>> *actions_out,
    SessionStateUpdateCriteria &update_criteria, const bool force_update) {
  for (auto &monitor_pair : monitor_map_) {
    auto &credit = monitor_pair.second->credit;
    auto credit_uc = get_credit_update(monitor_pair.first, update_criteria);
    auto action_type = credit.get_action(*credit_uc);
    if (action_type != CONTINUE_SERVICE) {
      auto action = std::make_unique<ServiceAction>(action_type);
      populate_output_actions(imsi, ip_addr, monitor_pair.first, static_rules,
                              dynamic_rules, action, actions_out);
    }
    auto update_type = credit.get_update_type();
    if (update_type != CREDIT_NO_UPDATE || force_update) {
      MLOG(MDEBUG) << "Subscriber " << imsi_ << " monitoring key "
                   << monitor_pair.first << " updating due to type "
                   << update_type;
      updates_out->push_back(get_monitor_update_from_struct(
          credit.get_usage_for_reporting(*credit_uc), monitor_pair.first,
          monitor_pair.second->level));
    }
  }
}

bool UsageMonitoringCreditPool::get_termination_updates(
    SessionTerminateRequest *termination_out,
    SessionStateUpdateCriteria &update_criteria) {
  for (auto &credit_pair : monitor_map_) {
    auto credit_uc = get_credit_update(credit_pair.first, update_criteria);
    termination_out->mutable_monitor_usages()->Add()->CopyFrom(
        get_monitor_update_from_struct(
            credit_pair.second->credit.get_all_unreported_usage_for_reporting(
                *credit_uc),
            credit_pair.first, credit_pair.second->level));
  }
}

void UsageMonitoringCreditPool::update_session_level_key(
    const UsageMonitoringUpdateResponse &update,
    SessionStateUpdateCriteria &update_criteria) {
  if (!update.success()) {
    return;
  }
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

bool UsageMonitoringCreditPool::init_new_credit(
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
  auto monitor = std::make_unique<UsageMonitoringCreditPool::Monitor>();
  monitor->level = update.credit().level();
  // validity time and final units not used for monitors
  // unless defined, volume is defined as the maximum possible value
  uint64_t default_volume = std::numeric_limits<uint64_t>::max();
  auto _ = SessionCreditUpdateCriteria{};
  receive_monitoring_credit_with_default(
      monitor->credit, update.credit().granted_units(), default_volume, _);
  update_criteria.monitor_credit_to_install[update.credit().monitoring_key()] =
      marshal_monitor(monitor);
  monitor_map_[update.credit().monitoring_key()] = std::move(monitor);
  return true;
}

bool UsageMonitoringCreditPool::receive_credit(
    const UsageMonitoringUpdateResponse &update,
    SessionStateUpdateCriteria &update_criteria) {
  if (update.credit().level() == MonitoringLevel::SESSION_LEVEL) {
    update_session_level_key(update, update_criteria);
  }
  auto it = monitor_map_.find(update.credit().monitoring_key());
  if (it == monitor_map_.end()) {
    // new credit
    return init_new_credit(update, update_criteria);
  }
  auto credit_uc =
      get_credit_update(update.credit().monitoring_key(), update_criteria);
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
  uint64_t default_volume = 0;
  receive_monitoring_credit_with_default(it->second->credit,
                                         update.credit().granted_units(),
                                         default_volume, *credit_uc);
  if (update.credit().action() == UsageMonitoringCredit::DISABLE) {
    monitor_map_.erase(update.credit().monitoring_key());
  }
  return true;
}

uint64_t UsageMonitoringCreditPool::get_credit(const std::string &key,
                                               Bucket bucket) const {
  auto it = monitor_map_.find(key);
  if (it == monitor_map_.end()) {
    return 0;
  }
  return it->second->credit.get_credit(bucket);
}

void UsageMonitoringCreditPool::add_monitor(
    const std::string &key,
    std::unique_ptr<UsageMonitoringCreditPool::Monitor> monitor,
    SessionStateUpdateCriteria &update_criteria) {
  update_criteria.monitor_credit_to_install[key] = marshal_monitor(monitor);
  monitor_map_[key] = std::move(monitor);
}

SessionCreditUpdateCriteria *
UsageMonitoringCreditPool::get_credit_update(const std::string &key,
                                             SessionStateUpdateCriteria &uc) {
  if (uc.monitor_credit_map.find(key) == uc.monitor_credit_map.end()) {
    uc.monitor_credit_map[key] =
        monitor_map_[key]->credit.get_update_criteria();
  }
  return &(uc.monitor_credit_map[key]);
}

void UsageMonitoringCreditPool::merge_credit_update(
    const std::string &key, SessionCreditUpdateCriteria &credit_update) {
  auto it = monitor_map_.find(key);
  if (it == monitor_map_.end()) {
    return;
  }

  it->second->credit.set_is_final_grant_and_final_action(
        credit_update.is_final, credit_update.final_action_info, credit_update);
  it->second->credit.set_reauth(credit_update.reauth_state, credit_update);
  it->second->credit.set_service_state(credit_update.service_state,
                                       credit_update);
  it->second->credit.set_expiry_time(credit_update.expiry_time, credit_update);
  for (int i = USED_TX; i != MAX_VALUES; i++) {
    Bucket bucket = static_cast<Bucket>(i);
    it->second->credit.add_credit(
        credit_update.bucket_deltas.find(bucket)->second, bucket,
        credit_update);
  }
}

uint32_t UsageMonitoringCreditPool::get_credit_key_count() const {
  return monitor_map_.size();
}

std::unique_ptr<std::string>
UsageMonitoringCreditPool::get_session_level_key() {
  if (session_level_key_ == nullptr)
    return nullptr;
  return std::make_unique<std::string>(*session_level_key_);
}

} // namespace magma
