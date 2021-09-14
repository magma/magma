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
#include <chrono>
#include <gtest/gtest.h>
#include <thread>

#include "feg/protos/s6a_proxy.pb.h"
#include "../mock_tasks/mock_tasks.h"
#include "mme_app_state_manager.h"
#include "mme_app_ip_imsi.h"
#include "proto_msg_to_itti_msg.h"

extern "C" {
#include "common_types.h"
#include "mme_config.h"
#include "mme_app_defs.h"
#include "mme_app_extern.h"
#include "mme_app_state.h"
#include "security_types.h"
}

using ::testing::_;
using ::testing::Return;

namespace magma {

namespace lte {

task_zmq_ctx_t task_zmq_ctx_main;

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    default: { } break; }

  itti_free_msg_content(received_message_p);
  free(received_message_p);
  return 0;
}
static void send_mme_app_uplink_data_ind(
    uint8_t* nas_msg, uint8_t nas_msg_length, const plmn_t& plmn) {
  MessageDef* message_p =
      itti_alloc_new_message(TASK_S1AP, MME_APP_UPLINK_DATA_IND);
  ITTI_MSG_LASTHOP_LATENCY(message_p)     = 0;
  MME_APP_UL_DATA_IND(message_p).ue_id    = 1;
  MME_APP_UL_DATA_IND(message_p).nas_msg  = blk2bstr(nas_msg, nas_msg_length);
  MME_APP_UL_DATA_IND(message_p).tai.plmn = plmn;
  MME_APP_UL_DATA_IND(message_p).tai.tac  = 1;
  MME_APP_UL_DATA_IND(message_p).cgi.plmn = plmn;
  MME_APP_UL_DATA_IND(message_p).cgi.cell_identity = {0, 0, 0};
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);
  return;
}

class MmeAppProcedureTest : public ::testing::Test {
  virtual void SetUp() {
    s1ap_handler = std::make_shared<MockS1apHandler>();
    s6a_handler  = std::make_shared<MockS6aHandler>();
    spgw_handler = std::make_shared<MockSpgwHandler>();

    itti_init(
        TASK_MAX, THREAD_MAX, MESSAGES_ID_MAX, tasks_info, messages_info, NULL,
        NULL);

    // initialize mme config
    mme_config_init(&mme_config);
    create_partial_lists(&mme_config);
    mme_config.use_stateless                              = false;
    mme_config.nas_config.prefered_integrity_algorithm[0] = EIA2_128_ALG_ID;

    task_id_t task_id_list[10] = {
        TASK_MME_APP,    TASK_HA,  TASK_S1AP,   TASK_S6A,      TASK_S11,
        TASK_SERVICE303, TASK_SGS, TASK_SGW_S8, TASK_SPGW_APP, TASK_SMS_ORC8R};
    init_task_context(
        TASK_MAIN, task_id_list, 10, handle_message, &task_zmq_ctx_main);

    std::thread task_ha(start_mock_ha_task);
    std::thread task_s1ap(start_mock_s1ap_task, s1ap_handler);
    std::thread task_s6a(start_mock_s6a_task, s6a_handler);
    std::thread task_s11(start_mock_s11_task);
    std::thread task_service303(start_mock_service303_task);
    std::thread task_sgs(start_mock_sgs_task);
    std::thread task_sgw_s8(start_mock_sgw_s8_task);
    std::thread task_sms_orc8r(start_mock_sms_orc8r_task);
    std::thread task_spgw(start_mock_spgw_task, spgw_handler);
    task_ha.detach();
    task_s1ap.detach();
    task_s6a.detach();
    task_s11.detach();
    task_service303.detach();
    task_sgs.detach();
    task_sgw_s8.detach();
    task_sms_orc8r.detach();
    task_spgw.detach();

    mme_app_init(&mme_config);
  }

  virtual void TearDown() {
    send_terminate_message_fatal(&task_zmq_ctx_main);
    destroy_task_context(&task_zmq_ctx_main);
    itti_free_desc_threads();
    // Sleep to ensure that messages are received and contexts are released
    std::this_thread::sleep_for(std::chrono::milliseconds(1000));
  }

