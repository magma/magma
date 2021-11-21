/**
 * Copyright 2021 The Magma Authors.
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

#include "../mock_tasks/mock_tasks.h"

extern "C" {
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/include/mme_config.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme_decoder.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme_handlers.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme_nas_procedures.h"
}

#include "lte/gateway/c/core/oai/test/s1ap_task/s1ap_mme_test_utils.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_state_manager.h"

extern bool hss_associated;

namespace magma {
namespace lte {

task_zmq_ctx_t task_zmq_ctx_main_s1ap;

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    default: { } break; }

  itti_free_msg_content(received_message_p);
  free(received_message_p);
  return 0;
}

class S1apMmeHandlersTest : public ::testing::Test {
  virtual void SetUp() {
    mme_app_handler = std::make_shared<MockMmeAppHandler>();
    sctp_handler    = std::make_shared<MockSctpHandler>();

    itti_init(
        TASK_MAX, THREAD_MAX, MESSAGES_ID_MAX, tasks_info, messages_info, NULL,
        NULL);

    // initialize mme config
    mme_config_init(&mme_config);
    create_partial_lists(&mme_config);
    mme_config.use_stateless = false;
    hss_associated           = true;

    task_id_t task_id_list[4] = {TASK_MME_APP, TASK_S1AP, TASK_SCTP,
                                 TASK_SERVICE303};
    init_task_context(
        TASK_MAIN, task_id_list, 4, handle_message, &task_zmq_ctx_main_s1ap);

    std::thread task_mme_app(start_mock_mme_app_task, mme_app_handler);
    std::thread task_sctp(start_mock_sctp_task, sctp_handler);
    task_mme_app.detach();
    task_sctp.detach();

    s1ap_mme_init(&mme_config);

    // Setup new association for testing
    state     = S1apStateManager::getInstance().get_state(false);
    assoc_id  = 1;
    stream_id = 0;
    setup_new_association(state, assoc_id);
    std::this_thread::sleep_for(std::chrono::milliseconds(500));
  }

  virtual void TearDown() {
    send_terminate_message_fatal(&task_zmq_ctx_main_s1ap);
    destroy_task_context(&task_zmq_ctx_main_s1ap);
    itti_free_desc_threads();

    free_mme_config(&mme_config);

    // Sleep to ensure that messages are received and contexts are released
    std::this_thread::sleep_for(std::chrono::milliseconds(1000));
  }

 protected:
  std::shared_ptr<MockMmeAppHandler> mme_app_handler;
  std::shared_ptr<MockSctpHandler> sctp_handler;
  s1ap_state_t* state;
  sctp_assoc_id_t assoc_id;
  sctp_stream_id_t stream_id;
};

TEST_F(S1apMmeHandlersTest, HandleS1SetupRequestFailureHss) {
  // Setup new association for testing
  s1ap_state_t* s          = S1apStateManager::getInstance().get_state(false);
  sctp_assoc_id_t assoc_id = 1;
  setup_new_association(s, assoc_id);

  EXPECT_CALL(*sctp_handler, sctpd_send_dl()).Times(1);

  hss_associated = false;

  S1ap_S1AP_PDU_t pdu_s1;
  memset(&pdu_s1, 0, sizeof(pdu_s1));
  status_code_e pdu_rc = generate_s1_setup_request_pdu(&pdu_s1);
  ASSERT_EQ(pdu_rc, RETURNok);

  sctp_stream_id_t stream_id = 0;
  status_code_e rc = s1ap_mme_handle_message(s, assoc_id, stream_id, &pdu_s1);
  ASSERT_EQ(rc, RETURNok);

  // State validation
  ASSERT_EQ(s->num_enbs, 0);

  // Freeing pdu and payload data
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_s1);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(1000));
}

TEST_F(S1apMmeHandlersTest, HandleS1SetupRequestFailureReseting) {
  EXPECT_CALL(*sctp_handler, sctpd_send_dl()).Times(1);

  enb_description_t* enb_associated = NULL;
  hashtable_ts_get(
      &state->enbs, (const hash_key_t) assoc_id,
      reinterpret_cast<void**>(&enb_associated));
  enb_associated->s1_state = S1AP_RESETING;

  S1ap_S1AP_PDU_t pdu_s1;
  memset(&pdu_s1, 0, sizeof(pdu_s1));
  status_code_e pdu_rc = generate_s1_setup_request_pdu(&pdu_s1);
  ASSERT_EQ(pdu_rc, RETURNok);

  sctp_stream_id_t stream_id = 0;
  status_code_e rc =
      s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_s1);
  ASSERT_EQ(rc, RETURNok);

  // State validation
  ASSERT_EQ(state->num_enbs, 0);

  // Freeing pdu and payload data
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_s1);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(1500));
}

TEST_F(S1apMmeHandlersTest, HandleICSResponseICSRelease) {
  ASSERT_EQ(task_zmq_ctx_main_s1ap.ready, true);

  bool is_state_same = false;

  EXPECT_CALL(*sctp_handler, sctpd_send_dl()).Times(2);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_initial_ue_message()).Times(1);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_s1ap_ue_context_release_req())
      .Times(1);

  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_INIT, 0));

  S1ap_S1AP_PDU_t pdu_s1;
  memset(&pdu_s1, 0, sizeof(pdu_s1));
  ASSERT_EQ(RETURNok, generate_s1_setup_request_pdu(&pdu_s1));
  ASSERT_EQ(
      RETURNok, s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_s1));

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

  bstring payload;
  payload = blk2bstr(&initial_ue_bytes, sizeof(initial_ue_bytes));
  S1ap_S1AP_PDU_t pdu;
  memset(&pdu, 0, sizeof(pdu));

  ASSERT_EQ(RETURNok, s1ap_mme_decode_pdu(&pdu, payload));
  ASSERT_EQ(
      RETURNok, s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu));

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

  bstring payload_nas;
  payload_nas = blk2bstr(&auth_bytes, sizeof(auth_bytes));
  S1ap_S1AP_PDU_t pdu_nas;
  memset(&pdu_nas, 0, sizeof(pdu_nas));

  ASSERT_EQ(s1ap_mme_decode_pdu(&pdu_nas, payload_nas), RETURNok);
  ASSERT_EQ(
      s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_nas), RETURNok);

  uint8_t ics_bytes[] = {0x20, 0x09, 0x00, 0x22, 0x00, 0x00, 0x03, 0x00,
                         0x00, 0x40, 0x02, 0x00, 0x07, 0x00, 0x08, 0x40,
                         0x02, 0x00, 0x01, 0x00, 0x33, 0x40, 0x0f, 0x00,
                         0x00, 0x32, 0x40, 0x0a, 0x0a, 0x1f, 0xc0, 0xa8,
                         0x3c, 0x8d, 0x0a, 0x00, 0x01, 0x28};

  bstring payload_ics;
  payload_ics = blk2bstr(&ics_bytes, sizeof(ics_bytes));
  S1ap_S1AP_PDU_t pdu_ics;
  memset(&pdu_ics, 0, sizeof(pdu_ics));

  ASSERT_EQ(s1ap_mme_decode_pdu(&pdu_ics, payload_ics), RETURNok);
  ASSERT_EQ(
      s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_ics), RETURNok);

  // Freeing pdu and payload data
  bdestroy_wrapper(&payload);
  bdestroy_wrapper(&payload_nas);
  bdestroy_wrapper(&payload_ics);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_s1);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_nas);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_ics);

  // State validation
  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_READY, 1));
  ASSERT_TRUE(is_num_enbs_valid(state, 1));

  uint8_t ics_release_bytes[] = {0x00, 0x12, 0x40, 0x15, 0x00, 0x00, 0x03,
                                 0x00, 0x00, 0x00, 0x02, 0x00, 0x07, 0x00,
                                 0x08, 0x00, 0x02, 0x00, 0x01, 0x00, 0x02,
                                 0x40, 0x02, 0x02, 0x80};

  bstring payload_ics_r;
  payload_ics_r = blk2bstr(&ics_release_bytes, sizeof(ics_release_bytes));
  S1ap_S1AP_PDU_t pdu_ics_r;
  memset(&pdu_ics_r, 0, sizeof(pdu_ics_r));

  ASSERT_EQ(s1ap_mme_decode_pdu(&pdu_ics_r, payload_ics_r), RETURNok);
  ASSERT_EQ(
      s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_ics_r),
      RETURNok);

  // State validation
  ASSERT_TRUE(is_num_enbs_valid(state, 1));

  bdestroy_wrapper(&payload_ics_r);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_ics_r);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(1000));
}

TEST_F(S1apMmeHandlersTest, HandleUECapIndication) {
  ASSERT_EQ(task_zmq_ctx_main_s1ap.ready, true);

  EXPECT_CALL(*sctp_handler, sctpd_send_dl()).Times(1);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_initial_ue_message()).Times(1);

  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_INIT, 0));

  S1ap_S1AP_PDU_t pdu_s1;
  memset(&pdu_s1, 0, sizeof(pdu_s1));
  ASSERT_EQ(RETURNok, generate_s1_setup_request_pdu(&pdu_s1));
  ASSERT_EQ(
      RETURNok, s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_s1));

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

  bstring payload;
  payload = blk2bstr(&initial_ue_bytes, sizeof(initial_ue_bytes));
  S1ap_S1AP_PDU_t pdu;
  memset(&pdu, 0, sizeof(pdu));

  ASSERT_EQ(RETURNok, s1ap_mme_decode_pdu(&pdu, payload));
  ASSERT_EQ(
      RETURNok, s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu));

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
  ASSERT_EQ(
      s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_cap), RETURNok);

  // State validation
  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_READY, 1));
  ASSERT_TRUE(is_num_enbs_valid(state, 1));

  // Freeing pdu and payload data
  bdestroy_wrapper(&payload);
  bdestroy_wrapper(&payload_ue_cap);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_cap);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_s1);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(500));
}

TEST_F(S1apMmeHandlersTest, GenerateUEContextReleaseCommand) {
  ue_description_t ue_ref_p = {
      .enb_ue_s1ap_id = 1,
      .mme_ue_s1ap_id = 1,
      .sctp_assoc_id  = assoc_id,
      .comp_s1ap_id   = S1AP_GENERATE_COMP_S1AP_ID(assoc_id, 1)};

  EXPECT_CALL(*sctp_handler, sctpd_send_dl()).Times(2);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_initial_ue_message()).Times(1);

  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_INIT, 0));

  S1ap_S1AP_PDU_t pdu_s1;
  memset(&pdu_s1, 0, sizeof(pdu_s1));
  ASSERT_EQ(RETURNok, generate_s1_setup_request_pdu(&pdu_s1));
  ASSERT_EQ(
      RETURNok, s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_s1));

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

  bstring payload;
  payload = blk2bstr(&initial_ue_bytes, sizeof(initial_ue_bytes));
  S1ap_S1AP_PDU_t pdu;
  memset(&pdu, 0, sizeof(pdu));

  ASSERT_EQ(RETURNok, s1ap_mme_decode_pdu(&pdu, payload));
  ASSERT_EQ(
      RETURNok, s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu));

  // Invalid S1 Cause returns error
  ASSERT_EQ(
      RETURNerror, s1ap_mme_generate_ue_context_release_command(
                       state, &ue_ref_p, S1AP_IMPLICIT_CONTEXT_RELEASE,
                       INVALID_IMSI64, assoc_id, stream_id, 1, 1));
  // Valid S1 Causes passess successfully
  ASSERT_EQ(
      RETURNok, s1ap_mme_generate_ue_context_release_command(
                    state, &ue_ref_p, S1AP_INITIAL_CONTEXT_SETUP_FAILED,
                    INVALID_IMSI64, assoc_id, stream_id, 1, 1));

  // State validation
  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_READY, 0));

  // Freeing pdu and payload data
  bdestroy_wrapper(&payload);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_s1);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(500));
}

TEST_F(S1apMmeHandlersTest, HandleUEContextReleaseCommand) {
  ASSERT_EQ(task_zmq_ctx_main_s1ap.ready, true);

  EXPECT_CALL(*sctp_handler, sctpd_send_dl()).Times(1);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_initial_ue_message()).Times(1);

  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_INIT, 0));

  S1ap_S1AP_PDU_t pdu_s1;
  memset(&pdu_s1, 0, sizeof(pdu_s1));
  ASSERT_EQ(RETURNok, generate_s1_setup_request_pdu(&pdu_s1));
  ASSERT_EQ(
      RETURNok, s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_s1));

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
  ASSERT_EQ(
      RETURNok, s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu));

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
  ASSERT_EQ(
      send_msg_to_task(&task_zmq_ctx_main_s1ap, TASK_S1AP, message_p),
      RETURNok);

  // Freeing pdu and payload data
  bdestroy_wrapper(&payload);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_s1);

  ASSERT_TRUE(is_num_enbs_valid(state, 1));

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(500));

  ASSERT_EQ(state->mmeid2associd.num_elements, 0);
}

TEST_F(S1apMmeHandlersTest, HandleConnectionEstCnf) {
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
  ASSERT_EQ(
      RETURNok, s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_s1));

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
  ASSERT_EQ(
      RETURNok, s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu));

  handle_mme_ue_id_notification(state, assoc_id);

  ASSERT_EQ(state->mmeid2associd.num_elements, 1);

  // Send UE connection establishment cnf mimicing MME_APP

  ASSERT_EQ(send_conn_establishment_cnf(7, true, true), RETURNok);

  // Freeing pdu and payload data
  bdestroy_wrapper(&payload);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_s1);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(500));
}

TEST_F(S1apMmeHandlersTest, HandleS1apErabRelCmd) {
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
  ASSERT_EQ(
      RETURNok, s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_s1));

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
  ASSERT_EQ(
      RETURNok, s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu));

  handle_mme_ue_id_notification(state, assoc_id);

  ASSERT_EQ(state->mmeid2associd.num_elements, 1);

  // Send S1AP_ERAB_REL_CMD mimicing MME_APP
  ASSERT_EQ(send_s1ap_erab_rel_cmd(7, 1), RETURNok);

  // Freeing pdu and payload data
  bdestroy_wrapper(&payload);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_s1);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(500));
}

TEST_F(S1apMmeHandlersTest, HandleS1apErabSetupReq) {
  ASSERT_EQ(task_zmq_ctx_main_s1ap.ready, true);

  EXPECT_CALL(*sctp_handler, sctpd_send_dl()).Times(2);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_initial_ue_message()).Times(1);

  ASSERT_TRUE(is_enb_state_valid(state, assoc_id, S1AP_INIT, 0));

  S1ap_S1AP_PDU_t pdu_s1;
  memset(&pdu_s1, 0, sizeof(pdu_s1));
  ASSERT_EQ(RETURNok, generate_s1_setup_request_pdu(&pdu_s1));
  ASSERT_EQ(
      RETURNok, s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_s1));

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
  ASSERT_EQ(
      RETURNok, s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu));

  handle_mme_ue_id_notification(state, assoc_id);

  ASSERT_EQ(state->mmeid2associd.num_elements, 1);

  // Send S1AP_ERAB_REL_CMD mimicing MME_APP
  ASSERT_EQ(send_s1ap_erab_setup_req(7, 1, 1), RETURNok);

  // Freeing pdu and payload data
  bdestroy_wrapper(&payload);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_s1);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(500));
}

TEST_F(S1apMmeHandlersTest, HandleS1apErabResetReq) {
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
  ASSERT_EQ(
      RETURNok, s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_s1));

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
  ASSERT_EQ(
      RETURNok, s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu));

  handle_mme_ue_id_notification(state, assoc_id);

  ASSERT_EQ(state->mmeid2associd.num_elements, 1);

  // Send S1AP_ERAB_REL_CMD mimicing MME_APP
  ASSERT_EQ(send_s1ap_erab_reset_req(assoc_id, stream_id, 1, 7), RETURNok);

  // Freeing pdu and payload data
  bdestroy_wrapper(&payload);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_s1);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(500));
}

TEST_F(S1apMmeHandlersTest, HandleS1apNasNonDeliveryFailure) {
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
  ASSERT_EQ(
      RETURNok, s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_s1));

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
  ASSERT_EQ(
      RETURNok, s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu));

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

  bstring payload_nas;
  payload_nas = blk2bstr(&auth_bytes, sizeof(auth_bytes));
  S1ap_S1AP_PDU_t pdu_nas;
  memset(&pdu_nas, 0, sizeof(pdu_nas));

  ASSERT_EQ(s1ap_mme_decode_pdu(&pdu_nas, payload_nas), RETURNok);
  ASSERT_EQ(
      s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_nas), RETURNok);

  uint8_t ics_bytes[] = {0x20, 0x09, 0x00, 0x22, 0x00, 0x00, 0x03, 0x00,
                         0x00, 0x40, 0x02, 0x00, 0x07, 0x00, 0x08, 0x40,
                         0x02, 0x00, 0x01, 0x00, 0x33, 0x40, 0x0f, 0x00,
                         0x00, 0x32, 0x40, 0x0a, 0x0a, 0x1f, 0xc0, 0xa8,
                         0x3c, 0x8d, 0x0a, 0x00, 0x01, 0x28};

  bstring payload_ics;
  payload_ics = blk2bstr(&ics_bytes, sizeof(ics_bytes));
  S1ap_S1AP_PDU_t pdu_ics;
  memset(&pdu_ics, 0, sizeof(pdu_ics));

  ASSERT_EQ(s1ap_mme_decode_pdu(&pdu_ics, payload_ics), RETURNok);
  ASSERT_EQ(
      s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu_ics), RETURNok);

  // Send NAS Non Delivery payload message
  uint8_t nas_non_delivery_bytes[] = {
      0x00, 0x10, 0x40, 0x28, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x02,
      0x00, 0x07, 0x00, 0x08, 0x00, 0x02, 0x00, 0x01, 0x00, 0x1a, 0x00,
      0x0f, 0x0e, 0x37, 0x2c, 0x71, 0xdc, 0xfa, 0x00, 0x07, 0x5d, 0x02,
      0x00, 0x02, 0xe0, 0xe0, 0xc1, 0x00, 0x02, 0x40, 0x02, 0x00, 0x60};

  bstring payload_nas_del;
  payload_nas_del =
      blk2bstr(&nas_non_delivery_bytes, sizeof(nas_non_delivery_bytes));
  S1ap_S1AP_PDU_t pdu_nas_del;
  memset(&pdu_nas_del, 0, sizeof(pdu_nas_del));

  ASSERT_EQ(s1ap_mme_decode_pdu(&pdu_nas_del, payload_nas_del), RETURNok);
  ASSERT_EQ(
      s1ap_mme_handle_message(state, assoc_id, 1, &pdu_nas_del), RETURNok);

  // Freeing pdu and payload data
  bdestroy_wrapper(&payload);
  bdestroy_wrapper(&payload_nas);
  bdestroy_wrapper(&payload_ics);
  bdestroy_wrapper(&payload_nas_del);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_s1);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_nas);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_ics);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_nas_del);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(500));
}

}  // namespace lte
}  // namespace magma
