// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <folly/json.h>
#include <gtest/gtest.h>

#include <devmand/devices/demo/Device.h>

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
  folly::dynamic data = devices::demo::Device::getDemoDatastore();
  std::cerr << folly::toJson(data) << std::endl;
}

} // namespace test
} // namespace devmand
