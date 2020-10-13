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
#include "PDUSessionEstablishmentRequest.h"
#include "CommonDefs.h"

namespace magma5g {
PDUSessionEstablishmentRequestMsg::PDUSessionEstablishmentRequestMsg(){};
PDUSessionEstablishmentRequestMsg::~PDUSessionEstablishmentRequestMsg(){};

// Decode PDUSessionEstablishmentRequest Message and its IEs
int PDUSessionEstablishmentRequestMsg::DecodePDUSessionEstablishmentRequestMsg(
    PDUSessionEstablishmentRequestMsg* pdusessionestablishmentrequest,
    uint8_t* buffer, uint32_t len) {
  uint32_t decoded  = 0;
  int decodedresult = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, PDU_SESSION_ESTABLISH_REQ_MIN_LEN, len);

  MLOG(MDEBUG) << "DecodePDUSessionEstablishmentRequestMsg : ";
  if ((decodedresult =
           pdusessionestablishmentrequest->extendedprotocoldiscriminator
               .DecodeExtendedProtocolDiscriminatorMsg(
                   &pdusessionestablishmentrequest
                        ->extendedprotocoldiscriminator,
                   0, buffer + decoded, len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult =
           pdusessionestablishmentrequest->pdusessionidentity
               .DecodePDUSessionIdentityMsg(
                   &pdusessionestablishmentrequest->pdusessionidentity, 0,
                   buffer + decoded, len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult = pdusessionestablishmentrequest->pti.DecodePTIMsg(
           &pdusessionestablishmentrequest->pti, 0, buffer + decoded,
           len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult =
           pdusessionestablishmentrequest->messagetype.DecodeMessageTypeMsg(
               &pdusessionestablishmentrequest->messagetype, 0,
               buffer + decoded, len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult =
           pdusessionestablishmentrequest->integrityprotmaxdatarate
               .DecodeIntegrityProtMaxDataRateMsg(
                   &pdusessionestablishmentrequest->integrityprotmaxdatarate, 0,
                   buffer + decoded, len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult =
           pdusessionestablishmentrequest->pdusessiontype
               .DecodePDUSessionTypeMsg(
                   &pdusessionestablishmentrequest->pdusessiontype,
                   PDUSESSIONTYPE, buffer + decoded, len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  if ((decodedresult = pdusessionestablishmentrequest->sscmode.DecodeSSCModeMsg(
           &pdusessionestablishmentrequest->sscmode, SSCMODE, buffer + decoded,
           len - decoded)) < 0)
    return decodedresult;
  else
    decoded += decodedresult;
  return decoded;
}

// Encode PDUSessionEstablishmentRequest Message and its IEs
int PDUSessionEstablishmentRequestMsg::EncodePDUSessionEstablishmentRequestMsg(
    PDUSessionEstablishmentRequestMsg* pdusessionestablishmentrequest,
    uint8_t* buffer, uint32_t len) {
  uint32_t Encode  = 0;
// Not yet implemented, will be supported POST MVC
}
}  // namespace magma5g
