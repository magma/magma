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

TEST_F(S1apMmeHandlersTest, HandleS1SetupRequestFailureHss) {
  EXPECT_CALL(*sctp_handler, sctpd_send_dl()).Times(1);

  hss_associated = false;

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
}

TEST_F(S1apMmeHandlersTest, HandleS1SetupRequestFailureReseting) {
  EXPECT_CALL(*sctp_handler, sctpd_send_dl()).Times(1);

  enb_description_t* enb_associated = NULL;
  hashtable_ts_get(&state->enbs, (const hash_key_t)assoc_id,
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
}

TEST_F(S1apMmeHandlersTest, HandleCloseSctpAssociation) {
  ASSERT_EQ(task_zmq_ctx_main_s1ap.ready, true);

  EXPECT_CALL(*sctp_handler, sctpd_send_dl()).Times(1);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_initial_ue_message()).Times(1);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_s1ap_ue_context_release_req())
      .Times(0);

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

  // Send SCTP_CLOSE_ASSOCIATION mimicing SCTP task
  ASSERT_EQ(send_s1ap_close_sctp_association(assoc_id), RETURNok);

  // Freeing pdu and payload data
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_s1);
}

TEST_F(S1apMmeHandlersTest, HandleEnbResetPartial) {
  ASSERT_EQ(task_zmq_ctx_main_s1ap.ready, true);

  bool is_state_same = true;

  EXPECT_CALL(*sctp_handler, sctpd_send_dl()).Times(2);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_initial_ue_message()).Times(1);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_s1ap_ue_context_release_req())
      .Times(0);
  EXPECT_CALL(*mme_app_handler, nas_proc_dl_transfer_rej()).Times(0);
  EXPECT_CALL(*mme_app_handler, mme_app_handle_enb_reset_req()).Times(1);

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

  // Simulate ENB Reset
  uint8_t enb_reset[] = {0x00, 0x0e, 0x00, 0x16, 0x00, 0x00, 0x02, 0x00, 0x02,
                         0x40, 0x01, 0x43, 0x00, 0x5c, 0x00, 0x0a, 0x40, 0x00,
                         0x00, 0x5b, 0x00, 0x04, 0x60, 0x07, 0x00, 0x01};

  ASSERT_EQ(simulate_pdu_s1_message(enb_reset, sizeof(enb_reset), state,
                                    assoc_id, stream_id),
            RETURNok);
}
}  // namespace lte
}  // namespace magma
