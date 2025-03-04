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

#include <sstream>
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GDLNASTransport.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
DLNASTransportMsg::DLNASTransportMsg() {};
DLNASTransportMsg::~DLNASTransportMsg() {};

// Decode DLNASTransport Message and its IEs
int DLNASTransportMsg::DecodeDLNASTransportMsg(
    DLNASTransportMsg* dl_nas_transport, uint8_t* buffer, uint32_t len) {
  uint32_t decoded = 0;
  int decoded_result = 0;

  // Checking Pointer
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(buffer, DL_NAS_TRANSPORT_MINIMUM_LENGTH,
                                       len);

  if ((decoded_result =
           dl_nas_transport->extended_protocol_discriminator
               .DecodeExtendedProtocolDiscriminatorMsg(
                   &dl_nas_transport->extended_protocol_discriminator, 0,
                   buffer + decoded, len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  if ((decoded_result =
           dl_nas_transport->spare_half_octet.DecodeSpareHalfOctetMsg(
               &dl_nas_transport->spare_half_octet, 0, buffer + decoded,
               len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  if ((decoded_result =
           dl_nas_transport->sec_header_type.DecodeSecurityHeaderTypeMsg(
               &dl_nas_transport->sec_header_type, 0, buffer + decoded,
               len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  if ((decoded_result = dl_nas_transport->message_type.DecodeMessageTypeMsg(
           &dl_nas_transport->message_type, 0, buffer + decoded,
           len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  if ((decoded_result =
           dl_nas_transport->spare_half_octet.DecodeSpareHalfOctetMsg(
               &dl_nas_transport->spare_half_octet, 0, buffer + decoded,
               len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  if ((decoded_result = dl_nas_transport->payload_container_type
                            .DecodePayloadContainerTypeMsg(
                                &dl_nas_transport->payload_container_type, 0,
                                buffer + decoded, len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  if ((decoded_result =
           dl_nas_transport->payload_container.DecodePayloadContainerMsg(
               &dl_nas_transport->payload_container, 0, buffer + decoded,
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
      case M5GIei::M5GMM_CAUSE: {
        if ((decoded_result = dl_nas_transport->m5gmm_cause.DecodeM5GMMCauseMsg(
                 &dl_nas_transport->m5gmm_cause,
                 static_cast<uint8_t>(M5GIei::M5GMM_CAUSE), buffer + decoded,
                 len - decoded)) < 0) {
          return decoded_result;
        } else {
          decoded += decoded_result;
        }
        break;
      }
      case M5GIei::PDU_SESSION_IDENTITY_2: {
        if ((decoded_result =
                 dl_nas_transport->pdu_session_identity
                     .DecodePDUSessionIdentityMsg(
                         &dl_nas_transport->pdu_session_identity,
                         static_cast<uint8_t>(M5GIei::PDU_SESSION_IDENTITY_2),
                         buffer + decoded, len - decoded)) < 0) {
          return decoded_result;
        } else {
          decoded += decoded_result;
        }
        break;
      }
      default:
        OAILOG_ERROR(LOG_NAS5G, "Unable to decode Optional Parameter");
        decoded_result = -1;
        break;
    }

    if (decoded_result < 0) {
      OAILOG_ERROR(LOG_NAS5G, "DLNASTransport Message Decoding FAILED");
      return decoded_result;
    }
  }
  return decoded;
}

// Encode DL NAS Transport Message and its IEs
int DLNASTransportMsg::EncodeDLNASTransportMsg(
    DLNASTransportMsg* dl_nas_transport, uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;

  int encoded_result = 0;

  // Check if we got a NDLL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, DL_NAS_TRANSPORT_MINIMUM_LENGTH,
                                       len);

  if ((encoded_result =
           dl_nas_transport->extended_protocol_discriminator
               .EncodeExtendedProtocolDiscriminatorMsg(
                   &dl_nas_transport->extended_protocol_discriminator, 0,
                   buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result =
           dl_nas_transport->spare_half_octet.EncodeSpareHalfOctetMsg(
               &dl_nas_transport->spare_half_octet, 0, buffer + encoded,
               len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result =
           dl_nas_transport->sec_header_type.EncodeSecurityHeaderTypeMsg(
               &dl_nas_transport->sec_header_type, 0, buffer + encoded,
               len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result = dl_nas_transport->message_type.EncodeMessageTypeMsg(
           &dl_nas_transport->message_type, 0, buffer + encoded,
           len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result =
           dl_nas_transport->spare_half_octet.EncodeSpareHalfOctetMsg(
               &dl_nas_transport->spare_half_octet, 0, buffer + encoded,
               len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result = dl_nas_transport->payload_container_type
                            .EncodePayloadContainerTypeMsg(
                                &dl_nas_transport->payload_container_type, 0,
                                buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result =
           dl_nas_transport->payload_container.EncodePayloadContainerMsg(
               &dl_nas_transport->payload_container, 0, buffer + encoded,
               len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if (dl_nas_transport->pdu_session_identity.pdu_session_id) {
    if ((encoded_result =
             dl_nas_transport->pdu_session_identity.EncodePDUSessionIdentityMsg(
                 &dl_nas_transport->pdu_session_identity,
                 static_cast<uint8_t>(M5GIei::PDU_SESSION_IDENTITY_2),
                 buffer + encoded, len - encoded)) < 0) {
      return encoded_result;
    } else {
      encoded += encoded_result;
    }
  }

  if (dl_nas_transport->m5gmm_cause.m5gmm_cause) {
    if ((encoded_result = dl_nas_transport->m5gmm_cause.EncodeM5GMMCauseMsg(
             &dl_nas_transport->m5gmm_cause,
             static_cast<uint8_t>(M5GIei::M5GMM_CAUSE), buffer + encoded,
             len - encoded)) < 0) {
      return encoded_result;
    } else {
      encoded += encoded_result;
    }
  }

  return encoded;
}
}  // namespace magma5g
