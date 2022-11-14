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

#pragma once

#include <pthread.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/lib/hashtable/hashtable.h"
#include "lte/gateway/c/core/oai/include/sgw_config.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/include/spgw_types.hpp"

// Initializes SGW state struct when task process starts.
int sgw_state_init(bool persist_state, const sgw_config_t* config);

// Function that frees sgw_state.
void sgw_state_exit(void);

// Function that returns a pointer to sgw_state.
sgw_state_t* get_sgw_state(bool read_from_db);

// Function that writes the sgw_state struct into db.
void put_sgw_state(void);

/**
 * Returns pointer to SGW UE state
 * @return hashtable_ts_t struct with SGW UE context structs as data
 */
hash_table_ts_t* get_sgw_ue_state(void);

/**
 * Populates SGW UE hashtable from db
 * @return response code
 */
int read_sgw_ue_state_db(void);

/**
 * Saves an UE context state to db
 * @param pointer to SGW's state
 * @param imsi64
 */
void put_sgw_ue_state(sgw_state_t* sgw_state, imsi64_t imsi64);

/**
 * Removes entry for SGW UE state in db
 * @param imsi64
 */
void delete_sgw_ue_state(imsi64_t imsi64);

/**
 * Callback function for s11_bearer_context_information hashtable freefunc
 * @param context_p sgw eps bearer context entry on map
 */
void sgw_free_s11_bearer_context_information(void** ptr);

/**
 * Callback function for imsi_ue_context hashtable's freefunc
 * @param spgw_ue_context_t
 */
void sgw_free_ue_context(spgw_ue_context_t** ue_context_p);
