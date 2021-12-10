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
extern "C" {
#include "lte/gateway/c/core/oai/tasks/sgw/pgw_procedures.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/include/mme_config.h"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_embedded_spgw.h"
}

extern task_zmq_ctx_t task_zmq_ctx_main_s8;
using ::testing::Test;
spgw_config_t spgw_config;
// TC validates updation of bearer context on reception of Create Session Rsp
TEST_F(SgwS8ConfigAndCreateMock, check_dedicated_bearer_creation_request) {
  sgw_eps_bearer_context_information_t* sgw_pdn_session = NULL;
  uint32_t temporary_create_session_procedure_id        = 0;

  sgw_state_t* sgw_state = create_and_get_contexts_on_cs_req(
      &temporary_create_session_procedure_id, &sgw_pdn_session);
  s8_create_session_response_t csresp = {0};
  fill_itti_csrsp(&csresp, temporary_create_session_procedure_id);
  sgw_s8_handle_create_session_response(sgw_state, &csresp, imsi64);

  s8_create_bearer_request_t cb_req = {0};
  fill_create_bearer_request(
      &cb_req, csresp.context_teid, default_eps_bearer_id);

  itti_gx_nw_init_actv_bearer_request_t itti_bearer_req = {0};
  s8_bearer_context_t bc_cbreq = cb_req.bearer_context[0];

  itti_bearer_req.lbi = cb_req.linked_eps_bearer_id;

  memcpy(
      &itti_bearer_req.ul_tft, &bc_cbreq.tft, sizeof(traffic_flow_template_t));
  memcpy(
      &itti_bearer_req.dl_tft, &bc_cbreq.tft, sizeof(traffic_flow_template_t));
  memcpy(&itti_bearer_req.eps_bearer_qos, &bc_cbreq.qos, sizeof(bearer_qos_t));
  teid_t s1_u_sgw_fteid = sgw_get_new_s1u_teid(sgw_state);
  // Validates temporary bearer context is created
  EXPECT_EQ(
      create_temporary_dedicated_bearer_context(
          sgw_pdn_session, &itti_bearer_req, IPv4,
          sgw_state->sgw_ip_address_S1u_S12_S4_up.s_addr,
          &sgw_state->sgw_ipv6_address_S1u_S12_S4_up, s1_u_sgw_fteid,
          cb_req.sequence_number, LOG_SGW_S8),
      RETURNok);
  EXPECT_EQ(
      update_pgw_info_to_temp_dedicated_bearer_context(
          sgw_pdn_session, s1_u_sgw_fteid, &bc_cbreq, sgw_state,
          cb_req.pgw_cp_address),
      RETURNok);
  // Validates sequence number matches with received create bearer request
  pgw_ni_cbr_proc_t* pgw_ni_cbr_proc =
      pgw_get_procedure_create_bearer(sgw_pdn_session);
  EXPECT_TRUE(pgw_ni_cbr_proc != nullptr);

  bool is_seq_number_updated                             = false;
  sgw_eps_bearer_entry_wrapper_t* sgw_eps_bearer_entry_p = nullptr;
  LIST_FOREACH(
      sgw_eps_bearer_entry_p, pgw_ni_cbr_proc->pending_eps_bearers, entries) {
    if ((sgw_eps_bearer_entry_p) &&
        (sgw_eps_bearer_entry_p->sgw_eps_bearer_entry->sgw_sequence_number ==
         cb_req.sequence_number)) {
      is_seq_number_updated = true;
      break;
    }
  }
  EXPECT_EQ(is_seq_number_updated, true);
  EXPECT_EQ(
      sgw_eps_bearer_entry_p->sgw_eps_bearer_entry->p_gw_teid_S5_S8_up,
      cb_req.bearer_context[0].pgw_s8_up.teid);
  free_wrapper(reinterpret_cast<void**>(&cb_req.pgw_cp_address));
}

