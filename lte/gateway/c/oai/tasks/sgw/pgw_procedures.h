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
#ifndef FILE_PGW_PROCEDURES_SEEN
#define FILE_PGW_PROCEDURES_SEEN

/*! \file pgw_procedures.h
  \brief  Just a workaround waiting for PCEF implementation
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include "queue.h"
#include "3gpp_23.401.h"
#include "common_types.h"
#include "pgw_pcef_emulation.h"
#include "sgw_context_manager.h"

typedef enum {
  // should introduce Gx IP CAN procedures, etc, here
  PGW_BASE_PROC_TYPE_NONE = 0,
  PGW_BASE_PROC_TYPE_NETWORK_INITATED_CREATE_BEARER_REQUEST,
} pgw_base_proc_type_t;

typedef struct pgw_base_proc_s {
  //..
  LIST_ENTRY(pgw_base_proc_s) entries; /* List. */
  pgw_base_proc_type_t type;
} pgw_base_proc_t;

typedef struct sgw_eps_bearer_entry_wrapper_s {
  LIST_ENTRY(sgw_eps_bearer_entry_wrapper_s) entries; /* List. */
  sgw_eps_bearer_ctxt_t* sgw_eps_bearer_entry;
} sgw_eps_bearer_entry_wrapper_t;

typedef struct pgw_ni_cbr_proc_s {
  pgw_base_proc_t proc;
  teid_t teid;
  sdf_id_t sdf_id;
  // a list of sgw_eps_bearer_entry_t
  LIST_HEAD(pending_eps_bearers_s, sgw_eps_bearer_entry_wrapper_s) *
      pending_eps_bearers;
} pgw_ni_cbr_proc_t;

void pgw_delete_procedures(s_plus_p_gw_eps_bearer_context_information_t* ctx_p);
pgw_ni_cbr_proc_t* pgw_create_procedure_create_bearer(
    s_plus_p_gw_eps_bearer_context_information_t* ctx_p);
pgw_ni_cbr_proc_t* pgw_get_procedure_create_bearer(
    s_plus_p_gw_eps_bearer_context_information_t* ctx_p);
void pgw_delete_procedure_create_bearer(
    s_plus_p_gw_eps_bearer_context_information_t* ctx_p);
void pgw_free_procedure_create_bearer(pgw_ni_cbr_proc_t** ni_cbr_proc);

#endif
