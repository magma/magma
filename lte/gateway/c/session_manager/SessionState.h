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
#pragma once

#include <functional>
#include <experimental/optional>
#include <utility>

#include <lte/protos/session_manager.grpc.pb.h>
#include <lte/protos/spgw_service.grpc.pb.h>

#include "RuleStore.h"
#include "SessionReporter.h"
#include "StoredState.h"
#include "CreditKey.h"
#include "SessionCredit.h"
#include "Monitor.h"
#include "ChargingGrant.h"

namespace magma {
typedef std::unordered_map<
    CreditKey, std::unique_ptr<ChargingGrant>, decltype(&ccHash),
    decltype(&ccEqual)>
    CreditMap;
typedef std::unordered_map<std::string, std::unique_ptr<Monitor>> MonitorMap;
static SessionStateUpdateCriteria UNUSED_UPDATE_CRITERIA;

struct RulesToProcess {
  std::vector<std::string> static_rules;
  std::vector<PolicyRule> dynamic_rules;
};

// Used to transform the proto message RuleSet into a more useful structure
struct RuleSetToApply {
  std::unordered_map<std::string, PolicyRule> dynamic_rules;
  std::unordered_set<std::string> static_rules;

  RuleSetToApply() {}
  RuleSetToApply(const magma::lte::RuleSet& rule_set);
  void combine_rule_set(const RuleSetToApply& other);
};

// Used to transform the proto message RulesPerSubscriber into a more useful
// structure
struct RuleSetBySubscriber {
  std::string imsi;
  std::unordered_map<std::string, RuleSetToApply> rule_set_by_apn;
  std::experimental::optional<RuleSetToApply> subscriber_wide_rule_set;

  RuleSetBySubscriber(const RulesPerSubscriber& rules_per_subscriber);
  std::experimental::optional<RuleSetToApply> get_combined_rule_set_for_apn(
      const std::string& apn);
};

struct BearerUpdate {
  bool needs_creation;
  CreateBearerRequest create_req;  // only valid if needs_creation is true
  bool needs_deletion;
  DeleteBearerRequest delete_req;  // only valid if needs_deletion is true
  BearerUpdate() : needs_creation(false), needs_deletion(false) {}
};

/**
 * SessionState keeps track of a current UE session in the PCEF, recording
 * usage and allowance for all charging keys
 */
class SessionState {
 public:
  struct SessionInfo {
    std::string imsi;
    std::string ip_addr;
    std::vector<std::string> static_rules;
    std::vector<PolicyRule> dynamic_rules;
    std::vector<PolicyRule> gy_dynamic_rules;
    std::vector<std::string> restrict_rules;
    std::experimental::optional<AggregatedMaximumBitrate> ambr;
  };
  struct TotalCreditUsage {
    uint64_t monitoring_tx;
    uint64_t monitoring_rx;
    uint64_t charging_tx;
    uint64_t charging_rx;
  };

 public:
  SessionState(
      const std::string& imsi, const std::string& session_id,
      const SessionConfig& cfg, StaticRuleStore& rule_store,
      const magma::lte::TgppContext& tgpp_context, uint64_t pdp_start_time);

  SessionState(
      const StoredSessionState& marshaled, StaticRuleStore& rule_store);

  static std::unique_ptr<SessionState> unmarshal(
      const StoredSessionState& marshaled, StaticRuleStore& rule_store);

  StoredSessionState marshal();

  /**
   * Updates rules to be scheduled, active, or removed, depending on the
   * specified time.
   *
   * NOTE: This function has undefined behavior if attempting to go backwards
   *       in time.
   */
  void sync_rules_to_time(
      std::time_t current_time, SessionStateUpdateCriteria& update_criteria);

  /**
   * add_rule_usage adds used TX/RX bytes to a particular rule
   */
  void add_rule_usage( const std::string& rule_id, uint64_t used_tx,
    uint64_t used_rx, SessionStateUpdateCriteria& update_criteria);

  /**
   * get_updates collects updates and adds them to a UpdateSessionRequest
   * for reporting.
   * Only updates request number
   * @param update_request (out) - request to add new updates to
   * @param actions (out) - actions to take on services
   */
  void get_updates(
      UpdateSessionRequest& update_request_out,
      std::vector<std::unique_ptr<ServiceAction>>* actions_out,
      SessionStateUpdateCriteria& update_criteria);

