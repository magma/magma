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

/*! \file mme_app_session_context.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
  \author Dincer Beken
  \company Blackned GmbH
  \email: dbeken@blackned.de
*/
#include <pthread.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "bstrlib.h"

#include "assertions.h"
#include "common_defs.h"
#include "common_types.h"
#include "conversions.h"
#include "dynamic_memory_check.h"
#include "esm_cause.h"
#include "intertask_interface.h"
#include "log.h"
#include "mme_app_apn_selection.h"
#include "mme_app_bearer_context.h"
#include "mme_app_defs.h"
#include "mme_app_esm_procedures.h"
#include "mme_app_extern.h"
#include "mme_app_pdn_context.h"
#include "mme_app_session_context.h"
#include "mme_config.h"
#include "msc.h"

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/
static void clear_session_pool(ue_session_pool_t* ue_session_pool);
static int mme_insert_ue_session_pool(
    mme_ue_session_pool_t* const mme_ue_session_pool_p,
    const struct ue_session_pool_s* const ue_session_pool);

// todo: check the locks here
//------------------------------------------------------------------------------
void mme_ue_session_pool_update_coll_keys(
    mme_ue_session_pool_t* const mme_ue_session_pool_p,
    ue_session_pool_t* const ue_session_pool,
    const mme_ue_s1ap_id_t mme_ue_s1ap_id, const s11_teid_t mme_teid_s11) {
  hashtable_rc_t h_rc = HASH_TABLE_OK;
  void* id            = NULL;

  OAILOG_FUNC_IN(LOG_MME_APP);

  OAILOG_TRACE(
      LOG_MME_APP,
      "Update ue_session_pool.mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT
      " teid " TEID_FMT ". \n",
      ue_session_pool->privates.mme_ue_s1ap_id,
      ue_session_pool->privates.fields.mme_teid_s11);

  AssertFatal(
      (ue_session_pool->privates.mme_ue_s1ap_id == mme_ue_s1ap_id) &&
          (INVALID_MME_UE_S1AP_ID != mme_ue_s1ap_id),
      "Mismatch in UE session pool mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT
      "/" MME_UE_S1AP_ID_FMT "\n",
      ue_session_pool->privates.mme_ue_s1ap_id, mme_ue_s1ap_id);

  if ((INVALID_MME_UE_S1AP_ID != mme_ue_s1ap_id) &&
      (ue_session_pool->privates.mme_ue_s1ap_id != mme_ue_s1ap_id)) {
    // new insertion of mme_ue_s1ap_id, not a change in the id
    h_rc = hashtable_ts_remove(
        mme_ue_session_pool_p->mme_ue_s1ap_id_ue_session_pool_htbl,
        (const hash_key_t) ue_session_pool->privates.mme_ue_s1ap_id,
        (void**) &ue_session_pool);
    h_rc = hashtable_ts_insert(
        mme_ue_session_pool_p->mme_ue_s1ap_id_ue_session_pool_htbl,
        (const hash_key_t) mme_ue_s1ap_id, (void*) ue_session_pool);

    if (HASH_TABLE_OK != h_rc) {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Error could not update this ue session pool "
          "mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT " %s\n",
          ue_session_pool, ue_session_pool->privates.mme_ue_s1ap_id,
          hashtable_rc_code2string(h_rc));
    }
    ue_session_pool->privates.mme_ue_s1ap_id = mme_ue_s1ap_id;
  }
  /** S11 Key. */
  if ((ue_session_pool->privates.fields.mme_teid_s11 != mme_teid_s11) ||
      (ue_session_pool->privates.mme_ue_s1ap_id != mme_ue_s1ap_id)) {
    h_rc = hashtable_uint64_ts_remove(
        mme_ue_session_pool_p->tun11_ue_session_pool_htbl,
        (const hash_key_t) ue_session_pool->privates.fields.mme_teid_s11);
    if (INVALID_MME_UE_S1AP_ID != mme_ue_s1ap_id &&
        INVALID_TEID != mme_teid_s11) {
      h_rc = hashtable_uint64_ts_insert(
          mme_ue_session_pool_p->tun11_ue_session_pool_htbl,
          (const hash_key_t) mme_teid_s11, (void*) (uintptr_t) mme_ue_s1ap_id);
    } else {
      h_rc = HASH_TABLE_KEY_NOT_EXISTS;
    }

    if (HASH_TABLE_OK != h_rc && INVALID_TEID != mme_teid_s11) {
      OAILOG_TRACE(
          LOG_MME_APP,
          "Error could not update this ue session pool %p "
          "mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT " mme_s11_teid " TEID_FMT
          " : %s\n",
          ue_session_pool, ue_session_pool->privates.mme_ue_s1ap_id,
          mme_teid_s11, hashtable_rc_code2string(h_rc));
    }
    ue_session_pool->privates.fields.mme_teid_s11 = mme_teid_s11;
  }

  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
void mme_ue_session_pool_dump_coll_keys(void) {
  bstring tmp = bfromcstr(" ");
  btrunc(tmp, 0);

  hashtable_uint64_ts_dump_content(
      mme_app_desc.mme_ue_session_pools.tun11_ue_session_pool_htbl, tmp);
  OAILOG_TRACE(LOG_MME_APP, "tun11_ue_session_pool_htbl %s\n", bdata(tmp));

  btrunc(tmp, 0);
  hashtable_ts_dump_content(
      mme_app_desc.mme_ue_session_pools.mme_ue_s1ap_id_ue_session_pool_htbl,
      tmp);
  OAILOG_TRACE(
      LOG_MME_APP, "mme_ue_s1ap_id_ue_session_pool_htbl %s\n", bdata(tmp));
}

