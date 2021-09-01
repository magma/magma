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
#include "M5GAuthenticationFailure.h"
#include "M5GCommonDefs.h"

namespace magma5g {
AuthenticationFailureMsg::AuthenticationFailureMsg(){};
AuthenticationFailureMsg::~AuthenticationFailureMsg(){};

// Decoding Authentication Failure Message and its IEs
int AuthenticationFailureMsg::DecodeAuthenticationFailureMsg(
    AuthenticationFailureMsg* auth_failure, uint8_t* buffer, uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, AUTHENTICATION_FAILURE_MINIMUM_LENGTH, len);

  if ((decoded_result = auth_failure->extended_protocol_discriminator
                            .DecodeExtendedProtocolDiscriminatorMsg(
                                &auth_failure->extended_protocol_discriminator,
                                0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = auth_failure->spare_half_octet.DecodeSpareHalfOctetMsg(
           &auth_failure->spare_half_octet, 0, buffer + decoded,
           len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           auth_failure->sec_header_type.DecodeSecurityHeaderTypeMsg(
               &auth_failure->sec_header_type, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = auth_failure->message_type.DecodeMessageTypeMsg(
           &auth_failure->message_type, 0, buffer + decoded, len - decoded)) <
      0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = auth_failure->m5gmm_cause.DecodeM5GMMCauseMsg(
           &auth_failure->m5gmm_cause, 0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  while (decoded < len) {
    uint8_t ieiDecoded = *(buffer + decoded);

    switch (ieiDecoded) {
      case AUTHENTICATION_FAILURE_PARAMETER_IEI_AUTH_CHALLENGE:
        if ((decoded_result =
                 auth_failure->auth_failure_ie.DecodeM5GAuthenticationFailureIE(
                     &auth_failure->auth_failure_ie,
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
int AuthenticationFailureMsg::EncodeAuthenticationFailureMsg(
    AuthenticationFailureMsg* auth_failure, uint8_t* buffer, uint32_t len) {
  uint32_t encoded  = 0;
  int encodedresult = 0;
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, AUTHENTICATION_FAILURE_MINIMUM_LENGTH, len);

  if ((encodedresult = auth_failure->extended_protocol_discriminator
                           .EncodeExtendedProtocolDiscriminatorMsg(
                               &auth_failure->extended_protocol_discriminator,
                               0, buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = auth_failure->spare_half_octet.EncodeSpareHalfOctetMsg(
           &auth_failure->spare_half_octet, 0, buffer + encoded,
           len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           auth_failure->sec_header_type.EncodeSecurityHeaderTypeMsg(
               &auth_failure->sec_header_type, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = auth_failure->message_type.EncodeMessageTypeMsg(
           &auth_failure->message_type, 0, buffer + encoded, len - encoded)) <
      0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = auth_failure->m5gmm_cause.EncodeM5GMMCauseMsg(
           &auth_failure->m5gmm_cause, 0, buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  return encoded;
}
}  // namespace magma5g
