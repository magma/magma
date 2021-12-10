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
#include <folly/io/async/EventBase.h>
#include <glog/logging.h>
#include <grpcpp/impl/codegen/status.h>
#include <grpcpp/impl/codegen/status_code_enum.h>
#include <experimental/optional>
#include <ostream>
#include <string>
#include <unordered_map>
#include <utility>
#include <vector>

#include "GrpcMagmaUtils.h"
#include "SessionReporter.h"
#include "SessionState.h"
#include "SessionStateEnforcer.h"
#include "SessionStore.h"
#include "SetMessageManagerHandler.h"
#include "lte/protos/session_manager.pb.h"
#include "lte/protos/subscriberdb.pb.h"
#include "magma_logging.h"

namespace google {
namespace protobuf {
class Message;
}  // namespace protobuf
}  // namespace google
namespace grpc {
class ServerContext;
}  // namespace grpc
namespace magma {
struct SessionStateUpdateCriteria;
}  // namespace magma

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
    SessionStore& session_store, SessionReporter* reporter,
    std::shared_ptr<EventsReporter> events_reporter)
    : session_store_(session_store),
      m5g_enforcer_(m5genforcer),
      reporter_(reporter),
      events_reporter_(events_reporter) {}

/* Building session config with required parameters
 * TODO Note: this function can be removed by implementing 5G specific
 * SeesionConfig constructor
 */
SessionConfig SetMessageManagerHandler::m5g_build_session_config(
    const SetSMSessionContext& request) {
  SessionConfig cfg;
  /*copying only 5G specific data to respective elements*/
  cfg.common_context = request.common_context();
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
    std::string imsi = request_cpy.common_context().sid().id();
    const auto rat_type = request_cpy.common_context().rat_type();
    if (rat_type != TGPP_NR) {
      // We don't support outside of 5G
      std::ostringstream failure_stream;
      failure_stream << "Received an invalid RAT type " << rat_type;
      std::string failure_msg = failure_stream.str();
      MLOG(MERROR) << failure_msg;
      Status status(grpc::FAILED_PRECONDITION, failure_msg);
      response_callback(status, SmContextVoid());
      return;
    }
    // Fetch PDU session ID from rat_specific_context and
    // pdu_id is unique to IMSI
    uint32_t pdu_id = request_cpy.rat_specific_context()
                          .m5gsm_session_context()
                          .pdu_session_id();
    // Fetch complete message from proto message
    SessionConfig cfg = m5g_build_session_config(request_cpy);

    /* Read the proto message and check for state. Get the config out of proto.
     * Code for relase state, then creating
     */
    // Requested message from AMF to release the session
    if (cfg.common_context.sm_session_state() == RELEASED_4) {
      if (cfg.common_context.sm_session_version() == 0) {
        MLOG(MERROR) << "Wrong version received from AMF for IMSI " << imsi
                     << " but continuing release request";
        Status status(grpc::OUT_OF_RANGE, "Version number Out of Range");
        response_callback(status, SmContextVoid());
        return;
      }
      MLOG(MINFO) << "Release request for session from IMSI: " << imsi
                  << " pdu_id " << pdu_id;
      /* Read the SessionMap from global session_store,
       * if it is not found, it will be added w.r.t imsi
       */
      auto session_map = session_store_.read_sessions({imsi});
      initiate_release_session(session_map, pdu_id, imsi);
      response_callback(Status::OK, SmContextVoid());
    } else {
      // The Event Based main_thread invocation and runs to handle session state
      MLOG(MDEBUG) << "Requested session from UE with IMSI: " << imsi
                   << " PDU ID " << pdu_id;

      /* Message may be intial or modification message. Only taken care
       * intial message. Check if it's initial message
       */
      if ((cfg.rat_specific_context.m5gsm_session_context().request_type() ==
           INITIAL_REQUEST) &&
          (cfg.common_context.sm_session_state() == CREATING_0)) {
        /* it's new UE establisment request and need to create the session
         * context
         */
        auto session_map = session_store_.read_sessions({imsi});
        send_create_session(session_map, imsi, cfg, pdu_id);
        response_callback(Status::OK, SmContextVoid());
        return;
      }
      MLOG(MERROR)
          << "AMF request type- Unhandled request type:"
          << cfg.rat_specific_context.m5gsm_session_context().request_type();
      Status status(grpc::UNKNOWN, "Unknown session state or request");
      response_callback(status, SmContextVoid());
    }
  });
}

