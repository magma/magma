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
#include "M5GExtendedProtocolDiscriminator.h"
#include "M5GPDUSessionIdentity.h"
#include "M5GPTI.h"
#include "M5GMessageType.h"
#include "M5GIntegrityProtMaxDataRate.h"
#include "M5GPDUSessionType.h"
#include "M5GSSCMode.h"

using namespace std;
namespace magma5g {
// PDUSessionEstablishmentRequest Message Class
class PDUSessionEstablishmentRequestMsg {
 public:
#define PDU_SESSION_ESTABLISH_REQ_MIN_LEN 4
  ExtendedProtocolDiscriminatorMsg extended_protocol_discriminator;
  PDUSessionIdentityMsg pdu_session_identity;
  PTIMsg pti;
  MessageTypeMsg message_type;
  IntegrityProtMaxDataRateMsg integrity_prot_max_data_rate;
  PDUSessionTypeMsg pdu_session_type;
  SSCModeMsg ssc_mode;

  PDUSessionEstablishmentRequestMsg();
  ~PDUSessionEstablishmentRequestMsg();
  int DecodePDUSessionEstablishmentRequestMsg(
      PDUSessionEstablishmentRequestMsg* pdu_session_estab_request,
      uint8_t* buffer, uint32_t len);
  int EncodePDUSessionEstablishmentRequestMsg(
      PDUSessionEstablishmentRequestMsg* pdu_session_estab_request,
      uint8_t* buffer, uint32_t len);
};
}  // namespace magma5g
