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

float SessionCredit::USAGE_REPORTING_THRESHOLD             = 0.8;
bool SessionCredit::TERMINATE_SERVICE_WHEN_QUOTA_EXHAUSTED = true;

SessionCredit SessionCredit::unmarshal(const StoredSessionCredit& marshaled) {
  SessionCredit credit;

  credit.reporting_           = marshaled.reporting;
  credit.credit_limit_type_   = marshaled.credit_limit_type;
  credit.grant_tracking_type_ = marshaled.grant_tracking_type;

  for (int bucket_int = USED_TX; bucket_int != MAX_VALUES; bucket_int++) {
    Bucket bucket = static_cast<Bucket>(bucket_int);
    if (marshaled.buckets.find(bucket) != marshaled.buckets.end()) {
      credit.buckets_[bucket] = marshaled.buckets.find(bucket)->second;
    }
  }
  return credit;
}

StoredSessionCredit SessionCredit::marshal() {
  StoredSessionCredit marshaled{};
  marshaled.reporting           = reporting_;
  marshaled.credit_limit_type   = credit_limit_type_;
  marshaled.grant_tracking_type = grant_tracking_type_;

  for (int bucket_int = USED_TX; bucket_int != MAX_VALUES; bucket_int++) {
    Bucket bucket             = static_cast<Bucket>(bucket_int);
    marshaled.buckets[bucket] = buckets_[bucket];
  }
  return marshaled;
}

SessionCreditUpdateCriteria SessionCredit::get_update_criteria() {
  SessionCreditUpdateCriteria uc{};
  uc.grant_tracking_type = grant_tracking_type_;
  for (int bucket_int = USED_TX; bucket_int != MAX_VALUES; bucket_int++) {
    Bucket bucket            = static_cast<Bucket>(bucket_int);
    uc.bucket_deltas[bucket] = 0;
  }
  return uc;
}

SessionCredit::SessionCredit(ServiceState start_state)
    : SessionCredit(start_state, FINITE) {}

SessionCredit::SessionCredit(
    ServiceState start_state, CreditLimitType credit_limit_type)
    : buckets_{},
      reporting_(false),
      credit_limit_type_(credit_limit_type),
      grant_tracking_type_(TOTAL_ONLY) {}

// by default, enable service & finite credit
SessionCredit::SessionCredit() : SessionCredit(SERVICE_ENABLED, FINITE) {}

void SessionCredit::add_used_credit(
    uint64_t used_tx, uint64_t used_rx, SessionCreditUpdateCriteria& uc) {
  buckets_[USED_TX] += used_tx;
  buckets_[USED_RX] += used_rx;
  uc.bucket_deltas[USED_TX] += used_tx;
  uc.bucket_deltas[USED_RX] += used_rx;

  log_quota_and_usage();
}

void SessionCredit::reset_reporting_credit(SessionCreditUpdateCriteria* uc) {
  buckets_[REPORTING_RX] = 0;
  buckets_[REPORTING_TX] = 0;
  reporting_             = false;
  if (uc != NULL) {
    uc->reporting = false;
  }
}

void SessionCredit::mark_failure(
    uint32_t code, SessionCreditUpdateCriteria* uc) {
  if (DiameterCodeHandler::is_transient_failure(code)) {
    buckets_[REPORTED_RX] += buckets_[REPORTING_RX];
    buckets_[REPORTED_TX] += buckets_[REPORTING_TX];
    if (uc != NULL) {
      uc->bucket_deltas[REPORTED_RX] += buckets_[REPORTING_RX];
      uc->bucket_deltas[REPORTED_TX] += buckets_[REPORTING_TX];
    }
  }
  reset_reporting_credit(uc);
}

void SessionCredit::receive_credit(
    const GrantedUnits& gsu, SessionCreditUpdateCriteria* uc) {
  grant_tracking_type_ = determine_grant_tracking_type(gsu);

  // Clear invalid values
  uint64_t total_volume = gsu.total().is_valid() ? gsu.total().volume() : 0;
  uint64_t tx_volume    = gsu.tx().is_valid() ? gsu.tx().volume() : 0;
  uint64_t rx_volume    = gsu.rx().is_valid() ? gsu.rx().volume() : 0;

  // Update allowed bytes
  buckets_[ALLOWED_TOTAL] += total_volume;
  buckets_[ALLOWED_TX] += tx_volume;
  buckets_[ALLOWED_RX] += rx_volume;

  MLOG(MINFO) << "Received the following credit"
              << " total_volume=" << total_volume << " tx_volume=" << tx_volume
              << " rx_volume=" << rx_volume << " grant_tracking_type="
              << grant_type_to_str(grant_tracking_type_);

  // Update reporting/reported bytes
  buckets_[REPORTED_RX] += buckets_[REPORTING_RX];
  buckets_[REPORTED_TX] += buckets_[REPORTING_TX];

  if (uc != NULL) {
    uc->grant_tracking_type = grant_tracking_type_;

    uc->bucket_deltas[ALLOWED_TOTAL] += total_volume;
    uc->bucket_deltas[ALLOWED_TX] += tx_volume;
    uc->bucket_deltas[ALLOWED_RX] += rx_volume;

    uc->bucket_deltas[REPORTED_RX] += buckets_[REPORTING_RX];
    uc->bucket_deltas[REPORTED_TX] += buckets_[REPORTING_TX];
  }

  reset_reporting_credit(uc);
  log_quota_and_usage();
}

