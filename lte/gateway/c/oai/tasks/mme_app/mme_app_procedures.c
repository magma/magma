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

/*! \file mme_app_procedures.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#include <stdio.h>
#include <stdlib.h>

#include "dynamic_memory_check.h"
#include "assertions.h"
#include "common_types.h"
#include "intertask_interface.h"
#include "mme_app_ue_context.h"
#include "sgw_ie_defs.h"
#include "common_defs.h"
#include "mme_app_procedures.h"
#include "3gpp_24.007.h"
#include "3gpp_29.274.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "s11_messages_types.h"

static void mme_app_free_s11_procedure_create_bearer(
  mme_app_s11_proc_t **s11_proc);

//------------------------------------------------------------------------------
void mme_app_delete_s11_procedures(ue_mm_context_t *const ue_context_p)
{
  if (ue_context_p->s11_procedures) {
    mme_app_s11_proc_t *s11_proc1 = NULL;
    mme_app_s11_proc_t *s11_proc2 = NULL;

    s11_proc1 =
      LIST_FIRST(ue_context_p->s11_procedures); /* Faster List Deletion. */
    while (s11_proc1) {
      s11_proc2 = LIST_NEXT(s11_proc1, entries);
      if (MME_APP_S11_PROC_TYPE_CREATE_BEARER == s11_proc1->type) {
        mme_app_free_s11_procedure_create_bearer(&s11_proc1);
      } // else ...
      s11_proc1 = s11_proc2;
    }
    LIST_INIT(ue_context_p->s11_procedures);
    free_wrapper((void **) &ue_context_p->s11_procedures);
  }
}
//------------------------------------------------------------------------------
mme_app_s11_proc_create_bearer_t *mme_app_create_s11_procedure_create_bearer(
  ue_mm_context_t *const ue_context_p)
{
  mme_app_s11_proc_create_bearer_t *s11_proc_create_bearer =
    calloc(1, sizeof(mme_app_s11_proc_create_bearer_t));
  s11_proc_create_bearer->proc.proc.type = MME_APP_BASE_PROC_TYPE_S11;
  s11_proc_create_bearer->proc.type = MME_APP_S11_PROC_TYPE_CREATE_BEARER;
  mme_app_s11_proc_t *s11_proc = (mme_app_s11_proc_t *) s11_proc_create_bearer;

  if (!ue_context_p->s11_procedures) {
    ue_context_p->s11_procedures = calloc(1, sizeof(struct s11_procedures_s));
    LIST_INIT(ue_context_p->s11_procedures);
  }
  LIST_INSERT_HEAD((ue_context_p->s11_procedures), s11_proc, entries);
  return s11_proc_create_bearer;
}

//------------------------------------------------------------------------------
mme_app_s11_proc_create_bearer_t *mme_app_get_s11_procedure_create_bearer(
  ue_mm_context_t *const ue_context_p)
{
  if (ue_context_p->s11_procedures) {
    mme_app_s11_proc_t *s11_proc = NULL;

    LIST_FOREACH(s11_proc, ue_context_p->s11_procedures, entries)
    {
      if (MME_APP_S11_PROC_TYPE_CREATE_BEARER == s11_proc->type) {
        return (mme_app_s11_proc_create_bearer_t *) s11_proc;
      }
    }
  }
  return NULL;
}
//------------------------------------------------------------------------------
void mme_app_delete_s11_procedure_create_bearer(
  ue_mm_context_t *const ue_context_p)
{
  if (ue_context_p->s11_procedures) {
    mme_app_s11_proc_t *s11_proc = NULL;

    LIST_FOREACH(s11_proc, ue_context_p->s11_procedures, entries)
    {
      if (MME_APP_S11_PROC_TYPE_CREATE_BEARER == s11_proc->type) {
        LIST_REMOVE(s11_proc, entries);
        mme_app_free_s11_procedure_create_bearer(&s11_proc);
        return;
      }
    }
  }
}
//------------------------------------------------------------------------------
static void mme_app_free_s11_procedure_create_bearer(
  mme_app_s11_proc_t **s11_proc)
{
  // DO here specific releases (memory,etc)
  // nothing to do actually
  free_wrapper((void **) s11_proc);
}

