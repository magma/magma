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

class AMFAppProcedureTest : public ::testing::Test {
  virtual void SetUp() {
    itti_init(
        TASK_MAX, THREAD_MAX, MESSAGES_ID_MAX, tasks_info, messages_info, NULL,
        NULL);

    // initialize mme config
    amf_config_init(&amf_config);
    amf_nas_state_init(&amf_config);
    create_state_matrix();

    init_task_context(TASK_MAIN, nullptr, 0, NULL, &amf_app_task_zmq_ctx);

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

bool validate_auth_procedure(
    amf_ue_ngap_id_t ue_id, uint32_t expected_retransmission_count) {
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

bool validate_smc_procedure(
    amf_ue_ngap_id_t ue_id, uint32_t expected_retransmission_count) {
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
  imsi64          = send_initial_ue_message_no_tmsi(
      amf_app_desc_p, 36, 1, 1, 0, plmn, initial_ue_message_hexbuf,
      sizeof(initial_ue_message_hexbuf));

  /* Check if UE Context is created with correct imsi */
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id));

  /* Send the authentication response message from subscriberdb */
  int rc = RETURNok;
  rc     = send_proc_authentication_info_answer(imsi, ue_id, true);
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

  /* Send uplinkg nas message for security mode complete response from UE */
  rc = send_uplink_nas_message_ue_smc_response(
      amf_app_desc_p, ue_id, plmn, ue_smc_response_hexbuf,
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
  int rc                 = RETURNerror;
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
  imsi64          = send_initial_ue_message_no_tmsi(
      amf_app_desc_p, 36, 1, 1, 0, plmn, initial_ue_message_hexbuf,
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

  /* Send uplinkg nas message for security mode complete response from UE */
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

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_TRUE(rc == RETURNok);
  EXPECT_TRUE(expected_Ids == AMFClientServicer::getInstance().msgtype_stack);
}

TEST_F(AMFAppProcedureTest, TestPDUSessionSetup) {
  int rc                 = RETURNerror;
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
  imsi64          = send_initial_ue_message_no_tmsi(
      amf_app_desc_p, 36, 1, 1, 0, plmn, initial_ue_message_hexbuf,
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
  int rc                 = RETURNerror;
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
  imsi64          = send_initial_ue_message_no_tmsi(
      amf_app_desc_p, 36, 1, 1, 0, plmn, initial_ue_message_hexbuf,
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
  int rc                 = RETURNerror;
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
  imsi64          = send_initial_ue_message_no_tmsi(
      amf_app_desc_p, 36, 1, 1, 0, plmn, initial_ue_message_hexbuf,
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

  ue_m5gmm_context_t* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  // ue context should exist
  ASSERT_NE(ue_context_p, nullptr);
  memset(
      &ue_context_p->amf_context.apn_config_profile, 0,
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
  int rc                 = RETURNerror;
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
  imsi64          = send_initial_ue_message_no_tmsi(
      amf_app_desc_p, 36, 1, 1, 0, plmn, initial_ue_message_hexbuf,
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
  int rc                 = RETURNerror;
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
  imsi64          = send_initial_ue_message_no_tmsi(
      amf_app_desc_p, 36, 1, 1, 0, plmn, initial_ue_message_hexbuf,
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

}  // namespace magma5g
