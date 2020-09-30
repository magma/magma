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

namespace magma {
// temp routine
void call_back_void_upf(grpc::Status, magma::UpfRes response) {
  // do nothinf but to only passing call back
}

/*constructor*/
SessionStateEnforcer::SessionStateEnforcer(
    std::shared_ptr<StaticRuleStore> rule_store, SessionStore& session_store,
    std::shared_ptr<PipelinedClient> pipelined_client,
    std::shared_ptr<AmfServiceClient> amf_srv_client,
    magma::mconfig::SessionD mconfig)
    : session_store_(session_store),
      pipelined_client_(pipelined_client),
      amf_srv_client_(amf_srv_client),
      retry_timeout_(1),
      mconfig_(mconfig) {
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
  MLOG(MDEBUG) << "Added a session (" << session_map[imsi].size()
               << ") for IMSI " << imsi << " with session context ID "
               << session_id;
  return true;
}

void SessionStateEnforcer::handle_session_init_rule_updates(
    SessionMap& session_map, const std::string& imsi,
    SessionState& session_state) {
  auto itp = pdr_map_.equal_range(imsi);
  for (auto itr = itp.first; itr != itp.second; itr++) {
    // Get the PDR numbers, now  get the rules from global static rule list
    SetGroupPDR rule;
    GlobalRuleList.get_rule(itr->second, &rule);
    session_state.insert_pdr(&rule);
  }
  auto itf = far_map_.equal_range(imsi);
  for (auto itr = itf.first; itr != itf.second; itr++) {
    // Get the PDR number, from global static rule list
    SetGroupFAR rule;
    GlobalRuleList.get_rule(itr->second, &rule);
    // Add to the the session vector
    session_state.insert_far(&rule);
  }
  auto ip_addr = session_state.get_config()
                     .rat_specific_context.m5gsm_session_context()
                     .pdu_address()
                     .redirect_server_address();
  SessionState::SessionInfo sess_info;
  sess_info.imsi       = imsi;
  sess_info.ip_addr    = ip_addr;
  sess_info.Pdr_rules_ = session_state.get_all_pdr_rules();
  sess_info.Far_rules_ = session_state.get_all_far_rules();
  session_state.sess_infocopy(&sess_info);

  /* session_state elments are filled with rules. State needs to be
   * moved to CREATED and sending message to UPF.
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
  session_state.set_fsm_state(CREATED, update_criteria);
  MLOG(MDEBUG) << "State of session changed to "
               << session_fsm_state_to_str(session_state.get_state());
  /* Update the m5gsm_cause and prepare for respone along with actual cause*/
  prepare_response_to_access(
      imsi, session_state, magma::lte::M5GSMCause::OPERATION_SUCCESS);
  pipelined_client_->set_upf_session(sess_info, call_back_void_upf);
  return;
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

bool SessionStateEnforcer::static_rule_init() {
  // Static PDR, FAR, QDR, URR and BAR mapping  and also define 1 PDR and FAR
  SetGroupPDR reqpdr1, reqpdr2;
  SetGroupFAR reqf;
  magma::PDI pdireq;
  reqpdr1.set_pdr_id(1);
  reqpdr1.set_sess_ver_no(2);
  reqpdr1.set_precedence(5);
  reqpdr1.set_far_id(1);
  pdireq.set_src_interface(3);
  pdireq.set_net_instance("downlink");
  pdireq.set_ue_ip_adr("10.10.1.1");
  GlobalRuleList.insert_rule(1, reqpdr1);
  reqpdr2.set_pdr_id(2);
  reqpdr2.set_sess_ver_no(1);
  reqpdr2.set_precedence(2);
  reqpdr2.set_far_id(2);
  pdireq.set_src_interface(3);
  pdireq.set_net_instance("uplink");
  pdireq.set_ue_ip_adr("10.10.1.2");
  GlobalRuleList.insert_rule(2, reqpdr2);
  SetGroupFAR far1;
  far1.set_far_id(1);
  far1.set_sess_ver_no(4);
  far1.set_bar_id(6);
  GlobalRuleList.insert_rule(1, far1);

  // subscriber Id 1 to PDR 1 and FAR 1
  pdr_map_.insert(std::pair<std::string, uint32_t>("IMSI00000001", 1));
  far_map_.insert(std::pair<std::string, uint32_t>("IMSI00000001", 1));
  // subscriber Id 2  PDR list shud 2 and 4, far list also 2 and 4
  pdr_map_.insert(std::pair<std::string, uint32_t>("IMSI00000002", 2));
  far_map_.insert(std::pair<std::string, uint32_t>("IMSI00000002", 2));
  return true;
}
}  // end namespace magma
