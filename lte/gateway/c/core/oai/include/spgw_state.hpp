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

#include <pthread.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/lib/hashtable/hashtable.h"
#include "lte/gateway/c/core/oai/include/gtpv1u_types.h"
#include "lte/gateway/c/core/oai/include/spgw_config.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/include/spgw_types.hpp"

// Initializes SGW state struct when task process starts.
int spgw_state_init(bool persist_state, const spgw_config_t* spgw_config_p);
// Function that frees spgw_state.
void spgw_state_exit(void);
// Function that returns a pointer to spgw_state.
spgw_state_t* get_spgw_state(bool read_from_db);
// Function that writes the spgw_state struct into db.
void put_spgw_state(void);

// retunrs pointer to proto map, map_uint64_spgw_ue_context_t
map_uint64_spgw_ue_context_t* get_spgw_ue_state(void);

state_teid_map_t* get_spgw_teid_state(void);

/**
 * Populates SPGW UE hashtable from db
 * @return response code
 */
int read_spgw_ue_state_db(void);

/**
 * Saves an UE context state to db
 * @param s11_bearer_context_info SPGW ue context pointer
 * @param imsi64
 */
void put_spgw_ue_state(imsi64_t imsi64);

/**
 * Removes entry for SPGW UE state in db
 * @param imsi64
 */
void delete_spgw_ue_state(imsi64_t imsi64);

/**
 * Callback function for s11_bearer_context_information hashtable freefunc
 * @param context_p spgw eps bearer context entry on map
 */
void spgw_free_s11_bearer_context_information(void**);
/**
 * Frees pdn connection and its contained objects
 * @param pdn_connection_p
 */
void sgw_free_pdn_connection(sgw_pdn_connection_t* pdn_connection_p);
/**
 * Frees sgw_eps_bearer_ctxt entry
 * @param sgw_eps_bearer_ctxt
 */
void sgw_free_eps_bearer_context(sgw_eps_bearer_ctxt_t** sgw_eps_bearer_ctxt);
/**
 * Callback function for imsi_ue_context hashtable's freefunc
 * @param spgw_ue_context_t
 */
void sgw_free_ue_context(void** ptr);
