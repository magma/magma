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
#include "PDUSessionEstablishmentReject.h"
#include "CommonDefs.h"

namespace magma5g {
PDUSessionEstablishmentRejectMsg::PDUSessionEstablishmentRejectMsg(){};
PDUSessionEstablishmentRejectMsg::~PDUSessionEstablishmentRejectMsg(){};

// Decode PDUSessionEstablishmentReject Message and its IEs
int PDUSessionEstablishmentRejectMsg::DecodePDUSessionEstablishmentRejectMsg(
    PDUSessionEstablishmentRejectMsg* pdusessionestablishmentreject,
    uint8_t* buffer, uint32_t len) {
  // Not yet Implemented, will be supported POST MVC
  return 0;
}

// Encode PDUSessionEstablishmentReject Message and its IEs
int PDUSessionEstablishmentRejectMsg::EncodePDUSessionEstablishmentRejectMsg(
    PDUSessionEstablishmentRejectMsg* pdusessionestablishmentreject,
    uint8_t* buffer, uint32_t len) {
  uint32_t encoded  = 0;
  int encodedresult = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, PDU_SESSION_ESTABLISHMENT_REJ_MIN_LEN, len);

  MLOG(MDEBUG) << "EncodePDUSessionEstablishmentRejectMsg : \n";
  if ((encodedresult =
           pdusessionestablishmentreject->extendedprotocoldiscriminator
               .EncodeExtendedProtocolDiscriminatorMsg(
                   &pdusessionestablishmentreject
                        ->extendedprotocoldiscriminator,
                   0, buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           pdusessionestablishmentreject->pdusessionidentity
               .EncodePDUSessionIdentityMsg(
                   &pdusessionestablishmentreject->pdusessionidentity, 0,
                   buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = pdusessionestablishmentreject->pti.EncodePTIMsg(
           &pdusessionestablishmentreject->pti, 0, buffer + encoded,
           len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           pdusessionestablishmentreject->messagetype.EncodeMessageTypeMsg(
               &pdusessionestablishmentreject->messagetype, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           pdusessionestablishmentreject->m5gsmcause.EncodeM5GSMCauseMsg(
               &pdusessionestablishmentreject->m5gsmcause, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  return encoded;
}
}  // Namespace magma5g
