// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <stdexcept>

#include <folly/GLog.h>

#include <devmand/channels/snmp/Engine.h>
#include <devmand/channels/snmp/Snmp.h>

namespace devmand {
namespace channels {
namespace snmp {

Engine::Engine(const std::string& appName) {
  init_snmp(appName.c_str());
}

void Engine::run() {
  while (not stopping) {
    int fds{0};
    int block{0};
    fd_set fdset;
    timeval timeout;
    timeout.tv_sec = 5;
    timeout.tv_usec = 0;

    FD_ZERO(&fdset);
    snmp_select_info(&fds, &fdset, &timeout, &block);
    fds = ::select(fds, &fdset, nullptr, nullptr, &timeout);
    if (fds < 0) {
      perror("select failed");
      throw std::runtime_error("select error");
    } else if (fds != 0) {
      snmp_read(&fdset);
    } else {
      snmp_timeout();
    }
  }
}

void Engine::stopEventually() {
  stopping = true;
}

} // namespace snmp
} // namespace channels
} // namespace devmand
