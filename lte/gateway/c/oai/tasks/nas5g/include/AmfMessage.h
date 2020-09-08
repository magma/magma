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
#include "RegistrationRequest.h"
#include "RegistrationAccept.h"
#if 0 // TBD
#include"registration_complete.h"
#include"registration_reject.h"
#include"identity_request.h"
#include"identity_response.h"
#include"Authentication_Request.h"
#include"Authentication_Response.h"
#include"Authentication_Reject.h"
#include"AuthenticationFailure.h"
#include "security_mode_command.h"
#include "security_mode_complete.h"
#include"deregistration_request.h"
#include"deregistration_accept.h"
#endif
using namespace std;

namespace magma5g
{
  // Amf NAS Msg Header Class
  class AmfMsgHeader
  {
    public:
      uint8_t extendedprotocoldiscriminator;
      uint8_t securityheadertype;
      uint8_t messagetype;
  };

  // Amf NAS Msg Class
  class AmfMsg
  {
    public:
      AmfMsg();
      ~AmfMsg();
      AmfMsgHeader header;
      RegistrationRequestMsg registrationrequestmsg;
      RegistrationAcceptMsg  registrationacceptmsg;
#if 0 // TBD
      RegistrationCompleteMsg registrationcompletemsg;
      RegistrationRejectMsg registrationrejectmsg;
      IdentityRequestMsg identityrequestmsg;
      IdentityResponseMsg identityresponsemsg;
      AuthenticationRequestMsg authenticationrequestmsg;
      AuthenticationResponseMsg authenticationresponsemsg;
      AuthenticationRejectMsg authenticationrejectmsg;
      authenticationFailureMsg authenticationfailuremsg;
      SecurityModeCommandMsg securitymodecommandmsg;
      SecurityModeCompleteMsg securitymodecompletemsg;
      DeregistrationRequestMsg deregistrationequesmsg;
      DregistrationAcceptMsg deregistrationacceptmsg;
#endif
      int M5gNasMessageEncodeMsg (AmfMsg *msg, uint8_t *buffer, uint32_t len);
      int M5gNasMessageDecodeMsg (AmfMsg *msg, uint8_t *buffer, uint32_t len);
      int AmfMsgDecodeHeaderMsg (AmfMsgHeader *header, uint8_t *buffer, uint32_t len);
      int AmfMsgEncodeHeaderMsg (AmfMsgHeader *header, uint8_t *buffer,uint32_t len);
      int AmfMsgDecodeMsg (AmfMsg *msg, uint8_t *buffer, uint32_t len);
      int AmfMsgEncodeMsg (AmfMsg *msg, uint8_t *buffer, uint32_t len);
  };
}
