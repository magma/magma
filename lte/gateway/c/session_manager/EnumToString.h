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

namespace magma {
std::string reauth_state_to_str(ReAuthState state);

std::string service_state_to_str(ServiceState state);

std::string final_action_to_str(ChargingCredit_FinalAction final_action);

std::string grant_type_to_str(GrantTrackingType grant_type);

std::string session_fsm_state_to_str(SessionFsmState state);
}  // namespace magma
