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

#pragma once

#include <gtest/gtest.h>
#include <czmq.h>
#include <mutex>
#include <thread>
#include <cstdint>
#include <condition_variable>

extern "C" {
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
}

#include "lte/gateway/c/core/oai/include/mme_app_desc.hpp"
#include "lte/gateway/c/core/oai/include/mme_config.hpp"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_extern.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/api/network/nas_message.hpp"
#include "lte/gateway/c/core/oai/test/mock_tasks/mock_tasks.hpp"
#include "lte/gateway/c/core/oai/test/mme_app_task/mme_app_test_util.hpp"

extern bool mme_hss_associated;
extern bool mme_sctp_bounded;
namespace magma {
namespace lte {

ACTION_P(ReturnFromAsyncTask, cv) { cv->notify_all(); }

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

MATCHER_P2(check_params_in_path_switch_req_ack, new_enb_ue_s1ap_id,
           new_sctp_assoc_id, "") {
  auto path_switch_ack_recv =
      static_cast<itti_s1ap_path_switch_request_ack_t>(arg);
  if ((path_switch_ack_recv.enb_ue_s1ap_id == new_enb_ue_s1ap_id) &&
      (path_switch_ack_recv.sctp_assoc_id == new_sctp_assoc_id)) {
    return true;
  }
  return false;
}

MATCHER_P2(check_params_in_path_switch_req_failure, new_enb_ue_s1ap_id,
           new_sctp_assoc_id, "") {
  auto path_switch_ack_recv =
      static_cast<itti_s1ap_path_switch_request_failure_t>(arg);
  if ((path_switch_ack_recv.enb_ue_s1ap_id == new_enb_ue_s1ap_id) &&
      (path_switch_ack_recv.sctp_assoc_id == new_sctp_assoc_id)) {
    return true;
  }
  return false;
}

MATCHER_P2(check_params_in_mme_app_handover_request, mme_ue_s1ap_id,
           new_sctp_assoc_id, "") {
  auto mme_app_handover_request_recv =
      static_cast<itti_mme_app_handover_request_t>(arg);
  if ((mme_app_handover_request_recv.mme_ue_s1ap_id == mme_ue_s1ap_id) &&
      (mme_app_handover_request_recv.target_sctp_assoc_id ==
       new_sctp_assoc_id)) {
    return true;
  }
  return false;
}

MATCHER_P2(check_params_in_mme_app_handover_command, mme_ue_s1ap_id, new_enb_id,
           "") {
  auto mme_app_handover_command_recv =
      static_cast<itti_mme_app_handover_command_t>(arg);
  if ((mme_app_handover_command_recv.target_enb_id == new_enb_id) &&
      (mme_app_handover_command_recv.mme_ue_s1ap_id == mme_ue_s1ap_id)) {
    return true;
  }
  return false;
}

class MmeAppProcedureTest : public ::testing::Test {
  virtual void SetUp() {
    mme_hss_associated = false;
    mme_sctp_bounded = false;
    s1ap_handler = std::make_shared<MockS1apHandler>();
    s6a_handler = std::make_shared<MockS6aHandler>();
    s8_handler = std::make_shared<MockS8Handler>();
    spgw_handler = std::make_shared<MockSpgwHandler>();
    service303_handler = std::make_shared<MockService303Handler>();
    itti_init(TASK_MAX, THREAD_MAX, MESSAGES_ID_MAX, tasks_info, messages_info,
              NULL, NULL);

    // initialize mme config
    mme_config_init(&mme_config);
    create_partial_lists(&mme_config);
    mme_config.use_stateless = true;
    mme_config.nas_config.prefered_integrity_algorithm[0] = EIA2_128_ALG_ID;

    task_id_t task_id_list[10] = {
        TASK_MME_APP,    TASK_HA,  TASK_S1AP,   TASK_S6A,      TASK_S11,
        TASK_SERVICE303, TASK_SGS, TASK_SGW_S8, TASK_SPGW_APP, TASK_SMS_ORC8R};
    init_task_context(TASK_MAIN, task_id_list, 10, handle_message,
                      &task_zmq_ctx_main);

    std::thread task_s6a(start_mock_s6a_task, s6a_handler);
    std::thread task_s1ap(start_mock_s1ap_task, s1ap_handler);
    std::thread task_spgw(start_mock_spgw_task, spgw_handler);
    std::thread task_ha(start_mock_ha_task);
    std::thread task_s11(start_mock_s11_task);
    std::thread task_service303(start_mock_service303_task, service303_handler);
    std::thread task_sgs(start_mock_sgs_task);
    std::thread task_sgw_s8(start_mock_sgw_s8_task, s8_handler);
    std::thread task_sms_orc8r(start_mock_sms_orc8r_task);

    task_s6a.detach();
    task_s1ap.detach();
    task_spgw.detach();
    task_ha.detach();
    task_s11.detach();
    task_service303.detach();
    task_sgs.detach();
    task_sgw_s8.detach();
    task_sms_orc8r.detach();

    // Sleep for 10 milliseconds to make sure all mock tasks
    // are running before the test starts
    std::this_thread::sleep_for(std::chrono::milliseconds(10));

    mme_app_init(&mme_config);
    // Fake initialize sctp server.
    // We can then send activate messages in each test
    // whenever we need to read mme_app state.
    send_sctp_mme_server_initialized();
  }

