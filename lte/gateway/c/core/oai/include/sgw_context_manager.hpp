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

/*! \file sgw_context_manager.hpp
 * \brief
 * \author Lionel Gauthier
 * \company Eurecom
 * \email: lionel.gauthier@eurecom.fr
 */
#pragma once

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.007.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/include/spgw_state.hpp"

#define INITIAL_SGW_S8_S1U_TEID 0x7FFFFFFF
void sgw_display_sgw_eps_bearer_context(
    const magma::lte::oai::SgwEpsBearerContext* eps_bearer_ctxt);
void sgw_display_s11_bearer_context_information(
    log_proto_t module,
    magma::lte::oai::S11BearerContext* sgw_context_information);

// TODO(rsarwad): to be removed while porting sgw_s8 task
void sgw_s8_display_sgw_eps_bearer_context(
    const sgw_eps_bearer_ctxt_t* eps_bearer_ctxt);

void sgw_s8_display_s11_bearer_context_information(
    log_proto_t module,
    sgw_eps_bearer_context_information_t* sgw_context_information);

sgw_eps_bearer_ctxt_t* sgw_s8_cm_get_eps_bearer_entry(
    sgw_pdn_connection_t* const sgw_pdn_connection, ebi_t ebi);

mme_sgw_tunnel_t* sgw_cm_create_s11_tunnel(teid_t remote_teid,
                                           teid_t local_teid);
magma::lte::oai::S11BearerContext*
sgw_cm_create_bearer_context_information_in_collection(teid_t teid);
magma::proto_map_rc_t sgw_cm_remove_bearer_context_information(teid_t teid,
                                                               imsi64_t imsi64);
sgw_eps_bearer_ctxt_t* sgw_cm_create_eps_bearer_ctxt_in_collection(
    sgw_pdn_connection_t* const sgw_pdn_connection, const ebi_t eps_bearer_idP);

magma::proto_map_rc_t sgw_cm_insert_eps_bearer_ctxt_in_collection(
    magma::lte::oai::SgwPdnConnection* const sgw_pdn_connection,
    magma::lte::oai::SgwEpsBearerContext* const sgw_eps_bearer_ctxt);

// Returns SPGW state pointer for given UE indexed by IMSI
magma::lte::oai::S11BearerContext* sgw_cm_get_spgw_context(teid_t teid);
spgw_ue_context_t* spgw_get_ue_context(imsi64_t imsi64);
spgw_ue_context_t* spgw_create_or_get_ue_context(imsi64_t imsi64);

status_code_e spgw_update_teid_in_ue_context(imsi64_t imsi64, teid_t teid);

magma::proto_map_rc_t sgw_cm_get_eps_bearer_entry(
    magma::lte::oai::SgwPdnConnection* const sgw_pdn_connection, ebi_t ebi,
    magma::lte::oai::SgwEpsBearerContext* bearer_context_p);

magma::proto_map_rc_t sgw_update_eps_bearer_entry(
    magma::lte::oai::SgwPdnConnection* sgw_pdn_connection_p, uint32_t ebi,
    magma::lte::oai::SgwEpsBearerContext* bearer_context_p);
