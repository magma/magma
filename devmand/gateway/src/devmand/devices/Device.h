// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <stdexcept>
#include <string>
#include <vector>

#include <folly/dynamic.h>
#include <folly/futures/Future.h>

#include <devmand/ErrorQueue.h>
#include <devmand/UnifiedView.h>
#include <devmand/YangUtils.h>
#include <devmand/cartography/DeviceConfig.h>
#include <devmand/devices/Id.h>
#include <devmand/devices/State.h>

namespace devmand {

class Application;

namespace devices {

enum class DeviceConfigType { YangJson, NativeConfigJson };

class Device {
 public:
  Device(Application& application, const Id& id_);
  Device() = delete;
  virtual ~Device() = default;
  Device(const Device&) = delete;
  Device& operator=(const Device&) = delete;
  Device(Device&&) = delete;
  Device& operator=(Device&&) = delete;

 public:
  /* This function returns in the future the state representing the device.
   * This state is a folly dynamic which can be serialized as json structured
   * by a yang datamodel. */
  virtual std::shared_ptr<State> getState() = 0;

  /* This function asynchronously modifies the shared unified view (the common
   * way of looking at and operating on the network) with the state provided by
   * get state, */
  virtual void updateSharedView(SharedUnifiedView& sharedUnifiedView);

  Id getId() const;

  /*
   * Given a string config this method parses the config and passing it on
   * to the correct handler to apply the config.
   *
   * TODO provide a path to signal errors
   */
  void applyConfig(const std::string& config);

  virtual DeviceConfigType getDeviceConfigType() const;

 protected:
  /*
   * Inherited method to override in device instances. This is called
   * by the json overload of apply config. This is normally what users will
   * implement.
   */
  virtual void setConfig(const folly::dynamic& config) = 0;

  /*
   * Inherited method to override in device instances. This is called by the
   * string override of apply config. It is not normally what should be
   * overridden as this means the data is not properly formated yang json. This
   * is useful for tempory proof of concepts for quick onboarding.
   */
  virtual void setNativeConfig(const std::string& config);

  folly::dynamic lookup(const YangPath& path) const;

 protected:
  Application& app;
  Id id;
  folly::dynamic lastConfig;
};

} // namespace devices
} // namespace devmand
