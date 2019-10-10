// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/Channel.h>
#include <devmand/channels/ping/Engine.h>

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
  RequestId getSequence();
  icmphdr makeIcmpPacket();

 private:
  Engine& engine;
  folly::IPAddress target;
  // TODO BOOTCAMP make this randomly initilized to minimize collisions.
  RequestId sequence{0};
};

} // namespace ping
} // namespace channels
} // namespace devmand
