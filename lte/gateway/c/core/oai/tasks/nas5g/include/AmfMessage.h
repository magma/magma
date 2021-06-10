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
#include "M5GRegistrationRequest.h"
#include "M5GRegistrationAccept.h"
#include "M5GRegistrationComplete.h"
#include "M5GRegistrationReject.h"
#include "M5GIdentityRequest.h"
#include "M5GIdentityResponse.h"
#include "M5GAuthenticationRequest.h"
#include "M5GAuthenticationResponse.h"
#include "M5GAuthenticationReject.h"
#include "M5GAuthenticationFailure.h"
#include "M5GSecurityModeCommand.h"
#include "M5GSecurityModeComplete.h"
#include "M5GSecurityModeReject.h"
#include "M5GDeRegistrationRequestUEInit.h"
#include "M5GDeRegistrationAcceptUEInit.h"
#include "M5GULNASTransport.h"
#include "M5GDLNASTransport.h"

using namespace std;
namespace magma5g {
// Amf NAS Msg Header
struct AmfMsgHeader_s {
  uint8_t extended_protocol_discriminator;
  uint8_t sec_header_type;
  uint8_t message_type;
};

union MMsg_u {
  RegistrationRequestMsg reg_request;
  RegistrationAcceptMsg reg_accept;
  RegistrationCompleteMsg reg_complete;
  RegistrationRejectMsg reg_reject;
  IdentityRequestMsg identity_request;
  IdentityResponseMsg identity_response;
  AuthenticationRequestMsg auth_request;
  AuthenticationResponseMsg auth_response;
  AuthenticationRejectMsg auth_reject;
  AuthenticationFailureMsg auth_failure;
  SecurityModeCommandMsg sec_mode_command;
  SecurityModeCompleteMsg sec_mode_complete;
  SecurityModeRejectMsg sec_mode_reject;
  DeRegistrationRequestUEInitMsg de_reg_request;
  DeRegistrationAcceptUEInitMsg de_reg_accept;
  ULNASTransportMsg ul_nas_transport;
  DLNASTransportMsg dl_nas_transport;
  MMsg_u();
  ~MMsg_u();
};

// Amf NAS Msg Class
class AmfMsg {
 public:
  AmfMsgHeader_s header;
  MMsg_u msg;

  AmfMsg();
  ~AmfMsg();
  int M5gNasMessageEncodeMsg(AmfMsg* msg, uint8_t* buffer, uint32_t len);
  int M5gNasMessageDecodeMsg(AmfMsg* msg, uint8_t* buffer, uint32_t len);
  int AmfMsgDecodeHeaderMsg(
      AmfMsgHeader_s* header, uint8_t* buffer, uint32_t len);
  int AmfMsgEncodeHeaderMsg(
      AmfMsgHeader_s* header, uint8_t* buffer, uint32_t len);
  int AmfMsgDecodeMsg(AmfMsg* msg, uint8_t* buffer, uint32_t len);
  int AmfMsgEncodeMsg(AmfMsg* msg, uint8_t* buffer, uint32_t len);
};
}  // namespace magma5g
