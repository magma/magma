// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#define LOG_WITH_GLOG

#include <devmand/channels/cli/SshSessionAsync.h>
#include <magma_logging.h>

using devmand::channels::cli::sshsession::SshSessionAsync;

namespace devmand {
namespace channels {
namespace cli {

class SshSocketReader {
 private:
  struct event_base* base;
  std::thread notificationThread;

 public:
  static SshSocketReader& getInstance(); // singleton
  SshSocketReader();
  SshSocketReader(SshSocketReader const&) = delete; // singleton
  void operator=(SshSocketReader const&) = delete; // singleton
  virtual ~SshSocketReader();
  struct event*
  addSshReader(event_callback_fn callbackFn, socket_t fd, void* ptr);
};

} // namespace cli
} // namespace channels
} // namespace devmand
