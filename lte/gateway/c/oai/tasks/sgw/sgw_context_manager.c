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

/*! \file sgw_context_manager.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#define SGW
#define SGW_CONTEXT_MANAGER_C

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>
#include <inttypes.h>

#include "bstrlib.h"
#include "dynamic_memory_check.h"
#include "assertions.h"
#include "hashtable.h"
#include "obj_hashtable.h"
#include "common_defs.h"
#include "log.h"
#include "sgw_context_manager.h"

//-----------------------------------------------------------------------------
static bool sgw_display_s11teid2mme_mapping(
  uint64_t keyP,
  void *dataP,
  void *unused_parameterP,
  void **unused_resultP)
//-----------------------------------------------------------------------------
{
  mme_sgw_tunnel_t *mme_sgw_tunnel = NULL;

  if (dataP) {
    mme_sgw_tunnel = (mme_sgw_tunnel_t *) dataP;
    OAILOG_DEBUG(
      LOG_SPGW_APP,
      "| " TEID_FMT "\t<------------->\t" TEID_FMT "\n",
      mme_sgw_tunnel->remote_teid,
      mme_sgw_tunnel->local_teid);
  } else {
    OAILOG_DEBUG(LOG_SPGW_APP, "INVALID S11 TEID MAPPING FOUND\n");
  }
  return false;
}

//-----------------------------------------------------------------------------
void sgw_display_s11teid2mme_mappings(spgw_state_t *state)
//-----------------------------------------------------------------------------
{
  OAILOG_DEBUG(LOG_SPGW_APP, "+--------------------------------------+\n");
  OAILOG_DEBUG(LOG_SPGW_APP, "| MME <--- S11 TE ID MAPPINGS ---> SGW |\n");
  OAILOG_DEBUG(LOG_SPGW_APP, "+--------------------------------------+\n");
  hashtable_ts_apply_callback_on_elements(
    state->sgw_state.s11teid2mme, sgw_display_s11teid2mme_mapping, NULL, NULL);
  OAILOG_DEBUG(LOG_SPGW_APP, "+--------------------------------------+\n");
}

//-----------------------------------------------------------------------------
void sgw_display_sgw_eps_bearer_context(
  const sgw_eps_bearer_ctxt_t *const eps_bearer_ctxt)
//-----------------------------------------------------------------------------
{
  if (eps_bearer_ctxt) {
    OAILOG_DEBUG(
      LOG_SPGW_APP,
      "|\t\t\t\tebi: %u, enb_teid_for_S1u: " TEID_FMT
      ", s_gw_teid_for_S1u_S12_S4_up: " TEID_FMT " (tbc)\n",
      eps_bearer_ctxt->eps_bearer_id,
      eps_bearer_ctxt->enb_teid_S1u,
      eps_bearer_ctxt->s_gw_teid_S1u_S12_S4_up);
  }
}

//-----------------------------------------------------------------------------
static bool sgw_display_s11_bearer_context_information(
  uint64_t keyP,
  void *dataP,
  void *unused_parameterP,
  void **unused_resultP)
//-----------------------------------------------------------------------------
{
  s_plus_p_gw_eps_bearer_context_information_t *sp_context_information = NULL;

  if (dataP) {
    sp_context_information =
      (s_plus_p_gw_eps_bearer_context_information_t *) dataP;
    OAILOG_DEBUG(LOG_SPGW_APP, "| KEY %" PRId64 ":      \n", keyP);
    OAILOG_DEBUG(
      LOG_SPGW_APP, "|\tsgw_eps_bearer_context_information:     |\n");
    //Imsi_t               imsi;                           ///< IMSI (International Mobile Subscriber Identity) is the subscriber permanent identity.
    OAILOG_DEBUG(
      LOG_SPGW_APP,
      "|\t\timsi_unauthenticated_indicator:\t%u\n",
      sp_context_information->sgw_eps_bearer_context_information
        .imsi_unauthenticated_indicator);
    //char                 msisdn[MSISDN_LENGTH];          ///< The basic MSISDN of the UE. The presence is dictated by its storage in the HSS.
    OAILOG_DEBUG(
      LOG_SPGW_APP,
      "|\t\tmme_teid_    S11:              \t" TEID_FMT "\n",
      sp_context_information->sgw_eps_bearer_context_information.mme_teid_S11);
    //ip_address_t         mme_ip_address_for_S11;         ///< MME IP address the S11 interface.
    OAILOG_DEBUG(
      LOG_SPGW_APP,
      "|\t\ts_gw_teid_S11_S4:              \t" TEID_FMT "\n",
      sp_context_information->sgw_eps_bearer_context_information
        .s_gw_teid_S11_S4);
    //ip_address_t         s_gw_ip_address_for_S11_S4;     ///< S-GW IP address for the S11 interface and the S4 Interface (control plane).
    //cgi_t                last_known_cell_Id;             ///< This is the last location of the UE known by the network
    OAILOG_DEBUG(LOG_SPGW_APP, "|\t\tpdn_connection:\n");
    OAILOG_DEBUG(
      LOG_SPGW_APP,
      "|\t\t\tapn_in_use:        %s\n",
      sp_context_information->sgw_eps_bearer_context_information.pdn_connection
        .apn_in_use);
    OAILOG_DEBUG(
      LOG_SPGW_APP,
      "|\t\t\tdefault_bearer:    %u\n",
      sp_context_information->sgw_eps_bearer_context_information.pdn_connection
        .default_bearer);
    OAILOG_DEBUG(LOG_SPGW_APP, "|\t\t\teps_bearers:\n");
    for (int ebix = 0; ebix < BEARERS_PER_UE; ebix++) {
      sgw_display_sgw_eps_bearer_context(
        sp_context_information->sgw_eps_bearer_context_information
          .pdn_connection.sgw_eps_bearers_array[ebix]);
    }
    //void                  *trxn;
    //uint32_t               peer_ip;
  } else {
    OAILOG_DEBUG(
      LOG_SPGW_APP,
      "INVALID s_plus_p_gw_eps_bearer_context_information FOUND\n");
  }
  return false;
}

//-----------------------------------------------------------------------------
void sgw_display_s11_bearer_context_information_mapping(spgw_state_t *state)
//-----------------------------------------------------------------------------
{
  OAILOG_DEBUG(LOG_SPGW_APP, "+-----------------------------------------+\n");
  OAILOG_DEBUG(LOG_SPGW_APP, "| S11 BEARER CONTEXT INFORMATION MAPPINGS |\n");
  OAILOG_DEBUG(LOG_SPGW_APP, "+-----------------------------------------+\n");
  hashtable_ts_apply_callback_on_elements(
    state->sgw_state.s11_bearer_context_information,
    sgw_display_s11_bearer_context_information,
    NULL,
    NULL);
  OAILOG_DEBUG(LOG_SPGW_APP, "+--------------------------------------+\n");
}

//-----------------------------------------------------------------------------
void pgw_lite_cm_free_apn(pgw_apn_t **apnP)
//-----------------------------------------------------------------------------
{
  if (*apnP) {
    if ((*apnP)->pdn_connections) {
      obj_hashtable_ts_destroy((*apnP)->pdn_connections);
    }
  }
}

//-----------------------------------------------------------------------------
teid_t sgw_get_new_S11_tunnel_id(spgw_state_t *state)
//-----------------------------------------------------------------------------
{
  return ++state->sgw_state.tunnel_id;
}

//-----------------------------------------------------------------------------
mme_sgw_tunnel_t *sgw_cm_create_s11_tunnel(
  spgw_state_t *state,
  teid_t remote_teid,
  teid_t local_teid)
//-----------------------------------------------------------------------------
{
  mme_sgw_tunnel_t *new_tunnel = NULL;

  new_tunnel = calloc(1, sizeof(mme_sgw_tunnel_t));

  if (new_tunnel == NULL) {
    /*
     * Malloc failed, may be ENOMEM error
     */
    OAILOG_ERROR(
      LOG_SPGW_APP,
      "Failed to create tunnel for remote_teid " TEID_FMT "\n",
      remote_teid);
    return NULL;
  }

  new_tunnel->remote_teid = remote_teid;
  new_tunnel->local_teid = local_teid;
  /*
   * Trying to insert the new tunnel into the tree.
   * * * * If collision_p is not NULL (0), it means tunnel is already present.
   */
  hashtable_ts_insert(state->sgw_state.s11teid2mme, local_teid, new_tunnel);

  return new_tunnel;
}