  /**
   * mark_as_awaiting_termination transitions the session state from
   * SESSION_ACTIVE to SESSION_TERMINATION_SCHEDULED
   */
  void mark_as_awaiting_termination(
      SessionStateUpdateCriteria& update_criteria);

  bool is_terminating();

  /**
   * complete_termination collects final usages for all credits into a
   * SessionTerminateRequest and calls the on termination callback with the
   * request.
   * Note that complete_termination will forcefully complete the termination
   * no matter the current state of the session. To properly complete the
   * termination, this function should only be called when
   * can_complete_termination returns true.
   */
  void complete_termination(
      SessionReporter& reporter, SessionStateUpdateCriteria& update_criteria);

  bool reset_reporting_charging_credit(const CreditKey &key,
                              SessionStateUpdateCriteria &update_criteria);

  /**
   * Receive the credit grant if the credit update was successful
   *
   * @param update
   * @param uc
   * @return True if usage for the charging key is allowed after receiving
   *         the credit response. This requires that the credit update was
   *         a success. Also it must either be an infinite credit rating group,
   *         or have associated credit grant.
   */
  bool receive_charging_credit(const CreditUpdateResponse &update,
                               SessionStateUpdateCriteria &update_criteria);

  uint64_t get_charging_credit(const CreditKey &key, Bucket bucket) const;

  ReAuthResult reauth_key(const CreditKey &charging_key,
                               SessionStateUpdateCriteria &update_criteria);

  ReAuthResult reauth_all(SessionStateUpdateCriteria &update_criteria);

  void set_charging_credit(const CreditKey &key,
    ChargingGrant charging_grant, SessionStateUpdateCriteria &uc);

  /**
   * get_total_credit_usage returns the tx and rx of the session,
   * accounting for all unique keys (charging and monitoring) used by all
   * rules (static and dynamic)
   * Should be called after complete_termination.
   */
  TotalCreditUsage get_total_credit_usage();

  std::string get_session_id() const;

  SubscriberQuotaUpdate_Type get_subscriber_quota_state() const;

  bool is_radius_cwf_session() const;

  void get_session_info(SessionState::SessionInfo& info);

  void set_tgpp_context(
      const magma::lte::TgppContext& tgpp_context,
      SessionStateUpdateCriteria& update_criteria);

  void set_config(const SessionConfig& config);

  SessionConfig get_config() const;

  void set_subscriber_quota_state(
      const magma::lte::SubscriberQuotaUpdate_Type state,
      SessionStateUpdateCriteria& update_criteria);

  bool active_monitored_rules_exist();

  uint32_t get_request_number();

  uint64_t get_pdp_start_time();

  void set_pdp_end_time(uint64_t epoch);

  uint64_t get_pdp_end_time();

  void increment_request_number(uint32_t incr);

  // Methods related to the session's static and dynamic rules
  /**
   * Infer the policy's type (STATIC or DYNAMIC)
   * @param rule_id
   * @return the type if the rule exists, {} otherwise.
   */
  std::experimental::optional<PolicyType> get_policy_type(
      const std::string& rule_id);

  bool is_dynamic_rule_installed(const std::string& rule_id);

  bool is_gy_dynamic_rule_installed(const std::string& rule_id);

  bool is_static_rule_installed(const std::string& rule_id);

  bool is_restrict_rule_installed(const std::string& rule_id);

  /**
   * Add a dynamic rule to the session which is currently active.
   */
  void insert_dynamic_rule(
      const PolicyRule& rule, RuleLifetime& lifetime,
      SessionStateUpdateCriteria& update_criteria);

  void insert_gy_dynamic_rule(
      const PolicyRule& rule, RuleLifetime& lifetime,
      SessionStateUpdateCriteria& update_criteria);

  /**
   * Add a static rule to the session which is currently active.
   */
  void activate_static_rule(
      const std::string& rule_id, RuleLifetime& lifetime,
      SessionStateUpdateCriteria& update_criteria);

  void activate_restrict_rule(
      const std::string& rule_id, RuleLifetime& lifetime,
      SessionStateUpdateCriteria& update_criteria);

  /**
   * Remove a currently active dynamic rule to mark it as deactivated.
   *
   * @param rule_id ID of the rule to be removed.
   * @param rule_out Will point to the removed rule if it is not nullptr
   * @param update_criteria Tracks updates to the session. To be passed back to
   *                        the SessionStore to resolve issues of concurrent
   *                        updates to a session.
   * @return True if successfully removed.
   */
  bool remove_dynamic_rule(
      const std::string& rule_id, PolicyRule* rule_out,
      SessionStateUpdateCriteria& update_criteria);

