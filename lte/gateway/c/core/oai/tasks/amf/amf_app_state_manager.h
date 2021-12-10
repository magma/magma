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

#pragma once
#include <sstream>
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/assertions.h"
#include "hashtable.h"
#include "obj_hashtable.h"
#ifdef __cplusplus
}
#endif

#include "amf_smfDefs.h"
#include "amf_app_defs.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_state_converter.h"
#include <lte/gateway/c/core/oai/include/state_manager.h>

namespace magma5g {
constexpr char AMF_NAS_STATE_KEY[] = "amf_nas_state";
constexpr char AMF_UE_ID_UE_CTXT_TABLE_NAME[] =
    "amf_app_amf_ue_ngap_id_ue_context_htbl";
constexpr char AMF_IMSI_UE_ID_TABLE_NAME[] = "amf_app_imsi_ue_context_htbl";
constexpr char AMF_TUN_UE_ID_TABLE_NAME[]  = "amf_app_tun11_ue_context_htbl";
constexpr char AMF_GUTI_UE_ID_TABLE_NAME[] = "amf_app_tun11_ue_context_htbl";
constexpr char AMF_GNB_UE_ID_AMF_UE_ID_TABLE_NAME[] =
    "anf_app_gnb_ue_ngap_id_ue_context_htbl";
constexpr char AMF_TASK_NAME[]  = "AMF";
const int NUM_MAX_UE_HTBL_LISTS = 6;
/**
 * When the process starts, initialize the in-memory AMF+NAS state and, if
 * persist state flag is set, load it from the data store.
 * This is only done by the amf_app task.
 */
int amf_nas_state_init(const amf_config_t* amf_config_p);

/**
 * Return pointer to the in-memory AMF/NAS state from state manager before
 * processing any message. This is a thread safe call
 * If the read_from_db flag is set to true, the state is loaded from data store
 * before returning the pointer.
 */
amf_app_desc_t* get_amf_nas_state(bool read_from_redis);

/**
 * Write the AMF/NAS state to data store after processing any message. This is a
 * thread safe call
 */
void put_amf_nas_state();

/**
 * Release the memory allocated for the AMF NAS state, this does not clean the
 * state persisted in data store
 */
void clear_amf_nas_state();

// Retrieving respective global hash table
map_uint64_ue_context_t get_amf_ue_state();

// Persists UE AMF state for subscriber into db
void put_amf_ue_state(
    magma5g::amf_app_desc_t* amf_app_desc_p, imsi64_t imsi64,
    bool force_ue_write);
// Deletes entry for UE AMF state on db
void delete_amf_ue_state(imsi64_t imsi64);

/**
 * AmfNasStateManager is a singleton (thread-safe, destruction guaranteed) class
 * that contains functions to maintain Amf and NAS state, i.e. for allocating
 * and freeing state structs, and writing/reading state to db.
 */
class AmfNasStateManager
    : public magma::lte::StateManager<
          amf_app_desc_t, ue_m5gmm_context_t, magma::lte::oai::MmeNasState,
          magma::lte::oai::UeContext, AmfNasStateConverter> {
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
  amf_app_desc_t* get_state(bool read_from_redis) override;

  // Retriving respective hash table from global data
  map_uint64_ue_context_t get_ue_state_map();

  /**
   * Copy constructor and assignment operator are marked as deleted functions.
   * Making them public for better debugging/logging.
   */
  AmfNasStateManager(AmfNasStateManager const&) = delete;
  AmfNasStateManager& operator=(AmfNasStateManager const&) = delete;

  // AMF state initializemanager flag
  uint32_t max_ue_htbl_lists_;
  uint32_t amf_statistic_timer_;
  map_uint64_ue_context_t state_ue_map;

  void free_state() override;

  void write_state_to_db() override;
  status_code_e read_state_from_db() override;

 private:
  AmfNasStateManager();
  ~AmfNasStateManager();

  // Initialize state that is non-persistent, e.g. mutex locks and timers
  void amf_nas_state_init_local_state();

  // Create in-memory hashtables for Amf NAS state
  void create_hashtables();

  /**
   * Initialize memory for Amf state before reading from data-store, the state
   * manager owns the memory allocated for Amf state and frees it when the
   * task terminates
   */
  void create_state() override;
};
}  // namespace magma5g
