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
#include "M5gNasMessage.h"
#include "M5GExtendedProtocolDiscriminator.h"
#include "M5GSecurityHeaderType.h"
#include "M5GSpareHalfOctet.h"
#include "M5GMessageType.h"
#include "M5GAuthenticationResponseParameter.h"

using namespace std;

namespace magma5g {
// AuthenticationResponse Message Class
class AuthenticationResponseMsg {
 public:
#define AUTHENTICATION_RESPONSE_MINIMUM_LENGTH 3
  ExtendedProtocolDiscriminatorMsg extended_protocol_discriminator;
  SecurityHeaderTypeMsg sec_header_type;
  SpareHalfOctetMsg spare_half_octet;
  MessageTypeMsg message_type;
  AuthenticationResponseParameterMsg autn_response_parameter;

  AuthenticationResponseMsg();
  ~AuthenticationResponseMsg();
  int DecodeAuthenticationResponseMsg(
      AuthenticationResponseMsg* auth_response, uint8_t* buffer, uint32_t len);
  int EncodeAuthenticationResponseMsg(
      AuthenticationResponseMsg* auth_response, uint8_t* buffer, uint32_t len);
};
}  // namespace magma5g

/******************************************************************************
   Table 8.2.2.1.1: AUTHENTICATION RESPONSE message content --- TS 24.501
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
|   |Authentication Response |Message type 9.7        |    M   |  V   |  1    |
|   |message identity        |                        |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
|2D |Authentication response |Authentication response |    O   |  TLV |  18   |
|   |parameter               |parameter 9.11.3.17     |        |              |
-------------------------------------------------------------------------------
******************************************************************************/
