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

extern "C" {
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/include/amf_config.h"
}
#include "lte/gateway/c/core/oai/tasks/amf/include/amf_client_servicer.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_state_manager.h"
#include "lte/gateway/c/core/oai/test/amf/amf_app_test_util.h"

using ::testing::Test;

namespace magma5g {

extern task_zmq_ctx_s amf_app_task_zmq_ctx;

TEST(test_state_converter, test_guti_to_string) {
  guti_m5_t guti1, guti2;
  guti1.guamfi.plmn.mcc_digit1 = 2;
  guti1.guamfi.plmn.mcc_digit2 = 2;
  guti1.guamfi.plmn.mcc_digit3 = 2;
  guti1.guamfi.plmn.mnc_digit1 = 4;
  guti1.guamfi.plmn.mnc_digit2 = 5;
  guti1.guamfi.plmn.mnc_digit3 = 6;
  guti1.guamfi.amf_regionid    = 1;
  guti1.guamfi.amf_set_id      = 1;
  guti1.guamfi.amf_pointer     = 0;
  guti1.m_tmsi                 = 0X212e5025;

  std::string guti1_str =
      AmfNasStateConverter::amf_app_convert_guti_m5_to_string(guti1);

  AmfNasStateConverter::amf_app_convert_string_to_guti_m5(guti1_str, &guti2);

  EXPECT_EQ(guti1.guamfi.plmn.mcc_digit1, guti2.guamfi.plmn.mcc_digit1);
  EXPECT_EQ(guti1.guamfi.plmn.mcc_digit2, guti2.guamfi.plmn.mcc_digit2);
  EXPECT_EQ(guti1.guamfi.plmn.mcc_digit3, guti2.guamfi.plmn.mcc_digit3);
  EXPECT_EQ(guti1.guamfi.plmn.mnc_digit1, guti2.guamfi.plmn.mnc_digit1);
  EXPECT_EQ(guti1.guamfi.plmn.mnc_digit2, guti2.guamfi.plmn.mnc_digit2);
  EXPECT_EQ(guti1.guamfi.plmn.mnc_digit3, guti2.guamfi.plmn.mnc_digit3);
  EXPECT_EQ(guti1.guamfi.amf_regionid, guti2.guamfi.amf_regionid);
  EXPECT_EQ(guti1.guamfi.amf_set_id, guti2.guamfi.amf_set_id);
  EXPECT_EQ(guti1.guamfi.amf_pointer, guti2.guamfi.amf_pointer);
  EXPECT_EQ(guti1.m_tmsi, guti2.m_tmsi);
}

TEST(test_state_converter, test_state_to_proto) {
  // Guti setup
  guti_m5_t guti1;
  memset(&guti1, 0, sizeof(guti1));

  guti1.guamfi.plmn.mcc_digit1 = 2;
  guti1.guamfi.plmn.mcc_digit2 = 2;
  guti1.guamfi.plmn.mcc_digit3 = 2;
  guti1.guamfi.plmn.mnc_digit1 = 4;
  guti1.guamfi.plmn.mnc_digit2 = 5;
  guti1.guamfi.plmn.mnc_digit3 = 6;
  guti1.guamfi.amf_regionid    = 1;
  guti1.guamfi.amf_set_id      = 1;
  guti1.guamfi.amf_pointer     = 0;
  guti1.m_tmsi                 = 556683301;

  amf_app_desc_t amf_app_desc1 = {}, amf_app_desc2 = {};
  magma::lte::oai::MmeNasState state_proto = magma::lte::oai::MmeNasState();
  uint64_t data                            = 0;

  amf_app_desc1.amf_app_ue_ngap_id_generator = 0x05;
  amf_app_desc1.amf_ue_contexts.imsi_amf_ue_id_htbl.insert(1, 10);
  amf_app_desc1.amf_ue_contexts.tun11_ue_context_htbl.insert(2, 20);
  amf_app_desc1.amf_ue_contexts.gnb_ue_ngap_id_ue_context_htbl.insert(3, 30);
  amf_app_desc1.amf_ue_contexts.guti_ue_context_htbl.insert(guti1, 40);

  AmfNasStateConverter::state_to_proto(&amf_app_desc1, &state_proto);

  AmfNasStateConverter::proto_to_state(state_proto, &amf_app_desc2);

  EXPECT_EQ(
      amf_app_desc1.amf_app_ue_ngap_id_generator,
      amf_app_desc2.amf_app_ue_ngap_id_generator);

  EXPECT_EQ(
      amf_app_desc2.amf_ue_contexts.imsi_amf_ue_id_htbl.get(1, &data),
      magma::MAP_OK);
  EXPECT_EQ(data, 10);
  data = 0;

  EXPECT_EQ(
      amf_app_desc2.amf_ue_contexts.tun11_ue_context_htbl.get(2, &data),
      magma::MAP_OK);
  EXPECT_EQ(data, 20);
  data = 0;

  EXPECT_EQ(
      amf_app_desc2.amf_ue_contexts.gnb_ue_ngap_id_ue_context_htbl.get(
          3, &data),
      magma::MAP_OK);
  EXPECT_EQ(data, 30);
  data = 0;

  EXPECT_EQ(
      amf_app_desc2.amf_ue_contexts.guti_ue_context_htbl.get(guti1, &data),
      magma::MAP_OK);
  EXPECT_EQ(data, 40);
}

class AMFAppStatelessTest : public ::testing::Test {
 protected:
  virtual void SetUp() {
    itti_init(
        TASK_MAX, THREAD_MAX, MESSAGES_ID_MAX, tasks_info, messages_info, NULL,
        NULL);

    // initialize amf config
    amf_config_init(&amf_config);
    amf_config.use_stateless = true;
    amf_nas_state_init(&amf_config);
    create_state_matrix();

    init_task_context(TASK_MAIN, nullptr, 0, NULL, &amf_app_task_zmq_ctx);

    amf_app_desc_p = get_amf_nas_state(true);
  }

