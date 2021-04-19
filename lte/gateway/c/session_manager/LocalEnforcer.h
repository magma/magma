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

#include <experimental/optional>
#include <folly/io/async/EventBaseManager.h>
#include <lte/protos/mconfig/mconfigs.pb.h>
#include <lte/protos/policydb.pb.h>
#include <lte/protos/session_manager.grpc.pb.h>

#include <iomanip>
#include <memory>
#include <string>
#include <unordered_map>
#include <unordered_set>
#include <utility>
#include <vector>

#include "AAAClient.h"
#include "DirectorydClient.h"
#include "PipelinedClient.h"
#include "RuleStore.h"
#include "SessionEvents.h"
#include "SessionReporter.h"
#include "SessionState.h"
#include "SessionStore.h"
#include "SpgwServiceClient.h"

namespace magma {
using std::experimental::optional;

typedef std::pair<std::string, std::string> ImsiAndSessionID;

struct RuleRecord_equal {
  bool operator()(const RuleRecord& l, const RuleRecord& r) const {
    return l.sid() == r.sid() && l.teid() == r.teid();
  }
};
struct RuleRecord_hash {
  std::size_t operator()(const RuleRecord& el) const {
    size_t h1 = std::hash<std::string>()(el.sid());
    size_t h2 = std::hash<uint32_t>()(el.teid());
    return h1 ^ h2;
  }
};
typedef std::unordered_set<RuleRecord, RuleRecord_hash, RuleRecord_equal>
    RuleRecordSet;

struct ImsiSessionIDAndCreditkey {
  std::string imsi;
  std::string session_id;
  CreditKey cKey;

  bool operator==(const ImsiSessionIDAndCreditkey& other) const {
    return this->imsi == other.imsi && this->session_id == other.session_id &&
           ccEqual(cKey, other.cKey);
  }
};

struct ImsiSessionIDAndCreditkey_hash {
  std::size_t operator()(const ImsiSessionIDAndCreditkey& el) const {
    size_t h1 = std::hash<std::string>()(el.imsi);
    size_t h2 = std::hash<std::string>()(el.session_id);
    size_t h3 = ccHash(el.cKey);
    return h1 ^ h2 ^ h3;
  }
};

struct UpdateChargingCreditActions {
  std::unordered_set<ImsiAndSessionID> sessions_to_terminate;
  std::unordered_set<ImsiSessionIDAndCreditkey, ImsiSessionIDAndCreditkey_hash>
      suspended_credits;
  std::unordered_set<ImsiSessionIDAndCreditkey, ImsiSessionIDAndCreditkey_hash>
      unsuspended_credits;
};

// Used to transform the UpdateSessionRequest proto message into a per-session
// structure
struct UpdateRequestsBySession {
  std::unordered_map<ImsiAndSessionID, UpdateRequests> requests_by_id;
  UpdateRequestsBySession() {}
  UpdateRequestsBySession(const magma::lte::UpdateSessionRequest& response);
};
class SessionNotFound : public std::exception {
 public:
  SessionNotFound() = default;
};

/**
 * LocalEnforcer can register traffic records and credits to track when a flow
 * has run out of credit
 */
class LocalEnforcer {
 public:
  LocalEnforcer();

  LocalEnforcer(
      std::shared_ptr<SessionReporter> reporter,
      std::shared_ptr<StaticRuleStore> rule_store, SessionStore& session_store,
      std::shared_ptr<PipelinedClient> pipelined_client,
      std::shared_ptr<EventsReporter> events_reporter,
      std::shared_ptr<SpgwServiceClient> spgw_client,
      std::shared_ptr<aaa::AAAClient> aaa_client,
      long session_force_termination_timeout_ms,
      long quota_exhaustion_termination_on_init_ms,
      magma::mconfig::SessionD mconfig);

  void attachEventBase(folly::EventBase* evb);

  // blocks
  void start();

  void stop();

  folly::EventBase& get_event_base();

