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

#include "lte/gateway/c/core/oai/tasks/sgw/spgw_state_manager.hpp"
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/common/dynamic_memory_check.h"

namespace magma {
namespace lte {

SpgwStateManager::SpgwStateManager() : config_(nullptr) {}

SpgwStateManager::~SpgwStateManager() { free_state(); }

SpgwStateManager& SpgwStateManager::getInstance() {
  static SpgwStateManager instance;
  return instance;
}

void SpgwStateManager::init(bool persist_state, const spgw_config_t* config) {
  log_task = LOG_SPGW_APP;
  task_name = SPGW_TASK_NAME;
  table_key = SPGW_STATE_TABLE_NAME;
  persist_state_enabled = persist_state;
  config_ = config;
  redis_client = std::make_unique<RedisClient>(persist_state);
  create_state();
  if (read_state_from_db() != RETURNok) {
    OAILOG_ERROR(LOG_SPGW_APP, "Failed to read state from redis");
  }
  is_initialized = true;
}

void SpgwStateManager::create_state() {
  // Allocating spgw_state_p
  state_cache_p = new oai::SpgwState();

  state_teid_map.map =
      new google::protobuf::Map<uint32_t, magma::lte::oai::S11BearerContext*>();
  state_teid_map.set_name(S11_BEARER_CONTEXT_INFO_HT_NAME);

  state_ue_map.map =
      new google::protobuf::Map<uint64_t, magma::lte::oai::SpgwUeContext*>();
  state_ue_map.set_name(SPGW_STATE_UE_MAP);
  state_ue_map.bind_callback(sgw_free_ue_context);
  char ip_str[INET_ADDRSTRLEN];
  inet_ntop(AF_INET, &(config_->sgw_config.ipv4.S1u_S12_S4_up.s_addr), ip_str,
            INET_ADDRSTRLEN);
  state_cache_p->mutable_sgw_s1u_ip_addr()->set_ipv4_addr(ip_str);

  char ip6_str[INET6_ADDRSTRLEN];
  inet_ntop(AF_INET6, &(config_->sgw_config.ipv6.S1u_S12_S4_up), ip6_str,
            INET6_ADDRSTRLEN);
  state_cache_p->mutable_sgw_s1u_ip_addr()->set_ipv6_addr(ip6_str);

  state_cache_p->mutable_gtpv1u_data()->mutable_sgw_s1u_ip_addr()->MergeFrom(
      state_cache_p->sgw_s1u_ip_addr());

  state_cache_p->set_gtpv1u_teid(0);
}

void SpgwStateManager::free_state() {
  AssertFatal(
      is_initialized,
      "SpgwStateManager init() function should be called to initialize state.");

  if (state_cache_p == nullptr) {
    return;
  }

  if (state_ue_map.map && state_ue_map.destroy_map() != PROTO_MAP_OK) {
    OAI_FPRINTF_ERR("An error occurred while destroying SPGW's state ue map ");
  }

  if (state_teid_map.destroy_map() != magma::PROTO_MAP_OK) {
    OAI_FPRINTF_ERR("An error occurred while destroying state_teid_map");
  }
  state_cache_p->Clear();
  free_cpp_wrapper((void**)&state_cache_p);
}

status_code_e SpgwStateManager::read_ue_state_from_db() {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  if (!persist_state_enabled) {
    return RETURNok;
  }
  state_teid_map_t* state_teid_map = get_spgw_teid_state();
  if (!state_teid_map) {
    OAILOG_ERROR(LOG_SPGW_APP, "Failed to get state_teid_map");
    return RETURNerror;
  }

  auto keys = redis_client->get_keys("IMSI*" + task_name + "*");
  for (const auto& key : keys) {
    oai::SpgwUeContext ue_proto = oai::SpgwUeContext();
    if (redis_client->read_proto(key, ue_proto) != RETURNok) {
      return RETURNerror;
    }
    OAILOG_DEBUG(log_task, "Reading UE state from db for key %s", key.c_str());
    oai::SpgwUeContext* ue_context_p = new oai::SpgwUeContext();
    // Update each UE state version from redis
    this->ue_state_version[key] = redis_client->read_version(table_key);

    ue_context_p->MergeFrom(ue_proto);

    state_ue_map.insert(get_imsi_from_key(key), ue_context_p);
    for (uint8_t idx = 0; idx < ue_context_p->s11_bearer_context_size();
         idx++) {
      state_teid_map->insert(ue_context_p->s11_bearer_context(idx)
                                 .sgw_eps_bearer_context()
                                 .sgw_teid_s11_s4(),
                             ue_context_p->mutable_s11_bearer_context(idx));
    }
  }
  return RETURNok;
}

state_teid_map_t* SpgwStateManager::get_state_teid_map() {
  AssertFatal(
      is_initialized,
      "StateManager init() function should be called to initialize state");

  return &state_teid_map;
}

map_uint64_spgw_ue_context_t* SpgwStateManager::get_spgw_ue_state_map() {
  return &state_ue_map;
}

void SpgwStateManager::write_ue_state_to_db(
    const oai::SpgwUeContext* ue_context, const std::string& imsi_str) {
  AssertFatal(
      is_initialized,
      "StateManager init() function should be called to initialize state");

  std::string proto_str;
  redis_client->serialize(*ue_context, proto_str);
  std::size_t new_hash = std::hash<std::string>{}(proto_str);
  if (new_hash != this->ue_state_hash[imsi_str]) {
    std::string key = IMSI_STR_PREFIX + imsi_str + ":" + task_name;
    if (redis_client->write_proto_str(key, proto_str,
                                      ue_state_version[imsi_str]) != RETURNok) {
      OAILOG_ERROR(LOG_SPGW_APP, "Failed to write UE state to db for IMSI %s",
                   imsi_str.c_str());
      return;
    }
    this->ue_state_version[imsi_str]++;
    this->ue_state_hash[imsi_str] = new_hash;
    OAILOG_DEBUG(log_task, "Finished writing UE state for IMSI %s",
                 imsi_str.c_str());
  }
}

void SpgwStateManager::write_spgw_state_to_db(void) {
  AssertFatal(
      is_initialized,
      "Spgw StateManager init() function should be called to initialize state");

  if (!state_dirty) {
    OAILOG_ERROR(log_task, "Tried to put state while it was not in use");
    return;
  }

  if (persist_state_enabled) {
    std::string proto_str;
    redis_client->serialize(*state_cache_p, proto_str);
    std::size_t new_hash = std::hash<std::string>{}(proto_str);

    if (new_hash != this->task_state_hash) {
      if (redis_client->write_proto_str(table_key, proto_str,
                                        this->task_state_version) != RETURNok) {
        OAILOG_ERROR(log_task, "Failed to write state to db");
        return;
      }
      OAILOG_DEBUG(log_task, "Finished writing state");
      this->task_state_version++;
      this->state_dirty = false;
      this->task_state_hash = new_hash;
    }
  }
}

status_code_e SpgwStateManager::read_state_from_db() {
#if !MME_UNIT_TEST
  if (persist_state_enabled) {
    state_cache_p->Clear();
    if (redis_client->read_proto(table_key, *state_cache_p) != RETURNok) {
      OAILOG_ERROR(LOG_MME_APP, "Failed to read state info from db \n");
      return RETURNerror;
    }

    // Update the state version from redis
    this->task_state_version = redis_client->read_version(table_key);
  }
#endif
  return RETURNok;
}

oai::SpgwState* SpgwStateManager::get_state(bool read_from_db) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  AssertFatal(
      is_initialized,
      "SpgwStateManager init() function should be called to initialize state");
  state_dirty = true;
  AssertFatal(state_cache_p != nullptr, " spgw State cache is NULL");
  if (persist_state_enabled && read_from_db) {
    read_state_from_db();
    read_ue_state_from_db();
  }
  OAILOG_FUNC_RETURN(LOG_S1AP, state_cache_p);
}

}  // namespace lte
}  // namespace magma
