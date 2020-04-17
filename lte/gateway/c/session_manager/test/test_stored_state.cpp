/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <memory>

#include <glog/logging.h>
#include <gtest/gtest.h>

#include "StoredState.h"
#include "magma_logging.h"

using ::testing::Test;

namespace magma {

class StoredStateTest : public ::testing::Test {
protected:
  QoSInfo get_stored_qos_info() {
    QoSInfo stored;
    stored.enabled = true;
    stored.qci = 123;
    return stored;
  }

  SessionConfig get_stored_session_config() {
    SessionConfig stored;
    stored.ue_ipv4 = "192.168.0.1";
    stored.spgw_ipv4 = "192.168.0.2";
    stored.msisdn = "a";
    stored.apn = "b";
    stored.imei = "c";
    stored.plmn_id = "d";
    stored.imsi_plmn_id = "e";
    stored.user_location = "f";
    stored.rat_type = RATType::TGPP_WLAN;
    stored.mac_addr = "g";      // MAC Address for WLAN
    stored.hardware_addr = "h"; // MAC Address for WLAN (binary)
    stored.radius_session_id = "i";
    stored.bearer_id = 321;
    stored.qos_info = get_stored_qos_info();
    return stored;
  }

  StoredRedirectServer get_stored_redirect_server() {
    StoredRedirectServer stored;
    stored.redirect_address_type = RedirectServer_RedirectAddressType::
        RedirectServer_RedirectAddressType_IPV6;
    stored.redirect_server_address = "redirect_server_address";
    return stored;
  }

  FinalActionInfo get_stored_final_action_info() {
    FinalActionInfo stored;
    stored.final_action =
        ChargingCredit_FinalAction::ChargingCredit_FinalAction_REDIRECT;
    stored.redirect_server.set_redirect_address_type(
        RedirectServer_RedirectAddressType::
            RedirectServer_RedirectAddressType_IPV6);
    stored.redirect_server.set_redirect_server_address(
        "redirect_server_address");
    return stored;
  }

  StoredSessionCredit get_stored_session_credit() {
    StoredSessionCredit stored;

    stored.reporting = true;
    stored.is_final = true;
    stored.unlimited_quota = true;

    stored.final_action_info.final_action =
        ChargingCredit_FinalAction::ChargingCredit_FinalAction_REDIRECT;
    stored.final_action_info.redirect_server.set_redirect_address_type(
        RedirectServer_RedirectAddressType::
            RedirectServer_RedirectAddressType_IPV6);
    stored.final_action_info.redirect_server.set_redirect_server_address(
        "redirect_server_address");

    stored.reauth_state = REAUTH_REQUIRED;
    stored.service_state = SERVICE_NEEDS_ACTIVATION;

    stored.expiry_time = 0;

    stored.buckets[USED_TX] = 12345;
    stored.buckets[ALLOWED_TOTAL] = 54321;

    stored.usage_reporting_limit = 4444;
    return stored;
  };

  StoredMonitor get_stored_monitor() {
    StoredMonitor stored;
    stored.credit = get_stored_session_credit();
    stored.level = MonitoringLevel::PCC_RULE_LEVEL;
    return stored;
  }

  StoredChargingCreditPool get_stored_charging_credit_pool() {
    StoredChargingCreditPool stored;
    stored.imsi = "IMSI1";
    auto credit_map =
        std::unordered_map<CreditKey, StoredSessionCredit, decltype(&ccHash),
                           decltype(&ccEqual)>(4, &ccHash, &ccEqual);
    credit_map[CreditKey(1, 2)] = get_stored_session_credit();
    stored.credit_map = credit_map;
    return stored;
  }

  StoredUsageMonitoringCreditPool get_stored_monitor_pool() {
    StoredUsageMonitoringCreditPool stored;
    stored.imsi = "IMSI1";
    stored.session_level_key = "session_level_key";
    stored.monitor_map["mk1"] = get_stored_monitor();
    return stored;
  }

