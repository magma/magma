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
  Source      	SetMessageManagerHandler.h
  Version     	0.1
  Date       	2020/08/08
  Product     	SessionD
  Subsystem   	5G Landing object in SessionD
  Author/Editor Sanjay Kumar Ojha
  Description 	Defines Access and Mobility Management Messages
*****************************************************************************/
#pragma once
#include <grpc++/grpc++.h>
#include <lte/protos/session_manager.grpc.pb.h>
#include <stdint.h>
#include <functional>
#include <memory>
#include <string>

#include "SessionID.h"
#include "SessionReporter.h"
#include "SessionStateEnforcer.h"
#include "SessionStore.h"
#include "StoreClient.h"
#include "Types.h"
#include "orc8r/protos/common.pb.h"

namespace grpc {
class ServerContext;
class Status;
}  // namespace grpc
namespace magma {
class SessionReporter;
class SessionStateEnforcer;
namespace lte {
class SessionStore;
class SetSMSessionContext;
class SetSmNotificationContext;
class SmContextVoid;
}  // namespace lte
}  // namespace magma

using grpc::ServerContext;
using grpc::Status;
namespace magma {
using namespace orc8r;

/* SetMessageManagerHandler processes gRPC requests for the sessionD
 * This composites the earlier LocalSessionManagerHandlerImpl and uses the
 * exiting functionalities.
 */

/* This the landing object of 5G gRPC call by set message*/
class SetMessageManager {
 public:
  virtual ~SetMessageManager() {}

  /* RPC call from AMF "rpc SetAmfSessionContext (SetSMSessionContext) returns
   * (SmContextVoid);" as its set-interface API, no need to send response back,
   * response is void and gRPC will take care on acknowledgement
   */
  virtual void SetSmfNotification(
      ServerContext* context, const SetSmNotificationContext* notif,
      std::function<void(Status, SmContextVoid)> response_callback) = 0;

  virtual void SetAmfSessionContext(
      ServerContext* context, const SetSMSessionContext* request,
      std::function<void(Status, SmContextVoid)> response_callback) = 0;
};  // end of abstract class

class SetMessageManagerHandler : public SetMessageManager {
 public:
  SetMessageManagerHandler(std::shared_ptr<SessionStateEnforcer> m5genforcer,
                           SessionStore& session_store,
                           SessionReporter* reporter,
                           std::shared_ptr<EventsReporter> events_reporter);
  ~SetMessageManagerHandler() {}

  /* Paging, idle state change notifcation receiving */
  virtual void SetSmfNotification(
      ServerContext* context, const SetSmNotificationContext* notif,
      std::function<void(Status, SmContextVoid)> response_callback);

  virtual void SetAmfSessionContext(
      ServerContext* context, const SetSMSessionContext* request,
      std::function<void(Status, SmContextVoid)> response_callback);

  /* When any specific IMIS/PDU id session is in-active */
  void pdu_session_inactive(
      const SetSmNotificationContext& notif,
      std::function<void(Status, SmContextVoid)> response_callback);

  /* When any IMSI is moved to inactive state */
  void idle_mode_change_sessions_handle(
      const SetSmNotificationContext& notif,
      std::function<void(Status, SmContextVoid)> response_callback);

  /* Handle service request received from AMF after paging request is sent
   */
  void service_handle_request_on_paging(
      const SetSmNotificationContext& notif,
      std::function<void(Status, SmContextVoid)> response_callback);

  /*
   * Send session creation related request to the CentralSessionController.
   * which is policy/QoS related. On successful, creates and populate,
   * session_map in memoery and response set message to AMF by gRPC.
   * It uses SessionStateEnforcer object to create new session state.
   */
  void send_create_session(SessionMap& session_map, const std::string& imsi,
                           SessionConfig& cfg, uint32_t& pdu_id);
  /*initialize the session message from proto message*/
  SessionConfig m5g_build_session_config(const SetSMSessionContext& request);

  /*Release request message handling*/
  void initiate_release_session(SessionMap& session_map, const uint32_t& pdu_id,
                                const std::string& imsi);

 private:
  SessionStore& session_store_;
  std::shared_ptr<SessionStateEnforcer> m5g_enforcer_;
  SessionReporter* reporter_;
  std::shared_ptr<EventsReporter> events_reporter_;
  SessionIDGenerator id_gen_;

  bool validate_session_request(const SessionConfig cfg);

};  // end of class SetMessageManagerHandlerImpl

}  // end namespace magma
