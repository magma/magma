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

#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GPDUSessionModificationCommandReject.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5gNasMessage.h"

namespace magma5g {
PDUSessionModificationCommandReject::PDUSessionModificationCommandReject() {}
PDUSessionModificationCommandReject::~PDUSessionModificationCommandReject() {}

int PDUSessionModificationCommandReject::
    EncodePDUSessionModificationCommandReject(
        PDUSessionModificationCommandReject* pdu_sess_mod_comd_rej,
        uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;
  uint32_t encoded_result = 0;

  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, PDU_SESSION_MODIFICATION_COMMAND_REJECT_MIN_LEN, len);

  OAILOG_DEBUG(LOG_NAS5G, "EncodePDUSessionModificationCommandReject");
  if ((encoded_result =
           pdu_sess_mod_comd_rej->extended_protocol_discriminator
               .EncodeExtendedProtocolDiscriminatorMsg(
                   &pdu_sess_mod_comd_rej->extended_protocol_discriminator, 0,
                   buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result = pdu_sess_mod_comd_rej->pdu_session_identity
                            .EncodePDUSessionIdentityMsg(
                                &pdu_sess_mod_comd_rej->pdu_session_identity, 0,
                                buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result = pdu_sess_mod_comd_rej->pti.EncodePTIMsg(
           &pdu_sess_mod_comd_rej->pti, 0, buffer + encoded, len - encoded)) <
      0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result =
           pdu_sess_mod_comd_rej->message_type.EncodeMessageTypeMsg(
               &pdu_sess_mod_comd_rej->message_type, 0, buffer + encoded,
               len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result = pdu_sess_mod_comd_rej->cause.EncodeM5GSMCauseMsg(
           &pdu_sess_mod_comd_rej->cause, M5GSM_CAUSE, buffer + encoded,
           len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }

  return encoded;
}

int PDUSessionModificationCommandReject::
    DecodePDUSessionModificationCommandReject(
        PDUSessionModificationCommandReject* pdu_sess_mod_comd_rej,
        uint8_t* buffer, uint32_t len) {
  uint32_t decoded = 0;
  uint32_t decoded_result = 0;

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, PDU_SESSION_MODIFICATION_COMMAND_REJECT_MIN_LEN, len);

  OAILOG_DEBUG(LOG_NAS5G, "DecodePDUSessionModificationCommandReject : ");
  if ((decoded_result =
           pdu_sess_mod_comd_rej->extended_protocol_discriminator
               .DecodeExtendedProtocolDiscriminatorMsg(
                   &pdu_sess_mod_comd_rej->extended_protocol_discriminator, 0,
                   buffer + decoded, len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  if ((decoded_result = pdu_sess_mod_comd_rej->pdu_session_identity
                            .DecodePDUSessionIdentityMsg(
                                &pdu_sess_mod_comd_rej->pdu_session_identity, 0,
                                buffer + decoded, len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  if ((decoded_result = pdu_sess_mod_comd_rej->pti.DecodePTIMsg(
           &pdu_sess_mod_comd_rej->pti, 0, buffer + decoded, len - decoded)) <
      0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  if ((decoded_result =
           pdu_sess_mod_comd_rej->message_type.DecodeMessageTypeMsg(
               &pdu_sess_mod_comd_rej->message_type, 0, buffer + decoded,
               len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  if ((decoded_result = pdu_sess_mod_comd_rej->cause.DecodeM5GSMCauseMsg(
           &pdu_sess_mod_comd_rej->cause, M5GSM_CAUSE, buffer + decoded,
           len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  if (decoded < len) {
    if ((decoded_result =
             pdu_sess_mod_comd_rej->extProtocolconfigurationoptions
                 .DecodeProtocolConfigurationOptions(
                     &pdu_sess_mod_comd_rej->extProtocolconfigurationoptions,
                     REQUEST_EXTENDED_PROTOCOL_CONFIGURATION_OPTIONS_TYPE,
                     buffer + decoded, len - decoded)) < 0) {
      return decoded_result;
    } else {
      decoded += decoded_result;
    }
  }

  return decoded;
}
}  // namespace magma5g
