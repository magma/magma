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

#include "lte/gateway/c/core/oai/test/mock_tasks/mock_tasks.h"

extern "C" {
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/include/amf_config.h"
#include "lte/gateway/c/core/oai/include/amf_app_messages_types.h"
}
#include "lte/gateway/c/core/oai/tasks/amf/include/amf_client_servicer.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_state_manager.h"
#include "lte/gateway/c/core/oai/test/amf/amf_app_test_util.h"
#include "lte/gateway/c/core/oai/test/amf/util_s6a_update_location.h"

using ::testing::Test;

namespace magma5g {

extern task_zmq_ctx_s amf_app_task_zmq_ctx;
extern std::unordered_map<amf_ue_ngap_id_t, ue_m5gmm_context_s*> ue_context_map;

class AMFAppProcedureTest : public ::testing::Test {
  virtual void SetUp() {
    itti_init(TASK_MAX, THREAD_MAX, MESSAGES_ID_MAX, tasks_info, messages_info,
              NULL, NULL);

    // initialize mme config
    amf_config_init(&amf_config);
    amf_nas_state_init(&amf_config);
    create_state_matrix();

    init_task_context(TASK_MAIN, nullptr, 0, NULL, &amf_app_task_zmq_ctx);

    char dnn[] = "internet";
    amf_config.default_dnn = bfromcstr(dnn);
    inet_aton("8.8.8.8", &(amf_config.ipv4.default_dns));
    inet_aton("8.8.10.10", &(amf_config.ipv4.default_dns_sec));

    /* Initialize the plmn */
    amf_config.guamfi.nb = 1;
    amf_config.guamfi.guamfi[0].plmn = {.mcc_digit2 = 2,
                                        .mcc_digit1 = 2,
                                        .mnc_digit3 = 6,
                                        .mcc_digit3 = 2,
                                        .mnc_digit2 = 5,
                                        .mnc_digit1 = 4};

    amf_app_desc_p = get_amf_nas_state(false);
    AMFClientServicer::getInstance().msgtype_stack.clear();
  }

  virtual void TearDown() {
    clear_amf_nas_state();
    clear_amf_config(&amf_config);
    destroy_task_context(&amf_app_task_zmq_ctx);
    itti_free_desc_threads();
    AMFClientServicer::getInstance().msgtype_stack.clear();
  }

 protected:
  amf_app_desc_t* amf_app_desc_p;
  std::string imsi = "222456000000001";
  plmn_t plmn = {.mcc_digit2 = 2,
                 .mcc_digit1 = 2,
                 .mnc_digit3 = 6,
                 .mcc_digit3 = 2,
                 .mnc_digit2 = 5,
                 .mnc_digit1 = 4};

  itti_amf_decrypted_imsi_info_ans_t decrypted_imsi = {
      .imsi = "222456000000001", .imsi_length = 15, .result = 1, .ue_id = 1};

  const uint8_t intital_ue_message_suci_ext_hexbuf[65] = {
      0x7e, 0x00, 0x41, 0x79, 0x00, 0x35, 0x01, 0x09, 0xf1, 0x07, 0x00,
      0x00, 0x01, 0x04, 0xc8, 0xfc, 0x0c, 0xe5, 0x47, 0x9a, 0x51, 0x5d,
      0xab, 0xf2, 0xf3, 0x45, 0xae, 0xb4, 0x66, 0x92, 0xd6, 0xff, 0x7a,
      0x5f, 0x4f, 0x57, 0x2a, 0x47, 0x99, 0xf2, 0x33, 0x69, 0x35, 0x16,
      0x40, 0x31, 0xbd, 0x3f, 0x84, 0x41, 0x26, 0xdf, 0x5b, 0x47, 0x06,
      0x41, 0xe2, 0xa9, 0x57, 0x2e, 0x04, 0xf0, 0xf0, 0xf0, 0xf0};

  const uint8_t initial_ue_message_hexbuf[25] = {
      0x7e, 0x00, 0x41, 0x79, 0x00, 0x0d, 0x01, 0x22, 0x62,
      0x54, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
      0x01, 0x2e, 0x04, 0xf0, 0xf0, 0xf0, 0xf0};

  /* Guti based initial registration message */
  const uint8_t guti_initial_ue_message_hexbuf[98] = {
      0x7e, 0x01, 0x67, 0xb7, 0x6f, 0xc6, 0x03, 0x7e, 0x00, 0x41, 0x09,
      0x00, 0x0b, 0xf2, 0x22, 0x62, 0x54, 0x01, 0x00, 0x40, 0x12, 0x70,
      0x41, 0x77, 0x2e, 0x04, 0x80, 0xe0, 0x80, 0xe0, 0x71, 0x00, 0x41,
      0x7e, 0x00, 0x41, 0x09, 0x00, 0x0b, 0xf2, 0x22, 0x62, 0x54, 0x01,
      0x00, 0x40, 0x12, 0x70, 0x41, 0x77, 0x10, 0x01, 0x03, 0x2e, 0x04,
      0x80, 0xe0, 0x80, 0xe0, 0x2f, 0x02, 0x01, 0x01, 0x52, 0x22, 0x62,
      0x54, 0x00, 0x00, 0x01, 0x17, 0x07, 0x80, 0xe0, 0xe0, 0x60, 0x00,
      0x1c, 0x30, 0x18, 0x01, 0x00, 0x74, 0x00, 0x0a, 0x09, 0x08, 0x69,
      0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x53, 0x01, 0x01};

  /* Mobile Termination as the initial ue message */
  const uint8_t mu_initial_ue_message_hexbuf[93] = {
      0x7e, 0x01, 0xa3, 0xcf, 0x4c, 0x7e, 0xd1, 0x7e, 0x00, 0x41, 0x02, 0x00,
      0x0b, 0xf2, 0x22, 0x62, 0x54, 0x01, 0x00, 0x40, 0x8f, 0x71, 0x9f, 0x0e,
      0x2e, 0x04, 0xf0, 0x70, 0xf0, 0x70, 0x71, 0x00, 0x3c, 0x7e, 0x00, 0x41,
      0x02, 0x00, 0x0b, 0xf2, 0x22, 0x62, 0x54, 0x01, 0x00, 0x40, 0x8f, 0x71,
      0x9f, 0x0e, 0x10, 0x01, 0x03, 0x2e, 0x04, 0xf0, 0x70, 0xf0, 0x70, 0x2f,
      0x02, 0x01, 0x01, 0x52, 0x22, 0x62, 0x54, 0x00, 0x00, 0x01, 0x17, 0x07,
      0xf0, 0x70, 0x00, 0x00, 0x18, 0x80, 0xb0, 0x50, 0x02, 0x00, 0x00, 0x18,
      0x01, 0x01, 0x74, 0x00, 0x00, 0x90, 0x53, 0x01, 0x01};

  const uint8_t identity_response[25] = {
      0x7e, 0x01, 0xe2, 0x1d, 0xd9, 0xe5, 0x04, 0x7e, 0x00,
      0x5c, 0x00, 0x0d, 0x01, 0x22, 0x62, 0x54, 0xf0, 0xff,
      0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01};

