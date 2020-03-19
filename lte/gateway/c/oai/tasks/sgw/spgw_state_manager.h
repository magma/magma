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

#include "state_manager.h"
#include "spgw_state.h"
#include "spgw_state_converter.h"

namespace {
constexpr int SGW_STATE_CONTEXT_HT_MAX_SIZE = 512;
constexpr int MAX_PREDEFINED_PCC_RULES_HT_SIZE = 32;
constexpr char S11_BEARER_CONTEXT_INFO_HT_NAME[] =
  "s11_bearer_context_information_htbl";
constexpr char SPGW_STATE_TABLE_NAME[] = "spgw_state";
constexpr char SPGW_TASK_NAME[] = "SPGW";
} // namespace

namespace magma {
namespace lte {

/**
 * SpgwStateManager is a singleton (thread-safe, destruction guaranteed)
 * that contains functions to maintain SGW and PGW state, allocating and
 * freeing state structs, and writing / reading state to db.
 */
class SpgwStateManager :
  public StateManager<
    spgw_state_t,
    s_plus_p_gw_eps_bearer_context_information_t,
    gateway::spgw::SpgwState,
    gateway::spgw::S11BearerContext,
    SpgwStateConverter> {
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
   * Frees all memory allocated on spgw_state_t.
   */
  void free_state() override;

  int read_ue_state_from_db() override;

 private:
  SpgwStateManager();
  ~SpgwStateManager();

  /**
   * Allocates a new spgw_state_t struct, and inits hashtables and state
   * structs to default values.
   */
  void create_state() override;

  const spgw_config_t* config_;
};

} // namespace lte
} // namespace magma
