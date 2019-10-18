// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/snmp/Engine.h>
#include <devmand/channels/snmp/EventHandler.h>
#include <devmand/channels/snmp/Snmp.h>

namespace devmand {
namespace channels {
namespace snmp {

EventHandler::EventHandler(Engine& engine_, int fd_)
    : folly::EventHandler(&engine_.getEventBase()), engine(engine_), fd(fd_) {
  folly::EventHandler::changeHandlerFD(folly::NetworkSocket::fromFd(fd));
  registerHandler(folly::EventHandler::READ | folly::EventHandler::PERSIST);
}

EventHandler::~EventHandler() {
  if (fd != -1) {
    unregisterHandler();
  }
}

int EventHandler::getFd() const {
  return fd;
}

void EventHandler::handlerReady(uint16_t) noexcept {
  fd_set fdset;
  FD_ZERO(&fdset);
  FD_SET(fd, &fdset);
  snmp_read(&fdset);
  engine.sync();
}

} // namespace snmp
} // namespace channels
} // namespace devmand
