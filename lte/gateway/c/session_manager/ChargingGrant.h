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
  // The number of seconds for which the credit is valid
  std::time_t expiry_time;
  ServiceState service_state;
  ReAuthState reauth_state;

  ChargingGrant() : credit(CreditType::CHARGING) {}

  StoredChargingGrant marshal();

  static ChargingGrant unmarshal(const StoredChargingGrant &marshaled);

  void set_final_action_info(const magma::lte::ChargingCredit &credit);
};

}  // namespace magma
