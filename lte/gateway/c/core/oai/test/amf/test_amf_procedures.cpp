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

#include "lte/gateway/c/core/oai/test/mock_tasks/mock_tasks.hpp"

extern "C" {
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/include/amf_config.hpp"
#include "lte/gateway/c/core/oai/include/amf_app_messages_types.h"
}
#include "lte/gateway/c/core/oai/tasks/amf/include/amf_client_servicer.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_state_manager.hpp"
#include "lte/gateway/c/core/oai/test/amf/amf_app_test_util.h"
#include "lte/gateway/c/core/oai/test/amf/util_s6a_update_location.hpp"

using ::testing::Test;

namespace magma5g {

extern task_zmq_ctx_s amf_app_task_zmq_ctx;

class AMFAppProcedureTest : public ::testing::Test {
  virtual void SetUp() {
    itti_init(TASK_MAX, THREAD_MAX, MESSAGES_ID_MAX, tasks_info, messages_info,
              NULL, NULL);

    amf_config_init(&amf_config);
    amf_config.guamfi.guamfi[0].plmn = plmn;
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

  itti_amf_decrypted_msin_info_ans_t decrypted_msin = {
      .msin = {'9', '7', '6', '5', '4', '5', '6', '6', '0'},
      .msin_length = 10,
      .result = 1,
      .ue_id = 1};

  std::string decrypted_imsi = "222456976545660";
  const uint8_t intital_ue_message_suci_ext_hexbuf[67] = {
      0x7e, 0x00, 0x41, 0x79, 0x00, 0x39, 0x01, 0x22, 0x62, 0x54, 0xf0, 0xff,
      0x01, 0x05, 0x25, 0xb6, 0xb6, 0xdf, 0x89, 0xaf, 0x58, 0xb0, 0xe7, 0x07,
      0x87, 0xfe, 0x52, 0x77, 0xa6, 0x31, 0x7c, 0x2c, 0xc4, 0x7d, 0x76, 0x4a,
      0x81, 0xaa, 0x3e, 0xcc, 0xbe, 0xa3, 0x7b, 0xd0, 0x57, 0x40, 0xae, 0xe0,
      0xd5, 0x54, 0x70, 0xbf, 0xf4, 0x7c, 0x08, 0xe3, 0x1d, 0xf9, 0xb8, 0x55,
      0x99, 0x12, 0x48, 0x2e, 0x02, 0xf0, 0xf0};

  const uint8_t initial_ue_message_hexbuf[29] = {
      0x7e, 0x00, 0x41, 0x79, 0x00, 0x0d, 0x01, 0x22, 0x62, 0x54,
      0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x2e,
      0x08, 0x80, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00};

  const uint8_t initial_ue_plmn_mismatch_message_hexbuf[23] = {
      0x7e, 0x00, 0x41, 0x79, 0x00, 0x0d, 0x01, 0x00, 0x91, 0x10, 0xf0, 0xff,
      0x00, 0x00, 0x79, 0x56, 0x54, 0x66, 0xf4, 0x2e, 0x02, 0xf0, 0xf0};

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

  const uint8_t guti_initial_ue_reregister_message_hexbuf[36] = {
      0x7e, 0x00, 0x41, 0x01, 0x00, 0x0b, 0xf2, 0x22, 0x62, 0x54, 0x01, 0x00,
      0x40, 0xd9, 0x8a, 0x4a, 0x7d, 0x10, 0x01, 0x00, 0x2e, 0x02, 0xc0, 0xc0,
      0x2f, 0x02, 0x01, 0x02, 0x17, 0x02, 0xc0, 0xc0, 0xb0, 0x2b, 0x01, 0x00};

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

  const uint8_t identity_response_profile_a[59] = {
      0x7e, 0x00, 0x5c, 0x00, 0x35, 0x01, 0x22, 0x62, 0x54, 0xf0, 0xff, 0x01,
      0x01, 0xa1, 0xae, 0x7e, 0x56, 0x1b, 0x54, 0x4b, 0x7c, 0x74, 0x03, 0x09,
      0x89, 0xa9, 0x21, 0x55, 0xd7, 0x8b, 0xbf, 0x28, 0x53, 0x5f, 0xb5, 0x94,
      0x23, 0xb4, 0xcb, 0x64, 0x0e, 0x6c, 0x1e, 0xd4, 0x37, 0x15, 0xe8, 0x3c,
      0x72, 0x39, 0x1f, 0x15, 0x8a, 0x9d, 0x7a, 0xcc, 0x09, 0x4b, 0x00};

  const uint8_t ue_auth_response_hexbuf[21] = {
      0x7e, 0x0,  0x57, 0x2d, 0x10, 0x25, 0x70, 0x6f, 0x9a, 0x5b, 0x90,
      0xb6, 0xc9, 0x57, 0x50, 0x6c, 0x88, 0x3d, 0x76, 0xcc, 0x63};

  const uint8_t ue_auth_response_security_capability_mismatch_hexbuf[4] = {
      0x7e, 0x0, 0x59, 0x17};

  const uint8_t ue_auth_response_security_mode_reject_hexbuf[4] = {0x7e, 0x0,
                                                                   0x59, 0x18};

  const uint8_t ue_smc_response_hexbuf[64] = {
      0x7e, 0x4,  0x54, 0xf6, 0xe1, 0x2a, 0x0,  0x7e, 0x0,  0x5e, 0x77,
      0x0,  0x9,  0x45, 0x73, 0x80, 0x61, 0x21, 0x85, 0x61, 0x51, 0xf1,
      0x71, 0x0,  0x23, 0x7e, 0x0,  0x41, 0x79, 0x0,  0xd,  0x1,  0x22,
      0x62, 0x54, 0x0,  0x0,  0x0,  0x0,  0x0,  0x0,  0x0,  0x0,  0xf1,
      0x10, 0x1,  0x0,  0x2e, 0x08, 0x80, 0x20, 0x00, 0x00, 0x00, 0x00,
      0x00, 0x00, 0x2f, 0x2,  0x1,  0x1,  0x53, 0x1,  0x0};

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

  const uint8_t
      initial_ue_msg_service_request_mobile_trminated_no_UL_data_status[40] = {
          0x7e, 0x01, 0x30, 0x6f, 0xf2, 0xae, 0x03, 0x7e, 0x00, 0x4c,
          0x20, 0x00, 0x07, 0xf4, 0x00, 0x40, 0xd9, 0x58, 0xf8, 0x3b,
          0x71, 0x00, 0x11, 0x7e, 0x00, 0x4c, 0x00, 0x00, 0x07, 0xf4,
          0x00, 0x40, 0xd9, 0x58, 0xf8, 0x3b, 0x50, 0x02, 0x20, 0x00};

  const uint8_t service_req_wrong_tmsi[44] = {
      0x7e, 0x01, 0xca, 0x3f, 0x92, 0xbe, 0x03, 0x7e, 0x00, 0x4c, 0x10,
      0x00, 0x07, 0xf4, 0x00, 0x40, 0xff, 0xff, 0xff, 0xff, 0x71, 0x00,
      0x15, 0x7e, 0x00, 0x4c, 0x10, 0x00, 0x07, 0xf4, 0x00, 0x40, 0xff,
      0xff, 0xff, 0xff, 0x40, 0x02, 0x20, 0x00, 0x50, 0x02, 0x20, 0x00};
  const uint8_t initial_ue_periodic_reg[36] = {
      0x7e, 0x00, 0x41, 0x03, 0x00, 0x0b, 0xf2, 0x22, 0xf2, 0x07, 0x01, 0x00,
      0x40, 0x42, 0xdb, 0x2a, 0x42, 0x10, 0x01, 0x00, 0x2e, 0x02, 0x00, 0x00,
      0x2f, 0x02, 0x01, 0x02, 0x17, 0x02, 0xc0, 0xc0, 0xb0, 0x2b, 0x01, 0x00};

  const uint8_t initial_ue_msg_service_request[20] = {
      0x7e, 0x01, 0xca, 0x3f, 0x92, 0xbe, 0x03, 0x7e, 0x00, 0x4c,
      0x00, 0x00, 0x07, 0xf4, 0x00, 0x40, 0xff, 0xff, 0xff, 0xff};

  const uint8_t initial_ue_msg_service_request_high_priority[17] = {
      0x7e, 0x00, 0x4c, 0x50, 0x00, 0x07, 0xf4, 0x01, 0x01,
      0x88, 0xa4, 0x5d, 0x2d, 0x50, 0x02, 0x00, 0x00};

  const uint8_t pdu_sess_modification_complete_hex_buff[13] = {
      0x7e, 0x00, 0x67, 0x01, 0x00, 0x04, 0x2e,
      0x01, 0x01, 0xcc, 0x12, 0x05, 0x82};
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

bool validate_identification_procedure(uint32_t expected_retransmission_count,
                                       amf_ue_ngap_id_t* ue_id) {
  // By this time we should have one entry in ue_id table
  map_uint64_ue_context_t* amf_state_ue_id_ht = get_amf_ue_state();
  for (auto& elem : amf_state_ue_id_ht->umap) {
    *ue_id = elem.first;

    // Found the ue_id
    if (*ue_id) {
      break;
    }
  }

  amf_context_t* amf_ctx = get_amf_context_by_ueid(*ue_id);
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

  amf_ue_context_exists_amf_ue_ngap_id(ident_proc->ue_id);
  if (ident_proc->ue_id != *ue_id) {
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

TEST_F(AMFAppProcedureTest, TestRegistrationPlmnMismatch) {
  amf_ue_ngap_id_t ue_id = 0;
  std::vector<MessagesIds> expected_Ids{
      AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,  // new registration notification
                                            // indication to ngap
      NGAP_NAS_DL_DATA_REQ,                 // Registration Reject
      NGAP_UE_CONTEXT_RELEASE_COMMAND       // UEContextReleaseCommand
  };
  /* Send the initial UE message */
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(
      amf_app_desc_p, 36, 1, 1, 0, plmn,
      initial_ue_plmn_mismatch_message_hexbuf,
      sizeof(initial_ue_plmn_mismatch_message_hexbuf));
  EXPECT_EQ(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id), 0);
  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, TestRegistrationProcNoTMSI) {
  amf_ue_ngap_id_t ue_id = 0;
  std::vector<MessagesIds> expected_Ids{
      AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,  // new registration notification
                                            // indication to ngap
      NGAP_NAS_DL_DATA_REQ,                 // Authentication Request to UE
      NGAP_NAS_DL_DATA_REQ,            // Security Command Mode Request to UE
      NGAP_INITIAL_CONTEXT_SETUP_REQ,  // Initial Conext Setup Request to UE &
      NGAP_NAS_DL_DATA_REQ,            // De-registration accept
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

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for registration complete response from UE
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

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

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

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

  m5gmm_state_t mm_state;
  rc = amf_get_ue_context_mm_state(ue_id, &mm_state);
  EXPECT_TRUE(rc == RETURNok);
  EXPECT_TRUE(mm_state == DEREGISTERED);

  n2cause_e ue_context_rel_cause;
  rc = amf_get_ue_context_rel_cause(ue_id, &ue_context_rel_cause);
  EXPECT_TRUE(rc == RETURNok);

  // After sending the de-registration request the rel cause changes
  // to normal
  EXPECT_TRUE(ue_context_rel_cause == NGAP_INVALID_CAUSE);

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);

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
  send_initial_ue_message_with_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn, m_tmsi,
                                    guti_initial_ue_message_hexbuf,
                                    sizeof(guti_initial_ue_message_hexbuf));

  EXPECT_TRUE(validate_identification_procedure(0, &ue_id));

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

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  amf_app_handle_deregistration_req(ue_id);
  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, TestRegistrationProcGutiBasedEncryption) {
  amf_ue_ngap_id_t ue_id = 0;
  send_initial_ue_message_no_tmsi(
      amf_app_desc_p, 36, 1, 1, 0, plmn,
      guti_initial_ue_reregister_message_hexbuf,
      sizeof(guti_initial_ue_reregister_message_hexbuf));

  EXPECT_TRUE(validate_identification_procedure(0, &ue_id));

  int rc = RETURNok;
  rc = send_uplink_nas_identity_response_message(
      amf_app_desc_p, ue_id, plmn, identity_response_profile_a,
      sizeof(identity_response_profile_a));
  EXPECT_TRUE(rc == RETURNok);

  ue_m5gmm_context_t* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  rc = amf_decrypt_msin_info_answer(&decrypted_msin);
  EXPECT_TRUE(rc == RETURNok);

  ue_m5gmm_context_s* context_encrypted_imsi = amf_get_ue_context_from_imsi(
      reinterpret_cast<char*>(const_cast<char*>(decrypted_imsi.c_str())));

  // Check if UE Context is created with correct imsi
  bool res = false;
  res = get_ue_id_from_imsi(amf_app_desc_p,
                            context_encrypted_imsi->amf_context.imsi64, &ue_id);
  EXPECT_TRUE(res == true);

  // Send the authentication response message from subscriberdb
  rc = send_proc_authentication_info_answer(decrypted_imsi, ue_id, true);
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
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(decrypted_imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

  // Send uplink nas message for registration complete response from UE
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for deregistration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_TRUE(rc == RETURNok);

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);
  ue_context_p = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  delete ue_context_p;
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

  EXPECT_TRUE(validate_identification_procedure(0, &ue_id));

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

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

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

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

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

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);

  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, TestPDUSessionSetupWithoutContext) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id = 0;
  std::vector<MessagesIds> expected_Ids{
      AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,  // new registration notification
                                            // indication to ngap
      NGAP_NAS_DL_DATA_REQ,                 // Authentication Request to UE
      NGAP_NAS_DL_DATA_REQ,                 // Send Registration Accept
      NGAP_INITIAL_CONTEXT_SETUP_REQ,  // Initial Conext Setup Request to UE &
                                       // Pdu session establishment accept
      NGAP_PDUSESSION_RESOURCE_SETUP_REQ,
      NGAP_PDUSESSIONRESOURCE_REL_REQ,  // PDU Session Resource Release Command
      NGAP_NAS_DL_DATA_REQ,             // Deregistaration Accept
      NGAP_UE_CONTEXT_RELEASE_COMMAND   // UEContextReleaseCommand
  };