bool SessionCredit::is_quota_exhausted(float threshold) const {
  if (credit_limit_type_ != FINITE) {
    return false;
  }
  uint64_t total_reported = buckets_[REPORTED_TX] + buckets_[REPORTED_RX];
  uint64_t total_usage    = buckets_[USED_TX] + buckets_[USED_RX];

  bool rx_exhausted = compute_quota_exhausted(
      buckets_[ALLOWED_RX], buckets_[REPORTED_RX], buckets_[USED_RX],
      threshold);
  bool tx_exhausted = compute_quota_exhausted(
      buckets_[ALLOWED_TX], buckets_[REPORTED_TX], buckets_[USED_TX],
      threshold);

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

SessionCredit::Usage SessionCredit::get_all_unreported_usage_for_reporting(
    SessionCreditUpdateCriteria& update_criteria) {
  auto usage = get_unreported_usage();
  buckets_[REPORTING_TX] += usage.bytes_tx;
  buckets_[REPORTING_RX] += usage.bytes_rx;
  reporting_                = true;
  update_criteria.reporting = true;
  log_usage_report(usage);
  return usage;
}

SessionCredit::Usage SessionCredit::get_usage_for_reporting(
    SessionCreditUpdateCriteria& update_criteria) {
  auto usage = get_unreported_usage();
  // Apply reporting limits since the user is not getting terminated.
  // We never want to report more than the amount we've received
  apply_reporting_limits(usage);

  buckets_[REPORTING_TX] += usage.bytes_tx;
  buckets_[REPORTING_RX] += usage.bytes_rx;
  reporting_                = true;
  update_criteria.reporting = true;

  log_usage_report(usage);
  return usage;
}

// Take the minimum of (grant - reported) and (reported - used)
void SessionCredit::apply_reporting_limits(SessionCredit::Usage& usage) {
  uint64_t tx_limit, rx_limit, total_limit, total_reported;

  tx_limit =
      compute_reporting_limit(buckets_[ALLOWED_TX], buckets_[REPORTED_TX]);
  rx_limit =
      compute_reporting_limit(buckets_[ALLOWED_TX], buckets_[REPORTED_TX]);

  switch (grant_tracking_type_) {
    case TX_ONLY:
      usage.bytes_tx = std::min(usage.bytes_tx, tx_limit);
      usage.bytes_rx = 0;  // Clear field that doesn't need to be reported
      MLOG(MDEBUG) << "Applying a TX reporting limit of " << tx_limit;
      break;
    case RX_ONLY:
      usage.bytes_rx = std::min(usage.bytes_rx, rx_limit);
      usage.bytes_tx = 0;  // Clear field that doesn't need to be reported
      MLOG(MDEBUG) << "Applying a RX reporting limit of " << rx_limit;
      break;
    case TX_AND_RX:
      usage.bytes_tx = std::min(usage.bytes_tx, tx_limit);
      usage.bytes_rx = std::min(usage.bytes_rx, rx_limit);
      MLOG(MDEBUG) << "Applying TX and RX reporting limits of " << tx_limit
                   << ", " << rx_limit;
      break;
    case TOTAL_ONLY:
      total_reported = buckets_[REPORTED_RX] + buckets_[REPORTED_TX];
      total_limit =
          compute_reporting_limit(buckets_[ALLOWED_TOTAL], total_reported);
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
  SessionCredit::Usage usage = {bytes_tx: 0, bytes_rx: 0};
  auto report                = buckets_[REPORTED_TX] + buckets_[REPORTING_TX];
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

uint64_t SessionCredit::compute_reporting_limit(
    const uint64_t allowed, const uint64_t reported) const {
  uint64_t limit = 0;
  if (allowed > reported) {
    limit = allowed - reported;
  }
  return limit;
}

bool SessionCredit::compute_quota_exhausted(
    const uint64_t allowed, const uint64_t reported, const uint64_t used,
    float threshold_ratio) const {
  uint64_t unreported_usage = 0;
  if (used > reported) {
    unreported_usage = used - reported;
  }
  uint64_t grant_left = allowed - reported;
  auto threshold      = std::max(0.0f, grant_left * threshold_ratio);
  return unreported_usage >= threshold;
}

bool SessionCredit::is_reporting() const {
  return reporting_;
}

uint64_t SessionCredit::get_credit(Bucket bucket) const {
  return buckets_[bucket];
}

void SessionCredit::set_grant_tracking_type(
    GrantTrackingType g_type, SessionCreditUpdateCriteria& uc) {
  grant_tracking_type_   = g_type;
  uc.grant_tracking_type = g_type;
}

void SessionCredit::add_credit(
    uint64_t credit, Bucket bucket,
    SessionCreditUpdateCriteria& update_criteria) {
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
    MLOG(MWARNING)
        << "Received a GSU with no valid grants, keeping the same type";
    return grant_tracking_type_;
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
               << grant_type_to_str(grant_tracking_type_)
               << ",  Reporting: " << reporting_;
}

void SessionCredit::log_usage_report(SessionCredit::Usage usage) const {
  MLOG(MDEBUG) << "===> Amount reporting for this report:"
               << " tx=" << usage.bytes_tx << " rx=" << usage.bytes_rx;
  MLOG(MDEBUG) << "===> The total amount currently being reported:"
               << " tx=" << buckets_[REPORTING_TX]
               << " rx=" << buckets_[REPORTING_RX];
}
}  // namespace magma
