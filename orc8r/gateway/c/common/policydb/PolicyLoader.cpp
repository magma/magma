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
#include "includes/PolicyLoader.h"
#include <glog/logging.h>             // for COMPACT_GOOGLE_LOG_INFO, LogMes
#include <yaml-cpp/yaml.h>            // IWYU pragma: keep
#include <chrono>                     // for seconds
#include <cpp_redis/core/client.hpp>  // for client, client::connect_state
#include <cpp_redis/misc/error.hpp>   // for redis_error
#include <cstdint>                    // for uint32_t
#include <memory>                     // for make_shared, __shared_ptr, ...
#include <ostream>                    // for operator<<, basic_ostream, ...
#include <string>                     // for string, char_traits, operator<<
#include <thread>                     // for sleep_for
#include "ObjectMap.h"                // for SUCCESS
#include "includes/RedisMap.hpp"      // for RedisMap
#include "includes/Serializers.h"     // for get_proto_deserializer, get_pro
#include "includes/ServiceConfigLoader.h"  // for ServiceConfigLoader
#include "lte/protos/policydb.pb.h"        // for PolicyRule
#include "magma_logging.h"                 // for MLOG, MERROR, MDEBUG, MINFO

namespace magma {

bool try_redis_connect(cpp_redis::client& client) {
  ServiceConfigLoader loader;
  auto config = loader.load_service_config("redis");
  auto port   = config["port"].as<uint32_t>();
  auto addr   = config["bind"].as<std::string>();
  try {
    client.connect(
        addr, port,
        [](const std::string& host, std::size_t port,
           cpp_redis::client::connect_state status) {
          if (status == cpp_redis::client::connect_state::dropped) {
            MLOG(MERROR) << "Client disconnected from " << host << ":" << port;
          }
        });
    return client.is_connected();
  } catch (const cpp_redis::redis_error& e) {
    MLOG(MERROR) << "Could not connect to redis: " << e.what();
    return false;
  }
}

bool do_loop(
    cpp_redis::client& client, RedisMap<PolicyRule>& policy_map,
    const std::function<void(std::vector<PolicyRule>)>& processor) {
  if (!client.is_connected()) {
    if (!try_redis_connect(client)) {
      return false;
    }
    MLOG(MINFO) << "Connected to redis server";
  }
  std::vector<PolicyRule> rules;
  auto result = policy_map.getall(rules);
  if (result != SUCCESS) {
    MLOG(MERROR) << "Failed to get rules from map because map error " << result;
    return false;
  }
  processor(rules);
  return true;
}

void PolicyLoader::start_loop(
    std::function<void(std::vector<PolicyRule>)> processor,
    uint32_t loop_interval_seconds) {
  is_running_     = true;
  auto client     = std::make_shared<cpp_redis::client>();
  auto policy_map = RedisMap<PolicyRule>(
      client, "policydb:rules", get_proto_serializer(),
      get_proto_deserializer());
  while (is_running_) {
    do_loop(*client, policy_map, processor);
    std::this_thread::sleep_for(std::chrono::seconds(loop_interval_seconds));
  }
}

void PolicyLoader::stop() {
  is_running_ = false;
}

}  // namespace magma
