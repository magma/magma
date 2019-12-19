// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <folly/IPAddress.h>
#include <folly/dynamic.h>

#include <devmand/channels/http/Channel.h>
#include <devmand/devices/Device.h>

namespace devmand {

class Application;

namespace devices {
namespace frinx {

class Device : public devices::Device {
 public:
  // TODO object
  Device(
      Application& application,
      const Id& id_,
      bool readonly_,
      const std::string& controllerHost,
      const int controllerPort,
      const folly::IPAddress& deviceIp_,
      const int devicePort_,
      const std::string& authorization_,
      const std::string& deviceId_,
      const std::string& transportType_,
      const std::string& deviceType_,
      const std::string& deviceVersion_,
      const std::string& deviceUsername_,
      const std::string& devicePassword_);
  Device() = delete;
  virtual ~Device();
  Device(const Device&) = delete;
  Device& operator=(const Device&) = delete;
  Device(Device&&) = delete;
  Device& operator=(Device&&) = delete;

  static std::shared_ptr<devices::Device> createDevice(
      Application& app,
      const cartography::DeviceConfig& deviceConfig);

 public:
  std::shared_ptr<Datastore> getOperationalDatastore() override;

 protected:
  void setIntendedDatastore(const folly::dynamic& config) override;

 private:
  void connect();
  void checkConnection();

 private:
  channels::http::Channel channel;
  bool connected{false};
  httplib::Headers headers;

  folly::IPAddress deviceIp;
  int devicePort;
  std::string deviceId;
  std::string transportType;
  std::string deviceType;
  std::string deviceVersion;
  std::string deviceUsername;
  std::string devicePassword;
};

} // namespace frinx
} // namespace devices
} // namespace devmand
