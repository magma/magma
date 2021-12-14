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

#include <sstream>
#include <thread>
#include <memory>
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"
#include "lte/gateway/c/core/oai/lib/secu/secu_defs.h"
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/common/common_defs.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5gNasMessage.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_defs.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_authentication.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_as.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_fsm.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_recv.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GDLNASTransport.h"
#include "lte/gateway/c/core/oai/lib/s6a_proxy/S6aClient.h"
#include "lte/gateway/c/core/oai/tasks/grpc_service/proto_msg_to_itti_msg.h"
#include "lte/gateway/c/core/oai/include/ngap_messages_types.h"
#include "lte/gateway/c/core/oai/lib/n11/M5GAuthenticationServiceClient.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GMMCause.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_38.401.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_common.h"
using magma5g::AsyncM5GAuthenticationServiceClient;

using namespace magma;

#define AMF_CAUSE_SUCCESS (1)
namespace magma5g {
/*forward declaration*/
extern task_zmq_ctx_t amf_app_task_zmq_ctx;
static int amf_as_establish_req(amf_as_establish_t* msg, int* amf_cause);
static int amf_as_security_req(
    const amf_as_security_t* msg, m5g_dl_info_transfer_req_t* as_msg);
static int amf_as_security_rej(
    const amf_as_security_t* msg, m5g_dl_info_transfer_req_t* as_msg);
static int amf_as_establish_rej(
    const amf_as_establish_t* msg, nas5g_establish_rsp_t* amf_msg);
// Setup the security header of the given NAS message
static AMFMsg* amf_as_set_header(
    amf_nas_message_t* msg, const amf_as_security_data_t* security);
static int amf_service_rejectmsg(
    const amf_as_establish_t* msg, ServiceRejectMsg* service_reject);
/**************************************************************************
**                                                                       **
** Name        : amf_as_send()                                           **
**                                                                       **
** Description : Processes the AMF-AS Service Access Point primitive.    **
**                                                                       **
** Inputs      : msg    :  The AMF-AS-SAP primitive to process           **
**               Others :  None                                          **
**                                                                       **
** Outputs     : None                                                    **
**      Return : RETURNok, RETURNerror                                   **
**      Others : None                                                    **
**                                                                       **
**************************************************************************/
int amf_as_send(amf_as_t* msg) {
  int rc                       = RETURNok;
  int amf_cause                = AMF_CAUSE_SUCCESS;
  amf_as_primitive_t primitive = msg->primitive;

  switch (primitive) {
    case _AMFAS_DATA_IND:
    case _AMFAS_ESTABLISH_REQ:
      // Process UE's establishment request
      rc = amf_as_establish_req(&msg->u.establish, &amf_cause);
      break;
    case _AMFAS_RELEASE_IND:
    default:
      // Other primitives are forwarded to NGAP
      rc = amf_as_send_ng(msg);
      break;
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/***************************************************************************
  uint8_t* lenPtr;
**                                                                        **
** Name:    amf_as_establish_req()                                        **
**                                                                        **
** Description: Processes the _AMFAS-SAP connection establish request      **
**      primitive                                                         **
**                                                                        **
** _AMFAS-SAP - AS->AMF : ESTABLISH_REQ - NAS signaling connection        **
**     The AS notifies the NAS that establishment of the signal-          **
**     ling connection has been requested to tranfer initial NAS          **
**     message from the UE.                                               **
**                                                                        **
**      Inputs:  msg:       The _AMFAS-SAP primitive to process            **
**      Others:    None                                                   **
**                                                                        **
**      Outputs:   amf_cause: AMF cause code                              **
**      Return:    RETURNok, RETURNerror                                  **
**      Others:    None                                                   **
**                                                                        **
***************************************************************************/
static int amf_as_establish_req(amf_as_establish_t* msg, int* amf_cause) {
  amf_security_context_t* amf_security_context = NULL;
  amf_nas_message_decode_status_t decode_status;
  memset(&decode_status, 0, sizeof(decode_status));
  int decoder_rc                       = 1;
  int rc                               = RETURNerror;
  tai_t originating_tai                = {0};
  amf_nas_message_t nas_msg            = {0};
  ue_m5gmm_context_s* ue_m5gmm_context = NULL;
  ue_m5gmm_context = amf_ue_context_exists_amf_ue_ngap_id(msg->ue_id);
  if (ue_m5gmm_context == NULL) {
    OAILOG_ERROR(
        LOG_AMF_APP,
        "ue context not found for the ue_id=" AMF_UE_NGAP_ID_FMT "\n",
        msg->ue_id);
    OAILOG_FUNC_RETURN(LOG_AMF_APP, rc);
  }

  amf_context_t* amf_ctx = NULL;
  amf_ctx                = &ue_m5gmm_context->amf_context;

  if (amf_ctx) {
    if (IS_AMF_CTXT_PRESENT_SECURITY(amf_ctx)) {
      amf_security_context = &amf_ctx->_security;
    }
  }

  if ((msg->nas_msg->data[1] != 0x0) && (msg->nas_msg->data[9] == 0x5c)) {
    for (int i = 0, j = 7; j < blength(msg->nas_msg); i++, j++) {
      msg->nas_msg->data[i] = msg->nas_msg->data[j];
    }
    msg->nas_msg->slen = msg->nas_msg->slen - 7;
  }

  // Decode initial NAS message
  decoder_rc = nas5g_message_decode(
      msg->nas_msg->data, &nas_msg, blength(msg->nas_msg), amf_security_context,
      &decode_status);

  bdestroy_wrapper(&msg->nas_msg);

  // conditional IE error
  if (decoder_rc < 0) {
    if (decoder_rc < TLV_FATAL_ERROR) {
      *amf_cause = AMF_CAUSE_PROTOCOL_ERROR;
    } else if (decoder_rc == TLV_MANDATORY_FIELD_NOT_PRESENT) {
      *amf_cause = AMF_CAUSE_INVALID_MANDATORY_INFO;
    } else if (decoder_rc == TLV_UNEXPECTED_IEI) {
      *amf_cause = AMF_CAUSE_IE_NOT_IMPLEMENTED;
    } else {
      *amf_cause = AMF_CAUSE_PROTOCOL_ERROR;
    }
  }

  // Process initial NAS message
  AMFMsg* amf_msg = &nas_msg.plain.amf;

  switch (amf_msg->header.message_type) {
    case REG_REQUEST:
      memcpy(&originating_tai, &msg->tai, sizeof(originating_tai));
      rc = amf_handle_registration_request(
          msg->ue_id, &originating_tai, &msg->ecgi,
          &amf_msg->msg.registrationrequestmsg, msg->is_initial,
          msg->is_amf_ctx_new, *amf_cause, decode_status);
      break;
    case M5G_SERVICE_REQUEST:  // SERVICE_REQUEST:
      rc = amf_handle_service_request(
          msg->ue_id, &amf_msg->msg.service_request, decode_status);
      break;
    case M5G_IDENTITY_RESPONSE:
      rc = amf_handle_identity_response(
          msg->ue_id, &amf_msg->msg.identityresponsemsg.m5gs_mobile_identity,
          *amf_cause, decode_status);
      break;
    case AUTH_RESPONSE:
      rc = amf_handle_authentication_response(
          msg->ue_id, &amf_msg->msg.authenticationresponsemsg, *amf_cause,
          decode_status);
      break;
    case AUTH_FAILURE:
      rc = amf_handle_authentication_failure(
          msg->ue_id, &amf_msg->msg.authenticationfailuremsg, *amf_cause,
          decode_status);
      break;
    case SEC_MODE_COMPLETE:
      rc = amf_handle_security_complete_response(msg->ue_id, decode_status);
      break;
    case SEC_MODE_REJECT:
      rc = amf_handle_security_mode_reject(
          msg->ue_id, &amf_msg->msg.securitymodereject, *amf_cause,
          decode_status);
      break;
    case REG_COMPLETE:
      rc = amf_handle_registration_complete_response(
          msg->ue_id, &amf_msg->msg.registrationcompletemsg, *amf_cause,
          decode_status);
      break;
    case DE_REG_REQUEST_UE_ORIGIN:
      rc = amf_handle_deregistration_ue_origin_req(
          msg->ue_id, &amf_msg->msg.deregistrationequesmsg, *amf_cause,
          decode_status);
      break;
    case ULNASTRANSPORT:
      rc = amf_smf_process_pdu_session_packet(
          msg->ue_id, &amf_msg->msg.uplinknas5gtransport, *amf_cause);
      break;
    default:
      OAILOG_DEBUG(
          LOG_NAS_AMF, "Unknown message type: %d, in %s ",
          amf_msg->header.message_type, __FUNCTION__);
  }
  return rc;
}

/**************************************************************************
 **                                                                      **
 ** Name       : amf_as_send_ng()                                        **
 **                                                                      **
 ** Description: Builds NAS message according to the given _AMFAS Service **
 **      Access Point primitive and sends it to the Access Stratum       **
 **      sublayer                                                        **
 **                                                                      **
 ** Inputs     : msg: The _AMFAS-SAP primitive to be sent                 **
 **      Others: None                                                    **
 **                                                                      **
 ** Outputs:     None                                                    **
 **      Return: RETURNok, RETURNerror                                   **
 **      Others: None                                                    **
 **                                                                      **
 *************************************************************************/
int amf_as_send_ng(const amf_as_t* msg) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  amf_as_message_t as_msg = {0};

  switch (msg->primitive) {
    case _AMFAS_DATA_REQ:
      as_msg.msg_id =
          amf_as_data_req(&msg->u.data, &as_msg.msg.dl_info_transfer_req);
      break;
    case _AMFAS_ESTABLISH_CNF:
      as_msg.msg_id = amf_as_establish_cnf(
          &msg->u.establish, &as_msg.msg.nas_establish_rsp);
      break;
    case _AMFAS_SECURITY_REQ:
      as_msg.msg_id = amf_as_security_req(
          &msg->u.security, &as_msg.msg.dl_info_transfer_req);
      break;
    case _AMFAS_SECURITY_REJ:
      as_msg.msg_id = amf_as_security_rej(
          &msg->u.security, &as_msg.msg.dl_info_transfer_req);
      break;
    case _AMFAS_ESTABLISH_REJ:
      as_msg.msg_id = amf_as_establish_rej(
          &msg->u.establish, &as_msg.msg.nas_establish_rsp);
      break;
    default:
      as_msg.msg_id = 0;
      break;
  }

  /*
   * Send the message to the Access Stratum or NGAP in case of AMF
   */
  if (as_msg.msg_id > 0) {
    switch (as_msg.msg_id) {
      case AS_DL_INFO_TRANSFER_REQ_: {
        amf_app_handle_nas_dl_req(
            as_msg.msg.dl_info_transfer_req.ue_id,
            as_msg.msg.dl_info_transfer_req.nas_msg,
            as_msg.msg.dl_info_transfer_req.err_code);
        OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
      } break;
      case AS_NAS_ESTABLISH_RSP_: {
        amf_app_handle_nas_dl_req(
            as_msg.msg.nas_establish_rsp.ue_id,
            as_msg.msg.nas_establish_rsp.nas_msg,
            as_msg.msg.nas_establish_rsp.err_code);
        as_msg.msg.nas_establish_rsp.nas_msg = NULL;
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
      } break;
      case AS_NAS_RELEASE_REQ_:
        OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
        break;
      case AS_NAS_ESTABLISH_CNF_:
        OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
        break;
      default:
        break;
    }
  }

  /* ICS goes as a separate message.
   * Return ok here
   */
  if (msg->primitive == _AMFAS_ESTABLISH_CNF) {
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
}

/************************************************************************
 **                                                                    **
 ** Name       : amf_as_encode()                                       **
 **                                                                    **
 ** Description: Encodes NAS message into NAS information container    **
 **                                                                    **
 ** Inputs     : msg : The NAS message to encode                       **
 **      length: The maximum length of the NAS message                 **
 **      Others: None                                                  **
 **                                                                    **
 ** Outputs    : info : The NAS information container                  **
 **      msg   : The NAS message to encode                             **
 **      Return: The number of bytes successfully encoded              **
 **      Others: None                                                  **
 **                                                                    **
 ***********************************************************************/
static int amf_as_encode(
    bstring* info, amf_nas_message_t* msg, size_t length,
    amf_security_context_t* amf_security_context) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int bytes = 1;

  /* Ciphering algorithms, EA1 and EA2 expects length to be mode of 4,
   * so length is modified such that it will be mode of 4 */
  AMF_GET_BYTE_ALIGNED_LENGTH(length);
  if (msg->header.security_header_type != SECURITY_HEADER_TYPE_NOT_PROTECTED) {
    amf_msg_header* header = &msg->security_protected.plain.amf.header;
    // Expand size of protected NAS message
    length += NAS_MESSAGE_SECURITY_HEADER_SIZE;
    // Set header of plain NAS message
    header->extended_protocol_discriminator = M5GS_MOBILITY_MANAGEMENT_MESSAGE;
    header->security_header_type = SECURITY_HEADER_TYPE_NOT_PROTECTED;
  }

  // Allocate memory to the NAS information container
  *info = bfromcstralloc(length, "\0");

  if (*info) {
    // Encode the NAS message
    bytes =
        nas5g_message_encode((*info)->data, msg, length, amf_security_context);

    if (bytes > 0) {
      (*info)->slen = bytes;
    } else {
      bdestroy_wrapper(info);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, bytes);
}

/****************************************************************************
 **                                                                        **
 ** Name:        amf_reg_acceptmsg()                                       **
 **                                                                        **
 ** Description: Builds Registration accept message                        **
 **                                                                        **
 **              The Registration Accept message is sent by the            **
 **              network to the UEi.                                       **
 **                                                                        **
 ** Inputs:      msg:           The AMFMsg    primitive to process         **
 **              Others:        None                                       **
 **                                                                        **
 ** Outputs:     amf_msg:       The AMF message to be sent                 **
 **              Return:        The size of the AMF message                **
 **              Others:        None                                       **
 **                                                                        **
 ***************************************************************************/
int amf_reg_acceptmsg(const guti_m5_t* guti, amf_nas_message_t* nas_msg) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int size = REGISTRATION_ACCEPT_MINIMUM_LENGTH;
  nas_msg->security_protected.plain.amf.header.message_type = REG_ACCEPT;
  nas_msg->security_protected.plain.amf.header.extended_protocol_discriminator =
      M5G_MOBILITY_MANAGEMENT_MESSAGES;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg
      .extended_protocol_discriminator.extended_proto_discriminator =
      M5G_MOBILITY_MANAGEMENT_MESSAGES;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg
      .sec_header_type.sec_hdr = SECURITY_HEADER_TYPE_NOT_PROTECTED;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.message_type
      .msg_type = REG_ACCEPT;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg
      .m5gs_reg_result.sms_allowed = NOT_ALLOWED;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg
      .m5gs_reg_result.reg_result_val = M3GPP_ACCESS;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.mobile_id
      .mobile_identity.guti.odd_even = EVEN_IDENTITY_DIGITS;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.mobile_id
      .iei = M5GS_MOBILE_IDENTITY;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.mobile_id
      .len = M5GSMobileIdentityMsg_GUTI_LENGTH;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.mobile_id
      .mobile_identity.guti.type_of_identity = M5GSMobileIdentityMsg_GUTI;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.mobile_id
      .mobile_identity.guti.mcc_digit1 = guti->guamfi.plmn.mcc_digit1;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.mobile_id
      .mobile_identity.guti.mcc_digit2 = guti->guamfi.plmn.mcc_digit2;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.mobile_id
      .mobile_identity.guti.mcc_digit3 = guti->guamfi.plmn.mcc_digit3;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.mobile_id
      .mobile_identity.guti.mnc_digit1 = guti->guamfi.plmn.mnc_digit1;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.mobile_id
      .mobile_identity.guti.mnc_digit2 = guti->guamfi.plmn.mnc_digit2;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.mobile_id
      .mobile_identity.guti.mnc_digit3 = guti->guamfi.plmn.mnc_digit3;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.mobile_id
      .mobile_identity.guti.amf_regionid = guti->guamfi.amf_regionid;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.mobile_id
      .mobile_identity.guti.amf_setid = guti->guamfi.amf_set_id;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.mobile_id
      .mobile_identity.guti.amf_pointer = guti->guamfi.amf_pointer;

  uint8_t* offset;
  uint32_t encoded_tmsi = ntohl(guti->m_tmsi);

  offset = reinterpret_cast<uint8_t*>(&encoded_tmsi);

  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.mobile_id
      .mobile_identity.guti.tmsi1 = *offset;
  offset++;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.mobile_id
      .mobile_identity.guti.tmsi2 = *offset;
  offset++;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.mobile_id
      .mobile_identity.guti.tmsi3 = *offset;
  offset++;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.mobile_id
      .mobile_identity.guti.tmsi4 = *offset;

  // TAI List, Allowed NSSAI and GPRS Timer 3 hardcoded
  nas_msg->header.security_header_type =
      SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED;  // sit_change
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.tai_list.iei =
      0x54;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.tai_list.len =
      0x7;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.tai_list
      .list_type = 0x0;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.tai_list
      .num_elements =
      0x0;  // 0 implies 1 as per ts_124501v1506 section 9.11.3.9
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.tai_list
      .mcc_digit1 = guti->guamfi.plmn.mcc_digit1;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.tai_list
      .mcc_digit2 = guti->guamfi.plmn.mcc_digit2;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.tai_list
      .mcc_digit3 = guti->guamfi.plmn.mcc_digit3;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.tai_list
      .mnc_digit1 = guti->guamfi.plmn.mnc_digit1;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.tai_list
      .mnc_digit2 = guti->guamfi.plmn.mnc_digit2;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.tai_list
      .mnc_digit3 = guti->guamfi.plmn.mnc_digit3;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.tai_list
      .tac[0] = 0x00;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.tai_list
      .tac[1] = 0x00;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.tai_list
      .tac[2] = 0x01;

  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.allowed_nssai
      .iei = ALLOWED_NSSAI;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.allowed_nssai
      .len = 2;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.allowed_nssai
      .nssai.len = 1;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.allowed_nssai
      .nssai.sst = 0x01;

  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.gprs_timer
      .len = 1;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.gprs_timer
      .unit = 0;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg.gprs_timer
      .timervalue = 6;

  size += MOBILE_IDENTITY_MAX_LENGTH;
  size += 20;
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, size);
}

/****************************************************************************
 **                                                                        **
 ** Name:        amf_service_acceptmsg()                                   **
 **                                                                        **
 ** Description: Builds Service accept message                             **
 **                                                                        **
 **              The Service Accept message is sent by the                 **
 **              network to the UEi.                                       **
 **                                                                        **
 ** Inputs:      msg:           The AMFMsg    primitive to process         **
 **              Others:        None                                       **
 **                                                                        **
 ** Outputs:     amf_msg:       The AMF message to be sent                 **
 **              Return:        The size of the AMF message                **
 **              Others:        None                                       **
 **                                                                        **
 ***************************************************************************/
int amf_service_acceptmsg(
    const amf_as_establish_t* msg, amf_nas_message_t* nas_msg) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int size = SERVICE_ACCEPT_MINIMUM_LENGTH;
  nas_msg->security_protected.plain.amf.header.message_type =
      M5G_SERVICE_ACCEPT;
  nas_msg->security_protected.plain.amf.header.extended_protocol_discriminator =
      M5G_MOBILITY_MANAGEMENT_MESSAGES;
  nas_msg->security_protected.plain.amf.msg.registrationacceptmsg
      .extended_protocol_discriminator.extended_proto_discriminator =
      M5G_MOBILITY_MANAGEMENT_MESSAGES;
  nas_msg->security_protected.plain.amf.msg.service_accept.sec_header_type
      .sec_hdr = SECURITY_HEADER_TYPE_NOT_PROTECTED;
  nas_msg->security_protected.plain.amf.msg.service_accept.message_type
      .msg_type = M5G_SERVICE_ACCEPT;
  nas_msg->header.security_header_type =
      SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED;

