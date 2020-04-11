/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#pragma once

#include <cpp_redis/cpp_redis>
#include <folly/Format.h>
#include <folly/dynamic.h>
#include <folly/json.h>

#include "StoreClient.h"
#include "StoredState.h"
#include "ServiceConfigLoader.h"

namespace magma {
namespace lte {

class RedisReadFailed : public std::exception {
 public:
  RedisReadFailed() = default;
};

/**
 * Persistent StoreClient used to allow stateless session_manager to function
 */
class RedisStoreClient final : public StoreClient {
 public:
  RedisStoreClient(
      std::shared_ptr<cpp_redis::client> client, const std::string& redis_table,
      std::shared_ptr<StaticRuleStore> rule_store);

  RedisStoreClient(RedisStoreClient const&) = delete;
  RedisStoreClient(RedisStoreClient&&)      = default;
  ~RedisStoreClient()                       = default;

  bool try_redis_connect();

  SessionMap read_sessions(std::set<std::string> subscriber_ids);

  SessionMap read_all_sessions();

  bool write_sessions(SessionMap session_map);

 private:
  std::shared_ptr<cpp_redis::client> client_;
  std::string redis_table_;
  std::shared_ptr<StaticRuleStore> rule_store_;

 private:
  std::string serialize_session_vec(
      std::vector<std::unique_ptr<SessionState>>& session_vec);

  std::vector<std::unique_ptr<SessionState>> deserialize_session_vec(
      std::string serialized);
};

}  // namespace lte
}  // namespace magma
