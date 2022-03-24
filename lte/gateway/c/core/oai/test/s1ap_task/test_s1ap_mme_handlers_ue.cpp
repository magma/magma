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
#include <gtest/gtest.h>
#include <thread>

#include "lte/gateway/c/core/oai/test/mock_tasks/mock_tasks.h"

extern "C" {
#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/include/mme_config.h"
#include "lte/gateway/c/core/oai/include/s1ap_state.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme_decoder.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme_handlers.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme_nas_procedures.h"
}

#include "lte/gateway/c/core/oai/test/s1ap_task/s1ap_mme_test_utils.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_state_manager.h"
#include "lte/gateway/c/core/oai/test/s1ap_task/s1ap_mme_handlers_test_fixture.h"

extern bool hss_associated;

namespace magma {
namespace lte {

TEST_F(S1apMmeHandlersTest, HandleICSResponseICSRelease) {
  ASSERT_EQ(task_zmq_ctx_main_s1ap.ready, true);

  bool is_state_same = false;

  EXPECT_CALL(*sctp_handler, sctpd_send_dl()).Times(2);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_initial_ue_message()).Times(1);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_s1ap_ue_context_release_req())
      .Times(1);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_initial_context_setup_failure())
      .Times(0);

  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_INIT, 0));

  S1ap_S1AP_PDU_t pdu_s1;
  memset(&pdu_s1, 0, sizeof(pdu_s1));
  ASSERT_EQ(RETURNok, generate_s1_setup_request_pdu(&pdu_s1));
  ASSERT_EQ(RETURNok,
            s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_s1));

  // State validation
  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_READY, 0));
  ASSERT_TRUE(is_num_enbs_valid(state, 1));

  uint8_t initial_ue_bytes[] = {
      0x00, 0x0c, 0x40, 0x48, 0x00, 0x00, 0x05, 0x00, 0x08, 0x00, 0x02,
      0x00, 0x01, 0x00, 0x1a, 0x00, 0x20, 0x1f, 0x07, 0x41, 0x71, 0x08,
      0x09, 0x10, 0x10, 0x00, 0x00, 0x00, 0x00, 0x10, 0x02, 0xe0, 0xe0,
      0x00, 0x04, 0x02, 0x01, 0xd0, 0x11, 0x40, 0x08, 0x04, 0x02, 0x60,
      0x04, 0x00, 0x02, 0x1c, 0x00, 0x00, 0x43, 0x00, 0x06, 0x00, 0x00,
      0xf1, 0x10, 0x00, 0x01, 0x00, 0x64, 0x40, 0x08, 0x00, 0x00, 0xf1,
      0x10, 0x00, 0x00, 0x00, 0xa0, 0x00, 0x86, 0x40, 0x01, 0x30};

  ASSERT_EQ(simulate_pdu_s1_message(initial_ue_bytes, sizeof(initial_ue_bytes),
                                    state, assoc_id, stream_id),
            RETURNok);

  handle_mme_ue_id_notification(state, assoc_id);

  // Generate downlink nas transport with dummy payload
  bstring p;
  std::string test_str = "test";
  STRING_TO_BSTRING(test_str, p);
  s1ap_generate_downlink_nas_transport(state, 1, 7, &p, 1, &is_state_same);
  bdestroy_wrapper(&p);

  // Authentication response proc packet bytes
  uint8_t auth_bytes[] = {
      0x00, 0x0d, 0x40, 0x3d, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x02,
      0x00, 0x07, 0x00, 0x08, 0x00, 0x02, 0x00, 0x01, 0x00, 0x1a, 0x00,
      0x14, 0x13, 0x07, 0x53, 0x10, 0x1e, 0x63, 0x7e, 0x5c, 0x58, 0xec,
      0x5a, 0xa8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
      0x64, 0x40, 0x08, 0x00, 0x00, 0xf1, 0x10, 0x00, 0x00, 0x00, 0xa0,
      0x00, 0x43, 0x40, 0x06, 0x00, 0x00, 0xf1, 0x10, 0x00, 0x01};

  ASSERT_EQ(simulate_pdu_s1_message(auth_bytes, sizeof(auth_bytes), state,
                                    assoc_id, stream_id),
            RETURNok);

  uint8_t ics_bytes[] = {0x20, 0x09, 0x00, 0x22, 0x00, 0x00, 0x03, 0x00,
                         0x00, 0x40, 0x02, 0x00, 0x07, 0x00, 0x08, 0x40,
                         0x02, 0x00, 0x01, 0x00, 0x33, 0x40, 0x0f, 0x00,
                         0x00, 0x32, 0x40, 0x0a, 0x0a, 0x1f, 0xc0, 0xa8,
                         0x3c, 0x8d, 0x0a, 0x00, 0x01, 0x28};

  ASSERT_EQ(simulate_pdu_s1_message(ics_bytes, sizeof(ics_bytes), state,
                                    assoc_id, stream_id),
            RETURNok);

  // Freeing pdu and payload data
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_s1);

  // State validation
  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_READY, 1));
  ASSERT_TRUE(is_num_enbs_valid(state, 1));

  uint8_t ics_release_bytes[] = {0x00, 0x12, 0x40, 0x15, 0x00, 0x00, 0x03,
                                 0x00, 0x00, 0x00, 0x02, 0x00, 0x07, 0x00,
                                 0x08, 0x00, 0x02, 0x00, 0x01, 0x00, 0x02,
                                 0x40, 0x02, 0x02, 0x80};

  ASSERT_EQ(
      simulate_pdu_s1_message(ics_release_bytes, sizeof(ics_release_bytes),
                              state, assoc_id, stream_id),
      RETURNok);

  // State validation
  ASSERT_TRUE(is_num_enbs_valid(state, 1));
}

