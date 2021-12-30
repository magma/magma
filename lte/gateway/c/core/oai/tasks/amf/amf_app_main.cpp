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

#include <iostream>
#include <sstream>
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/common/itti_free_defined_msg.h"
#include "lte/gateway/c/core/oai/lib/message_utils/service303_message_utils.h"
#include "lte/gateway/c/core/oai/include/amf_config.h"
#include "lte/gateway/c/core/oai/include/amf_as_message.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/include/amf_app_messages_types.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_fsm.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_data.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_defs.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_authentication.h"
#include "lte/gateway/c/core/oai/include/ngap_messages_types.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_state_manager.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_smfDefs.h"
#include "lte/gateway/c/core/oai/common/common_defs.h"

namespace magma5g {
task_zmq_ctx_t amf_app_task_zmq_ctx;
void amf_app_exit(void);

/****************************************************************************
 **                                                                        **
 ** Name:    handle_message()                                              **
 **                                                                        **
 ** Description: Handle Uplink UE messages                                 **
 **                                                                        **
 ** Inputs:  loop:    Read the packets in loop                             **
 **      reader:    Read the packets from other thread                     **
 **      arg:       Argument                                               **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);
  amf_app_desc_t* amf_app_desc_p = get_amf_nas_state(false);
  // imsi64_t imsi64                =
  // itti_get_associated_imsi(received_message_p);

  OAILOG_INFO(
      LOG_AMF_APP, "Received msg from :[%s] id:[%d] name:[%s]\n",
      ITTI_MSG_ORIGIN_NAME(received_message_p), ITTI_MSG_ID(received_message_p),
      ITTI_MSG_NAME(received_message_p));

  switch (ITTI_MSG_ID(received_message_p)) {
    /* Handle Initial UE message from NGAP */
    case NGAP_INITIAL_UE_MESSAGE:
      amf_app_handle_initial_ue_message(
          amf_app_desc_p, &NGAP_INITIAL_UE_MESSAGE(received_message_p));
      break;
    /* Handle uplink NAS message Recevied from the UE */
    case AMF_APP_UPLINK_DATA_IND:
      amf_app_handle_uplink_nas_message(
          amf_app_desc_p, AMF_APP_UL_DATA_IND(received_message_p).nas_msg,
          AMF_APP_UL_DATA_IND(received_message_p).ue_id,
          AMF_APP_UL_DATA_IND(received_message_p).tai);
      break;
    case AMF_APP_INITIAL_CONTEXT_SETUP_RSP:
      amf_app_handle_initial_context_setup_rsp(
          amf_app_desc_p,
          &AMF_APP_INITIAL_CONTEXT_SETUP_RSP(received_message_p));
      break;
    /* Handle PDU session Response from UE */
    case N11_CREATE_PDU_SESSION_RESPONSE:
      amf_app_handle_pdu_session_response(
          &N11_CREATE_PDU_SESSION_RESPONSE(received_message_p));
      break;
    case AMF_APP_SUBS_AUTH_INFO_RESP:
      amf_nas_proc_authentication_info_answer(
          &AMF_APP_AUTH_RESPONSE_DATA(received_message_p));
      break;

    case AMF_APP_DECRYPT_IMSI_INFO_RESP:
      amf_decrypt_imsi_info_answer(
          &AMF_APP_DECRYPT_IMSI_RESPONSE_DATA(received_message_p));
      break;

    case AMF_IP_ALLOCATION_RESPONSE:
      itti_amf_ip_allocation_response_t* response_p;
      response_p = &(received_message_p->ittiMsg.amf_ip_allocation_response);
      amf_smf_handle_ip_address_response(response_p);
      break;

    case S6A_UPDATE_LOCATION_ANS: {
      OAILOG_INFO(
          LOG_MME_APP,
          "Received S6A Update Location Answer from subscriberd\n");
      amf_handle_s6a_update_location_ans(
          &received_message_p->ittiMsg.s6a_update_location_ans);
    } break;

