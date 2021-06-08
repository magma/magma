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

/*! \file sgw_context_manager.h
 * \brief
 * \author Lionel Gauthier
 * \company Eurecom
 * \email: lionel.gauthier@eurecom.fr
 */
#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

#include "3gpp_24.007.h"

#include "spgw_state.h"

#define INITIAL_SGW_S8_S1U_TEID 0x7FFFFFFF
void sgw_display_sgw_eps_bearer_context(
    const sgw_eps_bearer_ctxt_t* eps_bearer_ctxt);
void sgw_display_s11_bearer_context_information(
    log_proto_t module,
    sgw_eps_bearer_context_information_t* sgw_context_information);
void pgw_lite_cm_free_apn(pgw_apn_t** apnP);

mme_sgw_tunnel_t* sgw_cm_create_s11_tunnel(
    teid_t remote_teid, teid_t local_teid);
s_plus_p_gw_eps_bearer_context_information_t*
sgw_cm_create_bearer_context_information_in_collection(teid_t teid);
int sgw_cm_remove_bearer_context_information(teid_t teid, imsi64_t imsi64);
sgw_eps_bearer_ctxt_t* sgw_cm_create_eps_bearer_ctxt_in_collection(
    sgw_pdn_connection_t* const sgw_pdn_connection, const ebi_t eps_bearer_idP);
sgw_eps_bearer_ctxt_t* sgw_cm_insert_eps_bearer_ctxt_in_collection(
    sgw_pdn_connection_t* const sgw_pdn_connection,
    sgw_eps_bearer_ctxt_t* const sgw_eps_bearer_ctxt);
sgw_eps_bearer_ctxt_t* sgw_cm_get_eps_bearer_entry(
    sgw_pdn_connection_t* const sgw_pdn_connection, ebi_t ebi);
int sgw_cm_remove_eps_bearer_entry(
    sgw_pdn_connection_t* const sgw_pdn_connection, ebi_t eps_bearer_idP);
// Returns SPGW state pointer for given UE indexed by IMSI
s_plus_p_gw_eps_bearer_context_information_t* sgw_cm_get_spgw_context(
    teid_t teid);
spgw_ue_context_t* spgw_create_or_get_ue_context(imsi64_t imsi64);

int spgw_update_teid_in_ue_context(imsi64_t imsi64, teid_t teid);

#ifdef __cplusplus
}
#endif
