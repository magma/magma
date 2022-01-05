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
#include "include/amf_client_servicer.h"
#include <sstream>
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/lib/secu/secu_defs.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.003.h"
#include "lte/gateway/c/core/oai/common/common_defs.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_authentication.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_recv.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_identity.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_sap.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_timer_management.h"

#define AMF_CAUSE_SUCCESS (1)
#define MAX_5G_AUTH_VECTORS 1

using magma5g::AMFClientServicer;

namespace magma5g {
extern task_zmq_ctx_t amf_app_task_zmq_ctx;

amf_as_data_t amf_data_sec_auth;
static int authenthication_t3560_handler(
    zloop_t* loop, int timer_id, void* output);

nas_amf_smc_proc_t* get_nas5g_common_procedure_smc(const amf_context_t* ctxt) {
  return (nas_amf_smc_proc_t*) get_nas5g_common_procedure(
      ctxt, AMF_COMM_PROC_SMC);
}

nas5g_cn_proc_t* get_nas5g_cn_procedure(
    const amf_context_t* ctxt, cn5g_proc_type_t proc_type) {
  if (ctxt) {
    if (ctxt->amf_procedures) {
      nas5g_cn_procedure_t* p1 = LIST_FIRST(&ctxt->amf_procedures->cn_procs);
      nas5g_cn_procedure_t* p2 = NULL;
      while (p1) {
        p2 = LIST_NEXT(p1, entries);
        if (p1->proc->type == proc_type) {
          return p1->proc;
        }
        p1 = p2;
      }
    }
  }
  return NULL;
}

static int calculate_amf_serving_network_name(
    amf_context_t* amf_ctx, uint8_t* snni);
/***************************************************************************
**                                                                        **
** Name:    get_nas5g_cn_procedure_auth_info()                            **
**                                                                        **
** Description: Invokes get_nas5g_cn_procedure                            **
**              to fetch new security context                             **
**                                                                        **
**                                                                        **
***************************************************************************/
nas5g_auth_info_proc_t* get_nas5g_cn_procedure_auth_info(
    const amf_context_t* ctxt) {
  return (nas5g_auth_info_proc_t*) get_nas5g_cn_procedure(
      ctxt, CN5G_PROC_AUTH_INFO);
}

/***************************************************************************
**                                                                        **
** Name:    start_authentication_information_procedure()                  **
**                                                                        **
** Description: Invokes get_nas5g_cn_proceduree_auth_info                 **
**              to fetch new security context                             **
**                                                                       **
**                                                                        **
***************************************************************************/
static int start_authentication_information_procedure(
    amf_context_t* amf_context, nas5g_amf_auth_proc_t* const auth_proc,
    const_bstring auts) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);

  int rc = RETURNerror;
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  uint8_t snni[40] = {0};

  amf_ue_ngap_id_t ue_id =
      PARENT_STRUCT(amf_context, ue_m5gmm_context_s, amf_context)
          ->amf_ue_ngap_id;
  // Upper layer to fetch new security context
  nas5g_auth_info_proc_t* auth_info_proc =
      get_nas5g_cn_procedure_auth_info(amf_context);
  if (!auth_info_proc) {
    auth_info_proc = nas5g_new_cn_auth_info_procedure(amf_context);
    auth_info_proc->request_sent = false;
  }

  auth_info_proc->cn_proc.base_proc.parent =
      &auth_proc->amf_com_proc.amf_proc.base_proc;
  auth_proc->amf_com_proc.amf_proc.base_proc.child =
      &auth_info_proc->cn_proc.base_proc;
  auth_info_proc->ue_id  = ue_id;
  auth_info_proc->resync = auth_info_proc->request_sent;

  bool is_initial_req          = !(auth_info_proc->request_sent);
  auth_info_proc->request_sent = true;

  IMSI64_TO_STRING(amf_context->imsi64, imsi_str, IMSI_LENGTH);

  rc = calculate_amf_serving_network_name(amf_context, snni);
  if (rc != RETURNok) {
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
  }

  if (is_initial_req) {
    OAILOG_INFO(
        LOG_AMF_APP,
        "Sending msg(grpc) to :[subscriberdb] for ue: [%s] auth-info\n",
        imsi_str);
    AMFClientServicer::getInstance().get_subs_auth_info(
        imsi_str, IMSI_LENGTH, reinterpret_cast<const char*>(snni), ue_id);
  } else if (auts->data) {
    OAILOG_INFO(
        LOG_AMF_APP,
        "Sending msg(grpc) to :[subscriberdb] for ue: [%s] auth-info-resync\n",
        imsi_str);
    AMFClientServicer::getInstance().get_subs_auth_info_resync(
        imsi_str, IMSI_LENGTH, reinterpret_cast<const char*>(snni), auts->data,
        RAND_LENGTH_OCTETS + AUTS_LENGTH, ue_id);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
}

