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

/*! \file mme_app_location.c
   \brief
   \author Sebastien ROUX, Lionel GAUTHIER
   \version 1.0
   \company Eurecom
   \email: lionel.gauthier@eurecom.fr
*/

#include <stdio.h>
#include <string.h>
#include <stdbool.h>
#include <stdint.h>
#include <inttypes.h>

#include "bstrlib.h"
#include "log.h"
#include "common_types.h"
#include "conversions.h"
#include "intertask_interface.h"
#include "common_defs.h"
#include "common_utility_funs.h"
#include "mme_config.h"
#include "mme_app_ue_context.h"
#include "mme_app_defs.h"
#include "timer.h"
#include "3gpp_23.003.h"
#include "3gpp_36.401.h"
#include "TrackingAreaIdentity.h"
#include "emm_data.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "mme_app_desc.h"
#include "s6a_messages_types.h"
#include "service303.h"
#include "sgs_messages_types.h"
#include "esm_proc.h"
#include "nas_proc.h"
#include "emm_cnDef.h"
#include "emm_proc.h"
#include "mme_app_timer.h"
#include "dynamic_memory_check.h"

//------------------------------------------------------------------------------
int mme_app_send_s6a_update_location_req(
    struct ue_mm_context_s* const ue_context_p) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  MessageDef* message_p                = NULL;
  s6a_update_location_req_t* s6a_ulr_p = NULL;
  int rc                               = RETURNok;

  OAILOG_INFO(
      TASK_MME_APP,
      "Sending S6A UPDATE LOCATION REQ to S6A, ue_id =" MME_UE_S1AP_ID_FMT
      " \n",
      ue_context_p->mme_ue_s1ap_id);
  message_p = itti_alloc_new_message(TASK_MME_APP, S6A_UPDATE_LOCATION_REQ);
  if (message_p == NULL) {
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  s6a_ulr_p = &message_p->ittiMsg.s6a_update_location_req;
  memset((void*) s6a_ulr_p, 0, sizeof(s6a_update_location_req_t));
  IMSI64_TO_STRING(
      (ue_context_p->emm_context._imsi64), s6a_ulr_p->imsi,
      ue_context_p->emm_context._imsi.length);
  s6a_ulr_p->imsi_length    = strlen(s6a_ulr_p->imsi);
  s6a_ulr_p->initial_attach = INITIAL_ATTACH;
  plmn_t visited_plmn       = {0};
  COPY_PLMN(visited_plmn, ue_context_p->emm_context.originating_tai.plmn);
  memcpy(&s6a_ulr_p->visited_plmn, &visited_plmn, sizeof(plmn_t));
  s6a_ulr_p->rat_type = RAT_EUTRAN;
  OAILOG_DEBUG(
      TASK_MME_APP, "S6A ULR: RAT TYPE = (%d) for (ue_id = %u)\n",
      s6a_ulr_p->rat_type, ue_context_p->mme_ue_s1ap_id);

  // Set regional_subscription flag
  s6a_ulr_p->supportedfeatures.regional_subscription = true;
  /*
   * Check if we already have UE data
   * set the skip subscriber data flag as true in case we are sending ULR
   * against recieved HSS Reset
   */
  if (ue_context_p->location_info_confirmed_in_hss == true) {
    s6a_ulr_p->skip_subscriber_data = 1;
    OAILOG_DEBUG(
        TASK_MME_APP,
        "S6A Location information confirmed in HSS (%d) for (ue_id = %u)\n",
        ue_context_p->location_info_confirmed_in_hss,
        ue_context_p->mme_ue_s1ap_id);
  } else {
    s6a_ulr_p->skip_subscriber_data = 0;
    OAILOG_DEBUG(
        TASK_MME_APP,
        "S6A Location information not confirmed in HSS (%d) for (ue_id = %u)\n",
        ue_context_p->location_info_confirmed_in_hss,
        ue_context_p->mme_ue_s1ap_id);
  }

  /*
   * Check if we have UE 5G-NR
   * connection supported in Attach request message.
   * This is done by checking either en_dc flag in ms network capability or
   * by checking  dcnr flag in ue network capability.
   */

  if (ue_context_p->emm_context._ue_network_capability.dcnr) {
    s6a_ulr_p->dual_regis_5g_ind = 1;
  } else {
    s6a_ulr_p->dual_regis_5g_ind = 0;
    OAILOG_DEBUG(
        TASK_MME_APP,
        "UE is connected on LTE, Dual registration with 5G-NR is disabled for "
        "(ue_id = %u)\n",
        ue_context_p->mme_ue_s1ap_id);
  }

  // Check if we have voice domain preference IE and send to S6a task
  if (ue_context_p->emm_context.volte_params.presencemask &
      VOICE_DOMAIN_PREF_UE_USAGE_SETTING) {
    s6a_ulr_p->voice_dom_pref_ue_usg_setting =
        ue_context_p->emm_context.volte_params
            .voice_domain_preference_and_ue_usage_setting;
    s6a_ulr_p->presencemask |= S6A_PDN_CONFIG_VOICE_DOM_PREF;
  }
  OAILOG_DEBUG(
      LOG_MME_APP,
      "0 S6A_UPDATE_LOCATION_REQ imsi %s with length %d for (ue_id = %u)\n",
      s6a_ulr_p->imsi, s6a_ulr_p->imsi_length, ue_context_p->mme_ue_s1ap_id);
  rc = send_msg_to_task(&mme_app_task_zmq_ctx, TASK_S6A, message_p);
  /*
   * Do not start this timer in case we are sending ULR after receiving HSS
   * reset
   */
  if (ue_context_p->location_info_confirmed_in_hss == false) {
    // Start ULR Response timer
    if ((ue_context_p->ulr_response_timer.id = mme_app_start_timer(
             ue_context_p->ulr_response_timer.sec * 1000, TIMER_REPEAT_ONCE,
             mme_app_handle_ulr_timer_expiry, ue_context_p->mme_ue_s1ap_id)) ==
        -1) {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Failed to start Update location update response timer for UE id "
          " " MME_UE_S1AP_ID_FMT "\n",
          ue_context_p->mme_ue_s1ap_id);
      ue_context_p->ulr_response_timer.id = MME_APP_TIMER_INACTIVE_ID;
    } else {
      OAILOG_DEBUG(
          LOG_MME_APP,
          "Started location update response timer for UE id  %d \n",
          ue_context_p->mme_ue_s1ap_id);
    }
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

int handle_ula_failure(struct ue_mm_context_s* ue_context_p) {
  int rc = RETURNok;

  OAILOG_FUNC_IN(LOG_MME_APP);

  // Stop ULR Response timer if running
  if (ue_context_p->ulr_response_timer.id != MME_APP_TIMER_INACTIVE_ID) {
    mme_app_stop_timer(ue_context_p->ulr_response_timer.id);
    ue_context_p->ulr_response_timer.id = MME_APP_TIMER_INACTIVE_ID;
  }
  increment_counter("mme_s6a_update_location_ans", 1, 1, "result", "failure");
  emm_cn_ula_or_csrsp_fail_t cn_ula_fail = {0};
  if (ue_context_p->emm_context.esm_ctx.esm_proc_data) {
    cn_ula_fail.pti = ue_context_p->emm_context.esm_ctx.esm_proc_data->pti;
  } else {
    OAILOG_ERROR(
        LOG_MME_APP, " esm_proc_data is NULL, so failed to fetch pti \n");
  }
  cn_ula_fail.ue_id = ue_context_p->mme_ue_s1ap_id;
  cn_ula_fail.cause = CAUSE_SYSTEM_FAILURE;
  rc                = nas_proc_ula_or_csrsp_fail(&cn_ula_fail);
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

//------------------------------------------------------------------------------
int mme_app_handle_s6a_update_location_ans(
    mme_app_desc_t* mme_app_desc_p,
    const s6a_update_location_ans_t* const ula_pP) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  uint64_t imsi64                       = 0;
  struct ue_mm_context_s* ue_mm_context = NULL;
  int rc                                = RETURNok;

  if (ula_pP == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Invalid S6a Update Location Answer ITTI message received\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  IMSI_STRING_TO_IMSI64((char*) ula_pP->imsi, &imsi64);
  OAILOG_DEBUG(LOG_MME_APP, "Handling imsi " IMSI_64_FMT "\n", imsi64);

  if ((ue_mm_context = mme_ue_context_exists_imsi(
           &mme_app_desc_p->mme_ue_contexts, imsi64)) == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "That's embarrassing as we don't know this IMSI " IMSI_64_FMT "\n",
        imsi64);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  if (ula_pP->result.present == S6A_RESULT_BASE) {
    if (ula_pP->result.choice.base != DIAMETER_SUCCESS) {
      /*
       * The update location procedure has failed. Notify the NAS module
       * and don't initiate the bearer creation on S-GW side.
       */
      OAILOG_ERROR(
          LOG_MME_APP,
          "ULR/ULA procedure returned non success "
          "(ULA.result.choice.base=%d)\n",
          ula_pP->result.choice.base);
      if (handle_ula_failure(ue_mm_context) != RETURNok) {
        OAILOG_ERROR(
            LOG_MME_APP,
            "Failed to handle Un-successful ULA message for "
            "ue_id " MME_UE_S1AP_ID_FMT "\n",
            ue_mm_context->mme_ue_s1ap_id);
        OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
      } else {
        OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
      }
    }
  } else {
    /*
     * The update location procedure has failed. Notify the NAS layer
     * and don't initiate the bearer creation on S-GW side.
     */
    OAILOG_ERROR(
        LOG_MME_APP,
        "ULR/ULA procedure returned non success (ULA.result.present=%d)\n",
        ula_pP->result.present);
    if (handle_ula_failure(ue_mm_context) == RETURNok) {
      OAILOG_DEBUG(
          LOG_MME_APP, "Sent PDN Connectivity failure to NAS for ue_id (%u)\n",
          ue_mm_context->mme_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
    } else {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Failed to send PDN Connectivity failure to NAS for "
          "ue_id " MME_UE_S1AP_ID_FMT "\n",
          ue_mm_context->mme_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
    }
  }

  // Stop ULR Response timer.
  // If expired its timer id should be MME_APP_TIMER_INACTIVE_ID and
  // it should be already treated as failure
  if (ue_mm_context->ulr_response_timer.id != MME_APP_TIMER_INACTIVE_ID) {
    mme_app_stop_timer(ue_mm_context->ulr_response_timer.id);
    ue_mm_context->ulr_response_timer.id = MME_APP_TIMER_INACTIVE_ID;
  } else {
    OAILOG_ERROR(
        LOG_MME_APP,
        "ULR Response Timer has invalid id for ue_id " MME_UE_S1AP_ID_FMT
        ". This implies that "
        "the timer has expired and ULR has been handled as failure. \n ",
        ue_mm_context->mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  ue_mm_context->subscription_known = SUBSCRIPTION_KNOWN;
  ue_mm_context->subscriber_status =
      ula_pP->subscription_data.subscriber_status;

  // Verify service area restriction
  if (ula_pP->subscription_data.num_zcs > 0) {
    if (verify_service_area_restriction(
            ue_mm_context->emm_context.originating_tai.tac,
            ula_pP->subscription_data.reg_sub,
            ula_pP->subscription_data.num_zcs) != RETURNok) {
      OAILOG_ERROR_UE(
          LOG_MME_APP, imsi64,
          "No suitable cells found for tac = %d, sending attach_reject "
          "message "
          "for ue_id " MME_UE_S1AP_ID_FMT " with emm cause = %d\n",
          ue_mm_context->emm_context.originating_tai.tac,
          ue_mm_context->mme_ue_s1ap_id, EMM_CAUSE_NO_SUITABLE_CELLS);
      if (emm_proc_attach_reject(
              ue_mm_context->mme_ue_s1ap_id, EMM_CAUSE_NO_SUITABLE_CELLS) !=
          RETURNok) {
        OAILOG_ERROR_UE(
            LOG_MME_APP, imsi64,
            "Sending of attach reject message failed for "
            "ue_id " MME_UE_S1AP_ID_FMT "\n",
            ue_mm_context->mme_ue_s1ap_id);
        OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
      }
      OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
    } else {
      // Store the zone codes in ue_mm_context
      ue_mm_context->num_reg_sub = ula_pP->subscription_data.num_zcs;
      for (uint8_t itr = 0; itr < ue_mm_context->num_reg_sub; itr++) {
        memcpy(
            ue_mm_context->reg_sub[itr].zone_code,
            ula_pP->subscription_data.reg_sub[itr].zone_code,
            strlen((const char*) ula_pP->subscription_data.reg_sub[itr]
                       .zone_code));
      }
    }
  }
  ue_mm_context->access_restriction_data =
      ula_pP->subscription_data.access_restriction;
  /*
   * Copy the subscribed UE AMBR (comes from data plan) to UE context
   */
  memcpy(
      &ue_mm_context->subscribed_ue_ambr,
      &ula_pP->subscription_data.subscribed_ambr, sizeof(ambr_t));
  OAILOG_DEBUG(
      LOG_MME_APP, "Received UL rate %" PRIu64 " and DL rate %" PRIu64 "\n",
      ue_mm_context->subscribed_ue_ambr.br_ul,
      ue_mm_context->subscribed_ue_ambr.br_dl);

  if (ula_pP->subscription_data.msisdn_length != 0) {
    ue_mm_context->msisdn = blk2bstr(
        ula_pP->subscription_data.msisdn,
        ula_pP->subscription_data.msisdn_length);
  } else {
    OAILOG_WARNING(
        LOG_MME_APP, "No MSISDN received for %s " IMSI_64_FMT "\n",
        __FUNCTION__, imsi64);
  }

  ue_mm_context->rau_tau_timer       = ula_pP->subscription_data.rau_tau_timer;
  ue_mm_context->network_access_mode = ula_pP->subscription_data.access_mode;
  memcpy(
      &ue_mm_context->apn_config_profile,
      &ula_pP->subscription_data.apn_config_profile,
      sizeof(apn_config_profile_t));
  memcpy(
      &ue_mm_context->default_charging_characteristics,
      &ula_pP->subscription_data.default_charging_characteristics,
      sizeof(charging_characteristics_t));

  /*
   * Set the value of  Mobile Reachability timer based on value of T3412
   * (Periodic TAU timer) sent in Attach accept /TAU accept. Set it to
   * MME_APP_DELTA_T3412_REACHABILITY_TIMER minutes greater than T3412. Set the
   * value of Implicit timer. Set it to
   * MME_APP_DELTA_REACHABILITY_IMPLICIT_DETACH_TIMER minutes greater than
   * Mobile Reachability timer
   */
  ue_mm_context->mobile_reachability_timer.id = MME_APP_TIMER_INACTIVE_ID;
  ue_mm_context->mobile_reachability_timer.sec =
      ((mme_config.nas_config.t3412_min) +
       MME_APP_DELTA_T3412_REACHABILITY_TIMER) *
      60;
  ue_mm_context->implicit_detach_timer.id = MME_APP_TIMER_INACTIVE_ID;
  ue_mm_context->implicit_detach_timer.sec =
      (ue_mm_context->mobile_reachability_timer.sec) +
      MME_APP_DELTA_REACHABILITY_IMPLICIT_DETACH_TIMER * 60;

  /*
   * Set the flag: send_ue_purge_request to indicate that
   * Update Location procedure is completed.
   * During UE initiated detach/Implicit detach this MME would send PUR to hss,
   * if this flag is true.
   */
  ue_mm_context->send_ue_purge_request = true;
  /*
   * Set the flag: location_info_confirmed_in_hss to false to indicate that
   * Update Location procedure is completed.
   * During HSS Reset
   * if this flag is true.
   */
  if (ue_mm_context->location_info_confirmed_in_hss == true) {
    ue_mm_context->location_info_confirmed_in_hss = false;
  }
  rc = nas_proc_ula_success(ue_mm_context->mme_ue_s1ap_id);

  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

int mme_app_handle_s6a_cancel_location_req(
    mme_app_desc_t* mme_app_desc_p,
    const s6a_cancel_location_req_t* const clr_pP) {
  uint64_t imsi                        = 0;
  struct ue_mm_context_s* ue_context_p = NULL;
  int cla_result                       = DIAMETER_SUCCESS;

  OAILOG_FUNC_IN(LOG_MME_APP);
  if (clr_pP == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Invalid S6a Cancel Location Request ITTI message received\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  IMSI_STRING_TO_IMSI64((char*) clr_pP->imsi, &imsi);
  OAILOG_DEBUG(
      LOG_MME_APP, "S6a Cancel Location Request for imsi " IMSI_64_FMT "\n",
      imsi);

  if ((mme_app_send_s6a_cancel_location_ans(
          cla_result, clr_pP->imsi, clr_pP->imsi_length, clr_pP->msg_cla_p)) !=
      RETURNok) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "S6a Cancel Location Request: Failed to send Cancel Location Answer "
        "from "
        "MME app for imsi " IMSI_64_FMT "\n",
        imsi);
  }

  if ((ue_context_p = mme_ue_context_exists_imsi(
           &mme_app_desc_p->mme_ue_contexts, imsi)) == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "IMSI is not present in the MME context for imsi " IMSI_64_FMT "\n",
        imsi);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  if (clr_pP->cancellation_type != SUBSCRIPTION_WITHDRAWL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "S6a Cancel Location Request: Cancellation_type not supported %d for"
        "imsi " IMSI_64_FMT "\n",
        clr_pP->cancellation_type, imsi);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  /*
   * set the flag: hss_initiated_detach to indicate that,
   * hss has initiated detach and MME shall not send PUR to hss
   */
  ue_context_p->hss_initiated_detach = true;

  /*
   * Check UE's S1 connection status.If UE is in connected state,
   * send Detach Request to UE. If UE is in idle state,
   * Page UE to bring it back to connected mode and then send Detach Request
   */
  if (ue_context_p->ecm_state == ECM_IDLE) {
    /* Page the UE to bring it back to connected mode
     * and then send Detach Request
     */
    mme_app_paging_request_helper(
        ue_context_p, true, false /* s-tmsi */, CN_DOMAIN_PS);
    // Set the flag and send detach to UE after receiving service req
    ue_context_p->emm_context.nw_init_bearer_deactv = true;
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
  } else {
    // Send N/W Initiated Detach Request to NAS module

    OAILOG_INFO(
        LOG_MME_APP,
        "Sending Detach Req to NAS module for (ue_id = " MME_UE_S1AP_ID_FMT
        "\n",
        ue_context_p->mme_ue_s1ap_id);
    if ((mme_app_handle_nw_initiated_detach_request(
            ue_context_p->mme_ue_s1ap_id, HSS_INITIATED_EPS_DETACH)) ==
        RETURNerror) {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Failed to handle network initiated Detach Request in nas module for "
          "ue-id: " MME_UE_S1AP_ID_FMT "\n",
          ue_context_p->mme_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
    } else {
      // Send SGS explicit network initiated Detach Ind to SGS
      if (ue_context_p->sgs_context) {
        mme_app_handle_sgs_detach_req(
            ue_context_p, EMM_SGS_NW_INITIATED_EPS_DETACH);
      }
    }
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

int mme_app_send_s6a_cancel_location_ans(
    int cla_result, const char* imsi, uint8_t imsi_length, void* msg_cla_p) {
  MessageDef* message_p                = NULL;
  s6a_cancel_location_ans_t* s6a_cla_p = NULL;
  int rc                               = RETURNok;

  OAILOG_FUNC_IN(LOG_MME_APP);

  message_p = itti_alloc_new_message(TASK_MME_APP, S6A_CANCEL_LOCATION_ANS);

  if (message_p == NULL) {
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  s6a_cla_p = &message_p->ittiMsg.s6a_cancel_location_ans;
  memset((void*) s6a_cla_p, 0, sizeof(s6a_cancel_location_ans_t));

  /* Using the IMSI details deom CLR */
  memcpy(s6a_cla_p->imsi, imsi, imsi_length);
  s6a_cla_p->imsi_length = imsi_length;

  s6a_cla_p->result    = cla_result;
  s6a_cla_p->msg_cla_p = msg_cla_p;
  rc = send_msg_to_task(&mme_app_task_zmq_ctx, TASK_S6A, message_p);
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

void mme_app_get_user_location_information(
    Uli_t* uli_t_p, const ue_mm_context_t* ue_context_p) {
  uli_t_p->present = uli_t_p->present | ULI_TAI;
  COPY_TAI(uli_t_p->s.tai, ue_context_p->emm_context.originating_tai);
  uli_t_p->present = uli_t_p->present | ULI_ECGI;
  COPY_PLMN(uli_t_p->s.ecgi.plmn, ue_context_p->e_utran_cgi.plmn);
  uli_t_p->s.ecgi.cell_identity.enb_id =
      ue_context_p->e_utran_cgi.cell_identity.enb_id;
  uli_t_p->s.ecgi.cell_identity.cell_id =
      ue_context_p->e_utran_cgi.cell_identity.cell_id;
}
