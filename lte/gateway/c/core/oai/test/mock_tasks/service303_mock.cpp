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

task_zmq_ctx_t task_zmq_ctx_service303;
static std::shared_ptr<MockService303Handler> service303_handler_;

void stop_mock_service303_task();

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    case TERMINATE_MESSAGE: {
      service303_handler_.reset();
      itti_free_msg_content(received_message_p);
      free(received_message_p);
      stop_mock_service303_task();
    } break;

    case APPLICATION_HEALTHY_MSG: {
      service303_handler_->service303_set_application_health();
    } break;

    default: {
    } break;
  }
  itti_free_msg_content(received_message_p);
  free(received_message_p);

  return 0;
}

void stop_mock_service303_task() {
  destroy_task_context(&task_zmq_ctx_service303);
  pthread_exit(NULL);
}

void start_mock_service303_task(
    std::shared_ptr<MockService303Handler> service303_handler) {
  service303_handler_ = service303_handler;
  init_task_context(TASK_SERVICE303, nullptr, 0, handle_message,
                    &task_zmq_ctx_service303);
  zloop_start(task_zmq_ctx_service303.event_loop);
}
