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

#include "sgw_s8_utility.h"
#include "ControllerMain.h"

const struct gtp_tunnel_ops* gtp_tunnel_ops;
void fill_imsi(char* imsi) {
  uint8_t idx = 0;
  imsi[idx++] = '0';
  imsi[idx++] = '0';
  imsi[idx++] = '1';
  imsi[idx++] = '0';
  imsi[idx++] = '1';
  imsi[idx++] = '0';
  imsi[idx++] = '0';
  imsi[idx++] = '0';
  imsi[idx++] = '0';
  imsi[idx++] = '0';
  imsi[idx++] = '0';
  imsi[idx++] = '0';
  imsi[idx++] = '0';
  imsi[idx++] = '0';
  imsi[idx++] = '1';
}

void fill_itti_csreq(
    itti_s11_create_session_request_t* session_req_pP,
    uint8_t default_eps_bearer_id) {
  uint8_t idx = 0;
  fill_imsi((reinterpret_cast<char*>(session_req_pP->imsi.digit)));
  session_req_pP->teid                                    = 0;
  session_req_pP->imsi.length                             = 15;
  idx                                                     = 0;
  session_req_pP->serving_network.mcc[idx++]              = 0;
  session_req_pP->serving_network.mcc[idx++]              = 0;
  session_req_pP->serving_network.mcc[idx]                = 0;
  idx                                                     = 0;
  session_req_pP->serving_network.mnc[idx++]              = 1;
  session_req_pP->serving_network.mnc[idx++]              = 1;
  session_req_pP->serving_network.mnc[idx]                = 15;
  session_req_pP->rat_type                                = RAT_EUTRAN;
  session_req_pP->sender_fteid_for_cp.teid                = 1;
  session_req_pP->sender_fteid_for_cp.ipv4_address.s_addr = 0x8e3ca8c0;
  session_req_pP->sender_fteid_for_cp.interface_type      = S11_MME_GTP_C;

  session_req_pP->default_ebi = 5;
  bearer_contexts_to_be_created_t* bc_to_be_created =
      &session_req_pP->bearer_contexts_to_be_created;
  bc_to_be_created->num_bearer_context               = 1;
  bc_to_be_created->bearer_contexts[0].eps_bearer_id = 5;
}

void fill_itti_csrsp(s8_create_session_response_t* csr_resp, uint32_t teid) {
  uint8_t idx = 0;
  fill_imsi((reinterpret_cast<char*>(csr_resp->imsi)));
  csr_resp->imsi_length = 15;

  csr_resp->pdn_type                = IPv4;
  csr_resp->paa.pdn_type            = IPv4;
  csr_resp->paa.ipv4_address.s_addr = 0xc0a87e1;
  csr_resp->context_teid            = teid;
  csr_resp->eps_bearer_id           = 5;

  csr_resp->bearer_context[0].eps_bearer_id                 = 5;
  csr_resp->bearer_context[0].pgw_s8_up.ipv4                = 1;
  csr_resp->bearer_context[0].pgw_s8_up.interface_type      = S5_S8_PGW_GTP_U;
  csr_resp->bearer_context[0].pgw_s8_up.teid                = 123;
  csr_resp->bearer_context[0].pgw_s8_up.ipv4_address.s_addr = 0xc0a87e19;

  csr_resp->pgw_s8_cp_teid.ipv4                = 1;
  csr_resp->pgw_s8_cp_teid.interface_type      = S5_S8_PGW_GTP_C;
  csr_resp->pgw_s8_cp_teid.teid                = 124;
  csr_resp->pgw_s8_cp_teid.ipv4_address.s_addr = 0xc0a87e20;

  csr_resp->cause = 16;
}

