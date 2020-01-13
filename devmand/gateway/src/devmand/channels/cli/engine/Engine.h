// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#define LOG_WITH_GLOG

#include <devmand/channels/Engine.h>
#include <devmand/channels/cli/CliThreadWheelTimekeeper.h>
#include <devmand/devices/cli/ModelRegistry.h>
#include <folly/Executor.h>
#include <folly/futures/ThreadWheelTimekeeper.h>
#include <magma_logging.h>
#include <atomic>

namespace devmand {
namespace channels {
namespace cli {

using namespace std;
using devmand::channels::cli::CliThreadWheelTimekeeper;
using devmand::devices::cli::ModelRegistry;

static atomic<bool> loggingInitialized(false);
static atomic<bool> sshInitialized(false);

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

  shared_ptr<CliThreadWheelTimekeeper> timekeeper;
  shared_ptr<folly::Executor> sshCliExecutor;
  shared_ptr<folly::Executor> commonExecutor;
  shared_ptr<folly::Executor> kaCliExecutor;

  shared_ptr<CliThreadWheelTimekeeper> getTimekeeper();

  enum executorRequestType {
    sshCli,
    paCli,
    rcCli,
    tcCli,
    ttCli,
    lCli,
    qCli,
    rCli,
    kaCli,
    plaintextCliDevice
  };

 private:
  shared_ptr<ModelRegistry> mreg;

 public:
  shared_ptr<folly::Executor> getExecutor(
      executorRequestType requestType) const;

  shared_ptr<ModelRegistry> getModelRegistry() const;
};

} // namespace cli
} // namespace channels
} // namespace devmand
