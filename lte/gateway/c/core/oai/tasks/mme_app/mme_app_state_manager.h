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

#pragma once

extern "C" {
#include "mme_config.h"
}

#include <state_manager.h>
#include "mme_app_state_converter.h"
#include "includes/ServiceConfigLoader.h"

namespace magma {
namespace lte {
/**
 * MmeNasStateManager is a singleton (thread-safe, destruction guaranteed) class
 * that contains functions to maintain MME and NAS state, i.e. for allocating
 * and freeing state structs, and writing/reading state to db.
 */

class MmeNasStateManager
    : public StateManager<
          mme_app_desc_t, ue_mm_context_t, oai::MmeNasState, oai::UeContext,
          MmeNasStateConverter> {
 public:
  /**
   * Returns an instance of MmeNasStateManager, guaranteed to be thread safe and
   * initialized only once.
   **/
  static MmeNasStateManager& getInstance();

  // Initialize the local in-memory state when MME app inits
  int initialize_state(const mme_config_t* mme_config_p);

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

  status_code_e read_ue_state_from_db() override;

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

  int max_ue_htbl_lists_;
  uint32_t mme_statistic_timer_;

  // Initialize state that is non-persistent, e.g. mutex locks and timers
  void mme_nas_state_init_local_state();

  // Create in-memory hashtables for MME NAS state
  void create_hashtables();

  // Write an empty value to data store, if needed for debugging
  void clear_db_state();

  /**
   * Initialize memory for MME state before reading from data-store, the state
   * manager owns the memory allocated for MME state and frees it when the
   * task terminates
   */
  void create_state() override;

  // Clean-up the in-memory hashtables
  void clear_mme_nas_hashtables();
};
}  // namespace lte
}  // namespace magma
