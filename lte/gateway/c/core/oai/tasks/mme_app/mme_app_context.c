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

/*! \file mme_app_context.c
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <string.h>
#include <inttypes.h>
#include <arpa/inet.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <stdint.h>
#include <sys/time.h>
#include <pthread.h>
#include <execinfo.h>
#include <netinet/in.h>
#include <sys/socket.h>
#include <sys/time.h>
#include <time.h>

#include "dynamic_memory_check.h"
#include "log.h"
#include "common_types.h"
#include "conversions.h"
#include "intertask_interface.h"
#include "mme_config.h"
#include "enum_string.h"
#include "mme_app_ue_context.h"
#include "mme_app_bearer_context.h"
#include "mme_app_pdn_context.h"
#include "mme_app_defs.h"
#include "mme_app_itti_messaging.h"
#include "mme_app_procedures.h"
#include "nas_proc.h"
#include "common_defs.h"
#include "esm_ebr.h"
#include "esm_ebr_context.h"
#include "timer.h"
#include "mme_app_statistics.h"
#include "directoryd.h"
#include "3gpp_23.003.h"
#include "3gpp_24.008.h"
#include "3gpp_29.274.h"
#include "3gpp_36.401.h"
#include "bstrlib.h"
#include "emm_data.h"
#include "esm_data.h"
#include "hashtable.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "mme_api.h"
#include "mme_app_state.h"
#include "mme_app_timer.h"
#include "nas_timer.h"
#include "obj_hashtable.h"
#include "s1ap_messages_types.h"

extern task_zmq_ctx_t mme_app_task_zmq_ctx;

/* Obtain a backtrace and print it to stdout. */
void print_trace(void) {
  void* array[10];
  size_t size;
  char** strings;
  size_t i;

  size    = backtrace(array, 10);
  strings = backtrace_symbols(array, size);

  printf("Obtained %zd stack frames.\n", size);

  for (i = 0; i < size; i++) printf("%s\n", strings[i]);

  free(strings);
}

typedef void (*mme_app_timer_callback_t)(void* args, imsi64_t* imsi64);
static void mme_app_handle_s1ap_ue_context_release(
    const mme_ue_s1ap_id_t mme_ue_s1ap_id,
    const enb_ue_s1ap_id_t enb_ue_s1ap_id, uint32_t enb_id, enum s1cause cause);

static bool mme_app_recover_timers_for_ue(
    const hash_key_t keyP, void* const ue_context_pP, void* unused_param_pP,
    void** unused_result_pP);

static void directoryd_report_location_a(uint64_t imsi, uint8_t imsi_len) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi, imsi_str, imsi_len);
  directoryd_report_location(imsi_str);
  OAILOG_INFO_UE(LOG_MME_APP, imsi, "Reported UE location to directoryd\n");
}

static void directoryd_remove_location_a(uint64_t imsi, uint8_t imsi_len) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi, imsi_str, imsi_len);
  directoryd_remove_location(imsi_str);
  OAILOG_INFO_UE(LOG_MME_APP, imsi, "Deleted UE location from directoryd\n");
}

static void mme_app_resume_esm_ebr_timer(ue_mm_context_t* ue_context_p);

static void mme_app_handle_timer_for_unregistered_ue(
    ue_mm_context_t* ue_context_p);

//------------------------------------------------------------------------------
// warning: lock the UE context
ue_mm_context_t* mme_create_new_ue_context(void) {
  ue_mm_context_t* new_p = calloc(1, sizeof(ue_mm_context_t));
  if (!new_p) {
    OAILOG_ERROR(LOG_MME_APP, "Failed to allocate memory for UE context \n");
    return NULL;
  }

  new_p->mme_ue_s1ap_id  = INVALID_MME_UE_S1AP_ID;
  new_p->enb_s1ap_id_key = INVALID_ENB_UE_S1AP_ID_KEY;
  emm_init_context(&new_p->emm_context, true);

  // Initialize timers to INVALID IDs
  new_p->mobile_reachability_timer.id = MME_APP_TIMER_INACTIVE_ID;
  new_p->implicit_detach_timer.id     = MME_APP_TIMER_INACTIVE_ID;

  new_p->initial_context_setup_rsp_timer = (struct mme_app_timer_t){
      MME_APP_TIMER_INACTIVE_ID, MME_APP_INITIAL_CONTEXT_SETUP_RSP_TIMER_VALUE};
  new_p->paging_response_timer = (struct mme_app_timer_t){
      MME_APP_TIMER_INACTIVE_ID, MME_APP_PAGING_RESPONSE_TIMER_VALUE};
  new_p->ulr_response_timer = (struct mme_app_timer_t){
      MME_APP_TIMER_INACTIVE_ID, MME_APP_ULR_RESPONSE_TIMER_VALUE};
  new_p->ue_context_modification_timer = (struct mme_app_timer_t){
      MME_APP_TIMER_INACTIVE_ID, MME_APP_UE_CONTEXT_MODIFICATION_TIMER_VALUE};

  new_p->ue_context_rel_cause = S1AP_INVALID_CAUSE;
  new_p->sgs_context          = NULL;
  new_p->paging_retx_count    = 0;
  return new_p;
}

//------------------------------------------------------------------------------
void mme_app_ue_sgs_context_free_content(
    sgs_context_t* const sgs_context_p, imsi64_t imsi) {
  if (sgs_context_p == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP, "Invalid SGS context received for IMSI: " IMSI_64_FMT "\n",
        imsi);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  // Stop SGS Location update timer if running
  if (sgs_context_p->ts6_1_timer.id != MME_APP_TIMER_INACTIVE_ID) {
    mme_app_stop_timer(sgs_context_p->ts6_1_timer.id);
    sgs_context_p->ts6_1_timer.id = MME_APP_TIMER_INACTIVE_ID;
  }
  // Stop SGS EPS Detach indication timer if running
  if (sgs_context_p->ts8_timer.id != MME_APP_TIMER_INACTIVE_ID) {
    mme_app_stop_timer(sgs_context_p->ts8_timer.id);
    sgs_context_p->ts8_timer.id = MME_APP_TIMER_INACTIVE_ID;
  }

  // Stop SGS IMSI Detach indication timer if running
  if (sgs_context_p->ts9_timer.id != MME_APP_TIMER_INACTIVE_ID) {
    mme_app_stop_timer(sgs_context_p->ts9_timer.id);
    sgs_context_p->ts9_timer.id = MME_APP_TIMER_INACTIVE_ID;
  }
  // Stop SGS Implicit IMSI Detach indication timer if running
  if (sgs_context_p->ts10_timer.id != MME_APP_TIMER_INACTIVE_ID) {
    mme_app_stop_timer(sgs_context_p->ts10_timer.id);
    sgs_context_p->ts10_timer.id = MME_APP_TIMER_INACTIVE_ID;
  }
  // Stop SGS Implicit EPS Detach indication timer if running
  if (sgs_context_p->ts13_timer.id != MME_APP_TIMER_INACTIVE_ID) {
    mme_app_stop_timer(sgs_context_p->ts13_timer.id);
    sgs_context_p->ts13_timer.id = MME_APP_TIMER_INACTIVE_ID;
  }
}

//------------------------------------------------------------------------------
void mme_app_ue_context_free_content(ue_mm_context_t* const ue_context_p) {
  bdestroy_wrapper(&ue_context_p->msisdn);
  bdestroy_wrapper(&ue_context_p->ue_radio_capability);
  bdestroy_wrapper(&ue_context_p->apn_oi_replacement);

  // Stop Mobile reachability timer,if running
  if (ue_context_p->mobile_reachability_timer.id != MME_APP_TIMER_INACTIVE_ID) {
    mme_app_stop_timer(ue_context_p->mobile_reachability_timer.id);
    ue_context_p->mobile_reachability_timer.id = MME_APP_TIMER_INACTIVE_ID;
  }
  // Stop Implicit detach timer,if running
  if (ue_context_p->implicit_detach_timer.id != MME_APP_TIMER_INACTIVE_ID) {
    mme_app_stop_timer(ue_context_p->implicit_detach_timer.id);
    ue_context_p->implicit_detach_timer.id = MME_APP_TIMER_INACTIVE_ID;
  }

  // Stop Initial context setup process guard timer,if running
  if (ue_context_p->initial_context_setup_rsp_timer.id !=
      MME_APP_TIMER_INACTIVE_ID) {
    mme_app_stop_timer(ue_context_p->initial_context_setup_rsp_timer.id);
    ue_context_p->time_ics_rsp_timer_started = 0;
    ue_context_p->initial_context_setup_rsp_timer.id =
        MME_APP_TIMER_INACTIVE_ID;
  }
  // Stop UE context modification process guard timer,if running
  if (ue_context_p->ue_context_modification_timer.id !=
      MME_APP_TIMER_INACTIVE_ID) {
    mme_app_stop_timer(ue_context_p->ue_context_modification_timer.id);
    ue_context_p->ue_context_modification_timer.id = MME_APP_TIMER_INACTIVE_ID;
  }

  // Stop ULR Response timer if running
  if (ue_context_p->ulr_response_timer.id != MME_APP_TIMER_INACTIVE_ID) {
    mme_app_stop_timer(ue_context_p->ulr_response_timer.id);
    ue_context_p->ulr_response_timer.id = MME_APP_TIMER_INACTIVE_ID;
  }

  // Stop paging response timer if running
  if (ue_context_p->paging_response_timer.id != MME_APP_TIMER_INACTIVE_ID) {
    mme_app_stop_timer(ue_context_p->paging_response_timer.id);
    ue_context_p->paging_response_timer.id = MME_APP_TIMER_INACTIVE_ID;
  }

  if (ue_context_p->sgs_context != NULL) {
    // free the sgs context
    mme_app_ue_sgs_context_free_content(
        ue_context_p->sgs_context, ue_context_p->emm_context._imsi64);
    free_wrapper((void**) &(ue_context_p->sgs_context));
  }
  ue_context_p->ue_context_rel_cause = S1AP_INVALID_CAUSE;

  ue_context_p->send_ue_purge_request = false;
  ue_context_p->hss_initiated_detach  = false;
  for (int i = 0; i < MAX_APN_PER_UE; i++) {
    if (ue_context_p->pdn_contexts[i]) {
      mme_app_free_pdn_context(
          &ue_context_p->pdn_contexts[i], ue_context_p->emm_context._imsi64);
    }
  }

  for (int i = 0; i < BEARERS_PER_UE; i++) {
    if (ue_context_p->bearer_contexts[i]) {
      mme_app_free_bearer_context(&ue_context_p->bearer_contexts[i]);
    }
  }

  if (ue_context_p->s11_procedures) {
    mme_app_delete_s11_procedures(ue_context_p);
  }
}

void mme_app_state_free_ue_context(void** ue_context_node) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  ue_mm_context_t* ue_context_p = (ue_mm_context_t*) (*ue_context_node);
  // clean up EMM context
  emm_context_t* emm_ctx = &ue_context_p->emm_context;
  free_emm_ctx_memory(emm_ctx, ue_context_p->mme_ue_s1ap_id);
  mme_app_ue_context_free_content(ue_context_p);
  free_wrapper((void**) &ue_context_p);
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
ue_mm_context_t* mme_ue_context_exists_enb_ue_s1ap_id(
    mme_ue_context_t* const mme_ue_context_p, const enb_s1ap_id_key_t enb_key) {
  hashtable_rc_t h_rc       = HASH_TABLE_OK;
  uint64_t mme_ue_s1ap_id64 = 0;

  hashtable_uint64_ts_get(
      mme_ue_context_p->enb_ue_s1ap_id_ue_context_htbl,
      (const hash_key_t) enb_key, &mme_ue_s1ap_id64);
  if (HASH_TABLE_OK == h_rc) {
    return mme_ue_context_exists_mme_ue_s1ap_id(
        (mme_ue_s1ap_id_t) mme_ue_s1ap_id64);
  }
  return NULL;
}

//------------------------------------------------------------------------------
ue_mm_context_t* mme_ue_context_exists_mme_ue_s1ap_id(
    const mme_ue_s1ap_id_t mme_ue_s1ap_id) {
  struct ue_mm_context_s* ue_context_p = NULL;
  hash_table_ts_t* state_imsi_ht       = get_mme_ue_state();

  hashtable_ts_get(
      state_imsi_ht, (const hash_key_t) mme_ue_s1ap_id, (void**) &ue_context_p);
  if (ue_context_p) {
    OAILOG_TRACE(
        LOG_MME_APP,
        "UE  " MME_UE_S1AP_ID_FMT " fetched MM state %s, ECM state %s\n ",
        mme_ue_s1ap_id,
        (ue_context_p->mm_state == UE_UNREGISTERED) ?
            "UE_UNREGISTERED" :
            (ue_context_p->mm_state == UE_REGISTERED) ? "UE_REGISTERED" :
                                                        "UNKNOWN",
        (ue_context_p->ecm_state == ECM_IDLE) ?
            "ECM_IDLE" :
            (ue_context_p->ecm_state == ECM_CONNECTED) ? "ECM_CONNECTED" :
                                                         "UNKNOWN");
  }
  return ue_context_p;
}

//------------------------------------------------------------------------------
struct ue_mm_context_s* mme_ue_context_exists_imsi(
    mme_ue_context_t* const mme_ue_context_p, imsi64_t imsi) {
  hashtable_rc_t h_rc       = HASH_TABLE_OK;
  uint64_t mme_ue_s1ap_id64 = 0;

  h_rc = hashtable_uint64_ts_get(
      mme_ue_context_p->imsi_mme_ue_id_htbl, (const hash_key_t) imsi,
      &mme_ue_s1ap_id64);

  if (HASH_TABLE_OK == h_rc) {
    return mme_ue_context_exists_mme_ue_s1ap_id(
        (mme_ue_s1ap_id_t) mme_ue_s1ap_id64);
  } else {
    OAILOG_WARNING_UE(LOG_MME_APP, imsi, " No IMSI hashtable for this IMSI\n");
  }
  return NULL;
}

