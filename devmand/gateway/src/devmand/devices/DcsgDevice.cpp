// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <experimental/filesystem>

#include <folly/Format.h>
#include <folly/GLog.h>

#include <devmand/Application.h>
#include <devmand/FileUtils.h>
#include <devmand/devices/DcsgDevice.h>

namespace devmand {
namespace devices {

const char* deviceConfigFilePathTemplate = "/var/www/{}.conf";

std::unique_ptr<devices::Device> DcsgDevice::createDevice(
    Application& app,
    const cartography::DeviceConfig& deviceConfig) {
  const auto& channelConfigs = deviceConfig.channelConfigs;
  const auto& otherKv = channelConfigs.at("other").kvPairs;
  Host host{deviceConfig.id,
            folly::MacAddress(otherKv.at("mac")),
            folly::IPAddress(deviceConfig.ip)};
  return std::make_unique<devices::DcsgDevice>(app, deviceConfig.id, host);
}

DcsgDevice::DcsgDevice(Application& application, const Id& id_, Host& host_)
    : Device(application, id_), host(host_) {
  app.getDhcpdConfig().add(host);
  FileUtils::mkdir(
      std::experimental::filesystem::path(deviceConfigFilePathTemplate)
          .parent_path());
}

DcsgDevice::~DcsgDevice() {
  app.getDhcpdConfig().remove(host);
}

DeviceConfigType DcsgDevice::getDeviceConfigType() const {
  return DeviceConfigType::NativeConfigJson;
}

std::shared_ptr<State> DcsgDevice::getState() {
  return State::make(app, getId());
}

void DcsgDevice::setConfig(const folly::dynamic&) {
  LOG(ERROR) << "set config called on device not supporting json configs";
}

void DcsgDevice::setNativeConfig(const std::string& config) {
  FileUtils::write(
      folly::sformat(deviceConfigFilePathTemplate, host.name), config);
}

} // namespace devices
} // namespace devmand
