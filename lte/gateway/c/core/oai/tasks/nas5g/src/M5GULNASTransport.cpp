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
#include <iomanip>
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GULNASTransport.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
ULNASTransportMsg::ULNASTransportMsg() {};
ULNASTransportMsg::~ULNASTransportMsg() {};

// Decode ULNASTransport Message and its IEs
int ULNASTransportMsg::DecodeULNASTransportMsg(
    ULNASTransportMsg* ul_nas_transport, uint8_t* buffer, uint32_t len) {
  uint32_t decoded = 0;
  int decoded_result = 0;
  uint8_t type_len = 0;
  uint8_t length_len = 0;

  // Checking Pointer
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(buffer, UL_NAS_TRANSPORT_MINIMUM_LENGTH,
                                       len);

  if ((decoded_result =
           ul_nas_transport->extended_protocol_discriminator
               .DecodeExtendedProtocolDiscriminatorMsg(
                   &ul_nas_transport->extended_protocol_discriminator, 0,
                   buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  if ((decoded_result =
           ul_nas_transport->spare_half_octet.DecodeSpareHalfOctetMsg(
               &ul_nas_transport->spare_half_octet, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  if ((decoded_result =
           ul_nas_transport->sec_header_type.DecodeSecurityHeaderTypeMsg(
               &ul_nas_transport->sec_header_type, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  if ((decoded_result = ul_nas_transport->message_type.DecodeMessageTypeMsg(
           &ul_nas_transport->message_type, 0, buffer + decoded,
           len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  if ((decoded_result = ul_nas_transport->payload_container_type
                            .DecodePayloadContainerTypeMsg(
                                &ul_nas_transport->payload_container_type, 0,
                                buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  if ((decoded_result =
           ul_nas_transport->payload_container.DecodePayloadContainerMsg(
               &ul_nas_transport->payload_container, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  while (decoded < len) {
    // Size is incremented for the unhandled types by 1 byte
    uint32_t type = *(buffer + decoded) >= 0x80 ? ((*(buffer + decoded)) & 0xf0)
                                                : (*(buffer + decoded));
    decoded_result = 0;

    switch (static_cast<M5GIei>(type)) {
      case M5GIei::REQUEST_TYPE: {
        if ((decoded_result = ul_nas_transport->request_type.DecodeRequestType(
                 &ul_nas_transport->request_type,
                 static_cast<uint8_t>(M5GIei::REQUEST_TYPE), buffer + decoded,
                 len - decoded)) < 0) {
          return decoded_result;
        } else {
          decoded += decoded_result;
        }
        break;
      }
      case M5GIei::PDU_SESSION_IDENTITY_2:
      case M5GIei::OLD_PDU_SESSION_IDENTITY_2:
        decoded_result += 2;
        decoded += decoded_result;
        break;
      case M5GIei::MA_PDU_SESSION_INFORMATION:
      case M5GIei::RELEASE_ASSISTANCE_INDICATION:
        decoded_result += 1;
        decoded += decoded_result;
        break;
      case M5GIei::DNN:
        if ((decoded_result = ul_nas_transport->dnn.DecodeDNNMsg(
                 &ul_nas_transport->dnn, static_cast<uint8_t>(M5GIei::DNN),
                 buffer + decoded, len - decoded)) < 0) {
          return decoded_result;
        } else {
          decoded += decoded_result;
        }
        break;
      case M5GIei::S_NSSAI:
        if ((decoded_result = ul_nas_transport->nssai.DecodeNSSAIMsg(
                 &ul_nas_transport->nssai,
                 static_cast<uint8_t>(M5GIei::S_NSSAI), buffer + decoded,
                 len - decoded)) < 0) {
          return decoded_result;
        } else {
          decoded += decoded_result;
        }
        break;
      case M5GIei::ADDITIONAL_INFORMATION:
        // TLV Types. 1 byte for Type and 1 Byte for size
        type_len = sizeof(uint8_t);
        length_len = sizeof(uint8_t);
        DECODE_U8(buffer + decoded + type_len, decoded_result, decoded);

        decoded += (length_len + decoded_result);
        break;
      default:
        decoded_result = -1;
        break;
    }
    if (decoded_result < 0) {
      OAILOG_ERROR(LOG_NAS5G, "ULNASTransportMsg Decoding FAILED");
      return decoded_result;
    }
  }

  return decoded;
}

// Encode DL NAS Transport Message and its IEs
int ULNASTransportMsg::EncodeULNASTransportMsg(
    ULNASTransportMsg* ul_nas_transport, uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;

  int encoded_result = 0;

  // Check if we got a NDLL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, UL_NAS_TRANSPORT_MINIMUM_LENGTH,
                                       len);

  if ((encoded_result =
           ul_nas_transport->extended_protocol_discriminator
               .EncodeExtendedProtocolDiscriminatorMsg(
                   &ul_nas_transport->extended_protocol_discriminator, 0,
                   buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           ul_nas_transport->spare_half_octet.EncodeSpareHalfOctetMsg(
               &ul_nas_transport->spare_half_octet, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           ul_nas_transport->sec_header_type.EncodeSecurityHeaderTypeMsg(
               &ul_nas_transport->sec_header_type, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = ul_nas_transport->message_type.EncodeMessageTypeMsg(
           &ul_nas_transport->message_type, 0, buffer + encoded,
           len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           ul_nas_transport->spare_half_octet.EncodeSpareHalfOctetMsg(
               &ul_nas_transport->spare_half_octet, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = ul_nas_transport->payload_container_type
                            .EncodePayloadContainerTypeMsg(
                                &ul_nas_transport->payload_container_type, 0,
                                buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           ul_nas_transport->payload_container.EncodePayloadContainerMsg(
               &ul_nas_transport->payload_container, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;

  if ((uint32_t)ul_nas_transport->request_type.type_val) {
    if ((encoded_result = ul_nas_transport->request_type.EncodeRequestType(
             &ul_nas_transport->request_type,
             static_cast<uint8_t>(M5GIei::REQUEST_TYPE), buffer + encoded,
             len - encoded)) < 0) {
      return encoded_result;
    } else {
      encoded += encoded_result;
    }
  }

  if ((uint32_t)ul_nas_transport->dnn.len) {
    if ((encoded_result = ul_nas_transport->dnn.EncodeDNNMsg(
             &ul_nas_transport->dnn, static_cast<uint8_t>(M5GIei::DNN),
             buffer + encoded, len - encoded)) < 0) {
      return encoded_result;
    } else {
      encoded += encoded_result;
    }
  }

  if ((uint32_t)ul_nas_transport->nssai.len) {
    if ((encoded_result = ul_nas_transport->nssai.EncodeNSSAIMsg(
             &ul_nas_transport->nssai, static_cast<uint8_t>(M5GIei::S_NSSAI),
             buffer + encoded, len - encoded)) < 0) {
      return encoded_result;
    } else {
      encoded += encoded_result;
    }
  }
  return encoded;
}
}  // namespace magma5g
