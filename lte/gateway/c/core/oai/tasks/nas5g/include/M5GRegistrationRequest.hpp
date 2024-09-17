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
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GSpareHalfOctet.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GSecurityHeaderType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GMessageType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GSRegistrationType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GNASKeySetIdentifier.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GSMobileIdentity.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GUESecurityCapability.hpp"

namespace magma5g {
// RegistrationRequest Message Class
class RegistrationRequestMsg {
 public:
  ExtendedProtocolDiscriminatorMsg extended_protocol_discriminator;
  SpareHalfOctetMsg spare_half_octet;
  SecurityHeaderTypeMsg sec_header_type;
  MessageTypeMsg message_type;
  M5GSRegistrationTypeMsg m5gs_reg_type;
  NASKeySetIdentifierMsg nas_key_set_identifier;
  M5GSMobileIdentityMsg m5gs_mobile_identity;
  UESecurityCapabilityMsg ue_sec_capability;
#define REGISTRATION_REQUEST_MINIMUM_LENGTH 5

#define REGISTRATION_REQUEST_5GMM_CAPABILITY_TYPE 0x10
#define REGISTRATION_REQUEST_UE_SECURITY_CAPABILITY_TYPE 0x2E
#define REGISTRATION_REQUEST_REQUESTED_NSSAI_TYPE 0x2F
#define REGISTRATION_REQUEST_LAST_VISITED_REGISTERED_TAI_TYPE 0x52
#define REGISTRATION_REQUEST_S1_UE_NETWORK_CAPABILITY_TYPE 0x17
#define REGISTRATION_REQUEST_UPLINK_DATA_STATUS_TYPE 0x40
#define REGISTRATION_REQUEST_PDU_SESSION_STATUS_TYPE 0x50
#define REGISTRATION_REQUEST_MICO_INDICATION_TYPE 0xB0
#define REGISTRATION_REQUEST_UE_STATUS_TYPE 0x2B
#define REGISTRATION_REQUEST_ADDITIONAL_GUTI_TYPE 0x77
#define REGISTRATION_REQUEST_ALLOWED_PDU_SESSION_STATUS_TYPE 0x25
#define REGISTRATION_REQUEST_UE_USAGE_SETTING_TYPE 0x18
#define REGISTRATION_REQUEST_REQUESTED_DRX_PARAMETERS_TYPE 0x51
#define REGISTRATION_REQUEST_EPS_NAS_MESSAGE_CONTAINER_TYPE 0x70
#define REGISTRATION_REQUEST_LADN_INDICATION_TYPE 0x74
#define REGISTRATION_REQUEST_PAYLOAD_CONTAINER_TYPE_TYPE 0x80
#define REGISTRATION_REQUEST_PAYLOAD_CONTAINER_TYPE 0x7B
#define REGISTRATION_REQUEST_NETWORK_SLICING_INDICATION_TYPE 0x90
#define REGISTRATION_REQUEST_5GS_UPDATE_TYPE_TYPE 0x53
#define REGISTRATION_REQUEST_MOBILE_STATION_CLASSMARK_2_TYPE 0x41
#define REGISTRATION_REQUEST_SUPPORTED_CODECS_TYPE 0x42
#define REGISTRATION_REQUEST_NAS_MESSAGE_CONTAINER_TYPE 0x71
#define REGISTRATION_REQUEST_EPS_BEARER_CONTEXT_STATUS_TYPE 0x60

  RegistrationRequestMsg();
  ~RegistrationRequestMsg();
  int DecodeRegistrationRequestMsg(RegistrationRequestMsg* reg_request,
                                   uint8_t* buffer, uint32_t len);
  int EncodeRegistrationRequestMsg(RegistrationRequestMsg* reg_request,
                                   uint8_t* buffer, uint32_t len);
};
}  // namespace magma5g

/******************************************************************************
   SPEC TS-24501  Table 8.2.6.1.1: REGISTRATION REQUEST message content
-------------------------------------------------------------------------------
|IEI|   Information Element  |    Type/Reference      |Presence|Format|Length |
|---|------------------------|------------------------|--------|------|-------|
|   |Extended protocol descr-|Extended Protocol descr-|    M   |  V   |  1    |
|   |-iminator               |-iminator 9.2           |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
|   |Security header type    |Security header type 9.3|    M   |  V   |  1/2  |
|---|------------------------|------------------------|--------|------|-------|
|   |Spare half octet        |Spare half octet 9.5    |    M   |  V   |  1/2  |
|---|------------------------|------------------------|--------|------|-------|
|   |Registration Request    |Message type 9.7        |    M   |  V   |  1    |
|   |message identity        |                        |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
|   |5GS registration type   |5GS registration type   |    M   |  V   |  1/2  |
|   |                        |9.11.3.7                |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
|   |ngKSI                   |NAS key set identifier  |    M   |  V   |  1/2  |
|   |                        |9.11.3.32               |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
|   |5GS mobile identity     |5GS mobile identity     |    M   | LV-E |  6-n  |
|---|------------------------|------------------------|--------|------|-------|
|2E |UE security capability  |UE security capability  |    O   |  TLV |  4-10 |
|---|------------------------|------------------------|--------|------|-------|
******************************************************************************/
