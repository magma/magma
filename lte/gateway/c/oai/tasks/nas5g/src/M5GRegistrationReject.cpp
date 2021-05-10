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
#include "M5GRegistrationReject.h"
#include "M5GCommonDefs.h"

namespace magma5g {
RegistrationRejectMsg::RegistrationRejectMsg(){};
RegistrationRejectMsg::~RegistrationRejectMsg(){};

// Decoding Registration Reject Message and its IEs
int RegistrationRejectMsg::DecodeRegistrationRejectMsg(
    RegistrationRejectMsg* reg_reject, uint8_t* buffer, uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, REGISTRATION_REJECT_MINIMUM_LENGTH, len);

  if ((decoded_result = reg_reject->extended_protocol_discriminator
                            .DecodeExtendedProtocolDiscriminatorMsg(
                                &reg_reject->extended_protocol_discriminator, 0,
                                buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = reg_reject->spare_half_octet.DecodeSpareHalfOctetMsg(
           &reg_reject->spare_half_octet, 0, buffer + decoded, len - decoded)) <
      0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = reg_reject->sec_header_type.DecodeSecurityHeaderTypeMsg(
           &reg_reject->sec_header_type, 0, buffer + decoded, len - decoded)) <
      0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = reg_reject->message_type.DecodeMessageTypeMsg(
           &reg_reject->message_type, 0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = reg_reject->m5gmm_cause.DecodeM5GMMCauseMsg(
           &reg_reject->m5gmm_cause, 0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  return decoded;
}

// Encoding Registration Reject Message and its IEs
int RegistrationRejectMsg::EncodeRegistrationRejectMsg(
    RegistrationRejectMsg* reg_reject, uint8_t* buffer, uint32_t len) {
  uint32_t encoded   = 0;
  int encoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, REGISTRATION_REJECT_MINIMUM_LENGTH, len);

  if ((encoded_result = reg_reject->extended_protocol_discriminator
                            .EncodeExtendedProtocolDiscriminatorMsg(
                                &reg_reject->extended_protocol_discriminator, 0,
                                buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = reg_reject->spare_half_octet.EncodeSpareHalfOctetMsg(
           &reg_reject->spare_half_octet, 0, buffer + encoded, len - encoded)) <
      0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = reg_reject->sec_header_type.EncodeSecurityHeaderTypeMsg(
           &reg_reject->sec_header_type, 0, buffer + encoded, len - encoded)) <
      0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = reg_reject->message_type.EncodeMessageTypeMsg(
           &reg_reject->message_type, 0, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = reg_reject->m5gmm_cause.EncodeM5GMMCauseMsg(
           &reg_reject->m5gmm_cause, 0, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  return encoded;
}
}  // namespace magma5g
