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
#include <devmand/devices/cli/PlaintextCliDevice.h>
#include <devmand/test/TestUtils.h>
#include <devmand/test/cli/utils/Log.h>
#include <devmand/test/cli/utils/MockCli.h>
#include <devmand/test/cli/utils/Ssh.h>
#include <gtest/gtest.h>

namespace devmand {
namespace test {
namespace cli {

using namespace devmand::channels::cli;
using namespace devmand::cartography;
using namespace devmand::devices;
using namespace devmand::devices::cli;
using namespace devmand::test::utils::cli;
using namespace std;
using namespace folly;
using namespace devmand::test::utils::ssh;
using devmand::channels::cli::IoConfigurationBuilder;
using devmand::channels::cli::ReadCachingCli;
using namespace std::chrono;

class CliTest : public ::testing::Test {
 protected:
  shared_ptr<server> ssh;
  unique_ptr<channels::cli::Engine> cliEngine;

  void SetUp() override {
    devmand::test::utils::log::initLog(MDEBUG);
    cliEngine = make_unique<channels::cli::Engine>();
    ssh = startSshServer();
  }

  void TearDown() override {
    ssh->close();
  }
};

static std::function<bool()> ensureConnected(const shared_ptr<Cli>& cli) {
  return [cli]() {
    try {
      cli->executeRead(ReadCommand::create("echo 123", true)).get();
        return true;
    } catch (const exception& e) {
        return false;
    }
  };
}

TEST_F(CliTest, writeMultipleTimesAllExecute) {
  cartography::DeviceConfig deviceConfig;
  devmand::cartography::ChannelConfig chnlCfg;
  std::map<std::string, std::string> kvPairs;
  kvPairs.insert(std::make_pair("port", "9999"));
  kvPairs.insert(std::make_pair("username", "root"));
  kvPairs.insert(std::make_pair("password", "root"));
  chnlCfg.kvPairs = kvPairs;
  deviceConfig.channelConfigs.insert(std::make_pair("cli", chnlCfg));
  deviceConfig.ip = "localhost";
  deviceConfig.id = "localhost-test-device";

  IoConfigurationBuilder ioConfigurationBuilder(deviceConfig, *cliEngine);
  const shared_ptr<Cli>& cli =
      ioConfigurationBuilder.createAll(ReadCachingCli::createCache());
  const function<bool()> &connectionTest = ensureConnected(cli);

  EXPECT_BECOMES_TRUE(connectionTest());

  cli->executeWrite(WriteCommand::create("sleep 2\nsleep2")).get();
  steady_clock::time_point begin = steady_clock::now();
  cli->executeWrite(WriteCommand::create("sleep 2")).get();
  steady_clock::time_point end = steady_clock::now();
  EXPECT_GE(duration_cast<milliseconds>(end - begin).count(), 2000);
}

} // namespace cli
} // namespace test
} // namespace devmand
