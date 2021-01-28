
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

  Source      amf_message.cpp

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#ifdef __cplusplus
extern "C" {
#endif
#include "log.h"
#ifdef __cplusplus
}
#endif
#include "amf_fsm.h"
#include "amf_app_ue_context_and_proc.h"
#include "M5gNasMessage.h"
namespace magma5g {
// namespace NR_amf_msg
//{
int AMFMsg::amf_msg_decode(AMFMsg* msg, uint8_t* buffer, uint32_t len) {
  int header_result = 0;
  int decode_result = 0;
  header_result     = amf_msg_decode_header(&msg->header, buffer, len);
  if (header_result < 0) {
    // some error msg put in log file.
  }
  buffer += header_result;
  len -= header_result;
#if 0  // TODO -  NEED-RECHECK
    switch (msg->header.message_type) {
      case REGISTRATION_REQEST:
        decode_result = registration_request_msg.decode_registration_request(
            &msg->registrationrequestmsg, buffer, len);
        break;
      case IDENTITY_REQUEST:
        decode_result = identity_request_msg.decode_identity_request(
            &msg->identityrequestmsg, buffer, len);
        break;
      case IDENTITY_RESPONSE:
        decode_result = identity_response_msg.decode_identity_response(
            &msg->identityresponsemsg, buffer, len);
        break;
      case AUTHENTICATION_REQUEST:
        decode_result =
            authentication_request_msg.decode_authentication_request(
                &msg->authenticationrequestmsg, buffer, len);
        break;
      case AUTHENTICATION_RESPONSE:
        decode_result =
            authentication_response_msg.decode_authentication_response(
                &msg->authenticationresponsemsg, buffer, len);
        break;
      case AUTHENTICATION_REJECT:
        decode_result = authentication_reject_msg.decode_authentication_reject(
            &msg->authenticationrejectmsg, buffer, len);
        break;
      case AUTHENTICATION_FAILURE:
        decode_result =
            authentication_failure_msg.decode_authentication_failure(
                &msg->authenticationfailuremsg, buffer, len);
        break;
      case SECURITY_MODE_COMMAND:
        decode_result = security_mode_command_msg.decode_security_mode_command(
            &msg->securitymodecommandmsg, buffer, len);
        break;
      case SECURITY_MODE_COMPLETE:
        decode_result =
            security_mode_complete_msg.decode_security_mode_complete(
                &msg->securitymodecompletemsg, buffer, len);
        break;
      case REGISTRATION_ACCEPT:
        decode_result = registration_accept_msg.decode_registration_accept(
            &msg->registrationacceptmsg, buffer, len);
        break;
      case REGISTRATION_COMPLETE:
        decode_result = registration_complete_msg.decode_registration_complete(
            &msg->registrationcompletemsg, buffer, len);
        break;
      default:
        // error logs
    } //end switch
#endif
  if (decode_result < 0) {
    // some error logs
  }
}
//}  // namespace NR_amf_msg
}  // namespace magma5g
