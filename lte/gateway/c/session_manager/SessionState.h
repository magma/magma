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
#include "Types.h"

namespace magma {
using std::experimental::optional;
typedef std::unordered_map<
    CreditKey, std::unique_ptr<ChargingGrant>, decltype(&ccHash),
    decltype(&ccEqual)>
    CreditMap;
typedef std::unordered_map<
    CreditKey, SessionCredit::Summary, decltype(&ccHash), decltype(&ccEqual)>
    ChargingCreditSummaries;
typedef std::unordered_map<std::string, std::unique_ptr<Monitor>> MonitorMap;

// Used to transform the proto message RuleSet into a more useful structure
struct RuleSetToApply {
  std::unordered_map<std::string, PolicyRule> dynamic_rules;
  std::unordered_set<std::string> static_rules;

  RuleSetToApply() {}
  RuleSetToApply(const magma::lte::RuleSet& rule_set);
  void combine_rule_set(const RuleSetToApply& other);
};

struct UpdateRequests {
  std::vector<UsageMonitoringUpdateRequest> monitor_requests;
  std::vector<CreditUsageUpdate> charging_requests;
};

// Used to transform the proto message RulesPerSubscriber into a more useful
// structure
struct RuleSetBySubscriber {
  std::string imsi;
  std::unordered_map<std::string, RuleSetToApply> rule_set_by_apn;
  optional<RuleSetToApply> subscriber_wide_rule_set;

