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

  Source      amf_as.cpp

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#include <thread>
#ifdef __cplusplus
extern "C" {
#endif
#include "log.h"
#include "3gpp_24.501.h"
#include "3gpp_24.007.h"
#include "3gpp_24.301.h"
#include "amf_as_message.h"
#include "conversions.h"
#ifdef __cplusplus
}
#endif
#include "M5gNasMessage.h"
#include "amf_app_defs.h"
#include "amf_app_ue_context_and_proc.h"
#include "amf_as.h"
#include "amf_fsm.h"
#include "amf_recv.h"
#include "nas5g_network.h"
#include "M5GDLNASTransport.h"
#include <grpcpp/impl/codegen/status.h>
#include "S6aClient.h"
#include "feg/protos/s6a_proxy.pb.h"
#include "intertask_interface_types.h"
#include "proto_msg_to_itti_msg.h"

using namespace magma;
using namespace magma::feg;
using namespace std;
typedef uint32_t amf_ue_ngap_id_t;
#define QUADLET 4
#define AMF_GET_BYTE_ALIGNED_LENGTH(LENGTH)                                    \
  LENGTH += QUADLET - (LENGTH % QUADLET)
#define AMF_CAUSE_SUCCESS (1)

namespace magma5g {
/*forward declaration*/
extern ue_m5gmm_context_s
    ue_m5gmm_global_context;  // TODO AMF_TEST global var to temporarily
                              // store context inserted to ht
static int amf_as_establish_req(amf_as_establish_t* msg, int* amf_cause);
static int amf_as_security_req(
    const amf_as_security_t* msg, m5g_dl_info_transfer_req_t* as_msg);
int amf_send_security_mode_command(
    const amf_as_security_t* msg, SecurityModeCommandMsg* amf_msg);
nas_network nas_networks;
amf_procedure_handler procedure_handler;
nas_proc nas_procedure_as;
amf_app_defs amf_app_def_as;
amf_as_dl_message as_dl_message;
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
int amf_as::amf_as_send(amf_as_t* msg) {
  int rc                       = RETURNok;
  int amf_cause                = AMF_CAUSE_SUCCESS;
  amf_as_primitive_t primitive = msg->primitive;
  amf_ue_ngap_id_t ue_id       = 0;

  switch (primitive) {
    case _AMFAS_DATA_IND:
      // TODO
      break;
    case _AMFAS_ESTABLISH_REQ:
      rc = amf_as_establish_req(
          &msg->u.establish,
          &amf_cause);  // registration request
      ue_id = msg->u.establish.ue_id;
      break;
    case _AMFAS_RELEASE_IND:
      // TODO
      break;
    default:
      /*
       * Other primitives are forwarded to NGAP
       */
      rc = amf_as::amf_as_send_ng(msg);  // TODO -  NEED-RECHECK
      break;
  }
}

/***************************************************************************
**                                                                        **
** Name:    amf_as_establish_req()                                        **
**                                                                        **
** Description: Processes the AMFAS-SAP connection establish request      **
**      primitive                                                         **
**                                                                        **
** AMFAS-SAP - AS->AMF : ESTABLISH_REQ - NAS signalling connection        **
**     The AS notifies the NAS that establishment of the signal-          **
**     ling connection has been requested to tranfer initial NAS          **
**     message from the UE.                                               **
**                                                                        **
**      Inputs:  msg:       The AMFAS-SAP primitive to process            **
**      Others:    None                                                   **
**                                                                        **
**      Outputs:   amf_cause: AMF cause code                              **
**      Return:    RETURNok, RETURNerror                                  **
**      Others:    None                                                   **
**                                                                        **
***************************************************************************/
static int amf_as_establish_req(amf_as_establish_t* msg, int* amf_cause) {
  amf_context_t* amf_ctx                       = NULL;
  amf_security_context_t* amf_security_context = NULL;
  amf_nas_message_decode_status_t decode_status;  //    = {0};
  int decoder_rc        = 1;                      // TODO enable
  int rc                = RETURNerror;
  tai_t originating_tai = {0};
  amf_nas_message_t nas_msg;  // TODO AMF_TEST // Union of nas messages
  // AMFMsg nas_msg;  // TODO AMF_TEST verify with Sanjay
  ue_m5gmm_context_s ue_m5gmm_context;
  ue_m5gmm_context.mm_state = UNREGISTERED;
  amf_ctx = &ue_m5gmm_global_context.amf_context;  // TODO AMF_TEST global var
                                                   // to temporarily store
                                                   // context inserted to ht

  if (amf_ctx) {
    if (IS_AMF_CTXT_PRESENT_SECURITY(amf_ctx)) {
      OAILOG_INFO(LOG_AMF_APP, "AMF_TEST:AMF_CTXT_PRESENT_SECURITY ");
      amf_security_context = &amf_ctx->_security;  // AMF_TEST
    }
  }

  /*
   * Decode initial NAS message
   */
  OAILOG_INFO(LOG_AMF_APP, "AMF_TEST: Decoding NAS Message");
  decoder_rc = nas5g_message_decode(
      msg->nas_msg->data, &nas_msg, blength(msg->nas_msg), amf_security_context,
      &decode_status);
  OAILOG_INFO(LOG_AMF_APP, "AMF_TEST: rc = %d", decoder_rc);
  nas_networks.bdestroy_wrapper(&msg->nas_msg);

  // TODO conditional IE error
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
  } else {
    OAILOG_INFO(LOG_AMF_APP, "AMF_TEST: NAS Decode Success");
  }