TEST_F(S1apMmeHandlersTest, HandleICSFailure) {
  ASSERT_EQ(task_zmq_ctx_main_s1ap.ready, true);

  bool is_state_same = false;

  EXPECT_CALL(*sctp_handler, sctpd_send_dl()).Times(2);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_initial_ue_message()).Times(1);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_s1ap_ue_context_release_req())
      .Times(0);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_initial_context_setup_failure())
      .Times(1);

  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_INIT, 0));

  S1ap_S1AP_PDU_t pdu_s1;
  memset(&pdu_s1, 0, sizeof(pdu_s1));
  ASSERT_EQ(RETURNok, generate_s1_setup_request_pdu(&pdu_s1));
  ASSERT_EQ(RETURNok,
            s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_s1));

  // State validation
  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_READY, 0));
  ASSERT_TRUE(is_num_enbs_valid(state, 1));

  uint8_t initial_ue_bytes[] = {
      0x00, 0x0c, 0x40, 0x48, 0x00, 0x00, 0x05, 0x00, 0x08, 0x00, 0x02,
      0x00, 0x01, 0x00, 0x1a, 0x00, 0x20, 0x1f, 0x07, 0x41, 0x71, 0x08,
      0x09, 0x10, 0x10, 0x00, 0x00, 0x00, 0x00, 0x10, 0x02, 0xe0, 0xe0,
      0x00, 0x04, 0x02, 0x01, 0xd0, 0x11, 0x40, 0x08, 0x04, 0x02, 0x60,
      0x04, 0x00, 0x02, 0x1c, 0x00, 0x00, 0x43, 0x00, 0x06, 0x00, 0x00,
      0xf1, 0x10, 0x00, 0x01, 0x00, 0x64, 0x40, 0x08, 0x00, 0x00, 0xf1,
      0x10, 0x00, 0x00, 0x00, 0xa0, 0x00, 0x86, 0x40, 0x01, 0x30};

  ASSERT_EQ(simulate_pdu_s1_message(initial_ue_bytes, sizeof(initial_ue_bytes),
                                    state, assoc_id, stream_id),
            RETURNok);

  handle_mme_ue_id_notification(state, assoc_id);

  // Generate downlink nas transport with dummy payload
  bstring p;
  std::string test_str = "test";
  STRING_TO_BSTRING(test_str, p);
  s1ap_generate_downlink_nas_transport(state, 1, 7, &p, 1, &is_state_same);
  bdestroy_wrapper(&p);

  // Simulate Authentication Rsp
  uint8_t auth_bytes[] = {
      0x00, 0x0d, 0x40, 0x3d, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x02,
      0x00, 0x07, 0x00, 0x08, 0x00, 0x02, 0x00, 0x01, 0x00, 0x1a, 0x00,
      0x14, 0x13, 0x07, 0x53, 0x10, 0x1e, 0x63, 0x7e, 0x5c, 0x58, 0xec,
      0x5a, 0xa8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
      0x64, 0x40, 0x08, 0x00, 0x00, 0xf1, 0x10, 0x00, 0x00, 0x00, 0xa0,
      0x00, 0x43, 0x40, 0x06, 0x00, 0x00, 0xf1, 0x10, 0x00, 0x01};

  ASSERT_EQ(simulate_pdu_s1_message(auth_bytes, sizeof(auth_bytes), state,
                                    assoc_id, stream_id),
            RETURNok);

  // Simulate ICS Failure
  uint8_t ics_fail[] = {0x40, 0x09, 0x00, 0x15, 0x00, 0x00, 0x03, 0x00, 0x00,
                        0x40, 0x02, 0x00, 0x07, 0x00, 0x08, 0x40, 0x02, 0x00,
                        0x01, 0x00, 0x02, 0x40, 0x02, 0x00, 0x00};

  ASSERT_EQ(simulate_pdu_s1_message(ics_fail, sizeof(ics_fail), state, assoc_id,
                                    stream_id),
            RETURNok);

  // Freeing pdu and payload data
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_s1);
}

TEST_F(S1apMmeHandlersTest, HandleUECapIndication) {
  ASSERT_EQ(task_zmq_ctx_main_s1ap.ready, true);

  EXPECT_CALL(*sctp_handler, sctpd_send_dl()).Times(1);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_initial_ue_message()).Times(1);

  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_INIT, 0));

  S1ap_S1AP_PDU_t pdu_s1;
  memset(&pdu_s1, 0, sizeof(pdu_s1));
  ASSERT_EQ(RETURNok, generate_s1_setup_request_pdu(&pdu_s1));
  ASSERT_EQ(RETURNok,
            s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_s1));

  // State validation
  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_READY, 0));
  ASSERT_TRUE(is_num_enbs_valid(state, 1));

  uint8_t initial_ue_bytes[] = {
      0x00, 0x0c, 0x40, 0x48, 0x00, 0x00, 0x05, 0x00, 0x08, 0x00, 0x02,
      0x00, 0x01, 0x00, 0x1a, 0x00, 0x20, 0x1f, 0x07, 0x41, 0x71, 0x08,
      0x09, 0x10, 0x10, 0x00, 0x00, 0x00, 0x00, 0x10, 0x02, 0xe0, 0xe0,
      0x00, 0x04, 0x02, 0x01, 0xd0, 0x11, 0x40, 0x08, 0x04, 0x02, 0x60,
      0x04, 0x00, 0x02, 0x1c, 0x00, 0x00, 0x43, 0x00, 0x06, 0x00, 0x00,
      0xf1, 0x10, 0x00, 0x01, 0x00, 0x64, 0x40, 0x08, 0x00, 0x00, 0xf1,
      0x10, 0x00, 0x00, 0x00, 0xa0, 0x00, 0x86, 0x40, 0x01, 0x30};

  ASSERT_EQ(simulate_pdu_s1_message(initial_ue_bytes, sizeof(initial_ue_bytes),
                                    state, assoc_id, stream_id),
            RETURNok);

  handle_mme_ue_id_notification(state, assoc_id);

  uint8_t ue_cap_bytes[] = {
      0x00, 0x16, 0x40, 0x53, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x02,
      0x00, 0x07, 0x00, 0x08, 0x00, 0x02, 0x00, 0x01, 0x00, 0x4a, 0x40,
      0x40, 0x3f, 0x01, 0xe8, 0x01, 0x03, 0xac, 0x59, 0x80, 0x07, 0x00,
      0x08, 0x20, 0x81, 0x83, 0x9b, 0x4e, 0x1c, 0x3f, 0xf8, 0x7f, 0xf0,
      0xff, 0xe1, 0xff, 0xc3, 0xff, 0x87, 0xff, 0x0f, 0xfe, 0x1f, 0xfd,
      0xf8, 0x37, 0x62, 0x78, 0x00, 0xa0, 0x18, 0x5f, 0x80, 0x00, 0x00,
      0x00, 0x1c, 0x07, 0xe0, 0xdd, 0x89, 0xe0, 0x00, 0x00, 0x00, 0x07,
      0x09, 0xf8, 0x37, 0x62, 0x78, 0x00, 0x00, 0x00, 0x00, 0x00};

  bstring payload_ue_cap;
  payload_ue_cap = blk2bstr(&ue_cap_bytes, sizeof(ue_cap_bytes));
  S1ap_S1AP_PDU_t pdu_cap;
  memset(&pdu_cap, 0, sizeof(pdu_cap));

  ASSERT_EQ(s1ap_mme_decode_pdu(&pdu_cap, payload_ue_cap), RETURNok);
  ASSERT_EQ(s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_cap),
            RETURNok);

  // State validation
  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_READY, 1));
  ASSERT_TRUE(is_num_enbs_valid(state, 1));

  // Freeing pdu and payload data
  bdestroy_wrapper(&payload_ue_cap);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_cap);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_s1);
}

