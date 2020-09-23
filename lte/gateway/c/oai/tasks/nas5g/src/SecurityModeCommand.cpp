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
#include "SecurityModeCommand.h"
#include "CommonDefs.h"

namespace magma5g {
SecurityModeCommandMsg::SecurityModeCommandMsg(){};
SecurityModeCommandMsg::~SecurityModeCommandMsg(){};

// Decode SecurityModeCommand Message and its IEs
int SecurityModeCommandMsg::DecodeSecurityModeCommandMsg(
    SecurityModeCommandMsg* securitymodecommand, uint8_t* buffer,
    uint32_t len) {
  uint32_t decoded  = 0;
  int decodedresult = 0;

  // Checking Pointer
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, SECURITY_MODE_COMMAND_MINIMUM_LENGTH, len);

  MLOG(MDEBUG) << "DecodeSecurityModeCommandMsg : \n";
  if ((decodedresult =
           securitymodecommand->extendedprotocoldiscriminator
               .DecodeExtendedProtocolDiscriminatorMsg(
                   &securitymodecommand->extendedprotocoldiscriminator, 0,
                   buffer + decoded, len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult =
           securitymodecommand->sparehalfoctet.DecodeSpareHalfOctetMsg(
               &securitymodecommand->sparehalfoctet, 0, buffer + decoded,
               len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult =
           securitymodecommand->securityheadertype.DecodeSecurityHeaderTypeMsg(
               &securitymodecommand->securityheadertype, 0, buffer + decoded,
               len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult = securitymodecommand->messagetype.DecodeMessageTypeMsg(
           &securitymodecommand->messagetype, 0, buffer + decoded,
           len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult = securitymodecommand->nassecurityalgorithms
                           .DecodeNASSecurityAlgorithmsMsg(
                               &securitymodecommand->nassecurityalgorithms, 0,
                               buffer + decoded, len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult =
           securitymodecommand->sparehalfoctet.DecodeSpareHalfOctetMsg(
               &securitymodecommand->sparehalfoctet, 0, buffer + decoded,
               len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult = securitymodecommand->naskeysetidentifier
                           .DecodeNASKeySetIdentifierMsg(
                               &securitymodecommand->naskeysetidentifier, 0,
                               buffer + decoded, len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult = securitymodecommand->uesecuritycapability
                           .DecodeUESecurityCapabilityMsg(
                               &securitymodecommand->uesecuritycapability, 0,
                               buffer + decoded, len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;

  return decoded;
}

// Encode Security Mode Command Message and its IEs
int SecurityModeCommandMsg::EncodeSecurityModeCommandMsg(
    SecurityModeCommandMsg* securitymodecommand, uint8_t* buffer,
    uint32_t len) {
  uint32_t encoded = 0;

  MLOG(MDEBUG) << "EncodeSecurityModeCommandMsg:";
  int encodedresult = 0;

  // Check if we got a NULL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, SECURITY_MODE_COMMAND_MINIMUM_LENGTH, len);

  if ((encodedresult =
           securitymodecommand->extendedprotocoldiscriminator
               .EncodeExtendedProtocolDiscriminatorMsg(
                   &securitymodecommand->extendedprotocoldiscriminator, 0,
                   buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           securitymodecommand->sparehalfoctet.EncodeSpareHalfOctetMsg(
               &securitymodecommand->sparehalfoctet, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           securitymodecommand->securityheadertype.EncodeSecurityHeaderTypeMsg(
               &securitymodecommand->securityheadertype, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = securitymodecommand->messagetype.EncodeMessageTypeMsg(
           &securitymodecommand->messagetype, 0, buffer + encoded,
           len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = securitymodecommand->nassecurityalgorithms
                           .EncodeNASSecurityAlgorithmsMsg(
                               &securitymodecommand->nassecurityalgorithms, 0,
                               buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           securitymodecommand->sparehalfoctet.EncodeSpareHalfOctetMsg(
               &securitymodecommand->sparehalfoctet, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = securitymodecommand->naskeysetidentifier
                           .EncodeNASKeySetIdentifierMsg(
                               &securitymodecommand->naskeysetidentifier, 0,
                               buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = securitymodecommand->uesecuritycapability
                           .EncodeUESecurityCapabilityMsg(
                               &securitymodecommand->uesecuritycapability, 0,
                               buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;

  return encoded;
}
}  // namespace magma5g
