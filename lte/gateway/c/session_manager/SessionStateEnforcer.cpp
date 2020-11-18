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

#include <string>
#include <time.h>
#include <utility>
#include <vector>

#include <google/protobuf/repeated_field.h>
#include <google/protobuf/timestamp.pb.h>
#include <google/protobuf/util/time_util.h>
#include <grpcpp/channel.h>
#include "magma_logging.h"
#include "EnumToString.h"
#include "SessionStateEnforcer.h"

std::shared_ptr<magma::SessionStateEnforcer> conv_session_enforcer;
namespace magma {
void call_back_upf(grpc::Status, magma::UPFSessionContextState response) {
  std::string imsi             = response.session_snapshot().subscriber_id();
  uint32_t version             = response.session_snapshot().session_version();
  uint32_t fteid               = response.session_snapshot().local_f_teid();
  const std::string session_id = response.session_snapshot().subscriber_id();
  MLOG(MINFO) << " Async Response received from UPF: imsi " << imsi
              << " local fteid : " << fteid;
  conv_session_enforcer->get_event_base().runInEventBaseThread(
      [imsi, fteid, version]() {
        /* Update the state change, and notifiy to AMF */
        // For now fteid will be zero in all cases
        conv_session_enforcer->m5g_update_session_state_to_amf(
            imsi, fteid, version);
      });
}

/*constructor*/
SessionStateEnforcer::SessionStateEnforcer(
    std::shared_ptr<StaticRuleStore> rule_store, SessionStore& session_store,
    std::shared_ptr<PipelinedClient> pipelined_client,
    std::shared_ptr<AmfServiceClient> amf_srv_client,
    magma::mconfig::SessionD mconfig, long session_force_termination_timeout_ms)
    : session_store_(session_store),
      pipelined_client_(pipelined_client),
      amf_srv_client_(amf_srv_client),
      retry_timeout_(1),
      mconfig_(mconfig),
      session_force_termination_timeout_ms_(
          session_force_termination_timeout_ms) {
  // for now this is the right place, need to move if find  anohter right place
  static_rule_init();
}

void SessionStateEnforcer::attachEventBase(folly::EventBase* evb) {
  evb_ = evb;
}

void SessionStateEnforcer::stop() {
  evb_->terminateLoopSoon();
}

folly::EventBase& SessionStateEnforcer::get_event_base() {
  return *evb_;
}

bool SessionStateEnforcer::m5g_init_session_credit(
    SessionMap& session_map, const std::string& imsi,
    const std::string& session_id, const SessionConfig& cfg) {
  /* creating SessionState object with state CREATING
   * This calls constructor and allocates memory*/
  auto session_state =
      std::make_unique<SessionState>(imsi, session_id, cfg, *rule_store_);
  MLOG(MINFO) << " New SessionState object created with IMSI: " << imsi
              << " session context id : " << session_id;
  handle_session_init_rule_updates(session_map, imsi, *session_state);

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
  MLOG(MINFO) << "Added a session (" << session_map[imsi].size()
              << ") for IMSI " << imsi << " with session context ID ";
  return true;
}

void SessionStateEnforcer::handle_session_init_rule_updates(
    SessionMap& session_map, const std::string& imsi,
    SessionState& session_state) {
  uint32_t l_teid = 0;
  auto itp        = pdr_map_.equal_range(imsi);
  for (auto itr = itp.first; itr != itp.second; itr++) {
    // Get the PDR numbers, now  get the rules from global static rule list
    SetGroupPDR rule;
    GlobalRuleList.get_rule(itr->second, &rule);
    session_state.insert_pdr(&rule);
    // Get new TEID
    if (!l_teid) l_teid = pipelined_client_->get_current_teid();
    rule.mutable_pdi()->set_local_f_teid(pipelined_client_->get_next_teid());
  }
  auto ip_addr = session_state.get_config()
                     .rat_specific_context.m5gsm_session_context()
                     .pdu_address()
                     .redirect_server_address();
  /* session_state elments are filled with rules. State needs to be
   * moved to CREATED, increment version and send message to UPF.
   * Note: charging and credit related info not taken care in drop-1
   *
   * TODO - will be taken care later
   * Thinking of adding the created session into session store after step
   * CREATING state and before rules are add to session. And adding the logic
   * of implementing policy rules into the event loop. This way, when PCF
   * component is introduced an asynchronous call into PCF (PolicyDB) can be
   * handled easily later on.
   *
   */
  auto update_criteria = get_default_update_criteria();
  session_state.set_local_teid(l_teid, update_criteria);
  session_state.set_fsm_state(CREATED, update_criteria);
  MLOG(MINFO) << " Teid " << session_state.get_local_teid();
  uint32_t cur_version = session_state.get_current_version();
  session_state.set_current_version(++cur_version, update_criteria);

  /* Update the m5gsm_cause and prepare for respone along with actual cause*/
  prepare_response_to_access(
      imsi, session_state, magma::lte::M5GSMCause::OPERATION_SUCCESS);

  SessionState::SessionInfo sess_info;
  sess_info.imsi         = imsi;
  sess_info.ip_addr      = ip_addr;
  sess_info.local_f_teid = session_state.get_local_teid();

  sess_info.Pdr_rules_ = session_state.get_all_pdr_rules();
  // sess_info.Far_rules_ = session_state.get_all_far_rules();
  session_state.sess_infocopy(&sess_info);
  pipelined_client_->set_upf_session(sess_info, call_back_upf);
  return;
}

/* Function to initiate release of the session in enforcer requested by AMF
 * Go over session map and find the respective session of imsi
 * Go over SessionState vector and find the respective dnn (apn)
 * start terminating session process
 */
bool SessionStateEnforcer::m5g_release_session(
    SessionMap& session_map, const std::string& imsi, const std::string& dnn,
    SessionUpdate& session_update) {
  /* Search with session search criteria of IMSI and apn/dnn and
   * find  respective sesion to release operation
   */
  SessionSearchCriteria criteria(imsi, IMSI_AND_APN, dnn);
  auto session_it = session_store_.find_session(session_map, criteria);
  if (!session_it) {
    MLOG(MERROR) << "No session found in SessionMap for IMSI " << imsi
                 << " with DNN " << dnn << " to release";
    return false;
  }
  // Found the respective session to be updated
  auto& session   = **session_it;
  auto session_id = session->get_session_id();
  /*Irrespective of any State of Session, release and terminate*/
  SessionStateUpdateCriteria& session_uc = session_update[imsi][session_id];
  MLOG(MINFO) << "Trying to release session with id " << session_id
              << " from state "
              << session_fsm_state_to_str(session->get_state());
  m5g_start_session_termination(imsi, session, dnn, session_uc);
  return true;
}

/*Start processing to terminate respective session requested from AMF*/
void SessionStateEnforcer::m5g_start_session_termination(
    const std::string& imsi, const std::unique_ptr<SessionState>& session,
    const std::string& dnn, SessionStateUpdateCriteria& session_uc) {
  const auto session_id = session->get_session_id();

  /* update respective session's state and return from here before timeout
   * to update session store with state and version
   */
  session->set_fsm_state(RELEASE, session_uc);
  uint32_t cur_version = session->get_current_version();
  session->set_current_version(++cur_version, session_uc);
  MLOG(MINFO) << "During release state of session changed to "
              << session_fsm_state_to_str(session->get_state());
  handle_state_update_to_amf(
      imsi, *session, magma::lte::M5GSMCause::OPERATION_SUCCESS);

  /* Call for all rules to be de-associated from session
   * and inform to UPF
   */
  m5g_remove_rules_and_update_upf(imsi, session, dnn, session_uc);

  /* Forcefully terminate session context on time out
   * time out = 5000ms from sessiond.yml config file
   */
  MLOG(MINFO) << "Scheduling a force termination timeout for session_id "
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
  auto session_map    = session_store_.read_sessions_for_deletion({imsi});
  auto session_update = SessionStore::get_default_session_update(session_map);
  bool marked_termination =
      session_update[imsi].find(session_id) != session_update[imsi].end();
  MLOG(MINFO) << "Forced termination timeout! Checking if termination has to "
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
      MLOG(MINFO) << "Updated session termination of " << session_id
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
  auto& session    = **session_it;
  auto& session_uc = session_update[imsi][session_id];
  if (!session->can_complete_termination(session_uc)) {
    return;  // error is logged in SessionState's complete_termination
  }
  // Now remove all rules
  session->remove_all_rules();
  // Release and maintain TEID trakcing data structure TODO
  session_uc.is_session_ended = true;
  /*Removing session from map*/
  session_map[imsi].erase(*session_it);
  MLOG(MINFO) << session_id << " deleted from SessionMap";
  /* If it is last session terminated and no session left for this IMSI
   * remove the imsi as well
   */
  if (session_map[imsi].size() == 0) {
    session_map.erase(imsi);
    MLOG(MINFO) << "All sessions terminated for IMSI " << imsi;
  }
  MLOG(MINFO) << "Successfully terminated session " << session_id;
  return;
}

/*Removing all associated rules to session*/
void SessionStateEnforcer::m5g_remove_rules_and_update_upf(
    const std::string& imsi, const std::unique_ptr<SessionState>& session,
    const std::string& dnn, SessionStateUpdateCriteria& uc) {
  // remove all rules;
  MLOG(MINFO) << "Will be removing all associated rules of session id "
              << session->get_session_id();
  auto ip_addr = session->get_config()
                     .rat_specific_context.m5gsm_session_context()
                     .pdu_address()
                     .redirect_server_address();
  // Update to UPF
  SessionState::SessionInfo sess_info;
  sess_info.imsi         = imsi;
  sess_info.ip_addr      = ip_addr;
  sess_info.local_f_teid = session->get_local_teid();
  session->sess_infocopy(&sess_info);
  // Set PDR state as  REMOVE for all PDRs
  session->set_remove_all_pdrs();
  sess_info.Pdr_rules_ = session->get_all_pdr_rules();
  pipelined_client_->set_upf_session(sess_info, call_back_upf);
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
void SessionStateEnforcer::m5g_update_session_state_to_amf(
    const std::string& imsi, uint32_t teid, uint32_t version) {
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
  bool amf_update_pending = false;
  uint32_t cur_version;
  auto& session       = **session_it;
  auto session_update = SessionStore::get_default_session_update(session_map);
  SessionStateUpdateCriteria& session_uc =
      session_update[imsi][session->get_session_id()];
  switch (session->get_state()) {
    case CREATED:
      session->set_fsm_state(ACTIVE, session_uc);
      cur_version = session->get_current_version();
      session->set_current_version(++cur_version, session_uc);
      amf_update_pending = true;
      break;
    case RELEASE:
      m5g_complete_termination(
          session_map, imsi, session->get_session_id(), session_update);
    default:
      break;
  }
  if (amf_update_pending) {
    bool update_success = session_store_.update_sessions(session_update);
    if (update_success) {
      MLOG(MINFO) << "Updated SessionStore SessionState based on UPF message "
                  << " with session_id" << session->get_session_id();
    } else {
      MLOG(MERROR) << "Failed to update SessionStore based on UPF message"
                   << " with session_id" << session->get_session_id();
    }
    /* Update the state change notification to AMF */
    handle_state_update_to_amf(
        imsi, *session, magma::lte::M5GSMCause::OPERATION_SUCCESS);
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
    const std::string& imsi, SessionState& session_state,
    const magma::lte::M5GSMCause m5gsm_cause) {
  magma::SetSMSessionContextAccess response;
  const auto& config = session_state.get_config();

  if (!config.rat_specific_context.has_m5gsm_session_context()) {
    MLOG(MWARNING) << "No M5G SM Session Context is specified for session";
    return;
  }

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
  rsp->set_m5gsm_congetion_re_attempt_indicator(true);
  rsp->mutable_pdu_address()->set_redirect_address_type(
      config.rat_specific_context.m5gsm_session_context()
          .pdu_address()
          .redirect_address_type());
  rsp->mutable_pdu_address()->set_redirect_server_address(
      config.rat_specific_context.m5gsm_session_context()
          .pdu_address()
          .redirect_server_address());
  /* TODO AMBR and QoS will be implemented later.
   * AMBR value need to compared from AMF and PCF, then fill the required
   * values and sent to AMF.
   */

  rsp_cmn->mutable_sid()->CopyFrom(config.common_context.sid());  // imsi
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
    const std::string& imsi, SessionState& session_state,
    const magma::lte::M5GSMCause m5gsm_cause) {
  magma::SetSmNotificationContext notif;
  const auto& config = session_state.get_config();

  if (!config.rat_specific_context.has_m5gsm_session_context()) {
    MLOG(MWARNING) << "No M5G SM Session Context is specified for session";
    return;
  }
  auto* req     = notif.mutable_rat_specific_notification();
  auto* req_cmn = notif.mutable_common_context();

  req_cmn->mutable_sid()->CopyFrom(config.common_context.sid());  // imsi
  req_cmn->set_sm_session_state(config.common_context.sm_session_state());
  req_cmn->set_sm_session_version(config.common_context.sm_session_version());
  req->set_pdu_session_id(
      config.rat_specific_context.m5gsm_session_context().pdu_session_id());
  req->set_pdu_session_type(
      config.rat_specific_context.m5gsm_session_context().pdu_session_type());
  req->set_access_type(
      config.rat_specific_context.m5gsm_session_context().access_type());
  req->set_request_type(EXISTING_PDU_SESSION);
  req->set_m5gsm_cause(m5gsm_cause);
  // Send message to AMF gRPC client handler.
  amf_srv_client_->handle_notification_to_access(notif);
}
bool SessionStateEnforcer::static_rule_init() {
  // Static PDR, FAR, QDR, URR and BAR mapping  and also define 1 PDR and FAR
  SetGroupPDR reqpdr1;
  Action Value   = FORW;
  uint32_t count = 0;

  reqpdr1.set_pdr_id(++count);
  reqpdr1.set_precedence(32);
  reqpdr1.set_pdr_version(1);
  reqpdr1.set_pdr_state(PdrState::INSTALL);

  reqpdr1.mutable_pdi()->set_ue_ip_adr("192.168.128.11");
  reqpdr1.mutable_pdi()->set_net_instance("uplink");
  reqpdr1.set_o_h_remo_desc(0);
  reqpdr1.mutable_set_gr_far()->add_far_action_to_apply(Value);
  reqpdr1.mutable_activate_flow_req()->mutable_request_origin()->set_type(
      RequestOriginType_OriginType_N4);
  // Add rule 1
  PolicyRule rule1;
  rule1.set_id("rule1");
  rule1.set_priority(10);
  FlowDescription fd1;
  fd1.mutable_match()->set_ipv4_dst("192.168.0.0/16");
  fd1.mutable_match()->set_direction(FlowMatch_Direction_UPLINK);
  fd1.set_action(FlowDescription_Action_PERMIT);
  rule1.mutable_flow_list()->Add()->CopyFrom(fd1);
  reqpdr1.mutable_activate_flow_req()->mutable_dynamic_rules()->Add()->CopyFrom(
      rule1);
  GlobalRuleList.insert_rule(1, reqpdr1);
  // PDR 2 details
  SetGroupPDR reqpdr2;
  reqpdr2.set_pdr_id(++count);
  reqpdr2.set_precedence(32);
  reqpdr2.set_pdr_version(1);
  reqpdr2.set_pdr_state(PdrState::INSTALL);
  reqpdr2.mutable_pdi()->set_src_interface(1);
  reqpdr2.mutable_pdi()->set_ue_ip_adr("192.168.128.11");
  reqpdr2.mutable_set_gr_far()->add_far_action_to_apply(Value);

  reqpdr2.mutable_set_gr_far()
      ->mutable_fwd_parm()
      ->mutable_outr_head_cr()
      ->set_o_teid(200);
  reqpdr2.mutable_set_gr_far()
      ->mutable_fwd_parm()
      ->mutable_outr_head_cr()
      ->set_gnb_ipv4_adr("192.168.60.141");
  // Filling qos params
  reqpdr2.mutable_pdi()->set_net_instance("downlink");
  reqpdr2.mutable_activate_flow_req()->mutable_request_origin()->set_type(
      RequestOriginType_OriginType_N4);
  // For rule 1 change the driection alone
  PolicyRule rule2;
  rule2.set_id("rule2");
  rule2.set_priority(10);
  FlowDescription fd2;
  fd2.mutable_match()->set_ipv4_src("192.168.0.0/16");
  fd2.mutable_match()->set_direction(FlowMatch_Direction_DOWNLINK);
  fd2.set_action(FlowDescription_Action_PERMIT);
  rule2.mutable_flow_list()->Add()->CopyFrom(fd2);
  reqpdr2.mutable_activate_flow_req()->mutable_dynamic_rules()->Add()->CopyFrom(
      rule2);
  GlobalRuleList.insert_rule(2, reqpdr2);

  // subscriber Id 1 to PDR 1 and FAR 1
  pdr_map_.insert(std::pair<std::string, uint32_t>("imsi00000000001", 1));
  pdr_map_.insert(std::pair<std::string, uint32_t>("imsi00000000001", 2));
  pdr_map_.insert(std::pair<std::string, uint32_t>("IMSI001222333", 1));
  pdr_map_.insert(std::pair<std::string, uint32_t>("IMSI001222333", 2));
  reqpdr1.set_pdr_id(++count);
  reqpdr1.mutable_pdi()->set_ue_ip_adr("192.168.128.12");
  GlobalRuleList.insert_rule(3, reqpdr1);
  reqpdr2.set_pdr_id(++count);
  reqpdr2.mutable_pdi()->set_ue_ip_adr("192.168.128.12");
  reqpdr2.mutable_set_gr_far()
      ->mutable_fwd_parm()
      ->mutable_outr_head_cr()
      ->set_o_teid(200 + count);
  GlobalRuleList.insert_rule(4, reqpdr2);
  pdr_map_.insert(std::pair<std::string, uint32_t>("imsi00000000002", 3));
  pdr_map_.insert(std::pair<std::string, uint32_t>("imsi00000000002", 4));
  reqpdr1.set_pdr_id(++count);
  reqpdr1.mutable_pdi()->set_ue_ip_adr("192.168.128.13");
  GlobalRuleList.insert_rule(5, reqpdr1);
  reqpdr2.set_pdr_id(++count);
  reqpdr2.mutable_pdi()->set_ue_ip_adr("192.168.128.13");
  reqpdr2.mutable_set_gr_far()
      ->mutable_fwd_parm()
      ->mutable_outr_head_cr()
      ->set_o_teid(200 + count);
  GlobalRuleList.insert_rule(6, reqpdr2);
  pdr_map_.insert(std::pair<std::string, uint32_t>("imsi00000000003", 5));
  pdr_map_.insert(std::pair<std::string, uint32_t>("imsi00000000003", 6));
  reqpdr1.set_pdr_id(++count);
  reqpdr1.mutable_pdi()->set_ue_ip_adr("192.168.128.14");
  GlobalRuleList.insert_rule(7, reqpdr1);
  reqpdr2.set_pdr_id(++count);
  reqpdr2.mutable_pdi()->set_ue_ip_adr("192.168.128.14");
  reqpdr2.mutable_set_gr_far()
      ->mutable_fwd_parm()
      ->mutable_outr_head_cr()
      ->set_o_teid(200 + count);
  GlobalRuleList.insert_rule(8, reqpdr2);
  return true;
}
}  // end namespace magma
