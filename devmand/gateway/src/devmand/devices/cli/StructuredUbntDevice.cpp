// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/cli/Cli.h>
#include <devmand/channels/cli/IoConfigurationBuilder.h>
#include <devmand/devices/Datastore.h>
#include <devmand/devices/cli/ParsingUtils.h>
#include <devmand/devices/cli/StructuredUbntDevice.h>
#include <devmand/devices/cli/schema/ModelRegistry.h>
#include <folly/futures/Future.h>
#include <folly/json.h>
#include <ydk_openconfig/iana_if_type.hpp>
#include <ydk_openconfig/openconfig_interfaces.hpp>
#include <ydk_openconfig/openconfig_network_instance.hpp>
#include <ydk_openconfig/openconfig_vlan_types.hpp>
#include <memory>
#include <regex>
#include <unordered_map>

namespace devmand {
namespace devices {
namespace cli {

using namespace devmand::channels::cli;
using namespace std;

using Nis = openconfig::openconfig_network_instance::NetworkInstances;
using Ni = Nis::NetworkInstance;
using Vlan = Ni::Vlans::Vlan;

using Ifcs = openconfig::openconfig_interfaces::Interfaces;
using Ifc = openconfig::openconfig_interfaces::Interfaces::Interface;
using IfcType = openconfig::iana_if_type::IanaInterfaceType;
using AdminState = Ifc::State::AdminStatus;
using OperState = Ifc::State::OperStatus;

using VlanMode = openconfig::openconfig_vlan_types::VlanModeType;

using OpenconfigInterfaces = openconfig::openconfig_interfaces::Interfaces;
using OpenconfigInterface = OpenconfigInterfaces::Interface;
using OpenconfigConfig = OpenconfigInterface::Config;
using folly::dynamic;

static WriteCommand createInterfaceCommand(string name, bool enabled) {
  string shutdownCmd = enabled ? "no shutdown" : "shutdown";
  return WriteCommand::create(
      "configure\ninterface " + name + "\n" + shutdownCmd + "\nend\n");
}

static const auto shutdown = regex("shutdown");
static const auto description = regex(R"(description '?(.+?)'?)");
static const auto mtu = regex(R"(mtu (.+))");
static const auto type = regex(R"(interface\s+(.+))");
static const auto ethernetIfc = regex(R"(\d+/\d+)");

static void parseConfig(
    Channel& channel,
    const string& ifcId,
    shared_ptr<Ifc::Config>& cfg) {
  const ReadCommand cmd =
      ReadCommand::create("show running-config interface " + ifcId);
  string output = channel.executeRead(cmd).get();

  cfg->name = ifcId;
  parseLeaf(output, mtu, cfg->mtu, 1, toUI16);
  parseLeaf<string>(output, description, cfg->description);
  cfg->enabled = true;
  parseLeaf<bool>(output, shutdown, cfg->enabled, 0, [](auto str) {
    return "shutdown" != str;
  });

  parseLeaf<IfcType>(output, type, cfg->type, 1, [&ifcId](auto str) -> IfcType {
    if (str.find("lag") == 0) {
      return openconfig::iana_if_type::Ieee8023adLag();
    } else if (str.find("vlan") == 0) {
      return openconfig::iana_if_type::L3ipvlan();
    } else if (regex_match(ifcId, ethernetIfc)) {
      return openconfig::iana_if_type::EthernetCsmacd();
    }
    return openconfig::iana_if_type::Other();
  });
}

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

static const auto counterReset =
    regex(R"(Time Since Counters Last Cleared.*?([^.]+).*)");

static void parseEthernetCounters(
    const string& output,
    shared_ptr<Ifc::State::Counters>& ctr) {
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

  parseLeaf<string>(output, counterReset, ctr->last_clear, 1);
}

static const auto mtuState = regex(R"(Max Frame Size.*?(\d+).*)");

static void parseState(
    Channel& channel,
    const string& ifcId,
    shared_ptr<Ifc::State>& state,
    const shared_ptr<Ifc::Config>& cfg) {
  const ReadCommand cmdState =
      ReadCommand::create("show interfaces description");
  string stateOut = channel.executeRead(cmdState).get();
  const auto regexIfcState = regex(ifcId + R"(\s+(\S+)\s+(\S+)\s*(.*))");

  state->name = cfg->name.get();
  state->type = cfg->type.get();
  state->enabled = cfg->enabled.get();
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

  if (cfg->type.get() == openconfig::iana_if_type::EthernetCsmacd().to_string()) {
    const ReadCommand cmdStateEth =
        ReadCommand::create("show interface ethernet " + ifcId);
    string outputStateEth = channel.executeRead(cmdStateEth).get();
    parseLeaf(outputStateEth, mtuState, state->mtu, 1, toUI16);
    parseEthernetCounters(outputStateEth, state->counters);
  }
}

static const auto vlanModeRegx = regex(R"(switchport mode (trunk|access).*)");
static const auto accessVlanRegx = regex(R"(switchport access vlan (\d+))");
static const auto trunkVlanRegx =
    regex(R"(switchport trunk allowed vlan (.+))");
static const auto vlanIdRegx = regex(R"(\d+)");

static void parseEthernet(
    Channel& channel,
    const string& ifcId,
    shared_ptr<Ifc::Ethernet>& eth) {
  const ReadCommand cmd =
      ReadCommand::create("show running-config interface " + ifcId);
  string output = channel.executeRead(cmd).get();

  parseValue(output, vlanModeRegx, 1, [&eth](auto str) {
    if (str == "trunk") {
      eth->switched_vlan->config->interface_mode = VlanMode::TRUNK;
    } else {
      eth->switched_vlan->config->interface_mode = VlanMode::ACCESS;
    }
  });

  parseValue(output, accessVlanRegx, 1, [&eth](auto str) {
    eth->switched_vlan->config->access_vlan = toUI16(str);
    // Set mode again in case the mode command was missing
    eth->switched_vlan->config->interface_mode = VlanMode::ACCESS;
  });

  parseValue(output, trunkVlanRegx, 1, [&eth](auto str) {
    // FIXME support vlan ranges, "all" and "except"
    for (auto vlan : parseLineKeys(str, vlanIdRegx, toUI16)) {
      eth->switched_vlan->config->trunk_vlans.append(vlan);
    }
    // Set mode again in case the mode command was missing
    eth->switched_vlan->config->interface_mode = VlanMode::TRUNK;
  });
}

static shared_ptr<Ifc> parseInterface(Channel& channel, const string& ifcId) {
  auto ifc = make_shared<Ifc>();
  ifc->name = ifcId;
  parseConfig(channel, ifcId, ifc->config);
  parseState(channel, ifcId, ifc->state, ifc->config);
  parseEthernet(channel, ifcId, ifc->ethernet);
  return ifc;
}

static const auto ifcIdRegex = regex(R"(^(\S+)\s+(\S+)\s+(\S+)\s*(.*))");

static shared_ptr<Ifcs> parseIfcs(Channel& channel) {
  auto cmd = ReadCommand::create("show interfaces description");
  string output = channel.executeRead(cmd).get();

  auto interfaces = make_shared<Ifcs>();
  for (auto& ifcId : parseKeys<string>(output, ifcIdRegex, 1, 4)) {
    interfaces->interface.append(parseInterface(channel, ifcId));
  }
  return interfaces;
}

static shared_ptr<Vlan> parseVlan(
    Channel& channel,
    ydk::uint16 vlanId,
    const string& vlanCfgOut,
    const string& vlanStateOut) {
  (void)channel;
  auto vlan = make_shared<Vlan>();
  vlan->vlan_id = vlanId;

  auto vlanNameRegx = regex("vlan name " + to_string(vlanId) + R"( \"(.+)\")");
  parseLeaf<string>(vlanCfgOut, vlanNameRegx, vlan->config->name);
  vlan->config->status = Vlan::Config::Status::ACTIVE;

  auto vlanStateRegx = regex(to_string(vlanId) + R"(\s+(\S+)\s+(\S+).*)");
  parseLeaf<string>(vlanStateOut, vlanStateRegx, vlan->state->name);
  vlan->state->status = Vlan::Config::Status::ACTIVE;

  return vlan;
}

static const auto vlanLineRegx = regex(R"(vlan ([\d,]+))");

static shared_ptr<Ni> parseDefaultNetwork(Channel& channel) {
  auto defaultNi = make_shared<Ni>();
  defaultNi->name = "default";

  auto cmd =
      ReadCommand::create("show running-config | section \"vlan database\"");
  string output = channel.executeRead(cmd).get();
  auto cmdState = ReadCommand::create("show vlan");
  string outputState = channel.executeRead(cmdState).get();

  auto vlanLineOpt = extractValue(output, vlanLineRegx, 1);
  if (vlanLineOpt) {
    auto configuredVlans =
        parseLineKeys(vlanLineOpt.value(), vlanIdRegx, toUI16);
    // Default vlan
    configuredVlans.push_back(1);
    for (auto vlanId : configuredVlans) {
      defaultNi->vlans->vlan.append(
          parseVlan(channel, vlanId, output, outputState));
    }
  }

  return defaultNi;
}

static shared_ptr<Nis> parseNetworks(Channel& channel) {
  (void)channel;
  auto nis = make_shared<Nis>();
  shared_ptr<Ni> defaultNi = parseDefaultNetwork(channel);
  nis->network_instance.append(defaultNi);
  return nis;
}

std::unique_ptr<devices::Device> StructuredUbntDevice::createDevice(
    Application& app,
    const cartography::DeviceConfig& deviceConfig) {
  return createDeviceWithEngine(app, deviceConfig, app.getCliEngine());
}

unique_ptr<devices::Device> StructuredUbntDevice::createDeviceWithEngine(
    Application& app,
    const cartography::DeviceConfig& deviceConfig,
    Engine& engine) {
  IoConfigurationBuilder ioConfigurationBuilder(deviceConfig, engine);
  auto cmdCache = ReadCachingCli::createCache();
  const std::shared_ptr<Channel>& channel = std::make_shared<Channel>(
      deviceConfig.id, ioConfigurationBuilder.createAll(cmdCache));

  return std::make_unique<StructuredUbntDevice>(
      app,
      deviceConfig.id,
      deviceConfig.readonly,
      channel,
      engine.getModelRegistry(),
      cmdCache);
}

StructuredUbntDevice::StructuredUbntDevice(
    Application& application,
    const Id id_,
    bool readonly_,
    const shared_ptr<Channel> _channel,
    const std::shared_ptr<ModelRegistry> _mreg,
    const shared_ptr<CliCache> _cmdCache)
    : Device(application, id_, readonly_),
      channel(_channel),
      cmdCache(_cmdCache),
      mreg(_mreg) {}

void StructuredUbntDevice::setIntendedDatastore(const dynamic& config) {
  const string& json = folly::toJson(config);
  auto& bundle = mreg->getBindingContext(Model::OPENCONFIG_2_4_3);
  const shared_ptr<OpenconfigInterfaces>& ydkModel =
      make_shared<OpenconfigInterfaces>();
  const shared_ptr<Entity> decodedIfcEntity =
      bundle.getCodec().decode(json, ydkModel);
  MLOG(MDEBUG) << decodedIfcEntity->get_segment_path();

  for (shared_ptr<Entity> entity : ydkModel->interface.entities()) {
    shared_ptr<OpenconfigInterface> iface =
        std::static_pointer_cast<OpenconfigInterface>(entity);
    shared_ptr<OpenconfigConfig> openConfig = iface->config;
    if (openConfig->type.get() !=
        openconfig::iana_if_type::EthernetCsmacd().to_string()) {
      continue;
    }
    string enabled =
        openConfig->enabled.get(); // TODO YLeaf does not support bool
    string name = openConfig->name;
    channel->executeWrite(createInterfaceCommand(name, enabled == "true"));
  }
}

shared_ptr<Datastore> StructuredUbntDevice::getOperationalDatastore() {
  MLOG(MINFO) << "[" << id << "] "
              << "Retrieving state";

  // Reset cache
  cmdCache->wlock()->clear();

  auto state = Datastore::make(*reinterpret_cast<MetricSink*>(&app), getId());
  state->setStatus(true);

  auto& bundle = mreg->getBindingContext(Model::OPENCONFIG_2_4_3);

  // TODO the conversion here is: Object -> Json -> folly:dynamic
  // the json step is unnecessary

  auto ifcs = parseIfcs(*channel);
  string json = bundle.getCodec().encode(*ifcs);
  folly::dynamic dynamicIfcs = folly::parseJson(json);

  auto networks = parseNetworks(*channel);
  json = bundle.getCodec().encode(*networks);
  folly::dynamic dynamicNis = folly::parseJson(json);

  state->update([&dynamicIfcs, &dynamicNis](folly::dynamic& lockedState) {
    lockedState.merge_patch(dynamicIfcs);
    lockedState.merge_patch(dynamicNis);
  });

  return state;
}

} // namespace cli
} // namespace devices
} // namespace devmand
