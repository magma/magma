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

#include <grpc++/grpc++.h>
#include <grpcpp/impl/codegen/status.h>
#include <lte/protos/session_manager.grpc.pb.h>
#include <stdint.h>
#include <chrono>
#include <experimental/optional>
#include <functional>
#include <memory>
#include <string>

#include "LocalEnforcer.h"
#include "SessionID.h"
#include "SessionReporter.h"
#include "SessionStore.h"
#include "StoreClient.h"
#include "orc8r/protos/common.pb.h"

namespace grpc {
class ServerContext;
}  // namespace grpc
namespace magma {
class DirectorydClient;
class LocalEnforcer;
class SessionReporter;
class SessionState;
namespace lte {
class EventsReporter;
class LocalCreateSessionRequest;
class LocalCreateSessionResponse;
class LocalEndSessionRequest;
class LocalEndSessionResponse;
class PolicyBearerBindingRequest;
class PolicyBearerBindingResponse;
class RuleRecordTable;
class SessionRules;
class SetupFlowsResult;
class SubscriberID;
class UpdateTunnelIdsRequest;
class UpdateTunnelIdsResponse;
}  // namespace lte
struct SessionConfig;
}  // namespace magma

using grpc::ServerContext;
using grpc::Status;

namespace magma {
using namespace orc8r;
using std::experimental::optional;

// Utility struct to capture all actions that could result from handling a
// LocalCreateSessionRequest
struct SessionActionOrStatus {
  bool create_new_session;
  // If true, end the existing session registered (in SessionD and policy
  // component) for the IMSI and APN in the request
  bool end_existing_session;
  // SessionID that needs to be sent back with LocalCreateSessionResponse
  std::string session_id_to_send_back;
  // If this value is set, we should respond back to access immediately
  optional<grpc::Status> status_back_to_access;
  SessionActionOrStatus()
      : create_new_session(false),
        end_existing_session(false),
        status_back_to_access({}) {}
  static SessionActionOrStatus create_new_session_action(
      const std::string session_id) {
    SessionActionOrStatus action;
    action.create_new_session = true;
    action.session_id_to_send_back = session_id;
    return action;
  }
  static SessionActionOrStatus status(grpc::Status status) {
    SessionActionOrStatus action;
    action.status_back_to_access = status;
    return action;
  }
  static SessionActionOrStatus OK(const std::string session_id) {
    SessionActionOrStatus action;
    action.status_back_to_access = grpc::Status::OK;
    action.session_id_to_send_back = session_id;
    return action;
  }
  void set_create_new_session(const std::string session_id) {
    create_new_session = true;
    session_id_to_send_back = session_id;
  }
  void set_end_existing_session() { create_new_session = true; }
  void set_status(grpc::Status status) { status_back_to_access = status; }
};

class LocalSessionManagerHandler {
 public:
  enum PipelineDState {
    // PipelineD restarted and has not been setup
    NOT_READY = 0,
    // Currently exchanging Setup requests
    SETTING_UP = 1,
    // PipelineD is setup and is ready to accept requests
    READY = 2,
  };

  virtual ~LocalSessionManagerHandler() {}

  /**
   * Report flow stats from pipelined and track the usage per rule
   */
  virtual void ReportRuleStats(
      ServerContext* context, const RuleRecordTable* request,
      std::function<void(Status, Void)> response_callback) = 0;

  /**
   * Create a new session, initializing credit monitoring and requesting credit
   * from the cloud
   */
  virtual void CreateSession(
      ServerContext* context, const LocalCreateSessionRequest* request,
      std::function<void(Status, LocalCreateSessionResponse)>
          response_callback) = 0;

  /**
   * Terminate a session, untracking credit and terminating in the cloud
   */
  virtual void EndSession(ServerContext* context,
                          const LocalEndSessionRequest* request,
                          std::function<void(Status, LocalEndSessionResponse)>
                              response_callback) = 0;

