/*
Copyright 2020 The Magma Authors.
This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*****************************************************************************
  Source      	SessionStateEnforcer.h
  Version     	0.1
  Date       	2020/08/08
  Product     	SessionD
  Subsystem   	5G managing & maintaining state & store of session of SessionD
                Fanout message to Access and UPF through respective client obj
  Author/Editor Sanjay Kumar Ojha
  Description 	Objects run in main thread context invoked by folly event
*****************************************************************************/

#pragma once

#include <folly/io/async/EventBaseManager.h>
#include <lte/protos/mconfig/mconfigs.pb.h>
#include <lte/protos/policydb.pb.h>
#include <stdint.h>
#include <chrono>
#include <map>
#include <memory>
#include <string>
#include <unordered_map>
#include <unordered_set>
#include <vector>

#include "AmfServiceClient.h"
#include "PipelinedClient.h"
#include "RuleStore.h"
#include "SessionState.h"
#include "SessionStore.h"
#include "StoreClient.h"
#include "SessionEvents.h"
#include "Types.h"
#include "lte/protos/pipelined.pb.h"
#include "lte/protos/session_manager.pb.h"

namespace folly {
class EventBase;
}  // namespace folly
namespace magma {
class AmfServiceClient;
class PipelinedClient;
class SessionState;
struct SessionStateUpdateCriteria;
}  // namespace magma

#define M5G_MIN_TEID (UINT32_MAX / 2)
#define DEFAULT_PDR_VERSION 1
#define DEFAULT_PDR_ID 0
#define DEFAULT_PDR_PRECEDENCE 32

namespace magma {

class SessionStateEnforcer {
 public:
  SessionStateEnforcer(std::shared_ptr<StaticRuleStore> rule_store,
                       SessionStore& session_store,
                       /*M5G specific parameter new objects to communicate UPF
                          and response to AMF*/
                       std::unordered_multimap<std::string, uint32_t> pdr_map,
                       std::shared_ptr<PipelinedClient> pipelined_client,
                       std::shared_ptr<AmfServiceClient> amf_srv_client,
                       SessionReporter* reporter,
                       std::shared_ptr<EventsReporter> events_reporter,
                       magma::mconfig::SessionD mconfig,
                       long session_force_termination_timeout_ms,
                       uint32_t session_max_rtx_count);

  ~SessionStateEnforcer() {}

  void attachEventBase(folly::EventBase* evb);

  void stop();

  folly::EventBase& get_event_base();

  /*Member functions*/
  bool m5g_init_session_credit(SessionMap& session_map, const std::string& imsi,
                               const std::string& session_id,
                               const SessionConfig& cfg);

  /*Member functions*/
  bool m5g_update_session_context(SessionMap& session_map,
                                  const std::string& imsi,
                                  std::unique_ptr<SessionState>& session_state,
                                  SessionUpdate& session_update);

  /*Charging & rule related*/
  bool handle_session_init_rule_updates(
      SessionMap& session_map, const std::string& imsi,
      std::unique_ptr<SessionState>& session_state);

  /* Move to idle state */
  void m5g_move_to_inactive_state(const std::string& imsi,
                                  std::unique_ptr<SessionState>& session,
                                  SetSmNotificationContext notif,
                                  SessionStateUpdateCriteria* session_uc);

  /* Move to active state */
  void m5g_move_to_active_state(std::unique_ptr<SessionState>& session,
                                SetSmNotificationContext notif,
                                SessionStateUpdateCriteria* session_uc);

  /*Release request handle*/
  bool m5g_release_session(SessionMap& session_map, const std::string& imsi,
                           const uint32_t& pdu_id,
                           SessionUpdate& session_update);

  /*Handle and update respective session upon receiving message from UPF*/
  void m5g_process_response_from_upf(const std::string& imsi, uint32_t teid,
                                     uint32_t version);

  /* Send session requst to upf */
  void m5g_send_session_request_to_upf(
      const std::unique_ptr<SessionState>& session,
      const RulesToProcess& pending_activation,
      const RulesToProcess& pending_deactivation);

  std::vector<std::string> static_rules;

  /* Get N3 ip  address of UPF */
  std::string get_upf_n3_addr() const;

  /* Get N3 ip  address of UPF */
  std::string get_upf_node_id() const;

  /* Initialize the upf node id and n3 address
   */
  bool set_upf_node(const std::string& node_id, const std::string& n3_addr);

  /**
   * If UPF received session version doesn't match with SMF local session
   * version no, then we continue to resend till SESSION_THROTTLE_CNT reaches.
   * Gets the incremented rtx counter and compares with SESSION_THROTTLE_CNT
   *
   * @param session : reference to SessionState
   * @return true : if incremented rtx counter < SESSION_THROTTLE_CNT
   * 	     Otherwise false
   */
  bool is_incremented_rtx_counter_within_max(
      const std::unique_ptr<SessionState>& session);

