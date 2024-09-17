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

#pragma once
#include <sstream>
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GExtendedProtocolDiscriminator.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GSecurityHeaderType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GSpareHalfOctet.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GMessageType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GNASKeySetIdentifier.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GServiceType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GSMobileIdentity.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GUplinkDataStatus.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GPDUSessionStatus.hpp"

namespace magma5g {
// ServiceRequest Message Class
class ServiceRequestMsg {
 public:
#define SERVICE_REQUEST_MINIMUM_LENGTH 13
  ExtendedProtocolDiscriminatorMsg extended_protocol_discriminator;
  SecurityHeaderTypeMsg sec_header_type;
  SpareHalfOctetMsg spare_half_octet;
  MessageTypeMsg message_type;
  NASKeySetIdentifierMsg nas_key_set_identifier;
  ServiceTypeMsg service_type;
  M5GSMobileIdentityMsg m5gs_mobile_identity;
  M5GUplinkDataStatus uplink_data_status;
  M5GPDUSessionStatus pdu_session_status;

  ServiceRequestMsg();
  ~ServiceRequestMsg();
  int DecodeServiceRequestMsg(ServiceRequestMsg* svc_request, uint8_t* buffer,
                              uint32_t len);
  int EncodeServiceRequestMsg(ServiceRequestMsg* svc_request, uint8_t* buffer,
                              uint32_t len);
};
}  // namespace magma5g
