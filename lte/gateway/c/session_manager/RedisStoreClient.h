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

#pragma once

#include <cpp_redis/cpp_redis>
#include <exception>      // IWYU pragma: keep
#include <memory>         // for shared_ptr
#include <set>            // for set
#include <string>         // for string
#include "StoreClient.h"  // for SessionMap, SessionVector, StoreClient
namespace magma {
class StaticRuleStore;
}

namespace magma {
namespace lte {

class RedisReadFailed : public std::exception {
 public:
  RedisReadFailed() = default;
};

class RedisWriteFailed : public std::exception {
 public:
  RedisWriteFailed() = default;
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

  bool is_ready() { return client_->is_connected(); };

  SessionMap read_sessions(std::set<std::string> subscriber_ids);

  SessionMap read_all_sessions();

  bool write_sessions(SessionMap session_map);

 private:
  std::shared_ptr<cpp_redis::client> client_;
  std::string redis_table_;
  std::shared_ptr<StaticRuleStore> rule_store_;

 private:
  std::string serialize_session_vec(SessionVector& session_vec);

  SessionVector deserialize_session_vec(std::string serialized);
};

}  // namespace lte
}  // namespace magma
