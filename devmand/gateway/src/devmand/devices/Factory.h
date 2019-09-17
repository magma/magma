// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

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
  using PlatformBuilder = std::function<std::unique_ptr<devices::Device>(
      Application& application,
      const cartography::DeviceConfig& deviceConfig)>;

  std::unique_ptr<devices::Device> createDevice(
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
