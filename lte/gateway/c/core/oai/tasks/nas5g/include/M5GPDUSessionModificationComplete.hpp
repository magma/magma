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
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GProtocolConfigurationOptions.hpp"

namespace magma5g {

// PDUSessionModificationCommandComplete Message Class
class PDUSessionModificationComplete {
 public:
#define PDU_SESSION_MODIFICATION_COMPLETE_MIN_LEN 4
  ExtendedProtocolDiscriminatorMsg extended_protocol_discriminator;
  PDUSessionIdentityMsg pdu_session_identity;
  PTIMsg pti;
  MessageTypeMsg message_type;
  ProtocolConfigurationOptions extProtocolconfigurationoptions;
  PDUSessionModificationComplete();
  ~PDUSessionModificationComplete();
  int DecodePDUSessionModificationComplete(
      PDUSessionModificationComplete* pdu_sess_mod_com, uint8_t* buffer,
      uint32_t len);
  int EncodePDUSessionModificationComplete(
      PDUSessionModificationComplete* pdu_sess_mod_com, uint8_t* buffer,
      uint32_t len);
};
}  // namespace magma5g