  if (msg->pdu_session_status_ie & AMF_AS_PDU_SESSION_STATUS) {
    nas_msg->security_protected.plain.amf.msg.service_accept.pdu_session_status
        .iei = PDU_SESSION_STATUS;
    nas_msg->security_protected.plain.amf.msg.service_accept.pdu_session_status
        .len = 0x02;
    nas_msg->security_protected.plain.amf.msg.service_accept.pdu_session_status
        .pduSessionStatus = msg->pdu_session_status;
  }

  if (msg->pdu_session_status_ie & AMF_AS_PDU_SESSION_REACTIVATION_STATUS) {
    nas_msg->security_protected.plain.amf.msg.service_accept
        .pdu_re_activation_status.iei = PDU_SESSION_REACTIVATION_RESULT;
    nas_msg->security_protected.plain.amf.msg.service_accept
        .pdu_re_activation_status.len = 0x02;

    nas_msg->security_protected.plain.amf.msg.service_accept
        .pdu_re_activation_status.pduSessionReActivationResult =
        msg->pdu_session_reactivation_status;
  }

  nas_msg->security_protected.header.message_type = M5G_SERVICE_ACCEPT;
  size += NAS5G_MESSAGE_CONTAINER_MAXIMUM_LENGTH;
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, size);
}
/****************************************************************************
 **                                                                        **
 ** Name:        amf_dl_nas_transport_msg()                            **
 **                                                                        **
 ** Description: Builds Downlink Nas Transport message                     **
 **                                                                        **
 **              The Downlink Nas Transport message is sent by the         **
 **              network to the UE to transfer the data in DL              **
 **              This function is used to send DL NAS Transport message    **
 **              via S1AP DL NAS Transport message.                        **
 **                                                                        **
 ** Inputs:      msg:           The AMFMsg    primitive to process         **
 **              Others:        None                                       **
 **                                                                        **
 ** Outputs:     amf_msg:       The AMF message to be sent                 **
 **              Return:        The size of the AMF message                **
 **              Others:        None                                       **
 **                                                                        **
 ***************************************************************************/
