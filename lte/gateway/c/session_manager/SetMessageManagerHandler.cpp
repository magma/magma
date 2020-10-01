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
  /*copying pnly 5G specific data to respective elements*/
  cfg.common_context       = request.common_context();
  cfg.rat_specific_context = request.rat_specific_context();

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
  response_callback(Status::OK, SmContextVoid());
  /*The Event Based main_thread invocation and runs to handle session state*/
  m5g_enforcer_->get_event_base().runInEventBaseThread(
      [this, response_callback, request_cpy]() {
        // extract values from proto
        auto imsi = request_cpy.common_context().sid().id();
        std::string dnn =
            request_cpy.common_context().apn();  // may not required for demo-1
        std::string session_id = id_gen_.gen_session_id(imsi);

        MLOG(MDEBUG) << "Requested session from UE with IMSI: " << imsi
                     << " Generated sessioncontext ID" << session_id;
        /*reach complete message from proto message*/
        SessionConfig cfg = m5g_build_session_config(request_cpy);

        /*check if it's initial message*/
        if ((cfg.rat_specific_context.m5gsm_session_context().rquest_type() ==
             INITIAL_REQUEST) &&
            (cfg.common_context.sm_session_state() == CREATING_0)) {
          /* it's new UE establisment request and need to create the session
           * context
           */
          MLOG(MDEBUG)
              << "AMF request type INITIAL_REQUEST and session state CREATING";
          /* Read the SessionMap from global session_store,
           * if it is not found, it will be added w.r.t imsi
           */
          auto session_map = session_store_.read_sessions({imsi});
          send_create_session(session_map, imsi, session_id, cfg);
        }
      });
}

/* Creeate respective SessionState and context*/
void SetMessageManagerHandler::send_create_session(
    SessionMap& session_map, const std::string& imsi,
    const std::string& session_id, const SessionConfig& cfg) {
  auto session_map_ptr = std::make_shared<SessionMap>(std::move(session_map));
  /* initialization of SessionState for IMSI by SessionStateEnforcer*/
  bool success = m5g_enforcer_->m5g_init_session_credit(
      *session_map_ptr, imsi, session_id, cfg);
  if (!success) {
    MLOG(MERROR) << "Failed to initialize SessionStore for 5G session "
                 << session_id << " IMSI "
                 << " imsi";
    return;
  } else {
    /* writing of SessionMap in memory through SessionStore object*/
    if (session_store_.create_sessions(
            imsi, std::move((*session_map_ptr)[imsi]))) {
      MLOG(MINFO)
          << "Successfully initialized 5G session for subscriber "
          << cfg.common_context.sid().id() << " with PDU session ID "
          << cfg.rat_specific_context.m5gsm_session_context().pdu_session_id();
    } else {
      MLOG(MERROR)
          << "Failed to initialize 5G session for subscriber"
          << cfg.common_context.sid().id() << " with PDU session ID "
          << cfg.rat_specific_context.m5gsm_session_context().pdu_session_id()
          << " due to failure writing to SessionStore.";
    }
  }
}

}  // end namespace magma
