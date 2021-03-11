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
#include "M5GPDUSessionModificationComplete.h"
#include "M5GCommonDefs.h"

namespace magma5g {
PDUSessionModificationCompleteMsg::PDUSessionModificationCompleteMsg(){};
PDUSessionModificationCompleteMsg::~PDUSessionModificationCompleteMsg(){};

// Decode PDUSessionModificationComplete Message and its IEs
int PDUSessionModificationCompleteMsg::DecodePDUSessionModificationCompleteMsg(
    PDUSessionModificationCompleteMsg* pdu_session_modif_complete,
    uint8_t* buffer, uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, PDU_SESSION_MODIFICATION_COMP_MIN_LEN, len);

  MLOG(MDEBUG) << "DecodePDUSessionModificationCompleteMsg\n";
  if ((decoded_result =
           pdu_session_modif_complete->extended_protocol_discriminator
               .DecodeExtendedProtocolDiscriminatorMsg(
                   &pdu_session_modif_complete->extended_protocol_discriminator,
                   0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           pdu_session_modif_complete->pdu_session_identity
               .DecodePDUSessionIdentityMsg(
                   &pdu_session_modif_complete->pdu_session_identity, 0,
                   buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = pdu_session_modif_complete->pti.DecodePTIMsg(
           &pdu_session_modif_complete->pti, 0, buffer + decoded,
           len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           pdu_session_modif_complete->message_type.DecodeMessageTypeMsg(
               &pdu_session_modif_complete->message_type, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  return decoded;
}

// Encode PDUSessionModificationComplete Message and its IEs
int PDUSessionModificationCompleteMsg::EncodePDUSessionModificationCompleteMsg(
    PDUSessionModificationCompleteMsg* pdu_session_modif_complete,
    uint8_t* buffer, uint32_t len) {
  MLOG(MDEBUG) << "EncodePDUSessionModificationCompleteMsg\n";
  // Not yet implemented, Will be supported POST MVC.
  return 0;
}

}  // namespace magma5g
