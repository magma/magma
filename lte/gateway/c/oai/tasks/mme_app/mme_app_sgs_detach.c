/*******************************************************************************
    OpenAirInterface
    Copyright(c) 1999 - 2014 Eurecom

    OpenAirInterface is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.


    OpenAirInterface is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with OpenAirInterface.The full GNU General Public License is
   included in this distribution in the file called "COPYING". If not,
   see <http://www.gnu.org/licenses/>.

  Contact Information
  OpenAirInterface Admin: openair_admin@eurecom.fr
  OpenAirInterface Tech : openair_tech@eurecom.fr
  OpenAirInterface Dev  : openair4g-devel@eurecom.fr

  Address      : Eurecom, Compus SophiaTech 450, route des chappes, 06451 Biot,
 France.

 *******************************************************************************/

#include <stdio.h>
#include <string.h>
#include <stdint.h>

#include "dynamic_memory_check.h"
#include "intertask_interface.h"
#include "mme_app_ue_context.h"
#include "mme_app_itti_messaging.h"
#include "mme_app_defs.h"
#include "timer.h"
#include "conversions.h"
#include "service303.h"
#include "3gpp_36.401.h"
#include "common_defs.h"
#include "common_types.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "log.h"
#include "mme_app_state.h"
#include "mme_app_sgs_fsm.h"
#include "s1ap_messages_types.h"
#include "sgs_messages_types.h"
#include "mme_app_timer.h"

/**
 * Function to send a SGS EPS detach indication to SGSAP in either the initial
 * case or the retransmission case.
 *
 * @param ue_context - ue_context pointer
 * @param sgs_detach_type - SGS EPS detach type
 */
static void mme_app_send_sgs_eps_detach_indication(
    ue_mm_context_t* ue_context_p, uint8_t detach_type) {
  MessageDef* message_p = NULL;

  OAILOG_FUNC_IN(LOG_MME_APP);
  OAILOG_INFO(
      LOG_MME_APP,
      "Send SGSAP_EPS_DETACH_IND to SGS, detach_type = %u for "
      "ue_id: " MME_UE_S1AP_ID_FMT "\n",
      detach_type, ue_context_p->mme_ue_s1ap_id);
  message_p = itti_alloc_new_message(TASK_MME_APP, SGSAP_EPS_DETACH_IND);
  if (!message_p) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Failed to allocate memory for SGSAP_EPS_DETACH_IND for "
        "ue_id: " MME_UE_S1AP_ID_FMT "\n",
        ue_context_p->mme_ue_s1ap_id);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  memset(
      (void*) &message_p->ittiMsg.sgsap_eps_detach_ind, 0,
      sizeof(itti_sgsap_eps_detach_ind_t));

  IMSI64_TO_STRING(
      ue_context_p->emm_context._imsi64, SGSAP_EPS_DETACH_IND(message_p).imsi,
      ue_context_p->emm_context._imsi.length);
  SGSAP_EPS_DETACH_IND(message_p).imsi_length =
      (uint8_t) strlen(SGSAP_IMSI_DETACH_IND(message_p).imsi);

  SGSAP_EPS_DETACH_IND(message_p).eps_detach_type = detach_type;

  send_msg_to_task(&mme_app_task_zmq_ctx, TASK_SGS, message_p);

  // Start SGS Implicit EPS Detach indication timer
  if (detach_type == SGS_NW_INITIATED_IMSI_DETACH_FROM_EPS) {
    if ((ue_context_p->sgs_context->ts13_timer.id = mme_app_start_timer(
             ue_context_p->sgs_context->ts13_timer.sec * 1000,
             TIMER_REPEAT_ONCE,
             mme_app_handle_sgs_implicit_eps_detach_timer_expiry,
             ue_context_p->mme_ue_s1ap_id)) == -1) {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Failed to start SGS Implicit EPS Detach indication timer for "
          "ue_id " MME_UE_S1AP_ID_FMT "\n",
          ue_context_p->mme_ue_s1ap_id);
      ue_context_p->sgs_context->ts13_timer.id = MME_APP_TIMER_INACTIVE_ID;
    } else {
      OAILOG_DEBUG(
          LOG_MME_APP,
          "Started SGS Implicit EPS Detach indication timer for "
          "ue_id " MME_UE_S1AP_ID_FMT "\n",
          ue_context_p->mme_ue_s1ap_id);
    }
  } else {
    // Start SGS EPS Detach indication timer
    if ((ue_context_p->sgs_context->ts8_timer.id = mme_app_start_timer(
             ue_context_p->sgs_context->ts8_timer.sec * 1000, TIMER_REPEAT_ONCE,
             mme_app_handle_sgs_eps_detach_timer_expiry,
             ue_context_p->mme_ue_s1ap_id)) == -1) {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Failed to start SGS EPS Detach indication timer for "
          "ue_id " MME_UE_S1AP_ID_FMT "\n",
          ue_context_p->mme_ue_s1ap_id);
      ue_context_p->sgs_context->ts8_timer.id = MME_APP_TIMER_INACTIVE_ID;
    } else {
      OAILOG_DEBUG(
          LOG_MME_APP,
          "Started SGS EPS Detach indication timer for "
          "ue_id " MME_UE_S1AP_ID_FMT "\n",
          ue_context_p->mme_ue_s1ap_id);
    }

    // Stop  and reset SGS Location Update Request timer if running
    if (ue_context_p->sgs_context->ts6_1_timer.id !=
        MME_APP_TIMER_INACTIVE_ID) {
      mme_app_stop_timer(ue_context_p->sgs_context->ts6_1_timer.id);
      ue_context_p->sgs_context->ts6_1_timer.id = MME_APP_TIMER_INACTIVE_ID;
    }
  }

  increment_counter(
      "mme_sgs_eps_detach_indication_sent", 1, 1, "result", "success");
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

