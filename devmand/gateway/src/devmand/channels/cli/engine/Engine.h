// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#define LOG_WITH_GLOG

#include <devmand/channels/Engine.h>
#include <magma_logging.h>

namespace devmand {
namespace channels {
namespace cli {

class Engine : public channels::Engine {
 public:
  Engine();
  ~Engine() override;
  Engine(const Engine&) = delete;
  Engine& operator=(const Engine&) = delete;
  Engine(Engine&&) = delete;
  Engine& operator=(Engine&&) = delete;

  static void initLogging(
      uint32_t verbosity = MINFO,
      bool callInitMlog = false);
  static void closeLogging();
  static void initSsh();
  static void closeSsh();
};

} // namespace cli
} // namespace channels
} // namespace devmand
