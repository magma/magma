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

#include "lte/gateway/c/core/oai/test/sgw_s8_task/sgw_s8_utility.h"

using ::testing::Test;

MATCHER_P2(check_cause_in_ds_rsp, cause, teid, "") {
  auto ds_rsp_rcvd_at_mme =
      static_cast<itti_s11_delete_session_response_t>(arg);
  if (ds_rsp_rcvd_at_mme.cause.cause_value == cause) {
    return true;
  }
  if (ds_rsp_rcvd_at_mme.teid == teid) {
    return true;
  }
  return false;
}

// TC validates successful handling of create session request
TEST_F(SgwS8ConfigAndCreateMock, create_context_on_cs_req_success) {
  sgw_state_t* sgw_state                        = get_sgw_state(false);
  itti_s11_create_session_request_t session_req = {0};
  fill_itti_csreq(&session_req, default_eps_bearer_id);
  memcpy(session_req.apn, "internet", sizeof("internet"));
  EXPECT_EQ(
      sgw_s8_handle_s11_create_session_request(sgw_state, &session_req, imsi64),
      RETURNok);
}

/* TC validates creation of ue context, pdn context and bearer context
 * on reception of Create Session Req
 */
TEST_F(SgwS8ConfigAndCreateMock, create_context_on_cs_req) {
  sgw_state_t* sgw_state                         = get_sgw_state(false);
  uint32_t temporary_create_session_procedure_id = 0;

  // validates creation of pdn session on reception of Create Session Req with
  // key as temporary session id
  sgw_eps_bearer_context_information_t* sgw_pdn_session =
      sgw_create_bearer_context_information_in_collection(
          sgw_state, &temporary_create_session_procedure_id);
  EXPECT_TRUE(sgw_pdn_session != nullptr);
  sgw_pdn_session = nullptr;
  EXPECT_EQ(
      hashtable_ts_get(
          sgw_state->temporary_create_session_procedure_id_htbl,
          temporary_create_session_procedure_id,
          reinterpret_cast<void**>(&sgw_pdn_session)),
      HASH_TABLE_OK);

  // validates creation of bearer context on reception of Create Session Req
  itti_s11_create_session_request_t session_req = {0};
  fill_itti_csreq(&session_req, default_eps_bearer_id);
  memcpy(session_req.apn, "internet", sizeof("internet"));
  EXPECT_EQ(
      sgw_update_bearer_context_information_on_csreq(
          sgw_state, sgw_pdn_session, &session_req, imsi64),
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
}

// TC validates updation of bearer context on reception of Create Session Rsp
TEST_F(SgwS8ConfigAndCreateMock, update_pdn_session_on_cs_rsp) {
  sgw_state_t* sgw_state = get_sgw_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  sgw_eps_bearer_context_information_t* sgw_pdn_session = NULL;
  uint32_t temporary_create_session_procedure_id        = 0;
  sgw_pdn_session = sgw_create_bearer_context_information_in_collection(
      sgw_state, &temporary_create_session_procedure_id);

  itti_s11_create_session_request_t session_req = {0};
  fill_itti_csreq(&session_req, default_eps_bearer_id);

  sgw_update_bearer_context_information_on_csreq(
      sgw_state, sgw_pdn_session, &session_req, imsi64);

  EXPECT_EQ(strcmp(sgw_pdn_session->pdn_connection.apn_in_use, "NO APN"), 0);

  s8_create_session_response_t csresp = {0};
  fill_itti_csrsp(&csresp, temporary_create_session_procedure_id);

  EXPECT_CALL(*mme_app_handler, mme_app_handle_create_sess_resp())
      .Times(1)
      .WillOnce(ReturnFromAsyncTask(&cv));
  EXPECT_EQ(
      sgw_s8_handle_create_session_response(sgw_state, &csresp, imsi64),
      RETURNok);
  cv.wait_for(lock, std::chrono::milliseconds(END_OF_TESTCASE_SLEEP_MS));

  EXPECT_TRUE((sgw_get_sgw_eps_bearer_context(csresp.context_teid)) != nullptr);
  sgw_eps_bearer_ctxt_t* bearer_ctx_p = sgw_cm_get_eps_bearer_entry(
      &sgw_pdn_session->pdn_connection, csresp.eps_bearer_id);
  EXPECT_TRUE(bearer_ctx_p != nullptr);

  EXPECT_TRUE(
      bearer_ctx_p->paa.ipv4_address.s_addr == csresp.paa.ipv4_address.s_addr);
  EXPECT_TRUE(
      bearer_ctx_p->p_gw_teid_S5_S8_up ==
      csresp.bearer_context[0].pgw_s8_up.teid);
  // Check pdn session is removed from
  // temporary_create_session_procedure_id_htbl
  EXPECT_EQ(
      hashtable_ts_get(
          sgw_state->temporary_create_session_procedure_id_htbl,
          temporary_create_session_procedure_id,
          reinterpret_cast<void**>(&sgw_pdn_session)),
      HASH_TABLE_KEY_NOT_EXISTS);
}

// TC indicates that SGW_S8 has received incorrect temporary session id in
// Create Session Rsp
TEST_F(
    SgwS8ConfigAndCreateMock,
    recv_different_temporary_create_session_procedure_id_on_cs_rsp) {
  sgw_state_t* sgw_state = get_sgw_state(false);

  sgw_eps_bearer_context_information_t* sgw_pdn_session = NULL;
  uint32_t temporary_create_session_procedure_id        = 0;
  sgw_pdn_session = sgw_create_bearer_context_information_in_collection(
      sgw_state, &temporary_create_session_procedure_id);

  itti_s11_create_session_request_t session_req = {0};
  fill_itti_csreq(&session_req, default_eps_bearer_id);

  sgw_update_bearer_context_information_on_csreq(
      sgw_state, sgw_pdn_session, &session_req, imsi64);

  s8_create_session_response_t csresp = {0};
  fill_itti_csrsp(&csresp, temporary_create_session_procedure_id + 1);
  EXPECT_EQ(
      sgw_s8_handle_create_session_response(sgw_state, &csresp, imsi64),
      RETURNerror);
}

// TC indicates that SGW_S8 has received incorrect sgw_s8_teid in Create Session
// Rsp
TEST_F(SgwS8ConfigAndCreateMock, recv_different_sgw_s8_teid) {
  sgw_state_t* sgw_state = get_sgw_state(false);

  sgw_eps_bearer_context_information_t* sgw_pdn_session = NULL;
  uint32_t temporary_create_session_procedure_id        = 0;
  sgw_pdn_session = sgw_create_bearer_context_information_in_collection(
      sgw_state, &temporary_create_session_procedure_id);

  itti_s11_create_session_request_t session_req = {0};
  fill_itti_csreq(&session_req, default_eps_bearer_id);

  sgw_update_bearer_context_information_on_csreq(
      sgw_state, sgw_pdn_session, &session_req, imsi64);

  s8_create_session_response_t csresp = {0};
  fill_itti_csrsp(&csresp, temporary_create_session_procedure_id);
  sgw_s8_handle_create_session_response(sgw_state, &csresp, imsi64);
  // validate with wrong sgw_s8_teid, fails to get sgw_pdn_session
  EXPECT_EQ(
      (sgw_get_sgw_eps_bearer_context((csresp.context_teid + 1))), nullptr);
}

// TC indicates that SGW_S8 has received incorrect eps bearer id in Create
// Session Rsp
TEST_F(SgwS8ConfigAndCreateMock, failed_to_get_bearer_context_on_cs_rsp) {
  uint32_t temporary_create_session_procedure_id = 0;

  sgw_state_t* sgw_state = get_sgw_state(false);
  sgw_eps_bearer_context_information_t* sgw_pdn_session =
      sgw_create_bearer_context_information_in_collection(
          sgw_state, &temporary_create_session_procedure_id);

  itti_s11_create_session_request_t session_req = {0};
  fill_itti_csreq(&session_req, default_eps_bearer_id);

  sgw_update_bearer_context_information_on_csreq(
      sgw_state, sgw_pdn_session, &session_req, imsi64);

  s8_create_session_response_t csresp = {0};
  fill_itti_csrsp(&csresp, temporary_create_session_procedure_id);
  csresp.eps_bearer_id = 7;  // Send wrong eps_bearer_id
  // fails to update bearer context
  EXPECT_EQ(
      sgw_update_bearer_context_information_on_csrsp(sgw_pdn_session, &csresp),
      RETURNerror);
}

// TC validates deletion of pdn session
TEST_F(SgwS8ConfigAndCreateMock, delete_session_req_handling) {
  sgw_state_t* sgw_state = get_sgw_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  sgw_eps_bearer_context_information_t* sgw_pdn_session = NULL;
  uint32_t temporary_create_session_procedure_id        = 0;
  sgw_pdn_session = sgw_create_bearer_context_information_in_collection(
      sgw_state, &temporary_create_session_procedure_id);

  itti_s11_create_session_request_t session_req = {0};
  fill_itti_csreq(&session_req, default_eps_bearer_id);

  sgw_update_bearer_context_information_on_csreq(
      sgw_state, sgw_pdn_session, &session_req, imsi64);

  // Since apn is not sent in CSReq, verify that "NO APN" is set for apn_in_use
  // within pdn context
  EXPECT_EQ(strcmp(sgw_pdn_session->pdn_connection.apn_in_use, "NO APN"), 0);

  s8_create_session_response_t csresp = {0};
  fill_itti_csrsp(&csresp, temporary_create_session_procedure_id);

  // Below steps validate that successful handling of create session response
  // which eventually sends message to MME
  EXPECT_CALL(*mme_app_handler, mme_app_handle_create_sess_resp())
      .Times(1)
      .WillOnce(ReturnFromAsyncTask(&cv));
  EXPECT_EQ(
      sgw_s8_handle_create_session_response(sgw_state, &csresp, imsi64),
      RETURNok);
  cv.wait_for(lock, std::chrono::milliseconds(END_OF_TESTCASE_SLEEP_MS));

  EXPECT_TRUE((sgw_get_sgw_eps_bearer_context(csresp.context_teid)) != nullptr);
  sgw_eps_bearer_ctxt_t* bearer_ctx_p = sgw_cm_get_eps_bearer_entry(
      &sgw_pdn_session->pdn_connection, csresp.eps_bearer_id);
  EXPECT_TRUE(bearer_ctx_p != nullptr);

  EXPECT_TRUE(
      bearer_ctx_p->paa.ipv4_address.s_addr == csresp.paa.ipv4_address.s_addr);
  EXPECT_TRUE(
      bearer_ctx_p->p_gw_teid_S5_S8_up ==
      csresp.bearer_context[0].pgw_s8_up.teid);
  itti_s11_delete_session_request_t ds_req = {0};
  fill_delete_session_request(
      &ds_req, csresp.context_teid, csresp.eps_bearer_id);

  // Validates the successful handling of delete session req
  EXPECT_EQ(
      sgw_s8_handle_s11_delete_session_request(sgw_state, &ds_req, imsi64),
      RETURNok);

  // Below steps validate that successful handling of delete session response
  // which eventually sends message to MME
  s8_delete_session_response_t ds_rsp = {0};
  fill_delete_session_response(&ds_rsp, csresp.context_teid, REQUEST_ACCEPTED);
  EXPECT_CALL(
      *mme_app_handler,
      mme_app_handle_delete_sess_rsp(check_cause_in_ds_rsp(
          REQUEST_ACCEPTED, session_req.sender_fteid_for_cp.teid)))
      .Times(1)
      .WillOnce(ReturnFromAsyncTask(&cv));
  EXPECT_EQ(
      sgw_s8_handle_delete_session_response(sgw_state, &ds_rsp, imsi64),
      RETURNok);
  cv.wait_for(lock, std::chrono::milliseconds(END_OF_TESTCASE_SLEEP_MS));

  // after handling of delete session response, makes sure that pdn context is
  // deleted
  EXPECT_TRUE((sgw_get_sgw_eps_bearer_context(csresp.context_teid)) == nullptr);
}

// TC validates that sgw_s8 module fails to find the pdn context based on
// invalid teid
TEST_F(SgwS8ConfigAndCreateMock, delete_session_req_handling_invalid_teid) {
  sgw_state_t* sgw_state = get_sgw_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  sgw_eps_bearer_context_information_t* sgw_pdn_session = NULL;
  uint32_t temporary_create_session_procedure_id        = 0;
  sgw_pdn_session = sgw_create_bearer_context_information_in_collection(
      sgw_state, &temporary_create_session_procedure_id);

  itti_s11_create_session_request_t session_req = {0};
  fill_itti_csreq(&session_req, default_eps_bearer_id);

  sgw_update_bearer_context_information_on_csreq(
      sgw_state, sgw_pdn_session, &session_req, imsi64);

  s8_create_session_response_t csresp = {0};
  fill_itti_csrsp(&csresp, temporary_create_session_procedure_id);
  sgw_s8_handle_create_session_response(sgw_state, &csresp, imsi64);
  // Validate that sgw_s8 fails to find the context for in-correct teid received
  // in delete session req
  itti_s11_delete_session_request_t ds_req = {0};
  fill_delete_session_request(
      &ds_req, csresp.context_teid + 1, csresp.eps_bearer_id);

  // Since sgw_s8 module fails to find pdn context, it returns with response
  // cause set to "context_not_found"
  EXPECT_CALL(
      *mme_app_handler,
      mme_app_handle_delete_sess_rsp(check_cause_in_ds_rsp(
          CONTEXT_NOT_FOUND, session_req.sender_fteid_for_cp.teid)))
      .Times(1)
      .WillOnce(ReturnFromAsyncTask(&cv));
  EXPECT_EQ(
      sgw_s8_handle_s11_delete_session_request(sgw_state, &ds_req, imsi64),
      RETURNerror);

  cv.wait_for(lock, std::chrono::milliseconds(END_OF_TESTCASE_SLEEP_MS));
}
