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
#include <iostream>
#include <sstream>
#include "AmfMessage.h"
#include "M5gNasMessage.h"
#include "CommonDefs.h"

using namespace std;
namespace magma5g
{
  AmfMsg::AmfMsg()
  {
  };

  AmfMsg::~AmfMsg()
  {
  };

  // Decode AMF NAS Header and Message 
  int AmfMsg::M5gNasMessageDecodeMsg(AmfMsg *msg, uint8_t *buffer, uint32_t len)
  {
    int headerresult = 0;
    int decoderesult = 0;

    if (len > 0 || buffer != NULL) {
      headerresult = msg->AmfMsgDecodeHeaderMsg(&msg->header, buffer, len);
      if (headerresult <= 0) { 
        MLOG(MERROR) << "   Error : Header Decoding Failed" << std::dec << RETURN_ERROR;
        return (RETURN_ERROR);
      }
    } else {
      MLOG(MERROR) << "Error : Buffer is Empty";
      return (RETURN_ERROR);
    }
    MLOG(MDEBUG) << "   epd = 0x" << hex << int(msg->header.extendedprotocoldiscriminator)  <<"\n"<< "   security hdr =  0x" << hex << int(msg->header.securityheadertype) <<"\n"<< "   hdr type = 0x" << hex << int(msg->header.messagetype)<<"\n";
    decoderesult = msg->AmfMsgDecodeMsg(msg, buffer, len);
    if (decoderesult <= 0) {
      MLOG(MERROR) << "decode result error ";
      return (RETURN_ERROR);
    }
    return (headerresult + decoderesult);       
  }

  // Encode AMF NAS  Header and Message
  int AmfMsg::M5gNasMessageEncodeMsg(AmfMsg *msg, uint8_t *buffer, uint32_t len)
  {
    int headerresult = 0;
    int encoderesult = 0;

    MLOG(MDEBUG) << "M5gNasMessageEncodeMsg:";
    if (len > 0 || buffer != NULL) {
      headerresult = msg->AmfMsgEncodeHeaderMsg(&msg->header, buffer, len);
      if (headerresult <= 0) {
        MLOG(MERROR) << "In M5gNasMessageEncodeMsg AmfMsgEncodeHeaderMsg ret error: "<< std::dec << RETURN_ERROR;
        return (RETURN_ERROR);
      }
    } else {
      MLOG(MERROR) << "Error : Buffer is empty "<<endl;
      return (RETURN_ERROR);
    }
    encoderesult = msg->AmfMsgEncodeMsg(msg, buffer, len);
    if (encoderesult <= 0) {
      MLOG(MERROR) << "Error : Encoding AMF Message Failed"<<endl;
      return (RETURN_ERROR);
    }
    return (headerresult + encoderesult);       
  }

  // Decode AMF Message Header
  int AmfMsg::AmfMsgDecodeHeaderMsg(AmfMsgHeader *hdr, uint8_t *buffer, uint32_t len)
  {
    int size = 0;

    MLOG(MDEBUG) << "AmfMsgDecodeHeaderMsg:"<<endl;
    if (len > 0  || buffer != NULL) {
      DECODE_U8(buffer + size, hdr->extendedprotocoldiscriminator, size);
      DECODE_U8(buffer + size, hdr->securityheadertype, size);
      DECODE_U8(buffer + size, hdr->messagetype, size);
      MLOG(MDEBUG) << "epd = 0x" << hex << int(hdr->extendedprotocoldiscriminator)  << "security hdr = 0x" << hex << int(hdr->securityheadertype) << " hdr type = 0x" << hex << int(hdr->messagetype);
    } else {
      MLOG(MERROR) << "Error : Buffer is Empty"<<endl;
      return(RETURN_ERROR);
    }

    if (hdr->extendedprotocoldiscriminator != M5G_MOBILITY_MANAGEMENT_MESSAGES) {
      MLOG(MERROR) << "Error : TLV not supported"<<endl;
      return (TLV_PROTOCOL_NOT_SUPPORTED);
    }
    return (size);
  }

  // Encode AMF Message Header
  int AmfMsg::AmfMsgEncodeHeaderMsg(AmfMsgHeader *hdr, uint8_t *buffer, uint32_t len)
  {
    int size = 0;

    MLOG(MDEBUG) << "AmfMsgEncodeHeaderMsg:";
    if (len > 0  || buffer != NULL) {
      ENCODE_U8(buffer + size, hdr->extendedprotocoldiscriminator, size);
      ENCODE_U8(buffer + size, hdr->securityheadertype, size);
      ENCODE_U8(buffer + size, hdr->messagetype, size);
      MLOG(MDEBUG) << "epd = 0x" << hex << int(hdr->extendedprotocoldiscriminator)  << "security hdr = 0x" << hex << int(hdr->securityheadertype) << "hdr type = 0x" << hex << int(hdr->messagetype);
    } else {
      MLOG(MERROR) << "Error : Buffer is Empty " ;
      return(RETURN_ERROR);
    }
    if ((unsigned char)hdr->extendedprotocoldiscriminator != M5G_MOBILITY_MANAGEMENT_MESSAGES) {
      MLOG(MERROR) << "Error : TLV not supported";
      return (TLV_PROTOCOL_NOT_SUPPORTED);
    }
    return (size);
  }

