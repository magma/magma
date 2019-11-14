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
#include <devmand/devices/State.h>
#include <devmand/devices/cli/PlaintextCliDevice.h>
#include <devmand/test/cli/utils/Log.h>
#include <devmand/test/cli/utils/MockCli.h>
#include <devmand/test/cli/utils/Ssh.h>
#include <folly/executors/CPUThreadPoolExecutor.h>
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

class PlaintextCliDeviceTest : public ::testing::Test {
 protected:
  shared_ptr<server> ssh;

  void SetUp() override {
    devmand::test::utils::log::initLog();
    devmand::test::utils::ssh::initSsh();
    ssh = startSshServer();
  }

  void TearDown() override {
    ssh->close();
  }
};

TEST_F(PlaintextCliDeviceTest, checkEcho) {
  devmand::Application app;
  cartography::DeviceConfig deviceConfig;
  devmand::cartography::ChannelConfig chnlCfg;
  std::map<std::string, std::string> kvPairs;
  kvPairs.insert(std::make_pair("stateCommand", "show interfaces brief"));
  chnlCfg.kvPairs = kvPairs;
  deviceConfig.channelConfigs.insert(std::make_pair("cli", chnlCfg));

  const std::shared_ptr<EchoCli>& echoCli = std::make_shared<EchoCli>();
  const std::shared_ptr<Channel>& channel =
      std::make_shared<Channel>("test", echoCli);
  std::unique_ptr<devices::Device> dev = std::make_unique<PlaintextCliDevice>(
      app, deviceConfig.id, "show interfaces brief", channel);

  std::shared_ptr<State> state = dev->getState();
  const folly::dynamic& stateResult = state->collect().get();

  std::stringstream buffer;
  buffer << stateResult["show interfaces brief"];
  EXPECT_EQ("show interfaces brief", buffer.str());
}

static DeviceConfig getConfig(string port) {
  DeviceConfig deviceConfig;
  ChannelConfig chnlCfg;
  std::map<std::string, std::string> kvPairs;
  kvPairs.insert(std::make_pair("stateCommand", "echo 123"));
  kvPairs.insert(std::make_pair("port", port));
  kvPairs.insert(std::make_pair("username", "root"));
  kvPairs.insert(std::make_pair("password", "root"));
  chnlCfg.kvPairs = kvPairs;
  deviceConfig.channelConfigs.insert(std::make_pair("cli", chnlCfg));
  deviceConfig.ip = "localhost";
  deviceConfig.id = "ubuntu-test-device";
  return deviceConfig;
}

TEST_F(PlaintextCliDeviceTest, plaintextCliDevicesError) {
  Application app;
  EXPECT_ANY_THROW(PlaintextCliDevice::createDevice(app, getConfig("9998")));
}

TEST_F(PlaintextCliDeviceTest, plaintextCliDevice) {
  Application app;

  std::vector<std::unique_ptr<Device>> ds;
  for (int i = 0; i < 1; i++) {
    ds.push_back(
        std::move(PlaintextCliDevice::createDevice(app, getConfig("9999"))));
  }

  for (const auto& dev : ds) {
    std::shared_ptr<State> state = dev->getState();
    auto t1 = std::chrono::high_resolution_clock::now();
    const folly::dynamic& stateResult = state->collect().get();
    auto t2 = std::chrono::high_resolution_clock::now();
    auto duration =
        std::chrono::duration_cast<std::chrono::microseconds>(t2 - t1).count();
    MLOG(MDEBUG) << "Retrieving state took: " << duration << " mu.";
    std::stringstream buffer;

    buffer << stateResult["echo 123"];
    EXPECT_EQ("123", boost::algorithm::trim_copy(buffer.str()));
  }
}

} // namespace cli
} // namespace test
} // namespace devmand
