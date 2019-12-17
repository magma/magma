// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <iostream>
#include <stdexcept>

#include <folly/Format.h>

#include <devmand/Application.h>
#include <devmand/channels/snmp/IfMib.h>
#include <devmand/devices/Datastore.h>
#include <devmand/devices/snmpv2/Device.h>
#include <devmand/error/ErrorHandler.h>
#include <devmand/models/interface/Model.h>

namespace devmand {
namespace devices {
namespace snmpv2 {

// TODO autogen this from:
//   https://www.iana.org/assignments/ianaiftype-mib/ianaiftype-mib
#define IANAIFTYPE_OTHER 1
#define IANAIFTYPE_ETHERNETCSMACD 6
#define IANAIFTYPE_SOFTWARELOOPBACK 24
#define IANAIFTYPE_BRIDGE 209

static bool isLoopBack(int ifMibType) {
  switch (ifMibType) {
    case IANAIFTYPE_SOFTWARELOOPBACK:
      return true;
    default:
      return false;
  }
}

static std::string getTypeString(int ifMibType) {
  std::string ts;
  switch (ifMibType) {
    case IANAIFTYPE_OTHER:
      ts = "other";
      break;
    case IANAIFTYPE_ETHERNETCSMACD:
      ts = "ethernetCsmacd";
      break;
    case IANAIFTYPE_SOFTWARELOOPBACK:
      ts = "softwareLoopback";
      break;
    case IANAIFTYPE_BRIDGE:
      ts = "bridge";
      break;
    default:
      ts = "unknown";
      break;
  }
  // TODO confirm if this is the true valid yang extension
  return folly::sformat("iana-if-type:{}", ts);
}

std::shared_ptr<devices::Device> Device::createDevice(
    Application& app,
    const cartography::DeviceConfig& deviceConfig) {
  const auto& channelConfigs = deviceConfig.channelConfigs;
  const auto& snmpKv = channelConfigs.at("snmp").kvPairs;
  return std::make_unique<devices::snmpv2::Device>(
      app,
      deviceConfig.id,
      deviceConfig.readonly,
      deviceConfig.ip,
      snmpKv.at("community"),
      snmpKv.at("version"));
}

Device::Device(
    Application& application,
    const Id& id_,
    bool readonly_,
    const channels::snmp::Peer& peer,
    const channels::snmp::Community& community,
    const channels::snmp::Version& version,
    const std::string& passphrase,
    const std::string& securityName,
    const channels::snmp::SecurityLevel& securityLevel,
    oid proto[])
    : ping::Device(application, id_, readonly_, folly::IPAddress(peer)),
      snmpChannel(
          application.getSnmpEngine(),
          peer,
          community,
          version,
          passphrase,
          securityName,
          securityLevel,
          proto) {}

std::shared_ptr<Datastore> Device::getOperationalDatastore() {
  using IfMib = devmand::channels::snmp::IfMib;
  using IModel = devmand::models::interface::Model;

  std::shared_ptr<Datastore> state = ping::Device::getOperationalDatastore();

  state->update([](auto& lockedState) {
    IModel::init(lockedState);
    lockedState["ietf-system:system"] = folly::dynamic::object;
  });

  state->addRequest(
      IfMib::getSystemName(snmpChannel).thenValue([state](auto v) {
        state->update([&v](auto& lockedState) {
          lockedState["ietf-system:system"]["hostname"] = v;
        });
      }));

  state->addRequest(
      IfMib::getSystemContact(snmpChannel).thenValue([state](auto v) {
        state->update([&v](auto& lockedState) {
          lockedState["ietf-system:system"]["contact"] = v;
        });
      }));
  state->addRequest(
      IfMib::getSystemLocation(snmpChannel).thenValue([state](auto v) {
        state->update([&v](auto& lockedState) {
          lockedState["ietf-system:system"]["location"] = v;
        });
      }));
  std::weak_ptr<devices::Device> weak(this->shared_from_this());
  // TODO: state::addRequest is for individual requests. We use it here so that
  // State::collect doesn't miss any requests, since they're nested. Figure out
  // a way to track individual requests again.
  state->addRequest(
      IfMib::getInterfaceIndicies(snmpChannel)
          // "this" pointer included in lambda so that it has snmpChannel in
          // scope. Don't use "this" without first locking on the weak pointer
          // so that you make sure the object still exists.
          .thenValue([this, weak, state](auto interfaceIndices) {
            auto sharedDevice = weak.lock();
            if (not sharedDevice) {
              LOG(INFO) << "Can't get state with destructed device object.";
              return folly::Future<folly::Unit>();
            }
            return this->addToStateWithInterfaceIndices(
                state, interfaceIndices);
          }));

  return state;
}

folly::Future<folly::Unit> Device::addToStateWithInterfaceIndices(
    std::shared_ptr<Datastore> state,
    const devmand::channels::snmp::InterfaceIndicies& interfaceIndices) {
  using IfMib = devmand::channels::snmp::IfMib;
  using IModel = devmand::models::interface::Model;

  std::vector<folly::Future<folly::Unit>> allFutures;
  allFutures.push_back(
      IfMib::getInterfaceNames(snmpChannel, interfaceIndices)
          .thenValue([state](auto results) {
            state->update([&results](auto& lockedState) {
              for (auto result : results) {
                IModel::updateInterface(
                    lockedState, result.index, "name", result.value);
                IModel::updateInterface(
                    lockedState, result.index, "state/name", result.value);
                IModel::updateInterface(
                    lockedState, result.index, "config/name", result.value);
              }
            });
          }));
  allFutures.push_back(
      IfMib::getInterfaceOperStatuses(snmpChannel, interfaceIndices)
          .thenValue([state](auto results) {
            state->update([&results](auto& lockedState) {
              for (auto result : results) {
                // TODO this is not valid according to the model but
                // we need to fix the front-end.
                IModel::updateInterface(
                    lockedState, result.index, "oper-status", result.value);
                IModel::updateInterface(
                    lockedState,
                    result.index,
                    "state/oper-status",
                    result.value);
              }
            });
          }));
  allFutures.push_back(
      IfMib::getInterfaceAdminStatuses(snmpChannel, interfaceIndices)
          .thenValue([state](auto results) {
            state->update([&results](auto& lockedState) {
              for (auto result : results) {
                bool isEnabled = result.value == "UP";
                IModel::updateInterface(
                    lockedState,
                    result.index,
                    "state/admin-status",
                    result.value);
                IModel::updateInterface(
                    lockedState, result.index, "config/enabled", isEnabled);
                IModel::updateInterface(
                    lockedState, result.index, "state/enabled", isEnabled);
              }
            });
          }));
  allFutures.push_back(
      IfMib::getInterfaceMtus(snmpChannel, interfaceIndices)
          .thenValue([state](auto results) {
            state->update([&results](auto& lockedState) {
              for (auto result : results) {
                IModel::updateInterface(
                    lockedState, result.index, "config/mtu", result.value);
                IModel::updateInterface(
                    lockedState, result.index, "state/mtu", result.value);
              }
            });
          }));
  allFutures.push_back(
      IfMib::getInterfaceTypes(snmpChannel, interfaceIndices)
          .thenValue([state](auto results) {
            state->update([&results](auto& lockedState) {
              for (auto result : results) {
                int ifMibType = folly::to<int>(result.value);
                std::string interfaceType = getTypeString(ifMibType);
                IModel::updateInterface(
                    lockedState, result.index, "config/type", interfaceType);
                IModel::updateInterface(
                    lockedState, result.index, "state/type", interfaceType);

                bool loopbackMode = isLoopBack(ifMibType);
                IModel::updateInterface(
                    lockedState,
                    result.index,
                    "config/loopback-mode",
                    loopbackMode);
                IModel::updateInterface(
                    lockedState,
                    result.index,
                    "state/loopback-mode",
                    loopbackMode);
              }
            });
          }));

  allFutures.push_back(
      IfMib::getInterfaceDescriptions(snmpChannel, interfaceIndices)
          .thenValue([state](auto results) {
            state->update([&results](auto& lockedState) {
              for (auto result : results) {
                IModel::updateInterface(
                    lockedState,
                    result.index,
                    "config/description",
                    result.value);
                IModel::updateInterface(
                    lockedState,
                    result.index,
                    "state/description",
                    result.value);
              }
            });
          }));

  allFutures.push_back(
      IfMib::getInterfaceLastChange(snmpChannel, interfaceIndices)
          .thenValue([state](auto results) {
            state->update([&results](auto& lockedState) {
              for (auto result : results) {
                IModel::updateInterface(
                    lockedState,
                    result.index,
                    "state/last-change",
                    result.value);
              }
            });
          }));

  auto addRequest = [this, &state, &allFutures, &interfaceIndices](
                        const std::string& oid, const std::string& path) {
    allFutures.push_back(
        IfMib::getInterfaceField(snmpChannel, interfaceIndices, oid)
            .thenValue([state, path](auto results) {
              state->update([&results, &state, &path](auto& lockedState) {
                for (auto result : results) {
                  IModel::updateInterface(
                      lockedState, result.index, path, result.value);
                  // TODO: instead of doing this per device type, move to
                  //   traversing the resulting device model and
                  //   creating metrics in a more general fashion
                  state->setGauge(
                      folly::sformat(
                          "/{}/interface[ifindex={}]/{}",
                          "openconfig-interfaces:interface",
                          result.index,
                          path),
                      folly::to<float>(result.value));
                }
              });
            }));
  };

  // TODO if devices don't support the 64 bit version should we revert to 32?
  addRequest(".1.3.6.1.2.1.31.1.1.1.6.", "state/counters/in-octets");
  addRequest(".1.3.6.1.2.1.31.1.1.1.7.", "state/counters/in-unicast-pkts");
  addRequest(".1.3.6.1.2.1.31.1.1.1.9.", "state/counters/in-broadcast-pkts");
  addRequest(".1.3.6.1.2.1.31.1.1.1.8.", "state/counters/in-multicast-pkts");
  addRequest(".1.3.6.1.2.1.2.2.1.13.", "state/counters/in-discards");
  addRequest(".1.3.6.1.2.1.2.2.1.14.", "state/counters/in-errors");
  addRequest(".1.3.6.1.2.1.2.2.1.15.", "state/counters/in-unknown-protos");

  addRequest(".1.3.6.1.2.1.31.1.1.1.6.", "state/counters/out-octets");
  addRequest(".1.3.6.1.2.1.31.1.1.1.11.", "state/counters/out-unicast-pkts");
  addRequest(".1.3.6.1.2.1.31.1.1.1.13.", "state/counters/out-broadcast-pkts");
  addRequest(".1.3.6.1.2.1.31.1.1.1.12.", "state/counters/out-multicast-pkts");
  addRequest(".1.3.6.1.2.1.2.2.1.19.", "state/counters/out-discards");
  addRequest(".1.3.6.1.2.1.2.2.1.20.", "state/counters/out-errors");

  // TODO how should I get state/logical? Perhaps based on type?

  /*
   *  TODO here are some fields I still don't know how to get. The in/out pkts
   *  could come from summing perhaps.
      |     +--ro in-pkts?               oc-yang:counter64
      |     +--ro in-fcs-errors?         oc-yang:counter64
      |     +--ro out-pkts?              oc-yang:counter64
      |     +--ro carrier-transitions?   oc-yang:counter64
      |     +--ro last-clear?            oc-types:timeticks64
   */
  return folly::collect(std::move(allFutures)).unit();
}

} // namespace snmpv2
} // namespace devices
} // namespace devmand
