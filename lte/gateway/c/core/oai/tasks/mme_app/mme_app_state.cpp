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

#include "lte/gateway/c/core/oai/include/mme_app_state.hpp"

#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_state_manager.hpp"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_ip_imsi.hpp"

using magma::lte::MmeNasStateManager;

/**
 * When the process starts, initialize the in-memory MME+NAS state and, if
 * persist state flag is set, load it from the data store.
 * This is only done by the mme_app task.
 */
int mme_nas_state_init(const mme_config_t* mme_config_p) {
  return MmeNasStateManager::getInstance().initialize_state(mme_config_p);
}

/**
 * Return pointer to the in-memory MME/NAS state from state manager before
 * processing any message. This is a thread safe call
 * If the read_from_db flag is set to true, the state is loaded from data store
 * before returning the pointer.
 */
mme_app_desc_t* get_mme_nas_state(bool read_from_db) {
  return MmeNasStateManager::getInstance().get_state(read_from_db);
}

/**
 * Write the MME/NAS state to data store after processing any message. This is
 * a thread safe call
 */
void put_mme_nas_state() {
  MmeNasStateManager::getInstance().write_state_to_db();
}

/**
 * Release the memory allocated for the MME NAS state, this does not clean the
 * state persisted in data store
 */
void clear_mme_nas_state() { MmeNasStateManager::getInstance().free_state(); }

proto_map_uint32_ue_context_t* get_mme_ue_state() {
  return MmeNasStateManager::getInstance().get_ue_state_map();
}

void put_mme_ue_state(mme_app_desc_t* mme_app_desc_p, imsi64_t imsi64,
                      bool force_ue_write) {
  if (MmeNasStateManager::getInstance().is_persist_state_enabled()) {
    if (imsi64 != INVALID_IMSI64) {
      ue_mm_context_t* ue_context = nullptr;
      ue_context =
          mme_ue_context_exists_imsi(&mme_app_desc_p->mme_ue_contexts, imsi64);
      // Only write MME UE state to redis if force flag is set or UE is in EMM
      // Registered state
      if ((ue_context && force_ue_write) ||
          (ue_context && ue_context->mm_state == UE_REGISTERED)) {
        auto imsi_str = MmeNasStateManager::getInstance().get_imsi_str(imsi64);
        MmeNasStateManager::getInstance().write_ue_state_to_db(ue_context,
                                                               imsi_str);
      }
    }
  }
}

void delete_mme_ue_state(imsi64_t imsi64) {
  auto imsi_str = MmeNasStateManager::getInstance().get_imsi_str(imsi64);
  MmeNasStateManager::getInstance().clear_ue_state_db(imsi_str);
}