  /*
   * To send session or UE notifcaiton
   * to AMF, this cud be session state or
   * paging notification
   */
  void handle_state_update_to_amf(SessionState& session_state,
                                  const magma::lte::M5GSMCause m5gsmcause,
                                  NotifyUeEvents event);

  /* Get next teid */
  uint32_t get_next_teid();

  /* Get current TEID */
  uint32_t get_current_teid();

  bool add_default_rules(std::unique_ptr<SessionState>& session_state,
                         const std::string& imsi);

  /* Pdr State change routine */
  void m5g_pdr_rules_change_and_update_upf(
      const std::unique_ptr<SessionState>& session, enum PdrState pdrstate);

  /*Start processing to terminate respective session requested from AMF*/
  void m5g_start_session_termination(
      SessionMap& session_map, const std::unique_ptr<SessionState>& session,
      const uint32_t& pdu_id, SessionStateUpdateCriteria* session_uc);

  /* Set new fsm state and increment version*/
  void set_new_fsm_state_and_increment_version(
      std::unique_ptr<SessionState>& session, SessionFsmState target_state,
      SessionStateUpdateCriteria* session_uc);

  /* update the GNB endpoint details in a rule */
  bool insert_pdr_from_core(std::unique_ptr<SessionState>& session,
                            SetGroupPDR& rule,
                            SessionStateUpdateCriteria* session_uc);

  uint32_t insert_pdr_from_access(std::unique_ptr<SessionState>& session,
                                  SetGroupPDR& rule,
                                  SessionStateUpdateCriteria* session_uc);

  /*
   * Acquire and update session rules based on the IMSI
   * @param session : reference to SessionState
   * @param get_gnb_teid : true if existing session has gNB TEID, else false
   * @param get_upf_teid : true if existing session has UPF TEID, else false
   * @return upf_teid : returns updated UPF TEID from insert_pdr_from_access()
   * */
  uint32_t update_session_rules(std::unique_ptr<SessionState>& session_state,
                                bool get_gnb_teid, bool get_upf_teid,
                                SessionStateUpdateCriteria* session_uc);

  /*Function will clean up all resources related to requested session*/
  void m5g_complete_termination(SessionMap& session_map,
                                const std::string& imsi,
                                const std::string& session_id,
                                SessionUpdate& session_update);

  std::unique_ptr<Timezone>& get_access_timezone() { return access_timezone_; }

  void update_session_with_policy(std::unique_ptr<SessionState>& session,
                                  const CreateSessionResponse& response,
                                  SessionStateUpdateCriteria* session_uc);

  std::vector<StaticRuleInstall> to_vec(
      const google::protobuf::RepeatedPtrField<magma::lte::StaticRuleInstall>
          static_rule_installs);

  std::vector<DynamicRuleInstall> to_vec(
      const google::protobuf::RepeatedPtrField<magma::lte::DynamicRuleInstall>
          dynamic_rule_installs);

 private:
  ConvergedRuleStore GlobalRuleList;
  std::shared_ptr<StaticRuleStore> rule_store_;
  SessionStore& session_store_;
  std::unordered_multimap<std::string, uint32_t> pdr_map_;
  std::shared_ptr<PipelinedClient> pipelined_client_;
  std::shared_ptr<AmfServiceClient> amf_srv_client_;
  SessionReporter* reporter_;
  std::shared_ptr<EventsReporter> events_reporter_;
  magma::mconfig::SessionD mconfig_;
  // Timer used to forcefully terminate session context on time out
  long session_force_termination_timeout_ms_;
  uint32_t session_max_rtx_count_;
  folly::EventBase* evb_;
  std::chrono::seconds retry_timeout_;
  std::string upf_node_id_;
  uint32_t teid_counter_;
  std::string upf_node_ip_addr_;
  std::unique_ptr<Timezone> access_timezone_;

  bool default_and_static_rule_init();

  /* To send response back to AMF
   * Fill the response structure and call rpc of AmfServiceClient
   */
  void prepare_response_to_access(SessionState& session_state,
                                  const magma::lte::M5GSMCause m5gsmcause,
                                  std::string upf_ip, uint32_t upf_teid);

  /* Function to handle termination if UPF doesn't send required report
   * As per current implementation, upf report is not in place and
   * termination on time out will be executed forcefully
   */
  void m5g_handle_termination_on_timeout(const std::string& imsi,
                                         const std::string& session_id);

  bool inc_rtx_counter(const std::unique_ptr<SessionState>& session);

  void set_pdr_attributes(const std::string& imsi,
                          std::unique_ptr<SessionState>& session_state,
                          SetGroupPDR* rule);

};  // End of class SessionStateEnforcer

}  // end namespace magma
