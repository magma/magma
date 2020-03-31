/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <memory>

#include <glog/logging.h>
#include <gtest/gtest.h>
#include <folly/io/async/EventBaseManager.h>

#include "MagmaService.h"
#include "LocalEnforcer.h"
#include "ProtobufCreators.h"
#include "RuleStore.h"
#include "ServiceRegistrySingleton.h"
#include "SessiondMocks.h"
#include "SessionStore.h"
#include "magma_logging.h"


using ::testing::Test;

namespace magma {

class SessionManagerHandlerTest : public ::testing::Test {
  protected:
  protected:
    virtual void SetUp() {
        reporter = std::make_shared<MockSessionReporter>();
        auto rule_store = std::make_shared<StaticRuleStore>();
        auto session_store = new SessionStore(rule_store);
        auto pipelined_client = std::make_shared<MockPipelinedClient>();
        auto directoryd_client = std::make_shared<MockDirectorydClient>();
        auto eventd_client = std::make_shared<MockEventdClient>();
        auto spgw_client = std::make_shared<MockSpgwServiceClient>();
        auto aaa_client = std::make_shared<MockAAAClient>();
        local_enforcer = std::make_shared<LocalEnforcer>(
                reporter,
                rule_store,
                pipelined_client,
                directoryd_client,
                eventd_client,
                spgw_client,
                aaa_client,
                0,
                0);
        evb = folly::EventBaseManager::get()->getEventBase();
        local_enforcer->attachEventBase(evb);
      session_map = SessionMap{};

      session_manager = std::make_shared<LocalSessionManagerHandlerImpl>(
        local_enforcer, reporter.get(), directoryd_client, session_map, *session_store);
    }

  protected:
    std::shared_ptr<LocalSessionManagerHandlerImpl> session_manager;
    std::shared_ptr<MockSessionReporter> reporter;
    std::shared_ptr <LocalEnforcer> local_enforcer;
    SessionIDGenerator id_gen_;
    folly::EventBase *evb;
    SessionMap session_map;
};

TEST_F(SessionManagerHandlerTest, test_create_session_cfg)
{
    LocalCreateSessionRequest request;
    CreateSessionResponse response;
    std::string hardware_addr_bytes = {0x0f,0x10,0x2e,0x12,0x3a,0x55};
    std::string imsi = "IMSI1";
    std::string msisdn = "5100001234";
    std::string radius_session_id = "AA-AA-AA-AA-AA-AA:TESTAP__"
                                    "0F-10-2E-12-3A-55";
    auto sid = id_gen_.gen_session_id(imsi);
    SessionState::Config cfg = {.ue_ipv4 = "",
            .spgw_ipv4 = "",
            .msisdn = msisdn,
            .apn = "apn1",
            .imei = "",
            .plmn_id = "",
            .imsi_plmn_id = "",
            .user_location = "",
            .rat_type = RATType::TGPP_WLAN,
            .mac_addr = "0f:10:2e:12:3a:55",
            .hardware_addr = hardware_addr_bytes,
            .radius_session_id = radius_session_id};

    local_enforcer->init_session_credit(session_map, imsi, sid, cfg, response);

    grpc::ServerContext create_context;
    request.mutable_sid()->set_id("IMSI1");
    request.set_rat_type(RATType::TGPP_WLAN);
    request.set_hardware_addr(hardware_addr_bytes);
    request.set_msisdn(msisdn);
    request.set_radius_session_id(radius_session_id);

    // Ensure session is not reported as its a duplicate
    EXPECT_CALL(*reporter, report_create_session(_, _)).Times(0);
    session_manager->CreateSession(&create_context, &request, [this](
            grpc::Status status, LocalCreateSessionResponse response_out) {});
}

int main(int argc, char **argv)
{
    ::testing::InitGoogleTest(&argc, argv);
    return RUN_ALL_TESTS();
}

} // namespace magma