  // Send the initial UE message
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi_no_ctx_req(
      amf_app_desc_p, 36, 1, 1, 0, plmn, initial_ue_message_hexbuf,
      sizeof(initial_ue_message_hexbuf));

  // Check if UE Context is created with correct imsi
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id));

  // Send the authentication response message from subscriberdb
  rc = send_proc_authentication_info_answer(imsi, ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for auth response from UE
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for security mode complete response from UE
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
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

  // Send ip address response  from pipelined
  rc = send_ip_address_response_itti(IPv4);
  EXPECT_TRUE(rc == RETURNok);

  // Send pdu session setup response  from smf
  rc = send_pdu_session_response_itti(IPv4);
  EXPECT_TRUE(rc == RETURNok);

  // Send pdu resource setup response  from UE
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

  // Send uplink nas message for pdu session release request from UE
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, ue_id, plmn, pdu_sess_release_hexbuf,
      sizeof(pdu_sess_release_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for pdu session release complete from UE
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, ue_id, plmn, pdu_sess_release_complete_hexbuf,
      sizeof(pdu_sess_release_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  EXPECT_EQ(ue_context_p->amf_context.smf_ctxt_map.size(), 0);

  rc = send_pdu_notification_response();
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for deregistration complete response from UE
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_TRUE(rc == RETURNok);

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);

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

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

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

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);
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

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

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

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);
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

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

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

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);
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

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

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

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);
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

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

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

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);
  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, TestRegistrationProcSUCIExt) {
  amf_ue_ngap_id_t ue_id = 0;

  // Send the initial UE message
  imsi64_t imsi64 = 0;
  int rc = RETURNok;
  imsi64 = send_initial_ue_message_no_tmsi(
      amf_app_desc_p, 36, 1, 1, 0, plmn, intital_ue_message_suci_ext_hexbuf,
      sizeof(intital_ue_message_suci_ext_hexbuf));

  rc = amf_decrypt_msin_info_answer(&decrypted_msin);
  EXPECT_TRUE(rc == RETURNok);

  ue_m5gmm_context_s* context_encrypted_imsi =
      amf_get_ue_context_from_imsi((char*)decrypted_imsi.c_str());

  // Check if UE Context is created with correct imsi
  bool res = false;
  res = get_ue_id_from_imsi(amf_app_desc_p,
                            context_encrypted_imsi->amf_context.imsi64, &ue_id);
  EXPECT_TRUE(res == true);

  // Send the authentication response message from subscriberdb
  rc = send_proc_authentication_info_answer(decrypted_imsi, ue_id, true);
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
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(decrypted_imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

  // Send uplink nas message for registration complete response from UE
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for deregistration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_TRUE(rc == RETURNok);

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);
}

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

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

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

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);

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

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

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

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);

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
                                        NGAP_NAS_DL_DATA_REQ,
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

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

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

  /*Send UE context release request to move to idle mode*/
  send_ue_context_release_request_message(amf_app_desc_p, 1, 1, ue_id);

  // Send the ue context release complete message
  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);

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

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);

  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, ServiceRequestMTWithoutUplinkDataStatus) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id = 0;
  std::vector<MessagesIds> expected_Ids{AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_INITIAL_CONTEXT_SETUP_REQ,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND,
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

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

  send_initial_context_response(amf_app_desc_p, ue_id);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /*Send UE context release request to move to idle mode*/
  send_ue_context_release_request_message(amf_app_desc_p, 1, 1, ue_id);

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);

  ue_m5gmm_context_s* ue_context = nullptr;
  ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  EXPECT_EQ(REGISTERED_IDLE, ue_context->mm_state);

  /*Send initial UE message Service Request with PDU*/
  amf_ue_ngap_id_t updated_ue_id = 0;
  imsi64 = 0;
  imsi64 = send_initial_ue_message_service_request(
      amf_app_desc_p, 36, 1, 2, ue_id, plmn,
      initial_ue_msg_service_request_mobile_trminated_no_UL_data_status,
      sizeof(initial_ue_msg_service_request_mobile_trminated_no_UL_data_status),
      4);

  // Check if UE Context is created with correct imsi
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &updated_ue_id));

  EXPECT_EQ(REGISTERED_CONNECTED, ue_context->mm_state);

  // Ue id will be freshly generated
  EXPECT_NE(ue_context->amf_ue_ngap_id, ue_id);

  send_initial_context_response(amf_app_desc_p, updated_ue_id);

  // Send uplink nas message for deregistration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, updated_ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_TRUE(rc == RETURNok);

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 2, updated_ue_id);

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

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

  send_initial_context_response(amf_app_desc_p, ue_id);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /*Send UE context release request to move to idle mode*/
  send_ue_context_release_request_message(amf_app_desc_p, 1, 1, ue_id);

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);

  ue_m5gmm_context_s* ue_context = nullptr;
  ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  EXPECT_EQ(REGISTERED_IDLE, ue_context->mm_state);

  /*Send initial UE message Service Request with PDU*/
  amf_ue_ngap_id_t updated_ue_id = 0;
  imsi64 = 0;
  imsi64 = send_initial_ue_message_service_request(
      amf_app_desc_p, 36, 1, 2, ue_id, plmn,
      initial_ue_msg_service_request_signaling,
      sizeof(initial_ue_msg_service_request_signaling), 4);

  // Check if UE Context is created with correct imsi
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &updated_ue_id));

  EXPECT_EQ(REGISTERED_CONNECTED, ue_context->mm_state);

  // Ue id will be freshly generated
  EXPECT_NE(ue_context->amf_ue_ngap_id, ue_id);

  send_initial_context_response(amf_app_desc_p, updated_ue_id);

  // Send uplink nas message for deregistration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, updated_ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_TRUE(rc == RETURNok);

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 2, updated_ue_id);

  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, ImplicitDeregDuplicateSuciReg) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t init_ue_id = 0;
  std::vector<MessagesIds> expected_Ids{AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_INITIAL_CONTEXT_SETUP_REQ,
                                        NGAP_PDUSESSION_RESOURCE_SETUP_REQ,
                                        AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_INITIAL_CONTEXT_SETUP_REQ,
                                        NGAP_PDUSESSION_RESOURCE_SETUP_REQ,
                                        NGAP_PDUSESSIONRESOURCE_REL_REQ,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND};

  // Send the initial UE message
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));

  // Check if UE Context is created with correct imsi
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &init_ue_id));

  // Send the authentication response message from subscriberdb
  rc = send_proc_authentication_info_answer(imsi, init_ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for auth response from UE
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, init_ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for security mode complete response from UE
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, init_ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

  send_initial_context_response(amf_app_desc_p, init_ue_id);

  // Send uplink nas message for registration complete response from UE
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, init_ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for pdu session establishment request from UE
  rc = send_uplink_nas_pdu_session_establishment_request(
      amf_app_desc_p, init_ue_id, plmn, ue_pdu_session_est_req_hexbuf,
      sizeof(ue_pdu_session_est_req_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Send ip address response  from pipelined
  rc = send_ip_address_response_itti(IPv4);
  EXPECT_TRUE(rc == RETURNok);

  // Send pdu session setup response  from smf
  rc = send_pdu_session_response_itti(IPv4);
  EXPECT_TRUE(rc == RETURNok);

  // Send pdu resource setup response  from UE
  rc = send_pdu_resource_setup_response(init_ue_id);
  EXPECT_TRUE(rc == RETURNok);

  rc = send_pdu_notification_response();
  EXPECT_TRUE(rc == RETURNok);

  // Duplicate Registration Request
  amf_ue_ngap_id_t updated_ue_id = 0;
  imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 2, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));

  // Check if UE Context is created with correct imsi
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &updated_ue_id));

  // Send the authentication response message from subscriberdb
  rc = send_proc_authentication_info_answer(imsi, updated_ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for auth response from UE
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, updated_ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for security mode complete response from UE
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, updated_ue_id,
                                               plmn, ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

  send_initial_context_response(amf_app_desc_p, updated_ue_id);

  // Send uplink nas message for registration complete response from UE
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, updated_ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for pdu session establishment request from UE
  rc = send_uplink_nas_pdu_session_establishment_request(
      amf_app_desc_p, updated_ue_id, plmn, ue_pdu_session_est_req_hexbuf,
      sizeof(ue_pdu_session_est_req_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Send ip address response  from pipelined
  rc = send_ip_address_response_itti(IPv4);
  EXPECT_TRUE(rc == RETURNok);
  // Send pdu session setup response  from smf
  rc = send_pdu_session_response_itti(IPv4);
  EXPECT_TRUE(rc == RETURNok);

  // Send pdu resource setup response  from UE
  rc = send_pdu_resource_setup_response(updated_ue_id);
  EXPECT_TRUE(rc == RETURNok);

  rc = send_pdu_notification_response();
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for pdu session release request from UE
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, updated_ue_id, plmn, pdu_sess_release_hexbuf,
      sizeof(pdu_sess_release_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for pdu session release complete from UE
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, updated_ue_id, plmn, pdu_sess_release_complete_hexbuf,
      sizeof(pdu_sess_release_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  ue_m5gmm_context_t* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id(updated_ue_id);
  ASSERT_NE(ue_context_p, nullptr);
  EXPECT_EQ(ue_context_p->amf_context.smf_ctxt_map.size(), 0);

  rc = send_pdu_notification_response();
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for deregistration complete response from UE
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, updated_ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_TRUE(rc == RETURNok);
  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, updated_ue_id);
  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, ServiceRequestWrongTMSI) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t init_ue_id = 0;
  std::vector<MessagesIds> expected_Ids{AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_INITIAL_CONTEXT_SETUP_REQ,
                                        NGAP_PDUSESSION_RESOURCE_SETUP_REQ,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND,
                                        AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND,
                                        AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_INITIAL_CONTEXT_SETUP_REQ,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND};

  // Send the initial UE message
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));

  // Check if UE Context is created with correct imsi
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &init_ue_id));

  // Send the authentication response message from subscriberdb
  rc = send_proc_authentication_info_answer(imsi, init_ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for auth response from UE
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, init_ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for security mode complete response from UE
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, init_ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

  send_initial_context_response(amf_app_desc_p, init_ue_id);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, init_ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  send_initial_context_response(amf_app_desc_p, init_ue_id);

  /* Send uplink nas message for pdu session establishment request from UE */
  rc = send_uplink_nas_pdu_session_establishment_request(
      amf_app_desc_p, init_ue_id, plmn, ue_pdu_session_est_req_hexbuf,
      sizeof(ue_pdu_session_est_req_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Send ip address response  from pipelined
  rc = send_ip_address_response_itti(IPv4);
  EXPECT_TRUE(rc == RETURNok);

  // Send pdu session setup response  from smf
  rc = send_pdu_session_response_itti(IPv4);
  EXPECT_TRUE(rc == RETURNok);

  // Send pdu resource setup response  from UE
  rc = send_pdu_resource_setup_response(init_ue_id);
  EXPECT_TRUE(rc == RETURNok);

  rc = send_pdu_notification_response();
  EXPECT_TRUE(rc == RETURNok);

  send_ue_context_release_request_message(amf_app_desc_p, 1, 1, init_ue_id);

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, init_ue_id);

  amf_ue_ngap_id_t service_request_ue_id = 0;
  send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 2, 0, plmn,
                                  service_req_wrong_tmsi,
                                  sizeof(service_req_wrong_tmsi));

  // Duplicate Registration Request
  amf_ue_ngap_id_t updated_ue_id = 0;
  imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 2, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));

  // Check if UE Context is created with correct imsi
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &updated_ue_id));

  // Send the authentication response message from subscriberdb
  rc = send_proc_authentication_info_answer(imsi, updated_ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for auth response from UE
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, updated_ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for security mode complete response from UE
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, updated_ue_id,
                                               plmn, ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

  send_initial_context_response(amf_app_desc_p, updated_ue_id);

  // Send uplink nas message for registration complete response from UE
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, updated_ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  ue_m5gmm_context_t* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id(updated_ue_id);
  ASSERT_NE(ue_context_p, nullptr);
  EXPECT_EQ(ue_context_p->amf_context.smf_ctxt_map.size(), 0);

  // Send uplink nas message for deregistration complete response from UE
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, updated_ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_TRUE(rc == RETURNok);
  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, updated_ue_id);
  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, ServiceRequestSignalWithPDU) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id = 0;
  std::vector<MessagesIds> expected_Ids{AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_INITIAL_CONTEXT_SETUP_REQ,
                                        NGAP_PDUSESSION_RESOURCE_SETUP_REQ,
                                        NGAP_PDUSESSIONRESOURCE_REL_REQ,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND,
                                        AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND,
                                        NGAP_PDUSESSIONRESOURCE_REL_REQ,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND};

  // Send the initial UE message
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));

  // Check if UE Context is created with correct imsi
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id));

  // Send the authentication response message from subscriberdb
  rc = send_proc_authentication_info_answer(imsi, ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for auth response from UE
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for security mode complete response from UE
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

  send_initial_context_response(amf_app_desc_p, ue_id);

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

  // Send ip address response  from pipelined
  rc = send_ip_address_response_itti(IPv4);
  EXPECT_TRUE(rc == RETURNok);

  // Send pdu session setup response  from smf
  rc = send_pdu_session_response_itti(IPv4);
  EXPECT_TRUE(rc == RETURNok);

  // Send pdu resource setup response  from UE
  rc = send_pdu_resource_setup_response(ue_id);
  EXPECT_TRUE(rc == RETURNok);

  rc = send_pdu_notification_response();
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for pdu session release request from UE
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, ue_id, plmn, pdu_sess_release_hexbuf,
      sizeof(pdu_sess_release_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for pdu session release complete from UE
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, ue_id, plmn, pdu_sess_release_complete_hexbuf,
      sizeof(pdu_sess_release_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  send_ue_context_release_request_message(amf_app_desc_p, 1, 1, ue_id);

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);

  // Service Request
  amf_ue_ngap_id_t updated_ue_id = 0;
  imsi64 = 0;
  imsi64 = send_initial_ue_message_service_request(
      amf_app_desc_p, 36, 1, 2, ue_id, plmn, initial_ue_msg_service_request,
      sizeof(initial_ue_msg_service_request), 4);

  // Check if UE Context is created with correct imsi
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &updated_ue_id));

  send_initial_context_response(amf_app_desc_p, updated_ue_id);

  // Send uplink nas message for pdu session establishment request from UE
  rc = send_uplink_nas_pdu_session_establishment_request(
      amf_app_desc_p, updated_ue_id, plmn, ue_pdu_session_est_req_hexbuf,
      sizeof(ue_pdu_session_est_req_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for pdu session release request from UE
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, updated_ue_id, plmn, pdu_sess_release_hexbuf,
      sizeof(pdu_sess_release_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for pdu session release complete from UE
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, updated_ue_id, plmn, pdu_sess_release_complete_hexbuf,
      sizeof(pdu_sess_release_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  rc = send_pdu_notification_response();
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for deregistration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, updated_ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_TRUE(rc == RETURNok);
  send_ue_context_release_complete_message(amf_app_desc_p, 1, 2, updated_ue_id);
  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, PeriodicRegistraionNoTmsi) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t init_ue_id = 0;
  std::vector<MessagesIds> expected_Ids{AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_INITIAL_CONTEXT_SETUP_REQ,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND,
                                        AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND};
  // Send the initial UE message
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));

  // Check if UE Context is created with correct imsi
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &init_ue_id));

  // Send the authentication response message from subscriberdb
  rc = send_proc_authentication_info_answer(imsi, init_ue_id, true);
  EXPECT_TRUE(rc == RETURNok);
  // Send uplink nas message for auth response from UE
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, init_ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);
  // Send uplink nas message for security mode complete response from UE
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, init_ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

  send_initial_context_response(amf_app_desc_p, init_ue_id);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, init_ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  send_initial_context_response(amf_app_desc_p, init_ue_id);

  send_ue_context_release_request_message(amf_app_desc_p, 1, 1, init_ue_id);

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, init_ue_id);
  ue_m5gmm_context_s* ue_context_p = nullptr;
  ue_context_p = amf_ue_context_exists_amf_ue_ngap_id(init_ue_id);
  delete ue_context_p;
  uint32_t m_tmsi = 0x42db2a42;
  send_initial_ue_message_with_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn, m_tmsi,
                                    initial_ue_periodic_reg,
                                    sizeof(initial_ue_periodic_reg));

  EXPECT_TRUE(validate_identification_procedure(0, &init_ue_id));

  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, init_ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  amf_app_handle_deregistration_req(init_ue_id);
  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, ReRegistraion) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t init_ue_id = 0;
  std::vector<MessagesIds> expected_Ids{AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_INITIAL_CONTEXT_SETUP_REQ,
                                        AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_INITIAL_CONTEXT_SETUP_REQ,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND};

  // Send the initial UE message
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));

  // Check if UE Context is created with correct imsi
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &init_ue_id));

  // Send the authentication response message from subscriberdb
  rc = send_proc_authentication_info_answer(imsi, init_ue_id, true);
  EXPECT_TRUE(rc == RETURNok);
  // Send uplink nas message for auth response from UE
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, init_ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);
  // Send uplink nas message for security mode complete response from UE
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, init_ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

  send_initial_context_response(amf_app_desc_p, init_ue_id);

  // Send uplink nas message for registration complete response from UE
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, init_ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  send_initial_context_response(amf_app_desc_p, init_ue_id);

  amf_ue_ngap_id_t updated_ue_id = 0;
  imsi64 = 0;

  // Replace the tmsi with tmsi generated during first registration

  imsi64 = send_initial_ue_message_no_tmsi_replace_mtmsi(
      amf_app_desc_p, 36, 1, 2, 0, plmn,
      guti_initial_ue_reregister_message_hexbuf,
      sizeof(guti_initial_ue_reregister_message_hexbuf), init_ue_id, 19);

  EXPECT_TRUE(validate_identification_procedure(0, &updated_ue_id));

  rc = send_uplink_nas_identity_response_message(amf_app_desc_p, updated_ue_id,
                                                 plmn, identity_response,
                                                 sizeof(identity_response));
  EXPECT_TRUE(rc == RETURNok);

  ue_m5gmm_context_t* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id(updated_ue_id);
  ASSERT_NE(ue_context_p, nullptr);
  EXPECT_TRUE(ue_context_p->amf_context.imsi64 == stoul(imsi));

  // Send the authentication response message from subscriberdb
  rc = send_proc_authentication_info_answer(imsi, updated_ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  // Validate if authentication procedure is initialized as expected
  EXPECT_TRUE(validate_auth_procedure(updated_ue_id, 0));

  // Send uplink nas message for auth response from UE
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, updated_ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Check whether security mode procedure is initiated
  EXPECT_TRUE(validate_smc_procedure(updated_ue_id, 0));

  // Send uplink nas message for security mode complete response from UE
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, updated_ue_id,
                                               plmn, ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

  // Send uplink nas message for registration complete response from UE
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, updated_ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);
  amf_app_handle_deregistration_req(updated_ue_id);
  ue_context_p = amf_ue_context_exists_amf_ue_ngap_id(init_ue_id);
  delete ue_context_p;
  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}
