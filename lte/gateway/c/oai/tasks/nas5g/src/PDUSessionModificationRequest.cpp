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
#include "PDUSessionModificationRequest.h"
#include "CommonDefs.h"

namespace magma5g {
PDUSessionModificationRequestMsg::PDUSessionModificationRequestMsg(){};
PDUSessionModificationRequestMsg::~PDUSessionModificationRequestMsg(){};

// Decode PDUSessionModificationRequest Message and its IEs
int PDUSessionModificationRequestMsg::DecodePDUSessionModificationRequestMsg(
    PDUSessionModificationRequestMsg* pdusessionmodificationrequest,
    uint8_t* buffer, uint32_t len) {
  uint32_t decoded  = 0;
  int decodedresult = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, PDU_SESSION_MODIFICATION_REQ_MIN_LEN, len);

  MLOG(MDEBUG) << "DecodePDUSessionModificationRequestMsg\n";
  if ((decodedresult =
           pdusessionmodificationrequest->extendedprotocoldiscriminator
               .DecodeExtendedProtocolDiscriminatorMsg(
                   &pdusessionmodificationrequest
                        ->extendedprotocoldiscriminator,
                   0, buffer + decoded, len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult =
           pdusessionmodificationrequest->pdusessionidentity
               .DecodePDUSessionIdentityMsg(
                   &pdusessionmodificationrequest->pdusessionidentity, 0,
                   buffer + decoded, len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult = pdusessionmodificationrequest->pti.DecodePTIMsg(
           &pdusessionmodificationrequest->pti, 0, buffer + decoded,
           len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult =
           pdusessionmodificationrequest->messagetype.DecodeMessageTypeMsg(
               &pdusessionmodificationrequest->messagetype, 0, buffer + decoded,
               len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;

  return 0;
}

// Encode PDUSessionModificationRequest Message and its IEs
int PDUSessionModificationRequestMsg::EncodePDUSessionModificationRequestMsg(
    PDUSessionModificationRequestMsg* pdusessionmodificationrequest,
    uint8_t* buffer, uint32_t len) {
  MLOG(MDEBUG) << "EncodePDUSessionModificationRequestMsg\n";
  // Not yet implemented, Will be supported POST MVC.
  return 0;
}

}  // namespace magma5g
