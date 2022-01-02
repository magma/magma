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
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_timer_management.h"

namespace magma5g {
extern task_zmq_ctx_t amf_app_task_zmq_ctx;
// TODO(panyogesh): To be removed as part of Fixing the global maps
extern std::unordered_map<imsi64_t, guti_and_amf_id_t> amf_supi_guti_map;
// Creating ue_context_map based on key:ue_id and value:ue_context
std::unordered_map<amf_ue_ngap_id_t, ue_m5gmm_context_s*> ue_context_map;

std::shared_ptr<smf_context_t> amf_insert_smf_context(ue_m5gmm_context_s*,
                                                      uint8_t);

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
  MessageDef* message_p = NULL;
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
  notification_p->sctp_assoc_id = ue_context_p->sctp_assoc_id_key;

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
status_code_e amf_insert_ue_context(
    amf_ue_context_t* const amf_ue_context_p,
    struct ue_m5gmm_context_s* const ue_context_p) {
  magma::map_rc_t m_rc                        = magma::MAP_OK;
  map_uint64_ue_context_t* amf_state_ue_id_ht = get_amf_ue_state();

  OAILOG_FUNC_IN(LOG_AMF_APP);
  if (amf_ue_context_p == NULL) {
    OAILOG_ERROR(LOG_MME_APP, "Invalid AMF UE context received\n");
    OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
  }

  if (ue_context_p == NULL) {
    OAILOG_ERROR(LOG_AMF_APP, "Invalid UE context received\n");
    OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
  }

