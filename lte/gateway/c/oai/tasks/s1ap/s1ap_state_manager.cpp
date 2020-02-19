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
 *------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#include "s1ap_state_manager.h"

namespace {
constexpr char S1AP_ENB_COLL[] = "s1ap_eNB_coll";
constexpr char S1AP_MME_ID2ASSOC_ID_COLL[] = "s1ap_mme_id2assoc_id_coll";
constexpr char S1AP_IMSI_MAP_TABLE_NAME[] = "s1ap_imsi_map";
} // namespace

namespace magma {
namespace lte {

S1apStateManager::S1apStateManager(): max_enbs_(0), max_ues_(0) {}

S1apStateManager::~S1apStateManager()
{
  free_state();
}

S1apStateManager& S1apStateManager::getInstance()
{
  static S1apStateManager instance;
  return instance;
}

void S1apStateManager::init(
  uint32_t max_ues,
  uint32_t max_enbs,
  bool use_stateless)
{
  log_task = LOG_S1AP;
  table_key = S1AP_STATE_TABLE;
  persist_state_enabled = use_stateless;
  max_ues_ = max_ues;
  max_enbs_ = max_enbs;
  create_state();
  if (read_state_from_db() != RETURNok) {
    OAILOG_ERROR(LOG_S1AP, "Failed to read state from redis");
  }
  is_initialized = true;
}

void S1apStateManager::create_state()
{
  bstring ht_name;

  state_cache_p = (s1ap_state_t*) calloc(1, sizeof(s1ap_state_t));

  ht_name = bfromcstr(S1AP_ENB_COLL);
  hashtable_ts_init(
    &state_cache_p->enbs, max_enbs_, nullptr, free_wrapper, ht_name);
  bdestroy(ht_name);

  ht_name = bfromcstr(S1AP_MME_ID2ASSOC_ID_COLL);
  hashtable_ts_init(
    &state_cache_p->mmeid2associd,
    max_ues_,
    nullptr,
    hash_free_int_func,
    ht_name);
  bdestroy(ht_name);

  state_cache_p->num_enbs = 0;

  create_s1ap_imsi_map();
}

void S1apStateManager::free_state()
{
  AssertFatal(
    is_initialized,
    "S1apStateManager init() function should be called to initialize state.");

  if (state_cache_p == nullptr) {
    return;
  }

  int i;
  hashtable_rc_t ht_rc;
  hashtable_key_array_t* keys;
  sctp_assoc_id_t assoc_id;
  enb_description_t* enb;

  keys = hashtable_ts_get_keys(&state_cache_p->enbs);
  if (!keys) {
    OAILOG_DEBUG(LOG_S1AP, "No keys in the enb hashtable");
  } else {
    for (i = 0; i < keys->num_keys; i++) {
      assoc_id = (sctp_assoc_id_t) keys->keys[i];
      ht_rc = hashtable_ts_get(
        &state_cache_p->enbs, (hash_key_t) assoc_id, (void**) &enb);
      AssertFatal(ht_rc == HASH_TABLE_OK, "enbueid not in assoc_id");

      if (hashtable_ts_destroy(&enb->ue_coll) != HASH_TABLE_OK) {
        OAI_FPRINTF_ERR("An error occured while destroying UE coll hash table");
      }
    }
    FREE_HASHTABLE_KEY_ARRAY(keys);
  }

  if (hashtable_ts_destroy(&state_cache_p->enbs) != HASH_TABLE_OK) {
    OAI_FPRINTF_ERR("An error occured while destroying s1 eNB hash table");
  }
  if (hashtable_ts_destroy(&state_cache_p->mmeid2associd) != HASH_TABLE_OK) {
    OAI_FPRINTF_ERR("An error occured while destroying assoc_id hash table");
  }
  free_wrapper((void**) &state_cache_p);

  clear_s1ap_imsi_map();
}

void S1apStateManager::create_s1ap_imsi_map()
{
  s1ap_imsi_map_ = (s1ap_imsi_map_t*) calloc(1, sizeof(s1ap_imsi_map_t));

  s1ap_imsi_map_->mme_ue_id_imsi_htbl =
    hashtable_uint64_ts_create(max_ues_, nullptr, nullptr);

  gateway::s1ap::S1apImsiMap imsi_proto = gateway::s1ap::S1apImsiMap();
  redis_client->read_proto(S1AP_IMSI_MAP_TABLE_NAME, imsi_proto);

  S1apStateConverter::proto_to_s1ap_imsi_map(imsi_proto, s1ap_imsi_map_);
}

void S1apStateManager::clear_s1ap_imsi_map() {
  if(!s1ap_imsi_map_) {
    return;
  }
  hashtable_uint64_ts_destroy(s1ap_imsi_map_->mme_ue_id_imsi_htbl);

  free_wrapper((void **) &s1ap_imsi_map_);
}

s1ap_imsi_map_t* S1apStateManager::get_s1ap_imsi_map()
{
  return s1ap_imsi_map_;
}

void S1apStateManager::put_s1ap_imsi_map() {
  gateway::s1ap::S1apImsiMap imsi_proto = gateway::s1ap::S1apImsiMap();
  S1apStateConverter::s1ap_imsi_map_to_proto(s1ap_imsi_map_, &imsi_proto);
  redis_client->write_proto(S1AP_IMSI_MAP_TABLE_NAME, imsi_proto);
}

} // namespace lte
} // namespace magma
