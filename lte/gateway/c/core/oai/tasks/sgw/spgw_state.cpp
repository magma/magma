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
  if (ptr) {
    s_plus_p_gw_eps_bearer_context_information_t* context_p =
        reinterpret_cast<s_plus_p_gw_eps_bearer_context_information_t*>(*ptr);
    if (context_p) {
      sgw_free_pdn_connection(
          &(context_p->sgw_eps_bearer_context_information.pdn_connection));
      clear_protocol_configuration_options(
          &(context_p->sgw_eps_bearer_context_information.saved_message.pco));
      delete_pending_procedures(
          &(context_p->sgw_eps_bearer_context_information));

      free_cpp_wrapper(reinterpret_cast<void**>(ptr));
    }
  }
}

void sgw_free_pdn_connection(sgw_pdn_connection_t* pdn_connection_p) {
  if (pdn_connection_p) {
    if (pdn_connection_p->apn_in_use) {
      free_wrapper((void**)&pdn_connection_p->apn_in_use);
    }

    for (auto& ebix : pdn_connection_p->sgw_eps_bearers_array) {
      sgw_free_eps_bearer_context(&ebix);
    }
  }
}

void sgw_free_eps_bearer_context(sgw_eps_bearer_ctxt_t** sgw_eps_bearer_ctxt) {
  if (*sgw_eps_bearer_ctxt) {
    if ((*sgw_eps_bearer_ctxt)->pgw_cp_ip_port) {
      free_wrapper(
          reinterpret_cast<void**>(&(*sgw_eps_bearer_ctxt)->pgw_cp_ip_port));
    }
    free_cpp_wrapper(reinterpret_cast<void**>(sgw_eps_bearer_ctxt));
  }
}

void sgw_free_ue_context(void** ptr) {
  if (!ptr) {
    OAILOG_ERROR(LOG_SPGW_APP,
                 "sgw_free_ue_context received invalid pointer for deletion");
    return;
  }
  spgw_ue_context_t* ue_context_p = reinterpret_cast<spgw_ue_context_t*>(*ptr);
  if (ue_context_p) {
    sgw_s11_teid_t* p1 = LIST_FIRST(&(ue_context_p)->sgw_s11_teid_list);
    sgw_s11_teid_t* p2 = nullptr;
    while (p1) {
      p2 = LIST_NEXT(p1, entries);
      LIST_REMOVE(p1, entries);
      free_cpp_wrapper(reinterpret_cast<void**>(&p1));
      p1 = p2;
    }
    free_cpp_wrapper(ptr);
  }
}