static int amf_dl_nas_transport_msg(
    const amf_as_data_t* msg, DLNASTransportMsg* amf_msg) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int size = AMF_HEADER_LENGTH;
  // Mandatory - Message type
  amf_msg->message_type.msg_type = DOWNLINK_NAS_TRANSPORT;
  // Mandatory - Nas message container
  size += NAS5G_MESSAGE_CONTAINER_MAXIMUM_LENGTH;
  memcpy(
      amf_msg->payload_container.contents, &(msg->nas_msg),
      sizeof(msg->nas_msg));
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, size);
}

/****************************************************************************
 **                                                                        **
 ** Name:        amf_de_reg_acceptmsg()                                    **
 **                                                                        **
 ** Description: Builds De-Registration Accept message                     **
 **                                                                        **
 **              The De-Registration accept message is sent by the network **
 **                                                                        **
 ** Inputs:      msg:           The AMFMsg    primitive to process         **
 **              Others:        None                                       **
 **                                                                        **
 ** Outputs:     amf_msg:       The AMF message to be sent                 **
 **              Return:        The size of the AMF message                **
 **              Others:        None                                       **
 **                                                                        **
 ***************************************************************************/
static int amf_de_reg_acceptmsg(
    const amf_as_data_t* msg, m5g_dl_info_transfer_req_t* as_msg,
    amf_nas_message_t* nas_msg, DeRegistrationAcceptUEInitMsg* amf_msg) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int size = AMF_HEADER_LENGTH;
  ue_m5gmm_context_s* ue_context;
  uint8_t seq_no = 0;

  ue_context = amf_ue_context_exists_amf_ue_ngap_id(as_msg->ue_id);

  if (ue_context) {
    seq_no = ue_context->amf_context._security.dl_count.seq_num;
  }

  nas_msg->security_protected.plain.amf.header.extended_protocol_discriminator =
      M5G_MOBILITY_MANAGEMENT_MESSAGES;
  nas_msg->security_protected.plain.amf.header.message_type =
      DE_REG_ACCEPT_UE_ORIGIN;
  nas_msg->header.security_header_type =
      SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED;
  nas_msg->header.extended_protocol_discriminator =
      M5G_MOBILITY_MANAGEMENT_MESSAGES;
  nas_msg->header.sequence_number = seq_no;

  // Mandatory - Message type
  amf_msg->extended_protocol_discriminator.extended_proto_discriminator =
      M5G_MOBILITY_MANAGEMENT_MESSAGES;
  amf_msg->spare_half_octet.spare  = 0x00;
  amf_msg->sec_header_type.sec_hdr = SECURITY_HEADER_TYPE_NOT_PROTECTED;
  amf_msg->message_type.msg_type   = DE_REG_ACCEPT_UE_ORIGIN;

  size += NAS5G_MESSAGE_CONTAINER_MAXIMUM_LENGTH;
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, size);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_as_data_req()                                             **
 **                                                                        **
 ** Description: Processes the AMFAS-SAP data transfer request             **
 **      primitive                                                         **
 **                                                                        **
 ** AMFAS-SAP - AMF->AS : DATA_REQ - Data transfer procedure               **
 **                                                                        **
 ** Inputs:  msg:       The AMFAS-SAP primitive to process                 **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     as_msg:    The message to send to the AS                  **
 **      Return:    The identifier of the AS message                       **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/
