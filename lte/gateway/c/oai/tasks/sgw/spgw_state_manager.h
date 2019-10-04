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

#ifdef __cplusplus
}
#endif

#include "spgw_state.h"

#define SGW_STATE_CONTEXT_HT_MAX_SIZE 512
#define SGW_S11_TEID_MME_HT_NAME "sgw_s11_teid2mme_htbl"
#define S11_BEARER_CONTEXT_INFO_HT_NAME "s11_bearer_context_information_htbl"
#define MAX_PREDEFINED_PCC_RULES_HT_SIZE 32

namespace magma {
namespace lte {

/**
 * SpgwStateManager is a singleton (thread-safe, destruction guaranteed)
 * that contains functions to maintain SGW and PGW state, allocating and
 * freeing state structs, and writing / reading state to db.
 */
class SpgwStateManager {
 public:
  /**
   * Returns an instance of SpgwStateManager, guaranteed to be thread safe and
   * initialized only once.
   * @return SpgwStateManager instance.
   */
  static SpgwStateManager& getInstance();
  /**
   * Initialization function to initialize member variables.
   * @param persist_state should read and write state from db
   * @param config SPGW config struct
   */
  void init(bool persist_state, const spgw_config_t* config);

  /**
   * Singleton class, copy constructor and assignment operator are marked
   * as deleted functions.
   */
  // Making them public for better debugging logging.
  SpgwStateManager(SpgwStateManager const&) = delete;
  SpgwStateManager& operator=(SpgwStateManager const&) = delete;

  /**
   * @return A pointer to spgw_state_cache
   */
  spgw_state_t* get_spgw_state();

  /**
   * Frees all memory allocated on spgw_state_t.
   */
  void free_spgw_state();

  // TODO: Implement redis r/w functions
  int read_state_from_db();
  void write_state_to_db();

 private:
  SpgwStateManager();
  ~SpgwStateManager() = default;

  /**
   * Allocates a new spgw_state_t struct, and inits hashtables and state
   * structs to default values.
   * @param config pointer to spgw_config
   * @return spgw_state pointer
   */
  spgw_state_t* create_spgw_state(const spgw_config_t* config);

  // Flag for check asserting if the state has been initialized.
  bool is_initialized_;
  // Flag for check asserting that write should be done after read.
  // TODO: Convert this to state versioning variable
  bool state_dirty_;
  // Flag for enabling writing and reading to db.
  bool persist_state_;
  // TODO: Make this a unique_ptr
  spgw_state_t* spgw_state_cache_p_;
};

} // namespace lte
} // namespace magma
