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
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GAuthenticationRequest.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
AuthenticationRequestMsg::AuthenticationRequestMsg(){};
AuthenticationRequestMsg::~AuthenticationRequestMsg(){};

// Decode AuthenticationRequest Messsage
int AuthenticationRequestMsg::DecodeAuthenticationRequestMsg(uint8_t* buffer,
                                                             uint32_t len) {
  uint32_t decoded = 0;
  /*** Not Implemented, will be supported POST MVC ***/
  return decoded;
};

// Encode AuthenticationRequest Messsage
int AuthenticationRequestMsg::EncodeAuthenticationRequestMsg(uint8_t* buffer,
                                                             uint32_t len) {
  uint32_t encoded = 0;
  int encoded_result = 0;

  // Check if we got a NULL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, AUTHENTICATION_REQUEST_MINIMUM_LENGTH, len);

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
  if ((encoded_result =
           this->nas_key_set_identifier.EncodeNASKeySetIdentifierMsg(
               0, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           this->abba.EncodeABBAMsg(0, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = this->auth_rand.EncodeAuthenticationParameterRANDMsg(
           AUTH_PARAM_RAND, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;

  if ((encoded_result = this->auth_autn.EncodeAuthenticationParameterAUTNMsg(
           AUTH_PARAM_AUTN, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;

  return encoded;
};
}  // namespace magma5g
