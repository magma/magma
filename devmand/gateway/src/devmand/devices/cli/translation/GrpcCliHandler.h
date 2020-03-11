// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <devmand/channels/cli/plugin/protocpp/Common.pb.h>
#include <devmand/channels/cli/plugin/protocpp/ReaderPlugin.grpc.pb.h>
#include <devmand/devices/cli/translation/DeviceAccess.h>
#include <folly/Executor.h>
#include <grpc++/grpc++.h>
#include <chrono>

namespace devmand {
namespace devices {
namespace cli {

using namespace folly;
using namespace std;
using namespace grpc;
using namespace devmand::channels::cli::plugin;
using namespace devmand::channels::cli;
using namespace std::chrono;

typedef high_resolution_clock clock;

class GrpcCliHandler {
 protected:
  const string id;

 private:
  shared_ptr<Executor> executor;

 public:
  GrpcCliHandler(const string _id, shared_ptr<Executor> _executor);

 private:
  // Handle single cli request, writing response back
  CliResponse* handleCliRequest(
      const DeviceAccess& device,
      const CliRequest& cliRequest,
      bool writingAllowed) const;

  // Loop while remote plugin sends CLI requests, exit when final response is
  // received
  template <class RequestClass, class ResponseClass>
  long handleCliRequests(
      const DeviceAccess& device,
      ClientReaderWriter<RequestClass, ResponseClass>* stream,
      ResponseClass* response,
      bool writingAllowed) const {
    long int spentInCliMillis = 0;
    while (stream->Read(response) && response->has_clirequest()) {
      auto cliStartTime = clock::now();
      CliResponse* cliResponse =
          handleCliRequest(device, response->clirequest(), writingAllowed);
      RequestClass nextRequest;
      nextRequest.set_allocated_cliresponse(cliResponse);
      spentInCliMillis +=
          (duration_cast<milliseconds>(clock::now() - cliStartTime)).count();
      stream->Write(nextRequest);
    }
    // final response is in `response`
    return spentInCliMillis;
  }

 protected:
  // Do the actual RPC request, handle Cli requests and transform final response
  // using closure resultTransformer.
  template <class RequestClass, class ResponseClass, class ResultClass>
  Future<ResultClass> finish(
      RequestClass request,
      const DeviceAccess& device,
      function<unique_ptr<ClientReaderWriter<RequestClass, ResponseClass>>(
          ClientContext*)> rpc,
      function<Future<ResultClass>(ResponseClass)> resultTransformer) const {
    auto startTime = clock::now();
    ClientContext context;
    unique_ptr<ClientReaderWriter<RequestClass, ResponseClass>> stream(
        rpc(&context)); // TODO async

    if (not stream) {
      MLOG(MWARNING) << "[" << id << "] Cannot connect";
      return makeFuture<ResultClass>(runtime_error("Cannot connect"));
    }

    // send the request
    stream->Write(request);

    // start reading responses
    ResponseClass response;
    long spentInCliMillis = handleCliRequests<RequestClass, ResponseClass>(
        device, stream.get(), &response, true);

    auto totalMillis =
        (duration_cast<milliseconds>(clock::now() - startTime)).count();
    MLOG(MDEBUG) << "[" << id << "] Total duration: " << totalMillis << " ms, "
                 << "in grpc: " << (totalMillis - spentInCliMillis) << " ms, "
                 << "in cli " << spentInCliMillis << " ms";
    Status status = stream->Finish();

    if (status.ok()) {
      return resultTransformer(response);
    } else {
      MLOG(MWARNING) << "[" << id << "] Error " << status.error_code() << ": "
                     << status.error_message();
      return makeFuture<ResultClass>(runtime_error("RPC failed"));
    }
  }
};

} // namespace cli
} // namespace devices
} // namespace devmand