  /**
   * Setup rules for all sessions in pipelined, used whenever pipelined
   * restarts and needs to recover state
   */
  void setup(
      SessionMap& session_map, const std::uint64_t& epoch,
      std::function<void(Status status, SetupFlowsResult)> callback);

  /**
   * Updates rules to be activated/deactivated based on the current time.
   * Also schedules future rule activation and deactivation callbacks to run
   * on the event loop.
   */
  void sync_sessions_on_restart(std::time_t current_time);

  /**
   * Insert a group of rule usage into the monitor and update credit manager
   * Assumes records are aggregates, as in the usages sent are cumulative and
   * not differences.
   *
   * @param records - a RuleRecordTable protobuf with a vector of RuleRecords
   */
  void aggregate_records(
      SessionMap& session_map, const RuleRecordTable& records,
      SessionUpdate& session_update);

  /**
   * handle_update_failure resets all of the charging keys / monitors being
   * updated in failed_request. This should only be called if the *entire*
   * request fails (i.e. the entire request to the cloud timed out). Individual
   * failures are handled when update_session_credits_and_rules is called.
   *
   * @param failed_request - UpdateRequestsBySession that couldn't be sent to
   * the cloud for whatever reason
   */
  void handle_update_failure(
      SessionMap& session_map, const UpdateRequestsBySession& failed_request,
      SessionUpdate& updates);

  /**
   * Collect any credit keys that are either exhausted, timed out, or terminated
   * and apply actions to the services if need be
   * @param updates_out (out) - vector to add usage updates to, if they exist
   */
  UpdateSessionRequest collect_updates(
      SessionMap& session_map,
      std::vector<std::unique_ptr<ServiceAction>>& actions,
      SessionUpdate& session_update) const;

  /**
   * Perform any rule installs/removals that need to be executed given a
   * CreateSessionResponse.
   */
  void handle_session_activate_rule_updates(
      SessionState& session, const CreateSessionResponse& response,
      std::unordered_set<uint32_t>& charging_credits_received);

  void schedule_session_init_dedicated_bearer_creations(
      const std::string& imsi, const std::string& session_id,
      BearerUpdate& bearer_updates);

  /**
   * Initialize session on session map. Adds some information comming from
   * the core (cloud). Rules will be installed by init_session_credit
   * @param credit_response - message from cloud containing initial credits
   */
  void init_session(
      SessionMap& session_map, const std::string& imsi,
      const std::string& session_id, const SessionConfig& cfg,
      const CreateSessionResponse& response);

  /**
   * Process the update response from the reporter and update the
   * monitoring/charging credits and attached rules.
   * @param credit_response - message from cloud containing new credits
   */
  void update_session_credits_and_rules(
      SessionMap& session_map, const UpdateSessionResponse& response,
      SessionUpdate& session_update);

  /**
   * handle_termination_from_access handles externally triggered session
   * termination. This assumes that the termination is coming from the access
   * component, so it does not notify the termination back to the access
   * component.
   * @param session_map
   * @param imsi
   * @param apn
   * @param session_update
   */
  bool handle_termination_from_access(
      SessionMap& session_map, const std::string& imsi, const std::string& apn,
      SessionUpdate& session_updates);

  /**
   * handle abort session - ungraceful termination
   * 1. Remove all rules attached to the session + update PipelineD
   * 2. Notify the access component
   * 3. Remove the session from SessionMap
   */
  bool handle_abort_session(
      SessionMap& session_map, const std::string& imsi,
      const std::string& session_id, SessionUpdate& session_updates);

  /**
   * Initialize reauth for a subscriber service. If the subscriber cannot be
   * found, the method returns SESSION_NOT_FOUND
   */
  ReAuthResult init_charging_reauth(
      SessionMap& session_map, ChargingReAuthRequest request,
      SessionUpdate& session_update);

