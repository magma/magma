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
#include "PDUSessionReleaseRequest.h"
#include "CommonDefs.h"

namespace magma5g {
PDUSessionReleaseRequestMsg::PDUSessionReleaseRequestMsg(){};
PDUSessionReleaseRequestMsg::~PDUSessionReleaseRequestMsg(){};

// Decode PDUSessionReleaseRequest Message and its IEs
int PDUSessionReleaseRequestMsg::DecodePDUSessionReleaseRequestMsg(
    PDUSessionReleaseRequestMsg* pdusessionreleaserequest, uint8_t* buffer,
    uint32_t len) {
  uint32_t decoded  = 0;
  int decodedresult = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, PDU_SESSION_RELEASE_REQ_MIN_LEN, len);

  MLOG(MDEBUG) << "DecodePDUSessionReleaseRequestMsg\n";
  if ((decodedresult =
           pdusessionreleaserequest->extendedprotocoldiscriminator
               .DecodeExtendedProtocolDiscriminatorMsg(
                   &pdusessionreleaserequest->extendedprotocoldiscriminator, 0,
                   buffer + decoded, len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult = pdusessionreleaserequest->pdusessionidentity
                           .DecodePDUSessionIdentityMsg(
                               &pdusessionreleaserequest->pdusessionidentity, 0,
                               buffer + decoded, len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult = pdusessionreleaserequest->pti.DecodePTIMsg(
           &pdusessionreleaserequest->pti, 0, buffer + decoded,
           len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult =
           pdusessionreleaserequest->messagetype.DecodeMessageTypeMsg(
               &pdusessionreleaserequest->messagetype, 0, buffer + decoded,
               len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;

  return decoded;
}

// Encode PDUSessionReleaseRequest Message and its IEs
int PDUSessionReleaseRequestMsg::EncodePDUSessionReleaseRequestMsg(
    PDUSessionReleaseRequestMsg* pdusessionreleaserequest, uint8_t* buffer,
    uint32_t len) {
  MLOG(MDEBUG) << "EncodePDUSessionReleaseRequestMsg\n";
  //Not yet implemented, Will be supported POST MVC.
  return 0;
}
}  // namespace magma5g
