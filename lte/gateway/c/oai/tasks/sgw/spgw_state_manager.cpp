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

#include "spgw_state_manager.h"

spgw_state_t *SpgwStateManager::create_spgw_state(spgw_config_t *config)
{
  // Allocating spgw_state_p
  spgw_state_t *state_p;
  state_p = (spgw_state_t *) calloc(1, sizeof(spgw_state_t));

  bstring b = bfromcstr(SGW_S11_TEID_MME_HT_NAME);
  state_p->sgw_state.s11teid2mme =
    hashtable_ts_create(SGW_STATE_CONTEXT_HT_MAX_SIZE, nullptr, nullptr, b);
  btrunc(b, 0);

  bassigncstr(b, S11_BEARER_CONTEXT_INFO_HT_NAME);
  state_p->sgw_state.s11_bearer_context_information = hashtable_ts_create(
    SGW_STATE_CONTEXT_HT_MAX_SIZE,
    nullptr,
    (void (*)(void **)) sgw_free_s11_bearer_context_information,
    b);
  bdestroy_wrapper(&b);

  state_p->sgw_state.sgw_ip_address_S1u_S12_S4_up.s_addr =
    config->sgw_config.ipv4.S1u_S12_S4_up.s_addr;

  // TODO: Refactor GTPv1u_data state
  state_p->sgw_state.gtpv1u_data.sgw_ip_address_for_S1u_S12_S4_up =
    state_p->sgw_state.sgw_ip_address_S1u_S12_S4_up;

  return state_p;
}

void SpgwStateManager::free_spgw_state()
{
  if (
    hashtable_ts_destroy(spgw_state_cache_p->sgw_state.s11teid2mme) !=
    HASH_TABLE_OK) {
    OAI_FPRINTF_ERR(
      "An error occurred while destroying SGW s11teid2mme hashtable");
  }

  if (
    hashtable_ts_destroy(
      spgw_state_cache_p->sgw_state.s11_bearer_context_information) !=
    HASH_TABLE_OK) {
    OAI_FPRINTF_ERR(
      "An error occurred while destroying SGW s11_bearer_context_information "
      "hashtable");
  }

  free(spgw_state_cache_p);
}

int SpgwStateManager::read_state_from_db()
{
  // TODO: Implement read from redis db
  return RETURNok;
}

void SpgwStateManager::write_state_to_db()
{
  // TODO: Implement put to redis db
  AssertFatal(
    this->state_accessed, "Tried to put SPGW state while it was not in use");
  this->state_accessed = false;
}
