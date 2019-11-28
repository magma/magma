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
#include <devmand/Config.h>
#include <devmand/ErrorHandler.h>

namespace devmand {
namespace devices {

Device::Device(Application& application, const Id& id_, bool readonly_)
    : app(application), id(id_), readonly(readonly_) {}

Device::~Device() {
  auto oldHostname =
      YangUtils::lookup(lastState, "ietf-system:system/hostname");
  if (oldHostname != nullptr) {
    app.getSyslogManager().removeIdentifier(oldHostname.asString(), id);
    app.getSyslogManager().restartTdAgentBitAsync();
  }
}

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

  std::weak_ptr<Device> weak(shared_from_this());
  ErrorHandler::thenError(
      getState()
          ->collect()
          .thenValue([weak](auto data) {
            if (auto shared = weak.lock()) {
              auto newHostname =
                  YangUtils::lookup(data, "ietf-system:system/hostname");
              auto oldHostname = YangUtils::lookup(
                  shared->lastState, "ietf-system:system/hostname");
              auto& sm = shared->app.getSyslogManager();
              if (newHostname == nullptr) {
                if (oldHostname != nullptr) {
                  sm.removeIdentifier(oldHostname.asString(), shared->id);
                  sm.restartTdAgentBitAsync();
                }
              } else if (oldHostname == nullptr) {
                sm.addIdentifier(newHostname.asString(), shared->id);
                sm.restartTdAgentBitAsync();
              } else if (oldHostname != newHostname) {
                sm.removeIdentifier(oldHostname.asString(), shared->id);
                sm.addIdentifier(newHostname.asString(), shared->id);
                sm.restartTdAgentBitAsync();
              }
              shared->lastState = data;
            } else {
              // The device is gone and it is its responsiblity to clean up ids.
            }
            return data;
          })
          .thenValue([idL, &sharedUnifiedView](auto data) {
            sharedUnifiedView.withULockPtr([&idL, &data](auto uUnifiedView) {
              auto unifiedView = uUnifiedView.moveFromUpgradeToWrite();

              // TODO this is an expensive hack... fix later. Prob. just store
              // in dyn and have the magma service convert it.
              folly::dynamic dyn =
                  unifiedView->find("devmand") != unifiedView->end()
                  ? folly::parseJson((*unifiedView)["devmand"])
                  : folly::dynamic::object;
              dyn[idL] = data;
              (*unifiedView)["devmand"] = folly::toJson(dyn);
              LOG(INFO) << "state for " << idL << " is " << folly::toJson(dyn);
            });
          }));
}

void Device::applyConfig(const std::string& config) {
  if (isReadonly()) {
    LOG(INFO) << "Not applying configuration on device " << id
              << " as the device is read only.";
    return;
  }

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

bool Device::isReadonly() const {
  return readonly or FLAGS_devices_readonly;
}

} // namespace devices
} // namespace devmand
