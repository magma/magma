// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/channels/cli/IoConfigurationBuilder.h>
#include <devmand/channels/cli/SshSessionAsync.h>
#include <devmand/test/TestUtils.h>
#include <devmand/test/cli/utils/Log.h>
#include <folly/Singleton.h>
#include <folly/futures/Future.h>
#include <gtest/gtest.h>
#include <libssh/libssh.h>
#include <chrono>
#include <ctime>
#include <thread>

namespace devmand {
namespace test {
namespace cli {

using namespace devmand::channels::cli;
using namespace std;
using namespace folly;

class CliScaleTest : public ::testing::Test {
 protected:
  channels::cli::Engine cliEngine{dynamic::object()};

  void SetUp() override {
    devmand::test::utils::log::initLog(MWARNING);
    //    cliEngine = channels::cli::Engine();
  }
};

const static int DEVICES = 20;
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
  function<shared_ptr<Cli>(int)> connect = [this](int port) {
    IoConfigurationBuilder ioConfigurationBuilder(
        IoConfigurationBuilder::makeConnectionParameters(
            to_string(port),
            "172.8.0.1",
            "cisco",
            "cisco",
            "ubiquiti",
            port,
            10s,
            60s,
            10s,
            30,
            cliEngine));
    return ioConfigurationBuilder.createAll(
        ReadCachingCli::createCache(),
        make_shared<TreeCache>(
            ioConfigurationBuilder.getConnectionParameters()->flavour));
  };
  const ReadCommand& cmd = ReadCommand::create("show running-config");

  {
    chrono::steady_clock::time_point begin = chrono::steady_clock::now();
    vector<Future<shared_ptr<Cli>>> connects;
    for (int cPort = START_PORT; cPort < DEVICES + START_PORT;
         cPort = cPort + 1) {
      MLOG(MWARNING) << "Connecting device at port: " << cPort;
      connects.emplace_back(connect(cPort));
    }

    vector<shared_ptr<Cli>> clis;
    for (uint i = 0; i < connects.size(); i++) {
      clis.push_back(move(connects.at(i)).get());

      // Wait till connected
      EXPECT_BECOMES_TRUE(not clis.at(clis.size() - 1)
                                  ->executeRead(ReadCommand::create("", true))
                                  .getTry()
                                  .hasException());

      MLOG(MWARNING) << "Connected device at port: " << START_PORT + i;
    }
    chrono::steady_clock::time_point end = chrono::steady_clock::now();

    MLOG(MWARNING) << "Connected devices count: " << clis.size();
    MLOG(MWARNING)
        << "Connected devices time: "
        << chrono::duration_cast<chrono::seconds>(end - begin).count();

    this_thread::sleep_for(chrono::seconds(10));

    begin = chrono::steady_clock::now();
    vector<SemiFuture<string>> requests;
    for (uint device = 0; device < connects.size(); device++) {
      for (int req = 0; req < REQUESTS; req++) {
        MLOG(MWARNING) << "Invoking device: " << START_PORT + device
                       << " request: " << req;
        requests.push_back(move(clis[device]->executeRead(cmd)));
      }
    }

    vector<string> responses;
    for (uint i = 0; i < requests.size(); i++) {
      responses.push_back(move(requests.at(i)).get());
    }
    end = chrono::steady_clock::now();

    MLOG(MWARNING) << "Response count: " << responses.size();
    MLOG(MWARNING)
        << "Request execution time: "
        << chrono::duration_cast<chrono::seconds>(end - begin).count();
    exit(1);
  }
  MLOG(MWARNING) << "Closed";
}
} // namespace cli
} // namespace test
} // namespace devmand
