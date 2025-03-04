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

#include "lte/gateway/c/core/oai/test/ngap/mock_utils.h"
#include "lte/gateway/c/core/oai/test/ngap/util_ngap_pkt.hpp"
#include <gtest/gtest.h>

extern "C" {
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/tasks/ngap/ngap_amf_handlers.h"
#include "lte/gateway/c/core/oai/include/amf_config.hpp"
}
#include "lte/gateway/c/core/oai/tasks/ngap/ngap_state_manager.hpp"

using ::testing::Test;

extern bool hss_associated;

namespace magma5g {

class NgapFlowTest : public testing::Test {
 protected:
  void SetUp() {
    itti_init(TASK_MAX, THREAD_MAX, MESSAGES_ID_MAX, tasks_info, messages_info,
              NULL, NULL);

    amf_config_init(&amf_config);
    amf_config.plmn_support_list.plmn_support_count = 1;
    amf_config.plmn_support_list.plmn_support[0].plmn = {.mcc_digit2 = 2,
                                                         .mcc_digit1 = 2,
                                                         .mnc_digit3 = 6,
                                                         .mcc_digit3 = 2,
                                                         .mnc_digit2 = 5,
                                                         .mnc_digit1 = 4};

    amf_config.plmn_support_list.plmn_support[0].s_nssai.sst = 0x1;
    amf_config.served_tai.plmn_mcc[0] = 222;
    amf_config.served_tai.plmn_mnc[0] = 456;
    amf_config.served_tai.plmn_mnc_len[0] = 3;

    ngap_state_init(2, 2, false);
    state = get_ngap_state(false);
    hss_associated = true;

    ran_cp_ipaddr = bfromcstr("\xc0\xa8\x3c\x8d");
    peerInfo = {
        .instreams = 1,
        .outstreams = 2,
        .assoc_id = 3,
        .ran_cp_ipaddr = ran_cp_ipaddr,
    };
    NGAPClientServicer::getInstance().msgtype_stack.clear();
  }

  void TearDown() {
    ngap_state_exit();
    bdestroy(ran_cp_ipaddr);
    itti_free_desc_threads();
    amf_config_free(&amf_config);
    NGAPClientServicer::getInstance().msgtype_stack.clear();
  }

  ngap_state_t* state = NULL;
  bstring ran_cp_ipaddr;
  sctp_new_peer_t peerInfo;
  const unsigned int AMF_UE_NGAP_ID = 0x05;
  const unsigned int gNB_UE_NGAP_ID = 0x09;
};

// Unit for Ng setup request message
TEST_F(NgapFlowTest, test_ngap_setup_request) {
  unsigned char ngap_setup_req_hexbuf[] = {
      0x00, 0x15, 0x00, 0x42, 0x00, 0x00, 0x04, 0x00, 0x1b, 0x00, 0x09, 0x00,
      0x22, 0x42, 0x65, 0x50, 0x00, 0x00, 0x00, 0x01, 0x00, 0x52, 0x40, 0x18,
      0x0a, 0x80, 0x55, 0x45, 0x52, 0x41, 0x4e, 0x53, 0x49, 0x4d, 0x2d, 0x67,
      0x6e, 0x62, 0x2d, 0x32, 0x32, 0x32, 0x2d, 0x34, 0x35, 0x36, 0x2d, 0x31,
      0x00, 0x66, 0x00, 0x0d, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x22, 0x42,
      0x65, 0x00, 0x00, 0x00, 0x08, 0x00, 0x15, 0x40, 0x01, 0x40};

  // Verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  // Verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);

  // Decode the PDU
  Ngap_NGAP_PDU_t decoded_pdu = {};
  uint16_t length = sizeof(ngap_setup_req_hexbuf) / sizeof(unsigned char);

  bstring ngap_setup_req_msg = blk2bstr(ngap_setup_req_hexbuf, length);

  // Check if the pdu can be decoded
  ASSERT_EQ(ngap_amf_decode_pdu(&decoded_pdu, ngap_setup_req_msg), RETURNok);

  sctp_stream_id_t stream_id = 0;

  gnb_description_t* gnb_association = NULL;
  gnb_association = ngap_state_get_gnb(state, peerInfo.assoc_id);
  gnb_association->ng_state = NGAP_SHUTDOWN;

  // check if Ng setup request message is handled successfully
  EXPECT_EQ(ngap_amf_handle_message(state, peerInfo.assoc_id, stream_id,
                                    &decoded_pdu),
            RETURNok);

  // As gnb in shutdown, the gnb should be with default values
  EXPECT_EQ(gnb_association->gnb_id, 0);

  gnb_association->ng_state = NGAP_READY;
  // check if Ng setup request message is handled successfully
  EXPECT_EQ(ngap_amf_handle_message(state, peerInfo.assoc_id, stream_id,
                                    &decoded_pdu),
            RETURNok);

  // As gnb in ready, the gnb id should be 1
  EXPECT_EQ(gnb_association->gnb_id, 1);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decoded_pdu);
  bdestroy(ngap_setup_req_msg);
}

// Unit for Initial UE message
TEST_F(NgapFlowTest, initial_ue_message_sunny_day) {
  unsigned char initial_ue_message_hexbuf[] = {
      0x00, 0x0f, 0x40, 0x48, 0x00, 0x00, 0x05, 0x00, 0x55, 0x00, 0x02,
      0x00, 0x01, 0x00, 0x26, 0x00, 0x1a, 0x19, 0x7e, 0x00, 0x41, 0x79,
      0x00, 0x0d, 0x01, 0x22, 0x62, 0x54, 0x00, 0x00, 0x00, 0x00, 0x00,
      0x00, 0x00, 0x00, 0x01, 0x2e, 0x04, 0xf0, 0xf0, 0xf0, 0xf0, 0x00,
      0x79, 0x00, 0x13, 0x50, 0x22, 0x42, 0x65, 0x00, 0x00, 0x00, 0x01,
      0x00, 0x22, 0x42, 0x65, 0x00, 0x00, 0x01, 0xe4, 0xf7, 0x04, 0x44,
      0x00, 0x5a, 0x40, 0x01, 0x18, 0x00, 0x70, 0x40, 0x01, 0x00};

  // Verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  // Verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);

  Ngap_NGAP_PDU_t decoded_pdu = {};
  uint16_t length = sizeof(initial_ue_message_hexbuf) / sizeof(unsigned char);

  bstring ngap_initial_ue_msg = blk2bstr(initial_ue_message_hexbuf, length);

  // Check if the pdu can be decoded
  ASSERT_EQ(ngap_amf_decode_pdu(&decoded_pdu, ngap_initial_ue_msg), RETURNok);

  xer_fprint(stdout, &asn_DEF_Ngap_NGAP_PDU, &decoded_pdu);

  // check if initial UE message is handled successfully
  EXPECT_EQ(ngap_amf_handle_message(state, peerInfo.assoc_id,
                                    peerInfo.instreams, &decoded_pdu),
            RETURNok);

  Ngap_InitialUEMessage_t* container;
  gnb_description_t* gNB_ref = NULL;
  m5g_ue_description_t* ue_ref = NULL;

  container =
      &(decoded_pdu.choice.initiatingMessage.value.choice.InitialUEMessage);
  Ngap_InitialUEMessage_IEs_t* ie = NULL;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_InitialUEMessage_IEs_t, ie,
                                      container,
                                      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID);

  // Check if Ran_UE_NGAP_ID is present in initial message
  ASSERT_TRUE(ie != NULL);

  // Check if gNB exists
  gNB_ref = ngap_state_get_gnb(state, peerInfo.assoc_id);
  ASSERT_TRUE(gNB_ref != NULL);

  gnb_ue_ngap_id_t gnb_ue_ngap_id = 0;
  gnb_ue_ngap_id = (gnb_ue_ngap_id_t)(ie->value.choice.RAN_UE_NGAP_ID);

  // Check for UE associated with gNB
  ue_ref = ngap_state_get_ue_gnbid(gNB_ref->sctp_assoc_id, gnb_ue_ngap_id);
  ASSERT_TRUE(ue_ref != NULL);

  // Check if UE is pointing to invalid ID if it is initial ue message
  EXPECT_EQ(ue_ref->amf_ue_ngap_id, INVALID_AMF_UE_NGAP_ID);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decoded_pdu);
  bdestroy(ngap_initial_ue_msg);
}

// Uplink Nas Transport With Auth Response
TEST_F(NgapFlowTest, uplink_nas_transport_sunny_day) {
  Ngap_UplinkNASTransport_t* container;
  gnb_description_t* gNB_ref = NULL;
  m5g_ue_description_t* ue_ref = NULL;

  unsigned char uplink_nas_transport_ue_message_hexbuf[] = {
      0x00, 0x2e, 0x40, 0x40, 0x00, 0x00, 0x04, 0x00, 0x0a, 0x00, 0x04, 0x40,
      0x01, 0x00, 0x02, 0x00, 0x55, 0x00, 0x04, 0x80, 0x01, 0x00, 0x01, 0x00,
      0x26, 0x00, 0x16, 0x15, 0x7e, 0x00, 0x57, 0x2d, 0x10, 0x8f, 0x17, 0xab,
      0x63, 0xde, 0x8b, 0xde, 0xba, 0x9a, 0x55, 0xe4, 0xc5, 0xdc, 0x12, 0xb1,
      0x54, 0x00, 0x79, 0x40, 0x0f, 0x40, 0x13, 0xf1, 0x84, 0x00, 0x02, 0x00,
      0x00, 0x00, 0x13, 0xf1, 0x84, 0x00, 0x00, 0x88};

  // verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  // verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);

  Ngap_NGAP_PDU_t decoded_pdu = {};
  uint16_t length =
      sizeof(uplink_nas_transport_ue_message_hexbuf) / sizeof(unsigned char);

  bstring uplink_nas_transport_msg =
      blk2bstr(uplink_nas_transport_ue_message_hexbuf, length);

  // Check if the uplink nas transport pdu is decoded successfully
  ASSERT_EQ(ngap_amf_decode_pdu(&decoded_pdu, uplink_nas_transport_msg),
            RETURNok);

  container =
      &(decoded_pdu.choice.initiatingMessage.value.choice.UplinkNASTransport);
  Ngap_UplinkNASTransport_IEs* ie = NULL;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_UplinkNASTransport_IEs, ie,
                                      container,
                                      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID);

  // Check if Ran_UE_NGAP_ID is present in Uplink nas transport
  ASSERT_TRUE(ie != NULL);

  // Check if gNB exists
  gNB_ref = ngap_state_get_gnb(state, peerInfo.assoc_id);
  ASSERT_TRUE(gNB_ref != NULL);

  gnb_ue_ngap_id_t gnb_ue_ngap_id = 0;
  amf_ue_ngap_id_t amf_ue_ngap_id = 0;
  gnb_ue_ngap_id = (gnb_ue_ngap_id_t)(ie->value.choice.RAN_UE_NGAP_ID);
  ue_ref = ngap_new_ue(state, peerInfo.assoc_id, gnb_ue_ngap_id);
  // Check if new UE is associated with gnb
  ASSERT_TRUE(ue_ref != NULL);
  ue_ref->ng_ue_state = NGAP_UE_CONNECTED;
  ue_ref->gnb_ue_ngap_id = gnb_ue_ngap_id;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_UplinkNASTransport_IEs, ie,
                                      container,
                                      Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID);

  // Check if AMF_UE_NGAP_ID present in Uplink_nas_transport
  ASSERT_TRUE(ie != NULL);
  asn_INTEGER2ulong(&ie->value.choice.AMF_UE_NGAP_ID,
                    reinterpret_cast<uint64_t*>(&amf_ue_ngap_id));
  ue_ref->amf_ue_ngap_id = amf_ue_ngap_id;

  // verify uplinkNasTransport handled successfully
  EXPECT_EQ(ngap_amf_handle_message(state, peerInfo.assoc_id,
                                    peerInfo.instreams, &decoded_pdu),
            RETURNok);

  // Check if UE is not invalid
  EXPECT_NE(ue_ref->amf_ue_ngap_id, INVALID_AMF_UE_NGAP_ID);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decoded_pdu);
  bdestroy(uplink_nas_transport_msg);
}

