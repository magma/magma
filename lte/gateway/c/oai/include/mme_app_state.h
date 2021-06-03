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
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include "mme_config.h"
#include "mme_app_desc.h"

/**
 * When the process starts, initialize the in-memory MME+NAS state and, if
 * persist state flag is set, load it from the data store.
 * This is only done by the mme_app task.
 */
int mme_nas_state_init(const mme_config_t* mme_config_p);

/**
 * Return pointer to the in-memory MME/NAS state from state manager before
 * processing any message. This is a thread safe call
 * If the read_from_db flag is set to true, the state is loaded from data store
 * before returning the pointer.
 */
mme_app_desc_t* get_mme_nas_state(bool read_from_db);

/**
 * Write the MME/NAS state to data store after processing any message. This is a
 * thread safe call
 */
void put_mme_nas_state(void);

/**
 * Release the memory allocated for the MME NAS state, this does not clean the
 * state persisted in data store
 */
void clear_mme_nas_state(void);

// Returns UE MME state hashtable, indexed by IMSI
hash_table_ts_t* get_mme_ue_state(void);
// Persists UE MME state for subscriber into db
void put_mme_ue_state(mme_app_desc_t* mme_app_desc_p, imsi64_t imsi64);
// Deletes entry for UE MME state on db
void delete_mme_ue_state(imsi64_t imsi64);

#ifdef __cplusplus
}
#endif
