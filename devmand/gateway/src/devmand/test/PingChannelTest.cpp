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
  folly::IPAddressV6 local3{"::1"};
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
  channels::ping::Engine engineV6(eventBase, channels::ping::IPVersion::v6);
  auto channel = std::make_shared<channels::ping::Channel>(engine, local);
  auto channel2 = std::make_shared<channels::ping::Channel>(engine, local2);
  auto channel3 = std::make_shared<channels::ping::Channel>(engineV6, local3);
  EXPECT_NE(0, channel->ping().get());
  EXPECT_NE(0, channel2->ping().get());
  EXPECT_NE(0, channel3->ping().get());
  EXPECT_NE(0, channel3->ping().get());
  EXPECT_NE(0, channel2->ping().get());
  EXPECT_NE(0, channel->ping().get());
  stop();
}

static bool isIncrementing(std::vector<uint16_t> nums, int increment) {
  for (size_t i = 1; i < nums.size(); i++) {
    uint16_t prev = nums[i - 1];
    uint16_t current = nums[i];
    if (current != prev + increment) {
      return false;
    }
  }
  return true;
}

static bool allSame(std::vector<uint16_t> nums) {
  if (nums.size() <= 1) {
    return true;
  }
  uint16_t first = nums[0];
  for (uint16_t n : nums) {
    if (first != n) {
      return false;
    }
  }
  return true;
}

TEST_F(PingChannelTest, checkSequenceIdGeneration) {
  std::vector<uint16_t> sequenceNumbers;
  static int NUM_GENERATIONS = 200;
  for (int i = 0; i < NUM_GENERATIONS; i++) {
    channels::ping::Engine engineV6(eventBase, channels::ping::IPVersion::v6);
    auto channel = std::make_shared<channels::ping::Channel>(engineV6, local3);
    uint16_t seqNum = channel->getSequence();
    sequenceNumbers.push_back(seqNum);
  }

  EXPECT_FALSE(allSame(sequenceNumbers));
  EXPECT_FALSE(isIncrementing(sequenceNumbers, 1));
  EXPECT_FALSE(isIncrementing(sequenceNumbers, -1));

  stop();
}

TEST_F(PingChannelTest, checkPingIpv6) {
  channels::ping::Engine engineV6(eventBase, channels::ping::IPVersion::v6);
  auto channel = std::make_shared<channels::ping::Channel>(engineV6, local3);
  EXPECT_NE(0, channel->ping().get());
  stop();
}

} // namespace test
} // namespace devmand
