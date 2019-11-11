/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include <atomic>
#include <functional>
#include <cpp_redis/cpp_redis>
#include <lte/protos/policydb.pb.h>

namespace magma {
using namespace lte;
/**
 * PolicyLoader is used to sync policies with Redis every so often
 */
class PolicyLoader {
public:

  /**
   * load_config_async initiates an async loop to load config. Based on
   * the given loop interval length, this function will load the policies from
   * redis, and call the processor callback.
   */
  void load_config_async(
      const std::function<void(std::vector<PolicyRule>)>& processor,
      uint32_t loop_interval_seconds);

  /**
   * Stop the config loop on the next loop. Blocks until the async thread
   * completes.
   */
  void stop();
private:
  std::atomic<bool> is_running_;
  std::thread redis_client_thread_;
};
}
