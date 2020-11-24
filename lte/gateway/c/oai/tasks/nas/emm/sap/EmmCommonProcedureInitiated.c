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

/*****************************************************************************

  Source      EmmCommonProcedureInitiated.c

  Version     0.1

  Date        2012/10/03

  Product     NAS stack

  Subsystem   EPS Mobility Management

  Author      Frederic Maurel, Lionel GAUTHIER

  Description Implements the EPS Mobility Management procedures executed
        when the EMM-SAP is in EMM-COMMON-PROCEDURE-INITIATED state.

        In EMM-COMMON-PROCEDURE-INITIATED state, the MME has started
        a common EMM procedure and is waiting for a response from the
        UE.

*****************************************************************************/

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
 ** Name:    EmmCommonProcedureInitiated()                             **
 **                                                                        **
 ** Description: Handles the behaviour of the MME while the EMM-SAP is in  **
 **      EMM_COMMON_PROCEDURE_INITIATED state.                     **
 **                                                                        **
 **              3GPP TS 24.301, section 5.1.3.4.2                         **
 **                                                                        **
 ** Inputs:  evt:       The received EMM-SAP event                 **
 **      Others:    emm_fsm_status                             **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    emm_fsm_status                             **
 **                                                                        **
 ***************************************************************************/
int EmmCommonProcedureInitiated(emm_reg_t* const evt) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc                 = RETURNerror;
  emm_context_t* emm_ctx = evt->ctx;

  assert(emm_fsm_get_state(emm_ctx) == EMM_COMMON_PROCEDURE_INITIATED);

  switch (evt->primitive) {
    case _EMMREG_COMMON_PROC_REQ:
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-FSM state EMM_COMMON_PROCEDURE_INITIATED - Primitive "
          "_EMMREG_COMMON_PROC_REQ is not valid\n");
      break;

    case _EMMREG_COMMON_PROC_CNF:

      /*
       * An EMM common procedure successfully completed;
       */
      if (evt->u.common.common_proc) {
        if (evt->u.common.common_proc->emm_proc.base_proc.parent) {
          rc = nas_unlink_procedures(
              evt->u.common.common_proc->emm_proc.base_proc.parent,
              (nas_base_proc_t*) &evt->u.common.common_proc->emm_proc
                  .base_proc);
        }

        rc = emm_fsm_set_state(
            evt->ue_id, emm_ctx,
            ((nas_emm_proc_t*) evt->u.common.common_proc)
                ->previous_emm_fsm_state);

        if ((rc != RETURNerror) && (emm_ctx) && (evt->notify)) {
          (*evt->u.common.common_proc->emm_proc.base_proc.success_notif)(
              emm_ctx);
        }
        if (evt->free_proc) {
          nas_delete_common_procedure(emm_ctx, &evt->u.common.common_proc);
        }
      }
      break;

    case _EMMREG_COMMON_PROC_REJ:
      if (evt->u.common.common_proc) {
        /*
         * An EMM common procedure failed;
         * enter state EMM-DEREGISTERED.
         */
        rc = emm_fsm_set_state(evt->ue_id, emm_ctx, EMM_DEREGISTERED);

        if ((emm_ctx) &&
            (evt->u.common.common_proc->emm_proc.base_proc.fail_out)) {
          rc = (*evt->u.common.common_proc->emm_proc.base_proc.fail_out)(
              emm_ctx, &evt->u.common.common_proc->emm_proc.base_proc);
        }

        if ((rc != RETURNerror) && (emm_ctx) && (evt->notify) &&
            (evt->u.common.common_proc->emm_proc.base_proc.failure_notif)) {
          rc = (*evt->u.common.common_proc->emm_proc.base_proc.failure_notif)(
              emm_ctx);
        }

        if (evt->free_proc) {
          nas_delete_common_procedure(emm_ctx, &evt->u.common.common_proc);
        }
      }

      break;

    case _EMMREG_COMMON_PROC_ABORT:
      if (evt->u.common.common_proc) {
        if (evt->u.common.common_proc->emm_proc.base_proc.parent) {
          rc = nas_unlink_procedures(
              evt->u.common.common_proc->emm_proc.base_proc.parent,
              (nas_base_proc_t*) &evt->u.common.common_proc->emm_proc
                  .base_proc);
        }

        if ((emm_ctx) &&
            (evt->u.common.common_proc->emm_proc.base_proc.abort)) {
          (*evt->u.common.common_proc->emm_proc.base_proc.abort)(
              emm_ctx, &evt->u.common.common_proc->emm_proc.base_proc);
        }

        rc = emm_fsm_set_state(
            evt->ue_id, emm_ctx,
            ((nas_emm_proc_t*) evt->u.common.common_proc)
                ->previous_emm_fsm_state);

        if ((rc != RETURNerror) && (emm_ctx) && (evt->notify) &&
            (evt->u.common.common_proc->emm_proc.base_proc.failure_notif)) {
          (*evt->u.common.common_proc->emm_proc.base_proc.failure_notif)(
              emm_ctx);
        }
        if (evt->free_proc) {
          nas_delete_common_procedure(emm_ctx, &evt->u.common.common_proc);
        }
      }

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
      if (evt->u.attach.proc) {
        rc = emm_fsm_set_state(evt->ue_id, emm_ctx, EMM_DEREGISTERED);
        if ((emm_ctx) &&
            (evt->u.attach.proc->emm_spec_proc.emm_proc.base_proc.abort)) {
          (*evt->u.attach.proc->emm_spec_proc.emm_proc.base_proc.abort)(
              emm_ctx, &evt->u.attach.proc->emm_spec_proc.emm_proc.base_proc);
        }

        if ((rc != RETURNerror) && (emm_ctx) && (evt->notify) &&
            (evt->u.attach.proc->emm_spec_proc.emm_proc.base_proc
                 .failure_notif)) {
          (*evt->u.attach.proc->emm_spec_proc.emm_proc.base_proc.failure_notif)(
              emm_ctx);
        }
        if (evt->free_proc) {
          nas_delete_attach_procedure(emm_ctx);
        }
      }

      break;

    case _EMMREG_DETACH_INIT:
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-FSM state EMM_COMMON_PROCEDURE_INITIATED - Primitive "
          "_EMMREG_DETACH_INIT is not valid\n");
      break;

    case _EMMREG_DETACH_REQ:
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-FSM state EMM_COMMON_PROCEDURE_INITIATED - Primitive "
          "_EMMREG_DETACH_REQ is not valid\n");
      break;

    case _EMMREG_DETACH_FAILED:
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-FSM state EMM_COMMON_PROCEDURE_INITIATED - Primitive "
          "_EMMREG_DETACH_FAILED is not valid\n");
      break;

    case _EMMREG_DETACH_CNF:
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-FSM state EMM_COMMON_PROCEDURE_INITIATED - Primitive "
          "_EMMREG_DETACH_CNF is not valid\n");
      break;

    case _EMMREG_TAU_REQ:
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-FSM state EMM_COMMON_PROCEDURE_INITIATED - Primitive "
          "_EMMREG_TAU_REQ is not valid\n");
      break;

    case _EMMREG_TAU_CNF:
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-FSM state EMM_COMMON_PROCEDURE_INITIATED - Primitive "
          "_EMMREG_TAU_CNF is not valid\n");
      break;

    case _EMMREG_TAU_REJ:
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-FSM state EMM_COMMON_PROCEDURE_INITIATED - Primitive "
          "_EMMREG_TAU_REJ is not valid\n");
      break;

    case _EMMREG_SERVICE_REQ:
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-FSM state EMM_COMMON_PROCEDURE_INITIATED - Primitive "
          "_EMMREG_SERVICE_REQ is not valid\n");
      break;

    case _EMMREG_SERVICE_CNF:
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-FSM state EMM_COMMON_PROCEDURE_INITIATED - Primitive "
          "_EMMREG_SERVICE_CNF is not valid\n");
      break;

    case _EMMREG_SERVICE_REJ:
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-FSM state EMM_COMMON_PROCEDURE_INITIATED - Primitive "
          "_EMMREG_SERVICE_REJ is not valid\n");
      break;

    case _EMMREG_LOWERLAYER_SUCCESS:
      /*
       * Data successfully delivered to the network
       */
      if (emm_ctx) {
        nas_emm_proc_t* emm_proc = nas_emm_find_procedure_by_msg_digest(
            emm_ctx, (const char*) evt->u.ll_success.msg_digest,
            evt->u.ll_success.digest_len, evt->u.ll_success.msg_len);
        if (emm_proc) {
          if ((evt->notify) && (emm_proc->delivered)) {
            rc = (*emm_proc->delivered)(emm_ctx, emm_proc);
          }
        }
      }
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
        rc = emm_fsm_set_state(evt->ue_id, emm_ctx, EMM_DEREGISTERED);
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
          }
        }
        rc = emm_fsm_set_state(evt->ue_id, emm_ctx, EMM_DEREGISTERED);
      }
      break;

    default:
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-FSM state EMM_COMMON_PROCEDURE_INITIATED - Primitive is not "
          "valid "
          "(%d)\n",
          evt->primitive);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/
