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
#include "IdentityResponse.h"
#include "CommonDefs.h"

namespace magma5g {
IdentityResponseMsg::IdentityResponseMsg(){};
IdentityResponseMsg::~IdentityResponseMsg(){};

// Decode IdentityResponse Message and its IEs
int IdentityResponseMsg::DecodeIdentityResponseMsg(
    IdentityResponseMsg* identityresponse, uint8_t* buffer, uint32_t len) {
  uint32_t decoded  = 0;
  int decodedresult = 0;

  // Checking Pointer
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, IDENTITY_RESPONSE_MINIMUM_LENGTH, len);

  MLOG(MDEBUG) << "DecodeIdentityResponseMsg : \n";
  if ((decodedresult = identityresponse->extendedprotocoldiscriminator
                           .DecodeExtendedProtocolDiscriminatorMsg(
                               &identityresponse->extendedprotocoldiscriminator,
                               0, buffer + decoded, len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult = identityresponse->sparehalfoctet.DecodeSpareHalfOctetMsg(
           &identityresponse->sparehalfoctet, 0, buffer + decoded,
           len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult =
           identityresponse->securityheadertype.DecodeSecurityHeaderTypeMsg(
               &identityresponse->securityheadertype, 0, buffer + decoded,
               len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult = identityresponse->messagetype.DecodeMessageTypeMsg(
           &identityresponse->messagetype, 0, buffer + decoded,
           len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult =
           identityresponse->m5gsmobileidentity.DecodeM5GSMobileIdentityMsg(
               &identityresponse->m5gsmobileidentity, 0, buffer + decoded,
               len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;

  return decoded;
}

// Encode Identity Response Message and its IEs
int IdentityResponseMsg::EncodeIdentityResponseMsg(
    IdentityResponseMsg* identityresponse, uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;

  MLOG(MDEBUG) << "EncodeIdentityResponseMsg:";
  int encodedresult = 0;

  // Check if we got a NULL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, IDENTITY_RESPONSE_MINIMUM_LENGTH, len);

  if ((encodedresult = identityresponse->extendedprotocoldiscriminator
                           .EncodeExtendedProtocolDiscriminatorMsg(
                               &identityresponse->extendedprotocoldiscriminator,
                               0, buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = identityresponse->sparehalfoctet.EncodeSpareHalfOctetMsg(
           &identityresponse->sparehalfoctet, 0, buffer + encoded,
           len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           identityresponse->securityheadertype.EncodeSecurityHeaderTypeMsg(
               &identityresponse->securityheadertype, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = identityresponse->messagetype.EncodeMessageTypeMsg(
           &identityresponse->messagetype, 0, buffer + encoded,
           len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           identityresponse->m5gsmobileidentity.EncodeM5GSMobileIdentityMsg(
               &identityresponse->m5gsmobileidentity, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;

  return encoded;
}
}  // namespace magma5g
