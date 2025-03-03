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

AmfMsg::~AmfMsg() {};
MMsg_u::MMsg_u() {};
MMsg_u::~MMsg_u() {};

// Decode AMF NAS Header and Message
int AmfMsg::M5gNasMessageDecodeMsg(AmfMsg* msg, uint8_t* buffer, uint32_t len) {
  int header_result = 0;
  int decode_result = 0;

  if (len > 0 || buffer != NULL) {
    header_result = msg->AmfMsgDecodeHeaderMsg(&msg->header, buffer, len);
    if (header_result <= 0) {
      OAILOG_ERROR(LOG_NAS5G, "Header Decoding Failed");
      return (RETURNerror);
    }
  } else {
    OAILOG_ERROR(LOG_NAS5G, "Buffer is Empty");
    return (RETURNerror);
  }
  OAILOG_DEBUG(LOG_NAS5G, "Header Decoded successfully");
  OAILOG_DEBUG(LOG_NAS5G,
               "EPD = 0x%X,  SecurityHeader =  0x%X, MessageType = 0x%X",
               static_cast<int>(msg->header.extended_protocol_discriminator),
               static_cast<int>(msg->header.sec_header_type),
               static_cast<int>(msg->header.message_type));
  decode_result = msg->AmfMsgDecodeMsg(msg, buffer, len);
  if (decode_result <= 0) {
    OAILOG_ERROR(LOG_NAS5G, "Decode result error");
    return (RETURNerror);
  }
  return (header_result + decode_result);
}

// Encode AMF NAS  Header and Message
int AmfMsg::M5gNasMessageEncodeMsg(AmfMsg* msg, uint8_t* buffer, uint32_t len) {
  int header_result = 0;
  int encode_result = 0;

  OAILOG_DEBUG(LOG_NAS5G, "Encoding NasMessage");
  if (len > 0 || buffer != NULL) {
    header_result = msg->AmfMsgEncodeHeaderMsg(&msg->header, buffer, len);
    if (header_result <= 0) {
      OAILOG_ERROR(LOG_NAS5G, "Header encoding error");
      return (RETURNerror);
    }
  } else {
    OAILOG_ERROR(LOG_NAS5G, "Buffer is empty");
    return (RETURNerror);
  }
  encode_result = msg->AmfMsgEncodeMsg(msg, buffer, len);
  if (encode_result <= 0) {
    OAILOG_ERROR(LOG_NAS5G, "Encoding AMF Message Failed");
    return (RETURNerror);
  }

  return (encode_result);
}

