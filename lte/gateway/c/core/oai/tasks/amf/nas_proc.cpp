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

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"

#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/common/common_defs.h"
#include <sstream>
#include "lte/gateway/c/core/oai/tasks/amf/amf_asDefs.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_authentication.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_sap.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_timer_management.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_recv.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_identity.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/include/amf_smf_session_context.hpp"

extern amf_config_t amf_config;
namespace magma5g {
extern task_zmq_ctx_s amf_app_task_zmq_ctx;
AmfMsg amf_msg_obj;
static int identification_t3570_handler(zloop_t* loop, int timer_id, void* arg);
static int subs_auth_retry(zloop_t* loop, int timer_id, void* arg);
status_code_e nas_proc_establish_ind(
    const amf_ue_ngap_id_t ue_id, const bool is_mm_ctx_new,
    const tai_t originating_tai, const ecgi_t ecgi,
    const m5g_rrc_establishment_cause_t as_cause, const s_tmsi_m5_t s_tmsi,
    bstring msg) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  amf_sap_t amf_sap = {};
  status_code_e rc = RETURNerror;
  if (msg) {
    /*
     * Notify the AMF procedure call manager that NAS signaling
     * connection establishment indication message has been received
     * from the Access-Stratum sublayer
     */
    amf_sap.primitive = AMFAS_ESTABLISH_REQ;
    amf_sap.u.amf_as.primitive = _AMFAS_ESTABLISH_REQ;
    amf_sap.u.amf_as.u.establish.ue_id = ue_id;
    amf_sap.u.amf_as.u.establish.is_initial = true;
    amf_sap.u.amf_as.u.establish.is_amf_ctx_new = is_mm_ctx_new;
    amf_sap.u.amf_as.u.establish.nas_msg = msg;
    amf_sap.u.amf_as.u.establish.ecgi = ecgi;
    amf_sap.u.amf_as.u.establish.tai = originating_tai;
    rc = amf_sap_send(&amf_sap);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/***************************************************************************
**                                                                        **
** Name:    nas_new_amf_procedures()                                     **
**                                                                        **
** Description: Generic function for new amf Procedures                   **
**                                                                        **
**                                                                        **
***************************************************************************/
amf_procedures_t* nas_new_amf_procedures(amf_context_t* const amf_context) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  amf_procedures_t* amf_procedures = new amf_procedures_t();
  LIST_INIT(&amf_procedures->amf_common_procs);
  OAILOG_FUNC_RETURN(LOG_AMF_APP, amf_procedures);
}

//-----------------------------------------------------------------------------
static void nas5g_delete_auth_info_procedure(
    struct amf_context_s* amf_context,
    nas5g_auth_info_proc_t** auth_info_proc) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  if (*auth_info_proc) {
    if ((*auth_info_proc)->cn_proc.base_proc.parent) {
      (*auth_info_proc)->cn_proc.base_proc.parent->child = NULL;
    }
    free_wrapper((void**)auth_info_proc);
  }
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

/***********************************************************************
 ** Name:    amf_delete_child_procedures()                            **
 **                                                                   **
 ** Description: deletes the nas registration specific child          **
 **              procedures                                           **
 **                                                                   **
 ** Inputs:  amf_ctx:   The amf context                               **
 **          parent_proc: nas 5g base proc                            **
 **                                                                   **
 ** Return:    void                                                   **
 **                                                                   **
 ***********************************************************************/
void amf_delete_child_procedures(amf_context_t* amf_ctx,
                                 struct nas5g_base_proc_t* const parent_proc) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  if (amf_ctx && amf_ctx->amf_procedures) {
    nas_amf_common_procedure_t* p1 =
        LIST_FIRST(&amf_ctx->amf_procedures->amf_common_procs);
    nas_amf_common_procedure_t* p2 = NULL;
    while (p1) {
      p2 = LIST_NEXT(p1, entries);
      if (((nas5g_base_proc_t*)p1->proc)->parent == parent_proc) {
        amf_delete_common_procedure(amf_ctx, &p1->proc);
      }
      p1 = p2;
    }
  }
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

//---------------------------------------------------------------------------------
static void delete_common_proc_by_type(nas_amf_common_proc_t* proc) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  if (proc) {
    switch (proc->type) {
      case AMF_COMM_PROC_AUTH: {
        delete (reinterpret_cast<nas5g_amf_auth_proc_t*>(proc));
      } break;
      case AMF_COMM_PROC_SMC: {
        delete (reinterpret_cast<nas_amf_smc_proc_t*>(proc));
      } break;
      case AMF_COMM_PROC_IDENT: {
        delete (reinterpret_cast<nas_amf_ident_proc_t*>(proc));
      } break;
      default: {
        OAILOG_ERROR(LOG_AMF_APP,
                     "Error: Function  received Invalid Procedure type \n");
      }
    }
  }
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

/***********************************************************************
 ** Name:    amf_delete_common_procedure()                            **
 **                                                                   **
 ** Description: deletes the nas common  procedures                   **
 **                                                                   **
 ** Inputs:  amf context                                              **
 **          proc: nas amf common proc                                **
 **                                                                   **
 **                                                                   **
 ** Return:    void                                                   **
 **                                                                   **
 ***********************************************************************/
void amf_delete_common_procedure(amf_context_t* amf_ctx,
                                 nas_amf_common_proc_t** proc) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  if (proc && *proc) {
    switch ((*proc)->type) {
      case AMF_COMM_PROC_AUTH: {
      } break;
      case AMF_COMM_PROC_SMC: {
      } break;
      case AMF_COMM_PROC_IDENT: {
      } break;
      default: {
        OAILOG_ERROR(LOG_AMF_APP,
                     "Error: Function  received Invalid Procedure type \n");
      }
    }
  }