TEST_F(S1apMmeHandlersTest, GenerateUEContextReleaseCommand) {
  ue_description_t ue_ref_p = {
      .enb_ue_s1ap_id = 1,
      .mme_ue_s1ap_id = 1,
      .sctp_assoc_id = assoc_id,
      .comp_s1ap_id = S1AP_GENERATE_COMP_S1AP_ID(assoc_id, 1)};

  ue_ref_p.s1ap_ue_context_rel_timer.id = -1;
  ue_ref_p.s1ap_ue_context_rel_timer.msec = 1000;
  EXPECT_CALL(*sctp_handler, sctpd_send_dl()).Times(2);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_initial_ue_message()).Times(1);

  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_INIT, 0));

  S1ap_S1AP_PDU_t pdu_s1;
  memset(&pdu_s1, 0, sizeof(pdu_s1));
  ASSERT_EQ(RETURNok, generate_s1_setup_request_pdu(&pdu_s1));
  ASSERT_EQ(RETURNok,
            s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_s1));

  // State validation
  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_READY, 0));
  ASSERT_TRUE(is_num_enbs_valid(state, 1));

  uint8_t initial_ue_bytes[] = {
      0x00, 0x0c, 0x40, 0x48, 0x00, 0x00, 0x05, 0x00, 0x08, 0x00, 0x02,
      0x00, 0x01, 0x00, 0x1a, 0x00, 0x20, 0x1f, 0x07, 0x41, 0x71, 0x08,
      0x09, 0x10, 0x10, 0x00, 0x00, 0x00, 0x00, 0x10, 0x02, 0xe0, 0xe0,
      0x00, 0x04, 0x02, 0x01, 0xd0, 0x11, 0x40, 0x08, 0x04, 0x02, 0x60,
      0x04, 0x00, 0x02, 0x1c, 0x00, 0x00, 0x43, 0x00, 0x06, 0x00, 0x00,
      0xf1, 0x10, 0x00, 0x01, 0x00, 0x64, 0x40, 0x08, 0x00, 0x00, 0xf1,
      0x10, 0x00, 0x00, 0x00, 0xa0, 0x00, 0x86, 0x40, 0x01, 0x30};

  ASSERT_EQ(simulate_pdu_s1_message(initial_ue_bytes, sizeof(initial_ue_bytes),
                                    state, assoc_id, stream_id),
            RETURNok);

  // Invalid S1 Cause returns error
  ASSERT_EQ(RETURNerror, s1ap_mme_generate_ue_context_release_command(
                             state, &ue_ref_p, S1AP_IMPLICIT_CONTEXT_RELEASE,
                             INVALID_IMSI64, assoc_id, stream_id, 1, 1));
  // Valid S1 Causes passess successfully
  ASSERT_EQ(RETURNok, s1ap_mme_generate_ue_context_release_command(
                          state, &ue_ref_p, S1AP_INITIAL_CONTEXT_SETUP_FAILED,
                          INVALID_IMSI64, assoc_id, stream_id, 1, 1));

  EXPECT_NE(ue_ref_p.s1ap_ue_context_rel_timer.id, S1AP_TIMER_INACTIVE_ID);

  // Freeing pdu and payload data
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_s1);
}

