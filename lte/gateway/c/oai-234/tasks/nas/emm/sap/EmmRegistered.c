/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#include <assert.h>

#include "log.h"
#include "common_defs.h"
#include "emm_fsm.h"
#include "emm_data.h"
#include "emm_regDef.h"
#include "nas_procedures.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/****************************************************************************
 **                                                                        **
 ** Name:    EmmRegistered()                                           **
 **                                                                        **
 ** Description: Handles the behaviour of the UE and the MME while the     **
 **      EMM-SAP is in EMM-REGISTERED state.                       **
 **                                                                        **
 ** Inputs:  evt:       The received EMM-SAP event                 **
 **      Others:    emm_fsm_status                             **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    emm_fsm_status                             **
 **                                                                        **
 ***************************************************************************/
int EmmRegistered(const emm_reg_t* evt) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc                 = RETURNerror;
  emm_context_t* emm_ctx = evt->ctx;

  assert(emm_fsm_get_state(emm_ctx) == EMM_REGISTERED);

  switch (evt->primitive) {
    case _EMMREG_COMMON_PROC_REQ:
      /*
       * An EMM common procedure has been initiated;
       * enter state EMM-COMMON-PROCEDURE-INITIATED.
       */
      rc = emm_fsm_set_state(
          evt->ue_id, evt->ctx, EMM_COMMON_PROCEDURE_INITIATED);
      break;

    case _EMMREG_COMMON_PROC_CNF:
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-FSM state EMM_REGISTERED - Primitive _EMMREG_COMMON_PROC_CNF is "
          "not valid\n");
      break;

    case _EMMREG_COMMON_PROC_REJ:
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-FSM state EMM_REGISTERED - Primitive _EMMREG_COMMON_PROC_REJ is "
          "not valid\n");
      break;

    case _EMMREG_COMMON_PROC_ABORT:
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-FSM state EMM_REGISTERED - Primitive _EMMREG_COMMON_PROC_ABORT "
          "is "
          "not valid\n");
      rc = emm_fsm_set_state(evt->ue_id, evt->ctx, EMM_DEREGISTERED);
      break;

    case _EMMREG_ATTACH_CNF:
      /*
       * Attach procedure successful and default EPS bearer
       * context activated;
       * enter state EMM-REGISTERED.
       */
      rc = emm_fsm_set_state(evt->ue_id, emm_ctx, EMM_REGISTERED);
      break;

    case _EMMREG_ATTACH_REJ:
      /*
       * Attach procedure failed;
       * enter state EMM-DEREGISTERED.
       */
      rc = emm_fsm_set_state(evt->ue_id, emm_ctx, EMM_DEREGISTERED);
      nas_delete_attach_procedure(emm_ctx);
      break;

    case _EMMREG_ATTACH_ABORT:
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-FSM state EMM_REGISTERED - Primitive _EMMREG_ATTACH_ABORT is "
          "not "
          "valid\n");
      break;

    case _EMMREG_DETACH_INIT:
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-FSM state EMM_REGISTERED - Primitive _EMMREG_DETACH_INIT is not "
          "valid\n");
      break;

    case _EMMREG_DETACH_REQ:
      rc = emm_fsm_set_state(evt->ue_id, emm_ctx, EMM_DEREGISTERED_INITIATED);
      break;

    case _EMMREG_DETACH_FAILED:
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-FSM state EMM_REGISTERED - Primitive _EMMREG_DETACH_FAILED is "
          "not "
          "valid\n");
      break;

    case _EMMREG_DETACH_CNF:
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-FSM state EMM_REGISTERED - Primitive _EMMREG_DETACH_CNF is not "
          "valid\n");
      break;

    case _EMMREG_TAU_REQ:
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-FSM state EMM_REGISTERED - Primitive _EMMREG_TAU_REQ is not "
          "valid\n");
      break;

    case _EMMREG_TAU_CNF:
      if (evt->u.tau.proc) {
        if ((emm_ctx) &&
            (evt->u.tau.proc->emm_spec_proc.emm_proc.base_proc.fail_out)) {
          rc = (*evt->u.tau.proc->emm_spec_proc.emm_proc.base_proc.fail_out)(
              emm_ctx, &evt->u.tau.proc->emm_spec_proc.emm_proc.base_proc);
        }

        if ((rc != RETURNerror) && (emm_ctx) && (evt->notify) &&
            (evt->u.tau.proc->emm_spec_proc.emm_proc.base_proc.failure_notif)) {
          rc = (*evt->u.tau.proc->emm_spec_proc.emm_proc.base_proc
                     .failure_notif)(emm_ctx);
        }

        if (evt->free_proc) {
          nas_delete_tau_procedure(emm_ctx);
        }
      }
      break;
    case _EMMREG_TAU_REJ:
      if (evt->u.tau.proc) {
        rc = emm_fsm_set_state(evt->ue_id, evt->ctx, EMM_DEREGISTERED);

        if ((emm_ctx) &&
            (evt->u.tau.proc->emm_spec_proc.emm_proc.base_proc.fail_out)) {
          rc = (*evt->u.tau.proc->emm_spec_proc.emm_proc.base_proc.fail_out)(
              emm_ctx, &evt->u.tau.proc->emm_spec_proc.emm_proc.base_proc);
        }

        if ((rc != RETURNerror) && (emm_ctx) && (evt->notify) &&
            (evt->u.tau.proc->emm_spec_proc.emm_proc.base_proc.failure_notif)) {
          rc = (*evt->u.tau.proc->emm_spec_proc.emm_proc.base_proc
                     .failure_notif)(emm_ctx);
        }

        if (evt->free_proc) {
          nas_delete_tau_procedure(emm_ctx);
        }
      }
      break;

    case _EMMREG_SERVICE_REQ:
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-FSM state EMM_REGISTERED - Primitive _EMMREG_SERVICE_REQ is not "
          "valid\n");
      break;

    case _EMMREG_SERVICE_CNF:
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-FSM state EMM_REGISTERED - Primitive _EMMREG_SERVICE_CNF is not "
          "valid\n");
      break;

    case _EMMREG_SERVICE_REJ:
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-FSM state EMM_REGISTERED - Primitive _EMMREG_SERVICE_REJ is not "
          "valid\n");
      break;

    case _EMMREG_LOWERLAYER_SUCCESS:
      /*
       * Data successfully delivered to the network
       */
      rc = RETURNok;
      break;

    case _EMMREG_LOWERLAYER_FAILURE:
      if (emm_ctx) {
        nas_emm_proc_t* emm_proc = nas_emm_find_procedure_by_msg_digest(
            emm_ctx, (const char*) evt->u.ll_failure.msg_digest,
            evt->u.ll_failure.digest_len, evt->u.ll_failure.msg_len);
        if (emm_proc) {
          if ((evt->notify) && (emm_proc->not_delivered)) {
            rc = (*emm_proc->not_delivered)(emm_ctx, emm_proc);
          }
        }
      }
      break;

    case _EMMREG_LOWERLAYER_RELEASE:
      nas_delete_all_emm_procedures(emm_ctx);
      rc = RETURNok;
      break;

    case _EMMREG_LOWERLAYER_NON_DELIVERY:
      if (emm_ctx) {
        nas_emm_proc_t* emm_proc = nas_emm_find_procedure_by_msg_digest(
            emm_ctx, (const char*) evt->u.non_delivery_ho.msg_digest,
            evt->u.non_delivery_ho.digest_len, evt->u.non_delivery_ho.msg_len);
        if (emm_proc) {
          if ((evt->notify) && (emm_proc->not_delivered)) {
            rc = (*emm_proc->not_delivered_ho)(emm_ctx, emm_proc);
          } else {
            rc = RETURNok;
          }
        }
      }
      break;

    default:
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-FSM state EMM_REGISTERED - Primitive is not valid (%d)\n",
          evt->primitive);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/