  bool remove_scheduled_dynamic_rule(
      const std::string& rule_id, PolicyRule* rule_out,
      SessionStateUpdateCriteria& update_criteria);

  bool remove_gy_dynamic_rule(
      const std::string& rule_id, PolicyRule *rule_out,
      SessionStateUpdateCriteria& update_criteria);

  /**
   * Remove a currently active static rule to mark it as deactivated.
   *
   * @param rule_id ID of the rule to be removed.
   * @param update_criteria Tracks updates to the session. To be passed back to
   *                        the SessionStore to resolve issues of concurrent
   *                        updates to a session.
   * @return True if successfully removed.
   */
  bool deactivate_static_rule(
      const std::string& rule_id, SessionStateUpdateCriteria& update_criteria);


  bool deactivate_scheduled_static_rule(
      const std::string& rule_id, SessionStateUpdateCriteria& update_criteria);

  bool deactivate_restrict_rule(
      const std::string& rule_id, SessionStateUpdateCriteria& update_criteria);


  std::vector<std::string>& get_static_rules();

  std::set<std::string>& get_scheduled_static_rules();

  std::vector<std::string>& get_restrict_rules();

  DynamicRuleStore& get_dynamic_rules();

  DynamicRuleStore& get_scheduled_dynamic_rules();

  /**
   * Schedule a dynamic rule for activation in the future.
   */
  void schedule_dynamic_rule(
      const PolicyRule& rule, RuleLifetime& lifetime,
      SessionStateUpdateCriteria& update_criteria);

  /**
   * Schedule a static rule for activation in the future.
   */
  void schedule_static_rule(
      const std::string& rule_id, RuleLifetime& lifetime,
      SessionStateUpdateCriteria& update_criteria);

  /**
   * Mark a scheduled dynamic rule as activated.
   */
  void install_scheduled_dynamic_rule(
      const std::string& rule_id, SessionStateUpdateCriteria& update_criteria);

  /**
   * Mark a scheduled static rule as activated.
   */
  void install_scheduled_static_rule(
      const std::string& rule_id, SessionStateUpdateCriteria& update_criteria);

  RuleLifetime& get_rule_lifetime(const std::string& rule_id);

  DynamicRuleStore& get_gy_dynamic_rules();

  uint32_t total_monitored_rules_count();

  bool is_active();

  uint32_t get_credit_key_count();

  void set_fsm_state(
    SessionFsmState new_state, SessionStateUpdateCriteria& uc);

  StaticRuleInstall get_static_rule_install(
    const std::string& rule_id, const RuleLifetime& lifetime);

  DynamicRuleInstall get_dynamic_rule_install(
    const std::string& rule_id, const RuleLifetime& lifetime);

  SessionFsmState get_state();

  // Event Triggers
  void add_new_event_trigger(
  magma::lte::EventTrigger trigger,
  SessionStateUpdateCriteria& update_criteria);

  void mark_event_trigger_as_triggered(
    magma::lte::EventTrigger trigger,
    SessionStateUpdateCriteria& update_criteria);

  void set_event_trigger(
    magma::lte::EventTrigger trigger, const EventTriggerState value,
    SessionStateUpdateCriteria& update_criteria);

  void remove_event_trigger(magma::lte::EventTrigger trigger,
    SessionStateUpdateCriteria& update_criteria);

  void set_revalidation_time(const google::protobuf::Timestamp& time,
                             SessionStateUpdateCriteria& update_criteria);

  google::protobuf::Timestamp get_revalidation_time() {return revalidation_time_;}

  EventTriggerStatus get_event_triggers() {return pending_event_triggers_;}

  bool is_credit_in_final_unit_state(const CreditKey &charging_key) const;

  // Monitors
  bool receive_monitor(const UsageMonitoringUpdateResponse &update,
                       SessionStateUpdateCriteria &uc);

  uint64_t get_monitor(const std::string &key, Bucket bucket) const;

  bool add_to_monitor(const std::string &key, uint64_t used_tx,
                      uint64_t used_rx, SessionStateUpdateCriteria &uc);

  void set_monitor(
      const std::string& key, Monitor monitor, SessionStateUpdateCriteria& uc);

  bool reset_reporting_monitor(
      const std::string& key, SessionStateUpdateCriteria& uc);

  void set_session_level_key(const std::string new_key);

  bool apply_update_criteria(SessionStateUpdateCriteria& uc);

