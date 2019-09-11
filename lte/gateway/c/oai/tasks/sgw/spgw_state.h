/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
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

#include <pthread.h>

#include "hashtable.h"

#include "gtpv1u_sgw_defs.h"
#include "sgw_context_manager.h"

typedef struct sgw_state_s {
  struct in_addr sgw_ip_address_S1u_S12_S4_up;

  // Maps teid (as uint32 key) to mme_sgw_tunnel
  hash_table_ts_t *s11teid2mme;
  // Maps teid (as uint32 key) to s11_eps_bearer_context_information
  hash_table_ts_t *s11_bearer_context_information;

  gtpv1u_data_t gtpv1u_data;

} sgw_state_t;

typedef struct spgw_state_s {
  sgw_state_t sgw_state;
} spgw_state_t;

// Initializes SGW state struct when task process starts.
int spgw_state_init(bool use_stateless);
// Function that frees spgw_state.
void spgw_state_exit(void);
// Function that returns a pointer to spgw_state.
spgw_state_t *get_spgw_state(void);

#ifdef __cplusplus
}
#endif
