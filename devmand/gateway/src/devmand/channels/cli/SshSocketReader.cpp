// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/cli/SshSessionAsync.h>
#include <devmand/channels/cli/SshSocketReader.h>
#include <event2/event.h>
#include <chrono>
#include <thread>

using devmand::channels::cli::sshsession::SshSessionAsync;

void sshReadNotificationThread(struct event_base* base);
void sshReadNotificationThread(struct event_base* base) {
  while (true) {
    int rv = event_base_dispatch(base);

    if (event_base_got_exit(base)) {
      MLOG(MDEBUG) << "event_base_exit called, terminating";
      return;
    }
    if (rv == 1) {
      std::this_thread::sleep_for(std::chrono::seconds(1));
      continue;
    } else {
      const char* error =
          "event_base_dispatch catastrophic failure, terminating";
      MLOG(MERROR) << error << " returned value: " << rv;
      throw std::runtime_error(error);
    }
  }
}

devmand::channels::cli::SshSocketReader::SshSocketReader()
    : base(event_base_new()),
      notificationThread(std::thread(sshReadNotificationThread, base)) {}

struct event* devmand::channels::cli::SshSocketReader::addSshReader(
    event_callback_fn callbackFn,
    socket_t fd,
    void* ptr) {
  struct event* event_on_heap =
      event_new(this->base, fd, EV_READ | EV_PERSIST, callbackFn, ptr);
  event_add(event_on_heap, nullptr);
  return event_on_heap;
}

devmand::channels::cli::SshSocketReader::~SshSocketReader() {
  event_base_loopexit(base, nullptr);
  notificationThread.join();
  event_base_free(base);
}

devmand::channels::cli::SshSocketReader&
devmand::channels::cli::SshSocketReader::getInstance() {
  static SshSocketReader instance;
  return instance;
}
