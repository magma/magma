/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include <lte/protos/mconfig/mconfigs.pb.h>
#include <lte/protos/session_manager.grpc.pb.h>
#include <lte/protos/pipelined.grpc.pb.h>

namespace magma {
using namespace lte;

void create_rule_record(
  const std::string& imsi,
  const std::string& rule_id,
  uint64_t bytes_rx,
  uint64_t bytes_tx,
  RuleRecord* rule_record);

void create_charging_credit(uint64_t volume, ChargingCredit* credit);

void create_credit_update_response(
  const std::string& imsi,
  uint32_t charging_key,
  uint64_t volume,
  CreditUpdateResponse* response);

void create_credit_update_response(
  const std::string& imsi,
  uint32_t charging_key,
  uint64_t volume,
  bool is_final,
  CreditUpdateResponse* response);

void create_monitor_credit(
  const std::string& m_key,
  MonitoringLevel level,
  uint64_t volume,
  UsageMonitoringCredit* response);

void create_monitor_update_response(
  const std::string& imsi,
  const std::string& m_key,
  MonitoringLevel level,
  uint64_t volume,
  UsageMonitoringUpdateResponse* response);

void create_monitor_update_response(
  const std::string& imsi,
  const std::string& m_key,
  MonitoringLevel level,
  uint64_t volume,
  const std::vector<EventTrigger>& event_triggers,
  const uint64_t revalidation_time_unix_ts,
  UsageMonitoringUpdateResponse* response);

void create_usage_update(
  const std::string& imsi,
  uint32_t charging_key,
  uint64_t bytes_rx,
  uint64_t bytes_tx,
  CreditUsage::UpdateType type,
  CreditUsageUpdate* update);

void create_policy_reauth_request(
  const std::string& session_id,
  const std::string& imsi,
  const std::vector<std::string>& rules_to_remove,
  const std::vector<StaticRuleInstall>& rules_to_install,
  const std::vector<DynamicRuleInstall>& dynamic_rules_to_install,
  const std::vector<EventTrigger>& event_triggers,
  const uint64_t revalidation_time_unix_ts,
  const std::vector<UsageMonitoringCredit>& usage_monitoring_credits,
  PolicyReAuthRequest* request);

void create_tgpp_context(
  const std::string& gx_dest_host,
  const std::string& gy_dest_host,
  TgppContext* context);

void create_subscriber_quota_update(
  const std::string& imsi,
  const std::string& ue_mac_addr,
  const SubscriberQuotaUpdate_Type state,
  SubscriberQuotaUpdate* update);

void create_session_create_response(
  const std::string& imsi,
  const std::string& monitoring_key,
  std::vector<std::string>& static_rules,
  CreateSessionResponse* response);

magma::mconfig::SessionD get_default_mconfig();
} // namespace magma
