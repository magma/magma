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
ServiceRequestMsg::ServiceRequestMsg(){};
ServiceRequestMsg::~ServiceRequestMsg(){};

// Decode ServiceRequest Messsage
int ServiceRequestMsg::DecodeServiceRequestMsg(uint8_t* buffer, uint32_t len) {
  uint32_t decoded = 0;
  int decoded_result = 0;

  // Check if we got a NULL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(buffer, SERVICE_REQUEST_MINIMUM_LENGTH,
                                       len);

  if ((decoded_result = this->extended_protocol_discriminator
                            .DecodeExtendedProtocolDiscriminatorMsg(
                                0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = this->spare_half_octet.DecodeSpareHalfOctetMsg(
           0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = this->sec_header_type.DecodeSecurityHeaderTypeMsg(
           0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = this->message_type.DecodeMessageTypeMsg(
           0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = this->service_type.DecodeServiceTypeMsg(
           0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           this->nas_key_set_identifier.DecodeNASKeySetIdentifierMsg(
               0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = this->m5gs_mobile_identity.DecodeM5GSMobileIdentityMsg(
           0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  while (decoded < len) {
    uint8_t type = *(buffer + decoded);
    switch (type) {
      case UP_LINK_DATA_STATUS: {
        if ((decoded_result = this->uplink_data_status.DecodeUplinkDataStatus(
                 UP_LINK_DATA_STATUS, buffer + decoded, len - decoded)) < 0)
          return decoded_result;
        else
          decoded += decoded_result;
      } break;
      case PDU_SESSION_STATUS: {
        if ((decoded_result = this->pdu_session_status.DecodePDUSessionStatus(
                 PDU_SESSION_STATUS, buffer + decoded, len - decoded)) < 0)
          return decoded_result;
        else
          decoded += decoded_result;
      } break;
      case NAS_MESSAGE_CONTAINER: {
        if ((decoded_result = this->DecodeServiceRequestMsg(
                 buffer + (decoded + 3), len - (decoded + 3))) < 0) {
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
int ServiceRequestMsg::EncodeServiceRequestMsg(uint8_t* buffer, uint32_t len) {
  /*** Not Implemented, will be supported POST MVC ***/
  return 0;
};
}  // namespace magma5g
