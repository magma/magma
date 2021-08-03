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

#include "../mock_tasks/mock_tasks.h"
#include "mme_app_state_manager.h"
#include "mme_app_ip_imsi.h"

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
    // log_init(MME_CONFIG_STRING_MME_CONFIG, OAILOG_LEVEL_DEBUG,
    // MAX_LOG_PROTOS);
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
    std::thread task_s1ap(start_mock_s1ap_task);
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
    // Sleep to ensure that messages are received and contexts are released
    std::this_thread::sleep_for(std::chrono::milliseconds(1000));
  }
};

TEST_F(MmeAppProcedureTest, TestInitialUeMessage) {
  MessageDef* message_p = NULL;
  message_p = itti_alloc_new_message(TASK_S1AP, S1AP_INITIAL_UE_MESSAGE);

  uint8_t nas_msg[]       = {0x72, 0x08, 0x09, 0x10, 0x10, 0x00, 0x00, 0x00,
                       0x00, 0x10, 0x02, 0xe0, 0xe0, 0x00, 0x04, 0x02,
                       0x01, 0xd0, 0x11, 0x40, 0x08, 0x04, 0x02, 0x60,
                       0x04, 0x00, 0x02, 0x1c, 0x00};
  uint32_t nas_msg_length = 29;

  S1AP_INITIAL_UE_MESSAGE(message_p).sctp_assoc_id  = 0;
  S1AP_INITIAL_UE_MESSAGE(message_p).enb_ue_s1ap_id = 0;
  S1AP_INITIAL_UE_MESSAGE(message_p).enb_id         = 0;
  S1AP_INITIAL_UE_MESSAGE(message_p).nas = blk2bstr(nas_msg, nas_msg_length);
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);
  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(1000));
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  OAILOG_INIT("MME", OAILOG_LEVEL_DEBUG, MAX_LOG_PROTOS);
  return RUN_ALL_TESTS();
}
