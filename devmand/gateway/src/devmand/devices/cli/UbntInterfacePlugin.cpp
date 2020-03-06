// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/cli/Cli.h>
#include <devmand/devices/cli/ParsingUtils.h>
#include <devmand/devices/cli/UbntInterfacePlugin.h>
#include <devmand/devices/cli/schema/ModelRegistry.h>
#include <devmand/devices/cli/translation/BindingReaderRegistry.h>
#include <devmand/devices/cli/translation/BindingWriterRegistry.h>
#include <devmand/devices/cli/translation/PluginRegistry.h>
#include <devmand/devices/cli/translation/ReaderRegistry.h>
#include <folly/IPAddressV4.h>
#include <folly/futures/Future.h>
#include <folly/json.h>
#include <ydk_openconfig/iana_if_type.hpp>
#include <ydk_openconfig/openconfig_if_ethernet.hpp>
#include <ydk_openconfig/openconfig_if_ip.hpp>
#include <ydk_openconfig/openconfig_interfaces.hpp>
#include <ydk_openconfig/openconfig_vlan_types.hpp>
#include <memory>
#include <regex>
#include <unordered_map>

namespace devmand {
namespace devices {
namespace cli {

using namespace devmand::channels::cli;
using namespace std;
using namespace folly;
using namespace ydk;

DeviceType UbntInterfacePlugin::getDeviceType() const {
  return {"ubiquiti", "*"};
}

static const regex mtuRegx = regex(R"(mtu (.+))");
static const regex descrRegx = regex(R"(description '?(.+?)'?)");
static const regex shutRegx = regex("shutdown");
static const regex typeRegx = regex(R"(interface\s+(.+))");
static const regex ethernetIfcRegx = regex(R"(\d+/\d+)");

static string parseIfcType(string ifcName) {
  if (ifcName.find("lag") == 0) {
    return "iana-if-type:ieee8023adLag";
  } else if (ifcName.find("vlan") == 0) {
    return "iana-if-type:l3ipvlan";
  } else if (regex_match(ifcName, ethernetIfcRegx)) {
    return "iana-if-type:ethernetCsmacd";
  }
  return "iana-if-type:other";
}

using Ifc = openconfig::openconfig_interfaces::Interfaces::Interface;
using SubIfc = Ifc::Subinterfaces::Subinterface;
using IfcType = openconfig::iana_if_type::IanaInterfaceType;
using IpType = openconfig::openconfig_if_ip::IpAddressOrigin;

static IfcType parseIfcTypeBinding(string str, string ifcName) {
  if (str.find("lag") == 0) {
    return openconfig::iana_if_type::Ieee8023adLag();
  } else if (str.find("vlan") == 0) {
    return openconfig::iana_if_type::L3ipvlan();
  } else if (regex_match(ifcName, ethernetIfcRegx)) {
    return openconfig::iana_if_type::EthernetCsmacd();
  }
  return openconfig::iana_if_type::Other();
}

static Future<string> invokeRead(
    const DeviceAccess& device,
    const string& cmd) {
  return device.cli()
      ->executeRead(ReadCommand::create(cmd))
      .via(device.executor().get());
}

static Future<Unit> invokeWrite(const DeviceAccess& device, const string& cmd) {
  return device.cli()
      ->executeWrite(WriteCommand::create(cmd))
      .via(device.executor().get())
      .thenValue([](auto output) {
        (void)output;
        return unit;
      });
}

static Future<vector<string>> invokeReads(
    const DeviceAccess& device,
    const vector<string>& cmds) {
  vector<Future<string>> cmdFutures;
  for (auto& cmd : cmds) {
    cmdFutures.emplace_back(device.cli()
                                ->executeRead(ReadCommand::create(cmd))
                                .via(device.executor().get()));
  }

  return collect(cmdFutures.begin(), cmdFutures.end());
}

class IfcConfigReader : public Reader {
 public:
  Future<dynamic> read(const Path& path, const DeviceAccess& device)
      const override {
    string ifcName = path.getKeysFromSegment("interface")["name"].getString();

    return invokeRead(device, "show running-config interface " + ifcName)
        .thenValue([ifcName](auto out) {
          dynamic config = dynamic::object;
          config["name"] = ifcName;
          parseValue(out, mtuRegx, 1, [&config](auto mtuAsString) {
            config["mtu"] = toUI16(mtuAsString);
          });
          parseValue(out, descrRegx, 1, [&config](auto descrAsString) {
            config["description"] = descrAsString;
          });
          config["enabled"] = true;
          parseValue(out, shutRegx, 0, [&config](auto shutdownAsString) {
            (void)shutdownAsString;
            config["enabled"] = false;
          });
          parseValue(out, typeRegx, 1, [&config](auto typeAsString) {
            config["type"] = parseIfcType(typeAsString);
          });
          return config;
        });
  }
};

static const auto vlanModeRegx = regex(R"(switchport mode (trunk|access).*)");
static const auto accessVlanRegx = regex(R"(switchport access vlan (\d+))");
static const auto trunkVlanRegx =
    regex(R"(switchport trunk allowed vlan (.+))");
static const auto vlanIdRegx = regex(R"(\d+)");

class IfcVlanCfgReader : public BindingReader {
  using Ifc = openconfig::openconfig_interfaces::Interfaces::Interface;
  using VlanMode = openconfig::openconfig_vlan_types::VlanModeType;

