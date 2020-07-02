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
#include "EnumToString.h"
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
    is_final_grant = credit.is_final();
    if (is_final_grant) {
      final_action_info.final_action = credit.final_action();
      if (credit.final_action() == ChargingCredit_FinalAction_REDIRECT) {
        final_action_info.redirect_server = credit.redirect_server();
      }
    }
  }

  SessionCreditUpdateCriteria ChargingGrant::get_update_criteria() {
    SessionCreditUpdateCriteria uc = credit.get_update_criteria();
    uc.is_final = is_final_grant;
    uc.final_action_info = final_action_info;
    uc.expiry_time = expiry_time;

    // todo overwrite these in later diffs
    uc.reauth_state = reauth_state;
    //uc.service_state = service_state_;
    return uc;
  }

  void ChargingGrant::set_expiry_time_as_timestamp(uint32_t delta_time_sec) {
    if (delta_time_sec == 0) {
      expiry_time = std::numeric_limits<std::time_t>::max();
    } else {
      expiry_time = std::time(nullptr) + delta_time_sec;
    }
  }

  bool ChargingGrant::get_update_type(CreditUsage::UpdateType* update_type) const {
    if (credit.is_reporting()) {
      return false; // No update
    }
    if (reauth_state == REAUTH_REQUIRED) {
      *update_type = CreditUsage::REAUTH_REQUIRED;
      return true;
    }
    if (credit.is_final_grant() && credit.is_quota_exhausted(1)) {
      // Don't request updates if this is the final grant
      return false;
    }
    if (credit.is_quota_exhausted(SessionCredit::USAGE_REPORTING_THRESHOLD)) {
      *update_type = CreditUsage::QUOTA_EXHAUSTED;
      return true;
    }
    if (credit.validity_timer_expired()) {
      *update_type = CreditUsage::VALIDITY_TIMER_EXPIRED;
      return true;
    }
    return false;
  }

  void ChargingGrant::set_reauth_state(
    const ReAuthState new_state, SessionCreditUpdateCriteria& uc) {
    if (reauth_state != new_state) {
      MLOG(MDEBUG) << "ReAuth state change from "
                   << reauth_state_to_str(reauth_state) << " to "
                   << reauth_state_to_str(new_state);
    }
    reauth_state = new_state;
    uc.reauth_state = new_state;
  }

} // namespace magma
