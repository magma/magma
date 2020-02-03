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
#include "sgw_config.h"
#include "sgw_context_manager.h"
#include "pgw_ue_ip_address_alloc.h"
#include "pgw_pcef_emulation.h"
#include "spgw_config.h"

spgw_config_t spgw_config;

extern __pid_t g_pid;

static void sgw_exit(void);

//------------------------------------------------------------------------------
static void* sgw_intertask_interface(void* args_p)
{
  itti_mark_task_ready(TASK_SPGW_APP);
  spgw_state_t* spgw_state_p;

  while (1) {
    MessageDef* received_message_p = NULL;
    itti_receive_msg(TASK_SPGW_APP, &received_message_p);

    imsi64_t imsi64 = itti_get_associated_imsi(received_message_p);
    OAILOG_DEBUG(
      LOG_SPGW_APP,
      "Received message with imsi: " IMSI_64_FMT,
      imsi64);

    spgw_state_p = get_spgw_state(false);

    switch (ITTI_MSG_ID(received_message_p)) {
      case GTPV1U_CREATE_TUNNEL_RESP: {
        OAILOG_DEBUG(
          LOG_SPGW_APP,
          "Received teid for S1-U: %u and status: %s\n",
          received_message_p->ittiMsg.gtpv1uCreateTunnelResp.S1u_teid,
          received_message_p->ittiMsg.gtpv1uCreateTunnelResp.status == 0 ?
            "Success" :
            "Failure");
        sgw_handle_gtpv1uCreateTunnelResp(
          spgw_state_p, &received_message_p->ittiMsg.gtpv1uCreateTunnelResp, imsi64);
      } break;

      case MESSAGE_TEST:
        OAILOG_DEBUG(LOG_SPGW_APP, "Received MESSAGE_TEST\n");
        break;

      case S11_CREATE_BEARER_RESPONSE: {
        sgw_handle_create_bearer_response(
          spgw_state_p,
          &received_message_p->ittiMsg.s11_create_bearer_response);
      } break;

      case S11_CREATE_SESSION_REQUEST: {
        /*
         * We received a create session request from MME (with GTP abstraction here)
         * * * * procedures might be:
         * * * *      E-UTRAN Initial Attach
         * * * *      UE requests PDN connectivity
         */
        sgw_handle_create_session_request(
          spgw_state_p,
          &received_message_p->ittiMsg.s11_create_session_request,
          imsi64);
      } break;

      case S11_DELETE_SESSION_REQUEST: {
        sgw_handle_delete_session_request(
          spgw_state_p,
          &received_message_p->ittiMsg.s11_delete_session_request,
          imsi64);
      } break;

      case S11_MODIFY_BEARER_REQUEST: {
        sgw_handle_modify_bearer_request(
          spgw_state_p, &received_message_p->ittiMsg.s11_modify_bearer_request,
          imsi64);
      } break;

      case S11_RELEASE_ACCESS_BEARERS_REQUEST: {
        sgw_handle_release_access_bearers_request(
          spgw_state_p,
          &received_message_p->ittiMsg.s11_release_access_bearers_request,
          imsi64);
      } break;

      case S11_SUSPEND_NOTIFICATION: {
        sgw_handle_suspend_notification(
          spgw_state_p, &received_message_p->ittiMsg.s11_suspend_notification, imsi64);
      } break;

      case GTPV1U_UPDATE_TUNNEL_RESP: {
        sgw_handle_gtpv1uUpdateTunnelResp(
          spgw_state_p, &received_message_p->ittiMsg.gtpv1uUpdateTunnelResp, imsi64);
      } break;

      case SGI_CREATE_ENDPOINT_RESPONSE: {
        sgw_handle_sgi_endpoint_created(
          spgw_state_p,
          &received_message_p->ittiMsg.sgi_create_end_point_response, imsi64);
      } break;

      case SGI_UPDATE_ENDPOINT_RESPONSE: {
        sgw_handle_sgi_endpoint_updated(
          spgw_state_p,
          &received_message_p->ittiMsg.sgi_update_end_point_response, imsi64);
      } break;

      case S5_CREATE_BEARER_RESPONSE: {
        sgw_handle_s5_create_bearer_response(
          spgw_state_p, &received_message_p->ittiMsg.s5_create_bearer_response,
          imsi64);
      } break;

      case S5_NW_INITIATED_ACTIVATE_BEARER_REQ: {
        //Handle Dedicated bearer activation from PCRF
        if (
          sgw_handle_nw_initiated_actv_bearer_req(
            spgw_state_p,
            &received_message_p->ittiMsg.s5_nw_init_actv_bearer_request,
            imsi64) != RETURNok) {
          // If request handling fails send reject to PGW
          send_activate_dedicated_bearer_rsp_to_pgw(
            spgw_state_p,
            REQUEST_REJECTED /*Cause*/,
            received_message_p->ittiMsg.s5_nw_init_actv_bearer_request
              .s_gw_teid_S11_S4, /*SGW C-plane teid to fetch spgw context*/
            0 /*EBI*/,
            0 /*enb teid*/,
            0 /*sgw teid*/,
            imsi64);
        }
      } break;

      case S11_NW_INITIATED_ACTIVATE_BEARER_RESP: {
        //Handle Dedicated bearer Activation Rsp from MME
        sgw_handle_nw_initiated_actv_bearer_rsp(
          spgw_state_p,
          &received_message_p->ittiMsg.s11_nw_init_actv_bearer_rsp,
          imsi64);
      } break;

      case S5_NW_INITIATED_DEACTIVATE_BEARER_REQ: {
        //Handle Dedicated bearer Deactivation Req from PGW
        sgw_handle_nw_initiated_deactv_bearer_req(
          &received_message_p->ittiMsg.s5_nw_init_deactv_bearer_request,
          imsi64);
      } break;

      case S11_NW_INITIATED_DEACTIVATE_BEARER_RESP: {
        //Handle Dedicated bearer deactivation Rsp from MME
        sgw_handle_nw_initiated_deactv_bearer_rsp(
          spgw_state_p,
          &received_message_p->ittiMsg.s11_nw_init_deactv_bearer_rsp,
          imsi64);
      } break;

      case TERMINATE_MESSAGE: {
        put_spgw_state();
        sgw_exit();
        OAI_FPRINTF_INFO("TASK_SGW terminated\n");
        itti_exit_task();
      } break;

      default: {
        OAILOG_DEBUG(
          LOG_SPGW_APP,
          "Unkwnon message ID %d:%s\n",
          ITTI_MSG_ID(received_message_p),
          ITTI_MSG_NAME(received_message_p));
      } break;
    }

    put_spgw_state();

    itti_free_msg_content(received_message_p);
    itti_free(ITTI_MSG_ORIGIN_ID(received_message_p), received_message_p);
    received_message_p = NULL;
  }

  return NULL;
}