  const uint8_t ue_auth_response_hexbuf[21] = {
      0x7e, 0x0,  0x57, 0x2d, 0x10, 0x25, 0x70, 0x6f, 0x9a, 0x5b, 0x90,
      0xb6, 0xc9, 0x57, 0x50, 0x6c, 0x88, 0x3d, 0x76, 0xcc, 0x63};

  const uint8_t ue_auth_response_security_capability_mismatch_hexbuf[4] = {
      0x7e, 0x0, 0x59, 0x17};

  const uint8_t ue_auth_response_security_mode_reject_hexbuf[4] = {0x7e, 0x0,
                                                                   0x59, 0x18};

  const uint8_t ue_smc_response_hexbuf[60] = {
      0x7e, 0x4,  0x54, 0xf6, 0xe1, 0x2a, 0x0,  0x7e, 0x0,  0x5e, 0x77, 0x0,
      0x9,  0x45, 0x73, 0x80, 0x61, 0x21, 0x85, 0x61, 0x51, 0xf1, 0x71, 0x0,
      0x23, 0x7e, 0x0,  0x41, 0x79, 0x0,  0xd,  0x1,  0x22, 0x62, 0x54, 0x0,
      0x0,  0x0,  0x0,  0x0,  0x0,  0x0,  0x0,  0xf1, 0x10, 0x1,  0x0,  0x2e,
      0x4,  0xf0, 0xf0, 0xf0, 0xf0, 0x2f, 0x2,  0x1,  0x1,  0x53, 0x1,  0x0};

  const uint8_t ue_registration_complete_hexbuf[10] = {
      0x7e, 0x02, 0x5d, 0x5f, 0x49, 0x18, 0x01, 0x7e, 0x00, 0x43};

  const uint8_t ue_pdu_session_est_req_hexbuf[44] = {
      0x7e, 0x00, 0x67, 0x01, 0x00, 0x15, 0x2e, 0x01, 0x01, 0xc1, 0xff,
      0xff, 0x91, 0xa1, 0x28, 0x01, 0x00, 0x7b, 0x00, 0x07, 0x80, 0x00,
      0x0a, 0x00, 0x00, 0x0d, 0x00, 0x12, 0x01, 0x81, 0x22, 0x01, 0x01,
      0x25, 0x09, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74};

  const uint8_t ue_pdu_session_est_req_no_dnn_no_slice_hexbuf[67] = {
      0x7e, 0x00, 0x67, 0x01, 0x00, 0x3a, 0x2e, 0x01, 0x01, 0xc1, 0x00, 0x00,
      0x91, 0x28, 0x01, 0x00, 0x7b, 0x00, 0x2d, 0x80, 0x80, 0x21, 0x10, 0x01,
      0x00, 0x00, 0x10, 0x81, 0x06, 0x00, 0x00, 0x00, 0x00, 0x83, 0x06, 0x00,
      0x00, 0x00, 0x00, 0x00, 0x0d, 0x00, 0x00, 0x0a, 0x00, 0x00, 0x05, 0x00,
      0x00, 0x10, 0x00, 0x00, 0x11, 0x00, 0x00, 0x17, 0x01, 0x01, 0x00, 0x23,
      0x00, 0x00, 0x24, 0x00, 0x12, 0x01, 0x81};

  const uint8_t ue_pdu_session_v6_est_req_hexbuf[44] = {
      0x7e, 0x00, 0x67, 0x01, 0x00, 0x15, 0x2e, 0x01, 0x01, 0xc1, 0xff,
      0xff, 0x92, 0xa1, 0x28, 0x01, 0x00, 0x7b, 0x00, 0x07, 0x80, 0x00,
      0x0a, 0x00, 0x00, 0x0d, 0x00, 0x12, 0x01, 0x81, 0x22, 0x01, 0x01,
      0x25, 0x09, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74};

  const uint8_t ue_pdu_session_v4v6_est_req_hexbuf[44] = {
      0x7e, 0x00, 0x67, 0x01, 0x00, 0x15, 0x2e, 0x01, 0x01, 0xc1, 0xff,
      0xff, 0x93, 0xa1, 0x28, 0x01, 0x00, 0x7b, 0x00, 0x07, 0x80, 0x00,
      0x0a, 0x00, 0x00, 0x0d, 0x00, 0x12, 0x01, 0x81, 0x22, 0x01, 0x01,
      0x25, 0x09, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74};

  const uint8_t ue_pdu_session_est_req_dnn_not_subscried_hexbuf[42] = {
      0x7e, 0x00, 0x67, 0x01, 0x00, 0x15, 0x2e, 0x01, 0x01, 0xc1, 0xff,
      0xff, 0x91, 0xa1, 0x28, 0x01, 0x00, 0x7b, 0x00, 0x07, 0x80, 0x00,
      0x0a, 0x00, 0x00, 0x0d, 0x00, 0x12, 0x01, 0x81, 0x22, 0x01, 0x01,
      0x25, 0x07, 0x06, 0x69, 0x6d, 0x73, 0x2d, 0x35, 0x67};

  const uint8_t ue_pdu_session_est_req_missing_dnn_hexbuf[33] = {
      0x7e, 0x00, 0x67, 0x01, 0x00, 0x15, 0x2e, 0x01, 0x01, 0xc1, 0xff,
      0xff, 0x91, 0xa1, 0x28, 0x01, 0x00, 0x7b, 0x00, 0x07, 0x80, 0x00,
      0x0a, 0x00, 0x00, 0x0d, 0x00, 0x12, 0x01, 0x81, 0x22, 0x01, 0x01,
  };

  const uint8_t ue_pdu_session_est_req_unknown_session_type_hexbuf[44] = {
      0x7e, 0x00, 0x67, 0x01, 0x00, 0x15, 0x2e, 0x01, 0x01, 0xc1, 0xff,
      0xff, 0x94, 0xa1, 0x28, 0x01, 0x00, 0x7b, 0x00, 0x07, 0x80, 0x00,
      0x0a, 0x00, 0x00, 0x0d, 0x00, 0x12, 0x01, 0x81, 0x22, 0x01, 0x01,
      0x25, 0x09, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74};

  const uint8_t pdu_sess_release_hexbuf[14] = {0x7e, 0x00, 0x67, 0x01, 0x00,
                                               0x06, 0x2e, 0x01, 0x01, 0xd1,
                                               0x59, 0x24, 0x12, 0x01};

  const uint8_t pdu_sess_release_complete_hexbuf[12] = {
      0x7e, 0x00, 0x67, 0x01, 0x00, 0x04, 0x2e, 0x01, 0x01, 0xd4, 0x12, 0x01};

  uint8_t ue_initiated_dereg_hexbuf[24] = {
      0x7e, 0x01, 0x41, 0x21, 0xe6, 0xe2, 0x03, 0x7e, 0x00, 0x45, 0x01, 0x00,
      0x0b, 0xf2, 0x22, 0x62, 0x54, 0x01, 0x00, 0x40, 0x0,  0x0,  0x0,  0x0};

  const uint8_t initial_ue_msg_service_request_with_pdu[28] = {
      // Security protected NAS 5G msg with
      // Mobility mgmt msg with security header and auth code
      0x7e, 0x01, 0xca, 0x3f, 0x92, 0xbe,
      // seq no
      0x03,
      // Plain NAS 5G msg with
      // Mobility mgmt msg with no security header and msg type
      0x7e, 0x00, 0x4c,
      // service type as Mobile Terminated
      0x20,
      // 5GS Mobile identity
      0x00, 0x07, 0xf4, 0x00, 0x40,
      // Replace TMSI value to be sent
      // during message tx at position[16]to [19]
      0xff, 0xff, 0xff, 0xff,
      // Uplink data status and pdu session status
      0x40, 0x02, 0x20, 0x00, 0x50, 0x02, 0x20, 0x00};