//------------------------------------------------------------------------------
struct ue_mm_context_s* mme_ue_context_exists_s11_teid(
    mme_ue_context_t* const mme_ue_context_p, const s11_teid_t teid) {
  hashtable_rc_t h_rc       = HASH_TABLE_OK;
  uint64_t mme_ue_s1ap_id64 = 0;

  h_rc = hashtable_uint64_ts_get(
      mme_ue_context_p->tun11_ue_context_htbl, (const hash_key_t) teid,
      &mme_ue_s1ap_id64);

  if (HASH_TABLE_OK == h_rc) {
    return mme_ue_context_exists_mme_ue_s1ap_id(
        (mme_ue_s1ap_id_t) mme_ue_s1ap_id64);
  } else {
    OAILOG_WARNING(
        LOG_MME_APP, " No S11 hashtable for S11 Teid " TEID_FMT "\n", teid);
  }
  return NULL;
}

//------------------------------------------------------------------------------
ue_mm_context_t* mme_ue_context_exists_guti(
    mme_ue_context_t* const mme_ue_context_p, const guti_t* const guti_p) {
  hashtable_rc_t h_rc       = HASH_TABLE_OK;
  uint64_t mme_ue_s1ap_id64 = 0;

  h_rc = obj_hashtable_uint64_ts_get(
      mme_ue_context_p->guti_ue_context_htbl, (const void*) guti_p,
      sizeof(*guti_p), &mme_ue_s1ap_id64);

  if (HASH_TABLE_OK == h_rc) {
    return mme_ue_context_exists_mme_ue_s1ap_id(
        (mme_ue_s1ap_id_t) mme_ue_s1ap_id64);
  } else {
    OAILOG_WARNING(LOG_MME_APP, " No GUTI hashtable for GUTI ");
  }

  return NULL;
}