  /*
   * Process initial NAS message
   */
  AMFMsg* amf_msg = &nas_msg.plain.amf;
  switch (amf_msg->header.message_type) {
    case REG_REQUEST:  // REGISTRATION_REQUEST:
      memcpy(&originating_tai, &msg->tai, sizeof(originating_tai));
      rc = procedure_handler.amf_handle_registration_request(
          msg->ue_id, &originating_tai, &msg->ecgi,
          &amf_msg->registrationrequestmsg, msg->is_initial,
          msg->is_amf_ctx_new, *amf_cause, decode_status);
      break;
    case M5G_IDENTITY_RESPONSE:
      rc = procedure_handler.amf_handle_identity_response(
          msg->ue_id, &amf_msg->identityresponsemsg.m5gs_mobile_identity,
          *amf_cause, decode_status);
      // msg->ue_id, &amf_msg->identityrequestmsg, *amf_cause, decode_status);
      break;
    case AUTH_RESPONSE:  // M5G_AUTHENTICATION_RESPONSE:
      rc = procedure_handler.amf_handle_authentication_response(
          msg->ue_id, &amf_msg->authenticationresponsemsg, *amf_cause,
          decode_status);
      break;
    case SEC_MODE_COMPLETE:  // M5G_SECURITY_MODE_COMPLETE:
      rc = procedure_handler.amf_handle_securitycomplete_response(
          msg->ue_id, decode_status);
      break;

    case REG_COMPLETE:  // REGISTRATION_COMPLETE:
      rc = procedure_handler.amf_handle_registrationcomplete_response(
          msg->ue_id, &amf_msg->registrationcompletemsg, *amf_cause,
          decode_status);
      break;
    case DE_REG_REQUEST_UE_ORIGIN:  // DEREGISTRATION Request from UE
      OAILOG_INFO(
          LOG_NAS_AMF,
          "AMF_TEST: Processing UE originated Deregistration procedure"
          " with NGAP UE ID %d \n",
          (uint32_t) msg->ue_id);
      rc = procedure_handler.amf_handle_deregistration_ue_origin_req(
          msg->ue_id, &amf_msg->deregistrationequesmsg, *amf_cause,
          decode_status);
      break;
    case ULNASTRANSPORT:
      OAILOG_INFO(LOG_NAS_AMF, "AMF_TEST: Processing UL NAS Transport Message");
      rc = procedure_handler.amf_smf_send(
          msg->ue_id, &amf_msg->uplinknas5gtransport, *amf_cause);
      break;
      // more case to come......
    default:
      OAILOG_INFO(
          LOG_NAS_AMF, "AMF_TEST: unknown header.message_type: %d, from %s\n",
          amf_msg->header.message_type, __FUNCTION__);
  }
}

/**************************************************************************
 **                                                                      **
 ** Name       : amf_as_send_ng()                                        **
 **                                                                      **
 ** Description: Builds NAS message according to the given AMFAS Service **
 **      Access Point primitive and sends it to the Access Stratum       **
 **      sublayer                                                        **
 **                                                                      **
 ** Inputs     : msg: The AMFAS-SAP primitive to be sent                 **
 **      Others: None                                                    **
 **                                                                      **
 ** Outputs:     None                                                    **
 **      Return: RETURNok, RETURNerror                                   **
 **      Others: None                                                    **
 **                                                                      **
 *************************************************************************/
int amf_as::amf_as_send_ng(const amf_as_t* msg) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  amf_as_message_t as_msg = {0};

  switch (msg->primitive) {
    case _AMFAS_DATA_REQ:
      as_msg.msg_id = as_dl_message.amf_as_data_req(
          &msg->u.data, &as_msg.msg.dl_info_transfer_req);
      break;
    case _AMFAS_ESTABLISH_CNF:
      as_msg.msg_id = as_dl_message.amf_as_establish_cnf(
          &msg->u.establish, &as_msg.msg.nas_establish_rsp);
      break;
    case _AMFAS_ESTABLISH_REJ:
      as_msg.msg_id = as_dl_message.amf_as_establish_rej(
          &msg->u.establish, &as_msg.msg.nas_establish_rsp);
      break;
    case _AMFAS_SECURITY_REQ:
      as_msg.msg_id = amf_as_security_req(
          &msg->u.security, &as_msg.msg.dl_info_transfer_req);
      break;
      // more case to wright......
    default:
      as_msg.msg_id = 0;
      break;
  }

  /*
   * Send the message to the Access Stratum or NGAP in case of AMF
   */
  if (as_msg.msg_id > 0) {
    OAILOG_DEBUG(
        LOG_NAS_AMF,
        "AMFAS-SAP - "
        "Sending msg with id 0x%x, primitive (%d) to NGAP layer for "
        "transmission\n",
        as_msg.msg_id, msg->primitive);

    switch (as_msg.msg_id) {
      case AS_DL_INFO_TRANSFER_REQ_: {
        amf_app_def_as.amf_app_handle_nas_dl_req(
            as_msg.msg.dl_info_transfer_req.ue_id,
            as_msg.msg.dl_info_transfer_req.nas_msg,
            as_msg.msg.dl_info_transfer_req.err_code);
        OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
      } break;
      case AS_NAS_ESTABLISH_RSP_:
      case AS_NAS_ESTABLISH_CNF_: {
        if (as_msg.msg.nas_establish_rsp.err_code == M5G_AS_SUCCESS) {
          // This flow is to release the UE context after sending the NAS
          // message.
          amf_app_def_as.amf_app_handle_nas_dl_req(
              as_msg.msg.nas_establish_rsp.ue_id,
              as_msg.msg.nas_establish_rsp.nas_msg,
              as_msg.msg.nas_establish_rsp.err_code);
          as_msg.msg.nas_establish_rsp.nas_msg = NULL;
          OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
        } else {
          OAILOG_DEBUG(
              LOG_NAS_AMF,
              "AMFAS-SAP - Sending establish_cnf to AMF-APP module for UE "
              "ID: " AMF_UE_NGAP_ID_FMT
              " selected eea "
              "0x%04X selected eia 0x%04X\n",
              as_msg.msg.nas_establish_rsp.ue_id,
              as_msg.msg.nas_establish_rsp.selected_encryption_algorithm,
              as_msg.msg.nas_establish_rsp.selected_integrity_algorithm);
          /*
           * Handle success case
           */
          OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
        }
      } break;
      case AS_NAS_RELEASE_REQ_:
        OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
        break;
      default:
        break;
    }
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
  int bytes = 1;  // todo enable

  /* Ciphering algorithms, EA1 and EA2 expects length to be mode of 4,
   * so length is modified such that it will be mode of 4
   */
  AMF_GET_BYTE_ALIGNED_LENGTH(length);
  if (msg->header.security_header_type != SECURITY_HEADER_TYPE_NOT_PROTECTED) {
    amf_msg_header* header = &msg->security_protected.plain.amf.header;
    /*
     * Expand size of protected NAS message
     */
    length += NAS_MESSAGE_SECURITY_HEADER_SIZE;
    /*
     * Set header of plain NAS message
     */
    header->extended_protocol_discriminator = M5GS_MOBILITY_MANAGEMENT_MESSAGE;
    header->security_header_type =
        SECURITY_HEADER_TYPE_NOT_PROTECTED;  // TODO, needs revisit, logic seems
                                             // off
  }

  /*
   * Allocate memory to the NAS information container
   */
  *info = bfromcstralloc(length, "\0");

  if (*info) {
    /*
     * Encode the NAS message
     */
    AmfMsg amf_msg_test;
    bytes =
        nas5g_message_encode((*info)->data, msg, length, amf_security_context);

    if (bytes > 0) {
      (*info)->slen = bytes;
    } else {
      nas_networks.bdestroy_wrapper(info);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, bytes);
}
/****************************************************************************
 **                                                                        **
 ** Name:        amf_send_dl_nas_transportmsg()                            **
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
static int amf_send_dl_nas_transportmsg(
    const amf_as_data_t* msg, DLNASTransportMsg* amf_msg) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int size = AMF_HEADER_MAXIMUM_LENGTH;
  /*
   * Mandatory - Message type
   */
  amf_msg->message_type.msg_type = DOWNLINK_NAS_TRANSPORT;
  /*
   * Mandatory - Nas message container
   */
  size += NAS5G_MESSAGE_CONTAINER_MAXIMUM_LENGTH;
  memcpy(
      amf_msg->payload_container.contents, &(msg->nas_msg),
      sizeof(msg->nas_msg));
  OAILOG_INFO(LOG_NAS_AMF, "AMFAS-SAP - Sending DL NAS - DL NAS5G Transport\n");
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, size);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_as_data_req()                                        **
 **                                                                        **
 ** Description: Processes the AMFAS-SAP data transfer request             **
 **      primitive                                                 **
 **                                                                        **
 ** AMFAS-SAP - AMF->AS : DATA_REQ - Data transfer procedure                **
 **                                                                        **
 ** Inputs:  msg:       The AMFAS-SAP primitive to process         **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     as_msg:    The message to send to the AS              **
 **      Return:    The identifier of the AS message           **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