 public:
  Future<shared_ptr<Entity>> read(const Path& path, const DeviceAccess& device)
      const override {
    string ifcName = path.getKeysFromSegment("interface")["name"].getString();

    return invokeRead(device, "show running-config interface " + ifcName)
        .thenValue([ifcName](auto output) {
          auto vlanCfg = make_shared<Ifc::Ethernet::SwitchedVlan::Config>();

          parseValue(output, vlanModeRegx, 1, [&vlanCfg](auto str) {
            if (str == "trunk") {
              vlanCfg->interface_mode = VlanMode::TRUNK;
            } else {
              vlanCfg->interface_mode = VlanMode::ACCESS;
            }
          });

          parseValue(output, accessVlanRegx, 1, [&vlanCfg](auto str) {
            vlanCfg->access_vlan = toUI16(str);
            // Set mode again in case the mode command was missing
            vlanCfg->interface_mode = VlanMode::ACCESS;
          });

          parseValue(output, trunkVlanRegx, 1, [&vlanCfg](auto str) {
            // FIXME support vlan ranges, "all" and "except"
            for (auto vlan : parseLineKeys(str, vlanIdRegx, toUI16)) {
              vlanCfg->trunk_vlans.append(vlan);
            }
            // Set mode again in case the mode command was missing
            vlanCfg->interface_mode = VlanMode::TRUNK;
          });
          return vlanCfg;
        });
  }
};

static bool isEthernet(string shRunOutput, string ifcName) {
  auto str = extractValue(shRunOutput, typeRegx, 1).value_or("");
  if (parseIfcTypeBinding(str, ifcName).to_string() ==
      openconfig::iana_if_type::EthernetCsmacd().to_string()) {
    return true;
  }

  return false;
}

class IfcPoeCfgReader : public BindingReader {
  using Ifc = openconfig::openconfig_interfaces::Interfaces::Interface;

 public:
  Future<shared_ptr<Entity>> read(const Path& path, const DeviceAccess& device)
      const override {
    string ifcName = path.getKeysFromSegment("interface")["name"].getString();

    return invokeRead(device, "show running-config interface " + ifcName)
        .thenValue([ifcName](auto output) {
          // Only ethernet ifcs can have poe
          if (!isEthernet(output, ifcName)) {
            return make_shared<Ifc::Ethernet::Poe::Config>();
          }

          auto poeCfg = make_shared<Ifc::Ethernet::Poe::Config>();
          poeCfg->enabled = output.find("poe opmode shutdown") == string::npos;
          return poeCfg;
        });
  }
};

static const auto mtuState = regex(R"(Max Frame Size.*?(\d+).*)");

class IfcStateReader : public BindingReader {
  using Ifc = openconfig::openconfig_interfaces::Interfaces::Interface;
  using AdminState = Ifc::State::AdminStatus;
  using OperState = Ifc::State::OperStatus;

