/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#ifdef __cplusplus
extern "C" {
#endif
#include "log.h"
#include "intertask_interface.h"
#include "directoryd.h"
#include "conversions.h"
#ifdef __cplusplus
}
#endif
#include "common_defs.h"
#include "dynamic_memory_check.h"
#include "amf_app_state_manager.h"

namespace magma5g {
extern task_zmq_ctx_t amf_app_task_zmq_ctx;
ue_m5gmm_context_s ue_m5gmm_global_context;  // TODO : Upcoming PR has MAP
                                             // implementation, will be removed
static void amf_directoryd_report_location(uint64_t imsi, uint8_t imsi_len) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi, imsi_str, imsi_len);
  directoryd_report_location(imsi_str);
}

amf_ue_ngap_id_t amf_app_ctx_get_new_ue_id(
    amf_ue_ngap_id_t* amf_app_ue_ngap_id_generator_p) {
  amf_ue_ngap_id_t tmp = 0;
  tmp = __sync_fetch_and_add(amf_app_ue_ngap_id_generator_p, 1);
  return tmp;
}

/****************************************************************************
 **                                                                        **
 ** Name:    notify_ngap_new_ue_amf_ngap_id_association()                  **
 **                                                                        **
 ** Description: Sends AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION to NGAP         **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
void notify_ngap_new_ue_amf_ngap_id_association(
    const ue_m5gmm_context_s* ue_context_p) {
  MessageDef* message_p                                      = NULL;
  itti_amf_app_ngap_amf_ue_id_notification_t* notification_p = NULL;

  OAILOG_FUNC_IN(LOG_AMF_APP);
  if (ue_context_p == NULL) {
    OAILOG_ERROR(LOG_AMF_APP, " NULL UE context!\n");
    OAILOG_FUNC_OUT(LOG_AMF_APP);
  }
  message_p =
      itti_alloc_new_message(TASK_AMF_APP, AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION);
  notification_p = &message_p->ittiMsg.amf_app_ngap_amf_ue_id_notification;
  memset(notification_p, 0, sizeof(itti_amf_app_ngap_amf_ue_id_notification_t));
  notification_p->gnb_ue_ngap_id = ue_context_p->gnb_ue_ngap_id;
  notification_p->amf_ue_ngap_id = ue_context_p->amf_ue_ngap_id;
  notification_p->sctp_assoc_id  = ue_context_p->sctp_assoc_id_key;

  OAILOG_INFO_UE(
      LOG_AMF_APP, ue_context_p->amf_context.imsi64,
      " Sent AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION to NGAP for (ue_id = %u)\n",
      notification_p->amf_ue_ngap_id);

  send_msg_to_task(&amf_app_task_zmq_ctx, TASK_NGAP, message_p);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

/****************************************************************************
 **                                                                        **
 ** Name:        amf_insert_ue_context()                                   **
 **                                                                        **
 ** Description: Registers the UE context                                  **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
int amf_insert_ue_context(
    const amf_ue_context_t* amf_ue_context_p,
    const ue_m5gmm_context_s* ue_context_p) {
  hashtable_rc_t h_rc = HASH_TABLE_OK;

  OAILOG_FUNC_IN(LOG_AMF_APP);
  if (amf_ue_context_p == NULL) {
    OAILOG_ERROR(LOG_AMF_APP, "Invalid AMF UE context received\n");
    OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
  }
  if (ue_context_p == NULL) {
    OAILOG_ERROR(LOG_AMF_APP, "Invalid UE context received\n");
    OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
  }
  if (ue_context_p->gnb_ngap_id_key == INVALID_GNB_UE_NGAP_ID_KEY) {
    OAILOG_ERROR(LOG_AMF_APP, "Invalid gnb_ngap_id_key received\n");
    OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
  }

  // filled GNB UE NGAP ID
  h_rc = hashtable_uint64_ts_is_key_exists(
      amf_ue_context_p->gnb_ue_ngap_id_ue_context_htbl,
      (const hash_key_t) ue_context_p->gnb_ngap_id_key);

  h_rc = hashtable_uint64_ts_insert(
      amf_ue_context_p->gnb_ue_ngap_id_ue_context_htbl,
      (const hash_key_t) ue_context_p->gnb_ngap_id_key,
      ue_context_p->amf_ue_ngap_id);

  if (INVALID_AMF_UE_NGAP_ID != ue_context_p->amf_ue_ngap_id) {
    // filled IMSI
    if (ue_context_p->amf_context.imsi64) {
      amf_directoryd_report_location(
          ue_context_p->amf_context.imsi64,
          ue_context_p->amf_context.imsi.length);
    }

    // filled guti
    if ((0 != ue_context_p->amf_context.m5_guti.guamfi.amf_set_id) ||
        (0 != ue_context_p->amf_context.m5_guti.guamfi.amf_regionid) ||
        (0 != ue_context_p->amf_context.m5_guti.m_tmsi) ||
        (0 != ue_context_p->amf_context.m5_guti.guamfi.plmn
                  .mcc_digit1) ||  // MCC 000 does not exist in ITU table
        (0 != ue_context_p->amf_context.m5_guti.guamfi.plmn.mcc_digit2) ||
        (0 != ue_context_p->amf_context.m5_guti.guamfi.plmn.mcc_digit3)) {
      h_rc = obj_hashtable_uint64_ts_insert(
          amf_ue_context_p->guti_ue_context_htbl,
          (const void* const) & ue_context_p->amf_context.m5_guti,
          sizeof(ue_context_p->amf_context.m5_guti),
          ue_context_p->amf_ue_ngap_id);

      if (HASH_TABLE_OK != h_rc) {
        OAILOG_WARNING(
            LOG_AMF_APP,
            "Error could not register the ue context\n"
            " gnb_ue_ngap_id " GNB_UE_NGAP_ID_FMT
            " amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT,
            ue_context_p->gnb_ue_ngap_id, ue_context_p->amf_ue_ngap_id);
        OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
      }
    }
  }
  ue_m5gmm_global_context =
      *ue_context_p;  // TODO This has been taken care in new PR with support
                      // for MultiUE feature
  OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNok);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_create_new_ue_context()                                   **
 **                                                                        **
 ** Description: Creates new UE context                                    **
 **                                                                        **
 ***************************************************************************/