// Decode AMF Message Header
int AmfMsg::AmfMsgDecodeHeaderMsg(AmfMsgHeader_s* hdr, uint8_t* buffer,
                                  uint32_t len) {
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
int AmfMsg::AmfMsgEncodeHeaderMsg(AmfMsgHeader_s* hdr, uint8_t* buffer,
                                  uint32_t len) {
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
int AmfMsg::AmfMsgDecodeMsg(AmfMsg* msg, uint8_t* buffer, uint32_t len) {
  int decode_result = 0;

  if (len <= 0 || buffer == NULL) {
    OAILOG_ERROR(LOG_NAS5G, "Buffer is Empty");
    return (RETURNerror);
  }

  OAILOG_DEBUG(
      LOG_NAS5G, "Decoding AMF message : %s",
      get_message_type_str(static_cast<uint8_t>(msg->header.message_type))
          .c_str());

  switch (
      static_cast<M5GMessageType>((unsigned char)msg->header.message_type)) {
    case M5GMessageType::REG_REQUEST:
      decode_result = msg->msg.reg_request.DecodeRegistrationRequestMsg(
          &msg->msg.reg_request, buffer, len);
      break;
    case M5GMessageType::REG_ACCEPT:
      decode_result = msg->msg.reg_accept.DecodeRegistrationAcceptMsg(
          &msg->msg.reg_accept, buffer, len);
      break;
    case M5GMessageType::REG_COMPLETE:
      decode_result = msg->msg.reg_complete.DecodeRegistrationCompleteMsg(
          &msg->msg.reg_complete, buffer, len);
      break;
    case M5GMessageType::REG_REJECT:
      decode_result = msg->msg.reg_reject.DecodeRegistrationRejectMsg(
          &msg->msg.reg_reject, buffer, len);
      break;
    case M5GMessageType::M5G_IDENTITY_REQUEST:
      decode_result = msg->msg.identity_request.DecodeIdentityRequestMsg(
          &msg->msg.identity_request, buffer, len);
      break;
    case M5GMessageType::M5G_IDENTITY_RESPONSE:
      decode_result = msg->msg.identity_response.DecodeIdentityResponseMsg(
          &msg->msg.identity_response, buffer, len);
      break;
    case M5GMessageType::AUTH_REQUEST:
      decode_result = msg->msg.auth_request.DecodeAuthenticationRequestMsg(
          &msg->msg.auth_request, buffer, len);
      break;
    case M5GMessageType::AUTH_RESPONSE:
      decode_result = msg->msg.auth_response.DecodeAuthenticationResponseMsg(
          &msg->msg.auth_response, buffer, len);
      break;
    case M5GMessageType::AUTH_REJECT:
      decode_result = msg->msg.auth_reject.DecodeAuthenticationRejectMsg(
          &msg->msg.auth_reject, buffer, len);
      break;
    case M5GMessageType::AUTH_FAILURE:
      decode_result = msg->msg.auth_failure.DecodeAuthenticationFailureMsg(
          &msg->msg.auth_failure, buffer, len);
      break;
    case M5GMessageType::SEC_MODE_COMMAND:
      decode_result = msg->msg.sec_mode_command.DecodeSecurityModeCommandMsg(
          &msg->msg.sec_mode_command, buffer, len);
      break;
    case M5GMessageType::SEC_MODE_COMPLETE:
      decode_result = msg->msg.sec_mode_complete.DecodeSecurityModeCompleteMsg(
          &msg->msg.sec_mode_complete, buffer, len);
      break;
    case M5GMessageType::SEC_MODE_REJECT:
      decode_result = msg->msg.sec_mode_reject.DecodeSecurityModeRejectMsg(
          &msg->msg.sec_mode_reject, buffer, len);
      break;
    case M5GMessageType::DE_REG_REQUEST_UE_ORIGIN:
      decode_result =
          msg->msg.de_reg_request.DecodeDeRegistrationRequestUEInitMsg(
              &msg->msg.de_reg_request, buffer, len);
      break;
    case M5GMessageType::DE_REG_ACCEPT_UE_ORIGIN:
      decode_result =
          msg->msg.de_reg_accept.DecodeDeRegistrationAcceptUEInitMsg(
              &msg->msg.de_reg_accept, buffer, len);
      break;
    case M5GMessageType::ULNASTRANSPORT:
      decode_result = msg->msg.ul_nas_transport.DecodeULNASTransportMsg(
          &msg->msg.ul_nas_transport, buffer, len);
      break;
    case M5GMessageType::DLNASTRANSPORT:
      decode_result = msg->msg.dl_nas_transport.DecodeDLNASTransportMsg(
          &msg->msg.dl_nas_transport, buffer, len);
      break;
    case M5GMessageType::M5G_SERVICE_REQUEST:
      decode_result = msg->msg.svc_req.DecodeServiceRequestMsg(
          &msg->msg.svc_req, buffer, len);
      break;
    default:
      decode_result = TLV_WRONG_MESSAGE_TYPE;
  }
  return (decode_result);
}

// Encode AMF Message
int AmfMsg::AmfMsgEncodeMsg(AmfMsg* msg, uint8_t* buffer, uint32_t len) {
  int encode_result = 0;

  if (len <= 0 || buffer == NULL) {
    OAILOG_ERROR(LOG_NAS5G, "Buffer is Empty");
    return (RETURNerror);
  }

  OAILOG_DEBUG(
      LOG_NAS5G, "Encoding AMF message : %s",
      get_message_type_str(static_cast<uint8_t>(msg->header.message_type))
          .c_str());
  switch (
      static_cast<M5GMessageType>((unsigned char)msg->header.message_type)) {
    case M5GMessageType::REG_REQUEST:
      encode_result = msg->msg.reg_request.EncodeRegistrationRequestMsg(
          &msg->msg.reg_request, buffer, len);
      break;
    case M5GMessageType::REG_ACCEPT:
      encode_result = msg->msg.reg_accept.EncodeRegistrationAcceptMsg(
          &msg->msg.reg_accept, buffer, len);
      break;
    case M5GMessageType::REG_COMPLETE:
      encode_result = msg->msg.reg_complete.EncodeRegistrationCompleteMsg(
          &msg->msg.reg_complete, buffer, len);
      break;
    case M5GMessageType::REG_REJECT:
      encode_result = msg->msg.reg_reject.EncodeRegistrationRejectMsg(
          &msg->msg.reg_reject, buffer, len);
      break;
    case M5GMessageType::M5G_IDENTITY_REQUEST:
      encode_result = msg->msg.identity_request.EncodeIdentityRequestMsg(
          &msg->msg.identity_request, buffer, len);
      break;
    case M5GMessageType::M5G_IDENTITY_RESPONSE:
      encode_result = msg->msg.identity_response.EncodeIdentityResponseMsg(
          &msg->msg.identity_response, buffer, len);
      break;
    case M5GMessageType::AUTH_REQUEST:
      encode_result = msg->msg.auth_request.EncodeAuthenticationRequestMsg(
          &msg->msg.auth_request, buffer, len);
      break;
    case M5GMessageType::AUTH_RESPONSE:
      encode_result = msg->msg.auth_response.EncodeAuthenticationResponseMsg(
          &msg->msg.auth_response, buffer, len);
      break;
    case M5GMessageType::AUTH_REJECT:
      encode_result = msg->msg.auth_reject.EncodeAuthenticationRejectMsg(
          &msg->msg.auth_reject, buffer, len);
      break;
    case M5GMessageType::AUTH_FAILURE:
      encode_result = msg->msg.auth_failure.EncodeAuthenticationFailureMsg(
          &msg->msg.auth_failure, buffer, len);
      break;
    case M5GMessageType::SEC_MODE_COMMAND:
      encode_result = msg->msg.sec_mode_command.EncodeSecurityModeCommandMsg(
          &msg->msg.sec_mode_command, buffer, len);
      break;
    case M5GMessageType::SEC_MODE_COMPLETE:
      encode_result = msg->msg.sec_mode_complete.EncodeSecurityModeCompleteMsg(
          &msg->msg.sec_mode_complete, buffer, len);
      break;
    case M5GMessageType::SEC_MODE_REJECT:
      encode_result = msg->msg.sec_mode_reject.EncodeSecurityModeRejectMsg(
          &msg->msg.sec_mode_reject, buffer, len);
      break;
    case M5GMessageType::DE_REG_ACCEPT_UE_ORIGIN:
      encode_result =
          msg->msg.de_reg_accept.EncodeDeRegistrationAcceptUEInitMsg(
              &msg->msg.de_reg_accept, buffer, len);
      break;
    case M5GMessageType::DLNASTRANSPORT:
      encode_result = msg->msg.dl_nas_transport.EncodeDLNASTransportMsg(
          &msg->msg.dl_nas_transport, buffer, len);
      break;
    case M5GMessageType::M5G_SERVICE_ACCEPT:
      encode_result = msg->msg.svc_acpt.EncodeServiceAcceptMsg(
          &msg->msg.svc_acpt, buffer, len);
      break;
    case M5GMessageType::M5G_SERVICE_REJECT:
      encode_result = msg->msg.svc_rej.EncodeServiceRejectMsg(&msg->msg.svc_rej,
                                                              buffer, len);
      break;
    case M5GMessageType::PDU_SESSION_MODIFICATION_COMMAND:
      encode_result =
          msg->msg.pdu_sess_mod_cmd.EncodePDUSessionModificationCommand(
              &msg->msg.pdu_sess_mod_cmd, buffer, len);
      break;
    default:
      encode_result = TLV_WRONG_MESSAGE_TYPE;
  }
  return (encode_result);
}
}  // namespace magma5g
