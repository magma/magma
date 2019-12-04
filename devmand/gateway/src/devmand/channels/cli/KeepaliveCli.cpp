// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/cli/Command.h>
#include <devmand/channels/cli/KeepaliveCli.h>
#include <magma_logging.h>

namespace devmand::channels::cli {

using devmand::channels::cli::Command;
using namespace std;
using namespace folly;

shared_ptr<KeepaliveCli> KeepaliveCli::make(
    string id,
    shared_ptr<Cli> cli,
    shared_ptr<folly::Executor> parentExecutor,
    shared_ptr<folly::Timekeeper> timekeeper,
    chrono::milliseconds heartbeatInterval,
    string keepAliveCommand,
    chrono::milliseconds backoffAfterKeepaliveTimeout) {
  const shared_ptr<KeepaliveCli>& result =
      shared_ptr<KeepaliveCli>(new KeepaliveCli(
          id,
          cli,
          parentExecutor,
          timekeeper,
          heartbeatInterval,
          move(keepAliveCommand),
          backoffAfterKeepaliveTimeout));
  return result;
}

KeepaliveCli::KeepaliveCli(
    string _id,
    shared_ptr<Cli> _cli,
    shared_ptr<Executor> _parentExecutor,
    shared_ptr<Timekeeper> _timekeeper,
    chrono::milliseconds _heartbeatInterval,
    string _keepAliveCommand,
    chrono::milliseconds _backoffAfterKeepaliveTimeout) {
  keepaliveParameters = shared_ptr<KeepaliveParameters>(new KeepaliveParameters{
      /* id */ _id,
      /* cli */ _cli,
      /* timekeeper */ _timekeeper,
      /* parentExecutor */ _parentExecutor,
      /* serialExecutorKeepAlive */
      SerialExecutor::create(
          Executor::getKeepAliveToken(_parentExecutor.get())),
      /* keepAliveCommand */ move(_keepAliveCommand),
      /* heartbeatInterval */ _heartbeatInterval,
      /* backoffAfterKeepaliveTimeout */ _backoffAfterKeepaliveTimeout,
      /* shutdown */ {false}});

  MLOG(MDEBUG) << "[" << _id << "] "
               << "initialized";
  triggerSendKeepAliveCommand(keepaliveParameters);
}

KeepaliveCli::~KeepaliveCli() {
  string id = keepaliveParameters->id;
  keepaliveParameters->shutdown = true;
  MLOG(MDEBUG) << "[" << id << "] "
               << "~KeepaliveCli started";
  while (keepaliveParameters.use_count() >
         1) { // TODO cancel currently running future
    MLOG(MDEBUG) << "[" << id << "] "
                 << "~KeepaliveCli sleeping";
    std::this_thread::sleep_for(std::chrono::seconds(1));
  }
  keepaliveParameters = nullptr;
  MLOG(MDEBUG) << "[" << id << "] "
               << "~KeepaliveCli done";
}

void KeepaliveCli::triggerSendKeepAliveCommand(
    shared_ptr<KeepaliveParameters> keepaliveParameters) {
  if (keepaliveParameters->shutdown) {
    MLOG(MDEBUG) << "[" << keepaliveParameters->id << "] "
                 << "triggerSendKeepAliveCommand shutting down";
    return;
  }
  ReadCommand cmd =
      ReadCommand::create(keepaliveParameters->keepAliveCommand, true);
  MLOG(MDEBUG) << "[" << keepaliveParameters->id << "] (" << cmd.getIdx()
               << ") "
               << "triggerSendKeepAliveCommand created new command";

  via(keepaliveParameters->serialExecutorKeepAlive)
      .thenValue([params = keepaliveParameters, cmd](auto) {
        MLOG(MDEBUG)
            << "[" << params->id << "] (" << cmd.getIdx() << ") "
            << "triggerSendKeepAliveCommand executing keepalive command";

        return params->cli->executeRead(cmd);
      })
      .thenValue([params = keepaliveParameters, cmd](auto) -> SemiFuture<Unit> {
        MLOG(MDEBUG) << "[" << params->id << "] (" << cmd.getIdx() << ") "
                     << "Creating sleep future";
        return futures::sleep(
            params->heartbeatInterval, params->timekeeper.get());
      })
      .thenValue([keepaliveParameters, cmd](auto) -> Unit {
        MLOG(MDEBUG) << "[" << keepaliveParameters->id << "] (" << cmd.getIdx()
                     << ") "
                     << "Woke up after sleep";
        triggerSendKeepAliveCommand(keepaliveParameters);
        return Unit{};
      })
      .thenError(
          folly::tag_t<std::exception>{},
          [params = keepaliveParameters,
           cmd](std::exception const& e) -> Future<Unit> {
            MLOG(MINFO)
                << "[" << params->id << "] (" << cmd.getIdx() << ") "
                << "Got error running keepalive, backing off "
                << e.what(); // FIXME: real exception is not propagated here

            return futures::sleep(
                       params->backoffAfterKeepaliveTimeout,
                       params->timekeeper.get())
                .via(params->serialExecutorKeepAlive)
                .thenValue([params, cmd](auto) -> Unit {
                  MLOG(MDEBUG)
                      << "[" << params->id << "] (" << cmd.getIdx() << ") "
                      << "Woke up after backing off";
                  triggerSendKeepAliveCommand(params);
                  return Unit{};
                });
          });
}

folly::SemiFuture<std::string> KeepaliveCli::executeRead(
    const ReadCommand cmd) {
  return keepaliveParameters->cli->executeRead(cmd);
}

folly::SemiFuture<std::string> KeepaliveCli::executeWrite(
    const WriteCommand cmd) {
  return keepaliveParameters->cli->executeWrite(cmd);
}

} // namespace devmand::channels::cli