static int start_authentication_information_procedure_synch(
    amf_context_t* amf_context, nas5g_amf_auth_proc_t* const auth_proc,
    const_bstring auts) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);

  // Ask upper layer to fetch new security context
  nas5g_auth_info_proc_t* auth_info_proc =
      get_nas5g_cn_procedure_auth_info(amf_context);

  if (!auth_info_proc) {
    auth_info_proc = nas5g_new_cn_auth_info_procedure(amf_context);
    auth_info_proc->request_sent = true;
    start_authentication_information_procedure(amf_context, auth_proc, auts);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
}

/***************************************************************************
**                                                                        **
** Name:    get_nas5g_common_procedure()                                  **
**                                                                        **
** Description:  Generic function to fetch security context and others    **
**                                                                        **
**                                                                        **
***************************************************************************/
nas_amf_common_proc_t* get_nas5g_common_procedure(
    const amf_context_t* const ctxt, amf_common_proc_type_t proc_type) {
  if (ctxt) {
    if (ctxt->amf_procedures) {
      nas_amf_common_procedure_t* p1 =
          LIST_FIRST(&ctxt->amf_procedures->amf_common_procs);
      nas_amf_common_procedure_t* p2 = NULL;
      while (p1) {
        p2 = LIST_NEXT(p1, entries);
        if (p1->proc->type == proc_type) {
          return p1->proc;
        }
        p1 = p2;
      }
    }
  }
  return NULL;
}