uint16_t amf_as_dl_message::amf_as_data_req(
    const amf_as_data_t* msg, m5g_dl_info_transfer_req_t* as_msg) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int size       = 0;
  int is_encoded = false;
  amf_nas_message_t nas_msg;
  nas_msg.security_protected.header           = {0};
  nas_msg.security_protected.plain.amf.header = {0};
  nas_msg.security_protected.plain.amf.header = {0};

  /*
   * Setup the AS message
   */
  if (msg->guti) {
    as_msg->s_tmsi.amf_code = msg->guti->guamfi.amf_code;
    as_msg->s_tmsi.m_tmsi   = msg->guti->m_tmsi;
  } else {
    as_msg->ue_id = msg->ue_id;
  }

  /*
   * Setup the NAS security header
   */
  AMFMsg* amf_msg = amf_as::amf_as_set_header(
      &nas_msg, &msg->sctx);  // all header part==> all mendatory field

  /*
   * Setup the NAS information message
   */
  if (amf_msg) {
    switch (msg->nas_info) {
      case AMF_AS_NAS_DATA_REGISTRATION_ACCEPT:
        size = amf_registration_procedure::amf_send_registration_accept_dl_nas(
            msg,
            &amf_msg->registrationacceptmsg);  // make the contents of
                                               // registration accept message
        break;
      case AMF_AS_NAS_DL_NAS_TRANSPORT:
        // DL messages to NGAP on Identity/Authentication request
        size =
            amf_send_dl_nas_transportmsg(msg, &amf_msg->downlinknas5gtransport);
        break;
      case AMF_AS_NAS_DATA_DEREGISTRATION_ACCEPT: {
        // fill DL NAS message of deregistration accept
        // 0 1 0 0 0 1 1 0 Deregistration accept (UE originating) 70
        OAILOG_INFO(
            LOG_AMF_APP, "AMF_TEST: Sending DEREGISTRATION_ACCEPT to UE\n");
        size = 5;

        amf_msg->deregistrationacceptmsg.message_type.msg_type =
            DEREGISTRATION_ACCEPT_UE_INIT;
        nas_msg.security_protected.plain.amf.header
            .extended_protocol_discriminator = 0x7e;
        nas_msg.security_protected.plain.amf.header.message_type =
            DEREGISTRATION_ACCEPT_UE_INIT;
        nas_msg.security_protected.plain.amf.deregistrationacceptmsg
            .extended_protocol_discriminator.extended_proto_discriminator =
            0x7e;
        nas_msg.security_protected.plain.amf.identityrequestmsg.message_type
            .msg_type = DEREGISTRATION_ACCEPT_UE_INIT;
      } break;
      default:
        /*
         * Send other NAS messages as already encoded SMF messages
         */
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
          amf_security_context = &amf_ctx->_security;
        }
      }
    } else {
      OAILOG_INFO(LOG_AMF_APP, "AMF_TEST: get the security context\n");
      amf_ctx =
          &ue_m5gmm_global_context.amf_context;  // TODO AMF_TEST global var to
                                                 // temporarily store context
                                                 // inserted to ht
      if (IS_AMF_CTXT_PRESENT_SECURITY(amf_ctx)) {
        OAILOG_INFO(LOG_AMF_APP, "AMF_TEST: get the security context\n");
        amf_security_context = &amf_ctx->_security;
      }
    }
    if (amf_security_context) {
      nas_msg.header.sequence_number = amf_security_context->dl_count.seq_num;
      OAILOG_DEBUG(
          LOG_NAS_AMF, "Set nas_msg.header.sequence_number -> %u\n",
          nas_msg.header.sequence_number);
    } else {
      OAILOG_ERROR(
          LOG_NAS_AMF, "Security context is NULL for UE -> %d\n", msg->ue_id);
      OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
    }
    OAILOG_INFO(LOG_AMF_APP, "AMF_TEST: start NAS encoding\n");
    if (!is_encoded) {
      /*
       * Encode the NAS information message
       */
      // TODO re-check with team on function name
      bytes =
          amf_as_encode(&as_msg->nas_msg, &nas_msg, size, amf_security_context);
    } else {
      // TODO
    }

    // Free any allocated data
    switch (msg->nas_info) {
      // amf_information message and downlink_nas_transtport is the only message
      // that has allocated data
      case AMF_AS_NAS_DATA_REGISTRATION_ACCEPT:
        // TODO
        break;
        // many more remain....
    }

    if (bytes > 0) {
      OAILOG_INFO(LOG_AMF_APP, "AMF_TEST: NAS encoding success\n");
      as_msg->err_code = M5G_AS_SUCCESS;
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
AMFMsg* amf_as::amf_as_set_header(
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

  /*
   * A valid 5G CN security context exists but NAS integrity key
   * * * * is not available
   */
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, NULL);
}

