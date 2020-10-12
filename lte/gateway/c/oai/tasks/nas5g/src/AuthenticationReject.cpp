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
#include "AuthenticationReject.h"
#include "CommonDefs.h"

namespace magma5g {
AuthenticationRejectMsg::AuthenticationRejectMsg(){};
AuthenticationRejectMsg::~AuthenticationRejectMsg(){};

// Decoding Authentication Reject Message and its IEs
int AuthenticationRejectMsg::DecodeAuthenticationRejectMsg(
    AuthenticationRejectMsg* authenticationreject, uint8_t* buffer, uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, AUTHENTICATION_REJECT_MINIMUM_LENGTH, len);

  if ((decoded_result =
           authenticationreject->extendedprotocoldiscriminator
               .DecodeExtendedProtocolDiscriminatorMsg(
                   &authenticationreject->extendedprotocoldiscriminator, 0,
                   buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           authenticationreject->sparehalfoctet.DecodeSpareHalfOctetMsg(
               &authenticationreject->sparehalfoctet, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           authenticationreject->securityheadertype.DecodeSecurityHeaderTypeMsg(
               &authenticationreject->securityheadertype, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = authenticationreject->messagetype.DecodeMessageTypeMsg(
           &authenticationreject->messagetype, 0, buffer + decoded,
           len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  return decoded;
}

// Encoding Authentication Reject Message and its IEs
int AuthenticationRejectMsg::EncodeAuthenticationRejectMsg(
    AuthenticationRejectMsg* authenticationreject, uint8_t* buffer, uint32_t len) {
  uint32_t encoded  = 0;
  int encodedresult = 0;
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, AUTHENTICATION_REJECT_MINIMUM_LENGTH, len);

  if ((encodedresult =
           authenticationreject->extendedprotocoldiscriminator
               .EncodeExtendedProtocolDiscriminatorMsg(
                   &authenticationreject->extendedprotocoldiscriminator, 0,
                   buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           authenticationreject->sparehalfoctet.EncodeSpareHalfOctetMsg(
               &authenticationreject->sparehalfoctet, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           authenticationreject->securityheadertype.EncodeSecurityHeaderTypeMsg(
               &authenticationreject->securityheadertype, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = authenticationreject->messagetype.EncodeMessageTypeMsg(
           &authenticationreject->messagetype, 0, buffer + encoded,
           len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;

  return encoded;
}
}  // namespace magma5g
