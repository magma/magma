#include <gtest/gtest.h>
#include <string>
#include <thread>

#include "../mock_tasks/mock_tasks.h"
#include "spgw_test_util.h"
#include "spgw_state.h"
#include "spgw_state_converter.h"
#include "mock_spgw_op.h"

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

class SPGWAppInjectedStateProcedureTest : public ::testing::Test {
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
    // might need to mock this function
    // spgw_state_init(false, &spgw_config);
    spgw_app_init(&spgw_config, mme_config.use_stateless);    
    // mock_spgw_app_init();
    mock_read_spgw_ue_state_db();
    
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
  std::string test_imsi_str          = "001010000000002";
  unsigned long long int test_imsi64 = 1010000000002;
  plmn_t test_plmn                   = {.mcc_digit2 = 0,
                      .mcc_digit1 = 0,
                      .mnc_digit3 = 0x00,
                      .mcc_digit3 = 0,
                      .mnc_digit2 = 0,
                      .mnc_digit1 = 0};
  bearer_context_to_be_created_t sample_default_bearer_context = {
      .eps_bearer_id    = 5,
      .bearer_level_qos = {.pci = 1,
                           .pl  = 15,
                          //  .pvi = 0,
                           .qci = 9,
                           .gbr = {},
                           }}; // .mbr = {} .br_ul = 200000000, .br_dl = 100000000
  int test_mme_s11_teid = 14;
  int test_bearer_index = 5;
  in_addr_t test_ue_ip = 0x01020305;

};

TEST_F(SPGWAppInjectedStateProcedureTest, TestIPAllocFailure) {
  spgw_state_t* spgw_state = get_spgw_state(false);

  // Verify that a UE context exists in SPGW state after CSR is received
  spgw_ue_context_t* ue_context_p = spgw_get_ue_context(test_imsi64);

  teid_t ue_sgw_teid =
      LIST_FIRST(&ue_context_p->sgw_s11_teid_list)->sgw_s11_teid;

  // Verify that no IP address is allocated for this UE
  s_plus_p_gw_eps_bearer_context_information_t* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p = sgw_cm_get_eps_bearer_entry(
      &spgw_eps_bearer_ctxt_info_p->sgw_eps_bearer_context_information
           .pdn_connection,
      DEFAULT_EPS_BEARER_ID);

  // send an IP alloc response to SPGW with status as failure
  itti_ip_allocation_response_t test_ip_alloc_resp = {};
  fill_ip_allocation_response(
      &test_ip_alloc_resp, SGI_STATUS_ERROR_ALL_DYNAMIC_ADDRESSES_OCCUPIED,
      ue_sgw_teid, DEFAULT_EPS_BEARER_ID, 0x0d80a8c0, DEFAULT_VLAN);
  status_code_e ip_alloc_rc = sgw_handle_ip_allocation_rsp(
      spgw_state, &test_ip_alloc_resp, test_imsi64);

  // check that IP address is not allocated
  ASSERT_TRUE(eps_bearer_ctxt_p->paa.ipv4_address.s_addr == 0x0d80a8c0);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppInjectedStateProcedureTest, TestCreateSessionPCEFFailure) {
  spgw_state_t* spgw_state = get_spgw_state(false);
  
  // Verify that a UE context exists in SPGW state after CSR is received
  spgw_ue_context_t* ue_context_p = spgw_get_ue_context(test_imsi64);

  // Verify that teid is created
  teid_t ue_sgw_teid =
      LIST_FIRST(&ue_context_p->sgw_s11_teid_list)->sgw_s11_teid;

  // Verify that no IP address is allocated for this UE
  s_plus_p_gw_eps_bearer_context_information_t* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p = sgw_cm_get_eps_bearer_entry(
      &spgw_eps_bearer_ctxt_info_p->sgw_eps_bearer_context_information
           .pdn_connection,
      DEFAULT_EPS_BEARER_ID);

  // send an IP alloc response to SPGW
  itti_ip_allocation_response_t test_ip_alloc_resp = {};
  fill_ip_allocation_response(
      &test_ip_alloc_resp, SGI_STATUS_OK, ue_sgw_teid, DEFAULT_EPS_BEARER_ID,
      0x0d80a8c0, DEFAULT_VLAN);
  status_code_e ip_alloc_rc = sgw_handle_ip_allocation_rsp(
      spgw_state, &test_ip_alloc_resp, test_imsi64);

  // check if IP address is allocated after this message is done
  ASSERT_TRUE(eps_bearer_ctxt_p->paa.ipv4_address.s_addr == 0x0d80a8c0);

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

TEST_F(SPGWAppInjectedStateProcedureTest, TestDeleteSessionSuccess) {
  spgw_state_t* spgw_state  = get_spgw_state(false);
  
  spgw_ue_context_t* ue_context_p = spgw_get_ue_context(test_imsi64);

  // create sample delete session request
  teid_t ue_sgw_teid =
      LIST_FIRST(&ue_context_p->sgw_s11_teid_list)->sgw_s11_teid;

  // check if IP address is allocated after this message is done
  s_plus_p_gw_eps_bearer_context_information_t* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p = sgw_cm_get_eps_bearer_entry(
      &spgw_eps_bearer_ctxt_info_p->sgw_eps_bearer_context_information
           .pdn_connection,
      DEFAULT_EPS_BEARER_ID);

  char str[INET_ADDRSTRLEN];
  
    // now get it back and print it
    inet_ntop(AF_INET, &(eps_bearer_ctxt_p->paa.ipv4_address.s_addr), str, INET_ADDRSTRLEN);

  std::cout << "======== IP: " << std::string(str) << std::endl;
  std::cout <<  "val : " << std::hex << eps_bearer_ctxt_p->paa.ipv4_address.s_addr << std::endl;
  ASSERT_TRUE(eps_bearer_ctxt_p->paa.ipv4_address.s_addr == 0x0d80a8c0); //c0a8800d    

  // verify that exactly one session exists in SPGW state
  ASSERT_TRUE(is_num_sessions_valid(spgw_state, test_imsi64, 1, 1));

  

  itti_s11_delete_session_request_t sample_delete_session_request = {};
  fill_delete_session_request(
      &sample_delete_session_request, test_mme_s11_teid, ue_sgw_teid,
      DEFAULT_EPS_BEARER_ID, test_plmn);

  EXPECT_CALL(*mme_app_handler, mme_app_handle_delete_sess_rsp()).Times(1);

  status_code_e return_code = RETURNerror;
  return_code = sgw_handle_delete_session_request(
      &sample_delete_session_request, test_imsi64);
  ASSERT_EQ(return_code, RETURNok);

  // verify SPGW state is cleared
  ASSERT_TRUE(is_num_sessions_valid(spgw_state, test_imsi64, 0, 0));
  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

}  // namespace lte
}  // namespace magma