static int amf_send_identity_request(
    const amf_as_security_t* msg, IdentityRequestMsg* amf_msg) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int size = AMF_HEADER_MAXIMUM_LENGTH;
  /*
   * Mandatory - Message type
   */
  amf_msg->message_type.msg_type = IDENTITY_REQUEST;
  /*
   *Mandatory - Nas message container
   */
// TODO keep the header where it belong
#include "3gpp_24.008.h"
  size += IDENTITY_TYPE_2_IE_MAX_LENGTH;
  if (msg->ident_type == IDENTITY_TYPE_2_IMSI) {
    amf_msg->m5gs_identity_type.toi = IDENTITY_TYPE_2_IMSI;
  } else {
    // TODO, handle else for timsi, imei;
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, size);
}

static void _s6a_handle_authentication_info_ans(
    const std::string& imsi, uint8_t imsi_length, const grpc::Status& status,
    feg::AuthenticationInformationAnswer response, s6a_auth_info_ans_t* aia_p) {
  OAILOG_INFO(LOG_AMF_APP, "AMF_TEST: callback for s6a_air invoked");
  strncpy(aia_p->imsi, imsi.c_str(), imsi_length);
  aia_p->imsi_length = imsi_length;

  if (status.ok()) {
    if (response.error_code() < feg::ErrorCode::COMMAND_UNSUPORTED) {
      magma::convert_proto_msg_to_itti_s6a_auth_info_ans(response, aia_p);
      std::cout << "[INFO] "
                << "Received S6A-AUTHENTICATION_INFORMATION_ANSWER for IMSI: "
                << imsi << "; Status: " << status.error_message()
                << "; StatusCode: " << response.error_code() << std::endl;
    }
  } else {
    OAILOG_INFO(
        LOG_AMF_APP,
        "AMF_TEST: S6A-AUTHENTICATION_INFORMATION_ANSWER failed with "
        "status:%d, StatusCode:%d\n",
        status.error_message(), response.error_code());
    std::cout << "[ERROR] " << status.error_code() << ": "
              << status.error_message() << std::endl;
    std::cout
        << "[ERROR] Received S6A-AUTHENTICATION_INFORMATION_ANSWER for IMSI: "
        << imsi << "; Status: " << status.error_message()
        << "; ErrorCode: " << response.error_code() << std::endl;
  }
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
  nas_msg.security_protected.header           = {0};
  nas_msg.security_protected.plain.amf.header = {0};
  nas_msg.security_protected.plain.amf.header = {0};

  /*
   * Setup the AS message
   */
  if (msg) {
    as_msg->s_tmsi.amf_code = msg->guti.guamfi.amf_code;
    as_msg->s_tmsi.m_tmsi   = msg->guti.m_tmsi;
    as_msg->ue_id = msg->ue_id;  // TODO AMF_TEST: Originally in "else"
  } else {
    as_msg->ue_id = msg->ue_id;
  }
  /*
   * Setup the NAS security header
   */
  AMFMsg* amf_msg = amf_as::amf_as_set_header(&nas_msg, &msg->sctx);
  /*
   * Setup the NAS security message
   */
  if (amf_msg) switch (msg->msg_type) {
      case AMF_AS_MSG_TYPE_IDENT:
        OAILOG_INFO(LOG_AMF_APP, "AMF_TEST: Sending IDENTITY_REQUEST to UE\n");
        size = amf_send_identity_request(msg, &amf_msg->identityrequestmsg);
        nas_msg.header.extended_protocol_discriminator           = 0x7E;
        nas_msg.header.security_header_type                      = 0x0;
        nas_msg.plain.amf.header.message_type                    = 0x5B;
        nas_msg.plain.amf.header.extended_protocol_discriminator = 0x7E;
        nas_msg.plain.amf.identityrequestmsg.extended_protocol_discriminator
            .extended_proto_discriminator                           = 0x7e;
        nas_msg.plain.amf.identityrequestmsg.message_type.msg_type  = 0x5b;
        nas_msg.plain.amf.identityrequestmsg.m5gs_identity_type.toi = 1;
        break;
      case AMF_AS_MSG_TYPE_AUTH: {
        s6a_auth_info_req_t air_t;
        memset(&air_t, 0, sizeof(s6a_auth_info_req_t));
        extern ue_m5gmm_context_s
            ue_m5gmm_global_context;  // TODO AMF_TEST global var to temporarily
                                      // store context inserted to ht
        ue_m5gmm_context_s* ue_context =
            amf_ue_context_exists_amf_ue_ngap_id(as_msg->ue_id);
        if (ue_context) {
          IMSI64_TO_STRING(
              ue_context->amf_context._imsi64, air_t.imsi,
              // ue_m5gmm_global_context.amf_context._imsi.length);
              15);
        } else {
          ue_context = &ue_m5gmm_global_context;  // TODO AMF_TEST global var to
                                                  // temporarily store context
                                                  // inserted to ht
          IMSI64_TO_STRING(
              ue_context->amf_context._imsi64, air_t.imsi,
              // ue_m5gmm_global_context.amf_context._imsi.length);
              15);
          OAILOG_INFO(
              LOG_AMF_APP, "AMF_TEST: from amf_context, imsi:%lu\n",
              ue_context->amf_context._imsi64);
          OAILOG_INFO(
              LOG_AMF_APP, "AMF_TEST: from amf_context, imsi:%s", air_t.imsi);
        }
        char temp_imsi[IMSI_BCD_DIGITS_MAX + 1] = "208950000000031";
        strcpy(air_t.imsi, temp_imsi);
        air_t.imsi_length             = 15;
        air_t.visited_plmn.mcc_digit1 = 0x2;
        air_t.visited_plmn.mcc_digit2 = 0x0;
        air_t.visited_plmn.mcc_digit3 = 0x8;
        air_t.visited_plmn.mnc_digit1 = 0x9;
        air_t.visited_plmn.mnc_digit2 = 0x5;
        air_t.visited_plmn.mnc_digit3 = 0x0;
        air_t.nb_of_vectors           = 1;
        air_t.re_synchronization      = 0;
        s6a_auth_info_ans_t aia_t;
        memset(&aia_t, 0, sizeof(s6a_auth_info_ans_t));
        auto imsi_len = air_t.imsi_length;
        OAILOG_INFO(
            LOG_AMF_APP,
            "AMF_TEST: Sending S6A-AUTHENTICATION_INFORMATION_REQUEST\n");
        magma::S6aClient::authentication_info_req(
            &air_t, [imsiStr = std::string(air_t.imsi), imsi_len, &aia_t](
                        grpc::Status status,
                        feg::AuthenticationInformationAnswer response) {
              _s6a_handle_authentication_info_ans(
                  imsiStr, imsi_len, status, response, &aia_t);
            });
        std::this_thread::sleep_for(
            std::chrono::milliseconds(60));  // TODO remove this blocking call
        OAILOG_INFO(
            LOG_AMF_APP,
            "AMF_TEST: after S6A-AUTHENTICATION_INFORMATION_REQUEST\n");
        OAILOG_INFO(LOG_AMF_APP, "AMF_TEST: imsi:%s\n", air_t.imsi);

        if (aia_t.auth_info.nb_of_vectors ==
            1) {  // TODO better conditional checks!!!!
          // if(aia_t.auth_info.nb_of_vectors != 1) { //bypassing s6a_AIR
          OAILOG_INFO(LOG_AMF_APP, "AMF_TEST: imsi in aia_t:%s\n", aia_t.imsi);
          nas_msg.plain.amf.authenticationrequestmsg.auth_autn.AUTN.assign(
              (const char*) aia_t.auth_info.eutran_vector[0].autn,
              AUTN_LENGTH_OCTETS);
          nas_msg.plain.amf.authenticationrequestmsg.auth_rand.rand_val.assign(
              (const char*) aia_t.auth_info.eutran_vector[0].rand,
              RAND_LENGTH_OCTETS);
          memcpy(
              ue_context->amf_context
                  ._vector
                      [ue_context->amf_context._security.eksi %
                       MAX_EPS_AUTH_VECTORS]
                  .autn,
              aia_t.auth_info.eutran_vector[0].autn, AUTN_LENGTH_OCTETS);
          memcpy(
              ue_context->amf_context
                  ._vector
                      [ue_context->amf_context._security.eksi %
                       MAX_EPS_AUTH_VECTORS]
                  .rand,
              aia_t.auth_info.eutran_vector[0].rand, RAND_LENGTH_OCTETS);
          for (int i = 0; i < AUTN_LENGTH_OCTETS; i++) {
            OAILOG_DEBUG(
                LOG_AMF_APP, "AMF_TEST: autn[%d]:%x", i,
                aia_t.auth_info.eutran_vector[0].autn[i]);
          }
          for (int i = 0; i < RAND_LENGTH_OCTETS; i++) {
            OAILOG_DEBUG(
                LOG_AMF_APP, "AMF_TEST: rand[%d]:%x", i,
                aia_t.auth_info.eutran_vector[0].rand[i]);
          }
        } else {
          // TODO register error, s6a_air failed and return
          OAILOG_INFO(LOG_AMF_APP, "s6a_air request failed\n");
          uint8_t autn_buff[] = {0x88, 0x21, 0x9a, 0x2b, 0xd5, 0x90,
                                 0x90, 0x01, 0x98, 0x1e, 0x81, 0x4f,
                                 0x29, 0x83, 0x21, 0xd2};
          nas_msg.security_protected.plain.amf.authenticationrequestmsg
              .auth_autn.AUTN.assign((const char*) autn_buff, 16);
          uint8_t rand_buff[] = {0xad, 0x7f, 0x25, 0x2e, 0x97, 0x48,
                                 0x57, 0x35, 0x70, 0xfe, 0x24, 0x5e,
                                 0x41, 0x84, 0x60, 0x40};
          nas_msg.security_protected.plain.amf.authenticationrequestmsg
              .auth_rand.rand_val.assign((const char*) rand_buff, 16);
        }
        OAILOG_INFO(
            LOG_AMF_APP, "AMF_TEST: Sending AUTHENTICATION_REQUEST to UE\n");
        size                                                     = 50;
        nas_msg.header.extended_protocol_discriminator           = 0x7E;
        nas_msg.header.security_header_type                      = 0x0;
        nas_msg.plain.amf.header.extended_protocol_discriminator = 0x7e;
        nas_msg.plain.amf.header.message_type                    = 0x56;
        nas_msg.plain.amf.authenticationrequestmsg
            .extended_protocol_discriminator.extended_proto_discriminator =
            0x7e;
        nas_msg.plain.amf.authenticationrequestmsg.message_type.msg_type = 0x56;
        nas_msg.plain.amf.authenticationrequestmsg.nas_key_set_identifier.tsc =
            0;
        nas_msg.plain.amf.authenticationrequestmsg.nas_key_set_identifier
            .nas_key_set_identifier = 0x0;
        uint8_t abba_buff[]         = {0x00, 0x00};
        nas_msg.plain.amf.authenticationrequestmsg.abba.contents.assign(
            (const char*) abba_buff, 2);
        nas_msg.plain.amf.authenticationrequestmsg.auth_rand.iei = 0x21;
        nas_msg.plain.amf.authenticationrequestmsg.auth_autn.iei = 0x20;
        nas_msg.plain.amf.authenticationrequestmsg.auth_rand
            .iei = 0x21;
        nas_msg.plain.amf.authenticationrequestmsg.auth_autn
            .iei = 0x20;
      }
      break;

      case AMF_AS_MSG_TYPE_SMC: {
        size = 8;
        OAILOG_INFO(
            LOG_AMF_APP, "AMF_TEST: Sending SECURITY_MODE_COMMAND to UE\n");
        nas_msg.security_protected.plain.amf.header
            .extended_protocol_discriminator                     = 0x7e;
        nas_msg.security_protected.plain.amf.header.message_type = 0x5d;
        nas_msg.security_protected.plain.amf.securitymodecommandmsg
            .extended_protocol_discriminator.extended_proto_discriminator =
            0x7e;
        nas_msg.security_protected.plain.amf.securitymodecommandmsg
            .sec_header_type.sec_hdr = 0;
        nas_msg.security_protected.plain.amf.securitymodecommandmsg
            .spare_half_octet.spare = 0;
        nas_msg.security_protected.plain.amf.securitymodecommandmsg.message_type
            .msg_type = 0x5D;
        ue_m5gmm_context_s* ue_context =
            amf_ue_context_exists_amf_ue_ngap_id(as_msg->ue_id);
        if (ue_context) {
          amf_security_context_t* amf_security_context =
              &ue_context->amf_context._security;
          amf_security_context->selected_algorithms.integrity =
              0;  // TODO get this computed
          amf_security_context->selected_algorithms.encryption =
              0;  // TODO get this computed
          nas_msg.security_protected.plain.amf.securitymodecommandmsg
              .nas_sec_algorithms.tca =
              amf_security_context->selected_algorithms.encryption;
          nas_msg.security_protected.plain.amf.securitymodecommandmsg
              .nas_sec_algorithms.tia =
              amf_security_context->selected_algorithms.integrity;
          // relay UE security capabilities saved to amf_context back to UE
          memcpy(
              &(nas_msg.security_protected.plain.amf.securitymodecommandmsg
                    .ue_sec_capability),
              &(ue_context->amf_context.ue_sec_capability),
              sizeof(UESecurityCapabilityMsg));

        } else {
          ue_context = &ue_m5gmm_global_context;  // TODO AMF_TEST global var to
                                                  // temporarily store context
                                                  // inserted to ht
          amf_security_context_t* amf_security_context =
              &ue_context->amf_context._security;
          amf_security_context->selected_algorithms.integrity =
              0;  // TODO get this computed
          amf_security_context->selected_algorithms.encryption =
              0;  // TODO get this computed
          nas_msg.security_protected.plain.amf.securitymodecommandmsg
              .nas_sec_algorithms.tca =
              amf_security_context->selected_algorithms.encryption;
          nas_msg.security_protected.plain.amf.securitymodecommandmsg
              .nas_sec_algorithms.tia =
              amf_security_context->selected_algorithms.integrity;
          // relay UE security capabilities saved to amf_context back to UE
          memcpy(
              &(nas_msg.security_protected.plain.amf.securitymodecommandmsg
                    .ue_sec_capability),
              &(ue_context->amf_context.ue_sec_capability),
              sizeof(UESecurityCapabilityMsg));
        }
        nas_msg.security_protected.plain.amf.securitymodecommandmsg
            .nas_key_set_identifier.tsc = 0;
        nas_msg.security_protected.plain.amf.securitymodecommandmsg
            .nas_key_set_identifier.nas_key_set_identifier = 0x0;
        nas_msg.security_protected.plain.amf.securitymodecommandmsg
            .spare_half_octet.spare = 0;
        nas_msg.security_protected.plain.amf.securitymodecommandmsg
            .imeisv_request.imeisv_request = 1;
      } break;
      default:
        OAILOG_WARNING(
            LOG_NAS_AMF,
            "AMFAS-SAP - Type of NAS security "
            "message 0x%.2x is not valid\n",
            msg->msg_type);
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
          OAILOG_DEBUG(
              LOG_NAS_AMF, "Set nas_msg.header.sequence_number -> %u\n",
              nas_msg.header.sequence_number);
        }
      }
    } else {
      // TODO
    }

    /*
     * Encode the NAS security message
     */
    OAILOG_INFO(LOG_AMF_APP, "AMF_TEST: Start NAS encoding");
    int bytes =
        amf_as_encode(&as_msg->nas_msg, &nas_msg, size, amf_security_context);
    // Free any allocated data
    switch (msg->msg_type) {
      // authentication_request is the only message with allocated amf
      case AMF_AS_MSG_TYPE_AUTH:
        amf_free_send_authentication_request(
            &amf_msg->authenticationrequestmsg);
        break;
        // Other cases to free resources of Identity and smc
    }

    if (bytes > 0) {
      OAILOG_INFO(LOG_AMF_APP, "AMF_TEST: NAS Encoding Success");
      as_msg->err_code = M5G_AS_SUCCESS;
      OAILOG_FUNC_RETURN(LOG_NAS_AMF, AS_DL_INFO_TRANSFER_REQ_);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, 0);
}

