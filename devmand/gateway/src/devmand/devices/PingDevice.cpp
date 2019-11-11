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
#include <devmand/ErrorHandler.h>
#include <devmand/devices/PingDevice.h>
#include <devmand/devices/State.h>
#include <devmand/models/device/Model.h>

namespace devmand {
namespace devices {

std::unique_ptr<devices::Device> PingDevice::createDevice(
    Application& app,
    const cartography::DeviceConfig& deviceConfig) {
  return std::make_unique<devices::PingDevice>(
      app, deviceConfig.id, folly::IPAddress(deviceConfig.ip));
}

PingDevice::PingDevice(
    Application& application,
    const Id& id_,
    const folly::IPAddress& ip_)
    : Device(application, id_), channel(application.getPingEngine(ip_), ip_) {}

std::shared_ptr<State> PingDevice::getState() {
  auto state = State::make(app, getId());
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

} // namespace devices
} // namespace devmand
