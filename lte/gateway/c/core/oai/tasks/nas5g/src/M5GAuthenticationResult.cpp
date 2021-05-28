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

#include <iostream>
#include <sstream>
#include "M5GAuthenticationResult.h"
#include "M5GCommonDefs.h"

using namespace std;
namespace magma5g {
AuthenticationResultMsg::AuthenticationResultMsg(){};
AuthenticationResultMsg::~AuthenticationResultMsg(){};

// Decode Authentication Result Message and its IEs
int AuthenticationResultMsg::DecodeAuthenticationResultMsg(
    AuthenticationResultMsg* auth_result, uint8_t* buffer, uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, AUTHENTICATION_RESULT_MINIMUM_LENGTH, len);

  MLOG(MDEBUG) << "\n\n---Decoding Authentication Result Message---\n" << endl;
  if ((decoded_result = auth_result->extended_protocol_discriminator
                            .DecodeExtendedProtocolDiscriminatorMsg(
                                &auth_result->extended_protocol_discriminator,
                                0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = auth_result->spare_half_octet.DecodeSpareHalfOctetMsg(
           &auth_result->spare_half_octet, 0, buffer + decoded,
           len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           auth_result->sec_header_type.DecodeSecurityHeaderTypeMsg(
               &auth_result->sec_header_type, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = auth_result->message_type.DecodeMessageTypeMsg(
           &auth_result->message_type, 0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = auth_result->spare_half_octet.DecodeSpareHalfOctetMsg(
           &auth_result->spare_half_octet, 0, buffer + decoded,
           len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           auth_result->nas_key_set_identifier.DecodeNASKeySetIdentifierMsg(
               &auth_result->nas_key_set_identifier, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = auth_result->eap_message.DecodeEAPMessageMsg(
           &auth_result->eap_message, 0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  return decoded;
};

// Encode Authentication Result Message and its IEs
int AuthenticationResultMsg::EncodeAuthenticationResultMsg(
    AuthenticationResultMsg* auth_result, uint8_t* buffer, uint32_t len) {
  uint32_t encoded   = 0;
  int encoded_result = 0;

  // Check if we got a NULL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, AUTHENTICATION_RESULT_MINIMUM_LENGTH, len);

  if ((encoded_result = auth_result->extended_protocol_discriminator
                            .EncodeExtendedProtocolDiscriminatorMsg(
                                &auth_result->extended_protocol_discriminator,
                                0, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = auth_result->spare_half_octet.EncodeSpareHalfOctetMsg(
           &auth_result->spare_half_octet, 0, buffer + encoded,
           len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           auth_result->sec_header_type.EncodeSecurityHeaderTypeMsg(
               &auth_result->sec_header_type, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = auth_result->message_type.EncodeMessageTypeMsg(
           &auth_result->message_type, 0, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = auth_result->spare_half_octet.EncodeSpareHalfOctetMsg(
           &auth_result->spare_half_octet, 0, buffer + encoded,
           len - encoded)) < 0)
    return encoded_result;
  if ((encoded_result =
           auth_result->nas_key_set_identifier.EncodeNASKeySetIdentifierMsg(
               &auth_result->nas_key_set_identifier, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = auth_result->eap_message.EncodeEAPMessageMsg(
           &auth_result->eap_message, 0, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  return encoded;
};
}  // namespace magma5g
