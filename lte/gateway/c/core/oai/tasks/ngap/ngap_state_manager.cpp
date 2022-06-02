/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/****************************************************************************
  Source      ngap_state_manager.cpp
  Date        2020/07/28
  Author      Ashish Prajapati
  Subsystem   Access and Mobility Management Function
  Description Defines NG Application Protocol Messages

*****************************************************************************/

#include "lte/gateway/c/core/oai/tasks/ngap/ngap_state_manager.hpp"
extern "C" {
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
}

typedef unsigned int uint32_t;
namespace {
constexpr char NGAP_GNB_COLL[] = "ngap_gNB_coll";
constexpr char NGAP_AMF_ID2ASSOC_ID_COLL[] = "ngap_amf_id2assoc_id_coll";
constexpr char NGAP_IMSI_MAP_TABLE_NAME[] = "ngap_imsi_map";
}  // namespace

using magma::lte::oai::Ngap_UeDescription;
using magma::lte::oai::NgapImsiMap;
using magma::lte::oai::NgapState;

namespace magma5g {

NgapStateManager::NgapStateManager() : max_ues_(0), max_gnbs_(0) {}

NgapStateManager::~NgapStateManager() { free_state(); }

NgapStateManager& NgapStateManager::getInstance() {
  static NgapStateManager instance;
  return instance;
}

void NgapStateManager::init(uint32_t max_ues, uint32_t max_gnbs,
                            bool use_stateless) {
  log_task = LOG_NGAP;
  table_key = NGAP_STATE_TABLE;
  task_name = NGAP_TASK_NAME;
  persist_state_enabled = use_stateless;
  max_ues_ = max_ues;
  max_gnbs_ = max_gnbs;

  OAILOG_FUNC_IN(LOG_NGAP);
#if !MME_UNIT_TEST
  redis_client = std::make_unique<RedisClient>(persist_state_enabled);
#endif
  create_state();
  if (read_state_from_db() != RETURNok) {
    OAILOG_ERROR(LOG_NGAP, "Failed to read state from redis");
  }
  read_ue_state_from_db();
  is_initialized = true;
  OAILOG_FUNC_OUT(LOG_NGAP);
}

ngap_state_t* create_ngap_state(uint32_t max_gnbs, uint32_t max_ues) {
  bstring ht_name;

  OAILOG_FUNC_IN(LOG_NGAP);
  ngap_state_t* state_cache_p =
      static_cast<ngap_state_t*>(calloc(1, sizeof(ngap_state_t)));

  ht_name = bfromcstr(NGAP_GNB_COLL);
  hashtable_ts_init(&state_cache_p->gnbs, max_gnbs, nullptr, free_wrapper,
                    ht_name);
  bdestroy(ht_name);

  ht_name = bfromcstr(NGAP_AMF_ID2ASSOC_ID_COLL);
  hashtable_ts_init(&state_cache_p->amfid2associd, max_ues, nullptr,
                    hash_free_int_func, ht_name);
  bdestroy(ht_name);

  state_cache_p->num_gnbs = 0;
  OAILOG_FUNC_RETURN(LOG_NGAP, state_cache_p);
}

void NgapStateManager::create_state() {
  OAILOG_FUNC_IN(LOG_NGAP);
  state_cache_p = create_ngap_state(max_gnbs_, max_ues_);

  bstring ht_name;
  ht_name = bfromcstr(NGAP_GNB_COLL);
  state_ue_ht = hashtable_ts_create(max_ues_, nullptr, free_wrapper, ht_name);
  bdestroy(ht_name);

  create_ngap_imsi_map();
  OAILOG_FUNC_OUT(LOG_NGAP);
}

status_code_e NgapStateManager::read_state_from_db() {
#if !MME_UNIT_TEST
  /* Data store is Redis db. In this case actual call is made to Redis db */
  StateManager::read_state_from_db();
#else
  /* Data store is a map defined in NGAPClientServicer.In this case call is NOT
   * made to Redis db */
  NgapState state_proto = NgapState();
  std::string proto_str;

  OAILOG_FUNC_IN(LOG_NGAP);
  // Reads from the map_ngap_state_proto_str Map
  if (NGAPClientServicer::getInstance().map_ngap_state_proto_str.get(
          table_key, &proto_str) != magma::MAP_OK) {
    OAILOG_DEBUG(log_task, "Failed to read proto from db \n");
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }
  // Deserialization Step
  if (!state_proto.ParseFromString(proto_str)) {
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }
  NgapStateConverter::proto_to_state(state_proto, state_cache_p);
#endif
  OAILOG_FUNC_RETURN(LOG_NGAP, RETURNok);
}

void NgapStateManager::write_state_to_db() {
#if !MME_UNIT_TEST
  /* Data store is Redis db. In this case actual call is made to Redis db */
  StateManager::write_state_to_db();
#else
  /* Data store is a map defined in NGAPClientServicer. In this case call is NOT
   * made to Redis db */
  NgapState state_proto = NgapState();
  NgapStateConverter::state_to_proto(state_cache_p, &state_proto);
  std::string proto_str;
  redis_client->serialize(state_proto, proto_str);
  std::size_t new_hash = std::hash<std::string>{}(proto_str);

  OAILOG_FUNC_IN(LOG_NGAP);
  if (new_hash != this->task_state_hash) {
    // Writes to the map_ngap_state_proto_str Map
    if (NGAPClientServicer::getInstance().map_ngap_state_proto_str.insert(
            table_key, proto_str) != magma::MAP_OK) {
      OAILOG_ERROR(log_task, "Failed to write state to db");
      OAILOG_FUNC_OUT(LOG_NGAP);
    }
    OAILOG_DEBUG(log_task, "Finished writing state");
    this->task_state_version++;
    this->state_dirty = false;
    this->task_state_hash = new_hash;
  }
#endif
  OAILOG_FUNC_OUT(LOG_NGAP);
}

void free_ngap_state(ngap_state_t* state_cache_p) {
  int i;
  hashtable_rc_t ht_rc;
  hashtable_key_array_t* keys;
  sctp_assoc_id_t assoc_id;
  gnb_description_t* gnb;

  keys = hashtable_ts_get_keys(&state_cache_p->gnbs);

  OAILOG_FUNC_IN(LOG_NGAP);
  if (!keys) {
    OAILOG_DEBUG(LOG_NGAP, "No keys in the amf hashtable");
  } else {
    for (i = 0; i < keys->num_keys; i++) {
      assoc_id = (sctp_assoc_id_t)keys->keys[i];
      ht_rc = hashtable_ts_get(&state_cache_p->gnbs, (hash_key_t)assoc_id,
                               (void**)&gnb);
      if (ht_rc != HASH_TABLE_OK) {
        OAILOG_ERROR(LOG_NGAP, "gNB entry not found in gNB NGP state");
      } else {
        hashtable_uint64_ts_destroy(&gnb->ue_id_coll);
      }

      AssertFatal(ht_rc == HASH_TABLE_OK, "eNB UE id not in assoc_id");
    }
    FREE_HASHTABLE_KEY_ARRAY(keys);
  }

  if (hashtable_ts_destroy(&state_cache_p->gnbs) != HASH_TABLE_OK) {
    OAI_FPRINTF_ERR("An error occurred while destroying s1 eNB hash table");
  }
  if (hashtable_ts_destroy(&state_cache_p->amfid2associd) != HASH_TABLE_OK) {
    OAI_FPRINTF_ERR("An error occurred while destroying assoc_id hash table");
  }

  free(state_cache_p);
  OAILOG_FUNC_OUT(LOG_NGAP);
}

void NgapStateManager::free_state() {
  AssertFatal(
      is_initialized,
      "NgapStateManager init() function should be called to initialize state.");

  OAILOG_FUNC_IN(LOG_NGAP);
  if (state_cache_p == nullptr) {
    OAILOG_FUNC_OUT(LOG_NGAP);
  }

  free_ngap_state(state_cache_p);
  state_cache_p = nullptr;

  if (hashtable_ts_destroy(state_ue_ht) != HASH_TABLE_OK) {
    OAI_FPRINTF_ERR("An error occurred while destroying assoc_id hash table");
  }

  clear_ngap_imsi_map();
  OAILOG_FUNC_OUT(LOG_NGAP);
}

status_code_e NgapStateManager::read_ue_state_from_db() {
  OAILOG_FUNC_IN(LOG_NGAP);
  if (!persist_state_enabled) {
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNok);
  }
#if !MME_UNIT_TEST
  /* Data store is Redis db. In this case actual call is made to Redis db */
  auto keys = redis_client->get_keys("IMSI*" + task_name + "*");

