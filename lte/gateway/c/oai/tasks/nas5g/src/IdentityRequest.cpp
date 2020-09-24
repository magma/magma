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
#include "IdentityRequest.h"
#include "CommonDefs.h"

namespace magma5g {
IdentityRequestMsg::IdentityRequestMsg(){};
IdentityRequestMsg::~IdentityRequestMsg(){};

// Decode IdentityRequest Message and its IEs
int IdentityRequestMsg::DecodeIdentityRequestMsg(
    IdentityRequestMsg* identityrequest, uint8_t* buffer, uint32_t len) {
  uint32_t decoded  = 0;
  int decodedresult = 0;

  // Checking Pointer
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, IDENTITY_REQUEST_MINIMUM_LENGTH, len);

  MLOG(MDEBUG) << "DecodeIdentityRequestMsg : \n";
  if ((decodedresult = identityrequest->extendedprotocoldiscriminator
                           .DecodeExtendedProtocolDiscriminatorMsg(
                               &identityrequest->extendedprotocoldiscriminator,
                               0, buffer + decoded, len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult = identityrequest->sparehalfoctet.DecodeSpareHalfOctetMsg(
           &identityrequest->sparehalfoctet, 0, buffer + decoded,
           len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult =
           identityrequest->securityheadertype.DecodeSecurityHeaderTypeMsg(
               &identityrequest->securityheadertype, 0, buffer + decoded,
               len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult = identityrequest->messagetype.DecodeMessageTypeMsg(
           &identityrequest->messagetype, 0, buffer + decoded, len - decoded)) <
      0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult = identityrequest->sparehalfoctet.DecodeSpareHalfOctetMsg(
           &identityrequest->sparehalfoctet, 0, buffer + decoded,
           len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult =
           identityrequest->m5gsidentitytype.DecodeM5GSIdentityTypeMsg(
               &identityrequest->m5gsidentitytype, 0, buffer + decoded,
               len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;

  return decoded;
}

// Encode Identity Request Message and its IEs
int IdentityRequestMsg::EncodeIdentityRequestMsg(
    IdentityRequestMsg* identityrequest, uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;

  MLOG(MDEBUG) << "EncodeIdentityRequestMsg:";
  int encodedresult = 0;

  // Check if we got a NULL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, IDENTITY_REQUEST_MINIMUM_LENGTH, len);

  if ((encodedresult = identityrequest->extendedprotocoldiscriminator
                           .EncodeExtendedProtocolDiscriminatorMsg(
                               &identityrequest->extendedprotocoldiscriminator,
                               0, buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = identityrequest->sparehalfoctet.EncodeSpareHalfOctetMsg(
           &identityrequest->sparehalfoctet, 0, buffer + encoded,
           len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           identityrequest->securityheadertype.EncodeSecurityHeaderTypeMsg(
               &identityrequest->securityheadertype, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = identityrequest->messagetype.EncodeMessageTypeMsg(
           &identityrequest->messagetype, 0, buffer + encoded, len - encoded)) <
      0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = identityrequest->sparehalfoctet.EncodeSpareHalfOctetMsg(
           &identityrequest->sparehalfoctet, 0, buffer + encoded,
           len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           identityrequest->m5gsidentitytype.EncodeM5GSIdentityTypeMsg(
               &identityrequest->m5gsidentitytype, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;

  return encoded;
}
}  // namespace magma5g
