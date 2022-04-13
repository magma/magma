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
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GSecurityHeaderType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GSpareHalfOctet.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GMessageType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GNASKeySetIdentifier.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GABBA.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GAuthenticationParameterRAND.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GAuthenticationParameterAUTN.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GEAPMessage.hpp"

namespace magma5g {
// AuthenticationRequest Message Class
class AuthenticationRequestMsg {
 public:
#define AUTHENTICATION_REQUEST_MINIMUM_LENGTH 3
  ExtendedProtocolDiscriminatorMsg extended_protocol_discriminator;
  SecurityHeaderTypeMsg sec_header_type;
  SpareHalfOctetMsg spare_half_octet;
  MessageTypeMsg message_type;
  NASKeySetIdentifierMsg nas_key_set_identifier;
  ABBAMsg abba;
  AuthenticationParameterAUTNMsg auth_autn;
  AuthenticationParameterRANDMsg auth_rand;
  EAPMessageMsg eap;

  AuthenticationRequestMsg();
  ~AuthenticationRequestMsg();
  int DecodeAuthenticationRequestMsg(AuthenticationRequestMsg* auth_request,
                                     uint8_t* buffer, uint32_t len);
  int EncodeAuthenticationRequestMsg(AuthenticationRequestMsg* auth_request,
                                     uint8_t* buffer, uint32_t len);
};
}  // namespace magma5g

/******************************************************************************
         AUTHENTICATION REQUEST message content --- TS 24.501 8.2.1.1
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
|   |Authentication Request  |Message type 9.7        |    M   |  V   |  1    |
|   |message identity        |                        |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
|   |ABBA                    |ABBA                    |    M   |  LV  |  3-n  |
|---|------------------------|------------------------|--------|------|-------|
|21 |Authentication parameter|Authentication parameter|    O   |  TV  |  17   |
|   |RAND                    |RAND 9.11.3.16          |        |      |       |
|---|------------------------|------------------------|--------|------|-------|
|20 |Authentication parameter|Authentication parameter|    O   |  TLV |  18   |
|   |AUTN                    |AUTN 9.11.3.15          |        |      |       |
-------------------------------------------------------------------------------
******************************************************************************/
