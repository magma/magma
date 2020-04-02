// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/devices/cli/translation/DeviceAccess.h>

namespace devmand {
namespace devices {
namespace cli {

DeviceAccess::DeviceAccess(
    shared_ptr<Cli> _cliChannel,
    string _deviceId,
    shared_ptr<Executor> _workerExecutor)
    : cliChannel(_cliChannel),
      deviceId(_deviceId),
      workerExecutor(_workerExecutor) {}

shared_ptr<Cli> DeviceAccess::cli() const {
  return cliChannel;
}

shared_ptr<Executor> DeviceAccess::executor() const {
  return workerExecutor;
}

string DeviceAccess::id() const {
  return deviceId;
}

} // namespace cli
} // namespace devices
} // namespace devmand
