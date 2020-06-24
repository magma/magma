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
#include "EnumToString.h"
#include "SessionCredit.h"
#include "magma_logging.h"

namespace magma {

float SessionCredit::USAGE_REPORTING_THRESHOLD = 0.8;
bool SessionCredit::TERMINATE_SERVICE_WHEN_QUOTA_EXHAUSTED = true;

SessionCredit SessionCredit::unmarshal(
  const StoredSessionCredit &marshaled, CreditType credit_type) {
  SessionCredit session_credit(credit_type);

  session_credit.reporting_ = marshaled.reporting;
  session_credit.is_final_grant_ = marshaled.is_final;
  session_credit.credit_limit_type_ = marshaled.credit_limit_type;

  // FinalActionInfo
  FinalActionInfo final_action_info;
  final_action_info.final_action = marshaled.final_action_info.final_action;
  final_action_info.redirect_server =
      marshaled.final_action_info.redirect_server;
  session_credit.final_action_info_ = final_action_info;

  session_credit.reauth_state_ = marshaled.reauth_state;
  session_credit.service_state_ = marshaled.service_state;
  session_credit.expiry_time_ = marshaled.expiry_time;
  session_credit.grant_tracking_type_ = marshaled.grant_tracking_type;

  for (int bucket_int = USED_TX; bucket_int != MAX_VALUES; bucket_int++) {
    Bucket bucket = static_cast<Bucket>(bucket_int);
    if (marshaled.buckets.find(bucket) != marshaled.buckets.end()) {
      session_credit.buckets_[bucket] = marshaled.buckets.find(bucket)->second;
    }
  }
  return session_credit;
}

StoredSessionCredit SessionCredit::marshal() {
  StoredSessionCredit marshaled{};
  marshaled.reporting = reporting_;
  marshaled.is_final = is_final_grant_;
  marshaled.credit_limit_type = credit_limit_type_;

  marshaled.final_action_info.final_action = final_action_info_.final_action;
  marshaled.final_action_info.redirect_server =
      final_action_info_.redirect_server;

  marshaled.reauth_state = reauth_state_;
  marshaled.service_state = service_state_;
  marshaled.expiry_time = expiry_time_;
  marshaled.grant_tracking_type = grant_tracking_type_;

  for (int bucket_int = USED_TX; bucket_int != MAX_VALUES; bucket_int++) {
    Bucket bucket = static_cast<Bucket>(bucket_int);
    marshaled.buckets[bucket] = buckets_[bucket];
  }
  return marshaled;
}

SessionCreditUpdateCriteria SessionCredit::get_update_criteria() {
  SessionCreditUpdateCriteria uc{};
  uc.is_final = is_final_grant_;
  uc.final_action_info = final_action_info_;
  uc.reauth_state = reauth_state_;
  uc.service_state = service_state_;
  uc.expiry_time = expiry_time_;
  uc.grant_tracking_type = grant_tracking_type_;
  for (int bucket_int = USED_TX; bucket_int != MAX_VALUES; bucket_int++) {
    Bucket bucket = static_cast<Bucket>(bucket_int);
    uc.bucket_deltas[bucket] = 0;
  }
  return uc;
}

SessionCredit::SessionCredit(CreditType credit_type, ServiceState start_state)
    : credit_type_(credit_type), reauth_state_(REAUTH_NOT_NEEDED),
      service_state_(start_state), buckets_{}, reporting_(false),
      credit_limit_type_(FINITE), is_final_grant_(false) {}

SessionCredit::SessionCredit(
  CreditType credit_type, ServiceState start_state,
  CreditLimitType credit_limit_type)
    : credit_type_(credit_type), reauth_state_(REAUTH_NOT_NEEDED),
      service_state_(start_state), buckets_{}, reporting_(false),
      credit_limit_type_(credit_limit_type), is_final_grant_(false) {}

// by default, enable service & finite credit
SessionCredit::SessionCredit(CreditType credit_type)
    : SessionCredit(credit_type, SERVICE_ENABLED, FINITE) {}

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

void SessionCredit::receive_credit(const GrantedUnits& gsu,
    uint32_t validity_time, bool is_final_grant, FinalActionInfo final_action,
    SessionCreditUpdateCriteria& uc) {
  grant_tracking_type_ = determine_grant_tracking_type(gsu);
  uc.grant_tracking_type = grant_tracking_type_;

  // Clear invalid values
  uint64_t total_volume = gsu.total().is_valid() ? gsu.total().volume() : 0;
  uint64_t tx_volume = gsu.tx().is_valid() ? gsu.tx().volume() : 0;
  uint64_t rx_volume = gsu.rx().is_valid() ? gsu.rx().volume() : 0;

  // Update allowed bytes
  buckets_[ALLOWED_TOTAL] += total_volume;
  buckets_[ALLOWED_TX] += tx_volume;
  buckets_[ALLOWED_RX] += rx_volume;
  uc.bucket_deltas[ALLOWED_TOTAL] += total_volume;
  uc.bucket_deltas[ALLOWED_TX] += tx_volume;
  uc.bucket_deltas[ALLOWED_RX] += rx_volume;

  MLOG(MINFO)  << "Received the following credit"
              << " total_volume=" << total_volume
              << " tx_volume=" << tx_volume
              << " rx_volume=" << rx_volume
              << " grant_tracking_type="
              << grant_type_to_str(grant_tracking_type_)
              << " w/ validity time=" << validity_time;
  if (is_final_grant) {
    MLOG(MINFO) << "This credit received is the final grant, with final "
                << "action="
                << final_action_to_str(final_action.final_action);
  }

  // Update reporting/reported bytes
  buckets_[REPORTED_RX] += buckets_[REPORTING_RX];
  buckets_[REPORTED_TX] += buckets_[REPORTING_TX];
  uc.bucket_deltas[REPORTED_RX] += buckets_[REPORTING_RX];
  uc.bucket_deltas[REPORTED_TX] += buckets_[REPORTING_TX];
  reset_reporting_credit(uc);

  set_expiry_time(validity_time, uc);
  set_is_final_grant_and_final_action(is_final_grant, final_action, uc);

  if (reauth_state_ == REAUTH_PROCESSING) {
    set_reauth(REAUTH_NOT_NEEDED, uc);
  }
  if (!is_quota_exhausted(1) && (service_state_ == SERVICE_DISABLED ||
                                 service_state_ == SERVICE_REDIRECTED ||
                                 service_state_ == SERVICE_NEEDS_DEACTIVATION)) {
    // if quota no longer exhausted, reenable services as needed
    MLOG(MINFO) << "Quota available. Activating service";
    set_service_state(SERVICE_NEEDS_ACTIVATION, uc);
  }
  log_quota_and_usage();
}

bool SessionCredit::is_quota_exhausted(float threshold) const {
  if (credit_limit_type_ != FINITE) {
    return false;
  }
  uint64_t total_reported = buckets_[REPORTED_TX] + buckets_[REPORTED_RX];
  uint64_t total_usage = buckets_[USED_TX] + buckets_[USED_RX];

  bool rx_exhausted = compute_quota_exhausted(buckets_[ALLOWED_RX],
    buckets_[REPORTED_RX], buckets_[USED_RX], threshold);
  bool tx_exhausted = compute_quota_exhausted(buckets_[ALLOWED_TX],
    buckets_[REPORTED_TX], buckets_[USED_TX], threshold);

  bool is_exhausted = false;
  switch (grant_tracking_type_) {
    case RX_ONLY:
      is_exhausted = rx_exhausted;
      break;
    case TX_ONLY:
      is_exhausted = tx_exhausted;
      break;
    case TX_AND_RX:
      is_exhausted = rx_exhausted || tx_exhausted;
      break;
    case TOTAL_ONLY:
      is_exhausted = compute_quota_exhausted(
        buckets_[ALLOWED_TOTAL], total_reported, total_usage, threshold);
      break;
    default:
      MLOG(MERROR) << "Invalid grant_tracking_type="
                   << grant_type_to_str(grant_tracking_type_);
      return false;
  }
  if (is_exhausted) {
    MLOG(MDEBUG) << grant_type_to_str(grant_tracking_type_)
                 << " grant is exhausted";
  }
  return is_exhausted;
}

bool SessionCredit::should_deactivate_service() const {
  if (credit_type_ != CreditType::CHARGING) {
    // we only terminate on charging quota exhaustion
    return false;
  }
  if (credit_limit_type_ != FINITE) {
    return false;
  }
  if ((final_action_info_.final_action ==
        ChargingCredit_FinalAction_TERMINATE) &&
      !SessionCredit::TERMINATE_SERVICE_WHEN_QUOTA_EXHAUSTED) {
      // configured in sessiond.yml
      return false;
  }
  if (service_state_ != SERVICE_ENABLED){
    // service is not enabled
    return false;
  }
  if (is_final_grant_ && is_quota_exhausted(1)) {
    // We only deactivate service when we receive a Final Unit
    // Indication (final Grant) and we've exhausted all quota
    MLOG(MINFO) << "Deactivating service because we have exhausted the given "
      << "quota and it is the final grant."
      << "action=" << final_action_to_str(final_action_info_.final_action);
    return true;
  }
  return false;
}

bool SessionCredit::is_service_redirected() const {
  return service_state_ == SERVICE_REDIRECTED;
}

bool SessionCredit::validity_timer_expired() const {
  return time(NULL) >= expiry_time_;
}

CreditUpdateType SessionCredit::get_update_type() const {
  if (is_reporting()) {
    return CREDIT_NO_UPDATE;
  } else if (is_reauth_required()) {
    return CREDIT_REAUTH_REQUIRED;
  } else if (is_final_grant_ && is_quota_exhausted(1)) {
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

// Todo Return Proto CreditUsage
SessionCredit::Usage SessionCredit::get_usage_for_reporting(
    SessionCreditUpdateCriteria &update_criteria) {
  if (reauth_state_ == REAUTH_REQUIRED) {
    set_reauth(REAUTH_PROCESSING, update_criteria);
  }

  if (is_final_grant_) {
    return get_all_unreported_usage_for_reporting(update_criteria);
  }

  auto usage = get_unreported_usage();
  // Apply reporting limits since the user is not getting terminated.
  // We never want to report more than the amount we've received
  apply_reporting_limits(usage);

  buckets_[REPORTING_TX] += usage.bytes_tx;
  buckets_[REPORTING_RX] += usage.bytes_rx;
  reporting_ = true;
  update_criteria.reporting = true;

  log_usage_report(usage);
  return usage;
}

// Take the minimum of (grant - reported) and (reported - used)
void SessionCredit::apply_reporting_limits(SessionCredit::Usage& usage) {
  uint64_t tx_limit, rx_limit, total_limit, total_reported;

  tx_limit = compute_reporting_limit(
    buckets_[ALLOWED_TX], buckets_[REPORTED_TX]);
  rx_limit = compute_reporting_limit(
    buckets_[ALLOWED_TX], buckets_[REPORTED_TX]);

  switch (grant_tracking_type_) {
    case TX_ONLY:
      usage.bytes_tx = std::min(usage.bytes_tx, tx_limit);
      usage.bytes_rx = 0; // Clear field that doesn't need to be reported
      MLOG(MDEBUG) << "Applying a TX reporting limit of " << tx_limit;
      break;
    case RX_ONLY:
      usage.bytes_rx = std::min(usage.bytes_rx, rx_limit);
      usage.bytes_tx = 0; // Clear field that doesn't need to be reported
      MLOG(MDEBUG) << "Applying a RX reporting limit of " << rx_limit;
      break;
    case TX_AND_RX:
      usage.bytes_tx = std::min(usage.bytes_tx, tx_limit);
      usage.bytes_rx = std::min(usage.bytes_rx, rx_limit);
      MLOG(MDEBUG) << "Applying TX and RX reporting limits of "
                   << tx_limit << ", " << rx_limit;
      break;
    case TOTAL_ONLY:
      total_reported = buckets_[REPORTED_RX] + buckets_[REPORTED_TX];
      total_limit = compute_reporting_limit(
        buckets_[ALLOWED_TOTAL], total_reported);
      usage.bytes_tx = std::min(usage.bytes_tx, total_limit);
      usage.bytes_rx = std::min(usage.bytes_rx, total_limit - usage.bytes_tx);
      MLOG(MDEBUG) << "Applying a total reporting limit of " << total_limit;
      break;
    default:
      MLOG(MERROR) << "Credit with this tracking type is probably unexpected: "
                   << grant_tracking_type_;
  }
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

uint64_t SessionCredit::compute_reporting_limit(const uint64_t allowed,
    const uint64_t reported) const {
  uint64_t limit = 0;
  if (allowed > reported) {
    limit = allowed - reported;
  }
  return limit;
}

bool SessionCredit::compute_quota_exhausted(const uint64_t allowed,
  const uint64_t reported, const uint64_t used, float threshold_ratio) const{
  uint64_t unreported_usage = 0;
  if (used > reported) {
    unreported_usage = used - reported;
  }
  uint64_t grant_left = allowed - reported;
  auto threshold = std::max(0.0f, grant_left * threshold_ratio);
  return unreported_usage >= threshold;
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
  if (!is_final_grant_) {
    return TERMINATE_SERVICE;
  }
  // TODO look into just reusing the proto defined enum
  switch (final_action_info_.final_action) {
    case ChargingCredit_FinalAction_REDIRECT:
      return REDIRECT;
    case ChargingCredit_FinalAction_RESTRICT_ACCESS:
      return RESTRICT_ACCESS;
    case ChargingCredit_FinalAction_TERMINATE:
    default:
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

void SessionCredit::set_grant_tracking_type(GrantTrackingType g_type,
  SessionCreditUpdateCriteria& uc) {
  grant_tracking_type_ = g_type;
  uc.grant_tracking_type = g_type;
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

// Determine the grant's tracking type by looking at which values are valid.
GrantTrackingType SessionCredit::determine_grant_tracking_type(
  const GrantedUnits& grant) {
  if (grant.total().is_valid()) {
    return TOTAL_ONLY;
  }
  bool tx_valid = grant.tx().is_valid();
  bool rx_valid = grant.rx().is_valid();
  if (tx_valid && rx_valid) {
    return TX_AND_RX;
  } else if (tx_valid) {
    return TX_ONLY;
  } else if (rx_valid) {
    return RX_ONLY;
  } else {
    MLOG(MWARNING) << "Received a GSU with no valid grants";
    return NO_TRACKING;
  }
}

void SessionCredit::log_quota_and_usage() const {
  MLOG(MDEBUG) << "===> Used     Tx: " << buckets_[USED_TX]
               << " Rx: " << buckets_[USED_RX]
               << " Total: " << buckets_[USED_TX] + buckets_[USED_RX];
  MLOG(MDEBUG) << "===> Allowed  Tx: " << buckets_[ALLOWED_TX]
               << " Rx: " << buckets_[ALLOWED_RX]
               << " Total: " << buckets_[ALLOWED_TOTAL];
  MLOG(MDEBUG) << "===> Reported Tx: " << buckets_[REPORTED_TX]
               << " Rx: " << buckets_[REPORTED_RX]
               << " Total: " << buckets_[REPORTED_RX] + buckets_[REPORTED_TX];
  MLOG(MDEBUG) << "===> Grant tracking type "
               << grant_type_to_str(grant_tracking_type_);

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

void SessionCredit::log_usage_report(SessionCredit::Usage usage) const {
  MLOG(MDEBUG) << "===> Amount reporting for this report:"
               << " tx=" << usage.bytes_tx << " rx=" << usage.bytes_rx;
  MLOG(MDEBUG) << "===> The total amount currently being reported:"
               << " tx=" << buckets_[REPORTING_TX]
               << " rx=" << buckets_[REPORTING_RX];
}
} // namespace magma
