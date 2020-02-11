// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/channels/cli/CliThreadWheelTimekeeper.h>
#include <devmand/channels/cli/Spd2Glog.h>
#include <devmand/channels/cli/engine/Engine.h>
#include <devmand/devices/cli/UbntInterfacePlugin.h>
#include <devmand/devices/cli/UbntNetworksPlugin.h>
#include <devmand/devices/cli/UbntStpPlugin.h>
#include <devmand/devices/cli/translation/GrpcPlugin.h>
#include <event2/thread.h>
#include <folly/executors/CPUThreadPoolExecutor.h>
#include <folly/executors/IOThreadPoolExecutor.h>
#include <folly/executors/thread_factory/NamedThreadFactory.h>
#include <libssh/callbacks.h>
#include <libssh/libssh.h>
#include <libyang/libyang.h>
#include <spdlog/spdlog.h>
#include <iostream>

namespace devmand {
namespace channels {
namespace cli {

using devmand::channels::cli::CliThreadWheelTimekeeper;

void Engine::closeSsh() {}

void Engine::closeLogging() {
  spdlog::drop("ydk");
}

void Engine::initSsh() {
  bool f = false;
  if (sshInitialized.compare_exchange_strong(f, true)) {
    ssh_threads_set_callbacks(ssh_threads_get_pthread());
    ssh_init();
    ssh_set_log_level(SSH_LOG_NOLOG);
    evthread_use_pthreads();
  } else {
    MLOG(MWARNING) << "SSH already initialized";
  }
}

void Engine::initLogging(uint32_t verbosity, bool callInitMlog) {
  bool f = false;
  if (loggingInitialized.compare_exchange_strong(f, true)) {
    if (callInitMlog) {
      ::magma::init_logging("devmand");
    }
    ::magma::set_verbosity(verbosity);
    // Initialize spd -> glog sink for YDK lib
    spdlog::create<Spd2Glog>("ydk");
    spdlog::set_level(spdlog::level::level_enum::info);

    // Disable libyang logs
    llly_log_options(0);
  } else {
    MLOG(MWARNING) << "Logging already initialized";
  }
}

static uint CPU_CORES = std::max(uint(4), std::thread::hardware_concurrency());

typedef string PluginId;
typedef string PluginEndpoint;
static map<PluginId, PluginEndpoint> parseGrpcPlugins(dynamic pluginConfig) {
  map<PluginId, PluginEndpoint> result{};
  dynamic grpcPlugins = pluginConfig["grpcPlugins"];
  if (grpcPlugins.isArray()) {
    for (long idx = 0; idx < (long)grpcPlugins.size(); idx++) {
      try {
        dynamic plugin = *(grpcPlugins.begin() + idx);
        string id = plugin["id"].asString();
        string endpoint = plugin["endpoint"].asString();
        MLOG(MDEBUG) << "Adding grpc plugin id: " << id
                     << ", endpoint: " << endpoint;
        result[id] = endpoint;
      } catch (runtime_error& e) {
        MLOG(MWARNING) << "Cannot add plugin at index: " << idx
                       << ", error: " << e.what();
      }
    }
  } else {
    MLOG(MWARNING) << "Plugin config does not contain 'grpcPlugins' array";
  }
  return result;
}

static unique_ptr<PluginRegistry> loadPlugins(
    shared_ptr<ModelRegistry> modelRegistry,
    dynamic pluginConfig,
    shared_ptr<Executor> executor) {
  unique_ptr<PluginRegistry> pReg = make_unique<PluginRegistry>();
  // all local plugins should be added here:
  pReg->registerPlugin(make_shared<UbntInterfacePlugin>(
      modelRegistry->getBindingContext(Model::OPENCONFIG_2_4_3)));
  pReg->registerPlugin(make_shared<UbntStpPlugin>(
      modelRegistry->getBindingContext(Model::OPENCONFIG_2_4_3)));
  pReg->registerPlugin(make_shared<UbntNetworksPlugin>(
      modelRegistry->getBindingContext(Model::OPENCONFIG_2_4_3)));
  // TODO read configuration, add remote plugins
  map<PluginId, PluginEndpoint> grpcPlugins = parseGrpcPlugins(pluginConfig);
  for (auto const& kv : grpcPlugins) {
    const string& id = kv.first;
    const string& endpoint = kv.second;
    shared_ptr<grpc::Channel> grpcChannel =
        grpc::CreateChannel(endpoint, grpc::InsecureChannelCredentials());
    shared_ptr<Plugin> plugin = GrpcPlugin::create(grpcChannel, id, executor);
    pReg->registerPlugin(plugin);
  }
  return move(pReg);
}

/*
 * Keep alive cli layer needs separate executor for protecting against resource
 * starvation - connection failures should always be detected.
 * SSH cli (SshSessionAsync) - separate executor for ssh layer.
 * All other layers share common threadpool executor.
 * Plugin executor is used used for remote plugins.
 * All executors are cpu threadpool executors, to mitigate possible deadlocks,
 * such as when cli destructor calls destroy.get() - this would block forever on
 * IO threadpool executor.
 */
Engine::Engine(folly::dynamic pluginConfig)
    : channels::Engine("Cli"),
      timekeeper(make_shared<CliThreadWheelTimekeeper>()),
      sshCliExecutor(std::make_shared<folly::CPUThreadPoolExecutor>(
          CPU_CORES,
          std::make_shared<folly::NamedThreadFactory>("sshCli"))),
      commonExecutor(std::make_shared<folly::CPUThreadPoolExecutor>(
          CPU_CORES,
          std::make_shared<folly::NamedThreadFactory>("commonCli"))),
      kaCliExecutor(std::make_shared<folly::CPUThreadPoolExecutor>(
          CPU_CORES,
          std::make_shared<folly::NamedThreadFactory>("kaCli"))),
      pluginExecutor(std::make_shared<folly::CPUThreadPoolExecutor>(
          CPU_CORES,
          std::make_shared<folly::NamedThreadFactory>("plugin"))),
      mreg(make_shared<ModelRegistry>()),
      pluginRegistry(loadPlugins(mreg, pluginConfig, pluginExecutor)) {
  // TODO use singleton instead of new ThreadWheelTimekeeper when folly is
  // initialized
  Engine::initSsh();
  Engine::initLogging();
  MLOG(MINFO) << "Cli engine started with concurrency set to " << CPU_CORES;
}

Engine::~Engine() {
  Engine::closeSsh();
  Engine::closeLogging();
  MLOG(MDEBUG) << "Cli engine closed";
}

shared_ptr<CliThreadWheelTimekeeper> Engine::getTimekeeper() {
  return timekeeper;
}

shared_ptr<folly::Executor> Engine::getExecutor(
    Engine::executorRequestType requestType) const {
  if (requestType == kaCli) {
    return kaCliExecutor;
  } else if (requestType == sshCli) {
    return sshCliExecutor;
  }
  return commonExecutor;
}

shared_ptr<ModelRegistry> Engine::getModelRegistry() const {
  return mreg;
}

shared_ptr<DeviceContext> Engine::getDeviceContext(const DeviceType& type) {
  return pluginRegistry->getDeviceContext(type);
}

// TODO should this be singleton/cached?
unique_ptr<ReaderRegistry> Engine::getReaderRegistry(
    shared_ptr<DeviceContext> deviceCtx) {
  ReaderRegistryBuilder rRegBuilder{
      mreg->getSchemaContext(Model::OPENCONFIG_2_4_3)};
  deviceCtx->provideReaders(rRegBuilder);
  return rRegBuilder.build();
}

// TODO should this be singleton/cached?
unique_ptr<WriterRegistry> Engine::getWriterRegistry(
    shared_ptr<DeviceContext> deviceCtx) {
  WriterRegistryBuilder wRegBuilder{
      mreg->getSchemaContext(Model::OPENCONFIG_2_4_3)};
  deviceCtx->provideWriters(wRegBuilder);
  return wRegBuilder.build();
}

} // namespace cli
} // namespace channels
} // namespace devmand
