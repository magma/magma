/**
 * Copyright 2022 The Magma Authors.
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

#include "lte/gateway/c/core/oai/test/spgw_task/spgw_procedures_test_fixture.hpp"

#include <gmock/gmock.h>
#include <gtest/gtest.h>
#include <netinet/in.h>

extern "C" {
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/common/queue.h"
#include "lte/gateway/c/core/oai/include/gx_messages_types.h"
#include "lte/gateway/c/core/oai/include/ip_forward_messages_types.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.401.h"
}

#include "lte/gateway/c/core/oai/tasks/sgw/sgw_defs.hpp"
#include "lte/gateway/c/core/oai/test/spgw_task/spgw_test_util.h"
#include "lte/gateway/c/core/oai/include/sgw_context_manager.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/pgw_handlers.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/pgw_procedures.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/sgw_handlers.hpp"

namespace magma {
namespace lte {
task_zmq_ctx_t task_zmq_ctx_main_spgw;

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    default: {
    } break;
  }

  itti_free_msg_content(received_message_p);
  free(received_message_p);
  return 0;
}

void SPGWAppProcedureTest::SetUp() {
  // setup mock MME app task
  mme_app_handler = std::make_shared<MockMmeAppHandler>();
  itti_init(TASK_MAX, THREAD_MAX, MESSAGES_ID_MAX, tasks_info, messages_info,
            NULL, NULL);

  // initialize configs
  mme_config_init(&mme_config);
  spgw_config_init(&spgw_config);
  create_partial_lists(&mme_config);
  mme_config.use_stateless = false;
  hss_associated = true;

  task_id_t task_id_list[2] = {TASK_MME_APP, TASK_SPGW_APP};
  init_task_context(TASK_MAIN, task_id_list, 2, handle_message,
                    &task_zmq_ctx_main_spgw);

  std::thread task_mme_app(start_mock_mme_app_task, mme_app_handler);
  task_mme_app.detach();

  std::cout << "Running setup" << std::endl;

  // initialize the SPGW task
  spgw_app_init(&spgw_config, mme_config.use_stateless);
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
  std::cout << "Setup done" << std::endl;
}

void SPGWAppProcedureTest::TearDown() {
  send_terminate_message_fatal(&task_zmq_ctx_main_spgw);
  destroy_task_context(&task_zmq_ctx_main_spgw);
  itti_free_desc_threads();

  free_mme_config(&mme_config);
  free_spgw_config(&spgw_config);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

teid_t SPGWAppProcedureTest::create_default_session(spgw_state_t* spgw_state) {
  status_code_e return_code = RETURNerror;
  itti_s11_create_session_request_t sample_session_req_p = {};
  fill_create_session_request(&sample_session_req_p, test_imsi_str,
                              DEFAULT_MME_S11_TEID, DEFAULT_BEARER_INDEX,
                              sample_default_bearer_context, test_plmn);

  // trigger create session req to SPGW
  return_code = sgw_handle_s11_create_session_request(
      spgw_state, &sample_session_req_p, test_imsi64);

  EXPECT_EQ(return_code, RETURNok);

  // Verify that a UE context exists in SPGW state after CSR is received
  spgw_ue_context_t* ue_context_p = spgw_get_ue_context(test_imsi64);
  EXPECT_TRUE(ue_context_p != nullptr);

  // Verify that teid is created
  EXPECT_FALSE(LIST_EMPTY(&ue_context_p->sgw_s11_teid_list));
  teid_t ue_sgw_teid =
      LIST_FIRST(&ue_context_p->sgw_s11_teid_list)->sgw_s11_teid;

  // Verify that no IP address is allocated for this UE
  magma::lte::oai::S11BearerContext* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt;
  magma::proto_map_rc_t rc = sgw_cm_get_eps_bearer_entry(
      spgw_eps_bearer_ctxt_info_p->mutable_sgw_eps_bearer_context()
          ->mutable_pdn_connection(),
      DEFAULT_EPS_BEARER_ID, &eps_bearer_ctxt);
  EXPECT_TRUE(eps_bearer_ctxt.ue_ip_paa().ipv4_addr().size() ==
              UNASSIGNED_UE_IP);

  // send an IP alloc response to SPGW
  itti_ip_allocation_response_t test_ip_alloc_resp = {};
  fill_ip_allocation_response(&test_ip_alloc_resp, SGI_STATUS_OK, ue_sgw_teid,
                              DEFAULT_EPS_BEARER_ID, DEFAULT_UE_IP,
                              DEFAULT_VLAN);
  return_code = sgw_handle_ip_allocation_rsp(spgw_state, &test_ip_alloc_resp,
                                             test_imsi64);

  EXPECT_EQ(return_code, RETURNok);

  // check if IP address is allocated after this message is done
  sgw_cm_get_eps_bearer_entry(
      spgw_eps_bearer_ctxt_info_p->mutable_sgw_eps_bearer_context()
          ->mutable_pdn_connection(),
      DEFAULT_EPS_BEARER_ID, &eps_bearer_ctxt);
  struct in_addr ue_ipv4 = {};
  uint32_t ue_ip = DEFAULT_UE_IP;
  inet_pton(AF_INET, eps_bearer_ctxt.ue_ip_paa().ipv4_addr().c_str(), &ue_ipv4);
  EXPECT_TRUE(!(memcmp(&ue_ipv4, &ue_ip, sizeof(DEFAULT_UE_IP))));

  // send pcef create session response to SPGW
  itti_pcef_create_session_response_t sample_pcef_csr_resp;
  fill_pcef_create_session_response(&sample_pcef_csr_resp, PCEF_STATUS_OK,
                                    ue_sgw_teid, DEFAULT_EPS_BEARER_ID,
                                    SGI_STATUS_OK);

  // check if MME gets a create session response
  EXPECT_CALL(*mme_app_handler, mme_app_handle_create_sess_resp()).Times(1);

  spgw_handle_pcef_create_session_response(spgw_state, &sample_pcef_csr_resp,
                                           test_imsi64);

  // create sample modify default bearer request
  itti_s11_modify_bearer_request_t sample_modify_bearer_req = {};
  fill_modify_bearer_request(&sample_modify_bearer_req, DEFAULT_MME_S11_TEID,
                             ue_sgw_teid, DEFAULT_ENB_GTP_TEID,
                             DEFAULT_BEARER_INDEX, DEFAULT_EPS_BEARER_ID);

  EXPECT_CALL(*mme_app_handler, mme_app_handle_modify_bearer_rsp()).Times(1);
  return_code =
      sgw_handle_modify_bearer_request(&sample_modify_bearer_req, test_imsi64);

  EXPECT_EQ(return_code, RETURNok);

  // verify that exactly one session exists in SPGW state
  EXPECT_TRUE(is_num_ue_contexts_valid(1));
  EXPECT_TRUE(is_num_cp_teids_valid(test_imsi64, 1));

  // verify that eNB address information exists
  EXPECT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 1));

  return ue_sgw_teid;
}

ebi_t SPGWAppProcedureTest ::activate_dedicated_bearer(
    spgw_state_t* spgw_state,
    magma::lte::oai::S11BearerContext* spgw_eps_bearer_ctxt_info_p,
    teid_t ue_sgw_teid) {
  status_code_e return_code = RETURNerror;
  // send network initiated dedicated bearer activation request from Session
  // Manager
  itti_gx_nw_init_actv_bearer_request_t sample_gx_nw_init_ded_bearer_actv_req =
      {};
  gtpv2c_cause_value_t failed_cause = REQUEST_ACCEPTED;
  fill_nw_initiated_activate_bearer_request(
      &sample_gx_nw_init_ded_bearer_actv_req, test_imsi_str,
      DEFAULT_EPS_BEARER_ID, sample_dedicated_bearer_qos);

  // check that MME gets a bearer activation request
  EXPECT_CALL(*mme_app_handler,
              mme_app_handle_nw_init_ded_bearer_actv_req(
                  check_params_in_actv_bearer_req(
                      sample_gx_nw_init_ded_bearer_actv_req.lbi,
                      sample_gx_nw_init_ded_bearer_actv_req.ul_tft)))
      .Times(1);

  return_code = spgw_handle_nw_initiated_bearer_actv_req(
      spgw_state, &sample_gx_nw_init_ded_bearer_actv_req, test_imsi64,
      &failed_cause);

  EXPECT_EQ(return_code, RETURNok);

  magma::lte::oai::SgwEpsBearerContextInfo* sgw_context_p =
      spgw_eps_bearer_ctxt_info_p->mutable_sgw_eps_bearer_context();
  // check number of pending procedures
  EXPECT_EQ(get_num_pending_create_bearer_procedures(sgw_context_p), 1);

  // fetch new SGW teid for the pending bearer procedure
  teid_t ue_ded_bearer_sgw_teid = 0;
  for (int proc_index = 0;
       proc_index < sgw_context_p->pending_procedures_size(); proc_index++) {
    magma::lte::oai::PgwCbrProcedure* pgw_ni_cbr_proc =
        sgw_context_p->mutable_pending_procedures(proc_index);
    EXPECT_TRUE(pgw_ni_cbr_proc->type() ==
                PGW_BASE_PROC_TYPE_NETWORK_INITATED_CREATE_BEARER_REQUEST);
    for (uint8_t bearer_index = 0;
         bearer_index < pgw_ni_cbr_proc->pending_eps_bearers_size();
         bearer_index++) {
      magma::lte::oai::SgwEpsBearerContext* bearer_context_proto =
          pgw_ni_cbr_proc->mutable_pending_eps_bearers(bearer_index);
      ue_ded_bearer_sgw_teid = bearer_context_proto->sgw_teid_s1u_s12_s4_up();
    }
  }

  // send bearer activation response from MME
  ebi_t ded_eps_bearer_id = DEFAULT_EPS_BEARER_ID + 1;
  itti_s11_nw_init_actv_bearer_rsp_t sample_nw_init_ded_bearer_actv_resp = {};
  fill_nw_initiated_activate_bearer_response(
      &sample_nw_init_ded_bearer_actv_resp, DEFAULT_MME_S11_TEID, ue_sgw_teid,
      ue_ded_bearer_sgw_teid, DEFAULT_ENB_GTP_TEID + 1, ded_eps_bearer_id,
      REQUEST_ACCEPTED, test_plmn);
  return_code = sgw_handle_nw_initiated_actv_bearer_rsp(
      &sample_nw_init_ded_bearer_actv_resp, test_imsi64);

  EXPECT_EQ(return_code, RETURNok);

  // check that bearer is created
  EXPECT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 2));

  // check that no pending procedure is left
  EXPECT_EQ(get_num_pending_create_bearer_procedures(
                spgw_eps_bearer_ctxt_info_p->mutable_sgw_eps_bearer_context()),
            0);
  return ded_eps_bearer_id;
}

void SPGWAppProcedureTest ::deactivate_dedicated_bearer(
    spgw_state_t* spgw_state, teid_t ue_sgw_teid, ebi_t ded_eps_bearer_id) {
  status_code_e return_code = RETURNerror;
  // send deactivate request for dedicated bearer from Session Manager
  itti_gx_nw_init_deactv_bearer_request_t
      sample_gx_nw_init_ded_bearer_deactv_req = {};
  fill_nw_initiated_deactivate_bearer_request(
      &sample_gx_nw_init_ded_bearer_deactv_req, test_imsi_str,
      DEFAULT_EPS_BEARER_ID, ded_eps_bearer_id);

  return_code = spgw_handle_nw_initiated_bearer_deactv_req(
      &sample_gx_nw_init_ded_bearer_deactv_req, test_imsi64);

  EXPECT_EQ(return_code, RETURNok);

  // send a delete dedicated bearer response from MME
  itti_s11_nw_init_deactv_bearer_rsp_t sample_nw_init_ded_bearer_deactv_resp =
      {};
  int num_bearers_to_delete = 1;
  ebi_t eps_bearer_id_array[] = {ded_eps_bearer_id};

  fill_nw_initiated_deactivate_bearer_response(
      &sample_nw_init_ded_bearer_deactv_resp, test_imsi64, false,
      REQUEST_ACCEPTED, eps_bearer_id_array, num_bearers_to_delete,
      ue_sgw_teid);
  return_code = sgw_handle_nw_initiated_deactv_bearer_rsp(
      spgw_state, &sample_nw_init_ded_bearer_deactv_resp, test_imsi64);
  EXPECT_EQ(return_code, RETURNok);

  // check that bearer is deleted
  EXPECT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 1));
}
}  // namespace lte
}  // namespace magma
