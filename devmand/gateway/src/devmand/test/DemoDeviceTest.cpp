// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

#include <folly/json.h>
#include <gtest/gtest.h>

#include <devmand/devices/DemoDevice.h>

namespace devmand {
namespace test {

class DemoDeviceTest : public ::testing::Test {
 public:
  DemoDeviceTest() = default;
  virtual ~DemoDeviceTest() = default;
  DemoDeviceTest(const DemoDeviceTest&) = delete;
  DemoDeviceTest& operator=(const DemoDeviceTest&) = delete;
  DemoDeviceTest(DemoDeviceTest&&) = delete;
  DemoDeviceTest& operator=(DemoDeviceTest&&) = delete;
};

TEST_F(DemoDeviceTest, jsonSample) {
  folly::dynamic data = devices::DemoDevice::getDemoState();
  std::cerr << folly::toJson(data) << std::endl;
}

} // namespace test
} // namespace devmand
