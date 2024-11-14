/*
   Copyright 2022 The Magma Authors.
   This source code is licensed under the BSD-style license found in the
   LICENSE file in the root directory of this source tree.
   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
 */
#include <sstream>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GPDUSessionModificationCommand.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5gNasMessage.h"

namespace magma5g {
PDUSessionModificationCommand::PDUSessionModificationCommand() {}
PDUSessionModificationCommand::~PDUSessionModificationCommand() {}

int PDUSessionModificationCommand::EncodePDUSessionModificationCommand(
    uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;
  uint32_t encoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, PDU_SESSION_MODIFICATION_COMMAND_MIN_LEN, len);

  OAILOG_DEBUG(LOG_NAS5G, "EncodePDUSessionModificationCommand");
  if ((encoded_result = this->extended_protocol_discriminator
                            .EncodeExtendedProtocolDiscriminatorMsg(
                                0, buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result = this->pdu_session_identity.EncodePDUSessionIdentityMsg(
           0, buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result =
           this->pti.EncodePTIMsg(0, buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result = this->message_type.EncodeMessageTypeMsg(
           0, buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if (this->cause.iei) {
    if ((encoded_result = this->cause.EncodeM5GSMCauseMsg(0, buffer + encoded,
                                                          len - encoded)) < 0) {
      return encoded_result;
    } else {
      encoded += encoded_result;
    }
  }
  if (this->sessionambr.iei) {
    if ((encoded_result = this->sessionambr.EncodeSessionAMBRMsg(
             PDU_SESSION_AMBR_IE_TYPE, buffer + encoded, len - encoded)) < 0) {
      return encoded_result;
    } else {
      encoded += encoded_result;
    }
  }

  if (blength(this->authorized_qosrules)) {
    // Encode the IE of Authorized QoS Rules
    *(buffer + encoded++) = PDU_SESSION_QOS_RULES_IE_TYPE;

    // Encode the length of the IE
    IES_ENCODE_U16(buffer, encoded, blength(this->authorized_qosrules));

    memcpy(buffer + encoded, this->authorized_qosrules->data,
           blength(this->authorized_qosrules));

    encoded += blength(this->authorized_qosrules);
  }

  if (blength(this->authorized_qosflowdescriptors)) {
    // Encode the IE of Authorized QOS Flow descriptions
    *(buffer + encoded++) = PDU_SESSION_QOS_FLOW_DESC_IE_TYPE;

    // Encode the length of the IE
    IES_ENCODE_U16(buffer, encoded,
                   blength(this->authorized_qosflowdescriptors));

    memcpy(buffer + encoded, this->authorized_qosflowdescriptors->data,
           blength(this->authorized_qosflowdescriptors));

    encoded += blength(this->authorized_qosflowdescriptors);
  }

  return encoded;
}

int PDUSessionModificationCommand::DecodePDUSessionModificationCommand(
    uint8_t* buffer, uint32_t len) {
  uint32_t decoded = 0;
  uint32_t decoded_result = 0;

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, PDU_SESSION_MODIFICATION_COMMAND_MIN_LEN, len);

  OAILOG_DEBUG(LOG_NAS5G, "DecodePDUSessionModificationCommandMsg");
  if ((decoded_result = this->extended_protocol_discriminator
                            .DecodeExtendedProtocolDiscriminatorMsg(
                                0, buffer + decoded, len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  if ((decoded_result = this->pdu_session_identity.DecodePDUSessionIdentityMsg(
           0, buffer + decoded, len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  if ((decoded_result =
           this->pti.DecodePTIMsg(0, buffer + decoded, len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  if ((decoded_result = this->message_type.DecodeMessageTypeMsg(
           0, buffer + decoded, len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  while (decoded < len) {
    uint8_t ie_type = *(buffer + decoded);

    switch (ie_type) {
      case M5GSM_CAUSE: {
        if ((decoded_result = this->cause.DecodeM5GSMCauseMsg(
                 M5GSM_CAUSE, buffer + decoded, len - decoded)) < 0) {
          return decoded_result;
        } else {
          decoded += decoded_result;
        }
      } break;
      case PDU_SESSION_AMBR_IE_TYPE: {
        if ((decoded_result = this->sessionambr.DecodeSessionAMBRMsg(
                 PDU_SESSION_AMBR_IE_TYPE, buffer + decoded, len - decoded)) <
            0) {
          return decoded_result;
        } else {
          decoded += decoded_result;
        }
      } break;
      case PDU_SESSION_QOS_RULES_IE_TYPE: {
        // on the IE
        decoded += sizeof(ie_type);

        // Tracking the data length
        uint16_t qos_rules_buf_len = 0;
        IES_DECODE_U16(buffer, decoded, qos_rules_buf_len);

        // Store the information in Bstring
        this->authorized_qosrules =
            blk2bstr(buffer + decoded, qos_rules_buf_len);
        decoded += qos_rules_buf_len;
      } break;
      case PDU_SESSION_QOS_FLOW_DESC_IE_TYPE: {
        // on the IE
        decoded += sizeof(ie_type);

        // Tracking the data length
        uint16_t qos_flow_desc_buf_len = 0;
        IES_DECODE_U16(buffer, decoded, qos_flow_desc_buf_len);

        // Store the information in Bstring
        this->authorized_qosflowdescriptors =
            blk2bstr(buffer + decoded, qos_flow_desc_buf_len);
        decoded += qos_flow_desc_buf_len;
      } break;
      case REQUEST_EXTENDED_PROTOCOL_CONFIGURATION_OPTIONS_TYPE: {
        if ((decoded_result =
                 this->extprotocolconfigurationoptions
                     .DecodeProtocolConfigurationOptions(
                         REQUEST_EXTENDED_PROTOCOL_CONFIGURATION_OPTIONS_TYPE,
                         buffer + decoded, len - decoded)) < 0) {
          return decoded_result;
        } else {
          decoded += decoded_result;
        }
      } break;
      default: {
      } break;
    }
  }
  return decoded;
}
}  // namespace magma5g
