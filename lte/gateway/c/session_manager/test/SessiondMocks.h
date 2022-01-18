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

#include <gmock/gmock.h>
#include <grpc++/grpc++.h>
#include <gtest/gtest.h>
#include <lte/protos/pipelined.grpc.pb.h>
#include <lte/protos/pipelined.pb.h>
#include <lte/protos/policydb.pb.h>
#include <lte/protos/session_manager.grpc.pb.h>
#include <orc8r/protos/eventd.pb.h>

#include <memory>
#include <string>
#include <vector>

#include "LocalSessionManagerHandler.h"
#include "UpfMsgManageHandler.h"
#include "PipelinedClient.h"
#include "RuleStore.h"
#include "SessionReporter.h"
#include "SessionState.h"
#include "SpgwServiceClient.h"

using grpc::Status;
using ::testing::_;
using ::testing::Return;

namespace magma {
/**
 * Mock handler to mock actual request handling and just test server
 */
class MockPipelined final : public Pipelined::Service {
 public:
  MockPipelined() : Pipelined::Service() {
    ON_CALL(*this, ActivateFlows(_, _, _)).WillByDefault(Return(Status::OK));
    ON_CALL(*this, DeactivateFlows(_, _, _)).WillByDefault(Return(Status::OK));
    ON_CALL(*this, SetupPolicyFlows(_, _, _)).WillByDefault(Return(Status::OK));
    ON_CALL(*this, SetupDefaultControllers(_, _, _))
        .WillByDefault(Return(Status::OK));
    ON_CALL(*this, SetupUEMacFlows(_, _, _)).WillByDefault(Return(Status::OK));
    ON_CALL(*this, SetupQuotaFlows(_, _, _)).WillByDefault(Return(Status::OK));
    ON_CALL(*this, GetStats(_, _, _)).WillByDefault(Return(Status::OK));
    ON_CALL(*this, SetSMFSessions(_, _, _)).WillByDefault(Return(Status::OK));
  }

  MOCK_METHOD3(ActivateFlows,
               Status(grpc::ServerContext*, const ActivateFlowsRequest*,
                      ActivateFlowsResult*));
  MOCK_METHOD3(DeactivateFlows,
               Status(grpc::ServerContext*, const DeactivateFlowsRequest*,
                      DeactivateFlowsResult*));
  MOCK_METHOD3(SetupPolicyFlows,
               Status(grpc::ServerContext*, const SetupPolicyRequest*,
                      SetupFlowsResult*));
  MOCK_METHOD3(SetupDefaultControllers,
               Status(grpc::ServerContext*, const SetupDefaultRequest*,
                      SetupFlowsResult*));
  MOCK_METHOD3(SetupUEMacFlows,
               Status(grpc::ServerContext*, const SetupUEMacRequest*,
                      SetupFlowsResult*));
  MOCK_METHOD3(SetupQuotaFlows,
               Status(grpc::ServerContext*, const SetupQuotaRequest*,
                      SetupFlowsResult*));
  MOCK_METHOD3(GetStats, Status(grpc::ServerContext*, const GetStatsRequest*,
                                RuleRecordTable*));
  MOCK_METHOD3(SetSMFSessions, Status(grpc::ServerContext*, const SessionSet*,
                                      UPFSessionContextState*));
};

class MockPipelinedClient : public PipelinedClient {
 public:
  MockPipelinedClient() {
    ON_CALL(*this, get_next_teid()).WillByDefault(Return(0));
    ON_CALL(*this, get_current_teid()).WillByDefault(Return(0));
  }

