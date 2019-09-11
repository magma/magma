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

#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include <assertions.h>
#include <dynamic_memory_check.h>
#include <cstdlib>
#include <common_defs.h>

#include "spgw_state.h"

#define SGW_STATE_CONTEXT_HT_MAX_SIZE 512
#define SGW_S11_TEID_MME_HT_NAME "sgw_s11_teid2mme_htbl"
#define S11_BEARER_CONTEXT_INFO_HT_NAME "s11_bearer_context_information_htbl"

class SpgwStateManager {
 private:
  spgw_state_t *spgw_state_cache_p {};
  bool state_accessed;

 public:
  bool persist_state;

  SpgwStateManager()
  {
    this->persist_state = false;
    this->state_accessed = false;
  }

  SpgwStateManager(
    bool persist_state,
    spgw_config_t *config)
  {
    this->persist_state = persist_state;
    this->state_accessed = false;
    spgw_state_cache_p = create_spgw_state(config);
  }

  spgw_state_t *get_spgw_state()
  {
    this->state_accessed = true;

    AssertFatal(
      spgw_state_cache_p != nullptr, "SPGW state cache is NULL");

    return spgw_state_cache_p;
  }

  static spgw_state_t *create_spgw_state(spgw_config_t *config);
  void free_spgw_state();

  // TODO: Implement redis r/w functions
  int read_state_from_db();
  void write_state_to_db();
};

#ifdef __cplusplus
}
#endif
