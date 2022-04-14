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
#include <gtest/gtest.h>
#include <glog/logging.h>
#include "endian.h"

extern "C" {
#include "lte/gateway/c/core/oai/lib/gtpv2-c/nwgtpv2c-0.11/include/NwGtpv2c.h"
#include "lte/gateway/c/core/oai/lib/gtpv2-c/nwgtpv2c-0.11/shared/NwGtpv2cMsg.h"
#include "lte/gateway/c/core/oai/lib/gtpv2-c/nwgtpv2c-0.11/include/NwGtpv2cPrivate.h"
#include "lte/gateway/c/core/oai/lib/gtpv2-c/gtpv2c_ie_formatter/shared/gtpv2c_ie_formatter.h"
}
using ::testing::Test;

namespace magma {
namespace lte {

TEST(test_create_session_request_pdn1, create_session_request_pdn1) {
  nw_gtpv2c_stack_handle_t s11_mme_stack_handle = 0;

  uint8_t shift_buffer[1024];
  uint8_t create_session_request_pdn1[] = {
      0x48, 0x20, 0x00, 0xa3, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46, 0x80, 0x00,
      0x03, 0x00, 0x01, 0x00, 0x00, 0x01, 0x00, 0x08, 0x00, 0x13, 0x41, 0x08,
      0x01, 0x00, 0x10, 0x01, 0xf1, 0x56, 0x00, 0x0d, 0x00, 0x18, 0x02, 0xf8,
      0x99, 0x00, 0x01, 0x02, 0xf8, 0x99, 0x0f, 0xff, 0xff, 0xd2, 0x52, 0x00,
      0x01, 0x00, 0x06, 0x63, 0x00, 0x01, 0x00, 0x01, 0x4f, 0x00, 0x05, 0x00,
      0x01, 0x00, 0x00, 0x00, 0x00, 0x7f, 0x00, 0x01, 0x00, 0x00, 0x48, 0x00,
      0x08, 0x00, 0x00, 0x00, 0x27, 0x10, 0x00, 0x00, 0x4e, 0x20, 0x4d, 0x00,
      0x03, 0x00, 0x00, 0x00, 0x00, 0x57, 0x00, 0x09, 0x00, 0x8a, 0x00, 0x00,
      0x00, 0x01, 0xc0, 0xa8, 0x41, 0x04, 0x47, 0x00, 0x09, 0x00, 0x03, 0x6f,
      0x61, 0x69, 0x04, 0x69, 0x70, 0x76, 0x34, 0x80, 0x00, 0x01, 0x00, 0x00,
      0x53, 0x00, 0x03, 0x00, 0x02, 0xf8, 0x99, 0x5d, 0x00, 0x24, 0x00, 0x49,
      0x00, 0x01, 0x00, 0x05, 0x50, 0x00, 0x16, 0x00, 0x3c, 0x09, 0x00, 0x00,
      0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
      0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x54, 0x00, 0x01, 0x00, 0x00};
  EXPECT_EQ(NW_OK, nwGtpv2cInitialize(&s11_mme_stack_handle));
  nw_gtpv2c_stack_t* pGtpv2c_stack =
      reinterpret_cast<nw_gtpv2c_stack_t*>(s11_mme_stack_handle);
  {
    EXPECT_TRUE(sizeof(shift_buffer) >=
                (128 + sizeof(create_session_request_pdn1)));
    // whatever is the byte alignment, should not fail
    for (int s = 0; s < 128; s++) {
      memset(shift_buffer, 0, sizeof(shift_buffer));
      memcpy(&shift_buffer[s], create_session_request_pdn1,
             sizeof(create_session_request_pdn1));
      nw_gtpv2c_msg_handle_t hMsg = 0;
      nw_gtpv2c_error_t error = {0};
      EXPECT_EQ(NW_OK, nwGtpv2cMsgFromBufferNew(
                           s11_mme_stack_handle, &shift_buffer[s],
                           sizeof(create_session_request_pdn1), &hMsg));

      nw_gtpv2c_msg_t* pMsg = reinterpret_cast<nw_gtpv2c_msg_t*>(hMsg);
      EXPECT_EQ(pMsg->msgLen, 167);
      EXPECT_EQ(pMsg->version, 2);
      EXPECT_TRUE(pMsg->teidPresent);
      EXPECT_EQ(pMsg->msgType, 32);
      EXPECT_EQ(pMsg->seqNum, 0x4680);
      NW_ASSERT(pGtpv2c_stack->pGtpv2cMsgIeParseInfo[pMsg->msgType]);
      EXPECT_EQ(NW_OK, nwGtpv2cMsgIeParse(
                           pGtpv2c_stack->pGtpv2cMsgIeParseInfo[pMsg->msgType],
                           hMsg, &error));
      EXPECT_TRUE(
          pMsg->isIeValid[NW_GTPV2C_IE_RECOVERY][NW_GTPV2C_IE_INSTANCE_ZERO]);
      EXPECT_TRUE(
          pMsg->isIeValid[NW_GTPV2C_IE_ULI][NW_GTPV2C_IE_INSTANCE_ZERO]);
      EXPECT_TRUE(
          pMsg->isIeValid[NW_GTPV2C_IE_RAT_TYPE][NW_GTPV2C_IE_INSTANCE_ZERO]);
      EXPECT_TRUE(
          pMsg->isIeValid[NW_GTPV2C_IE_PDN_TYPE][NW_GTPV2C_IE_INSTANCE_ZERO]);
      EXPECT_TRUE(
          pMsg->isIeValid[NW_GTPV2C_IE_PAA][NW_GTPV2C_IE_INSTANCE_ZERO]);
      EXPECT_TRUE(pMsg->isIeValid[NW_GTPV2C_IE_APN_RESTRICTION]
                                 [NW_GTPV2C_IE_INSTANCE_ZERO]);
      EXPECT_TRUE(
          pMsg->isIeValid[NW_GTPV2C_IE_AMBR][NW_GTPV2C_IE_INSTANCE_ZERO]);
      EXPECT_TRUE(
          pMsg->isIeValid[NW_GTPV2C_IE_INDICATION][NW_GTPV2C_IE_INSTANCE_ZERO]);
      EXPECT_TRUE(
          pMsg->isIeValid[NW_GTPV2C_IE_FTEID][NW_GTPV2C_IE_INSTANCE_ZERO]);
      fteid_t fteid;
      memset(reinterpret_cast<void*>(&fteid), 0xAB, sizeof(fteid));
      nw_gtpv2c_ie_tlv_t* pIe = reinterpret_cast<nw_gtpv2c_ie_tlv_t*>(
          pMsg->pIe[NW_GTPV2C_IE_FTEID][NW_GTPV2C_IE_INSTANCE_ZERO]);
      uint16_t ieLength = ntohs(pIe->l);
      EXPECT_EQ(
          NW_OK,
          gtpv2c_fteid_ie_get(
              NW_GTPV2C_IE_FTEID, ieLength, NW_GTPV2C_IE_INSTANCE_ZERO,
              &pMsg->pIe[NW_GTPV2C_IE_FTEID][NW_GTPV2C_IE_INSTANCE_ZERO][4],
              reinterpret_cast<void*>(&fteid)));
      EXPECT_EQ(fteid.interface_type, S11_MME_GTP_C);
      EXPECT_EQ(fteid.teid, 0x00000001);
      EXPECT_EQ(fteid.ipv4_address.s_addr, 0x0441a8c0);
      EXPECT_EQ(fteid.ipv4_address.s_addr, htonl(0xc0a84104));
      EXPECT_TRUE(fteid.ipv4);
      EXPECT_EQ(fteid.ipv6, 0);
      EXPECT_TRUE(
          pMsg->isIeValid[NW_GTPV2C_IE_APN][NW_GTPV2C_IE_INSTANCE_ZERO]);
      EXPECT_TRUE(pMsg->isIeValid[NW_GTPV2C_IE_SELECTION_MODE]
                                 [NW_GTPV2C_IE_INSTANCE_ZERO]);
      EXPECT_TRUE(pMsg->isIeValid[NW_GTPV2C_IE_SERVING_NETWORK]
                                 [NW_GTPV2C_IE_INSTANCE_ZERO]);
      EXPECT_TRUE(pMsg->isIeValid[NW_GTPV2C_IE_BEARER_CONTEXT]
                                 [NW_GTPV2C_IE_INSTANCE_ZERO]);
      EXPECT_EQ(NW_OK, nwGtpv2cMsgDelete(s11_mme_stack_handle, hMsg));
    }
  }

  EXPECT_EQ(NW_OK, nwGtpv2cFinalize(s11_mme_stack_handle));
}

TEST(test_create_session_response_pdn1, create_session_response_pdn1) {
  nw_gtpv2c_stack_handle_t s11_mme_stack_handle = 0;

  uint8_t shift_buffer[1024];
  uint8_t create_session_response_pdn1[] = {
      0x48, 0x21, 0x00, 0x51, 0x00, 0x00, 0x00, 0x01, 0x00, 0x46, 0x80,
      0x00, 0x02, 0x00, 0x02, 0x00, 0x10, 0x00, 0x57, 0x00, 0x09, 0x00,
      0x8b, 0x00, 0x00, 0x00, 0x01, 0xc0, 0xa8, 0x41, 0x05, 0x4f, 0x00,
      0x05, 0x00, 0x01, 0xc0, 0xa8, 0x1d, 0x02, 0x48, 0x00, 0x08, 0x00,
      0x00, 0x00, 0x27, 0x10, 0x00, 0x00, 0x4e, 0x20, 0x4e, 0x00, 0x01,
      0x00, 0x80, 0x5d, 0x00, 0x18, 0x00, 0x49, 0x00, 0x01, 0x00, 0x05,
      0x02, 0x00, 0x02, 0x00, 0x10, 0x00, 0x57, 0x00, 0x09, 0x00, 0x81,
      0x00, 0x00, 0x00, 0x01, 0xc0, 0xa8, 0x41, 0x06};

  EXPECT_EQ(NW_OK, nwGtpv2cInitialize(&s11_mme_stack_handle));
  nw_gtpv2c_stack_t* pGtpv2c_stack =
      reinterpret_cast<nw_gtpv2c_stack_t*>(s11_mme_stack_handle);
  {
    EXPECT_TRUE(sizeof(shift_buffer) >=
                (128 + sizeof(create_session_response_pdn1)));
    // whatever is the byte alignment, should not fail
    for (int s = 0; s < 128; s++) {
      memset(shift_buffer, 0, sizeof(shift_buffer));
      memcpy(&shift_buffer[s], create_session_response_pdn1,
             sizeof(create_session_response_pdn1));
      nw_gtpv2c_msg_handle_t hMsg = 0;
      nw_gtpv2c_error_t error = {0};
      EXPECT_EQ(NW_OK, nwGtpv2cMsgFromBufferNew(
                           s11_mme_stack_handle, &shift_buffer[s],
                           sizeof(create_session_response_pdn1), &hMsg));
      nw_gtpv2c_msg_t* pMsg = reinterpret_cast<nw_gtpv2c_msg_t*>(hMsg);
      EXPECT_EQ(pMsg->msgLen, 85);
      EXPECT_EQ(pMsg->version, 2);
      EXPECT_TRUE(pMsg->teidPresent);
      EXPECT_EQ(pMsg->msgType, 33);
      EXPECT_EQ(pMsg->seqNum, 0x4680);
      NW_ASSERT(pGtpv2c_stack->pGtpv2cMsgIeParseInfo[pMsg->msgType]);
      EXPECT_EQ(NW_OK, nwGtpv2cMsgIeParse(
                           pGtpv2c_stack->pGtpv2cMsgIeParseInfo[pMsg->msgType],
                           hMsg, &error));
      EXPECT_TRUE(
          pMsg->isIeValid[NW_GTPV2C_IE_CAUSE][NW_GTPV2C_IE_INSTANCE_ZERO]);
      EXPECT_TRUE(
          pMsg->isIeValid[NW_GTPV2C_IE_FTEID][NW_GTPV2C_IE_INSTANCE_ZERO]);
      fteid_t fteid;
      memset(reinterpret_cast<void*>(&fteid), 0xAB, sizeof(fteid));
      nw_gtpv2c_ie_tlv_t* pIe = reinterpret_cast<nw_gtpv2c_ie_tlv_t*>(
          pMsg->pIe[NW_GTPV2C_IE_FTEID][NW_GTPV2C_IE_INSTANCE_ZERO]);
      uint16_t ieLength = ntohs(pIe->l);
      EXPECT_EQ(
          NW_OK,
          gtpv2c_fteid_ie_get(
              NW_GTPV2C_IE_FTEID, ieLength, NW_GTPV2C_IE_INSTANCE_ZERO,
              &pMsg->pIe[NW_GTPV2C_IE_FTEID][NW_GTPV2C_IE_INSTANCE_ZERO][4],
              reinterpret_cast<void*>(&fteid)));
      EXPECT_EQ(fteid.interface_type, S11_SGW_GTP_C);
      EXPECT_EQ(fteid.teid, 0x00000001);
      EXPECT_EQ(fteid.ipv4_address.s_addr, 0x0541a8c0);
      EXPECT_EQ(fteid.ipv4_address.s_addr, htonl(0xc0a84105));
      EXPECT_TRUE(fteid.ipv4);
      EXPECT_EQ(fteid.ipv6, 0);
      EXPECT_TRUE(
          pMsg->isIeValid[NW_GTPV2C_IE_PAA][NW_GTPV2C_IE_INSTANCE_ZERO]);
      EXPECT_TRUE(
          pMsg->isIeValid[NW_GTPV2C_IE_AMBR][NW_GTPV2C_IE_INSTANCE_ZERO]);
      EXPECT_TRUE(
          pMsg->isIeValid[NW_GTPV2C_IE_PCO][NW_GTPV2C_IE_INSTANCE_ZERO]);
      EXPECT_TRUE(pMsg->isIeValid[NW_GTPV2C_IE_BEARER_CONTEXT]
                                 [NW_GTPV2C_IE_INSTANCE_ZERO]);
      EXPECT_EQ(NW_OK, nwGtpv2cMsgDelete(s11_mme_stack_handle, hMsg));
    }
  }
  EXPECT_EQ(NW_OK, nwGtpv2cFinalize(s11_mme_stack_handle));
}

TEST(test_modify_bearer_request_pdn1, modify_bearer_request_pdn1) {
  nw_gtpv2c_stack_handle_t s11_mme_stack_handle = 0;

  uint8_t shift_buffer[1024];
  uint8_t modify_bearer_request_pdn1[] = {
      0x48, 0x22, 0x00, 0x1e, 0x00, 0x00, 0x00, 0x01, 0x00, 0x46, 0x81, 0x00,
      0x5d, 0x00, 0x12, 0x00, 0x49, 0x00, 0x01, 0x00, 0x05, 0x57, 0x00, 0x09,
      0x00, 0x80, 0x00, 0x00, 0x00, 0x01, 0xc0, 0xa8, 0x96, 0x01};
  EXPECT_EQ(NW_OK, nwGtpv2cInitialize(&s11_mme_stack_handle));
  nw_gtpv2c_stack_t* pGtpv2c_stack =
      reinterpret_cast<nw_gtpv2c_stack_t*>(s11_mme_stack_handle);
  {
    EXPECT_TRUE(sizeof(shift_buffer) >=
                (128 + sizeof(modify_bearer_request_pdn1)));
    // whatever is the byte alignment, should not fail
    for (int s = 0; s < 128; s++) {
      memset(shift_buffer, 0, sizeof(shift_buffer));
      memcpy(&shift_buffer[s], modify_bearer_request_pdn1,
             sizeof(modify_bearer_request_pdn1));
      nw_gtpv2c_msg_handle_t hMsg = 0;
      nw_gtpv2c_error_t error = {0};
      EXPECT_EQ(NW_OK, nwGtpv2cMsgFromBufferNew(
                           s11_mme_stack_handle, &shift_buffer[s],
                           sizeof(modify_bearer_request_pdn1), &hMsg));
      nw_gtpv2c_msg_t* pMsg = reinterpret_cast<nw_gtpv2c_msg_t*>(hMsg);
      EXPECT_EQ(pMsg->msgLen, 34);
      EXPECT_EQ(pMsg->version, 2);
      EXPECT_TRUE(pMsg->teidPresent);
      EXPECT_EQ(pMsg->msgType, NW_GTP_MODIFY_BEARER_REQ);
      EXPECT_EQ(pMsg->seqNum, 0x4681);
      NW_ASSERT(pGtpv2c_stack->pGtpv2cMsgIeParseInfo[NW_GTP_MODIFY_BEARER_REQ]);
      EXPECT_EQ(
          NW_OK,
          nwGtpv2cMsgIeParse(
              pGtpv2c_stack->pGtpv2cMsgIeParseInfo[NW_GTP_MODIFY_BEARER_REQ],
              hMsg, &error));
      EXPECT_TRUE(pMsg->isIeValid[NW_GTPV2C_IE_BEARER_CONTEXT]
                                 [NW_GTPV2C_IE_INSTANCE_ZERO]);
      EXPECT_EQ(NW_OK, nwGtpv2cMsgDelete(s11_mme_stack_handle, hMsg));
    }
  }

  EXPECT_EQ(NW_OK, nwGtpv2cFinalize(s11_mme_stack_handle));
}

TEST(test_modify_bearer_response_pdn1, modify_bearer_response_pdn1) {
  nw_gtpv2c_stack_handle_t s11_mme_stack_handle = 0;
  uint8_t shift_buffer[1024];
  uint8_t modify_bearer_response_pdn1[] = {
      0x48, 0x23, 0x00, 0x2a, 0x00, 0x00, 0x00, 0x01, 0x00, 0x46, 0x81, 0x00,
      0x02, 0x00, 0x02, 0x00, 0x10, 0x00, 0x5d, 0x00, 0x18, 0x00, 0x49, 0x00,
      0x01, 0x00, 0x05, 0x02, 0x00, 0x02, 0x00, 0x10, 0x00, 0x57, 0x00, 0x09,
      0x00, 0x81, 0x00, 0x00, 0x00, 0x01, 0xc0, 0xa8, 0x41, 0x06};

  EXPECT_EQ(NW_OK, nwGtpv2cInitialize(&s11_mme_stack_handle));
  nw_gtpv2c_stack_t* pGtpv2c_stack =
      reinterpret_cast<nw_gtpv2c_stack_t*>(s11_mme_stack_handle);
  {
    EXPECT_TRUE(sizeof(shift_buffer) >=
                (128 + sizeof(modify_bearer_response_pdn1)));
    // whatever is the byte alignment, should not fail
    for (int s = 0; s < 128; s++) {
      memset(shift_buffer, 0, sizeof(shift_buffer));
      memcpy(&shift_buffer[s], modify_bearer_response_pdn1,
             sizeof(modify_bearer_response_pdn1));
      nw_gtpv2c_msg_handle_t hMsg = 0;
      nw_gtpv2c_error_t error = {0};
      EXPECT_EQ(NW_OK, nwGtpv2cMsgFromBufferNew(
                           s11_mme_stack_handle, &shift_buffer[s],
                           sizeof(modify_bearer_response_pdn1), &hMsg));
      nw_gtpv2c_msg_t* pMsg = reinterpret_cast<nw_gtpv2c_msg_t*>(hMsg);
      EXPECT_EQ(pMsg->msgLen, 46);
      EXPECT_EQ(pMsg->version, 2);
      EXPECT_TRUE(pMsg->teidPresent);
      EXPECT_EQ(pMsg->msgType, NW_GTP_MODIFY_BEARER_RSP);
      EXPECT_EQ(pMsg->seqNum, 0x4681);
      NW_ASSERT(pGtpv2c_stack->pGtpv2cMsgIeParseInfo[NW_GTP_MODIFY_BEARER_RSP]);
      EXPECT_EQ(
          NW_OK,
          nwGtpv2cMsgIeParse(
              pGtpv2c_stack->pGtpv2cMsgIeParseInfo[NW_GTP_MODIFY_BEARER_RSP],
              hMsg, &error));
      EXPECT_TRUE(
          pMsg->isIeValid[NW_GTPV2C_IE_CAUSE][NW_GTPV2C_IE_INSTANCE_ZERO]);
      EXPECT_TRUE(pMsg->isIeValid[NW_GTPV2C_IE_BEARER_CONTEXT]
                                 [NW_GTPV2C_IE_INSTANCE_ZERO]);
      EXPECT_EQ(NW_OK, nwGtpv2cMsgDelete(s11_mme_stack_handle, hMsg));
    }
  }
  EXPECT_EQ(NW_OK, nwGtpv2cFinalize(s11_mme_stack_handle));
}

TEST(test_delete_session_request_pdn1, delete_session_request_pdn1) {
  nw_gtpv2c_stack_handle_t s11_mme_stack_handle = 0;
  uint8_t shift_buffer[1024];
  uint8_t delete_session_request_pdn1[] = {
      0x48, 0x24, 0x00, 0x21, 0x00, 0x00, 0x00, 0x63, 0x00, 0x52,
      0xd6, 0x00, 0x57, 0x00, 0x09, 0x00, 0x8a, 0x00, 0x00, 0x00,
      0x62, 0xc0, 0xa8, 0x41, 0x04, 0x49, 0x00, 0x01, 0x00, 0x05,
      0x4d, 0x00, 0x03, 0x00, 0x08, 0x00, 0x00};

  EXPECT_EQ(NW_OK, nwGtpv2cInitialize(&s11_mme_stack_handle));
  nw_gtpv2c_stack_t* pGtpv2c_stack =
      reinterpret_cast<nw_gtpv2c_stack_t*>(s11_mme_stack_handle);
  {
    EXPECT_TRUE(sizeof(shift_buffer) >=
                (128 + sizeof(delete_session_request_pdn1)));
    // whatever is the byte alignment, should not fail
    for (int s = 0; s < 128; s++) {
      memset(shift_buffer, 0, sizeof(shift_buffer));
      memcpy(&shift_buffer[s], delete_session_request_pdn1,
             sizeof(delete_session_request_pdn1));
      nw_gtpv2c_msg_handle_t hMsg = 0;
      nw_gtpv2c_error_t error = {0};
      EXPECT_TRUE(NW_OK == nwGtpv2cMsgFromBufferNew(
                               s11_mme_stack_handle, &shift_buffer[s],
                               sizeof(delete_session_request_pdn1), &hMsg));
      nw_gtpv2c_msg_t* pMsg = reinterpret_cast<nw_gtpv2c_msg_t*>(hMsg);
      EXPECT_EQ(pMsg->msgLen, 37);
      EXPECT_EQ(pMsg->version, 2);
      EXPECT_TRUE(pMsg->teidPresent);
      EXPECT_EQ(pMsg->msgType, NW_GTP_DELETE_SESSION_REQ);
      EXPECT_EQ(pMsg->seqNum, 0x52D6);
      NW_ASSERT(
          pGtpv2c_stack->pGtpv2cMsgIeParseInfo[NW_GTP_DELETE_SESSION_REQ]);
      EXPECT_EQ(
          NW_OK,
          nwGtpv2cMsgIeParse(
              pGtpv2c_stack->pGtpv2cMsgIeParseInfo[NW_GTP_DELETE_SESSION_REQ],
              hMsg, &error));
      EXPECT_TRUE(
          pMsg->isIeValid[NW_GTPV2C_IE_FTEID][NW_GTPV2C_IE_INSTANCE_ZERO]);
      fteid_t fteid;
      memset(reinterpret_cast<void*>(&fteid), 0xAB, sizeof(fteid));
      nw_gtpv2c_ie_tlv_t* pIe = reinterpret_cast<nw_gtpv2c_ie_tlv_t*>(
          pMsg->pIe[NW_GTPV2C_IE_FTEID][NW_GTPV2C_IE_INSTANCE_ZERO]);
      uint16_t ieLength = ntohs(pIe->l);
      EXPECT_EQ(
          NW_OK,
          gtpv2c_fteid_ie_get(
              NW_GTPV2C_IE_FTEID, ieLength, NW_GTPV2C_IE_INSTANCE_ZERO,
              &pMsg->pIe[NW_GTPV2C_IE_FTEID][NW_GTPV2C_IE_INSTANCE_ZERO][4],
              reinterpret_cast<void*>(&fteid)));
      EXPECT_EQ(fteid.interface_type, S11_MME_GTP_C);
      EXPECT_EQ(fteid.teid, 0x00000062);
      EXPECT_EQ(fteid.ipv4_address.s_addr, 0x0441a8c0);
      EXPECT_EQ(fteid.ipv4_address.s_addr, htonl(0xc0a84104));
      EXPECT_TRUE(fteid.ipv4);
      EXPECT_EQ(fteid.ipv6, 0);
      EXPECT_TRUE(
          pMsg->isIeValid[NW_GTPV2C_IE_EBI][NW_GTPV2C_IE_INSTANCE_ZERO]);
      EXPECT_TRUE(
          pMsg->isIeValid[NW_GTPV2C_IE_INDICATION][NW_GTPV2C_IE_INSTANCE_ZERO]);
      EXPECT_EQ(NW_OK, nwGtpv2cMsgDelete(s11_mme_stack_handle, hMsg));
    }
  }
  EXPECT_EQ(NW_OK, nwGtpv2cFinalize(s11_mme_stack_handle));
}
}  // namespace lte
}  // namespace magma
