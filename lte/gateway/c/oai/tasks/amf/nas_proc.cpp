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

  Source      nas_proc.cpp

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#include "nas_proc.h"
#include "common_defs.h"
#include "as_message.h"
#include "amf_sap.h"
using namespace std;
#pragma once

namespace magma5g
{
    int nas_proc::nas_proc_establish_ind(const amf_ue_ngap_id_t ue_id, const bool is_mm_ctx_new,
                                        const tai_t originating_tai,const ecgi_t ecgi,
                                        const as_m5gcause_t as_cause,const s_tmsi_t s_tmsi,
                                        std::string *msg)
    {
        if (msg)
         {
            amf_sap_t amf_sap = {0};

            /*
            * Notify the AMF procedure call manager that NAS signalling
            * connection establishment indication message has been received
            * from the Access-Stratum sublayer
            */

            amf_sap.primitive = AMFAS_ESTABLISH_REQ;
            amf_sap.u.amf_as.u.establish.ue_id = ue_id;
            amf_sap.u.amf_as.u.establish.is_initial = true;
            amf_sap.u.amf_as.u.establish.is_mm_ctx_new = is_mm_ctx_new;

            amf_sap.u.amf_as.u.establish.nas_msg = *msg;
            *msg = NULL;
            amf_sap.u.amf_as.u.establish.tai = &originating_tai;
            //amf_sap.u.amf_as.u.establish.plmn_id            = &originating_tai.plmn;
            //amf_sap.u.amf_as.u.establish.tac                = originating_tai.tac;
            amf_sap.u.amf_as.u.establish.ecgi = ecgi;

            rc = amf_sap_c::amf_sap_send(&amf_sap);
         }
         OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);



    }
    static emm_procedures_t* _nas_new_amf_procedures(amf_context_t* const amf_context) 
    {
      amf_procedures_t* amf_procedures =new(*emm_context->emm_procedures);
      LIST_INIT(&amf_procedures->amf_common_procs);
      return amf_procedures;
}
    //-----------------------------------------------------------------------------
  nas_5g_auth_info_proc_t* nas_new_5gcn_auth_info_procedure(amf_context_t* const amf_context) {
    if (!(amf_context->amf_procedures)) {
      amf_context->amf_procedures = _nas_new_amf_procedures(amf_context);
    }

      nas_5g_auth_info_proc_t* auth_info_proc =new(nas_5g_auth_info_proc_t);
    auth_info_proc->cn_proc.base_proc.nas_puid =
        __sync_fetch_and_add(&nas_puid, 1);
    auth_info_proc->cn_proc.base_proc.type = NAS_PROC_TYPE_CN;
    auth_info_proc->cn_proc.type           = CN_PROC_AUTH_INFO;

    nas_cn_procedure_t* wrapper = new(*wrapper);
    if (wrapper) {
      wrapper->proc = &auth_info_proc->cn_proc;
      LIST_INSERT_HEAD(&amf_context->mf_procedures->cn_procs, wrapper, entries);
      OAILOG_TRACE(LOG_NAS_AMF, "New CN_PROC_AUTH_INFO\n");
      return auth_info_proc;
    } else {
      free_wrapper((void**) &auth_info_proc);
    }
    return NULL;
  }
//---------------------------------------------------------------------------------------------
    nas_amf_registration_proc_t* nas_proc::get_nas_specific_procedure_registration(const struct amf_context_s* const ctxt) 
    {
      if ((ctxt) && (ctxt->amf_procedures) && (ctxt->amf_procedures->amf_specific_proc) &&
         ((AMF_SPEC_PROC_TYPE_REGISTRATION == ctxt->amf_procedures->amf_specific_proc->type))){
           return (nas_amf_registration_proc_t*) ctxt->amf_procedures->amf_specific_proc;
         }
      
      return NULL;
    }
    //-----------------------------------------------------------------------------
static amf_procedures_t* nas_new_amf_procedures(amf_context_t* const amf_context)
 {
     amf_procedures_t* amf_procedures =new(*amf_context->amf_procedures);
  LIST_INIT(&amf_procedures->amf_common_procs);
  return amf_procedures;
}
    //-----------------------------------------------------------------------------
  nas_amf_ident_proc_t* nas5g_new_identification_procedure(amf_context_t* const amf_context) {
  if (!(amf_context->amf_procedures)) {
    amf_context->amf_procedures = nas_new_amf_procedures(amf_context);
  }

  nas_amf_ident_proc_t* ident_proc = new(nas_amf_ident_proc_t);

  ident_proc->amf_com_proc.amf_proc.base_proc.nas_puid = __sync_fetch_and_add(&nas_puid, 1);
  ident_proc->amf_com_proc.amf_proc.base_proc.type = NAS_PROC_TYPE_AMF;
  ident_proc->amf_com_proc.amf_proc.type           = NAS_AMF_PROC_TYPE_COMMON;
  ident_proc->amf_com_proc.type                    = AMF_COMM_PROC_IDENT;

  ident_proc->T3470.sec = amf_config.nas_config.t3470_sec;
  ident_proc->T3470.id  = NAS_TIMER_INACTIVE_ID;

  nas_amf_common_procedure_t* wrapper = new(wrapper);
  if (wrapper) {
    wrapper->proc = &ident_proc->amf_com_proc;
    LIST_INSERT_HEAD( &amf_context->amf_procedures->amf_common_procs, wrapper, entries);
    OAILOG_TRACE(LOG_NAS_AMF, "New AMF_COMM_PROC_IDENT\n");
    return ident_proc;
  } else {
    free_wrapper((void**) &ident_proc);
  }
  return ident_proc;
}
/*
   --------------------------------------------------------------------------
                AMF NAS specific local functions
   --------------------------------------------------------------------------
*/

