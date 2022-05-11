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

#pragma once
#include <sstream>
#include <vector>

#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GExtendedProtocolDiscriminator.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GPDUSessionIdentity.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GPTI.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GMessageType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GIntegrityProtMaxDataRate.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GPDUSessionType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GSSCMode.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GQOSRules.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GSessionAMBR.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GPDUAddress.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GDNN.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GNSSAI.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GProtocolConfigurationOptions.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GSMCause.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GGprsTimer.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GQosFlowDescriptor.hpp"

namespace magma5g {

// PDUSessionModificationCommand Message Class
class PDUSessionModificationCommandReject {
 public:
#define PDU_SESSION_MODIFICATION_COMMAND_REJECT_MIN_LEN 6
  ExtendedProtocolDiscriminatorMsg extended_protocol_discriminator;
  PDUSessionIdentityMsg pdu_session_identity;
  PTIMsg pti;
  MessageTypeMsg message_type;
  M5GSMCauseMsg cause;
  ProtocolConfigurationOptions extProtocolconfigurationoptions;

  PDUSessionModificationCommandReject();
  ~PDUSessionModificationCommandReject();
  int DecodePDUSessionModificationCommandReject(
      PDUSessionModificationCommandReject* pdu_sess_mod_comd_rej,
      uint8_t* buffer, uint32_t len);
  int EncodePDUSessionModificationCommandReject(
      PDUSessionModificationCommandReject* pdu_sess_mod_comd_rej,
      uint8_t* buffer, uint32_t len);
};
}  // namespace magma5g
