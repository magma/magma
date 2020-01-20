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

#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include <assertions.h>
#include <common_defs.h>
#include <cstdlib>
#include <log.h>

#ifdef __cplusplus
}
#endif

#include "redis_utils/redis_client.h"

namespace magma {
namespace lte {

template<typename StateType, typename ProtoType, typename StateConverter>
class StateManager {
 public:
  /**
   * @param read_from_db forces a read from db when true
   */
  virtual StateType* get_state(bool read_from_db)
  {
    AssertFatal(
      is_initialized,
      "StateManager init() function should be called to initialize state");

    //TODO: Add check for reentrant read/write function, to block multiple reads

    state_dirty = true;

    AssertFatal(state_cache_p != nullptr, "State cache is NULL");

    if (persist_state_enabled && read_from_db) {
      free_state();
      create_state();
      read_state_from_db();
    }

    return state_cache_p;
  }

  // TODO: Read state per IMSI
  /**
   * Reads and parses task state from db if persist_state is enabled
   * @return response code of operation
   */
  virtual int read_state_from_db()
  {
    if (persist_state_enabled) {
      ProtoType state_proto = ProtoType();
      if (redis_client->read_proto(table_key, state_proto) != RETURNok) {
        return RETURNerror;
      }
      StateConverter::proto_to_state(state_proto, state_cache_p);
    }
    return RETURNok;
  }

  /**
   * Writes task state to db if persist_state is enabled
   */
  virtual void write_state_to_db()
  {
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

      if (redis_client->write_proto(table_key, state_proto) != RETURNok) {
        OAILOG_ERROR(log_task, "Failed to write state to db");
        return;
      }
      OAILOG_DEBUG(log_task, "Finished writing state");
    }

    this->state_dirty = false;
  }

  /**
   * Virtual function for freeing state_cache_p
   */
  virtual void free_state() = 0;

 protected:
  StateManager():
    is_initialized(false),
    state_dirty(false),
    persist_state_enabled(false),
    state_cache_p(nullptr),
    log_task(LOG_MME_APP),
    redis_client(std::make_unique<RedisClient>())
  {
  }
  virtual ~StateManager() = default;

  /**
   * Virtual function for allocating state_cache_p
   */
  virtual void create_state() = 0;

  // TODO: Make this a unique_ptr
  StateType* state_cache_p;
  // TODO: Revisit one shared connection for all types of state
  std::unique_ptr<RedisClient> redis_client;
  // Flag for check asserting if the state has been initialized.
  bool is_initialized;
  // Flag for check asserting that write should be done after read.
  // TODO: Convert this to state versioning variable
  bool state_dirty;
  // Flag for enabling writing and reading to db.
  bool persist_state_enabled;
  std::string table_key;
  log_proto_t log_task;
};

} // namespace lte
} // namespace magma