// Downlink Nas Transport Test with Auth Req
TEST_F(NgapFlowTest, downlink_nas_transport_auth_req_sunny_day) {
  unsigned char dl_nas_auth_req_msg[] = {
      0x7e, 0x00, 0x56, 0x00, 0x02, 0x00, 0x00, 0x21, 0xb4, 0x74, 0x3d,
      0x51, 0x76, 0xb8, 0xe5, 0x45, 0xe1, 0xdc, 0x03, 0x68, 0x25, 0x9a,
      0x67, 0x6c, 0x20, 0x10, 0xa4, 0x9b, 0x6b, 0x3d, 0x65, 0x6d, 0x80,
      0x00, 0x41, 0xc5, 0x72, 0x9e, 0xd9, 0xe1, 0xf0, 0xd6};
  MessageDef* message_p = NULL;
  m5g_ue_description_t* ue_ref = NULL;
  bstring buffer;
  unsigned int len = sizeof(dl_nas_auth_req_msg) / sizeof(unsigned char);
  // verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  // verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);

  message_p = itti_alloc_new_message(TASK_AMF_APP, NGAP_NAS_DL_DATA_REQ);

  NGAP_NAS_DL_DATA_REQ(message_p).gnb_ue_ngap_id = gNB_UE_NGAP_ID;
  NGAP_NAS_DL_DATA_REQ(message_p).amf_ue_ngap_id = AMF_UE_NGAP_ID;
  message_p->ittiMsgHeader.imsi = 0x311480000000001;
  buffer = bfromcstralloc(len, "\0");
  memcpy(buffer->data, dl_nas_auth_req_msg, len);
  buffer->slen = len;
  NGAP_NAS_DL_DATA_REQ(message_p).nas_msg = bstrcpy(buffer);

  ue_ref = ngap_new_ue(state, peerInfo.assoc_id, gNB_UE_NGAP_ID);
  // Check if new UE is associated with gnb
  ASSERT_TRUE(ue_ref != NULL);
  ue_ref->gnb_ue_ngap_id = gNB_UE_NGAP_ID;
  ue_ref->amf_ue_ngap_id = AMF_UE_NGAP_ID;

  // verify downlink nas transport is encoded and
  // handled successfully
  EXPECT_EQ(RETURNok, ngap_generate_downlink_nas_transport(
                          state, gNB_UE_NGAP_ID, AMF_UE_NGAP_ID,
                          &NGAP_NAS_DL_DATA_REQ(message_p).nas_msg,
                          message_p->ittiMsgHeader.imsi));
  itti_free_msg_content(message_p);
  free(message_p);
  bdestroy(buffer);
}

// Initial Context Setup Request
TEST_F(NgapFlowTest, initial_context_setup_request_sunny_day) {
  unsigned char reg_accept_msg[] = {
      0x7e, 0x02, 0x63, 0x26, 0x1f, 0x59, 0x01, 0x7e, 0x00, 0x42, 0x01,
      0x01, 0x77, 0x00, 0x0b, 0xf2, 0x13, 0xf1, 0x84, 0x01, 0x00, 0x40,
      0x41, 0x26, 0xca, 0x16, 0x54, 0x07, 0x00, 0x13, 0xf1, 0x84, 0x00,
      0x00, 0x01, 0x15, 0x02, 0x01, 0x01, 0x5e, 0x01, 0x06};
  bstring buffer;
  m5g_ue_description_t* ue_ref = NULL;
  unsigned int len = sizeof(reg_accept_msg) / sizeof(unsigned char);
  MessageDef* message_p = nullptr;
  Ngap_initial_context_setup_request_t* req = nullptr;

  // Verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  // Verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);

  message_p =
      itti_alloc_new_message(TASK_AMF_APP, NGAP_INITIAL_CONTEXT_SETUP_REQ);

  // Veirfy if memory allocated
  ASSERT_TRUE(message_p != NULL);

  req = &message_p->ittiMsg.ngap_initial_context_setup_req;
  memset(req, 0, sizeof(Ngap_initial_context_setup_request_t));
  req->amf_ue_ngap_id = AMF_UE_NGAP_ID;
  req->ran_ue_ngap_id = gNB_UE_NGAP_ID;
  buffer = bfromcstralloc(len, "\0");
  memcpy(buffer->data, reg_accept_msg, len);
  buffer->slen = len;
  req->nas_pdu = bstrcpy(buffer);
  req->ue_security_capabilities.m5g_encryption_algo = 0xC000;
  req->ue_security_capabilities.m5g_integrity_protection_algo = 0xC000;
  req->Security_Key = (unsigned char *)
      "e4298c66a3d368b59db6d6defa1b4ddaf40b9bef7cb398b44f468ea2f531ded3";
  message_p->ittiMsgHeader.imsi = 0x311480000000001;

  // Check if new UE is associated with gnb
  ue_ref = ngap_new_ue(state, peerInfo.assoc_id, AMF_UE_NGAP_ID);
  ASSERT_TRUE(ue_ref != NULL);
  ue_ref->ng_ue_state = NGAP_UE_CONNECTED;
  ue_ref->gnb_ue_ngap_id = gNB_UE_NGAP_ID;
  ue_ref->amf_ue_ngap_id = AMF_UE_NGAP_ID;

  /* Without SSD */
  ngap_handle_conn_est_cnf(state, &NGAP_INITIAL_CONTEXT_SETUP_REQ(message_p));

  /* With SSD */
  amf_config.plmn_support_list.plmn_support[0].s_nssai.sd.v = 0x1;
  ngap_handle_conn_est_cnf(state, &NGAP_INITIAL_CONTEXT_SETUP_REQ(message_p));

  /* Reset the SSD */
  amf_config.plmn_support_list.plmn_support[0].s_nssai.sd.v =
      AMF_S_NSSAI_SD_INVALID_VALUE;

  itti_free_msg_content(message_p);
  free(message_p);
  bdestroy(buffer);
}

// Initial Context Setup Response
TEST_F(NgapFlowTest, initial_context_setup_response_sunny_day) {
  Ngap_InitialContextSetupResponse_t* container = NULL;
  gnb_description_t* gNB_ref = NULL;
  m5g_ue_description_t* ue_ref = NULL;

  unsigned char initial_context_setup_response_message_hexbuf[] = {
      0x20, 0x0e, 0x00, 0x11, 0x00, 0x00, 0x02, 0x00, 0x0a, 0x40, 0x02,
      0x00, 0x05, 0x00, 0x55, 0x40, 0x04, 0x80, 0x01, 0x00, 0x01};

  // Verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  // Verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);

  Ngap_NGAP_PDU_t decoded_pdu = {};
  uint16_t length = sizeof(initial_context_setup_response_message_hexbuf) /
                    sizeof(unsigned char);

  bstring initial_context_response_succ_msg =
      blk2bstr(initial_context_setup_response_message_hexbuf, length);

  // Check if initial context setup response is decoded successfully
  ASSERT_EQ(
      ngap_amf_decode_pdu(&decoded_pdu, initial_context_response_succ_msg),
      RETURNok);

  container = &(decoded_pdu.choice.successfulOutcome.value.choice
                    .InitialContextSetupResponse);
  Ngap_InitialContextSetupResponseIEs* ie = NULL;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_InitialContextSetupResponseIEs, ie,
                                      container,
                                      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID);

  // Check if RAN_UE_NGAP_ID present in initial_context_setup_rsp
  ASSERT_TRUE(ie != NULL);

  // Check if gNB exists
  gNB_ref = ngap_state_get_gnb(state, peerInfo.assoc_id);
  ASSERT_TRUE(gNB_ref != NULL);

  gnb_ue_ngap_id_t gnb_ue_ngap_id = 0;
  amf_ue_ngap_id_t amf_ue_ngap_id = 0;
  gnb_ue_ngap_id = (gnb_ue_ngap_id_t)(ie->value.choice.RAN_UE_NGAP_ID);
  ue_ref = ngap_new_ue(state, peerInfo.assoc_id, gnb_ue_ngap_id);

  // Check if new UE is associated with gnb
  ASSERT_TRUE(ue_ref != NULL);
  ue_ref->gnb_ue_ngap_id = gnb_ue_ngap_id;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_InitialContextSetupResponseIEs, ie,
                                      container,
                                      Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID);

  // Check if AMF_UE_NGAP_ID is present in initial_context_setup_rsp
  ASSERT_TRUE(ie != NULL);
  asn_INTEGER2ulong(&ie->value.choice.AMF_UE_NGAP_ID,
                    reinterpret_cast<uint64_t*>(&amf_ue_ngap_id));
  ue_ref->amf_ue_ngap_id = amf_ue_ngap_id;

  // Verify initial_context_setup_rsp is encoded correctly
  EXPECT_EQ(ngap_amf_handle_message(state, peerInfo.assoc_id,
                                    peerInfo.instreams, &decoded_pdu),
            RETURNok);

  // Check if AMF_UE_NGAP_ID is not invalid
  EXPECT_NE(ue_ref->amf_ue_ngap_id, INVALID_AMF_UE_NGAP_ID);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decoded_pdu);
  bdestroy(initial_context_response_succ_msg);
}

// Ue Context Release Request
TEST_F(NgapFlowTest, ue_context_release_request_sunny_day) {
  Ngap_UEContextReleaseRequest_t* container = NULL;
  gnb_description_t* gNB_ref = NULL;
  m5g_ue_description_t* ue_ref = NULL;

  unsigned char ue_context_release_req_message_hexbuf[] = {
      0x00, 0x2a, 0x40, 0x17, 0x00, 0x00, 0x03, 0x00, 0x0a,
      0x00, 0x02, 0x00, 0x05, 0x00, 0x55, 0x00, 0x04, 0x80,
      0x01, 0x00, 0x01, 0x00, 0x0f, 0x40, 0x02, 0x05, 0x00};

  // Verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  // Verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);

  Ngap_NGAP_PDU_t decoded_pdu = {};
  uint16_t length =
      sizeof(ue_context_release_req_message_hexbuf) / sizeof(unsigned char);

  bstring ue_context_release_req_msg =
      blk2bstr(ue_context_release_req_message_hexbuf, length);

  // Check if ue_context_release_request is decoded successfully
  ASSERT_EQ(ngap_amf_decode_pdu(&decoded_pdu, ue_context_release_req_msg),
            RETURNok);

  container = &(decoded_pdu.choice.initiatingMessage.value.choice
                    .UEContextReleaseRequest);
  Ngap_UEContextReleaseRequest_IEs* ie = NULL;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_UEContextReleaseRequest_IEs, ie,
                                      container,
                                      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID);

  // Check if RAN_UE_NGAP_ID is present ue_context_release_request
  ASSERT_TRUE(ie != NULL);

  // Check if gNB exists
  gNB_ref = ngap_state_get_gnb(state, peerInfo.assoc_id);
  ASSERT_TRUE(gNB_ref != NULL);

  gnb_ue_ngap_id_t gnb_ue_ngap_id = 0;
  amf_ue_ngap_id_t amf_ue_ngap_id = 0;
  gnb_ue_ngap_id = (gnb_ue_ngap_id_t)(ie->value.choice.RAN_UE_NGAP_ID);

  // Check if new UE is associated with gnb
  ue_ref = ngap_new_ue(state, peerInfo.assoc_id, gnb_ue_ngap_id);
  ASSERT_TRUE(ue_ref != NULL);
  ue_ref->gnb_ue_ngap_id = gnb_ue_ngap_id;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_UEContextReleaseRequest_IEs, ie,
                                      container,
                                      Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID);

  // Check if AMF_UE_NGAP_ID ie present in ue_context_release_request
  ASSERT_TRUE(ie != NULL);
  asn_INTEGER2ulong(&ie->value.choice.AMF_UE_NGAP_ID,
                    reinterpret_cast<uint64_t*>(&amf_ue_ngap_id));
  ue_ref->amf_ue_ngap_id = amf_ue_ngap_id;

  // check if Ngap_Cause IE present
  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_UEContextReleaseRequest_IEs, ie,
                                      container, Ngap_ProtocolIE_ID_id_Cause);
  ASSERT_TRUE(ie != NULL);

  // Verify if UE_Context_Release_Request is handled properly
  EXPECT_EQ(ngap_amf_handle_message(state, peerInfo.assoc_id,
                                    peerInfo.instreams, &decoded_pdu),
            RETURNok);

  // Check if UE is not invalid
  EXPECT_NE(ue_ref->amf_ue_ngap_id, INVALID_AMF_UE_NGAP_ID);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decoded_pdu);
  bdestroy(ue_context_release_req_msg);
}

