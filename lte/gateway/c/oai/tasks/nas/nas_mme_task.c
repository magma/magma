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

#include <stdio.h>

#include "log.h"
#include "intertask_interface.h"
#include "itti_free_defined_msg.h"
#include "mme_config.h"
#include "nas_defs.h"
#include "nas_network.h"
#include "nas_proc.h"
#include "nas_timer.h"
#include "intertask_interface_types.h"
#include "mme_app_messages_types.h"
#include "mme_app_state.h"
#include "nas_messages_types.h"
#include "s1ap_messages_types.h"
#include "s6a_messages_types.h"
#include "sgs_messages_types.h"
#include "timer_messages_types.h"

static void nas_exit(void);

//------------------------------------------------------------------------------
static void *nas_intertask_interface(void *args_p)
{
  itti_mark_task_ready(TASK_NAS_MME);
  mme_app_desc_t *mme_app_desc_p;

  while (1) {
    MessageDef *received_message_p = NULL;

    itti_receive_msg(TASK_NAS_MME, &received_message_p);
    mme_app_desc_p = get_locked_mme_nas_state(false);

    switch (ITTI_MSG_ID(received_message_p)) {
      case MESSAGE_TEST: {
        OAI_FPRINTF_INFO("TASK_NAS_MME received MESSAGE_TEST\n");
      } break;

      case S1AP_DEREGISTER_UE_REQ: {
        nas_proc_deregister_ue(
          S1AP_DEREGISTER_UE_REQ(received_message_p).mme_ue_s1ap_id);
      } break;

      case SGSAP_DOWNLINK_UNITDATA: {
        /*
         * We received the Downlink Unitdata from MSC, trigger a
         * Downlink Nas Transport message to UE.
         */
        nas_proc_downlink_unitdata(
          &SGSAP_DOWNLINK_UNITDATA(received_message_p));
      } break;

      case SGSAP_RELEASE_REQ: {
        /*
         * We received the SGS Release request from MSC,to indicate that there are no more NAS messages to be exchanged
         * between the VLR and the UE, or when a further exchange of NAS messages for the specified UE is not possible
         * due to an error.
         */
        nas_proc_sgs_release_req(&SGSAP_RELEASE_REQ(received_message_p));
      } break;
      case SGSAP_MM_INFORMATION_REQ: {
        /*Received SGSAP MM Information Request message from SGS task*/
        nas_proc_cs_domain_mm_information_request(
          &SGSAP_MM_INFORMATION_REQ(received_message_p));
      } break;
      case NAS_CS_SERVICE_NOTIFICATION: {
        nas_proc_cs_service_notification(
          &NAS_CS_SERVICE_NOTIFICATION(received_message_p));
      } break;

      case NAS_NOTIFY_SERVICE_REJECT: {
        nas_proc_notify_service_reject(
          &NAS_NOTIFY_SERVICE_REJECT(received_message_p));
      } break;

      case TERMINATE_MESSAGE: {
        put_mme_nas_state(&mme_app_desc_p);
        nas_exit();
        OAI_FPRINTF_INFO("TASK_NAS_MME terminated\n");
        itti_free_msg_content(received_message_p);
        itti_free(ITTI_MSG_ORIGIN_ID(received_message_p), received_message_p);
        itti_exit_task();
      } break;

      case TIMER_HAS_EXPIRED: {
        /*
         * Call the NAS timer api
         */
        nas_timer_handle_signal_expiry(
          TIMER_HAS_EXPIRED(received_message_p).timer_id,
          TIMER_HAS_EXPIRED(received_message_p).arg);
      } break;

      default: {
        OAILOG_DEBUG(
          LOG_NAS,
          "Unkwnon message ID %d:%s from %s\n",
          ITTI_MSG_ID(received_message_p),
          ITTI_MSG_NAME(received_message_p),
          ITTI_MSG_ORIGIN_NAME(received_message_p));
      } break;
    }

    put_mme_nas_state(&mme_app_desc_p);
    itti_free_msg_content(received_message_p);
    itti_free(ITTI_MSG_ORIGIN_ID(received_message_p), received_message_p);
    received_message_p = NULL;
  }

  return NULL;
}

//------------------------------------------------------------------------------
int nas_init(mme_config_t *mme_config_p)
{
  OAILOG_DEBUG(LOG_NAS, "Initializing NAS task interface\n");
  nas_network_initialize(mme_config_p);

  if (itti_create_task(TASK_NAS_MME, &nas_intertask_interface, NULL) < 0) {
    OAILOG_ERROR(LOG_NAS, "Create task failed");
    OAILOG_DEBUG(LOG_NAS, "Initializing NAS task interface: FAILED\n");
    return -1;
  }

  OAILOG_DEBUG(LOG_NAS, "Initializing NAS task interface: DONE\n");
  return 0;
}

//------------------------------------------------------------------------------
static void nas_exit(void)
{
  OAILOG_DEBUG(LOG_NAS, "Cleaning NAS task interface\n");
  nas_network_cleanup();
  OAILOG_DEBUG(LOG_NAS, "Cleaning NAS task interface: DONE\n");
}
