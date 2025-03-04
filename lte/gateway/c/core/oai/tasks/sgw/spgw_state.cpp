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
 *------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#include "lte/gateway/c/core/oai/include/spgw_state.hpp"

#include <cstdlib>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/include/sgw_context_manager.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/pgw_procedures.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/spgw_state_manager.hpp"

using magma::lte::SpgwStateManager;

int spgw_state_init(bool persist_state, const spgw_config_t* config) {
  SpgwStateManager::getInstance().init(persist_state, config);
  return RETURNok;
}

spgw_state_t* get_spgw_state(bool read_from_db) {
  return SpgwStateManager::getInstance().get_state(read_from_db);
}

map_uint64_spgw_ue_context_t* get_spgw_ue_state() {
  return SpgwStateManager::getInstance().get_spgw_ue_state_map();
}

state_teid_map_t* get_spgw_teid_state() {
  return SpgwStateManager::getInstance().get_state_teid_map();
}

int read_spgw_ue_state_db() {
  return SpgwStateManager::getInstance().read_ue_state_from_db();
}

void spgw_state_exit() { SpgwStateManager::getInstance().free_state(); }

void put_spgw_state() { SpgwStateManager::getInstance().write_state_to_db(); }

void put_spgw_ue_state(imsi64_t imsi64) {
  if (SpgwStateManager::getInstance().is_persist_state_enabled()) {
    spgw_ue_context_t* ue_context_p = nullptr;
    map_uint64_spgw_ue_context_t* spgw_ue_state = get_spgw_ue_state();
    if (!spgw_ue_state) {
      OAILOG_ERROR(LOG_SPGW_APP, "Failed to find spgw_ue_state");
      OAILOG_FUNC_OUT(LOG_SPGW_APP);
    }
    spgw_ue_state->get(imsi64, &ue_context_p);
    if (ue_context_p) {
      auto imsi_str = SpgwStateManager::getInstance().get_imsi_str(imsi64);
      SpgwStateManager::getInstance().write_ue_state_to_db(ue_context_p,
                                                           imsi_str);
    }
  }
}

void delete_spgw_ue_state(imsi64_t imsi64) {
  // Parsing IMSI with fixed digits length
  auto imsi_str = SpgwStateManager::getInstance().get_imsi_str(imsi64);
  SpgwStateManager::getInstance().clear_ue_state_db(imsi_str);
}

void spgw_free_s11_bearer_context_information(void** ptr) {
  if (!ptr) {
    OAILOG_ERROR(LOG_SPGW_APP, "Received null pointer");
    OAILOG_FUNC_OUT(LOG_SPGW_APP);
  }

  if (*ptr) {
    magma::lte::oai::S11BearerContext* context_p =
        reinterpret_cast<magma::lte::oai::S11BearerContext*>(*ptr);
    if (context_p) {
      magma::lte::oai::SgwEpsBearerContextInfo* sgw_context_p =
          context_p->mutable_sgw_eps_bearer_context();
      if (sgw_context_p) {
        sgw_free_pdn_connection(sgw_context_p->mutable_pdn_connection());
        if (sgw_context_p->saved_message().has_pco()) {
          if (sgw_context_p->saved_message().pco().pco_protocol_size()) {
            sgw_context_p->mutable_saved_message()
                ->mutable_pco()
                ->clear_pco_protocol();
          }
          sgw_context_p->mutable_saved_message()->clear_pco();
        }
        if (sgw_context_p->pending_procedures_size()) {
          sgw_context_p->clear_pending_procedures();
        }
        if (sgw_context_p->has_mme_cp_ip_address_s11()) {
          sgw_context_p->clear_mme_cp_ip_address_s11();
        }
      }  // end of sgw_context_p
      context_p->clear_sgw_eps_bearer_context();
      context_p->clear_pgw_eps_bearer_context();
    }  // end of s11_bearer_context_p
  }  // end of ptr
  free_cpp_wrapper(ptr);
}

