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
#include "ExtendedProtocolDiscriminator.h"
#include "PDUSessionIdentity.h"
#include "PTI.h"
#include "MessageType.h"
#include "5GSMCause.h"

using namespace std;
namespace magma5g {
// PDUSessionEstablishmentReject Message Class
class PDUSessionEstablishmentRejectMsg {
 public:
#define PDU_SESSION_ESTABLISHMENT_REJ_MIN_LEN 5
  ExtendedProtocolDiscriminatorMsg extendedprotocoldiscriminator;
  PDUSessionIdentityMsg pdusessionidentity;
  PTIMsg pti;
  MessageTypeMsg messagetype;
  M5GSMCauseMsg m5gsmcause;

  PDUSessionEstablishmentRejectMsg();
  ~PDUSessionEstablishmentRejectMsg();
  int DecodePDUSessionEstablishmentRejectMsg(
      PDUSessionEstablishmentRejectMsg* pdusessionestablishmentreject,
      uint8_t* buffer, uint32_t len);
  int EncodePDUSessionEstablishmentRejectMsg(
      PDUSessionEstablishmentRejectMsg* pdusessionestablishmentreject,
      uint8_t* buffer, uint32_t len);
};
}  // namespace magma5g
