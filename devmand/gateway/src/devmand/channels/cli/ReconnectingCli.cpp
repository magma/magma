// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/cli/ReconnectingCli.h>
#include <magma_logging.h>

namespace devmand {
namespace channels {
namespace cli {

using namespace std;
using namespace folly;

shared_ptr<ReconnectingCli> ReconnectingCli::make(
    string id,
    shared_ptr<Executor> executor,
    function<SemiFuture<shared_ptr<Cli>>()>&& createCliStack,
    shared_ptr<Timekeeper> timekeeper,
    chrono::milliseconds quietPeriod) {
  return std::make_shared<ReconnectingCli>(
      id, executor, move(createCliStack), move(timekeeper), move(quietPeriod));
}

ReconnectingCli::ReconnectingCli(
    string id,
    shared_ptr<Executor> executor,
    function<SemiFuture<shared_ptr<Cli>>()>&& createCliStack,
    shared_ptr<Timekeeper> timekeeper,
    chrono::milliseconds quietPeriod) {
  reconnectParameters = make_shared<ReconnectParameters>();
  reconnectParameters->id = id;
  reconnectParameters->executor = executor;
  reconnectParameters->createCliStack = move(createCliStack);
  reconnectParameters->maybeCli = nullptr;
  reconnectParameters->shutdown = false;
  reconnectParameters->isReconnecting = false;
  reconnectParameters->quietPeriod = quietPeriod;
  reconnectParameters->timekeeper = timekeeper;
  // start async (re)connect
  triggerReconnect(reconnectParameters);
}

ReconnectingCli::~ReconnectingCli() {
  string id = reconnectParameters->id;
  MLOG(MDEBUG) << "[" << id << "] "
               << "~RCli started";
  reconnectParameters->shutdown = true;
  while (reconnectParameters.use_count() >
         1) { // TODO cancel currently running future
    MLOG(MDEBUG) << "[" << id << "] "
                 << "~RCli sleeping";
    std::this_thread::sleep_for(std::chrono::seconds(1));
  }

  reconnectParameters = nullptr;
  MLOG(MDEBUG) << "[" << id << "] "
               << "~RCli done";
}

void ReconnectingCli::triggerReconnect(shared_ptr<ReconnectParameters> params) {
  if (params->shutdown)
    return;
  bool f = false;
  if (params->isReconnecting.compare_exchange_strong(f, true)) {
    via(params->executor.get(),
        [params]() -> Future<Unit> {
          MLOG(MDEBUG) << "[" << params->id << "] "
                       << "Recreating cli stack";
          {
            boost::mutex::scoped_lock scoped_lock(params->cliMutex);
            params->maybeCli = nullptr;
            MLOG(MDEBUG) << "[" << params->id << "] "
                         << "Recreating cli stack - destroyed old stack";
          }
          Future<shared_ptr<Cli>> newCliFuture =
              params->createCliStack().via(params->executor.get());

          return move(newCliFuture)
              .thenValue([params](shared_ptr<Cli> newCli) -> Unit {
                {
                  boost::mutex::scoped_lock scoped_lock(params->cliMutex);

                  params->maybeCli = std::move(newCli);
                }
                params->isReconnecting = false;
                MLOG(MDEBUG) << "[" << params->id << "] "
                             << "Recreating cli stack - done";
                return unit;
              });
        }) // TODO: Add onTimeout here to handle establish session timeouts?

        .thenError([params](const exception_wrapper& e) -> Future<Unit> {
          // quiet period
          MLOG(MDEBUG) << "[" << params->id << "] "
                       << "triggerReconnect started quiet period, got error : "
                       << e.what();

          return futures::sleep(params->quietPeriod, params->timekeeper.get())
              .via(params->executor.get())
              .thenValue([params](Unit) -> Future<Unit> {
                MLOG(MDEBUG)
                    << "[" << params->id << "] "
                    << "triggerReconnect reconnecting after quiet period";
                params->isReconnecting = false;
                triggerReconnect(params);
                return makeFuture();
              });
        });
  } else {
    MLOG(MDEBUG) << "[" << params->id << "] "
                 << "Already reconnecting";
  }
}

folly::SemiFuture<std::string> ReconnectingCli::executeRead(
    const ReadCommand cmd) {
  // capturing this is ok here - lambda is evaluated synchronously
  return executeSomething(
      "RCli.executeRead",
      [cmd](shared_ptr<Cli> cli) { return cli->executeRead(cmd); },
      cmd);
}

folly::SemiFuture<std::string> ReconnectingCli::executeWrite(
    const WriteCommand cmd) {
  // capturing this is ok here - lambda is evaluated synchronously
  return executeSomething(
      "RCli.executeWrite",
      [cmd](shared_ptr<Cli> cli) { return cli->executeWrite(cmd); },
      cmd);
}

SemiFuture<string> ReconnectingCli::executeSomething(
    const string&& loggingPrefix,
    const function<SemiFuture<string>(shared_ptr<Cli>)>& innerFunc,
    const Command cmd) {
  shared_ptr<Cli> cliOrNull = nullptr;

  if (reconnectParameters->isReconnecting) {
    return makeFuture<string>(DisconnectedException());
  }

  { // TODO: trylock and throw on already locked. This means reconnect is in
    // progress
    boost::mutex::scoped_lock scoped_lock(reconnectParameters->cliMutex);
    cliOrNull = reconnectParameters->maybeCli;
  }

  if (cliOrNull == nullptr) {
    return makeFuture<string>(DisconnectedException());
  }

  return innerFunc(cliOrNull)
      .via(reconnectParameters->executor.get())
      .thenTry([params = reconnectParameters, loggingPrefix, cmd](
                   const Try<string>& t) {
        // Reconnect in case of any CommandExecutionException
        // e.g. write failed, read failed, command timeout
        // using thenTry to preserve ex type
        if (t.hasException() &&
            t.exception().is_compatible_with<CommandExecutionException>()) {
          MLOG(MDEBUG) << "[" << params->id << "] (" << cmd << ") "
                       << loggingPrefix
                       << " failed due to: " << t.exception().what();

          // Using "this" raw pointer, however we have the
          // shared_ptr<params> to protect against destructor call
          triggerReconnect(params);
        }
        return t;
      })

      .semi();
}

} // namespace cli
} // namespace channels
} // namespace devmand
