/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include <grpc++/grpc++.h>
#include <gtest/gtest.h>
#include <gmock/gmock.h>

#include <lte/protos/policydb.pb.h>
#include <lte/protos/pipelined.grpc.pb.h>
#include <lte/protos/pipelined.pb.h>
#include <lte/protos/session_manager.grpc.pb.h>

#include <folly/io/async/EventBase.h>

#include "SessionReporter.h"
#include "LocalSessionManagerHandler.h"
#include "PipelinedClient.h"
#include "RuleStore.h"
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
  MockPipelined(): Pipelined::Service()
  {
    ON_CALL(*this, AddRule(_, _, _)).WillByDefault(Return(Status::OK));
    ON_CALL(*this, ActivateFlows(_, _, _)).WillByDefault(Return(Status::OK));
    ON_CALL(*this, DeactivateFlows(_, _, _)).WillByDefault(Return(Status::OK));
  }

  MOCK_METHOD3(
    AddRule,
    Status(grpc::ServerContext *, const PolicyRule *, Void *));
  MOCK_METHOD3(
    ActivateFlows,
    Status(
      grpc::ServerContext *,
      const ActivateFlowsRequest *,
      ActivateFlowsResult *));
  MOCK_METHOD3(
    DeactivateFlows,
    Status(
      grpc::ServerContext *,
      const DeactivateFlowsRequest *,
      DeactivateFlowsResult *));
};

class MockPipelinedClient : public PipelinedClient {
 public:
  MockPipelinedClient()
  {
    ON_CALL(*this, setup(_,_,_)).WillByDefault(Return(true));
    ON_CALL(*this, deactivate_all_flows(_)).WillByDefault(Return(true));
    ON_CALL(*this, deactivate_flows_for_rules(_, _, _))
      .WillByDefault(Return(true));
    ON_CALL(*this, activate_flows_for_rules(_, _, _, _))
      .WillByDefault(Return(true));
    ON_CALL(*this, add_ue_mac_flow(_, _, _, _, _)).WillByDefault(Return(true));
    ON_CALL(*this, delete_ue_mac_flow(_, _)).WillByDefault(Return(true));
    ON_CALL(*this, update_subscriber_quota_state(_))
      .WillByDefault(Return(true));
  }

  MOCK_METHOD3(setup,
    bool(
      const std::vector<SessionState::SessionInfo>& infos,
      const std::uint64_t& epoch,
      std::function<void(Status status, SetupFlowsResult)> callback));
  MOCK_METHOD1(deactivate_all_flows, bool(const std::string& imsi));
  MOCK_METHOD3(
    deactivate_flows_for_rules,
    bool(
      const std::string& imsi,
      const std::vector<std::string>& rule_ids,
      const std::vector<PolicyRule>& dynamic_rules));
  MOCK_METHOD4(
    activate_flows_for_rules,
    bool(
      const std::string& imsi,
      const std::string& ip_addr,
      const std::vector<std::string>& static_rules,
      const std::vector<PolicyRule>& dynamic_rules));
  MOCK_METHOD5(
    add_ue_mac_flow,
    bool(
      const SubscriberID &sid,
      const std::string &ue_mac_addr,
      const std::string &msisdn,
      const std::string &ap_mac_addr,
      const std::string &ap_name));
  MOCK_METHOD2(
    delete_ue_mac_flow,
    bool(
      const SubscriberID &sid,
      const std::string &ue_mac_addr));
  MOCK_METHOD1(
    update_subscriber_quota_state,
    bool(const std::vector<SubscriberQuotaUpdate>& updates));
};

class MockDirectorydClient : public AsyncDirectorydClient {
 public:
  MockDirectorydClient()
  {
    ON_CALL(*this, get_directoryd_ip_field(_,_)).WillByDefault(Return(true));
  }

  MOCK_METHOD2(get_directoryd_ip_field,
    bool(
      const std::string& imsi,
      std::function<void(Status status, DirectoryField)> callback));
};

/**
 * Mock handler to mock actual request handling and just test server
 */
class MockCentralController final : public CentralSessionController::Service {
 public:
  MOCK_METHOD3(
    CreateSession,
    Status(
      grpc::ServerContext *,
      const CreateSessionRequest *,
      CreateSessionResponse *));

  MOCK_METHOD3(
    UpdateSession,
    Status(
      grpc::ServerContext *,
      const UpdateSessionRequest *,
      UpdateSessionResponse *));

  MOCK_METHOD3(
    TerminateSession,
    Status(
      grpc::ServerContext *,
      const SessionTerminateRequest *,
      SessionTerminateResponse *));
};

class MockCallback {
 public:
  MOCK_METHOD2(
    update_callback,
    void(Status status, const UpdateSessionResponse& ));
  MOCK_METHOD2(
    create_callback,
    void(Status status, const CreateSessionResponse& ));
};

/**
 * Mock handler to mock actual request handling and just test server
 */
class MockSessionHandler final : public LocalSessionManagerHandler {
 public:
  ~MockSessionHandler() {}

  MOCK_METHOD3(
    ReportRuleStats,
    void(
      grpc::ServerContext *,
      const RuleRecordTable *,
      std::function<void(Status, Void)>));

  MOCK_METHOD3(
    CreateSession,
    void(
      grpc::ServerContext *,
      const LocalCreateSessionRequest *,
      std::function<void(Status, LocalCreateSessionResponse)>));

  MOCK_METHOD3(
    EndSession,
    void(
      grpc::ServerContext *,
      const LocalEndSessionRequest *,
      std::function<void(Status, LocalEndSessionResponse)>));
};

class MockSessionReporter : public SessionReporter {
  public:
    MOCK_METHOD2(
      report_updates,
      void(
        const UpdateSessionRequest& ,
        std::function<void(grpc::Status, UpdateSessionResponse)>));

    MOCK_METHOD2(
      report_create_session,
      void(
        const CreateSessionRequest& ,
        std::function<void(Status, CreateSessionResponse)>));

    MOCK_METHOD2(
      report_terminate_session,
      void(
        const SessionTerminateRequest& ,
        std::function<void(Status, SessionTerminateResponse)>));

};

class MockAAAClient : public aaa::AAAClient {
 public:
  MockAAAClient()
  {
    ON_CALL(*this, terminate_session(_, _)).WillByDefault(Return(true));
  }

  MOCK_METHOD2(
    terminate_session,
    bool(
      const std::string& radius_session_id,
      const std::string& imsi));
};

class MockSpgwServiceClient : public SpgwServiceClient {
  public:
    MockSpgwServiceClient()
    {
      ON_CALL(*this, delete_dedicated_bearer(_, _, _, _))
        .WillByDefault(Return(true));
      ON_CALL(*this, create_dedicated_bearer(_, _, _, _))
        .WillByDefault(Return(true));
    }

    MOCK_METHOD4(
      delete_dedicated_bearer,
      bool(
        const std::string& ,
        const std::string& ,
        const uint32_t,
        const std::vector<uint32_t>& ));
    MOCK_METHOD4(
      create_dedicated_bearer,
      bool(
        const std::string& ,
        const std::string& ,
        const uint32_t,
        const std::vector<PolicyRule>& ));
};

} // namespace magma