  // remove proc from list
  if (amf_ctx->amf_procedures) {
    nas_amf_common_procedure_t* p1 =
        LIST_FIRST(&amf_ctx->amf_procedures->amf_common_procs);
    nas_amf_common_procedure_t* p2 = NULL;

    // 2 methods: this one, the other: use parent struct macro and LIST_REMOVE
    // without searching matching element in the list
    while (p1) {
      p2 = LIST_NEXT(p1, entries);
      if (p1->proc == (nas_amf_common_proc_t*)(*proc)) {
        LIST_REMOVE(p1, entries);
        delete_common_proc_by_type(p1->proc);
        delete (p1);
        return;
      }
      p1 = p2;
    }
    nas_amf_procedure_gc(amf_ctx);
  }

  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

/***********************************************************************
 ** Name:    nas5g_delete_common_procedures()                         **
 **                                                                   **
 ** Description: deletes all nas common  procedures                   **
 **                                                                   **
 ** Inputs:  amf_context                                              **
 **                                                                   **
 **                                                                   **
 ** Outputs:     None                                                 **
 **      Return:    void                                              **
 **      Others:    None                                              **
 **                                                                   **
 ***********************************************************************/

static void nas5g_delete_common_procedures(amf_context_t* amf_context) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  // remove proc from list

