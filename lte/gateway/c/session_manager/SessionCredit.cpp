/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <limits>

#include "DiameterCodes.h"
#include "SessionCredit.h"
#include "magma_logging.h"

namespace magma {

float SessionCredit::USAGE_REPORTING_THRESHOLD = 0.8;
bool SessionCredit::TERMINATE_SERVICE_WHEN_QUOTA_EXHAUSTED = true;
std::string final_action_to_str(ChargingCredit_FinalAction final_action);
std::string service_state_to_str(ServiceState state);
std::string reauth_state_to_str(ReAuthState state);

std::unique_ptr<SessionCredit>
SessionCredit::unmarshal(const StoredSessionCredit &marshaled,
                         CreditType credit_type) {
  auto session_credit = std::make_unique<SessionCredit>(credit_type);

  session_credit->reporting_ = marshaled.reporting;
  session_credit->is_final_grant_ = marshaled.is_final;
  session_credit->unlimited_quota_ = marshaled.unlimited_quota;

  // FinalActionInfo
  FinalActionInfo final_action_info;
  final_action_info.final_action = marshaled.final_action_info.final_action;
  final_action_info.redirect_server =
      marshaled.final_action_info.redirect_server;
  session_credit->final_action_info_ = final_action_info;

  session_credit->reauth_state_ = marshaled.reauth_state;
  session_credit->service_state_ = marshaled.service_state;
  session_credit->expiry_time_ = marshaled.expiry_time;

  for (int bucket_int = USED_TX; bucket_int != MAX_VALUES; bucket_int++) {
    Bucket bucket = static_cast<Bucket>(bucket_int);
    if (marshaled.buckets.find(bucket) != marshaled.buckets.end()) {
      session_credit->buckets_[bucket] = marshaled.buckets.find(bucket)->second;
    }
  }

  session_credit->usage_reporting_limit_ = marshaled.usage_reporting_limit;

  return session_credit;
}

StoredSessionCredit SessionCredit::marshal() {
  StoredSessionCredit marshaled{};
  marshaled.reporting = reporting_;
  marshaled.is_final = is_final_grant_;
  marshaled.unlimited_quota = unlimited_quota_;

  marshaled.final_action_info.final_action = final_action_info_.final_action;
  marshaled.final_action_info.redirect_server =
      final_action_info_.redirect_server;

  marshaled.reauth_state = reauth_state_;
  marshaled.service_state = service_state_;
  marshaled.expiry_time = expiry_time_;

  for (int bucket_int = USED_TX; bucket_int != MAX_VALUES; bucket_int++) {
    Bucket bucket = static_cast<Bucket>(bucket_int);
    marshaled.buckets[bucket] = buckets_[bucket];
  }

  marshaled.usage_reporting_limit = usage_reporting_limit_;

  return marshaled;
}

SessionCreditUpdateCriteria SessionCredit::get_update_criteria() {
  SessionCreditUpdateCriteria uc{};
  uc.is_final = is_final_grant_;
  uc.final_action_info = final_action_info_;
  uc.reauth_state = reauth_state_;
  uc.service_state = service_state_;
  uc.expiry_time = expiry_time_;
  for (int bucket_int = USED_TX; bucket_int != MAX_VALUES; bucket_int++) {
    Bucket bucket = static_cast<Bucket>(bucket_int);
    uc.bucket_deltas[bucket] = 0;
  }
  return uc;
}

SessionCredit::SessionCredit(CreditType credit_type, ServiceState start_state)
    : credit_type_(credit_type), reporting_(false),
      reauth_state_(REAUTH_NOT_NEEDED), service_state_(start_state),
      unlimited_quota_(false), buckets_{}, is_final_grant_(false){}

SessionCredit::SessionCredit(CreditType credit_type, ServiceState start_state,
                             bool unlimited_quota)
    : credit_type_(credit_type), reporting_(false),
      reauth_state_(REAUTH_NOT_NEEDED), service_state_(start_state),
      unlimited_quota_(unlimited_quota), buckets_{}, is_final_grant_(false) {}


// by default, enable service
SessionCredit::SessionCredit(CreditType credit_type)
    : SessionCredit(credit_type, SERVICE_ENABLED, false) {}

void SessionCredit::set_expiry_time(uint32_t validity_time,
                                    SessionCreditUpdateCriteria &uc) {
  if (validity_time == 0) {
    // set as max possible time
    expiry_time_ = std::numeric_limits<std::time_t>::max();
    uc.expiry_time = expiry_time_;
    return;
  }
  expiry_time_ = std::time(nullptr) + validity_time;
  uc.expiry_time = expiry_time_;
}

void SessionCredit::add_used_credit(uint64_t used_tx, uint64_t used_rx,
                                    SessionCreditUpdateCriteria &uc) {
  buckets_[USED_TX] += used_tx;
  buckets_[USED_RX] += used_rx;
  uc.bucket_deltas[USED_TX] += used_tx;
  uc.bucket_deltas[USED_RX] += used_rx;

  log_quota_and_usage();
  if (should_deactivate_service()) {
    MLOG(MDEBUG) << "Quota exhausted. Deactivating service";
    set_service_state(SERVICE_NEEDS_DEACTIVATION, uc);
  }
}

void SessionCredit::reset_reporting_credit(
    SessionCreditUpdateCriteria &update_criteria) {
  buckets_[REPORTING_RX] = 0;
  buckets_[REPORTING_TX] = 0;
  reporting_ = false;
  update_criteria.reporting = false;
}

void SessionCredit::mark_failure(uint32_t code,
                                 SessionCreditUpdateCriteria &update_criteria) {
  if (DiameterCodeHandler::is_transient_failure(code)) {
    buckets_[REPORTED_RX] += buckets_[REPORTING_RX];
    buckets_[REPORTED_TX] += buckets_[REPORTING_TX];
    update_criteria.bucket_deltas[REPORTED_RX] += buckets_[REPORTING_RX];
    update_criteria.bucket_deltas[REPORTED_TX] += buckets_[REPORTING_TX];
  }
  reset_reporting_credit(update_criteria);
  if (should_deactivate_service()) {
    set_service_state(SERVICE_NEEDS_DEACTIVATION, update_criteria);
  }
}

void SessionCredit::receive_credit(
    uint64_t total_volume, uint64_t tx_volume, uint64_t rx_volume,
    uint32_t validity_time, bool is_final_grant,
    FinalActionInfo final_action_info,
    SessionCreditUpdateCriteria &update_criteria) {
  MLOG(MINFO) << "Received the following credit"
              << " total_volume=" << total_volume << " tx_volume=" << tx_volume
              << " rx_volume=" << rx_volume
              << " w/ validity time=" << validity_time;
  if (is_final_grant) {
    MLOG(MINFO) << "This credit received is the final grant, with final "
                << "action="
                << final_action_to_str(final_action_info.final_action);
  }

  buckets_[ALLOWED_TOTAL] += total_volume;
  buckets_[ALLOWED_TX] += tx_volume;
  buckets_[ALLOWED_RX] += rx_volume;
  update_criteria.bucket_deltas[ALLOWED_TOTAL] += total_volume;
  update_criteria.bucket_deltas[ALLOWED_TX] += tx_volume;
  update_criteria.bucket_deltas[ALLOWED_RX] += rx_volume;

  // transfer reporting usage to reported
  buckets_[REPORTED_RX] += buckets_[REPORTING_RX];
  buckets_[REPORTED_TX] += buckets_[REPORTING_TX];

  // Set the usage_reporting_limit so that we never report more than grant
  // we've received.
  update_criteria.bucket_deltas[REPORTED_RX] += buckets_[REPORTING_RX];
  update_criteria.bucket_deltas[REPORTED_TX] += buckets_[REPORTING_TX];
  auto reported_sum = buckets_[REPORTED_RX] + buckets_[REPORTED_TX];
  if (buckets_[ALLOWED_TOTAL] > reported_sum) {
    usage_reporting_limit_ = buckets_[ALLOWED_TOTAL] - reported_sum;
  } else if (usage_reporting_limit_ != 0) {
    MLOG(MINFO) << "We have reported data usage for all credit received, the "
                << "upper limit for reporting is now 0";
    usage_reporting_limit_ = 0;
    update_criteria.usage_reporting_limit = usage_reporting_limit_;
  }

  set_expiry_time(validity_time, update_criteria);
  reset_reporting_credit(update_criteria);

  set_is_final_grant_and_final_action(
      is_final_grant, final_action_info, update_criteria);

  if (reauth_state_ == REAUTH_PROCESSING) {
    set_reauth(REAUTH_NOT_NEEDED, update_criteria);
  }
  if (!is_quota_exhausted() && (service_state_ == SERVICE_DISABLED ||
                                service_state_ == SERVICE_NEEDS_DEACTIVATION)) {
    // if quota no longer exhausted, reenable services as needed
    MLOG(MINFO) << "Quota available. Activating service";
    set_service_state(SERVICE_NEEDS_ACTIVATION, update_criteria);
  }
  log_quota_and_usage();
}

void SessionCredit::log_quota_and_usage() const {
  auto reported_sum = buckets_[REPORTED_TX] + buckets_[REPORTED_RX];
  MLOG(MDEBUG) << "===> Used     Tx: " << buckets_[USED_TX]
               << " Rx: " << buckets_[USED_RX]
               << " Total: " << buckets_[USED_TX] + buckets_[USED_RX];
  MLOG(MDEBUG) << "===> Allowed  Tx: " << buckets_[ALLOWED_TX]
               << " Rx: " << buckets_[ALLOWED_RX]
               << " Total: " << buckets_[ALLOWED_TOTAL];
  MLOG(MDEBUG) << "===> Reported Tx: " << buckets_[REPORTED_TX]
               << " Rx: " << buckets_[REPORTED_RX]
               << " Total: " << buckets_[REPORTED_RX] + buckets_[REPORTED_TX];

  std::string final_action = "";
  if (is_final_grant_ ) {
    final_action += ", with final action: ";
    final_action +=  final_action_to_str(final_action_info_.final_action);
    if (final_action_info_.final_action ==
        ChargingCredit_FinalAction_REDIRECT) {
      final_action += ", redirect_server: ";
      final_action += final_action_info_.redirect_server.redirect_server_address();
    }
  }
  MLOG(MDEBUG) << "===> Is final grant: " << is_final_grant_
               << final_action;
}

bool SessionCredit::is_quota_exhausted(float usage_reporting_threshold) const {
  // used quota since last report
  uint64_t total_reported_usage = buckets_[REPORTED_TX] + buckets_[REPORTED_RX];
  uint64_t total_usage_since_report =
      std::max(uint64_t(0),
               buckets_[USED_TX] + buckets_[USED_RX] - total_reported_usage);
  uint64_t tx_usage_since_report =
      std::max(uint64_t(0), buckets_[USED_TX] - buckets_[REPORTED_TX]);
  uint64_t rx_usage_since_report =
      std::max(uint64_t(0), buckets_[USED_RX] - buckets_[REPORTED_RX]);

  // available quota since last report
  auto total_usage_reporting_threshold =
      std::max(0.0f, (buckets_[ALLOWED_TOTAL] - total_reported_usage) *
                usage_reporting_threshold);

  // reported tx/rx could be greater than allowed tx/rx
  // because some OCS/PCRF might not track tx/rx,
  // and 0 is added to the allowed credit when an credit update is received
  auto tx_usage_reporting_threshold =
      std::max(0.0f, (buckets_[ALLOWED_TX] - buckets_[REPORTED_TX]) *
                         usage_reporting_threshold);
  auto rx_usage_reporting_threshold =
      std::max(0.0f, (buckets_[ALLOWED_RX] - buckets_[REPORTED_RX]) *
                         usage_reporting_threshold);

  if (total_usage_since_report >= total_usage_reporting_threshold) {
    MLOG(MDEBUG) << "Total Quota exhausted";
    return true;
  } else if ((buckets_[ALLOWED_TX] > 0) &&
             (tx_usage_since_report >= tx_usage_reporting_threshold)) {
    MLOG(MDEBUG) << "Tx Quota exhausted";
    return true;
  } else if ((buckets_[ALLOWED_RX] > 0) &&
             (rx_usage_since_report >= rx_usage_reporting_threshold)) {
    MLOG(MDEBUG) << "Rx Quota exhausted";
    return true;
  }
  return false;
}

bool SessionCredit::should_deactivate_service() const {
  if (credit_type_ != CreditType::CHARGING) {
    // we only terminate on charging quota exhaustion
    return false;
  }
  if (unlimited_quota_) {
    return false;
  }
  if (!SessionCredit::TERMINATE_SERVICE_WHEN_QUOTA_EXHAUSTED) {
    // configured in sessiond.yml
    return false;
  }
  if (is_final_grant_ && is_quota_exhausted()) {
    // We only terminate when we receive a Final Unit Indication (final Grant)
    // and we've exhausted all quota
    MLOG(MINFO) << "Terminating service because we have exhausted the given "
                << "quota and it is the final grant";
    return true;
  }
  return false;
}

bool SessionCredit::validity_timer_expired() const {
  return time(NULL) >= expiry_time_;
}

CreditUpdateType SessionCredit::get_update_type() const {
  if (is_reporting()) {
    return CREDIT_NO_UPDATE;
  } else if (is_reauth_required()) {
    return CREDIT_REAUTH_REQUIRED;
  } else if (is_final_grant_ && is_quota_exhausted()) {
    // Don't request updates if there's no quota left
    return CREDIT_NO_UPDATE;
  } else if (is_quota_exhausted(SessionCredit::USAGE_REPORTING_THRESHOLD)) {
    return CREDIT_QUOTA_EXHAUSTED;
  } else if (validity_timer_expired()) {
    return CREDIT_VALIDITY_TIMER_EXPIRED;
  } else {
    return CREDIT_NO_UPDATE;
  }
}

SessionCredit::Usage SessionCredit::get_all_unreported_usage_for_reporting(
    SessionCreditUpdateCriteria &update_criteria) {
  auto usage = get_unreported_usage();
  buckets_[REPORTING_TX] += usage.bytes_tx;
  buckets_[REPORTING_RX] += usage.bytes_rx;
  reporting_ = true;
  update_criteria.reporting = true;

  log_usage_report(usage);
  return usage;
}

SessionCredit::Usage SessionCredit::get_usage_for_reporting(
    SessionCreditUpdateCriteria &update_criteria) {
  if (is_final_grant_) {
    return get_all_unreported_usage_for_reporting(update_criteria);
  }

  if (reauth_state_ == REAUTH_REQUIRED) {
    set_reauth(REAUTH_PROCESSING, update_criteria);
  }

  auto usage = get_unreported_usage();

  // Apply reporting limits since the user is not getting terminated.
  // The limits are applied on total usage (ie. tx + rx)
  usage.bytes_tx = std::min(usage.bytes_tx, usage_reporting_limit_);
  usage.bytes_rx =
      std::min(usage.bytes_rx, usage_reporting_limit_ - usage.bytes_tx);
  MLOG(MDEBUG) << "Since this is not the last report (final grant), we will "
               << "only report min(usage, usage_reporting_limit="
               << usage_reporting_limit_ << ")";

  buckets_[REPORTING_TX] += usage.bytes_tx;
  buckets_[REPORTING_RX] += usage.bytes_rx;
  reporting_ = true;
  update_criteria.reporting = true;

  log_usage_report(usage);
  return usage;
}

SessionCredit::Usage SessionCredit::get_unreported_usage() const {
  SessionCredit::Usage usage = {bytes_tx : 0, bytes_rx : 0};
  auto report = buckets_[REPORTED_TX] + buckets_[REPORTING_TX];
  if (buckets_[USED_TX] > report) {
    usage.bytes_tx = buckets_[USED_TX] - report;
  }
  report = buckets_[REPORTED_RX] - buckets_[REPORTING_RX];
  if (buckets_[USED_RX] > report) {
    usage.bytes_rx = buckets_[USED_RX] - report;
  }
  MLOG(MDEBUG) << "===> Data usage since last report is tx=" << usage.bytes_tx
               << " rx=" << usage.bytes_rx;
  return usage;
}

ServiceActionType
SessionCredit::get_action(SessionCreditUpdateCriteria &update_criteria) {
  if (service_state_ == SERVICE_NEEDS_DEACTIVATION) {
    set_service_state(SERVICE_DISABLED, update_criteria);
    return get_action_for_deactivating_service();
  } else if (service_state_ == SERVICE_NEEDS_ACTIVATION) {
    set_service_state(SERVICE_ENABLED, update_criteria);
    return ACTIVATE_SERVICE;
  }
  return CONTINUE_SERVICE;
}

ServiceActionType SessionCredit::get_action_for_deactivating_service() const {
  if (is_final_grant_ &&
      final_action_info_.final_action == ChargingCredit_FinalAction_REDIRECT) {
    return REDIRECT;
  } else if (is_final_grant_ &&
             final_action_info_.final_action ==
                 ChargingCredit_FinalAction_RESTRICT_ACCESS) {
    return RESTRICT_ACCESS;
  } else {
    return TERMINATE_SERVICE;
  }
}

bool SessionCredit::is_reporting() const { return reporting_; }

uint64_t SessionCredit::get_credit(Bucket bucket) const {
  return buckets_[bucket];
}

bool SessionCredit::is_reauth_required() const {
  return reauth_state_ == REAUTH_REQUIRED;
}

void SessionCredit::reauth(SessionCreditUpdateCriteria &update_criteria) {
  set_reauth(REAUTH_REQUIRED, update_criteria);
}

RedirectServer SessionCredit::get_redirect_server() const {
  return final_action_info_.redirect_server;
}

void SessionCredit::set_is_final_grant_and_final_action(
    bool is_final_grant, FinalActionInfo final_action_info,
    SessionCreditUpdateCriteria& update_criteria) {
  is_final_grant_          = is_final_grant;
  update_criteria.is_final = is_final_grant;
  final_action_info_                = final_action_info;
  update_criteria.final_action_info = final_action_info;
}

void SessionCredit::set_reauth(ReAuthState new_reauth_state,
                               SessionCreditUpdateCriteria &update_criteria) {
  if (reauth_state_ != new_reauth_state) {
    MLOG(MDEBUG) << "ReAuth state change from "
                 << reauth_state_to_str(reauth_state_) << " to "
                 << reauth_state_to_str(new_reauth_state);
  }
  reauth_state_ = new_reauth_state;
  update_criteria.reauth_state = new_reauth_state;
}

void SessionCredit::set_service_state(
    ServiceState new_service_state,
    SessionCreditUpdateCriteria &update_criteria) {
  if (service_state_ != new_service_state) {
    MLOG(MDEBUG) << "Service state change from "
                 << service_state_to_str(service_state_) << " to "
                 << service_state_to_str(new_service_state);
  }
  service_state_ = new_service_state;
  update_criteria.service_state = new_service_state;
}

void SessionCredit::set_expiry_time(
    std::time_t expiry_time, SessionCreditUpdateCriteria &update_criteria) {
  expiry_time_ = expiry_time;
  update_criteria.expiry_time = expiry_time;
}

void SessionCredit::add_credit(uint64_t credit, Bucket bucket,
                               SessionCreditUpdateCriteria &update_criteria) {
  buckets_[bucket] += credit;
  update_criteria.bucket_deltas[bucket] += credit;
}

void SessionCredit::log_usage_report(SessionCredit::Usage usage) const {
  MLOG(MDEBUG) << "===> Amount reporting for this report:"
               << " tx=" << usage.bytes_tx << " rx=" << usage.bytes_rx;
  MLOG(MDEBUG) << "===> The total amount currently being reported:"
               << " tx=" << buckets_[REPORTING_TX]
               << " rx=" << buckets_[REPORTING_RX];
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
    return "";
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
  default:
    return "INVALID SERVICE STATE";
  }
}

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
} // namespace magma
