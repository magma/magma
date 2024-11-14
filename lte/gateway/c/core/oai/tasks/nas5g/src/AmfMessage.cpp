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
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas5g/include/AmfMessage.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5gNasMessage.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_defs.hpp"

namespace magma5g {
AmfMsg::AmfMsg() { memset(&msg, 0, sizeof(MMsg_u)); };

AmfMsg::~AmfMsg(){};
MMsg_u::MMsg_u(){};
MMsg_u::~MMsg_u(){};

// Decode AMF NAS Header and Message
int AmfMsg::M5gNasMessageDecodeMsg(uint8_t* buffer, uint32_t len) {
  int header_result = 0;
  int decode_result = 0;

  if (len <= 0 || buffer == NULL) {
    OAILOG_ERROR(LOG_NAS5G, "Buffer is Empty");
    return (RETURNerror);
  }

  header_result = this->AmfMsgDecodeHeaderMsg(buffer, len);
  /** TODO: header bytes are being overwritten, this information is repeated in
   * the structure. */
  if (header_result <= 0) {
    OAILOG_ERROR(LOG_NAS5G, "Header Decoding Failed");
    return (RETURNerror);
  }

  OAILOG_DEBUG(LOG_NAS5G, "Header Decoded successfully");
  OAILOG_DEBUG(LOG_NAS5G,
               "EPD = 0x%X,  SecurityHeader =  0x%X, MessageType = 0x%X",
               static_cast<int>(this->header.extended_protocol_discriminator),
               static_cast<int>(this->header.sec_header_type),
               static_cast<int>(this->header.message_type));
  decode_result = this->AmfMsgDecodeMsg(buffer, len);
  if (decode_result <= 0) {
    OAILOG_ERROR(LOG_NAS5G, "Decode result error");
    return (RETURNerror);
  }
  return (header_result + decode_result);
}

// Encode AMF NAS  Header and Message
int AmfMsg::M5gNasMessageEncodeMsg(uint8_t* buffer, uint32_t len) {
  int header_result = 0;
  int encode_result = 0;

  /** TODO: header contents are being overwritten, consider optimizing this
   * structure */
  OAILOG_DEBUG(LOG_NAS5G, "Encoding NasMessage");
  if (len > 0 || buffer != NULL) {
    header_result = this->AmfMsgEncodeHeaderMsg(buffer, len);
    if (header_result <= 0) {
      OAILOG_ERROR(LOG_NAS5G, "Header encoding error");
      return (RETURNerror);
    }
  } else {
    OAILOG_ERROR(LOG_NAS5G, "Buffer is empty");
    return (RETURNerror);
  }
  encode_result = this->AmfMsgEncodeMsg(buffer, len);
  if (encode_result <= 0) {
    OAILOG_ERROR(LOG_NAS5G, "Encoding AMF Message Failed");
    return (RETURNerror);
  }

  return (encode_result);
}

// Decode AMF Message Header
int AmfMsg::AmfMsgDecodeHeaderMsg(uint8_t* buffer, uint32_t len) {
  AmfMsgHeader_s* hdr = &this->header;
  int size = 0;

  OAILOG_DEBUG(LOG_NAS5G, "Decoding AMF Message Header");
  if (len > 0 || buffer != NULL) {
    DECODE_U8(buffer + size, hdr->extended_protocol_discriminator, size);
    DECODE_U8(buffer + size, hdr->sec_header_type, size);
    DECODE_U8(buffer + size, hdr->message_type, size);
    OAILOG_DEBUG(LOG_NAS5G,
                 "EPD = 0x%X,  SecurityHeader =  0x%X, MessageType = 0x%X",
                 static_cast<int>(hdr->extended_protocol_discriminator),
                 static_cast<int>(hdr->sec_header_type),
                 static_cast<int>(hdr->message_type));
  } else {
    OAILOG_ERROR(LOG_NAS5G, "Buffer is Empty");
    return (RETURNerror);
  }

  if (hdr->extended_protocol_discriminator !=
      M5G_MOBILITY_MANAGEMENT_MESSAGES) {
    OAILOG_ERROR(LOG_NAS5G, "TLV not supported");
    return (TLV_PROTOCOL_NOT_SUPPORTED);
  }
  return (size);
}

// Encode AMF Message Header
int AmfMsg::AmfMsgEncodeHeaderMsg(uint8_t* buffer, uint32_t len) {
  AmfMsgHeader_s* hdr = &this->header;
  int size = 0;

  OAILOG_DEBUG(LOG_NAS5G, "Encoding AMF Message Header");

  if (len > 0 || buffer != NULL) {
    ENCODE_U8(buffer + size, hdr->extended_protocol_discriminator, size);
    ENCODE_U8(buffer + size, hdr->sec_header_type, size);
    ENCODE_U8(buffer + size, hdr->message_type, size);
    OAILOG_DEBUG(LOG_NAS5G,
                 "EPD = 0x%X,  SecurityHeader =  0x%X, MessageType = 0x%X",
                 static_cast<int>(hdr->extended_protocol_discriminator),
                 static_cast<int>(hdr->sec_header_type),
                 static_cast<int>(hdr->message_type));
  } else {
    OAILOG_ERROR(LOG_NAS5G, "Buffer is Empty");
    return (RETURNerror);
  }
  if ((unsigned char)hdr->extended_protocol_discriminator !=
      M5G_MOBILITY_MANAGEMENT_MESSAGES) {
    OAILOG_ERROR(LOG_NAS5G, "TLV not supported");
    return (TLV_PROTOCOL_NOT_SUPPORTED);
  }

  return (size);
}

// Decode AMF Message
int AmfMsg::AmfMsgDecodeMsg(uint8_t* buffer, uint32_t len) {
  int decode_result = 0;

  if (len <= 0 || buffer == NULL) {
    OAILOG_ERROR(LOG_NAS5G, "Buffer is Empty");
    return (RETURNerror);
  }

  OAILOG_DEBUG(
      LOG_NAS5G, "Decoding AMF message : %s",
      get_message_type_str(static_cast<uint8_t>(this->header.message_type))
          .c_str());

  switch (
      static_cast<M5GMessageType>((unsigned char)this->header.message_type)) {
    case M5GMessageType::REG_REQUEST:
      decode_result =
          this->msg.reg_request.DecodeRegistrationRequestMsg(buffer, len);
      break;
    case M5GMessageType::REG_ACCEPT:
      decode_result =
          this->msg.reg_accept.DecodeRegistrationAcceptMsg(buffer, len);
      break;
    case M5GMessageType::REG_COMPLETE:
      decode_result =
          this->msg.reg_complete.DecodeRegistrationCompleteMsg(buffer, len);
      break;
    case M5GMessageType::REG_REJECT:
      decode_result =
          this->msg.reg_reject.DecodeRegistrationRejectMsg(buffer, len);
      break;
    case M5GMessageType::M5G_IDENTITY_REQUEST:
      decode_result =
          this->msg.identity_request.DecodeIdentityRequestMsg(buffer, len);
      break;
    case M5GMessageType::M5G_IDENTITY_RESPONSE:
      decode_result =
          this->msg.identity_response.DecodeIdentityResponseMsg(buffer, len);
      break;
    case M5GMessageType::AUTH_REQUEST:
      decode_result =
          this->msg.auth_request.DecodeAuthenticationRequestMsg(buffer, len);
      break;
    case M5GMessageType::AUTH_RESPONSE:
      decode_result =
          this->msg.auth_response.DecodeAuthenticationResponseMsg(buffer, len);
      break;
    case M5GMessageType::AUTH_REJECT:
      decode_result =
          this->msg.auth_reject.DecodeAuthenticationRejectMsg(buffer, len);
      break;
    case M5GMessageType::AUTH_FAILURE:
      decode_result =
          this->msg.auth_failure.DecodeAuthenticationFailureMsg(buffer, len);
      break;
    case M5GMessageType::SEC_MODE_COMMAND:
      decode_result =
          this->msg.sec_mode_command.DecodeSecurityModeCommandMsg(buffer, len);
      break;
    case M5GMessageType::SEC_MODE_COMPLETE:
      decode_result = this->msg.sec_mode_complete.DecodeSecurityModeCompleteMsg(
          buffer, len);
      break;
    case M5GMessageType::SEC_MODE_REJECT:
      decode_result =
          this->msg.sec_mode_reject.DecodeSecurityModeRejectMsg(buffer, len);
      break;
    case M5GMessageType::DE_REG_REQUEST_UE_ORIGIN:
      decode_result =
          this->msg.de_reg_request.DecodeDeRegistrationRequestUEInitMsg(buffer,
                                                                        len);
      break;
    case M5GMessageType::DE_REG_ACCEPT_UE_ORIGIN:
      decode_result =
          this->msg.de_reg_accept.DecodeDeRegistrationAcceptUEInitMsg(buffer,
                                                                      len);
      break;
    case M5GMessageType::ULNASTRANSPORT:
      decode_result =
          this->msg.ul_nas_transport.DecodeULNASTransportMsg(buffer, len);
      break;
    case M5GMessageType::DLNASTRANSPORT:
      decode_result =
          this->msg.dl_nas_transport.DecodeDLNASTransportMsg(buffer, len);
      break;
    case M5GMessageType::M5G_SERVICE_REQUEST:
      decode_result = this->msg.svc_req.DecodeServiceRequestMsg(buffer, len);
      break;
    default:
      decode_result = TLV_WRONG_MESSAGE_TYPE;
  }
  return (decode_result);
}

// Encode AMF Message
int AmfMsg::AmfMsgEncodeMsg(uint8_t* buffer, uint32_t len) {
  int encode_result = 0;

  if (len <= 0 || buffer == NULL) {
    OAILOG_ERROR(LOG_NAS5G, "Buffer is Empty");
    return (RETURNerror);
  }

  OAILOG_DEBUG(
      LOG_NAS5G, "Encoding AMF message : %s",
      get_message_type_str(static_cast<uint8_t>(this->header.message_type))
          .c_str());
  switch (
      static_cast<M5GMessageType>((unsigned char)this->header.message_type)) {
    case M5GMessageType::REG_REQUEST:
      encode_result =
          this->msg.reg_request.EncodeRegistrationRequestMsg(buffer, len);
      break;
    case M5GMessageType::REG_ACCEPT:
      encode_result =
          this->msg.reg_accept.EncodeRegistrationAcceptMsg(buffer, len);
      break;
    case M5GMessageType::REG_COMPLETE:
      encode_result =
          this->msg.reg_complete.EncodeRegistrationCompleteMsg(buffer, len);
      break;
    case M5GMessageType::REG_REJECT:
      encode_result =
          this->msg.reg_reject.EncodeRegistrationRejectMsg(buffer, len);
      break;
    case M5GMessageType::M5G_IDENTITY_REQUEST:
      encode_result =
          this->msg.identity_request.EncodeIdentityRequestMsg(buffer, len);
      break;
    case M5GMessageType::M5G_IDENTITY_RESPONSE:
      encode_result =
          this->msg.identity_response.EncodeIdentityResponseMsg(buffer, len);
      break;
    case M5GMessageType::AUTH_REQUEST:
      encode_result =
          this->msg.auth_request.EncodeAuthenticationRequestMsg(buffer, len);
      break;
    case M5GMessageType::AUTH_RESPONSE:
      encode_result =
          this->msg.auth_response.EncodeAuthenticationResponseMsg(buffer, len);
      break;
    case M5GMessageType::AUTH_REJECT:
      encode_result =
          this->msg.auth_reject.EncodeAuthenticationRejectMsg(buffer, len);
      break;
    case M5GMessageType::AUTH_FAILURE:
      encode_result =
          this->msg.auth_failure.EncodeAuthenticationFailureMsg(buffer, len);
      break;
    case M5GMessageType::SEC_MODE_COMMAND:
      encode_result =
          this->msg.sec_mode_command.EncodeSecurityModeCommandMsg(buffer, len);
      break;
    case M5GMessageType::SEC_MODE_COMPLETE:
      encode_result = this->msg.sec_mode_complete.EncodeSecurityModeCompleteMsg(
          buffer, len);
      break;
    case M5GMessageType::SEC_MODE_REJECT:
      encode_result =
          this->msg.sec_mode_reject.EncodeSecurityModeRejectMsg(buffer, len);
      break;
    case M5GMessageType::DE_REG_ACCEPT_UE_ORIGIN:
      encode_result =
          this->msg.de_reg_accept.EncodeDeRegistrationAcceptUEInitMsg(buffer,
                                                                      len);
      break;
    case M5GMessageType::DLNASTRANSPORT:
      encode_result =
          this->msg.dl_nas_transport.EncodeDLNASTransportMsg(buffer, len);
      break;
    case M5GMessageType::M5G_SERVICE_ACCEPT:
      encode_result = this->msg.svc_acpt.EncodeServiceAcceptMsg(buffer, len);
      break;
    case M5GMessageType::M5G_SERVICE_REJECT:
      encode_result = this->msg.svc_rej.EncodeServiceRejectMsg(buffer, len);
      break;
    case M5GMessageType::PDU_SESSION_MODIFICATION_COMMAND:
      encode_result =
          this->msg.pdu_sess_mod_cmd.EncodePDUSessionModificationCommand(buffer,
                                                                         len);
      break;
    default:
      encode_result = TLV_WRONG_MESSAGE_TYPE;
  }
  return (encode_result);
}
}  // namespace magma5g
