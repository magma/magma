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
#include "M5GSecurityModeCommand.h"
#include "M5GCommonDefs.h"

namespace magma5g {
SecurityModeCommandMsg::SecurityModeCommandMsg(){};
SecurityModeCommandMsg::~SecurityModeCommandMsg(){};

// Decode SecurityModeCommand Message and its IEs
int SecurityModeCommandMsg::DecodeSecurityModeCommandMsg(
    SecurityModeCommandMsg* sec_mode_command, uint8_t* buffer, uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;

  // Checking Pointer
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, SECURITY_MODE_COMMAND_MINIMUM_LENGTH, len);

  MLOG(MDEBUG) << "DecodeSecurityModeCommandMsg : \n";
  if ((decoded_result =
           sec_mode_command->extended_protocol_discriminator
               .DecodeExtendedProtocolDiscriminatorMsg(
                   &sec_mode_command->extended_protocol_discriminator, 0,
                   buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           sec_mode_command->spare_half_octet.DecodeSpareHalfOctetMsg(
               &sec_mode_command->spare_half_octet, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           sec_mode_command->sec_header_type.DecodeSecurityHeaderTypeMsg(
               &sec_mode_command->sec_header_type, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = sec_mode_command->message_type.DecodeMessageTypeMsg(
           &sec_mode_command->message_type, 0, buffer + decoded,
           len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           sec_mode_command->nas_sec_algorithms.DecodeNASSecurityAlgorithmsMsg(
               &sec_mode_command->nas_sec_algorithms, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           sec_mode_command->spare_half_octet.DecodeSpareHalfOctetMsg(
               &sec_mode_command->spare_half_octet, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = sec_mode_command->nas_key_set_identifier
                            .DecodeNASKeySetIdentifierMsg(
                                &sec_mode_command->nas_key_set_identifier, 0,
                                buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           sec_mode_command->ue_sec_capability.DecodeUESecurityCapabilityMsg(
               &sec_mode_command->ue_sec_capability, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  return decoded;
}

// Encode Security Mode Command Message and its IEs
int SecurityModeCommandMsg::EncodeSecurityModeCommandMsg(
    SecurityModeCommandMsg* sec_mode_command, uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;

  MLOG(MDEBUG) << "EncodeSecurityModeCommandMsg:";
  int encoded_result = 0;

  // Check if we got a NULL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, SECURITY_MODE_COMMAND_MINIMUM_LENGTH, len);

  if ((encoded_result =
           sec_mode_command->extended_protocol_discriminator
               .EncodeExtendedProtocolDiscriminatorMsg(
                   &sec_mode_command->extended_protocol_discriminator, 0,
                   buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           sec_mode_command->spare_half_octet.EncodeSpareHalfOctetMsg(
               &sec_mode_command->spare_half_octet, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           sec_mode_command->sec_header_type.EncodeSecurityHeaderTypeMsg(
               &sec_mode_command->sec_header_type, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = sec_mode_command->message_type.EncodeMessageTypeMsg(
           &sec_mode_command->message_type, 0, buffer + encoded,
           len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           sec_mode_command->nas_sec_algorithms.EncodeNASSecurityAlgorithmsMsg(
               &sec_mode_command->nas_sec_algorithms, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           sec_mode_command->spare_half_octet.EncodeSpareHalfOctetMsg(
               &sec_mode_command->spare_half_octet, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = sec_mode_command->nas_key_set_identifier
                            .EncodeNASKeySetIdentifierMsg(
                                &sec_mode_command->nas_key_set_identifier, 0,
                                buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           sec_mode_command->ue_sec_capability.EncodeUESecurityCapabilityMsg(
               &sec_mode_command->ue_sec_capability, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = sec_mode_command->imeisv_request.EncodeImeisvRequestMsg(
           &sec_mode_command->imeisv_request, 0, buffer + encoded,
           len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;

  return encoded;
}
}  // namespace magma5g
