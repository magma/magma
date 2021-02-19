/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/*****************************************************************************
  Source      amf_app_state_manager.h
  Version     0.1
  Date        2020/10/26
  Product     AMF Core
  Subsystem   Access and Mobility Management Function
  Author      Sanjay Kumar Ojha
  Description Define and  access AMF/NAS state
*****************************************************************************/
#ifndef AMF_APP_STATE_MANAGER_SEEN
#define AMF_APP_STATE_MANAGER_SEEN

#include <sstream>
#include "amf_common_defs.h"
//#include "amf_app_desc.h"
#include "amf_app_defs.h"  //includes amf_app_ue_context.h
using namespace std;

namespace magma5g {
/**
 * Return pointer to the in-memory AMF/NAS state from state manager before
 * processing any message. This is a thread safe call
 * If the read_from_db flag is set to true, the state is loaded from data store
 * before returning the pointer.
 */
amf_app_desc_t* get_amf_nas_state(bool read_from_redis);

/* Retriving respective global hash table*/
hash_table_ts_t* get_amf_ue_state();

int amf_nas_state_init(const amf_config_t* amf_config_p);

/**
 * AmfNasStateManager is a singleton (thread-safe, destruction guaranteed) class
 * that contains functions to maintain Amf and NAS state, i.e. for allocating
 * and freeing state structs, and writing/reading state to db.
 */

class AmfNasStateManager {
 public:
  /**
   * Returns an instance of AmfNasStateManager, guaranteed to be thread safe and
   * initialized only once.
   **/
  static AmfNasStateManager& getInstance();

  // Initialize the local in-memory state when Amf app inits
  int initialize_state(const amf_config_t* amf_config_p);

  /**
   * Retrieve the state pointer from state manager. The read_from_db flag is a
   * debug flag; if set to true, the state is loaded from the data store on
   * every get.
   */
  amf_app_desc_t* get_state(bool read_from_redis);

  /**
   * Release the memory for Amf NAS state and destroy the read-write lock. This
   * is only called when task terminates
   */

  void free_state();

  // int read_ue_state_from_db(); TODO NEED-RECHECK during state machine impl

  /* Retriving respective hash table from global data */
  hash_table_ts_t* get_ue_state_ht();

  /**
   * Copy constructor and assignment operator are marked as deleted functions.
   * Making them public for better debugging/logging.
   */
  AmfNasStateManager(AmfNasStateManager const&) = delete;
  AmfNasStateManager& operator=(AmfNasStateManager const&) = delete;

  // UE state initializemanager flag
  bool persist_state_enabled_;
  bool is_initialized;
  bool state_dirty;
  std::string table_key;
  std::string task_name;
  log_proto_t log_task;
  uint32_t max_ue_htbl_lists_;
  uint32_t amf_statistic_timer_;
  hash_table_ts_t* state_ue_ht;
  amf_app_desc_t* state_cache_p;
  // std::shared_ptr<magma5g::amf_app_desc_t> state_cache_p;

 private:
  // Constructor for Amf NAS state manager
  AmfNasStateManager() {}

  // Destructor for Amf NAS state manager
  ~AmfNasStateManager() {}

  // Initialize state that is non-persistent, e.g. mutex locks and timers
  void amf_nas_state_init_local_state();

  // Create in-memory hashtables for Amf NAS state
  void create_hashtables();

  // Write an empty value to data store, if needed for debugging
  // void clear_db_state();

  /**
   * Initialize memory for Amf state before reading from data-store, the state
   * manager owns the memory allocated for Amf state and frees it when the
   * task terminates
   */
  void create_state();

  // Clean-up the in-memory hashtables
  void clear_amf_nas_hashtables();
};

}  // end namespace magma5g
#endif
