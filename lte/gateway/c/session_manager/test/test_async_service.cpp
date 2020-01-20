/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include <memory>
#include <chrono>
#include <thread>
#include <future>

#include <gtest/gtest.h>
#include <gmock/gmock.h>

#include "MagmaService.h"
#include "ServiceRegistrySingleton.h"
#include "SessiondMocks.h"
#include "SessionManagerServer.h"

using grpc::Status;
using ::testing::Test;

namespace magma {

class AsyncServiceTest : public ::testing::Test {
 protected:
  /**
   * Create magma service and run in separate thread
   */
  virtual void SetUp()
  {
    auto channel = ServiceRegistrySingleton::Instance()->GetGrpcChannel(
      "test_service", ServiceRegistrySingleton::LOCAL);
    magma_service =
      std::make_shared<service303::MagmaService>("test_service", "1.0");

    auto mock_handler_p = std::make_unique<MockSessionHandler>();
    mock_handler = mock_handler_p.get();

    async_service = std::make_shared<LocalSessionManagerAsyncService>(
      magma_service->GetNewCompletionQueue(), std::move(mock_handler_p));
    magma_service->AddServiceToServer(async_service.get());

    stub = LocalSessionManager::NewStub(channel);

    std::thread grpc_thread([&]() {
      std::cout << "Started grpc thread\n";
      magma_service->Start();
      async_service->wait_for_requests(); // block here instead of on server
      std::cout << "Stopped grpc thread\n";
    });
    // wait for server to start
    std::this_thread::sleep_for(std::chrono::milliseconds(100));
    grpc_thread.detach();
  }

  virtual void TearDown()
  {
    magma_service->Stop();
    async_service->stop();
  }

  void create_session()
  {
    grpc::ClientContext create_context;
    LocalCreateSessionResponse create_resp;
    LocalCreateSessionRequest request;
    request.mutable_sid()->set_id("IMSI1");
    auto status = stub->CreateSession(&create_context, request, &create_resp);
    EXPECT_TRUE(status.ok());
  }

  void end_session()
  {
    grpc::ClientContext end_context;
    LocalEndSessionResponse end_resp;
    LocalEndSessionRequest request;
    request.mutable_sid()->set_id("IMSI1");
    auto status = stub->EndSession(&end_context, request, &end_resp);
    EXPECT_TRUE(status.ok());
  }

  void report_rule_stats()
  {
    grpc::ClientContext create_context;
    Void void_resp;
    RuleRecordTable table;
    table.mutable_records()->Add();
    stub->ReportRuleStats(&create_context, table, &void_resp);
  }

 protected:
  std::shared_ptr<service303::MagmaService> magma_service;
  std::shared_ptr<LocalSessionManagerAsyncService> async_service;
  MockSessionHandler *mock_handler;
  std::unique_ptr<LocalSessionManager::Stub> stub;
};

// Test requests on single thread
TEST_F(AsyncServiceTest, test_single_thread)
{
  LocalCreateSessionResponse create_response;
  EXPECT_CALL(*mock_handler, CreateSession(testing::_, testing::_, testing::_))
    .Times(3)
    .WillRepeatedly(testing::InvokeArgument<2>(Status::OK, create_response));

  LocalEndSessionResponse end_resp;
  EXPECT_CALL(*mock_handler, EndSession(testing::_, testing::_, testing::_))
    .Times(1)
    .WillRepeatedly(testing::InvokeArgument<2>(Status::OK, end_resp));

  Void void_resp;
  EXPECT_CALL(
    *mock_handler, ReportRuleStats(testing::_, testing::_, testing::_))
    .Times(1)
    .WillRepeatedly(testing::InvokeArgument<2>(Status::OK, void_resp));

  create_session();
  create_session();
  create_session();
  end_session();
  report_rule_stats();
}

// Test multiple requests on multiple threads
TEST_F(AsyncServiceTest, test_multi_thread)
{
  LocalCreateSessionResponse response;
  EXPECT_CALL(*mock_handler, CreateSession(testing::_, testing::_, testing::_))
    .Times(3)
    .WillRepeatedly(testing::InvokeArgument<2>(Status::OK, response));
  Void void_resp;
  EXPECT_CALL(
    *mock_handler, ReportRuleStats(testing::_, testing::_, testing::_))
    .Times(3)
    .WillRepeatedly(testing::InvokeArgument<2>(Status::OK, void_resp));

  std::promise<void> ready_promise, t1_ready_promise, t2_ready_promise;
  std::shared_future<void> ready_future(ready_promise.get_future());
  auto fun1 = [&, ready_future]() {
    t1_ready_promise.set_value();
    ready_future.wait();
    create_session();
    create_session();
    report_rule_stats();
  };

  auto fun2 = [&, ready_future]() {
    t2_ready_promise.set_value();
    ready_future.wait();
    create_session();
    report_rule_stats();
    report_rule_stats();
  };
  auto result1 = std::async(std::launch::async, fun1);
  auto result2 = std::async(std::launch::async, fun2);

  t1_ready_promise.get_future().wait(); // wait until threads ready
  t2_ready_promise.get_future().wait();
  ready_promise.set_value();
  result1.get();
  result2.get();
}

int main(int argc, char **argv)
{
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

} // namespace magma
