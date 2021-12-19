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
#include "lte/gateway/c/core/oai/test/ngap/util_ngap_pkt.h"
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"
#include <gtest/gtest.h>
#include <thread>

using ::testing::Test;

namespace magma5g {

TEST(test_ngap_pkt_tests, test_ngap_unsuccess_outcome_asn_raw) {
  bstring stream_setup_failure;
  Ngap_NGAP_PDU_t decode_pdu;
  Ngap_NGSetupFailure_t* container;
  Ngap_NGSetupFailureIEs_t* ie;
  int ret        = 0;
  bool decode_op = false;

  ret = ngap_ng_setup_failure_stream(
      Ngap_Cause_PR_misc, Ngap_CauseMisc_unspecified, stream_setup_failure);
  EXPECT_TRUE(ret == EXIT_SUCCESS);

  decode_op = ng_setup_failure_decode(stream_setup_failure, &decode_pdu);
  EXPECT_TRUE(decode_op == true);

  container =
      &decode_pdu.choice.unsuccessfulOutcome.value.choice.NGSetupFailure;
  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(
      Ngap_NGSetupFailureIEs_t, ie, container, Ngap_ProtocolIE_ID_id_Cause);

  EXPECT_FALSE(ie == nullptr);
  EXPECT_TRUE(ie->value.choice.Cause.present == Ngap_Cause_PR_misc);
  EXPECT_TRUE(ie->value.choice.Cause.choice.misc == Ngap_CauseMisc_unspecified);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decode_pdu);
  bdestroy(stream_setup_failure);
}

TEST(test_ngap_pkt_tests, test_ngap_unsuccess_outcome_pdu) {
  bstring stream_setup_failure;
  Ngap_NGAP_PDU_t encode_pdu;
  Ngap_NGAP_PDU_t decode_pdu;
  Ngap_NGSetupFailure_t* container;
  Ngap_NGSetupFailureIEs_t* ie;
  uint8_t* buffer_p = NULL;
  uint32_t length   = 0;
  int ret           = 0;
  bool decode_op    = false;

  memset(&encode_pdu, 0, sizeof(encode_pdu));

  ngap_ng_setup_failure_pdu(
      Ngap_Cause_PR_misc, Ngap_CauseMisc_unspecified, encode_pdu);
  ret = ngap_amf_encode_pdu(&encode_pdu, &buffer_p, &length);

  EXPECT_TRUE(ret == 0);

  stream_setup_failure = blk2bstr(buffer_p, length);
  free(buffer_p);

  memset(&decode_pdu, 0, sizeof(decode_pdu));
  decode_op = ng_setup_failure_decode(stream_setup_failure, &decode_pdu);
  EXPECT_TRUE(decode_op == true);

  container =
      &decode_pdu.choice.unsuccessfulOutcome.value.choice.NGSetupFailure;
  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(
      Ngap_NGSetupFailureIEs_t, ie, container, Ngap_ProtocolIE_ID_id_Cause);

  EXPECT_FALSE(ie == nullptr);
  EXPECT_TRUE(ie->value.choice.Cause.present == Ngap_Cause_PR_misc);
  EXPECT_TRUE(ie->value.choice.Cause.choice.misc == Ngap_CauseMisc_unspecified);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decode_pdu);
  bdestroy(stream_setup_failure);
}

TEST(test_ngap_pkt_tests, test_set_ng_setup_failure_cause) {
  Ngap_Cause_t cause_p;
  const long cause_value = 20;
  status_code_e ret = RETURNok;

  ret = ngap_amf_set_cause(&cause_p, Ngap_Cause_PR_radioNetwork, cause_value);
  EXPECT_TRUE(ret == 0);

  ret = ngap_amf_set_cause(&cause_p, Ngap_Cause_PR_transport, cause_value);
  EXPECT_TRUE(ret == 0);

  ret = ngap_amf_set_cause(&cause_p, Ngap_Cause_PR_nas, cause_value);
  EXPECT_TRUE(ret == 0);

  ret = ngap_amf_set_cause(&cause_p, Ngap_Cause_PR_protocol, cause_value);
  EXPECT_TRUE(ret == 0);

  ret = ngap_amf_set_cause(&cause_p, Ngap_Cause_PR_misc, cause_value);
  EXPECT_TRUE(ret == 0);

  ret = ngap_amf_set_cause(&cause_p, Ngap_Cause_PR_NOTHING, cause_value);
  EXPECT_TRUE(ret == RETURNerror);
}

