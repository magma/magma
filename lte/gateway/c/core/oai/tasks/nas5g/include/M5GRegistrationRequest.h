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
#include "M5GSpareHalfOctet.h"
#include "M5GSecurityHeaderType.h"
#include "M5GMessageType.h"
#include "M5GSRegistrationType.h"
#include "M5GNASKeySetIdentifier.h"
#include "M5GSMobileIdentity.h"
#include "M5GUESecurityCapability.h"

using namespace std;
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

  RegistrationRequestMsg();
  ~RegistrationRequestMsg();
  int DecodeRegistrationRequestMsg(
      RegistrationRequestMsg* reg_request, uint8_t* buffer, uint32_t len);
  int EncodeRegistrationRequestMsg(
      RegistrationRequestMsg* reg_request, uint8_t* buffer, uint32_t len);
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
|   |Registartion Request    |Message type 9.7        |    M   |  V   |  1    |
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
