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

#include <gtest/gtest.h>
#include <thread>

#include "lte/gateway/c/core/oai/test/mock_tasks/mock_tasks.hpp"
#include "feg/protos/s6a_proxy.grpc.pb.h"
#include "lte/gateway/c/core/oai/tasks/async_grpc_service/grpc_async_service_task.hpp"
#include "lte/gateway/c/core/oai/include/grpc_service.hpp"

#include "lte/gateway/c/core/oai/include/mme_config.hpp"
#include "lte/gateway/c/core/oai/tasks/s6a/s6a_defs.hpp"

using grpc::ServerContext;
using ::testing::Test;
task_zmq_ctx_t task_zmq_ctx_main_s6a;
struct mme_config_s mme_config = {.rw_lock = PTHREAD_RWLOCK_INITIALIZER, 0};

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg);

class S6aMessagesTest : public ::testing::Test {
 protected:
  virtual void SetUp();
  virtual void TearDown();
  void build_grpc_cancel_location_req(magma::feg::CancelLocationRequest* clr);
  void build_grpc_reset_req(magma::feg::ResetRequest* reset);

 protected:
  std::shared_ptr<MockMmeAppHandler> mme_app_handler;
  std::shared_ptr<magma::S6aProxyAsyncResponderHandler> async_service_handler;
};

void S6aMessagesTest::SetUp() {
  mme_app_handler = std::make_shared<MockMmeAppHandler>();
  async_service_handler =
      std::make_shared<magma::S6aProxyAsyncResponderHandler>();
  itti_init(TASK_MAX, THREAD_MAX, MESSAGES_ID_MAX, tasks_info, messages_info,
            NULL, NULL);
  task_id_t task_id_list[2] = {TASK_S6A, TASK_MME_APP};
  init_task_context(TASK_MAIN, task_id_list, 2, handle_message,
                    &task_zmq_ctx_main_s6a);

  std::thread task_mme_app(start_mock_mme_app_task, mme_app_handler);
  task_mme_app.detach();

  s6a_init(&mme_config);
  std::this_thread::sleep_for(std::chrono::milliseconds(250));
}

void S6aMessagesTest::TearDown() {
  send_terminate_message_fatal(&task_zmq_ctx_main_s6a);
  destroy_task_context(&task_zmq_ctx_main_s6a);
  itti_free_desc_threads();
  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(1000));
}

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    default: {
    } break;
  }

  itti_free_msg_content(received_message_p);
  free(received_message_p);
  return 0;
}

void S6aMessagesTest::build_grpc_cancel_location_req(
    magma::feg::CancelLocationRequest* clr) {
  clr->set_user_name("1010000000001");
  clr->set_cancellation_type(
      magma::feg::CancelLocationRequest::SUBSCRIPTION_WITHDRAWAL);
  return;
}

void S6aMessagesTest::build_grpc_reset_req(magma::feg::ResetRequest* reset) {
  reset->add_user_id("1010000000001");
  return;
}

TEST_F(S6aMessagesTest, recv_cancel_location_req) {
  magma::feg::CancelLocationRequest clr;
  build_grpc_cancel_location_req(&clr);
  grpc::ServerContext server_context;
  EXPECT_CALL(*mme_app_handler, mme_app_handle_s6a_cancel_location_req())
      .Times(1);
  async_service_handler->CancelLocation(&server_context, &clr, nullptr);
}

TEST_F(S6aMessagesTest, recv_reset_req) {
  magma::feg::ResetRequest reset;
  build_grpc_reset_req(&reset);
  grpc::ServerContext server_context;
  EXPECT_CALL(*mme_app_handler, mme_app_handle_s6a_reset_req()).Times(1);
  async_service_handler->Reset(&server_context, &reset, nullptr);
}
