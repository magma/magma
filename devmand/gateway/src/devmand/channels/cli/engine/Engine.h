// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/channels/Engine.h>
#include <devmand/channels/cli/CliThreadWheelTimekeeper.h>
#include <devmand/devices/cli/schema/ModelRegistry.h>
#include <devmand/devices/cli/translation/PluginRegistry.h>
#include <devmand/devices/cli/translation/ReaderRegistry.h>
#include <devmand/devices/cli/translation/WriterRegistry.h>
#include <devmand/magma/DevConf.h>
#include <folly/Executor.h>
#include <folly/futures/ThreadWheelTimekeeper.h>
#include <atomic>

namespace devmand {
namespace channels {
namespace cli {

using namespace std;
using devmand::channels::cli::CliThreadWheelTimekeeper;
using devmand::devices::cli::ModelRegistry;
using namespace devmand::devices::cli;

static atomic<bool> loggingInitialized(false);
static atomic<bool> sshInitialized(false);

class Engine : public channels::Engine {
 public:
  Engine(folly::dynamic pluginConfig);
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
  shared_ptr<folly::Executor> pluginExecutor;

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
  unique_ptr<PluginRegistry> pluginRegistry;

 public:
  /*
   * Get executor for cli layer
   */
  shared_ptr<folly::Executor> getExecutor(
      executorRequestType requestType) const;

  shared_ptr<ModelRegistry> getModelRegistry() const;

  shared_ptr<DeviceContext> getDeviceContext(const DeviceType& type);

  unique_ptr<ReaderRegistry> getReaderRegistry(
      shared_ptr<DeviceContext> deviceCtx);

  unique_ptr<WriterRegistry> getWriterRegistry(
      shared_ptr<DeviceContext> deviceCtx);
};

} // namespace cli
} // namespace channels
} // namespace devmand
