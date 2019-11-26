// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <spdlog/sinks/base_sink.h>
#include <mutex>

namespace devmand {
namespace channels {
namespace cli {

using namespace spdlog::level;

// Usage:
//  auto logger = spdlog::create<Spd2Glog>("spd_logger");
//  logger->info("log me");
class Spd2Glog : public spdlog::sinks::base_sink<std::mutex> {
 public:
  void toGlog(const spdlog::details::log_msg& msg);

 protected:
  void _sink_it(const spdlog::details::log_msg& msg) override;
  void flush() override;
};
} // namespace cli
} // namespace channels
} // namespace devmand