// warning: lock the UE context
ue_m5gmm_context_s* amf_create_new_ue_context(void) {
  ue_m5gmm_context_s* new_p = NULL;  // new (ue_m5gmm_context_s);

  if (!new_p) {
    OAILOG_ERROR(LOG_AMF_APP, "Failed to allocate memory for UE context \n");
    return NULL;
  }

  new_p->amf_ue_ngap_id  = INVALID_AMF_UE_NGAP_ID;
  new_p->gnb_ngap_id_key = INVALID_GNB_UE_NGAP_ID_KEY;
  new_p->gnb_ue_ngap_id  = INVALID_GNB_UE_NGAP_ID;

  // Initialize timers to INVALID IDs
  new_p->m5_mobile_reachability_timer.id    = AMF_APP_TIMER_INACTIVE_ID;
  new_p->m5_implicit_detach_timer.id        = AMF_APP_TIMER_INACTIVE_ID;
  new_p->m5_initial_context_setup_rsp_timer = (amf_app_timer_t){
      AMF_APP_TIMER_INACTIVE_ID, AMF_APP_INITIAL_CONTEXT_SETUP_RSP_TIMER_VALUE};
  new_p->m5_ulr_response_timer = (amf_app_timer_t){
      AMF_APP_TIMER_INACTIVE_ID, AMF_APP_ULR_RESPONSE_TIMER_VALUE};
  new_p->m5_ue_context_modification_timer = (amf_app_timer_t){
      AMF_APP_TIMER_INACTIVE_ID, AMF_APP_UE_CONTEXT_MODIFICATION_TIMER_VALUE};

  return new_p;
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_context_get()                                             **
 **                                                                        **
 ** Description: Get the amf context based on ue identity                  **
 **                                                                        **
 ** Input : ue_id : user equipment identity value                          **
 **                                                                        **
 ** Return: amf_context structure,Success case                             **
 **         NULL ,Failure case                                             **
 ***************************************************************************/
amf_context_t* amf_context_get(const amf_ue_ngap_id_t ue_id) {
  amf_context_t* amf_context_p = nullptr;

  if (INVALID_AMF_UE_NGAP_ID != ue_id) {
    ue_m5gmm_context_s* ue_mm_context =
        &ue_m5gmm_global_context;  // TODO upcoming PR has MAP implementation
                                   // will be removed.
    if (ue_mm_context) {
      amf_context_p = &ue_mm_context->amf_context;
    }
    OAILOG_DEBUG(LOG_NAS_AMF, "Stored UE id " AMF_UE_NGAP_ID_FMT " \n", ue_id);
  }
  return amf_context_p;
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_ue_context_exists_amf_ue_ngap_id()                        **
 **                                                                        **
 ** Description: Checks if UE context exists already or not                **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
ue_m5gmm_context_s* amf_ue_context_exists_amf_ue_ngap_id(
    const amf_ue_ngap_id_t amf_ue_ngap_id) {
  ue_m5gmm_context_s* ue_context_p = NULL;
  hash_table_ts_t* state_imsi_ht   = get_amf_ue_state();
  hashtable_ts_get(
      state_imsi_ht, (const hash_key_t) amf_ue_ngap_id, (void**) &ue_context_p);

  if (ue_context_p) {
    OAILOG_TRACE(
        LOG_AMF_APP,
        "UE  " AMF_UE_NGAP_ID_FMT " fetched AMF state %s with state %s\n ",
        amf_ue_ngap_id,
        (ue_context_p->mm_state == UNREGISTERED) ?
            "UE_UNREGISTERED" :
            (ue_context_p->mm_state == REGISTERED_CONNECTED) ? "UE_REGISTERED" :
                                                               "UNKNOWN",
        (ue_context_p->mm_state == REGISTERED_IDLE) ?
            "REGISTERED_IDLE" :
            (ue_context_p->mm_state == REGISTERED_CONNECTED) ?
            "REGISTERED_CONNECTED" :
            "UNKNOWN");
  }
  return ue_context_p;
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_context_upsert_imsi()                                     **
 **                                                                        **
 ** Description: upsert imsi in amf context                                **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
// in upcoming PR with MAP implementation, this routine will be depricated
int amf_context_upsert_imsi(amf_context_t* elm) {
  hashtable_rc_t h_rc = HASH_TABLE_OK;
  amf_ue_ngap_id_t ue_id =
      (PARENT_STRUCT(elm, ue_m5gmm_context_s, amf_context))->amf_ue_ngap_id;
  amf_app_desc_t* amf_app_desc_p = get_amf_nas_state(false);
  h_rc                           = hashtable_uint64_ts_remove(
      amf_app_desc_p->amf_ue_contexts.imsi_amf_ue_id_htbl,
      (const hash_key_t) elm->imsi64);

  if (ue_id != INVALID_AMF_UE_NGAP_ID) {
    h_rc = hashtable_uint64_ts_insert(
        amf_app_desc_p->amf_ue_contexts.imsi_amf_ue_id_htbl,
        (const hash_key_t) elm->imsi64, ue_id);
  } else {
    h_rc = HASH_TABLE_KEY_NOT_EXISTS;
  }
  if (h_rc != HASH_TABLE_OK) {
    OAILOG_TRACE(
        LOG_AMF_APP,
        "Error could not update this ue context "
        "amf_ue_s1ap_id " AMF_UE_S1AP_ID_FMT " imsi " IMSI_64_FMT ": %s\n",
        ue_id, elm->imsi64, hashtable_rc_code2string(h_rc));
    return RETURNerror;
  }
  return RETURNok;
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_remove_ue_context()                                       **
 **                                                                        **
 ** Description: Removes UE context                                        **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/

void amf_remove_ue_context(
    amf_ue_context_t* const amf_ue_context_p,
    ue_m5gmm_context_s* ue_context_p) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  hashtable_rc_t hash_rc              = HASH_TABLE_OK;
  hash_table_ts_t* amf_state_ue_id_ht = NULL;

  if (!amf_ue_context_p) {
    OAILOG_ERROR(LOG_AMF_APP, "Invalid AMF UE context received\n");
    OAILOG_FUNC_OUT(LOG_AMF_APP);
  }
  if (!ue_context_p) {
    OAILOG_ERROR(LOG_AMF_APP, "Invalid UE context received\n");
    OAILOG_FUNC_OUT(LOG_AMF_APP);
  }

  // IMSI
  if (ue_context_p->amf_context.imsi64) {
    hash_rc = hashtable_uint64_ts_remove(
        amf_ue_context_p->imsi_amf_ue_id_htbl,
        (const hash_key_t) ue_context_p->amf_context.imsi64);
    if (HASH_TABLE_OK != hash_rc) {
      OAILOG_ERROR_UE(
          LOG_AMF_APP, ue_context_p->amf_context.imsi64,
          "UE context not found!\n"
          " gnb_ue_ngap_id " GNB_UE_NGAP_ID_FMT
          " amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT " not in IMSI collection\n",
          ue_context_p->gnb_ue_ngap_id, ue_context_p->amf_ue_ngap_id);
    }
  }

  // gNB UE NGAP UE ID
  hash_rc = hashtable_uint64_ts_remove(
      amf_ue_context_p->gnb_ue_ngap_id_ue_context_htbl,
      (const hash_key_t) ue_context_p->gnb_ngap_id_key);

  if (HASH_TABLE_OK != hash_rc)
    OAILOG_ERROR(
        LOG_AMF_APP,
        "UE context not found!\n"
        " gnb_ue_ngap_id " GNB_UE_NGAP_ID_FMT
        " amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT
        ", GNB_UE_NGAP_ID not in GNB_UE_NGAP_ID collection",
        ue_context_p->gnb_ue_ngap_id, ue_context_p->amf_ue_ngap_id);

  // filled N11 tun id
  if (ue_context_p->amf_teid_n11) {
    hash_rc = hashtable_uint64_ts_remove(
        amf_ue_context_p->tun11_ue_context_htbl,
        (const hash_key_t) ue_context_p->amf_teid_n11);
    if (HASH_TABLE_OK != hash_rc)
      OAILOG_ERROR(
          LOG_AMF_APP,
          "UE Context not found!\n"
          " gnb_ue_ngap_id " GNB_UE_NGAP_ID_FMT
          " amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT ", AMF S11 TEID  " TEID_FMT
          "  not in N11 collection\n",
          ue_context_p->gnb_ue_ngap_id, ue_context_p->amf_ue_ngap_id,
          ue_context_p->amf_teid_n11);
  }

  // filled guti
  if ((ue_context_p->amf_context.m5_guti.guamfi.amf_set_id) ||
      (ue_context_p->amf_context.m5_guti.guamfi.amf_regionid) ||
      (ue_context_p->amf_context.m5_guti.m_tmsi) ||
      (ue_context_p->amf_context.m5_guti.guamfi.plmn.mcc_digit1) ||
      (ue_context_p->amf_context.m5_guti.guamfi.plmn.mcc_digit2) ||
      (ue_context_p->amf_context.m5_guti.guamfi.plmn
           .mcc_digit3)) {  // MCC 000 does not exist in ITU table
    hash_rc = obj_hashtable_uint64_ts_remove(
        amf_ue_context_p->guti_ue_context_htbl,
        (const void* const) & ue_context_p->amf_context.m5_guti,
        sizeof(ue_context_p->amf_context.m5_guti));
    if (HASH_TABLE_OK != hash_rc)
      OAILOG_ERROR(
          LOG_AMF_APP,
          "UE Context not found!\n"
          " gnb_ue_ngap_id " GNB_UE_NGAP_ID_FMT
          " amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT
          ", GUTI  not in GUTI collection\n",
          ue_context_p->gnb_ue_ngap_id, ue_context_p->amf_ue_ngap_id);
  }

  // filled NAS UE ID/ AMF UE NGAP ID
  if (INVALID_AMF_UE_NGAP_ID != ue_context_p->amf_ue_ngap_id) {
    hash_rc = hashtable_ts_remove(
        amf_state_ue_id_ht, (const hash_key_t) ue_context_p->amf_ue_ngap_id,
        (void**) &ue_context_p);
    if (HASH_TABLE_OK != hash_rc)
      OAILOG_ERROR(
          LOG_AMF_APP,
          "UE context not found!\n"
          "  gnb_ue_ngap_id " GNB_UE_NGAP_ID_FMT
          ", amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT
          " not in AMF UE NGAP ID collection",
          ue_context_p->gnb_ue_ngap_id, ue_context_p->amf_ue_ngap_id);
  }

  free_wrapper((void**) &ue_context_p);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_app_state_free_ue_context()                               **
 **                                                                        **
 ** Description: Cleans up AMF context                                     **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
void amf_app_state_free_ue_context(void** ue_context_node) {
  OAILOG_FUNC_IN(LOG_AMF_APP);

  // TODO clean up AMF context. This has been taken care in new PR with support
  // for Multi UE.
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}
}  // namespace magma5g
