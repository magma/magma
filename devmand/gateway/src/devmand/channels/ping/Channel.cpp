// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <arpa/inet.h>
#include <folly/GLog.h>
#include <random>

#include <devmand/channels/ping/Channel.h>

namespace devmand {
namespace channels {
namespace ping {

Channel::Channel(Engine& engine_, folly::IPAddress target_)
    : engine(engine_), target(target_), sequence(genRandomRequestId()) {}

folly::Future<Rtt> Channel::ping() {
  auto pkt = IcmpPacket(target, getSequence());
  LOG(INFO) << "Sending ping to " << target.str() << " with sequence "
            << pkt.getSequence();
  return engine.ping(pkt);
}

RequestId Channel::genRandomRequestId() {
  std::random_device rd;
  std::mt19937 gen(rd());
  std::uniform_int_distribution<uint16_t> dis(0, UINT16_MAX);
  return dis(gen);
}

RequestId Channel::getSequence() {
  return ++sequence;
}

} // namespace ping
} // namespace channels
} // namespace devmand
