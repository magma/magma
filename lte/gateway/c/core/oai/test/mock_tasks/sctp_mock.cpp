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

task_zmq_ctx_t task_zmq_ctx_sctp;
static std::shared_ptr<MockSctpHandler> sctp_handler_;

void stop_mock_sctp_task();

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    case SCTP_INIT_MSG: {
    } break;

    case SCTP_CLOSE_ASSOCIATION: {
    } break;

    case SCTP_DATA_REQ: {
      sctp_handler_->sctpd_send_dl();
    } break;

    case MESSAGE_TEST: {
    } break;

    case TERMINATE_MESSAGE: {
      itti_free_msg_content(received_message_p);
      free(received_message_p);
      stop_mock_sctp_task();
    } break;

    default: {
    } break;
  }

  itti_free_msg_content(received_message_p);
  free(received_message_p);

  return 0;
}

void stop_mock_sctp_task() {
  destroy_task_context(&task_zmq_ctx_sctp);
  pthread_exit(NULL);
}

void start_mock_sctp_task(std::shared_ptr<MockSctpHandler> sctp_handler) {
  sctp_handler_ = sctp_handler;
  init_task_context(TASK_SCTP, nullptr, 0, handle_message, &task_zmq_ctx_sctp);
  zloop_start(task_zmq_ctx_sctp.event_loop);
}
