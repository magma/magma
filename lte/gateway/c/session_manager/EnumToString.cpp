/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "EnumToString.h"

namespace magma {
  std::string reauth_state_to_str(ReAuthState state) {
    switch (state) {
    case REAUTH_NOT_NEEDED:
      return "REAUTH_NOT_NEEDED";
    case REAUTH_REQUIRED:
      return "REAUTH_REQUIRED";
    case REAUTH_PROCESSING:
      return "REAUTH_PROCESSING";
    default:
      return "INVALID REAUTH STATE";
    }
  }

  std::string service_state_to_str(ServiceState state) {
    switch (state) {
    case SERVICE_ENABLED:
      return "SERVICE_ENABLED";
    case SERVICE_NEEDS_DEACTIVATION:
      return "SERVICE_NEEDS_DEACTIVATION";
    case SERVICE_DISABLED:
      return "SERVICE_DISABLED";
    case SERVICE_NEEDS_ACTIVATION:
      return "SERVICE_NEEDS_ACTIVATION";
    case SERVICE_REDIRECTED:
      return "SERVICE_REDIRECTED";
    case SERVICE_RESTRICTED:
      return "SERVICE_RESTRICTED";
    default:
      return "INVALID SERVICE STATE";
    }
  }

  std::string final_action_to_str(ChargingCredit_FinalAction final_action) {
    switch (final_action) {
    case ChargingCredit_FinalAction_TERMINATE:
      return "TERMINATE";
    case ChargingCredit_FinalAction_REDIRECT:
      return "REDIRECT";
    case ChargingCredit_FinalAction_RESTRICT_ACCESS:
      return "RESTRICT_ACCESS";
    default:
      return "INVALID FINAL ACTION";
    }
  }

  std::string session_fsm_state_to_str(SessionFsmState state) {
    switch (state) {
    case SESSION_ACTIVE:
      return "SESSION_ACTIVE";
    case SESSION_TERMINATING_FLOW_ACTIVE:
      return "SESSION_TERMINATING_FLOW_ACTIVE";
    case SESSION_TERMINATING_AGGREGATING_STATS:
      return "SESSION_TERMINATING_AGGREGATING_STATS";
    case SESSION_TERMINATING_FLOW_DELETED:
      return "SESSION_TERMINATING_FLOW_DELETED";
    case SESSION_TERMINATED:
      return "SESSION_TERMINATED";
    case SESSION_TERMINATION_SCHEDULED:
      return "SESSION_TERMINATION_SCHEDULED";
    default:
      return "INVALID SESSION FSM STATE";
    }
  }
} // namespace magma
