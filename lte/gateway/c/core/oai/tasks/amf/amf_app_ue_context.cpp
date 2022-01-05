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
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/lib/directoryd/directoryd.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#ifdef __cplusplus
}
#endif
#include <unordered_map>
#include "lte/gateway/c/core/oai/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_state_manager.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_recv.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_common.h"
#include "lte/gateway/c/core/oai/include/map.h"

namespace magma5g {
extern task_zmq_ctx_t amf_app_task_zmq_ctx;
// Creating ue_context_map based on key:ue_id and value:ue_context
std::unordered_map<amf_ue_ngap_id_t, ue_m5gmm_context_s*> ue_context_map;
// Creating smf_ctxt_map based on key:pdu_session_id and value:smf_context
std::unordered_map<uint8_t, std::shared_ptr<smf_context_t>> smf_ctxt_map;

std::shared_ptr<smf_context_t> amf_insert_smf_context(
    ue_m5gmm_context_s*, uint8_t);

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
    OAILOG_ERROR(LOG_AMF_APP, "UE context is null\n");
    OAILOG_FUNC_OUT(LOG_AMF_APP);
  }
  message_p =
      itti_alloc_new_message(TASK_AMF_APP, AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION);
  notification_p = &message_p->ittiMsg.amf_app_ngap_amf_ue_id_notification;
  memset(notification_p, 0, sizeof(itti_amf_app_ngap_amf_ue_id_notification_t));
  notification_p->gnb_ue_ngap_id = ue_context_p->gnb_ue_ngap_id;
  notification_p->amf_ue_ngap_id = ue_context_p->amf_ue_ngap_id;
  notification_p->sctp_assoc_id  = ue_context_p->sctp_assoc_id_key;

  amf_send_msg_to_task(&amf_app_task_zmq_ctx, TASK_NGAP, message_p);
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
    amf_ue_ngap_id_t ue_id, ue_m5gmm_context_s* ue_context_p) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  if (ue_context_p == NULL) {
    OAILOG_ERROR(LOG_AMF_APP, "Invalid UE context received\n");
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
      OAILOG_DEBUG(
          LOG_AMF_APP,
          "Overwriting the Existing entry UE_ID = " AMF_UE_NGAP_ID_FMT, ue_id);
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
 ** Description: Insert smf context in the map                             **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
std::shared_ptr<smf_context_t> amf_insert_smf_context(
    ue_m5gmm_context_s* ue_context, uint8_t pdu_session_id) {
  std::shared_ptr<smf_context_t> smf_context;
  smf_context =
      amf_get_smf_context_by_pdu_session_id(ue_context, pdu_session_id);
  if (smf_context) {
    return smf_context;
  } else {
    smf_context = std::make_shared<smf_context_t>();
    ue_context->amf_context.smf_ctxt_map[pdu_session_id] = smf_context;
  }
  return smf_context;
}

/****************************************************************************
 **                                                                        **
 ** Name:   amf_get_smf_context_by_pdu_session_id()                        **
 **                                                                        **
 ** Description: Get the smf context from the map                          **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
std::shared_ptr<smf_context_t> amf_get_smf_context_by_pdu_session_id(
    ue_m5gmm_context_s* ue_context, uint8_t id) {
  std::shared_ptr<smf_context_t> smf_context;
  for (const auto& it : ue_context->amf_context.smf_ctxt_map) {
    if (it.first == id) {
      smf_context = it.second;
      break;
    }
  }
  return smf_context;
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
  magma::map_rc_t m_rc = magma::MAP_OK;
  amf_ue_ngap_id_t ue_id =
      (PARENT_STRUCT(elm, ue_m5gmm_context_s, amf_context))->amf_ue_ngap_id;
  amf_app_desc_t* amf_app_desc_p = get_amf_nas_state(false);

  m_rc =
      amf_app_desc_p->amf_ue_contexts.imsi_amf_ue_id_htbl.remove(elm->imsi64);

  if (ue_id != INVALID_AMF_UE_NGAP_ID) {
    m_rc = amf_app_desc_p->amf_ue_contexts.imsi_amf_ue_id_htbl.insert(
        elm->imsi64, ue_id);
  } else {
    OAILOG_TRACE(
        LOG_AMF_APP,
        "Error could not update this ue context "
        "amf_ue_s1ap_id " AMF_UE_S1AP_ID_FMT " imsi " IMSI_64_FMT ": %s\n",
        ue_id, elm->imsi64, map_rc_code2string(m_rc).c_str());
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
 ** Name:    amf_lookup_guti_by_ueid()                                     **
 **                                                                        **
 ** Description:  Fetch the guti based on ue id                            **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
tmsi_t amf_lookup_guti_by_ueid(amf_ue_ngap_id_t ue_id) {
  amf_context_t* amf_ctxt = amf_context_get(ue_id);

  if (amf_ctxt == NULL) {
    return (0);
  }

  return amf_ctxt->m5_guti.m_tmsi;
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

  std::shared_ptr<smf_context_t> smf_ctx;
  for (auto& it : ue_context_p->amf_context.smf_ctxt_map) {
    smf_ctx                    = it.second;
    smf_ctx->pdu_session_state = INACTIVE;
  }

  amf_smf_notification_send(ue_id, ue_context_p, UE_IDLE_MODE_NOTIFY);

  OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNok);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_free_ue_context()                                         **
 **                                                                        **
 ** Description: Deletes the ue context                                    **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
void amf_free_ue_context(ue_m5gmm_context_s* ue_context_p) {
  hashtable_rc_t h_rc                = HASH_TABLE_OK;
  magma::map_rc_t m_rc               = magma::MAP_OK;
  amf_app_desc_t* amf_app_desc_p     = get_amf_nas_state(false);
  amf_ue_context_t* amf_ue_context_p = &amf_app_desc_p->amf_ue_contexts;
  OAILOG_DEBUG(LOG_NAS_AMF, "amf_free_ue_context \n");
  map_uint64_ue_context_t amf_state_ue_id_ht = get_amf_ue_state();
  if (!ue_context_p || !amf_ue_context_p) {
    return;
  }

  amf_remove_ue_context(ue_context_p);

  // Clean up the procedures
  amf_nas_proc_clean_up(ue_context_p);

  if (ue_context_p->gnb_ngap_id_key != INVALID_GNB_UE_NGAP_ID_KEY) {
    m_rc = amf_ue_context_p->gnb_ue_ngap_id_ue_context_htbl.remove(
        ue_context_p->gnb_ngap_id_key);

    if (m_rc != magma::MAP_OK)
      OAILOG_TRACE(LOG_AMF_APP, "Error Could not remove this ue context \n");
    ue_context_p->gnb_ngap_id_key = INVALID_GNB_UE_NGAP_ID_KEY;
  }

  if (ue_context_p->amf_ue_ngap_id != INVALID_AMF_UE_NGAP_ID) {
    m_rc = amf_state_ue_id_ht.remove(ue_context_p->amf_ue_ngap_id);
    if (m_rc != magma::MAP_OK)
      OAILOG_TRACE(LOG_AMF_APP, "Error Could not remove this ue context \n");
    ue_context_p->amf_ue_ngap_id = INVALID_AMF_UE_NGAP_ID;
  }

  amf_ue_context_p->imsi_amf_ue_id_htbl.remove(
      ue_context_p->amf_context.imsi64);

  amf_ue_context_p->tun11_ue_context_htbl.remove(ue_context_p->amf_teid_n11);

  amf_ue_context_p->guti_ue_context_htbl.remove(
      ue_context_p->amf_context.m5_guti);

  delete ue_context_p;
  ue_context_p = NULL;
}

}  // namespace magma5g
