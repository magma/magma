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

#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_state_manager.hpp"

#include <string>

extern "C" {
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/oai/common/log.h"
}

#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/emm_proc.hpp"

namespace {
constexpr char MME_NAS_STATE_KEY[] = "mme_nas_state";
const int NUM_MAX_UE_HTBL_LISTS = 6;
constexpr char MME_UE_ID2UE_CTXT_MAP_NAME[] =
    "mme_app_mme_ue_s1ap_id2ue_context_map";
constexpr char MME_IMSI2MME_UE_ID_MAP_NAME[] = "mme_imsi2ue_id_map";
constexpr char MME_S11_TEID2MME_UE_ID_MAP_NAME[] = "mme_s11_teid2ue_id_map";
constexpr char GUTI_UE_ID_TABLE_NAME[] = "mme_app_tun11_ue_context_htbl";
constexpr char MME_ENB_UE_S1AP_KEY2MME_UE_ID_MAP_NAME[] =
    "mme_enb_ue_s1ap_key2ue_id_map";
constexpr char MME_TASK_NAME[] = "MME";
constexpr char MME_UEIP_IMSI_MAP_NAME[] = "mme_ueip_imsi_map";
}  // namespace

namespace magma {
namespace lte {

/**
 * Getter function for singleton instance of the MmeNasStateManager class,
 * guaranteed to be thread-safe and initialized only once
 */
MmeNasStateManager& MmeNasStateManager::getInstance() {
  static MmeNasStateManager instance;
  return instance;
}

proto_map_uint32_ue_context_t* MmeNasStateManager::get_ue_state_map() {
  return &mme_app_state_ue_map;
}

// Constructor for MME NAS state object
MmeNasStateManager::MmeNasStateManager()
    : max_ue_htbl_lists_(NUM_MAX_UE_HTBL_LISTS), ueip_imsi_map{0} {}

// Destructor for MME NAS state object
MmeNasStateManager::~MmeNasStateManager() { free_state(); }

int MmeNasStateManager::initialize_state(const mme_config_t* mme_config_p) {
  persist_state_enabled = mme_config_p->use_stateless;
  max_ue_htbl_lists_ = mme_config_p->max_ues;
  log_task = LOG_MME_APP;
  task_name = MME_TASK_NAME;
  table_key = MME_NAS_STATE_KEY;

  // Allocate the local mme state
  create_state();
#if !MME_UNIT_TEST
  OAILOG_DEBUG(LOG_MME_APP, "MME_UNIT_TEST Flag is Disabled");
  redis_client = std::make_unique<RedisClient>(persist_state_enabled);
#else
  redis_client = std::make_unique<RedisClient>(false);
#endif
  int rc = read_state_from_db();
  read_ue_state_from_db();
  create_mme_ueip_imsi_map();
  is_initialized = true;
  return rc;
}

/**
 * Getter function to get the pointer to the in-memory user state. The
 * read_from_db flag is a debug flag to force read from data store instead of
 * just returning the pointer. In non-debug mode, the state is only read from
 * data store when initialize_state is called and get_state just returns the
 * pointer to that structure.
 */
mme_app_desc_t* MmeNasStateManager::get_state(bool read_from_db) {
  AssertFatal(is_initialized,
              "Calling get_state without initializing state manager");
  AssertFatal(state_cache_p, "mme_nas_state is NULL");
  OAILOG_DEBUG(LOG_MME_APP, "Inside get_state with read_from_db %d",
               read_from_db);

  state_dirty = true;
  if (persist_state_enabled && read_from_db) {
    // free up the memory allocated to protobuf maps
    OAILOG_DEBUG(LOG_MME_APP, "Freeing up in-memory protobuf maps");
    clear_mme_nas_protomaps();
    // allocate memory for protobuf maps
    OAILOG_DEBUG(LOG_MME_APP, "Allocating memory for new protobuf maps");
    create_protomaps();
    // read the state from data store
    int rc = read_state_from_db();
    if (rc != RETURNok) {
      OAILOG_ERROR(LOG_MME_APP, "Failed to read task state from redis");
    }
    read_ue_state_from_db();
    AssertFatal(state_cache_p, "mme_nas_state is NULL");
  }
  return state_cache_p;
}

// This is a helper function for debugging. If the state manager needs to clear
// the state in the data store, it can call this function to delete the key.
void MmeNasStateManager::clear_db_state() {
  OAILOG_DEBUG(LOG_MME_APP, "Clearing state in data store");
  std::vector<std::string> keys_to_del;
  keys_to_del.emplace_back(MME_NAS_STATE_KEY);

  if (redis_client->clear_keys(keys_to_del) != RETURNok) {
    OAILOG_ERROR(LOG_MME_APP, "Failed to clear the state in data store");
    return;
  }
}

// Initialize state that is non-persistent, e.g. timers
void MmeNasStateManager::mme_nas_state_init_local_state() {}

// Create the protobuf maps for MME NAS state
void MmeNasStateManager::create_protomaps() {
  state_cache_p->mme_ue_contexts.imsi2mme_ueid_map.map =
      new google::protobuf::Map<uint64_t, uint32_t>();
  state_cache_p->mme_ue_contexts.imsi2mme_ueid_map.set_name(
      MME_IMSI2MME_UE_ID_MAP_NAME);
  state_cache_p->mme_ue_contexts.s11_teid2mme_ueid_map.map =
      new google::protobuf::Map<uint32_t, uint32_t>();
  state_cache_p->mme_ue_contexts.s11_teid2mme_ueid_map.set_name(
      MME_S11_TEID2MME_UE_ID_MAP_NAME);

  mme_app_state_ue_map.map =
      new google::protobuf::Map<uint32_t, struct ue_mm_context_s*>();
  mme_app_state_ue_map.set_name(MME_UE_ID2UE_CTXT_MAP_NAME);
  mme_app_state_ue_map.bind_callback(mme_app_state_free_ue_context);

  state_cache_p->mme_ue_contexts.enb_ue_s1ap_key2mme_ueid_map.map =
      new google::protobuf::Map<uint64_t, uint32_t>();
  state_cache_p->mme_ue_contexts.enb_ue_s1ap_key2mme_ueid_map.set_name(
      MME_ENB_UE_S1AP_KEY2MME_UE_ID_MAP_NAME);

  bstring b = bfromcstr(GUTI_UE_ID_TABLE_NAME);
  state_cache_p->mme_ue_contexts.guti_ue_context_htbl =
      obj_hashtable_uint64_ts_create(max_ue_htbl_lists_, nullptr, nullptr, b);
  bdestroy_wrapper(&b);
}

// Initialize memory for MME state before reading from data-store
void MmeNasStateManager::create_state() {
  state_cache_p = (mme_app_desc_t*)calloc(1, sizeof(mme_app_desc_t));
  if (!state_cache_p) {
    return;
  }
  state_cache_p->mme_app_ue_s1ap_id_generator = 1;

  create_protomaps();
  // Initialize the local timers, which are non-persistent
  mme_nas_state_init_local_state();
}

// Delete the protobuf maps for MME NAS state
void MmeNasStateManager::clear_mme_nas_protomaps() {
  if (!state_cache_p) {
    return;
  }

  mme_app_state_ue_map.destroy_map();
  state_cache_p->mme_ue_contexts.imsi2mme_ueid_map.destroy_map();
  state_cache_p->mme_ue_contexts.s11_teid2mme_ueid_map.destroy_map();
  state_cache_p->mme_ue_contexts.enb_ue_s1ap_key2mme_ueid_map.destroy_map();
  obj_hashtable_uint64_ts_destroy(
      state_cache_p->mme_ue_contexts.guti_ue_context_htbl);
}

// Free the memory allocated to state pointer
void MmeNasStateManager::free_state() {
  if (!state_cache_p) {
    return;
  }
  clear_mme_nas_protomaps();
  free(state_cache_p);
  state_cache_p = nullptr;
}

status_code_e MmeNasStateManager::read_ue_state_from_db() {
#if !MME_UNIT_TEST
  if (persist_state_enabled) {
    auto keys = redis_client->get_keys("IMSI*" + task_name + "*");
    for (const auto& key : keys) {
      OAILOG_DEBUG(log_task, "Reading UE state from db for %s", key.c_str());
      oai::UeContext ue_proto = oai::UeContext();
      if (redis_client->read_proto(key, ue_proto) != RETURNok) {
        return RETURNerror;
      }
      auto* ue_context = reinterpret_cast<ue_mm_context_t*>(
          calloc(1, sizeof(ue_mm_context_t)));
      MmeNasStateConverter::proto_to_ue(ue_proto, ue_context);

      if (mme_app_state_ue_map.insert(ue_context->mme_ue_s1ap_id, ue_context) !=
          magma::PROTO_MAP_OK) {
        OAILOG_ERROR(log_task,
                     "Failed to insert UE state with key mme_ue_s1ap_id "
                     " " MME_UE_S1AP_ID_FMT,
                     ue_context->mme_ue_s1ap_id);
      } else {
        OAILOG_DEBUG(
            log_task,
            "Inserted UE state with key mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT,
            ue_context->mme_ue_s1ap_id);
      }
    }
  }
#endif
  return RETURNok;
}

void MmeNasStateManager::create_mme_ueip_imsi_map() {
#if !MME_UNIT_TEST
  if (!persist_state_enabled) {
    OAILOG_ERROR(log_task, "persist_state_enabled is not enabled \n");
    return;
  }
  oai::MmeUeIpImsiMap ueip_proto = oai::MmeUeIpImsiMap();
  redis_client->read_proto(MME_UEIP_IMSI_MAP_NAME, ueip_proto);

  MmeNasStateConverter::mme_app_proto_to_ueip_imsi_map(ueip_proto,
                                                       ueip_imsi_map);
#endif
  return;
}

void MmeNasStateManager::write_mme_ueip_imsi_map_to_db() {
  if (!persist_state_enabled) {
    OAILOG_ERROR(log_task, "persist_state_enabled is not enabled \n");
    return;
  }

  oai::MmeUeIpImsiMap ueip_proto = oai::MmeUeIpImsiMap();
  MmeNasStateConverter::mme_app_ueip_imsi_map_to_proto(ueip_imsi_map,
                                                       &ueip_proto);
  std::string proto_msg;
  redis_client->serialize(ueip_proto, proto_msg);

  // ueip_imsi_map is not state service synced, so version will not be updated
  redis_client->write_proto_str(MME_UEIP_IMSI_MAP_NAME, proto_msg, 0);
  return;
}

UeIpImsiMap& MmeNasStateManager::get_mme_ueip_imsi_map(void) {
  return ueip_imsi_map;
}

}  // namespace lte
}  // namespace magma