/* Handling set message from AMF
 * check if PDU_SESSION exists
 * if the PDU_SESSION doen't exists, log and ignore
 * If EXISTING_PDU_SESSION, get the session entry, check on incoming session
 * state and version  accordingly take action, write to memory by SessionStore
 */
void SetMessageManagerHandler::SetSmfNotification(
    ServerContext* context, const SetSmNotificationContext* notif,
    std::function<void(Status, SmContextVoid)> response_callback) {
  // Print the message from AMF
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(*notif));
  // Read the event type from the proto message
  auto& noti = *notif;
  m5g_enforcer_->get_event_base().runInEventBaseThread([this, response_callback,
                                                        noti]() {
    NotifyUeEvents Uevent = noti.rat_specific_notification().notify_ue_event();
    MLOG(MINFO) << "Notification of imsi: " << noti.common_context().sid().id()
                << " from AMF  Event value:" << Uevent;
    switch (Uevent) {
      case PDU_SESSION_INACTIVE_NOTIFY:
        pdu_session_inactive(noti, response_callback);
        return;
      case UE_IDLE_MODE_NOTIFY:
        idle_mode_change_sessions_handle(noti, response_callback);
        return;
      case UE_PAGING_NOTIFY:
        return;
      case UE_PERIODIC_REG_ACTIVE_MODE_NOTIFY:
        return;
      case UE_SERVICE_REQUEST_ON_PAGING:
        service_handle_request_on_paging(noti, response_callback);
        return;
      default:
        return;
    }
  });
}

static CreateSessionRequest make_create_session_request(
    const SessionConfig& cfg, const std::string& session_id,
    const std::unique_ptr<Timezone>& access_timezone) {
  CreateSessionRequest create_request;
  create_request.set_session_id(session_id);
  create_request.mutable_common_context()->CopyFrom(cfg.common_context);
  create_request.mutable_rat_specific_context()->CopyFrom(
      cfg.rat_specific_context);

  if (access_timezone != nullptr) {
    create_request.mutable_access_timezone()->CopyFrom(*access_timezone);
  }

  return create_request;
}

/* Creeate respective SessionState and context*/
void SetMessageManagerHandler::send_create_session(SessionMap& session_map,
                                                   const std::string& imsi,
                                                   SessionConfig& new_cfg,
                                                   uint32_t& pdu_id) {
  /* If it is new session to be created, check for same PDU_ID exists
   * for same IMSI, i.e if IMSI found and respective PDU_ID found in
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
              new_cfg.rat_specific_context.m5gsm_session_context()
                  .gnode_endpoint()
                  .end_ipv4_addr());
      cfg.rat_specific_context.mutable_m5gsm_session_context()
          ->mutable_gnode_endpoint()
          ->set_teid(new_cfg.rat_specific_context.m5gsm_session_context()
                         .gnode_endpoint()
                         .teid());
      MLOG(MDEBUG) << "2nd Request of session from UE with IMSI: " << imsi
                   << " PDU id " << pdu_id;
      session->set_config(cfg, nullptr);
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
  bool success = m5g_enforcer_->m5g_init_session_credit(*session_map_ptr, imsi,
                                                        session_id, new_cfg);
  if (!success) {
    std::ostringstream failure_stream;
    failure_stream << "Failed to initialize SessionStore for 5G session "
                   << session_id;
    std::string failure_msg = failure_stream.str();
    MLOG(MERROR) << failure_msg;
    events_reporter_->session_create_failure(new_cfg, failure_msg);
    return;
  } else {
    /* writing of SessionMap in memory through SessionStore object*/
    if (session_store_.create_sessions(imsi,
                                       std::move((*session_map_ptr)[imsi]))) {
      MLOG(MDEBUG) << "Successfully initialized 5G session for subscriber "
                   << new_cfg.common_context.sid().id()
                   << " with PDU session ID "
                   << new_cfg.rat_specific_context.m5gsm_session_context()
                          .pdu_session_id();
    } else {
      MLOG(MERROR) << "Failed to initialize 5G session for subscriber"
                   << new_cfg.common_context.sid().id()
                   << " with PDU session ID  from UE"
                   << new_cfg.rat_specific_context.m5gsm_session_context()
                          .pdu_session_id()
                   << " due to failure writing to SessionStore.";
    }
  }
  auto create_req = make_create_session_request(
      new_cfg, session_id, m5g_enforcer_->get_access_timezone());
  reporter_->report_create_session(
      create_req, [this, imsi, new_cfg, session_id](
                      Status status, CreateSessionResponse response) mutable {
        if (status.ok()) {
          MLOG(MINFO) << "Processing a CreateSessionResponse for "
                      << session_id;
          PrintGrpcMessage(
              static_cast<const google::protobuf::Message&>(response));
          auto session_map = session_store_.read_sessions({imsi});
          SessionUpdate update =
              SessionStore::get_default_session_update(session_map);
          SessionSearchCriteria criteria(imsi, IMSI_AND_SESSION_ID, session_id);
          auto session_it = session_store_.find_session(session_map, criteria);
          if (!session_it) {
            MLOG(MWARNING) << "Could not find session " << session_id;
            status = Status(grpc::ABORTED, "Session not found");
            return;
          }
          auto& session = **session_it;
          SessionStateUpdateCriteria* session_uc = &update[imsi][session_id];
          m5g_enforcer_->update_session_with_policy(session, response,
                                                    session_uc);
          session_store_.update_sessions(update);
        } else {
          MLOG(MINFO) << "Failed to initialize new session " << session_id
                      << " in SessionD for subscriber " << imsi
                      << " due to failure writing to SessionStore."
                      << " An earlier update may have invalidated it.";
          status = Status(grpc::ABORTED,
                          "Failed to write session to SessionD storage");
          events_reporter_->session_create_failure(
              new_cfg, "Failed to initialize session in SessionProxy/PolicyDB");
        }
      });
}

