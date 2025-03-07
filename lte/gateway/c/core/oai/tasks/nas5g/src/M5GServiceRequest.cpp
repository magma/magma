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

#include <iostream>
#include <sstream>
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GServiceRequest.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5gNasMessage.h"

namespace magma5g {
ServiceRequestMsg::ServiceRequestMsg() {};
ServiceRequestMsg::~ServiceRequestMsg() {};

// Decode ServiceRequest Messsage
int ServiceRequestMsg::DecodeServiceRequestMsg(ServiceRequestMsg* svc_req,
                                               uint8_t* buffer, uint32_t len) {
  uint32_t decoded = 0;
  int decoded_result = 0;

  // Check if we got a NULL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(buffer, SERVICE_REQUEST_MINIMUM_LENGTH,
                                       len);

  if ((decoded_result = svc_req->extended_protocol_discriminator
                            .DecodeExtendedProtocolDiscriminatorMsg(
                                &svc_req->extended_protocol_discriminator, 0,
                                buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = svc_req->spare_half_octet.DecodeSpareHalfOctetMsg(
           &svc_req->spare_half_octet, 0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = svc_req->sec_header_type.DecodeSecurityHeaderTypeMsg(
           &svc_req->sec_header_type, 0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = svc_req->message_type.DecodeMessageTypeMsg(
           &svc_req->message_type, 0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = svc_req->service_type.DecodeServiceTypeMsg(
           &svc_req->service_type, 0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           svc_req->nas_key_set_identifier.DecodeNASKeySetIdentifierMsg(
               &svc_req->nas_key_set_identifier, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           svc_req->m5gs_mobile_identity.DecodeM5GSMobileIdentityMsg(
               &svc_req->m5gs_mobile_identity, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  while (decoded < len) {
    uint8_t type = *(buffer + decoded);
    switch (type) {
      case UP_LINK_DATA_STATUS: {
        if ((decoded_result =
                 svc_req->uplink_data_status.DecodeUplinkDataStatus(
                     &svc_req->uplink_data_status, UP_LINK_DATA_STATUS,
                     buffer + decoded, len - decoded)) < 0)
          return decoded_result;
        else
          decoded += decoded_result;
      } break;
      case PDU_SESSION_STATUS: {
        if ((decoded_result =
                 svc_req->pdu_session_status.DecodePDUSessionStatus(
                     &svc_req->pdu_session_status, PDU_SESSION_STATUS,
                     buffer + decoded, len - decoded)) < 0)
          return decoded_result;
        else
          decoded += decoded_result;
      } break;
      case NAS_MESSAGE_CONTAINER: {
        if ((decoded_result = DecodeServiceRequestMsg(
                 svc_req, buffer + (decoded + 3), len - (decoded + 3))) < 0) {
          return decoded_result;
        } else {
          decoded += (decoded_result + 3);
        }
      } break;
      default: {
        return decoded;
      }
    }
  }

  return decoded;
};

// Encode ServiceRequest Messsage
int ServiceRequestMsg::EncodeServiceRequestMsg(ServiceRequestMsg* svc_req,
                                               uint8_t* buffer, uint32_t len) {
  /*** Not Implemented, will be supported POST MVC ***/
  return 0;
};
}  // namespace magma5g