void amf_free_send_authentication_request(
    AuthenticationRequestMsg* amf_msg_req) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  OAILOG_DEBUG(
      LOG_NAS_AMF, "AMFAS-SAP - Freeing Send Authentication Request message\n");
  OAILOG_FUNC_OUT(LOG_NAS_AMF);
}

/****************************************************************************
 **                                                                        **
 ** Name:             amf_as_establish_cnf()                               **
 **                                                                        **
 ** Description:      Processes the AMFAS-SAP connection establish confirm **
 **      primitive of PDU session                                          **
 **                                                                        **
 ** AMFAS-SAP-AMF->AS:ESTABLISH_CNF - NAS signalling connection            **
 **                                                                        **
 ** Inputs:   msg:    The AMFAS-SAP primitive to process                   **
 **           Others: None                                                 **
 **                                                                        **
 ** Outputs:  as_msg: The message to send to the AS                        **
 **           Return: The identifier of the AS message                     **
 **           Others: None                                                 **
 **                                                                        **
 ***************************************************************************/
uint16_t amf_as_dl_message::amf_as_establish_cnf(
    const amf_as_establish_t* msg, nas5g_establish_rsp_t* as_msg) {
  AMFMsg* amf_msg = NULL;
  int size        = 0;
  int ret_val     = 0;
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  OAILOG_INFO(
      LOG_NAS_AMF,
      "AMF_TEST: Send AS connection establish confirmation for (ue_id = "
      "%d)\n",
      msg->ue_id);
  amf_nas_message_t nas_msg;
  // Setting-up the AS message
  as_msg->ue_id = msg->ue_id;

  if (msg->pds_id.guti == NULL) {
    OAILOG_WARNING(LOG_NAS_AMF, "AMFAS-SAP - GUTI is NULL...");
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, ret_val);
  }

  as_msg->s_tmsi.amf_code  = msg->pds_id.guti->guamfi.amf_code;
  as_msg->s_tmsi.m_tmsi    = msg->pds_id.guti->m_tmsi;
  as_msg->nas_msg          = msg->nas_msg;
  as_msg->presencemask     = msg->presencemask;
  as_msg->m5g_service_type = msg->service_type;
  amf_context_t* amf_ctx   = NULL;
  amf_security_context_t* amf_security_context = NULL;
  amf_ctx                                      = amf_context_get(msg->ue_id);
  if (amf_ctx) {
    if (IS_AMF_CTXT_PRESENT_SECURITY(amf_ctx)) {
      amf_security_context                  = &amf_ctx->_security;
      as_msg->selected_encryption_algorithm = (uint16_t) htons(
          0x10000 >> amf_security_context->selected_algorithms.encryption);
      as_msg->selected_integrity_algorithm = (uint16_t) htons(
          0x10000 >> amf_security_context->selected_algorithms.integrity);
      as_msg->nas_ul_count =
          0x00000000 | (amf_security_context->ul_count.overflow << 8) |
          amf_security_context->ul_count.seq_num;  // This is sent to calculate
                                                   // KgNB OAILOG_INFO(
    }
  } else {
    OAILOG_WARNING(LOG_NAS_AMF, "AMFAS-SAP - AMF Context is NULL...!");
  }
  switch (msg->nas_info) {
    case AMF_AS_NAS_INFO_REGISTERD:
      /*
       * Setup the NAS security header
       */
      OAILOG_INFO(LOG_AMF_APP, "AMF_TEST: Sending REGISTRATION_ACCEPT to UE\n");
      amf_msg = amf_as::amf_as_set_header(&nas_msg, &msg->sctx);
      size    = 19;
      nas_msg.security_protected.plain.amf.header.message_type = 0x42;
      nas_msg.security_protected.plain.amf.registrationacceptmsg
          .extended_protocol_discriminator.extended_proto_discriminator = 0x7e;
      nas_msg.security_protected.plain.amf.header
          .extended_protocol_discriminator = 0x7e;
      nas_msg.security_protected.plain.amf.registrationacceptmsg
          .extended_protocol_discriminator.extended_proto_discriminator = 0x7e;
      nas_msg.security_protected.plain.amf.registrationacceptmsg.sec_header_type
          .sec_hdr = 0;
      nas_msg.security_protected.plain.amf.registrationacceptmsg
          .spare_half_octet.spare = 0;
      nas_msg.security_protected.plain.amf.registrationacceptmsg.message_type
          .msg_type = 0x42;
      nas_msg.security_protected.plain.amf.registrationacceptmsg.m5gs_reg_result
          .spare = 0;
      nas_msg.security_protected.plain.amf.registrationacceptmsg.m5gs_reg_result
          .sms_allowed = 0;
      nas_msg.security_protected.plain.amf.registrationacceptmsg.m5gs_reg_result
          .reg_result_val = 1;
      nas_msg.security_protected.plain.amf.registrationacceptmsg.mobile_id
          .mobile_identity.guti.odd_even = 0;
      nas_msg.security_protected.plain.amf.registrationacceptmsg.mobile_id.iei =
          0x77;
      nas_msg.security_protected.plain.amf.registrationacceptmsg.mobile_id.len =
          11;
      nas_msg.security_protected.plain.amf.registrationacceptmsg.mobile_id
          .mobile_identity.guti.type_of_identity = 2;
      // Filling GUTI from amf_as_establish msg
      nas_msg.security_protected.plain.amf.registrationacceptmsg.mobile_id
          .mobile_identity.guti.mcc_digit1 = msg->guti.guamfi.plmn.mcc_digit1;
      nas_msg.security_protected.plain.amf.registrationacceptmsg.mobile_id
          .mobile_identity.guti.mcc_digit2 = msg->guti.guamfi.plmn.mcc_digit2;
      nas_msg.security_protected.plain.amf.registrationacceptmsg.mobile_id
          .mobile_identity.guti.mcc_digit3 = msg->guti.guamfi.plmn.mcc_digit3;
      nas_msg.security_protected.plain.amf.registrationacceptmsg.mobile_id
          .mobile_identity.guti.mnc_digit1 = msg->guti.guamfi.plmn.mnc_digit1;
      nas_msg.security_protected.plain.amf.registrationacceptmsg.mobile_id
          .mobile_identity.guti.mnc_digit2 = msg->guti.guamfi.plmn.mnc_digit2;
      nas_msg.security_protected.plain.amf.registrationacceptmsg.mobile_id
          .mobile_identity.guti.mnc_digit3 = msg->guti.guamfi.plmn.mnc_digit3;
      uint8_t* offset;
      offset = (uint8_t*) &msg->guti.m_tmsi;
      nas_msg.security_protected.plain.amf.registrationacceptmsg.mobile_id
          .mobile_identity.guti.tmsi1 = *offset;
      offset++;
      nas_msg.security_protected.plain.amf.registrationacceptmsg.mobile_id
          .mobile_identity.guti.tmsi2 = *offset;
      offset++;
      nas_msg.security_protected.plain.amf.registrationacceptmsg.mobile_id
          .mobile_identity.guti.tmsi3 = *offset;
      offset++;
      nas_msg.security_protected.plain.amf.registrationacceptmsg.mobile_id
          .mobile_identity.guti.tmsi4 = *offset;
      break;
    case AMF_AS_NAS_INFO_TAU:
      /*
       * Setup the NAS security header
       */
      amf_msg = amf_as::amf_as_set_header(&nas_msg, &msg->sctx);
      if (amf_msg) {  // TODO-RECHECK later
      }
      break;
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

  if (size > 0) {
    nas_msg.header.sequence_number = amf_security_context->dl_count.seq_num;
  } else {
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, ret_val);
  }
  /*
   * Encode the initial NAS information message
   */
  OAILOG_INFO(LOG_AMF_APP, "AMF_TEST: start NAS encoding \n");
  int bytes =
      amf_as_encode(&as_msg->nas_msg, &nas_msg, size, amf_security_context);
  // Free any allocated data
  if (msg->nas_info == AMF_AS_NAS_INFO_REGISTERD) {
    // TODO
    // bdestroy_wrapper(&(amf_msg->RegistrationAcceptMsg.esamfssagecontainer));
  }

  if (bytes > 0) {
    OAILOG_INFO(LOG_AMF_APP, "AMF_TEST: NAS encoding success\n");
    as_msg->err_code = M5G_AS_SUCCESS;
    ret_val          = AS_NAS_ESTABLISH_CNF_;
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, ret_val);
}