uint16_t amf_as_data_req(
    const amf_as_data_t* msg, m5g_dl_info_transfer_req_t* as_msg) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int size                                    = 0;
  int is_encoded                              = false;
  amf_nas_message_t nas_msg                   = {0};
  nas_msg.security_protected.header           = {0};
  nas_msg.security_protected.plain.amf.header = {0};
  nas_msg.security_protected.plain.amf.header = {0};

  // Setup the AS message
  as_msg->ue_id = msg->ue_id;
  if (msg->guti) {
    as_msg->s_tmsi.amf_set_id = msg->guti->guamfi.amf_set_id;
    as_msg->s_tmsi.m_tmsi     = msg->guti->m_tmsi;
  }

  // Setup the NAS security header
  AMFMsg* amf_msg = amf_as_set_header(&nas_msg, &msg->sctx);

  // Setup the NAS information message
  if (amf_msg) {
    switch (msg->nas_info) {
      case AMF_AS_NAS_DATA_REGISTRATION_ACCEPT:
        size = amf_reg_acceptmsg(msg->guti, &nas_msg);

        if (msg->guti) {
          delete (msg->guti);
        }

        break;
      case AMF_AS_NAS_DL_NAS_TRANSPORT:
        // DL messages to NGAP on Identity/Authentication request
        size =
            amf_dl_nas_transport_msg(msg, &amf_msg->msg.downlinknas5gtransport);
        break;
      case AMF_AS_NAS_DATA_DEREGISTRATION_ACCEPT: {
        size = amf_de_reg_acceptmsg(
            msg, as_msg, &nas_msg, &amf_msg->msg.deregistrationacceptmsg);
      } break;
      default:
        // Send other NAS messages as already encoded SMF messages
        size = msg->nas_msg.length();
        break;
    }
  }

  if (size > 0) {
    int bytes                                    = 0;
    amf_security_context_t* amf_security_context = NULL;
    amf_context_t* amf_ctx                       = NULL;
    ue_m5gmm_context_s* ue_m5gmm_context =
        amf_ue_context_exists_amf_ue_ngap_id(msg->ue_id);

    if (ue_m5gmm_context) {
      amf_ctx = &ue_m5gmm_context->amf_context;
      if (amf_ctx) {
        if (IS_AMF_CTXT_PRESENT_SECURITY(amf_ctx)) {
          amf_security_context           = &amf_ctx->_security;
          nas_msg.header.sequence_number = amf_ctx->_security.dl_count.seq_num;
        }
      }
    } else {
      OAILOG_ERROR(
          LOG_AMF_APP,
          "ue context not found for the ue_id=" AMF_UE_NGAP_ID_FMT "\n",
          msg->ue_id);
      OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
    }

    if (!is_encoded) {
      /*
       * Encode the NAS information message
       */
      bytes =
          amf_as_encode(&as_msg->nas_msg, &nas_msg, size, amf_security_context);
    }

    if (bytes > 0) {
      as_msg->err_code = M5G_AS_SUCCESS;
    } else {
      OAILOG_ERROR(LOG_AMF_APP, "NAS encoding failed bytes=%d\n", bytes);
    }

    OAILOG_FUNC_RETURN(LOG_NAS_AMF, AS_DL_INFO_TRANSFER_REQ_);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, 0);
}

/***************************************************************************
 **                                                                       **
 ** Name:        amf_as_set_header()                                      **
 **                                                                       **
 ** Description: Setup the security header of the given NAS message       **
 **                                                                       **
 ** Inputs:      security: The NAS security data to use                   **
 **              Others:   None                                           **
 **                                                                       **
 ** Outputs:     msg:     The NAS message                                 **
 **              Return:  Pointer to the plain NAS message to be se-      **
 **                       curity protected if setting of the securi-      **
 **                       ty header succeed;                              **
 **                       NULL pointer otherwise                          **
 **              Others:  None                                            **
 **                                                                       **
 **************************************************************************/
AMFMsg* amf_as_set_header(
    amf_nas_message_t* msg, const amf_as_security_data_t* security) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  msg->header.extended_protocol_discriminator =
      M5GS_MOBILITY_MANAGEMENT_MESSAGE;

  if (security && (security->ksi != KSI_NO_KEY_AVAILABLE)) {
    /*
     * A valid 5G CN security context exists
     */
    if (security->is_new) {
      /*
       * New 5G CN security context is taken into use
       */
      if (security->is_knas_int_present) {
        if (security->is_knas_enc_present) {
          /*
           * NAS integrity and cyphering keys are available
           */
          msg->header.security_header_type =
              SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED_NEW;
        } else {
          /*
           * NAS integrity key only is available
           */
          msg->header.security_header_type =
              SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_NEW;
        }
        OAILOG_FUNC_RETURN(LOG_NAS_AMF, &msg->security_protected.plain.amf);
      }
    } else if (security->is_knas_int_present) {
      if (security->is_knas_enc_present) {
        /*
         * NAS integrity and cyphering keys are available
         */
        msg->header.security_header_type =
            SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED;
      } else {
        /*
         * NAS integrity key only is available
         */
        msg->header.security_header_type =
            SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED;
      }
      OAILOG_FUNC_RETURN(LOG_NAS_AMF, &msg->security_protected.plain.amf);
    } else {
      /*
       * No valid 5G CN security context exists
       */
      msg->header.security_header_type = SECURITY_HEADER_TYPE_NOT_PROTECTED;
      OAILOG_FUNC_RETURN(LOG_NAS_AMF, &msg->plain.amf);
    }
  } else {
    /*
     * No valid 5G CN security context exists
     */
    msg->header.security_header_type = SECURITY_HEADER_TYPE_NOT_PROTECTED;
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, &msg->plain.amf);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, NULL);
}

/***************************************************************************
 **                                                                       **
 ** Name:        amf_identity_request()                              **
 **                                                                       **
 ** Description: Send Identity Request to UE                              **
 **                                                                       **
 ** Inputs:      msg: Security msg                                        **
 **              amf_msg :   amf msg                                      **
 **                                                                       **
 ** Return:   size                                                        **
 **                                                                       **
 **************************************************************************/
static int amf_identity_request(
    const amf_as_security_t* msg, IdentityRequestMsg* amf_msg) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int size = AMF_HEADER_LENGTH;
  /*
   * Mandatory - Message type
   */
  amf_msg->message_type.msg_type  = IDENTITY_REQUEST;
  amf_msg->m5gs_identity_type.toi = M5G_IDENTITY_SUCI;
  size += IDENTITY_TYPE_2_IE_MAX_LENGTH;
  if (msg->ident_type == IDENTITY_TYPE_2_IMSI) {
    amf_msg->m5gs_identity_type.toi = IDENTITY_TYPE_2_IMSI;
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, size);
}

