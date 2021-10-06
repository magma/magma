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

using ::testing::Test;

/* TC validates creation of ue context, pdn context and bearer context
 * on reception of Create Session Req
 */
TEST_F(SgwS8Config, create_context_on_cs_req) {
  sgw_state_init(false, config);
  sgw_state_t* sgw_state = get_sgw_state(false);

  spgw_ue_context_t* ue_context_p = NULL;
  mme_sgw_tunnel_t sgw_s11_tunnel = {0};
  sgw_s11_tunnel.local_teid       = sgw_s8_generate_new_cp_teid();

  // validates creation of UE context on reception of Create Session Req
  EXPECT_EQ(
      sgw_update_teid_in_ue_context(
          sgw_state, imsi64, sgw_s11_tunnel.local_teid),
      RETURNok);
  EXPECT_EQ(
      hashtable_ts_get(
          sgw_state->imsi_ue_context_htbl, (const hash_key_t) imsi64,
          reinterpret_cast<void**>(&ue_context_p)),
      HASH_TABLE_OK);

  sgw_s11_teid_t* s11_teid_p        = NULL;
  bool sgw_local_teid_present_in_ue = false;
  LIST_FOREACH(s11_teid_p, &ue_context_p->sgw_s11_teid_list, entries) {
    if ((s11_teid_p) &&
        (s11_teid_p->sgw_s11_teid == sgw_s11_tunnel.local_teid)) {
      sgw_local_teid_present_in_ue = true;
    }
  }
  EXPECT_EQ(sgw_local_teid_present_in_ue, true);

  // validates creation of pdn session on reception of Create Session Req
  sgw_eps_bearer_context_information_t* sgw_pdn_session = NULL;
  sgw_pdn_session = sgw_create_bearer_context_information_in_collection(
      sgw_s11_tunnel.local_teid);
  EXPECT_TRUE(sgw_pdn_session != nullptr);
  sgw_pdn_session                = nullptr;
  hash_table_ts_t* state_imsi_ht = get_sgw_ue_state();
  EXPECT_EQ(
      hashtable_ts_get(
          state_imsi_ht, sgw_s11_tunnel.local_teid,
          reinterpret_cast<void**>(&sgw_pdn_session)),
      HASH_TABLE_OK);

  // validates creation of bearer context on reception of Create Session Req
  itti_s11_create_session_request_t session_req = {0};
  fill_itti_csreq(&session_req, default_eps_bearer_id);
  memcpy(session_req.apn, "internet", sizeof("internet"));
  sgw_s11_tunnel.remote_teid = session_req.sender_fteid_for_cp.teid;
  EXPECT_EQ(
      sgw_update_bearer_context_information_on_csreq(
          sgw_state, sgw_pdn_session, sgw_s11_tunnel, &session_req, imsi64),
      RETURNok);

  // Validates whether MME's control plane teid is set within pdn session
  EXPECT_EQ(
      sgw_pdn_session->mme_teid_S11, session_req.sender_fteid_for_cp.teid);
  EXPECT_EQ(strcmp(sgw_pdn_session->pdn_connection.apn_in_use, "internet"), 0);
  EXPECT_EQ(
      sgw_pdn_session->pdn_connection.default_bearer, session_req.default_ebi);

  // Validates whether bearer is created within pdn session
  bool bearer_id_inserted                  = false;
  sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p = nullptr;
  for (uint8_t idx = 0; idx < BEARERS_PER_UE; idx++) {
    eps_bearer_ctxt_p =
        sgw_pdn_session->pdn_connection.sgw_eps_bearers_array[idx];
    if (eps_bearer_ctxt_p &&
        (eps_bearer_ctxt_p->eps_bearer_id ==
         session_req.bearer_contexts_to_be_created.bearer_contexts[idx]
             .eps_bearer_id)) {
      bearer_id_inserted = true;
      break;
    }
  }
  EXPECT_EQ(bearer_id_inserted, true);
  // Validates whether userplane teids are created for s1-u and s8-u interfaces
  EXPECT_GT(eps_bearer_ctxt_p->s_gw_teid_S1u_S12_S4_up, 0);
  EXPECT_GT(eps_bearer_ctxt_p->s_gw_teid_S5_S8_up, 0);
  sgw_state_exit();
}

