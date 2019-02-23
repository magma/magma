/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include <lte/protos/session_manager.grpc.pb.h>

namespace magma {
using namespace lte;

void create_rule_record(
    const std::string& imsi,
    const std::string& rule_id,
    uint64_t bytes_rx,
    uint64_t bytes_tx,
    RuleRecord* rule_record);

void create_charging_credit(uint64_t volume, ChargingCredit* credit);

void create_update_response(
    const std::string& imsi,
    uint32_t charging_key,
    uint64_t volume,
    CreditUpdateResponse* response);

void create_update_response(
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

void create_usage_update(
    const std::string& imsi,
    uint32_t charging_key,
    uint64_t bytes_rx,
    uint64_t bytes_tx,
    CreditUsage::UpdateType type,
    CreditUsageUpdate* update);

}
