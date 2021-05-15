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

int s1ap_state_init(uint32_t max_ues, uint32_t max_enbs, bool use_stateless) {
  S1apStateManager::getInstance().init(max_ues, max_enbs, use_stateless);
  // remove UEs with unknown IMSI from eNB state
  remove_ues_without_imsi_from_ue_id_coll();
  return RETURNok;
}

s1ap_state_t* get_s1ap_state(bool read_from_db) {
  return S1apStateManager::getInstance().get_state(read_from_db);
}

void s1ap_state_exit() {
  S1apStateManager::getInstance().free_state();
}

void put_s1ap_state() {
  S1apStateManager::getInstance().write_state_to_db();
}

enb_description_t* s1ap_state_get_enb(
    s1ap_state_t* state, sctp_assoc_id_t assoc_id) {
  enb_description_t* enb = nullptr;

  hashtable_ts_get(&state->enbs, (const hash_key_t) assoc_id, (void**) &enb);

  return enb;
}

ue_description_t* s1ap_state_get_ue_enbid(
    sctp_assoc_id_t sctp_assoc_id, enb_ue_s1ap_id_t enb_ue_s1ap_id) {
  ue_description_t* ue = nullptr;

  hash_table_ts_t* state_ue_ht = get_s1ap_ue_state();
  uint64_t comp_s1ap_id =
      S1AP_GENERATE_COMP_S1AP_ID(sctp_assoc_id, enb_ue_s1ap_id);
  hashtable_ts_get(state_ue_ht, (const hash_key_t) comp_s1ap_id, (void**) &ue);

  return ue;
}

ue_description_t* s1ap_state_get_ue_mmeid(mme_ue_s1ap_id_t mme_ue_s1ap_id) {
  ue_description_t* ue = nullptr;

  hash_table_ts_t* state_ue_ht = get_s1ap_ue_state();
  hashtable_ts_apply_callback_on_elements(
      (hash_table_ts_t* const) state_ue_ht, s1ap_ue_compare_by_mme_ue_id_cb,
      &mme_ue_s1ap_id, (void**) &ue);

  return ue;
}

ue_description_t* s1ap_state_get_ue_imsi(imsi64_t imsi64) {
  ue_description_t* ue = nullptr;

  hash_table_ts_t* state_ue_ht = get_s1ap_ue_state();
  hashtable_ts_apply_callback_on_elements(
      (hash_table_ts_t* const) state_ue_ht, s1ap_ue_compare_by_imsi, &imsi64,
      (void**) &ue);

  return ue;
}

void put_s1ap_imsi_map() {
  S1apStateManager::getInstance().write_s1ap_imsi_map_to_db();
}

s1ap_imsi_map_t* get_s1ap_imsi_map() {
  return S1apStateManager::getInstance().get_s1ap_imsi_map();
}

bool s1ap_ue_compare_by_mme_ue_id_cb(
    __attribute__((unused)) const hash_key_t keyP, void* const elementP,
    void* parameterP, void** resultP) {
  mme_ue_s1ap_id_t* mme_ue_s1ap_id_p = (mme_ue_s1ap_id_t*) parameterP;
  ue_description_t* ue_ref           = (ue_description_t*) elementP;
  if (*mme_ue_s1ap_id_p == ue_ref->mme_ue_s1ap_id) {
    *resultP = elementP;
    OAILOG_TRACE(
        LOG_S1AP, "Found ue_ref %p mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT "\n",
        ue_ref, ue_ref->mme_ue_s1ap_id);
    return true;
  }
  return false;
}

bool s1ap_ue_compare_by_imsi(
    __attribute__((unused)) const hash_key_t keyP, void* const elementP,
    void* parameterP, void** resultP) {
  imsi64_t imsi64          = INVALID_IMSI64;
  imsi64_t* target_imsi64  = (imsi64_t*) parameterP;
  ue_description_t* ue_ref = (ue_description_t*) elementP;

  s1ap_imsi_map_t* imsi_map = get_s1ap_imsi_map();
  hashtable_uint64_ts_get(
      imsi_map->mme_ue_id_imsi_htbl, (const hash_key_t) ue_ref->mme_ue_s1ap_id,
      &imsi64);

  if (*target_imsi64 != INVALID_IMSI64 && *target_imsi64 == imsi64) {
    *resultP = elementP;
    OAILOG_DEBUG_UE(LOG_S1AP, imsi64, "Found ue_ref\n");
    return true;
  }
  return false;
}

