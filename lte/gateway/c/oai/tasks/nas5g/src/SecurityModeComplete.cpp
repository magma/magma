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
#include "SecurityModeComplete.h"
#include "CommonDefs.h"

namespace magma5g {
SecurityModeCompleteMsg::SecurityModeCompleteMsg(){};
SecurityModeCompleteMsg::~SecurityModeCompleteMsg(){};

// Decode SecurityModeComplete Message and its IEs
int SecurityModeCompleteMsg::DecodeSecurityModeCompleteMsg(
    SecurityModeCompleteMsg* securitymodecomplete, uint8_t* buffer,
    uint32_t len) {
  uint32_t decoded  = 0;
  int decodedresult = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, SECURITY_MODE_COMPLETE_MINIMUM_LENGTH, len);

  MLOG(MDEBUG) << "DecodeSecurityModeCompleteMsg : \n";
  if ((decodedresult =
           securitymodecomplete->extendedprotocoldiscriminator
               .DecodeExtendedProtocolDiscriminatorMsg(
                   &securitymodecomplete->extendedprotocoldiscriminator, 0,
                   buffer + decoded, len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult =
           securitymodecomplete->sparehalfoctet.DecodeSpareHalfOctetMsg(
               &securitymodecomplete->sparehalfoctet, 0, buffer + decoded,
               len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult =
           securitymodecomplete->securityheadertype.DecodeSecurityHeaderTypeMsg(
               &securitymodecomplete->securityheadertype, 0, buffer + decoded,
               len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult = securitymodecomplete->messagetype.DecodeMessageTypeMsg(
           &securitymodecomplete->messagetype, 0, buffer + decoded,
           len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  return decoded;
}

// Will be supported POST MVC
// Encode Security Mode Complete Message and its IEs
int SecurityModeCompleteMsg::EncodeSecurityModeCompleteMsg(
    SecurityModeCompleteMsg* securitymodecomplete, uint8_t* buffer,
    uint32_t len) {
  uint32_t encoded = 0;

#ifdef HANDLE_POST_MVC
  MLOG(MDEBUG) << "EncodeSecurityModeCompleteMsg:";
  int encodedresult = 0;

  // Check if we got a NULL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, SECURITY_MODE_COMPLETE_MINIMUM_LENGTH, len);

  if ((encodedresult =
           securitymodecomplete->EncodeExtendedProtocolDiscriminatorMsg(
               securitymodecomplete->extendedprotocoldiscriminator, 0,
               buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = securitymodecomplete->EncodeSecurityHeaderTypeMsg(
           securitymodecomplete->securityheadertype, 0, buffer + encoded,
           len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = securitymodecomplete->EncodeMessageTypeMsg(
           securitymodecomplete->messagetype, 0, buffer + encoded,
           len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
#endif
  return encoded;
}
}  // namespace magma5g
