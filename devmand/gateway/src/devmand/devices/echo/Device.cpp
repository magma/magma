// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/devices/echo/Device.h>

namespace devmand {
namespace devices {
namespace echo {

std::shared_ptr<devices::Device> Device::createDevice(
    Application& app,
    const cartography::DeviceConfig& deviceConfig) {
  return std::make_unique<devices::echo::Device>(
      app, deviceConfig.id, deviceConfig.readonly);
}

Device::Device(Application& application, const Id& id_, bool readonly_)
    : devices::Device(application, id_, readonly_) {}

void Device::setIntendedDatastore(const folly::dynamic& config) {
  state = config;
}

std::shared_ptr<Datastore> Device::getOperationalDatastore() {
  auto stateCopy =
      Datastore::make(*reinterpret_cast<MetricSink*>(&app), getId());
  stateCopy->update([this](auto& lockedDatastore) { lockedDatastore = state; });
  return stateCopy;
}

} // namespace echo
} // namespace devices
} // namespace devmand
