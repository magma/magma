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

#include "ChargingGrant.h"
#include "ProtobufCreators.h"
#include <gtest/gtest.h>

using ::testing::Test;

#define HIGH_CREDIT 1000000

namespace magma {

class ChargingGrantTest : public ::testing::Test {
protected:
  ChargingGrant get_default_grant() {
    FinalActionInfo final_action = {
    .final_action = ChargingCredit_FinalAction_TERMINATE,
    };
    final_action.redirect_server.set_redirect_address_type(
      RedirectServer_RedirectAddressType::RedirectServer_RedirectAddressType_IPV6);
    final_action.redirect_server.set_redirect_server_address("addr");

    ChargingGrant grant;
    grant.is_final_grant = true;
    grant.final_action_info = final_action;
    grant.expiry_time = time(NULL);
    grant.service_state = SERVICE_NEEDS_ACTIVATION;
    grant.reauth_state = REAUTH_PROCESSING;
    return grant;
  }
};

TEST_F(ChargingGrantTest, test_marshal) {
  ChargingGrant grant = get_default_grant();

  StoredChargingGrant stored = grant.marshal();
  EXPECT_EQ(grant.is_final_grant, stored.is_final);
  EXPECT_EQ(grant.final_action_info.final_action,
            stored.final_action_info.final_action);

  auto redirect = grant.final_action_info.redirect_server;
  auto stored_redirect = stored.final_action_info.redirect_server;
  EXPECT_EQ(redirect.redirect_server_address(),
            stored_redirect.redirect_server_address());
  EXPECT_EQ(redirect.redirect_address_type(),
            stored_redirect.redirect_address_type());

  EXPECT_EQ(grant.expiry_time, stored.expiry_time);
  EXPECT_EQ(grant.service_state, stored.service_state);
  EXPECT_EQ(grant.reauth_state, stored.reauth_state);
}

TEST_F(ChargingGrantTest, test_get_update_type) {
  ChargingGrant grant = get_default_grant();
  GrantedUnits gsu;
  uint64_t total_grant = 1000;
  FinalActionInfo final_action_info;

  auto uc = grant.get_update_criteria();
  create_granted_units(&total_grant, NULL, NULL, &gsu);
  grant.credit.receive_credit(gsu, 0, false, final_action_info, uc);
  EXPECT_EQ(uc.grant_tracking_type, TOTAL_ONLY);

  grant.credit.add_used_credit(2000, 0, uc);
  EXPECT_TRUE(grant.credit.is_quota_exhausted(0.8));
  // Credit is exhausted, we expect a quota exhaustion updat type
  CreditUsage::UpdateType update_type;
  EXPECT_TRUE(grant.get_update_type(&update_type));
  EXPECT_EQ(update_type, CreditUsage::QUOTA_EXHAUSTED);

  // Check how much we will report (as much as the total allowed, the rest
  // will be left unreported and reported in future updates)
  auto update = grant.credit.get_usage_for_reporting(uc);
  EXPECT_EQ(update.bytes_tx, 1000);
  EXPECT_TRUE(grant.credit.is_reporting());

  // Credit is being reported, no update necessary
  EXPECT_FALSE(grant.get_update_type(&update_type));

  // Receive a final grant
  grant.credit.receive_credit(gsu, 0, true, final_action_info, uc);
  EXPECT_EQ(uc.grant_tracking_type, TOTAL_ONLY);
  grant.credit.reset_reporting_credit(uc);

  EXPECT_FALSE(grant.get_update_type(&update_type));
}

int main(int argc, char **argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

} // namespace magma