 protected:
  std::shared_ptr<MockS1apHandler> s1ap_handler;
  std::shared_ptr<MockS6aHandler> s6a_handler;
  std::shared_ptr<MockSpgwHandler> spgw_handler;
};

TEST_F(MmeAppProcedureTest, TestInitialUeMessageFaultyNasMsg) {
  MessageDef* message_p = NULL;

  // Initial Attach
  message_p = itti_alloc_new_message(TASK_S1AP, S1AP_INITIAL_UE_MESSAGE);
  /* The following buffer just includes an attach request */
  uint8_t nas_msg[]       = {0x72, 0x08, 0x09, 0x10, 0x10, 0x00, 0x00, 0x00,
                       0x00, 0x10, 0x02, 0xe0, 0xe0, 0x00, 0x04, 0x02,
                       0x01, 0xd0, 0x11, 0x40, 0x08, 0x04, 0x02, 0x60,
                       0x04, 0x00, 0x02, 0x1c, 0x00};
  uint32_t nas_msg_length = 29;

  EXPECT_CALL(*s1ap_handler, s1ap_generate_downlink_nas_transport()).Times(1);

  ITTI_MSG_LASTHOP_LATENCY(message_p)               = 0;
  S1AP_INITIAL_UE_MESSAGE(message_p).sctp_assoc_id  = 0;
  S1AP_INITIAL_UE_MESSAGE(message_p).enb_ue_s1ap_id = 0;
  S1AP_INITIAL_UE_MESSAGE(message_p).enb_id         = 0;
  S1AP_INITIAL_UE_MESSAGE(message_p).nas = blk2bstr(nas_msg, nas_msg_length);
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(1000));
}