  const uint8_t initial_ue_msg_service_request_signaling[40] = {
      0x7e, 0x01, 0x30, 0x6f, 0xf2, 0xae, 0x03, 0x7e, 0x00, 0x4c,
      0x00, 0x00, 0x07, 0xf4, 0x00, 0x40, 0xd9, 0x58, 0xf8, 0x3b,
      0x71, 0x00, 0x11, 0x7e, 0x00, 0x4c, 0x00, 0x00, 0x07, 0xf4,
      0x00, 0x40, 0xd9, 0x58, 0xf8, 0x3b, 0x50, 0x02, 0x20, 0x00};
};

amf_context_t* get_amf_context_by_ueid(amf_ue_ngap_id_t ue_id) {
  /* Get UE Context */
  ue_m5gmm_context_s* ue_m5gmm_context =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  if (ue_m5gmm_context == NULL) {
    return NULL;
  }

  /* Get AMF Context */
  amf_context_t* amf_ctx = &ue_m5gmm_context->amf_context;

  return amf_ctx;
}

bool validate_identification_procedure(amf_ue_ngap_id_t ue_id,
                                       uint32_t expected_retransmission_count) {
  amf_context_t* amf_ctx = get_amf_context_by_ueid(ue_id);
  if (amf_ctx == NULL) {
    return false;
  }

  /* Fetch security mode control procedure */
  nas_amf_ident_proc_t* ident_proc =
      get_5g_nas_common_procedure_identification(amf_ctx);
  if (ident_proc == NULL) {
    return false;
  }

  if (ident_proc->retransmission_count != expected_retransmission_count) {
    return false;
  }

  if (ident_proc->ue_id != ue_id) {
    return false;
  }

  return true;
}

bool validate_auth_procedure(amf_ue_ngap_id_t ue_id,
                             uint32_t expected_retransmission_count) {
  amf_context_t* amf_ctx = get_amf_context_by_ueid(ue_id);
  if (amf_ctx == NULL) {
    return false;
  }

  /* Fetch authentication procedure */
  nas5g_amf_auth_proc_t* auth_proc =
      get_nas5g_common_procedure_authentication(amf_ctx);

  if (auth_proc == NULL) {
    return false;
  }

  if (auth_proc->retransmission_count != expected_retransmission_count) {
    return false;
  }

  return true;
}

bool validate_smc_procedure(amf_ue_ngap_id_t ue_id,
                            uint32_t expected_retransmission_count) {
  amf_context_t* amf_ctx = get_amf_context_by_ueid(ue_id);
  if (amf_ctx == NULL) {
    return false;
  }

  /* Fetch security mode control procedure */
  nas_amf_smc_proc_t* smc_proc = get_nas5g_common_procedure_smc(amf_ctx);
  if (smc_proc == NULL) {
    return false;
  }

  if (smc_proc->retransmission_count != expected_retransmission_count) {
    return false;
  }

  if (smc_proc->ue_id != ue_id) {
    return false;
  }

  return true;
}

TEST_F(AMFAppProcedureTest, TestRegistrationAuthSecurityModeReject) {
  amf_ue_ngap_id_t ue_id = 0;
  /* Send the initial UE message */
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));

  /* Check if UE Context is created with correct imsi */
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id));

  /* Send the authentication response message from subscriberdb */
  int rc = RETURNok;
  rc = send_proc_authentication_info_answer(imsi, ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  /* Validate if authentication procedure is initialized as expected */
  EXPECT_TRUE(validate_auth_procedure(ue_id, 0));

  /* Send uplink nas message for auth response from UE */
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, ue_id, plmn, ue_auth_response_security_mode_reject_hexbuf,
      sizeof(ue_auth_response_security_mode_reject_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  ue_m5gmm_context_t* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  // ue context should not exist
  EXPECT_TRUE(ue_context_p == nullptr);
}

TEST_F(AMFAppProcedureTest, TestRegistrationAuthSecurityCapabilityMismatch) {
  amf_ue_ngap_id_t ue_id = 0;
  /* Send the initial UE message */
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));

  /* Check if UE Context is created with correct imsi */
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id));

  /* Send the authentication response message from subscriberdb */
  int rc = RETURNok;
  rc = send_proc_authentication_info_answer(imsi, ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for auth response from UE */
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, ue_id, plmn,
      ue_auth_response_security_capability_mismatch_hexbuf,
      sizeof(ue_auth_response_security_capability_mismatch_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  ue_m5gmm_context_t* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  // ue context should not exist
  EXPECT_TRUE(ue_context_p == nullptr);
}

TEST_F(AMFAppProcedureTest, TestRegistrationProcNoTMSI) {
  amf_ue_ngap_id_t ue_id = 0;
  std::vector<MessagesIds> expected_Ids{
      AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,  // new registration notification
                                            // indication to ngap
      NGAP_NAS_DL_DATA_REQ,                 // Authentication Request to UE
      NGAP_NAS_DL_DATA_REQ,            // Security Command Mode Request to UE
      NGAP_INITIAL_CONTEXT_SETUP_REQ,  // Initial Conext Setup Request to UE &
                                       // Registration Accept
      NGAP_UE_CONTEXT_RELEASE_COMMAND  // UEContextReleaseCommand
  };

  /* Send the initial UE message */
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));

  /* Check if UE Context is created with correct imsi */
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id));

  /* Send the authentication response message from subscriberdb */
  int rc = RETURNok;
  rc = send_proc_authentication_info_answer(imsi, ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  /* Validate if authentication procedure is initialized as expected */
  EXPECT_TRUE(validate_auth_procedure(ue_id, 0));

  /* Send uplink nas message for auth response from UE */
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Check whether security mode procedure is initiated */
  EXPECT_TRUE(validate_smc_procedure(ue_id, 0));

  /* Send uplink nas message for security mode complete response from UE */
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  amf_app_handle_deregistration_req(ue_id);
  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, TestDeRegistration) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id = 0;
  std::vector<MessagesIds> expected_Ids{
      AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,  // new registration notification
                                            // indication to ngap
      NGAP_NAS_DL_DATA_REQ,                 // Authentication Request to UE
      NGAP_NAS_DL_DATA_REQ,            // Security Command Mode Request to UE
      NGAP_INITIAL_CONTEXT_SETUP_REQ,  // Initial Conext Setup Request to UE &
                                       // Registration Accept
      NGAP_NAS_DL_DATA_REQ,            // Deregistaration Accept
      NGAP_UE_CONTEXT_RELEASE_COMMAND  // UEContextReleaseCommand
  };

  /* Send the initial UE message */
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));

  /* Check if UE Context is created with correct imsi */
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id));

  /* Send the authentication response message from subscriberdb */
  rc = send_proc_authentication_info_answer(imsi, ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for auth response from UE */
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for security mode complete response from UE */
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  send_initial_context_response(amf_app_desc_p, ue_id);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_TRUE(rc == RETURNok);
  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, TestRegistrationProcGutiBased) {
  amf_ue_ngap_id_t ue_id = 0;
  std::vector<MessagesIds> expected_Ids{
      AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,  // new registration notification
                                            // indication to ngap
      NGAP_NAS_DL_DATA_REQ,                 // Authentication Request to UE
      NGAP_NAS_DL_DATA_REQ,  // Security Command Mode Request to UE
      NGAP_NAS_DL_DATA_REQ,
      NGAP_INITIAL_CONTEXT_SETUP_REQ,  // Initial Conext Setup Request to UE &
                                       // Registration Accept
      NGAP_UE_CONTEXT_RELEASE_COMMAND  // UEContextReleaseCommand
  };

  /* Send the initial UE message */
  imsi64_t imsi64 = 0;
  uint32_t m_tmsi = 0x12704177;
  ue_id = send_initial_ue_message_with_tmsi(
      amf_app_desc_p, 36, 1, 1, 0, plmn, m_tmsi, guti_initial_ue_message_hexbuf,
      sizeof(guti_initial_ue_message_hexbuf));

  EXPECT_TRUE(validate_identification_procedure(ue_id, 0));

  int rc = RETURNok;
  rc = send_uplink_nas_identity_response_message(amf_app_desc_p, ue_id, plmn,
                                                 identity_response,
                                                 sizeof(identity_response));
  EXPECT_TRUE(rc == RETURNok);

  ue_m5gmm_context_t* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  ASSERT_NE(ue_context_p, nullptr);
  EXPECT_TRUE(ue_context_p->amf_context.imsi64 == stoul(imsi));

  /* Send the authentication response message from subscriberdb */
  rc = send_proc_authentication_info_answer(imsi, ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  /* Validate if authentication procedure is initialized as expected */
  EXPECT_TRUE(validate_auth_procedure(ue_id, 0));

  /* Send uplink nas message for auth response from UE */
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Check whether security mode procedure is initiated */
  EXPECT_TRUE(validate_smc_procedure(ue_id, 0));

  /* Send uplink nas message for security mode complete response from UE */
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  amf_app_handle_deregistration_req(ue_id);
  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, TestMobileUpdatingRegistrationProcGutiBased) {
  amf_ue_ngap_id_t ue_id = 0;
  std::vector<MessagesIds> expected_Ids{
      AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,  // new registration notification
                                            // indication to ngap
      NGAP_NAS_DL_DATA_REQ,                 // Authentication Request to UE
      NGAP_NAS_DL_DATA_REQ,  // Security Command Mode Request to UE
      NGAP_NAS_DL_DATA_REQ,
      NGAP_INITIAL_CONTEXT_SETUP_REQ,  // Initial Conext Setup Request to UE &
                                       // Registration Accept
      NGAP_UE_CONTEXT_RELEASE_COMMAND  // UEContextReleaseCommand
  };

  /* Send the initial UE message */
  imsi64_t imsi64 = 0;
  uint32_t m_tmsi = 0x8f719f0e;
  ue_id = send_initial_ue_message_with_tmsi(
      amf_app_desc_p, 36, 1, 1, 0, plmn, m_tmsi, mu_initial_ue_message_hexbuf,
      sizeof(mu_initial_ue_message_hexbuf));

  EXPECT_TRUE(validate_identification_procedure(ue_id, 0));

  int rc = RETURNok;
  rc = send_uplink_nas_identity_response_message(amf_app_desc_p, ue_id, plmn,
                                                 identity_response,
                                                 sizeof(identity_response));
  EXPECT_TRUE(rc == RETURNok);

  ue_m5gmm_context_t* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  ASSERT_NE(ue_context_p, nullptr);
  EXPECT_TRUE(ue_context_p->amf_context.imsi64 == stoul(imsi));

  /* Send the authentication response message from subscriberdb */
  rc = send_proc_authentication_info_answer(imsi, ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  /* Validate if authentication procedure is initialized as expected */
  EXPECT_TRUE(validate_auth_procedure(ue_id, 0));

  /* Send uplink nas message for auth response from UE */
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Check whether security mode procedure is initiated */
  EXPECT_TRUE(validate_smc_procedure(ue_id, 0));

  /* Send uplink nas message for security mode complete response from UE */
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  amf_app_handle_deregistration_req(ue_id);
  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, TestPDUSessionSetup) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id = 0;
  std::vector<MessagesIds> expected_Ids{
      AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,  // new registration notification
                                            // indication to ngap
      NGAP_NAS_DL_DATA_REQ,                 // Authentication Request to UE
      NGAP_NAS_DL_DATA_REQ,            // Security Command Mode Request to UE
      NGAP_INITIAL_CONTEXT_SETUP_REQ,  // Initial Conext Setup Request to UE &
                                       // Registration Accept
      NGAP_PDUSESSION_RESOURCE_SETUP_REQ,  // PDU Resource Setup Request to GNB
                                           // & PDU Session Establishment Accept
      NGAP_PDUSESSIONRESOURCE_REL_REQ,  // PDU Session Resource Release Command
      NGAP_NAS_DL_DATA_REQ,             // Deregistaration Accept
      NGAP_UE_CONTEXT_RELEASE_COMMAND   // UEContextReleaseCommand
  };

  /* Send the initial UE message */
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));

  /* Check if UE Context is created with correct imsi */
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id));

  /* Send the authentication response message from subscriberdb */
  rc = send_proc_authentication_info_answer(imsi, ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for auth response from UE */
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for security mode complete response from UE */
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  send_initial_context_response(amf_app_desc_p, ue_id);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for pdu session establishment request from UE */
  rc = send_uplink_nas_pdu_session_establishment_request(
      amf_app_desc_p, ue_id, plmn, ue_pdu_session_est_req_hexbuf,
      sizeof(ue_pdu_session_est_req_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send ip address response  from pipelined */
  rc = send_ip_address_response_itti(IPv4);
  EXPECT_TRUE(rc == RETURNok);

  /* Send pdu session setup response  from smf */
  rc = send_pdu_session_response_itti(IPv4);
  EXPECT_TRUE(rc == RETURNok);

  /* Send pdu resource setup response  from UE */
  rc = send_pdu_resource_setup_response(ue_id);
  EXPECT_TRUE(rc == RETURNok);

  rc = send_pdu_notification_response();
  EXPECT_TRUE(rc == RETURNok);

  ue_m5gmm_context_t* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  ASSERT_NE(ue_context_p, nullptr);

  EXPECT_EQ(ue_context_p->amf_context.smf_ctxt_map.size(), 1);
  std::shared_ptr<smf_context_t> smf_ctxt =
      ue_context_p->amf_context.smf_ctxt_map.begin()->second;

  /* Send uplink nas message for pdu session release request from UE */
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, ue_id, plmn, pdu_sess_release_hexbuf,
      sizeof(pdu_sess_release_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for pdu session release complete from UE */
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, ue_id, plmn, pdu_sess_release_complete_hexbuf,
      sizeof(pdu_sess_release_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  EXPECT_EQ(ue_context_p->amf_context.smf_ctxt_map.size(), 0);

  rc = send_pdu_notification_response();
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for deregistration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_TRUE(rc == RETURNok);
  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, TestPDUSessionSetupNoDnnNoSlice) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id = 0;
  std::vector<MessagesIds> expected_Ids{
      AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,  // new registration notification
                                            // indication to ngap
      NGAP_NAS_DL_DATA_REQ,                 // Authentication Request to UE
      NGAP_NAS_DL_DATA_REQ,            // Security Command Mode Request to UE
      NGAP_INITIAL_CONTEXT_SETUP_REQ,  // Initial Conext Setup Request to UE &
                                       // Registration Accept
      NGAP_PDUSESSION_RESOURCE_SETUP_REQ,  // PDU Resource Setup Request to GNB
                                           // & PDU Session Establishment Accept
      NGAP_PDUSESSIONRESOURCE_REL_REQ,  // PDU Session Resource Release Command
      NGAP_NAS_DL_DATA_REQ,             // Deregistaration Accept
      NGAP_UE_CONTEXT_RELEASE_COMMAND   // UEContextReleaseCommand
  };

  /* Send the initial UE message */
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));

  /* Check if UE Context is created with correct imsi */
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id));

  /* Send the authentication response message from subscriberdb */
  rc = send_proc_authentication_info_answer(imsi, ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for auth response from UE */
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for security mode complete response from UE */
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  send_initial_context_response(amf_app_desc_p, ue_id);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for pdu session establishment request from UE */
  rc = send_uplink_nas_pdu_session_establishment_request(
      amf_app_desc_p, ue_id, plmn,
      ue_pdu_session_est_req_no_dnn_no_slice_hexbuf,
      sizeof(ue_pdu_session_est_req_no_dnn_no_slice_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send ip address response  from pipelined */
  rc = send_ip_address_response_itti(IPv4);
  EXPECT_TRUE(rc == RETURNok);

  /* Send pdu session setup response  from smf */
  rc = send_pdu_session_response_itti(IPv4);
  EXPECT_TRUE(rc == RETURNok);

  /* Send pdu resource setup response  from UE */
  rc = send_pdu_resource_setup_response(ue_id);
  EXPECT_TRUE(rc == RETURNok);

  rc = send_pdu_notification_response();
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for pdu session release request from UE */
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, ue_id, plmn, pdu_sess_release_hexbuf,
      sizeof(pdu_sess_release_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for pdu session release complete from UE */
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, ue_id, plmn, pdu_sess_release_complete_hexbuf,
      sizeof(pdu_sess_release_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  ue_m5gmm_context_t* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  ASSERT_NE(ue_context_p, nullptr);
  EXPECT_EQ(ue_context_p->amf_context.smf_ctxt_map.size(), 0);

  rc = send_pdu_notification_response();
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for deregistration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_TRUE(rc == RETURNok);
  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, TestPDUSessionFailure_dnn_not_subscribed) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id = 0;
  std::vector<MessagesIds> expected_Ids{
      AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,  // new registration notification
                                            // indication to ngap
      NGAP_NAS_DL_DATA_REQ,                 // Authentication Request to UE
      NGAP_NAS_DL_DATA_REQ,            // Security Command Mode Request to UE
      NGAP_INITIAL_CONTEXT_SETUP_REQ,  // Initial Conext Setup Request to UE &
                                       // Registration Accept
      NGAP_NAS_DL_DATA_REQ,  // PDU Session Establishment Request with failure
                             // mm cause
      NGAP_NAS_DL_DATA_REQ,  // Deregistaration Accept
      NGAP_UE_CONTEXT_RELEASE_COMMAND  // UEContextReleaseCommand
  };

  /* Send the initial UE message */
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));

  /* Check if UE Context is created with correct imsi */
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id));

  /* Send the authentication response message from subscriberdb */
  rc = send_proc_authentication_info_answer(imsi, ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for auth response from UE */
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for security mode complete response from UE */
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  send_initial_context_response(amf_app_desc_p, ue_id);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for pdu session establishment request from UE */
  rc = send_uplink_nas_pdu_session_establishment_request(
      amf_app_desc_p, ue_id, plmn,
      ue_pdu_session_est_req_dnn_not_subscried_hexbuf,
      sizeof(ue_pdu_session_est_req_dnn_not_subscried_hexbuf));
  EXPECT_EQ(rc, RETURNok);
  ue_m5gmm_context_t* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  // ue context should exist
  ASSERT_NE(ue_context_p, nullptr);
  // smf context should not be present
  EXPECT_EQ(ue_context_p->amf_context.smf_ctxt_map.size(), 0);

  /* Send uplink nas message for deregistration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_TRUE(rc == RETURNok);
  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, TestPDUSession_missing_dnn) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id = 0;
  std::vector<MessagesIds> expected_Ids{
      AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,  // new registration notification
                                            // indication to ngap
      NGAP_NAS_DL_DATA_REQ,                 // Authentication Request to UE
      NGAP_NAS_DL_DATA_REQ,            // Security Command Mode Request to UE
      NGAP_INITIAL_CONTEXT_SETUP_REQ,  // Initial Conext Setup Request to UE &
                                       // Registration Accept
      NGAP_NAS_DL_DATA_REQ,  // PDU Session Establishment Request with failure
                             // sm cause
      NGAP_NAS_DL_DATA_REQ,  // Deregistaration Accept
      NGAP_UE_CONTEXT_RELEASE_COMMAND  // UEContextReleaseCommand
  };

  /* Send the initial UE message */
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));

  /* Check if UE Context is created with correct imsi */
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id));

  /* Send the authentication response message from subscriberdb */
  rc = send_proc_authentication_info_answer(imsi, ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for auth response from UE */
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for security mode complete response from UE */
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  send_initial_context_response(amf_app_desc_p, ue_id);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  ue_m5gmm_context_t* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  // ue context should exist
  ASSERT_NE(ue_context_p, nullptr);
  memset(&ue_context_p->amf_context.apn_config_profile, 0,
         sizeof(ue_context_p->amf_context.apn_config_profile));

  /* Send uplink nas message for pdu session establishment request from UE */
  rc = send_uplink_nas_pdu_session_establishment_request(
      amf_app_desc_p, ue_id, plmn, ue_pdu_session_est_req_missing_dnn_hexbuf,
      sizeof(ue_pdu_session_est_req_missing_dnn_hexbuf));
  EXPECT_EQ(rc, RETURNok);

  ue_context_p = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  // ue context should exist
  ASSERT_NE(ue_context_p, nullptr);
  // smf context should be present
  EXPECT_EQ(ue_context_p->amf_context.smf_ctxt_map.size(), 0);

  /* Send uplink nas message for deregistration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_TRUE(rc == RETURNok);
  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, TestPDUSession_unknown_pdu_session_type) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id = 0;
  std::vector<MessagesIds> expected_Ids{
      AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,  // new registration notification
                                            // indication to ngap
      NGAP_NAS_DL_DATA_REQ,                 // Authentication Request to UE
      NGAP_NAS_DL_DATA_REQ,            // Security Command Mode Request to UE
      NGAP_INITIAL_CONTEXT_SETUP_REQ,  // Initial Conext Setup Request to UE &
                                       // Registration Accept
      NGAP_NAS_DL_DATA_REQ,  // PDU Session Establishment Request with failure
                             // sm cause
      NGAP_NAS_DL_DATA_REQ,  // Deregistaration Accept
      NGAP_UE_CONTEXT_RELEASE_COMMAND  // UEContextReleaseCommand
  };

  /* Send the initial UE message */
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));

  /* Check if UE Context is created with correct imsi */
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id));

  /* Send the authentication response message from subscriberdb */
  rc = send_proc_authentication_info_answer(imsi, ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for auth response from UE */
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for security mode complete response from UE */
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  send_initial_context_response(amf_app_desc_p, ue_id);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for pdu session establishment request from UE */
  rc = send_uplink_nas_pdu_session_establishment_request(
      amf_app_desc_p, ue_id, plmn,
      ue_pdu_session_est_req_unknown_session_type_hexbuf,
      sizeof(ue_pdu_session_est_req_unknown_session_type_hexbuf));
  EXPECT_EQ(rc, RETURNok);

  ue_m5gmm_context_t* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  // ue context should exist
  ASSERT_NE(ue_context_p, nullptr);
  // smf context should be present
  EXPECT_EQ(ue_context_p->amf_context.smf_ctxt_map.size(), 0);

  /* Send uplink nas message for deregistration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_TRUE(rc == RETURNok);
  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, TestPDUSession_Invalid_PDUSession_Identity) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id = 0;
  std::vector<MessagesIds> expected_Ids{
      AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,  // new registration notification
                                            // indication to ngap
      NGAP_NAS_DL_DATA_REQ,                 // Authentication Request to UE
      NGAP_NAS_DL_DATA_REQ,            // Security Command Mode Request to UE
      NGAP_INITIAL_CONTEXT_SETUP_REQ,  // Initial Conext Setup Request to UE &
                                       // Registration Accept
      NGAP_PDUSESSION_RESOURCE_SETUP_REQ,  // PDU Resource Setup Request to GNB
                                           // & PDU Session Establishment Accept
      NGAP_NAS_DL_DATA_REQ,  // PDU Session Establishment Request with failure
                             // sm cause
      NGAP_PDUSESSIONRESOURCE_REL_REQ,  // PDU Session Resource Release Command
      NGAP_NAS_DL_DATA_REQ,             // Deregistaration Accept
      NGAP_UE_CONTEXT_RELEASE_COMMAND   // UEContextReleaseCommand
  };

  /* Send the initial UE message */
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));

  /* Check if UE Context is created with correct imsi */
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id));

  ue_m5gmm_context_t* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  // ue context should exist
  ASSERT_NE(ue_context_p, nullptr);
  nas5g_auth_info_proc_t* auth_info_proc =
      get_nas5g_cn_procedure_auth_info(&ue_context_p->amf_context);

  ASSERT_NE(auth_info_proc, nullptr);

  /* Send the authentication response message from subscriberdb */
  rc = send_proc_authentication_info_answer(imsi, ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for auth response from UE */
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for security mode complete response from UE */
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  send_initial_context_response(amf_app_desc_p, ue_id);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for pdu session establishment request from UE */
  rc = send_uplink_nas_pdu_session_establishment_request(
      amf_app_desc_p, ue_id, plmn, ue_pdu_session_est_req_hexbuf,
      sizeof(ue_pdu_session_est_req_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send ip address response  from pipelined */
  rc = send_ip_address_response_itti(IPv4);
  EXPECT_TRUE(rc == RETURNok);

  /* Send pdu session setup response  from smf */
  rc = send_pdu_session_response_itti(IPv4);
  EXPECT_TRUE(rc == RETURNok);

  /* Send pdu resource setup response  from UE */
  rc = send_pdu_resource_setup_response(ue_id);
  EXPECT_TRUE(rc == RETURNok);

  rc = send_pdu_notification_response();
  EXPECT_TRUE(rc == RETURNok);

  ue_context_p = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  // ue context should exist
  ASSERT_NE(ue_context_p, nullptr);
  // smf context should be present
  EXPECT_EQ(ue_context_p->amf_context.smf_ctxt_map.size(), 1);
  std::shared_ptr<smf_context_t> smf_ctx =
      amf_get_smf_context_by_pdu_session_id(ue_context_p, 1);
  EXPECT_EQ(smf_ctx->duplicate_pdu_session_est_req_count, 0);
  EXPECT_EQ(smf_ctx->pdu_session_state, ACTIVE);

  for (uint8_t i = 1; i < 5; ++i) {
    /* Send duplicate uplink nas message for pdu session establishment request
     * from UE */
    rc = send_uplink_nas_pdu_session_establishment_request(
        amf_app_desc_p, ue_id, plmn, ue_pdu_session_est_req_hexbuf,
        sizeof(ue_pdu_session_est_req_hexbuf));
    EXPECT_EQ(rc, RETURNok);

    // smf context should be present
    EXPECT_EQ(ue_context_p->amf_context.smf_ctxt_map.size(), 1);
    smf_ctx = amf_get_smf_context_by_pdu_session_id(ue_context_p, 1);
    EXPECT_EQ(smf_ctx->duplicate_pdu_session_est_req_count, i);
  }

  /* Send duplicate uplink nas message for pdu session establishment request
   * from UE */
  rc = send_uplink_nas_pdu_session_establishment_request(
      amf_app_desc_p, ue_id, plmn, ue_pdu_session_est_req_hexbuf,
      sizeof(ue_pdu_session_est_req_hexbuf));
  EXPECT_EQ(rc, RETURNok);

  // smf context should be present
  EXPECT_EQ(ue_context_p->amf_context.smf_ctxt_map.size(), 1);
  smf_ctx = amf_get_smf_context_by_pdu_session_id(ue_context_p, 1);
  EXPECT_EQ(smf_ctx->duplicate_pdu_session_est_req_count, 4);

  /* Send uplink nas message for pdu session release request from UE */
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, ue_id, plmn, pdu_sess_release_hexbuf,
      sizeof(pdu_sess_release_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for pdu session release complete from UE */
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, ue_id, plmn, pdu_sess_release_complete_hexbuf,
      sizeof(pdu_sess_release_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  rc = send_pdu_notification_response();
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for deregistration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_TRUE(rc == RETURNok);
  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

// TODO: #11034 is the starting point of TestRegistrationProcSUCIExt failure.
// This issue is tracked here #11382
#if 0
TEST_F(AMFAppProcedureTest, TestRegistrationProcSUCIExt) {
  amf_ue_ngap_id_t ue_id = 0;

  // Send the initial UE message
  imsi64_t imsi64 = 0;
  int rc          = RETURNok;
  imsi64          = send_initial_ue_message_no_tmsi(
      amf_app_desc_p, 36, 1, 1, 0, plmn, intital_ue_message_suci_ext_hexbuf,
      sizeof(intital_ue_message_suci_ext_hexbuf));

  rc = amf_decrypt_imsi_info_answer(&decrypted_imsi);
  EXPECT_TRUE(rc == RETURNok);

  // Check if UE Context is created with correct imsi
  bool res = false;
  res      = get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id);
  EXPECT_TRUE(res == true);

  // Send the authentication response message from subscriberdb
  rc = send_proc_authentication_info_answer(imsi, ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  // Validate if authentication procedure is initialized as expected
  res = validate_auth_procedure(ue_id, 0);
  EXPECT_TRUE(res == true);

  // Send uplink nas message for auth response from UE
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Check whether security mode procedure is initiated
  res = validate_smc_procedure(ue_id, 0);
  EXPECT_TRUE(res == true);

  // Send uplink nas message for security mode complete response from UE
  rc = send_uplink_nas_message_ue_smc_response(
      amf_app_desc_p, ue_id, plmn, ue_smc_response_hexbuf,
      sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for registration complete response from UE
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  amf_app_handle_deregistration_req(ue_id);
}
#endif

TEST_F(AMFAppProcedureTest, TestAuthFailureFromSubscribeDb) {
  amf_ue_ngap_id_t ue_id = 0;
  ue_m5gmm_context_s* ue_5gmm_context_p = NULL;
  std::vector<MessagesIds> expected_Ids{
      AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,  // new registration notification
                                            // indication to ngap
      NGAP_NAS_DL_DATA_REQ,                 // Registration Reject
      NGAP_UE_CONTEXT_RELEASE_COMMAND       // UEContextReleaseCommand
  };

  /* Send the initial UE message */
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));

  /* Check if UE Context is created with correct imsi */
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id));

  /* Send the authentication response message from subscriberdb */
  int rc = RETURNok;
  rc = send_proc_authentication_info_answer(imsi, ue_id, false);
  EXPECT_TRUE(rc == RETURNok);
  ue_5gmm_context_p = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  EXPECT_TRUE(ue_5gmm_context_p == NULL);
}

TEST(test_t3592abort, test_pdu_session_release_notify_smf) {
  amf_ue_ngap_id_t ue_id = 1;
  uint8_t pdu_session_id = 1;
  int rc = RETURNerror;
  // creating ue_context
  ue_m5gmm_context_s* ue_context = amf_create_new_ue_context();
  ue_context_map.insert(
      std::pair<amf_ue_ngap_id_t, ue_m5gmm_context_s*>(ue_id, ue_context));
  std::shared_ptr<smf_context_t> smf_ctx =
      amf_insert_smf_context(ue_context, pdu_session_id);
  smf_ctx->pdu_session_state = ACTIVE;
  ue_context->mm_state = REGISTERED_CONNECTED;
  smf_ctx->n_active_pdus = 1;
  EXPECT_NE(ue_context, nullptr);
  EXPECT_NE(ue_context->amf_context.smf_ctxt_map.size(), 0);
  rc = t3592_abort_handler(ue_context, smf_ctx, pdu_session_id);
  EXPECT_TRUE(rc == RETURNok);
  EXPECT_EQ(ue_context->amf_context.smf_ctxt_map.size(), 0);
  ue_context_map.clear();
  delete ue_context;
}

TEST_F(AMFAppProcedureTest, TestPDUv6SessionSetup) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id = 0;
  std::vector<MessagesIds> expected_Ids{
      AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,  // new registration notification
                                            // indication to ngap
      NGAP_NAS_DL_DATA_REQ,                 // Authentication Request to UE
      NGAP_NAS_DL_DATA_REQ,            // Security Command Mode Request to UE
      NGAP_INITIAL_CONTEXT_SETUP_REQ,  // Initial Conext Setup Request to UE &
                                       // Registration Accept
      NGAP_PDUSESSION_RESOURCE_SETUP_REQ,  // PDU Resource Setup Request to GNB
                                           // & PDU Session Establishment Accept
      NGAP_PDUSESSIONRESOURCE_REL_REQ,  // PDU Session Resource Release Command
      NGAP_NAS_DL_DATA_REQ,             // Deregistaration Accept
      NGAP_UE_CONTEXT_RELEASE_COMMAND   // UEContextReleaseCommand
  };

  /* Send the initial UE message */
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));

  /* Check if UE Context is created with correct imsi */
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id));

  /* Send the authentication response message from subscriberdb */
  rc = send_proc_authentication_info_answer(imsi, ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for auth response from UE */
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for security mode complete response from UE */
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  send_initial_context_response(amf_app_desc_p, ue_id);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for pdu session establishment request from UE */
  rc = send_uplink_nas_pdu_session_establishment_request(
      amf_app_desc_p, ue_id, plmn, ue_pdu_session_v6_est_req_hexbuf,
      sizeof(ue_pdu_session_est_req_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send ip address response  from pipelined */
  rc = send_ip_address_response_itti(IPv6);
  EXPECT_TRUE(rc == RETURNok);

  /* Send pdu session setup response  from smf */
  rc = send_pdu_session_response_itti(IPv6);
  EXPECT_TRUE(rc == RETURNok);

  /* Send pdu resource setup response  from UE */
  rc = send_pdu_resource_setup_response(ue_id);
  EXPECT_TRUE(rc == RETURNok);

  rc = send_pdu_notification_response();
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for pdu session release request from UE */
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, ue_id, plmn, pdu_sess_release_hexbuf,
      sizeof(pdu_sess_release_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for pdu session release complete from UE */
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, ue_id, plmn, pdu_sess_release_complete_hexbuf,
      sizeof(pdu_sess_release_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  ue_m5gmm_context_t* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  ASSERT_NE(ue_context_p, nullptr);
  EXPECT_EQ(ue_context_p->amf_context.smf_ctxt_map.size(), 0);

  rc = send_pdu_notification_response();
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for deregistration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_TRUE(rc == RETURNok);
  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, TestPDUv4v6SessionSetup) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id = 0;
  std::vector<MessagesIds> expected_Ids{
      AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,  // new registration notification
                                            // indication to ngap
      NGAP_NAS_DL_DATA_REQ,                 // Authentication Request to UE
      NGAP_NAS_DL_DATA_REQ,            // Security Command Mode Request to UE
      NGAP_INITIAL_CONTEXT_SETUP_REQ,  // Initial Conext Setup Request to UE &
                                       // Registration Accept
      NGAP_PDUSESSION_RESOURCE_SETUP_REQ,  // PDU Resource Setup Request to GNB
                                           // & PDU Session Establishment Accept
      NGAP_PDUSESSIONRESOURCE_REL_REQ,  // PDU Session Resource Release Command
      NGAP_NAS_DL_DATA_REQ,             // Deregistaration Accept
      NGAP_UE_CONTEXT_RELEASE_COMMAND   // UEContextReleaseCommand
  };

  /* Send the initial UE message */
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));

  /* Check if UE Context is created with correct imsi */
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id));

  /* Send the authentication response message from subscriberdb */
  rc = send_proc_authentication_info_answer(imsi, ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for auth response from UE */
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for security mode complete response from UE */
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  send_initial_context_response(amf_app_desc_p, ue_id);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for pdu session establishment request from UE */
  rc = send_uplink_nas_pdu_session_establishment_request(
      amf_app_desc_p, ue_id, plmn, ue_pdu_session_v4v6_est_req_hexbuf,
      sizeof(ue_pdu_session_est_req_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send ip address response  from pipelined */
  rc = send_ip_address_response_itti(IPv4_AND_v6);
  EXPECT_TRUE(rc == RETURNok);

  /* Send pdu session setup response  from smf */
  rc = send_pdu_session_response_itti(IPv4_AND_v6);
  EXPECT_TRUE(rc == RETURNok);

  /* Send pdu resource setup response  from UE */
  rc = send_pdu_resource_setup_response(ue_id);
  EXPECT_TRUE(rc == RETURNok);

  rc = send_pdu_notification_response();
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for pdu session release request from UE */
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, ue_id, plmn, pdu_sess_release_hexbuf,
      sizeof(pdu_sess_release_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for pdu session release complete from UE */
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, ue_id, plmn, pdu_sess_release_complete_hexbuf,
      sizeof(pdu_sess_release_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  ue_m5gmm_context_t* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  ASSERT_NE(ue_context_p, nullptr);
  EXPECT_EQ(ue_context_p->amf_context.smf_ctxt_map.size(), 0);

  rc = send_pdu_notification_response();
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for deregistration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_TRUE(rc == RETURNok);
  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, ServiceRequestMTWithPDU) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id = 0;
  std::vector<MessagesIds> expected_Ids{AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_INITIAL_CONTEXT_SETUP_REQ,
                                        NGAP_PDUSESSION_RESOURCE_SETUP_REQ,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND,
                                        AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,
                                        AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,
                                        NGAP_INITIAL_CONTEXT_SETUP_REQ,
                                        NGAP_PDUSESSIONRESOURCE_REL_REQ,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND};

  /* Send the initial UE message */
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));

  /* Check if UE Context is created with correct imsi */
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id));

  /* Send the authentication response message from subscriberdb */
  rc = send_proc_authentication_info_answer(imsi, ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for auth response from UE */
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for security mode complete response from UE */
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  send_initial_context_response(amf_app_desc_p, ue_id);

  /* Send uplink nas message for pdu session establishment request from UE */
  rc = send_uplink_nas_pdu_session_establishment_request(
      amf_app_desc_p, ue_id, plmn, ue_pdu_session_est_req_hexbuf,
      sizeof(ue_pdu_session_est_req_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send ip address response  from pipelined */
  rc = send_ip_address_response_itti(IPv4);
  EXPECT_TRUE(rc == RETURNok);

  /* Send pdu session setup response  from smf */
  rc = send_pdu_session_response_itti(IPv4);
  EXPECT_TRUE(rc == RETURNok);

  /* Send pdu resource setup response  from UE */
  rc = send_pdu_resource_setup_response(ue_id);
  EXPECT_TRUE(rc == RETURNok);

  rc = send_pdu_notification_response();
  EXPECT_TRUE(rc == RETURNok);

  /*Send UE context release request to move to idle mode*/
  send_ue_context_release_request_message(amf_app_desc_p, 1, 1, ue_id);
  ue_m5gmm_context_s* ue_context = nullptr;
  ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  EXPECT_EQ(REGISTERED_IDLE, ue_context->mm_state);

  /*Send initial UE message Service Request with PDU*/
  imsi64 = 0;
  imsi64 = send_initial_ue_message_service_request(
      amf_app_desc_p, 36, 1, 2, ue_id, plmn,
      initial_ue_msg_service_request_with_pdu,
      sizeof(initial_ue_msg_service_request_with_pdu), 8);

  EXPECT_EQ(REGISTERED_IDLE, ue_context->mm_state);
  EXPECT_EQ(true, ue_context->pending_service_response);

  // Ue id will be freshly generated
  EXPECT_NE(ue_context->amf_ue_ngap_id, ue_id);

  // Update the ue_id
  ue_id = ue_context->amf_ue_ngap_id;

  /* Send pdu session setup response  from smf */
  rc = send_pdu_session_response_itti(IPv4);
  EXPECT_TRUE(rc == RETURNok);
  EXPECT_EQ(REGISTERED_CONNECTED, ue_context->mm_state);
  EXPECT_EQ(false, ue_context->pending_service_response);

  send_initial_context_response(amf_app_desc_p, ue_id);

  /* Send uplink nas message for pdu session release request from UE */
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, ue_id, plmn, pdu_sess_release_hexbuf,
      sizeof(pdu_sess_release_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for pdu session release complete from UE */
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, ue_id, plmn, pdu_sess_release_complete_hexbuf,
      sizeof(pdu_sess_release_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  rc = send_pdu_notification_response();
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for deregistration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_TRUE(rc == RETURNok);

  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, ServiceRequestSignalling) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id = 0;
  std::vector<MessagesIds> expected_Ids{AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_INITIAL_CONTEXT_SETUP_REQ,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND,
                                        AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,
                                        AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,
                                        NGAP_INITIAL_CONTEXT_SETUP_REQ,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND};

  /* Send the initial UE message */
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));
  /* Check if UE Context is created with correct imsi */
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id));

  /* Send the authentication response message from subscriberdb */
  rc = send_proc_authentication_info_answer(imsi, ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for auth response from UE */
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for security mode complete response from UE */
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  send_initial_context_response(amf_app_desc_p, ue_id);

  /*Send UE context release request to move to idle mode*/
  send_ue_context_release_request_message(amf_app_desc_p, 1, 1, ue_id);
  ue_m5gmm_context_s* ue_context = nullptr;
  ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  EXPECT_EQ(REGISTERED_IDLE, ue_context->mm_state);

  /*Send initial UE message Service Request with PDU*/
  imsi64 = 0;
  imsi64 = send_initial_ue_message_service_request(
      amf_app_desc_p, 36, 1, 2, ue_id, plmn,
      initial_ue_msg_service_request_signaling,
      sizeof(initial_ue_msg_service_request_signaling), 4);

  EXPECT_EQ(REGISTERED_CONNECTED, ue_context->mm_state);

  // Ue id will be freshly generated
  EXPECT_NE(ue_context->amf_ue_ngap_id, ue_id);

  // Update the ue_id
  ue_id = ue_context->amf_ue_ngap_id;

  send_initial_context_response(amf_app_desc_p, ue_id);

  /* Send uplink nas message for deregistration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_TRUE(rc == RETURNok);

  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, TestAuthFailureFromSubscribeDbLock) {
  amf_ue_ngap_id_t ue_id = 0;
  amf_context_t* amf_ctxt_p = nullptr;
  nas5g_auth_info_proc_t* auth_info_proc = nullptr;
  ue_m5gmm_context_s* ue_context_p = nullptr;
  std::vector<MessagesIds> expected_Ids{
      AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,  // new registration notification
                                            // indication to ngap
      NGAP_NAS_DL_DATA_REQ,                 // Registration Reject
      NGAP_UE_CONTEXT_RELEASE_COMMAND       // UEContextReleaseCommand
  };
  /* Send the initial UE message */
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));
  /* Check if UE Context is created with correct imsi */
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id));

  /* Send the authentication response message from subscriberdb */
  itti_amf_subs_auth_info_ans_t aia_itti_msg = {};
  strncpy(aia_itti_msg.imsi, imsi.c_str(), imsi.size());
  aia_itti_msg.imsi_length = imsi.size();
  aia_itti_msg.result = DIAMETER_TOO_BUSY;
  int rc = RETURNerror;
  rc = amf_nas_proc_authentication_info_answer(&aia_itti_msg);
  EXPECT_TRUE(rc == RETURNok);
  ue_context_p = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  EXPECT_NE(ue_context_p, nullptr);
  amf_ctxt_p = &ue_context_p->amf_context;
  EXPECT_NE(amf_ctxt_p, nullptr);
  auth_info_proc = get_nas5g_cn_procedure_auth_info(amf_ctxt_p);
  EXPECT_NE(auth_info_proc, nullptr);
  nas5g_delete_cn_procedure(amf_ctxt_p, &auth_info_proc->cn_proc);
  amf_free_ue_context(ue_context_p);
}

}  // namespace magma5g
