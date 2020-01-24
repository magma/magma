/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
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
#include "timer_messages_types.h"
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
#include "timer.h"
#include "s6a_client_api.h"
#include "s6a_c_iface.h"


struct session_handler *ts_sess_hdl;

void *s6a_thread(void *args);
static void s6a_exit(void);

//------------------------------------------------------------------------------
void *s6a_thread(void *args)
{
  itti_mark_task_ready(TASK_S6A);

  while (1) {
    MessageDef *received_message_p = NULL;
    int rc = RETURNerror;

    /*
     * Trying to fetch a message from the message queue.
     * * If the queue is empty, this function will block till a
     * * message is sent to the task.
     */
    itti_receive_msg(TASK_S6A, &received_message_p);
    DevAssert(received_message_p);

    switch (ITTI_MSG_ID(received_message_p)) {
      case MESSAGE_TEST: {
        OAI_FPRINTF_INFO("TASK_S6A received MESSAGE_TEST\n");
      } break;
      case S6A_UPDATE_LOCATION_REQ: {
        rc = s6a_viface_update_location_req(
          &received_message_p->ittiMsg.s6a_update_location_req);
        if (rc) {
          OAILOG_DEBUG(
            LOG_S6A,
            "Sending s6a ULR for imsi=%s\n",
            received_message_p->ittiMsg.s6a_update_location_req.imsi);
        } else {
          OAILOG_ERROR(
            LOG_S6A,
            "Failure in sending s6a ULR for imsi=%s\n",
            received_message_p->ittiMsg.s6a_update_location_req.imsi);
        }
      } break;
      case S6A_AUTH_INFO_REQ: {
        rc = s6a_viface_authentication_info_req(
          &received_message_p->ittiMsg.s6a_auth_info_req);
        if (rc) {
          OAILOG_DEBUG(
            LOG_S6A,
            "Sending s6a AIR for imsi=%s\n",
            received_message_p->ittiMsg.s6a_auth_info_req.imsi);
        } else {
          OAILOG_ERROR(
            LOG_S6A,
            "Failure in sending s6a AIR for imsi=%s\n",
            received_message_p->ittiMsg.s6a_auth_info_req.imsi);
        }
      } break;
      case TIMER_HAS_EXPIRED: {
        s6a_viface_timer_expired(received_message_p->ittiMsg.timer_has_expired.timer_id);
      } break;
      case TERMINATE_MESSAGE: {
        s6a_exit();
        itti_free_msg_content(received_message_p);
        itti_free(ITTI_MSG_ORIGIN_ID(received_message_p), received_message_p);
        OAI_FPRINTF_INFO("TASK_S6A terminated\n");
        itti_exit_task();
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
        rc = s6a_viface_purge_ue(received_message_p->ittiMsg.s6a_purge_ue_req.imsi);

        if (rc) {
          OAILOG_DEBUG(
            LOG_S6A,
            "Sending s6a pur for imsi=%s\n",
            received_message_p->ittiMsg.s6a_purge_ue_req.imsi);
        } else {
          OAILOG_ERROR(
            LOG_S6A,
            "Failure in sending s6a pur for imsi=%s\n",
            received_message_p->ittiMsg.s6a_purge_ue_req.imsi);
        }
      } break;
      default: {
        OAILOG_DEBUG(
          LOG_S6A,
          "Unkwnon message ID %d: %s\n",
          ITTI_MSG_ID(received_message_p),
          ITTI_MSG_NAME(received_message_p));
      } break;
    }
    itti_free_msg_content(received_message_p);
    itti_free(ITTI_MSG_ORIGIN_ID(received_message_p), received_message_p);
    received_message_p = NULL;
  }
  return NULL;
}

//------------------------------------------------------------------------------
int s6a_init(const mme_config_t *mme_config_p)
{
  OAILOG_DEBUG(LOG_S6A, "Initializing S6a interface\n");

  if (itti_create_task(TASK_S6A, &s6a_thread, NULL) < 0) {
    OAILOG_ERROR(LOG_S6A, "s6a create task\n");
    return RETURNerror;
  }

  if (s6a_viface_open(&mme_config_p->s6a_config)) return RETURNok;
  else return RETURNerror;
}

//------------------------------------------------------------------------------
static void s6a_exit(void)
{
  s6a_viface_close();
}
