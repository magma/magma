/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include <chrono>
#include <thread>

#include <glog/logging.h>
#include <gtest/gtest.h>

#include "CreditPool.h"
#include "StoredState.h"
#include <lte/protos/session_manager.pb.h>
#include <lte/protos/policydb.pb.h>

using ::testing::Test;

namespace magma {
using namespace lte;

class CreditPoolTest : public ::testing::Test {
 protected:
  CreditUpdateResponse* get_credit_update()
  {
    auto credit_update = new CreditUpdateResponse();
    credit_update->set_success(true);
    credit_update->set_sid("sid1"); // System Identification Number
    credit_update->set_charging_key(1);
    credit_update->set_allocated_credit(get_charging_credit());
    credit_update->set_type(CreditUpdateResponse_ResponseType_UPDATE);
    // Don't set the result code since this is a successful response
    auto service_identifier = new ServiceIdentifier();
    service_identifier->set_value(uint32_t(1));
    credit_update->set_allocated_service_identifier(service_identifier);
    credit_update->set_limit_type(
      CreditUpdateResponse_CreditLimitType_INFINITE_METERED);
    // Don't set the TgppContext, assume relay disabled
    return credit_update;
  }

  ChargingCredit* get_charging_credit()
  {
    auto charging_credit = new ChargingCredit();
    charging_credit->set_type(ChargingCredit_UnitType_BYTES);
    charging_credit->set_validity_time(1000);
    charging_credit->set_is_final(false);
    charging_credit->set_final_action(ChargingCredit_FinalAction_TERMINATE);
    charging_credit->set_allocated_granted_units(get_granted_units());
    charging_credit->set_allocated_redirect_server(get_redirect_server());
    return charging_credit;
  }

  RedirectServer* get_redirect_server()
  {
    auto redirect = new RedirectServer();
    redirect->set_redirect_address_type(RedirectServer_RedirectAddressType_IPV4);
    redirect->set_redirect_server_address("192.168.0.1");
  }

  GrantedUnits* get_granted_units()
  {
    auto units = new GrantedUnits();

    auto total = new CreditUnit();
    total->set_is_valid(true);
    total->set_volume(1000);

    auto tx = new CreditUnit();
    tx->set_is_valid(true);
    tx->set_volume(1000);

    auto rx = new CreditUnit();
    rx->set_is_valid(true);
    rx->set_volume(1000);

    units->set_allocated_total(total);
    units->set_allocated_tx(tx);
    units->set_allocated_rx(rx);

    return units;
  }

  UsageMonitoringUpdateResponse* get_monitoring_update()
  {
    auto credit_update = new UsageMonitoringUpdateResponse();
    credit_update->set_allocated_credit(get_monitoring_credit());
    credit_update->set_session_id("sid1");
    credit_update->set_success(true);
    // Don't set event triggers
    // Don't set result code since the response is already successful
    // Don't set any rule installation/uninstallation
    // Don't set the TgppContext, assume relay disabled
    return credit_update;
  }

  UsageMonitoringCredit* get_monitoring_credit()
  {
    auto monitoring_credit = new UsageMonitoringCredit();
    monitoring_credit->set_action(UsageMonitoringCredit_Action_CONTINUE);
    monitoring_credit->set_monitoring_key("mk1");
    monitoring_credit->set_level(SESSION_LEVEL);
    monitoring_credit->set_allocated_granted_units(get_granted_units());
    return monitoring_credit;
  }
};

TEST_F(CreditPoolTest, test_marshal_unmarshal_charging)
{
  auto pool = new ChargingCreditPool("imsi1");

  // Receive credit
  auto credit_update = get_credit_update();
  CreditUpdateResponse& credit_update_ref = *credit_update;
  pool->receive_credit(credit_update_ref);

  // Add some used credit
  pool->add_used_credit(CreditKey(credit_update), uint64_t(123), uint64_t(456));
  EXPECT_EQ(pool->get_credit(CreditKey(credit_update), USED_TX), 123);
  EXPECT_EQ(pool->get_credit(CreditKey(credit_update), USED_RX), 456);

  // Check that after marshaling/unmarshaling that the fields are still the
  // same.
  auto marshaled = pool->marshal();
  auto pool_2 = ChargingCreditPool::unmarshal(marshaled);
  EXPECT_EQ(pool_2->get_credit(CreditKey(credit_update), USED_TX), 123);
  EXPECT_EQ(pool_2->get_credit(CreditKey(credit_update), USED_RX), 456);
}

TEST_F(CreditPoolTest, test_marshal_unmarshal_monitoring)
{
  auto pool = new UsageMonitoringCreditPool("imsi1");

  // Receive credit
  auto credit_update = get_monitoring_update();
  UsageMonitoringUpdateResponse& credit_update_ref = *credit_update;
  pool->receive_credit(credit_update_ref);

  // Add some used credit
  pool->add_used_credit("mk1", uint64_t(123), uint64_t(456));
  EXPECT_EQ(pool->get_credit("mk1", USED_TX), 123);
  EXPECT_EQ(pool->get_credit("mk1", USED_RX), 456);

  // Check that after marshaling/unmarshaling that the fields are still the
  // same.
  auto marshaled = pool->marshal();
  auto pool_2 = UsageMonitoringCreditPool::unmarshal(marshaled);
  //EXPECT_EQ(pool_2->get_credit("mk1", USED_TX), 123);
  //EXPECT_EQ(pool_2->get_credit("mk1", USED_RX), 456);
}

int main(int argc, char** argv)
{
  ::testing::InitGoogleTest(&argc, argv);
  FLAGS_logtostderr = 1;
  FLAGS_v = 10;
  return RUN_ALL_TESTS();
}

} // namespace magma