TEST(test_ngap_pkt_tests, test_ngap_amf_encode_initiating) {
  Ngap_NGAP_PDU_t encode_pdu;
  uint8_t* buffer_p        = NULL;
  uint32_t length          = 0;
  int unknown_procedure_id = -1;
  int ret                  = -1;

  memset(&encode_pdu, 0, sizeof(encode_pdu));

  encode_pdu.present = Ngap_NGAP_PDU_PR_initiatingMessage;
  encode_pdu.choice.initiatingMessage.procedureCode =
	  Ngap_ProcedureCode_id_DownlinkNASTransport;

  ret = ngap_amf_encode_pdu(&encode_pdu, &buffer_p, &length);
  EXPECT_TRUE(ret == 0);

  free(buffer_p);

  memset(&encode_pdu, 0, sizeof(encode_pdu));

  encode_pdu.present = Ngap_NGAP_PDU_PR_initiatingMessage;
  encode_pdu.choice.initiatingMessage.procedureCode = unknown_procedure_id;

  ret = ngap_amf_encode_pdu(&encode_pdu, &buffer_p, &length);
  EXPECT_TRUE(ret == -1);

  free(buffer_p);
}

TEST(test_ngap_pkt_tests, test_ngap_amf_encode_successful_outcome) {
  Ngap_NGAP_PDU_t encode_pdu;
  uint8_t* buffer_p        = NULL;
  uint32_t length          = 0;
  int unknown_procedure_id = -1;
  int ret                  = -1;

  memset(&encode_pdu, 0, sizeof(encode_pdu));

  encode_pdu.present = Ngap_NGAP_PDU_PR_successfulOutcome;
  encode_pdu.choice.initiatingMessage.procedureCode = Ngap_ProcedureCode_id_NGSetup;

  ret = ngap_amf_encode_pdu(&encode_pdu, &buffer_p, &length);
  EXPECT_TRUE(ret == 0);

  free(buffer_p);

  memset(&encode_pdu, 0, sizeof(encode_pdu));

  encode_pdu.present = Ngap_NGAP_PDU_PR_successfulOutcome;
  encode_pdu.choice.initiatingMessage.procedureCode = unknown_procedure_id;

  ret = ngap_amf_encode_pdu(&encode_pdu, &buffer_p, &length);
  EXPECT_TRUE(ret == -1);

  free(buffer_p);
}

TEST(test_ngap_pkt_tests, test_ngap_amf_encode_unsuccessful_outcome) {
  Ngap_NGAP_PDU_t encode_pdu;
  uint8_t* buffer_p        = NULL;
  uint32_t length          = 0;
  int unknown_procedure_id = -1;
  int ret                  = -1;

  memset(&encode_pdu, 0, sizeof(encode_pdu));

  encode_pdu.present = Ngap_NGAP_PDU_PR_unsuccessfulOutcome;
  encode_pdu.choice.initiatingMessage.procedureCode = Ngap_ProcedureCode_id_NGSetup;

  ret = ngap_amf_encode_pdu(&encode_pdu, &buffer_p, &length);
  EXPECT_TRUE(ret == 0);

  free(buffer_p);

  memset(&encode_pdu, 0, sizeof(encode_pdu));

  encode_pdu.present = Ngap_NGAP_PDU_PR_unsuccessfulOutcome;
  encode_pdu.choice.initiatingMessage.procedureCode = unknown_procedure_id;

  ret = ngap_amf_encode_pdu(&encode_pdu, &buffer_p, &length);
  EXPECT_TRUE(ret == -1);

  free(buffer_p);
}

TEST(test_ngap_pkt_tests, test_ngap_amf_encode_unknown_msg_outcome) {
  Ngap_NGAP_PDU_t encode_pdu;
  uint8_t* buffer_p = NULL;
  uint32_t length   = 0;
  int ret           = 0;

  memset(&encode_pdu, 0, sizeof(encode_pdu));

  encode_pdu.present = Ngap_NGAP_PDU_PR_NOTHING;

  ret = ngap_amf_encode_pdu(&encode_pdu, &buffer_p, &length);
  EXPECT_TRUE(ret == -1);

  free(buffer_p);
}

TEST(test_ngap_pkt_tests, test_ng_gnb_state2str) {
  const char* state_to_string = NULL;

  state_to_string = ng_gnb_state2str(NGAP_INIT);
  EXPECT_STREQ("NGAP_INIT", state_to_string);

  state_to_string = ng_gnb_state2str(NGAP_RESETING);
  EXPECT_STREQ("NGAP_RESETING", state_to_string);

  state_to_string = ng_gnb_state2str(NGAP_READY);
  EXPECT_STREQ("NGAP_READY", state_to_string);

  state_to_string = ng_gnb_state2str(NGAP_SHUTDOWN);
  EXPECT_STREQ("NGAP_SHUTDOWN", state_to_string);
}