  /**
   * Handles the equivalent of a RAR.
   * For the matching session ID, activate and/or deactivate the specified
   * rules.
   * Afterwards, a bearer is created.
   * If a session is CWF and out of monitoring quota, it will trigger a session
   * terminate
   *
   * NOTE: If an empty session ID is specified, apply changes to all matching
   * sessions with the specified IMSI.
   */
  void init_policy_reauth(
      SessionMap& session_map, PolicyReAuthRequest request,
      PolicyReAuthAnswer& answer_out, SessionUpdate& session_update);

  /**
   * Set session config for the IMSI.
   * Should be only used for WIFI as it will apply it to all sessions with the
   * IMSI
   */
  void handle_cwf_roaming(
      SessionMap& session_map, const std::string& imsi,
      const magma::SessionConfig& config, SessionUpdate& session_update);

  /**
   * Execute actions on subscriber's service, eg. terminate, redirect data, or
   * just continue
   */
  void execute_actions(
      SessionMap& session_map,
      const std::vector<std::unique_ptr<ServiceAction>>& actions,
      SessionUpdate& session_update);

  /**
   * handle_set_session_rules takes SessionRules, which is a set message that
   * reflects the desired rule state, and apply the changes. The changes should
   * be propagated to PipelineD and MME if the session is 4G.
   * @param session_map
   * @param updates
   * @param session_update
   */
  void handle_set_session_rules(
      SessionMap& session_map, const SessionRules& rules,
      SessionUpdate& session_update);

  /**
   * Check if PolicyBearerBindingRequest has a non-zero dedicated bearer ID:
   * Update the policy to bearer map if non-zero
   * Delete the policy rule if zero
   * @return true if successfully processed the request
   */
  bool bind_policy_to_bearer(
      SessionMap& session_map, const PolicyBearerBindingRequest& request,
      SessionUpdate& session_update);

  /**
   * Sends enb_teid and agw_teid for a specific bearer to a flow for a specific
   * UE on pipelined. UE will be identified by pipelined using its IP
   * @param session_map
   * @param request
   * @return true if successfully processed the request
   */
  bool update_tunnel_ids(
      SessionMap& session_map, const UpdateTunnelIdsRequest& request);

  std::unique_ptr<Timezone>& get_access_timezone() { return access_timezone_; };

  static uint32_t REDIRECT_FLOW_PRIORITY;
  static uint32_t BEARER_CREATION_DELAY_ON_SESSION_INIT;
  // If this is set to true, we will send the timezone along with
  // CreateSessionRequest
  static bool SEND_ACCESS_TIMEZONE;
  // If true, for any rule reported as part of ReportRuleStats,
  // remove it if the rule's IMSI+TEIDs pair do no exist as
  // a session
  static bool CLEANUP_DANGLING_FLOWS;

 private:
  std::shared_ptr<SessionReporter> reporter_;
  std::shared_ptr<StaticRuleStore> rule_store_;
  std::shared_ptr<PipelinedClient> pipelined_client_;
  std::shared_ptr<EventsReporter> events_reporter_;
  std::shared_ptr<SpgwServiceClient> spgw_client_;
  std::shared_ptr<aaa::AAAClient> aaa_client_;
  SessionStore& session_store_;
  folly::EventBase* evb_;
  long session_force_termination_timeout_ms_;
  // [CWF-ONLY] This configures how long we should wait before terminating a
  // session after it is created without any monitoring quota
  long quota_exhaustion_termination_on_init_ms_;
  std::chrono::milliseconds retry_timeout_;
  magma::mconfig::SessionD mconfig_;
  std::unique_ptr<Timezone> access_timezone_;

 private:
  /**
   * complete_termination_for_released_sessions completes the termination
   * process for sessions whose flows have been removed in PipelineD. Since
   * PipelineD reports all rule records that exist in PipelineD with each
   * report, if the session is not included, that means the enforcement flows
   * have been removed.
   * @param session_map
   * @param sessions_with_active_flows: a set of IMSIs whose rules were reported
   * @param session_update
   */
  void complete_termination_for_released_sessions(
      SessionMap& session_map,
      std::unordered_set<ImsiAndSessionID> sessions_with_active_flows,
      SessionUpdate& session_update);

