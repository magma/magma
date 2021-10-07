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

#include "mock_utils.h"
#include "util_ngap_pkt.h"
#include <gtest/gtest.h>

#include "util_ngap_pkt.h"
extern "C" {
#include "log.h"
#include "ngap_amf_handlers.h"
#include "amf_config.h"
}
#include "ngap_state_manager.h"

using ::testing::Test;

namespace magma5g {

TEST(test_ngap_flow_handler, initial_ue_message) {
  unsigned char initial_ue_message_hexbuf[] = {
      0x00, 0x0f, 0x40, 0x48, 0x00, 0x00, 0x05, 0x00, 0x55, 0x00, 0x02,
      0x00, 0x01, 0x00, 0x26, 0x00, 0x1a, 0x19, 0x7e, 0x00, 0x41, 0x79,
      0x00, 0x0d, 0x01, 0x22, 0x62, 0x54, 0x00, 0x00, 0x00, 0x00, 0x00,
      0x00, 0x00, 0x00, 0x01, 0x2e, 0x04, 0xf0, 0xf0, 0xf0, 0xf0, 0x00,
      0x79, 0x00, 0x13, 0x48, 0x22, 0x42, 0x65, 0x00, 0x00, 0x00, 0x01,
      0x00, 0x22, 0x42, 0x65, 0x00, 0x00, 0x01, 0xe4, 0xf7, 0x04, 0x44,
      0x00, 0x5a, 0x40, 0x01, 0x18, 0x00, 0x70, 0x40, 0x01, 0x00};

  itti_init(
      TASK_MAX, THREAD_MAX, MESSAGES_ID_MAX, tasks_info, messages_info, NULL,
      NULL);

  amf_config_init(&amf_config);
  ngap_state_init(2, 2, false);
  ngap_state_t* state = NULL;
  state               = get_ngap_state(false);

  bstring ran_cp_ipaddr = bfromcstr("\xc0\xa8\x3c\x8d");
  sctp_new_peer_t p     = {
      .instreams     = 1,
      .outstreams    = 2,
      .assoc_id      = 3,
      .ran_cp_ipaddr = ran_cp_ipaddr,
  };

  EXPECT_EQ(ngap_handle_new_association(state, &p), RETURNok);
  EXPECT_EQ(state->gnbs.num_elements, 1);

  Ngap_NGAP_PDU_t decoded_pdu = {};
  unsigned short length =
      sizeof(initial_ue_message_hexbuf) / sizeof(unsigned char);

  bstring ngap_initial_ue_msg = blk2bstr(initial_ue_message_hexbuf, length);
  memcpy(ngap_initial_ue_msg->data, initial_ue_message_hexbuf, length);

  // Check if the pdu can be decoded
  ASSERT_EQ(ngap_amf_decode_pdu(&decoded_pdu, ngap_initial_ue_msg), RETURNok);

  // Walk through the initial UE message
  EXPECT_EQ(
      ngap_amf_handle_message(state, p.assoc_id, p.instreams, &decoded_pdu),
      RETURNok);

  Ngap_InitialUEMessage_t* container;
  gnb_description_t* gNB_ref   = NULL;
  m5g_ue_description_t* ue_ref = NULL;

  container =
      &(decoded_pdu.choice.initiatingMessage.value.choice.InitialUEMessage);
  Ngap_InitialUEMessage_IEs_t* ie = NULL;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(
      Ngap_InitialUEMessage_IEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID);

  // Check if ran UE ID was found
  ASSERT_TRUE(ie != NULL);

  // Check if GNB exists
  gNB_ref = ngap_state_get_gnb(state, p.assoc_id);
  ASSERT_TRUE(gNB_ref != NULL);

  gnb_ue_ngap_id_t gnb_ue_ngap_id = 0;
  gnb_ue_ngap_id = (gnb_ue_ngap_id_t)(ie->value.choice.RAN_UE_NGAP_ID);

  // Check for UE reference
  ue_ref = ngap_state_get_ue_gnbid(gNB_ref->sctp_assoc_id, gnb_ue_ngap_id);
  ASSERT_TRUE(ue_ref != NULL);

  // Check if UE is pointing to invalid ID
  EXPECT_EQ(ue_ref->amf_ue_ngap_id, INVALID_AMF_UE_NGAP_ID);

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, &decoded_pdu);
  bdestroy(ngap_initial_ue_msg);
  bdestroy(ran_cp_ipaddr);
}

}  // namespace magma5g