  if (amf_context->amf_procedures) {
    nas_amf_common_procedure_t* p1 =
        LIST_FIRST(&amf_context->amf_procedures->amf_common_procs);
    nas_amf_common_procedure_t* p2 = NULL;
    while (p1) {
      p2 = LIST_NEXT(p1, entries);
      LIST_REMOVE(p1, entries);

      switch (p1->proc->type) {
        case AMF_COMM_PROC_AUTH: {
          nas5g_amf_auth_proc_t* auth_proc = (nas5g_amf_auth_proc_t*)p1->proc;
          if (auth_proc->T3560.id != NAS5G_TIMER_INACTIVE_ID) {
            amf_app_stop_timer(auth_proc->T3560.id);
          }
        } break;
        case AMF_COMM_PROC_SMC: {
          nas_amf_smc_proc_t* smc_proc = (nas_amf_smc_proc_t*)(p1->proc);
          if (smc_proc->T3560.id != NAS5G_TIMER_INACTIVE_ID) {
            amf_app_stop_timer(smc_proc->T3560.id);
          }
        } break;
        case AMF_COMM_PROC_IDENT: {
          nas_amf_ident_proc_t* ident_proc = (nas_amf_ident_proc_t*)(p1->proc);
          if (ident_proc->T3570.id != NAS5G_TIMER_INACTIVE_ID) {
            amf_app_stop_timer(ident_proc->T3570.id);
          }
        } break;
        default:;
      }

      delete_common_proc_by_type(p1->proc);
      delete (p1);

      p1 = p2;
    }
    nas_amf_procedure_gc(amf_context);
  }
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

/***************************************************************************
**                                                                        **
** Name:    nas5g_delete_cn_procedure()                                   **
**                                                                        **
** Description: Generic function to delete core network procedure         **
** Input : Specifc cn type to be deleted                                  **
**                                                                        **
**                                                                        **
***************************************************************************/
void nas5g_delete_cn_procedure(struct amf_context_s* amf_context,
                               nas5g_cn_proc_t* cn_proc) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  if (amf_context->amf_procedures) {
    nas5g_cn_procedure_t* p1 =
        LIST_FIRST(&amf_context->amf_procedures->cn_procs);
    nas5g_cn_procedure_t* p2 = NULL;

    while (p1) {
      p2 = LIST_NEXT(p1, entries);
      if (p1->proc == cn_proc) {
        switch (cn_proc->type) {
          case CN5G_PROC_AUTH_INFO:
            nas5g_delete_auth_info_procedure(
                amf_context, (nas5g_auth_info_proc_t**)&cn_proc);
            break;
          case CN5G_PROC_NONE:
            free_wrapper((void**)&cn_proc);
            break;
          default:;
        }
        LIST_REMOVE(p1, entries);
        free_wrapper((void**)&p1);
        OAILOG_FUNC_OUT(LOG_AMF_APP);
      }
      p1 = p2;
    }
    nas_amf_procedure_gc(amf_context);
  }
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

/***************************************************************************
**                                                                        **
** Name:    nas5g_delete_cn_procedures()                                  **
**                                                                        **
** Description: Generic function to delete all cn procedures              **
**              at amf_context level                                      **
**                                                                        **
**                                                                        **
***************************************************************************/
static void nas5g_delete_cn_procedures(struct amf_context_s* amf_context) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  if (amf_context->amf_procedures) {
    nas5g_cn_procedure_t* p1 =
        LIST_FIRST(&amf_context->amf_procedures->cn_procs);
    nas5g_cn_procedure_t* p2 = NULL;
    while (p1) {
      p2 = LIST_NEXT(p1, entries);
      switch (p1->proc->type) {
        case CN5G_PROC_AUTH_INFO:
          nas5g_delete_auth_info_procedure(amf_context,
                                           (nas5g_auth_info_proc_t**)&p1->proc);
          break;
        default:
          break;
      }
      LIST_REMOVE(p1, entries);
      delete (p1);
      p1 = p2;
    }
    nas_amf_procedure_gc(amf_context);
  }
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

//-----------------------------------------------------------------------------
void nas_delete_all_amf_procedures(amf_context_t* const amf_context) {
  OAILOG_FUNC_IN(LOG_AMF_APP);

  if (amf_context->amf_procedures) {
    nas5g_delete_cn_procedures(amf_context);
    nas5g_delete_common_procedures(amf_context);

    amf_delete_registration_proc(amf_context);

    if (amf_context->amf_procedures) {
      delete amf_context->amf_procedures;
    }
  }
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

/***************************************************************************
**                                                                        **
** Name:    nas5g_new_cn_auth_info_procedure()                            **
**                                                                        **
** Description: Generic function for new auth info  Procedure             **
**                                                                        **
**                                                                        **
***************************************************************************/
nas5g_auth_info_proc_t* nas5g_new_cn_auth_info_procedure(
    amf_context_t* const amf_context) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  if (!(amf_context->amf_procedures)) {
    amf_context->amf_procedures = nas_new_amf_procedures(amf_context);
  }
  nas5g_auth_info_proc_t* auth_info_proc =
      (nas5g_auth_info_proc_t*)calloc(1, sizeof(nas5g_auth_info_proc_t));

  auth_info_proc->cn_proc.base_proc.type = NAS_PROC_TYPE_CN;
  auth_info_proc->cn_proc.type = CN5G_PROC_AUTH_INFO;

  nas5g_cn_procedure_t* wrapper =
      (nas5g_cn_procedure_t*)calloc(1, sizeof(*wrapper));
  if (wrapper) {
    wrapper->proc = &auth_info_proc->cn_proc;
    LIST_INSERT_HEAD(&amf_context->amf_procedures->cn_procs, wrapper, entries);
    OAILOG_FUNC_RETURN(LOG_AMF_APP, auth_info_proc);
  } else {
    free_wrapper((void**)&auth_info_proc);
  }
  OAILOG_FUNC_RETURN(LOG_AMF_APP, NULL);
}

/***************************************************************************
**                                                                        **
** Name:    get_nas_specific_procedure_registration()                     **
**                                                                        **
** Description: Function for NAS Specific Procedure Registration          **
**                                                                        **
**                                                                        **
***************************************************************************/
nas_amf_registration_proc_t* get_nas_specific_procedure_registration(
    const amf_context_t* ctxt) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  if ((ctxt) && (ctxt->amf_procedures) &&
      (ctxt->amf_procedures->amf_specific_proc) &&
      ((AMF_SPEC_PROC_TYPE_REGISTRATION ==
        ctxt->amf_procedures->amf_specific_proc->type))) {
    OAILOG_FUNC_RETURN(
        LOG_AMF_APP,
        (nas_amf_registration_proc_t*)ctxt->amf_procedures->amf_specific_proc);
  }
  OAILOG_FUNC_RETURN(LOG_AMF_APP, NULL);
}

/***************************************************************************
**                                                                        **
** Name:    is_nas_specific_procedure_registration_running()              **
**                                                                        **
** Description: Function to check if NAS procedure registration running   **
**                                                                        **
**                                                                        **
***************************************************************************/
bool is_nas_specific_procedure_registration_running(const amf_context_t* ctxt) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  if ((ctxt) && (ctxt->amf_procedures) &&
      (ctxt->amf_procedures->amf_specific_proc) &&
      ((AMF_SPEC_PROC_TYPE_REGISTRATION ==
        ctxt->amf_procedures->amf_specific_proc->type)))
    OAILOG_FUNC_RETURN(LOG_AMF_APP, true);
  OAILOG_FUNC_RETURN(LOG_AMF_APP, false);
}

