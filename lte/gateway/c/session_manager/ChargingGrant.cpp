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

#include <ctime>
#include <limits>

#include "ChargingGrant.h"
#include "EnumToString.h"
#include "magma_logging.h"

namespace magma {
ChargingGrant::ChargingGrant(const StoredChargingGrant& marshaled) {
  credit = SessionCredit(marshaled.credit);

  final_action_info.final_action = marshaled.final_action_info.final_action;
  final_action_info.redirect_server =
      marshaled.final_action_info.redirect_server;
  final_action_info.restrict_rules = marshaled.final_action_info.restrict_rules;

  reauth_state   = marshaled.reauth_state;
  service_state  = marshaled.service_state;
  expiry_time    = marshaled.expiry_time;
  is_final_grant = marshaled.is_final;
}

StoredChargingGrant ChargingGrant::marshal() {
  StoredChargingGrant marshaled{};
  marshaled.is_final                       = is_final_grant;
  marshaled.final_action_info.final_action = final_action_info.final_action;
  marshaled.final_action_info.redirect_server =
      final_action_info.redirect_server;
  marshaled.final_action_info.restrict_rules = final_action_info.restrict_rules;
  marshaled.reauth_state                     = reauth_state;
  marshaled.service_state                    = service_state;
  marshaled.expiry_time                      = expiry_time;
  marshaled.credit                           = credit.marshal();
  return marshaled;
}

void ChargingGrant::receive_charging_grant(
    const magma::lte::ChargingCredit& p_credit,
    SessionCreditUpdateCriteria* uc) {
  credit.receive_credit(p_credit.granted_units(), uc);

  // Final Action
  is_final_grant = p_credit.is_final();
  if (is_final_grant) {
    final_action_info.final_action = p_credit.final_action();
    switch (final_action_info.final_action) {
      case ChargingCredit_FinalAction_REDIRECT:
        final_action_info.redirect_server = p_credit.redirect_server();
        break;
      case ChargingCredit_FinalAction_RESTRICT_ACCESS:
        // Clear the previous restrict rules
        final_action_info.restrict_rules.clear();
        for (auto rule : p_credit.restrict_rules()) {
          final_action_info.restrict_rules.push_back(rule);
        }
        break;
      default:  // do nothing
        break;
    }
    log_final_action_info();
  }

  // Expiry Time
  const auto delta_time_sec = p_credit.validity_time();
  if (delta_time_sec == 0) {
    expiry_time = std::numeric_limits<std::time_t>::max();
  } else {
    expiry_time = std::time(nullptr) + delta_time_sec;
  }

  // Update the UpdateCriteria if not NULL
  if (uc != NULL) {
    uc->is_final          = is_final_grant;
    uc->final_action_info = final_action_info;
    uc->expiry_time       = expiry_time;
  }
}

SessionCreditUpdateCriteria ChargingGrant::get_update_criteria() {
  SessionCreditUpdateCriteria uc = credit.get_update_criteria();
  uc.is_final                    = is_final_grant;
  uc.final_action_info           = final_action_info;
  uc.expiry_time                 = expiry_time;
  uc.reauth_state                = reauth_state;
  uc.service_state               = service_state;
  return uc;
}

CreditUsage ChargingGrant::get_credit_usage(
    CreditUsage::UpdateType update_type, SessionCreditUpdateCriteria& uc,
    bool is_terminate) {
  CreditUsage p_usage;
  SessionCredit::Usage credit_usage;

  if (is_final_grant || is_terminate) {
    credit_usage = credit.get_all_unreported_usage_for_reporting(uc);
  } else {
    credit_usage = credit.get_usage_for_reporting(uc);
  }
  p_usage.set_bytes_tx(credit_usage.bytes_tx);
  p_usage.set_bytes_rx(credit_usage.bytes_rx);
  p_usage.set_type(update_type);

  // add the Requested-Service-Unit only if we are not on final grant
  RequestedUnits requestedUnits;
  if (is_final_grant) {
    requestedUnits.set_total(0);
    requestedUnits.set_tx(0);
    requestedUnits.set_rx(0);
  } else {
    requestedUnits = credit.get_requested_credits_units();
  }
  p_usage.mutable_requested_units()->CopyFrom(requestedUnits);
  return p_usage;
}

bool ChargingGrant::get_update_type(
    CreditUsage::UpdateType* update_type) const {
  if (credit.is_reporting()) {
    MLOG(MDEBUG) << "is_reporting is True , not sending update";
    return false;  // No update
  }
  if (reauth_state == REAUTH_REQUIRED) {
    *update_type = CreditUsage::REAUTH_REQUIRED;
    return true;
  }
  if (is_final_grant) {
    // Don't request updates if this is the final grant
    return false;
  }
  if (credit.is_quota_exhausted(SessionCredit::USAGE_REPORTING_THRESHOLD)) {
    *update_type = CreditUsage::QUOTA_EXHAUSTED;
    return true;
  }
  if (time(NULL) >= expiry_time) {
    *update_type = CreditUsage::VALIDITY_TIMER_EXPIRED;
    return true;
  }
  return false;
}

bool ChargingGrant::should_deactivate_service() const {
  if ((final_action_info.final_action ==
       ChargingCredit_FinalAction_TERMINATE) &&
      !SessionCredit::TERMINATE_SERVICE_WHEN_QUOTA_EXHAUSTED) {
    // configured in sessiond.yml
    return false;
  }
  if (service_state != SERVICE_ENABLED) {
    // service is not enabled
    return false;
  }
  if (is_final_grant && credit.is_quota_exhausted(1)) {
    // We only deactivate service when we receive a Final Unit
    // Indication (final Grant) and we've exhausted all quota
    MLOG(MINFO) << "Deactivating service because we have exhausted the given "
                << "quota and it is the final grant."
                << "action="
                << final_action_to_str(final_action_info.final_action);
    return true;
  }
  return false;
}

ServiceActionType ChargingGrant::get_action(SessionCreditUpdateCriteria& uc) {
  switch (service_state) {
    case SERVICE_NEEDS_DEACTIVATION:
      set_service_state(SERVICE_DISABLED, uc);
      if (!is_final_grant) {
        return TERMINATE_SERVICE;
      }
      return final_action_to_action(final_action_info.final_action);
    case SERVICE_NEEDS_ACTIVATION:
      set_service_state(SERVICE_ENABLED, uc);
      return ACTIVATE_SERVICE;
    default:
      return CONTINUE_SERVICE;
  }
}

ServiceActionType ChargingGrant::final_action_to_action(
    const ChargingCredit_FinalAction action) const {
  switch (action) {
    case ChargingCredit_FinalAction_REDIRECT:
      return REDIRECT;
    case ChargingCredit_FinalAction_RESTRICT_ACCESS:
      return RESTRICT_ACCESS;
    case ChargingCredit_FinalAction_TERMINATE:
    default:
      return TERMINATE_SERVICE;
  }
}

void ChargingGrant::set_reauth_state(
    const ReAuthState new_state, SessionCreditUpdateCriteria& uc) {
  if (reauth_state != new_state) {
    MLOG(MDEBUG) << "ReAuth state change from "
                 << reauth_state_to_str(reauth_state) << " to "
                 << reauth_state_to_str(new_state);
  }
  reauth_state    = new_state;
  uc.reauth_state = new_state;
}

void ChargingGrant::set_service_state(
    const ServiceState new_service_state, SessionCreditUpdateCriteria& uc) {
  if (service_state != new_service_state) {
    MLOG(MDEBUG) << "Service state change from "
                 << service_state_to_str(service_state) << " to "
                 << service_state_to_str(new_service_state);
  }
  service_state    = new_service_state;
  uc.service_state = new_service_state;
}

void ChargingGrant::log_final_action_info() const {
  std::string final_action = "";
  if (is_final_grant) {
    final_action += "final action: ";
    final_action += final_action_to_str(final_action_info.final_action);
    switch (final_action_info.final_action) {
      case ChargingCredit_FinalAction_REDIRECT:
        final_action += ", redirect_server: ";
        final_action +=
            final_action_info.redirect_server.redirect_server_address();
        break;
      case ChargingCredit_FinalAction_RESTRICT_ACCESS:
        final_action += ", restrict_rules: { ";
        for (auto rule : final_action_info.restrict_rules) {
          final_action += rule + " ";
        }
        final_action += "}";
        break;
      default:  // do nothing;
        break;
    }
  }
  MLOG(MINFO) << "This is a final credit, with " << final_action;
}

void ChargingGrant::set_reporting(bool reporting) {
  credit.set_reporting(reporting);
}
}  // namespace magma
