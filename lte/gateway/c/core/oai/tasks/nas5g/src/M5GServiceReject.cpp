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
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5gNasMessage.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GServiceReject.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
ServiceRejectMsg::ServiceRejectMsg() {};
ServiceRejectMsg::~ServiceRejectMsg() {};

// Decoding Service Reject Message and its IEs
int ServiceRejectMsg::DecodeServiceRejectMsg(ServiceRejectMsg* svc_rej,
                                             uint8_t* buffer, uint32_t len) {
  uint32_t decoded = 0;
  int decoded_result = 0;

  if ((decoded_result = svc_rej->extended_protocol_discriminator
                            .DecodeExtendedProtocolDiscriminatorMsg(
                                &svc_rej->extended_protocol_discriminator, 0,
                                buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = svc_rej->spare_half_octet.DecodeSpareHalfOctetMsg(
           &svc_rej->spare_half_octet, 0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = svc_rej->sec_header_type.DecodeSecurityHeaderTypeMsg(
           &svc_rej->sec_header_type, 0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = svc_rej->message_type.DecodeMessageTypeMsg(
           &svc_rej->message_type, 0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = svc_rej->cause.DecodeM5GMMCauseMsg(
           &svc_rej->cause, 0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if (decoded < len) {
    if ((decoded_result = svc_rej->pdu_session_status.DecodePDUSessionStatus(
             &svc_rej->pdu_session_status, PDU_SESSION_STATUS, buffer + decoded,
             len - decoded)) < 0)
      return decoded_result;
    else
      decoded += decoded_result;
  }

  return decoded;
}

// Encoding Service Reject Message and its IEs
int ServiceRejectMsg::EncodeServiceRejectMsg(ServiceRejectMsg* svc_rej,
                                             uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;
  int encoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer,
                                       M5G_SERVICE_REJECT_MINIMUM_LENGTH, len);

  if ((encoded_result = svc_rej->extended_protocol_discriminator
                            .EncodeExtendedProtocolDiscriminatorMsg(
                                &svc_rej->extended_protocol_discriminator, 0,
                                buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = svc_rej->spare_half_octet.EncodeSpareHalfOctetMsg(
           &svc_rej->spare_half_octet, 0, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = svc_rej->sec_header_type.EncodeSecurityHeaderTypeMsg(
           &svc_rej->sec_header_type, 0, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = svc_rej->message_type.EncodeMessageTypeMsg(
           &svc_rej->message_type, 0, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = svc_rej->cause.EncodeM5GMMCauseMsg(
           &svc_rej->cause, 0, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = svc_rej->pdu_session_status.EncodePDUSessionStatus(
           &svc_rej->pdu_session_status, 0, buffer + encoded, len - encoded)) <
      0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = svc_rej->t3346Value.EncodeGPRSTimer2Msg(
           &svc_rej->t3346Value, 0, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;

  return encoded;
}
}  // namespace magma5g