/***************************************************************************
**                                                                        **
** Name:    nas5g_new_identification_procedure()                          **
**                                                                        **
** Description: Invokes Function for new identification procedure         **
**                                                                        **
**                                                                        **
***************************************************************************/
nas_amf_ident_proc_t* nas5g_new_identification_procedure(
    amf_context_t* const amf_context) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  if (!(amf_context->amf_procedures)) {
    amf_context->amf_procedures = nas_new_amf_procedures(amf_context);
  }
  nas_amf_ident_proc_t* ident_proc = new nas_amf_ident_proc_t;
  ident_proc->amf_com_proc.amf_proc.type = NAS_AMF_PROC_TYPE_COMMON;
  ident_proc->T3570.msec = 1000 * amf_config.nas_config.t3570_sec;
  ident_proc->T3570.id = AMF_APP_TIMER_INACTIVE_ID;
  ident_proc->amf_com_proc.type = AMF_COMM_PROC_IDENT;
  nas_amf_common_procedure_t* wrapper = new nas_amf_common_procedure_t;
  if (wrapper) {
    wrapper->proc = &ident_proc->amf_com_proc;
    LIST_INSERT_HEAD(&amf_context->amf_procedures->amf_common_procs, wrapper,
                     entries);
    OAILOG_FUNC_RETURN(LOG_AMF_APP, ident_proc);
  } else {
    free_wrapper((void**)&ident_proc);
  }
  OAILOG_FUNC_RETURN(LOG_AMF_APP, ident_proc);
}