TEST_F(SgwS8ConfigAndCreateMock, dedicated_bearer_invalid_lbi) {
  sgw_eps_bearer_context_information_t* sgw_pdn_session = NULL;
  uint32_t temporary_create_session_procedure_id        = 0;

  sgw_state_t* sgw_state = create_and_get_contexts_on_cs_req(
      &temporary_create_session_procedure_id, &sgw_pdn_session);
  s8_create_session_response_t csresp = {0};
  fill_itti_csrsp(&csresp, temporary_create_session_procedure_id);
  sgw_s8_handle_create_session_response(sgw_state, &csresp, imsi64);

  s8_create_bearer_request_t cb_req = {0};
  // send invalid default eps bearer id
  fill_create_bearer_request(
      &cb_req, csresp.context_teid, default_eps_bearer_id + 1);
  gtpv2c_cause_value_t cause_value = REQUEST_ACCEPTED;
  EXPECT_EQ(
      sgw_s8_handle_create_bearer_request(sgw_state, &cb_req, &cause_value),
      INVALID_IMSI64);
  free_wrapper(reinterpret_cast<void**>(&cb_req.pgw_cp_address));
}

// TC validates temporary contexts are deleted reception of failed create
// bearer response
TEST_F(SgwS8ConfigAndCreateMock, check_failed_to_create_dedicated_bearer) {
  sgw_eps_bearer_context_information_t* sgw_pdn_session = NULL;
  uint32_t temporary_create_session_procedure_id        = 0;

  sgw_state_t* sgw_state = create_and_get_contexts_on_cs_req(
      &temporary_create_session_procedure_id, &sgw_pdn_session);
  s8_create_session_response_t csresp = {0};
  fill_itti_csrsp(&csresp, temporary_create_session_procedure_id);
  sgw_s8_handle_create_session_response(sgw_state, &csresp, imsi64);

  s8_create_bearer_request_t cb_req = {0};
  fill_create_bearer_request(
      &cb_req, csresp.context_teid, default_eps_bearer_id);

  itti_gx_nw_init_actv_bearer_request_t itti_bearer_req = {0};
  s8_bearer_context_t bc_cbreq = cb_req.bearer_context[0];

  itti_bearer_req.lbi = cb_req.linked_eps_bearer_id;

  memcpy(
      &itti_bearer_req.ul_tft, &bc_cbreq.tft, sizeof(traffic_flow_template_t));
  memcpy(
      &itti_bearer_req.dl_tft, &bc_cbreq.tft, sizeof(traffic_flow_template_t));
  memcpy(&itti_bearer_req.eps_bearer_qos, &bc_cbreq.qos, sizeof(bearer_qos_t));
  teid_t s1_u_sgw_fteid = sgw_get_new_s1u_teid(sgw_state);
  create_temporary_dedicated_bearer_context(
      sgw_pdn_session, &itti_bearer_req, IPv4,
      sgw_state->sgw_ip_address_S1u_S12_S4_up.s_addr,
      &sgw_state->sgw_ipv6_address_S1u_S12_S4_up, s1_u_sgw_fteid,
      cb_req.sequence_number, LOG_SGW_S8);
  update_pgw_info_to_temp_dedicated_bearer_context(
      sgw_pdn_session, s1_u_sgw_fteid, &bc_cbreq, sgw_state,
      cb_req.pgw_cp_address);

  itti_s11_nw_init_actv_bearer_rsp_t s11_actv_bearer_rsp;
  memset(&s11_actv_bearer_rsp, 0, sizeof(itti_s11_nw_init_actv_bearer_rsp_t));
  fill_create_bearer_response(
      &s11_actv_bearer_rsp, csresp.context_teid, default_eps_bearer_id,
      s1_u_sgw_fteid, REQUEST_REJECTED);
  handle_failed_create_bearer_response(
      sgw_pdn_session, s11_actv_bearer_rsp.cause.cause_value, imsi64,
      &s11_actv_bearer_rsp.bearer_contexts.bearer_contexts[0], NULL,
      LOG_SGW_S8);
  EXPECT_EQ(sgw_pdn_session->pending_procedures, nullptr);

  free_wrapper(reinterpret_cast<void**>(&cb_req.pgw_cp_address));
}

