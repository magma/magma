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

/*! \file mme_app_session_context.h
 *  \brief MME applicative layer
 *  \author Dincer Beken, Lionel Gauthier
 *  \date 2013
 *  \version 1.0
 *  \email: dbeken@blackned.de
 *  @defgroup _mme_app_impl_ MME applicative layer
 *  @ingroup _ref_implementation_
 *  @{
 */

#ifndef FILE_MME_APP_SESSION_CONTEXT_SEEN
#define FILE_MME_APP_SESSION_CONTEXT_SEEN
#include <inttypes.h> /* For sscanf formats */
#include <stdint.h>
#include <time.h> /* to provide time_t */

#include "bstrlib.h"
#include "common_types.h"
#include "esm_data.h"
#include "hashtable.h"
#include "mme_app_bearer_context.h"
#include "mme_app_messages_types.h"
#include "mme_app_procedures.h"
#include "obj_hashtable.h"
#include "queue.h"
#include "tree.h"

typedef int (*mme_app_ue_callback_t)(void*);

// TODO: (amar) only used in testing
#define IMSI_FORMAT "s"
#define IMSI_DATA(MME_APP_IMSI) (MME_APP_IMSI.data)

#define BEARER_STATE_NULL 0
#define BEARER_STATE_SGW_CREATED (1 << 0)
#define BEARER_STATE_MME_CREATED (1 << 1)
#define BEARER_STATE_ENB_CREATED (1 << 2)
#define BEARER_STATE_ACTIVE (1 << 3)
#define BEARER_STATE_S1_RELEASED (1 << 4)

#define MAX_NUM_BEARERS_UE 11 /**< Maximum number of bearers. */

/**
 * Bearer Pool elements.
 * Contains list bearer elements and PDN sessions.
 * A lock should be kept for all of it.
 */
typedef struct ue_session_pool_s {
  struct {
    pthread_mutex_t recmutex;  // mutex on the ue_context_t
    mme_ue_s1ap_id_t mme_ue_s1ap_id;
    /** Don't add them below, because they contain entries. */
    bearer_context_new_t bcs_ue[MAX_NUM_BEARERS_UE];
    pdn_context_t pdn_ue[MAX_APN_PER_UE];
    struct {
      /** Put field here only. */
      teid_t mme_teid_s11;
      teid_t saegw_teid_s11;
      int num_pdn_contexts;
      // Subscribed UE-AMBR: The Maximum Aggregated uplink and downlink MBR
      // values to be shared across all Non-GBR bearers according to the
      // subscription of the user. The used UE-AMBR will be calculated.
      ambr_t subscribed_ue_ambr;
      /** ESM Procedures : can be initialized together with the remaining
       * fields. */
      struct esm_procedures_s {
        LIST_HEAD(
            esm_pdn_connectivity_procedures_s,
            nas_esm_proc_pdn_connectivity_s) *
            pdn_connectivity_procedures;
        LIST_HEAD(
            esm_bearer_context_procedures_s, nas_esm_proc_bearer_context_s) *
            bearer_context_procedures;
      } esm_procedures;
      // todo: remove later
      ebi_t next_def_ebi_offset;
    } fields;
  } privates;

  /*
   * List of empty bearer context.
   * Take the bearer contexts from here and put them into the PDN context.
   */
  // LIST_HEAD(free_bearers_s, bearer_context_new_s) free_bearers;
  STAILQ_HEAD(free_bearers_s, bearer_context_new_s) free_bearers;
  /**
   * Map of PDN
   */
  RB_HEAD(PdnContexts, pdn_context_s) pdn_contexts;
  STAILQ_HEAD(free_pdn_s, pdn_context_s) free_pdn_contexts;
  LIST_HEAD(s1ap_procedures_s, mme_app_s1ap_proc_s) s1ap_procedures;

  /** Point to the next free session pool. */
  STAILQ_ENTRY(ue_session_pool_s) entries;
} ue_session_pool_t;

/* Declaration (prototype) of the function to store pdn and bearer contexts. */
RB_PROTOTYPE(
    PdnContexts, pdn_context_s, pdn_ctx_rbt_Node, mme_app_compare_pdn_context)

typedef struct mme_ue_session_pool_s {
  uint32_t nb_ue_session_pools_managed;
  uint32_t nb_ue_session_pools_since_last_stat;
  uint32_t nb_bearers_managed;
  uint32_t nb_bearers_since_last_stat;
  // todo: check if necessary hash_table_uint64_ts_t  *tun11_ue_context_htbl;//
  // data is mme_ue_s1ap_id_t
  hash_table_ts_t* mme_ue_s1ap_id_ue_session_pool_htbl;
  hash_table_uint64_ts_t*
      tun11_ue_session_pool_htbl;  // data is mme_ue_s1ap_id_t
} mme_ue_session_pool_t;

