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
   * start_loop is the main function to call to initiate a load loop. Based on
   * the given loop interval length, this function will load the policies from
   * redis, and call the processor callback.
   */
  void start_loop(
    std::function<void(std::vector<PolicyRule>)> processor,
    uint32_t loop_interval_seconds);

  /**
   * Stop the config loop on the next loop
   */
  void stop();
private:
  std::atomic<bool> is_running_;
};
}
