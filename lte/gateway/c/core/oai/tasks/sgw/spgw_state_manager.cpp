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

#include "spgw_state_manager.h"

extern "C" {
#include <dynamic_memory_check.h>
#include "common_defs.h"
}

namespace magma {
namespace lte {

SpgwStateManager::SpgwStateManager() : config_(nullptr) {}

SpgwStateManager::~SpgwStateManager() {
  free_state();
}

SpgwStateManager& SpgwStateManager::getInstance() {
  static SpgwStateManager instance;
  return instance;
}

void SpgwStateManager::init(bool persist_state, const spgw_config_t* config) {
  log_task              = LOG_SPGW_APP;
  task_name             = SPGW_TASK_NAME;
  table_key             = SPGW_STATE_TABLE_NAME;
  persist_state_enabled = persist_state;
  config_               = config;
  redis_client          = std::make_unique<RedisClient>(persist_state);
  create_state();
  if (read_state_from_db() != RETURNok) {
    OAILOG_ERROR(LOG_SPGW_APP, "Failed to read state from redis");
  }
  is_initialized = true;
}

void SpgwStateManager::create_state() {
  // Allocating spgw_state_p
  state_cache_p = (spgw_state_t*) calloc(1, sizeof(spgw_state_t));

  bstring b      = bfromcstr(S11_BEARER_CONTEXT_INFO_HT_NAME);
  state_teid_ht_ = hashtable_ts_create(
      SGW_STATE_CONTEXT_HT_MAX_SIZE, nullptr,
      (void (*)(void**)) spgw_free_s11_bearer_context_information, b);

  state_ue_ht = hashtable_ts_create(
      SGW_STATE_CONTEXT_HT_MAX_SIZE, nullptr,
      (void (*)(void**)) sgw_free_ue_context, nullptr);

  state_cache_p->sgw_ip_address_S1u_S12_S4_up.s_addr =
      config_->sgw_config.ipv4.S1u_S12_S4_up.s_addr;

  state_cache_p->gtpv1u_data.sgw_ip_address_for_S1u_S12_S4_up =
      state_cache_p->sgw_ip_address_S1u_S12_S4_up;

  // Creating PGW related state structs
  state_cache_p->deactivated_predefined_pcc_rules = hashtable_ts_create(
      MAX_PREDEFINED_PCC_RULES_HT_SIZE, nullptr, pgw_free_pcc_rule, nullptr);

  state_cache_p->predefined_pcc_rules = hashtable_ts_create(
      MAX_PREDEFINED_PCC_RULES_HT_SIZE, nullptr, pgw_free_pcc_rule, nullptr);

  state_cache_p->gtpv1u_teid = 0;

  bdestroy_wrapper(&b);
}

void SpgwStateManager::free_state() {
  AssertFatal(
      is_initialized,
      "SpgwStateManager init() function should be called to initialize state.");

  if (state_cache_p == nullptr) {
    return;
  }

  if (hashtable_ts_destroy(state_ue_ht) != HASH_TABLE_OK) {
    OAI_FPRINTF_ERR(
        "An error occurred while destroying SGW s11_bearer_context_information "
        "hashtable");
  }

  hashtable_ts_destroy(state_teid_ht_);

  if (state_cache_p->deactivated_predefined_pcc_rules) {
    hashtable_ts_destroy(state_cache_p->deactivated_predefined_pcc_rules);
  }

  if (state_cache_p->predefined_pcc_rules) {
    hashtable_ts_destroy(state_cache_p->predefined_pcc_rules);
  }
  free_wrapper((void**) &state_cache_p);
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
    spgw_ue_context_t* ue_context_p =
        (spgw_ue_context_t*) calloc(1, sizeof(spgw_ue_context_t));
    SpgwStateConverter::proto_to_ue(ue_proto, ue_context_p);
  }
  return RETURNok;
}

hash_table_ts_t* SpgwStateManager::get_state_teid_ht() {
  AssertFatal(
      is_initialized,
      "StateManager init() function should be called to initialize state");

  return state_teid_ht_;
}

}  // namespace lte
}  // namespace magma
