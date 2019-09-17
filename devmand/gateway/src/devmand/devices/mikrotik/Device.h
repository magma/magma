// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

#pragma once

#include <devmand/devices/Snmpv2Device.h>

#include <devmand/channels/mikrotik/Channel.h>

namespace devmand {
namespace devices {
namespace mikrotik {

class Device : public Snmpv2Device {
 public:
  Device(
      Application& application,
      const Id& id,
      const folly::IPAddress& _ip,
      const std::string& _username,
      const std::string& _password,
      const channels::snmp::Peer& peer,
      const channels::snmp::Community& community,
      const channels::snmp::Version& version,
      const std::string& passphrase = "",
      const std::string& securityName = "",
      const channels::snmp::SecurityLevel& securityLevel = "",
      oid proto[] = {});

  Device() = delete;
  virtual ~Device() = default;
  Device(const Device&) = delete;
  Device& operator=(const Device&) = delete;
  Device(Device&&) = delete;
  Device& operator=(Device&&) = delete;

  static std::unique_ptr<devices::Device> createDevice(
      Application& app,
      const cartography::DeviceConfig& deviceConfig);

 public:
  std::shared_ptr<State> getState() override;

 protected:
  void setConfig(const folly::dynamic& config) override {
    (void)config;
    LOG(ERROR) << "set config on unconfigurable device";
  }

 private:
  std::shared_ptr<channels::mikrotik::Channel> mikrotikCh;
};

} // namespace mikrotik
} // namespace devices
} // namespace devmand