// Ue Context Release Complete
TEST_F(NgapFlowTest, ue_context_release_complete_sunny_day) {
  Ngap_UEContextReleaseComplete_t* container = NULL;
  gnb_description_t* gNB_ref = NULL;
  m5g_ue_description_t* ue_ref = NULL;

  unsigned char ue_context_release_com_message_hexbuf[] = {
      0x20, 0x29, 0x00, 0x11, 0x00, 0x00, 0x02, 0x00, 0x0a, 0x40, 0x02,
      0x00, 0x05, 0x00, 0x55, 0x40, 0x04, 0x80, 0x01, 0x00, 0x01};

  // Verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  // Verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);

  Ngap_NGAP_PDU_t decoded_pdu = {};
  uint16_t length =
      sizeof(ue_context_release_com_message_hexbuf) / sizeof(unsigned char);

  bstring ue_context_release_com_msg =
      blk2bstr(ue_context_release_com_message_hexbuf, length);

  // Check if the ue_context_release_complete decoded successfully
  ASSERT_EQ(ngap_amf_decode_pdu(&decoded_pdu, ue_context_release_com_msg),
            RETURNok);

  container = &(decoded_pdu.choice.successfulOutcome.value.choice
                    .UEContextReleaseComplete);
  Ngap_UEContextReleaseComplete_IEs* ie = NULL;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_UEContextReleaseComplete_IEs, ie,
                                      container,
                                      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID);

  // Check if RAN_UE_NGAP_ID present
  ASSERT_TRUE(ie != NULL);

  // Check if gNB exists
  gNB_ref = ngap_state_get_gnb(state, peerInfo.assoc_id);
  ASSERT_TRUE(gNB_ref != NULL);

  gnb_ue_ngap_id_t gnb_ue_ngap_id = 0;
  amf_ue_ngap_id_t amf_ue_ngap_id = 0;
  gnb_ue_ngap_id = (gnb_ue_ngap_id_t)(ie->value.choice.RAN_UE_NGAP_ID);
  ue_ref = ngap_new_ue(state, peerInfo.assoc_id, gnb_ue_ngap_id);

  // Check if new UE is associated with gnb
  ASSERT_TRUE(ue_ref != NULL);
  ue_ref->gnb_ue_ngap_id = gnb_ue_ngap_id;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_UEContextReleaseComplete_IEs, ie,
                                      container,
                                      Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID);

  // Check if AMF_UE_NGAP_ID present
  ASSERT_TRUE(ie != NULL);
  asn_INTEGER2ulong(&ie->value.choice.AMF_UE_NGAP_ID,
                    reinterpret_cast<uint64_t*>(&amf_ue_ngap_id));
  ue_ref->amf_ue_ngap_id = amf_ue_ngap_id;

  // Verify ue_context_release_complete is handled successfully
  EXPECT_EQ(ngap_amf_handle_message(state, peerInfo.assoc_id,
                                    peerInfo.instreams, &decoded_pdu),
            RETURNok);

  // Check if UE is not invalid
  EXPECT_NE(ue_ref->amf_ue_ngap_id, INVALID_AMF_UE_NGAP_ID);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decoded_pdu);
  bdestroy(ue_context_release_com_msg);
}

// UE Context Release Command
TEST_F(NgapFlowTest, ue_context_release_command_sunny_day) {
  MessageDef* message_p = NULL;
  m5g_ue_description_t* ue_ref = NULL;

  // Verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  // Verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);

  message_p =
      itti_alloc_new_message(TASK_AMF_APP, NGAP_UE_CONTEXT_RELEASE_COMMAND);

  // verify memory allocation is successful
  ASSERT_TRUE(message_p != NULL);

  // Check if new UE is associated with gnb
  ue_ref = ngap_new_ue(state, peerInfo.assoc_id, gNB_UE_NGAP_ID);

  ASSERT_TRUE(ue_ref != NULL);
  ue_ref->ng_ue_state = NGAP_UE_CONNECTED;
  ue_ref->gnb_ue_ngap_id = gNB_UE_NGAP_ID;
  ue_ref->amf_ue_ngap_id = AMF_UE_NGAP_ID;

  message_p->ittiMsgHeader.imsi = 0x311480000000001;
  NGAP_UE_CONTEXT_RELEASE_COMMAND(message_p).amf_ue_ngap_id = AMF_UE_NGAP_ID;
  NGAP_UE_CONTEXT_RELEASE_COMMAND(message_p).gnb_ue_ngap_id = gNB_UE_NGAP_ID;
  NGAP_UE_CONTEXT_RELEASE_COMMAND(message_p).cause = NGAP_USER_INACTIVITY;

  // verify ue_context_release_command is encoded correctly
  EXPECT_EQ(RETURNok, ngap_handle_ue_context_release_command(
                          state, &NGAP_UE_CONTEXT_RELEASE_COMMAND(message_p),
                          message_p->ittiMsgHeader.imsi));
  itti_free_msg_content(message_p);
  free(message_p);
}

// Pdu Session Resource Setup Request
TEST_F(NgapFlowTest, pdu_sess_resource_setup_req_sunny_day) {
  MessageDef* message_p = NULL;
  uint8_t ip_buff[4] = {0xc0, 0xa8, 0x3c, 0x8e};
  m5g_ue_description_t* ue_ref = NULL;
  itti_ngap_pdusession_resource_setup_req_t* ngap_pdu_ses_setup_req = nullptr;
  pdu_session_resource_setup_request_transfer_t amf_pdu_ses_setup_transfer_req =
      {};

  // Verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  // Verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);

  message_p =
      itti_alloc_new_message(TASK_AMF_APP, NGAP_PDUSESSION_RESOURCE_SETUP_REQ);
  message_p->ittiMsgHeader.imsi = 0x311480000000001;

  ngap_pdu_ses_setup_req =
      &message_p->ittiMsg.ngap_pdusession_resource_setup_req;
  memset(ngap_pdu_ses_setup_req, 0,
         sizeof(itti_ngap_pdusession_resource_setup_req_t));

  ngap_pdu_ses_setup_req->gnb_ue_ngap_id = gNB_UE_NGAP_ID;
  ngap_pdu_ses_setup_req->amf_ue_ngap_id = AMF_UE_NGAP_ID;

  ngap_pdu_ses_setup_req->ue_aggregate_maximum_bit_rate.dl = 2048;
  ngap_pdu_ses_setup_req->ue_aggregate_maximum_bit_rate.ul = 2048;

  ngap_pdu_ses_setup_req->pduSessionResource_setup_list.no_of_items = 1;
  ngap_pdu_ses_setup_req->pduSessionResource_setup_list.item[0].Pdu_Session_ID =
      0x05;

  amf_pdu_ses_setup_transfer_req.pdu_aggregate_max_bit_rate.dl = 1024;
  amf_pdu_ses_setup_transfer_req.pdu_aggregate_max_bit_rate.ul = 1024;
  amf_pdu_ses_setup_transfer_req.up_transport_layer_info.gtp_tnl.gtp_tied[0] =
      0x80;
  amf_pdu_ses_setup_transfer_req.up_transport_layer_info.gtp_tnl.gtp_tied[1] =
      0x00;
  amf_pdu_ses_setup_transfer_req.up_transport_layer_info.gtp_tnl.gtp_tied[2] =
      0x00;
  amf_pdu_ses_setup_transfer_req.up_transport_layer_info.gtp_tnl.gtp_tied[3] =
      0x01;
  amf_pdu_ses_setup_transfer_req.up_transport_layer_info.gtp_tnl
      .endpoint_ip_address = blk2bstr(ip_buff, 4);
  amf_pdu_ses_setup_transfer_req.pdu_ip_type.pdn_type = IPv4;

  amf_pdu_ses_setup_transfer_req.qos_flow_add_or_mod_request_list
      .maxNumOfQosFlows = 1;
  amf_pdu_ses_setup_transfer_req.qos_flow_add_or_mod_request_list.item[0]
      .qos_flow_req_item.qos_flow_identifier = 5;
  amf_pdu_ses_setup_transfer_req.qos_flow_add_or_mod_request_list.item[0]
      .qos_flow_req_item.qos_flow_level_qos_param.qos_characteristic
      .non_dynamic_5QI_desc.fiveQI = 9;
  amf_pdu_ses_setup_transfer_req.qos_flow_add_or_mod_request_list.item[0]
      .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
      .priority_level = 1;
  amf_pdu_ses_setup_transfer_req.qos_flow_add_or_mod_request_list.item[0]
      .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
      .pre_emption_cap = SHALL_NOT_TRIGGER_PRE_EMPTION;
  amf_pdu_ses_setup_transfer_req.qos_flow_add_or_mod_request_list.item[0]
      .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
      .pre_emption_vul = NOT_PREEMPTABLE;
  ngap_pdu_ses_setup_req->pduSessionResource_setup_list.item[0]
      .PDU_Session_Resource_Setup_Request_Transfer =
      amf_pdu_ses_setup_transfer_req;

  // verify new ue is associated with gnb
  ue_ref = ngap_new_ue(state, peerInfo.assoc_id, gNB_UE_NGAP_ID);
  ASSERT_TRUE(ue_ref != NULL);
  ue_ref->ng_ue_state = NGAP_UE_CONNECTED;
  ue_ref->gnb_ue_ngap_id = gNB_UE_NGAP_ID;
  ue_ref->amf_ue_ngap_id = AMF_UE_NGAP_ID;

  // verify pdu_session_resource_setup_request is encoded correctly
  EXPECT_EQ(RETURNok,
            ngap_generate_ngap_pdusession_resource_setup_req(
                state, &NGAP_PDUSESSION_RESOURCE_SETUP_REQ(message_p)));
  itti_free_msg_content(message_p);
  free(message_p);
}