 public:
  Future<shared_ptr<Entity>> read(const Path& path, const DeviceAccess& device)
      const override {
    string ifcName = path.getKeysFromSegment("interface")["name"].getString();

    return invokeReads(
               device,
               {"show running-config interface " + ifcName,
                "show interfaces description"})
        .thenValue([ifcName](auto cmdOutputs) {
          string output = cmdOutputs[0];
          string stateOut = cmdOutputs[1];

          auto state = make_shared<Ifc::State>();
          state->name = ifcName;

          state->enabled = true;
          parseLeaf<bool>(output, shutRegx, state->enabled, 0, [](auto str) {
            (void)str;
            return false;
          });

          parseLeaf<IfcType>(
              output,
              typeRegx,
              state->type,
              1,
              [&ifcName](auto str) -> IfcType {
                return parseIfcTypeBinding(str, ifcName);
              });

          auto regexIfcState = regex(ifcName + R"(\s+(\S+)\s+(\S+)\s*(.*))");
          parseLeaf<string>(stateOut, regexIfcState, state->description, 3);

          parseValue(stateOut, regexIfcState, 1, [&state](string value) {
            if ("Enable" == value) {
              state->admin_status = AdminState::UP;
            } else if ("Disable" == value) {
              state->admin_status = AdminState::DOWN;
            }
          });

          parseValue(stateOut, regexIfcState, 2, [&state](string value) {
            if ("Up" == value) {
              state->oper_status = OperState::UP;
            } else if ("Down" == value) {
              state->oper_status = OperState::DOWN;
            } else {
              state->oper_status = OperState::UNKNOWN;
            }
          });
          return state;
        });
  }
};

static const auto inOct =
    regex(R"(Total Packets Received \(Octets\).*?(\d+).*)");
static const auto inUPkt = regex(R"(Unicast Packets Received.*?(\d+).*)");
static const auto inMPkt = regex(R"(Multicast Packets Received.*?(\d+).*)");
static const auto inBPkt = regex(R"(Broadcast Packets Received.*?(\d+).*)");
static const auto inDisc = regex(R"(Receive Packets Discarded.*?(\d+).*)");
static const auto inErrors =
    regex(R"(Total Packets Received with MAC Errors.*?(\d+).*)");

static const auto outOct =
    regex(R"(Total Packets Transmitted \(Octets\).*?(\d+).*)");
static const auto outUPkt = regex(R"(Unicast Packets Transmitted.*?(\d+).*)");
static const auto outMPkt = regex(R"(Multicast Packets Transmitted.*?(\d+).*)");
static const auto outBPkt = regex(R"(Broadcast Packets Transmitted.*?(\d+).*)");
static const auto outDisc = regex(R"(Transmit Packets Discarded.*?(\d+).*)");
static const auto outErrors = regex(R"(Total Transmit Errors.*?(\d+).*)");

// static const auto counterReset =
//    regex(R"(Time Since Counters Last Cleared.*?([^.]+).*)");

class IfcCountersReader : public BindingReader {
  using Ifc = openconfig::openconfig_interfaces::Interfaces::Interface;