void fill_create_bearer_request(
    s8_create_bearer_request_t* cb_req, uint32_t teid,
    uint8_t default_eps_bearer_id) {
#define IPV4_LEN 4
  cb_req->sequence_number      = 10;
  cb_req->context_teid         = teid;
  cb_req->linked_eps_bearer_id = default_eps_bearer_id;
  cb_req->pgw_cp_address = reinterpret_cast<char*>(calloc(1, IPV4_LEN + 1));
  cb_req->bearer_context[0].eps_bearer_id                 = 0;
  cb_req->bearer_context[0].pgw_s8_up.ipv4                = 1;
  cb_req->bearer_context[0].pgw_s8_up.interface_type      = S5_S8_PGW_GTP_U;
  cb_req->bearer_context[0].pgw_s8_up.teid                = 20;
  cb_req->bearer_context[0].pgw_s8_up.ipv4_address.s_addr = 0xac101496;
  cb_req->bearer_context[0].tft.tftoperationcode =
      TRAFFIC_FLOW_TEMPLATE_OPCODE_CREATE_NEW_TFT;
  cb_req->bearer_context[0].tft.numberofpacketfilters = 1;
  packet_filter_t* create_tft_packet_filter =
      &(cb_req->bearer_context[0].tft.packetfilterlist.createnewtft[0]);
  create_tft_packet_filter->direction       = TRAFFIC_FLOW_TEMPLATE_UPLINK_ONLY;
  create_tft_packet_filter->identifier      = 1;
  create_tft_packet_filter->eval_precedence = 200;
  create_tft_packet_filter->length          = 14;
  create_tft_packet_filter->packetfiltercontents.flags =
      TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG;
  for (uint8_t idx = 0; idx < TRAFFIC_FLOW_TEMPLATE_IPV4_ADDR_SIZE; idx++) {
    create_tft_packet_filter->packetfiltercontents.ipv4remoteaddr[0].addr =
        172 + idx;
    create_tft_packet_filter->packetfiltercontents.ipv4remoteaddr[0].mask = 255;
  }
  create_tft_packet_filter->packetfiltercontents.flags |=
      TRAFFIC_FLOW_TEMPLATE_PROTOCOL_NEXT_HEADER_FLAG;
  create_tft_packet_filter->packetfiltercontents.protocolidentifier_nextheader =
      17;
  create_tft_packet_filter->packetfiltercontents.flags |=
      TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT_FLAG;
  create_tft_packet_filter->packetfiltercontents.singleremoteport = 19166;
}

void fill_create_bearer_response(
    itti_s11_nw_init_actv_bearer_rsp_t* cb_response, uint32_t teid,
    uint8_t eps_bearer_id, teid_t s1_u_sgw_fteid, gtpv2c_cause_value_t cause) {
  cb_response->cause.cause_value                  = cause;
  cb_response->sgw_s11_teid                       = teid;
  cb_response->bearer_contexts.num_bearer_context = 1;
  bearer_context_within_create_bearer_response_t* bc_context =
      &(cb_response->bearer_contexts.bearer_contexts[0]);
  bc_context->eps_bearer_id                     = 6;
  bc_context->cause.cause_value                 = cause;
  bc_context->s1u_enb_fteid.ipv4                = 1;
  bc_context->s1u_enb_fteid.interface_type      = S1_U_ENODEB_GTP_U;
  bc_context->s1u_enb_fteid.teid                = 30;
  bc_context->s1u_enb_fteid.ipv4_address.s_addr = 0x0a160328;

  bc_context->s1u_sgw_fteid.teid = s1_u_sgw_fteid;
}

void fill_delete_bearer_response(
    itti_s11_nw_init_deactv_bearer_rsp_t* db_response,
    uint32_t s_gw_teid_s11_s4, uint8_t eps_bearer_id,
    gtpv2c_cause_value_t cause) {
  db_response->delete_default_bearer              = false;
  db_response->s_gw_teid_s11_s4                   = s_gw_teid_s11_s4;
  db_response->bearer_contexts.num_bearer_context = 1;
  for (uint8_t idx = 0; idx < db_response->bearer_contexts.num_bearer_context;
       idx++) {
    db_response->bearer_contexts.bearer_contexts[idx].eps_bearer_id =
        eps_bearer_id;
    db_response->bearer_contexts.bearer_contexts[idx].cause.cause_value = cause;
  }
}

sgw_state_t* SgwS8Config::create_ue_context(mme_sgw_tunnel_t* sgw_s11_tunnel) {
  sgw_state_init(false, config);
  sgw_state_t* sgw_state     = get_sgw_state(false);
  sgw_s11_tunnel->local_teid = sgw_s8_generate_new_cp_teid();
  sgw_update_teid_in_ue_context(sgw_state, imsi64, sgw_s11_tunnel->local_teid);
  return sgw_state;
}
