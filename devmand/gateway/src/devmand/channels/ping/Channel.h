// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <map>
#include <string>

#include <netinet/ip_icmp.h>

#include <folly/futures/Future.h>
#include <folly/io/async/AsyncSocket.h>

#include <devmand/channels/Channel.h>
#include <devmand/utils/Time.h>

namespace devmand {
namespace channels {
namespace ping {

using Rtt = uint64_t;

struct Request {
  utils::TimePoint start;
  folly::Promise<Rtt> promise;
};

using RequestId = uint16_t;

using OutstandingRequests = std::map<RequestId, Request>;

class Channel : public channels::Channel, public folly::EventHandler {
 public:
  Channel(folly::EventBase& _eventBase, folly::IPAddress target_);
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

  virtual void handlerReady(uint16_t events) noexcept override;

 private:
  folly::EventBase& eventBase;
  OutstandingRequests outstandingRequests;
  folly::IPAddress target;
  int icmpSocket{-1};
  // TODO BOOTCAMP make this randomly initilized to minimize collisions.
  RequestId sequence{0};
};

} // namespace ping
} // namespace channels
} // namespace devmand
