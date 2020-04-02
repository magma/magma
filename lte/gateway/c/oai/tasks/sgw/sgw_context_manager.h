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

/*! \file sgw_context_manager.h
* \brief
* \author Lionel Gauthier
* \company Eurecom
* \email: lionel.gauthier@eurecom.fr
*/
#ifndef FILE_SGW_CONTEXT_MANAGER_SEEN
#define FILE_SGW_CONTEXT_MANAGER_SEEN

#include <stdint.h>

#include "3gpp_24.007.h"

#include "spgw_state.h"

void sgw_display_sgw_eps_bearer_context(
  const sgw_eps_bearer_ctxt_t* eps_bearer_ctxt);
void sgw_display_s11teid2mme(mme_sgw_tunnel_t* mme_sgw_tunnel);
void sgw_display_s11_bearer_context_information(
  s_plus_p_gw_eps_bearer_context_information_t* sp_context_information);
void pgw_lite_cm_free_apn(pgw_apn_t **apnP);

teid_t sgw_get_new_S11_tunnel_id(spgw_state_t *state);
mme_sgw_tunnel_t *sgw_cm_create_s11_tunnel(
  teid_t remote_teid,
  teid_t local_teid);
s_plus_p_gw_eps_bearer_context_information_t*
sgw_cm_create_bearer_context_information_in_collection(
  spgw_state_t* spgw_state,
  teid_t teid,
  imsi64_t imsi64);
int sgw_cm_remove_bearer_context_information(teid_t teid, imsi64_t imsi64);
sgw_eps_bearer_ctxt_t *sgw_cm_create_eps_bearer_ctxt_in_collection(
  sgw_pdn_connection_t *const sgw_pdn_connection,
  const ebi_t eps_bearer_idP);
sgw_eps_bearer_ctxt_t *sgw_cm_insert_eps_bearer_ctxt_in_collection(
  sgw_pdn_connection_t *const sgw_pdn_connection,
  sgw_eps_bearer_ctxt_t *const sgw_eps_bearer_ctxt);
sgw_eps_bearer_ctxt_t *sgw_cm_get_eps_bearer_entry(
  sgw_pdn_connection_t *const sgw_pdn_connection,
  ebi_t ebi);
int sgw_cm_remove_eps_bearer_entry(
  sgw_pdn_connection_t* const sgw_pdn_connection,
  ebi_t eps_bearer_idP);
// Returns SPGW state pointer for given UE indexed by IMSI
s_plus_p_gw_eps_bearer_context_information_t* sgw_cm_get_spgw_context(
  teid_t teid);

#endif /* FILE_SGW_CONTEXT_MANAGER_SEEN */