/***************************************************************************
 **                                                                       **
 ** Name:        amf_auth_reject()                                        **
 **                                                                       **
 ** Description: Send authentication Reject to UE                         **
 **                                                                       **
 ** Inputs:      msg: Security msg                                        **
 **              amf_msg :   amf msg                                      **
 **                                                                       **
 ** Return:   size                                                        **
 **                                                                       **
 **************************************************************************/
static int amf_auth_reject(
    const amf_as_security_t* msg, AuthenticationRejectMsg* amf_msg) {
  int size = AUTHENTICATION_REJECT_MINIMUM_LENGTH;
  amf_msg->extended_protocol_discriminator.extended_proto_discriminator =
      M5G_MOBILITY_MANAGEMENT_MESSAGES;
  amf_msg->message_type.msg_type = AUTH_REJECT;
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, size);
}
/****************************************************************************
 **                                                                        **
 ** Name:              amf_as_security_req()                               **
 **                                                                        **
 ** Description:       Processes the AMFAS-SAP security request primitive  **
 **                                                                        **
 ** AMFAS-SAP-AMF->AS: SECURITY_REQ - Security mode control procedure      **
 **                                                                        **
 ** Inputs:  msg:      The AMFAS-SAP primitive to process                  **
 **          Others:   None                                                **
 **                                                                        **
 ** Outputs: as_msg:   The message to send to the AS                       **
 **          Return:   The identifier of the AS message                    **
 **          Others:   None                                                **
 **                                                                        **
 ***************************************************************************/
static int amf_as_security_req(
    const amf_as_security_t* msg, m5g_dl_info_transfer_req_t* as_msg) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int size = 0;
  amf_nas_message_t nas_msg;

  memset(&nas_msg, 0, sizeof(amf_nas_message_t));

  /*
   * Setup the AS message
   */
  if (msg) {
    as_msg->s_tmsi.amf_set_id = msg->guti.guamfi.amf_set_id;
    as_msg->s_tmsi.m_tmsi     = msg->guti.m_tmsi;
    as_msg->ue_id             = msg->ue_id;
  } else {
    as_msg->ue_id = msg->ue_id;
  }
  /*
   * Setup the NAS security header
   */
  AMFMsg* amf_msg = amf_as_set_header(&nas_msg, &msg->sctx);
  /*
   * Setup the NAS security message
   */
  if (amf_msg) {
    switch (msg->msg_type) {
      case AMF_AS_MSG_TYPE_IDENT:
        size = amf_identity_request(msg, &amf_msg->msg.identityrequestmsg);
        nas_msg.header.security_header_type =
            SECURITY_HEADER_TYPE_NOT_PROTECTED;
        nas_msg.header.extended_protocol_discriminator           = 0x7E;
        nas_msg.plain.amf.header.message_type                    = 0x5B;
        nas_msg.plain.amf.header.extended_protocol_discriminator = 0x7E;
        nas_msg.plain.amf.msg.identityrequestmsg.extended_protocol_discriminator
            .extended_proto_discriminator                               = 0x7e;
        nas_msg.plain.amf.msg.identityrequestmsg.message_type.msg_type  = 0x5b;
        nas_msg.plain.amf.msg.identityrequestmsg.m5gs_identity_type.toi = 1;

        break;
      case AMF_AS_MSG_TYPE_AUTH: {
        ue_m5gmm_context_s* ue_context =
            amf_ue_context_exists_amf_ue_ngap_id(as_msg->ue_id);
        amf_context_t* amf_ctx = NULL;

        amf_ctx = &ue_context->amf_context;

        size                                                     = 50;
        nas_msg.header.extended_protocol_discriminator           = 0x7E;
        nas_msg.header.security_header_type                      = 0x0;
        nas_msg.plain.amf.header.extended_protocol_discriminator = 0x7e;
        nas_msg.plain.amf.header.message_type                    = 0x56;
        nas_msg.plain.amf.msg.authenticationrequestmsg
            .extended_protocol_discriminator.extended_proto_discriminator =
            0x7e;
        nas_msg.plain.amf.msg.authenticationrequestmsg.message_type.msg_type =
            0x56;
        nas_msg.plain.amf.msg.authenticationrequestmsg.nas_key_set_identifier
            .tsc = 0;
        nas_msg.plain.amf.msg.authenticationrequestmsg.nas_key_set_identifier
            .nas_key_set_identifier = msg->ksi;
        memcpy(
            nas_msg.plain.amf.msg.authenticationrequestmsg.auth_rand.rand_val,
            amf_ctx->_vector[amf_ctx->_security.eksi % MAX_EPS_AUTH_VECTORS]
                .rand,
            RAND_LENGTH_OCTETS);

        memcpy(
            nas_msg.plain.amf.msg.authenticationrequestmsg.auth_autn.AUTN,
            amf_ctx->_vector[amf_ctx->_security.eksi % MAX_EPS_AUTH_VECTORS]
                .autn,
            AUTN_LENGTH_OCTETS);

        uint8_t abba_buff[] = {0x00, 0x00};
        memcpy(
            &(nas_msg.plain.amf.msg.authenticationrequestmsg.abba.contents),
            (const char*) abba_buff, 2);
        nas_msg.plain.amf.msg.authenticationrequestmsg.auth_rand.iei = 0x21;
        nas_msg.plain.amf.msg.authenticationrequestmsg.auth_autn.iei = 0x20;
      } break;
      case AMF_AS_MSG_TYPE_SMC: {
        size = 8;
        nas_msg.security_protected.plain.amf.header
            .extended_protocol_discriminator                     = 0x7e;
        nas_msg.security_protected.plain.amf.header.message_type = 0x5d;
        nas_msg.security_protected.plain.amf.msg.securitymodecommandmsg
            .extended_protocol_discriminator.extended_proto_discriminator =
            0x7e;
        nas_msg.security_protected.plain.amf.msg.securitymodecommandmsg
            .sec_header_type.sec_hdr = 0;
        nas_msg.security_protected.plain.amf.msg.securitymodecommandmsg
            .spare_half_octet.spare = 0;
        nas_msg.security_protected.plain.amf.msg.securitymodecommandmsg
            .message_type.msg_type = 0x5D;
        ue_m5gmm_context_s* ue_context =
            amf_ue_context_exists_amf_ue_ngap_id(as_msg->ue_id);
        if (ue_context) {
          amf_security_context_t* amf_security_context =
              &ue_context->amf_context._security;
          nas_msg.security_protected.plain.amf.msg.securitymodecommandmsg
              .nas_sec_algorithms.tca =
              amf_security_context->selected_algorithms.encryption;
          nas_msg.security_protected.plain.amf.msg.securitymodecommandmsg
              .nas_sec_algorithms.tia =
              amf_security_context->selected_algorithms.integrity;
          // relay UE security capabilities saved to amf_context back to UE
          memcpy(
              &(nas_msg.security_protected.plain.amf.msg.securitymodecommandmsg
                    .ue_sec_capability),
              &(ue_context->amf_context.ue_sec_capability),
              sizeof(UESecurityCapabilityMsg));
        } else {
          OAILOG_ERROR(
              LOG_AMF_APP, "UE not found : " AMF_UE_NGAP_ID_FMT "\n",
              as_msg->ue_id);
          return -2;
        }
        nas_msg.security_protected.plain.amf.msg.securitymodecommandmsg
            .nas_key_set_identifier.tsc = 0;
        nas_msg.security_protected.plain.amf.msg.securitymodecommandmsg
            .nas_key_set_identifier.nas_key_set_identifier =
            ue_context->amf_context._security.eksi;
        nas_msg.security_protected.plain.amf.msg.securitymodecommandmsg
            .spare_half_octet.spare = 0;
        nas_msg.security_protected.plain.amf.msg.securitymodecommandmsg
            .imeisv_request.imeisv_request = 1;
      } break;
      default:
        OAILOG_WARNING(
            LOG_NAS_AMF,
            "AMFAS-SAP - Type of NAS security "
            "message 0x%.2x is not valid\n",
            msg->msg_type);
    }
  }

  if (size > 0) {
    amf_context_t* amf_ctx                       = NULL;
    amf_security_context_t* amf_security_context = NULL;
    ue_m5gmm_context_s* ue_mm_context =
        amf_ue_context_exists_amf_ue_ngap_id(msg->ue_id);

    if (ue_mm_context) {
      amf_ctx = &ue_mm_context->amf_context;

      if (amf_ctx) {
        if (IS_AMF_CTXT_PRESENT_SECURITY(amf_ctx)) {
          amf_security_context           = &amf_ctx->_security;
          nas_msg.header.sequence_number = amf_ctx->_security.dl_count.seq_num;
        }
      }
    } else {
      OAILOG_ERROR(
          LOG_AMF_APP,
          "ue context not found for the ue_id=" AMF_UE_NGAP_ID_FMT "\n",
          msg->ue_id);
      OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
    }

    /*
     * Encode the NAS security message
     */
    int bytes =
        amf_as_encode(&as_msg->nas_msg, &nas_msg, size, amf_security_context);

    if (bytes > 0) {
      as_msg->err_code = M5G_AS_SUCCESS;
      OAILOG_FUNC_RETURN(LOG_NAS_AMF, AS_DL_INFO_TRANSFER_REQ_);
    } else {
      OAILOG_ERROR(LOG_AMF_APP, "NAS encoding failed");
      OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, 0);
}

