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
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5gNasMessage.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GExtendedProtocolDiscriminator.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GPDUSessionIdentity.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GPTI.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GMessageType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GIntegrityProtMaxDataRate.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GPDUSessionType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GSSCMode.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GProtocolConfigurationOptions.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GMaxNumOfSupportedPacketFilters.hpp"

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
  M5GMaxNumOfSupportedPacketFilters maxNumOfSuppPacketFilters;
  ProtocolConfigurationOptions protocolconfigurationoptions;

  PDUSessionEstablishmentRequestMsg();
  ~PDUSessionEstablishmentRequestMsg();
  int DecodePDUSessionEstablishmentRequestMsg(
      PDUSessionEstablishmentRequestMsg* pdu_session_estab_request,
      uint8_t* buffer, uint32_t len);
  int EncodePDUSessionEstablishmentRequestMsg(
      PDUSessionEstablishmentRequestMsg* pdu_session_estab_request,
      uint8_t* buffer, uint32_t len);
  void copy(const PDUSessionEstablishmentRequestMsg& p) {
    extended_protocol_discriminator.copy(p.extended_protocol_discriminator);
    pdu_session_identity.copy(p.pdu_session_identity);
    pti.copy(p.pti);
    message_type.copy(p.message_type);
    integrity_prot_max_data_rate.copy(p.integrity_prot_max_data_rate);
    pdu_session_type.copy(p.pdu_session_type);
    ssc_mode.copy(p.ssc_mode);
  }
  bool isEqual(const PDUSessionEstablishmentRequestMsg& p) {
    return (
        extended_protocol_discriminator.isEqual(
            p.extended_protocol_discriminator) &&
        pdu_session_identity.isEqual(p.pdu_session_identity) &&
        pti.isEqual(p.pti) && message_type.isEqual(p.message_type) &&
        integrity_prot_max_data_rate.isEqual(p.integrity_prot_max_data_rate) &&
        pdu_session_type.isEqual(p.pdu_session_type) &&
        ssc_mode.isEqual(p.ssc_mode));
  }
};
}  // namespace magma5g
/******************************************************************************
  TS-24.501 Table 8.3.1.1.1: PDU SESSION ESTABLISHMENT Request message content
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
|   |-NT Request message ide-|                        |        |      |       |
|   |-ntity                  |                        |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
|   |Integrity protection    |Integrity protection    |    M   |  V   |  2    |
|   |maximum data rate       |maximum data rate       |        |      |       |
|   |                        |9.11.4.7                |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
|9- |PDU session type        |PDU session type        |    O   |  TV  |  1    |
|   |                        |9.11.4.14               |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
|A- |SSC mode                |SSC mode 9.11.4,16      |    O   |  V   |  1/2  |
|---|------------------------|------------------------|--------|------|-------|
******************************************************************************/