  MOCK_METHOD9(
      setup_cwf,
      void(const std::vector<SessionState::SessionInfo>& infos,
           const std::vector<SubscriberQuotaUpdate>& quota_updates,
           const std::vector<std::string> ue_mac_addrs,
           const std::vector<std::string> msisdns,
           const std::vector<std::string> apn_mac_addrs,
           const std::vector<std::string> apn_names,
           const std::vector<std::uint64_t> pdp_start_times,
           const std::uint64_t& epoch,
           std::function<void(Status status, SetupFlowsResult)> callback));
  MOCK_METHOD3(
      setup_lte,
      void(const std::vector<SessionState::SessionInfo>& infos,
           const std::uint64_t& epoch,
           std::function<void(Status status, SetupFlowsResult)> callback));
  MOCK_METHOD6(deactivate_flows_for_rules,
               void(const std::string& imsi, const std::string& ip_addr,
                    const std::string& ipv6_addr, const Teids teids,
                    const RulesToProcess to_process,
                    const RequestOriginType_OriginType origin_type));
  MOCK_METHOD5(deactivate_flows_for_rules_for_termination,
               void(const std::string& imsi, const std::string& ip_addr,
                    const std::string& ipv6_addr,
                    const std::vector<Teids>& teids,
                    const RequestOriginType_OriginType origin_type));
  MOCK_METHOD8(
      activate_flows_for_rules,
      void(const std::string& imsi, const std::string& ip_addr,
           const std::string& ipv6_addr, const Teids teids,
           const std::string& msisdn,
           const std::experimental::optional<AggregatedMaximumBitrate>& ambr,
           const RulesToProcess to_process,
           std::function<void(Status status, ActivateFlowsResult)> callback));
  MOCK_METHOD6(add_ue_mac_flow,
               void(const SubscriberID& sid, const std::string& ue_mac_addr,
                    const std::string& msisdn, const std::string& ap_mac_addr,
                    const std::string& ap_name,
                    std::function<void(Status status, FlowResponse)> callback));
  MOCK_METHOD6(update_ipfix_flow,
               void(const SubscriberID& sid, const std::string& ue_mac_addr,
                    const std::string& msisdn, const std::string& ap_mac_addr,
                    const std::string& ap_name,
                    const uint64_t& pdp_start_time));
  MOCK_METHOD1(update_subscriber_quota_state,
               void(const std::vector<SubscriberQuotaUpdate>& updates));
  MOCK_METHOD2(delete_ue_mac_flow,
               void(const SubscriberID& sid, const std::string& ue_mac_addr));
  MOCK_METHOD6(add_gy_final_action_flow,
               void(const std::string& imsi, const std::string& ip_addr,
                    const std::string& ipv6_addr, const Teids teids,
                    const std::string& msisdn,
                    const RulesToProcess to_process));
  MOCK_METHOD4(set_upf_session,
               void(const SessionState::SessionInfo info,
                    const RulesToProcess to_activate_process,
                    const RulesToProcess to_deactivate_process,
                    std::function<void(Status status, UPFSessionContextState)>
                        callback));
  MOCK_METHOD0(get_next_teid, uint32_t());
  MOCK_METHOD0(get_current_teid, uint32_t());
  MOCK_METHOD3(poll_stats,
               void(int cookie, int cookie_mask,
                    std::function<void(Status, RuleRecordTable)> callback));
};

class MockDirectorydClient : public DirectorydClient {
 public:
  MOCK_METHOD2(update_directoryd_record,
               void(const UpdateRecordRequest& request,
                    std::function<void(Status status, Void)> callback));
  MOCK_METHOD2(delete_directoryd_record,
               void(const DeleteRecordRequest& request,
                    std::function<void(Status status, Void)> callback));
  MOCK_METHOD1(
      get_all_directoryd_records,
      void(std::function<void(Status status, AllDirectoryRecords)> callback));
};

class MockEventdClient : public EventdClient {
 public:
  MOCK_METHOD2(log_event,
               void(const Event& request,
                    std::function<void(Status status, Void)> callback));
};

/**
 * Mock handler to mock actual request handling and just test server
 */
class MockCentralController final : public CentralSessionController::Service {
 public:
  MOCK_METHOD3(CreateSession,
               Status(grpc::ServerContext*, const CreateSessionRequest*,
                      CreateSessionResponse*));

  MOCK_METHOD3(UpdateSession,
               Status(grpc::ServerContext*, const UpdateSessionRequest*,
                      UpdateSessionResponse*));

  MOCK_METHOD3(TerminateSession,
               Status(grpc::ServerContext*, const SessionTerminateRequest*,
                      SessionTerminateResponse*));
};

class MockCallback {
 public:
  MOCK_METHOD2(update_callback,
               void(Status status, const UpdateSessionResponse&));
  MOCK_METHOD2(create_callback,
               void(Status status, const CreateSessionResponse&));
};

/**
 * Mock handler to mock actual request handling and just test server
 */
class MockSessionHandler final : public LocalSessionManagerHandler {
 public:
  ~MockSessionHandler() {}

