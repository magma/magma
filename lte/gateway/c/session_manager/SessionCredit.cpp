/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <limits>

#include "SessionCredit.h"
#include "magma_logging.h"

namespace magma {

uint64_t SessionCredit::USAGE_REPORTING_LIMIT =
    std::numeric_limits<uint64_t>::max();

SessionCredit::SessionCredit(ServiceState start_state)
  : reporting_(false),
    reauth_state_(REAUTH_NOT_NEEDED),
    service_state_(start_state),
    buckets_{} {}

// by default, enable service
SessionCredit::SessionCredit()
  : SessionCredit(SERVICE_ENABLED) {}

void SessionCredit::set_expiry_time(uint32_t validity_time) {
  if (validity_time == 0) {
    // set as max possible time
    expiry_time_ = std::numeric_limits<std::time_t>::max();
    return;
  }
  expiry_time_ = std::time(nullptr) + validity_time;
}

void SessionCredit::add_used_credit(uint64_t used_tx, uint64_t used_rx) {
  buckets_[USED_TX] += used_tx;
  buckets_[USED_RX] += used_rx;
  if (quota_exhausted() && is_final_) {
    MLOG(MDEBUG) << "Quota exhausted. Deactivating service";
    service_state_ = SERVICE_NEEDS_DEACTIVATION;
  }
}

void SessionCredit::reset_reporting_credit() {
  buckets_[REPORTING_RX] = 0;
  buckets_[REPORTING_TX] = 0;
  reporting_ = false;
}

void SessionCredit::mark_failure() {
  reset_reporting_credit();
  service_state_ = SERVICE_NEEDS_DEACTIVATION;
}

void SessionCredit::receive_credit(
    uint64_t total_volume,
    uint64_t tx_volume,
    uint64_t rx_volume,
    uint32_t validity_time,
    bool is_final) {
  MLOG(MDEBUG) << "receive_credit:" << "total allowed octets:  "
       << buckets_[ALLOWED_TOTAL] << "total_tx allowed: "
       << buckets_[ALLOWED_TX] << "total_rx allowed "
       << buckets_[ALLOWED_RX];
  buckets_[ALLOWED_TOTAL] += total_volume;
  buckets_[ALLOWED_TX] += tx_volume;
  buckets_[ALLOWED_RX] += rx_volume;
  MLOG(MDEBUG) << "receive_credit:" << "total allowed octets "
       << buckets_[ALLOWED_TOTAL] << "total_tx allowed "
       << buckets_[ALLOWED_TX] << "total_rx allowed "
       << buckets_[ALLOWED_RX];
  // transfer reporting usage to reported
  MLOG(MDEBUG) << "receive_credit:" << "reported rx "
       << buckets_[REPORTED_RX] << "reported_tx "
       << buckets_[REPORTED_TX] << "reporting_rx "
       << buckets_[REPORTING_RX] << "reporting_tx "
       << buckets_[REPORTING_TX];
  buckets_[REPORTED_RX] += buckets_[REPORTING_RX];
  buckets_[REPORTED_TX] += buckets_[REPORTING_TX];
  MLOG(MDEBUG) << "receive_credit:" << "reported rx "
       << buckets_[REPORTED_RX] << "reported_tx "
       << buckets_[REPORTED_TX] << "reporting_rx "
       << buckets_[REPORTING_RX] << "reporting_tx "
       << buckets_[REPORTING_TX];
  set_expiry_time(validity_time);
  reset_reporting_credit();
  MLOG(MDEBUG) << "receive_credit:" << "reported rx "
       << buckets_[REPORTED_RX] << "reported_tx "
       << buckets_[REPORTED_TX] << "reporting_rx "
       << buckets_[REPORTING_RX] << "reporting_tx "
       << buckets_[REPORTING_TX];
  is_final_ = is_final;

  if (reauth_state_ == REAUTH_PROCESSING) {
    reauth_state_ = REAUTH_NOT_NEEDED; // done
  }
  if (!quota_exhausted() && (service_state_ == SERVICE_DISABLED
      || service_state_ == SERVICE_NEEDS_DEACTIVATION)) {
    // if quota no longer exhausted, reenable services as needed
    MLOG(MDEBUG) << "Quota available. Activating service";
    service_state_ = SERVICE_NEEDS_ACTIVATION;
  }
}

bool SessionCredit::quota_exhausted() {
  uint64_t used_tot = buckets_[USED_TX] + buckets_[USED_RX];
  bool is_exhausted = false;
  if ((buckets_[ALLOWED_TX] > 0) || (buckets_[ALLOWED_RX] > 0)) {
    is_exhausted = used_tot > buckets_[ALLOWED_TOTAL]
      || buckets_[USED_TX] > buckets_[ALLOWED_TX]
      || buckets_[USED_RX] > buckets_[ALLOWED_RX];
  } else {
    MLOG(MDEBUG) << " Is Quota exhausted ? " <<  "Total_used:  "
      << used_tot << " Allowed Total: " << buckets_[ALLOWED_TOTAL];
    is_exhausted =  used_tot > buckets_[ALLOWED_TOTAL];
  }
  if (is_exhausted == true) {
    MLOG(MDEBUG) << " YES Quota exhausted " <<  "Total_used:  "
      << used_tot << " Allowed Total: " << buckets_[ALLOWED_TOTAL];
  }
  return is_exhausted;
}

bool SessionCredit::validity_timer_expired() {
  return time(NULL) >= expiry_time_;
}

CreditUpdateType SessionCredit::get_update_type() {
  if (is_reporting()) {
    return CREDIT_NO_UPDATE;
  } else if (is_reauth_required()) {
    return CREDIT_REAUTH_REQUIRED;
  } else if (is_final_ && quota_exhausted()) {
    // Don't request updates if there's no quota left
    return CREDIT_NO_UPDATE;
  } else if (quota_exhausted()) {
    return CREDIT_QUOTA_EXHAUSTED;
  } else if (validity_timer_expired()){
    return CREDIT_VALIDITY_TIMER_EXPIRED;
  } else {
    return CREDIT_NO_UPDATE;
  }
}

SessionCredit::Usage SessionCredit::get_usage_for_reporting(
    bool is_termination) {
  // Send delta. If bytes are reporting, don't resend them
  uint64_t tx =
    buckets_[USED_TX] - buckets_[REPORTED_TX] - buckets_[REPORTING_TX];
  uint64_t rx =
    buckets_[USED_RX] - buckets_[REPORTED_RX] - buckets_[REPORTING_RX];

  if (!is_termination && !is_final_) {
    // Apply reporting limits since the user is not getting terminated.
    // The limits are applied on total usage (ie. tx + rx)
    uint64_t limit = SessionCredit::USAGE_REPORTING_LIMIT;
    tx = std::min(tx, limit);
    rx = std::min(rx, limit - tx);
  }

  if (get_update_type() == CREDIT_REAUTH_REQUIRED) {
    reauth_state_ = REAUTH_PROCESSING;
  }
  MLOG(MDEBUG) << "get Usage for reporting:" << " Used TX:  "
       << tx << " Used Rx: " << rx << "Reporting Tx: "
       << buckets_[REPORTING_TX] <<  "Reporting Tx: "
       << buckets_[REPORTING_RX];

  buckets_[REPORTING_TX] += tx;
  buckets_[REPORTING_RX] += rx;
  reporting_ = true;

  MLOG(MDEBUG) << "get Usage for reporting:" << " Used TX:  "
       << tx << " Used Rx: " << rx << "Reporting Tx: "
       << buckets_[REPORTING_TX] <<  "Reporting Tx: "
       << buckets_[REPORTING_RX];

  return SessionCredit::Usage{
    .bytes_tx = tx,
    .bytes_rx = rx};

}

ServiceActionType SessionCredit::get_action() {
  if (service_state_ == SERVICE_NEEDS_DEACTIVATION) {
    MLOG(MDEBUG) << "Service State: " << service_state_ ;
    // received used credits, but service should be disabled
    service_state_ = SERVICE_DISABLED;
    return TERMINATE_SERVICE;
  } else if (service_state_ == SERVICE_NEEDS_ACTIVATION) {
    MLOG(MDEBUG) << "Service State: " << service_state_ ;
    // didn't receive used credits, but service should be enabled
    service_state_ = SERVICE_ENABLED;
    return ACTIVATE_SERVICE;
  }
  return CONTINUE_SERVICE;
}

bool SessionCredit::is_reporting() {
  return reporting_;
}

uint64_t SessionCredit::get_credit(Bucket bucket) const {
  return buckets_[bucket];
}

bool SessionCredit::is_reauth_required() {
  return reauth_state_ == REAUTH_REQUIRED;
}

void SessionCredit::reauth() {
  reauth_state_ = REAUTH_REQUIRED;
}

}
