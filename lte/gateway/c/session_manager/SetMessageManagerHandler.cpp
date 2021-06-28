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
  Source      	SetMessageManagerHandler.cpp
  Version     	0.1
  Date       	2020/08/08
  Product     	SessionD
  Subsystem   	5G Landing object in SessionD
  Author/Editor Sanjay Kumar Ojha
  Description 	Acts as 5G Landing object in SessionD & start 5G related flow
*****************************************************************************/
#include <chrono>
#include <thread>

#include <google/protobuf/util/time_util.h>

#include "SetMessageManagerHandler.h"
#include "magma_logging.h"
#include "GrpcMagmaUtils.h"

using grpc::Status;

namespace magma {

/*
 * SetMessageManagerHandler processes gRPC requests for the sessionD
 * This composites the earlier LocalSessionManagerHandlerImpl and uses the
 * exiting functionalities.
 */

/*Constructor*/
SetMessageManagerHandler::SetMessageManagerHandler(
    std::shared_ptr<SessionStateEnforcer> m5genforcer,
    SessionStore& session_store)
    : session_store_(session_store), m5g_enforcer_(m5genforcer) {}

/* Building session config with required parameters
 * TODO Note: this function can be removed by implementing 5G specific
 * SeesionConfig constructor
 */
SessionConfig SetMessageManagerHandler::m5g_build_session_config(
    const SetSMSessionContext& request) {
  SessionConfig cfg;
  /*copying only 5G specific data to respective elements*/
  cfg.common_context       = request.common_context();
  cfg.rat_specific_context = request.rat_specific_context();

  /* As we dont have 5G polices defined yet, for now
   * for all new connection we set SSC  mode as SSC_MODE_3
   */
  cfg.rat_specific_context.mutable_m5gsm_session_context()->set_ssc_mode(
      SSC_MODE_3);

  return cfg;
}

/* Handling set message from AMF
 * check if it is INITIAL_REQUEST or EXISTING_PDU_SESSION
 * if it is INITIAL_REQUEST need to create the session context in SessionMap
 * and write to memory by SessionStore
 * If EXISTING_PDU_SESSION, check on incoming session state and version
 * accordingly take action.
 * As per ASYNC method, response_callback is included but not functional
 */

void SetMessageManagerHandler::SetAmfSessionContext(
    ServerContext* context, const SetSMSessionContext* request,
    std::function<void(Status, SmContextVoid)> response_callback) {
  auto& request_cpy = *request;
  // Print the message from AMF
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(request_cpy));

  // Requested message from AMF to release the session
  m5g_enforcer_->get_event_base().runInEventBaseThread([this, response_callback,
                                                        request_cpy]() {
    // extract values from proto
    auto imsi       = request_cpy.common_context().sid().id();
    std::string dnn = request_cpy.common_context().apn();
    // pdu_id is unique to IMSI
    auto pdu_id = request_cpy.rat_specific_context()
                      .m5gsm_session_context()
                      .pdu_session_id();

    // Fetch complete message from proto message
    SessionConfig cfg = m5g_build_session_config(request_cpy);

    /* Read the proto message and check for state. Get the config out of proto.
     * Code for relase state, then creating
     */
    // Requested message from AMF to release the session
    if (cfg.common_context.sm_session_state() == RELEASED_4) {
      if (cfg.common_context.sm_session_version() != 0) {
        MLOG(MERROR) << "Wrong version received from AMF for IMSI " << imsi
                     << " but continuing release request";
      }
      MLOG(MINFO) << "Release request for session from IMSI: " << imsi
                  << " DNN " << dnn;
      /* Read the SessionMap from global session_store,
       * if it is not found, it will be added w.r.t imsi
       */
      auto session_map = session_store_.read_sessions({imsi});
      initiate_release_session(session_map, dnn, imsi);
    } else {
      // The Event Based main_thread invocation and runs to handle session state
      MLOG(MDEBUG) << "Requested session from UE with IMSI: " << imsi
                   << " PDU ID " << pdu_id;

      /* Message may be intial or modification message. Only taken care
       * intial message. Check if it's initial message
       */
      if ((cfg.rat_specific_context.m5gsm_session_context().rquest_type() ==
           INITIAL_REQUEST) &&
          (cfg.common_context.sm_session_state() == CREATING_0)) {
        /* it's new UE establisment request and need to create the session
         * context
         */
        MLOG(MINFO)
            << "AMF request type INITIAL_REQUEST and session state CREATING";
        auto session_map = session_store_.read_sessions({imsi});
        send_create_session(session_map, imsi, cfg, pdu_id);
        response_callback(Status::OK, SmContextVoid());
        return;
      }
      MLOG(MERROR)
          << "AMF request type- Unhandled request type:"
          << cfg.rat_specific_context.m5gsm_session_context().rquest_type();
      Status status(grpc::UNKNOWN, "Unknown session state or request");
      response_callback(status, SmContextVoid());
    }
    return;
  });
}

