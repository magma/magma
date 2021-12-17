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
#include "lte/gateway/c/core/oai/tasks/amf/include/amf_ue_context_storage.h"

namespace magma5g {
extern task_zmq_ctx_t amf_app_task_zmq_ctx;

auto& context_store = AmfUeContextStorage::getUeContextStorage();

std::shared_ptr<smf_context_t> amf_insert_smf_context(
    std::shared_ptr<ue_m5gmm_context_t>, uint8_t);

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
    const std::shared_ptr<ue_m5gmm_context_t> ue_context_p) {
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
    auto ue_mm_context = amf_get_ue_context(ue_id);

    if (ue_mm_context) {
      amf_context_p = &ue_mm_context->amf_context;
    }
    OAILOG_DEBUG(LOG_NAS_AMF, "Stored UE id " AMF_UE_NGAP_ID_FMT " \n", ue_id);
  }
  return amf_context_p;
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
    std::shared_ptr<ue_m5gmm_context_t> ue_context, uint8_t pdu_session_id) {
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
    std::shared_ptr<ue_m5gmm_context_t> ue_context, uint8_t id) {
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
  auto ue_context_p      = amf_ctx->ue_context_p.lock();
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
 ** Description: Deletes the ue context information from all maps          **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
void amf_free_ue_context(std::shared_ptr<ue_m5gmm_context_t> ue_context_p) {
  hashtable_rc_t h_rc            = HASH_TABLE_OK;
  magma::map_rc_t m_rc           = magma::MAP_OK;
  amf_app_desc_t* amf_app_desc_p = get_amf_nas_state(false);
  OAILOG_DEBUG(LOG_NAS_AMF, "amf_free_ue_context \n");
  if (!ue_context_p) {
    return;
  }

  // Clean up the procedures
  amf_nas_proc_clean_up(ue_context_p);
  AmfUeContextStorage::getUeContextStorage().amf_remove_ue_context_from_cache(ue_context_p);
}

}  // namespace magma5g
