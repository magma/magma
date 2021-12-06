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
#include <cstdint>
#include <string>
#include <thread>

#include "lte/gateway/c/core/oai/test/mock_tasks/mock_tasks.h"
#include "lte/gateway/c/core/oai/test/spgw_task/spgw_test_util.h"

extern "C" {
#include "lte/gateway/c/core/oai/include/mme_config.h"
#include "lte/gateway/c/core/oai/tasks/sgw/pgw_handlers.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_29.274.h"
#include "lte/gateway/c/core/oai/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/lib/hashtable/hashtable.h"
#include "lte/gateway/c/core/oai/include/s11_messages_types.h"
#include "lte/gateway/c/core/oai/tasks/sgw/sgw_defs.h"
#include "lte/gateway/c/core/oai/tasks/sgw/sgw_handlers.h"
#include "lte/gateway/c/core/oai/include/spgw_config.h"
#include "lte/gateway/c/core/oai/include/spgw_types.h"
}

extern bool hss_associated;

namespace magma {
namespace lte {

task_zmq_ctx_t task_zmq_ctx_main_spgw;

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    default: { } break; }

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
  if ((memcmp(
          &cb_req_rcvd_at_mme.tft, &tft, sizeof(traffic_flow_template_t)))) {
    return false;
  }
  return true;
}

MATCHER_P2(
    check_params_in_deactv_bearer_req, num_bearers, eps_bearer_id_array, "") {
  auto db_req_rcvd_at_mme =
      static_cast<itti_s11_nw_init_deactv_bearer_request_t>(arg);
  if (db_req_rcvd_at_mme.no_of_bearers != num_bearers) {
    return false;
  }
  if (memcmp(
          db_req_rcvd_at_mme.ebi, eps_bearer_id_array,
          sizeof(db_req_rcvd_at_mme.ebi))) {
    return false;
  }
  return true;
}

MATCHER_P2(check_params_in_suspend_ack, return_val, teid, "") {
  auto suspend_ack_rcvd_at_mme =
      static_cast<itti_s11_suspend_acknowledge_t>(arg);
  if ((suspend_ack_rcvd_at_mme.cause.cause_value == return_val) &&
      (suspend_ack_rcvd_at_mme.teid == teid)) {
    return true;
  }
  return false;
}

