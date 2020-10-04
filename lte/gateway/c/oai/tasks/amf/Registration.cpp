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

  Source      Registration.cpp

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#include "amf_proc.h"
#include "amf_common_defs.h"
#include "amf_app_ue_context.h"
#include "amf_fsm.h"
#include "amf_data.h"
#include "nas_proc.h"
#include "amf_sap.hs"
#include "amf_msgdef.h"
#include "registration_accept.h"
#include "3gpp_24.501.h"
#include "5GSRegistrationResult.h"
using namespace std;
#pragma once
#define INVALID_IMSI64 (imsi64_t) 0
#define INVALID_AMF_UE_NGAP_ID 0x0

namespace magma5g
{
     static int amf_proc_registration_request(amf_ue_ngap_id_t ue_id, const bool is_mm_ctx_new,
                                                amf_registration_request_ies_t* const ies) {
        int rc = RETURNerror;
        ue_m5gmm_context_s ue_ctx;
        amf_fsm_state_t fsm_state = AMF_DEREGISTERED;
        bool clear_amf_ctxt = false;
        ue_m5gmm_context_s *ue_m5gmm_context = NULL;
        ue_m5gmm_context_s *guti_ue_m5gmm_ctx = NULL;
        ue_m5gmm_context_s *imsi_ue_m5gmm_ctx = NULL;
        amf_context_t *new_amf_ctx = NULL;
        imsi64_t imsi64 = INVALID_IMSI64;
        amf_ue_ngap_id_t old_ue_id = INVALID_AMF_UE_NGAP_ID;

        if (ies->imsi) 
        {
            imsi64 = imsi_to_imsi64(ies->imsi);
            OAILOG_INFO(LOG_NAS_AMF,"REGISTRATION REQ (ue_id = " AMF_UE_NGAP_ID_FMT ") (IMSI = " IMSI_64_FMT ") \n",ue_id,imsi64);
        } 
        else if (ies->guti) 
        {
            OAILOG_INFO(LOG_NAS_AMF,"REGISTRATION REQ (ue_id = " AMF_UE_NGAP_ID_FMT ") (GUTI = " GUTI_FMT ") \n",ue_id,GUTI_ARG(ies->guti));
        } 
        else if (ies->imei) 
        {
            char imei_str[16];
            IMEI_TO_STRING(ies->imei, imei_str, 16);
            OAILOG_INFO(LOG_NAS_AMF,"REGISTRATION REQ (ue_id = " AMF_UE_NGAP_ID_FMT ") (IMEI = %s ) \n",ue_id,imei_str);
        }

        OAILOG_INFO(LOG_NAS_AMF,"AMF-PROC:  Registration -Registration type = %s (%d)\n", _amf_registration_type_str[ies->type], ies->type);
        OAILOG_DEBUG(LOG_NAS_AMF,"is_initial request = %u\n (ue_id=" AMF_UE_NGAP_ID_FMT ") \n(imsi = " IMSI_64_FMT ") \n", ies->is_initial,ue_id,imsi64);
        /*
        * Initialize the temporary UE context
        */
        memset(&ue_ctx, 0, sizeof(ue_amf_context_t));
        ue_ctx.amf_context.is_dynamic = false;
        ue_ctx.amf_ue_ngap_id = ue_id;

        rc = amf_registration_run_procedure(&ue_mm_context->amf_context);
        OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);


    }
   
