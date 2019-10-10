// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/devices/Device.h>

#include <iostream>

#include <folly/json.h>

#include <devmand/Application.h>
#include <devmand/ErrorHandler.h>

namespace devmand {
namespace devices {

Device::Device(Application& application, const Id& id_)
    : app(application), id(id_) {}

Id Device::getId() const {
  return id;
}

DeviceConfigType Device::getDeviceConfigType() const {
  return DeviceConfigType::YangJson;
}

void Device::setNativeConfig(const std::string&) {
  LOG(ERROR) << "set native config called on device that doesn't implement.";
}

folly::dynamic Device::lookup(const YangPath& path) const {
  return YangUtils::lookup(lastConfig, path);
}

void Device::updateSharedView(SharedUnifiedView& sharedUnifiedView) {
  Id idL = id;
  ErrorHandler::thenError(
      getState()->collect().thenValue([idL, &sharedUnifiedView](auto data) {
        sharedUnifiedView.withULockPtr([&idL, &data](auto uUnifiedView) {
          auto unifiedView = uUnifiedView.moveFromUpgradeToWrite();

          // TODO this is an expensive hack... fix later. Prob. just store in
          // dyn and have the magma service convert it.
          folly::dynamic dyn =
              unifiedView->find("devmand") != unifiedView->end()
              ? folly::parseJson((*unifiedView)["devmand"])
              : folly::dynamic::object;
          dyn[idL] = data;
          (*unifiedView)["devmand"] = folly::toJson(dyn);
        });
      }));
}

void Device::applyConfig(const std::string& config) {
  LOG(INFO) << "Applying config '" << config;
  if (not config.empty()) {
    switch (getDeviceConfigType()) {
      case DeviceConfigType::YangJson: {
        folly::dynamic json = folly::parseJson(config);
        setConfig(json);
        lastConfig = json;
        break;
      }
      case DeviceConfigType::NativeConfigJson:
        folly::dynamic json = folly::parseJson(config);
        auto* nativeConfig = json.get_ptr("native_config");
        if (nativeConfig != nullptr and nativeConfig->isString()) {
          setNativeConfig(nativeConfig->asString());
        }
        lastConfig = json;
        break;
    }
  }
}

} // namespace devices
} // namespace devmand