/* Creeate respective SessionState and context*/
void SetMessageManagerHandler::send_create_session(
    SessionMap& session_map, const std::string& imsi, SessionConfig& cfg_new,
    uint32_t& pdu_id) {
  /* If it is new session to be created, check for same DNN exists
   * for same IMSI, i.e if IMSI found and respective DNN found in
   * SessionStore, then return from here and nothing to do
   * as already same session exist, its duplicate request
   */

  SessionSearchCriteria criteria(imsi, IMSI_AND_PDUID, pdu_id);
  auto session_it = session_store_.find_session(session_map, criteria);
  if (session_it) {
    auto& session = **session_it;
    /* check if session state is in "CREATING" state
     * then update the state, otherwise fire an error
     */
    if (session->get_state() == CREATING) {
      // Get the GNODEB teid and IP address from config
      SessionConfig cfg = session->get_config();
      cfg.rat_specific_context.mutable_m5gsm_session_context()
          ->mutable_gnode_endpoint()
          ->set_end_ipv4_addr(
              cfg_new.rat_specific_context.m5gsm_session_context()
                  .gnode_endpoint()
                  .end_ipv4_addr());
      cfg.rat_specific_context.mutable_m5gsm_session_context()
          ->mutable_gnode_endpoint()
          ->set_teid(cfg_new.rat_specific_context.m5gsm_session_context()
                         .gnode_endpoint()
                         .teid());
      MLOG(MDEBUG) << "2nd Request of session from UE with IMSI: " << imsi
                   << " PDU id " << pdu_id;
      session->set_config(cfg);
      SessionUpdate update =
          SessionStore::get_default_session_update(session_map);
      bool success = m5g_enforcer_->m5g_update_session_context(
          session_map, imsi, session, update);
      if (!success) {
        MLOG(MERROR) << "Failed to update  SessionStore for 5G session "
                     << session->get_session_id();
        return;
      }
      /* update the session changes back to sessionsotre
       */
      bool update_success = session_store_.update_sessions(update);
      if (!update_success) {
        MLOG(MERROR) << "Failed to update the session, re-get gnode teid and ip"
                     << imsi;
        return;
      }
      MLOG(MDEBUG) << " Successfully updated SessionStore of subscriber: "
                   << imsi;
      return;
    }
    // PDU ID found and return from here
    MLOG(MERROR) << "Duplicate request of same PDU_id " << pdu_id << " of IMSI "
                 << imsi << " nothing to do";
    return;
  }
  std::string session_id = id_gen_.gen_session_id(imsi);
  MLOG(MDEBUG) << "First Requested session from UE with IMSI: " << imsi
               << " Generated session " << session_id << " PDU id " << pdu_id;

  auto session_map_ptr = std::make_shared<SessionMap>(std::move(session_map));
  /* initialization of SessionState for IMSI by SessionStateEnforcer*/
  bool success = m5g_enforcer_->m5g_init_session_credit(
      *session_map_ptr, imsi, session_id, cfg_new);
  if (!success) {
    MLOG(MERROR) << "Failed to initialize SessionStore for 5G session "
                 << session_id;
    return;
  } else {
    /* writing of SessionMap in memory through SessionStore object*/
    if (session_store_.create_sessions(
            imsi, std::move((*session_map_ptr)[imsi]))) {
      MLOG(MDEBUG) << "Successfully initialized 5G session for subscriber "
                   << cfg_new.common_context.sid().id()
                   << " with PDU session ID "
                   << cfg_new.rat_specific_context.m5gsm_session_context()
                          .pdu_session_id();
    } else {
      MLOG(MERROR) << "Failed to initialize 5G session for subscriber"
                   << cfg_new.common_context.sid().id()
                   << " with PDU session ID  from UE"
                   << cfg_new.rat_specific_context.m5gsm_session_context()
                          .pdu_session_id()
                   << " due to failure writing to SessionStore.";
    }
  }
}

/* This starts releasing the session in main session enforcer thread context
 * Before startting it checks if respective session */
void SetMessageManagerHandler::initiate_release_session(
    SessionMap& session_map, const std::string& dnn, const std::string& imsi) {
  // TODO as modification and dynamic rules are not implemented this may
  // return empty map.
  auto update = SessionStore::get_default_session_update(session_map);
  bool exist =
      m5g_enforcer_->m5g_release_session(session_map, imsi, dnn, update);
  // If no entry found, nothing to do and return from here
  if (!exist) {
    MLOG(MERROR) << "Entry is not found in SessionStore for subscriber "
                 << imsi;
    return;
  }

  bool update_success = session_store_.update_sessions(update);
  /* No need to respond AMF through gRPC as AMF has informed SessionD
   * to release this session. if failed in update session & store
   * print the fatal error and move on
   */
  if (update_success) {
    MLOG(MINFO)
        << "Successfully released and updated SessionStore of subscriber"
        << imsi;
  }
}

}  // end namespace magma
