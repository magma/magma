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
// PDUSessionModificationRequest Message Class
class PDUSessionModificationRequestMsg {
 public:
#define PDU_SESSION_MODIFICATION_REQ_MIN_LEN 4
  ExtendedProtocolDiscriminatorMsg extended_protocol_discriminator;
  PDUSessionIdentityMsg pdu_session_identity;
  PTIMsg pti;
  MessageTypeMsg message_type;

  PDUSessionModificationRequestMsg();
  ~PDUSessionModificationRequestMsg();
  int DecodePDUSessionModificationRequestMsg(
      PDUSessionModificationRequestMsg* pdu_session_modif_request,
      uint8_t* buffer, uint32_t len);
  int EncodePDUSessionModificationRequestMsg(
      PDUSessionModificationRequestMsg* pdu_session_modif_request,
      uint8_t* buffer, uint32_t len);
};
}  // namespace magma5g

/******************************************************************************
   TS-24.501 Table 8.3.7.1.1: PDU SESSION MODIFICATION REQUEST message content
-------------------------------------------------------------------------------
|IEI|   Information Element  |    Type/Reference      |Presence|Format|Length |
|---|------------------------|------------------------|--------|------|-------|
|   |Extended protocol descr-|Extended Protocol descr-|    M   |  V   |  1    |
|   |-iminator               |-iminator 9.2           |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
|   |PDU session ID          |PDU session ID 9.4      |    M   |  V   |  1    |
|---|------------------------|------------------------|--------|------|-------|
|   |PTI                     |Procedure transacti     |    M   |  V   |  1    |
|   |                        |identity 9.6            |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
|   |PDU SESSION MODIFICATION|Message type 9.7        |    M   |  V   |  1    |
|   |REQUEST message identity|                        |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
******************************************************************************/
