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

#include <memory>
#include "lte/gateway/c/core/oai/test/sgw_s8_task/sgw_s8_utility.h"
#include "lte/gateway/c/core/oai/lib/openflow/controller/ControllerMain.h"

task_zmq_ctx_t task_zmq_ctx_main_s8;
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

void fill_itti_csrsp(
    s8_create_session_response_t* csr_resp,
    uint32_t temporary_create_session_procedure_id) {
  uint8_t idx = 0;
  fill_imsi((reinterpret_cast<char*>(csr_resp->imsi)));
  csr_resp->imsi_length = 15;

  csr_resp->pdn_type                = IPv4;
  csr_resp->paa.pdn_type            = IPv4;
  csr_resp->paa.ipv4_address.s_addr = 0xc0a87e1;
  csr_resp->context_teid  = 16;  // This teid would be allocated by orc8r
  csr_resp->eps_bearer_id = 5;
  csr_resp->temporary_create_session_procedure_id =
      temporary_create_session_procedure_id;

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
  bc_context->eps_bearer_id                     = eps_bearer_id;
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

void fill_delete_bearer_request(
    s8_delete_bearer_request_t* db_req, uint32_t teid, uint8_t eps_bearer_id) {
  db_req->context_teid      = teid;
  db_req->num_eps_bearer_id = 1;
  for (uint8_t idx = 0; idx < db_req->num_eps_bearer_id; idx++) {
    db_req->eps_bearer_id[idx] = eps_bearer_id;
  }
  db_req->sequence_number = 2;
}

void fill_delete_session_request(
    itti_s11_delete_session_request_t* ds_req_p, uint32_t teid, uint8_t lbi) {
  ds_req_p->lbi  = lbi;
  ds_req_p->teid = teid;

  ds_req_p->sender_fteid_for_cp.teid                = 1;
  ds_req_p->sender_fteid_for_cp.ipv4_address.s_addr = 0x8e3ca8c0;
  ds_req_p->sender_fteid_for_cp.interface_type      = S11_MME_GTP_C;
  uint8_t idx                                       = 0;
  ds_req_p->serving_network.mcc[idx++]              = 0;
  ds_req_p->serving_network.mcc[idx++]              = 0;
  ds_req_p->serving_network.mcc[idx]                = 0;
  idx                                               = 0;
  ds_req_p->serving_network.mnc[idx++]              = 1;
  ds_req_p->serving_network.mnc[idx++]              = 1;
  ds_req_p->serving_network.mnc[idx]                = 15;
}

void fill_delete_session_response(
    s8_delete_session_response_t* ds_rsp_p, uint32_t teid, uint8_t cause) {
  ds_rsp_p->context_teid = teid;
  ds_rsp_p->cause        = cause;
}

sgw_state_t* SgwS8ConfigAndCreateMock::create_and_get_contexts_on_cs_req(
    uint32_t* temporary_create_session_procedure_id,
    sgw_eps_bearer_context_information_t** sgw_pdn_session) {
  sgw_state_t* sgw_state = get_sgw_state(false);
  *sgw_pdn_session       = sgw_create_bearer_context_information_in_collection(
      sgw_state, temporary_create_session_procedure_id);

  itti_s11_create_session_request_t session_req = {0};
  fill_itti_csreq(&session_req, default_eps_bearer_id);
  memcpy(session_req.apn, "internet", sizeof("internet"));

  *sgw_pdn_session = sgw_create_bearer_context_information_in_collection(
      sgw_state, temporary_create_session_procedure_id);
  sgw_update_bearer_context_information_on_csreq(
      sgw_state, *sgw_pdn_session, &session_req, imsi64);

  return sgw_state;
}

sgw_state_t* SgwS8ConfigAndCreateMock::create_ue_context(
    mme_sgw_tunnel_t* sgw_s11_tunnel) {
  sgw_state_t* sgw_state = get_sgw_state(false);
  sgw_update_teid_in_ue_context(sgw_state, imsi64, sgw_s11_tunnel->local_teid);
  return sgw_state;
}

void SgwS8ConfigAndCreateMock::sgw_s8_config_init() {
  config->itti_config.queue_size     = 0;
  std::string file_string            = "/var/opt/magma/tmp/spgw.conf";
  config->itti_config.log_file       = bfromcstr(file_string.c_str());
  std::string s1u_if_name            = "eth1";
  config->ipv4.if_name_S1u_S12_S4_up = bfromcstr(s1u_if_name.c_str());
  config->ipv4.S1u_S12_S4_up.s_addr  = 0x8e3ca8c0;
  config->ipv4.netmask_S1u_S12_S4_up = 24;
  std::string s5s8u_if_name          = "eth0";
  config->ipv4.if_name_S5_S8_up      = bfromcstr(s5s8u_if_name.c_str());
  config->ipv4.S5_S8_up.s_addr       = 0xf02000a;
  config->ipv4.netmask_S5_S8_up      = 24;
  std::string s11                    = "lo";
  config->ipv4.if_name_S11           = bfromcstr(s11.c_str());
  config->ipv4.S11.s_addr            = 0x100007f;
  config->ipv4.netmask_S11           = 8;
  config->udp_port_S1u_S12_S4_up     = 2152;
  config->config_file                = bfromcstr(file_string.c_str());
}

static int handle_message_test_sgw_s8(
    zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    default: { } break; }

  itti_free_msg_content(received_message_p);
  free(received_message_p);
  return 0;
}

void SgwS8ConfigAndCreateMock::SetUp() {
  mme_app_handler = std::make_shared<MockMmeAppHandler>();

  itti_init(
      TASK_MAX, THREAD_MAX, MESSAGES_ID_MAX, tasks_info, messages_info, NULL,
      NULL);
  sgw_s8_config_init();
  task_id_t task_id_list[3] = {TASK_SGW_S8, TASK_MME_APP, TASK_GRPC_SERVICE};
  init_task_context(
      TASK_MAIN, task_id_list, 3, handle_message_test_sgw_s8,
      &task_zmq_ctx_main_s8);

  std::thread task_grpc(start_mock_grpc_task);
  std::thread task_mme_app(start_mock_mme_app_task, mme_app_handler);
  task_grpc.detach();
  task_mme_app.detach();

  sgw_s8_init(config);
  std::this_thread::sleep_for(
      std::chrono::milliseconds(SLEEP_AT_INITIALIZATION_TIME_MS));
}

void SgwS8ConfigAndCreateMock::TearDown() {
  sgw_state_exit();
  bdestroy_wrapper(&config->itti_config.log_file);
  bdestroy_wrapper(&config->ipv4.if_name_S1u_S12_S4_up);
  bdestroy_wrapper(&config->ipv4.if_name_S5_S8_up);
  bdestroy_wrapper(&config->ipv4.if_name_S11);
  bdestroy_wrapper(&config->config_file);
  free(config);

  send_terminate_message_fatal(&task_zmq_ctx_main_s8);
  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(500));
  destroy_task_context(&task_zmq_ctx_main_s8);
  itti_free_desc_threads();
}
