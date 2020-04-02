// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <boost/algorithm/string/trim.hpp>
#include <devmand/Application.h>
#include <devmand/channels/cli/IoConfigurationBuilder.h>
#include <devmand/channels/cli/KeepaliveCli.h>
#include <devmand/channels/cli/TimeoutTrackingCli.h>
#include <devmand/devices/cli/PlaintextCliDevice.h>
#include <devmand/test/cli/utils/Log.h>
#include <devmand/test/cli/utils/Ssh.h>
#include <folly/Singleton.h>
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
using namespace devmand::cartography;
using namespace devmand::devices;
using namespace devmand::devices::cli;

using namespace chrono_literals;

class ReconnectingSshTest : public ::testing::Test {
 protected:
  shared_ptr<server> ssh;
  unique_ptr<channels::cli::Engine> cliEngine;

  void SetUp() override {
    devmand::test::utils::log::initLog();
    ssh = startSshServer();
    cliEngine = make_unique<channels::cli::Engine>(dynamic::object());
  }

  void TearDown() override {
    ssh->close();
  }
};

static DeviceConfig getConfig(
    string port,
    chrono::seconds commandTimeout = defaultCommandTimeout,
    chrono::seconds keepaliveTimeout = 5s) {
  DeviceConfig deviceConfig;
  ChannelConfig chnlCfg;
  map<string, string> kvPairs;
  kvPairs.insert(make_pair("stateCommand", "echo 123"));
  kvPairs.insert(make_pair("port", port));
  kvPairs.insert(make_pair("username", "root"));
  kvPairs.insert(make_pair("password", "root"));
  kvPairs.insert(make_pair(
      configMaxCommandTimeoutSeconds, to_string(commandTimeout.count())));
  kvPairs.insert(make_pair(
      configKeepAliveIntervalSeconds, to_string(keepaliveTimeout.count())));
  chnlCfg.kvPairs = kvPairs;
  deviceConfig.channelConfigs.insert(make_pair("cli", chnlCfg));
  deviceConfig.ip = "localhost";
  deviceConfig.id = "ubuntu-test-device";
  return deviceConfig;
}

static void ensureConnected(const shared_ptr<Cli>& cli) {
  bool connected = false;
  int attempts = 0;
  while (!connected && attempts++ < 40) {
    MLOG(MDEBUG) << "Testing connection attempt:" << attempts;
    try {
      const string& echoResult =
          cli->executeRead(ReadCommand::create("echo 123", true)).get();
      EXPECT_EQ("123", boost::algorithm::trim_copy(echoResult));
      connected = true;
    } catch (const exception& e) {
      MLOG(MDEBUG) << "Not connected:" << e.what();
      this_thread::sleep_for(500ms);
    }
  }
  EXPECT_TRUE(connected);
}

TEST_F(ReconnectingSshTest, commandTimeout) {
  int cmdTimeout = 5;
  IoConfigurationBuilder ioConfigurationBuilder(
      getConfig(
          "9999", std::chrono::seconds(cmdTimeout), std::chrono::seconds(10)),
      *cliEngine,
      CliFlavour::getDefaultInstance());
  shared_ptr<Cli> cli = ioConfigurationBuilder.createAll(
      ReadCachingCli::createCache(),
      make_shared<TreeCache>(
          ioConfigurationBuilder.getConnectionParameters()->flavour));
  ensureConnected(cli);
  // sleep so that cli stack will be destroyed
  string sleepCommand = "sleep ";
  sleepCommand.append(to_string(cmdTimeout + 1));
  EXPECT_THROW(
      { cli->executeRead(ReadCommand::create(sleepCommand, true)).get(); },
      CommandTimeoutException);

  ssh->close();
  ssh = startSshServer();
  ensureConnected(cli);
}

TEST_F(ReconnectingSshTest, serverDisconnectSendCommands) {
  int cmdTimeout = 60;
  IoConfigurationBuilder ioConfigurationBuilder(
      getConfig(
          "9999", std::chrono::seconds(cmdTimeout), std::chrono::seconds(60)),
      *cliEngine,
      CliFlavour::getDefaultInstance());
  shared_ptr<Cli> cli = ioConfigurationBuilder.createAll(
      ReadCachingCli::createCache(),
      make_shared<TreeCache>(
          ioConfigurationBuilder.getConnectionParameters()->flavour));

  ensureConnected(cli);
  ssh->close();
  ssh = startSshServer();
  ensureConnected(cli);
}

TEST_F(ReconnectingSshTest, serverDisconnectWaitForKeepalive) {
  int cmdTimeout = 5;
  int keepaliveFreq = 5;
  IoConfigurationBuilder ioConfigurationBuilder(
      getConfig(
          "9999",
          std::chrono::seconds(cmdTimeout),
          std::chrono::seconds(keepaliveFreq)),
      *cliEngine,
      CliFlavour::getDefaultInstance());
  shared_ptr<Cli> cli = ioConfigurationBuilder.createAll(
      ReadCachingCli::createCache(),
      make_shared<TreeCache>(
          ioConfigurationBuilder.getConnectionParameters()->flavour));

  ensureConnected(cli);

  // Disconnect from server side
  ssh->close();
  ssh = startSshServer();

  MLOG(MDEBUG) << "Waiting for CLI to reconnect";

  int attempt = 0;
  while (true) {
    if (ssh->isConnected()) {
      break;
    }
    if ((attempt++) == 30) {
      FAIL() << "CLI did not reconnect, something went wrong";
    }

    this_thread::sleep_for(1s);
  }
}

TEST_F(ReconnectingSshTest, keepalive) {
  int cmdTimeout = 5;
  int keepaliveTimeout = 5;
  IoConfigurationBuilder ioConfigurationBuilder(
      getConfig(
          "9999",
          std::chrono::seconds(cmdTimeout),
          std::chrono::seconds(keepaliveTimeout)),
      *cliEngine,
      CliFlavour::getDefaultInstance());
  shared_ptr<Cli> cli = ioConfigurationBuilder.createAll(
      ReadCachingCli::createCache(),
      make_shared<TreeCache>(
          ioConfigurationBuilder.getConnectionParameters()->flavour));

  int attempt = 0;
  while (true) {
    auto receivedOnServer = ssh->getReceived();
    // Make sure at least 4 newlines have been executed
    // 2 for prompt resolution and 2 for keepalives
    if (count(receivedOnServer.begin(), receivedOnServer.end(), '\n') > 3) {
      break;
    }
    if ((attempt++) == 30) {
      FAIL()
          << "Keepalive did not occur, something went wrong. Server received: "
          << receivedOnServer;
    }

    this_thread::sleep_for(1s);
  }
}

} // namespace cli
} // namespace test
} // namespace devmand