// Pdu Session Resource Setup Response
TEST_F(NgapFlowTest, pdu_session_resource_setup_resp_sunny_day) {
  Ngap_PDUSessionResourceSetupResponse_t* container = NULL;
  gnb_description_t* gNB_ref = NULL;
  m5g_ue_description_t* ue_ref = NULL;
  Ngap_NGAP_PDU_t decoded_pdu = {};
  uint8_t pdu_ss_resource_setup_resp_hex_buff[] = {
      0x20, 0x1d, 0x00, 0x27, 0x00, 0x00, 0x03, 0x00, 0x0a, 0x40, 0x02,
      0x00, 0x05, 0x00, 0x55, 0x40, 0x04, 0x80, 0x01, 0x00, 0x01, 0x00,
      0x4b, 0x40, 0x12, 0x00, 0x00, 0x05, 0x0e, 0x00, 0x03, 0xe0, 0x05,
      0x05, 0x05, 0x02, 0x00, 0x00, 0x00, 0x0b, 0x01, 0x00, 0x00};
  uint16_t len = sizeof(pdu_ss_resource_setup_resp_hex_buff) / sizeof(uint8_t);

  // Verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  // Verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);

  bstring pdu_ss_resource_response_succ_msg =
      blk2bstr(pdu_ss_resource_setup_resp_hex_buff, len);

  // Check if the pdu_session_resource_setup_response decoded successfully
  ASSERT_EQ(
      ngap_amf_decode_pdu(&decoded_pdu, pdu_ss_resource_response_succ_msg),
      RETURNok);
  container = &(decoded_pdu.choice.successfulOutcome.value.choice
                    .PDUSessionResourceSetupResponse);
  Ngap_PDUSessionResourceSetupResponseIEs_t* ie = NULL;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_PDUSessionResourceSetupResponseIEs_t,
                                      ie, container,
                                      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID);

  // Check if RAN_UE_NGAP_ID Present
  ASSERT_TRUE(ie != NULL);

  // Check if gNB exists
  gNB_ref = ngap_state_get_gnb(state, peerInfo.assoc_id);
  ASSERT_TRUE(gNB_ref != NULL);

  gnb_ue_ngap_id_t gnb_ue_ngap_id = 0;
  amf_ue_ngap_id_t amf_ue_ngap_id = 0;
  gnb_ue_ngap_id = (gnb_ue_ngap_id_t)(ie->value.choice.RAN_UE_NGAP_ID);

  // verify new ue is associated with gnb
  ue_ref = ngap_new_ue(state, peerInfo.assoc_id, gnb_ue_ngap_id);
  ASSERT_TRUE(ue_ref != NULL);
  ue_ref->gnb_ue_ngap_id = gnb_ue_ngap_id;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_PDUSessionResourceSetupResponseIEs_t,
                                      ie, container,
                                      Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID);

  // Check if AMF_UE_NGAP_ID present
  ASSERT_TRUE(ie != NULL);
  asn_INTEGER2ulong(&ie->value.choice.AMF_UE_NGAP_ID,
                    reinterpret_cast<uint64_t*>(&amf_ue_ngap_id));
  ue_ref->amf_ue_ngap_id = amf_ue_ngap_id;

  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_PDUSessionResourceSetupResponseIEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_PDUSessionResourceSetupListSURes, false);
  // verify if Ngap_ProtocolIE_ID_id_PDUSessionResourceSetupListSURes present
  ASSERT_TRUE(ie != NULL);

  // verify pdu_session_resource_setup_response is handled_correctly
  EXPECT_EQ(ngap_amf_handle_message(state, peerInfo.assoc_id,
                                    peerInfo.instreams, &decoded_pdu),
            RETURNok);

  // Check if UE is not invalid ID
  EXPECT_NE(ue_ref->amf_ue_ngap_id, INVALID_AMF_UE_NGAP_ID);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decoded_pdu);
  bdestroy(pdu_ss_resource_response_succ_msg);
}

// Pdu Session Resource Release command
TEST_F(NgapFlowTest, pdu_sess_resource_rel_cmd_sunny_day) {
  unsigned char pdu_sess_resource_rel_cmd[] = {
      0x7e, 0x00, 0x56, 0x00, 0x02, 0x00, 0x00, 0x21, 0xb4, 0x74, 0x3d,
      0x51, 0x76, 0xb8, 0xe5, 0x45, 0xe1, 0xdc, 0x03, 0x68, 0x25, 0x9a,
      0x67, 0x6c, 0x20, 0x10, 0xa4, 0x9b, 0x6b, 0x3d, 0x65, 0x6d, 0x80,
      0x00, 0x41, 0xc5, 0x72, 0x9e, 0xd9, 0xe1, 0xf0, 0xd6};
  MessageDef* message_p = NULL;
  m5g_ue_description_t* ue_ref = NULL;
  bstring buffer;
  unsigned int len = sizeof(pdu_sess_resource_rel_cmd) / sizeof(unsigned char);

  // Verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  // Verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);

  message_p =
      itti_alloc_new_message(TASK_AMF_APP, NGAP_PDUSESSIONRESOURCE_REL_REQ);

  NGAP_PDUSESSIONRESOURCE_REL_REQ(message_p).gnb_ue_ngap_id = gNB_UE_NGAP_ID;
  NGAP_PDUSESSIONRESOURCE_REL_REQ(message_p).amf_ue_ngap_id = AMF_UE_NGAP_ID;
  message_p->ittiMsgHeader.imsi = 0x311480000000001;
  buffer = bfromcstralloc(len, "\0");
  memcpy(buffer->data, pdu_sess_resource_rel_cmd, len);
  buffer->slen = len;
  NGAP_PDUSESSIONRESOURCE_REL_REQ(message_p).nas_msg = bstrcpy(buffer);

  // verify new ue is associated with gnb
  ue_ref = ngap_new_ue(state, peerInfo.assoc_id, gNB_UE_NGAP_ID);
  ASSERT_TRUE(ue_ref != NULL);
  ue_ref->ng_ue_state = NGAP_UE_CONNECTED;
  ue_ref->gnb_ue_ngap_id = gNB_UE_NGAP_ID;
  ue_ref->amf_ue_ngap_id = AMF_UE_NGAP_ID;

  // verify pdu_session_resource_release_command handled
  // and encoded corrrectly
  EXPECT_EQ(RETURNok, ngap_generate_ngap_pdusession_resource_rel_cmd(
                          state, &NGAP_PDUSESSIONRESOURCE_REL_REQ(message_p)));
  itti_free_msg_content(message_p);
  free(message_p);
  bdestroy(buffer);
}

// Initial Context Setup  failure
TEST_F(NgapFlowTest, initial_context_setup_failure_rainy_day) {
  Ngap_InitialContextSetupFailure_t* container = NULL;
  gnb_description_t* gNB_ref = NULL;
  m5g_ue_description_t* ue_ref = NULL;

  unsigned char initial_context_setup_failure_message_hexbuf[] = {
      0x40, 0x0e, 0x00, 0x21, 0x00, 0x00, 0x04, 0x00, 0x0a, 0x40,
      0x02, 0x00, 0x02, 0x00, 0x55, 0x40, 0x04, 0x80, 0x01, 0x00,
      0x02, 0x00, 0x84, 0x40, 0x06, 0x00, 0x00, 0x05, 0x02, 0x00,
      0xe0, 0x00, 0x0f, 0x40, 0x02, 0x00, 0x40};

  // Verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  // Verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  EXPECT_EQ(state->gnbs.num_elements, 1);

  Ngap_NGAP_PDU_t decoded_pdu = {};
  uint16_t length = sizeof(initial_context_setup_failure_message_hexbuf) /
                    sizeof(unsigned char);

  bstring initial_context_setup_fail_msg =
      blk2bstr(initial_context_setup_failure_message_hexbuf, length);

  // Check if the initial_context_setup_failure decoded
  ASSERT_EQ(ngap_amf_decode_pdu(&decoded_pdu, initial_context_setup_fail_msg),
            RETURNok);

  container = &(decoded_pdu.choice.unsuccessfulOutcome.value.choice
                    .InitialContextSetupFailure);
  Ngap_InitialContextSetupFailureIEs* ie = NULL;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_InitialContextSetupFailureIEs, ie,
                                      container,
                                      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID);

  // Check if Ran_UE_NGAP_ID
  ASSERT_TRUE(ie != NULL);

  // Check if gNB exists
  gNB_ref = ngap_state_get_gnb(state, peerInfo.assoc_id);
  ASSERT_TRUE(gNB_ref != NULL);

  gnb_ue_ngap_id_t gnb_ue_ngap_id = 0;
  amf_ue_ngap_id_t amf_ue_ngap_id = 0;
  gnb_ue_ngap_id = (gnb_ue_ngap_id_t)(ie->value.choice.RAN_UE_NGAP_ID);

  // verify new ue is associated with gnb
  ue_ref = ngap_new_ue(state, peerInfo.assoc_id, gnb_ue_ngap_id);
  ASSERT_TRUE(ue_ref != NULL);
  ue_ref->gnb_ue_ngap_id = gnb_ue_ngap_id;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_InitialContextSetupFailureIEs, ie,
                                      container,
                                      Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID);

  // Check if AMF_UE_NGAP_ID present
  ASSERT_TRUE(ie != NULL);
  asn_INTEGER2ulong(&ie->value.choice.AMF_UE_NGAP_ID,
                    reinterpret_cast<uint64_t*>(&amf_ue_ngap_id));
  ue_ref->amf_ue_ngap_id = amf_ue_ngap_id;

  // check if NGAP_Cause IE present
  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_InitialContextSetupFailureIEs, ie,
                                      container, Ngap_ProtocolIE_ID_id_Cause);
  ASSERT_TRUE(ie != NULL);

  // vefiry initial_context_setup_failure is handled
  EXPECT_EQ(ngap_amf_handle_message(state, peerInfo.assoc_id,
                                    peerInfo.instreams, &decoded_pdu),
            RETURNok);

  // Check if UE is not invalid
  EXPECT_NE(ue_ref->amf_ue_ngap_id, INVALID_AMF_UE_NGAP_ID);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decoded_pdu);
  bdestroy(initial_context_setup_fail_msg);
}

// Uplink Nas Transport With Auth Response
TEST_F(NgapFlowTest, uplink_nas_transport_rainy_day) {
  Ngap_UplinkNASTransport_t* container;
  gnb_description_t* gNB_ref = NULL;
  m5g_ue_description_t* ue_ref = NULL;

  unsigned char uplink_nas_transport_ue_message_hexbuf[] = {
      0x00, 0x2e, 0x40, 0x40, 0x00, 0x00, 0x04, 0x00, 0x0a, 0x00, 0x04, 0x40,
      0x01, 0x00, 0x02, 0x00, 0x55, 0x00, 0x04, 0x80, 0x01, 0x00, 0x01, 0x00,
      0x26, 0x00, 0x16, 0x15, 0x7e, 0x00, 0x57, 0x2d, 0x10, 0x8f, 0x17, 0xab,
      0x63, 0xde, 0x8b, 0xde, 0xba, 0x9a, 0x55, 0xe4, 0xc5, 0xdc, 0x12, 0xb1,
      0x54, 0x00, 0x79, 0x40, 0x0f, 0x40, 0x13, 0xf1, 0x84, 0x00, 0x02, 0x00,
      0x00, 0x00, 0x13, 0xf1, 0x84, 0x00, 0x00, 0x88};

  // Verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  // Verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);

  Ngap_NGAP_PDU_t decoded_pdu = {};
  uint16_t length =
      sizeof(uplink_nas_transport_ue_message_hexbuf) / sizeof(unsigned char);

  bstring uplink_nas_transport_msg =
      blk2bstr(uplink_nas_transport_ue_message_hexbuf, length);

  // Check if the uplink_nas_transport_msg decoded
  ASSERT_EQ(ngap_amf_decode_pdu(&decoded_pdu, uplink_nas_transport_msg),
            RETURNok);

  container =
      &(decoded_pdu.choice.initiatingMessage.value.choice.UplinkNASTransport);
  Ngap_UplinkNASTransport_IEs* ie = NULL;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_UplinkNASTransport_IEs, ie,
                                      container,
                                      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID);

  // Check if RAN_UE_NGAP_ID present
  ASSERT_TRUE(ie != NULL);

  // Check if gNB exists
  gNB_ref = ngap_state_get_gnb(state, peerInfo.assoc_id);
  ASSERT_TRUE(gNB_ref != NULL);

  gnb_ue_ngap_id_t gnb_ue_ngap_id = 0;
  amf_ue_ngap_id_t amf_ue_ngap_id = 0;
  gnb_ue_ngap_id = (gnb_ue_ngap_id_t)(ie->value.choice.RAN_UE_NGAP_ID);
  // verify new ue is associated with gnb
  ue_ref = ngap_new_ue(state, peerInfo.assoc_id, gnb_ue_ngap_id);
  ASSERT_TRUE(ue_ref != NULL);
  ue_ref->gnb_ue_ngap_id = 0xffff;

  /* verify uplink_nas_transport is not handled as gnb_ue_ngap_id present
   * in ue_ref is not equal to the IE RAN_UE_NGAP_ID present in
   * uplink_nas_transport
   */
  EXPECT_EQ(ngap_amf_handle_message(state, peerInfo.assoc_id,
                                    peerInfo.instreams, &decoded_pdu),
            RETURNerror);

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_UplinkNASTransport_IEs, ie,
                                      container,
                                      Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID);

  // Check if AMF_UE_NGAP_ID present
  ASSERT_TRUE(ie != NULL);
  asn_INTEGER2ulong(&ie->value.choice.AMF_UE_NGAP_ID,
                    reinterpret_cast<uint64_t*>(&amf_ue_ngap_id));
  ue_ref->amf_ue_ngap_id = 0xffff;

  // verify uplink_nas_transport is not handled as amf_ue_ngap_id present
  // in ue_ref is not equal to IE AMF_UE_NGAP_ID present in uplink_nas_transport
  EXPECT_EQ(ngap_amf_handle_message(state, peerInfo.assoc_id,
                                    peerInfo.instreams, &decoded_pdu),
            RETURNerror);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decoded_pdu);
  bdestroy(uplink_nas_transport_msg);
}

