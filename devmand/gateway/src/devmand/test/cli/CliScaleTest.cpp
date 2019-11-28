// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/channels/cli/PromptAwareCli.h>
#include <devmand/channels/cli/QueuedCli.h>
#include <devmand/channels/cli/SshSessionAsync.h>
#include <devmand/channels/cli/SshSocketReader.h>
#include <devmand/test/cli/utils/Log.h>
#include <folly/Singleton.h>
#include <folly/executors/IOThreadPoolExecutor.h>
#include <folly/futures/Future.h>
#include <gtest/gtest.h>
#include <libssh/callbacks.h>
#include <libssh/libssh.h>
#include <magma_logging.h>
#include <chrono>
#include <ctime>
#include <thread>

namespace devmand {
namespace test {
namespace cli {

using namespace devmand::channels::cli;
using namespace std;
using namespace folly;
using devmand::channels::cli::sshsession::readCallback;
using devmand::channels::cli::sshsession::SshSession;
using devmand::channels::cli::sshsession::SshSessionAsync;
using folly::IOThreadPoolExecutor;

class CliScaleTest : public ::testing::Test {
 protected:
  void SetUp() override {
    devmand::test::utils::log::initLog();
  }
};

static const shared_ptr<IOThreadPoolExecutor> executor =
    std::make_shared<IOThreadPoolExecutor>(8);
static const shared_ptr<IOThreadPoolExecutor> testExecutor =
    std::make_shared<IOThreadPoolExecutor>(20);

const static int DEVICES = 10;
const static int START_PORT = 10000;
const static int REQUESTS = 10;

/*
 * Run cli-testtool manually:
 *
 * git clone https://github.com/ncouture/MockSSH.git
 * cd MockSSH
 * Change file MockSSH.py by removing following 2 lines
 *
 *
         def sendKexInit(self):
         # Don't send key exchange prematurely
-        if not self.gotVersion:
-            return
         transport.SSHServerTransport.sendKexInit(self)
 *
 * Install this as a python package:
 * sudo python ./setup.py install
 *
 * git clone https://github.com/FRINXio/cli-testtool.git
 * cd cli-testtool
 * python mockdevice.py 0.0.0.0 10000 10 ssh devices/cisco_IOS.json
 *
 * This will run 10 IOS like SSH servers.
 */

TEST_F(CliScaleTest, DISABLED_scale) {
  ssh_threads_set_callbacks(ssh_threads_get_pthread());
  ssh_init();

  function<shared_ptr<QueuedCli>(int)> connect = [](int port) {
    const std::shared_ptr<SshSessionAsync>& session =
        std::make_shared<SshSessionAsync>(to_string(port), executor);

    session->openShell("172.8.0.100", port, "cisco", "cisco").get();

    shared_ptr<CliFlavour> cl = CliFlavour::create("ubiquiti");

    const shared_ptr<PromptAwareCli>& cli =
        std::make_shared<PromptAwareCli>(session, cl);

    cli->initializeCli();
    cli->resolvePrompt();
    event* sessionEvent = SshSocketReader::getInstance().addSshReader(
        readCallback, session->getSshFd(), session.get());
    MLOG(MWARNING) << "Connecting device at port" << port
                   << " at FD: " << session->getSshFd();
    session->setEvent(sessionEvent);
    return std::make_shared<QueuedCli>(to_string(port), cli, executor);
  };
  const Command& cmd = Command::makeReadCommand("show running-config");

  {
    std::chrono::steady_clock::time_point begin =
        std::chrono::steady_clock::now();
    vector<Future<shared_ptr<QueuedCli>>> connects;
    for (int cPort = START_PORT; cPort < DEVICES + START_PORT;
         cPort = cPort + 1) {
      MLOG(MWARNING) << "Connecting device at port: " << cPort;
      Future<shared_ptr<QueuedCli>> future =
          via(testExecutor.get(),
              [&connects, connect, cPort]() { return connect(cPort); });
      connects.push_back(std::move(future));
    }

    vector<shared_ptr<QueuedCli>> clis;
    for (uint i = 0; i < connects.size(); i++) {
      clis.push_back(std::move(connects.at(i)).get());
      MLOG(MWARNING) << "Connected device at port: " << START_PORT + i;
    }
    chrono::steady_clock::time_point end = chrono::steady_clock::now();

    MLOG(MWARNING) << "Connected devices count: " << clis.size();
    MLOG(MWARNING)
        << "Connected devices time: "
        << chrono::duration_cast<chrono::seconds>(end - begin).count();
    this_thread::sleep_for(chrono::seconds(10));

    vector<Future<string>> requests;
    for (uint device = 0; device < connects.size(); device++) {
      for (int req = 0; req < REQUESTS; req++) {
        MLOG(MWARNING) << "Invoking device: " << START_PORT + device
                       << " request: " << req;
        requests.push_back(move(clis[device]->executeAndRead(cmd)));
      }
    }

    vector<string> responses;
    for (uint i = 0; i < requests.size(); i++) {
      responses.push_back(std::move(requests.at(i)).get());
    }

    MLOG(MWARNING) << "Response count: " << responses.size();
  }

  MLOG(MWARNING) << "Now everything should be closed";
  ssh_finalize();
}

} // namespace cli
} // namespace test
} // namespace devmand
