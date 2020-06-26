/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <limits>

#include "ChargingGrant.h"
#include "magma_logging.h"

namespace magma {
  StoredChargingGrant ChargingGrant::marshal() {
    StoredChargingGrant marshaled{};
    marshaled.is_final = is_final_grant;
    marshaled.final_action_info.final_action =
      final_action_info.final_action;
    marshaled.final_action_info.redirect_server =
      final_action_info.redirect_server;
    marshaled.reauth_state = reauth_state;
    marshaled.service_state = service_state;
    marshaled.expiry_time = expiry_time;
    marshaled.credit = credit.marshal();
    return marshaled;
  }

  ChargingGrant ChargingGrant::unmarshal(const StoredChargingGrant &marshaled) {
    ChargingGrant charging;
    charging.credit =
      SessionCredit::unmarshal(marshaled.credit, CreditType::CHARGING);

    FinalActionInfo final_action_info;
    final_action_info.final_action = marshaled.final_action_info.final_action;
    final_action_info.redirect_server =
        marshaled.final_action_info.redirect_server;
    charging.final_action_info = final_action_info;

    charging.reauth_state = marshaled.reauth_state;
    charging.service_state = marshaled.service_state;
    charging.expiry_time = marshaled.expiry_time;
    charging.is_final_grant = marshaled.is_final;

    return charging;
  }

  void ChargingGrant::set_final_action_info(const magma::lte::ChargingCredit &credit) {
    if (credit.is_final()) {
      final_action_info.final_action = credit.final_action();
      if (credit.final_action() == ChargingCredit_FinalAction_REDIRECT) {
        final_action_info.redirect_server = credit.redirect_server();
      }
    }
  }
} // namespace magma