  for (const auto& key : keys) {
    Ngap_UeDescription ue_proto = Ngap_UeDescription();
    m5g_ue_description_t* ue_context =
        (m5g_ue_description_t*)calloc(1, sizeof(m5g_ue_description_t));
    if (redis_client->read_proto(key.c_str(), ue_proto) != RETURNok) {
      OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
    }

    NgapStateConverter::proto_to_ue(ue_proto, ue_context);

    hashtable_ts_insert(state_ue_ht, ue_context->comp_ngap_id,
                        (void*)ue_context);
    OAILOG_DEBUG(log_task, "Reading UE state from db for %s", key.c_str());
  }
#else
  /* Data store is a map defined in NGAPClientServicer. In this case call is NOT
   * made to Redis db */
  for (const auto& kv :
       NGAPClientServicer::getInstance().map_ngap_uestate_proto_str.umap) {
    Ngap_UeDescription ue_proto = Ngap_UeDescription();
    std::string ue_proto_str;
    m5g_ue_description_t* ue_context = reinterpret_cast<m5g_ue_description_t*>(
        calloc(1, sizeof(m5g_ue_description_t)));
    // Reads from the map_ngap_uestate_proto_str Map
    if (NGAPClientServicer::getInstance().map_ngap_uestate_proto_str.get(
            kv.first, &ue_proto_str) != magma::MAP_OK) {
      OAILOG_DEBUG(log_task, "Failed to read UE proto from db \n");
      free(ue_context);
      OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
    }
    // Deserialization Step
    if (!ue_proto.ParseFromString(ue_proto_str)) {
      free(ue_context);
      OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
    }

    NgapStateConverter::proto_to_ue(ue_proto, ue_context);

    hashtable_ts_insert(state_ue_ht, ue_context->comp_ngap_id,
                        reinterpret_cast<void*>(ue_context));
    OAILOG_DEBUG(log_task, "Reading UE state from db");
  }
#endif
  OAILOG_FUNC_RETURN(LOG_NGAP, RETURNok);
}

