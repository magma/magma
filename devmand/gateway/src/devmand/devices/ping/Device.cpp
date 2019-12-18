// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <iostream>
#include <stdexcept>

#include <folly/Format.h>

#include <devmand/Application.h>
#include <devmand/devices/Datastore.h>
#include <devmand/devices/ping/Device.h>
#include <devmand/error/ErrorHandler.h>
#include <devmand/models/device/Model.h>

namespace devmand {
namespace devices {
namespace ping {

std::shared_ptr<devices::Device> Device::createDevice(
    Application& app,
    const cartography::DeviceConfig& deviceConfig) {
  return std::make_unique<devices::ping::Device>(
      app,
      deviceConfig.id,
      deviceConfig.readonly,
      folly::IPAddress(deviceConfig.ip));
}

Device::Device(
    Application& application,
    const Id& id_,
    bool readonly_,
    const folly::IPAddress& ip_)
    : devices::Device(application, id_, readonly_),
      channel(application.getPingEngine(ip_), ip_) {}

std::shared_ptr<Datastore> Device::getOperationalDatastore() {
  auto state = Datastore::make(app, getId());
  state->setStatus(false);
  state->update([](auto& lockedState) {
    devmand::models::device::Model::init(lockedState);
  });

  state->addRequest(channel.ping().thenValue([state](auto rtt) {
    state->update([rtt](auto& lockedState) {
      devmand::models::device::Model::addLatency(
          lockedState, "ping", "agent", "device", rtt);
    });

    state->setGauge<unsigned long int>(
        "/fbc-symphony-device:system/latencies/"
        "latency[type=ping and src=agent and dst=device]/rtt",
        rtt);
  }));
  return state;
}

} // namespace ping
} // namespace devices
} // namespace devmand
