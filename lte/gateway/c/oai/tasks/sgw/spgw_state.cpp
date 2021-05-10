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

#include "spgw_state.h"

#include <cstdlib>
#include <conversions.h>

extern "C" {
#include "assertions.h"
#include "bstrlib.h"
#include "dynamic_memory_check.h"
#include "pgw_procedures.h"
#include "sgw_context_manager.h"
}

#include "spgw_state_manager.h"

using magma::lte::SpgwStateManager;

int spgw_state_init(bool persist_state, const spgw_config_t* config) {
  SpgwStateManager::getInstance().init(persist_state, config);
  return RETURNok;
}

spgw_state_t* get_spgw_state(bool read_from_db) {
  return SpgwStateManager::getInstance().get_state(read_from_db);
}

hash_table_ts_t* get_spgw_ue_state() {
  OAILOG_DEBUG(
      LOG_SPGW_APP, "get_spgw_ue_state called by thread id %lu",
      pthread_self());
  return SpgwStateManager::getInstance().get_ue_state_ht();
}

hash_table_ts_t* get_spgw_teid_state() {
  return SpgwStateManager::getInstance().get_state_teid_ht();
}

int read_spgw_ue_state_db() {
  return SpgwStateManager::getInstance().read_ue_state_from_db();
}

void spgw_state_exit() {
  SpgwStateManager::getInstance().free_state();
}

void put_spgw_state() {
  SpgwStateManager::getInstance().write_state_to_db();
}

void put_spgw_ue_state(imsi64_t imsi64) {
  if (SpgwStateManager::getInstance().is_persist_state_enabled()) {
    spgw_ue_context_t* ue_context_p = nullptr;
    hash_table_ts_t* spgw_ue_state  = get_spgw_ue_state();
    hashtable_ts_get(
        spgw_ue_state, (const hash_key_t) imsi64, (void**) &ue_context_p);
    if (ue_context_p) {
      auto imsi_str = SpgwStateManager::getInstance().get_imsi_str(imsi64);
      SpgwStateManager::getInstance().write_ue_state_to_db(
          ue_context_p, imsi_str);
    }
  }
}

void delete_spgw_ue_state(imsi64_t imsi64) {
  // Parsing IMSI with fixed digits length
  auto imsi_str = SpgwStateManager::getInstance().get_imsi_str(imsi64);
  SpgwStateManager::getInstance().clear_ue_state_db(imsi_str);
}

void spgw_free_s11_bearer_context_information(
    s_plus_p_gw_eps_bearer_context_information_t** context_p) {
  if (*context_p) {
    sgw_free_pdn_connection(
        &(*context_p)->sgw_eps_bearer_context_information.pdn_connection);
    clear_protocol_configuration_options(
        &(*context_p)->sgw_eps_bearer_context_information.saved_message.pco);
    pgw_delete_procedures(*context_p);
    if ((*context_p)->pgw_eps_bearer_context_information.apns) {
      obj_hashtable_ts_destroy(
          (*context_p)->pgw_eps_bearer_context_information.apns);
    }

    free_wrapper((void**) context_p);
  }
}

void sgw_free_pdn_connection(sgw_pdn_connection_t* pdn_connection_p) {
  if (pdn_connection_p) {
    if (pdn_connection_p->apn_in_use) {
      free_wrapper((void**) &pdn_connection_p->apn_in_use);
    }

    for (auto& ebix : pdn_connection_p->sgw_eps_bearers_array) {
      sgw_free_eps_bearer_context(&ebix);
    }
  }
}

void sgw_free_eps_bearer_context(sgw_eps_bearer_ctxt_t** sgw_eps_bearer_ctxt) {
  if (*sgw_eps_bearer_ctxt) {
    free_wrapper((void**) sgw_eps_bearer_ctxt);
  }
}

void pgw_free_pcc_rule(void** rule) {
  if (rule) {
    auto* pcc_rule = (pcc_rule_t*) *rule;
    if (pcc_rule) {
      if (pcc_rule->name) {
        bdestroy_wrapper(&pcc_rule->name);
      }
      free_wrapper(rule);
    }
  }
}

void sgw_free_ue_context(spgw_ue_context_t** ue_context_p) {
  if (*ue_context_p) {
    sgw_s11_teid_t* p1 = LIST_FIRST(&(*ue_context_p)->sgw_s11_teid_list);
    sgw_s11_teid_t* p2 = nullptr;
    while (p1) {
      p2 = LIST_NEXT(p1, entries);
      LIST_REMOVE(p1, entries);
      free_wrapper((void**) &p1);
      p1 = p2;
    }
    free_wrapper((void**) ue_context_p);
  }
}
