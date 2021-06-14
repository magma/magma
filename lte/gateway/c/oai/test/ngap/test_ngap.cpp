/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#include "util_ngap_pkt.h"
#include "dynamic_memory_check.h"
#include <gtest/gtest.h>
#include <thread>

using ::testing::Test;

namespace magma {
namespace lte {

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

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace lte
}  // namespace magma
