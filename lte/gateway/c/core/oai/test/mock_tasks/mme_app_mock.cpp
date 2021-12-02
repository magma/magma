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

#include "lte/gateway/c/core/oai/include/mme_app_messages_types.h"

task_zmq_ctx_t task_zmq_ctx_mme;
static std::shared_ptr<MockMmeAppHandler> mme_app_handler_;

void stop_mock_mme_app_task();

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    case MME_APP_INITIAL_CONTEXT_SETUP_RSP: {
    } break;

    case S6A_CANCEL_LOCATION_REQ: {
    } break;

    case MME_APP_UPLINK_DATA_IND: {
    } break;

    case S11_CREATE_BEARER_REQUEST: {
    } break;

    case S6A_RESET_REQ: {
    } break;

    case S11_CREATE_SESSION_RESPONSE: {
      mme_app_handler_->mme_app_handle_create_sess_resp();
    } break;

    case S11_MODIFY_BEARER_RESPONSE: {
      mme_app_handler_->mme_app_handle_modify_bearer_rsp();
    } break;

    case S11_RELEASE_ACCESS_BEARERS_RESPONSE: {
      mme_app_handler_->mme_app_handle_release_access_bearers_resp();
    } break;

    case S11_DELETE_SESSION_RESPONSE: {
      mme_app_handler_->mme_app_handle_delete_sess_rsp();
    } break;

    case S11_SUSPEND_ACKNOWLEDGE: {
    } break;

    case S1AP_E_RAB_SETUP_RSP: {
      mme_app_handler_->mme_app_handle_e_rab_setup_rsp();
    } break;

    case S1AP_E_RAB_REL_RSP: {
    } break;

    case S1AP_E_RAB_MODIFICATION_IND: {
    } break;

    case S1AP_INITIAL_UE_MESSAGE: {
      mme_app_handler_->mme_app_handle_initial_ue_message();
      bdestroy_wrapper(&S1AP_INITIAL_UE_MESSAGE(received_message_p).nas);
    } break;

    case S6A_UPDATE_LOCATION_ANS: {
    } break;

    case S1AP_ENB_INITIATED_RESET_REQ: {
      mme_app_handler_->mme_app_handle_enb_reset_req();
      free_wrapper((void**) &S1AP_ENB_INITIATED_RESET_REQ(received_message_p)
                       .ue_to_reset_list);
    } break;

    case S11_PAGING_REQUEST: {
    } break;

    case MME_APP_INITIAL_CONTEXT_SETUP_FAILURE: {
      mme_app_handler_->mme_app_handle_initial_context_setup_failure();
    } break;

    case S1AP_UE_CAPABILITIES_IND: {
    } break;

    case S1AP_UE_CONTEXT_RELEASE_REQ: {
      mme_app_handler_->mme_app_handle_s1ap_ue_context_release_req();
    } break;

    case S1AP_UE_CONTEXT_MODIFICATION_RESPONSE: {
    } break;

    case S1AP_UE_CONTEXT_MODIFICATION_FAILURE: {
    } break;

    case S1AP_UE_CONTEXT_RELEASE_COMPLETE: {
    } break;

    case S1AP_ENB_DEREGISTERED_IND: {
    } break;

    case ACTIVATE_MESSAGE: {
    } break;

    case SCTP_MME_SERVER_INITIALIZED: {
    } break;

    case S6A_PURGE_UE_ANS: {
    } break;

    case SGSAP_LOCATION_UPDATE_ACC: {
    } break;

    case SGSAP_LOCATION_UPDATE_REJ: {
    } break;

    case SGSAP_ALERT_REQUEST: {
    } break;

    case SGSAP_VLR_RESET_INDICATION: {
    } break;

    case SGSAP_PAGING_REQUEST: {
    } break;

    case SGSAP_SERVICE_ABORT_REQ: {
    } break;

    case SGSAP_EPS_DETACH_ACK: {
    } break;

    case SGSAP_IMSI_DETACH_ACK: {
    } break;

    case S11_MODIFY_UE_AMBR_REQUEST: {
    } break;

    case S11_NW_INITIATED_ACTIVATE_BEARER_REQUEST: {
      mme_app_handler_->mme_app_handle_nw_init_ded_bearer_actv_req(
          received_message_p->ittiMsg.s11_nw_init_actv_bearer_request);
    } break;

    case SGSAP_STATUS: {
    } break;

    case S11_NW_INITIATED_DEACTIVATE_BEARER_REQUEST: {
      mme_app_handler_->mme_app_handle_nw_init_bearer_deactv_req(
          received_message_p->ittiMsg.s11_nw_init_deactv_bearer_request);
    } break;

    case S1AP_PATH_SWITCH_REQUEST: {
      mme_app_handler_->mme_app_handle_path_switch_request();
    } break;

    case S1AP_HANDOVER_REQUIRED: {
      mme_app_handler_->mme_app_handle_handover_required();
    } break;

    case S1AP_HANDOVER_REQUEST_ACK: {
    } break;

    case S1AP_HANDOVER_NOTIFY: {
    } break;

    case S6A_AUTH_INFO_ANS: {
    } break;

    case MME_APP_DOWNLINK_DATA_CNF: {
    } break;

    case MME_APP_DOWNLINK_DATA_REJ: {
      mme_app_handler_->nas_proc_dl_transfer_rej();
      bdestroy_wrapper(&MME_APP_DL_DATA_REJ(received_message_p).nas_msg);
    } break;

    case SGSAP_DOWNLINK_UNITDATA: {
    } break;

    case SGSAP_RELEASE_REQ: {
    } break;

    case SGSAP_MM_INFORMATION_REQ: {
    } break;

    case S1AP_REMOVE_STALE_UE_CONTEXT: {
    } break;

    case RECOVERY_MESSAGE: {
    } break;

    case TERMINATE_MESSAGE: {
      mme_app_handler_.reset();
      itti_free_msg_content(received_message_p);
      free(received_message_p);
      stop_mock_mme_app_task();
    } break;

    default: { } break; }

  itti_free_msg_content(received_message_p);
  free(received_message_p);

  return 0;
}

void stop_mock_mme_app_task() {
  destroy_task_context(&task_zmq_ctx_mme);
  pthread_exit(NULL);
}

void start_mock_mme_app_task(
    std::shared_ptr<MockMmeAppHandler> mme_app_handler) {
  mme_app_handler_ = mme_app_handler;
  init_task_context(
      TASK_MME_APP, nullptr, 0, handle_message, &task_zmq_ctx_mme);
  zloop_start(task_zmq_ctx_mme.event_loop);
}
