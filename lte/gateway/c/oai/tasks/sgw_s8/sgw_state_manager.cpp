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

extern "C" {
#include <dynamic_memory_check.h>
#include "backtrace.h"
}

#include "sgw_state_manager.h"

namespace magma {
namespace lte {

SgwStateManager::SgwStateManager() : config_(nullptr) {}

SgwStateManager::~SgwStateManager() {
  free_state();
}

SgwStateManager& SgwStateManager::getInstance() {
  static SgwStateManager instance;
  return instance;
}

void SgwStateManager::init(bool persist_state, const sgw_config_t* config) {
  log_task              = LOG_SGW_S8;
  task_name             = SGW_TASK_NAME;
  table_key             = SGW_STATE_TABLE_NAME;
  persist_state_enabled = persist_state;
  config_               = config;
  create_state();
  if (read_state_from_db() != RETURNok) {
    OAILOG_ERROR(LOG_SGW_S8, "Failed to read state from redis");
  }
  is_initialized = true;
}

void SgwStateManager::create_state() {
  // Allocating sgw_state_p
  state_cache_p = (sgw_state_t*) calloc(1, sizeof(sgw_state_t));
  display_backtrace();

  OAILOG_INFO(LOG_SGW_S8, "Creating SGW_S8 state ");
  bstring b   = bfromcstr(S11_BEARER_CONTEXT_INFO_HT_NAME);
  state_ue_ht = hashtable_ts_create(
      SGW_STATE_CONTEXT_HT_MAX_SIZE, nullptr,
      (void (*)(void**)) sgw_free_s11_bearer_context_information, b);

  state_cache_p->sgw_ip_address_S1u_S12_S4_up.s_addr =
      config_->ipv4.S1u_S12_S4_up.s_addr;

  state_cache_p->imsi_ue_context_htbl = hashtable_ts_create(
      SGW_STATE_CONTEXT_HT_MAX_SIZE, nullptr,
      (void (*)(void**)) spgw_free_ue_context, nullptr);

  state_cache_p->tunnel_id = 0;

  state_cache_p->gtpv1u_teid = 0;

  bdestroy_wrapper(&b);
}

void SgwStateManager::free_state() {
  AssertFatal(
      is_initialized,
      "SgwStateManager init() function should be called to initialize state.");

  if (state_cache_p == nullptr) {
    return;
  }

  if (hashtable_ts_destroy(state_ue_ht) != HASH_TABLE_OK) {
    OAI_FPRINTF_ERR(
        "An error occurred while destroying SGW s11_bearer_context_information "
        "hashtable");
  }

  hashtable_ts_destroy(state_cache_p->imsi_ue_context_htbl);

  free_wrapper((void**) &state_cache_p);
}

int SgwStateManager::read_ue_state_from_db() {
  /*TODO handle stateless */
  return RETURNok;
}

}  // namespace lte
}  // namespace magma