class SPGWAppProcedureTest : public ::testing::Test {
  virtual void SetUp() {
    // setup mock MME app task
    mme_app_handler = std::make_shared<MockMmeAppHandler>();
    itti_init(
        TASK_MAX, THREAD_MAX, MESSAGES_ID_MAX, tasks_info, messages_info, NULL,
        NULL);

    // initialize configs
    mme_config_init(&mme_config);
    spgw_config_init(&spgw_config);
    create_partial_lists(&mme_config);
    mme_config.use_stateless = false;
    hss_associated           = true;

    task_id_t task_id_list[2] = {TASK_MME_APP, TASK_SPGW_APP};
    init_task_context(
        TASK_MAIN, task_id_list, 2, handle_message, &task_zmq_ctx_main_spgw);

    std::thread task_mme_app(start_mock_mme_app_task, mme_app_handler);
    task_mme_app.detach();

    std::cout << "Running setup" << std::endl;

    // initialize the SPGW task
    spgw_app_init(&spgw_config, mme_config.use_stateless);
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
  std::string test_imsi_str = "001010000000001";
  uint64_t test_imsi64      = 1010000000001;
  plmn_t test_plmn          = {.mcc_digit2 = 0,
                      .mcc_digit1 = 0,
                      .mnc_digit3 = 0x0f,
                      .mcc_digit3 = 1,
                      .mnc_digit2 = 1,
                      .mnc_digit1 = 0};
  bearer_context_to_be_created_t sample_default_bearer_context = {
      .eps_bearer_id    = 5,
      .bearer_level_qos = {.pci = 1,
                           .pl  = 15,
                           .pvi = 0,
                           .qci = 9,
                           .gbr = {},
                           .mbr = {.br_ul = 200000000, .br_dl = 100000000}}};

  bearer_qos_t sample_dedicated_bearer_qos = {
      .pci = 1,
      .pl  = 1,
      .pvi = 0,
      .qci = 1,
      .gbr = {.br_ul = 200000000, .br_dl = 100000000},
      .mbr = {.br_ul = 200000000, .br_dl = 100000000}};
};

TEST_F(SPGWAppProcedureTest, TestCreateSessionSuccess) {
  spgw_state_t* spgw_state  = get_spgw_state(false);
  status_code_e return_code = RETURNerror;
  // expect call to MME create session response
  itti_s11_create_session_request_t sample_session_req_p = {};
  fill_create_session_request(
      &sample_session_req_p, test_imsi_str, DEFAULT_MME_S11_TEID,
      DEFAULT_BEARER_INDEX, sample_default_bearer_context, test_plmn);

  // trigger create session req to SPGW
  return_code = sgw_handle_s11_create_session_request(
      spgw_state, &sample_session_req_p, test_imsi64);

  EXPECT_EQ(return_code, RETURNok);

  // Verify that a UE context exists in SPGW state after CSR is received
  spgw_ue_context_t* ue_context_p = spgw_get_ue_context(test_imsi64);
  ASSERT_TRUE(ue_context_p != nullptr);

  // Verify that teid is created
  ASSERT_FALSE(LIST_EMPTY(&ue_context_p->sgw_s11_teid_list));
  teid_t ue_sgw_teid =
      LIST_FIRST(&ue_context_p->sgw_s11_teid_list)->sgw_s11_teid;

  // Verify that no IP address is allocated for this UE
  s_plus_p_gw_eps_bearer_context_information_t* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p = sgw_cm_get_eps_bearer_entry(
      &spgw_eps_bearer_ctxt_info_p->sgw_eps_bearer_context_information
           .pdn_connection,
      DEFAULT_EPS_BEARER_ID);

  ASSERT_TRUE(eps_bearer_ctxt_p->paa.ipv4_address.s_addr == UNASSIGNED_UE_IP);

  // send an IP alloc response to SPGW
  itti_ip_allocation_response_t test_ip_alloc_resp = {};
  fill_ip_allocation_response(
      &test_ip_alloc_resp, SGI_STATUS_OK, ue_sgw_teid, DEFAULT_EPS_BEARER_ID,
      DEFAULT_UE_IP, DEFAULT_VLAN);
  return_code = sgw_handle_ip_allocation_rsp(
      spgw_state, &test_ip_alloc_resp, test_imsi64);

  ASSERT_EQ(return_code, RETURNok);

  // check if IP address is allocated after this message is done
  ASSERT_TRUE(eps_bearer_ctxt_p->paa.ipv4_address.s_addr == DEFAULT_UE_IP);

  // send pcef create session response to SPGW
  itti_pcef_create_session_response_t sample_pcef_csr_resp;
  fill_pcef_create_session_response(
      &sample_pcef_csr_resp, PCEF_STATUS_OK, ue_sgw_teid, DEFAULT_EPS_BEARER_ID,
      SGI_STATUS_OK);

  // check if MME gets a create session response
  EXPECT_CALL(*mme_app_handler, mme_app_handle_create_sess_resp()).Times(1);

  spgw_handle_pcef_create_session_response(
      spgw_state, &sample_pcef_csr_resp, test_imsi64);

  // create sample modify default bearer request
  itti_s11_modify_bearer_request_t sample_modify_bearer_req = {};
  fill_modify_bearer_request(
      &sample_modify_bearer_req, DEFAULT_MME_S11_TEID, ue_sgw_teid,
      DEFAULT_ENB_GTP_TEID, DEFAULT_BEARER_INDEX, DEFAULT_EPS_BEARER_ID);

  EXPECT_CALL(*mme_app_handler, mme_app_handle_modify_bearer_rsp()).Times(1);
  return_code =
      sgw_handle_modify_bearer_request(&sample_modify_bearer_req, test_imsi64);

  ASSERT_EQ(return_code, RETURNok);

  // verify that exactly one session exists in SPGW state
  ASSERT_TRUE(is_num_sessions_valid(test_imsi64, 1, 1));

  // verify that eNB address information exists
  ASSERT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 1));

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppProcedureTest, TestCreateSessionIPAllocFailure) {
  spgw_state_t* spgw_state = get_spgw_state(false);
  itti_s11_create_session_request_t sample_session_req_p = {};
  fill_create_session_request(
      &sample_session_req_p, test_imsi_str, DEFAULT_MME_S11_TEID,
      DEFAULT_BEARER_INDEX, sample_default_bearer_context, test_plmn);

  // trigger create session req to SPGW
  status_code_e create_session_rc = sgw_handle_s11_create_session_request(
      spgw_state, &sample_session_req_p, test_imsi64);

  ASSERT_EQ(create_session_rc, RETURNok);

  // Verify that a UE context exists in SPGW state after CSR is received
  spgw_ue_context_t* ue_context_p = spgw_get_ue_context(test_imsi64);
  ASSERT_TRUE(ue_context_p != nullptr);

  // Verify that teid is created
  ASSERT_FALSE(LIST_EMPTY(&ue_context_p->sgw_s11_teid_list));
  teid_t ue_sgw_teid =
      LIST_FIRST(&ue_context_p->sgw_s11_teid_list)->sgw_s11_teid;

  // Verify that no IP address is allocated for this UE
  s_plus_p_gw_eps_bearer_context_information_t* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p = sgw_cm_get_eps_bearer_entry(
      &spgw_eps_bearer_ctxt_info_p->sgw_eps_bearer_context_information
           .pdn_connection,
      DEFAULT_EPS_BEARER_ID);

  ASSERT_TRUE(eps_bearer_ctxt_p->paa.ipv4_address.s_addr == UNASSIGNED_UE_IP);

  // send an IP alloc response to SPGW with status as failure
  itti_ip_allocation_response_t test_ip_alloc_resp = {};
  fill_ip_allocation_response(
      &test_ip_alloc_resp, SGI_STATUS_ERROR_ALL_DYNAMIC_ADDRESSES_OCCUPIED,
      ue_sgw_teid, DEFAULT_EPS_BEARER_ID, DEFAULT_UE_IP, DEFAULT_VLAN);
  status_code_e ip_alloc_rc = sgw_handle_ip_allocation_rsp(
      spgw_state, &test_ip_alloc_resp, test_imsi64);

  // check that IP address is not allocated
  ASSERT_TRUE(eps_bearer_ctxt_p->paa.ipv4_address.s_addr == UNASSIGNED_UE_IP);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppProcedureTest, TestCreateSessionPCEFFailure) {
  spgw_state_t* spgw_state = get_spgw_state(false);
  // expect call to MME create session response
  itti_s11_create_session_request_t sample_session_req_p = {};
  fill_create_session_request(
      &sample_session_req_p, test_imsi_str, DEFAULT_MME_S11_TEID,
      DEFAULT_BEARER_INDEX, sample_default_bearer_context, test_plmn);

  // trigger create session req to SPGW
  status_code_e create_session_rc = sgw_handle_s11_create_session_request(
      spgw_state, &sample_session_req_p, test_imsi64);

  ASSERT_EQ(create_session_rc, RETURNok);

  // Verify that a UE context exists in SPGW state after CSR is received
  spgw_ue_context_t* ue_context_p = spgw_get_ue_context(test_imsi64);
  ASSERT_TRUE(ue_context_p != nullptr);

  // Verify that teid is created
  ASSERT_FALSE(LIST_EMPTY(&ue_context_p->sgw_s11_teid_list));
  teid_t ue_sgw_teid =
      LIST_FIRST(&ue_context_p->sgw_s11_teid_list)->sgw_s11_teid;

  // Verify that no IP address is allocated for this UE
  s_plus_p_gw_eps_bearer_context_information_t* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p = sgw_cm_get_eps_bearer_entry(
      &spgw_eps_bearer_ctxt_info_p->sgw_eps_bearer_context_information
           .pdn_connection,
      DEFAULT_EPS_BEARER_ID);

  ASSERT_TRUE(eps_bearer_ctxt_p->paa.ipv4_address.s_addr == UNASSIGNED_UE_IP);

  // send an IP alloc response to SPGW
  itti_ip_allocation_response_t test_ip_alloc_resp = {};
  fill_ip_allocation_response(
      &test_ip_alloc_resp, SGI_STATUS_OK, ue_sgw_teid, DEFAULT_EPS_BEARER_ID,
      DEFAULT_UE_IP, DEFAULT_VLAN);
  status_code_e ip_alloc_rc = sgw_handle_ip_allocation_rsp(
      spgw_state, &test_ip_alloc_resp, test_imsi64);

  // check if IP address is allocated after this message is done
  ASSERT_TRUE(eps_bearer_ctxt_p->paa.ipv4_address.s_addr == DEFAULT_UE_IP);

  // send pcef create session response to SPGW
  itti_pcef_create_session_response_t sample_pcef_csr_resp;
  fill_pcef_create_session_response(
      &sample_pcef_csr_resp, PCEF_STATUS_FAILED, ue_sgw_teid,
      DEFAULT_EPS_BEARER_ID, SGI_STATUS_OK);

  // check if MME gets a create session response
  EXPECT_CALL(*mme_app_handler, mme_app_handle_create_sess_resp()).Times(1);

  spgw_handle_pcef_create_session_response(
      spgw_state, &sample_pcef_csr_resp, test_imsi64);

  // verify that spgw context for IMSI has been cleared
  ue_context_p = spgw_get_ue_context(test_imsi64);
  ASSERT_TRUE(ue_context_p == nullptr);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppProcedureTest, TestModifyBearerFailure) {
  spgw_state_t* spgw_state  = get_spgw_state(false);
  status_code_e return_code = RETURNerror;

  // create sample modify default bearer request
  itti_s11_modify_bearer_request_t sample_modify_bearer_req = {};
  fill_modify_bearer_request(
      &sample_modify_bearer_req, DEFAULT_MME_S11_TEID, ERROR_SGW_S11_TEID,
      DEFAULT_ENB_GTP_TEID, DEFAULT_BEARER_INDEX, DEFAULT_EPS_BEARER_ID);

  EXPECT_CALL(*mme_app_handler, mme_app_handle_modify_bearer_rsp()).Times(1);
  return_code =
      sgw_handle_modify_bearer_request(&sample_modify_bearer_req, test_imsi64);

  ASSERT_EQ(return_code, RETURNok);

  // verify that no session exists in SPGW state
  ASSERT_TRUE(is_num_sessions_valid(test_imsi64, 0, 0));

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppProcedureTest, TestDeleteSessionSuccess) {
  spgw_state_t* spgw_state  = get_spgw_state(false);
  status_code_e return_code = RETURNerror;
  // expect call to MME create session response
  itti_s11_create_session_request_t sample_session_req_p = {};
  fill_create_session_request(
      &sample_session_req_p, test_imsi_str, DEFAULT_MME_S11_TEID,
      DEFAULT_BEARER_INDEX, sample_default_bearer_context, test_plmn);

  // trigger create session req to SPGW
  return_code = sgw_handle_s11_create_session_request(
      spgw_state, &sample_session_req_p, test_imsi64);

  ASSERT_EQ(return_code, RETURNok);

  // Verify that a UE context exists in SPGW state after CSR is received
  spgw_ue_context_t* ue_context_p = spgw_get_ue_context(test_imsi64);
  ASSERT_TRUE(ue_context_p != nullptr);

  // Verify that teid is created
  ASSERT_FALSE(LIST_EMPTY(&ue_context_p->sgw_s11_teid_list));
  teid_t ue_sgw_teid =
      LIST_FIRST(&ue_context_p->sgw_s11_teid_list)->sgw_s11_teid;

  // Verify that no IP address is allocated for this UE
  s_plus_p_gw_eps_bearer_context_information_t* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p = sgw_cm_get_eps_bearer_entry(
      &spgw_eps_bearer_ctxt_info_p->sgw_eps_bearer_context_information
           .pdn_connection,
      DEFAULT_EPS_BEARER_ID);

  ASSERT_TRUE(eps_bearer_ctxt_p->paa.ipv4_address.s_addr == UNASSIGNED_UE_IP);

  // send an IP alloc response to SPGW
  itti_ip_allocation_response_t test_ip_alloc_resp = {};
  fill_ip_allocation_response(
      &test_ip_alloc_resp, SGI_STATUS_OK, ue_sgw_teid, DEFAULT_EPS_BEARER_ID,
      DEFAULT_UE_IP, DEFAULT_VLAN);
  return_code = sgw_handle_ip_allocation_rsp(
      spgw_state, &test_ip_alloc_resp, test_imsi64);

  ASSERT_EQ(return_code, RETURNok);

  // check if IP address is allocated after this message is done
  ASSERT_TRUE(eps_bearer_ctxt_p->paa.ipv4_address.s_addr == DEFAULT_UE_IP);

  // send pcef create session response to SPGW
  itti_pcef_create_session_response_t sample_pcef_csr_resp;
  fill_pcef_create_session_response(
      &sample_pcef_csr_resp, PCEF_STATUS_OK, ue_sgw_teid, DEFAULT_EPS_BEARER_ID,
      SGI_STATUS_OK);

  // check if MME gets a create session response
  EXPECT_CALL(*mme_app_handler, mme_app_handle_create_sess_resp()).Times(1);

  spgw_handle_pcef_create_session_response(
      spgw_state, &sample_pcef_csr_resp, test_imsi64);

  // create sample modify default bearer request
  itti_s11_modify_bearer_request_t sample_modify_bearer_req = {};
  fill_modify_bearer_request(
      &sample_modify_bearer_req, DEFAULT_MME_S11_TEID, ue_sgw_teid,
      DEFAULT_ENB_GTP_TEID, DEFAULT_BEARER_INDEX, DEFAULT_EPS_BEARER_ID);

  EXPECT_CALL(*mme_app_handler, mme_app_handle_modify_bearer_rsp()).Times(1);
  return_code =
      sgw_handle_modify_bearer_request(&sample_modify_bearer_req, test_imsi64);

  ASSERT_EQ(return_code, RETURNok);

  // verify that exactly one session exists in SPGW state
  ASSERT_TRUE(is_num_sessions_valid(test_imsi64, 1, 1));

  // verify that eNB address information exists
  ASSERT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 1));

  // create sample delete session request
  itti_s11_delete_session_request_t sample_delete_session_request = {};
  fill_delete_session_request(
      &sample_delete_session_request, DEFAULT_MME_S11_TEID, ue_sgw_teid,
      DEFAULT_EPS_BEARER_ID, test_plmn);

  EXPECT_CALL(*mme_app_handler, mme_app_handle_delete_sess_rsp()).Times(1);

  return_code = sgw_handle_delete_session_request(
      &sample_delete_session_request, test_imsi64);
  ASSERT_EQ(return_code, RETURNok);

  // verify SPGW state is cleared
  ASSERT_TRUE(is_num_sessions_valid(test_imsi64, 0, 0));
  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppProcedureTest, TestReleaseBearerSuccess) {
  spgw_state_t* spgw_state  = get_spgw_state(false);
  status_code_e return_code = RETURNerror;
  // expect call to MME create session response
  itti_s11_create_session_request_t sample_session_req_p = {};
  fill_create_session_request(
      &sample_session_req_p, test_imsi_str, DEFAULT_MME_S11_TEID,
      DEFAULT_BEARER_INDEX, sample_default_bearer_context, test_plmn);

  // trigger create session req to SPGW
  return_code = sgw_handle_s11_create_session_request(
      spgw_state, &sample_session_req_p, test_imsi64);

  ASSERT_EQ(return_code, RETURNok);

  // Verify that a UE context exists in SPGW state after CSR is received
  spgw_ue_context_t* ue_context_p = spgw_get_ue_context(test_imsi64);
  ASSERT_TRUE(ue_context_p != nullptr);

  // Verify that teid is created
  ASSERT_FALSE(LIST_EMPTY(&ue_context_p->sgw_s11_teid_list));
  teid_t ue_sgw_teid =
      LIST_FIRST(&ue_context_p->sgw_s11_teid_list)->sgw_s11_teid;

  // Verify that no IP address is allocated for this UE
  s_plus_p_gw_eps_bearer_context_information_t* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p = sgw_cm_get_eps_bearer_entry(
      &spgw_eps_bearer_ctxt_info_p->sgw_eps_bearer_context_information
           .pdn_connection,
      DEFAULT_EPS_BEARER_ID);

  ASSERT_TRUE(eps_bearer_ctxt_p->paa.ipv4_address.s_addr == UNASSIGNED_UE_IP);

  // send an IP alloc response to SPGW
  itti_ip_allocation_response_t test_ip_alloc_resp = {};
  fill_ip_allocation_response(
      &test_ip_alloc_resp, SGI_STATUS_OK, ue_sgw_teid, DEFAULT_EPS_BEARER_ID,
      DEFAULT_UE_IP, DEFAULT_VLAN);
  return_code = sgw_handle_ip_allocation_rsp(
      spgw_state, &test_ip_alloc_resp, test_imsi64);

  ASSERT_EQ(return_code, RETURNok);

  // check if IP address is allocated after this message is done
  ASSERT_TRUE(eps_bearer_ctxt_p->paa.ipv4_address.s_addr == DEFAULT_UE_IP);

  // send pcef create session response to SPGW
  itti_pcef_create_session_response_t sample_pcef_csr_resp;
  fill_pcef_create_session_response(
      &sample_pcef_csr_resp, PCEF_STATUS_OK, ue_sgw_teid, DEFAULT_EPS_BEARER_ID,
      SGI_STATUS_OK);

  // check if MME gets a create session response
  EXPECT_CALL(*mme_app_handler, mme_app_handle_create_sess_resp()).Times(1);

  spgw_handle_pcef_create_session_response(
      spgw_state, &sample_pcef_csr_resp, test_imsi64);

  // create sample modify default bearer request
  itti_s11_modify_bearer_request_t sample_modify_bearer_req = {};
  fill_modify_bearer_request(
      &sample_modify_bearer_req, DEFAULT_MME_S11_TEID, ue_sgw_teid,
      DEFAULT_ENB_GTP_TEID, DEFAULT_BEARER_INDEX, DEFAULT_EPS_BEARER_ID);

  EXPECT_CALL(*mme_app_handler, mme_app_handle_modify_bearer_rsp()).Times(1);
  return_code =
      sgw_handle_modify_bearer_request(&sample_modify_bearer_req, test_imsi64);

  ASSERT_EQ(return_code, RETURNok);

  // verify that exactly one session exists in SPGW state
  ASSERT_TRUE(is_num_sessions_valid(test_imsi64, 1, 1));

  // verify that eNB address information exists
  ASSERT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 1));

  // send release access bearer request
  itti_s11_release_access_bearers_request_t sample_release_bearer_req = {};
  fill_release_access_bearer_request(
      &sample_release_bearer_req, DEFAULT_MME_S11_TEID, ue_sgw_teid);

  EXPECT_CALL(*mme_app_handler, mme_app_handle_release_access_bearers_resp())
      .Times(1);

  sgw_handle_release_access_bearers_request(
      &sample_release_bearer_req, test_imsi64);

  // verify that eNB information has been cleared
  ASSERT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 0));

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppProcedureTest, TestReleaseBearerError) {
  spgw_state_t* spgw_state  = get_spgw_state(false);
  status_code_e return_code = RETURNerror;
  // expect call to MME create session response
  itti_s11_create_session_request_t sample_session_req_p = {};
  fill_create_session_request(
      &sample_session_req_p, test_imsi_str, DEFAULT_MME_S11_TEID,
      DEFAULT_BEARER_INDEX, sample_default_bearer_context, test_plmn);

  // trigger create session req to SPGW
  return_code = sgw_handle_s11_create_session_request(
      spgw_state, &sample_session_req_p, test_imsi64);

  ASSERT_EQ(return_code, RETURNok);

  // Verify that a UE context exists in SPGW state after CSR is received
  spgw_ue_context_t* ue_context_p = spgw_get_ue_context(test_imsi64);
  ASSERT_TRUE(ue_context_p != nullptr);

  // Verify that teid is created
  ASSERT_FALSE(LIST_EMPTY(&ue_context_p->sgw_s11_teid_list));
  teid_t ue_sgw_teid =
      LIST_FIRST(&ue_context_p->sgw_s11_teid_list)->sgw_s11_teid;

  // Verify that no IP address is allocated for this UE
  s_plus_p_gw_eps_bearer_context_information_t* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p = sgw_cm_get_eps_bearer_entry(
      &spgw_eps_bearer_ctxt_info_p->sgw_eps_bearer_context_information
           .pdn_connection,
      DEFAULT_EPS_BEARER_ID);

  ASSERT_TRUE(eps_bearer_ctxt_p->paa.ipv4_address.s_addr == UNASSIGNED_UE_IP);

  // send an IP alloc response to SPGW
  itti_ip_allocation_response_t test_ip_alloc_resp = {};
  fill_ip_allocation_response(
      &test_ip_alloc_resp, SGI_STATUS_OK, ue_sgw_teid, DEFAULT_EPS_BEARER_ID,
      DEFAULT_UE_IP, DEFAULT_VLAN);
  return_code = sgw_handle_ip_allocation_rsp(
      spgw_state, &test_ip_alloc_resp, test_imsi64);

  ASSERT_EQ(return_code, RETURNok);

  // check if IP address is allocated after this message is done
  ASSERT_TRUE(eps_bearer_ctxt_p->paa.ipv4_address.s_addr == DEFAULT_UE_IP);

  // send pcef create session response to SPGW
  itti_pcef_create_session_response_t sample_pcef_csr_resp;
  fill_pcef_create_session_response(
      &sample_pcef_csr_resp, PCEF_STATUS_OK, ue_sgw_teid, DEFAULT_EPS_BEARER_ID,
      SGI_STATUS_OK);

  // check if MME gets a create session response
  EXPECT_CALL(*mme_app_handler, mme_app_handle_create_sess_resp()).Times(1);

  spgw_handle_pcef_create_session_response(
      spgw_state, &sample_pcef_csr_resp, test_imsi64);

  // send modify default bearer request
  itti_s11_modify_bearer_request_t sample_modify_bearer_req = {};
  fill_modify_bearer_request(
      &sample_modify_bearer_req, DEFAULT_MME_S11_TEID, ue_sgw_teid,
      DEFAULT_ENB_GTP_TEID, DEFAULT_BEARER_INDEX, DEFAULT_EPS_BEARER_ID);

  EXPECT_CALL(*mme_app_handler, mme_app_handle_modify_bearer_rsp()).Times(1);
  return_code =
      sgw_handle_modify_bearer_request(&sample_modify_bearer_req, test_imsi64);

  ASSERT_EQ(return_code, RETURNok);

  // verify that exactly one session exists in SPGW state
  ASSERT_TRUE(is_num_sessions_valid(test_imsi64, 1, 1));

  // verify that eNB address information exists
  ASSERT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 1));

  // send release access bearer request
  itti_s11_release_access_bearers_request_t sample_release_bearer_req = {};
  fill_release_access_bearer_request(
      &sample_release_bearer_req, DEFAULT_MME_S11_TEID, ERROR_SGW_S11_TEID);

  EXPECT_CALL(*mme_app_handler, mme_app_handle_release_access_bearers_resp())
      .Times(1);

  sgw_handle_release_access_bearers_request(
      &sample_release_bearer_req, test_imsi64);

  // verify that eNB information has not been cleared
  ASSERT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 1));

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppProcedureTest, TestDedicatedBearerActivation) {
  spgw_state_t* spgw_state  = get_spgw_state(false);
  status_code_e return_code = RETURNerror;
  // expect call to MME create session response
  itti_s11_create_session_request_t sample_session_req_p = {};
  fill_create_session_request(
      &sample_session_req_p, test_imsi_str, DEFAULT_MME_S11_TEID,
      DEFAULT_BEARER_INDEX, sample_default_bearer_context, test_plmn);

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
  s_plus_p_gw_eps_bearer_context_information_t* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p = sgw_cm_get_eps_bearer_entry(
      &spgw_eps_bearer_ctxt_info_p->sgw_eps_bearer_context_information
           .pdn_connection,
      DEFAULT_EPS_BEARER_ID);

  EXPECT_TRUE(eps_bearer_ctxt_p->paa.ipv4_address.s_addr == UNASSIGNED_UE_IP);

  // send an IP alloc response to SPGW
  itti_ip_allocation_response_t test_ip_alloc_resp = {};
  fill_ip_allocation_response(
      &test_ip_alloc_resp, SGI_STATUS_OK, ue_sgw_teid, DEFAULT_EPS_BEARER_ID,
      DEFAULT_UE_IP, DEFAULT_VLAN);
  return_code = sgw_handle_ip_allocation_rsp(
      spgw_state, &test_ip_alloc_resp, test_imsi64);

  EXPECT_EQ(return_code, RETURNok);

  // check if IP address is allocated after this message is done
  EXPECT_TRUE(eps_bearer_ctxt_p->paa.ipv4_address.s_addr == DEFAULT_UE_IP);

  // send pcef create session response to SPGW
  itti_pcef_create_session_response_t sample_pcef_csr_resp;
  fill_pcef_create_session_response(
      &sample_pcef_csr_resp, PCEF_STATUS_OK, ue_sgw_teid, DEFAULT_EPS_BEARER_ID,
      SGI_STATUS_OK);

  // check if MME gets a create session response
  EXPECT_CALL(*mme_app_handler, mme_app_handle_create_sess_resp()).Times(1);

  spgw_handle_pcef_create_session_response(
      spgw_state, &sample_pcef_csr_resp, test_imsi64);

  // create sample modify default bearer request
  itti_s11_modify_bearer_request_t sample_modify_bearer_req = {};
  fill_modify_bearer_request(
      &sample_modify_bearer_req, DEFAULT_MME_S11_TEID, ue_sgw_teid,
      DEFAULT_ENB_GTP_TEID, DEFAULT_BEARER_INDEX, DEFAULT_EPS_BEARER_ID);

  EXPECT_CALL(*mme_app_handler, mme_app_handle_modify_bearer_rsp()).Times(1);
  return_code =
      sgw_handle_modify_bearer_request(&sample_modify_bearer_req, test_imsi64);

  EXPECT_EQ(return_code, RETURNok);

  // verify that exactly one session exists in SPGW state
  EXPECT_TRUE(is_num_sessions_valid(test_imsi64, 1, 1));

  // send network initiated dedicated bearer activation request from Session
  // Manager
  itti_gx_nw_init_actv_bearer_request_t sample_gx_nw_init_ded_bearer_actv_req =
      {};
  gtpv2c_cause_value_t failed_cause = REQUEST_ACCEPTED;
  fill_nw_initiated_activate_bearer_request(
      &sample_gx_nw_init_ded_bearer_actv_req, test_imsi_str,
      DEFAULT_EPS_BEARER_ID, sample_dedicated_bearer_qos);

  // check that MME gets a bearer activation request
  EXPECT_CALL(
      *mme_app_handler, mme_app_handle_nw_init_ded_bearer_actv_req(
                            check_params_in_actv_bearer_req(
                                sample_gx_nw_init_ded_bearer_actv_req.lbi,
                                sample_gx_nw_init_ded_bearer_actv_req.ul_tft)))
      .Times(1);

  return_code = spgw_handle_nw_initiated_bearer_actv_req(
      spgw_state, &sample_gx_nw_init_ded_bearer_actv_req, test_imsi64,
      &failed_cause);

  EXPECT_EQ(return_code, RETURNok);

  // check number of pending procedures
  EXPECT_EQ(
      get_num_pending_create_bearer_procedures(
          &spgw_eps_bearer_ctxt_info_p->sgw_eps_bearer_context_information),
      1);

  // fetch new SGW teid for the pending bearer procedure
  pgw_ni_cbr_proc_t* pgw_ni_cbr_proc = pgw_get_procedure_create_bearer(
      &spgw_eps_bearer_ctxt_info_p->sgw_eps_bearer_context_information);
  EXPECT_TRUE(pgw_ni_cbr_proc != nullptr);
  sgw_eps_bearer_entry_wrapper_t* spgw_eps_bearer_entry_p =
      LIST_FIRST(pgw_ni_cbr_proc->pending_eps_bearers);
  teid_t ue_ded_bearer_sgw_teid =
      spgw_eps_bearer_entry_p->sgw_eps_bearer_entry->s_gw_teid_S1u_S12_S4_up;

  // send bearer activation response from MME
  itti_s11_nw_init_actv_bearer_rsp_t sample_nw_init_ded_bearer_actv_resp = {};
  fill_nw_initiated_activate_bearer_response(
      &sample_nw_init_ded_bearer_actv_resp, DEFAULT_MME_S11_TEID, ue_sgw_teid,
      ue_ded_bearer_sgw_teid, DEFAULT_ENB_GTP_TEID + 1,
      DEFAULT_EPS_BEARER_ID + 1, REQUEST_ACCEPTED, test_plmn);
  return_code = sgw_handle_nw_initiated_actv_bearer_rsp(
      &sample_nw_init_ded_bearer_actv_resp, test_imsi64);

  EXPECT_EQ(return_code, RETURNok);

  // check that bearer is created
  EXPECT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 2));

  // check that no pending procedure is left
  EXPECT_EQ(
      get_num_pending_create_bearer_procedures(
          &spgw_eps_bearer_ctxt_info_p->sgw_eps_bearer_context_information),
      0);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppProcedureTest, TestDedicatedBearerDeactivation) {
  spgw_state_t* spgw_state  = get_spgw_state(false);
  status_code_e return_code = RETURNerror;
  // expect call to MME create session response
  itti_s11_create_session_request_t sample_session_req_p = {};
  fill_create_session_request(
      &sample_session_req_p, test_imsi_str, DEFAULT_MME_S11_TEID,
      DEFAULT_BEARER_INDEX, sample_default_bearer_context, test_plmn);

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
  s_plus_p_gw_eps_bearer_context_information_t* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p = sgw_cm_get_eps_bearer_entry(
      &spgw_eps_bearer_ctxt_info_p->sgw_eps_bearer_context_information
           .pdn_connection,
      DEFAULT_EPS_BEARER_ID);

  EXPECT_TRUE(eps_bearer_ctxt_p->paa.ipv4_address.s_addr == UNASSIGNED_UE_IP);

  // send an IP alloc response to SPGW
  itti_ip_allocation_response_t test_ip_alloc_resp = {};
  fill_ip_allocation_response(
      &test_ip_alloc_resp, SGI_STATUS_OK, ue_sgw_teid, DEFAULT_EPS_BEARER_ID,
      DEFAULT_UE_IP, DEFAULT_VLAN);
  return_code = sgw_handle_ip_allocation_rsp(
      spgw_state, &test_ip_alloc_resp, test_imsi64);

  EXPECT_EQ(return_code, RETURNok);

  // check if IP address is allocated after this message is done
  EXPECT_TRUE(eps_bearer_ctxt_p->paa.ipv4_address.s_addr == DEFAULT_UE_IP);

  // send pcef create session response to SPGW
  itti_pcef_create_session_response_t sample_pcef_csr_resp;
  fill_pcef_create_session_response(
      &sample_pcef_csr_resp, PCEF_STATUS_OK, ue_sgw_teid, DEFAULT_EPS_BEARER_ID,
      SGI_STATUS_OK);

  // check if MME gets a create session response
  EXPECT_CALL(*mme_app_handler, mme_app_handle_create_sess_resp()).Times(1);

  spgw_handle_pcef_create_session_response(
      spgw_state, &sample_pcef_csr_resp, test_imsi64);

  // create sample modify default bearer request
  itti_s11_modify_bearer_request_t sample_modify_bearer_req = {};
  fill_modify_bearer_request(
      &sample_modify_bearer_req, DEFAULT_MME_S11_TEID, ue_sgw_teid,
      DEFAULT_ENB_GTP_TEID, DEFAULT_BEARER_INDEX, DEFAULT_EPS_BEARER_ID);

  EXPECT_CALL(*mme_app_handler, mme_app_handle_modify_bearer_rsp()).Times(1);
  return_code =
      sgw_handle_modify_bearer_request(&sample_modify_bearer_req, test_imsi64);

  EXPECT_EQ(return_code, RETURNok);

  // verify that exactly one session exists in SPGW state
  EXPECT_TRUE(is_num_sessions_valid(test_imsi64, 1, 1));

  // send network initiated dedicated bearer activation request from Session
  // Manager
  itti_gx_nw_init_actv_bearer_request_t sample_gx_nw_init_ded_bearer_actv_req =
      {};
  gtpv2c_cause_value_t failed_cause = REQUEST_ACCEPTED;
  fill_nw_initiated_activate_bearer_request(
      &sample_gx_nw_init_ded_bearer_actv_req, test_imsi_str,
      DEFAULT_EPS_BEARER_ID, sample_dedicated_bearer_qos);

  // check that MME gets a bearer activation request
  EXPECT_CALL(
      *mme_app_handler, mme_app_handle_nw_init_ded_bearer_actv_req(
                            check_params_in_actv_bearer_req(
                                sample_gx_nw_init_ded_bearer_actv_req.lbi,
                                sample_gx_nw_init_ded_bearer_actv_req.ul_tft)))
      .Times(1);

  return_code = spgw_handle_nw_initiated_bearer_actv_req(
      spgw_state, &sample_gx_nw_init_ded_bearer_actv_req, test_imsi64,
      &failed_cause);

  EXPECT_EQ(return_code, RETURNok);

  // check number of pending procedures
  EXPECT_EQ(
      get_num_pending_create_bearer_procedures(
          &spgw_eps_bearer_ctxt_info_p->sgw_eps_bearer_context_information),
      1);

  // fetch new SGW teid for the pending bearer procedure
  pgw_ni_cbr_proc_t* pgw_ni_cbr_proc = pgw_get_procedure_create_bearer(
      &spgw_eps_bearer_ctxt_info_p->sgw_eps_bearer_context_information);
  EXPECT_TRUE(pgw_ni_cbr_proc != nullptr);
  sgw_eps_bearer_entry_wrapper_t* spgw_eps_bearer_entry_p =
      LIST_FIRST(pgw_ni_cbr_proc->pending_eps_bearers);
  teid_t ue_ded_bearer_sgw_teid =
      spgw_eps_bearer_entry_p->sgw_eps_bearer_entry->s_gw_teid_S1u_S12_S4_up;

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
  EXPECT_EQ(
      get_num_pending_create_bearer_procedures(
          &spgw_eps_bearer_ctxt_info_p->sgw_eps_bearer_context_information),
      0);

  // send deactivate request for dedicated bearer from Session Manager
  itti_gx_nw_init_deactv_bearer_request_t
      sample_gx_nw_init_ded_bearer_deactv_req = {};
  fill_nw_initiated_deactivate_bearer_request(
      &sample_gx_nw_init_ded_bearer_deactv_req, test_imsi_str,
      DEFAULT_EPS_BEARER_ID, ded_eps_bearer_id);

  // check that MME gets a bearer deactivation request
  EXPECT_CALL(
      *mme_app_handler,
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
  int num_bearers_to_delete   = 1;
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

TEST_F(
    SPGWAppProcedureTest, TestDedicatedBearerDeactivationDeleteDefaultBearer) {
  spgw_state_t* spgw_state  = get_spgw_state(false);
  status_code_e return_code = RETURNerror;
  // expect call to MME create session response
  itti_s11_create_session_request_t sample_session_req_p = {};
  fill_create_session_request(
      &sample_session_req_p, test_imsi_str, DEFAULT_MME_S11_TEID,
      DEFAULT_BEARER_INDEX, sample_default_bearer_context, test_plmn);

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
  s_plus_p_gw_eps_bearer_context_information_t* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p = sgw_cm_get_eps_bearer_entry(
      &spgw_eps_bearer_ctxt_info_p->sgw_eps_bearer_context_information
           .pdn_connection,
      DEFAULT_EPS_BEARER_ID);

  EXPECT_TRUE(eps_bearer_ctxt_p->paa.ipv4_address.s_addr == UNASSIGNED_UE_IP);

  // send an IP alloc response to SPGW
  itti_ip_allocation_response_t test_ip_alloc_resp = {};
  fill_ip_allocation_response(
      &test_ip_alloc_resp, SGI_STATUS_OK, ue_sgw_teid, DEFAULT_EPS_BEARER_ID,
      DEFAULT_UE_IP, DEFAULT_VLAN);
  return_code = sgw_handle_ip_allocation_rsp(
      spgw_state, &test_ip_alloc_resp, test_imsi64);

  EXPECT_EQ(return_code, RETURNok);

  // check if IP address is allocated after this message is done
  EXPECT_TRUE(eps_bearer_ctxt_p->paa.ipv4_address.s_addr == DEFAULT_UE_IP);

  // send pcef create session response to SPGW
  itti_pcef_create_session_response_t sample_pcef_csr_resp;
  fill_pcef_create_session_response(
      &sample_pcef_csr_resp, PCEF_STATUS_OK, ue_sgw_teid, DEFAULT_EPS_BEARER_ID,
      SGI_STATUS_OK);

  // check if MME gets a create session response
  EXPECT_CALL(*mme_app_handler, mme_app_handle_create_sess_resp()).Times(1);

  spgw_handle_pcef_create_session_response(
      spgw_state, &sample_pcef_csr_resp, test_imsi64);

  // create sample modify default bearer request
  itti_s11_modify_bearer_request_t sample_modify_bearer_req = {};
  fill_modify_bearer_request(
      &sample_modify_bearer_req, DEFAULT_MME_S11_TEID, ue_sgw_teid,
      DEFAULT_ENB_GTP_TEID, DEFAULT_BEARER_INDEX, DEFAULT_EPS_BEARER_ID);

  EXPECT_CALL(*mme_app_handler, mme_app_handle_modify_bearer_rsp()).Times(1);
  return_code =
      sgw_handle_modify_bearer_request(&sample_modify_bearer_req, test_imsi64);

  EXPECT_EQ(return_code, RETURNok);

  // verify that exactly one session exists in SPGW state
  EXPECT_TRUE(is_num_sessions_valid(test_imsi64, 1, 1));

  // send network initiated dedicated bearer activation request from Session
  // Manager
  itti_gx_nw_init_actv_bearer_request_t sample_gx_nw_init_ded_bearer_actv_req =
      {};
  gtpv2c_cause_value_t failed_cause = REQUEST_ACCEPTED;
  fill_nw_initiated_activate_bearer_request(
      &sample_gx_nw_init_ded_bearer_actv_req, test_imsi_str,
      DEFAULT_EPS_BEARER_ID, sample_dedicated_bearer_qos);

  // check that MME gets a bearer activation request
  EXPECT_CALL(
      *mme_app_handler, mme_app_handle_nw_init_ded_bearer_actv_req(
                            check_params_in_actv_bearer_req(
                                sample_gx_nw_init_ded_bearer_actv_req.lbi,
                                sample_gx_nw_init_ded_bearer_actv_req.ul_tft)))
      .Times(1);

  return_code = spgw_handle_nw_initiated_bearer_actv_req(
      spgw_state, &sample_gx_nw_init_ded_bearer_actv_req, test_imsi64,
      &failed_cause);

  EXPECT_EQ(return_code, RETURNok);

  // check number of pending procedures
  EXPECT_EQ(
      get_num_pending_create_bearer_procedures(
          &spgw_eps_bearer_ctxt_info_p->sgw_eps_bearer_context_information),
      1);

  // fetch new SGW teid for the pending bearer procedure
  pgw_ni_cbr_proc_t* pgw_ni_cbr_proc = pgw_get_procedure_create_bearer(
      &spgw_eps_bearer_ctxt_info_p->sgw_eps_bearer_context_information);
  EXPECT_TRUE(pgw_ni_cbr_proc != nullptr);
  sgw_eps_bearer_entry_wrapper_t* spgw_eps_bearer_entry_p =
      LIST_FIRST(pgw_ni_cbr_proc->pending_eps_bearers);
  teid_t ue_ded_bearer_sgw_teid =
      spgw_eps_bearer_entry_p->sgw_eps_bearer_entry->s_gw_teid_S1u_S12_S4_up;

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
  EXPECT_EQ(
      get_num_pending_create_bearer_procedures(
          &spgw_eps_bearer_ctxt_info_p->sgw_eps_bearer_context_information),
      0);

  // send deactivate request for dedicated bearer from Session Manager
  itti_gx_nw_init_deactv_bearer_request_t
      sample_gx_nw_init_ded_bearer_deactv_req = {};
  fill_nw_initiated_deactivate_bearer_request(
      &sample_gx_nw_init_ded_bearer_deactv_req, test_imsi_str,
      DEFAULT_EPS_BEARER_ID, ded_eps_bearer_id);

  // check that MME gets a bearer deactivation request
  EXPECT_CALL(
      *mme_app_handler,
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
  int num_bearers_to_delete   = 2;
  ebi_t eps_bearer_id_array[] = {DEFAULT_EPS_BEARER_ID, ded_eps_bearer_id};

  fill_nw_initiated_deactivate_bearer_response(
      &sample_nw_init_ded_bearer_deactv_resp, test_imsi64, true,
      REQUEST_ACCEPTED, eps_bearer_id_array, num_bearers_to_delete,
      ue_sgw_teid);
  return_code = sgw_handle_nw_initiated_deactv_bearer_rsp(
      spgw_state, &sample_nw_init_ded_bearer_deactv_resp, test_imsi64);
  EXPECT_EQ(return_code, RETURNok);

  // check that session is removed
  EXPECT_TRUE(is_num_sessions_valid(test_imsi64, 0, 0));

  free(sample_nw_init_ded_bearer_deactv_resp.lbi);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppProcedureTest, TestSuspendNotification) {
  spgw_state_t* spgw_state  = get_spgw_state(false);
  status_code_e return_code = RETURNerror;

  // expect call to MME create session response
  itti_s11_create_session_request_t sample_session_req_p = {};
  fill_create_session_request(
      &sample_session_req_p, test_imsi_str, DEFAULT_MME_S11_TEID,
      DEFAULT_BEARER_INDEX, sample_default_bearer_context, test_plmn);

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
  s_plus_p_gw_eps_bearer_context_information_t* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p = sgw_cm_get_eps_bearer_entry(
      &spgw_eps_bearer_ctxt_info_p->sgw_eps_bearer_context_information
           .pdn_connection,
      DEFAULT_EPS_BEARER_ID);

  EXPECT_TRUE(eps_bearer_ctxt_p->paa.ipv4_address.s_addr == UNASSIGNED_UE_IP);

  // send an IP alloc response to SPGW
  itti_ip_allocation_response_t test_ip_alloc_resp = {};
  fill_ip_allocation_response(
      &test_ip_alloc_resp, SGI_STATUS_OK, ue_sgw_teid, DEFAULT_EPS_BEARER_ID,
      DEFAULT_UE_IP, DEFAULT_VLAN);
  return_code = sgw_handle_ip_allocation_rsp(
      spgw_state, &test_ip_alloc_resp, test_imsi64);

  EXPECT_EQ(return_code, RETURNok);

  // check if IP address is allocated after this message is done
  EXPECT_TRUE(eps_bearer_ctxt_p->paa.ipv4_address.s_addr == DEFAULT_UE_IP);

  // send pcef create session response to SPGW
  itti_pcef_create_session_response_t sample_pcef_csr_resp;
  fill_pcef_create_session_response(
      &sample_pcef_csr_resp, PCEF_STATUS_OK, ue_sgw_teid, DEFAULT_EPS_BEARER_ID,
      SGI_STATUS_OK);

  // check if MME gets a create session response
  EXPECT_CALL(*mme_app_handler, mme_app_handle_create_sess_resp()).Times(1);

  spgw_handle_pcef_create_session_response(
      spgw_state, &sample_pcef_csr_resp, test_imsi64);

  // create sample modify default bearer request
  itti_s11_modify_bearer_request_t sample_modify_bearer_req = {};
  fill_modify_bearer_request(
      &sample_modify_bearer_req, DEFAULT_MME_S11_TEID, ue_sgw_teid,
      DEFAULT_ENB_GTP_TEID, DEFAULT_BEARER_INDEX, DEFAULT_EPS_BEARER_ID);

  EXPECT_CALL(*mme_app_handler, mme_app_handle_modify_bearer_rsp()).Times(1);
  return_code =
      sgw_handle_modify_bearer_request(&sample_modify_bearer_req, test_imsi64);

  EXPECT_EQ(return_code, RETURNok);

  // verify that exactly one session exists in SPGW state
  EXPECT_TRUE(is_num_sessions_valid(test_imsi64, 1, 1));

  // verify that eNB address information exists
  EXPECT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 1));

  // trigger suspend notification to SPGW task
  itti_s11_suspend_notification_t sample_suspend_notification = {};
  fill_s11_suspend_notification(
      &sample_suspend_notification, ue_sgw_teid, test_imsi_str,
      DEFAULT_EPS_BEARER_ID);

  // verify that mock MME app task receives an acknowledgement with
  // REQUEST_ACCEPTED
  EXPECT_CALL(
      *mme_app_handler,
      mme_app_handle_suspend_acknowledge(check_params_in_suspend_ack(
          REQUEST_ACCEPTED,
          spgw_eps_bearer_ctxt_info_p->sgw_eps_bearer_context_information
              .mme_teid_S11)))
      .Times(1);
  return_code = sgw_handle_suspend_notification(
      &sample_suspend_notification, test_imsi64);

  EXPECT_EQ(return_code, RETURNok);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

}  // namespace lte
}  // namespace magma