  StoredSessionState get_stored_session() {
    StoredSessionState stored;

    stored.config = get_stored_session_config();
    stored.charging_pool = get_stored_charging_credit_pool();
    stored.monitor_pool = get_stored_monitor_pool();
    stored.imsi = "IMSI1";
    stored.session_id = "session_id";
    stored.core_session_id = "core_session_id";
    stored.subscriber_quota_state = SubscriberQuotaUpdate_Type_VALID_QUOTA;

    magma::lte::TgppContext tgpp_context;
    tgpp_context.set_gx_dest_host("gx");
    tgpp_context.set_gy_dest_host("gy");
    stored.tgpp_context = tgpp_context;

    stored.request_number = 1;

    return stored;
  }
};

TEST_F(StoredStateTest, test_stored_qos_info) {
  auto stored = get_stored_qos_info();

  auto serialized = serialize_stored_qos_info(stored);
  auto deserialized = deserialize_stored_qos_info(serialized);

  EXPECT_EQ(deserialized.enabled, true);
  EXPECT_EQ(deserialized.qci, 123);
}

TEST_F(StoredStateTest, test_stored_session_config) {
  auto stored = get_stored_session_config();

  std::string serialized = serialize_stored_session_config(stored);
  SessionConfig deserialized = deserialize_stored_session_config(serialized);

  EXPECT_EQ(deserialized.ue_ipv4, "192.168.0.1");
  EXPECT_EQ(deserialized.spgw_ipv4, "192.168.0.2");
  EXPECT_EQ(deserialized.msisdn, "a");
  EXPECT_EQ(deserialized.apn, "b");
  EXPECT_EQ(deserialized.imei, "c");
  EXPECT_EQ(deserialized.plmn_id, "d");
  EXPECT_EQ(deserialized.imsi_plmn_id, "e");
  EXPECT_EQ(deserialized.user_location, "f");
  EXPECT_EQ(deserialized.rat_type, RATType::TGPP_WLAN);
  EXPECT_EQ(deserialized.mac_addr, "g");
  EXPECT_EQ(deserialized.hardware_addr, "h");
  EXPECT_EQ(deserialized.radius_session_id, "i");
  EXPECT_EQ(deserialized.bearer_id, 321);
  EXPECT_EQ(deserialized.qos_info.enabled, true);
  EXPECT_EQ(deserialized.qos_info.qci, 123);
}

TEST_F(StoredStateTest, test_stored_redirect_server) {
  auto stored = get_stored_redirect_server();

  auto serialized = serialize_stored_redirect_server(stored);
  auto deserialized = deserialize_stored_redirect_server(serialized);

  EXPECT_EQ(deserialized.redirect_address_type,
            RedirectServer_RedirectAddressType::
                RedirectServer_RedirectAddressType_IPV6);
  EXPECT_EQ(deserialized.redirect_server_address, "redirect_server_address");
}

TEST_F(StoredStateTest, test_stored_final_action_info) {
  auto stored = get_stored_final_action_info();

  auto serialized = serialize_stored_final_action_info(stored);
  auto deserialized = deserialize_stored_final_action_info(serialized);

  EXPECT_EQ(deserialized.final_action,
            ChargingCredit_FinalAction::ChargingCredit_FinalAction_REDIRECT);
  EXPECT_EQ(deserialized.redirect_server.redirect_address_type(),
            RedirectServer_RedirectAddressType::
                RedirectServer_RedirectAddressType_IPV6);
  EXPECT_EQ(deserialized.redirect_server.redirect_server_address(),
            "redirect_server_address");
}

TEST_F(StoredStateTest, test_stored_session_credit) {
  auto stored = get_stored_session_credit();

  auto serialized = serialize_stored_session_credit(stored);
  auto deserialized = deserialize_stored_session_credit(serialized);

  EXPECT_EQ(deserialized.reporting, true);
  EXPECT_EQ(deserialized.is_final, true);
  EXPECT_EQ(deserialized.unlimited_quota, true);

  EXPECT_EQ(deserialized.final_action_info.final_action,
            ChargingCredit_FinalAction::ChargingCredit_FinalAction_REDIRECT);
  EXPECT_EQ(
      deserialized.final_action_info.redirect_server.redirect_address_type(),
      RedirectServer_RedirectAddressType::
          RedirectServer_RedirectAddressType_IPV6);
  EXPECT_EQ(
      deserialized.final_action_info.redirect_server.redirect_server_address(),
      "redirect_server_address");

  EXPECT_EQ(deserialized.reauth_state, REAUTH_REQUIRED);
  EXPECT_EQ(deserialized.service_state, SERVICE_NEEDS_ACTIVATION);

  EXPECT_EQ(deserialized.expiry_time, 0);
  EXPECT_EQ(deserialized.buckets[USED_TX], 12345);
  EXPECT_EQ(deserialized.buckets[ALLOWED_TOTAL], 54321);

  EXPECT_EQ(deserialized.usage_reporting_limit, 4444);
}

TEST_F(StoredStateTest, test_stored_monitor) {
  auto stored = get_stored_monitor();

  auto serialized = serialize_stored_monitor(stored);
  auto deserialized = deserialize_stored_monitor(serialized);

  EXPECT_EQ(deserialized.credit.reporting, true);
  EXPECT_EQ(deserialized.credit.is_final, true);
  EXPECT_EQ(deserialized.credit.unlimited_quota, true);
  EXPECT_EQ(deserialized.credit.final_action_info.final_action,
            ChargingCredit_FinalAction::ChargingCredit_FinalAction_REDIRECT);
  EXPECT_EQ(deserialized.credit.final_action_info.redirect_server
                .redirect_address_type(),
            RedirectServer_RedirectAddressType::
                RedirectServer_RedirectAddressType_IPV6);
  EXPECT_EQ(deserialized.credit.final_action_info.redirect_server
                .redirect_server_address(),
            "redirect_server_address");
  EXPECT_EQ(deserialized.credit.reauth_state, REAUTH_REQUIRED);
  EXPECT_EQ(deserialized.credit.service_state, SERVICE_NEEDS_ACTIVATION);
  EXPECT_EQ(deserialized.credit.expiry_time, 0);
  EXPECT_EQ(deserialized.credit.buckets[USED_TX], 12345);
  EXPECT_EQ(deserialized.credit.buckets[ALLOWED_TOTAL], 54321);
  EXPECT_EQ(deserialized.credit.usage_reporting_limit, 4444);
  EXPECT_EQ(deserialized.level, MonitoringLevel::PCC_RULE_LEVEL);
}

TEST_F(StoredStateTest, test_stored_charging_credit_pool) {
  auto stored = get_stored_charging_credit_pool();

  auto serialized = serialize_stored_charging_credit_pool(stored);
  auto deserialized = deserialize_stored_charging_credit_pool(serialized);

  auto stored_credit = deserialized.credit_map[CreditKey(1, 2)];
  EXPECT_EQ(stored_credit.reporting, true);
  EXPECT_EQ(stored_credit.is_final, true);
  EXPECT_EQ(stored_credit.unlimited_quota, true);
  EXPECT_EQ(stored_credit.final_action_info.final_action,
            ChargingCredit_FinalAction::ChargingCredit_FinalAction_REDIRECT);
  EXPECT_EQ(
      stored_credit.final_action_info.redirect_server.redirect_address_type(),
      RedirectServer_RedirectAddressType::
          RedirectServer_RedirectAddressType_IPV6);
  EXPECT_EQ(
      stored_credit.final_action_info.redirect_server.redirect_server_address(),
      "redirect_server_address");
  EXPECT_EQ(stored_credit.reauth_state, REAUTH_REQUIRED);
  EXPECT_EQ(stored_credit.service_state, SERVICE_NEEDS_ACTIVATION);
  EXPECT_EQ(stored_credit.expiry_time, 0);
  EXPECT_EQ(stored_credit.buckets[USED_TX], 12345);
  EXPECT_EQ(stored_credit.buckets[ALLOWED_TOTAL], 54321);
  EXPECT_EQ(stored_credit.usage_reporting_limit, 4444);
}

TEST_F(StoredStateTest, test_stored_monitor_pool) {
  auto stored = get_stored_monitor_pool();

  auto serialized = serialize_stored_usage_monitoring_pool(stored);
  auto deserialized = deserialize_stored_usage_monitoring_pool(serialized);

  EXPECT_EQ(deserialized.imsi, "IMSI1");
  EXPECT_EQ(deserialized.session_level_key, "session_level_key");

  auto stored_monitor = deserialized.monitor_map["mk1"];
  EXPECT_EQ(stored_monitor.credit.reporting, true);
  EXPECT_EQ(stored_monitor.credit.is_final, true);
  EXPECT_EQ(stored_monitor.credit.unlimited_quota, true);
  EXPECT_EQ(stored_monitor.credit.final_action_info.final_action,
            ChargingCredit_FinalAction::ChargingCredit_FinalAction_REDIRECT);
  EXPECT_EQ(stored_monitor.credit.final_action_info.redirect_server
                .redirect_address_type(),
            RedirectServer_RedirectAddressType::
                RedirectServer_RedirectAddressType_IPV6);
  EXPECT_EQ(stored_monitor.credit.final_action_info.redirect_server
                .redirect_server_address(),
            "redirect_server_address");
  EXPECT_EQ(stored_monitor.credit.reauth_state, REAUTH_REQUIRED);
  EXPECT_EQ(stored_monitor.credit.service_state, SERVICE_NEEDS_ACTIVATION);
  EXPECT_EQ(stored_monitor.credit.expiry_time, 0);
  EXPECT_EQ(stored_monitor.credit.buckets[USED_TX], 12345);
  EXPECT_EQ(stored_monitor.credit.buckets[ALLOWED_TOTAL], 54321);
  EXPECT_EQ(stored_monitor.credit.usage_reporting_limit, 4444);
  EXPECT_EQ(stored_monitor.level, MonitoringLevel::PCC_RULE_LEVEL);
}

TEST_F(StoredStateTest, test_stored_session) {
  auto stored = get_stored_session();

  auto serialized = serialize_stored_session(stored);
  auto deserialized = deserialize_stored_session(serialized);

  auto config = deserialized.config;
  EXPECT_EQ(config.ue_ipv4, "192.168.0.1");
  EXPECT_EQ(config.spgw_ipv4, "192.168.0.2");
  EXPECT_EQ(config.msisdn, "a");
  EXPECT_EQ(config.apn, "b");
  EXPECT_EQ(config.imei, "c");
  EXPECT_EQ(config.plmn_id, "d");
  EXPECT_EQ(config.imsi_plmn_id, "e");
  EXPECT_EQ(config.user_location, "f");
  EXPECT_EQ(config.rat_type, RATType::TGPP_WLAN);
  EXPECT_EQ(config.mac_addr, "g");
  EXPECT_EQ(config.hardware_addr, "h");
  EXPECT_EQ(config.radius_session_id, "i");
  EXPECT_EQ(config.bearer_id, 321);
  EXPECT_EQ(config.qos_info.enabled, true);
  EXPECT_EQ(config.qos_info.qci, 123);

  auto charging_pool = deserialized.charging_pool;
  auto stored_credit = charging_pool.credit_map[CreditKey(1, 2)];
  EXPECT_EQ(stored_credit.reporting, true);
  EXPECT_EQ(stored_credit.is_final, true);
  EXPECT_EQ(stored_credit.unlimited_quota, true);
  EXPECT_EQ(stored_credit.final_action_info.final_action,
            ChargingCredit_FinalAction::ChargingCredit_FinalAction_REDIRECT);
  EXPECT_EQ(
      stored_credit.final_action_info.redirect_server.redirect_address_type(),
      RedirectServer_RedirectAddressType::
          RedirectServer_RedirectAddressType_IPV6);
  EXPECT_EQ(
      stored_credit.final_action_info.redirect_server.redirect_server_address(),
      "redirect_server_address");
  EXPECT_EQ(stored_credit.reauth_state, REAUTH_REQUIRED);
  EXPECT_EQ(stored_credit.service_state, SERVICE_NEEDS_ACTIVATION);
  EXPECT_EQ(stored_credit.expiry_time, 0);
  EXPECT_EQ(stored_credit.buckets[USED_TX], 12345);
  EXPECT_EQ(stored_credit.buckets[ALLOWED_TOTAL], 54321);
  EXPECT_EQ(stored_credit.usage_reporting_limit, 4444);

  auto monitor_pool = deserialized.monitor_pool;
  EXPECT_EQ(monitor_pool.imsi, "IMSI1");
  EXPECT_EQ(monitor_pool.session_level_key, "session_level_key");
  auto stored_monitor = monitor_pool.monitor_map["mk1"];
  EXPECT_EQ(stored_monitor.credit.reporting, true);
  EXPECT_EQ(stored_monitor.credit.is_final, true);
  EXPECT_EQ(stored_monitor.credit.unlimited_quota, true);
  EXPECT_EQ(stored_monitor.credit.final_action_info.final_action,
            ChargingCredit_FinalAction::ChargingCredit_FinalAction_REDIRECT);
  EXPECT_EQ(stored_monitor.credit.final_action_info.redirect_server
                .redirect_address_type(),
            RedirectServer_RedirectAddressType::
                RedirectServer_RedirectAddressType_IPV6);
  EXPECT_EQ(stored_monitor.credit.final_action_info.redirect_server
                .redirect_server_address(),
            "redirect_server_address");
  EXPECT_EQ(stored_monitor.credit.reauth_state, REAUTH_REQUIRED);
  EXPECT_EQ(stored_monitor.credit.service_state, SERVICE_NEEDS_ACTIVATION);
  EXPECT_EQ(stored_monitor.credit.expiry_time, 0);
  EXPECT_EQ(stored_monitor.credit.buckets[USED_TX], 12345);
  EXPECT_EQ(stored_monitor.credit.buckets[ALLOWED_TOTAL], 54321);
  EXPECT_EQ(stored_monitor.credit.usage_reporting_limit, 4444);
  EXPECT_EQ(stored_monitor.level, MonitoringLevel::PCC_RULE_LEVEL);

  EXPECT_EQ(stored.imsi, "IMSI1");
  EXPECT_EQ(stored.session_id, "session_id");
  EXPECT_EQ(stored.core_session_id, "core_session_id");
  EXPECT_EQ(stored.subscriber_quota_state,
            SubscriberQuotaUpdate_Type_VALID_QUOTA);

  EXPECT_EQ(stored.tgpp_context.gx_dest_host(), "gx");
  EXPECT_EQ(stored.tgpp_context.gy_dest_host(), "gy");

  EXPECT_EQ(stored.request_number, 1);
}

int main(int argc, char **argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

} // namespace magma
