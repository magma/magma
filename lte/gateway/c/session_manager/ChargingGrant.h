/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include "StoredState.h"
#include "SessionCredit.h"

namespace magma {

// ChargingGrant is a struct because all fields are public
struct ChargingGrant {
  // Keep track of used/reported/allowed bytes
  SessionCredit credit;
  // When this is true, final_action should be acted on upon credit exhaust
  bool is_final_grant;
  // Only valid if is_final_grant is true
  FinalActionInfo final_action_info;
  // The expiry time for the credit's validity
  std::time_t expiry_time;
  ServiceState service_state;
  ReAuthState reauth_state;

  ChargingGrant() : credit(CreditType::CHARGING) {}

  // ChargingGrant -> StoredChargingGrant
  StoredChargingGrant marshal();

  // StoredChargingGrant -> ChargingGrant
  static ChargingGrant unmarshal(const StoredChargingGrant &marshaled);

  // Set is_final_grant and final_action_info values
  void set_final_action_info(const magma::lte::ChargingCredit &credit);

  SessionCreditUpdateCriteria get_update_criteria();

  // Convert rel_time_sec, which is a delta value in seconds, into a timestamp
  // and assign it to expiry_time
  void set_expiry_time_as_timestamp(uint32_t delta_time_sec);

  // Determine whether the charging grant should send an update request
  // Return true if an update is required, with the update_type set to indicate
  // the reason.
  // Return false otherwise. In this case, update_type is untouched.
  bool get_update_type(CreditUsage::UpdateType* update_type) const;

  void set_reauth_state(const ReAuthState new_state, SessionCreditUpdateCriteria &uc);
};

}  // namespace magma
