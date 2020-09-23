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
#include "SecurityHeaderType.h"
#include "MessageType.h"
#include "SpareHalfOctet.h"

using namespace std;
namespace magma5g {
class DeRegistrationAcceptUEInitMsg {
 public:
  DeRegistrationAcceptUEInitMsg();
  ~DeRegistrationAcceptUEInitMsg();

  ExtendedProtocolDiscriminatorMsg extendedprotocoldiscriminator;
  SecurityHeaderTypeMsg securityheadertype;
  SpareHalfOctetMsg sparehalfoctet;
  MessageTypeMsg messagetype;
#define DEREGISTRATION_ACCEPT_UEINIT_MINIMUM_LENGTH 3
  int DecodeDeRegistrationAcceptUEInitMsg(
      DeRegistrationAcceptUEInitMsg* deregistrationacceptueinit,
      uint8_t* buffer, uint32_t len);
  int EncodeDeRegistrationAcceptUEInitMsg(
      DeRegistrationAcceptUEInitMsg* deregistrationacceptueinit,
      uint8_t* buffer, uint32_t len);
};
}  // namespace magma5g

/*
DEREGISTRATION ACCEPT UE Initiated message content
IEI  Information Element             Type/Reference Presence   Format     Length
     Extended protocol discriminator Extended protocol discriminator 9.2 M V 1
     Security header type            Security header type            9.3 M V 1/2
     Spare half octet                Spare half                      9.5 M V 1/2
     De-registration request message Message type                    9.7 M V 1
*/
