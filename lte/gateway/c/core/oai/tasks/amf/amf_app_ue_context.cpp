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
#include <unordered_map>
#include "common_defs.h"
#include "dynamic_memory_check.h"
#include "amf_app_state_manager.h"
#include "amf_recv.h"

namespace magma5g {
extern task_zmq_ctx_t amf_app_task_zmq_ctx;
// Creating ue_context_map based on key:ue_id and value:ue_context
std::unordered_map<amf_ue_ngap_id_t, ue_m5gmm_context_s*> ue_context_map;

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
    amf_ue_ngap_id_t ue_id, amf_ue_context_t* amf_ue_context_p,
    ue_m5gmm_context_s* ue_context_p) {
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

  if (ue_context_map.size() == 0) {
    // first entry.
    ue_context_map.insert(
        std::pair<amf_ue_ngap_id_t, ue_m5gmm_context_s*>(ue_id, ue_context_p));
  } else {
    /* already elements exist then check if given ue_id already present
     * if it exists, update/overwrite the element. Otherwise add a new element
     */
    auto found_ue_id = ue_context_map.find(ue_id);
    if (found_ue_id == ue_context_map.end()) {
      // it is new entry to map
      ue_context_map.insert(std::pair<amf_ue_ngap_id_t, ue_m5gmm_context_s*>(
          ue_id, ue_context_p));
    } else {
      // Overwrite the existing element.
      found_ue_id->second = ue_context_p;
      OAILOG_INFO(
          LOG_AMF_APP, "Overwriting the Existing entry UE_ID=%u\n", ue_id);
    }
  }

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
  // Make ue_context zero initialize
  ue_m5gmm_context_s* new_p = new ue_m5gmm_context_s();

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
  new_p->mm_state = DEREGISTERED;

  new_p->amf_context._security.eksi = KSI_NO_KEY_AVAILABLE;
  new_p->mm_state                   = DEREGISTERED;

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
        amf_ue_context_exists_amf_ue_ngap_id(ue_id);

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
  std::unordered_map<amf_ue_ngap_id_t, ue_m5gmm_context_s*>::iterator
      found_ue_id = ue_context_map.find(amf_ue_ngap_id);

  if (found_ue_id == ue_context_map.end()) {
    return NULL;
  } else {
    return found_ue_id->second;
  }
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_insert_smf_context()                                      **
 **                                                                        **
 ** Description: Insert smf context in the vector                          **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
smf_context_t* amf_insert_smf_context(
    ue_m5gmm_context_s* ue_context, uint8_t pdu_session_id) {
  smf_context_t smf_context = {};
  std::vector<smf_context_t>::iterator i;
  int j = 0;

  for (i = ue_context->amf_context.smf_ctxt_vector.begin();
       i != ue_context->amf_context.smf_ctxt_vector.end(); i++, j++) {
    OAILOG_INFO(LOG_AMF_APP, "insert smf_ctx %d", j);
    if (i->smf_proc_data.pdu_session_identity.pdu_session_id ==
        pdu_session_id) {
      ue_context->amf_context.smf_ctxt_vector.at(j) = smf_context;
      return ue_context->amf_context.smf_ctxt_vector.data() + j;
    }
  }

  // add new element to the vector
  ue_context->amf_context.smf_ctxt_vector.push_back(smf_context);
  // i = ue_context->amf_context.smf_ctxt_vector.begin();
  // ue_context->amf_context.smf_ctxt_vector.insert(i, smf_context);
  // ue_context->amf_context.smf_ctxt_vector.at(0) = smf_context;
  OAILOG_INFO(
      LOG_AMF_APP, "insert smf_ctx_vector_data %lu",
      ue_context->amf_context.smf_ctxt_vector.size());
  return ue_context->amf_context.smf_ctxt_vector.data() + j;
  // return ue_context->amf_context.smf_ctxt_vector.data();
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_smf_context_exists_pdu_session_id()                       **
 **                                                                        **
 ** Description: Get the smf context from the vector                       **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
smf_context_t* amf_smf_context_exists_pdu_session_id(
    ue_m5gmm_context_s* ue_context, uint8_t id) {
  std::vector<smf_context_t>::iterator i;
  int j = 0;
  for (i = ue_context->amf_context.smf_ctxt_vector.begin();
       i != ue_context->amf_context.smf_ctxt_vector.end(); i++, j++) {
    if (i->smf_proc_data.pdu_session_identity.pdu_session_id == id) {
      return ue_context->amf_context.smf_ctxt_vector.data() + j;
    }
  }
  return NULL;
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

/****************************************************************************
 **                                                                        **
 ** Name:    ue_context_loopkup_by_guti()                                  **
 **                                                                        **
 ** Description: Checks if UE context exists already or not                **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
ue_m5gmm_context_s* ue_context_loopkup_by_guti(tmsi_t tmsi_rcv) {
  tmsi_t tmsi_stored;
  ue_m5gmm_context_s* ue_context;

  for (auto i = ue_context_map.begin(); i != ue_context_map.end(); i++) {
    ue_context = i->second;
    if (ue_context == NULL) {
      continue;
    }

    tmsi_stored = ue_context->amf_context.m5_guti.m_tmsi;

    if (tmsi_rcv == tmsi_stored) {
      return ue_context;
    }
  }

  return NULL;
}

/****************************************************************************
 **                                                                        **
 ** Name:    ue_context_update_ue_id()                                     **
 **                                                                        **
 ** Description:  Update the UE_ID                                         **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
void ue_context_update_ue_id(
    ue_m5gmm_context_s* ue_context, amf_ue_ngap_id_t ue_id) {
  if (ue_id == ue_context->amf_ue_ngap_id) {
    return;
  }

  /* Erase the content with ue_id */
  ue_context_map.erase(ue_context->amf_ue_ngap_id);

  /* Re-Insert with same context but with updated Key */
  ue_context_map.insert(
      std::pair<amf_ue_ngap_id_t, ue_m5gmm_context_s*>(ue_id, ue_context));

  return;
}

/****************************************************************************
 **                                                                        **
 ** Name:    ue_context_lookup_by_gnb_ue_id()                              **
 **                                                                        **
 ** Description:  Fetch the ue_context by gnb ue id                        **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
ue_m5gmm_context_s* ue_context_lookup_by_gnb_ue_id(
    gnb_ue_ngap_id_t gnb_ue_ngap_id) {
  ue_m5gmm_context_s* ue_context;
  for (auto i = ue_context_map.begin(); i != ue_context_map.end(); i++) {
    ue_context = i->second;
    if (ue_context == NULL) {
      continue;
    }

    if (ue_context->gnb_ue_ngap_id == gnb_ue_ngap_id) {
      return ue_context;
    }
  }
  return NULL;
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_idle_mode_procedure()                                     **
 **                                                                        **
 ** Description:  Transition to idle mode                                  **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
int amf_idle_mode_procedure(amf_context_t* amf_ctx) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  ue_m5gmm_context_s* ue_context_p =
      PARENT_STRUCT(amf_ctx, ue_m5gmm_context_s, amf_context);
  amf_ue_ngap_id_t ue_id = ue_context_p->amf_ue_ngap_id;

  smf_context_t smf_ctx;

  for (auto it = ue_context_p->amf_context.smf_ctxt_vector.begin();
       it != ue_context_p->amf_context.smf_ctxt_vector.end(); it++) {
    smf_ctx                   = *it;
    smf_ctx.pdu_session_state = INACTIVE;
  }

  amf_smf_notification_send(ue_id, ue_context_p, UE_IDLE_MODE_NOTIFY);

  OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNok);
}
}  // namespace magma5g
