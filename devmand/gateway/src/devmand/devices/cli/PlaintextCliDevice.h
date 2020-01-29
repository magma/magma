// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/Application.h>
#include <devmand/channels/cli/Channel.h>
#include <devmand/channels/cli/Command.h>
#include <devmand/channels/cli/ReadCachingCli.h>
#include <devmand/devices/Device.h>

namespace devmand {
namespace devices {
namespace cli {

using namespace devmand::channels::cli;

class PlaintextCliDevice : public Device {
 public:
  PlaintextCliDevice(
      Application& application,
      Engine& engine,
      const Id id,
      const std::string stateCommand,
      const std::shared_ptr<Channel> channel,
      const std::shared_ptr<CliCache> cmdCache = ReadCachingCli::createCache());
  PlaintextCliDevice() = delete;
  virtual ~PlaintextCliDevice() = default;
  PlaintextCliDevice(const PlaintextCliDevice&) = delete;
  PlaintextCliDevice& operator=(const PlaintextCliDevice&) = delete;
  PlaintextCliDevice(PlaintextCliDevice&&) = delete;
  PlaintextCliDevice& operator=(PlaintextCliDevice&&) = delete;

  static std::unique_ptr<devices::Device> createDevice(
      Application& app,
      const cartography::DeviceConfig& deviceConfig);

  // visible for testing
  static std::unique_ptr<devices::Device> createDeviceWithEngine(
      Application& app,
      const cartography::DeviceConfig& deviceConfig,
      Engine& engine);

 public:
  std::shared_ptr<Datastore> getOperationalDatastore() override;

 protected:
  void setIntendedDatastore(const folly::dynamic& config) override {
    (void)config;
    MLOG(MERROR) << "[" << id << "] "
                 << "set config on unconfigurable device";
  }

 private:
  std::shared_ptr<Channel> channel;
  const ReadCommand stateCommand;
  std::shared_ptr<CliCache> cmdCache;
  std::shared_ptr<folly::Executor> executor;
};

} // namespace cli
} // namespace devices
} // namespace devmand
