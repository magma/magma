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
#include "PDUSessionModificationReject.h"
#include "CommonDefs.h"

namespace magma5g {
PDUSessionModificationRejectMsg::PDUSessionModificationRejectMsg(){};
PDUSessionModificationRejectMsg::~PDUSessionModificationRejectMsg(){};

// Decode PDUSessionModificationReject Message and its IEs
int PDUSessionModificationRejectMsg::DecodePDUSessionModificationRejectMsg(
    PDUSessionModificationRejectMsg* pdusessionmodificationreject,
    uint8_t* buffer, uint32_t len) {
  uint32_t decoded  = 0;
  //Not yet Implemented, will be supported POST MVC.
  return decoded;
}

// Encode PDUSessionModificationReject Message and its IEs
int PDUSessionModificationRejectMsg::EncodePDUSessionModificationRejectMsg(
    PDUSessionModificationRejectMsg* pdusessionmodificationreject,
    uint8_t* buffer, uint32_t len) {
  uint32_t encoded  = 0;
  int encodedresult = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, PDU_SESSION_MODIFICATION_REJECT_MIN_LEN, len);

  MLOG(MDEBUG) << "EncodePDUSessionModificationRejectMsg : \n";
  if ((encodedresult =
           pdusessionmodificationreject->extendedprotocoldiscriminator
               .EncodeExtendedProtocolDiscriminatorMsg(
                   &pdusessionmodificationreject->extendedprotocoldiscriminator,
                   0, buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           pdusessionmodificationreject->pdusessionidentity
               .EncodePDUSessionIdentityMsg(
                   &pdusessionmodificationreject->pdusessionidentity, 0,
                   buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = pdusessionmodificationreject->pti.EncodePTIMsg(
           &pdusessionmodificationreject->pti, 0, buffer + encoded,
           len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           pdusessionmodificationreject->messagetype.EncodeMessageTypeMsg(
               &pdusessionmodificationreject->messagetype, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           pdusessionmodificationreject->m5gsmcause.EncodeM5GSMCauseMsg(
               &pdusessionmodificationreject->m5gsmcause, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  return encoded;
}
}  // namespace magma5g
