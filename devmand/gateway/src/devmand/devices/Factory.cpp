// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <boost/algorithm/string.hpp>

#include <devmand/devices/Factory.h>
#include <devmand/error/ErrorHandler.h>
#include <devmand/utils/StringUtils.h>

namespace devmand {
namespace devices {

Factory::Factory(Application& application) : app(application) {}

void Factory::addPlatform(
    const std::string& platform,
    PlatformBuilder platformBuilder) {
  std::string platformLowerCase = platform;
  boost::algorithm::to_lower(platformLowerCase);
  if (not platformBuilders.emplace(platformLowerCase, platformBuilder).second) {
    LOG(ERROR) << "Failed to add platform " << platform;
    throw std::runtime_error("Failed to add device platform");
  }
}

void Factory::setDefaultPlatform(PlatformBuilder defaultPlatformBuilder_) {
  defaultPlatformBuilder = defaultPlatformBuilder_;
}

std::shared_ptr<devices::Device> Factory::createDevice(
    const cartography::DeviceConfig& deviceConfig) {
  LOG(INFO) << "Loading device " << deviceConfig.id << " on platform "
            << deviceConfig.platform << " with ip " << deviceConfig.ip;

  std::string platformLowerCase = deviceConfig.platform;
  boost::algorithm::to_lower(platformLowerCase);
  PlatformBuilder builder{nullptr};
  auto builderIt = platformBuilders.find(platformLowerCase);
  if (builderIt == platformBuilders.end()) {
    LOG(INFO) << "Didn't find matching platform so using default.";
    builder = defaultPlatformBuilder;
  } else {
    builder = builderIt->second;
  }

  std::shared_ptr<devices::Device> device{nullptr};
  if (builder != nullptr) {
    ErrorHandler::executeWithCatch([this, &builder, &deviceConfig, &device]() {
      device = builder(app, deviceConfig);
    });
  }

  if (device != nullptr) {
    return device;
  } else {
    LOG(ERROR) << "Error adding device " << deviceConfig;
    throw std::runtime_error("Error adding device");
  }
}

} // namespace devices
} // namespace devmand
