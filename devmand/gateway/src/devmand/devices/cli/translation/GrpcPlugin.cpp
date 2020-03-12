// Copyright (c) 2020-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/cli/plugin/protocpp/PluginRegistration.grpc.pb.h>
#include <devmand/devices/cli/translation/GrpcListReader.h>
#include <devmand/devices/cli/translation/GrpcPlugin.h>
#include <devmand/devices/cli/translation/GrpcReader.h>
#include <devmand/devices/cli/translation/GrpcWriter.h>
// TODO make async

namespace devmand {
namespace devices {
namespace cli {

using namespace grpc;
using namespace devmand::channels::cli::plugin;

shared_ptr<GrpcPlugin> GrpcPlugin::create(
    shared_ptr<grpc::Channel> channel,
    const string id,
    shared_ptr<Executor> executor) {
  // obtain capabilities via gRPC
  // TODO: remote server unavailability/reconnecting handling
  unique_ptr<devmand::channels::cli::plugin::PluginRegistration::Stub> stub(
      devmand::channels::cli::plugin::PluginRegistration::NewStub(channel));
  ClientContext context;
  CapabilitiesRequest request;
  request.set_id(id);
  CapabilitiesResponse response;
  Status status = stub->GetCapabilities(&context, request, &response);
  if (status.ok()) {
    return make_shared<GrpcPlugin>(GrpcPlugin(channel, id, executor, response));
  } else {
    MLOG(MERROR) << "[" << id << "] Error " << status.error_code() << ": "
                 << status.error_message();
    throw runtime_error("PluginRegistration RPC failed");
  }
}

GrpcPlugin::GrpcPlugin(
    shared_ptr<grpc::Channel> _channel,
    const string _id,
    shared_ptr<Executor> _executor,
    CapabilitiesResponse _capabilities)
    : channel(_channel),
      id(_id),
      executor(_executor),
      capabilities(_capabilities) {}

DeviceType GrpcPlugin::getDeviceType() const {
  if (capabilities.has_devicetype()) {
    return {capabilities.devicetype().device(),
            capabilities.devicetype().version()};
  }
  return DeviceType::getDefaultInstance();
}

void GrpcPlugin::provideReaders(ReaderRegistryBuilder& registry) const {
  for (int i = 0; i < capabilities.readers_size(); i++) {
    Path path(capabilities.readers().Get(i).path());
    auto remoteReaderPlugin = make_shared<GrpcReader>(channel, id, executor);
    registry.add(path, remoteReaderPlugin);
  }
  for (int i = 0; i < capabilities.listreaders_size(); i++) {
    Path path(capabilities.listreaders().Get(i).path());
    auto remoteReaderPlugin =
        make_shared<GrpcListReader>(channel, id, executor);
    registry.addList(path, remoteReaderPlugin);
  }
}

void GrpcPlugin::provideWriters(WriterRegistryBuilder& registry) const {
  for (int i = 0; i < capabilities.writers_size(); i++) {
    WriterCapability writerCapability = capabilities.writers().Get(i);
    Path path(writerCapability.path());
    vector<Path> dependencies;
    for (int depIdx = 0; depIdx < writerCapability.dependencies_size();
         depIdx++) {
      dependencies.push_back(Path(writerCapability.dependencies(depIdx)));
    }
    auto remoteWriterPlugin = make_shared<GrpcWriter>(channel, id, executor);
    registry.add(path, remoteWriterPlugin, dependencies);
  }
}

static Optional<char> sanitizeSingleIndentChar(const string& str) {
  if (str.size() == 1) {
    return Optional<char>(str[0]);
  } else if (not str.empty()) {
    MLOG(MWARNING)
        << "sanitizeSingleIndentChar expected 0 or 1 character, got '" << str
        << "'";
  }
  return folly::none;
}

Optional<CliFlavourParameters> GrpcPlugin::getCliFlavourParameters() const {
  if (capabilities.has_devicetype() &&
      capabilities.devicetype().has_cliflavourparams()) {
    CliFlavourParameters result = {
        capabilities.devicetype().cliflavourparams().newline(),
        regex(
            capabilities.devicetype().cliflavourparams().baseshowconfigregex()),
        capabilities.devicetype().cliflavourparams().baseshowconfigidx(),
        sanitizeSingleIndentChar(
            capabilities.devicetype().cliflavourparams().singleindentchar()),
        capabilities.devicetype().cliflavourparams().configsubsectionend()};
    return result;
  }
  return none;
}

} // namespace cli
} // namespace devices
} // namespace devmand