  virtual void TearDown() {
    clear_amf_nas_state();
    clear_amf_config(&amf_config);
    destroy_task_context(&amf_app_task_zmq_ctx);
    itti_free_desc_threads();
  }

  amf_app_desc_t* amf_app_desc_p;
  std::string imsi = "222456000000001";
  plmn_t plmn      = {.mcc_digit2 = 0,
                 .mcc_digit1 = 0,
                 .mnc_digit3 = 0x0f,
                 .mcc_digit3 = 1,
                 .mnc_digit2 = 1,
                 .mnc_digit1 = 0};

  const uint8_t initial_ue_message_hexbuf[25] = {
      0x7e, 0x00, 0x41, 0x79, 0x00, 0x0d, 0x01, 0x22, 0x62,
      0x54, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
      0x01, 0x2e, 0x04, 0xf0, 0xf0, 0xf0, 0xf0};

  const uint8_t ue_auth_response_hexbuf[21] = {
      0x7e, 0x0,  0x57, 0x2d, 0x10, 0x25, 0x70, 0x6f, 0x9a, 0x5b, 0x90,
      0xb6, 0xc9, 0x57, 0x50, 0x6c, 0x88, 0x3d, 0x76, 0xcc, 0x63};

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

  const uint8_t pdu_sess_release_hexbuf[14] = {0x7e, 0x00, 0x67, 0x01, 0x00,
                                               0x06, 0x2e, 0x01, 0x01, 0xd1,
                                               0x59, 0x24, 0x12, 0x01};

  const uint8_t pdu_sess_release_complete_hexbuf[12] = {
      0x7e, 0x00, 0x67, 0x01, 0x00, 0x04, 0x2e, 0x01, 0x01, 0xd4, 0x12, 0x01};

  uint8_t ue_initiated_dereg_hexbuf[24] = {
      0x7e, 0x01, 0x41, 0x21, 0xe6, 0xe2, 0x03, 0x7e, 0x00, 0x45, 0x01, 0x00,
      0x0b, 0xf2, 0x22, 0x62, 0x54, 0x01, 0x00, 0x40, 0x0,  0x0,  0x0,  0x0};
};

TEST_F(AMFAppStatelessTest, TestStateless) {
  int rc                 = RETURNerror;
  amf_ue_ngap_id_t ue_id = 0;

  /* Send the initial UE message */
  imsi64_t imsi64 = 0;
  imsi64          = send_initial_ue_message_no_tmsi(
      amf_app_desc_p, 36, 1, 1, 0, plmn, initial_ue_message_hexbuf,
      sizeof(initial_ue_message_hexbuf));
  AMFClientServicer::getInstance().map_tableKey_protoStr.clear();
  EXPECT_TRUE(AMFClientServicer::getInstance().map_tableKey_protoStr.isEmpty());
  // Writes the state to the data store
  put_amf_nas_state();
  EXPECT_EQ(
      AMFClientServicer::getInstance().map_tableKey_protoStr.isEmpty(), false);

  /* Check if UE Context is created with correct imsi */
  bool res = false;
  res      = get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id);
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
  rc = send_uplink_nas_message_ue_smc_response(
      amf_app_desc_p, ue_id, plmn, ue_smc_response_hexbuf,
      sizeof(ue_smc_response_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  send_initial_context_response(amf_app_desc_p, ue_id);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  // Clears the state
  AMFAppStatelessTest::TearDown();
  EXPECT_TRUE(amf_app_desc_p->amf_ue_contexts.imsi_amf_ue_id_htbl.isEmpty());
  EXPECT_TRUE(amf_app_desc_p->amf_ue_contexts.tun11_ue_context_htbl.isEmpty());
  EXPECT_TRUE(amf_app_desc_p->amf_ue_contexts.guti_ue_context_htbl.isEmpty());
  // Internally reads back the state
  AMFAppStatelessTest::SetUp();
  EXPECT_EQ(
      amf_app_desc_p->amf_ue_contexts.imsi_amf_ue_id_htbl.isEmpty(), false);
  EXPECT_EQ(
      amf_app_desc_p->amf_ue_contexts.tun11_ue_context_htbl.isEmpty(), false);
  EXPECT_EQ(
      amf_app_desc_p->amf_ue_contexts.guti_ue_context_htbl.isEmpty(), false);

  /* Send uplink nas message for pdu session establishment request from UE */
  rc = send_uplink_nas_pdu_session_establishment_request(
      amf_app_desc_p, ue_id, plmn, ue_pdu_session_est_req_hexbuf,
      sizeof(ue_pdu_session_est_req_hexbuf));
  EXPECT_TRUE(rc == RETURNok);

  /* Send ip address response  from pipelined */
  rc = send_ip_address_response_itti();
  EXPECT_TRUE(rc == RETURNok);

  /* Send pdu session setup response  from smf */
  rc = send_pdu_session_response_itti();
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

  rc = send_pdu_notification_response();
  EXPECT_TRUE(rc == RETURNok);

  /* Send uplink nas message for deregistration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_TRUE(rc == RETURNok);

  AMFClientServicer::getInstance().map_tableKey_protoStr.clear();
  EXPECT_TRUE(AMFClientServicer::getInstance().map_tableKey_protoStr.isEmpty());
}
}  // namespace magma5g
