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

  Source      amf_authentication.cpp

  Version     0.1

  Date        2020/10/01

  Product     AMF 

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#include "log.h"
#include "nas_proc.h"
#include "amf_app_ue_context.h"

using namespace std;


namespace magma5g
{
    /****************************************************************************
     **                                                                        **
    ** Name:    nas_itti_auth_info_req()                                     **
    **                                                                        **
    ** Description: Sends Authenticatio Req to UDM via S6a Task               **
    **                                                                        **
    ** Inputs: ue_idP: UE context Identifier                                  **
    **      imsiP: IMSI of UE                                                 **
    **      is_initial_reqP: Flag to indicate, whether Auth Req is sent       **
    **                      for first time or initited as part of             **
    **                      re-synchronisation                                **
    **      visited_plmnP : Visited PLMN                                      **
    **      num_vectorsP : Number of Auth vectors in case of                  **
    **                    re-synchronisation                                  **
    **      auts_pP : sent in case of re-synchronisation                      **
    ** Outputs:                                                               **
    **     Return: None                                                       **
    **                                                                        **
    ***************************************************************************/
    static void nas_itti_auth_info_req(
        const amf_ue_ngap_id_t ue_id, const imsi_t* const imsiP,
        const bool is_initial_reqP, plmn_t* const visited_plmnP,
        const uint8_t num_vectorsP, const_bstring const auts_pP) {
    OAILOG_FUNC_IN(LOG_NAS);
    MessageDef* message_p              = NULL;
    s6a_auth_info_req_t* auth_info_req = NULL;

    OAILOG_INFO( LOG_NAS_AMF, "Sending Authentication Information Request message to S6A"
                " for ue_id =" AMF_UE_NGAP_ID_FMT "\n", ue_id);

    message_p = itti_alloc_new_message(TASK_AMF_APP, S6A_AUTH_INFO_REQ);
    if (!message_p) {
        OAILOG_CRITICAL(LOG_NAS_AMF,"itti_alloc_new_message failed for Authentication"
            " Information Request message to S6A for"" ue-id = " AMF_UE_NGAP_ID_FMT "\n", ue_id);
        OAILOG_FUNC_OUT(LOG_NAS);
    }
    auth_info_req = &message_p->ittiMsg.s6a_auth_info_req;
    memset(auth_info_req, 0, sizeof(s6a_auth_info_req_t));

    IMSI_TO_STRING(imsiP, auth_info_req->imsi, IMSI_BCD_DIGITS_MAX + 1);
    auth_info_req->imsi_length = (uint8_t) strlen(auth_info_req->imsi);

    if (!(auth_info_req->imsi_length > 5) && (auth_info_req->imsi_length < 16)) {
        OAILOG_WARNING(
            LOG_NAS_AMF, "Bad IMSI length %d", auth_info_req->imsi_length);
        OAILOG_FUNC_OUT(LOG_NAS);
    }
    auth_info_req->visited_plmn  = *visited_plmnP;
    auth_info_req->nb_of_vectors = num_vectorsP;

    if (is_initial_reqP) {
        auth_info_req->re_synchronization = 0;
        memset(auth_info_req->resync_param, 0, sizeof auth_info_req->resync_param);
    } else {
        if (!auts_pP) {
        OAILOG_WARNING(LOG_NAS_EMM, "Auts Null during resynchronization \n");
        OAILOG_FUNC_OUT(LOG_NAS);
        }
        auth_info_req->re_synchronization = 1;
        memcpy(
            auth_info_req->resync_param, auts_pP->data,
            sizeof auth_info_req->resync_param);
    }
    send_msg_to_task(&amf_app_task_zmq_ctx, TASK_S6A, message_p);
    OAILOG_FUNC_OUT(LOG_NAS);
    }
    //------------------------------------------------------------------------------
static int start_authentication_information_procedure(amf_context_t* amf_context, 
                                nas_amf_auth_proc_t* const auth_proc,  const_bstring auts) 
    {
        OAILOG_FUNC_IN(LOG_NAS_AMF);
        amf_ue_ngap_id_t ue_id =  PARENT_STRUCT(amf_context, ue_m5gmm_context_s, amf_context)->amf_ue_ngap_id;
        // Ask upper layer to fetch new security context
        nas_5g_auth_info_proc_t* auth_info_proc = get_nas_cn_procedure_auth_info(amf_context);
        if (!auth_info_proc) {
            auth_info_proc               = nas_new_cn_auth_info_procedure(amf_context);
            auth_info_proc->request_sent = false;
        }

        auth_info_proc->cn_proc.base_proc.parent         =   &auth_proc->amf_com_proc.amf_proc.base_proc;
        auth_proc->amf_com_proc.amf_proc.base_proc.child =   &auth_info_proc->cn_proc.base_proc;
        auth_info_proc->success_notif = _auth_info_proc_success_cb;
        auth_info_proc->failure_notif = _auth_info_proc_failure_cb;
        auth_info_proc->ue_id  = ue_id;
        auth_info_proc->resync = auth_info_proc->request_sent;

        plmn_t visited_plmn     = {0};
        visited_plmn.mcc_digit1 = amf_context->originating_tai.mcc_digit1;
        visited_plmn.mcc_digit2 = amf_context->originating_tai.mcc_digit2;
        visited_plmn.mcc_digit3 = amf_context->originating_tai.mcc_digit3;
        visited_plmn.mnc_digit1 = amf_context->originating_tai.mnc_digit1;
        visited_plmn.mnc_digit2 = amf_context->originating_tai.mnc_digit2;
        visited_plmn.mnc_digit3 = amf_context->originating_tai.mnc_digit3;

        bool is_initial_req          = !(auth_info_proc->request_sent);
        auth_info_proc->request_sent = true;
        

        nas_itti_auth_info_req(
            ue_id, &amf_context->_imsi, is_initial_req, &visited_plmn,
            MAX_EPS_AUTH_VECTORS, auts);

        OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
}
    //------------------------------------------------------------------------------
    int authentication::amf_proc_authentication(amf_context_s* amf_context, nas_amf_specific_proc_t* const amf_specific_proc,
                                 success_cb_t success, failure_cb_t failure) 
    {
        OAILOG_FUNC_IN(LOG_NAS_EMM);
        int rc = RETURNerror;

        amf_ue_ngap_id_t ue_id = PARENT_STRUCT(amf_context, ue_m5gmm_context_s, amf_context)->amf_ue_ngap_id;
        nas_amf_auth_proc_t* auth_proc ;//= get_nas_common_procedure_authentication(amf_context);
        if (!auth_proc) {
            //auth_proc = nas_new_authentication_procedure(amf_context);
        }
        if (!auth_proc) {
            if (amf_specific_proc) {
            if (AMF_SPEC_PROC_TYPE_REGISTRATION == amf_specific_proc->type) {
                auth_proc->is_cause_is_registered = true;
            } else if (AMF_SPEC_PROC_TYPE_TAU == amf_specific_proc->type) {
                auth_proc->is_cause_is_registered = false;
            }
            }

            auth_proc->amf_cause            = AMF_CAUSE_SUCCESS;
            auth_proc->retransmission_count = 0;
            auth_proc->ue_id                = ue_id;
            ((nas_base_proc_t*) auth_proc)->parent = (nas_base_proc_t*) amf_specific_proc;
            auth_proc->amf_com_proc.amf_proc.delivered = NULL;
            //TODO
            //auth_proc->amf_com_proc.amf_proc.previous_amf_fsm_state = amf_fsm_get_state(amf_context);
            auth_proc->amf_com_proc.amf_proc.not_delivered           = NULL;
            auth_proc->amf_com_proc.amf_proc.not_delivered_ho        = NULL;
            auth_proc->amf_com_proc.amf_proc.base_proc.success_notif = success;
            auth_proc->amf_com_proc.amf_proc.base_proc.failure_notif = failure;
            auth_proc->amf_com_proc.amf_proc.base_proc.abort   = _authentication_abort;
            auth_proc->amf_com_proc.amf_proc.base_proc.fail_in = NULL;  // only response
            //TODO
            //auth_proc->amf_com_proc.amf_proc.base_proc.fail_out = _authentication_reject;
            auth_proc->amf_com_proc.amf_proc.base_proc.time_out = NULL;

            bool run_auth_info_proc = false;
            if (!IS_AMF_CTXT_VALID_AUTH_VECTORS(amf_context)) {
            // Ask upper layer to fetch new security context
            nas_5g_auth_info_proc_t* auth_info_proc ;//=  get_nas_cn_procedure_auth_info(amf_context);
            if (!auth_info_proc) {
                auth_info_proc = nas_new_5gcn_auth_info_procedure(amf_context);
            }
            if (!auth_info_proc->request_sent) {
                run_auth_info_proc = true;
            }
            rc = RETURNok;
            } else {
            ksi_t eksi = 0;
            if (amf_context->_security.eksi < KSI_NO_KEY_AVAILABLE) {
                eksi = (amf_context->_security.eksi + 1) % (EKSI_MAX_VALUE + 1);
            }
            for (; eksi < MAX_5G_AUTH_VECTORS; eksi++) {
                if (IS_AMF_CTXT_VALID_AUTH_VECTOR(
                        amf_context, (eksi % MAX_5G_AUTH_VECTORS))) {
                break;
                }
            }
            // eksi should always be 0
            if (!IS_AMF_CTXT_VALID_AUTH_VECTORS( amf_context, (eksi % MAX_5G_AUTH_VECTORS))) {
                run_auth_info_proc = true;
            } else {
                rc = authentication::amf_proc_authentication_ksi(
                    amf_context, amf_specific_proc, eksi,
                    amf_context->_vector[eksi % MAX_5G_AUTH_VECTORS].rand,
                    amf_context->_vector[eksi % MAX_5G_AUTH_VECTORS].autn, success,
                    failure);
                OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
            }
            }
            if (run_auth_info_proc) {
            rc = start_authentication_information_procedure( amf_context, auth_proc, NULL);
            }
        }

        OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
    }
    /*
        --------------------------------------------------------------------------
                Authentication procedure executed by the AMF
        --------------------------------------------------------------------------
        */
        /****************************************************************************
         **                                                                        **
        ** Name:    amf_proc_authentication()                                 **
        **                                                                        **
        ** Description: Initiates authentication procedure to establish partial   **
        **      native 5G CN security context in the UE and the AMF.        **
        **                                                                        **
        **              3GPP TS 24.501, section 5.4.1.3                           **
        **      The network initiates the authentication procedure by     **
        **      sending an AUTHENTICATION REQUEST message to the UE and   **
        **      starting the timer T3460. The AUTHENTICATION REQUEST mes- **
        **      sage contains the parameters necessary to calculate the   **
        **      authentication response.                                  **
        **                                                                        **
        ** Inputs:  ue_id:      UE lower layer identifier                  **
        **      ksi:       NAS key set identifier                     **
        **      rand:      Random challenge number                    **
        **      autn:      Authentication token                       **
        **      success:   Callback function executed when the authen-**
        **             tication procedure successfully completes  **
        **      reject:    Callback function executed when the authen-**
        **             tication procedure fails or is rejected    **
        **      failure:   Callback function executed whener a lower  **
        **             layer failure occured before the authenti- **
        **             cation procedure comnpletes                **
        **      Others:    None                                       **
        **                                                                        **
        ** Outputs:     None                                                      **
        **      Return:    RETURNok, RETURNerror                      **
        **      Others:    None                                       **
        **                                                                        **
        ***************************************************************************/
        int authentication::amf_proc_authentication_ksi(amf_context_t* amf_context, nas_amf_specific_proc_t* const amf_specific_proc, 
                                                    ksi_t ksi,const uint8_t* const rand, const uint8_t* const autn, 
                                                    success_cb_t success, failure_cb_t failure) 
        {
            OAILOG_FUNC_IN(LOG_NAS_EMM);
            int rc = RETURNerror;

            if ((emm_context) && ((AMF_DEREGISTERED == amf_context->amf_fsm_state) ||
                                    (AMF_REGISTERED == amf_context->amf_fsm_state))) 
                {
                    amf_ue_ngap_id_t ue_id = PARENT_STRUCT(amf_context, ue_m5gmm_context_s, amf_context)->amf_ue_ngap_id;
                    OAILOG_INFO(LOG_NAS_AMF,"ue_id=" AMF_UE_NGAP_ID_FMT " AMF-PROC  - Initiate Authentication KSI = %d\n", ue_id, ksi);
                    nas_amf_auth_proc_t* auth_proc = get_nas_common_procedure_authentication(amf_context);
                    if (!auth_proc) {
                    //auth_proc = nas_new_authentication_procedure(emm_context);
                }
                if (auth_proc) {
                
                    if (AMF_SPEC_PROC_TYPE_REGISTRATION == amf_specific_proc->type) 
                    auth_proc->is_cause_is_registered = true;
                    OAILOG_DEBUG(LOG_NAS_AMF,"Auth proc cause is AMF_SPEC_PROC_TYPE_REGISTRATION (%d) for ue_id " "(%u)\n",
                        amf_specific_proc->type, ue_id);
                    
                }
                // Set the RAND value
                auth_proc->ksi = ksi;
                if (rand) {
                    memcpy(auth_proc->rand, rand, AUTH_RAND_SIZE);
                }
                // Set the authentication token
                if (autn) {
                    memcpy(auth_proc->autn, autn, AUTH_AUTN_SIZE);
                }
                auth_proc->amf_cause            = AMF_CAUSE_SUCCESS;
                auth_proc->retransmission_count = 0;
                auth_proc->ue_id                = ue_id;
                ((nas_base_proc_t*) auth_proc)->parent = (nas_base_proc_t*) amf_specific_proc;
                auth_proc->amf_com_proc.amf_proc.delivered = NULL;
                auth_proc->amf_com_proc.amf_proc.previous_amf_fsm_state = amf_fsm_get_state(amf_context);
                auth_proc->amf_com_proc.amf_proc.not_delivered =  _authentication_ll_failure;
                auth_proc->amf_com_proc.amf_proc.not_delivered_ho = _authentication_non_delivered_ho;
                auth_proc->amf_com_proc.amf_proc.base_proc.success_notif = success;
                auth_proc->amf_com_proc.amf_proc.base_proc.failure_notif = failure;
                auth_proc->amf_com_proc.amf_proc.base_proc.abort = _authentication_abort;
                auth_proc->amf_com_proc.amf_proc.base_proc.fail_in = NULL;  // only response
                auth_proc->amf_com_proc.amf_proc.base_proc.fail_out = _authentication_reject;
                auth_proc->amf_com_proc.amf_proc.base_proc.time_out = _authentication_t3460_handler;
                }

                /*
                * Send authentication request message to the UE
                */
                rc = _authentication_request(amf_context, auth_proc);

                if (rc != RETURNerror) {
                /*
                * Notify EMM that common procedure has been initiated
                */
                amf_sap_t amf_sap = {0};

                amf_sap.primitive       = AMFREG_COMMON_PROC_REQ;
                amf_sap.u.amf_reg.ue_id = ue_id;
                amf_sap.u.amf_reg.ctx   = amf_context;
                rc                      = amf_sap_send(&amf_sap);
                }
            }
            OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
        }
}