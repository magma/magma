// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/ping/Channel.h>
#include <devmand/test/EventBaseTest.h>

namespace devmand {
namespace test {

class PingChannelTest : public EventBaseTest {
 public:
  PingChannelTest() = default;
  ~PingChannelTest() override = default;
  PingChannelTest(const PingChannelTest&) = delete;
  PingChannelTest& operator=(const PingChannelTest&) = delete;
  PingChannelTest(PingChannelTest&&) = delete;
  PingChannelTest& operator=(PingChannelTest&&) = delete;

 protected:
  folly::IPAddress local{"127.0.0.1"};
  folly::IPAddress google{"127.0.0.2"};
};

TEST_F(PingChannelTest, checkPing) {
  channels::ping::Engine engine(eventBase);
  auto channel = std::make_shared<channels::ping::Channel>(engine, local);
  EXPECT_NE(0, channel->ping().get());
}

TEST_F(PingChannelTest, checkPingGoogle) {
  channels::ping::Engine engine(eventBase);
  auto channel = std::make_shared<channels::ping::Channel>(engine, google);
  EXPECT_NE(0, channel->ping().get());
}

TEST_F(PingChannelTest, checkMultiPing) {
  channels::ping::Engine engine(eventBase);
  auto channel = std::make_shared<channels::ping::Channel>(engine, local);
  auto channel2 = std::make_shared<channels::ping::Channel>(engine, google);
  EXPECT_NE(0, channel->ping().get());
  EXPECT_NE(0, channel2->ping().get());
  EXPECT_NE(0, channel2->ping().get());
  EXPECT_NE(0, channel->ping().get());
}

} // namespace test
} // namespace devmand
