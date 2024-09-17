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
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GPDUSessionStatus.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GPDUSessionReActivationResult.hpp"

namespace magma5g {
// ServiceAccept Message Class
class ServiceAcceptMsg {
 public:
#define SERVICE_ACCEPT_MINIMUM_LENGTH 3
  ExtendedProtocolDiscriminatorMsg extended_protocol_discriminator;
  SecurityHeaderTypeMsg sec_header_type;
  SpareHalfOctetMsg spare_half_octet;
  MessageTypeMsg message_type;
  M5GPDUSessionStatus pdu_session_status;
  M5GPDUSessionReActivationResult pdu_re_activation_status;

  ServiceAcceptMsg();
  ~ServiceAcceptMsg();
  int DecodeServiceAcceptMsg(ServiceAcceptMsg* service_accept, uint8_t* buffer,
                             uint32_t len);
  int EncodeServiceAcceptMsg(ServiceAcceptMsg* service_accept, uint8_t* buffer,
                             uint32_t len);
};
}  // namespace magma5g
