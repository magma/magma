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

#include <unistd.h>
#include <sys/types.h>
#include <pwd.h>

#include <gtest/gtest.h>
#include <string>
#include <thread>
#include <cstdlib>

#include "lte/gateway/c/core/oai/test/mock_tasks/mock_tasks.hpp"
#include "lte/gateway/c/core/oai/test/spgw_task/spgw_test_util.h"
#include "lte/gateway/c/core/oai/test/spgw_task/mock_spgw_op.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/pgw_handlers.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/sgw_defs.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/sgw_handlers.hpp"
#include "lte/gateway/c/core/oai/include/mme_config.hpp"
#include "lte/gateway/c/core/oai/include/sgw_context_manager.hpp"
#include "lte/gateway/c/core/oai/include/spgw_state.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/spgw_state_converter.hpp"

extern "C" {
#include "lte/gateway/c/core/oai/include/spgw_config.h"
}

extern bool hss_associated;

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

MATCHER_P2(check_params_in_actv_bearer_req, lbi, tft, "") {
  auto cb_req_rcvd_at_mme =
      static_cast<itti_s11_nw_init_actv_bearer_request_t>(arg);
  if (cb_req_rcvd_at_mme.lbi != lbi) {
    return false;
  }
  if (!(cb_req_rcvd_at_mme.s1_u_sgw_fteid.teid)) {
    return false;
  }
  if ((memcmp(&cb_req_rcvd_at_mme.tft, &tft,
              sizeof(traffic_flow_template_t)))) {
    return false;
  }
  return true;
}

MATCHER_P2(check_params_in_deactv_bearer_req, num_bearers, eps_bearer_id_array,
           "") {
  auto db_req_rcvd_at_mme =
      static_cast<itti_s11_nw_init_deactv_bearer_request_t>(arg);
  if (db_req_rcvd_at_mme.no_of_bearers != num_bearers) {
    return false;
  }
  if (memcmp(db_req_rcvd_at_mme.ebi, eps_bearer_id_array,
             sizeof(db_req_rcvd_at_mme.ebi))) {
    return false;
  }
  return true;
}

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

class SPGWAppInjectedStateProcedureTest : public ::testing::Test {
  virtual void SetUp() {
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

    // add injection of state loaded in SPGW state manager
#ifdef CMAKE_BUILD  // if CMAKE is used the absolute path is needed
    std::string magma_root = std::getenv("MAGMA_ROOT");
    std::string path = magma_root + "/" + DEFAULT_SPGW_CONTEXT_DATA_PATH;
#else  // if bazel is used the relative path is used
    std::string path = DEFAULT_SPGW_CONTEXT_DATA_PATH;
#endif
    name_of_ue_samples =
        load_file_into_vector_of_line_content(path, path + "data_list.txt");
    mock_read_spgw_ue_state_db(name_of_ue_samples);

    std::this_thread::sleep_for(
        std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
    std::cout << "Setup done" << std::endl;
  }

  virtual void TearDown() {
    send_terminate_message_fatal(&task_zmq_ctx_main_spgw);
    destroy_task_context(&task_zmq_ctx_main_spgw);
    itti_free_desc_threads();

    free_mme_config(&mme_config);
    free_spgw_config(&spgw_config);

    // Sleep to ensure that messages are received and contexts are released
    std::this_thread::sleep_for(
        std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
  }

 protected:
  std::shared_ptr<MockMmeAppHandler> mme_app_handler;
  std::string test_imsi_str = "001010000000002";
  uint64_t test_imsi64 = 1010000000002;
  uint64_t test_imsi64_invalid = 1010000000001;

  plmn_t test_plmn = {.mcc_digit2 = 0,
                      .mcc_digit1 = 0,
                      .mnc_digit3 = 0x00,
                      .mcc_digit3 = 0,
                      .mnc_digit2 = 0,
                      .mnc_digit1 = 0};
  bearer_context_to_be_created_t sample_default_bearer_context = {
      .eps_bearer_id = 5,
      .bearer_level_qos = {
          .pci = 1,
          .pl = 15,
          .qci = 9,
          .gbr = {},
      }};

  bearer_qos_t sample_dedicated_bearer_qos = {
      .pci = 1, .pl = 15, .qci = 9, .gbr = {}};

  int test_mme_s11_teid = 4;
  int test_bearer_index = 5;
  in_addr_t test_ue_ip = 0x0d80a8c0;
  std::vector<std::string> name_of_ue_samples;
};

TEST_F(SPGWAppInjectedStateProcedureTest, TestDeleteSessionSuccess) {
  spgw_ue_context_t* ue_context_p = spgw_get_ue_context(test_imsi64);
  teid_t ue_sgw_teid =
      LIST_FIRST(&ue_context_p->sgw_s11_teid_list)->sgw_s11_teid;

  magma::lte::oai::S11BearerContext* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt;
  EXPECT_EQ((sgw_cm_get_eps_bearer_entry(
                spgw_eps_bearer_ctxt_info_p->mutable_sgw_eps_bearer_context()
                    ->mutable_pdn_connection(),
                DEFAULT_EPS_BEARER_ID, &eps_bearer_ctxt)),
            magma::PROTO_MAP_OK);
  struct in_addr ue_ipv4 = {};
  inet_pton(AF_INET, eps_bearer_ctxt.ue_ip_paa().ipv4_addr().c_str(), &ue_ipv4);
  ASSERT_TRUE(!(memcmp(&ue_ipv4, &test_ue_ip, sizeof(test_ue_ip))));

  // verify that exactly one session exists in SPGW state
  ASSERT_TRUE(is_num_ue_contexts_valid(name_of_ue_samples.size()));
  ASSERT_TRUE(is_num_cp_teids_valid(test_imsi64, 1));

  // verify that eNB address information exists
  ASSERT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 1));

