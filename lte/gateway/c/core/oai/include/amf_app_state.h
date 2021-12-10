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

#pragma once
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_defs.h"
#ifdef __cplusplus
extern "C" {
#endif

/**
 * When the process starts, initialize the in-memory AMF+NAS state and, if
 * persist state flag is set, load it from the data store.
 * This is only done by the amf_app task.
 */
int amf_nas_state_init(const amf_config_t* amf_config_p);

/**
 * Return pointer to the in-memory AMF/NAS state from state manager before
 * processing any message. This is a thread safe call
 * If the read_from_db flag is set to true, the state is loaded from data store
 * before returning the pointer.
 */
magma5g::amf_app_desc_t* get_amf_nas_state(bool read_from_db);

/**
 * Write the AMF/NAS state to data store after processing any message. This is a
 * thread safe call
 */
void put_amf_nas_state(void);

/**
 * Release the memory allocated for the AMF NAS state, this does not clean the
 * state persisted in data store
 */
void clear_amf_nas_state(void);

// Returns UE AMF state hashtable, indexed by IMSI
hash_table_ts_t* get_amf_ue_state(void);
// Persists UE AMF state for subscriber into db
void put_amf_ue_state(
    magma5g::amf_app_desc_t* amf_app_desc_p, imsi64_t imsi64,
    bool force_ue_write);
// Deletes entry for UE AMF state on db
void delete_amf_ue_state(imsi64_t imsi64);

#ifdef __cplusplus
}
#endif