// TC validates, failed to find PDN context on wrong sgw_s11_teid received in
// delete bearer response
TEST_F(SgwS8ConfigAndCreateMock, delete_bearer_response_invalid_teid) {
  sgw_eps_bearer_context_information_t* sgw_pdn_session = NULL;
  uint32_t temporary_create_session_procedure_id        = 0;

  sgw_state_t* sgw_state = create_and_get_contexts_on_cs_req(
      &temporary_create_session_procedure_id, &sgw_pdn_session);
  s8_create_session_response_t csresp = {0};
  fill_itti_csrsp(&csresp, temporary_create_session_procedure_id);
  sgw_s8_handle_create_session_response(sgw_state, &csresp, imsi64);

  s8_create_bearer_request_t cb_req = {0};
  fill_create_bearer_request(
      &cb_req, csresp.context_teid, default_eps_bearer_id);

  itti_gx_nw_init_actv_bearer_request_t itti_bearer_req = {0};
  s8_bearer_context_t bc_cbreq = cb_req.bearer_context[0];

  itti_bearer_req.lbi = cb_req.linked_eps_bearer_id;

  memcpy(
      &itti_bearer_req.ul_tft, &bc_cbreq.tft, sizeof(traffic_flow_template_t));
  memcpy(
      &itti_bearer_req.dl_tft, &bc_cbreq.tft, sizeof(traffic_flow_template_t));
  memcpy(&itti_bearer_req.eps_bearer_qos, &bc_cbreq.qos, sizeof(bearer_qos_t));
  teid_t s1_u_sgw_fteid = sgw_get_new_s1u_teid(sgw_state);
  // Validates temporary bearer context is created
  create_temporary_dedicated_bearer_context(
      sgw_pdn_session, &itti_bearer_req, IPv4,
      sgw_state->sgw_ip_address_S1u_S12_S4_up.s_addr,
      &sgw_state->sgw_ipv6_address_S1u_S12_S4_up, s1_u_sgw_fteid,
      cb_req.sequence_number, LOG_SGW_S8);
  update_pgw_info_to_temp_dedicated_bearer_context(
      sgw_pdn_session, s1_u_sgw_fteid, &bc_cbreq, sgw_state,
      cb_req.pgw_cp_address);

  itti_s11_nw_init_deactv_bearer_rsp_t s11_delete_bearer_response = {0};
  fill_delete_bearer_response(
      &s11_delete_bearer_response, csresp.context_teid + 1, 6,
      REQUEST_ACCEPTED);
  EXPECT_EQ(
      sgw_s8_handle_s11_delete_bearer_response(
          sgw_state, &s11_delete_bearer_response, imsi64),
      RETURNerror);
  free_wrapper(reinterpret_cast<void**>(&cb_req.pgw_cp_address));
}

TEST_F(SgwS8ConfigAndCreateMock, create_bearer_req_fails_to_find_ctxt) {
  mme_sgw_tunnel_t sgw_s11_tunnel = {0};
  sgw_state_t* sgw_state          = create_ue_context(&sgw_s11_tunnel);
  sgw_eps_bearer_context_information_t* sgw_pdn_session = NULL;
  uint32_t temporary_create_session_procedure_id        = 0;
  sgw_pdn_session = sgw_create_bearer_context_information_in_collection(
      sgw_state, &temporary_create_session_procedure_id);
  itti_s11_create_session_request_t session_req = {0};
  fill_itti_csreq(&session_req, default_eps_bearer_id);
  memcpy(session_req.apn, "internet", sizeof("internet"));
  sgw_s11_tunnel.remote_teid = session_req.sender_fteid_for_cp.teid;
  sgw_update_bearer_context_information_on_csreq(
      sgw_state, sgw_pdn_session, &session_req, imsi64);
  s8_create_bearer_request_t cb_req = {0};
  // Send wrong sgw_s11_teid
  fill_create_bearer_request(
      &cb_req, sgw_s11_tunnel.local_teid + 1, default_eps_bearer_id);
  gtpv2c_cause_value_t cause_value = REQUEST_REJECTED;
  Imsi_t imsi                      = {0};
  imsi64_t imsi64 =
      sgw_s8_handle_create_bearer_request(sgw_state, &cb_req, &cause_value);
  EXPECT_EQ(imsi64, INVALID_IMSI64);
  sgw_s8_send_failed_create_bearer_response(
      sgw_state, cb_req.sequence_number, cb_req.pgw_cp_address, cause_value,
      imsi, cb_req.bearer_context[0].pgw_s8_up.teid);
  itti_gx_nw_init_actv_bearer_request_t itti_bearer_req = {0};
  s8_bearer_context_t bc_cbreq = cb_req.bearer_context[0];

  itti_bearer_req.lbi = cb_req.linked_eps_bearer_id;

  memcpy(
      &itti_bearer_req.ul_tft, &bc_cbreq.tft, sizeof(traffic_flow_template_t));
  memcpy(
      &itti_bearer_req.dl_tft, &bc_cbreq.tft, sizeof(traffic_flow_template_t));
  memcpy(&itti_bearer_req.eps_bearer_qos, &bc_cbreq.qos, sizeof(bearer_qos_t));
  teid_t s1_u_sgw_fteid = sgw_get_new_s1u_teid(sgw_state);
  // Validates temporary bearer context is created
  EXPECT_EQ(
      create_temporary_dedicated_bearer_context(
          sgw_pdn_session, &itti_bearer_req, IPv4,
          sgw_state->sgw_ip_address_S1u_S12_S4_up.s_addr,
          &sgw_state->sgw_ipv6_address_S1u_S12_S4_up, s1_u_sgw_fteid,
          cb_req.sequence_number, LOG_SGW_S8),
      RETURNok);
  EXPECT_EQ(
      update_pgw_info_to_temp_dedicated_bearer_context(
          sgw_pdn_session, s1_u_sgw_fteid, &bc_cbreq, sgw_state,
          cb_req.pgw_cp_address),
      RETURNok);
  free_wrapper(reinterpret_cast<void**>(&cb_req.pgw_cp_address));
}

