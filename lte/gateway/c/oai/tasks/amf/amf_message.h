/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/*****************************************************************************

  Source      amf_message.h

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#include "amf_msgdef.h"
#include"registration_request.h"
#include"registration_accept.h"
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

using namespace std;
#pragma once
namespace magma5g
{
	
		class AMFMsg
		{
			public:
			AMFMsg();

			~AMFMsg();

			amf_msg_header header;

			registration_request_msg registrationrequestmsg;

			registration_accept_msg  registrationacceptmsg;

			registration_complete_msg registrationcompletemsg;

			registration_reject_msg registrationrejectmsg;

			identity_request_msg identityrequestmsg;

			identity_response_msg identityresponsemsg;

			authentication_request_msg authenticationrequestmsg;

			authentication_response_msg authenticationresponsemsg;

			authentication_reject_msg authenticationrejectmsg;

			authentication_failure_msg authenticationfailuremsg;

			security_mode_command_msg securitymodecommandmsg;

			security_mode_complete_msg securitymodecompletemsg;

			deregistration_request_msg deregistrationequesmsg;
			
			deregistration_accept_msg deregistrationacceptmsg;
			
			//SERVICE REQUEST
			int amf_msg_decode_header(amf_msg_header *header,const uint8_t *buffer, uint32_t len);
			
			int amf_msg_encode_header(const amf_msg_header *header, uint8_t *buffer,uint32_t len);

			int amf_msg_decode(amf_msg *msg, uint8_t *buffer, uint32_t len);

			int AmfMsgEncode(amf_msg *msg, uint8_t *buffer, uint32_t len);
		};

		
	

}
