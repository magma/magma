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
#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/include/map.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_state_manager.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/include/amf_client_servicer.hpp"

using magma::lte::oai::MmeNasState;
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
  OAILOG_FUNC_IN(LOG_AMF_APP);
  magma5g::AmfNasStateManager::getInstance().write_state_to_db();
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

/**
 * Release the memory allocated for the AMF NAS state, this does not clean the
 * state persisted in data store
 */
void clear_amf_nas_state() {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  AmfNasStateManager::getInstance().free_state();
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

map_uint64_ue_context_t* get_amf_ue_state() {
  return AmfNasStateManager::getInstance().get_ue_state_map();
}

void delete_amf_ue_state(imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  OAILOG_DEBUG(LOG_AMF_APP, "Delete AMF ue state, %lu", imsi64);
#if !MME_UNIT_TEST
  /* Data store is Redis db. In this case entry is removed from Redis db */
  auto imsi_str = AmfNasStateManager::getInstance().get_imsi_str(imsi64);
  AmfNasStateManager::getInstance().clear_ue_state_db(imsi_str);
#else
  /* Data store is a map defined in AmfClientServicer.In this case entry is
   * removed from map_imsi_ue_proto_str */
  auto imsi_str = AmfNasStateManager::getInstance().get_imsi_str(imsi64);
  std::string key = IMSI_PREFIX + imsi_str + ":" + AMF_TASK_NAME;
  AMFClientServicer::getInstance().map_imsi_ue_proto_str.remove(key);
#endif
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

/**
 * Getter function for singleton instance of the AmfNasStateManager class,
 * guaranteed to be thread-safe and initialized only once
 */
AmfNasStateManager& AmfNasStateManager::getInstance() {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  static AmfNasStateManager instance;
  OAILOG_FUNC_RETURN(LOG_AMF_APP, instance);
}

// Constructor for AMF NAS state object
AmfNasStateManager::AmfNasStateManager()
    : max_ue_htbl_lists_(NUM_MAX_UE_HTBL_LISTS) {}

// Destructor for AMF NAS state object
AmfNasStateManager::~AmfNasStateManager() { free_state(); }

// Singleton class initializer which calls to create new object of
// AmfNasStateManager
status_code_e AmfNasStateManager::initialize_state(
    const amf_config_t* amf_config_p) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  status_code_e rc = RETURNok;
  persist_state_enabled = amf_config_p->use_stateless;
  max_ue_htbl_lists_ = amf_config_p->max_ues;
  amf_statistic_timer_ = amf_config_p->amf_statistic_timer;
  log_task = LOG_AMF_APP;
  task_name = AMF_TASK_NAME;
  table_key = AMF_NAS_STATE_KEY;

  // Allocate the local AMF state and create respective single object
  create_state();
#if !MME_UNIT_TEST
  redis_client =
      std::make_unique<magma::lte::RedisClient>(persist_state_enabled);
#endif
  read_state_from_db();
  read_ue_state_from_db();
  is_initialized = true;
  OAILOG_FUNC_RETURN(LOG_AMF_APP, rc);
}

// Create an object of AmfNasStateManager and Initialize memory
// for AMF state before doing any operation from data store
void AmfNasStateManager::create_state() {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  state_cache_p = new (amf_app_desc_t);
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
  state_cache_p->nb_ue_connected = 0;
  state_cache_p->nb_ue_attached = 0;
  state_cache_p->nb_pdu_sessions = 0;
  state_cache_p->nb_ue_idle = 0;
  // Initialize the local timers, which are non-persistent
  amf_nas_state_init_local_state();
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

// Free the memory allocated to state pointer
void AmfNasStateManager::free_state() {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  if (!state_cache_p) {
    OAILOG_ERROR(LOG_AMF_APP, "state_cache_p is NULL");
    OAILOG_FUNC_OUT(LOG_AMF_APP);
  }
  state_ue_map.umap.clear();
  delete state_cache_p;
  state_cache_p = nullptr;
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

// Initialize state that is non-persistent, e.g. timers
void AmfNasStateManager::amf_nas_state_init_local_state() {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  // create statistic timer locally
  state_cache_p->m5_statistic_timer_period = amf_statistic_timer_;
  state_cache_p->m5_statistic_timer_id = 0;
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

/**
 * Getter function to get the pointer to the in-memory user state. The
 * read_from_db flag is a debug flag to force read from data store instead of
 * just returning the pointer. In non-debug mode, the state is only read from
 * data store when initialize_state is called and get_state just returns the
 * pointer to amf_app_desc_t structure.
 */
amf_app_desc_t* AmfNasStateManager::get_state(bool read_from_redis) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  state_dirty = true;
  if (persist_state_enabled && read_from_redis) {
    read_state_from_db();
    read_ue_state_from_db();
  }
  OAILOG_FUNC_RETURN(LOG_AMF_APP, state_cache_p);
}

map_uint64_ue_context_t* AmfNasStateManager::get_ue_state_map() {
  return &state_ue_map;
}

// This is a helper function for debugging. If the state manager needs to clear
// the state in the data store, it can call this function to delete the key.
void AmfNasStateManager::clear_db_state() {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  OAILOG_DEBUG(LOG_AMF_APP, "Clearing state in data store");
  std::vector<std::string> keys_to_del;
  keys_to_del.emplace_back(AMF_NAS_STATE_KEY);

  if (redis_client->clear_keys(keys_to_del) != RETURNok) {
    OAILOG_ERROR(LOG_AMF_APP, "Failed to clear the state in data store");
  }
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void put_amf_ue_state(amf_app_desc_t* amf_app_desc_p, imsi64_t imsi64,
                      bool force_ue_write) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  if ((!AmfNasStateManager::getInstance().is_persist_state_enabled()) ||
      (imsi64 == INVALID_IMSI64)) {
    OAILOG_FUNC_OUT(LOG_AMF_APP);
  }
  ue_m5gmm_context_t* ue_context_p = nullptr;
  amf_ue_ngap_id_t ue_id;
  get_amf_ue_id_from_imsi(&amf_app_desc_p->amf_ue_contexts, imsi64, &ue_id);
  ue_context_p = amf_ue_context_exists_amf_ue_ngap_id((amf_ue_ngap_id_t)ue_id);
  // Only write MME UE state to redis if force flag is set or UE is in EMM
  // Registered state
  if ((ue_context_p && force_ue_write) ||
      (ue_context_p && ue_context_p->mm_state == REGISTERED_CONNECTED)) {
    auto imsi_str = AmfNasStateManager::getInstance().get_imsi_str(imsi64);
    AmfNasStateManager::getInstance().write_ue_state_to_db(ue_context_p,
                                                           imsi_str);
  }
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateManager::write_ue_state_to_db(
    const ue_m5gmm_context_t* ue_context, const std::string& imsi_str) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
#if !MME_UNIT_TEST
  /* Data store is Redis db. In this case actual call is made to Redis db */
  StateManager::write_ue_state_to_db(ue_context, imsi_str);
#else
  /* Data store is a map defined in AmfClientServicer.In this case call is NOT
   * made to Redis db */
  std::string proto_str;
  magma::lte::oai::UeContext ue_proto = magma::lte::oai::UeContext();
  AmfNasStateConverter::ue_to_proto(ue_context, &ue_proto);
  redis_client->serialize(ue_proto, proto_str);
  std::size_t new_hash = std::hash<std::string>{}(proto_str);

  if (new_hash != this->ue_state_hash[imsi_str]) {
    std::string key = IMSI_PREFIX + imsi_str + ":" + task_name;
    if (AMFClientServicer::getInstance().map_imsi_ue_proto_str.insert(
            key, proto_str) != magma::MAP_OK) {
      OAILOG_ERROR(log_task, "Failed to write UE state to db for IMSI %s",
                   imsi_str.c_str());
      OAILOG_FUNC_OUT(LOG_AMF_APP);
    }

    this->ue_state_version[imsi_str]++;
    this->ue_state_hash[imsi_str] = new_hash;
    OAILOG_DEBUG(log_task, "Finished writing UE state for IMSI %s",
                 imsi_str.c_str());
  }
#endif
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

status_code_e AmfNasStateManager::read_ue_state_from_db() {
  OAILOG_FUNC_IN(LOG_AMF_APP);
#if !MME_UNIT_TEST
  /* Data store is Redis db. In this case actual call is made to Redis db */
  if (!persist_state_enabled) {
    OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNok);
  }
  auto keys = redis_client->get_keys("IMSI*" + task_name + "*");
  for (const auto& key : keys) {
    magma::lte::oai::UeContext ue_proto = magma::lte::oai::UeContext();
    if (redis_client->read_proto(key.c_str(), ue_proto) != RETURNok) {
      OAILOG_ERROR(log_task, "Failed to read UE state from db for %s",
                   key.c_str());
      OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
    }

    // Update each UE state version from redis
    this->ue_state_version[key] = redis_client->read_version(table_key);
    ue_m5gmm_context_t* ue_context_p = new ue_m5gmm_context_t();
    AmfNasStateConverter::proto_to_ue(ue_proto, ue_context_p);
    state_ue_map.insert(ue_context_p->amf_ue_ngap_id, ue_context_p);
    OAILOG_DEBUG(log_task, "Reading UE state from db for %s", key.c_str());
  }
  OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNok);
#else
  /* Data store is a map defined in AmfClientServicer.In this case call is NOT
   * made to Redis db */
  if (!persist_state_enabled) {
    OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNok);
  }
  for (const auto& kv :
       AMFClientServicer::getInstance().map_imsi_ue_proto_str.umap) {
    magma::lte::oai::UeContext ue_proto = magma::lte::oai::UeContext();
    if (!ue_proto.ParseFromString(kv.second)) {
      OAILOG_ERROR(log_task, "Failed to parse proto from string");
      OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
    }
    ue_m5gmm_context_t* ue_context_p = new ue_m5gmm_context_t();
    AmfNasStateConverter::proto_to_ue(ue_proto, ue_context_p);
    state_ue_map.insert(ue_context_p->amf_ue_ngap_id, ue_context_p);
    OAILOG_DEBUG(log_task, "Reading UE state from db for %s", kv.first.c_str());
  }
  OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNok);
#endif
}

