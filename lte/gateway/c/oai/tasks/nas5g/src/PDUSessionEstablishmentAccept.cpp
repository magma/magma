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
#include "PDUSessionEstablishmentAccept.h"
#include "CommonDefs.h"

namespace magma5g {
PDUSessionEstablishmentAcceptMsg::PDUSessionEstablishmentAcceptMsg(){};
PDUSessionEstablishmentAcceptMsg::~PDUSessionEstablishmentAcceptMsg(){};

// Decode PDUSessionEstablishmentAccept Message and its IEs
int PDUSessionEstablishmentAcceptMsg::DecodePDUSessionEstablishmentAcceptMsg(
    PDUSessionEstablishmentAcceptMsg* pdusessionestablishmentaccept,
    uint8_t* buffer, uint32_t len) {
  // Not yet implemented, will be supported POST MVC
  return 0;
}

// Encode PDUSessionEstablishmentAccept Message and its IEs
int PDUSessionEstablishmentAcceptMsg::EncodePDUSessionEstablishmentAcceptMsg(
    PDUSessionEstablishmentAcceptMsg* pdusessionestablishmentaccept,
    uint8_t* buffer, uint32_t len) {
  uint32_t encoded  = 0;
  int encodedresult = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, PDU_SESSION_ESTABLISH_ACPT_MIN_LEN, len);

  MLOG(MDEBUG) << "EncodePDUSessionEstablishmentAcceptMsg : \n";
  if ((encodedresult =
           pdusessionestablishmentaccept->extendedprotocoldiscriminator
               .EncodeExtendedProtocolDiscriminatorMsg(
                   &pdusessionestablishmentaccept
                        ->extendedprotocoldiscriminator,
                   0, buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           pdusessionestablishmentaccept->pdusessionidentity
               .EncodePDUSessionIdentityMsg(
                   &pdusessionestablishmentaccept->pdusessionidentity, 0,
                   buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = pdusessionestablishmentaccept->pti.EncodePTIMsg(
           &pdusessionestablishmentaccept->pti, 0, buffer + encoded,
           len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           pdusessionestablishmentaccept->messagetype.EncodeMessageTypeMsg(
               &pdusessionestablishmentaccept->messagetype, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = pdusessionestablishmentaccept->sscmode.EncodeSSCModeMsg(
           &pdusessionestablishmentaccept->sscmode, 0, buffer + encoded,
           len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult = pdusessionestablishmentaccept->pdusessiontype
                           .EncodePDUSessionTypeMsg(
                               &pdusessionestablishmentaccept->pdusessiontype,
                               0, buffer + encoded, len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           pdusessionestablishmentaccept->qosrules.EncodeQOSRulesMsg(
               &pdusessionestablishmentaccept->qosrules, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;
  if ((encodedresult =
           pdusessionestablishmentaccept->sessionambr.EncodeSessionAMBRMsg(
               &pdusessionestablishmentaccept->sessionambr, 0, buffer + encoded,
               len - encoded)) < 0)
    return encodedresult;
  else
    encoded += encodedresult;

  return encoded;
}
}  // namespace magma5g