TEST(test_ngap_pkt_tests, test_ngap_direction2str) {
  const char* dir_to_string = NULL;
  uint8_t invalid_case      = 20;

  dir_to_string = ngap_direction2str(Ngap_NGAP_PDU_PR_NOTHING);
  EXPECT_STREQ("<nothing>", dir_to_string);

  dir_to_string = ngap_direction2str(Ngap_NGAP_PDU_PR_initiatingMessage);
  EXPECT_STREQ("originating message", dir_to_string);

  dir_to_string = ngap_direction2str(Ngap_NGAP_PDU_PR_successfulOutcome);
  EXPECT_STREQ("successful outcome", dir_to_string);

  dir_to_string = ngap_direction2str(Ngap_NGAP_PDU_PR_unsuccessfulOutcome);
  EXPECT_STREQ("unsuccessful outcome", dir_to_string);

  dir_to_string = ngap_direction2str(invalid_case);
  EXPECT_STRNE("invalid case", dir_to_string);
}

TEST(test_ngap_pkt_tests, test_ngap_amf_handle_error_ind_message) {
  Ngap_NGAP_PDU_t message;
  ngap_state_t* state            = NULL;
  const sctp_assoc_id_t assoc_id = 1;
  const sctp_stream_id_t stream  = 1;
  int ret = -1;

  memset(&message, 0, sizeof(message));

  ret = ngap_amf_handle_error_ind_message(state, assoc_id, stream, &message);
  EXPECT_TRUE(ret == RETURNok);
}

TEST(test_ngap_pkt_tests, test_ngap_amf_compare_tac) {
  Ngap_TAC_t tac;
  int ret            = -1;
  uint8_t tac_buf[3] = {0x01, 0x02, 0x03};

  memset(&tac, 0, sizeof(Ngap_TAC_t));

  tac.buf = tac_buf;
  tac.size = sizeof(tac_buf);

  ret = ngap_amf_compare_tac(&tac);
  EXPECT_TRUE(ret == TA_LIST_AT_LEAST_ONE_MATCH);

  // TODO: Add UT case for TA_LIST_NO_MATCH once mme config
  //       parsing issue is fixed
}

TEST(test_ngap_pkt_tests, test_ngap_initiate_ue_message) {
  bool output    = false;
  int decode_ops = -1;
  bstring stream;
  Ngap_NGAP_PDU_t decode_pdu;

  output = ngap_initiate_ue_message(stream);

  // Check if encoding is successful
  EXPECT_TRUE(output == true);

  memset(&decode_pdu, 0, sizeof(decode_pdu));

  decode_ops = ngap_amf_decode_pdu(&decode_pdu, stream);
  EXPECT_TRUE(decode_ops == 0);

  bdestroy(stream);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decode_pdu);
}

TEST(test_ngap_pkt_tests, test_ngap_pdusession_resource_setup_req) {
  bool output    = false;
  int decode_ops = -1;
  bstring stream;
  Ngap_NGAP_PDU_t decode_pdu;

  output = generator_ngap_pdusession_resource_setup_req(stream);

  // Check if encoding is successful
  EXPECT_TRUE(output == true);
  EXPECT_TRUE(blength(stream) != 0);

  memset(&decode_pdu, 0, sizeof(decode_pdu));

  decode_ops = ngap_amf_decode_pdu(&decode_pdu, stream);
  EXPECT_TRUE(decode_ops == 0);
  bdestroy(stream);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decode_pdu);
}

TEST(test_ngap_pkt_tests, test_ngap_pdusession_resource_setup_stream) {
  bool output    = false;
  int decode_ops = -1;
  bstring stream;
  Ngap_NGAP_PDU_t decode_pdu;

  output = generator_itti_ngap_pdusession_resource_setup_req(stream);

  // Check if encoding is successful
  EXPECT_TRUE(output == true);
  EXPECT_TRUE(blength(stream) != 0);

  memset(&decode_pdu, 0, sizeof(decode_pdu));

  decode_ops = ngap_amf_decode_pdu(&decode_pdu, stream);
  EXPECT_TRUE(decode_ops == 0);
  bdestroy(stream);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decode_pdu);
}

TEST(test_ngap_pkt_tests, test_ngap_pdusession_resource_rel_cmd_stream) {
  bool output    = false;
  int decode_ops = -1;
  bstring stream;
  Ngap_NGAP_PDU_t decode_pdu;

  output = generator_ngap_pdusession_resource_rel_cmd_stream(stream);

  // Check if encoding is successful
  EXPECT_TRUE(output == true);
  EXPECT_TRUE(blength(stream) != 0);

  memset(&decode_pdu, 0, sizeof(decode_pdu));

  decode_ops = ngap_amf_decode_pdu(&decode_pdu, stream);
  EXPECT_TRUE(decode_ops == 0);
  bdestroy(stream);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decode_pdu);
}

