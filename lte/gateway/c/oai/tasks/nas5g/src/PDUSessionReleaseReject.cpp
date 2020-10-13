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
#include "PDUSessionReleaseReject.h"
#include "CommonDefs.h"

namespace magma5g {
PDUSessionReleaseRejectMsg::PDUSessionReleaseRejectMsg(){};
PDUSessionReleaseRejectMsg::~PDUSessionReleaseRejectMsg(){};

// Decode PDUSessionReleaseReject Message and its IEs
int PDUSessionReleaseRejectMsg::DecodePDUSessionReleaseRejectMsg(
    PDUSessionReleaseRejectMsg* pdusessionreleasereject, uint8_t* buffer,
    uint32_t len) {
  uint32_t decoded  = 0;
  //Not yet implemented, Will be supported POST MVC.
  return 0;
}

// Encode PDUSessionReleaseReject Message and its IEs
int PDUSessionReleaseRejectMsg::EncodePDUSessionReleaseRejectMsg(
    PDUSessionReleaseRejectMsg* pdusessionreleasereject, uint8_t* buffer,
    uint32_t len) {
  uint32_t encoded  = 0;
  int encodedresult = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, PDU_SESSION_RELEASE_REJECT_MIN_LEN, len);

  MLOG(MDEBUG) << "EncodePDUSessionReleaseRejectMsg : \n";
  if ((encodedresult =
           pdusessionreleasereject->extendedprotocoldiscriminator
               .EncodeExtendedProtocolDiscriminatorMsg(
                   &pdusessionreleasereject->extendedprotocoldiscriminator, 0,
                   buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = pdusessionreleasereject->pdusessionidentity
                           .EncodePDUSessionIdentityMsg(
                               &pdusessionreleasereject->pdusessionidentity, 0,
                               buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = pdusessionreleasereject->pti.EncodePTIMsg(
           &pdusessionreleasereject->pti, 0, buffer + encoded, len - encoded)) <
      0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           pdusessionreleasereject->messagetype.EncodeMessageTypeMsg(
               &pdusessionreleasereject->messagetype, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = pdusessionreleasereject->m5gsmcause.EncodeM5GSMCauseMsg(
           &pdusessionreleasereject->m5gsmcause, 0, buffer + encoded,
           len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  return encoded;
}
}  // namespace magma5g