// Downlink Nas Transport Test with Auth Req
TEST_F(NgapFlowTest, downlink_nas_transport_auth_rainy_day) {
  unsigned char dl_nas_auth_req_msg[] = {
      0x7e, 0x00, 0x56, 0x00, 0x02, 0x00, 0x00, 0x21, 0xb4, 0x74, 0x3d,
      0x51, 0x76, 0xb8, 0xe5, 0x45, 0xe1, 0xdc, 0x03, 0x68, 0x25, 0x9a,
      0x67, 0x6c, 0x20, 0x10, 0xa4, 0x9b, 0x6b, 0x3d, 0x65, 0x6d, 0x80,
      0x00, 0x41, 0xc5, 0x72, 0x9e, 0xd9, 0xe1, 0xf0, 0xd6};
  MessageDef* message_p = NULL;
  m5g_ue_description_t* ue_ref = NULL;
  bstring buffer;
  unsigned int len = sizeof(dl_nas_auth_req_msg) / sizeof(unsigned char);

  // Verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  // Verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);

  message_p = itti_alloc_new_message(TASK_AMF_APP, NGAP_NAS_DL_DATA_REQ);

  NGAP_NAS_DL_DATA_REQ(message_p).gnb_ue_ngap_id = gNB_UE_NGAP_ID;
  NGAP_NAS_DL_DATA_REQ(message_p).amf_ue_ngap_id = AMF_UE_NGAP_ID;
  message_p->ittiMsgHeader.imsi = 0x311480000000001;
  buffer = bfromcstralloc(len, "\0");
  memcpy(buffer->data, dl_nas_auth_req_msg, len);
  buffer->slen = len;
  NGAP_NAS_DL_DATA_REQ(message_p).nas_msg = bstrcpy(buffer);

  ue_ref = ngap_new_ue(state, peerInfo.assoc_id, gNB_UE_NGAP_ID);
  ASSERT_TRUE(ue_ref != NULL);
  ue_ref->gnb_ue_ngap_id = 0xffff;
  ue_ref->amf_ue_ngap_id = 0xffff;

  // verify downlink_nas_transport is not handled as gnb_ue_ngap_id,
  // amf_ue_nagp_id  in ue_ref does not match with ies present in
  // downlink_nas_transport
  EXPECT_EQ(RETURNerror, ngap_generate_downlink_nas_transport(
                             state, gNB_UE_NGAP_ID, AMF_UE_NGAP_ID,
                             &NGAP_NAS_DL_DATA_REQ(message_p).nas_msg,
                             message_p->ittiMsgHeader.imsi));
  itti_free_msg_content(message_p);
  free(message_p);
  bdestroy(buffer);
}

// Initial Context Setup Request
TEST_F(NgapFlowTest, initial_context_setup_request_rainy_day) {
  m5g_ue_description_t* ue_ref = NULL;
  MessageDef* message_p = nullptr;
  Ngap_initial_context_setup_request_t* req = nullptr;

  // Verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  // Verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);

  message_p =
      itti_alloc_new_message(TASK_AMF_APP, NGAP_INITIAL_CONTEXT_SETUP_REQ);
  ASSERT_TRUE(message_p != NULL);
  req = &message_p->ittiMsg.ngap_initial_context_setup_req;
  memset(req, 0, sizeof(Ngap_initial_context_setup_request_t));
  req->amf_ue_ngap_id = AMF_UE_NGAP_ID;
  req->ran_ue_ngap_id = gNB_UE_NGAP_ID;
  req->ue_security_capabilities.m5g_encryption_algo = 0xC000;
  req->ue_security_capabilities.m5g_integrity_protection_algo = 0xC000;
  req->Security_Key =
      (unsigned char *)
      "e4298c66a3d368b59db6d6defa1b4ddaf40b9bef7cb398b44f468ea2f531ded3";
  message_p->ittiMsgHeader.imsi = 0x311480000000001;

  // verify new ue is associated with gnb
  ue_ref = ngap_new_ue(state, peerInfo.assoc_id, AMF_UE_NGAP_ID);
  ASSERT_TRUE(ue_ref != NULL);
  ue_ref->ng_ue_state = NGAP_UE_CONNECTED;
  ue_ref->gnb_ue_ngap_id = 0xffff;
  ue_ref->amf_ue_ngap_id = 0xffff;

  ngap_handle_conn_est_cnf(state, &NGAP_INITIAL_CONTEXT_SETUP_REQ(message_p));
  itti_free_msg_content(message_p);
  free(message_p);
}

// Initial Context Setup Response
TEST_F(NgapFlowTest, initial_context_setup_response_rainy_day) {
  Ngap_InitialContextSetupResponse_t* container = NULL;
  gnb_description_t* gNB_ref = NULL;
  m5g_ue_description_t* ue_ref = NULL;

  unsigned char initial_context_setup_response_message_hexbuf[] = {
      0x20, 0x0e, 0x00, 0x11, 0x00, 0x00, 0x02, 0x00, 0x0a, 0x40, 0x02,
      0x00, 0x05, 0x00, 0x55, 0x40, 0x04, 0x80, 0x01, 0x00, 0x01};
  // Verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  // Verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);

  Ngap_NGAP_PDU_t decoded_pdu = {};
  uint16_t length = sizeof(initial_context_setup_response_message_hexbuf) /
                    sizeof(unsigned char);

  bstring initial_context_response_succ_msg =
      blk2bstr(initial_context_setup_response_message_hexbuf, length);

  // Check if initial_context_setup_response_message is decoded
  ASSERT_EQ(
      ngap_amf_decode_pdu(&decoded_pdu, initial_context_response_succ_msg),
      RETURNok);

  container = &(decoded_pdu.choice.successfulOutcome.value.choice
                    .InitialContextSetupResponse);
  Ngap_InitialContextSetupResponseIEs* ie = NULL;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_InitialContextSetupResponseIEs, ie,
                                      container,
                                      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID);

  // Check if RAN_UE_NGAP_ID is present
  ASSERT_TRUE(ie != NULL);

  // Check if gNB exists
  gNB_ref = ngap_state_get_gnb(state, peerInfo.assoc_id);
  ASSERT_TRUE(gNB_ref != NULL);

  gnb_ue_ngap_id_t gnb_ue_ngap_id = 0;
  amf_ue_ngap_id_t amf_ue_ngap_id = 0;
  gnb_ue_ngap_id = (gnb_ue_ngap_id_t)(ie->value.choice.RAN_UE_NGAP_ID);
  // verify new ue is associated with gnb
  ue_ref = ngap_new_ue(state, peerInfo.assoc_id, gnb_ue_ngap_id);
  ASSERT_TRUE(ue_ref != NULL);
  ue_ref->gnb_ue_ngap_id = gNB_UE_NGAP_ID;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_InitialContextSetupResponseIEs, ie,
                                      container,
                                      Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID);

  // Check if AMF_UE_NGAP_ID present
  ASSERT_TRUE(ie != NULL);
  asn_INTEGER2ulong(&ie->value.choice.AMF_UE_NGAP_ID,
                    reinterpret_cast<uint64_t*>(&amf_ue_ngap_id));
  ue_ref->amf_ue_ngap_id = amf_ue_ngap_id;

  // verify initial_context_setup_response is not handled as gnb_ue_ngap_id
  // present in ue_ref than is not equal to IE RAN_UE_NGAP_ID present in
  // initial_context_setup_response
  EXPECT_EQ(ngap_amf_handle_message(state, peerInfo.assoc_id,
                                    peerInfo.instreams, &decoded_pdu),
            RETURNerror);

  // Check if UE is not invalid
  EXPECT_NE(ue_ref->amf_ue_ngap_id, INVALID_AMF_UE_NGAP_ID);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decoded_pdu);
  bdestroy(initial_context_response_succ_msg);
}

// Ue Context Release Request
TEST_F(NgapFlowTest, ue_context_release_request_rainy_day) {
  Ngap_UEContextReleaseRequest_t* container = NULL;
  gnb_description_t* gNB_ref = NULL;
  m5g_ue_description_t* ue_ref = NULL;

  unsigned char ue_context_release_req_message_hexbuf[] = {
      0x00, 0x2a, 0x40, 0x17, 0x00, 0x00, 0x03, 0x00, 0x0a,
      0x00, 0x02, 0x00, 0x05, 0x00, 0x55, 0x00, 0x04, 0x80,
      0x01, 0x00, 0x01, 0x00, 0x0f, 0x40, 0x02, 0x05, 0x00};
  // Verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  // Verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);

  Ngap_NGAP_PDU_t decoded_pdu = {};
  uint16_t length =
      sizeof(ue_context_release_req_message_hexbuf) / sizeof(unsigned char);

  bstring ue_context_release_req_msg =
      blk2bstr(ue_context_release_req_message_hexbuf, length);

  // Check if ue_context_release_request is decoded
  ASSERT_EQ(ngap_amf_decode_pdu(&decoded_pdu, ue_context_release_req_msg),
            RETURNok);

  container = &(decoded_pdu.choice.initiatingMessage.value.choice
                    .UEContextReleaseRequest);
  Ngap_UEContextReleaseRequest_IEs* ie = NULL;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_UEContextReleaseRequest_IEs, ie,
                                      container,
                                      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID);

  // Check if RAN_UE_NGAP_ID is present
  ASSERT_TRUE(ie != NULL);

  // Check if gNB exists
  gNB_ref = ngap_state_get_gnb(state, peerInfo.assoc_id);
  ASSERT_TRUE(gNB_ref != NULL);

  gnb_ue_ngap_id_t gnb_ue_ngap_id = 0;
  amf_ue_ngap_id_t amf_ue_ngap_id = 0;
  gnb_ue_ngap_id = (gnb_ue_ngap_id_t)(ie->value.choice.RAN_UE_NGAP_ID);

  // verify new ue is associated with gnb
  ue_ref = ngap_new_ue(state, peerInfo.assoc_id, gnb_ue_ngap_id);
  ASSERT_TRUE(ue_ref != NULL);

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_UEContextReleaseRequest_IEs, ie,
                                      container,
                                      Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID);

  // Check if AMF_UE_NGAP_ID present
  ASSERT_TRUE(ie != NULL);
  asn_INTEGER2ulong(&ie->value.choice.AMF_UE_NGAP_ID,
                    reinterpret_cast<uint64_t*>(&amf_ue_ngap_id));
  // verify ue_context_release_request is not handled as gnb_ue_ngap_id
  // present in ue_ref is not equal to  IE RAN_UE_NGAP_ID present
  // in ue_context_release_request
  EXPECT_EQ(ngap_amf_handle_message(state, peerInfo.assoc_id,
                                    peerInfo.instreams, &decoded_pdu),
            RETURNerror);
  ue_ref->amf_ue_ngap_id = amf_ue_ngap_id;
  ue_ref->gnb_ue_ngap_id = amf_ue_ngap_id;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_UEContextReleaseRequest_IEs, ie,
                                      container, Ngap_ProtocolIE_ID_id_Cause);
  // check if NGAP_Cause IE present
  ASSERT_TRUE(ie != NULL);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decoded_pdu);
  bdestroy(ue_context_release_req_msg);
}