MATCHER_P2(check_params_in_db_req, num_eps_bearer_id, eps_bearer_id, "") {
  auto db_req_rcvd_at_mme =
      static_cast<itti_s11_nw_init_deactv_bearer_request_t>(arg);
  if (db_req_rcvd_at_mme.no_of_bearers != num_eps_bearer_id) {
    return false;
  }
  if (memcmp(
          db_req_rcvd_at_mme.ebi, eps_bearer_id,
          sizeof(db_req_rcvd_at_mme.ebi))) {
    return false;
  }
  return true;
}

MATCHER_P2(check_params_in_cb_req, linked_eps_bearer_id, tft, "") {
  auto cb_req_rcvd_at_mme =
      static_cast<itti_s11_nw_init_actv_bearer_request_t>(arg);
  if (cb_req_rcvd_at_mme.lbi != linked_eps_bearer_id) {
    return false;
  }
  if (!(cb_req_rcvd_at_mme.s1_u_sgw_fteid.teid)) {
    return false;
  }
  if ((memcmp(
          &cb_req_rcvd_at_mme.tft, &tft, sizeof(traffic_flow_template_t)))) {
    return false;
  }
  return true;
}

TEST_F(SgwS8ConfigAndCreateMock, send_create_bearer_req_to_mme) {
  ASSERT_EQ(task_zmq_ctx_main_s8.ready, true);
  sgw_eps_bearer_context_information_t* sgw_pdn_session = NULL;
  uint32_t temporary_create_session_procedure_id        = 0;

  sgw_state_t* sgw_state = create_and_get_contexts_on_cs_req(
      &temporary_create_session_procedure_id, &sgw_pdn_session);
  s8_create_session_response_t csresp = {0};
  fill_itti_csrsp(&csresp, temporary_create_session_procedure_id);
  sgw_s8_handle_create_session_response(sgw_state, &csresp, imsi64);

  s8_create_bearer_request_t cb_req = {0};
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);
  fill_create_bearer_request(
      &cb_req, csresp.context_teid, default_eps_bearer_id);
  gtpv2c_cause_value_t cause_value = REQUEST_REJECTED;

  EXPECT_CALL(
      *mme_app_handler,
      mme_app_handle_nw_init_ded_bearer_actv_req(check_params_in_cb_req(
          cb_req.linked_eps_bearer_id, cb_req.bearer_context[0].tft)))
      .Times(1)
      .WillOnce(ReturnFromAsyncTask(&cv));

  EXPECT_NE(
      sgw_s8_handle_create_bearer_request(sgw_state, &cb_req, &cause_value),
      INVALID_IMSI64);
  cv.wait_for(lock, std::chrono::milliseconds(END_OF_TESTCASE_SLEEP_MS));

  // Validates sequence number matches with received create bearer request
  pgw_ni_cbr_proc_t* pgw_ni_cbr_proc =
      pgw_get_procedure_create_bearer(sgw_pdn_session);
  EXPECT_TRUE(pgw_ni_cbr_proc != nullptr);

  bool is_seq_number_updated                             = false;
  sgw_eps_bearer_entry_wrapper_t* sgw_eps_bearer_entry_p = nullptr;
  LIST_FOREACH(
      sgw_eps_bearer_entry_p, pgw_ni_cbr_proc->pending_eps_bearers, entries) {
    if ((sgw_eps_bearer_entry_p) &&
        (sgw_eps_bearer_entry_p->sgw_eps_bearer_entry->sgw_sequence_number ==
         cb_req.sequence_number)) {
      is_seq_number_updated = true;
      break;
    }
  }
  EXPECT_TRUE(is_seq_number_updated == true);

  free_wrapper(reinterpret_cast<void**>(&cb_req.pgw_cp_address));
}

