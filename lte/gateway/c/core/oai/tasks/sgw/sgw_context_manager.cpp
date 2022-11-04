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

/*! \file sgw_context_manager.cpp
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#define SGW
#define SGW_CONTEXT_MANAGER_C

#include "lte/gateway/c/core/oai/include/sgw_context_manager.hpp"

#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <inttypes.h>

#include "lte/gateway/c/core/common/dynamic_memory_check.h"

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#include "lte/gateway/c/core/oai/lib/hashtable/hashtable.h"
#include "lte/gateway/c/core/oai/lib/hashtable/obj_hashtable.h"
#ifdef __cplusplus
}
#endif

//-----------------------------------------------------------------------------
void sgw_display_sgw_eps_bearer_context(
    const sgw_eps_bearer_ctxt_t* const eps_bearer_ctxt)
//-----------------------------------------------------------------------------
{
  if (eps_bearer_ctxt) {
    OAILOG_DEBUG(LOG_SPGW_APP,
                 "|\t\t\t\tebi: %u, enb_teid_for_S1u: " TEID_FMT
                 ", s_gw_teid_for_S1u_S12_S4_up: " TEID_FMT " (tbc)\n",
                 eps_bearer_ctxt->eps_bearer_id, eps_bearer_ctxt->enb_teid_S1u,
                 eps_bearer_ctxt->s_gw_teid_S1u_S12_S4_up);
  }
}

//-----------------------------------------------------------------------------
void sgw_display_s11_bearer_context_information(
    log_proto_t module,
    sgw_eps_bearer_context_information_t* sgw_context_information)
//-----------------------------------------------------------------------------
{
  OAILOG_DEBUG(module, "| KEY %" PRId64 ":      \n",
               sgw_context_information->imsi64);
  OAILOG_DEBUG(module, "|\tsgw_eps_bearer_context_information:     |\n");
  // Imsi_t               imsi;                           ///< IMSI
  // (International Mobile Subscriber Identity) is the subscriber permanent
  // identity.
  OAILOG_DEBUG(module, "|\t\timsi_unauthenticated_indicator:\t%u\n",
               sgw_context_information->imsi_unauthenticated_indicator);
  // char                 msisdn[MSISDN_LENGTH];          ///< The basic MSISDN
  // of the UE. The presence is dictated by its storage in the HSS.
  OAILOG_DEBUG(module, "|\t\tmme_teid_    S11:              \t" TEID_FMT "\n",
               sgw_context_information->mme_teid_S11);
  // ip_address_t         mme_ip_address_for_S11;         ///< MME IP address
  // the S11 interface.
  OAILOG_DEBUG(module, "|\t\ts_gw_teid_S11_S4:              \t" TEID_FMT "\n",
               sgw_context_information->s_gw_teid_S11_S4);
  // ip_address_t         s_gw_ip_address_for_S11_S4;     ///< S-GW IP address
  // for the S11 interface and the S4 Interface (control plane). cgi_t
  // last_known_cell_Id;             ///< This is the last location of the UE
  // known by the network
  OAILOG_DEBUG(module, "|\t\tpdn_connection:\n");
  OAILOG_DEBUG(module, "|\t\t\tapn_in_use:        %s\n",
               sgw_context_information->pdn_connection.apn_in_use);
  OAILOG_DEBUG(module, "|\t\t\tdefault_bearer:    %u\n",
               sgw_context_information->pdn_connection.default_bearer);
  OAILOG_DEBUG(module, "|\t\t\teps_bearers:\n");
  for (int ebix = 0; ebix < BEARERS_PER_UE; ebix++) {
    sgw_display_sgw_eps_bearer_context(
        sgw_context_information->pdn_connection.sgw_eps_bearers_array[ebix]);
  }
  // void                  *trxn;
  // uint32_t               peer_ip;
}

//-----------------------------------------------------------------------------
void pgw_lite_cm_free_apn(pgw_apn_t** apnP)
//-----------------------------------------------------------------------------
{
  if (*apnP) {
    if ((*apnP)->pdn_connections) {
      obj_hashtable_ts_destroy((*apnP)->pdn_connections);
    }
  }
}

//-----------------------------------------------------------------------------
mme_sgw_tunnel_t* sgw_cm_create_s11_tunnel(teid_t remote_teid,
                                           teid_t local_teid)
//-----------------------------------------------------------------------------
{
  mme_sgw_tunnel_t* new_tunnel = new mme_sgw_tunnel_t();
  new_tunnel->remote_teid = remote_teid;
  new_tunnel->local_teid = local_teid;

  return new_tunnel;
}

//-----------------------------------------------------------------------------
s_plus_p_gw_eps_bearer_context_information_t*
sgw_cm_create_bearer_context_information_in_collection(teid_t teid) {
  s_plus_p_gw_eps_bearer_context_information_t* new_bearer_context_information =
      nullptr;
  new_bearer_context_information =
      new s_plus_p_gw_eps_bearer_context_information_t();

  /*
   * Trying to insert the new tunnel into the map.
   * * * * If collision_p is not NULL (0), it means tunnel is already present.
   */
  state_teid_map_t* state_teid_map = get_spgw_teid_state();
  if (!state_teid_map) {
    OAILOG_ERROR(LOG_SPGW_APP, "Failed to get state_teid_map");
    return nullptr;
  }
  state_teid_map->insert(teid, new_bearer_context_information);

  OAILOG_DEBUG(LOG_SPGW_APP,
               "Added new s_plus_p_gw_eps_bearer_context_information_t in "
               "s11_bearer_context_information_hashtable key TEID " TEID_FMT
               "\n",
               teid);
  return new_bearer_context_information;
}

