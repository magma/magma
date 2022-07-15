/*
   Copyright 2022 The Magma Authors.
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

#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GPDUSessionModificationComplete.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5gNasMessage.h"

namespace magma5g {
PDUSessionModificationComplete::PDUSessionModificationComplete() {}
PDUSessionModificationComplete::~PDUSessionModificationComplete() {}

int PDUSessionModificationComplete::EncodePDUSessionModificationComplete(
    PDUSessionModificationComplete* pdu_sess_mod_com, uint8_t* buffer,
    uint32_t len) {
  uint32_t encoded = 0;
  uint32_t encoded_result = 0;

  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, PDU_SESSION_MODIFICATION_COMPLETE_MIN_LEN, len);

  OAILOG_DEBUG(LOG_NAS5G, "EncodePDUSessionModificationComplete");
  if ((encoded_result =
           pdu_sess_mod_com->extended_protocol_discriminator
               .EncodeExtendedProtocolDiscriminatorMsg(
                   &pdu_sess_mod_com->extended_protocol_discriminator, 0,
                   buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result =
           pdu_sess_mod_com->pdu_session_identity.EncodePDUSessionIdentityMsg(
               &pdu_sess_mod_com->pdu_session_identity, 0, buffer + encoded,
               len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result = pdu_sess_mod_com->pti.EncodePTIMsg(
           &pdu_sess_mod_com->pti, 0, buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result = pdu_sess_mod_com->message_type.EncodeMessageTypeMsg(
           &pdu_sess_mod_com->message_type, 0, buffer + encoded,
           len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }

  return encoded;
}

int PDUSessionModificationComplete::DecodePDUSessionModificationComplete(
    PDUSessionModificationComplete* pdu_sess_mod_com, uint8_t* buffer,
    uint32_t len) {
  uint32_t decoded = 0;
  uint32_t decoded_result = 0;

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, PDU_SESSION_MODIFICATION_COMPLETE_MIN_LEN, len);

  OAILOG_DEBUG(LOG_NAS5G, "DecodePDUSessionModificationComplete");
  if ((decoded_result =
           pdu_sess_mod_com->extended_protocol_discriminator
               .DecodeExtendedProtocolDiscriminatorMsg(
                   &pdu_sess_mod_com->extended_protocol_discriminator, 0,
                   buffer + decoded, len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  if ((decoded_result =
           pdu_sess_mod_com->pdu_session_identity.DecodePDUSessionIdentityMsg(
               &pdu_sess_mod_com->pdu_session_identity, 0, buffer + decoded,
               len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  if ((decoded_result = pdu_sess_mod_com->pti.DecodePTIMsg(
           &pdu_sess_mod_com->pti, 0, buffer + decoded, len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  if ((decoded_result = pdu_sess_mod_com->message_type.DecodeMessageTypeMsg(
           &pdu_sess_mod_com->message_type, 0, buffer + decoded,
           len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  return decoded;
}
}  // namespace magma5g
