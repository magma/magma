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
#include <devmand/utils/EventBaseUtils.h>

namespace devmand {
namespace channels {
namespace snmp {

// TODO make this configurable
const constexpr std::chrono::seconds timeoutCheckInterval{1};
const constexpr std::chrono::seconds timeoutInterval{50}; // polling interval 55

Engine::Engine(folly::EventBase& eventBase_, const std::string& appName)
    : channels::Engine("SNMP"), eventBase(eventBase_) {
  init_snmp(appName.c_str());

  eventBase.runInEventBaseThread([this]() {
    EventBaseUtils::scheduleEvery(
        eventBase, [this]() { this->timeout(); }, timeoutCheckInterval);
  });
}

// This function is unused but will enable a vast amount of debugging
// information. You can limit it with the tokens you register.
void Engine::enableDebug() {
  debug_register_tokens("ALL");
  snmp_set_do_debugging(1);
}

folly::EventBase& Engine::getEventBase() {
  return eventBase;
}

void Engine::timeout() {
  // LOG(INFO) << "Processing SNMP Timeouts";
  snmp_timeout();
  sync();
}

void Engine::sync() {
  int maxfd{0};
  int block{0};
  timeval timeout{};
  timeout.tv_sec = timeoutInterval.count();
  timeout.tv_usec = 0;

  fd_set fdset{};
  FD_ZERO(&fdset);
  snmp_select_info(&maxfd, &fdset, &timeout, &block);

  // compare new vs old fds
  for (auto handler = handlers.begin(); handler != handlers.end();) {
    auto fd = (*handler)->getFd();
    if (fd >= maxfd or (not FD_ISSET(fd, &fdset))) {
      handler = handlers.erase(handler);
    } else {
      FD_CLR(fd, &fdset);
      ++handler;
    }
  }

  for (int fd = 0; fd < maxfd; ++fd) {
    if (FD_ISSET(fd, &fdset)) {
      handlers.emplace_back(std::make_unique<EventHandler>(*this, fd));
    }
  }
}

} // namespace snmp
} // namespace channels
} // namespace devmand
