// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/cli/Spd2Glog.h>
#include <devmand/channels/cli/engine/Engine.h>
#include <libssh/callbacks.h>
#include <libssh/libssh.h>
#include <spdlog/spdlog.h>
#include <iostream>

namespace devmand {
namespace channels {
namespace cli {

void Engine::closeSsh() {
  ssh_finalize();
}

void Engine::closeLogging() {
  spdlog::drop("ydk");
}

void Engine::initSsh() {
  ssh_threads_set_callbacks(ssh_threads_get_pthread());
  ssh_init();
  ssh_set_log_level(SSH_LOG_NOLOG);
}

void Engine::initLogging(uint32_t verbosity, bool callInitMlog) {
  if (callInitMlog) {
    magma::init_logging("devmand");
  }
  magma::set_verbosity(verbosity);
  // IInitialize spd -> glog sink for YDK lib
  spdlog::create<Spd2Glog>("ydk");
  spdlog::set_level(spdlog::level::level_enum::info);
}

Engine::Engine() : channels::Engine("Cli") {
  Engine::initSsh();
  Engine::initLogging();
  MLOG(MERROR) << "Cli engine started";
}

Engine::~Engine() {
  Engine::closeSsh();
  Engine::closeLogging();
  MLOG(MDEBUG) << "Cli engine closed";
}

} // namespace cli
} // namespace channels
} // namespace devmand