TEST_F(MmeAppProcedureTest, TestImsiAttachEpsOnlyDetach) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  MessageDef* message_p = NULL;
  std::string imsi      = "001010000000001";
  plmn_t plmn           = {.mcc_digit2 = 0,
                 .mcc_digit1 = 0,
                 .mnc_digit3 = 0x0f,
                 .mcc_digit3 = 1,
                 .mnc_digit2 = 1,
                 .mnc_digit1 = 0};

  EXPECT_CALL(*s1ap_handler, s1ap_generate_downlink_nas_transport()).Times(3);
  EXPECT_CALL(*s1ap_handler, s1ap_handle_conn_est_cnf()).Times(1);
  EXPECT_CALL(*s1ap_handler, s1ap_handle_ue_context_release_command()).Times(1);
  EXPECT_CALL(*s6a_handler, s6a_viface_authentication_info_req()).Times(1);
  EXPECT_CALL(*s6a_handler, s6a_viface_update_location_req()).Times(1);
  EXPECT_CALL(*s6a_handler, s6a_viface_purge_ue()).Times(1);
  EXPECT_CALL(*spgw_handler, sgw_handle_s11_create_session_request()).Times(1);
  EXPECT_CALL(*spgw_handler, sgw_handle_delete_session_request()).Times(1);

  // Construction and sending Initial Attach Request to mme_app mimicing S1AP
  message_p = itti_alloc_new_message(TASK_S1AP, S1AP_INITIAL_UE_MESSAGE);

  uint8_t nas_msg[]       = {0x07, 0x41, 0x71, 0x08, 0x09, 0x10, 0x10, 0x00,
                       0x00, 0x00, 0x00, 0x10, 0x02, 0xe0, 0xe0, 0x00,
                       0x04, 0x02, 0x01, 0xd0, 0x11, 0x40, 0x08, 0x04,
                       0x02, 0x60, 0x04, 0x00, 0x02, 0x1c, 0x00};
  uint32_t nas_msg_length = 31;

  ITTI_MSG_LASTHOP_LATENCY(message_p)               = 0;
  S1AP_INITIAL_UE_MESSAGE(message_p).sctp_assoc_id  = 0;
  S1AP_INITIAL_UE_MESSAGE(message_p).enb_ue_s1ap_id = 0;
  S1AP_INITIAL_UE_MESSAGE(message_p).enb_id         = 0;
  S1AP_INITIAL_UE_MESSAGE(message_p).nas = blk2bstr(nas_msg, nas_msg_length);
  S1AP_INITIAL_UE_MESSAGE(message_p).tai.plmn           = plmn;
  S1AP_INITIAL_UE_MESSAGE(message_p).tai.tac            = 1;
  S1AP_INITIAL_UE_MESSAGE(message_p).ecgi.plmn          = plmn;
  S1AP_INITIAL_UE_MESSAGE(message_p).ecgi.cell_identity = {0, 0, 0};
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);

  // Sending AIA to mme_app mimicing successful S6A response for AIR
  message_p = itti_alloc_new_message(TASK_S6A, S6A_AUTH_INFO_ANS);
  s6a_auth_info_ans_t* itti_msg = &message_p->ittiMsg.s6a_auth_info_ans;
  strncpy(itti_msg->imsi, imsi.c_str(), imsi.size());
  itti_msg->imsi_length        = imsi.size();
  itti_msg->result.present     = S6A_RESULT_BASE;
  itti_msg->result.choice.base = DIAMETER_SUCCESS;
  magma::feg::AuthenticationInformationAnswer aia;
  magma::feg::AuthenticationInformationAnswer::EUTRANVector eutran_vector;
  uint8_t xres_buf[XRES_LENGTH_MAX]      = {0x66, 0xff, 0x47, 0x2d, 0xd4, 0x93,
                                       0xf1, 0x5a, 0x00, 0x00, 0x00, 0x00,
                                       0x00, 0x00, 0x00, 0x00};
  uint8_t rand_buf[RAND_LENGTH_OCTETS]   = {0x68, 0x16, 0xa1, 0x0c, 0x0f, 0xeb,
                                          0x44, 0xa5, 0x00, 0x5c, 0x9c, 0x9c,
                                          0x3c, 0x6f, 0xd6, 0x15};
  uint8_t autn_buf[AUTN_LENGTH_OCTETS]   = {0x4a, 0xe4, 0xe0, 0xd9, 0xaa, 0x4b,
                                          0x80, 0x00, 0xc4, 0x80, 0xa1, 0x97,
                                          0x70, 0x4b, 0x7b, 0x8f};
  uint8_t kasme_buf[KASME_LENGTH_OCTETS] = {
      0xc3, 0x5f, 0x03, 0x8f, 0x5f, 0xbe, 0xcc, 0x23, 0xc4, 0xd1, 0xa7,
      0xd6, 0x8a, 0xf7, 0x05, 0x32, 0xf2, 0x37, 0xf6, 0x40, 0x47, 0xdd,
      0x29, 0x6e, 0x7d, 0x0e, 0xf6, 0xe9, 0x26, 0x5f, 0x24, 0x39};
  eutran_vector.set_rand((const void*) rand_buf, RAND_LENGTH_OCTETS);
  eutran_vector.set_xres((const void*) xres_buf, XRES_LENGTH_MAX);
  eutran_vector.set_autn((const void*) autn_buf, AUTN_LENGTH_OCTETS);
  eutran_vector.set_kasme((const void*) kasme_buf, KASME_LENGTH_OCTETS);
  aia.set_error_code(magma::feg::ErrorCode::SUCCESS);
  auto eutran_vectors = aia.mutable_eutran_vectors();
  eutran_vectors->Add()->CopyFrom(eutran_vector);
  magma::convert_proto_msg_to_itti_s6a_auth_info_ans(aia, itti_msg);
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);

  // Constructing and sending Authentication Response to mme_app mimicing S1AP
  uint8_t nas_msg2[] = {0x07, 0x53, 0x10, 0x66, 0xff, 0x47, 0x2d,
                        0xd4, 0x93, 0xf1, 0x5a, 0x00, 0x00, 0x00,
                        0x00, 0x00, 0x00, 0x00, 0x00};
  nas_msg_length     = 19;
  send_mme_app_uplink_data_ind(&nas_msg2[0], nas_msg_length, plmn);

  // Constructing and sending SMC Response to mme_app mimicing S1AP
  uint8_t nas_msg3[] = {0x47, 0xc0, 0xb5, 0x35, 0x6b, 0x00, 0x07,
                        0x5e, 0x23, 0x09, 0x33, 0x08, 0x45, 0x86,
                        0x34, 0x12, 0x31, 0x71, 0xf2};
  nas_msg_length     = 19;
  send_mme_app_uplink_data_ind(&nas_msg3[0], nas_msg_length, plmn);

  // Sending ULA to mme_app mimicing successful S6A response for ULR
  message_p = itti_alloc_new_message(TASK_S6A, S6A_UPDATE_LOCATION_ANS);
  s6a_update_location_ans_t* itti_msg2 =
      &message_p->ittiMsg.s6a_update_location_ans;
  strncpy(itti_msg2->imsi, imsi.c_str(), imsi.size());
  itti_msg2->imsi_length        = imsi.size();
  itti_msg2->result.present     = S6A_RESULT_BASE;
  itti_msg2->result.choice.base = DIAMETER_SUCCESS;
  magma::feg::UpdateLocationAnswer ula;
  ula.set_default_context_id(0);
  auto total_ambr = ula.mutable_total_ambr();
  total_ambr->set_max_bandwidth_ul(100000000);
  total_ambr->set_max_bandwidth_dl(200000000);
  ula.set_all_apns_included(false);
  magma::feg::UpdateLocationAnswer::APNConfiguration apnconfig;
  apnconfig.set_context_id(0);
  apnconfig.set_service_selection("magma.ipv4");
  auto apn_qosprofile = apnconfig.mutable_qos_profile();
  apn_qosprofile->set_class_id(9);
  apn_qosprofile->set_priority_level(15);
  auto apn_ambr = apnconfig.mutable_ambr();
  apn_ambr->set_max_bandwidth_ul(10000000);
  apn_ambr->set_max_bandwidth_dl(75000000);
  apnconfig.set_pdn(magma::feg::UpdateLocationAnswer::APNConfiguration::IPV4);
  auto apns = ula.mutable_apn();
  apns->Add()->CopyFrom(apnconfig);
  magma::convert_proto_msg_to_itti_s6a_update_location_ans(ula, itti_msg2);
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);

  // Constructing and sending Create Session Response to mme_app mimicing SPGW
  message_p =
      itti_alloc_new_message(TASK_SPGW_APP, S11_CREATE_SESSION_RESPONSE);
  itti_s11_create_session_response_t* create_session_response_p =
      &message_p->ittiMsg.s11_create_session_response;
  create_session_response_p->teid                    = 1;
  create_session_response_p->cause.cause_value       = REQUEST_ACCEPTED;
  create_session_response_p->paa.pdn_type            = IPv4;
  create_session_response_p->paa.ipv4_address.s_addr = 1000;
  create_session_response_p->bearer_contexts_created.bearer_contexts[0]
      .cause.cause_value = REQUEST_ACCEPTED;
  create_session_response_p->bearer_contexts_created.bearer_contexts[0]
      .s1u_sgw_fteid.teid = 1000;
  create_session_response_p->bearer_contexts_created.bearer_contexts[0]
      .s1u_sgw_fteid.interface_type = S1_U_SGW_GTP_U;
  create_session_response_p->bearer_contexts_created.bearer_contexts[0]
      .s1u_sgw_fteid.ipv4 = 1;
  create_session_response_p->bearer_contexts_created.bearer_contexts[0]
      .s1u_sgw_fteid.ipv4_address.s_addr = 100;
  create_session_response_p->bearer_contexts_created.bearer_contexts[0]
      .eps_bearer_id                                                    = 5;
  create_session_response_p->bearer_contexts_created.num_bearer_context = 1;
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);

  // Constructing and sending ICS Response to mme_app mimicing S1AP
  message_p =
      itti_alloc_new_message(TASK_S1AP, MME_APP_INITIAL_CONTEXT_SETUP_RSP);
  MME_APP_INITIAL_CONTEXT_SETUP_RSP(message_p).ue_id                        = 1;
  MME_APP_INITIAL_CONTEXT_SETUP_RSP(message_p).e_rab_setup_list.no_of_items = 1;
  MME_APP_INITIAL_CONTEXT_SETUP_RSP(message_p)
      .e_rab_setup_list.item[0]
      .e_rab_id = 5;
  MME_APP_INITIAL_CONTEXT_SETUP_RSP(message_p)
      .e_rab_setup_list.item[0]
      .gtp_teid                     = 0;
  uint8_t transport_address_buff[4] = {0, 0, 0, 0};
  MME_APP_INITIAL_CONTEXT_SETUP_RSP(message_p)
      .e_rab_setup_list.item[0]
      .transport_layer_address = blk2bstr(transport_address_buff, 4);
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);

  // Constructing UE Capability Indication message to mme_app
  // mimicing S1AP with dummy radio capabilities
  message_p = itti_alloc_new_message(TASK_S1AP, S1AP_UE_CAPABILITIES_IND);
  itti_s1ap_ue_cap_ind_t* ue_cap_ind_p    = &message_p->ittiMsg.s1ap_ue_cap_ind;
  ue_cap_ind_p->enb_ue_s1ap_id            = 0;
  ue_cap_ind_p->mme_ue_s1ap_id            = 1;
  ue_cap_ind_p->radio_capabilities_length = 200;
  // using malloc to create uninitialized buffer
  ue_cap_ind_p->radio_capabilities =
      (uint8_t*) malloc(ue_cap_ind_p->radio_capabilities_length);
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);

  // Constructing and sending Attach Complete to mme_app
  // mimicing S1AP
  uint8_t nas_msg4[] = {0x27, 0xb6, 0x28, 0x5a, 0x49, 0x01, 0x07,
                        0x43, 0x00, 0x03, 0x52, 0x00, 0xc2};
  nas_msg_length     = 13;
  send_mme_app_uplink_data_ind(&nas_msg4[0], nas_msg_length, plmn);

  // Check MME state after attach complete
  // Sleep briefly to ensure processing my mme_app
  std::this_thread::sleep_for(std::chrono::milliseconds(200));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 1);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);

  // Constructing and sending Detach Request to mme_app
  // mimicing S1AP
  uint8_t nas_msg5[] = {0x27, 0x8f, 0xf4, 0x06, 0xe5, 0x02, 0x07,
                        0x45, 0x09, 0x0b, 0xf6, 0x00, 0xf1, 0x10,
                        0x00, 0x01, 0x01, 0x46, 0x93, 0xe8, 0xb8};
  nas_msg_length     = 21;
  send_mme_app_uplink_data_ind(&nas_msg5[0], nas_msg_length, plmn);

  // Constructing and sending Delete Session Response to mme_app
  // mimicing SPGW task
  message_p =
      itti_alloc_new_message(TASK_SPGW_APP, S11_DELETE_SESSION_RESPONSE);
  itti_s11_delete_session_response_t* delete_session_resp_p =
      &message_p->ittiMsg.s11_delete_session_response;
  delete_session_resp_p->cause.cause_value = REQUEST_ACCEPTED;
  delete_session_resp_p->teid              = 1;
  delete_session_resp_p->peer_ip.s_addr    = 100;
  delete_session_resp_p->lbi               = 5;
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);

  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  message_p =
      itti_alloc_new_message(TASK_S1AP, S1AP_UE_CONTEXT_RELEASE_COMPLETE);
  S1AP_UE_CONTEXT_RELEASE_COMPLETE(message_p).mme_ue_s1ap_id = 1;
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);

  // Check MME state after detach complete
  // Sleep briefly to ensure processing my mme_app
  std::this_thread::sleep_for(std::chrono::milliseconds(200));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 0);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 0);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(500));
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  OAILOG_INIT("MME", OAILOG_LEVEL_DEBUG, MAX_LOG_PROTOS);
  return RUN_ALL_TESTS();
}

}  // namespace lte
}  // namespace magma
