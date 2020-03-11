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
 *------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#include "spgw_state.h"

#include <cstdlib>

extern "C" {
#include "assertions.h"
#include "bstrlib.h"
#include "dynamic_memory_check.h"
}

#include "spgw_state_manager.h"

using magma::lte::SpgwStateManager;

int spgw_state_init(bool persist_state, const spgw_config_t* config)
{
  SpgwStateManager::getInstance().init(persist_state, config);
  return RETURNok;
}

spgw_state_t* get_spgw_state(bool read_from_db)
{
  return SpgwStateManager::getInstance().get_state(read_from_db);
}

void spgw_state_exit()
{
  SpgwStateManager::getInstance().free_state();
}

void put_spgw_state()
{
  SpgwStateManager::getInstance().write_state_to_db();
}

void put_spgw_imsi_map() {
  SpgwStateManager::getInstance().put_spgw_imsi_map();
}

spgw_imsi_map_t* get_spgw_imsi_map() {
  return SpgwStateManager::getInstance().get_spgw_imsi_map();
}

void sgw_free_s11_bearer_context_information(
  s_plus_p_gw_eps_bearer_context_information_t** context_p)
{
  if (*context_p) {
    sgw_free_pdn_connection(
      &(*context_p)->sgw_eps_bearer_context_information.pdn_connection);

    if ((*context_p)->pgw_eps_bearer_context_information.apns) {
      obj_hashtable_ts_destroy(
        (*context_p)->pgw_eps_bearer_context_information.apns);
    }

    free_wrapper((void**) context_p);
  }
}

void sgw_free_pdn_connection(sgw_pdn_connection_t* pdn_connection_p)
{
  if (pdn_connection_p) {
    if (pdn_connection_p->apn_in_use) {
      free_wrapper((void**) &pdn_connection_p->apn_in_use);
    }

    for (auto& ebix : pdn_connection_p->sgw_eps_bearers_array) {
      sgw_free_eps_bearer_context(&ebix);
    }
  }
}

void sgw_free_eps_bearer_context(sgw_eps_bearer_ctxt_t** sgw_eps_bearer_ctxt)
{
  if (*sgw_eps_bearer_ctxt) {
    free_wrapper((void**) sgw_eps_bearer_ctxt);
  }
}

void pgw_free_pcc_rule(void** rule)
{
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
