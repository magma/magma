// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <iostream>
#include <stdexcept>

#include <folly/Format.h>

#include <devmand/channels/cli/Channel.h>
#include <devmand/channels/cli/Cli.h>
#include <devmand/channels/cli/IoConfigurationBuilder.h>
#include <devmand/channels/cli/engine/Engine.h>
#include <devmand/devices/Datastore.h>
#include <devmand/devices/cli/PlaintextCliDevice.h>
#include <folly/executors/IOThreadPoolExecutor.h>

namespace devmand {
namespace devices {
namespace cli {

using namespace devmand::channels::cli;
using namespace devmand::channels::cli::sshsession;

std::unique_ptr<devices::Device> PlaintextCliDevice::createDevice(
    Application& app,
    const cartography::DeviceConfig& deviceConfig) {
  return createDeviceWithEngine(app, deviceConfig, app.getCliEngine());
}

std::unique_ptr<devices::Device> PlaintextCliDevice::createDeviceWithEngine(
    Application& app,
    const cartography::DeviceConfig& deviceConfig,
    Engine& engine) {
  IoConfigurationBuilder ioConfigurationBuilder(deviceConfig, engine);

  auto cmdCache = ReadCachingCli::createCache();

  const std::shared_ptr<Channel>& channel = std::make_shared<Channel>(
      deviceConfig.id, ioConfigurationBuilder.createAll(cmdCache));

  return std::make_unique<devices::cli::PlaintextCliDevice>(
      app,
      engine,
      deviceConfig.id,
      deviceConfig.channelConfigs.at("cli").kvPairs.at("stateCommand"),
      channel,
      cmdCache);
}

PlaintextCliDevice::PlaintextCliDevice(
    Application& application,
    Engine& engine,
    const Id id_,
    const std::string _stateCommand,
    const std::shared_ptr<Channel> _channel,
    const std::shared_ptr<CliCache> _cmdCache)
    : Device(application, id_, true),
      channel(_channel),
      stateCommand(ReadCommand::create(_stateCommand)),
      cmdCache(_cmdCache),
      executor(
          engine.getExecutor(Engine::executorRequestType::plaintextCliDevice)) {
}

std::shared_ptr<Datastore> PlaintextCliDevice::getOperationalDatastore() {
  MLOG(MINFO) << "[" << id << "] "
              << "Retrieving state";

  // Reset cache
  cmdCache->wlock()->clear();

  auto state = Datastore::make(*reinterpret_cast<MetricSink*>(&app), getId());

  state->addRequest(
      channel->executeRead(stateCommand)
          .via(executor.get())
          .thenValue([state, cmd = stateCommand](std::string v) {
            state->setStatus(true);
            state->update(
                [&v, &cmd](auto& lockedState) { lockedState[cmd.raw()] = v; });
          })
          .thenError(
              // TODO unify with ReconnectingCli
              // (DisconnectedException+CommandExecutionException)
              folly::tag_t<DisconnectedException>{},
              [state,
               id = this->id](DisconnectedException const& e) -> Future<Unit> {
                state->setStatus(false);
                throw e;
              })
          .thenError(
              folly::tag_t<std::exception>{},
              [state, id = this->id](std::exception const& e) {
                MLOG(MWARNING) << "[" << id << "] "
                               << "Retrieving state failed: " << e.what();
                state->addError(e.what());
              }));

  return state;
}

} // namespace cli
} // namespace devices
} // namespace devmand