/*
 * Description: Sends IDENTITY REQUEST message and start timer T3470.
 *
 * Inputs:  args:      handler parameters
 *      Others:    None
 *
 * Outputs:     None
 *      Return:    None
 *      Others:    T3470
 */
static int amf_identification_request(nas_amf_ident_proc_t* const proc) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  mf_sap_t amf_sap             = {0};
  int rc                        = RETURNok;
  amf_context_t* amf_ctx = NULL;

  ue_m5gmm_context_s* ue_5gmm_context =  amf_ue_context_exists_amf_ue_ngap_id(proc->ue_id);
  if (ue_5gmm_context) {
    amf_ctx = &ue_5gmm_context->amf_context;
  } else {
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
  }
  /*
   * Notify AMF-AS SAP that Identity Request message has to be sent
   * to the UE
   */
  amf_sap.primitive = EMMAS_SECURITY_REQ;
  amf_sap.u.amf_as.u.security.puid        =  proc->amf_com_proc.amf_proc.base_proc.nas_puid;
  amf_sap.u.amf_as.u.security.guti_m5_t   = NULL;
  amf_sap.u.amf_as.u.security.ue_id       = proc->ue_id;
  amf_sap.u.amf_as.u.security.msg_type    = AMF_AS_MSG_TYPE_IDENT;
  amf_sap.u.amf_as.u.security.ident_type  = proc->identity_type;

  /*
   * Setup 5G CN NAS security data
   */
  //TODO
  //amf_as_set_security_data(&amf_sap.u.amf_as.u.security.sctx, &amf_ctx->_security, false, true);
  rc = amf_sap_send(&amf_sap);

  if (rc != RETURNerror) {
    /*
     * Start T3470 timer
     */
    //TODO
    //nas_start_T3470(proc->ue_id, &proc->T3470,proc->amf_com_proc.amf_proc.base_proc.time_out, (void*) amf_ctx);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}
//-------------------------------------------------------------------------------------
      int amf_proc_identification(amf_context_t* const amf_context, nas_amf_proc_t* const amf_proc,
        const identity_type2_t type, success_cb_t success, failure_cb_t failure) 
        {
          OAILOG_FUNC_IN(LOG_NAS_AMF);
          int rc = RETURNerror;

          if ((amf_context) && ((AMF_DEREGISTERED == amf_context->_anf_fsm_state) ||
                                (AMF_REGISTERED == amf_context->_anf_fsm_state))) 
            {
            
              amf_ue_ngap_id_t ue_id = PARENT_STRUCT(amf_context, ue_m5gmm_context_s, amf_context)->amf_ue_ngap_id;

             OAILOG_INFO(LOG_NAS_AMF,"AMF-PROC  - Initiate identification type = %s (%d), ctx = %p\n",
                amf_identity_type_str[type], type, amf_context);

            nas_amf_ident_proc_t* ident_proc = nas5g_new_identification_procedure(amf_context);
            if (ident_proc) {
              if (amf_proc) {
                if ((NAS_AMF_PROC_TYPE_SPECIFIC == amf_proc->type) &&
                    (AMF_SPEC_PROC_TYPE_REGISTRATION == ((nas_amf_specific_proc_t*) amf_proc)->type)) {
                  ident_proc->is_cause_is_registered = true;
                }
              }
              ident_proc->identity_type                   = type;
              ident_proc->retransmission_count            = 0;
              ident_proc->ue_id                           = ue_id;
              ((nas_base_proc_t*) ident_proc)->parent     = (nas_base_proc_t*) amf_proc;
              ident_proc->amf_com_proc.amf_proc.delivered = NULL;
              //TODO ident_proc->amf_com_proc.amf_proc.previous_amf_fsm_state =  amf_fsm_get_state(amf_context);
              //ident_proc->amf_com_proc.amf_proc.not_delivered =identification_Nll_failure;
             // ident_proc->amf_com_proc.amf_proc.not_delivered_ho = _identification_non_delivered_ho;
              ident_proc->amf_com_proc.amf_proc.base_proc.success_notif = success;
              ident_proc->amf_com_proc.amf_proc.base_proc.failure_notif = failure;
              ident_proc->amf_com_proc.amf_proc.base_proc.abort = _identification_abort;
              ident_proc->amf_com_proc.amf_proc.base_proc.fail_in = NULL;  // only response
              //ident_proc->amf_com_proc.amf_proc.base_proc.time_out = _identification_t3470_handler;
            }

            rc = amf_identification_request(ident_proc);

            if (rc != RETURNerror) {
              /*
              * Notify 5G CN that common procedure has been initiated
              */
              amf_sap_t amf_sap = {0};

              amf_sap.primitive       = EMMREG_COMMON_PROC_REQ;
              amf_sap.u.amf_reg.ue_id = ue_id;
              amf_sap.u.amf_reg.ctx   = amf_context;
              rc                      = amf_sap_send(&amf_sap);
            }
          }

          OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
        }
}