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
#include "M5GPDUSessionModificationCommandReject.h"
#include "M5GCommonDefs.h"

namespace magma5g {
PDUSessionModificationCommandRejectMsg::
    PDUSessionModificationCommandRejectMsg(){};
PDUSessionModificationCommandRejectMsg::
    ~PDUSessionModificationCommandRejectMsg(){};

// Decode PDUSessionModificationCommandReject Message and its IEs
int PDUSessionModificationCommandRejectMsg::
    DecodePDUSessionModificationCommandRejectMsg(
        PDUSessionModificationCommandRejectMsg* pdu_session_modif_request,
        uint8_t* buffer, uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, PDU_SESSION_MODIFICATION_CMD_REJ_MIN_LEN, len);

  MLOG(MDEBUG) << "DecodePDUSessionModificationCommandRejectMsg\n";
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
   if ((decoded_result =
           pdu_session_modif_request->m5gsm_cause.DecodeM5GSMCauseMsg(
               &pdu_session_modif_request->m5gsm_cause, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  return decoded;
}

// Encode PDUSessionModificationCommandReject Message and its IEs
int PDUSessionModificationCommandRejectMsg::
    EncodePDUSessionModificationCommandRejectMsg(
        PDUSessionModificationCommandRejectMsg* pdu_session_modif_request,
        uint8_t* buffer, uint32_t len) {
  MLOG(MDEBUG) << "EncodePDUSessionModificationCommandRejectMsg\n";
  // Not yet implemented, Will be supported POST MVC.
  return 0;
}

}  // namespace magma5g
