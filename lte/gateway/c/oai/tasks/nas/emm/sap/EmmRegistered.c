/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under 
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.  
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
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
  Source      EmmRegistered.c

  Version     0.1

  Date        2012/10/03

  Product     NAS stack

  Subsystem   EPS Mobility Management

  Author      Frederic Maurel, Lionel GAUTHIER

  Description Implements the EPS Mobility Management procedures executed
        when the EMM-SAP is in EMM-REGISTERED state.

        In EMM-REGISTERED state, an EMM context has been established
        and a default EPS bearer context has been activated in the UE
        and the MME.
        The UE may initiate sending and receiving user data and signal-
        ling information and reply to paging. Additionally, tracking
        area updating or combined tracking area updating procedure is
        performed.

*****************************************************************************/
#include <pthread.h>
#include <inttypes.h>
#include <stdint.h>
#include <stdbool.h>
#include <string.h>
#include <stdlib.h>
#include <arpa/inet.h>
#include <assert.h>

#include "bstrlib.h"

#include "log.h"
#include "common_defs.h"
#include "emm_fsm.h"
#include "commonDef.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"
#include "3gpp_29.274.h"
#include "networkDef.h"
#include "emm_proc.h"

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
int EmmRegistered(const emm_reg_t *evt)
{
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;
  emm_context_t *emm_ctx = evt->ctx;

  assert(emm_fsm_get_state(emm_ctx) == EMM_REGISTERED);

  switch (evt->primitive) {
    case _EMMREG_COMMON_PROC_REQ:
      MSC_LOG_RX_MESSAGE(
        MSC_NAS_EMM_MME,
        MSC_NAS_EMM_MME,
        NULL,
        0,
        "_EMMREG_COMMON_PROC_REQ ue id " MME_UE_S1AP_ID_FMT " ",
        evt->ue_id);
      /*
     * An EMM common procedure has been initiated;
     * enter state EMM-COMMON-PROCEDURE-INITIATED.
     */
      rc =
        emm_fsm_set_state(evt->ue_id, evt->ctx, EMM_COMMON_PROCEDURE_INITIATED);
      break;

    case _EMMREG_COMMON_PROC_CNF:
      OAILOG_ERROR(
        LOG_NAS_EMM,
        "EMM-FSM state EMM_REGISTERED - Primitive _EMMREG_COMMON_PROC_CNF is "
        "not valid\n");
      MSC_LOG_RX_DISCARDED_MESSAGE(
        MSC_NAS_EMM_MME,
        MSC_NAS_EMM_MME,
        NULL,
        0,
        "_EMMREG_COMMON_PROC_CNF ue id " MME_UE_S1AP_ID_FMT " ",
        evt->ue_id);
      break;

    case _EMMREG_COMMON_PROC_REJ:
      OAILOG_ERROR(
        LOG_NAS_EMM,
        "EMM-FSM state EMM_REGISTERED - Primitive _EMMREG_COMMON_PROC_REJ is "
        "not valid\n");
      MSC_LOG_RX_DISCARDED_MESSAGE(
        MSC_NAS_EMM_MME,
        MSC_NAS_EMM_MME,
        NULL,
        0,
        "_EMMREG_COMMON_PROC_REJ ue id " MME_UE_S1AP_ID_FMT " ",
        evt->ue_id);
      break;

    case _EMMREG_COMMON_PROC_ABORT:
      OAILOG_ERROR(
        LOG_NAS_EMM,
        "EMM-FSM state EMM_REGISTERED - Primitive _EMMREG_COMMON_PROC_ABORT is "
        "not valid\n");
      MSC_LOG_RX_DISCARDED_MESSAGE(
        MSC_NAS_EMM_MME,
        MSC_NAS_EMM_MME,
        NULL,
        0,
        "_EMMREG_COMMON_PROC_ABORT ue id " MME_UE_S1AP_ID_FMT " ",
        evt->ue_id);
      rc = emm_fsm_set_state(evt->ue_id, evt->ctx, EMM_DEREGISTERED);
      break;

    case _EMMREG_ATTACH_CNF:
      MSC_LOG_RX_MESSAGE(
        MSC_NAS_EMM_MME,
        MSC_NAS_EMM_MME,
        NULL,
        0,
        "_EMMREG_ATTACH_CNF ue id " MME_UE_S1AP_ID_FMT " ",
        evt->ue_id);
      /*
     * Attach procedure successful and default EPS bearer
     * context activated;
     * enter state EMM-REGISTERED.
     */
      rc = emm_fsm_set_state(evt->ue_id, emm_ctx, EMM_REGISTERED);
      break;

    case _EMMREG_ATTACH_REJ:
      MSC_LOG_RX_MESSAGE(
        MSC_NAS_EMM_MME,
        MSC_NAS_EMM_MME,
        NULL,
        0,
        "_EMMREG_ATTACH_REJ ue id " MME_UE_S1AP_ID_FMT " ",
        evt->ue_id);
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
        "EMM-FSM state EMM_REGISTERED - Primitive _EMMREG_ATTACH_ABORT is not "
        "valid\n");
      MSC_LOG_RX_DISCARDED_MESSAGE(
        MSC_NAS_EMM_MME,
        MSC_NAS_EMM_MME,
        NULL,
        0,
        "_EMMREG_ATTACH_ABORT ue id " MME_UE_S1AP_ID_FMT " ",
        evt->ue_id);
      break;

    case _EMMREG_DETACH_INIT:
      OAILOG_ERROR(
        LOG_NAS_EMM,
        "EMM-FSM state EMM_REGISTERED - Primitive _EMMREG_DETACH_INIT is not "
        "valid\n");
      MSC_LOG_RX_DISCARDED_MESSAGE(
        MSC_NAS_EMM_MME,
        MSC_NAS_EMM_MME,
        NULL,
        0,
        "_EMMREG_DETACH_INIT ue id " MME_UE_S1AP_ID_FMT " ",
        evt->ue_id);
      break;

    case _EMMREG_DETACH_REQ:
      MSC_LOG_RX_MESSAGE(
        MSC_NAS_EMM_MME,
        MSC_NAS_EMM_MME,
        NULL,
        0,
        "_EMMREG_DETACH_REQ ue id " MME_UE_S1AP_ID_FMT " ",
        evt->ue_id);
      rc = emm_fsm_set_state(evt->ue_id, emm_ctx, EMM_DEREGISTERED_INITIATED);
      break;

    case _EMMREG_DETACH_FAILED:
      OAILOG_ERROR(
        LOG_NAS_EMM,
        "EMM-FSM state EMM_REGISTERED - Primitive _EMMREG_DETACH_FAILED is not "
        "valid\n");
      MSC_LOG_RX_DISCARDED_MESSAGE(
        MSC_NAS_EMM_MME,
        MSC_NAS_EMM_MME,
        NULL,
        0,
        "_EMMREG_DETACH_FAILED ue id " MME_UE_S1AP_ID_FMT " ",
        evt->ue_id);
      break;

    case _EMMREG_DETACH_CNF:
      OAILOG_ERROR(
        LOG_NAS_EMM,
        "EMM-FSM state EMM_REGISTERED - Primitive _EMMREG_DETACH_CNF is not "
        "valid\n");
      MSC_LOG_RX_DISCARDED_MESSAGE(
        MSC_NAS_EMM_MME,
        MSC_NAS_EMM_MME,
        NULL,
        0,
        "_EMMREG_DETACH_CNF ue id " MME_UE_S1AP_ID_FMT " ",
        evt->ue_id);
      break;

    case _EMMREG_TAU_REQ:
      OAILOG_ERROR(
        LOG_NAS_EMM,
        "EMM-FSM state EMM_REGISTERED - Primitive _EMMREG_TAU_REQ is not "
        "valid\n");
      MSC_LOG_RX_DISCARDED_MESSAGE(
        MSC_NAS_EMM_MME,
        MSC_NAS_EMM_MME,
        NULL,
        0,
        "_EMMREG_TAU_REQ ue id " MME_UE_S1AP_ID_FMT " ",
        evt->ue_id);
      break;

    case _EMMREG_TAU_CNF:
      if (evt->u.tau.proc) {
        MSC_LOG_RX_MESSAGE(
          MSC_NAS_EMM_MME,
          MSC_NAS_EMM_MME,
          NULL,
          0,
          "_EMMREG_TAU_CNF ue id " MME_UE_S1AP_ID_FMT " ",
          evt->ue_id);
        if (
          (emm_ctx) &&
          (evt->u.tau.proc->emm_spec_proc.emm_proc.base_proc.fail_out)) {
          rc = (*evt->u.tau.proc->emm_spec_proc.emm_proc.base_proc.fail_out)(
            emm_ctx, &evt->u.tau.proc->emm_spec_proc.emm_proc.base_proc);
        }

        if (
          (rc != RETURNerror) && (emm_ctx) && (evt->notify) &&
          (evt->u.tau.proc->emm_spec_proc.emm_proc.base_proc.failure_notif)) {
          rc =
            (*evt->u.tau.proc->emm_spec_proc.emm_proc.base_proc.failure_notif)(
              emm_ctx);
        }

        if (evt->free_proc) {
          nas_delete_tau_procedure(emm_ctx);
        }
      }
      break;
    case _EMMREG_TAU_REJ:
      if (evt->u.tau.proc) {
        MSC_LOG_RX_MESSAGE(
          MSC_NAS_EMM_MME,
          MSC_NAS_EMM_MME,
          NULL,
          0,
          "_EMMREG_TAU_REJ ue id " MME_UE_S1AP_ID_FMT " ",
          evt->ue_id);
        rc = emm_fsm_set_state(evt->ue_id, evt->ctx, EMM_DEREGISTERED);

        if (
          (emm_ctx) &&
          (evt->u.tau.proc->emm_spec_proc.emm_proc.base_proc.fail_out)) {
          rc = (*evt->u.tau.proc->emm_spec_proc.emm_proc.base_proc.fail_out)(
            emm_ctx, &evt->u.tau.proc->emm_spec_proc.emm_proc.base_proc);
        }

        if (
          (rc != RETURNerror) && (emm_ctx) && (evt->notify) &&
          (evt->u.tau.proc->emm_spec_proc.emm_proc.base_proc.failure_notif)) {
          rc =
            (*evt->u.tau.proc->emm_spec_proc.emm_proc.base_proc.failure_notif)(
              emm_ctx);
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
      MSC_LOG_RX_DISCARDED_MESSAGE(
        MSC_NAS_EMM_MME,
        MSC_NAS_EMM_MME,
        NULL,
        0,
        "_EMMREG_SERVICE_REQ ue id " MME_UE_S1AP_ID_FMT " ",
        evt->ue_id);
      break;

    case _EMMREG_SERVICE_CNF:
      OAILOG_ERROR(
        LOG_NAS_EMM,
        "EMM-FSM state EMM_REGISTERED - Primitive _EMMREG_SERVICE_CNF is not "
        "valid\n");
      MSC_LOG_RX_DISCARDED_MESSAGE(
        MSC_NAS_EMM_MME,
        MSC_NAS_EMM_MME,
        NULL,
        0,
        "_EMMREG_SERVICE_CNF ue id " MME_UE_S1AP_ID_FMT " ",
        evt->ue_id);
      break;

    case _EMMREG_SERVICE_REJ:
      OAILOG_ERROR(
        LOG_NAS_EMM,
        "EMM-FSM state EMM_REGISTERED - Primitive _EMMREG_SERVICE_REJ is not "
        "valid\n");
      MSC_LOG_RX_DISCARDED_MESSAGE(
        MSC_NAS_EMM_MME,
        MSC_NAS_EMM_MME,
        NULL,
        0,
        "_EMMREG_SERVICE_REJ ue id " MME_UE_S1AP_ID_FMT " ",
        evt->ue_id);
      break;

    case _EMMREG_LOWERLAYER_SUCCESS:
      MSC_LOG_RX_MESSAGE(
        MSC_NAS_EMM_MME,
        MSC_NAS_EMM_MME,
        NULL,
        0,
        "_EMMREG_LOWERLAYER_SUCCESS ue id " MME_UE_S1AP_ID_FMT " ",
        evt->ue_id);
      /*
     * Data successfully delivered to the network
     */
      rc = RETURNok;
      break;

    case _EMMREG_LOWERLAYER_FAILURE:
      MSC_LOG_RX_MESSAGE(
        MSC_NAS_EMM_MME,
        MSC_NAS_EMM_MME,
        NULL,
        0,
        "_EMMREG_LOWERLAYER_FAILURE ue id " MME_UE_S1AP_ID_FMT " ",
        evt->ue_id);

      if (emm_ctx) {
        nas_emm_proc_t *emm_proc = nas_emm_find_procedure_by_msg_digest(
          emm_ctx,
          (const char *) evt->u.ll_failure.msg_digest,
          evt->u.ll_failure.digest_len,
          evt->u.ll_failure.msg_len);
        if (emm_proc) {
          if ((evt->notify) && (emm_proc->not_delivered)) {
            rc = (*emm_proc->not_delivered)(emm_ctx, emm_proc);
          }
        }
      }
      break;

    case _EMMREG_LOWERLAYER_RELEASE:
      MSC_LOG_RX_MESSAGE(
        MSC_NAS_EMM_MME,
        MSC_NAS_EMM_MME,
        NULL,
        0,
        "_EMMREG_LOWERLAYER_RELEASE ue id " MME_UE_S1AP_ID_FMT " ",
        evt->ue_id);
      nas_delete_all_emm_procedures(emm_ctx);
      rc = RETURNok;
      break;

    case _EMMREG_LOWERLAYER_NON_DELIVERY:
      MSC_LOG_RX_MESSAGE(
        MSC_NAS_EMM_MME,
        MSC_NAS_EMM_MME,
        NULL,
        0,
        "_EMMREG_LOWERLAYER_NON_DELIVERY ue id " MME_UE_S1AP_ID_FMT " ",
        evt->ue_id);
      if (emm_ctx) {
        nas_emm_proc_t *emm_proc = nas_emm_find_procedure_by_msg_digest(
          emm_ctx,
          (const char *) evt->u.non_delivery_ho.msg_digest,
          evt->u.non_delivery_ho.digest_len,
          evt->u.non_delivery_ho.msg_len);
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
      MSC_LOG_RX_DISCARDED_MESSAGE(
        MSC_NAS_EMM_MME,
        MSC_NAS_EMM_MME,
        NULL,
        0,
        "_EMMREG_UNKNOWN(primitive id %d) ue id " MME_UE_S1AP_ID_FMT " ",
        evt->primitive,
        evt->ue_id);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/