TEST(test_ngap_pkt_tests, test_ngap_init_ue_msg_pdu) {
  Ngap_NGAP_PDU_t init_ue_pdu;
  bool res = false;
  gnb_description_t gNB_ref;
  m5g_ue_description_t ue_ref;

  memset(&init_ue_pdu, 0, sizeof(Ngap_NGAP_PDU_t));

  memset(&ue_ref, 0, sizeof(m5g_ue_description_t));
  memset(&gNB_ref, 0, sizeof(gnb_description_t));

  gNB_ref.sctp_assoc_id    = 1;
  gNB_ref.next_sctp_stream = 3;
  gNB_ref.gnb_id           = 2;

  memset(&ue_ref, 0, sizeof(m5g_ue_description_t));

  res = generate_guti_ngap_pdu(&init_ue_pdu);
  EXPECT_TRUE(res == true);

  res = validate_handle_initial_ue_message(&gNB_ref, &ue_ref, &init_ue_pdu);
  EXPECT_TRUE(res == true);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &init_ue_pdu);
}

TEST(test_ngap_pkt_tests, test_ngap_Five_G_TMSI) {
  Ngap_NGAP_PDU_t decode_pdu;
  int decode_ops = -1;
  bstring stream_fiveg_tmsi;
  gnb_description_t gNB_ref;
  m5g_ue_description_t ue_ref;
  bool res = false;

  uint8_t buffer[] = {
      0x00, 0x0f, 0x40, 0x67, 0x00, 0x00, 0x05, 0x00, 0x55, 0x00, 0x02, 0x00,
      0x0c, 0x00, 0x26, 0x00, 0x37, 0x36, 0x7e, 0x01, 0xdc, 0xcf, 0x83, 0xd1,
      0x02, 0x7e, 0x00, 0x41, 0x03, 0x00, 0x0b, 0xf2, 0x22, 0x62, 0x54, 0x00,
      0x00, 0x00, 0x3a, 0xc1, 0x8e, 0x00, 0x71, 0x00, 0x1b, 0x7e, 0x00, 0x41,
      0x03, 0x00, 0x0b, 0xf2, 0x22, 0x62, 0x54, 0x00, 0x00, 0x00, 0x3a, 0xc1,
      0x8e, 0x00, 0x50, 0x02, 0x00, 0x00, 0x18, 0x01, 0x00, 0x53, 0x01, 0x01,
      0x00, 0x79, 0x00, 0x0f, 0x40, 0x22, 0x42, 0x65, 0x00, 0x00, 0x00, 0x01,
      0x00, 0x22, 0x42, 0x65, 0x00, 0x00, 0x01, 0x00, 0x5a, 0x40, 0x01, 0x18,
      0x00, 0x1a, 0x00, 0x07, 0x00, 0x00, 0x00, 0x3a, 0xc1, 0x8e, 0x00};

  stream_fiveg_tmsi = blk2bstr(buffer, sizeof(buffer) / sizeof(uint8_t));

  memset(&decode_pdu, 0, sizeof(decode_pdu));
  memset(&ue_ref, 0, sizeof(m5g_ue_description_t));
  memset(&gNB_ref, 0, sizeof(gnb_description_t));

  decode_ops = ngap_amf_decode_pdu(&decode_pdu, stream_fiveg_tmsi);
  EXPECT_TRUE(decode_ops == 0);

  gNB_ref.sctp_assoc_id    = 1;
  gNB_ref.next_sctp_stream = 3;
  gNB_ref.gnb_id           = 2;

  res = validate_handle_initial_ue_message(&gNB_ref, &ue_ref, &decode_pdu);
  EXPECT_TRUE(res == true);

  bdestroy(stream_fiveg_tmsi);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decode_pdu);
}

TEST(test_ngap_pkt_tests, test_ngap_setup_request_sd) {
  Ngap_NGAP_PDU_t init_ue_pdu;
  bool res = false;
  gnb_description_t gNB_ref;
  m5g_ue_description_t ue_ref;

  memset(&init_ue_pdu, 0, sizeof(Ngap_NGAP_PDU_t));

  memset(&ue_ref, 0, sizeof(m5g_ue_description_t));
  memset(&gNB_ref, 0, sizeof(gnb_description_t));

  res = generate_ngap_request_msg(&init_ue_pdu);
  EXPECT_TRUE(res == true);

  res = validate_ngap_setup_request(&init_ue_pdu);
  EXPECT_TRUE(res == true);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &init_ue_pdu);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace magma5g
