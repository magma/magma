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
#include "M5GDeRegistrationRequestUEInit.h"
#include "M5GCommonDefs.h"

using namespace std;
namespace magma5g {
DeRegistrationRequestUEInitMsg::DeRegistrationRequestUEInitMsg(){};
DeRegistrationRequestUEInitMsg::~DeRegistrationRequestUEInitMsg(){};

// Decode De Registration Request(UE) Message and its IEs
int DeRegistrationRequestUEInitMsg::DecodeDeRegistrationRequestUEInitMsg(
    DeRegistrationRequestUEInitMsg* de_reg_request, uint8_t* buffer,
    uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, DEREGISTRATION_REQUEST_UEINIT_MINIMUM_LENGTH, len);

  MLOG(MDEBUG) << "\n\n---Decoding De-Registration Request Message---\n"
               << endl;
  if ((decoded_result =
           de_reg_request->extended_protocol_discriminator
               .DecodeExtendedProtocolDiscriminatorMsg(
                   &de_reg_request->extended_protocol_discriminator, 0,
                   buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           de_reg_request->spare_half_octet.DecodeSpareHalfOctetMsg(
               &de_reg_request->spare_half_octet, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           de_reg_request->sec_header_type.DecodeSecurityHeaderTypeMsg(
               &de_reg_request->sec_header_type, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = de_reg_request->message_type.DecodeMessageTypeMsg(
           &de_reg_request->message_type, 0, buffer + decoded, len - decoded)) <
      0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           de_reg_request->m5gs_de_reg_type.DecodeM5GSDeRegistrationTypeMsg(
               &de_reg_request->m5gs_de_reg_type, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           de_reg_request->nas_key_set_identifier.DecodeNASKeySetIdentifierMsg(
               &de_reg_request->nas_key_set_identifier, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           de_reg_request->m5gs_mobile_identity.DecodeM5GSMobileIdentityMsg(
               &de_reg_request->m5gs_mobile_identity, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  return decoded;
};

// Encode De Registration Request(UE) Message and its IEs
int DeRegistrationRequestUEInitMsg::EncodeDeRegistrationRequestUEInitMsg(
    DeRegistrationRequestUEInitMsg* de_reg_request, uint8_t* buffer,
    uint32_t len) {
  uint32_t encoded   = 0;
  int encoded_result = 0;

  // Check if we got a NULL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, DEREGISTRATION_REQUEST_UEINIT_MINIMUM_LENGTH, len);

  if ((encoded_result =
           de_reg_request->extended_protocol_discriminator
               .EncodeExtendedProtocolDiscriminatorMsg(
                   &de_reg_request->extended_protocol_discriminator, 0,
                   buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           de_reg_request->spare_half_octet.EncodeSpareHalfOctetMsg(
               &de_reg_request->spare_half_octet, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           de_reg_request->sec_header_type.EncodeSecurityHeaderTypeMsg(
               &de_reg_request->sec_header_type, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = de_reg_request->message_type.EncodeMessageTypeMsg(
           &de_reg_request->message_type, 0, buffer + encoded, len - encoded)) <
      0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           de_reg_request->m5gs_de_reg_type.EncodeM5GSDeRegistrationTypeMsg(
               &de_reg_request->m5gs_de_reg_type, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           de_reg_request->nas_key_set_identifier.EncodeNASKeySetIdentifierMsg(
               &de_reg_request->nas_key_set_identifier, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           de_reg_request->m5gs_mobile_identity.EncodeM5GSMobileIdentityMsg(
               &de_reg_request->m5gs_mobile_identity, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  return encoded;
};
}  // namespace magma5g
