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
#include "M5GPDUSessionEstablishmentRequest.h"
#include "M5GPDUSessionEstablishmentAccept.h"
#include "M5GPDUSessionEstablishmentReject.h"
#include "M5GPDUSessionReleaseRequest.h"
#include "M5GPDUSessionReleaseReject.h"
#include "M5GPDUSessionModificationRequest.h"
#include "M5GPDUSessionModificationReject.h"

using namespace std;
namespace magma5g {
// Smf NAS Header Class
class SmfMsgHeader {
 public:
  uint8_t extended_protocol_discriminator;
  uint8_t pdu_session_id;
  uint8_t procedure_transaction_id;
  uint8_t message_type;
};

// Smf NAS messages
union SMsg_u {
  PDUSessionEstablishmentRequestMsg pdu_session_estab_request;
  PDUSessionEstablishmentAcceptMsg pdu_session_estab_accept;
  PDUSessionEstablishmentRejectMsg pdu_session_estab_reject;
  PDUSessionReleaseRequestMsg pdu_session_release_request;
  PDUSessionReleaseRejectMsg pdu_session_release_reject;
  PDUSessionModificationRequestMsg pdu_session_modif_request;
  PDUSessionModificationRejectMsg pdu_session_modif_reject;
  SMsg_u();
  ~SMsg_u();
};

// Smf NAS Msg Class
class SmfMsg {
 public:
  SmfMsgHeader header;
  SMsg_u msg;

  SmfMsg();
  ~SmfMsg();
  int SmfMsgDecodeHeaderMsg(SmfMsgHeader* hdr, uint8_t* buffer, uint32_t len);
  int SmfMsgEncodeHeaderMsg(SmfMsgHeader* hdr, uint8_t* buffer, uint32_t len);
  int SmfMsgDecodeMsg(SmfMsg* msg, uint8_t* buffer, uint32_t len);
  int SmfMsgEncodeMsg(SmfMsg* msg, uint8_t* buffer, uint32_t len);
};
}  // namespace magma5g
