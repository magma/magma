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
#include "M5GRegistrationComplete.h"
#include "M5GCommonDefs.h"

namespace magma5g {
RegistrationCompleteMsg::RegistrationCompleteMsg(){};
RegistrationCompleteMsg::~RegistrationCompleteMsg(){};

// Decoding Registration Complete Message and its IEs
int RegistrationCompleteMsg::DecodeRegistrationCompleteMsg(
    RegistrationCompleteMsg* reg_complete, uint8_t* buffer, uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, REGISTRATION_COMPLETE_MINIMUM_LENGTH, len);

  if ((decoded_result = reg_complete->extended_protocol_discriminator
                            .DecodeExtendedProtocolDiscriminatorMsg(
                                &reg_complete->extended_protocol_discriminator,
                                0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = reg_complete->spare_half_octet.DecodeSpareHalfOctetMsg(
           &reg_complete->spare_half_octet, 0, buffer + decoded,
           len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           reg_complete->sec_header_type.DecodeSecurityHeaderTypeMsg(
               &reg_complete->sec_header_type, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = reg_complete->message_type.DecodeMessageTypeMsg(
           &reg_complete->message_type, 0, buffer + decoded, len - decoded)) <
      0)
    return decoded_result;
  else
    decoded += decoded_result;

  return decoded;
}

// Encoding Registration Complete Message and its IEs
int RegistrationCompleteMsg::EncodeRegistrationCompleteMsg(
    RegistrationCompleteMsg* reg_complete, uint8_t* buffer, uint32_t len) {
  uint32_t encoded   = 0;
  int encoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, REGISTRATION_COMPLETE_MINIMUM_LENGTH, len);

  if ((encoded_result = reg_complete->extended_protocol_discriminator
                            .EncodeExtendedProtocolDiscriminatorMsg(
                                &reg_complete->extended_protocol_discriminator,
                                0, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = reg_complete->spare_half_octet.EncodeSpareHalfOctetMsg(
           &reg_complete->spare_half_octet, 0, buffer + encoded,
           len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           reg_complete->sec_header_type.EncodeSecurityHeaderTypeMsg(
               &reg_complete->sec_header_type, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = reg_complete->message_type.EncodeMessageTypeMsg(
           &reg_complete->message_type, 0, buffer + encoded, len - encoded)) <
      0)
    return encoded_result;
  else
    encoded += encoded_result;

  return encoded;
}
}  // namespace magma5g
