// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/channels/cli/Channel.h>
#include <folly/executors/CPUThreadPoolExecutor.h>

namespace devmand {
namespace devices {
namespace cli {

using namespace std;
using namespace folly;
using namespace devmand::channels::cli;

class DeviceAccess {
 public:
  DeviceAccess(
      shared_ptr<Cli> _cliChannel,
      string _deviceId,
      shared_ptr<Executor> _workerExecutor);
  ~DeviceAccess() = default;

 public:
  shared_ptr<Cli> cli() const;
  string id() const;
  shared_ptr<Executor> executor() const;

 private:
  shared_ptr<Cli> cliChannel;
  string deviceId;
  shared_ptr<Executor> workerExecutor;
};

} // namespace cli
} // namespace devices
} // namespace devmand