TEST_F(S1apMmeHandlersTest, HandleUEContextRelease) {
  ASSERT_EQ(task_zmq_ctx_main_s1ap.ready, true);

  EXPECT_CALL(*sctp_handler, sctpd_send_dl()).Times(1);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_initial_ue_message()).Times(1);

  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_INIT, 0));

  S1ap_S1AP_PDU_t pdu_s1;
  memset(&pdu_s1, 0, sizeof(pdu_s1));
  ASSERT_EQ(RETURNok, generate_s1_setup_request_pdu(&pdu_s1));
  ASSERT_EQ(RETURNok,
            s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_s1));

  uint8_t initial_ue_bytes[] = {
      0x00, 0x0c, 0x40, 0x48, 0x00, 0x00, 0x05, 0x00, 0x08, 0x00, 0x02,
      0x00, 0x01, 0x00, 0x1a, 0x00, 0x20, 0x1f, 0x07, 0x41, 0x71, 0x08,
      0x09, 0x10, 0x10, 0x00, 0x00, 0x00, 0x00, 0x10, 0x02, 0xe0, 0xe0,
      0x00, 0x04, 0x02, 0x01, 0xd0, 0x11, 0x40, 0x08, 0x04, 0x02, 0x60,
      0x04, 0x00, 0x02, 0x1c, 0x00, 0x00, 0x43, 0x00, 0x06, 0x00, 0x00,
      0xf1, 0x10, 0x00, 0x01, 0x00, 0x64, 0x40, 0x08, 0x00, 0x00, 0xf1,
      0x10, 0x00, 0x00, 0x00, 0xa0, 0x00, 0x86, 0x40, 0x01, 0x30};

  ASSERT_EQ(simulate_pdu_s1_message(initial_ue_bytes, sizeof(initial_ue_bytes),
                                    state, assoc_id, stream_id),
            RETURNok);

  handle_mme_ue_id_notification(state, assoc_id);

  // State validation
  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_READY, 1));
  ASSERT_TRUE(is_num_enbs_valid(state, 1));
  ASSERT_EQ(state->mmeid2associd.num_elements, 1);

  // Send UE context release command mimicing MME_APP
  MessageDef* message_p;
  message_p =
      itti_alloc_new_message(TASK_MME_APP, S1AP_UE_CONTEXT_RELEASE_COMMAND);
  S1AP_UE_CONTEXT_RELEASE_COMMAND(message_p).mme_ue_s1ap_id = 7;
  S1AP_UE_CONTEXT_RELEASE_COMMAND(message_p).enb_ue_s1ap_id = 1;
  S1AP_UE_CONTEXT_RELEASE_COMMAND(message_p).cause =
      S1AP_SCTP_SHUTDOWN_OR_RESET;
  ASSERT_EQ(send_msg_to_task(&task_zmq_ctx_main_s1ap, TASK_S1AP, message_p),
            RETURNok);

  std::this_thread::sleep_for(std::chrono::milliseconds(500));
  ASSERT_TRUE(is_num_enbs_valid(state, 1));

  uint8_t rel_comp_bytes[] = {0x20, 0x17, 0x00, 0x0f, 0x00, 0x00, 0x02,
                              0x00, 0x00, 0x40, 0x02, 0x00, 0x07, 0x00,
                              0x08, 0x40, 0x02, 0x00, 0x01};

  bstring payload_rel;
  payload_rel = blk2bstr(&rel_comp_bytes, sizeof(rel_comp_bytes));
  S1ap_S1AP_PDU_t pdu_rel;
  memset(&pdu_rel, 0, sizeof(pdu_rel));

  ASSERT_EQ(RETURNok, s1ap_mme_decode_pdu(&pdu_rel, payload_rel));
  ASSERT_EQ(RETURNok,
            s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_rel));

  // Freeing pdu and payload data
  bdestroy_wrapper(&payload_rel);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_rel);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_s1);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(500));

  ASSERT_EQ(state->mmeid2associd.num_elements, 0);
}

TEST_F(S1apMmeHandlersTest, HandleConnectionEstCnf) {
  ASSERT_EQ(task_zmq_ctx_main_s1ap.ready, true);
  itti_mme_app_connection_establishment_cnf_t* establishment_cnf_p = NULL;

  EXPECT_CALL(*sctp_handler, sctpd_send_dl()).Times(2);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_initial_ue_message()).Times(1);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_s1ap_ue_context_release_req())
      .Times(0);
  EXPECT_CALL(*mme_app_handler, nas_proc_dl_transfer_rej()).Times(0);

  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_INIT, 0));

  S1ap_S1AP_PDU_t pdu_s1;
  memset(&pdu_s1, 0, sizeof(pdu_s1));
  ASSERT_EQ(RETURNok, generate_s1_setup_request_pdu(&pdu_s1));
  ASSERT_EQ(RETURNok,
            s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_s1));

  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_READY, 0));

  uint8_t initial_ue_bytes[] = {
      0x00, 0x0c, 0x40, 0x48, 0x00, 0x00, 0x05, 0x00, 0x08, 0x00, 0x02,
      0x00, 0x01, 0x00, 0x1a, 0x00, 0x20, 0x1f, 0x07, 0x41, 0x71, 0x08,
      0x09, 0x10, 0x10, 0x00, 0x00, 0x00, 0x00, 0x10, 0x02, 0xe0, 0xe0,
      0x00, 0x04, 0x02, 0x01, 0xd0, 0x11, 0x40, 0x08, 0x04, 0x02, 0x60,
      0x04, 0x00, 0x02, 0x1c, 0x00, 0x00, 0x43, 0x00, 0x06, 0x00, 0x00,
      0xf1, 0x10, 0x00, 0x01, 0x00, 0x64, 0x40, 0x08, 0x00, 0x00, 0xf1,
      0x10, 0x00, 0x00, 0x00, 0xa0, 0x00, 0x86, 0x40, 0x01, 0x30};

  ASSERT_EQ(simulate_pdu_s1_message(initial_ue_bytes, sizeof(initial_ue_bytes),
                                    state, assoc_id, stream_id),
            RETURNok);

  handle_mme_ue_id_notification(state, assoc_id);

  ASSERT_EQ(state->mmeid2associd.num_elements, 1);

  // Send UE connection establishment cnf mimicing MME_APP

  ASSERT_EQ(send_conn_establishment_cnf(7, false, true, true), RETURNok);

  // Freeing pdu and payload data
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_s1);
}

