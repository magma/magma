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
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GPDUSessionIdentity.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GPTI.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GMessageType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GSMCause.hpp"

namespace magma5g {
// PDUSessionReleaseCommand Message Class
class PDUSessionReleaseCommandMsg {
 public:
#define PDU_SESSION_RELEASE_COMMAND_MIN_LEN 5
  ExtendedProtocolDiscriminatorMsg extended_protocol_discriminator;
  PDUSessionIdentityMsg pdu_session_identity;
  PTIMsg pti;
  MessageTypeMsg message_type;
  M5GSMCauseMsg m5gsm_cause;

  PDUSessionReleaseCommandMsg();
  ~PDUSessionReleaseCommandMsg();
  int DecodePDUSessionReleaseCommandMsg(
      PDUSessionReleaseCommandMsg* pdu_session_release_command, uint8_t* buffer,
      uint32_t len);
  int EncodePDUSessionReleaseCommandMsg(
      PDUSessionReleaseCommandMsg* pdu_session_release_command, uint8_t* buffer,
      uint32_t len);
};
}  // namespace magma5g
