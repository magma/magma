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

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#include <thread>
#include "amf_as.h"
#include "nas_network.h"
#include "amf_recv.h"
#include "amf_message.h"
#include "amf_common_defs.h"
#include "3gpp_24.007.h"
#include "log.h"

using namespace std;
#pragma once
typedef uint32_t amf_ue_ngap_id_t;
#define QUADLET 4
#define AMF_GET_BYTE_ALIGNED_LENGTH(LENGTH) LENGTH += QUADLET - (LENGTH % QUADLET)

namespace magma5g
{
       /****************************************************************************
       **                                                                        **
      ** Name:    amf_as_send()                                             **
      **                                                                        **
      ** Description: Processes the AMF-AS Service Access Point primitive.       **
      **                                                                        **
      ** Inputs:  msg:       The AMF-AS-SAP primitive to process         **
      **      Others:    None                                       **
      **                                                                        **
      ** Outputs:     None                                                      **
      **      Return:    RETURNok, RETURNerror                      **
      **      Others:    None                                       **
      **                                                                        **
      ***************************************************************************/
    int amf_as::amf_as_send(amf_sap_t* msg)
    {
        int rc = RETURNok;
        int amf_cause = AMF_CAUSE_SUCCESS;
        amf_as_primitive_t primitive = msg->primitive;
        amf_ue_ngap_id_t ue_id = 0;

        switch (primitive) 
        {
            case _AMFAS_DATA_IND:
              //rc = _amf_as_data_ind(&msg->u.data, &amf_cause);
              //ue_id = msg->u.data.ue_id;
              break;

            case _AMFAS_ESTABLISH_REQ:
              rc = amf_as_establish_req(&msg->m5gu.establish, &amf_cause); //registration request
              ue_id = msg->u.establish.ue_id;
              break;

            case _AMFAS_RELEASE_IND:
              //rc = _amf_as_release_ind(&msg->u.release, &amf_cause);
              //ue_id = msg->u.release.ue_id;
              break;

            default:
              /*
            * Other primitives are forwarded to NGAP
            */
              rc = amf_as::amf_as_send_ng(msg);

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
      static int amf_as_establish_req(amf_as_establish_t *msg, int *amf_cause)
      {
        
        amf_context_t *amf_ctx = NULL;
        amf_security_context_t *AMF_security_context = NULL;
        nas_message_decode_status_t decode_status = {0};
        int decoder_rc = 0;
        int rc = RETURNerror;
        tai_t originating_tai = {0};

        
        nas_message_t nas_msg = {.security_protected.header = {0},
                                .security_protected.plain.amf.header = {0},
                                .security_protected.plain.smf.header = {0}};

        ue_m5gmm_context_s *ue_m5gmm_context = UE_UNREGISTERED;
        //amf_ue_context_exists_amf_ue_ngap_id(msg->ue_id);

        if (ue_mm_context) 
        {
          amf_ctx = &ue_m5gmm_context->amf_context;
          if (amf_ctx) 
          {
            if (IS_AMF_CTXT_PRESENT_SECURITY(AMF_ctx)) 
            {
              AMF_security_context = &amf_ctx->_security;
            }
          }
        } 
        else 
        {
          // some log info
        }

        /*
        * Decode initial NAS message
        */     
        decoder_rc = nas_message_decode(msg->nas_msg->data,&nas_msg,blength(msg->nas_msg),amf_security_context,&decode_status);
        
        nas_network::bdestroy_wrapper(&msg->nas_msg);

        // TODO conditional IE error
        if (decoder_rc < 0) 
        {
          if (decoder_rc < TLV_FATAL_ERROR) 
          {
            *amf_cause = AMF_CAUSE_PROTOCOL_ERROR;            
          }
          else if (decoder_rc == TLV_MANDATORY_FIELD_NOT_PRESENT) 
          {
            *amf_cause = AMF_CAUSE_INVALID_MANDATORY_INFO;            
          }
          else if (decoder_rc == TLV_UNEXPECTED_IEI) 
          {
            *amf_cause = AMF_CAUSE_IE_NOT_IMPLEMENTED;           
          } 
          else 
          {
            *amf_cause = AMF_CAUSE_PROTOCOL_ERROR;
          }
        }

        /*
        * Process initial NAS message
        */
       
        AMFMsg *amf_msg = &nas_msg.plain.emm;

        switch (amf_msg->header.message_type) 
        {
          case REGISTRATION_REQUEST:
            memcpy(&originating_tai, msg->tai, sizeof(originating_tai));
            rc = amf_procedure_handler::amf_handle_registration_request(msg->ue_id,&originating_tai,&msg->ecgi,&amf_msg->registration_request,msg->is_initial,msg->is_mm_ctx_new,amf_cause,&decode_status);
            break;
          // more case to wright......
          
        }

        
      }
      /****************************************************************************
       **                                                                        **
      ** Name:    amf_as_send_ng()                                            **
      **                                                                        **
      ** Description: Builds NAS message according to the given AMFAS Service   **
      **      Access Point primitive and sends it to the Access Stratum **
      **      sublayer                                                  **
      **                                                                        **
      ** Inputs:  msg:       The AMFAS-SAP primitive to be sent         **
      **      Others:    None                                       **
      **                                                                        **
      ** Outputs:     None                                                      **
      **      Return:    RETURNok, RETURNerror                      **
      **      Others:    None                                       **
      **                                                                        **
      ***************************************************************************/
      int amf_as::amf_as_send_ng(const amf_as_t *msg)
      {
        OAILOG_FUNC_IN(LOG_NAS_AMF);
        as_message_t as_msg = {0};

        switch (msg->primitive) 
        {
          case _AMFAS_DATA_REQ:
            as_msg.msg_id =amf_as_data_req(&msg->u.data, &as_msg.msg.dl_info_transfer_req);
            break;

          case _AMFAS_ESTABLISH_CNF:
            as_msg.msg_id = _amf_as_establish_cnf(&msg->u.establish, &as_msg.msg.nas_establish_rsp);
            break;

          case _AMFAS_ESTABLISH_REJ:
            as_msg.msg_id =_amf_as_establish_rej(&msg->u.establish, &as_msg.msg.nas_establish_rsp);
            break;
            

           // more case to wright......

          default: as_msg.msg_id = 0; break;
        }

        /*
        * Send the message to the Access Stratum or NGAP in case of AMF
        */
        if (as_msg.msg_id > 0) 
        {
          OAILOG_DEBUG(LOG_NAS_AMF,"AMFAS-SAP - ""Sending msg with id 0x%x, primitive %s (%d) to NGAP layer for " "transmission\n", as_msg.msg_id, _amf_as_primitive_str[msg->primitive - _AMFAS_START - 1], msg->primitive);

          switch (as_msg.msg_id) 
          {
            case AS_DL_INFO_TRANSFER_REQ: 
            {
                 amf_app_defs::amf_app_handle_nas_dl_req(as_msg.msg.dl_info_transfer_req.ue_id,as_msg.msg.dl_info_transfer_req.nas_msg,as_msg.msg.dl_info_transfer_req.err_code);
                OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
            } break;

           case AS_NAS_ESTABLISH_RSP:
           case AS_NAS_ESTABLISH_CNF: 
            {
              if (as_msg.msg.nas_establish_rsp.err_code == AS_SUCCESS) 
              {
                // This flow is to release the UE context after sending the NAS message.
                  amf_app_defs::amf_app_handle_nas_dl_req(as_msg.msg.nas_establish_rsp.ue_id, as_msg.msg.nas_establish_rsp.nas_msg, as_msg.msg.nas_establish_rsp.err_code);
                  as_msg.msg.nas_establish_rsp.nas_msg = NULL;
                  OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
              } 
              else 
              {
                  OAILOG_DEBUG(LOG_NAS_AMF, "AMFAS-SAP - Sending establish_cnf to AMF-APP module for UE ID: "
                  AMF_UE_NGAP_ID_FMT " selected eea ""0x%04X selected eia 0x%04X\n",
                  as_msg.msg.nas_establish_rsp.ue_id, as_msg.msg.nas_establish_rsp.selected_encryption_algorithm,
                  as_msg.msg.nas_establish_rsp.selected_integrity_algorithm);
                /*
                * Handle success case
                */
                  amf_app_handle_conn_est_cnf(&as_msg.msg.nas_establish_rsp);
                  OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
              }
            } break;

            case AS_NAS_RELEASE_REQ:
              amf_app_handle_deregister_req(as_msg.msg.nas_release_req.ue_id);
              OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
              break;

            default: break;
          }
        }

        OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
      }
      /****************************************************************************
       **                                                                        **
      ** Name:    amf_as_encode()                                          **
      **                                                                        **
      ** Description: Encodes NAS message into NAS information container        **
      **                                                                        **
      ** Inputs:  msg:       The NAS message to encode                  **
      **      length:    The maximum length of the NAS message      **
      **      Others:    None                                       **
      **                                                                        **
      ** Outputs:     info:      The NAS information container              **
      **      msg:       The NAS message to encode                  **
      **      Return:    The number of bytes successfully encoded   **
      **      Others:    None                                       **
      **                                                                        **
      ***************************************************************************/
      static int amf_as_encode(bstring *info, nas_message_t *msg, size_t length, emm_security_context_t *amf_security_context)
      {
        OAILOG_FUNC_IN(LOG_NAS_AMF);
        int bytes = 0;

        /* Ciphering algorithms, EEA1 and EEA2 expects length to be mode of 4,
        * so length is modified such that it will be mode of 4
        */
        AMF_GET_BYTE_ALIGNED_LENGTH(length);
        if (msg->header.security_header_type != SECURITY_HEADER_TYPE_NOT_PROTECTED) 
        {
          amf_msg_header *header = &msg->security_protected.plain.amf.header;

          /*
          * Expand size of protected NAS message
          */
          length += NAS_MESSAGE_SECURITY_HEADER_SIZE;
          /*
          * Set header of plain NAS message
          */
          header->extended_protocol_discriminator = M5GS_MOBILITY_MANAGEMENT_MESSAGE;
          header->security_header_type = SECURITY_HEADER_TYPE_NOT_PROTECTED;
        }

        /*
        * Allocate memory to the NAS information container
        */
        *info = bfromcstralloc(length, "\0");

        if (*info) {
          /*
          * Encode the NAS message
          */
          bytes = nas_message_encode((*info)->data, msg, length, amf_security_context);

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
      static int amf_as_data_req(const amf_as_data_t *msg, dl_info_transfer_req_t *as_msg)
      {
        OAILOG_FUNC_IN(LOG_NAS_AMF);
        int size = 0;
        int is_encoded = false;

        OAILOG_INFO(LOG_NAS_AMF, "AMFAS-SAP - Send AS data transfer request\n");
        /*nas_message_t nas_msg = {.security_protected.header = {0},
                                .security_protected.plain.emm.header = {0},
                                .security_protected.plain.esm.header = {0}};

        /*
        * Setup the AS message
        */
        if (msg->guti) {
          as_msg->s_tmsi.amf_code = msg->guti->gummei.amf_code;
          as_msg->s_tmsi.m_tmsi = msg->guti->m_tmsi;
        } 
        else 
        {
          as_msg->ue_id = msg->ue_id;
        }

        /*
        * Setup the NAS security header
        */
        AMFMsg *amf_msg = amf_as_set_header(&nas_msg, &msg->sctx); // all header part==> all mendatory field

        /*
        * Setup the NAS information message
        */
        if (amf_msg)switch (msg->nas_info) 
        {
            case AMF_AS_NAS_DATA_REGISTRATION_ACCEPT:
              size = amf_send_registration_accept_dl_nas(msg, &amf_msg->registration_accept); // make the contents of registration accept message
              break;
            case AMF_AS_NAS_DL_NAS_TRANSPORT:
               size= amf_send_dl_nas_transportmsg(msg, &amf_msg->downlink_nas_transport)
              // many more case ....
            default:
              /*
            * Send other NAS messages as already encoded SMF messages
            */
              size = msg->nas_msg->slen;
              //is_encoded = true;
              break;
          }

        if (size > 0) 
        {
          int bytes = 0;
          amf_security_context_t* amf_security_context = NULL;
          class amf_context_s* amf_ctx = NULL;
          ue_m5gmm_context_s* ue_m5gmm_context = amf_ue_context_exists_amf_ue_ngap_id(msg->ue_id);

          if (ue_m5gmm_context) 
          {
            amf_ctx = &ue_m5gmm_context->amf_context;
            if (amf_ctx) 
            {
              if (amf_msg->nw_deregister_request.nw_deregistertype ==NW_DEREGISTER_TYPE_IMSI_DEREGISTER) 
              {
                amf_ctx->is_imsi_only_deregister = true;
              }
              if (IS_AMF_CTXT_PRESENT_SECURITY(amf_ctx)) 
              {
                amf_security_context = &amf_ctx->_security;
                is_encoded = true;
              }
            }
          }

          if (amf_security_context) 
          {
            nas_msg.header.sequence_number = amf_security_context->dl_count.seq_num;
            OAILOG_DEBUG(LOG_NAS_AMF,"Set nas_msg.header.sequence_number -> %u\n", nas_msg.header.sequence_number);
          } 
          else 
          {
            OAILOG_ERROR(LOG_NAS_AMF, "Security context is NULL for UE -> %d\n", msg->ue_id);
            OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
          }

          if (!is_encoded) 
          {
            /*
            * Encode the NAS information message
            */
            bytes =amf_as_encode(&as_msg->nas_msg, &nas_msg, size, amf_security_context);
          } 
          else 
          {
            /*
            * Encrypt the NAS information message
            */
           // bytes = amf_as_encrypt(&as_msg->nas_msg,&nas_msg.header, msg->nas_msg->data, size,amf_security_context);
          }

          // Free any allocated data
          switch (msg->nas_info) 
          {
            // amf_information message and downlink_nas_transtport is the only message that has allocated data
            case AMF_AS_NAS_DATA_REGISTRATION_ACCEPT:
              nas_network::bdestroy_wrapper(&(amf_msg->registration_accept.smfmessagecontainer));
              break;
              // many more remain....
          }

          if (bytes > 0) 
          {
            as_msg->err_code = AS_SUCCESS;            
          }
            OAILOG_FUNC_RETURN(LOG_NAS_AMF, AS_DL_INFO_TRANSFER_REQ);
          
        }

        OAILOG_FUNC_RETURN(LOG_NAS_AMF, 0);
      }

    /****************************************************************************
     **                                                                        **
    ** Name:    amf_as_set_header()                                      **
    **                                                                        **
    ** Description: Setup the security header of the given NAS message        **
    **                                                                        **
    ** Inputs:  security:  The NAS security data to use               **
    **      Others:    None                                       **
    **                                                                        **
    ** Outputs:     msg:       The NAS message                            **
    **      Return:    Pointer to the plain NAS message to be se- **
    **             curity protected if setting of the securi- **
    **             ty header succeed;                         **
    **             NULL pointer otherwise                     **
    **      Others:    None                                       **
    **                                                                        **
    ***************************************************************************/
    static AMFMsg* amf_as_set_header(nas_message_t* msg, 
                                    const emm_as_security_data_t* security) {
      OAILOG_FUNC_IN(LOG_NAS_AMF);
      msg->header.extended_protocol_discriminator = M5GS_MOBILITY_MANAGEMENT_MESSAGE;

      if (security && (security->ksi != KSI_NO_KEY_AVAILABLE)) {
        /*
        * A valid EPS security context exists
        */
        if (security->is_new) {
          /*
          * New EPS security context is taken into use
          */
          if (security->is_knas_int_present) {
            if (security->is_knas_enc_present) {
              /*
              * NAS integrity and cyphering keys are available
              */
              msg->header.security_header_type = SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED_NEW;
            } else {
              /*
              * NAS integrity key only is available
              */
              msg->header.security_header_type = SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_NEW;
            }

            OAILOG_FUNC_RETURN(LOG_NAS_AMF, &msg->security_protected.plain.amf);
          }
        } else if (security->is_knas_int_present) {
          if (security->is_knas_enc_present) {
            /*
            * NAS integrity and cyphering keys are available
            */
            msg->header.security_header_type = SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED;
          } else {
            /*
            * NAS integrity key only is available
            */
            msg->header.security_header_type = SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED;
          }

          OAILOG_FUNC_RETURN(LOG_NAS_AMF, &msg->security_protected.plain.amf);
        } else {
          /*
          * No valid EPS security context exists
          */
          msg->header.security_header_type = SECURITY_HEADER_TYPE_NOT_PROTECTED;
          OAILOG_FUNC_RETURN(LOG_NAS_AMF, &msg->plain.amf);
        }
      } else {
        /*
        * No valid EPS security context exists
        */
        msg->header.security_header_type = SECURITY_HEADER_TYPE_NOT_PROTECTED;
        OAILOG_FUNC_RETURN(LOG_NAS_AMF, &msg->plain.amf);
      }

      /*
      * A valid EPS security context exists but NAS integrity key
      * * * * is not available
      */
      OAILOG_FUNC_RETURN(LOG_NAS_AMF, NULL);
    }
 
}
    
    