/* This starts releasing the session in main session enforcer thread context
 * Before startting it checks if respective session */
void SetMessageManagerHandler::initiate_release_session(
    SessionMap& session_map, const uint32_t& pdu_id, const std::string& imsi) {
  // TODO as modification and dynamic rules are not implemented this may
  // return empty map.
  auto update = SessionStore::get_default_session_update(session_map);
  bool exist =
      m5g_enforcer_->m5g_release_session(session_map, imsi, pdu_id, update);
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
    MLOG(MDEBUG)
        << "Successfully released and updated SessionStore of subscriber"
        << imsi;
  }
}

/* This function called when any specific  session moved to
 * idle state
 */
void SetMessageManagerHandler::pdu_session_inactive(
    const SetSmNotificationContext& notif,
    std::function<void(Status, SmContextVoid)> response_callback) {
  // extract values from proto
  uint32_t pdu_id = notif.rat_specific_notification().pdu_session_id();
  std::string imsi = notif.common_context().sid().id();
  /* Read the SessionMap from global session_store */
  SessionSearchCriteria criteria(imsi, IMSI_AND_PDUID, pdu_id);
  auto session_map = session_store_.read_sessions({imsi});
  auto session_it = session_store_.find_session(session_map, criteria);
  auto& session = **session_it;
  if (!session_it) {
    MLOG(MINFO) << " No session found for IMSI: " << imsi << " pdu id "
                << pdu_id;
    Status status(grpc::NOT_FOUND, "Sesion Not found");
    response_callback(status, SmContextVoid());
    return;
  }
  MLOG(MINFO) << "Received Session Notification from UE with IMSI: " << imsi
              << " of session Id " << session->get_session_id();
  if (notif.common_context().sm_session_version() == 0) {
    Status status(grpc::OUT_OF_RANGE, "Version number Out of Range");
    response_callback(status, SmContextVoid());
    return;
  }
  if (notif.rat_specific_notification().request_type() !=
      EXISTING_PDU_SESSION) {
    MLOG(MINFO) << " Wrong request type for sesion id"
                << session->get_session_id();
    Status status(grpc::UNKNOWN, "Unknown session request type");
    response_callback(status, SmContextVoid());
    return;
  }
  if (session->get_state() == RELEASE) {
    // Nothing to be done;
    MLOG(MINFO) << " No sessions to move to Idle state : " << imsi;
    response_callback(Status::OK, SmContextVoid());
    return;
  }
  auto session_update = SessionStore::get_default_session_update(session_map);
  auto session_id = session->get_session_id();
  SessionStateUpdateCriteria& session_uc = session_update[imsi][session_id];
  m5g_enforcer_->m5g_move_to_inactive_state(imsi, **session_it, notif,
                                            &session_uc);
  bool update_success = session_store_.update_sessions(session_update);
  if (update_success) {
    MLOG(MDEBUG) << "Successfully updated SessionStore "
                 << "of subscriber" << imsi;
    response_callback(Status::OK, SmContextVoid());
    return;
  }
  Status status(grpc::ABORTED, "update operation aborted in  middle");
  response_callback(status, SmContextVoid());
  return;
}
/* This function callled when specific IMSi and its associated
 * sessions moved to idle mode
 */
