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

/*! \file s6a_task.c
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <stdio.h>
#include <stdint.h>
#include <gnutls/gnutls.h>
#include <stdarg.h>
#include <string.h>

#include "bstrlib.h"
#include "3gpp_23.003.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "s6a_messages_types.h"
#include "service303.h"
#if HAVE_CONFIG_H
#include "config.h"
#endif
#include "log.h"
#include "assertions.h"
#include "intertask_interface.h"
#include "itti_free_defined_msg.h"
#include "common_defs.h"
#include "s6a_defs.h"
#include "s6a_messages.h"
#include "mme_config.h"
#include "s6a_client_api.h"
#include "s6a_c_iface.h"

static void s6a_exit(void);

struct session_handler* ts_sess_hdl;
task_zmq_ctx_t s6a_task_zmq_ctx;

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);
  int rc                         = RETURNerror;

  switch (ITTI_MSG_ID(received_message_p)) {
    case MESSAGE_TEST: {
      OAI_FPRINTF_INFO("TASK_S6A received MESSAGE_TEST\n");
    } break;
    case S6A_UPDATE_LOCATION_REQ: {
      rc = s6a_viface_update_location_req(
          &received_message_p->ittiMsg.s6a_update_location_req);
      if (rc) {
        OAILOG_DEBUG(
            LOG_S6A, "Sending s6a ULR for imsi=%s\n",
            received_message_p->ittiMsg.s6a_update_location_req.imsi);
      } else {
        OAILOG_ERROR(
            LOG_S6A, "Failure in sending s6a ULR for imsi=%s\n",
            received_message_p->ittiMsg.s6a_update_location_req.imsi);
      }
    } break;
    case S6A_AUTH_INFO_REQ: {
      rc = s6a_viface_authentication_info_req(
          &received_message_p->ittiMsg.s6a_auth_info_req);
      if (rc) {
        OAILOG_DEBUG(
            LOG_S6A, "Sending s6a AIR for imsi=%s\n",
            received_message_p->ittiMsg.s6a_auth_info_req.imsi);
      } else {
        OAILOG_ERROR(
            LOG_S6A, "Failure in sending s6a AIR for imsi=%s\n",
            received_message_p->ittiMsg.s6a_auth_info_req.imsi);
      }
    } break;
    case S6A_CANCEL_LOCATION_ANS: {
      s6a_viface_send_cancel_location_ans(
          &received_message_p->ittiMsg.s6a_cancel_location_ans);
    } break;
    case S6A_PURGE_UE_REQ: {
      uint8_t imsi_length;
      imsi_length = received_message_p->ittiMsg.s6a_purge_ue_req.imsi_length;
      if (imsi_length > IMSI_BCD_DIGITS_MAX) {
        OAILOG_ERROR(
            LOG_S6A, "imsi length exceeds IMSI_BCD_DIGITS_MAX length \n");
      }
      received_message_p->ittiMsg.s6a_purge_ue_req.imsi[imsi_length] = '\0';
      rc = s6a_viface_purge_ue(
          received_message_p->ittiMsg.s6a_purge_ue_req.imsi);
      if (rc) {
        OAILOG_DEBUG(
            LOG_S6A, "Sending s6a pur for imsi=%s\n",
            received_message_p->ittiMsg.s6a_purge_ue_req.imsi);
      } else {
        OAILOG_ERROR(
            LOG_S6A, "Failure in sending s6a pur for imsi=%s\n",
            received_message_p->ittiMsg.s6a_purge_ue_req.imsi);
      }
    } break;
    case TERMINATE_MESSAGE: {
      itti_free_msg_content(received_message_p);
      free(received_message_p);
      s6a_exit();
    } break;
    default: {
      OAILOG_DEBUG(
          LOG_S6A, "Unknown message ID %d: %s\n",
          ITTI_MSG_ID(received_message_p), ITTI_MSG_NAME(received_message_p));
    } break;
  }

  itti_free_msg_content(received_message_p);
  free(received_message_p);
  return 0;
}

//------------------------------------------------------------------------------
static void* s6a_thread(void* args) {
  itti_mark_task_ready(TASK_S6A);
  init_task_context(
      TASK_S6A, (task_id_t[]){TASK_MME_APP, TASK_S1AP}, 2, handle_message,
      &s6a_task_zmq_ctx);

  if (!s6a_viface_open((s6a_config_t*) args)) {
    OAILOG_ERROR(LOG_S6A, "Failed to initialize S6a interface");
    s6a_exit();
    return NULL;
  }

  zloop_start(s6a_task_zmq_ctx.event_loop);
  s6a_exit();
  return NULL;
}

//------------------------------------------------------------------------------
int s6a_init(const mme_config_t* mme_config_p) {
  OAILOG_DEBUG(LOG_S6A, "Initializing S6a interface\n");

  if (itti_create_task(
          TASK_S6A, &s6a_thread, (void*) &mme_config_p->s6a_config) < 0) {
    OAILOG_ERROR(LOG_S6A, "s6a create task\n");
    return RETURNerror;
  }

  return RETURNok;
}

//------------------------------------------------------------------------------
static void s6a_exit(void) {
  s6a_viface_close();
  destroy_task_context(&s6a_task_zmq_ctx);
  OAI_FPRINTF_INFO("TASK_S6A terminated\n");
  pthread_exit(NULL);
}
