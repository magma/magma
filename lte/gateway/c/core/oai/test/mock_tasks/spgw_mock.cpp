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

task_zmq_ctx_t task_zmq_ctx_spgw;
static std::shared_ptr<MockSpgwHandler> spgw_handler_;

void stop_mock_spgw_task();

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    case TERMINATE_MESSAGE: {
      spgw_handler_.reset();
      itti_free_msg_content(received_message_p);
      free(received_message_p);
      stop_mock_spgw_task();
    } break;

    case S11_CREATE_SESSION_REQUEST: {
      spgw_handler_->sgw_handle_s11_create_session_request();
    } break;

    case S11_DELETE_SESSION_REQUEST: {
      spgw_handler_->sgw_handle_delete_session_request();
    } break;

    case S11_MODIFY_BEARER_REQUEST: {
      spgw_handler_->sgw_handle_modify_bearer_request();
    } break;

    case S11_RELEASE_ACCESS_BEARERS_REQUEST: {
      spgw_handler_->sgw_handle_release_access_bearers_request();
    } break;

    case S11_NW_INITIATED_ACTIVATE_BEARER_RESP: {
      // Handle Dedicated bearer Activation Rsp from MME
      spgw_handler_->sgw_handle_nw_initiated_actv_bearer_rsp();
    } break;

    case S11_NW_INITIATED_DEACTIVATE_BEARER_RESP: {
      // Handle Dedicated bearer Deactivation Rsp from MME
      spgw_handler_->sgw_handle_nw_initiated_deactv_bearer_rsp();
    } break;

    default: {
    } break;
  }
  itti_free_msg_content(received_message_p);
  free(received_message_p);

  return 0;
}

void stop_mock_spgw_task() {
  destroy_task_context(&task_zmq_ctx_spgw);
  pthread_exit(NULL);
}

void start_mock_spgw_task(std::shared_ptr<MockSpgwHandler> spgw_handler) {
  spgw_handler_ = spgw_handler;
  init_task_context(TASK_SPGW_APP, nullptr, 0, handle_message,
                    &task_zmq_ctx_spgw);
  zloop_start(task_zmq_ctx_spgw.event_loop);
}
