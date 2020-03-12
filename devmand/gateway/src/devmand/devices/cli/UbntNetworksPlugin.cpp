// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/cli/Cli.h>
#include <devmand/devices/cli/ParsingUtils.h>
#include <devmand/devices/cli/UbntNetworksPlugin.h>
#include <devmand/devices/cli/schema/ModelRegistry.h>
#include <devmand/devices/cli/translation/BindingReaderRegistry.h>
#include <devmand/devices/cli/translation/BindingWriterRegistry.h>
#include <devmand/devices/cli/translation/PluginRegistry.h>
#include <devmand/devices/cli/translation/ReaderRegistry.h>
#include <folly/executors/GlobalExecutor.h>
#include <folly/futures/Future.h>
#include <folly/json.h>
#include <ydk_openconfig/openconfig_network_instance.hpp>
#include <ydk_openconfig/openconfig_network_instance_types.hpp>
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

DeviceType UbntNetworksPlugin::getDeviceType() const {
  return {"ubiquiti", "*"};
}

static Future<string> invokeRead(
    const DeviceAccess& device,
    const string& cmd) {
  return device.cli()
      ->executeRead(ReadCommand::create(cmd))
      .via(device.executor().get());
}

static const auto vlanLineRegx = regex(R"(vlan ([\d,]+))");
static const auto vlanIdRegx = regex(R"(\d+)");

void cli::UbntNetworksPlugin::provideReaders(ReaderRegistryBuilder& reg) const {
  // List reader is a Binding-aware lambda
  using Nis = openconfig::openconfig_network_instance::NetworkInstances;
  BINDING(reg, openconfigContext)
      .addList(
          Path(
              "/openconfig-network-instance:network-instances/network-instance"),
          [](const Path& path, const DeviceAccess& device) {
            (void)path;
            return invokeRead(device, "show running-config")
                .thenValue([](auto output) -> Future<vector<EntityKeys>> {
                  (void)output;
                  vector<EntityKeys> allKeys;
                  Nis::NetworkInstance instance;
                  instance.name = "default";
                  allKeys.push_back({instance.name});
                  return allKeys;
                });
          });

  //   Reader is a Binding-aware lambda
  using DEFAULT_NI =
      openconfig::openconfig_network_instance_types::DEFAULTINSTANCE;
  BINDING(reg, openconfigContext)
      .add(
          "/openconfig-network-instance:network-instances/network-instance/config",
          [](const Path& path,
             const DeviceAccess& device) -> Future<shared_ptr<Entity>> {
            (void)path;
            (void)device;
            auto cfg = make_shared<Nis::NetworkInstance::Config>();
            cfg->name =
                path.getKeysFromSegment("network-instance")["name"].asString();
            cfg->type = DEFAULT_NI();
            return cfg;
          });

  // List reader is a Binding-aware lambda
  using Vlan = Nis::NetworkInstance::Vlans::Vlan;
  BINDING(reg, openconfigContext)
      .addList(
          Path(
              "/openconfig-network-instance:network-instances/network-instance/vlans/vlan"),
          [](const Path& path, const DeviceAccess& device) {
            (void)path;
            return invokeRead(
                       device,
                       "show running-config | section \"vlan database\"")
                .thenValue([](auto output) -> Future<vector<EntityKeys>> {
                  vector<EntityKeys> allKeys;
                  auto vlanLine =
                      extractValue(output, vlanLineRegx, 1).value_or("");
                  vlanLine = "1," + vlanLine; // add default vlan
                  return parseLineKeys<EntityKeys>(
                      vlanLine, vlanIdRegx, [](string vlanId) -> EntityKeys {
                        Vlan instance;
                        instance.vlan_id = toUI16(vlanId);
                        return {instance.vlan_id};
                      });
                });
          });

  BINDING(reg, openconfigContext)
      .add(
          "/openconfig-network-instance:network-instances/network-instance/vlans/vlan/config",
          [](const Path& path,
             const DeviceAccess& device) -> Future<shared_ptr<Entity>> {
            return invokeRead(
                       device,
                       "show running-config | section \"vlan database\"")
                .thenValue([path](auto output) {
                  auto cfg = make_shared<Vlan::Config>();

                  string vlanId =
                      path.getKeysFromSegment("vlan")["vlan-id"].asString();
                  cfg->vlan_id = toUI16(vlanId);
                  cfg->status = Vlan::Config::Status::ACTIVE;

                  auto vlanNameRegx =
                      regex("vlan name " + vlanId + R"( \"(.+)\")");
                  parseLeaf<string>(output, vlanNameRegx, cfg->name);
                  return cfg;
                });
          });

  BINDING(reg, openconfigContext)
      .add(
          "/openconfig-network-instance:network-instances/network-instance/vlans/vlan/state",
          [](const Path& path,
             const DeviceAccess& device) -> Future<shared_ptr<Entity>> {
            return invokeRead(device, "show vlan")
                .thenValue([path](auto output) {
                  auto state = make_shared<Vlan::State>();

                  string vlanId =
                      path.getKeysFromSegment("vlan")["vlan-id"].asString();
                  state->vlan_id = toUI16(vlanId);
                  state->status = Vlan::Config::Status::ACTIVE;
                  auto vlanStateRegx = regex(vlanId + R"(\s+(\S+)\s+(\S+).*)");
                  parseLeaf<string>(output, vlanStateRegx, state->name);
                  return state;
                });
          });
}

void cli::UbntNetworksPlugin::provideWriters(WriterRegistryBuilder& reg) const {
  (void)reg;
}

UbntNetworksPlugin::UbntNetworksPlugin(BindingContext& _openconfigContext)
    : openconfigContext(_openconfigContext) {}

} // namespace cli
} // namespace devices
} // namespace devmand