TEST_F(S1apMmeHandlersTest, HandleConnectionEstCnfExtUEAMBR) {
  ASSERT_EQ(task_zmq_ctx_main_s1ap.ready, true);

  EXPECT_CALL(*sctp_handler, sctpd_send_dl()).Times(2);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_initial_ue_message()).Times(1);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_s1ap_ue_context_release_req())
      .Times(0);
  EXPECT_CALL(*mme_app_handler, nas_proc_dl_transfer_rej()).Times(0);

  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_INIT, 0));

  S1ap_S1AP_PDU_t pdu_s1;
  memset(&pdu_s1, 0, sizeof(pdu_s1));
  ASSERT_EQ(RETURNok, generate_s1_setup_request_pdu(&pdu_s1));
  ASSERT_EQ(RETURNok,
            s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_s1));

  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_READY, 0));

  uint8_t initial_ue_bytes[] = {
      0x00, 0x0c, 0x40, 0x48, 0x00, 0x00, 0x05, 0x00, 0x08, 0x00, 0x02,
      0x00, 0x01, 0x00, 0x1a, 0x00, 0x20, 0x1f, 0x07, 0x41, 0x71, 0x08,
      0x09, 0x10, 0x10, 0x00, 0x00, 0x00, 0x00, 0x10, 0x02, 0xe0, 0xe0,
      0x00, 0x04, 0x02, 0x01, 0xd0, 0x11, 0x40, 0x08, 0x04, 0x02, 0x60,
      0x04, 0x00, 0x02, 0x1c, 0x00, 0x00, 0x43, 0x00, 0x06, 0x00, 0x00,
      0xf1, 0x10, 0x00, 0x01, 0x00, 0x64, 0x40, 0x08, 0x00, 0x00, 0xf1,
      0x10, 0x00, 0x00, 0x00, 0xa0, 0x00, 0x86, 0x40, 0x01, 0x30};

  ASSERT_EQ(simulate_pdu_s1_message(initial_ue_bytes, sizeof(initial_ue_bytes),
                                    state, assoc_id, stream_id),
            RETURNok);

  handle_mme_ue_id_notification(state, assoc_id);

  ASSERT_EQ(state->mmeid2associd.num_elements, 1);

  // Send UE connection establishment cnf mimicing MME_APP

  ASSERT_EQ(send_conn_establishment_cnf(7, true, true, true), RETURNok);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_s1);
}

TEST_F(S1apMmeHandlersTest, HandleS1apUeCtxtModification) {
  ASSERT_EQ(task_zmq_ctx_main_s1ap.ready, true);

  EXPECT_CALL(*sctp_handler, sctpd_send_dl()).Times(2);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_initial_ue_message()).Times(1);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_s1ap_ue_context_release_req())
      .Times(0);
  EXPECT_CALL(*mme_app_handler, nas_proc_dl_transfer_rej()).Times(0);

  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_INIT, 0));

  S1ap_S1AP_PDU_t pdu_s1;
  memset(&pdu_s1, 0, sizeof(pdu_s1));
  ASSERT_EQ(RETURNok, generate_s1_setup_request_pdu(&pdu_s1));
  ASSERT_EQ(RETURNok,
            s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_s1));

  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_READY, 0));

  uint8_t initial_ue_bytes[] = {
      0x00, 0x0c, 0x40, 0x48, 0x00, 0x00, 0x05, 0x00, 0x08, 0x00, 0x02,
      0x00, 0x01, 0x00, 0x1a, 0x00, 0x20, 0x1f, 0x07, 0x41, 0x71, 0x08,
      0x09, 0x10, 0x10, 0x00, 0x00, 0x00, 0x00, 0x10, 0x02, 0xe0, 0xe0,
      0x00, 0x04, 0x02, 0x01, 0xd0, 0x11, 0x40, 0x08, 0x04, 0x02, 0x60,
      0x04, 0x00, 0x02, 0x1c, 0x00, 0x00, 0x43, 0x00, 0x06, 0x00, 0x00,
      0xf1, 0x10, 0x00, 0x01, 0x00, 0x64, 0x40, 0x08, 0x00, 0x00, 0xf1,
      0x10, 0x00, 0x00, 0x00, 0xa0, 0x00, 0x86, 0x40, 0x01, 0x30};

  ASSERT_EQ(simulate_pdu_s1_message(initial_ue_bytes, sizeof(initial_ue_bytes),
                                    state, assoc_id, stream_id),
            RETURNok);

  handle_mme_ue_id_notification(state, assoc_id);

  // State validation
  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_READY, 1));
  ASSERT_TRUE(is_num_enbs_valid(state, 1));
  ASSERT_EQ(state->mmeid2associd.num_elements, 1);

  // Send S1AP_UE_CONTEXT_MODIFICATION_REQUEST mimicing MME_APP
  ASSERT_EQ(send_s1ap_ue_ctxt_mod(1, 7), RETURNok);

  // Freeing pdu and payload data
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_s1);
}

