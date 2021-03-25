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
#include "M5GPDUSessionModificationRequest.h"
#include "M5GCommonDefs.h"

namespace magma5g {
PDUSessionModificationRequestMsg::PDUSessionModificationRequestMsg(){};
PDUSessionModificationRequestMsg::~PDUSessionModificationRequestMsg(){};

// Decode PDUSessionModificationRequest Message and its IEs
int PDUSessionModificationRequestMsg::DecodePDUSessionModificationRequestMsg(
    PDUSessionModificationRequestMsg* pdu_session_modif_request,
    uint8_t* buffer, uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, PDU_SESSION_MODIFICATION_REQ_MIN_LEN, len);

  MLOG(MDEBUG) << "DecodePDUSessionModificationRequestMsg\n";
  if ((decoded_result =
           pdu_session_modif_request->extended_protocol_discriminator
               .DecodeExtendedProtocolDiscriminatorMsg(
                   &pdu_session_modif_request->extended_protocol_discriminator,
                   0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           pdu_session_modif_request->pdu_session_identity
               .DecodePDUSessionIdentityMsg(
                   &pdu_session_modif_request->pdu_session_identity, 0,
                   buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = pdu_session_modif_request->pti.DecodePTIMsg(
           &pdu_session_modif_request->pti, 0, buffer + decoded,
           len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           pdu_session_modif_request->message_type.DecodeMessageTypeMsg(
               &pdu_session_modif_request->message_type, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  return 0;
}

// Encode PDUSessionModificationRequest Message and its IEs
int PDUSessionModificationRequestMsg::EncodePDUSessionModificationRequestMsg(
    PDUSessionModificationRequestMsg* pdu_session_modif_request,
    uint8_t* buffer, uint32_t len) {
  MLOG(MDEBUG) << "EncodePDUSessionModificationRequestMsg\n";
  // Not yet implemented, Will be supported POST MVC.
  return 0;
}

}  // namespace magma5g
