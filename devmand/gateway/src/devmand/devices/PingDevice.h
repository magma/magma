// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <devmand/channels/ping/Channel.h>
#include <devmand/devices/Device.h>

namespace devmand {
namespace devices {

class PingDevice : public Device {
 public:
  PingDevice(
      Application& application,
      const Id& id,
      bool readonly_,
      const folly::IPAddress& ip_);
  PingDevice() = delete;
  virtual ~PingDevice() = default;
  PingDevice(const PingDevice&) = delete;
  PingDevice& operator=(const PingDevice&) = delete;
  PingDevice(PingDevice&&) = delete;
  PingDevice& operator=(PingDevice&&) = delete;

  static std::shared_ptr<devices::Device> createDevice(
      Application& app,
      const cartography::DeviceConfig& deviceConfig);

 public:
  std::shared_ptr<State> getState() override;

 protected:
  void setConfig(const folly::dynamic& config) override {
    (void)config;
    LOG(ERROR) << "set config on unconfigurable device";
  }

 protected:
  channels::ping::Channel channel;
};

} // namespace devices
} // namespace devmand