//-----------------------------------------------------------------------------
int sgw_cm_remove_s11_tunnel(spgw_state_t *state, teid_t local_teid)
//-----------------------------------------------------------------------------
{
  int temp = 0;

  temp = hashtable_ts_free(state->sgw_state.s11teid2mme, local_teid);
  return temp;
}

//-----------------------------------------------------------------------------
s_plus_p_gw_eps_bearer_context_information_t *
sgw_cm_create_bearer_context_information_in_collection(
  spgw_state_t *state,
  teid_t teid)
{
  s_plus_p_gw_eps_bearer_context_information_t *new_bearer_context_information =
    NULL;

  new_bearer_context_information =
    calloc(1, sizeof(s_plus_p_gw_eps_bearer_context_information_t));

  if (new_bearer_context_information == NULL) {
    /*
     * Malloc failed, may be ENOMEM error
     */
    OAILOG_ERROR(
      LOG_SPGW_APP,
      "Failed to create new bearer context information object for S11 "
      "remote_teid " TEID_FMT "\n",
      teid);
    return NULL;
  }

  OAILOG_DEBUG(
    LOG_SPGW_APP,
    "sgw_cm_create_bearer_context_information_in_collection " TEID_FMT "\n",
    teid);

  bstring b = bfromcstr("pgw_eps_bearer_ctxt_info_apns");
  new_bearer_context_information->pgw_eps_bearer_context_information.apns =
    obj_hashtable_ts_create(
      32, NULL, NULL, (void (*)(void **)) pgw_lite_cm_free_apn, b);
  bdestroy_wrapper(&b);

  if (
    new_bearer_context_information->pgw_eps_bearer_context_information.apns ==
    NULL) {
    OAILOG_ERROR(
      LOG_SPGW_APP,
      "Failed to create APN collection object entry for EPS bearer S11 "
      "teid " TEID_FMT "\n",
      teid);
    sgw_free_s11_bearer_context_information(&new_bearer_context_information);
    return NULL;
  }

  /*
   * Trying to insert the new tunnel into the tree.
   * * * * If collision_p is not NULL (0), it means tunnel is already present.
   */
  hashtable_ts_insert(
    state->sgw_state.s11_bearer_context_information,
    teid,
    new_bearer_context_information);
  OAILOG_DEBUG(
    LOG_SPGW_APP,
    "Added new s_plus_p_gw_eps_bearer_context_information_t in "
    "s11_bearer_context_information_hashtable key teid " TEID_FMT "\n",
    teid);
  return new_bearer_context_information;
}