  // QoS Management
  /**
   * get_dedicated_bearer_updates processes the two rule update inputs and
   * produces a BearerUpdate based on whether a bearer has to be create/deleted.
   * @param rules_to_activate
   * @param rules_to_deactivate
   * @param uc: update criteria needs to be updated if the bearer mapping is
   * modified
   * @return BearerUpdate with Create/DeleteBearerRequest
   */
  BearerUpdate get_dedicated_bearer_updates(
      RulesToProcess& rules_to_activate, RulesToProcess& rules_to_deactivate,
      SessionStateUpdateCriteria& uc);

  /**
   *
   * @param rule_set
   * @param subscriber_wide_rule_set
   * @param rules_to_activate
   * @param rules_to_deactivate
   */
  void apply_session_rule_set(
      RuleSetToApply& rule_set, RulesToProcess& rules_to_activate,
      RulesToProcess& rules_to_deactivate, SessionStateUpdateCriteria& uc);

  /**
   * Add the association of policy -> bearerID into bearer_id_by_policy_
   * This assumes the bearerID is not 0
   */
  void bind_policy_to_bearer(
      const PolicyBearerBindingRequest& request,
      SessionStateUpdateCriteria& uc);

 private:
  std::string imsi_;
  std::string session_id_;
  uint32_t request_number_;
  SessionFsmState curr_state_;
  SessionConfig config_;
  uint64_t pdp_start_time_;
  uint64_t pdp_end_time_;
  // Used to keep track of whether the subscriber has valid quota.
  // (only used for CWF at the moment)
  magma::lte::SubscriberQuotaUpdate_Type subscriber_quota_state_;
  magma::lte::TgppContext tgpp_context_;

  // All static rules synced from policy DB
  StaticRuleStore& static_rules_;
  // Static rules that are currently installed for the session
  std::vector<std::string> active_static_rules_;
  // Dynamic GX rules that are currently installed for the session
  DynamicRuleStore dynamic_rules_;
  // Dynamic GY rules that are currently installed for the session
  DynamicRuleStore gy_dynamic_rules_;
  // Static rules that are currently installed when service is restricted
  std::vector<std::string> active_restrict_rules_;

  // Static rules that are scheduled for installation for the session
  std::set<std::string> scheduled_static_rules_;
  // Dynamic rules that are scheduled for installation for the session
  DynamicRuleStore scheduled_dynamic_rules_;
  // Activation & deactivation times for each rule that is either currently
  // installed, or scheduled for installation for this session
  std::unordered_map<std::string, RuleLifetime> rule_lifetimes_;

  // map of Gx event_triggers that are pending and its status (bool)
  // If the value is true, that means an update request for that event trigger
  // should be sent
  EventTriggerStatus pending_event_triggers_;
  // todo for stateless we will have to store a bit more information so we can
  // reschedule triggers
  google::protobuf::Timestamp revalidation_time_;
  /*
   * ChargingCreditPool manages a pool of credits for OCS-based charging. It is
   * keyed by rating groups & service Identity (uint32, [uint32]) and receives
   * CreditUpdateResponses to update
   * credit
   */
  CreditMap credit_map_;
  MonitorMap monitor_map_;
  std::string session_level_key_;

  // PolicyID->DedicatedBearerID used for 4G bearer/QoS management
  BearerIDByPolicyID bearer_id_by_policy_;

 private:
  /**
   * For this session, add the CreditUsageUpdate to the UpdateSessionRequest.
   * Also
   *
   * @param update_request_out Modified with added CreditUsageUpdate
   * @param actions_out Modified with additional actions to take on session
   */
  void get_charging_updates(
      UpdateSessionRequest& update_request_out,
      std::vector<std::unique_ptr<ServiceAction>>* actions_out,
      SessionStateUpdateCriteria& uc);

  void apply_charging_credit_update(
      const CreditKey &key, SessionCreditUpdateCriteria &credit_update);

  /**
   * Receive the credit grant if the credit update was successful
   *
   * @param update
   * @param uc
   * @return True if usage for the charging key is allowed after receiving
   *         the credit response. This requires that the credit update was
   *         a success. Also it must either be an infinite credit rating group,
   *         or have associated credit grant.
   */
  bool init_charging_credit(
    const CreditUpdateResponse &update, SessionStateUpdateCriteria &uc);

  /**
   * Return true if any credit unit is valid and has non-zero volume
   */
  bool contains_credit(const GrantedUnits& gsu);