// UE Context Release Command
TEST_F(NgapFlowTest, ue_context_release_command_rainy_day) {
  MessageDef* message_p = NULL;
  m5g_ue_description_t* ue_ref = NULL;

  // Verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  // Verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);

  message_p =
      itti_alloc_new_message(TASK_AMF_APP, NGAP_UE_CONTEXT_RELEASE_COMMAND);
  message_p->ittiMsgHeader.imsi = 0x311480000000001;

  // verify new ue is associated with gnb
  ue_ref = ngap_new_ue(state, peerInfo.assoc_id, gNB_UE_NGAP_ID);
  ASSERT_TRUE(ue_ref != NULL);
  ue_ref->ng_ue_state = NGAP_UE_CONNECTED;
  ue_ref->gnb_ue_ngap_id = gNB_UE_NGAP_ID;
  ue_ref->amf_ue_ngap_id = 0xffff;

  NGAP_UE_CONTEXT_RELEASE_COMMAND(message_p).amf_ue_ngap_id = AMF_UE_NGAP_ID;
  NGAP_UE_CONTEXT_RELEASE_COMMAND(message_p).gnb_ue_ngap_id = gNB_UE_NGAP_ID;
  NGAP_UE_CONTEXT_RELEASE_COMMAND(message_p).cause = NGAP_USER_INACTIVITY;

  // verify ue_context_release_command is not handled as amf_ue_ngap_id
  // present in ue_ref is not equla to IE AMF_UE_NGAP_ID present in
  // itti ue_context_release_command
  EXPECT_EQ(RETURNok, ngap_handle_ue_context_release_command(
                          state, &NGAP_UE_CONTEXT_RELEASE_COMMAND(message_p),
                          message_p->ittiMsgHeader.imsi));
  itti_free_msg_content(message_p);
  free(message_p);
}
// Pdu Session Resource Setup Response
TEST_F(NgapFlowTest, pdu_session_resource_setup_resp_rainy_day) {
  Ngap_PDUSessionResourceSetupResponse_t* container = NULL;
  gnb_description_t* gNB_ref = NULL;
  m5g_ue_description_t* ue_ref = NULL;
  Ngap_NGAP_PDU_t decoded_pdu = {};
  uint8_t pdu_ss_resource_setup_resp_hex_buff[] = {
      0x20, 0x1d, 0x00, 0x27, 0x00, 0x00, 0x03, 0x00, 0x0a, 0x40, 0x02,
      0x00, 0x05, 0x00, 0x55, 0x40, 0x04, 0x80, 0x01, 0x00, 0x01, 0x00,
      0x4b, 0x40, 0x12, 0x00, 0x00, 0x05, 0x0e, 0x00, 0x03, 0xe0, 0x05,
      0x05, 0x05, 0x02, 0x00, 0x00, 0x00, 0x0b, 0x01, 0x00, 0x00};
  uint16_t len = sizeof(pdu_ss_resource_setup_resp_hex_buff) / sizeof(uint8_t);

  // Verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  // Verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);

  bstring pdu_ss_resource_response_succ_msg =
      blk2bstr(pdu_ss_resource_setup_resp_hex_buff, len);

  // Verify pdu_session_resource_setup is decoded
  ASSERT_EQ(
      ngap_amf_decode_pdu(&decoded_pdu, pdu_ss_resource_response_succ_msg),
      RETURNok);
  container = &(decoded_pdu.choice.successfulOutcome.value.choice
                    .PDUSessionResourceSetupResponse);
  Ngap_PDUSessionResourceSetupResponseIEs_t* ie = NULL;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_PDUSessionResourceSetupResponseIEs_t,
                                      ie, container,
                                      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID);

  // Check if Ran_UE_NGAP_ID is present
  ASSERT_TRUE(ie != NULL);

  // Check if gNB exists
  gNB_ref = ngap_state_get_gnb(state, peerInfo.assoc_id);
  ASSERT_TRUE(gNB_ref != NULL);

  gnb_ue_ngap_id_t gnb_ue_ngap_id = 0;
  amf_ue_ngap_id_t amf_ue_ngap_id = 0;
  gnb_ue_ngap_id = (gnb_ue_ngap_id_t)(ie->value.choice.RAN_UE_NGAP_ID);
  ue_ref = ngap_new_ue(state, peerInfo.assoc_id, gnb_ue_ngap_id);

  // verify new ue is associated with gnb
  ASSERT_TRUE(ue_ref != NULL);
  ue_ref->gnb_ue_ngap_id = 0xffff;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_PDUSessionResourceSetupResponseIEs_t,
                                      ie, container,
                                      Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID);

  // Check if AMF_UE_NGAP_ID is present
  ASSERT_TRUE(ie != NULL);
  asn_INTEGER2ulong(&ie->value.choice.AMF_UE_NGAP_ID,
                    reinterpret_cast<uint64_t*>(&amf_ue_ngap_id));
  ue_ref->amf_ue_ngap_id = 0xffff;

  // verify pdu_session_resource_setup is not handled as amf_ue_ngap_id
  // present in ue_ref is not equal to ie AMF_UE_NGAP_ID
  // present in pdu_session_resource_setup message
  EXPECT_EQ(ngap_amf_handle_message(state, peerInfo.assoc_id,
                                    peerInfo.instreams, &decoded_pdu),
            RETURNerror);
  ue_ref->amf_ue_ngap_id = amf_ue_ngap_id;

  // verify pdu_session_resource_setup is not handled as ran_ue_ngap_id
  // present in ue_ref is not equal to ie RAN_UE_NGAP_ID
  // present in pdu_session_resource_setup message
  EXPECT_EQ(ngap_amf_handle_message(state, peerInfo.assoc_id,
                                    peerInfo.instreams, &decoded_pdu),
            RETURNerror);

  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_PDUSessionResourceSetupResponseIEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_PDUSessionResourceSetupListSURes, false);

  // verify Ngap_ProtocolIE_ID_id_PDUSessionResourceSetupListSURes ie is present
  ASSERT_TRUE(ie != NULL);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decoded_pdu);
  bdestroy(pdu_ss_resource_response_succ_msg);
}

// Unit for Initial UE message
TEST_F(NgapFlowTest, test_ue_notifications_from_amf) {
  unsigned char initial_ue_message_hexbuf[] = {
      0x00, 0x0f, 0x40, 0x48, 0x00, 0x00, 0x05, 0x00, 0x55, 0x00, 0x02,
      0x00, 0x01, 0x00, 0x26, 0x00, 0x1a, 0x19, 0x7e, 0x00, 0x41, 0x79,
      0x00, 0x0d, 0x01, 0x22, 0x62, 0x54, 0x00, 0x00, 0x00, 0x00, 0x00,
      0x00, 0x00, 0x00, 0x01, 0x2e, 0x04, 0xf0, 0xf0, 0xf0, 0xf0, 0x00,
      0x79, 0x00, 0x13, 0x50, 0x22, 0x42, 0x65, 0x00, 0x00, 0x00, 0x01,
      0x00, 0x22, 0x42, 0x65, 0x00, 0x00, 0x01, 0xe4, 0xf7, 0x04, 0x44,
      0x00, 0x5a, 0x40, 0x01, 0x18, 0x00, 0x70, 0x40, 0x01, 0x00};

  // Verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  // Verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);

  Ngap_NGAP_PDU_t decoded_pdu = {};
  uint16_t length = sizeof(initial_ue_message_hexbuf) / sizeof(unsigned char);

  bstring ngap_initial_ue_msg = blk2bstr(initial_ue_message_hexbuf, length);

  // Check if the pdu can be decoded
  ASSERT_EQ(ngap_amf_decode_pdu(&decoded_pdu, ngap_initial_ue_msg), RETURNok);

  // check if initial UE message is handled successfully
  EXPECT_EQ(ngap_amf_handle_message(state, peerInfo.assoc_id,
                                    peerInfo.instreams, &decoded_pdu),
            RETURNok);

  Ngap_InitialUEMessage_t* container;
  gnb_description_t* gNB_ref = NULL;
  m5g_ue_description_t* ue_ref = NULL;

  container =
      &(decoded_pdu.choice.initiatingMessage.value.choice.InitialUEMessage);
  Ngap_InitialUEMessage_IEs_t* ie = NULL;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_InitialUEMessage_IEs_t, ie,
                                      container,
                                      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID);

  // Check if Ran_UE_NGAP_ID is present in initial message
  ASSERT_TRUE(ie != NULL);

  // Check if gNB exists
  gNB_ref = ngap_state_get_gnb(state, peerInfo.assoc_id);
  ASSERT_TRUE(gNB_ref != NULL);

  gnb_ue_ngap_id_t gnb_ue_ngap_id = 0;
  gnb_ue_ngap_id = (gnb_ue_ngap_id_t)(ie->value.choice.RAN_UE_NGAP_ID);

  // Check for UE associated with gNB
  ue_ref = ngap_state_get_ue_gnbid(gNB_ref->sctp_assoc_id, gnb_ue_ngap_id);
  ASSERT_TRUE(ue_ref != NULL);

  // Check if UE is pointing to invalid ID if it is initial ue message
  EXPECT_EQ(ue_ref->amf_ue_ngap_id, INVALID_AMF_UE_NGAP_ID);

  itti_amf_app_ngap_amf_ue_id_notification_t notification_p;
  memset(&notification_p, 0,
         sizeof(itti_amf_app_ngap_amf_ue_id_notification_t));
  notification_p.gnb_ue_ngap_id = gnb_ue_ngap_id;
  notification_p.amf_ue_ngap_id = 10;
  notification_p.sctp_assoc_id = gNB_ref->sctp_assoc_id;

  ngap_handle_amf_ue_id_notification(state, &notification_p);
  ue_ref = ngap_state_get_ue_gnbid(gNB_ref->sctp_assoc_id, gnb_ue_ngap_id);
  ASSERT_TRUE(ue_ref != NULL);

  // Check if UE is pointing to invalid ID if it is initial ue message
  EXPECT_EQ(ue_ref->amf_ue_ngap_id, 10);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decoded_pdu);
  bdestroy(ngap_initial_ue_msg);
}

