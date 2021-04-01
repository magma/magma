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
#include "M5GAuthenticationReject.h"
#include "M5GCommonDefs.h"

namespace magma5g {
AuthenticationRejectMsg::AuthenticationRejectMsg(){};
AuthenticationRejectMsg::~AuthenticationRejectMsg(){};

// Decoding Authentication Reject Message and its IEs
int AuthenticationRejectMsg::DecodeAuthenticationRejectMsg(
    AuthenticationRejectMsg* auth_reject, uint8_t* buffer, uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, AUTHENTICATION_REJECT_MINIMUM_LENGTH, len);

  if ((decoded_result = auth_reject->extended_protocol_discriminator
                            .DecodeExtendedProtocolDiscriminatorMsg(
                                &auth_reject->extended_protocol_discriminator,
                                0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = auth_reject->spare_half_octet.DecodeSpareHalfOctetMsg(
           &auth_reject->spare_half_octet, 0, buffer + decoded,
           len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           auth_reject->sec_header_type.DecodeSecurityHeaderTypeMsg(
               &auth_reject->sec_header_type, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = auth_reject->message_type.DecodeMessageTypeMsg(
           &auth_reject->message_type, 0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  return decoded;
}

// Encoding Authentication Reject Message and its IEs
int AuthenticationRejectMsg::EncodeAuthenticationRejectMsg(
    AuthenticationRejectMsg* auth_reject, uint8_t* buffer, uint32_t len) {
  uint32_t encoded   = 0;
  int encoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, AUTHENTICATION_REJECT_MINIMUM_LENGTH, len);

  if ((encoded_result = auth_reject->extended_protocol_discriminator
                            .EncodeExtendedProtocolDiscriminatorMsg(
                                &auth_reject->extended_protocol_discriminator,
                                0, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = auth_reject->spare_half_octet.EncodeSpareHalfOctetMsg(
           &auth_reject->spare_half_octet, 0, buffer + encoded,
           len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           auth_reject->sec_header_type.EncodeSecurityHeaderTypeMsg(
               &auth_reject->sec_header_type, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = auth_reject->message_type.EncodeMessageTypeMsg(
           &auth_reject->message_type, 0, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;

  return encoded;
}
}  // namespace magma5g
