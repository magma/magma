// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/devices/EchoDevice.h>

namespace devmand {
namespace devices {

std::unique_ptr<devices::Device> EchoDevice::createDevice(
    Application& app,
    const cartography::DeviceConfig& deviceConfig) {
  return std::make_unique<devices::EchoDevice>(app, deviceConfig.id);
}

EchoDevice::EchoDevice(Application& application, const Id& id_)
    : Device(application, id_) {}

void EchoDevice::setConfig(const folly::dynamic& config) {
  state = config;
}

std::shared_ptr<State> EchoDevice::getState() {
  auto stateCopy = State::make(*reinterpret_cast<MetricSink*>(&app), getId());
  stateCopy->update([this](auto& lockedState) { lockedState = state; });
  return stateCopy;
}

} // namespace devices
} // namespace devmand
