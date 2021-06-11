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

extern "C" {
#include "assertions.h"
#include "dynamic_memory_check.h"
#include "emm_proc.h"
#include "log.h"
#include "timer.h"
}

#include "mme_app_state_manager.h"

namespace {
constexpr char MME_NAS_STATE_KEY[] = "mme_nas_state";
const int NUM_MAX_UE_HTBL_LISTS    = 6;
constexpr char UE_ID_UE_CTXT_TABLE_NAME[] =
    "mme_app_mme_ue_s1ap_id_ue_context_htbl";
constexpr char IMSI_UE_ID_TABLE_NAME[] = "mme_app_imsi_ue_context_htbl";
constexpr char TUN_UE_ID_TABLE_NAME[]  = "mme_app_tun11_ue_context_htbl";
constexpr char GUTI_UE_ID_TABLE_NAME[] = "mme_app_tun11_ue_context_htbl";
constexpr char ENB_UE_ID_MME_UE_ID_TABLE_NAME[] =
    "mme_app_enb_ue_s1ap_id_ue_context_htbl";
constexpr char MME_TASK_NAME[] = "MME";
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

// Constructor for MME NAS state object
MmeNasStateManager::MmeNasStateManager()
    : max_ue_htbl_lists_(NUM_MAX_UE_HTBL_LISTS), mme_statistic_timer_(10) {}

// Destructor for MME NAS state object
MmeNasStateManager::~MmeNasStateManager() {
  free_state();
}

int MmeNasStateManager::initialize_state(const mme_config_t* mme_config_p) {
  persist_state_enabled = mme_config_p->use_stateless;
  max_ue_htbl_lists_    = mme_config_p->max_ues;
  mme_statistic_timer_  = mme_config_p->mme_statistic_timer;
  log_task              = LOG_MME_APP;
  task_name             = MME_TASK_NAME;
  table_key             = MME_NAS_STATE_KEY;

  // Allocate the local mme state
  create_state();

  redis_client = std::make_unique<RedisClient>(persist_state_enabled);
  int rc       = read_state_from_db();
  read_ue_state_from_db();
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
  AssertFatal(
      is_initialized, "Calling get_state without initializing state manager");
  AssertFatal(state_cache_p, "mme_nas_state is NULL");
  OAILOG_DEBUG(
      LOG_MME_APP, "Inside get_state with read_from_db %d", read_from_db);

  state_dirty = true;
  if (persist_state_enabled && read_from_db) {
    // free up the memory allocated to hashtables
    OAILOG_DEBUG(LOG_MME_APP, "Freeing up in-memory hashtables");
    clear_mme_nas_hashtables();
    // allocate memory for hashtables
    OAILOG_DEBUG(LOG_MME_APP, "Allocating memory for new hashtables");
    create_hashtables();
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
void MmeNasStateManager::mme_nas_state_init_local_state() {
  // create statistic timer locally
  state_cache_p->statistic_timer_period = mme_statistic_timer_;

  // Request for periodic timer to print statistics in debug mode
  if (timer_setup(
          mme_statistic_timer_, 0, TASK_MME_APP, INSTANCE_DEFAULT,
          TIMER_PERIODIC, nullptr, 0,
          &(state_cache_p->statistic_timer_id)) < 0) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Failed to request new timer for statistics with %ds "
        "of periocidity\n",
        mme_statistic_timer_);
    state_cache_p->statistic_timer_id = 0;
  }
}

// Create the hashtables for MME NAS state
void MmeNasStateManager::create_hashtables() {
  bstring b = bfromcstr(IMSI_UE_ID_TABLE_NAME);
  state_cache_p->mme_ue_contexts.imsi_mme_ue_id_htbl =
      hashtable_uint64_ts_create(max_ue_htbl_lists_, nullptr, b);
  btrunc(b, 0);
  bassigncstr(b, TUN_UE_ID_TABLE_NAME);
  state_cache_p->mme_ue_contexts.tun11_ue_context_htbl =
      hashtable_uint64_ts_create(max_ue_htbl_lists_, nullptr, b);
  AssertFatal(
      sizeof(uintptr_t) >= sizeof(uint64_t),
      "Problem with mme_ue_s1ap_id_ue_context_htbl in MME_APP");
  btrunc(b, 0);
  bassigncstr(b, UE_ID_UE_CTXT_TABLE_NAME);
  state_ue_ht = hashtable_ts_create(
      max_ue_htbl_lists_, nullptr, mme_app_state_free_ue_context, b);

  if (!(state_ue_ht->lock_attr = (pthread_mutexattr_t*) calloc(
            max_ue_htbl_lists_, sizeof(pthread_mutexattr_t)))) {
    free_wrapper((void**) &state_ue_ht->lock_nodes);
    free_wrapper((void**) &state_ue_ht->nodes);
    free_wrapper((void**) &state_ue_ht->name);
    free_wrapper((void**) &state_ue_ht);
    return;
  }

  for (int i = 0; i < max_ue_htbl_lists_; i++) {
    pthread_mutexattr_init(&state_ue_ht->lock_attr[i]);
    pthread_mutexattr_settype(
        &state_ue_ht->lock_attr[i], PTHREAD_MUTEX_RECURSIVE);
    pthread_mutex_init(&state_ue_ht->lock_nodes[i], &state_ue_ht->lock_attr[i]);
  }

  btrunc(b, 0);
  bassigncstr(b, ENB_UE_ID_MME_UE_ID_TABLE_NAME);
  state_cache_p->mme_ue_contexts.enb_ue_s1ap_id_ue_context_htbl =
      hashtable_uint64_ts_create(max_ue_htbl_lists_, nullptr, b);
  btrunc(b, 0);
  bassigncstr(b, GUTI_UE_ID_TABLE_NAME);
  state_cache_p->mme_ue_contexts.guti_ue_context_htbl =
      obj_hashtable_uint64_ts_create(max_ue_htbl_lists_, nullptr, nullptr, b);
  bdestroy_wrapper(&b);
}

// Initialize memory for MME state before reading from data-store
void MmeNasStateManager::create_state() {
  state_cache_p = (mme_app_desc_t*) calloc(1, sizeof(mme_app_desc_t));
  if (!state_cache_p) {
    return;
  }
  state_cache_p->mme_app_ue_s1ap_id_generator = 1;

  create_hashtables();
  // Initialize the local timers, which are non-persistent
  mme_nas_state_init_local_state();
}

// Delete the hashtables for MME NAS state
void MmeNasStateManager::clear_mme_nas_hashtables() {
  if (!state_cache_p) {
    return;
  }

  hashtable_ts_destroy(state_ue_ht);
  hashtable_uint64_ts_destroy(
      state_cache_p->mme_ue_contexts.imsi_mme_ue_id_htbl);
  hashtable_uint64_ts_destroy(
      state_cache_p->mme_ue_contexts.tun11_ue_context_htbl);
  hashtable_uint64_ts_destroy(
      state_cache_p->mme_ue_contexts.enb_ue_s1ap_id_ue_context_htbl);
  obj_hashtable_uint64_ts_destroy(
      state_cache_p->mme_ue_contexts.guti_ue_context_htbl);
}

// Free the memory allocated to state pointer
void MmeNasStateManager::free_state() {
  if (!state_cache_p) {
    return;
  }
  clear_mme_nas_hashtables();
  timer_remove(state_cache_p->statistic_timer_id, nullptr);
  free(state_cache_p);
  state_cache_p = nullptr;
}

status_code_e MmeNasStateManager::read_ue_state_from_db() {
  if (persist_state_enabled) {
    auto keys = redis_client->get_keys("IMSI*" + task_name + "*");
    for (const auto& key : keys) {
      OAILOG_DEBUG(log_task, "Reading UE state from db for %s", key.c_str());
      oai::UeContext ue_proto = oai::UeContext();
      auto* ue_context =
          (ue_mm_context_t*) (calloc(1, sizeof(ue_mm_context_t)));
      if (redis_client->read_proto(key, ue_proto) != RETURNok) {
        return RETURNerror;
      }
      MmeNasStateConverter::proto_to_ue(ue_proto, ue_context);

      hashtable_rc_t h_rc = hashtable_ts_insert(
          state_ue_ht, ue_context->mme_ue_s1ap_id, (void*) ue_context);
      if (HASH_TABLE_OK != h_rc) {
        OAILOG_ERROR(
            log_task,
            "Failed to insert UE state with key mme_ue_s1ap_id "
            " " MME_UE_S1AP_ID_FMT " (Error Code: %s)\n",
            ue_context->mme_ue_s1ap_id, hashtable_rc_code2string(h_rc));
      } else {
        OAILOG_DEBUG(
            log_task,
            "Inserted UE state with key mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT,
            ue_context->mme_ue_s1ap_id);
      }
    }
  }
  return RETURNok;
}

}  // namespace lte
}  // namespace magma
