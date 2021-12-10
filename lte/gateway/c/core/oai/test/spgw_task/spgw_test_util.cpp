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

#include <cstdint>
#include <cstring>

extern "C" {
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.003.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.007.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_29.274.h"
#include "lte/gateway/c/core/oai/include/gx_messages_types.h"
#include "lte/gateway/c/core/oai/tasks/nas/api/mme/mme_api.h"
#include "lte/gateway/c/core/oai/include/s11_messages_types.h"
#include "lte/gateway/c/core/oai/common/itti_free_defined_msg.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
}

namespace magma {
namespace lte {

extern task_zmq_ctx_t task_zmq_ctx_main_spgw;

bool is_num_sessions_valid(
    uint64_t imsi64, int expected_num_ue_contexts, int expected_num_teids) {
  hash_table_ts_t* state_ue_ht = get_spgw_ue_state();
  if (state_ue_ht->num_elements != expected_num_ue_contexts) {
    return false;
  }
  if (expected_num_ue_contexts == 0) {
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
    return false;
  }
  return true;
}

bool is_num_s1_bearers_valid(
    teid_t context_teid, int expected_num_active_bearers) {
  s_plus_p_gw_eps_bearer_context_information_t* ctxt_p =
      sgw_cm_get_spgw_context(context_teid);
  if (ctxt_p == nullptr) {
    return false;
  }
  sgw_eps_bearer_context_information_t sgw_context_p =
      ctxt_p->sgw_eps_bearer_context_information;
  int num_active_bearers = 0;
  for (int ebx = 0; ebx < BEARERS_PER_UE; ebx++) {
    sgw_eps_bearer_ctxt_t* eps_bearer_ctxt =
        sgw_context_p.pdn_connection.sgw_eps_bearers_array[ebx];
    if ((eps_bearer_ctxt) &&
        (eps_bearer_ctxt->enb_ip_address_S1u.address.ipv4_address.s_addr !=
         0)) {
      num_active_bearers++;
    }
  }
  if (num_active_bearers == expected_num_active_bearers) {
    return true;
  }
  return false;
}

int get_num_pending_create_bearer_procedures(
    sgw_eps_bearer_context_information_t* ctxt_p) {
  if (ctxt_p == nullptr) {
    return 0;
  }

  int num_pending_create_procedures = 0;
  if (ctxt_p->pending_procedures) {
    pgw_base_proc_t* base_proc = NULL;

    LIST_FOREACH(base_proc, ctxt_p->pending_procedures, entries) {
      if (PGW_BASE_PROC_TYPE_NETWORK_INITATED_CREATE_BEARER_REQUEST ==
          base_proc->type) {
        num_pending_create_procedures++;
      }
    }
  }
  return num_pending_create_procedures;
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
  modify_bearer_req->bearer_contexts_to_be_modified.bearer_contexts[bearer_id]
      .s1_eNB_fteid.ipv4_address.s_addr = DEFAULT_ENB_IP;

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

void fill_release_access_bearer_request(
    itti_s11_release_access_bearers_request_t* release_access_bearers_req,
    teid_t mme_s11_teid, teid_t sgw_s11_context_teid) {
  release_access_bearers_req->local_teid = mme_s11_teid;
  release_access_bearers_req->teid       = sgw_s11_context_teid;
  release_access_bearers_req->edns_peer_ip.addr_v4.sin_addr.s_addr =
      DEFAULT_EDNS_IP;
  release_access_bearers_req->edns_peer_ip.addr_v4.sin_family = AF_INET;
  release_access_bearers_req->originating_node                = NODE_TYPE_MME;
}

void fill_packet_filter_content(packet_filter_contents_t* pf_content) {
  // TODO : Parameterize the protocol, IP Address and port numbers
  pf_content->flags = TRAFFIC_FLOW_TEMPLATE_PROTOCOL_NEXT_HEADER_FLAG |
                      TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG |
                      TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT_FLAG;

  pf_content->protocolidentifier_nextheader = IPPROTO_TCP;

  // iPerf server port
  pf_content->singleremoteport = 5001;

  // Remote address as 192.168.129.42/24
  pf_content->ipv4remoteaddr[0].addr = 192;
  pf_content->ipv4remoteaddr[1].addr = 168;
  pf_content->ipv4remoteaddr[2].addr = 129;
  pf_content->ipv4remoteaddr[3].addr = 42;
  pf_content->ipv4remoteaddr[0].mask = 255;
  pf_content->ipv4remoteaddr[1].mask = 255;
  pf_content->ipv4remoteaddr[2].mask = 255;
  pf_content->ipv4remoteaddr[3].mask = 0;
}

void fill_nw_initiated_activate_bearer_request(
    itti_gx_nw_init_actv_bearer_request_t* gx_nw_init_actv_req_p,
    const std::string& imsi_str, ebi_t lbi, bearer_qos_t qos) {
  gx_nw_init_actv_req_p->imsi_length = imsi_str.size();
  strncpy(gx_nw_init_actv_req_p->imsi, imsi_str.c_str(), imsi_str.size());
  gx_nw_init_actv_req_p->lbi            = lbi;
  gx_nw_init_actv_req_p->eps_bearer_qos = qos;

  strncpy(
      gx_nw_init_actv_req_p->policy_rule_name, DEFAULT_POLICY_RULE_NAME,
      DEFAULT_POLICY_RULE_NAME_LEN);
  gx_nw_init_actv_req_p->policy_rule_name[DEFAULT_POLICY_RULE_NAME_LEN] = '\0';
  gx_nw_init_actv_req_p->policy_rule_name_length = DEFAULT_POLICY_RULE_NAME_LEN;

  traffic_flow_template_t* ul_tft = &gx_nw_init_actv_req_p->ul_tft;
  traffic_flow_template_t* dl_tft = &gx_nw_init_actv_req_p->dl_tft;
  memset(ul_tft, 0, sizeof(traffic_flow_template_t));
  memset(dl_tft, 0, sizeof(traffic_flow_template_t));

  ul_tft->tftoperationcode = TRAFFIC_FLOW_TEMPLATE_OPCODE_CREATE_NEW_TFT;
  dl_tft->tftoperationcode = TRAFFIC_FLOW_TEMPLATE_OPCODE_CREATE_NEW_TFT;
  ul_tft->ebit = TRAFFIC_FLOW_TEMPLATE_PARAMETER_LIST_IS_NOT_INCLUDED;
  dl_tft->ebit = TRAFFIC_FLOW_TEMPLATE_PARAMETER_LIST_IS_NOT_INCLUDED;

  // create one uplink tft
  ul_tft->numberofpacketfilters = 1;
  ul_tft->packetfilterlist.createnewtft[0].direction =
      TRAFFIC_FLOW_TEMPLATE_UPLINK_ONLY;
  ul_tft->packetfilterlist.createnewtft[0].eval_precedence = qos.pl;
  fill_packet_filter_content(
      &ul_tft->packetfilterlist.createnewtft[0].packetfiltercontents);

  // create one downlink tft
  dl_tft->numberofpacketfilters = 1;
  dl_tft->packetfilterlist.createnewtft[0].direction =
      TRAFFIC_FLOW_TEMPLATE_DOWNLINK_ONLY;
  dl_tft->packetfilterlist.createnewtft[0].eval_precedence = qos.pl;
  fill_packet_filter_content(
      &dl_tft->packetfilterlist.createnewtft[0].packetfiltercontents);
}

void fill_nw_initiated_activate_bearer_response(
    itti_s11_nw_init_actv_bearer_rsp_t* nw_actv_bearer_resp,
    teid_t mme_s11_teid, teid_t sgw_s11_cp_teid, teid_t sgw_s11_ded_teid,
    teid_t s1u_enb_ded_teid, ebi_t eps_bearer_id, gtpv2c_cause_value_t cause,
    plmn_t plmn) {
  nw_actv_bearer_resp->sgw_s11_teid = sgw_s11_cp_teid;
  COPY_PLMN_IN_ARRAY_FMT(nw_actv_bearer_resp->serving_network, plmn);
  nw_actv_bearer_resp->cause.cause_value = cause;

  int msg_bearer_index = 0;
  nw_actv_bearer_resp->bearer_contexts.bearer_contexts[msg_bearer_index]
      .eps_bearer_id = eps_bearer_id;
  nw_actv_bearer_resp->bearer_contexts.bearer_contexts[msg_bearer_index]
      .cause.cause_value = REQUEST_ACCEPTED;

  // Fill eNB S1u Fteid with new teid for dedicated bearer
  nw_actv_bearer_resp->bearer_contexts.bearer_contexts[msg_bearer_index]
      .s1u_enb_fteid = {.ipv4           = true,
                        .interface_type = S1_U_ENODEB_GTP_U,
                        .teid           = s1u_enb_ded_teid,
                        .ipv4_address   = {.s_addr = DEFAULT_ENB_IP}};

  // Fill SGW S1u Fteid
  nw_actv_bearer_resp->bearer_contexts.bearer_contexts[msg_bearer_index]
      .s1u_sgw_fteid = {.ipv4           = true,
                        .interface_type = S1_U_SGW_GTP_U,
                        .teid           = sgw_s11_ded_teid,
                        .ipv4_address   = {.s_addr = DEFAULT_SGW_IP}};

  nw_actv_bearer_resp->bearer_contexts.num_bearer_context = 1;
}

void fill_nw_initiated_deactivate_bearer_request(
    itti_gx_nw_init_deactv_bearer_request_t* gx_nw_init_deactv_req_p,
    const std::string& imsi_str, ebi_t lbi, ebi_t eps_bearer_id) {
  gx_nw_init_deactv_req_p->imsi_length = imsi_str.size();
  strncpy(gx_nw_init_deactv_req_p->imsi, imsi_str.c_str(), imsi_str.size());
  gx_nw_init_deactv_req_p->lbi           = lbi;
  gx_nw_init_deactv_req_p->no_of_bearers = 1;
  gx_nw_init_deactv_req_p->ebi[0]        = eps_bearer_id;
}

void fill_nw_initiated_deactivate_bearer_response(
    itti_s11_nw_init_deactv_bearer_rsp_t* nw_deactv_bearer_resp,
    uint64_t test_imsi64, bool delete_default_bearer,
    gtpv2c_cause_value_t cause, ebi_t ebi[], unsigned int num_bearer_context,
    teid_t sgw_s11_context_teid) {
  nw_deactv_bearer_resp->delete_default_bearer = delete_default_bearer;
  nw_deactv_bearer_resp->cause.cause_value     = cause;

  if (delete_default_bearer) {
    nw_deactv_bearer_resp->lbi =
        reinterpret_cast<ebi_t*>(calloc(1, sizeof(ebi_t)));
    *nw_deactv_bearer_resp->lbi = ebi[0];
    nw_deactv_bearer_resp->bearer_contexts.bearer_contexts[0]
        .cause.cause_value = cause;
  } else {
    for (int i = 0; i < num_bearer_context; i++) {
      nw_deactv_bearer_resp->bearer_contexts.bearer_contexts[i].eps_bearer_id =
          ebi[i];
      nw_deactv_bearer_resp->bearer_contexts.bearer_contexts[i]
          .cause.cause_value = cause;
    }
  }
  nw_deactv_bearer_resp->bearer_contexts.num_bearer_context =
      num_bearer_context;
  nw_deactv_bearer_resp->imsi             = test_imsi64;
  nw_deactv_bearer_resp->s_gw_teid_s11_s4 = sgw_s11_context_teid;
}

void fill_s11_suspend_notification(
    itti_s11_suspend_notification_t* suspend_notif, teid_t sgw_s11_context_teid,
    const std::string& imsi_str, ebi_t link_bearer_id) {
  suspend_notif->teid        = sgw_s11_context_teid;
  suspend_notif->lbi         = link_bearer_id;
  suspend_notif->imsi.length = imsi_str.size();
  strncpy(
      (char*) suspend_notif->imsi.digit, imsi_str.c_str(),
      suspend_notif->imsi.length);
}

}  // namespace lte
}  // namespace magma
