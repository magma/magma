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
#ifndef FILE_MME_APP_PROCEDURES_SEEN
#define FILE_MME_APP_PROCEDURES_SEEN

/*! \file mme_app_procedures.h
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <stdint.h>

#include "queue.h"
#include "common_types.h"
#include "mme_app_ue_context.h"

// typedef int (*mme_app_pdu_in_resp_t)(void *arg);
// typedef int (*mme_app_pdu_in_rej_t)(void *arg);
// typedef int (*mme_app_time_out_t)(void *arg);
// typedef int (*mme_app_sdu_out_not_delivered_t)(void *arg);

typedef enum {
  MME_APP_BASE_PROC_TYPE_NONE = 0,
  MME_APP_BASE_PROC_TYPE_S1AP,
  MME_APP_BASE_PROC_TYPE_S11
} mme_app_base_proc_type_t;

typedef struct mme_app_base_proc_s {
  // PDU interface
  // pdu_in_resp_t              resp_in;
  // pdu_in_rej_t               fail_in;
  // time_out_t                 time_out;
  mme_app_base_proc_type_t type;
} mme_app_base_proc_t;

typedef enum {
  MME_APP_S11_PROC_TYPE_NONE = 0,
  MME_APP_S11_PROC_TYPE_CREATE_BEARER
} mme_app_s11_proc_type_t;

typedef struct mme_app_s11_proc_s {
  mme_app_base_proc_t proc;
  mme_app_s11_proc_type_t type;
  uintptr_t s11_trxn;
  LIST_ENTRY(mme_app_s11_proc_s) entries; /* List. */
} mme_app_s11_proc_t;

typedef enum {
  S11_PROC_BEARER_UNKNOWN = 0,
  S11_PROC_BEARER_PENDING = 1,
  S11_PROC_BEARER_FAILED  = 2,
  S11_PROC_BEARER_SUCCESS = 3
} s11_proc_bearer_status_t;

typedef struct mme_app_s11_proc_create_bearer_s {
  mme_app_s11_proc_t proc;
  int num_bearers;
  int num_status_received;
  // TODO here give a NAS/S1AP/.. reason -> GTPv2-C reason
  s11_proc_bearer_status_t bearer_status[BEARERS_PER_UE];
} mme_app_s11_proc_create_bearer_t;

typedef enum {
  MME_APP_S1AP_PROC_TYPE_NONE = 0,
  MME_APP_S1AP_PROC_TYPE_INITIAL
} mme_app_s1ap_proc_type_t;

typedef struct mme_app_s1ap_proc_s {
  mme_app_base_proc_t proc;
  mme_app_s1ap_proc_type_t type;
  LIST_ENTRY(mme_app_s1ap_proc_s) entries; /* List. */
} mme_app_s1ap_proc_t;

typedef struct mme_app_s11_proc_update_bearer_s {
  mme_app_s11_proc_t proc;
  int num_status_received;

  pdn_cid_t pci;
  pti_t pti;
  ambr_t new_used_ue_ambr;
  ambr_t apn_ambr;
  ebi_t linked_ebi;
  // TODO here give a NAS/S1AP/.. reason -> GTPv2-C reason
  bearer_contexts_to_be_updated_t*
      bcs_tbu; /**< Store the bearer contexts to be created here, and don't
                  register them yet in the MME_APP context. */
} mme_app_s11_proc_update_bearer_t;

typedef struct mme_app_s11_proc_delete_bearer_s {
  mme_app_s11_proc_t proc;
  ebi_t linked_ebi;
  int num_status_received;
  ebi_list_t ebis;

  // TODO here give a NAS/S1AP/.. reason -> GTPv2-C reason
  bearer_contexts_to_be_removed_t
      bcs_failed; /**< Store the bearer contexts to be created here, and don't
                     register them yet in the MME_APP context. */
} mme_app_s11_proc_delete_bearer_t;

typedef struct mme_app_s1ap_proc_e_rab_modify_bearer_ind_s {
  mme_app_s1ap_proc_t proc;
  mme_ue_s1ap_id_t mme_ue_s1ap_id;
  int num_status_received;
  e_rab_to_be_modified_bearer_mod_ind_list_t e_rab_to_be_modified_list;
  e_rab_not_to_be_modified_bearer_mod_ind_list_t e_rab_not_to_be_modified_list;

  e_rab_modify_bearer_mod_conf_list_t e_rab_modified_list;
  e_rab_list_t e_rab_failed_to_be_modified_list;
} mme_app_s1ap_proc_modify_bearer_ind_t;

void mme_app_delete_s11_procedures(ue_mm_context_t* const ue_context_p);
mme_app_s11_proc_create_bearer_t* mme_app_create_s11_procedure_create_bearer(
    ue_mm_context_t* const ue_context_p);
mme_app_s11_proc_create_bearer_t* mme_app_get_s11_procedure_create_bearer(
    ue_mm_context_t* const ue_context_p);
void mme_app_delete_s11_procedure_create_bearer(
    ue_mm_context_t* const ue_context_p);
void mme_app_s11_procedure_create_bearer_send_response(
    ue_mm_context_t* const ue_context_p,
    mme_app_s11_proc_create_bearer_t* s11_proc_create);

#endif