//------------------------------------------------------------------------------
void mme_app_move_context(ue_mm_context_t* dst, ue_mm_context_t* src) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  if ((dst) && (src)) {
    enb_s1ap_id_key_t enb_s1ap_id_key = dst->enb_s1ap_id_key;
    enb_ue_s1ap_id_t enb_ue_s1ap_id   = dst->enb_ue_s1ap_id;
    mme_ue_s1ap_id_t mme_ue_s1ap_id   = dst->mme_ue_s1ap_id;
    memcpy(dst, src, sizeof(*dst));
    dst->enb_s1ap_id_key = enb_s1ap_id_key;
    dst->enb_ue_s1ap_id  = enb_ue_s1ap_id;
    dst->mme_ue_s1ap_id  = mme_ue_s1ap_id;
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
void mme_ue_context_update_coll_keys(
    mme_ue_context_t* const mme_ue_context_p,
    ue_mm_context_t* const ue_context_p,
    const enb_s1ap_id_key_t enb_s1ap_id_key,
    const mme_ue_s1ap_id_t mme_ue_s1ap_id, imsi64_t imsi,
    const s11_teid_t mme_teid_s11,
    const guti_t* const guti_p)  //  never NULL, if none put &ue_context_p->guti
{
  hashtable_rc_t h_rc                 = HASH_TABLE_OK;
  hash_table_ts_t* mme_state_ue_id_ht = get_mme_ue_state();
  OAILOG_FUNC_IN(LOG_MME_APP);

  OAILOG_TRACE(
      LOG_MME_APP,
      "Update ue context.old_enb_ue_s1ap_id_key %ld ue "
      "context.old_mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT
      " ue context.old_IMSI " IMSI_64_FMT " ue context.old_GUTI " GUTI_FMT "\n",
      ue_context_p->enb_s1ap_id_key, ue_context_p->mme_ue_s1ap_id,
      ue_context_p->emm_context._imsi64,
      GUTI_ARG(&ue_context_p->emm_context._guti));

  OAILOG_TRACE(
      LOG_MME_APP,
      "Update ue context %p updated_enb_ue_s1ap_id_key %ld "
      "updated_mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT " updated_IMSI " IMSI_64_FMT
      " updated_GUTI " GUTI_FMT "\n",
      ue_context_p, enb_s1ap_id_key, mme_ue_s1ap_id, imsi, GUTI_ARG(guti_p));

  if ((INVALID_ENB_UE_S1AP_ID_KEY != enb_s1ap_id_key) &&
      (ue_context_p->enb_s1ap_id_key != enb_s1ap_id_key)) {
    // new insertion of enb_ue_s1ap_id_key,
    h_rc = hashtable_uint64_ts_remove(
        mme_ue_context_p->enb_ue_s1ap_id_ue_context_htbl,
        (const hash_key_t) ue_context_p->enb_s1ap_id_key);
    h_rc = hashtable_uint64_ts_insert(
        mme_ue_context_p->enb_ue_s1ap_id_ue_context_htbl,
        (const hash_key_t) enb_s1ap_id_key, mme_ue_s1ap_id);

    if (HASH_TABLE_OK != h_rc) {
      OAILOG_ERROR_UE(
          LOG_MME_APP, imsi,
          "Error could not update this ue context %p "
          "enb_ue_s1ap_ue_id " ENB_UE_S1AP_ID_FMT
          " mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT " %s\n",
          ue_context_p, ue_context_p->enb_ue_s1ap_id,
          ue_context_p->mme_ue_s1ap_id, hashtable_rc_code2string(h_rc));
    }
    ue_context_p->enb_s1ap_id_key = enb_s1ap_id_key;
  } else {
    OAILOG_DEBUG_UE(
        LOG_MME_APP, imsi,
        "Did not update enb_s1ap_id_key %ld in ue context %p "
        "enb_ue_s1ap_ue_id " ENB_UE_S1AP_ID_FMT
        " mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT "\n",
        enb_s1ap_id_key, ue_context_p, ue_context_p->enb_ue_s1ap_id,
        ue_context_p->mme_ue_s1ap_id);
  }

  if (INVALID_MME_UE_S1AP_ID != mme_ue_s1ap_id) {
    if (ue_context_p->mme_ue_s1ap_id != mme_ue_s1ap_id) {
      // new insertion of mme_ue_s1ap_id, not a change in the id
      h_rc = hashtable_ts_remove(
          mme_state_ue_id_ht, (const hash_key_t) ue_context_p->mme_ue_s1ap_id,
          (void**) &ue_context_p);
      h_rc = hashtable_ts_insert(
          mme_state_ue_id_ht, (const hash_key_t) mme_ue_s1ap_id,
          (void*) ue_context_p);

      if (HASH_TABLE_OK != h_rc) {
        OAILOG_ERROR_UE(
            LOG_MME_APP, imsi,
            "Error could not update this ue context %p "
            "enb_ue_s1ap_ue_id " ENB_UE_S1AP_ID_FMT
            " mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT " %s\n",
            ue_context_p, ue_context_p->enb_ue_s1ap_id,
            ue_context_p->mme_ue_s1ap_id, hashtable_rc_code2string(h_rc));
      }
      ue_context_p->mme_ue_s1ap_id = mme_ue_s1ap_id;
    }
  } else {
    OAILOG_DEBUG_UE(
        LOG_MME_APP, imsi,
        "Did not update hashtable  for ue context %p "
        "enb_ue_s1ap_ue_id " ENB_UE_S1AP_ID_FMT
        " mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT " imsi " IMSI_64_FMT " \n",
        ue_context_p, ue_context_p->enb_ue_s1ap_id,
        ue_context_p->mme_ue_s1ap_id, imsi);
  }

  h_rc = hashtable_uint64_ts_remove(
      mme_ue_context_p->imsi_mme_ue_id_htbl,
      (const hash_key_t) ue_context_p->emm_context._imsi64);
  if (INVALID_MME_UE_S1AP_ID != mme_ue_s1ap_id) {
    h_rc = hashtable_uint64_ts_insert(
        mme_ue_context_p->imsi_mme_ue_id_htbl, (const hash_key_t) imsi,
        mme_ue_s1ap_id);
  } else {
    h_rc = HASH_TABLE_KEY_NOT_EXISTS;
  }
  if (HASH_TABLE_OK != h_rc) {
    OAILOG_ERROR_UE(
        LOG_MME_APP, imsi,
        "Error could not update this ue context %p "
        "enb_ue_s1ap_ue_id " ENB_UE_S1AP_ID_FMT
        " mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT " imsi " IMSI_64_FMT ": %s\n",
        ue_context_p, ue_context_p->enb_ue_s1ap_id,
        ue_context_p->mme_ue_s1ap_id, imsi, hashtable_rc_code2string(h_rc));
  }

  if ((INVALID_MME_UE_S1AP_ID != mme_ue_s1ap_id) && (mme_teid_s11)) {
    h_rc = hashtable_uint64_ts_remove(
        mme_ue_context_p->tun11_ue_context_htbl,
        (const hash_key_t) ue_context_p->mme_teid_s11);
    h_rc = hashtable_uint64_ts_insert(
        mme_ue_context_p->tun11_ue_context_htbl,
        (const hash_key_t) mme_teid_s11, (uint64_t) mme_ue_s1ap_id);
    ue_context_p->mme_teid_s11 = mme_teid_s11;
  } else {
    h_rc = HASH_TABLE_KEY_NOT_EXISTS;
  }

  if ((HASH_TABLE_OK != h_rc) && (mme_teid_s11)) {
    OAILOG_ERROR_UE(
        LOG_MME_APP, imsi,
        "Error could not update this ue context %p "
        "enb_ue_s1ap_ue_id " ENB_UE_S1AP_ID_FMT
        " mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT " mme_teid_s11 " TEID_FMT
        " : %s\n",
        ue_context_p, ue_context_p->enb_ue_s1ap_id,
        ue_context_p->mme_ue_s1ap_id, mme_teid_s11,
        hashtable_rc_code2string(h_rc));
  }

  if (guti_p) {
    if ((guti_p->gummei.mme_code !=
         ue_context_p->emm_context._guti.gummei.mme_code) ||
        (guti_p->gummei.mme_gid !=
         ue_context_p->emm_context._guti.gummei.mme_gid) ||
        (guti_p->m_tmsi != ue_context_p->emm_context._guti.m_tmsi) ||
        (guti_p->gummei.plmn.mcc_digit1 !=
         ue_context_p->emm_context._guti.gummei.plmn.mcc_digit1) ||
        (guti_p->gummei.plmn.mcc_digit2 !=
         ue_context_p->emm_context._guti.gummei.plmn.mcc_digit2) ||
        (guti_p->gummei.plmn.mcc_digit3 !=
         ue_context_p->emm_context._guti.gummei.plmn.mcc_digit3) ||
        (ue_context_p->mme_ue_s1ap_id != mme_ue_s1ap_id)) {
      // may check guti_p with a kind of instanceof()?
      h_rc = obj_hashtable_uint64_ts_remove(
          mme_ue_context_p->guti_ue_context_htbl,
          &ue_context_p->emm_context._guti, sizeof(*guti_p));
      if (INVALID_MME_UE_S1AP_ID != mme_ue_s1ap_id) {
        h_rc = obj_hashtable_uint64_ts_insert(
            mme_ue_context_p->guti_ue_context_htbl, (const void* const) guti_p,
            sizeof(*guti_p), (uint64_t) mme_ue_s1ap_id);
      } else {
        h_rc = HASH_TABLE_KEY_NOT_EXISTS;
      }

      if (HASH_TABLE_OK != h_rc) {
        OAILOG_ERROR_UE(
            LOG_MME_APP, imsi,
            "Error could not update this ue context %p "
            "enb_ue_s1ap_ue_id " ENB_UE_S1AP_ID_FMT
            " mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT " guti " GUTI_FMT " %s\n",
            ue_context_p, ue_context_p->enb_ue_s1ap_id,
            ue_context_p->mme_ue_s1ap_id, GUTI_ARG(guti_p),
            hashtable_rc_code2string(h_rc));
      }
      ue_context_p->emm_context._guti = *guti_p;
    }
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
void mme_ue_context_dump_coll_keys(const mme_ue_context_t* mme_ue_contexts_p) {
  bstring tmp                         = bfromcstr(" ");
  hash_table_ts_t* mme_state_ue_id_ht = get_mme_ue_state();

  btrunc(tmp, 0);
  hashtable_uint64_ts_dump_content(mme_ue_contexts_p->imsi_mme_ue_id_htbl, tmp);
  OAILOG_DEBUG(LOG_MME_APP, "imsi_ue_context_htbl %s\n", bdata(tmp));

  btrunc(tmp, 0);
  hashtable_uint64_ts_dump_content(
      mme_ue_contexts_p->tun11_ue_context_htbl, tmp);
  OAILOG_DEBUG(LOG_MME_APP, "tun11_ue_context_htbl %s\n", bdata(tmp));

  btrunc(tmp, 0);
  hashtable_ts_dump_content(mme_state_ue_id_ht, tmp);
  OAILOG_DEBUG(LOG_MME_APP, "mme_ue_s1ap_id_ue_context_htbl %s\n", bdata(tmp));

  btrunc(tmp, 0);
  hashtable_uint64_ts_dump_content(
      mme_ue_contexts_p->enb_ue_s1ap_id_ue_context_htbl, tmp);
  OAILOG_DEBUG(LOG_MME_APP, "enb_ue_s1ap_id_ue_context_htbl %s\n", bdata(tmp));

  btrunc(tmp, 0);
  obj_hashtable_uint64_ts_dump_content(
      mme_ue_contexts_p->guti_ue_context_htbl, tmp);
  OAILOG_DEBUG(LOG_MME_APP, "guti_ue_context_htbl %s", bdata(tmp));

  bdestroy(tmp);
}

//------------------------------------------------------------------------------
int mme_insert_ue_context(
    mme_ue_context_t* const mme_ue_context_p,
    const struct ue_mm_context_s* const ue_context_p) {
  hashtable_rc_t h_rc                 = HASH_TABLE_OK;
  hash_table_ts_t* mme_state_ue_id_ht = get_mme_ue_state();

  OAILOG_FUNC_IN(LOG_MME_APP);
  if (mme_ue_context_p == NULL) {
    OAILOG_ERROR(LOG_MME_APP, "Invalid MME UE context received\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  if (ue_context_p == NULL) {
    OAILOG_ERROR(LOG_MME_APP, "Invalid UE context received\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  // filled ENB UE S1AP ID
  h_rc = hashtable_uint64_ts_is_key_exists(
      mme_ue_context_p->enb_ue_s1ap_id_ue_context_htbl,
      (const hash_key_t) ue_context_p->enb_s1ap_id_key);
  if (HASH_TABLE_OK == h_rc) {
    OAILOG_WARNING(
        LOG_MME_APP,
        "This ue context %p already exists enb_ue_s1ap_id " ENB_UE_S1AP_ID_FMT
        "\n",
        ue_context_p, ue_context_p->enb_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  h_rc = hashtable_uint64_ts_insert(
      mme_ue_context_p->enb_ue_s1ap_id_ue_context_htbl,
      (const hash_key_t) ue_context_p->enb_s1ap_id_key,
      ue_context_p->mme_ue_s1ap_id);

  if (HASH_TABLE_OK != h_rc) {
    OAILOG_WARNING(
        LOG_MME_APP,
        "Error could not register this ue context %p "
        "enb_ue_s1ap_id " ENB_UE_S1AP_ID_FMT " ue_id 0x%x\n",
        ue_context_p, ue_context_p->enb_ue_s1ap_id,
        ue_context_p->mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  if (INVALID_MME_UE_S1AP_ID != ue_context_p->mme_ue_s1ap_id) {
    h_rc = hashtable_ts_is_key_exists(
        mme_state_ue_id_ht, (const hash_key_t) ue_context_p->mme_ue_s1ap_id);

    if (HASH_TABLE_OK == h_rc) {
      OAILOG_WARNING(
          LOG_MME_APP,
          "This ue context %p already exists mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT
          "\n",
          ue_context_p, ue_context_p->mme_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
    }

    h_rc = hashtable_ts_insert(
        mme_state_ue_id_ht, (const hash_key_t) ue_context_p->mme_ue_s1ap_id,
        (void*) ue_context_p);

    if (HASH_TABLE_OK != h_rc) {
      OAILOG_WARNING(
          LOG_MME_APP,
          "Error could not register this ue context %p "
          "mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT "\n",
          ue_context_p, ue_context_p->mme_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
    }

    // filled IMSI
    if (ue_context_p->emm_context._imsi64) {
      h_rc = hashtable_uint64_ts_insert(
          mme_ue_context_p->imsi_mme_ue_id_htbl,
          (const hash_key_t) ue_context_p->emm_context._imsi64,
          ue_context_p->mme_ue_s1ap_id);

      if (HASH_TABLE_OK != h_rc) {
        OAILOG_WARNING_UE(
            LOG_MME_APP, ue_context_p->emm_context._imsi64,
            "Error could not register this ue context %p "
            "mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT " imsi " IMSI_64_FMT "\n",
            ue_context_p, ue_context_p->mme_ue_s1ap_id,
            ue_context_p->emm_context._imsi64);
        OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
      }
    }

    // filled S11 tun id
    if (ue_context_p->mme_teid_s11) {
      h_rc = hashtable_uint64_ts_insert(
          mme_ue_context_p->tun11_ue_context_htbl,
          (const hash_key_t) ue_context_p->mme_teid_s11,
          ue_context_p->mme_ue_s1ap_id);

      if (HASH_TABLE_OK != h_rc) {
        OAILOG_WARNING(
            LOG_MME_APP,
            "Error could not register this ue context %p "
            "mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT " mme_teid_s11 " TEID_FMT "\n",
            ue_context_p, ue_context_p->mme_ue_s1ap_id,
            ue_context_p->mme_teid_s11);
        OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
      }
    }

    // filled guti
    if ((0 != ue_context_p->emm_context._guti.gummei.mme_code) ||
        (0 != ue_context_p->emm_context._guti.gummei.mme_gid) ||
        (INVALID_TMSI != ue_context_p->emm_context._guti.m_tmsi) ||
        (0 != ue_context_p->emm_context._guti.gummei.plmn
                  .mcc_digit1) ||  // MCC 000 does not exist in ITU table
        (0 != ue_context_p->emm_context._guti.gummei.plmn.mcc_digit2) ||
        (0 != ue_context_p->emm_context._guti.gummei.plmn.mcc_digit3)) {
      h_rc = obj_hashtable_uint64_ts_insert(
          mme_ue_context_p->guti_ue_context_htbl,
          (const void* const) & ue_context_p->emm_context._guti,
          sizeof(ue_context_p->emm_context._guti),
          ue_context_p->mme_ue_s1ap_id);

      if (HASH_TABLE_OK != h_rc) {
        OAILOG_WARNING(
            LOG_MME_APP,
            "Error could not register this ue context %p "
            "mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT " guti " GUTI_FMT "\n",
            ue_context_p, ue_context_p->mme_ue_s1ap_id,
            GUTI_ARG(&ue_context_p->emm_context._guti));
        OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
      }
    }
  }

  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

//------------------------------------------------------------------------------
void mme_remove_ue_context(
    mme_ue_context_t* const mme_ue_context_p,
    struct ue_mm_context_s* ue_context_p) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  hashtable_rc_t hash_rc              = HASH_TABLE_OK;
  hash_table_ts_t* mme_state_ue_id_ht = get_mme_ue_state();

  if (!mme_ue_context_p) {
    OAILOG_ERROR(LOG_MME_APP, "Invalid MME UE context received\n");
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  if (!ue_context_p) {
    OAILOG_ERROR(LOG_MME_APP, "Invalid UE context received\n");
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }

  // First, notify directoryd of removal
  directoryd_remove_location_a(
      ue_context_p->emm_context._imsi64,
      ue_context_p->emm_context._imsi.length);

  // Release emm and esm context
  delete_mme_ue_state(ue_context_p->emm_context._imsi64);
  mme_app_ue_context_free_content(ue_context_p);
  // IMSI
  if (ue_context_p->emm_context._imsi64) {
    hash_rc = hashtable_uint64_ts_remove(
        mme_ue_context_p->imsi_mme_ue_id_htbl,
        (const hash_key_t) ue_context_p->emm_context._imsi64);
    if (HASH_TABLE_OK != hash_rc) {
      OAILOG_ERROR_UE(
          LOG_MME_APP, ue_context_p->emm_context._imsi64,
          "UE context not found!\n"
          " enb_ue_s1ap_id " ENB_UE_S1AP_ID_FMT
          " mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT " not in IMSI collection\n",
          ue_context_p->enb_ue_s1ap_id, ue_context_p->mme_ue_s1ap_id);
    }
  }
  // filled guti
  if ((ue_context_p->emm_context._guti.gummei.mme_code) ||
      (ue_context_p->emm_context._guti.gummei.mme_gid) ||
      (ue_context_p->emm_context._guti.m_tmsi) ||
      (ue_context_p->emm_context._guti.gummei.plmn.mcc_digit1) ||
      (ue_context_p->emm_context._guti.gummei.plmn.mcc_digit2) ||
      (ue_context_p->emm_context._guti.gummei.plmn
           .mcc_digit3)) {  // MCC 000 does not exist in ITU table
    hash_rc = obj_hashtable_uint64_ts_remove(
        mme_ue_context_p->guti_ue_context_htbl,
        (const void* const) & ue_context_p->emm_context._guti,
        sizeof(ue_context_p->emm_context._guti));
    if (HASH_TABLE_OK != hash_rc)
      OAILOG_ERROR(
          LOG_MME_APP,
          "UE Context not found!\n"
          " enb_ue_s1ap_id " ENB_UE_S1AP_ID_FMT
          " mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT
          ", GUTI  not in GUTI collection\n",
          ue_context_p->enb_ue_s1ap_id, ue_context_p->mme_ue_s1ap_id);
  }

  _clear_emm_ctxt(&ue_context_p->emm_context);
  // eNB UE S1P UE ID
  hash_rc = hashtable_uint64_ts_remove(
      mme_ue_context_p->enb_ue_s1ap_id_ue_context_htbl,
      (const hash_key_t) ue_context_p->enb_s1ap_id_key);
  if (HASH_TABLE_OK != hash_rc)
    OAILOG_ERROR(
        LOG_MME_APP,
        "UE context not found!\n"
        " enb_ue_s1ap_id " ENB_UE_S1AP_ID_FMT
        " mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT
        ", ENB_UE_S1AP_ID not in ENB_UE_S1AP_ID collection",
        ue_context_p->enb_ue_s1ap_id, ue_context_p->mme_ue_s1ap_id);

  // filled S11 tun id
  if (ue_context_p->mme_teid_s11) {
    hash_rc = hashtable_uint64_ts_remove(
        mme_ue_context_p->tun11_ue_context_htbl,
        (const hash_key_t) ue_context_p->mme_teid_s11);
    if (HASH_TABLE_OK != hash_rc)
      OAILOG_ERROR(
          LOG_MME_APP,
          "UE Context not found!\n"
          " enb_ue_s1ap_id " ENB_UE_S1AP_ID_FMT
          " mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT ", MME S11 TEID  " TEID_FMT
          "  not in S11 collection\n",
          ue_context_p->enb_ue_s1ap_id, ue_context_p->mme_ue_s1ap_id,
          ue_context_p->mme_teid_s11);
  }
  // filled NAS UE ID/ MME UE S1AP ID
  if (INVALID_MME_UE_S1AP_ID != ue_context_p->mme_ue_s1ap_id) {
    hash_rc = hashtable_ts_remove(
        mme_state_ue_id_ht, (const hash_key_t) ue_context_p->mme_ue_s1ap_id,
        (void**) &ue_context_p);
    if (HASH_TABLE_OK != hash_rc)
      OAILOG_ERROR(
          LOG_MME_APP,
          "UE context not found!\n"
          "  enb_ue_s1ap_id " ENB_UE_S1AP_ID_FMT
          ", mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT
          " not in MME UE S1AP ID collection",
          ue_context_p->enb_ue_s1ap_id, ue_context_p->mme_ue_s1ap_id);
  }

  free_wrapper((void**) &ue_context_p);
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
void mme_app_dump_protocol_configuration_options(
    const protocol_configuration_options_t* const pco,
    const bool ms2network_direction, const uint8_t indent_spaces,
    bstring bstr_dump) {
  int i = 0;

  if (pco) {
    bformata(bstr_dump, "        Protocol configuration options:\n");
    bformata(
        bstr_dump, "        Configuration protocol .......: %" PRIx8 "\n",
        pco->configuration_protocol);
    while (i < pco->num_protocol_or_container_id) {
      switch (pco->protocol_or_container_ids[i].id) {
        case PCO_PI_LCP:
          bformata(bstr_dump, "        Protocol ID .......: LCP\n");
          break;
        case PCO_PI_PAP:
          bformata(bstr_dump, "        Protocol ID .......: PAP\n");
          break;
        case PCO_PI_CHAP:
          bformata(bstr_dump, "        Protocol ID .......: CHAP\n");
          break;
        case PCO_PI_IPCP:
          bformata(bstr_dump, "        Protocol ID .......: IPCP\n");
          break;

        default:
          if (ms2network_direction) {
            switch (pco->protocol_or_container_ids[i].id) {
              case PCO_CI_P_CSCF_IPV6_ADDRESS_REQUEST:
                bformata(
                    bstr_dump,
                    "        Container ID .......: "
                    "P_CSCF_IPV6_ADDRESS_REQUEST\n");
                break;
              case PCO_CI_DNS_SERVER_IPV6_ADDRESS_REQUEST:
                bformata(
                    bstr_dump,
                    "        Container ID .......: "
                    "DNS_SERVER_IPV6_ADDRESS_REQUEST\n");
                break;
              case PCO_CI_MS_SUPPORT_OF_NETWORK_REQUESTED_BEARER_CONTROL_INDICATOR:
                bformata(
                    bstr_dump,
                    "        Container ID .......: "
                    "MS_SUPPORT_OF_NETWORK_REQUESTED_BEARER_CONTROL_"
                    "INDICATOR\n");
                break;
              case PCO_CI_DSMIPV6_HOME_AGENT_ADDRESS_REQUEST:
                bformata(
                    bstr_dump,
                    "        Container ID .......: "
                    "DSMIPV6_HOME_AGENT_ADDRESS_REQUEST\n");
                break;
              case PCO_CI_DSMIPV6_HOME_NETWORK_PREFIX_REQUEST:
                bformata(
                    bstr_dump,
                    "        Container ID .......: "
                    "DSMIPV6_HOME_NETWORK_PREFIX_REQUEST\n");
                break;
              case PCO_CI_DSMIPV6_IPV4_HOME_AGENT_ADDRESS_REQUEST:
                bformata(
                    bstr_dump,
                    "        Container ID .......: "
                    "DSMIPV6_IPV4_HOME_AGENT_ADDRESS_REQUEST\n");
                break;
              case PCO_CI_IP_ADDRESS_ALLOCATION_VIA_NAS_SIGNALLING:
                bformata(
                    bstr_dump,
                    "        Container ID .......: "
                    "IP_ADDRESS_ALLOCATION_VIA_NAS_SIGNALLING\n");
                break;
              case PCO_CI_IPV4_ADDRESS_ALLOCATION_VIA_DHCPV4:
                bformata(
                    bstr_dump,
                    "        Container ID .......: "
                    "IPV4_ADDRESS_ALLOCATION_VIA_DHCPV4\n");
                break;
              case PCO_CI_P_CSCF_IPV4_ADDRESS_REQUEST:
                bformata(
                    bstr_dump,
                    "        Container ID .......: "
                    "P_CSCF_IPV4_ADDRESS_REQUEST\n");
                break;
              case PCO_CI_DNS_SERVER_IPV4_ADDRESS_REQUEST:
                bformata(
                    bstr_dump,
                    "        Container ID .......: "
                    "DNS_SERVER_IPV4_ADDRESS_REQUEST\n");
                break;
              case PCO_CI_MSISDN_REQUEST:
                bformata(
                    bstr_dump,
                    "        Container ID .......: MSISDN_REQUEST\n");
                break;
              case PCO_CI_IFOM_SUPPORT_REQUEST:
                bformata(
                    bstr_dump,
                    "        Container ID .......: IFOM_SUPPORT_REQUEST\n");
                break;
              case PCO_CI_IPV4_LINK_MTU_REQUEST:
                bformata(
                    bstr_dump,
                    "        Container ID .......: IPV4_LINK_MTU_REQUEST\n");
                break;
              case PCO_CI_IM_CN_SUBSYSTEM_SIGNALING_FLAG:
                bformata(
                    bstr_dump,
                    "        Container ID .......: "
                    "IM_CN_SUBSYSTEM_SIGNALING_FLAG\n");
                break;
              default:
                bformata(
                    bstr_dump, "       Unhandled container id %u length %d\n",
                    pco->protocol_or_container_ids[i].id,
                    blength(pco->protocol_or_container_ids[i].contents));
            }
          } else {
            switch (pco->protocol_or_container_ids[i].id) {
              case PCO_CI_P_CSCF_IPV6_ADDRESS:
                bformata(
                    bstr_dump,
                    "        Container ID .......: P_CSCF_IPV6_ADDRESS\n");
                break;
              case PCO_CI_DNS_SERVER_IPV6_ADDRESS:
                bformata(
                    bstr_dump,
                    "        Container ID .......: DNS_SERVER_IPV6_ADDRESS\n");
                break;
              case PCO_CI_POLICY_CONTROL_REJECTION_CODE:
                bformata(
                    bstr_dump,
                    "        Container ID .......: "
                    "POLICY_CONTROL_REJECTION_CODE\n");
                break;
              case PCO_CI_SELECTED_BEARER_CONTROL_MODE:
                bformata(
                    bstr_dump,
                    "        Container ID .......: "
                    "SELECTED_BEARER_CONTROL_MODE\n");
                break;
              case PCO_CI_DSMIPV6_HOME_AGENT_ADDRESS:
                bformata(
                    bstr_dump,
                    "        Container ID .......: "
                    "DSMIPV6_HOME_AGENT_ADDRESS\n");
                break;
              case PCO_CI_DSMIPV6_HOME_NETWORK_PREFIX:
                bformata(
                    bstr_dump,
                    "        Container ID .......: "
                    "DSMIPV6_HOME_NETWORK_PREFIX\n");
                break;
              case PCO_CI_DSMIPV6_IPV4_HOME_AGENT_ADDRESS:
                bformata(
                    bstr_dump,
                    "        Container ID .......: "
                    "DSMIPV6_IPV4_HOME_AGENT_ADDRESS\n");
                break;
              case PCO_CI_P_CSCF_IPV4_ADDRESS:
                bformata(
                    bstr_dump,
                    "        Container ID .......: P_CSCF_IPV4_ADDRESS\n");
                break;
              case PCO_CI_DNS_SERVER_IPV4_ADDRESS:
                bformata(
                    bstr_dump,
                    "        Container ID .......: DNS_SERVER_IPV4_ADDRESS\n");
                break;
              case PCO_CI_MSISDN:
                bformata(bstr_dump, "        Container ID .......: MSISDN\n");
                break;
              case PCO_CI_IFOM_SUPPORT:
                bformata(
                    bstr_dump, "        Container ID .......: IFOM_SUPPORT\n");
                break;
              case PCO_CI_IPV4_LINK_MTU:
                bformata(
                    bstr_dump, "        Container ID .......: IPV4_LINK_MTU\n");
                break;
              case PCO_CI_IM_CN_SUBSYSTEM_SIGNALING_FLAG:
                bformata(
                    bstr_dump,
                    "        Container ID .......: "
                    "IM_CN_SUBSYSTEM_SIGNALING_FLAG\n");
                break;
              default:
                bformata(
                    bstr_dump, "       Unhandled container id %u length %d\n",
                    pco->protocol_or_container_ids[i].id,
                    blength(pco->protocol_or_container_ids[i].contents));
            }
          }
      }
      OAILOG_STREAM_HEX(
          OAILOG_LEVEL_DEBUG, LOG_MME_APP, "        Hex data: ",
          bdata(pco->protocol_or_container_ids[i].contents),
          blength(pco->protocol_or_container_ids[i].contents));
      i++;
    }
  }
}

//------------------------------------------------------------------------------
void mme_app_dump_bearer_context(
    const bearer_context_t* const bc, uint8_t indent_spaces,
    bstring bstr_dump) {
  bformata(
      bstr_dump, "%*s - Bearer id .......: %02u\n", indent_spaces, " ",
      bc->ebi);
  bformata(
      bstr_dump, "%*s - Transaction ID ..: %x\n", indent_spaces, " ",
      bc->transaction_identifier);
  if (bc->s_gw_fteid_s1u.ipv4) {
    char ipv4[INET_ADDRSTRLEN];
    inet_ntop(
        AF_INET, (void*) &bc->s_gw_fteid_s1u.ipv4_address.s_addr, ipv4,
        INET_ADDRSTRLEN);
    bformata(
        bstr_dump, "%*s - S-GW S1-U IPv4 Address...: [%s]\n", indent_spaces,
        " ", ipv4);
  } else if (bc->s_gw_fteid_s1u.ipv6) {
    char ipv6[INET6_ADDRSTRLEN];
    inet_ntop(
        AF_INET6, &bc->s_gw_fteid_s1u.ipv6_address, ipv6, INET6_ADDRSTRLEN);
    bformata(
        bstr_dump, "%*s - S-GW S1-U IPv6 Address...: [%s]\n", indent_spaces,
        " ", ipv6);
  }
  bformata(
      bstr_dump, "%*s - S-GW TEID (UP)...: " TEID_FMT "\n", indent_spaces, " ",
      bc->s_gw_fteid_s1u.teid);
  if (bc->p_gw_fteid_s5_s8_up.ipv4) {
    char ipv4[INET_ADDRSTRLEN];
    inet_ntop(
        AF_INET, (void*) &bc->p_gw_fteid_s5_s8_up.ipv4_address.s_addr, ipv4,
        INET_ADDRSTRLEN);
    bformata(bstr_dump, "%*s - P-GW S5-S8 IPv4..: [%s]\n", ipv4);
  } else if (bc->p_gw_fteid_s5_s8_up.ipv6) {
    char ipv6[INET6_ADDRSTRLEN];
    inet_ntop(
        AF_INET6, &bc->p_gw_fteid_s5_s8_up.ipv6_address, ipv6,
        INET6_ADDRSTRLEN);
    bformata(
        bstr_dump, "%*s - P-GW S5-S8 IPv6..: [%s]\n", indent_spaces, " ", ipv6);
  }
  bformata(
      bstr_dump, "%*s - P-GW TEID S5-S8..: " TEID_FMT "\n", indent_spaces, " ",
      bc->p_gw_fteid_s5_s8_up.teid);
  bformata(
      bstr_dump, "%*s - QCI .............: %u\n", indent_spaces, " ", bc->qci);
  bformata(
      bstr_dump, "%*s - Priority level ..: %u\n", indent_spaces, " ",
      bc->priority_level);
  bformata(
      bstr_dump, "%*s - Pre-emp vul .....: %s\n", indent_spaces, " ",
      (bc->preemption_vulnerability == PRE_EMPTION_VULNERABILITY_ENABLED) ?
          "ENABLED" :
          "DISABLED");
  bformata(
      bstr_dump, "%*s - Pre-emp cap .....: %s\n", indent_spaces, " ",
      (bc->preemption_capability == PRE_EMPTION_CAPABILITY_ENABLED) ?
          "ENABLED" :
          "DISABLED");
  bformata(
      bstr_dump, "%*s - GBR UL ..........: %010" PRIu64 "\n", indent_spaces,
      " ", bc->esm_ebr_context.gbr_ul);
  bformata(
      bstr_dump, "%*s - GBR DL ..........: %010" PRIu64 "\n", indent_spaces,
      " ", bc->esm_ebr_context.gbr_dl);
  bformata(
      bstr_dump, "%*s - MBR UL ..........: %010" PRIu64 "\n", indent_spaces,
      " ", bc->esm_ebr_context.mbr_ul);
  bformata(
      bstr_dump, "%*s - MBR DL ..........: %010" PRIu64 "\n", indent_spaces,
      " ", bc->esm_ebr_context.mbr_dl);
  bformata(
      bstr_dump, "%*s - " ANSI_COLOR_BOLD_ON "NAS ESM bearer private data .:\n",
      indent_spaces, " ");
  bformata(
      bstr_dump, "%*s -     ESM State .......: %s\n", indent_spaces, " ",
      esm_ebr_state2string(bc->esm_ebr_context.status));
  bformata(
      bstr_dump, "%*s -     Timer id ........: %lx\n", indent_spaces, " ",
      bc->esm_ebr_context.timer.id);
  bformata(
      bstr_dump, "%*s -     Timer TO(seconds): %ld\n", indent_spaces, " ",
      bc->esm_ebr_context.timer.sec);
  bformata(
      bstr_dump, "%*s - PDN id ..........: %u\n", indent_spaces, " ",
      bc->pdn_cx_id);
}

//------------------------------------------------------------------------------
void mme_app_dump_pdn_context(
    const struct ue_mm_context_s* const ue_mm_context,
    const pdn_context_t* const pdn_context, const pdn_cid_t pdn_cid,
    const uint8_t indent_spaces, bstring bstr_dump) {
  if (pdn_context) {
    bformata(bstr_dump, "%*s - PDN ID %u:\n", indent_spaces, " ", pdn_cid);
    bformata(
        bstr_dump, "%*s - Context Identifier .: %x\n", indent_spaces, " ",
        pdn_context->context_identifier);
    bformata(
        bstr_dump, "%*s - Is active          .: %s\n", indent_spaces, " ",
        (pdn_context->is_active) ? "yes" : "no");
    bformata(
        bstr_dump, "%*s - APN in use .........: %s\n", indent_spaces, " ",
        bdata(pdn_context->apn_in_use));
    bformata(
        bstr_dump, "%*s - APN subscribed......: %s\n", indent_spaces, " ",
        bdata(pdn_context->apn_subscribed));
    bformata(
        bstr_dump, "%*s - APN OI replacement .: %s\n", indent_spaces, " ",
        bdata(pdn_context->apn_oi_replacement));

    bformata(
        bstr_dump, "%*s - PDN type ...........: %s\n", indent_spaces, " ",
        PDN_TYPE_TO_STRING(pdn_context->paa.pdn_type));
    if (pdn_context->paa.pdn_type == IPv4) {
      char ipv4[INET_ADDRSTRLEN];
      inet_ntop(
          AF_INET, (void*) &pdn_context->paa.ipv4_address, ipv4,
          INET_ADDRSTRLEN);
      bformata(
          bstr_dump, "%*s - PAA (IPv4)..........: %s\n", indent_spaces, " ",
          ipv4);
    } else {
      char ipv6[INET6_ADDRSTRLEN];
      inet_ntop(
          AF_INET6, &pdn_context->paa.ipv6_address, ipv6, INET6_ADDRSTRLEN);
      bformata(
          bstr_dump, "%*s - PAA (IPv6)..........: %s\n", indent_spaces, " ",
          ipv6);
    }
    if (pdn_context->p_gw_address_s5_s8_cp.pdn_type == IPv4) {
      char ipv4[INET_ADDRSTRLEN];
      inet_ntop(
          AF_INET,
          (void*) &pdn_context->p_gw_address_s5_s8_cp.address.ipv4_address,
          ipv4, INET_ADDRSTRLEN);
      bformata(
          bstr_dump, "%*s - P-GW s5 s8 cp (IPv4): %s\n", indent_spaces, " ",
          ipv4);
    } else {
      char ipv6[INET6_ADDRSTRLEN];
      inet_ntop(
          AF_INET6, &pdn_context->p_gw_address_s5_s8_cp.address.ipv6_address,
          ipv6, INET6_ADDRSTRLEN);
      bformata(
          bstr_dump, "%*s - P-GW s5 s8 cp (IPv6): %s\n", indent_spaces, " ",
          ipv6);
    }
    bformata(
        bstr_dump, "%*s - P-GW TEID s5 s8 cp .: " TEID_FMT "\n", indent_spaces,
        " ", pdn_context->p_gw_teid_s5_s8_cp);
    if (pdn_context->s_gw_address_s11_s4.pdn_type == IPv4) {
      char ipv4[INET_ADDRSTRLEN];
      inet_ntop(
          AF_INET,
          (void*) &pdn_context->s_gw_address_s11_s4.address.ipv4_address, ipv4,
          INET_ADDRSTRLEN);
      bformata(
          bstr_dump, "%*s - S-GW s11_s4 (IPv4) .: %s\n", indent_spaces, " ",
          ipv4);
    } else {
      char ipv6[INET6_ADDRSTRLEN];
      inet_ntop(
          AF_INET6, &pdn_context->s_gw_address_s11_s4.address.ipv6_address,
          ipv6, INET6_ADDRSTRLEN);
      bformata(
          bstr_dump, "%*s - S-GW s11_s4 (IPv6) .: %s\n", indent_spaces, " ",
          indent_spaces, " ", ipv6);
    }
    bformata(
        bstr_dump, "%*s - S-GW TEID s5 s8 cp .: " TEID_FMT "\n", indent_spaces,
        " ", pdn_context->s_gw_teid_s11_s4);

    bformata(
        bstr_dump, "%*s - Default bearer eps subscribed qos profile:\n",
        indent_spaces, " ");
    bformata(
        bstr_dump, "%*s     - QCI ......................: %u\n", indent_spaces,
        " ", pdn_context->default_bearer_eps_subscribed_qos_profile.qci);
    bformata(
        bstr_dump, "%*s     - Priority level ...........: %u\n", indent_spaces,
        " ",
        pdn_context->default_bearer_eps_subscribed_qos_profile
            .allocation_retention_priority.priority_level);
    bformata(
        bstr_dump, "%*s     - Pre-emp vulnerabil .......: %s\n", indent_spaces,
        " ",
        (pdn_context->default_bearer_eps_subscribed_qos_profile
             .allocation_retention_priority.pre_emp_vulnerability ==
         PRE_EMPTION_VULNERABILITY_ENABLED) ?
            "ENABLED" :
            "DISABLED");
    bformata(
        bstr_dump, "%*s     - Pre-emp capability .......: %s\n", indent_spaces,
        " ",
        (pdn_context->default_bearer_eps_subscribed_qos_profile
             .allocation_retention_priority.pre_emp_capability ==
         PRE_EMPTION_CAPABILITY_ENABLED) ?
            "ENABLED" :
            "DISABLED");
    bformata(
        bstr_dump, "%*s     - APN-AMBR (bits/s) DL .....: %010" PRIu64 "\n",
        indent_spaces, " ", pdn_context->subscribed_apn_ambr.br_dl);
    bformata(
        bstr_dump, "%*s     - APN-AMBR (bits/s) UL .....: %010" PRIu64 "\n",
        indent_spaces, " ", pdn_context->subscribed_apn_ambr.br_ul);
    bformata(
        bstr_dump, "%*s     - P-GW-APN-AMBR (bits/s) DL : %010" PRIu64 "\n",
        indent_spaces, " ", pdn_context->p_gw_apn_ambr.br_dl);
    bformata(
        bstr_dump, "%*s     - P-GW-APN-AMBR (bits/s) UL : %010" PRIu64 "\n",
        indent_spaces, " ", pdn_context->p_gw_apn_ambr.br_ul);

    bformata(
        bstr_dump, "%*s     - P-GW-APN-AMBR (bits/s) UL : %010" PRIu64 "\n",
        indent_spaces, " ", pdn_context->p_gw_apn_ambr.br_ul,
        (pdn_context->p_gw_apn_ambr.br_unit) ? "KBPS" : "BPS");

    bformata(
        bstr_dump, "%*s     - Default EBI ..............: %u\n", indent_spaces,
        " ", pdn_context->default_ebi);
    bformata(bstr_dump, "%*s - NAS ESM private data:\n");
    bformata(
        bstr_dump, "%*s     - Procedure transaction ID .: %x\n", indent_spaces,
        " ", pdn_context->esm_data.pti);
    bformata(
        bstr_dump, "%*s     -  Is emergency .............: %s\n", indent_spaces,
        " ", (pdn_context->esm_data.is_emergency) ? "yes" : "no");
    bformata(
        bstr_dump, "%*s     -  APN AMBR .................: %u\n", indent_spaces,
        " ", pdn_context->esm_data.ambr);
    bformata(
        bstr_dump, "%*s     -  Addr realloc allowed......: %s\n", indent_spaces,
        " ", (pdn_context->esm_data.addr_realloc) ? "yes" : "no");
    bformata(
        bstr_dump, "%*s     -  Num allocated EPS bearers.: %d\n", indent_spaces,
        " ", pdn_context->esm_data.n_bearers);
    bformata(bstr_dump, "%*s - Bearer List:\n");
    for (int bindex = 0; bindex < 0; bindex++) {
      // should be equal to bindex if valid
      int bcindex = pdn_context->bearer_contexts[bindex];
      if ((0 <= bcindex) && (BEARERS_PER_UE > bcindex)) {
        if (bindex != bcindex) {
          OAILOG_ERROR_UE(
              LOG_MME_APP, ue_mm_context->emm_context._imsi64,
              "Mismatch in configuration. PDN index (%i) != Bearer index "
              "(%i)\n",
              bindex, bcindex);
          OAILOG_FUNC_OUT(LOG_MME_APP);
        }

        bearer_context_t* bc = ue_mm_context->bearer_contexts[bcindex];
        if (!bc) {
          OAILOG_ERROR_UE(
              LOG_MME_APP, ue_mm_context->emm_context._imsi64,
              "Mismatch in configuration. Bearer context is NULL\n");
          OAILOG_FUNC_OUT(LOG_MME_APP);
        }
        bformata(bstr_dump, "%*s - Bearer item ----------------------------\n");
        mme_app_dump_bearer_context(bc, indent_spaces + 4, bstr_dump);
      }
    }
  }
}

//-------------------------------------------------------------------------------------------------------
void mme_ue_context_update_ue_sig_connection_state(
    mme_ue_context_t* const mme_ue_context_p,
    struct ue_mm_context_s* ue_context_p, ecm_state_t new_ecm_state) {
  // Function is used to update UE's Signaling Connection State
  hashtable_rc_t hash_rc = HASH_TABLE_OK;

  OAILOG_FUNC_IN(LOG_MME_APP);
  if (mme_ue_context_p == NULL) {
    OAILOG_ERROR(LOG_MME_APP, "Invalid MME UE context received\n");
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  if (ue_context_p == NULL) {
    OAILOG_ERROR(LOG_MME_APP, "Invalid UE context received\n");
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  if (ue_context_p->ecm_state == ECM_CONNECTED && new_ecm_state == ECM_IDLE) {
    hash_rc = hashtable_uint64_ts_remove(
        mme_ue_context_p->enb_ue_s1ap_id_ue_context_htbl,
        (const hash_key_t) ue_context_p->enb_s1ap_id_key);
    if (HASH_TABLE_OK != hash_rc) {
      OAILOG_WARNING_UE(
          LOG_MME_APP, ue_context_p->emm_context._imsi64,
          "UE context enb_ue_s1ap_ue_id_key %ld "
          "mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT
          ", ENB_UE_S1AP_ID_KEY could not be found",
          ue_context_p->enb_s1ap_id_key, ue_context_p->mme_ue_s1ap_id);
    }
    ue_context_p->enb_s1ap_id_key = INVALID_ENB_UE_S1AP_ID_KEY;

    OAILOG_DEBUG_UE(
        LOG_MME_APP, ue_context_p->emm_context._imsi64,
        "MME_APP: UE Connection State changed to IDLE. mme_ue_s1ap_id "
        "= " MME_UE_S1AP_ID_FMT "\n",
        ue_context_p->mme_ue_s1ap_id);

    if (mme_config.nas_config.t3412_min > 0) {
      // Start Mobile reachability timer only if peroidic TAU timer is not
      // disabled
      if (ue_context_p->mobile_reachability_timer.id ==
          MME_APP_TIMER_INACTIVE_ID) {
        // Start Mobile Reachability timer only if it is not running
        if ((ue_context_p->mobile_reachability_timer.id = mme_app_start_timer(
                 ue_context_p->mobile_reachability_timer.sec * 1000,
                 TIMER_REPEAT_ONCE,
                 mme_app_handle_mobile_reachability_timer_expiry,
                 ue_context_p->mme_ue_s1ap_id)) == -1) {
          OAILOG_ERROR_UE(
              LOG_MME_APP, ue_context_p->emm_context._imsi64,
              "Failed to start Mobile Reachability timer for UE id "
              " " MME_UE_S1AP_ID_FMT "\n",
              ue_context_p->mme_ue_s1ap_id);
          ue_context_p->mobile_reachability_timer.id =
              MME_APP_TIMER_INACTIVE_ID;
        } else {
          ue_context_p->time_mobile_reachability_timer_started = time(NULL);
          OAILOG_DEBUG_UE(
              LOG_MME_APP, ue_context_p->emm_context._imsi64,
              "Started Mobile Reachability timer for UE id " MME_UE_S1AP_ID_FMT
              "\n",
              ue_context_p->mme_ue_s1ap_id);
        }
      } else {
        OAILOG_DEBUG_UE(
            LOG_MME_APP, ue_context_p->emm_context._imsi64,
            "Mobile Reachability timer is already started for UE id:"
            " " MME_UE_S1AP_ID_FMT "\n",
            ue_context_p->mme_ue_s1ap_id);
      }
    }
    ue_context_p->ecm_state = ECM_IDLE;
    // Update Stats
    update_mme_app_stats_connected_ue_sub();
    OAILOG_INFO_UE(
        LOG_MME_APP, ue_context_p->emm_context._imsi64, "UE STATE - IDLE.\n");

  } else if (
      (ue_context_p->ecm_state == ECM_IDLE) &&
      (new_ecm_state == ECM_CONNECTED)) {
    ue_context_p->ecm_state = ECM_CONNECTED;

    OAILOG_DEBUG_UE(
        LOG_MME_APP, ue_context_p->emm_context._imsi64,
        "MME_APP: UE Connection State changed to CONNECTED.enb_ue_s1ap_id "
        "=" ENB_UE_S1AP_ID_FMT ", mme_ue_s1ap_id = " MME_UE_S1AP_ID_FMT "\n",
        ue_context_p->enb_ue_s1ap_id, ue_context_p->mme_ue_s1ap_id);
    // Set PPF flag to true whenever UE moves from ECM_IDLE to ECM_CONNECTED
    // state
    ue_context_p->ppf = true;
    // Stop Mobile reachability timer,if running
    if (ue_context_p->mobile_reachability_timer.id !=
        MME_APP_TIMER_INACTIVE_ID) {
      mme_app_stop_timer(ue_context_p->mobile_reachability_timer.id);
      ue_context_p->mobile_reachability_timer.id = MME_APP_TIMER_INACTIVE_ID;
    }
    // Stop Implicit detach timer,if running
    if (ue_context_p->implicit_detach_timer.id != MME_APP_TIMER_INACTIVE_ID) {
      mme_app_stop_timer(ue_context_p->implicit_detach_timer.id);
      ue_context_p->implicit_detach_timer.id = MME_APP_TIMER_INACTIVE_ID;
      ue_context_p->time_implicit_detach_timer_started = 0;
    }
    // Update Stats
    update_mme_app_stats_connected_ue_add();
    OAILOG_INFO_UE(
        LOG_MME_APP, ue_context_p->emm_context._imsi64,
        "UE STATE - CONNECTED.\n");
  } else if (ue_context_p->ecm_state == ECM_IDLE && new_ecm_state == ECM_IDLE) {
    OAILOG_INFO_UE(
        LOG_MME_APP, ue_context_p->emm_context._imsi64,
        "Old UE ECM State (IDLE) is same as the new UE ECM state (IDLE)\n");
    hash_rc = hashtable_uint64_ts_remove(
        mme_ue_context_p->enb_ue_s1ap_id_ue_context_htbl,
        (const hash_key_t) ue_context_p->enb_s1ap_id_key);
    if (HASH_TABLE_OK != hash_rc) {
      OAILOG_WARNING_UE(
          LOG_MME_APP, ue_context_p->emm_context._imsi64,
          "UE context enb_ue_s1ap_ue_id_key %ld "
          "mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT
          ", ENB_UE_S1AP_ID_KEY could not be found",
          ue_context_p->enb_s1ap_id_key, ue_context_p->mme_ue_s1ap_id);
    }
    ue_context_p->enb_s1ap_id_key = INVALID_ENB_UE_S1AP_ID_KEY;
  } else {
    OAILOG_INFO_UE(
        LOG_MME_APP, ue_context_p->emm_context._imsi64,
        "Old UE ECM State (CONNECTED) is same as the new UE ECM state "
        "(CONNECTED)\n");
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}
//------------------------------------------------------------------------------
bool mme_app_dump_ue_context(
    const hash_key_t keyP, void* const ue_mm_context_pP, void* unused_param_pP,
    void** unused_result_pP)
//------------------------------------------------------------------------------
{
  struct ue_mm_context_s* const ue_mm_context =
      (struct ue_mm_context_s*) ue_mm_context_pP;
  uint8_t j = 0;

  if (ue_mm_context) {
    bstring bstr_dump =
        bfromcstralloc(4096, "\n-----------------------UE MM context ");
    bformata(bstr_dump, "%p --------------------\n", ue_mm_context);
    bformata(
        bstr_dump, "    - eNB UE s1ap ID .: %08x\n",
        ue_mm_context->enb_ue_s1ap_id);
    bformata(
        bstr_dump, "    - MME UE s1ap ID .: %08x\n",
        ue_mm_context->mme_ue_s1ap_id);
    bformata(
        bstr_dump, "    - MME S11 TEID ...: " TEID_FMT "\n",
        ue_mm_context->mme_teid_s11);
    bformata(
        bstr_dump, "                        | mcc | mnc | cell identity |\n");
    bformata(
        bstr_dump,
        "    - E-UTRAN CGI ....: | %u%u%u | %u%u%c | %05x.%02x    |\n",
        ue_mm_context->e_utran_cgi.plmn.mcc_digit1,
        ue_mm_context->e_utran_cgi.plmn.mcc_digit2,
        ue_mm_context->e_utran_cgi.plmn.mcc_digit3,
        ue_mm_context->e_utran_cgi.plmn.mnc_digit1,
        ue_mm_context->e_utran_cgi.plmn.mnc_digit2,
        (ue_mm_context->e_utran_cgi.plmn.mnc_digit3 > 9) ?
            ' ' :
            0x30 + ue_mm_context->e_utran_cgi.plmn.mnc_digit3,
        ue_mm_context->e_utran_cgi.cell_identity.enb_id,
        ue_mm_context->e_utran_cgi.cell_identity.cell_id);
    /*
     * Ctime return a \n in the string
     */
    bformata(
        bstr_dump, "    - Last acquired ..: %s",
        ctime(&ue_mm_context->cell_age));

    emm_context_dump(&ue_mm_context->emm_context, 4, bstr_dump);
    /*
     * Display UE info only if we know them
     */
    if (SUBSCRIPTION_KNOWN == ue_mm_context->subscription_known) {
      /* TODO bformata (bstr_dump, "    - Status .........: %s\n",
       * (ue_mm_context->subscriber_status == SS_SERVICE_GRANTED) ?
       * "Granted" : "Barred");
       */
#define DISPLAY_BIT_MASK_PRESENT(mASK)                                         \
  ((ue_mm_context->access_restriction_data & mASK) ? 'X' : 'O')
      bformata(
          bstr_dump,
          "    (O = allowed, X = !O) |UTRAN|GERAN|GAN|HSDPA EVO|E_UTRAN|HO TO "
          "NO "
          "3GPP|\n");
      bformata(
          bstr_dump,
          "    - Access restriction  |  %c  |  %c  | %c |    %c    |   %c   |  "
          "  "
          "  %c      |\n",
          DISPLAY_BIT_MASK_PRESENT(ARD_UTRAN_NOT_ALLOWED),
          DISPLAY_BIT_MASK_PRESENT(ARD_GERAN_NOT_ALLOWED),
          DISPLAY_BIT_MASK_PRESENT(ARD_GAN_NOT_ALLOWED),
          DISPLAY_BIT_MASK_PRESENT(ARD_I_HSDPA_EVO_NOT_ALLOWED),
          DISPLAY_BIT_MASK_PRESENT(ARD_E_UTRAN_NOT_ALLOWED),
          DISPLAY_BIT_MASK_PRESENT(ARD_HO_TO_NON_3GPP_NOT_ALLOWED));
      // TODO bformata (bstr_dump, "    - Access Mode ....: %s\n",
      // ACCESS_MODE_TO_STRING (ue_mm_context->access_mode));
      // TODO MSISDN
      // bformata (bstr_dump, "    - MSISDN .........: %s\n",
      // (ue_mm_context->msisdn) ? ue_mm_context->msisdn->data:"None");
      bformata(
          bstr_dump, "    - RAU/TAU timer ..: %u\n",
          ue_mm_context->rau_tau_timer);
      // TODO IMEISV
      // if (IS_EMM_CTXT_PRESENT_IMEISV(&ue_mm_context->nas_emm_context)) {
      //  bformata (bstr_dump, "    - IMEISV .........: %*s\n",
      //  IMEISV_DIGITS_MAX, ue_mm_context->nas_emm_context._imeisv);
      //}
      bformata(bstr_dump, "    - AMBR (bits/s)     ( Downlink |  Uplink  )\n");
      // TODO bformata (bstr_dump, "        Subscribed ...: (%010" PRIu64
      // "|%010" PRIu64 ")\n", ue_mm_context->subscribed_ambr.br_dl,
      // ue_mm_context->subscribed_ambr.br_ul);
      bformata(
          bstr_dump,
          "        Allocated ....: (%010" PRIu64 "|%010" PRIu64 ")\n",
          ue_mm_context->used_ambr.br_dl, ue_mm_context->used_ambr.br_ul);

      bformata(bstr_dump, "    - APN config list:\n");

      for (j = 0; j < ue_mm_context->apn_config_profile.nb_apns; j++) {
        struct apn_configuration_s* apn_config_p;

        apn_config_p = &ue_mm_context->apn_config_profile.apn_configuration[j];
        /*
         * Default APN ?
         */
        bformata(
            bstr_dump, "        - Default APN ...: %s\n",
            (apn_config_p->context_identifier ==
             ue_mm_context->apn_config_profile.context_identifier) ?
                "TRUE" :
                "FALSE");
        bformata(
            bstr_dump, "        - APN ...........: %s\n",
            apn_config_p->service_selection);
        bformata(
            bstr_dump, "        - AMBR (bits/s) ( Downlink |  Uplink  )\n");
        bformata(
            bstr_dump,
            "                        (%010" PRIu64 "|%010" PRIu64 ")\n",
            apn_config_p->ambr.br_dl, apn_config_p->ambr.br_ul);
        bformata(
            bstr_dump, "        - PDN type ......: %s\n",
            PDN_TYPE_TO_STRING(apn_config_p->pdn_type));
        bformata(bstr_dump, "        - QOS\n");
        bformata(
            bstr_dump, "            QCI .........: %u\n",
            apn_config_p->subscribed_qos.qci);
        bformata(
            bstr_dump, "            Prio level ..: %u\n",
            apn_config_p->subscribed_qos.allocation_retention_priority
                .priority_level);
        bformata(
            bstr_dump, "            Pre-emp vul .: %s\n",
            (apn_config_p->subscribed_qos.allocation_retention_priority
                 .pre_emp_vulnerability == PRE_EMPTION_VULNERABILITY_ENABLED) ?
                "ENABLED" :
                "DISABLED");
        bformata(
            bstr_dump, "            Pre-emp cap .: %s\n",
            (apn_config_p->subscribed_qos.allocation_retention_priority
                 .pre_emp_capability == PRE_EMPTION_CAPABILITY_ENABLED) ?
                "ENABLED" :
                "DISABLED");

        if (apn_config_p->nb_ip_address == 0) {
          bformata(
              bstr_dump, "            IP addr .....: Dynamic allocation\n");
        } else {
          int i;

          bformata(bstr_dump, "            IP addresses :\n");

          for (i = 0; i < apn_config_p->nb_ip_address; i++) {
            if (apn_config_p->ip_address[i].pdn_type == IPv4) {
              char ipv4[INET_ADDRSTRLEN];
              inet_ntop(
                  AF_INET,
                  (void*) &apn_config_p->ip_address[i].address.ipv4_address,
                  ipv4, INET_ADDRSTRLEN);
              bformata(bstr_dump, "                           [%s]\n", ipv4);
            } else {
              char ipv6[INET6_ADDRSTRLEN];
              inet_ntop(
                  AF_INET6, &apn_config_p->ip_address[i].address.ipv6_address,
                  ipv6, INET6_ADDRSTRLEN);
              bformata(bstr_dump, "                           [%s]\n", ipv6);
            }
          }
        }
        bformata(bstr_dump, "\n");
      }
      bformata(bstr_dump, "    - PDNs:\n");
      for (pdn_cid_t pdn_cid = 0; pdn_cid < MAX_APN_PER_UE; pdn_cid++) {
        pdn_context_t* pdn_context = ue_mm_context->pdn_contexts[pdn_cid];
        if (pdn_context) {
          mme_app_dump_pdn_context(
              ue_mm_context, pdn_context, pdn_cid, 8, bstr_dump);
        }
      }
    }
    bcatcstr(
        bstr_dump,
        "---------------------------------------------------------\n");
    OAILOG_DEBUG(LOG_MME_APP, "%s\n", bdata(bstr_dump));
    bdestroy_wrapper(&bstr_dump);
    return false;
  }
  return true;
}

//------------------------------------------------------------------------------
void mme_app_dump_ue_contexts()
//------------------------------------------------------------------------------
{
  hash_table_ts_t* mme_state_ue_id_ht = get_mme_ue_state();
  hashtable_ts_apply_callback_on_elements(
      mme_state_ue_id_ht, mme_app_dump_ue_context, NULL, NULL);
}

//------------------------------------------------------------------------------
void mme_app_handle_s1ap_ue_context_release_req(
    const itti_s1ap_ue_context_release_req_t* const s1ap_ue_context_release_req)

{
  mme_app_handle_s1ap_ue_context_release(
      s1ap_ue_context_release_req->mme_ue_s1ap_id,
      s1ap_ue_context_release_req->enb_ue_s1ap_id,
      s1ap_ue_context_release_req->enb_id,
      s1ap_ue_context_release_req->relCause);
}

void mme_app_handle_s1ap_ue_context_modification_fail(
    const itti_s1ap_ue_context_mod_resp_fail_t* const s1ap_ue_context_mod_fail)
//------------------------------------------------------------------------------
{
  struct ue_mm_context_s* ue_context_p = NULL;

  OAILOG_FUNC_IN(LOG_MME_APP);

  OAILOG_ERROR(
      LOG_MME_APP,
      " UE CONTEXT MODIFICATION FAILURE RECEIVED for UE-ID [%d] FAILURE_CAUSE "
      "[%ld]\n ",
      s1ap_ue_context_mod_fail->mme_ue_s1ap_id,
      s1ap_ue_context_mod_fail->cause);

  ue_context_p = mme_ue_context_exists_mme_ue_s1ap_id(
      s1ap_ue_context_mod_fail->mme_ue_s1ap_id);
  if (!ue_context_p) {
    OAILOG_ERROR(
        LOG_MME_APP,
        " UE CONTEXT MODIFICATION FAILURE RECEIVED, Failed to find UE context"
        "for mme_ue_s1ap_id 0x%06" PRIX32 " \n",
        s1ap_ue_context_mod_fail->mme_ue_s1ap_id);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  // Stop ue_context_modification  guard timer,if running
  if (ue_context_p->ue_context_modification_timer.id !=
      MME_APP_TIMER_INACTIVE_ID) {
    mme_app_stop_timer(ue_context_p->ue_context_modification_timer.id);
    ue_context_p->ue_context_modification_timer.id = MME_APP_TIMER_INACTIVE_ID;
  }
  if (ue_context_p->sgs_context != NULL) {
    handle_csfb_s1ap_procedure_failure(
        ue_context_p, "ue_context_modification_timer_expired",
        UE_CONTEXT_MODIFICATION_PROCEDURE_FAILED);
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

void mme_app_handle_s1ap_ue_context_modification_resp(
    const itti_s1ap_ue_context_mod_resp_t* const s1ap_ue_context_mod_resp)
//------------------------------------------------------------------------------
{
  struct ue_mm_context_s* ue_context_p = NULL;
  OAILOG_FUNC_IN(LOG_MME_APP);

  OAILOG_DEBUG(
      LOG_MME_APP,
      " UE CONTEXT MODIFICATION RESPONSE RECEIVED for UE-ID [%d] \n ",
      s1ap_ue_context_mod_resp->mme_ue_s1ap_id);

  ue_context_p = mme_ue_context_exists_mme_ue_s1ap_id(
      s1ap_ue_context_mod_resp->mme_ue_s1ap_id);
  if (!ue_context_p) {
    OAILOG_ERROR(
        LOG_MME_APP,
        " UE CONTEXT MODIFICATION RESPONSE RECEIVED, Failed to find UE context"
        "for mme_ue_s1ap_id 0x%06" PRIX32 " \n",
        s1ap_ue_context_mod_resp->mme_ue_s1ap_id);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }

  // Stop ue_context_modification  guard timer,if running
  if (ue_context_p->ue_context_modification_timer.id !=
      MME_APP_TIMER_INACTIVE_ID) {
    mme_app_stop_timer(ue_context_p->ue_context_modification_timer.id);
    ue_context_p->ue_context_modification_timer.id = MME_APP_TIMER_INACTIVE_ID;
  }

  OAILOG_FUNC_OUT(LOG_MME_APP);
}
//------------------------------------------------------------------------------
void mme_app_handle_enb_deregister_ind(
    const itti_s1ap_eNB_deregistered_ind_t* const eNB_deregistered_ind) {
  for (int i = 0; i < eNB_deregistered_ind->nb_ue_to_deregister; i++) {
    mme_app_handle_s1ap_ue_context_release(
        eNB_deregistered_ind->mme_ue_s1ap_id[i],
        eNB_deregistered_ind->enb_ue_s1ap_id[i], eNB_deregistered_ind->enb_id,
        S1AP_SCTP_SHUTDOWN_OR_RESET);
  }
}

//------------------------------------------------------------------------------
void mme_app_handle_enb_reset_req(
    const itti_s1ap_enb_initiated_reset_req_t* const enb_reset_req) {
  MessageDef* msg;
  itti_s1ap_enb_initiated_reset_ack_t* reset_ack;

  OAILOG_INFO(
      LOG_MME_APP,
      " eNB Reset request received. eNB id = %d, reset_type  %d \n ",
      enb_reset_req->enb_id, enb_reset_req->s1ap_reset_type);
  if (enb_reset_req->ue_to_reset_list == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP, "Invalid UE list received in eNB Reset Request\n");
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }

  for (int i = 0; i < enb_reset_req->num_ue; i++) {
    mme_app_handle_s1ap_ue_context_release(
        enb_reset_req->ue_to_reset_list[i].mme_ue_s1ap_id,
        enb_reset_req->ue_to_reset_list[i].enb_ue_s1ap_id,
        enb_reset_req->enb_id, S1AP_SCTP_SHUTDOWN_OR_RESET);
  }

  // Send Reset Ack to S1AP module
  msg = itti_alloc_new_message(TASK_MME_APP, S1AP_ENB_INITIATED_RESET_ACK);
  reset_ack = &S1AP_ENB_INITIATED_RESET_ACK(msg);

  // ue_to_reset_list needs to be freed by S1AP module
  reset_ack->ue_to_reset_list = enb_reset_req->ue_to_reset_list;
  reset_ack->s1ap_reset_type  = enb_reset_req->s1ap_reset_type;
  reset_ack->sctp_assoc_id    = enb_reset_req->sctp_assoc_id;
  reset_ack->sctp_stream_id   = enb_reset_req->sctp_stream_id;
  reset_ack->num_ue           = enb_reset_req->num_ue;

  send_msg_to_task(&mme_app_task_zmq_ctx, TASK_S1AP, msg);

  OAILOG_INFO(
      LOG_MME_APP, " Reset Ack sent to S1AP. eNB id = %d, reset_type  %d \n ",
      enb_reset_req->enb_id, enb_reset_req->s1ap_reset_type);

  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
/*
   From GPP TS 23.401 version 11.11.0 Release 11, section 5.3.5 S1 release
   procedure, point 6: The MME deletes any eNodeB related information ("eNodeB
   Address in Use for S1-MME" and "eNB UE S1AP ID") from the UE's MME context,
   but, retains the rest of the UE's MME context including the S-GW's S1-U
   configuration information (address and TEIDs). All non-GBR EPS bearers
   established for the UE are preserved in the MME and in the Serving GW. If the
   cause of S1 release is because of User is inactivity, Inter-RAT Redirection,
   the MME shall preserve the GBR bearers. If the cause of S1 release is because
   of CS Fallback triggered, further details about bearer handling are described
   in TS 23.272 [58]. Otherwise, e.g. Radio Connection With UE Lost, S1
   signalling connection lost, eNodeB failure the MME shall trigger the MME
   Initiated Dedicated Bearer Deactivation procedure (clause 5.4.4.2) for the
   GBR bearer(s) of the UE after the S1 Release procedure is completed.
*/
//------------------------------------------------------------------------------
void mme_app_handle_s1ap_ue_context_release_complete(
    mme_app_desc_t* mme_app_desc_p,
    const itti_s1ap_ue_context_release_complete_t* const
        s1ap_ue_context_release_complete)
//------------------------------------------------------------------------------
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  struct ue_mm_context_s* ue_context_p = NULL;

  ue_context_p = mme_ue_context_exists_mme_ue_s1ap_id(
      s1ap_ue_context_release_complete->mme_ue_s1ap_id);

  OAILOG_INFO(
      LOG_MME_APP,
      "Received UE context release complete message for "
      "ue_id: " MME_UE_S1AP_ID_FMT,
      s1ap_ue_context_release_complete->mme_ue_s1ap_id);

  if (!ue_context_p) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "UE context doesn't exist for enb_ue_s1ap_ue_id " ENB_UE_S1AP_ID_FMT
        " mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT "\n",
        s1ap_ue_context_release_complete->enb_ue_s1ap_id,
        s1ap_ue_context_release_complete->mme_ue_s1ap_id);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }

  mme_app_delete_s11_procedure_create_bearer(ue_context_p);

  if (ue_context_p->mm_state == UE_UNREGISTERED) {
    if (ue_context_p->nb_active_pdn_contexts == 0) {
      // No Session
      OAILOG_DEBUG_UE(
          LOG_MME_APP, ue_context_p->emm_context._imsi64,
          "Deleting UE context associated in MME for "
          "mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT "\n ",
          s1ap_ue_context_release_complete->mme_ue_s1ap_id);

      // Send PUR,before removal of ue contexts
      if ((ue_context_p->send_ue_purge_request == true) &&
          (ue_context_p->hss_initiated_detach == false)) {
        mme_app_send_s6a_purge_ue_req(mme_app_desc_p, ue_context_p);
      }
      mme_remove_ue_context(&mme_app_desc_p->mme_ue_contexts, ue_context_p);
      update_mme_app_stats_connected_ue_sub();
      OAILOG_FUNC_OUT(LOG_MME_APP);
    } else {
      // Send a DELETE_SESSION_REQUEST message to the SGW
      for (pdn_cid_t i = 0; i < MAX_APN_PER_UE; i++) {
        if (ue_context_p->pdn_contexts[i]) {
          // Send a DELETE_SESSION_REQUEST message to the SGW
          mme_app_send_delete_session_request(
              ue_context_p, ue_context_p->pdn_contexts[i]->default_ebi, i);
        }
      }
      // Move the UE to Idle state
      mme_ue_context_update_ue_sig_connection_state(
          &mme_app_desc_p->mme_ue_contexts, ue_context_p, ECM_IDLE);
    }
  } else {
    // Update keys and ECM state
    mme_ue_context_update_ue_sig_connection_state(
        &mme_app_desc_p->mme_ue_contexts, ue_context_p, ECM_IDLE);
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//-------------------------------------------------------------------------------------------------------
void mme_ue_context_update_ue_emm_state(
    mme_ue_s1ap_id_t mme_ue_s1ap_id, mm_state_t new_mm_state) {
  // Function is used to update UE's mobility management State-
  // Registered/Un-Registered

  struct ue_mm_context_s* ue_context_p = NULL;

  OAILOG_FUNC_IN(LOG_MME_APP);
  ue_context_p = mme_ue_context_exists_mme_ue_s1ap_id(mme_ue_s1ap_id);
  if (ue_context_p == NULL) {
    OAILOG_CRITICAL(LOG_MME_APP, "**** Abnormal- UE context is null.****\n");
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  if ((ue_context_p->mm_state == UE_UNREGISTERED) &&
      (new_mm_state == UE_REGISTERED)) {
    ue_context_p->mm_state = new_mm_state;

    // Report directoryd UE record
    directoryd_report_location_a(
        ue_context_p->emm_context._imsi64,
        ue_context_p->emm_context._imsi.length);

    // Update Stats
    update_mme_app_stats_attached_ue_add();
    OAILOG_INFO_UE(
        LOG_MME_APP, ue_context_p->emm_context._imsi64,
        "UE STATE - REGISTERED.\n");
  } else if (
      (ue_context_p->mm_state == UE_REGISTERED) &&
      (new_mm_state == UE_UNREGISTERED)) {
    ue_context_p->mm_state = new_mm_state;

    // Update Stats
    update_mme_app_stats_attached_ue_sub();
    OAILOG_INFO_UE(
        LOG_MME_APP, ue_context_p->emm_context._imsi64,
        "UE STATE - UNREGISTERED.\n");
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
static void mme_app_handle_s1ap_ue_context_release(
    const mme_ue_s1ap_id_t mme_ue_s1ap_id,
    const enb_ue_s1ap_id_t enb_ue_s1ap_id, uint32_t enb_id, enum s1cause cause)
//------------------------------------------------------------------------------
{
  struct ue_mm_context_s* ue_mm_context = NULL;
  enb_s1ap_id_key_t enb_s1ap_id_key     = INVALID_ENB_UE_S1AP_ID_KEY;

  OAILOG_FUNC_IN(LOG_MME_APP);
  mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
  ue_mm_context = mme_ue_context_exists_mme_ue_s1ap_id(mme_ue_s1ap_id);
  if (!ue_mm_context) {
    /*
     * Use enb_ue_s1ap_id_key to get the UE context - In case MME APP could not
     * update S1AP with valid mme_ue_s1ap_id before context release is triggered
     * from s1ap.
     */
    MME_APP_ENB_S1AP_ID_KEY(enb_s1ap_id_key, enb_id, enb_ue_s1ap_id);
    ue_mm_context = mme_ue_context_exists_enb_ue_s1ap_id(
        &mme_app_desc_p->mme_ue_contexts, enb_s1ap_id_key);

    OAILOG_WARNING(
        LOG_MME_APP,
        "Invalid mme_ue_s1ap_ue_id " MME_UE_S1AP_ID_FMT
        " received from S1AP. Using enb_s1ap_id_key %ld to get the context \n",
        mme_ue_s1ap_id, enb_s1ap_id_key);
  }
  if (!ue_mm_context) {
    OAILOG_ERROR(
        LOG_MME_APP,
        " UE Context Release Req: UE context doesn't exist for "
        "enb_ue_s1ap_ue_id " ENB_UE_S1AP_ID_FMT
        " mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT "\n",
        enb_ue_s1ap_id, mme_ue_s1ap_id);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  // Set the UE context release cause in UE context. This is used while
  // constructing UE Context Release Command
  ue_mm_context->ue_context_rel_cause = cause;

  if (ue_mm_context->ecm_state == ECM_IDLE) {
    // This case could happen during sctp reset, before the UE could move to
    // ECM_CONNECTED calling below function to set the enb_s1ap_id_key to
    // invalid
    if (ue_mm_context->ue_context_rel_cause == S1AP_SCTP_SHUTDOWN_OR_RESET) {
      mme_ue_context_update_ue_sig_connection_state(
          &mme_app_desc_p->mme_ue_contexts, ue_mm_context, ECM_IDLE);
      mme_app_itti_ue_context_release(
          ue_mm_context, ue_mm_context->ue_context_rel_cause);
      OAILOG_WARNING_UE(
          LOG_MME_APP, ue_mm_context->emm_context._imsi64,
          "UE Conetext Release Reqeust:Cause SCTP RESET/SHUTDOWN. UE state: "
          "IDLE. mme_ue_s1ap_id = %d, enb_ue_s1ap_id = %d Action -- Handle the "
          "message\n ",
          ue_mm_context->mme_ue_s1ap_id, ue_mm_context->enb_ue_s1ap_id);
      OAILOG_FUNC_OUT(LOG_MME_APP);
    }
    OAILOG_ERROR_UE(
        LOG_MME_APP, ue_mm_context->emm_context._imsi64,
        "ERROR: UE Context Release Request: UE state : IDLE. "
        "enb_ue_s1ap_ue_id " ENB_UE_S1AP_ID_FMT
        " mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT " Action--- Ignore the message\n",
        ue_mm_context->enb_ue_s1ap_id, ue_mm_context->mme_ue_s1ap_id);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  } else {
    // This case could happen during sctp reset, while attach procedure is
    // ongoing and ue is in ECM_CONNECTED calling below function to set the
    // enb_s1ap_id_key to invalid
    if (ue_mm_context->ue_context_rel_cause == S1AP_SCTP_SHUTDOWN_OR_RESET) {
      // Update keys and ECM state
      mme_ue_context_update_ue_sig_connection_state(
          &mme_app_desc_p->mme_ue_contexts, ue_mm_context, ECM_IDLE);
      OAILOG_WARNING_UE(
          LOG_MME_APP, ue_mm_context->emm_context._imsi64,
          "SCTP RESET/SHUTDOWN. UE state: CONNECTED. mme_ue_s1ap_id = %d, "
          "enb_ue_s1ap_id = %d"
          " Action -- Handle the message\n ",
          ue_mm_context->mme_ue_s1ap_id, ue_mm_context->enb_ue_s1ap_id);
    }
  }

  // Stop Initial context setup process guard timer,if running
  if (ue_mm_context->initial_context_setup_rsp_timer.id !=
      MME_APP_TIMER_INACTIVE_ID) {
    mme_app_stop_timer(ue_mm_context->initial_context_setup_rsp_timer.id);
    ue_mm_context->initial_context_setup_rsp_timer.id =
        MME_APP_TIMER_INACTIVE_ID;
    // Setting UE context release cause as Initial context setup failure
    ue_mm_context->ue_context_rel_cause = S1AP_INITIAL_CONTEXT_SETUP_FAILED;
  }
  // Stop UE context modification process guard timer,if running
  if (ue_mm_context->ue_context_modification_timer.id !=
      MME_APP_TIMER_INACTIVE_ID) {
    mme_app_stop_timer(ue_mm_context->ue_context_modification_timer.id);
    ue_mm_context->ue_context_modification_timer.id = MME_APP_TIMER_INACTIVE_ID;
  }
  // Stop esm timer if running
  for (int bid = 0; bid < BEARERS_PER_UE; bid++) {
    if (ue_mm_context->bearer_contexts[bid]) {
      esm_ebr_stop_timer(
          &ue_mm_context->emm_context,
          ue_mm_context->bearer_contexts[bid]->ebi);
    }
  }

  if (ue_mm_context->mm_state == UE_UNREGISTERED) {
    // Initiate Implicit Detach for the UE
    OAILOG_ERROR_UE(
        LOG_MME_APP, ue_mm_context->emm_context._imsi64,
        "UE context release request received while UE is in Deregistered state "
        "Perform implicit detach for ue-id" MME_UE_S1AP_ID_FMT "\n",
        ue_mm_context->mme_ue_s1ap_id);
    if (!ue_mm_context->emm_context.new_attach_info) {
      nas_proc_implicit_detach_ue_ind(ue_mm_context->mme_ue_s1ap_id);
    }
  } else {
    if (cause == S1AP_NAS_UE_NOT_AVAILABLE_FOR_PS) {
      for (pdn_cid_t i = 0; i < MAX_APN_PER_UE; i++) {
        if (ue_mm_context->pdn_contexts[i]) {
          if ((mme_app_send_s11_suspend_notification(ue_mm_context, i)) !=
              RETURNok) {
            OAILOG_ERROR_UE(
                LOG_MME_APP, ue_mm_context->emm_context._imsi64,
                "Failed to send S11 Suspend Notification for imsi\n");
          }
        }
      }
    } else {
      // release S1-U tunnel mapping in S_GW for all the active bearers for the
      // UE
      for (pdn_cid_t i = 0; i < MAX_APN_PER_UE; i++) {
        if (ue_mm_context->pdn_contexts[i]) {
          mme_app_send_s11_release_access_bearers_req(ue_mm_context, i);
        }
      }
    }
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

bool is_mme_ue_context_network_access_mode_packet_only(
    ue_mm_context_t* ue_context_p) {
  // Function is used to check the UE's Network Access Mode received in ULA from
  // HSS

  OAILOG_FUNC_IN(LOG_MME_APP);
  if (ue_context_p == NULL) {
    OAILOG_CRITICAL(LOG_MME_APP, "**** Abnormal- UE context is null.****\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  if (ue_context_p->network_access_mode == NAM_ONLY_PACKET) {
    OAILOG_FUNC_RETURN(LOG_MME_APP, true);
  } else {
    OAILOG_FUNC_RETURN(LOG_MME_APP, false);
  }
}

//-------------------------------------------------------------------------------------------------------
void mme_ue_context_update_ue_sgs_vlr_reliable(
    mme_ue_s1ap_id_t mme_ue_s1ap_id, bool vlr_reliable) {
  // Function is used to update the UE's SGS vlr reliable flag - true/false

  struct ue_mm_context_s* ue_context_p = NULL;

  OAILOG_FUNC_IN(LOG_MME_APP);
  ue_context_p = mme_ue_context_exists_mme_ue_s1ap_id(mme_ue_s1ap_id);
  if (ue_context_p == NULL) {
    OAILOG_CRITICAL(LOG_MME_APP, "**** Abnormal- UE context is null.****\n");
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  if ((ue_context_p->sgs_context) &&
      (ue_context_p->sgs_context->vlr_reliable != vlr_reliable)) {
    ue_context_p->sgs_context->vlr_reliable = vlr_reliable;
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//-------------------------------------------------------------------------------------------------------
bool mme_ue_context_get_ue_sgs_vlr_reliable(mme_ue_s1ap_id_t mme_ue_s1ap_id) {
  // Function is used to get the UE's SGS vlr reliable flag - true/false

  struct ue_mm_context_s* ue_context_p = NULL;
  bool vlr_reliable                    = false;

  OAILOG_FUNC_IN(LOG_MME_APP);
  ue_context_p = mme_ue_context_exists_mme_ue_s1ap_id(mme_ue_s1ap_id);
  if (ue_context_p == NULL) {
    OAILOG_CRITICAL(LOG_MME_APP, "**** Abnormal- UE context is null.****\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  if ((ue_context_p->sgs_context) &&
      (ue_context_p->sgs_context->vlr_reliable == true)) {
    vlr_reliable = true;
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, vlr_reliable);
}

//-------------------------------------------------------------------------------------------------------
void mme_ue_context_update_ue_sgs_neaf(
    mme_ue_s1ap_id_t mme_ue_s1ap_id, bool neaf) {
  // Function is used to update the UE's SGS neaf flag - true/false

  struct ue_mm_context_s* ue_context_p = NULL;

  OAILOG_FUNC_IN(LOG_MME_APP);
  ue_context_p = mme_ue_context_exists_mme_ue_s1ap_id(mme_ue_s1ap_id);
  if (ue_context_p == NULL) {
    OAILOG_CRITICAL(LOG_MME_APP, "**** Abnormal- UE context is null.****\n");
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  if ((ue_context_p->sgs_context) &&
      (ue_context_p->sgs_context->neaf != neaf)) {
    ue_context_p->sgs_context->neaf = neaf;
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//-------------------------------------------------------------------------------------------------------
bool mme_ue_context_get_ue_sgs_neaf(mme_ue_s1ap_id_t mme_ue_s1ap_id) {
  // Function is used to get the UE's SGS neaf flag - true/false
  struct ue_mm_context_s* ue_context_p = NULL;

  OAILOG_FUNC_IN(LOG_MME_APP);
  ue_context_p = mme_ue_context_exists_mme_ue_s1ap_id(mme_ue_s1ap_id);
  if (ue_context_p == NULL) {
    OAILOG_CRITICAL(LOG_MME_APP, "**** Abnormal- UE context is null.****\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  if ((ue_context_p->sgs_context) &&
      (ue_context_p->sgs_context->neaf == true)) {
    OAILOG_ERROR_UE(
        LOG_MME_APP, ue_context_p->emm_context._imsi64,
        "In MME APP NEAF is set to True\n");
    return true;
  } else {
    return false;
  }
}

void mme_app_recover_timers_for_all_ues(void) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  hash_table_ts_t* mme_state_imsi_ht = get_mme_ue_state();
  uint32_t num_unreg_ues             = 0;
  hash_key_t* mme_ue_id_unreg_list;
  mme_ue_id_unreg_list =
      (hash_key_t*) calloc(mme_state_imsi_ht->num_elements, sizeof(hash_key_t));
  hashtable_ts_apply_callback_on_elements(
      mme_state_imsi_ht, mme_app_recover_timers_for_ue, &num_unreg_ues,
      (void**) &mme_ue_id_unreg_list);

  // Handle timer for unregistered UEs here as it will modify the hashtable
  // entries
  struct ue_mm_context_s* ue_context_p = NULL;
  for (uint32_t i = 0; i < num_unreg_ues; i++) {
    ue_context_p =
        mme_ue_context_exists_mme_ue_s1ap_id(mme_ue_id_unreg_list[i]);
    mme_app_handle_timer_for_unregistered_ue(ue_context_p);
  }

  free(mme_ue_id_unreg_list);
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

static bool mme_app_recover_timers_for_ue(
    const hash_key_t keyP, void* const ue_context_pP, void* param_pP,
    void** result_pP) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  uint32_t* num_unreg_ues  = (uint32_t*) param_pP;
  hash_key_t** mme_id_list = (hash_key_t**) result_pP;
  struct ue_mm_context_s* const ue_mm_context_pP =
      (struct ue_mm_context_s*) ue_context_pP;

  if (!ue_mm_context_pP) {
    OAILOG_ERROR(LOG_MME_APP, "UE context is NULL\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, false);
  }

  if (ue_mm_context_pP->time_mobile_reachability_timer_started) {
    mme_app_resume_timer(
        ue_mm_context_pP,
        ue_mm_context_pP->time_mobile_reachability_timer_started,
        &ue_mm_context_pP->mobile_reachability_timer,
        mme_app_handle_mobile_reachability_timer_expiry, "Mobile Reachability");
  }
  if (ue_mm_context_pP->time_implicit_detach_timer_started) {
    mme_app_resume_timer(
        ue_mm_context_pP, ue_mm_context_pP->time_implicit_detach_timer_started,
        &ue_mm_context_pP->implicit_detach_timer,
        mme_app_handle_implicit_detach_timer_expiry, "Implicit Detach");
  }
  if (ue_mm_context_pP->time_paging_response_timer_started) {
    mme_app_resume_timer(
        ue_mm_context_pP, ue_mm_context_pP->time_paging_response_timer_started,
        &ue_mm_context_pP->paging_response_timer,
        mme_app_handle_paging_timer_expiry, "Paging Response");
  }
  if (ue_mm_context_pP->emm_context._emm_fsm_state == EMM_REGISTERED &&
      ue_mm_context_pP->time_ics_rsp_timer_started) {
    mme_app_resume_timer(
        ue_mm_context_pP, ue_mm_context_pP->time_ics_rsp_timer_started,
        &ue_mm_context_pP->initial_context_setup_rsp_timer,
        mme_app_handle_initial_context_setup_rsp_timer_expiry,
        "Initial Context Setup Response");
  }
  if (ue_mm_context_pP->emm_context._emm_fsm_state == EMM_REGISTERED) {
    mme_app_resume_esm_ebr_timer(ue_mm_context_pP);
  }

  // timer for network initiated detach procedure
  if (ue_mm_context_pP->emm_context.t3422_arg) {
    detach_t3422_handler(
        (void*) ue_mm_context_pP->emm_context.t3422_arg,
        &ue_mm_context_pP->emm_context._imsi64);
  }

  if (ue_mm_context_pP->emm_context._emm_fsm_state != EMM_REGISTERED) {
    (*mme_id_list)[*num_unreg_ues] = keyP;
    ++(*num_unreg_ues);
    OAILOG_DEBUG(LOG_MME_APP, "Added %u unreg UEs", *num_unreg_ues);
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, false);
}

static void mme_app_resume_esm_ebr_timer(ue_mm_context_t* ue_context_p) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  for (int idx = 0; idx < BEARERS_PER_UE; idx++) {
    bearer_context_t* bearer_context_p = ue_context_p->bearer_contexts[idx];
    if (bearer_context_p) {
      bearer_context_p->esm_ebr_context.timer.id = NAS_TIMER_INACTIVE_ID;
      pdn_cid_t pdn_cid                          = bearer_context_p->pdn_cx_id;
      // Below check is added to identify default and dedicated bearer
      if (ue_context_p->pdn_contexts[pdn_cid] &&
          (ue_context_p->pdn_contexts[pdn_cid]->default_ebi ==
           bearer_context_p->ebi)) {
        // Invoke callback registered for default bearer's activation
        if ((bearer_context_p->esm_ebr_context.args) &&
            (bearer_context_p->esm_ebr_context.status ==
             ESM_EBR_ACTIVE_PENDING)) {
          bearer_context_p->esm_ebr_context.timer.sec = T3485_DEFAULT_VALUE;
          default_eps_bearer_activate_t3485_handler(
              bearer_context_p->esm_ebr_context.args,
              &ue_context_p->emm_context._imsi64);
        } else {  // Invoke callback registered for default bearer's
                  // deactivation procedure
          if ((ue_context_p->bearer_contexts[idx]->esm_ebr_context.args) &&
              (ue_context_p->bearer_contexts[idx]->esm_ebr_context.status ==
               ESM_EBR_INACTIVE_PENDING)) {
            eps_bearer_deactivate_t3495_handler(
                ue_context_p->bearer_contexts[idx]->esm_ebr_context.args,
                &ue_context_p->emm_context._imsi64);
          }
        }
      } else {
        // Invoke callback registered for dedicated bearer's activation
        if ((bearer_context_p->esm_ebr_context.args) &&
            (bearer_context_p->esm_ebr_context.status ==
             ESM_EBR_ACTIVE_PENDING)) {
          bearer_context_p->esm_ebr_context.timer.sec = T3485_DEFAULT_VALUE;
          dedicated_eps_bearer_activate_t3485_handler(
              bearer_context_p->esm_ebr_context.args,
              &ue_context_p->emm_context._imsi64);
        } else {  // Invoke callback registered for dedicated bearer's
                  // deactivation procedure
          if ((ue_context_p->bearer_contexts[idx]->esm_ebr_context.args) &&
              (ue_context_p->bearer_contexts[idx]->esm_ebr_context.status ==
               ESM_EBR_INACTIVE_PENDING)) {
            eps_bearer_deactivate_t3495_handler(
                ue_context_p->bearer_contexts[idx]->esm_ebr_context.args,
                &ue_context_p->emm_context._imsi64);
          }
        }
      }
    }
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

static void mme_app_handle_timer_for_unregistered_ue(
    ue_mm_context_t* ue_context_p) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  nas_emm_attach_proc_t* attach_proc =
      get_nas_specific_procedure_attach(&ue_context_p->emm_context);
  if ((attach_proc) && (is_nas_attach_accept_sent(attach_proc))) {
    OAILOG_INFO_UE(
        LOG_MME_APP, ue_context_p->emm_context._imsi64,
        "Send nw initiated detach req for ue_id with re-attach "
        "required" MME_UE_S1AP_ID_FMT "\n",
        ue_context_p->mme_ue_s1ap_id);
    emm_proc_nw_initiated_detach_request(
        ue_context_p->mme_ue_s1ap_id, NW_DETACH_TYPE_RE_ATTACH_REQUIRED);
  } else {
    OAILOG_INFO_UE(
        LOG_MME_APP, ue_context_p->emm_context._imsi64,
        "Initiate implicit detach for ue_id" MME_UE_S1AP_ID_FMT "\n",
        ue_context_p->mme_ue_s1ap_id);
    nas_proc_implicit_detach_ue_ind(ue_context_p->mme_ue_s1ap_id);
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

void mme_app_remove_stale_ue_context(
    mme_app_desc_t* mme_app_desc_p,
    itti_s1ap_remove_stale_ue_context_t* s1ap_remove_stale_ue_context) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  enb_s1ap_id_key_t enb_s1ap_id_key = INVALID_ENB_UE_S1AP_ID_KEY;
  MME_APP_ENB_S1AP_ID_KEY(
      enb_s1ap_id_key, s1ap_remove_stale_ue_context->enb_id,
      s1ap_remove_stale_ue_context->enb_ue_s1ap_id);
  if (enb_s1ap_id_key == INVALID_ENB_UE_S1AP_ID_KEY) {
    OAILOG_ERROR(
        LOG_MME_APP, "Received invalid enb_s1ap_id_key :%lx\n",
        enb_s1ap_id_key);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  uint64_t mme_ue_s1ap_id = INVALID_MME_UE_S1AP_ID;
  if (hashtable_uint64_ts_get(
          mme_app_desc_p->mme_ue_contexts.enb_ue_s1ap_id_ue_context_htbl,
          (const hash_key_t) enb_s1ap_id_key,
          &mme_ue_s1ap_id) == HASH_TABLE_OK) {
    ue_mm_context_t* ue_context_p =
        mme_ue_context_exists_mme_ue_s1ap_id(mme_ue_s1ap_id);
    if (!ue_context_p) {
      hashtable_uint64_ts_remove(
          mme_app_desc_p->mme_ue_contexts.enb_ue_s1ap_id_ue_context_htbl,
          (const hash_key_t) enb_s1ap_id_key);
      OAILOG_INFO(
          LOG_MME_APP,
          "Removed stale UE context for mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT
          "enb_ue_s1ap_id " ENB_UE_S1AP_ID_FMT " enb_id: %d \n ",
          (mme_ue_s1ap_id_t) mme_ue_s1ap_id,
          s1ap_remove_stale_ue_context->enb_ue_s1ap_id,
          s1ap_remove_stale_ue_context->enb_id);
    }
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}
