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

#include "RedisStoreClient.h"

#include <cpp_redis/core/client.hpp>  // for client, client::connect_state
#include <cpp_redis/core/reply.hpp>   // for reply
#include <cpp_redis/misc/error.hpp>   // for redis_error
#include <folly/Range.h>              // for operator<<, StringPiece
#include <folly/dynamic.h>            // for dynamic
#include <folly/json.h>               // for parseJson, toJson
#include <glog/logging.h>             // for COMPACT_GOOGLE_LOG_INFO, LogMes...
#include <stddef.h>                   // for size_t
#include <stdint.h>                   // for uint32_t
#include <yaml-cpp/yaml.h>            // IWYU pragma: keep
#include <future>                     // for future
#include <ostream>                    // for operator<<, basic_ostream, size_t
#include <unordered_map>              // for _Node_iterator, unordered_map
#include <utility>                    // for move, pair
#include <vector>                     // for vector

#include "SessionState.h"  // for SessionState
#include "StoredState.h"   // for deserialize_stored_session, ser...
#include "includes/ServiceConfigLoader.h"  // for ServiceConfigLoader
#include "magma_logging.h"                 // for MERROR, MLOG

namespace magma {
class StaticRuleStore;
}

namespace magma {
namespace lte {

RedisStoreClient::RedisStoreClient(std::shared_ptr<cpp_redis::client> client,
                                   const std::string& redis_table,
                                   std::shared_ptr<StaticRuleStore> rule_store)
    : client_(client), redis_table_(redis_table), rule_store_(rule_store) {}

bool RedisStoreClient::try_redis_connect() {
  ServiceConfigLoader loader;
  auto config = loader.load_service_config("redis");
  auto port = config["port"].as<uint32_t>();
  auto addr = config["bind"].as<std::string>();
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
  if (!client_->is_connected()) {
    auto connected = try_redis_connect();
    if (!connected) {
      throw RedisReadFailed();
    }
  }

  std::unordered_map<std::string, std::future<cpp_redis::reply>> futures;
  for (const std::string& key : subscriber_ids) {
    futures[key] = client_->hget(redis_table_, key);
  }

  client_->sync_commit();

  SessionMap session_map;
  for (const std::string& key : subscriber_ids) {
    auto reply = futures[key].get();
    if (reply.is_error()) {
      MLOG(MERROR) << "RedisStoreClient: Unable to get value for key " << key;
      throw RedisReadFailed();
    }
    if (reply.is_null() || !reply.is_string()) {
      // value just doesn't exist
      session_map[key] = SessionVector{};
    } else {
      session_map[key] = deserialize_session_vec(reply.as_string());
    }
  }
  return session_map;
}

SessionMap RedisStoreClient::read_all_sessions() {
  if (!client_->is_connected()) {
    auto connected = try_redis_connect();
    if (!connected) {
      throw RedisReadFailed();
    }
  }
  SessionMap session_map;
  auto hgetall_future = client_->hgetall(redis_table_);
  client_->sync_commit();

  auto reply = hgetall_future.get();
  if (reply.is_error()) {
    MLOG(MERROR) << "unable to read all sessions from redis";
    return session_map;
  }
  auto array = reply.as_array();
  for (size_t i = 0; i < array.size(); i += 2) {
    auto key_reply = array[i];
    if (!key_reply.is_string()) {
      MLOG(MERROR) << "Non string key found in sessions from redis";
      continue;
    }
    auto key = key_reply.as_string();
    auto value_reply = array[i + 1];
    if (!value_reply.is_string()) {
      MLOG(MERROR) << "RedisStoreClient: Unable to get value for key " << key;
      session_map[key] = SessionVector{};
    } else {
      session_map[key] = deserialize_session_vec(value_reply.as_string());
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
  if (!client_->is_connected()) {
    auto connected = try_redis_connect();
    if (!connected) {
      throw RedisWriteFailed();
    }
  }
  std::vector<std::string> keys;
  for (auto& it : session_map) {
    keys.push_back(it.first);
  }
  client_->watch(keys);

  // Set MULTI command.
  // Subsequent commands end up being queued for atomic execution with EXEC.
  // Together with WATCH, if one of the keys we intend to set are modified, then
  // the entire EXEC does not execute.
  client_->multi();

  // Queue up HSET commands after we've set up some sort of safety
  // guarantees.
  std::vector<std::string> keys_to_delete;
  for (auto& it : session_map) {
    if (it.second.empty()) {
      // if session is empty we shouldn't write back this subs anymore
      keys_to_delete.push_back(it.first);
      continue;
    }
    client_->hset(redis_table_, it.first, serialize_session_vec(it.second));
  }
  if (!keys_to_delete.empty()) {
    client_->hdel(redis_table_, keys_to_delete);
  }
  auto exec_future = client_->exec();
  client_->sync_commit();

  auto reply = exec_future.get();
  if (!reply.ok()) {
    MLOG(MERROR) << "Failed to write sessions to Redis.";
    return false;
  }
  return true;
}

std::string RedisStoreClient::serialize_session_vec(
    SessionVector& session_vec) {
  folly::dynamic marshaled = folly::dynamic::array;
  for (auto& session_ptr : session_vec) {
    auto stored_session = session_ptr->marshal();
    marshaled.push_back(serialize_stored_session(stored_session));
  }
  std::string serialized = folly::toJson(marshaled);
  return serialized;
}

SessionVector RedisStoreClient::deserialize_session_vec(
    std::string serialized) {
  SessionVector session_vec;
  auto folly_serialized = folly::StringPiece(serialized);
  try {
    folly::dynamic marshaled = folly::parseJson(folly_serialized);
    for (auto& it : marshaled) {
      auto stored_session = deserialize_stored_session(it.getString());
      session_vec.push_back(
          SessionState::unmarshal(stored_session, *rule_store_));
    }
  } catch (std::exception const& e) {
    // Very rare but we've seen a crash here
    MLOG(MERROR) << "Exception " << e.what()
                 << " parsing serialized states as JSON " << folly_serialized;
  }
  return session_vec;
}

}  // namespace lte
}  // namespace magma