  bool is_infinite_credit(const CreditUpdateResponse& response);

  /**
   * For this session, add the UsageMonitoringUpdateRequest to the
   * UpdateSessionRequest.
   *
   * @param update_request_out Modified with added UsdageMonitoringUpdateRequest
   * @param actions_out Modified with additional actions to take on session.
   */
  void get_monitor_updates(
      UpdateSessionRequest& update_request_out,
      std::vector<std::unique_ptr<ServiceAction>>* actions_out,
      SessionStateUpdateCriteria& update_criteria);

  /**
   * Apply SessionCreditUpdateCriteria, a per-credit diff of an update, into
   * the SessionState object
   * @param key : monitoring key for the update
   * @param update : the diff that needs to be applied
   */
  void apply_monitor_updates(
      const std::string &key, SessionCreditUpdateCriteria &update);

  void add_common_fields_to_usage_monitor_update(UsageMonitoringUpdateRequest* req);

  SessionTerminateRequest make_termination_request(
    SessionStateUpdateCriteria& update_criteria);

  /**
   * Returns true if the specified rule should be active at that time
   */
  bool should_rule_be_active(const std::string& rule_id, std::time_t time);

  /**
   * Returns true if the specified rule should be deactivated by that time
   */
  bool should_rule_be_deactivated(const std::string& rule_id, std::time_t time);

  SessionCreditUpdateCriteria* get_credit_uc(
    const CreditKey &key, SessionStateUpdateCriteria &uc);

  CreditUsageUpdate make_credit_usage_update_req(CreditUsage& usage) const;

  bool init_new_monitor(
    const UsageMonitoringUpdateResponse &update,
    SessionStateUpdateCriteria &update_criteria);

  void update_session_level_key(
    const UsageMonitoringUpdateResponse &update,
    SessionStateUpdateCriteria &update_criteria);

  SessionCreditUpdateCriteria* get_monitor_uc(
    const std::string &key, SessionStateUpdateCriteria &uc);

  void fill_protos_tgpp_context(magma::lte::TgppContext* tgpp_context) const;

  void get_event_trigger_updates(
      UpdateSessionRequest& update_request_out,
      std::vector<std::unique_ptr<ServiceAction>>* actions_out,
      SessionStateUpdateCriteria& update_criteria);

  bool is_static_rule_scheduled(const std::string& rule_id);

  bool is_dynamic_rule_scheduled(const std::string& rule_id);

  /** apply static_rules which is the desired state for the session's rules **/
  void apply_session_static_rule_set(
      std::unordered_set<std::string> static_rules,
      RulesToProcess& rules_to_activate, RulesToProcess& rules_to_deactivate,
      SessionStateUpdateCriteria& uc);

  /** apply dynamic_rules which is the desired state for the session's rules **/
  void apply_session_dynamic_rule_set(
      std::unordered_map<std::string, PolicyRule> dynamic_rules,
      RulesToProcess& rules_to_activate, RulesToProcess& rules_to_deactivate,
      SessionStateUpdateCriteria& uc);

  /**
   * Check if a new bearer has to be created for the given policy. If a creation
   * is needed, fill in the BearerUpdate with required info.
   * @param policy_type
   * @param rule_id
   * @param update
   */
  void update_bearer_creation_req(
      const PolicyType policy_type, const std::string& rule_id,
      const SessionConfig& config, BearerUpdate& update);

  /**
   * Check if a bearer has to be deleted for the given policy. If a deletion is
   * needed, fill in the BearerUpdate with required info.
   * @param policy_type
   * @param rule_id
   * @param update
   */
  void update_bearer_deletion_req(
      const PolicyType policy_type, const std::string& rule_id,
      const SessionConfig& config, BearerUpdate& update,
      SessionStateUpdateCriteria& uc);

  /**
   * Set bearer_id_by_policy_ to the input
   * @param bearer_id_by_policy
   */
  void set_bearer_map(BearerIDByPolicyID bearer_id_by_policy) {
    bearer_id_by_policy_ = bearer_id_by_policy;
  }
  /**
   * @param policy_type
   * @param rule_id
   * @return true if the policy definition includes a QoS field
   */
  bool policy_has_qos(
      const PolicyType policy_type, const std::string& rule_id,
      PolicyRule* rule_out);

  /**
   * Increments data usage values for session
   * @param bytes_tx
   * @param bytes_rx
   */
  void update_data_usage_metrics(uint64_t bytes_tx, uint64_t bytes_rx);
};

}  // namespace magma
