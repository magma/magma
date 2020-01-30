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

extern "C" {
#include "mme_config.h"
}

#include <redis_utils/redis_client.h>
#include <state_manager.h>

#include "mme_app_state_converter.h"
#include "ServiceConfigLoader.h"

namespace magma {
namespace lte {
/**
 * MmeNasStateManager is a singleton (thread-safe, destruction guaranteed) class
 * that contains functions to maintain MME and NAS state, i.e. for allocating
 * and freeing state structs, and writing/reading state to db.
 */

class MmeNasStateManager :
  public StateManager<mme_app_desc_t, MmeNasState, MmeNasStateConverter> {
 public:
  /**
   * Returns an instance of MmeNasStateManager, guaranteed to be thread safe and
   * initialized only once.
   **/
  static MmeNasStateManager& getInstance();

  // Initialize the local in-memory state when MME app inits
  int initialize_state(const mme_config_t* mme_config_p);

  /**
   * Write MME NAS state to redis. This function releases the lock on MME NAS
   * state, which was acquired by get_locked_mme_nas_state.
   */
  void write_state_to_db(mme_app_desc_t** task_state_ptr);

  /**
    * This is a thread-safe call to lock the state and retrieve the pointer to
    * MME Nas state from state manager. The read_from_db flag is a debug flag;
    * if set to true, the state is loaded from the data store on every get.
    */
  mme_app_desc_t* get_locked_mme_nas_state(bool read_from_db);

  /**
    * Retrieve the state pointer from state manager. The read_from_db flag is a
    * debug flag; if set to true, the state is loaded from the data store on
    * every get.
    */
  mme_app_desc_t* get_state(bool read_from_db) override;

  /**
    * Release the memory for MME NAS state and destroy the read-write lock. This
    * is only called when task terminates
    */

  void free_state() override;

  /**
   * Copy constructor and assignment operator are marked as deleted functions.
   * Making them public for better debugging/logging.
   */
  MmeNasStateManager(MmeNasStateManager const&) = delete;
  MmeNasStateManager& operator=(MmeNasStateManager const&) = delete;

 private:
  // Constructor for MME NAS state manager
  MmeNasStateManager();

  // Destructor for MME NAS state manager
  ~MmeNasStateManager();

  // flag to assert if singleton instance has been initialized
  bool is_initialized_;
  mme_app_desc_t* mme_nas_state_p_; // TODO: convert to unique_ptr
  bool persist_state_;
  std::unique_ptr<cpp_redis::client> mme_nas_db_client_;
  int max_ue_htbl_lists_;
  uint32_t mme_statistic_timer_;
  bool mme_nas_state_dirty_; // TODO: convert this to version numbers

  // Initialize state that is non-persistent, e.g. mutex locks and timers
  void mme_nas_state_init_local_state();

  // Create in-memory hashtables for MME NAS state
  void create_hashtables();

  // Write an empty value to data store, if needed for debugging
  void clear_db_state();

  // Acquire lock on the complete MME NAS state
  void lock_mme_nas_state();

  // Release lock on the complete MME NAS state
  void unlock_mme_nas_state();

  /**
   * Initialize memory for MME state before reading from data-store, the state
   * manager owns the memory allocated for MME state and frees it when the
   * task terminates
   */
  void create_state() override;

  // Clean-up the in-memory hashtables
  void clear_mme_nas_hashtables();
};
} // namespace lte
} // namespace magma
