// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <devmand/DhcpdConfig.h>
#include <devmand/devices/Device.h>

namespace devmand {
namespace devices {

class DcsgDevice : public Device {
 public:
  DcsgDevice(Application& application, const Id& id_, Host& host_);

  DcsgDevice() = delete;
  virtual ~DcsgDevice();
  DcsgDevice(const DcsgDevice&) = delete;
  DcsgDevice& operator=(const DcsgDevice&) = delete;
  DcsgDevice(DcsgDevice&&) = delete;
  DcsgDevice& operator=(DcsgDevice&&) = delete;

  static std::unique_ptr<devices::Device> createDevice(
      Application& app,
      const cartography::DeviceConfig& deviceConfig);

 public:
  std::shared_ptr<State> getState() override;

  DeviceConfigType getDeviceConfigType() const override;

 protected:
  void setConfig(const folly::dynamic& config) override;
  void setNativeConfig(const std::string& config) override;

 private:
  Host host;
};

} // namespace devices
} // namespace devmand
