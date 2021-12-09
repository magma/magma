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
#include <vector>
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GExtendedProtocolDiscriminator.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GPDUSessionIdentity.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GPTI.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GMessageType.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GIntegrityProtMaxDataRate.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GPDUSessionType.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GSSCMode.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GQOSRules.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GSessionAMBR.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GPDUAddress.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GDNN.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GNSSAI.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GProtocolConfigurationOptions.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GSMCause.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GGprsTimer.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GQOSRules.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GQosFlowDescriptor.h"

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