TEST_F(AMFAppProcedureTest, SctpShutWithServiceRequest) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id = 0;
  std::vector<MessagesIds> expected_Ids{AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_INITIAL_CONTEXT_SETUP_REQ,
                                        NGAP_PDUSESSION_RESOURCE_SETUP_REQ,
                                        NGAP_PDUSESSIONRESOURCE_REL_REQ,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND,
                                        AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND,
                                        NGAP_PDUSESSIONRESOURCE_REL_REQ,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND};

  // Send the initial UE message
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));

  // Check if UE Context is created with correct imsi
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id));

  // Send the authentication response message from subscriberdb
  rc = send_proc_authentication_info_answer(imsi, ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for auth response from UE
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for security mode complete response from UE
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

  send_initial_context_response(amf_app_desc_p, ue_id);

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

  // Send ip address response  from pipelined
  rc = send_ip_address_response_itti(IPv4);
  EXPECT_TRUE(rc == RETURNok);

  // Send pdu session setup response  from smf
  rc = send_pdu_session_response_itti(IPv4);
  EXPECT_TRUE(rc == RETURNok);

  // Send pdu resource setup response  from UE
  rc = send_pdu_resource_setup_response(ue_id);
  EXPECT_TRUE(rc == RETURNok);

  rc = send_pdu_notification_response();
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for pdu session release request from UE
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, ue_id, plmn, pdu_sess_release_hexbuf,
      sizeof(pdu_sess_release_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for pdu session release complete from UE
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, ue_id, plmn, pdu_sess_release_complete_hexbuf,
      sizeof(pdu_sess_release_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  itti_ngap_gNB_deregistered_ind_t gnb_derg_message = {};
  gnb_derg_message.nb_ue_to_deregister = 1;
  gnb_derg_message.gnb_ue_ngap_id[0] = 1;
  gnb_derg_message.amf_ue_ngap_id[0] = ue_id;
  gnb_derg_message.gnb_id = 1;

  amf_app_handle_gnb_deregister_ind(&gnb_derg_message);

  send_ue_context_release_request_message(amf_app_desc_p, 1, 1, ue_id);

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);

  int res = RETURNerror;

  // Check the states of UE
  res = check_ue_context_state(ue_id, REGISTERED_IDLE, M5GCM_IDLE);
  EXPECT_TRUE(res == RETURNok);

  // Service Request
  amf_ue_ngap_id_t updated_ue_id = 0;
  imsi64 = 0;
  imsi64 = send_initial_ue_message_service_request(
      amf_app_desc_p, 36, 1, 2, ue_id, plmn, initial_ue_msg_service_request,
      sizeof(initial_ue_msg_service_request), 4);

  // Check if UE Context is created with correct imsi
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &updated_ue_id));

  send_initial_context_response(amf_app_desc_p, updated_ue_id);

  // Send uplink nas message for pdu session establishment request from UE
  rc = send_uplink_nas_pdu_session_establishment_request(
      amf_app_desc_p, updated_ue_id, plmn, ue_pdu_session_est_req_hexbuf,
      sizeof(ue_pdu_session_est_req_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for pdu session release request from UE
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, updated_ue_id, plmn, pdu_sess_release_hexbuf,
      sizeof(pdu_sess_release_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for pdu session release complete from UE
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, updated_ue_id, plmn, pdu_sess_release_complete_hexbuf,
      sizeof(pdu_sess_release_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  rc = send_pdu_notification_response();
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for deregistration complete response from UE
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, updated_ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_TRUE(rc == RETURNok);
  send_ue_context_release_complete_message(amf_app_desc_p, 1, 2, updated_ue_id);
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

TEST_F(AMFAppProcedureTest, TestPDUSession_LocationUpdateFail) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id = 0;
  std::vector<MessagesIds> expected_Ids{
      AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,  // new registration notification
                                            // indication to ngap
      NGAP_NAS_DL_DATA_REQ,                 // Authentication Request to UE
      NGAP_NAS_DL_DATA_REQ,  // Security Command Mode Request to UE
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

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  snprintf(ula_ans.imsi, sizeof(ula_ans.imsi), "%s", "123456789012345");
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNerror);

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
      amf_app_desc_p, ue_id, plmn, ue_pdu_session_est_req_hexbuf,
      sizeof(ue_pdu_session_est_req_hexbuf));
  EXPECT_EQ(rc, RETURNok);

  ue_context_p = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  // ue context should exist
  ASSERT_NE(ue_context_p, nullptr);
  // smf context should not be present
  EXPECT_EQ(ue_context_p->amf_context.smf_ctxt_map.size(), 0);

  /* Send uplink nas message for deregistration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_TRUE(rc == RETURNok);

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);
  EXPECT_EQ(expected_Ids, AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, TestPDUSessionResourceModify) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id = 0;

  /* Send the initial UE message */
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));

  /* Check if UE Context is created with correct imsi */
  bool res = false;
  res = get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id);
  EXPECT_TRUE(res == true);

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

  /* Send pdu session setup response  from smf */
  rc = send_pdu_session_modification_itti();
  EXPECT_TRUE(rc == RETURNok);

  /* Send pdu resource setup response  from gnb */
  rc = send_pdu_resource_modify_response(ue_id);
  EXPECT_TRUE(rc == RETURNok);

  // Send pdu session modification complete
  /* Send uplink nas message for pdu session complete from UE */
  rc = send_uplink_nas_pdu_session_modification_complete(
      amf_app_desc_p, ue_id, plmn, pdu_sess_modification_complete_hex_buff,
      sizeof(pdu_sess_modification_complete_hex_buff));
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

  rc = send_pdu_notification_response();
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for deregistration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);
}

