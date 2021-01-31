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
#include "SmfMessage.h"
#include "M5gNasMessage.h"
#include "M5GCommonDefs.h"

using namespace std;
namespace magma5g {
SmfMsg::SmfMsg(){};
SmfMsg::~SmfMsg(){};

// Decode SMF Message Header
int SmfMsg::SmfMsgDecodeHeaderMsg(
    SmfMsgHeader* hdr, uint8_t* buffer, uint32_t len) {
  int size = 0;

  MLOG(MDEBUG) << "SmfMsgDecodeHeaderMsg:" << endl;
  if (len > 0 || buffer != NULL) {
    DECODE_U8(buffer + size, hdr->extended_protocol_discriminator, size);
    DECODE_U8(buffer + size, hdr->pdu_session_id, size);
    DECODE_U8(buffer + size, hdr->procedure_transaction_id, size);
    DECODE_U8(buffer + size, hdr->message_type, size);
    MLOG(MDEBUG) << "epd = 0x" << hex
                 << int(hdr->extended_protocol_discriminator)
                 << "pdu session id = 0x" << hex << int(hdr->pdu_session_id)
                 << " procedure_transaction_id = 0x" << hex
                 << int(hdr->procedure_transaction_id) << " message_type = 0x"
                 << hex << int(hdr->message_type);
  } else {
    MLOG(MERROR) << "Error : Buffer is Empty" << endl;
    return (RETURN_ERROR);
  }

  if (hdr->extended_protocol_discriminator != M5G_SESSION_MANAGEMENT_MESSAGES) {
    MLOG(MERROR) << "Error : TLV not supported" << endl;
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
    MLOG(MDEBUG) << " epd = 0x" << hex
                 << int(hdr->extended_protocol_discriminator)
                 << " pdu session id = 0x" << hex << int(hdr->pdu_session_id)
                 << " procedure_transaction_id = 0x" << hex
                 << int(hdr->procedure_transaction_id) << " message_type = 0x"
                 << hex << int(hdr->message_type);
  } else {
    MLOG(MERROR) << "Error : Buffer is Empty ";
    return (RETURN_ERROR);
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

  MLOG(MDEBUG) << "SmfMsgDecodeMsg:" << endl;
  if (len <= 0 || buffer == NULL) {
    MLOG(MERROR) << "Error : Buffer is Empty" << endl;
    return (RETURN_ERROR);
  }

  header_result = msg->SmfMsgDecodeHeaderMsg(&msg->header, buffer, len);
  if (header_result <= 0) {
    MLOG(MERROR) << "   Error : Header Decoding Failed" << std::dec
                 << RETURN_ERROR;
    return (RETURN_ERROR);
  }

  // buffer        = buffer + header_result;
  // decode_result = decode_result + header_result;

  MLOG(MDEBUG) << "msg type = 0x" << hex << int(msg->header.message_type);
  switch ((unsigned char) msg->header.message_type) {
    case PDU_SESSION_ESTABLISHMENT_REQUEST:
      decode_result = msg->pdu_session_estab_request
                          .DecodePDUSessionEstablishmentRequestMsg(
                              &msg->pdu_session_estab_request, buffer, len);
      break;
    case PDU_SESSION_RELEASE_REQUEST:
      decode_result =
          msg->pdu_session_release_request.DecodePDUSessionReleaseRequestMsg(
              &msg->pdu_session_release_request, buffer, len);
      break;
    case PDU_SESSION_MODIFICATION_REQUEST:
      decode_result =
          msg->pdu_session_modif_request.DecodePDUSessionModificationRequestMsg(
              &msg->pdu_session_modif_request, buffer, len);
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

  MLOG(MDEBUG) << " SmfMsgEncodeMsg : " << endl;
  if (len <= 0 || buffer == NULL) {
    MLOG(MERROR) << "Error : Buffer is Empty";
    return (RETURN_ERROR);
  }

  header_result = msg->SmfMsgEncodeHeaderMsg(&msg->header, buffer, len);
  if (header_result <= 0) {
    MLOG(MERROR) << "   Error : Header Encoding Failed" << std::dec
                 << RETURN_ERROR;
    return (RETURN_ERROR);
  }

  //  buffer        = buffer + header_result;
  //  encode_result = encode_result + header_result;

  MLOG(MDEBUG) << "msg type = 0x" << hex << int(msg->header.message_type);
  switch ((unsigned char) msg->header.message_type) {
    case PDU_SESSION_ESTABLISHMENT_ACCEPT:
      encode_result =
          msg->pdu_session_estab_accept.EncodePDUSessionEstablishmentAcceptMsg(
              &msg->pdu_session_estab_accept, buffer, len);
      break;
    case PDU_SESSION_ESTABLISHMENT_REJECT:
      encode_result =
          msg->pdu_session_estab_reject.EncodePDUSessionEstablishmentRejectMsg(
              &msg->pdu_session_estab_reject, buffer, len);
      break;
    case PDU_SESSION_MODIFICATION_REJECT:
      encode_result =
          msg->pdu_session_modif_reject.EncodePDUSessionModificationRejectMsg(
              &msg->pdu_session_modif_reject, buffer, len);
      break;
    case PDU_SESSION_RELEASE_REJECT:
      encode_result =
          msg->pdu_session_release_reject.EncodePDUSessionReleaseRejectMsg(
              &msg->pdu_session_release_reject, buffer, len);
      break;
    default:
      encode_result = TLV_WRONG_MESSAGE_TYPE;
  }
  return (encode_result);
}
}  // namespace magma5g
