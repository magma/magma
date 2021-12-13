/*
 * Copyright 2020 The Magma Authors.
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*****************************************************************************
  Source      	SessionStateEnforcer.cpp
  Version     	0.1
  Date       	2020/08/08
  Product     	SessionD
  Subsystem   	5G managing & maintaining state & store of session of SessionD
                Fanout message to Access and UPF through respective client obj
  Author/Editor Sanjay Kumar Ojha
  Description 	Objects run in main thread context invoked by folly event
*****************************************************************************/

#include <folly/io/async/EventBase.h>
#include <glog/logging.h>
#include <grpcpp/impl/codegen/status.h>
#include <algorithm>
#include <cstdint>
#include <experimental/optional>
#include <memory>
#include <ostream>
#include <string>
#include <utility>
#include <vector>

#include "AmfServiceClient.h"
#include "EnumToString.h"
#include "PipelinedClient.h"
#include "SessionState.h"
#include "SessionStateEnforcer.h"
#include "StoredState.h"
#include "lte/protos/apn.pb.h"
#include "lte/protos/mconfig/mconfigs.pb.h"
#include "lte/protos/policydb.pb.h"
#include "lte/protos/subscriberdb.pb.h"
#include "magma_logging.h"

#define DEFAULT_AMBR_UNITS (1024)
#define DEFAULT_UP_LINK_PDR_ID 1
#define DEFAULT_DOWN_LINK_PDR_ID 2
#define DEFAULT_RULE_COUNT 2