    /* Handle PDU session resource setup response */
    case NGAP_PDUSESSIONRESOURCE_SETUP_RSP:
      /* This is non-nas message and can be handled directly to check if failure
       * or success messages are coming from NGAP
       */
      amf_app_handle_resource_setup_response(
          NGAP_PDUSESSIONRESOURCE_SETUP_RSP(received_message_p));
      break;
    /* Handle PDU session resource release response */
    case NGAP_PDUSESSIONRESOURCE_REL_RSP:
      /* This is non-nas message and can be handled directly to check if failure
       * or success messages are coming from NGAP
       */
      amf_app_handle_resource_release_response(
          NGAP_PDUSESSIONRESOURCE_REL_RSP(received_message_p));
      break;
    case N11_NOTIFICATION_RECEIVED:
      /* This case handles Notification Received for Paging or other events
       * or success messages are coming from NGAP
       */
      amf_app_handle_notification_received(
          &N11_NOTIFICATION_RECEIVED(received_message_p));
      break;

    /* Handle UE context Release Requests */
    case NGAP_UE_CONTEXT_RELEASE_REQ:
      /* This is non-nas message and handled directly from NGAP sent to AMF
       * on RRC-Inactive mode to change UE's CM-connected to CM-idle state.
       */
      amf_app_handle_cm_idle_on_ue_context_release(
          NGAP_UE_CONTEXT_RELEASE_REQ(received_message_p));
      break;
    /* Handle Terminate message */
    case TERMINATE_MESSAGE:
      itti_free_msg_content(received_message_p);
      free(received_message_p);
      amf_app_exit();
      break;
    default:
      OAILOG_DEBUG(LOG_AMF_APP, "default message received");
      break;
  }
  return RETURNok;
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_app_thread()                                              **
 **                                                                        **
 ** Description: Launching of the amf Thread                               **
 **                                                                        **
 ** Inputs:  args: arguments                                               **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
void* amf_app_thread(void* args) {
  itti_mark_task_ready(TASK_AMF_APP);
  const task_id_t tasks[] = {TASK_NGAP, TASK_SERVICE303};
  init_task_context(
      TASK_AMF_APP, tasks, 2, handle_message, &amf_app_task_zmq_ctx);
  // Service started, but not healthy yet
  send_app_health_to_service303(&amf_app_task_zmq_ctx, TASK_AMF_APP, false);
  zloop_start(amf_app_task_zmq_ctx.event_loop);
  amf_app_exit();
  return NULL;
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_app_init()                                                **
 **                                                                        **
 ** Description: Initialisation of amf application thread                  **
 **              based on the configurations                               **
 **                                                                        **
 ** Inputs:  amf_config_p: amf configuration read from the file            **
 **                                                                        **
 ** Return:    RETURNok, RETURNerror                                       **
 ***************************************************************************/
extern "C" int amf_app_init(const amf_config_t* amf_config_p) {
  if (amf_nas_state_init(amf_config_p)) {
    OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
  }
  /*Initialize UE state matrix */
  create_state_matrix();
  if (itti_create_task(TASK_AMF_APP, &amf_app_thread, NULL) < 0) {
    OAILOG_CRITICAL(LOG_AMF_APP, "Amf app create task failed\n");
    OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
  }
  OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNok);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_app_exit()                                                **
 **                                                                        **
 ** Description: Exit the amf app thread resources                         **
 **              allocated during the initialisation                       **
 **                                                                        **
 ** Inputs:  void: no arguments                                            **
 **                                                                        **
 ***************************************************************************/
void amf_app_exit(void) {
  clear_amf_nas_state();
  amf_config_exit();
  destroy_task_context(&amf_app_task_zmq_ctx);
  OAI_FPRINTF_INFO("TASK_AMF_APP terminated\n");
  pthread_exit(NULL);
}
}  // namespace magma5g
