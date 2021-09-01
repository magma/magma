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

#include <string.h>
#include <stdlib.h>

#include "common_types.h"
#include "3gpp_24.008.h"
#include "emm_data.h"
#include "nas_timer.h"
#include "nas_procedures.h"
#include "esm_data.h"
#include "esm_proc.h"
#include "log.h"
#include "dynamic_memory_check.h"
#include "common_defs.h"
#include "mme_app_ue_context.h"
#include "mme_config.h"
#include "3gpp_36.401.h"
#include "mme_app_defs.h"
#include "mme_app_timer.h"

// free allocated structs
//------------------------------------------------------------------------------
void free_esm_bearer_context(esm_ebr_context_t* esm_ebr_context) {
  if (esm_ebr_context) {
    if (esm_ebr_context->pco) {
      free_protocol_configuration_options(&esm_ebr_context->pco);
    }
    if (esm_ebr_context->tft) {
      free_traffic_flow_template(&esm_ebr_context->tft);
    }
    if (NAS_TIMER_INACTIVE_ID != esm_ebr_context->timer.id) {
      esm_ebr_timer_data_t* esm_ebr_timer_data = NULL;
      esm_ebr_context->timer.id                = nas_timer_stop(
          esm_ebr_context->timer.id, (void**) &esm_ebr_timer_data);
      /*
       * Release the retransmisison timer parameters
       */
      if (esm_ebr_timer_data) {
        if (esm_ebr_timer_data->msg) {
          bdestroy_wrapper(&esm_ebr_timer_data->msg);
        }
        free_wrapper((void**) &esm_ebr_timer_data);
      }
    }
  }
}

//------------------------------------------------------------------------------
void esm_bearer_context_init(esm_ebr_context_t* esm_ebr_context) {
  if (esm_ebr_context) {
    memset(esm_ebr_context, 0, sizeof(*esm_ebr_context));
    esm_ebr_context->status   = ESM_EBR_INACTIVE;
    esm_ebr_context->timer.id = NAS_TIMER_INACTIVE_ID;
  }
}

//------------------------------------------------------------------------------
void nas_start_T3489(
    const mme_ue_s1ap_id_t ue_id, struct nas_timer_s* const T3489,
    time_out_t time_out_cb) {
  if ((T3489) && (T3489->id == NAS_TIMER_INACTIVE_ID)) {
    T3489->id = mme_app_start_timer(
        T3489->sec * 1000, TIMER_REPEAT_ONCE, time_out_cb, ue_id);
    if (NAS_TIMER_INACTIVE_ID != T3489->id) {
      OAILOG_DEBUG(
          LOG_NAS_EMM, "T3489 started UE " MME_UE_S1AP_ID_FMT "\n", ue_id);
    } else {
      OAILOG_ERROR(
          LOG_NAS_EMM, "Could not start T3489 UE " MME_UE_S1AP_ID_FMT " ",
          ue_id);
    }
  }
}

//------------------------------------------------------------------------------
void nas_stop_T3489(esm_context_t* const esm_ctx) {
  if ((esm_ctx) && (esm_ctx->T3489.id != NAS_TIMER_INACTIVE_ID)) {
    emm_context_t* emm_context =
        PARENT_STRUCT(esm_ctx, struct emm_context_s, esm_ctx);
    ue_mm_context_t* ue_mm_context =
        PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context);
    mme_ue_s1ap_id_t ue_id = ue_mm_context->mme_ue_s1ap_id;
    mme_app_stop_timer(esm_ctx->T3489.id);
    esm_ctx->T3489.id = NAS_TIMER_INACTIVE_ID;

    OAILOG_DEBUG(
        LOG_NAS_EMM, "T3489 stopped UE " MME_UE_S1AP_ID_FMT "\n", ue_id);
    if (esm_ctx->t3489_arg) {
      esm_ebr_timer_data_t* data = (esm_ebr_timer_data_t*) esm_ctx->t3489_arg;
      data->ctx                  = NULL;
      bdestroy_wrapper(&data->msg);
      free_wrapper((void**) &data);
    }
  }
}

// free allocated structs
//------------------------------------------------------------------------------
void free_esm_context_content(esm_context_t* esm_ctx) {
  if (!esm_ctx) {
    return;
  }
  nas_stop_T3489(esm_ctx);
  if (esm_ctx->esm_proc_data) {
    OAILOG_DEBUG(LOG_NAS_ESM, "Free up esm_proc_data");
    bdestroy_wrapper(&esm_ctx->esm_proc_data->apn);
    if (esm_ctx->esm_proc_data->pco.num_protocol_or_container_id) {
      clear_protocol_configuration_options(&esm_ctx->esm_proc_data->pco);
    }
    free_wrapper((void**) &esm_ctx->esm_proc_data);
  }
}

//------------------------------------------------------------------------------
void esm_init_context(struct esm_context_s* esm_context) {
  emm_context_t* emm_context =
      PARENT_STRUCT(esm_context, struct emm_context_s, esm_ctx);
  ue_mm_context_t* ue_mm_context =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context);
  OAILOG_DEBUG(
      LOG_NAS_ESM, "ESM-CTX - Init UE id " MME_UE_S1AP_ID_FMT "\n",
      ue_mm_context->mme_ue_s1ap_id);
  memset(esm_context, 0, sizeof(*esm_context));
  esm_context->T3489.id  = NAS_TIMER_INACTIVE_ID;
  esm_context->T3489.sec = mme_config.nas_config.t3489_sec;
}
