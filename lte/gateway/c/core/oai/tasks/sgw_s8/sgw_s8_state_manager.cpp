/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

#include "lte/gateway/c/core/oai/tasks/sgw_s8/sgw_s8_state_manager.hpp"

extern "C" {
#include "lte/gateway/c/core/common/backtrace.h"
}

#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/include/sgw_context_manager.hpp"

namespace magma {
namespace lte {

SgwStateManager::SgwStateManager() : config_(nullptr) {}

SgwStateManager::~SgwStateManager() { free_state(); }

SgwStateManager& SgwStateManager::getInstance() {
  static SgwStateManager instance;
  return instance;
}

void SgwStateManager::init(bool persist_state, const sgw_config_t* config) {
  log_task = LOG_SGW_S8;
  task_name = SGW_TASK_NAME;
  table_key = SGW_STATE_TABLE_NAME;
  persist_state_enabled = persist_state;
  config_ = config;
  create_state();
  if (read_state_from_db() != RETURNok) {
    OAILOG_ERROR(LOG_SGW_S8, "Failed to read state from redis");
  }
  is_initialized = true;
}

void SgwStateManager::create_state() {
  state_cache_p = new sgw_state_t();
  if (!state_cache_p) {
    OAILOG_CRITICAL(LOG_SGW_S8,
                    "Failed to allocate memory for sgw_state_t structure \n ");
    return;
  }

  OAILOG_INFO(LOG_SGW_S8, "Creating SGW_S8 state ");

  // sgw_free_s11_bearer_context_information is called when destroy_map or
  // remove_map is invoked, so as to remove any contexts allocated within
  // sgw_bearer context
  s8_state_teid_map.map =
      new google::protobuf::Map<unsigned int,
                                struct sgw_eps_bearer_context_information_s*>();
  s8_state_teid_map.set_name(S11_BEARER_CONTEXT_INFO_MAP_NAME);
  s8_state_teid_map.bind_callback(sgw_free_s11_bearer_context_information);

  state_cache_p->sgw_ip_address_S1u_S12_S4_up.s_addr =
      config_->ipv4.S1u_S12_S4_up.s_addr;

  memcpy(&state_cache_p->sgw_ipv6_address_S1u_S12_S4_up,
         &config_->ipv4.S1u_S12_S4_up.s_addr,
         sizeof(state_cache_p->sgw_ipv6_address_S1u_S12_S4_up));

  state_cache_p->sgw_ip_address_S5S8_up.s_addr = config_->ipv4.S5_S8_up.s_addr;

  state_cache_p->imsi_ue_context_map.map =
      new google::protobuf::Map<uint64_t, struct spgw_ue_context_s*>();
  state_cache_p->imsi_ue_context_map.set_name(SGW_S8_STATE_UE_MAP_NAME);
  state_cache_p->imsi_ue_context_map.bind_callback(sgw_free_ue_context);

  state_cache_p->temporary_create_session_procedure_id_map.map =
      new google::protobuf::Map<uint32_t,
                                struct sgw_eps_bearer_context_information_s*>();
  state_cache_p->temporary_create_session_procedure_id_map.set_name(
      SGW_S8_CSR_PROC_ID_MAP);
  state_cache_p->temporary_create_session_procedure_id_map.bind_callback(
      sgw_free_s11_bearer_context_information);

  state_cache_p->s1u_teid = INITIAL_SGW_S8_S1U_TEID;
  state_cache_p->s5s8u_teid = 0;
}

void SgwStateManager::free_state() {
  AssertFatal(
      is_initialized,
      "SgwStateManager init() function should be called to initialize state.");

  if (state_cache_p == nullptr) {
    return;
  }

  if (s8_state_teid_map.map) {
    if (s8_state_teid_map.destroy_map() != PROTO_MAP_OK) {
      OAILOG_ERROR(LOG_SGW_S8,
                   "An error occurred while destroying "
                   "temporary_create_session_procedure_id_map ");
    }
  }

  if (state_cache_p->imsi_ue_context_map.map) {
    if (state_cache_p->imsi_ue_context_map.destroy_map() != PROTO_MAP_OK) {
      OAILOG_ERROR(LOG_SGW_S8,
                   "An error occurred while destroying "
                   "imsi_ue_context_map ");
    }
  }

  if (state_cache_p->temporary_create_session_procedure_id_map.map) {
    if (state_cache_p->temporary_create_session_procedure_id_map
            .destroy_map() != PROTO_MAP_OK) {
      OAILOG_ERROR(LOG_SGW_S8,
                   "An error occurred while destroying "
                   "temporary_create_session_procedure_id_map ");
    }
  }
  free_cpp_wrapper(reinterpret_cast<void**>(&state_cache_p));
}

status_code_e SgwStateManager::read_ue_state_from_db() {
  /* TODO handle stateless for SGW_S8 task */
  return RETURNok;
}

sgw_state_t* SgwStateManager::get_state(bool read_from_db) {
  return state_cache_p;
}

map_uint32_sgw_eps_bearer_context_t* SgwStateManager::get_s8_state_teid_map() {
  AssertFatal(
      is_initialized,
      "StateManager init() function should be called to initialize state");

  return &s8_state_teid_map;
}

}  // namespace lte
}  // namespace magma