  void filter_rule_installs(
      std::vector<StaticRuleInstall>& static_rule_installs,
      std::vector<DynamicRuleInstall>& dynamic_rule_installs,
      const std::unordered_set<uint32_t>& successful_credits);

  std::vector<StaticRuleInstall> to_vec(
      const google::protobuf::RepeatedPtrField<magma::lte::StaticRuleInstall>
          static_rule_installs);
  std::vector<DynamicRuleInstall> to_vec(
      const google::protobuf::RepeatedPtrField<magma::lte::DynamicRuleInstall>
          dynamic_rule_installs);

  /**
   * Processes the charging component of UpdateSessionResponse.
   * Updates charging credits according to the response.
   */
  void update_charging_credits(
      SessionMap& session_map, const UpdateSessionResponse& response,
      UpdateChargingCreditActions& actions, SessionUpdate& session_update);

  /**
   * Processes the monitoring component of UpdateSessionResponse.
   * Updates moniroting credits according to the response and updates rules
   * that are installed for this session.
   * If a session is CWF and out of monitoring quota, it will trigger a session
   * terminate
   */
  void update_monitoring_credits_and_rules(
      SessionMap& session_map, const UpdateSessionResponse& response,
      UpdateChargingCreditActions& actions, SessionUpdate& session_update);

  /**
   * @brief For rules mentioned in both static_rule_installs and
   * dynamic_rule_installs, classify them into the three RulesToProcess vectors.
   * pending_activation, pending_deactivation, pending_bearer_setup will not
   * intersect in the set of rules they contain pending_scheduling may contain
   * some deactivation scheduling for rules mentioned in the above three sets
   * @param session
   * @param static_rule_installs
   * @param dynamic_rule_installs
   * @param pending_activation contains rules that need to be activated now
   * @param pending_deactivation contains rules that need to be deactivated now
   * @param pending_bearer_setup contains rules that need to get dedicated
   * bearers before they can be activated. The rules will be activated once MME
   * sends a BindPolicy2Bearer with the dedicated bearer Teids.
   * @param pending_scheduling contains rules that need to be scheduled to be
   * activated/deactivated
   * @param session_uc
   */
  void process_rules_to_install(
      SessionState& session,
      const std::vector<StaticRuleInstall>& static_rule_installs,
      const std::vector<DynamicRuleInstall>& dynamic_rule_installs,
      RulesToProcess* pending_activation, RulesToProcess* pending_deactivation,
      RulesToProcess* pending_bearer_setup, RulesToSchedule* pending_scheduling,
      SessionStateUpdateCriteria* session_uc);

  /**
   * propagate_rule_updates_to_pipelined calls the PipelineD RPC calls to
   * install/uninstall flows
   * @param config
   * @param pending_activation
   * @param pending_deactivation
   * @param always_send_activate : if this is set activate call will be sent
   * even if pending_activation is empty
   */
  void propagate_rule_updates_to_pipelined(
      const SessionConfig& config, const RulesToProcess& pending_activation,
      const RulesToProcess& pending_deactivation, bool always_send_activate);

  /**
   * @brief for each element in RulesToSchedule, schedule rule
   * activation/deactivation on the event loop
   *
   * @param imsi
   * @param session_id
   * @param pending_scheduling
   */
  void handle_rule_scheduling(
      const std::string& imsi, const std::string& session_id,
      const RulesToSchedule& pending_scheduling);

  /**
   * For the matching session ID, activate and/or deactivate the specified
   * rules.
   * Also create a bearer for the session.
   */
  void init_policy_reauth_for_session(
      const PolicyReAuthRequest& request,
      const std::unique_ptr<SessionState>& session,
      SessionUpdate& session_update);

  /**
   * find_and_terminate_session call start_session_termination on
   * a session with IMSI + session id.
   * @return true if start_session_termination was called, false if session was
   * not found
   */
  bool find_and_terminate_session(
      SessionMap& session_map, const std::string& imsi,
      const std::string& session_id, SessionUpdate& session_updates);

