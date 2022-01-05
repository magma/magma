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

#include "lte/gateway/c/core/oai/common/assertions.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/common/itti_free_defined_msg.h"
#include "lte/gateway/c/core/oai/include/mme_config.h"
#include "lte/gateway/c/core/oai/tasks/nas/nas_network.h"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_extern.h"
#include "lte/gateway/c/core/oai/include/mme_app_ue_context.h"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_defs.h"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_ha.h"
#include "lte/gateway/c/core/oai/include/mme_app_statistics.h"
#include "lte/gateway/c/core/oai/lib/message_utils/service303_message_utils.h"
#include "lte/gateway/c/core/oai/common/common_defs.h"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_edns_emulation.h"
#include "lte/gateway/c/core/oai/tasks/nas/nas_proc.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/include/mme_app_messages_types.h"
#include "lte/gateway/c/core/oai/include/mme_app_state.h"
#include "lte/gateway/c/core/oai/include/s11_messages_types.h"
#include "lte/gateway/c/core/oai/include/s1ap_messages_types.h"

static void check_mme_healthy_and_notify_service(void);
static bool is_mme_app_healthy(void);
static void mme_app_exit(void);
static void start_stats_timer(void);

bool mme_hss_associated = false;
bool mme_sctp_bounded   = false;
task_zmq_ctx_t mme_app_task_zmq_ctx;
bool mme_congestion_control_enabled = true;
long mme_app_last_msg_latency;
long pre_mme_task_msg_latency;
static long epc_stats_timer_id;
static size_t epc_stats_timer_sec = 60;

