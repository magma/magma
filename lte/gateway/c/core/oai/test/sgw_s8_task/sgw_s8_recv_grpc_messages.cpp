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
#include "../mock_tasks/mock_tasks.h"
#include "S8ServiceImpl.h"
#include "feg/protos/s8_proxy.grpc.pb.h"
extern "C" {
#include "grpc_service.h"
}

using grpc::ServerContext;
using ::testing::Test;
task_zmq_ctx_t task_zmq_ctx_main_grpc;
static int handle_message_test_s8_grpc(
    zloop_t* loop, zsock_t* reader, void* arg);
class SgwS8MessagesTest : public ::testing::Test {
 protected:
  virtual void SetUp();
  virtual void TearDown();
  void build_grpc_create_bearer_req(magma::feg::CreateBearerRequestPgw* cb_req);
  void build_grpc_delete_bearer_req(magma::feg::DeleteBearerRequestPgw* db_req);

 protected:
  std::shared_ptr<magma::S8ServiceImpl> s8_message_receiver;
  std::shared_ptr<MockS8Handler> sgw_s8_handler;
};

void SgwS8MessagesTest::SetUp() {
  s8_message_receiver = std::make_shared<magma::S8ServiceImpl>();
  sgw_s8_handler      = std::make_shared<MockS8Handler>();
  itti_init(
      TASK_MAX, THREAD_MAX, MESSAGES_ID_MAX, tasks_info, messages_info, NULL,
      NULL);
  task_id_t task_id_list[2] = {TASK_GRPC_SERVICE, TASK_SGW_S8};
  init_task_context(
      TASK_MAIN, task_id_list, 2, handle_message_test_s8_grpc,
      &task_zmq_ctx_main_grpc);

  std::thread task_sgw_s8(start_mock_sgw_s8_task, sgw_s8_handler);
  task_sgw_s8.detach();
  grpc_service_init(TEST_GRPCSERVICES_SERVER_ADDRESS);
  std::this_thread::sleep_for(std::chrono::milliseconds(500));
}

void SgwS8MessagesTest::TearDown() {
  send_terminate_message_fatal(&task_zmq_ctx_main_grpc);
  destroy_task_context(&task_zmq_ctx_main_grpc);
  itti_free_desc_threads();
  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(1000));
}

void SgwS8MessagesTest::build_grpc_delete_bearer_req(
    magma::feg::DeleteBearerRequestPgw* db_req) {
  db_req->set_sequence_number(2);
  db_req->set_c_agw_teid(10);
  db_req->set_linked_bearer_id(5);
  db_req->add_eps_bearer_id(6);
}

void SgwS8MessagesTest::build_grpc_create_bearer_req(
    magma::feg::CreateBearerRequestPgw* cb_req) {
  cb_req->set_pgwaddrs("192.168.32.118:12342");
  cb_req->set_sequence_number(1);
  cb_req->set_c_agw_teid(10);
  cb_req->set_linked_bearer_id(5);
  return;
}

static int handle_message_test_s8_grpc(
    zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    default: { } break; }

  itti_free_msg_content(received_message_p);
  free(received_message_p);
  return 0;
}

MATCHER_P(check_s8_params_in_cb_req, cb_req, "") {
  auto cb_req_rcvd_at_sgw_s8 = static_cast<s8_create_bearer_request_t>(arg);
  if (cb_req_rcvd_at_sgw_s8.linked_eps_bearer_id != cb_req.linked_bearer_id()) {
    return false;
  }
  if (cb_req_rcvd_at_sgw_s8.context_teid != cb_req.c_agw_teid()) {
    return false;
  }
  if (cb_req_rcvd_at_sgw_s8.sequence_number != cb_req.sequence_number()) {
    return false;
  }
  if (memcmp(
          cb_req_rcvd_at_sgw_s8.pgw_cp_address, cb_req.pgwaddrs().c_str(),
          cb_req.pgwaddrs().size())) {
    return false;
  }
  return true;
}

MATCHER_P(check_s8_params_in_db_req, db_req, "") {
  auto db_req_rcvd_at_sgw_s8 = static_cast<s8_delete_bearer_request_t>(arg);
  if (db_req_rcvd_at_sgw_s8.linked_eps_bearer_id != db_req.linked_bearer_id()) {
    return false;
  }
  if (db_req_rcvd_at_sgw_s8.context_teid != db_req.c_agw_teid()) {
    return false;
  }
  if (db_req_rcvd_at_sgw_s8.sequence_number != db_req.sequence_number()) {
    return false;
  }
  if (db_req_rcvd_at_sgw_s8.eps_bearer_id[0] != db_req.eps_bearer_id(0)) {
    return false;
  }
  return true;
}

TEST_F(SgwS8MessagesTest, recv_create_bearer_req) {
  magma::feg::CreateBearerRequestPgw cb_req;
  build_grpc_create_bearer_req(&cb_req);
  grpc::ServerContext server_context;
  EXPECT_CALL(
      *sgw_s8_handler,
      sgw_s8_handle_create_bearer_request(check_s8_params_in_cb_req(cb_req)))
      .Times(1);
  magma::orc8r::Void response;
  grpc::Status status =
      s8_message_receiver->CreateBearer(&server_context, &cb_req, &response);
  EXPECT_TRUE(status.ok());
}

TEST_F(SgwS8MessagesTest, recv_grcp_delete_bearer_req) {
  magma::feg::DeleteBearerRequestPgw db_req;
  build_grpc_delete_bearer_req(&db_req);
  grpc::ServerContext server_context;
  EXPECT_CALL(
      *sgw_s8_handler,
      sgw_s8_handle_delete_bearer_request(check_s8_params_in_db_req(db_req)))
      .Times(1);
  magma::orc8r::Void response;
  grpc::Status status = s8_message_receiver->DeleteBearerRequest(
      &server_context, &db_req, &response);
  EXPECT_TRUE(status.ok());
}
