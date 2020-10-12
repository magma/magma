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
#include "AuthenticationResult.h"
#include "CommonDefs.h"

using namespace std;
namespace magma5g {
AuthenticationResultMsg::AuthenticationResultMsg(){};

AuthenticationResultMsg::~AuthenticationResultMsg(){};

// Decode Authentication Result Message and its IEs
int AuthenticationResultMsg::DecodeAuthenticationResultMsg(
    AuthenticationResultMsg* authenticationresult, uint8_t* buffer,
    uint32_t len) {
  uint32_t decoded  = 0;
  int decodedresult = 0;

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, AUTHENTICATION_RESULT_MINIMUM_LENGTH, len);

  MLOG(MDEBUG) << "\n\n---Decoding Authentication Result Message---\n"
               << endl;
  if ((decodedresult =
           authenticationresult->extendedprotocoldiscriminator
               .DecodeExtendedProtocolDiscriminatorMsg(
                   &authenticationresult->extendedprotocoldiscriminator, 0,
                   buffer + decoded, len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult =
           authenticationresult->sparehalfoctet.DecodeSpareHalfOctetMsg(
               &authenticationresult->sparehalfoctet, 0, buffer + decoded,
               len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult = authenticationresult->securityheadertype
                           .DecodeSecurityHeaderTypeMsg(
                               &authenticationresult->securityheadertype, 0,
                               buffer + decoded, len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult = authenticationresult->messagetype.DecodeMessageTypeMsg(
           &authenticationresult->messagetype, 0, buffer + decoded,
           len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult =
           authenticationresult->sparehalfoctet.DecodeSpareHalfOctetMsg(
               &authenticationresult->sparehalfoctet, 0, buffer + decoded,
               len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult = authenticationresult->naskeysetidentifier
                           .DecodeNASKeySetIdentifierMsg(
                               &authenticationresult->naskeysetidentifier, 0,
                               buffer + decoded, len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult = authenticationresult->eapmessage
                           .DecodeEAPMessageMsg(
                               &authenticationresult->eapmessage, 0,
                               buffer + decoded, len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  return decoded;
};

// Encode Authentication Result Message and its IEs
int AuthenticationResultMsg::EncodeAuthenticationResultMsg(
    AuthenticationResultMsg* authenticationresult, uint8_t* buffer,
    uint32_t len) {
  uint32_t encoded  = 0;
  int encodedresult = 0;

  // Check if we got a NULL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, AUTHENTICATION_RESULT_MINIMUM_LENGTH, len);

  if ((encodedresult =
           authenticationresult->extendedprotocoldiscriminator
               .EncodeExtendedProtocolDiscriminatorMsg(
                   &authenticationresult->extendedprotocoldiscriminator, 0,
                   buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           authenticationresult->sparehalfoctet.EncodeSpareHalfOctetMsg(
               &authenticationresult->sparehalfoctet, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = authenticationresult->securityheadertype
                           .EncodeSecurityHeaderTypeMsg(
                               &authenticationresult->securityheadertype, 0,
                               buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = authenticationresult->messagetype.EncodeMessageTypeMsg(
           &authenticationresult->messagetype, 0, buffer + encoded,
           len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           authenticationresult->sparehalfoctet.EncodeSpareHalfOctetMsg(
               &authenticationresult->sparehalfoctet, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  if ((encodedresult = authenticationresult->naskeysetidentifier
                           .EncodeNASKeySetIdentifierMsg(
                               &authenticationresult->naskeysetidentifier, 0,
                               buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = authenticationresult->eapmessage
                           .EncodeEAPMessageMsg(
                               &authenticationresult->eapmessage, 0,
                               buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  return encoded;
};
}  // namespace magma5g
