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

#include "lte/gateway/c/core/oai/include/s1ap_state.hpp"

#include <cstdlib>
#include <cstring>
#include <vector>

#include <memory.h>

extern "C" {
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/common/common_defs.h"
}

#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_state_manager.hpp"

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

void s1ap_state_exit() { S1apStateManager::getInstance().free_state(); }

void put_s1ap_state() { S1apStateManager::getInstance().write_state_to_db(); }

enb_description_t* s1ap_state_get_enb(s1ap_state_t* state,
                                      sctp_assoc_id_t assoc_id) {
  enb_description_t* enb = nullptr;

  state->enbs.get(assoc_id, &enb);

  return enb;
}

ue_description_t* s1ap_state_get_ue_enbid(sctp_assoc_id_t sctp_assoc_id,
                                          enb_ue_s1ap_id_t enb_ue_s1ap_id) {
  ue_description_t* ue = nullptr;

  hash_table_ts_t* state_ue_ht = get_s1ap_ue_state();
  uint64_t comp_s1ap_id =
      S1AP_GENERATE_COMP_S1AP_ID(sctp_assoc_id, enb_ue_s1ap_id);
  hashtable_ts_get(state_ue_ht, (const hash_key_t)comp_s1ap_id, (void**)&ue);

  return ue;
}

ue_description_t* s1ap_state_get_ue_mmeid(mme_ue_s1ap_id_t mme_ue_s1ap_id) {
  ue_description_t* ue = nullptr;

  hash_table_ts_t* state_ue_ht = get_s1ap_ue_state();
  hashtable_ts_apply_callback_on_elements((hash_table_ts_t* const)state_ue_ht,
                                          s1ap_ue_compare_by_mme_ue_id_cb,
                                          &mme_ue_s1ap_id, (void**)&ue);

  return ue;
}

ue_description_t* s1ap_state_get_ue_imsi(imsi64_t imsi64) {
  ue_description_t* ue = nullptr;

  hash_table_ts_t* state_ue_ht = get_s1ap_ue_state();
  hashtable_ts_apply_callback_on_elements((hash_table_ts_t* const)state_ue_ht,
                                          s1ap_ue_compare_by_imsi, &imsi64,
                                          (void**)&ue);

  return ue;
}

void put_s1ap_imsi_map() {
  S1apStateManager::getInstance().write_s1ap_imsi_map_to_db();
}

s1ap_imsi_map_t* get_s1ap_imsi_map() {
  return S1apStateManager::getInstance().get_s1ap_imsi_map();
}

bool s1ap_ue_compare_by_mme_ue_id_cb(__attribute__((unused))
                                     const hash_key_t keyP,
                                     void* const elementP, void* parameterP,
                                     void** resultP) {
  mme_ue_s1ap_id_t* mme_ue_s1ap_id_p = (mme_ue_s1ap_id_t*)parameterP;
  ue_description_t* ue_ref = (ue_description_t*)elementP;
  if (*mme_ue_s1ap_id_p == ue_ref->mme_ue_s1ap_id) {
    *resultP = elementP;
    OAILOG_TRACE(LOG_S1AP,
                 "Found ue_ref %p mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT "\n",
                 ue_ref, ue_ref->mme_ue_s1ap_id);
    return true;
  }
  return false;
}

bool s1ap_ue_compare_by_imsi(__attribute__((unused)) const hash_key_t keyP,
                             void* const elementP, void* parameterP,
                             void** resultP) {
  imsi64_t imsi64 = INVALID_IMSI64;
  imsi64_t* target_imsi64 = (imsi64_t*)parameterP;
  ue_description_t* ue_ref = (ue_description_t*)elementP;

  s1ap_imsi_map_t* imsi_map = get_s1ap_imsi_map();
  imsi_map->mme_ueid2imsi_map.get(ue_ref->mme_ue_s1ap_id, &imsi64);

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

void remove_ues_without_imsi_from_ue_id_coll() {
  s1ap_state_t* s1ap_state_p = get_s1ap_state(false);
  hash_table_ts_t* s1ap_ue_state = get_s1ap_ue_state();
  std::vector<uint32_t> mme_ue_id_no_imsi_list = {};
  if (!s1ap_state_p || (s1ap_state_p->enbs.isEmpty())) {
    return;
  }
  s1ap_imsi_map_t* s1ap_imsi_map = get_s1ap_imsi_map();
  ue_description_t* ue_ref_p = NULL;

  // get each eNB in s1ap_state
  for (auto itr = s1ap_state_p->enbs.map->begin();
       itr != s1ap_state_p->enbs.map->end(); itr++) {
    struct enb_description_s* enb_association_p = itr->second;
    if (!enb_association_p) {
      continue;
    }

    if (enb_association_p->ue_id_coll.isEmpty()) {
      continue;
    }

    // for each ue comp_s1ap_id in eNB->ue_id_coll, check if it has an S1ap
    // ue_context, if not delete it
    for (auto ue_itr = enb_association_p->ue_id_coll.map->begin();
         ue_itr != enb_association_p->ue_id_coll.map->end(); ue_itr++) {
      // Check if a UE reference exists for this comp_s1ap_id
      hashtable_ts_get(s1ap_ue_state, (const hash_key_t)ue_itr->second,
                       reinterpret_cast<void**>(&ue_ref_p));
      if (!ue_ref_p) {
        mme_ue_id_no_imsi_list.push_back(ue_itr->first);
        OAILOG_DEBUG(LOG_S1AP,
                     "Adding mme_ue_s1ap_id %u to eNB clean up list with "
                     "num_ues_checked "
                     "%lu",
                     ue_itr->first, mme_ue_id_no_imsi_list.size());
      }
    }
    // remove all the mme_ue_s1ap_ids
    for (uint32_t i = 0; i < mme_ue_id_no_imsi_list.size(); i++) {
      enb_association_p->ue_id_coll.remove(mme_ue_id_no_imsi_list[i]);

      s1ap_imsi_map->mme_ueid2imsi_map.remove(mme_ue_id_no_imsi_list[i]);
      enb_association_p->nb_ue_associated--;

      OAILOG_DEBUG(LOG_S1AP,
                   "Num UEs associated %u num elements in ue_id_coll %zu",
                   enb_association_p->nb_ue_associated,
                   enb_association_p->ue_id_coll.size());
    }
  }
}