TEST_F(S1apMmeHandlersTest, HandleErrorIndicationMessage) {
  ASSERT_EQ(task_zmq_ctx_main_s1ap.ready, true);

  bool is_state_same = true;

  EXPECT_CALL(*sctp_handler, sctpd_send_dl()).Times(2);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_initial_ue_message()).Times(1);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_s1ap_ue_context_release_req())
      .Times(1);
  EXPECT_CALL(*mme_app_handler, nas_proc_dl_transfer_rej()).Times(0);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_enb_reset_req()).Times(0);

  // Simulate S1Setup
  uint8_t s1_bytes[] = {0x00, 0x11, 0x00, 0x2f, 0x00, 0x00, 0x04, 0x00, 0x3b,
                        0x00, 0x09, 0x00, 0x00, 0xf1, 0x10, 0x40, 0x00, 0x00,
                        0x00, 0x10, 0x00, 0x3c, 0x40, 0x0b, 0x80, 0x09, 0x22,
                        0x52, 0x41, 0x44, 0x49, 0x53, 0x59, 0x53, 0x22, 0x00,
                        0x40, 0x00, 0x07, 0x00, 0x00, 0x00, 0x40, 0x00, 0xf1,
                        0x10, 0x00, 0x89, 0x40, 0x01, 0x00};

  ASSERT_EQ(simulate_pdu_s1_message(s1_bytes, sizeof(s1_bytes), state, assoc_id,
                                    stream_id),
            RETURNok);

  // Simulate InitialUEMessage - Attach Request
  uint8_t initial_ue_bytes[] = {
      0x00, 0x0c, 0x40, 0x48, 0x00, 0x00, 0x05, 0x00, 0x08, 0x00, 0x02,
      0x00, 0x01, 0x00, 0x1a, 0x00, 0x20, 0x1f, 0x07, 0x41, 0x71, 0x08,
      0x09, 0x10, 0x10, 0x00, 0x00, 0x00, 0x00, 0x10, 0x02, 0xe0, 0xe0,
      0x00, 0x04, 0x02, 0x01, 0xd0, 0x11, 0x40, 0x08, 0x04, 0x02, 0x60,
      0x04, 0x00, 0x02, 0x1c, 0x00, 0x00, 0x43, 0x00, 0x06, 0x00, 0x00,
      0xf1, 0x10, 0x00, 0x01, 0x00, 0x64, 0x40, 0x08, 0x00, 0x00, 0xf1,
      0x10, 0x00, 0x00, 0x00, 0xa0, 0x00, 0x86, 0x40, 0x01, 0x30};

  ASSERT_EQ(simulate_pdu_s1_message(initial_ue_bytes, sizeof(initial_ue_bytes),
                                    state, assoc_id, stream_id),
            RETURNok);

  handle_mme_ue_id_notification(state, assoc_id);

  // Simulate UECapabilityInfoIndication
  uint8_t ue_cap_bytes[] = {
      0x00, 0x16, 0x40, 0x53, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x02,
      0x00, 0x07, 0x00, 0x08, 0x00, 0x02, 0x00, 0x01, 0x00, 0x4a, 0x40,
      0x40, 0x3f, 0x01, 0xe8, 0x01, 0x03, 0xac, 0x59, 0x80, 0x07, 0x00,
      0x08, 0x20, 0x81, 0x83, 0x9b, 0x4e, 0x1c, 0x3f, 0xf8, 0x7f, 0xf0,
      0xff, 0xe1, 0xff, 0xc3, 0xff, 0x87, 0xff, 0x0f, 0xfe, 0x1f, 0xfd,
      0xf8, 0x37, 0x62, 0x78, 0x00, 0xa0, 0x18, 0x5f, 0x80, 0x00, 0x00,
      0x00, 0x1c, 0x07, 0xe0, 0xdd, 0x89, 0xe0, 0x00, 0x00, 0x00, 0x07,
      0x09, 0xf8, 0x37, 0x62, 0x78, 0x00, 0x00, 0x00, 0x00, 0x00};

  ASSERT_EQ(simulate_pdu_s1_message(ue_cap_bytes, sizeof(ue_cap_bytes), state,
                                    assoc_id, stream_id),
            RETURNok);

  // Generate downlink nas transport with dummy payload
  bstring p;
  std::string test_str = "test";
  STRING_TO_BSTRING(test_str, p);
  s1ap_generate_downlink_nas_transport(state, 1, 7, &p, 1, &is_state_same);
  bdestroy_wrapper(&p);

  // Authentication response proc packet bytes
  uint8_t auth_bytes[] = {
      0x00, 0x0d, 0x40, 0x3d, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x02,
      0x00, 0x07, 0x00, 0x08, 0x00, 0x02, 0x00, 0x01, 0x00, 0x1a, 0x00,
      0x14, 0x13, 0x07, 0x53, 0x10, 0x1e, 0x63, 0x7e, 0x5c, 0x58, 0xec,
      0x5a, 0xa8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
      0x64, 0x40, 0x08, 0x00, 0x00, 0xf1, 0x10, 0x00, 0x00, 0x00, 0xa0,
      0x00, 0x43, 0x40, 0x06, 0x00, 0x00, 0xf1, 0x10, 0x00, 0x01};

  ASSERT_EQ(simulate_pdu_s1_message(auth_bytes, sizeof(auth_bytes), state,
                                    assoc_id, stream_id),
            RETURNok);

  // Simulate InitialContextSetup
  uint8_t ics_bytes[] = {0x20, 0x09, 0x00, 0x22, 0x00, 0x00, 0x03, 0x00,
                         0x00, 0x40, 0x02, 0x00, 0x07, 0x00, 0x08, 0x40,
                         0x02, 0x00, 0x01, 0x00, 0x33, 0x40, 0x0f, 0x00,
                         0x00, 0x32, 0x40, 0x0a, 0x0a, 0x1f, 0xc0, 0xa8,
                         0x3c, 0x8d, 0x0a, 0x00, 0x01, 0x28};

  ASSERT_EQ(simulate_pdu_s1_message(ics_bytes, sizeof(ics_bytes), state,
                                    assoc_id, stream_id),
            RETURNok);

  ASSERT_TRUE(is_ue_state_valid(assoc_id, 1, S1AP_UE_CONNECTED));

  // Simulate Attach Complete
  uint8_t attach_compl_bytes[] = {
      0x00, 0x0d, 0x40, 0x37, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x02, 0x00,
      0x07, 0x00, 0x08, 0x00, 0x02, 0x00, 0x01, 0x00, 0x1a, 0x00, 0x0e, 0x0d,
      0x27, 0xeb, 0x9e, 0x7f, 0x7e, 0x01, 0x07, 0x43, 0x00, 0x03, 0x52, 0x00,
      0xc2, 0x00, 0x64, 0x40, 0x08, 0x00, 0x00, 0xf1, 0x10, 0x00, 0x00, 0x00,
      0x10, 0x00, 0x43, 0x40, 0x06, 0x00, 0x00, 0xf1, 0x10, 0x00, 0x01};

  ASSERT_EQ(
      simulate_pdu_s1_message(attach_compl_bytes, sizeof(attach_compl_bytes),
                              state, assoc_id, stream_id),
      RETURNok);

  // Simulate Uplink NAS Transport
  uint8_t uplink_nas_bytes[] = {
      0x00, 0x0d, 0x40, 0x33, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x02,
      0x00, 0x07, 0x00, 0x08, 0x00, 0x02, 0x00, 0x01, 0x00, 0x1a, 0x00,
      0x0a, 0x09, 0x27, 0xca, 0x02, 0x76, 0x29, 0x02, 0x62, 0x00, 0xc6,
      0x00, 0x64, 0x40, 0x08, 0x00, 0x00, 0xf1, 0x10, 0x00, 0x00, 0x00,
      0xa0, 0x00, 0x43, 0x40, 0x06, 0x00, 0x00, 0xf1, 0x10, 0x00, 0x01};

  ASSERT_EQ(simulate_pdu_s1_message(uplink_nas_bytes, sizeof(uplink_nas_bytes),
                                    state, assoc_id, stream_id),
            RETURNok);

  // Simulate Error Ind Message
  uint8_t error_ind_bytes[] = {0x00, 0x0f, 0x40, 0x15, 0x00, 0x00, 0x03,
                               0x00, 0x00, 0x40, 0x02, 0x00, 0x07, 0x00,
                               0x08, 0x40, 0x02, 0x00, 0x01, 0x00, 0x02,
                               0x40, 0x02, 0x01, 0xe0};

  ASSERT_EQ(simulate_pdu_s1_message(error_ind_bytes, sizeof(error_ind_bytes),
                                    state, assoc_id, stream_id),
            RETURNok);
}