TEST_F(NgapFlowTest, NgapHandleSctpDisconnection) {
  MessageDef* amf_message_p = NULL;
  MessageDef* sctp_message_p = NULL;
  m5g_ue_description_t* ue_ref = NULL;
  Ngap_InitialUEMessage_t* container;

  unsigned char initial_ue_message_hexbuf[] = {
      0x00, 0x0f, 0x40, 0x48, 0x00, 0x00, 0x05, 0x00, 0x55, 0x00, 0x02,
      0x00, 0x01, 0x00, 0x26, 0x00, 0x1a, 0x19, 0x7e, 0x00, 0x41, 0x79,
      0x00, 0x0d, 0x01, 0x22, 0x62, 0x54, 0x00, 0x00, 0x00, 0x00, 0x00,
      0x00, 0x00, 0x00, 0x01, 0x2e, 0x04, 0xf0, 0xf0, 0xf0, 0xf0, 0x00,
      0x79, 0x00, 0x13, 0x50, 0x22, 0x42, 0x65, 0x00, 0x00, 0x00, 0x01,
      0x00, 0x22, 0x42, 0x65, 0x00, 0x00, 0x01, 0xe4, 0xf7, 0x04, 0x44,
      0x00, 0x5a, 0x40, 0x01, 0x18, 0x00, 0x70, 0x40, 0x01, 0x00};

  std::vector<MessagesIds> expected_Ids{
      NGAP_INITIAL_UE_MESSAGE, NGAP_GNB_DEREGISTERED_IND, SCTP_DATA_REQ,
      NGAP_UE_CONTEXT_RELEASE_COMPLETE};

  // Verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  // Verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);

  gnb_description_t* gnb_association = NULL;
  gnb_association = ngap_state_get_gnb(state, peerInfo.assoc_id);
  ASSERT_TRUE(gnb_association != NULL);
  EXPECT_EQ(gnb_association->ng_state, NGAP_INIT);
  EXPECT_EQ(gnb_association->nb_ue_associated, 0);

  // Handling initial UE message
  Ngap_NGAP_PDU_t decoded_pdu = {};
  uint16_t length = sizeof(initial_ue_message_hexbuf) / sizeof(unsigned char);
  bstring ngap_initial_ue_msg = blk2bstr(initial_ue_message_hexbuf, length);

  // Check if the pdu can be decoded
  ASSERT_EQ(ngap_amf_decode_pdu(&decoded_pdu, ngap_initial_ue_msg), RETURNok);

  // Mocking the sending of Initial UE message to AMF
  EXPECT_EQ(ngap_amf_handle_message(state, peerInfo.assoc_id,
                                    peerInfo.instreams, &decoded_pdu),
            RETURNok);
  // Checking that UE is connected after Initial UE message
  EXPECT_EQ(gnb_association->nb_ue_associated, 1);

  container =
      &(decoded_pdu.choice.initiatingMessage.value.choice.InitialUEMessage);
  Ngap_InitialUEMessage_IEs_t* ie = NULL;
  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_InitialUEMessage_IEs_t, ie,
                                      container,
                                      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID);

  // Check if Ran_UE_NGAP_ID is present in initial message
  ASSERT_TRUE(ie != NULL);
  gnb_ue_ngap_id_t gnb_ue_ngap_id = 0;
  gnb_ue_ngap_id = (gnb_ue_ngap_id_t)(ie->value.choice.RAN_UE_NGAP_ID);

  // Mocking the AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION from AMF
  itti_amf_app_ngap_amf_ue_id_notification_t notification_p = {};
  notification_p.gnb_ue_ngap_id = gnb_ue_ngap_id;
  notification_p.amf_ue_ngap_id = 1;
  notification_p.sctp_assoc_id = gnb_association->sctp_assoc_id;

  ngap_handle_amf_ue_id_notification(state, &notification_p);

  // Mocking the SCTP_CLOSE_ASSOCIATION ITTI message
  sctp_message_p = itti_alloc_new_message(TASK_SCTP, SCTP_CLOSE_ASSOCIATION);

  SCTP_CLOSE_ASSOCIATION(sctp_message_p).assoc_id = peerInfo.assoc_id;
  // To test SCTP Shutdown for NGAP
  SCTP_CLOSE_ASSOCIATION(sctp_message_p).reset = false;

  EXPECT_EQ(ngap_handle_sctp_disconnection(
                state, SCTP_CLOSE_ASSOCIATION(sctp_message_p).assoc_id,
                SCTP_CLOSE_ASSOCIATION(sctp_message_p).reset),
            RETURNok);

  gnb_description_t* gnb_association_sd = NULL;
  gnb_association_sd = ngap_state_get_gnb(state, peerInfo.assoc_id);
  ASSERT_TRUE(gnb_association_sd != NULL);
  EXPECT_EQ(gnb_association_sd->ng_state, NGAP_SHUTDOWN);
  EXPECT_EQ(gnb_association_sd->nb_ue_associated, 1);

  // Mocking the UE_CONTEXT_RELEASE_COMMAND from AMF
  amf_message_p =
      itti_alloc_new_message(TASK_AMF_APP, NGAP_UE_CONTEXT_RELEASE_COMMAND);
  ASSERT_TRUE(amf_message_p != NULL);

  // Check if UE is associated with gnb
  ue_ref =
      ngap_state_get_ue_gnbid(gnb_association->sctp_assoc_id, gnb_ue_ngap_id);
  ASSERT_TRUE(ue_ref != NULL);

  amf_message_p->ittiMsgHeader.imsi = 0x311480000000001;
  NGAP_UE_CONTEXT_RELEASE_COMMAND(amf_message_p).amf_ue_ngap_id =
      ue_ref->amf_ue_ngap_id;
  NGAP_UE_CONTEXT_RELEASE_COMMAND(amf_message_p).gnb_ue_ngap_id =
      ue_ref->gnb_ue_ngap_id;
  NGAP_UE_CONTEXT_RELEASE_COMMAND(amf_message_p).cause = NGAP_USER_INACTIVITY;

  EXPECT_EQ(state->num_gnbs, 1);

  // verify ue_context_release_command is encoded correctly
  EXPECT_EQ(ngap_handle_ue_context_release_command(
                state, &NGAP_UE_CONTEXT_RELEASE_COMMAND(amf_message_p),
                amf_message_p->ittiMsgHeader.imsi),
            RETURNok);

  // Checking the number of gnbs after ue context release in NGAP.
  EXPECT_EQ(state->num_gnbs, 0);

  // Checking that gnb description should be null after ue context release.
  gnb_description_t* gnb_association_ue = NULL;
  gnb_association_ue = ngap_state_get_gnb(state, peerInfo.assoc_id);
  ASSERT_TRUE(gnb_association_ue == NULL);

  EXPECT_TRUE(expected_Ids == NGAPClientServicer::getInstance().msgtype_stack);

  itti_free_msg_content(sctp_message_p);
  free(sctp_message_p);
  itti_free_msg_content(amf_message_p);
  free(amf_message_p);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decoded_pdu);
  bdestroy(ngap_initial_ue_msg);
}

// Pdu Session Resource Modify Request
TEST_F(NgapFlowTest, pdu_sess_resource_modify_req_sunny_day) {
  MessageDef* message_p = NULL;
  uint8_t ip_buff[4] = {0xc0, 0xa8, 0x3c, 0x8e};
  m5g_ue_description_t* ue_ref = NULL;
  itti_ngap_pdu_session_resource_modify_request_t* ngap_pdu_ses_modify_req =
      nullptr;
  pdu_session_resource_modify_request_transfer_t
      amf_pdu_ses_modify_transfer_req = {};

  // Verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  // Verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);

  message_p = itti_alloc_new_message(TASK_AMF_APP,
                                     NGAP_PDU_SESSION_RESOURCE_MODIFY_REQ);
  message_p->ittiMsgHeader.imsi = 0x311480000000001;

  ngap_pdu_ses_modify_req =
      &message_p->ittiMsg.ngap_pdu_session_resource_modify_req;
  memset(ngap_pdu_ses_modify_req, 0,
         sizeof(itti_ngap_pdu_session_resource_modify_request_t));

  ngap_pdu_ses_modify_req->gnb_ue_ngap_id = gNB_UE_NGAP_ID;
  ngap_pdu_ses_modify_req->amf_ue_ngap_id = AMF_UE_NGAP_ID;

  ngap_pdu_ses_modify_req->pduSessResourceModReqList.no_of_items = 1;
  ngap_pdu_ses_modify_req->pduSessResourceModReqList.item[0].Pdu_Session_ID =
      0x05;

  amf_pdu_ses_modify_transfer_req.pdu_sess_aggregate_max_bit_rate.dl = 1024;
  amf_pdu_ses_modify_transfer_req.pdu_sess_aggregate_max_bit_rate.ul = 1024;
  amf_pdu_ses_modify_transfer_req.ul_ng_u_up_tnl_modify_list.numOfItems = 1;
  up_transport_layer_information_t* ul_ng_u_up_tnl_item =
      &amf_pdu_ses_modify_transfer_req.ul_ng_u_up_tnl_modify_list
           .ul_ng_u_up_tnl_modfy_item[0];

  ul_ng_u_up_tnl_item->gtp_tnl.gtp_tied[0] = 0x80;
  ul_ng_u_up_tnl_item->gtp_tnl.gtp_tied[1] = 0x00;
  ul_ng_u_up_tnl_item->gtp_tnl.gtp_tied[2] = 0x00;
  ul_ng_u_up_tnl_item->gtp_tnl.gtp_tied[3] = 0x01;
  ul_ng_u_up_tnl_item->gtp_tnl.endpoint_ip_address = blk2bstr(ip_buff, 4);

  amf_pdu_ses_modify_transfer_req.qos_flow_add_or_mod_request_list
      .maxNumOfQosFlows = 1;
  qos_flow_request_list_t* qos_flow_request =
      &amf_pdu_ses_modify_transfer_req.qos_flow_add_or_mod_request_list.item[0];
  qos_flow_request->qos_flow_req_item.qos_flow_identifier = 5;
  qos_flow_request->qos_flow_req_item.qos_flow_level_qos_param
      .qos_characteristic.non_dynamic_5QI_desc.fiveQI = 9;
  qos_flow_request->qos_flow_req_item.qos_flow_level_qos_param
      .alloc_reten_priority.priority_level = 1;
  qos_flow_request->qos_flow_req_item.qos_flow_level_qos_param
      .alloc_reten_priority.pre_emption_cap = SHALL_NOT_TRIGGER_PRE_EMPTION;
  qos_flow_request->qos_flow_req_item.qos_flow_level_qos_param
      .alloc_reten_priority.pre_emption_vul = NOT_PREEMPTABLE;
  ngap_pdu_ses_modify_req->pduSessResourceModReqList.item[0]
      .PDU_Session_Resource_Modify_Request_Transfer =
      amf_pdu_ses_modify_transfer_req;

  // verify new ue is associated with gnb
  ue_ref = ngap_new_ue(state, peerInfo.assoc_id, gNB_UE_NGAP_ID);
  ASSERT_TRUE(ue_ref != NULL);
  ue_ref->ng_ue_state = NGAP_UE_CONNECTED;
  ue_ref->gnb_ue_ngap_id = gNB_UE_NGAP_ID;
  ue_ref->amf_ue_ngap_id = AMF_UE_NGAP_ID;

  // verify pdu_session_resource_setup_request is encoded correctly
  EXPECT_EQ(RETURNok,
            ngap_generate_ngap_pdusession_resource_modify_req(
                state, &NGAP_PDU_SESSION_RESOURCE_MODIFY_REQ(message_p)));
  itti_free_msg_content(message_p);
  free(message_p);
}

