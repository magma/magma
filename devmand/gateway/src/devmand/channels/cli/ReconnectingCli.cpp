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
  reconnectParameters = shared_ptr<ReconnectParameters>(
      new ReconnectParameters{id,
                              timekeeper,
                              executor,
                              move(createCliStack),
                              {},
                              nullptr,
                              quietPeriod,
                              {false},
                              {false}});

  // start async (re)connect
  triggerReconnect(reconnectParameters);
}

Future<Unit> ReconnectingCli::setReconnecting(
    shared_ptr<ReconnectParameters> params) {
  bool f = false;
  if (params->isReconnecting.compare_exchange_strong(f, true)) {
    return makeFuture(unit);
  }
  return via(params->executor.get())
      .thenValue([params](auto) {
        MLOG(MDEBUG) << "[" << params->id << "] "
                     << "setReconnecting: sleeping";
      })
      .delayed(1s, params->timekeeper.get())
      .thenValue(
          [params](auto) -> Future<Unit> { return setReconnecting(params); });
}

SemiFuture<Unit> ReconnectingCli::destroy() {
  if (reconnectParameters->shutdown) {
    return makeSemiFuture(folly::unit);
  }
  reconnectParameters->shutdown = true;
  MLOG(MDEBUG) << "[" << reconnectParameters->id << "] "
               << "destroy: started";

  return setReconnecting(reconnectParameters)
      .thenValue([params = reconnectParameters](auto) {
        SemiFuture<Unit> innerDestroy = folly::makeSemiFuture(unit);
        {
          shared_ptr<Cli> cliOrNull = nullptr;
          {
            //
            boost::mutex::scoped_lock scoped_lock(params->cliMutex);
            cliOrNull = params->maybeCli;
          }
          if (cliOrNull != nullptr) {
            innerDestroy = cliOrNull->destroy();
          }
          // drop cliOrNull
        }
        MLOG(MDEBUG) << "[" << params->id << "] "
                     << "destroy: finished";
        return innerDestroy;
      });
}

ReconnectingCli::~ReconnectingCli() {
  MLOG(MDEBUG) << "[" << reconnectParameters->id << "] "
               << "~ReconnectingCli: started";
  destroy().get();

  reconnectParameters->executor = nullptr;
  reconnectParameters->createCliStack = nullptr;
  reconnectParameters->maybeCli = nullptr;
  reconnectParameters->timekeeper = nullptr;

  MLOG(MDEBUG) << "[" << reconnectParameters->id << "] "
               << "~ReconnectingCli: done";
}

void ReconnectingCli::triggerReconnect(shared_ptr<ReconnectParameters> params) {
  if (params->shutdown) {
    MLOG(MDEBUG) << "[" << params->id << "] "
                 << "triggerReconnect shutting down";
    return;
  }
  bool f = false;
  if (params->isReconnecting.compare_exchange_strong(f, true)) {
    via(params->executor.get(),
        [params]() -> Future<Unit> {
          MLOG(MDEBUG) << "[" << params->id << "] "
                       << "Recreating cli stack";

          bool firstRun = params->maybeCli == nullptr;
          if (not firstRun) {
            MLOG(MDEBUG)
                << "[" << params->id << "] "
                << "Recreating cli stack - calling cli.destroy() on old stack";
            return params->maybeCli->destroy()
                .via(params->executor.get())
                .thenValue([params](auto) -> Unit {
                  MLOG(MDEBUG)
                      << "[" << params->id << "] "
                      << "Recreating cli stack - called cli.destroy() on old stack";

                  boost::mutex::scoped_lock scoped_lock(params->cliMutex);
                  params->maybeCli = nullptr;
                  MLOG(MDEBUG) << "[" << params->id << "] "
                               << "Recreating cli stack - destroyed old stack";
                  return unit;
                });
          } else {
            return unit;
          }
        })
        .thenValue([params](auto) -> Future<Unit> {
          if (params->shutdown) {
            MLOG(MDEBUG)
                << "[" << params->id << "] "
                << "triggerReconnect shutting down before new stack creation";
            params->isReconnecting = false;
            return unit;
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
        }) // TODO: Add onTimeout here to handle prompt resolution timeouts?

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
                return unit;
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
    // Fail fast, does not have effect on correctness
    // Caller will probably get connection exception from inner layers.
    return makeFuture<string>(DisconnectedException());
  }

  {
    boost::try_mutex::scoped_try_lock lock(reconnectParameters->cliMutex);
    if (lock) {
      cliOrNull = reconnectParameters->maybeCli;
    }
  }
  if (reconnectParameters->shutdown) {
    // Fail fast, does not have effect on correctness.
    // Caller will probably get connection exception from inner layers.
    return makeFuture<string>(CliException("Object destroyed"));
  }

  if (cliOrNull != nullptr) {
    return innerFunc(cliOrNull)
        .via(reconnectParameters->executor.get())
        .thenTry([params = reconnectParameters, loggingPrefix, cmd](
                     const Try<string>& t) {
          // Reconnect in case of any CommandExecutionException
          // e.g. write failed, read failed, command timeout
          // TODO unify with PlaintextCliDevice
          // (DisconnectedException+CommandExecutionException) using thenTry to
          // preserve ex type
          if (t.hasException() &&
              t.exception().is_compatible_with<CommandExecutionException>()) {
            MLOG(MDEBUG) << "[" << params->id << "] (" << cmd << ") "
                         << loggingPrefix
                         << " failed due to: " << t.exception().what();
            triggerReconnect(params);
          }
          return t;
        })
        .semi();
  } else {
    return makeFuture<string>(DisconnectedException());
  }
}

} // namespace cli
} // namespace channels
} // namespace devmand
