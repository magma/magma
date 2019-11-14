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

/*! \file sgs_task.c
  \brief
  \author
  \company
  \email:
*/
#define SGG
#define SGS_TASK_C

#include <stdio.h>

#include "log.h"
#include "intertask_interface.h"
#include "mme_config.h"
#include "sgs_messages_types.h"
#include "csfb_client_api.h"
#include "common_defs.h"
#include "intertask_interface_types.h"

static void sgs_exit(void);

//------------------------------------------------------------------------------
static void *sgs_thread(void *args_p)
{
  itti_mark_task_ready(TASK_SGS);

  while (1) {
    MessageDef *received_message_p = NULL;

    itti_receive_msg(TASK_SGS, &received_message_p);

    switch (ITTI_MSG_ID(received_message_p)) {
      case SGSAP_LOCATION_UPDATE_REQ: {
        /*
         * We received a SGs location update request from MME
         * * * * procedures might be:
         * * * *      E-UTRAN Combined Attach
         * * * *      TAU
         */
        OAILOG_DEBUG(LOG_SGS, "Received SGSAP_LOCATION_UPDATE_REQ message \n");
        /* send Location Update Request message to FeG*/
        send_location_update_request(
          &received_message_p->ittiMsg.sgsap_location_update_req);
      } break;

      case SGSAP_UPLINK_UNITDATA: {
        /*
         * We received a SGs uplink unitdata message from NAS
         * * * * procedures might be:
         * * * *      Mobile origination SMS - Uplink Nas Transport message
         * * * *      Mobile terminating SMS - Uplink Nas Transport message
         */
        OAILOG_DEBUG(LOG_SGS, "Received SGSAP_UPLINK_UNITDATA message \n");
        send_uplink_unitdata(
          &received_message_p->ittiMsg.sgsap_uplink_unitdata);
      } break;

      case SGSAP_EPS_DETACH_IND: {
        /*
         * We received a SGs eps detach indication from MME
         * * * * procedures might be:
         * * * *      Ue initiated Detach
         * * * *      Network Initiated Detach
         */
        OAILOG_DEBUG(LOG_SGS, "Received SGSAP_EPS_DETACH_IND message \n");
        /* send EPS Detach Indication message to FeG*/
        send_eps_detach_indication(
          &received_message_p->ittiMsg.sgsap_eps_detach_ind);
      } break;

      case SGSAP_IMSI_DETACH_IND: {
        /*
         * We received a SGs imsi detach indication from MME
         * * * * procedures might be:
         * * * *      Ue initiated Detach
         * * * *      Network Initiated Detach
         */
        OAILOG_DEBUG(LOG_SGS, "Received SGSAP_IMSI_DETACH_IND message \n");
        /* send IMSI Detach Indication message to FeG*/
        send_imsi_detach_indication(
          &received_message_p->ittiMsg.sgsap_imsi_detach_ind);
      } break;

      case SGSAP_TMSI_REALLOC_COMP: {
        /*
         * We received a SGs tmsi reallocation complete from NAS
         * * * * procedures might be:
         * * * *      Attach Complete
         * * * *      Tracking Area Update Complete
         */
        OAILOG_DEBUG(LOG_SGS, "Received SGSAP_TMSI_REALLOC_COMP message \n");
        /* send tmsi reallocation complete message to FeG*/
        send_tmsi_reallocation_complete(
          &received_message_p->ittiMsg.sgsap_tmsi_realloc_comp);
      } break;

      case SGSAP_UE_ACTIVITY_IND: {
        /*
         * We received a SGs ue activity indication from NAS
         * * * * procedures might be:
         * * * *      Service Request for SMS or PS data
         * * * *      Extended Service Request for MT CSFB in connected mode
         */
        OAILOG_DEBUG(
          LOG_SGS, "Received SGSAP_UE_ACTIVITY_IND message from NAS");
        /* send sgsap ue activity indication message to FeG*/
        send_ue_activity_indication(
          &received_message_p->ittiMsg.sgsap_ue_activity_ind);
      } break;

      case SGSAP_ALERT_ACK: {
        /*
         * We received a SGs Alert Ack from MME-app
         * * * * Message sent as part of procedure:
         * * * * Non-eps alert
         */
        OAILOG_DEBUG(LOG_SGS, "Received SGSAP_ALERT_ACK message");
        /* send SGs Alert Ack to FeG*/
        send_alert_ack(&received_message_p->ittiMsg.sgsap_alert_ack);
      } break;

      case SGSAP_ALERT_REJECT: {
        /*
         * We received a SGs Alert Reject from MME-app
         * * * * Message sent as part of procedure:
         * * * * Non-eps alert
         */
        OAILOG_DEBUG(LOG_SGS, "Received SGSAP_ALERT_REJECT message");
        /* send SGs Alert Reject to FeG*/
        send_alert_reject(&received_message_p->ittiMsg.sgsap_alert_reject);
      } break;

      case SGSAP_SERVICE_REQUEST: {
        OAILOG_DEBUG(LOG_SGS, "Received SGSAP_SERVICE_REQUEST message \n");
        send_service_request(&SGSAP_SERVICE_REQUEST(received_message_p));
      } break;

      case SGSAP_PAGING_REJECT: {
        OAILOG_DEBUG(
          LOG_SGS, "Received  message SGSAP_PAGING_REJECT message \n");
        send_paging_reject(&SGSAP_PAGING_REJECT(received_message_p));
      } break;

      case SGSAP_UE_UNREACHABLE: {
        OAILOG_DEBUG(
          LOG_SGS, "Received  message SGSAP_UE_UNREACHABLE message \n");
        send_ue_unreachable(&SGSAP_UE_UNREACHABLE(received_message_p));
      } break;

      case TERMINATE_MESSAGE: {
        sgs_exit();
        itti_exit_task();
      } break;

      default: {
        OAILOG_DEBUG(
          LOG_SGS,
          "Unkwnon message ID %d:%s\n",
          ITTI_MSG_ID(received_message_p),
          ITTI_MSG_NAME(received_message_p));
      } break;
    }

    itti_free(ITTI_MSG_ORIGIN_ID(received_message_p), received_message_p);
    received_message_p = NULL;
  }

  return NULL;
}

//------------------------------------------------------------------------------
int sgs_init(const mme_config_t *mme_config_p)
{
  OAILOG_DEBUG(LOG_SGS, "Initializing SGS task interface\n");

  if (itti_create_task(TASK_SGS, &sgs_thread, NULL) < 0) {
    OAILOG_ERROR(LOG_SGS, "sgs create task\n");
    return RETURNerror;
  }
  OAILOG_DEBUG(LOG_SGS, "Initializing SGS task interface: DONE\n");
  return RETURNok;
}

//------------------------------------------------------------------------------
static void sgs_exit(void)
{
  /*Cleanup SGS context*/
}
