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
PDUSessionEstablishmentAcceptMsg::PDUSessionEstablishmentAcceptMsg() {};
PDUSessionEstablishmentAcceptMsg::~PDUSessionEstablishmentAcceptMsg() {};

// Decode PDUSessionEstablishmentAccept Message and its IEs
int PDUSessionEstablishmentAcceptMsg::DecodePDUSessionEstablishmentAcceptMsg(
    PDUSessionEstablishmentAcceptMsg* pdu_session_estab_accept, uint8_t* buffer,
    uint32_t len) {
  uint32_t decoded = 0;
  int decoded_result = 0;
  uint8_t type_len = sizeof(uint8_t);
  uint8_t length_len = sizeof(uint8_t);

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(buffer,
                                       PDU_SESSION_ESTABLISH_ACPT_MIN_LEN, len);

  if ((decoded_result =
           pdu_session_estab_accept->extended_protocol_discriminator
               .DecodeExtendedProtocolDiscriminatorMsg(
                   &pdu_session_estab_accept->extended_protocol_discriminator,
                   0, buffer + decoded, len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }

  if ((decoded_result = pdu_session_estab_accept->pdu_session_identity
                            .DecodePDUSessionIdentityMsg(
                                &pdu_session_estab_accept->pdu_session_identity,
                                0, buffer + decoded, len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }

  if ((decoded_result = pdu_session_estab_accept->pti.DecodePTIMsg(
           &pdu_session_estab_accept->pti, 0, buffer + decoded,
           len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }

  if ((decoded_result =
           pdu_session_estab_accept->message_type.DecodeMessageTypeMsg(
               &pdu_session_estab_accept->message_type, 0, buffer + decoded,
               len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }

  {
    if ((decoded_result = pdu_session_estab_accept->ssc_mode.DecodeSSCModeMsg(
             &pdu_session_estab_accept->ssc_mode, 0, buffer + decoded,
             len - decoded)) < 0) {
      return decoded_result;
    } else {
      decoded += decoded_result;
    }

    if ((decoded_result =
             pdu_session_estab_accept->pdu_session_type.DecodePDUSessionTypeMsg(
                 &pdu_session_estab_accept->pdu_session_type, 0,
                 buffer + decoded, len - decoded)) < 0) {
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
    pdu_session_estab_accept->authorized_qosrules = blk2bstr(buffer, length);
    decoded += length;
  }

  if ((decoded_result =
           pdu_session_estab_accept->session_ambr.DecodeSessionAMBRMsg(
               &pdu_session_estab_accept->session_ambr, 0, buffer + decoded,
               len - decoded)) < 0) {
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
        if ((decoded_result =
                 pdu_session_estab_accept->pdu_address.DecodePDUAddressMsg(
                     &pdu_session_estab_accept->pdu_address,
                     static_cast<uint8_t>(M5GIei::PDU_ADDRESS),
                     buffer + decoded, len - decoded)) < 0) {
          return decoded_result;
        } else {
          decoded += decoded_result;
        }
        break;

      case M5GIei::S_NSSAI:
        if ((decoded_result = pdu_session_estab_accept->nssai.DecodeNSSAIMsg(
                 &pdu_session_estab_accept->nssai,
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
        if ((decoded_result = pdu_session_estab_accept->dnn.DecodeDNNMsg(
                 &pdu_session_estab_accept->dnn,
                 static_cast<uint8_t>(M5GIei::DNN), buffer + decoded,
                 len - decoded)) < 0) {
          return decoded_result;
        } else {
          decoded += decoded_result;
        }
        break;
      case M5GIei::EXTENDED_PROTOCOL_CONFIGURATION_OPTIONS:
        if ((decoded_result =
                 pdu_session_estab_accept->protocolconfigurationoptions
                     .DecodeProtocolConfigurationOptions(
                         &pdu_session_estab_accept
                              ->protocolconfigurationoptions,
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
    PDUSessionEstablishmentAcceptMsg* pdu_session_estab_accept, uint8_t* buffer,
    uint32_t len) {
  uint32_t encoded = 0;
  uint32_t encoded_result = 0;

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(buffer,
                                       PDU_SESSION_ESTABLISH_ACPT_MIN_LEN, len);

  if ((encoded_result =
           pdu_session_estab_accept->extended_protocol_discriminator
               .EncodeExtendedProtocolDiscriminatorMsg(
                   &pdu_session_estab_accept->extended_protocol_discriminator,
                   0, buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result = pdu_session_estab_accept->pdu_session_identity
                            .EncodePDUSessionIdentityMsg(
                                &pdu_session_estab_accept->pdu_session_identity,
                                0, buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result = pdu_session_estab_accept->pti.EncodePTIMsg(
           &pdu_session_estab_accept->pti, 0, buffer + encoded,
           len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result =
           pdu_session_estab_accept->message_type.EncodeMessageTypeMsg(
               &pdu_session_estab_accept->message_type, 0, buffer + encoded,
               len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result = pdu_session_estab_accept->ssc_mode.EncodeSSCModeMsg(
           &pdu_session_estab_accept->ssc_mode, 0, buffer + encoded,
           len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result =
           pdu_session_estab_accept->pdu_session_type.EncodePDUSessionTypeMsg(
               &pdu_session_estab_accept->pdu_session_type, 0, buffer + encoded,
               len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }

  if (blength(pdu_session_estab_accept->authorized_qosrules)) {
    // Encode the IE of Authorized QoS Rules
    // Encode the length of the IE
    IES_ENCODE_U16(buffer, encoded,
                   blength(pdu_session_estab_accept->authorized_qosrules));

    memcpy(buffer + encoded,
           pdu_session_estab_accept->authorized_qosrules->data,
           blength(pdu_session_estab_accept->authorized_qosrules));

    encoded += blength(pdu_session_estab_accept->authorized_qosrules);
  }

  if ((encoded_result =
           pdu_session_estab_accept->session_ambr.EncodeSessionAMBRMsg(
               &pdu_session_estab_accept->session_ambr, 0, buffer + encoded,
               len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result =
           pdu_session_estab_accept->pdu_address.EncodePDUAddressMsg(
               &pdu_session_estab_accept->pdu_address,
               static_cast<uint8_t>(M5GIei::PDU_ADDRESS), buffer + encoded,
               len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }

  if ((encoded_result = pdu_session_estab_accept->nssai.EncodeNSSAIMsg(
           &pdu_session_estab_accept->nssai,
           static_cast<uint8_t>(M5GIei::S_NSSAI), buffer + encoded,
           len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }

  if (blength(pdu_session_estab_accept->authorized_qosflowdescriptors)) {
    // Encode the IE of Authorized QOS Flow descriptions
    *(buffer + encoded) = PDU_SESSION_QOS_FLOW_DESC_IE_TYPE;
    encoded++;

    // Encode the length of the IE
    IES_ENCODE_U16(
        buffer, encoded,
        blength(pdu_session_estab_accept->authorized_qosflowdescriptors));

    memcpy(buffer + encoded,
           pdu_session_estab_accept->authorized_qosflowdescriptors->data,
           blength(pdu_session_estab_accept->authorized_qosflowdescriptors));

    encoded += blength(pdu_session_estab_accept->authorized_qosflowdescriptors);
  }

  if ((encoded_result =
           pdu_session_estab_accept->protocolconfigurationoptions
               .EncodeProtocolConfigurationOptions(
                   &pdu_session_estab_accept->protocolconfigurationoptions,
                   static_cast<uint8_t>(
                       M5GIei::EXTENDED_PROTOCOL_CONFIGURATION_OPTIONS),
                   buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if (pdu_session_estab_accept->dnn.dnn[0]) {
    if ((encoded_result = pdu_session_estab_accept->dnn.EncodeDNNMsg(
             &pdu_session_estab_accept->dnn, static_cast<uint8_t>(M5GIei::DNN),
             buffer + encoded, len - encoded)) < 0) {
      return encoded_result;
    } else {
      encoded += encoded_result;
    }
  }

  return encoded;
}
}  // namespace magma5g
