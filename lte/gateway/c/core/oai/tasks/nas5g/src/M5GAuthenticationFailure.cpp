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
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GAuthenticationFailure.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
AuthenticationFailureMsg::AuthenticationFailureMsg(){};
AuthenticationFailureMsg::~AuthenticationFailureMsg(){};

// Decoding Authentication Failure Message and its IEs
int AuthenticationFailureMsg::DecodeAuthenticationFailureMsg(uint8_t* buffer,
                                                             uint32_t len) {
  uint32_t decoded = 0;
  int decoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, AUTHENTICATION_FAILURE_MINIMUM_LENGTH, len);

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
  if ((decoded_result = this->m5gmm_cause.DecodeM5GMMCauseMsg(
           0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  while (decoded < len) {
    uint8_t ieiDecoded = *(buffer + decoded);

    switch (ieiDecoded) {
      case AUTHENTICATION_FAILURE_PARAMETER_IEI_AUTH_CHALLENGE:
        if ((decoded_result =
                 this->auth_failure_ie.DecodeM5GAuthenticationFailureIE(
                     AUTHENTICATION_FAILURE_PARAMETER_IEI_AUTH_CHALLENGE,
                     buffer + decoded, len - decoded)) < 0)
          return decoded_result;

        decoded += decoded_result;
        break;

      default:
        return TLV_UNEXPECTED_IEI;
    }
  }

  return decoded;
}

// Encoding Authentication Failure Message and its IEs
int AuthenticationFailureMsg::EncodeAuthenticationFailureMsg(uint8_t* buffer,
                                                             uint32_t len) {
  uint32_t encoded = 0;
  int encodedresult = 0;
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, AUTHENTICATION_FAILURE_MINIMUM_LENGTH, len);

  if ((encodedresult = this->extended_protocol_discriminator
                           .EncodeExtendedProtocolDiscriminatorMsg(
                               0, buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = this->spare_half_octet.EncodeSpareHalfOctetMsg(
           0, buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = this->sec_header_type.EncodeSecurityHeaderTypeMsg(
           0, buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = this->message_type.EncodeMessageTypeMsg(
           0, buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = this->m5gmm_cause.EncodeM5GMMCauseMsg(
           0, buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  return encoded;
}
}  // namespace magma5g
