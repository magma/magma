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
 *------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_36.413.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_state_manager.hpp"
#include "lte/gateway/c/core/oai/include/proto_map.hpp"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme.hpp"

namespace {
constexpr char S1AP_ENB_COLL[] = "s1ap_eNB_coll";
constexpr char S1AP_MME_ID2ASSOC_ID_COLL[] = "s1ap_mme_id2assoc_id_coll";
constexpr char S1AP_MME_UEID2IMSI_MAP[] = "s1ap_mme_ueid2imsi_map";
constexpr char S1AP_IMSI_MAP_TABLE_NAME[] = "s1ap_imsi_map";
}  // namespace

using magma::lte::oai::UeDescription;

namespace magma {
namespace lte {

S1apStateManager::S1apStateManager()
    : max_ues_(0),
      max_enbs_(0),
      s1ap_imsi_map_hash_(0),
      s1ap_imsi_map_(nullptr) {}

S1apStateManager::~S1apStateManager() { free_state(); }

S1apStateManager& S1apStateManager::getInstance() {
  static S1apStateManager instance;
  return instance;
}

void S1apStateManager::init(uint32_t max_ues, uint32_t max_enbs,
                            bool persist_state) {
  log_task = LOG_S1AP;
  table_key = S1AP_STATE_TABLE;
  task_name = S1AP_TASK_NAME;
  persist_state_enabled = persist_state;
  max_ues_ = max_ues;
  max_enbs_ = max_enbs;
  redis_client = std::make_unique<RedisClient>(persist_state);
  create_state();
  if (read_state_from_db() != RETURNok) {
    OAILOG_ERROR(LOG_S1AP, "Failed to read state from redis");
  }
  read_ue_state_from_db();
  is_initialized = true;
}

s1ap_state_t* create_s1ap_state(void) {
  bstring ht_name;

  s1ap_state_t* state_cache_p = new s1ap_state_t();
  state_cache_p->enbs.map =
      new google::protobuf::Map<unsigned int, struct enb_description_s*>();
  state_cache_p->enbs.set_name(S1AP_ENB_COLL);
  state_cache_p->enbs.bind_callback(free_cpp_wrapper);

  state_cache_p->mmeid2associd.map =
      new google::protobuf::Map<uint32_t, uint32_t>();
  state_cache_p->mmeid2associd.set_name(S1AP_MME_ID2ASSOC_ID_COLL);

  state_cache_p->num_enbs = 0;
  return state_cache_p;
}

void S1apStateManager::create_state() {
  state_cache_p = create_s1ap_state();

  bstring ht_name = bfromcstr(S1AP_ENB_COLL);
  state_ue_ht = hashtable_ts_create(max_ues_, nullptr, free_wrapper, ht_name);
  bdestroy(ht_name);

  create_s1ap_imsi_map();
}

void free_s1ap_state(s1ap_state_t* state_cache_p) {
  AssertFatal(state_cache_p,
              "s1ap_state_t passed to free_s1ap_state must not be null");

  int i;
  hashtable_rc_t ht_rc;
  hashtable_key_array_t* keys;
  sctp_assoc_id_t assoc_id;
  enb_description_t* enb;

  if (state_cache_p->enbs.isEmpty()) {
    OAILOG_DEBUG(LOG_S1AP, "No keys in the enb hashtable");
  } else {
    for (auto itr = state_cache_p->enbs.map->begin();
         itr != state_cache_p->enbs.map->end(); itr++) {
      enb = itr->second;
      if (!enb) {
        OAILOG_ERROR(LOG_S1AP, "eNB entry not found in eNB S1AP state");
      } else {
        enb->ue_id_coll.destroy_map();
      }
    }
  }
  if (state_cache_p->enbs.destroy_map() != PROTO_MAP_OK) {
    OAILOG_ERROR(LOG_S1AP, "An error occurred while destroying s1 eNB map");
  }
  if ((state_cache_p->mmeid2associd.destroy_map()) != magma::PROTO_MAP_OK) {
    OAILOG_ERROR(LOG_S1AP,
                 "An error occurred while destroying mmeid2associd map");
  }
  delete state_cache_p;
}

void S1apStateManager::free_state() {
  AssertFatal(
      is_initialized,
      "S1apStateManager init() function should be called to initialize state.");

  if (state_cache_p == nullptr) {
    return;
  }
  free_s1ap_state(state_cache_p);
  state_cache_p = nullptr;

  if (hashtable_ts_destroy(state_ue_ht) != HASH_TABLE_OK) {
    OAILOG_ERROR(LOG_S1AP,
                 "An error occurred while destroying assoc_id hash table");
  }
  clear_s1ap_imsi_map();
}

status_code_e S1apStateManager::read_ue_state_from_db() {
#if !MME_UNIT_TEST
  if (!persist_state_enabled) {
    return RETURNok;
  }
  auto keys = redis_client->get_keys("IMSI*" + task_name + "*");

  for (const auto& key : keys) {
    OAILOG_DEBUG(log_task, "Reading UE state from db for %s", key.c_str());
    UeDescription ue_proto = UeDescription();
    auto* ue_context = (ue_description_t*)calloc(1, sizeof(ue_description_t));
    if (redis_client->read_proto(key, ue_proto) != RETURNok) {
      return RETURNerror;
    }

    S1apStateConverter::proto_to_ue(ue_proto, ue_context);

    hashtable_rc_t h_rc = hashtable_ts_insert(
        state_ue_ht, ue_context->comp_s1ap_id, (void*)ue_context);
    if (HASH_TABLE_OK != h_rc) {
      OAILOG_ERROR(
          log_task,
          "Failed to insert UE state with key comp_s1ap_id " COMP_S1AP_ID_FMT
          ", ENB UE S1AP Id: " ENB_UE_S1AP_ID_FMT
          ", MME UE S1AP Id: " MME_UE_S1AP_ID_FMT " (Error Code: %s)\n",
          ue_context->comp_s1ap_id, ue_context->enb_ue_s1ap_id,
          ue_context->mme_ue_s1ap_id, hashtable_rc_code2string(h_rc));
    } else {
      OAILOG_DEBUG(log_task,
                   "Inserted UE state with key comp_s1ap_id " COMP_S1AP_ID_FMT
                   ", ENB UE S1AP Id: " ENB_UE_S1AP_ID_FMT
                   ", MME UE S1AP Id: " MME_UE_S1AP_ID_FMT,
                   ue_context->comp_s1ap_id, ue_context->enb_ue_s1ap_id,
                   ue_context->mme_ue_s1ap_id);
    }
  }
#endif
  return RETURNok;
}

void S1apStateManager::create_s1ap_imsi_map() {
  s1ap_imsi_map_ = new s1ap_imsi_map_t();

  s1ap_imsi_map_->mme_ueid2imsi_map.map =
      new google::protobuf::Map<uint32_t, uint64_t>();
  s1ap_imsi_map_->mme_ueid2imsi_map.set_name(S1AP_MME_UEID2IMSI_MAP);

  if (persist_state_enabled) {
    oai::S1apImsiMap imsi_proto = oai::S1apImsiMap();
    redis_client->read_proto(S1AP_IMSI_MAP_TABLE_NAME, imsi_proto);

    S1apStateConverter::proto_to_s1ap_imsi_map(imsi_proto, s1ap_imsi_map_);
  }
}

void S1apStateManager::clear_s1ap_imsi_map() {
  if (!s1ap_imsi_map_) {
    return;
  }
  s1ap_imsi_map_->mme_ueid2imsi_map.destroy_map();
  delete s1ap_imsi_map_;
}

s1ap_imsi_map_t* S1apStateManager::get_s1ap_imsi_map() {
  return s1ap_imsi_map_;
}

void S1apStateManager::write_s1ap_imsi_map_to_db() {
  if (!persist_state_enabled) {
    return;
  }
  oai::S1apImsiMap imsi_proto = oai::S1apImsiMap();
  S1apStateConverter::s1ap_imsi_map_to_proto(s1ap_imsi_map_, &imsi_proto);
  std::string proto_msg;
  redis_client->serialize(imsi_proto, proto_msg);
  std::size_t new_hash = std::hash<std::string>{}(proto_msg);

  // s1ap_imsi_map is not state service synced, so version will not be updated
  if (new_hash != this->s1ap_imsi_map_hash_) {
    redis_client->write_proto_str(S1AP_IMSI_MAP_TABLE_NAME, proto_msg, 0);
    this->s1ap_imsi_map_hash_ = new_hash;
  }
}

}  // namespace lte
}  // namespace magma
