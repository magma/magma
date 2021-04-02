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
#include "M5GSecurityHeaderType.h"
#include "M5GMessageType.h"
#include "M5GSpareHalfOctet.h"
#include "M5GSRegistrationResult.h"
#include "M5GSMobileIdentity.h"
#include "M5GUESecurityCapability.h"

using namespace std;
namespace magma5g {
class RegistrationAcceptMsg {
 public:
  ExtendedProtocolDiscriminatorMsg extended_protocol_discriminator;
  SecurityHeaderTypeMsg sec_header_type;
  SpareHalfOctetMsg spare_half_octet;
  MessageTypeMsg message_type;
  M5GSRegistrationResultMsg m5gs_reg_result;
  M5GSMobileIdentityMsg mobile_id;
#define REGISTRATION_ACCEPT_MINIMUM_LENGTH 5

  RegistrationAcceptMsg();
  ~RegistrationAcceptMsg();
  int DecodeRegistrationAcceptMsg(
      RegistrationAcceptMsg* reg_accept, uint8_t* buffer, uint32_t len);
  int EncodeRegistrationAcceptMsg(
      RegistrationAcceptMsg* reg_accept, uint8_t* buffer, uint32_t len);
};
}  // namespace magma5g

/******************************************************************************
         REGISTRATION ACCEPT message content --- TS 24.501 8.2.7.1.1
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
|   |Registration accept     |Message type 9.7        |    M   |  V   |  1    |
|   |message identity        |                        |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
|   |5GS registration result |5GS registration result |    M   |  LV  |  2    |
|   |                        |9.11.3.6                |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
|77 |5G-GUTI                 |5GS mobile identity     |    O   | TLV-E|  14   |
|   |                        |9.11.3.4                |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
******************************************************************************/
