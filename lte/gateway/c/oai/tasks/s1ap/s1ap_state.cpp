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

#include "s1ap_state.h"

#include <memory.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

extern "C" {
#include "bstrlib.h"

#include "assertions.h"
#include "dynamic_memory_check.h"
}

#include "s1ap_mme.h"
#include "mme_config.h"

s1ap_state_t *s1ap_state_new(void);
void s1ap_state_free(s1ap_state_t *state);

s1ap_state_t *s1ap_state_from_redis(void);
void s1ap_state_to_redis(s1ap_state_t *state);

bool in_use;
s1ap_state_t *state_cache;

int s1ap_state_init(void)
{
  in_use = false;
  state_cache = s1ap_state_from_redis();

  return 0;
}

void s1ap_state_exit(void)
{
  AssertFatal(!in_use, "Exiting without committing s1ap state");

  s1ap_state_free(state_cache);
}

s1ap_state_t *s1ap_state_get(void)
{
  AssertFatal(state_cache != NULL, "s1ap state cache was NULL");
  AssertFatal(!in_use, "Tried to get s1ap_state twice without put'ing it");

  in_use = true;

  return state_cache;
}

void s1ap_state_put(s1ap_state_t *state)
{
  AssertFatal(in_use, "Tried to put s1ap_state while it was not in use");

  state_cache = state;

  s1ap_state_to_redis(state_cache);

  in_use = false;
}

s1ap_state_t *s1ap_state_new(void)
{
  s1ap_state_t *state;
  hash_table_ts_t *ht;
  bstring ht_name;

  state = (s1ap_state_t *) calloc(1, sizeof(*state));
  if (state == NULL) return NULL;

  ht_name = bfromcstr("s1ap_eNB_coll");
  ht = hashtable_ts_init(
    &state->enbs, mme_config.max_enbs, NULL, free_wrapper, ht_name);
  bdestroy(ht_name);

  if (ht == NULL) {
    free(state);
    return NULL;
  }

  ht_name = bfromcstr("s1ap_mme_id2assoc_id_coll");
  ht = hashtable_ts_init(
    &state->mmeid2associd,
    mme_config.max_ues,
    NULL,
    hash_free_int_func,
    ht_name);
  bdestroy(ht_name);

  if (ht == NULL) {
    hashtable_ts_destroy(&state->enbs);
    free(state);
    return NULL;
  }

  return state;
}

void s1ap_state_free(s1ap_state_t *state)
{
  if (hashtable_ts_destroy(&state->enbs) != HASH_TABLE_OK) {
    OAI_FPRINTF_ERR("An error occured while destroying s1 eNB hash table");
  }
  if (hashtable_ts_destroy(&state->mmeid2associd) != HASH_TABLE_OK) {
    OAI_FPRINTF_ERR("An error occured while destroying assoc_id hash table");
  }
  free(state);
}

s1ap_state_t *s1ap_state_from_redis(void)
{
  // todo implement me
  return s1ap_state_new();
}

void s1ap_state_to_redis(s1ap_state_t *state)
{
  // todo implement me
  return;
}
