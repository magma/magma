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
#include "M5GIdentityResponse.h"
#include "M5GCommonDefs.h"

namespace magma5g {
IdentityResponseMsg::IdentityResponseMsg(){};
IdentityResponseMsg::~IdentityResponseMsg(){};

// Decode IdentityResponse Message and its IEs
int IdentityResponseMsg::DecodeIdentityResponseMsg(
    IdentityResponseMsg* identity_response, uint8_t* buffer, uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;

  // Checking Pointer
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, IDENTITY_RESPONSE_MINIMUM_LENGTH, len);

  MLOG(MDEBUG) << "DecodeIdentityResponseMsg : \n";
  if ((decoded_result =
           identity_response->extended_protocol_discriminator
               .DecodeExtendedProtocolDiscriminatorMsg(
                   &identity_response->extended_protocol_discriminator, 0,
                   buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           identity_response->spare_half_octet.DecodeSpareHalfOctetMsg(
               &identity_response->spare_half_octet, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           identity_response->sec_header_type.DecodeSecurityHeaderTypeMsg(
               &identity_response->sec_header_type, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = identity_response->message_type.DecodeMessageTypeMsg(
           &identity_response->message_type, 0, buffer + decoded,
           len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           identity_response->m5gs_mobile_identity.DecodeM5GSMobileIdentityMsg(
               &identity_response->m5gs_mobile_identity, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  return decoded;
}

// Encode Identity Response Message and its IEs
int IdentityResponseMsg::EncodeIdentityResponseMsg(
    IdentityResponseMsg* identity_response, uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;

  MLOG(MDEBUG) << "EncodeIdentityResponseMsg:";
  int encoded_result = 0;

  // Check if we got a NULL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, IDENTITY_RESPONSE_MINIMUM_LENGTH, len);

  if ((encoded_result =
           identity_response->extended_protocol_discriminator
               .EncodeExtendedProtocolDiscriminatorMsg(
                   &identity_response->extended_protocol_discriminator, 0,
                   buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           identity_response->spare_half_octet.EncodeSpareHalfOctetMsg(
               &identity_response->spare_half_octet, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           identity_response->sec_header_type.EncodeSecurityHeaderTypeMsg(
               &identity_response->sec_header_type, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = identity_response->message_type.EncodeMessageTypeMsg(
           &identity_response->message_type, 0, buffer + encoded,
           len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           identity_response->m5gs_mobile_identity.EncodeM5GSMobileIdentityMsg(
               &identity_response->m5gs_mobile_identity, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;

  return encoded;
}
}  // namespace magma5g
