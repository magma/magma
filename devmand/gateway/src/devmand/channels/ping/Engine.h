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

#include <devmand/channels/Engine.h>
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

using OutstandingRequests =
    std::map<std::pair<folly::IPAddress, RequestId>, Request>;

class Engine : public channels::Engine, public folly::EventHandler {
 public:
  Engine(folly::EventBase& _eventBase);
  Engine() = delete;
  ~Engine() override = default;
  Engine(const Engine&) = delete;
  Engine& operator=(const Engine&) = delete;
  Engine(Engine&&) = delete;
  Engine& operator=(Engine&&) = delete;

 public:
  folly::Future<Rtt> ping(
      const icmphdr& hdr,
      const folly::IPAddress& destination);

 private:
  virtual void handlerReady(uint16_t events) noexcept override;

 private:
  folly::EventBase& eventBase;
  OutstandingRequests outstandingRequests;
  int icmpSocket{-1};
};

} // namespace ping
} // namespace channels
} // namespace devmand
