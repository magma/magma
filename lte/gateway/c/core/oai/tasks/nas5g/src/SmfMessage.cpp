/*
   Copyright 2020 The Magma Authors.
   This source code is licensed under the BSD-style license found in the
   LICENSE file in the root directory of this source tree.
   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
 */

#include <iostream>
#include <sstream>
#include "lte/gateway/c/core/oai/tasks/nas5g/include/SmfMessage.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5gNasMessage.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
SMsg_u::SMsg_u(){};
SMsg_u::~SMsg_u(){};
SmfMsg::SmfMsg(){};
SmfMsg::~SmfMsg(){};

// Decode SMF Message Header
int SmfMsg::SmfMsgDecodeHeaderMsg(
    SmfMsgHeader* hdr, uint8_t* buffer, uint32_t len) {
  int size = 0;

  MLOG(MDEBUG) << "SmfMsgDecodeHeaderMsg:" << std::endl;
  if (len > 0 || buffer != NULL) {
    DECODE_U8(buffer + size, hdr->extended_protocol_discriminator, size);
    DECODE_U8(buffer + size, hdr->pdu_session_id, size);
    DECODE_U8(buffer + size, hdr->procedure_transaction_id, size);
    DECODE_U8(buffer + size, hdr->message_type, size);
    MLOG(MDEBUG) << "epd = 0x" << std::hex
                 << int(hdr->extended_protocol_discriminator)
                 << "pdu session id = 0x" << std::hex
                 << int(hdr->pdu_session_id) << " procedure_transaction_id = 0x"
                 << std::hex << int(hdr->procedure_transaction_id)
                 << " message_type = 0x" << std::hex << int(hdr->message_type);
  } else {
    MLOG(MERROR) << "Error : Buffer is Empty" << std::endl;
    return (RETURNerror);
  }

  if (hdr->extended_protocol_discriminator != M5G_SESSION_MANAGEMENT_MESSAGES) {
    MLOG(MERROR) << "Error : TLV not supported" << std::endl;
    return (TLV_PROTOCOL_NOT_SUPPORTED);
  }
  return (size);
}

// Encode SMF Message Header
int SmfMsg::SmfMsgEncodeHeaderMsg(
    SmfMsgHeader* hdr, uint8_t* buffer, uint32_t len) {
  int size = 0;

  MLOG(MDEBUG) << "SmfMsgEncodeHeaderMsg:";
  if (len > 0 || buffer != NULL) {
    ENCODE_U8(buffer + size, hdr->extended_protocol_discriminator, size);
    ENCODE_U8(buffer + size, hdr->pdu_session_id, size);
    ENCODE_U8(buffer + size, hdr->procedure_transaction_id, size);
    ENCODE_U8(buffer + size, hdr->message_type, size);
    MLOG(MDEBUG) << " epd = 0x" << std::hex
                 << int(hdr->extended_protocol_discriminator)
                 << " pdu session id = 0x" << std::hex
                 << int(hdr->pdu_session_id) << " procedure_transaction_id = 0x"
                 << std::hex << int(hdr->procedure_transaction_id)
                 << " message_type = 0x" << std::hex << int(hdr->message_type);
  } else {
    MLOG(MERROR) << "Error : Buffer is Empty ";
    return (RETURNerror);
  }
  if ((unsigned char) hdr->extended_protocol_discriminator !=
      M5G_SESSION_MANAGEMENT_MESSAGES) {
    MLOG(MERROR) << "Error : TLV not supported";
    return (TLV_PROTOCOL_NOT_SUPPORTED);
  }
  return (size);
}

