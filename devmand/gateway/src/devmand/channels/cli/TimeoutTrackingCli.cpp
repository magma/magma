// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/cli/Command.h>
#include <devmand/channels/cli/TimeoutTrackingCli.h>
#include <magma_logging.h>

namespace devmand::channels::cli {

using namespace std;
using namespace folly;
using devmand::channels::cli::Command;

shared_ptr<TimeoutTrackingCli> TimeoutTrackingCli::make(
    string id,
    shared_ptr<Cli> cli,
    shared_ptr<folly::Timekeeper> timekeeper,
    shared_ptr<folly::Executor> executor,
    std::chrono::milliseconds timeoutInterval) {
  return std::make_shared<TimeoutTrackingCli>(
      id, cli, timekeeper, executor, timeoutInterval);
}

TimeoutTrackingCli::TimeoutTrackingCli(
    string _id,
    shared_ptr<Cli> _cli,
    shared_ptr<folly::Timekeeper> _timekeeper,
    shared_ptr<folly::Executor> _executor,
    std::chrono::milliseconds _timeoutInterval) {
  timeoutTrackingParameters = std::make_shared<TimeoutTrackingParameters>(
      _id, _cli, _timekeeper, _executor, _timeoutInterval, false);
}

TimeoutTrackingCli::~TimeoutTrackingCli() {
  string id = timeoutTrackingParameters->id;
  MLOG(MDEBUG) << "[" << id << "] "
               << "~TTCli started";
  timeoutTrackingParameters->shutdown = true;
  while (timeoutTrackingParameters.use_count() >
         1) { // TODO cancel currently running future
    MLOG(MDEBUG) << "[" << timeoutTrackingParameters->id << "] "
                 << "~TTCli sleeping";
    std::this_thread::sleep_for(std::chrono::seconds(1));
  }

  MLOG(MDEBUG) << "[" << id << "] "
               << "~TTCli nulling timeoutTrackingParameters.cli";
  timeoutTrackingParameters->cli = nullptr;
  MLOG(MDEBUG) << "[" << id << "] "
               << "~TTCli nulling timeoutTrackingParameters";
  timeoutTrackingParameters = nullptr;
  MLOG(MDEBUG) << "[" << id << "] "
               << "~TTCli done";
}

folly::SemiFuture<std::string> TimeoutTrackingCli::executeRead(
    const ReadCommand cmd) {
  return executeSomething(
             cmd,
             "TTCli.executeRead",
             [params = timeoutTrackingParameters, cmd]() {
               return params->cli->executeRead(cmd);
             })
      .semi();
}

folly::SemiFuture<std::string> TimeoutTrackingCli::executeWrite(
    const WriteCommand cmd) {
  return executeSomething(
             cmd,
             "TTCli.executeWrite",
             [params = timeoutTrackingParameters, cmd]() {
               return params->cli->executeWrite(cmd);
             })
      .semi();
}

Future<string> TimeoutTrackingCli::executeSomething(
    const Command& cmd,
    const string&& loggingPrefix,
    const function<SemiFuture<string>()>& innerFunc) {
  MLOG(MDEBUG) << "[" << timeoutTrackingParameters->id << "] (" << cmd << ") "
               << loggingPrefix << "('" << cmd << "') called";
  if (timeoutTrackingParameters->shutdown) {
    return Future<string>(runtime_error("TTCli Shutting down"));
  }
  SemiFuture<string> inner =
      innerFunc(); // we expect that this method does not block
  MLOG(MDEBUG) << "[" << timeoutTrackingParameters->id << "] (" << cmd << ") "
               << "Obtained future from underlying cli";
  return move(inner)
      .via(timeoutTrackingParameters->executor.get())
      .onTimeout(
          timeoutTrackingParameters->timeoutInterval,
          [params = timeoutTrackingParameters, cmd](...) -> Future<string> {
            // NOTE: timeoutTrackingParameters must be captured mainly for
            // executor
            MLOG(MDEBUG) << "[" << params->id << "] (" << cmd << ") "
                         << "timing out";
            throw CommandTimeoutException();
          },
          timeoutTrackingParameters->timekeeper.get())
      .thenValue(
          [params = timeoutTrackingParameters, cmd](string result) -> string {
            MLOG(MDEBUG) << "[" << params->id << "] (" << cmd << ") "
                         << "succeeded";
            return result;
          });
}

TimeoutTrackingCli::TimeoutTrackingParameters::TimeoutTrackingParameters(
    const string& _id,
    const shared_ptr<Cli>& _cli,
    const shared_ptr<folly::Timekeeper>& _timekeeper,
    const shared_ptr<folly::Executor>& _executor,
    const chrono::milliseconds& _timeoutInterval,
    const bool _shutdown)
    : id(_id),
      cli(_cli),
      timekeeper(_timekeeper),
      executor(_executor),
      timeoutInterval(_timeoutInterval),
      shutdown(ATOMIC_VAR_INIT(_shutdown)) {}
} // namespace devmand::channels::cli
