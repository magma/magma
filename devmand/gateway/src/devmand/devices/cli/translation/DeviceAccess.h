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

namespace devmand {
namespace devices {
namespace cli {

using namespace std;
using namespace devmand::channels::cli;

class DeviceAccess {
 public:
  DeviceAccess(shared_ptr<Cli> _cliChannel, string _deviceId);
  ~DeviceAccess() = default;

 public:
  shared_ptr<Cli> cli() const;
  string id() const;

 private:
  shared_ptr<Cli> cliChannel;
  string deviceId;
};

} // namespace cli
} // namespace devices
} // namespace devmand
