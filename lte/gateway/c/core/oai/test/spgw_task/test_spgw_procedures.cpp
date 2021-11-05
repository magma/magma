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
#include <thread>

#include "../mock_tasks/mock_tasks.h"
#include "common_defs.h"
#include "spgw_test_util.h"
#include "spgw_state.h"

extern "C" {
#include "mme_config.h"
#include "pgw_handlers.h"
#include "sgw_context_manager.h"
#include "sgw_defs.h"
#include "sgw_handlers.h"
#include "spgw_config.h"
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
  std::string test_imsi_str          = "001010000000001";
  unsigned long long int test_imsi64 = 1010000000001;
  plmn_t test_plmn                   = {.mcc_digit2 = 0,
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
};

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
      &sample_modify_bearer_req, DEFAULT_MME_S11_TEID, DEFAULT_BEARER_INDEX,
      DEFAULT_EPS_BEARER_ID, ue_sgw_teid, DEFAULT_ENB_GTP_TEID);

  EXPECT_CALL(*mme_app_handler, mme_app_handle_modify_bearer_rsp()).Times(1);
  return_code =
      sgw_handle_modify_bearer_request(&sample_modify_bearer_req, test_imsi64);

  ASSERT_EQ(return_code, RETURNok);

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

}  // namespace lte
}  // namespace magma