TEST_F(AMFAppProcedureTest, RegistrationAfterFourRegAcceptMsgs) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id = 0;
  std::vector<MessagesIds> expected_Ids{AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_INITIAL_CONTEXT_SETUP_REQ,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND,
                                        AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,
                                        NGAP_INITIAL_CONTEXT_SETUP_REQ,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND};

  // Send the initial UE message
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));
  // Check if UE Context is created with correct imsi
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id));

  // Send the authentication response message from subscriberdb
  rc = send_proc_authentication_info_answer(imsi, ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for auth response from UE
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for security mode complete response from UE
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

  // mark registration complete post 5th expiry of T3550

  rc = unit_test_registration_accept_t3550(ue_id);
  EXPECT_TRUE(rc == RETURNok);
  ue_m5gmm_context_s* ue_context = nullptr;
  ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  EXPECT_EQ(REGISTERED_CONNECTED, ue_context->mm_state);

  send_initial_context_response(amf_app_desc_p, ue_id);

  // Send initial UE message Service Request with PDU
  amf_ue_ngap_id_t updated_ue_id = 0;
  imsi64 = 0;
  imsi64 = send_initial_ue_message_service_request(
      amf_app_desc_p, 36, 1, 2, ue_id, plmn,
      initial_ue_msg_service_request_signaling,
      sizeof(initial_ue_msg_service_request_signaling), 4);

  // Check if UE Context is created with correct imsi
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &updated_ue_id));

  EXPECT_EQ(REGISTERED_CONNECTED, ue_context->mm_state);

  // Ue id will be freshly generated
  EXPECT_NE(ue_context->amf_ue_ngap_id, ue_id);

  send_initial_context_response(amf_app_desc_p, updated_ue_id);

  // Send uplink nas message for deregistration complete response from UE
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, updated_ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_TRUE(rc == RETURNok);

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 2, updated_ue_id);

  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, TestPDUSessionResourceModifyDeletion) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id = 0;

  /* Send the initial UE message */
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));

  /* Check if UE Context is created with correct imsi */
  bool res = false;
  res = get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id);
  EXPECT_TRUE(res == true);

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

  /* Send pdu session setup response  from smf */
  rc = send_pdu_session_modification_deletion_itti();
  EXPECT_TRUE(rc == RETURNok);

  /* Send pdu resource setup response  from gnb */
  rc = send_pdu_resource_modify_response(ue_id);
  EXPECT_TRUE(rc == RETURNok);

  // Send pdu session modification complete
  /* Send uplink nas message for pdu session complete from UE */
  rc = send_uplink_nas_pdu_session_modification_complete(
      amf_app_desc_p, ue_id, plmn, pdu_sess_modification_complete_hex_buff,
      sizeof(pdu_sess_modification_complete_hex_buff));
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

  rc = send_pdu_notification_response();
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for deregistration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);
}