  MOCK_METHOD3(ReportRuleStats,
               void(grpc::ServerContext*, const RuleRecordTable*,
                    std::function<void(Status, Void)>));

  MOCK_METHOD3(CreateSession,
               void(grpc::ServerContext*, const LocalCreateSessionRequest*,
                    std::function<void(Status, LocalCreateSessionResponse)>));

  MOCK_METHOD3(EndSession,
               void(grpc::ServerContext*, const LocalEndSessionRequest*,
                    std::function<void(Status, LocalEndSessionResponse)>));
};

class MockSessionReporter : public SessionReporter {
 public:
  MOCK_METHOD2(report_updates,
               void(const UpdateSessionRequest&,
                    std::function<void(grpc::Status, UpdateSessionResponse)>));

  MOCK_METHOD2(report_create_session,
               void(const CreateSessionRequest&,
                    std::function<void(Status, CreateSessionResponse)>));

  MOCK_METHOD2(report_terminate_session,
               void(const SessionTerminateRequest&,
                    std::function<void(Status, SessionTerminateResponse)>));
};

class MockAAAClient : public aaa::AAAClient {
 public:
  MockAAAClient() {
    ON_CALL(*this, terminate_session(_, _)).WillByDefault(Return(true));
    ON_CALL(*this, add_sessions(_)).WillByDefault(Return(true));
  }

  MOCK_METHOD2(terminate_session, bool(const std::string& radius_session_id,
                                       const std::string& imsi));

  MOCK_METHOD1(add_sessions, bool(const magma::lte::SessionMap& session_map));
};

class MockSpgwServiceClient : public SpgwServiceClient {
 public:
  MockSpgwServiceClient() {
    ON_CALL(*this, delete_default_bearer(_, _, _)).WillByDefault(Return(true));
    ON_CALL(*this, delete_dedicated_bearer(_)).WillByDefault(Return(true));
    ON_CALL(*this, create_dedicated_bearer(_)).WillByDefault(Return(true));
  }
  MOCK_METHOD3(delete_default_bearer,
               bool(const std::string&, const std::string&, const uint32_t));

  MOCK_METHOD1(delete_dedicated_bearer, bool(const DeleteBearerRequest&));
  MOCK_METHOD1(create_dedicated_bearer, bool(const CreateBearerRequest&));
};

class MockEventsReporter : public EventsReporter {
 public:
  MOCK_METHOD4(session_created,
               void(const std::string&, const std::string&,
                    const SessionConfig&,
                    const std::unique_ptr<SessionState>& session));
  MOCK_METHOD2(session_create_failure,
               void(const SessionConfig&, const std::string&));
  MOCK_METHOD3(session_updated, void(const std::string&, const SessionConfig&,
                                     const UpdateRequests&));
  MOCK_METHOD4(session_update_failure,
               void(const std::string&, const SessionConfig&,
                    const UpdateRequests&, const std::string&));
  MOCK_METHOD2(session_terminated,
               void(const std::string&, const std::unique_ptr<SessionState>&));
};

class MockSetInterfaceForUserPlane final
    : public SetInterfaceForUserPlane::Service {
 public:
  MockSetInterfaceForUserPlane() : SetInterfaceForUserPlane::Service() {}
  MOCK_METHOD3(SetUPFNodeState,
               Status(grpc::ServerContext*, const UPFNodeState*,
                      std::function<void(Status, SmContextVoid)>));
  MOCK_METHOD3(SetUPFSessionConfig,
               Status(grpc::ServerContext*, const UPFSessionConfigState*,
                      std::function<void(Status, SmContextVoid)>));
};

class MockAmfServiceClient : public AmfServiceClient {
 public:
  MockAmfServiceClient() {
    ON_CALL(*this, handle_response_to_access(_)).WillByDefault(Return(true));
  }

  MOCK_METHOD1(handle_response_to_access,
               bool(const magma::SetSMSessionContextAccess&));

  MOCK_METHOD1(handle_notification_to_access,
               bool(const magma::SetSmNotificationContext&));
};

class MockMobilitydClient : public MobilitydClient {
 public:
  MockMobilitydClient() {
    ON_CALL(*this, get_subscriberid_from_ipv4(_, _)).WillByDefault(Return());
  }
  MOCK_METHOD2(get_subscriberid_from_ipv4,
               void(const IPAddress&,
                    std::function<void(Status status, SubscriberID)>));
};

}  // namespace magma
