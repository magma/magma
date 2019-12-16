// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <functional>
#include <map>
#include <memory>

#include <devmand/cartography/DeviceConfig.h>
#include <devmand/devices/Device.h>

namespace devmand {

class Application;

namespace devices {

class Factory final {
 public:
  Factory(Application& application);
  Factory() = delete;
  ~Factory() = default;
  Factory(const Factory&) = delete;
  Factory& operator=(const Factory&) = delete;
  Factory(Factory&&) = delete;
  Factory& operator=(Factory&&) = delete;

 public:
  using PlatformBuilder = std::function<std::shared_ptr<devices::Device>(
      Application& application,
      const cartography::DeviceConfig& deviceConfig)>;

  std::shared_ptr<devices::Device> createDevice(
      const cartography::DeviceConfig& deviceConfig);

  void addPlatform(
      const std::string& platform,
      PlatformBuilder platformBuilder);

  void setDefaultPlatform(PlatformBuilder defaultPlatformBuilder_);

 private:
  Application& app;
  std::map<std::string, PlatformBuilder> platformBuilders;
  PlatformBuilder defaultPlatformBuilder{nullptr};
};

} // namespace devices
} // namespace devmand
