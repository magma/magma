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
uint64_t SessionCredit::EXTRA_QUOTA_MARGIN = 1024;
bool SessionCredit::TERMINATE_SERVICE_WHEN_QUOTA_EXHAUSTED = true;

std::unique_ptr<SessionCredit> SessionCredit::unmarshal(
  const StoredSessionCredit &marshaled,
  CreditType credit_type)
{
  auto session_credit = std::make_unique<SessionCredit>(credit_type);

  session_credit->reporting_ = marshaled.reporting;
  session_credit->is_final_ = marshaled.is_final;
  session_credit->unlimited_quota_ = marshaled.unlimited_quota;

  // FinalActionInfo
  FinalActionInfo final_action_info;
  final_action_info.final_action = marshaled.final_action_info.final_action;
  final_action_info.redirect_server = marshaled.final_action_info.redirect_server;
  session_credit->final_action_info_ = final_action_info;

  session_credit->reauth_state_ = marshaled.reauth_state;
  session_credit->service_state_ = marshaled.service_state;
  session_credit->expiry_time_ = marshaled.expiry_time;

  for ( int bucket_int = USED_TX; bucket_int != MAX_VALUES; bucket_int++ )
  {
    Bucket bucket = static_cast<Bucket>(bucket_int);
    if (marshaled.buckets.find(bucket) != marshaled.buckets.end()) {
      session_credit->buckets_[bucket] = marshaled.buckets.find(bucket)->second;
    }
  }

  session_credit->usage_reporting_limit_ = marshaled.usage_reporting_limit;

  return session_credit;
}

StoredSessionCredit SessionCredit::marshal()
{
  StoredSessionCredit marshaled{};
  marshaled.reporting = reporting_;
  marshaled.is_final = is_final_;
  marshaled.unlimited_quota = unlimited_quota_;

  marshaled.final_action_info.final_action = final_action_info_.final_action;
  marshaled.final_action_info.redirect_server = final_action_info_.redirect_server;

  marshaled.reauth_state = reauth_state_;
  marshaled.service_state = service_state_;
  marshaled.expiry_time = expiry_time_;

  for ( int bucket_int = USED_TX; bucket_int != MAX_VALUES; bucket_int++ )
  {
    Bucket bucket = static_cast<Bucket>(bucket_int);
    marshaled.buckets[bucket] = buckets_[bucket];
  }

  usage_reporting_limit_ = marshaled.usage_reporting_limit;

  return marshaled;
}

SessionCredit::SessionCredit(CreditType credit_type, ServiceState start_state):
  credit_type_(credit_type),
  reporting_(false),
  reauth_state_(REAUTH_NOT_NEEDED),
  service_state_(start_state),
  unlimited_quota_(false),
  buckets_ {}
{
}

SessionCredit::SessionCredit(
  CreditType credit_type,
  ServiceState start_state,
  bool unlimited_quota):
  credit_type_(credit_type),
  reporting_(false),
  reauth_state_(REAUTH_NOT_NEEDED),
  service_state_(start_state),
  unlimited_quota_(unlimited_quota),
  buckets_ {}
{
}

// by default, enable service
SessionCredit::SessionCredit(CreditType credit_type):
  SessionCredit(credit_type, SERVICE_ENABLED, false)
{
}

void SessionCredit::set_expiry_time(uint32_t validity_time)
{
  if (validity_time == 0) {
    // set as max possible time
    expiry_time_ = std::numeric_limits<std::time_t>::max();
    return;
  }
  expiry_time_ = std::time(nullptr) + validity_time;
}

void SessionCredit::add_used_credit(uint64_t used_tx, uint64_t used_rx)
{
  buckets_[USED_TX] += used_tx;
  buckets_[USED_RX] += used_rx;

  if (should_deactivate_service()) {
    MLOG(MDEBUG) << "Quota exhausted. Deactivating service";
    service_state_ = SERVICE_NEEDS_DEACTIVATION;
  }
}

void SessionCredit::reset_reporting_credit()
{
  buckets_[REPORTING_RX] = 0;
  buckets_[REPORTING_TX] = 0;
  reporting_ = false;
}

void SessionCredit::mark_failure(uint32_t code)
{
  if (DiameterCodeHandler::is_transient_failure(code)) {
    buckets_[REPORTED_RX] += buckets_[REPORTING_RX];
    buckets_[REPORTED_TX] += buckets_[REPORTING_TX];
  }
  reset_reporting_credit();
  if (should_deactivate_service()) {
    service_state_ = SERVICE_NEEDS_DEACTIVATION;
  }
}

