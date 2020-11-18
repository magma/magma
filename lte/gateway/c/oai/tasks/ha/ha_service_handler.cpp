/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

#include <string.h>
#include <sys/types.h>

extern "C" {
#include "intertask_interface.h"
#include "common_types.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "ha_defs.h"
#include "ha_messages_types.h"
#include "s1ap_state.h"
#include "S1ap-CauseRadioNetwork.h"
}

#include "mme_app_state_manager.h"
#include "s1ap_state_manager.h"

static bool process_ue_context(
    const hash_key_t keyP, void* const elementP, void* parameterP,
    void** resultP);

bool handle_agw_offload_req(ha_agw_offload_req_t* offload_req) {
  hash_table_ts_t* state_imsi_ht =
      magma::lte::MmeNasStateManager::getInstance().get_ue_state_ht();
  s1ap_state_t* s1ap_state =
      magma::lte::S1apStateManager::getInstance().get_state(false);
  hashtable_ts_apply_callback_on_elements(
      state_imsi_ht, process_ue_context, (void*) s1ap_state, NULL);
}

bool process_ue_context(
    const hash_key_t keyP, void* const elementP, void* parameterP,
    void** resultP) {
  s1ap_state_t* s1ap_state             = (s1ap_state_t*) parameterP;
  struct ue_mm_context_s* ue_context_p = (struct ue_mm_context_s*) elementP;
  enb_description_t* enb_ref_p =
      s1ap_state_get_enb(s1ap_state, ue_context_p->sctp_assoc_id_key);
  MessageDef* message_p =
      itti_alloc_new_message(TASK_HA, S1AP_UE_CONTEXT_RELEASE_REQ);
  S1AP_UE_CONTEXT_RELEASE_REQ(message_p).mme_ue_s1ap_id =
      ue_context_p->mme_ue_s1ap_id;
  S1AP_UE_CONTEXT_RELEASE_REQ(message_p).enb_ue_s1ap_id =
      ue_context_p->enb_ue_s1ap_id;
  S1AP_UE_CONTEXT_RELEASE_REQ(message_p).enb_id   = enb_ref_p->enb_id;
  S1AP_UE_CONTEXT_RELEASE_REQ(message_p).relCause = S1AP_NAS_MME_OFFLOADING;

  OAILOG_INFO(
      LOG_UTIL,
      "Processing MME UE ID: %d, ENB UE ID: %d, UE Context ENB ID: %d, UE "
      "Context cell id: %d, S1AP State ENB ID: %d",
      ue_context_p->mme_ue_s1ap_id, ue_context_p->enb_ue_s1ap_id,
      ue_context_p->e_utran_cgi.cell_identity.enb_id,
      ue_context_p->e_utran_cgi.cell_identity.cell_id, enb_ref_p->enb_id);

  send_msg_to_task(&ha_task_zmq_ctx, TASK_MME_APP, message_p);
  return false;
}
