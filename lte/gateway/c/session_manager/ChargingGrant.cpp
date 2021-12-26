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

#include <glog/logging.h>
#include <algorithm>
#include <ctime>
#include <limits>
#include <sstream>
#include <string>
#include <vector>

#include "ChargingGrant.h"
#include "CreditKey.h"
#include "DiameterCodes.h"
#include "EnumToString.h"
#include "magma_logging.h"

namespace magma {
ChargingGrant::ChargingGrant(const StoredChargingGrant& marshaled) {
  credit = SessionCredit(marshaled.credit);

  final_action_info.final_action = marshaled.final_action_info.final_action;
  final_action_info.redirect_server =
      marshaled.final_action_info.redirect_server;
  final_action_info.restrict_rules = marshaled.final_action_info.restrict_rules;

  reauth_state = marshaled.reauth_state;
  service_state = marshaled.service_state;
  expiry_time = marshaled.expiry_time;
  is_final_grant = marshaled.is_final;
  suspended = marshaled.suspended;
}

StoredChargingGrant ChargingGrant::marshal() {
  StoredChargingGrant marshaled{};
  marshaled.is_final = is_final_grant;
  marshaled.final_action_info.final_action = final_action_info.final_action;
  marshaled.final_action_info.redirect_server =
      final_action_info.redirect_server;
  marshaled.final_action_info.restrict_rules = final_action_info.restrict_rules;
  marshaled.reauth_state = reauth_state;
  marshaled.service_state = service_state;
  marshaled.expiry_time = expiry_time;
  marshaled.credit = credit.marshal();
  marshaled.suspended = suspended;
  return marshaled;
}

CreditValidity ChargingGrant::get_credit_response_validity_type(
    const CreditUpdateResponse& update) {
  const uint32_t key = update.charging_key();
  const std::string session_id = update.session_id();
  CreditValidity credit_validity = VALID_CREDIT;
  if (!update.success()) {
    if (DiameterCodeHandler::is_permanent_failure(update.result_code())) {
      MLOG(MERROR) << "Credit update failed RG:" << key
                   << " code:" << update.result_code() << " for " << session_id;
      return INVALID_CREDIT;
    } else if (DiameterCodeHandler::is_transient_failure(
                   update.result_code())) {
      MLOG(MDEBUG) << " Received a transient failure update for RG "
                   << update.charging_key() << ". Continuing service";
      credit_validity = TRANSIENT_ERROR;

    } else {
      MLOG(MDEBUG) << " Received an unknown on update for RG "
                   << update.charging_key() << ". Discarding";
      return INVALID_CREDIT;
    }
  }
  // For infinite credit, we do not care about the GSU value
  if (update.limit_type() == INFINITE_UNMETERED ||
      update.limit_type() == INFINITE_METERED) {
    return credit_validity;
  }
  const auto& gsu = update.credit().granted_units();
  bool gsu_all_invalid =
      !gsu.total().is_valid() && !gsu.rx().is_valid() && !gsu.tx().is_valid();
  if (gsu_all_invalid) {
    if (update.credit().is_final() || credit_validity == TRANSIENT_ERROR) {
      // TODO(@themarwhal): look into this case. Before I figure it out, I will
      // allow empty GSU credits with FUA or to be suspended
      // to be on the conservative side.
      MLOG(MWARNING)
          << "GSU for RG: " << key << " " << session_id
          << " is invalid, but accepting it as it has a final unit action or "
             "suspended credit ";
      return credit_validity;
    }
    MLOG(MERROR) << "Credit update failed RG:" << key
                 << " invalid, empty GSU and no FUA for " << session_id;
    return INVALID_CREDIT;
  }
  return credit_validity;
}

void ChargingGrant::receive_charging_grant(
    const CreditUpdateResponse& update,
    SessionCreditUpdateCriteria* credit_uc) {
  auto p_credit = update.credit();
  credit.receive_credit(p_credit.granted_units(), credit_uc);

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
  }

  // Expiry Time
  const auto delta_time_sec = p_credit.validity_time();
  if (delta_time_sec == 0) {
    expiry_time = std::numeric_limits<std::time_t>::max();
  } else {
    expiry_time = std::time(nullptr) + delta_time_sec;
  }
  log_received_grant(update);

