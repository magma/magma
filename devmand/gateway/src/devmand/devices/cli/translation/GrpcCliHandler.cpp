// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/devices/cli/translation/GrpcCliHandler.h>

namespace devmand {
namespace devices {
namespace cli {

// TODO: reconnect, error handling - connection issues, wrong services
// provided etc
GrpcCliHandler::GrpcCliHandler(const string _id, shared_ptr<Executor> _executor)
    : id(_id), executor(_executor) {}

CliResponse* GrpcCliHandler::handleCliRequest(
    const DeviceAccess& device,
    const CliRequest& cliRequest,
    bool writingAllowed) const {
  if (not writingAllowed && cliRequest.write()) {
    MLOG(MWARNING) << "[" << id << "] "
                   << "Plugin requested to write command which is forbidden: "
                   << cliRequest.cmd();
    throw runtime_error("Forbidden to execute write commands");
  }
  string cliOutput;
  if (cliRequest.write()) {
    const WriteCommand& command = WriteCommand::create(cliRequest.cmd());
    MLOG(MDEBUG) << "[" << id << "] "
                 << "Got cli request: " << command;
    cliOutput = device.cli()->executeWrite(command).via(executor.get()).get();
  } else {
    const ReadCommand& command = ReadCommand::create(cliRequest.cmd());
    MLOG(MDEBUG) << "[" << id << "] "
                 << "Got cli request: " << command;
    cliOutput = device.cli()->executeRead(command).via(executor.get()).get();
  }
  CliResponse* cliResponse = new CliResponse();
  cliResponse->set_output(cliOutput);
  cliResponse->set_id(cliRequest.id());
  return cliResponse;
}

} // namespace cli
} // namespace devices
} // namespace devmand