  RuleSetBySubscriber(const RulesPerSubscriber& rules_per_subscriber);
  optional<RuleSetToApply> get_combined_rule_set_for_apn(
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
  // SessionInfo is a struct used to bundle necessary information
  // PipelineDClient needs to make requests
  struct SessionInfo {
    enum upfNodeType {
      IPv4 = 0,
      IPv6 = 1,
      FQDN = 2,
    };

    typedef struct tNodeId {
      upfNodeType node_id_type;
      char node_id[40];
    } NodeId;

    typedef struct Fseid {
      uint64_t f_seid;
      NodeId Nid;
    } FSid;
    std::string imsi;
    std::string ip_addr;
    std::string ipv6_addr;
    Teids teids;

    uint32_t local_f_teid;
    std::string msisdn;
    RulesToProcess gx_rules;
    RulesToProcess gy_dynamic_rules;
    optional<AggregatedMaximumBitrate> ambr;
    // 5G specific extensions
    std::vector<SetGroupPDR> Pdr_rules_;
    magma::lte::Fsm_state_FsmState state;
    std::string subscriber_id;
    uint32_t ver_no;
    NodeId nodeId;
    FSid Seid;
    // 5G specific extension routines
  };

  /* To remove below routine  once the UPF node core logic
   * get completed
   */
  void sess_infocopy(struct SessionInfo*);

  magma::lte::Fsm_state_FsmState get_proto_fsm_state();

 public:
  SessionState(
      const std::string& imsi, const std::string& session_id,
      const SessionConfig& cfg, StaticRuleStore& rule_store,
      const magma::lte::TgppContext& tgpp_context, uint64_t pdp_start_time,
      const CreateSessionResponse& csr);

  SessionState(
      const StoredSessionState& marshaled, StaticRuleStore& rule_store);

  static std::unique_ptr<SessionState> unmarshal(
      const StoredSessionState& marshaled, StaticRuleStore& rule_store);

  StoredSessionState marshal();

  // 5G processing constructor without response contxt as set-interface msg
  SessionState(
      const std::string& imsi, const std::string& session_ctx_id,
      const SessionConfig& cfg, StaticRuleStore& rule_store);

  /* methods of new messages of 5G and handle other message*/
  uint32_t get_current_version();

  void set_current_version(
      int new_session_version, SessionStateUpdateCriteria& uc);

  void insert_pdr(SetGroupPDR* rule);

  void set_remove_all_pdrs();

  void insert_far(SetGroupFAR* rule);

  void remove_all_rules();

  std::vector<SetGroupPDR>& get_all_pdr_rules();

  std::vector<SetGroupFAR>& get_all_far_rules();

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
  void add_rule_usage(
      const std::string& rule_id, uint64_t used_tx, uint64_t used_rx,
      uint64_t dropped_tx, uint64_t dropped_rx,
      SessionStateUpdateCriteria& update_criteria);

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

  bool is_terminating();

  /**
   * complete_termination checks the FSM state and transitions the state to
   * TERMINATED, if it can. If the state is ACTIVE or TERMINATED, it will not do
   * anything.
   * This function will return true if the termination happened successfully.
   */
  bool can_complete_termination(SessionStateUpdateCriteria& update_criteria);

  void handle_update_failure(
      const UpdateRequests& failed_requests,
      SessionStateUpdateCriteria& session_uc);

  /**
   * Receive the credit grant if the credit update was successful
   *
   * @param update
   * @param session_uc
   * @return True if usage for the charging key is allowed after receiving
   *         the credit response. This requires that the credit update was
   *         a success.
   */
  bool receive_charging_credit(
      const CreditUpdateResponse& update,
      SessionStateUpdateCriteria& session_uc);

  uint64_t get_charging_credit(const CreditKey& key, Bucket bucket) const;

  bool set_credit_reporting(
      const CreditKey& key, bool reporting,
      SessionStateUpdateCriteria* update_criteria);

  ReAuthResult reauth_key(
      const CreditKey& charging_key,
      SessionStateUpdateCriteria& update_criteria);

  ReAuthResult reauth_all(SessionStateUpdateCriteria& update_criteria);

  void set_charging_credit(
      const CreditKey& key, ChargingGrant charging_grant,
      SessionStateUpdateCriteria& uc);

  std::vector<PolicyRule> get_all_final_unit_rules();

  /**
   * get_total_credit_usage returns the tx and rx of the session,
   * accounting for all unique keys (charging and monitoring) used by all
   * rules (static and dynamic)
   * Should be called after complete_termination.
   */
  SessionCredit::TotalCreditUsage get_total_credit_usage();

  ChargingCreditSummaries get_charging_credit_summaries();

  std::string get_session_id() const;
  uint32_t get_local_teid() const;
  void set_local_teid(uint32_t teid, SessionStateUpdateCriteria& uc);

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

  void set_pdp_end_time(uint64_t epoch, SessionStateUpdateCriteria& session_uc);

  uint64_t get_pdp_end_time();

  uint64_t get_active_duration_in_seconds();

  void increment_request_number(uint32_t incr);

  SessionTerminateRequest make_termination_request(
      SessionStateUpdateCriteria& uc);

  CreateSessionResponse get_create_session_response();

  void clear_create_session_response();

  // Methods related to the session's static and dynamic rules
  /**
   * Infer the policy's type (STATIC or DYNAMIC)
   * @param rule_id
   * @return the type if the rule exists, {} otherwise.
   */
  optional<PolicyType> get_policy_type(const std::string& rule_id);

  /**
   * @brief Get the current rule version object
   *
   * @param rule_id
   * @return uint32_t
   */
  uint32_t get_current_rule_version(const std::string& rule_id);

  bool is_dynamic_rule_installed(const std::string& rule_id);

  bool is_gy_dynamic_rule_installed(const std::string& rule_id);

  bool is_static_rule_installed(const std::string& rule_id);

  /**
   * @brief Add a dynamic rule into dynamic rule store. Increment the associated
   * version and return the new version.
   *
   * @param rule
   * @param lifetime
   * @param session_uc
   * @return uint32_t updated version
   */
  uint32_t insert_dynamic_rule(
      const PolicyRule& rule, RuleLifetime& lifetime,
      SessionStateUpdateCriteria& session_uc);

  /**
   * @brief Insert a static rule into active_static_rules_. Increment the
   * associated version and return the new version.
   *
   * @param rule_id
   * @param lifetime
   * @param session_uc
   * @return uint32_t updated version
   */
  uint32_t activate_static_rule(
      const std::string& rule_id, RuleLifetime& lifetime,
      SessionStateUpdateCriteria& session_uc);
  /**
   * @brief Insert a PolicyRule into gy_dynamic_rules_
   *
   * @param rule
   * @param lifetime
   * @param update_criteria
   * @return uint32_t updated version
   */
  uint32_t insert_gy_rule(
      const PolicyRule& rule, RuleLifetime& lifetime,
      SessionStateUpdateCriteria& session_uc);

  /**
   * Remove a currently active dynamic rule to mark it as deactivated.
   *
   * @param rule_id ID of the rule to be removed.
   * @param rule_out Will point to the removed rule if it is not nullptr
   * @param update_criteria Tracks updates to the session. To be passed back to
   *                        the SessionStore to resolve issues of concurrent
   *                        updates to a session.
   * @return optional<uint32_t> updated version if success, {} if failure
   */
  optional<uint32_t> remove_dynamic_rule(
      const std::string& rule_id, PolicyRule* rule_out,
      SessionStateUpdateCriteria& update_criteria);

  bool remove_scheduled_dynamic_rule(
      const std::string& rule_id, PolicyRule* rule_out,
      SessionStateUpdateCriteria& update_criteria);

  /**
   * @brief Remove a Gy rule from SessionState and increment the corresponding
   * version
   *
   * @param rule_id
   * @param rule_out
   * @param session_uc
   * @return optional<uint32_t> updated version if success, {} if failure
   */
  optional<uint32_t> remove_gy_rule(
      const std::string& rule_id, PolicyRule* rule_out,
      SessionStateUpdateCriteria& session_uc);

  /**
   * Remove a currently active static rule to mark it as deactivated.
   *
   * @param rule_id ID of the rule to be removed.
   * @param session_uc Tracks updates to the session. To be passed back to
   *                        the SessionStore to resolve issues of concurrent
   *                        updates to a session.
   * @return new version if successfully removed. otherwise returns {}
   */
  optional<uint32_t> deactivate_static_rule(
      const std::string& rule_id, SessionStateUpdateCriteria& session_uc);

  bool deactivate_scheduled_static_rule(const std::string& rule_id);

  std::vector<std::string>& get_static_rules();

  std::set<std::string>& get_scheduled_static_rules();

  DynamicRuleStore& get_dynamic_rules();

  DynamicRuleStore& get_scheduled_dynamic_rules();

  std::vector<PolicyRule> get_all_active_policies();

  /**
   * Schedule a dynamic rule for activation in the future.
   */
  void schedule_dynamic_rule(
      const PolicyRule& rule, RuleLifetime& lifetime,
      SessionStateUpdateCriteria& update_criteria);

  bool is_static_rule_scheduled(const std::string& rule_id);

  /**
   * Schedule a static rule for activation in the future.
   */
  void schedule_static_rule(
      const std::string& rule_id, RuleLifetime& lifetime,
      SessionStateUpdateCriteria& update_criteria);

  void set_suspend_credit(
      const CreditKey& charging_key, bool new_suspended,
      SessionStateUpdateCriteria& update_criteria);

  bool is_credit_suspended(const CreditKey& charging_key);

  RuleLifetime& get_rule_lifetime(const std::string& rule_id);

  DynamicRuleStore& get_gy_dynamic_rules();

  uint32_t total_monitored_rules_count();

  bool is_active();

  uint32_t get_credit_key_count();

  void set_fsm_state(SessionFsmState new_state, SessionStateUpdateCriteria& uc);

  StaticRuleInstall get_static_rule_install(
      const std::string& rule_id, const RuleLifetime& lifetime);

  DynamicRuleInstall get_dynamic_rule_install(
      const std::string& rule_id, const RuleLifetime& lifetime);

  SessionFsmState get_state();

  void get_rules_per_credit_key(
      CreditKey charging_key, RulesToProcess& rulesToProcess,
      SessionStateUpdateCriteria& session_uc);

  /**
   * Remove all active/scheduled static/dynamic rules and reflect the change in
   * session_uc
   * @param session_uc
   */
  void remove_all_rules_for_termination(SessionStateUpdateCriteria& session_uc);

  void set_teids(uint32_t enb_teid, uint32_t agw_teid);

  void set_teids(Teids teids);

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

  void remove_event_trigger(
      magma::lte::EventTrigger trigger,
      SessionStateUpdateCriteria& update_criteria);

  void set_revalidation_time(
      const google::protobuf::Timestamp& time,
      SessionStateUpdateCriteria& update_criteria);

  google::protobuf::Timestamp get_revalidation_time() {
    return revalidation_time_;
  }

  EventTriggerStatus get_event_triggers() { return pending_event_triggers_; }

  optional<FinalActionInfo> get_final_action_if_final_unit_state(
      const CreditKey& ckey) const;

  RulesToProcess remove_all_final_action_rules(
      const FinalActionInfo& final_action_info,
      SessionStateUpdateCriteria& session_uc);

  // Monitors
  bool receive_monitor(
      const UsageMonitoringUpdateResponse& update,
      SessionStateUpdateCriteria& session_uc);

  uint64_t get_monitor(const std::string& key, Bucket bucket) const;

  bool set_monitor_reporting(
      const std::string& key, bool reporting,
      SessionStateUpdateCriteria* update_criteria);

  bool add_to_monitor(
      const std::string& key, uint64_t used_tx, uint64_t used_rx,
      SessionStateUpdateCriteria& uc);

  void set_monitor(
      const std::string& key, Monitor monitor, SessionStateUpdateCriteria& uc);

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
   * Determine whether a policy with type+ID needs a bearer to be created
   * @param policy_type
   * @param rule_id
   * @param config
   * @return an optional wrapped PolicyRule if creation is needed, {} otherwise
   */
  optional<PolicyRule> policy_needs_bearer_creation(
      const PolicyType policy_type, const std::string& rule_id,
      const SessionConfig& config);
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

  /**
   * Returns true if all the credits are suspended
   */
  void suspend_service_if_needed_for_credit(
      CreditKey ckey, SessionStateUpdateCriteria& update_criteria);

  /**
   * Returns true if the specified rule should be active at that time
   */
  bool should_rule_be_active(const std::string& rule_id, std::time_t time);
  bool is_dynamic_rule_scheduled(const std::string& rule_id);

  /**
   * Clear all per-session metrics
   */
  void clear_session_metrics();

 private:
  std::string imsi_;
  std::string session_id_;
  uint32_t local_teid_;
  uint32_t request_number_;
  SessionFsmState curr_state_;
  SessionConfig config_;
  uint64_t pdp_start_time_;
  uint64_t pdp_end_time_;
  /*5G related message to handle session state context */
  uint32_t current_version_;  // To compare with incoming session version
  // All 5G specific rules
  // use as shared_ptr to check
  std::vector<SetGroupPDR> PdrList_;
  // Used to keep track of whether the subscriber has valid quota.
  // (only used for CWF at the moment)
  magma::lte::SubscriberQuotaUpdate_Type subscriber_quota_state_;
  magma::lte::TgppContext tgpp_context_;

  // Used between create session and activate session. Empty afterwards
  CreateSessionResponse create_session_response_;

  // Track version tracking information used for LTE/WLAN
  PolicyStatsMap policy_version_and_stats_;

  // All static rules synced from policy DB
  StaticRuleStore& static_rules_;
  // Static rules that are currently installed for the session
  std::vector<std::string> active_static_rules_;
  // Dynamic GX rules that are currently installed for the session
  DynamicRuleStore dynamic_rules_;
  // Dynamic GY rules that are currently installed for the session
  DynamicRuleStore gy_dynamic_rules_;

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
  const uint32_t REDIRECT_FLOW_PRIORITY = 2000;

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

  /**
   * @brief Get a CreditUsageUpdate for the case where we want to continue
   * service
   *
   * @param grant
   * @param session_uc
   * @return optional<CreditUsageUpdate>
   */
  optional<CreditUsageUpdate> get_update_for_continue_service(
      const CreditKey& key, std::unique_ptr<ChargingGrant>& grant,
      SessionStateUpdateCriteria& session_uc);

  void fill_service_action_for_activate(
      std::unique_ptr<ServiceAction>& action, const CreditKey& key,
      SessionStateUpdateCriteria& session_uc);

  void fill_service_action_for_restrict(
      std::unique_ptr<ServiceAction>& action_p, const CreditKey& key,
      std::unique_ptr<ChargingGrant>& grant,
      SessionStateUpdateCriteria& session_uc);

  PolicyRule make_redirect_rule(std::unique_ptr<ChargingGrant>& grant);

  void fill_service_action_for_redirect(
      std::unique_ptr<ServiceAction>& action_p, const CreditKey& key,
      std::unique_ptr<ChargingGrant>& grant, PolicyRule redirect_rule,
      SessionStateUpdateCriteria& session_uc);

  void fill_service_action_with_context(
      std::unique_ptr<ServiceAction>& action, ServiceActionType action_type,
      const CreditKey& key);

  void apply_charging_credit_update(
      const CreditKey& key, SessionCreditUpdateCriteria& credit_uc);

  /**
   * Receive the credit grant if the credit update was successful
   *
   * @param update
   * @param session_uc
   * @return True if usage for the charging key is allowed after receiving
   *         the credit response. This requires that the credit update was
   *         a success.
   */
  bool init_charging_credit(
      const CreditUpdateResponse& update,
      SessionStateUpdateCriteria& session_uc);

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
   */
  void get_monitor_updates(
      UpdateSessionRequest& update_request_out,
      SessionStateUpdateCriteria& update_criteria);

  /**
   * Apply SessionCreditUpdateCriteria, a per-credit diff of an update, into
   * the SessionState object
   * @param key : monitoring key for the update
   * @param update : the diff that needs to be applied
   */
  void apply_monitor_updates(
      const std::string& key, SessionStateUpdateCriteria& session_uc,
      SessionCreditUpdateCriteria& credit_uc);

  void add_common_fields_to_usage_monitor_update(
      UsageMonitoringUpdateRequest* req);

  /**
   * Returns true if the specified rule should be deactivated by that time
   */
  bool should_rule_be_deactivated(const std::string& rule_id, std::time_t time);

  SessionCreditUpdateCriteria* get_credit_uc(
      const CreditKey& key, SessionStateUpdateCriteria& uc);

  CreditUsageUpdate make_credit_usage_update_req(CreditUsage& usage) const;

  bool init_new_monitor(
      const UsageMonitoringUpdateResponse& update,
      SessionStateUpdateCriteria& update_criteria);

  void update_session_level_key(
      const UsageMonitoringUpdateResponse& update,
      SessionStateUpdateCriteria& update_criteria);

  SessionCreditUpdateCriteria* get_monitor_uc(
      const std::string& key, SessionStateUpdateCriteria& uc);

  void fill_protos_tgpp_context(magma::lte::TgppContext* tgpp_context) const;

  void get_event_trigger_updates(
      UpdateSessionRequest& update_request_out,
      SessionStateUpdateCriteria& update_criteria);

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
   * @param usage_label either UE_DROPPED_LABEL / UE_USED_LABEL
   * @param bytes_tx
   * @param bytes_rx
   */
  void update_data_metrics(
      const char* counter_name, uint64_t bytes_tx, uint64_t bytes_rx);

  // PolicyStatsMap functions
  /**
   *
   * @param rule_id
   * @param session_uc
   */
  void increment_rule_stats(
      const std::string& rule_id, SessionStateUpdateCriteria& session_uc);
};

}  // namespace magma
