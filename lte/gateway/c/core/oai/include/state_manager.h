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

#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include <assertions.h>
#include <common_defs.h>
#include <cstdlib>
#include <log.h>
#include <hashtable.h>

#ifdef __cplusplus
}
#endif

#include <unordered_map>
#include <conversions.h>
#include "redis_utils/redis_client.h"

namespace {
constexpr char IMSI_PREFIX[] = "IMSI";
}  // namespace

namespace magma {
namespace lte {

template<
    typename StateType, typename UeContextType, typename ProtoType,
    typename ProtoUe, typename StateConverter>
class StateManager {
 public:
  /**
   * @param read_from_db forces a read from db when true
   */
  virtual StateType* get_state(bool read_from_db) {
    AssertFatal(
        is_initialized,
        "StateManager init() function should be called to initialize state");

    // TODO: Add check for reentrant read/write function, to block multiple
    // reads

    state_dirty = true;

    AssertFatal(state_cache_p != nullptr, "State cache is NULL");

    if (persist_state_enabled && read_from_db) {
      free_state();
      create_state();
      read_state_from_db();
    }

    return state_cache_p;
  }

  virtual hash_table_ts_t* get_ue_state_ht() {
    AssertFatal(
        is_initialized,
        "StateManager init() function should be called to initialize state");

    return state_ue_ht;
  }

  /**
   * Reads and parses task state from db if persist_state is enabled
   * @return response code of operation
   */
  virtual status_code_e read_state_from_db() {
    if (persist_state_enabled) {
      ProtoType state_proto = ProtoType();
      if (redis_client->read_proto(table_key, state_proto) != RETURNok) {
        OAILOG_DEBUG(LOG_MME_APP, "Failed to read proto from db \n");
        return RETURNerror;
      }

      // Update the state version from redis
      this->task_state_version = redis_client->read_version(table_key);

      StateConverter::proto_to_state(state_proto, state_cache_p);
    }
    return RETURNok;
  }

  virtual status_code_e read_ue_state_from_db() {
    if (!persist_state_enabled) {
      return RETURNok;
    }
    auto keys = redis_client->get_keys("IMSI*" + task_name + "*");
    for (const auto& key : keys) {
      ProtoUe ue_proto = ProtoUe();
      auto* ue_context = (UeContextType*) (calloc(1, sizeof(UeContextType)));
      if (redis_client->read_proto(key.c_str(), ue_proto) != RETURNok) {
        return RETURNerror;
      }

      // Update each UE state version from redis
      this->ue_state_version[key] = redis_client->read_version(table_key);

      StateConverter::proto_to_ue(ue_proto, ue_context);

      hashtable_ts_insert(
          state_ue_ht, get_imsi_from_key(key), (void*) ue_context);
      OAILOG_DEBUG(log_task, "Reading UE state from db for %s", key.c_str());
    }
    return RETURNok;
  }

  /**
   * Writes task state to db if persist_state is enabled
   */
  virtual void write_state_to_db() {
    AssertFatal(
        is_initialized,
        "StateManager init() function should be called to initialize state");

    if (!state_dirty) {
      OAILOG_ERROR(log_task, "Tried to put state while it was not in use");
      return;
    }

    if (persist_state_enabled) {
      ProtoType state_proto = ProtoType();
      StateConverter::state_to_proto(state_cache_p, &state_proto);
      std::string proto_str;
      redis_client->serialize(state_proto, proto_str);
      std::size_t new_hash = std::hash<std::string>{}(proto_str);

      if (new_hash != this->task_state_hash) {
        if (redis_client->write_proto_str(
                table_key, proto_str, this->task_state_version) != RETURNok) {
          OAILOG_ERROR(log_task, "Failed to write state to db");
          return;
        }
        OAILOG_DEBUG(log_task, "Finished writing state");
        this->task_state_version++;
        this->state_dirty     = false;
        this->task_state_hash = new_hash;
      }
    }
  }

