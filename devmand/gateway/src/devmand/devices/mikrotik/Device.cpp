// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/devices/mikrotik/Device.h>

#include <devmand/Application.h>
#include <devmand/devices/mikrotik/Mib.h>
#include <devmand/models/device/Model.h>
#include <devmand/models/wifi/Model.h>
#include <devmand/utils/StringUtils.h>

namespace devmand {
namespace devices {
namespace mikrotik {

static constexpr const unsigned short mikrotikPort = 8728;

std::shared_ptr<devices::Device> Device::createDevice(
    Application& app,
    const cartography::DeviceConfig& deviceConfig) {
  const auto& channelConfigs = deviceConfig.channelConfigs;
  const auto& otherKv = channelConfigs.at("other").kvPairs;
  const auto& snmpKv = channelConfigs.at("snmp").kvPairs;
  return std::make_unique<devices::mikrotik::Device>(
      app,
      deviceConfig.id,
      deviceConfig.readonly,
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
    bool readonly_,
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
    : snmpv2::Device(
          application,
          id_,
          readonly_,
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

std::shared_ptr<Datastore> Device::getOperationalDatastore() {
  auto state = snmpv2::Device::getOperationalDatastore();

#define GEOL(x) x["fbc-symphony-device:system"]["geo-location"]

  // fbc-symphony-device ######################################################
  // TODO check units and make conversions for all of these.
  state->addRequest(
      Mib::getLongtitude(snmpChannel)
          .thenValue([state](auto v) {
            state->update([&v](auto& lockedState) {
              GEOL(lockedState)["longitude"] = v;
            });
          })
          .thenError(
              folly::tag_t<channels::snmp::Exception>{},
              [state, this](const channels::snmp::Exception&) {
                auto v =
                    lookup("fbc-symphony-device:system/geo-location/longitude");
                if (not v.isNull()) {
                  state->update([&v](auto& lockedState) {
                    GEOL(lockedState)["longitude"] = v.asString();
                  });
                }
              }));
  state->addRequest(
      Mib::getLatitude(snmpChannel)
          .thenValue([state](auto v) {
            state->update([&v](auto& lockedState) {
              GEOL(lockedState)["latitude"] = v;
            });
          })
          .thenError(
              folly::tag_t<channels::snmp::Exception>{},
              [state, this](const channels::snmp::Exception&) {
                auto v =
                    lookup("fbc-symphony-device:system/geo-location/latitude");
                if (not v.isNull()) {
                  state->update([&v](auto& lockedState) {
                    GEOL(lockedState)["latitude"] = v.asString();
                  });
                }
              }));
  state->addRequest(
      Mib::getAltitude(snmpChannel)
          .thenValue([state](auto v) {
            state->update([&v](auto& lockedState) {
              GEOL(lockedState)["height"] = v;
            });
          })
          .thenError(
              folly::tag_t<channels::snmp::Exception>{},
              [state, this](const channels::snmp::Exception&) {
                auto v =
                    lookup("fbc-symphony-device:system/geo-location/height");
                if (not v.isNull()) {
                  state->update([&v](auto& lockedState) {
                    GEOL(lockedState)["height"] = v.asString();
                  });
                }
              }));

#undef GEOL

  // openconfig-wifi ##########################################################

#define PAP(x) x["openconfig-ap-manager:provision-aps"]["provision-ap"][0]
#define PAPC(x) PAP(x)["config"]
#define PAPT(x) PAP(x)["state"]
#define JAP(x) x["openconfig-ap-manager:joined-aps"]["joined-ap"][0]
#define JAPT(x) JAP(x)["state"]

  state->update([](auto& lockedState) {
    devmand::models::wifi::Model::init(lockedState);

    folly::dynamic& japt = JAPT(lockedState);
    // always enabled
    japt["enabled"] = true;
    // just using uptime snmp success to indicate if the device is up.
    japt["opstate"] = "openconfig-wifi-types:DOWN";
  });

  state->addRequest(Mib::getBaseMac(snmpChannel).thenValue([state](auto v) {
    auto hex = StringUtils::asHexString(v, ":");
    state->update([&hex](auto& lockedState) {
      PAP(lockedState)["mac"] = hex;
      PAPC(lockedState)["mac"] = hex;
      PAPT(lockedState)["mac"] = hex;
      JAPT(lockedState)["mac"] = hex;
    });
  }));
  state->addRequest(
      Mib::getFirmwareVersion(snmpChannel).thenValue([state](auto v) {
        state->update([&v](auto& lockedState) {
          JAPT(lockedState)["software-version"] = v;
        });
      }));
  state->addRequest(
      Mib::getSerialNumber(snmpChannel).thenValue([state](auto v) {
        state->update([&v](auto& lockedState) {
          JAPT(lockedState)["serial"] = v;
        });
      }));
  state->addRequest(Mib::getUpTime(snmpChannel).thenValue([state](auto v) {
    state->update([&v](auto& lockedState) {
      JAPT(lockedState)["uptime"] = v;
      JAPT(lockedState)["opstate"] = "openconfig-wifi-types:UP";
    });
  }));
  state->addRequest(Mib::getModel(snmpChannel).thenValue([state](auto v) {
    state->update([&v](auto& lockedState) { JAPT(lockedState)["model"] = v; });
  }));

  // TODO so this will need to be a lot more complicated as we need to combine
  // ssids that match into one in the yang model. Perhaps have a shared data
  // struct that all the requests populate into and the put it into the real
  // state in a finally at the end that also has this shared data struct.

  state->addRequest(
      snmpChannel.walk(channels::snmp::Oid(".1.3.6.1.4.1.14988.1.1.1.3.1.4"))
          .thenValue([state](auto ssids) {
            state->update([&ssids](auto& lockedState) {
              int index{0};
              for (auto& ssid : ssids) {
                devmand::models::wifi::Model::updateSsid(
                    lockedState, index, "name", ssid.value.asString());
                devmand::models::wifi::Model::updateSsid(
                    lockedState, index, "config/name", ssid.value.asString());
                devmand::models::wifi::Model::updateSsid(
                    lockedState, index, "state/name", ssid.value.asString());
                ++index;
              }
            });
          }));

  state->addRequest(
      snmpChannel.walk(channels::snmp::Oid(".1.3.6.1.4.1.14988.1.1.1.3.1.5"))
          .thenValue([state](auto bssids) {
            state->update([&bssids](auto& lockedState) {
              int index{0};
              for (auto& bssid : bssids) {
                // Mikrotik sets the bssid to an empty string. idk why...
                devmand::models::wifi::Model::updateSsidBssid(
                    lockedState, index, 0, "bssid", bssid.value.asString());
                devmand::models::wifi::Model::updateSsidBssid(
                    lockedState,
                    index,
                    0,
                    "state/bssid",
                    bssid.value.asString());
                devmand::models::wifi::Model::updateSsidBssid(
                    lockedState, index, 0, "state/radio-id", index);
                ++index;
              }
            });
          }));

  state->addRequest(
      snmpChannel.walk(channels::snmp::Oid(".1.3.6.1.4.1.14988.1.1.1.3.1.7"))
          .thenValue([state](auto freqs) {
            state->update([&freqs](auto& lockedState) {
              int index{0};
              for (auto& freq : freqs) {
                devmand::models::wifi::Model::updateSsid(
                    lockedState,
                    index,
                    "config/operating-frequency",
                    freq.value.asString());
                devmand::models::wifi::Model::updateSsid(
                    lockedState,
                    index,
                    "state/operating-frequency",
                    freq.value.asString());
                devmand::models::wifi::Model::updateRadio(
                    lockedState,
                    index,
                    "config/operating-frequency",
                    freq.value.asString());
                devmand::models::wifi::Model::updateRadio(
                    lockedState,
                    index,
                    "state/operating-frequency",
                    freq.value.asString());
                ++index;
              }
            });
          }));

  state->addFinally([state]() {
    state->update([](auto& lockedState) {
      auto hostname =
          YangUtils::lookup(lockedState, "ietf-system:system/hostname");
      if (hostname != nullptr) {
        PAPC(lockedState)
        ["hostname"] = PAPT(lockedState)["hostname"] =
            JAP(lockedState)["hostname"] = JAPT(lockedState)["hostname"] =
                hostname.asString();
      }
    });
  });

  auto venue = lookup("fbc-symphony-device:system/venue");
  if (not venue.isNull()) {
    state->update([&venue](auto& lockedState) {
      YangUtils::set(lockedState, "fbc-symphony-device:system/venue", venue);
    });
  }

  state->addRequest(Mib::getIpv4Address(snmpChannel).thenValue([state](auto v) {
    state->update([&v](auto& lockedState) { JAPT(lockedState)["ipv4"] = v; });
  }));

  state->addRequest(Mib::getIpv6Address(snmpChannel).thenValue([state](auto v) {
    state->update([&v](auto& lockedState) {
      JAPT(lockedState)
      ["ipv6"] =
          folly::IPAddress::fromBinary(
              folly::ByteRange(
                  reinterpret_cast<const unsigned char*>(v.data()), v.size()))
              .str();
    });
  }));

#undef PAP
#undef PAPC
#undef PAPT
#undef JAP
#undef JAPT

  /* TODO more ones to add
  papc["country-code"] = "US";
  papt["country-code"] = "US";
  */

  // Note we don't have a way on mikrotik to query this.
  // japt["power-source"] = "UNKNOWN";

  return state;
}

// TODO convert the device to have the concept of an intended config
void Device::setIntendedDatastore(const folly::dynamic& config) {
  auto oldInterfaces = YangUtils::lookup(
      operationalDatastore, "openconfig-interfaces:interfaces/interface");
  auto newInterfaces =
      YangUtils::lookup(config, "openconfig-interfaces:interfaces/interface");

  // TODO eh, this is very inefficient but this will go away once we switch
  // to a crud engine so its not worth making it better.
  if (newInterfaces == nullptr) {
    return;
  }

  for (auto& interface : newInterfaces) {
    folly::dynamic state;
    if (oldInterfaces != nullptr) {
      for (auto& oldInt : oldInterfaces) {
        if (interface["name"] == oldInt["name"]) {
          state = oldInt;
          break;
        }
      }
    }

    auto enabled = YangUtils::lookup(interface, "config/enabled");
    folly::dynamic adminStatus = state == nullptr
        ? YangUtils::lookup(state, "state/admin-status")
        : "DNE";

    bool isUp{false};
    if (adminStatus != nullptr and adminStatus.isString()) {
      if (adminStatus.asString() == "UP") {
        isUp = true;
      }
    }

    if (enabled != nullptr and enabled.isBool()) {
      bool isEnabled = enabled.asBool();
      if (isEnabled and not isUp) {
        mikrotikCh->writeSentence(
            {"/interface/enable", "=numbers=" + interface["name"].asString()});
        LOG(INFO) << "Interface up " << interface["name"];
      } else if (not isEnabled and isUp) {
        mikrotikCh->writeSentence(
            {"/interface/disable", "=numbers=" + interface["name"].asString()});
        LOG(INFO) << "Interface down " << interface["name"];
      }
    }
  }
}

} // namespace mikrotik
} // namespace devices
} // namespace devmand