hash_table_ts_t* get_s1ap_ue_state(void) {
  return S1apStateManager::getInstance().get_ue_state_ht();
}

void put_s1ap_ue_state(imsi64_t imsi64) {
  if (S1apStateManager::getInstance().is_persist_state_enabled()) {
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

bool get_mme_ue_ids_no_imsi(
    const hash_key_t keyP, uint64_t const dataP, void* argP, void** resultP) {
  hash_key_t** mme_id_list   = (hash_key_t**) resultP;
  uint32_t* num_ues_checked  = (uint32_t*) argP;
  ue_description_t* ue_ref_p = NULL;

  // Check if a UE reference exists for this comp_s1ap_id
  hash_table_ts_t* s1ap_ue_state = get_s1ap_ue_state();
  hashtable_ts_get(s1ap_ue_state, (const hash_key_t) dataP, (void**) &ue_ref_p);
  if (!ue_ref_p) {
    (*mme_id_list)[*num_ues_checked] = keyP;
    ++(*num_ues_checked);
    OAILOG_DEBUG(
        LOG_S1AP,
        "Adding mme_ue_s1ap_id %lu to eNB clean up list with num_ues_checked "
        "%u",
        keyP, *num_ues_checked);
  }
  return false;  // always return false to make sure it runs on all elements
}

void remove_ues_without_imsi_from_ue_id_coll() {
  s1ap_state_t* s1ap_state_p     = get_s1ap_state(false);
  hashtable_key_array_t* ht_keys = hashtable_ts_get_keys(&s1ap_state_p->enbs);
  if (ht_keys == nullptr) {
    return;
  }

  hashtable_rc_t ht_rc;
  hash_key_t* mme_ue_id_no_imsi_list;
  s1ap_imsi_map_t* s1ap_imsi_map = get_s1ap_imsi_map();
  uint32_t num_ues_checked;

  // get each eNB in s1ap_state
  for (int i = 0; i < ht_keys->num_keys; i++) {
    enb_description_t* enb_association_p = nullptr;
    ht_rc                                = hashtable_ts_get(
        &s1ap_state_p->enbs, (hash_key_t) ht_keys->keys[i],
        (void**) &enb_association_p);
    if (ht_rc != HASH_TABLE_OK) {
      continue;
    }

    if (enb_association_p->ue_id_coll.num_elements == 0) {
      continue;
    }

    // for each ue comp_s1ap_id in eNB->ue_id_coll, check if it has an S1ap
    // ue_context, if not delete it
    num_ues_checked        = 0;
    mme_ue_id_no_imsi_list = (hash_key_t*) calloc(
        enb_association_p->ue_id_coll.num_elements, sizeof(hash_key_t));
    hashtable_uint64_ts_apply_callback_on_elements(
        &enb_association_p->ue_id_coll, get_mme_ue_ids_no_imsi,
        &num_ues_checked, (void**) &mme_ue_id_no_imsi_list);

    // remove all the mme_ue_s1ap_ids
    for (uint32_t i = 0; i < num_ues_checked; i++) {
      hashtable_uint64_ts_remove(
          &enb_association_p->ue_id_coll, mme_ue_id_no_imsi_list[i]);
      hashtable_uint64_ts_remove(
          s1ap_imsi_map->mme_ue_id_imsi_htbl, mme_ue_id_no_imsi_list[i]);
      enb_association_p->nb_ue_associated--;

      OAILOG_DEBUG(
          LOG_S1AP, "Num UEs associated %u num ue_id_coll %zu",
          enb_association_p->nb_ue_associated,
          enb_association_p->ue_id_coll.num_elements);
    }

    // free the list
    free(mme_ue_id_no_imsi_list);
  }

  FREE_HASHTABLE_KEY_ARRAY(ht_keys);
}