//-----------------------------------------------------------------------------
int sgw_cm_remove_bearer_context_information(spgw_state_t *state, teid_t teid)
{
  int temp = 0;

  temp =
    hashtable_ts_free(state->sgw_state.s11_bearer_context_information, teid);
  return temp;
}

//--- EPS Bearer Entry

//-----------------------------------------------------------------------------
sgw_eps_bearer_ctxt_t *sgw_cm_create_eps_bearer_ctxt_in_collection(
  sgw_pdn_connection_t *const sgw_pdn_connection,
  const ebi_t eps_bearer_idP)
{
  sgw_eps_bearer_ctxt_t *new_eps_bearer_entry = NULL;

  AssertFatal(sgw_pdn_connection, "Bad parameter sgw_pdn_connection");
  AssertFatal(
    (eps_bearer_idP >= EPS_BEARER_IDENTITY_FIRST) &&
      (eps_bearer_idP <= EPS_BEARER_IDENTITY_LAST),
    "Bad parameter ebi %u",
    eps_bearer_idP);

  if (!sgw_pdn_connection
         ->sgw_eps_bearers_array[EBI_TO_INDEX(eps_bearer_idP)]) {
    new_eps_bearer_entry = calloc(1, sizeof(sgw_eps_bearer_ctxt_t));

    if (new_eps_bearer_entry == NULL) {
      /*
       * Malloc failed, may be ENOMEM error
       */
      OAILOG_ERROR(
        LOG_SPGW_APP,
        "Failed to create EPS bearer entry for EPS bearer id %u \n",
        eps_bearer_idP);
      return NULL;
    }

    new_eps_bearer_entry->eps_bearer_id = eps_bearer_idP;
    sgw_pdn_connection->sgw_eps_bearers_array[EBI_TO_INDEX(eps_bearer_idP)] =
      new_eps_bearer_entry;
    OAILOG_DEBUG(
      LOG_SPGW_APP,
      "Inserted new EPS bearer entry for EPS bearer id  %u at address %p\n",
      eps_bearer_idP,
      (void *) sgw_pdn_connection
        ->sgw_eps_bearers_array[EBI_TO_INDEX(eps_bearer_idP)]);
  } else {
    OAILOG_WARNING(
      LOG_SPGW_APP,
      "Could not create mew EPS bearer ctxt for EPS bearer id %u : already "
      "exist\n",
      eps_bearer_idP);
  }
  return new_eps_bearer_entry;
}
//-----------------------------------------------------------------------------
sgw_eps_bearer_ctxt_t *sgw_cm_insert_eps_bearer_ctxt_in_collection(
  sgw_pdn_connection_t *const sgw_pdn_connection,
  sgw_eps_bearer_ctxt_t *const sgw_eps_bearer_ctxt)
{
  if (!sgw_eps_bearer_ctxt) {
    OAILOG_ERROR(
      LOG_SPGW_APP, "Failed to insert EPS bearer context : NULL context\n");
    return NULL;
  }

  if (!sgw_pdn_connection->sgw_eps_bearers_array[EBI_TO_INDEX(
        sgw_eps_bearer_ctxt->eps_bearer_id)]) {
    sgw_pdn_connection->sgw_eps_bearers_array[EBI_TO_INDEX(
      sgw_eps_bearer_ctxt->eps_bearer_id)] = sgw_eps_bearer_ctxt;
    OAILOG_DEBUG(
      LOG_SPGW_APP,
      "Inserted new EPS bearer entry for EPS bearer id %u \n",
      sgw_eps_bearer_ctxt->eps_bearer_id);
  } else {
    OAILOG_WARNING(
      LOG_SPGW_APP,
      "Could not create mew EPS bearer ctxt for EPS bearer id %u : already "
      "exist\n",
      sgw_eps_bearer_ctxt->eps_bearer_id);
  }
  return sgw_eps_bearer_ctxt;
}