mme_congestion_params_t mme_congestion_params;

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);
  imsi64_t imsi64                = itti_get_associated_imsi(received_message_p);
  mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);

  bool is_task_state_same = false;
  bool force_ue_write     = false;

  mme_app_last_msg_latency =
      ITTI_MSG_LATENCY(received_message_p);  // microseconds
  pre_mme_task_msg_latency = ITTI_MSG_LASTHOP_LATENCY(received_message_p);

  OAILOG_DEBUG(
      LOG_MME_APP, "MME APP ZMQ latency: %ld.", mme_app_last_msg_latency);

  switch (ITTI_MSG_ID(received_message_p)) {
    case MESSAGE_TEST: {
      OAI_FPRINTF_INFO("TASK_MME_APP received MESSAGE_TEST\n");
    } break;

    case MME_APP_INITIAL_CONTEXT_SETUP_RSP: {
      mme_app_handle_initial_context_setup_rsp(
          &MME_APP_INITIAL_CONTEXT_SETUP_RSP(received_message_p));
      is_task_state_same = true;
    } break;

    case S6A_CANCEL_LOCATION_REQ: {
      /*
       * Check cancellation-type and handle it if it is SUBSCRIPTION_WITHDRAWAL.
       * For any other cancellation-type log it and ignore it.
       */
      mme_app_handle_s6a_cancel_location_req(
          mme_app_desc_p, &received_message_p->ittiMsg.s6a_cancel_location_req);
      is_task_state_same = true;
    } break;

    case MME_APP_UPLINK_DATA_IND: {
      nas_proc_ul_transfer_ind(
          MME_APP_UL_DATA_IND(received_message_p).ue_id,
          MME_APP_UL_DATA_IND(received_message_p).tai,
          MME_APP_UL_DATA_IND(received_message_p).cgi,
          &MME_APP_UL_DATA_IND(received_message_p).nas_msg);
      force_ue_write     = true;
      is_task_state_same = true;
    } break;

    case S11_CREATE_BEARER_REQUEST: {
      mme_app_handle_s11_create_bearer_req(
          mme_app_desc_p,
          &received_message_p->ittiMsg.s11_create_bearer_request);
      is_task_state_same = true;
    } break;

    case S6A_RESET_REQ: {
      mme_app_handle_s6a_reset_req(&received_message_p->ittiMsg.s6a_reset_req);
      is_task_state_same = true;
    } break;

    case S11_CREATE_SESSION_RESPONSE: {
      mme_app_handle_create_sess_resp(
          mme_app_desc_p,
          &received_message_p->ittiMsg.s11_create_session_response);
    } break;

    case S11_MODIFY_BEARER_RESPONSE: {
      ue_mm_context_t* ue_context_p = NULL;
      OAILOG_INFO(
          LOG_MME_APP, "Received S11 MODIFY BEARER RESPONSE from SPGW\n");
      ue_context_p = mme_ue_context_exists_s11_teid(
          &mme_app_desc_p->mme_ue_contexts,
          received_message_p->ittiMsg.s11_modify_bearer_response.teid);

      if (ue_context_p == NULL) {
        OAILOG_WARNING(
            LOG_MME_APP, "We didn't find this teid in list of UE: %08x\n",
            received_message_p->ittiMsg.s11_modify_bearer_response.teid);
      } else {
        OAILOG_DEBUG(
            LOG_MME_APP,
            "S11 MODIFY BEARER RESPONSE local S11 teid = " TEID_FMT "\n",
            received_message_p->ittiMsg.s11_modify_bearer_response.teid);
        if ((!ue_context_p->path_switch_req) && (!ue_context_p->erab_mod_ind)) {
          /* Updating statistics */
          mme_app_handle_modify_bearer_rsp(
              &received_message_p->ittiMsg.s11_modify_bearer_response,
              ue_context_p);
          update_mme_app_stats_s1u_bearer_add();
        } else if (ue_context_p->path_switch_req) {
          mme_app_handle_path_switch_req_ack(
              &received_message_p->ittiMsg.s11_modify_bearer_response,
              ue_context_p);
          ue_context_p->path_switch_req = false;
        } else if (ue_context_p->erab_mod_ind) {
          mme_app_handle_modify_bearer_rsp_erab_mod_ind(
              &received_message_p->ittiMsg.s11_modify_bearer_response,
              ue_context_p);
          ue_context_p->erab_mod_ind = false;
        }

        // Check if an offloading request is pending for this UE as part
        // of HA implementation
        if (ue_context_p->ue_context_rel_cause ==
            S1AP_NAS_MME_PENDING_OFFLOADING) {
          OAILOG_INFO(
              LOG_MME_APP,
              "UE CONTEXT REL CAUSE is S1AP_NAS_MME_PENDING_OFFLOADING");
          // This will be again overwritten when a release request is received.
          // It is safe to set it any value other than
          // S1AP_NAS_MME_PENDING_OFFLOADING to allow a UE be able to reattach
          // to this AGW instance.
          ue_context_p->ue_context_rel_cause = S1AP_INVALID_CAUSE;
          mme_app_handle_ue_offload(ue_context_p);
        }
        force_ue_write = true;
      }
    } break;

    case S11_RELEASE_ACCESS_BEARERS_RESPONSE: {
      mme_app_handle_release_access_bearers_resp(
          mme_app_desc_p,
          &received_message_p->ittiMsg.s11_release_access_bearers_response);
    } break;

    case S11_DELETE_SESSION_RESPONSE: {
      mme_app_handle_delete_session_rsp(
          mme_app_desc_p,
          &received_message_p->ittiMsg.s11_delete_session_response);
    } break;

    case S11_SUSPEND_ACKNOWLEDGE: {
      mme_app_handle_suspend_acknowledge(
          mme_app_desc_p, &received_message_p->ittiMsg.s11_suspend_acknowledge);
    } break;

    case S1AP_E_RAB_SETUP_RSP: {
      mme_app_handle_e_rab_setup_rsp(&S1AP_E_RAB_SETUP_RSP(received_message_p));
      is_task_state_same = true;
    } break;

    case S1AP_E_RAB_REL_RSP: {
      mme_app_handle_e_rab_rel_rsp(&S1AP_E_RAB_REL_RSP(received_message_p));
      is_task_state_same = true;
    } break;

    case S1AP_E_RAB_MODIFICATION_IND: {
      mme_app_handle_e_rab_modification_ind(
          &S1AP_E_RAB_MODIFICATION_IND(received_message_p));
      is_task_state_same = true;
    } break;

    case S1AP_INITIAL_UE_MESSAGE: {
      imsi64 = mme_app_handle_initial_ue_message(
          mme_app_desc_p, &S1AP_INITIAL_UE_MESSAGE(received_message_p));
    } break;

    case S6A_UPDATE_LOCATION_ANS: {
      /*
       * We received the update location answer message from HSS -> Handle it
       */
      OAILOG_INFO(
          LOG_MME_APP, "Received S6A Update Location Answer from S6A\n");
      mme_app_handle_s6a_update_location_ans(
          mme_app_desc_p, &received_message_p->ittiMsg.s6a_update_location_ans);
      is_task_state_same = true;
    } break;

    case S1AP_ENB_INITIATED_RESET_REQ: {
      mme_app_handle_enb_reset_req(
          &S1AP_ENB_INITIATED_RESET_REQ(received_message_p));
      is_task_state_same = true;
    } break;

    case S11_PAGING_REQUEST: {
      OAILOG_DEBUG(LOG_MME_APP, "MME handling paging request \n");
      imsi64 = mme_app_handle_initial_paging_request(
          mme_app_desc_p, &received_message_p->ittiMsg.s11_paging_request);
    } break;

    case MME_APP_INITIAL_CONTEXT_SETUP_FAILURE: {
      mme_app_handle_initial_context_setup_failure(
          &MME_APP_INITIAL_CONTEXT_SETUP_FAILURE(received_message_p));
      is_task_state_same = true;
    } break;

    case S1AP_UE_CAPABILITIES_IND: {
      mme_app_handle_s1ap_ue_capabilities_ind(
          &received_message_p->ittiMsg.s1ap_ue_cap_ind);
      is_task_state_same = true;
    } break;

    case S1AP_UE_CONTEXT_RELEASE_REQ: {
      mme_app_handle_s1ap_ue_context_release_req(
          &received_message_p->ittiMsg.s1ap_ue_context_release_req);
    } break;

    case S1AP_UE_CONTEXT_MODIFICATION_RESPONSE: {
      mme_app_handle_s1ap_ue_context_modification_resp(
          &received_message_p->ittiMsg.s1ap_ue_context_mod_response);
      is_task_state_same = true;
    } break;

    case S1AP_UE_CONTEXT_MODIFICATION_FAILURE: {
      mme_app_handle_s1ap_ue_context_modification_fail(
          &received_message_p->ittiMsg.s1ap_ue_context_mod_failure);
      is_task_state_same = true;
    } break;
    case S1AP_UE_CONTEXT_RELEASE_COMPLETE: {
      mme_app_handle_s1ap_ue_context_release_complete(
          mme_app_desc_p,
          &received_message_p->ittiMsg.s1ap_ue_context_release_complete);
    } break;

    case S1AP_ENB_DEREGISTERED_IND: {
      mme_app_handle_enb_deregister_ind(
          &received_message_p->ittiMsg.s1ap_eNB_deregistered_ind);
    } break;

    case ACTIVATE_MESSAGE: {
      mme_hss_associated = true;
      is_task_state_same = true;
      check_mme_healthy_and_notify_service();
    } break;

    case SCTP_MME_SERVER_INITIALIZED: {
      mme_sctp_bounded =
          &received_message_p->ittiMsg.sctp_mme_server_initialized.successful;
      check_mme_healthy_and_notify_service();
      is_task_state_same = true;
    } break;

    case S6A_PURGE_UE_ANS: {
      mme_app_handle_s6a_purge_ue_ans(
          &received_message_p->ittiMsg.s6a_purge_ue_ans);
      is_task_state_same = true;
    } break;

    case SGSAP_LOCATION_UPDATE_ACC: {
      /*Received SGSAP Location Update Accept message from SGS task*/
      OAILOG_INFO(
          LOG_MME_APP, "Received SGSAP Location Update Accept from SGS\n");
      mme_app_handle_sgsap_location_update_acc(
          mme_app_desc_p,
          &received_message_p->ittiMsg.sgsap_location_update_acc);
      is_task_state_same = true;
    } break;

    case SGSAP_LOCATION_UPDATE_REJ: {
      /*Received SGSAP Location Update Reject message from SGS task*/
      mme_app_handle_sgsap_location_update_rej(
          mme_app_desc_p,
          &received_message_p->ittiMsg.sgsap_location_update_rej);
    } break;

    case SGSAP_ALERT_REQUEST: {
      /*Received SGSAP Alert Request message from SGS task*/
      mme_app_handle_sgsap_alert_request(
          mme_app_desc_p, &received_message_p->ittiMsg.sgsap_alert_request);
      is_task_state_same = true;
    } break;

    case SGSAP_VLR_RESET_INDICATION: {
      /*Received SGSAP Reset Indication from SGS task*/
      mme_app_handle_sgsap_reset_indication(
          &received_message_p->ittiMsg.sgsap_vlr_reset_indication);
      is_task_state_same = true;
    } break;

    case SGSAP_PAGING_REQUEST: {
      mme_app_handle_sgsap_paging_request(
          mme_app_desc_p, &received_message_p->ittiMsg.sgsap_paging_request);
      is_task_state_same = true;
    } break;

    case SGSAP_SERVICE_ABORT_REQ: {
      mme_app_handle_sgsap_service_abort_request(
          mme_app_desc_p, &received_message_p->ittiMsg.sgsap_service_abort_req);
      is_task_state_same = true;
    } break;

    case SGSAP_EPS_DETACH_ACK: {
      mme_app_handle_sgs_eps_detach_ack(
          mme_app_desc_p, &received_message_p->ittiMsg.sgsap_eps_detach_ack);
      is_task_state_same = true;
    } break;

    case SGSAP_IMSI_DETACH_ACK: {
      mme_app_handle_sgs_imsi_detach_ack(
          mme_app_desc_p, &received_message_p->ittiMsg.sgsap_imsi_detach_ack);
      is_task_state_same = true;
    } break;

    case S11_MODIFY_UE_AMBR_REQUEST: {
      mme_app_handle_modify_ue_ambr_request(
          mme_app_desc_p, &S11_MODIFY_UE_AMBR_REQUEST(received_message_p));
      is_task_state_same = true;
    } break;

    case S11_NW_INITIATED_ACTIVATE_BEARER_REQUEST: {
      mme_app_handle_nw_init_ded_bearer_actv_req(
          mme_app_desc_p,
          &received_message_p->ittiMsg.s11_nw_init_actv_bearer_request);
    } break;

    case SGSAP_STATUS: {
      mme_app_handle_sgs_status_message(
          mme_app_desc_p, &received_message_p->ittiMsg.sgsap_status);
      is_task_state_same = true;
    } break;

    case S11_NW_INITIATED_DEACTIVATE_BEARER_REQUEST: {
      mme_app_handle_nw_init_bearer_deactv_req(
          mme_app_desc_p,
          &received_message_p->ittiMsg.s11_nw_init_deactv_bearer_request);
      is_task_state_same = true;
    } break;

    case S1AP_PATH_SWITCH_REQUEST: {
      mme_app_handle_path_switch_request(
          mme_app_desc_p, &S1AP_PATH_SWITCH_REQUEST(received_message_p));
    } break;

    case S1AP_HANDOVER_REQUIRED: {
      mme_app_handle_handover_required(
          &S1AP_HANDOVER_REQUIRED(received_message_p));
      is_task_state_same = true;
    } break;

    case S1AP_HANDOVER_REQUEST_ACK: {
      mme_app_handle_handover_request_ack(
          &S1AP_HANDOVER_REQUEST_ACK(received_message_p));
      is_task_state_same = true;
    } break;

    case S1AP_HANDOVER_NOTIFY: {
      mme_app_handle_handover_notify(
          mme_app_desc_p, &S1AP_HANDOVER_NOTIFY(received_message_p));
    } break;

    case S6A_AUTH_INFO_ANS: {
      /*
       * We received the authentication vectors from HSS,
       * Normally should trigger an authentication procedure towards UE.
       */
      nas_proc_authentication_info_answer(
          mme_app_desc_p, &S6A_AUTH_INFO_ANS(received_message_p));
    } break;

    case MME_APP_DOWNLINK_DATA_CNF: {
      bstring nas_msg = NULL;
      nas_proc_dl_transfer_cnf(
          MME_APP_DL_DATA_CNF(received_message_p).ue_id,
          MME_APP_DL_DATA_CNF(received_message_p).err_code, &nas_msg);
      is_task_state_same = true;
    } break;

    case MME_APP_DOWNLINK_DATA_REJ: {
      nas_proc_dl_transfer_rej(
          MME_APP_DL_DATA_REJ(received_message_p).ue_id,
          MME_APP_DL_DATA_REJ(received_message_p).err_code,
          &MME_APP_DL_DATA_REJ(received_message_p).nas_msg);
      is_task_state_same = true;
    } break;

    case SGSAP_DOWNLINK_UNITDATA: {
      /* We received the Downlink Unitdata from MSC, trigger a
       * Downlink Nas Transport message to UE.
       */
      nas_proc_downlink_unitdata(&SGSAP_DOWNLINK_UNITDATA(received_message_p));
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

    case S1AP_REMOVE_STALE_UE_CONTEXT: {
      mme_app_remove_stale_ue_context(
          mme_app_desc_p, &S1AP_REMOVE_STALE_UE_CONTEXT(received_message_p));
    } break;

    case TERMINATE_MESSAGE: {
      itti_free_msg_content(received_message_p);
      free(received_message_p);
      mme_app_exit();
    } break;

    case RECOVERY_MESSAGE: {
      OAILOG_INFO(LOG_MME_APP, "Received RECOVERY_MESSAGE \n");
      mme_app_recover_timers_for_all_ues();
    } break;

    default: {
      OAILOG_ERROR(
          LOG_MME_APP, "Unknown message (%s) received with message Id: %d\n",
          ITTI_MSG_NAME(received_message_p), ITTI_MSG_ID(received_message_p));
    } break;
  }

  put_mme_ue_state(mme_app_desc_p, imsi64, force_ue_write);

  if (!is_task_state_same) {
    put_mme_nas_state();
  }

  itti_free_msg_content(received_message_p);
  free(received_message_p);
  return 0;
}