status_code_e AmfNasStateManager::read_state_from_db() {
  OAILOG_FUNC_IN(LOG_AMF_APP);
#if !MME_UNIT_TEST
  /* Data store is Redis db. In this case actual call is made to Redis db */
  StateManager::read_state_from_db();
#else
  /* Data store is a map defined in AmfClientServicer.In this case call is NOT
   * made to Redis db */
  if (persist_state_enabled) {
    MmeNasState state_proto = MmeNasState();
    std::string proto_str;
    // Reads from the AmfClientServicer DataStore Map(map_table_key_proto_str)
    if (AMFClientServicer::getInstance().map_table_key_proto_str.get(
            table_key, &proto_str) != magma::MAP_OK) {
      OAILOG_DEBUG(LOG_AMF_APP, "Failed to read proto from db \n");
      OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
    }
    // Deserialization Step
    if (!state_proto.ParseFromString(proto_str)) {
      OAILOG_ERROR(log_task, "Failed to parse proto from string");
      OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
    }
    AmfNasStateConverter::proto_to_state(state_proto, state_cache_p);
  }
#endif
  OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNok);
}

void AmfNasStateManager::write_state_to_db() {
  OAILOG_FUNC_IN(LOG_AMF_APP);
#if !MME_UNIT_TEST
  /* Data store is Redis db. In this case actual call is made to Redis db */
  StateManager::write_state_to_db();
#else
  /* Data store is a map defined in AmfClientServicer.In this case call is NOT
   * made to Redis db */
  if (persist_state_enabled) {
    MmeNasState state_proto = MmeNasState();
    AmfNasStateConverter::state_to_proto(state_cache_p, &state_proto);
    std::string proto_str;
    redis_client->serialize(state_proto, proto_str);
    std::size_t new_hash = std::hash<std::string>{}(proto_str);
    if (new_hash != this->task_state_hash) {
      // Writes to the AmfClientServicer DataStore Map(map_table_key_proto_str)
      if (AMFClientServicer::getInstance().map_table_key_proto_str.insert(
              table_key, proto_str) != magma::MAP_OK) {
        OAILOG_ERROR(log_task, "Failed to write state to db");
        OAILOG_FUNC_IN(LOG_AMF_APP);
      }
      OAILOG_DEBUG(log_task, "Finished writing state");
      this->task_state_version++;
      this->state_dirty = false;
      this->task_state_hash = new_hash;
    }
  }
#endif
  OAILOG_FUNC_IN(LOG_AMF_APP);
}

}  // namespace magma5g