//-----------------------------------------------------------------------------
sgw_eps_bearer_ctxt_t *sgw_cm_get_eps_bearer_entry(
  sgw_pdn_connection_t *const sgw_pdn_connection,
  ebi_t ebi)
{
  if ((ebi < EPS_BEARER_IDENTITY_FIRST) || (ebi > EPS_BEARER_IDENTITY_LAST)) {
    return NULL;
  }

  return sgw_pdn_connection->sgw_eps_bearers_array[EBI_TO_INDEX(ebi)];
}

//-----------------------------------------------------------------------------
int sgw_cm_remove_eps_bearer_entry(
  sgw_pdn_connection_t *const sgw_pdn_connection,
  ebi_t ebi)
//-----------------------------------------------------------------------------
{
  if ((ebi < EPS_BEARER_IDENTITY_FIRST) || (ebi > EPS_BEARER_IDENTITY_LAST)) {
    return RETURNerror;
  }
  /*sgw_eps_bearer_ctxt_t * sgw_eps_bearer_ctxt = sgw_pdn_connection->sgw_eps_bearers_array[EBI_TO_INDEX(ebi)];
  if (sgw_eps_bearer_ctxt) {
    sgw_free_sgw_eps_bearer_context(&sgw_eps_bearer_ctxt);
    return RETURNok;
  }*/
  return RETURNerror;
}