// Decode SMF Message
int SmfMsg::SmfMsgDecodeMsg(SmfMsg* msg, uint8_t* buffer, uint32_t len) {
  int decode_result = 0;
  int header_result = 0;

  MLOG(MDEBUG) << "SmfMsgDecodeMsg:" << std::endl;
  if (len <= 0 || buffer == NULL) {
    MLOG(MERROR) << "Error : Buffer is Empty" << std::endl;
    return (RETURNerror);
  }

  header_result = msg->SmfMsgDecodeHeaderMsg(&msg->header, buffer, len);
  if (header_result <= 0) {
    MLOG(MERROR) << "   Error : Header Decoding Failed" << std::dec
                 << RETURNerror;
    return (RETURNerror);
  }

  switch ((unsigned char) msg->header.message_type) {
    case PDU_SESSION_ESTABLISHMENT_REQUEST:
      decode_result = msg->msg.pdu_session_estab_request
                          .DecodePDUSessionEstablishmentRequestMsg(
                              &msg->msg.pdu_session_estab_request, buffer, len);
      break;
    case PDU_SESSION_RELEASE_REQUEST:
    case PDU_SESSION_RELEASE_COMPLETE:
      decode_result =
          msg->msg.pdu_session_release_request
              .DecodePDUSessionReleaseRequestMsg(
                  &msg->msg.pdu_session_release_request, buffer, len);
      break;
    case PDU_SESSION_MODIFICATION_REQUEST:
      decode_result = msg->msg.pdu_session_modif_request
                          .DecodePDUSessionModificationRequestMsg(
                              &msg->msg.pdu_session_modif_request, buffer, len);
      break;
    default:
      decode_result = TLV_WRONG_MESSAGE_TYPE;
  }
  return (decode_result);
}

// Encode SMF Message
int SmfMsg::SmfMsgEncodeMsg(SmfMsg* msg, uint8_t* buffer, uint32_t len) {
  int encode_result = 0;
  int header_result = 0;

  MLOG(MDEBUG) << " SmfMsgEncodeMsg : " << std::endl;
  if (len <= 0 || buffer == NULL) {
    MLOG(MERROR) << "Error : Buffer is Empty";
    return (RETURNerror);
  }

  header_result = msg->SmfMsgEncodeHeaderMsg(&msg->header, buffer, len);
  if (header_result <= 0) {
    MLOG(MERROR) << "   Error : Header Encoding Failed" << std::dec
                 << RETURNerror;
    return (RETURNerror);
  }

  switch ((unsigned char) msg->header.message_type) {
    case PDU_SESSION_ESTABLISHMENT_REQUEST:
      encode_result = msg->msg.pdu_session_estab_request
                          .EncodePDUSessionEstablishmentRequestMsg(
                              &msg->msg.pdu_session_estab_request, buffer, len);
      break;
    case PDU_SESSION_ESTABLISHMENT_ACCEPT:
      encode_result = msg->msg.pdu_session_estab_accept
                          .EncodePDUSessionEstablishmentAcceptMsg(
                              &msg->msg.pdu_session_estab_accept, buffer, len);
      break;
    case PDU_SESSION_ESTABLISHMENT_REJECT:
      encode_result = msg->msg.pdu_session_estab_reject
                          .EncodePDUSessionEstablishmentRejectMsg(
                              &msg->msg.pdu_session_estab_reject, buffer, len);
      break;
    case PDU_SESSION_MODIFICATION_REJECT:
      encode_result = msg->msg.pdu_session_modif_reject
                          .EncodePDUSessionModificationRejectMsg(
                              &msg->msg.pdu_session_modif_reject, buffer, len);
      break;
    case PDU_SESSION_RELEASE_REJECT:
      encode_result =
          msg->msg.pdu_session_release_reject.EncodePDUSessionReleaseRejectMsg(
              &msg->msg.pdu_session_release_reject, buffer, len);
      break;
    case PDU_SESSION_RELEASE_COMMAND:
      encode_result =
          msg->msg.pdu_session_release_command
              .EncodePDUSessionReleaseCommandMsg(
                  &msg->msg.pdu_session_release_command, buffer, len);
      break;
    default:
      encode_result = TLV_WRONG_MESSAGE_TYPE;
  }
  return (encode_result);
}
}  // namespace magma5g