/***************************************************************************
** Description: Sends IDENTITY REQUEST message.                           **
**                                                                        **
** Inputs:  args:      handler parameters                                 **
**      Others:    None                                                   **
**                                                                        **
** Outputs:        None                                                   **
**      Return:    None                                                   **
***************************************************************************/
static status_code_e amf_identification_request(
    nas_amf_ident_proc_t* const proc) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  amf_sap_t amf_sap = {};
  status_code_e rc = RETURNok;
  proc->T3570.id = NAS5G_TIMER_INACTIVE_ID;
  OAILOG_DEBUG(LOG_AMF_APP, "Sending AS IDENTITY_REQUEST\n");
  /*
   * Notify AMF-AS SAP that Identity Request message has to be sent
   * to the UE
   */
  amf_sap.primitive = AMFAS_SECURITY_REQ;
  amf_sap.u.amf_as.u.security.ue_id = proc->ue_id;
  amf_sap.u.amf_as.u.security.msg_type = AMF_AS_MSG_TYPE_IDENT;
  amf_sap.u.amf_as.u.security.ident_type = proc->identity_type;
  amf_sap.u.amf_as.u.security.sctx.is_knas_int_present = true;
  amf_sap.u.amf_as.u.security.sctx.is_knas_enc_present = true;
  amf_sap.u.amf_as.u.security.sctx.is_new = true;
  rc = amf_sap_send(&amf_sap);

  if (rc != RETURNerror) {
    /*
     * Start Identification T3570 timer
     */
    OAILOG_DEBUG(LOG_AMF_APP,
                 "AMF_TEST: Timer: Starting Identity timer T3570 \n");
    proc->T3570.id =
        amf_app_start_timer(IDENTITY_TIMER_EXPIRY_MSECS, TIMER_REPEAT_ONCE,
                            identification_t3570_handler, proc->ue_id);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/* Identification Timer T3570 Expiry Handler */
static int identification_t3570_handler(zloop_t* loop, int timer_id,
                                        void* arg) {
  amf_ue_ngap_id_t ue_id = 0;
  amf_context_t* amf_ctx = NULL;
  OAILOG_FUNC_IN(LOG_NAS_AMF);

  if (!amf_pop_timer_arg(timer_id, &ue_id)) {
    OAILOG_WARNING(LOG_AMF_APP,
                   "T3570: Invalid Timer Id expiration, timer Id: %u\n",
                   timer_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
  }

  ue_m5gmm_context_s* ue_amf_context =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  if (ue_amf_context == NULL) {
    OAILOG_DEBUG(LOG_AMF_APP,
                 "T3570: ue_amf_context is NULL for UE ID: " AMF_UE_NGAP_ID_FMT,
                 ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
  }

  amf_ctx = &ue_amf_context->amf_context;
  if (!(amf_ctx)) {
    OAILOG_ERROR(
        LOG_AMF_APP,
        "T3570: timer expired No AMF context for UE ID: " AMF_UE_NGAP_ID_FMT,
        ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
  }

  nas_amf_ident_proc_t* ident_proc =
      get_5g_nas_common_procedure_identification(amf_ctx);

  if (ident_proc) {
    OAILOG_WARNING(
        LOG_AMF_APP,
        "T3570: Timer expired for timer id %lu for UE ID " AMF_UE_NGAP_ID_FMT,
        ident_proc->T3570.id, ident_proc->ue_id);
    ident_proc->T3570.id = NAS5G_TIMER_INACTIVE_ID;
    /*
     * Increment the retransmission counter
     */
    ident_proc->retransmission_count += 1;
    OAILOG_ERROR(LOG_AMF_APP,
                 "T3570: Incrementing retransmission_count to %d\n",
                 ident_proc->retransmission_count);

    if (ident_proc->retransmission_count < IDENTIFICATION_COUNTER_MAX) {
      /*
       * Send identity request message to the UE
       */
      OAILOG_ERROR(
          LOG_AMF_APP,
          "T3570: timer has expired retransmitting Identification request \n");
      amf_identification_request(ident_proc);
    } else {
      /*
       * Abort the identification procedure
       */
      OAILOG_ERROR(LOG_AMF_APP,
                   "T3570: Maximum retires:%d, done hence Abort the "
                   "identification "
                   "procedure\n",
                   ident_proc->retransmission_count);
      amf_proc_registration_abort(amf_ctx, ue_amf_context);
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
}

//-------------------------------------------------------------------------------------
status_code_e amf_proc_identification(amf_context_t* const amf_context,
                                      nas_amf_proc_t* const amf_proc,
                                      const identity_type2_t type,
                                      success_cb_t success,
                                      failure_cb_t failure) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  status_code_e rc = RETURNerror;
  amf_context->amf_fsm_state = AMF_REGISTERED;
  if ((amf_context) && ((AMF_DEREGISTERED == amf_context->amf_fsm_state) ||
                        (AMF_REGISTERED == amf_context->amf_fsm_state))) {
    amf_ue_ngap_id_t ue_id =
        PARENT_STRUCT(amf_context, ue_m5gmm_context_s, amf_context)
            ->amf_ue_ngap_id;
    nas_amf_ident_proc_t* ident_proc =
        nas5g_new_identification_procedure(amf_context);
    if (ident_proc) {
      if (amf_proc) {
        if ((NAS_AMF_PROC_TYPE_SPECIFIC == amf_proc->type) &&
            (AMF_SPEC_PROC_TYPE_REGISTRATION ==
             ((nas_amf_specific_proc_t*)amf_proc)->type)) {
          ident_proc->is_cause_is_registered = true;
        }
      }
      ident_proc->identity_type = type;
      ident_proc->retransmission_count = 0;
      ident_proc->ue_id = ue_id;
      (reinterpret_cast<nas5g_base_proc_t*>(ident_proc))->parent =
          reinterpret_cast<nas5g_base_proc_t*>(amf_proc);
      ident_proc->amf_com_proc.amf_proc.delivered = NULL;
      ident_proc->amf_com_proc.amf_proc.base_proc.success_notif = success;
      ident_proc->amf_com_proc.amf_proc.base_proc.failure_notif = failure;
      ident_proc->amf_com_proc.amf_proc.base_proc.fail_in = NULL;
    }
    rc = amf_identification_request(ident_proc);

    if (rc != RETURNerror) {
      /*
       * Notify 5G CN that common procedure has been initiated
       */
      amf_sap_t amf_sap = {};
      amf_sap.primitive = AMFREG_COMMON_PROC_REQ;
      amf_sap.u.amf_reg.ue_id = ue_id;
      amf_sap.u.amf_reg.ctx = amf_context;
      rc = amf_sap_send(&amf_sap);
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/***************************************************************************
 **                                                                        **
 ** Name:    amf_nas_proc_implicit_deregister_ue_ind()                     **
 **                                                                        **
 ** Description: Nas CN procedure to send implicit delete message          **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
status_code_e amf_nas_proc_implicit_deregister_ue_ind(amf_ue_ngap_id_t ue_id) {
  status_code_e rc = RETURNerror;
  amf_sap_t amf_sap = {};

  OAILOG_FUNC_IN(LOG_AMF_APP);
  amf_sap.primitive = AMFCN_IMPLICIT_DEREGISTER_UE;
  amf_sap.u.amf_cn.u.amf_cn_implicit_deregister.ue_id = ue_id;
  rc = amf_sap_send(&amf_sap);
  OAILOG_FUNC_RETURN(LOG_AMF_APP, rc);
}

/***************************************************************************
 **                                                                        **
 ** Name:    amf_nas_proc_auth_param_res()                                 **
 **                                                                        **
 ** Description: Process the authentication response received from         **
 **              Subscriberdb                                              **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
status_code_e amf_nas_proc_auth_param_res(amf_ue_ngap_id_t amf_ue_ngap_id,
                                          uint8_t nb_vectors,
                                          m5gauth_vector_t* vectors) {
  OAILOG_FUNC_IN(LOG_AMF_APP);

  status_code_e rc = RETURNerror;
  amf_sap_t amf_sap = {};
  amf_cn_auth_res_t amf_cn_auth_res = {};

  amf_cn_auth_res.ue_id = amf_ue_ngap_id;
  amf_cn_auth_res.nb_vectors = nb_vectors;
  for (int i = 0; i < nb_vectors; i++) {
    amf_cn_auth_res.vector[i] = &vectors[i];
  }

  amf_sap.primitive = AMFCN_AUTHENTICATION_PARAM_RES;
  amf_sap.u.amf_cn.u.auth_res = &amf_cn_auth_res;
  rc = amf_sap_send(&amf_sap);

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

static int subs_auth_retry(zloop_t* loop, int timer_id, void* arg) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  amf_ue_ngap_id_t ue_id = 0;
  amf_context_t* amf_ctxt_p = nullptr;
  int amf_cause = -1;
  nas5g_auth_info_proc_t* auth_info_proc = nullptr;
  int rc = RETURNerror;
  ue_m5gmm_context_s* ue_mm_context = nullptr;
  if (!amf_pop_timer_arg(timer_id, &ue_id)) {
    OAILOG_WARNING(LOG_AMF_APP,
                   "auth_retry_timer: Invalid Timer Id expiration, Timer Id: "
                   "%d and UE id: " AMF_UE_NGAP_ID_FMT "\n",
                   timer_id, ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
  }
  ue_mm_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  if (!ue_mm_context) {
    OAILOG_WARNING(
        LOG_NAS_AMF,
        "AMF-PROC - Failed authentication request for UE id " AMF_UE_NGAP_ID_FMT
        "due to NULL"
        "ue_context\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }
  amf_ctxt_p = &ue_mm_context->amf_context;
  if (!amf_ctxt_p) {
    OAILOG_WARNING(LOG_NAS_AMF,
                   "AMF-PROC - Failed authentication request for UE "
                   "id= " AMF_UE_NGAP_ID_FMT
                   "due to NULL"
                   "amf_ctxt_p\n",
                   ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }
  if (amf_ctxt_p->auth_retry_count < amf_config.auth_retry_max_count) {
    amf_ctxt_p->auth_retry_count++;
    OAILOG_INFO(LOG_AMF_APP,
                "auth_retry_timer: Incrementing auth_retry_count to %u\n",
                amf_ctxt_p->auth_retry_count);
    rc = amf_authentication_request_sent(ue_id);
    if (rc != RETURNok) {
      OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
    }
    amf_ctxt_p->auth_retry_timer.id =
        amf_app_start_timer(amf_config.auth_retry_interval, TIMER_REPEAT_ONCE,
                            subs_auth_retry, ue_id);
  } else {
    auth_info_proc = get_nas5g_cn_procedure_auth_info(amf_ctxt_p);
    OAILOG_ERROR(
        LOG_NAS_AMF,
        "auth_retry_timer is expired . Authentication reject with cause "
        "AMF_UE_ILLEGAL for ue_id " AMF_UE_NGAP_ID_FMT "\n",
        ue_id);
    amf_cause = AMF_UE_ILLEGAL;

    rc = amf_proc_registration_reject(ue_id, amf_cause);
    if (rc != RETURNok) {
      OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
    }
    if (auth_info_proc) {
      nas5g_delete_cn_procedure(amf_ctxt_p, &auth_info_proc->cn_proc);
    }

    amf_free_ue_context(ue_mm_context);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
}

status_code_e amf_nas_proc_authentication_info_answer(
    itti_amf_subs_auth_info_ans_t* aia) {
  imsi64_t imsi64 = INVALID_IMSI64;
  status_code_e rc = RETURNok;
  amf_context_t* amf_ctxt_p = NULL;
  ue_m5gmm_context_s* ue_5gmm_context_p = NULL;
  int amf_cause = -1;
  nas5g_auth_info_proc_t* auth_info_proc = NULL;
  OAILOG_FUNC_IN(LOG_AMF_APP);

  IMSI_STRING_TO_IMSI64((char*)aia->imsi, &imsi64);

  OAILOG_DEBUG(LOG_AMF_APP, "Handling imsi " IMSI_64_FMT "\n", imsi64);

  ue_5gmm_context_p = lookup_ue_ctxt_by_imsi(imsi64);

  if (ue_5gmm_context_p) {
    amf_ctxt_p = &ue_5gmm_context_p->amf_context;
  }

  if (!(amf_ctxt_p)) {
    OAILOG_ERROR(LOG_NAS_AMF,
                 "That's embarrassing as we don't know this IMSI\n");
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }

  amf_ue_ngap_id_t amf_ue_ngap_id = ue_5gmm_context_p->amf_ue_ngap_id;

  OAILOG_DEBUG(
      LOG_NAS_AMF,
      "Received Authentication Information Answer from Subscriberdb for"
      " UE ID = " AMF_UE_NGAP_ID_FMT,
      amf_ue_ngap_id);
  if (aia->auth_info.nb_of_vectors) {
    nas5g_amf_auth_proc_t* auth_proc =
        get_nas5g_common_procedure_authentication(amf_ctxt_p);
    if (auth_proc) {
      free_wrapper(reinterpret_cast<void**>(&auth_proc->auts.data));
    }
    if ((NAS5G_TIMER_INACTIVE_ID != amf_ctxt_p->auth_retry_timer.id) &&
        (0 != amf_ctxt_p->auth_retry_timer.id)) {
      OAILOG_DEBUG(LOG_NAS_AMF, "Stopping: Timer auth_retry_timer.\n");
      amf_app_stop_timer(amf_ctxt_p->auth_retry_timer.id);
      amf_ctxt_p->auth_retry_timer.id = NAS5G_TIMER_INACTIVE_ID;
    }
    /*
     * Check that list is not empty and contain at most MAX_EPS_AUTH_VECTORS
     * elements
     */
    if (aia->auth_info.nb_of_vectors > MAX_EPS_AUTH_VECTORS) {
      OAILOG_WARNING(
          LOG_NAS_AMF,
          "nb_of_vectors should be lesser than max_eps_auth_vectors");
      OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
    }

    OAILOG_DEBUG(LOG_NAS_AMF,
                 "INFORMING NAS ABOUT AUTH RESP SUCCESS got %u vector(s)\n",
                 aia->auth_info.nb_of_vectors);
    rc = amf_nas_proc_auth_param_res(amf_ue_ngap_id,
                                     aia->auth_info.nb_of_vectors,
                                     aia->auth_info.m5gauth_vector);
  } else {
    /* Get Auth Info Pro */
    auth_info_proc = get_nas5g_cn_procedure_auth_info(amf_ctxt_p);
    amf_ctxt_p->auth_retry_count = 0;
    if (aia->result != DIAMETER_TOO_BUSY) {
      if ((NAS5G_TIMER_INACTIVE_ID != amf_ctxt_p->auth_retry_timer.id) &&
          (0 != amf_ctxt_p->auth_retry_timer.id)) {
        OAILOG_DEBUG(LOG_NAS_AMF, "Stopping: Timer auth_retry_timer.\n");
        amf_app_stop_timer(amf_ctxt_p->auth_retry_timer.id);
        amf_ctxt_p->auth_retry_timer.id = NAS5G_TIMER_INACTIVE_ID;
      }
      OAILOG_ERROR(
          LOG_NAS_AMF,
          "result=%d, nb_of_vectors received is zero from subscriberdb",
          aia->result);
      amf_cause = AMF_UE_ILLEGAL;
      rc = amf_proc_registration_reject(amf_ue_ngap_id, amf_cause);
      if (auth_info_proc) {
        nas5g_delete_cn_procedure(amf_ctxt_p, &auth_info_proc->cn_proc);
      }
      amf_free_ue_context(ue_5gmm_context_p);
      OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
    } else {
      amf_ctxt_p->auth_retry_timer.id =
          amf_app_start_timer(amf_config.auth_retry_interval, TIMER_REPEAT_ONCE,
                              subs_auth_retry, aia->ue_id);
      rc = RETURNok;
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

int amf_decrypt_msin_info_answer(itti_amf_decrypted_msin_info_ans_t* aia) {
  imsi64_t imsi64 = INVALID_IMSI64;
  status_code_e rc = RETURNerror;
  amf_context_t* amf_ctxt_p = NULL;
  ue_m5gmm_context_s* ue_context = NULL;

  // Local imsi to be put in imsi defined in 3gpp_23.003.h
  supi_as_imsi_t supi_imsi;
  amf_guti_m5g_t amf_guti;
  const bool is_amf_ctx_new = true;
  OAILOG_FUNC_IN(LOG_AMF_APP);

  ue_context = amf_ue_context_exists_amf_ue_ngap_id(aia->ue_id);

  if (ue_context) {
    amf_ctxt_p = &ue_context->amf_context;
  }

  if (!(amf_ctxt_p)) {
    OAILOG_ERROR(LOG_NAS_AMF,
                 "That's embarrassing as we don't know this IMSI\n");
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }

  amf_ue_ngap_id_t amf_ue_ngap_id = ue_context->amf_ue_ngap_id;

  OAILOG_DEBUG(
      LOG_NAS_AMF,
      "Received decrypted imsi Information Answer from Subscriberdb for"
      " UE ID = " AMF_UE_NGAP_ID_FMT,
      amf_ue_ngap_id);

  amf_registration_request_ies_t* params =
      new (amf_registration_request_ies_t)();

  params->imsi = new imsi_t();

  supi_imsi.plmn.mcc_digit1 =
      ue_context->amf_context.m5_guti.guamfi.plmn.mcc_digit1;
  supi_imsi.plmn.mcc_digit2 =
      ue_context->amf_context.m5_guti.guamfi.plmn.mcc_digit2;
  supi_imsi.plmn.mcc_digit3 =
      ue_context->amf_context.m5_guti.guamfi.plmn.mcc_digit3;
  supi_imsi.plmn.mnc_digit1 =
      ue_context->amf_context.m5_guti.guamfi.plmn.mnc_digit1;
  supi_imsi.plmn.mnc_digit2 =
      ue_context->amf_context.m5_guti.guamfi.plmn.mnc_digit2;
  supi_imsi.plmn.mnc_digit3 =
      ue_context->amf_context.m5_guti.guamfi.plmn.mnc_digit3;

  supi_imsi.msin[0] =
      (uint8_t)(((aia->msin[0] - '0') << 4) | (aia->msin[1] - '0'));
  supi_imsi.msin[1] =
      (uint8_t)(((aia->msin[2] - '0') << 4) | (aia->msin[3] - '0'));
  supi_imsi.msin[2] =
      (uint8_t)(((aia->msin[4] - '0') << 4) | (aia->msin[5] - '0'));
  supi_imsi.msin[3] =
      (uint8_t)(((aia->msin[6] - '0') << 4) | (aia->msin[7] - '0'));
  supi_imsi.msin[4] =
      (uint8_t)(((aia->msin[8] - '0') << 4) | (aia->msin[9] - '0'));

  // Copy entire supi_imsi to param->imsi->u.value
  memcpy(&params->imsi->u.value, &supi_imsi, IMSI_BCD8_SIZE);

  if (supi_imsi.plmn.mnc_digit3 != 0xf) {
    params->imsi->u.value[0] = ((supi_imsi.plmn.mcc_digit1 << 4) & 0xf0) |
                               (supi_imsi.plmn.mcc_digit2 & 0xf);
    params->imsi->u.value[1] = ((supi_imsi.plmn.mcc_digit3 << 4) & 0xf0) |
                               (supi_imsi.plmn.mnc_digit1 & 0xf);
    params->imsi->u.value[2] = ((supi_imsi.plmn.mnc_digit2 << 4) & 0xf0) |
                               (supi_imsi.plmn.mnc_digit3 & 0xf);
  }

  imsi64 = amf_imsi_to_imsi64(params->imsi);
  ue_context->amf_context.imsi64 = imsi64;

  amf_app_generate_guti_on_supi(&amf_guti, &supi_imsi);
  amf_ue_context_on_new_guti(ue_context,
                             reinterpret_cast<guti_m5_t*>(&amf_guti));

  ue_context->amf_context.m5_guti.m_tmsi = amf_guti.m_tmsi;
  ue_context->amf_context.m5_guti.guamfi = amf_guti.guamfi;

  OAILOG_DEBUG(LOG_AMF_APP, "Handling imsi" IMSI_64_FMT "\n", imsi64);

  params->decode_status = ue_context->amf_context.decode_status;
  imsi_t* p_imsi = params->imsi;
  imeisv_t* p_imeisv = NULL;
  tmsi_t* p_tmsi = NULL;
  imei_t* p_imei = NULL;
  guti_m5_t* amf_ctx_guti = reinterpret_cast<guti_m5_t*>(&amf_guti);
  nas_amf_ident_proc_t* ident_proc =
      get_5g_nas_common_procedure_identification(amf_ctxt_p);
  if (ident_proc != NULL) {
    rc = amf_proc_identification_complete(aia->ue_id, p_imsi, p_imei, p_imeisv,
                                          reinterpret_cast<uint32_t*>(p_tmsi),
                                          amf_ctx_guti);
    delete (params->imsi);
    delete (params);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
  } else {
    rc = amf_proc_registration_request(aia->ue_id, is_amf_ctx_new, params);
    if (rc == RETURNerror) {
      OAILOG_ERROR(LOG_AMF_APP,
                   "processing registration request failed for ue-id "
                   ": " AMF_UE_NGAP_ID_FMT,
                   aia->ue_id);
    }
  }
  OAILOG_FUNC_RETURN(LOG_AMF_APP, rc);
}

status_code_e amf_handle_s6a_update_location_ans(
    const s6a_update_location_ans_t* ula_pP) {
  imsi64_t imsi64 = INVALID_IMSI64;
  amf_context_t* amf_ctxt_p = NULL;
  ue_m5gmm_context_s* ue_mm_context = NULL;
  OAILOG_FUNC_IN(LOG_AMF_APP);

  IMSI_STRING_TO_IMSI64((char*)ula_pP->imsi, &imsi64);

  ue_mm_context = lookup_ue_ctxt_by_imsi(imsi64);

  if (ue_mm_context) {
    amf_ctxt_p = &ue_mm_context->amf_context;
  }

  if (!(amf_ctxt_p)) {
    OAILOG_ERROR(LOG_NAS_AMF, "IMSI is invalid\n");
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }

  amf_ue_ngap_id_t amf_ue_ngap_id = ue_mm_context->amf_ue_ngap_id;

  // Validating whether the apn_config sent from ue and saved in amf_ctx is
  // present in s6a_update_location_ans_t received from subscriberdb.
  memcpy(&amf_ctxt_p->apn_config_profile,
         &ula_pP->subscription_data.apn_config_profile,
         sizeof(apn_config_profile_t));
  OAILOG_DEBUG(LOG_NAS_AMF,
               "Received update location Answer from Subscriberdb for"
               " ue_id = " AMF_UE_NGAP_ID_FMT,
               amf_ue_ngap_id);

  amf_smf_context_ue_aggregate_max_bit_rate_set(
      amf_ctxt_p, ula_pP->subscription_data.subscribed_ambr);

  OAILOG_DEBUG(LOG_NAS_AMF,
               "Received UL rate %" PRIu64 " and DL rate %" PRIu64
               "and BR unit: %d \n",
               ula_pP->subscription_data.subscribed_ambr.br_ul,
               ula_pP->subscription_data.subscribed_ambr.br_dl,
               ula_pP->subscription_data.subscribed_ambr.br_unit);

  /* FSM takes care of sending registration accept */
  ue_state_handle_message_initial(COMMON_PROCEDURE_INITIATED1,
                                  STATE_EVENT_SEC_MODE_COMPLETE, SESSION_NULL,
                                  ue_mm_context, amf_ctxt_p);

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
}

/* Cleanup all procedures in amf_context */
void amf_nas_proc_clean_up(ue_m5gmm_context_s* ue_context_p) {
  // Delete registration procedures
  amf_delete_registration_proc(&(ue_context_p->amf_context));
}

void nas_amf_procedure_gc(amf_context_t* const amf_context) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  if (LIST_EMPTY(&amf_context->amf_procedures->amf_common_procs) &&
      LIST_EMPTY(&amf_context->amf_procedures->cn_procs) &&
      (!amf_context->amf_procedures->amf_specific_proc)) {
    delete amf_context->amf_procedures;
    amf_context->amf_procedures = nullptr;
  }
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

}  // namespace magma5g