  /**
   * Completes the session termination and executes the callback function
   * registered in terminate_session.
   * complete_termination is called some time after terminate_session
   * when the flows no longer appear in the usage report, meaning that they have
   * been deleted.
   * It is also called after a timeout to perform forced termination.
   * If the session cannot be found, either because it has already terminated,
   * or a new session for the subscriber has been created, then it will do
   * nothing.
   */
  void complete_termination(
      SessionMap& session_map, const std::string& imsi,
      const std::string& session_id, SessionUpdate& session_update);

  void schedule_static_rule_activation(
      const std::string& imsi, const std::string& session_id,
      const std::string& rule_id, const std::time_t activation_time);

  void schedule_dynamic_rule_activation(
      const std::string& imsi, const std::string& session_id,
      const std::string& rule_id, const std::time_t activation_time);

  void schedule_static_rule_deactivation(
      const std::string& imsi, const std::string& session_id,
      const std::string& rule_id, const std::time_t deactivation_time);

  void schedule_dynamic_rule_deactivation(
      const std::string& imsi, const std::string& session_id,
      const std::string& rule_id, const std::time_t deactivation_time);

  /**
   * Get the monitoring credits from PolicyReAuthRequest (RAR) message
   * and add the credits to UsageMonitoringCreditPool of the session
   */
  void receive_monitoring_credit_from_rar(
      const PolicyReAuthRequest& request,
      const std::unique_ptr<SessionState>& session,
      SessionStateUpdateCriteria& uc);

  /**
   * Send bearer creation request through the PGW client if rules were
   * activated successfully in pipelined
   */
  void create_bearer(
      const std::unique_ptr<SessionState>& session,
      const PolicyReAuthRequest& request,
      const std::vector<RuleToProcess>& dynamic_rules);

  /**
   * Check if REVALIDATION_TIMEOUT is one of the event triggers
   */
  bool revalidation_required(
      const google::protobuf::RepeatedField<int>& event_triggers);

  void schedule_revalidation(
      SessionState& session,
      const google::protobuf::Timestamp& revalidation_time,
      SessionStateUpdateCriteria& uc);

  void handle_add_ue_mac_flow_callback(
      const SubscriberID& sid, const std::string& ue_mac_addr,
      const std::string& msisdn, const std::string& ap_mac_addr,
      const std::string& ap_name, Status status, FlowResponse resp);

  void handle_activate_ue_flows_callback(
      const std::string& imsi, const std::string& ip_addr,
      const std::string& ipv6_addr, const Teids teids, Status status,
      ActivateFlowsResult resp);

  /**
   * start_session_termination starts the termination process. This includes:
   * 1. Update the Session FSM State to Terminating
   * 2. Remove all policies attached to the session
   * 3. If notify_access param is set, communicate to the access component
   * 4. Propagate subscriber wallet status
   * 5. Schedule a callback to force termination if termination is not completed
   *    in a set amount of time
   * @param session
   * @param notify_access: bool to determine whether the access component needs
   * notification
   * @param uc
   */
  void start_session_termination(
      const std::unique_ptr<SessionState>& session, bool notify_access,
      SessionStateUpdateCriteria& uc);

  /**
   * handle_force_termination_timeout is scheduled to run when a termination
   * process starts. If the session did not terminate itself properly within the
   * timeout, this function will force the termination to complete.
   * @param imsi
   * @param session_id
   */
  void handle_force_termination_timeout(
      const std::string& imsi, const std::string& session_id);

  /**
   * remove_all_rules_for_termination talks to PipelineD and removes all rules
   * (Gx/Gy/static/dynamic/everything) attached to the session
   * @param session
   * @param uc
   */
  void remove_all_rules_for_termination(
      const std::unique_ptr<SessionState>& session,
      SessionStateUpdateCriteria& uc);

