/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
#include <gtest/gtest.h>

#include <chrono>
#include <thread>

#include "ChargingGrant.h"
#include "ProtobufCreators.h"

using ::testing::Test;

namespace magma {

class ChargingGrantTest : public ::testing::Test {
 protected:
  ChargingGrant get_default_grant() {
    ChargingGrant grant;
    grant.is_final_grant = false;
    grant.expiry_time    = std::numeric_limits<std::time_t>::max();
    grant.service_state  = SERVICE_ENABLED;
    grant.reauth_state   = REAUTH_NOT_NEEDED;
    grant.suspended      = false;
    return grant;
  }

  ChargingGrant get_default_grant(ChargingCredit_FinalAction action) {
    auto grant              = get_default_grant();
    grant.is_final_grant    = true;
    grant.final_action_info = get_final_action_info(action);
    return grant;
  }

  FinalActionInfo get_final_action_info(ChargingCredit_FinalAction action) {
    FinalActionInfo fa;
    fa.final_action = action;
    if (action == ChargingCredit_FinalAction_REDIRECT) {
      fa.redirect_server.set_redirect_address_type(
          RedirectServer_RedirectAddressType::
              RedirectServer_RedirectAddressType_IPV6);
      fa.redirect_server.set_redirect_server_address("addr");
    } else if (action == ChargingCredit_FinalAction_RESTRICT_ACCESS) {
      fa.restrict_rules.push_back("restrict_rule");
    }
    return fa;
  }
};

TEST_F(ChargingGrantTest, test_marshal) {
  // Get a grant with REDIRECT final action since that fills out all the fields
  ChargingGrant grant = get_default_grant(ChargingCredit_FinalAction_REDIRECT);
  // Set ReAuth state and service state to non-0 values
  grant.service_state = SERVICE_DISABLED;
  grant.reauth_state  = REAUTH_REQUIRED;

  StoredChargingGrant stored = grant.marshal();
  EXPECT_EQ(grant.is_final_grant, stored.is_final);
  EXPECT_EQ(
      grant.final_action_info.final_action,
      stored.final_action_info.final_action);

  auto redirect        = grant.final_action_info.redirect_server;
  auto stored_redirect = stored.final_action_info.redirect_server;
  EXPECT_EQ(
      redirect.redirect_server_address(),
      stored_redirect.redirect_server_address());
  EXPECT_EQ(
      redirect.redirect_address_type(),
      stored_redirect.redirect_address_type());

  EXPECT_EQ(grant.expiry_time, stored.expiry_time);
  EXPECT_EQ(grant.service_state, stored.service_state);
  EXPECT_EQ(grant.reauth_state, stored.reauth_state);
}

TEST_F(ChargingGrantTest, test_get_update_type) {
  ChargingGrant grant = get_default_grant();
  GrantedUnits gsu;
  uint64_t total_grant = 1000;
  create_granted_units(&total_grant, NULL, NULL, &gsu);

  auto uc = grant.get_update_criteria();

  grant.credit.receive_credit(gsu, &uc);
  grant.is_final_grant = false;
  EXPECT_EQ(uc.grant_tracking_type, TOTAL_ONLY);

  grant.credit.add_used_credit(2000, 0, &uc);
  EXPECT_TRUE(grant.credit.is_quota_exhausted(0.8));
  // Credit is exhausted, we expect a quota exhaustion update type
  CreditUsage::UpdateType update_type;
  EXPECT_TRUE(grant.get_update_type(&update_type));
  EXPECT_EQ(update_type, CreditUsage::QUOTA_EXHAUSTED);

  // Check how much we will report (we will report everything, even if
  // we go have gone over the allowed quota)
  auto update = grant.credit.get_usage_for_reporting(&uc);
  EXPECT_EQ(update.bytes_tx, 2000);
  EXPECT_TRUE(grant.credit.is_reporting());

  // Credit is being reported, no update necessary
  EXPECT_FALSE(grant.get_update_type(&update_type));

  // Receive a final grant
  total_grant = 0;
  create_granted_units(&total_grant, NULL, NULL, &gsu);
  grant.credit.receive_credit(gsu, &uc);
  EXPECT_EQ(uc.grant_tracking_type, TOTAL_ONLY);
  grant.is_final_grant = true;
  grant.credit.reset_reporting_credit(&uc);

  EXPECT_FALSE(grant.get_update_type(&update_type));

  // Set reauth state to be REAUTH_REQUIRED && final_grant
  grant.is_final_grant = true;
  grant.reauth_state   = REAUTH_REQUIRED;
  EXPECT_TRUE(grant.get_update_type(&update_type));
  EXPECT_EQ(update_type, CreditUsage::REAUTH_REQUIRED);

  // Set final_grant && validity timer
  grant.is_final_grant = true;
  grant.reauth_state   = REAUTH_NOT_NEEDED;
  grant.expiry_time    = time(nullptr) - 50;  // 50 seconds ago
  EXPECT_TRUE(grant.get_update_type(&update_type));
  EXPECT_EQ(update_type, CreditUsage::VALIDITY_TIMER_EXPIRED);
}

TEST_F(ChargingGrantTest, test_should_deactivate_service) {
  // Create a grant with some quota
  ChargingGrant grant;
  GrantedUnits gsu;
  uint64_t total_grant = 1000;
  create_granted_units(&total_grant, NULL, NULL, &gsu);
  auto uc = grant.get_update_criteria();
  grant.credit.receive_credit(gsu, &uc);
  EXPECT_FALSE(grant.credit.is_quota_exhausted(0.8));

  // If quota is not exhausted, don't deactivate
  grant.is_final_grant                 = true;
  grant.final_action_info.final_action = ChargingCredit_FinalAction_TERMINATE;
  grant.service_state                  = SERVICE_ENABLED;
  EXPECT_FALSE(grant.should_deactivate_service());

  // Exhaust quota
  uc = grant.get_update_criteria();
  grant.credit.add_used_credit(2000, 0, &uc);
  EXPECT_TRUE(grant.credit.is_quota_exhausted(0.8));

  // Test TERMINATE_SERVICE_WHEN_QUOTA_EXHAUSTED flag && is_final
  SessionCredit::TERMINATE_SERVICE_WHEN_QUOTA_EXHAUSTED = false;
  grant.is_final_grant                                  = true;
  grant.final_action_info.final_action = ChargingCredit_FinalAction_TERMINATE;
  grant.service_state                  = SERVICE_ENABLED;
  EXPECT_FALSE(grant.should_deactivate_service());

  // Test TERMINATE_SERVICE_WHEN_QUOTA_EXHAUSTED is set && is_final
  SessionCredit::TERMINATE_SERVICE_WHEN_QUOTA_EXHAUSTED = true;
  grant.is_final_grant                                  = true;
  grant.final_action_info.final_action = ChargingCredit_FinalAction_TERMINATE;
  grant.service_state                  = SERVICE_ENABLED;
  EXPECT_TRUE(grant.should_deactivate_service());

  // If service state is not ENABLED we should not deactivate service for FUA
  // redirect / restrict
  SessionCredit::TERMINATE_SERVICE_WHEN_QUOTA_EXHAUSTED = true;
  grant.is_final_grant                                  = true;
  grant.final_action_info.final_action = ChargingCredit_FinalAction_REDIRECT;
  grant.service_state                  = SERVICE_DISABLED;
  EXPECT_FALSE(grant.should_deactivate_service());

  // If service state is not ENABLED we should deactivate service for FUA
  // terminate
  SessionCredit::TERMINATE_SERVICE_WHEN_QUOTA_EXHAUSTED = true;
  grant.is_final_grant                                  = true;
  grant.final_action_info.final_action = ChargingCredit_FinalAction_TERMINATE;
  grant.service_state                  = SERVICE_DISABLED;
  EXPECT_TRUE(grant.should_deactivate_service());
}

TEST_F(ChargingGrantTest, test_get_action) {
  ChargingGrant grant = get_default_grant();
  auto uc             = grant.get_update_criteria();
  GrantedUnits gsu;
  uint64_t total_grant = 1024;
  create_granted_units(&total_grant, NULL, NULL, &gsu);

  // Not a final grant
  grant.is_final_grant = false;
  grant.credit.receive_credit(gsu, &uc);
  grant.credit.add_used_credit(1024, 0, &uc);
  auto cont_action = grant.get_action(&uc);
  EXPECT_EQ(cont_action, CONTINUE_SERVICE);
  // Check both the grant's service_state and update criteria
  EXPECT_EQ(grant.service_state, CONTINUE_SERVICE);
  EXPECT_EQ(uc.service_state, CONTINUE_SERVICE);

  // Final with TERMINATE final action
  grant.is_final_grant = true;
  grant.final_action_info =
      get_final_action_info(ChargingCredit_FinalAction_TERMINATE);
  grant.credit.receive_credit(gsu, &uc);
  grant.credit.add_used_credit(2048, 0, &uc);
  grant.credit.add_used_credit(30, 20, &uc);
  grant.service_state = SERVICE_NEEDS_DEACTIVATION;
  auto term_action    = grant.get_action(&uc);
  // Check that the update criteria also includes the changes
  EXPECT_EQ(term_action, TERMINATE_SERVICE);

  // Termination action only returned once
  auto repeated_action = grant.get_action(&uc);
  EXPECT_EQ(repeated_action, CONTINUE_SERVICE);
}

TEST_F(ChargingGrantTest, test_get_action_redirect) {
  ChargingGrant grant = get_default_grant();
  auto uc             = grant.get_update_criteria();
  GrantedUnits gsu;
  uint64_t total_grant = 1024;
  create_granted_units(&total_grant, NULL, NULL, &gsu);

  // Final with REDIRECT final action
  grant.is_final_grant = true;
  grant.final_action_info =
      get_final_action_info(ChargingCredit_FinalAction_REDIRECT);
  grant.credit.receive_credit(gsu, &uc);
  grant.credit.add_used_credit(2048, 0, &uc);
  grant.credit.add_used_credit(30, 20, &uc);
  grant.service_state = SERVICE_NEEDS_DEACTIVATION;
  auto term_action    = grant.get_action(&uc);
  // Check that the update criteria also includes the changes
  EXPECT_EQ(term_action, REDIRECT);

  // Termination action only returned once
  auto repeated_action = grant.get_action(&uc);
  EXPECT_EQ(repeated_action, CONTINUE_SERVICE);
}

TEST_F(ChargingGrantTest, test_get_action_restrict) {
  ChargingGrant grant = get_default_grant();
  auto uc             = grant.get_update_criteria();
  GrantedUnits gsu;
  uint64_t total_grant = 1024;
  create_granted_units(&total_grant, NULL, NULL, &gsu);

  // Final with REDIRECT final action
  grant.is_final_grant = true;
  grant.final_action_info =
      get_final_action_info(ChargingCredit_FinalAction_RESTRICT_ACCESS);
  grant.credit.receive_credit(gsu, &uc);
  grant.credit.add_used_credit(2048, 0, &uc);
  grant.credit.add_used_credit(30, 20, &uc);
  grant.service_state = SERVICE_NEEDS_DEACTIVATION;
  auto term_action    = grant.get_action(&uc);
  // Check that the update criteria also includes the changes
  EXPECT_EQ(term_action, RESTRICT_ACCESS);

  // Termination action only returned once
  auto repeated_action = grant.get_action(&uc);
  EXPECT_EQ(repeated_action, CONTINUE_SERVICE);
}

// test_tolerance_quota_exhausted checks that user will not be terminated if
// quota is exhausted but not final unit indication is received.
// That can happen if the quota reported by pipeline is too big and we go over
// both the threshold (0.8) and the maximum allowed quota
TEST_F(ChargingGrantTest, test_tolerance_quota_exhausted) {
  auto grant   = get_default_grant();
  auto& credit = grant.credit;
  auto uc      = grant.get_update_criteria();
  GrantedUnits gsu;
  uint64_t total_grant = 1000;
  create_granted_units(&total_grant, NULL, NULL, &gsu);
  grant.credit.receive_credit(gsu, &uc);
  EXPECT_EQ(uc.grant_tracking_type, TOTAL_ONLY);

  // Not a final credit
  grant.is_final_grant = false;
  credit.add_used_credit(2000, 0, &uc);
  EXPECT_EQ(uc.bucket_deltas[USED_TX], 2000);
  EXPECT_TRUE(credit.is_quota_exhausted(0.8));
  // continue the service even we are over the quota (but not final unit)
  EXPECT_EQ(grant.get_action(&uc), CONTINUE_SERVICE);

  // Check how much we will report (we will report everything, even if
  // we go have gone over the allowed quota)
  CreditUsage::UpdateType update_type;
  EXPECT_TRUE(grant.get_update_type(&update_type));
  auto c_usage = grant.get_credit_usage(update_type, &uc, false);
  EXPECT_EQ(update_type, CreditUsage::QUOTA_EXHAUSTED);
  EXPECT_EQ(c_usage.bytes_tx(), 2000);
  EXPECT_EQ(c_usage.bytes_rx(), 0);
  EXPECT_EQ(credit.get_credit(USED_TX), 2000);
  EXPECT_EQ(credit.get_credit(REPORTING_TX), 2000);

  // Now receive new quota (not final unit)
  uc = grant.get_update_criteria();  // reset UC
  grant.credit.receive_credit(gsu, &uc);
  // we overused, so we take into consideration the 2000 we used plus granted
  // 1000
  EXPECT_EQ(credit.get_credit(ALLOWED_TOTAL), 3000);
  EXPECT_EQ(credit.get_credit(REPORTED_TX), 2000);
  EXPECT_EQ(credit.get_credit(USED_TX), 2000);
  // we overused, so the delta is the overusage
  EXPECT_EQ(uc.bucket_deltas[ALLOWED_TOTAL], 2000);

  // No update should be triggered as everything is reported
  uc = grant.get_update_criteria();  // reset UC
  EXPECT_FALSE(grant.get_update_type(&update_type));
  EXPECT_EQ(credit.get_credit(USED_TX), 2000);
  EXPECT_EQ(credit.get_credit(REPORTING_TX), 0);

  // receive some more FINAL grant that will go over part of the used and not
  // reported credit
  uc                   = grant.get_update_criteria();  // reset UC
  grant.is_final_grant = true;
  grant.final_action_info =
      get_final_action_info(ChargingCredit_FinalAction_TERMINATE);
  grant.credit.receive_credit(gsu, &uc);
  EXPECT_EQ(credit.get_credit(ALLOWED_TOTAL), 4000);
  EXPECT_EQ(credit.get_credit(REPORTED_TX), 2000);
  EXPECT_EQ(credit.get_credit(USED_TX), 2000);
  EXPECT_EQ(uc.bucket_deltas[ALLOWED_TOTAL], 1000);

  // Use enough credit to exceed the given quota
  uc = grant.get_update_criteria();  // reset UC
  credit.add_used_credit(2000, 0, &uc);
  EXPECT_EQ(uc.bucket_deltas[USED_TX], 2000);
  EXPECT_EQ(credit.get_credit(ALLOWED_TOTAL), 4000);
  EXPECT_EQ(credit.get_credit(REPORTED_TX), 2000);
  EXPECT_EQ(credit.get_credit(USED_TX), 4000);
  EXPECT_TRUE(credit.is_quota_exhausted(1));  // 100% exceeded
  EXPECT_TRUE(grant.should_deactivate_service());
  grant.set_service_state(SERVICE_NEEDS_DEACTIVATION, &uc);

  // Since this is the final grant, we should not report anything
  EXPECT_FALSE(grant.get_update_type(&update_type));
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace magma