//------------------------------------------------------------------------------
static void* mme_app_thread(__attribute__((unused)) void* args) {
  itti_mark_task_ready(TASK_MME_APP);
  init_task_context(
      TASK_MME_APP,
      (task_id_t[]){TASK_SPGW_APP, TASK_SGS, TASK_SMS_ORC8R, TASK_S11, TASK_S6A,
                    TASK_S1AP, TASK_SERVICE303, TASK_HA, TASK_SGW_S8},
      9, handle_message, &mme_app_task_zmq_ctx);

  // Service started, but not healthy yet
  send_app_health_to_service303(&mme_app_task_zmq_ctx, TASK_MME_APP, false);
  start_stats_timer();

  zloop_start(mme_app_task_zmq_ctx.event_loop);
  AssertFatal(
      0, "Asserting as mme_app_thread should not be exiting on its own!");
  return NULL;
}

static void mme_app_init_congestion_params(const mme_config_t* mme_config_p) {
  mme_congestion_params.mme_app_zmq_congest_th =
      (long) mme_config_p->mme_app_zmq_congest_th;
  mme_congestion_params.mme_app_zmq_auth_th =
      (long) mme_config_p->mme_app_zmq_auth_th;
  mme_congestion_params.mme_app_zmq_ident_th =
      (long) mme_config_p->mme_app_zmq_ident_th;
  mme_congestion_params.mme_app_zmq_smc_th =
      (long) mme_config_p->mme_app_zmq_smc_th;
}

