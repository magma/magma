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
#include <iostream>
#include <string.h>
#include <sys/types.h>

extern "C" {
#include "common_types.h"
#include "ha_defs.h"
#include "ha_messages_types.h"
#include "intertask_interface.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "log.h"
#include "s1ap_state.h"
#include "S1ap_CauseRadioNetwork.h"
}

#include "HaClient.h"
#include "mme_app_state_manager.h"
#include "s1ap_state_manager.h"

static bool trigger_agw_offload_for_ue(
    const hash_key_t keyP, void* const elementP, void* parameterP,
    void** resultP);

bool sync_up_with_orc8r(void) {
  magma::HaClient::get_eNB_offload_state(
      [](grpc::Status status,
         magma::lte::GetEnodebOffloadStateResponse response) {
        if (status.ok()) {
          OAILOG_INFO(
              LOG_UTIL, "Received eNodeB connection state with the primary.");
          // iterate over the eNodeB connection states
          ha_agw_offload_req_t offload_req = {0};
          for (auto const& item : response.enodeb_offload_states()) {
            if (item.second ==
                magma::lte::GetEnodebOffloadStateResponse::PRIMARY_CONNECTED) {
              offload_req.eNB_id = item.first;
              // Offload any UE to check if it can be camped on the primary.
              // The effect will be observed in the next sync up with the cloud.
              offload_req.enb_offload_type = ANY;
              handle_agw_offload_req(&offload_req);
            } else if (
                item.second == magma::lte::GetEnodebOffloadStateResponse::
                                   PRIMARY_CONNECTED_AND_SERVING_UES) {
              offload_req.eNB_id = item.first;
              // Primary looks healthy as UEs are camped on it, offload the rest
              // of UEs.
              offload_req.enb_offload_type = ALL;
              handle_agw_offload_req(&offload_req);
            }
          }
        } else {
          OAILOG_ERROR(
              LOG_UTIL, "GRPC Failure Message: %s Status Error Code: %d",
              status.error_message().c_str(), status.error_code());
        }
      });
  return true;
}

typedef struct callback_data_s {
  s1ap_state_t* s1ap_state;
  ha_agw_offload_req_t* request;
} callback_data_t;

void handle_agw_offload_req(ha_agw_offload_req_t* offload_req) {
  hash_table_ts_t* state_imsi_ht =
      magma::lte::MmeNasStateManager::getInstance().get_ue_state_ht();
  callback_data_t callback_data;
  callback_data.s1ap_state =
      magma::lte::S1apStateManager::getInstance().get_state(false);
  callback_data.request = offload_req;
  hashtable_ts_apply_callback_on_elements(
      state_imsi_ht, trigger_agw_offload_for_ue, (void*) &callback_data, NULL);
}

bool trigger_agw_offload_for_ue(
    const hash_key_t keyP, void* const elementP, void* parameterP,
    void** resultP) {
  imsi64_t imsi64                = INVALID_IMSI64;
  callback_data_t* callback_data = (callback_data_t*) parameterP;
  ha_agw_offload_req_t* offload_request =
      (ha_agw_offload_req_t*) callback_data->request;
  s1ap_state_t* s1ap_state = (s1ap_state_t*) callback_data->s1ap_state;
  struct ue_mm_context_s* ue_context_p = (struct ue_mm_context_s*) elementP;
  bool any_flag = false;  // true if we tried offloading any UE

  IMSI_STRING_TO_IMSI64(offload_request->imsi, &imsi64);

  enb_description_t* enb_ref_p =
      s1ap_state_get_enb(s1ap_state, ue_context_p->sctp_assoc_id_key);

  // Return if this UE does not satisfy any of the filtering criteria
  if ((imsi64 != ue_context_p->emm_context._imsi64) &&
      (offload_request->eNB_id != enb_ref_p->enb_id)) {
    return false;
  }

  offload_type_t enb_offtype = offload_request->enb_offload_type;
  // When a UE is in ECM_CONNECTED state, we can direcly start offloading.
  // For a UE in ECM_IDLE mode however, we need to first page the user and
  // then we can offload it.
  if ((ue_context_p->ecm_state == ECM_CONNECTED) &&
      ((enb_offtype == ALL) || (enb_offtype == ANY) ||
       (enb_offtype == ANY_CONNECTED))) {
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
        "Processing IMSI64: " IMSI_64_FMT
        " Requested IMSI: %s, MME UE ID: %d, ENB UE ID: %d, UE "
        "Context ENB ID: "
        "%d, UE "
        "Context cell id: %d, S1AP State ENB ID: %d",
        ue_context_p->emm_context._imsi64, offload_request->imsi,
        ue_context_p->mme_ue_s1ap_id, ue_context_p->enb_ue_s1ap_id,
        ue_context_p->e_utran_cgi.cell_identity.enb_id,
        ue_context_p->e_utran_cgi.cell_identity.cell_id, enb_ref_p->enb_id);
    OAILOG_INFO(
        LOG_UTIL, "UE Context Release procedure initiated for IMSI%s",
        offload_request->imsi);
    IMSI_STRING_TO_IMSI64(
        offload_request->imsi, &message_p->ittiMsgHeader.imsi);
    send_msg_to_task(&ha_task_zmq_ctx, TASK_MME_APP, message_p);
    any_flag = true;
  } else if (
      (ue_context_p->ecm_state == ECM_IDLE) &&
      (ue_context_p->mm_state == UE_REGISTERED) &&
      ((enb_offtype == ALL) || (enb_offtype == ANY) ||
       (enb_offtype == ANY_IDLE))) {
    // Upon connection re-establishment, this release cause value will
    // be checked and cleared by MME APP to send offload request.
    ue_context_p->ue_context_rel_cause = S1AP_NAS_MME_PENDING_OFFLOADING;

    char imsi[IMSI_BCD_DIGITS_MAX + 1] = {0};
    IMSI64_TO_STRING(
        ue_context_p->emm_context._imsi64, imsi,
        ue_context_p->emm_context._imsi.length);

    OAILOG_INFO(LOG_UTIL, "Paging procedure initiated for IMSI%s", imsi);
    MessageDef* message_p                       = NULL;
    itti_s11_paging_request_t* paging_request_p = NULL;

    message_p        = itti_alloc_new_message(TASK_HA, S11_PAGING_REQUEST);
    paging_request_p = &message_p->ittiMsg.s11_paging_request;
    memset((void*) paging_request_p, 0, sizeof(itti_s11_paging_request_t));
    paging_request_p->imsi        = strdup(imsi);
    message_p->ittiMsgHeader.imsi = ue_context_p->emm_context._imsi64;
    send_msg_to_task(&ha_task_zmq_ctx, TASK_MME_APP, message_p);
    any_flag = true;
  }

  // Check if iterations should be stopped as single match was
  // sufficient.
  if (any_flag && ((enb_offtype == ANY) || (enb_offtype == ANY_CONNECTED) ||
                   (enb_offtype == ANY_IDLE))) {
    return true;
  }
  return false;
}
