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
const char* MME_NAS_STATE_KEY = "mme_nas_state";
const int NUM_MAX_UE_HTBL_LISTS = 6;
const char* UE_ID_UE_CTXT_TABLE_NAME = "mme_app_mme_ue_s1ap_id_ue_context_htbl";
const char* IMSI_UE_ID_TABLE_NAME = "mme_app_imsi_ue_context_htbl";
const char* TUN_UE_ID_TABLE_NAME = "mme_app_tun11_ue_context_htbl";
const char* GUTI_UE_ID_TABLE_NAME = "mme_app_tun11_ue_context_htbl";
const char* ENB_UE_ID_MME_UE_ID_TABLE_NAME =
  "mme_app_enb_ue_s1ap_id_ue_context_htbl";
} // namespace

namespace magma {
namespace lte {

/**
 * Getter function for singleton instance of the MmeNasStateManager class,
 * guaranteed to be thread-safe and initialized only once
 */
MmeNasStateManager& MmeNasStateManager::getInstance()
{
  static MmeNasStateManager instance;
  return instance;
}

// Constructor for MME NAS state object
MmeNasStateManager::MmeNasStateManager():
  mme_nas_state_dirty_(false),
  max_ue_htbl_lists_(NUM_MAX_UE_HTBL_LISTS),
  mme_statistic_timer_(10)
{
}

// Destructor for MME NAS state object
MmeNasStateManager::~MmeNasStateManager()
{
  free_state();
}

int MmeNasStateManager::initialize_state(const mme_config_t* mme_config_p)
{
  persist_state_enabled = mme_config_p->use_stateless;
  max_ue_htbl_lists_ = mme_config_p->max_ues;
  mme_statistic_timer_ = mme_config_p->mme_statistic_timer;

  // Allocate the local mme state
  create_state();

  int rc = read_state_from_db();
  is_initialized = true;
  return rc;
}

void MmeNasStateManager::lock_mme_nas_state()
{
  AssertFatal(
    is_initialized, "Trying to lock state without initializing state manager");
  OAILOG_DEBUG(LOG_MME_APP, "Acquiring lock");
  pthread_rwlock_wrlock(&state_cache_p->rw_lock);
}

void MmeNasStateManager::unlock_mme_nas_state()
{
  AssertFatal(
    is_initialized,
    "Trying to unlock state without initializing state manager");
  OAILOG_DEBUG(LOG_MME_APP, "Releasing lock");
  pthread_rwlock_unlock(&state_cache_p->rw_lock);
  OAILOG_DEBUG(LOG_MME_APP, "Lock released");
}

void MmeNasStateManager::write_state_to_db(mme_app_desc_t** task_state_ptr)
{
  AssertFatal(
    is_initialized, "Calling write without initializing MME state manager");
  if (!mme_nas_state_dirty_) {
    OAILOG_ERROR(
      LOG_MME_APP, "Tried to put mme_nas_state without getting it first");
    return;
  }

  // check if the calling thread owns the lock on state
  AssertFatal(
    pthread_rwlock_trywrlock(&state_cache_p->rw_lock),
    "Thread trying to write state without locking");

  // clear up the local ptr of the task holding the state pointer
  *task_state_ptr = nullptr;

  if (persist_state_enabled) {
    std::string serialized_state;
    // convert the in-memory state to proto message
    MmeNasState state_proto = MmeNasState();
    MmeNasStateConverter::state_to_proto(state_cache_p, &state_proto);

    if (!state_proto.SerializeToString(&serialized_state)) {
      OAILOG_ERROR(LOG_MME_APP, "Failed to serialize MME state");
      goto error;
      return;
    }

    OAILOG_DEBUG(LOG_MME_APP, "Writing serialized MME state to redis");
    // write the proto to redis store
    if (
      redis_client->write(MME_NAS_STATE_KEY, serialized_state) !=
      RETURNok) {
      OAILOG_ERROR(LOG_MME_APP, "Failed to write state to db");
      return;
    }
    OAILOG_DEBUG(LOG_MME_APP, "MME NAS state written to redis");
  }
  mme_nas_state_dirty_ = false;
error:
  unlock_mme_nas_state();
}

/**
 * Getter function to lock the state before returning the pointer to in-memory
 * user state. The read_from_db flag is a debug flag to force read from the
 * data store instead of just returning the pointer.
 */
mme_app_desc_t* MmeNasStateManager::get_locked_mme_nas_state(bool read_from_db)
{
  AssertFatal(
    is_initialized,
    "Calling get_locked_mme_nas_state without initializing state manager");
  OAILOG_DEBUG(
    LOG_MME_APP,
    "Inside get_locked_mme_nas_state with read_from_db %d",
    read_from_db);
  lock_mme_nas_state();
  return get_state(read_from_db);
}

/**
 * Getter function to get the pointer to the in-memory user state. The
 * read_from_db flag is a debug flag to force read from data store instead of
 * just returning the pointer. In non-debug mode, the state is only read from
 * data store when initialize_state is called and get_state just returns the
 * pointer to that structure.
 */
mme_app_desc_t* MmeNasStateManager::get_state(bool read_from_db)
{
  AssertFatal(
    is_initialized,
    "Calling get_state without initializing state manager");
  AssertFatal(state_cache_p, "mme_nas_state is NULL");
  OAILOG_DEBUG(
    LOG_MME_APP, "Inside get_state with read_from_db %d", read_from_db);

  // check if the calling thread owns the lock on state
  AssertFatal(
    pthread_rwlock_trywrlock(&state_cache_p->rw_lock),
    "Thread trying to get state without locking");

  mme_nas_state_dirty_ = true;
  if (persist_state_enabled && read_from_db) {
    // free up the memory allocated to hashtables
    OAILOG_DEBUG(LOG_MME_APP, "Freeing up in-memory hashtables");
    clear_mme_nas_hashtables();
    // allocate memory for hashtables
    OAILOG_DEBUG(LOG_MME_APP, "Allocating memory for new hashtables");
    create_hashtables();
    // read the state from data store
    int rc = read_state_from_db();
    AssertFatal(state_cache_p, "mme_nas_state is NULL");
  }
  return state_cache_p;
}

// This is a helper function for debugging. If the state manager needs to clear
// the state in the data store, it can call this function to delete the key.
void MmeNasStateManager::clear_db_state()
{
  OAILOG_DEBUG(LOG_MME_APP, "Clearing state in data store");
  std::vector<std::string> keys_to_del;
  keys_to_del.emplace_back(MME_NAS_STATE_KEY);

  if (redis_client->clear_keys(keys_to_del) != RETURNok) {
    OAILOG_ERROR(LOG_MME_APP, "Failed to clear the state in data store");
    return;
  }
}

// Initialize state that is non-persistent, e.g. mutex locks and timers
void MmeNasStateManager::mme_nas_state_init_local_state()
{
  // create local lock for this state
  pthread_rwlock_init(&state_cache_p->rw_lock, nullptr);

  // create statistic timer locally
  state_cache_p->statistic_timer_period = mme_statistic_timer_;

  // Request for periodic timer to print statistics in debug mode
  if (
    timer_setup(
      mme_statistic_timer_,
      0,
      TASK_MME_APP,
      INSTANCE_DEFAULT,
      TIMER_PERIODIC,
      nullptr,
      0,
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
void MmeNasStateManager::create_hashtables()
{
  bstring b = bfromcstr(IMSI_UE_ID_TABLE_NAME);
  state_cache_p->mme_ue_contexts.imsi_ue_context_htbl =
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
  state_cache_p->mme_ue_contexts.mme_ue_s1ap_id_ue_context_htbl =
    hashtable_ts_create(
      max_ue_htbl_lists_, nullptr, mme_app_state_free_ue_context, b);
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
void MmeNasStateManager::create_state()
{
  state_cache_p = (mme_app_desc_t*) calloc(1, sizeof(mme_app_desc_t));
  if (!state_cache_p) {
    return;
  }

  mme_imsi_map_ = (mme_imsi_map_t*) calloc(1, sizeof(mme_imsi_map_t));

  create_hashtables();
  // Initialize the lock and local timers, which are non-persistent
  mme_nas_state_init_local_state();

  create_mme_imsi_map();
}

// Delete the hashtables for MME NAS state
void MmeNasStateManager::clear_mme_nas_hashtables()
{
  if (!state_cache_p) {
    return;
  }

  hashtable_uint64_ts_destroy(
    state_cache_p->mme_ue_contexts.imsi_ue_context_htbl);
  hashtable_uint64_ts_destroy(
    state_cache_p->mme_ue_contexts.tun11_ue_context_htbl);
  hashtable_ts_destroy(
    state_cache_p->mme_ue_contexts.mme_ue_s1ap_id_ue_context_htbl);
  hashtable_uint64_ts_destroy(
    state_cache_p->mme_ue_contexts.enb_ue_s1ap_id_ue_context_htbl);
  obj_hashtable_uint64_ts_destroy(
    state_cache_p->mme_ue_contexts.guti_ue_context_htbl);
}

// Free the memory allocated to state pointer and destroy the read/write lock
void MmeNasStateManager::free_state()
{
  if (!state_cache_p) {
    return;
  }
  lock_mme_nas_state();
  clear_mme_nas_hashtables();
  timer_remove(state_cache_p->statistic_timer_id, nullptr);
  pthread_rwlock_destroy(&state_cache_p->rw_lock);
  free(state_cache_p);
  state_cache_p = nullptr;

  clear_mme_imsi_map();
}

void MmeNasStateManager::create_mme_imsi_map()
{
  mme_imsi_map_ = (mme_imsi_map_t*) calloc(1, sizeof(mme_imsi_map_t));

  mme_imsi_map_->mme_ue_s1ap_imsi_htbl =
    hashtable_uint64_ts_create(max_ue_htbl_lists_, nullptr, nullptr);
  mme_imsi_map_->enb_s1ap_mme_ue_id_htbl =
    hashtable_uint64_ts_create(max_ue_htbl_lists_, nullptr, nullptr);

  MmeImsiMap imsi_proto = MmeImsiMap();
  redis_client->read_proto(MME_IMSI_MAP_TABLE_NAME, imsi_proto);

  MmeNasStateConverter::proto_to_mme_imsi_map(imsi_proto, mme_imsi_map_);
}

mme_imsi_map_t* MmeNasStateManager::get_mme_imsi_map()
{
  return mme_imsi_map_;
}

void MmeNasStateManager::put_mme_imsi_map() {
  MmeImsiMap imsi_proto = MmeImsiMap();
  MmeNasStateConverter::mme_imsi_map_to_proto(mme_imsi_map_, &imsi_proto);
  redis_client->write_proto(MME_IMSI_MAP_TABLE_NAME, imsi_proto);
}

void MmeNasStateManager::clear_mme_imsi_map() {
  if(!mme_imsi_map_) {
    return;
  }

  hashtable_uint64_ts_destroy(mme_imsi_map_->mme_ue_s1ap_imsi_htbl);
  hashtable_uint64_ts_destroy(mme_imsi_map_->enb_s1ap_mme_ue_id_htbl);

  free_wrapper((void **) &mme_imsi_map_);

}

} // namespace lte
} // namespace magma