// handle the SGS EPS detach timer expiry
int mme_app_handle_sgs_eps_detach_timer_expiry(
    zloop_t* loop, int timer_id, void* args) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  mme_ue_s1ap_id_t mme_ue_s1ap_id = 0;
  if (!mme_app_get_timer_arg(timer_id, &mme_ue_s1ap_id)) {
    OAILOG_WARNING(
        LOG_MME_APP, "Invalid Timer Id expiration, Timer Id: %u\n", timer_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  struct ue_mm_context_s* ue_context_p =
      mme_app_get_ue_context_for_timer(mme_ue_s1ap_id, "sgs eps detach timer");
  if (ue_context_p == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Invalid UE context received, MME UE S1AP Id: " MME_UE_S1AP_ID_FMT "\n",
        mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  if (ue_context_p->sgs_context == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "SGS EPS Detach Timer expired but no assoicated SGS context for UE "
        "id " MME_UE_S1AP_ID_FMT "\n",
        mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  /*
   * Increment the retransmission counter
   */
  ue_context_p->sgs_context->ts8_retransmission_count += 1;
  OAILOG_WARNING(
      LOG_MME_APP,
      "MME APP: Ts8 timer expired,retransmission "
      "counter = %u \n",
      ue_context_p->sgs_context->ts8_retransmission_count);

  ue_context_p->sgs_context->ts8_timer.id = MME_APP_TIMER_INACTIVE_ID;
  if (ue_context_p->sgs_context->ts8_retransmission_count <
      EPS_DETACH_RETRANSMISSION_COUNTER_MAX) {
    /*
     * Resend SGS EPS Detach Indication message to the SGS
     */
    mme_app_send_sgs_eps_detach_indication(
        ue_context_p, ue_context_p->sgs_detach_type);
  } else {
    OAILOG_DEBUG(
        LOG_MME_APP,
        "SGS EPS DETACH indication failed after %u retransmission and expiry "
        "\n",
        ue_context_p->sgs_context->ts8_retransmission_count);
    increment_counter(
        "sgs_eps_detach_timer_expired", 1, 1, "cause",
        "Ts8 timer expired after max tetransmission");
  }

  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

// handle the SGS Implicit EPS detach timer expiry
int mme_app_handle_sgs_implicit_eps_detach_timer_expiry(
    zloop_t* loop, int timer_id, void* args) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  mme_ue_s1ap_id_t mme_ue_s1ap_id = 0;
  if (!mme_app_get_timer_arg(timer_id, &mme_ue_s1ap_id)) {
    OAILOG_WARNING(
        LOG_MME_APP, "Invalid Timer Id expiration, Timer Id: %u\n", timer_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  struct ue_mm_context_s* ue_context_p = mme_app_get_ue_context_for_timer(
      mme_ue_s1ap_id, "sgs implicit eps detach timer");
  if (ue_context_p == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Invalid UE context received, MME UE S1AP Id: " MME_UE_S1AP_ID_FMT "\n",
        mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  if (ue_context_p->sgs_context == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "SGS EPS Detach Timer expired but no assoicated SGS context for UE "
        "id " MME_UE_S1AP_ID_FMT "\n",
        mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  /*
   * Increment the retransmission counter
   */
  ue_context_p->sgs_context->ts13_retransmission_count += 1;
  OAILOG_ERROR(
      LOG_NAS_EMM,
      "MME APP: Ts13 timer expired,retransmission "
      "counter = %u \n",
      ue_context_p->sgs_context->ts13_retransmission_count);

  ue_context_p->sgs_context->ts13_timer.id = MME_APP_TIMER_INACTIVE_ID;
  if (ue_context_p->sgs_context->ts13_retransmission_count <
      IMPLICIT_EPS_DETACH_RETRANSMISSION_COUNTER_MAX) {
    /*
     * Resend SGS Implicit EPS Detach Indication message to the SGS
     */
    mme_app_send_sgs_eps_detach_indication(
        ue_context_p, ue_context_p->sgs_detach_type);
  } else {
    OAILOG_DEBUG(
        LOG_MME_APP,
        "SGS Implicit EPS DETACH indication failed after %u retransmission and "
        "expiry \n",
        ue_context_p->sgs_context->ts13_retransmission_count);
    increment_counter(
        "sgs_eps_implicit_detach_timer_expired", 1, 1, "cause",
        "Ts13 timer expired after max tetransmission");
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}
//------------------------------------------------------------------------------
void mme_app_send_sgs_imsi_detach_indication(
    struct ue_mm_context_s* ue_context_p, uint8_t detach_type) {
  MessageDef* message_p = NULL;

  OAILOG_FUNC_IN(LOG_MME_APP);
  OAILOG_INFO(
      LOG_MME_APP,
      "Send SGSAP_IMSI_DETACH_IND to SGS, detach_type = %u for "
      "ue_id " MME_UE_S1AP_ID_FMT "\n",
      detach_type, ue_context_p->mme_ue_s1ap_id);
  message_p = itti_alloc_new_message(TASK_MME_APP, SGSAP_IMSI_DETACH_IND);
  if (!message_p) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Failed to allocate memory for SGSAP_IMSI_DETACH_IND for "
        "ue_id: " MME_UE_S1AP_ID_FMT "\n",
        ue_context_p->mme_ue_s1ap_id);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  memset(
      (void*) &message_p->ittiMsg.sgsap_imsi_detach_ind, 0,
      sizeof(itti_sgsap_imsi_detach_ind_t));
  IMSI64_TO_STRING(
      ue_context_p->emm_context._imsi64, SGSAP_IMSI_DETACH_IND(message_p).imsi,
      ue_context_p->emm_context._imsi.length);
  SGSAP_IMSI_DETACH_IND(message_p).imsi_length =
      (uint8_t) strlen(SGSAP_IMSI_DETACH_IND(message_p).imsi);
  SGSAP_IMSI_DETACH_IND(message_p).noneps_detach_type = detach_type;

  send_msg_to_task(&mme_app_task_zmq_ctx, TASK_SGS, message_p);
  if (detach_type == SGS_IMPLICIT_NW_INITIATED_IMSI_DETACH_FROM_EPS_N_NONEPS) {
    // Start SGS Implicit IMSI Detach indication timer
    if ((ue_context_p->sgs_context->ts10_timer.id = mme_app_start_timer(
             ue_context_p->sgs_context->ts10_timer.sec * 1000,
             TIMER_REPEAT_ONCE,
             mme_app_handle_sgs_implicit_imsi_detach_timer_expiry,
             ue_context_p->mme_ue_s1ap_id)) == -1) {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Failed to start SGS Implicit IMSI Detach indication timer for UE id "
          " "
          "%d \n",
          ue_context_p->mme_ue_s1ap_id);
      ue_context_p->sgs_context->ts10_timer.id = MME_APP_TIMER_INACTIVE_ID;
    } else {
      OAILOG_DEBUG(
          LOG_MME_APP,
          "Started SGS Implicit IMSI Detach indication timer for UE id  %d \n",
          ue_context_p->mme_ue_s1ap_id);
    }
  } else {
    // Start SGS IMSI Detach indication timer

    if ((ue_context_p->sgs_context->ts9_timer.id = mme_app_start_timer(
             ue_context_p->sgs_context->ts9_timer.sec * 1000, TIMER_REPEAT_ONCE,
             mme_app_handle_sgs_imsi_detach_timer_expiry,
             ue_context_p->mme_ue_s1ap_id)) == -1) {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Failed to start SGS IMSI Detach indication timer for UE id  %d \n",
          ue_context_p->mme_ue_s1ap_id);
      ue_context_p->sgs_context->ts9_timer.id = MME_APP_TIMER_INACTIVE_ID;
    } else {
      OAILOG_DEBUG(
          LOG_MME_APP,
          "Started SGS IMSI Detach indication timer for UE id  %d \n",
          ue_context_p->mme_ue_s1ap_id);
    }

    // Stop and reset SGS Location Update Request timer if running
    if (ue_context_p->sgs_context->ts6_1_timer.id !=
        MME_APP_TIMER_INACTIVE_ID) {
      mme_app_stop_timer(ue_context_p->sgs_context->ts6_1_timer.id);
      ue_context_p->sgs_context->ts6_1_timer.id = MME_APP_TIMER_INACTIVE_ID;
    }
  }
  increment_counter(
      "mme_sgs_imsi_detach_indication_sent", 1, 1, "result", "success");

  OAILOG_FUNC_OUT(LOG_MME_APP);
}

/* handle the SGS IMSI detach timer expiry. */
int mme_app_handle_sgs_imsi_detach_timer_expiry(
    zloop_t* loop, int timer_id, void* args) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  mme_ue_s1ap_id_t mme_ue_s1ap_id = 0;
  if (!mme_app_get_timer_arg(timer_id, &mme_ue_s1ap_id)) {
    OAILOG_WARNING(
        LOG_MME_APP, "Invalid Timer Id expiration, Timer Id: %u\n", timer_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  struct ue_mm_context_s* ue_context_p =
      mme_app_get_ue_context_for_timer(mme_ue_s1ap_id, "sgs imsi detach timer");
  if (ue_context_p == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Invalid UE context received, MME UE S1AP Id: " MME_UE_S1AP_ID_FMT "\n",
        mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  if (ue_context_p->sgs_context == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "SGS EPS Detach Timer expired but no associated SGS context for UE "
        "id " MME_UE_S1AP_ID_FMT "\n",
        mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  /*
   * Increment the retransmission counter
   */
  ue_context_p->sgs_context->ts9_retransmission_count += 1;
  OAILOG_WARNING(
      LOG_MME_APP,
      "MME APP: Ts9 timer expired,retransmission "
      "counter = %u \n",
      ue_context_p->sgs_context->ts9_retransmission_count);

  ue_context_p->sgs_context->ts9_timer.id = MME_APP_TIMER_INACTIVE_ID;
  if (ue_context_p->sgs_context->ts9_retransmission_count <
      IMSI_DETACH_RETRANSMISSION_COUNTER_MAX) {
    /* Send the Detach Accept to Ue after the Ts9 timer expired and maximum
     * retransmission */
    send_msg_to_task(
        &mme_app_task_zmq_ctx, TASK_S1AP, ue_context_p->sgs_context->message_p);
    ue_context_p->sgs_context->message_p = NULL;
    /*
     Notify S1AP to send UE Context Release Command to eNB or free s1 context
     locally, if the ue requested for combined EPS/IMSI detach if the ue is in
     idle state and requested for IMSI detach
    */
    if ((ue_context_p->sgs_detach_type ==
         SGS_COMBINED_UE_INITIATED_IMSI_DETACH_FROM_EPS_N_NONEPS) ||
        ((ue_context_p->sgs_detach_type ==
          SGS_EXPLICIT_UE_INITIATED_IMSI_DETACH_FROM_NONEPS) &&
         (ue_context_p->ue_context_rel_cause ==
          S1AP_RADIO_EUTRAN_GENERATED_REASON))) {
      mme_app_itti_ue_context_release(
          ue_context_p, ue_context_p->ue_context_rel_cause);
      ue_context_p->ue_context_rel_cause = S1AP_INVALID_CAUSE;
    }
    /*
     * Resend SGS IMSI Detach Indication message to the SGS
     */
    mme_app_send_sgs_imsi_detach_indication(
        ue_context_p, ue_context_p->sgs_detach_type);
  } else {
    OAILOG_DEBUG(
        LOG_MME_APP,
        "SGS IMSI DETACH indication failed after %u retransmission and expiry "
        "\n",
        ue_context_p->sgs_context->ts9_retransmission_count);
    // Free the UE SGS context
    mme_app_ue_sgs_context_free_content(
        ue_context_p->sgs_context, ue_context_p->emm_context._imsi64);
    free_wrapper((void**) &(ue_context_p->sgs_context));
    increment_counter(
        "sgs_imsi_detach_timer_expired", 1, 1, "cause",
        "Ts9 timer expired after max retransmission");
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

/* handle the SGS Implicit IMSI detach timer expiry. */
int mme_app_handle_sgs_implicit_imsi_detach_timer_expiry(
    zloop_t* loop, int timer_id, void* args) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  mme_ue_s1ap_id_t mme_ue_s1ap_id = 0;
  if (!mme_app_get_timer_arg(timer_id, &mme_ue_s1ap_id)) {
    OAILOG_WARNING(
        LOG_MME_APP, "Invalid Timer Id expiration, Timer Id: %u\n", timer_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  struct ue_mm_context_s* ue_context_p = mme_app_get_ue_context_for_timer(
      mme_ue_s1ap_id, "sgs implicit imsi detach timer");
  if (ue_context_p == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Invalid UE context received, MME UE S1AP Id: " MME_UE_S1AP_ID_FMT "\n",
        mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  if (ue_context_p->sgs_context == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "SGS IMPLICIT IMSI Detach Timer expired but no assoicated SGS context "
        "for"
        " ue_id " MME_UE_S1AP_ID_FMT "\n",
        mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  /*
   * Increment the retransmission counter
   */
  ue_context_p->sgs_context->ts10_retransmission_count += 1;
  OAILOG_WARNING(
      LOG_NAS_EMM,
      "MME APP: Ts10 timer expired,retransmission "
      "counter = %u \n",
      ue_context_p->sgs_context->ts10_retransmission_count);

  ue_context_p->sgs_context->ts10_timer.id = MME_APP_TIMER_INACTIVE_ID;
  if (ue_context_p->sgs_context->ts10_retransmission_count <
      IMPLICIT_IMSI_DETACH_RETRANSMISSION_COUNTER_MAX) {
    /*
     * Resend SGS Implicit IMSI Detach Indication message to the SGS
     */
    mme_app_send_sgs_imsi_detach_indication(
        ue_context_p, ue_context_p->sgs_detach_type);
  } else {
    OAILOG_DEBUG(
        LOG_MME_APP,
        "SGS Implicit IMSI DETACH indication failed after %u retransmission "
        "and "
        "expiry \n",
        ue_context_p->sgs_context->ts10_retransmission_count);
    increment_counter(
        "sgs_imsi_implicit_detach_timer_expired", 1, 1, "cause",
        "Ts10 timer expired after max tetransmission");
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

//------------------------------------------------------------------------------
void mme_app_handle_sgs_detach_req(
    ue_mm_context_t* ue_context_p, emm_proc_sgs_detach_type_t detach_type) {
  sgs_fsm_t evnt = {0};

  OAILOG_FUNC_IN(LOG_MME_APP);
  if (ue_context_p == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP, "UE context doesn't exist -> Nothing to do :-) \n");
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  if (ue_context_p->sgs_context) {
    evnt.ue_id = ue_context_p->mme_ue_s1ap_id;
    evnt.ctx   = ue_context_p->sgs_context;
    // check the SGS state and if it is null then do not send Detach towards SGS
    OAILOG_DEBUG(LOG_MME_APP, "SGS Detach type = ( %u )\n", detach_type);
    if (sgs_fsm_get_status(evnt.ue_id, evnt.ctx) != SGS_NULL) {
      switch (detach_type) {
        // Handle Ue initiated EPS detach towards SGS
        case EMM_SGS_UE_INITIATED_EPS_DETACH: {
          ue_context_p->sgs_detach_type = SGS_UE_INITIATED_IMSI_DETACH_FROM_EPS;
          mme_app_send_sgs_eps_detach_indication(
              ue_context_p, ue_context_p->sgs_detach_type);
          evnt.primitive = _SGS_EPS_DETACH_IND;
        } break;
        // Handle Ue initiated IMSI detach towards SGS
        case EMM_SGS_UE_INITIATED_EXPLICIT_NONEPS_DETACH: {
          ue_context_p->sgs_detach_type =
              SGS_EXPLICIT_UE_INITIATED_IMSI_DETACH_FROM_NONEPS;
          mme_app_send_sgs_imsi_detach_indication(
              ue_context_p, ue_context_p->sgs_detach_type);
          evnt.primitive = _SGS_IMSI_DETACH_IND;
        } break;
        // Handle Ue initiated Combined EPS/IMSI detach towards SGS
        case EMM_SGS_UE_INITIATED_COMBINED_DETACH: {
          ue_context_p->sgs_detach_type =
              SGS_COMBINED_UE_INITIATED_IMSI_DETACH_FROM_EPS_N_NONEPS;
          mme_app_send_sgs_imsi_detach_indication(
              ue_context_p, ue_context_p->sgs_detach_type);
          evnt.primitive = _SGS_IMSI_DETACH_IND;
        } break;
        // Handle Network initiated EPS detach towards SGS
        case EMM_SGS_NW_INITIATED_EPS_DETACH: {
          ue_context_p->sgs_detach_type = SGS_NW_INITIATED_IMSI_DETACH_FROM_EPS;
          mme_app_send_sgs_eps_detach_indication(
              ue_context_p, ue_context_p->sgs_detach_type);
          evnt.primitive = _SGS_EPS_DETACH_IND;
        } break;
        // Handle Network initiated Implicit IMSI detach towards SGS
        case EMM_SGS_NW_INITIATED_IMPLICIT_NONEPS_DETACH: {
          ue_context_p->sgs_detach_type =
              SGS_IMPLICIT_NW_INITIATED_IMSI_DETACH_FROM_EPS_N_NONEPS;
          mme_app_send_sgs_imsi_detach_indication(
              ue_context_p, ue_context_p->sgs_detach_type);
          evnt.primitive = _SGS_IMSI_DETACH_IND;
        } break;
        default:
          OAILOG_INFO(
              LOG_MME_APP,
              "SGS-DETACH REQ: Ue Id %u Invalid detach type : %u \n",
              ue_context_p->mme_ue_s1ap_id, ue_context_p->sgs_detach_type);
          break;
      }
      /*
       * Call the SGS FSM to process the respective message
       * in different state and update the SGS State based on event
       */
      sgs_fsm_process(&evnt);
    }
  } else {
    OAILOG_ERROR(
        LOG_MME_APP,
        "UE SGS context doesn't exist for ue-id" MME_UE_S1AP_ID_FMT
        "-> Nothing to do :-) \n",
        ue_context_p->mme_ue_s1ap_id);
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

int mme_app_handle_sgs_eps_detach_ack(
    mme_app_desc_t* mme_app_desc_p,
    const itti_sgsap_eps_detach_ack_t* const eps_detach_ack_p) {
  imsi64_t imsi64                      = INVALID_IMSI64;
  struct ue_mm_context_s* ue_context_p = NULL;

  OAILOG_FUNC_IN(LOG_MME_APP);
  if (eps_detach_ack_p == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Invalid EPS Detach Acknowledgement ITTI message received\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  IMSI_STRING_TO_IMSI64(eps_detach_ack_p->imsi, &imsi64);
  OAILOG_INFO(
      LOG_MME_APP, "Received SGS EPS DETACH ACK for imsi " IMSI_64_FMT "\n",
      imsi64);

  if ((ue_context_p = mme_ue_context_exists_imsi(
           &mme_app_desc_p->mme_ue_contexts, imsi64)) == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP, "SGS-EPS DETACH ACK: Unknown IMSI " IMSI_64_FMT "\n",
        imsi64);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  if (ue_context_p->sgs_context) {
    // Stop SGS EPS Detach timer, after recieving the SGS EPS Detach Ack,if
    // running
    if (ue_context_p->sgs_context->ts8_timer.id != MME_APP_TIMER_INACTIVE_ID) {
      mme_app_stop_timer(ue_context_p->sgs_context->ts8_timer.id);
      ue_context_p->sgs_context->ts8_timer.id = MME_APP_TIMER_INACTIVE_ID;
      OAILOG_INFO(LOG_MME_APP, "Stopped Ts8 timer \n");
    } else if (
        ue_context_p->sgs_context->ts13_timer.id != MME_APP_TIMER_INACTIVE_ID) {
      mme_app_stop_timer(ue_context_p->sgs_context->ts13_timer.id);
      ue_context_p->sgs_context->ts13_timer.id = MME_APP_TIMER_INACTIVE_ID;
      OAILOG_INFO(LOG_MME_APP, "Stopped Ts13 timer \n");
    }
    // Release SGS Context
    if (ue_context_p->sgs_context != NULL) {
      // free the sgs context
      mme_app_ue_sgs_context_free_content(ue_context_p->sgs_context, imsi64);
      free_wrapper((void**) &(ue_context_p->sgs_context));
    }
  } else {
    OAILOG_ERROR(
        LOG_MME_APP,
        "SGS context not found in mme_app_handle_sgs_eps_detach_ack\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

int mme_app_handle_sgs_imsi_detach_ack(
    mme_app_desc_t* mme_app_desc_p,
    const itti_sgsap_imsi_detach_ack_t* const imsi_detach_ack_p) {
  imsi64_t imsi64                      = INVALID_IMSI64;
  struct ue_mm_context_s* ue_context_p = NULL;
  int rc                               = RETURNok;

  OAILOG_FUNC_IN(LOG_MME_APP);
  if (imsi_detach_ack_p == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Invalid IMSI Detach Acknowledgement ITTI message received\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  IMSI_STRING_TO_IMSI64(imsi_detach_ack_p->imsi, &imsi64);
  OAILOG_DEBUG(
      LOG_MME_APP, "Received SGS IMSI DETACH ACK for imsi " IMSI_64_FMT "\n",
      imsi64);

  if ((ue_context_p = mme_ue_context_exists_imsi(
           &mme_app_desc_p->mme_ue_contexts, imsi64)) == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP, "SGS-IMSI DETACH ACK: Unknown IMSI " IMSI_64_FMT "\n",
        imsi64);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  if (ue_context_p->sgs_context) {
    // Stop SGS IMSI Detach timer, after recieving the SGS EPS Detach Ack,if
    // running
    if (ue_context_p->sgs_context->ts9_timer.id != MME_APP_TIMER_INACTIVE_ID) {
      mme_app_stop_timer(ue_context_p->sgs_context->ts9_timer.id);
      ue_context_p->sgs_context->ts9_timer.id = MME_APP_TIMER_INACTIVE_ID;
    }
    /*
     * Send the S1AP NAS DL DATA REQ in case of IMSI or combined EPS/IMSI detach
     * once the SGS IMSI Detach Ack recieved from SGS task.
     */
    if ((ue_context_p->sgs_detach_type ==
         SGS_EXPLICIT_UE_INITIATED_IMSI_DETACH_FROM_NONEPS) ||
        (ue_context_p->sgs_detach_type ==
         SGS_COMBINED_UE_INITIATED_IMSI_DETACH_FROM_EPS_N_NONEPS)) {
      if (!ue_context_p->sgs_context->message_p) {
        OAILOG_DEBUG(
            LOG_MME_APP,
            "Detach Accept has been sent already after ts9 timer expired for "
            "UE id " MME_UE_S1AP_ID_FMT ", ignore the IMSI detach Ack \n",
            ue_context_p->mme_ue_s1ap_id);
        OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
      } else {
        rc = send_msg_to_task(
            &mme_app_task_zmq_ctx, TASK_S1AP,
            ue_context_p->sgs_context->message_p);
        ue_context_p->sgs_context->message_p = NULL;
      }
      /*
       Notify S1AP to send UE Context Release Command to eNB or free s1 context
       locally, if the ue requested for combined EPS/IMSI detach if the ue is in
       idle state and requested for IMSI detach
      */
      if ((ue_context_p->sgs_detach_type ==
           SGS_EXPLICIT_UE_INITIATED_IMSI_DETACH_FROM_NONEPS) &&
          (ue_context_p->ue_context_rel_cause ==
           S1AP_RADIO_EUTRAN_GENERATED_REASON)) {
        mme_app_itti_ue_context_release(
            ue_context_p, ue_context_p->ue_context_rel_cause);
        ue_context_p->ue_context_rel_cause = S1AP_INVALID_CAUSE;
      } else if (
          ue_context_p->sgs_detach_type ==
          SGS_COMBINED_UE_INITIATED_IMSI_DETACH_FROM_EPS_N_NONEPS) {
        mme_app_itti_ue_context_release(
            ue_context_p, ue_context_p->ue_context_rel_cause);
        ue_context_p->ue_context_rel_cause = S1AP_INVALID_CAUSE;
      }
    }
    // Free the UE SGS context
    mme_app_ue_sgs_context_free_content(
        ue_context_p->sgs_context, ue_context_p->emm_context._imsi64);
    free_wrapper((void**) &(ue_context_p->sgs_context));
  } else {
    OAILOG_ERROR(
        LOG_MME_APP,
        "SGS context not found in mme_app_handle_sgs_imsi_detach_ack\n");
    rc = RETURNerror;
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}