  virtual void TearDown() {
    bdestroy_wrapper(&nas_msg);
    send_terminate_message_fatal(&task_zmq_ctx_main);
    // Sleep to ensure that messages are received and contexts are released
    std::this_thread::sleep_for(std::chrono::milliseconds(500));
    destroy_task_context(&task_zmq_ctx_main);
    itti_free_desc_threads();
  }

 protected:
  itti_s1ap_nas_dl_data_req_t msg_nas_dl_data = {0};
  bstring nas_msg = NULL;
  std::shared_ptr<MockS1apHandler> s1ap_handler;
  std::shared_ptr<MockS6aHandler> s6a_handler;
  std::shared_ptr<MockSpgwHandler> spgw_handler;
  std::shared_ptr<MockService303Handler> service303_handler;
  std::shared_ptr<MockS8Handler> s8_handler;
  const uint8_t nas_msg_imsi_attach_req[31] = {
      0x07, 0x41, 0x71, 0x08, 0x09, 0x10, 0x10, 0x00, 0x00, 0x00, 0x00,
      0x10, 0x02, 0xe0, 0xe0, 0x00, 0x04, 0x02, 0x01, 0xd0, 0x11, 0x40,
      0x08, 0x04, 0x02, 0x60, 0x04, 0x00, 0x02, 0x1c, 0x00};
  const uint8_t nas_msg_guti_attach_req[34] = {
      0x07, 0x41, 0x71, 0x0b, 0xf6, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
      0x00, 0x00, 0x00, 0x02, 0xe0, 0xe0, 0x00, 0x04, 0x02, 0x01, 0xd0, 0x11,
      0x40, 0x08, 0x04, 0x02, 0x60, 0x04, 0x00, 0x02, 0x1c, 0x00};
  const uint8_t nas_msg_auth_resp[19] = {
      0x07, 0x53, 0x10, 0x66, 0xff, 0x47, 0x2d, 0xd4, 0x93, 0xf1,
      0x5a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00};
  const uint8_t nas_msg_smc_resp[19] = {
      0x47, 0xc0, 0xb5, 0x35, 0x6b, 0x00, 0x07, 0x5e, 0x23, 0x09,
      0x33, 0x08, 0x45, 0x86, 0x34, 0x12, 0x31, 0x71, 0xf2};
  const uint8_t nas_msg_ident_resp[11] = {0x07, 0x56, 0x08, 0x09, 0x10, 0x10,
                                          0x00, 0x00, 0x00, 0x00, 0x10};
  const uint8_t nas_msg_attach_comp[13] = {0x27, 0xb6, 0x28, 0x5a, 0x49,
                                           0x01, 0x07, 0x43, 0x00, 0x03,
                                           0x52, 0x00, 0xc2};
  const uint8_t nas_msg_detach_req[21] = {
      0x27, 0x8f, 0xf4, 0x06, 0xe5, 0x02, 0x07, 0x45, 0x09, 0x0b, 0xf6,
      0x00, 0xf1, 0x10, 0x00, 0x01, 0x01, 0x46, 0x93, 0xe8, 0xb8};
  const uint8_t nas_msg_detach_accept[8] = {0x17, 0x88, 0x16, 0x67,
                                            0xd3, 0x02, 0x07, 0x46};
  const uint8_t nas_msg_service_req[4] = {0xc7, 0x02, 0x79, 0xe0};
  const uint8_t nas_msg_activate_ded_bearer_accept[9] = {
      0x27, 0xAA, 0x95, 0x47, 0x92, 0x02, 0x62, 0x00, 0xc6};
  const uint8_t nas_msg_activate_ded_bearer_accept_ebi_7[9] = {
      0x27, 0xe2, 0xbc, 0xb5, 0x3f, 0x04, 0x72, 0x00, 0xc6};
  const uint8_t nas_msg_deactivate_ded_bearer_accept[9] = {
      0x27, 0x66, 0x5f, 0x4e, 0x87, 0x03, 0x62, 0x00, 0xce};