TEST_F(SgwS8ConfigAndCreateMock, recv_create_bearer_response) {
  ASSERT_EQ(task_zmq_ctx_main_s8.ready, true);
  sgw_eps_bearer_context_information_t* sgw_pdn_session = NULL;
  uint32_t temporary_create_session_procedure_id        = 0;

  sgw_state_t* sgw_state = create_and_get_contexts_on_cs_req(
      &temporary_create_session_procedure_id, &sgw_pdn_session);
  s8_create_session_response_t csresp = {0};
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);
  fill_itti_csrsp(&csresp, temporary_create_session_procedure_id);
  sgw_s8_handle_create_session_response(sgw_state, &csresp, imsi64);

  s8_create_bearer_request_t cb_req = {0};
  fill_create_bearer_request(
      &cb_req, csresp.context_teid, default_eps_bearer_id);
  gtpv2c_cause_value_t cause_value = REQUEST_REJECTED;

  EXPECT_CALL(
      *mme_app_handler,
      mme_app_handle_nw_init_ded_bearer_actv_req(check_params_in_cb_req(
          cb_req.linked_eps_bearer_id, cb_req.bearer_context[0].tft)))
      .Times(1)
      .WillOnce(ReturnFromAsyncTask(&cv));
  EXPECT_NE(
      sgw_s8_handle_create_bearer_request(sgw_state, &cb_req, &cause_value),
      INVALID_IMSI64);
  cv.wait_for(lock, std::chrono::milliseconds(END_OF_TESTCASE_SLEEP_MS));

  // Validates sequence number matches with received create bearer request
  pgw_ni_cbr_proc_t* pgw_ni_cbr_proc =
      pgw_get_procedure_create_bearer(sgw_pdn_session);
  EXPECT_TRUE(pgw_ni_cbr_proc != nullptr);

  bool is_seq_number_updated                             = false;
  sgw_eps_bearer_entry_wrapper_t* sgw_eps_bearer_entry_p = nullptr;
  LIST_FOREACH(
      sgw_eps_bearer_entry_p, pgw_ni_cbr_proc->pending_eps_bearers, entries) {
    if ((sgw_eps_bearer_entry_p) &&
        (sgw_eps_bearer_entry_p->sgw_eps_bearer_entry->sgw_sequence_number ==
         cb_req.sequence_number)) {
      is_seq_number_updated = true;
      break;
    }
  }
  uint32_t s1_u_sgw_fteid = 0;
  if (is_seq_number_updated) {
    s1_u_sgw_fteid =
        sgw_eps_bearer_entry_p->sgw_eps_bearer_entry->s_gw_teid_S1u_S12_S4_up;
  }

  itti_s11_nw_init_actv_bearer_rsp_t s11_actv_bearer_rsp;
  memset(&s11_actv_bearer_rsp, 0, sizeof(itti_s11_nw_init_actv_bearer_rsp_t));
  fill_create_bearer_response(
      &s11_actv_bearer_rsp, csresp.context_teid, 6, s1_u_sgw_fteid,
      REQUEST_ACCEPTED);
  sgw_s8_handle_s11_create_bearer_response(
      sgw_state, &s11_actv_bearer_rsp, imsi64);

  // On successful creation of dedicated bearer, there shall be no pending
  // create bearer procedures
  pgw_ni_cbr_proc_t* pgw_ni_cbr_proc_after =
      pgw_get_procedure_create_bearer(sgw_pdn_session);
  EXPECT_TRUE(pgw_ni_cbr_proc_after == nullptr);
  uint8_t bearer_id_updated = false;
  for (uint8_t idx = 0; idx < BEARERS_PER_UE; idx++) {
    if (sgw_pdn_session->pdn_connection.sgw_eps_bearers_array[idx]
            ->eps_bearer_id ==
        s11_actv_bearer_rsp.bearer_contexts.bearer_contexts[0].eps_bearer_id) {
      bearer_id_updated = true;
      break;
    }
  }
  EXPECT_EQ(bearer_id_updated, true);
  free_wrapper(reinterpret_cast<void**>(&cb_req.pgw_cp_address));
}

