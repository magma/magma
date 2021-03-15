
/****************************************************************************
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
 ****************************************************************************/
/*****************************************************************************

  Source      amf_app_main.cpp

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <iostream>
#include <sstream>
#ifdef __cplusplus
extern "C" {
#endif
#include "log.h"
#include "intertask_interface_types.h"
#include "intertask_interface.h"
#include "itti_free_defined_msg.h"
#include "service303_message_utils.h"
#include "amf_as_message.h"
#ifdef __cplusplus
}
#endif
#include "amf_app_messages_types.h"
#include "amf_config.h"
#include "amf_data.h"
#include "amf_fsm.h"
#include "amf_app_ue_context_and_proc.h"
#include "amf_app_defs.h"
#include "ngap_messages_types.h"
#include "amf_app_state_manager.h"
using namespace std;

namespace magma5g {
task_zmq_ctx_t amf_app_task_zmq_ctx;
void amf_app_exit(void);
static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  zframe_t* msg_frame = zframe_recv(reader);
  amf_app_defs amf_defs;
  assert(msg_frame);
  MessageDef* received_message_p = (MessageDef*) zframe_data(msg_frame);
  imsi64_t imsi64                = itti_get_associated_imsi(received_message_p);
  amf_app_desc_t* amf_app_desc_p = get_amf_nas_state(false);

  switch (ITTI_MSG_ID(received_message_p)) {
    case NGAP_INITIAL_UE_MESSAGE:
      OAILOG_INFO(LOG_AMF_APP, "AMF_TEST: NGAP_INITIAL_UE_MESSAGE received\n");
      amf_defs.amf_app_handle_initial_ue_message(
          amf_app_desc_p, &NGAP_INITIAL_UE_MESSAGE(received_message_p));
      break;

      //[authentication response, Identy response, security mode complete,
      // registration complete ]
      // case AMF_APP_UL_DATA_IND:
    case AMF_APP_UPLINK_DATA_IND:
      OAILOG_INFO(LOG_AMF_APP, "AMF_TEST: UPLINK_NAS_MESSAGE received\n");
      amf_defs.amf_app_handle_uplink_nas_message(
          amf_app_desc_p, AMF_APP_UL_DATA_IND(received_message_p).nas_msg);
      break;
    case N11_CREATE_PDU_SESSION_RESPONSE:
      OAILOG_INFO(
          LOG_AMF_APP, "AMF_TEST: session created for imsi:%s with IP:%s \n",
          N11_CREATE_PDU_SESSION_RESPONSE(received_message_p).imsi,
          N11_CREATE_PDU_SESSION_RESPONSE(received_message_p)
              .pdu_address.redirect_server_address);
      for (int i = 0; (N11_CREATE_PDU_SESSION_RESPONSE(received_message_p)
                           .pdu_address.redirect_server_address != '\0') &&
                      (i < 6);
           i++) {
        OAILOG_INFO(
            LOG_AMF_APP, "AMF_TEST: IP:%x \n",
            N11_CREATE_PDU_SESSION_RESPONSE(received_message_p)
                .pdu_address.redirect_server_address[i]);
      }
      amf_defs.amf_app_handle_pdu_session_response(
          &N11_CREATE_PDU_SESSION_RESPONSE(received_message_p));
      break;
#if 0  // TODO -  NEED-RECHECK to be defined. 
            case AMF_APP_UL_DATA_IND:   //[authentication response, Identy response, security mode complete, registration complete ]
                int amf_defs.amf_app_handle_uplink_nas_message(amf_app_desc_p, 
                             &AMF_APP_UL_DATA_IND(received_message_p).nas_msg);
                break;
            case AMF_APP_INITIAL_CONTEXT_SETUP_RSP:  
                amf_defs.amf_app_handle_initianl_context_setup_response_message(amf_app_desc_p,
                          &AMF_APP_INITIAL_CONTEXT_SETUP_RSP(received_message_p));
                break;
#endif
    case N11_PDU_SESSION_MODIFICATION_COMMAND:
      OAILOG_INFO(
          LOG_AMF_APP, "AMF_TEST: session created for imsi:%s \n",
          N11_PDU_SESSION_MODIFICATION_COMMAND(received_message_p).imsi);
      amf_defs.amf_app_handle_pdu_session_modification_command(
          &N11_PDU_SESSION_MODIFICATION_COMMAND(received_message_p));
      break;
    case N11_PDU_SESSION_MODIFICATION_REJECT:
      OAILOG_INFO(
          LOG_AMF_APP, "AMF_TEST: session created for imsi:%s \n",
          N11_PDU_SESSION_MODIFICATION_REJECT(received_message_p).imsi);
      amf_defs.amf_app_handle_pdu_session_modification_reject(
          &N11_PDU_SESSION_MODIFICATION_REJECT(received_message_p));
      break;
    case NGAP_PDUSESSIONRESOURCE_SETUP_RSP:
      /* This is non-nas message and can be handled directly to check if failure
       * or success messages are coming from NGAP
       */
      OAILOG_INFO(
          LOG_AMF_APP,
          "AMF_TEST: NGAP_PDUSESSIONRESOURCE_SETUP_RSP received\n");
      amf_app_handle_resource_setup_response(
          NGAP_PDUSESSIONRESOURCE_SETUP_RSP(received_message_p));
      break;

    case NGAP_PDUSESSIONRESOURCE_REL_RSP:
      /* This is non-nas message and can be handled directly to check if failure
       * or success messages are coming from NGAP
       */
      OAILOG_INFO(
          LOG_AMF_APP, "AMF_TEST: NGAP_PDUSESSIONRESOURCE_REL_RSP received\n");
      amf_app_handle_resource_release_response(
          NGAP_PDUSESSIONRESOURCE_REL_RSP(received_message_p));
      break;

    case NGAP_UE_CONTEXT_RELEASE_REQ:
      /* This is non-nas message and handled directly from NGAP sent to AMF
       * on RRC-Inactive mode to change UE's CM-connected to CM-idle state.
       */
      OAILOG_INFO(
          LOG_AMF_APP, "AMF_TEST: NGAP UE context release message to AMF"
	               " when gNB experiences RRC-Inactive \n");
      amf_app_handle_cm_idle_on_ue_context_release(
          NGAP_UE_CONTEXT_RELEASE_REQ(received_message_p));
      break;

    case TERMINATE_MESSAGE:
      OAILOG_INFO(LOG_AMF_APP, "AMF_TEST : TERMINATE_MESSAGE received\n");
      itti_free_msg_content(received_message_p);
      zframe_destroy(&msg_frame);
      amf_app_exit();
      break;
    default:
      OAILOG_INFO(
          LOG_AMF_APP, "AMF_TEST : default message received returning\n");
      break;
  }
}

void* amf_app_thread(void* args) {
  OAILOG_ERROR(LOG_AMF_APP, "Only for testing - amf_app_thread entered\n");

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
extern "C" int amf_app_init(const amf_config_t* amf_config_p) {
  if (amf_nas_state_init(amf_config_p)) {
    OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
  }
  OAILOG_ERROR(LOG_AMF_APP, "Only for testing - amf_nas_state_init done\n");
  // amf_app_edns_init(amf_config_p);
  // nas_network_initialize(amf_config_p); // needs to create initialization
  // part
  if (itti_create_task(TASK_AMF_APP, &amf_app_thread, NULL) < 0) {
    OAILOG_ERROR(LOG_AMF_APP, "AMF APP create task failed\n");
    OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
  }

  OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNok);
}

// void amf_app_main::amf_app_exit(void)
void amf_app_exit(void) {
  destroy_task_context(&amf_app_task_zmq_ctx);
  // put_amf_nas_state();
  // amf_app_edns_exit();
  // clear_amf_nas_state();
  // Clean-up NAS module
  // nas_network_cleanup();
  // amf_config_exit();

  OAI_FPRINTF_INFO("TASK_AMF_APP terminated\n");
  pthread_exit(NULL);
}
}  // namespace magma5g