//------------------------------------------------------------------------------
void mme_app_s11_procedure_create_bearer_send_response(
  ue_mm_context_t *const ue_context_p,
  mme_app_s11_proc_create_bearer_t *s11_proc_create)
{
  MessageDef *message_p =
    itti_alloc_new_message(TASK_MME_APP, S11_CREATE_BEARER_RESPONSE);
  AssertFatal(message_p, "itti_alloc_new_message Failed");

  message_p->ittiMsgHeader.imsi = ue_context_p->emm_context._imsi64;

  itti_s11_create_bearer_response_t *s11_create_bearer_response =
    &message_p->ittiMsg.s11_create_bearer_response;
  s11_create_bearer_response->local_teid = ue_context_p->mme_teid_s11;
  s11_create_bearer_response->trxn = (void *) s11_proc_create->proc.s11_trxn;
  s11_create_bearer_response->cause.cause_value = 0;
  int msg_bearer_index = 0;
  int num_rejected = 0;

  for (int ebix = 0; ebix < BEARERS_PER_UE; ebix++) {
    ebi_t ebi = INDEX_TO_EBI(ebix);
    if (S11_PROC_BEARER_FAILED == s11_proc_create->bearer_status[ebix]) {
      bearer_context_t *bc = mme_app_get_bearer_context(ue_context_p, ebi);
      // should not fail (bc != NULL)
      // Find remote S11 teid == find pdn
      s11_create_bearer_response->teid =
        ue_context_p->pdn_contexts[bc->pdn_cx_id]->s_gw_teid_s11_s4;

      s11_create_bearer_response->bearer_contexts
        .bearer_contexts[msg_bearer_index]
        .eps_bearer_id = ebi;
      s11_create_bearer_response->bearer_contexts
        .bearer_contexts[msg_bearer_index]
        .cause.cause_value = REQUEST_REJECTED;
      //  FTEID eNB
      s11_create_bearer_response->bearer_contexts
        .bearer_contexts[msg_bearer_index]
        .s1u_enb_fteid = bc->enb_fteid_s1u;
      // FTEID SGW S1U
      s11_create_bearer_response->bearer_contexts
        .bearer_contexts[msg_bearer_index]
        .s1u_sgw_fteid =
        bc->s_gw_fteid_s1u; ///< This IE shall be sent on the S11 interface. It shall be used
      s11_create_bearer_response->bearer_contexts.num_bearer_context++;
    } else if (
      S11_PROC_BEARER_SUCCESS == s11_proc_create->bearer_status[ebix]) {
      bearer_context_t *bc = mme_app_get_bearer_context(ue_context_p, ebi);
      // should not fail (bc != NULL)
      // Find remote S11 teid == find pdn
      s11_create_bearer_response->teid =
        ue_context_p->pdn_contexts[bc->pdn_cx_id]->s_gw_teid_s11_s4;

      s11_create_bearer_response->bearer_contexts
        .bearer_contexts[msg_bearer_index]
        .eps_bearer_id = ebi;
      s11_create_bearer_response->bearer_contexts
        .bearer_contexts[msg_bearer_index]
        .cause.cause_value = REQUEST_ACCEPTED;
      //  FTEID eNB
      s11_create_bearer_response->bearer_contexts
        .bearer_contexts[msg_bearer_index]
        .s1u_enb_fteid = bc->enb_fteid_s1u;
      // FTEID SGW S1U
      s11_create_bearer_response->bearer_contexts
        .bearer_contexts[msg_bearer_index]
        .s1u_sgw_fteid =
        bc->s_gw_fteid_s1u; ///< This IE shall be sent on the S11 interface. It shall be used
      s11_create_bearer_response->bearer_contexts.num_bearer_context++;
    }
  }
  if (s11_proc_create->num_bearers == num_rejected) {
    s11_create_bearer_response->cause.cause_value = REQUEST_REJECTED;
  } else if (num_rejected) {
    s11_create_bearer_response->cause.cause_value = REQUEST_ACCEPTED_PARTIALLY;
  } else {
    s11_create_bearer_response->cause.cause_value = REQUEST_ACCEPTED;
  }
  itti_send_msg_to_task(TASK_S11, INSTANCE_DEFAULT, message_p);
}
