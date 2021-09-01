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

#ifdef __cplusplus
extern "C" {
#endif
#include "log.h"
#include "dynamic_memory_check.h"
#ifdef __cplusplus
}
#endif
#include "common_defs.h"
#include "amf_app_state_manager.h"

namespace magma5g {
constexpr char AMF_NAS_STATE_KEY[] = "amf_nas_state";
constexpr char AMF_UE_ID_UE_CTXT_TABLE_NAME[] =
    "amf_app_amf_ue_ngap_id_ue_context_htbl";
constexpr char AMF_IMSI_UE_ID_TABLE_NAME[] = "amf_app_imsi_ue_context_htbl";
constexpr char AMF_TUN_UE_ID_TABLE_NAME[]  = "amf_app_tun11_ue_context_htbl";
constexpr char AMF_GUTI_UE_ID_TABLE_NAME[] = "amf_app_tun11_ue_context_htbl";
constexpr char AMF_GNB_UE_ID_AMF_UE_ID_TABLE_NAME[] =
    "anf_app_gnb_ue_ngap_id_ue_context_htbl";
constexpr char AMF_TASK_NAME[] = "AMF";

/*hash function similar to default to initialize during hash table
 * initialization*/
static hash_size_t amf_def_hashfunc(const uint64_t keyP) {
  return (hash_size_t) keyP;
}

/**
 * Getter function for singleton instance of the AmfNasStateManager class,
 * guaranteed to be thread-safe and initialized only once
 */
AmfNasStateManager& AmfNasStateManager::getInstance() {
  static AmfNasStateManager instance;
  return instance;
}

// Singleton class initializer which calls to create new object of
// AmfNasStateManager
int AmfNasStateManager::initialize_state(const amf_config_t* amf_config_p) {
  uint32_t rc            = RETURNok;
  persist_state_enabled_ = amf_config_p->use_stateless;
  max_ue_htbl_lists_     = amf_config_p->max_ues;
  amf_statistic_timer_   = amf_config_p->amf_statistic_timer;
  log_task               = LOG_AMF_APP;
  task_name              = AMF_TASK_NAME;
  table_key              = AMF_NAS_STATE_KEY;

  // Allocate the local AMF state and create respective single object
  create_state();

  is_initialized = true;
  return rc;
}

/**
 * When the process starts, initialize the in-memory AMF/NAS state and, if
 * persist state flag is set, load it from the data store.
 * This is only done by the mme_app task.
 */
int amf_nas_state_init(const amf_config_t* amf_config_p) {
  return AmfNasStateManager::getInstance().initialize_state(amf_config_p);
}

// Create an object of AmfNasStateManager and Initialize memory
// for AMF state before doing any operation from data store
void AmfNasStateManager::create_state() {
  state_cache_p                               = new (amf_app_desc_t);
  state_cache_p->amf_app_ue_ngap_id_generator = 1;
  create_hashtables();

  // Initialize the local timers, which are non-persistent
  amf_nas_state_init_local_state();
}

// Create the hashtables for AMF and NAS state
void AmfNasStateManager::create_hashtables() {
  bstring b          = bfromcstr(AMF_IMSI_UE_ID_TABLE_NAME);
  max_ue_htbl_lists_ = 2;
  state_cache_p->amf_ue_contexts.imsi_amf_ue_id_htbl =
      hashtable_uint64_ts_create(max_ue_htbl_lists_, nullptr, b);
  btrunc(b, 0);
  bassigncstr(b, AMF_TUN_UE_ID_TABLE_NAME);
  state_cache_p->amf_ue_contexts.tun11_ue_context_htbl =
      hashtable_uint64_ts_create(max_ue_htbl_lists_, nullptr, b);
  btrunc(b, 0);
  bassigncstr(b, AMF_UE_ID_UE_CTXT_TABLE_NAME);
  state_ue_ht = hashtable_ts_create(
      max_ue_htbl_lists_, nullptr, amf_app_state_free_ue_context, b);
  btrunc(b, 0);
  bassigncstr(b, AMF_GNB_UE_ID_AMF_UE_ID_TABLE_NAME);
  state_cache_p->amf_ue_contexts.gnb_ue_ngap_id_ue_context_htbl =
      hashtable_uint64_ts_create(max_ue_htbl_lists_, amf_def_hashfunc, b);
  btrunc(b, 0);
  bassigncstr(b, AMF_GUTI_UE_ID_TABLE_NAME);
  state_cache_p->amf_ue_contexts.guti_ue_context_htbl =
      obj_hashtable_uint64_ts_create(max_ue_htbl_lists_, nullptr, nullptr, b);
  bdestroy_wrapper(&b);
}

// Initialize state that is non-persistent, e.g. timers
void AmfNasStateManager::amf_nas_state_init_local_state() {
  // create statistic timer locally
  state_cache_p->m5_statistic_timer_period = amf_statistic_timer_;
  state_cache_p->m5_statistic_timer_id     = 0;
}

/**
 * Getter function to get the pointer to the in-memory user state. The
 * read_from_db flag is a debug flag to force read from data store instead of
 * just returning the pointer. In non-debug mode, the state is only read from
 * data store when initialize_state is called and get_state just returns the
 * pointer to amf_app_desc_t structure.
 */
amf_app_desc_t* AmfNasStateManager::get_state(bool read_from_redis) {
  state_dirty = true;

  // if read_from_redis is false, no need to clear and create ht.
  if (persist_state_enabled_ && read_from_redis) {
    clear_amf_nas_hashtables();
    create_hashtables();
  }
  return state_cache_p;
}

// Delete the hashtables for AMF NAS state
// TODO in future PR, Hash table is replaced by MAP & hash table is depricated
void AmfNasStateManager::clear_amf_nas_hashtables() {
  if (!state_cache_p) {
    return;
  }
  hashtable_ts_destroy(state_ue_ht);
  hashtable_uint64_ts_destroy(
      state_cache_p->amf_ue_contexts.imsi_amf_ue_id_htbl);
  hashtable_uint64_ts_destroy(
      state_cache_p->amf_ue_contexts.tun11_ue_context_htbl);
  hashtable_uint64_ts_destroy(
      state_cache_p->amf_ue_contexts.gnb_ue_ngap_id_ue_context_htbl);
  obj_hashtable_uint64_ts_destroy(
      state_cache_p->amf_ue_contexts.guti_ue_context_htbl);
}

hash_table_ts_t* AmfNasStateManager::get_ue_state_ht() {
  return state_ue_ht;
}

hash_table_ts_t* get_amf_ue_state() {
  return AmfNasStateManager::getInstance().get_ue_state_ht();
}

/**
 * Return pointer to the in-memory Amf/NAS state from state manager before
 * processing any message. This is a thread safe call
 * If the read_from_db flag is set to true, the state is loaded from data store
 * before returning the pointer.
 */
amf_app_desc_t* get_amf_nas_state(bool read_from_redis) {
  return AmfNasStateManager::getInstance().get_state(read_from_redis);
}
}  // namespace magma5g