  // Decode AMF Message
  int AmfMsg::AmfMsgDecodeMsg(AmfMsg *msg, uint8_t *buffer, uint32_t len)
  {
    int decoderesult = 0;

    MLOG(MDEBUG) << "AmfMsgDecodeMsg:"<<endl;
    if (len <= 0 || buffer == NULL) 
    {
      MLOG(MERROR) << "Error : Buffer is Empty"<<endl;
      return(RETURN_ERROR);
    } 
    MLOG(MDEBUG) << "msg type = 0x" << hex << int(msg->header.messagetype);
    switch ((unsigned char)msg->header.messagetype) 
    {
      #ifdef HANDLE_POST_MVC 
      case REGISTRATION_REQUEST:
        MLOG(MDEBUG) << "Registraion request msg"<<endl;
        decoderesult = msg->registrationrequestmsg.DecodeRegistrationRequestMsg(&msg->registrationrequestmsg, buffer, len);
        break;
      case REGISTRATION_ACCEPT:
        MLOG(MDEBUG) << "AmfMsgDecodeMsg: Registraion accept msg"<<endl;
        decoderesult = msg->registrationacceptmsg.DecodeRegistrationAcceptMsg(&msg->registrationacceptmsg, buffer, len);
        break;
      case REGISTRATION_COMPLETE:
        decoderesult = msg->RegistrationCompleteMsg.DecodeRegistrationComplete(&msg->registrationcompletemsg, buffer, len);
        break;
      case REGISTRATION_REJECT:
        decoderesult=registration_complete_msg.decode_registration_reject(&msg->registrationrejectmsg,buffer,len);
        break;
      case IDENTITY_REQUEST:
        decoderesult=identity_request_msg.decode_identity_request(&msg->identityrequestmsg,buffer,len);
        break;
      case IDENTITY_RESPONSE:
        decoderesult =identity_response_msg.decode_identity_response(&msg->identityresponsemsg,buffer,len);
        break;
      case AUTHENTICATION_REQUEST:
        decoderesult=authentication_request_msg.decode_authentication_request(&msg->authenticationrequestmsg,buffer,len);
        break;
      case AUTHENTICATION_RESPONSE:
        decoderesult =authentication_response_msg.decode_authentication_response(&msg->authenticationresponsemsg,buffer,len);
        break;
      case AUTHENTICATION_REJECT:
        decoderesult =authentication_reject_msg.decode_authentication_reject(&msg->authenticationrejectmsg,buffer,len);
        break;
      case AUTHENTICATION_FAILURE:
        decoderesult =authentication_failure_msg.decode_authentication_failure(&msg->authenticationfailuremsg,buffer,len);
        break;
      case SECURITY_MODE_COMMAND:
        decoderesult =security_mode_command_msg.decode_security_mode_command(&msg->securitymodecommandmsg,buffer,len);
        break;
      case SECURITY_MODE_COMPLETE:
        decoderesult = security_mode_complete_msg.decode_security_mode_complete(&msg->securitymodecompletemsg,buffer,len);
        break;
      #endif
      default:
        decoderesult = TLV_WRONG_MESSAGE_TYPE; 
    }
    return (decoderesult);
  }

  // Encode AMF Message
  int AmfMsg::AmfMsgEncodeMsg(AmfMsg *msg, uint8_t *buffer, uint32_t len)
  {
    int encoderesult = 0;

    MLOG(MDEBUG) << " AmfMsgEncodeMsg : "<<endl;
    if (len <= 0 || buffer == NULL) 
    {
      MLOG(MERROR) << "Error : Buffer is Empty" ;
      return(RETURN_ERROR);
    } 
    switch ((unsigned char)msg->header.messagetype) 
    {
      #ifdef HANDLE_POST_MVC
      case REGISTRATION_REQUEST:
        encoderesult = msg->registrationrequestmsg.EncodeRegistrationRequestMsg(&msg->registrationrequestmsg, buffer, len);
        break;
  	  case REGISTRATION_ACCEPT:
        MLOG(MDEBUG) << " registraion accept msg \n";
        encoderesult = msg->registrationacceptmsg.EncodeRegistrationAcceptMsg(&msg->registrationacceptmsg, buffer, len);
        break;
  	  case REGISTRATION_COMPLETE:
        encoderesult = msg->RegistrationCompleteMsg.EncodeRegistrationComplete(&msg->registrationcompletemsg, buffer, len);
        break;
      case REGISTRATION_REJECT:
        encoderesult=registration_complete_msg.encode_registration_reject(&msg->registrationrejectmsg,buffer,len);
        break;
      case IDENTITY_REQUEST:
        encoderesult=identity_request_msg.encode_identity_request(&msg->identityrequestmsg,buffer,len);
        break;
      case IDENTITY_RESPONSE:
        encoderesult =identity_response_msg.encode_identity_response(&msg->identityresponsemsg,buffer,len);
        break;
      case AUTHENTICATION_REQUEST:
        encoderesult=authentication_request_msg.encode_authentication_request(&msg->authenticationrequestmsg,buffer,len);
        break;
      case AUTHENTICATION_RESPONSE:
        encoderesult =authentication_response_msg.encode_authentication_response(&msg->authenticationresponsemsg,buffer,len);
        break;
      case AUTHENTICATION_REJECT:
        encoderesult =authentication_reject_msg.encode_authentication_reject(&msg->authenticationrejectmsg,buffer,len);
        break;
      case AUTHENTICATION_FAILURE:
        encoderesult =authentication_failure_msg.encode_authentication_failure(&msg->authenticationfailuremsg,buffer,len);
        break;
      case SECURITY_MODE_COMMAND:
        encoderesult =security_mode_command_msg.encode_security_mode_command(&msg->securitymodecommandmsg,buffer,len);
        break;
      case SECURITY_MODE_COMPLETE:
        encoderesult = security_mode_complete_msg.encode_security_mode_complete(&msg->securitymodecompletemsg,buffer,len);
        break;
      default:
        encoderesult = TLV_WRONG_MESSAGE_TYPE; 
      #endif
    }
    return (encoderesult);
  }
}//namespace magma5g