 public:
  Future<shared_ptr<Entity>> read(const Path& path, const DeviceAccess& device)
      const override {
    string ifcName = path.getKeysFromSegment("interface")["name"].getString();

    return invokeReads(
               device,
               {"show running-config interface " + ifcName,
                "show interface ethernet " + ifcName})
        .thenValue([ifcName](auto cmdOutputs) -> Future<shared_ptr<Entity>> {
          auto shRunOutput = cmdOutputs[0];
          auto output = cmdOutputs[1];

          // If is not ethernet interface, do not read anything and return empty
          if (!isEthernet(shRunOutput, ifcName)) {
            return make_shared<Ifc::State::Counters>();
          }

          auto ctr = make_shared<Ifc::State::Counters>();

          parseLeaf(output, inOct, ctr->in_octets, 1, toUI64);
          parseLeaf(output, inUPkt, ctr->in_unicast_pkts, 1, toUI64);
          parseLeaf(output, inMPkt, ctr->in_multicast_pkts, 1, toUI64);
          parseLeaf(output, inBPkt, ctr->in_broadcast_pkts, 1, toUI64);
          parseLeaf(output, inDisc, ctr->in_discards, 1, toUI64);
          parseLeaf(output, inErrors, ctr->in_errors, 1, toUI64);

          parseLeaf(output, outOct, ctr->out_octets, 1, toUI64);
          parseLeaf(output, outUPkt, ctr->out_unicast_pkts, 1, toUI64);
          parseLeaf(output, outMPkt, ctr->out_multicast_pkts, 1, toUI64);
          parseLeaf(output, outBPkt, ctr->out_broadcast_pkts, 1, toUI64);
          parseLeaf(output, outDisc, ctr->out_discards, 1, toUI64);
          parseLeaf(output, outErrors, ctr->out_errors, 1, toUI64);

          //          TODO fix format of last clear
          //          parseLeaf<string>(output, counterReset, ctr->last_clear,
          //          1);

          return ctr;
        });
  }
};

void cli::UbntInterfacePlugin::provideReaders(
    ReaderRegistryBuilder& reg) const {
  // List reader is a DOM lambda
  reg.addList(
      "/openconfig-interfaces:interfaces/interface",
      [](const Path& path, const DeviceAccess& device) {
        (void)path;
        return invokeRead(device, "show interfaces description")
            .thenValue([](auto output) -> Future<vector<dynamic>> {
              return parseKeys<dynamic>(
                  output,
                  regex(R"(^(\S+)\s+(\S+)\s+(\S+)\s*(.*))"),
                  1,
                  4,
                  [](auto ifcName) {
                    return dynamic::object("name", ifcName);
                  });
            });
      });

  // Reader is a DOM class
  reg.add(
      "/openconfig-interfaces:interfaces/interface/config",
      make_shared<IfcConfigReader>());

  // Reader is a Binding-aware class
  BINDING(reg, openconfigContext)
      .add(
          "/openconfig-interfaces:interfaces/interface/state",
          make_shared<IfcStateReader>());

  // Reader is a Binding-aware class
  BINDING(reg, openconfigContext)
      .add(
          "/openconfig-interfaces:interfaces/interface/state/counters",
          make_shared<IfcCountersReader>());

  BINDING(reg, openconfigContext)
      .addList(
          "/openconfig-interfaces:interfaces/interface/subinterfaces/subinterface",
          [](const Path& path, const DeviceAccess& device) {
            (void)path;
            (void)device;
            vector<EntityKeys> allKeys;
            SubIfc subifc;
            subifc.index_ = 0; // add just logical 0 subinterface
            allKeys.push_back({subifc.index_});
            return allKeys;
          });

  BINDING(reg, openconfigContext)
      .add(
          "/openconfig-interfaces:interfaces/interface/subinterfaces/subinterface/config",
          [](const Path& path, const DeviceAccess& device) {
            (void)path;
            (void)device;
            auto cfg = make_shared<SubIfc::Config>();
            cfg->index_ =
                path.getKeysFromSegment("subinterface")["index"].asString();
            return cfg;
          });

  BINDING(reg, openconfigContext)
      .addList(
          "/openconfig-interfaces:interfaces/interface/subinterfaces/subinterface/openconfig-if-ip:ipv4/addresses/address",
          [](const Path& path, const DeviceAccess& device) {
            string ifcName =
                path.getKeysFromSegment("interface")["name"].getString();

            return invokeRead(
                       device, "show running-config interface " + ifcName)
                .thenValue([ifcName](auto output) {
                  return parseKeys<EntityKeys>(
                      output,
                      regex(R"(ip address ([^\s]+) ([^\s]+))"),
                      1,
                      1,
                      [](auto ip) -> EntityKeys {
                        auto tmp =
                            make_shared<SubIfc::Ipv4::Addresses::Address>();
                        tmp->ip = ip;
                        return {tmp->ip};
                      });
                });
          });

  BINDING(reg, openconfigContext)
      .add(
          "/openconfig-interfaces:interfaces/interface/subinterfaces/subinterface/openconfig-if-ip:ipv4/addresses/address/config",
          [](const Path& path, const DeviceAccess& device) {
            string ifcName =
                path.getKeysFromSegment("interface")["name"].getString();
            string ip = path.getKeysFromSegment("address")["ip"].getString();

            return invokeRead(
                       device, "show running-config interface " + ifcName)
                .thenValue(
                    [ifcName, ip](auto output) -> Future<shared_ptr<Entity>> {
                      auto cfg = make_shared<
                          SubIfc::Ipv4::Addresses::Address::Config>();
                      cfg->ip = ip;
                      parseLeaf<int>(
                          output,
                          regex("ip address " + ip + " ([^\\s]+)"),
                          cfg->prefix_length,
                          1,
                          [](string mask) {
                            return bitset<32>(IPAddressV4::toLong(mask))
                                .count();
                          });
                      return cfg;
                    });
          });
  BINDING(reg, openconfigContext)
      .add(
          "/openconfig-interfaces:interfaces/interface/subinterfaces/subinterface/openconfig-if-ip:ipv4/addresses/address/state",
          [](const Path& path, const DeviceAccess& device) {
            string ifcName =
                path.getKeysFromSegment("interface")["name"].getString();
            string ip = path.getKeysFromSegment("address")["ip"].getString();

            return invokeRead(
                       device, "show running-config interface " + ifcName)
                .thenValue(
                    [ifcName, ip](auto output) -> Future<shared_ptr<Entity>> {
                      auto state = make_shared<
                          SubIfc::Ipv4::Addresses::Address::State>();
                      state->ip = ip;
                      // Only static ips are supported as of now
                      state->origin = IpType::STATIC;
                      parseLeaf<int>(
                          output,
                          regex("ip address " + ip + " ([^\\s]+)"),
                          state->prefix_length,
                          1,
                          [](string mask) {
                            return bitset<32>(IPAddressV4::toLong(mask))
                                .count();
                          });
                      return state;
                    });
          });

  // Reader is a Binding-aware class
  BINDING(reg, openconfigContext)
      .add(
          "/openconfig-interfaces:interfaces/interface/openconfig-if-ethernet:ethernet/openconfig-vlan:switched-vlan/config",
          make_shared<IfcVlanCfgReader>());

  // Reader is a Binding-aware lambda
  BINDING(reg, openconfigContext)
      .add(
          "/openconfig-interfaces:interfaces/interface/openconfig-if-ethernet:ethernet/openconfig-if-poe:poe/config",
          make_shared<IfcPoeCfgReader>());
}

class InterfaceConfigWriter : public BindingWriter<Ifc::Config> {
 public:
  Future<Unit> create(
      const Path& path,
      shared_ptr<Ifc::Config> cfg,
      const DeviceAccess& device) const override {
    (void)device;
    (void)cfg;
    return runtime_error(
        "Interface creation is not supported yet. Called for: " + path.str());
  }