void free_eps_bearer_context(
    magma::lte::oai::SgwEpsBearerContext* bearer_context_p) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  if (!bearer_context_p) {
    OAILOG_ERROR(LOG_SPGW_APP, "Received nullptr for bearer context");
    OAILOG_FUNC_OUT(LOG_SPGW_APP);
  }
  if (bearer_context_p->has_eps_bearer_qos()) {
    if (bearer_context_p->eps_bearer_qos().has_mbr()) {
      bearer_context_p->mutable_eps_bearer_qos()->clear_mbr();
    }
    if (bearer_context_p->eps_bearer_qos().has_gbr()) {
      bearer_context_p->mutable_eps_bearer_qos()->clear_gbr();
    }
    bearer_context_p->clear_eps_bearer_qos();
  }
  if (bearer_context_p->has_tft()) {
    if (bearer_context_p->tft().has_packet_filter_list()) {
      bearer_context_p->mutable_tft()->clear_packet_filter_list();
    }
    if (bearer_context_p->tft().has_parameters_list()) {
      bearer_context_p->mutable_tft()->clear_parameters_list();
    }
  }
  if (bearer_context_p->sdf_ids_size()) {
    bearer_context_p->clear_sdf_ids();
  }
  if (bearer_context_p->has_enb_s1u_ip_addr()) {
    bearer_context_p->clear_enb_s1u_ip_addr();
  }
  if (bearer_context_p->has_ue_ip_paa()) {
    bearer_context_p->clear_ue_ip_paa();
  }
}

void sgw_free_pdn_connection(
    magma::lte::oai::SgwPdnConnection* pdn_connection_p) {
  if (pdn_connection_p) {
    if (pdn_connection_p->eps_bearer_map_size()) {
      map_uint32_spgw_eps_bearer_context_t eps_bearer_map;
      eps_bearer_map.map = pdn_connection_p->mutable_eps_bearer_map();
      for (auto itr = eps_bearer_map.map->begin();
           itr != eps_bearer_map.map->end(); itr++) {
        magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt = itr->second;
        free_eps_bearer_context(&eps_bearer_ctxt);
      }
      eps_bearer_map.clear();
    }
    pdn_connection_p->Clear();
  }
}

void sgw_remove_eps_bearer_context(
    magma::lte::oai::SgwPdnConnection* pdn_connection_p, uint32_t ebi) {
  magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt;
  map_uint32_spgw_eps_bearer_context_t eps_bearer_map;
  eps_bearer_map.map = pdn_connection_p->mutable_eps_bearer_map();
  if (eps_bearer_map.get(ebi, &eps_bearer_ctxt) == magma::PROTO_MAP_OK) {
    free_eps_bearer_context(&eps_bearer_ctxt);
  }
  eps_bearer_map.remove(ebi);
}

void sgw_free_ue_context(void** ptr) {
  if (!ptr) {
    OAILOG_ERROR(LOG_SPGW_APP,
                 "sgw_free_ue_context received invalid pointer for deletion");
    return;
  }

  spgw_ue_context_t* ue_context_p = reinterpret_cast<spgw_ue_context_t*>(*ptr);
  if (!ue_context_p) {
    OAILOG_ERROR(LOG_SPGW_APP,
                 "sgw_free_ue_context received invalid pointer for deletion");
    return;
  }
  sgw_s11_teid_t* p1 = LIST_FIRST(&(ue_context_p)->sgw_s11_teid_list);
  sgw_s11_teid_t* p2 = nullptr;
  while (p1) {
    p2 = LIST_NEXT(p1, entries);
    LIST_REMOVE(p1, entries);
    free_cpp_wrapper(reinterpret_cast<void**>(&p1));
    p1 = p2;
  }
  free_cpp_wrapper(reinterpret_cast<void**>(ptr));
  return;
}
