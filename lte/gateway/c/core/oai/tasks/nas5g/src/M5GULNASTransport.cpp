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
#include "M5GULNASTransport.h"
#include "M5GCommonDefs.h"

namespace magma5g {
ULNASTransportMsg::ULNASTransportMsg(){};
ULNASTransportMsg::~ULNASTransportMsg(){};

// Decode ULNASTransport Message and its IEs
int ULNASTransportMsg::DecodeULNASTransportMsg(
    ULNASTransportMsg* ul_nas_transport, uint8_t* buffer, uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;

  // Checking Pointer
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, UL_NAS_TRANSPORT_MINIMUM_LENGTH, len);

  MLOG(MDEBUG) << "DecodeULNASTransportMsg : \n";
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

  return decoded;
}

// Encode DL NAS Transport Message and its IEs
int ULNASTransportMsg::EncodeULNASTransportMsg(
    ULNASTransportMsg* ul_nas_transport, uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;

  MLOG(MDEBUG) << "EncodeULNASTransportMsg:";
  int encoded_result = 0;

  // Check if we got a NDLL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, UL_NAS_TRANSPORT_MINIMUM_LENGTH, len);

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

  return encoded;
}
}  // namespace magma5g
