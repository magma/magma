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
#include "M5GAuthenticationRequest.h"
#include "M5GCommonDefs.h"

using namespace std;
namespace magma5g {
AuthenticationRequestMsg::AuthenticationRequestMsg(){};
AuthenticationRequestMsg::~AuthenticationRequestMsg(){};

// Decode AuthenticationRequest Messsage
int AuthenticationRequestMsg::DecodeAuthenticationRequestMsg(
    AuthenticationRequestMsg* auth_request, uint8_t* buffer, uint32_t len) {
  uint32_t decoded = 0;
  /*** Not Implemented, will be supported POST MVC ***/
  return decoded;
};

// Encode AuthenticationRequest Messsage
int AuthenticationRequestMsg::EncodeAuthenticationRequestMsg(
    AuthenticationRequestMsg* auth_request, uint8_t* buffer, uint32_t len) {
  uint32_t encoded   = 0;
  int encoded_result = 0;

  // Check if we got a NULL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, AUTHENTICATION_REQUEST_MINIMUM_LENGTH, len);

  if ((encoded_result = auth_request->extended_protocol_discriminator
                            .EncodeExtendedProtocolDiscriminatorMsg(
                                &auth_request->extended_protocol_discriminator,
                                0, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = auth_request->spare_half_octet.EncodeSpareHalfOctetMsg(
           &auth_request->spare_half_octet, 0, buffer + encoded,
           len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           auth_request->sec_header_type.EncodeSecurityHeaderTypeMsg(
               &auth_request->sec_header_type, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = auth_request->message_type.EncodeMessageTypeMsg(
           &auth_request->message_type, 0, buffer + encoded, len - encoded)) <
      0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           auth_request->nas_key_set_identifier.EncodeNASKeySetIdentifierMsg(
               &auth_request->nas_key_set_identifier, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = auth_request->abba.EncodeABBAMsg(
           &auth_request->abba, 0, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           auth_request->auth_rand.EncodeAuthenticationParameterRANDMsg(
               &auth_request->auth_rand, AUTH_PARAM_RAND, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;

  if ((encoded_result =
           auth_request->auth_autn.EncodeAuthenticationParameterAUTNMsg(
               &auth_request->auth_autn, AUTH_PARAM_AUTN, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;

  return encoded;
};
}  // namespace magma5g
