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
    shared_ptr<Executor> parentExecutor,
    shared_ptr<Timekeeper> timekeeper,
    chrono::milliseconds heartbeatInterval,
    string keepAliveCommand,
    chrono::milliseconds backoffAfterKeepaliveTimeout) {
  const shared_ptr<KeepaliveCli>& result = std::make_shared<KeepaliveCli>(
      id,
      cli,
      parentExecutor,
      timekeeper,
      heartbeatInterval,
      move(keepAliveCommand),
      backoffAfterKeepaliveTimeout);

  return result;
}

KeepaliveCli::KeepaliveCli(
    string _id,
    shared_ptr<Cli> _cli,
    shared_ptr<Executor> _parentExecutor,
    shared_ptr<Timekeeper> _timekeeper,
    chrono::milliseconds _heartbeatInterval,
    string _keepAliveCommand,
    chrono::milliseconds _backoffAfterKeepaliveTimeout)
    : sharedCli(_cli),
      sharedTimekeeper(_timekeeper),
      parentExecutor(_parentExecutor),
      sharedSerialExecutorKeepAlive(
          make_shared<Executor::KeepAlive<SerialExecutor>>(
              SerialExecutor::create(
                  Executor::getKeepAliveToken(_parentExecutor.get())))) {
  keepaliveParameters = std::make_shared<KeepaliveParameters>(
      /* id */ _id,
      /* cli */ _cli,
      /* timekeeper */ _timekeeper,
      /* serialExecutorKeepAlive */ sharedSerialExecutorKeepAlive,
      /* keepAliveCommand */ move(_keepAliveCommand),
      /* heartbeatInterval */ _heartbeatInterval,
      /* backoffAfterKeepaliveTimeout */ _backoffAfterKeepaliveTimeout);

  MLOG(MDEBUG) << "[" << _id << "] "
               << "initialized";
  triggerSendKeepAliveCommand(keepaliveParameters);
}

SemiFuture<Unit> KeepaliveCli::destroy() {
  MLOG(MDEBUG) << "[" << keepaliveParameters->id << "] "
               << "destroy: started";
  // call underlying destroy()
  SemiFuture<Unit> innerDestroy = sharedCli->destroy();
  // TODO cancel timekeeper futures

  MLOG(MDEBUG) << "[" << keepaliveParameters->id << "] "
               << "destroy: done";
  return innerDestroy;
}

KeepaliveCli::~KeepaliveCli() {
  MLOG(MDEBUG) << "[" << keepaliveParameters->id << "] "
               << "~KeepaliveCli: started";
  destroy().get();
  MLOG(MDEBUG) << "[" << keepaliveParameters->id << "] "
               << "~KeepaliveCli: started";
}

void KeepaliveCli::triggerSendKeepAliveCommand(
    shared_ptr<KeepaliveParameters> keepaliveParameters) {
  if (shared_ptr<Executor::KeepAlive<SerialExecutor>> serialExecutor =
          keepaliveParameters->serialExecutorKeepAlive.lock()) {
    ReadCommand cmd =
        ReadCommand::create(keepaliveParameters->keepAliveCommand, true);
    MLOG(MDEBUG) << "[" << keepaliveParameters->id << "] (" << cmd << ") "
                 << "triggerSendKeepAliveCommand created new command";

    via(serialExecutor->get())
        .thenValue(
            [params = keepaliveParameters, cmd](auto) -> SemiFuture<string> {
              if (shared_ptr<Cli> cli = params->cli.lock()) {
                MLOG(MDEBUG)
                    << "[" << params->id << "] (" << cmd << ") "
                    << "triggerSendKeepAliveCommand executing keepalive command";
                return cli->executeRead(cmd);
              }
              throw CliException("Object destroyed");
            })
        .thenValue(
            [params = keepaliveParameters, cmd](auto) -> SemiFuture<Unit> {
              if (shared_ptr<Timekeeper> timekeeper =
                      params->timekeeper.lock()) {
                MLOG(MDEBUG) << "[" << params->id << "] (" << cmd << ") "
                             << "Creating sleep future";
                return futures::sleep(
                    params->heartbeatInterval, timekeeper.get());
              }
              throw CliException("Object destroyed");
            })
        .thenValue([keepaliveParameters, cmd](auto) -> Unit {
          MLOG(MDEBUG) << "[" << keepaliveParameters->id << "] (" << cmd << ") "
                       << "Woke up after sleep";
          triggerSendKeepAliveCommand(keepaliveParameters);
          return Unit{};
        })
        .thenError([params = keepaliveParameters,
                    cmd](const exception_wrapper& e) {
          MLOG(MINFO) << "[" << params->id << "] (" << cmd << ") "
                      << "Got error running keepalive, backing off "
                      << e.what();
          if (shared_ptr<Timekeeper> timekeeper = params->timekeeper.lock()) {
            if (shared_ptr<Executor::KeepAlive<SerialExecutor>>
                    serialExecutor2 = params->serialExecutorKeepAlive.lock()) {
              // destructor
              return futures::sleep(
                         params->backoffAfterKeepaliveTimeout, timekeeper.get())
                  .via(serialExecutor2->get())
                  .thenValue([params, cmd](auto) -> Unit {
                    MLOG(MDEBUG) << "[" << params->id << "] (" << cmd << ") "
                                 << "Woke up after backing off";
                    triggerSendKeepAliveCommand(params);
                    return Unit{};
                  });
            }
          }
          throw CliException("Object destroyed");
        });
  }
}

SemiFuture<std::string> KeepaliveCli::executeRead(const ReadCommand cmd) {
  return sharedCli->executeRead(cmd);
}

SemiFuture<std::string> KeepaliveCli::executeWrite(const WriteCommand cmd) {
  return sharedCli->executeWrite(cmd);
}

KeepaliveCli::KeepaliveParameters::KeepaliveParameters(
    const string& _id,
    weak_ptr<Cli> _cli,
    weak_ptr<folly::Timekeeper> _timekeeper,
    weak_ptr<folly::Executor::KeepAlive<folly::SerialExecutor>>
        _serialExecutorKeepAlive,
    const string& _keepAliveCommand,
    const chrono::milliseconds& _heartbeatInterval,
    const chrono::milliseconds& _backoffAfterKeepaliveTimeout)
    : id(_id),
      cli(_cli),
      timekeeper(_timekeeper),
      serialExecutorKeepAlive(_serialExecutorKeepAlive),
      keepAliveCommand(_keepAliveCommand),
      heartbeatInterval(_heartbeatInterval),
      backoffAfterKeepaliveTimeout(_backoffAfterKeepaliveTimeout) {}
} // namespace devmand::channels::cli
