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
#include "mme_config.h"
#include "mme_app_defs.h"
#include "mme_app_extern.h"
#include "mme_app_state.h"
}

task_zmq_ctx_t task_zmq_ctx_main;

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

class MmeAppProcedureTest : public ::testing::Test {
  virtual void SetUp() {
    s1ap_handler = std::make_shared<MockS1apHandler>();

    itti_init(
        TASK_MAX, THREAD_MAX, MESSAGES_ID_MAX, tasks_info, messages_info, NULL,
        NULL);

    // initialize mme config
    mme_config_init(&mme_config);
    mme_config.use_stateless = false;

    task_id_t task_id_list[10] = {
        TASK_MME_APP,    TASK_HA,  TASK_S1AP,   TASK_S6A,      TASK_S11,
        TASK_SERVICE303, TASK_SGS, TASK_SGW_S8, TASK_SPGW_APP, TASK_SMS_ORC8R};
    init_task_context(
        TASK_MAIN, task_id_list, 10, handle_message, &task_zmq_ctx_main);

    std::thread task_ha(start_mock_ha_task);
    std::thread task_s1ap(start_mock_s1ap_task, s1ap_handler);
    std::thread task_s6a(start_mock_s6a_task);
    std::thread task_s11(start_mock_s11_task);
    std::thread task_service303(start_mock_service303_task);
    std::thread task_sgs(start_mock_sgs_task);
    std::thread task_sgw_s8(start_mock_sgw_s8_task);
    std::thread task_sms_orc8r(start_mock_sms_orc8r_task);
    std::thread task_spgw(start_mock_spgw_task);
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
};

TEST_F(MmeAppProcedureTest, TestInitialUeMessageFaultyNasMsg) {
  MessageDef* message_p = NULL;
  message_p = itti_alloc_new_message(TASK_S1AP, S1AP_INITIAL_UE_MESSAGE);

  /* The following buffer just includes an attach request */
  uint8_t nas_msg[]       = {0x72, 0x08, 0x09, 0x10, 0x10, 0x00, 0x00, 0x00,
                       0x00, 0x10, 0x02, 0xe0, 0xe0, 0x00, 0x04, 0x02,
                       0x01, 0xd0, 0x11, 0x40, 0x08, 0x04, 0x02, 0x60,
                       0x04, 0x00, 0x02, 0x1c, 0x00};
  uint32_t nas_msg_length = 29;

  EXPECT_CALL(*s1ap_handler, s1ap_generate_downlink_nas_transport()).Times(1);

  S1AP_INITIAL_UE_MESSAGE(message_p).sctp_assoc_id  = 0;
  S1AP_INITIAL_UE_MESSAGE(message_p).enb_ue_s1ap_id = 0;
  S1AP_INITIAL_UE_MESSAGE(message_p).enb_id         = 0;
  S1AP_INITIAL_UE_MESSAGE(message_p).nas = blk2bstr(nas_msg, nas_msg_length);
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(1000));
}

TEST_F(MmeAppProcedureTest, TestInitialAttachEpsOnly) {
  MessageDef* message_p = NULL;
  message_p = itti_alloc_new_message(TASK_S1AP, S1AP_INITIAL_UE_MESSAGE);

  uint8_t nas_msg[]       = {0x07, 0x41, 0x71, 0x08, 0x09, 0x10, 0x10, 0x00,
                       0x00, 0x00, 0x00, 0x10, 0x02, 0xe0, 0xe0, 0x00,
                       0x04, 0x02, 0x01, 0xd0, 0x11, 0x40, 0x08, 0x04,
                       0x02, 0x60, 0x04, 0x00, 0x02, 0x1c, 0x00};
  uint32_t nas_msg_length = 31;

  std::string imsi = "001010000000001";

  EXPECT_CALL(*s1ap_handler, s1ap_generate_downlink_nas_transport()).Times(1);

  // Sending Initial Attach Request to mme_app mimicing S1AP
  S1AP_INITIAL_UE_MESSAGE(message_p).sctp_assoc_id  = 0;
  S1AP_INITIAL_UE_MESSAGE(message_p).enb_ue_s1ap_id = 0;
  S1AP_INITIAL_UE_MESSAGE(message_p).enb_id         = 0;
  S1AP_INITIAL_UE_MESSAGE(message_p).nas = blk2bstr(nas_msg, nas_msg_length);
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);

  // Sending AIA to mme_app mimicing successful S6A response
  message_p = itti_alloc_new_message(TASK_S6A, S6A_AUTH_INFO_ANS);
  s6a_auth_info_ans_t* itti_msg = &message_p->ittiMsg.s6a_auth_info_ans;
  strncpy(itti_msg->imsi, imsi.c_str(), imsi.size());
  itti_msg->imsi_length        = imsi.size();
  itti_msg->result.present     = S6A_RESULT_BASE;
  itti_msg->result.choice.base = DIAMETER_SUCCESS;
  magma::feg::AuthenticationInformationAnswer aia;
  magma::feg::AuthenticationInformationAnswer::EUTRANVector eutran_vector;
  uint8_t xres_buf[17]  = {0xe1, 0xda, 0xf7, 0x88, 0x9d, 0xf7, 0x82, 0x68,
                          0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00};
  xres_buf[17]          = '\0';
  uint8_t rand_buf[17]  = {0x99, 0xb0, 0x44, 0xec, 0xc8, 0x83, 0xfd, 0xa1,
                          0x87, 0x2a, 0xdc, 0xe4, 0xc6, 0xe8, 0x27, 0xe7};
  rand_buf[17]          = '\0';
  uint8_t autn_buf[17]  = {0x98, 0xff, 0xee, 0x81, 0xce, 0x05, 0x80, 0x00,
                          0x2d, 0x60, 0xe4, 0xc0, 0xf0, 0xf4, 0xa0, 0x7a};
  autn_buf[17]          = '\0';
  uint8_t kasme_buf[33] = {0xbc, 0x5b, 0x76, 0x5f, 0xf3, 0xa3, 0x1a, 0x64,
                           0x30, 0x32, 0x27, 0x82, 0x5b, 0xfd, 0xef, 0x24,
                           0x8b, 0x81, 0x4e, 0x97, 0x50, 0xe5, 0x89, 0x94,
                           0xd7, 0x17, 0x38, 0x97, 0xfc, 0xbe, 0xea, 0xe4};
  kasme_buf[33]         = '\0';
  eutran_vector.set_rand((const char*) rand_buf);
  eutran_vector.set_xres((const char*) xres_buf);
  eutran_vector.set_autn((const char*) autn_buf);
  eutran_vector.set_kasme((const char*) kasme_buf); 
  aia.set_error_code(magma::feg::ErrorCode::SUCCESS);
  auto eutran_vectors = aia.mutable_eutran_vectors();
  eutran_vectors->Add()->CopyFrom(eutran_vector);
  magma::convert_proto_msg_to_itti_s6a_auth_info_ans(aia, itti_msg);
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(1000));
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  OAILOG_INIT("MME", OAILOG_LEVEL_DEBUG, MAX_LOG_PROTOS);
  return RUN_ALL_TESTS();
}
