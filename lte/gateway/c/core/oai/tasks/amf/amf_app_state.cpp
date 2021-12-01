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

#include "lte/gateway/c/core/oai/tasks/amf/amf_app_state_manager.h"
#include "lte/gateway/c/core/oai/include/state_manager.h"
#include "lte/gateway/c/core/oai/include/amf_app_state.h"

using magma5g::AmfNasStateManager;

/**
 * When the process starts, initialize the in-memory AMF+NAS state and, if
 * persist state flag is set, load it from the data store.
 * This is only done by the amf_app task.
 */
int amf_nas_state_init(const amf_config_t* amf_config_p) {
  return magma5g::AmfNasStateManager::getInstance().initialize_state(
      amf_config_p);
}

/**
 * Return pointer to the in-memory AMF/NAS state from state manager before
 * processing any message. This is a thread safe call
 * If the read_from_db flag is set to true, the state is loaded from data store
 * before returning the pointer.
 */
magma5g::amf_app_desc_t* get_amf_nas_state(bool read_from_db) {
  return magma5g::AmfNasStateManager::getInstance().get_state(read_from_db);
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
  magma5g::AmfNasStateManager::getInstance().free_state();
}

hash_table_ts_t* get_amf_ue_state() {
  return magma5g::AmfNasStateManager::getInstance().get_ue_state_ht();
}

void put_amf_ue_state(
    magma5g::amf_app_desc_t* amf_app_desc_p, guti_m5_t* guti_p,
    bool force_ue_write) {
  if (AmfNasStateManager::getInstance().is_persist_state_enabled()) {
    if (guti_p != INVALID_IMSI64) {  // need to put appropriate enum
      magma5g::ue_m5gmm_context_t* amf_ue_context = nullptr;
      amf_ue_context =
          amf_ue_context_exists_guti(&amf_app_desc_p->amf_ue_contexts, guti_p);
      // Only write AMF UE state to redis if force flag is set or UE is in EMM
      // Registered state
      if ((amf_ue_context && force_ue_write) ||
          (amf_ue_context &&
           amf_ue_context->mm_state == magma5g::REGISTERED_CONNECTED)) {
        // auto imsi_str =
        // magma5g::AmfNasStateManager::getInstance().get_imsi_str(guti_p);
        // magma5g::AmfNasStateManager::getInstance().write_ue_state_to_db(
        //     amf_ue_context, imsi_str);
      }
    }
  }
}

void delete_amf_ue_state(imsi64_t imsi64) {
  auto imsi_str =
      magma5g::AmfNasStateManager::getInstance().get_imsi_str(imsi64);
  magma5g::AmfNasStateManager::getInstance().clear_ue_state_db(imsi_str);
}
