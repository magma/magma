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

#include <cstring>
#include <iostream>
#include <sstream>
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas5g/include/SmfMessage.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5gNasMessage.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_defs.hpp"

namespace magma5g {
SMsg_u::SMsg_u() { memset(this, 0, sizeof(SMsg_u)); };
SMsg_u::~SMsg_u(){};
SmfMsg::SmfMsg() { memset(this, 0, sizeof(SmfMsg)); };
SmfMsg::~SmfMsg(){};

// Decode SMF Message Header
int SmfMsg::SmfMsgDecodeHeaderMsg(uint8_t* buffer, uint32_t len) {
  int size = 0;

  OAILOG_DEBUG(LOG_NAS5G, "Decoding SMF message header");
  if (len > 0 || buffer != NULL) {
    DECODE_U8(buffer + size, this->header.extended_protocol_discriminator,
              size);
    DECODE_U8(buffer + size, this->header.pdu_session_id, size);
    DECODE_U8(buffer + size, this->header.procedure_transaction_id, size);
    DECODE_U8(buffer + size, this->header.message_type, size);
    OAILOG_DEBUG(
        LOG_NAS5G,
        "EPD = 0x%X, PDUSessionID = 0x%X, ProcedureTransactionID = 0x%X,  "
        "MessageType = 0x%X",
        static_cast<int>(this->header.extended_protocol_discriminator),
        static_cast<int>(this->header.pdu_session_id),
        static_cast<int>(this->header.procedure_transaction_id),
        static_cast<int>(this->header.message_type));
  } else {
    OAILOG_ERROR(LOG_NAS5G, "Buffer is Empty");
    return (RETURNerror);
  }

  if (this->header.extended_protocol_discriminator !=
      M5G_SESSION_MANAGEMENT_MESSAGES) {
    OAILOG_ERROR(LOG_NAS5G, "TLV not supported");
    return (TLV_PROTOCOL_NOT_SUPPORTED);
  }
  return (size);
}

// Encode SMF Message Header
int SmfMsg::SmfMsgEncodeHeaderMsg(uint8_t* buffer, uint32_t len) {
  int size = 0;

  OAILOG_DEBUG(LOG_NAS5G, "Encoding SMF message header");
  if (len > 0 || buffer != NULL) {
    ENCODE_U8(buffer + size, this->header.extended_protocol_discriminator,
              size);
    ENCODE_U8(buffer + size, this->header.pdu_session_id, size);
    ENCODE_U8(buffer + size, this->header.procedure_transaction_id, size);
    ENCODE_U8(buffer + size, this->header.message_type, size);
    OAILOG_DEBUG(
        LOG_NAS5G,
        "EPD = 0x%X, PDUSessionID = 0x%X, ProcedureTransactionID = 0x%X,  "
        "MessageType = 0x%X",
        static_cast<int>(this->header.extended_protocol_discriminator),
        static_cast<int>(this->header.pdu_session_id),
        static_cast<int>(this->header.procedure_transaction_id),
        static_cast<int>(this->header.message_type));
  } else {
    OAILOG_ERROR(LOG_NAS5G, "Buffer is Empty");
    return (RETURNerror);
  }
  if ((unsigned char)this->header.extended_protocol_discriminator !=
      M5G_SESSION_MANAGEMENT_MESSAGES) {
    OAILOG_ERROR(LOG_NAS5G, "TLV not supported");
    return (TLV_PROTOCOL_NOT_SUPPORTED);
  }
  return (size);
}

// Decode SMF Message
int SmfMsg::SmfMsgDecodeMsg(uint8_t* buffer, uint32_t len) {
  int decode_result = 0;
  int header_result = 0;

  if (len <= 0 || buffer == NULL) {
    OAILOG_ERROR(LOG_NAS5G, "Buffer is Empty");
    return (RETURNerror);
  }

  header_result = this->SmfMsgDecodeHeaderMsg(buffer, len);
  if (header_result <= 0) {
    OAILOG_ERROR(LOG_NAS5G, "Header Decoding Failed");
    return (RETURNerror);
  }

  OAILOG_DEBUG(
      LOG_NAS5G, "Decoding SMF message : %s",
      get_message_type_str(static_cast<uint8_t>(this->header.message_type))
          .c_str());

  switch (
      static_cast<M5GMessageType>((unsigned char)this->header.message_type)) {
    case M5GMessageType::PDU_SESSION_ESTABLISHMENT_REQUEST:
      decode_result = this->msg.pdu_session_estab_request
                          .DecodePDUSessionEstablishmentRequestMsg(buffer, len);
      break;
    case M5GMessageType::PDU_SESSION_ESTABLISHMENT_REJECT:
      decode_result = this->msg.pdu_session_estab_reject
                          .DecodePDUSessionEstablishmentRejectMsg(buffer, len);
      break;
    case M5GMessageType::PDU_SESSION_ESTABLISHMENT_ACCEPT:
      decode_result = this->msg.pdu_session_estab_accept
                          .DecodePDUSessionEstablishmentAcceptMsg(buffer, len);
      break;
    case M5GMessageType::PDU_SESSION_RELEASE_REQUEST:
    case M5GMessageType::PDU_SESSION_RELEASE_COMPLETE:
      decode_result = this->msg.pdu_session_release_request
                          .DecodePDUSessionReleaseRequestMsg(buffer, len);
      break;
    case M5GMessageType::PDU_SESSION_MODIFICATION_REQUEST:
      decode_result = this->msg.pdu_session_modif_request
                          .DecodePDUSessionModificationRequestMsg(buffer, len);
      break;
    case M5GMessageType::PDU_SESSION_MODIFICATION_COMPLETE:
      decode_result =
          this->msg.pdu_sess_mod_com.DecodePDUSessionModificationComplete(
              buffer, len);
      break;
    case M5GMessageType::PDU_SESSION_MODIFICATION_COMMAND_REJECT:
      decode_result =
          this->msg.pdu_sess_mod_cmd_rej
              .DecodePDUSessionModificationCommandReject(buffer, len);
      break;
    default:
      decode_result = TLV_WRONG_MESSAGE_TYPE;
  }
  return (decode_result);
}

// Encode SMF Message
int SmfMsg::SmfMsgEncodeMsg(uint8_t* buffer, uint32_t len) {
  int encode_result = 0;
  int header_result = 0;

  if (len <= 0 || buffer == NULL) {
    OAILOG_ERROR(LOG_NAS5G, "Buffer is Empty");
    return (RETURNerror);
  }

  header_result = this->SmfMsgEncodeHeaderMsg(buffer, len);
  if (header_result <= 0) {
    OAILOG_ERROR(LOG_NAS5G, "Header Encoding Failed");
    return (RETURNerror);
  }

  OAILOG_DEBUG(
      LOG_NAS5G, "Encoding SMF message : %s",
      get_message_type_str(static_cast<uint8_t>(this->header.message_type))
          .c_str());

  switch (
      static_cast<M5GMessageType>((unsigned char)this->header.message_type)) {
    case M5GMessageType::PDU_SESSION_ESTABLISHMENT_REQUEST:
      encode_result = this->msg.pdu_session_estab_request
                          .EncodePDUSessionEstablishmentRequestMsg(buffer, len);
      break;
    case M5GMessageType::PDU_SESSION_ESTABLISHMENT_ACCEPT:
      encode_result = this->msg.pdu_session_estab_accept
                          .EncodePDUSessionEstablishmentAcceptMsg(buffer, len);
      break;
    case M5GMessageType::PDU_SESSION_ESTABLISHMENT_REJECT:
      encode_result = this->msg.pdu_session_estab_reject
                          .EncodePDUSessionEstablishmentRejectMsg(buffer, len);
      break;
    case M5GMessageType::PDU_SESSION_MODIFICATION_REJECT:
      encode_result = this->msg.pdu_session_modif_reject
                          .EncodePDUSessionModificationRejectMsg(buffer, len);
      break;
    case M5GMessageType::PDU_SESSION_MODIFICATION_COMMAND:
      encode_result =
          this->msg.pdu_sess_mod_cmd.EncodePDUSessionModificationCommand(buffer,
                                                                         len);
      break;
    case M5GMessageType::PDU_SESSION_RELEASE_REJECT:
      encode_result =
          this->msg.pdu_session_release_reject.EncodePDUSessionReleaseRejectMsg(
              buffer, len);
      break;
    case M5GMessageType::PDU_SESSION_RELEASE_COMMAND:
      encode_result = this->msg.pdu_session_release_command
                          .EncodePDUSessionReleaseCommandMsg(buffer, len);
      break;
    default:
      encode_result = TLV_WRONG_MESSAGE_TYPE;
      break;
  }
  return (encode_result);
}
}  // namespace magma5g
