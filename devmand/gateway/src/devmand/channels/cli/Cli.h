// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/channels/cli/Command.h>
#include <folly/futures/Future.h>
#include <folly/futures/Promise.h>
#include <magma_logging.h>

#include <chrono>
#include <thread>
#include <vector>

namespace devmand {
namespace channels {
namespace cli {
class Cli {
 public:
  Cli() = default;
  virtual ~Cli() = default;
  Cli(const Cli&) = delete;
  Cli& operator=(const Cli&) = delete;
  Cli(Cli&&) = delete;
  Cli& operator=(Cli&&) = delete;

 public:
  virtual folly::Future<std::string> executeAndRead(const Command& cmd) = 0;

  virtual folly::Future<std::string> execute(const Command& cmd) = 0;
};

} // namespace cli
} // namespace channels
} // namespace devmand