  /**
   * Bind the returned bearer id to the policy for which it is created
   */
  virtual void BindPolicy2Bearer(
      ServerContext* context, const PolicyBearerBindingRequest* request,
      std::function<void(Status, PolicyBearerBindingResponse)>
          response_callback) = 0;

  /**
   * Updates eNB and AGW tunnels id on a existing session for a default bearer
   */
  virtual void UpdateTunnelIds(
      ServerContext* context, UpdateTunnelIdsRequest* request,
      std::function<void(Status, UpdateTunnelIdsResponse)>
          response_callback) = 0;

  /**
   * Update active rules for session
   */
  virtual void SetSessionRules(
      ServerContext* context, const SessionRules* request,
      std::function<void(Status, Void)> response_callback) = 0;
};

/**
 * LocalSessionManagerHandler processes proxied gRPC requests to the session
 * manager. The handler uses a monitor and reporter to keep track of state
 * and report to the cloud, respectively
 */
class LocalSessionManagerHandlerImpl : public LocalSessionManagerHandler {
 public:
  LocalSessionManagerHandlerImpl(
      std::shared_ptr<LocalEnforcer> monitor, SessionReporter* reporter,
      std::shared_ptr<DirectorydClient> directoryd_client,
      std::shared_ptr<EventsReporter> events_reporter,
      SessionStore& session_store);
  ~LocalSessionManagerHandlerImpl() {}
  /**
   * Report flow stats from pipelined and track the usage per rule
   */
  void ReportRuleStats(ServerContext* context, const RuleRecordTable* request,
                       std::function<void(Status, Void)> response_callback);

  /**
   * Create a new session, initializing credit monitoring and requesting credit
   * from the cloud
   */
  void CreateSession(ServerContext* context,
                     const LocalCreateSessionRequest* request,
                     std::function<void(Status, LocalCreateSessionResponse)>
                         response_callback);

  /**
   * Terminate a session, untracking credit and terminating in the cloud
   */
  void EndSession(
      ServerContext* context, const LocalEndSessionRequest* request,
      std::function<void(Status, LocalEndSessionResponse)> response_callback);

  /**
   * Bind the returned bearer id to the policy for which it is created; if
   * the returned bearer id is 0 then the dedicated bearer request is rejected
   */
  void BindPolicy2Bearer(
      ServerContext* context, const PolicyBearerBindingRequest* request,
      std::function<void(Status, PolicyBearerBindingResponse)>
          response_callback);

  /**
   * Updates eNB and AGW tunnels id on a existing session for a default bearer
   */
  void UpdateTunnelIds(
      ServerContext* context, UpdateTunnelIdsRequest* request,
      std::function<void(Status, UpdateTunnelIdsResponse)> response_callback);

  /**
   * Update active rules for session
   * Get the SessionMap for the updates, apply the set rules and update the
   * store. The rule updates should be also propagated to PipelineD
   */
  void SetSessionRules(ServerContext* context, const SessionRules* request,
                       std::function<void(Status, Void)> response_callback);

 private:
  SessionStore& session_store_;
  std::shared_ptr<LocalEnforcer> enforcer_;
  SessionReporter* reporter_;
  std::shared_ptr<DirectorydClient> directoryd_client_;
  std::shared_ptr<EventsReporter> events_reporter_;
  SessionIDGenerator id_gen_;
  uint64_t current_epoch_;
  uint64_t reported_epoch_;
  std::chrono::milliseconds retry_timeout_ms_;
  PipelineDState pipelined_state_;
  static const std::string hex_digit_;

 private:
  void check_usage_for_reporting(SessionMap session_map,
                                 SessionUpdate& session_uc);
  bool is_pipelined_restarted();
  void call_setup_pipelined(const std::uint64_t& epoch,
                            const bool update_rule_versions);