//------------------------------------------------------------------------------
static int mme_insert_ue_session_pool(
    mme_ue_session_pool_t* const mme_ue_session_pool_p,
    const struct ue_session_pool_s* const ue_session_pool) {
  hashtable_rc_t h_rc = HASH_TABLE_OK;

  OAILOG_FUNC_IN(LOG_MME_APP);
  DevAssert(mme_ue_session_pool_p);
  DevAssert(ue_session_pool);

  if (INVALID_MME_UE_S1AP_ID != ue_session_pool->privates.mme_ue_s1ap_id) {
    h_rc = hashtable_ts_is_key_exists(
        mme_ue_session_pool_p->mme_ue_s1ap_id_ue_session_pool_htbl,
        (const hash_key_t) ue_session_pool->privates.mme_ue_s1ap_id);

    if (HASH_TABLE_OK == h_rc) {
      OAILOG_DEBUG(
          LOG_MME_APP,
          "This ue session pool %p already exists "
          "mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT "\n",
          ue_session_pool, ue_session_pool->privates.mme_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
    }

    h_rc = hashtable_ts_insert(
        mme_ue_session_pool_p->mme_ue_s1ap_id_ue_session_pool_htbl,
        (const hash_key_t) ue_session_pool->privates.mme_ue_s1ap_id,
        (void*) ue_session_pool);

    if (HASH_TABLE_OK != h_rc) {
      OAILOG_DEBUG(
          LOG_MME_APP,
          "Error could not register this ue session pool %p "
          "mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT "\n",
          ue_session_pool, ue_session_pool->privates.mme_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
    }

    // filled S11 tun id
    if (ue_session_pool->privates.fields.mme_teid_s11) {
      h_rc = hashtable_uint64_ts_insert(
          mme_ue_session_pool_p->tun11_ue_session_pool_htbl,
          (const hash_key_t) ue_session_pool->privates.fields.mme_teid_s11,
          (void*) ((uintptr_t) ue_session_pool->privates.mme_ue_s1ap_id));

      if (HASH_TABLE_OK != h_rc) {
        OAILOG_DEBUG(
            LOG_MME_APP,
            "Error could not register this ue session pool %p "
            "mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT " mme_teid_s11 " TEID_FMT "\n",
            ue_session_pool, ue_session_pool->privates.mme_ue_s1ap_id,
            ue_session_pool->privates.fields.mme_teid_s11);
        OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
      }
    }
  }
  /*
   * Updating statistics
   */
  __sync_fetch_and_add(&mme_ue_session_pool_p->nb_ue_session_pools_managed, 1);
  __sync_fetch_and_add(
      &mme_ue_session_pool_p->nb_ue_session_pools_since_last_stat, 1);

  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

//------------------------------------------------------------------------------
void mme_remove_ue_session_pool(
    mme_ue_session_pool_t* const mme_ue_session_pool_p,
    struct ue_session_pool_s* ue_session_pool) {
  unsigned int* id       = NULL;
  hashtable_rc_t hash_rc = HASH_TABLE_OK;

  OAILOG_FUNC_IN(LOG_MME_APP);
  DevAssert(mme_ue_session_pool_p);
  DevAssert(ue_session_pool);

  // filled S11 tun id
  if (ue_session_pool->privates.fields.mme_teid_s11) {
    hash_rc = hashtable_uint64_ts_remove(
        mme_ue_session_pool_p->tun11_ue_session_pool_htbl,
        (const hash_key_t) ue_session_pool->privates.fields.mme_teid_s11);
    if (HASH_TABLE_OK != hash_rc)
      OAILOG_DEBUG(
          LOG_MME_APP,
          "UE session_pool mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT
          ", MME TEID_S11 " TEID_FMT "  not in S11 collection. \n",
          ue_session_pool->privates.mme_ue_s1ap_id,
          ue_session_pool->privates.fields.mme_teid_s11);
  }

  // filled NAS UE ID/ MME UE S1AP ID
  if (INVALID_MME_UE_S1AP_ID != ue_session_pool->privates.mme_ue_s1ap_id) {
    hash_rc = hashtable_ts_remove(
        mme_ue_session_pool_p->mme_ue_s1ap_id_ue_session_pool_htbl,
        (const hash_key_t) ue_session_pool->privates.mme_ue_s1ap_id,
        (void**) &ue_session_pool);
    if (HASH_TABLE_OK != hash_rc)
      OAILOG_DEBUG(
          LOG_MME_APP,
          "UE session_pool mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT
          " not in MME UE S1AP ID collection",
          ue_session_pool->privates.mme_ue_s1ap_id);
  }
  release_session_pool(&ue_session_pool);
  // todo: unlock?
  //  unlock_ue_contexts(ue_context);

  /*
   * Updating statistics
   */
  __sync_fetch_and_sub(&mme_ue_session_pool_p->nb_ue_session_pools_managed, 1);
  __sync_fetch_and_sub(
      &mme_ue_session_pool_p->nb_ue_session_pools_since_last_stat, 1);

  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
ue_session_pool_t* get_new_session_pool(mme_ue_s1ap_id_t ue_id) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  // todo: lock the mme_desc

  /** Check the first element in the list. If it is not empty, reject. */
  ue_session_pool_t* ue_session_pool =
      STAILQ_FIRST(&mme_app_desc.mme_ue_session_pools_list);
  DevAssert(ue_session_pool); /**< todo: with locks, it should be guaranteed,
                                 that this should exist. */
  if (ue_session_pool->privates.mme_ue_s1ap_id != INVALID_MME_UE_S1AP_ID) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "No free ue session pool left. Cannot allocate a new one.\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, NULL);
  }
  /** Found a free pool: Remove it from the head, add the ue_id and set it to
   * the end. */
  STAILQ_REMOVE_HEAD(
      &mme_app_desc.mme_ue_session_pools_list,
      entries); /**< free_sp is removed. */

  /** Initialize the bearers in the pool. */
  /** Remove the EMS-EBR context of the bearer-context. */
  OAILOG_INFO(
      LOG_MME_APP, "Clearing received current sp %p.\n", ue_session_pool);
  clear_session_pool(ue_session_pool);
  ue_session_pool->privates.mme_ue_s1ap_id = ue_id;
  /** Add it to the back of the list. */
  STAILQ_INSERT_TAIL(
      &mme_app_desc.mme_ue_session_pools_list, ue_session_pool, entries);
  DevAssert(
      mme_insert_ue_session_pool(
          &mme_app_desc.mme_ue_session_pools, ue_session_pool) == 0);
  // todo: unlock!
  OAILOG_FUNC_RETURN(LOG_MME_APP, ue_session_pool);
}

//------------------------------------------------------------------------------
void release_session_pool(ue_session_pool_t** ue_session_pool) {
  OAILOG_FUNC_IN(LOG_MME_APP);

  /** Clear the UE session pool. */
  OAILOG_INFO(
      LOG_MME_APP,
      "EMMCN-SAP  - "
      "Releasing session pool %p of UE " MME_UE_S1AP_ID_FMT ".\n",
      *ue_session_pool, (*ue_session_pool)->privates.mme_ue_s1ap_id);

  // todo: lock the mme_desc (ue_context should already be locked)
  /** Remove the ue_session pool from the list (probably at the back - must not
   * be at the very end. */
  STAILQ_REMOVE(
      &mme_app_desc.mme_ue_session_pools_list, (*ue_session_pool),
      ue_session_pool_s, entries);
  clear_session_pool(*ue_session_pool);
  /** Put it into the head. */
  STAILQ_INSERT_HEAD(
      &mme_app_desc.mme_ue_session_pools_list, (*ue_session_pool), entries);
  *ue_session_pool = NULL;
  // todo: unlock the mme_desc
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
void mme_app_esm_detach(mme_ue_s1ap_id_t ue_id) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  ue_session_pool_t* ue_session_pool =
      mme_ue_session_pool_exists_mme_ue_s1ap_id(
          &mme_app_desc.mme_ue_session_pools, ue_id);
  pdn_context_t* pdn_context = NULL;

  if (!ue_session_pool) {
    OAILOG_WARNING(
        LOG_MME_APP,
        "No UE session_pool could be found for UE: " MME_UE_S1AP_ID_FMT
        " to release ESM contexts. \n",
        ue_id);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  // todo: LOCK UE SESSION POOL
  mme_app_nas_esm_free_bearer_context_procedures(ue_session_pool);
  mme_app_nas_esm_free_pdn_connectivity_procedures(ue_session_pool);

  OAILOG_INFO(
      LOG_MME_APP,
      "Removed all ESM procedures of UE: " MME_UE_S1AP_ID_FMT " for detach. \n",
      ue_id);
  pdn_context = RB_MIN(PdnContexts, &ue_session_pool->pdn_contexts);
  while (pdn_context) {
    mme_app_delete_pdn_context(ue_session_pool, &pdn_context);
    pdn_context = RB_MIN(PdnContexts, &ue_session_pool->pdn_contexts);
  }
  OAILOG_INFO(
      LOG_MME_APP,
      "Removed all ESM contexts of UE: " MME_UE_S1AP_ID_FMT
      " for detach. Removing the session pool. \n",
      ue_id);
  // todo: UNLOCK UE SESSION POOL

  /** Release the session pool. */  // todo: locks
  mme_remove_ue_session_pool(
      &mme_app_desc.mme_ue_session_pools, ue_session_pool);
  OAILOG_INFO(
      LOG_MME_APP,
      "Removing the session pool for UE " MME_UE_S1AP_ID_FMT ". \n", ue_id);
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
int mme_app_pdn_process_session_creation(
    mme_ue_s1ap_id_t ue_id, imsi64_t imsi, mm_state_t mm_state,
    ambr_t subscribed_ue_ambr, ebi_t default_ebi, fteid_t* saegw_s11_fteid,
    gtpv2c_cause_t* cause, bearer_contexts_created_t* bcs_created, ambr_t* ambr,
    paa_t** paa, protocol_configuration_options_t* pco) {
  OAILOG_FUNC_IN(LOG_MME_APP);

  pdn_context_t /* * pdn_context1 = NULL, */* pdn_context = NULL;
  ue_session_pool_t* ue_session_pool =
      mme_ue_session_pool_exists_mme_ue_s1ap_id(
          &mme_app_desc.mme_ue_session_pools, ue_id);
  if (!ue_session_pool) {
    OAILOG_WARNING(
        LOG_MME_APP,
        "No MME_APP UE session pool could be found for UE: " MME_UE_S1AP_ID_FMT
        " to process CSResp. \n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  //  /** Get the first unestablished PDN context from the UE context. */
  //  RB_FOREACH (pdn_context1, PdnContexts, &ue_session_pool->pdn_contexts) {
  //	for(int num_ebi = 0; num_ebi < bcs_created->num_bearer_context;
  // num_ebi++){ 	    if(pdn_context1->default_ebi ==
  // bcs_created->bearer_context[num_ebi].eps_bearer_id){
  //	      /** Found. */
  //	      pdn_context = pdn_context1;
  //	      break;
  //	    }
  //	}
  //	if(pdn_context)
  //		break;
  //  }

  /** Get the first unestablished PDN context from the UE context. */
  RB_FOREACH(pdn_context, PdnContexts, &ue_session_pool->pdn_contexts) {
    if (!pdn_context->s_gw_teid_s11_s4) {
      /** Found. */
      break;
    }
  }
  if (!pdn_context || pdn_context->s_gw_teid_s11_s4) {
    OAILOG_WARNING(
        LOG_MME_APP,
        "No unestablished PDN context could be found for "
        "UE: " MME_UE_S1AP_ID_FMT ". \n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  // LOCK_UE_CONTEXT
  /** Set the S11 FTEID for each PDN connection. */
  if (saegw_s11_fteid->teid)
    pdn_context->s_gw_teid_s11_s4 = saegw_s11_fteid->teid;
  if (!ue_session_pool->privates.fields.saegw_teid_s11)
    ue_session_pool->privates.fields.saegw_teid_s11 =
        pdn_context->s_gw_teid_s11_s4;
  if (pco) {
    if (!pdn_context->pco) {
      pdn_context->pco = calloc(1, sizeof(protocol_configuration_options_t));
    } else {
      clear_protocol_configuration_options(pdn_context->pco);
    }
    copy_protocol_configuration_options(pdn_context->pco, pco);
  }
  if (saegw_s11_fteid->ipv4) {
    ((struct sockaddr_in*) &pdn_context->s_gw_addr_s11_s4)->sin_addr.s_addr =
        saegw_s11_fteid->ipv4_address.s_addr;
    ((struct sockaddr_in*) &pdn_context->s_gw_addr_s11_s4)->sin_family =
        AF_INET;
  } else {
    ((struct sockaddr_in6*) &pdn_context->s_gw_addr_s11_s4)->sin6_family =
        AF_INET6;
    memcpy(
        &((struct sockaddr_in6*) &pdn_context->s_gw_addr_s11_s4)->sin6_addr,
        &saegw_s11_fteid->ipv6_address, sizeof(saegw_s11_fteid->ipv6_address));
  }

  /** Check the received cause. */
  if (cause->cause_value != REQUEST_ACCEPTED &&
      cause->cause_value != REQUEST_ACCEPTED_PARTIALLY) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Received S11_CREATE_SESSION_RESPONSE REJECTION with cause "
        "value %d for ue " MME_UE_S1AP_ID_FMT "from S+P-GW. \n",
        cause->cause_value, ue_id);
    /** Destroy the PDN context (not in the ESM layer). */
    mme_app_esm_delete_pdn_context(
        ue_id, pdn_context->apn_subscribed, pdn_context->context_identifier,
        pdn_context->default_ebi); /**< Frees it & puts session bearers back to
                                      the pool. */
    // todo: UNLOCK
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
  }
  /** Process the success case, no bearer at this point. */
  if (*paa) {
    /** Set the PAA. */
    if (pdn_context->paa) {
      free_wrapper((void**) &pdn_context->paa);
    }
    pdn_context->paa = *paa;
    /** Decouple it from the message. */
    *paa = NULL;
  }

  //#define TEMPORARY_DEBUG 1
  //#if TEMPORARY_DEBUG
  // bstring b =
  // protocol_configuration_options_to_xml(&ue_context->pending_pdn_connectivity_req_pco);
  // OAILOG_DEBUG (LOG_MME_APP, "PCO %s\n", bdata(b));
  // bdestroy_wrapper(&b);
  //#endif

  /** Check the received APN-AMBR is in bounds (a subscription profile may exist
   * or not, but these values must exist). */
  DevAssert(subscribed_ue_ambr.br_dl && subscribed_ue_ambr.br_ul);
  if (ambr->br_dl && ambr->br_ul) { /**< New APN-AMBR received. */
    ambr_t current_apn_ambr = mme_app_total_p_gw_apn_ambr_rest(
        ue_session_pool, pdn_context->context_identifier);
    ambr_t total_apn_ambr;
    total_apn_ambr.br_dl = current_apn_ambr.br_dl + ambr->br_dl;
    total_apn_ambr.br_ul = current_apn_ambr.br_ul + ambr->br_ul;
    /** Actualized, used total APN-AMBR. */
    if (total_apn_ambr.br_dl > subscribed_ue_ambr.br_dl ||
        total_apn_ambr.br_ul > subscribed_ue_ambr.br_ul) {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Received new APN_AMBR for PDN \"%s\" (ctx_id=%d) for "
          "UE " MME_UE_S1AP_ID_FMT
          " exceeds the subscribed ue ambr (br_dl=%d,br_ul=%d). \n",
          bdata(pdn_context->apn_subscribed), pdn_context->context_identifier,
          ue_id, subscribed_ue_ambr.br_dl, subscribed_ue_ambr.br_ul);
      /** Will use the current APN bitrates, either from the APN configuration
       * received from the HSS or from S10. */
      // todo: inform PCRF, that default APN bitrates could not be updated.
      DevAssert(
          pdn_context->subscribed_apn_ambr.br_dl &&
          pdn_context->subscribed_apn_ambr.br_ul);
      if (mm_state == UE_REGISTERED) {
        /**
         * This case is only if it is a multi-APN, including the handover case.
         * Check if there is any AMBR left in the UE-AMBR, so the remaining and
         * inform the PGW about it. If not reject the request.
         */
        /** Check if any PDN connectivity procedure is running, if not reject
         * the request. */
        nas_esm_proc_pdn_connectivity_t* esm_proc_pdn_connectivity =
            mme_app_nas_esm_get_pdn_connectivity_procedure(
                ue_id, PROCEDURE_TRANSACTION_IDENTITY_UNASSIGNED);
        if (esm_proc_pdn_connectivity == NULL) {
          OAILOG_ERROR(
              LOG_MME_APP,
              "For APN \"%s\" (ctx_id=%d) for UE " MME_UE_S1AP_ID_FMT
              ", no PDN connectivity procedure is running.\n",
              bdata(pdn_context->apn_subscribed),
              pdn_context->context_identifier, ue_id);
          OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
        }
        if ((current_apn_ambr.br_dl >= subscribed_ue_ambr.br_dl) ||
            (current_apn_ambr.br_ul >= subscribed_ue_ambr.br_ul)) {
          OAILOG_ERROR(
              LOG_MME_APP,
              "No more AMBR left for UE " MME_UE_S1AP_ID_FMT
              ". Rejecting the request for an additional PDN "
              "(apn_subscribed=\"%s\", ctx_id=%d). \n",
              ue_id, bdata(pdn_context->apn_subscribed),
              pdn_context->context_identifier);
          cause->cause_value =
              NO_RESOURCES_AVAILABLE; /**< Reject the pdn connection in
                                         particular, remove the PDN context. */
          //    		  mme_app_esm_delete_pdn_context(ue_id,
          //    pdn_context->apn_subscribed, pdn_context->context_identifier,
          //    pdn_context->default_ebi); /**< Frees it & puts session bearers
          //    back to the pool. */
          OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
        }
        pdn_context->subscribed_apn_ambr.br_dl =
            subscribed_ue_ambr.br_dl - current_apn_ambr.br_dl;
        pdn_context->subscribed_apn_ambr.br_ul =
            subscribed_ue_ambr.br_ul - current_apn_ambr.br_ul;
        OAILOG_WARNING(
            LOG_MME_APP,
            "For the UE " MME_UE_S1AP_ID_FMT
            ", enforcing the remaining AMBR (br_dl=%d,br_ul=%d) on the "
            "additionally requested PDN (apn_subscribed=\"%s\", ctx_id=%d). \n",
            ue_id, pdn_context->subscribed_apn_ambr.br_dl,
            pdn_context->subscribed_apn_ambr.br_ul,
            bdata(pdn_context->apn_subscribed),
            pdn_context->context_identifier);
        /** Continue to process it. */
        esm_proc_pdn_connectivity->saegw_qos_modification = true;
      } else {
        /**
         * If it is a handover, or S10 triggered idle TAU, there should be no
         * subscription data at this point. So reject it.
         *
         * Else, it is either initial attach or an initial tau without the S10
         * procedure. In that case the PDN should also be the default PDN. Use
         * the UE-AMBR for the PDN AMBR. Mark it as changed to update the AMBR
         * in the PGW/PCRF.
         */
        subscription_data_t* subscription_data =
            mme_ue_subscription_data_exists_imsi(
                &mme_app_desc.mme_ue_contexts, imsi);
        if (subscription_data &&
            !(current_apn_ambr.br_dl | current_apn_ambr.br_ul)) {
          nas_esm_proc_pdn_connectivity_t* esm_proc_pdn_connectivity =
              mme_app_nas_esm_get_pdn_connectivity_procedure(
                  ue_id, PROCEDURE_TRANSACTION_IDENTITY_UNASSIGNED);
          if (esm_proc_pdn_connectivity == NULL) {
            OAILOG_ERROR(
                LOG_MME_APP,
                "For APN \"%s\" (ctx_id=%d) for UE " MME_UE_S1AP_ID_FMT
                ", no PDN connectivity procedure is running (initial "
                "attach/tau).\n",
                bdata(pdn_context->apn_subscribed),
                pdn_context->context_identifier, ue_id);
            OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
          }
          /* Put the remaining PDN AMBR as the UE AMBR. There should be no other
           * PDN context allocated (no remaining AMBR). */
          pdn_context->subscribed_apn_ambr.br_dl =
              subscription_data->subscribed_ambr.br_dl;
          pdn_context->subscribed_apn_ambr.br_ul =
              subscription_data->subscribed_ambr.br_ul;
          /** Mark the procedure as modified. */
          esm_proc_pdn_connectivity->saegw_qos_modification = true;
        } else {
          /**
           * No subscription data exists. Assuming a handover/S10 TAU procedure
           * (even if a (decoupled) subscription data exists. Rejecting the
           * handover or the S10 part of the idle TAU (continue with
           * initial-TAU).
           */
          if (pdn_context->subscribed_apn_ambr.br_dl &&
              pdn_context->subscribed_apn_ambr.br_ul) {
            OAILOG_WARNING(
                LOG_MME_APP,
                "Using the PDN-AMBR (br_dl=%d, br_ul=%d) from handover for APN "
                "(apn_subscribed=\"%s\", ctx_id=%d) exceeds UE AMBR for "
                "UE " MME_UE_S1AP_ID_FMT
                "."
                "Rejecting the PDN connectivity for the inter-MME mobility "
                "procedure. \n",
                pdn_context->subscribed_apn_ambr.br_dl,
                pdn_context->subscribed_apn_ambr.br_ul,
                bdata(pdn_context->apn_subscribed),
                pdn_context->context_identifier, ue_id);
            // todo: no ESM proc exists.
          } else {
            OAILOG_ERROR(
                LOG_MME_APP,
                "Requested PDN (apn_subscribed=\"%s\", ctx_id=%d) "
                "exceeds UE AMBR for UE " MME_UE_S1AP_ID_FMT
                " with IMSI " IMSI_64_FMT
                ". "
                "Rejecting the PDN connectivity for the inter-MME "
                "mobility procedure. \n",
                bdata(pdn_context->apn_subscribed),
                pdn_context->context_identifier, ue_id, imsi);
            cause->cause_value =
                NO_RESOURCES_AVAILABLE; /**< Reject the request for this PDN
                                           connectivity procedure. */
            mme_app_esm_delete_pdn_context(
                ue_id, pdn_context->apn_subscribed,
                pdn_context->context_identifier,
                pdn_context->default_ebi); /**< Frees it & puts session bearers
                                              back to the pool. */
            OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
          }
        }
      }
    } else {
      OAILOG_DEBUG(
          LOG_MME_APP,
          "Received new valid APN_AMBR for APN \"%s\" (ctx_id=%d) for "
          "UE " MME_UE_S1AP_ID_FMT ". Updating APN ambr. \n",
          bdata(pdn_context->apn_subscribed), pdn_context->context_identifier,
          ue_id);
      pdn_context->subscribed_apn_ambr.br_dl = ambr->br_dl;
      pdn_context->subscribed_apn_ambr.br_ul = ambr->br_ul;
    }
  }

  /** Updated all PDN (bearer generic) information. Traverse all bearers,
   * including the default bearer. */
  for (int i = 0; i < bcs_created->num_bearer_context; i++) {
    ebi_t bearer_id = bcs_created->bearer_context[i].eps_bearer_id;
    bearer_context_created_t* bc_created = &bcs_created->bearer_context[i];
    bearer_context_new_t* bearer_context =
        mme_app_get_session_bearer_context(pdn_context, bearer_id);
    /*
     * Depending on s11 result we have to send reject or accept for bearers
     */
    DevCheck(
        (bearer_id < BEARERS_PER_UE + 5) && (bearer_id >= 5), bearer_id,
        BEARERS_PER_UE, 0);
    DevAssert(
        bcs_created->bearer_context[i].s1u_sgw_fteid.interface_type ==
        S1_U_SGW_GTP_U);
    /** Check if the bearer could be created in the PGW. */
    if (bc_created->cause.cause_value != REQUEST_ACCEPTED) {
      /** Check that it is a dedicated bearer. */
      DevAssert(pdn_context->default_ebi != bearer_id);
      /*
       * Release all session bearers of the PDN context back into the UE pool.
       */
      if (bearer_context) {
        /** Initialize the new bearer context. */
        STAILQ_REMOVE(
            &pdn_context->session_bearers, bearer_context, bearer_context_new_s,
            entries);
        clear_bearer_context(ue_session_pool, bearer_context);
        OAILOG_WARNING(
            LOG_MME_APP,
            "Successfully deregistered the bearer context (ebi=%d) "
            "from PDN \"%s\" and for ue_id " MME_UE_S1AP_ID_FMT "\n",
            bearer_id, bdata(pdn_context->apn_subscribed), ue_id);
      }
      continue;
    }
    /** Bearer context could be established successfully. */
    if (bearer_context) {
      bearer_context->bearer_state |= BEARER_STATE_SGW_CREATED;
      /** No context identifier might be set yet (multiple might exist and all
       * might be 0. */
      //      AssertFatal((pdn_cx_id >= 0) && (pdn_cx_id < MAX_APN_PER_UE), "Bad
      //      pdn id for bearer");
      /*
       * Updating statistics
       */
      mme_app_desc.mme_ue_session_pools.nb_bearers_managed++;
      mme_app_desc.mme_ue_session_pools.nb_bearers_since_last_stat++;
      /** Update the FTEIDs of the SAE-GW. */
      memcpy(
          &bearer_context->s_gw_fteid_s1u,
          &bcs_created->bearer_context[i].s1u_sgw_fteid,
          sizeof(fteid_t)); /**< Also copying the IPv4/V6 address. */
      memcpy(
          &bearer_context->p_gw_fteid_s5_s8_up,
          &bcs_created->bearer_context[i].s5_s8_u_pgw_fteid, sizeof(fteid_t));
      /** Check if the bearer level QoS parameters have been modified by the
       * PGW. */
      if (bcs_created->bearer_context[i].bearer_level_qos.qci &&
          bcs_created->bearer_context[i].bearer_level_qos.pl) {
        /**
         * We set them here, since we may not have a NAS context in (S10)
         * mobility. We don't check the subscribed HSS values. The PGW may ask
         * for more.
         */
        bearer_context->bearer_level_qos.qci =
            bcs_created->bearer_context[i].bearer_level_qos.qci;
        bearer_context->bearer_level_qos.pl =
            bcs_created->bearer_context[i].bearer_level_qos.pl;
        bearer_context->bearer_level_qos.pvi =
            bcs_created->bearer_context[i].bearer_level_qos.pvi;
        bearer_context->bearer_level_qos.pci =
            bcs_created->bearer_context[i].bearer_level_qos.pci;
        OAILOG_DEBUG(
            LOG_MME_APP, "Set qci %u in bearer %u\n",
            bearer_context->bearer_level_qos.qci, bearer_id);
      }
      /** Not touching TFT, ESM-EBR-State here. */
    } else {
      DevMessage(
          "Bearer context that could be established successfully in the SAE-GW "
          "could not be found in the MME session bearers."); /**< This should
                                                                not happen,
                                                                since we lock
                                                                the
                                                                pdn_context.. */
                                                             //        continue;
    }
  }
  /** Assert that the default bearer at least exists. */
  DevAssert(mme_app_get_session_bearer_context(
      pdn_context, pdn_context->default_ebi));
  // todo: UNLOCK_UE_SESSION_POOL
  OAILOG_INFO(
      LOG_MME_APP,
      "Processed all %d bearer contexts for APN \"%s\" for "
      "ue_id " MME_UE_S1AP_ID_FMT ". \n",
      bcs_created->num_bearer_context, bdata(pdn_context->apn_subscribed),
      ue_id);
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

//------------------------------------------------------------------------------
ambr_t mme_app_total_p_gw_apn_ambr(ue_session_pool_t* ue_session_pool) {
  pdn_context_t* registered_pdn_ctx = NULL;
  ambr_t apn_ambr_sum               = {0, 0};
  RB_FOREACH(registered_pdn_ctx, PdnContexts, &ue_session_pool->pdn_contexts) {
    DevAssert(registered_pdn_ctx);
    apn_ambr_sum.br_dl += registered_pdn_ctx->subscribed_apn_ambr.br_dl;
    apn_ambr_sum.br_ul += registered_pdn_ctx->subscribed_apn_ambr.br_ul;
  }
  return apn_ambr_sum;
}

//------------------------------------------------------------------------------
ambr_t mme_app_total_p_gw_apn_ambr_rest(
    ue_session_pool_t* ue_session_pool, pdn_cid_t pci) {
  /** Get the total APN AMBR excluding the given PCI. */
  pdn_context_t* registered_pdn_ctx = NULL;
  ambr_t apn_ambr_sum               = {0, 0};
  RB_FOREACH(registered_pdn_ctx, PdnContexts, &ue_session_pool->pdn_contexts) {
    DevAssert(registered_pdn_ctx);
    if (registered_pdn_ctx->context_identifier == pci) continue;
    apn_ambr_sum.br_dl += registered_pdn_ctx->subscribed_apn_ambr.br_dl;
    apn_ambr_sum.br_ul += registered_pdn_ctx->subscribed_apn_ambr.br_ul;
  }
  return apn_ambr_sum;
}

//------------------------------------------------------------------------------
void mme_app_ue_session_pool_s1_release_enb_informations(
    mme_ue_s1ap_id_t ue_id) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  ue_session_pool_t* ue_session_pool =
      mme_ue_session_pool_exists_mme_ue_s1ap_id(
          &mme_app_desc.mme_ue_session_pools, ue_id);
  if (!ue_session_pool) {
    OAILOG_INFO(
        LOG_MME_APP,
        "No session pool exists for UE " MME_UE_S1AP_ID_FMT
        ". Cannot release bearer s1 info. \n",
        ue_id);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  // LOCK UE_SESSION
  pdn_context_t* registered_pdn_ctx = NULL;
  /** Update all bearers and get the pdn context id. */
  RB_FOREACH(registered_pdn_ctx, PdnContexts, &ue_session_pool->pdn_contexts) {
    DevAssert(registered_pdn_ctx);
    /*
     * Get the first PDN whose bearers are not established yet.
     * Do the MBR just one PDN at a time.
     */
    bearer_context_new_t *bearer_context_to_set_idle      = NULL,
                         *bearer_context_to_set_idle_safe = NULL;
    STAILQ_FOREACH(
        bearer_context_to_set_idle, &registered_pdn_ctx->session_bearers,
        entries) {
      DevAssert(bearer_context_to_set_idle);
      /** Add them to the bearears list of the MBR. */
      mme_app_bearer_context_s1_release_enb_informations(
          bearer_context_to_set_idle);
    }
  }

  // UNLOCK UE_SESSION
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
ue_session_pool_t* mme_ue_session_pool_exists_mme_ue_s1ap_id(
    mme_ue_session_pool_t* const mme_ue_session_pool_p,
    const mme_ue_s1ap_id_t mme_ue_s1ap_id) {
  struct ue_session_pool_s* ue_session_pool = NULL;

  hashtable_ts_get(
      mme_ue_session_pool_p->mme_ue_s1ap_id_ue_session_pool_htbl,
      (const hash_key_t) mme_ue_s1ap_id, (void**) &ue_session_pool);
  if (ue_session_pool) {
    //    lock_ue_contexts(ue_context);
    //    OAILOG_TRACE (LOG_MME_APP, "UE  " MME_UE_S1AP_ID_FMT " fetched MM
    //    state %s, ECM state %s\n ",mme_ue_s1ap_id,
    //        (ue_context->privates.fields.mm_state == UE_UNREGISTERED) ?
    //        "UE_UNREGISTERED":(ue_context->privates.fields.mm_state ==
    //        UE_REGISTERED) ? "UE_REGISTERED":"UNKNOWN",
    //        (ue_context->privates.fields.ecm_state == ECM_IDLE) ?
    //        "ECM_IDLE":(ue_context->privates.fields.ecm_state ==
    //        ECM_CONNECTED) ? "ECM_CONNECTED":"UNKNOWN");
  }
  return ue_session_pool;
}

//------------------------------------------------------------------------------
struct ue_session_pool_s* mme_ue_session_pool_exists_s11_teid(
    mme_ue_session_pool_t* const mme_ue_session_pool_p, const s11_teid_t teid) {
  hashtable_rc_t h_rc       = HASH_TABLE_OK;
  uint64_t mme_ue_s1ap_id64 = 0;

  h_rc = hashtable_uint64_ts_get(
      mme_ue_session_pool_p->tun11_ue_session_pool_htbl,
      (const hash_key_t) teid, &mme_ue_s1ap_id64);

  if (HASH_TABLE_OK == h_rc) {
    return mme_ue_session_pool_exists_mme_ue_s1ap_id(
        mme_ue_session_pool_p, (mme_ue_s1ap_id_t) mme_ue_s1ap_id64);
  }
  return NULL;
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/

//------------------------------------------------------------------------------
static void clear_session_pool(ue_session_pool_t* ue_session_pool) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  mme_ue_s1ap_id_t ue_id = ue_session_pool->privates.mme_ue_s1ap_id;
  ue_session_pool->privates.mme_ue_s1ap_id = INVALID_MME_UE_S1AP_ID;

  OAILOG_INFO(
      LOG_MME_APP, "Clearing UE session pool of UE " MME_UE_S1AP_ID_FMT ". \n",
      ue_id);
  /** Release the procedures. */
  mme_app_delete_s11_procedures(ue_session_pool);

  /** Free the ESM procedures. */
  if (ue_session_pool->privates.fields.esm_procedures
          .pdn_connectivity_procedures) {
    OAILOG_WARNING(
        LOG_MME_APP,
        "ESM PDN Connectivity procedures still exist for UE " MME_UE_S1AP_ID_FMT
        ". \n",
        ue_id);
    mme_app_nas_esm_free_pdn_connectivity_procedures(ue_session_pool);
  }

  if (ue_session_pool->privates.fields.esm_procedures
          .bearer_context_procedures) {
    OAILOG_WARNING(
        LOG_MME_APP,
        "ESM Bearer Context procedures still exist for UE " MME_UE_S1AP_ID_FMT
        ". \n",
        ue_id);
    mme_app_nas_esm_free_bearer_context_procedures(ue_session_pool);
  }

  DevAssert(RB_EMPTY(&ue_session_pool->pdn_contexts));

  /** No need to check if list is full, since it is stacked. Just clear the
   * allocations in each session context. */
  // LIST_INIT(&ue_session_pool->free_bearers);
  STAILQ_INIT(&ue_session_pool->free_bearers);
  STAILQ_INIT(&ue_session_pool->free_pdn_contexts);

  /** Initialize the bearer contexts. */
  for (int num_bearer = 0; num_bearer < MAX_NUM_BEARERS_UE; num_bearer++) {
    /** Set the EBI, if nothing set. */
    if (!ue_session_pool->privates.bcs_ue[num_bearer].ebi)
      ue_session_pool->privates.bcs_ue[num_bearer].ebi = num_bearer + 5;
    OAILOG_TRACE(
        LOG_MME_APP,
        "Clearing bearer context with ebi %d (%p) for UE " MME_UE_S1AP_ID_FMT
        ". \n",
        ue_session_pool->privates.bcs_ue[num_bearer].ebi,
        &ue_session_pool->privates.bcs_ue[num_bearer], ue_id);
    clear_bearer_context(
        ue_session_pool, &ue_session_pool->privates.bcs_ue[num_bearer]);
  }

  /** Initialize the PDN contexts. */
  for (int num_pdn = 0; num_pdn < MAX_APN_PER_UE; num_pdn++) {
    memset(
        &ue_session_pool->privates.pdn_ue[num_pdn], 0,
        sizeof(
            pdn_context_t)); /**< Sets the SM status to ESM_STATUS_INVALID. */
    /** Insert the list into the empty list. */
    STAILQ_INSERT_TAIL(
        &ue_session_pool->free_pdn_contexts,
        &ue_session_pool->privates.pdn_ue[num_pdn], entries);
  }

  /** Initialize the RB_MAP. */
  RB_INIT(&ue_session_pool->pdn_contexts);
  /** Re-initialize the empty list: it should not point to a next free variable:
   * determined when it is put back into the list (like LIST_INIT). */
  memset(
      &ue_session_pool->privates.fields, 0,
      sizeof(ue_session_pool->privates.fields));
  OAILOG_INFO(
      LOG_MME_APP,
      "Successfully cleared UE session pool of UE " MME_UE_S1AP_ID_FMT ". \n",
      ue_id);
  OAILOG_FUNC_OUT(LOG_MME_APP);
}
