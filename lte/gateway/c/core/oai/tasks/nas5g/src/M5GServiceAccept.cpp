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
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GServiceAccept.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
ServiceAcceptMsg::ServiceAcceptMsg() {};
ServiceAcceptMsg::~ServiceAcceptMsg() {};

// Decoding Service Accept Message and its IEs
int ServiceAcceptMsg::DecodeServiceAcceptMsg(ServiceAcceptMsg* svc_acpt,
                                             uint8_t* buffer, uint32_t len) {
  return 0;
}

// Encoding Service Accept Message and its IEs
int ServiceAcceptMsg::EncodeServiceAcceptMsg(ServiceAcceptMsg* svc_acpt,
                                             uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;
  int encoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, SERVICE_ACCEPT_MINIMUM_LENGTH,
                                       len);

  if ((encoded_result = svc_acpt->extended_protocol_discriminator
                            .EncodeExtendedProtocolDiscriminatorMsg(
                                &svc_acpt->extended_protocol_discriminator, 0,
                                buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = svc_acpt->spare_half_octet.EncodeSpareHalfOctetMsg(
           &svc_acpt->spare_half_octet, 0, buffer + encoded, len - encoded)) <
      0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = svc_acpt->sec_header_type.EncodeSecurityHeaderTypeMsg(
           &svc_acpt->sec_header_type, 0, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = svc_acpt->message_type.EncodeMessageTypeMsg(
           &svc_acpt->message_type, 0, buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = svc_acpt->pdu_session_status.EncodePDUSessionStatus(
           &svc_acpt->pdu_session_status, 0, buffer + encoded, len - encoded)) <
      0)
    return encoded_result;
  else
    encoded += encoded_result;
  if ((encoded_result = svc_acpt->pdu_re_activation_status
                            .EncodePDUSessionReActivationResult(
                                &svc_acpt->pdu_re_activation_status, 0,
                                buffer + encoded, len - encoded)) < 0)
    return encoded_result;
  else
    encoded += encoded_result;

  return encoded;
}
}  // namespace magma5g