//-----------------------------------------------------------------------------
magma::proto_map_rc_t sgw_cm_remove_bearer_context_information(
    teid_t teid, imsi64_t imsi64) {
  state_teid_map_t* state_teid_map = get_spgw_teid_state();
  if (!state_teid_map) {
    OAILOG_ERROR(LOG_SPGW_APP, "Failed to get state_teid_map");
    return magma::PROTO_MAP_KEY_NOT_EXISTS;
  }
  if (state_teid_map->remove(teid) != magma::PROTO_MAP_OK) {
    OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                    "Failed to free teid from state_teid_map \n");
    return magma::PROTO_MAP_REMOVE_KEY_FAILED;
  }
  spgw_ue_context_t* ue_context_p = nullptr;
  map_uint64_spgw_ue_context_t* spgw_ue_state = get_spgw_ue_state();
  if (!spgw_ue_state) {
    OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64, "Failed to find spgw_ue_state");
    return magma::PROTO_MAP_SEARCH_NO_RESULT;
  }
  spgw_ue_state->get(imsi64, &ue_context_p);
  if (ue_context_p) {
    sgw_s11_teid_t* p1 = LIST_FIRST(&(ue_context_p->sgw_s11_teid_list));
    while (p1) {
      if (p1->sgw_s11_teid == teid) {
        LIST_REMOVE(p1, entries);
        free_cpp_wrapper(reinterpret_cast<void**>(&p1));
        break;
      }
      p1 = LIST_NEXT(p1, entries);
    }
    if (LIST_EMPTY(&ue_context_p->sgw_s11_teid_list)) {
      if (spgw_ue_state->remove(imsi64) != magma::PROTO_MAP_OK) {
        OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                        "Failed to free imsi64 from imsi_ue_context_map \n");
        return magma::PROTO_MAP_REMOVE_KEY_FAILED;
      }
      delete_spgw_ue_state(imsi64);
    }
  }
  return magma::PROTO_MAP_OK;
}

