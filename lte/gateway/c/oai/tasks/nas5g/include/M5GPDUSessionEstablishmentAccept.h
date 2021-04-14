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
#include "M5GQOSRules.h"
#include "M5GSessionAMBR.h"
#include "M5GPDUAddress.h"
#include "M5GDNN.h"
using namespace std;
namespace magma5g {
// PDUSessionEstablishmentAccept Message Class
class PDUSessionEstablishmentAcceptMsg {
 public:
#define PDU_SESSION_ESTABLISH_ACPT_MIN_LEN 18
  ExtendedProtocolDiscriminatorMsg extended_protocol_discriminator;
  PDUSessionIdentityMsg pdu_session_identity;
  PTIMsg pti;
  MessageTypeMsg message_type;
  PDUSessionTypeMsg pdu_session_type;
  SSCModeMsg ssc_mode;
  QOSRulesMsg qos_rules;
  SessionAMBRMsg session_ambr;
  PDUAddressMsg pdu_address;
  DNNMsg dnn;

  PDUSessionEstablishmentAcceptMsg();
  ~PDUSessionEstablishmentAcceptMsg();
  int DecodePDUSessionEstablishmentAcceptMsg(
      PDUSessionEstablishmentAcceptMsg* pdu_session_estab_accept,
      uint8_t* buffer, uint32_t len);
  int EncodePDUSessionEstablishmentAcceptMsg(
      PDUSessionEstablishmentAcceptMsg* pdu_session_estab_accept,
      uint8_t* buffer, uint32_t len);
};
}  // namespace magma5g

/******************************************************************************
   TS-24.501 Table 8.3.2.1.1: PDU SESSION ESTABLISHMENT ACCEPT message content
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
|   |PDU SESSION ESTABLISHME-|Message type 9.7        |    M   |  V   |  1    |
|   |-NT ACCEPT message iden-|                        |        |      |       |
|   |-tity                   |                        |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
|   |Selected PDU session    |PDU session type        |    M   |  V   |  1/2  |
|   |Type                    |9.11.4.11               |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
|   |Selected SSC mode       |SSC mode 9.11.4,16      |    M   |  V   |  1/2  |
|---|------------------------|------------------------|--------|------|-------|
|   |Authorized QoS rules    |QoS rules 9.11.4.13     |    M   | LV-E |6-65538|
|---|------------------------|------------------------|--------|------|-------|
|   |Session AMBR            |Session-AMBR 9.11.4.14  |    M   |  LV  |  7    |
|---|------------------------|------------------------|--------|------|-------|
******************************************************************************/