//------------------------------------------------------------------------------
int sgw_init(spgw_config_t* spgw_config_pP, bool persist_state)
{
  OAILOG_DEBUG(LOG_SPGW_APP, "Initializing SPGW-APP  task interface\n");

  if (spgw_state_init(persist_state, spgw_config_pP) < 0) {
    OAILOG_ALERT(LOG_SPGW_APP, "Error while initializing SGW state\n");
    return RETURNerror;
  }

  spgw_state_t* spgw_state_p = get_spgw_state(false);

  if (gtpv1u_init(spgw_state_p, spgw_config_pP, persist_state) < 0) {
    OAILOG_ALERT(LOG_SPGW_APP, "Initializing GTPv1-U ERROR\n");
    return RETURNerror;
  }

  if (
    RETURNerror ==
    pgw_pcef_emulation_init(spgw_state_p, &spgw_config_pP->pgw_config)) {
    return RETURNerror;
  }

  if (itti_create_task(TASK_SPGW_APP, &sgw_intertask_interface, NULL) < 0) {
    perror("pthread_create");
    OAILOG_ALERT(LOG_SPGW_APP, "Initializing SPGW-APP task interface: ERROR\n");
    return RETURNerror;
  }

  // Initial write of state, due to init of PCC rules on pcef emulation init.
  put_spgw_state();

  FILE* fp = NULL;
  bstring filename = bformat("/tmp/spgw_%d.status", g_pid);
  fp = fopen(bdata(filename), "w+");
  bdestroy_wrapper(&filename);
  fprintf(fp, "STARTED\n");
  fflush(fp);
  fclose(fp);

  OAILOG_DEBUG(LOG_SPGW_APP, "Initializing SPGW-APP task interface: DONE\n");
  return RETURNok;
}

//------------------------------------------------------------------------------
static void sgw_exit(void)
{
  OAILOG_DEBUG(LOG_SPGW_APP, "Cleaning SGW\n");

  spgw_state_exit();

  OAILOG_DEBUG(LOG_SPGW_APP, "Finished cleaning up SGW\n");
}