/************************************************************************
 ** Name:    amf_as_establish_rej()                                    **
 ** Description: Processes the AMFAS-SAP connection establish reject   **
 **      primitive w.r.t PDU session                                   **
 **                                                                    **
 ** AMFAS-SAP - AMF->AS: ESTABLISH_REJ - NAS signalling connection     **
 **                                                                    **
 ** Inputs:  msg:       The AMFAS-SAP primitive to process             **
 **          Others:    None                                           **
 **                                                                    **
 ** Outputs: as_msg:    The message to send to the AS                  **
 **          Return:    The identifier of the AS message               **
 **          Others:    None                                           **
 ***********************************************************************/
uint16_t amf_as_dl_message::amf_as_establish_rej(
    const amf_as_establish_t* msg, nas5g_establish_rsp_t* as_msg) {
  AMFMsg* amf_msg = NULL;
  int size        = 0;
  amf_nas_message_t nas_msg;
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  OAILOG_INFO(
      LOG_NAS_AMF, "AMFAS-SAP - Send AS PDU connection establish reject\n");

  /*
   * Setup the AS message
   */
  if (msg->pds_id.guti) {
    as_msg->s_tmsi.amf_code = msg->pds_id.guti->guamfi.amf_code;
    as_msg->s_tmsi.m_tmsi   = msg->pds_id.guti->m_tmsi;
  } else {
    as_msg->ue_id = msg->ue_id;
  }

  /*
   * Setup the NAS security header
   */
  amf_msg = amf_as::amf_as_set_header(&nas_msg, &msg->sctx);

  /*
   * Setup the NAS information messag
   */
  if (amf_msg) {
    switch (msg->nas_info) {
      case AMF_AS_NAS_INFO_REGISTERD:
        // TODO
        break;
      case AMF_AS_NAS_INFO_TAU:
        // TODO - TA upadate rejection will be taken care later
        // size = amf_send_tracking_area_update_reject(
        //    msg, &amf_msg->tracking_area_update_reject);
        break;
      case AMF_AS_NAS_INFO_SR:
        // TODO - Network initiated rejection will be taken care later
        // size = amf_send_service_reject(msg->amf_cause,
        // &amf_msg->service_reject);
        break;
      default:
        OAILOG_WARNING(
            LOG_NAS_AMF,
            "AMFAS-SAP - Type of initial NAS "
            "message 0x%.2x is not valid\n",
            msg->nas_info);
    }
  }

  if (size > 0) {
    amf_context_t* amf_ctx                       = NULL;
    amf_security_context_t* amf_security_context = NULL;
    ue_m5gmm_context_s* ue_m5g_context =
        amf_ue_context_exists_amf_ue_ngap_id(msg->ue_id);
    if (ue_m5g_context) {
      amf_ctx = &ue_m5g_context->amf_context;
      if (amf_ctx) {
        if (IS_AMF_CTXT_PRESENT_SECURITY(amf_ctx)) {
          amf_security_context = &amf_ctx->_security;
          nas_msg.header.sequence_number =
              amf_security_context->dl_count.seq_num;
          OAILOG_DEBUG(
              LOG_NAS_AMF, "Set nas_msg.header.sequence_number -> %u\n",
              nas_msg.header.sequence_number);
        }
      }
    }

    /*
     * Encode the initial NAS information message
     */
    int bytes =
        amf_as_encode(&as_msg->nas_msg, &nas_msg, size, amf_security_context);
    if (bytes > 0) {
      // This is to indicate AMF-APP to release the S1AP UE context after
      // sending the message.
      as_msg->err_code = M5G_AS_TERMINATED_NAS;
      OAILOG_FUNC_RETURN(LOG_NAS_AMF, AS_NAS_ESTABLISH_RSP_);
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, 0);
}

/****************************************************************************
 **                                                                        **
 ** Name:        amf_send_security_mode_command()                          **
 **                                                                        **
 ** Description: Builds Security Mode Command message.                     **
 **      The Security Mode Command message is sent by the network          **
 **      to the UE to establish NAS signalling security.                   **
 **                                                                        **
 ** Inputs:      msg:     The AMFAS-SAP primitive to process               **
 **              Others:  None                                             **
 **                                                                        **
 ** Outputs:     amf_msg: The AMF message to be sent                       **
 **              Return:  The size of the AMF message                      **
 **              Others:  None                                             **
 **                                                                        **
 ***************************************************************************/
int amf_send_security_mode_command(
    const amf_as_security_t* msg, SecurityModeCommandMsg* amf_msg) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int size = AMF_HEADER_MAXIMUM_LENGTH;
  OAILOG_INFO(
      LOG_NAS_AMF,
      "AMF_TEST: Send Security Mode Command message for ue_id = (%u)\n",
      msg->ue_id);
  /*
   * Mandatory - Message type
   */
  amf_msg->message_type.msg_type = SECURITY_MODE_COMMAND;
  /*
   * Selected NAS security algorithms
   */
  size += NAS5G_SECURITY_ALGORITHMS_MAXIMUM_LENGTH;
  amf_msg->nas_sec_algorithms.M5GNasSecurityAlgorithms_
      .m5gtypeofcipheringalgorithm = M5G_NAS_SECURITY_ALGORITHMS_5G_EA0;
  amf_msg->nas_sec_algorithms.M5GNasSecurityAlgorithms_
      .m5gtypeofintegrityalgorithm = M5G_NAS_SECURITY_ALGORITHMS_5G_IA0;
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, size);
}
}  // namespace magma5g
