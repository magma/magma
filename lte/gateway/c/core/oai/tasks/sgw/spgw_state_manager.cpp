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
  state_cache_p = (spgw_state_t*)calloc(1, sizeof(spgw_state_t));

  state_teid_map.map =
      new google::protobuf::Map<uint32_t, magma::lte::oai::S11BearerContext*>();
  state_teid_map.set_name(S11_BEARER_CONTEXT_INFO_HT_NAME);
  state_teid_map.bind_callback(spgw_free_s11_bearer_context_information);

  state_ue_map.map = new google::protobuf::Map<uint64_t, spgw_ue_context_s*>();
  state_ue_map.set_name(SPGW_STATE_UE_MAP);
  state_ue_map.bind_callback(sgw_free_ue_context);

  state_cache_p->sgw_ip_address_S1u_S12_S4_up.s_addr =
      config_->sgw_config.ipv4.S1u_S12_S4_up.s_addr;

  memcpy(&state_cache_p->sgw_ipv6_address_S1u_S12_S4_up,
         &config_->sgw_config.ipv6.S1u_S12_S4_up,
         sizeof(state_cache_p->sgw_ipv6_address_S1u_S12_S4_up));

  state_cache_p->gtpv1u_data.sgw_ip_address_for_S1u_S12_S4_up =
      state_cache_p->sgw_ip_address_S1u_S12_S4_up;

  memcpy(&state_cache_p->gtpv1u_data.sgw_ipv6_address_for_S1u_S12_S4_up,
         &state_cache_p->sgw_ipv6_address_S1u_S12_S4_up,
         sizeof(state_cache_p->gtpv1u_data.sgw_ipv6_address_for_S1u_S12_S4_up));

  state_cache_p->gtpv1u_teid = 0;
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

  free_wrapper((void**)&state_cache_p);
}

status_code_e SpgwStateManager::read_ue_state_from_db() {
  if (!persist_state_enabled) {
    return RETURNok;
  }
  auto keys = redis_client->get_keys("IMSI*" + task_name + "*");
  for (const auto& key : keys) {
    oai::SpgwUeContext ue_proto = oai::SpgwUeContext();
    if (redis_client->read_proto(key, ue_proto) != RETURNok) {
      return RETURNerror;
    }
    OAILOG_DEBUG(log_task, "Reading UE state from db for key %s", key.c_str());
    spgw_ue_context_t* ue_context_p = new spgw_ue_context_t();
    SpgwStateConverter::proto_to_ue(ue_proto, ue_context_p);
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

}  // namespace lte
}  // namespace magma
