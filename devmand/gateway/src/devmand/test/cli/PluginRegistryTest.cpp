// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/devices/cli/translation/PluginRegistry.h>
#include <devmand/test/cli/utils/Log.h>
#include <gtest/gtest.h>

namespace devmand {
namespace test {
namespace cli {

using namespace devmand::devices::cli;
using namespace folly;

class PluginRegistryTest : public ::testing::Test {
 protected:
  void SetUp() override {
    devmand::test::utils::log::initLog();
  }
};

TEST_F(PluginRegistryTest, api) {
  PluginRegistry reg;
  static const DeviceType type = {"IOS", "15.*"};

  typedef class : public Plugin {
   public:
    DeviceType getDeviceType() const override {
      return type;
    };

    void provideReaders(ReaderRegistryBuilder& reg) const override {
      (void)reg;
    };
    void provideWriters(WriterRegistryBuilder& reg) const override{
      (void)reg;
    };
  } TestPlugin;

  reg.registerPlugin(make_shared<TestPlugin>());
  MLOG(MDEBUG) << reg;

  EXPECT_THROW(
      reg.getDeviceContext({"NONEXISTING", "VERSION"}),
      PluginRegistryException);

  shared_ptr<DeviceContext> deviceCtx = reg.getDeviceContext(type);

  DeviceType reportedType = deviceCtx->getDeviceType();
  ASSERT_EQ(type, reportedType);
}

} // namespace cli
} // namespace test
} // namespace devmand