TEST_F(SgwS8ConfigAndCreateMock, recv_delete_bearer_req) {
  ASSERT_EQ(task_zmq_ctx_main_s8.ready, true);
  sgw_eps_bearer_context_information_t* sgw_pdn_session = NULL;
  uint32_t temporary_create_session_procedure_id        = 0;
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  sgw_state_t* sgw_state = create_and_get_contexts_on_cs_req(
      &temporary_create_session_procedure_id, &sgw_pdn_session);
  s8_create_session_response_t csresp = {0};
  fill_itti_csrsp(&csresp, temporary_create_session_procedure_id);
  sgw_s8_handle_create_session_response(sgw_state, &csresp, imsi64);

  s8_create_bearer_request_t cb_req = {0};
  fill_create_bearer_request(
      &cb_req, csresp.context_teid, default_eps_bearer_id);
  gtpv2c_cause_value_t cause_value = REQUEST_REJECTED;

  EXPECT_NE(
      sgw_s8_handle_create_bearer_request(sgw_state, &cb_req, &cause_value),
      INVALID_IMSI64);

  // Validates sequence number matches with received create bearer request
  pgw_ni_cbr_proc_t* pgw_ni_cbr_proc =
      pgw_get_procedure_create_bearer(sgw_pdn_session);
  EXPECT_TRUE(pgw_ni_cbr_proc != nullptr);

  bool is_seq_number_updated                             = false;
  sgw_eps_bearer_entry_wrapper_t* sgw_eps_bearer_entry_p = nullptr;
  LIST_FOREACH(
      sgw_eps_bearer_entry_p, pgw_ni_cbr_proc->pending_eps_bearers, entries) {
    if ((sgw_eps_bearer_entry_p) &&
        (sgw_eps_bearer_entry_p->sgw_eps_bearer_entry->sgw_sequence_number ==
         cb_req.sequence_number)) {
      is_seq_number_updated = true;
      break;
    }
  }
  uint32_t s1_u_sgw_fteid = 0;
  if (is_seq_number_updated) {
    s1_u_sgw_fteid =
        sgw_eps_bearer_entry_p->sgw_eps_bearer_entry->s_gw_teid_S1u_S12_S4_up;
  }

  itti_s11_nw_init_actv_bearer_rsp_t s11_actv_bearer_rsp;
  memset(&s11_actv_bearer_rsp, 0, sizeof(itti_s11_nw_init_actv_bearer_rsp_t));
  fill_create_bearer_response(
      &s11_actv_bearer_rsp, csresp.context_teid, 6, s1_u_sgw_fteid,
      REQUEST_ACCEPTED);
  sgw_s8_handle_s11_create_bearer_response(
      sgw_state, &s11_actv_bearer_rsp, imsi64);

  s8_delete_bearer_request_t db_req = {0};
  fill_delete_bearer_request(
      &db_req, csresp.context_teid,
      s11_actv_bearer_rsp.bearer_contexts.bearer_contexts[0].eps_bearer_id);

  EXPECT_CALL(
      *mme_app_handler,
      mme_app_handle_nw_init_bearer_deactv_req(check_params_in_db_req(
          db_req.num_eps_bearer_id, db_req.eps_bearer_id)))
      .Times(1)
      .WillOnce(ReturnFromAsyncTask(&cv));
  EXPECT_EQ(sgw_s8_handle_delete_bearer_request(sgw_state, &db_req), RETURNok);
  cv.wait_for(lock, std::chrono::milliseconds(END_OF_TESTCASE_SLEEP_MS));
  free_wrapper(reinterpret_cast<void**>(&cb_req.pgw_cp_address));
}