  // create sample delete session request
  itti_s11_delete_session_request_t sample_delete_session_request = {};
  fill_delete_session_request(&sample_delete_session_request, test_mme_s11_teid,
                              ue_sgw_teid, DEFAULT_EPS_BEARER_ID, test_plmn);

  EXPECT_CALL(*mme_app_handler,
              mme_app_handle_delete_sess_rsp(
                  check_cause_in_ds_rsp(REQUEST_ACCEPTED, test_mme_s11_teid)))
      .Times(1);

  status_code_e return_code = sgw_handle_delete_session_request(
      &sample_delete_session_request, test_imsi64);
  ASSERT_EQ(return_code, RETURNok);

  // verify SPGW state is cleared
  ASSERT_TRUE(is_num_ue_contexts_valid(name_of_ue_samples.size() - 1));
  ASSERT_TRUE(is_num_cp_teids_valid(test_imsi64, 0));

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppInjectedStateProcedureTest, TestModifyBearerFailure) {
  status_code_e return_code = RETURNerror;

  // verify that sessions exist in SPGW state
  ASSERT_TRUE(is_num_ue_contexts_valid(name_of_ue_samples.size()));
  ASSERT_TRUE(is_num_cp_teids_valid(test_imsi64, 1));

  // create sample modify default bearer request
  itti_s11_modify_bearer_request_t sample_modify_bearer_req = {};
  fill_modify_bearer_request(&sample_modify_bearer_req, DEFAULT_MME_S11_TEID,
                             ERROR_SGW_S11_TEID, DEFAULT_ENB_GTP_TEID,
                             DEFAULT_BEARER_INDEX, DEFAULT_EPS_BEARER_ID);

  EXPECT_CALL(*mme_app_handler, mme_app_handle_modify_bearer_rsp()).Times(1);
  return_code =
      sgw_handle_modify_bearer_request(&sample_modify_bearer_req, test_imsi64);

  ASSERT_EQ(return_code, RETURNok);

  // verify that number of valid sessions do not change after the modify bearer
  // request
  ASSERT_TRUE(is_num_ue_contexts_valid(name_of_ue_samples.size()));
  ASSERT_TRUE(is_num_cp_teids_valid(test_imsi64, 1));

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppInjectedStateProcedureTest, TestReleaseBearerSuccess) {
  spgw_ue_context_t* ue_context_p = spgw_get_ue_context(test_imsi64);

  teid_t ue_sgw_teid =
      LIST_FIRST(&ue_context_p->sgw_s11_teid_list)->sgw_s11_teid;

  // verify that exactly one session exists in SPGW state
  ASSERT_TRUE(is_num_ue_contexts_valid(name_of_ue_samples.size()));
  ASSERT_TRUE(is_num_cp_teids_valid(test_imsi64, 1));

  // verify that eNB address information exists
  ASSERT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 1));

  // send release access bearer request
  itti_s11_release_access_bearers_request_t sample_release_bearer_req = {};
  fill_release_access_bearer_request(&sample_release_bearer_req,
                                     DEFAULT_MME_S11_TEID, ue_sgw_teid);

  EXPECT_CALL(*mme_app_handler, mme_app_handle_release_access_bearers_resp())
      .Times(1);

  sgw_handle_release_access_bearers_request(&sample_release_bearer_req,
                                            test_imsi64);

  // verify that eNB information has been cleared
  ASSERT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 0));

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppInjectedStateProcedureTest, TestReleaseBearerWithInvalidImsi64) {
  spgw_ue_context_t* ue_context_p = spgw_get_ue_context(test_imsi64);
  teid_t ue_sgw_teid =
      LIST_FIRST(&ue_context_p->sgw_s11_teid_list)->sgw_s11_teid;

  // verify that exactly one session exists in SPGW state
  ASSERT_TRUE(is_num_ue_contexts_valid(name_of_ue_samples.size()));
  ASSERT_TRUE(is_num_cp_teids_valid(test_imsi64, 1));

  // verify that eNB address information exists
  ASSERT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 1));

  // send release access bearer request
  itti_s11_release_access_bearers_request_t sample_release_bearer_req = {};
  fill_release_access_bearer_request(&sample_release_bearer_req,
                                     test_mme_s11_teid, ERROR_SGW_S11_TEID);

  EXPECT_CALL(*mme_app_handler, mme_app_handle_release_access_bearers_resp())
      .Times(1);

  // Send wrong IMSI so that spgw will not be able to fetch and delete
  // the context
  sgw_handle_release_access_bearers_request(&sample_release_bearer_req,
                                            test_imsi64_invalid);

  // verify that eNB information has not been cleared
  ASSERT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 1));

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppInjectedStateProcedureTest, TestDedicatedBearerActivation) {
  spgw_state_t* spgw_state = get_spgw_state(false);
  status_code_e return_code = RETURNerror;

  spgw_ue_context_t* ue_context_p = spgw_get_ue_context(test_imsi64);
  teid_t ue_sgw_teid =
      LIST_FIRST(&ue_context_p->sgw_s11_teid_list)->sgw_s11_teid;

  magma::lte::oai::S11BearerContext* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  // verify that exactly one session exists in SPGW state
  EXPECT_TRUE(is_num_ue_contexts_valid(name_of_ue_samples.size()));
  EXPECT_TRUE(is_num_cp_teids_valid(test_imsi64, 1));

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
    for (int bearer_index = 0;
         bearer_index < pgw_ni_cbr_proc->pending_eps_bearers_size();
         bearer_index++) {
      magma::lte::oai::SgwEpsBearerContext* bearer_context_proto =
          pgw_ni_cbr_proc->mutable_pending_eps_bearers(bearer_index);
      ue_ded_bearer_sgw_teid = bearer_context_proto->sgw_teid_s1u_s12_s4_up();
    }
  }

  // send bearer activation response from MME
  itti_s11_nw_init_actv_bearer_rsp_t sample_nw_init_ded_bearer_actv_resp = {};
  fill_nw_initiated_activate_bearer_response(
      &sample_nw_init_ded_bearer_actv_resp, test_mme_s11_teid, ue_sgw_teid,
      ue_ded_bearer_sgw_teid, DEFAULT_ENB_GTP_TEID + 1,
      DEFAULT_EPS_BEARER_ID + 1, REQUEST_ACCEPTED, test_plmn);
  return_code = sgw_handle_nw_initiated_actv_bearer_rsp(
      &sample_nw_init_ded_bearer_actv_resp, test_imsi64);

  EXPECT_EQ(return_code, RETURNok);

  // check that bearer is created
  EXPECT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 2));

  // check that no pending procedure is left
  EXPECT_EQ(get_num_pending_create_bearer_procedures(sgw_context_p), 0);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppInjectedStateProcedureTest, TestDedicatedBearerDeactivation) {
  spgw_state_t* spgw_state = get_spgw_state(false);
  status_code_e return_code = RETURNerror;

  spgw_ue_context_t* ue_context_p = spgw_get_ue_context(test_imsi64);
  teid_t ue_sgw_teid =
      LIST_FIRST(&ue_context_p->sgw_s11_teid_list)->sgw_s11_teid;

  magma::lte::oai::S11BearerContext* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  // verify that exactly one session exists in SPGW state
  EXPECT_TRUE(is_num_ue_contexts_valid(name_of_ue_samples.size()));
  EXPECT_TRUE(is_num_cp_teids_valid(test_imsi64, 1));

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
    magma::lte::oai::SgwEpsBearerContext* bearer_context_proto =
        pgw_ni_cbr_proc->mutable_pending_eps_bearers(0);
    ue_ded_bearer_sgw_teid = bearer_context_proto->sgw_teid_s1u_s12_s4_up();
  }

  // send bearer activation response from MME
  ebi_t ded_eps_bearer_id = DEFAULT_EPS_BEARER_ID + 1;
  itti_s11_nw_init_actv_bearer_rsp_t sample_nw_init_ded_bearer_actv_resp = {};
  fill_nw_initiated_activate_bearer_response(
      &sample_nw_init_ded_bearer_actv_resp, test_mme_s11_teid, ue_sgw_teid,
      ue_ded_bearer_sgw_teid, DEFAULT_ENB_GTP_TEID + 1, ded_eps_bearer_id,
      REQUEST_ACCEPTED, test_plmn);
  return_code = sgw_handle_nw_initiated_actv_bearer_rsp(
      &sample_nw_init_ded_bearer_actv_resp, test_imsi64);

  EXPECT_EQ(return_code, RETURNok);

  // check that bearer is created
  EXPECT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 2));

  // check that no pending procedure is left
  EXPECT_EQ(get_num_pending_create_bearer_procedures(sgw_context_p), 0);

  // send deactivate request for dedicated bearer from Session Manager
  itti_gx_nw_init_deactv_bearer_request_t
      sample_gx_nw_init_ded_bearer_deactv_req = {};
  fill_nw_initiated_deactivate_bearer_request(
      &sample_gx_nw_init_ded_bearer_deactv_req, test_imsi_str,
      DEFAULT_EPS_BEARER_ID, ded_eps_bearer_id);

  // check that MME gets a bearer deactivation request
  EXPECT_CALL(*mme_app_handler,
              mme_app_handle_nw_init_bearer_deactv_req(
                  check_params_in_deactv_bearer_req(
                      sample_gx_nw_init_ded_bearer_deactv_req.no_of_bearers,
                      sample_gx_nw_init_ded_bearer_deactv_req.ebi)))
      .Times(1);

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

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppInjectedStateProcedureTest,
       TestDedicatedBearerDeactivationDeleteDefaultBearer) {
  spgw_state_t* spgw_state = get_spgw_state(false);
  status_code_e return_code = RETURNerror;

  spgw_ue_context_t* ue_context_p = spgw_get_ue_context(test_imsi64);
  teid_t ue_sgw_teid =
      LIST_FIRST(&ue_context_p->sgw_s11_teid_list)->sgw_s11_teid;
  magma::lte::oai::S11BearerContext* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  // verify that exactly one session exists in SPGW state
  EXPECT_TRUE(is_num_ue_contexts_valid(name_of_ue_samples.size()));
  EXPECT_TRUE(is_num_cp_teids_valid(test_imsi64, 1));

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
  for (uint8_t proc_index = 0;
       proc_index < sgw_context_p->pending_procedures_size(); proc_index++) {
    magma::lte::oai::PgwCbrProcedure pgw_ni_cbr_proc =
        sgw_context_p->pending_procedures(proc_index);
    EXPECT_TRUE(pgw_ni_cbr_proc.type() ==
                PGW_BASE_PROC_TYPE_NETWORK_INITATED_CREATE_BEARER_REQUEST);
    for (uint8_t bearer_index = 0;
         bearer_index < pgw_ni_cbr_proc.pending_eps_bearers_size();
         bearer_index++) {
      magma::lte::oai::SgwEpsBearerContext bearer_context_proto =
          pgw_ni_cbr_proc.pending_eps_bearers(bearer_index);
      ue_ded_bearer_sgw_teid = bearer_context_proto.sgw_teid_s1u_s12_s4_up();
    }
  }
  // send bearer activation response from MME
  ebi_t ded_eps_bearer_id = DEFAULT_EPS_BEARER_ID + 1;
  itti_s11_nw_init_actv_bearer_rsp_t sample_nw_init_ded_bearer_actv_resp = {};
  fill_nw_initiated_activate_bearer_response(
      &sample_nw_init_ded_bearer_actv_resp, test_mme_s11_teid, ue_sgw_teid,
      ue_ded_bearer_sgw_teid, DEFAULT_ENB_GTP_TEID + 1, ded_eps_bearer_id,
      REQUEST_ACCEPTED, test_plmn);
  return_code = sgw_handle_nw_initiated_actv_bearer_rsp(
      &sample_nw_init_ded_bearer_actv_resp, test_imsi64);

  EXPECT_EQ(return_code, RETURNok);

  // check that bearer is created
  EXPECT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 2));

  // check that no pending procedure is left
  EXPECT_EQ(get_num_pending_create_bearer_procedures(sgw_context_p), 0);

  // send deactivate request for dedicated bearer from Session Manager
  itti_gx_nw_init_deactv_bearer_request_t
      sample_gx_nw_init_ded_bearer_deactv_req = {};
  fill_nw_initiated_deactivate_bearer_request(
      &sample_gx_nw_init_ded_bearer_deactv_req, test_imsi_str,
      DEFAULT_EPS_BEARER_ID, ded_eps_bearer_id);

  // check that MME gets a bearer deactivation request
  EXPECT_CALL(*mme_app_handler,
              mme_app_handle_nw_init_bearer_deactv_req(
                  check_params_in_deactv_bearer_req(
                      sample_gx_nw_init_ded_bearer_deactv_req.no_of_bearers,
                      sample_gx_nw_init_ded_bearer_deactv_req.ebi)))
      .Times(1);

  return_code = spgw_handle_nw_initiated_bearer_deactv_req(
      &sample_gx_nw_init_ded_bearer_deactv_req, test_imsi64);

  EXPECT_EQ(return_code, RETURNok);

  // send a delete dedicated bearer response from MME
  itti_s11_nw_init_deactv_bearer_rsp_t sample_nw_init_ded_bearer_deactv_resp =
      {};
  int num_bearers_to_delete = 2;
  ebi_t eps_bearer_id_array[] = {DEFAULT_EPS_BEARER_ID, ded_eps_bearer_id};

  fill_nw_initiated_deactivate_bearer_response(
      &sample_nw_init_ded_bearer_deactv_resp, test_imsi64, true,
      REQUEST_ACCEPTED, eps_bearer_id_array, num_bearers_to_delete,
      ue_sgw_teid);
  return_code = sgw_handle_nw_initiated_deactv_bearer_rsp(
      spgw_state, &sample_nw_init_ded_bearer_deactv_resp, test_imsi64);
  EXPECT_EQ(return_code, RETURNok);

  // check that session is removed
  EXPECT_TRUE(is_num_ue_contexts_valid(name_of_ue_samples.size() - 1));
  EXPECT_TRUE(is_num_cp_teids_valid(test_imsi64, 0));

  free(sample_nw_init_ded_bearer_deactv_resp.lbi);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}
}  // namespace lte
}  // namespace magma
