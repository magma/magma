/*
 *
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

#include "mme_config.h"
#include "s1ap_types.h"

#ifdef __cplusplus
}
#endif

#include "state_manager.h"
#include "s1ap_state_converter.h"

namespace {
constexpr char S1AP_STATE_TABLE[] = "s1ap_state";
}

namespace magma {
namespace lte {

/**
 * S1apStateManager is a thread safe singleton class that contains functions
 * to maintain S1AP task state, allocating and freeing related state structs.
 */
class S1apStateManager :
  public StateManager<
    s1ap_state_t,
    ue_description_t,
    magma::lte::gateway::s1ap::S1apState,
    magma::lte::gateway::s1ap::UeDescription,
    S1apStateConverter> {
 public:
  /**
   * Returns an instance of S1apStateManager, guaranteed to be thread safe and
   * initialized only once.
   * @return S1apStateManager instance
   */
  static S1apStateManager& getInstance();

  /**
   * Function to initialize member variables
   * @param mme_config mme_config_t struct
   */
  void init(uint32_t max_ues, uint32_t max_enbs, bool use_stateless);

  // Copy constructor and assignment operator are marked as deleted functions
  S1apStateManager(S1apStateManager const&) = delete;
  S1apStateManager& operator=(S1apStateManager const&) = delete;

  /**
   * Frees all memory allocated on s1ap_state cache struct
   */
  void free_state() override;

  /**
   * Serializes s1ap_imsi_map to proto and saves it into data store
   */
  void put_s1ap_imsi_map();

  /**
   * Returns a pointer to s1ap_imsi_map
   */
  s1ap_imsi_map_t* get_s1ap_imsi_map();

 private:
  S1apStateManager();
  ~S1apStateManager();

  /**
   * Allocates new s1ap_state_t struct and its properties
   */
  void create_state() override;

  void create_s1ap_imsi_map();
  void clear_s1ap_imsi_map();

  uint32_t max_ues_;
  uint32_t max_enbs_;
  s1ap_imsi_map_t* s1ap_imsi_map_;
};
} // namespace lte
} // namespace magma
