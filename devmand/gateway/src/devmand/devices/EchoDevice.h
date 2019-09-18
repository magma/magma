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