/****************************************************************************
 **                                                                        **
 ** Name:              amf_as_security_rej()                               **
 **                                                                        **
 ** Description:       Processes the AMFAS-SAP security request primitive  **
 **                                                                        **
 ** AMFAS-SAP-AMF->AS: SECURITY_REJ - Security mode control procedure      **
 **                                                                        **
 ** Inputs:  msg:      The AMFAS-SAP primitive to process                  **
 **          Others:   None                                                **
 **                                                                        **
 ** Outputs: as_msg:   The message to send to the AS                       **
 **          Return:   The identifier of the AS message                    **
 **          Others:   None                                                **
 **                                                                        **
 ***************************************************************************/
static int amf_as_security_rej(
    const amf_as_security_t* msg, m5g_dl_info_transfer_req_t* as_msg) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int size                  = 0;
  amf_nas_message_t nas_msg = {0};

  /*
   * Setup the AS message
   */
  if (msg) {
    as_msg->s_tmsi.amf_set_id = msg->guti.guamfi.amf_set_id;
    as_msg->s_tmsi.m_tmsi     = msg->guti.m_tmsi;
    as_msg->ue_id             = msg->ue_id;
  } else {
    as_msg->ue_id = msg->ue_id;
  }
  /*
   * Setup the NAS security header
   */
  AMFMsg* amf_msg = amf_as_set_header(&nas_msg, &msg->sctx);
  /*
   * Setup the NAS security message
   */
  if (amf_msg) {
    switch (msg->msg_type) {
      case AMF_AS_MSG_TYPE_AUTH: {
        size = amf_auth_reject(msg, &amf_msg->msg.authenticationrejectmsg);
        nas_msg.header.security_header_type =
            SECURITY_HEADER_TYPE_NOT_PROTECTED;
        nas_msg.header.extended_protocol_discriminator           = 0x7E;
        nas_msg.plain.amf.header.message_type                    = AUTH_REJECT;
        nas_msg.plain.amf.header.extended_protocol_discriminator = 0x7E;
        nas_msg.plain.amf.msg.authenticationrejectmsg
            .extended_protocol_discriminator.extended_proto_discriminator =
            0x7e;
        nas_msg.plain.amf.msg.authenticationrejectmsg.message_type.msg_type =
            AUTH_REJECT;
        break;
      }
      default: { OAILOG_DEBUG(LOG_AMF_APP, " Invalid as msg type \n"); }
    }
  }

  if (size > 0) {
    amf_context_t* amf_ctx                       = NULL;
    amf_security_context_t* amf_security_context = NULL;
    ue_m5gmm_context_s* ue_mm_context =
        amf_ue_context_exists_amf_ue_ngap_id(msg->ue_id);

    if (ue_mm_context) {
      amf_ctx = &ue_mm_context->amf_context;

      if (amf_ctx) {
        if (IS_AMF_CTXT_PRESENT_SECURITY(amf_ctx)) {
          amf_security_context           = &amf_ctx->_security;
          nas_msg.header.sequence_number = amf_ctx->_security.dl_count.seq_num;
        }
      }
    } else {
      OAILOG_ERROR(
          LOG_AMF_APP,
          "ue context not found for the ue_id=" AMF_UE_NGAP_ID_FMT "\n",
          msg->ue_id);
      OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
    }

    /*
     * Encode the NAS security message
     */
    int bytes =
        amf_as_encode(&as_msg->nas_msg, &nas_msg, size, amf_security_context);

    if (bytes > 0) {
      as_msg->err_code = M5G_AS_TERMINATED_NAS;
      OAILOG_FUNC_RETURN(LOG_NAS_AMF, AS_DL_INFO_TRANSFER_REQ_);
    } else {
      OAILOG_ERROR(LOG_AMF_APP, "NAS encoding failed");
      OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, 0);
}

int initial_context_setup_request(
    amf_ue_ngap_id_t ue_id, amf_context_t* amf_ctx, bstring nas_msg) {
  /*This message is sent by the AMF to NG-RAN node to request the setup of a UE
   * context before Registration Accept is sent to UE*/

  ue_m5gmm_context_s* ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  Ngap_initial_context_setup_request_t* req = nullptr;
  Ngap_PDUSession_Resource_Setup_Request_List_t* pdu_resource_transfer_ie =
      nullptr;
  MessageDef* message_p = nullptr;
  message_p =
      itti_alloc_new_message(TASK_AMF_APP, NGAP_INITIAL_CONTEXT_SETUP_REQ);
  if (message_p == NULL) {
    OAILOG_ERROR(
        LOG_AMF_APP,
        "Failed to allocate memory for ngap_initial_context_setup_req\n");
    return RETURNerror;
  }
  req = &message_p->ittiMsg.ngap_initial_context_setup_req;
  memset(req, 0, sizeof(Ngap_initial_context_setup_request_t));
  req->amf_ue_ngap_id = ue_id;
  gnb_ue_ngap_id_t gnb_ue_ngap_id =
      PARENT_STRUCT(amf_ctx, ue_m5gmm_context_s, amf_context)->gnb_ue_ngap_id;
  req->ran_ue_ngap_id = gnb_ue_ngap_id;
  req->ue_security_capabilities.m5g_encryption_algo |=
      (amf_ctx->ue_sec_capability.ea1 & 0001) << 15;
  req->ue_security_capabilities.m5g_encryption_algo |=
      (amf_ctx->ue_sec_capability.ea2 & 0001) << 14;
  req->ue_security_capabilities.m5g_encryption_algo |=
      (amf_ctx->ue_sec_capability.ea3 & 0001) << 13;
  req->ue_security_capabilities.m5g_encryption_algo =
      htons(req->ue_security_capabilities.m5g_encryption_algo);
  req->ue_security_capabilities.m5g_integrity_protection_algo |=
      (amf_ctx->ue_sec_capability.ia1 & 0001) << 15;
  req->ue_security_capabilities.m5g_integrity_protection_algo |=
      (amf_ctx->ue_sec_capability.ia2 & 0001) << 14;
  req->ue_security_capabilities.m5g_integrity_protection_algo |=
      (amf_ctx->ue_sec_capability.ia3 & 0001) << 13;
  req->ue_security_capabilities.m5g_integrity_protection_algo =
      htons(req->ue_security_capabilities.m5g_integrity_protection_algo);
  req->Security_Key = (unsigned char*) &amf_ctx->_security.kgnb;
  memcpy(&req->Ngap_guami, &amf_ctx->m5_guti.guamfi, sizeof(guamfi_t));

  req->ue_aggregate_max_bit_rate.dl = amf_ctx->subscribed_ue_ambr.br_dl;
  req->ue_aggregate_max_bit_rate.ul = amf_ctx->subscribed_ue_ambr.br_ul;

  for (const auto& it : ue_context->amf_context.smf_ctxt_map) {
    pdusession_setup_item_t* item = nullptr;
    pdu_resource_transfer_ie = &req->PDU_Session_Resource_Setup_Transfer_List;

    std::shared_ptr<smf_context_t> smf_context = it.second;
    if (smf_context->pdu_session_state == ACTIVE) {
      uint8_t item_num     = 0;
      uint64_t ul_pdu_ambr = 0;
      uint64_t dl_pdu_ambr = 0;
      pdu_resource_transfer_ie->no_of_items += 1;
      item_num = pdu_resource_transfer_ie->no_of_items - 1;
      item     = &pdu_resource_transfer_ie->item[item_num];
      ambr_calculation_pdu_session(smf_context, &dl_pdu_ambr, &ul_pdu_ambr);

      // pdu session id
      item->Pdu_Session_ID =
          smf_context->smf_proc_data.pdu_session_identity.pdu_session_id;

      // pdu ambr
      item->PDU_Session_Resource_Setup_Request_Transfer
          .pdu_aggregate_max_bit_rate.dl = dl_pdu_ambr;
      item->PDU_Session_Resource_Setup_Request_Transfer
          .pdu_aggregate_max_bit_rate.ul = ul_pdu_ambr;

      // pdu session type
      item->PDU_Session_Resource_Setup_Request_Transfer.pdu_ip_type.pdn_type =
          smf_context->pdu_address.pdn_type;

      // up transport info
      memcpy(
          &item->PDU_Session_Resource_Setup_Request_Transfer
               .up_transport_layer_info.gtp_tnl.gtp_tied,
          smf_context->gtp_tunnel_id.upf_gtp_teid, GNB_TEID_LEN);
      item->PDU_Session_Resource_Setup_Request_Transfer.up_transport_layer_info
          .gtp_tnl.endpoint_ip_address = blk2bstr(
          &smf_context->gtp_tunnel_id.upf_gtp_teid_ip_addr, GNB_IPV4_ADDR_LEN);

      // qos flow list
      memcpy(
          &item->PDU_Session_Resource_Setup_Request_Transfer
               .qos_flow_setup_request_list.qos_flow_req_item,
          &smf_context->pdu_resource_setup_req
               .pdu_session_resource_setup_request_transfer
               .qos_flow_setup_request_list.qos_flow_req_item,
          sizeof(qos_flow_setup_request_item));
    }
  }

  if (nas_msg) {
    req->nas_pdu = nas_msg;
  } else {
    OAILOG_DEBUG(LOG_AMF_APP, "Invalid nas_msg for registration accept");
    return RETURNerror;
  }

  amf_send_msg_to_task(&amf_app_task_zmq_ctx, TASK_NGAP, message_p);
  return RETURNok;
}