void NgapStateManager::write_ue_state_to_db(
    const m5g_ue_description_t* ue_context, const std::string& imsi_str) {
#if !MME_UNIT_TEST
  /* Data store is Redis db. In this case actual call is made to Redis db */
  StateManager::write_ue_state_to_db(ue_context, imsi_str);
#else
  /* Data store is a map defined in NGAPClientServicer. In this case call is NOT
   * made to Redis db */
  std::string proto_ue_str;
  Ngap_UeDescription ue_proto = Ngap_UeDescription();
  NgapStateConverter::ue_to_proto(ue_context, &ue_proto);
  redis_client->serialize(ue_proto, proto_ue_str);
  std::size_t new_hash = std::hash<std::string>{}(proto_ue_str);

  OAILOG_FUNC_IN(LOG_NGAP);
  if (new_hash != this->ue_state_hash[imsi_str]) {
    std::string key = IMSI_PREFIX + imsi_str + ":" + task_name;
    // Writes to the map_ngap_uestate_proto_str Map
    if (NGAPClientServicer::getInstance().map_ngap_uestate_proto_str.insert(
            key, proto_ue_str) != magma::MAP_OK) {
      OAILOG_ERROR(log_task, "Failed to write UE state to db for IMSI %s",
                   imsi_str.c_str());
      OAILOG_FUNC_OUT(LOG_NGAP);
    }

    this->ue_state_version[imsi_str]++;
    this->ue_state_hash[imsi_str] = new_hash;
    OAILOG_DEBUG(log_task, "Finished writing UE state for IMSI %s",
                 imsi_str.c_str());
  }
#endif
  OAILOG_FUNC_OUT(LOG_NGAP);
}

