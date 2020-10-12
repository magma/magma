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
#include "ExtendedProtocolDiscriminator.h"
#include "SpareHalfOctet.h"
#include "SecurityHeaderType.h"
#include "MessageType.h"
#include "M5GSIdentityType.h"

using namespace std;
namespace magma5g {
// 5GSIdentityRequest Message Class
class IdentityRequestMsg {
 public:
  ExtendedProtocolDiscriminatorMsg extendedprotocoldiscriminator;
  SpareHalfOctetMsg sparehalfoctet;
  SecurityHeaderTypeMsg securityheadertype;
  MessageTypeMsg messagetype;
  M5GSIdentityTypeMsg m5gsidentitytype;
#define IDENTITY_REQUEST_MINIMUM_LENGTH 4

  IdentityRequestMsg();
  ~IdentityRequestMsg();
  int DecodeIdentityRequestMsg(
      IdentityRequestMsg* identityrequest, uint8_t* buffer, uint32_t len);
  int EncodeIdentityRequestMsg(
      IdentityRequestMsg* identityrequest, uint8_t* buffer, uint32_t len);
};
}  // namespace magma5g

/*
   SPEC TS-24501_v150600
   Table 8.2.21.1.1: IDENTITY REQUEST message content

IEI         Information Element                    Type/Reference                     Presence     Format        Length

       Extended protocol discriminator        Extended protocol discriminator 9.2         M           V             1
       Security header type                   Security header type 9.3                    M           V             1/2
       Spare half octet                       Spare half octet 9.5                        M           V             1/2
       Identity request message identity      Message type 9.7                            M           V             1
       Identity type                          5GS identity type 9.11.3.3                  M           V             1/2
       Spare half octet                       Spare half octet 9.5                        M           V             1/2
*/
