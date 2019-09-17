// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

#pragma once

#include <devmand/devices/Device.h>

namespace devmand {
namespace devices {

class EchoDevice : public Device {
 public:
  EchoDevice(Application& application, const Id& id);

  EchoDevice() = delete;
  virtual ~EchoDevice() = default;
  EchoDevice(const EchoDevice&) = delete;
  EchoDevice& operator=(const EchoDevice&) = delete;
  EchoDevice(EchoDevice&&) = delete;
  EchoDevice& operator=(EchoDevice&&) = delete;

  static std::unique_ptr<devices::Device> createDevice(
      Application& app,
      const cartography::DeviceConfig& deviceConfig);

 public:
  std::shared_ptr<State> getState() override;

 protected:
  void setConfig(const folly::dynamic& config) override;

 private:
  folly::dynamic state;
};

} // namespace devices
} // namespace devmand
