/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

/*! \file sgw_task.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#define SGW
#define SGW_TASK_C

#include <stdio.h>
#include <netinet/in.h>
#include <sys/types.h>

#include "bstrlib.h"
#include "dynamic_memory_check.h"
#include "hashtable.h"
#include "log.h"
#include "common_defs.h"
#include "gtpv1_u_messages_types.h"
#include "gtpv1u_sgw_defs.h"
#include "intertask_interface.h"
#include "intertask_interface_types.h"
#include "itti_free_defined_msg.h"
#include "sgw_defs.h"
#include "sgw_handlers.h"
#include "pgw_handlers.h"
#include "sgw_config.h"
#include "sgw_context_manager.h"
#include "pgw_ue_ip_address_alloc.h"
#include "pgw_pcef_emulation.h"
#include "spgw_config.h"

static void spgw_app_exit(void);

spgw_config_t spgw_config;
task_zmq_ctx_t spgw_app_task_zmq_ctx;
extern __pid_t g_pid;

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  imsi64_t imsi64          = itti_get_associated_imsi(received_message_p);
  spgw_state_t* spgw_state = get_spgw_state(false);

  bool is_state_same = false;

  switch (ITTI_MSG_ID(received_message_p)) {
    case MESSAGE_TEST:
      is_state_same = true;  // task state is not changed
      OAILOG_DEBUG(LOG_SPGW_APP, "Received MESSAGE_TEST\n");
      break;

    case S11_CREATE_SESSION_REQUEST: {
      /*
       * We received a create session request from MME (with GTP abstraction
       * here)
       * * * * procedures might be:
       * * * *      E-UTRAN Initial Attach
       * * * *      UE requests PDN connectivity
       */
      sgw_handle_s11_create_session_request(
          spgw_state, &received_message_p->ittiMsg.s11_create_session_request,
          imsi64);
    } break;

    case S11_DELETE_SESSION_REQUEST: {
      sgw_handle_delete_session_request(
          &received_message_p->ittiMsg.s11_delete_session_request, imsi64);
      is_state_same = true;  // task state is not changed
    } break;

    case S11_MODIFY_BEARER_REQUEST: {
      sgw_handle_modify_bearer_request(
          &received_message_p->ittiMsg.s11_modify_bearer_request, imsi64);
      is_state_same = true;  // task state is not changed
    } break;

    case S11_RELEASE_ACCESS_BEARERS_REQUEST: {
      sgw_handle_release_access_bearers_request(
          &received_message_p->ittiMsg.s11_release_access_bearers_request,
          imsi64);
      is_state_same = true;  // task state is not changed
    } break;

    case S11_SUSPEND_NOTIFICATION: {
      sgw_handle_suspend_notification(
          &received_message_p->ittiMsg.s11_suspend_notification, imsi64);
      is_state_same = true;  // task state is not changed
    } break;

    case SGI_CREATE_ENDPOINT_RESPONSE: {
      sgw_handle_sgi_endpoint_created(
          spgw_state,
          &received_message_p->ittiMsg.sgi_create_end_point_response, imsi64);
    } break;

    case SGI_UPDATE_ENDPOINT_RESPONSE: {
      sgw_handle_sgi_endpoint_updated(
          &received_message_p->ittiMsg.sgi_update_end_point_response, imsi64);
      is_state_same = true;  // task state is not changed
    } break;

    case S11_NW_INITIATED_ACTIVATE_BEARER_RESP: {
      // Handle Dedicated bearer Activation Rsp from MME
      sgw_handle_nw_initiated_actv_bearer_rsp(
          &received_message_p->ittiMsg.s11_nw_init_actv_bearer_rsp, imsi64);
      is_state_same = true;  // task state is not changed
    } break;

    case S11_NW_INITIATED_DEACTIVATE_BEARER_RESP: {
      // Handle Dedicated bearer deactivation Rsp from MME
      sgw_handle_nw_initiated_deactv_bearer_rsp(
          spgw_state,
          &received_message_p->ittiMsg.s11_nw_init_deactv_bearer_rsp, imsi64);
      is_state_same = true;  // task state is not changed
    } break;

    case PCEF_CREATE_SESSION_RESPONSE: {
      spgw_handle_pcef_create_session_response(
          spgw_state, &received_message_p->ittiMsg.pcef_create_session_response,
          imsi64);
    } break;

    case GX_NW_INITIATED_ACTIVATE_BEARER_REQ: {
      /* TODO need to discuss as part sending response to PCEF,
       * should these errors need to be mapped to gx errors
       * or sessiond does mapping of these error codes to gx error codes
       */
      gtpv2c_cause_value_t failed_cause = REQUEST_ACCEPTED;
      int32_t rc = spgw_handle_nw_initiated_bearer_actv_req(
          spgw_state,
          &received_message_p->ittiMsg.gx_nw_init_actv_bearer_request, imsi64,
          &failed_cause);
      if (rc != RETURNok) {
        OAILOG_ERROR_UE(
            LOG_SPGW_APP, imsi64,
            "Send Create Bearer Failure Response to PCRF with cause :%d \n",
            failed_cause);
        // Send Reject to PCRF
        // TODO-Uncomment once implemented at PCRF
        /* rc = send_dedicated_bearer_actv_rsp(bearer_req_p->lbi,
         *    failed_cause);
         */
      }
    } break;

    case GX_NW_INITIATED_DEACTIVATE_BEARER_REQ: {
      int32_t rc = spgw_handle_nw_initiated_bearer_deactv_req(
          &received_message_p->ittiMsg.gx_nw_init_deactv_bearer_request,
          imsi64);
      is_state_same = true;  // task state is not changed
      if (rc != RETURNok) {
        OAILOG_ERROR_UE(
            LOG_SPGW_APP, imsi64,
            "Failed to handle NW_INITIATED_DEACTIVATE_BEARER_REQ, "
            "send bearer deactivation reject to SPGW service \n");
        // TODO-Uncomment once implemented at PCRF
        /* rc =
         * send_dedicated_bearer_deactv_rsp(invalid_bearer_id,REQUEST_REJECTED);
         */
      }
    } break;

    case IP_ALLOCATION_RESPONSE: {
      int32_t rc = sgw_handle_ip_allocation_rsp(
          spgw_state, &received_message_p->ittiMsg.ip_allocation_response,
          imsi64);
      if (rc != RETURNok) {
        OAILOG_ERROR_UE(
            LOG_SPGW_APP, imsi64,
            "Failed to handle IP_ALLOCATION_RESPONSE, \n");
      }
    } break;

    case TERMINATE_MESSAGE: {
      itti_free_msg_content(received_message_p);
      free(received_message_p);
      spgw_app_exit();
    } break;

    default: {
      OAILOG_DEBUG(
          LOG_SPGW_APP, "Unknown message ID %d:%s\n",
          ITTI_MSG_ID(received_message_p), ITTI_MSG_NAME(received_message_p));
    } break;
  }

  if (!is_state_same) {
    put_spgw_state();
  }
  put_spgw_ue_state(imsi64);

  itti_free_msg_content(received_message_p);
  free(received_message_p);
  return 0;
}