void SetMessageManagerHandler::idle_mode_change_sessions_handle(
    const SetSmNotificationContext& notif,
    std::function<void(Status, SmContextVoid)> response_callback) {
  // extract IMSI value from proto
  auto imsi = notif.common_context().sid().id();
  auto session_map = session_store_.read_sessions({imsi});
  int count = 0;
  auto session_update = SessionStore::get_default_session_update(session_map);
  for (auto& session : session_map[imsi]) {
    if (session->get_state() != RELEASE) {
      auto session_id = session->get_session_id();
      SessionStateUpdateCriteria& session_uc = session_update[imsi][session_id];
      m5g_enforcer_->m5g_move_to_inactive_state(imsi, session, notif,
                                                &session_uc);
      bool update_success = session_store_.update_sessions(session_update);
      if (!update_success) {
        MLOG(MINFO) << "Operation aborted in  middle"
                    << " session id: " << session->get_session_id() << imsi
                    << " of imsi";
        continue;
      }
      MLOG(MDEBUG) << "Successfully updated SessionStore "
                   << " session id: " << session->get_session_id() << imsi
                   << " of imsi";
      count++;
    }
  }
  if (!count) {
    MLOG(MINFO) << " No Valid session found for IMSI: " << imsi;
    Status status(grpc::NOT_FOUND, "Sesion Not found");
    response_callback(status, SmContextVoid());
    return;
  }
  bool update_success = session_store_.update_sessions(session_update);
  if (!update_success) {
    Status status(grpc::ABORTED, "update operation aborted in  middle");
    response_callback(status, SmContextVoid());
    return;
  }
  MLOG(MDEBUG) << "Successfully updated SessionStore "
               << "of subscriber :" << imsi;
  response_callback(Status::OK, SmContextVoid());
  return;
}

/* This function callled when specific IMSI service
 * request is received from AMF.
 */
void SetMessageManagerHandler::service_handle_request_on_paging(
    const SetSmNotificationContext& notif,
    std::function<void(Status, SmContextVoid)> response_callback) {
  // extract IMSI value from proto
  auto imsi = notif.common_context().sid().id();
  auto session_map = session_store_.read_sessions({imsi});
  int count = 0;
  auto session_update = SessionStore::get_default_session_update(session_map);
  for (auto& session : session_map[imsi]) {
    if (session->get_state() == INACTIVE) {
      auto session_id = session->get_session_id();
      SessionStateUpdateCriteria& session_uc = session_update[imsi][session_id];
      m5g_enforcer_->m5g_move_to_active_state(session, notif, &session_uc);
      bool update_success = session_store_.update_sessions(session_update);
      if (!update_success) {
        MLOG(MINFO) << "Operation aborted in  middle"
                    << " session id: " << session->get_session_id() << imsi
                    << " of imsi";
        continue;
      }
      MLOG(MDEBUG) << "Successfully updated SessionStore "
                   << " session id: " << session->get_session_id() << imsi
                   << " of imsi";
      count++;
    }
  }
  if (!count) {
    MLOG(MINFO) << " No Valid session found for IMSI: " << imsi;
    Status status(grpc::NOT_FOUND, "Sesion Not found");
    response_callback(status, SmContextVoid());
    return;
  }
  bool update_success = session_store_.update_sessions(session_update);
  if (!update_success) {
    Status status(grpc::ABORTED, "update operation aborted in  middle");
    response_callback(status, SmContextVoid());
    return;
  }
  // Send  Paging service request response back to AMF
  auto session = session_map[imsi].begin();
  m5g_enforcer_->handle_state_update_to_amf(
      **session, magma::lte::M5GSMCause::OPERATION_SUCCESS,
      UE_SERVICE_REQUEST_ON_PAGING);
  MLOG(MDEBUG) << "Successfully updated SessionStore "
               << "of subscriber :" << imsi;
  response_callback(Status::OK, SmContextVoid());
  return;
}
}  // end namespace magma