  virtual void write_ue_state_to_db(
      const UeContextType* ue_context, const std::string& imsi_str) {
    AssertFatal(
        is_initialized,
        "StateManager init() function should be called to initialize state");

    std::string proto_str;
    ProtoUe ue_proto = ProtoUe();
    StateConverter::ue_to_proto(ue_context, &ue_proto);
    redis_client->serialize(ue_proto, proto_str);
    std::size_t new_hash = std::hash<std::string>{}(proto_str);

    if (new_hash != this->ue_state_hash[imsi_str]) {
      std::string key = IMSI_PREFIX + imsi_str + ":" + task_name;
      if (redis_client->write_proto_str(
              key, proto_str, ue_state_version[imsi_str]) != RETURNok) {
        OAILOG_ERROR(
            log_task, "Failed to write UE state to db for IMSI %s",
            imsi_str.c_str());
        return;
      }

      this->ue_state_version[imsi_str]++;
      this->ue_state_hash[imsi_str] = new_hash;
      OAILOG_DEBUG(
          log_task, "Finished writing UE state for IMSI %s", imsi_str.c_str());
    }
  }

  std::string get_imsi_str(imsi64_t imsi64) {
    AssertFatal(
        is_initialized,
        "StateManager init() function should be called to initialize state");

    char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
    IMSI64_TO_STRING(imsi64, (char*) imsi_str, IMSI_BCD_DIGITS_MAX);
    return imsi_str;
  }

  void clear_ue_state_db(const std::string& imsi_str) {
    AssertFatal(
        is_initialized,
        "StateManager init() function should be called to initialize state");

    if (persist_state_enabled) {
      std::vector<std::string> keys = {IMSI_PREFIX + imsi_str + ":" +
                                       task_name};
      if (redis_client->clear_keys(keys) != RETURNok) {
        OAILOG_ERROR(log_task, "Failed to remove UE state from db");
        return;
      }
      OAILOG_DEBUG(log_task, "Removing UE state for IMSI %s", imsi_str.c_str());
    }
  }

  /**
   * Virtual function for freeing state_cache_p
   */
  virtual void free_state() = 0;

  bool is_persist_state_enabled() const { return persist_state_enabled; }

 protected:
  StateManager()
      : state_cache_p(nullptr),
        state_ue_ht(nullptr),
        redis_client(nullptr),
        is_initialized(false),
        state_dirty(false),
        persist_state_enabled(false),
        task_state_version(0),
        ue_state_version(0),
        task_state_hash(0),
        ue_state_hash(0),
        log_task(LOG_UTIL) {}
  virtual ~StateManager() = default;

  /**
   * Virtual function for allocating state_cache_p
   */
  virtual void create_state() = 0;

  imsi64_t get_imsi_from_key(const std::string& key) const {
    imsi64_t imsi64;
    std::string imsi_str_prefix = key.substr(0, key.find(':'));
    std::string imsi_str = imsi_str_prefix.substr(4, imsi_str_prefix.length());
    IMSI_STRING_TO_IMSI64(imsi_str.c_str(), &imsi64);
    return imsi64;
  }

  // TODO: Make this a unique_ptr
  StateType* state_cache_p;
  hash_table_ts_t* state_ue_ht;
  // TODO: Revisit one shared connection for all types of state
  std::unique_ptr<RedisClient> redis_client;
  // Flag for check asserting if the state has been initialized.
  bool is_initialized;
  // Flag for check asserting that write should be done after read.
  bool state_dirty;
  // Flag for enabling writing and reading to db.
  bool persist_state_enabled;
  // State version counters for task and ue context
  uint64_t task_state_version;
  std::unordered_map<std::string, uint64_t> ue_state_version;
  // Last written hash values for task and ue context
  std::size_t task_state_hash;
  std::unordered_map<std::string, std::size_t> ue_state_hash;

 protected:
  std::string table_key;
  std::string task_name;
  log_proto_t log_task;
};

}  // namespace lte
}  // namespace magma
