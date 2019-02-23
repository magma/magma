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

#include "3gpp_23.401.h"

/********************************
*     Paired contexts           *
*********************************/
// data entry for s11_bearer_context_information_hashtable
// like this if needed in future, the split of S and P GW should be easier.
typedef struct s_plus_p_gw_eps_bearer_context_information_s {
  sgw_eps_bearer_context_information_t sgw_eps_bearer_context_information;
  pgw_eps_bearer_context_information_t pgw_eps_bearer_context_information;
} s_plus_p_gw_eps_bearer_context_information_t;

// data entry for s11teid2mme_hashtable
typedef struct mme_sgw_tunnel_s {
  uint32_t local_teid;  ///< Tunnel endpoint Identifier
  uint32_t remote_teid; ///< Tunnel endpoint Identifier
} mme_sgw_tunnel_t;

// data entry for s1uteid2enb_hashtable
typedef struct enb_sgw_s1u_tunnel_s {
  uint32_t local_teid;                 ///< S-GW Tunnel endpoint Identifier
  uint32_t remote_teid;                ///< eNB Tunnel endpoint Identifier
  ip_address_t enb_ip_address_for_S11; ///< eNB IP address the S1U interface.
} enb_sgw_s1u_tunnel_t;

void sgw_display_s11teid2mme_mappings(void);
void sgw_display_sgw_eps_bearer_context(
  const sgw_eps_bearer_ctxt_t *const eps_bearer_ctxt);
void sgw_display_s11_bearer_context_information_mapping(void);
void pgw_lite_cm_free_apn(pgw_apn_t **apnP);

teid_t sgw_get_new_S11_tunnel_id(void);
mme_sgw_tunnel_t *sgw_cm_create_s11_tunnel(
  teid_t remote_teid,
  teid_t local_teid);
int sgw_cm_remove_s11_tunnel(teid_t local_teid);
sgw_eps_bearer_ctxt_t *sgw_cm_create_eps_bearer_context(void);
sgw_pdn_connection_t *sgw_cm_create_pdn_connection(void);
void sgw_cm_free_pdn_connection(sgw_pdn_connection_t *pdn_connectionP);
void sgw_free_sgw_eps_bearer_context(
  sgw_eps_bearer_ctxt_t **sgw_eps_bearer_ctxt);
s_plus_p_gw_eps_bearer_context_information_t *
sgw_cm_create_bearer_context_information_in_collection(teid_t teid);
void sgw_cm_free_s_plus_p_gw_eps_bearer_context_information(
  s_plus_p_gw_eps_bearer_context_information_t **contextP);
int sgw_cm_remove_bearer_context_information(teid_t teid);
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
  sgw_pdn_connection_t *const sgw_pdn_connection,
  ebi_t eps_bearer_idP);

#endif /* FILE_SGW_CONTEXT_MANAGER_SEEN */
