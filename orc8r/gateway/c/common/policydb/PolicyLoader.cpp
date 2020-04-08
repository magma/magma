/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include "RedisMap.hpp"
#include "Serializers.h"
#include "PolicyLoader.h"
#include "ServiceConfigLoader.h"
#include "magma_logging.h"

namespace magma {

bool try_redis_connect(cpp_redis::client& client)
{
  ServiceConfigLoader loader;
  auto config = loader.load_service_config("redis");
  auto port = config["port"].as<uint32_t>();
  try {
    client.connect(
      "127.0.0.1",
      port,
      [](
        const std::string& host,
        std::size_t port,
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
  cpp_redis::client& client,
  RedisMap<PolicyRule>& policy_map,
  const std::function<void(std::vector<PolicyRule>)>& processor)
{
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
  if (rules.size() > 0) {
    MLOG(MDEBUG) << rules.size() << " rules synced";
  }
  processor(rules);
  return true;
}

void PolicyLoader::start_loop(
  std::function<void(std::vector<PolicyRule>)> processor,
  uint32_t loop_interval_seconds)
{
  is_running_ = true;
  auto client = std::make_shared<cpp_redis::client>();
  auto policy_map = RedisMap<PolicyRule>(
    client, "policydb:rules", get_proto_serializer(), get_proto_deserializer());
  while (is_running_) {
    do_loop(*client, policy_map, processor);
    std::this_thread::sleep_for(std::chrono::seconds(loop_interval_seconds));
  }
}

void PolicyLoader::stop()
{
  is_running_ = false;
}

}