void SessionCredit::receive_credit(
  uint64_t total_volume,
  uint64_t tx_volume,
  uint64_t rx_volume,
  uint32_t validity_time,
  bool is_final,
  FinalActionInfo final_action_info)
{
  MLOG(MDEBUG) << "receive_credit:"
               << "total allowed octets: " << buckets_[ALLOWED_TOTAL]
               << "total_tx allowed: " << buckets_[ALLOWED_TX]
               << "total_rx allowed: " << buckets_[ALLOWED_RX]
               << "adding total: " << total_volume
               << "adding tx: " << tx_volume
               << "adding rx: " << rx_volume;
  buckets_[ALLOWED_TOTAL] += total_volume;
  buckets_[ALLOWED_TX] += tx_volume;
  buckets_[ALLOWED_RX] += rx_volume;
  MLOG(MDEBUG) << "receive_credit result:"
               << "total allowed octets " << buckets_[ALLOWED_TOTAL]
               << "total_tx allowed " << buckets_[ALLOWED_TX]
               << "total_rx allowed " << buckets_[ALLOWED_RX];
  // transfer reporting usage to reported
  MLOG(MDEBUG) << "receive_credit:"
               << "reported rx " << buckets_[REPORTED_RX] << "reported_tx "
               << buckets_[REPORTED_TX] << "reporting_rx "
               << buckets_[REPORTING_RX] << "reporting_tx "
               << buckets_[REPORTING_TX];
  buckets_[REPORTED_RX] += buckets_[REPORTING_RX];
  buckets_[REPORTED_TX] += buckets_[REPORTING_TX];
  auto reported_sum = buckets_[REPORTED_RX] + buckets_[REPORTED_TX];
  usage_reporting_limit_ = buckets_[ALLOWED_TOTAL] > reported_sum ?
    buckets_[ALLOWED_TOTAL] - reported_sum : 0;
  MLOG(MDEBUG) << "receive_credit:"
               << "reported rx " << buckets_[REPORTED_RX] << "reported_tx "
               << buckets_[REPORTED_TX] << "reporting_rx "
               << buckets_[REPORTING_RX] << "reporting_tx "
               << buckets_[REPORTING_TX];
  set_expiry_time(validity_time);
  reset_reporting_credit();
  MLOG(MDEBUG) << "receive_credit:"
               << "reported rx " << buckets_[REPORTED_RX] << "reported_tx "
               << buckets_[REPORTED_TX] << "reporting_rx "
               << buckets_[REPORTING_RX] << "reporting_tx "
               << buckets_[REPORTING_TX];
  is_final_ = is_final;
  final_action_info_ = final_action_info;

  if (reauth_state_ == REAUTH_PROCESSING) {
    reauth_state_ = REAUTH_NOT_NEEDED; // done
  }
  if (!quota_exhausted() && (service_state_ == SERVICE_DISABLED ||
                             service_state_ == SERVICE_NEEDS_DEACTIVATION)) {
    // if quota no longer exhausted, reenable services as needed
    MLOG(MDEBUG) << "Quota available. Activating service";
    service_state_ = SERVICE_NEEDS_ACTIVATION;
  }
}

bool SessionCredit::quota_exhausted(
  float usage_reporting_threshold, uint64_t extra_quota_margin)
{
  // used quota since last report
  uint64_t total_reported_usage =
    buckets_[REPORTED_TX] + buckets_[REPORTED_RX];
  uint64_t total_usage_since_report = buckets_[USED_TX] + buckets_[USED_RX];
  if (total_usage_since_report > total_reported_usage) {
    total_usage_since_report -= total_reported_usage;
  } else {
    total_usage_since_report = 0;
  }
  uint64_t tx_usage_since_report = buckets_[USED_TX] > buckets_[REPORTED_TX] ?
    buckets_[USED_TX] - buckets_[REPORTED_TX] : 0;
  uint64_t rx_usage_since_report = buckets_[USED_RX] > buckets_[REPORTED_RX] ?
    buckets_[USED_RX] - buckets_[REPORTED_RX] : 0;

  // available quota since last report
  auto total_usage_reporting_threshold = extra_quota_margin +
    (buckets_[ALLOWED_TOTAL] > total_reported_usage ?
      (buckets_[ALLOWED_TOTAL] -
       total_reported_usage) * usage_reporting_threshold : 0);

  // reported tx/rx could be greater than allowed tx/rx
  // because some OCS/PCRF might not track tx/rx,
  // and 0 is added to the allowed credit when an credit update is received
  auto tx_usage_reporting_threshold =
    buckets_[ALLOWED_TX] > buckets_[REPORTED_TX] ?
      (buckets_[ALLOWED_TX] -
       buckets_[REPORTED_TX]) * usage_reporting_threshold : 0;
  auto rx_usage_reporting_threshold =
    buckets_[ALLOWED_RX] > buckets_[REPORTED_RX] ?
      (buckets_[ALLOWED_RX] -
       buckets_[REPORTED_RX]) * usage_reporting_threshold : 0;

  tx_usage_reporting_threshold += extra_quota_margin;
  rx_usage_reporting_threshold += extra_quota_margin;

  MLOG(MDEBUG) << " Is Quota exhausted?"
               << "\n Total used: " << buckets_[USED_TX] + buckets_[USED_RX]
               << "\n Allowed total: " << buckets_[ALLOWED_TOTAL]
               << "\n Reported total: " << total_reported_usage;

  bool is_exhausted = false;
  is_exhausted = total_usage_since_report >= total_usage_reporting_threshold ||
    (buckets_[ALLOWED_TX] > 0) && (tx_usage_since_report >= tx_usage_reporting_threshold) ||
    (buckets_[ALLOWED_RX] > 0) && (rx_usage_since_report >= rx_usage_reporting_threshold);
  if (is_exhausted == true) {
    MLOG(MDEBUG) << " YES Quota exhausted ";
  }
  return is_exhausted;
}