TEST_F(SgwS8ConfigAndCreateMock, delete_bearer_response) {
  ASSERT_EQ(task_zmq_ctx_main_s8.ready, true);
  sgw_eps_bearer_context_information_t* sgw_pdn_session = NULL;
  uint32_t temporary_create_session_procedure_id        = 0;

  sgw_state_t* sgw_state = create_and_get_contexts_on_cs_req(
      &temporary_create_session_procedure_id, &sgw_pdn_session);
  s8_create_session_response_t csresp = {0};
  fill_itti_csrsp(&csresp, temporary_create_session_procedure_id);
  sgw_s8_handle_create_session_response(sgw_state, &csresp, imsi64);

  s8_create_bearer_request_t cb_req = {0};
  fill_create_bearer_request(
      &cb_req, csresp.context_teid, default_eps_bearer_id);
  gtpv2c_cause_value_t cause_value = REQUEST_REJECTED;

  EXPECT_NE(
      sgw_s8_handle_create_bearer_request(sgw_state, &cb_req, &cause_value),
      INVALID_IMSI64);

  // Validates sequence number matches with received create bearer request
  pgw_ni_cbr_proc_t* pgw_ni_cbr_proc =
      pgw_get_procedure_create_bearer(sgw_pdn_session);
  EXPECT_TRUE(pgw_ni_cbr_proc != nullptr);

  bool is_seq_number_updated                             = false;
  sgw_eps_bearer_entry_wrapper_t* sgw_eps_bearer_entry_p = nullptr;
  LIST_FOREACH(
      sgw_eps_bearer_entry_p, pgw_ni_cbr_proc->pending_eps_bearers, entries) {
    if ((sgw_eps_bearer_entry_p) &&
        (sgw_eps_bearer_entry_p->sgw_eps_bearer_entry->sgw_sequence_number ==
         cb_req.sequence_number)) {
      is_seq_number_updated = true;
      break;
    }
  }
  uint32_t s1_u_sgw_fteid = 0;
  if (is_seq_number_updated) {
    s1_u_sgw_fteid =
        sgw_eps_bearer_entry_p->sgw_eps_bearer_entry->s_gw_teid_S1u_S12_S4_up;
  }

  itti_s11_nw_init_actv_bearer_rsp_t s11_actv_bearer_rsp;
  memset(&s11_actv_bearer_rsp, 0, sizeof(itti_s11_nw_init_actv_bearer_rsp_t));
  fill_create_bearer_response(
      &s11_actv_bearer_rsp, csresp.context_teid, 6, s1_u_sgw_fteid,
      REQUEST_ACCEPTED);
  sgw_s8_handle_s11_create_bearer_response(
      sgw_state, &s11_actv_bearer_rsp, imsi64);

  s8_delete_bearer_request_t db_req = {0};
  fill_delete_bearer_request(
      &db_req, csresp.context_teid,
      s11_actv_bearer_rsp.bearer_contexts.bearer_contexts[0].eps_bearer_id);

  EXPECT_EQ(sgw_s8_handle_delete_bearer_request(sgw_state, &db_req), RETURNok);

  itti_s11_nw_init_deactv_bearer_rsp_t s11_delete_bearer_response = {0};

  fill_delete_bearer_response(
      &s11_delete_bearer_response, csresp.context_teid, db_req.eps_bearer_id[0],
      REQUEST_ACCEPTED);
  EXPECT_EQ(
      sgw_s8_handle_s11_delete_bearer_response(
          sgw_state, &s11_delete_bearer_response, imsi64),
      RETURNok);
  free_wrapper(reinterpret_cast<void**>(&cb_req.pgw_cp_address));
}
