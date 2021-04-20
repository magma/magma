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

using namespace std;
namespace magma5g {
class DeRegistrationAcceptUEInitMsg {
 public:
  DeRegistrationAcceptUEInitMsg();
  ~DeRegistrationAcceptUEInitMsg();

  ExtendedProtocolDiscriminatorMsg extended_protocol_discriminator;
  SecurityHeaderTypeMsg sec_header_type;
  SpareHalfOctetMsg spare_half_octet;
  MessageTypeMsg message_type;
#define DEREGISTRATION_ACCEPT_UEINIT_MINIMUM_LENGTH 3
  int DecodeDeRegistrationAcceptUEInitMsg(
      DeRegistrationAcceptUEInitMsg* de_reg_accept_ue, uint8_t* buffer,
      uint32_t len);
  int EncodeDeRegistrationAcceptUEInitMsg(
      DeRegistrationAcceptUEInitMsg* de_reg_accept_ue, uint8_t* buffer,
      uint32_t len);
};
}  // namespace magma5g

/******************************************************************************
               DEREGISTRATION ACCEPT UE Initiated message content
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
|   |De-registration accept  |Message type 9.7        |    M   |  V   |  1    |
|   |message identity        |                        |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
******************************************************************************/