/***************************************************************************
**                                                                        **
** Name:    get_nas5g_common_procedure_authentication() **
**                                                                        **
** Description:  Generic function to fetch security context and others    **
**                                                                        **
**                                                                        **
***************************************************************************/
nas5g_amf_auth_proc_t* get_nas5g_common_procedure_authentication(
    const amf_context_t* const ctxt) {
  return (nas5g_amf_auth_proc_t*) get_nas5g_common_procedure(
      ctxt, AMF_COMM_PROC_AUTH);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_authentication_abort()                                    **
 **                                                                        **
 ** Description: Aborts the authentication procedure currently in progress **
 **                                                                        **
 ** Inputs:  args:      Authentication data to be released                 **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     None                                                      **
 **     Return: None                                                       **
 **     Others: None                                                       **
 **                                                                        **
 ***************************************************************************/
static int amf_authentication_abort(
    amf_context_t* amf_ctx, struct nas5g_base_proc_t* base_proc) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;
  if ((base_proc) && (amf_ctx)) {
    ue_m5gmm_context_s* ue_mm_context =
        PARENT_STRUCT(amf_ctx, ue_m5gmm_context_s, amf_context);
    OAILOG_DEBUG(
        LOG_NAS_AMF,
        "AMF-PROC  - Abort authentication procedure invoked "
        "(ue_id= " AMF_UE_NGAP_ID_FMT ")\n",
        ue_mm_context->amf_ue_ngap_id);

    rc = RETURNok;
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/***************************************************************************
**                                                                        **
** Name:    nas5g_new_authentication_procedure()                          **
**                                                                        **
** Description:  Handler for nas5g Authenthication Procedure              **
**                                                                        **
**                                                                        **
***************************************************************************/
nas5g_amf_auth_proc_t* nas5g_new_authentication_procedure(
    amf_context_t* const amf_context) {
  if (!(amf_context->amf_procedures)) {
    amf_context->amf_procedures = nas_new_amf_procedures(amf_context);
  }
  nas5g_amf_auth_proc_t* auth_proc = new (nas5g_amf_auth_proc_t)();
  auth_proc->amf_com_proc.amf_proc.base_proc.type = NAS_PROC_TYPE_AMF;
  auth_proc->amf_com_proc.amf_proc.type           = NAS_AMF_PROC_TYPE_COMMON;
  auth_proc->amf_com_proc.type                    = AMF_COMM_PROC_AUTH;
  auth_proc->retry_sync_failure                   = 0;
  nas_amf_common_procedure_t* wrapper = new nas_amf_common_procedure_t();
  if (wrapper) {
    wrapper->proc = &auth_proc->amf_com_proc;
    LIST_INSERT_HEAD(
        &amf_context->amf_procedures->amf_common_procs, wrapper, entries);
    OAILOG_TRACE(LOG_NAS_AMF, "New AMF_COMM_PROC_AUTH\n");
    return auth_proc;
  } else {
    free_wrapper((void**) &auth_proc);
  }
  return NULL;
}

/***************************************************************************
**                                                                        **
** Name:    amf_proc_authentication                                       **
**                                                                        **
** Description:  Procedure to start Authentication procedure              **
**                                                                        **
**                                                                        **
***************************************************************************/
int amf_proc_authentication(
    amf_context_t* amf_context,
    nas_amf_specific_proc_t* const amf_specific_proc, success_cb_t success,
    failure_cb_t failure) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc                  = RETURNerror;
  bool run_auth_info_proc = false;
  ksi_t eksi              = 0;
  OAILOG_DEBUG(LOG_NGAP, "starting Authentication procedure");
  amf_ue_ngap_id_t ue_id =
      PARENT_STRUCT(amf_context, ue_m5gmm_context_s, amf_context)
          ->amf_ue_ngap_id;
  nas5g_amf_auth_proc_t* auth_proc =
      get_nas5g_common_procedure_authentication(amf_context);
  if (!auth_proc) {
    auth_proc = nas5g_new_authentication_procedure(amf_context);
  }
  if (auth_proc) {
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
    ((nas5g_base_proc_t*) auth_proc)->parent =
        (nas5g_base_proc_t*) amf_specific_proc;
    auth_proc->amf_com_proc.amf_proc.delivered               = NULL;
    auth_proc->amf_com_proc.amf_proc.not_delivered           = NULL;
    auth_proc->amf_com_proc.amf_proc.not_delivered_ho        = NULL;
    auth_proc->amf_com_proc.amf_proc.base_proc.success_notif = success;
    auth_proc->amf_com_proc.amf_proc.base_proc.failure_notif = failure;
    auth_proc->amf_com_proc.amf_proc.base_proc.abort = amf_authentication_abort;
    auth_proc->amf_com_proc.amf_proc.base_proc.fail_in = NULL;  // only response
    // TODO Negative Scenarios to be taken in future.
    auth_proc->amf_com_proc.amf_proc.base_proc.time_out = NULL;
    if (!IS_AMF_CTXT_VALID_AUTH_VECTORS(amf_context)) {
      // Upper layer to fetch new security context
      nas5g_auth_info_proc_t* auth_info_proc =
          get_nas5g_cn_procedure_auth_info(amf_context);
      if (!auth_info_proc) {
        auth_info_proc = nas5g_new_cn_auth_info_procedure(amf_context);
      }
      if (!auth_info_proc->request_sent) {
        run_auth_info_proc = true;
      }
      rc = RETURNok;
    } else {
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
      if (!IS_AMF_CTXT_VALID_AUTH_VECTOR(
              amf_context, (eksi % MAX_5G_AUTH_VECTORS))) {
        run_auth_info_proc = true;
      } else {
        rc = amf_proc_authentication_ksi(
            amf_context, amf_specific_proc, eksi,
            amf_context->_vector[eksi % MAX_5G_AUTH_VECTORS].rand,
            amf_context->_vector[eksi % MAX_5G_AUTH_VECTORS].autn, success,
            failure);
        OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
      }
    }
    if (run_auth_info_proc) {
      rc = start_authentication_information_procedure(
          amf_context, auth_proc, NULL);
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_proc_authentication_ksi()                                 **
 **                                                                        **
 ** Description: Initiates authentication procedure to establish partial   **
 **      native 5G CN security context in the UE and the AMF.              **
 **                                                                        **
 **              3GPP TS 24.501, section 5.4.1.3                           **
 **      The network initiates the authentication procedure by             **
 **      sending an AUTHENTICATION REQUEST message to the UE and           **
 **      starting the timer T3560. The AUTHENTICATION REQUEST mes-         **
 **      sage contains the parameters necessary to calculate the           **
 **      authentication response.                                          **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                         **
 **      ksi:       NAS key set identifier                                 **
 **      rand:      Random challenge number                                **
 **      autn:      Authentication token                                   **
 **      success:   Callback function executed when the authen-            **
 **             tication procedure successfully completes                  **
 **      reject:    Callback function executed when the authen-            **
 **             tication procedure fails or is rejected                    **
 **      failure:   Callback function executed whener a lower              **
 **             layer failure occurred before the authenti-                 **
 **             cation procedure comnpletes                                **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                                  **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/
int amf_proc_authentication_ksi(
    amf_context_t* amf_context,
    nas_amf_specific_proc_t* const amf_specific_proc, ksi_t ksi,
    const uint8_t* const rand, const uint8_t* const autn, success_cb_t success,
    failure_cb_t failure) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;
  nas5g_amf_auth_proc_t* auth_proc;
  amf_ue_ngap_id_t ue_id;
  auth_proc = get_nas5g_common_procedure_authentication(amf_context);
  if (!auth_proc) {
    auth_proc = nas5g_new_authentication_procedure(amf_context);
  }

  if ((amf_context) && ((AMF_DEREGISTERED == amf_context->amf_fsm_state) ||
                        (AMF_REGISTERED == amf_context->amf_fsm_state))) {
    ue_id = PARENT_STRUCT(amf_context, ue_m5gmm_context_s, amf_context)
                ->amf_ue_ngap_id;
    OAILOG_DEBUG(
        LOG_NAS_AMF,
        "ue_id= " AMF_UE_NGAP_ID_FMT
        " AMF-PROC  - Initiate Authentication KSI = %d\n",
        ue_id, ksi);
    if (auth_proc) {
      if (AMF_SPEC_PROC_TYPE_REGISTRATION == amf_specific_proc->type)
        auth_proc->is_cause_is_registered = true;
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
    ((nas5g_base_proc_t*) auth_proc)->parent =
        (nas5g_base_proc_t*) amf_specific_proc;
    auth_proc->amf_com_proc.amf_proc.delivered               = NULL;
    auth_proc->amf_com_proc.amf_proc.base_proc.success_notif = success;
    auth_proc->amf_com_proc.amf_proc.base_proc.failure_notif = failure;
    auth_proc->amf_com_proc.amf_proc.base_proc.abort = amf_authentication_abort;
    auth_proc->amf_com_proc.amf_proc.base_proc.fail_in = NULL;
  }

  /*
   * Send authentication request message to the UE
   */
  rc = amf_send_authentication_request(amf_context, auth_proc);

  if (rc != RETURNerror) {
    /*
     * Notify AMF that common procedure has been initiated
     */
    amf_sap_t amf_sap       = {};
    amf_sap.primitive       = AMFREG_COMMON_PROC_REQ;
    amf_sap.u.amf_reg.ue_id = ue_id;
    amf_sap.u.amf_reg.ctx   = amf_context;
    rc                      = amf_sap_send(&amf_sap);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_proc_authentication_complete()                            **
 **                                                                        **
 ** Description: Performs the authentication completion procedure executed **
 **      by the network.                                                   **
 **                                                                        **
 **              3GPP TS 24.501, section 5.4.1.3.4                         **
 **      Upon receiving the AUTHENTICATION RESPONSE message, the           **
 **      MME shall stop timer T3560 and check the correctness of           **
 **      the RES parameter.                                                **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                         **
 **      emm_cause: Authentication failure AMF cause code                  **
 **      res:       Authentication response parameter. or auts             **
 **                 in case of sync failure                                **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                                  **
 **      Others:    amf_data, T3560                                        **
 **                                                                        **
 ***************************************************************************/
int amf_proc_authentication_complete(
    amf_ue_ngap_id_t ue_id, AuthenticationResponseMsg* msg, int amf_cause,
    const unsigned char* res) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc                         = RETURNerror;
  bool is_xres_validation_failed = false;
  nas_amf_smc_proc_t nas_amf_smc_proc_autn;
  nas_amf_registration_proc_t* registration_proc = NULL;
  nas5g_amf_auth_proc_t* auth_proc               = NULL;

  OAILOG_DEBUG(
      LOG_NAS_AMF,
      "Authentication  procedures complete for "
      "(ue_id=" AMF_UE_NGAP_ID_FMT ")\n",
      ue_id);
  ue_m5gmm_context_s* ue_mm_context = NULL;

  amf_context_t* amf_ctx = NULL;
  ue_mm_context          = amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  if (!ue_mm_context) {
    OAILOG_WARNING(
        LOG_NAS_AMF,
        "AMF-PROC - Failed to authenticate the UE due to NULL"
        "ue_mm_context\n");
    amf_cause = AMF_UE_ILLEGAL;
    rc        = amf_proc_registration_reject(ue_id, amf_cause);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
  }

  amf_ctx = &ue_mm_context->amf_context;

  registration_proc = get_nas_specific_procedure_registration(amf_ctx);
  auth_proc         = get_nas5g_common_procedure_authentication(amf_ctx);

  if (auth_proc) {
    /*    Stop Timer T3560 */
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "Timer:  Stopping Authentication Timer T3560 with id = %lu\n",
        auth_proc->T3560.id);
    amf_app_stop_timer(auth_proc->T3560.id);
    auth_proc->T3560.id = NAS5G_TIMER_INACTIVE_ID;

    nas_amf_smc_proc_autn.amf_ctx_set_security_eksi(amf_ctx, auth_proc->ksi);
    registration_proc->ksi = auth_proc->ksi;

    rc = memcmp(
        amf_ctx->_vector[auth_proc->ksi % MAX_EPS_AUTH_VECTORS].xres_star,
        msg->autn_response_parameter.response_parameter, AUTH_XRES_SIZE);

    if (rc != RETURNok) {
      is_xres_validation_failed = true;
    }

    /* As per Spec 24.501 Sec 5.4.1.3.5 If the authentication response (RES*)
     * returned by the UE is not valid, the network response depends upon the
     * type of identity used by the UE in the initial NAS message.
     *  1. If GUTI was used then the network should initiate an identification
     * procedure
     *  2. If SUCI was used then the network may send an AUTHENTICATION REJECT
     * message to the UE
     */
    if (is_xres_validation_failed) {
      auth_proc->retransmission_count++;
      OAILOG_WARNING(
          LOG_NAS_AMF, "Authentication failure due to RES,XRES mismatch \n");
      if (registration_proc &&
          (amf_ctx->reg_id_type == M5GSMobileIdentityMsg_GUTI)) {
        rc = amf_proc_identification(
            amf_ctx, (nas_amf_proc_t*) registration_proc, IDENTITY_TYPE_2_IMSI,
            amf_registration_success_identification_cb,
            amf_registration_failure_identification_cb);
      } else {
        rc = RETURNerror;
      }

      if (RETURNok != rc) {
        /*
         * Notify AMF that the authentication procedure failed
         */
        amf_sap_t amf_sap                    = {};
        amf_sap.primitive                    = AMFAS_SECURITY_REJ;
        amf_sap.u.amf_as.u.security.ue_id    = ue_id;
        amf_sap.u.amf_as.u.security.msg_type = AMF_AS_MSG_TYPE_AUTH;
        rc                                   = amf_sap_send(&amf_sap);
      }
      OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
    }

    OAILOG_DEBUG(LOG_NAS_AMF, "Authentication of the UE is Successful\n");

    /*
     * Notify AMF that the authentication procedure successfully completed
     */
    amf_sap_t amf_sap               = {};
    amf_sap.primitive               = AMFREG_COMMON_PROC_CNF;
    amf_sap.u.amf_reg.ue_id         = ue_id;
    amf_sap.u.amf_reg.ctx           = amf_ctx;
    amf_sap.u.amf_reg.notify        = true;
    amf_sap.u.amf_reg.free_proc     = true;
    amf_sap.u.amf_reg.u.common_proc = &auth_proc->amf_com_proc;
    rc                              = amf_sap_send(&amf_sap);
  } else {
    OAILOG_ERROR(LOG_NAS_AMF, "Auth proc is null");
  }
  /* Completing Authentication response and invoking Security Request
   * Invoking success directly to handle security mode command
   * */
  rc = amf_registration_success_authentication_cb(amf_ctx);
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

inline void amf_ctx_clear_attribute_present(
    amf_context_t* const ctxt, const int attribute_bit_pos) {
  ctxt->member_present_mask &= ~attribute_bit_pos;
  ctxt->member_valid_mask &= ~attribute_bit_pos;
}

void amf_ctx_clear_auth_vectors(amf_context_t* const ctxt) {
  amf_ctx_clear_attribute_present(ctxt, AMF_CTXT_MEMBER_AUTH_VECTORS);

  for (int i = 0; i < MAX_EPS_AUTH_VECTORS; i++) {
    memset((void*) &ctxt->_vector[i], 0, sizeof(ctxt->_vector[i]));
    amf_ctx_clear_attribute_present(ctxt, AMF_CTXT_MEMBER_AUTH_VECTOR0 + i);
  }

  ctxt->_security.vector_index = AMF_SECURITY_VECTOR_INDEX_INVALID;
}

int amf_auth_auth_rej(amf_ue_ngap_id_t ue_id) {
  int rc                               = RETURNerror;
  ue_m5gmm_context_s* ue_mm_context    = nullptr;
  amf_sap_t amf_sap                    = {};
  amf_sap.primitive                    = AMFAS_SECURITY_REJ;
  amf_sap.u.amf_as.u.security.ue_id    = ue_id;
  amf_sap.u.amf_as.u.security.msg_type = AMF_AS_MSG_TYPE_AUTH;
  rc                                   = amf_sap_send(&amf_sap);
  ue_mm_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  amf_free_ue_context(ue_mm_context);
  return rc;
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_proc_authentication_failure()                             **
 **                                                                        **
 ** Description: Performs the authentication failure procedure executed    **
 **      by the network.                                                   **
 **                                                                        **
 **              3GPP TS 24.501, section 5.4.1.3.7                         **
 **      Upon receiving the AUTHENTICATION FAILURE message, the            **
 **      MME shall stop timer T3560 and check the correctness of           **
 **      the RES parameter.                                                **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                         **
 **      emm_cause: Authentication failure AMF cause code                  **
 **      res:       Authentication response parameter. or auts             **
 **                 in case of sync failure                                **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                                  **
 **      Others:    amf_data, T3560                                        **
 **                                                                        **
 ***************************************************************************/
int amf_proc_authentication_failure(
    amf_ue_ngap_id_t ue_id, AuthenticationFailureMsg* msg, int amf_cause) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);

  OAILOG_DEBUG(
      LOG_NAS_AMF,
      "Authentication  procedures failed for "
      "(ue_id=" AMF_UE_NGAP_ID_FMT ")\n",
      ue_id);

  int rc                            = RETURNerror;
  ue_m5gmm_context_s* ue_mm_context = NULL;
  amf_context_t* amf_ctx            = NULL;

  ue_mm_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  if (!ue_mm_context) {
    OAILOG_WARNING(
        LOG_NAS_AMF, "Sending Auth Reject as UE MM Context is not found\n");
    rc = amf_auth_auth_rej(ue_id);
    OAILOG_WARNING(
        LOG_NAS_AMF,
        "AMF-PROC - Failed to authenticate the UE due to NULL"
        "ue_mm_context\n");
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
  }

  amf_ctx = &ue_mm_context->amf_context;
  nas5g_amf_auth_proc_t* auth_proc =
      get_nas5g_common_procedure_authentication(amf_ctx);

  if (!auth_proc) {
    OAILOG_WARNING(LOG_NAS_AMF, "authentication procedure not present\n");
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
  }

  /*	Stop Timer T3560 */
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "Timer:  Stopping Authentication Timer T3560 with id = %lu\n",
      auth_proc->T3560.id);
  amf_app_stop_timer(auth_proc->T3560.id);
  auth_proc->T3560.id = NAS5G_TIMER_INACTIVE_ID;

  OAILOG_DEBUG(
      LOG_NAS_AMF, "authentication of the ue is failed with error : %d",
      msg->m5gmm_cause.m5gmm_cause);

  switch (msg->m5gmm_cause.m5gmm_cause) {
    case AMF_CAUSE_NGKSI_ALREADY_INUSE: {
      OAILOG_WARNING(
          LOG_NAS_AMF, "Authentication failure received with NGSKI\n");
      nas_amf_registration_proc_t* registration_proc =
          get_nas_specific_procedure_registration(amf_ctx);

      amf_ctx->_security.eksi = auth_proc->ksi;
      OAILOG_DEBUG(LOG_NAS_AMF, "Updated EKSI %d\n", amf_ctx->_security.eksi);
      rc = amf_start_registration_proc_authentication(
          amf_ctx, registration_proc);
      break;
    }
    case AMF_CAUSE_MAC_FAILURE: {
      auth_proc->retransmission_count++;
      nas_amf_registration_proc_t* registration_proc =
          get_nas_specific_procedure_registration(amf_ctx);
      OAILOG_DEBUG(
          LOG_NAS_AMF,
          "Authentication failure received with failure response\n");
      if (registration_proc &&
          (amf_ctx->reg_id_type == M5GSMobileIdentityMsg_GUTI)) {
        rc = amf_proc_identification(
            amf_ctx, (nas_amf_proc_t*) registration_proc, IDENTITY_TYPE_2_IMSI,
            amf_registration_success_identification_cb,
            amf_registration_failure_identification_cb);
      } else {
        /*
         * in case of SUCI BASED REGISTRATION Send AUTH_REJECT */
        rc = RETURNerror;
      }

      if (RETURNok != rc) {
        /*
         * Notify AMF that the authentication procedure successfully completed
         */
        OAILOG_ERROR(
            LOG_NAS_AMF,
            "Sending authentication reject with cause AMF_CAUSE_MAC_FAILURE\n");
        rc = amf_auth_auth_rej(ue_id);
      }
    } break;
    case AMF_CAUSE_SYNCH_FAILURE: {
      auth_proc->retry_sync_failure++;
      if (MAX_SYNC_FAILURES <= auth_proc->retry_sync_failure) {
        rc = amf_auth_auth_rej(ue_id);
      } else {
        struct tagbstring resync_param;
        resync_param.data = (unsigned char*) calloc(1, RESYNC_PARAM_LENGTH);
        if (resync_param.data == NULL) {
          OAILOG_ERROR(
              LOG_NAS_AMF,
              "Sending authentication reject with cause "
              "AMF_CAUSE_SYNCH_FAILURE\n");
          rc = amf_auth_auth_rej(ue_id);
          OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
        }

        memcpy(
            resync_param.data,
            (amf_ctx->_vector[amf_ctx->_security.vector_index].rand),
            RAND_LENGTH_OCTETS);

        memcpy(
            (resync_param.data + RAND_LENGTH_OCTETS),
            msg->auth_failure_ie.authentication_failure_info->data,
            AUTS_LENGTH);

        start_authentication_information_procedure_synch(
            amf_ctx, auth_proc, &resync_param);
        free_wrapper(reinterpret_cast<void**>(&resync_param.data));

        amf_ctx_clear_auth_vectors(amf_ctx);
      }

    } break;
    case AMF_NON_5G_AUTHENTICATION_UNACCEPTABLE: {
      OAILOG_ERROR(
          LOG_NAS_AMF,
          "Sending authentication reject with cause AMF_CAUSE_MAC_FAILURE\n");
      rc = amf_auth_auth_rej(ue_id);
    } break;

    default: {
      OAILOG_DEBUG(LOG_NAS_AMF, "Unsupported 5gmm cause\n");
      break;
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_send_authentication_request()                             **
 **                                                                        **
 ** Description: Sends AUTHENTICATION REQUEST message and start timer T3560**
 **                                                                        **
 ** Inputs:  args: pointer to amf context                                  **
 **                handler parameters                                      **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                                  **
 **                                                                        **
 ***************************************************************************/
int amf_send_authentication_request(
    amf_context_t* amf_ctx, nas5g_amf_auth_proc_t* auth_proc) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc              = RETURNerror;
  auth_proc->T3560.id = NAS5G_TIMER_INACTIVE_ID;
  if (auth_proc) {
    /*
     * Notify AMF-AS SAP that Authentication Request message has to be sent
     * to the UE
     */
    amf_sap_t amf_sap                = {};
    amf_sap.primitive                = AMFAS_SECURITY_REQ;
    amf_sap.u.amf_as.u.security.guti = {0};
    //    amf_sap.u.amf_as.u.security.ue_id    = auth_proc->ue_id;
    amf_sap.u.amf_as.u.security.ue_id =
        PARENT_STRUCT(amf_ctx, ue_m5gmm_context_s, amf_context)->amf_ue_ngap_id;

    amf_sap.u.amf_as.u.security.msg_type = AMF_AS_MSG_TYPE_AUTH;
    amf_sap.u.amf_as.u.security.ksi      = auth_proc->ksi;
    memcpy(amf_sap.u.amf_as.u.security.rand, auth_proc->rand, AUTH_RAND_SIZE);
    memcpy(amf_sap.u.amf_as.u.security.autn, auth_proc->autn, AUTH_AUTN_SIZE);

    /*
     * Setup 5GCN NAS security data
     */
    amf_data_sec_auth.amf_as_set_security_data(
        &amf_sap.u.amf_as.u.security.sctx, &amf_ctx->_security, false, true);

    rc = amf_sap_send(&amf_sap);

    if (rc != RETURNerror) {
      OAILOG_WARNING(
          LOG_NAS_AMF,
          " T3560: Start Authentication Timer for "
          "ue id: " AMF_UE_NGAP_ID_FMT "\n",
          auth_proc->ue_id);
      auth_proc->T3560.id = amf_app_start_timer(
          AUTHENTICATION_TIMER_EXPIRY_MSECS, TIMER_REPEAT_ONCE,
          authenthication_t3560_handler, auth_proc->ue_id);
    }
    if (rc != RETURNerror) {
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

// Fetch the serving network name
static int calculate_amf_serving_network_name(
    amf_context_t* amf_ctx, unsigned char* snni) {
  uint32_t mcc              = 0;
  uint32_t mnc              = 0;
  uint32_t mnc_digit_length = 0;
  char snni_buffer[40]      = {0};

  /* Building 32 bytes of string with serving network SN
   * SN value = 5G:mnc<mnc>.mcc<mcc>.3gppnetwork.org
   * mcc and mnc are retrieved from serving network PLMN
   */

  PLMN_T_TO_MCC_MNC(amf_ctx->originating_tai.plmn, mcc, mnc, mnc_digit_length);

  uint32_t snni_buf_len =
      snprintf(snni_buffer, 40, "5G:mnc%03d.mcc%03d.3gppnetwork.org", mnc, mcc);

  if (snni_buf_len != 32) {
    OAILOG_ERROR(LOG_NAS_AMF, "Failed to create proper SNNI String: %s ", snni);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  } else {
    memcpy(snni, snni_buffer, snni_buf_len);
    OAILOG_DEBUG(LOG_NAS_AMF, "Serving network name: %s\n", snni);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_authentication_proc_success()                             **
 **                                                                        **
 ** Description: Process Authentication Success response from Subsdb       **
 **                                                                        **
 ** Inputs:  args: pointer to amf context                                  **
 **                handler parameters                                      **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                                  **
 **                                                                        **
 ***************************************************************************/
int amf_authentication_proc_success(amf_context_t* amf_ctx) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);

  nas5g_amf_auth_proc_t* auth_proc       = NULL;
  nas5g_auth_info_proc_t* auth_info_proc = NULL;
  uint8_t snni[40]                       = {0};
  int rc                                 = RETURNerror;

  /* Get Auth Proc */
  auth_proc = get_nas5g_common_procedure_authentication(amf_ctx);
  if (auth_proc == NULL) {
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
  }

  /* Get Auth Info Pro */
  auth_info_proc = get_nas5g_cn_procedure_auth_info(amf_ctx);
  if (auth_info_proc == NULL) {
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
  }

  // compute next eksi
  ksi_t eksi = 0;
  if (amf_ctx->_security.eksi < KSI_NO_KEY_AVAILABLE) {
    eksi = (amf_ctx->_security.eksi + 1) % (EKSI_MAX_VALUE + 1);
  }

  OAILOG_DEBUG(
      LOG_AMF_APP, "Security eksi:%x, eksi=%x", amf_ctx->_security.eksi, eksi);

  rc = calculate_amf_serving_network_name(amf_ctx, snni);
  if (rc != RETURNok) {
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
  }

  memcpy(
      amf_ctx->_vector[amf_ctx->_security.eksi % MAX_EPS_AUTH_VECTORS].kseaf,
      auth_info_proc->vector[0]->kseaf, KSEAF_LENGTH_OCTETS);

  memcpy(
      amf_ctx->_vector[amf_ctx->_security.eksi % MAX_EPS_AUTH_VECTORS].autn,
      auth_info_proc->vector[0]->autn, AUTN_LENGTH_OCTETS);

  memcpy(
      amf_ctx->_vector[amf_ctx->_security.eksi % MAX_EPS_AUTH_VECTORS].rand,
      auth_info_proc->vector[0]->rand, RAND_LENGTH_OCTETS);

  amf_ctx->_vector[amf_ctx->_security.eksi % MAX_EPS_AUTH_VECTORS]
      .xres_star_length = auth_info_proc->vector[0]->xres_star.size;

  memcpy(
      amf_ctx->_vector[amf_ctx->_security.eksi % MAX_EPS_AUTH_VECTORS]
          .xres_star,
      auth_info_proc->vector[0]->xres_star.data, AUTH_XRES_SIZE);

  /* Set the vector and corresponding vectors */
  amf_ctx_set_attribute_valid(amf_ctx, AMF_CTXT_MEMBER_AUTH_VECTOR0);

  if (auth_info_proc->nb_vectors > 0) {
    amf_ctx_set_attribute_valid(amf_ctx, AMF_CTXT_MEMBER_AUTH_VECTORS);
  }

  auth_proc->ksi = eksi;

  /* Send the authentication request */
  amf_send_authentication_request(amf_ctx, auth_proc);

  nas5g_delete_cn_procedure(amf_ctx, &auth_info_proc->cn_proc);

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
}

/* Timer Expiry Handler for AUTHENTHICATION Timer T3560 */
static int authenthication_t3560_handler(
    zloop_t* loop, int timer_id, void* arg) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);

  amf_context_t* amf_ctx = NULL;
  amf_ue_ngap_id_t ue_id = 0;

  if (!amf_pop_timer_arg(timer_id, &ue_id)) {
    OAILOG_WARNING(
        LOG_AMF_APP, "T3560: Invalid Timer Id expiration, Timer Id: %u\n",
        timer_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
  }

  ue_m5gmm_context_s* ue_amf_context =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  if (ue_amf_context == NULL) {
    OAILOG_DEBUG(
        LOG_NAS_AMF,
        "T3560: ue_amf_context is NULL for "
        "ue id: " AMF_UE_NGAP_ID_FMT "\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
  }

  amf_ctx = &ue_amf_context->amf_context;

  if (!(amf_ctx)) {
    OAILOG_ERROR(
        LOG_AMF_APP,
        "T3560: Timer expired no amf context for "
        "ue id: " AMF_UE_NGAP_ID_FMT "\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
  }

  nas5g_amf_auth_proc_t* auth_proc =
      get_nas5g_common_procedure_authentication(amf_ctx);
  if (auth_proc) {
    OAILOG_WARNING(
        LOG_AMF_APP,
        "T3560: timer expired timer id %lu "
        "ue id " AMF_UE_NGAP_ID_FMT "\n",
        auth_proc->T3560.id, auth_proc->ue_id);
    auth_proc->T3560.id = -1;
    /*
     *Increment the retransmission counter
     */
    auth_proc->retransmission_count += 1;
    OAILOG_WARNING(
        LOG_NAS_AMF, "T3560: Timer expired, retransmission counter = %d\n",
        auth_proc->retransmission_count);

    if (auth_proc->retransmission_count < AUTHENTICATION_COUNTER_MAX) {
      OAILOG_ERROR(
          LOG_NAS_AMF,
          "T3560: Retransmitting amf_send_authentication_request\n");
      amf_send_authentication_request(amf_ctx, auth_proc);
    } else {
      OAILOG_ERROR(
          LOG_AMF_APP,
          "T3560: Maximum retires done hence Abort the authentication "
          "procedure\n");
      amf_proc_registration_abort(amf_ctx, ue_amf_context);
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
}

}  // namespace magma5g
