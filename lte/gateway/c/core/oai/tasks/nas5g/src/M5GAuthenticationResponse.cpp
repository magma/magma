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
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GAuthenticationResponse.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
AuthenticationResponseMsg::AuthenticationResponseMsg() {};
AuthenticationResponseMsg::~AuthenticationResponseMsg() {};

// Decode AuthenticationResponse Messsage
int AuthenticationResponseMsg::DecodeAuthenticationResponseMsg(
    AuthenticationResponseMsg* auth_response, uint8_t* buffer, uint32_t len) {
  uint32_t decoded = 0;
  int decoded_result = 0;

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, AUTHENTICATION_RESPONSE_MINIMUM_LENGTH, len);

  if ((decoded_result = auth_response->extended_protocol_discriminator
                            .DecodeExtendedProtocolDiscriminatorMsg(
                                &auth_response->extended_protocol_discriminator,
                                0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = auth_response->spare_half_octet.DecodeSpareHalfOctetMsg(
           &auth_response->spare_half_octet, 0, buffer + decoded,
           len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           auth_response->sec_header_type.DecodeSecurityHeaderTypeMsg(
               &auth_response->sec_header_type, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = auth_response->message_type.DecodeMessageTypeMsg(
           &auth_response->message_type, 0, buffer + decoded, len - decoded)) <
      0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = auth_response->autn_response_parameter
                            .DecodeAuthenticationResponseParameterMsg(
                                &auth_response->autn_response_parameter,
                                AUTH_RESPONSE_PARAMETER, buffer + decoded,
                                len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  return decoded;
};

// Encode AuthenticationResponse Messsage
int AuthenticationResponseMsg::EncodeAuthenticationResponseMsg(
    AuthenticationResponseMsg* auth_response, uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;
  // Not Implemented, Will be supported POST MVC
  return encoded;
};
}  // namespace magma5g
