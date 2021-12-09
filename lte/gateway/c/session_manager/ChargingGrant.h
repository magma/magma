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
#pragma once

#include <stdint.h>
#include <ctime>

#include "DiameterCodes.h"
#include "ServiceAction.h"
#include "SessionCredit.h"
#include "StoredState.h"
#include "Types.h"
#include "lte/protos/policydb.pb.h"
#include "lte/protos/session_manager.pb.h"

namespace magma {

enum CreditValidity {
  VALID_CREDIT    = 0,
  INVALID_CREDIT  = 1,
  TRANSIENT_ERROR = 2,
};

// Used to keep track of credit grants from Gy. Some features of this type of
// grants include:
//    1. Final Unit Handling: Instructions on what is to happen to the session
//       on final grant exhaustion. Relevant fields are is_final_grant and
//       final_action_info. We currently support TERMINATE and REDIRECT.
//    2. Expiry/Validity Time: This grant can be received with an int to
//       indicate the number of seconds the grant is valid for. The expiry_time
//       field indicates the time at which the grant is no longer valid.
// ChargingGrant is a struct because all fields are public
struct ChargingGrant {
  // Keep track of used/reported/allowed bytes
  SessionCredit credit;
  // When this is true, final_action should be acted on upon credit exhaust
  bool is_final_grant;
  // Only valid if is_final_grant is true
  FinalActionInfo final_action_info;
  // The expiry time for the credit's validity
  // https://tools.ietf.org/html/rfc4006#section-8.33
  std::time_t expiry_time;
  ServiceState service_state;
  ReAuthState reauth_state;
  // indicates the rules have been removed from pipelined
  bool suspended;

  const uint32_t REDIRECT_FLOW_PRIORITY = 2000;

  // Default states
  ChargingGrant()
      : credit(),
        is_final_grant(false),
        service_state(SERVICE_ENABLED),
        reauth_state(REAUTH_NOT_NEEDED),
        suspended(false) {}

  explicit ChargingGrant(const StoredChargingGrant& marshaled);

  // ChargingGrant -> StoredChargingGrant
  StoredChargingGrant marshal();

  void receive_charging_grant(
      const CreditUpdateResponse& update,
      SessionCreditUpdateCriteria* credit_uc = NULL);

  // Returns true if the credit returned from the Policy component is valid and
  // good to be installed.
  static CreditValidity get_credit_response_validity_type(
      const CreditUpdateResponse& update);

  // Returns a SessionCreditUpdateCriteria that reflects the current state
  SessionCreditUpdateCriteria get_update_criteria();

  // Determine whether the charging grant should send an update request
  // Return true if an update is required, with the update_type set to indicate
  // the reason.
  // Return false otherwise. In this case, update_type is untouched.
  bool get_update_type(CreditUsage::UpdateType* update_type) const;

  // get_action returns the action to take on the credit based on the last
  // update. If no action needs to take place, CONTINUE_SERVICE is returned.
  ServiceActionType get_action(SessionCreditUpdateCriteria* update_criteria);

  // Return if the service needs activation
  bool should_be_unsuspended() const;

  bool get_suspended() const { return suspended; }

  // Get unreported usage from credit and return as part of CreditUsage
  // The update_type is also included in CreditUsage
  // If the grant is final or is_terminate is true, we include all unreported
  // usage, otherwise we only include unreported usage up to the allocated
  // amount.
  CreditUsage get_credit_usage(
      CreditUsage::UpdateType update_type,
      SessionCreditUpdateCriteria* credit_uc, bool is_terminate);

  // Return true if the service needs to be deactivated.
  // In order to deactivate, a few things are considered in order.
  // 1. Credit must be exhausted
  // 2. TERMINATE_SERVICE_WHEN_QUOTA_EXHAUSTED is not set
  // 3. FUA must be set
  //    3a. For FUA-terminate, always deactivate
  //    3b. For FUA-redirect/restrict, deactivate if not already deactivated
  bool should_deactivate_service() const;

  // Convert FinalAction enum to ServiceActionType
  ServiceActionType final_action_to_action(
      const ChargingCredit_FinalAction action) const;

  ServiceActionType final_action_to_action_on_suspension(
      const ChargingCredit_FinalAction action) const;

  void set_suspended(bool suspended, SessionCreditUpdateCriteria* credit_uc);

  // Set the object and update criteria's reauth state to new_state.
  void set_reauth_state(
      const ReAuthState new_state, SessionCreditUpdateCriteria* credit_uc);

  // Set the object and update criteria's service state to new_state.
  void set_service_state(
      const ServiceState new_service_state,
      SessionCreditUpdateCriteria* credit_uc);

  // Set the flag reporting. Used to signal this credit is waiting to receive
  // a response from the core
  void set_reporting(bool reporting);

  // Rollback reporting changes for failed updates
  void reset_reporting_grant(SessionCreditUpdateCriteria* credit_uc);

  // Log information about the grant received
  void log_received_grant(const CreditUpdateResponse& update);

  PolicyRule make_redirect_rule();
};

}  // namespace magma
