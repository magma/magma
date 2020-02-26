// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/cli/Cli.h>
#include <devmand/devices/cli/ParsingUtils.h>
#include <devmand/devices/cli/UbntStpPlugin.h>
#include <devmand/devices/cli/schema/ModelRegistry.h>
#include <devmand/devices/cli/translation/BindingReaderRegistry.h>
#include <devmand/devices/cli/translation/BindingWriterRegistry.h>
#include <devmand/devices/cli/translation/PluginRegistry.h>
#include <devmand/devices/cli/translation/ReaderRegistry.h>
#include <folly/executors/GlobalExecutor.h>
#include <folly/futures/Future.h>
#include <folly/json.h>
#include <ydk_openconfig/openconfig_spanning_tree.hpp>
#include <ydk_openconfig/openconfig_spanning_tree_types.hpp>
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

DeviceType UbntStpPlugin::getDeviceType() const {
  return {"ubnt", "*"};
}

class StpGlobalCfgReader : public BindingReader {
  using Stp = openconfig::openconfig_spanning_tree::Stp;

 public:
  Future<shared_ptr<Entity>> read(const Path& path, const DeviceAccess& device)
      const override {
    (void)path;
    return device.cli()
        ->executeRead(ReadCommand::create("show running-config"))
        .via(device.executor().get())
        .thenValue([](auto output) {
          (void)output;
          auto cfg = make_shared<Stp::Global::Config>();
          return cfg;
        });
  }
};

static const auto stpModeRegx =
    regex(R"(Spanning Tree Protocol[\s\.]+(stp|mstp|rstp)\s*)");

static const auto bpduGuardRegx =
    regex(R"(BPDU Guard Mode[\s\.]+(Enabled|Disabled)\s*)");

static const auto bpduFilterRegx =
    regex(R"(BPDU Filter Mode[\s\.]+(Enabled|Disabled)\s*)");

static openconfig::openconfig_spanning_tree_types::STPPROTOCOL parseStpMode(
    string stpAsString) {
  if (stpAsString == "mstp") {
    return openconfig::openconfig_spanning_tree_types::MSTP();
  }
  if (stpAsString == "stp") {
    return openconfig::openconfig_spanning_tree_types::STPPROTOCOL();
  }
  if (stpAsString == "rstp") {
    return openconfig::openconfig_spanning_tree_types::RSTP();
  }
  // default, should not happen
  return openconfig::openconfig_spanning_tree_types::STPPROTOCOL();
}

class StpGlobalStateReader : public BindingReader {
  using Stp = openconfig::openconfig_spanning_tree::Stp;

 public:
  Future<shared_ptr<Entity>> read(const Path& path, const DeviceAccess& device)
      const override {
    (void)path;

    vector<SemiFuture<string>> cmdFutures;
    cmdFutures.emplace_back(
        device.cli()->executeRead(ReadCommand::create("show spanning-tree")));
    cmdFutures.emplace_back(device.cli()->executeRead(
        ReadCommand::create("show spanning-tree summary")));

    return collect(cmdFutures.begin(), cmdFutures.end())
        .thenValue([](auto cmdOutputs) {
          string output = cmdOutputs[0];
          string summaryOutput = cmdOutputs[1];

          auto state = make_shared<Stp::Global::State>();
          parseValue(output, stpModeRegx, 1, [&state](auto stpMode) {
            state->enabled_protocol.append(parseStpMode(stpMode));
          });

          parseLeaf<bool>(
              summaryOutput,
              bpduGuardRegx,
              state->bpdu_guard,
              1,
              [](auto bpdu) { return bpdu == "Enabled"; });

          parseLeaf<bool>(
              summaryOutput,
              bpduFilterRegx,
              state->bpdu_filter,
              1,
              [](auto bpdu) { return bpdu == "Enabled"; });

          return state;
        });
  }
};

void cli::UbntStpPlugin::provideReaders(ReaderRegistryBuilder& reg) const {
  BINDING(reg, openconfigContext)
      .add(
          "/openconfig-spanning-tree:stp/global/config",
          make_shared<StpGlobalCfgReader>());
  BINDING(reg, openconfigContext)
      .add(
          "/openconfig-spanning-tree:stp/global/state",
          make_shared<StpGlobalStateReader>());
}

void cli::UbntStpPlugin::provideWriters(WriterRegistryBuilder& reg) const {
  (void)reg;
}

UbntStpPlugin::UbntStpPlugin(BindingContext& _openconfigContext)
    : openconfigContext(_openconfigContext) {}

} // namespace cli
} // namespace devices
} // namespace devmand