bool SessionCredit::should_deactivate_service()
{
  return credit_type_ == CreditType::CHARGING &&
         !unlimited_quota_ &&
         SessionCredit::TERMINATE_SERVICE_WHEN_QUOTA_EXHAUSTED &&
         ((no_more_grant() && quota_exhausted()) ||
           quota_exhausted(1, SessionCredit::EXTRA_QUOTA_MARGIN));
}

bool SessionCredit::validity_timer_expired()
{
  return time(NULL) >= expiry_time_;
}

CreditUpdateType SessionCredit::get_update_type()
{
  if (is_reporting()) {
    return CREDIT_NO_UPDATE;
  } else if (is_reauth_required()) {
    return CREDIT_REAUTH_REQUIRED;
  } else if (is_final_ && quota_exhausted()) {
    // Don't request updates if there's no quota left
    return CREDIT_NO_UPDATE;
  } else if (quota_exhausted(SessionCredit::USAGE_REPORTING_THRESHOLD, 0)) {
    return CREDIT_QUOTA_EXHAUSTED;
  } else if (validity_timer_expired()) {
    return CREDIT_VALIDITY_TIMER_EXPIRED;
  } else {
    return CREDIT_NO_UPDATE;
  }
}

SessionCredit::Usage SessionCredit::get_usage_for_reporting(bool is_termination)
{
  // Send delta. If bytes are reporting, don't resend them
  auto report = buckets_[REPORTED_TX] + buckets_[REPORTING_TX];
  uint64_t tx = buckets_[USED_TX] > report ? buckets_[USED_TX] - report : 0;
  report = buckets_[REPORTED_RX] - buckets_[REPORTING_RX];
  uint64_t rx = buckets_[USED_RX] > report ? buckets_[USED_RX] - report : 0;

  if (!is_termination && !is_final_) {
    // Apply reporting limits since the user is not getting terminated.
    // The limits are applied on total usage (ie. tx + rx)
    tx = std::min(tx, usage_reporting_limit_);
    rx = std::min(rx, usage_reporting_limit_ - tx);
  }

  if (get_update_type() == CREDIT_REAUTH_REQUIRED) {
    reauth_state_ = REAUTH_PROCESSING;
  }
  MLOG(MDEBUG) << "get Usage for reporting:"
               << " Used TX:  " << tx << " Used Rx: " << rx
               << "Reporting Tx: " << buckets_[REPORTING_TX]
               << "Reporting Tx: " << buckets_[REPORTING_RX];

  buckets_[REPORTING_TX] += tx;
  buckets_[REPORTING_RX] += rx;
  reporting_ = true;

  MLOG(MDEBUG) << "get Usage for reporting:"
               << " Used TX:  " << tx << " Used Rx: " << rx
               << "Reporting Tx: " << buckets_[REPORTING_TX]
               << "Reporting Tx: " << buckets_[REPORTING_RX];

  return SessionCredit::Usage {.bytes_tx = tx, .bytes_rx = rx};
}

ServiceActionType SessionCredit::get_action()
{
  if (service_state_ == SERVICE_NEEDS_DEACTIVATION) {
    MLOG(MDEBUG) << "Service State: " << service_state_;
    service_state_ = SERVICE_DISABLED;
    return get_action_for_deactivating_service();
  } else if (service_state_ == SERVICE_NEEDS_ACTIVATION) {
    MLOG(MDEBUG) << "Service State: " << service_state_;
    service_state_ = SERVICE_ENABLED;
    return ACTIVATE_SERVICE;
  }
  return CONTINUE_SERVICE;
}

ServiceActionType SessionCredit::get_action_for_deactivating_service()
{
  if (no_more_grant() &&
    final_action_info_.final_action == ChargingCredit_FinalAction_REDIRECT) {
    return REDIRECT;
  } else if (no_more_grant() &&
    final_action_info_.final_action == ChargingCredit_FinalAction_RESTRICT_ACCESS) {
    return RESTRICT_ACCESS;
  } else {
    return TERMINATE_SERVICE;
  }
}

bool SessionCredit::is_reporting()
{
  return reporting_;
}

uint64_t SessionCredit::get_credit(Bucket bucket) const
{
  return buckets_[bucket];
}

bool SessionCredit::is_reauth_required()
{
  return reauth_state_ == REAUTH_REQUIRED;
}

void SessionCredit::reauth()
{
  reauth_state_ = REAUTH_REQUIRED;
}

bool SessionCredit::no_more_grant()
{
  return is_final_;
}

RedirectServer SessionCredit::get_redirect_server() {
  return final_action_info_.redirect_server;
}

} // namespace magma