/****************************************************************************
 **                                                                        **
 ** Name:             amf_as_establish_cnf()                               **
 **                                                                        **
 ** Description:      Processes the AMFAS-SAP connection establish confirm **
 **      primitive of PDU session                                          **
 **                                                                        **
 ** AMFAS-SAP-AMF->AS:ESTABLISH_CNF - NAS signaling connection            **
 **                                                                        **
 ** Inputs:   msg:    The AMFAS-SAP primitive to process                   **
 **           Others: None                                                 **
 **                                                                        **
 ** Outputs:  as_msg: The message to send to the AS                        **
 **           Return: The identifier of the AS message                     **
 **           Others: None                                                 **
 **                                                                        **
 ***************************************************************************/
uint16_t amf_as_establish_cnf(
    const amf_as_establish_t* msg, nas5g_establish_rsp_t* as_msg) {
  int size    = 0;
  int ret_val = 0;
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  OAILOG_DEBUG(
      LOG_NAS_AMF,
      "Send AS connection establish confirmation for (ue_id "
      "= " AMF_UE_NGAP_ID_FMT ")\n",
      msg->ue_id);
  amf_nas_message_t nas_msg = {0};
  // Setting-up the AS message
  as_msg->ue_id = msg->ue_id;

  as_msg->nas_msg                              = msg->nas_msg;
  as_msg->presencemask                         = msg->presencemask;
  as_msg->m5g_service_type                     = msg->service_type;
  amf_context_t* amf_ctx                       = NULL;
  amf_security_context_t* amf_security_context = NULL;
  amf_ctx                                      = amf_context_get(msg->ue_id);
  ue_m5gmm_context_s* ue_mm_context =
      amf_ue_context_exists_amf_ue_ngap_id(msg->ue_id);
  if (amf_ctx) {
    if (IS_AMF_CTXT_PRESENT_SECURITY(amf_ctx)) {
      amf_security_context                  = &amf_ctx->_security;
      as_msg->selected_encryption_algorithm = (uint16_t) htons(
          0x10000 >> amf_security_context->selected_algorithms.encryption);
      as_msg->selected_integrity_algorithm = (uint16_t) htons(
          0x10000 >> amf_security_context->selected_algorithms.integrity);
      as_msg->nas_ul_count = 0x00000000 |
                             (amf_security_context->ul_count.overflow << 8) |
                             amf_security_context->ul_count.seq_num;
    }
  } else {
    OAILOG_WARNING(LOG_NAS_AMF, "AMFAS-SAP - AMF Context is NULL...!");
  }

  nas_amf_registration_proc_t* registration_proc =
      get_nas_specific_procedure_registration(amf_ctx);
  /*
   * Setup the NAS security header
   */
  amf_as_set_header(&nas_msg, &msg->sctx);
  switch (msg->nas_info) {
    case AMF_AS_NAS_INFO_REGISTERED:
      size = amf_reg_acceptmsg(&(msg->guti), &nas_msg);
      nas_msg.header.security_header_type =
          SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED;
      /* TODO amf_as_set_header() is incorrectly setting the security header
       * type for Registration Accept. Fix it in that function*/
      break;
    case AMF_AS_NAS_INFO_SR:
      size = amf_service_acceptmsg(msg, &nas_msg);
      nas_msg.header.security_header_type =
          SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED;
      break;
    case AMF_AS_NAS_INFO_TAU:
    case AMF_AS_NAS_INFO_NONE:  // Response to SR
      as_msg->err_code = M5G_AS_SUCCESS;
      ret_val          = AS_NAS_ESTABLISH_CNF_;
      OAILOG_FUNC_RETURN(LOG_NAS_AMF, ret_val);
    default:
      OAILOG_WARNING(
          LOG_NAS_AMF,
          "AMFAS-SAP - Type of initial NAS "
          "message 0x%.2x is not valid\n",
          msg->nas_info);
      break;
  }

  if ((size > 0) && amf_security_context) {
    nas_msg.header.sequence_number = amf_security_context->dl_count.seq_num;
  } else {
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, ret_val);
  }
  /*
   * Encode the initial NAS information message
   */
  int bytes =
      amf_as_encode(&as_msg->nas_msg, &nas_msg, size, amf_security_context);

  if (bytes > 0) {
    m5gmm_state_t state =
        PARENT_STRUCT(amf_ctx, ue_m5gmm_context_s, amf_context)->mm_state;

    if (registration_proc && !(registration_proc->registration_accept_sent) &&
        (state != REGISTERED_CONNECTED)) {
      /*GNB key, generated in AMF from KAMF and shared with gNB as part of
       * InitialContextSetupRequest*/
      derive_5gkey_gnb(
          amf_security_context->kamf, as_msg->nas_ul_count,
          amf_security_context->kgnb);
      registration_proc->registration_accept_sent++;
    }

    if ((ue_mm_context->ue_context_request) ||
        ((msg->nas_info == AMF_AS_NAS_INFO_SR) &&
         (msg->pdu_session_status_ie &
          AMF_AS_PDU_SESSION_REACTIVATION_STATUS))) {
      initial_context_setup_request(as_msg->ue_id, amf_ctx, as_msg->nas_msg);
    } else {
      amf_app_handle_nas_dl_req(as_msg->ue_id, as_msg->nas_msg, M5G_AS_SUCCESS);
      ue_mm_context->mm_state = REGISTERED_CONNECTED;
    }

    as_msg->err_code = M5G_AS_SUCCESS;
    ret_val          = AS_NAS_ESTABLISH_CNF_;
  } else {
    OAILOG_ERROR(LOG_AMF_APP, "NAS encoding failed");
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, ret_val);
}

