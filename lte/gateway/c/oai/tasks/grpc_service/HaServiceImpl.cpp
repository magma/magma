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
#include <string>

#include "HaServiceImpl.h"
#include "lte/protos/ha_service.pb.h"
#include "mme_app_state_manager.h"
#include "s1ap_state_manager.h"
extern "C" {
#include "log.h"
#include "intertask_interface.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "s1ap_state.h"
#include "S1ap-CauseRadioNetwork.h"
// #include "s1ap_messages_types.h"
}
extern task_zmq_ctx_t grpc_service_task_zmq_ctx;

// namespace grpc {
// class ServerContext;
// }  // namespace grpc

// using grpc::ServerContext;
// using grpc::Status;
// using magma::HaService;
// using magma::StartAgwOffloadRequest;
// using magma::StartAgwOffloadResponse;

namespace magma {
using namespace lte;

HaServiceImpl::HaServiceImpl() {}
/*
 * StartAgwOffload is called by North Bound to release UEs
 * based on eNB identification.
 */
grpc::Status HaServiceImpl::StartAgwOffload(
    grpc::ServerContext* context, const StartAgwOffloadRequest* request,
    StartAgwOffloadResponse* response) {
  OAILOG_INFO(LOG_UTIL, "Received StartAgwOffloadRequest GRPC request\n");
  hash_table_ts_t* state_imsi_ht =
      MmeNasStateManager::getInstance().get_ue_state_ht();
  s1ap_state_t* s1ap_state = S1apStateManager::getInstance().get_state(false);
  hashtable_ts_apply_callback_on_elements(
      state_imsi_ht, process_ue_context, (void*) s1ap_state, NULL);
  return grpc::Status::OK;
}

}  // namespace magma

bool process_ue_context(
    const hash_key_t keyP, void* const elementP, void* parameterP,
    void** resultP) {
  s1ap_state_t* s1ap_state             = (s1ap_state_t*) parameterP;
  struct ue_mm_context_s* ue_context_p = (struct ue_mm_context_s*) elementP;
  enb_description_t* enb_ref_p =
      s1ap_state_get_enb(s1ap_state, ue_context_p->sctp_assoc_id_key);
  MessageDef* message_p =
      itti_alloc_new_message(TASK_GRPC_SERVICE, S1AP_UE_CONTEXT_RELEASE_REQ);
  S1AP_UE_CONTEXT_RELEASE_REQ(message_p).mme_ue_s1ap_id =
      ue_context_p->mme_ue_s1ap_id;
  S1AP_UE_CONTEXT_RELEASE_REQ(message_p).enb_ue_s1ap_id =
      ue_context_p->enb_ue_s1ap_id;
  S1AP_UE_CONTEXT_RELEASE_REQ(message_p).enb_id   = enb_ref_p->enb_id;
  S1AP_UE_CONTEXT_RELEASE_REQ(message_p).relCause = S1AP_NAS_MME_OFFLOADING;

  OAILOG_INFO(
      LOG_UTIL,
      "Processing MME UE ID: %d, ENB UE ID: %d, ENB ID: %d, ENB ID: %d",
      ue_context_p->mme_ue_s1ap_id, ue_context_p->enb_ue_s1ap_id,
      ue_context_p->e_utran_cgi.cell_identity.enb_id, enb_ref_p->enb_id);

  send_msg_to_task(&grpc_service_task_zmq_ctx, TASK_MME_APP, message_p);
  return false;
}
