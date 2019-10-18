// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/snmp/Channel.h>
#include <devmand/channels/snmp/Engine.h>
#include <devmand/test/EventBaseTest.h>
#include <devmand/test/TestUtils.h>

#include <folly/IPAddress.h>

namespace devmand {
namespace test {

class SnmpChannelTest : public EventBaseTest {
 public:
  SnmpChannelTest() = default;
  ~SnmpChannelTest() override = default;
  SnmpChannelTest(const SnmpChannelTest&) = delete;
  SnmpChannelTest& operator=(const SnmpChannelTest&) = delete;
  SnmpChannelTest(SnmpChannelTest&&) = delete;
  SnmpChannelTest& operator=(SnmpChannelTest&&) = delete;

 protected:
  std::string local{"127.0.0.1"};
  std::string community{"public"};
  std::string version{"v1"};
};

TEST_F(SnmpChannelTest, checkSnmpTimeout) {
  channels::snmp::Engine engine(eventBase, "checkSnmpTimeout");
  auto channel = std::make_shared<channels::snmp::Channel>(
      engine, local, community, version);
  EXPECT_THROW(
      channel->asyncGet(channels::snmp::Oid("iso.3.6.1.2.1.1.4.0")).get(),
      std::runtime_error);
  stop();
}

} // namespace test
} // namespace devmand
