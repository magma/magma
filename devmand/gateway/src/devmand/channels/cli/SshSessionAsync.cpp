// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <ErrorHandler.h>
#include <devmand/channels/cli/SshSession.h>
#include <devmand/channels/cli/SshSessionAsync.h>
#include <folly/executors/IOExecutor.h>
#include <folly/futures/Future.h>
#include <chrono>
#include <condition_variable>

namespace devmand {
namespace channels {
namespace cli {
namespace sshsession {

using boost::lockfree::spsc_queue;
using devmand::ErrorHandler;
using devmand::channels::cli::sshsession::SshSession;
using devmand::channels::cli::sshsession::SshSessionAsync;
using std::lock_guard;
using std::unique_lock;
using namespace std::chrono_literals;

SshSessionAsync::SshSessionAsync(
    string _id,
    shared_ptr<folly::Executor> _executor,
    shared_ptr<Timekeeper> _timekeeper)
    : id(_id),
      timekeeper(_timekeeper),
      executor(_executor),
      serialExecutor(SerialExecutor::create(
          folly::Executor::getKeepAliveToken(_executor.get()))),
      session(std::make_unique<SshSession>(_id)),
      readingMutex(),
      shutdown(false),
      callbackFinished(false),
      matchingExpectedOutput(false) {}

static const int EVENT_FINISH = 9999;

folly::SemiFuture<Unit> SshSessionAsync::destroy() {
  // idempotency
  if (shutdown) {
    return makeSemiFuture(folly::unit);
  }
  shutdown = true;
  MLOG(MDEBUG) << "[" << id << "] "
               << "destroy: started";

  // Let the NIO callback finish by injecting artificial event and waiting for
  // callback to finish
  Future<Unit> firstFuture;
  if (sessionEvent != nullptr && event_get_base(sessionEvent) != nullptr) {
    MLOG(MDEBUG) << "[" << id << "] "
                 << "destroy: setting EVENT_FINISH";
    event_active(this->sessionEvent, EVENT_FINISH, -1);
    firstFuture = waitForCallbacks();
  } else {
    firstFuture = makeFuture(folly::unit);
  }
  MLOG(MDEBUG) << "[" << id << "] "
               << "destroy: starting async cleanup";
  return std::move(firstFuture)
      .via(executor.get())
      .thenValue([this](auto) {
        MLOG(MDEBUG) << "[" << id << "] "
                     << "destroy: closing session";
        this->session->close();
      })
      .thenValue([this](auto) {
        MLOG(MDEBUG) << "[" << id << "] "
                     << "destroy: setReadingFlag";
        return this->setReadingFlag();
      })
      .thenValue([this](auto) {
        failCurrentRead(
            runtime_error("destroy:Session is closed"),
            this->readingState.promise);
        MLOG(MDEBUG) << "[" << id << "] "
                     << "destroy: finished";
      })
      .semi();
}

Future<Unit> SshSessionAsync::waitForCallbacks() {
  if (callbackFinished) {
    return makeFuture(folly::unit);
  }
  return via(executor.get())
      .delayed(100ms, timekeeper.get())
      .thenValue(
          [this](auto) -> Future<Unit> { return this->waitForCallbacks(); });
}

Future<Unit> SshSessionAsync::setReadingFlag() {
  // protect against race with matchExpectedOutput

  boost::try_mutex::scoped_try_lock lock(readingMutex);
  if (lock) {
    return makeFuture(folly::unit);
  }
  return via(executor.get())
      .delayed(100ms, timekeeper.get())
      .thenValue(
          [this](auto) -> Future<Unit> { return this->setReadingFlag(); });
}

SshSessionAsync::~SshSessionAsync() {
  MLOG(MDEBUG) << "[" << id << "] "
               << "~SshSessionAsync started";
  destroy().get();
  MLOG(MDEBUG) << "[" << id << "] "
               << "~SshSessionAsync finished";
}

void SshSessionAsync::unregisterEvent() {
  if (sessionEvent != nullptr && event_get_base(sessionEvent) != nullptr) {
    MLOG(MDEBUG) << "[" << id << "] "
                 << "SshSessionAsync unregisterEvent";
    event_free(sessionEvent);
  }
  callbackFinished = true;
}

Future<string> SshSessionAsync::read() {
  return via(serialExecutor.get(), [dis = shared_from_this()] {
    return dis->session->read();
  });
}

Future<Unit> SshSessionAsync::openShell(
    const string& ip,
    int port,
    const string& username,
    const string& password,
    const long sshConnectionTimeout) {
  return via(
      serialExecutor.get(),
      [dis = shared_from_this(),
       ip,
       port,
       username,
       password,
       sshConnectionTimeout] {
        dis->session->openShell(
            ip, port, username, password, sshConnectionTimeout);
      });
}

Future<Unit> SshSessionAsync::write(const string& command) {
  return via(serialExecutor.get(), [dis = shared_from_this(), command] {
    dis->session->write(command);
  });
}

Future<string> SshSessionAsync::readUntilOutput(const string& lastOutput) {
  this->readingState.currentLastOutput = lastOutput;
  this->readingState.promise = std::make_shared<Promise<string>>();
  this->readingState.outputSoFar = "";
  matchingExpectedOutput.store(true);
  processDataInBuffer(); // we could have had something already waiting in the
                         // queue
  return this->readingState.promise->getFuture();
}

void SshSessionAsync::setEvent(event* event) {
  this->sessionEvent = event;
}

void readCallback(evutil_socket_t fd, short what, void* ptr) {
  (void)fd;
  if (what == EVENT_FINISH) {
    ((SshSessionAsync*)ptr)->unregisterEvent();
    return;
  }
  ((SshSessionAsync*)ptr)->readSshDataToBuffer();
  ((SshSessionAsync*)ptr)->processDataInBuffer();
}

socket_t SshSessionAsync::getSshFd() {
  return this->session->getSshFd();
}

void SshSessionAsync::readSshDataToBuffer() {
  ErrorHandler::executeWithCatch([dis = shared_from_this()]() {
    const string& output = dis->session->read();
    if (!output.empty()) {
      dis->readQueue.push(output);
    }
  });
}

void SshSessionAsync::matchExpectedOutput() {
  if (not matchingExpectedOutput) { // we are not allowed to match unless
                                    // readUntilOutput is called because we
                                    // don't know against what to match
    return;
  }
  // protect against race with destruct
  boost::try_mutex::scoped_try_lock lock(readingMutex);
  if (lock) {
    if (shutdown) {
      return;
    }

    // If we are expecting empty string as output, we can complete right away
    // Why ? because keepalive is empty string
    if (this->readingState.currentLastOutput.empty()) {
      matchingExpectedOutput.store(false);
      this->readingState.promise->setValue("");
    }

    while (this->readQueue.read_available() != 0) {
      string output;
      readQueue.pop(output);
      this->readingState.outputSoFar.append(output);
      std::size_t found = this->readingState.outputSoFar.find(
          this->readingState.currentLastOutput);
      if (found != std::string::npos) {
        string final = this->readingState.outputSoFar.substr(0, found);

        // Check for any outstanding output
        size_t consumedLength =
            final.length() + this->readingState.currentLastOutput.length();
        if (consumedLength < this->readingState.outputSoFar.length()) {
          MLOG(MWARNING) << "[" << id << "] "
                         << "Unexpected output from device: ("
                         << this->readingState.outputSoFar.substr(
                                consumedLength)
                         << "). This output will be lost";
        }
        matchingExpectedOutput.store(false);
        this->readingState.promise->setValue(final);
      }
    }
  }
}

void SshSessionAsync::processDataInBuffer() {
  via(serialExecutor.get(),
      [dis = shared_from_this()] { dis->matchExpectedOutput(); });
}

void SshSessionAsync::failCurrentRead(
    runtime_error e,
    shared_ptr<Promise<string>> ptr) {
  if (ptr != nullptr and !ptr->isFulfilled()) {
    ptr->setException(folly::make_exception_wrapper<runtime_error>(e));
  }
}

} // namespace sshsession
} // namespace cli
} // namespace channels
} // namespace devmand
