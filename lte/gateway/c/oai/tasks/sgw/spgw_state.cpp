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

#include <stdlib.h>

extern "C" {
#include "bstrlib.h"
#include "assertions.h"
#include "dynamic_memory_check.h"
}

#include "common_defs.h"
#include "spgw_state_manager.h"

static SpgwStateManager spgw_state_mgr;

int spgw_state_init(bool use_stateless)
{
  spgw_state_mgr = SpgwStateManager(use_stateless);
  if (spgw_state_mgr.persist_state) {
    return spgw_state_mgr.read_state_from_db();
  }
  return RETURNok;
}

spgw_state_t *get_spgw_state(void)
{
  spgw_state_t *state_cache_p = spgw_state_mgr.get_spgw_state();
  AssertFatal(state_cache_p != nullptr, "SPGW state cache is NULL");

  return state_cache_p;
}
