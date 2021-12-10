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
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/common/common_defs.h"
#include "include/amf_client_servicer.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_state_manager.h"
#include "lte/gateway/c/core/oai/include/map.h"

namespace magma5g {
/**
 * When the process starts, initialize the in-memory AMF/NAS state and, if
 * persist state flag is set, load it from the data store.
 * This is only done by the amf_app task.
 */
int amf_nas_state_init(const amf_config_t* amf_config_p) {
  return AmfNasStateManager::getInstance().initialize_state(amf_config_p);
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

/**
 * Write the AMF/NAS state to data store after processing any message. This is
 * a thread safe call
 */
void put_amf_nas_state() {
  magma5g::AmfNasStateManager::getInstance().write_state_to_db();
}

/**
 * Release the memory allocated for the AMF NAS state, this does not clean the
 * state persisted in data store
 */
void clear_amf_nas_state() {
  AmfNasStateManager::getInstance().free_state();
}

map_uint64_ue_context_t get_amf_ue_state() {
  return AmfNasStateManager::getInstance().get_ue_state_map();
}

void delete_amf_ue_state(imsi64_t imsi64) {
  auto imsi_str =
      magma5g::AmfNasStateManager::getInstance().get_imsi_str(imsi64);
  magma5g::AmfNasStateManager::getInstance().clear_ue_state_db(imsi_str);
}

/**
 * Getter function for singleton instance of the AmfNasStateManager class,
 * guaranteed to be thread-safe and initialized only once
 */
AmfNasStateManager& AmfNasStateManager::getInstance() {
  static AmfNasStateManager instance;
  return instance;
}

// Constructor for AMF NAS state object
AmfNasStateManager::AmfNasStateManager()
    : max_ue_htbl_lists_(NUM_MAX_UE_HTBL_LISTS) {}

// Destructor for AMF NAS state object
AmfNasStateManager::~AmfNasStateManager() {
  free_state();
}

// Singleton class initializer which calls to create new object of
// AmfNasStateManager
int AmfNasStateManager::initialize_state(const amf_config_t* amf_config_p) {
  uint32_t rc           = RETURNok;
  persist_state_enabled = amf_config_p->use_stateless;
  max_ue_htbl_lists_    = amf_config_p->max_ues;
  amf_statistic_timer_  = amf_config_p->amf_statistic_timer;
  log_task              = LOG_AMF_APP;
  task_name             = AMF_TASK_NAME;
  table_key             = AMF_NAS_STATE_KEY;

  // Allocate the local AMF state and create respective single object
  create_state();
#if MME_UNIT_TEST
  read_state_from_db();
#endif
  is_initialized = true;
  return rc;
}

// Create an object of AmfNasStateManager and Initialize memory
// for AMF state before doing any operation from data store
void AmfNasStateManager::create_state() {
  state_cache_p                               = new (amf_app_desc_t);
  state_cache_p->amf_app_ue_ngap_id_generator = 1;
  state_cache_p->amf_ue_contexts.imsi_amf_ue_id_htbl.set_name(
      AMF_IMSI_UE_ID_TABLE_NAME);
  state_cache_p->amf_ue_contexts.tun11_ue_context_htbl.set_name(
      AMF_TUN_UE_ID_TABLE_NAME);
  state_cache_p->amf_ue_contexts.gnb_ue_ngap_id_ue_context_htbl.set_name(
      AMF_GNB_UE_ID_AMF_UE_ID_TABLE_NAME);
  state_cache_p->amf_ue_contexts.guti_ue_context_htbl.set_name(
      AMF_GUTI_UE_ID_TABLE_NAME);
  state_ue_map.set_name(AMF_UE_ID_UE_CTXT_TABLE_NAME);

  // Initialize the local timers, which are non-persistent
  amf_nas_state_init_local_state();
}

// Free the memory allocated to state pointer
void AmfNasStateManager::free_state() {
  if (!state_cache_p) {
    return;
  }
  delete state_cache_p;
  state_cache_p = nullptr;
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
#if MME_UNIT_TEST
  if (persist_state_enabled && read_from_redis) {
    read_state_from_db();
  }
#endif
  return state_cache_p;
}

map_uint64_ue_context_t AmfNasStateManager::get_ue_state_map() {
  return state_ue_map;
}

status_code_e AmfNasStateManager::read_state_from_db() {
#if !MME_UNIT_TEST
  StateManager::read_state_from_db();
#else
  if (persist_state_enabled) {
    magma::lte::oai::MmeNasState state_proto = magma::lte::oai::MmeNasState();
    std::string proto_str;
    // Reads from the AmfClientServicer DataStore Map(map_tableKey_protoStr)
    if (AMFClientServicer::getInstance().map_tableKey_protoStr.get(
            table_key, &proto_str) != magma::MAP_OK) {
      OAILOG_DEBUG(LOG_MME_APP, "Failed to read proto from db \n");
      return RETURNerror;
    }
    // Deserialization Step
    if (!state_proto.ParseFromString(proto_str)) {
      return RETURNerror;
    }
    AmfNasStateConverter::proto_to_state(state_proto, state_cache_p);
  }
#endif
  return RETURNok;
}

void AmfNasStateManager::write_state_to_db() {
#if !MME_UNIT_TEST
  StateManager::write_state_to_db();
#else
  if (persist_state_enabled) {
    magma::lte::oai::MmeNasState state_proto = magma::lte::oai::MmeNasState();
    AmfNasStateConverter::state_to_proto(state_cache_p, &state_proto);
    std::string proto_str;
    redis_client->serialize(state_proto, proto_str);
    std::size_t new_hash = std::hash<std::string>{}(proto_str);
    if (new_hash != this->task_state_hash) {
      // Writes to the AmfClientServicer DataStore Map(map_tableKey_protoStr)
      if (AMFClientServicer::getInstance().map_tableKey_protoStr.insert(
              table_key, proto_str) != magma::MAP_OK) {
        OAILOG_ERROR(log_task, "Failed to write state to db");
        return;
      }
      OAILOG_DEBUG(log_task, "Finished writing state");
      this->task_state_version++;
      this->state_dirty     = false;
      this->task_state_hash = new_hash;
    }
  }
#endif
}

}  // namespace magma5g
