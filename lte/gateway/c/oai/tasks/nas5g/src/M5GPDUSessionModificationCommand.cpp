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
#include "M5GPDUSessionModificationCommand.h"
#include "M5GCommonDefs.h"

namespace magma5g {
PDUSessionModificationCommandMsg::PDUSessionModificationCommandMsg(){};
PDUSessionModificationCommandMsg::~PDUSessionModificationCommandMsg(){};

// Decode PDUSessionModificationReject Message and its IEs
int PDUSessionModificationCommandMsg::DecodePDUSessionModificationCommandMsg(
    PDUSessionModificationCommandMsg* pdu_session_modif_cmd, uint8_t* buffer,
    uint32_t len) {
  // Not yet implemented, Will be supported POST MVC
  return 0;
}
// Encode PDUSessionModificationReject Message and its IEs
int PDUSessionModificationCommandMsg::EncodePDUSessionModificationCommandMsg(
    PDUSessionModificationCommandMsg* pdu_session_modif_cmd, uint8_t* buffer,
    uint32_t len) {
  uint32_t encoded   = 0;
  int encoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, PDU_SESSION_MODIFICATION_CMD_MIN_LEN, len);

  MLOG(MDEBUG) << "EncodePDUSessionModificationCommandMsg : \n";
  if ((encoded_result =
           pdu_session_modif_cmd->extended_protocol_discriminator
               .EncodeExtendedProtocolDiscriminatorMsg(
                   &pdu_session_modif_cmd->extended_protocol_discriminator, 0,
                   buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = pdu_session_modif_cmd->pdu_session_identity
                            .EncodePDUSessionIdentityMsg(
                                &pdu_session_modif_cmd->pdu_session_identity, 0,
                                buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = pdu_session_modif_cmd->pti.EncodePTIMsg(
           &pdu_session_modif_cmd->pti, 0, buffer + encoded, len - encoded)) <
      0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           pdu_session_modif_cmd->message_type.EncodeMessageTypeMsg(
               &pdu_session_modif_cmd->message_type, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  return encoded;
}
}  // namespace magma5g
