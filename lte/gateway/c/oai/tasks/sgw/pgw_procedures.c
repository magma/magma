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

/*! \file pgw_procedures.c
  \brief  Just a workaround waiting for PCEF implementation
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#include <stdio.h>
#include <stdlib.h>

#include "dynamic_memory_check.h"
#include "sgw_context_manager.h"
#include "pgw_procedures.h"

//------------------------------------------------------------------------------
void pgw_delete_procedures(
    s_plus_p_gw_eps_bearer_context_information_t* const ctx_p) {
  if (ctx_p->sgw_eps_bearer_context_information.pending_procedures) {
    pgw_base_proc_t* base_proc1 = NULL;
    pgw_base_proc_t* base_proc2 = NULL;

    base_proc1 =
        LIST_FIRST(ctx_p->sgw_eps_bearer_context_information
                       .pending_procedures); /* Faster List Deletion. */
    while (base_proc1) {
      base_proc2 = LIST_NEXT(base_proc1, entries);
      if (PGW_BASE_PROC_TYPE_NETWORK_INITATED_CREATE_BEARER_REQUEST ==
          base_proc1->type) {
        pgw_free_procedure_create_bearer((pgw_ni_cbr_proc_t**) &base_proc1);
      }  // else ...
      base_proc1 = base_proc2;
    }
    LIST_INIT(ctx_p->sgw_eps_bearer_context_information.pending_procedures);
    free_wrapper(
        (void**) &ctx_p->sgw_eps_bearer_context_information.pending_procedures);
  }
}
//------------------------------------------------------------------------------
pgw_ni_cbr_proc_t* pgw_create_procedure_create_bearer(
    s_plus_p_gw_eps_bearer_context_information_t* const ctx_p) {
  pgw_ni_cbr_proc_t* s11_proc_create_bearer =
      calloc(1, sizeof(pgw_ni_cbr_proc_t));
  s11_proc_create_bearer->proc.type =
      PGW_BASE_PROC_TYPE_NETWORK_INITATED_CREATE_BEARER_REQUEST;
  pgw_base_proc_t* base_proc = (pgw_base_proc_t*) s11_proc_create_bearer;

  if (!ctx_p->sgw_eps_bearer_context_information.pending_procedures) {
    ctx_p->sgw_eps_bearer_context_information.pending_procedures =
        calloc(1, sizeof(struct pending_eps_bearers_s));
    LIST_INIT(ctx_p->sgw_eps_bearer_context_information.pending_procedures);
  }
  LIST_INSERT_HEAD(
      (ctx_p->sgw_eps_bearer_context_information.pending_procedures), base_proc,
      entries);

  s11_proc_create_bearer->pending_eps_bearers =
      calloc(1, sizeof(struct pending_eps_bearers_s));
  LIST_INIT(s11_proc_create_bearer->pending_eps_bearers);

  return s11_proc_create_bearer;
}

//------------------------------------------------------------------------------
pgw_ni_cbr_proc_t* pgw_get_procedure_create_bearer(
    s_plus_p_gw_eps_bearer_context_information_t* const ctx_p) {
  if (ctx_p->sgw_eps_bearer_context_information.pending_procedures) {
    pgw_base_proc_t* base_proc = NULL;

    LIST_FOREACH(
        base_proc, ctx_p->sgw_eps_bearer_context_information.pending_procedures,
        entries) {
      if (PGW_BASE_PROC_TYPE_NETWORK_INITATED_CREATE_BEARER_REQUEST ==
          base_proc->type) {
        return (pgw_ni_cbr_proc_t*) base_proc;
      }
    }
  }
  return NULL;
}
//------------------------------------------------------------------------------
void pgw_delete_procedure_create_bearer(
    s_plus_p_gw_eps_bearer_context_information_t* const ctx_p) {
  if (ctx_p->sgw_eps_bearer_context_information.pending_procedures) {
    pgw_base_proc_t* base_proc = NULL;

    LIST_FOREACH(
        base_proc, ctx_p->sgw_eps_bearer_context_information.pending_procedures,
        entries) {
      if (PGW_BASE_PROC_TYPE_NETWORK_INITATED_CREATE_BEARER_REQUEST ==
          base_proc->type) {
        LIST_REMOVE(base_proc, entries);
        pgw_free_procedure_create_bearer((pgw_ni_cbr_proc_t**) &base_proc);
        return;
      }
    }
  }
}
//------------------------------------------------------------------------------
void pgw_free_procedure_create_bearer(pgw_ni_cbr_proc_t** ni_cbr_proc) {
  if (*ni_cbr_proc && (*ni_cbr_proc)->pending_eps_bearers) {
    sgw_eps_bearer_entry_wrapper_t* eps_bearer_entry_wrapper = NULL;
    LIST_FOREACH(
        eps_bearer_entry_wrapper, (*ni_cbr_proc)->pending_eps_bearers,
        entries) {
      if (eps_bearer_entry_wrapper) {
        LIST_REMOVE(eps_bearer_entry_wrapper, entries);
        free_wrapper((void**) &eps_bearer_entry_wrapper->sgw_eps_bearer_entry);
        free_wrapper((void**) &eps_bearer_entry_wrapper);
        if (LIST_EMPTY((*ni_cbr_proc)->pending_eps_bearers)) {
          free_wrapper((void**) &(*ni_cbr_proc)->pending_eps_bearers);
          break;
        }
      }
    }
  }
  free_wrapper((void**) ni_cbr_proc);
}
