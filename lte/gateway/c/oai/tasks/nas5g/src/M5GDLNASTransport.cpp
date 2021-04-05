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
#include "M5GDLNASTransport.h"
#include "M5GCommonDefs.h"

namespace magma5g {
DLNASTransportMsg::DLNASTransportMsg(){};
DLNASTransportMsg::~DLNASTransportMsg(){};

// Decode DLNASTransport Message and its IEs
int DLNASTransportMsg::DecodeDLNASTransportMsg(
    DLNASTransportMsg* dl_nas_transport, uint8_t* buffer, uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;

  // Checking Pointer
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, DL_NAS_TRANSPORT_MINIMUM_LENGTH, len);

  MLOG(MDEBUG) << "DecodeDLNASTransportMsg : \n";
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
  return decoded;
}

// Encode DL NAS Transport Message and its IEs
int DLNASTransportMsg::EncodeDLNASTransportMsg(
    DLNASTransportMsg* dl_nas_transport, uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;

  MLOG(MDEBUG) << "EncodeDLNASTransportMsg:";
  int encoded_result = 0;

  // Check if we got a NDLL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, DL_NAS_TRANSPORT_MINIMUM_LENGTH, len);

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
  if ((encoded_result =
           dl_nas_transport->pdu_session_identity.EncodePDUSessionIdentityMsg(
               &dl_nas_transport->pdu_session_identity, PDU_SESSION_IDENTITY,
               buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  return encoded;
}
}  // namespace magma5g
