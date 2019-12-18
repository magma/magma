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

#include <devmand/UnifiedView.h>
#include <devmand/cartography/DeviceConfig.h>
#include <devmand/devices/Datastore.h>
#include <devmand/devices/Id.h>
#include <devmand/error/ErrorQueue.h>
#include <devmand/utils/YangUtils.h>

namespace devmand {

class Application;

namespace devices {

enum class DeviceConfigType { YangJson, NativeConfigJson };

// Loosely modeled on RFC 8342
class Device : public std::enable_shared_from_this<Device> {
 public:
  Device(Application& application, const Id& id_, bool readonly_);
  Device() = delete;
  virtual ~Device();
  Device(const Device&) = delete;
  Device& operator=(const Device&) = delete;
  Device(Device&&) = delete;
  Device& operator=(Device&&) = delete;

 public:
  /* This function returns in the future the state representing the device.
   * This state is a folly dynamic which can be serialized as json structured
   * by a yang datamodel. */
  virtual std::shared_ptr<Datastore> getOperationalDatastore() = 0;

  /* This function asynchronously modifies the shared unified view (the common
   * way of looking at and operating on the network) with the state provided by
   * get state, */
  virtual void updateSharedView(SharedUnifiedView& sharedUnifiedView);

  Id getId() const;

  /*
   * Save the config to the running datastore..
   *
   * TODO provide a path to signal errors
   */
  void setRunningDatastore(const std::string& config);

  virtual DeviceConfigType getDeviceConfigType() const;

  void tryToApplyRunningDatastore();

 protected:
  /*
   * Inherited method to override in device instances. This is called
   * by the json overload of apply config. This is normally what users will
   * implement.
   */
  virtual void setIntendedDatastore(const folly::dynamic& config) = 0;

  /*
   * Inherited method to override in device instances. This is called by the
   * string override of apply config. It is not normally what should be
   * overridden as this means the data is not properly formated yang json. This
   * is useful for tempory proof of concepts for quick onboarding.
   */
  virtual void setNativeConfig(const std::string& config);

  folly::dynamic lookup(const YangPath& path) const;

  folly::dynamic getIntendedDatastore() const;

 private:
  bool isReadonly() const;

 protected:
  Application& app;
  Id id;
  const bool readonly;
  folly::dynamic runningDatastore;
  folly::dynamic intendedDatastore;
  folly::dynamic operationalDatastore;
  // TODO std::map<std::string, Platform> platforms;
};

} // namespace devices
} // namespace devmand
