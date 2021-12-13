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
#include "lte/gateway/c/core/oai/test/mock_tasks/mock_tasks.h"

task_zmq_ctx_t task_zmq_ctx_s1ap;
static std::shared_ptr<MockS1apHandler> s1ap_handler_;

void stop_mock_s1ap_task();

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    case ACTIVATE_MESSAGE: {
    } break;

    case MESSAGE_TEST: {
    } break;

    case SCTP_DATA_IND: {
    } break;

    case SCTP_DATA_CNF: {
    } break;

    case SCTP_CLOSE_ASSOCIATION: {
    } break;

    case SCTP_NEW_ASSOCIATION: {
    } break;

    case S1AP_NAS_DL_DATA_REQ: {
      s1ap_handler_->s1ap_generate_downlink_nas_transport(
          S1AP_NAS_DL_DATA_REQ(received_message_p));
    } break;

    case S1AP_E_RAB_SETUP_REQ: {
      s1ap_handler_->s1ap_generate_s1ap_e_rab_setup_req();
    } break;

    case S1AP_E_RAB_MODIFICATION_CNF: {
    } break;

    case S1AP_UE_CONTEXT_RELEASE_COMMAND: {
      s1ap_handler_->s1ap_handle_ue_context_release_command();
    } break;

    case MME_APP_CONNECTION_ESTABLISHMENT_CNF: {
      s1ap_handler_->s1ap_handle_conn_est_cnf(bstrcpy(
          MME_APP_CONNECTION_ESTABLISHMENT_CNF(received_message_p).nas_pdu[0]));
    } break;

    case MME_APP_S1AP_MME_UE_ID_NOTIFICATION: {
    } break;

    case S1AP_ENB_INITIATED_RESET_ACK: {
    } break;

    case S1AP_PAGING_REQUEST: {
      s1ap_handler_->s1ap_handle_paging_request();
    } break;

    case S1AP_UE_CONTEXT_MODIFICATION_REQUEST: {
    } break;

    case S1AP_E_RAB_REL_CMD: {
      s1ap_handler_->s1ap_generate_s1ap_e_rab_rel_cmd();
    } break;

    case S1AP_PATH_SWITCH_REQUEST_ACK: {
      s1ap_handler_->s1ap_handle_path_switch_req_ack(
          received_message_p->ittiMsg.s1ap_path_switch_request_ack);
    } break;

    case S1AP_PATH_SWITCH_REQUEST_FAILURE: {
      s1ap_handler_->s1ap_handle_path_switch_req_failure(
          received_message_p->ittiMsg.s1ap_path_switch_request_failure);
    } break;

    case MME_APP_HANDOVER_REQUEST: {
    } break;

    case MME_APP_HANDOVER_COMMAND: {
    } break;

    case TERMINATE_MESSAGE: {
      s1ap_handler_.reset();
      itti_free_msg_content(received_message_p);
      free(received_message_p);
      stop_mock_s1ap_task();
    } break;

    default: { } break; }
  itti_free_msg_content(received_message_p);
  free(received_message_p);

  return 0;
}

void stop_mock_s1ap_task() {
  destroy_task_context(&task_zmq_ctx_s1ap);
  pthread_exit(NULL);
}

void start_mock_s1ap_task(std::shared_ptr<MockS1apHandler> s1ap_handler) {
  s1ap_handler_ = s1ap_handler;
  init_task_context(TASK_S1AP, nullptr, 0, handle_message, &task_zmq_ctx_s1ap);
  zloop_start(task_zmq_ctx_s1ap.event_loop);
}
