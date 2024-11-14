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
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GRegistrationAccept.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5gNasMessage.h"

namespace magma5g {
RegistrationAcceptMsg::RegistrationAcceptMsg(){};
RegistrationAcceptMsg::~RegistrationAcceptMsg(){};

// Decoding Registration Accept Message and its IEs
int RegistrationAcceptMsg::DecodeRegistrationAcceptMsg(uint8_t* buffer,
                                                       uint32_t len) {
  uint32_t decoded = 0;
  int decoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(buffer,
                                       REGISTRATION_ACCEPT_MINIMUM_LENGTH, len);

  if ((decoded_result = this->extended_protocol_discriminator
                            .DecodeExtendedProtocolDiscriminatorMsg(
                                0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = this->spare_half_octet.DecodeSpareHalfOctetMsg(
           0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = this->sec_header_type.DecodeSecurityHeaderTypeMsg(
           0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = this->message_type.DecodeMessageTypeMsg(
           0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = this->m5gs_reg_result.DecodeM5GSRegistrationResultMsg(
           0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  return decoded;
}

// Encoding Registration Accept Message and its IEs
int RegistrationAcceptMsg::EncodeRegistrationAcceptMsg(uint8_t* buffer,
                                                       uint32_t len) {
  uint32_t encoded = 0;
  int encoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer,
                                       REGISTRATION_ACCEPT_MINIMUM_LENGTH, len);

  if ((encoded_result = this->extended_protocol_discriminator
                            .EncodeExtendedProtocolDiscriminatorMsg(
                                0, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = this->spare_half_octet.EncodeSpareHalfOctetMsg(
           0, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = this->sec_header_type.EncodeSecurityHeaderTypeMsg(
           0, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = this->message_type.EncodeMessageTypeMsg(
           0, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = this->m5gs_reg_result.EncodeM5GSRegistrationResultMsg(
           0, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = this->mobile_id.EncodeM5GSMobileIdentityMsg(
           0x77, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = this->tai_list.EncodeTAIListMsg(0x54, buffer + encoded,
                                                        len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = this->allowed_nssai.EncodeNSSAIMsgList(
           ALLOWED_NSSAI, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = this->network_feature.EncodeNetworkFeatureSupportMsg(
           0x21, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = this->gprs_timer.EncodeGPRSTimer3Msg(
           0x5E, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;

  return encoded;
}
}  // namespace magma5g
