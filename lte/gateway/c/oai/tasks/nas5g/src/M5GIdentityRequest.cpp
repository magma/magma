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
#include "M5GIdentityRequest.h"
#include "M5GCommonDefs.h"

namespace magma5g {
IdentityRequestMsg::IdentityRequestMsg(){};
IdentityRequestMsg::~IdentityRequestMsg(){};

// Decode IdentityRequest Message and its IEs
int IdentityRequestMsg::DecodeIdentityRequestMsg(
    IdentityRequestMsg* identity_request, uint8_t* buffer, uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;

  // Checking Pointer
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, IDENTITY_REQUEST_MINIMUM_LENGTH, len);

  MLOG(MDEBUG) << "DecodeIdentityRequestMsg : \n";
  if ((decoded_result =
           identity_request->extended_protocol_discriminator
               .DecodeExtendedProtocolDiscriminatorMsg(
                   &identity_request->extended_protocol_discriminator, 0,
                   buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           identity_request->spare_half_octet.DecodeSpareHalfOctetMsg(
               &identity_request->spare_half_octet, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           identity_request->sec_header_type.DecodeSecurityHeaderTypeMsg(
               &identity_request->sec_header_type, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = identity_request->message_type.DecodeMessageTypeMsg(
           &identity_request->message_type, 0, buffer + decoded,
           len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           identity_request->spare_half_octet.DecodeSpareHalfOctetMsg(
               &identity_request->spare_half_octet, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           identity_request->m5gs_identity_type.DecodeM5GSIdentityTypeMsg(
               &identity_request->m5gs_identity_type, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  return decoded;
}

// Encode Identity Request Message and its IEs
int IdentityRequestMsg::EncodeIdentityRequestMsg(
    IdentityRequestMsg* identity_request, uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;

  MLOG(MDEBUG) << "EncodeIdentityRequestMsg:";
  int encoded_result = 0;

  // Check if we got a NULL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, IDENTITY_REQUEST_MINIMUM_LENGTH, len);

  if ((encoded_result =
           identity_request->extended_protocol_discriminator
               .EncodeExtendedProtocolDiscriminatorMsg(
                   &identity_request->extended_protocol_discriminator, 0,
                   buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           identity_request->spare_half_octet.EncodeSpareHalfOctetMsg(
               &identity_request->spare_half_octet, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           identity_request->sec_header_type.EncodeSecurityHeaderTypeMsg(
               &identity_request->sec_header_type, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = identity_request->message_type.EncodeMessageTypeMsg(
           &identity_request->message_type, 0, buffer + encoded,
           len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           identity_request->spare_half_octet.EncodeSpareHalfOctetMsg(
               &identity_request->spare_half_octet, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           identity_request->m5gs_identity_type.EncodeM5GSIdentityTypeMsg(
               &identity_request->m5gs_identity_type, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;

  return encoded;
}
}  // namespace magma5g
