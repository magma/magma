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

#include "lte/gateway/c/core/oai/test/spgw_task/spgw_test_util.h"
#include <iostream>

#include <cstdint>

extern "C" {
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.003.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.007.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_29.274.h"
#include "lte/gateway/c/core/oai/tasks/nas/api/mme/mme_api.h"
#include "lte/gateway/c/core/oai/include/s11_messages_types.h"
#include "lte/gateway/c/core/oai/common/itti_free_defined_msg.h"
}

namespace magma {
namespace lte {

extern task_zmq_ctx_t task_zmq_ctx_main_spgw;

bool is_num_sessions_valid(
    spgw_state_t* spgw_state, uint64_t imsi64, int expected_num_ue_contexts,
    int expected_num_teids) {
  hash_table_ts_t* state_ue_ht = get_spgw_ue_state();
  if (state_ue_ht->num_elements != expected_num_ue_contexts) {
    std::cout << "is_num_sessions_valid: false 1: " << "real: " << state_ue_ht->num_elements << " | expected: " << expected_num_ue_contexts << std::endl;
    return false;
  }
  if (expected_num_ue_contexts == 0) {
    std::cout << "is_num_sessions_valid: true 1" << std::endl;
    return true;
  }

  // If one UE context exists, check that teids also exist in hashtable
  spgw_ue_context_t* ue_context_p = spgw_get_ue_context(imsi64);
  int num_teids                   = 0;
  sgw_s11_teid_t* s11_teid_p      = nullptr;
  LIST_FOREACH(s11_teid_p, &ue_context_p->sgw_s11_teid_list, entries) {
    if (s11_teid_p &&
        (sgw_cm_get_spgw_context(s11_teid_p->sgw_s11_teid) != nullptr)) {
      num_teids++;
    }
  }
  if (num_teids != expected_num_teids) {
    std::cout << "is_num_sessions_valid: real: " << num_teids << " | expected: " << expected_num_teids << std::endl;
    return false;
  }
  return true;
}

void fill_create_session_request(
    itti_s11_create_session_request_t* session_request_p,
    const std::string& imsi_str, teid_t mme_s11_teid, int bearer_id,
    bearer_context_to_be_created_t sample_bearer_context, plmn_t sample_plmn) {
  session_request_p->teid = 0;
  strncpy(
      (char*) session_request_p->imsi.digit, imsi_str.c_str(), imsi_str.size());
  session_request_p->imsi.length                        = imsi_str.size();
  session_request_p->sender_fteid_for_cp.teid           = mme_s11_teid;
  session_request_p->sender_fteid_for_cp.interface_type = S11_MME_GTP_C;

  session_request_p->uli.present = 0;
  session_request_p->rat_type    = RAT_EUTRAN;

  session_request_p->bearer_contexts_to_be_created.bearer_contexts[bearer_id]
      .eps_bearer_id = sample_bearer_context.eps_bearer_id;
  session_request_p->bearer_contexts_to_be_created.bearer_contexts[bearer_id]
      .bearer_level_qos.pci = sample_bearer_context.bearer_level_qos.pci;
  session_request_p->bearer_contexts_to_be_created.bearer_contexts[bearer_id]
      .bearer_level_qos.pl = sample_bearer_context.bearer_level_qos.pl;
  session_request_p->bearer_contexts_to_be_created.bearer_contexts[bearer_id]
      .bearer_level_qos.pvi = sample_bearer_context.bearer_level_qos.pvi;
  session_request_p->bearer_contexts_to_be_created.bearer_contexts[bearer_id]
      .bearer_level_qos.qci = sample_bearer_context.bearer_level_qos.qci;
  session_request_p->bearer_contexts_to_be_created.bearer_contexts[bearer_id]
      .bearer_level_qos.mbr.br_ul =
      sample_bearer_context.bearer_level_qos.mbr.br_ul;
  session_request_p->bearer_contexts_to_be_created.bearer_contexts[bearer_id]
      .bearer_level_qos.mbr.br_dl =
      sample_bearer_context.bearer_level_qos.mbr.br_dl;
  session_request_p->bearer_contexts_to_be_created.num_bearer_context = 1;

  session_request_p->sender_fteid_for_cp.teid           = (teid_t) 1;
  session_request_p->sender_fteid_for_cp.interface_type = S11_MME_GTP_C;
  session_request_p->sender_fteid_for_cp.ipv4_address.s_addr =
      0xc0a83c8e;  // 192.168.60.142
  session_request_p->sender_fteid_for_cp.ipv4 = 1;

  const char default_apn[] = "magma.ipv4";
  strncpy(session_request_p->apn, default_apn, 10);
  session_request_p->ambr.br_dl = 100000000;
  session_request_p->ambr.br_ul = 200000000;

  session_request_p->pdn_type                = IPv4;
  session_request_p->paa.pdn_type            = IPv4;
  session_request_p->paa.ipv4_address.s_addr = INADDR_ANY;
  session_request_p->paa.ipv6_address        = in6addr_any;

  session_request_p->serving_network.mcc[0] = sample_plmn.mcc_digit1;
  session_request_p->serving_network.mcc[1] = sample_plmn.mcc_digit2;
  session_request_p->serving_network.mcc[2] = sample_plmn.mcc_digit3;
  session_request_p->serving_network.mnc[0] = sample_plmn.mnc_digit1;
  session_request_p->serving_network.mnc[1] = sample_plmn.mnc_digit2;
  session_request_p->serving_network.mnc[2] = sample_plmn.mnc_digit3;
}

void fill_ip_allocation_response(
    itti_ip_allocation_response_t* ip_alloc_resp_p, SGIStatus_t status,
    teid_t sgw_s11_context_teid, ebi_t eps_bearer_id, unsigned long ue_ip,
    int vlan) {
  ip_alloc_resp_p->status                  = status;
  ip_alloc_resp_p->context_teid            = sgw_s11_context_teid;
  ip_alloc_resp_p->eps_bearer_id           = eps_bearer_id;
  ip_alloc_resp_p->paa.ipv4_address.s_addr = ue_ip;
  ip_alloc_resp_p->paa.pdn_type            = IPv4;
  ip_alloc_resp_p->paa.vlan                = vlan;
}

void fill_pcef_create_session_response(
    itti_pcef_create_session_response_t* pcef_csr_resp_p,
    PcefRpcStatus_t rpc_status, teid_t sgw_s11_context_teid,
    ebi_t eps_bearer_id, SGIStatus_t sgi_status) {
  pcef_csr_resp_p->rpc_status    = rpc_status;
  pcef_csr_resp_p->teid          = sgw_s11_context_teid;
  pcef_csr_resp_p->eps_bearer_id = eps_bearer_id;
  pcef_csr_resp_p->sgi_status    = sgi_status;
}

void fill_modify_bearer_request(
    itti_s11_modify_bearer_request_t* modify_bearer_req, teid_t mme_s11_teid,
    teid_t sgw_s11_context_teid, teid_t enb_gtp_teid, int bearer_id,
    ebi_t eps_bearer_id) {
  modify_bearer_req->local_teid                = mme_s11_teid;
  modify_bearer_req->delay_dl_packet_notif_req = 0;
  modify_bearer_req->bearer_contexts_to_be_modified.bearer_contexts[bearer_id]
      .eps_bearer_id = eps_bearer_id;

  modify_bearer_req->edns_peer_ip.addr_v4.sin_addr.s_addr = DEFAULT_EDNS_IP;

  modify_bearer_req->edns_peer_ip.addr_v4.sin_family = AF_INET;

  modify_bearer_req->teid = sgw_s11_context_teid;

  // populate the eNB FTEID
  modify_bearer_req->bearer_contexts_to_be_modified.bearer_contexts[bearer_id]
      .s1_eNB_fteid.teid = enb_gtp_teid;
  modify_bearer_req->bearer_contexts_to_be_modified.bearer_contexts[bearer_id]
      .s1_eNB_fteid.interface_type = S1_U_ENODEB_GTP_U;
  modify_bearer_req->bearer_contexts_to_be_modified.bearer_contexts[bearer_id]
      .s1_eNB_fteid.ipv4 = 1;

  // Only one bearer context to be sent for default PDN
  modify_bearer_req->bearer_contexts_to_be_modified.num_bearer_context = 1;
  modify_bearer_req->bearer_contexts_to_be_removed.num_bearer_context  = 0;
  modify_bearer_req->mme_fq_csid.node_id_type = GLOBAL_UNICAST_IPv4;
  modify_bearer_req->mme_fq_csid.csid         = 0;
  memset(
      &modify_bearer_req->indication_flags, 0,
      sizeof(modify_bearer_req->indication_flags));
  modify_bearer_req->rat_type = RAT_EUTRAN;
}

void fill_delete_session_request(
    itti_s11_delete_session_request_t* delete_session_req, teid_t mme_s11_teid,
    teid_t sgw_s11_context_teid, ebi_t eps_bearer_id, plmn_t test_plmn) {
  delete_session_req->local_teid = mme_s11_teid;
  delete_session_req->teid       = sgw_s11_context_teid;
  delete_session_req->noDelete   = true;
  delete_session_req->lbi        = eps_bearer_id;

  // EDNS address
  delete_session_req->edns_peer_ip.addr_v4.sin_family      = AF_INET;
  delete_session_req->edns_peer_ip.addr_v4.sin_addr.s_addr = DEFAULT_EDNS_IP;

  // Sender FTEID
  delete_session_req->sender_fteid_for_cp.teid           = mme_s11_teid;
  delete_session_req->sender_fteid_for_cp.interface_type = S11_MME_GTP_C;
  delete_session_req->sender_fteid_for_cp.ipv4           = 1;

  delete_session_req->indication_flags.oi = 1;
  delete_session_req->peer_ip.s_addr      = DEFAULT_SGW_IP;
  delete_session_req->trxn                = nullptr;

  // PLMN
  COPY_PLMN_IN_ARRAY_FMT(delete_session_req->serving_network, test_plmn);
}

}  // namespace lte
}  // namespace magma
