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
#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <chrono>
#include <future>
#include <memory>
#include <thread>

#include "includes/MagmaService.h"
#include "includes/ServiceRegistrySingleton.h"
#include "SessiondMocks.h"
#include "SessionReporter.h"

using grpc::Status;
using ::testing::_;
using ::testing::Test;

namespace magma {

class SessionReporterTest : public ::testing::Test {
 protected:
  /**
   * Create magma service and run in separate thread
   */
  virtual void SetUp() {
    evb = new folly::EventBase();

    auto channel = ServiceRegistrySingleton::Instance()->GetGrpcChannel(
        "test_service", ServiceRegistrySingleton::LOCAL);
    magma_service =
        std::make_shared<service303::MagmaService>("test_service", "1.0");
    mock_cloud = std::make_shared<MockCentralController>();
    magma_service->AddServiceToServer(mock_cloud.get());

    reporter = std::make_shared<SessionReporterImpl>(evb, channel);

    std::thread reporter_thread([&]() {
      std::cout << "Started reporter thread\n";
      reporter->rpc_response_loop();
      std::cout << "Stopped reporter thread\n";
    });

    std::thread cloud_thread([&]() {
      std::cout << "Started cloud thread\n";
      magma_service->Start();
      magma_service->WaitForShutdown();
      std::cout << "Stopped cloud thread\n";
    });

    // wait for server to start
    std::this_thread::sleep_for(std::chrono::milliseconds(10));
    reporter_thread.detach();
    cloud_thread.detach();
  }

  virtual void TearDown() {
    magma_service->Stop();
    reporter->stop();
    delete evb;
  }

  // Timeout to not block test
  void set_timeout(uint32_t ms) {
    std::thread([ms]() {
      std::this_thread::sleep_for(std::chrono::milliseconds(ms));
      EXPECT_TRUE(false);
    }).detach();
  }

 protected:
  std::shared_ptr<service303::MagmaService> magma_service;
  std::shared_ptr<MockCentralController> mock_cloud;
  std::shared_ptr<SessionReporter> reporter;
  folly::EventBase* evb;
  MockCallback mock_callback;
};

MATCHER_P(CheckCreateResponseRuleSize, size, "") {
  return arg.static_rules_size() == size;
}

// Test requests on single thread
TEST_F(SessionReporterTest, test_single_call) {
  EXPECT_CALL(mock_callback, create_callback(_, CheckCreateResponseRuleSize(1)))
      .Times(1);
  // add rule id for verification
  CreateSessionResponse response;
  response.mutable_static_rules()->Add()->mutable_rule_id();
  EXPECT_CALL(*mock_cloud, CreateSession(_, _, _))
      .Times(1)
      .WillOnce(testing::DoAll(testing::SetArgPointee<2>(response),
                               testing::Return(grpc::Status::OK)));

  std::promise<void> promise1;
  CreateSessionRequest request;
  reporter->report_create_session(
      request,
      [this, &promise1](Status status, CreateSessionResponse response_out) {
        mock_callback.create_callback(status, response_out);
        promise1.set_value();
      });

  // wait for one response
  std::thread([&]() {
    promise1.get_future().wait();
    evb->terminateLoopSoon();
  }).detach();

  set_timeout(1000);

  // wait for callback
  evb->loopForever();
}

// Test multiple calls at the same time, wait for all to finish
TEST_F(SessionReporterTest, test_multi_call) {
  EXPECT_CALL(mock_callback, create_callback(_, _)).Times(2);
  EXPECT_CALL(mock_callback, update_callback(_, _)).Times(1);
  CreateSessionResponse response;
  EXPECT_CALL(*mock_cloud, CreateSession(_, _, _))
      .Times(2)
      .WillRepeatedly(testing::DoAll(testing::SetArgPointee<2>(response),
                                     testing::Return(grpc::Status::OK)));
  UpdateSessionResponse update_response;
  EXPECT_CALL(*mock_cloud, UpdateSession(_, _, _))
      .Times(1)
      .WillRepeatedly(testing::DoAll(testing::SetArgPointee<2>(update_response),
                                     testing::Return(grpc::Status::OK)));

  std::promise<void> promise1, promise2, promise3;

  reporter->report_create_session(
      CreateSessionRequest(),
      [&](Status status, CreateSessionResponse response_out) {
        mock_callback.create_callback(status, response_out);
        promise1.set_value();
      });
  reporter->report_updates(
      UpdateSessionRequest(),
      [&](Status status, UpdateSessionResponse response_out) {
        mock_callback.update_callback(status, response_out);
        promise2.set_value();
      });
  reporter->report_create_session(
      CreateSessionRequest(),
      [&](Status status, CreateSessionResponse response_out) {
        mock_callback.create_callback(status, response_out);
        promise3.set_value();
      });

  // wait for all 3 responses
  std::thread([&]() {
    promise1.get_future().wait();
    promise2.get_future().wait();
    promise3.get_future().wait();
    evb->terminateLoopSoon();
  }).detach();

  set_timeout(1000);

  // wait for callback
  evb->loopForever();
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace magma
