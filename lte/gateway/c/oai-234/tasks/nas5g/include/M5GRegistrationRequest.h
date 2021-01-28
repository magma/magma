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

/*
   SPEC TS-24501_v150600
   Table 8.2.6.1.1: REGISTRATION REQUEST message content

   IEI          Information Element                            Type/Reference
   Presence     Format     Length

            Extended protocol discriminator               Extended Protocol
   discriminator 9.2          M           V          1 Security header type
   Security header type 9.3                     M           V          1/2 Spare
   half octet                              Spare half octet 9.5 M           V
   1/2 Registration request message identity         Message type 9.7 M V 1 5GS
   registration type                         5GS registration type 9.11.3.7 M V
   1/2 ngKSI                                         NAS key set
   identifier 9.11.3.32             M           V          1/2 5GS mobile
   identity                           5GS mobile identity 9.11.3.4 M LV-E 6-n C-
   Non-current native NAS key set identifier     NAS key set
   identifier 9.11.3.32             O           TV         1 10       5GMM
   capability                               5GMM capability 9.11.3.1 O TLV 3-15
   2E       UE security capability                        UE security
   capability 9.11.3.54             O           TLV        4-10 2F Requested
   NSSAI                               NSSAI 9.11.3.37 O           TLV 4-74 52
   Last visited registered TAI                   5GS tracking area
   identity 9.11.3.8          O           TV         7 17 S1       UE network
   capability                         S1 UE network capability 9.11.3.48 O TLV
   4-15 40       Uplink data status                            Uplink data
   status 9.11.3.57                 O           TLV        4-34 50       PDU
   session status                            PDU session status 9.11.3.44 O TLV
   4-34 B-       MICO indication                               MICO
   indication 9.11.3.31                    O           TV         1 2 B UE
   status                                     UE status 9.11.3.56 O TLV        3
   77       Additional GUTI                               5GS mobile
   identity 9.11.3.4                 O           TLV-E      14 25       Allowed
   PDU session status                    Allowed PDU session status 9.11.3.13 O
   TLV        4-34 18       UE's usage setting                            UE's
   usage setting 9.11.3.55                 O           TLV        3 51 Requested
   DRX parameters                      5GS DRX parameters 9.11.3.2A O TLV 3 70
   EPS NAS message container                     EPS NAS message
   container 9.11.3.24          O           TLV-E      4-n 74       LADN
   indication                               LADN indication 9.11.3.29 O TLV-E
   3-811 8-       Payload container type                        Payload
   container type 9.11.3.40             O           TV         1 7B Payload
   container                             Payload container 9.11.3.39 O TLV-E
   4-65538 9-       Network slicing indication                    Network
   slicing indication 9.11.3.36         O           TV         1 53       5GS
   update type                               5GS update type 9.11.3.9A O TLV 3
   71       NAS message container                         NAS message
   container 9.11.3.33              O           TLV-E      4-n 60       EPS
   bearer context status                     EPS bearer context
   status 9.11.3.23A         O           TLV        4
 */
