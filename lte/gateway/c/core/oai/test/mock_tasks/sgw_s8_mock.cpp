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
#include "lte/gateway/c/core/oai/test/mock_tasks/mock_tasks.hpp"

task_zmq_ctx_t task_zmq_ctx_sgw_s8;
static std::shared_ptr<MockS8Handler> sgw_s8_handler_;

void stop_mock_sgw_s8_task();

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    case TERMINATE_MESSAGE: {
      sgw_s8_handler_.reset();
      itti_free_msg_content(received_message_p);
      free(received_message_p);
      stop_mock_sgw_s8_task();
    } break;
    case S8_CREATE_BEARER_REQ: {
      sgw_s8_handler_->sgw_s8_handle_create_bearer_request(
          received_message_p->ittiMsg.s8_create_bearer_req);
      free_wrapper((void**)&received_message_p->ittiMsg.s8_create_bearer_req
                       .pgw_cp_address);
    } break;
    case S8_DELETE_BEARER_REQ: {
      sgw_s8_handler_->sgw_s8_handle_delete_bearer_request(
          received_message_p->ittiMsg.s8_delete_bearer_req);
    } break;

    default: {
    } break;
  }
  itti_free_msg_content(received_message_p);
  free(received_message_p);

  return 0;
}

void stop_mock_sgw_s8_task() {
  destroy_task_context(&task_zmq_ctx_sgw_s8);
  pthread_exit(NULL);
}

void start_mock_sgw_s8_task(std::shared_ptr<MockS8Handler> sgw_s8_handler) {
  sgw_s8_handler_ = sgw_s8_handler;
  init_task_context(TASK_SGW_S8, nullptr, 0, handle_message,
                    &task_zmq_ctx_sgw_s8);
  zloop_start(task_zmq_ctx_sgw_s8.event_loop);
}
