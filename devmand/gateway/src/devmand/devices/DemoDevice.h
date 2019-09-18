// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <devmand/devices/Device.h>

namespace devmand {
namespace devices {

class DemoDevice : public Device {
 public:
  DemoDevice(Application& application, const Id& id);

  DemoDevice() = delete;
  virtual ~DemoDevice() = default;
  DemoDevice(const DemoDevice&) = delete;
  DemoDevice& operator=(const DemoDevice&) = delete;
  DemoDevice(DemoDevice&&) = delete;
  DemoDevice& operator=(DemoDevice&&) = delete;

  static std::unique_ptr<devices::Device> createDevice(
      Application& app,
      const cartography::DeviceConfig& deviceConfig);

 public:
  std::shared_ptr<State> getState() override;

  static folly::dynamic getDemoState();

 protected:
  void setConfig(const folly::dynamic& config) override {
    (void)config;
    LOG(ERROR) << "set config on unconfigurable device";
  }
};

} // namespace devices
} // namespace devmand
