/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "SessionState.h"
#include "RedisStoreClient.h"
#include "magma_logging.h"

namespace magma {
namespace lte {

RedisStoreClient::RedisStoreClient(
    std::shared_ptr<cpp_redis::client> client, const std::string& redis_table,
    std::shared_ptr<StaticRuleStore> rule_store)
    : client_(client), redis_table_(redis_table), rule_store_(rule_store) {}

bool RedisStoreClient::try_redis_connect() {
  ServiceConfigLoader loader;
  auto config = loader.load_service_config("redis");
  auto port   = config["port"].as<uint32_t>();
  auto addr   = config["bind"].as<std::string>();
  try {
    client_->connect(
        addr, port,
        [](const std::string& host, std::size_t port,
           cpp_redis::client::connect_state status) {
          if (status == cpp_redis::client::connect_state::dropped) {
            MLOG(MERROR) << "Client disconnected from " << host << ":" << port;
          }
        });
    return client_->is_connected();
  } catch (const cpp_redis::redis_error& e) {
    MLOG(MERROR) << "Could not connect to redis: " << e.what();
    return false;
  }
}

SessionMap RedisStoreClient::read_sessions(
    std::set<std::string> subscriber_ids) {
  // The approach here is made assuming that the SessionStore only has one
  // call being processed at a time, and that the writes it makes are done
  // atomically. Based on that, reads can be done without using Redis
  // transactions, or EVAL.
  std::unordered_map<std::string, std::future<cpp_redis::reply>> futures;
  for (const std::string& key : subscriber_ids) {
    futures[key] = client_->hget(redis_table_, key);
  }

  client_->sync_commit();

  SessionMap session_map;
  for (const std::string& key : subscriber_ids) {
    auto reply = futures[key].get();
    if (reply.is_null()) {
      // value just doesn't exist
      session_map[key] = std::vector<std::unique_ptr<SessionState>>{};
    } else if (reply.is_error()) {
      MLOG(MERROR) << "RedisStoreClient: Unable to get value for key " << key;
      throw RedisReadFailed();
    } else if (!reply.is_string()) {
      session_map[key] = std::vector<std::unique_ptr<SessionState>>{};
    } else {
      session_map[key] = std::move(deserialize_session_vec(reply.as_string()));
    }
  }
  return session_map;
}

SessionMap RedisStoreClient::read_all_sessions() {
  SessionMap session_map;
  auto hgetall_future = client_->hgetall(redis_table_);
  client_->sync_commit();

  auto reply = hgetall_future.get();
  if (reply.is_error()) {
    MLOG(MERROR) << "unable to read all sessions from redis";
    return session_map;
  }
  auto array = reply.as_array();
  for (int i = 0; i < array.size(); i += 2) {
    auto key_reply = array[i];
    if (!key_reply.is_string()) {
      MLOG(MERROR) << "Non string key found in sessions from redis";
      continue;
    }
    auto key         = key_reply.as_string();
    auto value_reply = array[i + 1];
    if (!value_reply.is_string()) {
      MLOG(MERROR) << "RedisStoreClient: Unable to get value for key " << key;
      session_map[key] = std::vector<std::unique_ptr<SessionState>>{};
    } else {
      session_map[key] =
          std::move(deserialize_session_vec(value_reply.as_string()));
    }
  }
  return session_map;
}

bool RedisStoreClient::write_sessions(SessionMap session_map) {
  // Writes should happen via a transaction, otherwise the state inside in
  // Redis may not be recoverable or consistent.
  // For reference, see https://redis.io/topics/transactions

  // First we need to watch the keys that we intend to write to.
  // If we don't, then one HSET might succeed but another will fail.
  std::vector<std::string> keys;
  for (auto& it : session_map) {
    keys.push_back(it.first);
  }
  client_->watch(keys);

  // Now set the MULTI command.
  // Subsequent commands end up being queued for atomic execution with EXEC.
  // Together with WATCH, if one of the keys we intend to set are modified, then
  // the entire EXEC does not execute.
  client_->multi();

  // And now we queue up HSET commands after we've set up some sort of safety
  // guarantees.
  for (auto& it : session_map) {
    client_->hset(redis_table_, it.first, serialize_session_vec(it.second));
  }

  auto exec_future = client_->exec();
  client_->sync_commit();

  auto reply = exec_future.get();
  if (!reply.ok()) {
    MLOG(MERROR) << "Failed to write sessions to Redis.";
    return false;
  }
}

std::string RedisStoreClient::serialize_session_vec(
    std::vector<std::unique_ptr<SessionState>>& session_vec) {
  folly::dynamic marshaled = folly::dynamic::array;
  for (auto& session_ptr : session_vec) {
    auto stored_session = session_ptr->marshal();
    marshaled.push_back(serialize_stored_session(stored_session));
  }
  std::string serialized = folly::toJson(marshaled);
  return serialized;
}

std::vector<std::unique_ptr<SessionState>>
RedisStoreClient::deserialize_session_vec(std::string serialized) {
  auto folly_serialized    = folly::StringPiece(serialized);
  folly::dynamic marshaled = folly::parseJson(folly_serialized);
  std::vector<std::unique_ptr<SessionState>> session_vec;
  for (auto& it : marshaled) {
    auto stored_session = deserialize_stored_session(it.getString());
    session_vec.push_back(
        SessionState::unmarshal(stored_session, *rule_store_));
  }
  return session_vec;
}

}  // namespace lte
}  // namespace magma
