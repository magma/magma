// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/devices/mikrotik/Device.h>

#include <devmand/Application.h>
#include <devmand/StringUtils.h>
#include <devmand/devices/DemoDevice.h>
#include <devmand/devices/mikrotik/Mib.h>
#include <devmand/models/device/Model.h>
#include <devmand/models/wifi/Model.h>

namespace devmand {
namespace devices {
namespace mikrotik {

static constexpr const unsigned short mikrotikPort = 8728;

std::unique_ptr<devices::Device> Device::createDevice(
    Application& app,
    const cartography::DeviceConfig& deviceConfig) {
  const auto& channelConfigs = deviceConfig.channelConfigs;
  const auto& otherKv = channelConfigs.at("other").kvPairs;
  const auto& snmpKv = channelConfigs.at("snmp").kvPairs;
  return std::make_unique<Device>(
      app,
      deviceConfig.id,
      folly::IPAddress(deviceConfig.ip),
      otherKv.at("username"),
      otherKv.at("password"),
      deviceConfig.ip,
      snmpKv.at("community"),
      snmpKv.at("version"));
}

Device::Device(
    Application& application,
    const Id& id_,
    const folly::IPAddress& _ip,
    const std::string& _username,
    const std::string& _password,
    const channels::snmp::Peer& peer,
    const channels::snmp::Community& community,
    const channels::snmp::Version& version,
    const std::string& passphrase,
    const std::string& securityName,
    const channels::snmp::SecurityLevel& securityLevel,
    oid proto[])
    : Snmpv2Device(
          application,
          id_,
          peer,
          community,
          version,
          passphrase,
          securityName,
          securityLevel,
          proto),
      mikrotikCh(std::make_shared<channels::mikrotik::Channel>(
          application.getEventBase(),
          folly::SocketAddress(_ip, mikrotikPort),
          _username,
          _password)) {
  mikrotikCh->connect();
}

std::shared_ptr<State> Device::getState() {
  auto state = Snmpv2Device::getState();

  // fbc-symphony-device ######################################################
  devmand::models::device::Model::init(state->update());

  auto& geol = state->update()["fbc-symphony-device:system"]["geo-location"];

  // TODO check units and make conversions for all of these.
  state->addRequest(
      Mib::getLongtitude(channel)
          .thenValue([&geol](auto v) { geol["longitude"] = v; })
          .thenError(
              folly::tag_t<channels::snmp::Exception>{},
              [&geol, this](const channels::snmp::Exception&) {
                auto v = lookup(
                    "fbc-symphony-device:system/geo-location/longitude");
                if (not v.isNull()) {
                  geol["longitude"] = v.asString();
                }
              }));
  state->addRequest(
      Mib::getLatitude(channel)
          .thenValue([&geol](auto v) { geol["latitude"] = v; })
          .thenError(
              folly::tag_t<channels::snmp::Exception>{},
              [&geol, this](const channels::snmp::Exception&) {
                auto v =
                    lookup("fbc-symphony-device:system/geo-location/latitude");
                if (not v.isNull()) {
                  geol["latitude"] = v.asString();
                }
              }));
  state->addRequest(
      Mib::getAltitude(channel)
          .thenValue([&geol](auto v) { geol["height"] = v; })
          .thenError(
              folly::tag_t<channels::snmp::Exception>{},
              [&geol, this](const channels::snmp::Exception&) {
                auto v =
                    lookup("fbc-symphony-device:system/geo-location/height");
                if (not v.isNull()) {
                  geol["height"] = v.asString();
                }
              }));

  // openconfig-wifi ##########################################################
  devmand::models::wifi::Model::init(state->update());

  auto& papRoot = state->update()["openconfig-ap-manager:provision-aps"];
  auto& paps = papRoot["provision-ap"];
  folly::dynamic& pap = paps[0];
  auto& papc = pap["config"];
  auto& papt = pap["state"];
  auto& japRoot = state->update()["openconfig-ap-manager:joined-aps"];
  auto& japs = japRoot["joined-ap"];
  folly::dynamic& jap = japs[0];
  auto& japt = jap["state"];

  // always enabled
  japt["enabled"] = true;
  // just using uptime snmp success to indicate if the device is up.
  japt["opstate"] = "openconfig-wifi-types:DOWN";

  state->addRequest(
      Mib::getBaseMac(channel).thenValue([&pap, &papc, &papt, &japt](auto v) {
        auto hex = StringUtils::asHexString(v, ":");
        pap["mac"] = hex;
        papc["mac"] = hex;
        papt["mac"] = hex;
        japt["mac"] = hex;
      }));
  state->addRequest(Mib::getFirmwareVersion(channel).thenValue(
      [&japt](auto v) { japt["software-version"] = v; }));
  state->addRequest(Mib::getSerialNumber(channel).thenValue(
      [&japt](auto v) { japt["serial"] = v; }));
  state->addRequest(Mib::getUpTime(channel).thenValue([&japt](auto v) {
    japt["uptime"] = v;
    japt["opstate"] = "openconfig-wifi-types:UP";
  }));
  state->addRequest(
      Mib::getModel(channel).thenValue([&japt](auto v) { japt["model"] = v; }));

  state->addFinally([state, &papc, &papt, &jap, &japt]() {
    auto* field = state->update().get_ptr("ietf-system:system");
    if (field != nullptr and ((field = field->get_ptr("name")) != nullptr)) {
      auto hostname = field->asString();
      papc["hostname"] = hostname;
      papt["hostname"] = hostname;
      jap["hostname"] = hostname;
      japt["hostname"] = hostname;
    }
  });

  state->addRequest(Mib::getIpv4Address(channel).thenValue(
      [&japt](auto v) { japt["ipv4"] = v; }));

  state->addRequest(Mib::getIpv6Address(channel).thenValue([&japt](auto v) {
    japt["ipv6"] =
        folly::IPAddress::fromBinary(
            folly::ByteRange(
                reinterpret_cast<const unsigned char*>(v.data()), v.size()))
            .str();
  }));

  /* TODO more ones to add
  papc["country-code"] = "US";
  papt["country-code"] = "US";
  */

  // Note we don't have a way on mikrotik to query this.
  // japt["power-source"] = "UNKNOWN";

  return state;
}

} // namespace mikrotik
} // namespace devices
} // namespace devmand
