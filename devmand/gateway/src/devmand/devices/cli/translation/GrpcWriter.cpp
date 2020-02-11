// Copyright (c) 2020-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/devices/cli/translation/GrpcWriter.h>
// TODO: make async

namespace devmand {
namespace devices {
namespace cli {

using namespace folly;
using namespace std;
using namespace grpc;
using namespace devmand::channels::cli::plugin;
using namespace devmand::channels::cli;

GrpcWriter::GrpcWriter(
    shared_ptr<grpc::Channel> channel,
    const string _id,
    shared_ptr<Executor> _executor)
    : GrpcCliHandler(_id, _executor),
      stub_(devmand::channels::cli::plugin::WriterPlugin::NewStub(channel)) {}

Future<Unit> GrpcWriter::create(
    const Path& path,
    dynamic cfg,
    const DeviceAccess& device) const {
  ActualCreateRequest* actualRequest = new ActualCreateRequest();
  actualRequest->set_path(path.str());
  actualRequest->set_cfg(toJson(cfg));
  CreateRequest request;
  request.set_allocated_actualcreaterequest(actualRequest);

  return finish<CreateRequest, CreateResponse, Unit>(
      request,
      device,
      [this](auto context) { return stub_->Create(context); },
      [](auto) { return Future<Unit>(); });
}

Future<Unit> GrpcWriter::update(
    const Path& path,
    dynamic before,
    dynamic after,
    const DeviceAccess& device) const {
  ActualUpdateRequest* actualRequest = new ActualUpdateRequest();
  actualRequest->set_path(path.str());
  actualRequest->set_before(toJson(before));
  actualRequest->set_after(toJson(after));
  UpdateRequest request;
  request.set_allocated_actualupdaterequest(actualRequest);

  return finish<UpdateRequest, UpdateResponse, Unit>(
      request,
      device,
      [this](auto context) { return stub_->Update(context); },
      [](auto) { return Future<Unit>(); });
}

Future<Unit> GrpcWriter::remove(
    const Path& path,
    dynamic before,
    const DeviceAccess& device) const {
  ActualRemoveRequest* actualRequest = new ActualRemoveRequest();
  actualRequest->set_path(path.str());
  actualRequest->set_before(toJson(before));
  RemoveRequest request;
  request.set_allocated_actualremoverequest(actualRequest);

  return finish<RemoveRequest, RemoveResponse, Unit>(
      request,
      device,
      [this](auto context) { return stub_->Remove(context); },
      [](auto) { return Future<Unit>(); });
}

} // namespace cli
} // namespace devices
} // namespace devmand
