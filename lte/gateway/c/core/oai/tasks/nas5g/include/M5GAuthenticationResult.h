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
#include "M5GNASKeySetIdentifier.h"
#include "M5GEAPMessage.h"

using namespace std;
namespace magma5g {
class AuthenticationResultMsg {
 public:
  AuthenticationResultMsg();
  ~AuthenticationResultMsg();

  ExtendedProtocolDiscriminatorMsg extended_protocol_discriminator;
  SecurityHeaderTypeMsg sec_header_type;
  SpareHalfOctetMsg spare_half_octet;
  MessageTypeMsg message_type;
  NASKeySetIdentifierMsg nas_key_set_identifier;
  EAPMessageMsg eap_message;
#define AUTHENTICATION_RESULT_MINIMUM_LENGTH 10
  int DecodeAuthenticationResultMsg(
      AuthenticationResultMsg* auth_result, uint8_t* buffer, uint32_t len);
  int EncodeAuthenticationResultMsg(
      AuthenticationResultMsg* auth_result, uint8_t* buffer, uint32_t len);
};
}  // namespace magma5g

/******************************************************************************
                    AUTHENTICATION RESULT message content
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
|   |Authentication result   |Message type 9.7        |    M   |  V   |  1    |
|   |message identity        |                        |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
|   |ngKSI                   |NAS key set identifier  |    M   |  V   |  1/2  |
|   |                        |9.11.3.32               |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
|   |Spare half octet        |Spare half octet 9.5    |    M   |  V   |  1/2  |
|---|------------------------|------------------------|--------|------|-------|
|   |EAP message             |EAP message 9.11.2.2    |    M   |  LV-E| 6-1502|
|---|------------------------|------------------------|--------|------|-------|
******************************************************************************/
