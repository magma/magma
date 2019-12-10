// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <boost/range/size_type.hpp>
#include <devmand/channels/cli/QueuedCli.h>

namespace devmand {
namespace channels {
namespace cli {

using namespace std;
using namespace folly;

shared_ptr<QueuedCli> QueuedCli::make(
    string id,
    shared_ptr<Cli> cli,
    shared_ptr<Executor> parentExecutor) {
  return std::make_shared<QueuedCli>(id, cli, parentExecutor);
}

QueuedCli::QueuedCli(
    string _id,
    shared_ptr<Cli> _cli,
    shared_ptr<Executor> _parentExecutor) {
  queuedParameters = std::make_shared<QueuedParameters>(
      _id,
      _cli,
      _parentExecutor,
      SerialExecutor::create(
          Executor::getKeepAliveToken(_parentExecutor.get())));
}

QueuedCli::~QueuedCli() {
  string id = queuedParameters->id;
  MLOG(MDEBUG) << "[" << id << "] "
               << "~QCli started";
  queuedParameters->shutdown = true;
  MLOG(MDEBUG) << "[" << id << "] "
               << "~QCli: dequeuing " << queuedParameters->queue.size()
               << " items";
  {
    QueueEntry queueEntry;
    while (queuedParameters->queue.try_dequeue(queueEntry)) {
      MLOG(MDEBUG) << "[" << id << "] (" << queueEntry.command << ") "
                   << "~QCli: fulfilling promise with exception";
      queueEntry.promise->setException(runtime_error("QCli: Shutting down"));
    }
  } // drop queueEntry to release queuedParameters from its obtain.. function
  while (queuedParameters.use_count() >
         1) { // TODO cancel currently running future
    MLOG(MDEBUG) << "[" << id << "] "
                 << "~QCli sleeping";
    this_thread::sleep_for(chrono::seconds(1));
  }
  queuedParameters = nullptr;
  MLOG(MDEBUG) << "[" << id << "] "
               << "~QCli done";
}

SemiFuture<string> QueuedCli::executeRead(const ReadCommand cmd) {
  boost::recursive_mutex::scoped_lock scoped_lock(queuedParameters->mutex);
  MLOG(MDEBUG) << "[" << queuedParameters->id << "] "
               << "Invoking read command: \"" << cmd << "\"";
  return executeSomething(
             cmd,
             "QCli.executeRead",
             [params = queuedParameters, cmd]() {
               return params->cli->executeRead(cmd);
             })
      .via(queuedParameters->serialExecutorKeepAlive)
      .thenTry([params = queuedParameters, cmd](const Try<string>&& t) {
        // Logging callback
        if (t.hasException()) {
          MLOG(MWARNING) << "[" << params->id << "] "
                         << "Read command: \"" << cmd << "\""
                         << " Failed with: " << t.exception().what();
        } else {
          MLOG(MINFO) << "[" << params->id << "] "
                      << "Read command: \"" << cmd << "\""
                      << " Invoked successfully with output: \""
                      << Command::escape(t.value()) << "\"";
        }
        return t;
      })
      .semi();
}

SemiFuture<string> QueuedCli::executeWrite(const WriteCommand cmd) {
  boost::recursive_mutex::scoped_lock scoped_lock(queuedParameters->mutex);
  MLOG(MDEBUG) << "[" << queuedParameters->id << "] "
               << "Invoking write command: \"" << cmd << "\"";
  Command command = cmd;
  if (!command.isMultiCommand()) {
    // Single line config command, execute with read
    return executeRead(ReadCommand::create(cmd.raw(), true)); //skip cache
  }

  const vector<Command>& commands = command.splitMultiCommand();
  vector<Future<string>> commmandsFutures;

  for (unsigned long i = 0; i < (commands.size() - 1); i++) {
    commmandsFutures.emplace_back(
        executeSomething(
            commands.at(i),
            "QCli.executeWrite",
            [=]() {
              return queuedParameters->cli->executeWrite(
                  WriteCommand::create(commands.at(i)));
            })
            .via(queuedParameters->serialExecutorKeepAlive));
  }

  commmandsFutures.emplace_back(
      executeRead(ReadCommand::create(commands.back().raw(), true)) //skip cache
          .via(queuedParameters->serialExecutorKeepAlive));

  return reduce(
             commmandsFutures.begin(),
             commmandsFutures.end(),
             string(""),
             [](string s1, string s2) { return s1 + s2; })
      .thenTry([params = queuedParameters, cmd](const Try<string>&& t) {
        // Logging callback
        if (t.hasException()) {
          MLOG(MWARNING) << "[" << params->id << "] "
                         << "Write command: \"" << cmd << "\""
                         << " Failed with: " << t.exception().what();
        } else {
          MLOG(MINFO) << "[" << params->id << "] "
                      << "Write command: \"" << cmd << "\""
                      << " Invoked successfully with output: \""
                      << Command::escape(t.value()) << "\"";
        }
        return t;
      })
      .semi();
}

SemiFuture<string> QueuedCli::executeSomething(
    const Command& cmd,
    const string& prefix,
    function<SemiFuture<string>()> innerFunc) {
  shared_ptr<Promise<string>> promise = make_shared<Promise<string>>();
  QueueEntry queueEntry = QueueEntry{move(innerFunc), promise, cmd, prefix};
  MLOG(MDEBUG) << "[" << queuedParameters->id << "] (" << queueEntry.command
               << ") " << prefix << " adding to queue ('" << cmd << "')";

  queuedParameters->queue.enqueue(move(queueEntry));
  if (!queuedParameters->isProcessing) {
    triggerDequeue(queuedParameters);
  }
  return promise->getFuture().semi();
}

/*
 * Start queue reading on consumer thread if queue contains new items.
 * It is safe to call this method anytime, it is thread safe.
 */
void QueuedCli::triggerDequeue(shared_ptr<QueuedParameters> queuedParameters) {
  if (queuedParameters->shutdown) {
    return;
  }

  // switch to consumer thread
  via(queuedParameters->serialExecutorKeepAlive, [params = queuedParameters]() {
    MLOG(MDEBUG) << "[" << params->id << "] "
                 << "QCli.isProcessing:" << params->isProcessing
                 << ", queue size:" << params->queue.size();
    // do nothing if still waiting for remote device to respond
    if (params->isProcessing) {
      return;
    }
    QueueEntry queueEntry;
    if (!params->queue.try_dequeue(queueEntry)) {
      return;
    }

    if (params->shutdown) {
      MLOG(MDEBUG) << "[" << params->id << "] (" << queueEntry.command << ") "
                   << queueEntry.loggingPrefix << " Shutting down";
      queueEntry.promise->setException(runtime_error("QCli: Shutting down"));
      return;
    }

    params->isProcessing = true;
    SemiFuture<string> cliFuture = queueEntry.obtainFutureFromCli();
    MLOG(MDEBUG) << "[" << params->id << "] (" << queueEntry.command << ") "
                 << queueEntry.loggingPrefix
                 << " dequeued and cli future obtained";

    move(cliFuture)
        .via(params->serialExecutorKeepAlive)
        .then(
            params->serialExecutorKeepAlive,
            [params, queueEntry](string result) -> Future<Unit> {
              // after cliFuture completes, finish processing on consumer thread
              onDequeueSuccess(params, queueEntry, result);
              return unit;
            })
        .thenError([params, queueEntry](exception_wrapper e) -> Future<Unit> {
          onDequeueError(params, queueEntry, e);
          return unit;
        });
  });
}

void QueuedCli::onDequeueSuccess(
    const shared_ptr<QueuedParameters>& params,
    const QueuedCli::QueueEntry& queueEntry,
    const string& result) {
  MLOG(MDEBUG) << "[" << params->id << "] (" << queueEntry.command << ") "
               << queueEntry.loggingPrefix << " succeeded";
  params->isProcessing = false;
  queueEntry.promise->setValue(result);
  triggerDequeue(params);
}

void QueuedCli::onDequeueError(
    const shared_ptr<QueuedParameters>& params,
    const QueuedCli::QueueEntry& queueEntry,
    const exception_wrapper& e) {
  MLOG(MDEBUG) << "[" << params->id << "] (" << queueEntry.command << ") "
               << queueEntry.loggingPrefix << " failed with exception '"
               << e.what() << "'";
  if (!params->shutdown) {
    params->isProcessing = false;
  }
  queueEntry.promise->setException(e);
  triggerDequeue(params);
}

QueuedCli::QueuedParameters::QueuedParameters(
    const string& _id,
    const shared_ptr<Cli>& _cli,
    const shared_ptr<Executor>& _parentExecutor,
    const Executor::KeepAlive<SerialExecutor>& _serialExecutorKeepAlive)
    : id(_id),
      cli(_cli),
      parentExecutor(_parentExecutor),
      serialExecutorKeepAlive(_serialExecutorKeepAlive) {}
} // namespace cli
} // namespace channels
} // namespace devmand
