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
    std::chrono::milliseconds _timeoutInterval)
    : id(_id),
      sharedCli(_cli),
      sharedTimekeeper(_timekeeper),
      sharedExecutor(_executor),
      timeoutInterval(_timeoutInterval) {}

SemiFuture<Unit> TimeoutTrackingCli::destroy() {
  MLOG(MDEBUG) << "[" << id << "] "
               << "destroy: started";
  // call underlying destroy()
  SemiFuture<Unit> innerDestroy = sharedCli->destroy();

  // TODO cancel timekeeper futures

  MLOG(MDEBUG) << "[" << id << "] "
               << "destroy: done";
  return innerDestroy;
}

TimeoutTrackingCli::~TimeoutTrackingCli() {
  MLOG(MDEBUG) << "[" << id << "] "
               << "~TTCli: started";
  destroy().get();
  sharedExecutor = nullptr;
  sharedCli = nullptr;
  sharedTimekeeper = nullptr; // clean timekeeper as last to keep tests from
                              // failing on "No timekeeper available"
  MLOG(MDEBUG) << "[" << id << "] "
               << "~TTCli: done";
}

folly::SemiFuture<std::string> TimeoutTrackingCli::executeRead(
    const ReadCommand cmd) {
  return executeSomething(
             cmd,
             "TTCli.executeRead",
             [cmd](shared_ptr<Cli> cli) { return cli->executeRead(cmd); })
      .semi();
}

folly::SemiFuture<std::string> TimeoutTrackingCli::executeWrite(
    const WriteCommand cmd) {
  return executeSomething(
             cmd,
             "TTCli.executeWrite",
             [cmd](shared_ptr<Cli> cli) { return cli->executeWrite(cmd); })
      .semi();
}

Future<string> TimeoutTrackingCli::executeSomething(
    const Command& cmd,
    const string&& loggingPrefix,
    const function<SemiFuture<string>(shared_ptr<Cli> cli)>& innerFunc) {
  MLOG(MDEBUG) << "[" << id << "] (" << cmd << ") " << loggingPrefix << "('"
               << cmd << "') called";

  SemiFuture<string> inner =
      innerFunc(sharedCli); // we expect that this method does not block
  MLOG(MDEBUG) << "[" << id << "] (" << cmd << ") "
               << "Obtained future from underlying cli";
  return move(inner)
      .via(sharedExecutor.get())
      .onTimeout(
          timeoutInterval,
          [_id = id, cmd](...) -> Future<string> {
            MLOG(MDEBUG) << "[" << _id << "] (" << cmd << ") "
                         << "timing out";
            throw CommandTimeoutException();
          },
          sharedTimekeeper.get())
      .thenValue([_id = id, cmd](string result) -> string {
        MLOG(MDEBUG) << "[" << _id << "] (" << cmd << ") "
                     << "succeeded";
        return result;
      });
}

} // namespace devmand::channels::cli