// Pdu Session Resource Modify Response
TEST_F(NgapFlowTest, pdu_session_resource_modify_resp_sunny_day) {
  Ngap_PDUSessionResourceModifyResponse_t* container = NULL;
  gnb_description_t* gNB_ref = NULL;
  m5g_ue_description_t* ue_ref = NULL;
  Ngap_NGAP_PDU_t decoded_pdu = {};
  uint8_t pdu_ss_resource_modify_resp_hex_buff[] = {
      0x20, 0x1a, 0x00, 0x1e, 0x00, 0x00, 0x03, 0x00, 0x0a, 0x40, 0x04, 0x40,
      0x01, 0x00, 0x01, 0x00, 0x55, 0x40, 0x04, 0x80, 0x01, 0x00, 0x01, 0x00,
      0x41, 0x40, 0x07, 0x00, 0x00, 0x05, 0x03, 0x10, 0x00, 0x0c};
  uint16_t len = sizeof(pdu_ss_resource_modify_resp_hex_buff) / sizeof(uint8_t);

  // Verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);
  // Verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);

  bstring pdu_ss_modify_response_succ_msg =
      blk2bstr(pdu_ss_resource_modify_resp_hex_buff, len);

  // Check if the pdu_session_resource_modify_response decoded successfully
  ASSERT_EQ(ngap_amf_decode_pdu(&decoded_pdu, pdu_ss_modify_response_succ_msg),
            RETURNok);
  container = &(decoded_pdu.choice.successfulOutcome.value.choice
                    .PDUSessionResourceModifyResponse);
  Ngap_PDUSessionResourceModifyResponseIEs_t* ie = NULL;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(
      Ngap_PDUSessionResourceModifyResponseIEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID);

  // Check if RAN_UE_NGAP_ID Present
  ASSERT_TRUE(ie != NULL);

  // Check if gNB exists
  gNB_ref = ngap_state_get_gnb(state, peerInfo.assoc_id);
  ASSERT_TRUE(gNB_ref != NULL);

  gnb_ue_ngap_id_t gnb_ue_ngap_id = 0;
  amf_ue_ngap_id_t amf_ue_ngap_id = 0;
  gnb_ue_ngap_id = (gnb_ue_ngap_id_t)(ie->value.choice.RAN_UE_NGAP_ID);

  // verify new ue is associated with gnb
  ue_ref = ngap_new_ue(state, peerInfo.assoc_id, gnb_ue_ngap_id);
  ASSERT_TRUE(ue_ref != NULL);
  ue_ref->gnb_ue_ngap_id = gnb_ue_ngap_id;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(
      Ngap_PDUSessionResourceModifyResponseIEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID);

  // Check if AMF_UE_NGAP_ID present
  ASSERT_TRUE(ie != NULL);
  asn_INTEGER2ulong(&ie->value.choice.AMF_UE_NGAP_ID,
                    reinterpret_cast<uint64_t*>(&amf_ue_ngap_id));
  ue_ref->amf_ue_ngap_id = amf_ue_ngap_id;

  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_PDUSessionResourceModifyResponseIEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_PDUSessionResourceModifyListModRes, false);
  // verify if Ngap_ProtocolIE_ID_id_PDUSessionResourceModifyListModRes present
  ASSERT_TRUE(ie != NULL);

  // verify pdu_session_resource_modify_response is handled_correctly
  EXPECT_EQ(ngap_amf_handle_message(state, peerInfo.assoc_id,
                                    peerInfo.instreams, &decoded_pdu),
            RETURNok);

  // Check if UE is not invalid ID
  EXPECT_NE(ue_ref->amf_ue_ngap_id, INVALID_AMF_UE_NGAP_ID);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decoded_pdu);
  bdestroy(pdu_ss_modify_response_succ_msg);
}

TEST_F(NgapFlowTest, TestNgapReset) {
  m5g_ue_description_t* ue_ref = NULL;
  Ngap_InitialUEMessage_t* container;
  unsigned char ngap_setup_req_hexbuf[] = {
      0x00, 0x15, 0x00, 0x42, 0x00, 0x00, 0x04, 0x00, 0x1b, 0x00, 0x09, 0x00,
      0x22, 0x42, 0x65, 0x50, 0x00, 0x00, 0x00, 0x01, 0x00, 0x52, 0x40, 0x18,
      0x0a, 0x80, 0x55, 0x45, 0x52, 0x41, 0x4e, 0x53, 0x49, 0x4d, 0x2d, 0x67,
      0x6e, 0x62, 0x2d, 0x32, 0x32, 0x32, 0x2d, 0x34, 0x35, 0x36, 0x2d, 0x31,
      0x00, 0x66, 0x00, 0x0d, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x22, 0x42,
      0x65, 0x00, 0x00, 0x00, 0x08, 0x00, 0x15, 0x40, 0x01, 0x40};

  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);

  Ngap_NGAP_PDU_t decoded_pdu1 = {};
  uint16_t length1 = sizeof(ngap_setup_req_hexbuf) / sizeof(unsigned char);

  bstring ngap_setup_req_msg = blk2bstr(ngap_setup_req_hexbuf, length1);

  // Check if the pdu can be decoded
  ASSERT_EQ(ngap_amf_decode_pdu(&decoded_pdu1, ngap_setup_req_msg), RETURNok);
  bdestroy(ngap_setup_req_msg);

  sctp_stream_id_t stream_id = 0;

  gnb_description_t* gnb_association = NULL;
  gnb_association = ngap_state_get_gnb(state, peerInfo.assoc_id);

  EXPECT_EQ(ngap_amf_handle_message(state, peerInfo.assoc_id, stream_id,
                                    &decoded_pdu1),
            RETURNok);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decoded_pdu1);

  unsigned char initial_ue_message_hexbuf[] = {
      0x00, 0x0f, 0x40, 0x48, 0x00, 0x00, 0x05, 0x00, 0x55, 0x00, 0x02,
      0x00, 0x01, 0x00, 0x26, 0x00, 0x1a, 0x19, 0x7e, 0x00, 0x41, 0x79,
      0x00, 0x0d, 0x01, 0x22, 0x62, 0x54, 0x00, 0x00, 0x00, 0x00, 0x00,
      0x00, 0x00, 0x00, 0x01, 0x2e, 0x04, 0xf0, 0xf0, 0xf0, 0xf0, 0x00,
      0x79, 0x00, 0x13, 0x50, 0x22, 0x42, 0x65, 0x00, 0x00, 0x00, 0x01,
      0x00, 0x22, 0x42, 0x65, 0x00, 0x00, 0x01, 0xe4, 0xf7, 0x04, 0x44,
      0x00, 0x5a, 0x40, 0x01, 0x18, 0x00, 0x70, 0x40, 0x01, 0x00};

  EXPECT_EQ(state->gnbs.num_elements, 1);

  Ngap_NGAP_PDU_t decoded_pdu = {};
  uint16_t length = sizeof(initial_ue_message_hexbuf) / sizeof(unsigned char);
  bstring ngap_initial_ue_msg = blk2bstr(initial_ue_message_hexbuf, length);

  // Check if the pdu can be decoded
  ASSERT_EQ(ngap_amf_decode_pdu(&decoded_pdu, ngap_initial_ue_msg), RETURNok);
  bdestroy(ngap_initial_ue_msg);

  // check if initial UE message is handled successfully
  EXPECT_EQ(ngap_amf_handle_message(state, peerInfo.assoc_id,
                                    peerInfo.instreams, &decoded_pdu),
            RETURNok);

  container =
      &(decoded_pdu.choice.initiatingMessage.value.choice.InitialUEMessage);
  Ngap_InitialUEMessage_IEs_t* ie = NULL;
  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_InitialUEMessage_IEs_t, ie,
                                      container,
                                      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID);

  // Check if Ran_UE_NGAP_ID is present in initial message
  ASSERT_TRUE(ie != NULL);
  gnb_ue_ngap_id_t gnb_ue_ngap_id = 0;
  gnb_ue_ngap_id = (gnb_ue_ngap_id_t)(ie->value.choice.RAN_UE_NGAP_ID);

  // Mocking the AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION from AMF
  itti_amf_app_ngap_amf_ue_id_notification_t notification_p = {};
  notification_p.gnb_ue_ngap_id = gnb_ue_ngap_id;
  notification_p.amf_ue_ngap_id = 1;
  notification_p.sctp_assoc_id = gnb_association->sctp_assoc_id;

  ngap_handle_amf_ue_id_notification(state, &notification_p);

  ue_ref =
      ngap_state_get_ue_gnbid(gnb_association->sctp_assoc_id, gnb_ue_ngap_id);
  ASSERT_TRUE(ue_ref != NULL);

  unsigned char ng_reset_hexbuf[] = {0x00, 0x14, 0x00, 0x0e, 0x00, 0x00,
                                     0x02, 0x00, 0x0f, 0x40, 0x02, 0x00,
                                     0x00, 0x00, 0x58, 0x00, 0x01, 0x00};
  Ngap_NGAP_PDU_t pdu = {};
  uint16_t length_reset = sizeof(ng_reset_hexbuf) / sizeof(unsigned char);
  bstring ng_reset_msg = blk2bstr(ng_reset_hexbuf, length_reset);

  //   Check if the pdu can be decoded
  ASSERT_EQ(ngap_amf_decode_pdu(&pdu, ng_reset_msg), RETURNok);
  bdestroy(ng_reset_msg);

  int rc = ngap_amf_handle_message(state, peerInfo.assoc_id, peerInfo.instreams,
                                   &pdu);
  EXPECT_EQ(rc, RETURNok);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decoded_pdu);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &pdu);
}

TEST_F(NgapFlowTest, test_gNB_reset_ack) {
  unsigned char initial_ue_message_hexbuf[] = {
      0x00, 0x0f, 0x40, 0x48, 0x00, 0x00, 0x05, 0x00, 0x55, 0x00, 0x02,
      0x00, 0x01, 0x00, 0x26, 0x00, 0x1a, 0x19, 0x7e, 0x00, 0x41, 0x79,
      0x00, 0x0d, 0x01, 0x22, 0x62, 0x54, 0x00, 0x00, 0x00, 0x00, 0x00,
      0x00, 0x00, 0x00, 0x01, 0x2e, 0x04, 0xf0, 0xf0, 0xf0, 0xf0, 0x00,
      0x79, 0x00, 0x13, 0x50, 0x22, 0x42, 0x65, 0x00, 0x00, 0x00, 0x01,
      0x00, 0x22, 0x42, 0x65, 0x00, 0x00, 0x01, 0xe4, 0xf7, 0x04, 0x44,
      0x00, 0x5a, 0x40, 0x01, 0x18, 0x00, 0x70, 0x40, 0x01, 0x00};

  // Verify sctp association is successful
  EXPECT_EQ(ngap_handle_new_association(state, &peerInfo), RETURNok);

  // Verify number of connected gNB's is 1
  EXPECT_EQ(state->gnbs.num_elements, 1);

  Ngap_NGAP_PDU_t decoded_pdu = {};
  uint16_t length = sizeof(initial_ue_message_hexbuf) / sizeof(unsigned char);
  bstring ngap_initial_ue_msg = blk2bstr(initial_ue_message_hexbuf, length);

  // Check if the pdu can be decoded
  ASSERT_EQ(ngap_amf_decode_pdu(&decoded_pdu, ngap_initial_ue_msg), RETURNok);

  // check if initial UE message is handled successfully
  EXPECT_EQ(ngap_amf_handle_message(state, peerInfo.assoc_id,
                                    peerInfo.instreams, &decoded_pdu),
            RETURNok);

  // Send NGAP GNB Reset Acknowledgement
  EXPECT_EQ(send_ngap_gnb_reset_ack(), RETURNok);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decoded_pdu);
  bdestroy(ngap_initial_ue_msg);
}
}  // namespace magma5g