// TC validates updation of bearer context on reception of Create Session Rsp
TEST_F(SgwS8Config, update_pdn_session_on_cs_rsp) {
  mme_sgw_tunnel_t sgw_s11_tunnel = {0};
  sgw_state_t* sgw_state          = create_ue_context(&sgw_s11_tunnel);
  sgw_eps_bearer_context_information_t* sgw_pdn_session = NULL;
  sgw_pdn_session = sgw_create_bearer_context_information_in_collection(
      sgw_s11_tunnel.local_teid);

  itti_s11_create_session_request_t session_req = {0};
  fill_itti_csreq(&session_req, default_eps_bearer_id);

  sgw_s11_tunnel.remote_teid = session_req.sender_fteid_for_cp.teid;
  sgw_update_bearer_context_information_on_csreq(
      sgw_state, sgw_pdn_session, sgw_s11_tunnel, &session_req, imsi64);

  EXPECT_EQ(strcmp(sgw_pdn_session->pdn_connection.apn_in_use, "NO APN"), 0);
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
  mme_sgw_tunnel_t sgw_s11_tunnel = {0};
  sgw_state_t* sgw_state          = create_ue_context(&sgw_s11_tunnel);

  sgw_eps_bearer_context_information_t* sgw_pdn_session = NULL;
  sgw_pdn_session = sgw_create_bearer_context_information_in_collection(
      sgw_s11_tunnel.local_teid);

  itti_s11_create_session_request_t session_req = {0};
  fill_itti_csreq(&session_req, default_eps_bearer_id);

  sgw_s11_tunnel.remote_teid = session_req.sender_fteid_for_cp.teid;
  sgw_update_bearer_context_information_on_csreq(
      sgw_state, sgw_pdn_session, sgw_s11_tunnel, &session_req, imsi64);

  s8_create_session_response_t csresp = {0};
  fill_itti_csrsp(&csresp, sgw_s11_tunnel.local_teid);
  csresp.context_teid = 7;  // Send wrong context_teid
  EXPECT_TRUE((sgw_get_sgw_eps_bearer_context(csresp.context_teid)) == nullptr);

  sgw_state_exit();
}

// TC indicates that SGW_S8 has received incorrect eps bearer id in Create
// Session Rsp
TEST_F(SgwS8Config, failed_to_get_bearer_context_on_cs_rsp) {
  mme_sgw_tunnel_t sgw_s11_tunnel = {0};
  sgw_state_t* sgw_state          = create_ue_context(&sgw_s11_tunnel);

  sgw_eps_bearer_context_information_t* sgw_pdn_session = NULL;
  sgw_pdn_session = sgw_create_bearer_context_information_in_collection(
      sgw_s11_tunnel.local_teid);

  itti_s11_create_session_request_t session_req = {0};
  fill_itti_csreq(&session_req, default_eps_bearer_id);

  sgw_s11_tunnel.remote_teid = session_req.sender_fteid_for_cp.teid;
  sgw_update_bearer_context_information_on_csreq(
      sgw_state, sgw_pdn_session, sgw_s11_tunnel, &session_req, imsi64);

  s8_create_session_response_t csresp = {0};
  fill_itti_csrsp(&csresp, sgw_s11_tunnel.local_teid);
  csresp.eps_bearer_id = 7;  // Send wrong eps_bearer_id
  EXPECT_EQ(
      sgw_update_bearer_context_information_on_csrsp(sgw_pdn_session, &csresp),
      RETURNerror);
  sgw_state_exit();
}
