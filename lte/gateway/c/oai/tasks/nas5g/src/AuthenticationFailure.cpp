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
#include "AuthenticationFailure.h"
#include "CommonDefs.h"

namespace magma5g {
AuthenticationFailureMsg::AuthenticationFailureMsg(){};
AuthenticationFailureMsg::~AuthenticationFailureMsg(){};

// Decoding Authentication Failure Message and its IEs
int AuthenticationFailureMsg::DecodeAuthenticationFailureMsg(
    AuthenticationFailureMsg* authenticationfailure, uint8_t* buffer, uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, AUTHENTICATION_FAILURE_MINIMUM_LENGTH, len);

  if ((decoded_result =
           authenticationfailure->extendedprotocoldiscriminator
               .DecodeExtendedProtocolDiscriminatorMsg(
                   &authenticationfailure->extendedprotocoldiscriminator, 0,
                   buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           authenticationfailure->sparehalfoctet.DecodeSpareHalfOctetMsg(
               &authenticationfailure->sparehalfoctet, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           authenticationfailure->securityheadertype.DecodeSecurityHeaderTypeMsg(
               &authenticationfailure->securityheadertype, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = authenticationfailure->messagetype.DecodeMessageTypeMsg(
           &authenticationfailure->messagetype, 0, buffer + decoded,
           len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = authenticationfailure->m5gmmcause
                            .DecodeM5GMMCauseMsg(
                                &authenticationfailure->m5gmmcause, 0,
                                buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  return decoded;
}

// Encoding Authentication Failure Message and its IEs
int AuthenticationFailureMsg::EncodeAuthenticationFailureMsg(
    AuthenticationFailureMsg* authenticationfailure, uint8_t* buffer, uint32_t len) {
  uint32_t encoded  = 0;
  int encodedresult = 0;
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, AUTHENTICATION_FAILURE_MINIMUM_LENGTH, len);

  if ((encodedresult =
           authenticationfailure->extendedprotocoldiscriminator
               .EncodeExtendedProtocolDiscriminatorMsg(
                   &authenticationfailure->extendedprotocoldiscriminator, 0,
                   buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           authenticationfailure->sparehalfoctet.EncodeSpareHalfOctetMsg(
               &authenticationfailure->sparehalfoctet, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           authenticationfailure->securityheadertype.EncodeSecurityHeaderTypeMsg(
               &authenticationfailure->securityheadertype, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = authenticationfailure->messagetype.EncodeMessageTypeMsg(
           &authenticationfailure->messagetype, 0, buffer + encoded,
           len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = authenticationfailure->m5gmmcause
                           .EncodeM5GMMCauseMsg(
                               &authenticationfailure->m5gmmcause, 0,
                               buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  return encoded;
}
}  // namespace magma5g
