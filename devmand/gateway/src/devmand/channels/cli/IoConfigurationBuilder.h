// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <devmand/cartography/DeviceConfig.h>
#include <devmand/channels/cli/Cli.h>
#include <devmand/channels/cli/CliFlavour.h>
#include <devmand/channels/cli/ReadCachingCli.h>

using devmand::cartography::DeviceConfig;
using devmand::channels::cli::Cli;
using devmand::channels::cli::CliFlavour;
using std::shared_ptr;

namespace devmand::channels::cli {
class IoConfigurationBuilder {
 private:
  DeviceConfig deviceConfig;

 public:
  IoConfigurationBuilder(const DeviceConfig& deviceConfig);

  shared_ptr<Cli> getIo(
      shared_ptr<CliCache> commandCache = ReadCachingCli::createCache());
};
}