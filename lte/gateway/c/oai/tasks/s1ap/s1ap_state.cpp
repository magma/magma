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

#include <cstdlib>
#include <cstring>

#include <memory.h>

extern "C" {
#include "bstrlib.h"

#include "assertions.h"
#include "common_defs.h"
#include "dynamic_memory_check.h"
}

#include "s1ap_state_manager.h"

using magma::lte::S1apStateManager;

int s1ap_state_init(uint32_t max_ues, uint32_t max_enbs, bool use_stateless)
{
  S1apStateManager::getInstance().init(max_ues, max_enbs, use_stateless);
  return RETURNok;
}

s1ap_state_t* get_s1ap_state(bool read_from_db)
{
  return S1apStateManager::getInstance().get_state(read_from_db);
}

void s1ap_state_exit()
{
  S1apStateManager::getInstance().free_state();
}

void put_s1ap_state()
{
  S1apStateManager::getInstance().write_state_to_db();
}

enb_description_t* s1ap_state_get_enb(
  s1ap_state_t* state,
  sctp_assoc_id_t assoc_id)
{
  enb_description_t* enb = nullptr;

  hashtable_ts_get(&state->enbs, (const hash_key_t) assoc_id, (void**) &enb);

  return enb;
}

ue_description_t* s1ap_state_get_ue_enbid(
  sctp_assoc_id_t sctp_assoc_id,
  enb_ue_s1ap_id_t enb_ue_s1ap_id)
{
  ue_description_t* ue = nullptr;

  hash_table_ts_t* state_ue_ht = get_s1ap_ue_state();
  uint64_t comp_s1ap_id = (uint64_t) enb_ue_s1ap_id << 32 | sctp_assoc_id;
  hashtable_ts_get(
    state_ue_ht, (const hash_key_t) comp_s1ap_id, (void**) &ue);

  return ue;
}

ue_description_t* s1ap_state_get_ue_mmeid(mme_ue_s1ap_id_t mme_ue_s1ap_id)
{
  ue_description_t* ue = nullptr;

  hash_table_ts_t* state_ue_ht = get_s1ap_ue_state();
  hashtable_ts_apply_callback_on_elements(
      (hash_table_ts_t* const) state_ue_ht,
      s1ap_ue_compare_by_mme_ue_id_cb,
      &mme_ue_s1ap_id,
      (void**) &ue);

  return ue;
}

ue_description_t* s1ap_state_get_ue_imsi(imsi64_t imsi64)
{
  ue_description_t* ue = nullptr;

  hash_table_ts_t* state_ue_ht = get_s1ap_ue_state();
  hashtable_ts_apply_callback_on_elements(
      (hash_table_ts_t* const) state_ue_ht,
      s1ap_ue_compare_by_imsi,
      &imsi64,
      (void**) &ue);

  return ue;
}

uint64_t s1ap_get_comp_s1ap_id(
    sctp_assoc_id_t sctp_assoc_id,
    enb_ue_s1ap_id_t enb_ue_s1ap_id)
{
  return (uint64_t) enb_ue_s1ap_id << 32 | sctp_assoc_id;
}

void put_s1ap_imsi_map() {
  S1apStateManager::getInstance().put_s1ap_imsi_map();
}

s1ap_imsi_map_t* get_s1ap_imsi_map() {
  return S1apStateManager::getInstance().get_s1ap_imsi_map();
}

bool s1ap_ue_compare_by_mme_ue_id_cb(
  __attribute__((unused)) const hash_key_t keyP,
  void* const elementP,
  void* parameterP,
  void** resultP)
{
  mme_ue_s1ap_id_t* mme_ue_s1ap_id_p = (mme_ue_s1ap_id_t*) parameterP;
  ue_description_t* ue_ref = (ue_description_t*) elementP;
  if (*mme_ue_s1ap_id_p == ue_ref->mme_ue_s1ap_id) {
    *resultP = elementP;
    OAILOG_TRACE(
      LOG_S1AP,
      "Found ue_ref %p mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT "\n",
      ue_ref,
      ue_ref->mme_ue_s1ap_id);
    return true;
  }
  return false;
}

bool s1ap_ue_compare_by_imsi(
    __attribute__((unused)) const hash_key_t keyP,
    void* const elementP,
    void* parameterP,
    void** resultP) {
  imsi64_t imsi64 = INVALID_IMSI64;
  imsi64_t* target_imsi64 = (imsi64_t*) parameterP;
  ue_description_t* ue_ref = (ue_description_t*) elementP;

  s1ap_imsi_map_t* imsi_map = get_s1ap_imsi_map();
  hashtable_uint64_ts_get(
      imsi_map->mme_ue_id_imsi_htbl, (const hash_key_t) ue_ref->mme_ue_s1ap_id,
      &imsi64);

  if (*target_imsi64 != INVALID_IMSI64 && *target_imsi64 == imsi64) {
    *resultP = elementP;
    OAILOG_DEBUG_UE(
        LOG_S1AP, imsi64, "Found ue_ref\n");
    return true;
  }
  return false;
}

hash_table_ts_t* get_s1ap_ue_state(void) {
  return S1apStateManager::getInstance().get_ue_state_ht();
}

void put_s1ap_ue_state(imsi64_t imsi64) {
  if(S1apStateManager::getInstance().is_persist_state_enabled()) {
    ue_description_t* ue_ctxt = s1ap_state_get_ue_imsi(imsi64);
    if (ue_ctxt) {
      auto imsi_str = S1apStateManager::getInstance().get_imsi_str(imsi64);
      S1apStateManager::getInstance().write_ue_state_to_db(ue_ctxt, imsi_str);
    }
  }
}

void delete_s1ap_ue_state(imsi64_t imsi64) {
  auto imsi_str = S1apStateManager::getInstance().get_imsi_str(imsi64);
  S1apStateManager::getInstance().clear_ue_state_db(imsi_str);
}