  // Update the UpdateCriteria if not nullptr
  if (credit_uc) {
    credit_uc->is_final = is_final_grant;
    credit_uc->final_action_info = final_action_info;
    credit_uc->expiry_time = expiry_time;
    credit_uc->suspended = suspended;
  }
}

SessionCreditUpdateCriteria ChargingGrant::get_update_criteria() {
  SessionCreditUpdateCriteria credit_uc = credit.get_update_criteria();
  credit_uc.is_final = is_final_grant;
  credit_uc.final_action_info = final_action_info;
  credit_uc.expiry_time = expiry_time;
  credit_uc.reauth_state = reauth_state;
  credit_uc.service_state = service_state;
  credit_uc.suspended = suspended;
  return credit_uc;
}

CreditUsage ChargingGrant::get_credit_usage(
    CreditUsage::UpdateType update_type, SessionCreditUpdateCriteria* credit_uc,
    bool is_terminate) {
  CreditUsage p_usage;
  Usage credit_usage;

  if (is_final_grant || is_terminate) {
    credit_usage = credit.get_all_unreported_usage_for_reporting(credit_uc);
  } else {
    credit_usage = credit.get_usage_for_reporting(credit_uc);
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
  if (time(nullptr) >= expiry_time) {
    *update_type = CreditUsage::VALIDITY_TIMER_EXPIRED;
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
  return false;
}

bool ChargingGrant::should_deactivate_service() const {
  const bool final_action_is_terminate =
      final_action_info.final_action == ChargingCredit_FinalAction_TERMINATE;
  const bool is_final_and_credit_exhausted =
      is_final_grant && credit.is_quota_exhausted(1);

  if (!is_final_and_credit_exhausted) {
    return false;
  }

  if (final_action_is_terminate &&
      !SessionCredit::TERMINATE_SERVICE_WHEN_QUOTA_EXHAUSTED) {
    // configured in sessiond.yml
    return false;
  }

  // 1. Disable if the credit is out of quota and final action is terminate
  // 2. Disable if the credit is out of quota and final action is
  //    redirect/restrict IF if it hasn't been acted on already (state is
  //    ENABLED)
  if (final_action_is_terminate || service_state == SERVICE_ENABLED) {
    MLOG(MINFO) << "Deactivating service because we have exhausted the given "
                << "quota and it is the final grant."
                << "action="
                << final_action_to_str(final_action_info.final_action);
    return true;
  }
  return false;
}

ServiceActionType ChargingGrant::get_action(
    SessionCreditUpdateCriteria* credit_uc) {
  switch (service_state) {
    case SERVICE_NEEDS_DEACTIVATION:
      set_service_state(SERVICE_DISABLED, credit_uc);
      if (!is_final_grant) {
        return TERMINATE_SERVICE;
      }
      return final_action_to_action(final_action_info.final_action);
    case SERVICE_NEEDS_ACTIVATION:
      set_service_state(SERVICE_ENABLED, credit_uc);
      return ACTIVATE_SERVICE;
    case SERVICE_NEEDS_SUSPENSION:
      set_service_state(SERVICE_DISABLED, credit_uc);
      return final_action_to_action_on_suspension(
          final_action_info.final_action);
    default:
      return CONTINUE_SERVICE;
  }
}

bool ChargingGrant::should_be_unsuspended() const {
  // transitioning out of FUA-redirect/restrict
  if (service_state == SERVICE_NEEDS_ACTIVATION) {
    return true;
  }
  // transitioning out of credit suspension
  if (suspended && !credit.is_quota_exhausted(1)) {
    return true;
  }
  return false;
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

ServiceActionType ChargingGrant::final_action_to_action_on_suspension(
    const ChargingCredit_FinalAction action) const {
  switch (action) {
    case ChargingCredit_FinalAction_REDIRECT:
      return REDIRECT;
    case ChargingCredit_FinalAction_RESTRICT_ACCESS:
      return RESTRICT_ACCESS;
    case ChargingCredit_FinalAction_TERMINATE:
    default:
      return CONTINUE_SERVICE;
  }
}

void ChargingGrant::set_reauth_state(const ReAuthState new_state,
                                     SessionCreditUpdateCriteria* credit_uc) {
  if (reauth_state != new_state) {
    MLOG(MDEBUG) << "ReAuth state change from "
                 << reauth_state_to_str(reauth_state) << " to "
                 << reauth_state_to_str(new_state);
  }
  reauth_state = new_state;
  if (credit_uc) {
    credit_uc->reauth_state = new_state;
  }
}

void ChargingGrant::set_service_state(const ServiceState new_service_state,
                                      SessionCreditUpdateCriteria* credit_uc) {
  if (service_state != new_service_state) {
    MLOG(MDEBUG) << "Service state change from "
                 << service_state_to_str(service_state) << " to "
                 << service_state_to_str(new_service_state);
  }
  service_state = new_service_state;
  if (credit_uc) {
    credit_uc->service_state = new_service_state;
  }
}

void ChargingGrant::set_suspended(bool new_suspended,
                                  SessionCreditUpdateCriteria* credit_uc) {
  if (suspended != new_suspended) {
    MLOG(MDEBUG) << "Credit suspension set to: " << new_suspended;
  }
  suspended = new_suspended;
  if (credit_uc) {
    credit_uc->suspended = new_suspended;
  }
}

void ChargingGrant::reset_reporting_grant(
    SessionCreditUpdateCriteria* credit_uc) {
  credit.reset_reporting_credit(credit_uc);
  if (reauth_state == REAUTH_PROCESSING) {
    set_reauth_state(REAUTH_REQUIRED, credit_uc);
  }
}

void ChargingGrant::log_received_grant(const CreditUpdateResponse& update) {
  std::ostringstream log;
  log << update.session_id() << " received a credit " << CreditKey(update);
  if (is_final_grant) {
    log << " with final action "
        << final_action_to_str(final_action_info.final_action);
    switch (final_action_info.final_action) {
      case ChargingCredit_FinalAction_REDIRECT:
        log << ", redirect_server: "
            << final_action_info.redirect_server.redirect_server_address();
        break;
      case ChargingCredit_FinalAction_RESTRICT_ACCESS:
        log << ", restrict_rules: { ";
        for (auto rule : final_action_info.restrict_rules) {
          log << (rule + " ");
        }
        log << "}";
        break;
      default:  // do nothing;
        break;
    }
  }
  if (update.credit().validity_time() != 0) {
    log << " with expiry timer in " << update.credit().validity_time()
        << " seconds";
  }
  MLOG(MINFO) << log.str();
}

void ChargingGrant::set_reporting(bool reporting) {
  credit.set_reporting(reporting);
}

// TODO: make session_manager.proto and policydb.proto to use common field
static RedirectInformation_AddressType address_type_converter(
    RedirectServer_RedirectAddressType address_type) {
  switch (address_type) {
    case RedirectServer_RedirectAddressType_IPV4:
      return RedirectInformation_AddressType_IPv4;
    case RedirectServer_RedirectAddressType_IPV6:
      return RedirectInformation_AddressType_IPv6;
    case RedirectServer_RedirectAddressType_URL:
      return RedirectInformation_AddressType_URL;
    case RedirectServer_RedirectAddressType_SIP_URI:
      return RedirectInformation_AddressType_SIP_URI;
    default:
      MLOG(MWARNING) << "Unknown redirect address type!";
      return RedirectInformation_AddressType_IPv4;
  }
}

PolicyRule ChargingGrant::make_redirect_rule() {
  PolicyRule redirect_rule;
  redirect_rule.set_id("redirect");
  redirect_rule.set_priority(ChargingGrant::REDIRECT_FLOW_PRIORITY);
  RedirectInformation* redirect_info = redirect_rule.mutable_redirect();
  redirect_info->set_support(RedirectInformation_Support_ENABLED);

  auto redirect_server = final_action_info.redirect_server;
  redirect_info->set_address_type(
      address_type_converter(redirect_server.redirect_address_type()));
  redirect_info->set_server_address(redirect_server.redirect_server_address());
  return redirect_rule;
}
}  // namespace magma