  /**
   * notify_termination_to_access_service cases on the session's rat type and
   * communicates to the appropriate access client to notify the session's
   * termination.
   * LTE -> MME, WLAN -> AAA
   * @param session_id
   * @param config
   */
  void notify_termination_to_access_service(
      const std::string& session_id, const SessionConfig& config);
  /**
   * handle_subscriber_quota_state_change will update the session's wallet state
   * to the desired new_state and propagate that state PipelineD.
   * @param session
   * @param new_state
   * @param session_uc
   */
  void handle_subscriber_quota_state_change(
      SessionState& session, SubscriberQuotaUpdate_Type new_state,
      SessionStateUpdateCriteria* session_uc);

  /**
   * Start the termination process for multiple sessions
   */
  void terminate_multiple_sessions(
      SessionMap& session_map,
      const std::unordered_set<ImsiAndSessionID>& sessions,
      SessionUpdate& session_update);

  void remove_rules_for_multiple_suspended_credit(
      SessionMap& session_map,
      std::unordered_set<
          ImsiSessionIDAndCreditkey, ImsiSessionIDAndCreditkey_hash>&
          suspended_credits,
      SessionUpdate& session_update);

  void add_rules_for_multiple_unsuspended_credit(
      SessionMap& session_map,
      std::unordered_set<
          ImsiSessionIDAndCreditkey, ImsiSessionIDAndCreditkey_hash>&
          unsuspended_credits,
      SessionUpdate& session_update);

  void handle_activate_service_action(
      const std::unique_ptr<ServiceAction>& action_p);

  /**
   * Install final action flows through pipelined
   */
  void install_final_unit_action_flows(
      const std::unique_ptr<ServiceAction>& action);

  /**
   * Create redirection rule
   */
  PolicyRule create_redirect_rule(const std::unique_ptr<ServiceAction>& action);

  void report_subscriber_state_to_pipelined(
      const std::string& imsi, const std::string& ue_mac_addr,
      const SubscriberQuotaUpdate_Type state);

  void update_ipfix_flow(
      const std::string& imsi, const SessionConfig& config,
      const uint64_t pdp_start_time);

  /**
   * If the session has active monitored rules attached to it, then propagate
   * to pipelined that the subscriber has valid quota.
   * Otherwise, mark the subscriber as out of quota to pipelined, and schedule
   * the session to be terminated in a configured amount of time.
   */
  void handle_session_activate_subscriber_quota_state(SessionState& session);

  bool is_wallet_exhausted(SessionState& session);

  bool terminate_on_wallet_exhaust();

  void schedule_termination(std::unordered_set<ImsiAndSessionID>& sessions);

  void propagate_bearer_updates_to_mme(const BearerUpdate& updates);

  /**
   * Remove the specified rule from the session and propagate the change to
   * PipelineD
   * @param rule_id rule to be deleted
   * @param uc
   */
  void remove_rule_due_to_bearer_creation_failure(
      SessionState& session, const std::string& rule_id,
      SessionStateUpdateCriteria& uc);

  /**
   * @brief Activate the rule after successfully binding it to a dedicated
   * bearer
   *
   * @param session
   * @param request
   */
  void install_rule_after_bearer_creation(
      SessionState& session, const PolicyBearerBindingRequest& request);

  static std::unique_ptr<Timezone> compute_access_timezone();

  void remove_rules_for_suspended_credit(
      const std::unique_ptr<SessionState>& session, const CreditKey& ckey,
      SessionStateUpdateCriteria& session_uc);

  void add_rules_for_unsuspended_credit(
      const std::unique_ptr<SessionState>& session, const CreditKey& ckey,
      SessionStateUpdateCriteria& session_uc);

  /**
   * Given a set of IMSI+IPs that are no longer tracked in SessionD, send a
   * deactivate flows request to PipelineD for all flows associated with those
   * IDs. The function sends the request for all types (ANY), because the set
   * does not specify the origin type (Gx/Gy/N4).
   * @param dead_sessions_to_cleanup
   */
  void cleanup_dead_sessions(const RuleRecordSet dead_sessions_to_cleanup);
};

}  // namespace magma
