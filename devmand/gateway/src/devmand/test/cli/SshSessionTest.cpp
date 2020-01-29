// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/channels/cli/SshSessionAsync.h>
#include <devmand/channels/cli/SshSocketReader.h>
#include <devmand/test/cli/utils/Log.h>
#include <devmand/test/cli/utils/Ssh.h>
#include <event2/thread.h>
#include <folly/Singleton.h>
#include <folly/executors/IOThreadPoolExecutor.h>
#include <folly/futures/Future.h>
#include <folly/futures/ThreadWheelTimekeeper.h>
#include <gtest/gtest.h>
#include <chrono>
#include <ctime>
#include <thread>

namespace devmand {
namespace test {
namespace cli {

using namespace devmand::channels::cli;
using namespace devmand::test::utils::ssh;
using namespace std;
using namespace folly;
using devmand::channels::cli::sshsession::readCallback;
using devmand::channels::cli::sshsession::SshSession;
using devmand::channels::cli::sshsession::SshSessionAsync;
using folly::IOThreadPoolExecutor;

static const shared_ptr<IOThreadPoolExecutor> executor =
    make_shared<IOThreadPoolExecutor>(2);

class SshSessionTest : public ::testing::Test {
 protected:
  shared_ptr<server> ssh;
  shared_ptr<folly::ThreadWheelTimekeeper> timekeeper =
      make_shared<ThreadWheelTimekeeper>();

  void SetUp() override {
    devmand::test::utils::log::initLog();
    devmand::test::utils::ssh::initSsh();
    ssh = startSshServer();
  }

  void TearDown() override {
    ssh->close();
  }
};

TEST_F(SshSessionTest, sessionStopReading) {
  atomic_bool exceptionCaught(false);
  {
    const shared_ptr<SshSessionAsync>& session =
        make_shared<SshSessionAsync>("testConn", executor, timekeeper);

    session->openShell("127.0.0.1", 9999, "cisco", "cisco").get();

    event* sessionEvent = SshSocketReader::getInstance().addSshReader(
        readCallback, session->getSshFd(), session.get());
    session->setEvent(sessionEvent);

    MLOG(MDEBUG) << "Starting infinite read";
    session->readUntilOutput("never").thenError(
        tag_t<runtime_error>{}, [&exceptionCaught](runtime_error const& e) {
          MLOG(MDEBUG) << "Read completed with error: " << e.what();
          exceptionCaught.store(true);
          return Future<string>(e);
        });
    this_thread::sleep_for(chrono::seconds(2));
  }

  // Assert the future completed exceptionally
  for (int i = 0; i < 10; i++) {
    if (exceptionCaught.load()) {
      return;
    } else {
      this_thread::sleep_for(chrono::seconds(1));
    }
  }
  FAIL();
}

TEST_F(SshSessionTest, sessionWriteThenDestruct) {
  evthread_use_pthreads();

  const shared_ptr<SshSessionAsync>& session =
      make_shared<SshSessionAsync>("testConn", executor, timekeeper);

  session->openShell("127.0.0.1", 9999, "cisco", "cisco").get();

  event* sessionEvent = SshSocketReader::getInstance().addSshReader(
      readCallback, session->getSshFd(), session.get());
  session->setEvent(sessionEvent);

  session->write("echo 1").get();
  auto output = session->readUntilOutput("1").get();

  MLOG(MDEBUG) << "Server responded with " << output;
  // Now make sure that session can cleanly disconnect
}

TEST_F(SshSessionTest, sessionStop) {
  const shared_ptr<SshSessionAsync>& session =
      make_shared<SshSessionAsync>("testConn", executor, timekeeper);

  session->openShell("127.0.0.1", 9999, "cisco", "cisco").get();

  event* sessionEvent = SshSocketReader::getInstance().addSshReader(
      readCallback, session->getSshFd(), session.get());
  session->setEvent(sessionEvent);

  // Now make sure that session can cleanly disconnect
}

TEST_F(SshSessionTest, empty) {
  // This tests the SSH server, to make sure it can cleanly start and close
}

} // namespace cli
} // namespace test
} // namespace devmand