  Future<Unit> update(
      const Path& path,
      shared_ptr<Ifc::Config> before,
      shared_ptr<Ifc::Config> after,
      const DeviceAccess& device) const override {
    (void)path;
    (void)before;
    string shutdownCmd =
        after->enabled.value == "true" ? "no shutdown" : "shutdown";
    return invokeWrite(
        device,
        "configure\ninterface " + after->name.value + "\n" + shutdownCmd +
            "\nend\n");
  }

  Future<Unit> remove(
      const Path& path,
      shared_ptr<Ifc::Config> before,
      const DeviceAccess& device) const override {
    (void)device;
    (void)before;
    return runtime_error(
        "Interface deletion is not supported yet. Called for: " + path.str());
  }
};

void cli::UbntInterfacePlugin::provideWriters(
    WriterRegistryBuilder& reg) const {
  // Writer is a binding class
  BINDING_W(reg, openconfigContext)
      .add(
          //          "/openconfig-interfaces:interfaces/interface/openconfig-interfaces:config",
          "/openconfig-interfaces:interfaces/openconfig-interfaces:interface/openconfig-interfaces:config",
          make_shared<InterfaceConfigWriter>());
}

UbntInterfacePlugin::UbntInterfacePlugin(BindingContext& _openconfigContext)
    : openconfigContext(_openconfigContext) {}

} // namespace cli
} // namespace devices
} // namespace devmand
