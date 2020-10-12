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
#include "SecurityModeReject.h"
#include "CommonDefs.h"

namespace magma5g {
SecurityModeRejectMsg::SecurityModeRejectMsg(){};
SecurityModeRejectMsg::~SecurityModeRejectMsg(){};

// Decoding Security Mode Reject Message and its IEs
int SecurityModeRejectMsg::DecodeSecurityModeRejectMsg(
    SecurityModeRejectMsg* securitymodereject, uint8_t* buffer, uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, SECURITY_MODE_REJECT_MINIMUM_LENGTH, len);

  if ((decoded_result =
           securitymodereject->extendedprotocoldiscriminator
               .DecodeExtendedProtocolDiscriminatorMsg(
                   &securitymodereject->extendedprotocoldiscriminator, 0,
                   buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           securitymodereject->sparehalfoctet.DecodeSpareHalfOctetMsg(
               &securitymodereject->sparehalfoctet, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           securitymodereject->securityheadertype.DecodeSecurityHeaderTypeMsg(
               &securitymodereject->securityheadertype, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = securitymodereject->messagetype.DecodeMessageTypeMsg(
           &securitymodereject->messagetype, 0, buffer + decoded,
           len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = securitymodereject->m5gmmcause
                            .DecodeM5GMMCauseMsg(
                                &securitymodereject->m5gmmcause, 0,
                                buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  return decoded;
}

// Encoding Security Mode Reject Message and its IEs
int SecurityModeRejectMsg::EncodeSecurityModeRejectMsg(
    SecurityModeRejectMsg* securitymodereject, uint8_t* buffer, uint32_t len) {
  uint32_t encoded  = 0;
  int encodedresult = 0;
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, SECURITY_MODE_REJECT_MINIMUM_LENGTH, len);

  if ((encodedresult =
           securitymodereject->extendedprotocoldiscriminator
               .EncodeExtendedProtocolDiscriminatorMsg(
                   &securitymodereject->extendedprotocoldiscriminator, 0,
                   buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           securitymodereject->sparehalfoctet.EncodeSpareHalfOctetMsg(
               &securitymodereject->sparehalfoctet, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           securitymodereject->securityheadertype.EncodeSecurityHeaderTypeMsg(
               &securitymodereject->securityheadertype, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = securitymodereject->messagetype.EncodeMessageTypeMsg(
           &securitymodereject->messagetype, 0, buffer + encoded,
           len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = securitymodereject->m5gmmcause
                           .EncodeM5GMMCauseMsg(
                               &securitymodereject->m5gmmcause, 0,
                               buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  return encoded;
}
}  // namespace magma5g
