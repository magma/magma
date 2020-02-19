/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "ProtobufCreators.h"

namespace magma {

void create_rule_record(
  const std::string& imsi,
  const std::string& rule_id,
  uint64_t bytes_rx,
  uint64_t bytes_tx,
  RuleRecord* rule_record)
{
  rule_record->set_sid(imsi);
  rule_record->set_rule_id(rule_id);
  rule_record->set_bytes_rx(bytes_rx);
  rule_record->set_bytes_tx(bytes_tx);
}

void create_charging_credit(
  uint64_t volume,
  bool is_final,
  ChargingCredit* credit)
{
  credit->mutable_granted_units()->mutable_total()->set_volume(volume);
  credit->mutable_granted_units()->mutable_total()->set_is_valid(true);
  credit->set_type(ChargingCredit::BYTES);
  credit->set_is_final(is_final);
}

// defaults to not final credit
void create_credit_update_response(
  const std::string& imsi,
  uint32_t charging_key,
  uint64_t volume,
  CreditUpdateResponse* response)
{
  create_credit_update_response(imsi, charging_key, volume, false, response);
}

void create_credit_update_response(
  const std::string& imsi,
  uint32_t charging_key,
  uint64_t volume,
  bool is_final,
  CreditUpdateResponse* response)
{
  create_charging_credit(volume, is_final, response->mutable_credit());
  response->set_success(true);
  response->set_sid(imsi);
  response->set_charging_key(charging_key);
  response->set_type(CreditUpdateResponse::UPDATE);
}

void create_usage_update(
  const std::string& imsi,
  uint32_t charging_key,
  uint64_t bytes_rx,
  uint64_t bytes_tx,
  CreditUsage::UpdateType type,
  CreditUsageUpdate* update)
{
  auto usage = update->mutable_usage();
  update->set_sid(imsi);
  usage->set_charging_key(charging_key);
  usage->set_bytes_rx(bytes_rx);
  usage->set_bytes_tx(bytes_tx);
  usage->set_type(type);
}

void create_monitor_credit(
  const std::string& m_key,
  MonitoringLevel level,
  uint64_t volume,
  UsageMonitoringCredit* credit)
{
  if (volume == 0) {
    credit->set_action(UsageMonitoringCredit::DISABLE);
  } else {
    credit->set_action(UsageMonitoringCredit::CONTINUE);
  }
  credit->mutable_granted_units()->mutable_total()->set_volume(volume);
  credit->mutable_granted_units()->mutable_total()->set_is_valid(true);
  credit->set_level(level);
  credit->set_monitoring_key(m_key);
}

void create_monitor_update_response(
  const std::string& imsi,
  const std::string& m_key,
  MonitoringLevel level,
  uint64_t volume,
  UsageMonitoringUpdateResponse* response)
{
  std::vector<EventTrigger> event_triggers;
  create_monitor_update_response(
    imsi, m_key, level, volume, event_triggers, 0, response);
}

void create_monitor_update_response(
  const std::string& imsi,
  const std::string& m_key,
  MonitoringLevel level,
  uint64_t volume,
  const std::vector<EventTrigger>& event_triggers,
  const uint64_t revalidation_time_unix_ts,
  UsageMonitoringUpdateResponse* response)
{
  create_monitor_credit(m_key, level, volume, response->mutable_credit());
  response->set_success(true);
  response->set_sid(imsi);
  for (const auto& event_trigger : event_triggers) {
    response->add_event_triggers(event_trigger);
  }
  response->mutable_revalidation_time()->set_seconds(revalidation_time_unix_ts);
}

void create_policy_reauth_request(
  const std::string& session_id,
  const std::string& imsi,
  const std::vector<std::string>& rules_to_remove,
  const std::vector<StaticRuleInstall>& rules_to_install,
  const std::vector<DynamicRuleInstall>& dynamic_rules_to_install,
  const std::vector<EventTrigger>& event_triggers,
  const uint64_t revalidation_time_unix_ts,
  const std::vector<UsageMonitoringCredit>& usage_monitoring_credits,
  PolicyReAuthRequest* request)
{
  request->set_session_id(session_id);
  request->set_imsi(imsi);
  for (const auto& rule_id : rules_to_remove) {
    request->add_rules_to_remove(rule_id);
  }
  auto req_rules_to_install = request->mutable_rules_to_install();
  for (const auto& static_rule_to_install : rules_to_install) {
    req_rules_to_install->Add()->CopyFrom(static_rule_to_install);
  }
  auto req_dynamic_rules_to_install = request->mutable_dynamic_rules_to_install();
  for (const auto& dynamic_rule_to_install : dynamic_rules_to_install) {
    req_dynamic_rules_to_install->Add()->CopyFrom(dynamic_rule_to_install);
  }
  for (const auto& event_trigger : event_triggers) {
    request->add_event_triggers(event_trigger);
  }
  request->mutable_revalidation_time()->set_seconds(revalidation_time_unix_ts);
  auto req_credits = request->mutable_usage_monitoring_credits();
  for (const auto& credit : usage_monitoring_credits) {
    req_credits->Add()->CopyFrom(credit);
  }
}

void create_tgpp_context(
  const std::string& gx_dest_host,
  const std::string& gy_dest_host,
  TgppContext* context)
{
  context->set_gx_dest_host(gx_dest_host);
  context->set_gy_dest_host(gy_dest_host);
}

void create_subscriber_quota_update(
  const std::string& imsi,
  const std::string& ue_mac_addr,
  const SubscriberQuotaUpdate_Type state,
  SubscriberQuotaUpdate* update)
{
  auto sid = update->mutable_sid();
  sid->set_id(imsi);
  update->set_mac_addr(ue_mac_addr);
  update->set_update_type(state);
}

} // namespace magma
