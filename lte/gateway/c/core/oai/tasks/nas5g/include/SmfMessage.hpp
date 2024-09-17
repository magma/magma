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
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GPDUSessionEstablishmentRequest.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GPDUSessionEstablishmentAccept.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GPDUSessionEstablishmentReject.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GPDUSessionReleaseRequest.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GPDUSessionReleaseReject.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GPDUSessionReleaseCommand.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GPDUSessionModificationRequest.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GPDUSessionModificationReject.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GPDUSessionModificationCommand.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GPDUSessionModificationComplete.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GPDUSessionModificationCommandReject.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GNasEnums.h"

namespace magma5g {
// Smf NAS Header Class
class SmfMsgHeader {
 public:
  uint8_t extended_protocol_discriminator;
  uint8_t pdu_session_id;
  uint8_t procedure_transaction_id;
  uint8_t message_type;

  void copy(const SmfMsgHeader& s) {
    extended_protocol_discriminator = s.extended_protocol_discriminator;
    pdu_session_id = s.pdu_session_id;
    procedure_transaction_id = s.procedure_transaction_id;
    message_type = s.message_type;
  }
  bool isEqual(const SmfMsgHeader& s) {
    return ((extended_protocol_discriminator ==
             s.extended_protocol_discriminator) &&
            (pdu_session_id == s.pdu_session_id) &&
            (procedure_transaction_id == s.procedure_transaction_id) &&
            (message_type == s.message_type));
  }
};

// Smf NAS messages
union SMsg_u {
  PDUSessionEstablishmentRequestMsg pdu_session_estab_request;
  PDUSessionEstablishmentAcceptMsg pdu_session_estab_accept;
  PDUSessionEstablishmentRejectMsg pdu_session_estab_reject;
  PDUSessionReleaseRequestMsg pdu_session_release_request;
  PDUSessionReleaseRejectMsg pdu_session_release_reject;
  PDUSessionReleaseCommandMsg pdu_session_release_command;
  PDUSessionModificationRequestMsg pdu_session_modif_request;
  PDUSessionModificationRejectMsg pdu_session_modif_reject;
  PDUSessionModificationCommand pdu_sess_mod_cmd;
  PDUSessionModificationComplete pdu_sess_mod_com;
  PDUSessionModificationCommandReject pdu_sess_mod_cmd_rej;
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
  void copy(const SmfMsg& s) {
    header.copy(s.header);
    switch (static_cast<M5GMessageType>(s.header.message_type)) {
      case M5GMessageType::PDU_SESSION_ESTABLISHMENT_REQUEST:
        msg.pdu_session_estab_request.copy(s.msg.pdu_session_estab_request);
        break;
      default:
        break;
    }
  }
  bool isEqual(const SmfMsg& s) {
    if (!header.isEqual(s.header)) return false;
    bool status = false;
    switch (static_cast<M5GMessageType>(s.header.message_type)) {
      case M5GMessageType::PDU_SESSION_ESTABLISHMENT_REQUEST:
        status = msg.pdu_session_estab_request.isEqual(
            s.msg.pdu_session_estab_request);
        break;
      default:
        break;
    }
    return status;
  }
};
}  // namespace magma5g