TEST_F(AMFAppProcedureTest, GnbInitiatedNGReset) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id = 0;
  std::vector<MessagesIds> expected_Ids{AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_INITIAL_CONTEXT_SETUP_REQ,
                                        AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_INITIAL_CONTEXT_SETUP_REQ,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND,
                                        NGAP_GNB_INITIATED_RESET_ACK,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND};

  /* UE-1-Registration */
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

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  send_initial_context_response(amf_app_desc_p, ue_id);

  /* UE-2-Registration */
  const uint8_t initial_ue_message_hexbuf_temp[25] = {
      0x7e, 0x00, 0x41, 0x79, 0x00, 0x0d, 0x01, 0x22, 0x62,
      0x54, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
      0x02, 0x2e, 0x04, 0xf0, 0xf0, 0xf0, 0xf0};

  amf_ue_ngap_id_t ue_id_temp = 0;

  imsi64_t imsi64_temp = 0;
  imsi64_temp = send_initial_ue_message_no_tmsi(
      amf_app_desc_p, 36, 1, 2, 0, plmn, initial_ue_message_hexbuf_temp,
      sizeof(initial_ue_message_hexbuf_temp));

  /* Check if UE Context is created with correct imsi */
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64_temp, &ue_id_temp));

  std::string imsi_temp = "222456000000002";
  /* Send the authentication response message from subscriberdb */
  rc = send_proc_authentication_info_answer(imsi_temp, ue_id_temp, true);
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for auth response from UE */
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, ue_id_temp, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for security mode complete response from UE */
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, ue_id_temp, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  s6a_update_location_ans_t ula_ans1 = util_amf_send_s6a_ula(imsi_temp);
  rc = amf_handle_s6a_update_location_ans(&ula_ans1);
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id_temp, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  send_initial_context_response(amf_app_desc_p, ue_id_temp);

  /* Send GNB reset request */
  send_gnb_reset_req();
  EXPECT_EQ(
      amf_app_desc_p->amf_ue_contexts.gnb_ue_ngap_id_ue_context_htbl.size(), 1);

  /* UE-1-Deregistration */
  /* Send uplink nas message for deregistration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_TRUE(rc == RETURNok);

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);

  /* UE-2-Deregistration */
  /* Send uplink nas message for deregistration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, ue_id_temp, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  send_ue_context_release_complete_message(amf_app_desc_p, 2, 2, ue_id_temp);

  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, ServiceRequestHighPriority) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id = 0;
  std::vector<MessagesIds> expected_Ids{AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_INITIAL_CONTEXT_SETUP_REQ,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND,
                                        AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION,
                                        NGAP_INITIAL_CONTEXT_SETUP_REQ,
                                        NGAP_PDUSESSION_RESOURCE_SETUP_REQ,
                                        NGAP_PDUSESSIONRESOURCE_REL_REQ,
                                        NGAP_NAS_DL_DATA_REQ,
                                        NGAP_UE_CONTEXT_RELEASE_COMMAND};

  // Send the initial UE message
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));

  // Check if UE Context is created with correct imsi
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id));

  // Send the authentication response message from subscriberdb
  rc = send_proc_authentication_info_answer(imsi, ue_id, true);
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for auth response from UE
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for security mode complete response from UE
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

  send_initial_context_response(amf_app_desc_p, ue_id);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  send_ue_context_release_request_message(amf_app_desc_p, 1, 1, ue_id);

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);

  // Service Request
  amf_ue_ngap_id_t updated_ue_id = 0;
  imsi64 = 0;
  imsi64 = send_initial_ue_message_service_request(
      amf_app_desc_p, 36, 1, 2, ue_id, plmn,
      initial_ue_msg_service_request_high_priority,
      sizeof(initial_ue_msg_service_request_high_priority), 4);

  // Check if UE Context is created with correct imsi
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &updated_ue_id));

  send_initial_context_response(amf_app_desc_p, updated_ue_id);

  /* Send uplink nas message for pdu session establishment request from UE */
  rc = send_uplink_nas_pdu_session_establishment_request(
      amf_app_desc_p, updated_ue_id, plmn, ue_pdu_session_est_req_hexbuf,
      sizeof(ue_pdu_session_est_req_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Send ip address response  from pipelined
  rc = send_ip_address_response_itti(IPv4);
  EXPECT_TRUE(rc == RETURNok);

  // Send pdu session setup response  from smf
  rc = send_pdu_session_response_itti(IPv4);
  EXPECT_TRUE(rc == RETURNok);

  // Send pdu resource setup response  from UE
  rc = send_pdu_resource_setup_response(updated_ue_id);
  EXPECT_TRUE(rc == RETURNok);

  rc = send_pdu_notification_response();
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for pdu session release request from UE
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, updated_ue_id, plmn, pdu_sess_release_hexbuf,
      sizeof(pdu_sess_release_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Send uplink nas message for pdu session release complete from UE
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, updated_ue_id, plmn, pdu_sess_release_complete_hexbuf,
      sizeof(pdu_sess_release_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  rc = send_pdu_notification_response();
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for deregistration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, updated_ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 2, updated_ue_id);
  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}
}  // namespace magma5g