//------------------------------------------------------------------------------
static void* spgw_app_thread(__attribute__((unused)) void* args) {
  itti_mark_task_ready(TASK_SPGW_APP);
  init_task_context(
      TASK_SPGW_APP, (task_id_t[]){TASK_MME_APP}, 1, handle_message,
      &spgw_app_task_zmq_ctx);

  zloop_start(spgw_app_task_zmq_ctx.event_loop);
  spgw_app_exit();
  return NULL;
}

//------------------------------------------------------------------------------
int spgw_app_init(spgw_config_t* spgw_config_pP, bool persist_state) {
  OAILOG_DEBUG(LOG_SPGW_APP, "Initializing SPGW-APP  task interface\n");

  if (spgw_state_init(persist_state, spgw_config_pP) < 0) {
    OAILOG_ALERT(LOG_SPGW_APP, "Error while initializing SGW state\n");
    return RETURNerror;
  }

  spgw_state_t* spgw_state_p = get_spgw_state(false);

  // Read SPGW state for subscribers from db
  read_spgw_ue_state_db();

  if (gtpv1u_init(spgw_state_p, spgw_config_pP, persist_state) < 0) {
    OAILOG_ALERT(LOG_SPGW_APP, "Initializing GTPv1-U ERROR\n");
    return RETURNerror;
  }

  if (RETURNerror ==
      pgw_pcef_emulation_init(spgw_state_p, &spgw_config_pP->pgw_config)) {
    return RETURNerror;
  }

  if (itti_create_task(TASK_SPGW_APP, &spgw_app_thread, NULL) < 0) {
    perror("pthread_create");
    OAILOG_ALERT(LOG_SPGW_APP, "Initializing SPGW-APP task interface: ERROR\n");
    return RETURNerror;
  }

  OAILOG_DEBUG(LOG_SPGW_APP, "Initializing SPGW-APP task interface: DONE\n");
  return RETURNok;
}

//------------------------------------------------------------------------------
static void spgw_app_exit(void) {
  OAILOG_DEBUG(LOG_SPGW_APP, "Cleaning SGW\n");
  put_spgw_state();
  gtpv1u_exit();
  spgw_state_exit();
  destroy_task_context(&spgw_app_task_zmq_ctx);
  OAILOG_DEBUG(LOG_SPGW_APP, "Finished cleaning up SGW\n");
  OAI_FPRINTF_INFO("TASK_SPGW_APP terminated\n");
  pthread_exit(NULL);
}