ue_session_pool_t* get_new_session_pool(mme_ue_s1ap_id_t ue_id);
void release_session_pool(ue_session_pool_t** ue_session_pool);

void mme_ue_session_pool_update_coll_keys(
    mme_ue_session_pool_t* const mme_ue_session_pool_p,
    ue_session_pool_t* const ue_session_pool,
    const mme_ue_s1ap_id_t mme_ue_s1ap_id, const s11_teid_t mme_teid_s11);

void mme_ue_session_pool_dump_coll_keys(void);

void mme_app_esm_detach(mme_ue_s1ap_id_t ue_id);
int mme_app_pdn_process_session_creation(
    mme_ue_s1ap_id_t ue_id, imsi64_t imsi, mm_state_t mm_state,
    ambr_t subscribed_ue_ambr, ebi_t default_ebi, fteid_t* saegw_s11_fteid,
    gtpv2c_cause_t* cause, bearer_contexts_created_t* bcs_created, ambr_t* ambr,
    paa_t** paa, protocol_configuration_options_t* pco);

/** \brief Retrieve an UE context by selecting the provided mme_ue_s1ap_id
 * \param mme_ue_s1ap_id The UE id identifier used in S1AP MME (and NAS)
 * @returns an UE context matching the mme_ue_s1ap_id or NULL if the context
 *doesn't exists
 **/
ue_session_pool_t* mme_ue_session_pool_exists_mme_ue_s1ap_id(
    mme_ue_session_pool_t* const mme_ue_context,
    const mme_ue_s1ap_id_t mme_ue_s1ap_id);
struct ue_session_pool_s* mme_ue_session_pool_exists_s11_teid(
    mme_ue_session_pool_t* const mme_ue_session_pool_p, const s11_teid_t teid);

void mme_app_ue_session_pool_s1_release_enb_informations(
    mme_ue_s1ap_id_t ue_id);

ambr_t mme_app_total_p_gw_apn_ambr(ue_session_pool_t* ue_session_pool);

ambr_t mme_app_total_p_gw_apn_ambr_rest(
    ue_session_pool_t* ue_session_pool, pdn_cid_t pci);

/** Create & deallocate a bearer context. Will also initialize the bearer
 * contexts. */
void clear_bearer_context(
    struct ue_session_pool_s* ue_session_pool, struct bearer_context_new_s* bc);

/** Find an allocated PDN session bearer context. */
struct bearer_context_new_s* mme_app_get_session_bearer_context(
    struct pdn_context_s* const pdn_context, const ebi_t ebi);

void mme_app_get_free_bearer_context(
    struct ue_session_pool_s* const ue_sp, const ebi_t ebi,
    struct bearer_context_new_s** bc_pp);

// todo_: combine these two methods
// void mme_app_get_session_bearer_context_from_all(
// struct ue_session_pool_s* const ue_session_pool, const ebi_t ebi,
// struct bearer_context_new_s** bc_pp);

/*
 * Receive Bearer Context VOs to send in CSR/Handover Request, etc..
 * Will set bearer state, unless it is null.
 */
void mme_app_get_bearer_contexts_to_be_created(
    struct pdn_context_s* pdn_context, bearer_contexts_to_be_created_t* bc_tbc,
    mme_app_bearer_state_t bc_state);

mme_app_s11_proc_t* mme_app_get_s11_procedure(
    struct ue_session_pool_s* const ue_session_pool);

mme_app_s11_proc_update_bearer_t* mme_app_create_s11_procedure_update_bearer(
    struct ue_session_pool_s* const ue_session_pool);
mme_app_s11_proc_update_bearer_t* mme_app_get_s11_procedure_update_bearer(
    struct ue_session_pool_s* const ue_session_pool);
void mme_app_delete_s11_procedure_update_bearer(
    struct ue_session_pool_s* const ue_session_pool);

mme_app_s11_proc_delete_bearer_t* mme_app_create_s11_procedure_delete_bearer(
    struct ue_session_pool_s* const ue_session_pool);
mme_app_s11_proc_delete_bearer_t* mme_app_get_s11_procedure_delete_bearer(
    struct ue_session_pool_s* const ue_session_pool);
void mme_app_delete_s11_procedure_delete_bearer(
    struct ue_session_pool_s* const ue_session_pool);

void mme_app_delete_s1ap_procedures(ue_session_pool_t* const ue_session_pool);
mme_app_s1ap_proc_modify_bearer_ind_t*
mme_app_create_s1ap_procedure_modify_bearer_ind(
    ue_session_pool_t* const ue_session_pool);
mme_app_s1ap_proc_modify_bearer_ind_t*
mme_app_get_s1ap_procedure_modify_bearer_ind(
    ue_session_pool_t* const ue_session_pool);
void mme_app_delete_s1ap_procedure_modify_bearer_ind(
    ue_session_pool_t* const ue_session_pool);

#endif /* FILE_MME_APP_UE_CONTEXT_SEEN */

/* @} */
