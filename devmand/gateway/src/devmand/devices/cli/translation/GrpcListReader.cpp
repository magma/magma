// Copyright (c) 2020-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/devices/cli/translation/GrpcListReader.h>

namespace devmand {
namespace devices {
namespace cli {

using namespace folly;
using namespace std;
using namespace grpc;
using namespace devmand::channels::cli::plugin;
using namespace devmand::channels::cli;

GrpcListReader::GrpcListReader(
    shared_ptr<grpc::Channel> channel,
    const string _id,
    shared_ptr<Executor> _executor)
    :
    GrpcCliHandler(_id, _executor),
    stub_(devmand::channels::cli::plugin::ReaderPlugin::NewStub(channel))
{}

Future<vector<dynamic>> GrpcListReader::readKeys(
    const Path& path,
    const DeviceAccess& device) const {

  ActualReadRequest* actualRequest = new ActualReadRequest();
  actualRequest->set_path(path.str());
  ReadRequest request;
  request.set_allocated_actualreadrequest(actualRequest);

  return finish<ReadRequest, ReadResponse, vector<dynamic>>(
      request,
      device,
      [this](auto context) { return stub_->Read(context); },
      [this](auto response) -> vector<dynamic> {
        dynamic result = parseJson(response.actualreadresponse().json());
        if (not result.isArray()) {
          MLOG(MERROR) << "[" << id << "] Response is not json array:" << response.actualreadresponse().json() ;
          throw runtime_error("Response is not json array");
        }
        vector<dynamic> values;
        for (auto& k : result) {
          MLOG(MDEBUG) << "pushing " << toJson(k);
          values.push_back(k);
        }
        return values;
      });
}

} // namespace cli
} // namespace devices
} // namespace devmand