  const uint8_t nas_msg_periodic_tau_req_with_actv_flag[21] = {
      0x17, 0x1F, 0x2C, 0x60, 0x5E, 0x02, 0x07, 0x48, 0x0b, 0x0b, 0xf6,
      0x00, 0xf1, 0x10, 0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x01,
  };

  const uint8_t nas_msg_periodic_tau_req_without_actv_flag[21] = {
      0x17, 0xE6, 0xDE, 0x80, 0xD4, 0x02, 0x07, 0x48, 0x03, 0x0b, 0xf6,
      0x00, 0xf1, 0x10, 0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x01,
  };
  const uint8_t nas_msg_normal_tau_req_with_actv_flag[21] = {
      0x17, 0xA,  0x7D, 0x19, 0x5F, 0x02, 0x07, 0x48, 0x08, 0x0b, 0xf6,
      0x00, 0xf1, 0x10, 0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x01,
  };

  const uint8_t nas_msg_normal_tau_req_without_actv_flag[21] = {
      0x17, 0xEA, 0x57, 0xCB, 0xAD, 0x02, 0x07, 0x48, 0x00, 0x0b, 0xf6,
      0x00, 0xf1, 0x10, 0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x01,
  };

  const uint8_t nas_msg_tau_req_with_eps_bearer_ctx_sts_def_ber[25] = {
      0x17, 0x4A, 0xE7, 0x23, 0x9E, 0x05, 0x07, 0x48, 0x08,
      0x0b, 0xf6, 0x00, 0xf1, 0x10, 0x00, 0x01, 0x01, 0x67,
      0x41, 0x14, 0xf4, 0x57, 0x02, 0x20, 0x00};

  const uint8_t nas_msg_tau_req_with_eps_bearer_ctx_sts_ded_ber[25] = {
      0x17, 0x18, 0xCF, 0x11, 0x4E, 0x05, 0x07, 0x48, 0x08,
      0x0b, 0xf6, 0x00, 0xf1, 0x10, 0x00, 0x01, 0x01, 0x67,
      0x41, 0x14, 0xf4, 0x57, 0x02, 0x60, 0x00};

  const uint8_t nas_msg_pdn_connectivity_req_ims[16] = {
      0x27, 0x13, 0xd2, 0x79, 0xbe, 0x02, 0x02, 0x01,
      0xd0, 0x11, 0x28, 0x04, 0x03, 0x69, 0x6d, 0x73};

  const uint8_t nas_msg_act_def_bearer_acc[9] = {0x27, 0xaa, 0x95, 0x47, 0x92,
                                                 0x03, 0x62, 0x00, 0xc2};
  const uint8_t nas_msg_activate_ded_bearer_reject[10] = {
      0x27, 0xe8, 0xf5, 0xb5, 0xf1, 0x02, 0x62, 0x00, 0xc7, 0x00};

  std::string imsi = "001010000000001";
  plmn_t plmn = {.mcc_digit2 = 0,
                 .mcc_digit1 = 0,
                 .mnc_digit3 = 0x0f,
                 .mcc_digit3 = 1,
                 .mnc_digit2 = 1,
                 .mnc_digit1 = 0};
  guti_eps_mobile_identity_t guti = {0};

  void attach_ue(std::condition_variable& cv,
                 std::unique_lock<std::mutex>& lock,
                 mme_app_desc_t* mme_state_p, guti_eps_mobile_identity_t* guti);

  void detach_ue(std::condition_variable& cv,
                 std::unique_lock<std::mutex>& lock,
                 mme_app_desc_t* mme_state_p, guti_eps_mobile_identity_t guti,
                 bool is_initial_ue);
};

}  // namespace lte
}  // namespace magma
