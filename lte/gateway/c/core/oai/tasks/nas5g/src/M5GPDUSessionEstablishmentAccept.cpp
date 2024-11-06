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

#include <sstream>
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GPDUSessionEstablishmentAccept.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5gNasMessage.h"

namespace magma5g {
PDUSessionEstablishmentAcceptMsg::PDUSessionEstablishmentAcceptMsg(){};
PDUSessionEstablishmentAcceptMsg::~PDUSessionEstablishmentAcceptMsg(){};

// Decode PDUSessionEstablishmentAccept Message and its IEs
int PDUSessionEstablishmentAcceptMsg::DecodePDUSessionEstablishmentAcceptMsg(
    uint8_t* buffer, uint32_t len) {
  uint32_t decoded = 0;
  int decoded_result = 0;
  uint8_t type_len = sizeof(uint8_t);
  uint8_t length_len = sizeof(uint8_t);

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(buffer,
                                       PDU_SESSION_ESTABLISH_ACPT_MIN_LEN, len);

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

  {
    if ((decoded_result = this->ssc_mode.DecodeSSCModeMsg(0, buffer + decoded,
                                                          len - decoded)) < 0) {
      return decoded_result;
    } else {
      decoded += decoded_result;
    }

    if ((decoded_result = this->pdu_session_type.DecodePDUSessionTypeMsg(
             0, buffer + decoded, len - decoded)) < 0) {
      return decoded_result;
    } else {
      decoded += decoded_result;
    }
    decoded += 1;
  }

  // Decode Qos Rule msg
  {
    uint16_t length = 0;
    IES_DECODE_U16(buffer, decoded, length);
    this->authorized_qosrules = blk2bstr(buffer, length);
    decoded += length;
  }

  if ((decoded_result = this->session_ambr.DecodeSessionAMBRMsg(
           0, buffer + decoded, len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }

  while (decoded < len) {
    // Size is incremented for the unhandled types by 1 byte
    uint32_t type = *(buffer + decoded) >= 0x80 ? ((*(buffer + decoded)) & 0xf0)
                                                : (*(buffer + decoded));
    decoded_result = 0;

    switch (static_cast<M5GIei>(type)) {
      case M5GIei::PDU_ADDRESS:
        if ((decoded_result = this->pdu_address.DecodePDUAddressMsg(
                 static_cast<uint8_t>(M5GIei::PDU_ADDRESS), buffer + decoded,
                 len - decoded)) < 0) {
          return decoded_result;
        } else {
          decoded += decoded_result;
        }
        break;

      case M5GIei::S_NSSAI:
        if ((decoded_result = this->nssai.DecodeNSSAIMsg(
                 static_cast<uint8_t>(M5GIei::S_NSSAI), buffer + decoded,
                 len - decoded)) < 0) {
          return decoded_result;
        } else {
          decoded += decoded_result;
        }
        break;
      case M5GIei::QOS_FLOW_DESCRIPTIONS:
        // TLV Types.
        type_len = sizeof(uint16_t);
        length_len = sizeof(uint16_t);
        DECODE_U8(buffer + decoded + type_len, decoded_result, decoded);

        decoded += (length_len + decoded_result);
        break;
      case M5GIei::DNN:
        if ((decoded_result =
                 this->dnn.DecodeDNNMsg(static_cast<uint8_t>(M5GIei::DNN),
                                        buffer + decoded, len - decoded)) < 0) {
          return decoded_result;
        } else {
          decoded += decoded_result;
        }
        break;
      case M5GIei::EXTENDED_PROTOCOL_CONFIGURATION_OPTIONS:
        if ((decoded_result =
                 this->protocolconfigurationoptions
                     .DecodeProtocolConfigurationOptions(
                         static_cast<uint8_t>(
                             M5GIei::EXTENDED_PROTOCOL_CONFIGURATION_OPTIONS),
                         buffer + decoded, len - decoded)) < 0) {
          return decoded_result;
        } else {
          decoded += decoded_result;
        }
        break;

      default:
        decoded_result = -1;
        break;
    }
    if (decoded_result < 0) {
      return decoded_result;
    }
  }
  return decoded;
}

// Encode PDUSessionEstablishmentAccept Message and its IEs
int PDUSessionEstablishmentAcceptMsg::EncodePDUSessionEstablishmentAcceptMsg(
    uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;
  uint32_t encoded_result = 0;

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(buffer,
                                       PDU_SESSION_ESTABLISH_ACPT_MIN_LEN, len);

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
  if ((encoded_result = this->ssc_mode.EncodeSSCModeMsg(0, buffer + encoded,
                                                        len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result = this->pdu_session_type.EncodePDUSessionTypeMsg(
           0, buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }

  if (blength(this->authorized_qosrules)) {
    // Encode the IE of Authorized QoS Rules
    // Encode the length of the IE
    IES_ENCODE_U16(buffer, encoded, blength(this->authorized_qosrules));

    memcpy(buffer + encoded, this->authorized_qosrules->data,
           blength(this->authorized_qosrules));

    encoded += blength(this->authorized_qosrules);
  }

  if ((encoded_result = this->session_ambr.EncodeSessionAMBRMsg(
           0, buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result = this->pdu_address.EncodePDUAddressMsg(
           static_cast<uint8_t>(M5GIei::PDU_ADDRESS), buffer + encoded,
           len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }

  if ((encoded_result =
           this->nssai.EncodeNSSAIMsg(static_cast<uint8_t>(M5GIei::S_NSSAI),
                                      buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
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

  if ((encoded_result =
           this->protocolconfigurationoptions
               .EncodeProtocolConfigurationOptions(
                   static_cast<uint8_t>(
                       M5GIei::EXTENDED_PROTOCOL_CONFIGURATION_OPTIONS),
                   buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if (this->dnn.dnn[0]) {
    if ((encoded_result =
             this->dnn.EncodeDNNMsg(static_cast<uint8_t>(M5GIei::DNN),
                                    buffer + encoded, len - encoded)) < 0) {
      return encoded_result;
    } else {
      encoded += encoded_result;
    }
  }

  return encoded;
}
}  // namespace magma5g
