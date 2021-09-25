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
#include "mme_app_test_util.h"

extern "C" {
#include "common_types.h"
#include "mme_config.h"
#include "mme_app_extern.h"
#include "mme_app_state.h"
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
    mme_config.use_stateless                              = true;
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
  plmn_t plmn = {.mcc_digit2 = 0,
                 .mcc_digit1 = 0,
                 .mnc_digit3 = 0x0f,
                 .mcc_digit3 = 1,
                 .mnc_digit2 = 1,
                 .mnc_digit1 = 0};

  EXPECT_CALL(*s1ap_handler, s1ap_generate_downlink_nas_transport()).Times(1);

  // Construction and sending Initial Attach Request to mme_app mimicing S1AP
  // The following buffer just includes an attach request
  uint8_t nas_msg[]       = {0x72, 0x08, 0x09, 0x10, 0x10, 0x00, 0x00, 0x00,
                       0x00, 0x10, 0x02, 0xe0, 0xe0, 0x00, 0x04, 0x02,
                       0x01, 0xd0, 0x11, 0x40, 0x08, 0x04, 0x02, 0x60,
                       0x04, 0x00, 0x02, 0x1c, 0x00};
  uint32_t nas_msg_length = 29;
  send_mme_app_initial_ue_msg(&nas_msg[0], nas_msg_length, plmn);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(1000));
}

TEST_F(MmeAppProcedureTest, TestImsiAttachEpsOnlyDetach) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::string imsi = "001010000000001";
  plmn_t plmn      = {.mcc_digit2 = 0,
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
  uint8_t nas_msg[]       = {0x07, 0x41, 0x71, 0x08, 0x09, 0x10, 0x10, 0x00,
                       0x00, 0x00, 0x00, 0x10, 0x02, 0xe0, 0xe0, 0x00,
                       0x04, 0x02, 0x01, 0xd0, 0x11, 0x40, 0x08, 0x04,
                       0x02, 0x60, 0x04, 0x00, 0x02, 0x1c, 0x00};
  uint32_t nas_msg_length = 31;
  send_mme_app_initial_ue_msg(&nas_msg[0], nas_msg_length, plmn);

  // Sending AIA to mme_app mimicing successful S6A response for AIR
  send_authentication_info_resp(imsi);

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
  send_s6a_ula(imsi);

  // Constructing and sending Create Session Response to mme_app mimicing SPGW
  send_create_session_resp();

  // Constructing and sending ICS Response to mme_app mimicing S1AP
  send_ics_response();

  // Constructing UE Capability Indication message to mme_app
  // mimicing S1AP with dummy radio capabilities
  send_ue_capabilities_ind();

  // Constructing and sending Attach Complete to mme_app
  // mimicing S1AP
  uint8_t nas_msg4[] = {0x27, 0xb6, 0x28, 0x5a, 0x49, 0x01, 0x07,
                        0x43, 0x00, 0x03, 0x52, 0x00, 0xc2};
  nas_msg_length     = 13;
  send_mme_app_uplink_data_ind(&nas_msg4[0], nas_msg_length, plmn);

  // Check MME state after attach complete
  // Sleep briefly to ensure processing my mme_app
  std::this_thread::sleep_for(std::chrono::milliseconds(300));
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
  send_delete_session_resp();

  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

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

TEST_F(MmeAppProcedureTest, TestGutiAttachEpsOnlyDetach) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::string imsi = "001010000000001";
  plmn_t plmn      = {.mcc_digit2 = 0,
                 .mcc_digit1 = 0,
                 .mnc_digit3 = 0x0f,
                 .mcc_digit3 = 1,
                 .mnc_digit2 = 1,
                 .mnc_digit1 = 0};

  EXPECT_CALL(*s1ap_handler, s1ap_generate_downlink_nas_transport()).Times(4);
  EXPECT_CALL(*s1ap_handler, s1ap_handle_conn_est_cnf()).Times(1);
  EXPECT_CALL(*s1ap_handler, s1ap_handle_ue_context_release_command()).Times(1);
  EXPECT_CALL(*s6a_handler, s6a_viface_authentication_info_req()).Times(1);
  EXPECT_CALL(*s6a_handler, s6a_viface_update_location_req()).Times(1);
  EXPECT_CALL(*s6a_handler, s6a_viface_purge_ue()).Times(1);
  EXPECT_CALL(*spgw_handler, sgw_handle_s11_create_session_request()).Times(1);
  EXPECT_CALL(*spgw_handler, sgw_handle_delete_session_request()).Times(1);

  // Construction and sending Initial Attach Request to mme_app mimicing S1AP
  uint8_t nas_msg0[] = {0x07, 0x41, 0x71, 0x0b, 0xf6, 0x00, 0x00, 0x00, 0x00,
                        0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xe0, 0xe0,
                        0x00, 0x04, 0x02, 0x01, 0xd0, 0x11, 0x40, 0x08, 0x04,
                        0x02, 0x60, 0x04, 0x00, 0x02, 0x1c, 0x00};
  uint32_t nas_msg_length = 34;
  send_mme_app_initial_ue_msg(&nas_msg0[0], nas_msg_length, plmn);

  // Constructing and sending Identity Response to mme_app mimicing S1AP
  uint8_t nas_msg1[] = {0x07, 0x56, 0x08, 0x09, 0x10, 0x10,
                        0x00, 0x00, 0x00, 0x00, 0x10};
  nas_msg_length     = 11;
  send_mme_app_uplink_data_ind(&nas_msg1[0], nas_msg_length, plmn);

  // Sending AIA to mme_app mimicing successful S6A response for AIR
  send_authentication_info_resp(imsi);

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
  send_s6a_ula(imsi);

  // Constructing and sending Create Session Response to mme_app mimicing SPGW
  send_create_session_resp();

  // Constructing and sending ICS Response to mme_app mimicing S1AP
  send_ics_response();

  // Constructing UE Capability Indication message to mme_app
  // mimicing S1AP with dummy radio capabilities
  send_ue_capabilities_ind();

  // Constructing and sending Attach Complete to mme_app
  // mimicing S1AP
  uint8_t nas_msg4[] = {0x27, 0xb6, 0x28, 0x5a, 0x49, 0x01, 0x07,
                        0x43, 0x00, 0x03, 0x52, 0x00, 0xc2};
  nas_msg_length     = 13;
  send_mme_app_uplink_data_ind(&nas_msg4[0], nas_msg_length, plmn);

  // Check MME state after attach complete
  // Sleep briefly to ensure processing my mme_app
  std::this_thread::sleep_for(std::chrono::milliseconds(300));
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
  send_delete_session_resp();

  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

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

}  // namespace lte
}  // namespace magma
