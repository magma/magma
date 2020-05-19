/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include <ctime>
#include <memory>
#include <unordered_map>
#include <unordered_set>

#include <lte/protos/session_manager.grpc.pb.h>

#include "CreditKey.h"
#include "ServiceAction.h"
#include "StoredState.h"

namespace magma {

enum CreditUpdateType {
  CREDIT_NO_UPDATE = 0,
  CREDIT_QUOTA_EXHAUSTED = 1,
  CREDIT_VALIDITY_TIMER_EXPIRED = 2,
  CREDIT_REAUTH_REQUIRED = 3
};

enum CreditType {
  MONITORING = 0,
  CHARGING = 1,
};

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

  static std::unique_ptr<SessionCredit>
  unmarshal(const StoredSessionCredit &marshaled, CreditType credit_type);

  StoredSessionCredit marshal();

  SessionCreditUpdateCriteria get_update_criteria();

  SessionCredit(CreditType credit_type);

  SessionCredit(CreditType credit_type, ServiceState start_state);

  SessionCredit(CreditType credit_type, ServiceState start_state,
                bool unlimited_quota);

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
  void reset_reporting_credit(SessionCreditUpdateCriteria &update_criteria);

  /**
   * Credit update has failed to the OCS, so mark this credit as failed so it
   * can be cut off accordingly
   */
  void mark_failure(uint32_t code,
                    SessionCreditUpdateCriteria &update_criteria);
  /**
   * receive_credit increments ALLOWED* and moves the REPORTING_* credit to
   * the REPORTED_* credit
   */
  void receive_credit(uint64_t total_volume, uint64_t tx_volume,
                      uint64_t rx_volume, uint32_t validity_time,
                      bool is_final_grant, FinalActionInfo final_action_info,
                      SessionCreditUpdateCriteria &update_criteria);

  /**
   * get_update_type returns the type of update required for the credit. If no
   * update is required, it returns CREDIT_NO_UPDATE. Else, it returns an update
   * type
   */
  CreditUpdateType get_update_type() const;

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
   * get_action returns the action to take on the credit based on the last
   * update. If no action needs to take place, CONTINUE_SERVICE is returned.
   */
  ServiceActionType get_action(SessionCreditUpdateCriteria &update_criteria);

  /**
   * Returns true if either of REPORTING_* buckets are more than 0
   */
  bool is_reporting() const;

  /**
   * Helper function to get the credit in a particular bucket
   */
  uint64_t get_credit(Bucket bucket) const;

  /**
   * Mark the credit to be in the REAUTH_REQUIRED state. The next time
   * get_update is called, this credit will report its usage.
   */
  void reauth(SessionCreditUpdateCriteria &update_criteria);

  /**
   * Returns
   */
  RedirectServer get_redirect_server() const;

  /**
   * Mark SessionCredit as having been given the final grant.
   * NOTE: Use only for merging updates into SessionStore
   * @param is_final_grant
   */
  void set_is_final_grant(bool is_final_grant,
                          SessionCreditUpdateCriteria &update_criteria);

  /**
   * Set ReAuthState.
   * NOTE: Use only for merging updates into SessionStore
   * @param reauth_state
   */
  void set_reauth(ReAuthState reauth_state,
                  SessionCreditUpdateCriteria &update_criteria);

  /**
   * Set ServiceState.
   * NOTE: Use only for merging updates into SessionStore
   * @param service_state
   */
  void set_service_state(ServiceState new_service_state,
                         SessionCreditUpdateCriteria &update_criteria);

  /**
   * Set expiry time of SessionCredit
   * NOTE: Use only for merging updates into SessionStore
   * @param expiry_time
   */
  void set_expiry_time(std::time_t expiry_time,
                       SessionCreditUpdateCriteria &update_criteria);

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
   * A threshold represented as a ratio for triggering usage update before
   * an user completely used up the quota
   * Session manager will send usage update when
   * (available bytes since last update) * USAGE_REPORTING_THRESHOLD >=
   * (used bytes since last update)
   */
  static float USAGE_REPORTING_THRESHOLD;

  /**
   * Extra number of bytes an user could use after the quota is exhausted.
   * Session manager will deactivate the service when
   * used quota >= (granted quota + EXTRA_QUOTA_MARGIN)
   */
  static uint64_t EXTRA_QUOTA_MARGIN;

  /**
   * Set to true to terminate service when the quota of a session is exhausted.
   * An user can still use up to the extra margin.
   * Set to false to allow users to use without any constraint.
   */
  static bool TERMINATE_SERVICE_WHEN_QUOTA_EXHAUSTED;

private:
  bool reporting_;
  bool is_final_grant_;
  bool unlimited_quota_;
  FinalActionInfo final_action_info_;
  ReAuthState reauth_state_;
  ServiceState service_state_;
  std::time_t expiry_time_;
  uint64_t buckets_[MAX_VALUES];
  /**
   * Limit for the total usage (tx + rx) in credit updates to prevent
   * session manager from reporting more usage than granted
   */
  uint64_t usage_reporting_limit_;
  CreditType credit_type_;

private:
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
   * We will also add a extra_quota_margin on to the available quota if it is
   * specified.
   * Check if the session has exhausted its quota granted since the last report.
   *
   * @param usage_reporting_threshold
   * @param extra_quota_margin Extra bytes of usage allowable before quota
   *        is exhausted.
   * @return true if quota is exhausted for the session
   */
  bool is_quota_exhausted(float usage_reporting_threshold = 1,
                          uint64_t extra_quota_margin = 0) const;

  void log_quota_and_usage() const;

  bool should_deactivate_service() const;

  bool validity_timer_expired() const;

  void set_expiry_time(uint32_t validity_time,
                       SessionCreditUpdateCriteria &update_criteria);

  bool is_reauth_required() const;

  ServiceActionType get_action_for_deactivating_service() const;

  SessionCredit::Usage get_unreported_usage() const;

  void log_usage_report(SessionCredit::Usage) const;
};

} // namespace magma
