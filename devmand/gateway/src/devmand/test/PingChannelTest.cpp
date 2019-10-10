// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/ping/Channel.h>
#include <devmand/test/EventBaseTest.h>
#include <devmand/test/TestUtils.h>

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
  folly::IPAddress local2{"127.0.0.2"};
  folly::IPAddress dne{"203.0.113.0"};
};

TEST_F(PingChannelTest, checkPing) {
  channels::ping::Engine engine(eventBase);
  auto channel = std::make_shared<channels::ping::Channel>(engine, local);
  EXPECT_NE(0, channel->ping().get());
  stop();
}

TEST_F(PingChannelTest, checkPingTimeout) {
  channels::ping::Engine engine(
      eventBase,
      std::chrono::milliseconds(100),
      std::chrono::milliseconds(100));
  EXPECT_BECOMES_TRUE(eventBase.isRunning());
  engine.start();

  auto channel = std::make_shared<channels::ping::Channel>(engine, dne);
  auto channel2 = std::make_shared<channels::ping::Channel>(engine, local);
  folly::Future<channels::ping::Rtt> toFuture = channel->ping();

  while (not toFuture.isReady()) {
    EXPECT_NE(0, channel2->ping().get());
    std::chrono::milliseconds step(10);
    std::this_thread::sleep_for(step);
  }
  EXPECT_EQ(0, std::move(toFuture).get());
  stop();
}

TEST_F(PingChannelTest, checkMultiPing) {
  channels::ping::Engine engine(eventBase);
  auto channel = std::make_shared<channels::ping::Channel>(engine, local);
  auto channel2 = std::make_shared<channels::ping::Channel>(engine, local2);
  EXPECT_NE(0, channel->ping().get());
  EXPECT_NE(0, channel2->ping().get());
  EXPECT_NE(0, channel2->ping().get());
  EXPECT_NE(0, channel->ping().get());
  stop();
}

} // namespace test
} // namespace devmand
