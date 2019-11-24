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

#include <cpp_redis/cpp_redis>

#include "ServiceConfigLoader.h"

// TODO: Move this to redis wrapper
namespace {
constexpr char LOCALHOST[] = "127.0.0.1";
}

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

  // TODO: Move redis connection related functions to common wrapper
  /**
   * Reads and parses task state from db if persist_state is enabled
   * @return response code of operation
   */
  virtual int read_state_from_db()
  {
    if (persist_state_enabled) {
      auto db_read_fut = db_client->get(table_key);
      db_client->sync_commit();
      auto db_read_reply = db_read_fut.get();

      if (db_read_reply.is_null()) {
        OAILOG_ERROR(log_task, "Empty state in data store");
      } else if (db_read_reply.is_error() || !db_read_reply.is_string()) {
        OAILOG_ERROR(log_task, "Failed to read state from db");
      } else {
        ProtoType state_proto = ProtoType();
        if (!state_proto.ParseFromString(db_read_reply.as_string())) {
          OAILOG_ERROR(log_task, "Failed to parse state");
        }

        StateConverter::proto_to_state(state_proto, state_cache_p);
        OAILOG_DEBUG(log_task, "Finished reading state");
      }
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
      std::string serialized_state_s;
      ProtoType state_proto = ProtoType();
      StateConverter::state_to_proto(state_cache_p, &state_proto);

      if (!state_proto.SerializeToString(&serialized_state_s)) {
        OAILOG_ERROR(log_task, "Failed to serialize state protobuf");
        return;
      }

      auto db_write_fut = db_client->set(table_key, serialized_state_s);
      db_client->sync_commit();
      auto db_write_reply = db_write_fut.get();

      if (db_write_reply.is_error()) {
        OAILOG_ERROR(log_task, "Failed to write state to db");
        return;
      }

      OAILOG_DEBUG(log_task, "Finished writing state");
    }

    this->state_dirty = false;
  }

  /**
   * Initializes a connection to redis datastore.
   * @param addr is IP address of redis server
   * @return response code of success / error with db connection
   */
  int init_db_connection(const std::string& addr)
  {
    // Init db client service config
    magma::ServiceConfigLoader loader;

    auto config = loader.load_service_config("redis");
    auto port = config["port"].as<uint32_t>();

    db_client = std::make_unique<cpp_redis::client>();
    // Make connection to db
    db_client->connect(addr, port, nullptr);

    if (!db_client->is_connected()) {
      OAILOG_ERROR(log_task, "Failed to connect to redis");
      return RETURNerror;
    }

    OAILOG_INFO(
      log_task, "Connected to redis datastore on %s:%u\n", addr.c_str(), port);

    return RETURNok;
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
    log_task(LOG_MME_APP)
  {
  }
  virtual ~StateManager() = default;

  /**
   * Virtual function for allocating state_cache_p
   */
  virtual void create_state() = 0;

  // TODO: Make this a unique_ptr
  StateType* state_cache_p;
  std::unique_ptr<cpp_redis::client> db_client;
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
