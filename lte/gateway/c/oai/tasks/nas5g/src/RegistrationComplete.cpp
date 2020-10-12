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
#include "RegistrationComplete.h"
#include "CommonDefs.h"

namespace magma5g {
RegistrationCompleteMsg::RegistrationCompleteMsg(){};
RegistrationCompleteMsg::~RegistrationCompleteMsg(){};

// Decoding Registration Complete Message and its IEs
int RegistrationCompleteMsg::DecodeRegistrationCompleteMsg(
    RegistrationCompleteMsg* registrationcomplete, uint8_t* buffer, uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, REGISTRATION_COMPLETE_MINIMUM_LENGTH, len);

  if ((decoded_result =
           registrationcomplete->extendedprotocoldiscriminator
               .DecodeExtendedProtocolDiscriminatorMsg(
                   &registrationcomplete->extendedprotocoldiscriminator, 0,
                   buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           registrationcomplete->sparehalfoctet.DecodeSpareHalfOctetMsg(
               &registrationcomplete->sparehalfoctet, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           registrationcomplete->securityheadertype.DecodeSecurityHeaderTypeMsg(
               &registrationcomplete->securityheadertype, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = registrationcomplete->messagetype.DecodeMessageTypeMsg(
           &registrationcomplete->messagetype, 0, buffer + decoded,
           len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  return decoded;
}

// Encoding Registration Complete Message and its IEs
int RegistrationCompleteMsg::EncodeRegistrationCompleteMsg(
    RegistrationCompleteMsg* registrationcomplete, uint8_t* buffer, uint32_t len) {
  uint32_t encoded  = 0;
  int encodedresult = 0;
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, REGISTRATION_COMPLETE_MINIMUM_LENGTH, len);

  if ((encodedresult =
           registrationcomplete->extendedprotocoldiscriminator
               .EncodeExtendedProtocolDiscriminatorMsg(
                   &registrationcomplete->extendedprotocoldiscriminator, 0,
                   buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           registrationcomplete->sparehalfoctet.EncodeSpareHalfOctetMsg(
               &registrationcomplete->sparehalfoctet, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           registrationcomplete->securityheadertype.EncodeSecurityHeaderTypeMsg(
               &registrationcomplete->securityheadertype, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = registrationcomplete->messagetype.EncodeMessageTypeMsg(
           &registrationcomplete->messagetype, 0, buffer + encoded,
           len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;

  return encoded;
}
}  // namespace magma5g
