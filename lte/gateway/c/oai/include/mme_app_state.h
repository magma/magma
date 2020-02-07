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

#include "mme_config.h"
#include "mme_app_desc.h"

/**
 * When the process starts, initialize the in-memory MME+NAS state and, if
 * persist state flag is set, load it from the data store.
 * This is only done by the mme_app task.
 */
int mme_nas_state_init(const mme_config_t* mme_config_p);

/**
 * Lock the MME NAS state and return pointer to the in-memory MME/NAS state from
 * state manager before processing any message. This is a thread safe call. If
 * the read_from_db flag is set to true, the state is loaded from data store
 * before returning the pointer. The lock acquired here is released by
 * put_mme_nas_state()
*/
mme_app_desc_t* get_locked_mme_nas_state(bool read_from_db);

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
void put_mme_nas_state(mme_app_desc_t** task_state_ptr);

/**
 * Release the memory allocated for the MME NAS state, this does not clean the
 * state persisted in data store
*/
void clear_mme_nas_state(void);

/**
 * Converts mme_imsi_map to protobuf and saves it into data store
 */
void put_mme_imsi_map(void);

/**
 * @return mme_imsi_map_t pointer
 */
mme_imsi_map_t* get_mme_imsi_map(void);

#ifdef __cplusplus
}
#endif
