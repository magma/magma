// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <devmand/channels/Channel.h>
#include <devmand/channels/ping/Engine.h>
#include <gtest/gtest_prod.h>
namespace devmand {
  namespace test {
    class PingChannelTest_checkSequenceIdGeneration_Test;
  }
}

namespace devmand {
namespace channels {
namespace ping {

class Channel : public channels::Channel {
 public:
  Channel(Engine& engine, folly::IPAddress target_);
  Channel() = delete;
  ~Channel() override = default;
  Channel(const Channel&) = delete;
  Channel& operator=(const Channel&) = delete;
  Channel(Channel&&) = delete;
  Channel& operator=(Channel&&) = delete;

 public:
  folly::Future<Rtt> ping();

 private:
  friend devmand::test::PingChannelTest_checkSequenceIdGeneration_Test;
  // this was not working?
  // FRIEND_TEST(PingChannelTest, checkSequenceIdGeneration);
  RequestId getSequence();
  icmphdr makeIcmpPacket();
  RequestId genRandomRequestId();

 private:
  Engine& engine;
  folly::IPAddress target;
  RequestId sequence;
};

} // namespace ping
} // namespace channels
} // namespace devmand
