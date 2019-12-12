// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#define LOG_WITH_GLOG

#include <folly/executors/CPUThreadPoolExecutor.h>
#include <folly/futures/Future.h>
#include <libssh/callbacks.h>
#include <libssh/libssh.h>
#include <libssh/server.h>
#include <magma_logging.h>
#include <unistd.h>

namespace devmand {
namespace test {
namespace utils {
namespace ssh {

using namespace std;
using namespace folly;

extern atomic_bool sshInitialized;

extern void initSsh();

extern shared_ptr<CPUThreadPoolExecutor> testExecutor;

struct server {
 public:
  string id;
  ssh_bind sshbind = nullptr;
  ssh_session session = nullptr;
  Future<Unit> serverFuture;
  mutex received_guard;
  string received = "";

  bool isConnected() {
    return session != nullptr and ssh_is_connected(session);
  }

  string getReceived() {
    lock_guard<std::mutex> lg(received_guard);
    return received;
  }

  void close() {
    MLOG(MDEBUG) << "Closing server: " << id;
    if (session != nullptr) {
      shutdown(ssh_get_fd(session), SHUT_RDWR);
    }
    if (sshbind != nullptr) {
      shutdown(ssh_bind_get_fd(sshbind), SHUT_RDWR);
    }

    // Make sure the server thread finished before calling free
    move(serverFuture).wait();

    if (session != nullptr) {
      ssh_free(session);
      session = nullptr;
    }

    if (sshbind != nullptr) {
      ssh_bind_free(sshbind);
      sshbind = nullptr;
    }
  }
};

/**
 * Creates a disposable (use only once !) ssh server.
 * Blocks 1 executor thread.
 *
 * Can be closed by closing return value or by sending ^C.
 */
extern shared_ptr<server> startSshServer(
    shared_ptr<CPUThreadPoolExecutor> executor = testExecutor,
    string address = "0.0.0.0",
    uint port = 9999,
    string rsaKey = "/etc/ssh/ssh_host_rsa_key",
    string prompt = "PROMPT>");

} // namespace ssh
} // namespace utils
} // namespace test
} // namespace devmand
