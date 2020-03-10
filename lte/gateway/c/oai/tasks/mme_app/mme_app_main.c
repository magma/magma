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

/*! \file mme_app_main.c
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <stdio.h>
#include <string.h>
#include <stdbool.h>
#include <stdint.h>
#include <pthread.h>

#include "bstrlib.h"
#include "dynamic_memory_check.h"
#include "log.h"
#include "intertask_interface.h"
#include "itti_free_defined_msg.h"
#include "mme_config.h"
#include "nas_network.h"
#include "timer.h"
#include "mme_app_extern.h"
#include "mme_app_ue_context.h"
#include "mme_app_defs.h"
#include "mme_app_statistics.h"
#include "service303_message_utils.h"
#include "service303.h"
#include "common_defs.h"
#include "mme_app_edns_emulation.h"
#include "nas_proc.h"
#include "3gpp_36.401.h"
#include "common_types.h"
#include "hashtable.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "mme_app_messages_types.h"
#include "mme_app_state.h"
#include "obj_hashtable.h"
#include "s11_messages_types.h"
#include "s1ap_messages_types.h"
#include "sctp_messages_types.h"
#include "timer_messages_types.h"

bool mme_hss_associated = false;
bool mme_sctp_bounded = false;

void *mme_app_thread(void *args);
static void _check_mme_healthy_and_notify_service(void);
static bool _is_mme_app_healthy(void);

//------------------------------------------------------------------------------
void *mme_app_thread(void *args)
{
  struct ue_mm_context_s *ue_context_p = NULL;
  itti_mark_task_ready(TASK_MME_APP);
  mme_app_desc_t *mme_app_desc_p;

  while (1) {
    MessageDef *received_message_p = NULL;

    /*
     * Trying to fetch a message from the message queue.
     * If the queue is empty, this function will block till a
     * message is sent to the task.
     */
    itti_receive_msg(TASK_MME_APP, &received_message_p);
    if (received_message_p == NULL) {
      OAILOG_ERROR(
        TASK_MME_APP, "Received an invalid Message from ITTI message queue\n");
      continue;
    }

    imsi64_t imsi64 = itti_get_associated_imsi(received_message_p);
    OAILOG_DEBUG(
      LOG_MME_APP, "Received message with imsi: " IMSI_64_FMT, imsi64);

    OAILOG_DEBUG(LOG_MME_APP, "Getting mme_nas_state");
    mme_app_desc_p = get_mme_nas_state(false);

    switch (ITTI_MSG_ID(received_message_p)) {
      case MESSAGE_TEST: {
        OAI_FPRINTF_INFO("TASK_MME_APP received MESSAGE_TEST\n");
      } break;

      case MME_APP_INITIAL_CONTEXT_SETUP_RSP: {
        mme_app_handle_initial_context_setup_rsp(mme_app_desc_p,
          &MME_APP_INITIAL_CONTEXT_SETUP_RSP(received_message_p));
      } break;

      case S6A_CANCEL_LOCATION_REQ: {
        /*
         * Check cancellation-type and handle it if it is SUBSCRIPTION_WITHDRAWAL.
         * For any other cancellation-type log it and ignore it.
         */
        mme_app_handle_s6a_cancel_location_req(mme_app_desc_p,
          &received_message_p->ittiMsg.s6a_cancel_location_req);
      } break;

      case MME_APP_UPLINK_DATA_IND: {
        nas_proc_ul_transfer_ind(
          MME_APP_UL_DATA_IND(received_message_p).ue_id,
          MME_APP_UL_DATA_IND(received_message_p).tai,
          MME_APP_UL_DATA_IND(received_message_p).cgi,
          &MME_APP_UL_DATA_IND(received_message_p).nas_msg);
      } break;

      case S11_CREATE_BEARER_REQUEST: {
        mme_app_handle_s11_create_bearer_req(mme_app_desc_p,
          &received_message_p->ittiMsg.s11_create_bearer_request);
      } break;

      case S6A_RESET_REQ: {
        mme_app_handle_s6a_reset_req(mme_app_desc_p,
          &received_message_p->ittiMsg.s6a_reset_req);
      } break;

      case S11_CREATE_SESSION_RESPONSE: {
        mme_app_handle_create_sess_resp(mme_app_desc_p,
          &received_message_p->ittiMsg.s11_create_session_response);
      } break;

      case S11_MODIFY_BEARER_RESPONSE: {
        OAILOG_INFO(
          TASK_MME_APP, "Received S11 MODIFY BEARER RESPONSE from SPGW\n");
        ue_context_p = mme_ue_context_exists_s11_teid(
          &mme_app_desc_p->mme_ue_contexts,
          received_message_p->ittiMsg.s11_modify_bearer_response.teid);

        if (ue_context_p == NULL) {
          OAILOG_WARNING(
            LOG_MME_APP,
            "We didn't find this teid in list of UE: %08x\n",
            received_message_p->ittiMsg.s11_modify_bearer_response.teid);
        } else {
          OAILOG_DEBUG(
            TASK_MME_APP, "S11 MODIFY BEARER RESPONSE local S11 teid = " TEID_FMT"\n",
            received_message_p->ittiMsg.s11_modify_bearer_response.teid);

          if (!ue_context_p->path_switch_req) {
            /* Updating statistics */
            update_mme_app_stats_s1u_bearer_add();
          } else {
            mme_app_handle_path_switch_req_ack(
              &received_message_p->ittiMsg.s11_modify_bearer_response,
              ue_context_p);
            ue_context_p->path_switch_req = false;
          }
        }
      } break;

      case S11_RELEASE_ACCESS_BEARERS_RESPONSE: {
        mme_app_handle_release_access_bearers_resp(mme_app_desc_p,
          &received_message_p->ittiMsg.s11_release_access_bearers_response);
      } break;

      case S11_DELETE_SESSION_RESPONSE: {
        mme_app_handle_delete_session_rsp(mme_app_desc_p,
          &received_message_p->ittiMsg.s11_delete_session_response);
      } break;

      case S11_SUSPEND_ACKNOWLEDGE: {
        mme_app_handle_suspend_acknowledge(mme_app_desc_p,
          &received_message_p->ittiMsg.s11_suspend_acknowledge);
      } break;

      case S1AP_E_RAB_SETUP_RSP: {
        mme_app_handle_e_rab_setup_rsp(mme_app_desc_p,
          &S1AP_E_RAB_SETUP_RSP(received_message_p));
      } break;

      case S1AP_E_RAB_REL_RSP: {
        mme_app_handle_e_rab_rel_rsp(&S1AP_E_RAB_REL_RSP(received_message_p));
      } break;

      case S1AP_INITIAL_UE_MESSAGE: {
        mme_app_handle_initial_ue_message(mme_app_desc_p,
          &S1AP_INITIAL_UE_MESSAGE(received_message_p));
      } break;

      case S6A_UPDATE_LOCATION_ANS: {
        /*
         * We received the update location answer message from HSS -> Handle it
         */
        OAILOG_INFO(LOG_MME_APP, "Received S6A Update Location Answer from S6A\n");
        mme_app_handle_s6a_update_location_ans(mme_app_desc_p,
          &received_message_p->ittiMsg.s6a_update_location_ans);
      } break;

      case S1AP_ENB_INITIATED_RESET_REQ: {
        mme_app_handle_enb_reset_req(
          &S1AP_ENB_INITIATED_RESET_REQ(received_message_p));
      } break;

      case S11_PAGING_REQUEST: {
        const char *imsi = received_message_p->ittiMsg.s11_paging_request.imsi;
        OAILOG_DEBUG(
          TASK_MME_APP, "MME handling paging request for IMSI%s\n", imsi);
        if (mme_app_handle_initial_paging_request(mme_app_desc_p, imsi)!=
            RETURNok) {
          OAILOG_ERROR(
            TASK_MME_APP,
            "Failed to send paging request to S1AP for IMSI%s\n",
            imsi);
        }
      } break;

      case MME_APP_INITIAL_CONTEXT_SETUP_FAILURE: {
        mme_app_handle_initial_context_setup_failure(mme_app_desc_p,
          &MME_APP_INITIAL_CONTEXT_SETUP_FAILURE(received_message_p));
      } break;

      case TIMER_HAS_EXPIRED: {
        /*
         * Check statistic timer
         */
        if (!timer_exists(
              received_message_p->ittiMsg.timer_has_expired.timer_id)) {
          OAILOG_WARNING(
            LOG_MME_APP,
            "Timer expiry signal received for timer \
            %lu, but it has already been deleted\n",
            received_message_p->ittiMsg.timer_has_expired.timer_id);
          break;
        }
        if (
          received_message_p->ittiMsg.timer_has_expired.timer_id ==
          mme_app_desc_p->statistic_timer_id) {
          mme_app_statistics_display();
        } else if (received_message_p->ittiMsg.timer_has_expired.arg != NULL) {
          mme_app_nas_timer_handle_signal_expiry(
            TIMER_HAS_EXPIRED(received_message_p).timer_id,
            TIMER_HAS_EXPIRED(received_message_p).arg);
        }
        timer_handle_expired(
          received_message_p->ittiMsg.timer_has_expired.timer_id);
      } break;

      case S1AP_UE_CAPABILITIES_IND: {
        mme_app_handle_s1ap_ue_capabilities_ind(mme_app_desc_p,
          &received_message_p->ittiMsg.s1ap_ue_cap_ind);
      } break;

      case S1AP_UE_CONTEXT_RELEASE_REQ: {
        mme_app_handle_s1ap_ue_context_release_req(
          &received_message_p->ittiMsg.s1ap_ue_context_release_req);
      } break;

      case S1AP_UE_CONTEXT_MODIFICATION_RESPONSE: {
        mme_app_handle_s1ap_ue_context_modification_resp(
          &mme_app_desc_p->mme_ue_contexts,
          &received_message_p->ittiMsg.s1ap_ue_context_mod_response);
      } break;

      case S1AP_UE_CONTEXT_MODIFICATION_FAILURE: {
        mme_app_handle_s1ap_ue_context_modification_fail(
          &mme_app_desc_p->mme_ue_contexts,
          &received_message_p->ittiMsg.s1ap_ue_context_mod_failure);
      } break;
      case S1AP_UE_CONTEXT_RELEASE_COMPLETE: {
        mme_app_handle_s1ap_ue_context_release_complete(mme_app_desc_p,
          &received_message_p->ittiMsg.s1ap_ue_context_release_complete);
      } break;

      case S1AP_ENB_DEREGISTERED_IND: {
        mme_app_handle_enb_deregister_ind(
          &received_message_p->ittiMsg.s1ap_eNB_deregistered_ind);
      } break;

      case ACTIVATE_MESSAGE: {
        mme_hss_associated = true;
        _check_mme_healthy_and_notify_service();
      } break;

      case SCTP_MME_SERVER_INITIALIZED: {
        mme_sctp_bounded =
          &received_message_p->ittiMsg.sctp_mme_server_initialized.successful;
        _check_mme_healthy_and_notify_service();
      } break;

      case S6A_PURGE_UE_ANS: {
        mme_app_handle_s6a_purge_ue_ans(
          &received_message_p->ittiMsg.s6a_purge_ue_ans);
      } break;

      case SGSAP_LOCATION_UPDATE_ACC: {
        /*Received SGSAP Location Update Accept message from SGS task*/
        OAILOG_INFO(
          TASK_MME_APP, "Received SGSAP Location Update Accept from SGS\n");
        mme_app_handle_sgsap_location_update_acc(mme_app_desc_p,
          &received_message_p->ittiMsg.sgsap_location_update_acc);
      } break;

      case SGSAP_LOCATION_UPDATE_REJ: {
        /*Received SGSAP Location Update Reject message from SGS task*/
        mme_app_handle_sgsap_location_update_rej(mme_app_desc_p,
          &received_message_p->ittiMsg.sgsap_location_update_rej);
      } break;

      case SGSAP_ALERT_REQUEST: {
        /*Received SGSAP Alert Request message from SGS task*/
        mme_app_handle_sgsap_alert_request(mme_app_desc_p,
          &received_message_p->ittiMsg.sgsap_alert_request);
      } break;

      case SGSAP_VLR_RESET_INDICATION: {
        /*Received SGSAP Reset Indication from SGS task*/
        mme_app_handle_sgsap_reset_indication(mme_app_desc_p,
          &received_message_p->ittiMsg.sgsap_vlr_reset_indication);
      } break;

      case SGSAP_PAGING_REQUEST: {
        mme_app_handle_sgsap_paging_request(mme_app_desc_p,
          &received_message_p->ittiMsg.sgsap_paging_request);
      } break;

      case SGSAP_SERVICE_ABORT_REQ: {
        mme_app_handle_sgsap_service_abort_request(mme_app_desc_p,
          &received_message_p->ittiMsg.sgsap_service_abort_req);
      } break;

      case SGSAP_EPS_DETACH_ACK: {
        mme_app_handle_sgs_eps_detach_ack(mme_app_desc_p,
          &received_message_p->ittiMsg.sgsap_eps_detach_ack);
      } break;

      case SGSAP_IMSI_DETACH_ACK: {
        mme_app_handle_sgs_imsi_detach_ack(mme_app_desc_p,
          &received_message_p->ittiMsg.sgsap_imsi_detach_ack);
      } break;

      case S11_MODIFY_UE_AMBR_REQUEST: {
        mme_app_handle_modify_ue_ambr_request(mme_app_desc_p,
          &S11_MODIFY_UE_AMBR_REQUEST(received_message_p));
      } break;

      case S11_NW_INITIATED_ACTIVATE_BEARER_REQUEST: {
        mme_app_handle_nw_init_ded_bearer_actv_req(mme_app_desc_p,
          &received_message_p->ittiMsg.s11_nw_init_actv_bearer_request);
      } break;

      case SGSAP_STATUS: {
        mme_app_handle_sgs_status_message(mme_app_desc_p,
          &received_message_p->ittiMsg.sgsap_status);
      } break;

      case S11_NW_INITIATED_DEACTIVATE_BEARER_REQUEST: {
        mme_app_handle_nw_init_bearer_deactv_req(mme_app_desc_p,
          &received_message_p->ittiMsg.s11_nw_init_deactv_bearer_request);
      } break;

      case S1AP_PATH_SWITCH_REQUEST: {
        mme_app_handle_path_switch_request(mme_app_desc_p,
          &S1AP_PATH_SWITCH_REQUEST(received_message_p));
      } break;

      case S6A_AUTH_INFO_ANS: {
        /*
         * We received the authentication vectors from HSS,
         * Normaly should trigger an authentication procedure towards UE.
         */
        nas_proc_authentication_info_answer(
          mme_app_desc_p, &S6A_AUTH_INFO_ANS(received_message_p));
      } break;

      case MME_APP_DOWNLINK_DATA_CNF: {
        nas_proc_dl_transfer_cnf(
          MME_APP_DL_DATA_CNF(received_message_p).ue_id,
          MME_APP_DL_DATA_CNF(received_message_p).err_code,
          &MME_APP_DL_DATA_REJ(received_message_p).nas_msg);
      } break;

      case MME_APP_DOWNLINK_DATA_REJ: {
        nas_proc_dl_transfer_rej(
          MME_APP_DL_DATA_REJ(received_message_p).ue_id,
          MME_APP_DL_DATA_REJ(received_message_p).err_code,
          &MME_APP_DL_DATA_REJ(received_message_p).nas_msg);
      } break;

      case SGSAP_DOWNLINK_UNITDATA: {
        /* We received the Downlink Unitdata from MSC, trigger a
         * Downlink Nas Transport message to UE.
         */
        nas_proc_downlink_unitdata(
          &SGSAP_DOWNLINK_UNITDATA(received_message_p));
      } break;

      case SGSAP_RELEASE_REQ: {
        /* We received the SGS Release request from MSC,to indicate that there
         * are no more NAS messages to be exchanged between the VLR and the UE,
         * or when a further exchange of NAS messages for the specified UE is
         * not possible due to an error.
         */
        nas_proc_sgs_release_req(&SGSAP_RELEASE_REQ(received_message_p));
      } break;

      case SGSAP_MM_INFORMATION_REQ: {
        // Received SGSAP MM Information Request message from SGS task
        nas_proc_cs_domain_mm_information_request(
          &SGSAP_MM_INFORMATION_REQ(received_message_p));
      } break;

     case TERMINATE_MESSAGE: {
       // Termination message received TODO -> release any data allocated
        put_mme_nas_state();
        mme_app_exit();
        itti_free_msg_content(received_message_p);
        itti_free(ITTI_MSG_ORIGIN_ID(received_message_p), received_message_p);
        OAI_FPRINTF_INFO("TASK_MME_APP terminated\n");
        itti_exit_task();
      } break;

      default: {
        OAILOG_ERROR(
          LOG_MME_APP,
          "Unknown message (%s) received with message Id: %d\n",
          ITTI_MSG_NAME(received_message_p),
          ITTI_MSG_ID(received_message_p));
      } break;
    }

    put_mme_nas_state();
    itti_free_msg_content(received_message_p);
    itti_free(ITTI_MSG_ORIGIN_ID(received_message_p), received_message_p);
    received_message_p = NULL;
  }

  return NULL;
}

//------------------------------------------------------------------------------
int mme_app_init(const mme_config_t *mme_config_p)
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  if (mme_nas_state_init(mme_config_p)) {
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  if (mme_app_edns_init(mme_config_p)) {
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  // Initialise NAS module
  nas_network_initialize(mme_config_p);
  /*
   * Create the thread associated with MME applicative layer
   */
  if (itti_create_task(TASK_MME_APP, &mme_app_thread, NULL) < 0) {
    OAILOG_ERROR(LOG_MME_APP, "MME APP create task failed\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  OAILOG_DEBUG(LOG_MME_APP, "Initializing MME applicative layer: DONE\n");
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

static void _check_mme_healthy_and_notify_service(void)
{
  if (_is_mme_app_healthy()) {
    send_app_health_to_service303(TASK_MME_APP, true);
  }
}

static bool _is_mme_app_healthy(void)
{
  return mme_hss_associated && mme_sctp_bounded;
}

//------------------------------------------------------------------------------
void mme_app_exit(void)
{
  mme_app_edns_exit();
  clear_mme_nas_state();
  // Clean-up NAS module
  nas_network_cleanup();
  mme_config_exit();
}
