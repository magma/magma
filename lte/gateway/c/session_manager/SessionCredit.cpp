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
#include <stdlib.h>
#include <cstdint>
#include <ostream>
#include <string>
#include <unordered_map>
#include <utility>

#include "DiameterCodes.h"
#include "EnumToString.h"
#include "SessionCredit.h"
#include "Utilities.h"
#include "magma_logging.h"
#include "magma_logging_init.h"

namespace magma {

float SessionCredit::USAGE_REPORTING_THRESHOLD = 0.8;
bool SessionCredit::TERMINATE_SERVICE_WHEN_QUOTA_EXHAUSTED = true;
uint64_t SessionCredit::DEFAULT_REQUESTED_UNITS = 10000000;

// by default, enable service & finite credit
SessionCredit::SessionCredit() : SessionCredit(SERVICE_ENABLED, FINITE) {}

SessionCredit::SessionCredit(ServiceState start_state)
    : SessionCredit(start_state, FINITE) {}

SessionCredit::SessionCredit(ServiceState start_state,
                             CreditLimitType credit_limit_type)
    : buckets_{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
      reporting_(false),
      credit_limit_type_(credit_limit_type),
      grant_tracking_type_(TRACKING_UNSET),
      report_last_credit_(false),
      time_of_first_usage_(0),
      time_of_last_usage_(0) {}

SessionCredit::SessionCredit(const StoredSessionCredit& marshaled) {
  reporting_ = marshaled.reporting;
  credit_limit_type_ = marshaled.credit_limit_type;
  grant_tracking_type_ = marshaled.grant_tracking_type;
  received_granted_units_ = marshaled.received_granted_units;
  report_last_credit_ = marshaled.report_last_credit;
  time_of_first_usage_ = marshaled.time_of_first_usage;
  time_of_last_usage_ = marshaled.time_of_last_usage;

  for (int bucket_int = USED_TX; bucket_int != BUCKET_ENUM_MAX_VALUE;
       bucket_int++) {
    Bucket bucket = static_cast<Bucket>(bucket_int);
    if (marshaled.buckets.find(bucket) != marshaled.buckets.end()) {
      buckets_[bucket] = marshaled.buckets.find(bucket)->second;
    }
  }
}

StoredSessionCredit SessionCredit::marshal() const {
  StoredSessionCredit marshaled{};
  marshaled.reporting = reporting_;
  marshaled.credit_limit_type = credit_limit_type_;
  marshaled.grant_tracking_type = grant_tracking_type_;
  marshaled.received_granted_units = received_granted_units_;
  marshaled.report_last_credit = report_last_credit_;
  marshaled.time_of_first_usage = time_of_first_usage_;
  marshaled.time_of_last_usage = time_of_last_usage_;

  for (int bucket_int = USED_TX; bucket_int != BUCKET_ENUM_MAX_VALUE;
       bucket_int++) {
    Bucket bucket = static_cast<Bucket>(bucket_int);
    marshaled.buckets[bucket] = buckets_[bucket];
  }
  return marshaled;
}

SessionCreditUpdateCriteria SessionCredit::get_update_criteria() const {
  SessionCreditUpdateCriteria credit_uc{};
  credit_uc.deleted = false;
  credit_uc.grant_tracking_type = grant_tracking_type_;
  credit_uc.received_granted_units = received_granted_units_;
  credit_uc.report_last_credit = report_last_credit_;
  credit_uc.time_of_first_usage = time_of_first_usage_;
  credit_uc.time_of_last_usage = time_of_last_usage_;

  for (int bucket_int = USED_TX; bucket_int != BUCKET_ENUM_MAX_VALUE;
       bucket_int++) {
    Bucket bucket = static_cast<Bucket>(bucket_int);
    credit_uc.bucket_deltas[bucket] = 0;
  }
  return credit_uc;
}

void SessionCredit::add_used_credit(uint64_t used_tx, uint64_t used_rx,
                                    SessionCreditUpdateCriteria* credit_uc) {
  if (used_tx > 0 || used_rx > 0) {
    buckets_[USED_TX] += used_tx;
    buckets_[USED_RX] += used_rx;
    if (credit_uc) {
      credit_uc->bucket_deltas[USED_TX] += used_tx;
      credit_uc->bucket_deltas[USED_RX] += used_rx;
    }
    update_usage_timestamps(credit_uc);
  }

  log_quota_and_usage();
}

void SessionCredit::update_usage_timestamps(
    SessionCreditUpdateCriteria* credit_uc) {
  auto now = magma::get_time_in_sec_since_epoch();
  if (time_of_first_usage_ == 0) {
    time_of_first_usage_ = now;
  }
  time_of_last_usage_ = now;

  if (credit_uc) {
    credit_uc->time_of_first_usage = time_of_first_usage_;
    credit_uc->time_of_last_usage = time_of_last_usage_;
  }
}

void SessionCredit::reset_reporting_credit(
    SessionCreditUpdateCriteria* credit_uc) {
  buckets_[REPORTING_RX] = 0;
  buckets_[REPORTING_TX] = 0;
  reporting_ = false;
  if (credit_uc) {
    credit_uc->reporting = false;
  }
}

void SessionCredit::mark_failure(uint32_t code,
                                 SessionCreditUpdateCriteria* credit_uc) {
  if (DiameterCodeHandler::is_transient_failure(code)) {
    MLOG(MDEBUG) << "Found transient failure code in mark_failure. Resetting "
                    "'REPORTING' values";
    buckets_[REPORTED_RX] += buckets_[REPORTING_RX];
    buckets_[REPORTED_TX] += buckets_[REPORTING_TX];
    if (credit_uc) {
      credit_uc->bucket_deltas[REPORTED_RX] += buckets_[REPORTING_RX];
      credit_uc->bucket_deltas[REPORTED_TX] += buckets_[REPORTING_TX];
    }
  }
  reset_reporting_credit(credit_uc);
}

// receive_credit will add received grant to current credits. Note that if
// there is over-usage, the extra amount will be added to the counters by
// calculate_delta_allowed_floor and calculate_delta_allowed
void SessionCredit::receive_credit(const GrantedUnits& gsu,
                                   SessionCreditUpdateCriteria* credit_uc) {
  // only decide the grant tracking type on the very first refill
  if (grant_tracking_type_ == TRACKING_UNSET) {
    grant_tracking_type_ = determine_grant_tracking_type(gsu);
  }
  received_granted_units_ = gsu;
  uint64_t bucket_used_total = buckets_[USED_TX] + buckets_[USED_RX];

  // Floor represent the previous value of ALLOWED counters before grant is
  // applied They will only be updated if a valid grant is received
  // In case we have overused (used > allowed) we will reset allowed_floor to
  // used
  uint64_t delta_allowed_floor_total = calculate_delta_allowed_floor(
      gsu.total(), ALLOWED_TOTAL, ALLOWED_FLOOR_TOTAL, bucket_used_total);
  uint64_t delta_allowed_floor_tx = calculate_delta_allowed_floor(
      gsu.total(), ALLOWED_TX, ALLOWED_FLOOR_TX, buckets_[USED_TX]);
  uint64_t delta_allowed_floor_rx = calculate_delta_allowed_floor(
      gsu.total(), ALLOWED_RX, ALLOWED_FLOOR_RX, buckets_[USED_RX]);

  buckets_[ALLOWED_FLOOR_TOTAL] += delta_allowed_floor_total;
  buckets_[ALLOWED_FLOOR_TX] += delta_allowed_floor_tx;
  buckets_[ALLOWED_FLOOR_RX] += delta_allowed_floor_rx;

  // Clear invalid values
  uint64_t total_volume = gsu.total().is_valid() ? gsu.total().volume() : 0;
  uint64_t tx_volume = gsu.tx().is_valid() ? gsu.tx().volume() : 0;
  uint64_t rx_volume = gsu.rx().is_valid() ? gsu.rx().volume() : 0;

  // in case we have overused (used > allowed) we will reset allowed to used +
  // gsu_volume
  uint64_t delta_allowed_total =
      calculate_delta_allowed(total_volume, ALLOWED_TOTAL, bucket_used_total);
  uint64_t delta_allowed_tx =
      calculate_delta_allowed(tx_volume, ALLOWED_TX, buckets_[USED_TX]);
  uint64_t delta_allowed_rx =
      calculate_delta_allowed(rx_volume, ALLOWED_RX, buckets_[USED_RX]);

  // Update allowed bytes
  buckets_[ALLOWED_TOTAL] += delta_allowed_total;
  buckets_[ALLOWED_TX] += delta_allowed_tx;
  buckets_[ALLOWED_RX] += delta_allowed_rx;

  MLOG(MINFO) << "Received the following credit"
              << " total_volume=" << total_volume << " tx_volume=" << tx_volume
              << " rx_volume=" << rx_volume << " grant_tracking_type="
              << grant_type_to_str(grant_tracking_type_);

  // Update reporting/reported bytes
  buckets_[REPORTED_RX] += buckets_[REPORTING_RX];
  buckets_[REPORTED_TX] += buckets_[REPORTING_TX];

  if (credit_uc) {
    credit_uc->received_granted_units = gsu;
    credit_uc->grant_tracking_type = grant_tracking_type_;

    credit_uc->bucket_deltas[ALLOWED_TOTAL] += delta_allowed_total;
    credit_uc->bucket_deltas[ALLOWED_TX] += delta_allowed_tx;
    credit_uc->bucket_deltas[ALLOWED_RX] += delta_allowed_rx;

    credit_uc->bucket_deltas[REPORTED_RX] += buckets_[REPORTING_RX];
    credit_uc->bucket_deltas[REPORTED_TX] += buckets_[REPORTING_TX];

    credit_uc->bucket_deltas[ALLOWED_FLOOR_TOTAL] += delta_allowed_floor_total;
    credit_uc->bucket_deltas[ALLOWED_FLOOR_TX] += delta_allowed_floor_tx;
    credit_uc->bucket_deltas[ALLOWED_FLOOR_RX] += delta_allowed_floor_rx;
  }

  reset_reporting_credit(credit_uc);
  log_quota_and_usage();
}

uint64_t SessionCredit::calculate_delta_allowed_floor(
    CreditUnit cu, Bucket allowed, Bucket floor, uint64_t volume_used) const {
  if (cu.is_valid() && cu.volume() != 0) {
    // only advance floor if there is new grant
    if (buckets_[allowed] < buckets_[floor]) {
      MLOG(MERROR) << "Error in calculate_delta_allowed_floor, "
                      "floor bigger than allowed "
                   << buckets_[floor] << ">" << buckets_[allowed];
      return 0;
    }

    if (volume_used > buckets_[allowed]) {
      // if we overused and received a grant that means that credit was already
      // counted. So add that to the grant
      // allowed_floor = used
      return volume_used - buckets_[floor];
    } else {
      // if we haven't overused, advance the pointer to the old allowed value
      // allowed_floor = allowed
      return buckets_[allowed] - buckets_[floor];
    }
  } else {
    // do NOT advance allowed_floor when grant do not exist or is 0
    return 0;
  }
}

uint64_t SessionCredit::calculate_delta_allowed(uint64_t gsu_volume,
                                                Bucket allowed,
                                                uint64_t volume_used) const {
  if (volume_used > buckets_[allowed]) {
    // if we overused and received a grant that means that credit was already
    // counted.
    // allowed = used + gsu_volume
    return gsu_volume + (volume_used - buckets_[allowed]);
  } else {
    // if we haven't overused, then just add the value
    // allowed = allowed + gsu_volume
    return gsu_volume;
  }
}

bool SessionCredit::is_quota_exhausted(float threshold) const {
  if (credit_limit_type_ != FINITE) {
    return false;
  }

  bool rx_exhausted =
      compute_quota_exhausted(buckets_[ALLOWED_RX], buckets_[USED_RX],
                              threshold, buckets_[ALLOWED_FLOOR_RX]);
  bool tx_exhausted =
      compute_quota_exhausted(buckets_[ALLOWED_TX], buckets_[USED_TX],
                              threshold, buckets_[ALLOWED_FLOOR_TX]);
  bool total_exhausted = compute_quota_exhausted(
      buckets_[ALLOWED_TOTAL], buckets_[USED_TX] + buckets_[USED_RX], threshold,
      buckets_[ALLOWED_FLOOR_TOTAL]);

  bool is_exhausted = false;
  switch (grant_tracking_type_) {
    case TRACKING_UNSET:
      // in case we haven't even initialized the credit at all but we have
      // received traffic, then the session should be marked as exhausted
      is_exhausted = true;
      break;
    case ALL_TOTAL_TX_RX:
      is_exhausted = rx_exhausted || tx_exhausted || total_exhausted;
      break;
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
      is_exhausted = total_exhausted;
      break;
    default:
      MLOG(MERROR) << "Invalid grant_tracking_type="
                   << grant_type_to_str(grant_tracking_type_);
      return false;
  }
  if (is_exhausted) {
    if (threshold == 1) {
      MLOG(MDEBUG) << grant_type_to_str(grant_tracking_type_)
                   << " grant is totally exhausted";
    } else {
      MLOG(MDEBUG) << grant_type_to_str(grant_tracking_type_)
                   << " grant is partially exhausted (threshold " << threshold
                   << ")";
    }
  }
  return is_exhausted;
}

Usage SessionCredit::get_all_unreported_usage_for_reporting(
    SessionCreditUpdateCriteria* credit_uc) {
  auto usage = get_unreported_usage();
  buckets_[REPORTING_TX] += usage.bytes_tx;
  buckets_[REPORTING_RX] += usage.bytes_rx;
  reporting_ = true;
  if (credit_uc) {
    credit_uc->reporting = true;
  }
  log_usage_report(usage);
  return usage;
}

SessionCredit::Summary SessionCredit::get_credit_summary() const {
  return SessionCredit::Summary{
      .usage =
          Usage{
              .bytes_tx = buckets_[USED_TX],
              .bytes_rx = buckets_[USED_RX],
          },
      .time_of_first_usage = time_of_first_usage_,
      .time_of_last_usage = time_of_last_usage_,
  };
}

Usage SessionCredit::get_usage_for_reporting(
    SessionCreditUpdateCriteria* credit_uc) {
  auto usage = get_unreported_usage();
  // Apply reporting limits since the user is not getting terminated.
  // We never want to report more than the amount we've received

  buckets_[REPORTING_TX] += usage.bytes_tx;
  buckets_[REPORTING_RX] += usage.bytes_rx;
  reporting_ = true;
  if (credit_uc) {
    credit_uc->reporting = true;
  }

  log_usage_report(usage);
  return usage;
}

RequestedUnits SessionCredit::get_initial_requested_credits_units() {
  RequestedUnits requestedUnits;
  requestedUnits.set_total(SessionCredit::DEFAULT_REQUESTED_UNITS);
  requestedUnits.set_rx(SessionCredit::DEFAULT_REQUESTED_UNITS);
  requestedUnits.set_tx(SessionCredit::DEFAULT_REQUESTED_UNITS);
  return requestedUnits;
}

RequestedUnits SessionCredit::get_requested_credits_units() const {
  RequestedUnits requestedUnits;
  uint64_t buckets_used_total = buckets_[USED_TX] + buckets_[USED_RX];

  uint64_t total_requested =
      calculate_requested_unit(received_granted_units_.total(), ALLOWED_TOTAL,
                               ALLOWED_FLOOR_TOTAL, buckets_used_total);

  uint64_t tx_requested =
      calculate_requested_unit(received_granted_units_.tx(), ALLOWED_TX,
                               ALLOWED_FLOOR_TX, buckets_[USED_TX]);

  uint64_t rx_requested =
      calculate_requested_unit(received_granted_units_.rx(), ALLOWED_RX,
                               ALLOWED_FLOOR_RX, buckets_[USED_RX]);

  requestedUnits.set_total(total_requested);
  requestedUnits.set_tx(tx_requested);
  requestedUnits.set_rx(rx_requested);

  return requestedUnits;
}

// returns either the last grant, or the difference between the last grant
// and the credit remaining. Prevents over requesting in case we still have
// credit available from the previous request
uint64_t SessionCredit::calculate_requested_unit(CreditUnit cu, Bucket allowed,
                                                 Bucket allowed_floor,
                                                 uint64_t used) const {
  if (cu.is_valid() == false) {
    return 0;
  }
  // get the current volume grant, or infer it in case of 0
  int64_t grant = cu.volume() != 0
                      ? cu.volume()
                      : buckets_[allowed] - buckets_[allowed_floor];
  int64_t remaining = buckets_[allowed] - used;
  if (remaining >= 0 && grant >= remaining) {
    // request just partial of a grant since we still have some credit left
    return grant - remaining;
  }
  return grant;
}

Usage SessionCredit::get_unreported_usage() const {
  Usage usage = {0, 0};
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

bool SessionCredit::compute_quota_exhausted(const uint64_t allowed,
                                            const uint64_t used,
                                            float threshold_ratio,
                                            const uint64_t floor) const {
  // current_granted_units: difference between current allowed and past
  // allowed (floor). This value is equivalent to the grant received.
  // Credit will be considered exhausted if the remaining credit is below
  // a percentage of the current granted units

  if (floor > allowed) {
    // This should never happen because floor should always be the previous
    // allowed value
    MLOG(MERROR) << "Error: Floor value is higher than allowed value " << floor
                 << ">" << allowed;
    return true;
  }

  if (used >= allowed) {
    // if we used more that we are allowed, we are for sure out of quota
    return true;
  }

  uint64_t remaining_credit = allowed - used;
  uint64_t current_granted_units = allowed - floor;
  // this step is necessary to avoid precision issues of float
  int integer_threshold_ratio = 100 - static_cast<int>(threshold_ratio * 100);
  uint64_t threshold = (current_granted_units * integer_threshold_ratio) / 100;

  return remaining_credit <= threshold;
}

bool SessionCredit::is_reporting() const { return reporting_; }

uint64_t SessionCredit::get_credit(Bucket bucket) const {
  return buckets_[bucket];
}

void SessionCredit::set_report_last_credit(
    bool report_last_credit, SessionCreditUpdateCriteria* credit_uc) {
  report_last_credit_ = report_last_credit;
  if (credit_uc) {
    credit_uc->report_last_credit = report_last_credit;
  }
}

void SessionCredit::set_reporting(bool reporting) { reporting_ = reporting; }

bool SessionCredit::is_report_last_credit() const {
  return report_last_credit_;
}

void SessionCredit::apply_update_criteria(
    const SessionCreditUpdateCriteria& credit_uc) {
  grant_tracking_type_ = credit_uc.grant_tracking_type;
  received_granted_units_ = credit_uc.received_granted_units;
  report_last_credit_ = credit_uc.report_last_credit;
  time_of_first_usage_ = credit_uc.time_of_first_usage;
  time_of_last_usage_ = credit_uc.time_of_last_usage;
  // DO NOT UPDATE reporting_. (done by LocalSessionManagerHandler)

  // add credit
  for (int i = USED_TX; i != BUCKET_ENUM_MAX_VALUE; i++) {
    Bucket bucket = static_cast<Bucket>(i);
    auto credit = credit_uc.bucket_deltas.find(bucket)->second;
    buckets_[bucket] += credit;
  }
}

// Determine the grant's tracking type by looking at which values are valid.
GrantTrackingType SessionCredit::determine_grant_tracking_type(
    const GrantedUnits& grant) const {
  bool total_valid = grant.total().is_valid() && grant.total().volume() != 0;
  bool tx_valid = grant.tx().is_valid() && grant.tx().volume() != 0;
  bool rx_valid = grant.rx().is_valid() && grant.rx().volume() != 0;

  if (total_valid && tx_valid && rx_valid) {
    return ALL_TOTAL_TX_RX;
  }
  if (total_valid) {
    return TOTAL_ONLY;
  }
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

bool SessionCredit::current_grant_contains_zero() const {
  switch (grant_tracking_type_) {
    case ALL_TOTAL_TX_RX:
      // Monitors should not have this mode enabled
      MLOG(MWARNING) << "Possible monitor with ALL_TOTAL_TX_RX enabled";
      return is_received_grented_unit_zero(received_granted_units_.total()) ||
             is_received_grented_unit_zero(received_granted_units_.tx()) ||
             is_received_grented_unit_zero(received_granted_units_.rx());
      break;
    case RX_ONLY:
      return is_received_grented_unit_zero(received_granted_units_.rx());
      break;
    case TX_ONLY:
      return is_received_grented_unit_zero(received_granted_units_.tx());
      break;
    case TX_AND_RX:
      return is_received_grented_unit_zero(received_granted_units_.tx()) ||
             is_received_grented_unit_zero(received_granted_units_.rx());
      break;
    case TOTAL_ONLY:
      return is_received_grented_unit_zero(received_granted_units_.total());
      break;
    default:
      MLOG(MERROR) << "Can't determine if grant has zeroes. "
                      "tracking type is probably unexpected: "
                   << grant_tracking_type_;
      return true;
      break;
  }
}

bool SessionCredit::is_received_grented_unit_zero(const CreditUnit& cu) const {
  if (!cu.is_valid() || cu.volume() == 0) {
    return true;
  }
  return false;
}

void SessionCredit::log_quota_and_usage() const {
  if (magma::get_verbosity() != MDEBUG) {
    return;
  }
  MLOG(MDEBUG) << "===> Used     Tx: " << buckets_[USED_TX]
               << " Rx: " << buckets_[USED_RX]
               << " Total: " << buckets_[USED_TX] + buckets_[USED_RX];
  MLOG(MDEBUG) << "===> Reported Tx: " << buckets_[REPORTED_TX]
               << " Rx: " << buckets_[REPORTED_RX]
               << " Total: " << buckets_[REPORTED_RX] + buckets_[REPORTED_TX];
  MLOG(MDEBUG) << "===> Allowed  Tx: " << buckets_[ALLOWED_TX]
               << " Rx: " << buckets_[ALLOWED_RX]
               << " Total: " << buckets_[ALLOWED_TOTAL];
  MLOG(MDEBUG) << "===> A_Floor  Tx: " << buckets_[ALLOWED_FLOOR_TX]
               << " Rx: " << buckets_[ALLOWED_FLOOR_RX]
               << " Total: " << buckets_[ALLOWED_FLOOR_TOTAL];
  MLOG(MDEBUG) << "===> (%used)  Tx: "
               << get_percentage_usage(buckets_[ALLOWED_TX],
                                       buckets_[ALLOWED_FLOOR_TX],
                                       buckets_[USED_TX])
               << " Rx: "
               << get_percentage_usage(buckets_[ALLOWED_RX],
                                       buckets_[ALLOWED_FLOOR_RX],
                                       buckets_[USED_RX])
               << " Total: "
               << get_percentage_usage(buckets_[ALLOWED_TOTAL],
                                       buckets_[ALLOWED_FLOOR_TOTAL],
                                       buckets_[USED_TX] + buckets_[USED_RX]);

  MLOG(MDEBUG) << "===> Grant tracking type "
               << grant_type_to_str(grant_tracking_type_)
               << ",  Reporting: " << reporting_;
  MLOG(MDEBUG) << "===> Last Granted Units Received (tx/rx/total) "
               << received_granted_units_.tx().volume() << "/"
               << received_granted_units_.rx().volume() << "/"
               << received_granted_units_.total().volume();
}

std::string SessionCredit::get_percentage_usage(uint64_t allowed,
                                                uint64_t floor,
                                                uint64_t used) const {
  if (allowed <= floor) {
    return "_%";
  }
  int64_t currentGrant = allowed - floor;
  int64_t currentUsage = used - floor;
  int currentPercent = static_cast<int>(100 * currentUsage / currentGrant);
  // cap % in case it grows too much
  if (abs(currentPercent) >= 1000) {
    currentPercent = 999 * currentPercent / abs(currentPercent);
  }
  return std::to_string(currentPercent) + "%";
}

void SessionCredit::log_usage_report(Usage usage) const {
  if (magma::get_verbosity() != MDEBUG) {
    return;
  }
  MLOG(MDEBUG) << "===> Amount reporting for this report:"
               << " tx=" << usage.bytes_tx << " rx=" << usage.bytes_rx;
  MLOG(MDEBUG) << "===> The total amount currently being reported:"
               << " tx=" << buckets_[REPORTING_TX]
               << " rx=" << buckets_[REPORTING_RX];
}
}  // namespace magma
