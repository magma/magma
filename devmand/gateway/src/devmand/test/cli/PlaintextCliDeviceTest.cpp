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
#include <folly/json.h>
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
  Application app;
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

TEST_F(PlaintextCliDeviceTest, checkEcho) {
  cartography::DeviceConfig deviceConfig;
  devmand::cartography::ChannelConfig chnlCfg;
  map<string, string> kvPairs;
  kvPairs.insert(make_pair("stateCommand", "show interfaces brief"));
  chnlCfg.kvPairs = kvPairs;
  deviceConfig.channelConfigs.insert(make_pair("cli", chnlCfg));

  const shared_ptr<EchoCli>& echoCli = make_shared<EchoCli>();
  const shared_ptr<Channel>& channel = make_shared<Channel>("test", echoCli);
  unique_ptr<devices::Device> dev = make_unique<PlaintextCliDevice>(
      app, *cliEngine, deviceConfig.id, "show interfaces brief", channel);

  shared_ptr<State> state = dev->getState();
  const folly::dynamic& stateResult = state->collect().get();

  stringstream buffer;
  buffer << stateResult["show interfaces brief"];
  EXPECT_EQ("show interfaces brief", buffer.str());
}

static DeviceConfig getConfig(string port) {
  DeviceConfig deviceConfig;
  ChannelConfig chnlCfg;
  map<string, string> kvPairs;
  kvPairs.insert(make_pair("stateCommand", "echo 123"));
  kvPairs.insert(make_pair("port", port));
  kvPairs.insert(make_pair("username", "root"));
  kvPairs.insert(make_pair("password", "root"));
  kvPairs.insert(make_pair("keepAliveIntervalSeconds", "5"));
  kvPairs.insert(make_pair("maxCommandTimeoutSeconds", "60"));
  chnlCfg.kvPairs = kvPairs;
  deviceConfig.channelConfigs.insert(make_pair("cli", chnlCfg));
  deviceConfig.ip = "localhost";
  deviceConfig.id = "ubuntu-test-device";
  return deviceConfig;
}

TEST_F(PlaintextCliDeviceTest, plaintextCliDevicesError) {
  auto plaintextDevice = PlaintextCliDevice::createDeviceWithEngine(
      app, getConfig("9998"), *cliEngine);
  const shared_ptr<State>& ptr = plaintextDevice->getState();
  auto state = ptr->collect().get();

  // FIXME State class overrides status set from device to always say UP
  //  EXPECT_EQ(state["fbc-symphony-device:system"]["status"], "DOWN");
  EXPECT_EQ(state["fbc-symphony-device:errors"][0], "Not connected");
}

TEST_F(PlaintextCliDeviceTest, plaintextCliDevice) {
  unique_ptr<Device> dev = PlaintextCliDevice::createDeviceWithEngine(
      app, getConfig("9999"), *cliEngine);

  int i = 0;
  string output = "";
  do {
    if (i > 0) {
      this_thread::sleep_for(chrono::seconds(1));
    }

    i++;

    shared_ptr<State> state = dev->getState();
    auto t1 = chrono::high_resolution_clock::now();
    const folly::dynamic& stateResult = state->collect().get();
    auto t2 = chrono::high_resolution_clock::now();
    auto duration =
        chrono::duration_cast<chrono::microseconds>(t2 - t1).count();
    MLOG(MDEBUG) << "Retrieving state took: " << duration << " mu.";

    // FIXME instead of checking if there is output, check status once properly
    // reported by State class
    output = stateResult.getDefault("echo 123", "").asString();

    if (i > 20) {
      FAIL()
          << "Unable to execute command, probably not connected, last state: "
          << folly::toJson(stateResult);
    }
  } while (output == "");

  EXPECT_EQ(string("123"), boost::algorithm::trim_copy(output));
}

} // namespace cli
} // namespace test
} // namespace devmand
