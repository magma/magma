// Copyright (c) 2020-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/channels/cli/plugin/protocpp/WriterPlugin.grpc.pb.h>
#include <devmand/devices/cli/translation/GrpcCliHandler.h>
#include <devmand/devices/cli/translation/PluginRegistry.h>
#include <grpc++/grpc++.h>

namespace devmand {
namespace devices {
namespace cli {

using namespace std;
using namespace folly;
using namespace devmand::channels::cli;

class GrpcWriter : public Writer, public GrpcCliHandler {
 private:
  unique_ptr<devmand::channels::cli::plugin::WriterPlugin::Stub> stub_;

 public:
  GrpcWriter(
      shared_ptr<grpc::Channel> channel,
      const string id,
      shared_ptr<Executor> executor);

  Future<Unit> create(const Path& path, dynamic cfg, const DeviceAccess& device)
      const override;

  Future<Unit> update(
      const Path& path,
      dynamic before,
      dynamic after,
      const DeviceAccess& device) const override;

  Future<Unit> remove(
      const Path& path,
      dynamic before,
      const DeviceAccess& device) const override;
};

} // namespace cli
} // namespace devices
} // namespace devmand