//------------------------------------------------------------------------------
int amf_proc_registration_reject(amf_ue_ngap_id_t ue_id, amf_cause_t amf_cause) 
{
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;

  amf_context_t* amf_ctx = amf_context_get(&_amf_data, ue_id);
  if (amf_ctx) {
    if (is_nas_specific_procedure_registration_running(amf_ctx)) {
      nas_amf_registration_proc_t* registration_proc =
          (nas_amf_registration_proc_t*) (amf_ctx->amf_procedures->amf_specific_proc);
      registration_proc->amf_cause = amf_cause;

      // TODO could be in callback of attach procedure triggered by
      // AMF REG__REJ
      rc = amf_registration_reject(amf_ctx, (struct nas_base_proc_s*) registration_proc);
      amf_sap_t amf_sap               = {0};
      amf_sap.primitive               = AMFREG_REGISTRATION_REJ;
      amf_sap.u.amf_reg.ue_id         = ue_id;
      amf_sap.u.amf_reg.ctx           = amf_ctx;
      amf_sap.u.amf_reg.notify        = false;
      amf_sap.u.amf_reg.free_proc     = true;
      amf_sap.u.amf_reg.u.attach.proc = registration_proc;
      rc                              = emm_sap_send(&emm_sap);
    } else {
      nas_amf_registration_proc_t no_registration_proc = {0};
      no_registration_proc.ue_id                 = ue_id;
      no_registration_proc.amf_cause             = amf_cause;
      no_registration_proc.amf_msg_out           = NULL;
      rc =amf_registration_reject(amf_ctx, (struct nas_base_proc_s*) &no_registration_proc);
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}
int amf_registration_reject(amf_context_t* amf_context, struct nas_base_proc_s* nas_base_proc) 
{
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;

  amf_sap_t amf_sap = {0};
  nas_amf_registration_proc_t* registration_proc = ( nas_amf_registration_proc_t*) nas_base_proc;

  OAILOG_WARNING(LOG_NAS_AMF,"AMF-PROC  - AMF Registration procedure not accepted "
      "by the network (ue_id=" AMF_UE_NGAP_ID_FMT ", cause=%d)\n",
      registration_proc->ue_id, registration_proc->amf_cause);
  /*
   * Notify AMF-AS SAP that Registration Reject message has to be sent
   * onto the network
   */
  amf_sap.primitive               = AMFREG_REGISTRATION_REJ;
  amf_sap.u.amf_as.u.establish.ue_id       = registration_proc->ue_id;
  amf_sap.u.amf_as.u.establish.amf_cause = registration_proc->amf_cause;
  amf_sap.u.amf_as.u.establish.nas_info  = AMF_AS_NAS_INFO_REGISTRATION;

  if (registration_proc->amf_cause != AMF_CAUSE_SMF_FAILURE) {
    amf_sap.u.amf_as.u.establish.nas_msg = NULL;
  } else if (registration_proc->amf_msg_out) {
    amf_sap.u.amf_as.u.establish.nas_msg = registration_proc->amf_msg_out;
  } else {
    OAILOG_ERROR(LOG_NAS_EMM, "AMF-PROC  - SMF message is missing\n");
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }

  /*
   * Setup 5G CN NAS security data
   */
  if (amf_context) {
    amf_as_set_security_data( &amf_sap.u.amf_as.u.establish.sctx, &amf_context->_security, false,  false);
  } else {
    amf_as_set_security_data(&amf_sap.u.amf_as.u.establish.sctx, NULL, false, false);
  }
  rc = amf_sap_send(&amf_sap);
  increment_counter("ue_Registration", 1, 1, "action", "Registration_reject_sent");
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}
    /*
    * --------------------------------------------------------------------------
    * Functions that may initiate AMF common procedures
    * --------------------------------------------------------------------------
    */

    //------------------------------------------------------------------------------
    static int amf_registration_run_procedure(amf_context_t *amf_context)
    {
        OAILOG_FUNC_IN(LOG_NAS_AMF);
        int rc = RETURNerror;
        nas_amf_registration_proc_t *registration_proc = nas_proc::get_nas_specific_procedure_registration(amf_context);

        if (registration_proc) 
        {
            

            if (registration_proc->ies->last_visited_registered_tai)
            {
            // amf_ctx_set_valid_lvr_tai(amf_context, registration_proc->ies->last_visited_registered_tai);
                //amf_ctx_set_valid_ue_nw_cap(amf_context, &registration_proc->ies->ue_network_capability);

            }
                
           // if (registration_proc->ies->ms_network_capability) 
           // {
            //amf_ctx_set_valid_ms_nw_cap(amf_context, registration_proc->ies->ms_network_capability);
           // }
            //amf_context->originating_tai = *registration_proc->ies->originating_tai;

            
            // temporary choice to clear security context if it exist
            //amf_ctx_clear_security(amf_context);

            if (registration_proc->ies->imsi) {
            if (
                (registration_proc->ies->decode_status.mac_matched) ||
                !(registration_proc->ies->decode_status.integrity_protected_message)) {
                // force authentication, even if not necessary
                imsi64_t imsi64 = imsi_to_imsi64(registration_proc->ies->imsi);
                amf_ctx_set_valid_imsi(amf_context, registration_proc->ies->imsi, imsi64);
               //TODO amf_context_upsert_imsi(&_amf_data, amf_context);
                rc = amf_start_registration_proc_authentication(amf_context, registration_proc);
                if (rc != RETURNok) {
                    OAILOG_ERROR(LOG_NAS_AMF, "Failed to start registration authentication procedure!\n");
                }
             } else {
                //force identification, even if not necessary
                rc = amf_proc_identification(amf_context,(nas_amf_proc_t *) registration_proc, IDENTITY_TYPE_2_IMSI,_amf_registration_success_identification_cb,amf_registration_failure_identification_cb);
              }
            } 
            else if (registration_proc->ies->guti) 
            {
            rc = amf_proc_identification(amf_context,(nas_amf_proc_t *) registration_proc, IDENTITY_TYPE_2_IMSI,_amf_registration_success_identification_cb,amf_registration_failure_identification_cb);
            }
            else if (registration_proc->ies->imei) 
            {
                // emergency allowed if go here, but have to be implemented...
                AssertFatal(0, "TODO emergency");
            }
        }
        amf_registration_procedure::amf_registration_success_identification_cb(amf_context);
        OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
    }
    //------------------------------------------------------------------------------
    int amf_registration_procedure::amf_registration_success_identification_cb(amf_context_t *amf_context)
    {
        OAILOG_FUNC_IN(LOG_NAS_AMF);
        int rc = RETURNerror;

        OAILOG_INFO(LOG_NAS_AMF, "registration - Identification procedure success!\n");
        nas_amf_registration_proc_t *registration_proc = get_nas_specific_procedure_registration(amf_context);

        if (registration_proc) 
        {
            
            rc = amf_start_registration_proc_authentication(amf_context,registration_proc); 
            
        }
        OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
    }

    //------------------------------------------------------------------------------
      static int _amf_registration_failure_identification_cb(amf_context_t *amf_context)
    {
        OAILOG_FUNC_IN(LOG_NAS_AMF);
        int rc = RETURNerror;

        OAILOG_ERROR(LOG_NAS_AMF, "registration - Identification procedure failed!\n");

        AssertFatal(0, "Cannot happen...\n");
        OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
    }

    //------------------------------------------------------------------------------
    static int amf_start_registration_proc_authentication(amf_context_t *amf_context, nas_amf_registration_proc_t *registration_proc)
    {
        OAILOG_FUNC_IN(LOG_NAS_AMF);
        int rc = RETURNerror;

        if ((amf_context) && (registration_proc)) 
        {
           rc = authentication::amf_proc_authentication(amf_context,&registration_proc->amf_spec_proc,amf_registration_success_authentication_cb, amf_registration_failure_authentication_cb);
        }
        //amf_registration_success_authentication_cb(amf_context);
        OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
    }

    //------------------------------------------------------------------------------
    static int amf_registration_success_authentication_cb(amf_context_t *amf_context)
    {
        OAILOG_FUNC_IN(LOG_NAS_AMF);
        int rc = RETURNerror;

        OAILOG_INFO( LOG_NAS_AMF, "REGISTRATION - Authentication procedure success!\n");
        nas_amf_registration_proc_t *registration_proc =  get_nas_specific_procedure_registration(amf_context);

        if (registration_proc) 
        {
            //REQUIREMENT_3GPP_24_501(R15_5_5_1_2_3__1);
            rc = amf_start_registration_proc_security(amf_context, registration_proc);
        }
        OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
    }

    //------------------------------------------------------------------------------
     static int amf_registration_failure_authentication_cb(amf_context_t *amf_context)
    {
        OAILOG_FUNC_IN(LOG_NAS_AMF);
        int rc = RETURNerror;
        OAILOG_ERROR(LOG_NAS_AMF, "REGISTRATION - Authentication procedure failed!\n");
        nas_amf_registration_proc_t *registration_proc = get_nas_specific_procedure_registration(amf_context);

        if (registration_proc) 
        {
            registration_proc->amf_cause = amf_context->amf_cause;

            amf_sap_t amf_sap = {0};
            amf_sap.primitive = EMMREG_REGISTRATION_REJ;
            amf_sap.u.amf_reg.ue_id = registration_proc->ue_id;
            amf_sap.u.amf_reg.ctx = amf_context;
            amf_sap.u.amf_reg.notify = true;
            amf_sap.u.amf_reg.free_proc = true;
            amf_sap.u.amf_reg.u.registration.proc = registration_proc;
            // dont' care amf_sap.u.amf_reg.u.registration.is_emergency = false;
            rc = amf_sap_send(&amf_sap);
        }
        OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
    }

    //------------------------------------------------------------------------------
    static int amf_start_registration_proc_security(amf_context_t *amf_context, nas_amf_registration_proc_t *registration_proc)
    {
        OAILOG_FUNC_IN(LOG_NAS_AMF);
        int rc = RETURNerror;

        if ((amf_context) && (registration_proc)) 
        {
            //REQUIREMENT_3GPP_24_501(R15_5_5_1_2_3__1);
            amf_ue_ngap_id_t ue_id = PARENT_STRUCT(amf_context, struct ue_m5gmm_context_s, amf_context)->amf_ue_ngap_id;
            
            /*
            * Create new NAS security context
            */
            //amf_ctx_clear_security(amf_context);
            rc = amf_proc_security_mode_control(amf_context, &registration_proc->amf_spec_proc,registration_proc->ksi,_amf_registration_success_security_cb,_amf_registration_failure_security_cb);
            amf_registration_success_security_cb(amf_context);
        }
        OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
    }

    //------------------------------------------------------------------------------
    static int amf_registration_success_security_cb(amf_context_t *amf_context)
    {
        OAILOG_FUNC_IN(LOG_NAS_AMF);
        int rc = RETURNerror;

        OAILOG_INFO(LOG_NAS_AMF, "REGISTRATION - Security procedure success!\n");
        nas_amf_registration_proc_t *registration_proc = get_nas_specific_procedure_registration(amf_context);

        if (registration_proc) 
        {
            rc = amf_registration(amf_context);
        }
        OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
    }

    //------------------------------------------------------------------------------
     static int _amf_registration_failure_security_cb(amf_context_t *amf_context)
    {
        OAILOG_FUNC_IN(LOG_NAS_AMF);
        int rc = RETURNerror;
        OAILOG_ERROR(LOG_NAS_AMF, "REGISTRATION - Security procedure failed!\n");
        nas_amf_registration_proc_t *registration_proc = get_nas_specific_procedure_registration(amf_context);

        if (registration_proc) 
        {
            _amf_registration_release(amf_context);
        }
        OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
    }


    /*
    *
    * Name:        amf_registration_security()
    *
    * Description: Initiates security mode control AMF common procedure.
    *
    * Inputs:          args:      security argument parameters
    *                  Others:    None
    *
    * Outputs:     None
    *                  Return:    RETURNok, RETURNerror
    *                  Others:    _amf_data
    */
    //------------------------------------------------------------------------------
    int amf_registration_security(struct amf_context_s *amf_context)
    {
    return _amf_registration_security(amf_context);
    
    }


    /*
    --------------------------------------------------------------------------
                    AMF specific local functions
    --------------------------------------------------------------------------
    */

    /*
    *
    * Name:    amf_registration()
    *
    * Description: Performs the registration signalling procedure while a context
    *      exists for the incoming UE in the network.
    *
    *              3GPP TS 24.501, section 5.5.1.2.4
    *      Upon receiving the REGISTRATION REQUEST message, the AMF shall
    *      send an REGISTRATION ACCEPT message to the UE and start timer
    *      T3450.
    *
    * Inputs:  args:      registration argument parameters
    *      Others:    None
    *
    * Outputs:     None
    *      Return:    RETURNok, RETURNerror
    *      Others:    _amf_data
    *
    */
    //------------------------------------------------------------------------------
    static int amf_registration(amf_context_t *amf_context)
    {
        OAILOG_FUNC_IN(LOG_NAS_AMF);
        int rc = RETURNerror;
        amf_ue_ngap_id_t ue_id =
            PARENT_STRUCT(amf_context, struct ue_m5gmm_context_s, amf_context)
            ->amf_ue_ngap_id;

        OAILOG_INFO(LOG_NAS_AMF,"ue_id=" AMF_UE_NGAP_ID_FMT " AMF-PROC  - REGISTRATION UE \n", ue_id);

        nas_amf_registration_proc_t * get_nas_specific_procedure_registrationistration_proc =  get_nas_specific_procedure_registration(amf_context);

        //if (registration_proc) 
        {
            if (registration_proc->ies->smf_msg) 
            {
                smf_sap_t smf_sap = {0};
                smf_sap.primitive = SMF_UNITDATA_IND;
                smf_sap.is_standalone = false;
                smf_sap.ue_id = ue_id;
                smf_sap.ctx = amf_context;
                smf_sap.recv = registration_proc->ies->smf_msg;
                //rc = smf_sap_send(&smf_sap);
                if ((rc != RETURNerror) && (smf_sap.err == SMF_SAP_SUCCESS)) 
                {
                    rc = RETURNok;
                } 
                /*else if (smf_sap.err != SMF_SAP_DISCARDED) 
                {
                    /*
                    * Theregistration procedure failed due to an SMF procedure failure
                    */
                    registration_proc->amf_cause = AMF_CAUSE_SMF_FAILURE;

                    /*
                    * Setup the SMF message container to include pdu session Connectivity Reject
                    * message within the REGISTRATION Reject message
                    */
                   /* bdestroy_wrapper(&registration_proc->ies->smf_msg);
                    registration_proc->smf_msg_out = smf_sap.send;
                    OAILOG_ERROR(LOG_NAS_AMF, "Sending Registration Reject to UE ue_id = (%u), amf_cause = (%d)\n", ue_id, registration_proc->amf_cause);
                    rc = _amf_registration_reject(amf_context, &registration_proc->amf_spec_proc.amf_proc.base_proc);
                }*/
                else 
                {
                    /*
                    * SMF procedure failed and, received message has been discarded or
                    * Status message has been returned; ignore SMF procedure failure
                    */
                    OAILOG_WARNING(LOG_NAS_AMF,"Ignore SMF procedure failure &""received message has been discarded for ue_id = (%u)\n", ue_id);
                    rc = RETURNok;
                }
            } 
            else 
            {
                rc = RETURNok;
                rc = amf_send_registration_accept(amf_context);
            }
        }

        if (rc != RETURNok) 
        {
            /*
            * The registration procedure failed
            */
            OAILOG_ERROR(LOG_NAS_AMF,"ue_id=" AMF_UE_NGAP_ID_FMT " AMF-PROC  - Failed to respond to Registration Request\n", ue_id);
            registration_proc->amf_cause = AMF_CAUSE_PROTOCOL_ERROR;
            /*
            * Do not accept the UE to registration to the network
            */
            OAILOG_ERROR(LOG_NAS_AMF,"Sending Registration Reject to UE ue_id = (%u), amf_cause = (%d)\n", ue_id, registration_proc->amf_cause);
           // rc = _amf_registration_reject(amf_context, &registration_proc->amf_spec_proc.amf_proc.base_proc);
            //increment_counter("ue_registration", 1, 2, "result", "failure", "cause", "protocol_error");
        }

        OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
    }

    /****************************************************************************
     **                                                                        **
    ** Name:    amf_send_registration_accept()                                      **
    **                                                                        **
    ** Description: Sends REGISTRATION ACCEPT message and start timer T3450         **
    **                                                                        **
    ** Inputs:  data:      Registration accept retransmission data          **
    **      Others:    None                                       **
    **                                                                        **
    ** Outputs:     None                                                      **
    **      Return:    RETURNok, RETURNerror                      **
    **      Others:    T3450                                      **
    **                                                                        **
    ***************************************************************************/
    static int amf_send_registration_accept(amf_context_t *amf_context)
    {
        OAILOG_FUNC_IN(LOG_NAS_AMF);
        int rc = RETURNerror;

        // may be caused by timer not stopped when deleted context
        if (amf_context) 
        {
            amf_sap_t amf_sap = {0};
            nas_amf_registration_proc_t *registration_proc = get_nas_specific_procedure_registration(amf_context);
            ue_m5gmm_context_s *ue_m5gmm_context_p = PARENT_STRUCT(amf_context, class ue_m5gmm_context_s, amf_context);
            amf_ue_ngap_id_t ue_id = ue_m5gmm_context_p->amf_ue_ngap_id;

            if (registration_proc) 
            {
                _amf_registration_update(amf_context, registration_proc->ies);
                /*
                * Notify AMF-AS SAP that Registaration Accept message together with an Activate
                * Pdu session Context Request message has to be sent to the UE
                */
                amf_sap.primitive = AMFAS_ESTABLISH_CNF;
                amf_sap.u.amf_as.u.establish.puid = registration_proc->amf_spec_proc.amf_proc.base_proc.nas_puid;
                amf_sap.u.amf_as.u.establish.ue_id = ue_id;
                amf_sap.u.amf_as.u.establish.nas_info = AMF_AS_NAS_INFO_REGISTRATION;

                //NO_REQUIREMENT_3GPP_24_501(R15_5_5_1_2_4__3);
                //bdestroy_wrapper(&ue_m5gmm_context_s->ue_radio_capability);
                //----------------------------------------
                //REQUIREMENT_3GPP_24_501(R15_5_5_1_2_4__4);
                //amf_ctx_set_attribute_valid(amf_context, AMF_CTXT_MEMBER_UE_NETWORK_CAPABILITY_IE);
                //amf_ctx_set_attribute_valid(amf_context, AMF_CTXT_MEMBER_MS_NETWORK_CAPABILITY_IE);
                //----------------------------------------
                if (attach_proc->ies->drx_parameter) 
                {
                    //REQUIREMENT_3GPP_24_501(R15_5_5_1_2_4__5);
                    amf_ctx_set_valid_drx_parameter(amf_context, registration_proc->ies->drx_parameter);
                }
                //----------------------------------------
                //REQUIREMENT_3GPP_24_501(R15_5_5_1_2_4__9);
                // the set of amf_sap.u.amf_as.u.establish.new_guti is for including the GUTI in the registration accept message
                //ONLY ONE MME NOW NO S10
                /*if (!IS_AMF_CTXT_PRESENT_GUTI(amf_context)) 
                {
                    // Sure it is an unknown GUTI in this AMF
                    guti_t old_guti = amf_context->_old_guti;
                    guti_t guti = {.gummei.plmn = {0},
                                .gummei.mme_gid = 0,
                                .gummei.mme_code = 0,
                                .m_tmsi = INVALID_M_TMSI};
                    clear_guti(&guti);

                    rc = amf_api_new_guti(
                    &amf_context->_imsi,
                    &old_guti,
                    &guti,
                    &amf_context->originating_tai,
                    &amf_context->_tai_list);
                    if (RETURNok == rc) {
                    amf_ctx_set_guti(amf_context, &guti);
                    amf_ctx_set_attribute_valid(amf_context, AMF_CTXT_MEMBER_TAI_LIST);
                    //----------------------------------------
                    REQUIREMENT_3GPP_24_501(R15_5_5_1_2_4__6);
                    REQUIREMENT_3GPP_24_501(R15_5_5_1_2_4__10);
                    memcpy(
                        &amf_sap.u.amf_as.u.establish.tai_list,
                        &amf_context->_tai_list,
                        sizeof(tai_list_t));
                    }
                    else 
                    {
                        OAILOG_ERROR(LOG_NAS_AMF,"Failed to assign amf api new guti for ue_id = %u\n",ue_id);
                        OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
                    }
                } 
                else 
                {
                    // Set the TAI attributes from the stored context for resends.
                    memcpy(
                    &amf_sap.u.amf_as.u.establish.tai_list,
                    &amf_context->_tai_list,
                    sizeof(tai_list_t));
                }*/
            }

            //amf_sap.u.amf_as.u.establish.eps_id.guti = &amf_context->_guti;

            /*if (!IS_AMF_CTXT_VALID_GUTI(amf_context) && IS_AMF_CTXT_PRESENT_GUTI(amf_context) &&   IS_AMF_CTXT_PRESENT_OLD_GUTI(amf_context)) 
            {
                /*
                * Implicit GUTI reallocation;
                * include the new assigned GUTI in the Registration Accept message
                */
                /*OAILOG_DEBUG(LOG_NAS_AMF, "ue_id=" AMF_UE_NGAP_ID_FMT " AMF-PROC  - Implicit GUTI reallocation, include the new assigned " "GUTI in the Registration Accept message\n", ue_id);
                amf_sap.u.amf_as.u.establish.new_guti = &amf->_guti;
            } */
            /*else if (!IS_AMF_CTXT_VALID_GUTI(amf_context) && IS_AMF_CTXT_PRESENT_GUTI(amf_context)) 
            {
                /*
                * include the new assigned GUTI in the Attach Accept message
                */
                /*OAILOG_DEBUG(LOG_NAS_AMF,"ue_id=" AMF_UE_NGAP_ID_FMT " AMF-PROC  - Include the new assigned GUTI in the Registration Accept ""message\n", ue_id);
                amf_sap.u.amf_as.u.establish.new_guti = &amf_context->_guti;
            } */
            //else 
            { // IS_AMF_CTXT_VALID_GUTI(ue_mm_context) is true
                amf_sap.u.amf_as.u.establish.new_guti = NULL;
            }
            //----------------------------------------
            //REQUIREMENT_3GPP_24_501(R15_5_5_1_2_4__14);
            //amf_sap.u.amf_as.u.establish.eps_network_feature_support = &_amf_data.conf.eps_network_feature_support;

            /*
            * Delete any preexisting UE radio capabilities, pursuant to
            * GPP 24.5R15:5.5.1.2.4
            */
            // Note: this is safe from double-free errors because it sets to NULL
            // after freeing, which free treats as a no-op.
            bdestroy_wrapper(&ue_m5gmm_context_p->ue_radio_capability);

            /*
            * Setup EPS NAS security data
            
            amf_as_set_security_data(&amf_sap.u.amf_as.u.establish.sctx, &amf_context->_security, false, true);
            amf_sap.u.amf_as.u.establish.encryption =amf_context->_security.selected_algorithms.encryption;
            amf_sap.u.amf_as.u.establish.integrity =amf_context->_security.selected_algorithms.integrity;
            OAILOG_DEBUG(LOG_NAS_AMF,"ue_id=" AMF_UE_NGAP_ID_FMT " AMF-PROC  - encryption = 0x%X (0x%X)\n",ue_id,amf_sap.u.amf_as.u.establish.encryption,amf_context->_security.selected_algorithms.encryption);
            OAILOG_DEBUG(LOG_NAS_AMF,"ue_id=" AMF_UE_NGAP_ID_FMT " AMF-PROC  - integrity  = 0x%X (0x%X)\n",ue_id,amf_sap.u.amf_as.u.establish.integrity,amf_context->_security.selected_algorithms.integrity);
            /*
            * Get the activate default 5GMM PDu Session context request message to
            * transfer within the SMF container of the Registration accept message
            */
            amf_sap.u.amf_as.u.establish.nas_msg = registration_proc->smf_msg_out;
            OAILOG_TRACE(LOG_NAS_AMF,"ue_id=" AMF_UE_NGAP_ID_FMT " AMF-PROC  - nas_msg  src size = %d nas_msg  dst size = %d \n", ue_id,blength(registration_proc->smf_msg_out),blength(amf_sap.u.amf_as.u.establish.nas_msg));

            // Send T3402
            //amf_sap.u.amf_as.u.establish.t3402 = &amf_config.nas_config.t3402_min;

            //Encode CSFB parameters
           // _encode_csfb_parameters_attach_accept(amf_context, &amf_sap.u.amf_as.u.establish);

            //REQUIREMENT_3GPP_24_501(R15_5_5_1_2_4__2);
            rc = amf_sap_send(&amf_sap);

            if (RETURNerror != rc) 
            {
                /*
                * Start T3450 timer
                */
               // nas_stop_T3450(registration_proc->ue_id, &registration_proc->T3450, NULL);
               // nas_start_T3450(registration_proc->ue_id,&registration_proc->T3450,registration_proc->amf_spec_proc.amf_proc.base_proc.time_out,(void *) amf_context);
            }
        } 
        else 
        {
            OAILOG_WARNING(LOG_NAS_AMF, "ue_mm_context NULL\n");
        }
        //increment_counter("ue_registration", 1, 1, "action", "registration_accept_sent");
        OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
    }
    /****************************************************************************
     **                                                                       **
    ** Name:    amf_send_registration_accept_dl_nas()                         **
    **                                                                        **
    ** Description: Builds Registration Accept message to be sent
    ** is NGAP : DL NAS Tx **
    **                                                                        **
    **      The registration Accept message is sent by the network to the     **
    **      UE to indicate that the corresponding attach request has          **
    **      been accepted.                                                    **
    **                                                                        **
    ** Inputs:  msg:       The AMFAS-SAP primitive to process                 **
    **      Others:    None                                                   **
    **                                                                        **
    ** Outputs:     amf_msg:   The AMF message to be sent                     **
    **      Return:    The size of the AMF message                            **
    **      Others:    None                                                   **
    **                                                                        **
    ***************************************************************************/
    int amf_send_registration_accept_dl_nas(const amf_as_data_t* msg, 
                                            registration_accept_msg *amf_msg) {
    OAILOG_FUNC_IN(LOG_NAS_AMF);
    int size = AMF_HEADER_MAXIMUM_LENGTH;

    // Get the UE context
    amf_context_t* amf_ctx;// = amf_context_get(&_amf_data, msg->ue_id);
    //DevAssert(amf_ctx);
    ue_m5gmm_context_s* ue_m5gmm_context_p =
        PARENT_STRUCT(amf_ctx, class ue_m5gmm_context_s, amf_context);
    amf_ue_ngap_id_t ue_id = ue_m5gmm_context_p->amf_ue_ngap_id;
    DevAssert(msg->ue_id == ue_id);

    OAILOG_DEBUG(LOG_NAS_AMF, "AMFAS-SAP - Send Regisration Accept message\n");
    OAILOG_DEBUG(LOG_NAS_AMF, "AMFAS-SAP - size = AMF_HEADER_MAXIMUM_LENGTH(%d)\n", size);
    /*
    * Mandatory - Message type
    */
    amf_msg->m5gsregistrationtype = REGISTRATION_ACCEPT;
    /*
    * Mandatory - 5GS Registration result
    */
    size += M5GS_REGISTRATION_RESULT_MAXIMUM_LENGTH;
    OAILOG_INFO(LOG_NAS_AMF,
        "AMFAS-SAP - size += AMF_REGISTRATION_RESULT_MAXIMUM_LENGTH(%d)  (%d)\n",
        M5GS_REGISTRATION_RESULT_MAXIMUM_LENGTH, size);
    switch (amf_ctx->m5gsregistrationtype) {
        case AMF_REGISTRATION_TYPE_INITIAL:
        amf_msg->m5ggsgsregistrationresult = M5GS_REGISTRATION_RESULT_3GPP_ACCESS;
        OAILOG_DEBUG(LOG_NAS_AMF, "AMFAS-SAP - M5GS_REGISTRATION_RESULT_3GPP_ACCESS\n");
        break;
        case AMF_REGISTRATION_TYPE_EMERGENCY:  // We should not reach here
        OAILOG_ERROR(LOG_NAS_AMF,"AMFAS-SAP - M5GS emergency Registration, currently unsupported\n");
        OAILOG_FUNC_RETURN(LOG_NAS_AMF, 0);  // TODO: fix once supported
        break;
    }
    
     /*
    * Optional - Mobile Identity
    */
    if (msg->m5ggsmobileidentity) {
        size += M5GS_MOBILE_IDENTITY_MAXIMUM_LENGTH;
        amf_msg->presencemask |= REGISTRATION_ACCEPT_UE_IDENTITY_PRESENT;
        if (msg->msidentity->imsi.typeofidentity == MOBILE_IDENTITY_IMSI) {
        memcpy(
            &amf_msg->msidentity.imsi, &msg->ms_identity->imsi,
            sizeof(amf_msg->msidentity.imsi));
        } else if (msg->ms_identity->imsi.typeofidentity == MOBILE_IDENTITY_TMSI) {
        memcpy(&amf_msg->msidentity.tmsi, &msg->ms_identity->tmsi,
        sizeof(amf_msg->msidentity.tmsi));
        }
    }

    /*
    * Optional - Additional Update Result
   
    if (msg->additional_update_result) {
        size += ADDITIONAL_UPDATE_RESULT_MAXIMUM_LENGTH;
        amf_msg->presencemask |= REGISTRATION_ACCEPT_ADDITIONAL_UPDATE_RESULT_PRESENT;
        amf_msg->additionalupdateresult = SMS_ONLY;
    }
     */
     OAILOG_FUNC_RETURN(LOG_NAS_AMF, size);
    }


}