std::shared_ptr<magma::SessionStateEnforcer> conv_session_enforcer;
namespace magma {

void call_back_upf(grpc::Status, magma::UPFSessionContextState response) {
  std::string imsi = response.session_snapshot().subscriber_id();
  uint32_t version = response.session_snapshot().session_version();
  uint32_t fteid = response.session_snapshot().local_f_teid();
  const std::string session_id = response.session_snapshot().subscriber_id();
  MLOG(MDEBUG) << " Async Response received from UPF: imsi " << imsi
               << " local fteid : " << fteid;
  conv_session_enforcer->get_event_base().runInEventBaseThread([imsi, fteid,
                                                                version]() {
    /* Update the state change, and notifiy to AMF */
    // For now fteid will be zero in all cases
    conv_session_enforcer->m5g_process_response_from_upf(imsi, fteid, version);
  });
}

/*constructor*/
SessionStateEnforcer::SessionStateEnforcer(
    std::shared_ptr<StaticRuleStore> rule_store, SessionStore& session_store,
    std::unordered_multimap<std::string, uint32_t> pdr_map,
    std::shared_ptr<PipelinedClient> pipelined_client,
    std::shared_ptr<AmfServiceClient> amf_srv_client, SessionReporter* reporter,
    std::shared_ptr<EventsReporter> events_reporter,
    magma::mconfig::SessionD mconfig, long session_force_termination_timeout_ms,
    uint32_t session_max_rtx_count)
    : rule_store_(rule_store),
      session_store_(session_store),
      pdr_map_(pdr_map),
      pipelined_client_(pipelined_client),
      amf_srv_client_(amf_srv_client),
      reporter_(reporter),
      events_reporter_(events_reporter),
      mconfig_(mconfig),
      session_force_termination_timeout_ms_(
          session_force_termination_timeout_ms),
      session_max_rtx_count_(session_max_rtx_count),
      retry_timeout_(1) {
  default_and_static_rule_init();
  teid_counter_ = M5G_MIN_TEID;
}

void SessionStateEnforcer::attachEventBase(folly::EventBase* evb) {
  evb_ = evb;
}

void SessionStateEnforcer::stop() { evb_->terminateLoopSoon(); }

folly::EventBase& SessionStateEnforcer::get_event_base() { return *evb_; }

bool SessionStateEnforcer::m5g_init_session_credit(
    SessionMap& session_map, const std::string& imsi,
    const std::string& session_id, const SessionConfig& cfg) {
  /* creating SessionState object with state CREATING
   * This calls constructor and allocates memory*/
  // TODO(veshkemburu): Figure out how to handle sharding(object)
  auto session_state =
      std::make_unique<SessionState>(imsi, session_id, cfg, *rule_store_);
  MLOG(MDEBUG) << " New SessionState object created with IMSI: " << imsi
               << " session context id : " << session_id;

  if (!handle_session_init_rule_updates(session_map, imsi, session_state)) {
    MLOG(MERROR) << " Session is not updated the rules IMSI::" << imsi;
    return false;
  }

  /* Find same UE or imsi already present, if not add
   * TODO - Need to check if same DNN/APN already exist
   */
  auto exist_imsi = session_map.find(imsi);
  if (exist_imsi == session_map.end()) {
    // First time a session is created for IMSI in the SessionMap
    session_map[imsi] = std::vector<std::unique_ptr<SessionState>>();
  } else {
    session_map[imsi].push_back(std::move(session_state));
  }
  MLOG(MDEBUG) << "Added a session (" << session_map[imsi].size()
               << ") for IMSI " << imsi << " with session context ID ";
  return true;
}

bool SessionStateEnforcer::m5g_update_session_context(
    SessionMap& session_map, const std::string& imsi,
    std::unique_ptr<SessionState>& session_state,
    SessionUpdate& session_update) {
  uint32_t upf_teid;
  bool get_gnb_teid = true;
  bool get_upf_teid = false;
  RulesToProcess pending_activation, pending_deactivation;
  RulesToSchedule pending_scheduling;

  /* Check and update latest session rules
   * we get gnodeb TEID, and IP address details here
   */
  auto session_id = session_state->get_session_id();
  SessionStateUpdateCriteria& session_uc = session_update[imsi][session_id];
  upf_teid = update_session_rules(session_state, get_gnb_teid, get_upf_teid,
                                  &session_uc);
  if (!upf_teid) {
    return false;
  }
  SessionSearchCriteria criteria(imsi, IMSI_AND_SESSION_ID, session_id);
  auto session_it = session_store_.find_session(session_map, criteria);
  if (!session_it) {
    MLOG(MERROR) << "No session found in SessionMap for IMSI " << imsi
                 << " with session_id " << session_id;
  }
  auto& session = **session_it;
  const CreateSessionResponse& csr = session->get_create_session_response();
  uint32_t cur_version = session_state->get_current_version();
  session_state->set_fsm_state(CREATED, &session_uc);
  session_state->set_current_version(++cur_version, &session_uc);
  std::vector<StaticRuleInstall> static_rule_installs =
      to_vec(csr.static_rules());
  std::vector<DynamicRuleInstall> dynamic_rule_installs =
      to_vec(csr.dynamic_rules());
  session->process_rules_to_install(static_rule_installs, dynamic_rule_installs,
                                    &pending_activation, &pending_deactivation,
                                    nullptr, &pending_scheduling, &session_uc);

  std::unordered_set<uint32_t> charging_credits_received;
  for (const auto& credit : csr.credits()) {
    if (session->receive_charging_credit(credit, &session_uc)) {
      charging_credits_received.insert(credit.charging_key());
    }
  }
  for (const auto& monitor : csr.usage_monitors()) {
    auto uc = get_default_update_criteria();
    session->receive_monitor(monitor, &session_uc);
  }
  /* Reset the upf resend retransmission counter counter, send the session
   * creation request to UPF
   */
  session_state->reset_rtx_counter();
  m5g_send_session_request_to_upf(session_state, pending_activation,
                                  pending_deactivation);
  events_reporter_->session_created(imsi, session_id, session->get_config(),
                                    session);
  return true;
}

uint32_t SessionStateEnforcer::update_session_rules(
    std::unique_ptr<SessionState>& session, bool get_gnb_teid,
    bool get_upf_teid, SessionStateUpdateCriteria* session_uc) {
  uint32_t upf_teid = 0;
  SetGroupPDR rule;
  const std::string& imsi = session->get_imsi();

  // Get the latest config
  const auto& config = session->get_config();
  auto itp = pdr_map_.equal_range(imsi);
  // Lets take local_teid of the session if already exists.
  // Not needed in 2nd AMF request processing though
  if (!get_upf_teid) {
    upf_teid = session->get_upf_local_teid();
  }
  for (auto itr = itp.first; itr != itp.second; itr++) {
    // Get the PDR numbers, now  get the rules from global static rule list
    GlobalRuleList.get_rule(itr->second, &rule);
    set_pdr_attributes(imsi, session, &rule);
    switch (rule.pdi().src_interface()) {
      case ACCESS:
        // Get new UPF TEID
        if (get_upf_teid) {
          upf_teid = insert_pdr_from_access(session, rule, session_uc);
        }
        break;
      case CORE:
        if (get_gnb_teid) {
          if (!insert_pdr_from_core(session, rule, session_uc)) {
            return 0;
          }
        }
        break;
    }
  }
  return upf_teid;
}

bool SessionStateEnforcer::handle_session_init_rule_updates(
    SessionMap& session_map, const std::string& imsi,
    std::unique_ptr<SessionState>& session_state) {
  uint32_t upf_teid = 0;
  bool get_gnb_teid = false;
  bool get_upf_teid = true;
  /* Add default rules to the IMSI set rules */
  add_default_rules(session_state, imsi);

  auto session_id = session_state->get_session_id();
  /* Irrespective of any State of Session, release and terminate*/
  SessionUpdate session_update =
      SessionStore::get_default_session_update(session_map);
  SessionStateUpdateCriteria& session_uc = session_update[imsi][session_id];
  const std::string upf_ip = get_upf_n3_addr();

  /* Attach rules to the session */
  upf_teid = update_session_rules(session_state, get_gnb_teid, get_upf_teid,
                                  &session_uc);
  /* session_state elments are filled with rules. State needs to be
   * moved to CREATING, increment version and send TEID details to AMF
   */
  if (!upf_teid) {
    MLOG(MERROR) << " Fail to get valide UPF teid end point"
                 << " for imsi " << imsi
                 << " session id: " << session_state->get_session_id();
    return false;
  }
  session_state->set_upf_teid_endpoint(upf_ip, upf_teid, &session_uc);
  session_state->set_fsm_state(CREATING, &session_uc);
  uint32_t cur_version = session_state->get_current_version();
  session_state->set_current_version(++cur_version, &session_uc);

  return true;
}

/* Function to initiate release of the session in enforcer requested by AMF
 * Go over session map and find the respective session of imsi
 * Go over SessionState vector and find the respective dnn (apn)
 * start terminating session process
 */
bool SessionStateEnforcer::m5g_release_session(SessionMap& session_map,
                                               const std::string& imsi,
                                               const uint32_t& pdu_id,
                                               SessionUpdate& session_update) {
  /* Search with session search criteria of IMSI and apn/dnn and
   * find  respective sesion to release operation
   * Note: DNN is optiona field, so find session from PDU_session_id
   */
  SessionSearchCriteria criteria(imsi, IMSI_AND_PDUID, pdu_id);
  auto session_it = session_store_.find_session(session_map, criteria);
  if (!session_it) {
    MLOG(MERROR) << "No session found in SessionMap for IMSI " << imsi
                 << " with pdu-id " << pdu_id << " to release";
    return false;
  }
  // Found the respective session to be updated
  auto& session = **session_it;
  auto session_id = session->get_session_id();
  /*Irrespective of any State of Session, release and terminate*/
  SessionStateUpdateCriteria& session_uc = session_update[imsi][session_id];
  MLOG(MDEBUG) << "Trying to release " << session_id << " from state "
               << session_fsm_state_to_str(session->get_state());

  m5g_start_session_termination(session_map, session, pdu_id, &session_uc);
  return true;
}

/*Start processing to terminate respective session requested from AMF*/
void SessionStateEnforcer::m5g_start_session_termination(
    SessionMap& session_map, const std::unique_ptr<SessionState>& session,
    const uint32_t& pdu_id, SessionStateUpdateCriteria* session_uc) {
  const auto session_id = session->get_session_id();
  const std::string& imsi = session->get_imsi();
  const auto previous_state = session->get_state();

  /* update respective session's state and return from here before timeout
   * to update session store with state and version
   */
  session->set_fsm_state(RELEASE, session_uc);
  uint32_t cur_version = session->get_current_version();
  session->set_current_version(++cur_version, session_uc);
  MLOG(MDEBUG) << "During release state of session changed to "
               << session_fsm_state_to_str(session->get_state());
  handle_state_update_to_amf(*session,
                             magma::lte::M5GSMCause::OPERATION_SUCCESS,
                             PDU_SESSION_STATE_NOTIFY);

  if (previous_state != CREATING) {
    /* Call for all rules to be de-associated from session
     * and inform to UPF
     */
    MLOG(MDEBUG) << "Will be removing all associated rules of session id "
                 << session->get_session_id();
    m5g_pdr_rules_change_and_update_upf(session, PdrState::REMOVE);
    if (session_map[imsi].size() == 0) {
      // delete the rules
      pdr_map_.erase(imsi);
    }
  }
  session->remove_all_rules_for_termination(session_uc);
  /* Forcefully terminate session context on time out
   * time out = 5000ms from sessiond.yml config file
   */
  MLOG(MDEBUG) << "Scheduling a force termination timeout for session_id "
               << session_id << " in " << session_force_termination_timeout_ms_
               << "ms";

  evb_->runAfterDelay(
      [this, imsi, session_id] {
        m5g_handle_termination_on_timeout(imsi, session_id);
      },
      session_force_termination_timeout_ms_);
}

/*Function to handle termination if UPF doesn't send required report
 * As per current implementation, upf report is not in place and
 * termination on time out will be executed forcefully
 */
void SessionStateEnforcer::m5g_handle_termination_on_timeout(
    const std::string& imsi, const std::string& session_id) {
  auto session_map = session_store_.read_sessions_for_deletion({imsi});
  auto session_update = SessionStore::get_default_session_update(session_map);
  bool marked_termination =
      session_update[imsi].find(session_id) != session_update[imsi].end();
  MLOG(MDEBUG) << "Forced termination timeout! Checking if termination has to "
               << "be forced for " << session_id << "... => "
               << (marked_termination ? "YES" : "NO");
  /* If the session doesn't exist in the session_update, then the session was
   * already released and terminated
   */
  if (marked_termination) {
    /*call to remove session from map*/
    m5g_complete_termination(session_map, imsi, session_id, session_update);

    bool update_success = session_store_.update_sessions(session_update);
    if (update_success) {
      MLOG(MDEBUG) << "Updated session termination of " << session_id
                   << " in to SessionStore";
    } else {
      MLOG(MERROR) << "Failed to update session termination of " << session_id
                   << " in to SessionStore";
    }
  } else {
    MLOG(MERROR) << "Nothing to remove as no respective entry found for "
                 << "session id " << session_id << " of IMSI " << imsi;
  }
}

/*Function will clean up all resources related to requested session
 * if it is last session entry, then delete the imsi
 * This function can be invoked from 2 different sources
 * 1. Time out and forcefully terminates session
 * 2. Once UPF sends report to SessionD
 * The 2nd one we are not taking care now.
 */
void SessionStateEnforcer::m5g_complete_termination(
    SessionMap& session_map, const std::string& imsi,
    const std::string& session_id, SessionUpdate& session_update) {
  // If the session cannot be found in session_map, or a new session has
  // already begun, do nothing.
  SessionSearchCriteria criteria(imsi, IMSI_AND_SESSION_ID, session_id);
  auto session_it = session_store_.find_session(session_map, criteria);
  if (!session_it) {
    // Session is already deleted, or new session already began, ignore.
    MLOG(MDEBUG) << "Could not find session for IMSI " << imsi
                 << " and session ID " << session_id
                 << ". Skipping termination.";
  }
  auto& session = **session_it;
  auto& session_uc = session_update[imsi][session_id];
  if (!session->can_complete_termination(&session_uc)) {
    return;  // error is logged in SessionState's complete_termination
  }
  auto termination_req = session->make_termination_request(&session_uc);
  auto logging_cb = SessionReporter::get_terminate_logging_cb(termination_req);
  reporter_->report_terminate_session(termination_req, logging_cb);
  events_reporter_->session_terminated(imsi, session);
  // Now remove all rules
  session->remove_all_rules(&session_uc);
  // Release and maintain TEID trakcing data structure TODO
  session_uc.is_session_ended = true;
  /*Removing session from map*/
  session_map[imsi].erase(*session_it);
  MLOG(MDEBUG) << session_id << " deleted from SessionMap";
  /* If it is last session terminated and no session left for this IMSI
   * remove the imsi as well
   */
  if (session_map[imsi].size() == 0) {
    session_map.erase(imsi);
    MLOG(MDEBUG) << "All sessions terminated for IMSI " << imsi;
  }
  MLOG(MDEBUG) << "Successfully terminated session " << session_id;
}

void SessionStateEnforcer::m5g_move_to_inactive_state(
    const std::string& imsi, std::unique_ptr<SessionState>& session,
    SetSmNotificationContext notif, SessionStateUpdateCriteria* session_uc) {
  set_new_fsm_state_and_increment_version(session, INACTIVE, session_uc);
  /* Call for all rules to be de-associated from session
   * and inform to UPF
   */
  m5g_pdr_rules_change_and_update_upf(session, PdrState::IDLE);
}

void SessionStateEnforcer::m5g_move_to_active_state(
    std::unique_ptr<SessionState>& session, SetSmNotificationContext notif,
    SessionStateUpdateCriteria* session_uc) {
  /* Reattach or get rules to the session */
  uint32_t upf_teid = update_session_rules(session, false, false, session_uc);
  /* As we got rules again, move the state to creating */
  session->set_fsm_state(CREATING, session_uc);
  uint32_t cur_version = session->get_current_version();
  session->set_current_version(++cur_version, session_uc);
  /* Send the UPF (local TEID) info to AMF which are going to
   * be used by GnodeB
   */
  prepare_response_to_access(*session,
                             magma::lte::M5GSMCause::OPERATION_SUCCESS,
                             get_upf_n3_addr(), upf_teid);
}

void SessionStateEnforcer::set_new_fsm_state_and_increment_version(
    std::unique_ptr<SessionState>& session, SessionFsmState target_state,
    SessionStateUpdateCriteria* session_uc) {
  auto stateStr = session_fsm_state_to_str(session->get_state());
  session->set_fsm_state(target_state, session_uc);
  uint32_t cur_version = session->get_current_version();
  session->set_current_version(++cur_version, session_uc);
  MLOG(MDEBUG) << "During " << stateStr << " state of session"
               << "of imsi: " << session->get_imsi()
               << "/teid:" << session->get_upf_local_teid() << "changed to "
               << session_fsm_state_to_str(session->get_state());
  return;
}

void SessionStateEnforcer::m5g_pdr_rules_change_and_update_upf(
    const std::unique_ptr<SessionState>& session, PdrState pdr_state) {
  // update criteria status not needed
  session->set_all_pdrs(pdr_state);
  session->reset_rtx_counter();
  RulesToProcess pending_activation;
  RulesToProcess pending_deactivation;
  const CreateSessionResponse& csr = session->get_create_session_response();

  std::vector<StaticRuleInstall> static_rule_installs =
      to_vec(csr.static_rules());
  std::vector<DynamicRuleInstall> dynamic_rule_installs =
      to_vec(csr.dynamic_rules());
  session->process_get_5g_rule_installs(
      static_rule_installs, dynamic_rule_installs, &pending_activation,
      &pending_deactivation);
  m5g_send_session_request_to_upf(session, pending_activation,
                                  pending_deactivation);
  return;
}

void SessionStateEnforcer::m5g_send_session_request_to_upf(
    const std::unique_ptr<SessionState>& session,
    const RulesToProcess& pending_activation,
    const RulesToProcess& pending_deactivation) {
  // Update to UPF
  SessionState::SessionInfo sess_info;
  session->sess_infocopy(&sess_info);
  // Set the node Id
  sess_info.nodeId.node_id = get_upf_node_id();
  pipelined_client_->set_upf_session(sess_info, pending_activation,
                                     pending_deactivation, call_back_upf);
  return;
}

/* This function will change the state of respective PDU session,
 * upon receving message or notification from UPF or due to
 * any other event or internal even/change causes state change,
 * then we further update state change to AMF module
 * imsi - from UPF handler to find respective SessionMap
 * seid - session context id to find respective session
 * new_state - state change required w.r.t. UPF message
 */
void SessionStateEnforcer::m5g_process_response_from_upf(
    const std::string& imsi, uint32_t teid, uint32_t version) {
  uint32_t cur_version;
  bool amf_update_pending = false;
  auto session_map = session_store_.read_sessions({imsi});
  /* Search with session search criteria of IMSI and session_id and
   * find  respective sesion to operate
   */
  SessionSearchCriteria criteria(imsi, IMSI_AND_TEID, teid);
  auto session_it = session_store_.find_session(session_map, criteria);
  if (!session_it) {
    MLOG(MERROR) << "No session found in SessionMap for IMSI " << imsi
                 << " with teid " << teid;
    return;
  }
  auto& session = **session_it;
  cur_version = session->get_current_version();
  auto session_update = SessionStore::get_default_session_update(session_map);
  auto session_id = session->get_session_id();
  SessionStateUpdateCriteria& session_uc = session_update[imsi][session_id];

  if (version < cur_version) {
    MLOG(MDEBUG) << "UPF verions of session imsi " << imsi << " session id "
                 << session->get_session_id() << " of  teid " << teid
                 << " recevied version " << version
                 << " SMF latest version: " << cur_version << " Resending";
    if (inc_rtx_counter(session)) {
      RulesToProcess pending_activation, pending_deactivation;
      RulesToSchedule pending_scheduling;
      const CreateSessionResponse& csr = session->get_create_session_response();
      std::vector<StaticRuleInstall> static_rule_installs =
          to_vec(csr.static_rules());
      std::vector<DynamicRuleInstall> dynamic_rule_installs =
          to_vec(csr.dynamic_rules());

      session->process_rules_to_install(
          static_rule_installs, dynamic_rule_installs, &pending_activation,
          &pending_deactivation, nullptr, &pending_scheduling, &session_uc);
      m5g_send_session_request_to_upf(session, pending_activation,
                                      pending_deactivation);
    }
    return;
  }
  switch (session->get_state()) {
    case CREATED:
      session->set_fsm_state(ACTIVE, &session_uc);
      /* As there is no config change, just state change as we
       * got response from UPF, so we dont bump up the session version
       * number here
       */
      amf_update_pending = true;
      break;
    case RELEASE:
      m5g_complete_termination(session_map, imsi, session_id, session_update);
    default:
      break;
  }
  if (amf_update_pending) {
    bool update_success = session_store_.update_sessions(session_update);
    if (update_success) {
      MLOG(MDEBUG) << "Updated SessionStore SessionState based on UPF message "
                   << " with session_id" << session->get_session_id();
    } else {
      MLOG(MERROR) << "Failed to update SessionStore based on UPF message"
                   << " with session_id" << session->get_session_id();
    }
    /* Update the state change notification to AMF */
    handle_state_update_to_amf(*session,
                               magma::lte::M5GSMCause::OPERATION_SUCCESS,
                               PDU_SESSION_STATE_NOTIFY);
  } else {
    session_store_.update_sessions(session_update);
  }
}

/* To prepare response back to AMF
 * Fill the response structure from session context message
 * and call rpc of AmfServiceClient.
 * TODO Recheck, if this can be part of AmfServiceClient and
 * move this function to AmfServiceClient object context.
 */
void SessionStateEnforcer::prepare_response_to_access(
    SessionState& session_state, const magma::lte::M5GSMCause m5gsm_cause,
    std::string upf_ip, uint32_t upf_teid) {
  magma::SetSMSessionContextAccess response;
  const auto& config = session_state.get_config();
  if (!config.rat_specific_context.has_m5gsm_session_context()) {
    MLOG(MWARNING) << "No M5G SM Session Context is specified for session";
    return;
  }

  RulesToProcess pending_activation, pending_deactivation;
  const CreateSessionResponse& csr =
      session_state.get_create_session_response();
  std::vector<StaticRuleInstall> static_rule_installs =
      to_vec(csr.static_rules());
  std::vector<DynamicRuleInstall> dynamic_rule_installs =
      to_vec(csr.dynamic_rules());

  session_state.process_get_5g_rule_installs(
      static_rule_installs, dynamic_rule_installs, &pending_activation,
      &pending_deactivation);

  /* Filing response proto message*/
  auto* rsp = response.mutable_rat_specific_context()
                  ->mutable_m5g_session_context_rsp();
  auto* rsp_cmn = response.mutable_common_context();

  rsp->set_pdu_session_id(
      config.rat_specific_context.m5gsm_session_context().pdu_session_id());
  rsp->set_pdu_session_type(
      config.rat_specific_context.m5gsm_session_context().pdu_session_type());
  rsp->set_selected_ssc_mode(
      config.rat_specific_context.m5gsm_session_context().ssc_mode());
  rsp->set_allowed_ssc_mode(
      config.rat_specific_context.m5gsm_session_context().ssc_mode());
  rsp->set_m5gsm_cause(m5gsm_cause);
  rsp->set_always_on_pdu_session_indication(
      config.rat_specific_context.m5gsm_session_context()
          .pdu_session_req_always_on());
  rsp->set_m5g_sm_congestion_reattempt_indicator(true);
  rsp->set_procedure_trans_identity(
      config.rat_specific_context.m5gsm_session_context()
          .procedure_trans_identity());

  /* AMBR value need to compared from AMF and PCF, then fill the required
   * values and sent to AMF.
   */
  // For now its default QOS, AMBR has default values
  rsp->mutable_session_ambr()->set_br_unit(
      config.rat_specific_context.m5gsm_session_context()
          .default_ambr()
          .br_unit());
  rsp->mutable_session_ambr()->set_max_bandwidth_ul(
      config.rat_specific_context.m5gsm_session_context()
          .default_ambr()
          .max_bandwidth_ul());
  rsp->mutable_session_ambr()->set_max_bandwidth_dl(
      config.rat_specific_context.m5gsm_session_context()
          .default_ambr()
          .max_bandwidth_ul());
  /* This flag is used for sending defult qos value or getting from policy
   *  value to AMF.
   */
  bool flag_set = false;
  for (auto& val : pending_activation) {
    if (val.rule.qos().max_req_bw_dl() && val.rule.qos().max_req_bw_ul()) {
      MLOG(MDEBUG) << "value set for pending_activation"
                   << val.rule.qos().max_req_bw_ul();
      rsp->mutable_session_ambr()->set_max_bandwidth_dl(
          val.rule.qos().max_req_bw_dl());
      rsp->mutable_session_ambr()->set_max_bandwidth_ul(
          val.rule.qos().max_req_bw_ul());
      rsp->mutable_qos()->CopyFrom(val.rule.qos());
      flag_set = true;
      break;
    }
  }
  if (!flag_set) {
    auto* convg_qos = rsp->mutable_qos();
    convg_qos->set_qci(FlowQos_Qci_QCI_9);
    convg_qos->mutable_arp()->set_pre_vulnerability(
        QosArp_PreVul_PRE_VUL_ENABLED);
    convg_qos->mutable_arp()->set_pre_capability(QosArp_PreCap_PRE_CAP_ENABLED);
    convg_qos->mutable_arp()->set_priority_level(1);
  }
  rsp->mutable_upf_endpoint()->set_teid(
      config.rat_specific_context.m5gsm_session_context()
          .upf_endpoint()
          .teid());

  rsp->mutable_upf_endpoint()->set_end_ipv4_addr(
      config.rat_specific_context.m5gsm_session_context()
          .upf_endpoint()
          .end_ipv4_addr());

  rsp_cmn->mutable_sid()->CopyFrom(config.common_context.sid());  // imsi
  if (!config.common_context.ue_ipv4().empty()) {
    rsp_cmn->set_ue_ipv4(config.common_context.ue_ipv4());
  }
  if (!config.common_context.ue_ipv6().empty()) {
    rsp_cmn->set_ue_ipv6(config.common_context.ue_ipv6());
  }
  rsp_cmn->set_apn(config.common_context.apn());
  rsp_cmn->set_sm_session_state(config.common_context.sm_session_state());
  rsp_cmn->set_sm_session_version(config.common_context.sm_session_version());
  // Send message to AMF gRPC client handler.
  amf_srv_client_->handle_response_to_access(response);
}

/* To update state change notification to AMF
 * Fill the notification structure from session context message
 * and call rpc of AmfServiceClient.
 * TODO Recheck, if this can be part of AmfServiceClient
 * move this function to AmfServiceClient object context.
 */
void SessionStateEnforcer::handle_state_update_to_amf(
    SessionState& session_state, const magma::lte::M5GSMCause m5gsm_cause,
    NotifyUeEvents event) {
  magma::SetSmNotificationContext notif;
  const auto& config = session_state.get_config();

  if (!config.rat_specific_context.has_m5gsm_session_context()) {
    MLOG(MWARNING) << "No M5G SM Session Context is specified for session";
    return;
  }
  auto* req = notif.mutable_rat_specific_notification();
  auto* req_cmn = notif.mutable_common_context();
  // Fill the imsi
  req_cmn->mutable_sid()->CopyFrom(config.common_context.sid());  // imsi
  req->set_notify_ue_event(event);
  // Fill the cause
  req->set_m5gsm_cause(m5gsm_cause);
  if (event == PDU_SESSION_STATE_NOTIFY) {
    req_cmn->set_sm_session_state(config.common_context.sm_session_state());
    req_cmn->set_sm_session_version(config.common_context.sm_session_version());
    req->set_pdu_session_id(
        config.rat_specific_context.m5gsm_session_context().pdu_session_id());
    req->set_pdu_session_type(
        config.rat_specific_context.m5gsm_session_context().pdu_session_type());
    req->set_request_type(EXISTING_PDU_SESSION);
    req->set_m5gsm_cause(m5gsm_cause);
  }
  // Send message to AMF gRPC client handler.
  amf_srv_client_->handle_notification_to_access(notif);
  return;
}

bool SessionStateEnforcer::default_and_static_rule_init() {
  // Static PDR, FAR, QDR, URR and BAR mapping  and also define 1 PDR and FAR
  SetGroupPDR reqpdr1;
  Action Value = FORW;
  uint32_t count = DEFAULT_PDR_ID;

  reqpdr1.set_pdr_id(++count);
  reqpdr1.set_precedence(DEFAULT_PDR_PRECEDENCE);
  reqpdr1.set_pdr_version(DEFAULT_PDR_VERSION);
  reqpdr1.set_pdr_state(PdrState::INSTALL);
  reqpdr1.mutable_pdi()->set_src_interface(ACCESS);

  reqpdr1.mutable_pdi()->set_net_instance("uplink");
  reqpdr1.set_o_h_remo_desc(0);
  reqpdr1.mutable_set_gr_far()->add_far_action_to_apply(Value);
  reqpdr1.mutable_activate_flow_req()->mutable_request_origin()->set_type(
      RequestOriginType_OriginType_N4);
  GlobalRuleList.insert_rule(DEFAULT_UP_LINK_PDR_ID, reqpdr1);
  // PDR 2 details
  SetGroupPDR reqpdr2;
  reqpdr2.set_pdr_id(++count);
  reqpdr2.set_precedence(DEFAULT_PDR_PRECEDENCE);
  reqpdr2.set_pdr_version(DEFAULT_PDR_VERSION);
  reqpdr2.set_pdr_state(PdrState::INSTALL);
  reqpdr2.mutable_pdi()->set_src_interface(CORE);
  reqpdr2.mutable_set_gr_far()->add_far_action_to_apply(Value);

  // Filling qos params
  reqpdr2.mutable_pdi()->set_net_instance("downlink");
  reqpdr2.mutable_activate_flow_req()->mutable_request_origin()->set_type(
      RequestOriginType_OriginType_N4);
  GlobalRuleList.insert_rule(DEFAULT_DOWN_LINK_PDR_ID, reqpdr2);

  return true;
}

uint32_t SessionStateEnforcer::get_next_teid() {
  /* For now TEID we use current no, increment for next, later we plan to
     maintain  release/alloc table for reu sing */
  uint32_t allocated_teid = teid_counter_;
  teid_counter_++;
  return allocated_teid;
}

bool SessionStateEnforcer::set_upf_node(const std::string& node_id,
                                        const std::string& addr) {
  upf_node_id_ = node_id;
  upf_node_ip_addr_ = addr;
  MLOG(MDEBUG) << "Set_upf_node_id: " << upf_node_id_;
  MLOG(MDEBUG) << "Set_upf_n3_addr: " << upf_node_ip_addr_;
  return true;
}

bool SessionStateEnforcer::is_incremented_rtx_counter_within_max(
    const std::unique_ptr<SessionState>& session) {
  uint32_t rtx_counter = session->get_incremented_rtx_counter();
  return rtx_counter < session_max_rtx_count_;
}

std::string SessionStateEnforcer::get_upf_n3_addr() const {
  return upf_node_ip_addr_;
}

std::string SessionStateEnforcer ::get_upf_node_id() const {
  return upf_node_id_;
}

/* Add defualt rules to the passed IMSI
 */
bool SessionStateEnforcer::add_default_rules(
    std::unique_ptr<SessionState>& session, const std::string& imsi) {
  // Check 2 default rules are added
  if (pdr_map_.count(imsi) == DEFAULT_RULE_COUNT) return true;
  /*
   * Lets add the default uplink and downlink PDR rule to this
   * imsi
   */
  if (!session->contains_pdr(DEFAULT_DOWN_LINK_PDR_ID)) {
    pdr_map_.insert(
        std::pair<std::string, uint32_t>(imsi, DEFAULT_DOWN_LINK_PDR_ID));
  }
  if (!session->contains_pdr(DEFAULT_UP_LINK_PDR_ID)) {
    pdr_map_.insert(
        std::pair<std::string, uint32_t>(imsi, DEFAULT_UP_LINK_PDR_ID));
  }
  /*
   * Plus policy DB added anymore it would be part of
   * this map
   */
  return true;
}

bool SessionStateEnforcer::insert_pdr_from_core(
    std::unique_ptr<SessionState>& session, SetGroupPDR& rule,
    SessionStateUpdateCriteria* session_uc) {
  const auto& config = session->get_config();
  uint32_t teid = 0;
  std::string ip_addr;

  // Get the latest session configuration
  teid = config.rat_specific_context.m5gsm_session_context()
             .gnode_endpoint()
             .teid();

  ip_addr = config.rat_specific_context.m5gsm_session_context()
                .gnode_endpoint()
                .end_ipv4_addr();
  if (!teid) {
    MLOG(MERROR) << " valid GNB endpoint details are not recv'd from AMF "
                 << " for imsi " << session->get_imsi()
                 << " session id: " << session->get_session_id()
                 << " gnodeb teid: " << teid;
  }
  rule.mutable_set_gr_far()
      ->mutable_fwd_parm()
      ->mutable_outr_head_cr()
      ->set_o_teid(teid);
  rule.mutable_set_gr_far()
      ->mutable_fwd_parm()
      ->mutable_outr_head_cr()
      ->set_gnb_ipv4_adr(ip_addr);
  rule.mutable_activate_flow_req()->set_downlink_tunnel(teid);
  rule.mutable_deactivate_flow_req()->set_downlink_tunnel(teid);

  MLOG(MINFO) << " AMF teid: " << teid << " ip address " << ip_addr
              << " of imsi: " << session->get_imsi()
              << " fteid: " << session->get_upf_local_teid()
              << " pdu id: " << session->get_pdu_id() << " UE ip address "
              << rule.pdi().ue_ipv4();
  // Insert the PDR along with teid into the session
  session->insert_pdr(&rule, session_uc);
  return true;
}

uint32_t SessionStateEnforcer::insert_pdr_from_access(
    std::unique_ptr<SessionState>& session, SetGroupPDR& rule,
    SessionStateUpdateCriteria* session_uc) {
  uint32_t upf_teid = get_next_teid();
  MLOG(MDEBUG) << "Acquried Teid: " << upf_teid;
  rule.mutable_pdi()->set_local_f_teid(upf_teid);
  rule.mutable_activate_flow_req()->set_uplink_tunnel(upf_teid);
  rule.mutable_deactivate_flow_req()->set_uplink_tunnel(upf_teid);
  // Insert the PDR along with teid into the session
  session->insert_pdr(&rule, session_uc);
  return upf_teid;
}

uint32_t SessionStateEnforcer::get_current_teid() { return teid_counter_; }

bool SessionStateEnforcer::inc_rtx_counter(
    const std::unique_ptr<SessionState>& session) {
  uint32_t rtx_counter = session->get_incremented_rtx_counter();
  return rtx_counter < session_max_rtx_count_;
}

void SessionStateEnforcer::set_pdr_attributes(
    const std::string& imsi, std::unique_ptr<SessionState>& session_state,
    SetGroupPDR* rule) {
  const auto& config = session_state->get_config();
  auto ue_ipv4 = config.common_context.ue_ipv4();
  auto ue_ipv6 = config.common_context.ue_ipv6();

  rule->mutable_pdi()->set_ue_ipv4(ue_ipv4);
  rule->mutable_pdi()->set_ue_ipv6(ue_ipv6);
  rule->mutable_activate_flow_req()->mutable_sid()->set_id(imsi);
  rule->mutable_deactivate_flow_req()->mutable_sid()->set_id(imsi);
  rule->mutable_activate_flow_req()->set_ip_addr(
      config.common_context.ue_ipv4());
  rule->mutable_activate_flow_req()->set_ipv6_addr(
      config.common_context.ue_ipv6());
  rule->mutable_deactivate_flow_req()->set_ip_addr(
      config.common_context.ue_ipv4());
  rule->mutable_deactivate_flow_req()->set_ipv6_addr(
      config.common_context.ue_ipv6());
}

std::vector<StaticRuleInstall> SessionStateEnforcer::to_vec(
    const google::protobuf::RepeatedPtrField<magma::lte::StaticRuleInstall>
        static_rule_installs) {
  std::vector<StaticRuleInstall> out;
  for (const auto& install : static_rule_installs) {
    out.push_back(install);
  }
  return out;
}

std::vector<DynamicRuleInstall> SessionStateEnforcer::to_vec(
    const google::protobuf::RepeatedPtrField<magma::lte::DynamicRuleInstall>
        dynamic_rule_installs) {
  std::vector<DynamicRuleInstall> out;
  for (const auto& install : dynamic_rule_installs) {
    out.push_back(install);
  }
  return out;
}

void SessionStateEnforcer::update_session_with_policy(
    std::unique_ptr<SessionState>& session,
    const CreateSessionResponse& response,
    SessionStateUpdateCriteria* session_uc) {
  session->set_tgpp_context(response.tgpp_ctx(), session_uc);
  session->set_create_session_response(response, session_uc);

  prepare_response_to_access(*session,
                             magma::lte::M5GSMCause::OPERATION_SUCCESS,
                             get_upf_n3_addr(), session->get_upf_local_teid());
}

}  // end namespace magma
