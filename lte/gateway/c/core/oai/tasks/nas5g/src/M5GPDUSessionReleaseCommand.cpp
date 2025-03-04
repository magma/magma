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
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GPDUSessionReleaseCommand.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
PDUSessionReleaseCommandMsg::PDUSessionReleaseCommandMsg() {};
PDUSessionReleaseCommandMsg::~PDUSessionReleaseCommandMsg() {};

// Decode PDUSessionReleaseCommand Message and its IEs
int PDUSessionReleaseCommandMsg::DecodePDUSessionReleaseCommandMsg(
    PDUSessionReleaseCommandMsg* pdu_session_release_command, uint8_t* buffer,
    uint32_t len) {
  // Not yet implemented, will be supported POST MVC
  return 0;
}

// Encode PDUSessionReleaseCommand Message and its IEs
int PDUSessionReleaseCommandMsg::EncodePDUSessionReleaseCommandMsg(
    PDUSessionReleaseCommandMsg* pdu_session_release_command, uint8_t* buffer,
    uint32_t len) {
  uint32_t encoded = 0;
  int encoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, PDU_SESSION_RELEASE_COMMAND_MIN_LEN, len);

  if ((encoded_result =
           pdu_session_release_command->extended_protocol_discriminator
               .EncodeExtendedProtocolDiscriminatorMsg(
                   &pdu_session_release_command
                        ->extended_protocol_discriminator,
                   0, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           pdu_session_release_command->pdu_session_identity
               .EncodePDUSessionIdentityMsg(
                   &pdu_session_release_command->pdu_session_identity, 0,
                   buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = pdu_session_release_command->pti.EncodePTIMsg(
           &pdu_session_release_command->pti, 0, buffer + encoded,
           len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           pdu_session_release_command->message_type.EncodeMessageTypeMsg(
               &pdu_session_release_command->message_type, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result =
           pdu_session_release_command->m5gsm_cause.EncodeM5GSMCauseMsg(
               &pdu_session_release_command->m5gsm_cause, 0, buffer + encoded,
               len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  return encoded;
}
}  // namespace magma5g
