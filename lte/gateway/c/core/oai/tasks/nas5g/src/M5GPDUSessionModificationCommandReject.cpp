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
    EncodePDUSessionModificationCommandReject(uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;
  uint32_t encoded_result = 0;

  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, PDU_SESSION_MODIFICATION_COMMAND_REJECT_MIN_LEN, len);

  OAILOG_DEBUG(LOG_NAS5G, "EncodePDUSessionModificationCommandReject");
  if ((encoded_result = this->extended_protocol_discriminator
                            .EncodeExtendedProtocolDiscriminatorMsg(
                                0, buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result = this->pdu_session_identity.EncodePDUSessionIdentityMsg(
           0, buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result =
           this->pti.EncodePTIMsg(0, buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result = this->message_type.EncodeMessageTypeMsg(
           0, buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result = this->cause.EncodeM5GSMCauseMsg(
           M5GSM_CAUSE, buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }

  return encoded;
}

int PDUSessionModificationCommandReject::
    DecodePDUSessionModificationCommandReject(uint8_t* buffer, uint32_t len) {
  uint32_t decoded = 0;
  uint32_t decoded_result = 0;

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, PDU_SESSION_MODIFICATION_COMMAND_REJECT_MIN_LEN, len);

  OAILOG_DEBUG(LOG_NAS5G, "DecodePDUSessionModificationCommandReject : ");
  if ((decoded_result = this->extended_protocol_discriminator
                            .DecodeExtendedProtocolDiscriminatorMsg(
                                0, buffer + decoded, len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  if ((decoded_result = this->pdu_session_identity.DecodePDUSessionIdentityMsg(
           0, buffer + decoded, len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  if ((decoded_result =
           this->pti.DecodePTIMsg(0, buffer + decoded, len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  if ((decoded_result = this->message_type.DecodeMessageTypeMsg(
           0, buffer + decoded, len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  if ((decoded_result = this->cause.DecodeM5GSMCauseMsg(
           M5GSM_CAUSE, buffer + decoded, len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  if (decoded < len) {
    if ((decoded_result =
             this->extProtocolconfigurationoptions
                 .DecodeProtocolConfigurationOptions(
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