/***************************************************************************
**                                                                        **
** Name:    amf_send_registration_reject()                                **
**                                                                        **
** Description: To fill Registration reject message                       **
**                                                                        **
** Inputs:  amf_as_establish_t : msg                                      **
**          RegistrationRejectMsg : amf_msg                               **
**                                                                        **
**      Others:    None                                                   **
**                                                                        **
** Outputs:                                                               **
**      Return:    size                                                   **
**      Others:    None                                                   **
**                                                                        **
***************************************************************************/
int amf_send_registration_reject(
    const amf_as_establish_t* msg, RegistrationRejectMsg* amf_msg) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int size = AMF_HEADER_LENGTH;

  OAILOG_INFO(
      LOG_NAS_AMF, "AMFAS-SAP - Send Registration Reject message (cause=%d)\n",
      msg->amf_cause);
  amf_msg->extended_protocol_discriminator.extended_proto_discriminator =
      M5G_MOBILITY_MANAGEMENT_MESSAGES;
  amf_msg->sec_header_type.sec_hdr = SECURITY_HEADER_TYPE_NOT_PROTECTED;
  /*
   * Mandatory - Message type
   */
  amf_msg->message_type.msg_type = REG_REJECT;
  /*
   * Mandatory - AMF cause
   */
  size += AMF_CAUSE_LENGTH;
  amf_msg->m5gmm_cause.m5gmm_cause = msg->amf_cause;

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, size);
}

/****************************************************************************
 **                                                                        **
 ** Name:             amf_as_establish_Rej()                               **
 **                                                                        **
 ** Description:      Processes the AMFAS-SAP connection establish Reject  **
 **      primitive of PDU session                                          **
 **                                                                        **
 ** AMFAS-SAP-AMF->AS:ESTABLISH_REJ - NAS signaling connection             **
 **                                                                        **
 ** Inputs:   msg:    The AMFAS-SAP primitive to process                   **
 **           Others: None                                                 **
 **                                                                        **
 ** Outputs:  as_msg: The message to send to the AS                        **
 **           Return: The identifier of the AS message                     **
 **           Others: None                                                 **
 **                                                                        **
 ***************************************************************************/
static int amf_as_establish_rej(
    const amf_as_establish_t* msg, nas5g_establish_rsp_t* as_msg) {
  int size    = 0;
  int ret_val = 0;
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  OAILOG_DEBUG(
      LOG_NAS_AMF,
      "Send AS connection establish Reject for UE ID: " AMF_UE_NGAP_ID_FMT,
      msg->ue_id);
  amf_nas_message_t nas_msg = {0};
  // Setting-up the AS message
  as_msg->ue_id = msg->ue_id;

  as_msg->nas_msg                              = msg->nas_msg;
  as_msg->presencemask                         = msg->presencemask;
  as_msg->m5g_service_type                     = msg->service_type;
  amf_context_t* amf_ctx                       = NULL;
  amf_security_context_t* amf_security_context = NULL;
  amf_ctx                                      = amf_context_get(msg->ue_id);
  if (amf_ctx) {
    if (IS_AMF_CTXT_PRESENT_SECURITY(amf_ctx)) {
      amf_security_context                  = &amf_ctx->_security;
      as_msg->selected_encryption_algorithm = (uint16_t) htons(
          0x10000 >> amf_security_context->selected_algorithms.encryption);
      as_msg->selected_integrity_algorithm = (uint16_t) htons(
          0x10000 >> amf_security_context->selected_algorithms.integrity);
      as_msg->nas_ul_count = 0x00000000 |
                             (amf_security_context->ul_count.overflow << 8) |
                             amf_security_context->ul_count.seq_num;
    }
  } else {
    OAILOG_WARNING(LOG_NAS_AMF, "AMFAS-SAP - AMF Context is NULL...!");
  }

  AMFMsg* amf_msg = amf_as_set_header(&nas_msg, &msg->sctx);

  switch (msg->nas_info) {
    case AMF_AS_NAS_INFO_REGISTERED:
      size = amf_send_registration_reject(
          msg, &amf_msg->msg.registrationrejectmsg);
      nas_msg.plain.amf.header.message_type = REG_REJECT;
      nas_msg.plain.amf.header.extended_protocol_discriminator =
          M5G_MOBILITY_MANAGEMENT_MESSAGES;
      break;
    case AMF_AS_NAS_INFO_SR:
      size = amf_service_rejectmsg(msg, &amf_msg->msg.service_reject);
      nas_msg.plain.amf.header.message_type = M5G_SERVICE_REJECT;
      nas_msg.plain.amf.header.extended_protocol_discriminator =
          M5G_MOBILITY_MANAGEMENT_MESSAGES;
      break;
    case AMF_AS_NAS_INFO_TAU:
    case AMF_AS_NAS_INFO_NONE:
      as_msg->err_code = M5G_AS_SUCCESS;
      ret_val          = AS_NAS_ESTABLISH_CNF_;
      OAILOG_FUNC_RETURN(LOG_NAS_AMF, ret_val);
    default:
      OAILOG_WARNING(
          LOG_NAS_AMF,
          "AMFAS-SAP - Type of initial NAS "
          "message 0x%.2x is not valid\n",
          msg->nas_info);
      break;
  }

  /*
   * Encode the NAS security message
   */
  int bytes =
      amf_as_encode(&as_msg->nas_msg, &nas_msg, size, amf_security_context);

  if (bytes > 0) {
    as_msg->err_code = M5G_AS_TERMINATED_NAS;
    ret_val          = AS_NAS_ESTABLISH_RSP_;
  } else {
    OAILOG_ERROR(LOG_AMF_APP, "NAS encoding failed");
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, ret_val);
}

/****************************************************************************
 **                                                                        **
 ** Name:        amf_service_rejectmsg()                                   **
 **                                                                        **
 ** Description: Builds Service reject message                             **
 **                                                                        **
 **              The Service reject message is sent by the                 **
 **              network to the UEi.                                       **
 **                                                                        **
 ** Inputs:      msg:           The AMFMsg    primitive to process         **
 **              Others:        None                                       **
 **                                                                        **
 ** Outputs:     amf_msg:       The AMF message to be sent                 **
 **              Return:        The size of the AMF message                **
 **              Others:        None                                       **
 **                                                                        **
 ***************************************************************************/
static int amf_service_rejectmsg(
    const amf_as_establish_t* msg, ServiceRejectMsg* service_reject) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int size = M5G_SERVICE_REJECT_MINIMUM_LENGTH;

  if (msg->pdu_session_status_ie & AMF_AS_PDU_SESSION_STATUS) {
    service_reject->pdu_session_status.iei = PDU_SESSION_STATUS;
    service_reject->pdu_session_status.len = 0x02;
    service_reject->pdu_session_status.pduSessionStatus =
        msg->pdu_session_status;
  }

  service_reject->extended_protocol_discriminator.extended_proto_discriminator =
      M5G_MOBILITY_MANAGEMENT_MESSAGES;
  service_reject->message_type.msg_type   = M5G_SERVICE_REJECT;
  service_reject->sec_header_type.sec_hdr = SECURITY_HEADER_TYPE_NOT_PROTECTED;

  service_reject->cause.iei         = static_cast<uint8_t>(M5GIei::M5GMM_CAUSE);
  service_reject->cause.m5gmm_cause = msg->amf_cause;

  if (msg->amf_cause == AMF_CAUSE_CONGESTION) {
    service_reject->t3346Value.iei        = GPRS_TIMER2;
    service_reject->t3346Value.len        = 1;
    service_reject->t3346Value.timervalue = 60;
  }

  size += NAS5G_MESSAGE_CONTAINER_MAXIMUM_LENGTH;
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, size);
}
}  // namespace magma5g
