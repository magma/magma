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

task_zmq_ctx_t task_zmq_ctx_s6a;
static std::shared_ptr<MockS6aHandler> s6a_handler_;

void stop_mock_s6a_task();

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    case TERMINATE_MESSAGE: {
      s6a_handler_.reset();
      itti_free_msg_content(received_message_p);
      free(received_message_p);
      stop_mock_s6a_task();
    } break;
    case S6A_AUTH_INFO_REQ: {
      s6a_handler_->s6a_viface_authentication_info_req();
    } break;
    case S6A_UPDATE_LOCATION_REQ: {
      s6a_handler_->s6a_viface_update_location_req();
    } break;
    case S6A_PURGE_UE_REQ: {
      s6a_handler_->s6a_viface_purge_ue();
    } break;
    case S6A_CANCEL_LOCATION_ANS: {
      s6a_handler_->s6a_cancel_location_ans();
    } break;
    default: {
    } break;
  }
  itti_free_msg_content(received_message_p);
  free(received_message_p);

  return 0;
}

void stop_mock_s6a_task() {
  destroy_task_context(&task_zmq_ctx_s6a);
  pthread_exit(NULL);
}

void start_mock_s6a_task(std::shared_ptr<MockS6aHandler> s6a_handler) {
  s6a_handler_ = s6a_handler;
  init_task_context(TASK_S6A, nullptr, 0, handle_message, &task_zmq_ctx_s6a);
  zloop_start(task_zmq_ctx_s6a.event_loop);
}
