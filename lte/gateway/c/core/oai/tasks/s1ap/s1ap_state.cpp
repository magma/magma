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

namespace magma {
namespace lte {

int s1ap_state_init(bool use_stateless) {
  S1apStateManager::getInstance().init(use_stateless);
  // remove UEs with unknown IMSI from eNB state
  remove_ues_without_imsi_from_ue_id_coll();
  return RETURNok;
}

oai::S1apState* get_s1ap_state(bool read_from_db) {
  return S1apStateManager::getInstance().get_state(read_from_db);
}

void s1ap_state_exit() { S1apStateManager::getInstance().free_state(); }

void put_s1ap_state() {
  S1apStateManager::getInstance().write_s1ap_state_to_db();
}

proto_map_rc_t s1ap_state_get_enb(oai::S1apState* state,
                                  sctp_assoc_id_t assoc_id,
                                  oai::EnbDescription* enb) {
  proto_map_uint32_enb_description_t enb_map;
  enb_map.map = state->mutable_enbs();
  return enb_map.get(assoc_id, enb);
}

proto_map_rc_t s1ap_state_update_enb_map(oai::S1apState* state,
                                         sctp_assoc_id_t assoc_id,
                                         oai::EnbDescription* enb) {
  proto_map_uint32_enb_description_t enb_map;
  enb_map.map = state->mutable_enbs();
  return enb_map.update_val(assoc_id, enb);
}

oai::UeDescription* s1ap_state_get_ue_enbid(sctp_assoc_id_t sctp_assoc_id,
                                            enb_ue_s1ap_id_t enb_ue_s1ap_id) {
  oai::UeDescription* ue = nullptr;

  map_uint64_ue_description_t* state_ue_map = get_s1ap_ue_state();
  if (!state_ue_map) {
    OAILOG_ERROR(LOG_S1AP, "Failed to get s1ap_ue_state");
    return ue;
  }
  uint64_t comp_s1ap_id =
      S1AP_GENERATE_COMP_S1AP_ID(sctp_assoc_id, enb_ue_s1ap_id);
  state_ue_map->get(comp_s1ap_id, &ue);

  return ue;
}

oai::UeDescription* s1ap_state_get_ue_mmeid(mme_ue_s1ap_id_t mme_ue_s1ap_id) {
  oai::UeDescription* ue = nullptr;

  map_uint64_ue_description_t* state_ue_map = get_s1ap_ue_state();
  if (!state_ue_map) {
    OAILOG_ERROR(LOG_S1AP, "Failed to get s1ap_ue_state");
    return ue;
  }
  state_ue_map->map_apply_callback_on_all_elements(
      s1ap_ue_compare_by_mme_ue_id_cb, reinterpret_cast<void*>(&mme_ue_s1ap_id),
      reinterpret_cast<void**>(&ue));

  return ue;
}

oai::UeDescription* s1ap_state_get_ue_imsi(imsi64_t imsi64) {
  oai::UeDescription* ue = nullptr;

  map_uint64_ue_description_t* state_ue_map = get_s1ap_ue_state();
  if (!state_ue_map) {
    OAILOG_ERROR(LOG_S1AP, "Failed to get s1ap_ue_state");
    return ue;
  }
  state_ue_map->map_apply_callback_on_all_elements(
      s1ap_ue_compare_by_imsi, reinterpret_cast<void*>(&imsi64),
      reinterpret_cast<void**>(&ue));

  return ue;
}

void put_s1ap_imsi_map() {
  S1apStateManager::getInstance().write_s1ap_imsi_map_to_db();
}

oai::S1apImsiMap* get_s1ap_imsi_map() {
  return S1apStateManager::getInstance().get_s1ap_imsi_map();
}

bool s1ap_ue_compare_by_mme_ue_id_cb(__attribute__((unused)) uint64_t keyP,
                                     oai::UeDescription* elementP,
                                     void* parameterP, void** resultP) {
  mme_ue_s1ap_id_t* mme_ue_s1ap_id_p = (mme_ue_s1ap_id_t*)parameterP;
  oai::UeDescription* ue_ref = reinterpret_cast<oai::UeDescription*>(elementP);
  if (*mme_ue_s1ap_id_p == ue_ref->mme_ue_s1ap_id()) {
    *resultP = elementP;
    OAILOG_TRACE(LOG_S1AP,
                 "Found ue_ref %p mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT "\n",
                 ue_ref, ue_ref->mme_ue_s1ap_id());
    return true;
  }
  return false;
}

bool s1ap_ue_compare_by_imsi(__attribute__((unused)) uint64_t keyP,
                             oai::UeDescription* elementP, void* parameterP,
                             void** resultP) {
  imsi64_t imsi64 = INVALID_IMSI64;
  imsi64_t* target_imsi64 = (imsi64_t*)parameterP;
  oai::UeDescription* ue_ref = reinterpret_cast<oai::UeDescription*>(elementP);

  magma::proto_map_uint32_uint64_t ueid_imsi_map;
  get_s1ap_ueid_imsi_map(&ueid_imsi_map);
  ueid_imsi_map.get(ue_ref->mme_ue_s1ap_id(), &imsi64);

  if (*target_imsi64 != INVALID_IMSI64 && *target_imsi64 == imsi64) {
    *resultP = elementP;
    OAILOG_DEBUG_UE(LOG_S1AP, imsi64, "Found ue_ref\n");
    return true;
  }
  return false;
}

map_uint64_ue_description_t* get_s1ap_ue_state(void) {
  return S1apStateManager::getInstance().get_s1ap_ue_state();
}

void put_s1ap_ue_state(imsi64_t imsi64) {
  if (S1apStateManager::getInstance().is_persist_state_enabled()) {
    oai::UeDescription* ue_ctxt = s1ap_state_get_ue_imsi(imsi64);
    if (ue_ctxt) {
      auto imsi_str = S1apStateManager::getInstance().get_imsi_str(imsi64);
      S1apStateManager::getInstance().s1ap_write_ue_state_to_db(ue_ctxt,
                                                                imsi_str);
    }
  }
}

void delete_s1ap_ue_state(imsi64_t imsi64) {
  auto imsi_str = S1apStateManager::getInstance().get_imsi_str(imsi64);
  S1apStateManager::getInstance().clear_ue_state_db(imsi_str);
}

void remove_ues_without_imsi_from_ue_id_coll() {
  oai::S1apState* s1ap_state_p = get_s1ap_state(false);
  if (!s1ap_state_p) {
    OAILOG_ERROR(LOG_S1AP, "Failed to get s1ap_state");
    return;
  }

  map_uint64_ue_description_t* s1ap_ue_state = get_s1ap_ue_state();
  if (!(s1ap_ue_state)) {
    OAILOG_ERROR(LOG_S1AP, "Failed to get s1ap_ue_state");
    return;
  }
  proto_map_uint32_enb_description_t enb_map;
  enb_map.map = s1ap_state_p->mutable_enbs();
  if ((enb_map.isEmpty())) {
    return;
  }

  std::vector<uint32_t> mme_ue_id_no_imsi_list = {};
  magma::proto_map_uint32_uint64_t ueid_imsi_map;
  oai::S1apImsiMap* s1ap_imsi_map = get_s1ap_imsi_map();
  if (!s1ap_imsi_map) {
    OAILOG_ERROR(LOG_S1AP, "Failed to get s1ap_imsi_map");
    return;
  }
  oai::UeDescription* ue_ref_p = nullptr;

  // get each eNB in s1ap_state
  for (auto itr = enb_map.map->begin(); itr != enb_map.map->end(); itr++) {
    struct oai::EnbDescription enb_association_p = itr->second;
    if (!enb_association_p.sctp_assoc_id()) {
      continue;
    }

    magma::proto_map_uint32_uint64_t ue_id_coll;
    ue_id_coll.map = enb_association_p.mutable_ue_id_map();
    if (ue_id_coll.isEmpty()) {
      continue;
    }

    // for each ue comp_s1ap_id in eNB->ue_id_coll, check if it has an S1ap
    // ue_context, if not delete it
    mme_ue_id_no_imsi_list.clear();
    for (auto ue_itr = ue_id_coll.map->begin(); ue_itr != ue_id_coll.map->end();
         ue_itr++) {
      // Check if a UE reference exists for this comp_s1ap_id
      s1ap_ue_state->get(ue_itr->second, &ue_ref_p);
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
    ueid_imsi_map.map = s1ap_imsi_map->mutable_mme_ue_s1ap_id_imsi_map();
    for (uint32_t i = 0; i < mme_ue_id_no_imsi_list.size(); i++) {
      ue_id_coll.remove(mme_ue_id_no_imsi_list[i]);

      ueid_imsi_map.remove(mme_ue_id_no_imsi_list[i]);
      DevAssert(enb_association_p.nb_ue_associated() > 0);
      enb_association_p.set_nb_ue_associated(
          (enb_association_p.nb_ue_associated() - 1));

      OAILOG_DEBUG(LOG_S1AP,
                   "Num UEs associated %u num elements in ue_id_coll %zu",
                   enb_association_p.nb_ue_associated(), ue_id_coll.size());
    }
    s1ap_state_update_enb_map(s1ap_state_p, enb_association_p.sctp_assoc_id(),
                              &enb_association_p);
  }
}

}  // namespace lte
}  // namespace magma
