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

#include "lte/gateway/c/core/oai/tasks/ngap/ngap_state_manager.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/common/common_defs.h"
typedef unsigned int uint32_t;
namespace {
constexpr char NGAP_GNB_COLL[]             = "ngap_gNB_coll";
constexpr char NGAP_AMF_ID2ASSOC_ID_COLL[] = "ngap_amf_id2assoc_id_coll";
constexpr char NGAP_IMSI_MAP_TABLE_NAME[]  = "ngap_imsi_map";
}  // namespace

using magma::lte::oai::UeDescription;

namespace magma5g {

NgapStateManager::NgapStateManager() : max_ues_(0), max_gnbs_(0) {}

NgapStateManager::~NgapStateManager() {
  free_state();
}

NgapStateManager& NgapStateManager::getInstance() {
  static NgapStateManager instance;
  return instance;
}

void NgapStateManager::init(
    uint32_t max_ues, uint32_t max_gnbs, bool use_stateless) {
  log_task              = LOG_NGAP;
  table_key             = NGAP_STATE_TABLE;
  task_name             = NGAP_TASK_NAME;
  persist_state_enabled = false;
  max_ues_              = max_ues;
  max_gnbs_             = max_gnbs;
  create_state();
  if (read_state_from_db() != RETURNok) {
    OAILOG_ERROR(LOG_NGAP, "Failed to read state from redis");
  }
  read_ue_state_from_db();
  is_initialized = true;
}

ngap_state_t* create_ngap_state(uint32_t max_gnbs, uint32_t max_ues) {
  bstring ht_name;

  ngap_state_t* state_cache_p =
      static_cast<ngap_state_t*>(calloc(1, sizeof(ngap_state_t)));

  ht_name = bfromcstr(NGAP_GNB_COLL);
  hashtable_ts_init(
      &state_cache_p->gnbs, max_gnbs, nullptr, free_wrapper, ht_name);
  bdestroy(ht_name);

  ht_name = bfromcstr(NGAP_AMF_ID2ASSOC_ID_COLL);
  hashtable_ts_init(
      &state_cache_p->amfid2associd, max_ues, nullptr, hash_free_int_func,
      ht_name);
  bdestroy(ht_name);

  state_cache_p->num_gnbs = 0;
  return state_cache_p;
}

void NgapStateManager::create_state() {
  state_cache_p = create_ngap_state(max_gnbs_, max_ues_);

  bstring ht_name;
  ht_name     = bfromcstr(NGAP_GNB_COLL);
  state_ue_ht = hashtable_ts_create(max_ues_, nullptr, free_wrapper, ht_name);
  bdestroy(ht_name);

  create_ngap_imsi_map();
}

void free_ngap_state(ngap_state_t* state_cache_p) {
  int i;
  hashtable_rc_t ht_rc;
  hashtable_key_array_t* keys;
  sctp_assoc_id_t assoc_id;
  gnb_description_t* gnb;

  keys = hashtable_ts_get_keys(&state_cache_p->gnbs);
  if (!keys) {
    OAILOG_DEBUG(LOG_NGAP, "No keys in the amf hashtable");
  } else {
    for (i = 0; i < keys->num_keys; i++) {
      assoc_id = (sctp_assoc_id_t) keys->keys[i];
      ht_rc    = hashtable_ts_get(
          &state_cache_p->gnbs, (hash_key_t) assoc_id, (void**) &gnb);
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
}

void NgapStateManager::free_state() {
  AssertFatal(
      is_initialized,
      "NgapStateManager init() function should be called to initialize state.");

  if (state_cache_p == nullptr) {
    return;
  }

  free_ngap_state(state_cache_p);
  state_cache_p = nullptr;

  if (hashtable_ts_destroy(state_ue_ht) != HASH_TABLE_OK) {
    OAI_FPRINTF_ERR("An error occurred while destroying assoc_id hash table");
  }

  clear_ngap_imsi_map();
}

status_code_e NgapStateManager::read_ue_state_from_db() {
  if (!persist_state_enabled) {
    return RETURNok;
  }
  auto keys = redis_client->get_keys("IMSI*" + task_name + "*");

  for (const auto& key : keys) {
    UeDescription ue_proto = UeDescription();
    m5g_ue_description_t* ue_context =
        (m5g_ue_description_t*) calloc(1, sizeof(m5g_ue_description_t));
    if (redis_client->read_proto(key.c_str(), ue_proto) != RETURNok) {
      return RETURNerror;
    }

    NgapStateConverter::proto_to_ue(ue_proto, ue_context);

    hashtable_ts_insert(
        state_ue_ht, ue_context->comp_ngap_id, (void*) ue_context);
    OAILOG_DEBUG(log_task, "Reading UE state from db for %s", key.c_str());
  }
  return RETURNok;
}

void NgapStateManager::create_ngap_imsi_map() {
  ngap_imsi_map_ = (ngap_imsi_map_t*) calloc(1, sizeof(ngap_imsi_map_t));

  ngap_imsi_map_->amf_ue_id_imsi_htbl =
      hashtable_uint64_ts_create(max_ues_, nullptr, nullptr);
  if (!persist_state_enabled) {
    return;
  }
  oai::S1apImsiMap imsi_proto = oai::S1apImsiMap();
  redis_client->read_proto(NGAP_IMSI_MAP_TABLE_NAME, imsi_proto);

  NgapStateConverter::proto_to_ngap_imsi_map(imsi_proto, ngap_imsi_map_);
}

void NgapStateManager::clear_ngap_imsi_map() {
  if (!ngap_imsi_map_) {
    return;
  }
  hashtable_uint64_ts_destroy(ngap_imsi_map_->amf_ue_id_imsi_htbl);

  free_wrapper((void**) &ngap_imsi_map_);
}

ngap_imsi_map_t* NgapStateManager::get_ngap_imsi_map() {
  return ngap_imsi_map_;
}

void NgapStateManager::put_ngap_imsi_map() {
  if (!persist_state_enabled) {
    return;
  }
  oai::S1apImsiMap imsi_proto = oai::S1apImsiMap();
  NgapStateConverter::ngap_imsi_map_to_proto(ngap_imsi_map_, &imsi_proto);
  std::string proto_msg;
  redis_client->serialize(imsi_proto, proto_msg);
  std::size_t new_hash = std::hash<std::string>{}(proto_msg);
  if (new_hash != this->ngap_imsi_map_hash_) {
    redis_client->write_proto_str(NGAP_IMSI_MAP_TABLE_NAME, proto_msg, 0);
    this->ngap_imsi_map_hash_ = new_hash;
  }
}

}  // namespace magma5g