//--- EPS Bearer Entry
//-----------------------------------------------------------------------------
sgw_eps_bearer_ctxt_t* sgw_cm_create_eps_bearer_ctxt_in_collection(
    sgw_pdn_connection_t* const sgw_pdn_connection,
    const ebi_t eps_bearer_idP) {
  sgw_eps_bearer_ctxt_t* new_eps_bearer_entry = nullptr;

  AssertFatal(sgw_pdn_connection, "Bad parameter sgw_pdn_connection");
  AssertFatal((eps_bearer_idP >= EPS_BEARER_IDENTITY_FIRST) &&
                  (eps_bearer_idP <= EPS_BEARER_IDENTITY_LAST),
              "Bad parameter ebi %u", eps_bearer_idP);

  if (!sgw_pdn_connection
           ->sgw_eps_bearers_array[EBI_TO_INDEX(eps_bearer_idP)]) {
    new_eps_bearer_entry = new sgw_eps_bearer_ctxt_t();

    new_eps_bearer_entry->eps_bearer_id = eps_bearer_idP;
    sgw_pdn_connection->sgw_eps_bearers_array[EBI_TO_INDEX(eps_bearer_idP)] =
        new_eps_bearer_entry;
    OAILOG_DEBUG(
        LOG_SPGW_APP,
        "Inserted new EPS bearer entry for EPS bearer id  %u at address %p\n",
        eps_bearer_idP,
        (void*)sgw_pdn_connection
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
sgw_eps_bearer_ctxt_t* sgw_cm_insert_eps_bearer_ctxt_in_collection(
    sgw_pdn_connection_t* const sgw_pdn_connection,
    sgw_eps_bearer_ctxt_t* const sgw_eps_bearer_ctxt) {
  if (!sgw_eps_bearer_ctxt) {
    OAILOG_ERROR(LOG_SPGW_APP,
                 "Failed to insert EPS bearer context : nullptr context\n");
    return nullptr;
  }

  if (!sgw_pdn_connection->sgw_eps_bearers_array[EBI_TO_INDEX(
          sgw_eps_bearer_ctxt->eps_bearer_id)]) {
    sgw_pdn_connection->sgw_eps_bearers_array[EBI_TO_INDEX(
        sgw_eps_bearer_ctxt->eps_bearer_id)] = sgw_eps_bearer_ctxt;
    OAILOG_DEBUG(LOG_SPGW_APP,
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
sgw_eps_bearer_ctxt_t* sgw_cm_get_eps_bearer_entry(
    sgw_pdn_connection_t* const sgw_pdn_connection, ebi_t ebi) {
  if ((ebi < EPS_BEARER_IDENTITY_FIRST) || (ebi > EPS_BEARER_IDENTITY_LAST)) {
    return nullptr;
  }

  return sgw_pdn_connection->sgw_eps_bearers_array[EBI_TO_INDEX(ebi)];
}

s_plus_p_gw_eps_bearer_context_information_t* sgw_cm_get_spgw_context(
    teid_t teid) {
  s_plus_p_gw_eps_bearer_context_information_t* spgw_bearer_context_info =
      nullptr;
  state_teid_map_t* state_teid_map = get_spgw_teid_state();
  if (!state_teid_map) {
    OAILOG_ERROR(LOG_SPGW_APP, "Failed to get state_teid_map");
    return nullptr;
  }

  state_teid_map->get(teid, &spgw_bearer_context_info);
  return spgw_bearer_context_info;
}

spgw_ue_context_t* spgw_get_ue_context(imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  spgw_ue_context_t* ue_context_p = nullptr;
  map_uint64_spgw_ue_context_t* state_ue_map = get_spgw_ue_state();
  if (!state_ue_map) {
    OAILOG_ERROR(LOG_SPGW_APP, "Failed to find spgw_ue_state");
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, nullptr);
  }
  state_ue_map->get(imsi64, &ue_context_p);
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, ue_context_p);
}

spgw_ue_context_t* spgw_create_or_get_ue_context(imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  spgw_ue_context_t* ue_context_p = nullptr;
  map_uint64_spgw_ue_context_t* state_ue_map = get_spgw_ue_state();
  if (!state_ue_map) {
    OAILOG_ERROR(LOG_SPGW_APP, "Failed to find spgw_ue_state");
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, nullptr);
  }
  state_ue_map->get(imsi64, &ue_context_p);
  if (!ue_context_p) {
    ue_context_p = new spgw_ue_context_t();
    if (ue_context_p) {
      LIST_INIT(&ue_context_p->sgw_s11_teid_list);
      state_ue_map->insert(imsi64, ue_context_p);
    } else {
      OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                      "Failed to allocate memory for UE context \n");
    }
  }
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, ue_context_p);
}

status_code_e spgw_update_teid_in_ue_context(imsi64_t imsi64, teid_t teid) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  spgw_ue_context_t* ue_context_p = spgw_create_or_get_ue_context(imsi64);
  if (!ue_context_p) {
    OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                    "Failed to get UE context for sgw_s11_teid " TEID_FMT "\n",
                    teid);
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
  }

  sgw_s11_teid_t* sgw_s11_teid_p = new sgw_s11_teid_t();
  sgw_s11_teid_p->sgw_s11_teid = teid;
  LIST_INSERT_HEAD(&ue_context_p->sgw_s11_teid_list, sgw_s11_teid_p, entries);
  OAILOG_DEBUG(LOG_SPGW_APP,
               "Inserted sgw_s11_teid to list of teids of UE context" TEID_FMT
               "\n",
               teid);
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNok);
}
