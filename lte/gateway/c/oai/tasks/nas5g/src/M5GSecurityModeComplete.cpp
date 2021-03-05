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
#include "M5GSecurityModeComplete.h"
#include "M5GCommonDefs.h"

namespace magma5g {
SecurityModeCompleteMsg::SecurityModeCompleteMsg(){};
SecurityModeCompleteMsg::~SecurityModeCompleteMsg(){};

// Decode SecurityModeComplete Message and its IEs
int SecurityModeCompleteMsg::DecodeSecurityModeCompleteMsg(
    SecurityModeCompleteMsg* sec_mode_complete, uint8_t* buffer, uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, SECURITY_MODE_COMPLETE_MINIMUM_LENGTH, len);

  MLOG(MDEBUG) << "DecodeSecurityModeCompleteMsg : \n";
  if ((decoded_result =
           sec_mode_complete->extended_protocol_discriminator
               .DecodeExtendedProtocolDiscriminatorMsg(
                   &sec_mode_complete->extended_protocol_discriminator, 0,
                   buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           sec_mode_complete->spare_half_octet.DecodeSpareHalfOctetMsg(
               &sec_mode_complete->spare_half_octet, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           sec_mode_complete->sec_header_type.DecodeSecurityHeaderTypeMsg(
               &sec_mode_complete->sec_header_type, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = sec_mode_complete->message_type.DecodeMessageTypeMsg(
           &sec_mode_complete->message_type, 0, buffer + decoded,
           len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  return decoded;
}

// Will be supported POST MVC
// Encode Security Mode Complete Message and its IEs
int SecurityModeCompleteMsg::EncodeSecurityModeCompleteMsg(
    SecurityModeCompleteMsg* sec_mode_complete, uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;

#ifdef HANDLE_POST_MVC
  MLOG(MDEBUG) << "EncodeSecurityModeCompleteMsg:";
  int encoded_result = 0;

  // Check if we got a NULL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, SECURITY_MODE_COMPLETE_MINIMUM_LENGTH, len);

  if ((encoded_result =
           sec_mode_complete->EncodeExtendedProtocolDiscriminatorMsg(
               sec_mode_complete->extended_protocol_discriminator, 0,
               buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = sec_mode_complete->EncodeSecurityHeaderTypeMsg(
           sec_mode_complete->sec_header_type, 0, buffer + encoded,
           len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = sec_mode_complete->EncodeMessageTypeMsg(
           sec_mode_complete->message_type, 0, buffer + encoded,
           len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
#endif
  return encoded;
}
}  // namespace magma5g
