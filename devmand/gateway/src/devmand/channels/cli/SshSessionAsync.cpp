// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <ErrorHandler.h>
#include <devmand/channels/cli/SshSession.h>
#include <devmand/channels/cli/SshSessionAsync.h>
#include <folly/executors/IOThreadPoolExecutor.h>
#include <folly/futures/Future.h>
#include <chrono>
#include <condition_variable>
#include <mutex>

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

SshSessionAsync::SshSessionAsync(string _id, shared_ptr<IOThreadPoolExecutor> _executor)
    : executor(_executor), session(_id), reading(false) {}

SshSessionAsync::~SshSessionAsync() {
  if (this->sessionEvent != nullptr &&
      event_get_base(this->sessionEvent) != nullptr) {
    event_free(this->sessionEvent);
  }
  this->session.close();

  while (reading.load()) {
    // waiting for any pending read to run out
  }
}

Future<string> SshSessionAsync::read(int timeoutMillis) {
  return via(executor.get(), [this, timeoutMillis] {
    return session.read(timeoutMillis);
  });
}

Future<Unit> SshSessionAsync::openShell(
    const string& ip,
    int port,
    const string& username,
    const string& password) {
  return via(executor.get(), [this, ip, port, username, password] {
    session.openShell(ip, port, username, password);
  });
}

Future<Unit> SshSessionAsync::write(const string& command) {
  return via(executor.get(), [this, command] { session.write(command); });
}

Future<string> SshSessionAsync::readUntilOutput(const string& lastOutput) {
  return via(executor.get(), [this, lastOutput] {
    return this->readUntilOutputBlocking(lastOutput);
  });
}

void SshSessionAsync::setEvent(event* event) {
  this->sessionEvent = event;
}

void readCallback(evutil_socket_t fd, short what, void* ptr) {
  (void)fd;
  (void)what;
  ((SshSessionAsync*)ptr)->readToBuffer();
}

socket_t SshSessionAsync::getSshFd() {
  return this->session.getSshFd();
}

void SshSessionAsync::readToBuffer() {
  {
    std::lock_guard<mutex> guard(mutex1);
    ErrorHandler::executeWithCatch([this]() {
      const string& output = this->session.read();
      if (!output.empty()) {
        readQueue.push(output);
      }
    });
  }
  condition.notify_one();
}

using namespace std::chrono_literals;

string SshSessionAsync::readUntilOutputBlocking(string lastOutput) {
  reading.store(true);
  string result;
  while (this->session.isOpen()) {
    unique_lock<mutex> lck(mutex1);
    bool satisfied = condition.wait_for(
        lck, 1000ms, [this] { return this->readQueue.read_available() != 0; });

    if (!satisfied) {
      continue;
    }

    string output;
    if (!readQueue.pop(output) ||
        output.empty()) { // sometimes we get the string "" or " " back, we can
                          // ignore that ...
      continue;
    }

    result.append(output);
    std::size_t found = result.find(lastOutput);
    if (found != std::string::npos) {
      // TODO check for any additional output after lastOutput
      reading.store(false);
      return result.substr(0, found);
    }
  }
  reading.store(false);
  throw std::runtime_error("Session is closed");
}

} // namespace sshsession
} // namespace cli
} // namespace channels
}