TEST_F(S1apMmeHandlersTest, HandleS1apPagingRequest) {
  ASSERT_EQ(task_zmq_ctx_main_s1ap.ready, true);

  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_INIT, 0));

  EXPECT_CALL(*sctp_handler, sctpd_send_dl()).Times(1);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_initial_ue_message()).Times(1);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_s1ap_ue_context_release_req())
      .Times(0);
  EXPECT_CALL(*mme_app_handler, nas_proc_dl_transfer_rej()).Times(0);

  S1ap_S1AP_PDU_t pdu_s1;
  memset(&pdu_s1, 0, sizeof(pdu_s1));
  ASSERT_EQ(RETURNok, generate_s1_setup_request_pdu(&pdu_s1));
  ASSERT_EQ(RETURNok,
            s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_s1));

  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_READY, 0));

  uint8_t initial_ue_bytes[] = {
      0x00, 0x0c, 0x40, 0x48, 0x00, 0x00, 0x05, 0x00, 0x08, 0x00, 0x02,
      0x00, 0x01, 0x00, 0x1a, 0x00, 0x20, 0x1f, 0x07, 0x41, 0x71, 0x08,
      0x09, 0x10, 0x10, 0x00, 0x00, 0x00, 0x00, 0x10, 0x02, 0xe0, 0xe0,
      0x00, 0x04, 0x02, 0x01, 0xd0, 0x11, 0x40, 0x08, 0x04, 0x02, 0x60,
      0x04, 0x00, 0x02, 0x1c, 0x00, 0x00, 0x43, 0x00, 0x06, 0x00, 0x00,
      0xf1, 0x10, 0x00, 0x01, 0x00, 0x64, 0x40, 0x08, 0x00, 0x00, 0xf1,
      0x10, 0x00, 0x00, 0x00, 0xa0, 0x00, 0x86, 0x40, 0x01, 0x30};

  bstring payload;
  payload = blk2bstr(&initial_ue_bytes, sizeof(initial_ue_bytes));
  S1ap_S1AP_PDU_t pdu;
  memset(&pdu, 0, sizeof(pdu));

  ASSERT_EQ(RETURNok, s1ap_mme_decode_pdu(&pdu, payload));
  ASSERT_EQ(RETURNok,
            s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu));

  handle_mme_ue_id_notification(state, assoc_id);

  // State validation
  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_READY, 1));
  ASSERT_TRUE(is_num_enbs_valid(state, 1));
  ASSERT_EQ(state->mmeid2associd.num_elements, 1);

  // Send S1AP_PAGING_REQUEST mimicing MME_APP
  ASSERT_EQ(send_s1ap_paging_request(assoc_id), RETURNok);

  // Freeing pdu and payload data
  bdestroy_wrapper(&payload);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_s1);
}

TEST_F(S1apMmeHandlersTest, HandleS1apErabModificationCnf) {
  ASSERT_EQ(task_zmq_ctx_main_s1ap.ready, true);

  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_INIT, 0));

  EXPECT_CALL(*sctp_handler, sctpd_send_dl()).Times(2);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_initial_ue_message()).Times(1);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_s1ap_ue_context_release_req())
      .Times(0);
  EXPECT_CALL(*mme_app_handler, nas_proc_dl_transfer_rej()).Times(0);

  S1ap_S1AP_PDU_t pdu_s1;
  memset(&pdu_s1, 0, sizeof(pdu_s1));
  ASSERT_EQ(RETURNok, generate_s1_setup_request_pdu(&pdu_s1));
  ASSERT_EQ(RETURNok,
            s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_s1));

  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_READY, 0));

  uint8_t initial_ue_bytes[] = {
      0x00, 0x0c, 0x40, 0x48, 0x00, 0x00, 0x05, 0x00, 0x08, 0x00, 0x02,
      0x00, 0x01, 0x00, 0x1a, 0x00, 0x20, 0x1f, 0x07, 0x41, 0x71, 0x08,
      0x09, 0x10, 0x10, 0x00, 0x00, 0x00, 0x00, 0x10, 0x02, 0xe0, 0xe0,
      0x00, 0x04, 0x02, 0x01, 0xd0, 0x11, 0x40, 0x08, 0x04, 0x02, 0x60,
      0x04, 0x00, 0x02, 0x1c, 0x00, 0x00, 0x43, 0x00, 0x06, 0x00, 0x00,
      0xf1, 0x10, 0x00, 0x01, 0x00, 0x64, 0x40, 0x08, 0x00, 0x00, 0xf1,
      0x10, 0x00, 0x00, 0x00, 0xa0, 0x00, 0x86, 0x40, 0x01, 0x30};

  ASSERT_EQ(simulate_pdu_s1_message(initial_ue_bytes, sizeof(initial_ue_bytes),
                                    state, assoc_id, stream_id),
            RETURNok);

  handle_mme_ue_id_notification(state, assoc_id);

  // State validation
  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_READY, 1));
  ASSERT_TRUE(is_num_enbs_valid(state, 1));
  ASSERT_EQ(state->mmeid2associd.num_elements, 1);

  // Send S1AP_E_RAB_MODIFICATION_CNF mimicing MME_APP
  ASSERT_EQ(send_s1ap_erab_mod_confirm(1, 7), RETURNok);

  // Freeing pdu and payload data
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_s1);
}