//------------------------------------------------------------------------------
status_code_e mme_app_init(const mme_config_t* mme_config_p) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  if (mme_nas_state_init(mme_config_p)) {
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  if (mme_app_edns_init(mme_config_p)) {
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  // Initialise NAS module
  nas_network_initialize(mme_config_p);

  // Initialize task global congestion parameters
  mme_congestion_control_enabled = mme_config_p->enable_congestion_control;
  mme_app_init_congestion_params(mme_config_p);

  // Initialize global stats timer
  epc_stats_timer_sec = (size_t) mme_config_p->stats_timer_sec;

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

static int handle_stats_timer(zloop_t* loop, int id, void* arg) {
  mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
  application_mme_app_stats_msg_t stats_msg;
  stats_msg.nb_ue_attached         = mme_app_desc_p->nb_ue_attached;
  stats_msg.nb_ue_connected        = mme_app_desc_p->nb_ue_connected;
  stats_msg.nb_default_eps_bearers = mme_app_desc_p->nb_default_eps_bearers;
  stats_msg.nb_s1u_bearers         = mme_app_desc_p->nb_s1u_bearers;
  stats_msg.nb_mme_app_last_msg_latency = mme_app_last_msg_latency;

  return send_mme_app_stats_to_service303(
      &mme_app_task_zmq_ctx, TASK_MME_APP, &stats_msg);
}

static void start_stats_timer(void) {
  epc_stats_timer_id = start_timer(
      &mme_app_task_zmq_ctx, 1000 * epc_stats_timer_sec, TIMER_REPEAT_FOREVER,
      handle_stats_timer, NULL);
}

static void check_mme_healthy_and_notify_service(void) {
  if (is_mme_app_healthy()) {
    send_app_health_to_service303(&mme_app_task_zmq_ctx, TASK_MME_APP, true);
  }
}

static bool is_mme_app_healthy(void) {
  return mme_hss_associated && mme_sctp_bounded;
}

//------------------------------------------------------------------------------
static void mme_app_exit(void) {
  stop_timer(&mme_app_task_zmq_ctx, epc_stats_timer_id);
  mme_app_edns_exit();
  clear_mme_nas_state();
  // Clean-up NAS module
  nas_network_cleanup();
  mme_config_exit();
  destroy_task_context(&mme_app_task_zmq_ctx);
  OAI_FPRINTF_INFO("TASK_MME_APP terminated\n");
  pthread_exit(NULL);
}