  void end_session(
      SessionMap& session_map, const SubscriberID& sid, const std::string& apn,
      std::function<void(Status, LocalEndSessionResponse)> response_callback);

  void add_session_to_directory_record(const std::string& imsi,
                                       const std::string& session_id,
                                       const std::string& msisdn);

  /**
   * handle_create_session_cwf handles a sequence of actions needed for the
   * RATType=WLAN case.
   * If there is an existing session for the IMSI that is ACTIVE, we will
   * simply update its SessionConfig with the new context. In this case, we will
   * NOT send a CreateSession request into FeG/PolicyDB.
   * Otherwise, we will go through the procedure of creating a new context.
   * @param session_map - SessionMap that contains all sessions with IMSI
   * @param session_id - newly created SessionID
   * @param cfg - newly created SessionConfig from the LocalCreateSessionRequest
   * @returns if we should immediately respond to the access component without
   * creating a new session in the policy component, it'll return a struct with
   * grpc::Status set. otherwise, it'll return a struct with fields set to
   * indicate what actions need to be taken.
   */
  SessionActionOrStatus handle_create_session_cwf(
      SessionMap& session_map, const std::string& session_id,
      const SessionConfig& cfg) const;

  /**
   * @brief Handle the logic to recycle an existing CWF session. This involves
   * updating the existing SessionConfig with the new one. This function is
   * responsible for responding to the original LocalCreateSession call with the
   * cb.
   *
   * @param cfg
   * @param session_map
   * @return SessionActionOrStatus with grpc::Status set
   */
  SessionActionOrStatus recycle_cwf_session(
      std::unique_ptr<SessionState>& session, const SessionConfig& cfg,
      SessionMap& session_map) const;
  /**
   * handle_create_session_lte handles a sequence of actions needed for the
   * RATType=LTE case. It is responsible for responding to the original
   * LocalCreateSession request.
   * If there is an existing identical session, same SessionConfig, for the IMSI
   * that is ACTIVE, we will reuse this session. In this case, we will NOT send
   * a CreateSession request into FeG/PolicyDB.
   * Otherwise, we will go through the procedure of creating a new context.
   * @param session_map - SessionMap that contains all sessions with IMSI
   * @param session_id - newly created SessionID
   * @param cfg - newly created SessionConfig from the LocalCreateSessionRequest
   * @returns if we should immediately respond to the access component without
   * creating a new session in the policy component, it'll return a struct with
   * grpc::Status set. otherwise, it'll return a struct with fields set to
   * indicate what actions need to be taken.
   */
  SessionActionOrStatus handle_create_session_lte(SessionMap& session_map,
                                                  const std::string& session_id,
                                                  const SessionConfig& cfg);

  /**
   * Send session creation request to the CentralSessionController.
   * If it is successful, create a session in session_map, and respond to
   * gRPC caller.
   */
  void send_create_session(
      SessionMap& session_map, const std::string& session_id,
      const SessionConfig& cfg,
      std::function<void(grpc::Status, LocalCreateSessionResponse)> cb);

  void handle_setup_callback(const std::uint64_t& epoch, Status status,
                             SetupFlowsResult resp);

  void send_local_create_session_response(
      Status status, const std::string& sid,
      std::function<void(Status, LocalCreateSessionResponse)>
          response_callback);

  void log_create_session(const SessionConfig& cfg);

  /**
   * @param cfg
   * @return status::OK if the request can be processed, otherwise the status
   * that should be sent back
   */
  grpc::Status validate_create_session_request(const SessionConfig cfg);

  /**
   * @brief SessionD will not process any requests if
   * 1. SessionStore is not ready
   * 2. PipelineD is not ready
   *
   * @return status::OK if SessionD is ready to accept requests
   */
  grpc::Status check_sessiond_is_ready();

  bool initialize_session(SessionMap& session_map,
                          const std::string& session_id,
                          const SessionConfig& cfg);
};

}  // namespace magma
