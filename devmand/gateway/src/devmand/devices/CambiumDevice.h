// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <folly/IPAddress.h>
#include <folly/dynamic.h>

#include <devmand/channels/cnmaestro/Channel.h>
#include <devmand/devices/Device.h>

namespace devmand {

class Application;

namespace devices {

class CambiumDevice : public Device {
 public:
  CambiumDevice(
      Application& application,
      const Id& id_,
      const folly::IPAddress& deviceIp_,
      const std::string& clientMac_,
      const std::string& clientId_,
      const std::string& clientSecret_);
  CambiumDevice() = delete;
  virtual ~CambiumDevice();
  CambiumDevice(const CambiumDevice&) = delete;
  CambiumDevice& operator=(const CambiumDevice&) = delete;
  CambiumDevice(CambiumDevice&&) = delete;
  CambiumDevice& operator=(CambiumDevice&&) = delete;

  static std::unique_ptr<devices::Device> createDevice(
      Application& app,
      const cartography::DeviceConfig& deviceConfig);

 public:
  std::shared_ptr<State> getState() override;
  void setConfig(const folly::dynamic& config) override;

 private:
  void connect();
  void checkConnection();
  folly::dynamic setupReturnData();
  void updateDevice(
      const folly::dynamic& yangIn,
      std::vector<std::string>& path,
      folly::dynamic& updateJson);
  void updateYang(
      const folly::dynamic& config,
      std::vector<std::string>& path,
      long unsigned int index,
      const folly::dynamic& original,
      folly::dynamic& updateJson);

  // TODO: Pull these out into their own file/library
  void setupOpenconfig(folly::dynamic& dataIn);
  void setupOpenconfigInterfaces(folly::dynamic& dataIn);
  void addOpenconfigInterface(folly::dynamic& data);
  void setupIpv4(folly::dynamic& dataIn, int interfaceNum);

 private:
  channels::cnmaestro::Channel channel;
  bool connected{false};
  std::string deviceId;
  folly::IPAddress deviceIp;
  std::string clientMac;
  folly::dynamic lastUpdate;
};

} // namespace devices
} // namespace devmand
