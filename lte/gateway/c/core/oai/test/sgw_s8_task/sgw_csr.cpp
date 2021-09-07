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

#include <gtest/gtest.h>
#include <string>
#include "sgw_s8_state_manager.h"
#include "sgw_s8_state.h"

extern "C" {
#include "log.h"
#include "sgw_s8_s11_handlers.h"
#include "spgw_types.h"
#include "s11_messages_types.h"
#include "common_types.h"
#include "sgw_config.h"
#include "dynamic_memory_check.h"
#include "sgw_context_manager.h"
}

using ::testing::Test;

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

void fill_itti_csreq(itti_s11_create_session_request_t* session_req_pP) {
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

// Initialize config params
class SgwS8Config : public ::testing::Test {
 protected:
  sgw_config_t* config =
      reinterpret_cast<sgw_config_t*>(calloc(1, sizeof(sgw_config_t)));
  uint64_t imsi64 = 1010000000001;
  virtual void SetUp() {
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
  virtual void TearDown() {
    bdestroy_wrapper(&config->itti_config.log_file);
    bdestroy_wrapper(&config->ipv4.if_name_S1u_S12_S4_up);
    bdestroy_wrapper(&config->ipv4.if_name_S5_S8_up);
    bdestroy_wrapper(&config->ipv4.if_name_S11);
    bdestroy_wrapper(&config->config_file);
    free(config);
  }
};

// TC validates creation of UE context on reception of Create Session Req
TEST_F(SgwS8Config, create_ue_context_cs_req) {
  sgw_state_init(false, config);
  sgw_state_t* sgw_state = get_sgw_state(false);

  spgw_ue_context_t* ue_context_p = NULL;
  mme_sgw_tunnel_t sgw_s11_tunnel = {0};
  sgw_s11_tunnel.local_teid       = sgw_s8_generate_new_cp_teid();
  EXPECT_EQ(
      sgw_update_teid_in_ue_context(
          sgw_state, imsi64, sgw_s11_tunnel.local_teid),
      RETURNok);
  EXPECT_EQ(
      hashtable_ts_get(
          sgw_state->imsi_ue_context_htbl, (const hash_key_t) imsi64,
          reinterpret_cast<void**>(&ue_context_p)),
      HASH_TABLE_OK);

  sgw_state_exit();
}

// TC validates creation of pdn session on reception of Create Session Req
TEST_F(SgwS8Config, create_pdn_session_cs_req) {
  sgw_state_init(false, config);
  sgw_state_t* sgw_state = get_sgw_state(false);

  spgw_ue_context_t* ue_context_p = NULL;
  mme_sgw_tunnel_t sgw_s11_tunnel = {0};
  sgw_s11_tunnel.local_teid       = sgw_s8_generate_new_cp_teid();

  EXPECT_EQ(
      sgw_update_teid_in_ue_context(
          sgw_state, imsi64, sgw_s11_tunnel.local_teid),
      RETURNok);
  EXPECT_EQ(
      hashtable_ts_get(
          sgw_state->imsi_ue_context_htbl, (const hash_key_t) imsi64,
          reinterpret_cast<void**>(&ue_context_p)),
      HASH_TABLE_OK);

  sgw_eps_bearer_context_information_t* sgw_pdn_session = NULL;
  sgw_pdn_session = sgw_create_bearer_context_information_in_collection(
      sgw_s11_tunnel.local_teid);
  EXPECT_TRUE(sgw_pdn_session != nullptr);

  sgw_state_exit();
}

// TC validates creation of bearer context on reception of Create Session Req
TEST_F(SgwS8Config, create_bearer_within_pdn_session_cs_req) {
  sgw_state_init(false, config);
  sgw_state_t* sgw_state = get_sgw_state(false);

  spgw_ue_context_t* ue_context_p = NULL;
  mme_sgw_tunnel_t sgw_s11_tunnel = {0};
  sgw_s11_tunnel.local_teid       = sgw_s8_generate_new_cp_teid();
  EXPECT_EQ(
      sgw_update_teid_in_ue_context(
          sgw_state, imsi64, sgw_s11_tunnel.local_teid),
      RETURNok);
  EXPECT_EQ(
      hashtable_ts_get(
          sgw_state->imsi_ue_context_htbl, (const hash_key_t) imsi64,
          reinterpret_cast<void**>(&ue_context_p)),
      HASH_TABLE_OK);

  sgw_eps_bearer_context_information_t* sgw_pdn_session = NULL;
  sgw_pdn_session = sgw_create_bearer_context_information_in_collection(
      sgw_s11_tunnel.local_teid);
  EXPECT_TRUE(sgw_pdn_session != nullptr);
  itti_s11_create_session_request_t session_req = {0};
  fill_itti_csreq(&session_req);
  sgw_s11_tunnel.remote_teid = session_req.sender_fteid_for_cp.teid;
  EXPECT_EQ(
      sgw_update_bearer_context_information_on_csreq(
          sgw_state, sgw_pdn_session, sgw_s11_tunnel, &session_req, imsi64),
      RETURNok);
  sgw_state_exit();
}

// TC validates updation of bearer context on reception of Create Session Rsp
TEST_F(SgwS8Config, update_pdn_session_on_cs_rsp) {
  sgw_state_init(false, config);
  sgw_state_t* sgw_state = get_sgw_state(false);

  spgw_ue_context_t* ue_context_p = NULL;
  mme_sgw_tunnel_t sgw_s11_tunnel = {0};
  sgw_s11_tunnel.local_teid       = sgw_s8_generate_new_cp_teid();

  EXPECT_EQ(
      sgw_update_teid_in_ue_context(
          sgw_state, imsi64, sgw_s11_tunnel.local_teid),
      RETURNok);
  EXPECT_EQ(
      hashtable_ts_get(
          sgw_state->imsi_ue_context_htbl, (const hash_key_t) imsi64,
          reinterpret_cast<void**>(&ue_context_p)),
      HASH_TABLE_OK);

  sgw_eps_bearer_context_information_t* sgw_pdn_session = NULL;
  sgw_pdn_session = sgw_create_bearer_context_information_in_collection(
      sgw_s11_tunnel.local_teid);
  EXPECT_TRUE(sgw_pdn_session != nullptr);

  itti_s11_create_session_request_t session_req = {0};
  fill_itti_csreq(&session_req);

  sgw_s11_tunnel.remote_teid = session_req.sender_fteid_for_cp.teid;
  EXPECT_EQ(
      sgw_update_bearer_context_information_on_csreq(
          sgw_state, sgw_pdn_session, sgw_s11_tunnel, &session_req, imsi64),
      RETURNok);

  EXPECT_TRUE(
      (sgw_get_sgw_eps_bearer_context(sgw_s11_tunnel.local_teid)) != nullptr);
  s8_create_session_response_t csresp = {0};
  fill_itti_csrsp(&csresp, sgw_s11_tunnel.local_teid);

  EXPECT_EQ(
      sgw_update_bearer_context_information_on_csrsp(sgw_pdn_session, &csresp),
      RETURNok);

  sgw_eps_bearer_ctxt_t* bearer_ctx_p = sgw_cm_get_eps_bearer_entry(
      &sgw_pdn_session->pdn_connection, csresp.eps_bearer_id);
  EXPECT_TRUE(bearer_ctx_p != nullptr);

  EXPECT_TRUE(
      bearer_ctx_p->paa.ipv4_address.s_addr == csresp.paa.ipv4_address.s_addr);
  EXPECT_TRUE(
      bearer_ctx_p->p_gw_teid_S5_S8_up ==
      csresp.bearer_context[0].pgw_s8_up.teid);

  sgw_state_exit();
}

// TC indicates that SGW_S8 has received incorrect sgw_s8_teid in Create Session
// Rsp
TEST_F(SgwS8Config, recv_different_cp_teid_on_cs_rsp) {
  sgw_state_init(false, config);
  sgw_state_t* sgw_state          = get_sgw_state(false);
  spgw_ue_context_t* ue_context_p = NULL;

  mme_sgw_tunnel_t sgw_s11_tunnel = {0};
  sgw_s11_tunnel.local_teid       = sgw_s8_generate_new_cp_teid();
  EXPECT_EQ(
      sgw_update_teid_in_ue_context(
          sgw_state, imsi64, sgw_s11_tunnel.local_teid),
      RETURNok);
  EXPECT_EQ(
      hashtable_ts_get(
          sgw_state->imsi_ue_context_htbl, (const hash_key_t) imsi64,
          reinterpret_cast<void**>(&ue_context_p)),
      HASH_TABLE_OK);

  sgw_eps_bearer_context_information_t* sgw_pdn_session = NULL;
  sgw_pdn_session = sgw_create_bearer_context_information_in_collection(
      sgw_s11_tunnel.local_teid);
  EXPECT_TRUE(sgw_pdn_session != nullptr);

  itti_s11_create_session_request_t session_req = {0};
  fill_itti_csreq(&session_req);

  sgw_s11_tunnel.remote_teid = session_req.sender_fteid_for_cp.teid;
  EXPECT_EQ(
      sgw_update_bearer_context_information_on_csreq(
          sgw_state, sgw_pdn_session, sgw_s11_tunnel, &session_req, imsi64),
      RETURNok);

  s8_create_session_response_t csresp = {0};
  fill_itti_csrsp(&csresp, sgw_s11_tunnel.local_teid);
  csresp.context_teid = 7;  // Send wrong context_teid
  EXPECT_TRUE((sgw_get_sgw_eps_bearer_context(csresp.context_teid)) == nullptr);

  sgw_state_exit();
}

// TC indicates that SGW_S8 has received incorrect eps bearer id in Create
// Session Rsp
TEST_F(SgwS8Config, failed_to_get_bearer_context_on_cs_rsp) {
  sgw_state_init(false, config);
  sgw_state_t* sgw_state          = get_sgw_state(false);
  spgw_ue_context_t* ue_context_p = NULL;

  mme_sgw_tunnel_t sgw_s11_tunnel = {0};
  sgw_s11_tunnel.local_teid       = sgw_s8_generate_new_cp_teid();
  EXPECT_EQ(
      sgw_update_teid_in_ue_context(
          sgw_state, imsi64, sgw_s11_tunnel.local_teid),
      RETURNok);
  EXPECT_EQ(
      hashtable_ts_get(
          sgw_state->imsi_ue_context_htbl, (const hash_key_t) imsi64,
          reinterpret_cast<void**>(&ue_context_p)),
      HASH_TABLE_OK);

  sgw_eps_bearer_context_information_t* sgw_pdn_session = NULL;
  sgw_pdn_session = sgw_create_bearer_context_information_in_collection(
      sgw_s11_tunnel.local_teid);
  EXPECT_TRUE(sgw_pdn_session != nullptr);

  itti_s11_create_session_request_t session_req = {0};
  fill_itti_csreq(&session_req);

  sgw_s11_tunnel.remote_teid = session_req.sender_fteid_for_cp.teid;
  EXPECT_EQ(
      sgw_update_bearer_context_information_on_csreq(
          sgw_state, sgw_pdn_session, sgw_s11_tunnel, &session_req, imsi64),
      RETURNok);

  EXPECT_TRUE(
      (sgw_get_sgw_eps_bearer_context(sgw_s11_tunnel.local_teid)) != nullptr);
  s8_create_session_response_t csresp = {0};
  fill_itti_csrsp(&csresp, sgw_s11_tunnel.local_teid);
  csresp.eps_bearer_id = 7;  // Send wrong eps_bearer_id
  EXPECT_EQ(
      sgw_update_bearer_context_information_on_csrsp(sgw_pdn_session, &csresp),
      RETURNerror);
  sgw_state_exit();
}
