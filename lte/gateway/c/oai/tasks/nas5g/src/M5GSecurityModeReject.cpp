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
#include "M5GSecurityModeReject.h"
#include "M5GCommonDefs.h"

namespace magma5g {
SecurityModeRejectMsg::SecurityModeRejectMsg(){};
SecurityModeRejectMsg::~SecurityModeRejectMsg(){};

// Decoding Security Mode Reject Message and its IEs
int SecurityModeRejectMsg::DecodeSecurityModeRejectMsg(
    SecurityModeRejectMsg* sec_mode_reject, uint8_t* buffer, uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, SECURITY_MODE_REJECT_MINIMUM_LENGTH, len);

  if ((decoded_result =
           sec_mode_reject->extended_protocol_discriminator
               .DecodeExtendedProtocolDiscriminatorMsg(
                   &sec_mode_reject->extended_protocol_discriminator, 0,
                   buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           sec_mode_reject->spare_half_octet.DecodeSpareHalfOctetMsg(
               &sec_mode_reject->spare_half_octet, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           sec_mode_reject->sec_header_type.DecodeSecurityHeaderTypeMsg(
               &sec_mode_reject->sec_header_type, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = sec_mode_reject->message_type.DecodeMessageTypeMsg(
           &sec_mode_reject->message_type, 0, buffer + decoded,
           len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = sec_mode_reject->m5gmm_cause.DecodeM5GMMCauseMsg(
           &sec_mode_reject->m5gmm_cause, 0, buffer + decoded, len - decoded)) <
      0)
    return decoded_result;
  else
    decoded += decoded_result;

  return decoded;
}

// Encoding Security Mode Reject Message and its IEs
int SecurityModeRejectMsg::EncodeSecurityModeRejectMsg(
    SecurityModeRejectMsg* sec_mode_reject, uint8_t* buffer, uint32_t len) {
  uint32_t encoded   = 0;
  int encoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, SECURITY_MODE_REJECT_MINIMUM_LENGTH, len);

  if ((encoded_result =
           sec_mode_reject->extended_protocol_discriminator
               .EncodeExtendedProtocolDiscriminatorMsg(
                   &sec_mode_reject->extended_protocol_discriminator, 0,
                   buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           sec_mode_reject->spare_half_octet.EncodeSpareHalfOctetMsg(
               &sec_mode_reject->spare_half_octet, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           sec_mode_reject->sec_header_type.EncodeSecurityHeaderTypeMsg(
               &sec_mode_reject->sec_header_type, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = sec_mode_reject->message_type.EncodeMessageTypeMsg(
           &sec_mode_reject->message_type, 0, buffer + encoded,
           len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = sec_mode_reject->m5gmm_cause.EncodeM5GMMCauseMsg(
           &sec_mode_reject->m5gmm_cause, 0, buffer + encoded, len - encoded)) <
      0)
    return encoded_result;
  else
    encoded += encoded_result;
  return encoded;
}
}  // namespace magma5g
