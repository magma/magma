// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/cli/CliThreadWheelTimekeeper.h>
#include <devmand/channels/cli/Spd2Glog.h>
#include <devmand/channels/cli/engine/Engine.h>
#include <event2/thread.h>
#include <folly/executors/CPUThreadPoolExecutor.h>
#include <folly/executors/IOThreadPoolExecutor.h>
#include <folly/executors/thread_factory/NamedThreadFactory.h>
#include <libssh/callbacks.h>
#include <libssh/libssh.h>
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
      magma::init_logging("devmand");
    }
    magma::set_verbosity(verbosity);
    // IInitialize spd -> glog sink for YDK lib
    spdlog::create<Spd2Glog>("ydk");
    spdlog::set_level(spdlog::level::level_enum::info);
  } else {
    MLOG(MWARNING) << "Logging already initialized";
  }
}

static uint CPU_CORES = std::max(uint(4), std::thread::hardware_concurrency());

Engine::Engine()
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
      mreg(make_shared<ModelRegistry>()) {
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

} // namespace cli
} // namespace channels
} // namespace devmand