  uint64_t amf_ue_ngap_id64 = 0;
  m_rc = amf_ue_context_p->gnb_ue_ngap_id_ue_context_htbl.get(
      ue_context_p->gnb_ngap_id_key, &amf_ue_ngap_id64);
  if (m_rc == magma::MAP_OK) {
    OAILOG_WARNING(
        LOG_AMF_APP,
        "This ue context %p already exists gnb_ue_ngap_id " GNB_UE_NGAP_ID_FMT
        "\n",
        ue_context_p, ue_context_p->gnb_ue_ngap_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  m_rc = amf_ue_context_p->gnb_ue_ngap_id_ue_context_htbl.insert(
      ue_context_p->gnb_ngap_id_key, ue_context_p->amf_ue_ngap_id);

  if (m_rc != magma::MAP_OK) {
    OAILOG_WARNING(
        LOG_AMF_APP,
        "Failed to insert ue context entry  " GNB_UE_NGAP_ID_FMT
        "in gnb_ue_ngap_id_ue_context_htbl \n",
        ue_context_p->gnb_ue_ngap_id);
  }

  if (INVALID_AMF_UE_NGAP_ID != ue_context_p->amf_ue_ngap_id) {
    ue_m5gmm_context_s* tmp_ue_context_p = NULL;
    if (amf_state_ue_id_ht->get(
            ue_context_p->amf_ue_ngap_id, &tmp_ue_context_p) == magma::MAP_OK) {
      OAILOG_WARNING(
          LOG_MME_APP,
          "This ue context %p already exists amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT
          "\n",
          tmp_ue_context_p, ue_context_p->amf_ue_ngap_id);
      OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
    }

    m_rc =
        amf_state_ue_id_ht->insert(ue_context_p->amf_ue_ngap_id, ue_context_p);
    if (m_rc != magma::MAP_OK) {
      OAILOG_WARNING(
          LOG_MME_APP,
          "Error could not register this ue context %p "
          "amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT "\n",
          ue_context_p, ue_context_p->amf_ue_ngap_id);
      OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
    }

    // filled IMSI
    if (ue_context_p->amf_context.imsi64) {
      m_rc = amf_ue_context_p->imsi_amf_ue_id_htbl.insert(
          ue_context_p->amf_context.imsi64, ue_context_p->amf_ue_ngap_id);

      if (m_rc != magma::MAP_OK) {
        OAILOG_WARNING_UE(
            LOG_AMF_APP, ue_context_p->amf_context.imsi64,
            "Error could not register this ue context %p "
            "amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT " imsi " IMSI_64_FMT "\n",
            ue_context_p, ue_context_p->amf_ue_ngap_id,
            ue_context_p->amf_context.imsi64);
        OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
      }
    }

    // filled guti
    if ((0 != ue_context_p->amf_context.m5_guti.guamfi.amf_regionid) ||
        (0 != ue_context_p->amf_context.m5_guti.guamfi.amf_set_id) ||
        (INVALID_TMSI != ue_context_p->amf_context.m5_guti.m_tmsi) ||
        (0 != ue_context_p->amf_context.m5_guti.guamfi.plmn
                  .mcc_digit1) ||  // MCC 000 does not exist in ITU table
        (0 != ue_context_p->amf_context.m5_guti.guamfi.plmn.mcc_digit2) ||
        (0 != ue_context_p->amf_context.m5_guti.guamfi.plmn.mcc_digit3)) {
      m_rc = amf_ue_context_p->guti_ue_context_htbl.insert(
          ue_context_p->amf_context.m5_guti, ue_context_p->amf_ue_ngap_id);

      if (m_rc != magma::MAP_OK) {
        OAILOG_WARNING(
            LOG_AMF_APP,
            "Error could not register this ue context %p "
            "amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT " \n",
            ue_context_p, ue_context_p->amf_ue_ngap_id);

        OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
      }
    }
  }

  OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNok);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_init_amf_context()                                        **
 **                                                                        **
 ** Description: Init amf context                                          **
 **                                                                        **
 ***************************************************************************/
void amf_init_amf_context(amf_context_t* amf_ctx) {
  amf_ctx->_security.eksi        = KSI_NO_KEY_AVAILABLE;
  amf_ctx->m5_guti.m_tmsi        = INVALID_TMSI;
  amf_ctx->new_registration_info = NULL;
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

  new_p->amf_ue_ngap_id = INVALID_AMF_UE_NGAP_ID;
  new_p->gnb_ngap_id_key = INVALID_GNB_UE_NGAP_ID_KEY;
  new_p->gnb_ue_ngap_id = INVALID_GNB_UE_NGAP_ID;

  // Initialize timers to INVALID IDs
  new_p->m5_mobile_reachability_timer.id = AMF_APP_TIMER_INACTIVE_ID;
  new_p->m5_implicit_deregistration_timer.id = AMF_APP_TIMER_INACTIVE_ID;
  new_p->m5_initial_context_setup_rsp_timer  = (amf_app_timer_t){
      AMF_APP_TIMER_INACTIVE_ID, AMF_APP_INITIAL_CONTEXT_SETUP_RSP_TIMER_VALUE};
  new_p->m5_ulr_response_timer = (amf_app_timer_t){
      AMF_APP_TIMER_INACTIVE_ID, AMF_APP_ULR_RESPONSE_TIMER_VALUE};
  new_p->m5_ue_context_modification_timer = (amf_app_timer_t){
      AMF_APP_TIMER_INACTIVE_ID, AMF_APP_UE_CONTEXT_MODIFICATION_TIMER_VALUE};
  new_p->mm_state = DEREGISTERED;

  new_p->amf_context._security.eksi = KSI_NO_KEY_AVAILABLE;
  new_p->mm_state = DEREGISTERED;
  new_p->pending_service_response = false;

  amf_init_amf_context(&new_p->amf_context);

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
 ** Name:    amf_ue_context_exists_imsi()                                  **
 **                                                                        **
 ** Description: Checks if UE context exists for IMSI or not               **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
struct ue_m5gmm_context_s* amf_ue_context_exists_imsi(
    amf_ue_context_t* const amf_ue_context_p, imsi64_t imsi64) {
  magma::map_rc_t m_rc      = magma::MAP_OK;
  uint64_t amf_ue_ngap_id64 = 0;

  amf_ue_context_p->imsi_amf_ue_id_htbl.get(imsi64, &amf_ue_ngap_id64);
  if (m_rc == magma::MAP_OK) {
    return amf_ue_context_exists_amf_ue_ngap_id(
        (amf_ue_ngap_id_t) amf_ue_ngap_id64);
  } else {
    OAILOG_WARNING_UE(
        LOG_AMF_APP, imsi64, " No IMSI hashtable for this IMSI\n");
  }
  return NULL;
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_get_ue_context_from_imsi()                                **
 **                                                                        **
 ** Description: Fettch the UE context from IMSI                           **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
ue_m5gmm_context_s* amf_get_ue_context_from_imsi(char* imsi) {
  amf_context_t* amf_context_p = nullptr;
  imsi64_t imsi64              = INVALID_IMSI64;

  IMSI_STRING_TO_IMSI64((char*) imsi, &imsi64);

  amf_app_desc_t* amf_app_desc_p = get_amf_nas_state(false);

  return (amf_ue_context_exists_imsi(&amf_app_desc_p->amf_ue_contexts, imsi64));
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_ue_context_exists_guti()                                  **
 **                                                                        **
 ** Description: Checks if UE context exists for GUTI or not               **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
ue_m5gmm_context_s* amf_ue_context_exists_guti(
    amf_ue_context_t* const amf_ue_context_p, const guti_m5_t* const guti_p) {
  magma::map_rc_t m_rc           = magma::MAP_OK;
  uint64_t amf_ue_ngap_id64      = 0;
  ue_m5gmm_context_t* ue_context = NULL;

  m_rc = amf_ue_context_p->guti_ue_context_htbl.get(*guti_p, &amf_ue_ngap_id64);

  if (m_rc == magma::MAP_OK) {
    ue_context = amf_ue_context_exists_amf_ue_ngap_id(
        (amf_ue_ngap_id_t) amf_ue_ngap_id64);
    if (ue_context) {
      return ue_context;
    }
  } else {
    OAILOG_WARNING(LOG_AMF_APP, " No GUTI hashtable for GUTI ");
  }

  return NULL;
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
  struct ue_m5gmm_context_s* ue_context_p     = NULL;
  map_uint64_ue_context_t* amf_state_ue_id_ht = get_amf_ue_state();
  magma::map_rc_t m_rc                        = magma::MAP_OK;

  if (amf_state_ue_id_ht->get(amf_ue_ngap_id, &ue_context_p) != magma::MAP_OK) {
    OAILOG_WARNING(
        LOG_AMF_APP, " amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT "does not exist\n",
        amf_ue_ngap_id);
    return ue_context_p;
  }

  return ue_context_p;
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_ue_context_exists_gnb_ue_ngap_id()                        **
 **                                                                        **
 ** Description: Checks if UE context exists already or not based on       **
 **              gnb_key                                                   **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
ue_m5gmm_context_s* amf_ue_context_exists_gnb_ue_ngap_id(
    amf_ue_context_t* const amf_ue_context_p, const gnb_ngap_id_key_t gnb_key) {
  magma::map_rc_t m_rc      = magma::MAP_OK;
  uint64_t amf_ue_ngap_id64 = 0;

  m_rc = amf_ue_context_p->gnb_ue_ngap_id_ue_context_htbl.get(
      gnb_key, &amf_ue_ngap_id64);
  if (m_rc == magma::MAP_OK) {
    return amf_ue_context_exists_amf_ue_ngap_id(amf_ue_ngap_id64);
  }

  return NULL;
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
    OAILOG_TRACE(LOG_AMF_APP,
                 "Error could not update this ue context "
                 "amf_ue_s1ap_id " AMF_UE_S1AP_ID_FMT " imsi " IMSI_64_FMT
                 ": %s\n",
                 ue_id, elm->imsi64, map_rc_code2string(m_rc).c_str());
    return RETURNerror;
  }
  return RETURNok;
}

/****************************************************************************
 **                                                                        **
 ** Name:    lookup_ue_ctxt_by_imsi()                                      **
 **                                                                        **
 ** Description: Lookupif imsi exists                                      **
 **                                                                        **
 ** Inputs:  imsi64: imsi value                                            **
 **                                                                        **
 ** Outputs: ue_m5gmm_context_s: pointer to ue context                     **
 **                                                                        **
 ***************************************************************************/
ue_m5gmm_context_s* lookup_ue_ctxt_by_imsi(imsi64_t imsi64) {
  amf_app_desc_t* amf_app_desc_p = get_amf_nas_state(false);

  return (amf_ue_context_exists_imsi(&amf_app_desc_p->amf_ue_contexts, imsi64));
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
    smf_ctx = it.second;
    smf_ctx->pdu_session_state = INACTIVE;
    amf_smf_notification_send(ue_id, ue_context_p, UE_IDLE_MODE_NOTIFY,
                              it.first);
  }

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
  amf_app_desc_t* amf_app_desc_p = get_amf_nas_state(false);
  amf_ue_context_t* amf_ue_context_p = &amf_app_desc_p->amf_ue_contexts;
  OAILOG_DEBUG(LOG_NAS_AMF, "amf_free_ue_context \n");
  map_uint64_ue_context_t* amf_state_ue_id_ht = get_amf_ue_state();
  if (!ue_context_p || !amf_ue_context_p) {
    return;
  }

  amf_remove_ue_context(&amf_app_desc_p->amf_ue_contexts, ue_context_p);
}

/****************************************************************************
 **                                                                        **
 ** Name:    proc_new_registration_req()                                   **
 **                                                                        **
 ** Description: Restarts the new registration procedure for stored        **
 **              attached registration information                         **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
void proc_new_registration_req(
    amf_ue_context_t* const amf_ue_context_p,
    struct ue_m5gmm_context_s* ue_context_p) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);

  OAILOG_INFO(
      LOG_NAS_AMF,
      "Process new Registration Request for ue_id " AMF_UE_NGAP_ID_FMT "\n",
      ue_context_p->amf_ue_ngap_id);

  new_registration_info_t registration_info = {};
  memcpy(
      &registration_info, ue_context_p->amf_context.new_registration_info,
      sizeof(new_registration_info_t));

  delete (ue_context_p->amf_context.new_registration_info);

  /* The new Registration Request is received in ngap initial ue message,
   * So release previous Registration Request's contexts
   */
  if (registration_info.is_mm_ctx_new) {
    ue_context_p->ue_context_rel_cause = NGAP_NAS_DEREGISTER;

    amf_ctx_release_ue_context(
        ue_context_p, ue_context_p->ue_context_rel_cause);

    ue_context_p->ue_context_rel_cause = NGAP_INVALID_CAUSE;
  } else {
    uint64_t amf_ue_ngap_id64 = 0;
    magma::map_rc_t m_rc      = magma::MAP_OK;

    m_rc = amf_ue_context_p->guti_ue_context_htbl.get(
        ue_context_p->amf_context.m5_guti, &amf_ue_ngap_id64);

    if (m_rc == magma::MAP_OK) {
      // While processing new attach req, remove GUTI from hashtable
      if ((ue_context_p->amf_context.m5_guti.guamfi.amf_regionid) ||
          (ue_context_p->amf_context.m5_guti.guamfi.amf_set_id) ||
          (ue_context_p->amf_context.m5_guti.m_tmsi) ||
          (ue_context_p->amf_context.m5_guti.guamfi.plmn.mcc_digit1) ||
          (ue_context_p->amf_context.m5_guti.guamfi.plmn.mcc_digit2) ||
          (ue_context_p->amf_context.m5_guti.guamfi.plmn.mcc_digit3)) {
        amf_ue_context_p->guti_ue_context_htbl.remove(
            ue_context_p->amf_context.m5_guti);
      }
    }
  }

  amf_remove_ue_context(amf_ue_context_p, ue_context_p);

  // Proceed with new attach request
  ue_m5gmm_context_t* ue_m5gmm_context =
      amf_ue_context_exists_amf_ue_ngap_id(registration_info.amf_ue_ngap_id);
  amf_context_t* new_amf_ctx = &ue_m5gmm_context->amf_context;

  amf_init_amf_context(new_amf_ctx);

  new_amf_ctx->is_dynamic = true;

  if (!is_nas_specific_procedure_registration_running(
          &ue_m5gmm_context->amf_context)) {
    amf_proc_create_procedure_registration_request(
        ue_m5gmm_context, registration_info.ies);
  }
  amf_registration_run_procedure(&ue_m5gmm_context->amf_context);

  OAILOG_FUNC_OUT(LOG_NAS_AMF);
}

//------------------------------------------------------------------------------
int amf_app_handle_implicit_deregistration_timer_expiry(
    zloop_t* loop, int timer_id, void* args) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);

  amf_context_t* amf_ctx = NULL;
  amf_ue_ngap_id_t ue_id = 0;

  if (!amf_pop_timer_arg(timer_id, &ue_id)) {
    OAILOG_WARNING(
        LOG_AMF_APP,
        "Implicit Deregistration: Invalid Timer Id expiration, Timer Id: %u\n",
        timer_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
  }

  ue_m5gmm_context_s* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  if (ue_context_p == NULL) {
    OAILOG_DEBUG(
        LOG_NAS_AMF,
        "Implicit Deregistration: ue_amf_context is NULL for "
        "ue id: " AMF_UE_NGAP_ID_FMT "\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
  }

  amf_ctx = &ue_context_p->amf_context;

  if (!(amf_ctx)) {
    OAILOG_ERROR(
        LOG_AMF_APP,
        "Implicit Deregistration: Timer expired no amf context for "
        "ue id: " AMF_UE_NGAP_ID_FMT "\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
  }

  ue_context_p->m5_implicit_deregistration_timer.id = AMF_APP_TIMER_INACTIVE_ID;

  // Initiate Implicit Detach for the UE
  amf_nas_proc_implicit_deregister_ue_ind(ue_id);
  OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNok);
}

//------------------------------------------------------------------------------
static int amf_app_handle_mobile_reachability_timer_expiry(
    zloop_t* loop, int timer_id, void* args) {
  OAILOG_FUNC_IN(LOG_AMF_APP);

  amf_context_t* amf_ctx = NULL;
  amf_ue_ngap_id_t ue_id = 0;

  if (!amf_pop_timer_arg(timer_id, &ue_id)) {
    OAILOG_WARNING(
        LOG_AMF_APP,
        "Mobile Rechability timer: Invalid Timer Id expiration, Timer Id: %u\n",
        timer_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
  }

  ue_m5gmm_context_s* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  if (ue_context_p == NULL) {
    OAILOG_DEBUG(
        LOG_NAS_AMF,
        "Mobile Reachability Timer: ue_amf_context is NULL for "
        "ue id: " AMF_UE_NGAP_ID_FMT "\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
  }

  amf_ctx = &ue_context_p->amf_context;

  if (!(amf_ctx)) {
    OAILOG_ERROR(
        LOG_AMF_APP,
        "Mobile Reachability Timer: Timer expired no amf context for "
        "ue id: " AMF_UE_NGAP_ID_FMT "\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
  }

  ue_context_p->m5_mobile_reachability_timer.id = AMF_APP_TIMER_INACTIVE_ID;

  // Start Implicit Deregister timer
  ue_context_p->m5_implicit_deregistration_timer.id = amf_app_start_timer(
      ue_context_p->m5_implicit_deregistration_timer.sec * 1000,
      TIMER_REPEAT_ONCE, amf_app_handle_implicit_deregistration_timer_expiry,
      ue_id);

  ue_context_p->ppf = false;
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

void amf_ue_context_update_ue_sig_connection_state(
    amf_ue_context_t* const amf_ue_context_p,
    struct ue_m5gmm_context_s* ue_context_p, m5gcm_state_t new_cm_state) {
  // Function is used to update UE's Signaling Connection State
  hashtable_rc_t hash_rc = HASH_TABLE_OK;

  OAILOG_FUNC_IN(LOG_AMF_APP);
  if (amf_ue_context_p == NULL) {
    OAILOG_ERROR(LOG_AMF_APP, "Invalid AMF UE context received\n");
    OAILOG_FUNC_OUT(LOG_AMF_APP);
  }

  if (ue_context_p == NULL) {
    OAILOG_ERROR(LOG_AMF_APP, "Invalid UE context received\n");
    OAILOG_FUNC_OUT(LOG_AMF_APP);
  }

  if (ue_context_p->cm_state == M5GCM_CONNECTED && new_cm_state == M5GCM_IDLE) {
    magma::map_rc_t m_rc = magma::MAP_OK;
    m_rc = amf_ue_context_p->gnb_ue_ngap_id_ue_context_htbl.remove(
        ue_context_p->gnb_ngap_id_key);

    if (m_rc != magma::MAP_OK) {
      OAILOG_WARNING_UE(
          LOG_AMF_APP, ue_context_p->amf_context.imsi64,
          "UE context gnb_ue_ngap_ue_id_key %ld "
          "amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT
          ", GNB_UE_NGAP_ID_KEY could not be found",
          ue_context_p->gnb_ngap_id_key, ue_context_p->amf_ue_ngap_id);
    }
    ue_context_p->gnb_ngap_id_key = INVALID_GNB_UE_NGAP_ID_KEY;

    OAILOG_DEBUG_UE(
        LOG_AMF_APP, ue_context_p->amf_context.imsi64,
        "AMF_APP: UE Connection State changed to IDLE. amf_ue_ngap_id "
        "= " AMF_UE_NGAP_ID_FMT "\n",
        ue_context_p->amf_ue_ngap_id);

    if (amf_config.nas_config.t3512_min &&
        ue_context_p->m5_mobile_reachability_timer.id ==
            AMF_APP_TIMER_INACTIVE_ID) {
      ue_context_p->m5_mobile_reachability_timer.sec =
          (amf_config.nas_config.t3512_min + 4 * 60);

      if (ue_context_p->m5_mobile_reachability_timer.sec ==
          AMF_APP_TIMER_INACTIVE_ID) {
        ue_context_p->m5_mobile_reachability_timer.id = amf_app_start_timer(
            ue_context_p->m5_mobile_reachability_timer.sec * 1000,
            TIMER_REPEAT_ONCE, amf_app_handle_mobile_reachability_timer_expiry,
            ue_context_p->amf_ue_ngap_id);
      }
    }

    ue_context_p->cm_state = M5GCM_IDLE;

    // Update Stats
    OAILOG_INFO_UE(
        LOG_MME_APP, ue_context_p->amf_context.imsi64, "UE STATE - IDLE.\n");

  } else if (
      (ue_context_p->cm_state == M5GCM_IDLE) &&
      (new_cm_state == M5GCM_CONNECTED)) {
    ue_context_p->cm_state = M5GCM_CONNECTED;

    OAILOG_DEBUG_UE(
        LOG_AMF_APP, ue_context_p->amf_context.imsi64,
        "AMF_APP: UE Connection State changed to CONNECTED.enb_ue_s1ap_id "
        "=" GNB_UE_NGAP_ID_FMT ", mme_ue_s1ap_id = " AMF_UE_NGAP_ID_FMT "\n",
        ue_context_p->gnb_ue_ngap_id, ue_context_p->amf_ue_ngap_id);
    // Set PPF flag to true whenever UE moves from M5GCM_IDLE to M5GCM_CONNECTED
    // state
    ue_context_p->ppf = true;

    OAILOG_INFO_UE(
        LOG_AMF_APP, ue_context_p->amf_context.imsi64,
        "UE STATE - CONNECTED.\n");
  } else if (
      ue_context_p->cm_state == M5GCM_IDLE && new_cm_state == M5GCM_IDLE) {
    OAILOG_INFO_UE(
        LOG_AMF_APP, ue_context_p->amf_context.imsi64,
        "Old UE CM State (IDLE) is same as the new UE CM state (IDLE)\n");
    magma::map_rc_t m_rc = magma::MAP_OK;
    m_rc = amf_ue_context_p->gnb_ue_ngap_id_ue_context_htbl.remove(
        ue_context_p->gnb_ngap_id_key);

    if (m_rc != magma::MAP_OK) {
      OAILOG_WARNING_UE(
          LOG_AMF_APP, ue_context_p->amf_context.imsi64,
          "UE context gnb_ue_ngap_ue_id_key %ld "
          "amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT
          ", GNB_UE_NGAP_ID_KEY could not be found",
          ue_context_p->gnb_ngap_id_key, ue_context_p->amf_ue_ngap_id);
    }

    ue_context_p->gnb_ngap_id_key = INVALID_GNB_UE_NGAP_ID_KEY;
  } else {
    OAILOG_INFO_UE(
        LOG_AMF_APP, ue_context_p->amf_context.imsi64,
        "Old UE CM State (CONNECTED) is same as the new UE CM state "
        "(CONNECTED)\n");
  }
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

// Context release timer expiry
static int amf_ue_context_release_complete_timer_handler(
    zloop_t* loop, int timer_id, void* output) {
  amf_context_t* amf_ctx = NULL;
  amf_ue_ngap_id_t ue_id = 0;

  if (!amf_pop_timer_arg(timer_id, &ue_id)) {
    OAILOG_WARNING(
        LOG_AMF_APP,
        "Context Release Timer: Invalid Timer Id expiration, Timer Id: %u\n",
        timer_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
  }

  ue_m5gmm_context_s* ue_amf_context =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  if (ue_amf_context == NULL) {
    OAILOG_ERROR(
        LOG_NAS_AMF,
        "Ue Context Release Timer: ue_amf_context is NULL for "
        "ue id: " AMF_UE_NGAP_ID_FMT "\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
  }

  amf_free_ue_context(ue_amf_context);
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
}

// Api for release UE context
void amf_ctx_release_ue_context(
    ue_m5gmm_context_s* ue_context_p, enum n2cause n2_cause) {
  if (!ue_context_p) {
    OAILOG_ERROR(LOG_AMF_APP, "Ue contex is null");
    OAILOG_FUNC_OUT(LOG_AMF_APP);
  }

  amf_app_itti_ue_context_release(
      ue_context_p, ue_context_p->ue_context_rel_cause);

  ue_context_p->m5_implicit_deregistration_timer.id = amf_app_start_timer(
      amf_config.nas_config.implicit_dereg_sec, TIMER_REPEAT_ONCE,
      amf_ue_context_release_complete_timer_handler,
      ue_context_p->amf_ue_ngap_id);
}

// Get the ue_context release cause
status_code_e amf_get_ue_context_rel_cause(
    amf_ue_ngap_id_t ue_id, enum n2cause* ue_context_rel_cause) {
  ue_m5gmm_context_s* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  if (ue_context_p == NULL) {
    return RETURNerror;
  }

  *ue_context_rel_cause = ue_context_p->ue_context_rel_cause;
  return RETURNok;
}

// Get the ue_context mm state
status_code_e amf_get_ue_context_mm_state(
    amf_ue_ngap_id_t ue_id, m5gmm_state_t* mm_state) {
  ue_m5gmm_context_s* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  if (ue_context_p == NULL) {
    return RETURNerror;
  }

  *mm_state = ue_context_p->mm_state;
  return RETURNok;
}

// Get the ue_context cm state
status_code_e amf_get_ue_context_cm_state(
    amf_ue_ngap_id_t ue_id, m5gcm_state_t* cm_state) {
  ue_m5gmm_context_s* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  if (ue_context_p == NULL) {
    return RETURNerror;
  }

  *cm_state = ue_context_p->cm_state;
  return RETURNok;
}

/* Get the ue id from IMSI */
bool get_amf_ue_id_from_imsi(amf_ue_context_t* amf_ue_context_p,
                             imsi64_t imsi64, amf_ue_ngap_id_t* ue_id) {
  magma::map_rc_t rc_map = magma::MAP_OK;
  rc_map = amf_ue_context_p->imsi_amf_ue_id_htbl.get(imsi64, ue_id);
  if (rc_map != magma::MAP_OK) {
    return (false);
  }
  return true;
}

}  // namespace magma5g