void NgapStateManager::create_ngap_imsi_map() {
  ngap_imsi_map_ = (ngap_imsi_map_t*)calloc(1, sizeof(ngap_imsi_map_t));

  OAILOG_FUNC_IN(LOG_NGAP);
  ngap_imsi_map_->amf_ue_id_imsi_htbl =
      hashtable_uint64_ts_create(max_ues_, nullptr, nullptr);
  if (!persist_state_enabled) {
    OAILOG_FUNC_OUT(LOG_NGAP);
  }
  NgapImsiMap imsi_proto = NgapImsiMap();
#if !MME_UNIT_TEST
  /* Data store is Redis db. In this case actual call is made to Redis db */
  redis_client->read_proto(NGAP_IMSI_MAP_TABLE_NAME, imsi_proto);
#else
  /* Data store is a map defined in NGAPClientServicer.In this case call is NOT
   * made to Redis db */
  // Reads from the NGAPClientServicer DataStore Map
  // (map_imsi_table_proto_str)
  std::string proto_msg;
  if (NGAPClientServicer::getInstance().map_imsi_table_proto_str.get(
          NGAP_IMSI_MAP_TABLE_NAME, &proto_msg) != magma::MAP_OK) {
    OAILOG_DEBUG(log_task, "Failed to read ngap_imsi_map proto from db \n");
    OAILOG_FUNC_OUT(LOG_NGAP);
  }
  // Deserialization Step
  if (!imsi_proto.ParseFromString(proto_msg)) {
    OAILOG_FUNC_OUT(LOG_NGAP);
  }
#endif

  NgapStateConverter::proto_to_ngap_imsi_map(imsi_proto, ngap_imsi_map_);
  OAILOG_FUNC_OUT(LOG_NGAP);
}

void NgapStateManager::clear_ngap_imsi_map() {
  OAILOG_FUNC_IN(LOG_NGAP);
  if (!ngap_imsi_map_) {
    OAILOG_FUNC_OUT(LOG_NGAP);
  }
  hashtable_uint64_ts_destroy(ngap_imsi_map_->amf_ue_id_imsi_htbl);

  free_wrapper((void**)&ngap_imsi_map_);
  OAILOG_FUNC_OUT(LOG_NGAP);
}

ngap_imsi_map_t* NgapStateManager::get_ngap_imsi_map() {
  return ngap_imsi_map_;
}

void NgapStateManager::put_ngap_imsi_map() {
  OAILOG_FUNC_IN(LOG_NGAP);
  if (!persist_state_enabled) {
    OAILOG_FUNC_OUT(LOG_NGAP);
  }
  NgapImsiMap imsi_proto = NgapImsiMap();
  NgapStateConverter::ngap_imsi_map_to_proto(ngap_imsi_map_, &imsi_proto);
  std::string proto_msg;
  redis_client->serialize(imsi_proto, proto_msg);
  std::size_t new_hash = std::hash<std::string>{}(proto_msg);
  if (new_hash != this->ngap_imsi_map_hash_) {
#if !MME_UNIT_TEST
    /* Data store is Redis db. In this case actual call is made to Redis db */
    redis_client->write_proto_str(NGAP_IMSI_MAP_TABLE_NAME, proto_msg, 0);
#else
    /* Data store is a map defined in NGAPClientServicer.In this case call is
     * NOT made to Redis db */
    // Writes to the map_imsi_table_proto_str Map
    if (NGAPClientServicer::getInstance().map_imsi_table_proto_str.insert(
            NGAP_IMSI_MAP_TABLE_NAME, proto_msg) != magma::MAP_OK) {
      OAILOG_ERROR(log_task, "Failed to write ngap_imsi_map state to db \n");
      OAILOG_FUNC_OUT(LOG_NGAP);
    }
#endif
    this->ngap_imsi_map_hash_ = new_hash;
  }
  OAILOG_FUNC_OUT(LOG_NGAP);
}

}  // namespace magma5g
