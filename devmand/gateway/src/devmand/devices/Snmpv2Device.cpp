// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

#include <iostream>
#include <stdexcept>

#include <devmand/ErrorHandler.h>
#include <devmand/channels/snmp/IfMib.h>
#include <devmand/devices/Snmpv2Device.h>
#include <devmand/devices/State.h>

namespace devmand {
namespace devices {

std::unique_ptr<devices::Device> Snmpv2Device::createDevice(
    Application& app,
    const cartography::DeviceConfig& deviceConfig) {
  const auto& channelConfigs = deviceConfig.channelConfigs;
  const auto& snmpKv = channelConfigs.at("snmp").kvPairs;
  return std::make_unique<devices::Snmpv2Device>(
      app,
      deviceConfig.id,
      deviceConfig.ip,
      snmpKv.at("community"),
      snmpKv.at("version"));
}

Snmpv2Device::Snmpv2Device(
    Application& application,
    const Id& id_,
    const channels::snmp::Peer& peer,
    const channels::snmp::Community& community,
    const channels::snmp::Version& version,
    const std::string& passphrase,
    const std::string& securityName,
    const channels::snmp::SecurityLevel& securityLevel,
    oid proto[])
    : Device(application, id_),
      channel(
          peer,
          community,
          version,
          passphrase,
          securityName,
          securityLevel,
          proto) {}

std::shared_ptr<State> Snmpv2Device::getState() {
  auto state = State::make(app, *this);
  state->setStatus(false);
  auto& system = state->update()["ietf-system:system"] = folly::dynamic::object;
  auto& interfaces = state->update()["openconfig-interfaces:interfaces"] =
      folly::dynamic::object;
  interfaces["interface"] = folly::dynamic::array;

  using IfMib = devmand::channels::snmp::IfMib;

  state->addRequest(IfMib::getSystemName(channel).thenValue(
      [&system](auto v) { system["name"] = v; }));
  state->addRequest(IfMib::getSystemContact(channel).thenValue(
      [&system](auto v) { system["contact"] = v; }));
  state->addRequest(IfMib::getSystemLocation(channel).thenValue(
      [&system](auto v) { system["location"] = v; }));
  state->addRequest(
      IfMib::getInterfaceNames(channel).thenValue([&interfaces](auto results) {
        for (auto result : results) {
          updateInterface(interfaces, result.index, "name", result.name);
        }
      }));
  state->addRequest(IfMib::getInterfaceStatuses(channel).thenValue(
      [&interfaces](auto results) {
        for (auto result : results) {
          updateInterface(
              interfaces, result.index, "oper-status", result.status);
        }
      }));
  return state;
}

void Snmpv2Device::updateInterface(
    folly::dynamic& interfaces,
    int index,
    const std::string& key,
    const std::string& value) {
  for (auto& interface : interfaces["interface"]) {
    if (interface["ifindex"] == index) {
      interface[key] = value;
      return;
    }
  }

  folly::dynamic interface = folly::dynamic::object;
  auto& state = interface["state"] = folly::dynamic::object;
  state["ifindex"] = index;
  state[key] = value;
  // TODO this is wrong... make sure i dont break ui first tho...
  interfaces["interface"].push_back(state);
}

} // namespace devices
} // namespace devmand
