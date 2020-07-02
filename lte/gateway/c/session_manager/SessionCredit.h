/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include <lte/protos/session_manager.grpc.pb.h>

#include "StoredState.h"

namespace magma {
/**
 * SessionCredit tracks all the credit volumes associated with a charging key
 * for a user. It can receive used credit, add allowed credit, and check if
 * there is an update (quota exhausted, etc)
 */
class SessionCredit {
public:
  struct Usage {
    uint64_t bytes_tx;
    uint64_t bytes_rx;
  };

  static SessionCredit unmarshal(const StoredSessionCredit &marshaled);

  StoredSessionCredit marshal();

  SessionCreditUpdateCriteria get_update_criteria();

  SessionCredit();

  SessionCredit(ServiceState start_state);

  SessionCredit(ServiceState start_state,
                CreditLimitType credit_limit_type);

  /**
   * add_used_credit increments USED_TX and USED_RX
   * as being recently updated
   */
  void add_used_credit(uint64_t used_tx, uint64_t used_rx,
                       SessionCreditUpdateCriteria &update_criteria);

  /**
   * reset_reporting_credit resets the REPORTING_* to 0
   * Also marks the session as not in reporting.
   */
  void reset_reporting_credit(SessionCreditUpdateCriteria* uc);

  /**
   * Credit update has failed to the OCS, so mark this credit as failed so it
   * can be cut off accordingly
   */
  void mark_failure(uint32_t code, SessionCreditUpdateCriteria* uc);
  /**
   * receive_credit increments ALLOWED* and moves the REPORTING_* credit to
   * the REPORTED_* credit
   */
  void receive_credit(const GrantedUnits& gsu, SessionCreditUpdateCriteria* uc);

  /**
   * get_update returns a filled-in CreditUsage if an update exists, and a blank
   * one if no update exists. Check has_update before calling.
   * This method also sets the REPORTING_* credit buckets
   */
  SessionCredit::Usage
  get_usage_for_reporting(SessionCreditUpdateCriteria &update_criteria);

  SessionCredit::Usage get_all_unreported_usage_for_reporting(
      SessionCreditUpdateCriteria &update_criteria);

  /**
   * Returns true if either of REPORTING_* buckets are more than 0
   */
  bool is_reporting() const;

  /**
   * Helper function to get the credit in a particular bucket
   */
  uint64_t get_credit(Bucket bucket) const;

  void set_grant_tracking_type(GrantTrackingType g_type,
    SessionCreditUpdateCriteria& uc);

  /**
   * Add credit to the specified bucket. This does not necessarily correspond
   * to allowed or used credit.
   * NOTE: Use only for merging updates into SessionStore
   * @param credit
   * @param bucket
   */
  void add_credit(uint64_t credit, Bucket bucket,
                  SessionCreditUpdateCriteria &update_criteria);
  /**
   * is_quota_exhausted checks if any of the tx, rx, or combined tx+rx are
   * exhausted. The exception to this is if ALLOWED_RX or ALLOWED_TX are 0,
   * which occurs for an OCS/PCRF which do not individually track tx/rx. In this
   * scenario, only the total matters.
   *
   * Quota usage is measured by reporting from pipelined since the last
   * SessionUpdate.
   * We mark quota as exhausted if usage_reporting_threshold * available quota
   * is reached. (so the default is 100% of quota)
   * Check if the session has exhausted its quota granted since the last report.
   *
   * @param usage_reporting_threshold
   * @return true if quota is exhausted for the session
   */
  bool is_quota_exhausted(float usage_reporting_threshold) const;

  /**
   * A threshold represented as a ratio for triggering usage update before
   * an user completely used up the quota
   * Session manager will send usage update when
   * (available bytes since last update) * USAGE_REPORTING_THRESHOLD >=
   * (used bytes since last update)
   */
  static float USAGE_REPORTING_THRESHOLD;

  /**
   * Set to true to terminate service when the quota of a session is exhausted.
   * An user can still use up to the extra margin.
   * Set to false to allow users to use without any constraint.
   */
  static bool TERMINATE_SERVICE_WHEN_QUOTA_EXHAUSTED;
private:
  uint64_t buckets_[MAX_VALUES];
  bool reporting_;
  CreditLimitType credit_limit_type_;
  GrantTrackingType grant_tracking_type_;

private:
  void log_quota_and_usage() const;

  SessionCredit::Usage get_unreported_usage() const;

  void log_usage_report(SessionCredit::Usage) const;

  GrantTrackingType determine_grant_tracking_type(const GrantedUnits& grant);

  bool compute_quota_exhausted(const uint64_t allowed,
    const uint64_t reported, const uint64_t used, float threshold_ratio) const;

  uint64_t compute_reporting_limit(
    const uint64_t allowed, const uint64_t reported) const;

  void apply_reporting_limits(SessionCredit::Usage& usage);
};

} // namespace magma
