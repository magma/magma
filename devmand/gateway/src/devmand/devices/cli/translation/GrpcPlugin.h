// Copyright (c) 2020-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/channels/cli/plugin/protocpp/PluginRegistration.pb.h>
#include <devmand/channels/cli/plugin/protocpp/ReaderPlugin.grpc.pb.h>
#include <devmand/devices/cli/translation/PluginRegistry.h>
#include <grpc++/grpc++.h>

namespace devmand {
namespace devices {
namespace cli {

using namespace std;
using namespace folly;
using namespace devmand::channels::cli;

class GrpcPlugin : public Plugin {
 private:
  shared_ptr<grpc::Channel> channel;
  const string id;
  shared_ptr<Executor> executor;
  devmand::channels::cli::plugin::CapabilitiesResponse capabilities;

  GrpcPlugin(
      shared_ptr<grpc::Channel> channel,
      const string id,
      shared_ptr<Executor> executor,
      devmand::channels::cli::plugin::CapabilitiesResponse capabilities);

 public:
  static shared_ptr<GrpcPlugin> create(
      shared_ptr<grpc::Channel> channel,
      const string id,
      shared_ptr<Executor> executor);

  DeviceType getDeviceType() const override;

  void provideReaders(ReaderRegistryBuilder& registry) const override;

  void provideWriters(WriterRegistryBuilder& registry) const override;
};

} // namespace cli
} // namespace devices
} // namespace devmand
