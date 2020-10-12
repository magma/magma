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
#include "RegistrationReject.h"
#include "CommonDefs.h"

namespace magma5g {
RegistrationRejectMsg::RegistrationRejectMsg(){};
RegistrationRejectMsg::~RegistrationRejectMsg(){};

// Decoding Registration Reject Message and its IEs
int RegistrationRejectMsg::DecodeRegistrationRejectMsg(
    RegistrationRejectMsg* registrationreject, uint8_t* buffer, uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, REGISTRATION_REJECT_MINIMUM_LENGTH, len);

  if ((decoded_result =
           registrationreject->extendedprotocoldiscriminator
               .DecodeExtendedProtocolDiscriminatorMsg(
                   &registrationreject->extendedprotocoldiscriminator, 0,
                   buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           registrationreject->sparehalfoctet.DecodeSpareHalfOctetMsg(
               &registrationreject->sparehalfoctet, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           registrationreject->securityheadertype.DecodeSecurityHeaderTypeMsg(
               &registrationreject->securityheadertype, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = registrationreject->messagetype.DecodeMessageTypeMsg(
           &registrationreject->messagetype, 0, buffer + decoded,
           len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = registrationreject->m5gmmcause
                            .DecodeM5GMMCauseMsg(
                                &registrationreject->m5gmmcause, 0,
                                buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  return decoded;
}

// Encoding Registration Reject Message and its IEs
int RegistrationRejectMsg::EncodeRegistrationRejectMsg(
    RegistrationRejectMsg* registrationreject, uint8_t* buffer, uint32_t len) {
  uint32_t encoded  = 0;
  int encodedresult = 0;
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, REGISTRATION_REJECT_MINIMUM_LENGTH, len);

  if ((encodedresult =
           registrationreject->extendedprotocoldiscriminator
               .EncodeExtendedProtocolDiscriminatorMsg(
                   &registrationreject->extendedprotocoldiscriminator, 0,
                   buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           registrationreject->sparehalfoctet.EncodeSpareHalfOctetMsg(
               &registrationreject->sparehalfoctet, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           registrationreject->securityheadertype.EncodeSecurityHeaderTypeMsg(
               &registrationreject->securityheadertype, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = registrationreject->messagetype.EncodeMessageTypeMsg(
           &registrationreject->messagetype, 0, buffer + encoded,
           len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = registrationreject->m5gmmcause
                           .EncodeM5GMMCauseMsg(
                               &registrationreject->m5gmmcause, 0,
                               buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  return encoded;
}
}  // namespace magma5g
