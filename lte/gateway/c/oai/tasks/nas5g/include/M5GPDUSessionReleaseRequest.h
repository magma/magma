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

using namespace std;
namespace magma5g {
// PDUSessionReleaseRequest Message Class
class PDUSessionReleaseRequestMsg {
 public:
#define PDU_SESSION_RELEASE_REQ_MIN_LEN 4
  ExtendedProtocolDiscriminatorMsg extended_protocol_discriminator;
  PDUSessionIdentityMsg pdu_session_identity;
  PTIMsg pti;
  MessageTypeMsg message_type;

  PDUSessionReleaseRequestMsg();
  ~PDUSessionReleaseRequestMsg();
  int DecodePDUSessionReleaseRequestMsg(
      PDUSessionReleaseRequestMsg* pdu_session_release_request, uint8_t* buffer,
      uint32_t len);
  int EncodePDUSessionReleaseRequestMsg(
      PDUSessionReleaseRequestMsg* pdu_session_release_request, uint8_t* buffer,
      uint32_t len);
};
}  // namespace magma5g
