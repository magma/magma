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
constexpr char S1AP_STATE_UE_MAP[] = "s1ap_state_ue_map";
}  // namespace

namespace magma {
namespace lte {

S1apStateManager::S1apStateManager()
    : s1ap_imsi_map_hash_(0), s1ap_imsi_map_(nullptr) {}

S1apStateManager::~S1apStateManager() { free_state(); }

S1apStateManager& S1apStateManager::getInstance() {
  static S1apStateManager instance;
  return instance;
}

void S1apStateManager::init(bool persist_state) {
  log_task = LOG_S1AP;
  table_key = S1AP_STATE_TABLE;
  task_name = S1AP_TASK_NAME;
  persist_state_enabled = persist_state;
  redis_client = std::make_unique<RedisClient>(persist_state);
  create_state();
  if (read_state_from_db() != RETURNok) {
    OAILOG_ERROR(LOG_S1AP, "Failed to read state from redis");
  }
  read_ue_state_from_db();
  is_initialized = true;
}

oai::S1apState* create_s1ap_state(void) {
  proto_map_uint32_enb_description_t enb_map;

  oai::S1apState* state_cache_p = new oai::S1apState();
  enb_map.map = state_cache_p->mutable_enbs();
  enb_map.set_name(S1AP_ENB_COLL);
  enb_map.bind_callback(free_enb_description);

  magma::proto_map_uint32_uint32_t mmeid2associd;
  mmeid2associd.map = state_cache_p->mutable_mmeid2associd();
  mmeid2associd.set_name(S1AP_MME_ID2ASSOC_ID_COLL);

  return state_cache_p;
}

void S1apStateManager::create_state() {
  state_cache_p = create_s1ap_state();
  if (!state_cache_p) {
    OAILOG_ERROR(LOG_S1AP, "Failed to create s1ap state");
    return;
  }

  state_ue_map.map = new google::protobuf::Map<uint64_t, oai::UeDescription*>();
  if (!(state_ue_map.map)) {
    OAILOG_ERROR(LOG_S1AP, "Failed to allocate memory for state_ue_map ");
    return;
  }
  state_ue_map.set_name(S1AP_STATE_UE_MAP);
  state_ue_map.bind_callback(free_ue_description);

  create_s1ap_imsi_map();
}

void free_s1ap_state(oai::S1apState* state_cache_p) {
  AssertFatal(state_cache_p,
              "S1apState passed to free_s1ap_state must not be null");

  proto_map_uint32_enb_description_t enb_map;
  enb_map.map = state_cache_p->mutable_enbs();

  if (enb_map.isEmpty()) {
    OAILOG_DEBUG(LOG_S1AP, "No keys in the enb map");
  } else {
    for (auto itr = enb_map.map->begin(); itr != enb_map.map->end(); itr++) {
      oai::EnbDescription enb = itr->second;
      enb.clear_ue_id_map();
    }
  }
  state_cache_p->Clear();
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

  if (state_ue_map.destroy_map() != PROTO_MAP_OK) {
    OAILOG_ERROR(LOG_S1AP, "An error occurred while destroying state_ue_map");
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
    oai::UeDescription ue_proto = oai::UeDescription();
    auto* ue_context = new oai::UeDescription();
    if (!ue_context) {
      OAILOG_ERROR(log_task, "Failed to allocate memory for ue context");
      return RETURNerror;
    }
    if (redis_client->read_proto(key, ue_proto) != RETURNok) {
      return RETURNerror;
    }

    // Update each UE state version from redis
    this->ue_state_version[key] = redis_client->read_version(table_key);

    ue_context->MergeFrom(ue_proto);

    proto_map_rc_t rc =
        state_ue_map.insert(ue_context->comp_s1ap_id(), ue_context);
    if (rc != PROTO_MAP_OK) {
      OAILOG_ERROR(
          log_task,
          "Failed to insert UE state with key comp_s1ap_id " COMP_S1AP_ID_FMT
          ", ENB UE S1AP Id: " ENB_UE_S1AP_ID_FMT
          ", MME UE S1AP Id: " MME_UE_S1AP_ID_FMT " (Error Code: %s)\n",
          ue_context->comp_s1ap_id(), ue_context->enb_ue_s1ap_id(),
          ue_context->mme_ue_s1ap_id(), magma::map_rc_code2string(rc));
    } else {
      OAILOG_DEBUG(log_task,
                   "Inserted UE state with key comp_s1ap_id " COMP_S1AP_ID_FMT
                   ", ENB UE S1AP Id: " ENB_UE_S1AP_ID_FMT
                   ", MME UE S1AP Id: " MME_UE_S1AP_ID_FMT,
                   ue_context->comp_s1ap_id(), ue_context->enb_ue_s1ap_id(),
                   ue_context->mme_ue_s1ap_id());
    }
  }
#endif
  return RETURNok;
}

void S1apStateManager::create_s1ap_imsi_map() {
  proto_map_uint32_uint64_t imsi_map;
  s1ap_imsi_map_ = new oai::S1apImsiMap();

  imsi_map.map = s1ap_imsi_map_->mutable_mme_ue_s1ap_id_imsi_map();
  imsi_map.set_name(S1AP_MME_UEID2IMSI_MAP);

  if (persist_state_enabled) {
    redis_client->read_proto(S1AP_IMSI_MAP_TABLE_NAME, *s1ap_imsi_map_);
  }
}

void S1apStateManager::clear_s1ap_imsi_map() {
  if (!s1ap_imsi_map_) {
    return;
  }
  s1ap_imsi_map_->Clear();
  delete s1ap_imsi_map_;
}

oai::S1apImsiMap* S1apStateManager::get_s1ap_imsi_map() {
  return s1ap_imsi_map_;
}

void S1apStateManager::write_s1ap_imsi_map_to_db() {
  if (!persist_state_enabled) {
    return;
  }
  std::string proto_msg;
  redis_client->serialize(*s1ap_imsi_map_, proto_msg);
  std::size_t new_hash = std::hash<std::string>{}(proto_msg);

  // s1ap_imsi_map is not state service synced, so version will not be updated
  if (new_hash != this->s1ap_imsi_map_hash_) {
    redis_client->write_proto_str(S1AP_IMSI_MAP_TABLE_NAME, proto_msg, 0);
    this->s1ap_imsi_map_hash_ = new_hash;
  }
}

map_uint64_ue_description_t* S1apStateManager::get_s1ap_ue_state() {
  return &state_ue_map;
}

oai::S1apState* S1apStateManager::get_state(bool read_from_db) {
  OAILOG_FUNC_IN(LOG_S1AP);
  AssertFatal(
      is_initialized,
      "S1apStateManager init() function should be called to initialize state");
  state_dirty = true;
  AssertFatal(state_cache_p != nullptr, " S1ap State cache is NULL");
  if (persist_state_enabled && read_from_db) {
    read_state_from_db();
    read_ue_state_from_db();
  }
  OAILOG_FUNC_RETURN(LOG_S1AP, state_cache_p);
}

void S1apStateManager::write_s1ap_state_to_db() {
  AssertFatal(
      is_initialized,
      "S1ap StateManager init() function should be called to initialize state");

  if (!state_dirty) {
    OAILOG_ERROR(log_task, "Tried to put state while it was not in use");
    return;
  }

  if (persist_state_enabled) {
    std::string proto_str;
    redis_client->serialize(*state_cache_p, proto_str);
    std::size_t new_hash = std::hash<std::string>{}(proto_str);

    if (new_hash != this->task_state_hash) {
      if (redis_client->write_proto_str(table_key, proto_str,
                                        this->task_state_version) != RETURNok) {
        OAILOG_ERROR(log_task, "Failed to write state to db");
        return;
      }
      OAILOG_DEBUG(log_task, "Finished writing state");
      this->task_state_version++;
      this->state_dirty = false;
      this->task_state_hash = new_hash;
    }
  }
}

void S1apStateManager::s1ap_write_ue_state_to_db(
    const oai::UeDescription* ue_context, const std::string& imsi_str) {
  AssertFatal(
      is_initialized,
      "StateManager init() function should be called to initialize state");

  std::string proto_str;
  redis_client->serialize(*ue_context, proto_str);
  std::size_t new_hash = std::hash<std::string>{}(proto_str);
  if (new_hash != this->ue_state_hash[imsi_str]) {
    std::string key = IMSI_STR_PREFIX + imsi_str + ":" + task_name;
    if (redis_client->write_proto_str(key, proto_str,
                                      ue_state_version[imsi_str]) != RETURNok) {
      OAILOG_ERROR(log_task, "Failed to write UE state to db for IMSI %s",
                   imsi_str.c_str());
      return;
    }
    this->ue_state_version[imsi_str]++;
    this->ue_state_hash[imsi_str] = new_hash;
    OAILOG_DEBUG(log_task, "Finished writing UE state for IMSI %s",
                 imsi_str.c_str());
  }
}

status_code_e S1apStateManager::read_state_from_db() {
#if !MME_UNIT_TEST
  if (persist_state_enabled) {
    oai::S1apState state_proto = oai::S1apState();
    if (redis_client->read_proto(table_key, state_proto) != RETURNok) {
      OAILOG_ERROR(LOG_MME_APP, "Failed to read proto from db \n");
      return RETURNerror;
    }

    // Update the state version from redis
    this->task_state_version = redis_client->read_version(table_key);

    state_cache_p->Clear();
    state_cache_p->MergeFrom(state_proto);
  }
#endif
  return RETURNok;
}

}  // namespace lte
}  // namespace magma