TEST_F(S1apMmeHandlersTest, HandleS1apNasNonDelivery) {
  ASSERT_EQ(task_zmq_ctx_main_s1ap.ready, true);

  EXPECT_CALL(*mme_app_handler, mme_app_handle_initial_ue_message()).Times(1);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_s1ap_ue_context_release_req())
      .Times(0);
  EXPECT_CALL(*mme_app_handler, nas_proc_dl_transfer_rej()).Times(1);

  bool is_state_same = false;

  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_INIT, 0));

  S1ap_S1AP_PDU_t pdu_s1;
  memset(&pdu_s1, 0, sizeof(pdu_s1));
  ASSERT_EQ(RETURNok, generate_s1_setup_request_pdu(&pdu_s1));
  ASSERT_EQ(RETURNok,
            s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_s1));

  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_READY, 0));

  uint8_t initial_ue_bytes[] = {
      0x00, 0x0c, 0x40, 0x48, 0x00, 0x00, 0x05, 0x00, 0x08, 0x00, 0x02,
      0x00, 0x01, 0x00, 0x1a, 0x00, 0x20, 0x1f, 0x07, 0x41, 0x71, 0x08,
      0x09, 0x10, 0x10, 0x00, 0x00, 0x00, 0x00, 0x10, 0x02, 0xe0, 0xe0,
      0x00, 0x04, 0x02, 0x01, 0xd0, 0x11, 0x40, 0x08, 0x04, 0x02, 0x60,
      0x04, 0x00, 0x02, 0x1c, 0x00, 0x00, 0x43, 0x00, 0x06, 0x00, 0x00,
      0xf1, 0x10, 0x00, 0x01, 0x00, 0x64, 0x40, 0x08, 0x00, 0x00, 0xf1,
      0x10, 0x00, 0x00, 0x00, 0xa0, 0x00, 0x86, 0x40, 0x01, 0x30};

  ASSERT_EQ(simulate_pdu_s1_message(initial_ue_bytes, sizeof(initial_ue_bytes),
                                    state, assoc_id, stream_id),
            RETURNok);

  handle_mme_ue_id_notification(state, assoc_id);

  // Generate downlink nas transport with dummy payload
  bstring p;
  std::string test_str = "test";
  STRING_TO_BSTRING(test_str, p);
  s1ap_generate_downlink_nas_transport(state, 1, 7, &p, 1, &is_state_same);
  bdestroy_wrapper(&p);

  // Authentication response proc packet bytes
  uint8_t auth_bytes[] = {
      0x00, 0x0d, 0x40, 0x3d, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x02,
      0x00, 0x07, 0x00, 0x08, 0x00, 0x02, 0x00, 0x01, 0x00, 0x1a, 0x00,
      0x14, 0x13, 0x07, 0x53, 0x10, 0x1e, 0x63, 0x7e, 0x5c, 0x58, 0xec,
      0x5a, 0xa8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
      0x64, 0x40, 0x08, 0x00, 0x00, 0xf1, 0x10, 0x00, 0x00, 0x00, 0xa0,
      0x00, 0x43, 0x40, 0x06, 0x00, 0x00, 0xf1, 0x10, 0x00, 0x01};

  ASSERT_EQ(simulate_pdu_s1_message(auth_bytes, sizeof(auth_bytes), state,
                                    assoc_id, stream_id),
            RETURNok);

  uint8_t ics_bytes[] = {0x20, 0x09, 0x00, 0x22, 0x00, 0x00, 0x03, 0x00,
                         0x00, 0x40, 0x02, 0x00, 0x07, 0x00, 0x08, 0x40,
                         0x02, 0x00, 0x01, 0x00, 0x33, 0x40, 0x0f, 0x00,
                         0x00, 0x32, 0x40, 0x0a, 0x0a, 0x1f, 0xc0, 0xa8,
                         0x3c, 0x8d, 0x0a, 0x00, 0x01, 0x28};

  ASSERT_EQ(simulate_pdu_s1_message(ics_bytes, sizeof(ics_bytes), state,
                                    assoc_id, stream_id),
            RETURNok);

  // Send NAS Non Delivery payload message
  uint8_t nas_non_delivery_bytes[] = {
      0x00, 0x10, 0x40, 0x28, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x02,
      0x00, 0x07, 0x00, 0x08, 0x00, 0x02, 0x00, 0x01, 0x00, 0x1a, 0x00,
      0x0f, 0x0e, 0x37, 0x2c, 0x71, 0xdc, 0xfa, 0x00, 0x07, 0x5d, 0x02,
      0x00, 0x02, 0xe0, 0xe0, 0xc1, 0x00, 0x02, 0x40, 0x02, 0x00, 0x60};

  ASSERT_EQ(simulate_pdu_s1_message(nas_non_delivery_bytes,
                                    sizeof(nas_non_delivery_bytes), state,
                                    assoc_id, 1),
            RETURNok);

  // State validation
  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_READY, 1));
  ASSERT_TRUE(is_num_enbs_valid(state, 1));
  ASSERT_EQ(state->mmeid2associd.num_elements, 1);

  // Freeing pdu and payload data
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_s1);
}

}  // namespace lte
}  // namespace magma
