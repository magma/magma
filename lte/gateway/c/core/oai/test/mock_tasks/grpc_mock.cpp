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

#include <grpcpp/grpcpp.h>
#include <grpcpp/security/server_credentials.h>
#include "lte/gateway/c/core/oai/include/grpc_service.hpp"
#include "lte/gateway/c/core/oai/tasks/grpc_service/S8ServiceImpl.hpp"
#include "lte/gateway/c/core/oai/test/mock_tasks/mock_tasks.hpp"

task_zmq_ctx_t task_zmq_ctx_grpc;
grpc_service_data_t grpc_service_config = {0};
using grpc::InsecureServerCredentials;
using grpc::Server;
using grpc::ServerBuilder;
using magma::S8ServiceImpl;

static S8ServiceImpl s8_service;
static std::unique_ptr<Server> server;
static void stop_mock_grpc_task();
static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    case TERMINATE_MESSAGE: {
      itti_free_msg_content(received_message_p);
      free(received_message_p);
      stop_mock_grpc_task();
    } break;

    default: {
    } break;
  }
  itti_free_msg_content(received_message_p);
  free(received_message_p);

  return 0;
}

static void stop_mock_grpc_task() {
  bdestroy_wrapper(&grpc_service_config.server_address);
  server->Shutdown();
  server->Wait();
  destroy_task_context(&task_zmq_ctx_grpc);
  pthread_exit(NULL);
}

static void start_grpc_s8_service(bstring server_address) {
  ServerBuilder builder;
  builder.AddListeningPort(bdata(server_address),
                           grpc::InsecureServerCredentials());
  builder.RegisterService(&s8_service);
  server = builder.BuildAndStart();
}

void start_mock_grpc_task() {
  grpc_service_config.server_address =
      bfromcstr(TEST_GRPCSERVICES_SERVER_ADDRESS);

  init_task_context(TASK_GRPC_SERVICE, nullptr, 0, handle_message,
                    &task_zmq_ctx_grpc);
  start_grpc_s8_service(grpc_service_config.server_address);
  zloop_start(task_zmq_ctx_grpc.event_loop);
}
